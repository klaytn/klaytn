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
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
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
	if err != nil {
		t.Fatal(err)
	}

	key, value := randBytes(32), randBytes(5*1024*1024) // 5MB value
	redis.Set(key, value)

	retValue := redis.Get(key)
	assert.Equal(t, bytes.Compare(value, retValue), 0)
}

// testNewRedisCache_Timeout test timout feature of redis client.
// INFO: Enable it just when you want to test.
func testNewRedisCache_Timeout(t *testing.T) {
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

	var redis TrieNodeCache = &RedisCache{redis.NewClient(&redis.Options{
		Addr:         "localhost:11234",
		DialTimeout:  redisCacheDialTimeout,
		ReadTimeout:  redisCacheTimeout,
		WriteTimeout: redisCacheTimeout,
		MaxRetries:   0,
	})}

	key, value := randBytes(32), randBytes(500)

	start := time.Now()
	redis.Set(key, value)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(time.Second))

	start = time.Now()
	_ = redis.Get(key)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(time.Second))

	start = time.Now()
	_, _ = redis.Has(key)
	assert.Equal(t, redisCacheTimeout, time.Since(start).Round(time.Second))
}
