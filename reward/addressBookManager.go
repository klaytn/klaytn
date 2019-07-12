package reward

import (
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/reward/contract"
	"github.com/klaytn/klaytn/params"
	"math/big"
	"strings"
)

// addressType defined in AddressBook
const (
	addressTypeNodeID = iota
	addressTypeStakingAddr
	addressTypeRewardAddr
	addressTypePoCAddr
	addressTypeKIRAddr
)

var (
	errAddressBookIncomplete = errors.New("incomplete node information from AddressBook")
)

type addressBookManager struct {
	bc                         *blockchain.BlockChain
	chainConfig                *params.ChainConfig
	governanceHelper           governanceHelper
	addressBookABI             string
	addressBookContractAddress common.Address
}

func newAddressBookManager(bc *blockchain.BlockChain, governanceHelper governanceHelper) *addressBookManager {
	return &addressBookManager{
		bc:                         bc,
		chainConfig:                bc.Config(),
		governanceHelper:           governanceHelper,
		addressBookABI:             contract.AddressBookABI,
		addressBookContractAddress: common.HexToAddress(contract.AddressBookContractAddress),
	}
}

func (abm *addressBookManager) makeMsgToAddressBook() (*types.Transaction, error) {
	abii, err := abi.JSON(strings.NewReader(abm.addressBookABI))
	if err != nil {
		return nil, err
	}

	data, err := abii.Pack("getAllAddress")
	if err != nil {
		return nil, err
	}

	intrinsicGas, err := types.IntrinsicGas(data, false, true)
	if err != nil {
		return nil, err
	}

	// Create new call message
	// TODO-Klaytn-Issue1166 Decide who will be sender(i.e. from)
	msg := types.NewMessage(common.Address{}, &abm.addressBookContractAddress, 0, big.NewInt(0), 10000000, big.NewInt(0), data, false, intrinsicGas)

	return msg, nil
}

func (abm *addressBookManager) getAllAddressFromAddressBook(result []byte) ([]common.Address, []common.Address, []common.Address, common.Address, common.Address, error) {
	if result == nil {
		return nil, nil, nil, common.Address{}, common.Address{}, errAddressBookIncomplete
	}

	abii, err := abi.JSON(strings.NewReader(abm.addressBookABI))
	if err != nil {
		return nil, nil, nil, common.Address{}, common.Address{}, err
	}

	var (
		allTypeList    = new([]uint8)
		allAddressList = new([]common.Address)
	)
	out := &[]interface{}{
		allTypeList,
		allAddressList,
	}

	err = abii.Unpack(out, "getAllAddress", result)
	if err != nil {
		logger.Trace("abii.Unpack failed for getAllAddress in AddressBook")
		return nil, nil, nil, common.Address{}, common.Address{}, err
	}

	if len(*allTypeList) != len(*allAddressList) {
		return nil, nil, nil, common.Address{}, common.Address{}, errors.New(fmt.Sprintf("length of type list and address list differ. len(type)=%d, len(addrs)=%d", len(*allTypeList), len(*allAddressList)))
	}

	// Parse and construct node information
	nodeIds := []common.Address{}
	stakingAddrs := []common.Address{}
	rewardAddrs := []common.Address{}
	pocAddr := common.Address{}
	kirAddr := common.Address{}
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
			return nil, nil, nil, common.Address{}, common.Address{}, errors.New(fmt.Sprintf("invalid type from AddressBook: %d", addrType))
		}
	}

	// validate parsed node information
	if len(nodeIds) != len(stakingAddrs) ||
		len(nodeIds) != len(rewardAddrs) ||
		isEmptyAddress(pocAddr) ||
		isEmptyAddress(kirAddr) {
		// This is expected behavior when bootstrapping
		//logger.Trace("Incomplete node information from AddressBook.",
		//	"# of nodeIds", len(nodeIds),
		//	"# of stakingAddrs", len(stakingAddrs),
		//	"# of rewardAddrs", len(rewardAddrs),
		//	"PoC address", pocAddr.String(),
		//	"KIR address", kirAddr.String())

		return nil, nil, nil, common.Address{}, common.Address{}, errAddressBookIncomplete
	}

	return nodeIds, stakingAddrs, rewardAddrs, pocAddr, kirAddr, nil
}

// getStakingInfoFromAddressBook returns staking info when calling AddressBook
// succeeded. It returns an error otherwise.
func (abm *addressBookManager) getStakingInfoFromAddressBook(blockNum uint64) (*StakingInfo, error) {
	var nodeIds []common.Address
	var stakingAddrs []common.Address
	var rewardAddrs []common.Address
	var KIRAddr = common.Address{}
	var PoCAddr = common.Address{}
	var err error

	if !params.IsStakingUpdateInterval(blockNum) {
		return nil, errors.New(fmt.Sprintf("not staking block number. blockNum: %d", blockNum))
	}

	// Prepare a message
	msg, err := abm.makeMsgToAddressBook()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to make message for AddressBook. root err: %s", err))
	}

	intervalBlock := abm.bc.GetBlockByNumber(blockNum)
	if intervalBlock == nil {
		return nil, errors.New("stateDB is not ready for staking info")
	}
	statedb, err := abm.bc.StateAt(intervalBlock.Root())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to make a state for interval block. blockNum: %d, root err: %s", blockNum, err))
	}

	// Create a new context to be used in the EVM environment
	context := blockchain.NewEVMContext(msg, intervalBlock.Header(), abm.bc, nil)
	evm := vm.NewEVM(context, statedb, abm.chainConfig, &vm.Config{})

	res, gas, kerr := blockchain.ApplyMessage(evm, msg)
	logger.Trace("Call AddressBook contract", "used gas", gas, "kerr", kerr)
	err = kerr.ErrTxInvalid
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to call AddressBook contract. root err: %s", err))
	}

	nodeIds, stakingAddrs, rewardAddrs, PoCAddr, KIRAddr, err = abm.getAllAddressFromAddressBook(res)
	if err != nil {
		if err == errAddressBookIncomplete {
			// This is expected behavior when smart contract is not setup yet.
			logger.Info("Use empty staking info instead of info from AddressBook", "reason", err)
		} else {
			logger.Error("Failed to parse result from AddressBook contract. Use empty staking info", "err", err)
		}
		return newEmptyStakingInfo(blockNum)
	}

	return newStakingInfo(abm.bc, abm.governanceHelper, blockNum, nodeIds, stakingAddrs, rewardAddrs, KIRAddr, PoCAddr)
}
