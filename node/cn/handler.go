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
// This file is derived from eth/handler.go (2018/06/04).
// Modified and improved for the klaytn development.

package cn

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"runtime/debug"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/datasync/fetcher"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/klaytn/klaytn/work"
)

const (
	softResponseLimit = 2 * 1024 * 1024 // Target maximum size of returned blocks, headers or node data.
	estHeaderRlpSize  = 500             // Approximate size of an RLP encoded block header

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096

	concurrentPerPeer  = 3
	channelSizePerPeer = 20

	blockReceivingPNLimit  = 5 // maximum number of PNs that a CN broadcasts block.
	minNumPeersToSendBlock = 3 // minimum number of peers that a node broadcasts block.

	// DefaultMaxResendTxCount is the number of resending transactions to peer in order to prevent the txs from missing.
	DefaultMaxResendTxCount = 1000

	// DefaultTxResendInterval is the second of resending transactions period.
	DefaultTxResendInterval = 4
)

// errIncompatibleConfig is returned if the requested protocols and configs are
// not compatible (low protocol version restrictions and high requirements).
var errIncompatibleConfig = errors.New("incompatible configuration")
var errUnknownProcessingError = errors.New("unknown error during the msg processing")

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

type ProtocolManager struct {
	networkId uint64

	fastSync  uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	acceptTxs uint32 // Flag whether we're considered synchronised (enables transaction processing)

	txpool      work.TxPool
	blockchain  work.BlockChain
	chainconfig *params.ChainConfig
	maxPeers    int

	downloader ProtocolManagerDownloader
	fetcher    ProtocolManagerFetcher
	peers      PeerSet

	SubProtocols []p2p.Protocol

	eventMux      *event.TypeMux
	txsCh         chan blockchain.NewTxsEvent
	txsSub        event.Subscription
	minedBlockSub *event.TypeMuxSubscription

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan Peer
	txsyncCh    chan *txsync
	quitSync    chan struct{}
	noMorePeers chan struct{}

	quitResendCh chan struct{}
	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg sync.WaitGroup
	// istanbul BFT
	engine consensus.Engine

	rewardbase   common.Address
	rewardwallet accounts.Wallet

	wsendpoint string

	nodetype          common.ConnType
	txResendUseLegacy bool

	//syncStop is a flag to stop peer sync
	syncStop int32
}

// NewProtocolManager returns a new Klaytn sub protocol manager. The Klaytn sub protocol manages peers capable
// with the Klaytn network.
func NewProtocolManager(config *params.ChainConfig, mode downloader.SyncMode, networkId uint64, mux *event.TypeMux,
	txpool work.TxPool, engine consensus.Engine, blockchain work.BlockChain, chainDB database.DBManager, cacheLimit int,
	nodetype common.ConnType, cnconfig *Config) (*ProtocolManager, error) {
	// Create the protocol maanger with the base fields
	manager := &ProtocolManager{
		networkId:         networkId,
		eventMux:          mux,
		txpool:            txpool,
		blockchain:        blockchain,
		chainconfig:       config,
		peers:             newPeerSet(),
		newPeerCh:         make(chan Peer),
		noMorePeers:       make(chan struct{}),
		txsyncCh:          make(chan *txsync),
		quitSync:          make(chan struct{}),
		quitResendCh:      make(chan struct{}),
		engine:            engine,
		nodetype:          nodetype,
		txResendUseLegacy: cnconfig.TxResendUseLegacy,
	}

	// istanbul BFT
	if handler, ok := engine.(consensus.Handler); ok {
		handler.SetBroadcaster(manager, manager.nodetype)
	}

	// Figure out whether to allow fast sync or not
	if mode == downloader.FastSync && blockchain.CurrentBlock().NumberU64() > 0 {
		logger.Error("Blockchain not empty, fast sync disabled")
		mode = downloader.FullSync
	}
	if mode == downloader.FastSync {
		manager.fastSync = uint32(1)
	}
	// istanbul BFT
	protocol := engine.Protocol()
	// Initiate a sub-protocol for every implemented version we can handle
	manager.SubProtocols = make([]p2p.Protocol, 0, len(protocol.Versions))
	for i, version := range protocol.Versions {
		// Skip protocol version if incompatible with the mode of operation
		if mode == downloader.FastSync && version < klay63 {
			continue
		}
		// Compatible; initialise the sub-protocol
		version := version
		manager.SubProtocols = append(manager.SubProtocols, p2p.Protocol{
			Name:    protocol.Name,
			Version: version,
			Length:  protocol.Lengths[i],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := manager.newPeer(int(version), p, rw)
				pubKey, err := p.ID().Pubkey()
				if err != nil {
					if p.ConnType() == common.CONSENSUSNODE {
						return err
					}
					peer.SetAddr(common.Address{})
				} else {
					addr := crypto.PubkeyToAddress(*pubKey)
					peer.SetAddr(addr)
				}
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					return manager.handle(peer)
				case <-manager.quitSync:
					return p2p.DiscQuitting
				}
			},
			RunWithRWs: func(p *p2p.Peer, rws []p2p.MsgReadWriter) error {
				peer, err := manager.newPeerWithRWs(int(version), p, rws)
				if err != nil {
					return err
				}
				pubKey, err := p.ID().Pubkey()
				if err != nil {
					if p.ConnType() == common.CONSENSUSNODE {
						return err
					}
					peer.SetAddr(common.Address{})
				} else {
					addr := crypto.PubkeyToAddress(*pubKey)
					peer.SetAddr(addr)
				}
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					return peer.Handle(manager)
				case <-manager.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return manager.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := manager.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}

	if len(manager.SubProtocols) == 0 {
		return nil, errIncompatibleConfig
	}

	// Create and set downloader
	if cnconfig.DownloaderDisable {
		manager.downloader = downloader.NewFakeDownloader()
	} else {
		// Construct the downloader (long sync) and its backing state bloom if fast
		// sync is requested. The downloader is responsible for deallocating the state

		// bloom when it's done.
		var stateBloom *statedb.SyncBloom
		if atomic.LoadUint32(&manager.fastSync) == 1 {
			stateBloom = statedb.NewSyncBloom(uint64(cacheLimit), chainDB.GetStateTrieDB())
		}
		manager.downloader = downloader.New(mode, chainDB, stateBloom, manager.eventMux, blockchain, nil, manager.removePeer)
	}

	// Create and set fetcher
	if cnconfig.FetcherDisable {
		manager.fetcher = fetcher.NewFakeFetcher()
	} else {
		validator := func(header *types.Header) error {
			return engine.VerifyHeader(blockchain, header, true)
		}
		heighter := func() uint64 {
			return blockchain.CurrentBlock().NumberU64()
		}
		inserter := func(blocks types.Blocks) (int, error) {
			// If fast sync is running, deny importing weird blocks
			if atomic.LoadUint32(&manager.fastSync) == 1 {
				logger.Warn("Discarded bad propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
				return 0, nil
			}
			atomic.StoreUint32(&manager.acceptTxs, 1) // Mark initial sync done on any fetcher import
			return manager.blockchain.InsertChain(blocks)
		}
		manager.fetcher = fetcher.New(blockchain.GetBlockByHash, validator, manager.BroadcastBlock, manager.BroadcastBlockHash, heighter, inserter, manager.removePeer)
	}

	if manager.useTxResend() {
		go manager.txResendLoop(cnconfig.TxResendInterval, cnconfig.TxResendCount)
	}
	return manager, nil
}

// istanbul BFT
func (pm *ProtocolManager) RegisterValidator(connType common.ConnType, validator p2p.PeerTypeValidator) {
	pm.peers.RegisterValidator(connType, validator)
}

func (pm *ProtocolManager) getWSEndPoint() string {
	return pm.wsendpoint
}

func (pm *ProtocolManager) SetRewardbase(addr common.Address) {
	pm.rewardbase = addr
}

func (pm *ProtocolManager) SetRewardbaseWallet(wallet accounts.Wallet) {
	pm.rewardwallet = wallet
}

func (pm *ProtocolManager) removePeer(id string) {
	// Short circuit if the peer was already removed
	peer := pm.peers.Peer(id)
	if peer == nil {
		return
	}
	logger.Debug("Removing Klaytn peer", "peer", id)

	// Unregister the peer from the downloader and peer set
	pm.downloader.UnregisterPeer(id)
	if err := pm.peers.Unregister(id); err != nil {
		logger.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.GetP2PPeer().Disconnect(p2p.DiscUselessPeer)
	}
}

// getChainID returns the current chain id.
func (pm *ProtocolManager) getChainID() *big.Int {
	return pm.blockchain.Config().ChainID
}

func (pm *ProtocolManager) Start(maxPeers int) {
	pm.maxPeers = maxPeers

	// broadcast transactions
	pm.txsCh = make(chan blockchain.NewTxsEvent, txChanSize)
	pm.txsSub = pm.txpool.SubscribeNewTxsEvent(pm.txsCh)
	go pm.txBroadcastLoop()

	// broadcast mined blocks
	pm.minedBlockSub = pm.eventMux.Subscribe(blockchain.NewMinedBlockEvent{})
	go pm.minedBroadcastLoop()

	// start sync handlers
	go pm.syncer()
	go pm.txsyncLoop()
}

func (pm *ProtocolManager) Stop() {
	logger.Info("Stopping Klaytn protocol")

	pm.txsSub.Unsubscribe()        // quits txBroadcastLoop
	pm.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	// Quit the sync loop.
	// After this send has completed, no new peers will be accepted.
	pm.noMorePeers <- struct{}{}

	if pm.useTxResend() {
		// Quit resend loop
		pm.quitResendCh <- struct{}{}
	}

	// Quit fetcher, txsyncLoop.
	close(pm.quitSync)

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to pm.peers yet
	// will exit when they try to register.
	pm.peers.Close()

	// Wait for all peer handler goroutines and the loops to come down.
	pm.wg.Wait()

	logger.Info("Klaytn protocol stopped")
}

// SetSyncStop sets value of syncStop flag. If it's true, peer sync process does not proceed.
func (pm *ProtocolManager) SetSyncStop(flag bool) {
	var i int32 = 0
	if flag {
		i = 1
	}
	atomic.StoreInt32(&(pm.syncStop), int32(i))
}

func (pm *ProtocolManager) GetSyncStop() bool {
	if atomic.LoadInt32(&(pm.syncStop)) != 0 {
		return true
	}
	return false
}

func (pm *ProtocolManager) newPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter) Peer {
	return newPeer(pv, p, newMeteredMsgWriter(rw))
}

// newPeerWithRWs creates a new Peer object with a slice of p2p.MsgReadWriter.
func (pm *ProtocolManager) newPeerWithRWs(pv int, p *p2p.Peer, rws []p2p.MsgReadWriter) (Peer, error) {
	meteredRWs := make([]p2p.MsgReadWriter, 0, len(rws))
	for _, rw := range rws {
		meteredRWs = append(meteredRWs, newMeteredMsgWriter(rw))
	}
	return newPeerWithRWs(pv, p, meteredRWs)
}

// handle is the callback invoked to manage the life cycle of a Klaytn peer. When
// this function terminates, the peer is disconnected.
func (pm *ProtocolManager) handle(p Peer) error {
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
	if rw, ok := p.GetRW().(*meteredMsgReadWriter); ok {
		rw.Init(p.GetVersion())
	}

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

	p.GetP2PPeer().Log().Info("Added a single channel P2P Peer", "peerID", p.GetP2PPeerID())

	pubKey, err := p.GetP2PPeerID().Pubkey()
	if err != nil {
		return err
	}
	addr := crypto.PubkeyToAddress(*pubKey)

	// TODO-Klaytn check global worker and peer worker
	messageChannel := make(chan p2p.Msg, channelSizePerPeer)
	defer close(messageChannel)
	errChannel := make(chan error, channelSizePerPeer)
	for w := 1; w <= concurrentPerPeer; w++ {
		go pm.processMsg(messageChannel, p, addr, errChannel)
	}

	// main loop. handle incoming messages.
	for {
		msg, err := p.GetRW().ReadMsg()
		if err != nil {
			p.GetP2PPeer().Log().Warn("ProtocolManager failed to read msg", "err", err)
			return err
		}
		if msg.Size > ProtocolMaxMsgSize {
			err := errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
			p.GetP2PPeer().Log().Warn("ProtocolManager over max msg size", "err", err)
			return err
		}

		select {
		case err := <-errChannel:
			return err
		case messageChannel <- msg:
		}
		//go pm.handleMsg(p, addr, msg)

		//if err := pm.handleMsg(p); err != nil {
		//	p.Log().Debug("Klaytn message handling failed", "err", err)
		//	return err
		//}
	}
}

func (pm *ProtocolManager) processMsg(msgCh <-chan p2p.Msg, p Peer, addr common.Address, errCh chan<- error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("stacktrace from panic: \n" + string(debug.Stack()))
			logger.Warn("the panic is recovered", "panicErr", err)
			errCh <- errUnknownProcessingError
		}
	}()

	_, fakeF := pm.fetcher.(*fetcher.FakeFetcher)
	_, fakeD := pm.downloader.(*downloader.FakeDownloader)
	if fakeD || fakeF {
		p.GetP2PPeer().Log().Warn("ProtocolManager does not handle p2p messages", "fakeFetcher", fakeF, "fakeDownloader", fakeD)
		for msg := range msgCh {
			msg.Discard()
		}
	} else {
		for msg := range msgCh {
			if err := pm.handleMsg(p, addr, msg); err != nil {
				p.GetP2PPeer().Log().Error("ProtocolManager failed to handle message", "msg", msg, "err", err)
				errCh <- err
				return
			}
			msg.Discard()
		}
	}

	p.GetP2PPeer().Log().Debug("ProtocolManager.processMsg closed", "PeerName", p.GetP2PPeer().Name())
}

// processConsensusMsg processes the consensus message.
func (pm *ProtocolManager) processConsensusMsg(msgCh <-chan p2p.Msg, p Peer, addr common.Address, errCh chan<- error) {
	for msg := range msgCh {
		if handler, ok := pm.engine.(consensus.Handler); ok {
			_, err := handler.HandleMsg(addr, msg)
			// if msg is istanbul msg, handled is true and err is nil if handle msg is successful.
			if err != nil {
				p.GetP2PPeer().Log().Warn("ProtocolManager failed to handle consensus message. This can happen during block synchronization.", "msg", msg, "err", err)
				errCh <- err
				return
			}
		}
		msg.Discard()
	}
	p.GetP2PPeer().Log().Info("ProtocolManager.processConsensusMsg closed", "PeerName", p.GetP2PPeer().Name())
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (pm *ProtocolManager) handleMsg(p Peer, addr common.Address, msg p2p.Msg) error {
	// Below message size checking is done by handle().
	// Read the next message from the remote peer, and ensure it's fully consumed
	//msg, err := p.rw.ReadMsg()
	//if err != nil {
	//	return err
	//}
	//if msg.Size > ProtocolMaxMsgSize {
	//	return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	//}
	//defer msg.Discard()

	// istanbul BFT
	if handler, ok := pm.engine.(consensus.Handler); ok {
		//pubKey, err := p.ID().Pubkey()
		//if err != nil {
		//	return err
		//}
		//addr := crypto.PubkeyToAddress(*pubKey)
		handled, err := handler.HandleMsg(addr, msg)
		// if msg is istanbul msg, handled is true and err is nil if handle msg is successful.
		if handled {
			return err
		}
	}

	// Handle the message depending on its contents
	switch {
	case msg.Code == StatusMsg:
		// Status messages should never arrive after the handshake
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")

		// Block header query, collect the requested headers and reply
	case msg.Code == BlockHeadersRequestMsg:
		if err := handleBlockHeadersRequestMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == BlockHeadersMsg:
		if err := handleBlockHeadersMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == BlockBodiesRequestMsg:
		if err := handleBlockBodiesRequestMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == BlockBodiesMsg:
		if err := handleBlockBodiesMsg(pm, p, msg); err != nil {
			return err
		}

	case p.GetVersion() >= klay63 && msg.Code == NodeDataRequestMsg:
		if err := handleNodeDataRequestMsg(pm, p, msg); err != nil {
			return err
		}

	case p.GetVersion() >= klay63 && msg.Code == NodeDataMsg:
		if err := handleNodeDataMsg(pm, p, msg); err != nil {
			return err
		}

	case p.GetVersion() >= klay63 && msg.Code == ReceiptsRequestMsg:
		if err := handleReceiptsRequestMsg(pm, p, msg); err != nil {
			return err
		}

	case p.GetVersion() >= klay63 && msg.Code == ReceiptsMsg:
		if err := handleReceiptsMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == NewBlockHashesMsg:
		if err := handleNewBlockHashesMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == BlockHeaderFetchRequestMsg:
		if err := handleBlockHeaderFetchRequestMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == BlockHeaderFetchResponseMsg:
		if err := handleBlockHeaderFetchResponseMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == BlockBodiesFetchRequestMsg:
		if err := handleBlockBodiesFetchRequestMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == BlockBodiesFetchResponseMsg:
		if err := handleBlockBodiesFetchResponseMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == NewBlockMsg:
		if err := handleNewBlockMsg(pm, p, msg); err != nil {
			return err
		}

	case msg.Code == TxMsg:
		if err := handleTxMsg(pm, p, msg); err != nil {
			return err
		}

	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}
	return nil
}

// handleBlockHeadersRequestMsg handles block header request message.
func handleBlockHeadersRequestMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// Decode the complex header query
	var query getBlockHeadersData
	if err := msg.Decode(&query); err != nil {
		return errResp(ErrDecode, "%v: %v", msg, err)
	}
	hashMode := query.Origin.Hash != (common.Hash{})

	// Gather headers until the fetch or network limits is reached
	var (
		bytes   common.StorageSize
		headers []*types.Header
		unknown bool
	)
	for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit && len(headers) < downloader.MaxHeaderFetch {
		// Retrieve the next header satisfying the query
		var origin *types.Header
		if hashMode {
			origin = pm.blockchain.GetHeaderByHash(query.Origin.Hash)
		} else {
			origin = pm.blockchain.GetHeaderByNumber(query.Origin.Number)
		}
		if origin == nil {
			break
		}
		number := origin.Number.Uint64()
		headers = append(headers, origin)
		bytes += estHeaderRlpSize

		// Advance to the next header of the query
		switch {
		case query.Origin.Hash != (common.Hash{}) && query.Reverse:
			// Hash based traversal towards the genesis block
			for i := 0; i < int(query.Skip)+1; i++ {
				if header := pm.blockchain.GetHeader(query.Origin.Hash, number); header != nil {
					query.Origin.Hash = header.ParentHash
					number--
				} else {
					unknown = true
					break
				}
			}
		case query.Origin.Hash != (common.Hash{}) && !query.Reverse:
			// Hash based traversal towards the leaf block
			var (
				current = origin.Number.Uint64()
				next    = current + query.Skip + 1
			)
			if next <= current {
				infos, _ := json.MarshalIndent(p.GetP2PPeer().Info(), "", "  ")
				p.GetP2PPeer().Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
				unknown = true
			} else {
				if header := pm.blockchain.GetHeaderByNumber(next); header != nil {
					if pm.blockchain.GetBlockHashesFromHash(header.Hash(), query.Skip+1)[query.Skip] == query.Origin.Hash {
						query.Origin.Hash = header.Hash()
					} else {
						unknown = true
					}
				} else {
					unknown = true
				}
			}
		case query.Reverse:
			// Number based traversal towards the genesis block
			if query.Origin.Number >= query.Skip+1 {
				query.Origin.Number -= query.Skip + 1
			} else {
				unknown = true
			}

		case !query.Reverse:
			// Number based traversal towards the leaf block
			query.Origin.Number += query.Skip + 1
		}
	}
	return p.SendBlockHeaders(headers)
}

// handleBlockHeadersMsg handles block header response message.
func handleBlockHeadersMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// A batch of headers arrived to one of our previous requests
	var headers []*types.Header
	if err := msg.Decode(&headers); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	if err := pm.downloader.DeliverHeaders(p.GetID(), headers); err != nil {
		logger.Debug("Failed to deliver headers", "err", err)
	}
	return nil
}

// handleBlockBodiesRequest handles common part for handleBlockBodiesRequest and
// handleBlockBodiesFetchRequestMsg. It decodes the message to get list of hashes
// and then send block bodies corresponding to those hashes.
func handleBlockBodiesRequest(pm *ProtocolManager, p Peer, msg p2p.Msg) ([]rlp.RawValue, error) {
	// Decode the retrieval message
	msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
	if _, err := msgStream.List(); err != nil {
		return nil, err
	}
	// Gather blocks until the fetch or network limits is reached
	var (
		hash   common.Hash
		bytes  int
		bodies []rlp.RawValue
	)
	for bytes < softResponseLimit && len(bodies) < downloader.MaxBlockFetch {
		// Retrieve the hash of the next block
		if err := msgStream.Decode(&hash); err == rlp.EOL {
			break
		} else if err != nil {
			return nil, errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Retrieve the requested block body, stopping if enough was found
		if data := pm.blockchain.GetBodyRLP(hash); len(data) != 0 {
			bodies = append(bodies, data)
			bytes += len(data)
		}
	}

	return bodies, nil
}

// handleBlockBodiesRequestMsg handles block body request message.
func handleBlockBodiesRequestMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	if bodies, err := handleBlockBodiesRequest(pm, p, msg); err != nil {
		return err
	} else {
		return p.SendBlockBodiesRLP(bodies)
	}
}

// handleGetBlockBodiesMsg handles block body response message.
func handleBlockBodiesMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// A batch of block bodies arrived to one of our previous requests
	var request blockBodiesData
	if err := msg.Decode(&request); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// Deliver them all to the downloader for queuing
	transactions := make([][]*types.Transaction, len(request))

	for i, body := range request {
		transactions[i] = body.Transactions
	}

	err := pm.downloader.DeliverBodies(p.GetID(), transactions)
	if err != nil {
		logger.Debug("Failed to deliver bodies", "err", err)
	}

	return nil
}

// handleNodeDataRequestMsg handles node data request message.
func handleNodeDataRequestMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// Decode the retrieval message
	msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
	if _, err := msgStream.List(); err != nil {
		return err
	}
	// Gather state data until the fetch or network limits is reached
	var (
		hash  common.Hash
		bytes int
		data  [][]byte
	)
	for bytes < softResponseLimit && len(data) < downloader.MaxStateFetch {
		// Retrieve the hash of the next state entry
		if err := msgStream.Decode(&hash); err == rlp.EOL {
			break
		} else if err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Retrieve the requested state entry, stopping if enough was found
		if entry, err := pm.blockchain.TrieNode(hash); err == nil {
			data = append(data, entry)
			bytes += len(entry)
		}
	}
	return p.SendNodeData(data)
}

// handleNodeDataMsg handles node data response message.
func handleNodeDataMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// A batch of node state data arrived to one of our previous requests
	var data [][]byte
	if err := msg.Decode(&data); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// Deliver all to the downloader
	if err := pm.downloader.DeliverNodeData(p.GetID(), data); err != nil {
		logger.Debug("Failed to deliver node state data", "err", err)
	}
	return nil
}

// handleGetReceiptsMsg handles receipt request message.
func handleReceiptsRequestMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// Decode the retrieval message
	msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
	if _, err := msgStream.List(); err != nil {
		return err
	}
	// Gather state data until the fetch or network limits is reached
	var (
		hash     common.Hash
		bytes    int
		receipts []rlp.RawValue
	)
	for bytes < softResponseLimit && len(receipts) < downloader.MaxReceiptFetch {
		// Retrieve the hash of the next block
		if err := msgStream.Decode(&hash); err == rlp.EOL {
			break
		} else if err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Retrieve the requested block's receipts, skipping if unknown to us
		results := pm.blockchain.GetReceiptsByBlockHash(hash)
		if results == nil {
			if header := pm.blockchain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
				continue
			}
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(results); err != nil {
			logger.Error("Failed to encode receipt", "err", err)
		} else {
			receipts = append(receipts, encoded)
			bytes += len(encoded)
		}
	}
	return p.SendReceiptsRLP(receipts)
}

// handleReceiptsMsg handles receipt response message.
func handleReceiptsMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// A batch of receipts arrived to one of our previous requests
	var receipts [][]*types.Receipt
	if err := msg.Decode(&receipts); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// Deliver all to the downloader
	if err := pm.downloader.DeliverReceipts(p.GetID(), receipts); err != nil {
		logger.Debug("Failed to deliver receipts", "err", err)
	}
	return nil
}

// handleNewBlockHashesMsg handles new block hashes message.
func handleNewBlockHashesMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	var (
		announces     newBlockHashesData
		maxTD         uint64
		candidateHash *common.Hash
	)
	if err := msg.Decode(&announces); err != nil {
		return errResp(ErrDecode, "%v: %v", msg, err)
	}
	// Mark the hashes as present at the remote node
	// Schedule all the unknown hashes for retrieval
	for _, block := range announces {
		p.AddToKnownBlocks(block.Hash)

		if maxTD < block.Number {
			maxTD = block.Number
			candidateHash = &block.Hash
		}
		if !pm.blockchain.HasBlock(block.Hash, block.Number) {
			pm.fetcher.Notify(p.GetID(), block.Hash, block.Number, time.Now(), p.FetchBlockHeader, p.FetchBlockBodies)
		}
	}
	blockTD := big.NewInt(int64(maxTD))
	if _, td := p.Head(); blockTD.Cmp(td) > 0 && candidateHash != nil {
		p.SetHead(*candidateHash, blockTD)
	}
	return nil
}

// handleBlockHeaderFetchRequestMsg handles block header fetch request message.
// It will send a header that the peer requested.
// If the peer requests a header which does not exist, error will be returned.
func handleBlockHeaderFetchRequestMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	var hash common.Hash
	if err := msg.Decode(&hash); err != nil {
		return errResp(ErrDecode, "%v: %v", msg, err)
	}

	header := pm.blockchain.GetHeaderByHash(hash)
	if header == nil {
		return fmt.Errorf("peer requested header for non-existing hash. peer: %v, hash: %v", p.GetID(), hash)
	}

	return p.SendFetchedBlockHeader(header)
}

// handleBlockHeaderFetchResponseMsg handles new block header response message.
// This message should contain only one header.
func handleBlockHeaderFetchResponseMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	var header *types.Header
	if err := msg.Decode(&header); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}

	headers := pm.fetcher.FilterHeaders(p.GetID(), []*types.Header{header}, time.Now())
	if len(headers) != 0 {
		logger.Debug("Failed to filter header", "peer", p.GetID(),
			"num", header.Number.Uint64(), "hash", header.Hash(), "len(headers)", len(headers))
	}

	return nil
}

// handleBlockBodiesFetchRequestMsg handles block bodies fetch request message.
// If the peer requests bodies which do not exist, error will be returned.
func handleBlockBodiesFetchRequestMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	if bodies, err := handleBlockBodiesRequest(pm, p, msg); err != nil {
		return err
	} else {
		return p.SendFetchedBlockBodiesRLP(bodies)
	}
}

// handleBlockBodiesFetchResponseMsg handles block bodies fetch response message.
func handleBlockBodiesFetchResponseMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// A batch of block bodies arrived to one of our previous requests
	var request blockBodiesData
	if err := msg.Decode(&request); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// Deliver them all to the downloader for queuing
	transactions := make([][]*types.Transaction, len(request))

	for i, body := range request {
		transactions[i] = body.Transactions
	}

	transactions = pm.fetcher.FilterBodies(p.GetID(), transactions, time.Now())

	if len(transactions) > 0 {
		logger.Warn("Failed to filter bodies", "peer", p.GetID(), "lenTxs", len(transactions))
	}
	return nil
}

// handleNewBlockMsg handles new block message.
func handleNewBlockMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// Retrieve and decode the propagated block
	var request newBlockData
	if err := msg.Decode(&request); err != nil {
		return errResp(ErrDecode, "%v: %v", msg, err)
	}
	request.Block.ReceivedAt = msg.ReceivedAt
	request.Block.ReceivedFrom = p

	// Mark the peer as owning the block and schedule it for import
	p.AddToKnownBlocks(request.Block.Hash())
	pm.fetcher.Enqueue(p.GetID(), request.Block)

	// Assuming the block is importable by the peer, but possibly not yet done so,
	// calculate the head hash and TD that the peer truly must have.
	var (
		trueHead = request.Block.ParentHash()
		trueTD   = new(big.Int).Sub(request.TD, request.Block.BlockScore())
	)
	// Update the peers total blockscore if better than the previous
	if _, td := p.Head(); trueTD.Cmp(td) > 0 {
		p.SetHead(trueHead, trueTD)

		// Schedule a sync if above ours. Note, this will not fire a sync for a gap of
		// a singe block (as the true TD is below the propagated block), however this
		// scenario should easily be covered by the fetcher.
		currentBlock := pm.blockchain.CurrentBlock()
		if trueTD.Cmp(pm.blockchain.GetTd(currentBlock.Hash(), currentBlock.NumberU64())) > 0 {
			go pm.synchronise(p)
		}
	}
	return nil
}

// handleTxMsg handles transaction-propagating message.
func handleTxMsg(pm *ProtocolManager, p Peer, msg p2p.Msg) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if atomic.LoadUint32(&pm.acceptTxs) == 0 {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs types.Transactions
	if err := msg.Decode(&txs); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// Only valid txs should be pushed into the pool.
	validTxs := make(types.Transactions, 0, len(txs))
	var err error
	for i, tx := range txs {
		// Validate and mark the remote transaction
		if tx == nil {
			err = errResp(ErrDecode, "transaction %d is nil", i)
			continue
		}
		p.AddToKnownTxs(tx.Hash())
		validTxs = append(validTxs, tx)
		txReceiveCounter.Inc(1)
	}
	pm.txpool.HandleTxMsg(validTxs)
	return err
}

// sampleSize calculates the number of peers to send block.
// If calcSampleSize is smaller than minNumPeersToSendBlock, it returns minNumPeersToSendBlock.
// Otherwise, it returns calcSampleSize.
func sampleSize(peers []Peer) int {
	if len(peers) < minNumPeersToSendBlock {
		return len(peers)
	}

	calcSampleSize := int(math.Sqrt(float64(len(peers))))
	if calcSampleSize > minNumPeersToSendBlock {
		return calcSampleSize
	} else {
		return minNumPeersToSendBlock
	}
}

// BroadcastBlock will propagate a block to a subset of its peers.
// If current node is CN, it will send block to all PN peers + sampled CN peers without block.
// However, if there are more than 5 PN peers, it will sample 5 PN peers.
// If current node is not CN, it will send block to sampled peers except CNs.
func (pm *ProtocolManager) BroadcastBlock(block *types.Block) {
	if parent := pm.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent == nil {
		logger.Error("Propagating dangling block", "number", block.Number(), "hash", block.Hash())
		return
	}
	// TODO-Klaytn only send all validators + sub(peer) except subset for this block
	//transfer := peers[:int(math.Sqrt(float64(len(peers))))]

	// Calculate the TD of the block (it's not imported yet, so block.Td is not valid)
	td := new(big.Int).Add(block.BlockScore(), pm.blockchain.GetTd(block.ParentHash(), block.NumberU64()-1))
	peersToSendBlock := pm.peers.SamplePeersToSendBlock(block, pm.nodetype)
	for _, peer := range peersToSendBlock {
		peer.AsyncSendNewBlock(block, td)
	}
}

// BroadcastBlockHash will propagate a blockHash to a subset of its peers.
func (pm *ProtocolManager) BroadcastBlockHash(block *types.Block) {
	if !pm.blockchain.HasBlock(block.Hash(), block.NumberU64()) {
		return
	}

	// Otherwise if the block is indeed in out own chain, announce it
	peersWithoutBlock := pm.peers.PeersWithoutBlock(block.Hash())
	for _, peer := range peersWithoutBlock {
		//peer.SendNewBlockHashes([]common.Hash{hash}, []uint64{block.NumberU64()})
		peer.AsyncSendNewBlockHash(block)
	}
	logger.Trace("Announced block", "hash", block.Hash(),
		"recipients", len(peersWithoutBlock), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
}

// BroadcastTxs propagates a batch of transactions to its peers which are not known to
// already have the given transaction.
func (pm *ProtocolManager) BroadcastTxs(txs types.Transactions) {
	// This function calls sendTransaction() to broadcast the transactions for each peer.
	// In that case, transactions are sorted for each peer in sendTransaction().
	// Therefore, it prevents sorting transactions by each peer.
	if !sort.IsSorted(types.TxByPriceAndTime(txs)) {
		sort.Sort(types.TxByPriceAndTime(txs))
	}

	switch pm.nodetype {
	case common.CONSENSUSNODE:
		pm.broadcastTxsFromCN(txs)
	case common.PROXYNODE:
		pm.broadcastTxsFromPN(txs)
	case common.ENDPOINTNODE:
		pm.broadcastTxsFromEN(txs)
	default:
		logger.Error("Unexpected nodeType of ProtocolManager", "nodeType", pm.nodetype)
	}
}

func (pm *ProtocolManager) broadcastTxsFromCN(txs types.Transactions) {
	cnPeersWithoutTxs := make(map[Peer]types.Transactions)
	for _, tx := range txs {
		peers := pm.peers.CNWithoutTx(tx.Hash())
		if len(peers) == 0 {
			logger.Trace("No peer to broadcast transaction", "hash", tx.Hash(), "recipients", len(peers))
			continue
		}

		// TODO-Klaytn Code Check
		//peers = peers[:int(math.Sqrt(float64(len(peers))))]
		half := (len(peers) / 2) + 2
		peers = samplingPeers(peers, half)
		for _, peer := range peers {
			cnPeersWithoutTxs[peer] = append(cnPeersWithoutTxs[peer], tx)
		}
		logger.Trace("Broadcast transaction", "hash", tx.Hash(), "recipients", len(peers))
	}

	propTxPeersGauge.Update(int64(len(cnPeersWithoutTxs)))
	// FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for peer, txs2 := range cnPeersWithoutTxs {
		//peer.SendTransactions(txs)
		peer.AsyncSendTransactions(txs2)
	}
}

func (pm *ProtocolManager) broadcastTxsFromPN(txs types.Transactions) {
	cnPeersWithoutTxs := make(map[Peer]types.Transactions)
	peersWithoutTxs := make(map[Peer]types.Transactions)
	for _, tx := range txs {
		// TODO-Klaytn drop or missing tx
		cnPeers := pm.peers.CNWithoutTx(tx.Hash())
		if len(cnPeers) > 0 {
			cnPeers = samplingPeers(cnPeers, 2) // TODO-Klaytn optimize pickSize or propagation way
			for _, peer := range cnPeers {
				cnPeersWithoutTxs[peer] = append(cnPeersWithoutTxs[peer], tx)
			}
			logger.Trace("Broadcast transaction", "hash", tx.Hash(), "recipients", len(cnPeers))
		}
		pm.peers.UpdateTypePeersWithoutTxs(tx, common.PROXYNODE, peersWithoutTxs)
		txSendCounter.Inc(1)
	}

	propTxPeersGauge.Update(int64(len(peersWithoutTxs) + len(cnPeersWithoutTxs)))
	sendTransactions(cnPeersWithoutTxs)
	sendTransactions(peersWithoutTxs)
}

func (pm *ProtocolManager) broadcastTxsFromEN(txs types.Transactions) {
	peersWithoutTxs := make(map[Peer]types.Transactions)
	for _, tx := range txs {
		pm.peers.UpdateTypePeersWithoutTxs(tx, common.CONSENSUSNODE, peersWithoutTxs)
		pm.peers.UpdateTypePeersWithoutTxs(tx, common.PROXYNODE, peersWithoutTxs)
		pm.peers.UpdateTypePeersWithoutTxs(tx, common.ENDPOINTNODE, peersWithoutTxs)
		txSendCounter.Inc(1)
	}

	propTxPeersGauge.Update(int64(len(peersWithoutTxs)))
	sendTransactions(peersWithoutTxs)
}

// ReBroadcastTxs sends transactions, not considering whether the peer has the transaction or not.
// Only PN and EN rebroadcast transactions to its peers, a CN does not rebroadcast transactions.
func (pm *ProtocolManager) ReBroadcastTxs(txs types.Transactions) {
	// A consensus node does not rebroadcast transactions, hence return here.
	if pm.nodetype == common.CONSENSUSNODE {
		return
	}

	// This function calls sendTransaction() to broadcast the transactions for each peer.
	// In that case, transactions are sorted for each peer in sendTransaction().
	// Therefore, it prevents sorting transactions by each peer.
	if !sort.IsSorted(types.TxByPriceAndTime(txs)) {
		sort.Sort(types.TxByPriceAndTime(txs))
	}

	peersWithoutTxs := make(map[Peer]types.Transactions)
	for _, tx := range txs {
		peers := pm.peers.SampleResendPeersByType(pm.nodetype)
		for _, peer := range peers {
			peersWithoutTxs[peer] = append(peersWithoutTxs[peer], tx)
		}
		txResendCounter.Inc(1)
	}

	propTxPeersGauge.Update(int64(len(peersWithoutTxs)))
	sendTransactions(peersWithoutTxs)
}

// sendTransactions iterates the given map with the key-value pair of Peer and Transactions
// and sends the paired transactions to the peer in synchronised way.
func sendTransactions(txsSet map[Peer]types.Transactions) {
	for peer, txs := range txsSet {
		if err := peer.SendTransactions(txs); err != nil {
			logger.Error("Failed to send txs", "peer", peer.GetAddr(), "peerType", peer.ConnType(), "numTxs", len(txs), "err", err)
		}
	}
}

func samplingPeers(peers []Peer, pickSize int) []Peer {
	if len(peers) <= pickSize {
		return peers
	}

	picker := rand.New(rand.NewSource(time.Now().Unix()))
	peerCount := len(peers)
	for i := 0; i < peerCount; i++ {
		randIndex := picker.Intn(peerCount)
		peers[i], peers[randIndex] = peers[randIndex], peers[i]
	}

	return peers[:pickSize]
}

// Mined broadcast loop
func (pm *ProtocolManager) minedBroadcastLoop() {
	// automatically stops if unsubscribe
	for obj := range pm.minedBlockSub.Chan() {
		switch ev := obj.Data.(type) {
		case blockchain.NewMinedBlockEvent:
			pm.BroadcastBlock(ev.Block)     // First propagate block to peers
			pm.BroadcastBlockHash(ev.Block) // Only then announce to the rest
		}
	}
}

func (pm *ProtocolManager) txBroadcastLoop() {
	for {
		select {
		case event := <-pm.txsCh:
			pm.BroadcastTxs(event.Txs)
			// Err() channel will be closed when unsubscribing.
		case <-pm.txsSub.Err():
			return
		}
	}
}

func (pm *ProtocolManager) txResendLoop(period uint64, maxTxCount int) {
	tick := time.Duration(period) * time.Second
	resend := time.NewTicker(tick)
	defer resend.Stop()

	logger.Debug("txResendloop started", "period", tick.Seconds())

	for {
		select {
		case <-resend.C:
			pending := pm.txpool.CachedPendingTxsByCount(maxTxCount)
			pm.txResend(pending)
		case <-pm.quitResendCh:
			logger.Debug("txResendloop stopped")
			return
		}
	}
}

func (pm *ProtocolManager) txResend(pending types.Transactions) {
	txResendRoutineGauge.Update(txResendRoutineGauge.Value() + 1)
	defer txResendRoutineGauge.Update(txResendRoutineGauge.Value() - 1)
	// TODO-Klaytn drop or missing tx
	if len(pending) > 0 {
		logger.Debug("Tx Resend", "count", len(pending))
		pm.ReBroadcastTxs(pending)
	}
}

func (pm *ProtocolManager) useTxResend() bool {
	if pm.nodetype != common.CONSENSUSNODE && !pm.txResendUseLegacy {
		return true
	}
	return false
}

// NodeInfo represents a short summary of the Klaytn sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	// TODO-Klaytn describe predefined network ID below
	Network    uint64              `json:"network"`    // Klaytn network ID
	BlockScore *big.Int            `json:"blockscore"` // Total blockscore of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // SHA3 hash of the host's best owned block
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (pm *ProtocolManager) NodeInfo() *NodeInfo {
	currentBlock := pm.blockchain.CurrentBlock()
	return &NodeInfo{
		Network:    pm.networkId,
		BlockScore: pm.blockchain.GetTd(currentBlock.Hash(), currentBlock.NumberU64()),
		Genesis:    pm.blockchain.Genesis().Hash(),
		Config:     pm.blockchain.Config(),
		Head:       currentBlock.Hash(),
	}
}

// Below functions are used in Istanbul BFT consensus.
// Enqueue wraps fetcher's Enqueue function to insert the given block.
func (pm *ProtocolManager) Enqueue(id string, block *types.Block) {
	pm.fetcher.Enqueue(id, block)
}

func (pm *ProtocolManager) FindPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	m := make(map[common.Address]consensus.Peer)
	for _, p := range pm.peers.Peers() {
		addr := p.GetAddr()
		if addr == (common.Address{}) {
			pubKey, err := p.GetP2PPeerID().Pubkey()
			if err != nil {
				continue
			}
			addr = crypto.PubkeyToAddress(*pubKey)
			p.SetAddr(addr)
		}
		if targets[addr] {
			m[addr] = p
		}
	}
	return m
}

func (pm *ProtocolManager) GetCNPeers() map[common.Address]consensus.Peer {
	m := make(map[common.Address]consensus.Peer)
	for addr, p := range pm.peers.CNPeers() {
		m[addr] = p
	}
	return m
}

func (pm *ProtocolManager) FindCNPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	m := make(map[common.Address]consensus.Peer)
	for addr, p := range pm.peers.CNPeers() {
		if targets[addr] {
			m[addr] = p
		}
	}
	return m
}

func (pm *ProtocolManager) GetENPeers() map[common.Address]consensus.Peer {
	m := make(map[common.Address]consensus.Peer)
	for addr, p := range pm.peers.ENPeers() {
		m[addr] = p
	}
	return m
}

func (pm *ProtocolManager) GetPeers() []common.Address {
	addrs := make([]common.Address, 0)
	for _, p := range pm.peers.Peers() {
		addr := p.GetAddr()
		if addr == (common.Address{}) {
			pubKey, err := p.GetP2PPeerID().Pubkey()
			if err != nil {
				continue
			}
			addr = crypto.PubkeyToAddress(*pubKey)
			p.SetAddr(addr)
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

func (pm *ProtocolManager) Downloader() ProtocolManagerDownloader {
	return pm.downloader
}

func (pm *ProtocolManager) SetWsEndPoint(wsep string) {
	pm.wsendpoint = wsep
}

func (pm *ProtocolManager) GetSubProtocols() []p2p.Protocol {
	return pm.SubProtocols
}

func (pm *ProtocolManager) ProtocolVersion() int {
	return int(pm.SubProtocols[0].Version)
}

func (pm *ProtocolManager) SetAcceptTxs() {
	atomic.StoreUint32(&pm.acceptTxs, 1)
}

func (pm *ProtocolManager) NodeType() common.ConnType {
	return pm.nodetype
}
