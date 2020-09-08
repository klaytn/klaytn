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

type TrieNodeCacheType string

const (
	// Available stateDB cache types
	LocalCache  TrieNodeCacheType = "LocalCache"
	RemoteCache                   = "RemoteCache"
	HybridCache                   = "HybridCache"
)

var (
	errNotSupportedCacheType = errors.New("not supported stateDB TrieNodeCache type")
)

// TrieNodeCache interface the cache of stateDB
type TrieNodeCache interface {
	Set(k, v []byte)
	Get(k []byte) []byte
	Has(k []byte) ([]byte, bool)
}

// NewTrieNodeCache creates one type of any supported trie node caches.
// TODO-Klaytn: refine input parameters after setting node flags
func NewTrieNodeCache(cacheType TrieNodeCacheType, maxBytes int, redisEndpoint []string, redisCluster bool) (TrieNodeCache, error) {
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
