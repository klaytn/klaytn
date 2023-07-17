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
	"github.com/stretchr/testify/require"
)

func checkHasherHash(t *testing.T, name string, tc testNodeEncodingTC, opts *hasherOpts, onRoot bool) {
	common.ResetExtHashCounterForTest(0xccccddddeeee00)
	memDB := database.NewMemoryDBManager()
	db := NewDatabase(memDB)

	h := newHasher(opts)
	defer returnHasherToPool(h)

	hashed, cached := h.hashNode(tc.expanded, db, false, onRoot)
	t.Logf("tc[%s] %s", name, hashed)
	assert.Equal(t, hashNode(tc.hash), hashed, name)

	cachedHash, _ := cached.cache()
	assert.Equal(t, hashNode(tc.hash), cachedHash, name)

	hash := common.BytesToExtHash(tc.hash)
	inserted := db.nodes[hash]
	require.NotNil(t, inserted)
	assert.Equal(t, tc.inserted, inserted.node, name)

	db.Cap(0)
	encoded, _ := memDB.ReadTrieNode(hash)
	assert.Equal(t, tc.encoded, encoded, name)
}

func TestHasherHashTC(t *testing.T) {
	optsLegacy := &hasherOpts{}
	optsState := &hasherOpts{pruning: true}
	optsStorage := &hasherOpts{pruning: true, storageRoot: true}

	for name, tc := range collapsedNodeTCs_legacy() {
		checkHasherHash(t, name, tc, optsLegacy, true)
		checkHasherHash(t, name, tc, optsLegacy, false)
	}
	for name, tc := range resolvedNodeTCs_legacy() {
		checkHasherHash(t, name, tc, optsLegacy, true)
		checkHasherHash(t, name, tc, optsLegacy, false)
	}

	for name, tc := range collapsedNodeTCs_extroot() {
		checkHasherHash(t, name, tc, optsState, true)
	}
	for name, tc := range resolvedNodeTCs_extroot() {
		checkHasherHash(t, name, tc, optsState, true)
	}

	for name, tc := range collapsedNodeTCs_exthash() {
		checkHasherHash(t, name, tc, optsState, false)
		checkHasherHash(t, name, tc, optsStorage, true)
		checkHasherHash(t, name, tc, optsStorage, false)
	}
	for name, tc := range resolvedNodeTCs_exthash() {
		checkHasherHash(t, name, tc, optsState, false)
		checkHasherHash(t, name, tc, optsStorage, true)
		checkHasherHash(t, name, tc, optsStorage, false)
	}
}
