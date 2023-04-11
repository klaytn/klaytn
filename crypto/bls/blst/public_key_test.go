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
	// deserialization_G2/deserialization_fails_not_in_G2.json
	testBadPublicKeyBytes = common.FromHex("8123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	// aggregate/aggregate_0xabababababababababababababababababababababababababababababababab.json
	testAggPublicKeyBytes1 = common.FromHex("b6ed936746e01f8ecf281f020953fbf1f01debd5657c4a383940b020b26507f6076334f91e2366c96e9ab279fb5158090352ea1c5b0c9274504f4f0e7053af24802e51e4568d164fe986834f41e55c8e850ce1f98458c0cfc9ab380b55285a55")
	testAggPublicKeyBytes2 = common.FromHex("b23c46be3a001c63ca711f87a005c200cc550b9429d5f4eb38d74322144f1b63926da3388979e5321012fb1a0526bcd100b5ef5fe72628ce4cd5e904aeaa3279527843fae5ca9ca675f4f51ed8f83bbf7155da9ecc9663100a885d5dc6df96d9")
	testAggPublicKeyBytes3 = common.FromHex("948a7cb99f76d616c2c564ce9bf4a519f1bea6b0a624a02276443c245854219fabb8d4ce061d255af5330b078d5380681751aa7053da2c98bae898edc218c75f07e24d8802a17cd1f6833b71e58f5eb5b94208b4d0bb3848cecb075ea21be115")
	testAggPublicKeyBytes  = common.FromHex("9683b3e6701f9a4b706709577963110043af78a5b41991b998475a3d3fd62abf35ce03b33908418efc95a058494a8ae504354b9f626231f6b3f3c849dfdeaf5017c4780e2aee1850ceaf4b4d9ce70971a3d2cfcd97b7e5ecf6759f8da5f76d31")
)

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
	pk1, _ := PublicKeyFromBytes(testAggPublicKeyBytes1)
	pk2, _ := PublicKeyFromBytes(testAggPublicKeyBytes2)
	pk3, _ := PublicKeyFromBytes(testAggPublicKeyBytes3)
	pks := []types.PublicKey{pk1, pk2, pk3}

	apk, err := AggregatePublicKeys(pks)
	assert.Nil(t, err)
	assert.Equal(t, testAggPublicKeyBytes, apk.Marshal())

	_, err = AggregatePublicKeys(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)
}

func TestAggregatePublicKeysFromBytes(t *testing.T) {
	bs := [][]byte{testAggPublicKeyBytes1, testAggPublicKeyBytes2, testAggPublicKeyBytes3}

	apk, err := AggregatePublicKeysFromBytes(bs)
	assert.Nil(t, err)
	assert.Equal(t, testAggPublicKeyBytes, apk.Marshal())

	_, err = AggregatePublicKeysFromBytes(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)

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
