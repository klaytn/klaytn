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

// https://github.com/ethereum/bls12-381-tests
// sign/sign_case_142f678a8d05fcd1.json
var testSecretKeyBytes = common.FromHex("0x47b8192d77bf871b62e87859d653922725724a5c031afeabc60bcef5ff665138")

func TestRandKey(t *testing.T) {
	sk, err := RandKey()
	assert.Nil(t, err)
	assert.Equal(t, types.SecretKeyLength, len(sk.Marshal()))
}

func TestSecretKeyFromBytes(t *testing.T) {
	b := testSecretKeyBytes
	sk, err := SecretKeyFromBytes(b)
	assert.Nil(t, err)
	assert.Equal(t, b, sk.Marshal())

	_, err = SecretKeyFromBytes([]byte{1, 2, 3, 4})
	assert.NotNil(t, err)

	// Valid secret key must be 1 <= SK < r
	// as per draft-irtf-cfrg-bls-signature-05#2.3. KeyGen
	zero := make([]byte, types.SecretKeyLength)
	_, err = SecretKeyFromBytes(zero)
	assert.Equal(t, types.ErrSecretKeyUnmarshal, err)

	order := common.FromHex(types.CurveOrderHex)
	_, err = SecretKeyFromBytes(order)
	assert.Equal(t, types.ErrSecretKeyUnmarshal, err)
}

func TestSecretKeyPublicKey(t *testing.T) {
	sk, err := SecretKeyFromBytes(testSecretKeyBytes)
	assert.Nil(t, err)

	// Correctness check is done in Sign() and Verify() tests
	pk := sk.PublicKey()
	assert.NotNil(t, pk)
}
