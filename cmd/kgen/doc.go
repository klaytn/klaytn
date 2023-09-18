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
kgen can be used to generate a nodekey and a derived address and kni information.

By providing some options in the command line, a user can assign IP and port for kni.

Options

All available options are as follows.
   --file        Generate a nodekey and a Klaytn node information as files
   --ip value    Specify an IP address (default: "0.0.0.0")
   --port value  Specify a tcp port number (default: 32323)
   --help, -h    Show help
*/
package main
