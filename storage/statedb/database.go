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
// This file is derived from trie/database.go (2018/06/04).
// Modified and improved for the klaytn development.

package statedb

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/pbnjay/memory"
	"github.com/rcrowley/go-metrics"
)

var (
	logger = log.NewModuleLogger(log.StorageStateDB)

	// metrics for Cap state
	memcacheFlushTimeGauge  = metrics.NewRegisteredGauge("trie/memcache/flush/time", nil)
	memcacheFlushNodesGauge = metrics.NewRegisteredGauge("trie/memcache/flush/nodes", nil)
	memcacheFlushSizeGauge  = metrics.NewRegisteredGauge("trie/memcache/flush/size", nil)

	// metrics for GC
	memcacheGCTimeGauge  = metrics.NewRegisteredGauge("trie/memcache/gc/time", nil)
	memcacheGCNodesMeter = metrics.NewRegisteredMeter("trie/memcache/gc/nodes", nil)
	memcacheGCSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/gc/size", nil)

	// metrics for commit state
	memcacheCommitTimeGauge  = metrics.NewRegisteredGauge("trie/memcache/commit/time", nil)
	memcacheCommitNodesMeter = metrics.NewRegisteredMeter("trie/memcache/commit/nodes", nil)
	memcacheCommitSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/commit/size", nil)
	memcacheUncacheTimeGauge = metrics.NewRegisteredGauge("trie/memcache/uncache/time", nil)

	// metrics for state trie cache db
	memcacheCleanHitMeter          = metrics.NewRegisteredMeter("trie/memcache/clean/hit", nil)
	memcacheCleanMissMeter         = metrics.NewRegisteredMeter("trie/memcache/clean/miss", nil)
	memcacheCleanPrefetchMissMeter = metrics.NewRegisteredMeter("trie/memcache/clean/prefetch/miss", nil)
	memcacheCleanReadMeter         = metrics.NewRegisteredMeter("trie/memcache/clean/read", nil)
	memcacheCleanWriteMeter        = metrics.NewRegisteredMeter("trie/memcache/clean/write", nil)

	// metric of total node number
	memcacheNodesGauge = metrics.NewRegisteredGauge("trie/memcache/nodes", nil)
)

// commitResultChSizeLimit limits the size of channel used for commitResult.
const commitResultChSizeLimit = 100 * 10000

// AutoScaling is for auto-scaling cache size. If cacheSize is set to this value,
// cache size is set scaling to physical memeory
const AutoScaling = -1

type DatabaseReader interface {
	// Get retrieves the value associated with key from the database.
	Get(key []byte) (value []byte, err error)

	// Has retrieves whether a key is present in the database.
	Has(key []byte) (bool, error)
}

// Database is an intermediate write layer between the trie data structures and
// the disk database. The aim is to accumulate trie writes in-memory and only
// periodically flush a couple tries to disk, garbage collecting the remainder.
type Database struct {
	diskDB database.DBManager // Persistent storage for matured trie nodes

	nodes  map[common.ExtHash]*cachedNode // Data and references relationships of a trie node
	oldest common.ExtHash                 // Oldest tracked node, flush-list head
	newest common.ExtHash                 // Newest tracked node, flush-list tail

	preimages    map[common.Hash][]byte // Preimages of nodes from the secure trie
	pruningMarks []database.PruningMark // Trie node pruning marks from the pruning trie

	gctime  time.Duration      // Time spent on garbage collection since last commit
	gcnodes uint64             // Nodes garbage collected since last commit
	gcsize  common.StorageSize // Data storage garbage collected since last commit
	gcLock  sync.RWMutex       // Lock for preventing to garbage collect cachedNode without flushing

	flushtime  time.Duration      // Time spent on data flushing since last commit
	flushnodes uint64             // Nodes flushed since last commit
	flushsize  common.StorageSize // Data storage flushed since last commit

	nodesSize     common.StorageSize // Storage size of the nodes cache
	preimagesSize common.StorageSize // Storage size of the preimages cache

	lock sync.RWMutex

	trieNodeCache                TrieNodeCache        // GC friendly memory cache of trie node RLPs
	trieNodeCacheConfig          *TrieNodeCacheConfig // Configuration of trieNodeCache
	savingTrieNodeCacheTriggered bool                 // Whether saving trie node cache has been triggered or not
}

// rawNode is a simple binary blob used to differentiate between collapsed trie
// nodes and already encoded RLP binary blobs (while at the same time store them
// in the same cache fields).
type rawNode []byte

func (n rawNode) canUnload(uint16, uint16) bool { panic("this should never end up in a live trie") }
func (n rawNode) cache() (hashNode, bool)       { panic("this should never end up in a live trie") }
func (n rawNode) fstring(ind string) string     { panic("this should never end up in a live trie") }
func (n rawNode) lenEncoded() uint16            { panic("this should never end up in a live trie") }
func (n rawNode) EncodeRLP(w io.Writer) error {
	_, err := w.Write([]byte(n))
	return err
}

// rawFullNode represents only the useful data content of a full node, with the
// caches and flags stripped out to minimize its data database. This type honors
// the same RLP encoding as the original parent.
type rawFullNode [17]node

func (n rawFullNode) canUnload(uint16, uint16) bool { panic("this should never end up in a live trie") }
func (n rawFullNode) cache() (hashNode, bool)       { panic("this should never end up in a live trie") }
func (n rawFullNode) fstring(ind string) string     { panic("this should never end up in a live trie") }
func (n rawFullNode) lenEncoded() uint16            { panic("this should never end up in a live trie") }

func (n rawFullNode) EncodeRLP(w io.Writer) error {
	var nodes [17]node

	for i, child := range n {
		if child != nil {
			nodes[i] = child
		} else {
			nodes[i] = nilValueNode
		}
	}
	return rlp.Encode(w, nodes)
}

// rawShortNode represents only the useful data content of a short node, with the
// caches and flags stripped out to minimize its data database. This type honors
// the same RLP encoding as the original parent.
type rawShortNode struct {
	Key []byte
	Val node
}

func (n rawShortNode) canUnload(uint16, uint16) bool {
	panic("this should never end up in a live trie")
}
func (n rawShortNode) cache() (hashNode, bool)   { panic("this should never end up in a live trie") }
func (n rawShortNode) fstring(ind string) string { panic("this should never end up in a live trie") }
func (n rawShortNode) lenEncoded() uint16        { panic("this should never end up in a live trie") }

// cachedNode is all the information we know about a single cached trie node
// in the memory database write layer.
type cachedNode struct {
	node node // Cached collapsed trie node, or raw rlp data
	// TODO-Klaytn: need to change data type of this if we increase the code size limit
	size uint16 // Byte size of the useful cached data

	parents  uint64                    // Number of live nodes referencing this one
	children map[common.ExtHash]uint64 // External children referenced by this node

	flushPrev common.ExtHash // Previous node in the flush-list
	flushNext common.ExtHash // Next node in the flush-list
}

// rlp returns the raw rlp encoded blob of the cached trie node, either directly
// from the cache, or by regenerating it from the collapsed node.
func (n *cachedNode) rlp() []byte {
	if node, ok := n.node.(rawNode); ok {
		return node
	}
	return nodeToBytes(n.node)
}

// obj returns the decoded and expanded trie node, either directly from the cache,
// or by regenerating it from the rlp encoded blob.
func (n *cachedNode) obj(hash common.ExtHash) node {
	if node, ok := n.node.(rawNode); ok {
		return mustDecodeNode(hash[:], node)
	}
	return expandNode(hash[:], n.node)
}

// childs returns all the tracked children of this node, both the implicit ones
// from inside the node as well as the explicit ones from outside the node.
func (n *cachedNode) childs() []common.ExtHash {
	children := make([]common.ExtHash, 0, 16)
	for child := range n.children {
		children = append(children, child)
	}
	if _, ok := n.node.(rawNode); !ok {
		gatherChildren(n.node, &children)
	}
	return children
}

// gatherChildren traverses the node hierarchy of a collapsed database node and
// retrieves all the hashnode children.
func gatherChildren(n node, children *[]common.ExtHash) {
	switch n := n.(type) {
	case *rawShortNode:
		gatherChildren(n.Val, children)

	case rawFullNode:
		for i := 0; i < 16; i++ {
			gatherChildren(n[i], children)
		}
	case hashNode:
		*children = append(*children, common.BytesToExtHash(n))

	case valueNode, nil, rawNode:

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

// simplifyNode traverses the hierarchy of an expanded memory node and discards
// all the internal caches, returning a node that only contains the raw data.
func simplifyNode(n node) node {
	switch n := n.(type) {
	case *shortNode:
		// Short nodes discard the flags and cascade
		return &rawShortNode{Key: n.Key, Val: simplifyNode(n.Val)}

	case *fullNode:
		// Full nodes discard the flags and cascade
		node := rawFullNode(n.Children)
		for i := 0; i < len(node); i++ {
			if node[i] != nil {
				node[i] = simplifyNode(node[i])
			}
		}
		return node

	case valueNode, hashNode, rawNode:
		return n

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

// expandNode traverses the node hierarchy of a collapsed database node and converts
// all fields and keys into expanded memory form.
func expandNode(hash hashNode, n node) node {
	switch n := n.(type) {
	case *rawShortNode:
		// Short nodes need key and child expansion
		return &shortNode{
			Key: compactToHex(n.Key),
			Val: expandNode(nil, n.Val),
			flags: nodeFlag{
				hash: hash,
			},
		}

	case rawFullNode:
		// Full nodes need child expansion
		node := &fullNode{
			flags: nodeFlag{
				hash: hash,
			},
		}
		for i := 0; i < len(node.Children); i++ {
			if n[i] != nil {
				node.Children[i] = expandNode(nil, n[i])
			}
		}
		return node

	case valueNode, hashNode:
		return n

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

// NewDatabase creates a new trie database to store ephemeral trie content before
// its written out to disk or garbage collected.
func NewDatabase(diskDB database.DBManager) *Database {
	return NewDatabaseWithNewCache(diskDB, GetEmptyTrieNodeCacheConfig())
}

// NewDatabaseWithNewCache creates a new trie database to store ephemeral trie content
// before its written out to disk or garbage collected. It also acts as a read cache
// for nodes loaded from disk.
func NewDatabaseWithNewCache(diskDB database.DBManager, cacheConfig *TrieNodeCacheConfig) *Database {
	trieNodeCache, err := NewTrieNodeCache(cacheConfig)
	if err != nil {
		logger.Error("Invalid trie node cache config", "err", err, "config", cacheConfig)
	}

	return &Database{
		diskDB:              diskDB,
		nodes:               map[common.ExtHash]*cachedNode{{}: {}},
		preimages:           make(map[common.Hash][]byte),
		trieNodeCache:       trieNodeCache,
		trieNodeCacheConfig: cacheConfig,
	}
}

// NewDatabaseWithExistingCache creates a new trie database to store ephemeral trie content
// before its written out to disk or garbage collected. It also acts as a read cache
// for nodes loaded from disk.
func NewDatabaseWithExistingCache(diskDB database.DBManager, cache TrieNodeCache) *Database {
	return &Database{
		diskDB:        diskDB,
		nodes:         map[common.ExtHash]*cachedNode{{}: {}},
		preimages:     make(map[common.Hash][]byte),
		trieNodeCache: cache,
	}
}

func getTrieNodeCacheSizeMiB() int {
	totalPhysicalMemMiB := float64(memory.TotalMemory() / 1024 / 1024)

	if totalPhysicalMemMiB < 10*1024 {
		return 0
	} else if totalPhysicalMemMiB < 20*1024 {
		return 1 * 1024 // allocate 1G for small memory
	} else if totalPhysicalMemMiB < 30*1024 {
		return 6 * 1024 // allocate 6G for medium memory
	} else {
		return 10 * 1024 // allocate 10G for large memory (>20G)
	} 
}

// DiskDB retrieves the persistent database backing the trie database.
func (db *Database) DiskDB() database.DBManager {
	return db.diskDB
}

// TrieNodeCache retrieves the trieNodeCache of the trie database.
func (db *Database) TrieNodeCache() TrieNodeCache {
	return db.trieNodeCache
}

// GetTrieNodeCacheConfig returns the configuration of TrieNodeCache.
func (db *Database) GetTrieNodeCacheConfig() *TrieNodeCacheConfig {
	return db.trieNodeCacheConfig
}

// GetTrieNodeLocalCacheByteLimit returns the byte size of trie node cache.
func (db *Database) GetTrieNodeLocalCacheByteLimit() uint64 {
	return uint64(db.trieNodeCacheConfig.LocalCacheSizeMiB) * 1024 * 1024
}

// RLockGCCachedNode locks the GC lock of CachedNode.
func (db *Database) RLockGCCachedNode() {
	db.gcLock.RLock()
}

// RUnlockGCCachedNode unlocks the GC lock of CachedNode.
func (db *Database) RUnlockGCCachedNode() {
	db.gcLock.RUnlock()
}

// NodeChildren retrieves the children of the given hash trie
func (db *Database) NodeChildren(hash common.ExtHash) ([]common.ExtHash, error) {
	childrenHash := make([]common.ExtHash, 0, 16)

	if common.EmptyExtHash(hash) {
		return childrenHash, ErrZeroHashNode
	}

	n, _ := db.node(hash)
	if n == nil {
		return childrenHash, nil
	}

	children := make([]node, 0, 16)

	switch n := (n).(type) {
	case *shortNode:
		children = []node{n.Val}
	case *fullNode:
		for i := 0; i < 17; i++ {
			if n.Children[i] != nil {
				children = append(children, n.Children[i])
			}
		}
	}

	for _, child := range children {
		n, ok := child.(hashNode)
		if ok {
			hash := common.BytesToExtHash(n)
			childrenHash = append(childrenHash, hash)
		}
	}

	return childrenHash, nil
}

// insert inserts a collapsed trie node into the memory database.
// The blob size must be specified to allow proper size tracking.
// All nodes inserted by this function will be reference tracked
// and in theory should only used for **trie nodes** insertion.
func (db *Database) insert(hash common.ExtHash, lenEncoded uint16, node node) {
	// If the node's already cached, skip
	if _, ok := db.nodes[hash]; ok {
		return
	}
	// Create the cached entry for this node
	entry := &cachedNode{
		node:      simplifyNode(node),
		size:      lenEncoded,
		flushPrev: db.newest,
	}
	for _, child := range entry.childs() {
		if c := db.nodes[child]; c != nil {
			c.parents++
		}
	}
	db.nodes[hash] = entry

	// Update the flush-list endpoints
	if common.EmptyExtHash(db.oldest) {
		db.oldest, db.newest = hash, hash
	} else {
		if _, ok := db.nodes[db.newest]; !ok {
			missingNewest := db.newest
			db.newest = db.getLastNodeHashInFlushList()
			db.nodes[db.newest].flushNext = common.ExtHash{}
			logger.Error("Found a newest node for missingNewest", "oldNewest", missingNewest, "newNewest", db.newest)
		}
		db.nodes[db.newest].flushNext, db.newest = hash, hash
	}
	db.nodesSize += common.StorageSize(common.HashLength + entry.size)
}

// insertPreimage writes a new trie node pre-image to the memory database if it's
// yet unknown. The method will make a copy of the slice.
//
// Note, this method assumes that the database's lock is held!
func (db *Database) insertPreimage(hash common.Hash, preimage []byte) {
	if _, ok := db.preimages[hash]; ok {
		return
	}
	db.preimages[hash] = common.CopyBytes(preimage)
	db.preimagesSize += common.StorageSize(common.HashLength + len(preimage))
}

// insertPruningMark writes a new pruning mark to the memory database.
// Note, this method assumes that the database's lock is held!
func (db *Database) insertPruningMark(hash common.ExtHash, blockNum uint64) {
	db.pruningMarks = append(db.pruningMarks, database.PruningMark{
		Number: blockNum,
		Hash:   hash,
	})
}

// getCachedNode finds an encoded node in the trie node cache if enabled.
func (db *Database) getCachedNode(hash common.ExtHash) []byte {
	if db.trieNodeCache != nil {
		if enc := db.trieNodeCache.Get(hash[:]); enc != nil {
			memcacheCleanHitMeter.Mark(1)
			memcacheCleanReadMeter.Mark(int64(len(enc)))
			return enc
		}
	}
	return nil
}

// setCachedNode stores an encoded node to the trie node cache if enabled.
func (db *Database) setCachedNode(hash common.ExtHash, enc []byte) {
	if db.trieNodeCache != nil {
		db.trieNodeCache.Set(hash[:], enc)
		memcacheCleanWriteMeter.Mark(int64(len(enc)))
	}
}

func recordTrieCacheMiss() {
	memcacheCleanMissMeter.Mark(1)
}

// node retrieves a cached trie node from memory, or returns nil if node can be
// found in the memory cache.
func (db *Database) node(hash common.ExtHash) (n node, fromDB bool) {
	// Retrieve the node from the trie node cache if available
	if enc := db.getCachedNode(hash); enc != nil {
		if dec, err := decodeNode(hash[:], enc); err == nil {
			return dec, false
		} else {
			logger.Error("node from cached trie node fails to be decoded!", "err", err)
		}
	}

	// Retrieve the node from the state cache if available
	db.lock.RLock()
	node := db.nodes[hash]
	db.lock.RUnlock()
	if node != nil {
		return node.obj(hash), false
	}

	// Content unavailable in memory, attempt to retrieve from disk
	enc, err := db.diskDB.ReadTrieNode(hash)
	if err != nil || enc == nil {
		return nil, true
	}
	db.setCachedNode(hash, enc)
	recordTrieCacheMiss()
	return mustDecodeNode(hash[:], enc), true
}

// Node retrieves an encoded cached trie node from memory. If it cannot be found
// cached, the method queries the persistent database for the content.
func (db *Database) Node(hash common.ExtHash) ([]byte, error) {
	if common.EmptyExtHash(hash) {
		return nil, ErrZeroHashNode
	}
	// Retrieve the node from the trie node cache if available
	if enc := db.getCachedNode(hash); enc != nil {
		return enc, nil
	}

	// Retrieve the node from cache if available
	db.lock.RLock()
	node := db.nodes[hash]
	db.lock.RUnlock()

	if node != nil {
		return node.rlp(), nil
	}
	// Content unavailable in memory, attempt to retrieve from disk
	enc, err := db.diskDB.ReadTrieNode(hash)
	if err == nil && enc != nil {
		db.setCachedNode(hash, enc)
		recordTrieCacheMiss()
	}
	return enc, err
}

// NodeFromOld retrieves an encoded cached trie node from memory. If it cannot be found
// cached, the method queries the old persistent database for the content.
func (db *Database) NodeFromOld(hash common.ExtHash) ([]byte, error) {
	if common.EmptyExtHash(hash) {
		return nil, ErrZeroHashNode
	}
	// Retrieve the node from the trie node cache if available
	if enc := db.getCachedNode(hash); enc != nil {
		return enc, nil
	}

	// Retrieve the node from cache if available
	db.lock.RLock()
	node := db.nodes[hash]
	db.lock.RUnlock()

	if node != nil {
		return node.rlp(), nil
	}
	// Content unavailable in memory, attempt to retrieve from disk
	enc, err := db.diskDB.ReadTrieNodeFromOld(hash)
	if err == nil && enc != nil {
		db.setCachedNode(hash, enc)
		recordTrieCacheMiss()
	}
	return enc, err
}

// DoesExistCachedNode returns if the node exists on cached trie node in memory.
func (db *Database) DoesExistCachedNode(hash common.ExtHash) bool {
	// Retrieve the node from cache if available
	db.lock.RLock()
	_, ok := db.nodes[hash]
	db.lock.RUnlock()
	return ok
}

// DoesExistNodeInPersistent returns if the node exists on the persistent database or its cache.
func (db *Database) DoesExistNodeInPersistent(hash common.ExtHash) bool {
	// Retrieve the node from DB cache if available
	if enc := db.getCachedNode(hash); enc != nil {
		return true
	}

	// Content unavailable in DB cache, attempt to retrieve from disk
	enc, err := db.diskDB.ReadTrieNode(hash)
	if err == nil && enc != nil {
		return true
	}

	return false
}

// preimage retrieves a cached trie node pre-image from memory. If it cannot be
// found cached, the method queries the persistent database for the content.
func (db *Database) preimage(hash common.Hash) []byte {
	// Retrieve the node from cache if available
	db.lock.RLock()
	preimage := db.preimages[hash]
	db.lock.RUnlock()

	if preimage != nil {
		return preimage
	}
	// Content unavailable in memory, attempt to retrieve from disk
	return db.diskDB.ReadPreimage(hash)
}

// Nodes retrieves the hashes of all the nodes cached within the memory database.
// This method is extremely expensive and should only be used to validate internal
// states in test code.
func (db *Database) Nodes() []common.ExtHash {
	db.lock.RLock()
	defer db.lock.RUnlock()

	hashes := make([]common.ExtHash, 0, len(db.nodes))
	for hash := range db.nodes {
		if !common.EmptyExtHash(hash) { // Special case for "root" references/nodes
			hashes = append(hashes, hash)
		}
	}
	return hashes
}

// Reference adds a new reference from a parent node to a child node.
// This function is used to add reference between internal trie node
// and external node(e.g. storage trie root), all internal trie nodes
// are referenced together by database itself.
// Use ReferenceRoot to reference a state root, otherwise use Reference.
func (db *Database) ReferenceRoot(root common.Hash) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.reference(root.ExtendZero(), common.ExtHash{})
}

// Reference adds a new reference from a parent node to a child node.
// This function is used to add reference between internal trie node
// and external node(e.g. storage trie root), all internal trie nodes
// are referenced together by database itself.
// Use ReferenceRoot to reference a state root, otherwise use Reference.
func (db *Database) Reference(child common.ExtHash, parent common.ExtHash) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.reference(child, parent)
}

// reference is the private locked version of Reference.
func (db *Database) reference(child common.ExtHash, parent common.ExtHash) {
	// If the node does not exist, it's a node pulled from disk, skip
	node, ok := db.nodes[child]
	if !ok {
		return
	}
	// If the reference already exists, only duplicate for roots
	if db.nodes[parent].children == nil {
		db.nodes[parent].children = make(map[common.ExtHash]uint64)
	} else if _, ok = db.nodes[parent].children[child]; ok && !common.EmptyExtHash(parent) {
		return
	}
	node.parents++
	db.nodes[parent].children[child]++
}

// Dereference removes an existing reference from a state root node.
func (db *Database) Dereference(root common.Hash) {
	// Sanity check to ensure that the meta-root is not removed
	if common.EmptyHash(root) {
		logger.Error("Attempted to dereference the trie cache meta root")
		return
	}

	db.gcLock.Lock()
	defer db.gcLock.Unlock()

	db.lock.Lock()
	defer db.lock.Unlock()

	nodes, storage, start := len(db.nodes), db.nodesSize, time.Now()
	db.dereference(root.ExtendZero(), common.ExtHash{})

	db.gcnodes += uint64(nodes - len(db.nodes))
	db.gcsize += storage - db.nodesSize
	db.gctime += time.Since(start)

	memcacheGCTimeGauge.Update(int64(time.Since(start)))
	memcacheGCSizeMeter.Mark(int64(storage - db.nodesSize))
	memcacheGCNodesMeter.Mark(int64(nodes - len(db.nodes)))

	logger.Debug("Dereferenced trie from memory database", "nodes", nodes-len(db.nodes), "size", storage-db.nodesSize, "time", time.Since(start),
		"gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime, "livenodes", len(db.nodes), "livesize", db.nodesSize)
}

// dereference is the private locked version of Dereference.
func (db *Database) dereference(child common.ExtHash, parent common.ExtHash) {
	// Dereference the parent-child
	node := db.nodes[parent]

	if node.children != nil && node.children[child] > 0 {
		node.children[child]--
		if node.children[child] == 0 {
			delete(node.children, child)
		}
	}
	// If the node does not exist, it's a previously committed node.
	node, ok := db.nodes[child]
	if !ok {
		return
	}
	// If there are no more references to the child, delete it and cascade
	if node.parents > 0 {
		// This is a special cornercase where a node loaded from disk (i.e. not in the
		// memcache any more) gets reinjected as a new node (short node split into full,
		// then reverted into short), causing a cached node to have no parents. That is
		// no problem in itself, but don't make maxint parents out of it.
		node.parents--
	}
	if node.parents == 0 {
		// Remove the node from the flush-list
		db.removeNodeInFlushList(child)
		// Dereference all children and delete the node
		for _, hash := range node.childs() {
			db.dereference(hash, child)
		}
		delete(db.nodes, child)
		db.nodesSize -= common.StorageSize(common.HashLength + int(node.size))
	}
}

// Cap iteratively flushes old but still referenced trie nodes until the total
// memory usage goes below the given threshold.
func (db *Database) Cap(limit common.StorageSize) error {
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent database). This is ensured
	// by only uncaching existing data when the database write finalizes.
	db.lock.RLock()

	nodes, nodeSize, start := len(db.nodes), db.nodesSize, time.Now()
	preimagesSize := db.preimagesSize

	// db.nodesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted. For every useful node, we track 2 extra hashes as the flushlist.
	size := db.nodesSize + common.StorageSize((len(db.nodes)-1)*2*common.HashLength)

	// If the preimage cache got large enough, push to disk. If it's still small
	// leave for later to deduplicate writes.
	flushPreimages := db.preimagesSize > 4*1024*1024
	numPreimages := 0
	if flushPreimages {
		db.diskDB.WritePreimages(0, db.preimages)
		numPreimages = len(db.preimages)
	}
	db.diskDB.WritePruningMarks(db.pruningMarks)
	numPruningMarks := len(db.pruningMarks)

	// Keep committing nodes from the flush-list until we're below allowance
	oldest := db.oldest
	batch := db.diskDB.NewBatch(database.StateTrieDB)
	defer batch.Release()
	for size > limit && !common.EmptyExtHash(oldest) {
		// Fetch the oldest referenced node and push into the batch
		node := db.nodes[oldest]
		enc := node.rlp()
		db.diskDB.PutTrieNodeToBatch(batch, oldest, enc)
		if _, err := database.WriteBatchesOverThreshold(batch); err != nil {
			db.lock.RUnlock()
			return err
		}

		db.setCachedNode(oldest, enc)
		// Iterate to the next flush item, or abort if the size cap was achieved. Size
		// is the total size, including both the useful cached data (hash -> blob), as
		// well as the flushlist metadata (2*hash). When flushing items from the cache,
		// we need to reduce both.
		size -= common.StorageSize(3*common.HashLength + int(node.size))
		oldest = node.flushNext
	}
	// Flush out any remainder data from the last batch
	if _, err := database.WriteBatches(batch); err != nil {
		logger.Error("Failed to write flush list to disk", "err", err)
		db.lock.RUnlock()
		return err
	}

	db.lock.RUnlock()

	// Write successful, clear out the flushed data
	db.lock.Lock()
	defer db.lock.Unlock()

	if flushPreimages {
		db.preimages = make(map[common.Hash][]byte)
		db.preimagesSize = 0
	}
	db.pruningMarks = []database.PruningMark{}

	for db.oldest != oldest {
		node := db.nodes[db.oldest]
		delete(db.nodes, db.oldest)
		db.oldest = node.flushNext

		db.nodesSize -= common.StorageSize(common.HashLength + int(node.size))
	}
	if !common.EmptyExtHash(db.oldest) {
		db.nodes[db.oldest].flushPrev = common.ExtHash{}
	} else {
		db.newest = common.ExtHash{}
	}
	db.flushnodes += uint64(nodes - len(db.nodes))
	db.flushsize += nodeSize - db.nodesSize
	db.flushtime += time.Since(start)

	memcacheFlushTimeGauge.Update(int64(time.Since(start)))
	memcacheFlushSizeGauge.Update(int64(nodeSize - db.nodesSize))
	memcacheFlushNodesGauge.Update(int64(nodes - len(db.nodes)))

	logger.Info("Persisted nodes from memory database by Cap", "nodes", nodes-len(db.nodes),
		"size", nodeSize-db.nodesSize, "preimagesSize", preimagesSize-db.preimagesSize, "time", time.Since(start),
		"flushnodes", db.flushnodes, "flushsize", db.flushsize, "flushtime", db.flushtime, "livenodes", len(db.nodes),
		"livesize", db.nodesSize, "preimages", numPreimages, "pruningMarks", numPruningMarks)
	return nil
}

// commitResult contains the result from concurrent commit calls.
// key and val are nil if the commitResult indicates the end of
// concurrentCommit goroutine.
type commitResult struct {
	hash common.ExtHash
	val  []byte
}

func (db *Database) writeBatchNodes(node common.ExtHash) error {
	rootNode, ok := db.nodes[node]
	if !ok {
		return nil
	}

	// To limit the size of commitResult channel, we use commitResultChSizeLimit here.
	var resultCh chan commitResult
	if len(db.nodes) > commitResultChSizeLimit {
		resultCh = make(chan commitResult, commitResultChSizeLimit)
	} else {
		resultCh = make(chan commitResult, len(db.nodes))
	}
	numGoRoutines := len(rootNode.childs())
	for i, child := range rootNode.childs() {
		go db.concurrentCommit(child, resultCh, i)
	}

	batch := db.diskDB.NewBatch(database.StateTrieDB)
	defer batch.Release()
	for numGoRoutines > 0 {
		result := <-resultCh
		if common.EmptyExtHash(result.hash) && result.val == nil {
			numGoRoutines--
			continue
		}

		db.diskDB.PutTrieNodeToBatch(batch, result.hash, result.val)
		if _, err := database.WriteBatchesOverThreshold(batch); err != nil {
			return err
		}
	}

	enc := rootNode.rlp()
	db.diskDB.PutTrieNodeToBatch(batch, node, enc)
	if err := batch.Write(); err != nil {
		logger.Error("Failed to write trie to disk", "err", err)
		return err
	}
	db.setCachedNode(node, enc)

	return nil
}

func (db *Database) concurrentCommit(hash common.ExtHash, resultCh chan<- commitResult, childIndex int) {
	logger.Trace("concurrentCommit start", "childIndex", childIndex)
	defer logger.Trace("concurrentCommit end", "childIndex", childIndex)
	db.commit(hash, resultCh)
	resultCh <- commitResult{common.ExtHash{}, nil}
}

// Commit iterates over all the children of a particular node, writes them out
// to disk, forcefully tearing down all references in both directions.
// The root must be a state root.
//
// As a side effect, all pre-images accumulated up to this point are also written.
func (db *Database) Commit(root common.Hash, report bool, blockNum uint64) error {
	hash := root.ExtendZero()
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent database). This is ensured
	// by only uncaching existing data when the database write finalizes.
	db.lock.RLock()

	commitStart := time.Now()
	db.diskDB.WritePreimages(0, db.preimages)
	db.diskDB.WritePruningMarks(db.pruningMarks)
	numPreimages := len(db.preimages)
	numPruningMarks := len(db.pruningMarks)

	// Move the trie itself into the batch, flushing if enough data is accumulated
	numNodes, nodesSize := len(db.nodes), db.nodesSize
	if err := db.writeBatchNodes(hash); err != nil {
		db.lock.RUnlock()
		return err
	}

	db.lock.RUnlock()

	// Write successful, clear out the flushed data
	db.lock.Lock()
	defer db.lock.Unlock()

	db.preimages = make(map[common.Hash][]byte)
	db.preimagesSize = 0
	db.pruningMarks = []database.PruningMark{}

	uncacheStart := time.Now()
	db.uncache(hash)
	commitEnd := time.Now()

	memcacheUncacheTimeGauge.Update(int64(commitEnd.Sub(uncacheStart)))
	memcacheCommitTimeGauge.Update(int64(commitEnd.Sub(commitStart)))
	memcacheCommitSizeMeter.Mark(int64(nodesSize - db.nodesSize))
	memcacheCommitNodesMeter.Mark(int64(numNodes - len(db.nodes)))

	localLogger := logger.Info
	if !report {
		localLogger = logger.Debug
	}
	localLogger("Persisted trie from memory database", "blockNum", blockNum,
		"updated nodes", numNodes-len(db.nodes), "updated nodes size", nodesSize-db.nodesSize,
		"time", commitEnd.Sub(commitStart), "gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime,
		"livenodes", len(db.nodes), "livesize", db.nodesSize, "preimages", numPreimages, "pruningMarks", numPruningMarks)

	// Reset the garbage collection statistics
	db.gcnodes, db.gcsize, db.gctime = 0, 0, 0
	db.flushnodes, db.flushsize, db.flushtime = 0, 0, 0
	return nil
}

// commit iteratively encodes nodes from parents to child nodes.
func (db *Database) commit(hash common.ExtHash, resultCh chan<- commitResult) {
	node, ok := db.nodes[hash]
	if !ok {
		return
	}
	for _, child := range node.childs() {
		db.commit(child, resultCh)
	}
	enc := node.rlp()
	resultCh <- commitResult{hash, enc}

	db.setCachedNode(hash, enc)
}

// uncache is the post-processing step of a commit operation where the already
// persisted trie is removed from the cache. The reason behind the two-phase
// commit is to ensure consistent data availability while moving from memory
// to disk.
func (db *Database) uncache(hash common.ExtHash) {
	// If the node does not exists, we're done on this path
	node, ok := db.nodes[hash]
	if !ok {
		return
	}
	// Node still exists, remove it from the flush-list
	db.removeNodeInFlushList(hash)
	// Uncache the node's subtries and remove the node itself too
	for _, child := range node.childs() {
		db.uncache(child)
	}
	delete(db.nodes, hash)
	db.nodesSize -= common.StorageSize(common.HashLength + int(node.size))
}

// Size returns the current database size of the memory cache in front of the
// persistent database layer.
func (db *Database) Size() (common.StorageSize, common.StorageSize, common.StorageSize) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	// db.nodesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted. For every useful node, we track 2 extra hashes as the flushlist.
	flushlistSize := common.StorageSize((len(db.nodes) - 1) * 2 * common.HashLength)
	return db.nodesSize + flushlistSize, db.nodesSize, db.preimagesSize
}

// verifyIntegrity is a debug method to iterate over the entire trie stored in
// memory and check whether every node is reachable from the meta root. The goal
// is to find any errors that might cause memory leaks and or trie nodes to go
// missing.
//
// This method is extremely CPU and memory intensive, only use when must.
func (db *Database) verifyIntegrity() {
	// Iterate over all the cached nodes and accumulate them into a set
	reachable := map[common.ExtHash]struct{}{{}: {}}

	for child := range db.nodes[common.ExtHash{}].children {
		db.accumulate(child, reachable)
	}
	// Find any unreachable but cached nodes
	unreachable := []string{}
	for hash, node := range db.nodes {
		if _, ok := reachable[hash]; !ok {
			unreachable = append(unreachable, fmt.Sprintf("%x: {Node: %v, Parents: %d, Prev: %x, Next: %x}",
				hash, node.node, node.parents, node.flushPrev, node.flushNext))
		}
	}
	if len(unreachable) != 0 {
		panic(fmt.Sprintf("trie cache memory leak: %v", unreachable))
	}
}

// accumulate iterates over the trie defined by hash and accumulates all the
// cached children found in memory.
func (db *Database) accumulate(hash common.ExtHash, reachable map[common.ExtHash]struct{}) {
	// Mark the node reachable if present in the memory cache
	node, ok := db.nodes[hash]
	if !ok {
		return
	}
	reachable[hash] = struct{}{}

	// Iterate over all the children and accumulate them too
	for _, child := range node.childs() {
		db.accumulate(child, reachable)
	}
}

func (db *Database) removeNodeInFlushList(hash common.ExtHash) {
	node, ok := db.nodes[hash]
	if !ok {
		return
	}

	if hash == db.oldest && hash == db.newest {
		db.oldest = common.ExtHash{}
		db.newest = common.ExtHash{}
	} else if hash == db.oldest {
		db.oldest = node.flushNext
		db.nodes[node.flushNext].flushPrev = common.ExtHash{}
	} else if hash == db.newest {
		db.newest = node.flushPrev
		db.nodes[node.flushPrev].flushNext = common.ExtHash{}
	} else {
		db.nodes[node.flushPrev].flushNext = node.flushNext
		db.nodes[node.flushNext].flushPrev = node.flushPrev
	}
}

func (db *Database) getLastNodeHashInFlushList() common.ExtHash {
	var lastNodeHash common.ExtHash
	nodeHash := db.oldest
	for {
		if _, ok := db.nodes[nodeHash]; ok {
			lastNodeHash = nodeHash
		} else {
			logger.Debug("not found next noode in map of flush list")
			break
		}

		if !common.EmptyExtHash(db.nodes[nodeHash].flushNext) {
			nodeHash = db.nodes[nodeHash].flushNext
		} else {
			logger.Debug("found last noode in map of flush list")
			break
		}
	}
	return lastNodeHash
}

// UpdateMetricNodes updates the size of Database.nodes
func (db *Database) UpdateMetricNodes() {
	memcacheNodesGauge.Update(int64(len(db.nodes)))
	if db.trieNodeCache != nil {
		db.trieNodeCache.UpdateStats()
	}
}

var (
	errDisabledTrieNodeCache         = errors.New("trie node cache is disabled, nothing to save to file")
	errSavingTrieNodeCacheInProgress = errors.New("saving trie node cache has been triggered already")
)

func (db *Database) CanSaveTrieNodeCacheToFile() error {
	if db.trieNodeCache == nil {
		return errDisabledTrieNodeCache
	}
	if db.savingTrieNodeCacheTriggered {
		return errSavingTrieNodeCacheInProgress
	}
	return nil
}

// SaveTrieNodeCacheToFile saves the current cached trie nodes to file to reuse when the node restarts
func (db *Database) SaveTrieNodeCacheToFile(filePath string, concurrency int) {
	db.savingTrieNodeCacheTriggered = true
	start := time.Now()
	logger.Info("start saving cache to file",
		"filePath", filePath, "concurrency", concurrency)
	if err := db.trieNodeCache.SaveToFile(filePath, concurrency); err != nil {
		logger.Error("failed to save cache to file",
			"filePath", filePath, "elapsed", time.Since(start), "err", err)
	} else {
		logger.Info("successfully saved cache to file",
			"filePath", filePath, "elapsed", time.Since(start))
	}
	db.savingTrieNodeCacheTriggered = false
}

// DumpPeriodically atomically saves fast cache data to the given dir with the specified interval.
func (db *Database) SaveCachePeriodically(c *TrieNodeCacheConfig, stopCh <-chan struct{}) {
	rand.Seed(time.Now().UnixNano())
	randomVal := 0.5 + rand.Float64()/2.0 // 0.5 <= randomVal < 1.0
	startTime := time.Duration(int(randomVal * float64(c.FastCacheSavePeriod)))
	logger.Info("first periodic cache saving will be triggered", "after", startTime)

	timer := time.NewTimer(startTime)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := db.CanSaveTrieNodeCacheToFile(); err != nil {
				logger.Warn("failed to trigger periodic cache saving", "err", err)
				continue
			}
			db.SaveTrieNodeCacheToFile(c.FastCacheFileDir, 1)
			timer.Reset(c.FastCacheSavePeriod)
		case <-stopCh:
			return
		}
	}
}

// NodeInfo is a struct used for collecting trie statistics
type NodeInfo struct {
	Depth    int  // 0 if not a leaf node
	Finished bool // true if the uppermost call is finished
}

// CollectChildrenStats collects the depth of the trie recursively
func (db *Database) CollectChildrenStats(node common.ExtHash, depth int, resultCh chan<- NodeInfo) {
	n, _ := db.node(node)
	if n == nil {
		return
	}
	// retrieve the children of the given node
	childrenNodes, err := db.NodeChildren(node)
	if err != nil {
		logger.Error("failed to retrieve the children nodes",
			"node", node.String(), "err", err)
		return
	}
	// write the depth of the node only if the node is a leaf node, otherwise set 0
	resultDepth := 0
	if len(childrenNodes) == 0 {
		resultDepth = depth
	}
	// send the result to the channel and iterate its children
	resultCh <- NodeInfo{Depth: resultDepth}
	for _, child := range childrenNodes {
		db.CollectChildrenStats(child, depth+1, resultCh)
	}
}
