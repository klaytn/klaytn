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
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/rlp"
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
func makeTestState(t *testing.T) (Database, common.ExtHash, []*testAccount) {
	// Create an empty state
	db := NewDatabase(database.NewMemoryDBManager())
	statedb, err := New(common.InitExtHash(), db, nil)
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

			obj.SetCode(crypto.Keccak256Hash([]byte{i, i, i, i, i}).ToRootExtHash(), []byte{i, i, i, i, i})
			acc.code = []byte{i, i, i, i, i}
			if i == 0 {
				// to test emptyCodeHash
				obj.SetCode(crypto.Keccak256Hash([]byte{}).ToRootExtHash(), []byte{})
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

	if err := checkStateConsistency(db.TrieDB().DiskDB(), root); err != nil {
		t.Fatalf("inconsistent state trie at %x: %v", root, err)
	}

	// Return the generated state
	return db, root, accounts
}

// checkStateAccounts cross references a reconstructed state with an expected
// account array.
func checkStateAccounts(t *testing.T, newDB database.DBManager, root common.ExtHash, accounts []*testAccount) {
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
func checkTrieConsistency(db database.DBManager, root common.ExtHash) error {
	if v, _ := db.ReadStateTrieNode(root.ToHash().Bytes()); v == nil {
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
func checkStateConsistency(db database.DBManager, root common.ExtHash) error {
	// Create and iterate a state trie rooted in a sub-node
	if _, err := db.ReadStateTrieNode(root.ToHash().Bytes()); err != nil {
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
		db := database.NewMemoryDBManager()
		sync := NewStateSync(empty, db, statedb.NewSyncBloom(1, db.GetMemDB()), nil, nil)
		if nodes, paths, codes := sync.Missing(1); len(nodes) != 0 || len(paths) != 0 || len(codes) != 0 {
			t.Errorf("content requested for empty state: %v", sync)
		}
	}

	// only lru
	{
		lruCache, _ := lru.New(int(1 * units.MB / common.HashLength))
		db := database.NewMemoryDBManager()
		sync := NewStateSync(empty, db, statedb.NewSyncBloom(1, db.GetMemDB()), lruCache, nil)
		if nodes, paths, codes := sync.Missing(1); len(nodes) != 0 || len(paths) != 0 || len(codes) != 0 {
			t.Errorf("content requested for empty state: %v", sync)
		}
	}

	// no bloom lru
	{
		db := database.NewMemoryDBManager()
		sync := NewStateSync(empty, db, nil, nil, nil)
		if nodes, paths, codes := sync.Missing(1); len(nodes) != 0 || len(paths) != 0 || len(codes) != 0 {
			t.Errorf("content requested for empty state: %v", sync)
		}
	}

	// both bloom, lru
	{
		bloom := statedb.NewSyncBloom(1, database.NewMemDB())
		lruCache, _ := lru.New(int(1 * units.MB / common.HashLength))
		db := database.NewMemoryDBManager()
		sync := NewStateSync(empty, db, bloom, lruCache, nil)
		if nodes, paths, codes := sync.Missing(1); len(nodes) != 0 || len(paths) != 0 || len(codes) != 0 {
			t.Errorf("content requested for empty state: %v", sync)
		}
	}
}

// Tests that given a root hash, a state can sync iteratively on a single thread,
// requesting retrieval tasks and returning all of them in one go.
func TestIterativeStateSyncIndividual(t *testing.T) {
	testIterativeStateSync(t, 1, false, false)
}

func TestIterativeStateSyncBatched(t *testing.T) {
	testIterativeStateSync(t, 100, false, false)
}

func TestIterativeStateSyncIndividualFromDisk(t *testing.T) {
	testIterativeStateSync(t, 1, true, false)
}

func TestIterativeStateSyncBatchedFromDisk(t *testing.T) {
	testIterativeStateSync(t, 100, true, false)
}

func TestIterativeStateSyncIndividualByPath(t *testing.T) {
	testIterativeStateSync(t, 1, false, true)
}

func TestIterativeStateSyncBatchedByPath(t *testing.T) {
	testIterativeStateSync(t, 100, false, true)
}

func testIterativeStateSync(t *testing.T, count int, commit bool, bypath bool) {
	// Create a random state to copy
	srcState, srcRoot, srcAccounts := makeTestState(t)
	if commit {
		srcState.TrieDB().Commit(srcRoot, false, 0)
	}
	srcTrie, _ := statedb.NewTrie(srcRoot, srcState.TrieDB())

	// Create a destination state and sync with the scheduler
	dstDiskDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDiskDb)
	sched := NewStateSync(srcRoot.ToHash(), dstDiskDb, statedb.NewSyncBloom(1, dstDiskDb.GetMemDB()), nil, nil)

	nodes, paths, codes := sched.Missing(count)
	var (
		hashQueue []common.ExtHash
		pathQueue []statedb.SyncPath
	)
	if !bypath {
		hashQueue = append(append(hashQueue[:0], nodes...), codes...)
	} else {
		hashQueue = append(hashQueue[:0], codes...)
		pathQueue = append(pathQueue[:0], paths...)
	}
	for len(hashQueue)+len(pathQueue) > 0 {
		results := make([]statedb.SyncResult, len(hashQueue)+len(pathQueue))
		for i, hash := range hashQueue {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				data, err = srcState.ContractCode(hash)
			}
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = statedb.SyncResult{Hash: hash, Data: data}
		}
		for i, path := range pathQueue {
			if len(path) == 1 {
				data, _, err := srcTrie.TryGetNode(path[0])
				if err != nil {
					t.Fatalf("failed to retrieve node data for path %x: %v", path, err)
				}
				tmpData, _ := common.RlpPaddingFilter(data)
				results[len(hashQueue)+i] = statedb.SyncResult{Hash: crypto.Keccak256Hash(tmpData).ToRootExtHash(), Data: data}
			} else {
				serializer := account.NewAccountSerializer()
				if err := rlp.DecodeBytes(srcTrie.Get(path[0]), serializer); err != nil {
					t.Fatalf("failed to decode account on path %x: %v", path, err)
				}
				acc := serializer.GetAccount()
				pacc := account.GetProgramAccount(acc)
				if pacc == nil {
					t.Errorf("failed to get contract")
				}
				stTrie, err := statedb.NewTrie(pacc.GetStorageRoot(), srcState.TrieDB())
				if err != nil {
					t.Fatalf("failed to retriev storage trie for path %x: %v", path, err)
				}
				data, _, err := stTrie.TryGetNode(path[1])
				if err != nil {
					t.Fatalf("failed to retrieve node data for path %x: %v", path, err)
				}
				tmpData, _ := common.RlpPaddingFilter(data)
				results[len(hashQueue)+i] = statedb.SyncResult{Hash: crypto.Keccak256Hash(tmpData).ToRootExtHash(), Data: data}
			}
		}
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := dstDiskDb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()

		nodes, paths, codes = sched.Missing(count)
		if !bypath {
			hashQueue = append(append(hashQueue[:0], nodes...), codes...)
		} else {
			hashQueue = append(hashQueue[:0], codes...)
			pathQueue = append(pathQueue[:0], paths...)
		}
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
	srcState, srcRoot, srcAccounts := makeTestState(t)
	newState, newRoot, _ := makeTestState(t)
	// commit stateTrie to DB
	srcState.TrieDB().Commit(srcRoot, false, 0)
	newState.TrieDB().Commit(newRoot, false, 0)

	isCode := func(hash common.Hash) bool {
		for _, acc := range srcAccounts {
			if hash == crypto.Keccak256Hash(acc.code) {
				return true
			}
		}
		return false
	}

	srcStateDB, err := New(srcRoot, srcState, nil)
	assert.NoError(t, err)

	it := NewNodeIterator(srcStateDB)
	it.Next() // skip trie root node

	for it.Next() {
		if !common.EmptyHash(it.Hash.ToHash()) {
			hash := it.Hash
			var (
				data []byte
				code = isCode(hash.ToHash())
				err  error
			)
			srcDiskDB := srcState.TrieDB().DiskDB()
			newDiskDB := newState.TrieDB().DiskDB()
			// Delete trie nodes or codes
			if code {
				data = srcDiskDB.ReadCode(hash)
				srcState.DeleteCode(hash)
				newState.DeleteCode(hash)
			} else {
				data, _ = srcDiskDB.ReadCachedTrieNode(hash)
				srcDiskDB.GetMemDB().Delete(hash.ToHash().Bytes())
				newDiskDB.GetMemDB().Delete(hash.ToHash().Bytes())
			}
			// Check consistency : errIterator
			err = CheckStateConsistency(srcState, newState, srcRoot, 100, nil)
			if !errors.Is(err, errIterator) {
				t.Log("mismatched err", "err", err, "expErr", errIterator)
				t.FailNow()
			}

			// Recover nodes
			srcDiskDB.GetMemDB().Put(hash.ToHash().Bytes(), data)
			newDiskDB.GetMemDB().Put(hash.ToHash().Bytes(), data)
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
	srcState.TrieDB().Commit(srcRoot, false, 0)

	// Create a destination state and sync with the scheduler
	dstDiskDB := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDiskDB)
	sched := NewStateSync(srcRoot.ToHash(), dstDiskDB, statedb.NewSyncBloom(1, dstDiskDB.GetMemDB()), nil, nil)

	nodes, _, codes := sched.Missing(0)
	queue := append(append([]common.ExtHash{}, nodes...), codes...)

	for len(queue) > 0 {
		// Sync only half of the scheduled nodes
		results := make([]statedb.SyncResult, len(queue)/2+1)
		for i, hash := range queue[:len(results)] {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				data, err = srcState.ContractCode(hash)
			}
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = statedb.SyncResult{Hash: hash, Data: data}
		}
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := dstDiskDB.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		nodes, _, codes := sched.Missing(0)
		queue = append(append(queue[len(results):], nodes...), codes...)
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
	srcState.TrieDB().Commit(srcRoot, false, 0)

	// Create a destination state and sync with the scheduler
	dstDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDb)
	sched := NewStateSync(srcRoot.ToHash(), dstDb, statedb.NewSyncBloom(1, dstDb.GetMemDB()), nil, nil)

	queue := make(map[common.ExtHash]struct{})
	nodes, _, codes := sched.Missing(count)
	for _, hash := range append(nodes, codes...) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
		// Fetch all the queued nodes in a random order
		results := make([]statedb.SyncResult, 0, len(queue))
		for hash := range queue {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				data, err = srcState.ContractCode(hash)
			}
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results = append(results, statedb.SyncResult{Hash: hash, Data: data})
		}
		// Feed the retrieved results back and queue new tasks
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := dstDb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		queue = make(map[common.ExtHash]struct{})
		nodes, _, codes := sched.Missing(0)
		for _, hash := range append(nodes, codes...) {
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
	srcState.TrieDB().Commit(srcRoot, false, 0)

	// Create a destination state and sync with the scheduler
	dstDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDb)
	sched := NewStateSync(srcRoot.ToHash(), dstDb, statedb.NewSyncBloom(1, dstDb.GetMemDB()), nil, nil)

	queue := make(map[common.ExtHash]struct{})
	nodes, _, codes := sched.Missing(0)
	for _, hash := range append(nodes, codes...) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
		// Sync only half of the scheduled nodes, even those in random order
		results := make([]statedb.SyncResult, 0, len(queue)/2+1)
		for hash := range queue {
			delete(queue, hash)

			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				data, err = srcState.ContractCode(hash)
			}
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results = append(results, statedb.SyncResult{Hash: hash, Data: data})

			if len(results) >= cap(results) {
				break
			}
		}
		// Feed the retrieved results back and queue new tasks
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := dstDb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()
		for _, result := range results {
			delete(queue, result.Hash)
		}
		nodes, _, codes = sched.Missing(0)
		for _, hash := range append(nodes, codes...) {
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
	srcState.TrieDB().Commit(srcRoot, false, 0)

	// isCode reports whether the hash is contract code hash.
	isCode := func(hash common.Hash) bool {
		for _, acc := range srcAccounts {
			if hash == crypto.Keccak256Hash(acc.code) {
				return true
			}
		}
		return false
	}
	checkTrieConsistency(srcState.TrieDB().DiskDB().(database.DBManager), srcRoot)

	// Create a destination state and sync with the scheduler
	dstDb := database.NewMemoryDBManager()
	dstState := NewDatabase(dstDb)
	sched := NewStateSync(srcRoot.ToHash(), dstDb, statedb.NewSyncBloom(1, dstDb.GetMemDB()), nil, nil)

	var added []common.ExtHash

	nodes, _, codes := sched.Missing(1)
	queue := append(append([]common.ExtHash{}, nodes...), codes...)

	for len(queue) > 0 {
		// Fetch a batch of state nodes
		results := make([]statedb.SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcState.TrieDB().Node(hash)
			if err != nil {
				data, err = srcState.ContractCode(hash)
			}
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = statedb.SyncResult{Hash: hash, Data: data}
		}
		// Process each of the state nodes
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
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
		for _, hash := range added {
			if isCode(hash.ToHash()) {
				continue
			}
			// Can't use checkStateConsistency here because subtrie keys may have odd
			// length and crash in LeafKey.
			if err := checkTrieConsistency(dstDb, hash); err != nil {
				t.Fatalf("state inconsistent: %v", err)
			}
		}
		// Fetch the next batch to retrieve
		nodes, _, codes = sched.Missing(1)
		queue = append(append(queue[:0], nodes...), codes...)
	}
	// Sanity check that removing any node from the database is detected
	for _, node := range added[1:] {
		var (
			key  = node.Bytes()
			code = isCode(node.ToHash())
			val  []byte
		)
		if code {
			val = dstDb.ReadCode(node)
			dstState.DeleteCode(node)
		} else {
			val, _ = dstDb.ReadCachedTrieNode(node)
			dstDb.GetMemDB().Delete(node.ToHash().Bytes())
		}

		if err := checkStateConsistency(dstDb, added[0]); err == nil {
			t.Fatalf("trie inconsistency not caught, missing: %x", key)
		}

		err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
		assert.Error(t, err)

		err = CheckStateConsistencyParallel(srcState, dstState, srcRoot, nil)
		assert.Error(t, err)
		if code {
			dstDb.WriteCode(node, val)
		} else {
			// insert a trie node to memory database
			dstDb.GetMemDB().Put(node.ToHash().Bytes(), val)
		}
	}

	err := CheckStateConsistency(srcState, dstState, srcRoot, 100, nil)
	assert.NoError(t, err)

	err = CheckStateConsistencyParallel(srcState, dstState, srcRoot, nil)
	assert.NoError(t, err)
}
