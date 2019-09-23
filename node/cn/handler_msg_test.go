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
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/networks/p2p"
	mocks2 "github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/work/mocks"
	"github.com/stretchr/testify/assert"
	"math/big"
	"sync/atomic"
	"testing"
)

// generateMsg creates a message struct for message handling tests.
func generateMsg(t *testing.T, msgCode uint64, data interface{}) p2p.Msg {
	size, r, err := rlp.EncodeToReader(data)
	if err != nil {
		t.Fatal(err)
	}
	return p2p.Msg{
		Code:    msgCode,
		Size:    uint32(size),
		Payload: r,
	}
}

// prepareTestHandleNewBlockMsg creates structs for TestHandleNewBlockMsg_ tests.
func prepareTestHandleNewBlockMsg(t *testing.T, mockCtrl *gomock.Controller, blockNum int) (*types.Block, p2p.Msg, *MockPeer, *mocks2.MockProtocolManagerFetcher) {
	mockPeer := NewMockPeer(mockCtrl)

	newBlock := newBlock(blockNum)
	newBlock.ReceivedFrom = mockPeer
	msg := generateMsg(t, NewBlockMsg, newBlockData{Block: newBlock, TD: big.NewInt(int64(blockNum))})

	mockPeer.EXPECT().AddToKnownBlocks(newBlock.Hash()).Times(1)
	mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()

	mockFetcher := mocks2.NewMockProtocolManagerFetcher(mockCtrl)
	mockFetcher.EXPECT().Enqueue(nodeids[0].String(), newBlock).Times(1)

	return newBlock, msg, mockPeer, mockFetcher
}

func TestHandleNewBlockMsg_LargeLocalPeerBlockScore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	_, msg, mockPeer, mockFetcher := prepareTestHandleNewBlockMsg(t, mockCtrl, blockNum1)

	pm := &ProtocolManager{}
	pm.fetcher = mockFetcher

	mockPeer.EXPECT().Head().Return(hash1, big.NewInt(blockNum1+1)).AnyTimes()

	assert.NoError(t, handleNewBlockMsg(pm, mockPeer, msg))
}

func TestHandleNewBlockMsg_SmallLocalPeerBlockScore_NoSynchronise(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	block, msg, mockPeer, mockFetcher := prepareTestHandleNewBlockMsg(t, mockCtrl, blockNum1)

	pm := &ProtocolManager{}
	pm.fetcher = mockFetcher

	mockPeer.EXPECT().Head().Return(hash1, big.NewInt(blockNum1-2)).AnyTimes()
	mockPeer.EXPECT().SetHead(block.ParentHash(), big.NewInt(blockNum1-1)).Times(1)

	currBlock := newBlock(blockNum1 - 1)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().CurrentBlock().Return(currBlock).Times(1)
	mockBlockChain.EXPECT().GetTd(currBlock.Hash(), currBlock.NumberU64()).Return(big.NewInt(blockNum1)).Times(1)

	pm.blockchain = mockBlockChain

	assert.NoError(t, handleNewBlockMsg(pm, mockPeer, msg))
}

func TestHandleTxMsg(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPeer := NewMockPeer(mockCtrl)

	txs := types.Transactions{tx1}
	msg := generateMsg(t, TxMsg, txs)

	// If pm.acceptTxs == 0, nothing happens.
	{
		assert.NoError(t, handleTxMsg(pm, mockPeer, msg))
	}
	// If pm.acceptTxs == 1, TxPool.HandleTxMsg is called.
	{
		atomic.StoreUint32(&pm.acceptTxs, 1)
		mockTxPool := mocks.NewMockTxPool(mockCtrl)
		mockTxPool.EXPECT().HandleTxMsg(gomock.Eq(txs)).AnyTimes()
		pm.txpool = mockTxPool

		mockPeer.EXPECT().AddToKnownTxs(txs[0].Hash()).Times(1)
		assert.NoError(t, handleTxMsg(pm, mockPeer, msg))
	}
}
