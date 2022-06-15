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
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

var (
	childHash  = common.HexToHash("1341655") // 20190805 in hexadecimal
	parentHash = common.HexToHash("1343A3F") // 20199999 in hexadecimal
)

func TestDatabase_Reference(t *testing.T) {
	memDB := database.NewMemoryDBManager()
	db := NewDatabaseWithNewCache(memDB, &TrieNodeCacheConfig{CacheType: CacheTypeLocal, LocalCacheSizeMiB: 128})

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
	db := NewDatabaseWithNewCache(memDB, &TrieNodeCacheConfig{CacheType: CacheTypeLocal, LocalCacheSizeMiB: 128})
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
	db := NewDatabaseWithNewCache(memDB, &TrieNodeCacheConfig{CacheType: CacheTypeLocal, LocalCacheSizeMiB: 128})

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
}

func TestDatabase_SecureKey(t *testing.T) {
	secKey1 := secureKey(childHash)
	copiedSecKey := make([]byte, len(secKey1))
	copy(copiedSecKey, secKey1)

	secKey2 := secureKey(parentHash)

	assert.Equal(t, secKey1, copiedSecKey) // after the next call of secureKey, secKey1 became different from the copied
	assert.NotEqual(t, secKey1, secKey2)   // secKey1 has changed into secKey2 as they are created from the different buffer
}

func TestCache(t *testing.T) {
	memDB := database.NewMemoryDBManager()
	db := NewDatabaseWithNewCache(memDB, &TrieNodeCacheConfig{CacheType: CacheTypeLocal, LocalCacheSizeMiB: 10})

	for i := 0; i < 100; i++ {
		key, value := common.MakeRandomBytes(256), common.MakeRandomBytes(63*1024) // fastcache can store entrie under 64KB
		db.trieNodeCache.Set(key, value)
		rValue, found := db.trieNodeCache.Has(key)

		assert.Equal(t, true, found)
		assert.Equal(t, value, rValue)
	}
}
