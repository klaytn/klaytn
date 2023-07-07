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
// This file is derived from trie/hasher.go (2018/06/04).
// Modified and improved for the klaytn development.

package statedb

import (
	"hash"
	"sync"

	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/rlp"
)

type hasherOpts struct {
	onleaf      LeafCallback
	pruning     bool // If pruning is true, non-root nodes are attached a fresh nonce.
	storageRoot bool // If both pruning and storageRoot are true, the root node is attached a fresh nonce.
}

type hasher struct {
	hasherOpts
	tmp sliceBuffer
	sha KeccakState
}

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

type sliceBuffer []byte

func (b *sliceBuffer) Write(data []byte) (n int, err error) {
	*b = append(*b, data...)
	return len(data), nil
}

func (b *sliceBuffer) Reset() {
	*b = (*b)[:0]
}

// hashers live in a global db.
var hasherPool = sync.Pool{
	New: func() interface{} {
		return &hasher{
			tmp: make(sliceBuffer, 0, 550), // cap is as large as a full fullNode.
			sha: sha3.NewKeccak256().(KeccakState),
		}
	},
}

func newHasher(opts *hasherOpts) *hasher {
	h := hasherPool.Get().(*hasher)
	if opts == nil {
		opts = &hasherOpts{}
	}
	h.hasherOpts = *opts
	return h
}

func returnHasherToPool(h *hasher) {
	hasherPool.Put(h)
}

// hashRoot is similar to hashNode() but adds special treatment for the root node.
func (h *hasher) hashRoot(n node, db *Database, force bool) (node, node) {
	return h.hashNode(n, db, force, true)
}

// hash is similar to hashNode() but assumes that the node is not a root node.
func (h *hasher) hash(n node, db *Database, force bool) (node, node) {
	return h.hashNode(n, db, force, false)
}

// hashNode collapses a node down into a hash node, also returning a copy of the
// original node initialized with the computed hash to replace the original one.
//
// hashNode is for hasher's internal use only.
// Please use hashRoot() or hash() for readability.
func (h *hasher) hashNode(n node, db *Database, force bool, onRoot bool) (node, node) {
	// If we're not storing the node, just hashing, use available cached data
	if hash, dirty := n.cache(); hash != nil {
		if db == nil {
			return hash, n
		}
		if !dirty {
			switch n.(type) {
			case *fullNode, *shortNode:
				return hash, hash
			default:
				return hash, n
			}
		}
	}
	// Trie not processed yet or needs storage, walk the children
	collapsed, cached := h.hashChildren(n, db, onRoot)
	hashed, lenEncoded := h.store(collapsed, db, force, onRoot)
	// Cache the hash of the node for later reuse and remove
	// the dirty flag in commit mode. It's fine to assign these values directly
	// without copying the node first because hashChildren copies it.
	cachedHash, _ := hashed.(hashNode)
	switch cn := cached.(type) {
	case *shortNode:
		cn.flags.hash = cachedHash
		cn.flags.lenEncoded = lenEncoded
		if db != nil {
			cn.flags.dirty = false
		}
	case *fullNode:
		cn.flags.hash = cachedHash
		cn.flags.lenEncoded = lenEncoded
		if db != nil {
			cn.flags.dirty = false
		}
	}
	return hashed, cached
}

// hashChildren replaces the children of a node with their hashes if the encoded
// size of the child is larger than a hash, returning the collapsed node as well
// as a replacement for the original node with the child hashes cached in.
func (h *hasher) hashChildren(original node, db *Database, onRoot bool) (node, node) {
	switch n := original.(type) {
	case *shortNode:
		// Hash the short node's child, caching the newly hashed subtree
		collapsed, cached := n.copy(), n.copy()
		collapsed.Key = hexToCompact(n.Key)
		cached.Key = common.CopyBytes(n.Key)

		if _, ok := n.Val.(valueNode); !ok {
			collapsed.Val, cached.Val = h.hash(n.Val, db, false)
		}
		return collapsed, cached

	case *fullNode:
		// Hash the full node's children, caching the newly hashed subtrees
		collapsed, cached := n.copy(), n.copy()

		if onRoot {
			var wg sync.WaitGroup
			wg.Add(16)
			for i := 0; i < 16; i++ {
				if n.Children[i] != nil {
					go func(i int) {
						childHasher := newHasher(&h.hasherOpts)
						collapsed.Children[i], cached.Children[i] = childHasher.hash(n.Children[i], db, false)
						returnHasherToPool(childHasher)
						wg.Done()
					}(i)
				} else {
					wg.Done()
				}
			}
			wg.Wait()
		} else {
			for i := 0; i < 16; i++ {
				if n.Children[i] != nil {
					collapsed.Children[i], cached.Children[i] = h.hash(n.Children[i], db, false)
				}
			}
		}
		cached.Children[16] = n.Children[16]
		return collapsed, cached

	default:
		// Value and hash nodes don't have children so they're left as were
		return n, original
	}
}

// store hashes the node n and if we have a storage layer specified, it writes
// the key/value pair to it and tracks any node->child references as well as any
// node->external trie references.
func (h *hasher) store(n node, db *Database, force bool, onRoot bool) (node, uint16) {
	// Don't store hashes or empty nodes.
	if _, isHash := n.(hashNode); n == nil || isHash {
		return n, 0
	}
	// hash is for the merkle proof. hash = Keccak(rlp.Encode(nodeForHashing(n)))
	// lenEncoded is for Database size accounting. lenEncoded = len(rlp.Encode(nodeForStoring(n)))
	hash, _ := n.cache()
	lenEncoded := n.lenEncoded()

	// Calculate lenEncoded if not set
	if hash == nil || lenEncoded == 0 {
		// Generate the RLP encoding of the node for database storing
		h.tmp.Reset()
		if err := rlp.Encode(&h.tmp, h.nodeForStoring(n)); err != nil {
			panic("encode error: " + err.Error())
		}
		lenEncoded = uint16(len(h.tmp))
	}
	if lenEncoded < 32 && !force {
		return n, lenEncoded // Nodes smaller than 32 bytes are stored inside their parent
	}

	// Calculate hash if not set
	if hash == nil {
		// Generate the RLP encoding of the node for Merkle hashing
		h.tmp.Reset()
		if err := rlp.Encode(&h.tmp, h.nodeForHashing(n)); err != nil {
			panic("encode error: " + err.Error())
		}
		hash = h.makeHashNode(h.tmp, onRoot)
	}

	if db != nil {
		// We are pooling the trie nodes into an intermediate memory cache
		hash := common.BytesToExtHash(hash)

		db.lock.Lock()
		db.insert(hash, lenEncoded, h.nodeForStoring(n))
		db.lock.Unlock()

		// Track external references from account->storage trie
		if h.onleaf != nil {
			switch n := n.(type) {
			case *shortNode:
				if child, ok := n.Val.(valueNode); ok {
					h.onleaf(nil, nil, child, hash, 0)
				}
			case *fullNode:
				for i := 0; i < 16; i++ {
					if child, ok := n.Children[i].(valueNode); ok {
						h.onleaf(nil, nil, child, hash, 0)
					}
				}
			}
		}
	}
	return hash, lenEncoded
}

func (h *hasher) makeHashNode(data []byte, onRoot bool) hashNode {
	var hash common.Hash
	h.sha.Reset()
	h.sha.Write(data)
	h.sha.Read(hash[:])
	if h.pruning && (h.storageRoot || !onRoot) {
		return hash.Extend().Bytes()
	} else {
		return hash.ExtendLegacy().Bytes()
	}
}

func (h *hasher) nodeForHashing(original node) node {
	return unextendNode(original, false)
}

func (h *hasher) nodeForStoring(original node) node {
	return unextendNode(original, true)
}

func unextendNode(original node, preserveExtHash bool) node {
	switch n := original.(type) {
	case *shortNode:
		stored := n.copy()
		stored.Val = unextendNode(n.Val, preserveExtHash)
		return stored
	case *fullNode:
		stored := n.copy()
		for i, child := range stored.Children {
			stored.Children[i] = unextendNode(child, preserveExtHash)
		}
		return stored
	case hashNode:
		exthash := common.BytesToExtHash(n)
		if exthash.IsLegacy() { // Always unextend ExtHashLegacy
			return hashNode(exthash.Unextend().Bytes())
		} else if !preserveExtHash { // It's ExtHash and not preserving ExtHash (for merkle hash)
			return hashNode(exthash.Unextend().Bytes())
		} else { // It's ExtHash and preserving ExtHash (for storing)
			return n
		}
	case valueNode:
		if !preserveExtHash {
			return valueNode(account.UnextendSerializedAccount(n))
		} else {
			// ExtHashLegacy should have been always unextended by AccountSerializer,
			// hence no need to check IsLegacy() here.
			return n
		}
	default:
		return n
	}
}
