// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	logger = log.NewModuleLogger(log.StorageDatabase)

	errGovIdxAlreadyExist = errors.New("a governance idx of the more recent or the same block exist")

	HeadBlockQ backupHashQueue
	FastBlockQ backupHashQueue
)

type DBManager interface {
	IsParallelDBWrite() bool
	IsSingle() bool
	InMigration() bool
	MigrationBlockNumber() uint64
	getStateTrieMigrationInfo() uint64

	Close()
	NewBatch(dbType DBEntryType) Batch
	getDBDir(dbEntry DBEntryType) string
	setDBDir(dbEntry DBEntryType, newDBDir string)
	setStateTrieMigrationStatus(uint64)
	GetMemDB() *MemDB
	GetDBConfig() *DBConfig
	getDatabase(DBEntryType) Database
	CreateMigrationDBAndSetStatus(blockNum uint64) error
	FinishStateMigration(succeed bool) chan struct{}
	GetStateTrieDB() Database
	GetStateTrieMigrationDB() Database
	GetMiscDB() Database
	GetSnapshotDB() Database
	GetProperty(dt DBEntryType, name string) string

	// from accessors_chain.go
	ReadCanonicalHash(number uint64) common.Hash
	WriteCanonicalHash(hash common.Hash, number uint64)
	DeleteCanonicalHash(number uint64)

	ReadHeadHeaderHash() common.Hash
	WriteHeadHeaderHash(hash common.Hash)

	ReadHeadBlockHash() common.Hash
	ReadHeadBlockBackupHash() common.Hash
	WriteHeadBlockHash(hash common.Hash)

	ReadHeadFastBlockHash() common.Hash
	ReadHeadFastBlockBackupHash() common.Hash
	WriteHeadFastBlockHash(hash common.Hash)

	ReadFastTrieProgress() uint64
	WriteFastTrieProgress(count uint64)

	HasHeader(hash common.Hash, number uint64) bool
	ReadHeader(hash common.Hash, number uint64) *types.Header
	ReadHeaderRLP(hash common.Hash, number uint64) rlp.RawValue
	WriteHeader(header *types.Header)
	DeleteHeader(hash common.Hash, number uint64)
	ReadHeaderNumber(hash common.Hash) *uint64

	HasBody(hash common.Hash, number uint64) bool
	ReadBody(hash common.Hash, number uint64) *types.Body
	ReadBodyInCache(hash common.Hash) *types.Body
	ReadBodyRLP(hash common.Hash, number uint64) rlp.RawValue
	ReadBodyRLPByHash(hash common.Hash) rlp.RawValue
	WriteBody(hash common.Hash, number uint64, body *types.Body)
	PutBodyToBatch(batch Batch, hash common.Hash, number uint64, body *types.Body)
	WriteBodyRLP(hash common.Hash, number uint64, rlp rlp.RawValue)
	DeleteBody(hash common.Hash, number uint64)

	ReadTd(hash common.Hash, number uint64) *big.Int
	WriteTd(hash common.Hash, number uint64, td *big.Int)
	DeleteTd(hash common.Hash, number uint64)

	ReadReceipt(txHash common.Hash) (*types.Receipt, common.Hash, uint64, uint64)
	ReadReceipts(blockHash common.Hash, number uint64) types.Receipts
	ReadReceiptsByBlockHash(hash common.Hash) types.Receipts
	WriteReceipts(hash common.Hash, number uint64, receipts types.Receipts)
	PutReceiptsToBatch(batch Batch, hash common.Hash, number uint64, receipts types.Receipts)
	DeleteReceipts(hash common.Hash, number uint64)

	ReadBlock(hash common.Hash, number uint64) *types.Block
	ReadBlockByHash(hash common.Hash) *types.Block
	ReadBlockByNumber(number uint64) *types.Block
	HasBlock(hash common.Hash, number uint64) bool
	WriteBlock(block *types.Block)
	DeleteBlock(hash common.Hash, number uint64)

	ReadBadBlock(hash common.Hash) *types.Block
	WriteBadBlock(block *types.Block)
	ReadAllBadBlocks() ([]*types.Block, error)
	DeleteBadBlocks()

	FindCommonAncestor(a, b *types.Header) *types.Header

	ReadIstanbulSnapshot(hash common.Hash) ([]byte, error)
	WriteIstanbulSnapshot(hash common.Hash, blob []byte) error
	DeleteIstanbulSnapshot(hash common.Hash)

	WriteMerkleProof(key, value []byte)

	// Bytecodes related operations
	ReadCode(hash common.Hash) []byte
	ReadCodeWithPrefix(hash common.Hash) []byte
	WriteCode(hash common.Hash, code []byte)
	PutCodeToBatch(batch Batch, hash common.Hash, code []byte)
	DeleteCode(hash common.Hash)
	HasCode(hash common.Hash) bool

	// State Trie Database related operations
	ReadTrieNode(hash common.ExtHash) ([]byte, error)
	HasTrieNode(hash common.ExtHash) (bool, error)
	HasCodeWithPrefix(hash common.Hash) bool
	ReadPreimage(hash common.Hash) []byte

	// Read StateTrie from new DB
	ReadTrieNodeFromNew(hash common.ExtHash) ([]byte, error)
	HasTrieNodeFromNew(hash common.ExtHash) (bool, error)
	HasCodeWithPrefixFromNew(hash common.Hash) bool
	ReadPreimageFromNew(hash common.Hash) []byte

	// Read StateTrie from old DB
	ReadTrieNodeFromOld(hash common.ExtHash) ([]byte, error)
	HasTrieNodeFromOld(hash common.ExtHash) (bool, error)
	HasCodeWithPrefixFromOld(hash common.Hash) bool
	ReadPreimageFromOld(hash common.Hash) []byte

	// Write StateTrie
	WriteTrieNode(hash common.ExtHash, node []byte)
	PutTrieNodeToBatch(batch Batch, hash common.ExtHash, node []byte)
	DeleteTrieNode(hash common.ExtHash)
	WritePreimages(number uint64, preimages map[common.Hash][]byte)

	// Trie pruning
	ReadPruningEnabled() bool
	WritePruningEnabled()
	DeletePruningEnabled()

	WritePruningMarks(marks []PruningMark)
	ReadPruningMarks(startNumber, endNumber uint64) []PruningMark
	DeletePruningMarks(marks []PruningMark)
	PruneTrieNodes(marks []PruningMark)
	WriteLastPrunedBlockNumber(blockNumber uint64)
	ReadLastPrunedBlockNumber() (uint64, error)

	// from accessors_indexes.go
	ReadTxLookupEntry(hash common.Hash) (common.Hash, uint64, uint64)
	WriteTxLookupEntries(block *types.Block)
	WriteAndCacheTxLookupEntries(block *types.Block) error
	PutTxLookupEntriesToBatch(batch Batch, block *types.Block)
	DeleteTxLookupEntry(hash common.Hash)

	ReadTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)

	NewSenderTxHashToTxHashBatch() Batch
	PutSenderTxHashToTxHashToBatch(batch Batch, senderTxHash, txHash common.Hash) error
	ReadTxHashFromSenderTxHash(senderTxHash common.Hash) common.Hash

	ReadBloomBits(bloomBitsKey []byte) ([]byte, error)
	WriteBloomBits(bloomBitsKey []byte, bits []byte) error

	ReadValidSections() ([]byte, error)
	WriteValidSections(encodedSections []byte)

	ReadSectionHead(encodedSection []byte) ([]byte, error)
	WriteSectionHead(encodedSection []byte, hash common.Hash)
	DeleteSectionHead(encodedSection []byte)

	// from accessors_metadata.go
	ReadDatabaseVersion() *uint64
	WriteDatabaseVersion(version uint64)

	ReadChainConfig(hash common.Hash) *params.ChainConfig
	WriteChainConfig(hash common.Hash, cfg *params.ChainConfig)

	// from accessors_snapshot.go
	ReadSnapshotJournal() []byte
	WriteSnapshotJournal(journal []byte)
	DeleteSnapshotJournal()

	ReadSnapshotGenerator() []byte
	WriteSnapshotGenerator(generator []byte)
	DeleteSnapshotGenerator()

	ReadSnapshotDisabled() bool
	WriteSnapshotDisabled()
	DeleteSnapshotDisabled()

	ReadSnapshotRecoveryNumber() *uint64
	WriteSnapshotRecoveryNumber(number uint64)
	DeleteSnapshotRecoveryNumber()

	ReadSnapshotSyncStatus() []byte
	WriteSnapshotSyncStatus(status []byte)
	DeleteSnapshotSyncStatus()

	ReadSnapshotRoot() common.Hash
	WriteSnapshotRoot(root common.Hash)
	DeleteSnapshotRoot()

	ReadAccountSnapshot(hash common.Hash) []byte
	WriteAccountSnapshot(hash common.Hash, entry []byte)
	DeleteAccountSnapshot(hash common.Hash)

	ReadStorageSnapshot(accountHash, storageHash common.Hash) []byte
	WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte)
	DeleteStorageSnapshot(accountHash, storageHash common.Hash)

	NewSnapshotDBIterator(prefix []byte, start []byte) Iterator

	NewSnapshotDBBatch() SnapshotDBBatch

	// below operations are used in parent chain side, not child chain side.
	WriteChildChainTxHash(ccBlockHash common.Hash, ccTxHash common.Hash)
	ConvertChildChainBlockHashToParentChainTxHash(scBlockHash common.Hash) common.Hash

	WriteLastIndexedBlockNumber(blockNum uint64)
	GetLastIndexedBlockNumber() uint64

	// below operations are used in child chain side, not parent chain side.
	WriteAnchoredBlockNumber(blockNum uint64)
	ReadAnchoredBlockNumber() uint64

	WriteReceiptFromParentChain(blockHash common.Hash, receipt *types.Receipt)
	ReadReceiptFromParentChain(blockHash common.Hash) *types.Receipt

	WriteHandleTxHashFromRequestTxHash(rTx, hTx common.Hash)
	ReadHandleTxHashFromRequestTxHash(rTx common.Hash) common.Hash

	WriteParentOperatorFeePayer(feePayer common.Address)
	WriteChildOperatorFeePayer(feePayer common.Address)
	ReadParentOperatorFeePayer() common.Address
	ReadChildOperatorFeePayer() common.Address

	// cacheManager related functions.
	ClearHeaderChainCache()
	ClearBlockChainCache()
	ReadTxAndLookupInfoInCache(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	ReadBlockReceiptsInCache(blockHash common.Hash) types.Receipts
	ReadTxReceiptInCache(txHash common.Hash) *types.Receipt

	// snapshot in clique(ConsensusClique) consensus
	WriteCliqueSnapshot(snapshotBlockHash common.Hash, encodedSnapshot []byte) error
	ReadCliqueSnapshot(snapshotBlockHash common.Hash) ([]byte, error)

	// Governance related functions
	WriteGovernance(data map[string]interface{}, num uint64) error
	WriteGovernanceIdx(num uint64) error
	ReadGovernance(num uint64) (map[string]interface{}, error)
	ReadRecentGovernanceIdx(count int) ([]uint64, error)
	ReadGovernanceAtNumber(num uint64, epoch uint64) (uint64, map[string]interface{}, error)
	WriteGovernanceState(b []byte) error
	ReadGovernanceState() ([]byte, error)
	DeleteGovernance(num uint64)
	// TODO-Klaytn implement governance DB deletion methods.

	// StakingInfo related functions
	ReadStakingInfo(blockNum uint64) ([]byte, error)
	WriteStakingInfo(blockNum uint64, stakingInfo []byte) error
	HasStakingInfo(blockNum uint64) (bool, error)
	DeleteStakingInfo(blockNum uint64)

	// DB migration related function
	StartDBMigration(DBManager) error

	// ChainDataFetcher checkpoint function
	WriteChainDataFetcherCheckpoint(checkpoint uint64) error
	ReadChainDataFetcherCheckpoint() (uint64, error)

	TryCatchUpWithPrimary() error
}

type DBEntryType uint8

const (
	MiscDB DBEntryType = iota // Do not move MiscDB which has the path of others DB.
	headerDB
	BodyDB
	ReceiptsDB
	StateTrieDB
	StateTrieMigrationDB
	TxLookUpEntryDB
	bridgeServiceDB
	SnapshotDB
	// databaseEntryTypeSize should be the last item in this list!!
	databaseEntryTypeSize
)

type backupHashQueue struct {
	backupHashes [backupHashCnt]common.Hash
	idx          int
}

func (b *backupHashQueue) push(h common.Hash) {
	b.backupHashes[b.idx%backupHashCnt] = h
	b.idx = (b.idx + 1) % backupHashCnt
}

func (b *backupHashQueue) pop() common.Hash {
	if b.backupHashes[b.idx] == (common.Hash{}) {
		return common.Hash{}
	}
	return b.backupHashes[b.idx]
}

func (et DBEntryType) String() string {
	return dbBaseDirs[et]
}

const (
	notInMigrationFlag = 0
	inMigrationFlag    = 1
	backupHashCnt      = 128
)

var dbBaseDirs = [databaseEntryTypeSize]string{
	"misc", // do not move misc
	"header",
	"body",
	"receipts",
	"statetrie",
	"statetrie_migrated", // "statetrie_migrated_#N" path will be used. (#N is a migrated block number.)
	"txlookup",
	"bridgeservice",
	"snapshot",
}

// Sum of dbConfigRatio should be 100.
// Otherwise, logger.Crit will be called at checkDBEntryConfigRatio.
var dbConfigRatio = [databaseEntryTypeSize]int{
	2,  // MiscDB
	5,  // headerDB
	5,  // BodyDB
	5,  // ReceiptsDB
	40, // StateTrieDB
	37, // StateTrieMigrationDB
	2,  // TXLookUpEntryDB
	1,  // bridgeServiceDB
	3,  // SnapshotDB
}

// checkDBEntryConfigRatio checks if sum of dbConfigRatio is 100.
// If it isn't, logger.Crit is called.
func checkDBEntryConfigRatio() {
	entryConfigRatioSum := 0
	for i := 0; i < int(databaseEntryTypeSize); i++ {
		entryConfigRatioSum += dbConfigRatio[i]
	}
	if entryConfigRatioSum != 100 {
		logger.Crit("Sum of dbConfigRatio elements should be 100", "actual", entryConfigRatioSum)
	}
}

// getDBEntryConfig returns a new DBConfig with original DBConfig, DBEntryType and dbDir.
// It adjusts configuration according to the ratio specified in dbConfigRatio and dbDirs.
func getDBEntryConfig(originalDBC *DBConfig, i DBEntryType, dbDir string) *DBConfig {
	newDBC := *originalDBC
	ratio := dbConfigRatio[i]

	newDBC.LevelDBCacheSize = originalDBC.LevelDBCacheSize * ratio / 100
	newDBC.OpenFilesLimit = originalDBC.OpenFilesLimit * ratio / 100

	// Update dir to each Database specific directory.
	newDBC.Dir = filepath.Join(originalDBC.Dir, dbDir)
	// Update dynmao table name to Database specific name.
	if newDBC.DynamoDBConfig != nil {
		newDynamoDBConfig := *originalDBC.DynamoDBConfig
		newDynamoDBConfig.TableName += "-" + dbDir
		newDBC.DynamoDBConfig = &newDynamoDBConfig
	}

	if newDBC.RocksDBConfig != nil {
		newRocksDBConfig := *originalDBC.RocksDBConfig
		newRocksDBConfig.CacheSize = originalDBC.RocksDBConfig.CacheSize * uint64(ratio) / 100
		newRocksDBConfig.MaxOpenFiles = originalDBC.RocksDBConfig.MaxOpenFiles * ratio / 100
		newDBC.RocksDBConfig = &newRocksDBConfig
	}

	return &newDBC
}

type databaseManager struct {
	config *DBConfig
	dbs    []Database
	cm     *cacheManager

	// TODO-Klaytn need to refine below.
	// -merge status variable
	lockInMigration      sync.RWMutex
	inMigration          bool
	migrationBlockNumber uint64
}

func NewMemoryDBManager() DBManager {
	dbc := &DBConfig{DBType: MemoryDB}

	dbm := databaseManager{
		config: dbc,
		dbs:    make([]Database, 1, 1),
		cm:     newCacheManager(),
	}
	dbm.dbs[0] = NewMemDB()

	return &dbm
}

// DBConfig handles database related configurations.
type DBConfig struct {
	// General configurations for all types of DB.
	Dir                 string
	DBType              DBType
	SingleDB            bool // whether dbs (such as MiscDB, headerDB and etc) share one physical DB
	NumStateTrieShards  uint // the number of shards of state trie db
	ParallelDBWrite     bool
	OpenFilesLimit      int
	EnableDBPerfMetrics bool // If true, read and write performance will be logged

	// LevelDB related configurations.
	LevelDBCacheSize   int // LevelDBCacheSize = BlockCacheCapacity + WriteBuffer
	LevelDBCompression LevelDBCompressionType
	LevelDBBufferPool  bool

	// RocksDB related configurations
	RocksDBConfig *RocksDBConfig

	// DynamoDB related configurations
	DynamoDBConfig *DynamoDBConfig
}

const dbMetricPrefix = "klay/db/chaindata/"

// singleDatabaseDBManager returns DBManager which handles one single Database.
// Each Database will share one common Database.
func singleDatabaseDBManager(dbc *DBConfig) (DBManager, error) {
	dbm := newDatabaseManager(dbc)
	db, err := newDatabase(dbc, 0)
	if err != nil {
		return nil, err
	}

	db.Meter(dbMetricPrefix)
	for i := 0; i < int(databaseEntryTypeSize); i++ {
		dbm.dbs[i] = db
	}
	return dbm, nil
}

// newMiscDB returns misc DBManager. If not exist, the function create DB before returning.
func newMiscDB(dbc *DBConfig) Database {
	newDBC := getDBEntryConfig(dbc, MiscDB, dbBaseDirs[MiscDB])
	db, err := newDatabase(newDBC, MiscDB)
	if err != nil {
		logger.Crit("Failed while generating a MISC database", "err", err)
	}

	db.Meter(dbMetricPrefix + dbBaseDirs[MiscDB] + "/")
	return db
}

// databaseDBManager returns DBManager which handles Databases.
// Each Database will have its own separated Database.
func databaseDBManager(dbc *DBConfig) (*databaseManager, error) {
	dbm := newDatabaseManager(dbc)
	var db Database
	var err error

	// Create Misc DB first to get the DB directory of stateTrieDB.
	miscDB := newMiscDB(dbc)
	dbm.dbs[MiscDB] = miscDB

	// Create other DBs
	for et := int(MiscDB) + 1; et < int(databaseEntryTypeSize); et++ {
		entryType := DBEntryType(et)
		dir := dbm.getDBDir(entryType)

		switch entryType {
		case StateTrieMigrationDB:
			if dir == dbBaseDirs[StateTrieMigrationDB] {
				// If there is no migration DB, skip to set.
				continue
			}
			fallthrough
		case StateTrieDB:
			newDBC := getDBEntryConfig(dbc, entryType, dir)
			if dbc.NumStateTrieShards > 1 && !dbc.DBType.selfShardable() { // make non-sharding db if the db is sharding itself
				db, err = newShardedDB(newDBC, entryType, dbc.NumStateTrieShards)
			} else {
				db, err = newDatabase(newDBC, entryType)
			}
		default:
			newDBC := getDBEntryConfig(dbc, entryType, dir)
			db, err = newDatabase(newDBC, entryType)
		}

		if err != nil {
			logger.Crit("Failed while generating databases", "DBType", dbBaseDirs[et], "err", err)
		}

		dbm.dbs[et] = db
		db.Meter(dbMetricPrefix + dbBaseDirs[et] + "/") // Each database collects metrics independently.
	}
	return dbm, nil
}

// newDatabase returns Database interface with given DBConfig.
func newDatabase(dbc *DBConfig, entryType DBEntryType) (Database, error) {
	switch dbc.DBType {
	case LevelDB:
		return NewLevelDB(dbc, entryType)
	case RocksDB:
		return NewRocksDB(dbc.Dir, dbc.RocksDBConfig)
	case BadgerDB:
		return NewBadgerDB(dbc.Dir)
	case MemoryDB:
		return NewMemDB(), nil
	case DynamoDB:
		return NewDynamoDB(dbc.DynamoDBConfig)
	default:
		logger.Info("database type is not set, fall back to default LevelDB")
		return NewLevelDB(dbc, 0)
	}
}

// newDatabaseManager returns the pointer of databaseManager with default configuration.
func newDatabaseManager(dbc *DBConfig) *databaseManager {
	return &databaseManager{
		config: dbc,
		dbs:    make([]Database, databaseEntryTypeSize),
		cm:     newCacheManager(),
	}
}

// NewDBManager returns DBManager interface.
// If SingleDB is false, each Database will have its own DB.
// If not, each Database will share one common DB.
func NewDBManager(dbc *DBConfig) DBManager {
	if dbc.SingleDB {
		logger.Info("Single database is used for persistent storage", "DBType", dbc.DBType)
		if dbm, err := singleDatabaseDBManager(dbc); err != nil {
			logger.Crit("Failed to create a single database", "DBType", dbc.DBType, "err", err)
		} else {
			return dbm
		}
	} else {
		checkDBEntryConfigRatio()
		logger.Info("Non-single database is used for persistent storage", "DBType", dbc.DBType)
		dbm, err := databaseDBManager(dbc)
		if err != nil {
			logger.Crit("Failed to create databases", "DBType", dbc.DBType, "err", err)
		}
		if migrationBlockNum := dbm.getStateTrieMigrationInfo(); migrationBlockNum > 0 {
			mdb := dbm.getDatabase(StateTrieMigrationDB)
			if mdb == nil {
				logger.Error("Failed to load StateTrieMigrationDB database", "migrationBlockNumber", migrationBlockNum)
			} else {
				dbm.inMigration = true
				dbm.migrationBlockNumber = migrationBlockNum
			}
		}
		return dbm
	}
	logger.Crit("Must not reach here!")
	return nil
}

func (dbm *databaseManager) IsParallelDBWrite() bool {
	return dbm.config.ParallelDBWrite
}

func (dbm *databaseManager) IsSingle() bool {
	return dbm.config.SingleDB
}

func (dbm *databaseManager) InMigration() bool {
	dbm.lockInMigration.RLock()
	defer dbm.lockInMigration.RUnlock()

	return dbm.inMigration
}

func (dbm *databaseManager) MigrationBlockNumber() uint64 {
	return dbm.migrationBlockNumber
}

func (dbm *databaseManager) NewBatch(dbEntryType DBEntryType) Batch {
	if dbEntryType == StateTrieDB {
		dbm.lockInMigration.RLock()
		defer dbm.lockInMigration.RUnlock()

		if dbm.inMigration {
			newDBBatch := dbm.getDatabase(StateTrieMigrationDB).NewBatch()
			oldDBBatch := dbm.getDatabase(StateTrieDB).NewBatch()
			return NewStateTrieDBBatch([]Batch{oldDBBatch, newDBBatch})
		}
	} else if dbEntryType == StateTrieMigrationDB {
		return dbm.GetStateTrieMigrationDB().NewBatch()
	}
	return dbm.getDatabase(dbEntryType).NewBatch()
}

func NewStateTrieDBBatch(batches []Batch) Batch {
	return &stateTrieDBBatch{batches: batches}
}

type stateTrieDBBatch struct {
	batches []Batch
}

func (stdBatch *stateTrieDBBatch) Put(key []byte, value []byte) error {
	var errResult error
	for _, batch := range stdBatch.batches {
		if err := batch.Put(key, value); err != nil {
			errResult = err
		}
	}
	return errResult
}

func (stdBatch *stateTrieDBBatch) Delete(key []byte) error {
	var errResult error
	for _, batch := range stdBatch.batches {
		if err := batch.Delete(key); err != nil {
			errResult = err
		}
	}
	return errResult
}

// ValueSize is called to determine whether to write batches when it exceeds
// certain limit. stdBatch returns the largest size of its batches to
// write all batches at once when one of batch exceeds the limit.
func (stdBatch *stateTrieDBBatch) ValueSize() int {
	maxSize := 0
	for _, batch := range stdBatch.batches {
		if batch.ValueSize() > maxSize {
			maxSize = batch.ValueSize()
		}
	}

	return maxSize
}

// Write passes the list of batch to WriteBatchesParallel for writing batches.
func (stdBatch *stateTrieDBBatch) Write() error {
	_, err := WriteBatchesParallel(stdBatch.batches...)
	return err
}

func (stdBatch *stateTrieDBBatch) Reset() {
	for _, batch := range stdBatch.batches {
		batch.Reset()
	}
}

func (stdBatch *stateTrieDBBatch) Release() {
	for _, batch := range stdBatch.batches {
		batch.Release()
	}
}

func (stdBatch *stateTrieDBBatch) Replay(w KeyValueWriter) error {
	var errResult error
	for _, batch := range stdBatch.batches {
		if err := batch.Replay(w); err != nil {
			errResult = err
		}
	}
	return errResult
}

func (dbm *databaseManager) getDBDir(dbEntry DBEntryType) string {
	miscDB := dbm.getDatabase(MiscDB)

	enc, _ := miscDB.Get(databaseDirKey(uint64(dbEntry)))
	if len(enc) == 0 {
		return dbBaseDirs[dbEntry]
	}
	return string(enc)
}

func (dbm *databaseManager) setDBDir(dbEntry DBEntryType, newDBDir string) {
	miscDB := dbm.getDatabase(MiscDB)
	if err := miscDB.Put(databaseDirKey(uint64(dbEntry)), []byte(newDBDir)); err != nil {
		logger.Crit("Failed to put DB dir", "err", err)
	}
}

func (dbm *databaseManager) getStateTrieMigrationInfo() uint64 {
	miscDB := dbm.getDatabase(MiscDB)

	enc, _ := miscDB.Get(migrationStatusKey)
	if len(enc) != 8 {
		return 0
	}

	blockNum := binary.BigEndian.Uint64(enc)
	return blockNum
}

func (dbm *databaseManager) setStateTrieMigrationStatus(blockNum uint64) {
	miscDB := dbm.getDatabase(MiscDB)
	if err := miscDB.Put(migrationStatusKey, common.Int64ToByteBigEndian(blockNum)); err != nil {
		logger.Crit("Failed to set state trie migration status", "err", err)
	}

	if blockNum == 0 {
		dbm.inMigration = false
		return
	}

	dbm.inMigration, dbm.migrationBlockNumber = true, blockNum
}

func newStateTrieMigrationDB(dbc *DBConfig, blockNum uint64) (Database, string) {
	dbDir := dbBaseDirs[StateTrieMigrationDB] + "_" + strconv.FormatUint(blockNum, 10)
	newDBConfig := getDBEntryConfig(dbc, StateTrieMigrationDB, dbDir)
	var newDB Database
	var err error
	if newDBConfig.NumStateTrieShards > 1 {
		newDB, err = newShardedDB(newDBConfig, StateTrieMigrationDB, newDBConfig.NumStateTrieShards)
	} else {
		newDB, err = newDatabase(newDBConfig, StateTrieMigrationDB)
	}
	if err != nil {
		logger.Crit("Failed to create a new database for state trie migration", "err", err)
	}

	newDB.Meter(dbMetricPrefix + dbBaseDirs[StateTrieMigrationDB] + "/") // Each database collects metrics independently.
	logger.Info("Created a new database for state trie migration", "newStateTrieDB", newDBConfig.Dir)

	return newDB, dbDir
}

// CreateMigrationDBAndSetStatus create migrationDB and set migration status.
func (dbm *databaseManager) CreateMigrationDBAndSetStatus(blockNum uint64) error {
	if dbm.InMigration() {
		logger.Warn("Failed to set a new state trie migration db. Already in migration")
		return errors.New("already in migration")
	}
	if dbm.config.SingleDB {
		logger.Warn("Setting a new database for state trie migration is allowed for non-single database only")
		return errors.New("singleDB does not support state trie migration")
	}

	logger.Info("Start setting a new database for state trie migration", "blockNum", blockNum)

	// Create a new database for migration process.
	newDB, newDBDir := newStateTrieMigrationDB(dbm.config, blockNum)

	// lock to prevent from a conflict of reading state DB and changing state DB
	dbm.lockInMigration.Lock()
	defer dbm.lockInMigration.Unlock()

	// Store migration db path in misc db
	dbm.setDBDir(StateTrieMigrationDB, newDBDir)

	// Set migration db
	dbm.dbs[StateTrieMigrationDB] = newDB

	// Store the migration status
	dbm.setStateTrieMigrationStatus(blockNum)

	return nil
}

// FinishStateMigration updates stateTrieDB and removes the old one.
// The function should be called only after when state trie migration is finished.
// It returns a channel that closes when removeDB is finished.
func (dbm *databaseManager) FinishStateMigration(succeed bool) chan struct{} {
	// lock to prevent from a conflict of reading state DB and changing state DB
	dbm.lockInMigration.Lock()
	defer dbm.lockInMigration.Unlock()

	dbRemoved := StateTrieDB
	dbUsed := StateTrieMigrationDB

	if !succeed {
		dbRemoved, dbUsed = dbUsed, dbRemoved
	}

	dbToBeRemoved := dbm.dbs[dbRemoved]
	dbToBeUsed := dbm.dbs[dbUsed]
	dbDirToBeRemoved := dbm.getDBDir(dbRemoved)
	dbDirToBeUsed := dbm.getDBDir(dbUsed)

	// Replace StateTrieDB with new one
	dbm.setDBDir(StateTrieDB, dbDirToBeUsed)
	dbm.dbs[StateTrieDB] = dbToBeUsed

	dbm.setStateTrieMigrationStatus(0)

	dbm.dbs[StateTrieMigrationDB] = nil
	dbm.setDBDir(StateTrieMigrationDB, "")

	dbPathToBeRemoved := filepath.Join(dbm.config.Dir, dbDirToBeRemoved)
	dbToBeRemoved.Close()

	endCheck := make(chan struct{})
	go removeDB(dbPathToBeRemoved, endCheck)
	return endCheck
}

func removeDB(dbPath string, endCheck chan struct{}) {
	defer func() {
		if endCheck != nil {
			close(endCheck)
		}
	}()
	if err := os.RemoveAll(dbPath); err != nil {
		logger.Error("Failed to remove the database due to an error", "err", err, "dir", dbPath)
		return
	}
	logger.Info("Successfully removed database", "path", dbPath)
}

func (dbm *databaseManager) GetStateTrieDB() Database {
	return dbm.dbs[StateTrieDB]
}

func (dbm *databaseManager) GetStateTrieMigrationDB() Database {
	return dbm.dbs[StateTrieMigrationDB]
}

func (dbm *databaseManager) GetMiscDB() Database {
	return dbm.dbs[MiscDB]
}

func (dbm *databaseManager) GetSnapshotDB() Database {
	return dbm.getDatabase(SnapshotDB)
}

func (dbm *databaseManager) GetProperty(dt DBEntryType, name string) string {
	return dbm.getDatabase(dt).GetProperty(name)
}

func (dbm *databaseManager) TryCatchUpWithPrimary() error {
	for _, db := range dbm.dbs {
		if db != nil {
			if err := db.TryCatchUpWithPrimary(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (dbm *databaseManager) GetMemDB() *MemDB {
	if dbm.config.DBType == MemoryDB {
		if memDB, ok := dbm.dbs[0].(*MemDB); ok {
			return memDB
		} else {
			logger.Error("DBManager is set as memory DBManager, but actual value is not set as memory DBManager.")
			return nil
		}
	}
	logger.Error("GetMemDB() call to non memory DBManager object.")
	return nil
}

// GetDBConfig returns DBConfig of the DB manager.
func (dbm *databaseManager) GetDBConfig() *DBConfig {
	return dbm.config
}

func (dbm *databaseManager) getDatabase(dbEntryType DBEntryType) Database {
	if dbm.config.DBType == MemoryDB {
		return dbm.dbs[0]
	} else {
		return dbm.dbs[dbEntryType]
	}
}

func (dbm *databaseManager) Close() {
	// If single DB, only close the first database.
	if dbm.config.SingleDB {
		dbm.dbs[0].Close()
		return
	}

	// If not single DB, close all databases.
	for _, db := range dbm.dbs {
		if db != nil {
			db.Close()
		}
	}
}

// TODO-Klaytn Some of below need to be invisible outside database package
// Canonical Hash operations.
// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
func (dbm *databaseManager) ReadCanonicalHash(number uint64) common.Hash {
	if cached := dbm.cm.readCanonicalHashCache(number); !common.EmptyHash(cached) {
		return cached
	}

	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headerHashKey(number))
	if len(data) == 0 {
		return common.Hash{}
	}

	hash := common.BytesToHash(data)
	dbm.cm.writeCanonicalHashCache(number, hash)
	return hash
}

// WriteCanonicalHash stores the hash assigned to a canonical block number.
func (dbm *databaseManager) WriteCanonicalHash(hash common.Hash, number uint64) {
	db := dbm.getDatabase(headerDB)
	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		logger.Crit("Failed to store number to hash mapping", "err", err)
	}
	dbm.cm.writeCanonicalHashCache(number, hash)
}

// DeleteCanonicalHash removes the number to hash canonical mapping.
func (dbm *databaseManager) DeleteCanonicalHash(number uint64) {
	db := dbm.getDatabase(headerDB)
	if err := db.Delete(headerHashKey(number)); err != nil {
		logger.Crit("Failed to delete number to hash mapping", "err", err)
	}
	dbm.cm.writeCanonicalHashCache(number, common.Hash{})
}

// Head Header Hash operations.
// ReadHeadHeaderHash retrieves the hash of the current canonical head header.
func (dbm *databaseManager) ReadHeadHeaderHash() common.Hash {
	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headHeaderKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadHeaderHash stores the hash of the current canonical head header.
func (dbm *databaseManager) WriteHeadHeaderHash(hash common.Hash) {
	db := dbm.getDatabase(headerDB)
	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		logger.Crit("Failed to store last header's hash", "err", err)
	}
}

// Block Hash operations.
func (dbm *databaseManager) ReadHeadBlockHash() common.Hash {
	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// Block Backup Hash operations.
func (dbm *databaseManager) ReadHeadBlockBackupHash() common.Hash {
	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headBlockBackupKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadBlockHash stores the head block's hash.
func (dbm *databaseManager) WriteHeadBlockHash(hash common.Hash) {
	HeadBlockQ.push(hash)

	db := dbm.getDatabase(headerDB)
	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		logger.Crit("Failed to store last block's hash", "err", err)
	}

	backupHash := HeadBlockQ.pop()
	if backupHash == (common.Hash{}) {
		return
	}
	if err := db.Put(headBlockBackupKey, backupHash.Bytes()); err != nil {
		logger.Crit("Failed to store last block's backup hash", "err", err)
	}
}

// Head Fast Block Hash operations.
// ReadHeadFastBlockHash retrieves the hash of the current fast-sync head block.
func (dbm *databaseManager) ReadHeadFastBlockHash() common.Hash {
	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headFastBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// Head Fast Block Backup Hash operations.
// ReadHeadFastBlockBackupHash retrieves the hash of the current fast-sync head block.
func (dbm *databaseManager) ReadHeadFastBlockBackupHash() common.Hash {
	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headFastBlockBackupKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadFastBlockHash stores the hash of the current fast-sync head block.
func (dbm *databaseManager) WriteHeadFastBlockHash(hash common.Hash) {
	FastBlockQ.push(hash)

	db := dbm.getDatabase(headerDB)
	if err := db.Put(headFastBlockKey, hash.Bytes()); err != nil {
		logger.Crit("Failed to store last fast block's hash", "err", err)
	}

	backupHash := FastBlockQ.pop()
	if backupHash == (common.Hash{}) {
		return
	}
	if err := db.Put(headFastBlockBackupKey, backupHash.Bytes()); err != nil {
		logger.Crit("Failed to store last fast block's backup hash", "err", err)
	}
}

// Fast Trie Progress operations.
// ReadFastTrieProgress retrieves the number of tries nodes fast synced to allow
// reporting correct numbers across restarts.
func (dbm *databaseManager) ReadFastTrieProgress() uint64 {
	db := dbm.getDatabase(MiscDB)
	data, _ := db.Get(fastTrieProgressKey)
	if len(data) == 0 {
		return 0
	}
	return new(big.Int).SetBytes(data).Uint64()
}

// WriteFastTrieProgress stores the fast sync trie process counter to support
// retrieving it across restarts.
func (dbm *databaseManager) WriteFastTrieProgress(count uint64) {
	db := dbm.getDatabase(MiscDB)
	if err := db.Put(fastTrieProgressKey, new(big.Int).SetUint64(count).Bytes()); err != nil {
		logger.Crit("Failed to store fast sync trie progress", "err", err)
	}
}

// (Block)Header operations.
// HasHeader verifies the existence of a block header corresponding to the hash.
func (dbm *databaseManager) HasHeader(hash common.Hash, number uint64) bool {
	if dbm.cm.hasHeaderInCache(hash) {
		return true
	}

	db := dbm.getDatabase(headerDB)
	if has, err := db.Has(headerKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadHeader retrieves the block header corresponding to the hash.
func (dbm *databaseManager) ReadHeader(hash common.Hash, number uint64) *types.Header {
	if cachedHeader := dbm.cm.readHeaderCache(hash); cachedHeader != nil {
		return cachedHeader
	}

	data := dbm.ReadHeaderRLP(hash, number)
	if len(data) == 0 {
		return nil
	}
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		logger.Error("Invalid block header RLP", "hash", hash, "err", err)
		return nil
	}

	// Write to cache before returning found value.
	dbm.cm.writeHeaderCache(hash, header)
	return header
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
func (dbm *databaseManager) ReadHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headerKey(number, hash))
	return data
}

// WriteHeader stores a block header into the database and also stores the hash-
// to-number mapping.
func (dbm *databaseManager) WriteHeader(header *types.Header) {
	db := dbm.getDatabase(headerDB)
	// Write the hash -> number mapping
	var (
		hash    = header.Hash()
		number  = header.Number.Uint64()
		encoded = common.Int64ToByteBigEndian(number)
	)
	key := headerNumberKey(hash)
	if err := db.Put(key, encoded); err != nil {
		logger.Crit("Failed to store hash to number mapping", "err", err)
	}
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		logger.Crit("Failed to RLP encode header", "err", err)
	}
	key = headerKey(number, hash)
	if err := db.Put(key, data); err != nil {
		logger.Crit("Failed to store header", "err", err)
	}

	// Write to cache at the end of successful write.
	dbm.cm.writeHeaderCache(hash, header)
	dbm.cm.writeBlockNumberCache(hash, number)
}

// DeleteHeader removes all block header data associated with a hash.
func (dbm *databaseManager) DeleteHeader(hash common.Hash, number uint64) {
	db := dbm.getDatabase(headerDB)
	if err := db.Delete(headerKey(number, hash)); err != nil {
		logger.Crit("Failed to delete header", "err", err)
	}
	if err := db.Delete(headerNumberKey(hash)); err != nil {
		logger.Crit("Failed to delete hash to number mapping", "err", err)
	}

	// Delete cache at the end of successful delete.
	dbm.cm.deleteHeaderCache(hash)
	dbm.cm.deleteBlockNumberCache(hash)
}

// Head Number operations.
// ReadHeaderNumber returns the header number assigned to a hash.
func (dbm *databaseManager) ReadHeaderNumber(hash common.Hash) *uint64 {
	if cachedHeaderNumber := dbm.cm.readBlockNumberCache(hash); cachedHeaderNumber != nil {
		return cachedHeaderNumber
	}

	db := dbm.getDatabase(headerDB)
	data, _ := db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)

	// Write to cache before returning found value.
	dbm.cm.writeBlockNumberCache(hash, number)
	return &number
}

// (Block)Body operations.
// HasBody verifies the existence of a block body corresponding to the hash.
func (dbm *databaseManager) HasBody(hash common.Hash, number uint64) bool {
	db := dbm.getDatabase(BodyDB)
	if has, err := db.Has(blockBodyKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadBody retrieves the block body corresponding to the hash.
func (dbm *databaseManager) ReadBody(hash common.Hash, number uint64) *types.Body {
	if cachedBody := dbm.cm.readBodyCache(hash); cachedBody != nil {
		return cachedBody
	}

	data := dbm.ReadBodyRLP(hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new(types.Body)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		logger.Error("Invalid block body RLP", "hash", hash, "err", err)
		return nil
	}

	// Write to cache at the end of successful read.
	dbm.cm.writeBodyCache(hash, body)
	return body
}

// ReadBodyInCache retrieves the block body in bodyCache.
// It only searches cache.
func (dbm *databaseManager) ReadBodyInCache(hash common.Hash) *types.Body {
	return dbm.cm.readBodyCache(hash)
}

// ReadBodyRLP retrieves the block body (transactions) in RLP encoding.
func (dbm *databaseManager) ReadBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	// Short circuit if the rlp encoded body's already in the cache, retrieve otherwise
	if cachedBodyRLP := dbm.readBodyRLPInCache(hash); cachedBodyRLP != nil {
		return cachedBodyRLP
	}

	// find cached body and encode it to return
	if cachedBody := dbm.ReadBodyInCache(hash); cachedBody != nil {
		if bodyRLP, err := rlp.EncodeToBytes(cachedBody); err != nil {
			dbm.cm.writeBodyRLPCache(hash, bodyRLP)
			return bodyRLP
		}
	}

	// not found in cache, find body in database
	db := dbm.getDatabase(BodyDB)
	data, _ := db.Get(blockBodyKey(number, hash))

	// Write to cache at the end of successful read.
	dbm.cm.writeBodyRLPCache(hash, data)
	return data
}

// ReadBodyRLPByHash retrieves the block body (transactions) in RLP encoding.
func (dbm *databaseManager) ReadBodyRLPByHash(hash common.Hash) rlp.RawValue {
	// Short circuit if the rlp encoded body's already in the cache, retrieve otherwise
	if cachedBodyRLP := dbm.readBodyRLPInCache(hash); cachedBodyRLP != nil {
		return cachedBodyRLP
	}

	// find cached body and encode it to return
	if cachedBody := dbm.ReadBodyInCache(hash); cachedBody != nil {
		if bodyRLP, err := rlp.EncodeToBytes(cachedBody); err != nil {
			dbm.cm.writeBodyRLPCache(hash, bodyRLP)
			return bodyRLP
		}
	}

	// not found in cache, find body in database
	number := dbm.ReadHeaderNumber(hash)
	if number == nil {
		return nil
	}

	db := dbm.getDatabase(BodyDB)
	data, _ := db.Get(blockBodyKey(*number, hash))

	// Write to cache at the end of successful read.
	dbm.cm.writeBodyRLPCache(hash, data)
	return data
}

// readBodyRLPInCache retrieves the block body (transactions) in RLP encoding
// in bodyRLPCache. It only searches cache.
func (dbm *databaseManager) readBodyRLPInCache(hash common.Hash) rlp.RawValue {
	return dbm.cm.readBodyRLPCache(hash)
}

// WriteBody storea a block body into the database.
func (dbm *databaseManager) WriteBody(hash common.Hash, number uint64, body *types.Body) {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		logger.Crit("Failed to RLP encode body", "err", err)
	}
	dbm.WriteBodyRLP(hash, number, data)
}

func (dbm *databaseManager) PutBodyToBatch(batch Batch, hash common.Hash, number uint64, body *types.Body) {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		logger.Crit("Failed to RLP encode body", "err", err)
	}

	if err := batch.Put(blockBodyKey(number, hash), data); err != nil {
		logger.Crit("Failed to store block body", "err", err)
	}
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func (dbm *databaseManager) WriteBodyRLP(hash common.Hash, number uint64, rlp rlp.RawValue) {
	dbm.cm.writeBodyRLPCache(hash, rlp)

	db := dbm.getDatabase(BodyDB)
	if err := db.Put(blockBodyKey(number, hash), rlp); err != nil {
		logger.Crit("Failed to store block body", "err", err)
	}
}

// DeleteBody removes all block body data associated with a hash.
func (dbm *databaseManager) DeleteBody(hash common.Hash, number uint64) {
	db := dbm.getDatabase(BodyDB)
	if err := db.Delete(blockBodyKey(number, hash)); err != nil {
		logger.Crit("Failed to delete block body", "err", err)
	}
	dbm.cm.deleteBodyCache(hash)
}

// TotalDifficulty operations.
// ReadTd retrieves a block's total blockscore corresponding to the hash.
func (dbm *databaseManager) ReadTd(hash common.Hash, number uint64) *big.Int {
	if cachedTd := dbm.cm.readTdCache(hash); cachedTd != nil {
		return cachedTd
	}

	db := dbm.getDatabase(MiscDB)
	data, _ := db.Get(headerTDKey(number, hash))
	if len(data) == 0 {
		return nil
	}
	td := new(big.Int)
	if err := rlp.Decode(bytes.NewReader(data), td); err != nil {
		logger.Error("Invalid block total blockscore RLP", "hash", hash, "err", err)
		return nil
	}

	// Write to cache before returning found value.
	dbm.cm.writeTdCache(hash, td)
	return td
}

// WriteTd stores the total blockscore of a block into the database.
func (dbm *databaseManager) WriteTd(hash common.Hash, number uint64, td *big.Int) {
	db := dbm.getDatabase(MiscDB)
	data, err := rlp.EncodeToBytes(td)
	if err != nil {
		logger.Crit("Failed to RLP encode block total blockscore", "err", err)
	}
	if err := db.Put(headerTDKey(number, hash), data); err != nil {
		logger.Crit("Failed to store block total blockscore", "err", err)
	}

	// Write to cache at the end of successful write.
	dbm.cm.writeTdCache(hash, td)
}

// DeleteTd removes all block total blockscore data associated with a hash.
func (dbm *databaseManager) DeleteTd(hash common.Hash, number uint64) {
	db := dbm.getDatabase(MiscDB)
	if err := db.Delete(headerTDKey(number, hash)); err != nil {
		logger.Crit("Failed to delete block total blockscore", "err", err)
	}
	// Delete cache at the end of successful delete.
	dbm.cm.deleteTdCache(hash)
}

// Receipts operations.
// ReadReceipt retrieves a receipt, blockHash, blockNumber and receiptIndex found by the given txHash.
func (dbm *databaseManager) ReadReceipt(txHash common.Hash) (*types.Receipt, common.Hash, uint64, uint64) {
	blockHash, blockNumber, receiptIndex := dbm.ReadTxLookupEntry(txHash)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}
	receipts := dbm.ReadReceipts(blockHash, blockNumber)
	if len(receipts) <= int(receiptIndex) {
		logger.Error("Receipt refereced missing", "number", blockNumber, "txHash", blockHash, "index", receiptIndex)
		return nil, common.Hash{}, 0, 0
	}
	return receipts[receiptIndex], blockHash, blockNumber, receiptIndex
}

// ReadReceipts retrieves all the transaction receipts belonging to a block.
func (dbm *databaseManager) ReadReceipts(blockHash common.Hash, number uint64) types.Receipts {
	db := dbm.getDatabase(ReceiptsDB)
	// Retrieve the flattened receipt slice
	data, _ := db.Get(blockReceiptsKey(number, blockHash))
	if len(data) == 0 {
		return nil
	}
	// Convert the revceipts from their database form to their internal representation
	storageReceipts := []*types.ReceiptForStorage{}
	if err := rlp.DecodeBytes(data, &storageReceipts); err != nil {
		logger.Error("Invalid receipt array RLP", "blockHash", blockHash, "err", err)
		return nil
	}
	receipts := make(types.Receipts, len(storageReceipts))
	for i, receipt := range storageReceipts {
		receipts[i] = (*types.Receipt)(receipt)
	}
	return receipts
}

func (dbm *databaseManager) ReadReceiptsByBlockHash(hash common.Hash) types.Receipts {
	receipts := dbm.ReadBlockReceiptsInCache(hash)
	if receipts != nil {
		return receipts
	}
	number := dbm.ReadHeaderNumber(hash)
	if number == nil {
		return nil
	}
	return dbm.ReadReceipts(hash, *number)
}

// WriteReceipts stores all the transaction receipts belonging to a block.
func (dbm *databaseManager) WriteReceipts(hash common.Hash, number uint64, receipts types.Receipts) {
	dbm.cm.writeBlockReceiptsCache(hash, receipts)

	db := dbm.getDatabase(ReceiptsDB)
	// When putReceiptsToPutter is called from WriteReceipts, txReceipt is cached.
	dbm.putReceiptsToPutter(db, hash, number, receipts, true)
}

func (dbm *databaseManager) PutReceiptsToBatch(batch Batch, hash common.Hash, number uint64, receipts types.Receipts) {
	// When putReceiptsToPutter is called from PutReceiptsToBatch, txReceipt is not cached.
	dbm.putReceiptsToPutter(batch, hash, number, receipts, false)
}

func (dbm *databaseManager) putReceiptsToPutter(putter KeyValueWriter, hash common.Hash, number uint64, receipts types.Receipts, addToCache bool) {
	// Convert the receipts into their database form and serialize them
	storageReceipts := make([]*types.ReceiptForStorage, len(receipts))
	for i, receipt := range receipts {
		storageReceipts[i] = (*types.ReceiptForStorage)(receipt)

		if addToCache {
			dbm.cm.writeTxReceiptCache(receipt.TxHash, receipt)
		}
	}
	bytes, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		logger.Crit("Failed to encode block receipts", "err", err)
	}
	// Store the flattened receipt slice
	if err := putter.Put(blockReceiptsKey(number, hash), bytes); err != nil {
		logger.Crit("Failed to store block receipts", "err", err)
	}
}

// DeleteReceipts removes all receipt data associated with a block hash.
func (dbm *databaseManager) DeleteReceipts(hash common.Hash, number uint64) {
	receipts := dbm.ReadReceipts(hash, number)

	db := dbm.getDatabase(ReceiptsDB)
	if err := db.Delete(blockReceiptsKey(number, hash)); err != nil {
		logger.Crit("Failed to delete block receipts", "err", err)
	}

	// Delete blockReceiptsCache and txReceiptCache.
	dbm.cm.deleteBlockReceiptsCache(hash)
	if receipts != nil {
		for _, receipt := range receipts {
			dbm.cm.deleteTxReceiptCache(receipt.TxHash)
		}
	}
}

// Block operations.
// ReadBlock retrieves an entire block corresponding to the hash, assembling it
// back from the stored header and body. If either the header or body could not
// be retrieved nil is returned.
//
// Note, due to concurrent download of header and block body the header and thus
// canonical hash can be stored in the database but the body data not (yet).
func (dbm *databaseManager) ReadBlock(hash common.Hash, number uint64) *types.Block {
	if cachedBlock := dbm.cm.readBlockCache(hash); cachedBlock != nil {
		return cachedBlock
	}

	header := dbm.ReadHeader(hash, number)
	if header == nil {
		return nil
	}

	body := dbm.ReadBody(hash, number)
	if body == nil {
		return nil
	}

	block := types.NewBlockWithHeader(header).WithBody(body.Transactions)

	// Write to cache at the end of successful write.
	dbm.cm.writeBlockCache(hash, block)
	return block
}

func (dbm *databaseManager) ReadBlockByHash(hash common.Hash) *types.Block {
	if cachedBlock := dbm.cm.readBlockCache(hash); cachedBlock != nil {
		return cachedBlock
	}

	number := dbm.ReadHeaderNumber(hash)
	if number == nil {
		return nil
	}

	header := dbm.ReadHeader(hash, *number)
	if header == nil {
		return nil
	}

	body := dbm.ReadBody(hash, *number)
	if body == nil {
		return nil
	}

	block := types.NewBlockWithHeader(header).WithBody(body.Transactions)

	// Write to cache at the end of successful write.
	dbm.cm.writeBlockCache(hash, block)
	return block
}

func (dbm *databaseManager) ReadBlockByNumber(number uint64) *types.Block {
	hash := dbm.ReadCanonicalHash(number)
	if hash == (common.Hash{}) {
		return nil
	}
	return dbm.ReadBlock(hash, number)
}

func (dbm *databaseManager) HasBlock(hash common.Hash, number uint64) bool {
	if dbm.cm.hasBlockInCache(hash) {
		return true
	}
	return dbm.HasBody(hash, number)
}

func (dbm *databaseManager) WriteBlock(block *types.Block) {
	dbm.cm.writeBodyCache(block.Hash(), block.Body())
	dbm.cm.blockCache.Add(block.Hash(), block)

	dbm.WriteBody(block.Hash(), block.NumberU64(), block.Body())
	dbm.WriteHeader(block.Header())
}

func (dbm *databaseManager) DeleteBlock(hash common.Hash, number uint64) {
	dbm.DeleteReceipts(hash, number)
	dbm.DeleteHeader(hash, number)
	dbm.DeleteBody(hash, number)
	dbm.DeleteTd(hash, number)
	dbm.cm.deleteBlockCache(hash)
}

const badBlockToKeep = 100

type badBlock struct {
	Header *types.Header
	Body   *types.Body
}

// badBlockList implements the sort interface to allow sorting a list of
// bad blocks by their number in the reverse order.
type badBlockList []*badBlock

func (s badBlockList) Len() int { return len(s) }
func (s badBlockList) Less(i, j int) bool {
	return s[i].Header.Number.Uint64() < s[j].Header.Number.Uint64()
}
func (s badBlockList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ReadBadBlock retrieves the bad block with the corresponding block hash.
func (dbm *databaseManager) ReadBadBlock(hash common.Hash) *types.Block {
	db := dbm.getDatabase(MiscDB)
	blob, err := db.Get(badBlockKey)
	if err != nil {
		return nil
	}
	var badBlocks badBlockList
	if err := rlp.DecodeBytes(blob, &badBlocks); err != nil {
		return nil
	}
	for _, bad := range badBlocks {
		if bad.Header.Hash() == hash {
			return types.NewBlockWithHeader(bad.Header).WithBody(bad.Body.Transactions)
		}
	}
	return nil
}

// ReadAllBadBlocks retrieves all the bad blocks in the database.
// All returned blocks are sorted in reverse order by number.
func (dbm *databaseManager) ReadAllBadBlocks() ([]*types.Block, error) {
	var badBlocks badBlockList
	db := dbm.getDatabase(MiscDB)
	blob, err := db.Get(badBlockKey)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(blob, &badBlocks); err != nil {
		return nil, err
	}
	blocks := make([]*types.Block, len(badBlocks))
	for i, bad := range badBlocks {
		blocks[i] = types.NewBlockWithHeader(bad.Header).WithBody(bad.Body.Transactions)
	}
	return blocks, nil
}

// WriteBadBlock serializes the bad block into the database. If the cumulated
// bad blocks exceed the capacity, the oldest will be dropped.
func (dbm *databaseManager) WriteBadBlock(block *types.Block) {
	db := dbm.getDatabase(MiscDB)
	blob, err := db.Get(badBlockKey)
	if err != nil {
		logger.Warn("Failed to load old bad blocks", "error", err)
	}
	var badBlocks badBlockList
	if len(blob) > 0 {
		if err := rlp.DecodeBytes(blob, &badBlocks); err != nil {
			logger.Error("failed to decode old bad blocks")
			return
		}
	}

	for _, badblock := range badBlocks {
		if badblock.Header.Hash() == block.Hash() && badblock.Header.Number.Uint64() == block.NumberU64() {
			logger.Info("There is already corresponding badblock in db.", "badblock number", block.NumberU64())
			return
		}
	}

	badBlocks = append(badBlocks, &badBlock{
		Header: block.Header(),
		Body:   block.Body(),
	})
	sort.Sort(sort.Reverse(badBlocks))
	if len(badBlocks) > badBlockToKeep {
		badBlocks = badBlocks[:badBlockToKeep]
	}
	data, err := rlp.EncodeToBytes(badBlocks)
	if err != nil {
		logger.Crit("Failed to encode bad blocks", "err", err)
		return
	}
	if err := db.Put(badBlockKey, data); err != nil {
		logger.Crit("Failed to write bad blocks", "err", err)
		return
	}
}

func (dbm *databaseManager) DeleteBadBlocks() {
	db := dbm.getDatabase(MiscDB)
	if err := db.Delete(badBlockKey); err != nil {
		logger.Error("Failed to delete bad blocks", "err", err)
	}
}

// Find Common Ancestor operation
// FindCommonAncestor returns the last common ancestor of two block headers
func (dbm *databaseManager) FindCommonAncestor(a, b *types.Header) *types.Header {
	for bn := b.Number.Uint64(); a.Number.Uint64() > bn; {
		a = dbm.ReadHeader(a.ParentHash, a.Number.Uint64()-1)
		if a == nil {
			return nil
		}
	}
	for an := a.Number.Uint64(); an < b.Number.Uint64(); {
		b = dbm.ReadHeader(b.ParentHash, b.Number.Uint64()-1)
		if b == nil {
			return nil
		}
	}
	for a.Hash() != b.Hash() {
		a = dbm.ReadHeader(a.ParentHash, a.Number.Uint64()-1)
		if a == nil {
			return nil
		}
		b = dbm.ReadHeader(b.ParentHash, b.Number.Uint64()-1)
		if b == nil {
			return nil
		}
	}
	return a
}

// Istanbul Snapshot operations.
func (dbm *databaseManager) ReadIstanbulSnapshot(hash common.Hash) ([]byte, error) {
	db := dbm.getDatabase(MiscDB)
	return db.Get(snapshotKey(hash))
}

func (dbm *databaseManager) WriteIstanbulSnapshot(hash common.Hash, blob []byte) error {
	db := dbm.getDatabase(MiscDB)
	return db.Put(snapshotKey(hash), blob)
}

func (dbm *databaseManager) DeleteIstanbulSnapshot(hash common.Hash) {
	db := dbm.getDatabase(MiscDB)
	if err := db.Delete(snapshotKey(hash)); err != nil {
		logger.Crit("Failed to delete snpahost", "err", err)
	}
}

// Merkle Proof operation.
func (dbm *databaseManager) WriteMerkleProof(key, value []byte) {
	db := dbm.getDatabase(MiscDB)
	if err := db.Put(key, value); err != nil {
		logger.Crit("Failed to write merkle proof", "err", err)
	}
}

// ReadCode retrieves the contract code of the provided code hash.
func (dbm *databaseManager) ReadCode(hash common.Hash) []byte {
	// Try with the legacy code scheme first, if not then try with current
	// scheme. Since most of the code will be found with legacy scheme.
	//
	// TODO-Klaytn-Snapsync change the order when we forcibly upgrade the code scheme with snapshot.
	db := dbm.getDatabase(StateTrieDB)
	if data, _ := db.Get(hash[:]); len(data) > 0 {
		return data
	}

	return dbm.ReadCodeWithPrefix(hash)
}

// ReadCodeWithPrefix retrieves the contract code of the provided code hash.
// The main difference between this function and ReadCode is this function
// will only check the existence with latest scheme(with prefix).
func (dbm *databaseManager) ReadCodeWithPrefix(hash common.Hash) []byte {
	db := dbm.getDatabase(StateTrieDB)
	data, _ := db.Get(CodeKey(hash))
	return data
}

// HasCode checks if the contract code corresponding to the
// provided code hash is present in the db.
func (dbm *databaseManager) HasCode(hash common.Hash) bool {
	// Try with the prefixed code scheme first, if not then try with legacy
	// scheme.
	//
	// TODO-Klaytn-Snapsync change the order when we forcibly upgrade the code scheme with snapshot.
	db := dbm.getDatabase(StateTrieDB)
	if ok, _ := db.Has(hash.Bytes()); ok {
		return true
	}
	return dbm.HasCodeWithPrefix(hash)
}

// HasCodeWithPrefix checks if the contract code corresponding to the
// provided code hash is present in the db. This function will only check
// presence using the prefix-scheme.
func (dbm *databaseManager) HasCodeWithPrefix(hash common.Hash) bool {
	db := dbm.getDatabase(StateTrieDB)
	ok, _ := db.Has(CodeKey(hash))
	return ok
}

// WriteCode writes the provided contract code database.
func (dbm *databaseManager) WriteCode(hash common.Hash, code []byte) {
	dbm.lockInMigration.RLock()
	defer dbm.lockInMigration.RUnlock()

	dbs := make([]Database, 0, 2)
	if dbm.inMigration {
		dbs = append(dbs, dbm.getDatabase(StateTrieMigrationDB))
	}
	dbs = append(dbs, dbm.getDatabase(StateTrieDB))
	for _, db := range dbs {
		if err := db.Put(CodeKey(hash), code); err != nil {
			logger.Crit("Failed to store contract code", "err", err)
		}
	}
}

func (dbm *databaseManager) PutCodeToBatch(batch Batch, hash common.Hash, code []byte) {
	if err := batch.Put(CodeKey(hash), code); err != nil {
		logger.Crit("Failed to store contract code", "err", err)
	}
}

// DeleteCode deletes the specified contract code from the database.
func (dbm *databaseManager) DeleteCode(hash common.Hash) {
	db := dbm.getDatabase(StateTrieDB)
	if err := db.Delete(CodeKey(hash)); err != nil {
		logger.Crit("Failed to delete contract code", "err", err)
	}
}

func (dbm *databaseManager) ReadTrieNode(hash common.ExtHash) ([]byte, error) {
	dbm.lockInMigration.RLock()
	defer dbm.lockInMigration.RUnlock()

	if dbm.inMigration {
		if val, err := dbm.ReadTrieNodeFromNew(hash); err == nil {
			return val, nil
		} else if err != dataNotFoundErr {
			// TODO-Klaytn-Database Need to be properly handled
			logger.Error("Unexpected error while reading cached trie node from state migration database", "err", err)
		}
	}
	val, err := dbm.ReadTrieNodeFromOld(hash)
	if err != nil && err != dataNotFoundErr {
		// TODO-Klaytn-Database Need to be properly handled
		logger.Error("Unexpected error while reading cached trie node", "err", err)
	}
	return val, err
}

func (dbm *databaseManager) HasTrieNode(hash common.ExtHash) (bool, error) {
	val, err := dbm.ReadTrieNode(hash)
	if val == nil || err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// ReadPreimage retrieves a single preimage of the provided hash.
func (dbm *databaseManager) ReadPreimage(hash common.Hash) []byte {
	dbm.lockInMigration.RLock()
	defer dbm.lockInMigration.RUnlock()

	if dbm.inMigration {
		if val, err := dbm.GetStateTrieMigrationDB().Get(preimageKey(hash)); err == nil {
			return val
		}
	}
	return dbm.ReadPreimageFromOld(hash)
}

func (dbm *databaseManager) ReadTrieNodeFromNew(hash common.ExtHash) ([]byte, error) {
	return dbm.GetStateTrieMigrationDB().Get(TrieNodeKey(hash))
}

func (dbm *databaseManager) HasTrieNodeFromNew(hash common.ExtHash) (bool, error) {
	val, err := dbm.ReadTrieNodeFromNew(hash)
	if val == nil || err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (dbm *databaseManager) HasCodeWithPrefixFromNew(hash common.Hash) bool {
	db := dbm.GetStateTrieMigrationDB()
	ok, _ := db.Has(CodeKey(hash))
	return ok
}

// ReadPreimage retrieves a single preimage of the provided hash.
func (dbm *databaseManager) ReadPreimageFromNew(hash common.Hash) []byte {
	data, _ := dbm.GetStateTrieMigrationDB().Get(preimageKey(hash))
	return data
}

func (dbm *databaseManager) ReadTrieNodeFromOld(hash common.ExtHash) ([]byte, error) {
	db := dbm.getDatabase(StateTrieDB)
	return db.Get(TrieNodeKey(hash))
}

func (dbm *databaseManager) HasTrieNodeFromOld(hash common.ExtHash) (bool, error) {
	val, err := dbm.ReadTrieNodeFromOld(hash)
	if val == nil || err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (dbm *databaseManager) HasCodeWithPrefixFromOld(hash common.Hash) bool {
	db := dbm.getDatabase(StateTrieDB)
	ok, _ := db.Has(CodeKey(hash))
	return ok
}

// ReadPreimage retrieves a single preimage of the provided hash.
func (dbm *databaseManager) ReadPreimageFromOld(hash common.Hash) []byte {
	db := dbm.getDatabase(StateTrieDB)
	data, _ := db.Get(preimageKey(hash))
	return data
}

func (dbm *databaseManager) WriteTrieNode(hash common.ExtHash, node []byte) {
	dbm.lockInMigration.RLock()
	defer dbm.lockInMigration.RUnlock()

	if dbm.inMigration {
		if err := dbm.getDatabase(StateTrieMigrationDB).Put(TrieNodeKey(hash), node); err != nil {
			logger.Crit("Failed to store trie node", "err", err)
		}
	}
	if err := dbm.getDatabase(StateTrieDB).Put(TrieNodeKey(hash), node); err != nil {
		logger.Crit("Failed to store trie node", "err", err)
	}
}

func (dbm *databaseManager) PutTrieNodeToBatch(batch Batch, hash common.ExtHash, node []byte) {
	if err := batch.Put(TrieNodeKey(hash), node); err != nil {
		logger.Crit("Failed to store trie node", "err", err)
	}
}

// DeleteTrieNode deletes a trie node having a specific hash. It is used only for testing.
func (dbm *databaseManager) DeleteTrieNode(hash common.ExtHash) {
	if err := dbm.getDatabase(StateTrieDB).Delete(TrieNodeKey(hash)); err != nil {
		logger.Crit("Failed to delete trie node", "err", err)
	}
}

// WritePreimages writes the provided set of preimages to the database. `number` is the
// current block number, and is used for debug messages only.
func (dbm *databaseManager) WritePreimages(number uint64, preimages map[common.Hash][]byte) {
	batch := dbm.NewBatch(StateTrieDB)
	defer batch.Release()
	for hash, preimage := range preimages {
		if err := batch.Put(preimageKey(hash), preimage); err != nil {
			logger.Crit("Failed to store trie preimage", "err", err)
		}
		if _, err := WriteBatchesOverThreshold(batch); err != nil {
			logger.Crit("Failed to store trie preimage", "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to batch write trie preimage", "err", err, "blockNumber", number)
	}
	preimageCounter.Inc(int64(len(preimages)))
	preimageHitCounter.Inc(int64(len(preimages)))
}

// ReadPruningEnabled reads if the live pruning flag is stored in database.
func (dbm *databaseManager) ReadPruningEnabled() bool {
	ok, _ := dbm.getDatabase(MiscDB).Has(pruningEnabledKey)
	return ok
}

// WritePruningEnabled writes the live pruning flag to the database.
func (dbm *databaseManager) WritePruningEnabled() {
	if err := dbm.getDatabase(MiscDB).Put(pruningEnabledKey, []byte("42")); err != nil {
		logger.Crit("Failed to store pruning enabled flag", "err", err)
	}
}

// DeletePruningEnabled deletes the live pruning flag. It is used only for testing.
func (dbm *databaseManager) DeletePruningEnabled() {
	if err := dbm.getDatabase(MiscDB).Delete(pruningEnabledKey); err != nil {
		logger.Crit("Failed to remove pruning enabled flag", "err", err)
	}
}

// WritePruningMarks writes the provided set of pruning marks to the database.
func (dbm *databaseManager) WritePruningMarks(marks []PruningMark) {
	batch := dbm.NewBatch(MiscDB)
	defer batch.Release()
	for _, mark := range marks {
		if err := batch.Put(pruningMarkKey(mark), pruningMarkValue); err != nil {
			logger.Crit("Failed to store trie pruning mark", "err", err)
		}
		if _, err := WriteBatchesOverThreshold(batch); err != nil {
			logger.Crit("Failed to store trie pruning mark", "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to batch write pruning mark", "err", err)
	}
}

// ReadPruningMarks reads the pruning marks in the block number range [startNumber, endNumber).
func (dbm *databaseManager) ReadPruningMarks(startNumber, endNumber uint64) []PruningMark {
	prefix := pruningMarkPrefix
	startKey := pruningMarkKey(PruningMark{startNumber, common.ExtHash{}})
	it := dbm.getDatabase(MiscDB).NewIterator(prefix, startKey[len(prefix):])
	defer it.Release()

	var marks []PruningMark
	for it.Next() {
		mark := parsePruningMarkKey(it.Key())
		if endNumber != 0 && mark.Number >= endNumber {
			break
		}
		marks = append(marks, mark)
	}
	return marks
}

// DeletePruningMarks deletes the provided set of pruning marks from the database.
// Note that trie nodes are not deleted by this function. To prune trie nodes, use
// the PruneTrieNodes or DeleteTrieNode functions.
func (dbm *databaseManager) DeletePruningMarks(marks []PruningMark) {
	batch := dbm.NewBatch(MiscDB)
	defer batch.Release()
	for _, mark := range marks {
		if err := batch.Delete(pruningMarkKey(mark)); err != nil {
			logger.Crit("Failed to delete trie pruning mark", "err", err)
		}
		if _, err := WriteBatchesOverThreshold(batch); err != nil {
			logger.Crit("Failed to delete trie pruning mark", "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to batch delete pruning mark", "err", err)
	}
}

// PruneTrieNodes deletes the trie nodes according to the provided set of pruning marks.
func (dbm *databaseManager) PruneTrieNodes(marks []PruningMark) {
	batch := dbm.NewBatch(StateTrieDB)
	defer batch.Release()
	for _, mark := range marks {
		if err := batch.Delete(TrieNodeKey(mark.Hash)); err != nil {
			logger.Crit("Failed to prune trie node", "err", err)
		}
		if _, err := WriteBatchesOverThreshold(batch); err != nil {
			logger.Crit("Failed to prune trie node", "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to batch prune trie node", "err", err)
	}
}

// WriteLastPrunedBlockNumber records a block number of the most recent pruning block
func (dbm *databaseManager) WriteLastPrunedBlockNumber(blockNumber uint64) {
	db := dbm.getDatabase(MiscDB)
	if err := db.Put(lastPrunedBlockNumberKey, common.Int64ToByteLittleEndian(blockNumber)); err != nil {
		logger.Crit("Failed to record the last pruned block number", "err", err)
	}
}

// ReadLastPrunedBlockNumber reads a block number of the most recent pruning block
func (dbm *databaseManager) ReadLastPrunedBlockNumber() (uint64, error) {
	db := dbm.getDatabase(MiscDB)
	lastPruned, err := db.Get(lastPrunedBlockNumberKey)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(lastPruned), nil
}

// ReadTxLookupEntry retrieves the positional metadata associated with a transaction
// hash to allow retrieving the transaction or receipt by hash.
func (dbm *databaseManager) ReadTxLookupEntry(hash common.Hash) (common.Hash, uint64, uint64) {
	db := dbm.getDatabase(TxLookUpEntryDB)
	data, _ := db.Get(TxLookupKey(hash))
	if len(data) == 0 {
		return common.Hash{}, 0, 0
	}
	var entry TxLookupEntry
	if err := rlp.DecodeBytes(data, &entry); err != nil {
		logger.Error("Invalid transaction lookup entry RLP", "hash", hash, "err", err)
		return common.Hash{}, 0, 0
	}
	return entry.BlockHash, entry.BlockIndex, entry.Index
}

// WriteTxLookupEntries stores a positional metadata for every transaction from
// a block, enabling hash based transaction and receipt lookups.
func (dbm *databaseManager) WriteTxLookupEntries(block *types.Block) {
	db := dbm.getDatabase(TxLookUpEntryDB)
	putTxLookupEntriesToPutter(db, block)
}

func (dbm *databaseManager) WriteAndCacheTxLookupEntries(block *types.Block) error {
	batch := dbm.NewBatch(TxLookUpEntryDB)
	defer batch.Release()
	for i, tx := range block.Transactions() {
		entry := TxLookupEntry{
			BlockHash:  block.Hash(),
			BlockIndex: block.NumberU64(),
			Index:      uint64(i),
		}
		data, err := rlp.EncodeToBytes(entry)
		if err != nil {
			logger.Crit("Failed to encode transaction lookup entry", "err", err)
		}
		if err := batch.Put(TxLookupKey(tx.Hash()), data); err != nil {
			logger.Crit("Failed to store transaction lookup entry", "err", err)
		}

		// Write to cache at the end of successful Put.
		dbm.cm.writeTxAndLookupInfoCache(tx.Hash(), &TransactionLookup{tx, &entry})
	}
	if err := batch.Write(); err != nil {
		logger.Error("Failed to write TxLookupEntries in batch", "err", err, "blockNumber", block.Number())
		return err
	}
	return nil
}

func (dbm *databaseManager) PutTxLookupEntriesToBatch(batch Batch, block *types.Block) {
	putTxLookupEntriesToPutter(batch, block)
}

func putTxLookupEntriesToPutter(putter KeyValueWriter, block *types.Block) {
	for i, tx := range block.Transactions() {
		entry := TxLookupEntry{
			BlockHash:  block.Hash(),
			BlockIndex: block.NumberU64(),
			Index:      uint64(i),
		}
		data, err := rlp.EncodeToBytes(entry)
		if err != nil {
			logger.Crit("Failed to encode transaction lookup entry", "err", err)
		}
		if err := putter.Put(TxLookupKey(tx.Hash()), data); err != nil {
			logger.Crit("Failed to store transaction lookup entry", "err", err)
		}
	}
}

// DeleteTxLookupEntry removes all transaction data associated with a hash.
func (dbm *databaseManager) DeleteTxLookupEntry(hash common.Hash) {
	db := dbm.getDatabase(TxLookUpEntryDB)
	db.Delete(TxLookupKey(hash))
}

// ReadTxAndLookupInfo retrieves a specific transaction from the database, along with
// its added positional metadata.
func (dbm *databaseManager) ReadTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	blockHash, blockNumber, txIndex := dbm.ReadTxLookupEntry(hash)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}
	body := dbm.ReadBody(blockHash, blockNumber)
	if body == nil || len(body.Transactions) <= int(txIndex) {
		logger.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash, "index", txIndex)
		return nil, common.Hash{}, 0, 0
	}
	return body.Transactions[txIndex], blockHash, blockNumber, txIndex
}

// NewSenderTxHashToTxHashBatch returns a batch to write senderTxHash to txHash mapping information.
func (dbm *databaseManager) NewSenderTxHashToTxHashBatch() Batch {
	return dbm.NewBatch(MiscDB) // batch.Release should be called from caller
}

// PutSenderTxHashToTxHashToBatch 1) puts the given senderTxHash and txHash to the given batch and
// 2) writes the information to the cache.
func (dbm *databaseManager) PutSenderTxHashToTxHashToBatch(batch Batch, senderTxHash, txHash common.Hash) error {
	if err := batch.Put(SenderTxHashToTxHashKey(senderTxHash), txHash.Bytes()); err != nil {
		return err
	}

	dbm.cm.writeSenderTxHashToTxHashCache(senderTxHash, txHash)

	if batch.ValueSize() > IdealBatchSize {
		batch.Write()
		batch.Reset()
	}

	return nil
}

// ReadTxHashFromSenderTxHash retrieves a txHash corresponding to the given senderTxHash.
func (dbm *databaseManager) ReadTxHashFromSenderTxHash(senderTxHash common.Hash) common.Hash {
	if txHash := dbm.cm.readSenderTxHashToTxHashCache(senderTxHash); !common.EmptyHash(txHash) {
		return txHash
	}

	data, _ := dbm.getDatabase(MiscDB).Get(SenderTxHashToTxHashKey(senderTxHash))
	if len(data) == 0 {
		return common.Hash{}
	}

	txHash := common.BytesToHash(data)
	dbm.cm.writeSenderTxHashToTxHashCache(senderTxHash, txHash)
	return txHash
}

// BloomBits operations.
// ReadBloomBits retrieves the compressed bloom bit vector belonging to the given
// section and bit index from the.
func (dbm *databaseManager) ReadBloomBits(bloomBitsKey []byte) ([]byte, error) {
	db := dbm.getDatabase(MiscDB)
	return db.Get(bloomBitsKey)
}

// WriteBloomBits stores the compressed bloom bits vector belonging to the given
// section and bit index.
func (dbm *databaseManager) WriteBloomBits(bloomBitsKey, bits []byte) error {
	db := dbm.getDatabase(MiscDB)
	return db.Put(bloomBitsKey, bits)
}

// ValidSections operation.
func (dbm *databaseManager) ReadValidSections() ([]byte, error) {
	db := dbm.getDatabase(MiscDB)
	return db.Get(validSectionKey)
}

func (dbm *databaseManager) WriteValidSections(encodedSections []byte) {
	db := dbm.getDatabase(MiscDB)
	db.Put(validSectionKey, encodedSections)
}

// SectionHead operation.
func (dbm *databaseManager) ReadSectionHead(encodedSection []byte) ([]byte, error) {
	db := dbm.getDatabase(MiscDB)
	return db.Get(sectionHeadKey(encodedSection))
}

func (dbm *databaseManager) WriteSectionHead(encodedSection []byte, hash common.Hash) {
	db := dbm.getDatabase(MiscDB)
	db.Put(sectionHeadKey(encodedSection), hash.Bytes())
}

func (dbm *databaseManager) DeleteSectionHead(encodedSection []byte) {
	db := dbm.getDatabase(MiscDB)
	db.Delete(sectionHeadKey(encodedSection))
}

// ReadDatabaseVersion retrieves the version number of the database.
func (dbm *databaseManager) ReadDatabaseVersion() *uint64 {
	db := dbm.getDatabase(MiscDB)
	var version uint64

	enc, _ := db.Get(databaseVerisionKey)
	if len(enc) == 0 {
		return nil
	}

	if err := rlp.DecodeBytes(enc, &version); err != nil {
		logger.Error("Failed to decode database version", "err", err)
		return nil
	}

	return &version
}

// WriteDatabaseVersion stores the version number of the database
func (dbm *databaseManager) WriteDatabaseVersion(version uint64) {
	db := dbm.getDatabase(MiscDB)
	enc, err := rlp.EncodeToBytes(version)
	if err != nil {
		logger.Crit("Failed to encode database version", "err", err)
	}
	if err := db.Put(databaseVerisionKey, enc); err != nil {
		logger.Crit("Failed to store the database version", "err", err)
	}
}

// ReadChainConfig retrieves the consensus settings based on the given genesis hash.
func (dbm *databaseManager) ReadChainConfig(hash common.Hash) *params.ChainConfig {
	db := dbm.getDatabase(MiscDB)
	data, _ := db.Get(configKey(hash))
	if len(data) == 0 {
		return nil
	}
	var config params.ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		logger.Error("Invalid chain config JSON", "hash", hash, "err", err)
		return nil
	}
	return &config
}

func (dbm *databaseManager) WriteChainConfig(hash common.Hash, cfg *params.ChainConfig) {
	db := dbm.getDatabase(MiscDB)
	if cfg == nil {
		return
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		logger.Crit("Failed to JSON encode chain config", "err", err)
	}
	if err := db.Put(configKey(hash), data); err != nil {
		logger.Crit("Failed to store chain config", "err", err)
	}
}

// ReadSnapshotJournal retrieves the serialized in-memory diff layers saved at
// the last shutdown. The blob is expected to be max a few 10s of megabytes.
func (dbm *databaseManager) ReadSnapshotJournal() []byte {
	db := dbm.getDatabase(SnapshotDB)
	data, _ := db.Get(snapshotJournalKey)
	return data
}

// WriteSnapshotJournal stores the serialized in-memory diff layers to save at
// shutdown. The blob is expected to be max a few 10s of megabytes.
func (dbm *databaseManager) WriteSnapshotJournal(journal []byte) {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Put(snapshotJournalKey, journal); err != nil {
		logger.Crit("Failed to store snapshot journal", "err", err)
	}
}

// DeleteSnapshotJournal deletes the serialized in-memory diff layers saved at
// the last shutdown
func (dbm *databaseManager) DeleteSnapshotJournal() {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Delete(snapshotJournalKey); err != nil {
		logger.Crit("Failed to remove snapshot journal", "err", err)
	}
}

// ReadSnapshotGenerator retrieves the serialized snapshot generator saved at
// the last shutdown.
func (dbm *databaseManager) ReadSnapshotGenerator() []byte {
	db := dbm.getDatabase(SnapshotDB)
	data, _ := db.Get(SnapshotGeneratorKey)
	return data
}

// WriteSnapshotGenerator stores the serialized snapshot generator to save at
// shutdown.
func (dbm *databaseManager) WriteSnapshotGenerator(generator []byte) {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Put(SnapshotGeneratorKey, generator); err != nil {
		logger.Crit("Failed to store snapshot generator", "err", err)
	}
}

// DeleteSnapshotGenerator deletes the serialized snapshot generator saved at
// the last shutdown
func (dbm *databaseManager) DeleteSnapshotGenerator() {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Delete(SnapshotGeneratorKey); err != nil {
		logger.Crit("Failed to remove snapshot generator", "err", err)
	}
}

// ReadSnapshotDisabled retrieves if the snapshot maintenance is disabled.
func (dbm *databaseManager) ReadSnapshotDisabled() bool {
	db := dbm.getDatabase(SnapshotDB)
	disabled, _ := db.Has(snapshotDisabledKey)
	return disabled
}

// WriteSnapshotDisabled stores the snapshot pause flag.
func (dbm *databaseManager) WriteSnapshotDisabled() {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Put(snapshotDisabledKey, []byte("42")); err != nil {
		logger.Crit("Failed to store snapshot disabled flag", "err", err)
	}
}

// DeleteSnapshotDisabled deletes the flag keeping the snapshot maintenance disabled.
func (dbm *databaseManager) DeleteSnapshotDisabled() {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Delete(snapshotDisabledKey); err != nil {
		logger.Crit("Failed to remove snapshot disabled flag", "err", err)
	}
}

// ReadSnapshotRecoveryNumber retrieves the block number of the last persisted
// snapshot layer.
func (dbm *databaseManager) ReadSnapshotRecoveryNumber() *uint64 {
	db := dbm.getDatabase(SnapshotDB)
	data, _ := db.Get(snapshotRecoveryKey)
	if len(data) == 0 {
		return nil
	}
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// WriteSnapshotRecoveryNumber stores the block number of the last persisted
// snapshot layer.
func (dbm *databaseManager) WriteSnapshotRecoveryNumber(number uint64) {
	db := dbm.getDatabase(SnapshotDB)
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], number)
	if err := db.Put(snapshotRecoveryKey, buf[:]); err != nil {
		logger.Crit("Failed to store snapshot recovery number", "err", err)
	}
}

// DeleteSnapshotRecoveryNumber deletes the block number of the last persisted
// snapshot layer.
func (dbm *databaseManager) DeleteSnapshotRecoveryNumber() {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Delete(snapshotRecoveryKey); err != nil {
		logger.Crit("Failed to remove snapshot recovery number", "err", err)
	}
}

// ReadSnapshotSyncStatus retrieves the serialized sync status saved at shutdown.
func (dbm *databaseManager) ReadSnapshotSyncStatus() []byte {
	db := dbm.getDatabase(SnapshotDB)
	data, _ := db.Get(snapshotSyncStatusKey)
	return data
}

// WriteSnapshotSyncStatus stores the serialized sync status to save at shutdown.
func (dbm *databaseManager) WriteSnapshotSyncStatus(status []byte) {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Put(snapshotSyncStatusKey, status); err != nil {
		logger.Crit("Failed to store snapshot sync status", "err", err)
	}
}

// DeleteSnapshotSyncStatus deletes the serialized sync status saved at the last
// shutdown
func (dbm *databaseManager) DeleteSnapshotSyncStatus() {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Delete(snapshotSyncStatusKey); err != nil {
		logger.Crit("Failed to remove snapshot sync status", "err", err)
	}
}

// ReadSnapshotRoot retrieves the root of the block whose state is contained in
// the persisted snapshot.
func (dbm *databaseManager) ReadSnapshotRoot() common.Hash {
	db := dbm.getDatabase(SnapshotDB)
	data, _ := db.Get(snapshotRootKey)
	if len(data) != common.HashLength {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteSnapshotRoot stores the root of the block whose state is contained in
// the persisted snapshot.
func (dbm *databaseManager) WriteSnapshotRoot(root common.Hash) {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Put(snapshotRootKey, root[:]); err != nil {
		logger.Crit("Failed to store snapshot root", "err", err)
	}
}

// DeleteSnapshotRoot deletes the hash of the block whose state is contained in
// the persisted snapshot. Since snapshots are not immutable, this  method can
// be used during updates, so a crash or failure will mark the entire snapshot
// invalid.
func (dbm *databaseManager) DeleteSnapshotRoot() {
	db := dbm.getDatabase(SnapshotDB)
	if err := db.Delete(snapshotRootKey); err != nil {
		logger.Crit("Failed to remove snapshot root", "err", err)
	}
}

// ReadAccountSnapshot retrieves the snapshot entry of an account trie leaf.
func (dbm *databaseManager) ReadAccountSnapshot(hash common.Hash) []byte {
	db := dbm.getDatabase(SnapshotDB)
	data, _ := db.Get(AccountSnapshotKey(hash))
	return data
}

// WriteAccountSnapshot stores the snapshot entry of an account trie leaf.
func (dbm *databaseManager) WriteAccountSnapshot(hash common.Hash, entry []byte) {
	db := dbm.getDatabase(SnapshotDB)
	writeAccountSnapshot(db, hash, entry)
}

// DeleteAccountSnapshot removes the snapshot entry of an account trie leaf.
func (dbm *databaseManager) DeleteAccountSnapshot(hash common.Hash) {
	db := dbm.getDatabase(SnapshotDB)
	deleteAccountSnapshot(db, hash)
}

// ReadStorageSnapshot retrieves the snapshot entry of an storage trie leaf.
func (dbm *databaseManager) ReadStorageSnapshot(accountHash, storageHash common.Hash) []byte {
	db := dbm.getDatabase(SnapshotDB)
	data, _ := db.Get(StorageSnapshotKey(accountHash, storageHash))
	return data
}

// WriteStorageSnapshot stores the snapshot entry of an storage trie leaf.
func (dbm *databaseManager) WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte) {
	db := dbm.getDatabase(SnapshotDB)
	writeStorageSnapshot(db, accountHash, storageHash, entry)
}

// DeleteStorageSnapshot removes the snapshot entry of an storage trie leaf.
func (dbm *databaseManager) DeleteStorageSnapshot(accountHash, storageHash common.Hash) {
	db := dbm.getDatabase(SnapshotDB)
	deleteStorageSnapshot(db, accountHash, storageHash)
}

func (dbm *databaseManager) NewSnapshotDBIterator(prefix []byte, start []byte) Iterator {
	db := dbm.getDatabase(SnapshotDB)
	return db.NewIterator(prefix, start)
}

// WriteChildChainTxHash writes stores a transaction hash of a transaction which contains
// AnchoringData, with the key made with given child chain block hash.
func (dbm *databaseManager) WriteChildChainTxHash(ccBlockHash common.Hash, ccTxHash common.Hash) {
	key := childChainTxHashKey(ccBlockHash)
	db := dbm.getDatabase(bridgeServiceDB)
	if err := db.Put(key, ccTxHash.Bytes()); err != nil {
		logger.Crit("Failed to store ChildChainTxHash", "ccBlockHash", ccBlockHash.String(), "ccTxHash", ccTxHash.String(), "err", err)
	}
}

// ConvertChildChainBlockHashToParentChainTxHash returns a transaction hash of a transaction which contains
// AnchoringData, with the key made with given child chain block hash.
func (dbm *databaseManager) ConvertChildChainBlockHashToParentChainTxHash(scBlockHash common.Hash) common.Hash {
	key := childChainTxHashKey(scBlockHash)
	db := dbm.getDatabase(bridgeServiceDB)
	data, _ := db.Get(key)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteLastIndexedBlockNumber writes the block number which is indexed lastly.
func (dbm *databaseManager) WriteLastIndexedBlockNumber(blockNum uint64) {
	key := lastIndexedBlockKey
	db := dbm.getDatabase(bridgeServiceDB)
	if err := db.Put(key, common.Int64ToByteBigEndian(blockNum)); err != nil {
		logger.Crit("Failed to store LastIndexedBlockNumber", "blockNumber", blockNum, "err", err)
	}
}

// GetLastIndexedBlockNumber returns the last block number which is indexed.
func (dbm *databaseManager) GetLastIndexedBlockNumber() uint64 {
	key := lastIndexedBlockKey
	db := dbm.getDatabase(bridgeServiceDB)
	data, _ := db.Get(key)
	if len(data) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(data)
}

// WriteAnchoredBlockNumber writes the block number whose data has been anchored to the parent chain.
func (dbm *databaseManager) WriteAnchoredBlockNumber(blockNum uint64) {
	key := lastServiceChainTxReceiptKey
	db := dbm.getDatabase(bridgeServiceDB)
	if err := db.Put(key, common.Int64ToByteBigEndian(blockNum)); err != nil {
		logger.Crit("Failed to store LatestServiceChainBlockNum", "blockNumber", blockNum, "err", err)
	}
}

// ReadAnchoredBlockNumber returns the latest block number whose data has been anchored to the parent chain.
func (dbm *databaseManager) ReadAnchoredBlockNumber() uint64 {
	key := lastServiceChainTxReceiptKey
	db := dbm.getDatabase(bridgeServiceDB)
	data, _ := db.Get(key)
	if len(data) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(data)
}

// WriteHandleTxHashFromRequestTxHash writes handle value transfer tx hash
// with corresponding request value transfer tx hash.
func (dbm *databaseManager) WriteHandleTxHashFromRequestTxHash(rTx, hTx common.Hash) {
	db := dbm.getDatabase(bridgeServiceDB)
	key := valueTransferTxHashKey(rTx)
	if err := db.Put(key, hTx.Bytes()); err != nil {
		logger.Crit("Failed to store handle value transfer tx hash", "request tx hash", rTx.String(), "handle tx hash", hTx.String(), "err", err)
	}
}

// ReadHandleTxHashFromRequestTxHash returns handle value transfer tx hash
// with corresponding the given request value transfer tx hash.
func (dbm *databaseManager) ReadHandleTxHashFromRequestTxHash(rTx common.Hash) common.Hash {
	key := valueTransferTxHashKey(rTx)
	db := dbm.getDatabase(bridgeServiceDB)
	data, _ := db.Get(key)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteReceiptFromParentChain writes a receipt received from parent chain to child chain
// with corresponding block hash. It assumes that a child chain has only one parent chain.
func (dbm *databaseManager) WriteReceiptFromParentChain(blockHash common.Hash, receipt *types.Receipt) {
	receiptForStorage := (*types.ReceiptForStorage)(receipt)
	db := dbm.getDatabase(bridgeServiceDB)
	byte, err := rlp.EncodeToBytes(receiptForStorage)
	if err != nil {
		logger.Crit("Failed to RLP encode receipt received from parent chain", "receipt.TxHash", receipt.TxHash, "err", err)
	}
	key := receiptFromParentChainKey(blockHash)
	if err = db.Put(key, byte); err != nil {
		logger.Crit("Failed to store receipt received from parent chain", "receipt.TxHash", receipt.TxHash, "err", err)
	}
}

// ReadReceiptFromParentChain returns a receipt received from parent chain to child chain
// with corresponding block hash. It assumes that a child chain has only one parent chain.
func (dbm *databaseManager) ReadReceiptFromParentChain(blockHash common.Hash) *types.Receipt {
	db := dbm.getDatabase(bridgeServiceDB)
	key := receiptFromParentChainKey(blockHash)
	data, _ := db.Get(key)
	if data == nil || len(data) == 0 {
		return nil
	}
	serviceChainTxReceipt := new(types.ReceiptForStorage)
	if err := rlp.Decode(bytes.NewReader(data), serviceChainTxReceipt); err != nil {
		logger.Error("Invalid Receipt RLP received from parent chain", "err", err)
		return nil
	}
	return (*types.Receipt)(serviceChainTxReceipt)
}

// WriteParentOperatorFeePayer writes a fee payer of parent operator.
func (dbm *databaseManager) WriteParentOperatorFeePayer(feePayer common.Address) {
	key := parentOperatorFeePayerPrefix
	db := dbm.getDatabase(bridgeServiceDB)

	if err := db.Put(key, feePayer.Bytes()); err != nil {
		logger.Crit("Failed to store parent operator fee payer", "feePayer", feePayer.String(), "err", err)
	}
}

// ReadParentOperatorFeePayer returns a fee payer of parent operator.
func (dbm *databaseManager) ReadParentOperatorFeePayer() common.Address {
	key := parentOperatorFeePayerPrefix
	db := dbm.getDatabase(bridgeServiceDB)
	data, _ := db.Get(key)
	if data == nil || len(data) == 0 {
		return common.Address{}
	}
	return common.BytesToAddress(data)
}

// WriteChildOperatorFeePayer writes a fee payer of child operator.
func (dbm *databaseManager) WriteChildOperatorFeePayer(feePayer common.Address) {
	key := childOperatorFeePayerPrefix
	db := dbm.getDatabase(bridgeServiceDB)

	if err := db.Put(key, feePayer.Bytes()); err != nil {
		logger.Crit("Failed to store parent operator fee payer", "feePayer", feePayer.String(), "err", err)
	}
}

// ReadChildOperatorFeePayer returns a fee payer of child operator.
func (dbm *databaseManager) ReadChildOperatorFeePayer() common.Address {
	key := childOperatorFeePayerPrefix
	db := dbm.getDatabase(bridgeServiceDB)
	data, _ := db.Get(key)
	if data == nil || len(data) == 0 {
		return common.Address{}
	}
	return common.BytesToAddress(data)
}

// ClearHeaderChainCache calls cacheManager.clearHeaderChainCache to flush out caches of HeaderChain.
func (dbm *databaseManager) ClearHeaderChainCache() {
	dbm.cm.clearHeaderChainCache()
}

// ClearBlockChainCache calls cacheManager.clearBlockChainCache to flush out caches of BlockChain.
func (dbm *databaseManager) ClearBlockChainCache() {
	dbm.cm.clearBlockChainCache()
}

func (dbm *databaseManager) ReadTxAndLookupInfoInCache(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return dbm.cm.readTxAndLookupInfoInCache(hash)
}

func (dbm *databaseManager) ReadBlockReceiptsInCache(blockHash common.Hash) types.Receipts {
	return dbm.cm.readBlockReceiptsInCache(blockHash)
}

func (dbm *databaseManager) ReadTxReceiptInCache(txHash common.Hash) *types.Receipt {
	return dbm.cm.readTxReceiptInCache(txHash)
}

func (dbm *databaseManager) WriteCliqueSnapshot(snapshotBlockHash common.Hash, encodedSnapshot []byte) error {
	db := dbm.getDatabase(MiscDB)
	return db.Put(snapshotKey(snapshotBlockHash), encodedSnapshot)
}

func (dbm *databaseManager) ReadCliqueSnapshot(snapshotBlockHash common.Hash) ([]byte, error) {
	db := dbm.getDatabase(MiscDB)
	return db.Get(snapshotKey(snapshotBlockHash))
}

func (dbm *databaseManager) WriteGovernance(data map[string]interface{}, num uint64) error {
	db := dbm.getDatabase(MiscDB)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := dbm.WriteGovernanceIdx(num); err != nil {
		if err == errGovIdxAlreadyExist {
			// Overwriting existing data is not allowed, but the attempt is not considered as a failure.
			return nil
		}
		return err
	}
	return db.Put(makeKey(governancePrefix, num), b)
}

func (dbm *databaseManager) DeleteGovernance(num uint64) {
	db := dbm.getDatabase(MiscDB)
	if err := dbm.deleteLastGovernance(num); err != nil {
		logger.Crit("Failed to delete Governance index", "err", err)
	}
	if err := db.Delete(makeKey(governancePrefix, num)); err != nil {
		logger.Crit("Failed to delete Governance", "err", err)
	}
}

func (dbm *databaseManager) WriteGovernanceIdx(num uint64) error {
	db := dbm.getDatabase(MiscDB)
	newSlice := make([]uint64, 0)

	if data, err := db.Get(governanceHistoryKey); err == nil {
		if err = json.Unmarshal(data, &newSlice); err != nil {
			return err
		}
	}

	if len(newSlice) > 0 && num <= newSlice[len(newSlice)-1] {
		logger.Error("The same or more recent governance index exist. Skip writing governance index",
			"newIdx", num, "govIdxes", newSlice)
		return errGovIdxAlreadyExist
	}

	newSlice = append(newSlice, num)

	data, err := json.Marshal(newSlice)
	if err != nil {
		return err
	}
	return db.Put(governanceHistoryKey, data)
}

// deleteLastGovernance deletes the last governanceIdx only if it is equal to `num`
func (dbm *databaseManager) deleteLastGovernance(num uint64) error {
	db := dbm.getDatabase(MiscDB)
	idxHistory, err := dbm.ReadRecentGovernanceIdx(0)
	if err != nil {
		return nil // Do nothing and return nil if no recent index found
	}
	end := len(idxHistory)
	if idxHistory[end-1] != num {
		return nil // Do nothing and return nil if the target number does not match the tip number
	}
	data, err := json.Marshal(idxHistory[0 : end-1])
	if err != nil {
		return err
	}
	return db.Put(governanceHistoryKey, data)
}

func (dbm *databaseManager) ReadGovernance(num uint64) (map[string]interface{}, error) {
	db := dbm.getDatabase(MiscDB)

	if data, err := db.Get(makeKey(governancePrefix, num)); err != nil {
		return nil, err
	} else {
		result := make(map[string]interface{})
		if e := json.Unmarshal(data, &result); e != nil {
			return nil, e
		}
		return result, nil
	}
}

// ReadRecentGovernanceIdx returns latest `count` number of indices. If `count` is 0, it returns all indices.
func (dbm *databaseManager) ReadRecentGovernanceIdx(count int) ([]uint64, error) {
	db := dbm.getDatabase(MiscDB)

	if history, err := db.Get(governanceHistoryKey); err != nil {
		return nil, err
	} else {
		idxHistory := make([]uint64, 0)
		if e := json.Unmarshal(history, &idxHistory); e != nil {
			return nil, e
		}

		// Make sure idxHistory should be in ascending order
		sort.Slice(idxHistory, func(i, j int) bool {
			return idxHistory[i] < idxHistory[j]
		})

		max := 0
		leng := len(idxHistory)
		if leng < count || count == 0 {
			max = leng
		} else {
			max = count
		}
		if count > 0 {
			return idxHistory[leng-max:], nil
		}
		return idxHistory, nil
	}
}

// ReadGovernanceAtNumber returns the block number and governance information which to be used for the block `num`
func (dbm *databaseManager) ReadGovernanceAtNumber(num uint64, epoch uint64) (uint64, map[string]interface{}, error) {
	minimum := num - (num % epoch)
	if minimum >= epoch {
		minimum -= epoch
	}
	totalIdx, _ := dbm.ReadRecentGovernanceIdx(0)
	for i := len(totalIdx) - 1; i >= 0; i-- {
		if totalIdx[i] <= minimum {
			result, err := dbm.ReadGovernance(totalIdx[i])
			return totalIdx[i], result, err
		}
	}
	return 0, nil, errors.New("No governance data found")
}

func (dbm *databaseManager) WriteGovernanceState(b []byte) error {
	db := dbm.getDatabase(MiscDB)
	return db.Put(governanceStateKey, b)
}

func (dbm *databaseManager) ReadGovernanceState() ([]byte, error) {
	db := dbm.getDatabase(MiscDB)
	return db.Get(governanceStateKey)
}

func (dbm *databaseManager) WriteChainDataFetcherCheckpoint(checkpoint uint64) error {
	db := dbm.getDatabase(MiscDB)
	return db.Put(chaindatafetcherCheckpointKey, common.Int64ToByteBigEndian(checkpoint))
}

func (dbm *databaseManager) ReadChainDataFetcherCheckpoint() (uint64, error) {
	db := dbm.getDatabase(MiscDB)
	data, err := db.Get(chaindatafetcherCheckpointKey)
	if err != nil {
		// if the key is not in the database, 0 is returned as the checkpoint
		if err == leveldb.ErrNotFound || err == badger.ErrKeyNotFound ||
			strings.Contains(err.Error(), "not found") { // memoryDB
			return 0, nil
		}
		return 0, err
	}
	// in case that error is nil, but the data does not exist
	if len(data) != 8 {
		logger.Warn("the returned error is nil, but the data is wrong", "len(data)", len(data))
		return 0, nil
	}
	return binary.BigEndian.Uint64(data), nil
}

func (dbm *databaseManager) NewSnapshotDBBatch() SnapshotDBBatch {
	return &snapshotDBBatch{dbm.NewBatch(SnapshotDB)}
}

type SnapshotDBBatch interface {
	Batch

	WriteSnapshotRoot(root common.Hash)
	DeleteSnapshotRoot()

	WriteAccountSnapshot(hash common.Hash, entry []byte)
	DeleteAccountSnapshot(hash common.Hash)

	WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte)
	DeleteStorageSnapshot(accountHash, storageHash common.Hash)

	WriteSnapshotJournal(journal []byte)
	DeleteSnapshotJournal()

	WriteSnapshotGenerator(generator []byte)
	DeleteSnapshotGenerator()

	WriteSnapshotDisabled()
	DeleteSnapshotDisabled()

	WriteSnapshotRecoveryNumber(number uint64)
	DeleteSnapshotRecoveryNumber()
}

type snapshotDBBatch struct {
	Batch
}

func (batch *snapshotDBBatch) WriteSnapshotRoot(root common.Hash) {
	writeSnapshotRoot(batch, root)
}

func (batch *snapshotDBBatch) DeleteSnapshotRoot() {
	deleteSnapshotRoot(batch)
}

func (batch *snapshotDBBatch) WriteAccountSnapshot(hash common.Hash, entry []byte) {
	writeAccountSnapshot(batch, hash, entry)
}

func (batch *snapshotDBBatch) DeleteAccountSnapshot(hash common.Hash) {
	deleteAccountSnapshot(batch, hash)
}

func (batch *snapshotDBBatch) WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte) {
	writeStorageSnapshot(batch, accountHash, storageHash, entry)
}

func (batch *snapshotDBBatch) DeleteStorageSnapshot(accountHash, storageHash common.Hash) {
	deleteStorageSnapshot(batch, accountHash, storageHash)
}

func (batch *snapshotDBBatch) WriteSnapshotJournal(journal []byte) {
	writeSnapshotJournal(batch, journal)
}

func (batch *snapshotDBBatch) DeleteSnapshotJournal() {
	deleteSnapshotJournal(batch)
}

func (batch *snapshotDBBatch) WriteSnapshotGenerator(generator []byte) {
	writeSnapshotGenerator(batch, generator)
}

func (batch *snapshotDBBatch) DeleteSnapshotGenerator() {
	deleteSnapshotGenerator(batch)
}

func (batch *snapshotDBBatch) WriteSnapshotDisabled() {
	writeSnapshotDisabled(batch)
}

func (batch *snapshotDBBatch) DeleteSnapshotDisabled() {
	deleteSnapshotDisabled(batch)
}

func (batch *snapshotDBBatch) WriteSnapshotRecoveryNumber(number uint64) {
	writeSnapshotRecoveryNumber(batch, number)
}

func (batch *snapshotDBBatch) DeleteSnapshotRecoveryNumber() {
	deleteSnapshotRecoveryNumber(batch)
}

func writeSnapshotRoot(db KeyValueWriter, root common.Hash) {
	if err := db.Put(snapshotRootKey, root[:]); err != nil {
		logger.Crit("Failed to store snapshot root", "err", err)
	}
}

func deleteSnapshotRoot(db KeyValueWriter) {
	if err := db.Delete(snapshotRootKey); err != nil {
		logger.Crit("Failed to remove snapshot root", "err", err)
	}
}

func writeAccountSnapshot(db KeyValueWriter, hash common.Hash, entry []byte) {
	if err := db.Put(AccountSnapshotKey(hash), entry); err != nil {
		logger.Crit("Failed to store account snapshot", "err", err)
	}
}

func deleteAccountSnapshot(db KeyValueWriter, hash common.Hash) {
	if err := db.Delete(AccountSnapshotKey(hash)); err != nil {
		logger.Crit("Failed to delete account snapshot", "err", err)
	}
}

func writeStorageSnapshot(db KeyValueWriter, accountHash, storageHash common.Hash, entry []byte) {
	if err := db.Put(StorageSnapshotKey(accountHash, storageHash), entry); err != nil {
		logger.Crit("Failed to store storage snapshot", "err", err)
	}
}

func deleteStorageSnapshot(db KeyValueWriter, accountHash, storageHash common.Hash) {
	if err := db.Delete(StorageSnapshotKey(accountHash, storageHash)); err != nil {
		logger.Crit("Failed to delete storage snapshot", "err", err)
	}
}

func writeSnapshotJournal(db KeyValueWriter, journal []byte) {
	if err := db.Put(snapshotJournalKey, journal); err != nil {
		logger.Crit("Failed to store snapshot journal", "err", err)
	}
}

func deleteSnapshotJournal(db KeyValueWriter) {
	if err := db.Delete(snapshotJournalKey); err != nil {
		logger.Crit("Failed to remove snapshot journal", "err", err)
	}
}

func writeSnapshotGenerator(db KeyValueWriter, generator []byte) {
	if err := db.Put(SnapshotGeneratorKey, generator); err != nil {
		logger.Crit("Failed to store snapshot generator", "err", err)
	}
}

func deleteSnapshotGenerator(db KeyValueWriter) {
	if err := db.Delete(SnapshotGeneratorKey); err != nil {
		logger.Crit("Failed to remove snapshot generator", "err", err)
	}
}

func writeSnapshotDisabled(db KeyValueWriter) {
	if err := db.Put(snapshotDisabledKey, []byte("42")); err != nil {
		logger.Crit("Failed to store snapshot disabled flag", "err", err)
	}
}

func deleteSnapshotDisabled(db KeyValueWriter) {
	if err := db.Delete(snapshotDisabledKey); err != nil {
		logger.Crit("Failed to remove snapshot disabled flag", "err", err)
	}
}

func writeSnapshotRecoveryNumber(db KeyValueWriter, number uint64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], number)
	if err := db.Put(snapshotRecoveryKey, buf[:]); err != nil {
		logger.Crit("Failed to store snapshot recovery number", "err", err)
	}
}

func deleteSnapshotRecoveryNumber(db KeyValueWriter) {
	if err := db.Delete(snapshotRecoveryKey); err != nil {
		logger.Crit("Failed to remove snapshot recovery number", "err", err)
	}
}
