// Modifications Copyright 2019 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from eth/protocol.go (2018/06/04).
// Modified and improved for the klaytn development.

package sc

import (
	"math/big"

	"github.com/klaytn/klaytn/common"
)

const ProtocolMaxMsgSize = 10 * 1024 * 1024 // Maximum cap on the size of a protocol message

const (
	// Protocol messages belonging to servicechain/1
	StatusMsg = 0x00

	// Below message can be deprecated.
	ServiceChainTxsMsg                     = 0x01
	ServiceChainReceiptResponseMsg         = 0x02
	ServiceChainReceiptRequestMsg          = 0x03
	ServiceChainParentChainInfoResponseMsg = 0x04
	ServiceChainParentChainInfoRequestMsg  = 0x05

	ServiceChainCall     = 0x06
	ServiceChainResponse = 0x07
	ServiceChainNotify   = 0x08
)

var (
	SCProtocolName    = "servicechain"
	SCProtocolVersion = []uint{1}
	SCProtocolLength  = []uint64{9}
)

// Protocol defines the protocol of the consensus
type SCProtocol struct {
	// Official short name of the protocol used during capability negotiation.
	Name string
	// Supported versions of the Klaytn protocol (first is primary).
	Versions []uint
	// Number of implemented message corresponding to different protocol versions.
	Lengths []uint64
}

type errCode int

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrNoStatusMsg
	ErrUnexpectedTxType
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

// XXX change once legacy code is out
var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIdMismatch:       "NetworkId mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrUnexpectedTxType:        "Unexpected tx type",
}

// statusData is the network packet for the status message.
type statusData struct {
	ProtocolVersion uint32
	NetworkId       uint64
	TD              *big.Int
	CurrentBlock    common.Hash
	ChainID         *big.Int // A child chain must know parent chain's ChainID to sign a transaction.
}
