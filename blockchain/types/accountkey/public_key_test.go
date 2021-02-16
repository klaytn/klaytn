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
	"encoding/json"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// TestPublicKeyRLP tests RLP encoding/decoding of PublicKeySerializable.
func TestPublicKeyRLP(t *testing.T) {
	k := newPublicKeySerializable()
	k.X.SetUint64(10)
	k.Y.SetUint64(20)

	b, err := rlp.EncodeToBytes(k)
	assert.Equal(t, err, errNotS256Curve)

	prv, _ := crypto.GenerateKey()
	k = (*PublicKeySerializable)(&prv.PublicKey)
	b, err = rlp.EncodeToBytes(k)
	if err != nil {
		t.Fatal(err)
	}

	dec := newPublicKeySerializable()

	if err := rlp.DecodeBytes(b, &dec); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, k, dec)
	if !k.Equal(dec) {
		t.Fatal("k != dec")
	}
}

// TestPublicKeyRLP tests JSON encoding/decoding of PublicKeySerializable.
func TestPublicKeyJSON(t *testing.T) {
	k := newPublicKeySerializable()
	k.X.SetUint64(10)
	k.Y.SetUint64(20)

	b, err := json.Marshal(k)
	assert.Equal(t, err.(*json.MarshalerError).Err, errNotS256Curve)

	prv, _ := crypto.GenerateKey()
	k = (*PublicKeySerializable)(&prv.PublicKey)
	b, err = json.Marshal(k)
	if err != nil {
		t.Fatal(err)
	}

	dec := newPublicKeySerializable()

	if err := json.Unmarshal(b, &dec); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, k, dec)
	if !k.Equal(dec) {
		t.Fatal("k != dec")
	}
}

// TestPublicKeyRLP tests DeepCopy() of PublicKeySerializable.
func TestPublicKeyDeepCopy(t *testing.T) {
	k := newPublicKeySerializable()
	k.X.SetUint64(10)
	k.Y.SetUint64(20)

	newK := k.DeepCopy()

	newK.X.SetUint64(30)
	newK.Y.SetUint64(40)

	assert.Equal(t, k.X, big.NewInt(10))
	assert.Equal(t, k.Y, big.NewInt(20))
	assert.Equal(t, newK.X, big.NewInt(30))
	assert.Equal(t, newK.Y, big.NewInt(40))
}
