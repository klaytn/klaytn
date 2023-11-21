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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/klaytn/klaytn/blockchain/system"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/log"
	uuid "github.com/satori/go.uuid"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/crypto/pbkdf2"
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

func GenerateKeysFromMnemonic(num int, mnemonic, path string) (keys []*ecdsa.PrivateKey, nodekeys []string, addrs []common.Address) {
	var key *bip32.Key

	for _, level := range strings.Split(path, "/") {
		if len(level) == 0 {
			continue
		}

		if level == "m" {
			seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)
			key, _ = bip32.NewMasterKey(seed)
		} else {
			num := uint64(0)
			var err error
			if strings.HasSuffix(level, "'") {
				num, err = strconv.ParseUint(level[:len(level)-1], 10, 32)
				num += 0x80000000
			} else {
				num, err = strconv.ParseUint(level, 10, 32)
			}
			if err != nil {
				logger.Error("Failed to parse path", "err", err)
				return nil, nil, nil
			}
			key, _ = key.NewChildKey(uint32(num))
		}
	}

	for i := 0; i < num; i++ {
		derived, _ := key.NewChildKey(uint32(i))
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

func GenerateKip113Init(privKeys []*ecdsa.PrivateKey, owner common.Address) system.AllocKip113Init {
	init := system.AllocKip113Init{}
	init.Infos = make(map[common.Address]system.BlsPublicKeyInfo)

	for i, key := range privKeys {
		blsKey, err := bls.GenerateKey(crypto.FromECDSA(privKeys[i]))
		if err != nil {
			logger.Error("Failed to generate bls key", "err", err)
			continue
		}
		addr := crypto.PubkeyToAddress(key.PublicKey)
		init.Infos[addr] = system.BlsPublicKeyInfo{
			PublicKey: blsKey.PublicKey().Marshal(),
			Pop:       bls.PopProve(blsKey).Marshal(),
		}
	}

	init.Owner = owner

	return init
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
	data, err := os.ReadFile(src)
	if err != nil {
		logger.Error("Failed to read file", "file", src, "err", err)
		return
	}
	err = os.WriteFile(dst, data, 0o644)
	if err != nil {
		logger.Error("Failed to write file", "file", dst, "err", err)
		return
	}
}
