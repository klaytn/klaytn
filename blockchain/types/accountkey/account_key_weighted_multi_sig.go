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

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
)

const (
	// TODO-Klaytn-MultiSig: Need to fix the maximum number of keys allowed for an account.
	// NOTE-Klaytn-MultiSig: This value should not be reduced. If it is reduced, there is a case:
	// - the tx validation will be failed if the sender has larger keys.
	MaxNumKeysForMultiSig = uint64(10)
)

// AccountKeyWeightedMultiSig is an account key type containing a threshold and `WeightedPublicKeys`.
// `WeightedPublicKeys` contains a slice of {weight and key}.
// To be a valid tx for an account associated with `AccountKeyWeightedMultiSig`,
// the weighted sum of signed public keys should be larger than the threshold.
// Refer to AccountKeyWeightedMultiSig.Validate().
type AccountKeyWeightedMultiSig struct {
	Threshold uint               `json:"threshold"`
	Keys      WeightedPublicKeys `json:"keys"`
}

func NewAccountKeyWeightedMultiSig() *AccountKeyWeightedMultiSig {
	return &AccountKeyWeightedMultiSig{}
}

func NewAccountKeyWeightedMultiSigWithValues(threshold uint, keys WeightedPublicKeys) *AccountKeyWeightedMultiSig {
	return &AccountKeyWeightedMultiSig{threshold, keys}
}

func (a *AccountKeyWeightedMultiSig) Type() AccountKeyType {
	return AccountKeyTypeWeightedMultiSig
}

func (a *AccountKeyWeightedMultiSig) IsCompositeType() bool {
	return false
}

func (a *AccountKeyWeightedMultiSig) DeepCopy() AccountKey {
	return &AccountKeyWeightedMultiSig{
		a.Threshold, a.Keys.DeepCopy(),
	}
}

func (a *AccountKeyWeightedMultiSig) Equal(b AccountKey) bool {
	tb, ok := b.(*AccountKeyWeightedMultiSig)
	if !ok {
		return false
	}

	return a.Threshold == tb.Threshold &&
		a.Keys.Equal(tb.Keys)
}

func (a *AccountKeyWeightedMultiSig) Validate(r RoleType, recoveredKeys []*ecdsa.PublicKey, from common.Address) bool {
	weightedSum := uint(0)

	// To prohibit making a signature with the same key, make a map.
	// TODO-Klaytn: find another way for better performance
	pMap := make(map[string]*ecdsa.PublicKey)
	for _, bk := range recoveredKeys {
		b, err := rlp.EncodeToBytes((*PublicKeySerializable)(bk))
		if err != nil {
			logger.Warn("Failed to encode recovered public keys of the tx", "recoveredKeys", recoveredKeys)
			continue
		}
		pMap[string(b)] = bk
	}

	for _, k := range a.Keys {
		b, err := rlp.EncodeToBytes(k.Key)
		if err != nil {
			logger.Warn("Failed to encode public keys in the account", "AccountKey", a.String())
			continue
		}

		if _, ok := pMap[string(b)]; ok {
			weightedSum += k.Weight
		}
	}

	if weightedSum >= a.Threshold {
		return true
	}

	logger.Debug("AccountKeyWeightedMultiSig validation is failed", "recoveredKeys", recoveredKeys,
		"accountKeys", a.String(), "threshold", a.Threshold, "weighted sum", weightedSum)

	return false
}

func (a *AccountKeyWeightedMultiSig) String() string {
	serializer := NewAccountKeySerializerWithAccountKey(a)
	b, _ := json.Marshal(serializer)
	return string(b)
}

func (a *AccountKeyWeightedMultiSig) AccountCreationGas(currentBlockNumber uint64) (uint64, error) {
	numKeys := uint64(len(a.Keys))
	if numKeys > MaxNumKeysForMultiSig {
		return 0, kerrors.ErrMaxKeysExceed
	}

	return numKeys * params.TxAccountCreationGasPerKey, nil
}

func (a *AccountKeyWeightedMultiSig) SigValidationGas(currentBlockNumber uint64, r RoleType) (uint64, error) {
	numKeys := uint64(len(a.Keys))
	if numKeys > MaxNumKeysForMultiSig {
		logger.Warn("validation failed due to the number of keys in the account is larger than the limit.",
			"account", a.String())
		return 0, kerrors.ErrMaxKeysExceedInValidation
	}
	if numKeys == 0 {
		logger.Error("should not happen! numKeys is equal to zero!")
		return 0, kerrors.ErrZeroLength
	}
	return (numKeys - 1) * params.TxValidationGasPerKey, nil
}

func (a *AccountKeyWeightedMultiSig) CheckInstallable(currentBlockNumber uint64) error {
	sum := uint(0)
	prevSum := uint(0)
	if len(a.Keys) == 0 {
		return kerrors.ErrZeroLength
	}
	if uint64(len(a.Keys)) > MaxNumKeysForMultiSig {
		return kerrors.ErrMaxKeysExceed
	}
	keyMap := make(map[string]bool)
	for _, k := range a.Keys {
		// Do not allow zero weight.
		if k.Weight == 0 {
			return kerrors.ErrZeroKeyWeight
		}
		sum += k.Weight

		b, err := rlp.EncodeToBytes(k.Key)
		if err != nil {
			// Do not allow unserializable keys.
			return kerrors.ErrUnserializableKey
		}
		if _, ok := keyMap[string(b)]; ok {
			// Do not allow duplicated keys.
			return kerrors.ErrDuplicatedKey
		}
		keyMap[string(b)] = true

		// Do not allow overflow of weighted sum.
		if prevSum > sum {
			return kerrors.ErrWeightedSumOverflow
		}
		prevSum = sum
	}
	// The weighted sum should be larger than the threshold.
	if sum < a.Threshold {
		return kerrors.ErrUnsatisfiableThreshold
	}
	return nil
}

func (a *AccountKeyWeightedMultiSig) CheckUpdatable(newKey AccountKey, currentBlockNumber uint64) error {
	if newKey, ok := newKey.(*AccountKeyWeightedMultiSig); ok {
		return newKey.CheckInstallable(currentBlockNumber)
	}
	// Update is not possible if the type is different.
	return kerrors.ErrDifferentAccountKeyType
}

func (a *AccountKeyWeightedMultiSig) Update(newKey AccountKey, currentBlockNumber uint64) error {
	if err := a.CheckUpdatable(newKey, currentBlockNumber); err != nil {
		return err
	}
	newMultiKey, _ := newKey.(*AccountKeyWeightedMultiSig)
	a.Threshold = newMultiKey.Threshold
	a.Keys = make(WeightedPublicKeys, len(newMultiKey.Keys))
	copy(a.Keys, newMultiKey.Keys)
	return nil
}

// WeightedPublicKey contains a public key and its weight.
// The weight is used to check whether the weighted sum of public keys are larger than
// the threshold of the AccountKeyWeightedMultiSig object.
type WeightedPublicKey struct {
	Weight uint                   `json:"weight"`
	Key    *PublicKeySerializable `json:"key"`
}

func (w *WeightedPublicKey) Equal(b *WeightedPublicKey) bool {
	return w.Weight == b.Weight &&
		w.Key.Equal(b.Key)
}

func NewWeightedPublicKey(weight uint, key *PublicKeySerializable) *WeightedPublicKey {
	return &WeightedPublicKey{weight, key}
}

// WeightedPublicKeys is a slice of WeightedPublicKey objects.
type WeightedPublicKeys []*WeightedPublicKey

func (w WeightedPublicKeys) DeepCopy() WeightedPublicKeys {
	keys := make(WeightedPublicKeys, len(w))

	for i, v := range w {
		keys[i] = NewWeightedPublicKey(v.Weight, v.Key.DeepCopy())
	}

	return keys
}

func (w WeightedPublicKeys) Equal(b WeightedPublicKeys) bool {
	if len(w) != len(b) {
		return false
	}

	for i, wv := range w {
		if !wv.Equal(b[i]) {
			return false
		}
	}

	return true
}
