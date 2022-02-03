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
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
)

// TODO-snapshot remove the following function after porting diffToDisk method
func createTestDiskLayer(db database.DBManager, accounts map[common.Hash][]byte, storages map[common.Hash]map[common.Hash][]byte) *diskLayer {
	batch := db.NewSnapshotDBBatch()

	for hash, data := range accounts {
		batch.WriteAccountSnapshot(hash, data)
		batch.Write()
		batch.Reset()
	}

	for accountHash, storage := range storages {
		for storageHash, data := range storage {
			if len(data) > 0 {
				batch.WriteStorageSnapshot(accountHash, storageHash, data)
			} else {
				batch.DeleteStorageSnapshot(accountHash, storageHash)
			}
			batch.Write()
			batch.Reset()
		}
	}

	return &diskLayer{
		diskdb: db,
	}
}

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
	// TODO-snapshot uncomment after porting difflayer.go
	//// Add some (identical) layers on top
	//diffLayer := newDiffLayer(emptyLayer(), common.Hash{}, copyDestructs(destructs), copyAccounts(accounts), copyStorage(storage))
	//it := diffLayer.AccountIterator(common.Hash{})
	//verifyIterator(t, 100, it, verifyNothing) // Nil is allowed for single layer iterator
	//
	//diskLayer := diffToDisk(diffLayer)

	// TODO-snapshot the following disklayer creation is temporary and will be removed after porting difflayer.go
	dir, err := ioutil.TempDir("", "klaytn-test-snapshot-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)
	dbc := &database.DBConfig{Dir: dir, DBType: database.LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	db := database.NewDBManager(dbc)
	diskLayer := createTestDiskLayer(db, accounts, storage)

	it := diskLayer.AccountIterator(common.Hash{})
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
	// TODO-snapshot uncomment after porting difflayer.go
	//// Add some (identical) layers on top
	//diffLayer := newDiffLayer(emptyLayer(), common.Hash{}, nil, copyAccounts(accounts), copyStorage(storage))
	//for account := range accounts {
	//	it, _ := diffLayer.StorageIterator(account, common.Hash{})
	//	verifyIterator(t, 100, it, verifyNothing) // Nil is allowed for single layer iterator
	//}

	// TODO-snapshot the following disklayer creation is temporary and will be removed after porting difflayer.go
	dir, err := ioutil.TempDir("", "klaytn-test-snapshot-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)
	dbc := &database.DBConfig{Dir: dir, DBType: database.LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	db := database.NewDBManager(dbc)
	diskLayer := createTestDiskLayer(db, accounts, storage)

	for account := range accounts {
		it, _ := diskLayer.StorageIterator(account, common.Hash{})
		verifyIterator(t, 100-nilStorage[account], it, verifyNothing) // Nil is allowed for single layer iterator
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
