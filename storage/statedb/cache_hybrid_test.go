package statedb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO-Klaytn: Enable tests when redis is prepared on CI

var (
	testMaxBytes       = 1024 * 1024
	testRedisEndpoints = []string{"localhost:6379"}
	testRedisCluster   = false
)

// TestHybridCache_Set tests whether a hybrid cache can set an item into both of local and remote caches.
func _TestHybridCache_Set(t *testing.T) {
	cache, err := NewHybridCache("", testMaxBytes, testRedisEndpoints, testRedisCluster)
	if err != nil {
		t.Fatal(err)
	}

	// Set an item
	key, value := randBytes(32), randBytes(500)
	cache.Set(key, value)

	// Type assertion to check both of local cache and remote cache
	hybrid, ok := cache.(*hybridCache)
	assert.Equal(t, ok, true)

	// Check whether the item is set in the local cache
	localVal := hybrid.local.Get(key)
	assert.Equal(t, bytes.Compare(localVal, value), 0)

	// Check whether the item is set in the remote cache
	remoteVal := hybrid.remote.Get(key)
	assert.Equal(t, bytes.Compare(remoteVal, value), 0)
}

// TestHybridCache_Get tests whether a hybrid cache can get an item from both of local and remote caches.
func _TestHybridCache_Get(t *testing.T) {
	// Prepare caches to be integrated with a hybrid cache
	localCache := NewFastCache("", testMaxBytes)
	remoteCache, err := NewRedisCache(testRedisEndpoints, testRedisCluster)
	if err != nil {
		t.Fatal(err)
	}

	var hybrid Cache = &hybridCache{
		local:  localCache,
		remote: remoteCache,
	}

	// Test local cache of the hybrid cache
	{
		// Store an item into local cache
		key, value := randBytes(32), randBytes(500)
		localCache.Set(key, value)

		// Get the item from the hybrid cache and check the validity
		returnedVal := hybrid.Get(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
	}

	// Test remote cache of the hybrid cache
	{
		// Store an item into remote cache
		key, value := randBytes(32), randBytes(500)
		remoteCache.Set(key, value)

		// Get the item from the hybrid cache and check the validity
		returnedVal := hybrid.Get(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
	}

	// Test the priority of local and remote caches
	{
		// Store an item into the remote cache
		key, value := randBytes(32), randBytes(500)
		localCache.Set(key, value)
		remoteCache.Set(key, []byte{0x11})

		// Get the item from the hybrid cache and check the validity
		returnedVal := hybrid.Get(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
	}
}

// TestHybridCache_Has tests whether a hybrid cache can check an item from both of local and remote caches.
func _TestHybridCache_Has(t *testing.T) {
	// Prepare caches to be integrated with a hybrid cache
	localCache := NewFastCache("", testMaxBytes)
	remoteCache, err := NewRedisCache(testRedisEndpoints, testRedisCluster)
	if err != nil {
		t.Fatal(err)
	}

	var hybrid Cache = &hybridCache{
		local:  localCache,
		remote: remoteCache,
	}

	// Test local cache of the hybrid cache
	{
		// Store an item into local cache
		key, value := randBytes(32), randBytes(500)
		localCache.Set(key, value)

		// Get the item from the hybrid cache and check the validity
		returnedVal, returnedExist := hybrid.Has(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
		assert.Equal(t, returnedExist, true)
	}

	// Test remote cache of the hybrid cache
	{
		// Store an item into remote cache
		key, value := randBytes(32), randBytes(500)
		remoteCache.Set(key, value)

		// Get the item from the hybrid cache and check the validity
		returnedVal, returnedExist := hybrid.Has(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
		assert.Equal(t, returnedExist, true)
	}

	// Test the priority of local and remote caches
	{
		// Store an item into the remote cache
		key, value := randBytes(32), randBytes(500)
		localCache.Set(key, value)
		remoteCache.Set(key, []byte{0x11})

		// Get the item from the hybrid cache and check the validity
		returnedVal, returnedExist := hybrid.Has(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
		assert.Equal(t, returnedExist, true)
	}
}
