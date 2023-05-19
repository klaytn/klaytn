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
