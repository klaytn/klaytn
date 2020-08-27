package statedb

import "github.com/VictoriaMetrics/fastcache"

// Cache interface the cache of stateDB
type Cache interface {
	Set(k, v []byte)
	Get(dst, k []byte) []byte
	HasGet(dst, k []byte) ([]byte, bool)
	UpdateStats(s *fastcache.Stats)
}
