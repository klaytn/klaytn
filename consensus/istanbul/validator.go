// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/validator.go (2018/06/04).
// Modified and improved for the klaytn development.

package istanbul

import (
	"strings"

	"github.com/klaytn/klaytn/params"

	"github.com/klaytn/klaytn/common"
)

type Validator interface {
	// Address returns address
	Address() common.Address

	// String representation of Validator
	String() string

	RewardAddress() common.Address
	VotingPower() uint64
	Weight() uint64
}

// ----------------------------------------------------------------------------

type Validators []Validator

func (slice Validators) Len() int {
	return len(slice)
}

func (slice Validators) Less(i, j int) bool {
	return strings.Compare(slice[i].String(), slice[j].String()) < 0
}

func (slice Validators) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice Validators) AddressStringList() []string {
	var stringAddrs []string
	for _, val := range slice {
		stringAddrs = append(stringAddrs, val.Address().String())
	}
	return stringAddrs
}

// ----------------------------------------------------------------------------

type ValidatorSet interface {
	// Calculate the proposer
	CalcProposer(lastProposer common.Address, round uint64)
	// Return the validator size
	Size() uint64
	// Return the sub validator group size
	SubGroupSize() uint64
	// Set the sub validator group size
	SetSubGroupSize(size uint64)
	// Return the validator array
	List() []Validator
	// Return the demoted validator array
	DemotedList() []Validator
	// SubList composes a committee after setting a proposer with a default value.
	SubList(prevHash common.Hash, view *View) []Validator
	// Return whether the given address is one of sub-list
	CheckInSubList(prevHash common.Hash, view *View, addr common.Address) bool
	// SubListWithProposer composes a committee with given parameters.
	SubListWithProposer(prevHash common.Hash, proposer common.Address, view *View) []Validator
	// Get validator by index
	GetByIndex(i uint64) Validator
	// Get validator by given address
	GetByAddress(addr common.Address) (int, Validator)
	// Get demoted validator by given address
	GetDemotedByAddress(addr common.Address) (int, Validator)
	// Get current proposer
	GetProposer() Validator
	// Check whether the validator with given address is a proposer
	IsProposer(address common.Address) bool
	// Add validator
	AddValidator(address common.Address) bool
	// Remove validator
	RemoveValidator(address common.Address) bool
	// Copy validator set
	Copy() ValidatorSet
	// Get the maximum number of faulty nodes
	F() int
	// Get proposer policy
	Policy() ProposerPolicy

	IsSubSet() bool

	// Refreshes a list of candidate proposers with given hash and blockNum
	Refresh(hash common.Hash, blockNum uint64, config *params.ChainConfig, isSingle bool, governingNode common.Address, minStaking uint64) error

	SetBlockNum(blockNum uint64)

	Proposers() []Validator // TODO-Klaytn-Issue1166 For debugging

	TotalVotingPower() uint64

	Selector(valSet ValidatorSet, lastProposer common.Address, round uint64) Validator
}

// ----------------------------------------------------------------------------

type ProposalSelector func(ValidatorSet, common.Address, uint64) Validator
