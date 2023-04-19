// Copyright 2023 The klaytn Authors
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

package blst

import (
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto/bls/types"
	"github.com/stretchr/testify/assert"
)

var (
	// https://github.com/ethereum/bls12-381-tests
	// verify/verify_valid_case_195246ee3bd3b6ec.json
	testPublicKeyBytes = common.FromHex("0xb53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
	// deserialization_G1/deserialization_fails_not_in_G1.json
	testBadPublicKeyBytes = common.FromHex("0x8123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
)

func TestPublicKeyCopy(t *testing.T) {
	pk, _ := PublicKeyFromBytes(testPublicKeyBytes)
	assert.Equal(t, pk.Marshal(), pk.Copy().Marshal())
}

func TestPublicKeyFromBytes(t *testing.T) {
	b := testPublicKeyBytes
	pk, err := PublicKeyFromBytes(b)
	assert.Nil(t, err)
	assert.Equal(t, b, pk.Marshal())

	_, err = PublicKeyFromBytes([]byte{1, 2, 3, 4})
	assert.NotNil(t, err)

	zero := make([]byte, types.PublicKeyLength)
	_, err = PublicKeyFromBytes(zero)
	assert.Equal(t, types.ErrPublicKeyUnmarshal, err)

	_, err = PublicKeyFromBytes(testBadPublicKeyBytes)
	assert.Equal(t, types.ErrPublicKeyUnmarshal, err)
}

func TestMultiplePublicKeysFromBytes(t *testing.T) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)
	bs := tc.pkbs

	// all uncached
	pks, err := MultiplePublicKeysFromBytes(tc.pkbs[:L/2])
	assert.Nil(t, err)
	for i := 0; i < L/2; i++ {
		assert.Equal(t, tc.pks[i].Marshal(), pks[i].Marshal())
	}

	// [0:L/2] are cached and [L/2:L] are uncached
	pks, err = MultiplePublicKeysFromBytes(tc.pkbs)
	assert.Nil(t, err)
	for i := 0; i < L; i++ {
		assert.Equal(t, tc.pks[i].Marshal(), pks[i].Marshal())
	}

	bs[0] = []byte{1, 2, 3, 4} // len?
	_, err = MultiplePublicKeysFromBytes(bs)
	assert.NotNil(t, err)

	bs[0] = make([]byte, types.PublicKeyLength) // is_inf?
	_, err = MultiplePublicKeysFromBytes(bs)
	assert.Equal(t, types.ErrPublicKeyUnmarshal, err)

	bs[0] = testBadPublicKeyBytes // in_g2?
	_, err = MultiplePublicKeysFromBytes(bs)
	assert.Equal(t, types.ErrPublicKeyUnmarshal, err)
}

func TestAggregatePublicKeys(t *testing.T) {
	var (
		// import milagro_bls_binding as bls
		// from binascii import unhexlify as unhex
		// pk1 = bls.SkToPk(unhex('263dbd792f5b1be47ed85f8938c0f29586af0d3ac7b977f21c278fe1462040e3'))
		// pk2 = bls.SkToPk(unhex('47b8192d77bf871b62e87859d653922725724a5c031afeabc60bcef5ff665138'))
		// pk3 = bls.SkToPk(unhex('328388aff0d4a5b7dc9205abd374e7e98f3cd9f3418edb4eafda5fb16473d216'))
		// apk = bls._AggregatePKs([pk1,pk2,pk3])
		// print('\n'.join(map(lambda k: k.hex(), [pk1,pk2,pk3,apk])))
		pkb1 = common.FromHex("a491d1b0ecd9bb917989f0e74f0dea0422eac4a873e5e2644f368dffb9a6e20fd6e10c1b77654d067c0618f6e5a7f79a")
		pkb2 = common.FromHex("b301803f8b5ac4a1133581fc676dfedc60d891dd5fa99028805e5ea5b08d3491af75d0707adab3b70c6a6a580217bf81")
		pkb3 = common.FromHex("b53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
		apkb = common.FromHex("a095608b35495ca05002b7b5966729dd1ed096568cf2ff24f3318468e0f3495361414a78ebc09574489bc79e48fca969")

		pk1, _ = PublicKeyFromBytes(pkb1)
		pk2, _ = PublicKeyFromBytes(pkb2)
		pk3, _ = PublicKeyFromBytes(pkb3)
	)

	apk, err := AggregatePublicKeys([]types.PublicKey{pk1, pk2, pk3})
	assert.Nil(t, err)
	assert.Equal(t, apkb, apk.Marshal())

	_, err = AggregatePublicKeys(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)
}

func TestAggregatePublicKeysFromBytes(t *testing.T) {
	var (
		// See TestAggregatePublicKeys for the data source
		pkb1 = common.FromHex("a491d1b0ecd9bb917989f0e74f0dea0422eac4a873e5e2644f368dffb9a6e20fd6e10c1b77654d067c0618f6e5a7f79a")
		pkb2 = common.FromHex("b301803f8b5ac4a1133581fc676dfedc60d891dd5fa99028805e5ea5b08d3491af75d0707adab3b70c6a6a580217bf81")
		pkb3 = common.FromHex("b53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
		apkb = common.FromHex("a095608b35495ca05002b7b5966729dd1ed096568cf2ff24f3318468e0f3495361414a78ebc09574489bc79e48fca969")
	)

	apk, err := AggregatePublicKeysFromBytes([][]byte{pkb1, pkb2, pkb3})
	assert.Nil(t, err)
	assert.Equal(t, apkb, apk.Marshal())

	_, err = AggregatePublicKeysFromBytes(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)

	bs := make([][]byte, 1)
	bs[0] = []byte{1, 2, 3, 4} // len?
	_, err = AggregatePublicKeysFromBytes(bs)
	assert.NotNil(t, err)

	bs[0] = make([]byte, types.PublicKeyLength) // is_inf?
	_, err = AggregatePublicKeysFromBytes(bs)
	assert.Equal(t, types.ErrPublicKeyUnmarshal, err)

	bs[0] = testBadPublicKeyBytes // in_g2?
	_, err = AggregatePublicKeysFromBytes(bs)
	assert.Equal(t, types.ErrPublicKeyUnmarshal, err)
}
