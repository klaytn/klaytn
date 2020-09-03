package statedb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testHosts = []string{"localhost:6379"}

// TODO-Klaytn: Enable tests when redis is prepared on CI

// TestNewRedisCache tests basic operations of redis cache
func _TestNewRedisCache(t *testing.T) {
	redis, err := NewRedisCache(testHosts, false)
	assert.Nil(t, err)

	key, value := randBytes(32), randBytes(500)
	redis.Set(key, value)

	getValue := redis.Get(key)
	assert.Equal(t, bytes.Compare(value, getValue), 0)

	hasValue, ok := redis.Has(key)
	assert.Equal(t, ok, true)
	assert.Equal(t, bytes.Compare(value, hasValue), 0)
}

// TestNewRedisCache_Set_LargeData check whether redis cache can store an large data (5MB).
func _TestNewRedisCache_Set_LargeData(t *testing.T) {
	redis, err := NewRedisCache(testHosts, false)
	assert.Nil(t, err)

	key, value := randBytes(32), randBytes(5*1024*1024) // 5MB value
	redis.Set(key, value)

	retValue := redis.Get(key)
	assert.Equal(t, bytes.Compare(value, retValue), 0)
}
