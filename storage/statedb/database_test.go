// Copyright 2019 The klaytn Authors
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
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

var childHash = common.HexToHash("1341655")  // 20190805 in hexadecimal
var parentHash = common.HexToHash("1343A3F") // 20199999 in hexadecimal

func TestDatabase_Reference(t *testing.T) {
	memDB := database.NewMemoryDBManager()
	db := NewDatabaseWithCache(memDB, 128, 0)

	assert.Equal(t, memDB, db.DiskDB())
	assert.Equal(t, 1, len(db.nodes)) // {} : {}

	db.Reference(childHash, parentHash)
	assert.Equal(t, 1, len(db.nodes)) // {} : {}

	child := &cachedNode{}
	parent := &cachedNode{}
	db.nodes[childHash] = child
	db.nodes[parentHash] = parent

	// Call Reference after updating db.nodes
	db.Reference(childHash, parentHash)
	assert.Equal(t, 3, len(db.nodes))
	assert.Equal(t, uint64(1), child.parents)
	assert.Equal(t, uint64(1), parent.children[childHash])

	// Just calling Reference does not have effect
	db.Reference(childHash, parentHash)
	assert.Equal(t, 3, len(db.nodes))
	assert.Equal(t, uint64(1), child.parents)
	assert.Equal(t, uint64(1), parent.children[childHash])
}

func TestDatabase_DeReference(t *testing.T) {
	memDB := database.NewMemoryDBManager()
	db := NewDatabaseWithCache(memDB, 128, 0)
	assert.Equal(t, 1, len(db.nodes)) // {} : {}

	db.Dereference(parentHash)
	assert.Equal(t, 1, len(db.nodes)) // {} : {}
	assert.Equal(t, uint64(0), db.gcnodes)
	assert.Equal(t, common.StorageSize(0), db.gcsize)

	child := &cachedNode{}
	parent := &cachedNode{}
	db.nodes[childHash] = child
	db.nodes[parentHash] = parent

	db.Reference(childHash, parentHash)
	assert.Equal(t, 3, len(db.nodes))
	assert.Equal(t, uint64(1), child.parents)
	assert.Equal(t, uint64(1), parent.children[childHash])
	assert.Equal(t, uint64(0), db.gcnodes)
	assert.Equal(t, common.StorageSize(0), db.gcsize)

	db.Dereference(parentHash)
	assert.Equal(t, 1, len(db.nodes))
	assert.Equal(t, uint64(0), child.parents)
	assert.Equal(t, uint64(0), parent.children[childHash])
	assert.Equal(t, uint64(2), db.gcnodes)
	assert.Equal(t, common.StorageSize(64), db.gcsize)
}

func TestDatabase_Size(t *testing.T) {
	memDB := database.NewMemoryDBManager()
	cacheSizeMB := 128
	db := NewDatabaseWithCache(memDB, cacheSizeMB, 0)

	totalMemorySize, preimagesSize := db.Size()
	assert.Equal(t, common.StorageSize(0), totalMemorySize)
	assert.Equal(t, common.StorageSize(0), preimagesSize)

	child := &cachedNode{}
	parent := &cachedNode{}
	db.nodes[childHash] = child
	db.nodes[parentHash] = parent

	db.Reference(childHash, parentHash)

	totalMemorySize, preimagesSize = db.Size()
	assert.Equal(t, common.StorageSize(128), totalMemorySize)
	assert.Equal(t, common.StorageSize(0), preimagesSize)

	db.preimagesSize += 100
	totalMemorySize, preimagesSize = db.Size()
	assert.Equal(t, common.StorageSize(128), totalMemorySize)
	assert.Equal(t, common.StorageSize(100), preimagesSize)

	cacheSize := db.CacheSize()
	assert.Equal(t, cacheSizeMB*1024*1024, cacheSize)
}

func TestDatabase_SecureKey(t *testing.T) {
	memDB := database.NewMemoryDBManager()
	db := NewDatabaseWithCache(memDB, 128, 0)

	secKey1 := db.secureKey(childHash[:])
	copiedSecKey := make([]byte, 0, len(secKey1))
	copy(copiedSecKey, secKey1)

	secKey2 := db.secureKey(parentHash[:])

	assert.NotEqual(t, secKey1, copiedSecKey) // after the next call of secureKey, secKey1 became different from the copied
	assert.Equal(t, secKey1, secKey2)         // secKey1 has changed into secKey2 as they are created from the same buffer
}
