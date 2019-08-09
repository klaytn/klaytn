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
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(301)), common.BigToAddress(big.NewInt(302)), common.BigToAddress(big.NewInt(303))},
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000},
			},
			[]float64{5000000, 5000000, 5000000},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(301)), common.BigToAddress(big.NewInt(302)), common.BigToAddress(big.NewInt(303))},
				CouncilStakingAmounts: []uint64{7000000, 5000000, 10000000},
			},
			[]float64{7000000, 5000000, 10000000},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(301)), common.BigToAddress(big.NewInt(302)), common.BigToAddress(big.NewInt(303)), common.BigToAddress(big.NewInt(304)), common.BigToAddress(big.NewInt(301))},
				CouncilStakingAmounts: []uint64{5000000, 5000000, 5000000, 5000000, 5000000},
			},
			[]float64{10000000, 5000000, 5000000, 5000000},
		},
		{
			[]common.Address{common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(101))},
			&reward.StakingInfo{
				CouncilNodeAddrs:      []common.Address{common.BigToAddress(big.NewInt(101)), common.BigToAddress(big.NewInt(102)), common.BigToAddress(big.NewInt(103)), common.BigToAddress(big.NewInt(104)), common.BigToAddress(big.NewInt(901)), common.BigToAddress(big.NewInt(902))},
				CouncilRewardAddrs:    []common.Address{common.BigToAddress(big.NewInt(301)), common.BigToAddress(big.NewInt(302)), common.BigToAddress(big.NewInt(303)), common.BigToAddress(big.NewInt(304)), common.BigToAddress(big.NewInt(301)), common.BigToAddress(big.NewInt(302))},
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

func TestWeightedCouncil_calcWeight(t *testing.T) {

}

func TestWeightedCouncil_validatorWeightWithStakingInfo(t *testing.T) {

}
