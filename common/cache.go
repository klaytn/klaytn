// Copyright 2018 The klaytn Authors
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

package common

import (
	"errors"
	"math"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/log"
	"github.com/pbnjay/memory"
)

type CacheType int

const (
	LRUCacheType CacheType = iota
	LRUShardCacheType
	FIFOCacheType
	ARCChacheType
)

const (
	cacheLevelSaving  = "saving"
	cacheLevelNormal  = "normal"
	cacheLevelExtreme = "extreme"
)
const (
	minimumMemorySize      = 16
	defaultCacheUsageLevel = cacheLevelSaving
)

// it's set by flag.
var DefaultCacheType CacheType = FIFOCacheType
var logger = log.NewModuleLogger(log.Common)
var CacheScale int = 100                             // Cache size = preset size * CacheScale / 100. Only used when IsScaled == true
var ScaleByCacheUsageLevel int = 100                 // Scale according to cache usage level (%). Only used when IsScaled == true
var TotalPhysicalMemGB int = getPhysicalMemorySize() // Convert Byte to GByte

// getPhysicalMemorySize returns the system's physical memory value.
// It internally returns a minimumMemorySize if it is an os that does not support using the system call to obtain it,
// or if the system call fails.
func getPhysicalMemorySize() int {
	TotalMemGB := int(memory.TotalMemory() / 1024 / 1024 / 1024)
	if TotalMemGB >= minimumMemorySize {
		return TotalMemGB
	} else {
		if TotalMemGB != 0 {
			logger.Warn("The system's physical memory is less than minimum physical memory size", "physicalSystemMemory(GB)", TotalMemGB, "minimumMemorySize(GB)", minimumMemorySize)
		} else {
			logger.Error("Failed to get the physical memory of the system. Minimum physical memory size is used", "minimumMemorySize(GB)", minimumMemorySize)
		}
		return minimumMemorySize
	}
}

type CacheKey interface {
	getShardIndex(shardMask int) int
}

type Cache interface {
	Add(key CacheKey, value interface{}) (evicted bool)
	Get(key CacheKey) (value interface{}, ok bool)
	Contains(key CacheKey) bool
	Purge()
}

type lruCache struct {
	lru *lru.Cache
}

type CacheKeyUint64 uint64

func (key CacheKeyUint64) getShardIndex(shardMask int) int {
	return int(uint64(key) & uint64(shardMask))
}

func (cache *lruCache) Add(key CacheKey, value interface{}) (evicted bool) {
	return cache.lru.Add(key, value)
}

func (cache *lruCache) Get(key CacheKey) (value interface{}, ok bool) {
	value, ok = cache.lru.Get(key)
	return
}

func (cache *lruCache) Contains(key CacheKey) bool {
	return cache.lru.Contains(key)
}

func (cache *lruCache) Purge() {
	cache.lru.Purge()
}

func (cache *lruCache) Keys() []interface{} {
	return cache.lru.Keys()
}

func (cache *lruCache) Peek(key CacheKey) (value interface{}, ok bool) {
	return cache.lru.Peek(key)
}

func (cache *lruCache) Remove(key CacheKey) {
	cache.lru.Remove(key)
}

func (cache *lruCache) Len() int {
	return cache.lru.Len()
}

type arcCache struct {
	arc *lru.ARCCache
}

func (cache *arcCache) Add(key CacheKey, value interface{}) (evicted bool) {
	cache.arc.Add(key, value)
	//TODO-Klaytn-RemoveLater need to be removed or should be added according to usage of evicted flag
	return true
}

func (cache *arcCache) Get(key CacheKey) (value interface{}, ok bool) {
	return cache.arc.Get(key)
}

func (cache *arcCache) Contains(key CacheKey) bool {
	return cache.arc.Contains(key)
}

func (cache *arcCache) Purge() {
	cache.arc.Purge()
}

func (cache *arcCache) Keys() []interface{} {
	return cache.arc.Keys()
}

func (cache *arcCache) Peek(key CacheKey) (value interface{}, ok bool) {
	return cache.arc.Peek(key)
}

func (cache *arcCache) Remove(key CacheKey) {
	cache.arc.Remove(key)
}

func (cache *arcCache) Len() int {
	return cache.arc.Len()
}

type lruShardCache struct {
	shards         []*lru.Cache
	shardIndexMask int
}

func (cache *lruShardCache) Add(key CacheKey, val interface{}) (evicted bool) {
	shardIndex := key.getShardIndex(cache.shardIndexMask)
	return cache.shards[shardIndex].Add(key, val)
}

func (cache *lruShardCache) Get(key CacheKey) (value interface{}, ok bool) {
	shardIndex := key.getShardIndex(cache.shardIndexMask)
	return cache.shards[shardIndex].Get(key)
}

func (cache *lruShardCache) Contains(key CacheKey) bool {
	shardIndex := key.getShardIndex(cache.shardIndexMask)
	return cache.shards[shardIndex].Contains(key)
}

func (cache *lruShardCache) Purge() {
	for _, shard := range cache.shards {
		s := shard
		go s.Purge()
	}
}

func NewCache(config CacheConfiger) Cache {
	if config == nil {
		logger.Crit("config shouldn't be nil!")
	}

	cache, err := config.newCache()
	if err != nil {
		logger.Crit("Failed to allocate cache!", "err", err)
	}
	return cache
}

type CacheConfiger interface {
	newCache() (Cache, error)
}

type LRUConfig struct {
	CacheSize int
	IsScaled  bool
}

func (c LRUConfig) newCache() (Cache, error) {
	cacheSize := c.CacheSize
	if c.IsScaled {
		cacheSize *= calculateScale()
	}
	lru, err := lru.New(cacheSize)
	return &lruCache{lru}, err
}

type LRUShardConfig struct {
	CacheSize int
	// Hash, and Address type can not generate as many shard indexes as the maximum (2 ^ 16 = 65536),
	// so it is meaningless to set the NumShards larger than this.
	NumShards int
	IsScaled  bool
}

const (
	minShardSize = 10
	minNumShards = 2
)

//If key is not common.Hash nor common.Address then you should set numShard 1 or use LRU Cache
//The number of shards is readjusted to meet the minimum shard size.
func (c LRUShardConfig) newCache() (Cache, error) {
	cacheSize := c.CacheSize
	if c.IsScaled {
		cacheSize *= calculateScale()
	}

	if cacheSize < 1 {
		logger.Error("Negative Cache Size Error", "Cache Size", cacheSize, "Cache Scale", CacheScale)
		return nil, errors.New("Must provide a positive size ")
	}

	numShards := c.makeNumShardsPowOf2()

	if c.NumShards != numShards {
		logger.Warn("numShards is ", "Expected", c.NumShards, "Actual", numShards)
	}
	if cacheSize%numShards != 0 {
		logger.Warn("Cache size is ", "Expected", cacheSize, "Actual", cacheSize-(cacheSize%numShards))
	}

	lruShard := &lruShardCache{shards: make([]*lru.Cache, numShards), shardIndexMask: numShards - 1}
	shardsSize := cacheSize / numShards
	var err error
	for i := 0; i < numShards; i++ {
		lruShard.shards[i], err = lru.NewWithEvict(shardsSize, nil)

		if err != nil {
			return nil, err
		}
	}
	return lruShard, nil
}

func (c LRUShardConfig) makeNumShardsPowOf2() int {
	maxNumShards := float64(c.CacheSize * calculateScale() / minShardSize)
	numShards := int(math.Min(float64(c.NumShards), maxNumShards))

	preNumShards := minNumShards
	for numShards > minNumShards {
		preNumShards = numShards
		numShards = numShards & (numShards - 1)
	}

	return preNumShards
}

// FIFOCacheConfig is a implementation of CacheConfiger interface for fifoCache.
type FIFOCacheConfig struct {
	CacheSize int
	IsScaled  bool
}

// newCache creates a Cache interface whose implementation is fifoCache.
func (c FIFOCacheConfig) newCache() (Cache, error) {
	cacheSize := c.CacheSize
	if c.IsScaled {
		cacheSize *= calculateScale()
	}

	lru, err := lru.New(cacheSize)
	return &fifoCache{&lruCache{lru}}, err
}

// fifoCache internally has a lruCache.
// All methods are the same as lruCache, but we override Get function, not to update the lifetime of data.
type fifoCache struct {
	*lruCache
}

// Get returns the value corresponding to the cache key.
func (cache *fifoCache) Get(key CacheKey) (value interface{}, ok bool) {
	return cache.Peek(key)
}

type ARCConfig struct {
	CacheSize int
	IsScaled  bool
}

func (c ARCConfig) newCache() (Cache, error) {
	cacheSize := c.CacheSize
	if c.IsScaled {
		cacheSize *= calculateScale()
	}
	arc, err := lru.NewARC(cacheSize)
	return &arcCache{arc}, err
}

// calculateScale returns the scale of the cache.
// The scale of the cache is obtained by multiplying (MemorySize / minimumMemorySize), (scaleByCacheUsageLevel / 100), and (CacheScale / 100).
func calculateScale() int {
	return CacheScale * ScaleByCacheUsageLevel * TotalPhysicalMemGB / minimumMemorySize / 100 / 100
}

// GetScaleByCacheUsageLevel returns the scale according to cacheUsageLevel
func GetScaleByCacheUsageLevel(cacheUsageLevelFlag string) (int, error) {
	switch cacheUsageLevelFlag {
	case cacheLevelSaving:
		return 100, nil
	case cacheLevelNormal:
		return 200, nil
	case cacheLevelExtreme:
		return 300, nil
	default:
		return 100, errors.New("input string does not meet the given format. expected: ('saving', 'normal, 'extreme')")
	}
}

type GovernanceCacheKey string

func (g GovernanceCacheKey) getShardIndex(shardMask int) int {
	return 0
}
