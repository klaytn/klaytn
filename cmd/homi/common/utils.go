// Copyright 2017 AMIS Technologies
// This file is part of the go-ethereum library.
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

package common

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/log"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/crypto/pbkdf2"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultLocalDir  = "/tmp/kdata"
	clientIdentifier = "klay"
)

var logger = log.NewModuleLogger(log.CMDIstanbul)

func GenerateRandomDir() (string, error) {
	err := os.MkdirAll(filepath.Join(defaultLocalDir), 0o700)
	if err != nil {
		logger.Error("Failed to create dir", "dir", defaultLocalDir, "err", err)
		return "", err
	}

	instanceDir := filepath.Join(defaultLocalDir, fmt.Sprintf("%s-%s", clientIdentifier, uuid.NewV4().String()))
	if err := os.MkdirAll(instanceDir, 0o700); err != nil {
		logger.Error("Failed to create dir", "dir", instanceDir, "err", err)
		return "", err
	}

	return instanceDir, nil
}

func GenerateKeys(num int) (keys []*ecdsa.PrivateKey, nodekeys []string, addrs []common.Address) {
	for i := 0; i < num; i++ {
		nodekey := RandomHex()[2:]
		nodekeys = append(nodekeys, nodekey)

		key, err := crypto.HexToECDSA(nodekey)
		if err != nil {
			logger.Error("Failed to generate key", "err", err)
			return nil, nil, nil
		}
		keys = append(keys, key)

		addr := crypto.PubkeyToAddress(key.PublicKey)
		addrs = append(addrs, addr)
	}

	return keys, nodekeys, addrs
}

func GenerateKeysFromMnemonic(num int, mnemonic string) (keys []*ecdsa.PrivateKey, nodekeys []string, addrs []common.Address) {
	// Ethereum key derivation path: m/44'/60'/0'/0/
	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)
	mk, _ := bip32.NewMasterKey(seed)
	m44h, _ := mk.NewChildKey(0x80000000 + 44)
	m44h_60h, _ := m44h.NewChildKey(0x80000000 + 60)
	m44h_60h_0h, _ := m44h_60h.NewChildKey(0x80000000)
	m44h_60h_0h_0, _ := m44h_60h_0h.NewChildKey(0)

	for i := 0; i < num; i++ {
		derived, _ := m44h_60h_0h_0.NewChildKey(uint32(i))
		nodekey := hexutil.Encode(derived.Key)[2:]
		nodekeys = append(nodekeys, nodekey)

		key, err := crypto.HexToECDSA(nodekey)
		if err != nil {
			logger.Error("Failed to generate key", "err", err)
			return nil, nil, nil
		}
		keys = append(keys, key)

		addr := crypto.PubkeyToAddress(key.PublicKey)
		addrs = append(addrs, addr)
	}

	return keys, nodekeys, addrs
}

func RandomHex() string {
	b, _ := RandomBytes(32)
	return common.BytesToHash(b).Hex()
}

func RandomBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	_, _ = rand.Read(b)

	return b, nil
}

func copyFile(src string, dst string) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		logger.Error("Failed to read file", "file", src, "err", err)
		return
	}
	err = ioutil.WriteFile(dst, data, 0o644)
	if err != nil {
		logger.Error("Failed to write file", "file", dst, "err", err)
		return
	}
}
