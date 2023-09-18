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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package api implements the general Klaytn API functions.

Overview of api package

This package provides various APIs to access the data of the Klaytn node.
Remote users can interact with Klyatn by calling these APIs instead of IPC.
APIs are grouped by access modifiers (public or private) and namespaces (klay, txpool, debug and personal).

Source Files

  - addrlock.go                    : implements Addrlocker which prevents another tx getting the same nonce through API.
  - api_private_account.go         : provides private APIs to access accounts managed by the node.
  - api_private_debug.go           : provides private APIs exposed over the debugging node.
  - api_public_account.go          : provides public APIs to access accounts managed by the node.
  - api_public_blockchain.go       : provides public APIs to access the Klaytn blockchain.
  - api_public_cypress.go          : provides public APIs to return specific information of Klaytn Cypress network.
  - api_public_debug.go            : provides public APIs exposed over the debugging node.
  - api_public_klay.go             : provides public APIs to access Klaytn related data.
  - api_public_net.go              : provides public APIs to offer network related RPC methods.
  - api_public_transaction_pool.go : provides public APIs having "klay" namespace to access transaction pool data.
  - api_public_tx_pool.go          : provides public APIs having "txpool" namespace to access transaction pool data.
  - backend.go                     : provides the common API services.
  - tx_args.go                     : provides API argument structures and functions.
*/
package api
