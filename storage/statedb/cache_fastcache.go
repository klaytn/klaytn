package statedb

import "github.com/VictoriaMetrics/fastcache"

type FastCache struct {
	cache *fastcache.Cache
}

func NewFastCache(maxBytes int) Cache {
	return &FastCache{
		cache: fastcache.New(maxBytes),
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
