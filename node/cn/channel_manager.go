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
	"fmt"

	"github.com/klaytn/klaytn/networks/p2p"
)

const (
	BlockChannel uint = iota
	TxChannel
	ConsensusChannel
	MiscChannel
	MaxChannel
)

type ChannelManager struct {
	msgChannels [][]chan p2p.Msg
	msgCodes    map[uint64]uint
}

// NewChannelManager returns a new ChannelManager.
// The ChannelManager manages the channel for the msg code.
func NewChannelManager(channelSize int) *ChannelManager {
	channelMgr := &ChannelManager{
		msgChannels: make([][]chan p2p.Msg, 0, channelSize),
		msgCodes:    make(map[uint64]uint),
	}

	for i := 0; i < channelSize; i++ {
		channelMgr.msgChannels = append(channelMgr.msgChannels, make([]chan p2p.Msg, MaxChannel, MaxChannel))
	}

	// register channelType to MsgCode
	channelMgr.RegisterMsgCode(BlockChannel, NewBlockHashesMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockHeaderFetchRequestMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockHeaderFetchResponseMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockBodiesFetchRequestMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockBodiesFetchResponseMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockHeadersRequestMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockHeadersMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockBodiesRequestMsg)
	channelMgr.RegisterMsgCode(BlockChannel, BlockBodiesMsg)
	channelMgr.RegisterMsgCode(BlockChannel, NewBlockMsg)

	channelMgr.RegisterMsgCode(TxChannel, TxMsg)

	channelMgr.RegisterMsgCode(MiscChannel, ReceiptsRequestMsg)
	channelMgr.RegisterMsgCode(MiscChannel, ReceiptsMsg)
	channelMgr.RegisterMsgCode(MiscChannel, StatusMsg)
	channelMgr.RegisterMsgCode(MiscChannel, NodeDataRequestMsg)
	channelMgr.RegisterMsgCode(MiscChannel, NodeDataMsg)

	return channelMgr
}

// RegisterChannelWithIndex registers the channel corresponding to network and channel ID.
func (cm *ChannelManager) RegisterChannelWithIndex(idx int, channelId uint, channel chan p2p.Msg) {
	cm.msgChannels[idx][channelId] = channel
}

// RegisterMsgCode registers the channel id corresponding to msgCode.
func (cm *ChannelManager) RegisterMsgCode(channelId uint, msgCode uint64) {
	cm.msgCodes[msgCode] = channelId
}

// GetChannelWithMsgCode returns the channel corresponding to msgCode.
func (cm *ChannelManager) GetChannelWithMsgCode(idx int, msgCode uint64) (chan p2p.Msg, error) {
	if channelID, ok := cm.msgCodes[msgCode]; ok {
		return cm.msgChannels[idx][channelID], nil
	} else {
		return nil, fmt.Errorf("there is no channel for idx:%v, msgCode:%v", idx, msgCode)
	}
}
