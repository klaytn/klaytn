package statedb

// Cache interface the cache of stateDB
type Cache interface {
	Set(k, v []byte)
	Get(k []byte) []byte
	Has(k []byte) ([]byte, bool)
}
