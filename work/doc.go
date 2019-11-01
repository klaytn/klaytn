// Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from miner/miner.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package work implements Klaytn block creation and mining.

work package is in charge of generating blocks and it contains `miner`, `worker` and `agents` which performs building blocks with transactions.
`miner` contains `worker` and a consensus engine, so it controls the worker and passes the result to the consensus engine.
`worker` contains `agent`s and each `agent` performs mining activity.

Source Files
 - agent.go		: Provides CpuAgent and accompanying functions which works as an agent of a miner. Agent is in charge of creating a block
 - remote_agent.go	: Provides RemoteAgent working as an another agent for a miner and can be controlled by RPC calls
 - work.go		: Provides Miner struct and interfaces through which the miner communicate with other objects
 - worker.go		: Provides Worker and performs the main part of block creation
*/
package work
