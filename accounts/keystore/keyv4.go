// Copyright 2020 The klaytn Authors
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
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/pborman/uuid"
)

type KeyV4 struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// We only store privkey as pubkey/address can be derived from it.
	PrivateKeys [][]*ecdsa.PrivateKey
}

type plainKeyJSONV4 struct {
	Address     string     `json:"address"`
	PrivateKeys [][]string `json:"privatekeys"`
	Id          string     `json:"id"`
	Version     int        `json:"version"`
}

type encryptedKeyJSONV4 struct {
	Address string         `json:"address"`
	Keyring [][]cryptoJSON `json:"keyring"`
	Id      string         `json:"id"`
	Version int            `json:"version"`
}

type encryptedKeyJSONV4Single struct {
	Address string       `json:"address"`
	Keyring []cryptoJSON `json:"keyring"`
	Id      string       `json:"id"`
	Version int          `json:"version"`
}

func (k *KeyV4) MarshalJSON() (j []byte, err error) {
	privateKeys := make([][]string, len(k.PrivateKeys))

	for i, keys := range k.PrivateKeys {
		privateKeys[i] = make([]string, len(keys))
		for j, key := range keys {
			privateKeys[i][j] = common.Bytes2Hex(crypto.FromECDSA(key))
		}
	}

	jStruct := plainKeyJSONV4{
		Address:     hex.EncodeToString(k.Address.Bytes()),
		PrivateKeys: privateKeys,
		Id:          k.Id.String(),
		Version:     4,
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *KeyV4) UnmarshalJSON(j []byte) (err error) {
	keyJSON := new(plainKeyJSONV4)
	if err := json.Unmarshal(j, keyJSON); err != nil {
		return err
	}

	k.Id = uuid.Parse(keyJSON.Id)

	addr, err := hex.DecodeString(keyJSON.Address)
	if err != nil {
		return err
	}

	k.PrivateKeys = make([][]*ecdsa.PrivateKey, len(keyJSON.PrivateKeys))
	for i, keys := range keyJSON.PrivateKeys {
		k.PrivateKeys[i] = make([]*ecdsa.PrivateKey, len(keys))
		for j, key := range keys {
			k.PrivateKeys[i][j], err = crypto.HexToECDSA(key)
			if err != nil {
				return err
			}
		}
	}

	k.Address = common.BytesToAddress(addr)

	return nil
}

func (k *KeyV4) GetId() uuid.UUID {
	return k.Id
}

func (k *KeyV4) GetAddress() common.Address {
	return k.Address
}

func (k *KeyV4) GetPrivateKey() *ecdsa.PrivateKey {
	if len(k.PrivateKeys) == 0 || len(k.PrivateKeys[0]) == 0 {
		return nil
	}

	return k.PrivateKeys[0][0]
}

func (k *KeyV4) GetPrivateKeys() [][]*ecdsa.PrivateKey {
	return k.PrivateKeys
}

func (k *KeyV4) GetPrivateKeysWithRole(role int) []*ecdsa.PrivateKey {
	if len(k.PrivateKeys) <= role {
		return nil
	}
	return k.PrivateKeys[role]
}

func (k *KeyV4) ResetPrivateKey() {
	for _, keys := range k.PrivateKeys {
		for _, key := range keys {
			zeroKey(key)
		}
	}
}
