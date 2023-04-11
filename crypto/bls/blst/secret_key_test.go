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
	// https://github.com/Wind4Greg/BBS-Draft-Checks/blob/main/JavaScript/BBSPubKey.js
	testSecretKeyBytes = common.FromHex("4a39afffd624d69e81808b2e84385cc80bf86adadf764e030caa46c231f2a8d7")
	testPublicKeyBytes = common.FromHex("aaff983278257afc45fa9d44d156c454d716fb1a250dfed132d65b2009331f618c623c14efa16245f50cc92e60334051087f1ae92669b89690f5feb92e91568f95a8e286d110b011e9ac9923fd871238f57d1295395771331ff6edee43e4ccc6")
)

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

	pk := sk.PublicKey()
	assert.Equal(t, testPublicKeyBytes, pk.Marshal())
}
