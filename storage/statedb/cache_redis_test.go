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
	"github.com/stretchr/testify/assert"
)

func getTestRedisConfig() *TrieNodeCacheConfig {
	return &TrieNodeCacheConfig{
		CacheType:          CacheTypeRedis,
		LocalCacheSizeMB:   1024 * 1024,
		RedisEndpoints:     []string{"localhost:6379"},
		RedisClusterEnable: false,
	}
}

func TestSubscription(t *testing.T) {
	msg1 := "testMessage1"
	msg2 := "testMessage2"

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		cache, err := NewRedisCache(getTestRedisConfig())
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
	time.Sleep(100 * time.Millisecond)

	cache, err := NewRedisCache(getTestRedisConfig())
	assert.Nil(t, err)

	if err := cache.PublishBlock(msg1); err != nil {
		t.Fatal(err)
	}

	if err := cache.PublishBlock(msg2); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

// TestNewRedisCache tests basic operations of redis cache
func TestNewRedisCache(t *testing.T) {
	cache, err := NewRedisCache(getTestRedisConfig())
	assert.Nil(t, err)

	key, value := randBytes(32), randBytes(500)
	cache.Set(key, value)

	getValue := cache.Get(key)
	assert.Equal(t, bytes.Compare(value, getValue), 0)

	hasValue, ok := cache.Has(key)
	assert.Equal(t, ok, true)
	assert.Equal(t, bytes.Compare(value, hasValue), 0)
}

// TestNewRedisCache_Set_LargeData check whether redis cache can store an large data (5MB).
func TestNewRedisCache_Set_LargeData(t *testing.T) {
	cache, err := NewRedisCache(getTestRedisConfig())
	if err != nil {
		t.Fatal(err)
	}

	key, value := randBytes(32), randBytes(5*1024*1024) // 5MB value
	cache.Set(key, value)

	retValue := cache.Get(key)
	assert.Equal(t, bytes.Compare(value, retValue), 0)
}

// TestNewRedisCache_Timeout test timout feature of redis client.
func TestNewRedisCache_Timeout(t *testing.T) {
	go func() {
		tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:11234")
		if err != nil {
			t.Fatal(err)
		}

		listen, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer listen.Close()

		for {
			if err := listen.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
				t.Fatal(err)
			}
			_, err := listen.AcceptTCP()
			if err != nil {
				if strings.Contains(err.Error(), "timeout") {
					return
				}
				t.Fatal(err)
			}
		}
	}()

	var cache TrieNodeCache = &RedisCache{redis.NewClient(&redis.Options{
		Addr:         "localhost:11234",
		DialTimeout:  redisCacheDialTimeout,
		ReadTimeout:  redisCacheTimeout,
		WriteTimeout: redisCacheTimeout,
		MaxRetries:   0,
	}), nil}

	key, value := randBytes(32), randBytes(500)

	start := time.Now()
	cache.Set(key, value)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(redisCacheTimeout/2))

	start = time.Now()
	_ = cache.Get(key)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(redisCacheTimeout/2))

	start = time.Now()
	_, _ = cache.Has(key)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(redisCacheTimeout/2))
}
