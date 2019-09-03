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
	"github.com/klaytn/klaytn/consensus"
	mocks2 "github.com/klaytn/klaytn/consensus/mocks"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var num1 = 20190902

func newMocks(t *testing.T) (*gomock.Controller, *mocks2.MockEngine, *mocks.MockBlockChain, *mocks.MockTxPool) {
	mockCtrl := gomock.NewController(t)
	mockEngine := mocks2.NewMockEngine(mockCtrl)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)

	return mockCtrl, mockEngine, mockBlockChain, mockTxPool
}

func newBlock(blockNum int) *types.Block {
	header := &types.Header{Number: big.NewInt(int64(blockNum))}
	header.Hash()
	return types.NewBlockWithHeader(header)
}

func TestNewProtocolManager(t *testing.T) {
	//1. If consensus.Engine returns an empty Protocol, NewProtocolManager throws an error.
	{
		mockCtrl, mockEngine, mockBlockChain, mockTxPool := newMocks(t)
		defer mockCtrl.Finish()

		block := newBlock(num1)
		mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)
		mockEngine.EXPECT().Protocol().Return(consensus.Protocol{}).Times(1)

		pm, err := NewProtocolManager(nil, downloader.FastSync, 0, nil, mockTxPool,
			mockEngine, mockBlockChain, nil, -1, &Config{})

		assert.Nil(t, pm)
		assert.Equal(t, errIncompatibleConfig, err)
	}
}

func TestSampleSize(t *testing.T) {
	peers := make([]Peer, minNumPeersToSendBlock-1)
	assert.Equal(t, len(peers), sampleSize(peers))

	peers = make([]Peer, 4)
	assert.Equal(t, minNumPeersToSendBlock, sampleSize(peers))

	peers = make([]Peer, 16)
	assert.Equal(t, 4, sampleSize(peers))
}

func TestSamplePeersToSendBlock(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	block := newBlock(123)
	hash := block.Hash()

	// 1. When the given node is a consensus node.
	{
		// len(pnsWithoutBlock) <= blockReceivingPNLimit
		mockPeers := NewMockPeerSet(mockCtrl)
		pm.nodetype = node.CONSENSUSNODE
		pm.peers = mockPeers

		mockCN := NewMockPeer(mockCtrl)
		mockPN := NewMockPeer(mockCtrl)

		mockPeers.EXPECT().CNWithoutBlock(hash).Return([]Peer{mockCN}).Times(1)
		mockPeers.EXPECT().PNWithoutBlock(hash).Return([]Peer{mockPN}).Times(1)

		sampledPeers := pm.samplePeersToSendBlock(block)
		assert.Equal(t, 2, len(sampledPeers))
		assert.Equal(t, mockCN, sampledPeers[0])
		assert.Equal(t, mockPN, sampledPeers[1])

		// len(pnsWithoutBlock) > blockReceivingPNLimit
		mockPeers = NewMockPeerSet(mockCtrl)
		pm.peers = mockPeers

		mockCN = NewMockPeer(mockCtrl)
		mockPN = NewMockPeer(mockCtrl)

		var pnWithoutBlock []Peer
		for i := 0; i <= blockReceivingPNLimit; i++ {
			pnWithoutBlock = append(pnWithoutBlock, mockPN)
		}

		mockPeers.EXPECT().CNWithoutBlock(hash).Return([]Peer{mockCN}).Times(1)
		mockPeers.EXPECT().PNWithoutBlock(hash).Return(pnWithoutBlock).Times(1)

		sampledPeers = pm.samplePeersToSendBlock(block)
		assert.Equal(t, 1+blockReceivingPNLimit, len(sampledPeers))
		assert.Equal(t, mockCN, sampledPeers[0])
	}

	// 2. When the given node is a proxy node.
	{
		mockPeers := NewMockPeerSet(mockCtrl)
		pm.nodetype = node.PROXYNODE
		pm.peers = mockPeers

		mockPeer := NewMockPeer(mockCtrl)
		mockPeers.EXPECT().PeersWithoutBlockExceptCN(hash).Return([]Peer{mockPeer}).Times(1)

		sampledPeers := pm.samplePeersToSendBlock(block)
		assert.Equal(t, 1, len(sampledPeers))
		assert.Equal(t, mockPeer, sampledPeers[0])
	}

	// 3. When the given node is an end point node.
	{
		mockPeers := NewMockPeerSet(mockCtrl)
		pm.nodetype = node.ENDPOINTNODE
		pm.peers = mockPeers

		mockPeer := NewMockPeer(mockCtrl)
		mockPeers.EXPECT().ENWithoutBlock(hash).Return([]Peer{mockPeer}).Times(1)

		sampledPeers := pm.samplePeersToSendBlock(block)
		assert.Equal(t, 1, len(sampledPeers))
		assert.Equal(t, mockPeer, sampledPeers[0])
	}

	// 4. When the given node is a boot node.
	{
		pm.nodetype = node.BOOTNODE
		assert.Equal(t, []Peer{}, pm.samplePeersToSendBlock(block))
	}

	// 5. When the given node is an unknown node.
	{
		pm.nodetype = node.UNKNOWNNODE
		assert.Equal(t, []Peer{}, pm.samplePeersToSendBlock(block))
	}
}

func TestSamplingPeers(t *testing.T) {
	peers := make([]Peer, 10)
	assert.Equal(t, peers, samplingPeers(peers, 20))
	assert.Equal(t, peers[:5], samplingPeers(peers, 5))
}

func TestBroadcastBlock_NoParentExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	block := newBlock(123)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().GetBlock(block.ParentHash(), block.NumberU64()-1).Return(nil).Times(1)

	pm.blockchain = mockBlockChain
	pm.BroadcastBlock(block)
}

func TestBroadcastBlock_ParentExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	block := newBlock(123)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().GetBlock(block.ParentHash(), block.NumberU64()-1).Return(block).Times(1)
	mockBlockChain.EXPECT().GetTd(block.ParentHash(), block.NumberU64()-1).Return(big.NewInt(100)).Times(1)
	pm.blockchain = mockBlockChain

	mockPeers := NewMockPeerSet(mockCtrl)
	pm.peers = mockPeers

	mockPeer := NewMockPeer(mockCtrl)
	mockPeers.EXPECT().ENWithoutBlock(block.Hash()).Return([]Peer{mockPeer}).Times(1)
	mockPeer.EXPECT().AsyncSendNewBlock(block, new(big.Int).Add(block.BlockScore(), big.NewInt(100)))

	pm.BroadcastBlock(block)
}

func TestBroadcastBlockHash(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	block := newBlock(123)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	{
		mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
		mockBlockChain.EXPECT().HasBlock(block.Hash(), block.NumberU64()).Return(false).Times(1)
		pm.blockchain = mockBlockChain
		pm.BroadcastBlockHash(block)
	}

	{
		mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
		mockBlockChain.EXPECT().HasBlock(block.Hash(), block.NumberU64()).Return(true) //.MinTimes(100)
		pm.blockchain = mockBlockChain

		mockPeer := NewMockPeer(mockCtrl)
		mockPeer.EXPECT().AsyncSendNewBlockHash(block).Times(1)

		mockPeers := NewMockPeerSet(mockCtrl)
		mockPeers.EXPECT().PeersWithoutBlock(block.Hash()).Return([]Peer{mockPeer}).Times(1)
		pm.peers = mockPeers

		pm.BroadcastBlockHash(block)
	}
}
