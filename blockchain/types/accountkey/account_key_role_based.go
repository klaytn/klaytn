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
	"encoding/json"
	"errors"
	"io"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/rlp"
)

type RoleType int

const (
	RoleTransaction RoleType = iota
	RoleAccountUpdate
	RoleFeePayer
	// TODO-Klaytn-Accounts: more roles can be listed here.
	RoleLast
)

var (
	errKeyLengthZero                    = errors.New("key length is zero")
	errKeyShouldNotBeNilOrCompositeType = errors.New("key should not be nil or a composite type")
)

// AccountKeyRoleBased represents a role-based key.
// The roles are defined like below:
// RoleTransaction   - this key is used to verify transactions transferring values.
// RoleAccountUpdate - this key is used to update keys in the account when using TxTypeAccountUpdate.
// RoleFeePayer      - this key is used to pay tx fee when using fee-delegated transactions.
//                     If an account has a key of this role and wants to pay tx fee,
//                     fee-delegated transactions should be signed by this key.
//
// If RoleAccountUpdate or RoleFeePayer is not set, RoleTransaction will be used instead by default.
type AccountKeyRoleBased []AccountKey

func NewAccountKeyRoleBased() *AccountKeyRoleBased {
	return &AccountKeyRoleBased{}
}

func NewAccountKeyRoleBasedWithValues(keys []AccountKey) *AccountKeyRoleBased {
	return (*AccountKeyRoleBased)(&keys)
}

func (a *AccountKeyRoleBased) Type() AccountKeyType {
	return AccountKeyTypeRoleBased
}

func (a *AccountKeyRoleBased) IsCompositeType() bool {
	return true
}

func (a *AccountKeyRoleBased) DeepCopy() AccountKey {
	n := make(AccountKeyRoleBased, len(*a))

	for i, k := range *a {
		n[i] = k.DeepCopy()
	}

	return &n
}

func (a *AccountKeyRoleBased) Equal(b AccountKey) bool {
	tb, ok := b.(*AccountKeyRoleBased)
	if !ok {
		return false
	}

	if len(*a) != len(*tb) {
		return false
	}

	for i, tbi := range *tb {
		if (*a)[i].Equal(tbi) == false {
			return false
		}
	}

	return true
}

func (a *AccountKeyRoleBased) EncodeRLP(w io.Writer) error {
	enc := make([][]byte, len(*a))

	for i, k := range *a {
		enc[i], _ = rlp.EncodeToBytes(NewAccountKeySerializerWithAccountKey(k))
	}

	return rlp.Encode(w, enc)
}

func (a *AccountKeyRoleBased) DecodeRLP(s *rlp.Stream) error {
	enc := [][]byte{}
	if err := s.Decode(&enc); err != nil {
		return err
	}

	keys := make([]AccountKey, len(enc))
	for i, b := range enc {
		serializer := NewAccountKeySerializer()
		if err := rlp.DecodeBytes(b, &serializer); err != nil {
			return err
		}
		keys[i] = serializer.key
	}

	*a = (AccountKeyRoleBased)(keys)

	return nil
}

func (a *AccountKeyRoleBased) MarshalJSON() ([]byte, error) {
	serializers := make([]*AccountKeySerializer, len(*a))

	for i, k := range *a {
		serializers[i] = NewAccountKeySerializerWithAccountKey(k)
	}

	return json.Marshal(serializers)
}

func (a *AccountKeyRoleBased) UnmarshalJSON(b []byte) error {
	var serializers []*AccountKeySerializer
	if err := json.Unmarshal(b, &serializers); err != nil {
		return err
	}

	*a = make(AccountKeyRoleBased, len(serializers))
	for i, s := range serializers {
		(*a)[i] = s.key
	}

	return nil
}

func (a *AccountKeyRoleBased) Validate(currentBlockNumber uint64, r RoleType, recoveredKeys []*ecdsa.PublicKey, from common.Address) bool {
	if len(*a) > int(r) {
		return (*a)[r].Validate(currentBlockNumber, r, recoveredKeys, from)
	}
	return a.getDefaultKey().Validate(currentBlockNumber, r, recoveredKeys, from)
}

func (a *AccountKeyRoleBased) getDefaultKey() AccountKey {
	return (*a)[RoleTransaction]
}

func (a *AccountKeyRoleBased) String() string {
	serializer := NewAccountKeySerializerWithAccountKey(a)
	b, _ := json.Marshal(serializer)
	return string(b)
}

func (a *AccountKeyRoleBased) AccountCreationGas(currentBlockNumber uint64) (uint64, error) {
	gas := uint64(0)
	for _, k := range *a {
		gasK, err := k.AccountCreationGas(currentBlockNumber)
		if err != nil {
			return 0, err
		}
		gas += gasK
	}

	return gas, nil
}

func (a *AccountKeyRoleBased) SigValidationGas(currentBlockNumber uint64, r RoleType, numSigs int) (uint64, error) {
	var key AccountKey
	// Set the key used to sign for validation.
	if len(*a) > int(r) {
		key = (*a)[r]
	} else {
		key = a.getDefaultKey()
	}

	gas, err := key.SigValidationGas(currentBlockNumber, r, numSigs)
	if err != nil {
		return 0, err
	}

	return gas, nil
}

func (a *AccountKeyRoleBased) CheckInstallable(currentBlockNumber uint64) error {
	// A zero-role key is not allowed.
	if len(*a) == 0 {
		return kerrors.ErrZeroLength
	}
	// Do not allow undefined roles.
	if len(*a) > (int)(RoleLast) {
		return kerrors.ErrLengthTooLong
	}
	for i := 0; i < len(*a); i++ {
		// A composite key is not allowed.
		if (*a)[i].IsCompositeType() {
			return kerrors.ErrNestedCompositeType
		}
		// If any key in the role cannot be initialized, return an error.
		if err := (*a)[i].CheckInstallable(currentBlockNumber); err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountKeyRoleBased) CheckUpdatable(newKey AccountKey, currentBlockNumber uint64) error {
	if newKey, ok := newKey.(*AccountKeyRoleBased); ok {
		lenOldKey := len(*a)
		lenNewKey := len(*newKey)
		// If no key is to be replaced, it is regarded as a fail.
		if lenNewKey == 0 {
			return kerrors.ErrZeroLength
		}
		// Do not allow undefined roles.
		if lenNewKey > (int)(RoleLast) {
			return kerrors.ErrLengthTooLong
		}
		for i := 0; i < lenNewKey; i++ {
			switch {
			// A composite key is not allowed.
			case (*newKey)[i].IsCompositeType():
				return kerrors.ErrNestedCompositeType
			// If newKey is longer than oldKey, init the new attributes.
			case i >= lenOldKey:
				if err := (*newKey)[i].CheckInstallable(currentBlockNumber); err != nil {
					return err
				}
			// Do nothing for AccountKeyTypeNil
			case (*newKey)[i].Type() == AccountKeyTypeNil:

			// Check whether the newKey is replacable or not
			default:
				if err := CheckReplacable((*a)[i], (*newKey)[i], currentBlockNumber); err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Update is not possible if the type is different.
	return kerrors.ErrDifferentAccountKeyType
}

func (a *AccountKeyRoleBased) Update(newKey AccountKey, currentBlockNumber uint64) error {
	if err := a.CheckUpdatable(newKey, currentBlockNumber); err != nil {
		return err
	}
	newRoleKey, _ := newKey.(*AccountKeyRoleBased)
	lenNewKey := len(*newRoleKey)
	lenOldKey := len(*a)
	if lenOldKey < lenNewKey {
		*a = append(*a, (*newRoleKey)[lenOldKey:]...)
	}
	for i := 0; i < lenNewKey; i++ {
		if (*newRoleKey)[i].Type() == AccountKeyTypeNil {
			continue
		}
		(*a)[i] = (*newRoleKey)[i]
	}
	return nil
}
