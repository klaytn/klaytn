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
// This file is derived from ethstats/ethstats.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package statedb implements the Merkle Patricia Trie structure used for state object trie.
This package is used to read/write data from/to the state object trie.

Overview of statedb package

There are 3 key struct in this package: Trie, SecureTrie and Database.

Trie struct represents a Merkle Patricia Trie.

SecureTrie is a basically same as Trie but it wraps a trie with key hashing.
In a SecureTrie, all access operations hash the key using keccak256.
This prevents calling code from creating long chains of nodes that increase the access time.

Database is an intermediate write layer between the trie data structures and
the disk database. The aim is to accumulate trie writes in-memory and only
periodically flush a couple tries to disk, garbage collecting the remainder.


Source Files

Related functions and variables are defined in the files listed below
  - database.go     : Implementation of Database struct
  - db_migration.go : Implementation of DB migration
  - derive_sha.go   : Implementation of DeriveShaOrig used in Klaytn
  - encoding.go     : Implementation of 3 encodings: KEYBYTES, HEX and COMPACT
  - errors.go       : Errors used in this package
  - hasher.go       : Implementation of recursive and bottom-up hashing
  - iterator.go     : Implementation of key-value trie iterator that traverses a Trie
  - node.go         : Implementation of 4 types of nodes, used in Merkle Patricia Trie
  - proof.go        : Functions which construct a Merkle Patricia Proof for the given key
  - secure_trie.go  : Implementation of Merkle Patricia Trie with key hashing
  - sync.go         : Implementation of state trie sync
  - trie.go         : Implementation of Merkle Patricia Trie
*/
package statedb
