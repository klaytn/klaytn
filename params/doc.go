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
Package params contains configuration parameters for Klaytn.

Source Files

Each file contains following parameters.

	- bootnodes.go			: Provides boot nodes information for Cypress and Baobab
	- computation_cost_params.go	: Defines computation costs for each opcode
	- config.go			: Defines various structs for different settings of a network. Also provides getters for those settings
	- denomination.go		: Defines units of KLAY
	- gas_table.go			: Organizes gas prices for different Klaytn phases. Currently prices for Cypress is defined
	- governance_params.go		: Defines constants for governance and reward system. Also provides setters and getters for reward releated variables
	- network_params.go		: Defines network parameters that need to be constant between clients. Only `BloomBitsBlocks` is defined at the moment
	- protocol_params.go		: Defines fee schedule, total time limit and maximum computation cost
	- version.go			: Defines release and version number. Also provides a getter for the version
*/
package params
