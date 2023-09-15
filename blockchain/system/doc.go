// Copyright 2023 The klaytn Authors
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

package system

/*
Package system deals with system contracts in Klaytn.

The system contracts are smart contracts that controls the blockchain consensus.

- AddressBook: the list of consensus nodes. It stores their nodeId, staking contracts and reward
  addresses.
- GovParam: the governance parameter storage, meant to be modified by on-chain governance vote.
  It overrides the existing header-based governance (i.e. block header votes) and it can be
  optionally configured after the Kore hardfork.
- TreasuryRebalance: the KFF and KCF rebalance configuration, can be configured by the optional
  KIP103 hardfork.

*/
