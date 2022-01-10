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
// This file is derived from core/vm/errors.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"errors"
	"fmt"

	"github.com/klaytn/klaytn/params"
)

// List execution errors
var (
	ErrCodeStoreOutOfGas                 = errors.New("contract creation code storage out of gas")
	ErrDepth                             = errors.New("max call depth exceeded")
	ErrTraceLimitReached                 = errors.New("the number of logs reached the specified limit")
	ErrInsufficientBalance               = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision          = errors.New("contract address collision")
	ErrTotalTimeLimitReached             = errors.New("reached the total execution time limit for txs in a block")
	ErrOpcodeComputationCostLimitReached = errors.New(fmt.Sprintf("reached the opcode computation cost limit (%d) for tx", params.OpcodeComputationCostLimit))
	ErrFailedOnSetCode                   = errors.New("failed on setting code to an account")

	// EVM internal errors
	ErrWriteProtection       = errors.New("evm: write protection")
	ErrReturnDataOutOfBounds = errors.New("evm: return data out of bounds")
	ErrExecutionReverted     = errors.New("evm: execution reverted")
	ErrMaxCodeSizeExceeded   = errors.New("evm: max code size exceeded")
	ErrInvalidJump           = errors.New("evm: invalid jump destination")
)
