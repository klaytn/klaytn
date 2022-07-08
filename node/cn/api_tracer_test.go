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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common/hexutil"
	mocks2 "github.com/klaytn/klaytn/consensus/mocks"
	"github.com/klaytn/klaytn/networks/rpc"
	mocks3 "github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/klaytn/klaytn/work/mocks"
	"github.com/stretchr/testify/assert"
)

func createCNMocks(t *testing.T) (*gomock.Controller, *PrivateDebugAPI, *mocks2.MockEngine, *mocks.MockBlockChain, *mocks3.MockMiner) {
	mockCtrl, mockEngine, mockBlockChain, _ := newMocks(t)
	mockMiner := mocks3.NewMockMiner(mockCtrl)
	api := NewPrivateDebugAPI(nil, &CN{miner: mockMiner, blockchain: mockBlockChain, engine: mockEngine})
	return mockCtrl, api, mockEngine, mockBlockChain, mockMiner
}

func TestPrivateDebugAPI_TraceChain(t *testing.T) {
	endBlockNumber := rpc.BlockNumber(123)
	block := newBlock(123)
	// from == nil
	{
		mockCtrl, api, _, _, mockMiner := createCNMocks(t)
		mockMiner.EXPECT().PendingBlock().Return(nil).Times(2)
		sub, err := api.TraceChain(context.Background(), rpc.PendingBlockNumber, rpc.PendingBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().CurrentBlock().Return(nil).Times(2)
		sub, err := api.TraceChain(context.Background(), rpc.LatestBlockNumber, rpc.LatestBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().GetBlockByNumber(uint64(endBlockNumber)).Return(nil).Times(2)
		sub, err := api.TraceChain(context.Background(), endBlockNumber, endBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	// from != nil
	{
		mockCtrl, api, _, mockBlockChain, mockMiner := createCNMocks(t)
		mockMiner.EXPECT().PendingBlock().Return(block).Times(1)
		mockBlockChain.EXPECT().CurrentBlock().Return(nil).Times(1)
		sub, err := api.TraceChain(context.Background(), rpc.PendingBlockNumber, rpc.LatestBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		endBlockNumber = rpc.BlockNumber(123)
		mockCtrl, api, _, mockBlockChain, mockMiner := createCNMocks(t)
		mockMiner.EXPECT().PendingBlock().Return(block).Times(1)
		mockBlockChain.EXPECT().GetBlockByNumber(uint64(endBlockNumber)).Return(nil).Times(1)
		sub, err := api.TraceChain(context.Background(), rpc.PendingBlockNumber, endBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		startBlockNumber := rpc.BlockNumber(130)
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().GetBlockByNumber(uint64(startBlockNumber)).Return(newBlock(130))
		mockBlockChain.EXPECT().GetBlockByNumber(uint64(endBlockNumber)).Return(newBlock(123))
		sub, err := api.TraceChain(context.Background(), startBlockNumber, endBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Equal(t, fmt.Errorf("end block #%d needs to come after start block #%d", endBlockNumber, startBlockNumber), err)
		mockCtrl.Finish()
	}
}

func TestPrivateDebugAPI_TraceBlockByNumber(t *testing.T) {
	blockNumber := rpc.BlockNumber(123)
	{
		mockCtrl, api, _, _, mockMiner := createCNMocks(t)
		mockMiner.EXPECT().PendingBlock().Return(nil).Times(1)
		sub, err := api.TraceBlockByNumber(context.Background(), rpc.PendingBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().CurrentBlock().Return(nil).Times(1)
		sub, err := api.TraceBlockByNumber(context.Background(), rpc.LatestBlockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().GetBlockByNumber(uint64(blockNumber)).Return(nil).Times(1)
		sub, err := api.TraceBlockByNumber(context.Background(), blockNumber, nil)
		assert.Nil(t, sub)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
}

func TestPrivateDebugAPI_TraceBlockByHash(t *testing.T) {
	mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
	mockBlockChain.EXPECT().GetBlockByHash(hashes[0]).Return(nil).Times(1)
	sub, err := api.TraceBlockByHash(context.Background(), hashes[0], nil)
	assert.Nil(t, sub)
	assert.Error(t, err)
	mockCtrl.Finish()
}

func TestPrivateDebugAPI_TraceBlock(t *testing.T) {
	mockCtrl, api, _, _, _ := createCNMocks(t)
	sub, err := api.TraceBlock(context.Background(), hexutil.Bytes{}, nil)
	assert.Nil(t, sub)
	assert.Error(t, err)
	mockCtrl.Finish()
}

func TestPrivateDebugAPI_TraceBlockFromFile(t *testing.T) {
	mockCtrl, api, _, _, _ := createCNMocks(t)
	sub, err := api.TraceBlockFromFile(context.Background(), "12345", nil)
	assert.Nil(t, sub)
	assert.Error(t, err)
	mockCtrl.Finish()
}

func TestPrivateDebugAPI_TraceBadBlock(t *testing.T) {
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().BadBlocks().Return(nil, expectedErr).Times(1)
		sub, returnedErr := api.TraceBadBlock(context.Background(), hashes[0], nil)
		assert.Nil(t, sub)
		assert.Equal(t, expectedErr, returnedErr)
		mockCtrl.Finish()
	}
	{
		block := newBlock(123)
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().BadBlocks().Return([]blockchain.BadBlockArgs{{Hash: block.Hash(), Block: block}}, expectedErr).Times(1)
		sub, returnedErr := api.TraceBadBlock(context.Background(), hashes[0], nil)
		assert.Nil(t, sub)
		assert.Equal(t, expectedErr, returnedErr)
		mockCtrl.Finish()
	}
}

func TestPrivateDebugAPI_StandardTraceBlockToFile(t *testing.T) {
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().GetBlockByHash(hashes[0]).Return(nil).Times(1)
		sub, returnedErr := api.StandardTraceBlockToFile(context.Background(), hashes[0], nil)
		assert.Nil(t, sub)
		assert.Error(t, returnedErr)
		mockCtrl.Finish()
	}
}

func TestPrivateDebugAPI_StandardTraceBadBlockToFile(t *testing.T) {
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().BadBlocks().Return(nil, expectedErr).Times(1)
		sub, returnedErr := api.StandardTraceBadBlockToFile(context.Background(), hashes[0], nil)
		assert.Nil(t, sub)
		assert.Equal(t, expectedErr, returnedErr)
		mockCtrl.Finish()
	}
	{
		block := newBlock(123)
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().BadBlocks().Return([]blockchain.BadBlockArgs{{Hash: block.Hash(), Block: block}}, expectedErr).Times(1)
		sub, returnedErr := api.StandardTraceBadBlockToFile(context.Background(), hashes[0], nil)
		assert.Nil(t, sub)
		assert.Equal(t, expectedErr, returnedErr)
		mockCtrl.Finish()
	}
}

func TestPrivateDebugAPI_computeTxEnv(t *testing.T) {
	parentBlock := newBlock(122)
	block := newBlockWithParentHash(123, parentBlock.Hash())
	blockHash := block.Hash()
	txIndex := 0
	reexec := uint64(0)
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().GetBlockByHash(blockHash).Return(nil).Times(1)
		msg, ctx, stateDB, err := api.computeTxEnv(blockHash, txIndex, reexec)
		assert.Nil(t, msg)
		assert.Equal(t, vm.Context{}, ctx)
		assert.Nil(t, stateDB)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, api, _, mockBlockChain, _ := createCNMocks(t)
		mockBlockChain.EXPECT().GetBlockByHash(blockHash).Return(block).Times(1)
		mockBlockChain.EXPECT().GetBlock(block.ParentHash(), block.NumberU64()-1).Return(nil).Times(1)
		msg, ctx, stateDB, err := api.computeTxEnv(blockHash, txIndex, reexec)
		assert.Nil(t, msg)
		assert.Equal(t, vm.Context{}, ctx)
		assert.Nil(t, stateDB)
		assert.Error(t, err)
		mockCtrl.Finish()
	}
}
