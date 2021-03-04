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

package sc

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
)

var (
	errClosed            = errors.New("peer set is closed")
	errAlreadyRegistered = errors.New("peer is already registered")
	errNotRegistered     = errors.New("peer is not registered")
)

const (
	maxKnownTxs = 32768 // Maximum transactions hashes to keep in the known list (prevent DOS)

	handshakeTimeout = 5 * time.Second
)

// BridgePeerInfo represents a short summary of the Klaytn Bridge sub-protocol metadata known
// about a connected peer.
type BridgePeerInfo struct {
	Version int    `json:"version"` // Klaytn Bridge protocol version negotiated
	Head    string `json:"head"`    // SHA3 hash of the peer's best owned block
}

type PeerSetManager interface {
	BridgePeerSet() *bridgePeerSet
}

//go:generate mockgen -destination=node/sc/bridgepeer_mock_test.go -package=sc github.com/klaytn/klaytn/node/sc BridgePeer
type BridgePeer interface {
	// Close signals the broadcast goroutine to terminate.
	Close()

	// Info gathers and returns a collection of metadata known about a peer.
	Info() *BridgePeerInfo

	Head() (hash common.Hash, td *big.Int)

	// AddToKnownTxs adds a transaction hash to knownTxsCache for the peer, ensuring that it
	// will never be propagated to this particular peer.
	AddToKnownTxs(hash common.Hash)

	// Send writes an RLP-encoded message with the given code.
	// data should have been encoded as an RLP list.
	Send(msgcode uint64, data interface{}) error

	// Handshake executes the Klaytn protocol handshake, negotiating version number,
	// network IDs, difficulties, head, genesis blocks, and onChildChain(if the node is on child chain for the peer)
	// and returning if the peer on the same chain or not and error.
	Handshake(network uint64, chainID, td *big.Int, head common.Hash) error

	// ConnType returns the conntype of the peer.
	ConnType() common.ConnType

	// GetID returns the id of the peer.
	GetID() string

	// GetP2PPeerID returns the id of the p2p.Peer.
	GetP2PPeerID() discover.NodeID

	// GetChainID returns the chain id of the peer.
	GetChainID() *big.Int

	// GetAddr returns the address of the peer.
	GetAddr() common.Address

	// SetAddr sets the address of the peer.
	SetAddr(addr common.Address)

	// GetVersion returns the version of the peer.
	GetVersion() int

	// KnowsTx returns if the peer is known to have the transaction, based on knownTxsCache.
	KnowsTx(hash common.Hash) bool

	// GetP2PPeer returns the p2p.
	GetP2PPeer() *p2p.Peer

	// GetRW returns the MsgReadWriter of the peer.
	GetRW() p2p.MsgReadWriter

	// Handle is the callback invoked to manage the life cycle of a Klaytn Peer. When
	// this function terminates, the Peer is disconnected.
	Handle(bn *MainBridge) error

	SendRequestRPC(data []byte) error
	SendResponseRPC(data []byte) error

	// SendServiceChainTxs sends child chain tx data to from child chain to parent chain.
	SendServiceChainTxs(txs types.Transactions) error

	// SendServiceChainInfoRequest sends a parentChainInfo request from child chain to parent chain.
	SendServiceChainInfoRequest(addr *common.Address) error

	// SendServiceChainInfoResponse sends a parentChainInfo from parent chain to child chain.
	// parentChainInfo includes nonce of an account and gasPrice in the parent chain.
	SendServiceChainInfoResponse(pcInfo *parentChainInfo) error

	// SendServiceChainReceiptRequest sends a receipt request from child chain to parent chain.
	SendServiceChainReceiptRequest(txHashes []common.Hash) error

	// SendServiceChainReceiptResponse sends a receipt as a response to request from child chain.
	SendServiceChainReceiptResponse(receipts []*types.ReceiptForStorage) error
}

// baseBridgePeer is a common data structure used by implementation of Peer.
type baseBridgePeer struct {
	id string

	addr common.Address

	*p2p.Peer
	rw p2p.MsgReadWriter

	version int // Protocol version negotiated

	head common.Hash
	td   *big.Int
	lock sync.RWMutex

	knownTxsCache common.Cache  // FIFO cache of transaction hashes known to be known by this peer
	term          chan struct{} // Termination channel to stop the broadcaster

	chainID *big.Int // A child chain must know parent chain's ChainID to sign a transaction.

	respCh chan *big.Int
}

// newKnownTxCache returns an empty cache for knownTxsCache.
func newKnownTxCache() common.Cache {
	return common.NewCache(common.FIFOCacheConfig{CacheSize: maxKnownTxs, IsScaled: true})
}

// newPeer returns new Peer interface.
func newBridgePeer(version int, p *p2p.Peer, rw p2p.MsgReadWriter) BridgePeer {
	id := p.ID()

	return &singleChannelPeer{
		baseBridgePeer: &baseBridgePeer{
			Peer:          p,
			rw:            rw,
			version:       version,
			id:            fmt.Sprintf("%x", id[:8]),
			knownTxsCache: newKnownTxCache(),
			term:          make(chan struct{}),
			respCh:        make(chan *big.Int),
		},
	}
}

// Close signals the broadcast goroutine to terminate.
func (p *baseBridgePeer) Close() {
	close(p.term)
}

// Info gathers and returns a collection of metadata known about a peer.
func (p *baseBridgePeer) Info() *BridgePeerInfo {
	hash, _ := p.Head()

	return &BridgePeerInfo{
		Version: p.version,
		Head:    hash.Hex(),
	}
}

// Head retrieves a copy of the current head hash and total blockscore of the
// peer.
func (p *baseBridgePeer) Head() (hash common.Hash, td *big.Int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	copy(hash[:], p.head[:])
	return hash, new(big.Int).Set(p.td)
}

// SetHead updates the head hash and total blockscore of the peer.
func (p *baseBridgePeer) SetHead(hash common.Hash, td *big.Int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.head[:], hash[:])
	p.td.Set(td)
}

// AddToKnownTxs adds a transaction hash to knownTxsCache for the peer, ensuring that it
// will never be propagated to this particular peer.
func (p *baseBridgePeer) AddToKnownTxs(hash common.Hash) {
	p.knownTxsCache.Add(hash, struct{}{})
}

// Send writes an RLP-encoded message with the given code.
// data should have been encoded as an RLP list.
func (p *baseBridgePeer) Send(msgcode uint64, data interface{}) error {
	return p2p.Send(p.rw, msgcode, data)
}

func (p *baseBridgePeer) SendRequestRPC(data []byte) error {
	return p2p.Send(p.rw, ServiceChainCall, data)
}

func (p *baseBridgePeer) SendResponseRPC(data []byte) error {
	return p2p.Send(p.rw, ServiceChainResponse, data)
}

func (p *baseBridgePeer) SendServiceChainTxs(txs types.Transactions) error {
	return p2p.Send(p.rw, ServiceChainTxsMsg, txs)
}

func (p *baseBridgePeer) SendServiceChainInfoRequest(addr *common.Address) error {
	return p2p.Send(p.rw, ServiceChainParentChainInfoRequestMsg, addr)
}

func (p *baseBridgePeer) SendServiceChainInfoResponse(pcInfo *parentChainInfo) error {
	return p2p.Send(p.rw, ServiceChainParentChainInfoResponseMsg, pcInfo)
}

func (p *baseBridgePeer) SendServiceChainReceiptRequest(txHashes []common.Hash) error {
	return p2p.Send(p.rw, ServiceChainReceiptRequestMsg, txHashes)
}

func (p *baseBridgePeer) SendServiceChainReceiptResponse(receipts []*types.ReceiptForStorage) error {
	return p2p.Send(p.rw, ServiceChainReceiptResponseMsg, receipts)
}

// Handshake executes the Klaytn protocol handshake, negotiating version number,
// network IDs, difficulties, head and genesis blocks.
func (p *baseBridgePeer) Handshake(network uint64, chainID, td *big.Int, head common.Hash) error {
	// Send out own handshake in a new thread
	errc := make(chan error, 2)
	var status statusData // safe to read after two values have been received from errc

	go func() {
		errc <- p2p.Send(p.rw, StatusMsg, &statusData{
			ProtocolVersion: uint32(p.version),
			NetworkId:       network,
			TD:              td,
			CurrentBlock:    head,
			ChainID:         chainID,
		})
	}()
	go func() {
		e := p.readStatus(network, &status)
		if e != nil {
			errc <- e
			return
		}
		errc <- e
	}()
	timeout := time.NewTimer(handshakeTimeout)
	defer timeout.Stop()
	for i := 0; i < 2; i++ {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return p2p.DiscReadTimeout
		}
	}
	p.td, p.head, p.chainID = status.TD, status.CurrentBlock, status.ChainID
	return nil
}

func (p *baseBridgePeer) readStatus(network uint64, status *statusData) error {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Code != StatusMsg {
		return errResp(ErrNoStatusMsg, "first msg has code %x (!= %x)", msg.Code, StatusMsg)
	}
	if msg.Size > ProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	}
	// Decode the handshake and make sure everything matches
	if err := msg.Decode(&status); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	if status.NetworkId != network {
		return errResp(ErrNetworkIdMismatch, "%d (!= %d)", status.NetworkId, network)
	}
	if int(status.ProtocolVersion) != p.version {
		return errResp(ErrProtocolVersionMismatch, "%d (!= %d)", status.ProtocolVersion, p.version)
	}
	return nil
}

// String implements fmt.Stringer.
func (p *baseBridgePeer) String() string {
	return fmt.Sprintf("Peer %s [%s]", p.id,
		fmt.Sprintf("klay/%2d", p.version),
	)
}

// ConnType returns the conntype of the peer.
func (p *baseBridgePeer) ConnType() common.ConnType {
	return p.Peer.ConnType()
}

// GetID returns the id of the peer.
func (p *baseBridgePeer) GetID() string {
	return p.id
}

// GetP2PPeerID returns the id of the p2p.Peer.
func (p *baseBridgePeer) GetP2PPeerID() discover.NodeID {
	return p.Peer.ID()
}

// GetChainID returns the chain id of the peer.
func (p *baseBridgePeer) GetChainID() *big.Int {
	return p.chainID
}

// GetAddr returns the address of the peer.
func (p *baseBridgePeer) GetAddr() common.Address {
	return p.addr
}

// SetAddr sets the address of the peer.
func (p *baseBridgePeer) SetAddr(addr common.Address) {
	p.addr = addr
}

// GetVersion returns the version of the peer.
func (p *baseBridgePeer) GetVersion() int {
	return p.version
}

// KnowsTx returns if the peer is known to have the transaction, based on knownTxsCache.
func (p *baseBridgePeer) KnowsTx(hash common.Hash) bool {
	_, ok := p.knownTxsCache.Get(hash)
	return ok
}

// GetP2PPeer returns the p2p.Peer.
func (p *baseBridgePeer) GetP2PPeer() *p2p.Peer {
	return p.Peer
}

// GetRW returns the MsgReadWriter of the peer.
func (p *baseBridgePeer) GetRW() p2p.MsgReadWriter {
	return p.rw
}

// Handle is the callback invoked to manage the life cycle of a Klaytn Peer. When
// this function terminates, the Peer is disconnected.
func (p *baseBridgePeer) Handle(bn *MainBridge) error {
	return bn.handle(p)
}

// singleChannelPeer is a peer that uses a single channel.
type singleChannelPeer struct {
	*baseBridgePeer
}

// bridgePeerSet represents the collection of active peers currently participating in
// the Klaytn sub-protocol.
type bridgePeerSet struct {
	peers  map[string]BridgePeer
	lock   sync.RWMutex
	closed bool
}

// newBridgePeerSet creates a new peer set to track the active participants.
func newBridgePeerSet() *bridgePeerSet {
	peerSet := &bridgePeerSet{
		peers: make(map[string]BridgePeer),
	}

	return peerSet
}

// Register injects a new peer into the working set, or returns an error if the
// peer is already known.
func (ps *bridgePeerSet) Register(p BridgePeer) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return errClosed
	}
	if _, ok := ps.peers[p.GetID()]; ok {
		return errAlreadyRegistered
	}
	ps.peers[p.GetID()] = p

	return nil
}

// Unregister removes a remote peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *bridgePeerSet) Unregister(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	p, ok := ps.peers[id]
	if !ok {
		return errNotRegistered
	}
	delete(ps.peers, id)
	p.Close()

	return nil
}

// istanbul BFT
func (ps *bridgePeerSet) Peers() map[string]BridgePeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	set := make(map[string]BridgePeer)
	for id, p := range ps.peers {
		set[id] = p
	}
	return set
}

// Peer retrieves the registered peer with the given id.
func (ps *bridgePeerSet) Peer(id string) BridgePeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}

// Len returns if the current number of peers in the set.
func (ps *bridgePeerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// PeersWithoutTx retrieves a list of peers that do not have a given transaction
// in their set of known hashes.
func (ps *bridgePeerSet) PeersWithoutTx(hash common.Hash) []BridgePeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]BridgePeer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.KnowsTx(hash) {
			list = append(list, p)
		}
	}
	return list
}

// BestPeer retrieves the known peer with the currently highest total blockscore.
func (ps *bridgePeerSet) BestPeer() BridgePeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	var (
		bestPeer BridgePeer
		bestTd   *big.Int
	)
	for _, p := range ps.peers {
		if _, td := p.Head(); bestPeer == nil || td.Cmp(bestTd) > 0 {
			bestPeer, bestTd = p, td
		}
	}
	return bestPeer
}

// Close disconnects all peers.
// No new peers can be registered after Close has returned.
func (ps *bridgePeerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.GetP2PPeer().Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
