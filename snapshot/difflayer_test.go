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
	"bytes"
	"math/rand"
	"testing"

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

// TestMergeBasics tests some simple merges
func TestMergeBasics(t *testing.T) {
	var (
		destructs = make(map[common.Hash]struct{})
		accounts  = make(map[common.Hash][]byte)
		storage   = make(map[common.Hash]map[common.Hash][]byte)
	)
	// Fill up a parent
	for i := 0; i < 100; i++ {
		h := randomHash()
		data := randomAccount()

		accounts[h] = data
		if rand.Intn(4) == 0 {
			destructs[h] = struct{}{}
		}
		if rand.Intn(2) == 0 {
			accStorage := make(map[common.Hash][]byte)
			value := make([]byte, 32)
			rand.Read(value)
			accStorage[randomHash()] = value
			storage[h] = accStorage
		}
	}
	// Add some (identical) layers on top
	parent := newDiffLayer(emptyLayer(), common.Hash{}, copyDestructs(destructs), copyAccounts(accounts), copyStorage(storage))
	child := newDiffLayer(parent, common.Hash{}, copyDestructs(destructs), copyAccounts(accounts), copyStorage(storage))
	child = newDiffLayer(child, common.Hash{}, copyDestructs(destructs), copyAccounts(accounts), copyStorage(storage))
	child = newDiffLayer(child, common.Hash{}, copyDestructs(destructs), copyAccounts(accounts), copyStorage(storage))
	child = newDiffLayer(child, common.Hash{}, copyDestructs(destructs), copyAccounts(accounts), copyStorage(storage))
	// And flatten
	merged := (child.flatten()).(*diffLayer)

	{ // Check account lists
		if have, want := len(merged.accountList), 0; have != want {
			t.Errorf("accountList wrong: have %v, want %v", have, want)
		}
		if have, want := len(merged.AccountList()), len(accounts); have != want {
			t.Errorf("AccountList() wrong: have %v, want %v", have, want)
		}
		if have, want := len(merged.accountList), len(accounts); have != want {
			t.Errorf("accountList [2] wrong: have %v, want %v", have, want)
		}
	}
	{ // Check account drops
		if have, want := len(merged.destructSet), len(destructs); have != want {
			t.Errorf("accountDrop wrong: have %v, want %v", have, want)
		}
	}
	{ // Check storage lists
		i := 0
		for aHash, sMap := range storage {
			if have, want := len(merged.storageList), i; have != want {
				t.Errorf("[1] storageList wrong: have %v, want %v", have, want)
			}
			list, _ := merged.StorageList(aHash)
			if have, want := len(list), len(sMap); have != want {
				t.Errorf("[2] StorageList() wrong: have %v, want %v", have, want)
			}
			if have, want := len(merged.storageList[aHash]), len(sMap); have != want {
				t.Errorf("storageList wrong: have %v, want %v", have, want)
			}
			i++
		}
	}
}

// TestMergeDelete tests some deletion
func TestMergeDelete(t *testing.T) {
	storage := make(map[common.Hash]map[common.Hash][]byte)
	// Fill up a parent
	h1 := common.HexToHash("0x01")
	h2 := common.HexToHash("0x02")

	flipDrops := func() map[common.Hash]struct{} {
		return map[common.Hash]struct{}{
			h2: {},
		}
	}
	flipAccs := func() map[common.Hash][]byte {
		return map[common.Hash][]byte{
			h1: randomAccount(),
		}
	}
	flopDrops := func() map[common.Hash]struct{} {
		return map[common.Hash]struct{}{
			h1: {},
		}
	}
	flopAccs := func() map[common.Hash][]byte {
		return map[common.Hash][]byte{
			h2: randomAccount(),
		}
	}
	// Add some flipAccs-flopping layers on top
	parent := newDiffLayer(emptyLayer(), common.Hash{}, flipDrops(), flipAccs(), storage)
	child := parent.Update(common.Hash{}, flopDrops(), flopAccs(), storage)
	child = child.Update(common.Hash{}, flipDrops(), flipAccs(), storage)
	child = child.Update(common.Hash{}, flopDrops(), flopAccs(), storage)
	child = child.Update(common.Hash{}, flipDrops(), flipAccs(), storage)
	child = child.Update(common.Hash{}, flopDrops(), flopAccs(), storage)
	child = child.Update(common.Hash{}, flipDrops(), flipAccs(), storage)

	if data, _ := child.Account(h1); data == nil {
		t.Errorf("last diff layer: expected %x account to be non-nil", h1)
	}
	if data, _ := child.Account(h2); data != nil {
		t.Errorf("last diff layer: expected %x account to be nil", h2)
	}
	if _, ok := child.destructSet[h1]; ok {
		t.Errorf("last diff layer: expected %x drop to be missing", h1)
	}
	if _, ok := child.destructSet[h2]; !ok {
		t.Errorf("last diff layer: expected %x drop to be present", h1)
	}
	// And flatten
	merged := (child.flatten()).(*diffLayer)

	if data, _ := merged.Account(h1); data == nil {
		t.Errorf("merged layer: expected %x account to be non-nil", h1)
	}
	if data, _ := merged.Account(h2); data != nil {
		t.Errorf("merged layer: expected %x account to be nil", h2)
	}
	if _, ok := merged.destructSet[h1]; !ok { // Note, drops stay alive until persisted to disk!
		t.Errorf("merged diff layer: expected %x drop to be present", h1)
	}
	if _, ok := merged.destructSet[h2]; !ok { // Note, drops stay alive until persisted to disk!
		t.Errorf("merged diff layer: expected %x drop to be present", h1)
	}
	// If we add more granular metering of memory, we can enable this again,
	// but it's not implemented for now
	//if have, want := merged.memory, child.memory; have != want {
	//	t.Errorf("mem wrong: have %d, want %d", have, want)
	//}
}

// This tests that if we create a new account, and set a slot, and then merge
// it, the lists will be correct.
func TestInsertAndMerge(t *testing.T) {
	// Fill up a parent
	var (
		acc    = common.HexToHash("0x01")
		slot   = common.HexToHash("0x02")
		parent *diffLayer
		child  *diffLayer
	)
	{
		var (
			destructs = make(map[common.Hash]struct{})
			accounts  = make(map[common.Hash][]byte)
			storage   = make(map[common.Hash]map[common.Hash][]byte)
		)
		parent = newDiffLayer(emptyLayer(), common.Hash{}, destructs, accounts, storage)
	}
	{
		var (
			destructs = make(map[common.Hash]struct{})
			accounts  = make(map[common.Hash][]byte)
			storage   = make(map[common.Hash]map[common.Hash][]byte)
		)
		accounts[acc] = randomAccount()
		storage[acc] = make(map[common.Hash][]byte)
		storage[acc][slot] = []byte{0x01}
		child = newDiffLayer(parent, common.Hash{}, destructs, accounts, storage)
	}
	// And flatten
	merged := (child.flatten()).(*diffLayer)
	{ // Check that slot value is present
		have, _ := merged.Storage(acc, slot)
		if want := []byte{0x01}; !bytes.Equal(have, want) {
			t.Errorf("merged slot value wrong: have %x, want %x", have, want)
		}
	}
}

func emptyLayer() *diskLayer {
	return &diskLayer{
		diskdb: database.NewMemoryDBManager(),
		cache:  fastcache.New(500 * 1024),
	}
}
