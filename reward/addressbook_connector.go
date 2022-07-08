// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package reward

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/reward/contract"
	"github.com/klaytn/klaytn/params"
)

// addressType defined in AddressBook
const (
	addressTypeNodeID = iota
	addressTypeStakingAddr
	addressTypeRewardAddr
	addressTypePoCAddr
	addressTypeKIRAddr
)

var errAddressBookIncomplete = errors.New("incomplete node information from AddressBook")

type addressBookConnector struct {
	bc              blockChain
	gh              governanceHelper
	abi             string
	contractAddress common.Address
}

// create and return addressBookConnector
func newAddressBookConnector(bc blockChain, gh governanceHelper) *addressBookConnector {
	return &addressBookConnector{
		bc:              bc,
		gh:              gh,
		abi:             contract.AddressBookABI,
		contractAddress: common.HexToAddress(contract.AddressBookContractAddress),
	}
}

// make a message to the addressBook contract for executing getAllAddress function of the addressBook contract
func (ac *addressBookConnector) makeMsgToAddressBook(r params.Rules) (*types.Transaction, error) {
	abiInstance, err := abi.JSON(strings.NewReader(ac.abi))
	if err != nil {
		return nil, err
	}

	data, err := abiInstance.Pack("getAllAddress")
	if err != nil {
		return nil, err
	}

	intrinsicGas, err := types.IntrinsicGas(data, nil, false, r)
	if err != nil {
		return nil, err
	}

	// Create new call message
	// TODO-Klaytn-Issue1166 Decide who will be sender(i.e. from)
	msg := types.NewMessage(common.Address{}, &ac.contractAddress, 0, big.NewInt(0), 10000000, big.NewInt(0), data, false, intrinsicGas)

	return msg, nil
}

// It parses the result bytes of calling addressBook to addresses.
func (ac *addressBookConnector) parseAllAddresses(result []byte) (nodeIds []common.Address, stakingAddrs []common.Address, rewardAddrs []common.Address, pocAddr common.Address, kirAddr common.Address, err error) {
	nodeIds = []common.Address{}
	stakingAddrs = []common.Address{}
	rewardAddrs = []common.Address{}
	pocAddr = common.Address{}
	kirAddr = common.Address{}

	if result == nil {
		err = errAddressBookIncomplete
		return
	}

	abiInstance, err := abi.JSON(strings.NewReader(ac.abi))
	if err != nil {
		return
	}

	var (
		allTypeList    = new([]uint8)
		allAddressList = new([]common.Address)
	)
	out := &[]interface{}{
		allTypeList,
		allAddressList,
	}

	err = abiInstance.Unpack(out, "getAllAddress", result)
	if err != nil {
		logger.Trace("abiInstance.Unpack failed for getAllAddress in AddressBook")
		return
	}

	if len(*allTypeList) != len(*allAddressList) {
		err = errors.New(fmt.Sprintf("length of type list and address list differ. len(type)=%d, len(addrs)=%d", len(*allTypeList), len(*allAddressList)))
		return
	}

	// Parse and construct node information
	for i, addrType := range *allTypeList {
		switch addrType {
		case addressTypeNodeID:
			nodeIds = append(nodeIds, (*allAddressList)[i])
		case addressTypeStakingAddr:
			stakingAddrs = append(stakingAddrs, (*allAddressList)[i])
		case addressTypeRewardAddr:
			rewardAddrs = append(rewardAddrs, (*allAddressList)[i])
		case addressTypePoCAddr:
			pocAddr = (*allAddressList)[i]
		case addressTypeKIRAddr:
			kirAddr = (*allAddressList)[i]
		default:
			err = errors.New(fmt.Sprintf("invalid type from AddressBook: %d", addrType))
			return
		}
	}

	// validate parsed node information
	if len(nodeIds) != len(stakingAddrs) ||
		len(nodeIds) != len(rewardAddrs) ||
		common.EmptyAddress(pocAddr) ||
		common.EmptyAddress(kirAddr) {
		err = errAddressBookIncomplete
		return
	}

	return
}

// getStakingInfoFromAddressBook returns stakingInfo when calling AddressBook succeeded.
// If addressBook is not activated, emptyStakingInfo is returned.
// After addressBook is activated, it returns stakingInfo with addresses and stakingAmount.
// Otherwise, it returns an error.
func (ac *addressBookConnector) getStakingInfoFromAddressBook(blockNum uint64) (*StakingInfo, error) {
	if !params.IsStakingUpdateInterval(blockNum) {
		return nil, errors.New(fmt.Sprintf("not staking block number. blockNum: %d", blockNum))
	}

	// Prepare a message
	msg, err := ac.makeMsgToAddressBook(ac.bc.Config().Rules(new(big.Int).SetUint64(blockNum)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to make message for AddressBook. root err: %s", err))
	}

	intervalBlock := ac.bc.GetBlockByNumber(blockNum)
	if intervalBlock == nil {
		return nil, errors.New("stateDB is not ready for staking info")
	}
	statedb, err := ac.bc.StateAt(intervalBlock.Root())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to make a state for interval block. blockNum: %d, root err: %s", blockNum, err))
	}

	// Create a new context to be used in the EVM environment
	context := blockchain.NewEVMContext(msg, intervalBlock.Header(), ac.bc, nil)
	evm := vm.NewEVM(context, statedb, ac.bc.Config(), &vm.Config{})

	res, gas, kerr := blockchain.ApplyMessage(evm, msg)
	logger.Trace("Call AddressBook contract", "used gas", gas, "kerr", kerr)
	err = kerr.ErrTxInvalid
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to call AddressBook contract. root err: %s", err))
	}

	nodeAddrs, stakingAddrs, rewardAddrs, PoCAddr, KIRAddr, err := ac.parseAllAddresses(res)
	if err != nil {
		if err == errAddressBookIncomplete {
			// This is an expected behavior when the addressBook contract is not activated yet.
			logger.Info("The addressBook is not yet activated. Use empty stakingInfo", "reason", err)
		} else {
			logger.Error("Fail while parsing a result from the addressBook. Use empty staking info", "err", err)
		}
		return newEmptyStakingInfo(blockNum), nil
	}

	return newStakingInfo(ac.bc, ac.gh, blockNum, nodeAddrs, stakingAddrs, rewardAddrs, KIRAddr, PoCAddr)
}
