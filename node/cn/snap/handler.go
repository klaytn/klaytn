// Modifications Copyright 2022 The klaytn Authors
// Copyright 2020 The go-ethereum Authors
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
// This file is derived from eth/protocols/snap/handler.go (2022/06/29).
// Modified and improved for the klaytn development.

package snap

import (
	"bytes"
	"fmt"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/snapshot"
	"github.com/klaytn/klaytn/storage/statedb"
)

const (
	// softResponseLimit is the target maximum size of replies to data retrievals.
	softResponseLimit = 2 * 1024 * 1024

	// maxCodeLookups is the maximum number of bytecodes to serve. This number is
	// there to limit the number of disk lookups.
	maxCodeLookups = 1024

	// stateLookupSlack defines the ratio by how much a state response can exceed
	// the requested limit in order to try and avoid breaking up contracts into
	// multiple packages and proving them.
	stateLookupSlack = 0.1

	// maxTrieNodeLookups is the maximum number of state trie nodes to serve. This
	// number is there to limit the number of disk lookups.
	maxTrieNodeLookups = 1024

	// maxTrieNodeTimeSpent is the maximum time we should spend on looking up trie nodes.
	// If we spend too much time, then it's a fairly high chance of timing out
	// at the remote side, which means all the work is in vain.
	maxTrieNodeTimeSpent = 5 * time.Second
)

// Handler is a callback to invoke from an outside runner after the boilerplate
// exchanges have passed.
type Handler func(peer *Peer) error

type SnapshotReader interface {
	StateCache() state.Database
	Snapshots() *snapshot.Tree
	ContractCode(hash common.ExtHash) ([]byte, error)
	ContractCodeWithPrefix(hash common.ExtHash) ([]byte, error)
}

type SnapshotDownloader interface {
	// DeliverSnapPacket is invoked from a peer's message handler when it transmits a
	// data packet for the local node to consume.
	DeliverSnapPacket(peer *Peer, packet Packet) error
}

// Handle is the callback invoked to manage the life cycle of a `snap` peer.
// When this function terminates, the peer is disconnected.
func Handle(reader SnapshotReader, downloader SnapshotDownloader, peer *Peer) error {
	for {
		if err := HandleMessage(reader, downloader, peer); err != nil {
			peer.Log().Debug("Message handling failed in `snap`", "err", err)
			return err
		}
	}
}

// HandleMessage is invoked whenever an inbound message is received from a
// remote peer on the `snap` protocol. The remote connection is torn down upon
// returning any error.
func HandleMessage(reader SnapshotReader, downloader SnapshotDownloader, peer *Peer) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := peer.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > maxMessageSize {
		return fmt.Errorf("%w: %v > %v", errMsgTooLarge, msg.Size, maxMessageSize)
	}
	defer msg.Discard()
	start := time.Now()
	// TODO-Klaytn-SnapSync add metric to track the amount of time it takes to serve the request and run the manager
	// Handle the message depending on its contents
	switch {
	case msg.Code == GetAccountRangeMsg:
		// Decode the account retrieval request
		var req GetAccountRangePacket
		if err := msg.Decode(&req); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		// Service the request, potentially returning nothing in case of errors
		accounts, proofs := ServiceGetAccountRangeQuery(reader, &req)

		// Send back anything accumulated (or empty in case of errors)
		return p2p.Send(peer.rw, AccountRangeMsg, &AccountRangePacket{
			ID:       req.ID,
			Accounts: accounts,
			Proof:    proofs,
		})

	case msg.Code == AccountRangeMsg:
		// A range of accounts arrived to one of our previous requests
		res := new(AccountRangePacket)
		if err := msg.Decode(res); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		// Ensure the range is monotonically increasing
		for i := 1; i < len(res.Accounts); i++ {
			if bytes.Compare(res.Accounts[i-1].Hash[:], res.Accounts[i].Hash[:]) >= 0 {
				return fmt.Errorf("accounts not monotonically increasing: #%d [%x] vs #%d [%x]", i-1, res.Accounts[i-1].Hash[:], i, res.Accounts[i].Hash[:])
			}
		}
		requestTracker.Fulfil(peer.id, peer.version, AccountRangeMsg, res.ID)

		return downloader.DeliverSnapPacket(peer, res)

	case msg.Code == GetStorageRangesMsg:
		// Decode the storage retrieval request
		var req GetStorageRangesPacket
		if err := msg.Decode(&req); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		// Service the request, potentially returning nothing in case of errors
		slots, proofs := ServiceGetStorageRangesQuery(reader, &req)

		// Send back anything accumulated (or empty in case of errors)
		return p2p.Send(peer.rw, StorageRangesMsg, &StorageRangesPacket{
			ID:    req.ID,
			Slots: slots,
			Proof: proofs,
		})

	case msg.Code == StorageRangesMsg:
		// A range of storage slots arrived to one of our previous requests
		res := new(StorageRangesPacket)
		if err := msg.Decode(res); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		// Ensure the ranges are monotonically increasing
		for i, slots := range res.Slots {
			for j := 1; j < len(slots); j++ {
				if bytes.Compare(slots[j-1].Hash[:], slots[j].Hash[:]) >= 0 {
					return fmt.Errorf("storage slots not monotonically increasing for account #%d: #%d [%x] vs #%d [%x]", i, j-1, slots[j-1].Hash[:], j, slots[j].Hash[:])
				}
			}
		}
		requestTracker.Fulfil(peer.id, peer.version, StorageRangesMsg, res.ID)

		return downloader.DeliverSnapPacket(peer, res)

	case msg.Code == GetByteCodesMsg:
		// Decode bytecode retrieval request
		var req GetByteCodesPacket
		if err := msg.Decode(&req); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		// Service the request, potentially returning nothing in case of errors
		codes := ServiceGetByteCodesQuery(reader, &req)

		// Send back anything accumulated (or empty in case of errors)
		return p2p.Send(peer.rw, ByteCodesMsg, &ByteCodesPacket{
			ID:    req.ID,
			Codes: codes,
		})

	case msg.Code == ByteCodesMsg:
		// A batch of byte codes arrived to one of our previous requests
		res := new(ByteCodesPacket)
		if err := msg.Decode(res); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		requestTracker.Fulfil(peer.id, peer.version, ByteCodesMsg, res.ID)

		return downloader.DeliverSnapPacket(peer, res)

	case msg.Code == GetTrieNodesMsg:
		// Decode trie node retrieval request
		var req GetTrieNodesPacket
		if err := msg.Decode(&req); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		// Service the request, potentially returning nothing in case of errors
		nodes, err := ServiceGetTrieNodesQuery(reader, &req, start)
		if err != nil {
			return err
		}
		// Send back anything accumulated (or empty in case of errors)
		return p2p.Send(peer.rw, TrieNodesMsg, &TrieNodesPacket{
			ID:    req.ID,
			Nodes: nodes,
		})

	case msg.Code == TrieNodesMsg:
		// A batch of trie nodes arrived to one of our previous requests
		res := new(TrieNodesPacket)
		if err := msg.Decode(res); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		requestTracker.Fulfil(peer.id, peer.version, TrieNodesMsg, res.ID)

		return downloader.DeliverSnapPacket(peer, res)

	default:
		return fmt.Errorf("%w: %v", errInvalidMsgCode, msg.Code)
	}
}

// ServiceGetAccountRangeQuery assembles the response to an account range query.
// It is exposed to allow external packages to test protocol behavior.
func ServiceGetAccountRangeQuery(chain SnapshotReader, req *GetAccountRangePacket) ([]*AccountData, [][]byte) {
	if req.Bytes > softResponseLimit {
		req.Bytes = softResponseLimit
	}
	// TODO-Klaytn-SnapSync investigate the cache pollution
	// Retrieve the requested state and bail out if non existent
	tr, err := statedb.NewTrie(req.Root.ToRootExtHash(), chain.StateCache().TrieDB())
	if err != nil {
		return nil, nil
	}
	it, err := chain.Snapshots().AccountIterator(req.Root.ToRootExtHash(), req.Origin.ToRootExtHash())
	if err != nil {
		return nil, nil
	}
	// Iterate over the requested range and pile accounts up
	var (
		accounts []*AccountData
		size     uint64
		last     common.ExtHash
	)
	for it.Next() {
		hash, account := it.Hash(), common.CopyBytes(it.Account())

		// Track the returned interval for the Merkle proofs
		last = hash

		// Assemble the reply item
		size += uint64(common.HashLength + len(account))
		accounts = append(accounts, &AccountData{
			Hash: hash.ToHash(),
			Body: account,
		})
		// If we've exceeded the request threshold, abort
		if bytes.Compare(hash[:], req.Limit[:]) >= 0 {
			break
		}
		// TODO-Klaytn-SnapSync check if the size is much larger than soft response limit
		if size > req.Bytes {
			break
		}
	}
	it.Release()

	// Generate the Merkle proofs for the first and last account
	proof := NewNodeSet()
	if err := tr.Prove(req.Origin[:], 0, proof); err != nil {
		logger.Warn("Failed to prove account range", "origin", req.Origin, "err", err)
		return nil, nil
	}
	if last.ToHash() != (common.Hash{}) {
		if err := tr.Prove(last[:], 0, proof); err != nil {
			logger.Warn("Failed to prove account range", "last", last, "err", err)
			return nil, nil
		}
	}
	nodeList := proof.NodeList()
	proofs := make([][]byte, 0, len(nodeList))
	for _, blob := range nodeList {
		proofs = append(proofs, blob)
	}
	return accounts, proofs
}

func ServiceGetStorageRangesQuery(chain SnapshotReader, req *GetStorageRangesPacket) ([][]*StorageData, [][]byte) {
	if req.Bytes > softResponseLimit {
		req.Bytes = softResponseLimit
	}
	// TODO-Klaytn-SnapSync Do we want to enforce > 0 accounts and 1 account if origin is set?
	// TODO-Klaytn-SnapSync   - Logging locally is not ideal as remote faulst annoy the local user
	// TODO-Klaytn-SnapSync   - Dropping the remote peer is less flexible wrt client bugs (slow is better than non-functional)

	// Calculate the hard limit at which to abort, even if mid storage trie
	hardLimit := uint64(float64(req.Bytes) * (1 + stateLookupSlack))

	// Retrieve storage ranges until the packet limit is reached
	var (
		slots  = make([][]*StorageData, 0, len(req.Accounts))
		proofs [][]byte
		size   uint64
	)
	for _, accountHash := range req.Accounts {
		// If we've exceeded the requested data limit, abort without opening
		// a new storage range (that we'd need to prove due to exceeded size)
		if size >= req.Bytes {
			break
		}
		// The first account might start from a different origin and end sooner
		var origin common.Hash
		if len(req.Origin) > 0 {
			origin, req.Origin = common.BytesToHash(req.Origin), nil
		}
		limit := common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
		if len(req.Limit) > 0 {
			limit, req.Limit = common.BytesToHash(req.Limit), nil
		}
		// Retrieve the requested state and bail out if non existent
		it, err := chain.Snapshots().StorageIterator(req.Root.ToRootExtHash(), accountHash.ToRootExtHash(), origin.ToRootExtHash()) // 2.3M_BAD_BLOCK_CODE
		// it, err := chain.Snapshots().StorageIterator(req.Root.ToRootExtHash(), accountHash.LegacyToExtHash(), origin.ToRootExtHash())
		if err != nil {
			return nil, nil
		}
		// Iterate over the requested range and pile slots up
		var (
			storage []*StorageData
			last    common.ExtHash
			abort   bool
		)
		for it.Next() {
			if size >= hardLimit {
				abort = true
				break
			}
			hash, slot := it.Hash(), common.CopyBytes(it.Slot())

			// Track the returned interval for the Merkle proofs
			last = hash

			// Assemble the reply item
			size += uint64(common.HashLength + len(slot))
			storage = append(storage, &StorageData{
				Hash: hash.ToHash(),
				Body: slot,
			})
			// If we've exceeded the request threshold, abort
			if bytes.Compare(hash[:], limit[:]) >= 0 {
				break
			}
		}
		slots = append(slots, storage)
		it.Release()

		// Generate the Merkle proofs for the first and last storage slot, but
		// only if the response was capped. If the entire storage trie included
		// in the response, no need for any proofs.
		if origin != (common.Hash{}) || abort {
			// Request started at a non-zero hash or was capped prematurely, add
			// the endpoint Merkle proofs
			accTrie, err := statedb.NewTrie(req.Root.ToRootExtHash(), chain.StateCache().TrieDB())
			if err != nil {
				return nil, nil
			}
			serializer := account.NewAccountSerializer()
			if err := rlp.DecodeBytes(accTrie.Get(accountHash[:]), serializer); err != nil {
				return nil, nil
			}
			acc := serializer.GetAccount()
			pacc := account.GetProgramAccount(acc)
			if pacc == nil {
				// TODO-Klaytn-SnapSync it would be better to continue rather than return. Do not waste the completed job until now.
				return nil, nil
			}
			stTrie, err := statedb.NewTrie(pacc.GetStorageRoot(), chain.StateCache().TrieDB())
			if err != nil {
				return nil, nil
			}
			proof := NewNodeSet()
			if err := stTrie.Prove(origin[:], 0, proof); err != nil {
				logger.Warn("Failed to prove storage range", "origin", req.Origin, "err", err)
				return nil, nil
			}
			if last.ToHash() != (common.Hash{}) {
				if err := stTrie.Prove(last[:], 0, proof); err != nil {
					logger.Warn("Failed to prove storage range", "last", last, "err", err)
					return nil, nil
				}
			}
			for _, blob := range proof.NodeList() {
				proofs = append(proofs, blob)
			}
			// Proof terminates the reply as proofs are only added if a node
			// refuses to serve more data (exception when a contract fetch is
			// finishing, but that's that).
			break
		}
	}
	return slots, proofs
}

// ServiceGetByteCodesQuery assembles the response to a byte codes query.
// It is exposed to allow external packages to test protocol behavior.
func ServiceGetByteCodesQuery(chain SnapshotReader, req *GetByteCodesPacket) [][]byte {
	if req.Bytes > softResponseLimit {
		req.Bytes = softResponseLimit
	}
	if len(req.Hashes) > maxCodeLookups {
		req.Hashes = req.Hashes[:maxCodeLookups]
	}
	// Retrieve bytecodes until the packet size limit is reached
	var (
		codes [][]byte
		bytes uint64
	)
	for _, hash := range req.Hashes {
		if hash == emptyCode {
			// Peers should not request the empty code, but if they do, at
			// least sent them back a correct response without db lookups
			codes = append(codes, []byte{})
		} else if blob, err := chain.ContractCode(hash.ToExtHash()); err == nil {
			codes = append(codes, blob)
			bytes += uint64(len(blob))
		}
		if bytes > req.Bytes {
			break
		}
	}
	return codes
}

// ServiceGetTrieNodesQuery assembles the response to a trie nodes query.
// It is exposed to allow external packages to test protocol behavior.
func ServiceGetTrieNodesQuery(chain SnapshotReader, req *GetTrieNodesPacket, start time.Time) ([][]byte, error) {
	if req.Bytes > softResponseLimit {
		req.Bytes = softResponseLimit
	}
	// Make sure we have the state associated with the request
	triedb := chain.StateCache().TrieDB()

	accTrie, err := statedb.NewSecureTrie(req.Root.ToRootExtHash(), triedb)
	if err != nil {
		// We don't have the requested state available, bail out
		return nil, nil
	}
	snap := chain.Snapshots().Snapshot(req.Root.ToRootExtHash())
	if snap == nil {
		// We don't have the requested state snapshotted yet, bail out.
		// In reality we could still serve using the account and storage
		// tries only, but let's protect the node a bit while it's doing
		// snapshot generation.
		return nil, nil
	}
	// Retrieve trie nodes until the packet size limit is reached
	var (
		nodes [][]byte
		bytes uint64
		loads int // Trie hash expansions to count database reads
	)
	for _, pathset := range req.Paths {
		switch len(pathset) {
		case 0:
			// Ensure we penalize invalid requests
			return nil, fmt.Errorf("%w: zero-item pathset requested", errBadRequest)

		case 1:
			// If we're only retrieving an account trie node, fetch it directly
			blob, resolved, err := accTrie.TryGetNode(pathset[0])
			loads += resolved // always account database reads, even for failures
			if err != nil {
				break
			}
			nodes = append(nodes, blob)
			bytes += uint64(len(blob))

		default:
			// Storage slots requested, open the storage trie and retrieve from there
			acc, err := snap.Account(common.BytesToExtHash(pathset[0]))
			loads++ // always account database reads, even for failures
			if err != nil || acc == nil {
				break
			}
			pacc := account.GetProgramAccount(acc)
			if pacc == nil {
				break
			}
			stTrie, err := statedb.NewSecureTrie(pacc.GetStorageRoot(), triedb)
			loads++ // always account database reads, even for failures
			if err != nil {
				break
			}
			for _, path := range pathset[1:] {
				blob, resolved, err := stTrie.TryGetNode(path)
				loads += resolved // always account database reads, even for failures
				if err != nil {
					break
				}
				nodes = append(nodes, blob)
				bytes += uint64(len(blob))

				// Sanity check limits to avoid DoS on the store trie loads
				if bytes > req.Bytes || loads > maxTrieNodeLookups || time.Since(start) > maxTrieNodeTimeSpent {
					break
				}
			}
		}
		// Abort request processing if we've exceeded our limits
		if bytes > req.Bytes || loads > maxTrieNodeLookups || time.Since(start) > maxTrieNodeTimeSpent {
			break
		}
	}
	return nodes, nil
}

// NodeInfo represents a short summary of the `snap` sub-protocol metadata
// known about the host peer.
type NodeInfo struct{}

// nodeInfo retrieves some `snap` protocol metadata about the running host node.
func nodeInfo(chain *blockchain.BlockChain) *NodeInfo {
	return &NodeInfo{}
}
