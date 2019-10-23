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
	"context"
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/work/mocks"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var cc = &params.ChainConfig{UnitPrice: 12345}

func prepareServiceChainAPIBackendTest(t *testing.T) (*ServiceChainAPIBackend, *gomock.Controller, *mocks.MockBlockChain, *mocks.MockTxPool) {
	mockCtrl := gomock.NewController(t)
	api := &ServiceChainAPIBackend{sc: &ServiceChain{chainConfig: cc}}
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	api.sc.blockchain = mockBlockChain

	mockTxPool := mocks.NewMockTxPool(mockCtrl)
	api.sc.txPool = mockTxPool

	return api, mockCtrl, mockBlockChain, mockTxPool
}

func TestServiceChainAPIBackend_GetNonceInCache(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	addr := addrs[0]
	expectedCachedNonce := uint64(123)
	expectedExist := true

	mockBlockChain.EXPECT().GetNonceInCache(addr).Return(expectedCachedNonce, expectedExist).AnyTimes()
	cachedNonce, exist := api.GetNonceInCache(addr)

	assert.Equal(t, expectedCachedNonce, cachedNonce)
	assert.Equal(t, expectedExist, exist)
}

func TestServiceChainAPIBackend_GetTxLookupInfoAndReceipt(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	rct := newReceipt(250000)

	expectedTx := tx1
	expectedHash := hashes[0]
	expectedBlockNum := uint64(123)
	expectedIndex := uint64(321)
	expectedReceipt := rct

	mockBlockChain.EXPECT().GetTxLookupInfoAndReceipt(hashes[0]).Return(tx1, hashes[0], uint64(123), uint64(321), rct).Times(1)

	tx, hash, blockNum, index, receipt := api.GetTxLookupInfoAndReceipt(context.Background(), hashes[0])

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedHash, hash)
	assert.Equal(t, expectedBlockNum, blockNum)
	assert.Equal(t, expectedIndex, index)
	assert.Equal(t, expectedReceipt, receipt)
}

func TestServiceChainAPIBackend_GetTxAndLookupInfoInCache(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expectedTx := tx1
	expectedHash := hashes[0]
	expectedBlockNum := uint64(123)
	expectedIndex := uint64(321)

	mockBlockChain.EXPECT().GetTxAndLookupInfoInCache(hashes[0]).Return(tx1, hashes[0], uint64(123), uint64(321)).Times(1)

	tx, hash, blockNum, index := api.GetTxAndLookupInfoInCache(hashes[0])

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedHash, hash)
	assert.Equal(t, expectedBlockNum, blockNum)
	assert.Equal(t, expectedIndex, index)
}

func TestServiceChainAPIBackend_GetBlockReceiptsInCache(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	rct1 := newReceipt(111111)
	rct2 := newReceipt(222222)
	rct3 := newReceipt(333333)

	receipts := types.Receipts{rct1, rct2, rct3}

	mockBlockChain.EXPECT().GetBlockReceiptsInCache(hashes[0]).Return(receipts).Times(1)
	returnedReceipts := api.GetBlockReceiptsInCache(hashes[0])

	assert.Equal(t, receipts, returnedReceipts)
}

func TestServiceChainAPIBackend_GetTxLookupInfoAndReceiptInCache(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	rct := newReceipt(250000)

	expectedTx := tx1
	expectedHash := hashes[0]
	expectedBlockNum := uint64(123)
	expectedIndex := uint64(321)
	expectedReceipt := rct

	mockBlockChain.EXPECT().GetTxLookupInfoAndReceiptInCache(hashes[0]).Return(tx1, hashes[0], uint64(123), uint64(321), rct).Times(1)

	tx, hash, blockNum, index, receipt := api.GetTxLookupInfoAndReceiptInCache(hashes[0])

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedHash, hash)
	assert.Equal(t, expectedBlockNum, blockNum)
	assert.Equal(t, expectedIndex, index)
	assert.Equal(t, expectedReceipt, receipt)
}

func TestServiceChainAPIBackend_ChainConfig(t *testing.T) {
	api, mockCtrl, _, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	assert.Equal(t, cc, api.ChainConfig())
}

func TestServiceChainAPIBackend_CurrentBlock(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	mockBlockChain.EXPECT().CurrentBlock().Return(blocks[0]).AnyTimes()
	assert.Equal(t, blocks[0], api.CurrentBlock())
}

func TestServiceChainAPIBackEnd_SetHead(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	headNumber := uint64(12345)
	mockBlockChain.EXPECT().SetHead(headNumber).Times(1)
	api.SetHead(headNumber)
}

func TestServiceChainAPIBackend_GetBlock(t *testing.T) {
	{
		api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
		mockBlockChain.EXPECT().GetBlockByHash(hashes[0]).Return(nil).Times(1)

		block, err := api.GetBlock(context.Background(), hashes[0])
		assert.Nil(t, block)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
		mockBlockChain.EXPECT().GetBlockByHash(hashes[0]).Return(blocks[0]).Times(1)

		block, err := api.GetBlock(context.Background(), hashes[0])
		assert.Equal(t, blocks[0], block)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestServiceChainAPIBackend_GetTxAndLookupInfo(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expectedTx := tx1
	expectedHash := hashes[0]
	expectedBlockNum := uint64(123)
	expectedIndex := uint64(321)

	mockBlockChain.EXPECT().GetTxAndLookupInfo(hashes[0]).Return(tx1, hashes[0], uint64(123), uint64(321)).Times(1)

	tx, hash, blockNum, index := api.GetTxAndLookupInfo(hashes[0])

	assert.Equal(t, expectedTx, tx)
	assert.Equal(t, expectedHash, hash)
	assert.Equal(t, expectedBlockNum, blockNum)
	assert.Equal(t, expectedIndex, index)
}

func TestServiceChainAPIBackend_GetBlockReceipts(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	rct1 := newReceipt(111111)
	rct2 := newReceipt(222222)
	rct3 := newReceipt(333333)

	receipts := types.Receipts{rct1, rct2, rct3}

	mockBlockChain.EXPECT().GetReceiptsByBlockHash(hashes[0]).Return(receipts).Times(1)
	returnedReceipts := api.GetBlockReceipts(context.Background(), hashes[0])

	assert.Equal(t, receipts, returnedReceipts)
}

func TestServiceChainAPIBackend_GetLogs(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expectedLogs := [][]*types.Log{{{BlockNumber: 123}, {BlockNumber: 321}}}
	mockBlockChain.EXPECT().GetLogsByHash(hashes[0]).Return(expectedLogs).Times(1)

	returnedLogs, err := api.GetLogs(context.Background(), hashes[0])
	assert.Equal(t, expectedLogs, returnedLogs)
	assert.NoError(t, err)
}

func TestServiceChainAPIBackend_GetTd(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expectedTD := big.NewInt(12345)
	mockBlockChain.EXPECT().GetTdByHash(hashes[0]).Return(expectedTD).Times(1)
	assert.Equal(t, expectedTD, api.GetTd(hashes[0]))
}

func TestServiceChainAPIBackend_SubscribeRemovedLogsEvent(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	ch := make(chan<- blockchain.RemovedLogsEvent, 10)

	mockBlockChain.EXPECT().SubscribeRemovedLogsEvent(ch).Times(1)
	api.SubscribeRemovedLogsEvent(ch)
}

func TestServiceChainAPIBackend_SubscribeChainEvent(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	ch := make(chan<- blockchain.ChainEvent, 10)

	mockBlockChain.EXPECT().SubscribeChainEvent(ch).Times(1)
	api.SubscribeChainEvent(ch)
}

func TestServiceChainAPIBackend_SubscribeChainHeadEvent(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	ch := make(chan<- blockchain.ChainHeadEvent, 10)

	mockBlockChain.EXPECT().SubscribeChainHeadEvent(ch).Times(1)
	api.SubscribeChainHeadEvent(ch)
}

func TestServiceChainAPIBackend_SubscribeChainSideEvent(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	ch := make(chan<- blockchain.ChainSideEvent, 10)

	mockBlockChain.EXPECT().SubscribeChainSideEvent(ch).Times(1)
	api.SubscribeChainSideEvent(ch)
}

func TestServiceChainAPIBackend_SubscribeLogsEvent(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	ch := make(chan<- []*types.Log, 10)

	mockBlockChain.EXPECT().SubscribeLogsEvent(ch).Times(1)
	api.SubscribeLogsEvent(ch)
}

func TestServiceChainAPIBackend_SendTx(t *testing.T) {
	{
		api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
		mockTxPool.EXPECT().AddLocal(tx1).Times(1).Return(nil)
		assert.NoError(t, api.SendTx(context.Background(), tx1))
		mockCtrl.Finish()
	}
	{
		api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
		mockTxPool.EXPECT().AddLocal(tx1).Times(1).Return(err)
		assert.Error(t, api.SendTx(context.Background(), tx1))
		mockCtrl.Finish()
	}
}

func TestServiceChainAPIBackend_GetPoolTransactions(t *testing.T) {
	{
		api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
		returnedTxs := map[common.Address]types.Transactions{addrs[0]: {tx1}}
		mockTxPool.EXPECT().Pending().Return(returnedTxs, err).Times(1)
		txs, err := api.GetPoolTransactions()
		assert.Nil(t, txs)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
		returnedTxs := map[common.Address]types.Transactions{addrs[0]: {tx1}}
		expectedTxs := types.Transactions{tx1}
		mockTxPool.EXPECT().Pending().Return(returnedTxs, nil).Times(1)
		txs, err := api.GetPoolTransactions()
		assert.Equal(t, expectedTxs, txs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestServiceChainAPIBackend_GetPoolTransaction(t *testing.T) {
	api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expectedTx := &*tx1
	mockTxPool.EXPECT().Get(hashes[0]).Return(tx1).Times(1)
	assert.Equal(t, expectedTx, api.GetPoolTransaction(hashes[0]))
}

func TestServiceChainAPIBackend_GetPoolNonce(t *testing.T) {
	api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expected := uint64(123)
	mockTxPool.EXPECT().GetPendingNonce(addrs[0]).Return(expected).Times(1)
	assert.Equal(t, expected, api.GetPoolNonce(context.Background(), addrs[0]))
}

func TestServiceChainAPIBackend_Stats(t *testing.T) {
	api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expectedPending := 123
	expectedQueued := 321

	mockTxPool.EXPECT().Stats().Return(expectedPending, expectedQueued).Times(1)
	pending, queued := api.Stats()

	assert.Equal(t, expectedPending, pending)
	assert.Equal(t, expectedQueued, queued)
}

func TestServiceChainAPIBackend_TxPoolContent(t *testing.T) {
	api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	expectedPending := map[common.Address]types.Transactions{addrs[0]: {tx1}}
	expectedQueued := map[common.Address]types.Transactions{addrs[1]: {tx1}}

	mockTxPool.EXPECT().Content().Return(expectedPending, expectedQueued).Times(1)

	returnedPending, returnedQueued := api.TxPoolContent()
	assert.Equal(t, expectedPending, returnedPending)
	assert.Equal(t, expectedQueued, returnedQueued)
}

func TestServiceChainAPIBackend_SubscribeNewTxsEvent(t *testing.T) {
	api, mockCtrl, _, mockTxPool := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	ch := make(chan<- blockchain.NewTxsEvent, 10)
	sub := &filters.Subscription{ID: rpc.NewID()}
	mockTxPool.EXPECT().SubscribeNewTxsEvent(ch).Return(sub).Times(1)
	assert.Equal(t, sub, api.SubscribeNewTxsEvent(ch))
}

func TestServiceChainAPIBackend_IsParallelDBWrite(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	mockBlockChain.EXPECT().IsParallelDBWrite().Return(true).Times(1)
	assert.True(t, api.IsParallelDBWrite())
}

func TestServiceChainAPIBackend_IsSenderTxHashIndexingEnabled(t *testing.T) {
	api, mockCtrl, mockBlockChain, _ := prepareServiceChainAPIBackendTest(t)
	defer mockCtrl.Finish()

	mockBlockChain.EXPECT().IsSenderTxHashIndexingEnabled().Return(true).Times(1)
	assert.True(t, api.IsSenderTxHashIndexingEnabled())
}
