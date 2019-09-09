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
// This file is derived from crypto/secp256k1/secp256.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package secp256k1 wraps the bitcoin secp256k1 C library.

secp256k1 refers to the parameters of the elliptic curve used in Bitcoin's public-key cryptography and is defined in Standards for Efficient Cryptography (SEC)(Certicom Research, http://www.secg.org/sec2-v2.pdf).

Package secp256k1 provides wrapper functions to utilize the library functions in Go.

Source Files

Each source file has the following contents
 - secp256.go  : Provides wrapper functions to utilize the secp256k1 library written in C
 - curve.go    : Implements Koblitz elliptic curves
 - panic_cb.go : Provides callbacks for converting libsecp256k1 internal faults into recoverable Go panics
 - schnorr.go  : Implements Schnorr signature algorithm. It is planned to be used in Klaytn
*/
package secp256k1
