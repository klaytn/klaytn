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

/*
Multiple Staking Contracts

Validators can deploy multiple staking contracts.
If a validator wants to deploy additional staking contracts, those staking contracts should have same rewardAddress.
StakingAmounts of staking contracts with a same rewardAddress will be added and it is reflected to a probability of becoming a block proposer.

Testing

StakingInfos are data from addressBook.
A StakingInfo has lists of addresses and stakingAmount.
They are matched by an index. Values of the lists with a same index are from a same staking contract.

All addresses used in tests are made by 3 digits number.
NodeAddress : begin with 1
rewardAddress : begin with 2
NodeAddress of additional staking contract : begin with 9
*/
package validator

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/reward"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func newTestWeightedCouncil(nodeAddrs []common.Address) *weightedCouncil {
	return NewWeightedCouncil(nodeAddrs, nil, make([]uint64, len(nodeAddrs)), nil, istanbul.WeightedRandom, 0, 0, 0, nil)
}

// check if validators and stakingAmount from stakingInfo is matched well.
// stakingAmounts of additional staking contracts will be added to stakingAmounts of validators which have the same reward address.
// input
//  - validator and stakingInfo is matched by a nodeAddress.
// output
//  - weightedValidators are sorted by nodeAddress
//  - stakingAmounts should be same as expectedStakingAmounts
func TestWeightedCouncil_getStakingAmountsOfValidators(t *testing.T) {
	testCases := []struct {
		validators             []common.Address
		stakingInfo            *reward.StakingInfo
		expectedStakingAmounts []float64
	}{
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203))},
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000},
			},
			[]float64{5000000, 5000000, 5000000},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203))},
				CouncilStakingAmounts: []uint64{7000000, 5000000, 10000000},
			},
			[]float64{7000000, 5000000, 10000000},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204)), common.BigToAddress(big.NewInt(201))},
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000, 5000000, 5000000},
			},
			[]float64{10000000, 5000000, 5000000, 5000000},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(101))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901)), common.BigToAddress(big.NewInt(902))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204)), common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202))},
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000, 5000000, 5000000, 5000000},
			},
			[]float64{10000000, 10000000, 5000000, 5000000},
		},
	}
	for _, testCase := range testCases {
		council := newTestWeightedCouncil(testCase.validators)

		weightedValidators, stakingAmounts, err := council.getStakingAmountsOfValidators(testCase.stakingInfo)

		assert.NoError(t, err)
		assert.Equal(t, len(testCase.validators), len(weightedValidators))
		for i := 0; i < len(testCase.validators); i++ {
			assert.Contains(t, testCase.validators, weightedValidators[i].address)
		}
		assert.Equal(t, testCase.expectedStakingAmounts, stakingAmounts)
	}
}

// calcTotalAmount calculates totalAmount of stakingAmounts and gini coefficient if UseGini is true.
// if UseGini is true, gini is calculated and reflected to stakingAmounts.
func TestCalcTotalAmount(t *testing.T) {
	testCases := []struct {
		stakingInfo            *reward.StakingInfo
		stakingAmounts         []float64
		expectedGini           float64
		expectedTotalAmount    float64
		expectedStakingAmounts []float64
	}{
		{
			&reward.StakingInfo{
				CouncilNodeAddrs: []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
				UseGini:          false,
				Gini:             reward.DefaultGiniCoefficient,
			},
			[]float64{5000000, 5000000, 5000000},
			reward.DefaultGiniCoefficient,
			15000000,
			[]float64{5000000, 5000000, 5000000},
		},
		{
			&reward.StakingInfo{
				CouncilNodeAddrs: []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
				UseGini:          true,
				Gini:             reward.DefaultGiniCoefficient,
			},
			[]float64{5000000, 5000000, 5000000},
			0,
			15000000,
			[]float64{5000000, 5000000, 5000000},
		},

		{
			&reward.StakingInfo{
				CouncilNodeAddrs: []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(105))},
				UseGini:          true,
				Gini:             reward.DefaultGiniCoefficient,
			},
			[]float64{10000000, 20000000, 30000000, 40000000, 50000000},
			0.27,
			3779508,
			[]float64{324946, 560845, 771786, 967997, 1153934},
		},
	}
	for _, testCase := range testCases {
		stakingAmounts := testCase.stakingAmounts
		totalAmount := calcTotalAmount(testCase.stakingInfo, stakingAmounts)

		assert.Equal(t, testCase.expectedGini, testCase.stakingInfo.Gini)
		assert.Equal(t, testCase.expectedTotalAmount, totalAmount)
		assert.Equal(t, testCase.expectedStakingAmounts, stakingAmounts)
	}
}

// calcWeight calculates weights and saves them to validators.
// weights are the ratio of each stakingAmount to totalStaking
func TestCalcWeight(t *testing.T) {
	testCases := []struct {
		weightedValidators []*weightedValidator
		stakingAmounts     []float64
		totalStaking       float64
		expectedWeights    []int64
	}{
		{
			[]*weightedValidator{
				{}, {}, {},
			},
			[]float64{0, 0, 0},
			0,
			[]int64{0, 0, 0},
		},
		{
			[]*weightedValidator{
				{}, {}, {},
			},
			[]float64{5000000, 5000000, 5000000},
			15000000,
			[]int64{33, 33, 33},
		},
		{
			[]*weightedValidator{
				{}, {}, {}, {},
			},
			[]float64{5000000, 10000000, 5000000, 5000000},
			25000000,
			[]int64{20, 40, 20, 20},
		},
		{
			[]*weightedValidator{
				{}, {}, {}, {}, {},
			},
			[]float64{324946, 560845, 771786, 967997, 1153934},
			3779508,
			[]int64{9, 15, 20, 26, 31},
		},
	}
	for _, testCase := range testCases {
		calcWeight(testCase.weightedValidators, testCase.stakingAmounts, testCase.totalStaking)
		for i, weight := range testCase.expectedWeights {
			assert.Equal(t, weight, testCase.weightedValidators[i].Weight())
		}
	}
}

// The test is union of above tests.
// Weight should be calculated exactly by a validator list and a stakingInfo given
func TestWeightedCouncil_validatorWeightWithStakingInfo(t *testing.T) {
	testCases := []struct {
		validators      []common.Address
		stakingInfo     *reward.StakingInfo
		expectedWeights []int64
	}{
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203))},
				UseGini:               false,
				Gini:                  reward.DefaultGiniCoefficient,
				CouncilStakingAmounts: []uint64{0, 0, 0},
			},
			[]int64{0, 0, 0},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204))},
				UseGini:               true,
				Gini:                  reward.DefaultGiniCoefficient,
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000, 5000000},
			},
			[]int64{25, 25, 25, 25},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(105))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(105))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204)), common.BigToAddress(big.NewInt(205))},
				UseGini:               true,
				Gini:                  reward.DefaultGiniCoefficient,
				CouncilStakingAmounts: []uint64{10000000, 20000000, 30000000, 40000000, 50000000},
			},
			[]int64{9, 15, 20, 26, 31},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204)), common.BigToAddress(big.NewInt(201))},
				UseGini:               false,
				Gini:                  reward.DefaultGiniCoefficient,
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000, 5000000, 5000000},
			},
			[]int64{40, 20, 20, 20},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204)), common.BigToAddress(big.NewInt(201))},
				UseGini:               true,
				Gini:                  reward.DefaultGiniCoefficient,
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000, 5000000, 5000000},
			},
			[]int64{38, 21, 21, 21},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(101))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901)), common.BigToAddress(big.NewInt(902))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204)), common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202))},
				UseGini:               true,
				Gini:                  reward.DefaultGiniCoefficient,
				CouncilStakingAmounts: []uint64{10000000, 5000000, 20000000, 5000000, 5000000, 5000000},
			},
			[]int64{29, 21, 37, 12},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(105))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901)), common.BigToAddress(big.NewInt(902))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202)), common.BigToAddress(big.NewInt(203)), common.BigToAddress(big.NewInt(204)), common.BigToAddress(big.NewInt(201)), common.BigToAddress(big.NewInt(202))},
				UseGini:               true,
				Gini:                  reward.DefaultGiniCoefficient,
				CouncilStakingAmounts: []uint64{10000000, 5000000, 20000000, 5000000, 5000000, 5000000},
			},
			[]int64{29, 22, 36, 13, 1},
		},
	}
	for _, testCase := range testCases {
		council := newTestWeightedCouncil(testCase.validators)

		weightedValidators, stakingAmounts, err := council.getStakingAmountsOfValidators(testCase.stakingInfo)
		assert.NoError(t, err)
		totalStaking := calcTotalAmount(testCase.stakingInfo, stakingAmounts)
		calcWeight(weightedValidators, stakingAmounts, totalStaking)

		for i, weight := range testCase.expectedWeights {
			assert.Equal(t, weight, weightedValidators[i].Weight())
		}
	}
}
