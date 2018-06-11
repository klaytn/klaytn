package gxp

import (
	"errors"
	"fmt"
	"ground-x/go-gxplatform/accounts"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/common/hexutil"
	"ground-x/go-gxplatform/consensus"
	"ground-x/go-gxplatform/consensus/gxhash"
	"ground-x/go-gxplatform/core"
	"ground-x/go-gxplatform/core/bloombits"
	"ground-x/go-gxplatform/core/rawdb"
	"ground-x/go-gxplatform/core/types"
	"ground-x/go-gxplatform/core/vm"
	"ground-x/go-gxplatform/event"
	"ground-x/go-gxplatform/gxdb"
	"ground-x/go-gxplatform/gxp/downloader"
	"ground-x/go-gxplatform/gxp/filters"
	"ground-x/go-gxplatform/gxp/gasprice"
	"ground-x/go-gxplatform/internal/gxapi"
	"ground-x/go-gxplatform/log"
	"ground-x/go-gxplatform/miner"
	"ground-x/go-gxplatform/node"
	"ground-x/go-gxplatform/p2p"
	"ground-x/go-gxplatform/params"
	"ground-x/go-gxplatform/rlp"
	"ground-x/go-gxplatform/rpc"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"ground-x/go-gxplatform/consensus/istanbul"
	istanbulBackend "ground-x/go-gxplatform/consensus/istanbul/backend"
	"ground-x/go-gxplatform/crypto"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// GXP implements the GXP full node service.
type GXP struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the GXP

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb gxdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *GxpAPIBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	etherbase common.Address

	networkId     uint64
	netRPCService *gxapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)
}

func (s *GXP) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new Ethereum object (including the
// initialisation of the common Ethereum object)
func New(ctx *node.ServiceContext, config *Config) (*GXP, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run gxp.GXP in light sync mode, use les.LightGXP")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	gxp := &GXP{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, config , chainConfig, chainDb),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		gasPrice:       config.GasPrice,
		etherbase:      config.Gxbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	// istanbul BFT. force to set the istanbul etherbase to node key address
	if chainConfig.Istanbul != nil {
		gxp.etherbase = crypto.PubkeyToAddress(ctx.NodeKey().PublicKey)
	}

	log.Info("Initialising GXP protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run geth upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	gxp.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, gxp.chainConfig, gxp.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		gxp.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	gxp.bloomIndexer.Start(gxp.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	gxp.txPool = core.NewTxPool(config.TxPool, gxp.chainConfig, gxp.blockchain)

	if gxp.protocolManager, err = NewProtocolManager(gxp.chainConfig, config.SyncMode, config.NetworkId, gxp.eventMux, gxp.txPool, gxp.engine, gxp.blockchain, chainDb); err != nil {
		return nil, err
	}
	gxp.miner = miner.New(gxp, gxp.chainConfig, gxp.EventMux(), gxp.engine)
	// istanbul BFT
	gxp.miner.SetExtra(makeExtraData(config.ExtraData, gxp.chainConfig.IsBFT))

	gxp.APIBackend = &GxpAPIBackend{gxp, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	gxp.APIBackend.gpo = gasprice.NewOracle(gxp.APIBackend, gpoParams)

	return gxp, nil
}

// istanbul BFT
func makeExtraData(extra []byte, isBFT bool) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"gxp",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.GetMaximumExtraDataSize(isBFT) {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.GetMaximumExtraDataSize(isBFT))
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (gxdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*gxdb.LDBDatabase); ok {
		db.Meter("gxp/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an GXPlatform service
func CreateConsensusEngine(ctx *node.ServiceContext, config *Config, chainConfig *params.ChainConfig, db gxdb.Database) consensus.Engine {
	// If proof-of-authority is requested, set it up
	//if chainConfig.Clique != nil {
	//	return clique.New(chainConfig.Clique, db)
	//}
	if chainConfig.Istanbul != nil {
		if chainConfig.Istanbul.Epoch != 0 {
			config.Istanbul.Epoch = chainConfig.Istanbul.Epoch
		}
		config.Istanbul.ProposerPolicy = istanbul.ProposerPolicy(chainConfig.Istanbul.ProposerPolicy)
		return istanbulBackend.New(&config.Istanbul, ctx.NodeKey(), db)
	}
	// Otherwise assume proof-of-work
	switch {
	case config.Gxhash.PowMode == gxhash.ModeFake:
		log.Warn("Gxhash used in fake mode")
		return gxhash.NewFaker()
	case config.Gxhash.PowMode == gxhash.ModeTest:
		log.Warn("Gxhash used in test mode")
		return gxhash.NewTester()
	case config.Gxhash.PowMode == gxhash.ModeShared:
		log.Warn("Gxhash used in shared mode")
		return gxhash.NewShared()
	default:
		engine := gxhash.New(gxhash.Config{
			CacheDir:       ctx.ResolvePath(config.Gxhash.CacheDir),
			CachesInMem:    config.Gxhash.CachesInMem,
			CachesOnDisk:   config.Gxhash.CachesOnDisk,
			DatasetDir:     config.Gxhash.DatasetDir,
			DatasetsInMem:  config.Gxhash.DatasetsInMem,
			DatasetsOnDisk: config.Gxhash.DatasetsOnDisk,
		})
		engine.SetThreads(-1) // Disable CPU mining
		return engine
	}
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *GXP) APIs() []rpc.API {
	apis := gxapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "gxp",
			Version:   "1.0",
			Service:   NewPublicGXPAPI(s),
			Public:    true,
		}, {
			Namespace: "gxp",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "gxp",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "gxp",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *GXP) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *GXP) Etherbase() (eb common.Address, err error) {
	s.lock.RLock()
	etherbase := s.etherbase
	s.lock.RUnlock()

	if etherbase != (common.Address{}) {
		return etherbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			etherbase := accounts[0].Address

			s.lock.Lock()
			s.etherbase = etherbase
			s.lock.Unlock()

			log.Info("Etherbase automatically configured", "address", etherbase)
			return etherbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("etherbase must be explicitly specified")
}

// SetEtherbase sets the mining reward address.
func (s *GXP) SetEtherbase(etherbase common.Address) {
	s.lock.Lock()
	// istanbul BFT
	if _, ok := s.engine.(consensus.Istanbul); ok {
		log.Error("Cannot set etherbase in Istanbul consensus")
		return
	}
	s.etherbase = etherbase
	s.lock.Unlock()

	s.miner.SetEtherbase(etherbase)
}

func (s *GXP) StartMining(local bool) error {
	eb, err := s.Etherbase()
	if err != nil {
		log.Error("Cannot start mining without etherbase", "err", err)
		return fmt.Errorf("etherbase missing: %v", err)
	}
	//if clique, ok := s.engine.(*clique.Clique); ok {
	//	wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
	//	if wallet == nil || err != nil {
	//		log.Error("Etherbase account unavailable locally", "err", err)
	//		return fmt.Errorf("signer missing: %v", err)
	//	}
	//	clique.Authorize(eb, wallet.SignHash)
	//}
	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start(eb)
	return nil
}

func (s *GXP) StopMining()         { s.miner.Stop() }
func (s *GXP) IsMining() bool      { return s.miner.Mining() }
func (s *GXP) Miner() *miner.Miner { return s.miner }

func (s *GXP) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *GXP) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *GXP) TxPool() *core.TxPool               { return s.txPool }
func (s *GXP) EventMux() *event.TypeMux           { return s.eventMux }
func (s *GXP) Engine() consensus.Engine           { return s.engine }
func (s *GXP) ChainDb() gxdb.Database             { return s.chainDb }
func (s *GXP) IsListening() bool                  { return true } // Always listening
func (s *GXP) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *GXP) NetVersion() uint64                 { return s.networkId }
func (s *GXP) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *GXP) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// GXP protocol implementation.
func (s *GXP) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = gxapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// GXP protocol.
func (s *GXP) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
