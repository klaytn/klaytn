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
	"fmt"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
)

const numKeys = 1

// AccountKeyPublic is used for accounts having one public key.
// In this case, verifying the signature of a transaction is performed as following:
// 1. The public key is derived from the hash and the signature of the tx.
// 2. Check that the derived public key is the same as the corresponding account's public key.
type AccountKeyPublic struct {
	*PublicKeySerializable
}

func NewAccountKeyPublicWithValue(pk *ecdsa.PublicKey) *AccountKeyPublic {
	return &AccountKeyPublic{(*PublicKeySerializable)(pk)}
}

func NewAccountKeyPublic() *AccountKeyPublic {
	return &AccountKeyPublic{newPublicKeySerializable()}
}

func (a *AccountKeyPublic) Type() AccountKeyType {
	return AccountKeyTypePublic
}

func (a *AccountKeyPublic) IsCompositeType() bool {
	return false
}

func (a *AccountKeyPublic) DeepCopy() AccountKey {
	return &AccountKeyPublic{
		a.PublicKeySerializable.DeepCopy(),
	}
}
func (a *AccountKeyPublic) Equal(b AccountKey) bool {
	tb, ok := b.(*AccountKeyPublic)
	if !ok {
		return false
	}
	return a.PublicKeySerializable.Equal(tb.PublicKeySerializable)
}

func (a *AccountKeyPublic) Validate(currentBlockNumber uint64, r RoleType, recoveredKeys []*ecdsa.PublicKey, from common.Address) bool {
	// AccountKeyPublic has only one public key.
	if len(recoveredKeys) != 1 {
		return false
	}

	return a.PublicKeySerializable.Equal((*PublicKeySerializable)(recoveredKeys[0]))
}

func (a *AccountKeyPublic) String() string {
	return fmt.Sprintf("AccountKeyPublic: %s", a.PublicKeySerializable.String())
}

func (a *AccountKeyPublic) AccountCreationGas(currentBlockNumber uint64) (uint64, error) {
	return numKeys * params.TxAccountCreationGasPerKey, nil
}

func (a *AccountKeyPublic) SigValidationGas(currentBlockNumber uint64, r RoleType, validSigNum int) (uint64, error) {
	return (numKeys - 1) * params.TxValidationGasPerKey, nil
}

func (a *AccountKeyPublic) CheckInstallable(currentBlockNumber uint64) error {
	// If the point is not on the curve, return an error.
	if a.IsOnCurve(a.X, a.Y) == false {
		return kerrors.ErrNotOnCurve
	}
	return nil
}

func (a *AccountKeyPublic) CheckUpdatable(newKey AccountKey, currentBlockNumber uint64) error {
	if newKey, ok := newKey.(*AccountKeyPublic); ok {
		return newKey.CheckInstallable(currentBlockNumber)
	}
	// Update is not possible if the type is different.
	return kerrors.ErrDifferentAccountKeyType
}

func (a *AccountKeyPublic) Update(newKey AccountKey, currentBlockNumber uint64) error {
	if err := a.CheckUpdatable(newKey, currentBlockNumber); err != nil {
		return err
	}
	newPubKey, _ := newKey.(*AccountKeyPublic)
	a.X = newPubKey.X
	a.Y = newPubKey.Y
	return nil
}
