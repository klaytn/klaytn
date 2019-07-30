// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from cmd/utils/flags.go (2018/06/04).
// Modified and improved for the klaytn development.

package utils

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/dbsyncer"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/metrics"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/nat"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/node/sc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func InitHelper() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
}

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	//app.Authors = nil
	app.Email = ""
	app.Version = params.Version
	if len(gitCommit) >= 8 {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	return app
}

var (
	// General settings
	NetworkTypeFlag = cli.StringFlag{
		Name:  "networktype",
		Usage: "Klaytn network type (main-net (mn), service chain-net (scn))",
		Value: "mn",
	}
	DbTypeFlag = cli.StringFlag{
		Name:  "dbtype",
		Usage: `Blockchain storage database type ("leveldb", "badger")`,
		Value: "leveldb",
	}
	SrvTypeFlag = cli.StringFlag{
		Name:  "srvtype",
		Usage: `json rpc server type ("http", "fasthttp")`,
		Value: "fasthttp",
	}
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DirectoryString{node.DefaultDataDir()},
	}
	KeyStoreDirFlag = DirectoryFlag{
		Name:  "keystore",
		Usage: "Directory for the keystore (default = inside the datadir)",
	}
	// TODO-Klaytn-Bootnode: redefine networkid
	NetworkIdFlag = cli.Uint64Flag{
		Name:  "networkid",
		Usage: "Network identifier (integer, 1=MainNet (Not yet launched), 1000=Aspen, 1001=Baobab)",
		Value: cn.DefaultConfig.NetworkId,
	}
	IdentityFlag = cli.StringFlag{
		Name:  "identity",
		Usage: "Custom node name",
	}
	DocRootFlag = DirectoryFlag{
		Name:  "docroot",
		Usage: "Document Root for HTTPClient file scheme",
		Value: DirectoryString{homeDir()},
	}
	defaultSyncMode = cn.DefaultConfig.SyncMode
	SyncModeFlag    = TextMarshalerFlag{
		Name:  "syncmode",
		Usage: `Blockchain sync mode (only "full" is supported)`,
		Value: &defaultSyncMode,
	}
	GCModeFlag = cli.StringFlag{
		Name:  "gcmode",
		Usage: `Blockchain garbage collection mode ("full", "archive")`,
		Value: "full",
	}
	LightKDFFlag = cli.BoolFlag{
		Name:  "lightkdf",
		Usage: "Reduce key-derivation RAM & CPU usage at some expense of KDF strength",
	}
	// Transaction pool settings
	TxPoolNoLocalsFlag = cli.BoolFlag{
		Name:  "txpool.nolocals",
		Usage: "Disables price exemptions for locally submitted transactions",
	}
	TxPoolJournalFlag = cli.StringFlag{
		Name:  "txpool.journal",
		Usage: "Disk journal for local transaction to survive node restarts",
		Value: blockchain.DefaultTxPoolConfig.Journal,
	}
	TxPoolJournalIntervalFlag = cli.DurationFlag{
		Name:  "txpool.journal-interval",
		Usage: "Time interval to regenerate the local transaction journal",
		Value: blockchain.DefaultTxPoolConfig.JournalInterval,
	}
	TxPoolPriceLimitFlag = cli.Uint64Flag{
		Name:  "txpool.pricelimit",
		Usage: "Minimum gas price limit to enforce for acceptance into the pool",
		Value: cn.DefaultConfig.TxPool.PriceLimit,
	}
	TxPoolPriceBumpFlag = cli.Uint64Flag{
		Name:  "txpool.pricebump",
		Usage: "Price bump percentage to replace an already existing transaction",
		Value: cn.DefaultConfig.TxPool.PriceBump,
	}
	TxPoolExecSlotsAccountFlag = cli.Uint64Flag{
		Name:  "txpool.exec-slots.account",
		Usage: "Number of executable transaction slots guaranteed per account",
		Value: cn.DefaultConfig.TxPool.ExecSlotsAccount,
	}
	TxPoolExecSlotsAllFlag = cli.Uint64Flag{
		Name:  "txpool.exec-slots.all",
		Usage: "Maximum number of executable transaction slots for all accounts",
		Value: cn.DefaultConfig.TxPool.ExecSlotsAll,
	}
	TxPoolNonExecSlotsAccountFlag = cli.Uint64Flag{
		Name:  "txpool.nonexec-slots.account",
		Usage: "Maximum number of non-executable transaction slots permitted per account",
		Value: cn.DefaultConfig.TxPool.NonExecSlotsAccount,
	}
	TxPoolNonExecSlotsAllFlag = cli.Uint64Flag{
		Name:  "txpool.nonexec-slots.all",
		Usage: "Maximum number of non-executable transaction slots for all accounts",
		Value: cn.DefaultConfig.TxPool.NonExecSlotsAll,
	}
	KeepLocalsFlag = cli.BoolFlag{
		Name:  "txpool.keeplocals",
		Usage: "Disables removing timed-out local transactions",
	}
	TxPoolLifetimeFlag = cli.DurationFlag{
		Name:  "txpool.lifetime",
		Usage: "Maximum amount of time non-executable transaction are queued",
		Value: cn.DefaultConfig.TxPool.Lifetime,
	}
	// Performance tuning settings
	StateDBCachingFlag = cli.BoolFlag{
		Name:  "statedb.use-cache",
		Usage: "Enables caching of state objects in stateDB",
	}
	NoPartitionedDBFlag = cli.BoolFlag{
		Name:  "db.no-partitioning",
		Usage: "Disable partitioned databases for persistent storage",
	}
	NumStateTriePartitionsFlag = cli.UintFlag{
		Name:  "db.num-statetrie-partitions",
		Usage: "Number of internal partitions of state trie partition. Should be power of 2",
		Value: 4,
	}
	LevelDBCacheSizeFlag = cli.IntFlag{
		Name:  "db.leveldb.cache-size",
		Usage: "Size of in-memory cache in LevelDB (MiB)",
		Value: 768,
	}
	// TODO-Klaytn-Database LevelDBCompressionTypeFlag should be removed before main-net release.
	LevelDBCompressionTypeFlag = cli.IntFlag{
		Name:  "db.leveldb.compression",
		Usage: "Determines the compression method for LevelDB. 0=AllNoCompression, 1=ReceiptOnlySnappyCompression, 2=StateTrieOnlyNoCompression, 3=AllSnappyCompression",
		Value: 0,
	}
	LevelDBNoBufferPoolFlag = cli.BoolFlag{
		Name:  "db.leveldb.no-buffer-pool",
		Usage: "Disables using buffer pool for LevelDB's block allocation",
	}
	NoParallelDBWriteFlag = cli.BoolFlag{
		Name:  "db.no-parallel-write",
		Usage: "Disables parallel writes of block data to persistent database",
	}
	TrieMemoryCacheSizeFlag = cli.IntFlag{
		Name:  "state.cache-size",
		Usage: "Size of in-memory cache of the global state (in MiB) to flush matured singleton trie nodes to disk",
		Value: 512,
	}
	TrieBlockIntervalFlag = cli.UintFlag{
		Name:  "state.block-interval",
		Usage: "An interval in terms of block number to commit the global state to disk",
		Value: blockchain.DefaultBlockInterval,
	}
	CacheTypeFlag = cli.IntFlag{
		Name:  "cache.type",
		Usage: "Cache Type: 0=LRUCache, 1=LRUShardCache, 2=FIFOCache",
		Value: int(common.DefaultCacheType),
	}
	CacheScaleFlag = cli.IntFlag{
		Name:  "cache.scale",
		Usage: "Scale of cache (cache size = preset size * scale of cache(%))",
	}
	CacheUsageLevelFlag = cli.StringFlag{
		Name:  "cache.level",
		Usage: "Set the cache usage level ('saving', 'normal', 'extreme')",
	}
	MemorySizeFlag = cli.IntFlag{
		Name:  "cache.memory",
		Usage: "Set the physical RAM size (GB, Default: 16GB)",
	}
	CacheWriteThroughFlag = cli.BoolFlag{
		Name:  "cache.writethrough",
		Usage: "Enables write-through writing to database and cache for certain types of cache.",
	}
	TxPoolStateCacheFlag = cli.BoolFlag{
		Name:  "statedb.use-txpool-cache",
		Usage: "Enables caching of nonce and balance for txpool.",
	}
	TrieCacheLimitFlag = cli.IntFlag{
		Name:  "state.trie-cache-limit",
		Usage: "Memory allowance (MB) to use for caching trie nodes in memory",
		Value: 4096,
	}

	SenderTxHashIndexingFlag = cli.BoolFlag{
		Name:  "sendertxhashindexing",
		Usage: "Enables storing mapping information of senderTxHash to txHash",
	}
	ChildChainIndexingFlag = cli.BoolFlag{
		Name:  "childchainindexing",
		Usage: "Enables storing transaction hash of child chain transaction for fast access to child chain data",
	}
	TargetGasLimitFlag = cli.Uint64Flag{
		Name:  "targetgaslimit",
		Usage: "Target gas limit sets the artificial target gas floor for the blocks to mine",
		Value: params.GenesisGasLimit,
	}
	ServiceChainSignerFlag = cli.StringFlag{
		Name:  "scsigner",
		Usage: "Public address for signing blocks in the service chain (default = first account created)",
		Value: "0",
	}
	RewardbaseFlag = cli.StringFlag{
		Name:  "rewardbase",
		Usage: "Public address for block consensus rewards (default = first account created)",
		Value: "0",
	}
	ExtraDataFlag = cli.StringFlag{
		Name:  "extradata",
		Usage: "Block extra data set by the work (default = client version)",
	}

	TxResendIntervalFlag = cli.Uint64Flag{
		Name:  "txresend.interval",
		Usage: "Set the transaction resend interval in seconds",
		Value: uint64(cn.DefaultTxResendInterval),
	}
	TxResendCountFlag = cli.IntFlag{
		Name:  "txresend.max-count",
		Usage: "Set the max count of resending transactions",
		Value: cn.DefaultMaxResendTxCount,
	}
	//TODO-Klaytn-RemoveLater Remove this flag when we are confident with the new transaction resend logic
	TxResendUseLegacyFlag = cli.BoolFlag{
		Name:  "txresend.use-legacy",
		Usage: "Enable the legacy transaction resend logic (For testing only)",
	}
	// Account settings
	UnlockedAccountFlag = cli.StringFlag{
		Name:  "unlock",
		Usage: "Comma separated list of accounts to unlock",
		Value: "",
	}
	PasswordFileFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Password file to use for non-interactive password input",
		Value: "",
	}

	VMEnableDebugFlag = cli.BoolFlag{
		Name:  "vmdebug",
		Usage: "Record information useful for VM and contract debugging",
	}
	VMLogTargetFlag = cli.IntFlag{
		Name:  "vmlog",
		Usage: "Set the output target of vmlog precompiled contract (0: no output, 1: file, 2: stdout, 3: both)",
		Value: 0,
	}

	// Logging and debug settings
	MetricsEnabledFlag = cli.BoolFlag{
		Name:  metrics.MetricsEnabledFlag,
		Usage: "Enable metrics collection and reporting",
	}
	PrometheusExporterFlag = cli.BoolFlag{
		Name:  metrics.PrometheusExporterFlag,
		Usage: "Enable prometheus exporter",
	}
	PrometheusExporterPortFlag = cli.IntFlag{
		Name:  metrics.PrometheusExporterPortFlag,
		Usage: "Prometheus exporter listening port",
		Value: 61001,
	}
	// RPC settings
	RPCEnabledFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Enable the HTTP-RPC server",
	}
	RPCListenAddrFlag = cli.StringFlag{
		Name:  "rpcaddr",
		Usage: "HTTP-RPC server listening interface",
		Value: node.DefaultHTTPHost,
	}
	RPCPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "HTTP-RPC server listening port",
		Value: node.DefaultHTTPPort,
	}
	RPCCORSDomainFlag = cli.StringFlag{
		Name:  "rpccorsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value: "",
	}
	RPCVirtualHostsFlag = cli.StringFlag{
		Name:  "rpcvhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		Value: strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
	}
	RPCApiFlag = cli.StringFlag{
		Name:  "rpcapi",
		Usage: "API's offered over the HTTP-RPC interface",
		Value: "",
	}
	IPCDisabledFlag = cli.BoolFlag{
		Name:  "ipcdisable",
		Usage: "Disable the IPC-RPC server",
	}
	IPCPathFlag = DirectoryFlag{
		Name:  "ipcpath",
		Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
	}
	WSEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable the WS-RPC server",
	}
	WSListenAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS-RPC server listening interface",
		Value: node.DefaultWSHost,
	}
	WSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS-RPC server listening port",
		Value: node.DefaultWSPort,
	}
	GRPCEnabledFlag = cli.BoolFlag{
		Name:  "grpc",
		Usage: "Enable the gRPC server",
	}
	GRPCListenAddrFlag = cli.StringFlag{
		Name:  "grpcaddr",
		Usage: "gRPC server listening interface",
		Value: node.DefaultGRPCHost,
	}
	GRPCPortFlag = cli.IntFlag{
		Name:  "grpcport",
		Usage: "gRPC server listening port",
		Value: node.DefaultGRPCPort,
	}
	WSApiFlag = cli.StringFlag{
		Name:  "wsapi",
		Usage: "API's offered over the WS-RPC interface",
		Value: "",
	}
	WSAllowedOriginsFlag = cli.StringFlag{
		Name:  "wsorigins",
		Usage: "Origins from which to accept websockets requests",
		Value: "",
	}
	ExecFlag = cli.StringFlag{
		Name:  "exec",
		Usage: "Execute JavaScript statement",
	}
	PreloadJSFlag = cli.StringFlag{
		Name:  "preload",
		Usage: "Comma separated list of JavaScript files to preload into the console",
	}

	// Network Settings
	NodeTypeFlag = cli.StringFlag{
		Name:  "nodetype",
		Usage: "Klaytn node type (consensus node (cn), proxy node (pn), endpoint node (en))",
		Value: "en",
	}
	MaxConnectionsFlag = cli.IntFlag{
		Name:  "maxconnections",
		Usage: "Maximum number of physical connections. All single channel peers can be maxconnections peers. All multi channel peers can be maxconnections/2 peers. (network disabled if set to 0)",
		Value: node.DefaultMaxPhysicalConnections,
	}
	MaxPendingPeersFlag = cli.IntFlag{
		Name:  "maxpendpeers",
		Usage: "Maximum number of pending connection attempts (defaults used if set to 0)",
		Value: 0,
	}
	ListenPortFlag = cli.IntFlag{
		Name:  "port",
		Usage: "Network listening port",
		Value: node.DefaultP2PPort,
	}
	SubListenPortFlag = cli.IntFlag{
		Name:  "subport",
		Usage: "Network sub listening port",
		Value: node.DefaultP2PSubPort,
	}
	MultiChannelUseFlag = cli.BoolFlag{
		Name:  "multichannel",
		Usage: "Create a dedicated channel for block propagation",
	}
	BootnodesFlag = cli.StringFlag{
		Name:  "bootnodes",
		Usage: "Comma separated kni URLs for P2P discovery bootstrap",
		Value: "",
	}
	NodeKeyFileFlag = cli.StringFlag{
		Name:  "nodekey",
		Usage: "P2P node key file",
	}
	NodeKeyHexFlag = cli.StringFlag{
		Name:  "nodekeyhex",
		Usage: "P2P node key as hex (for testing)",
	}
	NATFlag = cli.StringFlag{
		Name:  "nat",
		Usage: "NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>)",
		Value: "any",
	}
	NoDiscoverFlag = cli.BoolFlag{
		Name:  "nodiscover",
		Usage: "Disables the peer discovery mechanism (manual peer addition)",
	}
	NetrestrictFlag = cli.StringFlag{
		Name:  "netrestrict",
		Usage: "Restricts network communication to the given IP network (CIDR masks)",
	}
	AnchoringPeriodFlag = cli.Uint64Flag{
		Name:  "chaintxperiod",
		Usage: "The period to make and send a chain transaction to the main chain",
		Value: 1,
	}
	SentChainTxsLimit = cli.Uint64Flag{
		Name:  "chaintxlimit",
		Usage: "Number of service chain transactions stored for resending",
		Value: 100,
	}
	RWTimerIntervalFlag = cli.Uint64Flag{
		Name:  "rwtimerinterval",
		Usage: "Interval of using rw timer to check if it works well",
		Value: 1000,
	}
	RWTimerWaitTimeFlag = cli.DurationFlag{
		Name:  "rwtimerwaittime",
		Usage: "Wait time the rw timer waits for message writing",
		Value: 15 * time.Second,
	}

	// ATM the url is left to the user and deployment to
	JSpathFlag = cli.StringFlag{
		Name:  "jspath",
		Usage: "JavaScript root path for `loadScript`",
		Value: ".",
	}
	CypressFlag = cli.BoolFlag{
		Name:  "cypress",
		Usage: "Pre-configured Klaytn Cypress network",
	}
	// Baobab bootnodes setting
	BaobabFlag = cli.BoolFlag{
		Name:  "baobab",
		Usage: "Pre-configured Klaytn baobab network",
	}
	// Bootnode's settings
	AuthorizedNodesFlag = cli.StringFlag{
		Name:  "authorized-nodes",
		Usage: "Comma separated kni URLs for authorized nodes list",
		Value: "",
	}
	//TODO-Klaytn-Bootnode the boodnode flags should be updated when it is implemented
	BNAddrFlag = cli.StringFlag{
		Name:  "bnaddr",
		Usage: `udp address to use node discovery`,
		Value: ":32323",
	}
	GenKeyFlag = cli.StringFlag{
		Name:  "genkey",
		Usage: "generate a node private key and write to given filename",
	}
	WriteAddressFlag = cli.BoolFlag{
		Name:  "writeaddress",
		Usage: `write out the node's public key which is given by "--nodekeyfile" or "--nodekeyhex"`,
	}
	// ServiceChain's settings
	MainBridgeFlag = cli.BoolFlag{
		Name:  "mainbridge",
		Usage: "Enable main bridge service for service chain",
	}
	SubBridgeFlag = cli.BoolFlag{
		Name:  "subbridge",
		Usage: "Enable sub bridge service for service chain",
	}
	MainBridgeListenPortFlag = cli.IntFlag{
		Name:  "mainbridgeport",
		Usage: "main bridge listen port",
		Value: 50505,
	}
	SubBridgeListenPortFlag = cli.IntFlag{
		Name:  "subbridgeport",
		Usage: "sub bridge listen port",
		Value: 50506,
	}
	ParentChainIDFlag = cli.IntFlag{
		Name:  "parentchainid",
		Usage: "parent chain ID",
		Value: 8217, // Klaytn mainnet chain ID
	}
	VTRecoveryFlag = cli.BoolFlag{
		Name:  "vtrecovery",
		Usage: "Enable value transfer recovery (default: false)",
	}
	VTRecoveryIntervalFlag = cli.Uint64Flag{
		Name:  "vtrecoveryinterval",
		Usage: "Set the value transfer recovery interval (seconds)",
		Value: 60,
	}
	ServiceChainNewAccountFlag = cli.BoolFlag{
		Name:  "scnewaccount",
		Usage: "Enable account creation for the service chain (default: false). If set true, generated account can't be synced with the main chain.",
	}
	ServiceChainConsensusFlag = cli.StringFlag{
		Name:  "scconsensus",
		Usage: "Set the service chain consensus (\"clique\", \"istanbul\")",
		Value: "clique",
	}
	// DBSyncer
	EnableDBSyncerFlag = cli.BoolFlag{
		Name:  "dbsyncer",
		Usage: "Enable the DBSyncer",
	}
	DBHostFlag = cli.StringFlag{
		Name:  "dbsyncer.db.host",
		Usage: "db.host in dbsyncer",
	}
	DBPortFlag = cli.StringFlag{
		Name:  "dbsyncer.db.port",
		Usage: "db.port in dbsyncer",
		Value: "3306",
	}
	DBNameFlag = cli.StringFlag{
		Name:  "dbsyncer.db.name",
		Usage: "db.name in dbsyncer",
	}
	DBUserFlag = cli.StringFlag{
		Name:  "dbsyncer.db.user",
		Usage: "db.user in dbsyncer",
	}
	DBPasswordFlag = cli.StringFlag{
		Name:  "dbsyncer.db.password",
		Usage: "db.password in dbsyncer",
	}
	EnabledLogModeFlag = cli.BoolFlag{
		Name:  "dbsyncer.logmode",
		Usage: "Enable the dbsyncer logmode",
	}
	MaxIdleConnsFlag = cli.IntFlag{
		Name:  "dbsyncer.db.max.idle",
		Usage: "The maximum number of connections in the idle connection pool",
		Value: 50,
	}
	MaxOpenConnsFlag = cli.IntFlag{
		Name:  "dbsyncer.db.max.open",
		Usage: "The maximum number of open connections to the database",
		Value: 30,
	}
	ConnMaxLifeTimeFlag = cli.DurationFlag{
		Name:  "dbsyncer.db.max.lifetime",
		Usage: "The maximum amount of time a connection may be reused (default : 1h), ex: 300ms, 2h45m, 60s, ...",
		Value: 1 * time.Hour,
	}
	BlockSyncChannelSizeFlag = cli.IntFlag{
		Name:  "dbsyncer.block.channel.size",
		Usage: "Block received channel size",
		Value: 5,
	}
	DBSyncerModeFlag = cli.StringFlag{
		Name:  "dbsyncer.mode",
		Usage: "The mode of dbsyncer is way which handle block/tx data to insert db (multi, single, context)",
		Value: "multi",
	}
	GenQueryThreadFlag = cli.IntFlag{
		Name:  "dbsyncer.genquery.th",
		Usage: "The amount of thread of generation query in multi mode",
		Value: 50,
	}
	InsertThreadFlag = cli.IntFlag{
		Name:  "dbsyncer.insert.th",
		Usage: "The amount of thread of insert operation in multi mode",
		Value: 30,
	}
	BulkInsertSizeFlag = cli.IntFlag{
		Name:  "dbsyncer.bulk.size",
		Usage: "The amount of row for bulk-insert",
		Value: 200,
	}
	EventModeFlag = cli.StringFlag{
		Name:  "dbsyncer.event.mode",
		Usage: "The way how to sync all block or last block (block, head)",
		Value: "head",
	}
	MaxBlockDiffFlag = cli.Uint64Flag{
		Name:  "dbsyncer.max.block.diff",
		Usage: "The maximum difference between current block and event block. 0 means off",
		Value: 0,
	}

	// TODO-Klaytn-Bootnode: Add bootnode's metric options
	// TODO-Klaytn-Bootnode: Implements bootnode's RPC
)

// MakeDataDir retrieves the currently requested data directory, terminating
// if none (or the empty string) is specified. If the node is starting a baobab,
// the a subdirectory of the specified datadir will be used.
func MakeDataDir(ctx *cli.Context) string {
	if path := ctx.GlobalString(DataDirFlag.Name); path != "" {
		if ctx.GlobalBool(BaobabFlag.Name) {
			return filepath.Join(path, "baobab")
		}
		return path
	}
	log.Fatalf("Cannot determine default data directory, please set manually (--datadir)")
	return ""
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

// setNodeUserIdent creates the user identifier from CLI flags.
func setNodeUserIdent(ctx *cli.Context, cfg *node.Config) {
	if identity := ctx.GlobalString(IdentityFlag.Name); len(identity) > 0 {
		cfg.UserIdent = identity
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
		logger.Info("Cypress bootnodes are set")
		urls = params.MainnetBootnodes[cfg.ConnectionType].Addrs
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

// splitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func splitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
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
		cfg.HTTPCors = splitAndTrim(ctx.GlobalString(RPCCORSDomainFlag.Name))
	}
	if ctx.GlobalIsSet(RPCApiFlag.Name) {
		cfg.HTTPModules = splitAndTrim(ctx.GlobalString(RPCApiFlag.Name))
	}
	if ctx.GlobalIsSet(RPCVirtualHostsFlag.Name) {
		cfg.HTTPVirtualHosts = splitAndTrim(ctx.GlobalString(RPCVirtualHostsFlag.Name))
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
		cfg.WSOrigins = splitAndTrim(ctx.GlobalString(WSAllowedOriginsFlag.Name))
	}
	if ctx.GlobalIsSet(WSApiFlag.Name) {
		cfg.WSModules = splitAndTrim(ctx.GlobalString(WSApiFlag.Name))
	}
}

// setIPC creates an IPC path configuration from the set command line flags,
// returning an empty string if IPC was explicitly disabled, or the set path.
func setIPC(ctx *cli.Context, cfg *node.Config) {
	checkExclusive(ctx, IPCDisabledFlag, IPCPathFlag)
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

// MakeAddress converts an account specified directly as a hex encoded string or
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

// MakePasswordList reads password lines from the file specified by the global --password flag.
func MakePasswordList(ctx *cli.Context) []string {
	path := ctx.GlobalString(PasswordFileFlag.Name)
	if path == "" {
		return nil
	}
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read password file: %v", err)
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines
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
	if cfg.ConnectionType == node.UNKNOWNNODE {
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

	if ctx.GlobalIsSet(ListenPortFlag.Name) && (NodeTypeFlag.Value == "spn" || NodeTypeFlag.Value == "sen") {
		if !ctx.GlobalIsSet(NetworkIdFlag.Name) {
			log.Fatalf("Missing network id for the nodetype: %v", NodeTypeFlag.Value)
		}
		cfg.NoDiscovery = true
	}

	cfg.NetworkID, _ = getNetworkId(ctx)
}

func convertNodeType(nodetype string) p2p.ConnType {
	switch strings.ToLower(nodetype) {
	case "cn", "scn":
		return node.CONSENSUSNODE
	case "pn", "spn":
		return node.PROXYNODE
	case "en", "sen":
		return node.ENDPOINTNODE
	default:
		return node.UNKNOWNNODE
	}
}

// SetNodeConfig applies node-related command line flags to the config.
func SetNodeConfig(ctx *cli.Context, cfg *node.Config) {
	SetP2PConfig(ctx, &cfg.P2P)
	setIPC(ctx, cfg)

	// httptype is http or fasthttp
	if ctx.GlobalIsSet(SrvTypeFlag.Name) {
		cfg.HTTPServerType = ctx.GlobalString(SrvTypeFlag.Name)
	}

	setHTTP(ctx, cfg)
	setWS(ctx, cfg)
	setgRPC(ctx, cfg)
	setNodeUserIdent(ctx, cfg)

	cfg.DBType = ctx.GlobalString(DbTypeFlag.Name)
	cfg.DataDir = ctx.GlobalString(DataDirFlag.Name)

	if ctx.GlobalIsSet(KeyStoreDirFlag.Name) {
		cfg.KeyStoreDir = ctx.GlobalString(KeyStoreDirFlag.Name)
	}
	if ctx.GlobalIsSet(LightKDFFlag.Name) {
		cfg.UseLightweightKDF = ctx.GlobalBool(LightKDFFlag.Name)
	}
}

func setTxPool(ctx *cli.Context, cfg *blockchain.TxPoolConfig) {
	if ctx.GlobalIsSet(TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(TxPoolNoLocalsFlag.Name)
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

	cfg.KeepLocals = ctx.GlobalIsSet(KeepLocalsFlag.Name)

	if ctx.GlobalIsSet(TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(TxPoolLifetimeFlag.Name)
	}
}

// checkExclusive verifies that only a single instance of the provided flags was
// set by the user. Each flag might optionally be followed by a string type to
// specialize it further.
func checkExclusive(ctx *cli.Context, args ...interface{}) {
	set := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		// Make sure the next argument is a flag and skip if not set
		flag, ok := args[i].(cli.Flag)
		if !ok {
			panic(fmt.Sprintf("invalid argument, not cli.Flag type: %T", args[i]))
		}
		// Check if next arg extends current and expand its name if so
		name := flag.GetName()

		if i+1 < len(args) {
			switch option := args[i+1].(type) {
			case string:
				// Extended flag, expand the name and shift the arguments
				if ctx.GlobalString(flag.GetName()) == option {
					name += "=" + option
				}
				i++

			case cli.Flag:
			default:
				panic(fmt.Sprintf("invalid argument, not cli.Flag or string extension: %T", args[i+1]))
			}
		}
		// Mark the flag if it's set
		if ctx.GlobalIsSet(flag.GetName()) {
			set = append(set, "--"+name)
		}
	}
	if len(set) > 1 {
		log.Fatalf("Flags %v can't be used at the same time", strings.Join(set, ", "))
	}
}

// SetKlayConfig applies klay-related command line flags to the config.
func SetKlayConfig(ctx *cli.Context, stack *node.Node, cfg *cn.Config) {
	// TODO-Klaytn-Bootnode: better have to check conflicts about network flags when we add Klaytn's `mainnet` parameter
	// checkExclusive(ctx, DeveloperFlag, TestnetFlag, RinkebyFlag)

	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	setServiceChainSigner(ctx, ks, cfg)
	setRewardbase(ctx, ks, cfg)
	setTxPool(ctx, &cfg.TxPool)

	if ctx.GlobalIsSet(SyncModeFlag.Name) {
		cfg.SyncMode = *GlobalTextMarshaler(ctx, SyncModeFlag.Name).(*downloader.SyncMode)
		if cfg.SyncMode != downloader.FullSync {
			log.Fatalf("only syncmode=full can be used for syncmode!")
		}
	}

	cfg.NetworkId, cfg.IsPrivate = getNetworkId(ctx)

	cfg.PartitionedDB = !ctx.GlobalIsSet(NoPartitionedDBFlag.Name)
	cfg.NumStateTriePartitions = ctx.GlobalUint(NumStateTriePartitionsFlag.Name)
	if !database.IsPow2(cfg.NumStateTriePartitions) {
		log.Fatalf("--db.num-statetrie-partitions should be power of 2 but %v is not!", cfg.NumStateTriePartitions)
	}

	cfg.LevelDBCompression = database.LevelDBCompressionType(ctx.GlobalInt(LevelDBCompressionTypeFlag.Name))
	cfg.LevelDBBufferPool = !ctx.GlobalIsSet(LevelDBNoBufferPoolFlag.Name)
	cfg.LevelDBCacheSize = ctx.GlobalInt(LevelDBCacheSizeFlag.Name)

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

	common.WriteThroughCaching = ctx.GlobalIsSet(CacheWriteThroughFlag.Name)
	cfg.TxPoolStateCache = ctx.GlobalIsSet(TxPoolStateCacheFlag.Name)

	if ctx.GlobalIsSet(DocRootFlag.Name) {
		cfg.DocRoot = ctx.GlobalString(DocRootFlag.Name)
	}
	if ctx.GlobalIsSet(ExtraDataFlag.Name) {
		cfg.ExtraData = []byte(ctx.GlobalString(ExtraDataFlag.Name))
	}

	cfg.SenderTxHashIndexing = ctx.GlobalIsSet(SenderTxHashIndexingFlag.Name)
	cfg.ParallelDBWrite = !ctx.GlobalIsSet(NoParallelDBWriteFlag.Name)
	cfg.StateDBCaching = ctx.GlobalIsSet(StateDBCachingFlag.Name)
	cfg.TrieCacheLimit = ctx.GlobalInt(TrieCacheLimitFlag.Name)

	if ctx.GlobalIsSet(VMEnableDebugFlag.Name) {
		// TODO(fjl): force-enable this in --dev mode
		cfg.EnablePreimageRecording = ctx.GlobalBool(VMEnableDebugFlag.Name)
	}
	if ctx.GlobalIsSet(VMLogTargetFlag.Name) {
		if _, err := debug.Handler.SetVMLogTarget(ctx.GlobalInt(VMLogTargetFlag.Name)); err != nil {
			logger.Warn("Incorrect vmlog value", "err", err)
		}
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

// RegisterCNService adds a CN client to the stack.
func RegisterCNService(stack *node.Node, cfg *cn.Config) {
	// TODO-Klaytn add syncMode.LightSync func and add LesServer
	err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		cfg.WsEndpoint = stack.WSEndpoint()
		fullNode, err := cn.New(ctx, cfg)
		return fullNode, err
	})
	if err != nil {
		log.Fatalf("Failed to register the CN service: %v", err)
	}
}

// RegisterServiceChainService adds a ServiceChain node to the stack.
func RegisterServiceChainService(stack *node.Node, cfg *cn.Config) {
	err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		cfg.WsEndpoint = stack.WSEndpoint()
		fullNode, err := cn.NewServiceChain(ctx, cfg)
		return fullNode, err
	})
	if err != nil {
		log.Fatalf("Failed to register the SCN service: %v", err)
	}
}

func RegisterService(stack *node.Node, cfg *sc.SCConfig) {
	if cfg.EnabledMainBridge {
		err := stack.RegisterSubService(func(ctx *node.ServiceContext) (node.Service, error) {
			mainBridge, err := sc.NewMainBridge(ctx, cfg)
			return mainBridge, err
		})
		if err != nil {
			log.Fatalf("Failed to register the main bridge service: %v", err)
		}
	}

	if cfg.EnabledSubBridge {
		err := stack.RegisterSubService(func(ctx *node.ServiceContext) (node.Service, error) {
			subBridge, err := sc.NewSubBridge(ctx, cfg)
			return subBridge, err
		})
		if err != nil {
			log.Fatalf("Failed to register the sub bridge service: %v", err)
		}
	}
}

// RegisterDBSyncerService adds a DBSyncer to the stack
func RegisterDBSyncerService(stack *node.Node, cfg *dbsyncer.DBConfig) {
	if cfg.EnabledDBSyncer {
		err := stack.RegisterSubService(func(ctx *node.ServiceContext) (node.Service, error) {
			dbImporter, err := dbsyncer.NewDBSyncer(ctx, cfg)
			return dbImporter, err
		})
		if err != nil {
			log.Fatalf("Failed to register the service: %v", err)
		}
	}
}

// SetupNetwork configures the system for either the main net or some test network.
func SetupNetwork(ctx *cli.Context) {
	// TODO(fjl): move target gas limit into config
	params.TargetGasLimit = ctx.GlobalUint64(TargetGasLimitFlag.Name)
}

// MakeConsolePreloads retrieves the absolute paths for the console JavaScript
// scripts to preload before starting.
func MakeConsolePreloads(ctx *cli.Context) []string {
	// Skip preloading if there's nothing to preload
	if ctx.GlobalString(PreloadJSFlag.Name) == "" {
		return nil
	}
	// Otherwise resolve absolute paths and return them
	preloads := []string{}

	assets := ctx.GlobalString(JSpathFlag.Name)
	for _, file := range strings.Split(ctx.GlobalString(PreloadJSFlag.Name), ",") {
		preloads = append(preloads, common.AbsolutePath(assets, strings.TrimSpace(file)))
	}
	return preloads
}

// MigrateFlags sets the global flag from a local flag when it's set.
// This is a temporary function used for migrating old command/flags to the
// new format.
//
// e.g. ken account new --keystore /tmp/mykeystore --lightkdf
//
// is equivalent after calling this method with:
//
// ken --keystore /tmp/mykeystore --lightkdf account new
//
// This allows the use of the existing configuration functionality.
// When all flags are migrated this function can be removed and the existing
// configuration functionality must be changed that is uses local flags
func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))
			}
		}
		return action(ctx)
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

// getNetworkID returns the associated network ID with whether or not the network is private.
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
		logger.Info("Cypress network ID is set", "networkid", params.CypressNetworkId)
		return params.CypressNetworkId, false
	}
}
