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
)

func TestSignatureCopy(t *testing.T) {
	sig, _ := SignatureFromBytes(testSignatureBytes)
	assert.Equal(t, sig.Marshal(), sig.Copy().Marshal())
}

func TestSignatureFromBytes(t *testing.T) {
	b := testSignatureBytes
	s, err := SignatureFromBytes(b)
	assert.Nil(t, err)
	assert.Equal(t, b, s.Marshal())

	_, err = SignatureFromBytes([]byte{1, 2, 3, 4}) // len?
	assert.NotNil(t, err)

	zero := make([]byte, types.SignatureLength)
	_, err = SignatureFromBytes(zero) // is_inf?
	assert.Equal(t, types.ErrSignatureUnmarshal, err)

	_, err = SignatureFromBytes(testBadSignatureBytes) // in_g2?
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

	sigs, err = MultipleSignaturesFromBytes(nil) // empty
	assert.Nil(t, err)
	assert.Nil(t, sigs)

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
	var (
		// https://github.com/ethereum/bls12-381-tests
		// aggregate/aggregate_0xabababababababababababababababababababababababababababababababab.json
		sigb1 = common.FromHex("0x91347bccf740d859038fcdcaf233eeceb2a436bcaaee9b2aa3bfb70efe29dfb2677562ccbea1c8e061fb9971b0753c240622fab78489ce96768259fc01360346da5b9f579e5da0d941e4c6ba18a0e64906082375394f337fa1af2b7127b0d121")
		sigb2 = common.FromHex("0x9674e2228034527f4c083206032b020310face156d4a4685e2fcaec2f6f3665aa635d90347b6ce124eb879266b1e801d185de36a0a289b85e9039662634f2eea1e02e670bc7ab849d006a70b2f93b84597558a05b879c8d445f387a5d5b653df")
		sigb3 = common.FromHex("0xae82747ddeefe4fd64cf9cedb9b04ae3e8a43420cd255e3c7cd06a8d88b7c7f8638543719981c5d16fa3527c468c25f0026704a6951bde891360c7e8d12ddee0559004ccdbe6046b55bae1b257ee97f7cdb955773d7cf29adf3ccbb9975e4eb9")
		asigb = common.FromHex("0x9712c3edd73a209c742b8250759db12549b3eaf43b5ca61376d9f30e2747dbcf842d8b2ac0901d2a093713e20284a7670fcf6954e9ab93de991bb9b313e664785a075fc285806fa5224c82bde146561b446ccfc706a64b8579513cfc4ff1d930")

		sig1, _ = SignatureFromBytes(sigb1)
		sig2, _ = SignatureFromBytes(sigb2)
		sig3, _ = SignatureFromBytes(sigb3)
	)

	asig, err := AggregateSignatures([]types.Signature{sig1, sig2, sig3})
	assert.Nil(t, err)
	assert.Equal(t, asigb, asig.Marshal())

	_, err = AggregateSignatures(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)
}

func TestAggregateSignaturesFromBytes(t *testing.T) {
	var (
		// https://github.com/ethereum/bls12-381-tests
		// aggregate/aggregate_0xabababababababababababababababababababababababababababababababab.json
		sigb1 = common.FromHex("0x91347bccf740d859038fcdcaf233eeceb2a436bcaaee9b2aa3bfb70efe29dfb2677562ccbea1c8e061fb9971b0753c240622fab78489ce96768259fc01360346da5b9f579e5da0d941e4c6ba18a0e64906082375394f337fa1af2b7127b0d121")
		sigb2 = common.FromHex("0x9674e2228034527f4c083206032b020310face156d4a4685e2fcaec2f6f3665aa635d90347b6ce124eb879266b1e801d185de36a0a289b85e9039662634f2eea1e02e670bc7ab849d006a70b2f93b84597558a05b879c8d445f387a5d5b653df")
		sigb3 = common.FromHex("0xae82747ddeefe4fd64cf9cedb9b04ae3e8a43420cd255e3c7cd06a8d88b7c7f8638543719981c5d16fa3527c468c25f0026704a6951bde891360c7e8d12ddee0559004ccdbe6046b55bae1b257ee97f7cdb955773d7cf29adf3ccbb9975e4eb9")
		asigb = common.FromHex("0x9712c3edd73a209c742b8250759db12549b3eaf43b5ca61376d9f30e2747dbcf842d8b2ac0901d2a093713e20284a7670fcf6954e9ab93de991bb9b313e664785a075fc285806fa5224c82bde146561b446ccfc706a64b8579513cfc4ff1d930")
	)

	asig, err := AggregateSignaturesFromBytes([][]byte{sigb1, sigb2, sigb3})
	assert.Nil(t, err)
	assert.Equal(t, asigb, asig.Marshal())

	_, err = AggregateSignaturesFromBytes(nil) // empty
	assert.Equal(t, types.ErrEmptyArray, err)

	bs := make([][]byte, 1)
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
	var (
		// https://github.com/ethereum/bls12-381-tests
		// sign/sign_case_84d45c9c7cca6b92.json
		// verify/verify_valid_case_195246ee3bd3b6ec.json
		skb  = common.FromHex("0x328388aff0d4a5b7dc9205abd374e7e98f3cd9f3418edb4eafda5fb16473d216")
		msg  = common.FromHex("0xabababababababababababababababababababababababababababababababab")
		pkb  = common.FromHex("0xb53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
		sigb = common.FromHex("0xae82747ddeefe4fd64cf9cedb9b04ae3e8a43420cd255e3c7cd06a8d88b7c7f8638543719981c5d16fa3527c468c25f0026704a6951bde891360c7e8d12ddee0559004ccdbe6046b55bae1b257ee97f7cdb955773d7cf29adf3ccbb9975e4eb9")

		sk, _  = SecretKeyFromBytes(skb)
		pk, _  = PublicKeyFromBytes(pkb)
		sig, _ = SignatureFromBytes(sigb)
	)

	assert.Equal(t, sigb, Sign(sk, msg).Marshal())
	assert.True(t, Verify(sig, msg, pk))

	sk2, _ := RandKey()
	pk2 := sk2.PublicKey()
	msg2 := []byte("test message")
	assert.True(t, Verify(Sign(sk2, msg2), msg2, pk2))
	assert.False(t, Verify(Sign(sk2, msg), msg, pk))
	assert.False(t, Verify(Sign(sk, msg2), msg, pk))
	assert.False(t, Verify(Sign(sk, msg), msg2, pk))
	assert.False(t, Verify(Sign(sk, msg), msg, pk2))
}

func TestPopProveVerify(t *testing.T) {
	var (
		// https://github.com/Chia-Network/bls-signatures/blob/2.0.2/src/test.cpp
		// "Chia test vector 3 (PoP)"
		// Note: `skb` is the result of PopSchemeMPL().KeyGen(seed1) in the testcase
		skb  = common.FromHex("0x258787ef728c898e43bc76244d70f468c9c7e1338a107b18b42da0d86b663c26")
		popb = common.FromHex("0x84f709159435f0dc73b3e8bf6c78d85282d19231555a8ee3b6e2573aaf66872d9203fefa1ef700e34e7c3f3fb28210100558c6871c53f1ef6055b9f06b0d1abe22ad584ad3b957f3018a8f58227c6c716b1e15791459850f2289168fa0cf9115")

		sk, _ = SecretKeyFromBytes(skb)
		pk    = sk.PublicKey()
		pop   = PopProve(sk)
	)
	assert.Equal(t, popb, pop.Marshal())
	assert.True(t, PopVerify(pk, pop))

	sk2, _ := RandKey()
	pk2 := sk2.PublicKey()
	assert.True(t, PopVerify(pk2, PopProve(sk2)))
	assert.False(t, PopVerify(pk2, PopProve(sk)))
	assert.False(t, PopVerify(pk, PopProve(sk2)))
}

func TestAggregateVerify(t *testing.T) {
	var (
		// https://github.com/ethereum/bls12-381-tests
		// fast_aggregate_verify/fast_aggregate_verify_valid_3d7576f3c0e3570a.json
		pkb1  = common.FromHex("0xa491d1b0ecd9bb917989f0e74f0dea0422eac4a873e5e2644f368dffb9a6e20fd6e10c1b77654d067c0618f6e5a7f79a")
		pkb2  = common.FromHex("0xb301803f8b5ac4a1133581fc676dfedc60d891dd5fa99028805e5ea5b08d3491af75d0707adab3b70c6a6a580217bf81")
		pkb3  = common.FromHex("0xb53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
		msg   = common.FromHex("0xabababababababababababababababababababababababababababababababab")
		asigb = common.FromHex("0x9712c3edd73a209c742b8250759db12549b3eaf43b5ca61376d9f30e2747dbcf842d8b2ac0901d2a093713e20284a7670fcf6954e9ab93de991bb9b313e664785a075fc285806fa5224c82bde146561b446ccfc706a64b8579513cfc4ff1d930")

		sig, _ = SignatureFromBytes(asigb)
	)

	apk, _ := AggregatePublicKeysFromBytes([][]byte{pkb1, pkb2, pkb3})
	assert.True(t, Verify(sig, msg, apk))
}

func TestVerifyMultiple(t *testing.T) {
	var (
		// https://github.com/ethereum/bls12-381-tests
		// batch_verify/batch_verify_valid_simple_signature_set.json
		pkb1  = common.FromHex("0xa491d1b0ecd9bb917989f0e74f0dea0422eac4a873e5e2644f368dffb9a6e20fd6e10c1b77654d067c0618f6e5a7f79a")
		pkb2  = common.FromHex("0xb301803f8b5ac4a1133581fc676dfedc60d891dd5fa99028805e5ea5b08d3491af75d0707adab3b70c6a6a580217bf81")
		pkb3  = common.FromHex("0xb53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
		msg1  = common.BytesToHash(common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000000"))
		msg2  = common.BytesToHash(common.FromHex("0x5656565656565656565656565656565656565656565656565656565656565656"))
		msg3  = common.BytesToHash(common.FromHex("0xabababababababababababababababababababababababababababababababab"))
		sigb1 = common.FromHex("0xb6ed936746e01f8ecf281f020953fbf1f01debd5657c4a383940b020b26507f6076334f91e2366c96e9ab279fb5158090352ea1c5b0c9274504f4f0e7053af24802e51e4568d164fe986834f41e55c8e850ce1f98458c0cfc9ab380b55285a55")
		sigb2 = common.FromHex("0xaf1390c3c47acdb37131a51216da683c509fce0e954328a59f93aebda7e4ff974ba208d9a4a2a2389f892a9d418d618418dd7f7a6bc7aa0da999a9d3a5b815bc085e14fd001f6a1948768a3f4afefc8b8240dda329f984cb345c6363272ba4fe")
		sigb3 = common.FromHex("0xae82747ddeefe4fd64cf9cedb9b04ae3e8a43420cd255e3c7cd06a8d88b7c7f8638543719981c5d16fa3527c468c25f0026704a6951bde891360c7e8d12ddee0559004ccdbe6046b55bae1b257ee97f7cdb955773d7cf29adf3ccbb9975e4eb9")

		sigbs   = [][]byte{sigb1, sigb2, sigb3}
		msgs    = [][32]byte{msg1, msg2, msg3}
		pks, _  = MultiplePublicKeysFromBytes([][]byte{pkb1, pkb2, pkb3})
		sigs, _ = MultipleSignaturesFromBytes(sigbs)
	)

	ok, err := VerifyMultipleSignatures(sigbs, msgs, pks)
	assert.Nil(t, err)
	assert.True(t, ok)

	// Verify individually
	assert.True(t, Verify(sigs[0], msgs[0][:], pks[0]))
	assert.True(t, Verify(sigs[1], msgs[1][:], pks[1]))
	assert.True(t, Verify(sigs[2], msgs[2][:], pks[2]))
}
