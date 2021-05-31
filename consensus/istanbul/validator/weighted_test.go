// Copyright 2019 The klaytn Authors
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

package validator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto"
	"github.com/stretchr/testify/assert"
)

func TestNewWeightedCouncil(t *testing.T) {
	const ValCnt = 100
	var validators []istanbul.Validator
	var rewardAddrs []common.Address
	var votingPowers []uint64

	// Create 100 validators with random addresses
	b := []byte{}
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		val := New(addr)
		validators = append(validators, val)
		b = append(b, val.Address().Bytes()...)

		rewardKey, _ := crypto.GenerateKey()
		rewardAddr := crypto.PubkeyToAddress(rewardKey.PublicKey)
		rewardAddrs = append(rewardAddrs, rewardAddr)

		votingPowers = append(votingPowers, uint64(1))
	}

	// Create ValidatorSet
	valSet := NewWeightedCouncil(ExtractValidators(b), nil, rewardAddrs, votingPowers, nil, istanbul.WeightedRandom, 21, 0, 0, nil)
	if valSet == nil {
		t.Errorf("the validator byte array cannot be parsed")
		t.FailNow()
	}

	// Check validators sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		val := valSet.GetByIndex(uint64(i))
		nextVal := valSet.GetByIndex(uint64(i + 1))
		if strings.Compare(val.String(), nextVal.String()) >= 0 {
			t.Errorf("validator set is not sorted in descending order")
		}
	}
}

func TestNormalWeightedCouncil(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)

	rewardKey1, _ := crypto.GenerateKey()
	rewardAddr1 := crypto.PubkeyToAddress(rewardKey1.PublicKey)

	rewardKey2, _ := crypto.GenerateKey()
	rewardAddr2 := crypto.PubkeyToAddress(rewardKey2.PublicKey)

	votingPower1 := uint64(1)
	votingPower2 := uint64(2)

	weight1 := uint64(1)
	weight2 := uint64(2)

	val1 := newWeightedValidator(addr1, rewardAddr1, votingPower1, weight1)
	val2 := newWeightedValidator(addr2, rewardAddr2, votingPower2, weight2)

	valSet := NewWeightedCouncil([]common.Address{addr1, addr2}, nil, []common.Address{rewardAddr1, rewardAddr2}, []uint64{votingPower1, votingPower2}, []uint64{weight1, weight2}, istanbul.WeightedRandom, 21, 0, 0, nil)
	if valSet == nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	// check size
	if size := valSet.Size(); size != 2 {
		t.Errorf("the size of validator set is wrong: have %v, want 2", size)
	}

	// test get by index
	if val := valSet.GetByIndex(uint64(0)); !reflect.DeepEqual(val, val1) {
		t.Errorf("validator mismatch:")
		t.Errorf("  Address(): %v, %v", val.Address(), val1.Address())
		t.Errorf("  String(): %v, %v", val.String(), val1.String())
		t.Errorf("  RewardAddresS(): %v, %v", val.RewardAddress(), val1.RewardAddress())
		t.Errorf("  VotingPower(): %v, %v", val.VotingPower(), val1.VotingPower())
		t.Errorf("  Weight(): %v, %v", val.Weight(), val1.Weight())
	}

	// test get by invalid index
	if val := valSet.GetByIndex(uint64(2)); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}

	// test get by address
	if _, val := valSet.GetByAddress(addr2); !reflect.DeepEqual(val, val2) {
		t.Errorf("validator mismatch: have %v, want %v", val, val2)
	}

	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if _, val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}

	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}

	// test calculate proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}

	valSet.CalcProposer(lastProposer, uint64(1))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}

	valSet.CalcProposer(lastProposer, uint64(2))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}

	valSet.CalcProposer(lastProposer, uint64(5))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}

	// test empty last proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}

func TestEmptyWeightedCouncil(t *testing.T) {
	valSet := NewWeightedCouncil(ExtractValidators([]byte{}), nil, nil, nil, nil, istanbul.WeightedRandom, 0, 0, 0, &blockchain.BlockChain{})
	if valSet == nil {
		t.Errorf("validator set should not be nil")
	}
}

func TestNewWeightedCouncil_InvalidPolicy(t *testing.T) {
	// Invalid proposer policy
	valSet := NewWeightedCouncil(ExtractValidators([]byte{}), nil, nil, nil, nil, istanbul.Sticky, 0, 0, 0, &blockchain.BlockChain{})
	assert.Equal(t, (*weightedCouncil)(nil), valSet)

	valSet = NewWeightedCouncil(ExtractValidators([]byte{}), nil, nil, nil, nil, istanbul.RoundRobin, 0, 0, 0, &blockchain.BlockChain{})
	assert.Equal(t, (*weightedCouncil)(nil), valSet)
}

func TestNewWeightedCouncil_IncompleteParams(t *testing.T) {
	const ValCnt = 3
	var validators []istanbul.Validator
	var rewardAddrs []common.Address
	var votingPowers []uint64
	var weights []uint64

	// Create 3 validators with random addresses
	b := []byte{}
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		val := New(addr)
		validators = append(validators, val)
		b = append(b, val.Address().Bytes()...)

		rewardKey, _ := crypto.GenerateKey()
		rewardAddr := crypto.PubkeyToAddress(rewardKey.PublicKey)
		rewardAddrs = append(rewardAddrs, rewardAddr)

		votingPowers = append(votingPowers, uint64(1))
		weights = append(weights, uint64(1))
	}

	// No validator address
	valSet := NewWeightedCouncil(ExtractValidators([]byte{}), nil, rewardAddrs, votingPowers, weights, istanbul.WeightedRandom, 0, 0, 0, &blockchain.BlockChain{})
	assert.Equal(t, (*weightedCouncil)(nil), valSet)

	// Incomplete rewardAddrs
	incompleteRewardAddrs := make([]common.Address, 1)
	valSet = NewWeightedCouncil(ExtractValidators(b), nil, incompleteRewardAddrs, nil, nil, istanbul.WeightedRandom, 0, 0, 0, &blockchain.BlockChain{})
	assert.Equal(t, (*weightedCouncil)(nil), valSet)

	// Incomplete rewardAddrs
	incompleteVotingPowers := make([]uint64, 1)
	valSet = NewWeightedCouncil(ExtractValidators(b), nil, nil, incompleteVotingPowers, nil, istanbul.WeightedRandom, 0, 0, 0, &blockchain.BlockChain{})
	assert.Equal(t, (*weightedCouncil)(nil), valSet)

	// Incomplete rewardAddrs
	incompleteWeights := make([]uint64, 1)
	valSet = NewWeightedCouncil(ExtractValidators(b), nil, nil, nil, incompleteWeights, istanbul.WeightedRandom, 0, 0, 0, &blockchain.BlockChain{})
	assert.Equal(t, (*weightedCouncil)(nil), valSet)
}
