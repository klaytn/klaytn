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
	"github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var num1 = uint64(20190902)

func newMocks(t *testing.T) (*mocks2.MockEngine, *mocks.MockBlockChain, *mocks.MockTxPool) {
	mockCtrl := gomock.NewController(t)
	mockEngine := mocks2.NewMockEngine(mockCtrl)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)

	return mockEngine, mockBlockChain, mockTxPool
}

func TestNewProtocolManager(t *testing.T) {
	//1. If consensus.Engine returns an empty Protocol, NewProtocolManager throws an error.
	{
		mockEngine, mockBlockChain, mockTxPool := newMocks(t)

		header := &types.Header{Number: big.NewInt(int64(num1))}
		header.Hash()
		block := types.NewBlockWithHeader(header)

		mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(2)
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
