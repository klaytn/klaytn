// Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum library.
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
// This file is derived from eth/filters/filter_system.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package filters implements a Klaytn filtering system for blocks, transactions and log events.

Source Files

  - api.go           : provides public filter API functions to generate filters and use them to filter the result
  - filter.go        : implements basic filtering system based on bloom filter
  - filter_system.go : provides subscription scheme to register and filter the specific events
*/
package filters
