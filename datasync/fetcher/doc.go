// Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from eth/fetcher/fetcher.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package fetcher contains the block announcement based synchronisation.

How fetcher works

If the node receives a whole block, fetcher inserts the block into the chain and broadcast it to its peers.
If a block hash is received instead of a block, node requests and reflects the header and body from the peer who sent the block hash.

Source Files

Functions and variables related to fetcher are defined in the files listed below.
  - fetcher.go      : It includes functions for fetching the received block, header, body, and a queue data structure for the fetch operation.
  - fetcher_test.go : Functions for testing the fetcher's functions.
  - metrics.go      : Metric variables for packet header, body and blocks.
*/
package fetcher
