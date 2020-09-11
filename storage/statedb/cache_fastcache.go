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

import "github.com/VictoriaMetrics/fastcache"

type FastCache struct {
	cache *fastcache.Cache
}

// NewFastCache creates a FastCache with given cache size.
// If you want auto-scaled cache size, set cacheSizeMB to AutoScaling.
// It returns nil if the cache size is zero.
func NewFastCache(cacheSizeMB int) TrieNodeCache {
	if cacheSizeMB == AutoScaling {
		cacheSizeMB = getTrieNodeCacheSizeMB()
	}

	if cacheSizeMB <= 0 {
		return nil
	}

	logger.Info("Initialize local trie node cache (fastCache)", "MaxMB", cacheSizeMB)
	return &FastCache{cache: fastcache.New(cacheSizeMB * 1024 * 1024)} // Covert MB to Byte
}

func (l *FastCache) Get(k []byte) []byte {
	return l.cache.Get(nil, k)
}

func (l *FastCache) Set(k, v []byte) {
	l.cache.Set(k, v)
}

func (l *FastCache) Has(k []byte) ([]byte, bool) {
	return l.cache.HasGet(nil, k)
}

func (l *FastCache) UpdateStats() fastcache.Stats {
	var stats fastcache.Stats
	l.cache.UpdateStats(&stats)

	return stats
}
