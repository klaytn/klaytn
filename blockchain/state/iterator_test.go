// Copyright 2016 The go-ethereum Authors
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

// This file is derived from ethdb/iterator_test.go (2020/05/20).
// Modified and improved for the klaytn development.

package state

import (
	"bytes"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/storage/database"
)

// Tests that the node iterator indeed walks over the entire database contents.
func TestNodeIteratorCoverage(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlCrit)
	// Create some arbitrary test state to iterate
	db, root, _ := makeTestState(t)
	db.TrieDB().Commit(root, false, 0)

	state, err := New(root, db, nil, nil)
	if err != nil {
		t.Fatalf("failed to create state trie at %x: %v", root, err)
	}
	// Gather all the node hashes found by the iterator
	iterated := make(map[common.Hash]struct{})
	for it := NewNodeIterator(state); it.Next(); {
		if it.Hash != (common.Hash{}) {
			iterated[it.Hash] = struct{}{}
		}
	}

	// Cross check the iterated hashes and the database/nodepool content
	// (TrieDB.Node + ContractCode) contains all iterated hashes
	for itHash := range iterated {
		_, err1 := db.TrieDB().Node(itHash.ExtendZero())
		_, err2 := db.ContractCode(itHash)
		if err1 != nil && err2 != nil { // both failed
			t.Errorf("failed to retrieve reported node %x", itHash)
		}
	}
	// iterated hashes contains all TrieDB.Nodes
	for _, exthash := range db.TrieDB().Nodes() {
		hash := exthash.Unextend()
		if hash == emptyCode {
			continue // skip emptyCode
		}
		if _, ok := iterated[hash]; !ok {
			t.Errorf("state entry not reported %x", hash)
		}
	}
	// iterated hashes contains all DiskDB keys
	// StateTrieDB contains preimages, codes and nodes
	it := db.TrieDB().DiskDB().GetMemDB().NewIterator(nil, nil)
	for it.Next() {
		key := it.Key()
		if bytes.HasPrefix(key, []byte("secure-key-")) {
			continue // skip preimages
		}
		if isCode, _ := database.IsCodeKey(key); isCode {
			continue // skip codes
		}
		hash := common.BytesToExtHash(key).Unextend()
		if hash == emptyCode {
			continue // skip emptyCode
		}
		if _, ok := iterated[hash]; !ok {
			t.Errorf("state entry not reported %x", hash)
		}
	}
	it.Release()
}
