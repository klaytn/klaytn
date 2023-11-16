// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// Modified by Prysmatic Labs 2018
// Modified by the klaytn Authors 2023
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package keystore

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/pborman/uuid"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

// KeyEIP2335 is a decrypted BLS12-381 keypair.
type KeyEIP2335 struct {
	ID uuid.UUID // Version 4 "random" for unique id not derived from key data

	PublicKey bls.PublicKey // Represents the public key of the user.
	SecretKey bls.SecretKey // Represents the private key of the user.
}

// https://eips.ethereum.org/EIPS/eip-2335
// Required fields: crypto, path, uuid, version
type encryptedKeyEIP2335JSON struct {
	Crypto      map[string]interface{} `json:"crypto"`
	Description string                 `json:"description"`
	Pubkey      string                 `json:"pubkey"`
	Path        string                 `json:"path"`
	ID          string                 `json:"uuid"`
	Version     int                    `json:"version"`
}

// NewKeyEIP2335 creates a new EIP-2335 keystore Key type using a BLS private key.
func NewKeyEIP2335(blsKey bls.SecretKey) *KeyEIP2335 {
	return &KeyEIP2335{
		ID:        uuid.NewRandom(),
		PublicKey: blsKey.PublicKey(),
		SecretKey: blsKey,
	}
}

// DecryptKeyEIP2335 decrypts a key from an EIP-2335 JSON blob, returning the BLS private key.
func DecryptKeyEIP2335(keyJSON []byte, password string) (*KeyEIP2335, error) {
	k := new(encryptedKeyEIP2335JSON)
	if err := json.Unmarshal(keyJSON, k); err != nil {
		return nil, err
	}

	id := uuid.Parse(k.ID)
	if id == nil {
		return nil, errors.New("Invalid UUID")
	}

	decryptor := keystorev4.New()
	keyBytes, err := decryptor.Decrypt(k.Crypto, password)
	if err != nil {
		return nil, err
	}

	secretKey, err := bls.SecretKeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return &KeyEIP2335{
		ID:        id,
		PublicKey: secretKey.PublicKey(),
		SecretKey: secretKey,
	}, nil
}

// EncryptKeyEIP2335 encrypts a BLS key using the specified scrypt parameters into a JSON
// blob that can be decrypted later on.
func EncryptKeyEIP2335(key *KeyEIP2335, password string, scryptN, scryptP int) ([]byte, error) {
	keyBytes := key.SecretKey.Marshal()
	encryptor := keystorev4.New()
	cryptoObj, err := encryptor.Encrypt(keyBytes, password)
	if err != nil {
		return nil, err
	}

	encryptedJSON := encryptedKeyEIP2335JSON{
		Crypto:      cryptoObj,
		Description: "",
		Pubkey:      hex.EncodeToString(key.PublicKey.Marshal()),
		Path:        "", // EIP-2335: if no path is known or the path is not relevant, the empty string, "" indicates this
		ID:          key.ID.String(),
		Version:     4,
	}
	return json.Marshal(encryptedJSON)
}
