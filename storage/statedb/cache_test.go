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
	"reflect"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

// TODO-Klaytn: Enable tests when redis is prepared on CI

// TestNewCache tests creating all kinds of supported stateDB caches.
func _TestNewCache(t *testing.T) {
	testCases := []struct {
		cacheType    CacheType
		expectedType reflect.Type
	}{
		{"LocalCache", reflect.TypeOf(&FastCache{})},
		{"RemoteCache", reflect.TypeOf(&RedisCache{})},
		{"HybridCache", reflect.TypeOf(&hybridCache{})},
	}

	for _, tc := range testCases {
		cache, err := NewCache(tc.cacheType, testMaxBytes, testRedisEndpoints, testRedisCluster)
		assert.NilError(t, err)
		assert.Equal(t, reflect.TypeOf(cache), tc.expectedType)
	}
}
