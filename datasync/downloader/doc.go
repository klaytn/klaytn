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
// This file is derived from eth/downloader/downloader.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package downloader contains the manual full chain synchronisation.

How downloader works

The downloader is responsible for synchronizing up-to-date status from the peers connected to it. To do this,
download "headers", "bodies", and "receipts" in parallel, merge them through the pipeline, and reflect them in the state trie.

Source Files

Downloader related functions and variables are defined in the files listed below
  - api.go              : Console api to get synchronization information.
  - downloader.go       : Functions and variables to sync peer and block. And modules for QOS.
  - downloader_test.go  : Functions for testing the downloader package.
  - events.go           : Define event type.
  - metrics.go          : Metric variables for packet transmissions and receptions..
  - modes.go            : Defines the type for SyncMode. SyncMode includes "FullSync", "FastSync", and "LightSync".
  - peer.go             : Functions that request a packet to the peer and check and set the network status of the peer.
  - queue.go            : Functions for managing and scheduling received headers, bodies, and receipts.
  - types.go            : Defines the type of downloaded packet
*/
package downloader
