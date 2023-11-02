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
	"github.com/urfave/cli/v2"
)

func init() {
	cli.FlagStringer = FlagString
}

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	// app.Author = ""
	// app.Authors = nil
	// app.Email = ""
	app.Version = params.Version
	if len(gitCommit) >= 8 {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	return app
}

var (
	// General settings
	ConfFlag = &cli.StringFlag{
		Name: "conf",
	}
	NtpDisableFlag = &cli.BoolFlag{
		Name:     "ntp.disable",
		Usage:    "Disable checking if the local time is synchronized with ntp server. If this flag is not set, the local time is checked with the time of the server specified by ntp.server.",
		Value:    false,
		Aliases:  []string{"common.ntp.disable"},
		EnvVars:  []string{"KLAYTN_NTP_DISABLE"},
		Category: "KLAY",
	}
	NtpServerFlag = &cli.StringFlag{
		Name:     "ntp.server",
		Usage:    "Remote ntp server:port to get the time",
		Value:    "pool.ntp.org:123",
		Aliases:  []string{"common.ntp.server", "ns"},
		EnvVars:  []string{"KLAYTN_NTP_SERVER"},
		Category: "KLAY",
	}
	NetworkTypeFlag = &cli.StringFlag{
		Name:    "networktype",
		Usage:   "Klaytn network type (main-net (mn), service chain-net (scn))",
		Value:   "mn",
		Aliases: []string{},
		EnvVars: []string{"KLAYTN_NETWORKTYPE"},
	}
	DbTypeFlag = &cli.StringFlag{
		Name:     "dbtype",
		Usage:    `Blockchain storage database type ("LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3")`,
		Value:    "LevelDB",
		Aliases:  []string{"db.type", "migration.src.dbtype"},
		EnvVars:  []string{"KLAYTN_DBTYPE"},
		Category: "KLAY",
	}
	SrvTypeFlag = &cli.StringFlag{
		Name:     "srvtype",
		Usage:    `json rpc server type ("http", "fasthttp")`,
		Value:    "fasthttp",
		Aliases:  []string{"common.srvtype"},
		EnvVars:  []string{"KLAYTN_SRVTYPE"},
		Category: "KLAY",
	}
	DataDirFlag = &cli.PathFlag{
		Name:     "datadir",
		Value:    node.DefaultDataDir(),
		Usage:    "Data directory for the databases and keystore. This value is only used in local DB.",
		Aliases:  []string{"common.datadir", "migration.src.datadir"},
		EnvVars:  []string{"KLAYTN_DATADIR"},
		Category: "KLAY",
	}
	ChainDataDirFlag = &cli.PathFlag{
		Name:     "chaindatadir",
		Value:    "",
		Usage:    "Data directory for chaindata. If this is not specified, chaindata is stored in datadir",
		Aliases:  []string{"common.chaindatadir"},
		EnvVars:  []string{"KLAYTN_CHAINDATADIR"},
		Category: "KLAY",
	}
	KeyStoreDirFlag = &cli.PathFlag{
		Name:     "keystore",
		Usage:    "Directory for the keystore (default = inside the datadir)",
		Aliases:  []string{"common.keystore"},
		EnvVars:  []string{"KLAYTN_KEYSTORE"},
		Category: "ACCOUNT",
	}
	// TODO-Klaytn-Bootnode: redefine networkid
	NetworkIdFlag = &cli.Uint64Flag{
		Name:     "networkid",
		Usage:    "Network identifier (integer, 8217=Cypress (Mainnet) , 1000=Aspen, 1001=Baobab)",
		Value:    cn.GetDefaultConfig().NetworkId,
		Aliases:  []string{"p2p.network-id"},
		EnvVars:  []string{"KLAYTN_NETWORKID"},
		Category: "NETWORK",
	}
	IdentityFlag = &cli.StringFlag{
		Name:     "identity",
		Usage:    "Custom node name",
		Aliases:  []string{"common.identity"},
		EnvVars:  []string{"KLAYTN_IDENTITY"},
		Category: "KLAY",
	}
	DocRootFlag = &cli.PathFlag{
		Name:  "docroot",
		Usage: "Document Root for HTTPClient file scheme",
		// Value:   DirectoryString{homeDir()},
		Value:    homeDir(),
		Aliases:  []string{"common.docroot"},
		EnvVars:  []string{"KLAYTN_DOCROOT"},
		Category: "KLAY",
	}
	defaultSyncMode = cn.GetDefaultConfig().SyncMode
	SyncModeFlag    = &TextMarshalerFlag{
		Name:     "syncmode",
		Usage:    `Blockchain sync mode ("full" or "snap")`,
		Value:    &defaultSyncMode,
		Aliases:  []string{"common.syncmode"},
		EnvVars:  []string{"KLAYTN_SYNCMODE"},
		Category: "KLAY",
	}
	GCModeFlag = &cli.StringFlag{
		Name:     "gcmode",
		Usage:    `Blockchain garbage collection mode ("full", "archive")`,
		Value:    "full",
		Aliases:  []string{"common.garbage-collection-mode"},
		EnvVars:  []string{"KLAYTN_GCMODE"},
		Category: "KLAY",
	}
	LightKDFFlag = &cli.BoolFlag{
		Name:     "lightkdf",
		Usage:    "Reduce key-derivation RAM & CPU usage at some expense of KDF strength",
		Aliases:  []string{"common.light-kdf"},
		EnvVars:  []string{"KLAYTN_LIGHTKDF"},
		Category: "ACCOUNT",
	}
	OverwriteGenesisFlag = &cli.BoolFlag{
		Name:     "overwrite-genesis",
		Usage:    "Overwrites genesis block with the given new genesis block for testing purpose",
		Aliases:  []string{"common.overwrite-genesis"},
		EnvVars:  []string{"KLAYTN_OVERWRITE_GENESIS"},
		Category: "KLAY",
	}
	StartBlockNumberFlag = &cli.Uint64Flag{
		Name:     "start-block-num",
		Usage:    "Starts the node from the given block number. Starting from 0 is not supported.",
		Aliases:  []string{"common.start-block-num"},
		EnvVars:  []string{"KLAYTN_START_BLOCK_NUM"},
		Category: "KLAY",
	}
	// Transaction pool settings
	TxPoolNoLocalsFlag = &cli.BoolFlag{
		Name:     "txpool.nolocals",
		Usage:    "Disables price exemptions for locally submitted transactions",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_NOLOCALS"},
		Category: "TXPOOL",
	}
	TxPoolAllowLocalAnchorTxFlag = &cli.BoolFlag{
		Name:     "txpool.allow-local-anchortx",
		Usage:    "Allow locally submitted anchoring transactions",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_ALLOW_LOCAL_ANCHORTX"},
		Category: "TXPOOL",
	}
	TxPoolDenyRemoteTxFlag = &cli.BoolFlag{
		Name:     "txpool.deny.remotetx",
		Usage:    "Deny remote transaction receiving from other peers. Use only for emergency cases",
		Aliases:  []string{"txpool.deny-remote-tx"},
		EnvVars:  []string{"KLAYTN_TXPOOL_DENY_REMOTETX"},
		Category: "TXPOOL",
	}
	TxPoolJournalFlag = &cli.StringFlag{
		Name:     "txpool.journal",
		Usage:    "Disk journal for local transaction to survive node restarts",
		Value:    blockchain.DefaultTxPoolConfig.Journal,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_JOURNAL"},
		Category: "TXPOOL",
	}
	TxPoolJournalIntervalFlag = &cli.DurationFlag{
		Name:     "txpool.journal-interval",
		Usage:    "Time interval to regenerate the local transaction journal",
		Value:    blockchain.DefaultTxPoolConfig.JournalInterval,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_JOURNAL_INTERVAL"},
		Category: "TXPOOL",
	}
	TxPoolPriceLimitFlag = &cli.Uint64Flag{
		Name:     "txpool.pricelimit",
		Usage:    "Minimum gas price limit to enforce for acceptance into the pool",
		Value:    cn.GetDefaultConfig().TxPool.PriceLimit,
		Aliases:  []string{"txpool.price-limit"},
		EnvVars:  []string{"KLAYTN_TXPOOL_PRICELIMIT"},
		Category: "TXPOOL",
	}
	TxPoolPriceBumpFlag = &cli.Uint64Flag{
		Name:     "txpool.pricebump",
		Usage:    "Price bump percentage to replace an already existing transaction",
		Value:    cn.GetDefaultConfig().TxPool.PriceBump,
		Aliases:  []string{"txpool.price-bump"},
		EnvVars:  []string{"KLAYTN_TXPOOL_PRICEBUMP"},
		Category: "TXPOOL",
	}
	TxPoolExecSlotsAccountFlag = &cli.Uint64Flag{
		Name:     "txpool.exec-slots.account",
		Usage:    "Number of executable transaction slots guaranteed per account",
		Value:    cn.GetDefaultConfig().TxPool.ExecSlotsAccount,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_EXEC_SLOTS_ACCOUNT"},
		Category: "TXPOOL",
	}
	TxPoolExecSlotsAllFlag = &cli.Uint64Flag{
		Name:     "txpool.exec-slots.all",
		Usage:    "Maximum number of executable transaction slots for all accounts",
		Value:    cn.GetDefaultConfig().TxPool.ExecSlotsAll,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_EXEC_SLOTS_ALL"},
		Category: "TXPOOL",
	}
	TxPoolNonExecSlotsAccountFlag = &cli.Uint64Flag{
		Name:     "txpool.nonexec-slots.account",
		Usage:    "Maximum number of non-executable transaction slots permitted per account",
		Value:    cn.GetDefaultConfig().TxPool.NonExecSlotsAccount,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_NONEXEC_SLOTS_ACCOUNT"},
		Category: "TXPOOL",
	}
	TxPoolNonExecSlotsAllFlag = &cli.Uint64Flag{
		Name:     "txpool.nonexec-slots.all",
		Usage:    "Maximum number of non-executable transaction slots for all accounts",
		Value:    cn.GetDefaultConfig().TxPool.NonExecSlotsAll,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_NONEXEC_SLOTS_ALL"},
		Category: "TXPOOL",
	}
	TxPoolKeepLocalsFlag = &cli.BoolFlag{
		Name:     "txpool.keeplocals",
		Usage:    "Disables removing timed-out local transactions",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_KEEPLOCALS"},
		Category: "TXPOOL",
	}
	TxPoolLifetimeFlag = &cli.DurationFlag{
		Name:     "txpool.lifetime",
		Usage:    "Maximum amount of time non-executable transaction are queued",
		Value:    cn.GetDefaultConfig().TxPool.Lifetime,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXPOOL_LIFETIME"},
		Category: "TXPOOL",
	}
	// PN specific txpool settings
	TxPoolSpamThrottlerDisableFlag = &cli.BoolFlag{
		Name:    "txpool.spamthrottler.disable",
		Usage:   "Disable txpool spam throttler prototype",
		Aliases: []string{},
		EnvVars: []string{"KLAYTN_TXPOOL_SPAMTHROTTLER_DISABLE"},
	}

	// KES
	KESNodeTypeServiceFlag = &cli.BoolFlag{
		Name:     "kes.nodetype.service",
		Usage:    "Run as a KES Service Node (Disable fetcher, downloader, and worker)",
		Aliases:  []string{"common.kes-nodetype-service"},
		EnvVars:  []string{"KLAYTN_KES_NODETYPE_SERVICE"},
		Category: "MISC",
	}
	SingleDBFlag = &cli.BoolFlag{
		Name:     "db.single",
		Usage:    "Create a single persistent storage. MiscDB, headerDB and etc are stored in one DB.",
		Aliases:  []string{"migration.src.single"},
		EnvVars:  []string{"KLAYTN_DB_SINGLE"},
		Category: "DATABASE",
	}
	NumStateTrieShardsFlag = &cli.UintFlag{
		Name:     "db.num-statetrie-shards",
		Usage:    "Number of internal shards of state trie DB shards. Should be power of 2",
		Value:    4,
		Aliases:  []string{"migration.src.db.leveldb.num-statetrie-shards"},
		EnvVars:  []string{"KLAYTN_DB_NUM_STATETRIE_SHARDS"},
		Category: "DATABASE",
	}
	LevelDBCacheSizeFlag = &cli.IntFlag{
		Name:     "db.leveldb.cache-size",
		Usage:    "Size of in-memory cache in LevelDB (MiB)",
		Value:    768,
		Aliases:  []string{"migration.src.db.leveldb.cache-size"},
		EnvVars:  []string{"KLAYTN_DB_LEVELDB_CACHE_SIZE"},
		Category: "DATABASE",
	}
	// TODO-Klaytn-Database LevelDBCompressionTypeFlag should be removed before main-net release.
	LevelDBCompressionTypeFlag = &cli.IntFlag{
		Name:     "db.leveldb.compression",
		Usage:    "Determines the compression method for LevelDB. 0=AllNoCompression, 1=ReceiptOnlySnappyCompression, 2=StateTrieOnlyNoCompression, 3=AllSnappyCompression",
		Value:    0,
		Aliases:  []string{"migration.src.db.leveldb.compression"},
		EnvVars:  []string{"KLAYTN_DB_LEVELDB_COMPRESSION"},
		Category: "DATABASE",
	}
	LevelDBNoBufferPoolFlag = &cli.BoolFlag{
		Name:     "db.leveldb.no-buffer-pool",
		Usage:    "Disables using buffer pool for LevelDB's block allocation",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_DB_LEVELDB_NO_BUFFER_POOL"},
		Category: "DATABASE",
	}
	RocksDBSecondaryFlag = &cli.BoolFlag{
		Name:     "db.rocksdb.secondary",
		Usage:    "Enable rocksdb secondary mode (read-only and catch-up with primary node dynamically)",
		Aliases:  []string{"migration.src.db.rocksdb.secondary"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_SECONDARY"},
		Category: "DATABASE",
	}
	RocksDBCacheSizeFlag = &cli.Uint64Flag{
		Name:     "db.rocksdb.cache-size",
		Usage:    "Size of in-memory cache in RocksDB (MiB)",
		Value:    768,
		Aliases:  []string{"migration.src.db.rocksdb.cache-size"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_CACHE_SIZE"},
		Category: "DATABASE",
	}
	RocksDBDumpMallocStatFlag = &cli.BoolFlag{
		Name:     "db.rocksdb.dump-malloc-stat",
		Usage:    "Enable to print memory stat together with rocksdb.stat. Works with Jemalloc only.",
		Aliases:  []string{"migration.src.db.rocksdb.dump-malloc-stat"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_DUMP_MALLOC_STAT"},
		Category: "DATABASE",
	}
	RocksDBCompressionTypeFlag = &cli.StringFlag{
		Name:     "db.rocksdb.compression-type",
		Usage:    "RocksDB block compression type. Supported values are 'no', 'snappy', 'zlib', 'bz', 'lz4', 'lz4hc', 'xpress', 'zstd'",
		Value:    database.GetDefaultRocksDBConfig().CompressionType,
		Aliases:  []string{"migration.src.db.rocksdb.compression-type"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_COMPRESSION_TYPE"},
		Category: "DATABASE",
	}
	RocksDBBottommostCompressionTypeFlag = &cli.StringFlag{
		Name:     "db.rocksdb.bottommost-compression-type",
		Usage:    "RocksDB bottommost block compression type. Supported values are 'no', 'snappy', 'zlib', 'bz2', 'lz4', 'lz4hc', 'xpress', 'zstd'",
		Value:    database.GetDefaultRocksDBConfig().BottommostCompressionType,
		Aliases:  []string{"migration.src.db.rocksdb.bottommost-compression-type"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_BOTTOMMOST_COMPRESSION_TYPE"},
		Category: "DATABASE",
	}
	RocksDBFilterPolicyFlag = &cli.StringFlag{
		Name:     "db.rocksdb.filter-policy",
		Usage:    "RocksDB filter policy. Supported values are 'no', 'bloom', 'ribbon'",
		Value:    database.GetDefaultRocksDBConfig().FilterPolicy,
		Aliases:  []string{"migration.src.db.rocksdb.filter-policy"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_FILTER_POLICY"},
		Category: "DATABASE",
	}
	RocksDBDisableMetricsFlag = &cli.BoolFlag{
		Name:     "db.rocksdb.disable-metrics",
		Usage:    "Disable RocksDB metrics",
		Aliases:  []string{"migration.src.db.rocksdb.disable-metrics"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_DISABLE_METRICS"},
		Category: "DATABASE",
	}
	RocksDBMaxOpenFilesFlag = &cli.IntFlag{
		Name:     "db.rocksdb.max-open-files",
		Usage:    "Set RocksDB max open files. (the value should be greater than 16)",
		Value:    database.GetDefaultRocksDBConfig().MaxOpenFiles,
		Aliases:  []string{"migration.src.db.rocksdb.max-open-files"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_MAX_OPEN_FILES"},
		Category: "DATABASE",
	}
	RocksDBCacheIndexAndFilterFlag = &cli.BoolFlag{
		Name:     "db.rocksdb.cache-index-and-filter",
		Usage:    "Use block cache for index and filter blocks.",
		Aliases:  []string{"migration.src.db.rocksdb.cache-index-and-filter"},
		EnvVars:  []string{"KLAYTN_DB_ROCKSDB_CACHE_INDEX_AND_FILTER"},
		Category: "DATABASE",
	}
	DynamoDBTableNameFlag = &cli.StringFlag{
		Name:     "db.dynamo.tablename",
		Usage:    "Specifies DynamoDB table name. This is mandatory to use dynamoDB. (Set dbtype to use DynamoDBS3)",
		Aliases:  []string{"migration.src.db.dynamo.table-name", "db.dynamo.table-name"},
		EnvVars:  []string{"KLAYTN_DB_DYNAMO_TABLENAME"},
		Category: "DATABASE",
	}
	DynamoDBRegionFlag = &cli.StringFlag{
		Name:     "db.dynamo.region",
		Usage:    "AWS region where the DynamoDB will be created.",
		Value:    database.GetDefaultDynamoDBConfig().Region,
		Aliases:  []string{"migration.src.db.dynamo.region"},
		EnvVars:  []string{"KLAYTN_DB_DYNAMO_REGION"},
		Category: "DATABASE",
	}
	DynamoDBIsProvisionedFlag = &cli.BoolFlag{
		Name:     "db.dynamo.is-provisioned",
		Usage:    "Set DynamoDB billing mode to provision. The default billing mode is on-demand.",
		Aliases:  []string{"migration.src.db.dynamo.is-provisioned"},
		EnvVars:  []string{"KLAYTN_DB_DYNAMO_IS_PROVISIONED"},
		Category: "DATABASE",
	}
	DynamoDBReadCapacityFlag = &cli.Int64Flag{
		Name:     "db.dynamo.read-capacity",
		Usage:    "Read capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied.",
		Value:    database.GetDefaultDynamoDBConfig().ReadCapacityUnits,
		Aliases:  []string{"migration.src.db.dynamo.read-capacity"},
		EnvVars:  []string{"KLAYTN_DB_DYNAMO_READ_CAPACITY"},
		Category: "DATABASE",
	}
	DynamoDBWriteCapacityFlag = &cli.Int64Flag{
		Name:     "db.dynamo.write-capacity",
		Usage:    "Write capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied",
		Value:    database.GetDefaultDynamoDBConfig().WriteCapacityUnits,
		Aliases:  []string{"migration.src.db.dynamo.write-capacity"},
		EnvVars:  []string{"KLAYTN_DB_DYNAMO_WRITE_CAPACITY"},
		Category: "DATABASE",
	}
	DynamoDBReadOnlyFlag = &cli.BoolFlag{
		Name:     "db.dynamo.read-only",
		Usage:    "Disables write to DynamoDB. Only read is possible.",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_DB_DYNAMO_READ_ONLY"},
		Category: "DATABASE",
	}
	NoParallelDBWriteFlag = &cli.BoolFlag{
		Name:     "db.no-parallel-write",
		Usage:    "Disables parallel writes of block data to persistent database",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_DB_NO_PARALLEL_WRITE"},
		Category: "DATABASE",
	}
	DBNoPerformanceMetricsFlag = &cli.BoolFlag{
		Name:     "db.no-perf-metrics",
		Usage:    "Disables performance metrics of database's read and write operations",
		Value:    false,
		Aliases:  []string{"migration.no-perf-metrics"},
		EnvVars:  []string{"KLAYTN_DB_NO_PERF_METRICS"},
		Category: "DATABASE",
	}
	SnapshotFlag = &cli.BoolFlag{
		Name:     "snapshot",
		Usage:    "Enables snapshot-database mode",
		Aliases:  []string{"snapshot-database.enable"},
		EnvVars:  []string{"KLAYTN_SNAPSHOT"},
		Category: "MISC",
	}
	SnapshotCacheSizeFlag = &cli.IntFlag{
		Name:     "snapshot.cache-size",
		Usage:    "Size of in-memory cache of the state snapshot cache (in MiB)",
		Value:    512,
		Aliases:  []string{"snapshot-database.cache-size"},
		EnvVars:  []string{"KLAYTN_SNAPSHOT_CACHE_SIZE"},
		Category: "MISC",
	}
	SnapshotAsyncGen = &cli.BoolFlag{
		Name:     "snapshot.async-gen",
		Usage:    "Enables snapshot data generation in background",
		Value:    true,
		Aliases:  []string{"snapshot-database.async-gen"},
		EnvVars:  []string{"KLAYTN_SNAPSHOT_BACKGROUND_GENERATION"},
		Category: "MISC",
	}
	TrieMemoryCacheSizeFlag = &cli.IntFlag{
		Name:     "state.cache-size",
		Usage:    "Size of in-memory cache of the global state (in MiB) to flush matured singleton trie nodes to disk",
		Value:    512,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATE_CACHE_SIZE"},
		Category: "STATE",
	}
	TrieBlockIntervalFlag = &cli.UintFlag{
		Name:     "state.block-interval",
		Usage:    "An interval in terms of block number to commit the global state to disk",
		Value:    blockchain.DefaultBlockInterval,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATE_BLOCK_INTERVAL"},
		Category: "STATE",
	}
	TriesInMemoryFlag = &cli.Uint64Flag{
		Name:     "state.tries-in-memory",
		Usage:    "The number of recent state tries residing in the memory",
		Value:    blockchain.DefaultTriesInMemory,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATE_TRIES_IN_MEMORY"},
		Category: "STATE",
	}
	LivePruningFlag = &cli.BoolFlag{
		Name:     "state.live-pruning",
		Usage:    "Enable trie live pruning",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATE_LIVE_PRUNING"},
		Category: "STATE",
	}
	LivePruningRetentionFlag = &cli.Uint64Flag{
		Name:     "state.live-pruning-retention",
		Usage:    "Number of blocks from the latest block that are not to be pruned",
		Value:    blockchain.DefaultLivePruningRetention,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATE_LIVE_PRUNING_RETENTION"},
		Category: "STATE",
	}
	CacheTypeFlag = &cli.IntFlag{
		Name:     "cache.type",
		Usage:    "Cache Type: 0=LRUCache, 1=LRUShardCache, 2=FIFOCache",
		Value:    int(common.DefaultCacheType),
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_CACHE_TYPE"},
		Category: "CACHE",
	}
	CacheScaleFlag = &cli.IntFlag{
		Name:     "cache.scale",
		Usage:    "Scale of cache (cache size = preset size * scale of cache(%))",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_CACHE_SCALE"},
		Category: "CACHE",
	}
	CacheUsageLevelFlag = &cli.StringFlag{
		Name:     "cache.level",
		Usage:    "Set the cache usage level ('saving', 'normal', 'extreme')",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_CACHE_LEVEL"},
		Category: "CACHE",
	}
	MemorySizeFlag = &cli.IntFlag{
		Name:     "cache.memory",
		Usage:    "Set the physical RAM size (GB, Default: 16GB)",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_CACHE_MEMORY"},
		Category: "CACHE",
	}
	TrieNodeCacheTypeFlag = &cli.StringFlag{
		Name: "statedb.cache.type",
		Usage: "Set trie node cache type ('LocalCache', 'RemoteCache', " +
			"'HybridCache') (default = 'LocalCache')",
		Value:    string(statedb.CacheTypeLocal),
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATEDB_CACHE_TYPE"},
		Category: "CACHE",
	}
	NumFetcherPrefetchWorkerFlag = &cli.IntFlag{
		Name:     "statedb.cache.num-fetcher-prefetch-worker",
		Usage:    "Number of workers used to prefetch block when fetcher fetches block",
		Value:    32,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATEDB_CACHE_NUM_FETCHER_PREFETCH_WORKER"},
		Category: "CACHE",
	}
	UseSnapshotForPrefetchFlag = &cli.BoolFlag{
		Name:     "statedb.cache.use-snapshot-for-prefetch",
		Usage:    "Use state snapshot functionality while prefetching",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATEDB_CACHE_USE_SNAPSHOT_FOR_PREFETCH"},
		Category: "CACHE",
	}
	TrieNodeCacheRedisEndpointsFlag = &cli.StringSliceFlag{
		Name:     "statedb.cache.redis.endpoints",
		Usage:    "Set endpoints of redis trie node cache. More than one endpoints can be set",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATEDB_CACHE_REDIS_ENDPOINTS"},
		Category: "CACHE",
	}
	TrieNodeCacheRedisClusterFlag = &cli.BoolFlag{
		Name:     "statedb.cache.redis.cluster",
		Usage:    "Enables cluster-enabled mode of redis trie node cache",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATEDB_CACHE_REDIS_CLUSTER"},
		Category: "CACHE",
	}
	TrieNodeCacheRedisPublishBlockFlag = &cli.BoolFlag{
		Name:     "statedb.cache.redis.publish",
		Usage:    "Publishes every committed block to redis trie node cache",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATEDB_CACHE_REDIS_PUBLISH"},
		Category: "CACHE",
	}
	TrieNodeCacheRedisSubscribeBlockFlag = &cli.BoolFlag{
		Name:     "statedb.cache.redis.subscribe",
		Usage:    "Subscribes blocks from redis trie node cache",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATEDB_CACHE_REDIS_SUBSCRIBE"},
		Category: "CACHE",
	}
	TrieNodeCacheLimitFlag = &cli.IntFlag{
		Name:     "state.trie-cache-limit",
		Usage:    "Memory allowance (MiB) to use for caching trie nodes in memory. -1 is for auto-scaling",
		Value:    -1,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATE_TRIE_CACHE_LIMIT"},
		Category: "CACHE",
	}
	TrieNodeCacheSavePeriodFlag = &cli.DurationFlag{
		Name:     "state.trie-cache-save-period",
		Usage:    "Period of saving in memory trie cache to file if fastcache is used, 0 means disabled",
		Value:    0,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_STATE_TRIE_CACHE_SAVE_PERIOD"},
		Category: "CACHE",
	}

	SenderTxHashIndexingFlag = &cli.BoolFlag{
		Name:     "sendertxhashindexing",
		Usage:    "Enables storing mapping information of senderTxHash to txHash",
		Aliases:  []string{"common.sender-tx-hash-indexing"},
		EnvVars:  []string{"KLAYTN_SENDERTXHASHINDEXING"},
		Category: "DATABASE",
	}
	ChildChainIndexingFlag = &cli.BoolFlag{
		Name:     "childchainindexing",
		Usage:    "Enables storing transaction hash of child chain transaction for fast access to child chain data",
		Aliases:  []string{"common.child-chain-indexing"},
		EnvVars:  []string{"KLAYTN_CHILDCHAININDEXING"},
		Category: "SERVICECHAIN",
	}
	TargetGasLimitFlag = &cli.Uint64Flag{
		Name:     "targetgaslimit",
		Usage:    "Target gas limit sets the artificial target gas floor for the blocks to mine",
		Value:    params.GenesisGasLimit,
		Aliases:  []string{"common.target-gaslimit"},
		EnvVars:  []string{"KLAYTN_TARGETGASLIMIT"},
		Category: "NETWORK",
	}
	ServiceChainSignerFlag = &cli.StringFlag{
		Name:     "scsigner",
		Usage:    "Public address for signing blocks in the service chain (default = first account created)",
		Value:    "0",
		Aliases:  []string{"common.scsigner"},
		EnvVars:  []string{"KLAYTN_SCSIGNER"},
		Category: "CONSENSUS",
	}
	RewardbaseFlag = &cli.StringFlag{
		Name:     "rewardbase",
		Usage:    "Public address for block consensus rewards (default = first account created)",
		Value:    "0",
		Aliases:  []string{"common.rewardbase"},
		EnvVars:  []string{"KLAYTN_REWARDBASE"},
		Category: "CONSENSUS",
	}
	ExtraDataFlag = &cli.StringFlag{
		Name:     "extradata",
		Usage:    "Block extra data set by the work (default = client version)",
		Aliases:  []string{"common.block-extra-data"},
		EnvVars:  []string{"KLAYTN_EXTRADATA"},
		Category: "KLAY",
	}

	TxResendIntervalFlag = &cli.Uint64Flag{
		Name:     "txresend.interval",
		Usage:    "Set the transaction resend interval in seconds",
		Value:    uint64(cn.DefaultTxResendInterval),
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXRESEND_INTERVAL"},
		Category: "TXPOOL",
	}
	TxResendCountFlag = &cli.IntFlag{
		Name:     "txresend.max-count",
		Usage:    "Set the max count of resending transactions",
		Value:    cn.DefaultMaxResendTxCount,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXRESEND_MAX_COUNT"},
		Category: "TXPOOL",
	}
	// TODO-Klaytn-RemoveLater Remove this flag when we are confident with the new transaction resend logic
	TxResendUseLegacyFlag = &cli.BoolFlag{
		Name:     "txresend.use-legacy",
		Usage:    "Enable the legacy transaction resend logic (For testing only)",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_TXRESEND_USE_LEGACY"},
		Category: "TXPOOL",
	}
	// Account settings
	UnlockedAccountFlag = &cli.StringFlag{
		Name:     "unlock",
		Usage:    "Comma separated list of accounts to unlock",
		Value:    "",
		Aliases:  []string{"account-update.unlock"},
		EnvVars:  []string{"KLAYTN_UNLOCK"},
		Category: "ACCOUNT",
	}
	PasswordFileFlag = &cli.StringFlag{
		Name:     "password",
		Usage:    "Password file to use for non-interactive password input",
		Value:    "",
		Aliases:  []string{"account-update.password"},
		EnvVars:  []string{"KLAYTN_PASSWORD"},
		Category: "ACCOUNT",
	}

	VMEnableDebugFlag = &cli.BoolFlag{
		Name:     "vmdebug",
		Usage:    "Record information useful for VM and contract debugging",
		Aliases:  []string{"vm.debug"},
		EnvVars:  []string{"KLAYTN_VMDEBUG"},
		Category: "VIRTUAL MACHINE",
	}
	VMLogTargetFlag = &cli.IntFlag{
		Name:     "vmlog",
		Usage:    "Set the output target of vmlog precompiled contract (0: no output, 1: file, 2: stdout, 3: both)",
		Value:    0,
		Aliases:  []string{"vm.log"},
		EnvVars:  []string{"KLAYTN_VMLOG"},
		Category: "VIRTUAL MACHINE",
	}
	VMTraceInternalTxFlag = &cli.BoolFlag{
		Name:     "vm.internaltx",
		Usage:    "Collect internal transaction data while processing a block",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_VM_INTERNALTX"},
		Category: "VIRTUAL MACHINE",
	}

	// Logging and debug settings
	MetricsEnabledFlag = &cli.BoolFlag{
		Name:     metricutils.MetricsEnabledFlag,
		Usage:    "Enable metrics collection and reporting",
		Aliases:  []string{"metrics-collection-reporting.enable"},
		EnvVars:  []string{"KLAYTN_METRICUTILS_METRICSENABLEDFLAG"},
		Category: "METRIC",
	}
	PrometheusExporterFlag = &cli.BoolFlag{
		Name:     metricutils.PrometheusExporterFlag,
		Usage:    "Enable prometheus exporter",
		Aliases:  []string{"metrics-collection-reporting.prometheus"},
		EnvVars:  []string{"KLAYTN_METRICUTILS_PROMETHEUSEXPORTERFLAG"},
		Category: "METRIC",
	}
	PrometheusExporterPortFlag = &cli.IntFlag{
		Name:     metricutils.PrometheusExporterPortFlag,
		Usage:    "Prometheus exporter listening port",
		Value:    61001,
		Aliases:  []string{"metrics-collection-reporting.prometheus-port"},
		EnvVars:  []string{"KLAYTN_METRICUTILS_PROMETHEUSEXPORTERPORTFLAG"},
		Category: "METRIC",
	}

	// RPC settings
	RPCEnabledFlag = &cli.BoolFlag{
		Name:     "rpc",
		Usage:    "Enable the HTTP-RPC server",
		Aliases:  []string{"http-rpc.enable"},
		EnvVars:  []string{"KLAYTN_RPC"},
		Category: "API AND CONSOLE",
	}
	RPCListenAddrFlag = &cli.StringFlag{
		Name:     "rpcaddr",
		Usage:    "HTTP-RPC server listening interface",
		Value:    node.DefaultHTTPHost,
		Aliases:  []string{"http-rpc.addr"},
		EnvVars:  []string{"KLAYTN_RPCADDR"},
		Category: "API AND CONSOLE",
	}
	RPCPortFlag = &cli.IntFlag{
		Name:     "rpcport",
		Usage:    "HTTP-RPC server listening port",
		Value:    node.DefaultHTTPPort,
		Aliases:  []string{"http-rpc.port"},
		EnvVars:  []string{"KLAYTN_RPCPORT"},
		Category: "API AND CONSOLE",
	}
	RPCCORSDomainFlag = &cli.StringFlag{
		Name:     "rpccorsdomain",
		Usage:    "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value:    "",
		Aliases:  []string{"http-rpc.cors-domain"},
		EnvVars:  []string{"KLAYTN_RPCCORSDOMAIN"},
		Category: "API AND CONSOLE",
	}
	RPCVirtualHostsFlag = &cli.StringFlag{
		Name:     "rpcvhosts",
		Usage:    "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		Value:    strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
		Aliases:  []string{"http-rpc.vhosts"},
		EnvVars:  []string{"KLAYTN_RPCVHOSTS"},
		Category: "API AND CONSOLE",
	}
	RPCApiFlag = &cli.StringFlag{
		Name:     "rpcapi",
		Usage:    "API's offered over the HTTP-RPC interface",
		Value:    "",
		Aliases:  []string{"http-rpc.api"},
		EnvVars:  []string{"KLAYTN_RPCAPI"},
		Category: "API AND CONSOLE",
	}
	RPCGlobalGasCap = &cli.Uint64Flag{
		Name:     "rpc.gascap",
		Usage:    "Sets a cap on gas that can be used in klay_call/estimateGas",
		Aliases:  []string{"http-rpc.gascap"},
		EnvVars:  []string{"KLAYTN_RPC_GASCAP"},
		Category: "API AND CONSOLE",
	}
	RPCGlobalEVMTimeoutFlag = &cli.DurationFlag{
		Name:     "rpc.evmtimeout",
		Usage:    "Sets a timeout used for eth_call (0=infinite)",
		Aliases:  []string{"http-rpc.evmtimeout"},
		EnvVars:  []string{"KLAYTN_RPC_EVMTIMEOUT"},
		Category: "API AND CONSOLE",
	}
	RPCGlobalEthTxFeeCapFlag = &cli.Float64Flag{
		Name:     "rpc.ethtxfeecap",
		Usage:    "Sets a cap on transaction fee (in klay) that can be sent via the eth namespace RPC APIs (0 = no cap)",
		Aliases:  []string{"http-rpc.eth-tx-feecap"},
		EnvVars:  []string{"KLAYTN_RPC_ETHTXFEECAP"},
		Category: "API AND CONSOLE",
	}
	RPCConcurrencyLimit = &cli.IntFlag{
		Name:     "rpc.concurrencylimit",
		Usage:    "Sets a limit of concurrent connection number of HTTP-RPC server",
		Value:    rpc.ConcurrencyLimit,
		Aliases:  []string{"http-rpc.concurrency-limit"},
		EnvVars:  []string{"KLAYTN_RPC_CONCURRENCYLIMIT"},
		Category: "API AND CONSOLE",
	}
	RPCNonEthCompatibleFlag = &cli.BoolFlag{
		Name:     "rpc.eth.noncompatible",
		Usage:    "Disables the eth namespace API return formatting for compatibility",
		Aliases:  []string{"http-rpc.eth-noncompatible"},
		EnvVars:  []string{"KLAYTN_RPC_ETH_NONCOMPATIBLE"},
		Category: "API AND CONSOLE",
	}
	RPCReadTimeout = &cli.IntFlag{
		Name:     "rpcreadtimeout",
		Usage:    "HTTP-RPC server read timeout (seconds)",
		Value:    int(rpc.DefaultHTTPTimeouts.ReadTimeout / time.Second),
		Aliases:  []string{"http-rpc.read-timeout"},
		EnvVars:  []string{"KLAYTN_RPCREADTIMEOUT"},
		Category: "API AND CONSOLE",
	}
	RPCWriteTimeoutFlag = &cli.IntFlag{
		Name:     "rpcwritetimeout",
		Usage:    "HTTP-RPC server write timeout (seconds)",
		Value:    int(rpc.DefaultHTTPTimeouts.WriteTimeout / time.Second),
		Aliases:  []string{"http-rpc.write-timeout"},
		EnvVars:  []string{"KLAYTN_RPCWRITETIMEOUT"},
		Category: "API AND CONSOLE",
	}
	RPCIdleTimeoutFlag = &cli.IntFlag{
		Name:     "rpcidletimeout",
		Usage:    "HTTP-RPC server idle timeout (seconds)",
		Value:    int(rpc.DefaultHTTPTimeouts.IdleTimeout / time.Second),
		Aliases:  []string{"http-rpc.idle-timeout"},
		EnvVars:  []string{"KLAYTN_RPCIDLETIMEOUT"},
		Category: "API AND CONSOLE",
	}
	RPCExecutionTimeoutFlag = &cli.IntFlag{
		Name:     "rpcexecutiontimeout",
		Usage:    "HTTP-RPC server execution timeout (seconds)",
		Value:    int(rpc.DefaultHTTPTimeouts.ExecutionTimeout / time.Second),
		Aliases:  []string{"http-rpc.execution-timeout"},
		EnvVars:  []string{"KLAYTN_RPCEXECUTIONTIMEOUT"},
		Category: "API AND CONSOLE",
	}
	RPCUpstreamArchiveENFlag = &cli.StringFlag{
		Name:     "upstream-en",
		Usage:    "upstream archive mode EN endpoint",
		Aliases:  []string{"rpc.upstream-en"},
		EnvVars:  []string{"KLAYTN_RPC_UPSTREAM_EN"},
		Category: "API AND CONSOLE",
	}

	WSEnabledFlag = &cli.BoolFlag{
		Name:     "ws",
		Usage:    "Enable the WS-RPC server",
		Aliases:  []string{"ws-rpc.enable"},
		EnvVars:  []string{"KLAYTN_WS"},
		Category: "API AND CONSOLE",
	}
	WSListenAddrFlag = &cli.StringFlag{
		Name:     "wsaddr",
		Usage:    "WS-RPC server listening interface",
		Value:    node.DefaultWSHost,
		Aliases:  []string{"ws-rpc.addr"},
		EnvVars:  []string{"KLAYTN_WSADDR"},
		Category: "API AND CONSOLE",
	}
	WSPortFlag = &cli.IntFlag{
		Name:     "wsport",
		Usage:    "WS-RPC server listening port",
		Value:    node.DefaultWSPort,
		Aliases:  []string{"ws-rpc.port"},
		EnvVars:  []string{"KLAYTN_WSPORT"},
		Category: "API AND CONSOLE",
	}
	WSApiFlag = &cli.StringFlag{
		Name:     "wsapi",
		Usage:    "API's offered over the WS-RPC interface",
		Value:    "",
		Aliases:  []string{"ws-rpc.api"},
		EnvVars:  []string{"KLAYTN_WSAPI"},
		Category: "API AND CONSOLE",
	}
	WSAllowedOriginsFlag = &cli.StringFlag{
		Name:     "wsorigins",
		Usage:    "Origins from which to accept websockets requests",
		Value:    "",
		Aliases:  []string{"ws-rpc.origins"},
		EnvVars:  []string{"KLAYTN_WSORIGINS"},
		Category: "API AND CONSOLE",
	}
	WSMaxSubscriptionPerConn = &cli.IntFlag{
		Name:     "wsmaxsubscriptionperconn",
		Usage:    "Allowed maximum subscription number per a websocket connection",
		Value:    int(rpc.MaxSubscriptionPerWSConn),
		Aliases:  []string{"ws-rpc.max-subscription-per-conn"},
		EnvVars:  []string{"KLAYTN_WSMAXSUBSCRIPTIONPERCONN"},
		Category: "API AND CONSOLE",
	}
	WSReadDeadLine = &cli.Int64Flag{
		Name:     "wsreaddeadline",
		Usage:    "Set the read deadline on the underlying network connection in seconds. 0 means read will not timeout",
		Value:    rpc.WebsocketReadDeadline,
		Aliases:  []string{"ws-rpc.read-deadline"},
		EnvVars:  []string{"KLAYTN_WSREADDEADLINE"},
		Category: "API AND CONSOLE",
	}
	WSWriteDeadLine = &cli.Int64Flag{
		Name:     "wswritedeadline",
		Usage:    "Set the Write deadline on the underlying network connection in seconds. 0 means write will not timeout",
		Value:    rpc.WebsocketWriteDeadline,
		Aliases:  []string{"ws-rpc.write-deadline"},
		EnvVars:  []string{"KLAYTN_WSWRITEDEADLINE"},
		Category: "API AND CONSOLE",
	}
	WSMaxConnections = &cli.IntFlag{
		Name:     "wsmaxconnections",
		Usage:    "Allowed maximum websocket connection number",
		Value:    3000,
		Aliases:  []string{"ws-rpc.max-connections"},
		EnvVars:  []string{"KLAYTN_WSMAXCONNECTIONS"},
		Category: "API AND CONSOLE",
	}
	GRPCEnabledFlag = &cli.BoolFlag{
		Name:     "grpc",
		Usage:    "Enable the gRPC server",
		Aliases:  []string{"g-rpc.enable"},
		EnvVars:  []string{"KLAYTN_GRPC"},
		Category: "API AND CONSOLE",
	}
	GRPCListenAddrFlag = &cli.StringFlag{
		Name:     "grpcaddr",
		Usage:    "gRPC server listening interface",
		Value:    node.DefaultGRPCHost,
		Aliases:  []string{"g-rpc.addr"},
		EnvVars:  []string{"KLAYTN_GRPCADDR"},
		Category: "API AND CONSOLE",
	}
	GRPCPortFlag = &cli.IntFlag{
		Name:     "grpcport",
		Usage:    "gRPC server listening port",
		Value:    node.DefaultGRPCPort,
		Aliases:  []string{"g-rpc.port"},
		EnvVars:  []string{"KLAYTN_GRPCPORT"},
		Category: "API AND CONSOLE",
	}
	IPCDisabledFlag = &cli.BoolFlag{
		Name:     "ipcdisable",
		Usage:    "Disable the IPC-RPC server",
		Aliases:  []string{"ipc.disable"},
		EnvVars:  []string{"KLAYTN_IPCDISABLE"},
		Category: "API AND CONSOLE",
	}
	IPCPathFlag = &cli.PathFlag{
		Name:     "ipcpath",
		Usage:    "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
		Aliases:  []string{"ipc.path"},
		EnvVars:  []string{"KLAYTN_IPCPATH"},
		Category: "API AND CONSOLE",
	}

	// ATM the url is left to the user and deployment to
	JSpathFlag = &cli.StringFlag{
		Name:     "jspath",
		Usage:    "JavaScript root path for `loadScript`",
		Value:    ".",
		Aliases:  []string{"console.js-path"},
		EnvVars:  []string{"KLAYTN_JSPATH"},
		Category: "API AND CONSOLE",
	}
	ExecFlag = &cli.StringFlag{
		Name:     "exec",
		Usage:    "Execute JavaScript statement",
		Aliases:  []string{"console.exec"},
		EnvVars:  []string{"KLAYTN_EXEC"},
		Category: "API AND CONSOLE",
	}
	PreloadJSFlag = &cli.StringFlag{
		Name:     "preload",
		Usage:    "Comma separated list of JavaScript files to preload into the console",
		Aliases:  []string{"console.preload"},
		EnvVars:  []string{"KLAYTN_PRELOAD"},
		Category: "API AND CONSOLE",
	}
	APIFilterGetLogsDeadlineFlag = &cli.DurationFlag{
		Name:     "api.filter.getLogs.deadline",
		Usage:    "Execution deadline for log collecting filter APIs",
		Value:    filters.GetLogsDeadline,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_API_FILTER_GETLOGS_DEADLINE"},
		Category: "API AND CONSOLE",
	}
	APIFilterGetLogsMaxItemsFlag = &cli.IntFlag{
		Name:     "api.filter.getLogs.maxitems",
		Usage:    "Maximum allowed number of return items for log collecting filter API",
		Value:    filters.GetLogsMaxItems,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_API_FILTER_GETLOGS_MAXITEMS"},
		Category: "API AND CONSOLE",
	}
	UnsafeDebugDisableFlag = &cli.BoolFlag{
		Name:     "rpc.unsafe-debug.disable",
		Usage:    "Disable unsafe debug APIs (traceTransaction, traceChain, ...).",
		Aliases:  []string{"http-rpc.unsafe-debug.disable"},
		EnvVars:  []string{"KLAYTN_RPC_UNSAFE_DEBUG_DISABLE"},
		Category: "API AND CONSOLE",
	}
	// TODO-klaytn: Consider limiting the non-debug heavy apis.
	HeavyDebugRequestLimitFlag = &cli.IntFlag{
		Name:     "rpc.unsafe-debug.heavy-debug.request-limit",
		Usage:    "Limit the maximum number of heavy debug api requests. Works with unsafe-debug only.",
		Value:    50,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_RPC_UNSAFE_DEBUG_HEAVY_DEBUG_REQUEST_LIMIT"},
		Category: "API AND CONSOLE",
	}
	StateRegenerationTimeLimitFlag = &cli.DurationFlag{
		Name:     "rpc.unsafe-debug.state-regeneration.time-limit",
		Usage:    "Limit the state regeneration time. Works with unsafe-debug only.",
		Value:    60 * time.Second,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_RPC_UNSAFE_DEBUG_STATE_REGENERATION_TIME_LIMIT"},
		Category: "API AND CONSOLE",
	}

	// Network Settings
	NodeTypeFlag = &cli.StringFlag{
		Name:    "nodetype",
		Usage:   "Klaytn node type (consensus node (cn), proxy node (pn), endpoint node (en))",
		Value:   "en",
		Aliases: []string{},
		EnvVars: []string{"KLAYTN_NODETYPE"},
	}
	MaxConnectionsFlag = &cli.IntFlag{
		Name:     "maxconnections",
		Usage:    "Maximum number of physical connections. All single channel peers can be maxconnections peers. All multi channel peers can be maxconnections/2 peers. (network disabled if set to 0)",
		Value:    node.DefaultMaxPhysicalConnections,
		Aliases:  []string{"p2p.max-connections"},
		EnvVars:  []string{"KLAYTN_MAXCONNECTIONS"},
		Category: "NETWORK",
	}
	MaxPendingPeersFlag = &cli.IntFlag{
		Name:     "maxpendpeers",
		Usage:    "Maximum number of pending connection attempts (defaults used if set to 0)",
		Value:    0,
		Aliases:  []string{"p2p.max-pend-peers"},
		EnvVars:  []string{"KLAYTN_MAXPENDPEERS"},
		Category: "NETWORK",
	}
	ListenPortFlag = &cli.IntFlag{
		Name:     "port",
		Usage:    "Network listening port",
		Value:    node.DefaultP2PPort,
		Aliases:  []string{"p2p.port"},
		EnvVars:  []string{"KLAYTN_PORT"},
		Category: "NETWORK",
	}
	SubListenPortFlag = &cli.IntFlag{
		Name:     "subport",
		Usage:    "Network sub listening port",
		Value:    node.DefaultP2PSubPort,
		Aliases:  []string{"p2p.sub-port"},
		EnvVars:  []string{"KLAYTN_SUBPORT"},
		Category: "NETWORK",
	}
	MultiChannelUseFlag = &cli.BoolFlag{
		Name:     "multichannel",
		Usage:    "Create a dedicated channel for block propagation",
		Aliases:  []string{"p2p.multi-channel"},
		EnvVars:  []string{"KLAYTN_MULTICHANNEL"},
		Category: "NETWORK",
	}
	BootnodesFlag = &cli.StringFlag{
		Name:     "bootnodes",
		Usage:    "Comma separated kni URLs for P2P discovery bootstrap",
		Value:    "",
		Aliases:  []string{"p2p.bootnodes"},
		EnvVars:  []string{"KLAYTN_BOOTNODES"},
		Category: "NETWORK",
	}
	NodeKeyFileFlag = &cli.StringFlag{
		Name:     "nodekey",
		Usage:    "P2P node key file",
		Aliases:  []string{"p2p.node-key"},
		EnvVars:  []string{"KLAYTN_NODEKEY"},
		Category: "NETWORK",
	}
	NodeKeyHexFlag = &cli.StringFlag{
		Name:     "nodekeyhex",
		Usage:    "P2P node key as hex (for testing)",
		Aliases:  []string{"p2p.node-key-hex"},
		EnvVars:  []string{"KLAYTN_NODEKEYHEX"},
		Category: "NETWORK",
	}
	BlsNodeKeyFileFlag = &cli.StringFlag{
		Name:     "bls-nodekey",
		Usage:    "Consensus BLS node key file",
		Aliases:  []string{"p2p.bls-node-key"},
		EnvVars:  []string{"KLAYTN_BLS_NODEKEY"},
		Category: "NETWORK",
	}
	BlsNodeKeyHexFlag = &cli.StringFlag{
		Name:     "bls-nodekeyhex",
		Usage:    "Consensus BLS node key in hex (for testing)",
		Aliases:  []string{"p2p.bls-node-key-hex"},
		EnvVars:  []string{"KLAYTN_BLS_NODEKEYHEX"},
		Category: "NETWORK",
	}
	BlsNodeKeystoreFileFlag = &cli.StringFlag{
		Name:     "bls-nodekeystore",
		Usage:    "Consensus BLS node keystore JSON file",
		Aliases:  []string{"p2p.bls-node-keystore"},
		EnvVars:  []string{"KLAYTN_BLS_NODEKEYSTORE"},
		Category: "NETWORK",
	}
	NATFlag = &cli.StringFlag{
		Name:     "nat",
		Usage:    "NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>)",
		Value:    "any",
		Aliases:  []string{"p2p.nat"},
		EnvVars:  []string{"KLAYTN_NAT"},
		Category: "NETWORK",
	}
	NoDiscoverFlag = &cli.BoolFlag{
		Name:     "nodiscover",
		Usage:    "Disables the peer discovery mechanism (manual peer addition)",
		Aliases:  []string{"p2p.no-discover"},
		EnvVars:  []string{"KLAYTN_NODISCOVER"},
		Category: "NETWORK",
	}
	NetrestrictFlag = &cli.StringFlag{
		Name:     "netrestrict",
		Usage:    "Restricts network communication to the given IP network (CIDR masks)",
		Aliases:  []string{"p2p.net-restrict"},
		EnvVars:  []string{"KLAYTN_NETRESTRICT"},
		Category: "NETWORK",
	}
	RWTimerIntervalFlag = &cli.Uint64Flag{
		Name:     "rwtimerinterval",
		Usage:    "Interval of using rw timer to check if it works well",
		Value:    1000,
		Aliases:  []string{"p2p.rw-timer-interval"},
		EnvVars:  []string{"KLAYTN_RWTIMERINTERVAL"},
		Category: "NETWORK",
	}
	RWTimerWaitTimeFlag = &cli.DurationFlag{
		Name:     "rwtimerwaittime",
		Usage:    "Wait time the rw timer waits for message writing",
		Value:    15 * time.Second,
		Aliases:  []string{"p2p.rw-timer-wait-time"},
		EnvVars:  []string{"KLAYTN_RWTIMERWAITTIME"},
		Category: "NETWORK",
	}
	MaxRequestContentLengthFlag = &cli.IntFlag{
		Name:     "maxRequestContentLength",
		Usage:    "Max request content length in byte for http, websocket and gRPC",
		Value:    common.MaxRequestContentLength,
		Aliases:  []string{"p2p.max-request-content-length"},
		EnvVars:  []string{"KLAYTN_MAXREQUESTCONTENTLENGTH"},
		Category: "API AND CONSOLE",
	}

	CypressFlag = &cli.BoolFlag{
		Name:     "cypress",
		Usage:    "Pre-configured Klaytn Cypress network",
		Aliases:  []string{"p2p.cypress"},
		EnvVars:  []string{"KLAYTN_CYPRESS"},
		Category: "NETWORK",
	}
	// Baobab bootnodes setting
	BaobabFlag = &cli.BoolFlag{
		Name:     "baobab",
		Usage:    "Pre-configured Klaytn baobab network",
		Aliases:  []string{"p2p.baobab"},
		EnvVars:  []string{"KLAYTN_BAOBAB"},
		Category: "NETWORK",
	}
	// Bootnode's settings
	AuthorizedNodesFlag = &cli.StringFlag{
		Name:    "authorized-nodes",
		Usage:   "Comma separated kni URLs for authorized nodes list",
		Value:   "",
		Aliases: []string{"common.authorized-nodes"},
		EnvVars: []string{"KLAYTN_AUTHORIZED_NODES"},
	}
	// TODO-Klaytn-Bootnode the boodnode flags should be updated when it is implemented
	BNAddrFlag = &cli.StringFlag{
		Name:    "bnaddr",
		Usage:   `udp address to use node discovery`,
		Value:   ":32323",
		Aliases: []string{"p2p.bn-addr"},
		EnvVars: []string{"KLAYTN_BNADDR"},
	}
	GenKeyFlag = &cli.StringFlag{
		Name:     "genkey",
		Usage:    "generate a node private key and write to given filename",
		Aliases:  []string{"common.gen-key-path"},
		EnvVars:  []string{"KLAYTN_GENKEY"},
		Category: "MISC",
	}
	WriteAddressFlag = &cli.BoolFlag{
		Name:     "writeaddress",
		Usage:    `write out the node's public key which is given by "--nodekey" or "--nodekeyhex"`,
		Aliases:  []string{"common.write-address"},
		EnvVars:  []string{"KLAYTN_WRITEADDRESS"},
		Category: "MISC",
	}
	// ServiceChain's settings
	AnchoringPeriodFlag = &cli.Uint64Flag{
		Name:     "chaintxperiod",
		Usage:    "The period to make and send a chain transaction to the parent chain",
		Value:    1,
		Aliases:  []string{"servicechain.chain-tx-period"},
		EnvVars:  []string{"KLAYTN_CHAINTXPERIOD"},
		Category: "SERVICECHAIN",
	}
	SentChainTxsLimit = &cli.Uint64Flag{
		Name:     "chaintxlimit",
		Usage:    "Number of service chain transactions stored for resending",
		Value:    100,
		Aliases:  []string{"servicechain.chain-tx-limit"},
		EnvVars:  []string{"KLAYTN_CHAINTXLIMIT"},
		Category: "SERVICECHAIN",
	}
	MainBridgeFlag = &cli.BoolFlag{
		Name:     "mainbridge",
		Usage:    "Enable main bridge service for service chain",
		Aliases:  []string{"servicechain.mainbridge"},
		EnvVars:  []string{"KLAYTN_MAINBRIDGE"},
		Category: "SERVICECHAIN",
	}
	SubBridgeFlag = &cli.BoolFlag{
		Name:     "subbridge",
		Usage:    "Enable sub bridge service for service chain",
		Aliases:  []string{"servicechain.subbridge"},
		EnvVars:  []string{"KLAYTN_SUBBRIDGE"},
		Category: "SERVICECHAIN",
	}
	MainBridgeListenPortFlag = &cli.IntFlag{
		Name:     "mainbridgeport",
		Usage:    "main bridge listen port",
		Value:    50505,
		Aliases:  []string{"servicechain.mainbridge-port"},
		EnvVars:  []string{"KLAYTN_MAINBRIDGEPORT"},
		Category: "SERVICECHAIN",
	}
	SubBridgeListenPortFlag = &cli.IntFlag{
		Name:     "subbridgeport",
		Usage:    "sub bridge listen port",
		Value:    50506,
		Aliases:  []string{"servicechain.subbridge-port"},
		EnvVars:  []string{"KLAYTN_SUBBRIDGEPORT"},
		Category: "SERVICECHAIN",
	}
	ParentChainIDFlag = &cli.IntFlag{
		Name:     "parentchainid",
		Usage:    "parent chain ID",
		Value:    8217, // Klaytn mainnet chain ID
		Aliases:  []string{"servicechain.parent-chainid"},
		EnvVars:  []string{"KLAYTN_PARENTCHAINID"},
		Category: "SERVICECHAIN",
	}
	VTRecoveryFlag = &cli.BoolFlag{
		Name:     "vtrecovery",
		Usage:    "Enable value transfer recovery (default: false)",
		Aliases:  []string{"servicechain.vt-recovery"},
		EnvVars:  []string{"KLAYTN_VTRECOVERY"},
		Category: "SERVICECHAIN",
	}
	VTRecoveryIntervalFlag = &cli.Uint64Flag{
		Name:     "vtrecoveryinterval",
		Usage:    "Set the value transfer recovery interval (seconds)",
		Value:    5,
		Aliases:  []string{"servicechain.vt-recovery-interval"},
		EnvVars:  []string{"KLAYTN_VTRECOVERYINTERVAL"},
		Category: "SERVICECHAIN",
	}
	ServiceChainParentOperatorTxGasLimitFlag = &cli.Uint64Flag{
		Name:     "sc.parentoperator.gaslimit",
		Usage:    "Set the default value of gas limit for transactions made by bridge parent operator",
		Value:    10000000,
		Aliases:  []string{"servicechain.parent-operator-gaslimit"},
		EnvVars:  []string{"KLAYTN_SC_PARENTOPERATOR_GASLIMIT"},
		Category: "SERVICECHAIN",
	}
	ServiceChainChildOperatorTxGasLimitFlag = &cli.Uint64Flag{
		Name:     "sc.childoperator.gaslimit",
		Usage:    "Set the default value of gas limit for transactions made by bridge child operator",
		Value:    10000000,
		Aliases:  []string{"servicechain.child-operator-gaslimit"},
		EnvVars:  []string{"KLAYTN_SC_CHILDOPERATOR_GASLIMIT"},
		Category: "SERVICECHAIN",
	}
	ServiceChainNewAccountFlag = &cli.BoolFlag{
		Name:     "scnewaccount",
		Usage:    "Enable account creation for the service chain (default: false). If set true, generated account can't be synced with the parent chain.",
		Aliases:  []string{"servicechain.new-account"},
		EnvVars:  []string{"KLAYTN_SCNEWACCOUNT"},
		Category: "SERVICECHAIN",
	}
	ServiceChainAnchoringFlag = &cli.BoolFlag{
		Name:     "anchoring",
		Usage:    "Enable anchoring for service chain",
		Aliases:  []string{"servicechain.anchoring"},
		EnvVars:  []string{"KLAYTN_ANCHORING"},
		Category: "SERVICECHAIN",
	}
	// TODO-klaytn: need to check if deprecated.
	ServiceChainConsensusFlag = &cli.StringFlag{
		Name:    "scconsensus",
		Usage:   "Set the service chain consensus (\"istanbul\", \"clique\")",
		Value:   "istanbul",
		Aliases: []string{"servicechain.consensus"},
		EnvVars: []string{"KLAYTN_SCCONSENSUS"},
	}

	// KAS
	KASServiceChainAnchorFlag = &cli.BoolFlag{
		Name:     "kas.sc.anchor",
		Usage:    "Enable KAS anchoring for service chain",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_SC_ANCHOR"},
		Category: "SERVICECHAIN",
	}
	KASServiceChainAnchorPeriodFlag = &cli.Uint64Flag{
		Name:     "kas.sc.anchor.period",
		Usage:    "The period to anchor service chain blocks to KAS",
		Value:    1,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_SC_ANCHOR_PERIOD"},
		Category: "SERVICECHAIN",
	}
	KASServiceChainAnchorUrlFlag = &cli.StringFlag{
		Name:     "kas.sc.anchor.url",
		Usage:    "The url for KAS anchor",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_SC_ANCHOR_URL"},
		Category: "SERVICECHAIN",
	}
	KASServiceChainAnchorOperatorFlag = &cli.StringFlag{
		Name:     "kas.sc.anchor.operator",
		Usage:    "The operator address for KAS anchor",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_SC_ANCHOR_OPERATOR"},
		Category: "SERVICECHAIN",
	}
	KASServiceChainAnchorRequestTimeoutFlag = &cli.DurationFlag{
		Name:     "kas.sc.anchor.request.timeout",
		Usage:    "The reuqest timeout for KAS Anchoring API call",
		Value:    500 * time.Millisecond,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_SC_ANCHOR_REQUEST_TIMEOUT"},
		Category: "SERVICECHAIN",
	}
	KASServiceChainXChainIdFlag = &cli.StringFlag{
		Name:     "kas.x-chain-id",
		Usage:    "The x-chain-id for KAS",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_X_CHAIN_ID"},
		Category: "SERVICECHAIN",
	}
	KASServiceChainAccessKeyFlag = &cli.StringFlag{
		Name:     "kas.accesskey",
		Usage:    "The access key id for KAS",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_ACCESSKEY"},
		Category: "SERVICECHAIN",
	}
	KASServiceChainSecretKeyFlag = &cli.StringFlag{
		Name:     "kas.secretkey",
		Usage:    "The secret key for KAS",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_KAS_SECRETKEY"},
		Category: "SERVICECHAIN",
	}

	// ChainDataFetcher
	EnableChainDataFetcherFlag = &cli.BoolFlag{
		Name:     "chaindatafetcher",
		Usage:    "Enable the ChainDataFetcher Service",
		Aliases:  []string{"chain-data-fetcher.enable"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherMode = &cli.StringFlag{
		Name:     "chaindatafetcher.mode",
		Usage:    "The mode of chaindatafetcher (\"kas\", \"kafka\")",
		Value:    "kas",
		Aliases:  []string{"chain-data-fetcher.mode"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_MODE"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherNoDefault = &cli.BoolFlag{
		Name:     "chaindatafetcher.no.default",
		Usage:    "Turn off the starting of the chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.no-default"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_NO_DEFAULT"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherNumHandlers = &cli.IntFlag{
		Name:     "chaindatafetcher.num.handlers",
		Usage:    "Number of chaindata handlers",
		Value:    chaindatafetcher.DefaultNumHandlers,
		Aliases:  []string{"chain-data-fetcher.num-handlers"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_NUM_HANDLERS"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherJobChannelSize = &cli.IntFlag{
		Name:     "chaindatafetcher.job.channel.size",
		Usage:    "Job channel size",
		Value:    chaindatafetcher.DefaultJobChannelSize,
		Aliases:  []string{"chain-data-fetcher.job-channel-size"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_JOB_CHANNEL_SIZE"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherChainEventSizeFlag = &cli.IntFlag{
		Name:     "chaindatafetcher.block.channel.size",
		Usage:    "Block received channel size",
		Value:    chaindatafetcher.DefaultJobChannelSize,
		Aliases:  []string{"chain-data-fetcher.block-channel-size"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_BLOCK_CHANNEL_SIZE"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherMaxProcessingDataSize = &cli.IntFlag{
		Name:     "chaindatafetcher.max.processing.data.size",
		Usage:    "Maximum size of processing data before requesting range fetching of blocks (in MB)",
		Value:    chaindatafetcher.DefaultMaxProcessingDataSize,
		Aliases:  []string{"chain-data-fetcher.max-processing-data-size"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_MAX_PROCESSING_DATA_SIZE"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASDBHostFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.db.host",
		Usage:    "KAS specific DB host in chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.kas.db.host"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_DB_HOST"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASDBPortFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.db.port",
		Usage:    "KAS specific DB port in chaindatafetcher",
		Value:    chaindatafetcher.DefaultDBPort,
		Aliases:  []string{"chain-data-fetcher.kas.db.port"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_DB_PORT"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASDBNameFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.db.name",
		Usage:    "KAS specific DB name in chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.kas.db.name"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_DB_NAME"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASDBUserFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.db.user",
		Usage:    "KAS specific DB user in chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.kas.db.user"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_DB_USER"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASDBPasswordFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.db.password",
		Usage:    "KAS specific DB password in chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.kas.db.password"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_DB_PASSWORD"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASCacheUse = &cli.BoolFlag{
		Name:     "chaindatafetcher.kas.cache.use",
		Usage:    "Enable KAS cache invalidation",
		Aliases:  []string{"chain-data-fetcher.kas.cache.use"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_CACHE_USE"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASCacheURLFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.cache.url",
		Usage:    "KAS specific cache invalidate API endpoint in chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.kas.cache.url"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_CACHE_URL"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASXChainIdFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.xchainid",
		Usage:    "KAS specific header x-chain-id in chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.kas.xchainid"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_XCHAINID"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKASBasicAuthParamFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kas.basic.auth.param",
		Usage:    "KAS specific header basic authorization parameter in chaindatafetcher",
		Aliases:  []string{"chain-data-fetcher.kas.basic.auth.param"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAS_BASIC_AUTH_PARAM"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaBrokersFlag = &cli.StringSliceFlag{
		Name:     "chaindatafetcher.kafka.brokers",
		Usage:    "Kafka broker URL list",
		Aliases:  []string{"chain-data-fetcher.kafka.brokers"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_BROKERS"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaTopicEnvironmentFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kafka.topic.environment",
		Usage:    "Kafka topic environment prefix",
		Value:    kafka.DefaultTopicEnvironmentName,
		Aliases:  []string{"chain-data-fetcher.kafka.topic.environment"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_TOPIC_ENVIRONMENT"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaTopicResourceFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kafka.topic.resource",
		Usage:    "Kafka topic resource name",
		Value:    kafka.DefaultTopicResourceName,
		Aliases:  []string{"chain-data-fetcher.kafka.topic.resource"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_TOPIC_RESOURCE"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaReplicasFlag = &cli.Int64Flag{
		Name:     "chaindatafetcher.kafka.replicas",
		Usage:    "Kafka partition replication factor",
		Value:    kafka.DefaultReplicas,
		Aliases:  []string{"chain-data-fetcher.kafka.replicas"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_REPLICAS"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaPartitionsFlag = &cli.IntFlag{
		Name:     "chaindatafetcher.kafka.partitions",
		Usage:    "The number of partitions in a topic",
		Value:    kafka.DefaultPartitions,
		Aliases:  []string{"chain-data-fetcher.kafka.partitions"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_PARTITIONS"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaMaxMessageBytesFlag = &cli.IntFlag{
		Name:     "chaindatafetcher.kafka.max.message.bytes",
		Usage:    "The max size of a message produced by Kafka producer ",
		Value:    kafka.DefaultMaxMessageBytes,
		Aliases:  []string{"chain-data-fetcher.kafka.max-message-bytes"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_MAX_MESSAGE_BYTES"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaSegmentSizeBytesFlag = &cli.IntFlag{
		Name:     "chaindatafetcher.kafka.segment.size",
		Usage:    "The kafka data segment size (in byte)",
		Value:    kafka.DefaultSegmentSizeBytes,
		Aliases:  []string{"chain-data-fetcher.kafka.segment-size"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_SEGMENT_SIZE"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaRequiredAcksFlag = &cli.IntFlag{
		Name:     "chaindatafetcher.kafka.required.acks",
		Usage:    "The level of acknowledgement reliability needed from Kafka broker (0: NoResponse, 1: WaitForLocal, -1: WaitForAll)",
		Value:    kafka.DefaultRequiredAcks,
		Aliases:  []string{"chain-data-fetcher.kafka.required-acks"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_REQUIRED_ACKS"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaMessageVersionFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kafka.msg.version",
		Usage:    "The version of Kafka message",
		Value:    kafka.DefaultKafkaMessageVersion,
		Aliases:  []string{"chain-data-fetcher.kafka.msg-version"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_MSG_VERSION"},
		Category: "CHAINDATAFETCHER",
	}
	ChainDataFetcherKafkaProducerIdFlag = &cli.StringFlag{
		Name:     "chaindatafetcher.kafka.producer.id",
		Usage:    "The identifier of kafka message producer",
		Value:    kafka.GetDefaultProducerId(),
		Aliases:  []string{"chain-data-fetcher.kafka.producer-id"},
		EnvVars:  []string{"KLAYTN_CHAINDATAFETCHER_KAFKA_PRODUCER_ID"},
		Category: "CHAINDATAFETCHER",
	}
	// DBSyncer
	EnableDBSyncerFlag = &cli.BoolFlag{
		Name:     "dbsyncer",
		Usage:    "Enable the DBSyncer",
		Aliases:  []string{"db-syncer.enable"},
		EnvVars:  []string{"KLAYTN_DBSYNCER"},
		Category: "DATABASE SYNCER",
	}
	DBHostFlag = &cli.StringFlag{
		Name:     "dbsyncer.db.host",
		Usage:    "db.host in dbsyncer",
		Aliases:  []string{"db-syncer.db.host"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_HOST"},
		Category: "DATABASE SYNCER",
	}
	DBPortFlag = &cli.StringFlag{
		Name:     "dbsyncer.db.port",
		Usage:    "db.port in dbsyncer",
		Value:    "3306",
		Aliases:  []string{"db-syncer.db.port"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_PORT"},
		Category: "DATABASE SYNCER",
	}
	DBNameFlag = &cli.StringFlag{
		Name:     "dbsyncer.db.name",
		Usage:    "db.name in dbsyncer",
		Aliases:  []string{"db-syncer.db.name"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_NAME"},
		Category: "DATABASE SYNCER",
	}
	DBUserFlag = &cli.StringFlag{
		Name:     "dbsyncer.db.user",
		Usage:    "db.user in dbsyncer",
		Aliases:  []string{"db-syncer.db.user"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_USER"},
		Category: "DATABASE SYNCER",
	}
	DBPasswordFlag = &cli.StringFlag{
		Name:     "dbsyncer.db.password",
		Usage:    "db.password in dbsyncer",
		Aliases:  []string{"db-syncer.db.password"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_PASSWORD"},
		Category: "DATABASE SYNCER",
	}
	EnabledLogModeFlag = &cli.BoolFlag{
		Name:     "dbsyncer.logmode",
		Usage:    "Enable the dbsyncer logmode",
		Aliases:  []string{"db-syncer.log-mode"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_LOGMODE"},
		Category: "DATABASE SYNCER",
	}
	MaxIdleConnsFlag = &cli.IntFlag{
		Name:     "dbsyncer.db.max.idle",
		Usage:    "The maximum number of connections in the idle connection pool",
		Value:    50,
		Aliases:  []string{"db-syncer.db.max-idle"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_MAX_IDLE"},
		Category: "DATABASE SYNCER",
	}
	MaxOpenConnsFlag = &cli.IntFlag{
		Name:     "dbsyncer.db.max.open",
		Usage:    "The maximum number of open connections to the database",
		Value:    30,
		Aliases:  []string{"db-syncer.db.max-open"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_MAX_OPEN"},
		Category: "DATABASE SYNCER",
	}
	ConnMaxLifeTimeFlag = &cli.DurationFlag{
		Name:     "dbsyncer.db.max.lifetime",
		Usage:    "The maximum amount of time a connection may be reused (default : 1h), ex: 300ms, 2h45m, 60s, ...",
		Value:    1 * time.Hour,
		Aliases:  []string{"db-syncer.db.max-lifetime"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_DB_MAX_LIFETIME"},
		Category: "DATABASE SYNCER",
	}
	BlockSyncChannelSizeFlag = &cli.IntFlag{
		Name:     "dbsyncer.block.channel.size",
		Usage:    "Block received channel size",
		Value:    5,
		Aliases:  []string{"db-syncer.block-channel-size"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_BLOCK_CHANNEL_SIZE"},
		Category: "DATABASE SYNCER",
	}
	DBSyncerModeFlag = &cli.StringFlag{
		Name:     "dbsyncer.mode",
		Usage:    "The mode of dbsyncer is way which handle block/tx data to insert db (multi, single, context)",
		Value:    "multi",
		Aliases:  []string{"db-syncer.mode"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_MODE"},
		Category: "DATABASE SYNCER",
	}
	GenQueryThreadFlag = &cli.IntFlag{
		Name:     "dbsyncer.genquery.th",
		Usage:    "The amount of thread of generation query in multi mode",
		Value:    50,
		Aliases:  []string{"db-syncer.genquery-th"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_GENQUERY_TH"},
		Category: "DATABASE SYNCER",
	}
	InsertThreadFlag = &cli.IntFlag{
		Name:     "dbsyncer.insert.th",
		Usage:    "The amount of thread of insert operation in multi mode",
		Value:    30,
		Aliases:  []string{"db-syncer.insert-thread"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_INSERT_TH"},
		Category: "DATABASE SYNCER",
	}
	BulkInsertSizeFlag = &cli.IntFlag{
		Name:     "dbsyncer.bulk.size",
		Usage:    "The amount of row for bulk-insert",
		Value:    200,
		Aliases:  []string{"db-syncer.bulk-size"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_BULK_SIZE"},
		Category: "DATABASE SYNCER",
	}
	EventModeFlag = &cli.StringFlag{
		Name:     "dbsyncer.event.mode",
		Usage:    "The way how to sync all block or last block (block, head)",
		Value:    "head",
		Aliases:  []string{"db-syncer.event-mode"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_EVENT_MODE"},
		Category: "DATABASE SYNCER",
	}
	MaxBlockDiffFlag = &cli.Uint64Flag{
		Name:     "dbsyncer.max.block.diff",
		Usage:    "The maximum difference between current block and event block. 0 means off",
		Value:    0,
		Aliases:  []string{"db-syncer.max-block-diff"},
		EnvVars:  []string{"KLAYTN_DBSYNCER_MAX_BLOCK_DIFF"},
		Category: "DATABASE SYNCER",
	}
	AutoRestartFlag = &cli.BoolFlag{
		Name:     "autorestart.enable",
		Usage:    "Node can restart itself when there is a problem in making consensus",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_AUTORESTART_ENABLE"},
		Category: "MISC",
	}
	RestartTimeOutFlag = &cli.DurationFlag{
		Name:     "autorestart.timeout",
		Usage:    "The elapsed time to wait auto restart (minutes)",
		Value:    15 * time.Minute,
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_AUTORESTART_TIMEOUT"},
		Category: "MISC",
	}
	DaemonPathFlag = &cli.StringFlag{
		Name:     "autorestart.daemon.path",
		Usage:    "Path of node daemon. Used to give signal to kill",
		Value:    "~/klaytn/bin/kcnd",
		Aliases:  []string{"autorestart.daemon-path"},
		EnvVars:  []string{"KLAYTN_AUTORESTART_DAEMON_PATH"},
		Category: "MISC",
	}

	// db migration vars
	DstDbTypeFlag = &cli.StringFlag{
		Name:     "dst.dbtype",
		Usage:    `Blockchain storage database type ("LevelDB", "BadgerDB", "DynamoDBS3")`,
		Value:    "LevelDB",
		Aliases:  []string{"migration.dst.dbtype"},
		EnvVars:  []string{"KLAYTN_DST_DBTYPE"},
		Category: "DATABASE MIGRATION",
	}
	DstDataDirFlag = &cli.PathFlag{
		Name:     "dst.datadir",
		Usage:    "Data directory for the databases and keystore. This value is only used in local DB.",
		Aliases:  []string{"migration.dst.datadir"},
		EnvVars:  []string{"KLAYTN_DST_DATADIR"},
		Category: "DATABASE MIGRATION",
	}
	DstSingleDBFlag = &cli.BoolFlag{
		Name:     "db.dst.single",
		Usage:    "Create a single persistent storage. MiscDB, headerDB and etc are stored in one DB.",
		Aliases:  []string{"migration.dst.single"},
		EnvVars:  []string{"KLAYTN_DB_DST_SINGLE"},
		Category: "DATABASE MIGRATION",
	}
	DstLevelDBCacheSizeFlag = &cli.IntFlag{
		Name:     "db.dst.leveldb.cache-size",
		Usage:    "Size of in-memory cache in LevelDB (MiB)",
		Value:    768,
		Aliases:  []string{"migration.dst.db.leveldb.cache-size"},
		EnvVars:  []string{"KLAYTN_DB_DST_LEVELDB_CACHE_SIZE"},
		Category: "DATABASE MIGRATION",
	}
	DstLevelDBCompressionTypeFlag = &cli.IntFlag{
		Name:     "db.dst.leveldb.compression",
		Usage:    "Determines the compression method for LevelDB. 0=AllNoCompression, 1=ReceiptOnlySnappyCompression, 2=StateTrieOnlyNoCompression, 3=AllSnappyCompression",
		Value:    0,
		Aliases:  []string{"migration.dst.db.leveldb.compression"},
		EnvVars:  []string{"KLAYTN_DB_DST_LEVELDB_COMPRESSION"},
		Category: "DATABASE MIGRATION",
	}
	DstNumStateTrieShardsFlag = &cli.UintFlag{
		Name:     "db.dst.num-statetrie-shards",
		Usage:    "Number of internal shards of state trie DB shards. Should be power of 2",
		Value:    4,
		Aliases:  []string{"migration.dst.db.leveldb.num-statetrie-shards"},
		EnvVars:  []string{"KLAYTN_DB_DST_NUM_STATETRIE_SHARDS"},
		Category: "DATABASE MIGRATION",
	}
	DstDynamoDBTableNameFlag = &cli.StringFlag{
		Name:     "db.dst.dynamo.tablename",
		Usage:    "Specifies DynamoDB table name. This is mandatory to use dynamoDB. (Set dbtype to use DynamoDBS3). If dstDB is singleDB, tableName should be in form of 'PREFIX-TABLENAME'.(e.g. 'klaytn-misc', 'klaytn-statetrie')",
		Aliases:  []string{"migration.dst.db.dynamo.table-name"},
		EnvVars:  []string{"KLAYTN_DB_DST_DYNAMO_TABLENAME"},
		Category: "DATABASE MIGRATION",
	}
	DstDynamoDBRegionFlag = &cli.StringFlag{
		Name:     "db.dst.dynamo.region",
		Usage:    "AWS region where the DynamoDB will be created.",
		Value:    database.GetDefaultDynamoDBConfig().Region,
		Aliases:  []string{"migration.dst.db.dynamo.region"},
		EnvVars:  []string{"KLAYTN_DB_DST_DYNAMO_REGION"},
		Category: "DATABASE MIGRATION",
	}
	DstDynamoDBIsProvisionedFlag = &cli.BoolFlag{
		Name:     "db.dst.dynamo.is-provisioned",
		Usage:    "Set DynamoDB billing mode to provision. The default billing mode is on-demand.",
		Aliases:  []string{"migration.dst.db.dynamo.is-provisioned"},
		EnvVars:  []string{"KLAYTN_DB_DST_DYNAMO_IS_PROVISIONED"},
		Category: "DATABASE MIGRATION",
	}
	DstDynamoDBReadCapacityFlag = &cli.Int64Flag{
		Name:     "db.dst.dynamo.read-capacity",
		Usage:    "Read capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied.",
		Value:    database.GetDefaultDynamoDBConfig().ReadCapacityUnits,
		Aliases:  []string{"migration.dst.db.dynamo.read-capacity"},
		EnvVars:  []string{"KLAYTN_DB_DST_DYNAMO_READ_CAPACITY"},
		Category: "DATABASE MIGRATION",
	}
	DstDynamoDBWriteCapacityFlag = &cli.Int64Flag{
		Name:     "db.dst.dynamo.write-capacity",
		Usage:    "Write capacity unit of dynamoDB. If is-provisioned is not set, this flag will not be applied",
		Value:    database.GetDefaultDynamoDBConfig().WriteCapacityUnits,
		Aliases:  []string{"migration.dst.db.dynamo.write-capacity"},
		EnvVars:  []string{"KLAYTN_DB_DST_DYNAMO_WRITE_CAPACITY"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBSecondaryFlag = &cli.BoolFlag{
		Name:     "db.dst.rocksdb.secondary",
		Usage:    "Enable rocksdb secondary mode (read-only and catch-up with primary node dynamically)",
		Aliases:  []string{"migration.dst.db.rocksdb.secondary"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_SECONDARY"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBCacheSizeFlag = &cli.Uint64Flag{
		Name:     "db.dst.rocksdb.cache-size",
		Usage:    "Size of in-memory cache in RocksDB (MiB)",
		Value:    768,
		Aliases:  []string{"migration.dst.db.rocksdb.cache-size"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_CACHE_SIZE"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBDumpMallocStatFlag = &cli.BoolFlag{
		Name:     "db.dst.rocksdb.dump-memory-stat",
		Usage:    "Enable to print memory stat together with rocksdb.stat. Works with Jemalloc only.",
		Aliases:  []string{"migration.dst.db.rocksdb.dump-malloc-stat"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_DUMP_MALLOC_STAT"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBCompressionTypeFlag = &cli.StringFlag{
		Name:     "db.dst.rocksdb.compression-type",
		Usage:    "RocksDB block compression type. Supported values are 'no', 'snappy', 'zlib', 'bz', 'lz4', 'lz4hc', 'xpress', 'zstd'",
		Value:    database.GetDefaultRocksDBConfig().CompressionType,
		Aliases:  []string{"migration.dst.db.rocksdb.compression-type"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_COMPRESSION_TYPE"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBBottommostCompressionTypeFlag = &cli.StringFlag{
		Name:     "db.dst.rocksdb.bottommost-compression-type",
		Usage:    "RocksDB bottommost block compression type. Supported values are 'no', 'snappy', 'zlib', 'bz2', 'lz4', 'lz4hc', 'xpress', 'zstd'",
		Value:    database.GetDefaultRocksDBConfig().BottommostCompressionType,
		Aliases:  []string{"migration.dst.db.rocksdb.bottommost-compression-type"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_BOTTOMMOST_COMPRESSION_TYPE"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBFilterPolicyFlag = &cli.StringFlag{
		Name:     "db.dst.rocksdb.filter-policy",
		Usage:    "RocksDB filter policy. Supported values are 'no', 'bloom', 'ribbon'",
		Value:    database.GetDefaultRocksDBConfig().FilterPolicy,
		Aliases:  []string{"migration.dst.db.rocksdb.filter-policy"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_FILTER_POLICY"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBDisableMetricsFlag = &cli.BoolFlag{
		Name:     "db.dst.rocksdb.disable-metrics",
		Usage:    "Disable RocksDB metrics",
		Aliases:  []string{"migration.dst.db.rocksdb.disable-metrics"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_DISABLE_METRICS"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBMaxOpenFilesFlag = &cli.IntFlag{
		Name:     "db.dst.rocksdb.max-open-files",
		Usage:    "Set RocksDB max open files. (the value should be greater than 16)",
		Value:    database.GetDefaultRocksDBConfig().MaxOpenFiles,
		Aliases:  []string{"migration.dst.db.rocksdb.max-open-files"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_MAX_OPEN_FILES"},
		Category: "DATABASE MIGRATION",
	}
	DstRocksDBCacheIndexAndFilterFlag = &cli.BoolFlag{
		Name:     "db.dst.rocksdb.cache-index-and-filter",
		Usage:    "Use block cache for index and filter blocks.",
		Aliases:  []string{"migration.dst.db.rocksdb.cache-index-and-filter"},
		EnvVars:  []string{"KLAYTN_DB_DST_ROCKSDB_CACHE_INDEX_AND_FILTER"},
		Category: "DATABASE MIGRATION",
	}

	// Config
	ConfigFileFlag = &cli.StringFlag{
		Name:     "config",
		Usage:    "TOML configuration file",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_CONFIG"},
		Category: "KLAY",
	}
	BlockGenerationIntervalFlag = &cli.Int64Flag{
		Name: "block-generation-interval",
		Usage: "(experimental option) Set the block generation interval in seconds. " +
			"It should be equal or larger than 1. This flag is only applicable to CN.",
		Value:    params.DefaultBlockGenerationInterval,
		Aliases:  []string{"experimental.block-generation-interval"},
		EnvVars:  []string{"KLAYTN_BLOCK_GENERATION_INTERVAL"},
		Category: "KLAY",
	}
	BlockGenerationTimeLimitFlag = &cli.DurationFlag{
		Name: "block-generation-time-limit",
		Usage: "(experimental option) Set the vm execution time limit during block generation. " +
			"Less than half of the block generation interval is recommended for this value. " +
			"This flag is only applicable to CN",
		Value:    params.DefaultBlockGenerationTimeLimit,
		Aliases:  []string{"experimental.block-generation-time-limit"},
		EnvVars:  []string{"KLAYTN_BLOCK_GENERATION_TIME_LIMIT"},
		Category: "KLAY",
	}
	OpcodeComputationCostLimitFlag = &cli.Uint64Flag{
		Name: "opcode-computation-cost-limit",
		Usage: "(experimental option) Set the computation cost limit for a tx. " +
			"Should set the same value within the network",
		Value:    params.DefaultOpcodeComputationCostLimit,
		Aliases:  []string{"experimental.opcode-computation-cost-limit"},
		EnvVars:  []string{"KLAYTN_OPCODE_COMPUTATION_COST_LIMIT"},
		Category: "KLAY",
	}

	// TODO-Klaytn-Bootnode: Add bootnode's metric options
	// TODO-Klaytn-Bootnode: Implements bootnode's RPC
)

// MakeDataDir retrieves the currently requested data directory, terminating
// if none (or the empty string) is specified. If the node is starting a baobab,
// the a subdirectory of the specified datadir will be used.
func MakeDataDir(ctx *cli.Context) string {
	if path := ctx.String(DataDirFlag.Name); path != "" {
		if ctx.Bool(BaobabFlag.Name) {
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
	path := ctx.String(PasswordFileFlag.Name)
	if path == "" {
		return nil
	}
	text, err := os.ReadFile(path)
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
	if ctx.String(PreloadJSFlag.Name) == "" {
		return nil
	}
	// Otherwise resolve absolute paths and return them
	var preloads []string

	assets := ctx.String(JSpathFlag.Name)
	for _, file := range strings.Split(ctx.String(PreloadJSFlag.Name), ",") {
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
// Deprecated: urfave/cli/v2 doesn't support local scope flag and there is only global scope.
func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.Set(name, ctx.String(name))
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
		name := flag.Names()[0]

		if i+1 < len(args) {
			switch option := args[i+1].(type) {
			case string:
				// Extended flag, expand the name and shift the arguments
				if ctx.String(flag.Names()[0]) == option {
					name += "=" + option
				}
				i++

			case cli.Flag:
			default:
				panic(fmt.Sprintf("invalid argument, not cli.Flag or string extension: %T", args[i+1]))
			}
		}
		// Mark the flag if it's set
		if ctx.IsSet(flag.Names()[0]) {
			set = append(set, "--"+name)
		}
	}
	if len(set) > 1 {
		log.Fatalf("Flags %v can't be used at the same time", strings.Join(set, ", "))
	}
}

// FlagString prints a single flag in help.
func FlagString(f cli.Flag) string {
	df, ok := f.(cli.DocGenerationFlag)
	if !ok {
		return ""
	}

	needsPlaceholder := df.TakesValue()
	placeholder := ""
	if needsPlaceholder {
		placeholder = "value"
	}

	// namesText := cli.FlagNamePrefixer([]string{df.Names()[0]}, placeholder)
	// if len(df.Names()) > 1 {
	// 	namesText = cli.FlagNamePrefixer([]string{df.Names()[1]}, placeholder)
	// }
	namesText := cli.FlagNamePrefixer(df.Names(), placeholder)

	defaultValueString := ""
	if s := df.GetDefaultText(); s != "" {
		defaultValueString = " (default: " + s + ")"
	}

	usage := strings.TrimSpace(df.GetUsage())
	envHint := strings.TrimSpace(cli.FlagEnvHinter(df.GetEnvVars(), ""))
	if len(envHint) > 0 {
		usage += "\nEnvVar: " + envHint
	}

	// usage = wordWrap(usage, 150)
	usage = indent(usage, 10)
	return fmt.Sprintf("\n    %s%s\n%s", namesText, defaultValueString, usage)
}

func indent(s string, nspace int) string {
	ind := strings.Repeat(" ", nspace)
	return ind + strings.ReplaceAll(s, "\n", "\n"+ind)
}
