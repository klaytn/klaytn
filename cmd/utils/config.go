// Modifications Copyright 2022 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/utils/flags.go (2022/10/19).
// Modified and improved for the klaytn development.

package utils

import (
	"bufio"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Shopify/sarama"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/fdlimit"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas"
	"github.com/klaytn/klaytn/datasync/dbsyncer"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/nat"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/node/sc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/naoina/toml"
	"github.com/urfave/cli/v2"
)

const (
	ClientIdentifier = "klay" // Client identifier to advertise over the network
	SCNNetworkType   = "scn"  // Service Chain Network
	MNNetworkType    = "mn"   // Mainnet Network
	gitCommit        = ""
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var TomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type KlayConfig struct {
	CN               cn.Config
	Node             node.Config
	DB               dbsyncer.DBConfig
	ChainDataFetcher chaindatafetcher.ChainDataFetcherConfig
	ServiceChain     sc.SCConfig
}

func LoadConfig(file string, cfg *KlayConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = TomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func DefaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = ClientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit)
	cfg.HTTPModules = append(cfg.HTTPModules, "klay", "shh", "eth")
	cfg.WSModules = append(cfg.WSModules, "klay", "shh", "eth")
	cfg.IPCPath = "klay.ipc"
	return cfg
}

func MakeConfigNode(ctx *cli.Context) (*node.Node, KlayConfig) {
	// Load defaults.
	cfg := KlayConfig{
		CN:               *cn.GetDefaultConfig(),
		Node:             DefaultNodeConfig(),
		DB:               *dbsyncer.DefaultDBConfig(),
		ChainDataFetcher: *chaindatafetcher.DefaultChainDataFetcherConfig(),
		ServiceChain:     *sc.DefaultServiceChainConfig(),
	}

	// NOTE-Klaytn : klaytn loads the flags from yaml, not toml
	// Load config file.
	// if file := ctx.String(ConfigFileFlag.Name); file != "" {
	// 	if err := LoadConfig(file, &cfg); err != nil {
	// 		log.Fatalf("%v", err)
	// 	}
	// }

	// Apply flags.
	cfg.SetNodeConfig(ctx)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		log.Fatalf("Failed to create the protocol stack: %v", err)
	}
	cfg.SetKlayConfig(ctx, stack)

	cfg.SetDBSyncerConfig(ctx)
	cfg.SetChainDataFetcherConfig(ctx)
	cfg.SetServiceChainConfig(ctx)

	// SetShhConfig(ctx, stack, &cfg.Shh)
	// SetDashboardConfig(ctx, &cfg.Dashboard)

	return stack, cfg
}

func SetP2PConfig(ctx *cli.Context, cfg *p2p.Config) {
	setNodeKey(ctx, cfg)
	setNAT(ctx, cfg)
	setListenAddress(ctx, cfg)

	var nodeType string
	if ctx.IsSet(NodeTypeFlag.Name) {
		nodeType = ctx.String(NodeTypeFlag.Name)
	} else {
		nodeType = NodeTypeFlag.Value
	}

	cfg.ConnectionType = convertNodeType(nodeType)
	if cfg.ConnectionType == common.UNKNOWNNODE {
		logger.Crit("Unknown node type", "nodetype", nodeType)
	}
	logger.Info("Setting connection type", "nodetype", nodeType, "conntype", cfg.ConnectionType)

	// set bootnodes via this function by check specified parameters
	setBootstrapNodes(ctx, cfg)

	if ctx.IsSet(MaxConnectionsFlag.Name) {
		cfg.MaxPhysicalConnections = ctx.Int(MaxConnectionsFlag.Name)
	}
	logger.Info("Setting MaxPhysicalConnections", "MaxPhysicalConnections", cfg.MaxPhysicalConnections)

	if ctx.IsSet(MaxPendingPeersFlag.Name) {
		cfg.MaxPendingPeers = ctx.Int(MaxPendingPeersFlag.Name)
	}

	cfg.NoDiscovery = ctx.Bool(NoDiscoverFlag.Name)

	cfg.RWTimerConfig = p2p.RWTimerConfig{}
	cfg.RWTimerConfig.Interval = ctx.Uint64(RWTimerIntervalFlag.Name)
	cfg.RWTimerConfig.WaitTime = ctx.Duration(RWTimerWaitTimeFlag.Name)

	if netrestrict := ctx.String(NetrestrictFlag.Name); netrestrict != "" {
		list, err := netutil.ParseNetlist(netrestrict)
		if err != nil {
			log.Fatalf("Option %q: %v", NetrestrictFlag.Name, err)
		}
		cfg.NetRestrict = list
	}

	common.MaxRequestContentLength = ctx.Int(MaxRequestContentLengthFlag.Name)

	cfg.NetworkID, _ = getNetworkId(ctx)
}

// setNodeKey creates a node key from set command line flags, either loading it
// from a file or as a specified hex value. If neither flags were provided, this
// method returns nil and an emphemeral key is to be generated.
func setNodeKey(ctx *cli.Context, cfg *p2p.Config) {
	var (
		hex  = ctx.String(NodeKeyHexFlag.Name)
		file = ctx.String(NodeKeyFileFlag.Name)
		key  *ecdsa.PrivateKey
		err  error
	)
	switch {
	case file != "" && hex != "":
		log.Fatalf("Options %q and %q are mutually exclusive", NodeKeyFileFlag.Name, NodeKeyHexFlag.Name)
	case file != "":
		if key, err = crypto.LoadECDSA(file); err != nil {
			log.Fatalf("Option %q: %v", NodeKeyFileFlag.Name, err)
		}
		cfg.PrivateKey = key
	case hex != "":
		if key, err = crypto.HexToECDSA(hex); err != nil {
			log.Fatalf("Option %q: %v", NodeKeyHexFlag.Name, err)
		}
		cfg.PrivateKey = key
	}
}

// setNAT creates a port mapper from command line flags.
func setNAT(ctx *cli.Context, cfg *p2p.Config) {
	if ctx.IsSet(NATFlag.Name) {
		natif, err := nat.Parse(ctx.String(NATFlag.Name))
		if err != nil {
			log.Fatalf("Option %s: %v", NATFlag.Name, err)
		}
		cfg.NAT = natif
	}
}

// setListenAddress creates a TCP listening address string from set command
// line flags.
func setListenAddress(ctx *cli.Context, cfg *p2p.Config) {
	if ctx.IsSet(ListenPortFlag.Name) {
		cfg.ListenAddr = fmt.Sprintf(":%d", ctx.Int(ListenPortFlag.Name))
	}

	if ctx.Bool(MultiChannelUseFlag.Name) {
		cfg.EnableMultiChannelServer = true
		SubListenAddr := fmt.Sprintf(":%d", ctx.Int(SubListenPortFlag.Name))
		cfg.SubListenAddr = []string{SubListenAddr}
	}
}

func convertNodeType(nodetype string) common.ConnType {
	switch strings.ToLower(nodetype) {
	case "cn", "scn":
		return common.CONSENSUSNODE
	case "pn", "spn":
		return common.PROXYNODE
	case "en", "sen":
		return common.ENDPOINTNODE
	default:
		return common.UNKNOWNNODE
	}
}

// setBootstrapNodes creates a list of bootstrap nodes from the command line
// flags, reverting to pre-configured ones if none have been specified.
func setBootstrapNodes(ctx *cli.Context, cfg *p2p.Config) {
	var urls []string
	switch {
	case ctx.IsSet(BootnodesFlag.Name):
		logger.Info("Customized bootnodes are set")
		urls = strings.Split(ctx.String(BootnodesFlag.Name), ",")
	case ctx.Bool(CypressFlag.Name):
		logger.Info("Cypress bootnodes are set")
		urls = params.MainnetBootnodes[cfg.ConnectionType].Addrs
	case ctx.Bool(BaobabFlag.Name):
		logger.Info("Baobab bootnodes are set")
		// set pre-configured bootnodes when 'baobab' option was enabled
		urls = params.BaobabBootnodes[cfg.ConnectionType].Addrs
	case cfg.BootstrapNodes != nil:
		return // already set, don't apply defaults.
	case !ctx.IsSet(NetworkIdFlag.Name):
		if NodeTypeFlag.Value != "scn" && NodeTypeFlag.Value != "spn" && NodeTypeFlag.Value != "sen" {
			logger.Info("Cypress bootnodes are set")
			urls = params.MainnetBootnodes[cfg.ConnectionType].Addrs
		}
	}

	cfg.BootstrapNodes = make([]*discover.Node, 0, len(urls))
	for _, url := range urls {
		node, err := discover.ParseNode(url)
		if err != nil {
			logger.Error("Bootstrap URL invalid", "kni", url, "err", err)
			continue
		}
		if node.NType == discover.NodeTypeUnknown {
			logger.Debug("setBootstrapNode: set nodetype as bn from unknown", "nodeid", node.ID)
			node.NType = discover.NodeTypeBN
		}
		logger.Info("Bootnode - Add Seed", "Node", node)
		cfg.BootstrapNodes = append(cfg.BootstrapNodes, node)
	}
}

// setNodeConfig applies node-related command line flags to the config.
func (kCfg *KlayConfig) SetNodeConfig(ctx *cli.Context) {
	cfg := &kCfg.Node
	// ntp check enable with remote server
	if ctx.Bool(NtpDisableFlag.Name) {
		cfg.NtpRemoteServer = ""
	} else {
		cfg.NtpRemoteServer = ctx.String(NtpServerFlag.Name)
	}

	// disable unsafe debug APIs
	cfg.DisableUnsafeDebug = ctx.Bool(UnsafeDebugDisableFlag.Name)

	SetP2PConfig(ctx, &cfg.P2P)
	setIPC(ctx, cfg)

	// httptype is http
	// fasthttp type is deprecated
	if ctx.IsSet(SrvTypeFlag.Name) {
		cfg.HTTPServerType = ctx.String(SrvTypeFlag.Name)

		if cfg.HTTPServerType == "fasthttp" {
			logger.Warn("The fasthttp option is deprecated. Instead, the server will start with the http type")
			cfg.HTTPServerType = "http"
		}
	}

	setHTTP(ctx, cfg)
	setWS(ctx, cfg)
	setgRPC(ctx, cfg)
	setAPIConfig(ctx)
	setNodeUserIdent(ctx, cfg)

	if dbtype := database.DBType(ctx.String(DbTypeFlag.Name)).ToValid(); len(dbtype) != 0 {
		cfg.DBType = dbtype
	} else {
		logger.Crit("invalid dbtype", "dbtype", ctx.String(DbTypeFlag.Name))
	}
	cfg.DataDir = ctx.String(DataDirFlag.Name)
	cfg.ChainDataDir = ctx.String(ChainDataDirFlag.Name)

	if ctx.IsSet(KeyStoreDirFlag.Name) {
		cfg.KeyStoreDir = ctx.String(KeyStoreDirFlag.Name)
	}
	if ctx.IsSet(LightKDFFlag.Name) {
		cfg.UseLightweightKDF = ctx.Bool(LightKDFFlag.Name)
	}
	if ctx.IsSet(RPCNonEthCompatibleFlag.Name) {
		rpc.NonEthCompatible = ctx.Bool(RPCNonEthCompatibleFlag.Name)
	}
}

// setHTTP creates the HTTP RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setHTTP(ctx *cli.Context, cfg *node.Config) {
	if ctx.Bool(RPCEnabledFlag.Name) && cfg.HTTPHost == "" {
		cfg.HTTPHost = "127.0.0.1"
		if ctx.IsSet(RPCListenAddrFlag.Name) {
			cfg.HTTPHost = ctx.String(RPCListenAddrFlag.Name)
		}
	}

	if ctx.IsSet(RPCPortFlag.Name) {
		cfg.HTTPPort = ctx.Int(RPCPortFlag.Name)
	}
	if ctx.IsSet(RPCCORSDomainFlag.Name) {
		cfg.HTTPCors = SplitAndTrim(ctx.String(RPCCORSDomainFlag.Name))
	}
	if ctx.IsSet(RPCApiFlag.Name) {
		cfg.HTTPModules = SplitAndTrim(ctx.String(RPCApiFlag.Name))
	}
	if ctx.IsSet(RPCVirtualHostsFlag.Name) {
		cfg.HTTPVirtualHosts = SplitAndTrim(ctx.String(RPCVirtualHostsFlag.Name))
	}
	if ctx.IsSet(RPCConcurrencyLimit.Name) {
		rpc.ConcurrencyLimit = ctx.Int(RPCConcurrencyLimit.Name)
		logger.Info("Set the concurrency limit of RPC-HTTP server", "limit", rpc.ConcurrencyLimit)
	}
	if ctx.IsSet(RPCReadTimeout.Name) {
		cfg.HTTPTimeouts.ReadTimeout = time.Duration(ctx.Int(RPCReadTimeout.Name)) * time.Second
	}
	if ctx.IsSet(RPCWriteTimeoutFlag.Name) {
		cfg.HTTPTimeouts.WriteTimeout = time.Duration(ctx.Int(RPCWriteTimeoutFlag.Name)) * time.Second
	}
	if ctx.IsSet(RPCIdleTimeoutFlag.Name) {
		cfg.HTTPTimeouts.IdleTimeout = time.Duration(ctx.Int(RPCIdleTimeoutFlag.Name)) * time.Second
	}
	if ctx.IsSet(RPCExecutionTimeoutFlag.Name) {
		cfg.HTTPTimeouts.ExecutionTimeout = time.Duration(ctx.Int(RPCExecutionTimeoutFlag.Name)) * time.Second
	}
}

// setWS creates the WebSocket RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setWS(ctx *cli.Context, cfg *node.Config) {
	if ctx.Bool(WSEnabledFlag.Name) && cfg.WSHost == "" {
		cfg.WSHost = "127.0.0.1"
		if ctx.IsSet(WSListenAddrFlag.Name) {
			cfg.WSHost = ctx.String(WSListenAddrFlag.Name)
		}
	}

	if ctx.IsSet(WSPortFlag.Name) {
		cfg.WSPort = ctx.Int(WSPortFlag.Name)
	}
	if ctx.IsSet(WSAllowedOriginsFlag.Name) {
		cfg.WSOrigins = SplitAndTrim(ctx.String(WSAllowedOriginsFlag.Name))
	}
	if ctx.IsSet(WSApiFlag.Name) {
		cfg.WSModules = SplitAndTrim(ctx.String(WSApiFlag.Name))
	}
	rpc.MaxSubscriptionPerWSConn = int32(ctx.Int(WSMaxSubscriptionPerConn.Name))
	rpc.WebsocketReadDeadline = ctx.Int64(WSReadDeadLine.Name)
	rpc.WebsocketWriteDeadline = ctx.Int64(WSWriteDeadLine.Name)
	rpc.MaxWebsocketConnections = int32(ctx.Int(WSMaxConnections.Name))
}

// setIPC creates an IPC path configuration from the set command line flags,
// returning an empty string if IPC was explicitly disabled, or the set path.
func setIPC(ctx *cli.Context, cfg *node.Config) {
	CheckExclusive(ctx, IPCDisabledFlag, IPCPathFlag)
	switch {
	case ctx.Bool(IPCDisabledFlag.Name):
		cfg.IPCPath = ""
	case ctx.IsSet(IPCPathFlag.Name):
		cfg.IPCPath = ctx.String(IPCPathFlag.Name)
	}
}

// setgRPC creates the gRPC listener interface string from the set
// command line flags, returning empty if the gRPC endpoint is disabled.
func setgRPC(ctx *cli.Context, cfg *node.Config) {
	if ctx.Bool(GRPCEnabledFlag.Name) && cfg.GRPCHost == "" {
		cfg.GRPCHost = "127.0.0.1"
		if ctx.IsSet(GRPCListenAddrFlag.Name) {
			cfg.GRPCHost = ctx.String(GRPCListenAddrFlag.Name)
		}
	}

	if ctx.IsSet(GRPCPortFlag.Name) {
		cfg.GRPCPort = ctx.Int(GRPCPortFlag.Name)
	}
}

// setAPIConfig sets configurations for specific APIs.
func setAPIConfig(ctx *cli.Context) {
	filters.GetLogsDeadline = ctx.Duration(APIFilterGetLogsDeadlineFlag.Name)
	filters.GetLogsMaxItems = ctx.Int(APIFilterGetLogsMaxItemsFlag.Name)
}

// setNodeUserIdent creates the user identifier from CLI flags.
func setNodeUserIdent(ctx *cli.Context, cfg *node.Config) {
	if identity := ctx.String(IdentityFlag.Name); len(identity) > 0 {
		cfg.UserIdent = identity
	}
}

// setKlayConfig applies klay-related command line flags to the config.
func (kCfg *KlayConfig) SetKlayConfig(ctx *cli.Context, stack *node.Node) {
	// TODO-Klaytn-Bootnode: better have to check conflicts about network flags when we add Klaytn's `mainnet` parameter
	// checkExclusive(ctx, DeveloperFlag, TestnetFlag, RinkebyFlag)
	cfg := &kCfg.CN
	raiseFDLimit()

	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	setServiceChainSigner(ctx, ks, cfg)
	setRewardbase(ctx, ks, cfg)
	setTxPool(ctx, &cfg.TxPool)

	if ctx.IsSet(SyncModeFlag.Name) {
		cfg.SyncMode = *GlobalTextMarshaler(ctx, SyncModeFlag.Name).(*downloader.SyncMode)
		if cfg.SyncMode != downloader.FullSync && cfg.SyncMode != downloader.SnapSync {
			log.Fatalf("Full Sync or Snap Sync (prototype) is supported only!")
		}
		if cfg.SyncMode == downloader.SnapSync {
			logger.Info("Snap sync requested, enabling --snapshot")
			ctx.Set(SnapshotFlag.Name, "true")
		} else {
			cfg.SnapshotCacheSize = 0 // Disabled
		}
	}

	if ctx.Bool(KESNodeTypeServiceFlag.Name) {
		cfg.FetcherDisable = true
		cfg.DownloaderDisable = true
		cfg.WorkerDisable = true
	}

	if NetworkTypeFlag.Value == SCNNetworkType && kCfg.ServiceChain.EnabledSubBridge {
		cfg.NoAccountCreation = !ctx.Bool(ServiceChainNewAccountFlag.Name)
	}

	cfg.NetworkId, cfg.IsPrivate = getNetworkId(ctx)

	if dbtype := database.DBType(ctx.String(DbTypeFlag.Name)).ToValid(); len(dbtype) != 0 {
		cfg.DBType = dbtype
	} else {
		logger.Crit("invalid dbtype", "dbtype", ctx.String(DbTypeFlag.Name))
	}
	cfg.SingleDB = ctx.Bool(SingleDBFlag.Name)
	cfg.NumStateTrieShards = ctx.Uint(NumStateTrieShardsFlag.Name)
	if !database.IsPow2(cfg.NumStateTrieShards) {
		log.Fatalf("%v should be power of 2 but %v is not!", NumStateTrieShardsFlag.Name, cfg.NumStateTrieShards)
	}

	cfg.OverwriteGenesis = ctx.Bool(OverwriteGenesisFlag.Name)
	cfg.StartBlockNumber = ctx.Uint64(StartBlockNumberFlag.Name)

	cfg.LevelDBCompression = database.LevelDBCompressionType(ctx.Int(LevelDBCompressionTypeFlag.Name))
	cfg.LevelDBBufferPool = !ctx.Bool(LevelDBNoBufferPoolFlag.Name)
	cfg.EnableDBPerfMetrics = !ctx.Bool(DBNoPerformanceMetricsFlag.Name)
	cfg.LevelDBCacheSize = ctx.Int(LevelDBCacheSizeFlag.Name)

	cfg.RocksDBConfig.Secondary = ctx.Bool(RocksDBSecondaryFlag.Name)
	if cfg.RocksDBConfig.Secondary {
		cfg.FetcherDisable = true
		cfg.DownloaderDisable = true
		cfg.WorkerDisable = true
	}
	cfg.RocksDBConfig.CacheSize = ctx.Uint64(RocksDBCacheSizeFlag.Name)
	cfg.RocksDBConfig.DumpMallocStat = ctx.Bool(RocksDBDumpMemoryStatFlag.Name)
	cfg.RocksDBConfig.CompressionType = ctx.String(RocksDBCompressionTypeFlag.Name)
	cfg.RocksDBConfig.BottommostCompressionType = ctx.String(RocksDBBottommostCompressionTypeFlag.Name)
	cfg.RocksDBConfig.FilterPolicy = ctx.String(RocksDBFilterPolicyFlag.Name)
	cfg.RocksDBConfig.DisableMetrics = ctx.Bool(RocksDBDisableMetricsFlag.Name)

	cfg.DynamoDBConfig.TableName = ctx.String(DynamoDBTableNameFlag.Name)
	cfg.DynamoDBConfig.Region = ctx.String(DynamoDBRegionFlag.Name)
	cfg.DynamoDBConfig.IsProvisioned = ctx.Bool(DynamoDBIsProvisionedFlag.Name)
	cfg.DynamoDBConfig.ReadCapacityUnits = ctx.Int64(DynamoDBReadCapacityFlag.Name)
	cfg.DynamoDBConfig.WriteCapacityUnits = ctx.Int64(DynamoDBWriteCapacityFlag.Name)
	cfg.DynamoDBConfig.ReadOnly = ctx.Bool(DynamoDBReadOnlyFlag.Name)

	if gcmode := ctx.String(GCModeFlag.Name); gcmode != "full" && gcmode != "archive" {
		log.Fatalf("--%s must be either 'full' or 'archive'", GCModeFlag.Name)
	}
	cfg.NoPruning = ctx.String(GCModeFlag.Name) == "archive"
	logger.Info("Archiving mode of this node", "isArchiveMode", cfg.NoPruning)

	cfg.AnchoringPeriod = ctx.Uint64(AnchoringPeriodFlag.Name)
	cfg.SentChainTxsLimit = ctx.Uint64(SentChainTxsLimit.Name)

	cfg.TrieCacheSize = ctx.Int(TrieMemoryCacheSizeFlag.Name)
	common.DefaultCacheType = common.CacheType(ctx.Int(CacheTypeFlag.Name))
	cfg.TrieBlockInterval = ctx.Uint(TrieBlockIntervalFlag.Name)
	cfg.TriesInMemory = ctx.Uint64(TriesInMemoryFlag.Name)
	cfg.LivePruning = ctx.Bool(LivePruningFlag.Name)
	cfg.LivePruningRetention = ctx.Uint64(LivePruningRetentionFlag.Name)

	if ctx.IsSet(CacheScaleFlag.Name) {
		common.CacheScale = ctx.Int(CacheScaleFlag.Name)
	}
	if ctx.IsSet(CacheUsageLevelFlag.Name) {
		cacheUsageLevelFlag := ctx.String(CacheUsageLevelFlag.Name)
		if scaleByCacheUsageLevel, err := common.GetScaleByCacheUsageLevel(cacheUsageLevelFlag); err != nil {
			logger.Crit("Incorrect CacheUsageLevelFlag value", "error", err, "CacheUsageLevelFlag", cacheUsageLevelFlag)
		} else {
			common.ScaleByCacheUsageLevel = scaleByCacheUsageLevel
		}
	}
	if ctx.IsSet(MemorySizeFlag.Name) {
		physicalMemory := common.TotalPhysicalMemGB
		common.TotalPhysicalMemGB = ctx.Int(MemorySizeFlag.Name)
		logger.Info("Physical memory has been replaced by user settings", "PhysicalMemory(GB)", physicalMemory, "UserSetting(GB)", common.TotalPhysicalMemGB)
	} else {
		logger.Debug("Memory settings", "PhysicalMemory(GB)", common.TotalPhysicalMemGB)
	}

	if ctx.IsSet(DocRootFlag.Name) {
		cfg.DocRoot = ctx.String(DocRootFlag.Name)
	}
	if ctx.IsSet(ExtraDataFlag.Name) {
		cfg.ExtraData = []byte(ctx.String(ExtraDataFlag.Name))
	}

	cfg.SenderTxHashIndexing = ctx.Bool(SenderTxHashIndexingFlag.Name)
	cfg.ParallelDBWrite = !ctx.Bool(NoParallelDBWriteFlag.Name)
	cfg.TrieNodeCacheConfig = statedb.TrieNodeCacheConfig{
		CacheType: statedb.TrieNodeCacheType(ctx.String(TrieNodeCacheTypeFlag.
			Name)).ToValid(),
		NumFetcherPrefetchWorker:  ctx.Int(NumFetcherPrefetchWorkerFlag.Name),
		UseSnapshotForPrefetch:    ctx.Bool(UseSnapshotForPrefetchFlag.Name),
		LocalCacheSizeMiB:         ctx.Int(TrieNodeCacheLimitFlag.Name),
		FastCacheFileDir:          ctx.String(DataDirFlag.Name) + "/fastcache",
		FastCacheSavePeriod:       ctx.Duration(TrieNodeCacheSavePeriodFlag.Name),
		RedisEndpoints:            ctx.StringSlice(TrieNodeCacheRedisEndpointsFlag.Name),
		RedisClusterEnable:        ctx.Bool(TrieNodeCacheRedisClusterFlag.Name),
		RedisPublishBlockEnable:   ctx.Bool(TrieNodeCacheRedisPublishBlockFlag.Name),
		RedisSubscribeBlockEnable: ctx.Bool(TrieNodeCacheRedisSubscribeBlockFlag.Name),
	}

	if ctx.IsSet(VMEnableDebugFlag.Name) {
		// TODO(fjl): force-enable this in --dev mode
		cfg.EnablePreimageRecording = ctx.Bool(VMEnableDebugFlag.Name)
	}
	if ctx.IsSet(VMLogTargetFlag.Name) {
		if _, err := debug.Handler.SetVMLogTarget(ctx.Int(VMLogTargetFlag.Name)); err != nil {
			logger.Warn("Incorrect vmlog value", "err", err)
		}
	}
	cfg.EnableInternalTxTracing = ctx.Bool(VMTraceInternalTxFlag.Name)

	cfg.AutoRestartFlag = ctx.Bool(AutoRestartFlag.Name)
	cfg.RestartTimeOutFlag = ctx.Duration(RestartTimeOutFlag.Name)
	cfg.DaemonPathFlag = ctx.String(DaemonPathFlag.Name)

	if ctx.IsSet(RPCGlobalGasCap.Name) {
		cfg.RPCGasCap = new(big.Int).SetUint64(ctx.Uint64(RPCGlobalGasCap.Name))
	}
	if ctx.IsSet(RPCGlobalEVMTimeoutFlag.Name) {
		cfg.RPCEVMTimeout = ctx.Duration(RPCGlobalEVMTimeoutFlag.Name)
	}
	if ctx.IsSet(RPCGlobalEthTxFeeCapFlag.Name) {
		cfg.RPCTxFeeCap = ctx.Float64(RPCGlobalEthTxFeeCapFlag.Name)
	}

	// Only CNs could set BlockGenerationIntervalFlag and BlockGenerationTimeLimitFlag
	if ctx.IsSet(BlockGenerationIntervalFlag.Name) {
		params.BlockGenerationInterval = ctx.Int64(BlockGenerationIntervalFlag.Name)
		if params.BlockGenerationInterval < 1 {
			logger.Crit("Block generation interval should be equal or larger than 1", "interval", params.BlockGenerationInterval)
		}
	}
	if ctx.IsSet(BlockGenerationTimeLimitFlag.Name) {
		params.BlockGenerationTimeLimit = ctx.Duration(BlockGenerationTimeLimitFlag.Name)
	}

	params.OpcodeComputationCostLimit = ctx.Uint64(OpcodeComputationCostLimitFlag.Name)

	if ctx.IsSet(SnapshotFlag.Name) {
		cfg.SnapshotCacheSize = ctx.Int(SnapshotCacheSizeFlag.Name)
		if cfg.StartBlockNumber != 0 {
			logger.Crit("State snapshot should not be used with --start-block-num", "num", cfg.StartBlockNumber)
		}
		logger.Info("State snapshot is enabled", "cache-size (MB)", cfg.SnapshotCacheSize)
		cfg.SnapshotAsyncGen = ctx.Bool(SnapshotAsyncGen.Name)
	} else {
		cfg.SnapshotCacheSize = 0 // snapshot disabled
	}

	// disable unsafe debug APIs
	cfg.DisableUnsafeDebug = ctx.Bool(UnsafeDebugDisableFlag.Name)

	// Override any default configs for hard coded network.
	// TODO-Klaytn-Bootnode: Discuss and add `baobab` test network's genesis block
	/*
		if ctx.Bool(TestnetFlag.Name) {
			if !ctx.IsSet(NetworkIdFlag.Name) {
				cfg.NetworkId = 3
			}
			cfg.Genesis = blockchain.DefaultBaobabGenesisBlock()
		}
	*/
	// Set the Tx resending related configuration variables
	setTxResendConfig(ctx, cfg)
}

// raiseFDLimit increases the file descriptor limit to process's maximum value
func raiseFDLimit() {
	limit, err := fdlimit.Maximum()
	if err != nil {
		logger.Error("Failed to read maximum fd. you may suffer fd exhaustion", "err", err)
		return
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		logger.Warn("Failed to increase fd limit. you may suffer fd exhaustion", "err", err)
		return
	}
	logger.Info("Raised fd limit to process's maximum value", "fd", raised)
}

// setServiceChainSigner retrieves the service chain signer either from the directly specified
// command line flags or from the keystore if CLI indexed.
func setServiceChainSigner(ctx *cli.Context, ks *keystore.KeyStore, cfg *cn.Config) {
	if ctx.IsSet(ServiceChainSignerFlag.Name) {
		account, err := MakeAddress(ks, ctx.String(ServiceChainSignerFlag.Name))
		if err != nil {
			log.Fatalf("Option %q: %v", ServiceChainSignerFlag.Name, err)
		}
		cfg.ServiceChainSigner = account.Address
	}
}

// setRewardbase retrieves the rewardbase either from the directly specified
// command line flags or from the keystore if CLI indexed.
func setRewardbase(ctx *cli.Context, ks *keystore.KeyStore, cfg *cn.Config) {
	if ctx.IsSet(RewardbaseFlag.Name) {
		account, err := MakeAddress(ks, ctx.String(RewardbaseFlag.Name))
		if err != nil {
			log.Fatalf("Option %q: %v", RewardbaseFlag.Name, err)
		}
		cfg.Rewardbase = account.Address
	}
}

// makeAddress converts an account specified directly as a hex encoded string or
// a key index in the key store to an internal account representation.
func MakeAddress(ks *keystore.KeyStore, account string) (accounts.Account, error) {
	// If the specified account is a valid address, return it
	if common.IsHexAddress(account) {
		return accounts.Account{Address: common.HexToAddress(account)}, nil
	}
	// Otherwise try to interpret the account as a keystore index
	index, err := strconv.Atoi(account)
	if err != nil || index < 0 {
		return accounts.Account{}, fmt.Errorf("invalid account address or index %q", account)
	}
	logger.Warn("Use explicit addresses! Referring to accounts by order in the keystore folder is dangerous and will be deprecated!")

	accs := ks.Accounts()
	if len(accs) <= index {
		return accounts.Account{}, fmt.Errorf("index %d higher than number of accounts %d", index, len(accs))
	}
	return accs[index], nil
}

func setTxPool(ctx *cli.Context, cfg *blockchain.TxPoolConfig) {
	if ctx.IsSet(TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.Bool(TxPoolNoLocalsFlag.Name)
	}
	if ctx.IsSet(TxPoolAllowLocalAnchorTxFlag.Name) {
		cfg.AllowLocalAnchorTx = ctx.Bool(TxPoolAllowLocalAnchorTxFlag.Name)
	}
	if ctx.IsSet(TxPoolDenyRemoteTxFlag.Name) {
		cfg.DenyRemoteTx = ctx.Bool(TxPoolDenyRemoteTxFlag.Name)
	}
	if ctx.IsSet(TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.String(TxPoolJournalFlag.Name)
	}
	if ctx.IsSet(TxPoolJournalIntervalFlag.Name) {
		cfg.JournalInterval = ctx.Duration(TxPoolJournalIntervalFlag.Name)
	}
	if ctx.IsSet(TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.Uint64(TxPoolPriceLimitFlag.Name)
	}
	if ctx.IsSet(TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.Uint64(TxPoolPriceBumpFlag.Name)
	}
	if ctx.IsSet(TxPoolExecSlotsAccountFlag.Name) {
		cfg.ExecSlotsAccount = ctx.Uint64(TxPoolExecSlotsAccountFlag.Name)
	}
	if ctx.IsSet(TxPoolExecSlotsAllFlag.Name) {
		cfg.ExecSlotsAll = ctx.Uint64(TxPoolExecSlotsAllFlag.Name)
	}
	if ctx.IsSet(TxPoolNonExecSlotsAccountFlag.Name) {
		cfg.NonExecSlotsAccount = ctx.Uint64(TxPoolNonExecSlotsAccountFlag.Name)
	}
	if ctx.IsSet(TxPoolNonExecSlotsAllFlag.Name) {
		cfg.NonExecSlotsAll = ctx.Uint64(TxPoolNonExecSlotsAllFlag.Name)
	}

	cfg.KeepLocals = ctx.Bool(TxPoolKeepLocalsFlag.Name)

	if ctx.IsSet(TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.Duration(TxPoolLifetimeFlag.Name)
	}

	// PN specific txpool setting
	if NodeTypeFlag.Value == "pn" {
		cfg.EnableSpamThrottlerAtRuntime = !ctx.Bool(TxPoolSpamThrottlerDisableFlag.Name)
	}
}

// getNetworkId returns the associated network ID with whether or not the network is private.
func getNetworkId(ctx *cli.Context) (uint64, bool) {
	if ctx.Bool(BaobabFlag.Name) && ctx.Bool(CypressFlag.Name) {
		log.Fatalf("--baobab and --cypress must not be set together")
	}
	if ctx.Bool(BaobabFlag.Name) && ctx.IsSet(NetworkIdFlag.Name) {
		log.Fatalf("--baobab and --networkid must not be set together")
	}
	if ctx.Bool(CypressFlag.Name) && ctx.IsSet(NetworkIdFlag.Name) {
		log.Fatalf("--cypress and --networkid must not be set together")
	}

	switch {
	case ctx.Bool(CypressFlag.Name):
		logger.Info("Cypress network ID is set", "networkid", params.CypressNetworkId)
		return params.CypressNetworkId, false
	case ctx.Bool(BaobabFlag.Name):
		logger.Info("Baobab network ID is set", "networkid", params.BaobabNetworkId)
		return params.BaobabNetworkId, false
	case ctx.IsSet(NetworkIdFlag.Name):
		networkId := ctx.Uint64(NetworkIdFlag.Name)
		logger.Info("A private network ID is set", "networkid", networkId)
		return networkId, true
	default:
		if NodeTypeFlag.Value == "scn" || NodeTypeFlag.Value == "spn" || NodeTypeFlag.Value == "sen" {
			logger.Info("A Service Chain default network ID is set", "networkid", params.ServiceChainDefaultNetworkId)
			return params.ServiceChainDefaultNetworkId, true
		}
		logger.Info("Cypress network ID is set", "networkid", params.CypressNetworkId)
		return params.CypressNetworkId, false
	}
}

func setTxResendConfig(ctx *cli.Context, cfg *cn.Config) {
	// Set the Tx resending related configuration variables
	cfg.TxResendInterval = ctx.Uint64(TxResendIntervalFlag.Name)
	if cfg.TxResendInterval == 0 {
		cfg.TxResendInterval = cn.DefaultTxResendInterval
	}

	cfg.TxResendCount = ctx.Int(TxResendCountFlag.Name)
	if cfg.TxResendCount < cn.DefaultMaxResendTxCount {
		cfg.TxResendCount = cn.DefaultMaxResendTxCount
	}
	cfg.TxResendUseLegacy = ctx.Bool(TxResendUseLegacyFlag.Name)
	logger.Debug("TxResend config", "Interval", cfg.TxResendInterval, "TxResendCount", cfg.TxResendCount, "UseLegacy", cfg.TxResendUseLegacy)
}

func (kCfg *KlayConfig) SetChainDataFetcherConfig(ctx *cli.Context) {
	cfg := &kCfg.ChainDataFetcher
	if ctx.Bool(EnableChainDataFetcherFlag.Name) {
		cfg.EnabledChainDataFetcher = true

		if ctx.Bool(ChainDataFetcherNoDefault.Name) {
			cfg.NoDefaultStart = true
		}
		if ctx.IsSet(ChainDataFetcherNumHandlers.Name) {
			cfg.NumHandlers = ctx.Int(ChainDataFetcherNumHandlers.Name)
		}
		if ctx.IsSet(ChainDataFetcherJobChannelSize.Name) {
			cfg.JobChannelSize = ctx.Int(ChainDataFetcherJobChannelSize.Name)
		}
		if ctx.IsSet(ChainDataFetcherChainEventSizeFlag.Name) {
			cfg.BlockChannelSize = ctx.Int(ChainDataFetcherChainEventSizeFlag.Name)
		}
		if ctx.IsSet(ChainDataFetcherMaxProcessingDataSize.Name) {
			cfg.MaxProcessingDataSize = ctx.Int(ChainDataFetcherMaxProcessingDataSize.Name)
		}

		mode := ctx.String(ChainDataFetcherMode.Name)
		mode = strings.ToLower(mode)
		switch mode {
		case "kas": // kas option is not used.
			cfg.Mode = chaindatafetcher.ModeKAS
			cfg.KasConfig = makeKASConfig(ctx)
		case "kafka":
			cfg.Mode = chaindatafetcher.ModeKafka
			cfg.KafkaConfig = makeKafkaConfig(ctx)
		default:
			logger.Crit("unsupported chaindatafetcher mode (\"kas\", \"kafka\")", "mode", cfg.Mode)
		}
	}
}

// NOTE-klaytn
// Deprecated: KASConfig is not used anymore.
func checkKASDBConfigs(ctx *cli.Context) {
	if !ctx.IsSet(ChainDataFetcherKASDBHostFlag.Name) {
		logger.Crit("DBHost must be set !", "key", ChainDataFetcherKASDBHostFlag.Name)
	}
	if !ctx.IsSet(ChainDataFetcherKASDBUserFlag.Name) {
		logger.Crit("DBUser must be set !", "key", ChainDataFetcherKASDBUserFlag.Name)
	}
	if !ctx.IsSet(ChainDataFetcherKASDBPasswordFlag.Name) {
		logger.Crit("DBPassword must be set !", "key", ChainDataFetcherKASDBPasswordFlag.Name)
	}
	if !ctx.IsSet(ChainDataFetcherKASDBNameFlag.Name) {
		logger.Crit("DBName must be set !", "key", ChainDataFetcherKASDBNameFlag.Name)
	}
}

// NOTE-klaytn
// Deprecated: KASConfig is not used anymore.
func checkKASCacheInvalidationConfigs(ctx *cli.Context) {
	if !ctx.IsSet(ChainDataFetcherKASCacheURLFlag.Name) {
		logger.Crit("The cache invalidation url is not set")
	}
	if !ctx.IsSet(ChainDataFetcherKASBasicAuthParamFlag.Name) {
		logger.Crit("The authorization is not set")
	}
	if !ctx.IsSet(ChainDataFetcherKASXChainIdFlag.Name) {
		logger.Crit("The x-chain-id is not set")
	}
}

// NOTE-klaytn
// Deprecated: KASConfig is not used anymore.
func makeKASConfig(ctx *cli.Context) *kas.KASConfig {
	kasConfig := kas.DefaultKASConfig

	checkKASDBConfigs(ctx)
	kasConfig.DBHost = ctx.String(ChainDataFetcherKASDBHostFlag.Name)
	kasConfig.DBPort = ctx.String(ChainDataFetcherKASDBPortFlag.Name)
	kasConfig.DBUser = ctx.String(ChainDataFetcherKASDBUserFlag.Name)
	kasConfig.DBPassword = ctx.String(ChainDataFetcherKASDBPasswordFlag.Name)
	kasConfig.DBName = ctx.String(ChainDataFetcherKASDBNameFlag.Name)

	if ctx.Bool(ChainDataFetcherKASCacheUse.Name) {
		checkKASCacheInvalidationConfigs(ctx)
		kasConfig.CacheUse = true
		kasConfig.CacheInvalidationURL = ctx.String(ChainDataFetcherKASCacheURLFlag.Name)
		kasConfig.BasicAuthParam = ctx.String(ChainDataFetcherKASBasicAuthParamFlag.Name)
		kasConfig.XChainId = ctx.String(ChainDataFetcherKASXChainIdFlag.Name)
	}
	return kasConfig
}

func makeKafkaConfig(ctx *cli.Context) *kafka.KafkaConfig {
	kafkaConfig := kafka.GetDefaultKafkaConfig()
	if ctx.IsSet(ChainDataFetcherKafkaBrokersFlag.Name) {
		kafkaConfig.Brokers = ctx.StringSlice(ChainDataFetcherKafkaBrokersFlag.Name)
	} else {
		logger.Crit("The kafka brokers must be set")
	}
	kafkaConfig.TopicEnvironmentName = ctx.String(ChainDataFetcherKafkaTopicEnvironmentFlag.Name)
	kafkaConfig.TopicResourceName = ctx.String(ChainDataFetcherKafkaTopicResourceFlag.Name)
	kafkaConfig.Partitions = int32(ctx.Int64(ChainDataFetcherKafkaPartitionsFlag.Name))
	kafkaConfig.Replicas = int16(ctx.Int64(ChainDataFetcherKafkaReplicasFlag.Name))
	kafkaConfig.SaramaConfig.Producer.MaxMessageBytes = ctx.Int(ChainDataFetcherKafkaMaxMessageBytesFlag.Name)
	kafkaConfig.SegmentSizeBytes = ctx.Int(ChainDataFetcherKafkaSegmentSizeBytesFlag.Name)
	kafkaConfig.MsgVersion = ctx.String(ChainDataFetcherKafkaMessageVersionFlag.Name)
	kafkaConfig.ProducerId = ctx.String(ChainDataFetcherKafkaProducerIdFlag.Name)
	requiredAcks := sarama.RequiredAcks(ctx.Int(ChainDataFetcherKafkaRequiredAcksFlag.Name))
	if requiredAcks != sarama.NoResponse && requiredAcks != sarama.WaitForLocal && requiredAcks != sarama.WaitForAll {
		logger.Crit("not supported requiredAcks. it must be NoResponse(0), WaitForLocal(1), or WaitForAll(-1)", "given", requiredAcks)
	}
	kafkaConfig.SaramaConfig.Producer.RequiredAcks = requiredAcks
	return kafkaConfig
}

func (kCfg *KlayConfig) SetDBSyncerConfig(ctx *cli.Context) {
	cfg := &kCfg.DB
	if ctx.Bool(EnableDBSyncerFlag.Name) {
		cfg.EnabledDBSyncer = true

		if ctx.IsSet(DBHostFlag.Name) {
			dbhost := ctx.String(DBHostFlag.Name)
			cfg.DBHost = dbhost
		} else {
			logger.Crit("DBHost must be set !", "key", DBHostFlag.Name)
		}
		if ctx.IsSet(DBPortFlag.Name) {
			dbports := ctx.String(DBPortFlag.Name)
			cfg.DBPort = dbports
		}
		if ctx.IsSet(DBUserFlag.Name) {
			dbuser := ctx.String(DBUserFlag.Name)
			cfg.DBUser = dbuser
		} else {
			logger.Crit("DBUser must be set !", "key", DBUserFlag.Name)
		}
		if ctx.IsSet(DBPasswordFlag.Name) {
			dbpasswd := ctx.String(DBPasswordFlag.Name)
			cfg.DBPassword = dbpasswd
		} else {
			logger.Crit("DBPassword must be set !", "key", DBPasswordFlag.Name)
		}
		if ctx.IsSet(DBNameFlag.Name) {
			dbname := ctx.String(DBNameFlag.Name)
			cfg.DBName = dbname
		} else {
			logger.Crit("DBName must be set !", "key", DBNameFlag.Name)
		}
		if ctx.Bool(EnabledLogModeFlag.Name) {
			cfg.EnabledLogMode = true
		}
		if ctx.IsSet(MaxIdleConnsFlag.Name) {
			cfg.MaxIdleConns = ctx.Int(MaxIdleConnsFlag.Name)
		}
		if ctx.IsSet(MaxOpenConnsFlag.Name) {
			cfg.MaxOpenConns = ctx.Int(MaxOpenConnsFlag.Name)
		}
		if ctx.IsSet(ConnMaxLifeTimeFlag.Name) {
			cfg.ConnMaxLifetime = ctx.Duration(ConnMaxLifeTimeFlag.Name)
		}
		if ctx.IsSet(DBSyncerModeFlag.Name) {
			cfg.Mode = strings.ToLower(ctx.String(DBSyncerModeFlag.Name))
		}
		if ctx.IsSet(GenQueryThreadFlag.Name) {
			cfg.GenQueryThread = ctx.Int(GenQueryThreadFlag.Name)
		}
		if ctx.IsSet(InsertThreadFlag.Name) {
			cfg.InsertThread = ctx.Int(InsertThreadFlag.Name)
		}
		if ctx.IsSet(BulkInsertSizeFlag.Name) {
			cfg.BulkInsertSize = ctx.Int(BulkInsertSizeFlag.Name)
		}
		if ctx.IsSet(EventModeFlag.Name) {
			cfg.EventMode = strings.ToLower(ctx.String(EventModeFlag.Name))
		}
		if ctx.IsSet(MaxBlockDiffFlag.Name) {
			cfg.MaxBlockDiff = ctx.Uint64(MaxBlockDiffFlag.Name)
		}
		if ctx.IsSet(BlockSyncChannelSizeFlag.Name) {
			cfg.BlockChannelSize = ctx.Int(BlockSyncChannelSizeFlag.Name)
		}
	}
}

func (kCfg *KlayConfig) SetServiceChainConfig(ctx *cli.Context) {
	cfg := &kCfg.ServiceChain

	// bridge service
	if ctx.Bool(MainBridgeFlag.Name) {
		cfg.EnabledMainBridge = true
		cfg.MainBridgePort = fmt.Sprintf(":%d", ctx.Int(MainBridgeListenPortFlag.Name))
	}

	if ctx.Bool(SubBridgeFlag.Name) {
		cfg.EnabledSubBridge = true
		cfg.SubBridgePort = fmt.Sprintf(":%d", ctx.Int(SubBridgeListenPortFlag.Name))
	}

	cfg.Anchoring = ctx.Bool(ServiceChainAnchoringFlag.Name)
	cfg.ChildChainIndexing = ctx.Bool(ChildChainIndexingFlag.Name)
	cfg.AnchoringPeriod = ctx.Uint64(AnchoringPeriodFlag.Name)
	cfg.SentChainTxsLimit = ctx.Uint64(SentChainTxsLimit.Name)
	cfg.ParentChainID = ctx.Uint64(ParentChainIDFlag.Name)
	cfg.VTRecovery = ctx.Bool(VTRecoveryFlag.Name)
	cfg.VTRecoveryInterval = ctx.Uint64(VTRecoveryIntervalFlag.Name)
	cfg.ServiceChainConsensus = ServiceChainConsensusFlag.Value
	cfg.ServiceChainParentOperatorGasLimit = ctx.Uint64(ServiceChainParentOperatorTxGasLimitFlag.Name)
	cfg.ServiceChainChildOperatorGasLimit = ctx.Uint64(ServiceChainChildOperatorTxGasLimitFlag.Name)

	cfg.KASAnchor = ctx.Bool(KASServiceChainAnchorFlag.Name)
	if cfg.KASAnchor {
		cfg.KASAnchorPeriod = ctx.Uint64(KASServiceChainAnchorPeriodFlag.Name)
		if cfg.KASAnchorPeriod == 0 {
			cfg.KASAnchorPeriod = 1
			logger.Warn("KAS anchor period is set by 1")
		}

		cfg.KASAnchorUrl = ctx.String(KASServiceChainAnchorUrlFlag.Name)
		if cfg.KASAnchorUrl == "" {
			logger.Crit("KAS anchor url should be set", "key", KASServiceChainAnchorUrlFlag.Name)
		}

		cfg.KASAnchorOperator = ctx.String(KASServiceChainAnchorOperatorFlag.Name)
		if cfg.KASAnchorOperator == "" {
			logger.Crit("KAS anchor operator should be set", "key", KASServiceChainAnchorOperatorFlag.Name)
		}

		cfg.KASAccessKey = ctx.String(KASServiceChainAccessKeyFlag.Name)
		if cfg.KASAccessKey == "" {
			logger.Crit("KAS access key should be set", "key", KASServiceChainAccessKeyFlag.Name)
		}

		cfg.KASSecretKey = ctx.String(KASServiceChainSecretKeyFlag.Name)
		if cfg.KASSecretKey == "" {
			logger.Crit("KAS secret key should be set", "key", KASServiceChainSecretKeyFlag.Name)
		}

		cfg.KASXChainId = ctx.String(KASServiceChainXChainIdFlag.Name)
		if cfg.KASXChainId == "" {
			logger.Crit("KAS x-chain-id should be set", "key", KASServiceChainXChainIdFlag.Name)
		}

		cfg.KASAnchorRequestTimeout = ctx.Duration(KASServiceChainAnchorRequestTimeoutFlag.Name)
	}

	cfg.DataDir = kCfg.Node.DataDir
	cfg.Name = kCfg.Node.Name
}
