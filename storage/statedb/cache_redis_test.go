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

import (
	"bytes"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/klaytn/klaytn/storage"
	"github.com/stretchr/testify/assert"
)

const sleepDurationForAsyncBehavior = 100 * time.Millisecond

func getTestRedisConfig() *TrieNodeCacheConfig {
	return &TrieNodeCacheConfig{
		CacheType:          CacheTypeRedis,
		LocalCacheSizeMiB:  100,
		RedisEndpoints:     []string{"localhost:6379"},
		RedisClusterEnable: false,
	}
}

func TestSubscription(t *testing.T) {
	storage.SkipLocalTest(t)

	msg1 := "testMessage1"
	msg2 := "testMessage2"

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		cache, err := newRedisCache(getTestRedisConfig())
		assert.Nil(t, err)

		ch := cache.SubscribeBlockCh()

		select {
		case actualMsg := <-ch:
			assert.Equal(t, msg1, actualMsg.Payload)
		case <-time.After(time.Second):
			panic("timeout")
		}

		select {
		case actualMsg := <-ch:
			assert.Equal(t, msg2, actualMsg.Payload)
		case <-time.After(time.Second):
			panic("timeout")
		}

		wg.Done()
	}()
	time.Sleep(sleepDurationForAsyncBehavior)

	cache, err := newRedisCache(getTestRedisConfig())
	assert.Nil(t, err)

	if err := cache.PublishBlock(msg1); err != nil {
		t.Fatal(err)
	}

	if err := cache.PublishBlock(msg2); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

// TestRedisCache tests basic operations of redis cache
func TestRedisCache(t *testing.T) {
	storage.SkipLocalTest(t)

	cache, err := newRedisCache(getTestRedisConfig())
	assert.Nil(t, err)

	key, value := randBytes(32), randBytes(500)
	cache.Set(key, value)

	getValue := cache.Get(key)
	assert.Equal(t, bytes.Compare(value, getValue), 0)

	hasValue, ok := cache.Has(key)
	assert.Equal(t, ok, true)
	assert.Equal(t, bytes.Compare(value, hasValue), 0)
}

// TestRedisCache_Set_LargeData check whether redis cache can store an large data (5MB).
func TestRedisCache_Set_LargeData(t *testing.T) {
	storage.SkipLocalTest(t)

	cache, err := newRedisCache(getTestRedisConfig())
	if err != nil {
		t.Fatal(err)
	}

	key, value := randBytes(32), randBytes(5*1024*1024) // 5MB value
	cache.Set(key, value)

	retValue := cache.Get(key)
	assert.Equal(t, bytes.Compare(value, retValue), 0)
}

// TestRedisCache_SetAsync tests basic operations of redis cache using SetAsync instead of Set.
func TestRedisCache_SetAsync(t *testing.T) {
	storage.SkipLocalTest(t)

	cache, err := newRedisCache(getTestRedisConfig())
	assert.Nil(t, err)

	key, value := randBytes(32), randBytes(500)
	cache.SetAsync(key, value)
	time.Sleep(sleepDurationForAsyncBehavior)

	getValue := cache.Get(key)
	assert.Equal(t, bytes.Compare(value, getValue), 0)

	hasValue, ok := cache.Has(key)
	assert.Equal(t, ok, true)
	assert.Equal(t, bytes.Compare(value, hasValue), 0)
}

// TestRedisCache_SetAsync_LargeData check whether redis cache can store an large data asynchronously (5MB).
func TestRedisCache_SetAsync_LargeData(t *testing.T) {
	storage.SkipLocalTest(t)

	cache, err := newRedisCache(getTestRedisConfig())
	if err != nil {
		t.Fatal(err)
	}

	key, value := randBytes(32), randBytes(5*1024*1024) // 5MB value
	cache.SetAsync(key, value)
	time.Sleep(sleepDurationForAsyncBehavior)

	retValue := cache.Get(key)
	assert.Equal(t, bytes.Compare(value, retValue), 0)
}

// TestRedisCache_SetAsync_LargeNumberItems asynchronously sets lots of items exceeding channel size.
func TestRedisCache_SetAsync_LargeNumberItems(t *testing.T) {
	storage.SkipLocalTest(t)

	cache, err := newRedisCache(getTestRedisConfig())
	if err != nil {
		t.Fatal(err)
	}

	itemsLen := redisSetItemChannelSize * 2
	items := make([]setItem, itemsLen)
	for i := 0; i < itemsLen; i++ {
		items[i].key = randBytes(32)
		items[i].value = randBytes(500)
	}

	go func() {
		// wait for a while to avoid redis setItem channel full
		time.Sleep(sleepDurationForAsyncBehavior)

		for i := 0; i < itemsLen; i++ {
			if i == redisSetItemChannelSize {
				// sleep for a while because Set command can drop an item if setItem channel is full
				time.Sleep(2 * time.Second)
			}
			// set writes items asynchronously
			cache.SetAsync(items[i].key, items[i].value)
		}
	}()

	start := time.Now()
	for i := 0; i < itemsLen; i++ {
		// terminate if test lasts long
		if time.Since(start) > 5*time.Second {
			t.Fatalf("timeout checking %dth item", i+1)
		}

		v := cache.Get(items[i].key)
		if v == nil {
			// if the item is not set yet, wait and retry
			time.Sleep(sleepDurationForAsyncBehavior)
			i--
		} else {
			assert.Equal(t, v, items[i].value)
		}
	}
}

// TestRedisCache_Timeout tests timout feature of redis client.
func TestRedisCache_Timeout(t *testing.T) {
	storage.SkipLocalTest(t)

	go func() {
		tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:11234")
		if err != nil {
			t.Error(err)
			return
		}

		listen, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			t.Error(err)
			return
		}
		defer listen.Close()

		for {
			if err := listen.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
				t.Error(err)
				return
			}
			_, err := listen.AcceptTCP()
			if err != nil {
				if strings.Contains(err.Error(), "timeout") {
					return
				}
				t.Error(err)
				return
			}
		}
	}()

	var cache TrieNodeCache = &RedisCache{redis.NewClient(&redis.Options{
		Addr:         "localhost:11234",
		DialTimeout:  redisCacheDialTimeout,
		ReadTimeout:  redisCacheTimeout,
		WriteTimeout: redisCacheTimeout,
		MaxRetries:   0,
	}), nil, nil}

	key, value := randBytes(32), randBytes(500)

	start := time.Now()
	redisCache := cache.(*RedisCache) // Because RedisCache.Set writes item asynchronously, use RedisCache.set
	redisCache.Set(key, value)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(redisCacheTimeout/2))

	start = time.Now()
	_ = cache.Get(key)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(redisCacheTimeout/2))

	start = time.Now()
	_, _ = cache.Has(key)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(redisCacheTimeout/2))
}
