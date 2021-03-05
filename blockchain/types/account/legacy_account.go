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

package account

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/go-stack/stack"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
)

// LegacyAccount is the Klaytn consensus representation of legacy accounts.
// These objects are stored in the main account trie.
type LegacyAccount struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
}

// newEmptyLegacyAccount returns an empty legacy account.
// This object will be used for RLP-decoding.
func newEmptyLegacyAccount() *LegacyAccount {
	return &LegacyAccount{}
}

// newLegacyAccount returns a LegacyAccount object whose all
// attributes are initialized.
// This object is used when an account is created.
// Refer to StateDB.createObject().
func newLegacyAccount() *LegacyAccount {
	logger.CritWithStack("Legacy account is deprecated.")
	return &LegacyAccount{
		0, new(big.Int), common.Hash{}, emptyCodeHash,
	}
}

func newLegacyAccountWithMap(values map[AccountValueKeyType]interface{}) *LegacyAccount {
	acc := newLegacyAccount()

	if v, ok := values[AccountValueKeyNonce].(uint64); ok {
		acc.Nonce = v
	}

	if v, ok := values[AccountValueKeyBalance].(*big.Int); ok {
		acc.Balance.Set(v)
	}

	if v, ok := values[AccountValueKeyStorageRoot].(common.Hash); ok {
		acc.Root = v
	}

	if v, ok := values[AccountValueKeyCodeHash].([]byte); ok {
		acc.CodeHash = v
	}

	return acc
}

func (a *LegacyAccount) Type() AccountType {
	return LegacyAccountType
}

func (a *LegacyAccount) GetNonce() uint64 {
	return a.Nonce
}

func (a *LegacyAccount) GetBalance() *big.Int {
	return new(big.Int).Set(a.Balance)
}

func (a *LegacyAccount) GetHumanReadable() bool {
	// LegacyAccount cannot be a human-readable account.
	return false
}

func (a *LegacyAccount) GetStorageRoot() common.Hash {
	return a.Root
}

func (a *LegacyAccount) GetCodeHash() []byte {
	return a.CodeHash
}

func (a *LegacyAccount) SetNonce(n uint64) {
	a.Nonce = n
}

func (a *LegacyAccount) SetBalance(b *big.Int) {
	a.Balance.Set(b)
}

func (a *LegacyAccount) SetHumanReadable(b bool) {
	// DO NOTHING.
	logger.Warn("LegacyAccount.SetHumanReadable() should not be called. Please check the call stack.",
		"callstack", stack.Caller(0).String())
}

func (a *LegacyAccount) SetStorageRoot(h common.Hash) {
	a.Root = h
}

func (a *LegacyAccount) SetCodeHash(h []byte) {
	a.CodeHash = h
}

func (a *LegacyAccount) Empty() bool {
	return a.GetNonce() == 0 && a.GetBalance().Sign() == 0 && bytes.Equal(a.GetCodeHash(), emptyCodeHash)
}

func (a *LegacyAccount) UpdateKey(newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	return ErrAccountKeyNotModifiable
}

func (a *LegacyAccount) DeepCopy() Account {
	return &LegacyAccount{
		a.Nonce,
		new(big.Int).Set(a.Balance),
		a.Root,
		a.CodeHash,
	}
}

func (a *LegacyAccount) Dump() {
	fmt.Println(a.String())
}

func (a *LegacyAccount) String() string {
	return fmt.Sprintf(`
	Nonce:        %d
	Balance:      %v
	StorageRoot:  %s
	CodeHash:     %s`,
		a.Nonce,
		a.Balance,
		a.Root.String(),
		string(a.CodeHash))
}

func (a *LegacyAccount) Equal(b Account) bool {
	tb, ok := b.(*LegacyAccount)
	if !ok {
		return false
	}

	return a.Nonce == tb.Nonce &&
		a.Balance.Cmp(tb.Balance) == 0 &&
		bytes.Equal(a.Root.Bytes(), tb.Root.Bytes()) &&
		bytes.Equal(a.CodeHash, tb.CodeHash)
}
