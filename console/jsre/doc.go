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
// This file is derived from internal/jsre/jsre.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package jsre provides an execution environment for JavaScript.

JSRE is a generic JS runtime environment embedding the otto Javascript interpreter. <https://github.com/robertkrimen/otto>
 It provides some helper functions to
 - load code from files
 - run code snippets
 - require libraries
 - bind native go objects

Because of JSRE, an user can utilize JavaScript in the console as needed.
JSRE also provides two JavaScript libraries, bignumber.js and web3.js, for users to easily access Klaytn.

Source Files

Each file provides following features
 - completion.go	: Provides functions for keyword completion
 - jsre.go	: Wraps otto JavaScript interpreter and provides an event loop
 - pretty.go	: Prints results to the standard output in more readable way
*/
package jsre
