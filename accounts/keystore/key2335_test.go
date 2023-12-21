// Copyright 2018 The klaytn Authors
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

package keystore

import (
	"os"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptEIP2335(t *testing.T) {
	password := "password"
	sk, err := bls.RandKey()
	require.Nil(t, err)

	plain := NewKeyEIP2335(sk)

	encrypted, err := EncryptKeyEIP2335(plain, password, LightScryptN, LightScryptP)
	require.Nil(t, err)

	decrypted, err := DecryptKeyEIP2335(encrypted, password)
	require.Nil(t, err)

	assert.Equal(t, plain.ID, decrypted.ID)
	assert.Equal(t, plain.SecretKey.Marshal(), decrypted.SecretKey.Marshal())
	assert.Equal(t, plain.PublicKey.Marshal(), decrypted.PublicKey.Marshal())
}

func TestDecryptEIP2335(t *testing.T) {
	var (
		// https://eips.ethereum.org/EIPS/eip-2335 test vectors
		passwordBytes, _ = os.ReadFile("testdata/eip2335_password.txt")
		password         = strings.TrimSpace(string(passwordBytes))
		keyBytes         = hexutil.MustDecode("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")

		scryptJSON, _ = os.ReadFile("testdata/eip2335_scrypt.json")
		pbkdf2JSON, _ = os.ReadFile("testdata/eip2335_pbkdf2.json")
	)

	k, err := DecryptKeyEIP2335(scryptJSON, password)
	require.Nil(t, err)
	assert.Equal(t, keyBytes, k.SecretKey.Marshal())

	k, err = DecryptKeyEIP2335(pbkdf2JSON, password)
	require.Nil(t, err)
	assert.Equal(t, keyBytes, k.SecretKey.Marshal())
}
