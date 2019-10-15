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
	"errors"
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

var err = errors.New("some error")

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

func prepareTestHandleBlockHeaderFetchRequestMsg(t *testing.T) (*gomock.Controller, *MockPeer, *mocks.MockBlockChain, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockPeer := NewMockPeer(mockCtrl)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)

	return mockCtrl, mockPeer, mockBlockChain, &ProtocolManager{blockchain: mockBlockChain}
}

func TestHandleBlockHeaderFetchRequestMsg(t *testing.T) {
	// Decoding the message failed, an error is returned.
	{
		mockCtrl, mockPeer, _, pm := prepareTestHandleBlockHeaderFetchRequestMsg(t)

		msg := generateMsg(t, BlockHeaderFetchRequestMsg, newBlock(blockNum1)) // use message data as a block, not a hash

		assert.Error(t, handleBlockHeaderFetchRequestMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// GetHeaderByHash returns nil, an error is returned.
	{
		mockCtrl, mockPeer, mockBlockChain, pm := prepareTestHandleBlockHeaderFetchRequestMsg(t)
		mockBlockChain.EXPECT().GetHeaderByHash(hash1).Return(nil).AnyTimes()
		mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchRequestMsg, hash1)

		assert.Error(t, handleBlockHeaderFetchRequestMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// GetHeaderByHash returns a header, p.SendFetchedBlockHeader(header) should be called.
	{
		mockCtrl, mockPeer, mockBlockChain, pm := prepareTestHandleBlockHeaderFetchRequestMsg(t)

		header := newBlock(blockNum1).Header()

		mockBlockChain.EXPECT().GetHeaderByHash(hash1).Return(header).AnyTimes()
		mockPeer.EXPECT().SendFetchedBlockHeader(header).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchRequestMsg, hash1)
		assert.NoError(t, handleBlockHeaderFetchRequestMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
}

func prepareTestHandleBlockHeaderFetchResponseMsg(t *testing.T) (*gomock.Controller, *MockPeer, *mocks2.MockProtocolManagerFetcher, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockPeer := NewMockPeer(mockCtrl)

	mockFetcher := mocks2.NewMockProtocolManagerFetcher(mockCtrl)
	pm := &ProtocolManager{fetcher: mockFetcher}

	return mockCtrl, mockPeer, mockFetcher, pm
}

func TestHandleBlockHeaderFetchResponseMsg(t *testing.T) {
	header := newBlock(blockNum1).Header()
	// Decoding the message failed, an error is returned.
	{
		mockCtrl := gomock.NewController(t)
		mockPeer := NewMockPeer(mockCtrl)
		pm := &ProtocolManager{}
		msg := generateMsg(t, BlockHeaderFetchResponseMsg, newBlock(blockNum1)) // use message data as a block, not a header
		assert.Error(t, handleBlockHeaderFetchResponseMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// FilterHeaders returns nil, error is not returned.
	{
		mockCtrl, mockPeer, mockFetcher, pm := prepareTestHandleBlockHeaderFetchResponseMsg(t)
		mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()
		mockFetcher.EXPECT().FilterHeaders(nodeids[0].String(), gomock.Eq([]*types.Header{header}), gomock.Any()).Return(nil).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchResponseMsg, header)
		assert.NoError(t, handleBlockHeaderFetchResponseMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// FilterHeaders returns not-nil, peer.GetID() is called twice to leave a log.
	{
		mockCtrl, mockPeer, mockFetcher, pm := prepareTestHandleBlockHeaderFetchResponseMsg(t)
		mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()
		mockFetcher.EXPECT().FilterHeaders(nodeids[0].String(), gomock.Eq([]*types.Header{header}), gomock.Any()).Return([]*types.Header{header}).AnyTimes()

		msg := generateMsg(t, BlockHeaderFetchResponseMsg, header)
		assert.NoError(t, handleBlockHeaderFetchResponseMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
}

func preparePeerAndDownloader(t *testing.T) (*gomock.Controller, *MockPeer, *mocks2.MockProtocolManagerDownloader, *ProtocolManager) {
	mockCtrl := gomock.NewController(t)
	mockPeer := NewMockPeer(mockCtrl)
	mockPeer.EXPECT().GetID().Return(nodeids[0].String()).AnyTimes()

	mockDownloader := mocks2.NewMockProtocolManagerDownloader(mockCtrl)
	pm := &ProtocolManager{downloader: mockDownloader}

	return mockCtrl, mockPeer, mockDownloader, pm
}

func TestHandleReceiptMsg(t *testing.T) {
	// Decoding the message failed, an error is returned.
	{
		mockCtrl := gomock.NewController(t)
		mockPeer := NewMockPeer(mockCtrl)
		pm := &ProtocolManager{}
		msg := generateMsg(t, ReceiptsMsg, newBlock(blockNum1)) // use message data as a block, not a header
		assert.Error(t, handleReceiptsMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// DeliverReceipts returns nil, error is not returned.
	{
		receipts := make([][]*types.Receipt, 1)
		receipts[0] = []*types.Receipt{newReceipt(123)}

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverReceipts(nodeids[0].String(), gomock.Eq(receipts)).Times(1).Return(nil)

		msg := generateMsg(t, ReceiptsMsg, receipts)
		assert.NoError(t, handleReceiptsMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// DeliverReceipts returns an error, but the error is not returned.
	{
		receipts := make([][]*types.Receipt, 1)
		receipts[0] = []*types.Receipt{newReceipt(123)}

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverReceipts(nodeids[0].String(), gomock.Eq(receipts)).Times(1).Return(err)

		msg := generateMsg(t, ReceiptsMsg, receipts)
		assert.NoError(t, handleReceiptsMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
}

func TestHandleNodeDataMsg(t *testing.T) {
	// Decoding the message failed, an error is returned.
	{
		mockCtrl := gomock.NewController(t)
		mockPeer := NewMockPeer(mockCtrl)
		pm := &ProtocolManager{}
		msg := generateMsg(t, NodeDataMsg, newBlock(blockNum1)) // use message data as a block, not a node data
		assert.Error(t, handleNodeDataMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// DeliverNodeData returns nil, error is not returned.
	{
		nodeData := make([][]byte, 1)
		nodeData[0] = hash1[:]

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverNodeData(nodeids[0].String(), gomock.Eq(nodeData)).Times(1).Return(nil)

		msg := generateMsg(t, ReceiptsMsg, nodeData)
		assert.NoError(t, handleNodeDataMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
	// DeliverNodeData returns an error, but the error is not returned.
	{
		nodeData := make([][]byte, 1)
		nodeData[0] = hash1[:]

		mockCtrl, mockPeer, mockDownloader, pm := preparePeerAndDownloader(t)
		mockDownloader.EXPECT().DeliverNodeData(nodeids[0].String(), gomock.Eq(nodeData)).Times(1).Return(err)

		msg := generateMsg(t, ReceiptsMsg, nodeData)
		assert.NoError(t, handleNodeDataMsg(pm, mockPeer, msg))
		mockCtrl.Finish()
	}
}
