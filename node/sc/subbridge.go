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
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"path"
	"sync"
	"time"

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
	"github.com/klaytn/klaytn/node/sc/kas"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/work"
)

const (
	forceSyncCycle = 10 * time.Second // Time interval to force syncs, even if few peers are available

	chanReqVTevanSize    = 10000
	chanHandleVTevanSize = 10000

	resetBridgeCycle   = 3 * time.Second
	restoreBridgeCycle = 3 * time.Second
)

// RemoteBackendInterface wraps methods for remote backend
type RemoteBackendInterface interface {
	bind.ContractBackend
	TransactionReceiptRpcOutput(ctx context.Context, txHash common.Hash) (map[string]interface{}, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

// Backend wraps all methods for local and remote backend
type Backend interface {
	bind.ContractBackend
	CurrentBlockNumber(context.Context) (uint64, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
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

//go:generate mockgen -destination=bridgeTxPool_mock_test.go -package=sc github.com/klaytn/klaytn/node/sc BridgeTxPool
type BridgeTxPool interface {
	GetMaxTxNonce(from *common.Address) uint64
	AddLocal(tx *types.Transaction) error
	Stats() int
	Pending() map[common.Address]types.Transactions
	Get(hash common.Hash) *types.Transaction
	RemoveTx(tx *types.Transaction) error
	PendingTxHashesByAddress(from *common.Address, limit int) []common.Hash
	PendingTxsByAddress(from *common.Address, limit int) types.Transactions
	Stop()
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
	newPeerCh    chan BridgePeer
	addPeerCh    chan struct{}
	noMorePeers  chan struct{}
	removePeerCh chan struct{}
	quitSync     chan struct{}

	// wait group is used for graceful shutdowns during downloading and processing
	pmwg sync.WaitGroup

	blockchain   *blockchain.BlockChain
	txPool       *blockchain.TxPool
	bridgeTxPool BridgeTxPool

	// chain event
	chainCh  chan blockchain.ChainEvent
	chainSub event.Subscription
	logsCh   chan []*types.Log
	logsSub  event.Subscription

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

	chanReqVTev        chan RequestValueTransferEvent
	chanReqVTencodedEv chan RequestValueTransferEncodedEvent
	reqVTevSub         event.Subscription
	reqVTencodedEvSub  event.Subscription
	chanHandleVTev     chan *HandleValueTransferEvent
	handleVTevSub      event.Subscription

	bridgeAccounts *BridgeAccounts

	bootFail bool

	// service on/off
	onAnchoringTx bool

	rpcConn   net.Conn
	rpcSendCh chan []byte

	//KAS Anchor
	kasAnchor *kas.Anchor
}

// New creates a new CN object (including the
// initialisation of the common CN object)
func NewSubBridge(ctx *node.ServiceContext, config *SCConfig) (*SubBridge, error) {
	chainDB := CreateDB(ctx, config, "subbridgedata")

	sb := &SubBridge{
		config:         config,
		chainDB:        chainDB,
		peers:          newBridgePeerSet(),
		newPeerCh:      make(chan BridgePeer),
		addPeerCh:      make(chan struct{}),
		removePeerCh:   make(chan struct{}),
		noMorePeers:    make(chan struct{}),
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		networkId:      config.NetworkId,
		ctx:            ctx,
		chainCh:        make(chan blockchain.ChainEvent, chainEventChanSize),
		logsCh:         make(chan []*types.Log, chainLogChanSize),
		// txCh:            make(chan blockchain.NewTxsEvent, transactionChanSize),
		chanReqVTev:        make(chan RequestValueTransferEvent, chanReqVTevanSize),
		chanReqVTencodedEv: make(chan RequestValueTransferEncodedEvent, chanReqVTevanSize),
		chanHandleVTev:     make(chan *HandleValueTransferEvent, chanHandleVTevanSize),
		quitSync:           make(chan struct{}),
		maxPeers:           config.MaxPeer,
		onAnchoringTx:      config.Anchoring,
		bootFail:           false,
		rpcSendCh:          make(chan []byte),
	}
	// TODO-Klaytn change static config to user define config
	bridgetxConfig := bridgepool.BridgeTxPoolConfig{
		ParentChainID: new(big.Int).SetUint64(config.ParentChainID),
		Journal:       path.Join(config.DataDir, "bridge_transactions.rlp"),
		Rejournal:     time.Hour,
		GlobalQueue:   8192,
	}

	logger.Info("Initialising Klaytn-Bridge protocol", "network", config.NetworkId)
	sb.APIBackend = &SubBridgeAPI{sb}

	sb.bridgeTxPool = bridgepool.NewBridgeTxPool(bridgetxConfig)

	var err error
	sb.bridgeAccounts, err = NewBridgeAccounts(sb.accountManager, config.DataDir, chainDB, sb.config.ServiceChainParentOperatorGasLimit, sb.config.ServiceChainChildOperatorGasLimit)
	if err != nil {
		return nil, err
	}
	sb.handler, err = NewSubBridgeHandler(sb)
	if err != nil {
		return nil, err
	}
	sb.eventhandler, err = NewChildChainEventHandler(sb, sb.handler)
	if err != nil {
		return nil, err
	}
	sb.bridgeAccounts.pAccount.SetChainID(new(big.Int).SetUint64(config.ParentChainID))

	return sb, nil
}

func (sb *SubBridge) SetRPCConn(conn net.Conn) {
	sb.rpcConn = conn

	go func() {
		for {
			data := make([]byte, rpcBufferSize)
			rlen, err := sb.rpcConn.Read(data)
			if err != nil {
				if err == io.EOF {
					logger.Trace("EOF from the rpc pipe")
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					// If no one closes the pipe, this situation should not happen.
					logger.Error("failed to read from the rpc pipe", "err", err, "rlen", rlen)
					return
				}
			}
			sb.rpcSendCh <- data[:rlen]
		}
	}()
}

func (sb *SubBridge) SendRPCData(data []byte) error {
	peers := sb.BridgePeerSet().peers
	logger.Trace("send rpc message from the subbridge", "len", len(data), "peers", len(peers))
	for _, peer := range peers {
		err := peer.SendRequestRPC(data)
		if err != nil {
			logger.Error("SendRPCData Error", "err", err)
		}
		return err
	}
	logger.Trace("send rpc message from the subbridge, done")

	return nil
}

// implement PeerSetManager
func (sb *SubBridge) BridgePeerSet() *bridgePeerSet {
	return sb.peers
}

func (sb *SubBridge) GetBridgeTxPool() BridgeTxPool {
	return sb.bridgeTxPool
}

func (sb *SubBridge) GetAnchoringTx() bool {
	return sb.onAnchoringTx
}

func (sb *SubBridge) SetAnchoringTx(flag bool) bool {
	if sb.onAnchoringTx != flag && flag {
		sb.handler.txCountStartingBlockNumber = 0
	}
	sb.onAnchoringTx = flag
	return sb.GetAnchoringTx()
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (sb *SubBridge) APIs() []rpc.API {
	// Append all the local APIs and return
	return []rpc.API{
		{
			Namespace: "subbridge",
			Version:   "1.0",
			Service:   sb.APIBackend,
			Public:    true,
		},
		{
			Namespace: "subbridge",
			Version:   "1.0",
			Service:   sb.netRPCService,
			Public:    true,
		},
	}
}

func (sb *SubBridge) AccountManager() *accounts.Manager { return sb.accountManager }
func (sb *SubBridge) EventMux() *event.TypeMux          { return sb.eventMux }
func (sb *SubBridge) ChainDB() database.DBManager       { return sb.chainDB }
func (sb *SubBridge) IsListening() bool                 { return true } // Always listening
func (sb *SubBridge) ProtocolVersion() int              { return int(sb.SCProtocol().Versions[0]) }
func (sb *SubBridge) NetVersion() uint64                { return sb.networkId }

func (sb *SubBridge) Components() []interface{} {
	return nil
}

func (sb *SubBridge) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			sb.blockchain = v

			kasConfig := &kas.KASConfig{
				Url:            sb.config.KASAnchorUrl,
				XChainId:       sb.config.KASXChainId,
				User:           sb.config.KASAccessKey,
				Pwd:            sb.config.KASSecretKey,
				Operator:       common.HexToAddress(sb.config.KASAnchorOperator),
				Anchor:         sb.config.KASAnchor,
				AnchorPeriod:   sb.config.KASAnchorPeriod,
				RequestTimeout: sb.config.KASAnchorRequestTimeout,
			}
			sb.kasAnchor = kas.NewKASAnchor(kasConfig, sb.chainDB, v)

			// event from core-service
			sb.chainSub = sb.blockchain.SubscribeChainEvent(sb.chainCh)
			sb.logsSub = sb.blockchain.SubscribeLogsEvent(sb.logsCh)
			sb.bridgeAccounts.cAccount.SetChainID(v.Config().ChainID)
		case *blockchain.TxPool:
			sb.txPool = v
			// event from core-service
			// sb.txSub = sb.txPool.SubscribeNewTxsEvent(sb.txCh)
		// TODO-Klaytn if need pending block, should use miner
		case *work.Miner:
		}
	}

	var err error
	if sb.config.EnabledSubBridge {
		sb.remoteBackend, err = NewRemoteBackend(sb)
		if err != nil {
			logger.Error("fail to initialize RemoteBackend", "err", err)
			sb.bootFail = true
			return
		}
	}
	sb.localBackend, err = NewLocalBackend(sb)
	if err != nil {
		logger.Error("fail to initialize LocalBackend", "err", err)
		sb.bootFail = true
		return
	}

	sb.bridgeManager, err = NewBridgeManager(sb)
	if err != nil {
		logger.Error("fail to initialize BridgeManager", "err", err)
		sb.bootFail = true
		return
	}
	sb.reqVTevSub = sb.bridgeManager.SubscribeReqVTev(sb.chanReqVTev)
	sb.reqVTencodedEvSub = sb.bridgeManager.SubscribeReqVTencodedEv(sb.chanReqVTencodedEv)
	sb.handleVTevSub = sb.bridgeManager.SubscribeHandleVTev(sb.chanHandleVTev)

	sb.pmwg.Add(1)
	go sb.restoreBridgeLoop()

	sb.pmwg.Add(1)
	go sb.resetBridgeLoop()

	sb.bridgeAccounts.cAccount.SetNonce(sb.txPool.GetPendingNonce(sb.bridgeAccounts.cAccount.address))

	sb.pmwg.Add(1)
	go sb.loop()
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (sb *SubBridge) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (sb *SubBridge) SCProtocol() SCProtocol {
	return SCProtocol{
		Name:     SCProtocolName,
		Versions: SCProtocolVersion,
		Lengths:  SCProtocolLength,
	}
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (sb *SubBridge) NodeInfo() *SubBridgeInfo {
	currentBlock := sb.blockchain.CurrentBlock()
	return &SubBridgeInfo{
		Network: sb.networkId,
		Genesis: sb.blockchain.Genesis().Hash(),
		Config:  sb.blockchain.Config(),
		Head:    currentBlock.Hash(),
		ChainID: sb.blockchain.Config().ChainID,
	}
}

// getChainID returns the current chain id.
func (sb *SubBridge) getChainID() *big.Int {
	return sb.blockchain.Config().ChainID
}

// Start implements node.Service, starting all internal goroutines needed by the
// Klaytn protocol implementation.
func (sb *SubBridge) Start(srvr p2p.Server) error {

	if sb.bootFail {
		return errors.New("subBridge node fail to start")
	}

	serverConfig := p2p.Config{}
	serverConfig.PrivateKey = sb.ctx.NodeKey()
	serverConfig.Name = sb.ctx.NodeType().String()
	serverConfig.Logger = logger
	serverConfig.NoListen = true
	serverConfig.MaxPhysicalConnections = sb.maxPeers
	serverConfig.NoDiscovery = true
	serverConfig.EnableMultiChannelServer = false

	// connect to mainbridge as outbound
	serverConfig.StaticNodes = sb.config.MainBridges()

	p2pServer := p2p.NewServer(serverConfig)

	sb.bridgeServer = p2pServer

	scprotocols := make([]p2p.Protocol, 0, len(sb.SCProtocol().Versions))
	for i, version := range sb.SCProtocol().Versions {
		// Compatible; initialise the sub-protocol
		version := version
		scprotocols = append(scprotocols, p2p.Protocol{
			Name:    sb.SCProtocol().Name,
			Version: version,
			Length:  sb.SCProtocol().Lengths[i],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := sb.newPeer(int(version), p, rw)
				pubKey, _ := p.ID().Pubkey()
				addr := crypto.PubkeyToAddress(*pubKey)
				peer.SetAddr(addr)
				select {
				case sb.newPeerCh <- peer:
					return sb.handle(peer)
				case <-sb.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return sb.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := sb.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}
	sb.bridgeServer.AddProtocols(scprotocols)

	if err := p2pServer.Start(); err != nil {
		return errors.New("fail to bridgeserver start")
	}

	// Start the RPC service
	sb.netRPCService = api.NewPublicNetAPI(sb.bridgeServer, sb.NetVersion())

	// Figure out a max peers count based on the server limits
	//sb.maxPeers = sb.bridgeServer.MaxPhysicalConnections()
	//validator := func(header *types.Header) error {
	//	return nil
	//}
	//heighter := func() uint64 {
	//	return sb.blockchain.CurrentBlock().NumberU64()
	//}
	//inserter := func(blocks types.Blocks) (int, error) {
	//	return 0, nil
	//}
	//sb.fetcher = fetcher.New(sb.GetBlockByHash, validator, sb.BroadcastBlock, heighter, inserter, sb.removePeer)

	go sb.syncer()

	return nil
}

func (sb *SubBridge) newPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter) BridgePeer {
	return newBridgePeer(pv, p, newMeteredMsgWriter(rw))
}

func (sb *SubBridge) handle(p BridgePeer) error {
	// Ignore maxPeers if this is a trusted peer
	if sb.peers.Len() >= sb.maxPeers && !p.GetP2PPeer().Info().Networks[p2p.ConnDefault].Trusted {
		return p2p.DiscTooManyPeers
	}
	p.GetP2PPeer().Log().Debug("Klaytn peer connected", "name", p.GetP2PPeer().Name())

	// Execute the handshake
	var (
		head   = sb.blockchain.CurrentHeader()
		hash   = head.Hash()
		number = head.Number.Uint64()
		td     = sb.blockchain.GetTd(hash, number)
	)

	err := p.Handshake(sb.networkId, sb.getChainID(), td, hash)
	if err != nil {
		p.GetP2PPeer().Log().Debug("Klaytn peer handshake failed", "err", err)
		fmt.Println(err)
		return err
	}

	// Register the peer locally
	if err := sb.peers.Register(p); err != nil {
		// if starting node with unlock account, can't register peer until finish unlock
		p.GetP2PPeer().Log().Info("Klaytn peer registration failed", "err", err)
		fmt.Println(err)
		return err
	}
	defer sb.removePeer(p.GetID())

	sb.handler.RegisterNewPeer(p)

	p.GetP2PPeer().Log().Info("Added a P2P Peer", "peerID", p.GetP2PPeerID())

	// main loop. handle incoming messages.
	for {
		if err := sb.handleMsg(p); err != nil {
			p.GetP2PPeer().Log().Debug("Klaytn message handling failed", "err", err)
			return err
		}
	}
}

func (sb *SubBridge) resetBridgeLoop() {
	defer sb.pmwg.Done()

	ticker := time.NewTicker(resetBridgeCycle)
	defer ticker.Stop()

	peerCount := 0
	needResetSubscription := false

	for {
		select {
		case <-sb.quitSync:
			return
		case <-sb.addPeerCh:
			peerCount++
		case <-sb.removePeerCh:
			peerCount--
			if peerCount == 0 {
				needResetSubscription = true
				sb.handler.setParentOperatorNonceSynced(false)
			}
		case <-ticker.C:
			if needResetSubscription && peerCount > 0 {
				err := sb.bridgeManager.ResetAllSubscribedEvents()
				if err == nil {
					needResetSubscription = false
				}
			}
		}
	}
}

func (sb *SubBridge) restoreBridgeLoop() {
	defer sb.pmwg.Done()

	ticker := time.NewTicker(restoreBridgeCycle)
	defer ticker.Stop()

	for {
		select {
		case <-sb.quitSync:
			return
		case <-ticker.C:
			if err := sb.bridgeManager.RestoreBridges(); err != nil {
				logger.Debug("failed to sb.bridgeManager.RestoreBridges()", "err", err)
				continue
			}
			return
		}
	}
}

func (sb *SubBridge) loop() {
	defer sb.pmwg.Done()

	// Keep waiting for and reacting to the various events
	for {
		select {
		case sendData := <-sb.rpcSendCh:
			sb.SendRPCData(sendData)
		// Handle ChainHeadEvent
		case ev := <-sb.chainCh:
			if ev.Block != nil {
				if err := sb.eventhandler.HandleChainHeadEvent(ev.Block); err != nil {
					logger.Error("subbridge block event", "err", err)
				}

				sb.kasAnchor.AnchorPeriodicBlock(ev.Block)
			} else {
				logger.Error("subbridge block event is nil")
			}
		// Handle NewTexsEvent
		//case ev := <-sb.txCh:
		//	if ev.Txs != nil {
		//		if err := sb.eventhandler.HandleTxsEvent(ev.Txs); err != nil {
		//			logger.Error("subbridge tx event", "err", err)
		//		}
		//	} else {
		//		logger.Error("subbridge tx event is nil")
		//	}
		// Handle ChainLogsEvent
		case logs := <-sb.logsCh:
			if err := sb.eventhandler.HandleLogsEvent(logs); err != nil {
				logger.Error("subbridge log event", "err", err)
			}
		// Handle Bridge Event
		case ev := <-sb.chanReqVTev:
			vtRequestEventMeter.Mark(1)
			if err := sb.eventhandler.ProcessRequestEvent(ev); err != nil {
				logger.Error("fail to process request value transfer event ", "err", err)
			}
		case ev := <-sb.chanReqVTencodedEv:
			vtRequestEventMeter.Mark(1)
			if err := sb.eventhandler.ProcessRequestEvent(ev); err != nil {
				logger.Error("fail to process request value transfer event ", "err", err)
			}
		case ev := <-sb.chanHandleVTev:
			vtHandleEventMeter.Mark(1)
			if err := sb.eventhandler.ProcessHandleEvent(ev); err != nil {
				logger.Error("fail to process handle value transfer event ", "err", err)
			}
		case err := <-sb.chainSub.Err():
			if err != nil {
				logger.Error("subbridge block subscription ", "err", err)
			}
			return
		//case err := <-sb.txSub.Err():
		//	if err != nil {
		//		logger.Error("subbridge tx subscription ", "err", err)
		//	}
		//	return
		case err := <-sb.logsSub.Err():
			if err != nil {
				logger.Error("subbridge log subscription ", "err", err)
			}
			return
		case err := <-sb.reqVTevSub.Err():
			if err != nil {
				logger.Error("subbridge token-received subscription ", "err", err)
			}
			return
		case err := <-sb.reqVTencodedEvSub.Err():
			if err != nil {
				logger.Error("subbridge token-received subscription ", "err", err)
			}
			return
		case err := <-sb.handleVTevSub.Err():
			if err != nil {
				logger.Error("subbridge token-transfer subscription ", "err", err)
			}
			return
		}
	}
}

func (sb *SubBridge) removePeer(id string) {
	sb.removePeerCh <- struct{}{}

	// Short circuit if the peer was already removed
	peer := sb.peers.Peer(id)
	if peer == nil {
		return
	}
	logger.Debug("Removing Klaytn peer", "peer", id)

	if err := sb.peers.Unregister(id); err != nil {
		logger.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.GetP2PPeer().Disconnect(p2p.DiscUselessPeer)
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (sb *SubBridge) handleMsg(p BridgePeer) error {
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

	return sb.handler.HandleMainMsg(p, msg)
}

func (sb *SubBridge) syncer() {
	// Start and ensure cleanup of sync mechanisms
	//pm.fetcher.Start()
	//defer pm.fetcher.Stop()
	//defer pm.downloader.Terminate()

	// Wait for different events to fire synchronisation operations
	forceSync := time.NewTicker(forceSyncCycle)
	defer forceSync.Stop()

	for {
		select {
		case peer := <-sb.newPeerCh:
			go sb.synchronise(peer)

		case <-forceSync.C:
			// Force a sync even if not enough peers are present
			go sb.synchronise(sb.peers.BestPeer())

		case <-sb.noMorePeers:
			return
		}
	}
}

func (sb *SubBridge) synchronise(peer BridgePeer) {
	// @TODO Klaytn ServiceChain Sync
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Klaytn protocol.
func (sb *SubBridge) Stop() error {

	close(sb.quitSync)
	sb.bridgeManager.stopAllRecoveries()

	sb.chainSub.Unsubscribe()
	//sb.txSub.Unsubscribe()
	sb.logsSub.Unsubscribe()
	sb.reqVTevSub.Unsubscribe()
	sb.reqVTencodedEvSub.Unsubscribe()
	sb.handleVTevSub.Unsubscribe()
	sb.eventMux.Stop()
	sb.chainDB.Close()

	sb.bridgeManager.Stop()
	sb.bridgeTxPool.Stop()
	sb.bridgeServer.Stop()

	return nil
}
