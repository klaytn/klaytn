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
// This file is derived from core/state/sync.go (2018/06/04).
// Modified and improved for the klaytn development.

package state

import (
	"bytes"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/statedb"
)

// NewStateSync create a new state trie download scheduler.
// LRU cache is mendatory when state syncing and block processing are executed simultaneously
func NewStateSync(root common.Hash, database statedb.StateTrieReadDB, bloom *statedb.SyncBloom, lruCache *lru.Cache) *statedb.TrieSync {
	var syncer *statedb.TrieSync
	callback := func(leaf []byte, parent common.Hash, parentDepth int) error {
		serializer := account.NewAccountSerializer()
		if err := rlp.Decode(bytes.NewReader(leaf), serializer); err != nil {
			return err
		}
		obj := serializer.GetAccount()
		if pa := account.GetProgramAccount(obj); pa != nil {
			syncer.AddSubTrie(pa.GetStorageRoot(), parentDepth+1, parent, nil)
			syncer.AddRawEntry(common.BytesToHash(pa.GetCodeHash()), parentDepth+1, parent)
		}
		return nil
	}
	syncer = statedb.NewTrieSync(root, database, callback, bloom, lruCache)
	return syncer
}
