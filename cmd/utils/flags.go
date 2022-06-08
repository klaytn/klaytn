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
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/fdlimit"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka"
	"github.com/klaytn/klaytn/datasync/dbsyncer"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/log"
	metricutils "github.com/klaytn/klaytn/metrics/utils"
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
	"gopkg.in/urfave/cli.v1"
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
		Usage: `Blockchain storage database type ("LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3")`,
		Value: "LevelDB",
	}
	SrvTypeFlag = cli.StringFlag{
		Name:  "srvtype",
		Usage: `json rpc server type ("http", "fasthttp")`,
		Value: "fasthttp",
	}
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore. This value is only used in local DB.",
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
		Value: cn.GetDefaultConfig().NetworkId,
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
	defaultSyncMode = cn.GetDefaultConfig().SyncMode
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
	OverwriteGenesisFlag = cli.BoolFlag{
		Name:  "overwrite-genesis",
		Usage: "Overwrites genesis block with the given new genesis block for testing purpose",
	}
	StartBlockNumberFlag = cli.Uint64Flag{
		Name:  "start-block-num",
		Usage: "Starts the node from the given block number. Starting from 0 is not supported.",
	}
	// Transaction pool settings
	TxPoolNoLocalsFlag = cli.BoolFlag{
		Name:  "txpool.nolocals",
		Usage: "Disables price exemptions for locally submitted transactions",
	}
	TxPoolAllowLocalAnchorTxFlag = cli.BoolFlag{
		Name:  "txpool.allow-local-anchortx",
		Usage: "Allow locally submitted anchoring transactions",
	}
	TxPoolDenyRemoteTxFlag = cli.BoolFlag{
		Name:  "txpool.deny.remotetx",
		Usage: "Deny remote transaction receiving from other peers. Use only for emergency cases",
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
		Value: cn.GetDefaultConfig().TxPool.PriceLimit,
	}
	TxPoolPriceBumpFlag = cli.Uint64Flag{
		Name:  "txpool.pricebump",
		Usage: "Price bump percentage to replace an already existing transaction",
		Value: cn.GetDefaultConfig().TxPool.PriceBump,
	}
	TxPoolExecSlotsAccountFlag = cli.Uint64Flag{
		Name:  "txpool.exec-slots.account",
		Usage: "Number of executable transaction slots guaranteed per account",
		Value: cn.GetDefaultConfig().TxPool.ExecSlotsAccount,
	}
	TxPoolExecSlotsAllFlag = cli.Uint64Flag{
		Name:  "txpool.exec-slots.all",
		Usage: "Maximum number of executable transaction slots for all accounts",
		Value: cn.GetDefaultConfig().TxPool.ExecSlotsAll,
	}
	TxPoolNonExecSlotsAccountFlag = cli.Uint64Flag{
		Name:  "txpool.nonexec-slots.account",
		Usage: "Maximum number of non-executable transaction slots permitted per account",
		Value: cn.GetDefaultConfig().TxPool.NonExecSlotsAccount,
	}
	TxPoolNonExecSlotsAllFlag = cli.Uint64Flag{
		Name:  "txpool.nonexec-slots.all",
		Usage: "Maximum number of non-executable transaction slots for all accounts",
		Value: cn.GetDefaultConfig().TxPool.NonExecSlotsAll,
	}
	TxPoolKeepLocalsFlag = cli.BoolFlag{
		Name:  "txpool.keeplocals",
		Usage: "Disables removing timed-out local transactions",
	}
	TxPoolLifetimeFlag = cli.DurationFlag{
		Name:  "txpool.lifetime",
		Usage: "Maximum amount of time non-executable transaction are queued",
		Value: cn.GetDefaultConfig().TxPool.Lifetime,
	}
	// PN specific txpool settings
	TxPoolSpamThrottlerDisableFlag = cli.BoolFlag{
		Name:  "txpool.spamthrottler.disable",
		Usage: "Disable txpool spam throttler prototype",
	}

	// KES
	KESNodeTypeServiceFlag = cli.BoolFlag{
		Name:  "kes.nodetype.service",
		Usage: "Run as a KES Service Node (Disable fetcher, downloader, and worker)",
	}
	SingleDBFlag = cli.BoolFlag{
		Name:  "db.single",
		Usage: "Create a single persistent storage. MiscDB, headerDB and etc are stored in one DB.",
	}
	NumStateTrieShardsFlag = cli.UintFlag{
		Name:  "db.num-statetrie-shards",
		Usage: "Number of internal shards of state trie DB shards. Should be power of 2",
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
	DynamoDBTableNameFlag = cli.StringFlag{
		Name:  "db.dynamo.tablename",
		Usage: "Specifies DynamoDB table name. This is mandatory to use dynamoDB. (Set dbtype to use DynamoDBS3)",
	}
	DynamoDBRegionFlag = cli.StringFlag{
		Name:  "db.dynamo.region",
		Usage: "AWS region where the DynamoDB will be created.",
		Value: database.GetDefaultDynamoDBConfig().Region,
	}
	DynamoDBIsProvisionedFlag = cli.BoolFlag{
		Name:  "db.dynamo.is-provisioned",
		Usage: "Set DynamoDB billing mode to provision. The default billing mode is on-demand.",
	}
	DynamoDBReadCapacityFlag = cli.Int64Flag{
		Name:  "db.dynamo.read-capacity",
		Usage: "Read capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied.",
		Value: database.GetDefaultDynamoDBConfig().ReadCapacityUnits,
	}
	DynamoDBWriteCapacityFlag = cli.Int64Flag{
		Name:  "db.dynamo.write-capacity",
		Usage: "Write capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied",
		Value: database.GetDefaultDynamoDBConfig().WriteCapacityUnits,
	}
	DynamoDBReadOnlyFlag = cli.BoolFlag{
		Name:  "db.dynamo.read-only",
		Usage: "Disables write to DynamoDB. Only read is possible.",
	}
	NoParallelDBWriteFlag = cli.BoolFlag{
		Name:  "db.no-parallel-write",
		Usage: "Disables parallel writes of block data to persistent database",
	}
	DBNoPerformanceMetricsFlag = cli.BoolFlag{
		Name:  "db.no-perf-metrics",
		Usage: "Disables performance metrics of database's read and write operations",
	}
	SnapshotFlag = cli.BoolFlag{
		Name:  "snapshot",
		Usage: "Enables snapshot-database mode",
	}
	SnapshotCacheSizeFlag = cli.IntFlag{
		Name:  "snapshot.cache-size",
		Usage: "Size of in-memory cache of the state snapshot cache (in MiB)",
		Value: 512,
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
	TriesInMemoryFlag = cli.Uint64Flag{
		Name:  "state.tries-in-memory",
		Usage: "The number of recent state tries residing in the memory",
		Value: blockchain.DefaultTriesInMemory,
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
	TrieNodeCacheTypeFlag = cli.StringFlag{
		Name: "statedb.cache.type",
		Usage: "Set trie node cache type ('LocalCache', 'RemoteCache', " +
			"'HybridCache') (default = 'LocalCache')",
		Value: string(statedb.CacheTypeLocal),
	}
	NumFetcherPrefetchWorkerFlag = cli.IntFlag{
		Name:  "statedb.cache.num-fetcher-prefetch-worker",
		Usage: "Number of workers used to prefetch block when fetcher fetches block",
		Value: 32,
	}
	UseSnapshotForPrefetchFlag = cli.BoolFlag{
		Name:  "statedb.cache.use-snapshot-for-prefetch",
		Usage: "Use state snapshot functionality while prefetching",
	}
	TrieNodeCacheRedisEndpointsFlag = cli.StringSliceFlag{
		Name:  "statedb.cache.redis.endpoints",
		Usage: "Set endpoints of redis trie node cache. More than one endpoints can be set",
	}
	TrieNodeCacheRedisClusterFlag = cli.BoolFlag{
		Name:  "statedb.cache.redis.cluster",
		Usage: "Enables cluster-enabled mode of redis trie node cache",
	}
	TrieNodeCacheRedisPublishBlockFlag = cli.BoolFlag{
		Name:  "statedb.cache.redis.publish",
		Usage: "Publishes every committed block to redis trie node cache",
	}
	TrieNodeCacheRedisSubscribeBlockFlag = cli.BoolFlag{
		Name:  "statedb.cache.redis.subscribe",
		Usage: "Subscribes blocks from redis trie node cache",
	}
	TrieNodeCacheLimitFlag = cli.IntFlag{
		Name:  "state.trie-cache-limit",
		Usage: "Memory allowance (MiB) to use for caching trie nodes in memory. -1 is for auto-scaling",
		Value: -1,
	}
	TrieNodeCacheSavePeriodFlag = cli.DurationFlag{
		Name:  "state.trie-cache-save-period",
		Usage: "Period of saving in memory trie cache to file if fastcache is used, 0 means disabled",
		Value: 0,
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
	VMTraceInternalTxFlag = cli.BoolFlag{
		Name:  "vm.internaltx",
		Usage: "Collect internal transaction data while processing a block",
	}

	// Logging and debug settings
	MetricsEnabledFlag = cli.BoolFlag{
		Name:  metricutils.MetricsEnabledFlag,
		Usage: "Enable metrics collection and reporting",
	}
	PrometheusExporterFlag = cli.BoolFlag{
		Name:  metricutils.PrometheusExporterFlag,
		Usage: "Enable prometheus exporter",
	}
	PrometheusExporterPortFlag = cli.IntFlag{
		Name:  metricutils.PrometheusExporterPortFlag,
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
	RPCGlobalGasCap = cli.Uint64Flag{
		Name:  "rpc.gascap",
		Usage: "Sets a cap on gas that can be used in klay_call/estimateGas",
	}
	RPCGlobalEthTxFeeCapFlag = cli.Float64Flag{
		Name:  "rpc.ethtxfeecap",
		Usage: "Sets a cap on transaction fee (in klay) that can be sent via the eth namespace RPC APIs (0 = no cap)",
	}
	RPCConcurrencyLimit = cli.IntFlag{
		Name:  "rpc.concurrencylimit",
		Usage: "Sets a limit of concurrent connection number of HTTP-RPC server",
		Value: rpc.ConcurrencyLimit,
	}
	RPCNonEthCompatibleFlag = cli.BoolFlag{
		Name:  "rpc.eth.noncompatible",
		Usage: "Disables the eth namespace API return formatting for compatibility",
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
	WSMaxSubscriptionPerConn = cli.IntFlag{
		Name:  "wsmaxsubscriptionperconn",
		Usage: "Allowed maximum subscription number per a websocket connection",
		Value: int(rpc.MaxSubscriptionPerWSConn),
	}
	WSReadDeadLine = cli.Int64Flag{
		Name:  "wsreaddeadline",
		Usage: "Set the read deadline on the underlying network connection in seconds. 0 means read will not timeout",
		Value: rpc.WebsocketReadDeadline,
	}
	WSWriteDeadLine = cli.Int64Flag{
		Name:  "wswritedeadline",
		Usage: "Set the Write deadline on the underlying network connection in seconds. 0 means write will not timeout",
		Value: rpc.WebsocketWriteDeadline,
	}
	WSMaxConnections = cli.IntFlag{
		Name:  "wsmaxconnections",
		Usage: "Allowed maximum websocket connection number",
		Value: 3000,
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
	IPCDisabledFlag = cli.BoolFlag{
		Name:  "ipcdisable",
		Usage: "Disable the IPC-RPC server",
	}
	IPCPathFlag = DirectoryFlag{
		Name:  "ipcpath",
		Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
	}
	ExecFlag = cli.StringFlag{
		Name:  "exec",
		Usage: "Execute JavaScript statement",
	}
	PreloadJSFlag = cli.StringFlag{
		Name:  "preload",
		Usage: "Comma separated list of JavaScript files to preload into the console",
	}
	APIFilterGetLogsDeadlineFlag = cli.DurationFlag{
		Name:  "api.filter.getLogs.deadline",
		Usage: "Execution deadline for log collecting filter APIs",
		Value: filters.GetLogsDeadline,
	}
	APIFilterGetLogsMaxItemsFlag = cli.IntFlag{
		Name:  "api.filter.getLogs.maxitems",
		Usage: "Maximum allowed number of return items for log collecting filter API",
		Value: filters.GetLogsMaxItems,
	}
	RPCReadTimeout = cli.IntFlag{
		Name:  "rpcreadtimeout",
		Usage: "HTTP-RPC server read timeout (seconds)",
		Value: int(rpc.DefaultHTTPTimeouts.ReadTimeout / time.Second),
	}
	RPCWriteTimeoutFlag = cli.IntFlag{
		Name:  "rpcwritetimeout",
		Usage: "HTTP-RPC server write timeout (seconds)",
		Value: int(rpc.DefaultHTTPTimeouts.WriteTimeout / time.Second),
	}
	RPCIdleTimeoutFlag = cli.IntFlag{
		Name:  "rpcidletimeout",
		Usage: "HTTP-RPC server idle timeout (seconds)",
		Value: int(rpc.DefaultHTTPTimeouts.IdleTimeout / time.Second),
	}
	RPCExecutionTimeoutFlag = cli.IntFlag{
		Name:  "rpcexecutiontimeout",
		Usage: "HTTP-RPC server execution timeout (seconds)",
		Value: int(rpc.DefaultHTTPTimeouts.ExecutionTimeout / time.Second),
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
		Usage: "The period to make and send a chain transaction to the parent chain",
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
	MaxRequestContentLengthFlag = cli.IntFlag{
		Name:  "maxRequestContentLength",
		Usage: "Max request content length in byte for http, websocket and gRPC",
		Value: common.MaxRequestContentLength,
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
		Value: 5,
	}
	ServiceChainParentOperatorTxGasLimitFlag = cli.Uint64Flag{
		Name:  "sc.parentoperator.gaslimit",
		Usage: "Set the default value of gas limit for transactions made by bridge parent operator",
		Value: 10000000,
	}
	ServiceChainChildOperatorTxGasLimitFlag = cli.Uint64Flag{
		Name:  "sc.childoperator.gaslimit",
		Usage: "Set the default value of gas limit for transactions made by bridge child operator",
		Value: 10000000,
	}
	ServiceChainNewAccountFlag = cli.BoolFlag{
		Name:  "scnewaccount",
		Usage: "Enable account creation for the service chain (default: false). If set true, generated account can't be synced with the parent chain.",
	}
	ServiceChainConsensusFlag = cli.StringFlag{
		Name:  "scconsensus",
		Usage: "Set the service chain consensus (\"istanbul\", \"clique\")",
		Value: "istanbul",
	}
	ServiceChainAnchoringFlag = cli.BoolFlag{
		Name:  "anchoring",
		Usage: "Enable anchoring for service chain",
	}

	// KAS
	KASServiceChainAnchorFlag = cli.BoolFlag{
		Name:  "kas.sc.anchor",
		Usage: "Enable KAS anchoring for service chain",
	}
	KASServiceChainAnchorPeriodFlag = cli.Uint64Flag{
		Name:  "kas.sc.anchor.period",
		Usage: "The period to anchor service chain blocks to KAS",
		Value: 1,
	}
	KASServiceChainAnchorUrlFlag = cli.StringFlag{
		Name:  "kas.sc.anchor.url",
		Usage: "The url for KAS anchor",
	}
	KASServiceChainAnchorOperatorFlag = cli.StringFlag{
		Name:  "kas.sc.anchor.operator",
		Usage: "The operator address for KAS anchor",
	}
	KASServiceChainAnchorRequestTimeoutFlag = cli.DurationFlag{
		Name:  "kas.sc.anchor.request.timeout",
		Usage: "The reuqest timeout for KAS Anchoring API call",
		Value: 500 * time.Millisecond,
	}
	KASServiceChainXChainIdFlag = cli.StringFlag{
		Name:  "kas.x-chain-id",
		Usage: "The x-chain-id for KAS",
	}
	KASServiceChainAccessKeyFlag = cli.StringFlag{
		Name:  "kas.accesskey",
		Usage: "The access key id for KAS",
	}
	KASServiceChainSecretKeyFlag = cli.StringFlag{
		Name:  "kas.secretkey",
		Usage: "The secret key for KAS",
	}

	// ChainDataFetcher
	EnableChainDataFetcherFlag = cli.BoolFlag{
		Name:  "chaindatafetcher",
		Usage: "Enable the ChainDataFetcher Service",
	}
	ChainDataFetcherMode = cli.StringFlag{
		Name:  "chaindatafetcher.mode",
		Usage: "The mode of chaindatafetcher (\"kas\", \"kafka\")",
		Value: "kas",
	}
	ChainDataFetcherNoDefault = cli.BoolFlag{
		Name:  "chaindatafetcher.no.default",
		Usage: "Turn off the starting of the chaindatafetcher",
	}
	ChainDataFetcherNumHandlers = cli.IntFlag{
		Name:  "chaindatafetcher.num.handlers",
		Usage: "Number of chaindata handlers",
		Value: chaindatafetcher.DefaultNumHandlers,
	}
	ChainDataFetcherJobChannelSize = cli.IntFlag{
		Name:  "chaindatafetcher.job.channel.size",
		Usage: "Job channel size",
		Value: chaindatafetcher.DefaultJobChannelSize,
	}
	ChainDataFetcherChainEventSizeFlag = cli.IntFlag{
		Name:  "chaindatafetcher.block.channel.size",
		Usage: "Block received channel size",
		Value: chaindatafetcher.DefaultJobChannelSize,
	}
	ChainDataFetcherKASDBHostFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.db.host",
		Usage: "KAS specific DB host in chaindatafetcher",
	}
	ChainDataFetcherKASDBPortFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.db.port",
		Usage: "KAS specific DB port in chaindatafetcher",
		Value: chaindatafetcher.DefaultDBPort,
	}
	ChainDataFetcherKASDBNameFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.db.name",
		Usage: "KAS specific DB name in chaindatafetcher",
	}
	ChainDataFetcherKASDBUserFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.db.user",
		Usage: "KAS specific DB user in chaindatafetcher",
	}
	ChainDataFetcherKASDBPasswordFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.db.password",
		Usage: "KAS specific DB password in chaindatafetcher",
	}
	ChainDataFetcherKASCacheUse = cli.BoolFlag{
		Name:  "chaindatafetcher.kas.cache.use",
		Usage: "Enable KAS cache invalidation",
	}
	ChainDataFetcherKASCacheURLFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.cache.url",
		Usage: "KAS specific cache invalidate API endpoint in chaindatafetcher",
	}
	ChainDataFetcherKASXChainIdFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.xchainid",
		Usage: "KAS specific header x-chain-id in chaindatafetcher",
	}
	ChainDataFetcherKASBasicAuthParamFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kas.basic.auth.param",
		Usage: "KAS specific header basic authorization parameter in chaindatafetcher",
	}
	ChainDataFetcherKafkaBrokersFlag = cli.StringSliceFlag{
		Name:  "chaindatafetcher.kafka.brokers",
		Usage: "Kafka broker URL list",
	}
	ChainDataFetcherKafkaTopicEnvironmentFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kafka.topic.environment",
		Usage: "Kafka topic environment prefix",
		Value: kafka.DefaultTopicEnvironmentName,
	}
	ChainDataFetcherKafkaTopicResourceFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kafka.topic.resource",
		Usage: "Kafka topic resource name",
		Value: kafka.DefaultTopicResourceName,
	}
	ChainDataFetcherKafkaReplicasFlag = cli.Int64Flag{
		Name:  "chaindatafetcher.kafka.replicas",
		Usage: "Kafka partition replication factor",
		Value: kafka.DefaultReplicas,
	}
	ChainDataFetcherKafkaPartitionsFlag = cli.IntFlag{
		Name:  "chaindatafetcher.kafka.partitions",
		Usage: "The number of partitions in a topic",
		Value: kafka.DefaultPartitions,
	}
	ChainDataFetcherKafkaMaxMessageBytesFlag = cli.Int64Flag{
		Name:  "chaindatafetcher.kafka.max.message.bytes",
		Usage: "The max size of a message produced by Kafka producer ",
		Value: kafka.DefaultMaxMessageBytes,
	}
	ChainDataFetcherKafkaSegmentSizeBytesFlag = cli.IntFlag{
		Name:  "chaindatafetcher.kafka.segment.size",
		Usage: "The kafka data segment size (in byte)",
		Value: kafka.DefaultSegmentSizeBytes,
	}
	ChainDataFetcherKafkaRequiredAcksFlag = cli.IntFlag{
		Name:  "chaindatafetcher.kafka.required.acks",
		Usage: "The level of acknowledgement reliability needed from Kafka broker (0: NoResponse, 1: WaitForLocal, -1: WaitForAll)",
		Value: kafka.DefaultRequiredAcks,
	}
	ChainDataFetcherKafkaMessageVersionFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kafka.msg.version",
		Usage: "The version of Kafka message",
		Value: kafka.DefaultKafkaMessageVersion,
	}
	ChainDataFetcherKafkaProducerIdFlag = cli.StringFlag{
		Name:  "chaindatafetcher.kafka.producer.id",
		Usage: "The identifier of kafka message producer",
		Value: kafka.GetDefaultProducerId(),
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
	AutoRestartFlag = cli.BoolFlag{
		Name:  "autorestart.enable",
		Usage: "Node can restart itself when there is a problem in making consensus",
	}
	RestartTimeOutFlag = cli.DurationFlag{
		Name:  "autorestart.timeout",
		Usage: "The elapsed time to wait auto restart (minutes)",
		Value: 15 * time.Minute,
	}
	DaemonPathFlag = cli.StringFlag{
		Name:  "autorestart.daemon.path",
		Usage: "Path of node daemon. Used to give signal to kill",
		Value: "~/klaytn/bin/kcnd",
	}

	// db migration vars
	DstDbTypeFlag = cli.StringFlag{
		Name:  "dst.dbtype",
		Usage: `Blockchain storage database type ("LevelDB", "BadgerDB", "DynamoDBS3")`,
		Value: "LevelDB",
	}
	DstDataDirFlag = DirectoryFlag{
		Name:  "dst.datadir",
		Usage: "Data directory for the databases and keystore. This value is only used in local DB.",
	}
	DstSingleDBFlag = cli.BoolFlag{
		Name:  "db.dst.single",
		Usage: "Create a single persistent storage. MiscDB, headerDB and etc are stored in one DB.",
	}
	DstLevelDBCacheSizeFlag = cli.IntFlag{
		Name:  "db.dst.leveldb.cache-size",
		Usage: "Size of in-memory cache in LevelDB (MiB)",
		Value: 768,
	}
	DstLevelDBCompressionTypeFlag = cli.IntFlag{
		Name:  "db.dst.leveldb.compression",
		Usage: "Determines the compression method for LevelDB. 0=AllNoCompression, 1=ReceiptOnlySnappyCompression, 2=StateTrieOnlyNoCompression, 3=AllSnappyCompression",
		Value: 0,
	}
	DstNumStateTrieShardsFlag = cli.UintFlag{
		Name:  "db.dst.num-statetrie-shards",
		Usage: "Number of internal shards of state trie DB shards. Should be power of 2",
		Value: 4,
	}
	DstDynamoDBTableNameFlag = cli.StringFlag{
		Name:  "db.dst.dynamo.tablename",
		Usage: "Specifies DynamoDB table name. This is mandatory to use dynamoDB. (Set dbtype to use DynamoDBS3). If dstDB is singleDB, tableName should be in form of 'PREFIX-TABLENAME'.(e.g. 'klaytn-misc', 'klaytn-statetrie')",
	}
	DstDynamoDBRegionFlag = cli.StringFlag{
		Name:  "db.dst.dynamo.region",
		Usage: "AWS region where the DynamoDB will be created.",
		Value: database.GetDefaultDynamoDBConfig().Region,
	}
	DstDynamoDBIsProvisionedFlag = cli.BoolFlag{
		Name:  "db.dst.dynamo.is-provisioned",
		Usage: "Set DynamoDB billing mode to provision. The default billing mode is on-demand.",
	}
	DstDynamoDBReadCapacityFlag = cli.Int64Flag{
		Name:  "db.dst.dynamo.read-capacity",
		Usage: "Read capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied.",
		Value: database.GetDefaultDynamoDBConfig().ReadCapacityUnits,
	}
	DstDynamoDBWriteCapacityFlag = cli.Int64Flag{
		Name:  "db.dst.dynamo.write-capacity",
		Usage: "Write capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied",
		Value: database.GetDefaultDynamoDBConfig().WriteCapacityUnits,
	}

	// Config
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
	BlockGenerationIntervalFlag = cli.Int64Flag{
		Name: "block-generation-interval",
		Usage: "(experimental option) Set the block generation interval in seconds. " +
			"It should be equal or larger than 1. This flag is only applicable to CN.",
		Value: params.DefaultBlockGenerationInterval,
	}
	BlockGenerationTimeLimitFlag = cli.DurationFlag{
		Name: "block-generation-time-limit",
		Usage: "(experimental option) Set the vm execution time limit during block generation. " +
			"Less than half of the block generation interval is recommended for this value. " +
			"This flag is only applicable to CN",
		Value: params.DefaultBlockGenerationTimeLimit,
	}
	OpcodeComputationCostLimitFlag = cli.Uint64Flag{
		Name: "opcode-computation-cost-limit",
		Usage: "(experimental option) Set the computation cost limit for a tx. " +
			"Should set the same value within the network",
		Value: params.DefaultOpcodeComputationCostLimit,
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
		cfg.WSOrigins = splitAndTrim(ctx.GlobalString(WSAllowedOriginsFlag.Name))
	}
	if ctx.GlobalIsSet(WSApiFlag.Name) {
		cfg.WSModules = splitAndTrim(ctx.GlobalString(WSApiFlag.Name))
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

// CheckExclusive verifies that only a single instance of the provided flags was
// set by the user. Each flag might optionally be followed by a string type to
// specialize it further.
func CheckExclusive(ctx *cli.Context, args ...interface{}) {
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

// SetKlayConfig applies klay-related command line flags to the config.
func SetKlayConfig(ctx *cli.Context, stack *node.Node, cfg *cn.Config) {
	// TODO-Klaytn-Bootnode: better have to check conflicts about network flags when we add Klaytn's `mainnet` parameter
	// checkExclusive(ctx, DeveloperFlag, TestnetFlag, RinkebyFlag)
	raiseFDLimit()

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

	if ctx.GlobalBool(KESNodeTypeServiceFlag.Name) {
		cfg.FetcherDisable = true
		cfg.DownloaderDisable = true
		cfg.WorkerDisable = true
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

func MakeGenesis(ctx *cli.Context) *blockchain.Genesis {
	var genesis *blockchain.Genesis
	switch {
	case ctx.GlobalBool(CypressFlag.Name):
		genesis = blockchain.DefaultGenesisBlock()
	case ctx.GlobalBool(BaobabFlag.Name):
		genesis = blockchain.DefaultBaobabGenesisBlock()
	}
	return genesis
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

// RegisterChainDataFetcherService adds a ChainDataFetcher to the stack
func RegisterChainDataFetcherService(stack *node.Node, cfg *chaindatafetcher.ChainDataFetcherConfig) {
	if cfg.EnabledChainDataFetcher {
		err := stack.RegisterSubService(func(ctx *node.ServiceContext) (node.Service, error) {
			chainDataFetcher, err := chaindatafetcher.NewChainDataFetcher(ctx, cfg)
			return chainDataFetcher, err
		})
		if err != nil {
			log.Fatalf("Failed to register the service: %v", err)
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
	var preloads []string

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
		if NodeTypeFlag.Value == "scn" || NodeTypeFlag.Value == "spn" || NodeTypeFlag.Value == "sen" {
			logger.Info("A Service Chain default network ID is set", "networkid", params.ServiceChainDefaultNetworkId)
			return params.ServiceChainDefaultNetworkId, true
		}
		logger.Info("Cypress network ID is set", "networkid", params.CypressNetworkId)
		return params.CypressNetworkId, false
	}
}
