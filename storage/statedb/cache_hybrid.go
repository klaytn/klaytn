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

func NewHybridCache(filePath string, maxBytes int, endpoint []string, isCluster bool) (Cache, error) {
	redis, err := NewRedisCache(endpoint, isCluster)
	if err != nil {
		return nil, err
	}
	return &hybridCache{
		local:  NewFastCache(filePath, maxBytes),
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

func (cache *hybridCache) SaveToFile(filePath string, concurrency int) error {
	if err := cache.local.SaveToFile(filePath, concurrency); err != nil {
		logger.Error("failed to save local cache to file",
			"filePath", filePath, "concurrency", concurrency, "err", err)
		return err
	}
	if err := cache.remote.SaveToFile(filePath, concurrency); err != nil {
		logger.Error("failed to save remote cache to file",
			"filePath", filePath, "concurrency", concurrency, "err", err)
		return err
	}
	return nil
}
