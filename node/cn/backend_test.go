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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/klaytn/klaytn/params"
	mocks2 "github.com/klaytn/klaytn/work/mocks"
	"github.com/stretchr/testify/assert"
)

func newCN(t *testing.T) (*gomock.Controller, *MockBackendProtocolManager, *mocks.MockMiner, *CN) {
	mockCtrl := gomock.NewController(t)
	mockProtocolManager := NewMockBackendProtocolManager(mockCtrl)
	mockMiner := mocks.NewMockMiner(mockCtrl)
	return mockCtrl, mockProtocolManager, mockMiner,
		&CN{protocolManager: mockProtocolManager, miner: mockMiner}
}

func TestCN_AddLesServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cn := &CN{}
	les := mocks.NewMockLesServer(mockCtrl)
	les.EXPECT().SetBloomBitsIndexer(nil).Times(1)
	cn.AddLesServer(les)
}

func TestCN_CheckSyncMode(t *testing.T) {
	c := &Config{SyncMode: downloader.FastSync}
	assert.NoError(t, checkSyncMode(c))

	c.SyncMode = downloader.FullSync
	assert.NoError(t, checkSyncMode(c))

	c.SyncMode = downloader.LightSync
	assert.Equal(t, errCNLightSync, checkSyncMode(c))
}

func TestCN_SetEngineType(t *testing.T) {
	cc := &params.ChainConfig{}
	originalEngineType := types.EngineType

	setEngineType(cc)
	assert.Equal(t, originalEngineType, types.EngineType)

	cc.Clique = &params.CliqueConfig{}
	setEngineType(cc)
	assert.Equal(t, types.Engine_Clique, types.EngineType)

	cc.Istanbul = &params.IstanbulConfig{}
	setEngineType(cc)
	assert.Equal(t, types.Engine_IBFT, types.EngineType)
}

func TestCN_SetAcceptTxs(t *testing.T) {
	{
		mockCtrl, _, _, cn := newCN(t)
		cn.chainConfig = &params.ChainConfig{}
		assert.NoError(t, cn.setAcceptTxs())
		mockCtrl.Finish()
	}
}

func TestCN_ResetWithGenesisBlock(t *testing.T) {
	mockCtrl, _, _, cn := newCN(t)
	defer mockCtrl.Finish()
	mockBlockChain := mocks2.NewMockBlockChain(mockCtrl)
	cn.blockchain = mockBlockChain
	block := blocks[0]

	mockBlockChain.EXPECT().ResetWithGenesisBlock(block).Times(1)
	cn.ResetWithGenesisBlock(block)
}

func TestCN_Rewardbase(t *testing.T) {
	{
		mockCtrl, _, _, cn := newCN(t)
		cn.rewardbase = addrs[0]

		rb, err := cn.Rewardbase()
		assert.Equal(t, addrs[0], rb)
		assert.Nil(t, err)

		mockCtrl.Finish()
	}
}

func TestCN_StartMining(t *testing.T) {
	{
		mockCtrl, _, mockMiner, cn := newCN(t)
		mockMiner.EXPECT().Start().Times(1)
		assert.Nil(t, cn.StartMining(false))
		time.Sleep(100 * time.Millisecond)
		mockCtrl.Finish()
	}
	{
		mockCtrl, mockPM, mockMiner, cn := newCN(t)
		mockPM.EXPECT().SetAcceptTxs().Times(1)
		mockMiner.EXPECT().Start().Times(1)
		assert.Nil(t, cn.StartMining(true))
		time.Sleep(100 * time.Millisecond)
		mockCtrl.Finish()
	}
}

func TestCN_StopMining(t *testing.T) {
	mockCtrl, _, mockMiner, cn := newCN(t)
	mockMiner.EXPECT().Stop().Times(1)
	cn.StopMining()
	mockCtrl.Finish()
}

func TestCN_IsMining(t *testing.T) {
	mockCtrl, _, mockMiner, cn := newCN(t)
	mockMiner.EXPECT().Mining().Times(1)
	cn.IsMining()
	mockCtrl.Finish()
}

func TestCN_ReBroadcastTxs(t *testing.T) {
	mockCtrl, mockPM, _, cn := newCN(t)
	defer mockCtrl.Finish()
	txs := types.Transactions{tx1}
	mockPM.EXPECT().ReBroadcastTxs(txs).Times(1)
	cn.ReBroadcastTxs(txs)
}
