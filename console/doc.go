// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

/*
Package console implements JavaScript console.

Console is a JavaScript interpreted runtime environment. It is a fully fledged
JavaScript console attached to a running node via an external or in-process RPC
client.

Source Files

Each file provides following features
 - bridge.go	: bridge is a collection of JavaScript utility methods to bridge the .js runtime environment and the Go RPC connection backing the remote method calls
 - console.go	: Implements a console which supports JavaScript runtime environment
 - prompter.go	: Provides UserPrompter which defines the methods needed by the console to prompt the user for various types of inputs, such as normal text, a password and a confirmation
*/
package console
