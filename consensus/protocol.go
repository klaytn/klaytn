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
// This file is derived from quorum/consensus/protocol.go (2018/06/04).
// Modified and improved for the klaytn development.

package consensus

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
)

// Constants to match up protocol versions and messages
const (
	Klay62 = 62
	Klay63 = 63
)

var (
	KlayProtocol = Protocol{
		Name:     "klay",
		Versions: []uint{Klay63, Klay62},
		Lengths:  []uint64{17, 8},
	}
)

// Protocol defines the protocol of the consensus
type Protocol struct {
	// Official short name of the protocol used during capability negotiation.
	Name string
	// Supported versions of the klaytn protocol (first is primary).
	Versions []uint
	// Number of implemented message corresponding to different protocol versions.
	Lengths []uint64
}

// istanbul BFT
// Broadcaster defines the interface to enqueue blocks to fetcher and find peer
type Broadcaster interface {
	// Enqueue add a block into fetcher queue
	Enqueue(id string, block *types.Block)
	// FindPeers retrives peers by addresses
	FindPeers(map[common.Address]bool) map[common.Address]Peer

	FindCNPeers(map[common.Address]bool) map[common.Address]Peer

	GetCNPeers() map[common.Address]Peer

	GetENPeers() map[common.Address]Peer

	RegisterValidator(conType common.ConnType, validator p2p.PeerTypeValidator)
}

// Peer defines the interface to communicate with peer
type Peer interface {
	// Send sends the message to this peer
	Send(msgcode uint64, data interface{}) error

	// RegisterConsensusMsgCode registers the channel of consensus msg.
	RegisterConsensusMsgCode(msgCode uint64) error
}
