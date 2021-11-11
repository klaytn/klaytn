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
// This file is derived from core/state/snapshot/iterator_test.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/klaytn/klaytn/common"
)

// TestAccountIteratorBasics tests some simple single-layer(diff and disk) iteration
func TestAccountIteratorBasics(t *testing.T) {
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
	diffLayer := newDiffLayer(emptyLayer(), common.Hash{}, copyDestructs(destructs), copyAccounts(accounts), copyStorage(storage))
	it := diffLayer.AccountIterator(common.Hash{})
	verifyIterator(t, 100, it, verifyNothing) // Nil is allowed for single layer iterator

	diskLayer := diffToDisk(diffLayer)
	it = diskLayer.AccountIterator(common.Hash{})
	verifyIterator(t, 100, it, verifyNothing) // Nil is allowed for single layer iterator
}

// TestStorageIteratorBasics tests some simple single-layer(diff and disk) iteration for storage
func TestStorageIteratorBasics(t *testing.T) {
	var (
		nilStorage = make(map[common.Hash]int)
		accounts   = make(map[common.Hash][]byte)
		storage    = make(map[common.Hash]map[common.Hash][]byte)
	)
	// Fill some random data
	for i := 0; i < 10; i++ {
		h := randomHash()
		accounts[h] = randomAccount()

		accStorage := make(map[common.Hash][]byte)
		value := make([]byte, 32)

		var nilstorage int
		for i := 0; i < 100; i++ {
			rand.Read(value)
			if rand.Intn(2) == 0 {
				accStorage[randomHash()] = common.CopyBytes(value)
			} else {
				accStorage[randomHash()] = nil // delete slot
				nilstorage += 1
			}
		}
		storage[h] = accStorage
		nilStorage[h] = nilstorage
	}
	// Add some (identical) layers on top
	diffLayer := newDiffLayer(emptyLayer(), common.Hash{}, nil, copyAccounts(accounts), copyStorage(storage))
	for account := range accounts {
		it, _ := diffLayer.StorageIterator(account, common.Hash{})
		verifyIterator(t, 100, it, verifyNothing) // Nil is allowed for single layer iterator
	}

	diskLayer := diffToDisk(diffLayer)
	for account := range accounts {
		it, _ := diskLayer.StorageIterator(account, common.Hash{})
		verifyIterator(t, 100-nilStorage[account], it, verifyNothing) // Nil is allowed for single layer iterator
	}
}

type testIterator struct {
	values []byte
}

func newTestIterator(values ...byte) *testIterator {
	return &testIterator{values}
}

func (ti *testIterator) Seek(common.Hash) {
	panic("implement me")
}

func (ti *testIterator) Next() bool {
	ti.values = ti.values[1:]
	return len(ti.values) > 0
}

func (ti *testIterator) Error() error {
	return nil
}

func (ti *testIterator) Hash() common.Hash {
	return common.BytesToHash([]byte{ti.values[0]})
}

func (ti *testIterator) Account() []byte {
	return nil
}

func (ti *testIterator) Slot() []byte {
	return nil
}

func (ti *testIterator) Release() {}

func TestFastIteratorBasics(t *testing.T) {
	type testCase struct {
		lists   [][]byte
		expKeys []byte
	}
	for i, tc := range []testCase{
		{lists: [][]byte{{0, 1, 8}, {1, 2, 8}, {2, 9}, {4},
			{7, 14, 15}, {9, 13, 15, 16}},
			expKeys: []byte{0, 1, 2, 4, 7, 8, 9, 13, 14, 15, 16}},
		{lists: [][]byte{{0, 8}, {1, 2, 8}, {7, 14, 15}, {8, 9},
			{9, 10}, {10, 13, 15, 16}},
			expKeys: []byte{0, 1, 2, 7, 8, 9, 10, 13, 14, 15, 16}},
	} {
		var iterators []*weightedIterator
		for i, data := range tc.lists {
			it := newTestIterator(data...)
			iterators = append(iterators, &weightedIterator{it, i})
		}
		fi := &fastIterator{
			iterators: iterators,
			initiated: false,
		}
		count := 0
		for fi.Next() {
			if got, exp := fi.Hash()[31], tc.expKeys[count]; exp != got {
				t.Errorf("tc %d, [%d]: got %d exp %d", i, count, got, exp)
			}
			count++
		}
	}
}

type verifyContent int

const (
	verifyNothing verifyContent = iota
	verifyAccount
	verifyStorage
)

func verifyIterator(t *testing.T, expCount int, it Iterator, verify verifyContent) {
	t.Helper()

	var (
		count = 0
		last  = common.Hash{}
	)
	for it.Next() {
		hash := it.Hash()
		if bytes.Compare(last[:], hash[:]) >= 0 {
			t.Errorf("wrong order: %x >= %x", last, hash)
		}
		count++
		if verify == verifyAccount && len(it.(AccountIterator).Account()) == 0 {
			t.Errorf("iterator returned nil-value for hash %x", hash)
		} else if verify == verifyStorage && len(it.(StorageIterator).Slot()) == 0 {
			t.Errorf("iterator returned nil-value for hash %x", hash)
		}
		last = hash
	}
	if count != expCount {
		t.Errorf("iterator count mismatch: have %d, want %d", count, expCount)
	}
	if err := it.Error(); err != nil {
		t.Errorf("iterator failed: %v", err)
	}
}
