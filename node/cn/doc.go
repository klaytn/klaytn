// Copyright 2018 The klaytn Authors
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

/*
Package cn implements components related to network management and message handling.
CN implements the Klaytn consensus node service.
ProtocolManager handles the messages from the peer nodes and manages its peers.
Peer is the interface used for peer nodes and has two different kinds of implementation
depending on single or multi channel usage.

Source Files

  - api.go              : provides private debug API related to block and state
  - api_backend.go      : implements CNAPIBackend which is a wrapper of CN to serve API requests
  - api_tracer.go       : provides private debug API related to trace chain, block and state
  - backend.go          : implements CN struct used for the Klaytn consensus node service
  - bloombits.go        : implements BloomIndexer, an indexer built with bloom bits for fast filtering
  - channel_manager.go  : implements ChannelManager struct, which is used to manage channel for each message
  - config.go           : defines the configuration used by CN struct
  - gen_config.go       : is automatically generated from config.go
  - handler.go          : implements ProtocolManager which handles the message and manages network peers
  - metrics.go          : includes statistics used in cn package
  - peer.go             : provides the interface and implementation of Peer interface
  - peer_set.go         : provides the interface and implementation of PeerSet interface
  - protocol.go         : defines the protocol version of Klaytn network and includes errors in cn package
  - sync.go             : includes syncing features of ProtocolManager
*/
package cn
