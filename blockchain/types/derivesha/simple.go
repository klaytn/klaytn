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
// This function generates a hash of `DerivableList` by simulating merkle tree generation
type DeriveShaSimple struct{}

func (d DeriveShaSimple) DeriveSha(list types.DerivableList) common.Hash {
	hasher := sha3.NewKeccak256()

	encoded := make([][]byte, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		hasher.Reset()
		hasher.Write(list.GetRlp((i)))
		encoded = append(encoded, hasher.Sum(nil))
	}

	for len(encoded) > 1 {
		// make even numbers
		if len(encoded)%2 == 1 {
			encoded = append(encoded, encoded[len(encoded)-1])
		}

		for i := 0; i < len(encoded)/2; i++ {
			hasher.Reset()
			hasher.Write(encoded[2*i])
			hasher.Write(encoded[2*i+1])

			encoded[i] = hasher.Sum(nil)
		}

		encoded = encoded[0 : len(encoded)/2]
	}

	if len(encoded) == 0 {
		hasher.Reset()
		hasher.Write(nil)
		return common.BytesToHash(hasher.Sum(nil))
	}

	return common.BytesToHash(encoded[0])
}
