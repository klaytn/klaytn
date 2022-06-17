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
	"io"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/klaytn/klaytn/accounts"
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
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

const (
	// chainEventChanSize is the size of channel listening to ChainHeadEvent.
	chainEventChanSize  = 10000
	chainLogChanSize    = 10000
	transactionChanSize = 10000
	rpcBufferSize       = 1024
)

// MainBridgeInfo represents a short summary of the Klaytn sub-protocol metadata
// known about the host peer.
type MainBridgeInfo struct {
	Network    uint64              `json:"network"`    // Klaytn network ID
	BlockScore *big.Int            `json:"blockscore"` // Total blockscore of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // SHA3 hash of the host's best owned block
}

// MainBridge implements the main bridge of service chain.
type MainBridge struct {
	config *SCConfig

	// DB interfaces
	chainDB database.DBManager // Blockchain database

	eventMux       *event.TypeMux
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *api.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price)

	bridgeServer p2p.Server
	ctx          *node.ServiceContext
	maxPeers     int

	APIBackend *MainBridgeAPI

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan BridgePeer
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg   sync.WaitGroup
	pmwg sync.WaitGroup

	blockchain *blockchain.BlockChain
	txPool     *blockchain.TxPool

	chainHeadCh  chan blockchain.ChainHeadEvent
	chainHeadSub event.Subscription
	logsCh       chan []*types.Log
	logsSub      event.Subscription
	txCh         chan blockchain.NewTxsEvent
	txSub        event.Subscription

	peers        *bridgePeerSet
	handler      *MainBridgeHandler
	eventhandler *MainChainEventHandler

	rpcServer     *rpc.Server
	rpcConn       net.Conn
	rpcResponseCh chan []byte
}

// NewMainBridge creates a new MainBridge object (including the
// initialisation of the common MainBridge object)
func NewMainBridge(ctx *node.ServiceContext, config *SCConfig) (*MainBridge, error) {
	chainDB := CreateDB(ctx, config, "scchaindata")

	mb := &MainBridge{
		config:         config,
		chainDB:        chainDB,
		peers:          newBridgePeerSet(),
		newPeerCh:      make(chan BridgePeer),
		noMorePeers:    make(chan struct{}),
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		networkId:      config.NetworkId,
		ctx:            ctx,
		chainHeadCh:    make(chan blockchain.ChainHeadEvent, chainEventChanSize),
		logsCh:         make(chan []*types.Log, chainLogChanSize),
		txCh:           make(chan blockchain.NewTxsEvent, transactionChanSize),
		quitSync:       make(chan struct{}),
		maxPeers:       config.MaxPeer,
		rpcResponseCh:  make(chan []byte),
	}

	logger.Info("Initialising Klaytn-Bridge protocol", "network", config.NetworkId)
	mb.APIBackend = &MainBridgeAPI{mb}

	var err error
	mb.handler, err = NewMainBridgeHandler(mb.config, mb)
	if err != nil {
		return nil, err
	}
	mb.eventhandler, err = NewMainChainEventHandler(mb, mb.handler)
	if err != nil {
		return nil, err
	}

	mb.rpcServer = rpc.NewServer()
	p1, p2 := net.Pipe()
	mb.rpcConn = p1
	go mb.rpcServer.ServeCodec(rpc.NewJSONCodec(p2), rpc.OptionMethodInvocation|rpc.OptionSubscriptions)

	go func() {
		for {
			data := make([]byte, rpcBufferSize)
			rlen, err := mb.rpcConn.Read(data)
			if err != nil {
				if err == io.EOF {
					logger.Trace("EOF from the rpc server pipe")
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					// If no one closes the pipe, this situation should not happen.
					logger.Error("failed to read from the rpc pipe", "err", err, "rlen", rlen)
					return
				}
			}
			logger.Trace("mainbridge message from rpc server pipe", "rlen", rlen)
			err = mb.SendRPCResponseData(data[:rlen])
			if err != nil {
				logger.Error("failed to send response data from RPC server pipe", err)
			}
		}
	}()

	return mb, nil
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *SCConfig, name string) database.DBManager {
	// OpenFilesLimit and LevelDBCacheSize are used by minimum value.
	dbc := &database.DBConfig{Dir: name, DBType: database.LevelDB}
	return ctx.OpenDatabase(dbc)
}

// BridgePeerSet implements PeerSetManager
func (mb *MainBridge) BridgePeerSet() *bridgePeerSet {
	return mb.peers
}

// APIs returns the collection of RPC services the Klaytn sc package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (mb *MainBridge) APIs() []rpc.API {
	// Append all the local APIs and return
	return []rpc.API{
		{
			Namespace: "mainbridge",
			Version:   "1.0",
			Service:   mb.APIBackend,
			Public:    true,
		},
		{
			Namespace: "mainbridge",
			Version:   "1.0",
			Service:   mb.netRPCService,
			Public:    true,
		},
	}
}

func (mb *MainBridge) AccountManager() *accounts.Manager { return mb.accountManager }
func (mb *MainBridge) EventMux() *event.TypeMux          { return mb.eventMux }
func (mb *MainBridge) ChainDB() database.DBManager       { return mb.chainDB }
func (mb *MainBridge) IsListening() bool                 { return true } // Always listening
func (mb *MainBridge) ProtocolVersion() int              { return int(mb.SCProtocol().Versions[0]) }
func (mb *MainBridge) NetVersion() uint64                { return mb.networkId }

func (mb *MainBridge) Components() []interface{} {
	return nil
}

func (mb *MainBridge) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			mb.blockchain = v
			// event from core-service
			mb.chainHeadSub = mb.blockchain.SubscribeChainHeadEvent(mb.chainHeadCh)
			mb.logsSub = mb.blockchain.SubscribeLogsEvent(mb.logsCh)
		case *blockchain.TxPool:
			mb.txPool = v
			// event from core-service
			mb.txSub = mb.txPool.SubscribeNewTxsEvent(mb.txCh)
		case []rpc.API:
			logger.Debug("p2p rpc registered", "len(v)", len(v))
			for _, api := range v {
				if api.Public && api.Namespace == "klay" {
					logger.Error("p2p rpc registered", "namespace", api.Namespace)
					if err := mb.rpcServer.RegisterName(api.Namespace, api.Service); err != nil {
						logger.Error("pRPC failed to register", "namespace", api.Namespace)
					}
				}
			}
		}
	}

	mb.pmwg.Add(1)
	go mb.loop()
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (mb *MainBridge) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (mb *MainBridge) SCProtocol() SCProtocol {
	return SCProtocol{
		Name:     SCProtocolName,
		Versions: SCProtocolVersion,
		Lengths:  SCProtocolLength,
	}
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (mb *MainBridge) NodeInfo() *MainBridgeInfo {
	currentBlock := mb.blockchain.CurrentBlock()
	return &MainBridgeInfo{
		Network:    mb.networkId,
		BlockScore: mb.blockchain.GetTd(currentBlock.Hash(), currentBlock.NumberU64()),
		Genesis:    mb.blockchain.Genesis().Hash(),
		Config:     mb.blockchain.Config(),
		Head:       currentBlock.Hash(),
	}
}

// getChainID returns the current chain id.
func (mb *MainBridge) getChainID() *big.Int {
	return mb.blockchain.Config().ChainID
}

// Start implements node.Service, starting all internal goroutines needed by the
// Klaytn protocol implementation.
func (mb *MainBridge) Start(srvr p2p.Server) error {

	serverConfig := p2p.Config{}
	serverConfig.PrivateKey = mb.ctx.NodeKey()
	serverConfig.Name = mb.ctx.NodeType().String()
	serverConfig.Logger = logger
	serverConfig.ListenAddr = mb.config.MainBridgePort
	serverConfig.MaxPhysicalConnections = mb.maxPeers
	serverConfig.NoDiscovery = true
	serverConfig.EnableMultiChannelServer = false
	serverConfig.NoDial = true

	scprotocols := make([]p2p.Protocol, 0, len(mb.SCProtocol().Versions))
	for i, version := range mb.SCProtocol().Versions {
		// Compatible; initialise the sub-protocol
		version := version
		scprotocols = append(scprotocols, p2p.Protocol{
			Name:    mb.SCProtocol().Name,
			Version: version,
			Length:  mb.SCProtocol().Lengths[i],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := mb.newPeer(int(version), p, rw)
				pubKey, _ := p.ID().Pubkey()
				addr := crypto.PubkeyToAddress(*pubKey)
				peer.SetAddr(addr)
				select {
				case mb.newPeerCh <- peer:
					mb.wg.Add(1)
					defer mb.wg.Done()
					return mb.handle(peer)
				case <-mb.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return mb.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := mb.peers.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}
	mb.bridgeServer = p2p.NewServer(serverConfig)
	mb.bridgeServer.AddProtocols(scprotocols)

	if err := mb.bridgeServer.Start(); err != nil {
		return errors.New("fail to bridgeserver start")
	}

	// Start the RPC service
	mb.netRPCService = api.NewPublicNetAPI(mb.bridgeServer, mb.NetVersion())

	// Figure out a max peers count based on the server limits
	//s.maxPeers = s.bridgeServer.MaxPhysicalConnections()

	go mb.syncer()

	return nil
}

func (mb *MainBridge) newPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter) BridgePeer {
	return newBridgePeer(pv, p, newMeteredMsgWriter(rw))
}

func (mb *MainBridge) handle(p BridgePeer) error {
	// Ignore maxPeers if this is a trusted peer
	if mb.peers.Len() >= mb.maxPeers && !p.GetP2PPeer().Info().Networks[p2p.ConnDefault].Trusted {
		return p2p.DiscTooManyPeers
	}
	p.GetP2PPeer().Log().Debug("Klaytn peer connected", "name", p.GetP2PPeer().Name())

	// Execute the handshake
	var (
		head   = mb.blockchain.CurrentHeader()
		hash   = head.Hash()
		number = head.Number.Uint64()
		td     = mb.blockchain.GetTd(hash, number)
	)

	err := p.Handshake(mb.networkId, mb.getChainID(), td, hash)
	if err != nil {
		p.GetP2PPeer().Log().Debug("Klaytn peer handshake failed", "err", err)
		return err
	}

	// Register the peer locally
	if err := mb.peers.Register(p); err != nil {
		// if starting node with unlock account, can't register peer until finish unlock
		p.GetP2PPeer().Log().Info("Klaytn peer registration failed", "err", err)
		return err
	}
	defer mb.removePeer(p.GetID())

	p.GetP2PPeer().Log().Info("Added a P2P Peer", "peerID", p.GetP2PPeerID())

	// main loop. handle incoming messages.
	for {
		if err := mb.handleMsg(p); err != nil {
			p.GetP2PPeer().Log().Debug("Klaytn message handling failed", "err", err)
			return err
		}
	}
}

func (mb *MainBridge) SendRPCResponseData(data []byte) error {
	peers := mb.BridgePeerSet().peers
	logger.Trace("mainbridge send rpc response data to peers", "data len", len(data), "peers", len(peers))
	for _, peer := range peers {
		err := peer.SendResponseRPC(data)
		if err != nil {
			logger.Error("failed to send rpc response to the peer", "err", err)
		}
		return err
	}

	return nil
}

func (mb *MainBridge) loop() {
	defer mb.pmwg.Done()

	report := time.NewTicker(1 * time.Second)
	defer report.Stop()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle ChainHeadEvent
		case ev := <-mb.chainHeadCh:
			if ev.Block != nil {
				mb.eventhandler.HandleChainHeadEvent(ev.Block)
			} else {
				logger.Error("mainbridge block event is nil")
			}
		// Handle NewTxsEvent
		case ev := <-mb.txCh:
			if ev.Txs != nil {
				mb.eventhandler.HandleTxsEvent(ev.Txs)
			} else {
				logger.Error("mainbridge tx event is nil")
			}
		// Handle ChainLogsEvent
		case logs := <-mb.logsCh:
			mb.eventhandler.HandleLogsEvent(logs)
		case <-report.C:
			// report status
		case err := <-mb.chainHeadSub.Err():
			if err != nil {
				logger.Error("mainbridge block subscription ", "err", err)
			}
			return
		case err := <-mb.txSub.Err():
			if err != nil {
				logger.Error("mainbridge tx subscription ", "err", err)
			}
			return
		case err := <-mb.logsSub.Err():
			if err != nil {
				logger.Error("mainbridge log subscription ", "err", err)
			}
			return
		}
	}
}

func (mb *MainBridge) removePeer(id string) {
	// Short circuit if the peer was already removed
	peer := mb.peers.Peer(id)
	if peer == nil {
		return
	}
	logger.Debug("Removing Klaytn peer", "peer", id)

	if err := mb.peers.Unregister(id); err != nil {
		logger.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.GetP2PPeer().Disconnect(p2p.DiscUselessPeer)
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func (mb *MainBridge) handleMsg(p BridgePeer) error {
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

	return mb.handler.HandleSubMsg(p, msg)
}

func (mb *MainBridge) syncer() {
	// Start and ensure cleanup of sync mechanisms
	//pm.fetcher.Start()
	//defer pm.fetcher.Stop()
	//defer pm.downloader.Terminate()

	// Wait for different events to fire synchronisation operations
	forceSync := time.NewTicker(forceSyncCycle)
	defer forceSync.Stop()

	for {
		select {
		case peer := <-mb.newPeerCh:
			go mb.synchronise(peer)

		case <-forceSync.C:
			// Force a sync even if not enough peers are present
			go mb.synchronise(mb.peers.BestPeer())

		case <-mb.noMorePeers:
			return
		}
	}
}

func (mb *MainBridge) synchronise(peer BridgePeer) {
	// @TODO Klaytn ServiceChain Sync
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Klaytn protocol.
func (mb *MainBridge) Stop() error {

	close(mb.quitSync)

	mb.chainHeadSub.Unsubscribe()
	mb.txSub.Unsubscribe()
	mb.logsSub.Unsubscribe()
	mb.eventMux.Stop()
	mb.chainDB.Close()

	mb.bridgeServer.Stop()

	return nil
}

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}
