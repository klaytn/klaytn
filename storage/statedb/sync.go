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
// This file is derived from trie/sync.go (2018/06/04).
// Modified and improved for the klaytn development.

package statedb

import (
	"errors"
	"fmt"
	"strconv"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/prque"
	"github.com/klaytn/klaytn/storage/database"
)

// ErrNotRequested is returned by the trie sync when it's requested to process a
// node it did not request.
var ErrNotRequested = errors.New("not requested")

// ErrAlreadyProcessed is returned by the trie sync when it's requested to process a
// node it already processed previously.
var ErrAlreadyProcessed = errors.New("already processed")

// request represents a scheduled or already in-flight state retrieval request.
type request struct {
	hash common.Hash // Hash of the node data content to retrieve
	data []byte      // Data content of the node, cached until all subtrees complete
	raw  bool        // Whether this is a raw entry (code) or a trie node

	parents []*request // Parent state nodes referencing this entry (notify all upon completion)
	depth   int        // Depth level within the trie the node is located to prioritise DFS
	deps    int        // Number of dependencies before allowed to commit this node

	callback LeafCallback // Callback to invoke if a leaf node it reached on this branch
}

// SyncResult is a simple list to return missing nodes along with their request
// hashes.
type SyncResult struct {
	Hash common.Hash // Hash of the originally unknown trie node
	Data []byte      // Data content of the retrieved node
	Err  error
}

// syncMemBatch is an in-memory buffer of successfully downloaded but not yet
// persisted data items.
type syncMemBatch struct {
	batch map[common.Hash][]byte // In-memory membatch of recently completed items
	order []common.Hash          // Order of completion to prevent out-of-order data loss
}

// newSyncMemBatch allocates a new memory-buffer for not-yet persisted trie nodes.
func newSyncMemBatch() *syncMemBatch {
	return &syncMemBatch{
		batch: make(map[common.Hash][]byte),
		order: make([]common.Hash, 0, 256),
	}
}

type StateTrieReadDB interface {
	ReadStateTrieNode(key []byte) ([]byte, error)
	HasStateTrieNode(key []byte) (bool, error)
}

// TrieSync is the main state trie synchronisation scheduler, which provides yet
// unknown trie hashes to retrieve, accepts node data associated with said hashes
// and reconstructs the trie step by step until all is done.
type TrieSync struct {
	database         StateTrieReadDB          // Persistent database to check for existing entries
	membatch         *syncMemBatch            // Memory buffer to avoid frequent database writes
	requests         map[common.Hash]*request // Pending requests pertaining to a key hash
	queue            *prque.Prque             // Priority queue with the pending requests
	retrievedByDepth map[int]int              // Retrieved trie node number counted by depth
	committedByDepth map[int]int              // Committed trie nodes number counted by depth
	bloom            *SyncBloom               // Bloom filter for fast node existence checks
	exist            *lru.Cache               // exist to check if the trie node is already written or not
}

// NewTrieSync creates a new trie data download scheduler.
// If both bloom and cache are set, only cache is used.
func NewTrieSync(root common.Hash, database StateTrieReadDB, callback LeafCallback, bloom *SyncBloom, lruCache *lru.Cache) *TrieSync {
	ts := &TrieSync{
		database:         database,
		membatch:         newSyncMemBatch(),
		requests:         make(map[common.Hash]*request),
		queue:            prque.New(),
		retrievedByDepth: make(map[int]int),
		committedByDepth: make(map[int]int),
		bloom:            bloom,
		exist:            lruCache,
	}
	ts.AddSubTrie(root, 0, common.Hash{}, callback)
	return ts
}

// AddSubTrie registers a new trie to the sync code, rooted at the designated parent.
func (s *TrieSync) AddSubTrie(root common.Hash, depth int, parent common.Hash, callback LeafCallback) {
	// Short circuit if the trie is empty or already known
	if root == emptyRoot {
		return
	}
	if _, ok := s.membatch.batch[root]; ok {
		return
	}
	if s.exist != nil {
		if _, ok := s.exist.Get(root); ok {
			// already written in migration, skip the node
			return
		}
	} else if s.bloom == nil || s.bloom.Contains(root[:]) {
		// Bloom filter says this might be a duplicate, double check
		if ok, _ := s.database.HasStateTrieNode(root[:]); ok {
			logger.Info("skip write node in migration by ReadStateTrieNode", "AddSubTrie", root.String())
			return
		}
		// False positive, bump fault meter
		bloomFaultMeter.Mark(1)
	}
	// Assemble the new sub-trie sync request
	req := &request{
		hash:     root,
		depth:    depth,
		callback: callback,
	}
	// If this sub-trie has a designated parent, link them together
	if parent != (common.Hash{}) {
		ancestor := s.requests[parent]
		if ancestor == nil {
			panic(fmt.Sprintf("sub-trie ancestor not found: %x", parent))
		}
		ancestor.deps++
		req.parents = append(req.parents, ancestor)
	}
	s.schedule(req)
}

// AddRawEntry schedules the direct retrieval of a state entry that should not be
// interpreted as a trie node, but rather accepted and stored into the database
// as is. This method's goal is to support misc state metadata retrievals (e.g.
// contract code).
func (s *TrieSync) AddRawEntry(hash common.Hash, depth int, parent common.Hash) {
	// Short circuit if the entry is empty or already known
	if hash == emptyState {
		return
	}
	if _, ok := s.membatch.batch[hash]; ok {
		return
	}
	if s.exist != nil {
		if _, ok := s.exist.Get(hash); ok {
			// already written in migration, skip the node
			return
		}
	} else if s.bloom == nil || s.bloom.Contains(hash[:]) {
		// Bloom filter says this might be a duplicate, double check
		if ok, _ := s.database.HasStateTrieNode(hash.Bytes()); ok {
			return
		}
		// False positive, bump fault meter
		bloomFaultMeter.Mark(1)
	}
	// Assemble the new sub-trie sync request
	req := &request{
		hash:  hash,
		raw:   true,
		depth: depth,
	}
	// If this sub-trie has a designated parent, link them together
	if parent != (common.Hash{}) {
		ancestor := s.requests[parent]
		if ancestor == nil {
			panic(fmt.Sprintf("raw-entry ancestor not found: %x", parent))
		}
		ancestor.deps++
		req.parents = append(req.parents, ancestor)
	}
	s.schedule(req)
}

// Missing retrieves the known missing nodes from the trie for retrieval.
func (s *TrieSync) Missing(max int) []common.Hash {
	requests := []common.Hash{}
	for !s.queue.Empty() && (max == 0 || len(requests) < max) {
		requests = append(requests, s.queue.PopItem().(common.Hash))
	}
	return requests
}

// Process injects a batch of retrieved trie nodes data, returning if something
// was committed to the database and also the index of an entry if processing of
// it failed.
func (s *TrieSync) Process(results []SyncResult) (bool, int, error) {
	committed := false

	for i, item := range results {
		// If the item was not requested, bail out
		request := s.requests[item.Hash]
		if request == nil {
			return committed, i, ErrNotRequested
		}
		if request.data != nil {
			return committed, i, ErrAlreadyProcessed
		}
		// If the item is a raw entry request, commit directly
		if request.raw {
			request.data = item.Data
			s.commit(request)
			committed = true
			continue
		}
		// Decode the node data content and update the request
		node, err := decodeNode(item.Hash[:], item.Data)
		if err != nil {
			return committed, i, err
		}
		request.data = item.Data

		// Create and schedule a request for all the children nodes
		requests, err := s.children(request, node)
		if err != nil {
			return committed, i, err
		}
		if len(requests) == 0 && request.deps == 0 {
			s.commit(request)
			committed = true
			continue
		}
		request.deps += len(requests)
		for _, child := range requests {
			s.schedule(child)
		}
	}
	return committed, 0, nil
}

// Commit flushes the data stored in the internal membatch out to persistent
// storage, returning the number of items written and any occurred error.
func (s *TrieSync) Commit(dbw database.KeyValueWriter) (int, error) {
	// Dump the membatch into a database dbw
	for i, key := range s.membatch.order {
		if err := dbw.Put(key[:], s.membatch.batch[key]); err != nil {
			return i, err
		}

		if s.bloom != nil {
			s.bloom.Add(key[:])
		}

		if s.exist != nil {
			s.exist.Add(key, nil)
		}
	}
	written := len(s.membatch.order)

	// Drop the membatch data and return
	s.membatch = newSyncMemBatch()
	return written, nil
}

// Pending returns the number of state entries currently pending for download.
func (s *TrieSync) Pending() int {
	return len(s.requests)
}

// schedule inserts a new state retrieval request into the fetch queue. If there
// is already a pending request for this node, the new request will be discarded
// and only a parent reference added to the old one.
func (s *TrieSync) schedule(req *request) {
	// If we're already requesting this node, add a new reference and stop
	if old, ok := s.requests[req.hash]; ok {
		old.parents = append(old.parents, req.parents...)
		return
	}

	// Count the retrieved trie by depth
	s.retrievedByDepth[req.depth]++

	// Schedule the request for future retrieval
	s.queue.Push(req.hash, int64(req.depth))
	s.requests[req.hash] = req
}

// children retrieves all the missing children of a state trie entry for future
// retrieval scheduling.
func (s *TrieSync) children(req *request, object node) ([]*request, error) {
	// Gather all the children of the node, irrelevant whether known or not
	type child struct {
		node  node
		depth int
	}
	children := []child{}

	switch node := (object).(type) {
	case *shortNode:
		children = []child{{
			node:  node.Val,
			depth: req.depth + len(node.Key),
		}}
	case *fullNode:
		for i := 0; i < 17; i++ {
			if node.Children[i] != nil {
				children = append(children, child{
					node:  node.Children[i],
					depth: req.depth + 1,
				})
			}
		}
	default:
		panic(fmt.Sprintf("unknown node: %+v", node))
	}
	// Iterate over the children, and request all unknown ones
	requests := make([]*request, 0, len(children))
	for _, child := range children {
		// Notify any external watcher of a new key/value node
		if req.callback != nil {
			if node, ok := (child.node).(valueNode); ok {
				if err := req.callback(node, req.hash, child.depth); err != nil {
					return nil, err
				}
			}
		}
		// If the child references another node, resolve or schedule
		if node, ok := (child.node).(hashNode); ok {
			// Try to resolve the node from the local database
			hash := common.BytesToHash(node)
			if _, ok := s.membatch.batch[hash]; ok {
				continue
			}
			if s.exist != nil {
				if _, ok := s.exist.Get(hash); ok {
					// already written in migration, skip the node
					continue
				}
			} else if s.bloom == nil || s.bloom.Contains(node) {
				// Bloom filter says this might be a duplicate, double check
				if ok, _ := s.database.HasStateTrieNode(node); ok {
					continue
				}
				// False positive, bump fault meter
				bloomFaultMeter.Mark(1)
			}

			// Locally unknown node, schedule for retrieval
			requests = append(requests, &request{
				hash:     hash,
				parents:  []*request{req},
				depth:    child.depth,
				callback: req.callback,
			})
		}
	}
	return requests, nil
}

// commit finalizes a retrieval request and stores it into the membatch. If any
// of the referencing parent requests complete due to this commit, they are also
// committed themselves.
func (s *TrieSync) commit(req *request) (err error) {
	// Count the committed trie by depth and Clear the counts of lower depth
	s.committedByDepth[req.depth]++

	// Write the node content to the membatch
	s.membatch.batch[req.hash] = req.data
	s.membatch.order = append(s.membatch.order, req.hash)

	delete(s.requests, req.hash)

	// Check all parents for completion
	for _, parent := range req.parents {
		parent.deps--
		if parent.deps == 0 {
			if err := s.commit(parent); err != nil {
				return err
			}
		}
	}
	return nil
}

// RetrievedByDepth returns the retrieved trie count by given depth.
// This number is same as the number of nodes that needs to be committed to complete trie sync.
func (s *TrieSync) RetrievedByDepth(depth int) int {
	return s.retrievedByDepth[depth]
}

// CommittedByDepth returns the committed trie count by given depth.
func (s *TrieSync) CommittedByDepth(depth int) int {
	return s.committedByDepth[depth]
}

// CalcProgressPercentage returns the progress percentage.
func (s *TrieSync) CalcProgressPercentage() float64 {
	var progress float64
	// depth	max trie	resolution (%)
	// 0	 	1 	 		100.00000
	// 1	 	16 	 		6.25000
	// 2	 	256 	 	0.39063
	// 3	 	4,096 	 	0.02441
	// 4	 	65,536 	 	0.00153
	// 5	 	1,048,576 	0.00010

	for i := 0; i < 20; i++ {
		c, r := s.CommittedByDepth(i), s.RetrievedByDepth(i)

		var progressByDepth float64

		if r == 0 {
			break
		}

		if r > 0 {
			progressByDepth = float64(c) / float64(r) * 100
			if progressByDepth > progress && i < 4 { // Scan depth 0 ~ 3 for accuracy
				progress = progressByDepth
			}
		}

		logger.Debug("Trie sync progress by depth #"+strconv.Itoa(i), "committed", c, "retrieved", r, "progress", progressByDepth)
	}

	logger.Debug("Trie sync progress ", "progress", strconv.FormatFloat(progress, 'f', -1, 64)+"%")

	return progress
}
