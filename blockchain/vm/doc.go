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
// This file is derived from core/vm/doc.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package vm implements the Ethereum Virtual Machine.

The vm package implements one EVM, a byte code VM.
The Byte Code VM loops over a set of bytes and executes them according to the set of rules defined in the Ethereum yellow paper.

As well as the original functionality of the EVM, this package implemented additional pre-compiled contracts to support the native features of Klaytn.
For more information about pre-compiled contracts, see KlaytnDocs (https://docs.klaytn.com/smart-contract/precompiled-contracts).
*/
package vm
