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
	testSignatureBytes = common.FromHex("0xae82747ddeefe4fd64cf9cedb9b04ae3e8a43420cd255e3c7cd06a8d88b7c7f8638543719981c5d16fa3527c468c25f0026704a6951bde891360c7e8d12ddee0559004ccdbe6046b55bae1b257ee97f7cdb955773d7cf29adf3ccbb9975e4eb9")
	// deserialization_G2/deserialization_fails_not_in_curve.json
	testBadSignatureBytes = common.FromHex("0x8123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde0")
	// fast_aggregate_verify/fast_aggregate_verify_valid_3d7576f3c0e3570a.json
	testAggPublicKeyBytes1 = common.FromHex("0xa491d1b0ecd9bb917989f0e74f0dea0422eac4a873e5e2644f368dffb9a6e20fd6e10c1b77654d067c0618f6e5a7f79a")
	testAggPublicKeyBytes2 = common.FromHex("0xb301803f8b5ac4a1133581fc676dfedc60d891dd5fa99028805e5ea5b08d3491af75d0707adab3b70c6a6a580217bf81")
	testAggPublicKeyBytes3 = common.FromHex("0xb53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
	testMessage            = common.FromHex("0xabababababababababababababababababababababababababababababababab")
	testAggSignatureBytes  = common.FromHex("0x9712c3edd73a209c742b8250759db12549b3eaf43b5ca61376d9f30e2747dbcf842d8b2ac0901d2a093713e20284a7670fcf6954e9ab93de991bb9b313e664785a075fc285806fa5224c82bde146561b446ccfc706a64b8579513cfc4ff1d930")
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

func TestAggregateVerify(t *testing.T) {
	apk, _ := AggregatePublicKeysFromBytes([][]byte{
		testAggPublicKeyBytes1,
		testAggPublicKeyBytes2,
		testAggPublicKeyBytes3,
	})
	msg := testMessage
	sig, _ := SignatureFromBytes(testAggSignatureBytes)

	assert.True(t, Verify(sig, msg, apk))
}
