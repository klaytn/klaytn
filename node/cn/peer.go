// Modifications Copyright 2018 The klaytn Authors
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
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/rlp"
)

var (
	errClosed                 = errors.New("peer set is closed")
	errAlreadyRegistered      = errors.New("peer is already registered")
	errNotRegistered          = errors.New("peer is not registered")
	errUnexpectedNodeType     = errors.New("unexpected node type of peer")
	errNotSupportedByBasePeer = errors.New("not supported by basePeer")
)

const (
	maxKnownTxs    = 32768 // Maximum transactions hashes to keep in the known list (prevent DOS)
	maxKnownBlocks = 1024  // Maximum block hashes to keep in the known list (prevent DOS)

	// maxQueuedTxs is the maximum number of transaction lists to queue up before
	// dropping broadcasts. This is a sensitive number as a transaction list might
	// contain a single transaction, or thousands.
	maxQueuedTxs = 128

	// maxQueuedProps is the maximum number of block propagations to queue up before
	// dropping broadcasts. There's not much point in queueing stale blocks, so a few
	// that might cover uncles should be enough.
	// TODO-Klaytn-Refactoring Look into the usage of maxQueuedProps and remove it if needed
	maxQueuedProps = 4

	// maxQueuedAnns is the maximum number of block announcements to queue up before
	// dropping broadcasts. Similarly to block propagations, there's no point to queue
	// above some healthy uncle limit, so use that.
	// TODO-Klaytn-Refactoring Look into the usage of maxQueuedAnns and remove it if needed
	maxQueuedAnns = 4

	handshakeTimeout = 5 * time.Second
)

// PeerInfo represents a short summary of the Klaytn sub-protocol metadata known
// about a connected peer.
type PeerInfo struct {
	Version    int      `json:"version"`    // Klaytn protocol version negotiated
	BlockScore *big.Int `json:"blockscore"` // Total blockscore of the peer's blockchain
	Head       string   `json:"head"`       // SHA3 hash of the peer's best owned block
}

// propEvent is a block propagation, waiting for its turn in the broadcast queue.
type propEvent struct {
	block *types.Block
	td    *big.Int
}

//go:generate mockgen -destination=node/cn/peer_mock_test.go -package=cn github.com/klaytn/klaytn/node/cn Peer
type Peer interface {
	// Broadcast is a write loop that multiplexes block propagations, announcements
	// and transaction broadcasts into the remote peer. The goal is to have an async
	// writer that does not lock up node internals.
	Broadcast()

	// Close signals the broadcast goroutine to terminate.
	Close()

	// Info gathers and returns a collection of metadata known about a peer.
	Info() *PeerInfo

	// SetHead updates the head hash and total blockscore of the peer.
	SetHead(hash common.Hash, td *big.Int)

	// AddToKnownBlocks adds a block hash to knownBlocksCache for the peer, ensuring that the block will
	// never be propagated to this particular peer.
	AddToKnownBlocks(hash common.Hash)

	// AddToKnownTxs adds a transaction hash to knownTxsCache for the peer, ensuring that it
	// will never be propagated to this particular peer.
	AddToKnownTxs(hash common.Hash)

	// Send writes an RLP-encoded message with the given code.
	// data should have been encoded as an RLP list.
	Send(msgcode uint64, data interface{}) error

	// SendTransactions sends transactions to the peer and includes the hashes
	// in its transaction hash set for future reference.
	SendTransactions(txs types.Transactions) error

	// ReSendTransactions sends txs to a peer in order to prevent the txs from missing.
	ReSendTransactions(txs types.Transactions) error

	// AsyncSendTransactions sends transactions asynchronously to the peer.
	AsyncSendTransactions(txs types.Transactions)

	// SendNewBlockHashes announces the availability of a number of blocks through
	// a hash notification.
	SendNewBlockHashes(hashes []common.Hash, numbers []uint64) error

	// AsyncSendNewBlockHash queues the availability of a block for propagation to a
	// remote peer. If the peer's broadcast queue is full, the event is silently
	// dropped.
	AsyncSendNewBlockHash(block *types.Block)

	// SendNewBlock propagates an entire block to a remote peer.
	SendNewBlock(block *types.Block, td *big.Int) error

	// AsyncSendNewBlock queues an entire block for propagation to a remote peer. If
	// the peer's broadcast queue is full, the event is silently dropped.
	AsyncSendNewBlock(block *types.Block, td *big.Int)

	// SendBlockHeaders sends a batch of block headers to the remote peer.
	SendBlockHeaders(headers []*types.Header) error

	// SendFetchedBlockHeader sends a block header to the remote peer, requested by fetcher.
	SendFetchedBlockHeader(header *types.Header) error

	// SendBlockBodies sends a batch of block contents to the remote peer.
	SendBlockBodies(bodies []*blockBody) error

	// SendBlockBodiesRLP sends a batch of block contents to the remote peer from
	// an already RLP encoded format.
	SendBlockBodiesRLP(bodies []rlp.RawValue) error

	// SendFetchedBlockBodiesRLP sends a batch of block contents to the remote peer from
	// an already RLP encoded format, requested by fetcher.
	SendFetchedBlockBodiesRLP(bodies []rlp.RawValue) error

	// SendNodeDataRLP sends a batch of arbitrary internal data, corresponding to the
	// hashes requested.
	SendNodeData(data [][]byte) error

	// SendReceiptsRLP sends a batch of transaction receipts, corresponding to the
	// ones requested from an already RLP encoded format.
	SendReceiptsRLP(receipts []rlp.RawValue) error

	// FetchBlockHeader is a wrapper around the header query functions to fetch a
	// single header. It is used solely by the fetcher.
	FetchBlockHeader(hash common.Hash) error

	// FetchBlockBodies fetches a batch of blocks' bodies corresponding to the hashes
	// specified. If uses different message type from RequestBodies.
	// It is used solely by the fetcher.
	FetchBlockBodies(hashes []common.Hash) error

	// Handshake executes the Klaytn protocol handshake, negotiating version number,
	// network IDs, difficulties, head, and genesis blocks and returning error.
	Handshake(network uint64, chainID, td *big.Int, head common.Hash, genesis common.Hash) error

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

	// KnowsBlock returns if the peer is known to have the block, based on knownBlocksCache.
	KnowsBlock(hash common.Hash) bool

	// KnowsTx returns if the peer is known to have the transaction, based on knownTxsCache.
	KnowsTx(hash common.Hash) bool

	// GetP2PPeer returns the p2p.
	GetP2PPeer() *p2p.Peer

	// DisconnectP2PPeer disconnects the p2p peer with the given reason.
	DisconnectP2PPeer(discReason p2p.DiscReason)

	// GetRW returns the MsgReadWriter of the peer.
	GetRW() p2p.MsgReadWriter

	// Handle is the callback invoked to manage the life cycle of a Klaytn Peer. When
	// this function terminates, the Peer is disconnected.
	Handle(pm *ProtocolManager) error

	// UpdateRWImplementationVersion updates the version of the implementation of RW.
	UpdateRWImplementationVersion()

	// Peer encapsulates the methods required to synchronise with a remote full peer.
	downloader.Peer

	// RegisterConsensusMsgCode registers the channel of consensus msg.
	RegisterConsensusMsgCode(msgCode uint64) error
}

// basePeer is a common data structure used by implementation of Peer.
type basePeer struct {
	id string

	addr common.Address

	*p2p.Peer
	rw p2p.MsgReadWriter

	version  int         // Protocol version negotiated
	forkDrop *time.Timer // Timed connection dropper if forks aren't validated in time

	head common.Hash
	td   *big.Int
	lock sync.RWMutex

	knownTxsCache    common.Cache              // FIFO cache of transaction hashes known to be known by this peer
	knownBlocksCache common.Cache              // FIFO cache of block hashes known to be known by this peer
	queuedTxs        chan []*types.Transaction // Queue of transactions to broadcast to the peer
	queuedProps      chan *propEvent           // Queue of blocks to broadcast to the peer
	queuedAnns       chan *types.Block         // Queue of blocks to announce to the peer
	term             chan struct{}             // Termination channel to stop the broadcaster

	chainID *big.Int // ChainID to sign a transaction
}

// newKnownBlockCache returns an empty cache for knownBlocksCache.
func newKnownBlockCache() common.Cache {
	return common.NewCache(common.FIFOCacheConfig{CacheSize: maxKnownBlocks, IsScaled: true})
}

// newKnownTxCache returns an empty cache for knownTxsCache.
func newKnownTxCache() common.Cache {
	return common.NewCache(common.FIFOCacheConfig{CacheSize: maxKnownTxs, IsScaled: true})
}

// newPeer returns new Peer interface.
func newPeer(version int, p *p2p.Peer, rw p2p.MsgReadWriter) Peer {
	id := p.ID()

	return &singleChannelPeer{
		basePeer: &basePeer{
			Peer:             p,
			rw:               rw,
			version:          version,
			id:               fmt.Sprintf("%x", id[:8]),
			knownTxsCache:    newKnownTxCache(),
			knownBlocksCache: newKnownBlockCache(),
			queuedTxs:        make(chan []*types.Transaction, maxQueuedTxs),
			queuedProps:      make(chan *propEvent, maxQueuedProps),
			queuedAnns:       make(chan *types.Block, maxQueuedAnns),
			term:             make(chan struct{}),
		},
	}
}

// ChannelOfMessage is a map with the index of the channel per message
var ChannelOfMessage = map[uint64]int{
	StatusMsg:                   p2p.ConnDefault, //StatusMsg's Channel should to be set ConnDefault
	NewBlockHashesMsg:           p2p.ConnDefault,
	BlockHeaderFetchRequestMsg:  p2p.ConnDefault,
	BlockHeaderFetchResponseMsg: p2p.ConnDefault,
	BlockBodiesFetchRequestMsg:  p2p.ConnDefault,
	BlockBodiesFetchResponseMsg: p2p.ConnDefault,
	TxMsg:                       p2p.ConnTxMsg,
	BlockHeadersRequestMsg:      p2p.ConnDefault,
	BlockHeadersMsg:             p2p.ConnDefault,
	BlockBodiesRequestMsg:       p2p.ConnDefault,
	BlockBodiesMsg:              p2p.ConnDefault,
	NewBlockMsg:                 p2p.ConnDefault,

	// Protocol messages belonging to klay/63
	NodeDataRequestMsg: p2p.ConnDefault,
	NodeDataMsg:        p2p.ConnDefault,
	ReceiptsRequestMsg: p2p.ConnDefault,
	ReceiptsMsg:        p2p.ConnDefault,
}

var ConcurrentOfChannel = []int{
	p2p.ConnDefault: 1,
	p2p.ConnTxMsg:   3,
}

// newPeerWithRWs creates a new Peer object with a slice of p2p.MsgReadWriter.
func newPeerWithRWs(version int, p *p2p.Peer, rws []p2p.MsgReadWriter) (Peer, error) {
	id := p.ID()

	lenRWs := len(rws)
	if lenRWs == 1 {
		return newPeer(version, p, rws[p2p.ConnDefault]), nil
	} else if lenRWs > 1 {
		bPeer := &basePeer{
			Peer:             p,
			rw:               rws[p2p.ConnDefault],
			version:          version,
			id:               fmt.Sprintf("%x", id[:8]),
			knownTxsCache:    newKnownTxCache(),
			knownBlocksCache: newKnownBlockCache(),
			queuedTxs:        make(chan []*types.Transaction, maxQueuedTxs),
			queuedProps:      make(chan *propEvent, maxQueuedProps),
			queuedAnns:       make(chan *types.Block, maxQueuedAnns),
			term:             make(chan struct{}),
		}
		return &multiChannelPeer{
			basePeer: bPeer,
			rws:      rws,
			chMgr:    NewChannelManager(len(rws)),
		}, nil
	} else {
		return nil, errors.New("len(rws) should be greater than zero.")
	}
}

// Broadcast is a write loop that multiplexes block propagations, announcements
// and transaction broadcasts into the remote peer. The goal is to have an async
// writer that does not lock up node internals.
func (p *basePeer) Broadcast() {
	for {
		select {
		case txs := <-p.queuedTxs:
			if err := p.SendTransactions(txs); err != nil {
				logger.Error("fail to SendTransactions", "peer", p.id, "err", err)
				continue
				//return
			}
			p.Log().Trace("Broadcast transactions", "peer", p.id, "count", len(txs))

		case prop := <-p.queuedProps:
			if err := p.SendNewBlock(prop.block, prop.td); err != nil {
				logger.Error("fail to SendNewBlock", "peer", p.id, "err", err)
				continue
				//return
			}
			p.Log().Trace("Propagated block", "peer", p.id, "number", prop.block.Number(), "hash", prop.block.Hash(), "td", prop.td)

		case block := <-p.queuedAnns:
			if err := p.SendNewBlockHashes([]common.Hash{block.Hash()}, []uint64{block.NumberU64()}); err != nil {
				logger.Error("fail to SendNewBlockHashes", "peer", p.id, "err", err)
				continue
				//return
			}
			p.Log().Trace("Announced block", "peer", p.id, "number", block.Number(), "hash", block.Hash())

		case <-p.term:
			p.Log().Debug("Peer broadcast loop end", "peer", p.id)
			return
		}
	}
}

// Close signals the broadcast goroutine to terminate.
func (p *basePeer) Close() {
	close(p.term)
}

// Info gathers and returns a collection of metadata known about a peer.
func (p *basePeer) Info() *PeerInfo {
	hash, td := p.Head()

	return &PeerInfo{
		Version:    p.version,
		BlockScore: td,
		Head:       hash.Hex(),
	}
}

// Head retrieves a copy of the current head hash and total blockscore of the
// peer.
func (p *basePeer) Head() (hash common.Hash, td *big.Int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	copy(hash[:], p.head[:])
	return hash, new(big.Int).Set(p.td)
}

// SetHead updates the head hash and total blockscore of the peer.
func (p *basePeer) SetHead(hash common.Hash, td *big.Int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.head[:], hash[:])
	p.td.Set(td)
}

// AddToKnownBlocks adds a block hash to knownBlocksCache for the peer, ensuring that the block will
// never be propagated to this particular peer.
func (p *basePeer) AddToKnownBlocks(hash common.Hash) {
	p.knownBlocksCache.Add(hash, struct{}{})
}

// AddToKnownTxs adds a transaction hash to knownTxsCache for the peer, ensuring that it
// will never be propagated to this particular peer.
func (p *basePeer) AddToKnownTxs(hash common.Hash) {
	p.knownTxsCache.Add(hash, struct{}{})
}

// Send writes an RLP-encoded message with the given code.
// data should have been encoded as an RLP list.
func (p *basePeer) Send(msgcode uint64, data interface{}) error {
	return p2p.Send(p.rw, msgcode, data)
}

// SendTransactions sends transactions to the peer and includes the hashes
// in its transaction hash set for future reference.
func (p *basePeer) SendTransactions(txs types.Transactions) error {
	// Before sending transactions, sort transactions in ascending order by time.
	if !sort.IsSorted(types.TxByPriceAndTime(txs)) {
		sort.Sort(types.TxByPriceAndTime(txs))
	}

	for _, tx := range txs {
		p.AddToKnownTxs(tx.Hash())
	}
	return p2p.Send(p.rw, TxMsg, txs)
}

// ReSendTransactions sends txs to a peer in order to prevent the txs from missing.
func (p *basePeer) ReSendTransactions(txs types.Transactions) error {
	// Before sending transactions, sort transactions in ascending order by time.
	if !sort.IsSorted(types.TxByPriceAndTime(txs)) {
		sort.Sort(types.TxByPriceAndTime(txs))
	}

	return p2p.Send(p.rw, TxMsg, txs)
}

func (p *basePeer) AsyncSendTransactions(txs types.Transactions) {
	select {
	case p.queuedTxs <- txs:
		for _, tx := range txs {
			p.AddToKnownTxs(tx.Hash())
		}
	default:
		p.Log().Trace("Dropping transaction propagation", "count", len(txs))
	}
}

// SendNewBlockHashes announces the availability of a number of blocks through
// a hash notification.
func (p *basePeer) SendNewBlockHashes(hashes []common.Hash, numbers []uint64) error {
	for _, hash := range hashes {
		p.AddToKnownBlocks(hash)
	}
	request := make(newBlockHashesData, len(hashes))
	for i := 0; i < len(hashes); i++ {
		request[i].Hash = hashes[i]
		request[i].Number = numbers[i]
	}
	return p2p.Send(p.rw, NewBlockHashesMsg, request)
}

// AsyncSendNewBlockHash queues the availability of a block for propagation to a
// remote peer. If the peer's broadcast queue is full, the event is silently
// dropped.
func (p *basePeer) AsyncSendNewBlockHash(block *types.Block) {
	select {
	case p.queuedAnns <- block:
		p.AddToKnownBlocks(block.Hash())
	default:
		p.Log().Debug("Dropping block announcement", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// SendNewBlock propagates an entire block to a remote peer.
func (p *basePeer) SendNewBlock(block *types.Block, td *big.Int) error {
	p.AddToKnownBlocks(block.Hash())
	return p2p.Send(p.rw, NewBlockMsg, []interface{}{block, td})
}

// AsyncSendNewBlock queues an entire block for propagation to a remote peer. If
// the peer's broadcast queue is full, the event is silently dropped.
func (p *basePeer) AsyncSendNewBlock(block *types.Block, td *big.Int) {
	select {
	case p.queuedProps <- &propEvent{block: block, td: td}:
		p.AddToKnownBlocks(block.Hash())
	default:
		p.Log().Debug("Dropping block propagation", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// SendBlockHeaders sends a batch of block headers to the remote peer.
func (p *basePeer) SendBlockHeaders(headers []*types.Header) error {
	return p2p.Send(p.rw, BlockHeadersMsg, headers)
}

// SendFetchedBlockHeader sends a block header to the remote peer, requested by fetcher.
func (p *basePeer) SendFetchedBlockHeader(header *types.Header) error {
	return p2p.Send(p.rw, BlockHeaderFetchResponseMsg, header)
}

// SendBlockBodies sends a batch of block contents to the remote peer.
func (p *basePeer) SendBlockBodies(bodies []*blockBody) error {
	return p2p.Send(p.rw, BlockBodiesMsg, blockBodiesData(bodies))
}

// SendBlockBodiesRLP sends a batch of block contents to the remote peer from
// an already RLP encoded format.
func (p *basePeer) SendBlockBodiesRLP(bodies []rlp.RawValue) error {
	return p2p.Send(p.rw, BlockBodiesMsg, bodies)
}

// SendFetchedBlockBodiesRLP sends a batch of block contents to the remote peer from
// an already RLP encoded format.
func (p *basePeer) SendFetchedBlockBodiesRLP(bodies []rlp.RawValue) error {
	return p2p.Send(p.rw, BlockBodiesFetchResponseMsg, bodies)
}

// SendNodeDataRLP sends a batch of arbitrary internal data, corresponding to the
// hashes requested.
func (p *basePeer) SendNodeData(data [][]byte) error {
	return p2p.Send(p.rw, NodeDataMsg, data)
}

// SendReceiptsRLP sends a batch of transaction receipts, corresponding to the
// ones requested from an already RLP encoded format.
func (p *basePeer) SendReceiptsRLP(receipts []rlp.RawValue) error {
	return p2p.Send(p.rw, ReceiptsMsg, receipts)
}

// FetchBlockHeader is a wrapper around the header query functions to fetch a
// single header. It is used solely by the fetcher.
func (p *basePeer) FetchBlockHeader(hash common.Hash) error {
	p.Log().Debug("Fetching a new block header", "hash", hash)
	return p2p.Send(p.rw, BlockHeaderFetchRequestMsg, hash)
}

// RequestHeadersByHash fetches a batch of blocks' headers corresponding to the
// specified header query, based on the hash of an origin block.
func (p *basePeer) RequestHeadersByHash(origin common.Hash, amount int, skip int, reverse bool) error {
	p.Log().Debug("Fetching batch of headers", "count", amount, "fromhash", origin, "skip", skip, "reverse", reverse)
	return p2p.Send(p.rw, BlockHeadersRequestMsg, &getBlockHeadersData{Origin: hashOrNumber{Hash: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestHeadersByNumber fetches a batch of blocks' headers corresponding to the
// specified header query, based on the number of an origin block.
func (p *basePeer) RequestHeadersByNumber(origin uint64, amount int, skip int, reverse bool) error {
	p.Log().Debug("Fetching batch of headers", "count", amount, "fromnum", origin, "skip", skip, "reverse", reverse)
	return p2p.Send(p.rw, BlockHeadersRequestMsg, &getBlockHeadersData{Origin: hashOrNumber{Number: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestBodies fetches a batch of blocks' bodies corresponding to the hashes
// specified.
func (p *basePeer) RequestBodies(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of block bodies", "count", len(hashes))
	return p2p.Send(p.rw, BlockBodiesRequestMsg, hashes)
}

// FetchBlockBodies fetches a batch of blocks' bodies corresponding to the hashes
// specified. If uses different message type from RequestBodies.
func (p *basePeer) FetchBlockBodies(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of new block bodies", "count", len(hashes))
	return p2p.Send(p.rw, BlockBodiesFetchRequestMsg, hashes)
}

// RequestNodeData fetches a batch of arbitrary data from a node's known state
// data, corresponding to the specified hashes.
func (p *basePeer) RequestNodeData(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of state data", "count", len(hashes))
	return p2p.Send(p.rw, NodeDataRequestMsg, hashes)
}

// RequestReceipts fetches a batch of transaction receipts from a remote node.
func (p *basePeer) RequestReceipts(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of receipts", "count", len(hashes))
	return p2p.Send(p.rw, ReceiptsRequestMsg, hashes)
}

// Handshake executes the Klaytn protocol handshake, negotiating version number,
// network IDs, difficulties, head and genesis blocks.
func (p *basePeer) Handshake(network uint64, chainID, td *big.Int, head common.Hash, genesis common.Hash) error {
	// Send out own handshake in a new thread
	errc := make(chan error, 2)
	var status statusData // safe to read after two values have been received from errc

	go func() {
		errc <- p2p.Send(p.rw, StatusMsg, &statusData{
			ProtocolVersion: uint32(p.version),
			NetworkId:       network,
			TD:              td,
			CurrentBlock:    head,
			GenesisBlock:    genesis,
			ChainID:         chainID,
		})
	}()
	go func() {
		errc <- p.readStatus(network, &status, genesis, chainID)
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

func (p *basePeer) readStatus(network uint64, status *statusData, genesis common.Hash, chainID *big.Int) error {
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
	if status.GenesisBlock != genesis {
		return errResp(ErrGenesisBlockMismatch, "%x (!= %x)", status.GenesisBlock[:8], genesis[:8])
	}
	if status.NetworkId != network {
		return errResp(ErrNetworkIdMismatch, "%d (!= %d)", status.NetworkId, network)
	}
	if status.ChainID.Cmp(chainID) != 0 {
		return errResp(ErrChainIDMismatch, "%v (!= %v)", status.ChainID.String(), chainID.String())
	}
	if int(status.ProtocolVersion) != p.version {
		return errResp(ErrProtocolVersionMismatch, "%d (!= %d)", status.ProtocolVersion, p.version)
	}
	return nil
}

// String implements fmt.Stringer.
func (p *basePeer) String() string {
	return fmt.Sprintf("Peer %s [%s]", p.id,
		fmt.Sprintf("klay/%2d", p.version),
	)
}

// ConnType returns the conntype of the peer.
func (p *basePeer) ConnType() common.ConnType {
	return p.Peer.ConnType()
}

// GetID returns the id of the peer.
func (p *basePeer) GetID() string {
	return p.id
}

// GetP2PPeerID returns the id of the p2p.Peer.
func (p *basePeer) GetP2PPeerID() discover.NodeID {
	return p.Peer.ID()
}

// GetChainID returns the chain id of the peer.
func (p *basePeer) GetChainID() *big.Int {
	return p.chainID
}

// GetAddr returns the address of the peer.
func (p *basePeer) GetAddr() common.Address {
	return p.addr
}

// SetAddr sets the address of the peer.
func (p *basePeer) SetAddr(addr common.Address) {
	p.addr = addr
}

// GetVersion returns the version of the peer.
func (p *basePeer) GetVersion() int {
	return p.version
}

// KnowsBlock returns if the peer is known to have the block, based on knownBlocksCache.
func (p *basePeer) KnowsBlock(hash common.Hash) bool {
	_, ok := p.knownBlocksCache.Get(hash)
	return ok
}

// KnowsTx returns if the peer is known to have the transaction, based on knownTxsCache.
func (p *basePeer) KnowsTx(hash common.Hash) bool {
	_, ok := p.knownTxsCache.Get(hash)
	return ok
}

// GetP2PPeer returns the p2p.Peer.
func (p *basePeer) GetP2PPeer() *p2p.Peer {
	return p.Peer
}

// DisconnectP2PPeer disconnects the p2p peer with the given reason.
func (p *basePeer) DisconnectP2PPeer(discReason p2p.DiscReason) {
	p.GetP2PPeer().Disconnect(discReason)
}

// GetRW returns the MsgReadWriter of the peer.
func (p *basePeer) GetRW() p2p.MsgReadWriter {
	return p.rw
}

// Handle is the callback invoked to manage the life cycle of a Klaytn Peer. When
// this function terminates, the Peer is disconnected.
func (p *basePeer) Handle(pm *ProtocolManager) error {
	return pm.handle(p)
}

// UpdateRWImplementationVersion updates the version of the implementation of RW.
func (p *basePeer) UpdateRWImplementationVersion() {
	if rw, ok := p.rw.(*meteredMsgReadWriter); ok {
		rw.Init(p.GetVersion())
	}
}

// RegisterConsensusMsgCode is not supported by this peer.
func (p *basePeer) RegisterConsensusMsgCode(msgCode uint64) error {
	return fmt.Errorf("%v peerID: %v ", errNotSupportedByBasePeer, p.GetID())
}

// singleChannelPeer is a peer that uses a single channel.
type singleChannelPeer struct {
	*basePeer
}

// multiChannelPeer is a peer that uses a multi channel.
type multiChannelPeer struct {
	*basePeer                     // basePeer is a set of data structures that the peer implementation has in common
	rws       []p2p.MsgReadWriter // rws is a slice of p2p.MsgReadWriter for peer-to-peer transmission and reception

	chMgr *ChannelManager
}

// RegisterMsgCode registers the channel id corresponding to msgCode.
func (p *multiChannelPeer) RegisterMsgCode(channelId uint, msgCode uint64) {
	p.chMgr.RegisterMsgCode(channelId, msgCode)
}

// RegisterConsensusMsgCode registers the channel of consensus msg.
func (p *multiChannelPeer) RegisterConsensusMsgCode(msgCode uint64) error {
	p.chMgr.RegisterMsgCode(ConsensusChannel, msgCode)
	return nil
}

// Broadcast is a write loop that multiplexes block propagations, announcements
// and transaction broadcasts into the remote peer. The goal is to have an async
// writer that does not lock up node internals.
func (p *multiChannelPeer) Broadcast() {
	for {
		select {
		case txs := <-p.queuedTxs:
			if err := p.SendTransactions(txs); err != nil {
				logger.Error("fail to SendTransactions", "peer", p.id, "err", err)
				continue
				//return
			}
			p.Log().Trace("Broadcast transactions", "peer", p.id, "count", len(txs))

		case prop := <-p.queuedProps:
			if err := p.SendNewBlock(prop.block, prop.td); err != nil {
				logger.Error("fail to SendNewBlock", "peer", p.id, "err", err)
				continue
				//return
			}
			p.Log().Trace("Propagated block", "peer", p.id, "number", prop.block.Number(), "hash", prop.block.Hash(), "td", prop.td)

		case block := <-p.queuedAnns:
			if err := p.SendNewBlockHashes([]common.Hash{block.Hash()}, []uint64{block.NumberU64()}); err != nil {
				logger.Error("fail to SendNewBlockHashes", "peer", p.id, "err", err)
				continue
				//return
			}
			p.Log().Trace("Announced block", "peer", p.id, "number", block.Number(), "hash", block.Hash())

		case <-p.term:
			p.Log().Debug("Peer broadcast loop end", "peer", p.id)
			return
		}
	}
}

// SendTransactions sends transactions to the peer and includes the hashes
// in its transaction hash set for future reference.
func (p *multiChannelPeer) SendTransactions(txs types.Transactions) error {
	// Before sending transactions, sort transactions in ascending order by time.
	if !sort.IsSorted(types.TxByPriceAndTime(txs)) {
		sort.Sort(types.TxByPriceAndTime(txs))
	}

	for _, tx := range txs {
		p.AddToKnownTxs(tx.Hash())
	}
	return p.msgSender(TxMsg, txs)
}

// ReSendTransactions sends txs to a peer in order to prevent the txs from missing.
func (p *multiChannelPeer) ReSendTransactions(txs types.Transactions) error {
	// Before sending transactions, sort transactions in ascending order by time.
	if !sort.IsSorted(types.TxByPriceAndTime(txs)) {
		sort.Sort(types.TxByPriceAndTime(txs))
	}

	return p.msgSender(TxMsg, txs)
}

// SendNewBlockHashes announces the availability of a number of blocks through
// a hash notification.
func (p *multiChannelPeer) SendNewBlockHashes(hashes []common.Hash, numbers []uint64) error {
	for _, hash := range hashes {
		p.AddToKnownBlocks(hash)
	}
	request := make(newBlockHashesData, len(hashes))
	for i := 0; i < len(hashes); i++ {
		request[i].Hash = hashes[i]
		request[i].Number = numbers[i]
	}
	return p.msgSender(NewBlockHashesMsg, request)
}

// SendNewBlock propagates an entire block to a remote peer.
func (p *multiChannelPeer) SendNewBlock(block *types.Block, td *big.Int) error {
	p.AddToKnownBlocks(block.Hash())
	return p.msgSender(NewBlockMsg, []interface{}{block, td})
}

// SendBlockHeaders sends a batch of block headers to the remote peer.
func (p *multiChannelPeer) SendBlockHeaders(headers []*types.Header) error {
	return p.msgSender(BlockHeadersMsg, headers)
}

// SendFetchedBlockHeader sends a block header to the remote peer, requested by fetcher.
func (p *multiChannelPeer) SendFetchedBlockHeader(header *types.Header) error {
	return p.msgSender(BlockHeaderFetchResponseMsg, header)
}

// SendBlockBodies sends a batch of block contents to the remote peer.
func (p *multiChannelPeer) SendBlockBodies(bodies []*blockBody) error {
	return p.msgSender(BlockBodiesMsg, blockBodiesData(bodies))
}

// SendBlockBodiesRLP sends a batch of block contents to the remote peer from
// an already RLP encoded format.
func (p *multiChannelPeer) SendBlockBodiesRLP(bodies []rlp.RawValue) error {
	return p.msgSender(BlockBodiesMsg, bodies)
}

// SendFetchedBlockBodiesRLP sends a batch of block contents to the remote peer from
// an already RLP encoded format.
func (p *multiChannelPeer) SendFetchedBlockBodiesRLP(bodies []rlp.RawValue) error {
	return p.msgSender(BlockBodiesFetchResponseMsg, bodies)
}

// SendNodeDataRLP sends a batch of arbitrary internal data, corresponding to the
// hashes requested.
func (p *multiChannelPeer) SendNodeData(data [][]byte) error {
	return p.msgSender(NodeDataMsg, data)
}

// SendReceiptsRLP sends a batch of transaction receipts, corresponding to the
// ones requested from an already RLP encoded format.
func (p *multiChannelPeer) SendReceiptsRLP(receipts []rlp.RawValue) error {
	return p.msgSender(ReceiptsMsg, receipts)
}

// FetchBlockHeader is a wrapper around the header query functions to fetch a
// single header. It is used solely by the fetcher.
func (p *multiChannelPeer) FetchBlockHeader(hash common.Hash) error {
	p.Log().Debug("Fetching a new block header", "hash", hash)
	return p.msgSender(BlockHeaderFetchRequestMsg, hash)
}

// RequestHeadersByHash fetches a batch of blocks' headers corresponding to the
// specified header query, based on the hash of an origin block.
func (p *multiChannelPeer) RequestHeadersByHash(origin common.Hash, amount int, skip int, reverse bool) error {
	p.Log().Debug("Fetching batch of headers", "count", amount, "fromhash", origin, "skip", skip, "reverse", reverse)
	return p.msgSender(BlockHeadersRequestMsg, &getBlockHeadersData{Origin: hashOrNumber{Hash: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestHeadersByNumber fetches a batch of blocks' headers corresponding to the
// specified header query, based on the number of an origin block.
func (p *multiChannelPeer) RequestHeadersByNumber(origin uint64, amount int, skip int, reverse bool) error {
	p.Log().Debug("Fetching batch of headers", "count", amount, "fromnum", origin, "skip", skip, "reverse", reverse)
	return p.msgSender(BlockHeadersRequestMsg, &getBlockHeadersData{Origin: hashOrNumber{Number: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestBodies fetches a batch of blocks' bodies corresponding to the hashes
// specified.
func (p *multiChannelPeer) RequestBodies(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of block bodies", "count", len(hashes))
	return p.msgSender(BlockBodiesRequestMsg, hashes)
}

// FetchBlockBodies fetches a batch of blocks' bodies corresponding to the hashes
// specified.
func (p *multiChannelPeer) FetchBlockBodies(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of new block bodies", "count", len(hashes))
	return p.msgSender(BlockBodiesFetchRequestMsg, hashes)
}

// RequestNodeData fetches a batch of arbitrary data from a node's known state
// data, corresponding to the specified hashes.
func (p *multiChannelPeer) RequestNodeData(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of state data", "count", len(hashes))
	return p.msgSender(NodeDataRequestMsg, hashes)
}

// RequestReceipts fetches a batch of transaction receipts from a remote node.
func (p *multiChannelPeer) RequestReceipts(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of receipts", "count", len(hashes))
	return p.msgSender(ReceiptsRequestMsg, hashes)
}

// msgSender sends data to the peer.
func (p *multiChannelPeer) msgSender(msgcode uint64, data interface{}) error {
	if ch, ok := ChannelOfMessage[msgcode]; ok && len(p.rws) > ch {
		return p2p.Send(p.rws[ch], msgcode, data)
	} else {
		return errors.New("RW not found for message")
	}
}

// GetRW returns the MsgReadWriter of the peer.
func (p *multiChannelPeer) GetRW() p2p.MsgReadWriter {
	return p.rw //TODO-Klaytn check this function usage
}

// UpdateRWImplementationVersion updates the version of the implementation of RW.
func (p *multiChannelPeer) UpdateRWImplementationVersion() {
	for _, rw := range p.rws {
		if rw, ok := rw.(*meteredMsgReadWriter); ok {
			rw.Init(p.GetVersion())
		}
	}
	p.basePeer.UpdateRWImplementationVersion()
}

func (p *multiChannelPeer) ReadMsg(rw p2p.MsgReadWriter, connectionOrder int, errCh chan<- error, wg *sync.WaitGroup, closed <-chan struct{}) {
	defer wg.Done()
	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			p.GetP2PPeer().Log().Warn("ProtocolManager failed to read msg", "err", err)
			errCh <- err
			return
		}

		msgCh, err := p.chMgr.GetChannelWithMsgCode(connectionOrder, msg.Code)
		if err != nil {
			p.GetP2PPeer().Log().Warn("ProtocolManager failed to get msg channel", "err", err)
			errCh <- err
			return
		}

		if msg.Size > ProtocolMaxMsgSize {
			err := errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
			p.GetP2PPeer().Log().Warn("ProtocolManager over max msg size", "err", err)
			errCh <- err
			return
		}

		select {
		case msgCh <- msg:
		case <-closed:
			return
		}
	}
}

// Handle is the callback invoked to manage the life cycle of a Klaytn Peer. When
// this function terminates, the Peer is disconnected.
func (p *multiChannelPeer) Handle(pm *ProtocolManager) error {
	// Ignore maxPeers if this is a trusted peer
	if pm.peers.Len() >= pm.maxPeers && !p.GetP2PPeer().Info().Networks[p2p.ConnDefault].Trusted {
		return p2p.DiscTooManyPeers
	}
	p.GetP2PPeer().Log().Debug("Klaytn peer connected", "name", p.GetP2PPeer().Name())

	// Execute the handshake
	var (
		genesis = pm.blockchain.Genesis()
		head    = pm.blockchain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = pm.blockchain.GetTd(hash, number)
	)

	err := p.Handshake(pm.networkId, pm.getChainID(), td, hash, genesis.Hash())
	if err != nil {
		p.GetP2PPeer().Log().Debug("Klaytn peer handshake failed", "err", err)
		return err
	}

	p.UpdateRWImplementationVersion()

	// Register the peer locally
	if err := pm.peers.Register(p); err != nil {
		// if starting node with unlock account, can't register peer until finish unlock
		p.GetP2PPeer().Log().Info("Klaytn peer registration failed", "err", err)
		return err
	}
	defer pm.removePeer(p.GetID())

	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := pm.downloader.RegisterPeer(p.GetID(), p.GetVersion(), p); err != nil {
		return err
	}
	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	pm.syncTransactions(p)

	p.GetP2PPeer().Log().Info("Added a multichannel P2P Peer", "peerID", p.GetP2PPeerID())

	pubKey, err := p.GetP2PPeerID().Pubkey()
	if err != nil {
		return err
	}
	addr := crypto.PubkeyToAddress(*pubKey)
	lenRWs := len(p.rws)

	var wg sync.WaitGroup
	// TODO-GX check global worker and peer worker
	messageChannels := make([]chan p2p.Msg, 0, lenRWs)
	var consensusChannel chan p2p.Msg
	isCN := false

	if _, ok := pm.engine.(consensus.Handler); ok && pm.nodetype == common.CONSENSUSNODE {
		consensusChannel = make(chan p2p.Msg, channelSizePerPeer)
		defer close(consensusChannel)
		pm.engine.(consensus.Handler).RegisterConsensusMsgCode(p)
		isCN = true
	}

	for idx := range p.rws {
		channel := make(chan p2p.Msg, channelSizePerPeer)
		defer close(channel)
		messageChannels = append(messageChannels, channel)

		p.chMgr.RegisterChannelWithIndex(idx, BlockChannel, channel)
		p.chMgr.RegisterChannelWithIndex(idx, TxChannel, channel)
		p.chMgr.RegisterChannelWithIndex(idx, MiscChannel, channel)

		if isCN {
			p.chMgr.RegisterChannelWithIndex(idx, ConsensusChannel, consensusChannel)
		}
	}

	sumOfGoroutineForProcessMessage := 1 // 1 is for consensusChannel
	for connIdx := range messageChannels {
		sumOfGoroutineForProcessMessage += ConcurrentOfChannel[connIdx]
	}
	errChannel := make(chan error, lenRWs+sumOfGoroutineForProcessMessage) // errChannel size should be set to count of goroutine use errChannel
	closed := make(chan struct{})

	if isCN {
		go pm.processConsensusMsg(consensusChannel, p, addr, errChannel)
	}

	for connIdx, messageChannel := range messageChannels {
		for i := 0; i < ConcurrentOfChannel[connIdx]; i++ {
			go pm.processMsg(messageChannel, p, addr, errChannel)
		}
	}

	for idx, rw := range p.rws {
		wg.Add(1)
		go p.ReadMsg(rw, idx, errChannel, &wg, closed)
	}

	err = <-errChannel
	close(closed)
	wg.Wait()
	p.GetP2PPeer().Log().Info("Disconnected a multichannel P2P Peer", "peerID", p.GetP2PPeerID(), "peerName", p.GetP2PPeer().Name(), "err", err)
	return err
}

type ByPassValidator struct{}

func (v ByPassValidator) ValidatePeerType(addr common.Address) error {
	return nil
}
