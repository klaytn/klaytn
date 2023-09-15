// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/evm.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/params"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}

// NewEVMBlockContext creates a new context for use in the EVM.
func NewEVMBlockContext(header *types.Header, chain ChainContext, author *common.Address) vm.BlockContext {
	// If we don't have an explicit author (i.e. not mining), extract from the header
	var (
		beneficiary common.Address
		rewardBase  common.Address
		baseFee     *big.Int
	)

	if author == nil {
		beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}

	rewardBase = header.Rewardbase

	if header.BaseFee != nil {
		baseFee = header.BaseFee
	} else {
		// before magma hardfork, base fee is 0, effectiveGasPrice is unitPrice
		baseFee = new(big.Int).SetUint64(params.ZeroBaseFee)
	}

	return vm.BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, chain),
		Coinbase:    beneficiary,
		Rewardbase:  rewardBase,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).Set(header.Time),
		BlockScore:  new(big.Int).Set(header.BlockScore),
		BaseFee:     baseFee,
	}
}

// NewEVMTxContext creates a new transaction context for a single transaction.
func NewEVMTxContext(msg Message, header *types.Header) vm.TxContext {
	effectiveGasPrice := msg.GasPrice()
	if header.BaseFee != nil {
		effectiveGasPrice = header.BaseFee
	}

	return vm.TxContext{
		Origin:   msg.ValidatedSender(),
		GasPrice: new(big.Int).Set(effectiveGasPrice),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	// Cache will initially contain [refHash.parent],
	// Then fill up with [refHash.p, refHash.pp, refHash.ppp, ...]
	var cache []common.Hash

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		if len(cache) == 0 {
			cache = append(cache, ref.ParentHash)
		}
		if idx := ref.Number.Uint64() - n - 1; idx < uint64(len(cache)) {
			return cache[idx]
		}
		// No luck in the cache, but we can start iterating from the last element we already know
		lastKnownHash := cache[len(cache)-1]

		lastKnownNumber := ref.Number.Uint64() - uint64(len(cache))

		for {
			header := chain.GetHeader(lastKnownHash, lastKnownNumber)
			if header == nil {
				break
			}
			cache = append(cache, header.ParentHash)
			lastKnownHash = header.ParentHash
			lastKnownNumber = header.Number.Uint64() - 1
			if n == lastKnownNumber {
				return lastKnownHash
			}
		}
		return common.Hash{}
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

func DoEstimateGas(ctx context.Context, gasLimit, rpcGasCap uint64, txValue, gasPrice, balance *big.Int, test func(gas uint64) (bool, *ExecutionResult, error)) (hexutil.Uint64, error) {
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64 = params.UpperGasLimit
		cap uint64
	)

	// Initialize nil params
	if txValue == nil {
		txValue = big.NewInt(0)
	}
	if gasPrice == nil {
		gasPrice = big.NewInt(0)
	}
	if balance == nil {
		balance = big.NewInt(0)
	}

	if gasLimit >= params.TxGas {
		hi = gasLimit
	}

	// recap the highest gas limit with account's available balance.
	if gasPrice.BitLen() != 0 {
		available := new(big.Int).Set(balance)
		if txValue.Cmp(available) >= 0 {
			return 0, errors.New("insufficient funds for transfer")
		}
		available.Sub(available, txValue)
		allowance := new(big.Int).Div(available, gasPrice)

		// If the allowance is larger than maximum uint64, skip checking
		if allowance.IsUint64() && hi > allowance.Uint64() {
			logger.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", txValue, "maxFeePerGas", gasPrice, "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	// Recap the highest gas allowance with specified gascap.
	if rpcGasCap != 0 && hi > rpcGasCap {
		logger.Warn("Caller gas above allowance, capping", "requested", hi, "cap", rpcGasCap)
		hi = rpcGasCap
	}
	cap = hi

	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := test(mid)
		if err != nil {
			return 0, err
		}

		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		failed, result, err := test(hi)
		if err != nil {
			return 0, err
		}
		if failed {
			if result != nil && result.VmExecutionStatus != types.ReceiptStatusErrOutOfGas {
				if len(result.Revert()) > 0 {
					return 0, NewRevertError(result)
				}
				return 0, result.Unwrap()
			}
			// Otherwise, the specified gas cap is too low
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}
	return hexutil.Uint64(hi), nil
}

func NewRevertError(result *ExecutionResult) *RevertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &RevertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// RevertError is an API error that encompassas an EVM revertal with JSON error
// code and a binary data blob.
type RevertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *RevertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *RevertError) ErrorData() interface{} {
	return e.reason
}
