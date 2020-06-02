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

	// If the iteration's done, return no available data
	if it.state == nil {
		return false
	}
	// Otherwise retrieve the current entry
	switch {
	case it.dataIt != nil:
		it.Hash, it.Parent, it.Path = it.dataIt.Hash(), it.dataIt.Parent(), it.dataIt.Path()
		if it.Parent == (common.Hash{}) {
			it.Parent = it.accountHash
		}
	case it.Code != nil:
		it.Hash, it.Parent = it.codeHash, it.accountHash
	case it.stateIt != nil:
		it.Hash, it.Parent, it.Path = it.stateIt.Hash(), it.stateIt.Parent(), it.stateIt.Path()
	}
	return true
}
