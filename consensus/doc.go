// Copyright 2018 The klaytn Authors
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
// This file is derived from consensus/consensus.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package consensus defines interfaces for consensus engines and ChainReader.

Klaytn currently uses istanbul BFT engine on the mainnet, but PoA based clique engine also can be used.
Traditional PoW engine(gxhash) is used by legacy test codes but it is deprecated and is not recommended to use.

By implementing the Engine interface, new consensus engine can be added and used in Klaytn.

ChainReader interface defines a small collection of methods needed to access the local blockchain during a block header verification.
*/
package consensus
