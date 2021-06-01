// Copyright 2019 The klaytn Authors
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

package crypto

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
)

// BenchmarkCreateAddress measures performance of two address generation methods:
// createAddressUsingCode: hash(address, nonce, code)
// createAddressUsingCodeHash: hash(address, nonce, codeHash)
//
// This benchmark is created to measure performance difference of hashing twice with shorter data
// and hashing once with long data.
func BenchmarkCreateAddress(b *testing.B) {
	addr := common.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d")
	nonce := uint64(1023)

	fns := []func(addr common.Address, nonce uint64, code []byte) common.Address{
		createAddressUsingCode,
		createAddressUsingCodeHash,
	}

	const MaxCodeSize uint = 10 * 1024

	for codeSize := uint(0); codeSize <= MaxCodeSize; codeSize += 640 {
		code := randCode(codeSize)

		for _, f := range fns {
			fname := getFunctionName(f)
			fname = fname[strings.LastIndex(fname, ".")+1:]

			benchName := fmt.Sprintf("%s/%d", fname, codeSize)

			b.Run(benchName, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = f(addr, nonce, code)
				}
			})
		}
	}
}

func randCode(codeSize uint) []byte {
	code := make([]byte, 0, codeSize)
	rand.Seed(0)
	r := rand.Uint64()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, r)
	code = append(code, b...)

	return code
}

func createAddressUsingCodeHash(addr common.Address, nonce uint64, code []byte) common.Address {
	codeHash := Keccak256Hash(code)

	data, _ := rlp.EncodeToBytes(struct {
		Addr     common.Address
		nonce    uint64
		codeHash common.Hash
	}{addr, nonce, codeHash})

	return common.BytesToAddress(Keccak256(data)[12:])
}

func createAddressUsingCode(addr common.Address, nonce uint64, code []byte) common.Address {
	data, _ := rlp.EncodeToBytes(struct {
		Addr     common.Address
		nonce    uint64
		codeHash []byte
	}{addr, nonce, code})

	return common.BytesToAddress(Keccak256(data)[12:])
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
