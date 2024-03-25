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

package cn

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/consensus/istanbul/backend"
	mocks3 "github.com/klaytn/klaytn/event/mocks"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn/gasprice"
	mocks2 "github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/work/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func newCNAPIBackend(t *testing.T) (*gomock.Controller, *mocks.MockBlockChain, *mocks2.MockMiner, *CNAPIBackend) {
	mockCtrl := gomock.NewController(t)

	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockMiner := mocks2.NewMockMiner(mockCtrl)

	cn := &CN{blockchain: mockBlockChain, miner: mockMiner}

	return mockCtrl, mockBlockChain, mockMiner, &CNAPIBackend{cn: cn}
}

func TestCNAPIBackend_GetTxAndLookupInfoInCache(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	hash := hashes[0]

	expectedTx := tx1
	expectedBlockHash := hashes[1]
	expectedBlockNum := uint64(111)
	expectedIndex := uint64(999)

	mockBlockChain.EXPECT().GetTxAndLookupInfoInCache(hash).Times(1).Return(expectedTx, expectedBlockHash, expectedBlockNum, expectedIndex)
	tx, blockHash, blockNumber, index := api.GetTxAndLookupInfoInCache(hash)

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedBlockHash, blockHash)
	assert.Equal(t, expectedBlockNum, blockNumber)
	assert.Equal(t, expectedIndex, index)
}

func TestCNAPIBackend_GetTxLookupInfoAndReceipt(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	hash := hashes[0]

	expectedTx := tx1
	expectedBlockHash := hashes[1]
	expectedBlockNum := uint64(111)
	expectedIndex := uint64(999)
	expectedReceipt := newReceipt(123)

	mockBlockChain.EXPECT().GetTxLookupInfoAndReceipt(hash).Times(1).Return(expectedTx, expectedBlockHash, expectedBlockNum, expectedIndex, expectedReceipt)
	tx, blockHash, blockNumber, index, receipt := api.GetTxLookupInfoAndReceipt(context.Background(), hash)

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedBlockHash, blockHash)
	assert.Equal(t, expectedBlockNum, blockNumber)
	assert.Equal(t, expectedIndex, index)
	assert.Equal(t, expectedReceipt, receipt)
}

func TestCNAPIBackend_GetBlockReceiptsInCache(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	hash := hashes[0]
	expectedReceipts := types.Receipts{newReceipt(111), newReceipt(222)}

	mockBlockChain.EXPECT().GetBlockReceiptsInCache(hash).Return(expectedReceipts).Times(1)

	assert.Equal(t, expectedReceipts, api.GetBlockReceiptsInCache(hash))
}

func TestCNAPIBackend_GetTxLookupInfoAndReceiptInCache(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	hash := hashes[0]

	expectedTx := tx1
	expectedBlockHash := hashes[1]
	expectedBlockNum := uint64(111)
	expectedIndex := uint64(999)
	expectedReceipt := newReceipt(123)

	mockBlockChain.EXPECT().GetTxLookupInfoAndReceiptInCache(hash).Times(1).Return(expectedTx, expectedBlockHash, expectedBlockNum, expectedIndex, expectedReceipt)
	tx, blockHash, blockNumber, index, receipt := api.GetTxLookupInfoAndReceiptInCache(hash)

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedBlockHash, blockHash)
	assert.Equal(t, expectedBlockNum, blockNumber)
	assert.Equal(t, expectedIndex, index)
	assert.Equal(t, expectedReceipt, receipt)
}

func TestCNAPIBackend_ChainConfig(t *testing.T) {
	mockCtrl, _, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	assert.Nil(t, api.ChainConfig())

	emptyConfig := &params.ChainConfig{}
	api.cn.chainConfig = &*emptyConfig
	assert.Equal(t, emptyConfig, api.ChainConfig())

	nonEmptyConfig := &params.ChainConfig{ChainID: big.NewInt(123)}
	api.cn.chainConfig = &*nonEmptyConfig
	assert.Equal(t, nonEmptyConfig, api.ChainConfig())
}

func TestCNAPIBackend_CurrentBlock(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	block := newBlock(123)
	mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)

	assert.Equal(t, block, api.CurrentBlock())
}

func getTestConfig() *params.ChainConfig {
	config := params.TestChainConfig.Copy()
	config.Governance = params.GetDefaultGovernanceConfig()
	config.Istanbul = params.GetDefaultIstanbulConfig()
	return config
}

func testGov() *governance.MixedEngine {
	db := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	config := params.TestChainConfig.Copy()
	config.Governance = params.GetDefaultGovernanceConfig()
	config.Istanbul = params.GetDefaultIstanbulConfig()
	return governance.NewMixedEngine(config, db)
}

func TestCNAPIBackend_SetHead(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	mockDownloader := mocks2.NewMockProtocolManagerDownloader(mockCtrl)
	mockDownloader.EXPECT().Cancel().Times(1)
	pm := &ProtocolManager{downloader: mockDownloader}
	api.cn.protocolManager = pm
	api.cn.engine = gxhash.NewFullFaker()
	api.cn.governance = testGov()
	api.gpo = gasprice.NewOracle(api, gasprice.Config{}, nil)

	number := uint64(123)
	mockBlockChain.EXPECT().SetHead(number).Times(1)

	api.SetHead(number)
	block := newBlock(int(number))
	expectedHeader := block.Header()
	mockBlockChain.EXPECT().CurrentHeader().Return(expectedHeader).Times(1)
	assert.Equal(t, expectedHeader, mockBlockChain.CurrentHeader())
}

func TestCNAPIBackend_HeaderByNumber(t *testing.T) {
	blockNum := uint64(123)
	block := newBlock(int(blockNum))
	expectedHeader := block.Header()
	{
		mockCtrl, _, mockMiner, api := newCNAPIBackend(t)
		mockMiner.EXPECT().PendingBlock().Return(block).Times(1)

		header, err := api.HeaderByNumber(context.Background(), rpc.PendingBlockNumber)

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)

		header, err := api.HeaderByNumber(context.Background(), rpc.LatestBlockNumber)

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetHeaderByNumber(blockNum).Return(nil).Times(1)

		header, err := api.HeaderByNumber(context.Background(), rpc.BlockNumber(blockNum))

		assert.Nil(t, header)
		assert.Error(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetHeaderByNumber(blockNum).Return(expectedHeader).Times(1)

		header, err := api.HeaderByNumber(context.Background(), rpc.BlockNumber(blockNum))

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_HeaderByNumberOrHash(t *testing.T) {
	block := newBlock(123)
	expectedHeader := block.Header()
	{
		mockCtrl, _, mockMiner, api := newCNAPIBackend(t)
		mockMiner.EXPECT().PendingBlock().Return(block).Times(1)

		header, err := api.HeaderByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithNumber(rpc.PendingBlockNumber))

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)

		header, err := api.HeaderByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber))

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetHeaderByNumber(uint64(123)).Return(expectedHeader).Times(1)

		header, err := api.HeaderByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithNumber(rpc.BlockNumber(123)))

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetHeaderByHash(hash1).Return(expectedHeader).Times(1)

		header, err := api.HeaderByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithHash(hash1, false))

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_HeaderByHash(t *testing.T) {
	{
		blockNum := uint64(123)
		block := newBlock(int(blockNum))
		expectedHeader := block.Header()

		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetHeaderByHash(hash1).Return(expectedHeader).Times(1)

		header, err := api.HeaderByHash(context.Background(), hash1)

		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_BlockByNumber(t *testing.T) {
	blockNum := uint64(123)
	block := newBlock(int(blockNum))
	expectedBlock := block
	{
		mockCtrl, _, mockMiner, api := newCNAPIBackend(t)
		mockMiner.EXPECT().PendingBlock().Return(block).Times(1)

		block, err := api.BlockByNumber(context.Background(), rpc.PendingBlockNumber)

		assert.Equal(t, expectedBlock, block)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)

		block, err := api.BlockByNumber(context.Background(), rpc.LatestBlockNumber)

		assert.Equal(t, expectedBlock, block)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetBlockByNumber(blockNum).Return(nil).Times(1)

		block, err := api.BlockByNumber(context.Background(), rpc.BlockNumber(blockNum))

		assert.Nil(t, block)
		assert.Error(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetBlockByNumber(blockNum).Return(expectedBlock).Times(1)

		block, err := api.BlockByNumber(context.Background(), rpc.BlockNumber(blockNum))

		assert.Equal(t, expectedBlock, block)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_BlockByNumberOrHash(t *testing.T) {
	blockNum := uint64(123)
	block := newBlock(int(blockNum))
	expectedBlock := block
	{
		mockCtrl, _, mockMiner, api := newCNAPIBackend(t)
		mockMiner.EXPECT().PendingBlock().Return(block).Times(1)

		block, err := api.BlockByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithNumber(rpc.PendingBlockNumber))

		assert.Equal(t, expectedBlock, block)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().CurrentBlock().Return(expectedBlock).Times(1)

		block, err := api.BlockByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber))

		assert.Equal(t, expectedBlock, block)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetBlockByNumber(uint64(123)).Return(nil).Times(1)

		block, err := api.BlockByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithNumber(rpc.BlockNumber(123)))

		assert.Nil(t, block)
		assert.Error(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetBlockByHash(hash1).Return(expectedBlock).Times(1)

		block, err := api.BlockByNumberOrHash(context.Background(), rpc.NewBlockNumberOrHashWithHash(hash1, false))

		assert.Equal(t, expectedBlock, block)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_StateAndHeaderByNumber(t *testing.T) {
	blockNum := uint64(123)
	block := newBlock(int(blockNum))

	stateDB, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	stateDB.SetNonce(addrs[0], 123)
	stateDB.SetNonce(addrs[1], 321)

	expectedHeader := block.Header()
	{
		mockCtrl, _, mockMiner, api := newCNAPIBackend(t)
		mockMiner.EXPECT().Pending().Return(block, stateDB).Times(1)

		returnedStateDB, header, err := api.StateAndHeaderByNumber(context.Background(), rpc.PendingBlockNumber)

		assert.Equal(t, stateDB, returnedStateDB)
		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)

		mockBlockChain.EXPECT().GetHeaderByNumber(blockNum).Return(nil).Times(1)
		returnedStateDB, header, err := api.StateAndHeaderByNumber(context.Background(), rpc.BlockNumber(blockNum))

		assert.Nil(t, returnedStateDB)
		assert.Nil(t, header)
		assert.Error(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)

		mockBlockChain.EXPECT().GetHeaderByNumber(blockNum).Return(expectedHeader).Times(1)
		mockBlockChain.EXPECT().StateAt(expectedHeader.Root).Return(stateDB, nil).Times(1)
		returnedStateDB, header, err := api.StateAndHeaderByNumber(context.Background(), rpc.BlockNumber(blockNum))

		assert.Equal(t, stateDB, returnedStateDB)
		assert.Equal(t, expectedHeader, header)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_GetBlock(t *testing.T) {
	block := newBlock(123)
	hash := hashes[0]
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetBlockByHash(hash).Return(nil).Times(1)

		returnedBlock, err := api.BlockByHash(context.Background(), hash)
		assert.Nil(t, returnedBlock)
		assert.Error(t, err)

		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
		mockBlockChain.EXPECT().GetBlockByHash(hash).Return(block).Times(1)

		returnedBlock, err := api.BlockByHash(context.Background(), hash)
		assert.Equal(t, block, returnedBlock)
		assert.NoError(t, err)

		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_GetTxAndLookupInfo(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	hash := hashes[0]

	expectedTx := tx1
	expectedBlockHash := hashes[1]
	expectedBlockNum := uint64(111)
	expectedIndex := uint64(999)

	mockBlockChain.EXPECT().GetTxAndLookupInfo(hash).Times(1).Return(expectedTx, expectedBlockHash, expectedBlockNum, expectedIndex)
	tx, blockHash, blockNumber, index := api.GetTxAndLookupInfo(hash)

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedBlockHash, blockHash)
	assert.Equal(t, expectedBlockNum, blockNumber)
	assert.Equal(t, expectedIndex, index)
}

func TestCNAPIBackend_GetBlockReceipts(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	hash := hashes[0]
	expectedReceipts := types.Receipts{newReceipt(111), newReceipt(222)}

	mockBlockChain.EXPECT().GetReceiptsByBlockHash(hash).Return(expectedReceipts).Times(1)

	assert.Equal(t, expectedReceipts, api.GetBlockReceipts(context.Background(), hash))
}

func TestCNAPIBackend_GetLogs(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	hash := hashes[0]
	expectedLogs := [][]*types.Log{{{BlockNumber: 123}}, {{BlockNumber: 321}}}
	mockBlockChain.EXPECT().GetLogsByHash(hash).Return(expectedLogs).Times(1)

	logs, err := api.GetLogs(context.Background(), hash)
	assert.Equal(t, expectedLogs, logs)
	assert.NoError(t, err)
}

func TestCNAPIBackend_GetTd(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	td := big.NewInt(123)
	hash := hashes[0]
	mockBlockChain.EXPECT().GetTdByHash(hash).Return(td).Times(1)

	assert.Equal(t, td, api.GetTd(hash))
}

func TestCNAPIBackend_SubscribeEvents(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	api.cn.txPool = mockTxPool
	defer mockCtrl.Finish()

	rmCh := make(chan<- blockchain.RemovedLogsEvent)
	ceCh := make(chan<- blockchain.ChainEvent)
	chCh := make(chan<- blockchain.ChainHeadEvent)
	csCh := make(chan<- blockchain.ChainSideEvent)
	leCh := make(chan<- []*types.Log)
	txCh := make(chan<- blockchain.NewTxsEvent)

	sub := mocks3.NewMockSubscription(mockCtrl)

	mockBlockChain.EXPECT().SubscribeRemovedLogsEvent(rmCh).Return(sub).Times(1)
	mockBlockChain.EXPECT().SubscribeChainEvent(ceCh).Return(sub).Times(1)
	mockBlockChain.EXPECT().SubscribeChainHeadEvent(chCh).Return(sub).Times(1)
	mockBlockChain.EXPECT().SubscribeChainSideEvent(csCh).Return(sub).Times(1)
	mockBlockChain.EXPECT().SubscribeLogsEvent(leCh).Return(sub).Times(1)

	mockTxPool.EXPECT().SubscribeNewTxsEvent(txCh).Return(sub).Times(1)

	assert.Equal(t, sub, api.SubscribeRemovedLogsEvent(rmCh))
	assert.Equal(t, sub, api.SubscribeChainEvent(ceCh))
	assert.Equal(t, sub, api.SubscribeChainHeadEvent(chCh))
	assert.Equal(t, sub, api.SubscribeChainSideEvent(csCh))
	assert.Equal(t, sub, api.SubscribeLogsEvent(leCh))

	assert.Equal(t, sub, api.SubscribeNewTxsEvent(txCh))
}

func TestCNAPIBackend_SendTx(t *testing.T) {
	mockCtrl, _, _, api := newCNAPIBackend(t)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	mockTxPool.EXPECT().AddLocal(tx1).Return(expectedErr).Times(1)
	api.cn.txPool = mockTxPool

	defer mockCtrl.Finish()

	assert.Equal(t, expectedErr, api.SendTx(context.Background(), tx1))
}

func TestCNAPIBackend_GetPoolTransactions(t *testing.T) {
	{
		mockCtrl, _, _, api := newCNAPIBackend(t)
		mockTxPool := mocks.NewMockTxPool(mockCtrl)
		mockTxPool.EXPECT().Pending().Return(nil, expectedErr).Times(1)
		api.cn.txPool = mockTxPool

		txs, ReturnedErr := api.GetPoolTransactions()
		assert.Nil(t, txs)
		assert.Equal(t, expectedErr, ReturnedErr)
		mockCtrl.Finish()
	}
	{
		mockCtrl, _, _, api := newCNAPIBackend(t)
		mockTxPool := mocks.NewMockTxPool(mockCtrl)

		pendingTxs := map[common.Address]types.Transactions{addrs[0]: {tx1}}
		mockTxPool.EXPECT().Pending().Return(pendingTxs, nil).Times(1)
		api.cn.txPool = mockTxPool

		txs, ReturnedErr := api.GetPoolTransactions()
		assert.Equal(t, types.Transactions{tx1}, txs)
		assert.NoError(t, ReturnedErr)
		mockCtrl.Finish()
	}
}

func TestCNAPIBackend_GetPoolTransaction(t *testing.T) {
	hash := hashes[0]

	mockCtrl, _, _, api := newCNAPIBackend(t)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	mockTxPool.EXPECT().Get(hash).Return(tx1).Times(1)
	api.cn.txPool = mockTxPool

	defer mockCtrl.Finish()

	assert.Equal(t, tx1, api.GetPoolTransaction(hash))
}

func TestCNAPIBackend_GetPoolNonce(t *testing.T) {
	addr := addrs[0]
	nonce := uint64(123)

	mockCtrl, _, _, api := newCNAPIBackend(t)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	mockTxPool.EXPECT().GetPendingNonce(addr).Return(nonce).Times(1)
	api.cn.txPool = mockTxPool

	defer mockCtrl.Finish()

	assert.Equal(t, nonce, api.GetPoolNonce(context.Background(), addr))
}

func TestCNAPIBackend_Stats(t *testing.T) {
	pending := 123
	queued := 321

	mockCtrl, _, _, api := newCNAPIBackend(t)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	mockTxPool.EXPECT().Stats().Return(pending, queued).Times(1)
	api.cn.txPool = mockTxPool

	defer mockCtrl.Finish()

	p, q := api.Stats()
	assert.Equal(t, pending, p)
	assert.Equal(t, queued, q)
}

func TestCNAPIBackend_TxPoolContent(t *testing.T) {
	pending := map[common.Address]types.Transactions{addrs[0]: {tx1}}
	queued := map[common.Address]types.Transactions{addrs[1]: {tx1}}

	mockCtrl, _, _, api := newCNAPIBackend(t)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	mockTxPool.EXPECT().Content().Return(pending, queued).Times(1)
	api.cn.txPool = mockTxPool

	defer mockCtrl.Finish()

	p, q := api.TxPoolContent()
	assert.Equal(t, pending, p)
	assert.Equal(t, queued, q)
}

func TestCNAPIBackend_IsParallelDBWrite(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	mockBlockChain.EXPECT().IsParallelDBWrite().Return(true).Times(1)
	assert.True(t, api.IsParallelDBWrite())
}

func TestCNAPIBackend_IsSenderTxHashIndexingEnabled(t *testing.T) {
	mockCtrl, mockBlockChain, _, api := newCNAPIBackend(t)
	defer mockCtrl.Finish()

	mockBlockChain.EXPECT().IsSenderTxHashIndexingEnabled().Return(true).Times(1)
	assert.True(t, api.IsSenderTxHashIndexingEnabled())
}

type rewindTest struct {
	canonicalBlocks int // Number of blocks to generate for the canonical chain (heavier)

	setheadBlock       uint64 // Block number to set head back to
	expCanonicalBlocks int    // Number of canonical blocks expected to remain in the database (excl. genesis)
	expHeadHeader      uint64 // Block number of the expected head header
	expHeadFastBlock   uint64 // Block number of the expected head fast sync block
	expHeadBlock       uint64 // Block number of the expected head full block
}

func newCanonical(engine consensus.Engine, n int, full bool) (database.DBManager, *blockchain.BlockChain, error) {
	var (
		canonicalSeed = 1
		db            = database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
		// TODO: Remove me
		// db      = database.NewMemoryDBManager()
		genesis = new(blockchain.Genesis).MustCommit(db)
	)

	// Initialize a fresh chain with only a genesis block
	bc, _ := blockchain.NewBlockChain(db, nil, params.AllGxhashProtocolChanges, engine, vm.Config{})
	// Create and inject the requested chain
	if n == 0 {
		return db, bc, nil
	}
	if full {
		// Full block-chain requested
		blocks := blockchain.MakeBlockChain(genesis, n, engine, db, canonicalSeed)
		_, err := bc.InsertChain(blocks)
		return db, bc, err
	}
	// Header-only chain requested
	headers := blockchain.MakeHeaderChain(genesis.Header(), n, engine, db, canonicalSeed)
	_, err := bc.InsertHeaderChain(headers, 1)
	return db, bc, err
}

func expectedGovMap(t *testing.T, gov *governance.MixedEngine, num uint64, item string, value interface{}, expectedCacheSize int) {
	_, govMap, err := gov.ReadGovernance(num)
	assert.Nil(t, err)
	assert.Equal(t, govMap[item], value)
	assert.Equal(t, len(gov.IdxCache())-1, expectedCacheSize)
}

func TestSetHead(t *testing.T) {
	headerGovTest(t, &rewindTest{
		// `params.CheckpointInterval` is constant value of 1024.
		// Make it longer to include its working coverage
		canonicalBlocks:    1500,
		setheadBlock:       6,
		expCanonicalBlocks: 6,
		expHeadHeader:      6,
		expHeadFastBlock:   6,
		expHeadBlock:       6,
	})
}

func testCfg(epoch uint64) *params.ChainConfig {
	config := params.TestChainConfig.Copy()
	config.Governance = params.GetDefaultGovernanceConfig()
	config.Istanbul = params.GetDefaultIstanbulConfig()
	config.Istanbul.Epoch = epoch
	return config
}

func headerGovTest(t *testing.T, tt *rewindTest) {
	db, chain, err := newCanonical(gxhash.NewFullFaker(), 0, true)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer chain.Stop()

	_, _, _, api := newCNAPIBackend(t)

	var (
		epoch                 uint64 = 5
		govBlockNum                  = 10
		appliedGovBlockNum    uint64 = 20
		stakingUpdateInterval uint64 = 1
		stakingUpdateBlockNum uint64 = 15
		gov                          = governance.NewMixedEngine(testCfg(epoch), db)
		gpo                          = gasprice.NewOracle(api, gasprice.Config{}, nil)
	)
	chain.Config().Istanbul = &params.IstanbulConfig{Epoch: epoch, ProposerPolicy: params.WeightedRandom}

	canonblocks, _ := blockchain.GenerateChain(params.TestChainConfig, chain.CurrentBlock(), gxhash.NewFaker(), db, tt.canonicalBlocks, func(i int, b *blockchain.BlockGen) {
		if i == govBlockNum-1 { // Subtract 1, because the callback starts to enumerate from zero
			// "reward.mintingamount" = 123
			govData := hexutil.MustDecode("0x9e7b227265776172642e6d696e74696e67616d6f756e74223a22313233227d")
			b.SetGovData(govData)
			gov.WriteGovernanceForNextEpoch(uint64(govBlockNum), govData)
		}
	})

	if _, err := chain.InsertChain(canonblocks[:tt.canonicalBlocks]); err != nil {
		t.Fatalf("Failed to import canonical chain start: %v", err)
	}

	// Store snapshot
	snap := backend.Snapshot{Number: params.CheckpointInterval, Hash: chain.GetHeaderByNumber(params.CheckpointInterval).Hash()}
	blob, err := json.Marshal(snap)
	assert.Nil(t, err)
	db.WriteIstanbulSnapshot(snap.Hash, blob)
	_, err = db.ReadIstanbulSnapshot(snap.Hash)
	assert.Nil(t, err)

	// Initiailize staking info manager
	dummy := reward.StakingInfo{BlockNum: stakingUpdateInterval}
	blob, err = json.Marshal(dummy)
	assert.Nil(t, err)
	reward.SetTestStakingManagerWithStakingInfoCache(&dummy)
	assert.NotNil(t, reward.GetStakingManager())
	params.SetStakingUpdateInterval(stakingUpdateInterval)
	// Write a value to DB
	err = db.WriteStakingInfo(stakingUpdateBlockNum, blob)
	assert.Nil(t, err)
	_, err = db.ReadStakingInfo(stakingUpdateBlockNum)
	assert.Nil(t, err)
	assert.Equal(t, reward.TestGetStakingCacheSize(), 1)

	// Before setHead
	expectedGovMap(t, gov, appliedGovBlockNum, "reward.mintingamount", "123", 1)

	// Set the head of the chain back to the requested number
	err = doSetHead(chain, chain.Engine(), gov, gpo, tt.setheadBlock)
	assert.Nil(t, err)

	if head := chain.CurrentHeader(); head.Number.Uint64() != tt.expHeadHeader {
		t.Errorf("Head header mismatch!: have %d, want %d", head.Number, tt.expHeadHeader)
	}
	if head := chain.CurrentFastBlock(); head.NumberU64() != tt.expHeadFastBlock {
		t.Errorf("Head fast block mismatch: have %d, want %d", head.NumberU64(), tt.expHeadFastBlock)
	}
	if head := chain.CurrentBlock(); head.NumberU64() != tt.expHeadBlock {
		t.Errorf("Head block mismatch!!: have %d, want %d", head.NumberU64(), tt.expHeadBlock)
	}
	// After setHead
	// governance db and cachelookup
	expectedGovMap(t, gov, appliedGovBlockNum, "reward.mintingamount", "0", 0)

	// staking db and cache lookup
	assert.Equal(t, reward.TestGetStakingCacheSize(), 0)
	_, err = db.ReadStakingInfo(stakingUpdateBlockNum)
	assert.Equal(t, err.Error(), "data is not found with the given key")

	// snapshot db lookup
	_, err = db.ReadIstanbulSnapshot(snap.Hash)
	assert.Equal(t, err.Error(), "data is not found with the given key")

	for _, b := range canonblocks[tt.expCanonicalBlocks:] {
		if _, err := chain.InsertChain(types.Blocks{b}); err != nil {
			t.Fatalf("Failed to import canonical chain start: %v", err)
		}
		if len(b.Header().Governance) > 0 {
			assert.Equal(t, b.Header().Number.Uint64()%uint64(epoch), uint64(0))
			gov.WriteGovernanceForNextEpoch(uint64(govBlockNum), b.Header().Governance)
		}
	}
	if head := chain.CurrentBlock(); head.NumberU64() != uint64(tt.canonicalBlocks) {
		t.Errorf("Head block mismatch!!: have %d, want %d", head.NumberU64(), tt.expHeadBlock)
	}
	// After setHead and sync
	expectedGovMap(t, gov, appliedGovBlockNum, "reward.mintingamount", "123", 1)
}
