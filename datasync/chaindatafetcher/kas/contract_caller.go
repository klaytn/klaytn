// Copyright 2020 The klaytn Authors
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

package kas

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/vm"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/contracts/kip13"
	"github.com/klaytn/klaytn/networks/rpc"
)

// TODO-ChainDataFetcher extract the call timeout c as a configuration
const callTimeout = 300 * time.Millisecond

var (
	// KIP 13: Interface Query Standard - https://kips.klaytn.com/KIPs/kip-13
	IKIP13Id  = [4]byte{0x01, 0xff, 0xc9, 0xa7}
	InvalidId = [4]byte{0xff, 0xff, 0xff, 0xff}

	// KIP 7: Fungible Token Standard - https://kips.klaytn.com/KIPs/kip-7
	IKIP7Id         = [4]byte{0x65, 0x78, 0x73, 0x71}
	IKIP7MetadataId = [4]byte{0xa2, 0x19, 0xa0, 0x25}

	// KIP 17: Non-fungible Token Standard - https://kips.klaytn.com/KIPs/kip-17
	IKIP17Id         = [4]byte{0x80, 0xac, 0x58, 0xcd}
	IKIP17MetadataId = [4]byte{0x5b, 0x5e, 0x13, 0x9f}

	errMsgEmptyOutput = "abi: unmarshalling empty output"
)

//go:generate mockgen -destination=./mocks/blockchain_api_mock.go -package=mocks github.com/klaytn/klaytn/datasync/chaindatafetcher/kas BlockchainAPI
// BlockchainAPI interface is for testing purpose.
type BlockchainAPI interface {
	GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error)
	Call(ctx context.Context, args api.CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error)
}

// contractCaller performs kip13 method `supportsInterface` to detect the deployed contracts are KIP7 or KIP17.
type contractCaller struct {
	blockchainAPI BlockchainAPI
	callTimeout   time.Duration
}

func newContractCaller(api BlockchainAPI) *contractCaller {
	return &contractCaller{blockchainAPI: api, callTimeout: callTimeout}
}

func (f *contractCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	num := rpc.LatestBlockNumber
	if blockNumber != nil {
		num = rpc.BlockNumber(blockNumber.Int64())
	}
	return f.blockchainAPI.GetCode(ctx, contract, rpc.NewBlockNumberOrHashWithNumber(num))
}

func (f *contractCaller) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	num := rpc.LatestBlockNumber
	if blockNumber != nil {
		num = rpc.BlockNumber(blockNumber.Int64())
	}
	callArgs := api.CallArgs{
		From: call.From,
		To:   call.To,
		Data: hexutil.Bytes(call.Data),
	}
	return f.blockchainAPI.Call(ctx, callArgs, rpc.NewBlockNumberOrHashWithNumber(num))
}

func getCallOpts(blockNumber *big.Int, timeout time.Duration) (*bind.CallOpts, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &bind.CallOpts{Context: ctx, BlockNumber: blockNumber}, cancel
}

// supportsInterface returns true if the given interfaceID is supported, otherwise returns false.
func (f *contractCaller) supportsInterface(contract common.Address, opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	caller, err := kip13.NewInterfaceIdentifierCaller(contract, f)
	if err != nil {
		logger.Error("NewInterfaceIdentifierCaller is failed", "contract", contract.String(), "interfaceID", hexutil.Encode(interfaceID[:]))
		return false, err
	}
	isSupported, err := caller.SupportsInterface(opts, interfaceID)
	if err != nil {
		// removed handle error case for SupportsInterface contract call.
		// 1. The error cases are too many to be handled.
		// 2. There is no case to return an error (e.g. network error, ...) other than internal errors.
		if !strings.Contains(err.Error(), errMsgEmptyOutput) && err != vm.ErrExecutionReverted && err != bind.ErrNoCode && err != blockchain.ErrVMDefault {
			logger.Warn("supports interface returns an abnormal error", "err", err, "contract", contract.String(), "interfaceID", hexutil.Encode(interfaceID[:]))
		}
		return false, nil
	}
	return isSupported, nil
}

// isKIP13 checks if the given contract implements KIP13 interface or not at the given block.
func (f *contractCaller) isKIP13(contract common.Address, blockNumber *big.Int) (bool, error) {
	var (
		opts   *bind.CallOpts
		cancel context.CancelFunc
	)
	opts, cancel = getCallOpts(blockNumber, f.callTimeout)
	defer cancel()
	if isKIP13, err := f.supportsInterface(contract, opts, IKIP13Id); err != nil {
		logger.Error("supportsInterface is failed", "contract", contract.String(), "blockNumber", blockNumber, "interfaceID", hexutil.Encode(IKIP13Id[:]))
		return false, err
	} else if !isKIP13 {
		return false, nil
	}

	opts, cancel = getCallOpts(blockNumber, f.callTimeout)
	defer cancel()
	if isInvalid, err := f.supportsInterface(contract, opts, InvalidId); err != nil {
		logger.Error("supportsInterface is failed", "contract", contract.String(), "blockNumber", blockNumber, "interfaceID", hexutil.Encode(InvalidId[:]))
		return false, err
	} else if isInvalid {
		return false, nil
	}

	return true, nil
}

// isKIP7 checks if the given contract implements IKIP7 and IKIP7Metadata interface or not at the given block.
func (f *contractCaller) isKIP7(contract common.Address, blockNumber *big.Int) (bool, error) {
	var (
		opts   *bind.CallOpts
		cancel context.CancelFunc
	)
	opts, cancel = getCallOpts(blockNumber, f.callTimeout)
	defer cancel()
	if isIKIP7, err := f.supportsInterface(contract, opts, IKIP7Id); err != nil {
		logger.Error("supportsInterface is failed", "contract", contract.String(), "blockNumber", blockNumber, "interfaceID", hexutil.Encode(IKIP7Id[:]))
		return false, err
	} else if !isIKIP7 {
		return false, nil
	}

	opts, cancel = getCallOpts(blockNumber, f.callTimeout)
	defer cancel()
	if isIKIP7Metadata, err := f.supportsInterface(contract, opts, IKIP7MetadataId); err != nil {
		logger.Error("supportsInterface is failed", "contract", contract.String(), "blockNumber", blockNumber, "interfaceID", hexutil.Encode(IKIP7MetadataId[:]))
		return false, err
	} else if !isIKIP7Metadata {
		return false, nil
	}

	return true, nil
}

// isKIP17 checks if the given contract implements IKIP17 and IKIP17Metadata interface or not at the given block.
func (f *contractCaller) isKIP17(contract common.Address, blockNumber *big.Int) (bool, error) {
	var (
		opts   *bind.CallOpts
		cancel context.CancelFunc
	)
	opts, cancel = getCallOpts(blockNumber, f.callTimeout)
	defer cancel()
	if isIKIP17, err := f.supportsInterface(contract, opts, IKIP17Id); err != nil {
		logger.Error("supportsInterface is failed", "contract", contract.String(), "blockNumber", blockNumber, "interfaceID", hexutil.Encode(IKIP17Id[:]))
		return false, err
	} else if !isIKIP17 {
		return false, nil
	}

	opts, cancel = getCallOpts(blockNumber, f.callTimeout)
	defer cancel()
	if isIKIP17Metadata, err := f.supportsInterface(contract, opts, IKIP17MetadataId); err != nil {
		logger.Error("supportsInterface is failed", "contract", contract.String(), "blockNumber", blockNumber, "interfaceID", hexutil.Encode(IKIP17MetadataId[:]))
		return false, err
	} else if !isIKIP17Metadata {
		return false, nil
	}

	return true, nil
}
