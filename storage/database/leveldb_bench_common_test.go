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

package database

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

const (
	// numDataInsertions is the amount of data to be pre-stored in the DB before it is read
	numDataInsertions = 1000 * 1000 * 2
	// readArea is the range of data to be attempted to read from the data stored in the db.
	readArea = numDataInsertions / 512
	// To read the data cached in the DB, the most recently input data of the DB should be read.
	// This is the offset to specify the starting point of the data to read.
	readCacheOffset = numDataInsertions - readArea
	// primeNumber is used as an interval for reading data.
	// It is a number used to keep the data from being read consecutively.
	// This interval is used because the prime number is small relative to any number except for its multiple,
	// so you can access the whole range of data by cycling the data by a small number of intervals.
	primeNumber = 71887
	// levelDBMemDBSize is the size of internal memory database of LevelDB. Data is first saved to memDB and then moved to persistent storage.
	levelDBMemDBSize = 64
	// numOfGet is the number read per goroutine for the parallel read test.
	numberOfGet = 10000
)

// initTestDB creates the db and inputs the data in db for valSize
func initTestDB(valueSize int) (string, Database, [][]byte, error) {
	dir, err := ioutil.TempDir("", "bench-DB")
	if err != nil {
		return "", nil, nil, errors.New(fmt.Sprintf("can't create temporary directory: %v", err))
	}
	dbc := &DBConfig{Dir: dir, DBType: LevelDB, LevelDBCacheSize: levelDBMemDBSize, OpenFilesLimit: 0}
	db, err := newDatabase(dbc, 0)
	if err != nil {
		return "", nil, nil, errors.New(fmt.Sprintf("can't create database: %v", err))
	}
	keys, values := genKeysAndValues(valueSize, numDataInsertions)

	for i, key := range keys {
		if err := db.Put(key, values[i]); err != nil {
			return "", nil, nil, errors.New(fmt.Sprintf("fail to put data to db: %v", err))
		}
	}

	return dir, db, keys, nil
}

// benchmarkReadDBFromFile is a benchmark function that reads the data stored in the ldb file.
// Reads the initially entered data to read the value stored in the file.
func benchmarkReadDBFromFile(b *testing.B, valueSize int) {
	dir, db, keys, err := initTestDB(valueSize)
	defer os.RemoveAll(dir)
	defer db.Close()
	if err != nil {
		b.Fatalf("database initialization error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Get(keys[(i*primeNumber)%readArea])
	}
}

// benchmarkReadDBFromMemDB is a benchmark function that reads data stored in memDB.
// Read the data entered later to read the value stored in memDB, not in the disk storage.
func benchmarkReadDBFromMemDB(b *testing.B, valueSize int) {
	dir, db, keys, err := initTestDB(valueSize)
	defer os.RemoveAll(dir)
	defer db.Close()
	if err != nil {
		b.Fatalf("database initialization error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Get(keys[(readCacheOffset)+(i*primeNumber)%readArea])
	}
}

// getReadDataOptions is a test case for measuring read performance.
var getReadDataOptions = [...]struct {
	name        string
	valueLength int
	testFunc    func(b *testing.B, valueSize int)
}{
	{"DBFromFile1", 1, benchmarkReadDBFromFile},
	{"DBFromFile128", 128, benchmarkReadDBFromFile},
	{"DBFromFile256", 256, benchmarkReadDBFromFile},
	{"DBFromFile512", 512, benchmarkReadDBFromFile},

	{"DBFromMem1", 1, benchmarkReadDBFromMemDB},
	{"DBFromMem128", 128, benchmarkReadDBFromMemDB},
	{"DBFromMem256", 256, benchmarkReadDBFromMemDB},
	{"DBFromMem512", 512, benchmarkReadDBFromMemDB},

	{"LRUCacheSingle1", 1, benchmarkLruCacheGetSingle},
	{"LRUCacheSingle128", 128, benchmarkLruCacheGetSingle},
	{"LRUCacheSingle256", 256, benchmarkLruCacheGetSingle},
	{"LRUCacheSingle512", 512, benchmarkLruCacheGetSingle},

	{"FIFOCacheSingle1", 1, benchmarkFIFOCacheGetSingle},
	{"FIFOCacheSingle128", 128, benchmarkFIFOCacheGetSingle},
	{"FIFOCacheSingle256", 256, benchmarkFIFOCacheGetSingle},
	{"FIFOCacheSingle512", 512, benchmarkFIFOCacheGetSingle},

	{"LRUCacheParallel1", 1, benchmarkLruCacheCacheGetParallel},
	{"LRUCacheParallel128", 128, benchmarkLruCacheCacheGetParallel},
	{"LRUCacheParallel256", 256, benchmarkLruCacheCacheGetParallel},
	{"LRUCacheParallel512", 512, benchmarkLruCacheCacheGetParallel},

	{"FIFOCacheParallel1", 1, benchmarkFIFOCacheGetParallel},
	{"FIFOCacheParallel128", 128, benchmarkFIFOCacheGetParallel},
	{"FIFOCacheParallel256", 256, benchmarkFIFOCacheGetParallel},
	{"FIFOCacheParallel512", 512, benchmarkFIFOCacheGetParallel},
}

// Benchmark_read_data is a benchmark that measures data read performance in DB and cache.
func Benchmark_read_data(b *testing.B) {
	for _, bm := range getReadDataOptions {
		b.Run(bm.name, func(b *testing.B) {
			bm.testFunc(b, bm.valueLength)
		})
	}
}

// benchmarkFIFOCacheGetParallel measures the performance of the fifoCache when reading data in parallel
func benchmarkFIFOCacheGetParallel(b *testing.B, valueSize int) {
	cache := common.NewCache(common.FIFOCacheConfig{CacheSize: numDataInsertions})
	benchmarkCacheGetParallel(b, cache, valueSize)
}

// benchmarkLruCacheCacheGetParallel measures the performance of the lruCache when reading data in parallel
func benchmarkLruCacheCacheGetParallel(b *testing.B, valueSize int) {
	cache := common.NewCache(common.LRUConfig{CacheSize: numDataInsertions})
	benchmarkCacheGetParallel(b, cache, valueSize)
}

// benchmarkCacheGetParallel is a benchmark for measuring performance
// when cacheSize data is entered into the cache and then read in parallel.
func benchmarkCacheGetParallel(b *testing.B, cache common.Cache, valueSize int) {
	hashKeys := initCacheData(cache, valueSize)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < numberOfGet; i++ {
				key := hashKeys[(i*primeNumber)%numDataInsertions]
				cache.Get(key)
			}
		}
	})
}

// benchmarkFIFOCacheGetSingle is a benchmark to read fifoCache serially.
func benchmarkFIFOCacheGetSingle(b *testing.B, valueSize int) {
	cache := common.NewCache(common.FIFOCacheConfig{CacheSize: numDataInsertions})
	benchmarkCacheGetSingle(b, cache, valueSize)
}

// benchmarkLruCacheGetSingle is a benchmark to read lruCache serially.
func benchmarkLruCacheGetSingle(b *testing.B, valueSize int) {
	cache := common.NewCache(common.LRUConfig{CacheSize: numDataInsertions})
	benchmarkCacheGetSingle(b, cache, valueSize)
}

// benchmarkCacheGetSingle is a benchmark for measuring performance
// when cacheSize data is entered into the cache and then serially read.
func benchmarkCacheGetSingle(b *testing.B, cache common.Cache, valueSize int) {
	hashKeys := initCacheData(cache, valueSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := hashKeys[(i*primeNumber)%numDataInsertions]
		cache.Get(key)
	}
}

// initCacheData initializes the cache by entering test data into the cache.
func initCacheData(cache common.Cache, valueSize int) []common.Hash {
	keys, values := genKeysAndValues(valueSize, numDataInsertions)
	hashKeys := make([]common.Hash, 0, numDataInsertions)

	for i, key := range keys {
		var hashKey common.Hash
		copy(hashKey[:], key[:32])
		hashKeys = append(hashKeys, hashKey)
		cache.Add(hashKey, values[i])
	}

	return hashKeys
}

// TestLRUShardCacheAddressKey is a test to make sure add and get commands work when using Address as key.
// Cache hit for all data.
func TestLRUShardCacheAddressKey(t *testing.T) {
	cache := common.NewCache(common.LRUShardConfig{CacheSize: 40960, NumShards: 4096})

	for i := 0; i < 4096; i++ {
		cache.Add(getAddressKey(i), i)
	}

	for i := 0; i < 4096; i++ {
		data, ok := cache.Get(getAddressKey(i))
		assert.True(t, ok)
		assert.Equal(t, i, data.(int))
	}
}

// TestLRUShardCacheHashKey is a test to check whether the add and get commands are working
// when the Address created by SetBytesFromFront is used as key. Cache hit for all data.
func TestLRUShardCacheAddressKeyFromFront(t *testing.T) {
	cache := common.NewCache(common.LRUShardConfig{CacheSize: 40960, NumShards: 4096})

	for i := 0; i < 4096; i++ {
		cache.Add(getAddressKeyFromFront(i), i)
	}

	for i := 0; i < 4096; i++ {
		data, ok := cache.Get(getAddressKeyFromFront(i))
		assert.True(t, ok)
		assert.Equal(t, i, data.(int))
	}
}

// TestLRUShardCacheHashKey is a test to see if add and get commands work when using Hash as key.
// Cache hit for all data.
func TestLRUShardCacheHashKey(t *testing.T) {
	cache := common.NewCache(common.LRUShardConfig{CacheSize: 40960, NumShards: 4096})

	for i := 0; i < 4096; i++ {
		cache.Add(getHashKey(i), i)
	}

	for i := 0; i < 4096; i++ {
		data, ok := cache.Get(getHashKey(i))
		assert.True(t, ok)
		assert.Equal(t, i, data.(int))
	}
}

// getHashKey returns an int converted to a Hash.
func getHashKey(i int) common.Hash {
	var byteArray interface{}
	if i>>8 == 0 {
		byteArray = []byte{byte(i)}
	} else {
		byteArray = []byte{byte(i >> 8), byte(i)}
	}
	return common.BytesToHash(byteArray.([]byte))
}

// getAddressKey returns an int converted to a Address.
func getAddressKey(i int) common.Address {
	var byteArray interface{}
	if i>>8 == 0 {
		byteArray = []byte{byte(i)}
	} else {
		byteArray = []byte{byte(i >> 8), byte(i)}
	}
	return common.BytesToAddress(byteArray.([]byte))
}

// getAddressKeyFromFront converts an int to an Address and returns it from 0 in the array.
func getAddressKeyFromFront(i int) common.Address {
	var addr common.Address
	var byteArray interface{}
	if i>>8 == 0 {
		byteArray = []byte{byte(i)}
	} else {
		byteArray = []byte{byte(i), byte(i >> 8)}
	}
	addr.SetBytesFromFront(byteArray.([]byte))
	return addr
}
