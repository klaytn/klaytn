// Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from accounts/abi/bind/bind.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package bind generates Go bindings for Klaytn contracts.

Source Files

Each file provides following features
 - auth.go : Provides functions to create transaction signer from a private key or a keystore wallet
 - backend.go : Defines various interfaces which work as a backend in deploying / calling a contract and filtering logs
 - base.go : Provides functions to deploy a contract, filter / watch logs, invoke a transaction and so on
 - bind.go : Provides functions to generate a wrapper around a contract ABI
 - template.go : Provides templates to build a binding to use a contract in Go and Java
 - topics.go : Provides functions for making and parsing topics
 - util.go : Provides utility functions to wait for a transaction to be mined
*/
package bind
