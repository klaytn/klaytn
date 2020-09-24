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
	"github.com/VictoriaMetrics/fastcache"
	"github.com/rcrowley/go-metrics"
)

var (
	// metrics
	memcacheFastMisses                 = metrics.NewRegisteredGauge("trie/memcache/fast/misses", nil)
	memcacheFastCollisions             = metrics.NewRegisteredGauge("trie/memcache/fast/collisions", nil)
	memcacheFastCorruptions            = metrics.NewRegisteredGauge("trie/memcache/fast/corruptions", nil)
	memcacheFastEntriesCount           = metrics.NewRegisteredGauge("trie/memcache/fast/entries", nil)
	memcacheFastBytesSize              = metrics.NewRegisteredGauge("trie/memcache/fast/size", nil)
	memcacheFastGetBigCalls            = metrics.NewRegisteredGauge("trie/memcache/fast/get", nil)
	memcacheFastSetBigCalls            = metrics.NewRegisteredGauge("trie/memcache/fast/set", nil)
	memcacheFastTooBigKeyErrors        = metrics.NewRegisteredGauge("trie/memcache/fast/error/too/bigkey", nil)
	memcacheFastInvalidMetavalueErrors = metrics.NewRegisteredGauge("trie/memcache/fast/error/invalid/matal", nil)
	memcacheFastInvalidValueLenErrors  = metrics.NewRegisteredGauge("trie/memcache/fast/error/invalid/valuelen", nil)
	memcacheFastInvalidValueHashErrors = metrics.NewRegisteredGauge("trie/memcache/fast/error/invalid/hash", nil)
)

type FastCache struct {
	cache *fastcache.Cache
}

// NewFastCache creates a FastCache with given cache size.
// If you want auto-scaled cache size, set config.LocalCacheSizeMB to AutoScaling.
// It returns nil if the cache size is zero.
func NewFastCache(config TrieNodeCacheConfig) TrieNodeCache {
	if config.LocalCacheSizeMB == AutoScaling {
		config.LocalCacheSizeMB = getTrieNodeCacheSizeMB()
	}

	if config.LocalCacheSizeMB <= 0 {
		return nil
	}

	logger.Info("Initialize local trie node cache (fastCache)",
		"MaxMB", config.LocalCacheSizeMB, "FilePath", config.FastCacheFileDir)
	return &FastCache{cache: fastcache.LoadFromFileOrNew(config.FastCacheFileDir, config.LocalCacheSizeMB*1024*1024)} // Convert MB to Byte
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

func (l *FastCache) UpdateStats() interface{} {
	var stats fastcache.Stats
	l.cache.UpdateStats(&stats)

	memcacheFastMisses.Update(int64(stats.Misses))
	memcacheFastCollisions.Update(int64(stats.Collisions))
	memcacheFastCorruptions.Update(int64(stats.Corruptions))
	memcacheFastEntriesCount.Update(int64(stats.EntriesCount))
	memcacheFastBytesSize.Update(int64(stats.BytesSize))
	memcacheFastGetBigCalls.Update(int64(stats.GetBigCalls))
	memcacheFastSetBigCalls.Update(int64(stats.SetBigCalls))
	memcacheFastTooBigKeyErrors.Update(int64(stats.TooBigKeyErrors))
	memcacheFastInvalidMetavalueErrors.Update(int64(stats.InvalidMetavalueErrors))
	memcacheFastInvalidValueLenErrors.Update(int64(stats.InvalidValueLenErrors))
	memcacheFastInvalidValueHashErrors.Update(int64(stats.InvalidValueHashErrors))

	return stats
}

func (l *FastCache) SaveToFile(filePath string, concurrency int) error {
	return l.cache.SaveToFileConcurrent(filePath, concurrency)
}
