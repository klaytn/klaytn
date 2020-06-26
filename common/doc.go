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
// This file is derived from common/bytes.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package common contains various helper functions, commonly used data types and constants.

Source Files

Each source file provides functions and types as listed below
  - big.go   		: Defines common big integers often used such as Big1
  - bytes.go		: Provides conversion functions between a byte slice and other types such as string in Hex format and int
  - cache.go		: Defines Cache interface and constants such as CacheType
  - debug.go		: Provides a function to let an user know where to report a bug. It is being used by datasync/downloader
  - format.go		: Provides a function to print a time.Duration in more readable format
  - path.go  		: Provides functions to check a file's existence and to get an absolute path
  - size.go 		: Provides StorageSize type and its stringer functions
  - types.go		: Provides commonly used Hash and Address types and its methods
  - utils.go		: Provides LoadJSON function to read a JSON file
  - variables.go	: Provides configuration values
*/
package common
