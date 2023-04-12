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
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	// Correctness check is done in Sign() and Verify() tests
	_, err := AggregatePublicKeys(tc.pks)
	assert.Nil(t, err)

	_, err = AggregatePublicKeys(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)
}

func TestAggregatePublicKeysFromBytes(t *testing.T) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)
	bs := tc.pkbs

	// Correctness check is done in Sign() and Verify() tests
	_, err := AggregatePublicKeysFromBytes(bs)
	assert.Nil(t, err)

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
