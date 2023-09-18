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
Package nodecmd contains command definitions and related functions used for node cmds, such as kcn, kpn, and ken.

Source Files

Each file contains following contents
 - accountcmd.go		: Provides functions for creating, updating and importing an account.
 - chaincmd.go		: Provides functions to `init` a block chain,
 - consolecmd.go		: Provides console functions `attach` and `console`
 - migrationcmd.go		: Provides functions of DB migration
 - defaultcmd.go		: Provides functions to start a node
 - dumpconfigcmd.go		: Provides functions to dump and print current config to stdout
 - nodeflags.go		: Defines various flags that configure the node
 - versioncmd.go		: Provides functions to print application's version
*/
package nodecmd
