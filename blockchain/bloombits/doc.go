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
// This file is derived from core/bloombits/doc.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package bloombits implements bloom filtering on batches of data.

Source Files

Each file provides the following features.
  - generator.go : generates the rotated bloom bits to be used for batched filtering.
  - matcher.go : creates a stream of potential blocks to inspect for data content through binary AND/OR operations on the bit-streams.
  - scheduler.go : handles the scheduling of bloom-filter retrieval operations.
*/
package bloombits
