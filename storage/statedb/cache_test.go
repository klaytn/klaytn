// Copyright 2020 The klaytn Authors
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
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/klaytn/klaytn/common"

	"github.com/docker/docker/pkg/testutil/assert"
)

// TestNewTrieNodeCache tests creating all kinds of supported trie node caches.
func TestNewTrieNodeCache(t *testing.T) {
	testCases := []struct {
		cacheConfig  *TrieNodeCacheConfig
		expectedType reflect.Type
		err          error
	}{
		{getTestFastCacheConfig(), reflect.TypeOf(&FastCache{}), nil},
		{getTestRedisConfig(), reflect.TypeOf(&RedisCache{}), nil},
		{getTestHybridConfig(), reflect.TypeOf(&HybridCache{}), nil},
		{nil, nil, errNilTrieNodeCacheConfig},
	}

	for _, tc := range testCases {
		cache, err := NewTrieNodeCache(tc.cacheConfig)
		assert.Equal(t, err, tc.err)
		assert.Equal(t, reflect.TypeOf(cache), tc.expectedType)
	}
}

func TestFastCache_SaveAndLoad(t *testing.T) {
	// Create test directory
	dirName, err := ioutil.TempDir(os.TempDir(), "fastcache_saveandload")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirName)

	// Generate test data
	var keys [][]byte
	var vals [][]byte
	for i := 0; i < 10; i++ {
		keys = append(keys, common.MakeRandomBytes(128))
		vals = append(vals, common.MakeRandomBytes(128))
	}

	config := getTestHybridConfig()
	config.FastCacheFileDir = dirName

	// Create a fastcache from the file and save the data to the cache
	fastCache := newFastCache(config)
	for idx, key := range keys {
		assert.DeepEqual(t, fastCache.Get(key), []byte(nil))
		fastCache.Set(key, vals[idx])
		assert.DeepEqual(t, fastCache.Get(key), vals[idx])
	}
	// Save the cache to the file
	assert.NilError(t, fastCache.SaveToFile(dirName, runtime.NumCPU()))

	// Create a fastcache from the file and check if the data exists
	fastCacheFromFile := newFastCache(config)
	for idx, key := range keys {
		assert.DeepEqual(t, fastCacheFromFile.Get(key), vals[idx])
	}
}
