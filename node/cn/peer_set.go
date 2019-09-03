// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from eth/peer.go (2018/06/04).
// Modified and improved for the klaytn development.

package cn

import (
	"fmt"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/node"
	"math/big"
	"sync"
)

//go:generate mockgen -destination=node/cn/peer_set_mock.go -package=cn github.com/klaytn/klaytn/node/cn PeerSet
type PeerSet interface {
	Register(p Peer) error
	Unregister(id string) error

	Peers() map[string]Peer
	CNPeers() map[common.Address]Peer
	ENPeers() map[common.Address]Peer
	PNPeers() map[common.Address]Peer
	Peer(id string) Peer
	Len() int

	PeersWithoutBlock(hash common.Hash) []Peer
	TypePeersWithoutBlock(hash common.Hash, nodetype p2p.ConnType) []Peer
	PeersWithoutBlockExceptCN(hash common.Hash) []Peer

	CNWithoutBlock(hash common.Hash) []Peer
	PNWithoutBlock(hash common.Hash) []Peer
	ENWithoutBlock(hash common.Hash) []Peer
	TypePeers(nodetype p2p.ConnType) []Peer

	PeersWithoutTx(hash common.Hash) []Peer
	TypePeersWithoutTx(hash common.Hash, nodetype p2p.ConnType) []Peer
	TypePeersWithTx(hash common.Hash, nodetype p2p.ConnType) []Peer
	AnotherTypePeersWithoutTx(hash common.Hash, nodetype p2p.ConnType) []Peer
	AnotherTypePeersWithTx(hash common.Hash, nodetype p2p.ConnType) []Peer
	CNWithoutTx(hash common.Hash) []Peer

	BestPeer() Peer
	RegisterValidator(connType p2p.ConnType, validator p2p.PeerTypeValidator)
	Close()
}

// peerSet represents the collection of active peers currently participating in
// the Klaytn sub-protocol.
type peerSet struct {
	peers   map[string]Peer
	cnpeers map[common.Address]Peer
	pnpeers map[common.Address]Peer
	enpeers map[common.Address]Peer
	lock    sync.RWMutex
	closed  bool

	validator map[p2p.ConnType]p2p.PeerTypeValidator
}

// newPeerSet creates a new peer set to track the active participants.
func newPeerSet() *peerSet {
	peerSet := &peerSet{
		peers:     make(map[string]Peer),
		cnpeers:   make(map[common.Address]Peer),
		pnpeers:   make(map[common.Address]Peer),
		enpeers:   make(map[common.Address]Peer),
		validator: make(map[p2p.ConnType]p2p.PeerTypeValidator),
	}

	peerSet.validator[node.CONSENSUSNODE] = ByPassValidator{}
	peerSet.validator[node.PROXYNODE] = ByPassValidator{}
	peerSet.validator[node.ENDPOINTNODE] = ByPassValidator{}

	return peerSet
}

// Register injects a new peer into the working set, or returns an error if the
// peer is already known.
func (ps *peerSet) Register(p Peer) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return errClosed
	}
	if _, ok := ps.peers[p.GetID()]; ok {
		return errAlreadyRegistered
	}

	var peersByNodeType map[common.Address]Peer
	var peerTypeValidator p2p.PeerTypeValidator

	switch p.ConnType() {
	case node.CONSENSUSNODE:
		peersByNodeType = ps.cnpeers
		peerTypeValidator = ps.validator[node.CONSENSUSNODE]
	case node.PROXYNODE:
		peersByNodeType = ps.pnpeers
		peerTypeValidator = ps.validator[node.PROXYNODE]
	case node.ENDPOINTNODE:
		peersByNodeType = ps.enpeers
		peerTypeValidator = ps.validator[node.ENDPOINTNODE]
	default:
		return fmt.Errorf("undefined peer type entered, p.ConnType(): %v", p.ConnType())
	}

	if _, ok := peersByNodeType[p.GetAddr()]; ok {
		return errAlreadyRegistered
	}

	if err := peerTypeValidator.ValidatePeerType(p.GetAddr()); err != nil {
		return fmt.Errorf("fail to validate peer type: %s", err)
	}

	peersByNodeType[p.GetAddr()] = p // add peer to its node type peer map.
	ps.peers[p.GetID()] = p          // add peer to entire peer map.

	cnPeerCountGauge.Update(int64(len(ps.cnpeers)))
	pnPeerCountGauge.Update(int64(len(ps.pnpeers)))
	enPeerCountGauge.Update(int64(len(ps.enpeers)))
	go p.Broadcast()

	return nil
}

// Unregister removes a remote peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *peerSet) Unregister(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	p, ok := ps.peers[id]
	if !ok {
		return errNotRegistered
	}
	if p.ConnType() == node.CONSENSUSNODE {
		delete(ps.cnpeers, p.GetAddr())
	} else if p.ConnType() == node.PROXYNODE {
		delete(ps.pnpeers, p.GetAddr())
	} else if p.ConnType() == node.ENDPOINTNODE {
		delete(ps.enpeers, p.GetAddr())
	}
	delete(ps.peers, id)
	p.Close()

	cnPeerCountGauge.Update(int64(len(ps.cnpeers)))
	pnPeerCountGauge.Update(int64(len(ps.pnpeers)))
	enPeerCountGauge.Update(int64(len(ps.enpeers)))
	return nil
}

// istanbul BFT
func (ps *peerSet) Peers() map[string]Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	set := make(map[string]Peer)
	for id, p := range ps.peers {
		set[id] = p
	}
	return set
}

func (ps *peerSet) CNPeers() map[common.Address]Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	set := make(map[common.Address]Peer)
	for addr, p := range ps.cnpeers {
		set[addr] = p
	}
	return set
}

func (ps *peerSet) ENPeers() map[common.Address]Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	set := make(map[common.Address]Peer)
	for addr, p := range ps.enpeers {
		set[addr] = p
	}
	return set
}

func (ps *peerSet) PNPeers() map[common.Address]Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	set := make(map[common.Address]Peer)
	for addr, p := range ps.pnpeers {
		set[addr] = p
	}
	return set
}

// Peer retrieves the registered peer with the given id.
func (ps *peerSet) Peer(id string) Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}

// Len returns if the current number of peers in the set.
func (ps *peerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// PeersWithoutBlock retrieves a list of peers that do not have a given block in
// their set of known hashes.
func (ps *peerSet) PeersWithoutBlock(hash common.Hash) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.KnowsBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) TypePeersWithoutBlock(hash common.Hash, nodetype p2p.ConnType) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.ConnType() == nodetype && !p.KnowsBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PeersWithoutBlockExceptCN(hash common.Hash) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.ConnType() != node.CONSENSUSNODE && !p.KnowsBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) CNWithoutBlock(hash common.Hash) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.cnpeers))
	for _, p := range ps.cnpeers {
		if !p.KnowsBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) PNWithoutBlock(hash common.Hash) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.pnpeers))
	for _, p := range ps.pnpeers {
		if !p.KnowsBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) ENWithoutBlock(hash common.Hash) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.enpeers))
	for _, p := range ps.enpeers {
		if !p.KnowsBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) TypePeers(nodetype p2p.ConnType) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.ConnType() == nodetype {
			list = append(list, p)
		}
	}
	return list
}

// PeersWithoutTx retrieves a list of peers that do not have a given transaction
// in their set of known hashes.
func (ps *peerSet) PeersWithoutTx(hash common.Hash) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.KnowsTx(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) TypePeersWithoutTx(hash common.Hash, nodetype p2p.ConnType) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.ConnType() == nodetype && !p.KnowsTx(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) TypePeersWithTx(hash common.Hash, nodetype p2p.ConnType) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.ConnType() == nodetype && p.KnowsTx(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) AnotherTypePeersWithoutTx(hash common.Hash, nodetype p2p.ConnType) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.ConnType() != nodetype && !p.KnowsTx(hash) {
			list = append(list, p)
		}
	}
	return list
}

// TODO-Klaytn drop or missing tx
func (ps *peerSet) AnotherTypePeersWithTx(hash common.Hash, nodetype p2p.ConnType) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.ConnType() != nodetype && p.KnowsTx(hash) {
			list = append(list, p)
		}
	}
	return list
}

func (ps *peerSet) CNWithoutTx(hash common.Hash) []Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]Peer, 0, len(ps.cnpeers))
	for _, p := range ps.cnpeers {
		if !p.KnowsTx(hash) {
			list = append(list, p)
		}
	}
	return list
}

// BestPeer retrieves the known peer with the currently highest total blockscore.
func (ps *peerSet) BestPeer() Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	var (
		bestPeer Peer
		bestTd   *big.Int
	)
	for _, p := range ps.peers {
		if _, td := p.Head(); bestPeer == nil || td.Cmp(bestTd) > 0 {
			bestPeer, bestTd = p, td
		}
	}
	return bestPeer
}

// RegisterValidator registers a validator.
func (ps *peerSet) RegisterValidator(connType p2p.ConnType, validator p2p.PeerTypeValidator) {
	ps.validator[connType] = validator
}

// Close disconnects all peers.
// No new peers can be registered after Close has returned.
func (ps *peerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.GetP2PPeer().Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
