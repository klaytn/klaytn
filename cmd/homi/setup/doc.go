// Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// Copyright 2017 AMIS Technologies
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

/*
setup package generates config files following the given deployment options.
It creates the given number of genesis.json and nodekeys

Source Files

Each file contains the following contents
 - cmd.go : Provides functions to generate config files with given deployment configuration
 - flags.go : Defines command line flags which can be used in `setup` command
 - klaytn_config.go : Defines `KlaytnConfig` and provides a template to build it
 - prometheus_config.go : Defines `PrometheusConfig` and provides a template to build it
*/
package setup
