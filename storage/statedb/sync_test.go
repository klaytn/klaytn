// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from trie/sync_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package statedb

import (
	"bytes"
	"testing"

	"github.com/alecthomas/units"
	lru "github.com/hashicorp/golang-lru"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/storage/database"
)

// makeTestTrie create a sample test trie to test node-wise reconstruction.
func makeTestTrie() (*Database, *SecureTrie, map[string][]byte) {
	// Create an empty trie
	triedb := NewDatabase(database.NewMemoryDBManager())
	trie, _ := NewSecureTrie(common.Hash{}, triedb, nil)

	// Fill it with some arbitrary data
	content := make(map[string][]byte)
	for i := byte(0); i < 255; i++ {
		// Map the same data under multiple keys
		key, val := common.LeftPadBytes([]byte{1, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

		key, val = common.LeftPadBytes([]byte{2, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

		// Add some other data to inflate the trie
		for j := byte(3); j < 13; j++ {
			key, val = common.LeftPadBytes([]byte{j, i}, 32), []byte{j, i}
			content[string(key)] = val
			trie.Update(key, val)
		}
	}
	trie.Commit(nil)

	// Return the generated trie
	return triedb, trie, content
}

// checkTrieContents cross references a reconstructed trie with an expected data
// content map.
func checkTrieContents(t *testing.T, db *Database, root []byte, content map[string][]byte) {
	// Check root availability and trie contents
	trie, err := NewSecureTrie(common.BytesToHash(root), db, nil)
	if err != nil {
		t.Fatalf("failed to create trie at %x: %v", root, err)
	}
	if err := checkTrieConsistency(db, common.BytesToHash(root)); err != nil {
		t.Fatalf("inconsistent trie at %x: %v", root, err)
	}
	for key, val := range content {
		if have := trie.Get([]byte(key)); !bytes.Equal(have, val) {
			t.Errorf("entry %x: content mismatch: have %x, want %x", key, have, val)
		}
	}
}

// checkTrieConsistency checks that all nodes in a trie are indeed present.
func checkTrieConsistency(db *Database, root common.Hash) error {
	// Create and iterate a trie rooted in a subnode
	trie, err := NewSecureTrie(root, db, nil)
	if err != nil {
		return nil // Consider a non existent state consistent
	}
	it := trie.NodeIterator(nil)
	for it.Next(true) {
	}
	return it.Error()
}

// Tests that an empty trie is not scheduled for syncing.
func TestEmptyTrieSync(t *testing.T) {
	memDBManagerA := database.NewMemoryDBManager()
	memDBManagerB := database.NewMemoryDBManager()
	dbA := NewDatabase(memDBManagerA)
	dbB := NewDatabase(memDBManagerB)
	emptyA, _ := NewTrie(common.Hash{}, dbA, nil)
	emptyB, _ := NewTrie(emptyRoot, dbB, nil)

	for i, trie := range []*Trie{emptyA, emptyB} {
		sync := NewTrieSync(trie.Hash(), database.NewMemoryDBManager(), nil, NewSyncBloom(1, database.NewMemDB()), nil)
		if nodes, paths, codes := sync.Missing(1); len(nodes) != 0 || len(paths) != 0 || len(codes) != 0 {
			t.Errorf("test %d: content requested for empty trie: %v, %v, %v", i, nodes, paths, codes)
		}
	}
}

// Tests that given a root hash, a trie can sync iteratively on a single thread,
// requesting retrieval tasks and returning all of them in one go.
func TestIterativeSyncIndividual(t *testing.T)       { testIterativeTrieSync(t, 1, false) }
func TestIterativeSyncBatched(t *testing.T)          { testIterativeTrieSync(t, 100, false) }
func TestIterativeSyncIndividualByPath(t *testing.T) { testIterativeTrieSync(t, 1, true) }
func TestIterativeSyncBatchedByPath(t *testing.T)    { testIterativeTrieSync(t, 100, true) }

func trieSyncLoop(t *testing.T, count int, srcTrie *SecureTrie, sched *TrieSync, srcDB *Database, diskDB database.Database, bypath bool) {
	nodes, paths, codes := sched.Missing(count)
	var (
		hashQueue []common.Hash
		pathQueue []SyncPath
	)
	if !bypath {
		hashQueue = append(append(hashQueue[:0], nodes...), codes...)
	} else {
		hashQueue = append(hashQueue[:0], codes...)
		pathQueue = append(pathQueue[:0], paths...)
	}
	for len(hashQueue)+len(pathQueue) > 0 {
		results := make([]SyncResult, len(hashQueue)+len(pathQueue))
		for i, hash := range hashQueue {
			data, err := srcDB.Node(hash.ExtendZero())
			if err != nil {
				t.Fatalf("failed to retrieve node data for hash %x: %v", hash, err)
			}
			results[i] = SyncResult{hash, data, nil}
		}
		for i, path := range pathQueue {
			data, _, err := srcTrie.TryGetNode(path[0])
			if err != nil {
				t.Fatalf("failed to retrieve node data for path %x: %v", path, err)
			}
			results[len(hashQueue)+i] = SyncResult{crypto.Keccak256Hash(data), data, nil}
		}
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := diskDB.NewBatch()
		if index, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
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
}

func testIterativeTrieSync(t *testing.T, count int, bypath bool) {
	// Create a random trie to copy
	srcDb, srcTrie, srcData := makeTestTrie()

	// test with bloom filter
	{
		memDBManager := database.NewMemoryDBManager()
		diskdb := memDBManager.GetMemDB()
		triedb := NewDatabase(memDBManager)
		sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, NewSyncBloom(1, diskdb), nil)

		trieSyncLoop(t, count, srcTrie, sched, srcDb, diskdb, bypath)
		// Cross check that the two tries are in sync
		checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
	}

	// test with lru cache
	{
		memDBManager := database.NewMemoryDBManager()
		diskdb := memDBManager.GetMemDB()
		triedb := NewDatabase(memDBManager)
		lruCache, _ := lru.New(int(1 * units.MB / common.HashLength))
		sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, nil, lruCache)

		trieSyncLoop(t, count, srcTrie, sched, srcDb, diskdb, bypath)
		// Cross check that the two tries are in sync
		checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
	}

	// test without bloom, lru cache
	{
		memDBManager := database.NewMemoryDBManager()
		diskdb := memDBManager.GetMemDB()
		triedb := NewDatabase(memDBManager)
		sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, nil, nil)

		trieSyncLoop(t, count, srcTrie, sched, srcDb, diskdb, bypath)
		// Cross check that the two tries are in sync
		checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
	}

	// test with bloom, lru cache
	{
		memDBManager := database.NewMemoryDBManager()
		diskdb := memDBManager.GetMemDB()
		triedb := NewDatabase(memDBManager)
		bloom := NewSyncBloom(1, diskdb)
		lruCache, _ := lru.New(int(1 * units.MB / common.HashLength))
		sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, bloom, lruCache)

		trieSyncLoop(t, count, srcTrie, sched, srcDb, diskdb, bypath)
		// Cross check that the two tries are in sync
		checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
	}
}

// Tests that the trie scheduler can correctly reconstruct the state even if only
// partial results are returned, and the others sent only later.
func TestIterativeDelayedTrieSync(t *testing.T) {
	// Create a random trie to copy
	srcDb, srcTrie, srcData := makeTestTrie()

	// Create a destination trie and sync with the scheduler
	memDBManager := database.NewMemoryDBManager()
	diskdb := memDBManager.GetMemDB()
	triedb := NewDatabase(memDBManager)
	sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, NewSyncBloom(1, diskdb), nil)

	nodes, _, codes := sched.Missing(10000)
	queue := append(append([]common.Hash{}, nodes...), codes...)
	for len(queue) > 0 {
		// Sync only half of the scheduled nodes
		results := make([]SyncResult, len(queue)/2+1)
		for i, hash := range queue[:len(results)] {
			data, err := srcDb.Node(hash.ExtendZero())
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x: %v", hash, err)
			}
			results[i] = SyncResult{hash, data, nil}
		}
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := diskdb.NewBatch()
		if index, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		batch.Write()
		nodes, _, codes = sched.Missing(10000)
		queue = append(append(queue[len(results):], nodes...), codes...)
	}
	// Cross check that the two tries are in sync
	checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
}

// Tests that given a root hash, a trie can sync iteratively on a single thread,
// requesting retrieval tasks and returning all of them in one go, however in a
// random order.
func TestIterativeRandomTrieSyncIndividual(t *testing.T) { testIterativeRandomTrieSync(t, 1) }
func TestIterativeRandomTrieSyncBatched(t *testing.T)    { testIterativeRandomTrieSync(t, 100) }

func testIterativeRandomTrieSync(t *testing.T, count int) {
	// Create a random trie to copy
	srcDb, srcTrie, srcData := makeTestTrie()

	// Create a destination trie and sync with the scheduler
	memDBManager := database.NewMemoryDBManager()
	diskdb := memDBManager.GetMemDB()
	triedb := NewDatabase(memDBManager)
	sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, NewSyncBloom(1, diskdb), nil)

	queue := make(map[common.Hash]struct{})
	nodes, _, codes := sched.Missing(count)
	for _, hash := range append(nodes, codes...) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
		// Fetch all the queued nodes in a random order
		results := make([]SyncResult, 0, len(queue))
		for hash := range queue {
			data, err := srcDb.Node(hash.ExtendZero())
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x: %v", hash, err)
			}
			results = append(results, SyncResult{hash, data, nil})
		}
		// Feed the retrieved results back and queue new tasks
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := diskdb.NewBatch()
		if index, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		batch.Write()
		queue = make(map[common.Hash]struct{})
		nodes, _, codes = sched.Missing(count)
		for _, hash := range append(nodes, codes...) {
			queue[hash] = struct{}{}
		}
	}
	// Cross check that the two tries are in sync
	checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
}

// Tests that the trie scheduler can correctly reconstruct the state even if only
// partial results are returned (Even those randomly), others sent only later.
func TestIterativeRandomDelayedTrieSync(t *testing.T) {
	// Create a random trie to copy
	srcDb, srcTrie, srcData := makeTestTrie()

	// Create a destination trie and sync with the scheduler
	memDBManager := database.NewMemoryDBManager()
	diskdb := memDBManager.GetMemDB()
	triedb := NewDatabase(memDBManager)
	sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, NewSyncBloom(1, diskdb), nil)

	queue := make(map[common.Hash]struct{})
	nodes, _, codes := sched.Missing(10000)
	for _, hash := range append(nodes, codes...) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
		// Sync only half of the scheduled nodes, even those in random order
		results := make([]SyncResult, 0, len(queue)/2+1)
		for hash := range queue {
			data, err := srcDb.Node(hash.ExtendZero())
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x: %v", hash, err)
			}
			results = append(results, SyncResult{hash, data, nil})

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
		batch := diskdb.NewBatch()
		if index, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		batch.Write()
		for _, result := range results {
			delete(queue, result.Hash)
		}
		nodes, _, codes := sched.Missing(10000)
		for _, hash := range append(nodes, codes...) {
			queue[hash] = struct{}{}
		}
	}
	// Cross check that the two tries are in sync
	checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
}

// Tests that a trie sync will not request nodes multiple times, even if they
// have such references.
func TestDuplicateAvoidanceTrieSync(t *testing.T) {
	// Create a random trie to copy
	srcDb, srcTrie, srcData := makeTestTrie()

	// Create a destination trie and sync with the scheduler
	memDBManager := database.NewMemoryDBManager()
	diskdb := memDBManager.GetMemDB()
	triedb := NewDatabase(memDBManager)
	sched := NewTrieSync(srcTrie.Hash(), memDBManager, nil, NewSyncBloom(1, diskdb), nil)

	nodes, _, codes := sched.Missing(0)
	queue := append(append([]common.Hash{}, nodes...), codes...)
	requested := make(map[common.Hash]struct{})

	for len(queue) > 0 {
		results := make([]SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcDb.Node(hash.ExtendZero())
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x: %v", hash, err)
			}
			if _, ok := requested[hash]; ok {
				t.Errorf("hash %x already requested once", hash)
			}
			requested[hash] = struct{}{}

			results[i] = SyncResult{hash, data, nil}
		}
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := diskdb.NewBatch()
		if index, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		batch.Write()
		nodes, _, codes = sched.Missing(0)
		queue = append(append(queue[:0], nodes...), codes...)
	}
	// Cross check that the two tries are in sync
	checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)
}

// Tests that at any point in time during a sync, only complete sub-tries are in
// the database.
func TestIncompleteTrieSync(t *testing.T) {
	// Create a random trie to copy
	srcDb, srcTrie, _ := makeTestTrie()

	// Create a destination trie and sync with the scheduler
	dbm := database.NewMemoryDBManager()
	triedb := NewDatabase(dbm)
	sched := NewTrieSync(srcTrie.Hash(), dbm, nil, NewSyncBloom(1, dbm.GetMemDB()), nil)

	var added []common.Hash
	nodes, _, codes := sched.Missing(1)
	queue := append(append([]common.Hash{}, nodes...), codes...)
	for len(queue) > 0 {
		// Fetch a batch of trie nodes
		results := make([]SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcDb.Node(hash.ExtendZero())
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x: %v", hash, err)
			}
			results[i] = SyncResult{hash, data, nil}
		}
		// Process each of the trie nodes
		for index, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result #%d: %v", index, err)
			}
		}
		batch := dbm.NewBatch(database.StateTrieDB)
		if index, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		batch.Write()
		for _, result := range results {
			added = append(added, result.Hash)
		}
		// Check that all known sub-tries in the synced trie are complete
		for _, root := range added {
			if err := checkTrieConsistency(triedb, root); err != nil {
				t.Fatalf("trie inconsistent: %v", err)
			}
		}
		// Fetch the next batch to retrieve
		nodes, _, codes = sched.Missing(1)
		queue = append(append(queue[:0], nodes...), codes...)
	}
	// Sanity check that removing any node from the database is detected
	for _, hash := range added[1:] {
		nodehash := hash.ExtendZero()
		value, _ := dbm.ReadTrieNode(nodehash)

		dbm.DeleteTrieNode(nodehash)
		if err := checkTrieConsistency(triedb, added[0]); err == nil {
			t.Fatalf("trie inconsistency not caught, missing: %x", hash)
		}
		dbm.WriteTrieNode(nodehash, value)
	}
}

// Tests that trie nodes get scheduled lexicographically when having the same
// depth.
func TestSyncOrdering(t *testing.T) {
	// Create a random trie to copy
	srcDb, srcTrie, srcData := makeTestTrie()

	// Create a destination trie and sync with the scheduler, tracking the requests
	diskdb := database.NewMemoryDBManager()
	triedb := NewDatabase(diskdb)
	sched := NewTrieSync(srcTrie.Hash(), diskdb, nil, NewSyncBloom(1, diskdb.GetMemDB()), nil)

	nodes, paths, _ := sched.Missing(1)
	queue := append([]common.Hash{}, nodes...)
	reqs := append([]SyncPath{}, paths...)

	for len(queue) > 0 {
		results := make([]SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcDb.Node(hash.ExtendZero())
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x: %v", hash, err)
			}
			results[i] = SyncResult{hash, data, nil}
		}
		for _, result := range results {
			if err := sched.Process(result); err != nil {
				t.Fatalf("failed to process result %v", err)
			}
		}
		batch := diskdb.NewBatch(database.StateTrieDB)
		if _, err := sched.Commit(batch); err != nil {
			t.Fatalf("failed to commit data: %v", err)
		}
		batch.Write()

		nodes, paths, _ = sched.Missing(1)
		queue = append(queue[:0], nodes...)
		reqs = append(reqs, paths...)
	}
	// Cross check that the two tries are in sync
	checkTrieContents(t, triedb, srcTrie.Hash().Bytes(), srcData)

	// Check that the trie nodes have been requested path-ordered
	for i := 0; i < len(reqs)-1; i++ {
		if len(reqs[i]) > 1 || len(reqs[i+1]) > 1 {
			// In the case of the trie tests, there's no storage so the tuples
			// must always be single items. 2-tuples should be tested in state.
			t.Errorf("Invalid request tuples: len(%v) or len(%v) > 1", reqs[i], reqs[i+1])
		}
		if bytes.Compare(compactToHex(reqs[i][0]), compactToHex(reqs[i+1][0])) > 0 {
			t.Errorf("Invalid request order: %v before %v", compactToHex(reqs[i][0]), compactToHex(reqs[i+1][0]))
		}
	}
}
