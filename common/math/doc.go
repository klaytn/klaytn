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
// This file is derived from common/math/big.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package math provides convenience functions to use big.Int and to parse string into an integer

`big.go` provides Max, Min, Pow, Parse functions for big.Int types and it also implements encoding.TextMarshaler and encoding.TextUnmarshaler.

`integer.go` provides functions to parse string into unsigned int and to calculate safely by detecting overflow. It also provides implementation of encoding.TextMarshaler and encoding.TextUnmarshaler
*/
package math
