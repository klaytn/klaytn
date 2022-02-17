// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/error.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"errors"

	"github.com/klaytn/klaytn/blockchain/types"
)

var (
	// ErrKnownBlock is returned when a block to import is already known locally.
	ErrKnownBlock = errors.New("block already known")

	// ErrGasLimitReached is returned by the gas pool if the amount of gas required
	// by a transaction is higher than what's left in the block.
	ErrGasLimitReached = errors.New("gas limit reached")

	// ErrBlacklistedHash is returned if a block to import is on the blacklist.
	ErrBlacklistedHash = errors.New("blacklisted hash")

	// ErrNonceTooHigh is returned if the nonce of a transaction is higher than the
	// next one expected based on the local chain.
	ErrNonceTooHigh = errors.New("nonce too high")

	// tx_pool

	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")

	// ErrInvalidFeePayer is returned if the transaction contains an invalid signature of the fee payer.
	ErrInvalidFeePayer = errors.New("invalid fee payer")

	// ErrNonceTooLow is returned if the nonce of a transaction is lower than the
	// one present in the local chain.
	ErrNonceTooLow = errors.New("nonce too low")

	// ErrUnderpriced is returned if a transaction's gas price is below the minimum
	// configured for the transaction pool.
	ErrUnderpriced = errors.New("transaction underpriced")

	// ErrReplaceUnderpriced is returned if a transaction is attempted to be replaced
	// with a different one without the required price bump.
	ErrReplaceUnderpriced = errors.New("replacement transaction underpriced")

	// ErrAlreadyNonceExistInPool is returned if there is another tx with the same nonce in the tx pool.
	ErrAlreadyNonceExistInPool = errors.New("there is another tx which has the same nonce in the tx pool")

	// ErrInsufficientFunds is returned if the total cost of executing a transaction
	// is higher than the balance of the user's account.
	ErrInsufficientFunds = errors.New("insufficient funds for gas * price + value")

	// ErrInsufficientFundsFrom is returned if the value of a transaction is higher than
	// the balance of the user's account.
	ErrInsufficientFundsFrom = errors.New("insufficient funds of the sender for value ")

	// ErrInsufficientFundsFeePayer is returned if the fee of a transaction is higher than
	// the balance of the fee payer's account.
	ErrInsufficientFundsFeePayer = errors.New("insufficient funds of the fee payer for gas * price")

	// ErrIntrinsicGas is returned if the transaction is specified to use less gas
	// than required to start the invocation.
	ErrIntrinsicGas = errors.New("intrinsic gas too low")

	// ErrGasLimit is returned if a transaction's requested gas limit exceeds the
	// maximum allowance of the current block.
	ErrGasLimit = errors.New("exceeds block gas limit")

	// ErrNegativeValue is a sanity error to ensure noone is able to specify a
	// transaction with a negative value.
	ErrNegativeValue = errors.New("negative value")

	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")

	// ErrInvlidUnitPrice is returned if gas price of transaction is not equal to UnitPrice
	ErrInvalidUnitPrice = errors.New("invalid unit price")

	// ErrInvalidChainId is returned if the chain id of transaction is not equal to the chain id of the chain config.
	ErrInvalidChainId = errors.New("invalid chain id")

	// ErrNotYetImplementedAPI is returned if API is not yet implemented
	ErrNotYetImplementedAPI = errors.New("not yet implemented API")

	// Errors returned from GetVMerrFromReceiptStatus

	// ErrInvalidReceiptStatus is returned if status of receipt is invalid from GetVMerrFromReceiptStatus
	ErrInvalidReceiptStatus = errors.New("unknown receipt status")

	// ErrTxTypeNotSupported is returned if a transaction is not supported in the
	// current network configuration.
	ErrTxTypeNotSupported = types.ErrTxTypeNotSupported

	// ErrVMDefault is returned if status of receipt is ReceiptStatusErrDefault from GetVMerrFromReceiptStatus
	ErrVMDefault = errors.New("VM error occurs while running smart contract")

	// ErrAccountCreationPrevented is returned if account creation is inserted in the service chain's txpool.
	ErrAccountCreationPrevented = errors.New("account creation is prevented for the service chain")

	// ErrInvalidTracer is returned if the tracer type is not vm.InternalTxTracer
	ErrInvalidTracer = errors.New("tracer type is invalid for internal transaction tracing")

	// ErrTipVeryHigh is a sanity error to avoid extremely big numbers specified in the tip field.
	ErrTipVeryHigh = errors.New("max priority fee per gas higher than 2^256-1")

	// ErrFeeCapVeryHigh is a sanity error to avoid extremely big numbers specified in the fee cap field.
	ErrFeeCapVeryHigh = errors.New("max fee per gas higher than 2^256-1")

	// ErrTipAboveFeeCap is a sanity error to ensure no one is able to specify a
	// transaction with a tip higher than the total fee cap.
	ErrTipAboveFeeCap = errors.New("max fee per gas higher than max priority fee per gas")

	// ErrInvalidGasFeeCap is returned if gas fee cap of transaction is not equal to UnitPrice
	ErrInvalidGasFeeCap = errors.New("Invalid gas fee cap. It must be set to the same value as gas unit price.")

	// ErrInvalidGasTipCap is returned if gas tip cap of transaction is not equal to UnitPrice
	ErrInvalidGasTipCap = errors.New("Invalid gas tip cap. It must be set to the same value as gas unit price.")
)
