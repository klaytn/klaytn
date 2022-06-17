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

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/kerrors"
)

// AccountKeyLegacy is used for accounts having no keys.
// In this case, verifying the signature of a transaction uses the legacy scheme.
// 1. The address comes from the public key which is derived from txhash and the tx's signature.
// 2. Check that the address is the same as the address in the tx.
// It is implemented to support LegacyAccounts.
type AccountKeyLegacy struct {
}

var globalLegacyKey = &AccountKeyLegacy{}

// NewAccountKeyLegacy creates a new AccountKeyLegacy object.
// Since AccountKeyLegacy has no attributes, use one global variable for all allocations.
func NewAccountKeyLegacy() *AccountKeyLegacy { return globalLegacyKey }

func (a *AccountKeyLegacy) Type() AccountKeyType {
	return AccountKeyTypeLegacy
}

func (a *AccountKeyLegacy) IsCompositeType() bool {
	return false
}

func (a *AccountKeyLegacy) Equal(b AccountKey) bool {
	if _, ok := b.(*AccountKeyLegacy); !ok {
		return false
	}

	// if b is a type of AccountKeyLegacy, just return true.
	return true
}

func (a *AccountKeyLegacy) Validate(currentBlockNumber uint64, r RoleType, recoveredKeys []*ecdsa.PublicKey, from common.Address) bool {
	// A valid legacy key generates only one signature
	if len(recoveredKeys) != 1 {
		return false
	}
	return from == crypto.PubkeyToAddress(*recoveredKeys[0])
}

func (a *AccountKeyLegacy) String() string {
	return "AccountKeyLegacy"
}

func (a *AccountKeyLegacy) DeepCopy() AccountKey {
	return NewAccountKeyLegacy()
}

func (a *AccountKeyLegacy) AccountCreationGas(currentBlockNumber uint64) (uint64, error) {
	// No gas required to make an account with a nil key.

	return 0, nil
}

func (a *AccountKeyLegacy) SigValidationGas(currentBlockNumber uint64, r RoleType, validSigNum int) (uint64, error) {
	// No gas required to make an account with a nil key.
	return 0, nil
}

func (a *AccountKeyLegacy) CheckInstallable(currentBlockNumber uint64) error {
	// Since it has no data and it can be assigned to an account, it always returns nil.
	return nil
}

func (a *AccountKeyLegacy) CheckUpdatable(newKey AccountKey, currentBlockNumber uint64) error {
	if _, ok := newKey.(*AccountKeyLegacy); ok {
		return a.CheckInstallable(currentBlockNumber)
	}
	// Update is not possible if the type is different.
	return kerrors.ErrDifferentAccountKeyType
}

func (a *AccountKeyLegacy) Update(newKey AccountKey, currentBlockNumber uint64) error {
	return a.CheckUpdatable(newKey, currentBlockNumber)
}
