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
// This file is derived from accounts/abi/doc.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package abi implements the Klaytn ABI (Application Binary
Interface).

The Klaytn ABI is strongly typed, known at compile time
and static. This ABI will handle basic type casting; unsigned
to signed and visa versa. It does not handle slice casting such
as unsigned slice to signed slice. Bit size type casting is also
handled. ints with a bit size of 32 will be properly cast to int256,
etc.

Source Files

Each file provides the following features
 - abi.go	: Provides `ABI` struct which holds information about a contract's context and available invokable methods. It will allow you to type check function calls and packs data accordingly.
 - argument.go	: Provides `Argument` which holds the name of the argument and the corresponding type.
 - error.go	: Provides type check functions
 - event.go	: Provides `Event` struct which is an event potentially triggered by the EVM's LOG mechanism
 - method.go	: Provides `Method` struct which represents a callable given a `Name` and whether the method is a constant.
 - numbers.go	: Provides U256 function which converts a big Int into a 256 bit EVM number.
 - pack.go	: Provides functions which pack bytes slice, element and number
 - reflect.go	: Provides functions to map an ABI to a struct and detect the type and the kind of a field using reflection
 - type.go	: Provides `Type` struct which is the reflection of the supported argument type
 - unpack.go	: Provides functions which read values based on their kind
*/
package abi
