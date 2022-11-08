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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka"
	"github.com/klaytn/klaytn/datasync/dbsyncer"
	"github.com/klaytn/klaytn/log"
	metricutils "github.com/klaytn/klaytn/metrics/utils"
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
	// app.Authors = nil
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
	ConfFlag = cli.StringFlag{
		Name: "conf",
	}
	NtpDisableFlag = cli.BoolFlag{
		Name:   "ntp.disable",
		Usage:  "Disable checking if the local time is synchronized with ntp server. If this flag is not set, the local time is checked with the time of the server specified by ntp.server.",
		EnvVar: "KLAYTN_NTP_DISABLE",
	}
	NtpServerFlag = cli.StringFlag{
		Name:   "ntp.server",
		Usage:  "Remote ntp server:port to get the time",
		Value:  "pool.ntp.org:123",
		EnvVar: "KLAYTN_NTP_SERVER",
	}
	NetworkTypeFlag = cli.StringFlag{
		Name:   "networktype",
		Usage:  "Klaytn network type (main-net (mn), service chain-net (scn))",
		Value:  "mn",
		EnvVar: "KLAYTN_NETWORKTYPE",
	}
	DbTypeFlag = cli.StringFlag{
		Name:   "dbtype",
		Usage:  `Blockchain storage database type ("LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3")`,
		Value:  "LevelDB",
		EnvVar: "KLAYTN_DBTYPE",
	}
	SrvTypeFlag = cli.StringFlag{
		Name:   "srvtype",
		Usage:  `json rpc server type ("http", "fasthttp")`,
		Value:  "fasthttp",
		EnvVar: "KLAYTN_SRVTYPE",
	}
	DataDirFlag = DirectoryFlag{
		Name:   "datadir",
		Value:  DirectoryString{node.DefaultDataDir()},
		Usage:  "Data directory for the databases and keystore. This value is only used in local DB.",
		EnvVar: "KLAYTN_DATADIR",
	}
	KeyStoreDirFlag = DirectoryFlag{
		Name:   "keystore",
		Usage:  "Directory for the keystore (default = inside the datadir)",
		EnvVar: "KLAYTN_KEYSTORE",
	}
	// TODO-Klaytn-Bootnode: redefine networkid
	NetworkIdFlag = cli.Uint64Flag{
		Name:   "networkid",
		Usage:  "Network identifier (integer, 8217=Cypress (Mainnet) , 1000=Aspen, 1001=Baobab)",
		Value:  cn.GetDefaultConfig().NetworkId,
		EnvVar: "KLAYTN_NETWORKID",
	}
	IdentityFlag = cli.StringFlag{
		Name:   "identity",
		Usage:  "Custom node name",
		EnvVar: "KLAYTN_IDENTITY",
	}
	DocRootFlag = DirectoryFlag{
		Name:   "docroot",
		Usage:  "Document Root for HTTPClient file scheme",
		Value:  DirectoryString{homeDir()},
		EnvVar: "KLAYTN_DOCROOT",
	}
	defaultSyncMode = cn.GetDefaultConfig().SyncMode
	SyncModeFlag    = TextMarshalerFlag{
		Name:   "syncmode",
		Usage:  `Blockchain sync mode ("full" or "snap")`,
		Value:  &defaultSyncMode,
		EnvVar: "KLAYTN_SYNCMODE",
	}
	GCModeFlag = cli.StringFlag{
		Name:   "gcmode",
		Usage:  `Blockchain garbage collection mode ("full", "archive")`,
		Value:  "full",
		EnvVar: "KLAYTN_GCMODE",
	}
	LightKDFFlag = cli.BoolFlag{
		Name:   "lightkdf",
		Usage:  "Reduce key-derivation RAM & CPU usage at some expense of KDF strength",
		EnvVar: "KLAYTN_LIGHTKDF",
	}
	OverwriteGenesisFlag = cli.BoolFlag{
		Name:   "overwrite-genesis",
		Usage:  "Overwrites genesis block with the given new genesis block for testing purpose",
		EnvVar: "KLAYTN_OVERWRITE_GENESIS",
	}
	StartBlockNumberFlag = cli.Uint64Flag{
		Name:   "start-block-num",
		Usage:  "Starts the node from the given block number. Starting from 0 is not supported.",
		EnvVar: "KLAYTN_START_BLOCK_NUM",
	}
	// Transaction pool settings
	TxPoolNoLocalsFlag = cli.BoolFlag{
		Name:   "txpool.nolocals",
		Usage:  "Disables price exemptions for locally submitted transactions",
		EnvVar: "KLAYTN_TXPOOL_NOLOCALS",
	}
	TxPoolAllowLocalAnchorTxFlag = cli.BoolFlag{
		Name:   "txpool.allow-local-anchortx",
		Usage:  "Allow locally submitted anchoring transactions",
		EnvVar: "KLAYTN_TXPOOL_ALLOW_LOCAL_ANCHORTX",
	}
	TxPoolDenyRemoteTxFlag = cli.BoolFlag{
		Name:   "txpool.deny.remotetx",
		Usage:  "Deny remote transaction receiving from other peers. Use only for emergency cases",
		EnvVar: "KLAYTN_TXPOOL_DENY_REMOTETX",
	}
	TxPoolJournalFlag = cli.StringFlag{
		Name:   "txpool.journal",
		Usage:  "Disk journal for local transaction to survive node restarts",
		Value:  blockchain.DefaultTxPoolConfig.Journal,
		EnvVar: "KLAYTN_TXPOOL_JOURNAL",
	}
	TxPoolJournalIntervalFlag = cli.DurationFlag{
		Name:   "txpool.journal-interval",
		Usage:  "Time interval to regenerate the local transaction journal",
		Value:  blockchain.DefaultTxPoolConfig.JournalInterval,
		EnvVar: "KLAYTN_TXPOOL_JOURNAL_INTERVAL",
	}
	TxPoolPriceLimitFlag = cli.Uint64Flag{
		Name:   "txpool.pricelimit",
		Usage:  "Minimum gas price limit to enforce for acceptance into the pool",
		Value:  cn.GetDefaultConfig().TxPool.PriceLimit,
		EnvVar: "KLAYTN_TXPOOL_PRICELIMIT",
	}
	TxPoolPriceBumpFlag = cli.Uint64Flag{
		Name:   "txpool.pricebump",
		Usage:  "Price bump percentage to replace an already existing transaction",
		Value:  cn.GetDefaultConfig().TxPool.PriceBump,
		EnvVar: "KLAYTN_TXPOOL_PRICEBUMP",
	}
	TxPoolExecSlotsAccountFlag = cli.Uint64Flag{
		Name:   "txpool.exec-slots.account",
		Usage:  "Number of executable transaction slots guaranteed per account",
		Value:  cn.GetDefaultConfig().TxPool.ExecSlotsAccount,
		EnvVar: "KLAYTN_TXPOOL_EXEC_SLOTS_ACCOUNT",
	}
	TxPoolExecSlotsAllFlag = cli.Uint64Flag{
		Name:   "txpool.exec-slots.all",
		Usage:  "Maximum number of executable transaction slots for all accounts",
		Value:  cn.GetDefaultConfig().TxPool.ExecSlotsAll,
		EnvVar: "KLAYTN_TXPOOL_EXEC_SLOTS_ALL",
	}
	TxPoolNonExecSlotsAccountFlag = cli.Uint64Flag{
		Name:   "txpool.nonexec-slots.account",
		Usage:  "Maximum number of non-executable transaction slots permitted per account",
		Value:  cn.GetDefaultConfig().TxPool.NonExecSlotsAccount,
		EnvVar: "KLAYTN_TXPOOL_NONEXEC_SLOTS_ACCOUNT",
	}
	TxPoolNonExecSlotsAllFlag = cli.Uint64Flag{
		Name:   "txpool.nonexec-slots.all",
		Usage:  "Maximum number of non-executable transaction slots for all accounts",
		Value:  cn.GetDefaultConfig().TxPool.NonExecSlotsAll,
		EnvVar: "KLAYTN_TXPOOL_NONEXEC_SLOTS_ALL",
	}
	TxPoolKeepLocalsFlag = cli.BoolFlag{
		Name:   "txpool.keeplocals",
		Usage:  "Disables removing timed-out local transactions",
		EnvVar: "KLAYTN_TXPOOL_KEEPLOCALS",
	}
	TxPoolLifetimeFlag = cli.DurationFlag{
		Name:   "txpool.lifetime",
		Usage:  "Maximum amount of time non-executable transaction are queued",
		Value:  cn.GetDefaultConfig().TxPool.Lifetime,
		EnvVar: "KLAYTN_TXPOOL_LIFETIME",
	}
	// PN specific txpool settings
	TxPoolSpamThrottlerDisableFlag = cli.BoolFlag{
		Name:   "txpool.spamthrottler.disable",
		Usage:  "Disable txpool spam throttler prototype",
		EnvVar: "KLAYTN_TXPOOL_SPAMTHROTTLER_DISABLE",
	}

	// KES
	KESNodeTypeServiceFlag = cli.BoolFlag{
		Name:   "kes.nodetype.service",
		Usage:  "Run as a KES Service Node (Disable fetcher, downloader, and worker)",
		EnvVar: "KLAYTN_KES_NODETYPE_SERVICE",
	}
	SingleDBFlag = cli.BoolFlag{
		Name:   "db.single",
		Usage:  "Create a single persistent storage. MiscDB, headerDB and etc are stored in one DB.",
		EnvVar: "KLAYTN_DB_SINGLE",
	}
	NumStateTrieShardsFlag = cli.UintFlag{
		Name:   "db.num-statetrie-shards",
		Usage:  "Number of internal shards of state trie DB shards. Should be power of 2",
		Value:  4,
		EnvVar: "KLAYTN_DB_NUM_STATETRIE_SHARDS",
	}
	LevelDBCacheSizeFlag = cli.IntFlag{
		Name:   "db.leveldb.cache-size",
		Usage:  "Size of in-memory cache in LevelDB (MiB)",
		Value:  768,
		EnvVar: "KLAYTN_DB_LEVELDB_CACHE_SIZE",
	}
	// TODO-Klaytn-Database LevelDBCompressionTypeFlag should be removed before main-net release.
	LevelDBCompressionTypeFlag = cli.IntFlag{
		Name:   "db.leveldb.compression",
		Usage:  "Determines the compression method for LevelDB. 0=AllNoCompression, 1=ReceiptOnlySnappyCompression, 2=StateTrieOnlyNoCompression, 3=AllSnappyCompression",
		Value:  0,
		EnvVar: "KLAYTN_DB_LEVELDB_COMPRESSION",
	}
	LevelDBNoBufferPoolFlag = cli.BoolFlag{
		Name:   "db.leveldb.no-buffer-pool",
		Usage:  "Disables using buffer pool for LevelDB's block allocation",
		EnvVar: "KLAYTN_DB_LEVELDB_NO_BUFFER_POOL",
	}
	DynamoDBTableNameFlag = cli.StringFlag{
		Name:   "db.dynamo.tablename",
		Usage:  "Specifies DynamoDB table name. This is mandatory to use dynamoDB. (Set dbtype to use DynamoDBS3)",
		EnvVar: "KLAYTN_DB_DYNAMO_TABLENAME",
	}
	DynamoDBRegionFlag = cli.StringFlag{
		Name:   "db.dynamo.region",
		Usage:  "AWS region where the DynamoDB will be created.",
		Value:  database.GetDefaultDynamoDBConfig().Region,
		EnvVar: "KLAYTN_DB_DYNAMO_REGION",
	}
	DynamoDBIsProvisionedFlag = cli.BoolFlag{
		Name:   "db.dynamo.is-provisioned",
		Usage:  "Set DynamoDB billing mode to provision. The default billing mode is on-demand.",
		EnvVar: "KLAYTN_DB_DYNAMO_IS_PROVISIONED",
	}
	DynamoDBReadCapacityFlag = cli.Int64Flag{
		Name:   "db.dynamo.read-capacity",
		Usage:  "Read capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied.",
		Value:  database.GetDefaultDynamoDBConfig().ReadCapacityUnits,
		EnvVar: "KLAYTN_DB_DYNAMO_READ_CAPACITY",
	}
	DynamoDBWriteCapacityFlag = cli.Int64Flag{
		Name:   "db.dynamo.write-capacity",
		Usage:  "Write capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied",
		Value:  database.GetDefaultDynamoDBConfig().WriteCapacityUnits,
		EnvVar: "KLAYTN_DB_DYNAMO_WRITE_CAPACITY",
	}
	DynamoDBReadOnlyFlag = cli.BoolFlag{
		Name:   "db.dynamo.read-only",
		Usage:  "Disables write to DynamoDB. Only read is possible.",
		EnvVar: "KLAYTN_DB_DYNAMO_READ_ONLY",
	}
	NoParallelDBWriteFlag = cli.BoolFlag{
		Name:   "db.no-parallel-write",
		Usage:  "Disables parallel writes of block data to persistent database",
		EnvVar: "KLAYTN_DB_NO_PARALLEL_WRITE",
	}
	DBNoPerformanceMetricsFlag = cli.BoolFlag{
		Name:   "db.no-perf-metrics",
		Usage:  "Disables performance metrics of database's read and write operations",
		EnvVar: "KLAYTN_DB_NO_PERF_METRICS",
	}
	SnapshotFlag = cli.BoolFlag{
		Name:   "snapshot",
		Usage:  "Enables snapshot-database mode",
		EnvVar: "KLAYTN_SNAPSHOT",
	}
	SnapshotCacheSizeFlag = cli.IntFlag{
		Name:   "snapshot.cache-size",
		Usage:  "Size of in-memory cache of the state snapshot cache (in MiB)",
		Value:  512,
		EnvVar: "KLAYTN_SNAPSHOT_CACHE_SIZE",
	}
	SnapshotAsyncGen = cli.BoolTFlag{
		Name:   "snapshot.async-gen",
		Usage:  "Enables snapshot data generation in background",
		EnvVar: "KLAYTN_SNAPSHOT_BACKGROUND_GENERATION",
	}
	TrieMemoryCacheSizeFlag = cli.IntFlag{
		Name:   "state.cache-size",
		Usage:  "Size of in-memory cache of the global state (in MiB) to flush matured singleton trie nodes to disk",
		Value:  512,
		EnvVar: "KLAYTN_STATE_CACHE_SIZE",
	}
	TrieBlockIntervalFlag = cli.UintFlag{
		Name:   "state.block-interval",
		Usage:  "An interval in terms of block number to commit the global state to disk",
		Value:  blockchain.DefaultBlockInterval,
		EnvVar: "KLAYTN_STATE_BLOCK_INTERVAL",
	}
	TriesInMemoryFlag = cli.Uint64Flag{
		Name:   "state.tries-in-memory",
		Usage:  "The number of recent state tries residing in the memory",
		Value:  blockchain.DefaultTriesInMemory,
		EnvVar: "KLAYTN_STATE_TRIES_IN_MEMORY",
	}
	CacheTypeFlag = cli.IntFlag{
		Name:   "cache.type",
		Usage:  "Cache Type: 0=LRUCache, 1=LRUShardCache, 2=FIFOCache",
		Value:  int(common.DefaultCacheType),
		EnvVar: "KLAYTN_CACHE_TYPE",
	}
	CacheScaleFlag = cli.IntFlag{
		Name:   "cache.scale",
		Usage:  "Scale of cache (cache size = preset size * scale of cache(%))",
		EnvVar: "KLAYTN_CACHE_SCALE",
	}
	CacheUsageLevelFlag = cli.StringFlag{
		Name:   "cache.level",
		Usage:  "Set the cache usage level ('saving', 'normal', 'extreme')",
		EnvVar: "KLAYTN_CACHE_LEVEL",
	}
	MemorySizeFlag = cli.IntFlag{
		Name:   "cache.memory",
		Usage:  "Set the physical RAM size (GB, Default: 16GB)",
		EnvVar: "KLAYTN_CACHE_MEMORY",
	}
	TrieNodeCacheTypeFlag = cli.StringFlag{
		Name: "statedb.cache.type",
		Usage: "Set trie node cache type ('LocalCache', 'RemoteCache', " +
			"'HybridCache') (default = 'LocalCache')",
		Value:  string(statedb.CacheTypeLocal),
		EnvVar: "KLAYTN_STATEDB_CACHE_TYPE",
	}
	NumFetcherPrefetchWorkerFlag = cli.IntFlag{
		Name:   "statedb.cache.num-fetcher-prefetch-worker",
		Usage:  "Number of workers used to prefetch block when fetcher fetches block",
		Value:  32,
		EnvVar: "KLAYTN_STATEDB_CACHE_NUM_FETCHER_PREFETCH_WORKER",
	}
	UseSnapshotForPrefetchFlag = cli.BoolFlag{
		Name:   "statedb.cache.use-snapshot-for-prefetch",
		Usage:  "Use state snapshot functionality while prefetching",
		EnvVar: "KLAYTN_STATEDB_CACHE_USE_SNAPSHOT_FOR_PREFETCH",
	}
	TrieNodeCacheRedisEndpointsFlag = cli.StringSliceFlag{
		Name:   "statedb.cache.redis.endpoints",
		Usage:  "Set endpoints of redis trie node cache. More than one endpoints can be set",
		EnvVar: "KLAYTN_STATEDB_CACHE_REDIS_ENDPOINTS",
	}
	TrieNodeCacheRedisClusterFlag = cli.BoolFlag{
		Name:   "statedb.cache.redis.cluster",
		Usage:  "Enables cluster-enabled mode of redis trie node cache",
		EnvVar: "KLAYTN_STATEDB_CACHE_REDIS_CLUSTER",
	}
	TrieNodeCacheRedisPublishBlockFlag = cli.BoolFlag{
		Name:   "statedb.cache.redis.publish",
		Usage:  "Publishes every committed block to redis trie node cache",
		EnvVar: "KLAYTN_STATEDB_CACHE_REDIS_PUBLISH",
	}
	TrieNodeCacheRedisSubscribeBlockFlag = cli.BoolFlag{
		Name:   "statedb.cache.redis.subscribe",
		Usage:  "Subscribes blocks from redis trie node cache",
		EnvVar: "KLAYTN_STATEDB_CACHE_REDIS_SUBSCRIBE",
	}
	TrieNodeCacheLimitFlag = cli.IntFlag{
		Name:   "state.trie-cache-limit",
		Usage:  "Memory allowance (MiB) to use for caching trie nodes in memory. -1 is for auto-scaling",
		Value:  -1,
		EnvVar: "KLAYTN_STATE_TRIE_CACHE_LIMIT",
	}
	TrieNodeCacheSavePeriodFlag = cli.DurationFlag{
		Name:   "state.trie-cache-save-period",
		Usage:  "Period of saving in memory trie cache to file if fastcache is used, 0 means disabled",
		Value:  0,
		EnvVar: "KLAYTN_STATE_TRIE_CACHE_SAVE_PERIOD",
	}
	SenderTxHashIndexingFlag = cli.BoolFlag{
		Name:   "sendertxhashindexing",
		Usage:  "Enables storing mapping information of senderTxHash to txHash",
		EnvVar: "KLAYTN_SENDERTXHASHINDEXING",
	}
	ChildChainIndexingFlag = cli.BoolFlag{
		Name:   "childchainindexing",
		Usage:  "Enables storing transaction hash of child chain transaction for fast access to child chain data",
		EnvVar: "KLAYTN_CHILDCHAININDEXING",
	}
	TargetGasLimitFlag = cli.Uint64Flag{
		Name:   "targetgaslimit",
		Usage:  "Target gas limit sets the artificial target gas floor for the blocks to mine",
		Value:  params.GenesisGasLimit,
		EnvVar: "KLAYTN_TARGETGASLIMIT",
	}
	ServiceChainSignerFlag = cli.StringFlag{
		Name:   "scsigner",
		Usage:  "Public address for signing blocks in the service chain (default = first account created)",
		Value:  "0",
		EnvVar: "KLAYTN_SCSIGNER",
	}
	RewardbaseFlag = cli.StringFlag{
		Name:   "rewardbase",
		Usage:  "Public address for block consensus rewards (default = first account created)",
		Value:  "0",
		EnvVar: "KLAYTN_REWARDBASE",
	}
	ExtraDataFlag = cli.StringFlag{
		Name:   "extradata",
		Usage:  "Block extra data set by the work (default = client version)",
		EnvVar: "KLAYTN_EXTRADATA",
	}

	TxResendIntervalFlag = cli.Uint64Flag{
		Name:   "txresend.interval",
		Usage:  "Set the transaction resend interval in seconds",
		Value:  uint64(cn.DefaultTxResendInterval),
		EnvVar: "KLAYTN_TXRESEND_INTERVAL",
	}
	TxResendCountFlag = cli.IntFlag{
		Name:   "txresend.max-count",
		Usage:  "Set the max count of resending transactions",
		Value:  cn.DefaultMaxResendTxCount,
		EnvVar: "KLAYTN_TXRESEND_MAX_COUNT",
	}
	// TODO-Klaytn-RemoveLater Remove this flag when we are confident with the new transaction resend logic
	TxResendUseLegacyFlag = cli.BoolFlag{
		Name:   "txresend.use-legacy",
		Usage:  "Enable the legacy transaction resend logic (For testing only)",
		EnvVar: "KLAYTN_TXRESEND_USE_LEGACY",
	}
	// Account settings
	UnlockedAccountFlag = cli.StringFlag{
		Name:   "unlock",
		Usage:  "Comma separated list of accounts to unlock",
		Value:  "",
		EnvVar: "KLAYTN_UNLOCK",
	}
	PasswordFileFlag = cli.StringFlag{
		Name:   "password",
		Usage:  "Password file to use for non-interactive password input",
		Value:  "",
		EnvVar: "KLAYTN_PASSWORD",
	}

	VMEnableDebugFlag = cli.BoolFlag{
		Name:   "vmdebug",
		Usage:  "Record information useful for VM and contract debugging",
		EnvVar: "KLAYTN_VMDEBUG",
	}
	VMLogTargetFlag = cli.IntFlag{
		Name:   "vmlog",
		Usage:  "Set the output target of vmlog precompiled contract (0: no output, 1: file, 2: stdout, 3: both)",
		Value:  0,
		EnvVar: "KLAYTN_VMLOG",
	}
	VMTraceInternalTxFlag = cli.BoolFlag{
		Name:   "vm.internaltx",
		Usage:  "Collect internal transaction data while processing a block",
		EnvVar: "KLAYTN_VM_INTERNALTX",
	}

	// Logging and debug settings
	MetricsEnabledFlag = cli.BoolFlag{
		Name:   metricutils.MetricsEnabledFlag,
		Usage:  "Enable metrics collection and reporting",
		EnvVar: "KLAYTN_METRICUTILS_METRICSENABLEDFLAG",
	}
	PrometheusExporterFlag = cli.BoolFlag{
		Name:   metricutils.PrometheusExporterFlag,
		Usage:  "Enable prometheus exporter",
		EnvVar: "KLAYTN_METRICUTILS_PROMETHEUSEXPORTERFLAG",
	}
	PrometheusExporterPortFlag = cli.IntFlag{
		Name:   metricutils.PrometheusExporterPortFlag,
		Usage:  "Prometheus exporter listening port",
		Value:  61001,
		EnvVar: "KLAYTN_METRICUTILS_PROMETHEUSEXPORTERPORTFLAG",
	}
	// RPC settings
	RPCEnabledFlag = cli.BoolFlag{
		Name:   "rpc",
		Usage:  "Enable the HTTP-RPC server",
		EnvVar: "KLAYTN_RPC",
	}
	RPCListenAddrFlag = cli.StringFlag{
		Name:   "rpcaddr",
		Usage:  "HTTP-RPC server listening interface",
		Value:  node.DefaultHTTPHost,
		EnvVar: "KLAYTN_RPCADDR",
	}
	RPCPortFlag = cli.IntFlag{
		Name:   "rpcport",
		Usage:  "HTTP-RPC server listening port",
		Value:  node.DefaultHTTPPort,
		EnvVar: "KLAYTN_RPCPORT",
	}
	RPCCORSDomainFlag = cli.StringFlag{
		Name:   "rpccorsdomain",
		Usage:  "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value:  "",
		EnvVar: "KLAYTN_RPCCORSDOMAIN",
	}
	RPCVirtualHostsFlag = cli.StringFlag{
		Name:   "rpcvhosts",
		Usage:  "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		Value:  strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
		EnvVar: "KLAYTN_RPCVHOSTS",
	}
	RPCApiFlag = cli.StringFlag{
		Name:   "rpcapi",
		Usage:  "API's offered over the HTTP-RPC interface",
		Value:  "",
		EnvVar: "KLAYTN_RPCAPI",
	}
	RPCGlobalGasCap = cli.Uint64Flag{
		Name:   "rpc.gascap",
		Usage:  "Sets a cap on gas that can be used in klay_call/estimateGas",
		EnvVar: "KLAYTN_RPC_GASCAP",
	}
	RPCGlobalEthTxFeeCapFlag = cli.Float64Flag{
		Name:   "rpc.ethtxfeecap",
		Usage:  "Sets a cap on transaction fee (in klay) that can be sent via the eth namespace RPC APIs (0 = no cap)",
		EnvVar: "KLAYTN_RPC_ETHTXFEECAP",
	}
	RPCConcurrencyLimit = cli.IntFlag{
		Name:   "rpc.concurrencylimit",
		Usage:  "Sets a limit of concurrent connection number of HTTP-RPC server",
		Value:  rpc.ConcurrencyLimit,
		EnvVar: "KLAYTN_RPC_CONCURRENCYLIMIT",
	}
	RPCNonEthCompatibleFlag = cli.BoolFlag{
		Name:   "rpc.eth.noncompatible",
		Usage:  "Disables the eth namespace API return formatting for compatibility",
		EnvVar: "KLAYTN_RPC_ETH_NONCOMPATIBLE",
	}
	WSEnabledFlag = cli.BoolFlag{
		Name:   "ws",
		Usage:  "Enable the WS-RPC server",
		EnvVar: "KLAYTN_WS",
	}
	WSListenAddrFlag = cli.StringFlag{
		Name:   "wsaddr",
		Usage:  "WS-RPC server listening interface",
		Value:  node.DefaultWSHost,
		EnvVar: "KLAYTN_WSADDR",
	}
	WSPortFlag = cli.IntFlag{
		Name:   "wsport",
		Usage:  "WS-RPC server listening port",
		Value:  node.DefaultWSPort,
		EnvVar: "KLAYTN_WSPORT",
	}
	WSApiFlag = cli.StringFlag{
		Name:   "wsapi",
		Usage:  "API's offered over the WS-RPC interface",
		Value:  "",
		EnvVar: "KLAYTN_WSAPI",
	}
	WSAllowedOriginsFlag = cli.StringFlag{
		Name:   "wsorigins",
		Usage:  "Origins from which to accept websockets requests",
		Value:  "",
		EnvVar: "KLAYTN_WSORIGINS",
	}
	WSMaxSubscriptionPerConn = cli.IntFlag{
		Name:   "wsmaxsubscriptionperconn",
		Usage:  "Allowed maximum subscription number per a websocket connection",
		Value:  int(rpc.MaxSubscriptionPerWSConn),
		EnvVar: "KLAYTN_WSMAXSUBSCRIPTIONPERCONN",
	}
	WSReadDeadLine = cli.Int64Flag{
		Name:   "wsreaddeadline",
		Usage:  "Set the read deadline on the underlying network connection in seconds. 0 means read will not timeout",
		Value:  rpc.WebsocketReadDeadline,
		EnvVar: "KLAYTN_WSREADDEADLINE",
	}
	WSWriteDeadLine = cli.Int64Flag{
		Name:   "wswritedeadline",
		Usage:  "Set the Write deadline on the underlying network connection in seconds. 0 means write will not timeout",
		Value:  rpc.WebsocketWriteDeadline,
		EnvVar: "KLAYTN_WSWRITEDEADLINE",
	}
	WSMaxConnections = cli.IntFlag{
		Name:   "wsmaxconnections",
		Usage:  "Allowed maximum websocket connection number",
		Value:  3000,
		EnvVar: "KLAYTN_WSMAXCONNECTIONS",
	}
	GRPCEnabledFlag = cli.BoolFlag{
		Name:   "grpc",
		Usage:  "Enable the gRPC server",
		EnvVar: "KLAYTN_GRPC",
	}
	GRPCListenAddrFlag = cli.StringFlag{
		Name:   "grpcaddr",
		Usage:  "gRPC server listening interface",
		Value:  node.DefaultGRPCHost,
		EnvVar: "KLAYTN_GRPCADDR",
	}
	GRPCPortFlag = cli.IntFlag{
		Name:   "grpcport",
		Usage:  "gRPC server listening port",
		Value:  node.DefaultGRPCPort,
		EnvVar: "KLAYTN_GRPCPORT",
	}
	IPCDisabledFlag = cli.BoolFlag{
		Name:   "ipcdisable",
		Usage:  "Disable the IPC-RPC server",
		EnvVar: "KLAYTN_IPCDISABLE",
	}
	IPCPathFlag = DirectoryFlag{
		Name:   "ipcpath",
		Usage:  "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
		EnvVar: "KLAYTN_IPCPATH",
	}
	ExecFlag = cli.StringFlag{
		Name:   "exec",
		Usage:  "Execute JavaScript statement",
		EnvVar: "KLAYTN_EXEC",
	}
	PreloadJSFlag = cli.StringFlag{
		Name:   "preload",
		Usage:  "Comma separated list of JavaScript files to preload into the console",
		EnvVar: "KLAYTN_PRELOAD",
	}
	APIFilterGetLogsDeadlineFlag = cli.DurationFlag{
		Name:   "api.filter.getLogs.deadline",
		Usage:  "Execution deadline for log collecting filter APIs",
		Value:  filters.GetLogsDeadline,
		EnvVar: "KLAYTN_API_FILTER_GETLOGS_DEADLINE",
	}
	APIFilterGetLogsMaxItemsFlag = cli.IntFlag{
		Name:   "api.filter.getLogs.maxitems",
		Usage:  "Maximum allowed number of return items for log collecting filter API",
		Value:  filters.GetLogsMaxItems,
		EnvVar: "KLAYTN_API_FILTER_GETLOGS_MAXITEMS",
	}
	RPCReadTimeout = cli.IntFlag{
		Name:   "rpcreadtimeout",
		Usage:  "HTTP-RPC server read timeout (seconds)",
		Value:  int(rpc.DefaultHTTPTimeouts.ReadTimeout / time.Second),
		EnvVar: "KLAYTN_RPCREADTIMEOUT",
	}
	RPCWriteTimeoutFlag = cli.IntFlag{
		Name:   "rpcwritetimeout",
		Usage:  "HTTP-RPC server write timeout (seconds)",
		Value:  int(rpc.DefaultHTTPTimeouts.WriteTimeout / time.Second),
		EnvVar: "KLAYTN_RPCWRITETIMEOUT",
	}
	RPCIdleTimeoutFlag = cli.IntFlag{
		Name:   "rpcidletimeout",
		Usage:  "HTTP-RPC server idle timeout (seconds)",
		Value:  int(rpc.DefaultHTTPTimeouts.IdleTimeout / time.Second),
		EnvVar: "KLAYTN_RPCIDLETIMEOUT",
	}
	RPCExecutionTimeoutFlag = cli.IntFlag{
		Name:   "rpcexecutiontimeout",
		Usage:  "HTTP-RPC server execution timeout (seconds)",
		Value:  int(rpc.DefaultHTTPTimeouts.ExecutionTimeout / time.Second),
		EnvVar: "KLAYTN_RPCEXECUTIONTIMEOUT",
	}

	// Network Settings
	NodeTypeFlag = cli.StringFlag{
		Name:   "nodetype",
		Usage:  "Klaytn node type (consensus node (cn), proxy node (pn), endpoint node (en))",
		Value:  "en",
		EnvVar: "KLAYTN_NODETYPE",
	}
	MaxConnectionsFlag = cli.IntFlag{
		Name:   "maxconnections",
		Usage:  "Maximum number of physical connections. All single channel peers can be maxconnections peers. All multi channel peers can be maxconnections/2 peers. (network disabled if set to 0)",
		Value:  node.DefaultMaxPhysicalConnections,
		EnvVar: "KLAYTN_MAXCONNECTIONS",
	}
	MaxPendingPeersFlag = cli.IntFlag{
		Name:   "maxpendpeers",
		Usage:  "Maximum number of pending connection attempts (defaults used if set to 0)",
		Value:  0,
		EnvVar: "KLAYTN_MAXPENDPEERS",
	}
	ListenPortFlag = cli.IntFlag{
		Name:   "port",
		Usage:  "Network listening port",
		Value:  node.DefaultP2PPort,
		EnvVar: "KLAYTN_PORT",
	}
	SubListenPortFlag = cli.IntFlag{
		Name:   "subport",
		Usage:  "Network sub listening port",
		Value:  node.DefaultP2PSubPort,
		EnvVar: "KLAYTN_SUBPORT",
	}
	MultiChannelUseFlag = cli.BoolFlag{
		Name:   "multichannel",
		Usage:  "Create a dedicated channel for block propagation",
		EnvVar: "KLAYTN_MULTICHANNEL",
	}
	BootnodesFlag = cli.StringFlag{
		Name:   "bootnodes",
		Usage:  "Comma separated kni URLs for P2P discovery bootstrap",
		Value:  "",
		EnvVar: "KLAYTN_BOOTNODES",
	}
	NodeKeyFileFlag = cli.StringFlag{
		Name:   "nodekey",
		Usage:  "P2P node key file",
		EnvVar: "KLAYTN_NODEKEY",
	}
	NodeKeyHexFlag = cli.StringFlag{
		Name:   "nodekeyhex",
		Usage:  "P2P node key as hex (for testing)",
		EnvVar: "KLAYTN_NODEKEYHEX",
	}
	NATFlag = cli.StringFlag{
		Name:   "nat",
		Usage:  "NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>)",
		Value:  "any",
		EnvVar: "KLAYTN_NAT",
	}
	NoDiscoverFlag = cli.BoolFlag{
		Name:   "nodiscover",
		Usage:  "Disables the peer discovery mechanism (manual peer addition)",
		EnvVar: "KLAYTN_NODISCOVER",
	}
	NetrestrictFlag = cli.StringFlag{
		Name:   "netrestrict",
		Usage:  "Restricts network communication to the given IP network (CIDR masks)",
		EnvVar: "KLAYTN_NETRESTRICT",
	}
	AnchoringPeriodFlag = cli.Uint64Flag{
		Name:   "chaintxperiod",
		Usage:  "The period to make and send a chain transaction to the parent chain",
		Value:  1,
		EnvVar: "KLAYTN_CHAINTXPERIOD",
	}
	SentChainTxsLimit = cli.Uint64Flag{
		Name:   "chaintxlimit",
		Usage:  "Number of service chain transactions stored for resending",
		Value:  100,
		EnvVar: "KLAYTN_CHAINTXLIMIT",
	}
	RWTimerIntervalFlag = cli.Uint64Flag{
		Name:   "rwtimerinterval",
		Usage:  "Interval of using rw timer to check if it works well",
		Value:  1000,
		EnvVar: "KLAYTN_RWTIMERINTERVAL",
	}
	RWTimerWaitTimeFlag = cli.DurationFlag{
		Name:   "rwtimerwaittime",
		Usage:  "Wait time the rw timer waits for message writing",
		Value:  15 * time.Second,
		EnvVar: "KLAYTN_RWTIMERWAITTIME",
	}
	MaxRequestContentLengthFlag = cli.IntFlag{
		Name:   "maxRequestContentLength",
		Usage:  "Max request content length in byte for http, websocket and gRPC",
		Value:  common.MaxRequestContentLength,
		EnvVar: "KLAYTN_MAXREQUESTCONTENTLENGTH",
	}

	// ATM the url is left to the user and deployment to
	JSpathFlag = cli.StringFlag{
		Name:   "jspath",
		Usage:  "JavaScript root path for `loadScript`",
		Value:  ".",
		EnvVar: "KLAYTN_JSPATH",
	}
	CypressFlag = cli.BoolFlag{
		Name:   "cypress",
		Usage:  "Pre-configured Klaytn Cypress network",
		EnvVar: "KLAYTN_CYPRESS",
	}
	// Baobab bootnodes setting
	BaobabFlag = cli.BoolFlag{
		Name:   "baobab",
		Usage:  "Pre-configured Klaytn baobab network",
		EnvVar: "KLAYTN_BAOBAB",
	}
	// Bootnode's settings
	AuthorizedNodesFlag = cli.StringFlag{
		Name:   "authorized-nodes",
		Usage:  "Comma separated kni URLs for authorized nodes list",
		Value:  "",
		EnvVar: "KLAYTN_AUTHORIZED_NODES",
	}
	// TODO-Klaytn-Bootnode the boodnode flags should be updated when it is implemented
	BNAddrFlag = cli.StringFlag{
		Name:   "bnaddr",
		Usage:  `udp address to use node discovery`,
		Value:  ":32323",
		EnvVar: "KLAYTN_BNADDR",
	}
	GenKeyFlag = cli.StringFlag{
		Name:   "genkey",
		Usage:  "generate a node private key and write to given filename",
		EnvVar: "KLAYTN_GENKEY",
	}
	WriteAddressFlag = cli.BoolFlag{
		Name:   "writeaddress",
		Usage:  `write out the node's public key which is given by "--nodekeyfile" or "--nodekeyhex"`,
		EnvVar: "KLAYTN_WRITEADDRESS",
	}
	// ServiceChain's settings
	MainBridgeFlag = cli.BoolFlag{
		Name:   "mainbridge",
		Usage:  "Enable main bridge service for service chain",
		EnvVar: "KLAYTN_MAINBRIDGE",
	}
	SubBridgeFlag = cli.BoolFlag{
		Name:   "subbridge",
		Usage:  "Enable sub bridge service for service chain",
		EnvVar: "KLAYTN_SUBBRIDGE",
	}
	MainBridgeListenPortFlag = cli.IntFlag{
		Name:   "mainbridgeport",
		Usage:  "main bridge listen port",
		Value:  50505,
		EnvVar: "KLAYTN_MAINBRIDGEPORT",
	}
	SubBridgeListenPortFlag = cli.IntFlag{
		Name:   "subbridgeport",
		Usage:  "sub bridge listen port",
		Value:  50506,
		EnvVar: "KLAYTN_SUBBRIDGEPORT",
	}
	ParentChainIDFlag = cli.IntFlag{
		Name:   "parentchainid",
		Usage:  "parent chain ID",
		Value:  8217, // Klaytn mainnet chain ID
		EnvVar: "KLAYTN_PARENTCHAINID",
	}
	VTRecoveryFlag = cli.BoolFlag{
		Name:   "vtrecovery",
		Usage:  "Enable value transfer recovery (default: false)",
		EnvVar: "KLAYTN_VTRECOVERY",
	}
	VTRecoveryIntervalFlag = cli.Uint64Flag{
		Name:   "vtrecoveryinterval",
		Usage:  "Set the value transfer recovery interval (seconds)",
		Value:  5,
		EnvVar: "KLAYTN_VTRECOVERYINTERVAL",
	}
	ServiceChainParentOperatorTxGasLimitFlag = cli.Uint64Flag{
		Name:   "sc.parentoperator.gaslimit",
		Usage:  "Set the default value of gas limit for transactions made by bridge parent operator",
		Value:  10000000,
		EnvVar: "KLAYTN_SC_PARENTOPERATOR_GASLIMIT",
	}
	ServiceChainChildOperatorTxGasLimitFlag = cli.Uint64Flag{
		Name:   "sc.childoperator.gaslimit",
		Usage:  "Set the default value of gas limit for transactions made by bridge child operator",
		Value:  10000000,
		EnvVar: "KLAYTN_SC_CHILDOPERATOR_GASLIMIT",
	}
	ServiceChainNewAccountFlag = cli.BoolFlag{
		Name:   "scnewaccount",
		Usage:  "Enable account creation for the service chain (default: false). If set true, generated account can't be synced with the parent chain.",
		EnvVar: "KLAYTN_SCNEWACCOUNT",
	}
	ServiceChainConsensusFlag = cli.StringFlag{
		Name:   "scconsensus",
		Usage:  "Set the service chain consensus (\"istanbul\", \"clique\")",
		Value:  "istanbul",
		EnvVar: "KLAYTN_SCCONSENSUS",
	}
	ServiceChainAnchoringFlag = cli.BoolFlag{
		Name:   "anchoring",
		Usage:  "Enable anchoring for service chain",
		EnvVar: "KLAYTN_ANCHORING",
	}

	// KAS
	KASServiceChainAnchorFlag = cli.BoolFlag{
		Name:   "kas.sc.anchor",
		Usage:  "Enable KAS anchoring for service chain",
		EnvVar: "KLAYTN_KAS_SC_ANCHOR",
	}
	KASServiceChainAnchorPeriodFlag = cli.Uint64Flag{
		Name:   "kas.sc.anchor.period",
		Usage:  "The period to anchor service chain blocks to KAS",
		Value:  1,
		EnvVar: "KLAYTN_KAS_SC_ANCHOR_PERIOD",
	}
	KASServiceChainAnchorUrlFlag = cli.StringFlag{
		Name:   "kas.sc.anchor.url",
		Usage:  "The url for KAS anchor",
		EnvVar: "KLAYTN_KAS_SC_ANCHOR_URL",
	}
	KASServiceChainAnchorOperatorFlag = cli.StringFlag{
		Name:   "kas.sc.anchor.operator",
		Usage:  "The operator address for KAS anchor",
		EnvVar: "KLAYTN_KAS_SC_ANCHOR_OPERATOR",
	}
	KASServiceChainAnchorRequestTimeoutFlag = cli.DurationFlag{
		Name:   "kas.sc.anchor.request.timeout",
		Usage:  "The reuqest timeout for KAS Anchoring API call",
		Value:  500 * time.Millisecond,
		EnvVar: "KLAYTN_KAS_SC_ANCHOR_REQUEST_TIMEOUT",
	}
	KASServiceChainXChainIdFlag = cli.StringFlag{
		Name:   "kas.x-chain-id",
		Usage:  "The x-chain-id for KAS",
		EnvVar: "KLAYTN_KAS_X_CHAIN_ID",
	}
	KASServiceChainAccessKeyFlag = cli.StringFlag{
		Name:   "kas.accesskey",
		Usage:  "The access key id for KAS",
		EnvVar: "KLAYTN_KAS_ACCESSKEY",
	}
	KASServiceChainSecretKeyFlag = cli.StringFlag{
		Name:   "kas.secretkey",
		Usage:  "The secret key for KAS",
		EnvVar: "KLAYTN_KAS_SECRETKEY",
	}

	// ChainDataFetcher
	EnableChainDataFetcherFlag = cli.BoolFlag{
		Name:   "chaindatafetcher",
		Usage:  "Enable the ChainDataFetcher Service",
		EnvVar: "KLAYTN_CHAINDATAFETCHER",
	}
	ChainDataFetcherMode = cli.StringFlag{
		Name:   "chaindatafetcher.mode",
		Usage:  "The mode of chaindatafetcher (\"kas\", \"kafka\")",
		Value:  "kas",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_MODE",
	}
	ChainDataFetcherNoDefault = cli.BoolFlag{
		Name:   "chaindatafetcher.no.default",
		Usage:  "Turn off the starting of the chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_NO_DEFAULT",
	}
	ChainDataFetcherNumHandlers = cli.IntFlag{
		Name:   "chaindatafetcher.num.handlers",
		Usage:  "Number of chaindata handlers",
		Value:  chaindatafetcher.DefaultNumHandlers,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_NUM_HANDLERS",
	}
	ChainDataFetcherJobChannelSize = cli.IntFlag{
		Name:   "chaindatafetcher.job.channel.size",
		Usage:  "Job channel size",
		Value:  chaindatafetcher.DefaultJobChannelSize,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_JOB_CHANNEL_SIZE",
	}
	ChainDataFetcherChainEventSizeFlag = cli.IntFlag{
		Name:   "chaindatafetcher.block.channel.size",
		Usage:  "Block received channel size",
		Value:  chaindatafetcher.DefaultJobChannelSize,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_BLOCK_CHANNEL_SIZE",
	}
	ChainDataFetcherMaxProcessingDataSize = cli.IntFlag{
		Name:   "chaindatafetcher.max.processing.data.size",
		Usage:  "Maximum size of processing data before requesting range fetching of blocks (in MB)",
		Value:  chaindatafetcher.DefaultMaxProcessingDataSize,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_MAX_PROCESSING_DATA_SIZE",
	}
	ChainDataFetcherKASDBHostFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.db.host",
		Usage:  "KAS specific DB host in chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_DB_HOST",
	}
	ChainDataFetcherKASDBPortFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.db.port",
		Usage:  "KAS specific DB port in chaindatafetcher",
		Value:  chaindatafetcher.DefaultDBPort,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_DB_PORT",
	}
	ChainDataFetcherKASDBNameFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.db.name",
		Usage:  "KAS specific DB name in chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_DB_NAME",
	}
	ChainDataFetcherKASDBUserFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.db.user",
		Usage:  "KAS specific DB user in chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_DB_USER",
	}
	ChainDataFetcherKASDBPasswordFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.db.password",
		Usage:  "KAS specific DB password in chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_DB_PASSWORD",
	}
	ChainDataFetcherKASCacheUse = cli.BoolFlag{
		Name:   "chaindatafetcher.kas.cache.use",
		Usage:  "Enable KAS cache invalidation",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_CACHE_USE",
	}
	ChainDataFetcherKASCacheURLFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.cache.url",
		Usage:  "KAS specific cache invalidate API endpoint in chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_CACHE_URL",
	}
	ChainDataFetcherKASXChainIdFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.xchainid",
		Usage:  "KAS specific header x-chain-id in chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_XCHAINID",
	}
	ChainDataFetcherKASBasicAuthParamFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kas.basic.auth.param",
		Usage:  "KAS specific header basic authorization parameter in chaindatafetcher",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAS_BASIC_AUTH_PARAM",
	}
	ChainDataFetcherKafkaBrokersFlag = cli.StringSliceFlag{
		Name:   "chaindatafetcher.kafka.brokers",
		Usage:  "Kafka broker URL list",
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_BROKERS",
	}
	ChainDataFetcherKafkaTopicEnvironmentFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kafka.topic.environment",
		Usage:  "Kafka topic environment prefix",
		Value:  kafka.DefaultTopicEnvironmentName,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_TOPIC_ENVIRONMENT",
	}
	ChainDataFetcherKafkaTopicResourceFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kafka.topic.resource",
		Usage:  "Kafka topic resource name",
		Value:  kafka.DefaultTopicResourceName,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_TOPIC_RESOURCE",
	}
	ChainDataFetcherKafkaReplicasFlag = cli.Int64Flag{
		Name:   "chaindatafetcher.kafka.replicas",
		Usage:  "Kafka partition replication factor",
		Value:  kafka.DefaultReplicas,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_REPLICAS",
	}
	ChainDataFetcherKafkaPartitionsFlag = cli.IntFlag{
		Name:   "chaindatafetcher.kafka.partitions",
		Usage:  "The number of partitions in a topic",
		Value:  kafka.DefaultPartitions,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_PARTITIONS",
	}
	ChainDataFetcherKafkaMaxMessageBytesFlag = cli.Int64Flag{
		Name:   "chaindatafetcher.kafka.max.message.bytes",
		Usage:  "The max size of a message produced by Kafka producer ",
		Value:  kafka.DefaultMaxMessageBytes,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_MAX_MESSAGE_BYTES",
	}
	ChainDataFetcherKafkaSegmentSizeBytesFlag = cli.IntFlag{
		Name:   "chaindatafetcher.kafka.segment.size",
		Usage:  "The kafka data segment size (in byte)",
		Value:  kafka.DefaultSegmentSizeBytes,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_SEGMENT_SIZE",
	}
	ChainDataFetcherKafkaRequiredAcksFlag = cli.IntFlag{
		Name:   "chaindatafetcher.kafka.required.acks",
		Usage:  "The level of acknowledgement reliability needed from Kafka broker (0: NoResponse, 1: WaitForLocal, -1: WaitForAll)",
		Value:  kafka.DefaultRequiredAcks,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_REQUIRED_ACKS",
	}
	ChainDataFetcherKafkaMessageVersionFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kafka.msg.version",
		Usage:  "The version of Kafka message",
		Value:  kafka.DefaultKafkaMessageVersion,
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_MSG_VERSION",
	}
	ChainDataFetcherKafkaProducerIdFlag = cli.StringFlag{
		Name:   "chaindatafetcher.kafka.producer.id",
		Usage:  "The identifier of kafka message producer",
		Value:  kafka.GetDefaultProducerId(),
		EnvVar: "KLAYTN_CHAINDATAFETCHER_KAFKA_PRODUCER_ID",
	}
	// DBSyncer
	EnableDBSyncerFlag = cli.BoolFlag{
		Name:   "dbsyncer",
		Usage:  "Enable the DBSyncer",
		EnvVar: "KLAYTN_DBSYNCER",
	}
	DBHostFlag = cli.StringFlag{
		Name:   "dbsyncer.db.host",
		Usage:  "db.host in dbsyncer",
		EnvVar: "KLAYTN_DBSYNCER_DB_HOST",
	}
	DBPortFlag = cli.StringFlag{
		Name:   "dbsyncer.db.port",
		Usage:  "db.port in dbsyncer",
		Value:  "3306",
		EnvVar: "KLAYTN_DBSYNCER_DB_PORT",
	}
	DBNameFlag = cli.StringFlag{
		Name:   "dbsyncer.db.name",
		Usage:  "db.name in dbsyncer",
		EnvVar: "KLAYTN_DBSYNCER_DB_NAME",
	}
	DBUserFlag = cli.StringFlag{
		Name:   "dbsyncer.db.user",
		Usage:  "db.user in dbsyncer",
		EnvVar: "KLAYTN_DBSYNCER_DB_USER",
	}
	DBPasswordFlag = cli.StringFlag{
		Name:   "dbsyncer.db.password",
		Usage:  "db.password in dbsyncer",
		EnvVar: "KLAYTN_DBSYNCER_DB_PASSWORD",
	}
	EnabledLogModeFlag = cli.BoolFlag{
		Name:   "dbsyncer.logmode",
		Usage:  "Enable the dbsyncer logmode",
		EnvVar: "KLAYTN_DBSYNCER_LOGMODE",
	}
	MaxIdleConnsFlag = cli.IntFlag{
		Name:   "dbsyncer.db.max.idle",
		Usage:  "The maximum number of connections in the idle connection pool",
		Value:  50,
		EnvVar: "KLAYTN_DBSYNCER_DB_MAX_IDLE",
	}
	MaxOpenConnsFlag = cli.IntFlag{
		Name:   "dbsyncer.db.max.open",
		Usage:  "The maximum number of open connections to the database",
		Value:  30,
		EnvVar: "KLAYTN_DBSYNCER_DB_MAX_OPEN",
	}
	ConnMaxLifeTimeFlag = cli.DurationFlag{
		Name:   "dbsyncer.db.max.lifetime",
		Usage:  "The maximum amount of time a connection may be reused (default : 1h), ex: 300ms, 2h45m, 60s, ...",
		Value:  1 * time.Hour,
		EnvVar: "KLAYTN_DBSYNCER_DB_MAX_LIFETIME",
	}
	BlockSyncChannelSizeFlag = cli.IntFlag{
		Name:   "dbsyncer.block.channel.size",
		Usage:  "Block received channel size",
		Value:  5,
		EnvVar: "KLAYTN_DBSYNCER_BLOCK_CHANNEL_SIZE",
	}
	DBSyncerModeFlag = cli.StringFlag{
		Name:   "dbsyncer.mode",
		Usage:  "The mode of dbsyncer is way which handle block/tx data to insert db (multi, single, context)",
		Value:  "multi",
		EnvVar: "KLAYTN_DBSYNCER_MODE",
	}
	GenQueryThreadFlag = cli.IntFlag{
		Name:   "dbsyncer.genquery.th",
		Usage:  "The amount of thread of generation query in multi mode",
		Value:  50,
		EnvVar: "KLAYTN_DBSYNCER_GENQUERY_TH",
	}
	InsertThreadFlag = cli.IntFlag{
		Name:   "dbsyncer.insert.th",
		Usage:  "The amount of thread of insert operation in multi mode",
		Value:  30,
		EnvVar: "KLAYTN_DBSYNCER_INSERT_TH",
	}
	BulkInsertSizeFlag = cli.IntFlag{
		Name:   "dbsyncer.bulk.size",
		Usage:  "The amount of row for bulk-insert",
		Value:  200,
		EnvVar: "KLAYTN_DBSYNCER_BULK_SIZE",
	}
	EventModeFlag = cli.StringFlag{
		Name:   "dbsyncer.event.mode",
		Usage:  "The way how to sync all block or last block (block, head)",
		Value:  "head",
		EnvVar: "KLAYTN_DBSYNCER_EVENT_MODE",
	}
	MaxBlockDiffFlag = cli.Uint64Flag{
		Name:   "dbsyncer.max.block.diff",
		Usage:  "The maximum difference between current block and event block. 0 means off",
		Value:  0,
		EnvVar: "KLAYTN_DBSYNCER_MAX_BLOCK_DIFF",
	}
	AutoRestartFlag = cli.BoolFlag{
		Name:   "autorestart.enable",
		Usage:  "Node can restart itself when there is a problem in making consensus",
		EnvVar: "KLAYTN_AUTORESTART_ENABLE",
	}
	RestartTimeOutFlag = cli.DurationFlag{
		Name:   "autorestart.timeout",
		Usage:  "The elapsed time to wait auto restart (minutes)",
		Value:  15 * time.Minute,
		EnvVar: "KLAYTN_AUTORESTART_TIMEOUT",
	}
	DaemonPathFlag = cli.StringFlag{
		Name:   "autorestart.daemon.path",
		Usage:  "Path of node daemon. Used to give signal to kill",
		Value:  "~/klaytn/bin/kcnd",
		EnvVar: "KLAYTN_AUTORESTART_DAEMON_PATH",
	}

	// db migration vars
	DstDbTypeFlag = cli.StringFlag{
		Name:   "dst.dbtype",
		Usage:  `Blockchain storage database type ("LevelDB", "BadgerDB", "DynamoDBS3")`,
		Value:  "LevelDB",
		EnvVar: "KLAYTN_DST_DBTYPE",
	}
	DstDataDirFlag = DirectoryFlag{
		Name:   "dst.datadir",
		Usage:  "Data directory for the databases and keystore. This value is only used in local DB.",
		EnvVar: "KLAYTN_DST_DATADIR",
	}
	DstSingleDBFlag = cli.BoolFlag{
		Name:   "db.dst.single",
		Usage:  "Create a single persistent storage. MiscDB, headerDB and etc are stored in one DB.",
		EnvVar: "KLAYTN_DB_DST_SINGLE",
	}
	DstLevelDBCacheSizeFlag = cli.IntFlag{
		Name:   "db.dst.leveldb.cache-size",
		Usage:  "Size of in-memory cache in LevelDB (MiB)",
		Value:  768,
		EnvVar: "KLAYTN_DB_DST_LEVELDB_CACHE_SIZE",
	}
	DstLevelDBCompressionTypeFlag = cli.IntFlag{
		Name:   "db.dst.leveldb.compression",
		Usage:  "Determines the compression method for LevelDB. 0=AllNoCompression, 1=ReceiptOnlySnappyCompression, 2=StateTrieOnlyNoCompression, 3=AllSnappyCompression",
		Value:  0,
		EnvVar: "KLAYTN_DB_DST_LEVELDB_COMPRESSION",
	}
	DstNumStateTrieShardsFlag = cli.UintFlag{
		Name:   "db.dst.num-statetrie-shards",
		Usage:  "Number of internal shards of state trie DB shards. Should be power of 2",
		Value:  4,
		EnvVar: "KLAYTN_DB_DST_NUM_STATETRIE_SHARDS",
	}
	DstDynamoDBTableNameFlag = cli.StringFlag{
		Name:   "db.dst.dynamo.tablename",
		Usage:  "Specifies DynamoDB table name. This is mandatory to use dynamoDB. (Set dbtype to use DynamoDBS3). If dstDB is singleDB, tableName should be in form of 'PREFIX-TABLENAME'.(e.g. 'klaytn-misc', 'klaytn-statetrie')",
		EnvVar: "KLAYTN_DB_DST_DYNAMO_TABLENAME",
	}
	DstDynamoDBRegionFlag = cli.StringFlag{
		Name:   "db.dst.dynamo.region",
		Usage:  "AWS region where the DynamoDB will be created.",
		Value:  database.GetDefaultDynamoDBConfig().Region,
		EnvVar: "KLAYTN_DB_DST_DYNAMO_REGION",
	}
	DstDynamoDBIsProvisionedFlag = cli.BoolFlag{
		Name:   "db.dst.dynamo.is-provisioned",
		Usage:  "Set DynamoDB billing mode to provision. The default billing mode is on-demand.",
		EnvVar: "KLAYTN_DB_DST_DYNAMO_IS_PROVISIONED",
	}
	DstDynamoDBReadCapacityFlag = cli.Int64Flag{
		Name:   "db.dst.dynamo.read-capacity",
		Usage:  "Read capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied.",
		Value:  database.GetDefaultDynamoDBConfig().ReadCapacityUnits,
		EnvVar: "KLAYTN_DB_DST_DYNAMO_READ_CAPACITY",
	}
	DstDynamoDBWriteCapacityFlag = cli.Int64Flag{
		Name:   "db.dst.dynamo.write-capacity",
		Usage:  "Write capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied",
		Value:  database.GetDefaultDynamoDBConfig().WriteCapacityUnits,
		EnvVar: "KLAYTN_DB_DST_DYNAMO_WRITE_CAPACITY",
	}

	// Config
	ConfigFileFlag = cli.StringFlag{
		Name:   "config",
		Usage:  "TOML configuration file",
		EnvVar: "KLAYTN_CONFIG",
	}
	BlockGenerationIntervalFlag = cli.Int64Flag{
		Name: "block-generation-interval",
		Usage: "(experimental option) Set the block generation interval in seconds. " +
			"It should be equal or larger than 1. This flag is only applicable to CN.",
		Value:  params.DefaultBlockGenerationInterval,
		EnvVar: "KLAYTN_BLOCK_GENERATION_INTERVAL",
	}
	BlockGenerationTimeLimitFlag = cli.DurationFlag{
		Name: "block-generation-time-limit",
		Usage: "(experimental option) Set the vm execution time limit during block generation. " +
			"Less than half of the block generation interval is recommended for this value. " +
			"This flag is only applicable to CN",
		Value:  params.DefaultBlockGenerationTimeLimit,
		EnvVar: "KLAYTN_BLOCK_GENERATION_TIME_LIMIT",
	}
	OpcodeComputationCostLimitFlag = cli.Uint64Flag{
		Name: "opcode-computation-cost-limit",
		Usage: "(experimental option) Set the computation cost limit for a tx. " +
			"Should set the same value within the network",
		Value:  params.DefaultOpcodeComputationCostLimit,
		EnvVar: "KLAYTN_OPCODE_COMPUTATION_COST_LIMIT",
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

// splitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func SplitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
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
