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
Package validator implements the types related to the validators participating in istanbul consensus.

Types in validator package implement `Validator` and `ValidatorSet` interface in istanbul/validator.go file.
Those types are used for validating blocks to make consensus.

Validator

`Validator` is a node which has 2 features for consensus, proposing and validating. Klaytn uses weightedValidator for Klaytn consensus.

Propose: A node can propose a block, if it is a proposer. Only validator nodes can be a proposer.

Validate: A validator node can validate blocks from proposers. A block is valid only if more than 2/3 of validators approve the given block.

ValidatorSet

`ValidatorSet` is a group of validators. It is also called as a council.
ValidatorSet calculates the block proposer of a upcoming block.
A validator selected as a block proposer will have a chance to make a block.
Only a validator in ValidatorSet can propose a block.
A block made by a validator which is not in the ValidatorSet will not be accepted by other validators.
Klaytn uses weightedCouncil for Klaytn consensus.

Files

- default.go   : Validator and ValidatorSet for roundRobin policy is implemented.

- weighted.go  : Validator and ValidatorSet for weightedRandom policy is implemented.

- validator.go : common functions for Validator and ValidatorSet are implemented.
*/
package validator
