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

import "errors"

type CacheType string

const (
	// Available stateDB cache types
	LocalCache  CacheType = "LocalCache"
	RemoteCache           = "RemoteCache"
	HybridCache           = "HybridCache"
)

var (
	errNotSupportedCacheType = errors.New("not supported stateDB Cache type")
)

// Cache interface the cache of stateDB
type Cache interface {
	Set(k, v []byte)
	Get(k []byte) []byte
	Has(k []byte) ([]byte, bool)
}

// NewCache creates one type of any supported stateDB caches.
// TODO-Klaytn: refine input parameters after setting node flags
func NewCache(cacheType CacheType, maxBytes int, redisEndpoint []string, redisCluster bool) (Cache, error) {
	switch cacheType {
	case LocalCache:
		return NewFastCache(maxBytes), nil
	case RemoteCache:
		return NewRedisCache(redisEndpoint, redisCluster)
	case HybridCache:
		return NewHybridCache(maxBytes, redisEndpoint, redisCluster)
	default:
	}
	return nil, errNotSupportedCacheType
}

func NewHybridCache(maxBytes int, endpoint []string, isCluster bool) (Cache, error) {
	redis, err := NewRedisCache(endpoint, isCluster)
	if err != nil {
		return nil, err
	}
	return &hybridCache{
		local:  NewFastCache(maxBytes),
		remote: redis,
	}, nil
}

// hybridCache integrates two kinds of caches: local, remote.
// local cache uses memory of the local machine and remote cache uses memory of the remote machine.
type hybridCache struct {
	local  Cache
	remote Cache
}

func (cache *hybridCache) Set(k, v []byte) {
	cache.local.Set(k, v)
	cache.remote.Set(k, v)
}

func (cache *hybridCache) Get(k []byte) []byte {
	ret := cache.local.Get(k)
	if ret != nil {
		return ret
	}
	return cache.remote.Get(k)
}

func (cache *hybridCache) Has(k []byte) ([]byte, bool) {
	ret, has := cache.local.Has(k)
	if has {
		return ret, has
	}
	return cache.remote.Has(k)
}
