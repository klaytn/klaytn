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
// This file is derived from quorum/consensus/istanbul/validator/default.go (2018/06/04).
// Modified and improved for the klaytn development.

package validator

import (
	"math"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/params"
)

const (
	defaultSubSetLength = 21
)

type defaultValidator struct {
	address common.Address
}

func (val *defaultValidator) Address() common.Address {
	return val.address
}

func (val *defaultValidator) String() string {
	return val.Address().String()
}

func (val *defaultValidator) Equal(val2 *defaultValidator) bool {
	return val.address == val2.address
}

func (val *defaultValidator) Hash() int64 {
	return val.address.Hash().Big().Int64()
}

func (val *defaultValidator) RewardAddress() common.Address { return common.Address{} }
func (val *defaultValidator) VotingPower() uint64           { return 1000 }
func (val *defaultValidator) Weight() uint64                { return 0 }

type defaultSet struct {
	subSize uint64

	validators istanbul.Validators
	policy     istanbul.ProposerPolicy

	proposer    atomic.Value
	validatorMu sync.RWMutex
	selector    istanbul.ProposalSelector
}

func newDefaultSet(addrs []common.Address, policy istanbul.ProposerPolicy) *defaultSet {
	valSet := &defaultSet{}

	valSet.subSize = defaultSubSetLength
	valSet.policy = policy
	// init validators
	valSet.validators = make([]istanbul.Validator, len(addrs))
	for i, addr := range addrs {
		valSet.validators[i] = New(addr)
	}
	// sort validator
	sort.Sort(valSet.validators)
	// init proposer
	if valSet.Size() > 0 {
		valSet.proposer.Store(valSet.GetByIndex(0))
	}
	valSet.selector = roundRobinProposer
	if policy == istanbul.Sticky {
		valSet.selector = stickyProposer
	}

	return valSet
}

func newDefaultSubSet(addrs []common.Address, policy istanbul.ProposerPolicy, subSize uint64) *defaultSet {
	valSet := &defaultSet{}

	valSet.subSize = subSize
	valSet.policy = policy
	// init validators
	valSet.validators = make([]istanbul.Validator, len(addrs))
	for i, addr := range addrs {
		valSet.validators[i] = New(addr)
	}
	// sort validator
	sort.Sort(valSet.validators)
	// init proposer
	if valSet.Size() > 0 {
		valSet.proposer.Store(valSet.GetByIndex(0))
	}
	valSet.selector = roundRobinProposer
	if policy == istanbul.Sticky {
		valSet.selector = stickyProposer
	}

	return valSet
}

func (valSet *defaultSet) Size() uint64 {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return uint64(len(valSet.validators))
}

func (valSet *defaultSet) SubGroupSize() uint64 {
	return valSet.subSize
}

// SetSubGroupSize sets committee size of the valSet.
func (valSet *defaultSet) SetSubGroupSize(size uint64) {
	if size == 0 {
		logger.Error("cannot assign committee size to 0")
		return
	}
	valSet.subSize = size
}

func (valSet *defaultSet) List() []istanbul.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return valSet.validators
}

func (valSet *defaultSet) DemotedList() []istanbul.Validator {
	return nil
}

// SubList composes a committee after setting a proposer with a default value.
// This functions returns whole validators if it failed to compose a committee.
func (valSet *defaultSet) SubList(prevHash common.Hash, view *istanbul.View) []istanbul.Validator {
	// TODO-Klaytn-Istanbul: investigate whether `valSet.GetProposer().Address()` is a proper value or the proposer should be calculated based on `view`
	proposer := valSet.GetProposer()
	if proposer == nil {
		return valSet.List()
	}
	return valSet.SubListWithProposer(prevHash, proposer.Address(), view)
}

// SubListWithProposer composes a committee with given parameters.
// The first member of the committee is set to the given proposer without calculating proposer with the given `view`.
// The second member of the committee is calculated with a round number of the given view and `valSet.blockNum`.
// The reset of the committee is selected with a random seed derived from `prevHash`.
// This functions returns whole validators if it failed to compose a committee.
func (valSet *defaultSet) SubListWithProposer(prevHash common.Hash, proposerAddr common.Address, view *istanbul.View) []istanbul.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	validators := valSet.validators
	validatorSize := uint64(len(validators))
	committeeSize := valSet.subSize

	// return early if the committee size is equal or larger than the validator size
	if committeeSize >= validatorSize {
		return validators
	}

	// find the proposer
	proposerIdx, proposer := valSet.GetByAddress(proposerAddr)
	if proposerIdx < 0 {
		logger.Error("invalid index of the proposer",
			"addr", proposerAddr.String(), "index", proposerIdx)
		return validators
	}

	// return early if the committee size is 1
	if committeeSize == 1 {
		return []istanbul.Validator{proposer}
	}

	// find the next proposer
	nextProposer := valSet.selector(valSet, proposer.Address(), view.Round.Uint64())
	nextProposerIdx, _ := valSet.GetByAddress(nextProposer.Address())
	if nextProposerIdx < 0 {
		logger.Error("invalid index of the next proposer",
			"addr", nextProposer.Address().String(), "index", nextProposerIdx)
		return validators
	}

	// seed will be used to select a random committee
	seed, err := ConvertHashToSeed(prevHash)
	if fork.Rules(view.Sequence).IsIstanbul {
		seed += view.Round.Int64()
	}
	if err != nil {
		logger.Error("failed to convert hash to seed", "prevHash", prevHash, "err", err)
		return validators
	}

	// select a random committee
	committee := SelectRandomCommittee(validators, committeeSize, seed, proposerIdx, nextProposerIdx)
	if committee == nil {
		committee = validators
	}

	logger.Trace("composed committee", "prevHash", prevHash.Hex(), "proposerAddr", proposerAddr,
		"committee", committee, "committee size", len(committee), "valSet.subSize", committeeSize)

	return committee
}

func (valSet *defaultSet) CheckInSubList(prevHash common.Hash, view *istanbul.View, addr common.Address) bool {
	for _, val := range valSet.SubList(prevHash, view) {
		if val.Address() == addr {
			return true
		}
	}
	return false
}

func (valSet *defaultSet) IsSubSet() bool {
	return valSet.Size() > valSet.subSize
}

func (valSet *defaultSet) GetByIndex(i uint64) istanbul.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	if i < uint64(len(valSet.validators)) {
		return valSet.validators[i]
	}
	return nil
}

func (valSet *defaultSet) GetByAddress(addr common.Address) (int, istanbul.Validator) {
	for i, val := range valSet.List() {
		if addr == val.Address() {
			return i, val
		}
	}
	// TODO-Klaytn-Istanbul: Enable this log when non-committee nodes don't call `core.startNewRound()`
	// logger.Warn("failed to find an address in the validator list",
	// 	"address", addr, "validatorAddrs", valSet.validators.AddressStringList())
	return -1, nil
}

func (valSet *defaultSet) GetDemotedByAddress(addr common.Address) (int, istanbul.Validator) {
	return -1, nil
}

func (valSet *defaultSet) GetProposer() istanbul.Validator {
	proposer := valSet.proposer.Load()
	if proposer == nil {
		logger.Error("Proposer is nil", "validators", valSet.validators)
		return nil
	}
	return proposer.(istanbul.Validator)
}

func (valSet *defaultSet) IsProposer(address common.Address) bool {
	_, val := valSet.GetByAddress(address)
	return reflect.DeepEqual(valSet.GetProposer(), val)
}

func (valSet *defaultSet) CalcProposer(lastProposer common.Address, round uint64) {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	if len(valSet.validators) == 0 {
		logger.Error("len of validators is 0, Proposer is nil", "validators", valSet.validators)
		return
	}

	valSet.proposer.Store(valSet.selector(valSet, lastProposer, round))
}

func calcSeed(valSet istanbul.ValidatorSet, proposer common.Address, round uint64) uint64 {
	offset := 0
	if idx, val := valSet.GetByAddress(proposer); val != nil {
		offset = idx
	}
	return uint64(offset) + round
}

func emptyAddress(addr common.Address) bool {
	return addr == common.Address{}
}

func roundRobinProposer(valSet istanbul.ValidatorSet, proposer common.Address, round uint64) istanbul.Validator {
	seed := uint64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round) + 1
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func stickyProposer(valSet istanbul.ValidatorSet, proposer common.Address, round uint64) istanbul.Validator {
	seed := uint64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round)
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func (valSet *defaultSet) AddValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()
	for _, v := range valSet.validators {
		if v.Address() == address {
			return false
		}
	}
	valSet.validators = append(valSet.validators, New(address))
	// TODO: we may not need to re-sort it again
	// sort validator
	sort.Sort(valSet.validators)
	return true
}

func (valSet *defaultSet) RemoveValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	for i, v := range valSet.validators {
		if v.Address() == address {
			valSet.validators = append(valSet.validators[:i], valSet.validators[i+1:]...)
			return true
		}
	}
	return false
}

func (valSet *defaultSet) Copy() istanbul.ValidatorSet {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	addresses := make([]common.Address, 0, len(valSet.validators))
	for _, v := range valSet.validators {
		addresses = append(addresses, v.Address())
	}

	newValSet := NewSubSet(addresses, valSet.policy, valSet.subSize).(*defaultSet)
	_, proposer := newValSet.GetByAddress(valSet.GetProposer().Address())
	newValSet.proposer.Store(proposer)
	return newValSet
}

func (valSet *defaultSet) F() int {
	if valSet.Size() > valSet.subSize {
		return int(math.Ceil(float64(valSet.subSize)/3)) - 1
	} else {
		return int(math.Ceil(float64(valSet.Size())/3)) - 1
	}
}

func (valSet *defaultSet) Policy() istanbul.ProposerPolicy { return valSet.policy }

func (valSet *defaultSet) Refresh(hash common.Hash, blockNum uint64, config *params.ChainConfig, isSingle bool, governingNode common.Address, minStaking uint64) error {
	return nil
}
func (valSet *defaultSet) SetBlockNum(blockNum uint64)     { /* Do nothing */ }
func (valSet *defaultSet) Proposers() []istanbul.Validator { return nil }
func (valSet *defaultSet) TotalVotingPower() uint64 {
	sum := uint64(0)
	for _, v := range valSet.List() {
		sum += v.VotingPower()
	}
	return sum
}

func (valSet *defaultSet) Selector(valS istanbul.ValidatorSet, lastProposer common.Address, round uint64) istanbul.Validator {
	return valSet.selector(valS, lastProposer, round)
}
