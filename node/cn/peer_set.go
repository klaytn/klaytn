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
	"math/big"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
)

//go:generate mockgen -destination=node/cn/peer_set_mock_test.go -package=cn github.com/klaytn/klaytn/node/cn PeerSet
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

	SamplePeersToSendBlock(block *types.Block, nodeType common.ConnType) []Peer
	SampleResendPeersByType(nodeType common.ConnType) []Peer

	PeersWithoutTx(hash common.Hash) []Peer
	TypePeersWithoutTx(hash common.Hash, nodetype common.ConnType) []Peer
	CNWithoutTx(hash common.Hash) []Peer
	UpdateTypePeersWithoutTxs(tx *types.Transaction, nodeType common.ConnType, peersWithoutTxsMap map[Peer]types.Transactions)

	BestPeer() Peer
	RegisterValidator(connType common.ConnType, validator p2p.PeerTypeValidator)
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

	validator map[common.ConnType]p2p.PeerTypeValidator
}

// newPeerSet creates a new peer set to track the active participants.
func newPeerSet() *peerSet {
	peerSet := &peerSet{
		peers:     make(map[string]Peer),
		cnpeers:   make(map[common.Address]Peer),
		pnpeers:   make(map[common.Address]Peer),
		enpeers:   make(map[common.Address]Peer),
		validator: make(map[common.ConnType]p2p.PeerTypeValidator),
	}

	peerSet.validator[common.CONSENSUSNODE] = ByPassValidator{}
	peerSet.validator[common.PROXYNODE] = ByPassValidator{}
	peerSet.validator[common.ENDPOINTNODE] = ByPassValidator{}

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
	case common.CONSENSUSNODE:
		peersByNodeType = ps.cnpeers
		peerTypeValidator = ps.validator[common.CONSENSUSNODE]
	case common.PROXYNODE:
		peersByNodeType = ps.pnpeers
		peerTypeValidator = ps.validator[common.PROXYNODE]
	case common.ENDPOINTNODE:
		peersByNodeType = ps.enpeers
		peerTypeValidator = ps.validator[common.ENDPOINTNODE]
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
	delete(ps.peers, id)
	p.Close()

	switch p.ConnType() {
	case common.CONSENSUSNODE:
		delete(ps.cnpeers, p.GetAddr())
	case common.PROXYNODE:
		delete(ps.pnpeers, p.GetAddr())
	case common.ENDPOINTNODE:
		delete(ps.enpeers, p.GetAddr())
	default:
		return errUnexpectedNodeType
	}

	cnPeerCountGauge.Update(int64(len(ps.cnpeers)))
	pnPeerCountGauge.Update(int64(len(ps.pnpeers)))
	enPeerCountGauge.Update(int64(len(ps.enpeers)))
	return nil
}

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

func (ps *peerSet) typePeersWithoutBlock(hash common.Hash, nodetype common.ConnType) []Peer {
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
		if p.ConnType() != common.CONSENSUSNODE && !p.KnowsBlock(hash) {
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

func (ps *peerSet) typePeers(nodetype common.ConnType) []Peer {
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

func (ps *peerSet) TypePeersWithoutTx(hash common.Hash, nodetype common.ConnType) []Peer {
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
		bestPeer       Peer
		bestBlockScore *big.Int
	)
	for _, p := range ps.peers {
		if _, currBlockScore := p.Head(); bestPeer == nil || currBlockScore.Cmp(bestBlockScore) > 0 {
			bestPeer, bestBlockScore = p, currBlockScore
		}
	}
	return bestPeer
}

// RegisterValidator registers a validator.
func (ps *peerSet) RegisterValidator(connType common.ConnType, validator p2p.PeerTypeValidator) {
	ps.validator[connType] = validator
}

// Close disconnects all peers.
// No new peers can be registered after Close has returned.
func (ps *peerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.DisconnectP2PPeer(p2p.DiscQuitting)
	}
	ps.closed = true
}

// samplePeersToSendBlock samples peers from peers without block.
// It uses different sampling policy for different node type.
func (peers *peerSet) SamplePeersToSendBlock(block *types.Block, nodeType common.ConnType) []Peer {
	var peersWithoutBlock []Peer
	hash := block.Hash()

	switch nodeType {
	case common.CONSENSUSNODE:
		// If currNode is CN, sends block to sampled peers from (CN + PN), not to EN.
		cnsWithoutBlock := peers.CNWithoutBlock(hash)
		sampledCNsWithoutBlock := samplingPeers(cnsWithoutBlock, sampleSize(cnsWithoutBlock))

		// CN always broadcasts a block to its PN peers, unless the number of PN peers exceeds the limit.
		pnsWithoutBlock := peers.PNWithoutBlock(hash)
		if len(pnsWithoutBlock) > blockReceivingPNLimit {
			pnsWithoutBlock = samplingPeers(pnsWithoutBlock, blockReceivingPNLimit)
		}

		logger.Trace("Propagated block", "hash", hash,
			"CN recipients", len(sampledCNsWithoutBlock), "PN recipients", len(pnsWithoutBlock), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))

		return append(cnsWithoutBlock, pnsWithoutBlock...)
	case common.PROXYNODE:
		// If currNode is PN, sends block to sampled peers from (PN + EN), not to CN.
		peersWithoutBlock = peers.PeersWithoutBlockExceptCN(hash)

	case common.ENDPOINTNODE:
		// If currNode is EN, sends block to sampled EN peers, not to EN nor CN.
		peersWithoutBlock = peers.ENWithoutBlock(hash)

	default:
		logger.Error("Undefined nodeType of protocolManager! nodeType: %v", nodeType)
		return []Peer{}
	}

	sampledPeersWithoutBlock := samplingPeers(peersWithoutBlock, sampleSize(peersWithoutBlock))
	logger.Trace("Propagated block", "hash", hash,
		"recipients", len(sampledPeersWithoutBlock), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))

	return sampledPeersWithoutBlock
}

func (peers *peerSet) SampleResendPeersByType(nodeType common.ConnType) []Peer {
	// TODO-Klaytn Need to tune pickSize. Currently use 2 for availability and efficiency.
	var sampledPeers []Peer
	switch nodeType {
	case common.ENDPOINTNODE:
		sampledPeers = peers.typePeers(common.CONSENSUSNODE)
		if len(sampledPeers) < 2 {
			sampledPeers = append(sampledPeers, samplingPeers(peers.typePeers(common.PROXYNODE), 2-len(sampledPeers))...)
		}
		if len(sampledPeers) < 2 {
			sampledPeers = append(sampledPeers, samplingPeers(peers.typePeers(common.ENDPOINTNODE), 2-len(sampledPeers))...)
		}
		sampledPeers = samplingPeers(sampledPeers, 2)
	case common.PROXYNODE:
		sampledPeers = peers.typePeers(common.CONSENSUSNODE)
		if len(sampledPeers) == 0 {
			sampledPeers = peers.typePeers(common.PROXYNODE)
		}
		sampledPeers = samplingPeers(sampledPeers, 2)
	default:
		logger.Warn("Not supported nodeType", "nodeType", nodeType)
		return nil
	}
	return sampledPeers
}

func (peers *peerSet) UpdateTypePeersWithoutTxs(tx *types.Transaction, nodeType common.ConnType, peersWithoutTxsMap map[Peer]types.Transactions) {
	typePeers := peers.TypePeersWithoutTx(tx.Hash(), nodeType)
	for _, peer := range typePeers {
		peersWithoutTxsMap[peer] = append(peersWithoutTxsMap[peer], tx)
	}
	logger.Trace("Broadcast transaction", "hash", tx.Hash(), "recipients", len(typePeers))
}
