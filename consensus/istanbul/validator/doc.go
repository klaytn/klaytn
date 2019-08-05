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
Package validator implements the types related to the validators participating in consensus.

`Validator` and `ValidatorSet` interfaces are used to validate blocks and make a consensus.
Implementations of these interfaces are different, depending on proposer-choosing policy.

Validator

`Validator` is a node which has 2 features to make a consensus: proposing and validating.

Propose: A node can propose a block if it is a proposer. Only validators can be a proposer.

Validate: A validator node can validate blocks from proposers. A block is valid only if more than 2/3 of validators approve the given block.

ValidatorSet

`ValidatorSet` is a group of validators. It is also called as a council.
ValidatorSet calculates the block proposer of an upcoming block.
A validator selected as a block proposer will have a chance to make a block.

Implementation in Klaytn

Klaytn implements `Validator` and `ValidatorSet` interface for Klaytn consensus.
Klaytn reflects the ratio of staking amounts to the probability of selecting a proposer.
This is called weightedRandom policy.
Detailed information can be found in https://docs.klaytn.com/klaytn/token_economy#klaytn-governance-council-reward.
Implementation structures are weightedValidator and weightedCouncil in weighted.go file.

Files

- default.go   : Validator and ValidatorSet for roundRobin policy is implemented.

- weighted.go  : Validator and ValidatorSet for weightedRandom policy is implemented.

- validator.go : common functions for Validator and ValidatorSet are implemented.
*/
package validator
