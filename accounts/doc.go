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
// This file is derived from accounts/accounts.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package accounts implements high level Klaytn account management.

There are several important data structures and their hierarchy is as below
 Manager -> Backend -> Wallet -> Account

Source Files

Each file provides the following features
 - accounts.go	: Provides `Account` struct and Wallet/Backend interfaces. `Account` represents a Klaytn account located at a specific location defined by the optional URL field. `Wallet` represents a software or hardware wallet and `Backend` is a wallet provider that may contain a batch of accounts they can sign transactions with
 - errors.go	: Provides various account related error variables and helper functions
 - hd.go		: Defines derivation paths for Klaytn and parser function to derive the path from a path string. Klaytn uses 8217 as its coin type
 - manager.go 	: Provides `Manager` which is an overarching account manager that can communicate with various backends for signing transactions
 - url.go 	: Provides `URL` struct which represents the canonical identification URL of a wallet or account
*/
package accounts
