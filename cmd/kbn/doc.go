// Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/bootnode/main.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
kbn runs a bootstrap node for the Klaytn Node Discovery Protocol.

A bootstrap node is a kind of registry service which has every nodes' information and delivers it to a querying node to help the node to join the network.

Source Files

Each file contains the following contents
 - api.go	: Provides various APIs to use bootnode services
 - backend.go	: Provides supporting functions for APIs
 - config.go	: Provides `bootnodeConfig` which contains a configuration and accompanying setter and parser functions
 - main.go	: Main entry point of the application
 - node.go	: Provides `Node` struct which defines what kind of APIs can be provided through which port and protocols

*/
package main
