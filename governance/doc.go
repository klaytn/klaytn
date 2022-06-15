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
Package governance contains functions and variables used for voting and reflecting vote results in Klaytn.
In Klaytn, various settings such as the amount of KLAY minted as a block reward can be changed by using governance vote.
These votes can be casted by nodes of Governance Council and detailed introduction can be found at https://docs.klaytn.com/klaytn/design/governance

How to cast a vote

To cast a vote, a node have to be a member of the Governance Council.
If the governance mode is "single", only one designated node (the governing node) can vote.
In the console of the node, "governance.vote(key, value)" API can be used to cast a vote.

Keys for the voting API

Following keys can be handled as of 7/20/2019.
  - "governance.governancemode"   : To change the governance mode
  - "governance.governingnode"    : To change the governing node if the governance mode is "single"
  - "governance.unitprice"        : To change the unitprice of Klaytn (Unit price is same as gasprice in Ethereum)
  - "governance.addvalidator"     : To add new node as a council node
  - "governance.removevalidator"  : To remove a node from the governance council
  - "istanbul.epoch"              : To change Epoch, the period to gather votes
  - "istanbul.committeesize"      : To change the size of the committee
  - "reward.mintingamount"        : To change the amount of block generation reward
  - "reward.ratio"                : To change the ratio used to distribute the reward between block proposer node, PoC and KIR
  - "reward.useginicoeff"         : To change the application of gini coefficient to reduce gap between CCOs
  - "reward.deferredtxfee"        : To change the way of distributing tx fee
  - "reward.minimumstake"         : To change the minimum amount of stake to participate in the governance council


How governance works

Governance package contains a governance struct which stores current system configurations and voting status.
If a vote passed, the governance struct is updated to provide new information to related packages and users.
The API documentation can be found at https://docs.klaytn.com/bapp/json-rpc/api-references/governance

When a CN (consensus node which is managed by CCO) proposes a block, it writes its vote on the block header and other nodes
parse the header and handle it. This process is handled by snapshot.go in the consensus engine and processed by functions in handler.go

If a vote satisfies the requirement (more than 50% of votes in favor of), it will update the governance struct and many other packages
like "reward", "txpool" and so on will reference it.


Source Files

Governance related functions and variables are defined in the files listed below
  - default.go    : the governance struct, cache and persistence
  - handler.go    : functions to handle votes and its application
  - api.go        : console APIs to get governance information and to cast a vote
  - interface.go  : Abstract interfaces to various underlying implementations
  - mixed.go      : Wrapper for multiple engine implementations

*/
package governance
