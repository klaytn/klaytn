// Copyright 2019 The klaytn Authors
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
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

var (
	dbManagers  []DBManager
	dbConfigs   = make([]*DBConfig, 0, len(baseConfigs)*3)
	baseConfigs = []*DBConfig{
		{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 1, ParallelDBWrite: false},
		{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 1, ParallelDBWrite: true},
		{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 4, ParallelDBWrite: false},
		{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 4, ParallelDBWrite: true},

		{DBType: LevelDB, SingleDB: true, NumStateTrieShards: 1, ParallelDBWrite: false},
		{DBType: LevelDB, SingleDB: true, NumStateTrieShards: 1, ParallelDBWrite: true},
		{DBType: LevelDB, SingleDB: true, NumStateTrieShards: 4, ParallelDBWrite: false},
		{DBType: LevelDB, SingleDB: true, NumStateTrieShards: 4, ParallelDBWrite: true},
	}
)

var (
	num1 = uint64(20190815)
	num2 = uint64(20199999)
	num3 = uint64(12345678)
	num4 = uint64(87654321)
)

var (
	hash1 = common.HexToHash("1341655") // 20190805 in hexadecimal
	hash2 = common.HexToHash("1343A3F") // 20199999 in hexadecimal
	hash3 = common.HexToHash("BC614E")  // 12345678 in hexadecimal
	hash4 = common.HexToHash("5397FB1") // 87654321 in hexadecimal
)

var (
	key    *ecdsa.PrivateKey
	addr   common.Address
	signer types.Signer
)

func init() {
	GetOpenFilesLimit()

	key, _ = crypto.GenerateKey()
	addr = crypto.PubkeyToAddress(key.PublicKey)
	signer = types.LatestSignerForChainID(big.NewInt(18))

	for _, bc := range baseConfigs {
		badgerConfig := *bc
		badgerConfig.DBType = BadgerDB
		memoryConfig := *bc
		memoryConfig.DBType = MemoryDB

		dbConfigs = append(dbConfigs, bc)
		dbConfigs = append(dbConfigs, &badgerConfig)
		dbConfigs = append(dbConfigs, &memoryConfig)
	}

	dbManagers = createDBManagers(dbConfigs)
}

// createDBManagers generates a list of DBManagers to test various combinations of DBConfig.
func createDBManagers(configs []*DBConfig) []DBManager {
	dbManagers := make([]DBManager, 0, len(configs))

	for _, c := range configs {
		c.Dir, _ = ioutil.TempDir(os.TempDir(), "test-db-manager")
		dbManagers = append(dbManagers, NewDBManager(c))
	}

	return dbManagers
}

// TestDBManager_IsParallelDBWrite compares the return value of IsParallelDBWrite with the value in the config.
func TestDBManager_IsParallelDBWrite(t *testing.T) {
	for i, dbm := range dbManagers {
		c := dbConfigs[i]
		assert.Equal(t, c.ParallelDBWrite, dbm.IsParallelDBWrite())
	}
}

// TestDBManager_CanonicalHash tests read, write and delete operations of canonical hash.
func TestDBManager_CanonicalHash(t *testing.T) {
	for _, dbm := range dbManagers {
		// 1. Read from empty database, shouldn't be found.
		assert.Equal(t, common.Hash{}, dbm.ReadCanonicalHash(0))
		assert.Equal(t, common.Hash{}, dbm.ReadCanonicalHash(100))

		// 2. Write a row to the database.
		dbm.WriteCanonicalHash(hash1, num1)

		// 3. Read from the database, only written key-value pair should be found.
		assert.Equal(t, common.Hash{}, dbm.ReadCanonicalHash(0))
		assert.Equal(t, common.Hash{}, dbm.ReadCanonicalHash(100))
		assert.Equal(t, hash1, dbm.ReadCanonicalHash(num1)) // should be found

		// 4. Overwrite existing key with different value, value should be changed.
		hash2 := common.HexToHash("1343A3F")                // 20199999 in hexadecimal
		dbm.WriteCanonicalHash(hash2, num1)                 // overwrite hash1 by hash2 with same key
		assert.Equal(t, hash2, dbm.ReadCanonicalHash(num1)) // should be hash2

		// 5. Delete non-existing value.
		dbm.DeleteCanonicalHash(num2)
		assert.Equal(t, hash2, dbm.ReadCanonicalHash(num1)) // should be hash2, not deleted

		// 6. Delete existing value.
		dbm.DeleteCanonicalHash(num1)
		assert.Equal(t, common.Hash{}, dbm.ReadCanonicalHash(num1)) // shouldn't be found
	}
}

// TestDBManager_HeadHeaderHash tests read and write operations of head header hash.
func TestDBManager_HeadHeaderHash(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Equal(t, common.Hash{}, dbm.ReadHeadHeaderHash())

		dbm.WriteHeadHeaderHash(hash1)
		assert.Equal(t, hash1, dbm.ReadHeadHeaderHash())

		dbm.WriteHeadHeaderHash(hash2)
		assert.Equal(t, hash2, dbm.ReadHeadHeaderHash())
	}
}

// TestDBManager_HeadBlockHash tests read and write operations of head block hash.
func TestDBManager_HeadBlockHash(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Equal(t, common.Hash{}, dbm.ReadHeadBlockHash())

		dbm.WriteHeadBlockHash(hash1)
		assert.Equal(t, hash1, dbm.ReadHeadBlockHash())

		dbm.WriteHeadBlockHash(hash2)
		assert.Equal(t, hash2, dbm.ReadHeadBlockHash())
	}
}

// TestDBManager_HeadFastBlockHash tests read and write operations of head fast block hash.
func TestDBManager_HeadFastBlockHash(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Equal(t, common.Hash{}, dbm.ReadHeadFastBlockHash())

		dbm.WriteHeadFastBlockHash(hash1)
		assert.Equal(t, hash1, dbm.ReadHeadFastBlockHash())

		dbm.WriteHeadFastBlockHash(hash2)
		assert.Equal(t, hash2, dbm.ReadHeadFastBlockHash())
	}
}

// TestDBManager_FastTrieProgress tests read and write operations of fast trie progress.
func TestDBManager_FastTrieProgress(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Equal(t, uint64(0), dbm.ReadFastTrieProgress())

		dbm.WriteFastTrieProgress(num1)
		assert.Equal(t, num1, dbm.ReadFastTrieProgress())

		dbm.WriteFastTrieProgress(num2)
		assert.Equal(t, num2, dbm.ReadFastTrieProgress())
	}
}

// TestDBManager_Header tests read, write and delete operations of blockchain headers.
func TestDBManager_Header(t *testing.T) {
	header := &types.Header{Number: big.NewInt(int64(num1))}
	headerHash := header.Hash()

	encodedHeader, err := rlp.EncodeToBytes(header)
	if err != nil {
		t.Fatal("Failed to encode header!", "err", err)
	}

	for _, dbm := range dbManagers {
		assert.False(t, dbm.HasHeader(headerHash, num1))
		assert.Nil(t, dbm.ReadHeader(headerHash, num1))
		assert.Nil(t, dbm.ReadHeaderNumber(headerHash))

		dbm.WriteHeader(header)

		assert.True(t, dbm.HasHeader(headerHash, num1))
		assert.Equal(t, header, dbm.ReadHeader(headerHash, num1))
		assert.Equal(t, rlp.RawValue(encodedHeader), dbm.ReadHeaderRLP(headerHash, num1))
		assert.Equal(t, num1, *dbm.ReadHeaderNumber(headerHash))

		dbm.DeleteHeader(headerHash, num1)

		assert.False(t, dbm.HasHeader(headerHash, num1))
		assert.Nil(t, dbm.ReadHeader(headerHash, num1))
		assert.Nil(t, dbm.ReadHeaderNumber(headerHash))
	}
}

// TestDBManager_Body tests read, write and delete operations of blockchain bodies.
func TestDBManager_Body(t *testing.T) {
	body := &types.Body{Transactions: types.Transactions{}}
	encodedBody, err := rlp.EncodeToBytes(body)
	if err != nil {
		t.Fatal("Failed to encode body!", "err", err)
	}

	for _, dbm := range dbManagers {
		assert.False(t, dbm.HasBody(hash1, num1))
		assert.Nil(t, dbm.ReadBody(hash1, num1))
		assert.Nil(t, dbm.ReadBodyInCache(hash1))
		assert.Nil(t, dbm.ReadBodyRLP(hash1, num1))
		assert.Nil(t, dbm.ReadBodyRLPByHash(hash1))

		dbm.WriteBody(hash1, num1, body)

		assert.True(t, dbm.HasBody(hash1, num1))
		assert.Equal(t, body, dbm.ReadBody(hash1, num1))
		assert.Equal(t, body, dbm.ReadBodyInCache(hash1))
		assert.Equal(t, rlp.RawValue(encodedBody), dbm.ReadBodyRLP(hash1, num1))
		assert.Equal(t, rlp.RawValue(encodedBody), dbm.ReadBodyRLPByHash(hash1))

		dbm.DeleteBody(hash1, num1)

		assert.False(t, dbm.HasBody(hash1, num1))
		assert.Nil(t, dbm.ReadBody(hash1, num1))
		assert.Nil(t, dbm.ReadBodyInCache(hash1))
		assert.Nil(t, dbm.ReadBodyRLP(hash1, num1))
		assert.Nil(t, dbm.ReadBodyRLPByHash(hash1))

		dbm.WriteBodyRLP(hash1, num1, encodedBody)

		assert.True(t, dbm.HasBody(hash1, num1))
		assert.Equal(t, body, dbm.ReadBody(hash1, num1))
		assert.Equal(t, body, dbm.ReadBodyInCache(hash1))
		assert.Equal(t, rlp.RawValue(encodedBody), dbm.ReadBodyRLP(hash1, num1))
		assert.Equal(t, rlp.RawValue(encodedBody), dbm.ReadBodyRLPByHash(hash1))

		bodyBatch := dbm.NewBatch(BodyDB)
		dbm.PutBodyToBatch(bodyBatch, hash2, num2, body)
		if err := bodyBatch.Write(); err != nil {
			t.Fatal("Failed to write batch!", "err", err)
		}

		assert.True(t, dbm.HasBody(hash2, num2))
		assert.Equal(t, body, dbm.ReadBody(hash2, num2))
		assert.Equal(t, body, dbm.ReadBodyInCache(hash2))
		assert.Equal(t, rlp.RawValue(encodedBody), dbm.ReadBodyRLP(hash2, num2))
		assert.Equal(t, rlp.RawValue(encodedBody), dbm.ReadBodyRLPByHash(hash2))
	}
}

// TestDBManager_Td tests read, write and delete operations of blockchain headers' total difficulty.
func TestDBManager_Td(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Nil(t, dbm.ReadTd(hash1, num1))

		dbm.WriteTd(hash1, num1, big.NewInt(12345))
		assert.Equal(t, big.NewInt(12345), dbm.ReadTd(hash1, num1))

		dbm.WriteTd(hash1, num1, big.NewInt(54321))
		assert.Equal(t, big.NewInt(54321), dbm.ReadTd(hash1, num1))

		dbm.DeleteTd(hash1, num1)
		assert.Nil(t, dbm.ReadTd(hash1, num1))
	}
}

// TestDBManager_Receipts read, write and delete operations of blockchain receipts.
func TestDBManager_Receipts(t *testing.T) {
	header := &types.Header{Number: big.NewInt(int64(num1))}
	headerHash := header.Hash()
	receipts := types.Receipts{genReceipt(111)}

	for _, dbm := range dbManagers {
		assert.Nil(t, dbm.ReadReceipts(headerHash, num1))
		assert.Nil(t, dbm.ReadReceiptsByBlockHash(headerHash))

		dbm.WriteReceipts(headerHash, num1, receipts)
		dbm.WriteHeader(header)

		assert.Equal(t, receipts, dbm.ReadReceipts(headerHash, num1))
		assert.Equal(t, receipts, dbm.ReadReceiptsByBlockHash(headerHash))

		dbm.DeleteReceipts(headerHash, num1)

		assert.Nil(t, dbm.ReadReceipts(headerHash, num1))
		assert.Nil(t, dbm.ReadReceiptsByBlockHash(headerHash))
	}
}

// TestDBManager_Block read, write and delete operations of blockchain blocks.
func TestDBManager_Block(t *testing.T) {
	header := &types.Header{Number: big.NewInt(int64(num1))}
	headerHash := header.Hash()
	block := types.NewBlockWithHeader(header)

	for _, dbm := range dbManagers {
		assert.False(t, dbm.HasBlock(headerHash, num1))
		assert.Nil(t, dbm.ReadBlock(headerHash, num1))
		assert.Nil(t, dbm.ReadBlockByHash(headerHash))
		assert.Nil(t, dbm.ReadBlockByNumber(num1))

		dbm.WriteBlock(block)
		dbm.WriteCanonicalHash(headerHash, num1)

		assert.True(t, dbm.HasBlock(headerHash, num1))

		blockFromDB1 := dbm.ReadBlock(headerHash, num1)
		blockFromDB2 := dbm.ReadBlockByHash(headerHash)
		blockFromDB3 := dbm.ReadBlockByNumber(num1)

		assert.Equal(t, headerHash, blockFromDB1.Hash())
		assert.Equal(t, headerHash, blockFromDB2.Hash())
		assert.Equal(t, headerHash, blockFromDB3.Hash())

		dbm.DeleteBlock(headerHash, num1)
		dbm.DeleteCanonicalHash(num1)

		assert.False(t, dbm.HasBlock(headerHash, num1))
		assert.Nil(t, dbm.ReadBlock(headerHash, num1))
		assert.Nil(t, dbm.ReadBlockByHash(headerHash))
		assert.Nil(t, dbm.ReadBlockByNumber(num1))
	}
}

// TestDBManager_IstanbulSnapshot tests read and write operations of istanbul snapshots.
func TestDBManager_IstanbulSnapshot(t *testing.T) {
	for _, dbm := range dbManagers {
		snapshot, _ := dbm.ReadIstanbulSnapshot(hash3)
		assert.Nil(t, snapshot)

		dbm.WriteIstanbulSnapshot(hash3, hash2[:])
		snapshot, _ = dbm.ReadIstanbulSnapshot(hash3)
		assert.Equal(t, hash2[:], snapshot)

		dbm.WriteIstanbulSnapshot(hash3, hash1[:])
		snapshot, _ = dbm.ReadIstanbulSnapshot(hash3)
		assert.Equal(t, hash1[:], snapshot)
	}
}

// TestDBManager_TrieNode tests read and write operations of state trie nodes.
func TestDBManager_TrieNode(t *testing.T) {
	for _, dbm := range dbManagers {
		cachedNode, _ := dbm.ReadCachedTrieNode(hash1)
		assert.Nil(t, cachedNode)
		hasStateTrieNode, _ := dbm.HasStateTrieNode(hash1[:])
		assert.False(t, hasStateTrieNode)

		batch := dbm.NewBatch(StateTrieDB)
		if err := batch.Put(hash1[:], hash2[:]); err != nil {
			t.Fatal("Failed putting a row into the batch", "err", err)
		}
		if _, err := WriteBatches(batch); err != nil {
			t.Fatal("Failed writing batch", "err", err)
		}

		cachedNode, _ = dbm.ReadCachedTrieNode(hash1)
		assert.Equal(t, hash2[:], cachedNode)

		if err := batch.Put(hash1[:], hash1[:]); err != nil {
			t.Fatal("Failed putting a row into the batch", "err", err)
		}
		if _, err := WriteBatches(batch); err != nil {
			t.Fatal("Failed writing batch", "err", err)
		}

		cachedNode, _ = dbm.ReadCachedTrieNode(hash1)
		assert.Equal(t, hash1[:], cachedNode)

		stateTrieNode, _ := dbm.ReadStateTrieNode(hash1[:])
		assert.Equal(t, hash1[:], stateTrieNode)

		hasStateTrieNode, _ = dbm.HasStateTrieNode(hash1[:])
		assert.True(t, hasStateTrieNode)

		if dbm.IsSingle() {
			continue
		}
		err := dbm.CreateMigrationDBAndSetStatus(123)
		assert.NoError(t, err)

		cachedNode, _ = dbm.ReadCachedTrieNode(hash1)
		oldCachedNode, _ := dbm.ReadCachedTrieNodeFromOld(hash1)
		assert.Equal(t, hash1[:], cachedNode)
		assert.Equal(t, hash1[:], oldCachedNode)

		stateTrieNode, _ = dbm.ReadStateTrieNode(hash1[:])
		oldStateTrieNode, _ := dbm.ReadStateTrieNodeFromOld(hash1[:])
		assert.Equal(t, hash1[:], stateTrieNode)
		assert.Equal(t, hash1[:], oldStateTrieNode)

		hasStateTrieNode, _ = dbm.HasStateTrieNode(hash1[:])
		hasOldStateTrieNode, _ := dbm.HasStateTrieNodeFromOld(hash1[:])
		assert.True(t, hasStateTrieNode)
		assert.True(t, hasOldStateTrieNode)

		batch = dbm.NewBatch(StateTrieDB)
		if err := batch.Put(hash2[:], hash2[:]); err != nil {
			t.Fatal("Failed putting a row into the batch", "err", err)
		}
		if _, err := WriteBatches(batch); err != nil {
			t.Fatal("Failed writing batch", "err", err)
		}

		cachedNode, _ = dbm.ReadCachedTrieNode(hash2)
		oldCachedNode, _ = dbm.ReadCachedTrieNodeFromOld(hash2)
		assert.Equal(t, hash2[:], cachedNode)
		assert.Equal(t, hash2[:], oldCachedNode)

		stateTrieNode, _ = dbm.ReadStateTrieNode(hash2[:])
		oldStateTrieNode, _ = dbm.ReadStateTrieNodeFromOld(hash2[:])
		assert.Equal(t, hash2[:], stateTrieNode)
		assert.Equal(t, hash2[:], oldStateTrieNode)

		hasStateTrieNode, _ = dbm.HasStateTrieNode(hash2[:])
		hasOldStateTrieNode, _ = dbm.HasStateTrieNodeFromOld(hash2[:])
		assert.True(t, hasStateTrieNode)
		assert.True(t, hasOldStateTrieNode)

		dbm.FinishStateMigration(true)
	}
}

// TestDBManager_TxLookupEntry tests read, write and delete operations of TxLookupEntries.
func TestDBManager_TxLookupEntry(t *testing.T) {
	tx, err := genTransaction(num1)
	assert.NoError(t, err, "Failed to generate a transaction")

	body := &types.Body{Transactions: types.Transactions{tx}}
	for _, dbm := range dbManagers {
		blockHash, blockIndex, entryIndex := dbm.ReadTxLookupEntry(tx.Hash())
		assert.Equal(t, common.Hash{}, blockHash)
		assert.Equal(t, uint64(0), blockIndex)
		assert.Equal(t, uint64(0), entryIndex)

		header := &types.Header{Number: big.NewInt(int64(num1))}
		block := types.NewBlockWithHeader(header)
		block = block.WithBody(body.Transactions)

		dbm.WriteTxLookupEntries(block)

		blockHash, blockIndex, entryIndex = dbm.ReadTxLookupEntry(tx.Hash())
		assert.Equal(t, block.Hash(), blockHash)
		assert.Equal(t, block.NumberU64(), blockIndex)
		assert.Equal(t, uint64(0), entryIndex)

		dbm.DeleteTxLookupEntry(tx.Hash())

		blockHash, blockIndex, entryIndex = dbm.ReadTxLookupEntry(tx.Hash())
		assert.Equal(t, common.Hash{}, blockHash)
		assert.Equal(t, uint64(0), blockIndex)
		assert.Equal(t, uint64(0), entryIndex)

		dbm.WriteAndCacheTxLookupEntries(block)
		blockHash, blockIndex, entryIndex = dbm.ReadTxLookupEntry(tx.Hash())
		assert.Equal(t, block.Hash(), blockHash)
		assert.Equal(t, block.NumberU64(), blockIndex)
		assert.Equal(t, uint64(0), entryIndex)

		batch := dbm.NewSenderTxHashToTxHashBatch()
		if err := dbm.PutSenderTxHashToTxHashToBatch(batch, hash1, hash2); err != nil {
			t.Fatal("Failed while calling PutSenderTxHashToTxHashToBatch", "err", err)
		}

		if err := batch.Write(); err != nil {
			t.Fatal("Failed writing SenderTxHashToTxHashToBatch", "err", err)
		}

		assert.Equal(t, hash2, dbm.ReadTxHashFromSenderTxHash(hash1))
	}
}

// TestDBManager_BloomBits tests read, write and delete operations of bloom bits
func TestDBManager_BloomBits(t *testing.T) {
	for _, dbm := range dbManagers {
		hash1 := common.HexToHash("123456")
		hash2 := common.HexToHash("654321")

		sh, _ := dbm.ReadBloomBits(hash1[:])
		assert.Nil(t, sh)

		err := dbm.WriteBloomBits(hash1[:], hash1[:])
		if err != nil {
			t.Fatal("Failed to write bloom bits", "err", err)
		}

		sh, err = dbm.ReadBloomBits(hash1[:])
		if err != nil {
			t.Fatal("Failed to read bloom bits", "err", err)
		}
		assert.Equal(t, hash1[:], sh)

		err = dbm.WriteBloomBits(hash1[:], hash2[:])
		if err != nil {
			t.Fatal("Failed to write bloom bits", "err", err)
		}

		sh, err = dbm.ReadBloomBits(hash1[:])
		if err != nil {
			t.Fatal("Failed to read bloom bits", "err", err)
		}
		assert.Equal(t, hash2[:], sh)
	}
}

// TestDBManager_Sections tests read, write and delete operations of ValidSections and SectionHead.
func TestDBManager_Sections(t *testing.T) {
	for _, dbm := range dbManagers {
		// ValidSections
		vs, _ := dbm.ReadValidSections()
		assert.Nil(t, vs)

		dbm.WriteValidSections(hash1[:])

		vs, _ = dbm.ReadValidSections()
		assert.Equal(t, hash1[:], vs)

		// SectionHead
		sh, _ := dbm.ReadSectionHead(hash1[:])
		assert.Nil(t, sh)

		dbm.WriteSectionHead(hash1[:], hash1)

		sh, _ = dbm.ReadSectionHead(hash1[:])
		assert.Equal(t, hash1[:], sh)

		dbm.DeleteSectionHead(hash1[:])

		sh, _ = dbm.ReadSectionHead(hash1[:])
		assert.Nil(t, sh)
	}
}

// TestDBManager_DatabaseVersion tests read/write operations of database version.
func TestDBManager_DatabaseVersion(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Nil(t, dbm.ReadDatabaseVersion())

		dbm.WriteDatabaseVersion(uint64(1))
		assert.NotNil(t, dbm.ReadDatabaseVersion())
		assert.Equal(t, uint64(1), *dbm.ReadDatabaseVersion())

		dbm.WriteDatabaseVersion(uint64(2))
		assert.NotNil(t, dbm.ReadDatabaseVersion())
		assert.Equal(t, uint64(2), *dbm.ReadDatabaseVersion())
	}
}

// TestDBManager_ChainConfig tests read/write operations of chain configuration.
func TestDBManager_ChainConfig(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Nil(t, nil, dbm.ReadChainConfig(hash1))

		cc1 := &params.ChainConfig{UnitPrice: 12345}
		cc2 := &params.ChainConfig{UnitPrice: 54321}

		dbm.WriteChainConfig(hash1, cc1)
		assert.Equal(t, cc1, dbm.ReadChainConfig(hash1))
		assert.NotEqual(t, cc2, dbm.ReadChainConfig(hash1))

		dbm.WriteChainConfig(hash1, cc2)
		assert.NotEqual(t, cc1, dbm.ReadChainConfig(hash1))
		assert.Equal(t, cc2, dbm.ReadChainConfig(hash1))
	}
}

// TestDBManager_Preimage tests read/write operations of preimages.
func TestDBManager_Preimage(t *testing.T) {
	for _, dbm := range dbManagers {
		assert.Nil(t, nil, dbm.ReadPreimage(hash1))

		preimages1 := map[common.Hash][]byte{hash1: hash2[:], hash2: hash1[:]}
		dbm.WritePreimages(num1, preimages1)

		assert.Equal(t, hash2[:], dbm.ReadPreimage(hash1))
		assert.Equal(t, hash1[:], dbm.ReadPreimage(hash2))

		preimages2 := map[common.Hash][]byte{hash1: hash1[:], hash2: hash2[:]}
		dbm.WritePreimages(num1, preimages2)

		assert.Equal(t, hash1[:], dbm.ReadPreimage(hash1))
		assert.Equal(t, hash2[:], dbm.ReadPreimage(hash2))
	}
}

// TestDBManager_ParentChain tests service chain related database operations, used in the parent chain.
func TestDBManager_ParentChain(t *testing.T) {
	for _, dbm := range dbManagers {
		// 1. Read/Write SerivceChainTxHash
		assert.Equal(t, common.Hash{}, dbm.ConvertChildChainBlockHashToParentChainTxHash(hash1))

		dbm.WriteChildChainTxHash(hash1, hash1)
		assert.Equal(t, hash1, dbm.ConvertChildChainBlockHashToParentChainTxHash(hash1))

		dbm.WriteChildChainTxHash(hash1, hash2)
		assert.Equal(t, hash2, dbm.ConvertChildChainBlockHashToParentChainTxHash(hash1))

		// 2. Read/Write LastIndexedBlockNumber
		assert.Equal(t, uint64(0), dbm.GetLastIndexedBlockNumber())

		dbm.WriteLastIndexedBlockNumber(num1)
		assert.Equal(t, num1, dbm.GetLastIndexedBlockNumber())

		dbm.WriteLastIndexedBlockNumber(num2)
		assert.Equal(t, num2, dbm.GetLastIndexedBlockNumber())
	}
}

// TestDBManager_ChildChain tests service chain related database operations, used in the child chain.
func TestDBManager_ChildChain(t *testing.T) {
	for _, dbm := range dbManagers {
		// 1. Read/Write AnchoredBlockNumber
		assert.Equal(t, uint64(0), dbm.ReadAnchoredBlockNumber())

		dbm.WriteAnchoredBlockNumber(num1)
		assert.Equal(t, num1, dbm.ReadAnchoredBlockNumber())

		dbm.WriteAnchoredBlockNumber(num2)
		assert.Equal(t, num2, dbm.ReadAnchoredBlockNumber())

		// 2. Read/Write ReceiptFromParentChain
		// TODO-Klaytn-Database Implement this!

		// 3. Read/Write HandleTxHashFromRequestTxHash
		assert.Equal(t, common.Hash{}, dbm.ReadHandleTxHashFromRequestTxHash(hash1))

		dbm.WriteHandleTxHashFromRequestTxHash(hash1, hash1)
		assert.Equal(t, hash1, dbm.ReadHandleTxHashFromRequestTxHash(hash1))

		dbm.WriteHandleTxHashFromRequestTxHash(hash1, hash2)
		assert.Equal(t, hash2, dbm.ReadHandleTxHashFromRequestTxHash(hash1))
	}
}

// TestDBManager_CliqueSnapshot tests read and write operations of clique snapshots.
func TestDBManager_CliqueSnapshot(t *testing.T) {
	for _, dbm := range dbManagers {
		data, err := dbm.ReadCliqueSnapshot(hash1)
		assert.NotNil(t, err)
		assert.Nil(t, data)

		err = dbm.WriteCliqueSnapshot(hash1, hash1[:])
		assert.Nil(t, err)

		data, _ = dbm.ReadCliqueSnapshot(hash1)
		assert.Equal(t, hash1[:], data)

		err = dbm.WriteCliqueSnapshot(hash1, hash2[:])
		assert.Nil(t, err)

		data, _ = dbm.ReadCliqueSnapshot(hash1)
		assert.Equal(t, hash2[:], data)
	}
}

func TestDBManager_Governance(t *testing.T) {
	// TODO-Klaytn-Database Implement this!
}

func TestDatabaseManager_CreateMigrationDBAndSetStatus(t *testing.T) {
	for i, dbm := range dbManagers {
		if dbConfigs[i].DBType == MemoryDB {
			continue
		}

		// check if migration fails on single DB
		if dbm.IsSingle() {
			migrationBlockNum := uint64(12345)

			// check if not in migration
			assert.False(t, dbManagers[i].InMigration(), "migration status should be not set before testing")

			// check if create migration fails
			err := dbm.CreateMigrationDBAndSetStatus(migrationBlockNum)
			assert.Error(t, err, "error expected on single DB") // expect error

			continue
		}

		// check if migration fails when in migration
		{
			migrationBlockNum := uint64(34567)

			// check if not in migration
			assert.False(t, dbManagers[i].InMigration(), "migration status should be not set before testing")

			// set migration status
			dbm.setStateTrieMigrationStatus(migrationBlockNum)
			assert.True(t, dbManagers[i].InMigration())

			// check if create migration fails
			err := dbm.CreateMigrationDBAndSetStatus(migrationBlockNum)
			assert.Error(t, err, "error expected when in migration state") // expect error

			// reset migration status for next test
			dbm.setStateTrieMigrationStatus(0)
		}

		// check if CreateMigrationDBAndSetStatus works as expected
		{
			migrationBlockNum := uint64(56789)

			// check if not in migration state
			assert.False(t, dbManagers[i].InMigration(), "migration status should be not set before testing")

			err := dbm.CreateMigrationDBAndSetStatus(migrationBlockNum)
			assert.NoError(t, err)

			// check if in migration state
			assert.True(t, dbm.InMigration())

			// check migration DB path in MiscDB
			migrationDBPathKey := append(databaseDirPrefix, common.Int64ToByteBigEndian(uint64(StateTrieMigrationDB))...)
			fetchedMigrationPath, err := dbm.getDatabase(MiscDB).Get(migrationDBPathKey)
			assert.NoError(t, err)
			expectedMigrationPath := "statetrie_migrated_" + strconv.FormatUint(migrationBlockNum, 10)
			assert.Equal(t, expectedMigrationPath, string(fetchedMigrationPath))

			// check block number in MiscDB
			fetchedBlockNum, err := dbm.getDatabase(MiscDB).Get(migrationStatusKey)
			assert.NoError(t, err)
			assert.Equal(t, common.Int64ToByteBigEndian(migrationBlockNum), fetchedBlockNum)

			// reset migration status for next test
			dbm.FinishStateMigration(false) // migration fail
		}
	}
}

func TestDatabaseManager_FinishStateMigration(t *testing.T) {
	for i, dbm := range dbManagers {
		if dbm.IsSingle() || dbConfigs[i].DBType == MemoryDB {
			continue
		}

		migrationBlockNum := uint64(12345)
		migrationBlockNum2 := uint64(23456)

		// check status in miscDB on state migration failure
		{
			// check if not in migration state
			assert.False(t, dbManagers[i].InMigration(), "migration status should be not set before testing")
			// fetch state trie db name before migration
			initialDirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 1, len(initialDirNames), "migration status should be not set before testing")

			// finish migration with failure
			err := dbm.CreateMigrationDBAndSetStatus(migrationBlockNum)
			assert.NoError(t, err)
			endCheck := dbm.FinishStateMigration(false) // migration fail
			select {
			case <-endCheck: // wait for removing DB
			case <-time.After(1 * time.Second):
				t.Log("Take too long for a DB to be removed")
				t.FailNow()
			}

			// check if in migration state
			assert.False(t, dbm.InMigration())

			// check if state DB Path is set to old DB in MiscDB
			statDBPathKey := append(databaseDirPrefix, common.Int64ToByteBigEndian(uint64(StateTrieDB))...)
			fetchedStateDBPath, err := dbm.getDatabase(MiscDB).Get(statDBPathKey)
			assert.NoError(t, err)
			dirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 1, len(dirNames)) // check if DB is removed
			assert.Equal(t, initialDirNames[0], string(fetchedStateDBPath), "old DB should remain")

			// check if migration DB Path is not set in MiscDB
			migrationDBPathKey := append(databaseDirPrefix, common.Int64ToByteBigEndian(uint64(StateTrieMigrationDB))...)
			fetchedMigrationPath, err := dbm.getDatabase(MiscDB).Get(migrationDBPathKey)
			assert.NoError(t, err)
			assert.Equal(t, "", string(fetchedMigrationPath))

			// check if block number is not set in MiscDB
			fetchedBlockNum, err := dbm.getDatabase(MiscDB).Get(migrationStatusKey)
			assert.NoError(t, err)
			assert.Equal(t, common.Int64ToByteBigEndian(0), fetchedBlockNum)
		}

		// check status in miscDB on successful state migration
		{
			// check if not in migration state
			assert.False(t, dbManagers[i].InMigration(), "migration status should be not set before testing")

			// finish migration successfully
			err := dbm.CreateMigrationDBAndSetStatus(migrationBlockNum2)
			assert.NoError(t, err)
			endCheck := dbm.FinishStateMigration(true) // migration succeed
			select {
			case <-endCheck: // wait for removing DB
			case <-time.After(1 * time.Second):
				t.Log("Take too long for a DB to be removed")
				t.FailNow()
			}

			// check if in migration state
			assert.False(t, dbm.InMigration())

			// check if state DB Path is set to new DB in MiscDB
			statDBPathKey := append(databaseDirPrefix, common.Int64ToByteBigEndian(uint64(StateTrieDB))...)
			fetchedStateDBPath, err := dbm.getDatabase(MiscDB).Get(statDBPathKey)
			assert.NoError(t, err)
			dirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 1, len(dirNames))                                                         // check if DB is removed
			expectedStateDBPath := "statetrie_migrated_" + strconv.FormatUint(migrationBlockNum2, 10) // new DB format
			assert.Equal(t, expectedStateDBPath, string(fetchedStateDBPath), "new DB should remain")

			// check if migration DB Path is not set in MiscDB
			migrationDBPathKey := append(databaseDirPrefix, common.Int64ToByteBigEndian(uint64(StateTrieMigrationDB))...)
			fetchedMigrationPath, err := dbm.getDatabase(MiscDB).Get(migrationDBPathKey)
			assert.NoError(t, err)
			assert.Equal(t, "", string(fetchedMigrationPath))

			// check if block number is not set in MiscDB
			fetchedBlockNum, err := dbm.getDatabase(MiscDB).Get(migrationStatusKey)
			assert.NoError(t, err)
			assert.Equal(t, common.Int64ToByteBigEndian(0), fetchedBlockNum)
		}
	}
}

// While state trie migration, directory should be created with expected name
func TestDBManager_StateMigrationDBPath(t *testing.T) {
	for i, dbm := range dbManagers {
		if dbm.IsSingle() || dbConfigs[i].DBType == MemoryDB {
			continue
		}

		// check directory creation on successful migration
		{
			migrationBlockNum := uint64(12345)
			NewMigrationPath := dbBaseDirs[StateTrieMigrationDB] + "_" + strconv.FormatUint(migrationBlockNum, 10)

			// check if there is only one state trie db before migration
			initialDirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 1, len(initialDirNames), "migration status should be not set before testing")

			// check if new db is created
			err := dbm.CreateMigrationDBAndSetStatus(migrationBlockNum)
			assert.NoError(t, err)
			dirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 2, len(dirNames))
			assert.True(t, dirNames[0] == NewMigrationPath || dirNames[1] == NewMigrationPath)

			// check if old db is deleted on migration success
			endCheck := dbm.FinishStateMigration(true) // migration succeed
			select {
			case <-endCheck: // wait for removing DB
			case <-time.After(1 * time.Second):
				t.Log("Take too long for a DB to be removed")
				t.FailNow()
			}

			newDirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 1, len(newDirNames)) // check if DB is removed
			assert.Equal(t, NewMigrationPath, newDirNames[0], "new DB should remain")
		}

		// check directory creation on failed migration
		{
			migrationBlockNum := uint64(54321)
			NewMigrationPath := dbBaseDirs[StateTrieMigrationDB] + "_" + strconv.FormatUint(migrationBlockNum, 10)

			// check if there is only one state trie db before migration
			initialDirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 1, len(initialDirNames), "migration status should be not set before testing")

			// check if new db is created
			err := dbm.CreateMigrationDBAndSetStatus(migrationBlockNum)
			assert.NoError(t, err)
			dirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, "statetrie")
			assert.Equal(t, 2, len(dirNames))

			assert.True(t, dirNames[0] == NewMigrationPath || dirNames[1] == NewMigrationPath)

			// check if new db is deleted on migration fail
			endCheck := dbm.FinishStateMigration(false) // migration fail
			select {
			case <-endCheck: // wait for removing DB
			case <-time.After(1 * time.Second):
				t.Log("Take too long for a DB to be removed")
				t.FailNow()
			}

			newDirNames := getFilesInDir(t, dbm.GetDBConfig().Dir, dbm.getDBDir(StateTrieDB))
			assert.Equal(t, 1, len(newDirNames)) // check if DB is removed
			assert.Equal(t, initialDirNames[0], newDirNames[0], "old DB should remain")
		}
	}
}

func TestDBManager_WriteGovernanceIdx(t *testing.T) {
	testIdxes := []uint64{100, 200, 300}

	for _, dbm := range dbManagers {
		// normal case
		{
			// write test indexes
			for _, idx := range testIdxes {
				assert.Nil(t, dbm.WriteGovernanceIdx(idx))
			}

			// get the stored indexes
			encodedIdxes, err := dbm.GetMiscDB().Get(governanceHistoryKey)
			assert.Nil(t, err)

			// read and check the indexes from the database
			actualIdxes := make([]uint64, 0)
			assert.Nil(t, json.Unmarshal(encodedIdxes, &actualIdxes))
			assert.Equal(t, testIdxes, actualIdxes)
		}

		// unexpected case: try to write a governance index not in ascending order
		{
			assert.NotNil(t, dbm.WriteGovernanceIdx(testIdxes[0]))
		}

		// remove test data from database
		_ = dbm.GetMiscDB().Delete(governanceHistoryKey)
	}
}

func TestDBManager_ReadRecentGovernanceIdx(t *testing.T) {
	testIdxes := []uint64{100, 200, 300}

	for _, dbm := range dbManagers {
		// check empty
		idxes, err := dbm.ReadRecentGovernanceIdx(0)
		assert.Nil(t, idxes)
		assert.NotNil(t, err)

		// normal case
		{
			// write indexes on the database
			data, err := json.Marshal(testIdxes)
			assert.Nil(t, err)
			assert.Nil(t, dbm.GetMiscDB().Put(governanceHistoryKey, data))

			// read and check the indexes from the database
			idxes, err = dbm.ReadRecentGovernanceIdx(0)
			assert.Equal(t, testIdxes, idxes)
			assert.Nil(t, err)
		}

		// unexpected case: the governance indexes in the database is not in ascending order
		{
			invalidTestIdxes := append(testIdxes, testIdxes[0])
			expectedIdxes := append([]uint64{testIdxes[0]}, testIdxes...)

			// write invalid indexes on the database
			data, err := json.Marshal(invalidTestIdxes)
			assert.Nil(t, err)
			assert.Nil(t, dbm.GetMiscDB().Put(governanceHistoryKey, data))

			// read and check the indexes from the database
			idxes, err = dbm.ReadRecentGovernanceIdx(0)
			assert.Nil(t, err)
			assert.Equal(t, expectedIdxes, idxes)
		}

		// remove test data from database
		_ = dbm.GetMiscDB().Delete(governanceHistoryKey)
	}
}

func genReceipt(gasUsed int) *types.Receipt {
	log := &types.Log{Topics: []common.Hash{}, Data: []uint8{}, BlockNumber: uint64(gasUsed)}
	log.Topics = append(log.Topics, common.HexToHash(strconv.Itoa(gasUsed)))
	return &types.Receipt{
		TxHash:  common.HexToHash(strconv.Itoa(gasUsed)),
		GasUsed: uint64(gasUsed),
		Status:  types.ReceiptStatusSuccessful,
		Logs:    []*types.Log{log},
	}
}

func genTransaction(val uint64) (*types.Transaction, error) {
	return types.SignTx(
		types.NewTransaction(0, addr,
			big.NewInt(int64(val)), 0, big.NewInt(int64(val)), nil), signer, key)
}

// getFilesInDir returns all file names containing the substring in the directory
func getFilesInDir(t *testing.T, dirPath string, substr string) []string {
	files, err := ioutil.ReadDir(dirPath)
	assert.NoError(t, err)

	var dirNames []string
	for _, f := range files {
		if strings.Contains(f.Name(), substr) {
			dirNames = append(dirNames, f.Name())
		}
	}

	return dirNames
}

func TestDBManager_WriteAndReadAccountSnapshot(t *testing.T) {
	var (
		hash     common.Hash
		expected []byte
		actual   []byte
	)

	for _, dbm := range dbManagers {
		// read unknown key
		hash, _ = genRandomData()
		actual = dbm.ReadAccountSnapshot(hash)
		assert.Nil(t, actual)

		// write and read with empty hash
		_, expected = genRandomData()
		dbm.WriteAccountSnapshot(common.Hash{}, expected)
		actual = dbm.ReadAccountSnapshot(common.Hash{})
		assert.Equal(t, expected, actual)

		// write and read with empty data
		hash, _ = genRandomData()
		dbm.WriteAccountSnapshot(hash, []byte{})
		actual = dbm.ReadAccountSnapshot(hash)
		assert.Equal(t, []byte{}, actual)

		// write and read with random hash and data
		hash, expected = genRandomData()
		dbm.WriteAccountSnapshot(hash, expected)
		actual = dbm.ReadAccountSnapshot(hash)
		assert.Equal(t, expected, actual)
	}
}

func TestDBManager_DeleteAccountSnapshot(t *testing.T) {
	var (
		hash     common.Hash
		expected []byte
		actual   []byte
	)

	for _, dbm := range dbManagers {
		// delete unknown key
		hash, _ = genRandomData()
		dbm.DeleteAccountSnapshot(hash)
		actual = dbm.ReadAccountSnapshot(hash)
		assert.Nil(t, actual)

		// delete empty hash
		_, expected = genRandomData()
		dbm.WriteAccountSnapshot(common.Hash{}, expected)
		dbm.DeleteAccountSnapshot(common.Hash{})
		actual = dbm.ReadAccountSnapshot(hash)
		assert.Nil(t, actual)

		// write and read with empty data
		hash, _ = genRandomData()
		dbm.WriteAccountSnapshot(hash, []byte{})
		dbm.DeleteAccountSnapshot(hash)
		actual = dbm.ReadAccountSnapshot(hash)
		assert.Nil(t, actual)

		// write and read with random hash and data
		hash, expected = genRandomData()
		dbm.WriteAccountSnapshot(hash, expected)
		dbm.DeleteAccountSnapshot(hash)
		actual = dbm.ReadAccountSnapshot(hash)
		assert.Nil(t, actual)
	}
}

func genRandomData() (common.Hash, []byte) {
	rb := common.MakeRandomBytes(common.HashLength)
	hash := common.BytesToHash(rb)
	data := common.MakeRandomBytes(100)
	return hash, data
}
