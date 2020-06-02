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

// This file is derived from ethdb/iterator.go (2020/05/20).
// Modified and improved for the klaytn development.

package state

import (
	"bytes"
	"fmt"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
)

// NodeIterator is an iterator to traverse the entire state trie post-order,
// including all of the contract code and contract state tries.
type NodeIterator struct {
	state *StateDB // State being iterated

	stateIt statedb.NodeIterator // Primary iterator for the global state trie
	dataIt  statedb.NodeIterator // Secondary iterator for the data trie of a contract

	accountHash common.Hash // Hash of the node containing the account
	codeHash    common.Hash // Hash of the contract source code
	Code        []byte      // Source code associated with a contract

	Type   string
	Hash   common.Hash // Hash of the current entry being iterated (nil if not standalone)
	Parent common.Hash // Hash of the first full ancestor node (nil if current is the root)
	Path   []byte      // the hex-encoded path to the current node.

	Error error // Failure set in case of an internal error in the iteratord
}

// NewNodeIterator creates an post-order state node iterator.
func NewNodeIterator(state *StateDB) *NodeIterator {
	return &NodeIterator{
		state: state,
	}
}

// Next moves the iterator to the next node, returning whether there are any
// further nodes. In case of an internal error this method returns false and
// sets the Error field to the encountered failure.
func (it *NodeIterator) Next() bool {
	// If the iterator failed previously, don't do anything
	if it.Error != nil {
		return false
	}
	// Otherwise step forward with the iterator and report any errors
	if err := it.step(); err != nil {
		it.Error = err
		return false
	}
	return it.retrieve()
}

// step moves the iterator to the next entry of the state trie.
func (it *NodeIterator) step() error {
	// Abort if we reached the end of the iteration
	if it.state == nil {
		return nil
	}
	// Initialize the iterator if we've just started
	if it.stateIt == nil {
		it.stateIt = it.state.trie.NodeIterator(nil)
	}
	// If we had data nodes previously, we surely have at least state nodes
	if it.dataIt != nil {
		if cont := it.dataIt.Next(true); !cont {
			if it.dataIt.Error() != nil {
				return it.dataIt.Error()
			}
			it.dataIt = nil
		}
		return nil
	}
	// If we had source code previously, discard that
	if it.Code != nil {
		it.Code = nil
		return nil
	}
	// Step to the next state trie node, terminating if we're out of nodes
	if cont := it.stateIt.Next(true); !cont {
		if it.stateIt.Error() != nil {
			return it.stateIt.Error()
		}
		it.state, it.stateIt = nil, nil
		return nil
	}
	// If the state trie node is an internal entry, leave as is
	if !it.stateIt.Leaf() {
		return nil
	}
	// Otherwise we've reached an account node, initiate data iteration

	serializer := account.NewAccountSerializer()
	if err := rlp.Decode(bytes.NewReader(it.stateIt.LeafBlob()), serializer); err != nil {
		return err
	}
	obj := serializer.GetAccount()

	if pa := account.GetProgramAccount(obj); pa != nil {
		dataTrie, err := it.state.db.OpenStorageTrie(pa.GetStorageRoot())
		if err != nil {
			return err
		}
		it.dataIt = dataTrie.NodeIterator(nil)
		if !it.dataIt.Next(true) {
			it.dataIt = nil
		}

		it.codeHash = common.BytesToHash(pa.GetCodeHash())
		//addrHash := common.BytesToHash(it.stateIt.LeafKey())
		it.Code, err = it.state.db.ContractCode(common.BytesToHash(pa.GetCodeHash()))
		if err != nil {
			return fmt.Errorf("code %x: %v", pa.GetCodeHash(), err)
		}
	}
	it.accountHash = it.stateIt.Parent()
	return nil
}

// retrieve pulls and caches the current state entry the iterator is traversing.
// The method returns whether there are any more data left for inspection.
func (it *NodeIterator) retrieve() bool {
	// Clear out any previously set values
	it.Hash = common.Hash{}
	it.Path = []byte{}

	// If the iteration's done, return no available data
	if it.state == nil {
		return false
	}
	// Otherwise retrieve the current entry
	switch {
	case it.dataIt != nil:
		it.Type = "storage"
		if it.dataIt.Leaf() {
			it.Type = "storage_leaf"
		}

		it.Hash, it.Parent, it.Path = it.dataIt.Hash(), it.dataIt.Parent(), it.dataIt.Path()

		if it.Parent == (common.Hash{}) {
			it.Parent = it.accountHash
		}
	case it.Code != nil:
		it.Type = "code"
		it.Hash, it.Parent = it.codeHash, it.accountHash
	case it.stateIt != nil:
		it.Type = "state"
		if it.stateIt.Leaf() {
			it.Type = "state_leaf"
		}

		it.Hash, it.Parent, it.Path = it.stateIt.Hash(), it.stateIt.Parent(), it.stateIt.Path()
	}
	return true
}

// CheckStateConsistency checks the consistency of all state/storage trie of given two state database.
func CheckStateConsistency(oldDB database.DBManager, newDB database.DBManager, root common.Hash) error {
	// Create and iterate a state trie rooted in a sub-node
	oldState, err := New(root, NewDatabase(oldDB))
	if err != nil {
		return err
	}

	newState, err := New(root, NewDatabase(newDB))
	if err != nil {
		return err
	}

	oldIt := NewNodeIterator(oldState)
	newIt := NewNodeIterator(newState)

	cnt := 0
	nodes := make(map[common.Hash]bool)

	for oldIt.Next() {
		cnt++

		if !newIt.Next() {
			return fmt.Errorf("newDB iterator finished earlier : oldIt.Hash(%v) oldIt.Parent(%v)", oldIt.Hash, oldIt.Parent)
		}

		if oldIt.Hash != newIt.Hash {
			return fmt.Errorf("mismatched hash oldIt.Hash : oldIt.Hash(%v) newIt.Hash(%v)", oldIt.Hash, newIt.Hash)
		}

		if oldIt.Parent != newIt.Parent {
			return fmt.Errorf("mismatched parent hash : oldIt.Parent(%v) newIt.Parent(%v)", oldIt.Parent, newIt.Parent)
		}

		if !bytes.Equal(oldIt.Path, newIt.Path) {
			return fmt.Errorf("mismatched path : oldIt.path(%v) newIt.path(%v)",
				statedb.HexKeyPathToHashString(oldIt.Path), statedb.HexKeyPathToHashString(newIt.Path))
		}

		if oldIt.Code != nil {
			if newIt.Code != nil {
				if !bytes.Equal(oldIt.Code, newIt.Code) {
					return fmt.Errorf("mismatched code : oldIt.Code(%v) newIt.Code(%v)", oldIt.Code, newIt.Code)
				}
			} else {
				return fmt.Errorf("mismatched code : oldIt.Code(%v) newIt.Code(nil)", string(oldIt.Code))
			}
		} else {
			if newIt.Code != nil {
				return fmt.Errorf("mismatched code : oldIt.Code(nil) newIt.Code(%v)", string(newIt.Code))
			}
		}

		if !common.EmptyHash(oldIt.Hash) {
			nodes[oldIt.Hash] = true
		}

		logger.Trace("CheckStateConsistency next",
			"type", oldIt.Type,
			"hash", oldIt.Hash.String(),
			"parent", oldIt.Parent.String(),
			"path", statedb.HexKeyPathToHashString(oldIt.Path))
	}

	if newIt.Next() {
		return fmt.Errorf("oldDB iterator finished earlier  : newIt.Hash(%v) newIt.Parent(%v)", newIt.Hash, newIt.Parent)
	}

	logger.Info("CheckStateConsistency is completed", "cnt", cnt, "cnt without duplication", len(nodes))
	return nil
}
