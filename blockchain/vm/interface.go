// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/vm/interface.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateAccount(common.Address)
	CreateSmartContractAccount(addr common.Address, format params.CodeFormat, r params.Rules)
	CreateSmartContractAccountWithKey(addr common.Address, humanReadable bool, key accountkey.AccountKey, format params.CodeFormat, r params.Rules)
	CreateEOA(addr common.Address, humanReadable bool, key accountkey.AccountKey)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	IncNonce(common.Address)
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte) error
	GetCodeSize(common.Address) int
	GetVmVersion(common.Address) (params.VmVersion, bool)

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(common.Address, common.Hash) common.Hash
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// UpdateKey updates the account's key with the given key.
	UpdateKey(addr common.Address, newKey accountkey.AccountKey, currentBlockNumber uint64) error

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	// IsProgramAccount returns true if the account implements ProgramAccount.
	IsProgramAccount(address common.Address) bool
	IsContractAvailable(address common.Address) bool
	IsValidCodeFormat(addr common.Address) bool

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)

	GetTxHash() common.Hash

	GetKey(address common.Address) accountkey.AccountKey
}
