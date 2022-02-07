package statedb

import (
	"bytes"
	"testing"
	"time"

	"github.com/klaytn/klaytn/storage"
	"github.com/stretchr/testify/assert"
)

func getTestHybridConfig() *TrieNodeCacheConfig {
	return &TrieNodeCacheConfig{
		CacheType:          CacheTypeHybrid,
		LocalCacheSizeMiB:  100,
		FastCacheFileDir:   "",
		RedisEndpoints:     []string{"localhost:6379"},
		RedisClusterEnable: false,
	}
}

// TestHybridCache_Set tests whether a hybrid cache can set an item into both of local and remote caches.
func TestHybridCache_Set(t *testing.T) {
	storage.SkipLocalTest(t)

	cache, err := newHybridCache(getTestHybridConfig())
	if err != nil {
		t.Fatal(err)
	}

	// Set an item
	key, value := randBytes(32), randBytes(500)
	cache.Set(key, value)
	time.Sleep(sleepDurationForAsyncBehavior)

	// Type assertion to check both of local cache and remote cache
	hybrid, ok := cache.(*HybridCache)
	assert.Equal(t, ok, true)

	// Check whether the item is set in the local cache
	localVal := hybrid.local.Get(key)
	assert.Equal(t, bytes.Compare(localVal, value), 0)

	// Check whether the item is set in the remote cache
	remoteVal := hybrid.remote.Get(key)
	assert.Equal(t, bytes.Compare(remoteVal, value), 0)
}

// TestHybridCache_Get tests whether a hybrid cache can get an item from both of local and remote caches.
func TestHybridCache_Get(t *testing.T) {
	storage.SkipLocalTest(t)

	// Prepare caches to be integrated with a hybrid cache
	localCache := newFastCache(getTestHybridConfig())
	remoteCache, err := newRedisCache(getTestHybridConfig())
	if err != nil {
		t.Fatal(err)
	}

	var hybrid TrieNodeCache = &HybridCache{
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
		remoteCache.SetAsync(key, value)
		time.Sleep(sleepDurationForAsyncBehavior)

		// Make sure the item is not stored in the local cache.
		assert.Equal(t, len(localCache.Get(key)), 0)

		// Get the item from the hybrid cache and check the validity
		returnedVal := hybrid.Get(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)

		// Make sure that the item retrieved from the remote cache is also stored in the local cache
		assert.Equal(t, bytes.Compare(localCache.Get(key), value), 0)
	}

	// Test the priority of local and remote caches
	{
		// Store an item into the remote cache
		key, value := randBytes(32), randBytes(500)
		localCache.Set(key, value)
		remoteCache.SetAsync(key, []byte{0x11})
		time.Sleep(sleepDurationForAsyncBehavior)

		// Get the item from the hybrid cache and check the validity
		returnedVal := hybrid.Get(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
	}
}

// TestHybridCache_Has tests whether a hybrid cache can check an item from both of local and remote caches.
func TestHybridCache_Has(t *testing.T) {
	storage.SkipLocalTest(t)

	// Prepare caches to be integrated with a hybrid cache
	localCache := newFastCache(getTestHybridConfig())
	remoteCache, err := newRedisCache(getTestHybridConfig())
	if err != nil {
		t.Fatal(err)
	}

	var hybrid TrieNodeCache = &HybridCache{
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
		remoteCache.SetAsync(key, value)
		time.Sleep(sleepDurationForAsyncBehavior)

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
		remoteCache.SetAsync(key, []byte{0x11})
		time.Sleep(sleepDurationForAsyncBehavior)

		// Get the item from the hybrid cache and check the validity
		returnedVal, returnedExist := hybrid.Has(key)
		assert.Equal(t, bytes.Compare(returnedVal, value), 0)
		assert.Equal(t, returnedExist, true)
	}
}
