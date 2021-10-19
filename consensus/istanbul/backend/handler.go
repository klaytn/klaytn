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

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/networks/p2p"
)

const (
	IstanbulMsg = 0x11
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed       = errors.New("fail to decode istanbul message")
	errNoChainReader      = errors.New("sb.chain is nil! --mine option might be missing")
	errInvalidPeerAddress = errors.New("invalid address")

	// TODO-Klaytn-Istanbul: define Versions and Lengths with correct values.
	istanbulProtocol = consensus.Protocol{
		Name:     "istanbul",
		Versions: []uint{64},
		Lengths:  []uint64{21},
	}
)

// Protocol implements consensus.Engine.Protocol
func (sb *backend) Protocol() consensus.Protocol {
	return istanbulProtocol
}

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if msg.Code == IstanbulMsg {
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
		var m *lru.ARCCache
		ms, ok := sb.recentMessages.Get(addr)
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
		return errNoChainReader
	}
	validators := sb.getValidators(sb.chain.CurrentHeader().Number.Uint64(), sb.chain.CurrentHeader().Hash())
	for _, val := range validators.List() {
		if addr == val.Address() {
			return nil
		}
	}
	for _, val := range validators.DemotedList() {
		if addr == val.Address() {
			return nil
		}
	}
	return errInvalidPeerAddress
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *backend) SetBroadcaster(broadcaster consensus.Broadcaster, nodetype common.ConnType) {
	sb.broadcaster = broadcaster
	if nodetype == common.CONSENSUSNODE {
		sb.broadcaster.RegisterValidator(common.CONSENSUSNODE, sb)
	}
}

// RegisterConsensusMsgCode registers the channel of consensus msg.
func (sb *backend) RegisterConsensusMsgCode(peer consensus.Peer) {
	if err := peer.RegisterConsensusMsgCode(IstanbulMsg); err != nil {
		logger.Error("RegisterConsensusMsgCode failed", "err", err)
	}
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
