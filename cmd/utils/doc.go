// Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/utils/cmd.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package utils contains internal helper functions for klaytn commands.

utils package provides various helper functions especially for handling various commands and options.

Source Files

Each file contains the following contents
 - app.go          : Provides NewCLI() function but it is not being used.
 - cmd.go          : Provide import/export chain functions but it is not being used.
 - customflags.go  : Provides `DirectoryString`, `DirectoryFlags` and marshaling functions to support custom flags
 - files.go        : Provides `WriteFile` function to store contents in a given file
 - flags.go        : Defines various flags which can be used in running a node
 - flaggroup.go    : Categorizes flags into groups to print structured help descriptions.
 - helptemplate.go : Provides a template for help contents which explains option names and its usages
 - strings.go      : Provides helper functions to handle string manipulations
 - testcmd.go      : Provides test functions to test command arguments
 - usage.go        : Provides help printer which prints help contents neatly
*/
package utils
