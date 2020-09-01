package statedb

import "github.com/VictoriaMetrics/fastcache"

type LocalCache struct {
	cache *fastcache.Cache
}

func NewLocalCache(maxBytes int) Cache {
	return &LocalCache{
		cache: fastcache.New(maxBytes),
	}
}

func (l *LocalCache) Get(k []byte) []byte {
	return l.cache.Get(nil, k)
}

func (l *LocalCache) Set(k, v []byte) {
	l.cache.Set(k, v)
}

func (l *LocalCache) Has(k []byte) ([]byte, bool) {
	return l.cache.HasGet(nil, k)
}

func (l *LocalCache) UpdateStats() fastcache.Stats {
	var stats fastcache.Stats
	l.cache.UpdateStats(&stats)

	return stats
}
