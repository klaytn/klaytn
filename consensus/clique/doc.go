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
Package clique implements PoA (Proof of Authority) consensus engine which is mainly for private chains.

Consensus Engine Overview

In clique, only pre-selected nodes can propose a block. The list of these nodes are described in the header's extra field and the list can be changed by calling Propose() API.

The sequence of proposers is decided by the block number and the number of nodes which can propose a block. If a proposer proposes a block in its turn (in-turn), it will have a blockscore(formerly known as difficulty) of 2.
If a proposer proposes a block in other node's turn (out-of-turn), it will have a blockscore of 1.

If an in-turn proposer didn't propose a block for some time, other nodes start to propose their block with a blockscore of 1. So the block generation can continue.
But there can be a fork because several nodes can propose a block at the same time.
In this case, a block reorganization will happen based on total blockscore (the sum of all blockscore from the genesis block to the current one).
*/
package clique
