// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/statedb"
)

type DeriveShaOrig struct{}

func (d DeriveShaOrig) DeriveSha(list types.DerivableList) common.Hash {
	trie := statedb.NewStackTrie(nil)
	trie.Reset()
	var buf []byte

	// StackTrie requires values to be inserted in increasing
	// hash order, which is not the order that `list` provides
	// hashes in. This insertion sequence ensures that the
	// order is correct.
	for i := 1; i < list.Len() && i <= 0x7f; i++ {
		buf = rlp.AppendUint64(buf[:0], uint64(i))
		trie.Update(buf, list.GetRlp(i))
	}
	if list.Len() > 0 {
		buf = rlp.AppendUint64(buf[:0], 0)
		trie.Update(buf, list.GetRlp(0))
	}
	for i := 0x80; i < list.Len(); i++ {
		buf = rlp.AppendUint64(buf[:0], uint64(i))
		trie.Update(buf, list.GetRlp(i))
	}
	return trie.Hash()
}
