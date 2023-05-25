// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/types/derive_sha.go (2018/06/04).
// Modified and improved for the klaytn development.

package derivesha

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto/sha3"
)

// An alternative implementation of DeriveSha()
// This function generates a hash of `DerivableList` as below:
// 1. make a byte slice by concatenating RLP-encoded items
// 2. make a hash of the byte slice.
type DeriveShaConcat struct{}

func (d DeriveShaConcat) DeriveSha(list types.DerivableList) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	for i := 0; i < list.Len(); i++ {
		hasher.Write(list.GetRlp(i))
	}
	hasher.Sum(hash[:0])

	return hash
}
