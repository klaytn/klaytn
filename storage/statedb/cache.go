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
	SaveToFile(filePath string, concurrency int) error
}

type CacheConfig struct {
	Type CacheType
	// FastCache related configurations
	CacheFilePath string
	MaxBytes      int
	// RedisCache related configurations
	RedisEndpoints   []string
	IsClusteredRedis bool
}

// NewCache creates one type of any supported stateDB caches.
// TODO-Klaytn: refine input parameters after setting node flags
func NewCache(c *CacheConfig) (Cache, error) {
	switch c.Type {
	case LocalCache:
		return NewFastCache(c.CacheFilePath, c.MaxBytes), nil
	case RemoteCache:
		return NewRedisCache(c.RedisEndpoints, c.IsClusteredRedis)
	case HybridCache:
		return NewHybridCache(c.CacheFilePath, c.MaxBytes, c.RedisEndpoints, c.IsClusteredRedis)
	default:
	}
	return nil, errNotSupportedCacheType
}
