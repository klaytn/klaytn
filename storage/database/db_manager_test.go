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
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"testing"
)

var dbManagers []DBManager
var dbConfigs = make([]*DBConfig, 0, len(baseConfigs)*3)
var baseConfigs = []*DBConfig{
	{DBType: LevelDB, Partitioned: false, NumStateTriePartitions: 1, ParallelDBWrite: false},
	{DBType: LevelDB, Partitioned: false, NumStateTriePartitions: 1, ParallelDBWrite: true},
	{DBType: LevelDB, Partitioned: false, NumStateTriePartitions: 4, ParallelDBWrite: false},
	{DBType: LevelDB, Partitioned: false, NumStateTriePartitions: 4, ParallelDBWrite: true},

	{DBType: LevelDB, Partitioned: true, NumStateTriePartitions: 1, ParallelDBWrite: false},
	{DBType: LevelDB, Partitioned: true, NumStateTriePartitions: 1, ParallelDBWrite: true},
	{DBType: LevelDB, Partitioned: true, NumStateTriePartitions: 4, ParallelDBWrite: false},
	{DBType: LevelDB, Partitioned: true, NumStateTriePartitions: 4, ParallelDBWrite: true},
}

var num1 = uint64(20190815)
var num2 = uint64(20199999)
var hash1 = common.HexToHash("1341655") // 20190805 in hexadecimal
var hash2 = common.HexToHash("1343A3F") // 20199999 in hexadecimal

func init() {
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
		assert.Equal(t, num1, *dbm.ReadHeaderNumber(headerHash)) // headerNumber still exists in the cache
	}
}

// TestDBManager_Header tests read, write and delete operations of blockchain bodies.
func TestDBManager_Body(t *testing.T) {
	body := &types.Body{Transactions:types.Transactions{}}
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
	receipts := types.Receipts{generateReceipt(111)}

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
		snapshot, _ := dbm.ReadIstanbulSnapshot(hash1)
		assert.Nil(t, snapshot)

		dbm.WriteIstanbulSnapshot(hash1, hash2[:])
		snapshot, _ = dbm.ReadIstanbulSnapshot(hash1)
		assert.Equal(t, hash2[:], snapshot)

		dbm.WriteIstanbulSnapshot(hash1, hash1[:])
		snapshot, _ = dbm.ReadIstanbulSnapshot(hash1)
		assert.Equal(t, hash1[:], snapshot)
	}
}

// TestDBManager_DatabaseVersion tests read/write operations of database version.
func TestDBManager_DatabaseVersion(t *testing.T) {
	// TODO-Klaytn-Database DatabaseVersion should be handled carefully.
	//for i, dbm := range dbManagers {
	//	c := dbConfigs[i]
	//
	//	assert.Equal(t, uint64(0), dbm.ReadDatabaseVersion())
	//
	//	dbm.WriteDatabaseVersion(uint64(1))
	//	assert.Equal(t, uint64(1), dbm.ReadDatabaseVersion())
	//
	//	dbm.WriteDatabaseVersion(uint64(2))
	//	assert.Equal(t, uint64(2), dbm.ReadDatabaseVersion())
	//}
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
		assert.Equal(t, common.Hash{}, dbm.ConvertServiceChainBlockHashToMainChainTxHash(hash1))

		dbm.WriteChildChainTxHash(hash1, hash1)
		assert.Equal(t, hash1, dbm.ConvertServiceChainBlockHashToMainChainTxHash(hash1))

		dbm.WriteChildChainTxHash(hash1, hash2)
		assert.Equal(t, hash2, dbm.ConvertServiceChainBlockHashToMainChainTxHash(hash1))

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

// TestDBManager_FastTrieProgress tests read and write operations of clique snapshots.
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

func generateReceipt(gasUsed int) *types.Receipt {
	log := &types.Log{Topics: []common.Hash{}, Data: []uint8{}, BlockNumber: uint64(gasUsed)}
	log.Topics = append(log.Topics, common.HexToHash(strconv.Itoa(gasUsed)))
	return &types.Receipt{
		TxHash:  common.HexToHash(strconv.Itoa(gasUsed)),
		GasUsed: uint64(gasUsed),
		Status:  types.ReceiptStatusSuccessful,
		Logs:    []*types.Log{log},
	}
}
