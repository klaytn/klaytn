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
	"github.com/klaytn/klaytn/consensus/istanbul/backend"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestChannelManager tests registering and retrieving of message channels.
func TestChannelManager(t *testing.T) {
	cm := NewChannelManager(1)
	cm.RegisterMsgCode(ConsensusChannel, backend.IstanbulMsg)

	channel := make(chan p2p.Msg, channelSizePerPeer)
	consensusChannel := make(chan p2p.Msg, channelSizePerPeer)

	// Before calling RegisterChannelWithIndex,
	// calling GetChannelWithMsgCode with registered MsgCode should return no channel and no error.
	for i := NewBlockHashesMsg; i < MsgCodeEnd; i++ {
		ch, err := cm.GetChannelWithMsgCode(0, uint64(i))
		assert.Nil(t, ch)
		assert.NoError(t, err)
	}
	ch, err := cm.GetChannelWithMsgCode(0, backend.IstanbulMsg)
	assert.Nil(t, ch)
	assert.NoError(t, err)

	// Before calling RegisterChannelWithIndex,
	// calling GetChannelWithMsgCode with not-registered MsgCode should return no channel but an error.
	ch, err = cm.GetChannelWithMsgCode(0, 0xff)
	assert.Nil(t, ch)
	assert.Error(t, err)

	// Register channels with the port index.
	cm.RegisterChannelWithIndex(0, BlockChannel, channel)
	cm.RegisterChannelWithIndex(0, TxChannel, channel)
	cm.RegisterChannelWithIndex(0, MiscChannel, channel)
	cm.RegisterChannelWithIndex(0, ConsensusChannel, consensusChannel)

	// After calling RegisterChannelWithIndex,
	// calling GetChannelWithMsgCode with registered MsgCode should return a channel but no error.
	for i := NewBlockHashesMsg; i < MsgCodeEnd; i++ {
		ch, err := cm.GetChannelWithMsgCode(0, uint64(i))
		assert.Equal(t, channel, ch)
		assert.NoError(t, err)
	}
	ch, err = cm.GetChannelWithMsgCode(0, backend.IstanbulMsg)
	assert.Equal(t, consensusChannel, ch)
	assert.NoError(t, err)

	// After calling RegisterChannelWithIndex,
	// calling GetChannelWithMsgCode with not-registered MsgCode should return no channel but an error.
	ch, err = cm.GetChannelWithMsgCode(0, 0xff)
	assert.Nil(t, ch)
	assert.Error(t, err)
}
