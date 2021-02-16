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
	"testing"

	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/rlp"
)

func TestAccountKeySerialization(t *testing.T) {
	var keys = []struct {
		Name string
		k    AccountKey
	}{
		{"Nil", genAccountKeyNil()},
		{"Legacy", genAccountKeyLegacy()},
		{"Public", genAccountKeyPublic()},
		{"Fail", genAccountKeyFail()},
		{"WeightedMultisig", genAccountKeyWeightedMultisig()},
		{"RoleBased", genAccountKeyRoleBased()},
	}

	var testcases = []struct {
		Name string
		fn   func(t *testing.T, k AccountKey)
	}{
		{"RLP", testAccountKeyRLP},
		{"JSON", testAccountKeyJSON},
	}
	for _, test := range testcases {
		for _, key := range keys {
			Name := test.Name + "/" + key.Name
			t.Run(Name, func(t *testing.T) {
				test.fn(t, key.k)
			})
		}
	}
}

func testAccountKeyRLP(t *testing.T, k AccountKey) {
	enc := NewAccountKeySerializerWithAccountKey(k)

	b, err := rlp.EncodeToBytes(enc)
	if err != nil {
		t.Fatal(err)
	}

	dec := NewAccountKeySerializer()

	if err := rlp.DecodeBytes(b, &dec); err != nil {
		t.Fatal(err)
	}

	switch k.Type() {
	case AccountKeyTypeFail:
		if k.Equal(dec.key) {
			t.Errorf("AlwaysFail key returns true! k != dec.key\nk=%v\ndec.key=%v", k, dec.key)
		}
	default:
		if !k.Equal(dec.key) {
			t.Errorf("AlwaysFail key returns true! k != dec.key\nk=%v\ndec.key=%v", k, dec.key)
		}
	}
}

func testAccountKeyJSON(t *testing.T, k AccountKey) {
	enc := NewAccountKeySerializerWithAccountKey(k)

	b, err := json.Marshal(enc)
	if err != nil {
		t.Fatal(err)
	}

	dec := NewAccountKeySerializer()

	if err := json.Unmarshal(b, &dec); err != nil {
		t.Fatal(err)
	}

	switch k.Type() {
	case AccountKeyTypeFail:
		if k.Equal(dec.key) {
			t.Errorf("AlwaysFail key returns true! k != dec.key\nk=%v\ndec.key=%v", k, dec.key)
		}
	default:
		if !k.Equal(dec.key) {
			t.Errorf("AlwaysFail key returns true! k != dec.key\nk=%v\ndec.key=%v", k, dec.key)
		}
	}
}

func genAccountKeyNil() AccountKey {
	return NewAccountKeyNil()
}

func genAccountKeyLegacy() AccountKey {
	return NewAccountKeyLegacy()
}

func genAccountKeyPublic() AccountKey {
	k, _ := crypto.GenerateKey()
	return NewAccountKeyPublicWithValue(&k.PublicKey)
}

func genAccountKeyFail() AccountKey {
	return NewAccountKeyFail()
}

func genAccountKeyWeightedMultisig() AccountKey {
	threshold := uint(3)
	numKeys := 4
	keys := make(WeightedPublicKeys, numKeys)

	for i := 0; i < numKeys; i++ {
		k, _ := crypto.GenerateKey()
		keys[i] = NewWeightedPublicKey(1, (*PublicKeySerializable)(&k.PublicKey))
	}

	return NewAccountKeyWeightedMultiSigWithValues(threshold, keys)
}

func genAccountKeyRoleBased() AccountKey {
	k1, err := crypto.HexToECDSA("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	if err != nil {
		panic(err)
	}
	txKey := NewAccountKeyPublicWithValue(&k1.PublicKey)

	k2, err := crypto.HexToECDSA("c64f2cd1196e2a1791365b00c4bc07ab8f047b73152e4617c6ed06ac221a4b0c")
	if err != nil {
		panic(err)
	}
	threshold := uint(2)
	keys := WeightedPublicKeys{
		NewWeightedPublicKey(1, (*PublicKeySerializable)(&k1.PublicKey)),
		NewWeightedPublicKey(1, (*PublicKeySerializable)(&k2.PublicKey)),
	}
	updateKey := NewAccountKeyWeightedMultiSigWithValues(threshold, keys)

	k3, err := crypto.HexToECDSA("ed580f5bd71a2ee4dae5cb43e331b7d0318596e561e6add7844271ed94156b20")
	if err != nil {
		panic(err)
	}
	feeKey := NewAccountKeyPublicWithValue(&k3.PublicKey)

	return NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{txKey, updateKey, feeKey})
}
