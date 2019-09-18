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
	"fmt"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/reward"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
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

	proposer    istanbul.Validator
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
		valSet.proposer = valSet.GetByIndex(0)
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
		valSet.proposer = valSet.GetByIndex(0)
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

func (valSet *defaultSet) SetSubGroupSize(size uint64) {
	valSet.subSize = size
}

func (valSet *defaultSet) List() []istanbul.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return valSet.validators
}

func (valSet *defaultSet) SubList(prevHash common.Hash, view *istanbul.View) []istanbul.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	if uint64(len(valSet.validators)) <= valSet.subSize {
		return valSet.validators
	}
	hashstring := strings.TrimPrefix(prevHash.Hex(), "0x")
	if len(hashstring) > 15 {
		hashstring = hashstring[:15]
	}
	seed, err := strconv.ParseInt(hashstring, 16, 64)
	if err != nil {
		logger.Error("input", "hash", prevHash.Hex())
		logger.Error("fail to make sub-list of validators", "seed", seed, "err", err)
		return valSet.validators
	}

	// shuffle
	subset := make([]istanbul.Validator, valSet.subSize)
	subset[0] = valSet.GetProposer()
	// next proposer
	subset[1] = valSet.selector(valSet, subset[0].Address(), view.Round.Uint64())

	proposerIdx, _ := valSet.GetByAddress(subset[0].Address())
	nextproposerIdx, _ := valSet.GetByAddress(subset[1].Address())

	if proposerIdx == nextproposerIdx {
		logger.Error("fail to make propser", "current proposer idx", proposerIdx, "next idx", nextproposerIdx)
	}

	limit := len(valSet.validators)
	picker := rand.New(rand.NewSource(seed))

	pickSize := limit - 2
	indexs := make([]int, pickSize)
	idx := 0
	for i := 0; i < limit; i++ {
		if i != proposerIdx && i != nextproposerIdx {
			indexs[idx] = i
			idx++
		}
	}
	for i := 0; i < pickSize; i++ {
		randIndex := picker.Intn(pickSize)
		indexs[i], indexs[randIndex] = indexs[randIndex], indexs[i]
	}

	for i := uint64(0); i < valSet.subSize-2; i++ {
		subset[i+2] = valSet.validators[indexs[i]]
	}

	if prevHash.Hex() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		logger.Error("### subList", "prevHash", prevHash.Hex())
	}

	return subset
}

func (valSet *defaultSet) SubListWithProposer(prevHash common.Hash, proposer common.Address, view *istanbul.View) []istanbul.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	if uint64(len(valSet.validators)) <= valSet.subSize {
		return valSet.validators
	}
	hashstring := strings.TrimPrefix(prevHash.Hex(), "0x")
	if len(hashstring) > 15 {
		hashstring = hashstring[:15]
	}
	seed, err := strconv.ParseInt(hashstring, 16, 64)
	if err != nil {
		logger.Error("input", "hash", prevHash.Hex())
		logger.Error("fail to make sub-list of validators", "seed", seed, "err", err)
		return valSet.validators
	}

	// shuffle
	subset := make([]istanbul.Validator, valSet.subSize)
	subset[0] = New(proposer)
	// next proposer
	subset[1] = valSet.selector(valSet, subset[0].Address(), view.Round.Uint64())

	proposerIdx, _ := valSet.GetByAddress(subset[0].Address())
	nextproposerIdx, _ := valSet.GetByAddress(subset[1].Address())

	// TODO-Klaytn-RemoveLater remove this check code if the implementation is stable.
	if proposerIdx < 0 || nextproposerIdx < 0 {
		vals := "["
		for _, v := range valSet.validators {
			vals += fmt.Sprintf("%s,", v.Address().Hex())
		}
		vals += "]"
		logger.Error("idx should not be negative!", "proposerIdx", proposerIdx, "nextproposerIdx", nextproposerIdx, "proposer", subset[0].Address().Hex(),
			"nextproposer", subset[1].Address().Hex(), "validators", vals)
	}

	if proposerIdx == nextproposerIdx {
		logger.Error("fail to make propser", "current proposer idx", proposerIdx, "next idx", nextproposerIdx)
	}

	limit := len(valSet.validators)
	picker := rand.New(rand.NewSource(seed))

	pickSize := limit - 2
	indexs := make([]int, pickSize)
	idx := 0
	for i := 0; i < limit; i++ {
		if i != proposerIdx && i != nextproposerIdx {
			indexs[idx] = i
			idx++
		}
	}
	for i := 0; i < pickSize; i++ {
		randIndex := picker.Intn(pickSize)
		indexs[i], indexs[randIndex] = indexs[randIndex], indexs[i]
	}

	for i := uint64(0); i < valSet.subSize-2; i++ {
		subset[i+2] = valSet.validators[indexs[i]]
	}

	if prevHash.Hex() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		logger.Error("### subList", "prevHash", prevHash.Hex())
	}

	return subset
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
	return -1, nil
}

func (valSet *defaultSet) GetProposer() istanbul.Validator {
	return valSet.proposer
}

func (valSet *defaultSet) IsProposer(address common.Address) bool {
	_, val := valSet.GetByAddress(address)
	return reflect.DeepEqual(valSet.GetProposer(), val)
}

func (valSet *defaultSet) CalcProposer(lastProposer common.Address, round uint64) {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	if len(valSet.validators) == 0 {
		valSet.proposer = nil
		return
	}

	valSet.proposer = valSet.selector(valSet, lastProposer, round)
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
	_, proposer := newValSet.GetByAddress(valSet.proposer.Address())
	newValSet.proposer = proposer
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

func (valSet *defaultSet) Refresh(hash common.Hash, blockNum uint64, stakingManager *reward.StakingManager) error {
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
