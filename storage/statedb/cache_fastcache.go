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

func NewFastCache(filePath string, maxBytes int) Cache {
	return &FastCache{
		cache: fastcache.LoadFromFileOrNew(filePath, maxBytes),
	}
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

func (l *FastCache) SaveToFile(filePath string, concurrency int) error {
	return l.cache.SaveToFileConcurrent(filePath, concurrency)
}
