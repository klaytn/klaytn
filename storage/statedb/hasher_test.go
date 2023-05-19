// Copyright 2023 The klaytn Authors
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

package statedb

import (
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

func checkHasherHashFunc(t *testing.T, idx int, tc *testNodeEncodingTC, hashFunc func(*Database) (node, node)) {
	hash := common.BytesToHash(tc.hash)

	memDB := database.NewMemoryDBManager()
	db := NewDatabase(memDB)

	hashed, cached := hashFunc(db)
	t.Logf("tc[%02d] %s", idx, hashed)
	assert.Equal(t, hashNode(tc.hash), hashed, idx)

	cachedHash, _ := cached.cache()
	assert.Equal(t, hashNode(tc.hash), cachedHash, idx)

	inserted := db.nodes[hash].node
	assert.Equal(t, tc.inserted, inserted, idx)

	db.Cap(0)
	encoded, _ := memDB.ReadCachedTrieNode(hash)
	assert.Equal(t, tc.encoded, encoded, idx)
}

func checkHasherHash(t *testing.T, idx int, tc *testNodeEncodingTC) {
	checkHasherHashFunc(t, idx, tc, func(db *Database) (node, node) {
		h := newHasher(nil)
		defer returnHasherToPool(h)
		return h.hash(tc.expanded, db, false)
	})

	checkHasherHashFunc(t, idx, tc, func(db *Database) (node, node) {
		h := newHasher(nil)
		defer returnHasherToPool(h)
		return h.hashRoot(tc.expanded, db, false)
	})
}

func TestHasherHashTC(t *testing.T) {
	for idx, tc := range collapsedNodeTCs() {
		checkHasherHash(t, idx, tc)
	}
	for idx, tc := range resolvedNodeTCs() {
		checkHasherHash(t, idx, tc)
	}
}
