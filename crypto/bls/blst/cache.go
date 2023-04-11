// Copyright 2023 The klaytn Authors
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

package blst

import (
	"sync"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
)

type blstComputeCache struct {
	mu             sync.Mutex
	init           bool
	publicKeyCache common.Cache // PublicKey Uncompress
	signatureCache common.Cache // Signature Uncompress
}

var computeCache = &blstComputeCache{}

func cacheKey(b []byte) common.CacheKey {
	return crypto.Keccak256Hash(b)
}

func initCacheLocked() {
	computeCache.mu.Lock()
	defer computeCache.mu.Unlock()
	if !computeCache.init {
		cacheConfig := common.LRUConfig{CacheSize: 200}

		computeCache.publicKeyCache = common.NewCache(cacheConfig)
		computeCache.signatureCache = common.NewCache(cacheConfig)

		computeCache.init = true
	}
}

func publicKeyCache() common.Cache {
	initCacheLocked()
	return computeCache.publicKeyCache
}

func signatureCache() common.Cache {
	initCacheLocked()
	return computeCache.signatureCache
}
