// Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/state/statdb.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package state provides an uppermost caching layer of the Klaytn state trie.
This package is used to read/write stateObject from/to StateDB and it also acts as a bridge between
the objects and the persistent layer.

Overview of state package

stateObject represents a Klaytn account identified by its address.
Once it is loaded from the persistent layer, it is cached and managed by StateDB.

StateDB caches stateObjects and mediates the operations to them.


Source Files

Related functions and variables are defined in the files listed below
  - database.go              : Defines Database and other interfaces used in the package
  - dump.go                  : Functions to dump the contents of StateDB both in raw format and indented format
  - journal.go               : journal and state changes to track the list of state modifications since the last state commit
  - state_object.go          : Implementation of stateObject
  - state_object_encoder.go  : stateObjectEncoder is used to encode stateObject in parallel manner
  - statedb.go               : Implementation of StateDB
  - sync.go                  : Functions to schedule a state trie download
*/
package state
