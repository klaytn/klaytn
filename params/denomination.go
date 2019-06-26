// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from params/denomination.go (2018/06/04).
// Modified and improved for the klaytn development.

package params

// These are the multipliers for KLAY denominations.
// Example: To get the peb value of an amount in 'ston', use
//
//    new(big.Int).Mul(value, big.NewInt(params.Ston))
//
const (
	Peb      = 1    // official notation 'peb'
	Kpeb     = 1e3  // official notation 'kpeb'
	Mpeb     = 1e6  // same
	Gpeb     = 1e9  // same
	Ston     = 1e9  // official notation 'ston'
	UKLAY    = 1e12 // official notation 'uKLAY'
	MiliKLAY = 1e15 // official notation 'mKLAY'
	KLAY     = 1e18 // same
	KKLAY    = 1e21 // official notation 'kKLAY'
	MKLAY    = 1e24 // same
	GKLAY    = 1e27 // same
)
