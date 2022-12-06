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
	"gopkg.in/urfave/cli.v1"
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

	// Load config file.
	if file := ctx.GlobalString(ConfigFileFlag.Name); file != "" {
		if err := LoadConfig(file, &cfg); err != nil {
			log.Fatalf("%v", err)
		}
	}

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
	if ctx.GlobalIsSet(NodeTypeFlag.Name) {
		nodeType = ctx.GlobalString(NodeTypeFlag.Name)
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

	if ctx.GlobalIsSet(MaxConnectionsFlag.Name) {
		cfg.MaxPhysicalConnections = ctx.GlobalInt(MaxConnectionsFlag.Name)
	}
	logger.Info("Setting MaxPhysicalConnections", "MaxPhysicalConnections", cfg.MaxPhysicalConnections)

	if ctx.GlobalIsSet(MaxPendingPeersFlag.Name) {
		cfg.MaxPendingPeers = ctx.GlobalInt(MaxPendingPeersFlag.Name)
	}

	cfg.NoDiscovery = ctx.GlobalIsSet(NoDiscoverFlag.Name)

	cfg.RWTimerConfig = p2p.RWTimerConfig{}
	cfg.RWTimerConfig.Interval = ctx.GlobalUint64(RWTimerIntervalFlag.Name)
	cfg.RWTimerConfig.WaitTime = ctx.GlobalDuration(RWTimerWaitTimeFlag.Name)

	if netrestrict := ctx.GlobalString(NetrestrictFlag.Name); netrestrict != "" {
		list, err := netutil.ParseNetlist(netrestrict)
		if err != nil {
			log.Fatalf("Option %q: %v", NetrestrictFlag.Name, err)
		}
		cfg.NetRestrict = list
	}

	common.MaxRequestContentLength = ctx.GlobalInt(MaxRequestContentLengthFlag.Name)

	cfg.NetworkID, _ = getNetworkId(ctx)
}

// setNodeKey creates a node key from set command line flags, either loading it
// from a file or as a specified hex value. If neither flags were provided, this
// method returns nil and an emphemeral key is to be generated.
func setNodeKey(ctx *cli.Context, cfg *p2p.Config) {
	var (
		hex  = ctx.GlobalString(NodeKeyHexFlag.Name)
		file = ctx.GlobalString(NodeKeyFileFlag.Name)
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
	if ctx.GlobalIsSet(NATFlag.Name) {
		natif, err := nat.Parse(ctx.GlobalString(NATFlag.Name))
		if err != nil {
			log.Fatalf("Option %s: %v", NATFlag.Name, err)
		}
		cfg.NAT = natif
	}
}

// setListenAddress creates a TCP listening address string from set command
// line flags.
func setListenAddress(ctx *cli.Context, cfg *p2p.Config) {
	if ctx.GlobalIsSet(ListenPortFlag.Name) {
		cfg.ListenAddr = fmt.Sprintf(":%d", ctx.GlobalInt(ListenPortFlag.Name))
	}

	if ctx.GlobalBool(MultiChannelUseFlag.Name) {
		cfg.EnableMultiChannelServer = true
		SubListenAddr := fmt.Sprintf(":%d", ctx.GlobalInt(SubListenPortFlag.Name))
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
	case ctx.GlobalIsSet(BootnodesFlag.Name):
		logger.Info("Customized bootnodes are set")
		urls = strings.Split(ctx.GlobalString(BootnodesFlag.Name), ",")
	case ctx.GlobalIsSet(CypressFlag.Name):
		logger.Info("Cypress bootnodes are set")
		urls = params.MainnetBootnodes[cfg.ConnectionType].Addrs
	case ctx.GlobalIsSet(BaobabFlag.Name):
		logger.Info("Baobab bootnodes are set")
		// set pre-configured bootnodes when 'baobab' option was enabled
		urls = params.BaobabBootnodes[cfg.ConnectionType].Addrs
	case cfg.BootstrapNodes != nil:
		return // already set, don't apply defaults.
	case !ctx.GlobalIsSet(NetworkIdFlag.Name):
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
	if ctx.GlobalBool(NtpDisableFlag.Name) {
		cfg.NtpRemoteServer = ""
	} else {
		cfg.NtpRemoteServer = ctx.GlobalString(NtpServerFlag.Name)
	}
	SetP2PConfig(ctx, &cfg.P2P)
	setIPC(ctx, cfg)

	// httptype is http or fasthttp
	if ctx.GlobalIsSet(SrvTypeFlag.Name) {
		cfg.HTTPServerType = ctx.GlobalString(SrvTypeFlag.Name)
	}

	setHTTP(ctx, cfg)
	setWS(ctx, cfg)
	setgRPC(ctx, cfg)
	setAPIConfig(ctx)
	setNodeUserIdent(ctx, cfg)

	if dbtype := database.DBType(ctx.GlobalString(DbTypeFlag.Name)).ToValid(); len(dbtype) != 0 {
		cfg.DBType = dbtype
	} else {
		logger.Crit("invalid dbtype", "dbtype", ctx.GlobalString(DbTypeFlag.Name))
	}
	cfg.DataDir = ctx.GlobalString(DataDirFlag.Name)

	if ctx.GlobalIsSet(KeyStoreDirFlag.Name) {
		cfg.KeyStoreDir = ctx.GlobalString(KeyStoreDirFlag.Name)
	}
	if ctx.GlobalIsSet(LightKDFFlag.Name) {
		cfg.UseLightweightKDF = ctx.GlobalBool(LightKDFFlag.Name)
	}
	if ctx.GlobalIsSet(RPCNonEthCompatibleFlag.Name) {
		rpc.NonEthCompatible = ctx.GlobalBool(RPCNonEthCompatibleFlag.Name)
	}
}

// setHTTP creates the HTTP RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setHTTP(ctx *cli.Context, cfg *node.Config) {
	if ctx.GlobalBool(RPCEnabledFlag.Name) && cfg.HTTPHost == "" {
		cfg.HTTPHost = "127.0.0.1"
		if ctx.GlobalIsSet(RPCListenAddrFlag.Name) {
			cfg.HTTPHost = ctx.GlobalString(RPCListenAddrFlag.Name)
		}
	}

	if ctx.GlobalIsSet(RPCPortFlag.Name) {
		cfg.HTTPPort = ctx.GlobalInt(RPCPortFlag.Name)
	}
	if ctx.GlobalIsSet(RPCCORSDomainFlag.Name) {
		cfg.HTTPCors = SplitAndTrim(ctx.GlobalString(RPCCORSDomainFlag.Name))
	}
	if ctx.GlobalIsSet(RPCApiFlag.Name) {
		cfg.HTTPModules = SplitAndTrim(ctx.GlobalString(RPCApiFlag.Name))
	}
	if ctx.GlobalIsSet(RPCVirtualHostsFlag.Name) {
		cfg.HTTPVirtualHosts = SplitAndTrim(ctx.GlobalString(RPCVirtualHostsFlag.Name))
	}
	if ctx.GlobalIsSet(RPCConcurrencyLimit.Name) {
		rpc.ConcurrencyLimit = ctx.GlobalInt(RPCConcurrencyLimit.Name)
		logger.Info("Set the concurrency limit of RPC-HTTP server", "limit", rpc.ConcurrencyLimit)
	}
	if ctx.GlobalIsSet(RPCReadTimeout.Name) {
		cfg.HTTPTimeouts.ReadTimeout = time.Duration(ctx.GlobalInt(RPCReadTimeout.Name)) * time.Second
	}
	if ctx.GlobalIsSet(RPCWriteTimeoutFlag.Name) {
		cfg.HTTPTimeouts.WriteTimeout = time.Duration(ctx.GlobalInt(RPCWriteTimeoutFlag.Name)) * time.Second
	}
	if ctx.GlobalIsSet(RPCIdleTimeoutFlag.Name) {
		cfg.HTTPTimeouts.IdleTimeout = time.Duration(ctx.GlobalInt(RPCIdleTimeoutFlag.Name)) * time.Second
	}
	if ctx.GlobalIsSet(RPCExecutionTimeoutFlag.Name) {
		cfg.HTTPTimeouts.ExecutionTimeout = time.Duration(ctx.GlobalInt(RPCExecutionTimeoutFlag.Name)) * time.Second
	}
}

// setWS creates the WebSocket RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setWS(ctx *cli.Context, cfg *node.Config) {
	if ctx.GlobalBool(WSEnabledFlag.Name) && cfg.WSHost == "" {
		cfg.WSHost = "127.0.0.1"
		if ctx.GlobalIsSet(WSListenAddrFlag.Name) {
			cfg.WSHost = ctx.GlobalString(WSListenAddrFlag.Name)
		}
	}

	if ctx.GlobalIsSet(WSPortFlag.Name) {
		cfg.WSPort = ctx.GlobalInt(WSPortFlag.Name)
	}
	if ctx.GlobalIsSet(WSAllowedOriginsFlag.Name) {
		cfg.WSOrigins = SplitAndTrim(ctx.GlobalString(WSAllowedOriginsFlag.Name))
	}
	if ctx.GlobalIsSet(WSApiFlag.Name) {
		cfg.WSModules = SplitAndTrim(ctx.GlobalString(WSApiFlag.Name))
	}
	rpc.MaxSubscriptionPerWSConn = int32(ctx.GlobalInt(WSMaxSubscriptionPerConn.Name))
	rpc.WebsocketReadDeadline = ctx.GlobalInt64(WSReadDeadLine.Name)
	rpc.WebsocketWriteDeadline = ctx.GlobalInt64(WSWriteDeadLine.Name)
	rpc.MaxWebsocketConnections = int32(ctx.GlobalInt(WSMaxConnections.Name))
}

// setIPC creates an IPC path configuration from the set command line flags,
// returning an empty string if IPC was explicitly disabled, or the set path.
func setIPC(ctx *cli.Context, cfg *node.Config) {
	CheckExclusive(ctx, IPCDisabledFlag, IPCPathFlag)
	switch {
	case ctx.GlobalBool(IPCDisabledFlag.Name):
		cfg.IPCPath = ""
	case ctx.GlobalIsSet(IPCPathFlag.Name):
		cfg.IPCPath = ctx.GlobalString(IPCPathFlag.Name)
	}
}

// setgRPC creates the gRPC listener interface string from the set
// command line flags, returning empty if the gRPC endpoint is disabled.
func setgRPC(ctx *cli.Context, cfg *node.Config) {
	if ctx.GlobalBool(GRPCEnabledFlag.Name) && cfg.GRPCHost == "" {
		cfg.GRPCHost = "127.0.0.1"
		if ctx.GlobalIsSet(GRPCListenAddrFlag.Name) {
			cfg.GRPCHost = ctx.GlobalString(GRPCListenAddrFlag.Name)
		}
	}

	if ctx.GlobalIsSet(GRPCPortFlag.Name) {
		cfg.GRPCPort = ctx.GlobalInt(GRPCPortFlag.Name)
	}
}

// setAPIConfig sets configurations for specific APIs.
func setAPIConfig(ctx *cli.Context) {
	filters.GetLogsDeadline = ctx.GlobalDuration(APIFilterGetLogsDeadlineFlag.Name)
	filters.GetLogsMaxItems = ctx.GlobalInt(APIFilterGetLogsMaxItemsFlag.Name)
}

// setNodeUserIdent creates the user identifier from CLI flags.
func setNodeUserIdent(ctx *cli.Context, cfg *node.Config) {
	if identity := ctx.GlobalString(IdentityFlag.Name); len(identity) > 0 {
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

	if ctx.GlobalIsSet(SyncModeFlag.Name) {
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

	if ctx.GlobalBool(KESNodeTypeServiceFlag.Name) {
		cfg.FetcherDisable = true
		cfg.DownloaderDisable = true
		cfg.WorkerDisable = true
	}

	if NetworkTypeFlag.Value == SCNNetworkType && kCfg.ServiceChain.EnabledSubBridge {
		cfg.NoAccountCreation = !ctx.GlobalBool(ServiceChainNewAccountFlag.Name)
	}

	cfg.NetworkId, cfg.IsPrivate = getNetworkId(ctx)

	if dbtype := database.DBType(ctx.GlobalString(DbTypeFlag.Name)).ToValid(); len(dbtype) != 0 {
		cfg.DBType = dbtype
	} else {
		logger.Crit("invalid dbtype", "dbtype", ctx.GlobalString(DbTypeFlag.Name))
	}
	cfg.SingleDB = ctx.GlobalIsSet(SingleDBFlag.Name)
	cfg.NumStateTrieShards = ctx.GlobalUint(NumStateTrieShardsFlag.Name)
	if !database.IsPow2(cfg.NumStateTrieShards) {
		log.Fatalf("%v should be power of 2 but %v is not!", NumStateTrieShardsFlag.Name, cfg.NumStateTrieShards)
	}

	cfg.OverwriteGenesis = ctx.GlobalBool(OverwriteGenesisFlag.Name)
	cfg.StartBlockNumber = ctx.GlobalUint64(StartBlockNumberFlag.Name)

	cfg.LevelDBCompression = database.LevelDBCompressionType(ctx.GlobalInt(LevelDBCompressionTypeFlag.Name))
	cfg.LevelDBBufferPool = !ctx.GlobalIsSet(LevelDBNoBufferPoolFlag.Name)
	cfg.EnableDBPerfMetrics = !ctx.GlobalIsSet(DBNoPerformanceMetricsFlag.Name)
	cfg.LevelDBCacheSize = ctx.GlobalInt(LevelDBCacheSizeFlag.Name)

	cfg.DynamoDBConfig.TableName = ctx.GlobalString(DynamoDBTableNameFlag.Name)
	cfg.DynamoDBConfig.Region = ctx.GlobalString(DynamoDBRegionFlag.Name)
	cfg.DynamoDBConfig.IsProvisioned = ctx.GlobalBool(DynamoDBIsProvisionedFlag.Name)
	cfg.DynamoDBConfig.ReadCapacityUnits = ctx.GlobalInt64(DynamoDBReadCapacityFlag.Name)
	cfg.DynamoDBConfig.WriteCapacityUnits = ctx.GlobalInt64(DynamoDBWriteCapacityFlag.Name)
	cfg.DynamoDBConfig.ReadOnly = ctx.GlobalBool(DynamoDBReadOnlyFlag.Name)

	if gcmode := ctx.GlobalString(GCModeFlag.Name); gcmode != "full" && gcmode != "archive" {
		log.Fatalf("--%s must be either 'full' or 'archive'", GCModeFlag.Name)
	}
	cfg.NoPruning = ctx.GlobalString(GCModeFlag.Name) == "archive"
	logger.Info("Archiving mode of this node", "isArchiveMode", cfg.NoPruning)

	cfg.AnchoringPeriod = ctx.GlobalUint64(AnchoringPeriodFlag.Name)
	cfg.SentChainTxsLimit = ctx.GlobalUint64(SentChainTxsLimit.Name)

	cfg.TrieCacheSize = ctx.GlobalInt(TrieMemoryCacheSizeFlag.Name)
	common.DefaultCacheType = common.CacheType(ctx.GlobalInt(CacheTypeFlag.Name))
	cfg.TrieBlockInterval = ctx.GlobalUint(TrieBlockIntervalFlag.Name)
	cfg.TriesInMemory = ctx.GlobalUint64(TriesInMemoryFlag.Name)

	if ctx.GlobalIsSet(CacheScaleFlag.Name) {
		common.CacheScale = ctx.GlobalInt(CacheScaleFlag.Name)
	}
	if ctx.GlobalIsSet(CacheUsageLevelFlag.Name) {
		cacheUsageLevelFlag := ctx.GlobalString(CacheUsageLevelFlag.Name)
		if scaleByCacheUsageLevel, err := common.GetScaleByCacheUsageLevel(cacheUsageLevelFlag); err != nil {
			logger.Crit("Incorrect CacheUsageLevelFlag value", "error", err, "CacheUsageLevelFlag", cacheUsageLevelFlag)
		} else {
			common.ScaleByCacheUsageLevel = scaleByCacheUsageLevel
		}
	}
	if ctx.GlobalIsSet(MemorySizeFlag.Name) {
		physicalMemory := common.TotalPhysicalMemGB
		common.TotalPhysicalMemGB = ctx.GlobalInt(MemorySizeFlag.Name)
		logger.Info("Physical memory has been replaced by user settings", "PhysicalMemory(GB)", physicalMemory, "UserSetting(GB)", common.TotalPhysicalMemGB)
	} else {
		logger.Debug("Memory settings", "PhysicalMemory(GB)", common.TotalPhysicalMemGB)
	}

	if ctx.GlobalIsSet(DocRootFlag.Name) {
		cfg.DocRoot = ctx.GlobalString(DocRootFlag.Name)
	}
	if ctx.GlobalIsSet(ExtraDataFlag.Name) {
		cfg.ExtraData = []byte(ctx.GlobalString(ExtraDataFlag.Name))
	}

	cfg.SenderTxHashIndexing = ctx.GlobalIsSet(SenderTxHashIndexingFlag.Name)
	cfg.ParallelDBWrite = !ctx.GlobalIsSet(NoParallelDBWriteFlag.Name)
	cfg.TrieNodeCacheConfig = statedb.TrieNodeCacheConfig{
		CacheType: statedb.TrieNodeCacheType(ctx.GlobalString(TrieNodeCacheTypeFlag.
			Name)).ToValid(),
		NumFetcherPrefetchWorker:  ctx.GlobalInt(NumFetcherPrefetchWorkerFlag.Name),
		UseSnapshotForPrefetch:    ctx.GlobalBool(UseSnapshotForPrefetchFlag.Name),
		LocalCacheSizeMiB:         ctx.GlobalInt(TrieNodeCacheLimitFlag.Name),
		FastCacheFileDir:          ctx.GlobalString(DataDirFlag.Name) + "/fastcache",
		FastCacheSavePeriod:       ctx.GlobalDuration(TrieNodeCacheSavePeriodFlag.Name),
		RedisEndpoints:            ctx.GlobalStringSlice(TrieNodeCacheRedisEndpointsFlag.Name),
		RedisClusterEnable:        ctx.GlobalBool(TrieNodeCacheRedisClusterFlag.Name),
		RedisPublishBlockEnable:   ctx.GlobalBool(TrieNodeCacheRedisPublishBlockFlag.Name),
		RedisSubscribeBlockEnable: ctx.GlobalBool(TrieNodeCacheRedisSubscribeBlockFlag.Name),
	}

	if ctx.GlobalIsSet(VMEnableDebugFlag.Name) {
		// TODO(fjl): force-enable this in --dev mode
		cfg.EnablePreimageRecording = ctx.GlobalBool(VMEnableDebugFlag.Name)
	}
	if ctx.GlobalIsSet(VMLogTargetFlag.Name) {
		if _, err := debug.Handler.SetVMLogTarget(ctx.GlobalInt(VMLogTargetFlag.Name)); err != nil {
			logger.Warn("Incorrect vmlog value", "err", err)
		}
	}
	cfg.EnableInternalTxTracing = ctx.GlobalIsSet(VMTraceInternalTxFlag.Name)

	cfg.AutoRestartFlag = ctx.GlobalBool(AutoRestartFlag.Name)
	cfg.RestartTimeOutFlag = ctx.GlobalDuration(RestartTimeOutFlag.Name)
	cfg.DaemonPathFlag = ctx.GlobalString(DaemonPathFlag.Name)

	if ctx.GlobalIsSet(RPCGlobalGasCap.Name) {
		cfg.RPCGasCap = new(big.Int).SetUint64(ctx.GlobalUint64(RPCGlobalGasCap.Name))
	}
	if ctx.GlobalIsSet(RPCGlobalEVMTimeoutFlag.Name) {
		cfg.RPCEVMTimeout = ctx.GlobalDuration(RPCGlobalEVMTimeoutFlag.Name)
	}
	if ctx.GlobalIsSet(RPCGlobalEthTxFeeCapFlag.Name) {
		cfg.RPCTxFeeCap = ctx.GlobalFloat64(RPCGlobalEthTxFeeCapFlag.Name)
	}

	// Only CNs could set BlockGenerationIntervalFlag and BlockGenerationTimeLimitFlag
	if ctx.GlobalIsSet(BlockGenerationIntervalFlag.Name) {
		params.BlockGenerationInterval = ctx.GlobalInt64(BlockGenerationIntervalFlag.Name)
		if params.BlockGenerationInterval < 1 {
			logger.Crit("Block generation interval should be equal or larger than 1", "interval", params.BlockGenerationInterval)
		}
	}
	if ctx.GlobalIsSet(BlockGenerationTimeLimitFlag.Name) {
		params.BlockGenerationTimeLimit = ctx.GlobalDuration(BlockGenerationTimeLimitFlag.Name)
	}

	params.OpcodeComputationCostLimit = ctx.GlobalUint64(OpcodeComputationCostLimitFlag.Name)

	if ctx.GlobalIsSet(SnapshotFlag.Name) {
		cfg.SnapshotCacheSize = ctx.GlobalInt(SnapshotCacheSizeFlag.Name)
		if cfg.StartBlockNumber != 0 {
			logger.Crit("State snapshot should not be used with --start-block-num", "num", cfg.StartBlockNumber)
		}
		logger.Info("State snapshot is enabled", "cache-size (MB)", cfg.SnapshotCacheSize)
		cfg.SnapshotAsyncGen = ctx.GlobalBool(SnapshotAsyncGen.Name)
	} else {
		cfg.SnapshotCacheSize = 0 // snapshot disabled
	}

	// Override any default configs for hard coded network.
	// TODO-Klaytn-Bootnode: Discuss and add `baobab` test network's genesis block
	/*
		if ctx.GlobalBool(TestnetFlag.Name) {
			if !ctx.GlobalIsSet(NetworkIdFlag.Name) {
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
	if ctx.GlobalIsSet(ServiceChainSignerFlag.Name) {
		account, err := MakeAddress(ks, ctx.GlobalString(ServiceChainSignerFlag.Name))
		if err != nil {
			log.Fatalf("Option %q: %v", ServiceChainSignerFlag.Name, err)
		}
		cfg.ServiceChainSigner = account.Address
	}
}

// setRewardbase retrieves the rewardbase either from the directly specified
// command line flags or from the keystore if CLI indexed.
func setRewardbase(ctx *cli.Context, ks *keystore.KeyStore, cfg *cn.Config) {
	if ctx.GlobalIsSet(RewardbaseFlag.Name) {
		account, err := MakeAddress(ks, ctx.GlobalString(RewardbaseFlag.Name))
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
	if ctx.GlobalIsSet(TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolAllowLocalAnchorTxFlag.Name) {
		cfg.AllowLocalAnchorTx = ctx.GlobalBool(TxPoolAllowLocalAnchorTxFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolDenyRemoteTxFlag.Name) {
		cfg.DenyRemoteTx = ctx.GlobalBool(TxPoolDenyRemoteTxFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolJournalIntervalFlag.Name) {
		cfg.JournalInterval = ctx.GlobalDuration(TxPoolJournalIntervalFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.GlobalUint64(TxPoolPriceLimitFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolExecSlotsAccountFlag.Name) {
		cfg.ExecSlotsAccount = ctx.GlobalUint64(TxPoolExecSlotsAccountFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolExecSlotsAllFlag.Name) {
		cfg.ExecSlotsAll = ctx.GlobalUint64(TxPoolExecSlotsAllFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolNonExecSlotsAccountFlag.Name) {
		cfg.NonExecSlotsAccount = ctx.GlobalUint64(TxPoolNonExecSlotsAccountFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolNonExecSlotsAllFlag.Name) {
		cfg.NonExecSlotsAll = ctx.GlobalUint64(TxPoolNonExecSlotsAllFlag.Name)
	}

	cfg.KeepLocals = ctx.GlobalIsSet(TxPoolKeepLocalsFlag.Name)

	if ctx.GlobalIsSet(TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(TxPoolLifetimeFlag.Name)
	}

	// PN specific txpool setting
	if NodeTypeFlag.Value == "pn" {
		cfg.EnableSpamThrottlerAtRuntime = !ctx.GlobalIsSet(TxPoolSpamThrottlerDisableFlag.Name)
	}
}

// getNetworkId returns the associated network ID with whether or not the network is private.
func getNetworkId(ctx *cli.Context) (uint64, bool) {
	if ctx.GlobalIsSet(BaobabFlag.Name) && ctx.GlobalIsSet(CypressFlag.Name) {
		log.Fatalf("--baobab and --cypress must not be set together")
	}
	if ctx.GlobalIsSet(BaobabFlag.Name) && ctx.GlobalIsSet(NetworkIdFlag.Name) {
		log.Fatalf("--baobab and --networkid must not be set together")
	}
	if ctx.GlobalIsSet(CypressFlag.Name) && ctx.GlobalIsSet(NetworkIdFlag.Name) {
		log.Fatalf("--cypress and --networkid must not be set together")
	}

	switch {
	case ctx.GlobalIsSet(CypressFlag.Name):
		logger.Info("Cypress network ID is set", "networkid", params.CypressNetworkId)
		return params.CypressNetworkId, false
	case ctx.GlobalIsSet(BaobabFlag.Name):
		logger.Info("Baobab network ID is set", "networkid", params.BaobabNetworkId)
		return params.BaobabNetworkId, false
	case ctx.GlobalIsSet(NetworkIdFlag.Name):
		networkId := ctx.GlobalUint64(NetworkIdFlag.Name)
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
	cfg.TxResendInterval = ctx.GlobalUint64(TxResendIntervalFlag.Name)
	if cfg.TxResendInterval == 0 {
		cfg.TxResendInterval = cn.DefaultTxResendInterval
	}

	cfg.TxResendCount = ctx.GlobalInt(TxResendCountFlag.Name)
	if cfg.TxResendCount < cn.DefaultMaxResendTxCount {
		cfg.TxResendCount = cn.DefaultMaxResendTxCount
	}
	cfg.TxResendUseLegacy = ctx.GlobalBool(TxResendUseLegacyFlag.Name)
	logger.Debug("TxResend config", "Interval", cfg.TxResendInterval, "TxResendCount", cfg.TxResendCount, "UseLegacy", cfg.TxResendUseLegacy)
}

func (kCfg *KlayConfig) SetChainDataFetcherConfig(ctx *cli.Context) {
	cfg := &kCfg.ChainDataFetcher
	if ctx.GlobalBool(EnableChainDataFetcherFlag.Name) {
		cfg.EnabledChainDataFetcher = true

		if ctx.GlobalIsSet(ChainDataFetcherNoDefault.Name) {
			cfg.NoDefaultStart = true
		}
		if ctx.GlobalIsSet(ChainDataFetcherNumHandlers.Name) {
			cfg.NumHandlers = ctx.GlobalInt(ChainDataFetcherNumHandlers.Name)
		}
		if ctx.GlobalIsSet(ChainDataFetcherJobChannelSize.Name) {
			cfg.JobChannelSize = ctx.GlobalInt(ChainDataFetcherJobChannelSize.Name)
		}
		if ctx.GlobalIsSet(ChainDataFetcherChainEventSizeFlag.Name) {
			cfg.BlockChannelSize = ctx.GlobalInt(ChainDataFetcherChainEventSizeFlag.Name)
		}
		if ctx.GlobalIsSet(ChainDataFetcherMaxProcessingDataSize.Name) {
			cfg.MaxProcessingDataSize = ctx.GlobalInt(ChainDataFetcherMaxProcessingDataSize.Name)
		}

		mode := ctx.GlobalString(ChainDataFetcherMode.Name)
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
	if !ctx.GlobalIsSet(ChainDataFetcherKASDBHostFlag.Name) {
		logger.Crit("DBHost must be set !", "key", ChainDataFetcherKASDBHostFlag.Name)
	}
	if !ctx.GlobalIsSet(ChainDataFetcherKASDBUserFlag.Name) {
		logger.Crit("DBUser must be set !", "key", ChainDataFetcherKASDBUserFlag.Name)
	}
	if !ctx.GlobalIsSet(ChainDataFetcherKASDBPasswordFlag.Name) {
		logger.Crit("DBPassword must be set !", "key", ChainDataFetcherKASDBPasswordFlag.Name)
	}
	if !ctx.GlobalIsSet(ChainDataFetcherKASDBNameFlag.Name) {
		logger.Crit("DBName must be set !", "key", ChainDataFetcherKASDBNameFlag.Name)
	}
}

// NOTE-klaytn
// Deprecated: KASConfig is not used anymore.
func checkKASCacheInvalidationConfigs(ctx *cli.Context) {
	if !ctx.GlobalIsSet(ChainDataFetcherKASCacheURLFlag.Name) {
		logger.Crit("The cache invalidation url is not set")
	}
	if !ctx.GlobalIsSet(ChainDataFetcherKASBasicAuthParamFlag.Name) {
		logger.Crit("The authorization is not set")
	}
	if !ctx.GlobalIsSet(ChainDataFetcherKASXChainIdFlag.Name) {
		logger.Crit("The x-chain-id is not set")
	}
}

// NOTE-klaytn
// Deprecated: KASConfig is not used anymore.
func makeKASConfig(ctx *cli.Context) *kas.KASConfig {
	kasConfig := kas.DefaultKASConfig

	checkKASDBConfigs(ctx)
	kasConfig.DBHost = ctx.GlobalString(ChainDataFetcherKASDBHostFlag.Name)
	kasConfig.DBPort = ctx.GlobalString(ChainDataFetcherKASDBPortFlag.Name)
	kasConfig.DBUser = ctx.GlobalString(ChainDataFetcherKASDBUserFlag.Name)
	kasConfig.DBPassword = ctx.GlobalString(ChainDataFetcherKASDBPasswordFlag.Name)
	kasConfig.DBName = ctx.GlobalString(ChainDataFetcherKASDBNameFlag.Name)

	if ctx.GlobalBool(ChainDataFetcherKASCacheUse.Name) {
		checkKASCacheInvalidationConfigs(ctx)
		kasConfig.CacheUse = true
		kasConfig.CacheInvalidationURL = ctx.GlobalString(ChainDataFetcherKASCacheURLFlag.Name)
		kasConfig.BasicAuthParam = ctx.GlobalString(ChainDataFetcherKASBasicAuthParamFlag.Name)
		kasConfig.XChainId = ctx.GlobalString(ChainDataFetcherKASXChainIdFlag.Name)
	}
	return kasConfig
}

func makeKafkaConfig(ctx *cli.Context) *kafka.KafkaConfig {
	kafkaConfig := kafka.GetDefaultKafkaConfig()
	if ctx.GlobalIsSet(ChainDataFetcherKafkaBrokersFlag.Name) {
		kafkaConfig.Brokers = ctx.GlobalStringSlice(ChainDataFetcherKafkaBrokersFlag.Name)
	} else {
		logger.Crit("The kafka brokers must be set")
	}
	kafkaConfig.TopicEnvironmentName = ctx.GlobalString(ChainDataFetcherKafkaTopicEnvironmentFlag.Name)
	kafkaConfig.TopicResourceName = ctx.GlobalString(ChainDataFetcherKafkaTopicResourceFlag.Name)
	kafkaConfig.Partitions = int32(ctx.GlobalInt64(ChainDataFetcherKafkaPartitionsFlag.Name))
	kafkaConfig.Replicas = int16(ctx.GlobalInt64(ChainDataFetcherKafkaReplicasFlag.Name))
	kafkaConfig.SaramaConfig.Producer.MaxMessageBytes = ctx.GlobalInt(ChainDataFetcherKafkaMaxMessageBytesFlag.Name)
	kafkaConfig.SegmentSizeBytes = ctx.GlobalInt(ChainDataFetcherKafkaSegmentSizeBytesFlag.Name)
	kafkaConfig.MsgVersion = ctx.GlobalString(ChainDataFetcherKafkaMessageVersionFlag.Name)
	kafkaConfig.ProducerId = ctx.GlobalString(ChainDataFetcherKafkaProducerIdFlag.Name)
	requiredAcks := sarama.RequiredAcks(ctx.GlobalInt(ChainDataFetcherKafkaRequiredAcksFlag.Name))
	if requiredAcks != sarama.NoResponse && requiredAcks != sarama.WaitForLocal && requiredAcks != sarama.WaitForAll {
		logger.Crit("not supported requiredAcks. it must be NoResponse(0), WaitForLocal(1), or WaitForAll(-1)", "given", requiredAcks)
	}
	kafkaConfig.SaramaConfig.Producer.RequiredAcks = requiredAcks
	return kafkaConfig
}

func (kCfg *KlayConfig) SetDBSyncerConfig(ctx *cli.Context) {
	cfg := &kCfg.DB
	if ctx.GlobalBool(EnableDBSyncerFlag.Name) {
		cfg.EnabledDBSyncer = true

		if ctx.GlobalIsSet(DBHostFlag.Name) {
			dbhost := ctx.GlobalString(DBHostFlag.Name)
			cfg.DBHost = dbhost
		} else {
			logger.Crit("DBHost must be set !", "key", DBHostFlag.Name)
		}
		if ctx.GlobalIsSet(DBPortFlag.Name) {
			dbports := ctx.GlobalString(DBPortFlag.Name)
			cfg.DBPort = dbports
		}
		if ctx.GlobalIsSet(DBUserFlag.Name) {
			dbuser := ctx.GlobalString(DBUserFlag.Name)
			cfg.DBUser = dbuser
		} else {
			logger.Crit("DBUser must be set !", "key", DBUserFlag.Name)
		}
		if ctx.GlobalIsSet(DBPasswordFlag.Name) {
			dbpasswd := ctx.GlobalString(DBPasswordFlag.Name)
			cfg.DBPassword = dbpasswd
		} else {
			logger.Crit("DBPassword must be set !", "key", DBPasswordFlag.Name)
		}
		if ctx.GlobalIsSet(DBNameFlag.Name) {
			dbname := ctx.GlobalString(DBNameFlag.Name)
			cfg.DBName = dbname
		} else {
			logger.Crit("DBName must be set !", "key", DBNameFlag.Name)
		}
		if ctx.GlobalBool(EnabledLogModeFlag.Name) {
			cfg.EnabledLogMode = true
		}
		if ctx.GlobalIsSet(MaxIdleConnsFlag.Name) {
			cfg.MaxIdleConns = ctx.GlobalInt(MaxIdleConnsFlag.Name)
		}
		if ctx.GlobalIsSet(MaxOpenConnsFlag.Name) {
			cfg.MaxOpenConns = ctx.GlobalInt(MaxOpenConnsFlag.Name)
		}
		if ctx.GlobalIsSet(ConnMaxLifeTimeFlag.Name) {
			cfg.ConnMaxLifetime = ctx.GlobalDuration(ConnMaxLifeTimeFlag.Name)
		}
		if ctx.GlobalIsSet(DBSyncerModeFlag.Name) {
			cfg.Mode = strings.ToLower(ctx.GlobalString(DBSyncerModeFlag.Name))
		}
		if ctx.GlobalIsSet(GenQueryThreadFlag.Name) {
			cfg.GenQueryThread = ctx.GlobalInt(GenQueryThreadFlag.Name)
		}
		if ctx.GlobalIsSet(InsertThreadFlag.Name) {
			cfg.InsertThread = ctx.GlobalInt(InsertThreadFlag.Name)
		}
		if ctx.GlobalIsSet(BulkInsertSizeFlag.Name) {
			cfg.BulkInsertSize = ctx.GlobalInt(BulkInsertSizeFlag.Name)
		}
		if ctx.GlobalIsSet(EventModeFlag.Name) {
			cfg.EventMode = strings.ToLower(ctx.GlobalString(EventModeFlag.Name))
		}
		if ctx.GlobalIsSet(MaxBlockDiffFlag.Name) {
			cfg.MaxBlockDiff = ctx.GlobalUint64(MaxBlockDiffFlag.Name)
		}
		if ctx.GlobalIsSet(BlockSyncChannelSizeFlag.Name) {
			cfg.BlockChannelSize = ctx.GlobalInt(BlockSyncChannelSizeFlag.Name)
		}
	}
}

func (kCfg *KlayConfig) SetServiceChainConfig(ctx *cli.Context) {
	cfg := &kCfg.ServiceChain

	// bridge service
	if ctx.GlobalBool(MainBridgeFlag.Name) {
		cfg.EnabledMainBridge = true
		cfg.MainBridgePort = fmt.Sprintf(":%d", ctx.GlobalInt(MainBridgeListenPortFlag.Name))
	}

	if ctx.GlobalBool(SubBridgeFlag.Name) {
		cfg.EnabledSubBridge = true
		cfg.SubBridgePort = fmt.Sprintf(":%d", ctx.GlobalInt(SubBridgeListenPortFlag.Name))
	}

	cfg.Anchoring = ctx.GlobalBool(ServiceChainAnchoringFlag.Name)
	cfg.ChildChainIndexing = ctx.GlobalIsSet(ChildChainIndexingFlag.Name)
	cfg.AnchoringPeriod = ctx.GlobalUint64(AnchoringPeriodFlag.Name)
	cfg.SentChainTxsLimit = ctx.GlobalUint64(SentChainTxsLimit.Name)
	cfg.ParentChainID = ctx.GlobalUint64(ParentChainIDFlag.Name)
	cfg.VTRecovery = ctx.GlobalBool(VTRecoveryFlag.Name)
	cfg.VTRecoveryInterval = ctx.GlobalUint64(VTRecoveryIntervalFlag.Name)
	cfg.ServiceChainConsensus = ServiceChainConsensusFlag.Value
	cfg.ServiceChainParentOperatorGasLimit = ctx.GlobalUint64(ServiceChainParentOperatorTxGasLimitFlag.Name)
	cfg.ServiceChainChildOperatorGasLimit = ctx.GlobalUint64(ServiceChainChildOperatorTxGasLimitFlag.Name)

	cfg.KASAnchor = ctx.GlobalBool(KASServiceChainAnchorFlag.Name)
	if cfg.KASAnchor {
		cfg.KASAnchorPeriod = ctx.GlobalUint64(KASServiceChainAnchorPeriodFlag.Name)
		if cfg.KASAnchorPeriod == 0 {
			cfg.KASAnchorPeriod = 1
			logger.Warn("KAS anchor period is set by 1")
		}

		cfg.KASAnchorUrl = ctx.GlobalString(KASServiceChainAnchorUrlFlag.Name)
		if cfg.KASAnchorUrl == "" {
			logger.Crit("KAS anchor url should be set", "key", KASServiceChainAnchorUrlFlag.Name)
		}

		cfg.KASAnchorOperator = ctx.GlobalString(KASServiceChainAnchorOperatorFlag.Name)
		if cfg.KASAnchorOperator == "" {
			logger.Crit("KAS anchor operator should be set", "key", KASServiceChainAnchorOperatorFlag.Name)
		}

		cfg.KASAccessKey = ctx.GlobalString(KASServiceChainAccessKeyFlag.Name)
		if cfg.KASAccessKey == "" {
			logger.Crit("KAS access key should be set", "key", KASServiceChainAccessKeyFlag.Name)
		}

		cfg.KASSecretKey = ctx.GlobalString(KASServiceChainSecretKeyFlag.Name)
		if cfg.KASSecretKey == "" {
			logger.Crit("KAS secret key should be set", "key", KASServiceChainSecretKeyFlag.Name)
		}

		cfg.KASXChainId = ctx.GlobalString(KASServiceChainXChainIdFlag.Name)
		if cfg.KASXChainId == "" {
			logger.Crit("KAS x-chain-id should be set", "key", KASServiceChainXChainIdFlag.Name)
		}

		cfg.KASAnchorRequestTimeout = ctx.GlobalDuration(KASServiceChainAnchorRequestTimeoutFlag.Name)
	}

	cfg.DataDir = kCfg.Node.DataDir
	cfg.Name = kCfg.Node.Name
}
