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

package cn

import (
	"fmt"
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/bloombits"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/bitutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/clique"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/node/cn/gasprice"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/work"
	"math/big"
	"sync"
)

// ServiceChain implements the Klaytn servicechain node service.
type ServiceChain struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the CN

	// Handlers
	txPool          work.TxPool
	blockchain      work.BlockChain
	protocolManager BackendProtocolManager

	// DB interfaces
	chainDB database.DBManager // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *blockchain.ChainIndexer       // Bloom indexer operating during block imports

	APIBackend *ServiceChainAPIBackend

	miner    Miner
	gasPrice *big.Int
	signer   common.Address

	networkId     uint64
	netRPCService *api.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and coinbase)

	components []interface{}
}

// New creates a new ServiceChain object (including the
// initialisation of the common ServiceChain object)
func NewServiceChain(ctx *node.ServiceContext, config *Config) (*ServiceChain, error) {
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDB := CreateDB(ctx, config, "chaindata")

	chainConfig, genesisHash, genesisErr := blockchain.SetupGenesisBlock(chainDB, config.Genesis, config.NetworkId, false)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	// NOTE-Klaytn Now we use ChainConfig.UnitPrice from genesis.json.
	//         So let's update cn.Config.GasPrice using ChainConfig.UnitPrice.
	config.GasPrice = new(big.Int).SetUint64(chainConfig.UnitPrice)

	logger.Info("Initialised chain configuration", "config", chainConfig)

	cn := &ServiceChain{
		config:         config,
		chainDB:        chainDB,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateCliqueEngine(ctx, config, chainConfig, chainDB),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		gasPrice:       config.GasPrice,
		signer:         config.ServiceChainSigner,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDB, params.BloomBitsBlocks),
	}

	logger.Info("Initialising Klaytn protocol", "versions", cn.engine.Protocol().Versions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		if err := blockchain.CheckBlockChainVersion(chainDB); err != nil {
			return nil, err
		}
	}
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &blockchain.CacheConfig{StateDBCaching: config.StateDBCaching, ArchiveMode: config.NoPruning, CacheSize: config.TrieCacheSize, BlockInterval: config.TrieBlockInterval}
	)
	bc, err := blockchain.NewBlockChain(chainDB, cacheConfig, cn.chainConfig, cn.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	cn.blockchain = bc
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		logger.Error("Rewinding chain to upgrade configuration", "err", compat)
		cn.blockchain.SetHead(compat.RewindTo)
		chainDB.WriteChainConfig(genesisHash, chainConfig)
	}
	cn.bloomIndexer.Start(cn.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	// TODO-Klaytn-ServiceChain: add account creation prevention in the TxPool if TxTypeAccountCreation is supported.
	config.TxPool.NoAccountCreation = config.NoAccountCreation
	cn.txPool = blockchain.NewTxPool(config.TxPool, cn.chainConfig, bc)

	if cn.protocolManager, err = NewProtocolManager(cn.chainConfig, config.SyncMode, config.NetworkId, cn.eventMux, cn.txPool, cn.engine, cn.blockchain, chainDB, ctx.NodeType(), config); err != nil {
		return nil, err
	}

	cn.protocolManager.SetWsEndPoint(config.WsEndpoint)

	// TODO-Klaytn improve to handle drop transaction on network traffic in PN and EN
	cn.miner = work.New(cn, cn.chainConfig, cn.EventMux(), cn.engine, ctx.NodeType(), config.ServiceChainSigner, config.TxResendUseLegacy)

	cn.APIBackend = &ServiceChainAPIBackend{cn, nil}

	gpoParams := config.GPO

	// NOTE-Klaytn Now we use ChainConfig.UnitPrice from genesis.json and updated config.GasPrice with same value.
	//         So let's override gpoParams.Default with config.GasPrice
	gpoParams.Default = config.GasPrice

	cn.APIBackend.gpo = gasprice.NewOracle(cn.APIBackend, gpoParams)

	//TODO-Klaytn add core component
	cn.addComponent(cn.blockchain)
	cn.addComponent(cn.txPool)

	return cn, nil
}

// add component which may be used in another service component
func (s *ServiceChain) addComponent(component interface{}) {
	s.components = append(s.components, component)
}

func (s *ServiceChain) Components() []interface{} {
	return s.components
}

func (s *ServiceChain) SetComponents(component []interface{}) {
	// do nothing
}

// CreateConsensusEngine creates the required type of consensus engine instance for a Klaytn service
func CreateCliqueEngine(ctx *node.ServiceContext, config *Config, chainConfig *params.ChainConfig, db database.DBManager) consensus.Engine {
	// If proof-of-authority is requested, set it up
	if chainConfig.Clique != nil {
		return clique.New(chainConfig.Clique, db)
	} else {
		logger.Crit("Clique Engine needs Genesis file with clique information.")
	}
	return nil
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *ServiceChain) APIs() []rpc.API {
	apis := api.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "klay",
			Version:   "1.0",
			Service:   NewPublicKlayServiceChainAPI(s),
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   NewPublicServiceChainMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.Downloader(), s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateServiceChainMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateServiceChainAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicServiceChainDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateServiceChainDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *ServiceChain) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *ServiceChain) Signer() (eb common.Address, err error) {
	s.lock.RLock()
	signer := s.signer
	s.lock.RUnlock()

	if signer != (common.Address{}) {
		return signer, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			coinbase := accounts[0].Address

			s.lock.Lock()
			s.signer = coinbase
			s.lock.Unlock()

			logger.Info("Signer automatically configured", "address", coinbase)
			return coinbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("Signer must be explicitly specified")
}

func (s *ServiceChain) StartMining(local bool) error {
	eb, err := s.Signer()
	if err != nil {
		logger.Error("Cannot start mining without a signer", "err", err)
		return fmt.Errorf("signer missing: %v", err)
	}
	if clique, ok := s.engine.(*clique.Clique); ok {
		rewardwallet, err := s.accountManager.Find(accounts.Account{Address: eb})
		if rewardwallet == nil || err != nil {
			logger.Error("Signer account unavailable locally", "err", err)
			return fmt.Errorf("signer missing: %v", err)
		}
		clique.Authorize(eb, rewardwallet.SignHash)
	}
	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		s.protocolManager.SetAcceptTxs()
	}
	go s.miner.Start()
	return nil
}

func (s *ServiceChain) StopMining()    { s.miner.Stop() }
func (s *ServiceChain) IsMining() bool { return s.miner.Mining() }
func (s *ServiceChain) Miner() Miner   { return s.miner }

func (s *ServiceChain) AccountManager() *accounts.Manager { return s.accountManager }
func (s *ServiceChain) BlockChain() work.BlockChain       { return s.blockchain }
func (s *ServiceChain) TxPool() work.TxPool               { return s.txPool }
func (s *ServiceChain) EventMux() *event.TypeMux          { return s.eventMux }
func (s *ServiceChain) Engine() consensus.Engine          { return s.engine }
func (s *ServiceChain) ChainDB() database.DBManager       { return s.chainDB }
func (s *ServiceChain) IsListening() bool                 { return true } // Always listening
func (s *ServiceChain) ProtocolVersion() int              { return s.protocolManager.ProtocolVersion() }
func (s *ServiceChain) NetVersion() uint64                { return s.networkId }
func (s *ServiceChain) Progress() klaytn.SyncProgress {
	return s.protocolManager.Downloader().Progress()
}

func (s *ServiceChain) ReBroadcastTxs(transactions types.Transactions) {
	s.protocolManager.ReBroadcastTxs(transactions)
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *ServiceChain) Protocols() []p2p.Protocol {
	return s.protocolManager.GetSubProtocols()
}

// Start implements node.Service, starting all internal goroutines needed by the
// Klaytn protocol implementation.
func (s *ServiceChain) Start(srvr p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = api.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers()
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Klaytn protocol.
func (s *ServiceChain) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDB.Close()
	close(s.shutdownChan)

	return nil
}

// startBloomHandlers starts a batch of goroutines to accept bloom bit database
// retrievals from possibly a range of filters and serving the data to satisfy.
func (cn *ServiceChain) startBloomHandlers() {
	for i := 0; i < bloomServiceThreads; i++ {
		go func() {
			for {
				select {
				case <-cn.shutdownChan:
					return

				case request := <-cn.bloomRequests:
					task := <-request
					task.Bitsets = make([][]byte, len(task.Sections))
					for i, section := range task.Sections {
						head := cn.chainDB.ReadCanonicalHash((section+1)*params.BloomBitsBlocks - 1)
						if compVector, err := cn.chainDB.ReadBloomBits(database.BloomBitsKey(task.Bit, section, head)); err == nil {
							if blob, err := bitutil.DecompressBytes(compVector, int(params.BloomBitsBlocks)/8); err == nil {
								task.Bitsets[i] = blob
							} else {
								task.Error = err
							}
						} else {
							task.Error = err
						}
					}
					request <- task
				}
			}
		}()
	}
}
