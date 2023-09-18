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
// This file is derived from internal/web3ext/web3ext.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package web3ext contains the Klaytn specific web3.js extensions.

web3ext defines `Modules` which defines APIs for each category. This `Modules` is used by console to provide APIs to users.

API Categories

APIs are categorized as follows. If you want to know more detail, please refer to https://docs.klaytn.com/bapp/json-rpc
  - admin
  - debug
  - klay
  - miner
  - net
  - personal
  - rpc
  - txpool
  - istanbul
  - mainbridge
  - subbridge
  - clique
  - governance
  - bootnode
  - chaindatafetcher
  - eth
*/
package web3ext
