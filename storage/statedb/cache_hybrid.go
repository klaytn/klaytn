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

import "github.com/go-redis/redis/v7"

func newHybridCache(config *TrieNodeCacheConfig) (TrieNodeCache, error) {
	redis, err := newRedisCache(config)
	if err != nil {
		return nil, err
	}

	return &HybridCache{
		local:  newFastCache(config),
		remote: redis,
	}, nil
}

// HybridCache integrates two kinds of caches: local, remote.
// Local cache uses memory of the local machine and remote cache uses memory of the remote machine.
// When it sets data to both caches, only remote cache is set asynchronously
type HybridCache struct {
	local  TrieNodeCache
	remote *RedisCache
}

func (cache *HybridCache) Local() TrieNodeCache {
	return cache.local
}

func (cache *HybridCache) Remote() *RedisCache {
	return cache.remote
}

// Set writes data to local cache synchronously and to remote cache asynchronously.
func (cache *HybridCache) Set(k, v []byte) {
	cache.local.Set(k, v)
	cache.remote.SetAsync(k, v)
}

func (cache *HybridCache) Get(k []byte) []byte {
	ret := cache.local.Get(k)
	if ret != nil {
		return ret
	}
	ret = cache.remote.Get(k)
	if ret != nil {
		cache.local.Set(k, ret)
	}
	return ret
}

func (cache *HybridCache) Has(k []byte) ([]byte, bool) {
	ret, has := cache.local.Has(k)
	if has {
		return ret, has
	}
	return cache.remote.Has(k)
}

func (cache *HybridCache) UpdateStats() interface{} {
	type stats struct {
		local  interface{}
		remote interface{}
	}
	return stats{cache.local.UpdateStats(), cache.remote.UpdateStats()}
}

func (cache *HybridCache) SaveToFile(filePath string, concurrency int) error {
	if err := cache.local.SaveToFile(filePath, concurrency); err != nil {
		logger.Error("failed to save local cache to file",
			"filePath", filePath, "concurrency", concurrency, "err", err)
		return err
	}
	return nil
}

func (cache *HybridCache) PublishBlock(msg string) error {
	return cache.remote.PublishBlock(msg)
}

func (cache *HybridCache) SubscribeBlockCh() <-chan *redis.Message {
	return cache.remote.SubscribeBlockCh()
}

func (cache *HybridCache) UnsubscribeBlock() error {
	return cache.remote.UnsubscribeBlock()
}

func (cache *HybridCache) Close() error {
	err := cache.local.Close()
	if err != nil {
		return err
	}
	return cache.remote.Close()
}
