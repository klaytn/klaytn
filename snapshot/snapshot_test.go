// Modifications Copyright 2021 The klaytn Authors
// Copyright 2019 The go-ethereum Authors
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
//
// This file is derived from core/state/snapshot/snapshot_test.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"math/big"
	"math/rand"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
)

// randomHash generates a random blob of data and returns it as a hash.
func randomHash() common.Hash {
	var hash common.Hash
	if n, err := rand.Read(hash[:]); n != common.HashLength || err != nil {
		panic(err)
	}
	return hash
}

// randomAccount generates a random account and returns it RLP encoded.
func randomAccount() []byte {
	root := randomHash()
	a := Account{
		Balance:  big.NewInt(rand.Int63()),
		Nonce:    rand.Uint64(),
		Root:     root[:],
		CodeHash: emptyCode[:],
	}
	data, _ := rlp.EncodeToBytes(a)
	return data
}
