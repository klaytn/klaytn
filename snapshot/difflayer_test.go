// Modifications Copyright 2021 The klaytn Authors
// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/state/snapshot/difflayer_test.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
)

func copyDestructs(destructs map[common.Hash]struct{}) map[common.Hash]struct{} {
	copy := make(map[common.Hash]struct{})
	for hash := range destructs {
		copy[hash] = struct{}{}
	}
	return copy
}

func copyAccounts(accounts map[common.Hash][]byte) map[common.Hash][]byte {
	copy := make(map[common.Hash][]byte)
	for hash, blob := range accounts {
		copy[hash] = blob
	}
	return copy
}

func copyStorage(storage map[common.Hash]map[common.Hash][]byte) map[common.Hash]map[common.Hash][]byte {
	copy := make(map[common.Hash]map[common.Hash][]byte)
	for accHash, slots := range storage {
		copy[accHash] = make(map[common.Hash][]byte)
		for slotHash, blob := range slots {
			copy[accHash][slotHash] = blob
		}
	}
	return copy
}

func emptyLayer() *diskLayer {
	return &diskLayer{
		diskdb: database.NewMemoryDBManager(),
		cache:  fastcache.New(500 * 1024),
	}
}
