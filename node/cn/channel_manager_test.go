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

	"github.com/klaytn/klaytn/consensus/istanbul/backend"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/stretchr/testify/assert"
)

// TestChannelManager_ChannelSize_1 tests registering and retrieving of
// message channels with the channel size 1.
func TestChannelManager_ChannelSize_1(t *testing.T) {
	testChannelManager(t, 1)
}

// TestChannelManager_ChannelSize_3 tests registering and retrieving of
// message channels with the channel size 3.
func TestChannelManager_ChannelSize_3(t *testing.T) {
	testChannelManager(t, 3)
}

// TestChannelManager_ChannelSize_5 tests registering and retrieving of
// message channels with the channel size 5.
func TestChannelManager_ChannelSize_5(t *testing.T) {
	testChannelManager(t, 5)
}

// testChannelManager tests registering and retrieving of
// message channels with the given channel size.
func testChannelManager(t *testing.T, chSize int) {
	cm := NewChannelManager(chSize)
	cm.RegisterMsgCode(ConsensusChannel, backend.IstanbulMsg)

	channel := make(chan p2p.Msg, channelSizePerPeer)
	consensusChannel := make(chan p2p.Msg, channelSizePerPeer)

	for chIdx := 0; chIdx < chSize; chIdx++ {
		// Before calling RegisterChannelWithIndex,
		// calling GetChannelWithMsgCode with registered MsgCode should return no channel and no error.
		for i := StatusMsg; i < MsgCodeEnd; i++ {
			ch, err := cm.GetChannelWithMsgCode(chIdx, uint64(i))
			assert.Nil(t, ch)
			assert.NoError(t, err)
		}
		ch, err := cm.GetChannelWithMsgCode(chIdx, backend.IstanbulMsg)
		assert.Nil(t, ch)
		assert.NoError(t, err)

		// Before calling RegisterChannelWithIndex,
		// calling GetChannelWithMsgCode with not-registered MsgCode should return no channel but an error.
		ch, err = cm.GetChannelWithMsgCode(chIdx, MsgCodeEnd)
		assert.Nil(t, ch)
		assert.Error(t, err)
	}

	// Register channels with the port index.
	for chIdx := 0; chIdx < chSize; chIdx++ {
		cm.RegisterChannelWithIndex(chIdx, BlockChannel, channel)
		cm.RegisterChannelWithIndex(chIdx, TxChannel, channel)
		cm.RegisterChannelWithIndex(chIdx, MiscChannel, channel)
		cm.RegisterChannelWithIndex(chIdx, ConsensusChannel, consensusChannel)
	}

	for chIdx := 0; chIdx < chSize; chIdx++ {
		// After calling RegisterChannelWithIndex,
		// calling GetChannelWithMsgCode with registered MsgCode should return a channel but no error.
		for i := StatusMsg; i < MsgCodeEnd; i++ {
			ch, err := cm.GetChannelWithMsgCode(chIdx, uint64(i))
			assert.Equal(t, channel, ch)
			assert.NoError(t, err)
		}
		ch, err := cm.GetChannelWithMsgCode(chIdx, backend.IstanbulMsg)
		assert.Equal(t, consensusChannel, ch)
		assert.NoError(t, err)

		// After calling RegisterChannelWithIndex,
		// calling GetChannelWithMsgCode with not-registered MsgCode should return no channel but an error.
		ch, err = cm.GetChannelWithMsgCode(0, MsgCodeEnd)
		assert.Nil(t, ch)
		assert.Error(t, err)
	}
}
