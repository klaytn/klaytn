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

package tests

import (
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

type testEntry struct {
	name     string
	dbConfig *database.DBConfig
}

var testEntries = []testEntry{
	{"BadgerDB-Single", &database.DBConfig{DBType: database.BadgerDB, SingleDB: true}},
	{"BadgerDB", &database.DBConfig{DBType: database.BadgerDB, SingleDB: false, NumStateTrieShards: 4}},
	{"MemoryDB-Single", &database.DBConfig{DBType: database.MemoryDB, SingleDB: true}},
	{"MemoryDB", &database.DBConfig{DBType: database.MemoryDB, SingleDB: false, NumStateTrieShards: 4}},
	{"LevelDB-Single", &database.DBConfig{DBType: database.LevelDB, SingleDB: true, LevelDBCacheSize: 128, OpenFilesLimit: 32}},
	{"LevelDB", &database.DBConfig{DBType: database.LevelDB, SingleDB: false, LevelDBCacheSize: 128, OpenFilesLimit: 32, NumStateTrieShards: 4}},
}

// TestDBManager_WriteAndRead_Functional checks basic functionality of database.DBManager interface
// with underlying Database with combination of various configurations. Test cases usually include
// 1) read before write, 2) read after write, 3) read after overwrite and 4) read after delete.
// Sometimes 3) and 4) are omitted if such operation is not possible due to some reasons.
func TestDBManager_WriteAndRead_Functional(t *testing.T) {
	for _, entry := range testEntries {
		tempDir, err := ioutil.TempDir("", "klaytn-db-manager-test")
		if err != nil {
			t.Fatalf("cannot create temporary directory: %v", err)
		}
		entry.dbConfig.Dir = tempDir
		dbManager := database.NewDBManager(entry.dbConfig)

		t.Run(entry.name, func(t *testing.T) {
			testWriteAndReadCanonicalHash(t, dbManager)
			testWriteAndReadHeaderHash(t, dbManager)
			testWriteAndReadBlockHash(t, dbManager)
			testWriteAndReadFastBlockHash(t, dbManager)
			testWriteAndReadFastTrieProgress(t, dbManager)
			testWriteAndReadHeader(t, dbManager)
			testWriteAndReadBody(t, dbManager)
			testWriteAndReadTd(t, dbManager)
			testWriteAndReadReceipts(t, dbManager)
			testWriteAndReadBlock(t, dbManager)
			// TODO-Klaytn-Storage To implement this test case, error shouldn't be returned.
			// testWriteAndReadIstanbulSnapshot(t, dbManager)
		})

		dbManager.Close()
		os.RemoveAll(tempDir)
	}
}

func testWriteAndReadCanonicalHash(t *testing.T, dbManager database.DBManager) {
	hash1 := common.HexToHash("111")
	hash2 := common.HexToHash("222")

	// 1. Before write, empty hash should be returned.
	assert.Equal(t, common.Hash{}, dbManager.ReadCanonicalHash(111))

	// 2. After write, written hash should be returned.
	dbManager.WriteCanonicalHash(hash1, 111)
	assert.Equal(t, hash1, dbManager.ReadCanonicalHash(111))

	// 3. After overwrite, overwritten hash should be returned.
	dbManager.WriteCanonicalHash(hash2, 111)
	assert.Equal(t, hash2, dbManager.ReadCanonicalHash(111))

	// 4. After delete, empty hash should be returned.
	dbManager.DeleteCanonicalHash(111)
	assert.Equal(t, common.Hash{}, dbManager.ReadCanonicalHash(111))
}

func testWriteAndReadHeaderHash(t *testing.T, dbManager database.DBManager) {
	hash1 := common.HexToHash("111")
	hash2 := common.HexToHash("222")

	// 1. Before write, empty hash should be returned.
	assert.Equal(t, common.Hash{}, dbManager.ReadHeadHeaderHash())

	// 2. After write, written hash should be returned.
	dbManager.WriteHeadHeaderHash(hash1)
	assert.Equal(t, hash1, dbManager.ReadHeadHeaderHash())

	// 3. After overwrite, overwritten hash should be returned.
	dbManager.WriteHeadHeaderHash(hash2)
	assert.Equal(t, hash2, dbManager.ReadHeadHeaderHash())

	// 4. Delete operation is not supported.
}

func testWriteAndReadBlockHash(t *testing.T, dbManager database.DBManager) {
	hash1 := common.HexToHash("111")
	hash2 := common.HexToHash("222")

	// 1. Before write, empty hash should be returned.
	assert.Equal(t, common.Hash{}, dbManager.ReadHeadBlockHash())

	// 2. After write, written hash should be returned.
	dbManager.WriteHeadBlockHash(hash1)
	assert.Equal(t, hash1, dbManager.ReadHeadBlockHash())

	// 3. After overwrite, overwritten hash should be returned.
	dbManager.WriteHeadBlockHash(hash2)
	assert.Equal(t, hash2, dbManager.ReadHeadBlockHash())

	// 4. Delete operation is not supported.
}

func testWriteAndReadFastBlockHash(t *testing.T, dbManager database.DBManager) {
	hash1 := common.HexToHash("111")
	hash2 := common.HexToHash("222")

	// 1. Before write, empty hash should be returned.
	assert.Equal(t, common.Hash{}, dbManager.ReadHeadFastBlockHash())

	// 2. After write, written hash should be returned.
	dbManager.WriteHeadFastBlockHash(hash1)
	assert.Equal(t, hash1, dbManager.ReadHeadFastBlockHash())

	// 3. After overwrite, overwritten hash should be returned.
	dbManager.WriteHeadFastBlockHash(hash2)
	assert.Equal(t, hash2, dbManager.ReadHeadFastBlockHash())

	// 4. Delete operation is not supported.
}

func testWriteAndReadFastTrieProgress(t *testing.T, dbManager database.DBManager) {
	// 1. Before write, empty hash should be returned.
	assert.Equal(t, uint64(0), dbManager.ReadFastTrieProgress())

	// 2. After write, written hash should be returned.
	dbManager.WriteFastTrieProgress(111)
	assert.Equal(t, uint64(111), dbManager.ReadFastTrieProgress())

	// 3. After overwrite, overwritten hash should be returned.
	dbManager.WriteFastTrieProgress(222)
	assert.Equal(t, uint64(222), dbManager.ReadFastTrieProgress())

	// 4. Delete operation is not supported.
}

func testWriteAndReadHeader(t *testing.T, dbManager database.DBManager) {
	hash, blockNum, header := generateHeaderWithBlockNum(111)

	// 1. Before write, nil header should be returned.
	assert.Equal(t, (*types.Header)(nil), dbManager.ReadHeader(hash, blockNum))

	// 2. After write, written header should be returned.
	dbManager.WriteHeader(header)
	assert.Equal(t, header, dbManager.ReadHeader(hash, blockNum))

	// 3. Overwriting header with identical hash and blockNumber is not possible.

	// 4. After delete, deleted header should not be returned.
	dbManager.DeleteHeader(hash, blockNum)
	assert.Equal(t, (*types.Header)(nil), dbManager.ReadHeader(hash, blockNum))
}

func testWriteAndReadBody(t *testing.T, dbManager database.DBManager) {
	body := &types.Body{}

	tx := generateTx(t)
	txs := types.Transactions{tx}
	body.Transactions = txs

	hash, blockNum, _ := generateHeaderWithBlockNum(111)

	// 1. Before write, nil should be returned.
	assert.Equal(t, (*types.Body)(nil), dbManager.ReadBody(hash, blockNum))

	// 2. After write, written body should be returned.
	dbManager.WriteBody(hash, blockNum, body)
	assert.Equal(t, body.Transactions[0].Hash(), dbManager.ReadBody(hash, blockNum).Transactions[0].Hash())

	// 3. Overwriting body with identical hash and blockNumber is not possible.

	// 4. After delete, nil should be returned.
	dbManager.DeleteBody(hash, blockNum)
	assert.Equal(t, (*types.Body)(nil), dbManager.ReadBody(hash, blockNum))
}

func testWriteAndReadTd(t *testing.T, dbManager database.DBManager) {
	hash1 := common.HexToHash("111")
	blockNumber1 := uint64(111)

	td1 := new(big.Int).SetUint64(111)
	td2 := new(big.Int).SetUint64(222)

	// 1. Before write, empty td should be returned.
	assert.Equal(t, (*big.Int)(nil), dbManager.ReadTd(hash1, blockNumber1))

	// 2. After write, written td should be returned.
	dbManager.WriteTd(hash1, blockNumber1, td1)
	assert.Equal(t, td1, dbManager.ReadTd(hash1, blockNumber1))

	// 3. After overwrite, overwritten td should be returned.
	dbManager.WriteTd(hash1, blockNumber1, td2)
	assert.Equal(t, td2, dbManager.ReadTd(hash1, blockNumber1))

	// 4. After delete, deleted td should not be returned.
	dbManager.DeleteTd(hash1, blockNumber1)
	assert.Equal(t, (*big.Int)(nil), dbManager.ReadTd(hash1, blockNumber1))
}

func testWriteAndReadReceipts(t *testing.T, dbManager database.DBManager) {
	receipts := types.Receipts{
		generateReceipt(111),
		generateReceipt(222),
		generateReceipt(333),
	}

	hash := common.HexToHash("111")
	blockNumber := uint64(111)

	// 1. Before write, nil should be returned.
	assert.Equal(t, (types.Receipts)(nil), dbManager.ReadReceipts(hash, blockNumber))

	// 2. After write, written receipts should be returned.
	dbManager.WriteReceipts(hash, blockNumber, receipts)
	receiptsFromDB := dbManager.ReadReceipts(hash, blockNumber)
	assert.Equal(t, len(receipts), len(receiptsFromDB))
	assert.Equal(t, receipts, receiptsFromDB)

	// 3. After overwrite, overwritten receipts should be returned.
	receipts2 := types.Receipts{
		generateReceipt(444),
		generateReceipt(555),
		generateReceipt(666),
	}
	dbManager.WriteReceipts(hash, blockNumber, receipts2)
	receiptsFromDB = dbManager.ReadReceipts(hash, blockNumber)
	assert.Equal(t, receipts2, receiptsFromDB)
	assert.NotEqual(t, receipts, receiptsFromDB)

	// 4. After delete, nil should be returned.
	dbManager.DeleteReceipts(hash, blockNumber)
	assert.Equal(t, (types.Receipts)(nil), dbManager.ReadReceipts(hash, blockNumber))
}

func testWriteAndReadBlock(t *testing.T, dbManager database.DBManager) {
	hash, blockNum, header := generateHeaderWithBlockNum(111)
	body := generateBody(t)
	receipts := types.Receipts{generateReceipt(111)}

	blockchain.InitDeriveSha(types.ImplDeriveShaOriginal)
	block := types.NewBlock(header, body.Transactions, receipts)

	// 1. Before write, nil should be returned.
	assert.Equal(t, (*types.Block)(nil), dbManager.ReadBlock(hash, blockNum))

	// 2. After write, written block should be returned.
	dbManager.WriteBlock(block)
	blockFromDB := dbManager.ReadBlock(block.Hash(), block.NumberU64())
	assert.Equal(t, block.Header(), blockFromDB.Header())
	assert.Equal(t, block.Transactions()[0].Hash(), blockFromDB.Transactions()[0].Hash())

	// 3. Overwriting block with exact hash and blockNumber is not possible.

	// 4. After delete, nil should be returned.
	dbManager.DeleteBlock(block.Hash(), block.NumberU64())
	assert.Equal(t, (*types.Block)(nil), dbManager.ReadBlock(hash, blockNum))
}

func testWriteAndReadIstanbulSnapshot(t *testing.T, dbManager database.DBManager) {
	// TODO-Klaytn-Storage To implement this test case, error shouldn't be returned.
}

func generateHeaderWithBlockNum(blockNum int) (common.Hash, uint64, *types.Header) {
	parentHash := common.HexToHash(strconv.Itoa(blockNum))
	blockNumber := new(big.Int).SetUint64(uint64(blockNum))

	header := &types.Header{
		ParentHash: parentHash,
		Number:     blockNumber,
		BlockScore: blockNumber,
		Time:       blockNumber,
		Extra:      []byte{'a', 'b', 'c'},
		Governance: common.Hex2Bytes("b8dc7b22676f7665726e696e676e6f6465223a22307865373333636234643237396461363936663330643437306638633034646563623534666362306432222c22676f7665726e616e63656d6f6465223a2273696e676c65222c22726577617264223a7b226d696e74696e67616d6f756e74223a393630303030303030303030303030303030302c22726174696f223a2233342f33332f3333227d2c22626674223a7b2265706f6368223a33303030302c22706f6c696379223a302c22737562223a32317d2c22756e69745072696365223a32353030303030303030307d"),
		Vote:       common.Hex2Bytes("e194e733cb4d279da696f30d470f8c04decb54fcb0d28565706f6368853330303030"),
	}
	return header.Hash(), uint64(blockNum), header
}

func generateBody(t *testing.T) *types.Body {
	body := &types.Body{}

	tx := generateTx(t)
	txs := types.Transactions{tx}
	body.Transactions = txs

	return body
}

func generateTx(t *testing.T) *types.Transaction {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	signer := types.LatestSignerForChainID(big.NewInt(18))
	tx1, err := types.SignTx(types.NewTransaction(0, addr, new(big.Int), 0, new(big.Int), nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}
	return tx1
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
