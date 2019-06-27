// Modifications Copyright 2019 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from eth/backend.go (2018/06/04).
// Modified and improved for the klaytn development.

package sc

import (
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/sc/bridgepool"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/work"
	"math/big"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

const (
	forceSyncCycle = 10 * time.Second // Time interval to force syncs, even if few peers are available

	tokenReceivedChanSize = 10000
	tokenTransferChanSize = 10000
)

// Backend wraps all methods for local and remote backend
type Backend interface {
	bind.ContractBackend
}

// NodeInfo represents a short summary of the ServiceChain sub-protocol metadata
// known about the host peer.
type SubBridgeInfo struct {
	Network uint64              `json:"network"` // Klaytn network ID
	Genesis common.Hash         `json:"genesis"` // SHA3 hash of the host's genesis block
	Config  *params.ChainConfig `json:"config"`  // Chain configuration for the fork rules
	Head    common.Hash         `json:"head"`    // SHA3 hash of the host's best owned block
	ChainID *big.Int            `json:"chainid"` // ChainID
}

// SubBridge implements the Klaytn consensus node service.
type SubBridge struct {
	config *SCConfig

	// DB interfaces
	chainDB database.DBManager // Block chain database

	eventMux       *event.TypeMux
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *api.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and coinbase)

	bridgeServer p2p.Server
	ctx          *node.ServiceContext
	maxPeers     int

	APIBackend *SubBridgeAPI

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan BridgePeer
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg   sync.WaitGroup
	pmwg sync.WaitGroup

	blockchain   *blockchain.BlockChain
	txPool       *blockchain.TxPool
	bridgeTxPool *bridgepool.BridgeTxPool

	// chain event
	chainHeadCh  chan blockchain.ChainHeadEvent
	chainHeadSub event.Subscription
	logsCh       chan []*types.Log
	logsSub      event.Subscription

	// If this channel can't be read immediately, it can lock service chain tx pool.
	// Commented out because for now, it doesn't need.
	//txCh         chan blockchain.NewTxsEvent
	//txSub        event.Subscription

	peers        *bridgePeerSet
	handler      *SubBridgeHandler
	eventhandler *ChildChainEventHandler

	// bridgemanager for value exchange
	localBackend  Backend
	remoteBackend Backend
	bridgeManager *BridgeManager

	requestEventCh  chan RequestValueTransferEvent
	requestEventSub event.Subscription
	handleEventCh   chan HandleValueTransferEvent
	handleEventSub  event.Subscription

	bridgeAccountManager *BridgeAccountManager

	bootFail bool

	// service on/off
	onAnchoringTx bool

	checkConnection int64
}

// New creates a new CN object (including the
// initialisation of the common CN object)
func NewSubBridge(ctx *node.ServiceContext, config *SCConfig) (*SubBridge, error) {
	chainDB := CreateDB(ctx, config, "subbridgedata")

	// initialize private keys for servicechain
	config.chainkey = config.ChainKey()
	config.nodekey = config.NodeKey()

	sc := &SubBridge{
		config:         config,
		chainDB:        chainDB,
		peers:          newBridgePeerSet(),
		newPeerCh:      make(chan BridgePeer),
		noMorePeers:    make(chan struct{}),
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		networkId:      config.NetworkId,
		ctx:            ctx,
		chainHeadCh:    make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
		logsCh:         make(chan []*types.Log, chainLogChanSize),
		//txCh:            make(chan blockchain.NewTxsEvent, transactionChanSize),
		quitSync:       make(chan struct{}),
		maxPeers:       config.MaxPeer,
		requestEventCh: make(chan RequestValueTransferEvent, tokenReceivedChanSize),
		handleEventCh:  make(chan HandleValueTransferEvent, tokenTransferChanSize),
		onAnchoringTx:  false,
		bootFail:       false,
	}
	// TODO-Klaytn change static config to user define config
	bridgetxConfig := bridgepool.BridgeTxPoolConfig{
		ParentChainID: new(big.Int).SetUint64(config.ParentChainID),
		Journal:       path.Join(config.DataDir, "bridge_transactions.rlp"),
		Rejournal:     time.Hour,
		GlobalQueue:   8192,
	}

	logger.Info("Initialising Klaytn-Bridge protocol", "network", config.NetworkId)

	bcVersion := chainDB.ReadDatabaseVersion()
	if bcVersion != blockchain.BlockChainVersion && bcVersion != 0 {
		return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d).\n", bcVersion, blockchain.BlockChainVersion)
	}
	chainDB.WriteDatabaseVersion(blockchain.BlockChainVersion)

	sc.APIBackend = &SubBridgeAPI{sc}

	sc.bridgeTxPool = bridgepool.NewBridgeTxPool(bridgetxConfig)

	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		return nil, err
	}
	sc.eventhandler, err = NewChildChainEventHandler(sc, sc.handler)
	if err != nil {
		return nil, err
	}
	sc.bridgeAccountManager, err = NewBridgeAccountManager(sc.config.chainkey, sc.config.nodekey)
	if err != nil {
		return nil, err
	}
	sc.bridgeAccountManager.mcAccount.SetChainID(new(big.Int).SetUint64(config.ParentChainID))

	return sc, nil
}

// implement PeerSetManager
func (sb *SubBridge) BridgePeerSet() *bridgePeerSet {
	return sb.peers
}

func (sb *SubBridge) GetBridgeTxPool() *bridgepool.BridgeTxPool {
	return sb.bridgeTxPool
}

func (sb *SubBridge) GetAnchoringTx() bool {
	return sb.onAnchoringTx
}

func (sb *SubBridge) SetAnchoringTx(flag bool) bool {
	sb.onAnchoringTx = flag
	return sb.GetAnchoringTx()
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *SubBridge) APIs() []rpc.API {
	// Append all the local APIs and return
	return []rpc.API{
		{
			Namespace: "subbridge",
			Version:   "1.0",
			Service:   s.APIBackend,
			Public:    true,
		},
		{
			Namespace: "subbridge",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}
}

func (s *SubBridge) AccountManager() *accounts.Manager { return s.accountManager }
func (s *SubBridge) EventMux() *event.TypeMux          { return s.eventMux }
func (s *SubBridge) ChainDB() database.DBManager       { return s.chainDB }
func (s *SubBridge) IsListening() bool                 { return true } // Always listening
func (s *SubBridge) ProtocolVersion() int              { return int(s.SCProtocol().Versions[0]) }
func (s *SubBridge) NetVersion() uint64                { return s.networkId }

func (s *SubBridge) Components() []interface{} {
	return nil
}

func (sc *SubBridge) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			sc.blockchain = v
			// event from core-service
			sc.chainHeadSub = sc.blockchain.SubscribeChainHeadEvent(sc.chainHeadCh)
			sc.logsSub = sc.blockchain.SubscribeLogsEvent(sc.logsCh)
		case *blockchain.TxPool:
			sc.txPool = v
			// event from core-service
			// sc.txSub = sc.txPool.SubscribeNewTxsEvent(sc.txCh)
			// TODO-Klaytn if need pending block, should use miner
		case *work.Miner:
		}
	}

	var err error
	if sc.config.MainChainURL != "" {
		sc.remoteBackend, err = NewRemoteBackend(sc, sc.config.MainChainURL)
		if err != nil {
			logger.Error("fail to initialize RemoteBackend", "err", err, "url", sc.config.MainChainURL)
			sc.bootFail = true
			return
		}
	}
	sc.localBackend, err = NewLocalBackend(sc)
	if err != nil {
		logger.Error("fail to initialize LocalBackend", "err", err)
		sc.bootFail = true
		return
	}

	sc.bridgeManager, err = NewBridgeManager(sc)
	if err != nil {
		logger.Error("fail to initialize BridgeManager", "err", err)
		sc.bootFail = true
		return
	}
	sc.requestEventSub = sc.bridgeManager.SubscribeTokenReceived(sc.requestEventCh)
	sc.handleEventSub = sc.bridgeManager.SubscribeTokenWithDraw(sc.handleEventCh)

	if err := sc.bridgeManager.RestoreBridges(); err != nil {
		logger.Error("failed to sc.bridgeManager.RestoreBridges()", "err", err)
	}

	sc.bridgeAccountManager.scAccount.SetNonce(sc.txPool.GetPendingNonce(sc.bridgeAccountManager.scAccount.address))

	sc.pmwg.Add(1)
	go sc.loop()
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *SubBridge) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (s *SubBridge) SCProtocol() SCProtocol {
	return SCProtocol{
		Name:     SCProtocolName,
		Versions: SCProtocolVersion,
		Lengths:  SCProtocolLength,
	}
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (pm *SubBridge) NodeInfo() *SubBridgeInfo {
	currentBlock := pm.blockchain.CurrentBlock()
	return &SubBridgeInfo{
		Network: pm.networkId,
		Genesis: pm.blockchain.Genesis().Hash(),
		Config:  pm.blockchain.Config(),
		Head:    currentBlock.Hash(),
		ChainID: pm.blockchain.Config().ChainID,
	}
}

// getChainID returns the current chain id.
func (pm *SubBridge) getChainID() *big.Int {
	return pm.blockchain.Config().ChainID
}

// Start implements node.Service, starting all internal goroutines needed by the
// Klaytn protocol implementation.
func (s *SubBridge) Start(srvr p2p.Server) error {

	if s.bootFail {
		return errors.New("subBridge node fail to start")
	}

	serverConfig := p2p.Config{}
	serverConfig.PrivateKey = s.ctx.NodeKey()
	serverConfig.Name = s.ctx.NodeType().String()
	serverConfig.Logger = logger
	serverConfig.NoListen = true
	serverConfig.MaxPhysicalConnections = s.maxPeers
	serverConfig.NoDiscovery = true
	serverConfig.EnableMultiChannelServer = false

	// connect to mainbridge as outbound
	serverConfig.StaticNodes = s.config.MainBridges()

	p2pServer := p2p.NewServer(serverConfig)

	s.bridgeServer = p2pServer

	scprotocols := make([]p2p.Protocol, 0, len(s.SCProtocol().Versions))
	for i, version := range s.SCProtocol().Versions {
		// Compatible; initialise the sub-protocol
		version := version
		scprotocols = append(scprotocols, p2p.Protocol{
			Name:    s.SCProtocol().Name,
			Version: version,
			Length:  s.SCProtocol().Lengths[i],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := s.newPeer(int(version), p, rw)
				pubKey, _ := p.ID().Pubkey()
				addr := crypto.PubkeyToAddress(*pubKey)
				peer.SetAddr(addr)
				select {
				case s.newPeerCh <- peer:
					s.wg.Add(1)
					defer s.wg.Done()
					return s.handle(peer)
				case <-s.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return s.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := s.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}
	s.bridgeServer.AddProtocols(scprotocols)

	if err := p2pServer.Start(); err != nil {
		return errors.New("fail to bridgeserver start")
	}

	// Start the RPC service
	s.netRPCService = api.NewPublicNetAPI(s.bridgeServer, s.NetVersion())

	// Figure out a max peers count based on the server limits
	//s.maxPeers = s.bridgeServer.MaxPhysicalConnections()
	//validator := func(header *types.Header) error {
	//	return nil
	//}
	//heighter := func() uint64 {
	//	return s.blockchain.CurrentBlock().NumberU64()
	//}
	//inserter := func(blocks types.Blocks) (int, error) {
	//	return 0, nil
	//}
	//s.fetcher = fetcher.New(s.GetBlockByHash, validator, s.BroadcastBlock, heighter, inserter, s.removePeer)

	go s.syncer()

	return nil
}

func (pm *SubBridge) newPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter) BridgePeer {
	return newBridgePeer(pv, p, newMeteredMsgWriter(rw))
}

func (pm *SubBridge) handle(p BridgePeer) error {
	// Ignore maxPeers if this is a trusted peer
	if pm.peers.Len() >= pm.maxPeers && !p.GetP2PPeer().Info().Networks[p2p.ConnDefault].Trusted {
		return p2p.DiscTooManyPeers
	}
	p.GetP2PPeer().Log().Debug("Klaytn peer connected", "name", p.GetP2PPeer().Name())

	// Execute the handshake
	var (
		head   = pm.blockchain.CurrentHeader()
		hash   = head.Hash()
		number = head.Number.Uint64()
		td     = pm.blockchain.GetTd(hash, number)
	)

	err := p.Handshake(pm.networkId, pm.getChainID(), td, hash)
	if err != nil {
		p.GetP2PPeer().Log().Debug("Klaytn peer handshake failed", "err", err)
		fmt.Println(err)
		return err
	}

	// Register the peer locally
	if err := pm.peers.Register(p); err != nil {
		// if starting node with unlock account, can't register peer until finish unlock
		p.GetP2PPeer().Log().Info("Klaytn peer registration failed", "err", err)
		fmt.Println(err)
		return err
	}
	defer pm.removePeer(p.GetID())

	pm.handler.RegisterNewPeer(p)

	p.GetP2PPeer().Log().Info("Added a P2P Peer", "peerID", p.GetP2PPeerID())

	//pubKey, err := p.GetP2PPeerID().Pubkey()
	//if err != nil {
	//	return err
	//}
	//addr := crypto.PubkeyToAddress(*pubKey)

	// main loop. handle incoming messages.
	for {
		if err := pm.handleMsg(p); err != nil {
			p.GetP2PPeer().Log().Debug("Klaytn message handling failed", "err", err)

			if pm.peers.Len() == 1 {
				pm.handler.setMainChainAccountNonceSynced(false)
				atomic.StoreInt64(&pm.checkConnection, 1)
			}
			return err
		}
	}
}

func (sc *SubBridge) loop() {
	defer sc.pmwg.Done()

	report := time.NewTicker(1 * time.Second)
	defer report.Stop()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle ChainHeadEvent
		case ev := <-sc.chainHeadCh:
			if ev.Block != nil {
				if err := sc.eventhandler.HandleChainHeadEvent(ev.Block); err != nil {
					logger.Error("subbridge block event", "err", err)
				}
			} else {
				logger.Error("subbridge block event is nil")
			}
		// Handle NewTexsEvent
		//case ev := <-sc.txCh:
		//	if ev.Txs != nil {
		//		if err := sc.eventhandler.HandleTxsEvent(ev.Txs); err != nil {
		//			logger.Error("subbridge tx event", "err", err)
		//		}
		//	} else {
		//		logger.Error("subbridge tx event is nil")
		//	}
		// Handle ChainLogsEvent
		case logs := <-sc.logsCh:
			if err := sc.eventhandler.HandleLogsEvent(logs); err != nil {
				logger.Error("subbridge log event", "err", err)
			}
		// Handle Bridge Event
		case ev := <-sc.requestEventCh:
			vtRequestEventMeter.Mark(1)
			if err := sc.eventhandler.ProcessRequestEvent(ev); err != nil {
				logger.Error("fail to process request value transfer event ", "err", err)
			}
		case ev := <-sc.handleEventCh:
			vtHandleEventMeter.Mark(1)
			if err := sc.eventhandler.ProcessHandleEvent(ev); err != nil {
				logger.Error("fail to process handle value transfer event ", "err", err)
			}
		case <-report.C:
			// report status
		case err := <-sc.chainHeadSub.Err():
			if err != nil {
				logger.Error("subbridge block subscription ", "err", err)
			}
			return
		//case err := <-sc.txSub.Err():
		//	if err != nil {
		//		logger.Error("subbridge tx subscription ", "err", err)
		//	}
		//	return
		case err := <-sc.logsSub.Err():
			if err != nil {
				logger.Error("subbridge log subscription ", "err", err)
			}
			return
		case err := <-sc.requestEventSub.Err():
			if err != nil {
				logger.Error("subbridge token-received subscription ", "err", err)
			}
			return
		case err := <-sc.handleEventSub.Err():
			if err != nil {
				logger.Error("subbridge token-transfer subscription ", "err", err)
			}
			return
		}
	}
}

func (pm *SubBridge) removePeer(id string) {
	// Short circuit if the peer was already removed
	peer := pm.peers.Peer(id)
	if peer == nil {
		return
	}
	logger.Debug("Removing Klaytn peer", "peer", id)

	if err := pm.peers.Unregister(id); err != nil {
		logger.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.GetP2PPeer().Disconnect(p2p.DiscUselessPeer)
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (pm *SubBridge) handleMsg(p BridgePeer) error {
	//Below message size checking is done by handle().
	//Read the next message from the remote peer, and ensure it's fully consumed
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
	defer msg.Discard()

	return pm.handler.HandleMainMsg(p, msg)
}

func (pm *SubBridge) syncer() {
	// Start and ensure cleanup of sync mechanisms
	//pm.fetcher.Start()
	//defer pm.fetcher.Stop()
	//defer pm.downloader.Terminate()

	// Wait for different events to fire synchronisation operations
	forceSync := time.NewTicker(forceSyncCycle)
	defer forceSync.Stop()

	for {
		select {
		case peer := <-pm.newPeerCh:
			go pm.synchronise(peer)

		case <-forceSync.C:
			// Force a sync even if not enough peers are present
			go pm.synchronise(pm.peers.BestPeer())

		case <-pm.noMorePeers:
			return
		}
	}
}

func (pm *SubBridge) synchronise(peer BridgePeer) {
	// @TODO Klaytn ServiceChain Sync
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Klaytn protocol.
func (s *SubBridge) Stop() error {

	close(s.quitSync)

	s.chainHeadSub.Unsubscribe()
	//s.txSub.Unsubscribe()
	s.logsSub.Unsubscribe()
	s.requestEventSub.Unsubscribe()
	s.handleEventSub.Unsubscribe()
	s.eventMux.Stop()
	s.chainDB.Close()

	s.bridgeManager.Stop()
	s.bridgeTxPool.Stop()
	s.bridgeServer.Stop()

	return nil
}
