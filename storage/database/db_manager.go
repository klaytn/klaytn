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
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/pkg/errors"
	"math/big"
	"path/filepath"
)

var logger = log.NewModuleLogger(log.StorageDatabase)

type DBManager interface {
	IsParallelDBWrite() bool

	Close()
	NewBatch(dbType DBEntryType) Batch
	GetMemDB() *MemDB
	GetDBConfig() *DBConfig

	// from accessors_chain.go
	ReadCanonicalHash(number uint64) common.Hash
	WriteCanonicalHash(hash common.Hash, number uint64)
	DeleteCanonicalHash(number uint64)

	ReadHeadHeaderHash() common.Hash
	WriteHeadHeaderHash(hash common.Hash)

	ReadHeadBlockHash() common.Hash
	WriteHeadBlockHash(hash common.Hash)

	ReadHeadFastBlockHash() common.Hash
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

	FindCommonAncestor(a, b *types.Header) *types.Header

	ReadIstanbulSnapshot(hash common.Hash) ([]byte, error)
	WriteIstanbulSnapshot(hash common.Hash, blob []byte) error

	WriteMerkleProof(key, value []byte)

	ReadCachedTrieNode(hash common.Hash) ([]byte, error)
	ReadCachedTrieNodePreimage(secureKey []byte) ([]byte, error)

	ReadStateTrieNode(key []byte) ([]byte, error)
	HasStateTrieNode(key []byte) (bool, error)

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

	ReadPreimage(hash common.Hash) []byte
	WritePreimages(number uint64, preimages map[common.Hash][]byte)

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
}

type DBEntryType uint8

const (
	headerDB DBEntryType = iota
	BodyDB
	ReceiptsDB
	StateTrieDB
	TxLookUpEntryDB
	MiscDB
	bridgeServiceDB

	// databaseEntryTypeSize should be the last item in this list!!
	databaseEntryTypeSize
)

var dbDirs = [databaseEntryTypeSize]string{
	"header",
	"body",
	"receipts",
	"statetrie",
	"txlookup",
	"misc",
	"bridgeservice",
}

// Sum of dbConfigRatio should be 100.
// Otherwise, logger.Crit will be called at checkDBEntryConfigRatio.
var dbConfigRatio = [databaseEntryTypeSize]int{
	6,  // headerDB
	21, // BodyDB
	21, // ReceiptsDB
	23, // StateTrieDB
	21, // TXLookUpEntryDB
	3,  // MiscDB
	5,  // bridgeServiceDB
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

// getDBEntryConfig returns a new DBConfig with original DBConfig and DBEntryType.
// It adjusts configuration according to the ratio specified in dbConfigRatio and dbDirs.
func getDBEntryConfig(originalDBC *DBConfig, i DBEntryType) *DBConfig {
	newDBC := *originalDBC
	ratio := dbConfigRatio[i]

	newDBC.LevelDBCacheSize = originalDBC.LevelDBCacheSize * ratio / 100
	newDBC.OpenFilesLimit = originalDBC.OpenFilesLimit * ratio / 100

	// Update dir to each Database specific directory.
	newDBC.Dir = filepath.Join(originalDBC.Dir, dbDirs[i])

	return &newDBC
}

type databaseManager struct {
	config *DBConfig
	dbs    []Database
	cm     *cacheManager
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
	Dir                    string
	DBType                 DBType
	Partitioned            bool
	NumStateTriePartitions uint
	ParallelDBWrite        bool
	OpenFilesLimit         int

	// LevelDB related configurations.
	LevelDBCacheSize   int // LevelDBCacheSize = BlockCacheCapacity + WriteBuffer
	LevelDBCompression LevelDBCompressionType
	LevelDBBufferPool  bool
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

// partitionedDatabaseDBManager returns DBManager which handles partitioned Database.
// Each Database will have its own separated Database.
func partitionedDatabaseDBManager(dbc *DBConfig) (DBManager, error) {
	dbm := newDatabaseManager(dbc)
	var db Database
	var err error
	for et := 0; et < int(databaseEntryTypeSize); et++ {
		entryType := DBEntryType(et)

		newDBC := getDBEntryConfig(dbc, entryType)

		if entryType == StateTrieDB && dbc.NumStateTriePartitions > 1 {
			db, err = newPartitionedDB(newDBC, entryType, dbc.NumStateTriePartitions)
		} else {
			db, err = newDatabase(newDBC, entryType)
		}

		if err != nil {
			logger.Crit("Failed while generating a partition of LevelDB", "partition", dbDirs[et], "err", err)
		}
		dbm.dbs[et] = db
		db.Meter(dbMetricPrefix + dbDirs[et] + "/") // Each partition collects metrics independently.
	}
	return dbm, nil
}

// newDatabase returns Database interface with given DBConfig.
func newDatabase(dbc *DBConfig, entryType DBEntryType) (Database, error) {
	switch dbc.DBType {
	case LevelDB:
		return NewLevelDB(dbc, entryType)
	case BadgerDB:
		return NewBadgerDB(dbc.Dir)
	case MemoryDB:
		return NewMemDB(), nil
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
// If Partitioned is true, each Database will have its own LevelDB.
// If not, each Database will share one common LevelDB.
func NewDBManager(dbc *DBConfig) DBManager {
	if !dbc.Partitioned {
		logger.Info("Non-partitioned database is used for persistent storage", "DBType", dbc.DBType)
		if dbm, err := singleDatabaseDBManager(dbc); err != nil {
			logger.Crit("Failed to create non-partitioned database", "DBType", dbc.DBType, "err", err)
		} else {
			return dbm
		}
	} else {
		checkDBEntryConfigRatio()
		logger.Info("Partitioned database is used for persistent storage", "DBType", dbc.DBType)
		if dbm, err := partitionedDatabaseDBManager(dbc); err != nil {
			logger.Crit("Failed to partitioned database", "DBType", dbc.DBType, "err", err)
		} else {
			return dbm
		}
	}
	logger.Crit("Must not reach here!")
	return nil
}

func (dbm *databaseManager) IsParallelDBWrite() bool {
	return dbm.config.ParallelDBWrite
}

func (dbm *databaseManager) NewBatch(dbEntryType DBEntryType) Batch {
	return dbm.getDatabase(dbEntryType).NewBatch()
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
		return dbm.dbs[headerDB]
	} else {
		return dbm.dbs[dbEntryType]
	}
}

func (dbm *databaseManager) Close() {
	// If not partitioned, only close the first database.
	if !dbm.config.Partitioned {
		dbm.dbs[0].Close()
		return
	}

	// If partitioned, close all databases.
	for _, db := range dbm.dbs {
		db.Close()
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

// WriteHeadBlockHash stores the head block's hash.
func (dbm *databaseManager) WriteHeadBlockHash(hash common.Hash) {
	db := dbm.getDatabase(headerDB)
	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		logger.Crit("Failed to store last block's hash", "err", err)
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

// WriteHeadFastBlockHash stores the hash of the current fast-sync head block.
func (dbm *databaseManager) WriteHeadFastBlockHash(hash common.Hash) {
	db := dbm.getDatabase(headerDB)
	if err := db.Put(headFastBlockKey, hash.Bytes()); err != nil {
		logger.Crit("Failed to store last fast block's hash", "err", err)
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
		encoded = encodeBlockNumber(number)
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
	db := dbm.getDatabase(BodyDB)
	if err := db.Put(blockBodyKey(number, hash), rlp); err != nil {
		logger.Crit("Failed to store block body", "err", err)
	}
	if common.WriteThroughCaching {
		dbm.cm.writeBodyRLPCache(hash, rlp)
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
	db := dbm.getDatabase(ReceiptsDB)
	// When putReceiptsToPutter is called from WriteReceipts, txReceipt is cached.
	dbm.putReceiptsToPutter(db, hash, number, receipts, true)
	if common.WriteThroughCaching {
		// TODO-Klaytn goroutine for performance
		dbm.cm.writeBlockReceiptsCache(hash, receipts)
	}
}

func (dbm *databaseManager) PutReceiptsToBatch(batch Batch, hash common.Hash, number uint64, receipts types.Receipts) {
	// When putReceiptsToPutter is called from PutReceiptsToBatch, txReceipt is not cached.
	dbm.putReceiptsToPutter(batch, hash, number, receipts, false)
}

func (dbm *databaseManager) putReceiptsToPutter(putter Putter, hash common.Hash, number uint64, receipts types.Receipts, addToCache bool) {
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
	dbm.WriteBody(block.Hash(), block.NumberU64(), block.Body())
	dbm.WriteHeader(block.Header())

	if common.WriteThroughCaching {
		dbm.cm.writeBodyCache(block.Hash(), block.Body())
		dbm.cm.blockCache.Add(block.Hash(), block)
	}
}

func (dbm *databaseManager) DeleteBlock(hash common.Hash, number uint64) {
	dbm.DeleteReceipts(hash, number)
	dbm.DeleteHeader(hash, number)
	dbm.DeleteBody(hash, number)
	dbm.DeleteTd(hash, number)
	dbm.cm.deleteBlockCache(hash)
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

// Merkle Proof operation.
func (dbm *databaseManager) WriteMerkleProof(key, value []byte) {
	db := dbm.getDatabase(MiscDB)
	if err := db.Put(key, value); err != nil {
		logger.Crit("Failed to write merkle proof", "err", err)
	}
}

// Cached Trie Node operation.
func (dbm *databaseManager) ReadCachedTrieNode(hash common.Hash) ([]byte, error) {
	db := dbm.getDatabase(StateTrieDB)
	return db.Get(hash[:])
}

// Cached Trie Node Preimage operation.
func (dbm *databaseManager) ReadCachedTrieNodePreimage(secureKey []byte) ([]byte, error) {
	db := dbm.getDatabase(StateTrieDB)
	return db.Get(secureKey)
}

// State Trie Related operations.
func (dbm *databaseManager) ReadStateTrieNode(key []byte) ([]byte, error) {
	db := dbm.getDatabase(StateTrieDB)
	return db.Get(key)
}

func (dbm *databaseManager) HasStateTrieNode(key []byte) (bool, error) {
	val, err := dbm.ReadStateTrieNode(key)
	if val == nil || err != nil {
		return false, err
	}
	return true, nil
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

func putTxLookupEntriesToPutter(putter Putter, block *types.Block) {
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
	return dbm.NewBatch(MiscDB)
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

// ReadPreimage retrieves a single preimage of the provided hash.
func (dbm *databaseManager) ReadPreimage(hash common.Hash) []byte {
	db := dbm.getDatabase(StateTrieDB)
	data, _ := db.Get(preimageKey(hash))
	return data
}

// WritePreimages writes the provided set of preimages to the database. `number` is the
// current block number, and is used for debug messages only.
func (dbm *databaseManager) WritePreimages(number uint64, preimages map[common.Hash][]byte) {
	batch := dbm.getDatabase(StateTrieDB).NewBatch()
	for hash, preimage := range preimages {
		if err := batch.Put(preimageKey(hash), preimage); err != nil {
			logger.Crit("Failed to store trie preimage", "err", err)
		}
	}
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to batch write trie preimage", "err", err, "blockNumber", number)
	}
	preimageCounter.Inc(int64(len(preimages)))
	preimageHitCounter.Inc(int64(len(preimages)))
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
	if err := db.Put(key, encodeBlockNumber(blockNum)); err != nil {
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
	if err := db.Put(key, encodeBlockNumber(blockNum)); err != nil {
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
		return err
	}
	return db.Put(governanceKey(num), b)
}

func (dbm *databaseManager) WriteGovernanceIdx(num uint64) error {
	db := dbm.getDatabase(MiscDB)
	newSlice := make([]uint64, 0)

	if data, err := db.Get(governanceHistoryKey); err == nil {
		if err = json.Unmarshal(data, &newSlice); err != nil {
			return err
		}
	}
	newSlice = append(newSlice, num)

	data, err := json.Marshal(newSlice)
	if err != nil {
		return err
	}
	return db.Put(governanceHistoryKey, data)
}

func (dbm *databaseManager) ReadGovernance(num uint64) (map[string]interface{}, error) {
	db := dbm.getDatabase(MiscDB)

	if data, err := db.Get(governanceKey(num)); err != nil {
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
	var minimum = num - (num % epoch)
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
