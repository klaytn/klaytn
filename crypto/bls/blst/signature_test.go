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
	// deserialization_G1/deserialization_succeeds_correct_point.json
	testSignatureBytes = common.FromHex("a491d1b0ecd9bb917989f0e74f0dea0422eac4a873e5e2644f368dffb9a6e20fd6e10c1b77654d067c0618f6e5a7f79a")
	// deserialization_G1/deserialization_fails_not_in_curve.json
	testBadSignatureBytes = common.FromHex("8123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde0")
	// sign/sign_case_142f678a8d05fcd1.json
	testMessage = common.FromHex("0x5656565656565656565656565656565656565656565656565656565656565656")
)

func TestSignatureFromBytes(t *testing.T) {
	b := testSignatureBytes
	s, err := SignatureFromBytes(b)
	assert.Nil(t, err)
	assert.Equal(t, b, s.Marshal())

	_, err = SignatureFromBytes([]byte{1, 2, 3, 4})
	assert.NotNil(t, err)

	zero := make([]byte, types.SignatureLength)
	_, err = SignatureFromBytes(zero)
	assert.Equal(t, types.ErrSignatureUnmarshal, err)

	_, err = SignatureFromBytes(testBadSignatureBytes)
	assert.Equal(t, types.ErrSignatureUnmarshal, err)
}

func TestMultipleSignaturesFromBytes(t *testing.T) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)
	bs := tc.sigbs

	// all uncached
	sigs, err := MultipleSignaturesFromBytes(tc.sigbs[:L/2])
	assert.Nil(t, err)
	for i := 0; i < L/2; i++ {
		assert.Equal(t, tc.sigs[i].Marshal(), sigs[i].Marshal())
	}

	// [0:L/2] are cached and [L/2:L] are uncached
	sigs, err = MultipleSignaturesFromBytes(tc.sigbs)
	assert.Nil(t, err)
	for i := 0; i < L; i++ {
		assert.Equal(t, tc.sigs[i].Marshal(), sigs[i].Marshal())
	}

	bs[0] = []byte{1, 2, 3, 4} // len?
	_, err = MultipleSignaturesFromBytes(bs)
	assert.NotNil(t, err)

	bs[0] = make([]byte, types.SignatureLength) // is_inf?
	_, err = MultipleSignaturesFromBytes(bs)
	assert.Equal(t, types.ErrSignatureUnmarshal, err)

	bs[0] = testBadSignatureBytes // in_g2?
	_, err = MultipleSignaturesFromBytes(bs)
	assert.Equal(t, types.ErrSignatureUnmarshal, err)
}

func TestAggregateSignatures(t *testing.T) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	// Correctness check is done in Sign() and Verify() tests
	_, err := AggregateSignatures(tc.sigs)
	assert.Nil(t, err)

	_, err = AggregateSignatures(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)
}

func TestAggregateSignaturesFromBytes(t *testing.T) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)
	bs := tc.sigbs

	// Correctness check is done in Sign() and Verify() tests
	_, err := AggregateSignaturesFromBytes(tc.sigbs)
	assert.Nil(t, err)

	_, err = AggregateSignaturesFromBytes(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)

	bs[0] = []byte{1, 2, 3, 4} // len?
	_, err = AggregateSignaturesFromBytes(bs)
	assert.NotNil(t, err)

	bs[0] = make([]byte, types.SignatureLength) // is_inf?
	_, err = AggregateSignaturesFromBytes(bs)
	assert.Equal(t, types.ErrSignatureUnmarshal, err)

	bs[0] = testBadSignatureBytes // in_g2?
	_, err = AggregateSignaturesFromBytes(bs)
	assert.Equal(t, types.ErrSignatureUnmarshal, err)
}

func TestSignVerify(t *testing.T) {
	sk, _ := RandKey()
	pk := sk.PublicKey()
	msg := testMessage

	sig := Sign(sk, msg)
	assert.NotNil(t, sig)
	assert.True(t, Verify(sig, msg, pk))

	sk2, _ := RandKey()
	pk2 := sk2.PublicKey()
	msg2 := make([]byte, 32)
	assert.False(t, Verify(Sign(sk2, msg), msg, pk))
	assert.False(t, Verify(Sign(sk, msg2), msg, pk))
	assert.False(t, Verify(Sign(sk, msg), msg2, pk))
	assert.False(t, Verify(Sign(sk, msg), msg, pk2))

}
