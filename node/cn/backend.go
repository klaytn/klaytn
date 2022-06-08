// Modifications Copyright 2018 The klaytn Authors
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
	"errors"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/bloombits"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	istanbulBackend "github.com/klaytn/klaytn/consensus/istanbul/backend"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/node/cn/gasprice"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/work"
)

var errCNLightSync = errors.New("can't run cn.CN in light sync mode")

//go:generate mockgen -destination=node/cn/mocks/lesserver_mock.go -package=mocks github.com/klaytn/klaytn/node/cn LesServer
type LesServer interface {
	Start(srvr p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *blockchain.ChainIndexer)
}

//go:generate mockgen -destination=node/cn/mocks/miner_mock.go -package=mocks github.com/klaytn/klaytn/node/cn Miner
// Miner is an interface of work.Miner used by ServiceChain.
type Miner interface {
	Start()
	Stop()
	Register(agent work.Agent)
	Mining() bool
	HashRate() (tot int64)
	SetExtra(extra []byte) error
	Pending() (*types.Block, *state.StateDB)
	PendingBlock() *types.Block
}

//go:generate mockgen -destination=node/cn/protocolmanager_mock_test.go github.com/klaytn/klaytn/node/cn BackendProtocolManager
// BackendProtocolManager is an interface of cn.ProtocolManager used from cn.CN and cn.ServiceChain.
type BackendProtocolManager interface {
	Downloader() ProtocolManagerDownloader
	SetWsEndPoint(wsep string)
	GetSubProtocols() []p2p.Protocol
	ProtocolVersion() int
	ReBroadcastTxs(transactions types.Transactions)
	SetAcceptTxs()
	SetRewardbase(addr common.Address)
	SetRewardbaseWallet(wallet accounts.Wallet)
	NodeType() common.ConnType
	Start(maxPeers int)
	Stop()
	SetSyncStop(flag bool)
}

// CN implements the Klaytn consensus node service.
type CN struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Handlers
	txPool          work.TxPool
	blockchain      work.BlockChain
	protocolManager BackendProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDB database.DBManager // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager accounts.AccountManager

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *blockchain.ChainIndexer       // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

	APIBackend *CNAPIBackend

	miner    Miner
	gasPrice *big.Int

	rewardbase common.Address

	networkId     uint64
	netRPCService *api.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price)

	components []interface{}

	governance governance.Engine
}

func (s *CN) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// senderTxHashIndexer subscribes chainEvent and stores senderTxHash to txHash mapping information.
func senderTxHashIndexer(db database.DBManager, chainEvent <-chan blockchain.ChainEvent, subscription event.Subscription) {
	defer subscription.Unsubscribe()

	for {
		select {
		case event := <-chainEvent:
			var err error
			batch := db.NewSenderTxHashToTxHashBatch()
			for _, tx := range event.Block.Transactions() {
				senderTxHash, ok := tx.SenderTxHash()

				// senderTxHash and txHash are the same if tx is not a fee-delegated tx.
				// Do not store mapping between senderTxHash and txHash in this case.
				if !ok {
					continue
				}

				txHash := tx.Hash()

				if err = db.PutSenderTxHashToTxHashToBatch(batch, senderTxHash, txHash); err != nil {
					logger.Error("Failed to store senderTxHash to txHash mapping to database",
						"blockNum", event.Block.Number(), "senderTxHash", senderTxHash, "txHash", txHash, "err", err)
					break
				}
			}

			if err == nil {
				batch.Write()
			}

		case <-subscription.Err():
			return
		}
	}
}

func checkSyncMode(config *Config) error {
	if !config.SyncMode.IsValid() {
		return fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	if config.SyncMode == downloader.LightSync {
		return errCNLightSync
	}
	return nil
}

func setEngineType(chainConfig *params.ChainConfig) {
	if chainConfig.Clique != nil {
		types.EngineType = types.Engine_Clique
	}
	if chainConfig.Istanbul != nil {
		types.EngineType = types.Engine_IBFT
	}
}

// New creates a new CN object (including the
// initialisation of the common CN object)
func New(ctx *node.ServiceContext, config *Config) (*CN, error) {
	if err := checkSyncMode(config); err != nil {
		return nil, err
	}

	chainDB := CreateDB(ctx, config, "chaindata")

	chainConfig, genesisHash, genesisErr := blockchain.SetupGenesisBlock(chainDB, config.Genesis, config.NetworkId, config.IsPrivate, false)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	setEngineType(chainConfig)

	// load governance state
	governance := governance.NewMixedEngine(chainConfig, chainDB)

	// Set latest unitPrice/gasPrice
	chainConfig.UnitPrice = governance.UnitPrice()
	config.GasPrice = new(big.Int).SetUint64(chainConfig.UnitPrice)
	logger.Info("Initialised chain configuration", "config", chainConfig)

	cn := &CN{
		config:            config,
		chainDB:           chainDB,
		chainConfig:       chainConfig,
		eventMux:          ctx.EventMux,
		accountManager:    ctx.AccountManager,
		engine:            CreateConsensusEngine(ctx, config, chainConfig, chainDB, governance, ctx.NodeType()),
		networkId:         config.NetworkId,
		gasPrice:          config.GasPrice,
		rewardbase:        config.Rewardbase,
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		bloomIndexer:      NewBloomIndexer(chainDB, params.BloomBitsBlocks),
		closeBloomHandler: make(chan struct{}),
		governance:        governance,
	}

	// istanbul BFT. Derive and set node's address using nodekey
	if cn.chainConfig.Istanbul != nil {
		governance.SetNodeAddress(crypto.PubkeyToAddress(ctx.NodeKey().PublicKey))
	}

	logger.Info("Initialising Klaytn protocol", "versions", cn.engine.Protocol().Versions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		if err := blockchain.CheckBlockChainVersion(chainDB); err != nil {
			return nil, err
		}
	}
	var (
		vmConfig    = config.getVMConfig()
		cacheConfig = &blockchain.CacheConfig{ArchiveMode: config.NoPruning, CacheSize: config.TrieCacheSize,
			BlockInterval: config.TrieBlockInterval, TriesInMemory: config.TriesInMemory,
			TrieNodeCacheConfig: &config.TrieNodeCacheConfig, SenderTxHashIndexing: config.SenderTxHashIndexing, SnapshotCacheSize: config.SnapshotCacheSize}
	)

	bc, err := blockchain.NewBlockChain(chainDB, cacheConfig, cn.chainConfig, cn.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	bc.SetCanonicalBlock(config.StartBlockNumber)

	cn.blockchain = bc
	governance.SetBlockchain(cn.blockchain)
	// Synchronize proposerpolicy & useGiniCoeff
	if cn.blockchain.Config().Istanbul != nil {
		cn.blockchain.Config().Istanbul.ProposerPolicy = governance.ProposerPolicy()
	}
	if cn.blockchain.Config().Governance.Reward != nil {
		cn.blockchain.Config().Governance.Reward.UseGiniCoeff = governance.UseGiniCoeff()
	}

	if config.SenderTxHashIndexing {
		ch := make(chan blockchain.ChainEvent, 255)
		chainEventSubscription := cn.blockchain.SubscribeChainEvent(ch)
		go senderTxHashIndexer(chainDB, ch, chainEventSubscription)
	}

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		logger.Error("Rewinding chain to upgrade configuration", "err", compat)
		cn.blockchain.SetHead(compat.RewindTo)
		chainDB.WriteChainConfig(genesisHash, cn.chainConfig)
	}
	cn.bloomIndexer.Start(cn.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	// TODO-Klaytn-ServiceChain: add account creation prevention in the txPool if TxTypeAccountCreation is supported.
	config.TxPool.NoAccountCreation = config.NoAccountCreation
	cn.txPool = blockchain.NewTxPool(config.TxPool, cn.chainConfig, bc)
	governance.SetTxPool(cn.txPool)

	// Permit the downloader to use the trie cache allowance during fast sync
	cacheLimit := cacheConfig.TrieNodeCacheConfig.LocalCacheSizeMiB
	if cn.protocolManager, err = NewProtocolManager(cn.chainConfig, config.SyncMode, config.NetworkId, cn.eventMux, cn.txPool, cn.engine, cn.blockchain, chainDB, cacheLimit, ctx.NodeType(), config); err != nil {
		return nil, err
	}

	if err := cn.setAcceptTxs(); err != nil {
		logger.Error("Failed to decode IstanbulExtra", "err", err)
	}

	cn.protocolManager.SetWsEndPoint(config.WsEndpoint)

	if err := cn.setRewardWallet(); err != nil {
		logger.Error("Error happened while setting the reward wallet", "err", err)
	}

	if governance.ProposerPolicy() == uint64(istanbul.WeightedRandom) {
		// NewStakingManager is called with proper non-nil parameters
		reward.NewStakingManager(cn.blockchain, governance, cn.chainDB)
	}

	// set worker
	if config.WorkerDisable {
		cn.miner = work.NewFakeWorker()
		// Istanbul backend can be accessed by APIs to call its methods even though the core of the
		// consensus engine doesn't run.
		istBackend, ok := cn.engine.(consensus.Istanbul)
		if ok {
			istBackend.SetChain(cn.blockchain)
		}
	} else {
		// TODO-Klaytn improve to handle drop transaction on network traffic in PN and EN
		cn.miner = work.New(cn, cn.chainConfig, cn.EventMux(), cn.engine, ctx.NodeType(), crypto.PubkeyToAddress(ctx.NodeKey().PublicKey), cn.config.TxResendUseLegacy)
	}

	// istanbul BFT
	cn.miner.SetExtra(makeExtraData(config.ExtraData))

	cn.APIBackend = &CNAPIBackend{cn, nil}

	gpoParams := config.GPO

	// NOTE-Klaytn Now we use latest unitPrice
	//         So let's override gpoParams.Default with config.GasPrice
	gpoParams.Default = config.GasPrice

	cn.APIBackend.gpo = gasprice.NewOracle(cn.APIBackend, gpoParams, cn.txPool)
	//@TODO Klaytn add core component
	cn.addComponent(cn.blockchain)
	cn.addComponent(cn.txPool)
	cn.addComponent(cn.APIs())
	cn.addComponent(cn.ChainDB())

	if config.AutoRestartFlag {
		daemonPath := config.DaemonPathFlag
		restartInterval := config.RestartTimeOutFlag
		if restartInterval <= time.Second {
			logger.Crit("Invalid auto-restart timeout", "timeout", restartInterval)
		}

		// Restarts the node with the same configuration if blockNumber is not changed for a specific time.
		restartTimer := time.AfterFunc(restartInterval, func() {
			logger.Warn("Restart node", "command", daemonPath+" restart")
			cmd := exec.Command(daemonPath, "restart")
			cmd.Run()
		})
		logger.Info("Initialize auto-restart feature", "timeout", restartInterval, "daemonPath", daemonPath)

		go func() {
			blockChecker := time.NewTicker(time.Second)
			prevBlockNum := cn.blockchain.CurrentBlock().NumberU64()

			for range blockChecker.C {
				currentBlockNum := cn.blockchain.CurrentBlock().NumberU64()

				if prevBlockNum != currentBlockNum {
					prevBlockNum = currentBlockNum
					restartTimer.Reset(restartInterval)
				}
			}
		}()
	}

	// Only for KES nodes
	if config.TrieNodeCacheConfig.RedisSubscribeBlockEnable {
		go cn.blockchain.BlockSubscriptionLoop(cn.txPool.(*blockchain.TxPool))
	}

	return cn, nil
}

// setAcceptTxs sets AcceptTxs flag in 1CN case to receive tx propagation.
func (s *CN) setAcceptTxs() error {
	if s.chainConfig.Istanbul != nil {
		istanbulExtra, err := types.ExtractIstanbulExtra(s.blockchain.Genesis().Header())
		if err != nil {
			return err
		} else {
			if len(istanbulExtra.Validators) == 1 {
				s.protocolManager.SetAcceptTxs()
			}
		}
	}
	return nil
}

// setRewardWallet sets reward base and reward base wallet if the node is CN.
func (s *CN) setRewardWallet() error {
	if s.protocolManager.NodeType() == common.CONSENSUSNODE {
		wallet, err := s.RewardbaseWallet()
		if err != nil {
			return err
		} else {
			s.protocolManager.SetRewardbaseWallet(wallet)
		}
		s.protocolManager.SetRewardbase(s.rewardbase)
	}
	return nil
}

// add component which may be used in another service component
func (s *CN) addComponent(component interface{}) {
	s.components = append(s.components, component)
}

func (s *CN) Components() []interface{} {
	return s.components
}

func (s *CN) SetComponents(component []interface{}) {
	// do nothing
}

// istanbul BFT
func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"klay",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.GetMaximumExtraDataSize() {
		logger.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.GetMaximumExtraDataSize())
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) database.DBManager {
	dbc := &database.DBConfig{Dir: name, DBType: config.DBType, ParallelDBWrite: config.ParallelDBWrite, SingleDB: config.SingleDB, NumStateTrieShards: config.NumStateTrieShards,
		LevelDBCacheSize: config.LevelDBCacheSize, OpenFilesLimit: database.GetOpenFilesLimit(), LevelDBCompression: config.LevelDBCompression,
		LevelDBBufferPool: config.LevelDBBufferPool, EnableDBPerfMetrics: config.EnableDBPerfMetrics, DynamoDBConfig: &config.DynamoDBConfig}
	return ctx.OpenDatabase(dbc)
}

// CreateConsensusEngine creates the required type of consensus engine instance for a Klaytn service
func CreateConsensusEngine(ctx *node.ServiceContext, config *Config, chainConfig *params.ChainConfig, db database.DBManager, gov governance.Engine, nodetype common.ConnType) consensus.Engine {
	// Only istanbul  BFT is allowed in the main net. PoA is supported by service chain
	if chainConfig.Governance == nil {
		chainConfig.Governance = params.GetDefaultGovernanceConfig(params.UseIstanbul)
	}
	return istanbulBackend.New(config.Rewardbase, &config.Istanbul, ctx.NodeKey(), db, gov, nodetype)
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *CN) APIs() []rpc.API {
	apis, ethAPI := api.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	publicFilterAPI := filters.NewPublicFilterAPI(s.APIBackend, false)
	governanceKlayAPI := governance.NewGovernanceKlayAPI(s.governance, s.blockchain)
	publicGovernanceAPI := governance.NewGovernanceAPI(s.governance)
	publicDownloaderAPI := downloader.NewPublicDownloaderAPI(s.protocolManager.Downloader(), s.eventMux)

	ethAPI.SetPublicFilterAPI(publicFilterAPI)
	ethAPI.SetGovernanceKlayAPI(governanceKlayAPI)
	ethAPI.SetPublicGovernanceAPI(publicGovernanceAPI)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "klay",
			Version:   "1.0",
			Service:   NewPublicKlayAPI(s),
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   publicDownloaderAPI,
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   publicFilterAPI,
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   publicDownloaderAPI,
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "governance",
			Version:   "1.0",
			Service:   governance.NewGovernanceAPI(s.governance),
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   governanceKlayAPI,
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   ethAPI,
			Public:    true,
		},
	}...)
}

func (s *CN) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *CN) Rewardbase() (eb common.Address, err error) {
	s.lock.RLock()
	rewardbase := s.rewardbase
	s.lock.RUnlock()

	if rewardbase != (common.Address{}) {
		return rewardbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			rewardbase := accounts[0].Address

			s.lock.Lock()
			s.rewardbase = rewardbase
			s.lock.Unlock()

			logger.Info("Rewardbase automatically configured", "address", rewardbase)
			return rewardbase, nil
		}
	}

	return common.Address{}, fmt.Errorf("rewardbase must be explicitly specified")
}

func (s *CN) RewardbaseWallet() (accounts.Wallet, error) {
	rewardBase, err := s.Rewardbase()
	if err != nil {
		return nil, err
	}

	account := accounts.Account{Address: rewardBase}
	wallet, err := s.AccountManager().Find(account)
	if err != nil {
		logger.Error("find err", "err", err)
		return nil, err
	}
	return wallet, nil
}

func (s *CN) SetRewardbase(rewardbase common.Address) {
	s.lock.Lock()
	s.rewardbase = rewardbase
	s.lock.Unlock()
	wallet, err := s.RewardbaseWallet()
	if err != nil {
		logger.Error("find err", "err", err)
	}
	s.protocolManager.SetRewardbase(rewardbase)
	s.protocolManager.SetRewardbaseWallet(wallet)
}

func (s *CN) StartMining(local bool) error {
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

func (s *CN) StopMining()    { s.miner.Stop() }
func (s *CN) IsMining() bool { return s.miner.Mining() }
func (s *CN) Miner() Miner   { return s.miner }

func (s *CN) AccountManager() accounts.AccountManager { return s.accountManager }
func (s *CN) BlockChain() work.BlockChain             { return s.blockchain }
func (s *CN) TxPool() work.TxPool                     { return s.txPool }
func (s *CN) EventMux() *event.TypeMux                { return s.eventMux }
func (s *CN) Engine() consensus.Engine                { return s.engine }
func (s *CN) ChainDB() database.DBManager             { return s.chainDB }
func (s *CN) IsListening() bool                       { return true } // Always listening
func (s *CN) ProtocolVersion() int                    { return s.protocolManager.ProtocolVersion() }
func (s *CN) NetVersion() uint64                      { return s.networkId }
func (s *CN) Progress() klaytn.SyncProgress           { return s.protocolManager.Downloader().Progress() }

func (s *CN) ReBroadcastTxs(transactions types.Transactions) {
	s.protocolManager.ReBroadcastTxs(transactions)
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *CN) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.GetSubProtocols()
	}
	return append(s.protocolManager.GetSubProtocols(), s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Klaytn protocol implementation.
func (s *CN) Start(srvr p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = api.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers()
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}

	reward.StakingManagerSubscribe()

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Klaytn protocol.
func (s *CN) Stop() error {
	// Stop all the peer-related stuff first.
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}

	// Then stop everything else.
	s.bloomIndexer.Close()
	close(s.closeBloomHandler)
	s.txPool.Stop()
	s.miner.Stop()
	reward.StakingManagerUnsubscribe()
	s.blockchain.Stop()
	s.chainDB.Close()
	s.eventMux.Stop()

	return nil
}
