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
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

func TestChildChainData_ReadAndWrite_ChildChainTxHash(t *testing.T) {
	dir, err := ioutil.TempDir("", "klaytn-test-child-chain-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	dbc := &DBConfig{Dir: dir, DBType: LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	dbm := NewDBManager(dbc)
	defer dbm.Close()

	ccBlockHash := common.HexToHash("0x0e0e0e0e0e0e0e0e0e0e0e0e0e0e0e0e")
	ccTxHash := common.HexToHash("0x0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f")

	// Before writing the data into DB, nil should be returned.
	ccTxHashFromDB := dbm.ConvertChildChainBlockHashToParentChainTxHash(ccBlockHash)
	assert.Equal(t, common.Hash{}, ccTxHashFromDB)

	// After writing the data into DB, data should be returned.
	dbm.WriteChildChainTxHash(ccBlockHash, ccTxHash)
	ccTxHashFromDB = dbm.ConvertChildChainBlockHashToParentChainTxHash(ccBlockHash)
	assert.NotNil(t, ccTxHashFromDB)
	assert.Equal(t, ccTxHash, ccTxHashFromDB)

	ccBlockHashFake := common.HexToHash("0x0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a")
	// Invalid information should not return the data.
	ccTxHashFromDB = dbm.ConvertChildChainBlockHashToParentChainTxHash(ccBlockHashFake)
	assert.Equal(t, common.Hash{}, ccTxHashFromDB)
}

func TestLastIndexedBlockData_ReadAndWrite_AnchoredBlockNumber(t *testing.T) {
	dir, err := ioutil.TempDir("", "klaytn-test-child-chain-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	dbc := &DBConfig{Dir: dir, DBType: LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	dbm := NewDBManager(dbc)
	defer dbm.Close()

	blockNum := uint64(123)

	blockNumFromDB := dbm.GetLastIndexedBlockNumber()
	assert.Equal(t, uint64(0), blockNumFromDB)

	dbm.WriteLastIndexedBlockNumber(blockNum)
	blockNumFromDB = dbm.GetLastIndexedBlockNumber()
	assert.Equal(t, blockNum, blockNumFromDB)

	newBlockNum := uint64(321)
	dbm.WriteLastIndexedBlockNumber(newBlockNum)
	blockNumFromDB = dbm.GetLastIndexedBlockNumber()
	assert.Equal(t, newBlockNum, blockNumFromDB)
}

func TestChildChainData_ReadAndWrite_AnchoredBlockNumber(t *testing.T) {
	dir, err := ioutil.TempDir("", "klaytn-test-child-chain-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	dbc := &DBConfig{Dir: dir, DBType: LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	dbm := NewDBManager(dbc)
	defer dbm.Close()

	blockNum := uint64(123)

	blockNumFromDB := dbm.ReadAnchoredBlockNumber()
	assert.Equal(t, uint64(0), blockNumFromDB)

	dbm.WriteAnchoredBlockNumber(blockNum)
	blockNumFromDB = dbm.ReadAnchoredBlockNumber()
	assert.Equal(t, blockNum, blockNumFromDB)

	newBlockNum := uint64(321)
	dbm.WriteAnchoredBlockNumber(newBlockNum)
	blockNumFromDB = dbm.ReadAnchoredBlockNumber()
	assert.Equal(t, newBlockNum, blockNumFromDB)
}

func TestChildChainData_ReadAndWrite_ReceiptFromParentChain(t *testing.T) {
	dir, err := ioutil.TempDir("", "klaytn-test-child-chain-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	dbc := &DBConfig{Dir: dir, DBType: LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	dbm := NewDBManager(dbc)
	defer dbm.Close()

	blockHash := common.HexToHash("0x0e0e0e0e0e0e0e0e0e0e0e0e0e0e0e0e")
	rct := &types.Receipt{}
	rct.TxHash = common.BigToHash(big.NewInt(12345))
	rct.GasUsed = uint64(12345)
	rct.Status = types.ReceiptStatusSuccessful

	rctFromDB := dbm.ReadReceiptFromParentChain(blockHash)
	assert.Nil(t, rctFromDB)

	dbm.WriteReceiptFromParentChain(blockHash, rct)
	rctFromDB = dbm.ReadReceiptFromParentChain(blockHash)

	assert.Equal(t, rct.Status, rctFromDB.Status)
	assert.Equal(t, rct.GasUsed, rctFromDB.GasUsed)
	assert.Equal(t, rct.TxHash, rctFromDB.TxHash)

	newBlockHash := common.HexToHash("0x0f0f0e0e0e0e0e0e0e0e0e0e0e0e0f0f")
	rctFromDB = dbm.ReadReceiptFromParentChain(newBlockHash)
	assert.Nil(t, rctFromDB)
}

func TestChildChainData_ReadAndWrite_ValueTransferTxHash(t *testing.T) {
	dir, err := ioutil.TempDir("", "klaytn-test-child-chain-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	dbc := &DBConfig{Dir: dir, DBType: LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	dbm := NewDBManager(dbc)
	defer dbm.Close()

	rTxHash := common.HexToHash("0x0e0e0e0e0e0e0e0e0e0e0e0e0e0e0e0e")
	hTxHash := common.HexToHash("0x0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f")

	// Before writing the data into DB, nil should be returned.
	hTxHashFromDB := dbm.ReadHandleTxHashFromRequestTxHash(rTxHash)
	assert.Equal(t, common.Hash{}, hTxHashFromDB)

	// After writing the data into DB, data should be returned.
	dbm.WriteHandleTxHashFromRequestTxHash(rTxHash, hTxHash)
	hTxHashFromDB = dbm.ReadHandleTxHashFromRequestTxHash(rTxHash)
	assert.NotNil(t, hTxHashFromDB)
	assert.Equal(t, hTxHash, hTxHashFromDB)

	ccBlockHashFake := common.HexToHash("0x0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a")
	// Invalid information should not return the data.
	hTxHashFromDB = dbm.ReadHandleTxHashFromRequestTxHash(ccBlockHashFake)
	assert.Equal(t, common.Hash{}, hTxHashFromDB)
}

func TestChildChainData_ReadAndWrite_OperatorFeePayer(t *testing.T) {
	dir, err := ioutil.TempDir("", "klaytn-test-child-chain-data")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	firstAddr := common.HexToAddress("0x1")
	secondAddr := common.HexToAddress("0x2")

	dbc := &DBConfig{Dir: dir, DBType: LevelDB, LevelDBCacheSize: 32, OpenFilesLimit: 32}
	dbm := NewDBManager(dbc)
	defer dbm.Close()

	// check initial value
	{
		feePayer := dbm.ReadParentOperatorFeePayer()
		assert.Equal(t, common.Address{}, feePayer)
	}
	{
		feePayer := dbm.ReadChildOperatorFeePayer()
		assert.Equal(t, common.Address{}, feePayer)
	}

	// check write/read
	{
		dbm.WriteParentOperatorFeePayer(firstAddr)
		feePayer := dbm.ReadParentOperatorFeePayer()
		assert.Equal(t, firstAddr, feePayer)
	}
	{
		dbm.WriteChildOperatorFeePayer(secondAddr)
		feePayer := dbm.ReadChildOperatorFeePayer()
		assert.Equal(t, secondAddr, feePayer)
	}

	// check write zero address and read
	{
		dbm.WriteParentOperatorFeePayer(common.Address{})
		feePayer := dbm.ReadParentOperatorFeePayer()
		assert.Equal(t, common.Address{}, feePayer)
	}
	{
		dbm.WriteChildOperatorFeePayer(common.Address{})
		feePayer := dbm.ReadChildOperatorFeePayer()
		assert.Equal(t, common.Address{}, feePayer)
	}
}
