// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/state/sync_test.go (2020/05/20).
// Modified and improved for the klaytn development.

package state

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/alecthomas/units"
	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/stretchr/testify/assert"
)

// testAccount is the data associated with an account used by the state tests.
type testAccount struct {
	address    common.Address
	balance    *big.Int
	nonce      uint64
	code       []byte
	storageMap map[common.Hash]common.Hash
}

// makeTestState create a sample test state to test node-wise reconstruction.
func makeTestState(t *testing.T) (Database, common.Hash, []*testAccount) {
	// Create an empty state
	db := NewDatabase(database.NewMemoryDBManager())
	statedb, err := New(common.Hash{}, db, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Fill it with some arbitrary data
	var accounts []*testAccount
	for i := byte(0); i < 96; i++ {
		var obj *stateObject
		acc := &testAccount{
			address:    common.BytesToAddress([]byte{i}),
			storageMap: make(map[common.Hash]common.Hash),
		}

		if i%3 > 0 {
			obj = statedb.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		} else {
			obj = statedb.GetOrNewSmartContract(common.BytesToAddress([]byte{i}))

			obj.SetCode(crypto.Keccak256Hash([]byte{i, i, i, i, i}), []byte{i, i, i, i, i})
			acc.code = []byte{i, i, i, i, i}
			if i == 0 {
				// to test emptyCodeHash
				obj.SetCode(crypto.Keccak256Hash([]byte{}), []byte{})
				acc.code = []byte{}
			}

			for j := 0; j < int(i)%10; j++ {
				key := common.Hash{i + byte(j)}
				value := common.Hash{i*2 + 1}
				acc.storageMap[key] = value

				obj.SetState(db, key, value)
			}
		}

		obj.AddBalance(big.NewInt(int64(11 * i)))
		acc.balance = big.NewInt(int64(11 * i))

		obj.SetNonce(uint64(42 * i))
		acc.nonce = uint64(42 * i)

		statedb.updateStateObject(obj)
		accounts = append(accounts, acc)
	}
	root, _ := statedb.Commit(false)

	// commit stateTrie to DB
	statedb.db.TrieDB().Commit(root, false, 0)

	if err := checkStateConsistency(db.TrieDB().DiskDB(), root); err != nil {
		t.Fatalf("inconsistent state trie at %x: %v", root, err)
	}

	// Return the generated state
	return db, root, accounts
}

// checkStateAccounts cross references a reconstructed state with an expected
// account array.
func checkStateAccounts(t *testing.T, newDB database.DBManager, root common.Hash, accounts []*testAccount) {
	// Check root availability and state contents
	state, err := New(root, NewDatabase(newDB), nil)
	if err != nil {
		t.Fatalf("failed to create state trie at %x: %v", root, err)
	}
	if err := checkStateConsistency(newDB, root); err != nil {
		t.Fatalf("inconsistent state trie at %x: %v", root, err)
	}
	for i, acc := range accounts {
		if balance := state.GetBalance(acc.address); balance.Cmp(acc.balance) != 0 {
			t.Errorf("account %d: balance mismatch: have %v, want %v", i, balance, acc.balance)
		}
		if nonce := state.GetNonce(acc.address); nonce != acc.nonce {
			t.Errorf("account %d: nonce mismatch: have %v, want %v", i, nonce, acc.nonce)
		}
		if code := state.GetCode(acc.address); !bytes.Equal(code, acc.code) {
			t.Errorf("account %d: code mismatch: have %x, want %x", i, code, acc.code)
		}

		// check storage trie
		st := state.StorageTrie(acc.address)
		it := statedb.NewIterator(st.NodeIterator(nil))
		storageMapWithHashedKey := make(map[common.Hash]common.Hash)
		for it.Next() {
			storageMapWithHashedKey[common.BytesToHash(it.Key)] = common.BytesToHash(it.Value)
		}
		if len(storageMapWithHashedKey) != len(acc.storageMap) {
			t.Errorf("account %d: stroage trie number mismatch: have %x, want %x", i, len(storageMapWithHashedKey), len(acc.storageMap))
		}
		for key, value := range acc.storageMap {
			hk := crypto.Keccak256Hash(key[:])
			if storageMapWithHashedKey[hk] != value {
				t.Errorf("account %d: stroage trie (%v) mismatch: have %x, want %x", i, key.String(), acc.storageMap[key], value)
			}
		}
	}
}

// checkTrieConsistency checks that all nodes in a (sub-)trie are indeed present.
func checkTrieConsistency(db database.DBManager, root common.Hash) error {
	if v, _ := db.ReadStateTrieNode(root[:]); v == nil {
		return nil // Consider a non existent state consistent.
	}
	trie, err := statedb.NewTrie(root, statedb.NewDatabase(db))
	if err != nil {
		return err
	}
	it := trie.NodeIterator(nil)
	for it.Next(true) {
	}
	return it.Error()
}

// checkStateConsistency checks that all data of a state root is present.
func checkStateConsistency(db database.DBManager, root common.Hash) error {
	// Create and iterate a state trie rooted in a sub-node
	if _, err := db.ReadStateTrieNode(root.Bytes()); err != nil {
		return nil // Consider a non existent state consistent.
	}
	state, err := New(root, NewDatabase(db), nil)
	if err != nil {
		return err
	}
	it := NewNodeIterator(state)
	for it.Next() {
	}
	return it.Error
}

// Tests that an empty state is not scheduled for syncing.
func TestEmptyStateSync(t *testing.T) {
	empty := common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// only bloom
	{
		bloom := statedb.NewSyncBloom(1, database.NewMemDB())
		if req := NewStateSync(empty, database.NewMemoryDBManager(), bloom, nil).Missing(1); len(req) != 0 {
			t.Errorf("content requested for empty state: %v", req)
		}
	}

	// only lru
	{
		lruCache, _ := lru.New(int(1 * units.MB / common.HashLength))
		if req := NewStateSync(empty, database.NewMemoryDBManager(), nil, lruCache).Missing(1); len(req) != 0 {
			t.Errorf("content requested for empty state: %v", req)
		}
	}

	// no bloom lru
	{
		if req := NewStateSync(empty, database.NewMemoryDBManager(), nil, nil).Missing(1); len(req) != 0 {
			t.Errorf("content requested for empty state: %v", req)
		}
	}

	// both bloom, lru
	{
		bloom := statedb.NewSyncBloom(1, database.NewMemDB())
		lruCache, _ := lru.New(int(1 * units.MB / common.HashLength))
		if req := NewStateSync(empty, database.NewMemoryDBManager(), bloom, lruCache).Missing(1); len(req) != 0 {
			t.Errorf("content requested for empty state: %v", req)
		}
	}
}

// Tests that given a root hash, a state can sync iteratively on a single thread,
// requesting retrieval tasks and returning all of them in one go.
func TestIterativeStateSyncIndividual(t *testing.T) { testIterativeStateSync(t, 1) }
func TestIterativeStateSyncBatched(t *testing.T)    { testIterativeStateSync(t, 100) }

func testIterativeStateSync(t *testing.T, count int) {
	// Create a random state to copy
	srcState, srcRoot, srcAccounts := makeTestState(t)

	// Create a destination state and sync with the scheduler
	dstDiskDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDiskDb)
	sched := NewStateSync(srcRoot, dstDiskDb, statedb.NewSyncBloom(1, dstDiskDb.GetMemDB()), nil)

	queue := append([]common.Hash{}, sched.Missing(count)...)
	for len(queue) > 0 {
		results := make([]statedb.SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = statedb.SyncResult{Hash: hash, Data: data}
		}
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		batch := dstDiskDb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		queue = append(queue[:0], sched.Missing(count)...)
	}
	// Cross check that the two states are in sync
	checkStateAccounts(t, dstDiskDb, srcRoot, srcAccounts)

	err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
	assert.NoError(t, err)

	// Test with quit channel
	quit := make(chan struct{})

	// normal
	err = CheckStateConsistency(srcState, dstState, srcRoot, 100, quit)
	assert.NoError(t, err)

	// quit
	close(quit)
	err = CheckStateConsistency(srcState, dstState, srcRoot, 100, quit)
	assert.Error(t, err, errStopByQuit)
}

func TestCheckStateConsistencyMissNode(t *testing.T) {
	// Create a random state to copy
	srcState, srcRoot, _ := makeTestState(t)
	newState, _, _ := makeTestState(t)

	srcStateDB, err := New(srcRoot, srcState, nil)
	assert.NoError(t, err)

	it := NewNodeIterator(srcStateDB)
	it.Next() // skip trie root node

	for it.Next() {
		if !common.EmptyHash(it.Hash) {
			hash := it.Hash
			data, _ := srcStateDB.db.TrieDB().DiskDB().ReadStateTrieNode(hash[:])

			// Remove nodes
			srcState.TrieDB().DiskDB().GetMemDB().Delete(hash[:])
			newState.TrieDB().DiskDB().GetMemDB().Delete(hash[:])

			// Check consistency : errIterator
			err = CheckStateConsistency(srcState, newState, srcRoot, 100, nil)
			if !errors.Is(err, errIterator) {
				t.Log("mismatched err", "err", err, "expErr", errIterator)
				t.FailNow()
			}

			// Recover nodes
			srcState.TrieDB().DiskDB().GetMemDB().Put(hash[:], data)
			newState.TrieDB().DiskDB().GetMemDB().Put(hash[:], data)
		}
	}

	// Check consistency : no error
	err = CheckStateConsistency(srcState, newState, srcRoot, 100, nil)
	assert.NoError(t, err)

	err = CheckStateConsistencyParallel(srcState, newState, srcRoot, nil)
	assert.NoError(t, err)
}

// Tests that the trie scheduler can correctly reconstruct the state even if only
// partial results are returned, and the others sent only later.
func TestIterativeDelayedStateSync(t *testing.T) {
	// Create a random state to copy
	srcState, srcRoot, srcAccounts := makeTestState(t)

	// Create a destination state and sync with the scheduler
	dstDiskDB := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDiskDB)
	sched := NewStateSync(srcRoot, dstDiskDB, statedb.NewSyncBloom(1, dstDiskDB.GetMemDB()), nil)

	queue := append([]common.Hash{}, sched.Missing(0)...)
	for len(queue) > 0 {
		// Sync only half of the scheduled nodes
		results := make([]statedb.SyncResult, len(queue)/2+1)
		for i, hash := range queue[:len(results)] {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = statedb.SyncResult{Hash: hash, Data: data}
		}
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		batch := dstDiskDB.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		queue = append(queue[len(results):], sched.Missing(0)...)
	}
	// Cross check that the two states are in sync
	checkStateAccounts(t, dstDiskDB, srcRoot, srcAccounts)

	err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
	assert.NoError(t, err)

	err = CheckStateConsistencyParallel(srcState, dstState, srcRoot, nil)
	assert.NoError(t, err)
}

// Tests that given a root hash, a trie can sync iteratively on a single thread,
// requesting retrieval tasks and returning all of them in one go, however in a
// random order.
func TestIterativeRandomStateSyncIndividual(t *testing.T) { testIterativeRandomStateSync(t, 1) }
func TestIterativeRandomStateSyncBatched(t *testing.T)    { testIterativeRandomStateSync(t, 100) }

func testIterativeRandomStateSync(t *testing.T, count int) {
	// Create a random state to copy
	srcState, srcRoot, srcAccounts := makeTestState(t)

	// Create a destination state and sync with the scheduler
	dstDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDb)
	sched := NewStateSync(srcRoot, dstDb, statedb.NewSyncBloom(1, dstDb.GetMemDB()), nil)

	queue := make(map[common.Hash]struct{})
	for _, hash := range sched.Missing(count) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
		// Fetch all the queued nodes in a random order
		results := make([]statedb.SyncResult, 0, len(queue))
		for hash := range queue {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results = append(results, statedb.SyncResult{Hash: hash, Data: data})
		}
		// Feed the retrieved results back and queue new tasks
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		batch := dstDb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		queue = make(map[common.Hash]struct{})
		for _, hash := range sched.Missing(count) {
			queue[hash] = struct{}{}
		}
	}
	// Cross check that the two states are in sync
	checkStateAccounts(t, dstDb, srcRoot, srcAccounts)

	err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
	assert.NoError(t, err)

	err = CheckStateConsistencyParallel(srcState, dstState, srcRoot, nil)
	assert.NoError(t, err)
}

// Tests that the trie scheduler can correctly reconstruct the state even if only
// partial results are returned (Even those randomly), others sent only later.
func TestIterativeRandomDelayedStateSync(t *testing.T) {
	// Create a random state to copy
	srcState, srcRoot, srcAccounts := makeTestState(t)

	// Create a destination state and sync with the scheduler
	dstDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDb)
	sched := NewStateSync(srcRoot, dstDb, statedb.NewSyncBloom(1, dstDb.GetMemDB()), nil)

	queue := make(map[common.Hash]struct{})
	for _, hash := range sched.Missing(0) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
		// Sync only half of the scheduled nodes, even those in random order
		results := make([]statedb.SyncResult, 0, len(queue)/2+1)
		for hash := range queue {
			delete(queue, hash)

			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results = append(results, statedb.SyncResult{Hash: hash, Data: data})

			if len(results) >= cap(results) {
				break
			}
		}
		// Feed the retrieved results back and queue new tasks
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		batch := dstDb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		for _, hash := range sched.Missing(0) {
			queue[hash] = struct{}{}
		}
	}
	// Cross check that the two states are in sync
	checkStateAccounts(t, dstDb, srcRoot, srcAccounts)

	err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
	assert.NoError(t, err)

	err = CheckStateConsistencyParallel(srcState, dstState, srcRoot, nil)
	assert.NoError(t, err)
}

// Tests that at any point in time during a sync, only complete sub-tries are in
// the database.
func TestIncompleteStateSync(t *testing.T) {
	// Create a random state to copy
	srcState, srcRoot, srcAccounts := makeTestState(t)

	checkTrieConsistency(srcState.TrieDB().DiskDB().(database.DBManager), srcRoot)

	// Create a destination state and sync with the scheduler
	dstDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDb)
	sched := NewStateSync(srcRoot, dstDb, statedb.NewSyncBloom(1, dstDb.GetMemDB()), nil)

	var added []common.Hash
	queue := append([]common.Hash{}, sched.Missing(1)...)
	for len(queue) > 0 {
		// Fetch a batch of state nodes
		results := make([]statedb.SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = statedb.SyncResult{Hash: hash, Data: data}
		}
		// Process each of the state nodes
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		batch := dstDb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		for _, result := range results {
			added = append(added, result.Hash)
		}
		// Check that all known sub-tries added so far are complete or missing entirely.
	checkSubtries:
		for _, hash := range added {
			for _, acc := range srcAccounts {
				if hash == crypto.Keccak256Hash(acc.code) {
					continue checkSubtries // skip trie check of code nodes.
				}
			}
			// Can't use checkStateConsistency here because subtrie keys may have odd
			// length and crash in LeafKey.
			if err := checkTrieConsistency(dstDb, hash); err != nil {
				t.Fatalf("state inconsistent: %v", err)
			}
		}
		// Fetch the next batch to retrieve
		queue = append(queue[:0], sched.Missing(1)...)
	}
	// Sanity check that removing any node from the database is detected
	for _, node := range added[1:] {
		key := node.Bytes()
		value, _ := dstDb.GetMemDB().Get(key)

		dstDb.GetMemDB().Delete(key)
		if err := checkStateConsistency(dstDb, added[0]); err == nil {
			t.Fatalf("trie inconsistency not caught, missing: %x", key)
		}

		err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
		assert.Error(t, err)

		err = CheckStateConsistencyParallel(srcState, dstState, srcRoot, nil)
		assert.Error(t, err)

		dstDb.GetMemDB().Put(key, value)
	}

	err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
	assert.NoError(t, err)

	err = CheckStateConsistencyParallel(srcState, dstState, srcRoot, nil)
	assert.NoError(t, err)
}
