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

package accountkey

import (
	"crypto/ecdsa"
	"io"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/rlp"
)

// AccountKeyNil represents a key having nothing.
// This object is used when executing TxTypeAccountUpdate with a role-based key.
// If an item is an AccountKeyNil object, the corresponding key will not be updated.
// For example, if TxTypeAccountUpdate tries to update the account's key to
// [AccountKeyPublic, AccountKeyNil, AccountKeyPublic], the second item will not be updated in the account.
type AccountKeyNil struct {
}

var globalNilKey = &AccountKeyNil{}

// NewAccountKeyNil creates a new AccountKeyNil object.
// Since AccountKeyNil has no attributes, use one global variable for all allocations.
func NewAccountKeyNil() *AccountKeyNil { return globalNilKey }

func (a *AccountKeyNil) Type() AccountKeyType {
	return AccountKeyTypeNil
}

func (a *AccountKeyNil) IsCompositeType() bool {
	return false
}

func (a *AccountKeyNil) Equal(b AccountKey) bool {
	if _, ok := b.(*AccountKeyNil); !ok {
		return false
	}

	// if b is a type of AccountKeyNil, just return true.
	return true
}

func (a *AccountKeyNil) EncodeRLP(w io.Writer) error {
	return nil
}

func (a *AccountKeyNil) DecodeRLP(s *rlp.Stream) error {
	return nil
}

func (a *AccountKeyNil) Validate(currentBlockNumber uint64, r RoleType, recoveredKeys []*ecdsa.PublicKey, from common.Address) bool {
	logger.ErrorWithStack("this function should not be called. Validation should be done at ValidateSender or ValidateFeePayer")
	return false
}

func (a *AccountKeyNil) String() string {
	return "AccountKeyNil"
}

func (a *AccountKeyNil) DeepCopy() AccountKey {
	return NewAccountKeyNil()
}

func (a *AccountKeyNil) AccountCreationGas(currentBlockNumber uint64) (uint64, error) {
	// No gas required to make an account with a nil key.
	return 0, nil
}

func (a *AccountKeyNil) SigValidationGas(currentBlockNumber uint64, r RoleType, validSigNum int) (uint64, error) {
	// No gas required to make an account with a nil key.
	return 0, nil
}

func (a *AccountKeyNil) CheckInstallable(currentBlockNumber uint64) error {
	// Since AccountKeyNil cannot be assigned to an account, it always returns error.
	return kerrors.ErrAccountKeyNilUninitializable
}

func (a *AccountKeyNil) CheckUpdatable(newKey AccountKey, currentBlockNumber uint64) error {
	// Since AccountKeyNil cannot be assigned to an account, it should not be called.
	return kerrors.ErrAccountKeyNilUninitializable
}

func (a *AccountKeyNil) Update(newKey AccountKey, currentBlockNumber uint64) error {
	// Since AccountKeyNil cannot be assigned to an account, it should not be called.
	return kerrors.ErrDifferentAccountKeyType
}
