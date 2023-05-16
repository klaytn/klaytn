package statedb

import (
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

func checkHasherHash(t *testing.T, idx int, tc *testNodeEncodingTC) {
	hash := common.BytesToHash(tc.hash)

	memDB := database.NewMemoryDBManager()
	db := NewDatabase(memDB)
	h := newHasher(nil)
	defer returnHasherToPool(h)

	hashed, cached := h.hash(tc.expanded, db, false)
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

func TestHasherHashTC(t *testing.T) {
	for idx, tc := range collapsedNodeTCs() {
		checkHasherHash(t, idx, tc)
	}
	for idx, tc := range resolvedNodeTCs() {
		checkHasherHash(t, idx, tc)
	}
}
