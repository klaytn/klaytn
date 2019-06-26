// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/backend/handler.go (2018/06/04).
// Modified and improved for the klaytn development.

package backend

import (
	"errors"
	"github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/node"
)

const (
	istanbulMsg = 0x11
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode istanbul message")
)

// Protocol implements consensus.Engine.Protocol
func (sb *backend) Protocol() consensus.Protocol {
	return consensus.Protocol{
		Name:     "istanbul",
		Versions: []uint{64},
		//Lengths:  []uint64{18},
		//Lengths:  []uint64{19},  // add PoRMsg
		Lengths: []uint64{21},
	}
}

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if msg.Code == istanbulMsg {
		if !sb.coreStarted {
			return true, istanbul.ErrStoppedEngine
		}

		var cmsg istanbul.ConsensusMsg

		//var data []byte
		if err := msg.Decode(&cmsg); err != nil {
			return true, errDecodeFailed
		}
		data := cmsg.Payload
		hash := istanbul.RLPHash(data)

		// Mark peer's message
		ms, ok := sb.recentMessages.Get(addr)
		var m *lru.ARCCache
		if ok {
			m, _ = ms.(*lru.ARCCache)
		} else {
			m, _ = lru.NewARC(inmemoryMessages)
			sb.recentMessages.Add(addr, m)
		}
		m.Add(hash, true)

		// Mark self known message
		if _, ok := sb.knownMessages.Get(hash); ok {
			return true, nil
		}
		sb.knownMessages.Add(hash, true)

		go sb.istanbulEventMux.Post(istanbul.MessageEvent{
			Payload: data,
			Hash:    cmsg.PrevHash,
		})

		return true, nil
	}
	return false, nil
}

func (sb *backend) ValidatePeerType(addr common.Address) error {
	// istanbul.Start vs try to connect by peer
	for sb.chain == nil {
		return errors.New("sb.chain is nil! --mine option might be missing")
	}
	for _, val := range sb.getValidators(sb.chain.CurrentHeader().Number.Uint64(), sb.chain.CurrentHeader().Hash()).List() {
		if addr == val.Address() {
			return nil
		}
	}
	return errors.New("invalid address")
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *backend) SetBroadcaster(broadcaster consensus.Broadcaster, nodetype p2p.ConnType) {
	sb.broadcaster = broadcaster
	if nodetype == node.CONSENSUSNODE {
		sb.broadcaster.RegisterValidator(node.CONSENSUSNODE, sb)
	}
}

// RegisterConsensusMsgCode registers the channel of consensus msg.
func (sb *backend) RegisterConsensusMsgCode(peer consensus.Peer) {
	peer.RegisterConsensusMsgCode(istanbulMsg)
}

func (sb *backend) NewChainHead() error {
	sb.coreMu.RLock()
	defer sb.coreMu.RUnlock()
	if !sb.coreStarted {
		return istanbul.ErrStoppedEngine
	}

	go sb.istanbulEventMux.Post(istanbul.FinalCommittedEvent{})
	return nil
}
