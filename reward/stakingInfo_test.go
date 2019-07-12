package reward

import (
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestStakingInfo_GetIndexByNodeId(t *testing.T) {
	testdata := []common.Address{
		common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"),
		common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"),
		common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"),
		common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"),
	}
	testCases := []struct {
		address common.Address
		index   int
	}{
		{common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"), 0},
		{common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"), 1},
		{common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"), 2},
		{common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"), 3},
	}

	stakingInfo, error := newEmptyStakingInfo(0)
	if error != nil {
		t.Errorf("error has occurred. error : %v", error)
		t.FailNow()
	}
	stakingInfo.CouncilNodeIds = testdata

	for i := 0; i < len(testCases); i++ {
		result := stakingInfo.GetIndexByNodeId(testCases[i].address)
		if result == AddrNotFoundInCouncilNodes {
			t.Errorf("Address is not found in stakingInfo")
		}
		assert.Equal(t, testCases[i].index, result)
	}
}

func TestStakingInfo_GetStakingAmountByNodeId(t *testing.T) {
	testdata := struct {
		address       []common.Address
		stakingAmount []uint64
	}{
		[]common.Address{
			common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"),
			common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"),
			common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"),
			common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"),
		},
		[]uint64{
			100, 200, 300, 400,
		},
	}
	testCases := []struct {
		address       common.Address
		stakingAmount uint64
	}{
		{common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"), 100},
		{common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"), 200},
		{common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"), 300},
		{common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"), 400},
	}

	stakingInfo, error := newEmptyStakingInfo(0)
	if error != nil {
		t.Errorf("error has occurred. error : %v", error)
		t.FailNow()
	}
	stakingInfo.CouncilNodeIds = testdata.address
	stakingInfo.CouncilStakingAmounts = testdata.stakingAmount

	for i := 0; i < len(testCases); i++ {
		result := stakingInfo.GetStakingAmountByNodeId(testCases[i].address)
		assert.Equal(t, testCases[i].stakingAmount, result)
	}
}

func TestCalcGiniCoefficient(t *testing.T) {
	testCase := []struct {
		testdata []uint64
		result   float64
	}{
		{[]uint64{1, 1, 1}, 0.0},
		{[]uint64{0, 8, 0, 0, 0}, 0.8},
		{[]uint64{5, 4, 3, 2, 1}, 0.27},
	}

	for i := 0; i < len(testCase); i++ {
		result := CalcGiniCoefficient(testCase[i].testdata)

		if result != testCase[i].result {
			t.Errorf("The result is different from the expected result. result : %v, expected : %v", result, testCase[i].result)
		}
	}
}

func TestGiniReflectToExpectedCCO(t *testing.T) {
	testCase := []struct {
		ccoToken        []uint64
		beforeReflected []float64
		adjustment      []float64
		afterReflected  []float64
	}{
		{[]uint64{66666667, 233333333, 5000000, 5000000, 5000000,
			77777778, 5000000, 33333333, 20000000, 16666667,
			10000000, 5000000, 5000000, 5000000, 5000000,
			5000000, 5000000, 5000000, 5000000, 5000000,
			5000000,
		},
			[]float64{13, 44, 1, 1, 1, 15, 1, 6, 4, 3, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]float64{42612, 89426, 9202, 9202, 9202, 46682, 9202, 28275, 20900, 18762, 13868, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202},
			[]float64{11, 23, 2, 2, 2, 12, 2, 7, 5, 5, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		},
		{[]uint64{400000000, 233333333, 233333333, 150000000, 108333333,
			83333333, 66666667, 33333333, 20000000, 16666667,
			10000000, 5000000, 5000000, 5000000, 5000000,
			5000000, 5000000, 5000000, 5000000, 5000000,
			5000000,
		},
			[]float64{28, 17, 17, 11, 8, 6, 5, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]float64{123020, 89426, 89426, 68853, 56793, 48627, 42612, 28275, 20900, 18762, 13868, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202},
			[]float64{18, 13, 13, 10, 8, 7, 6, 4, 3, 3, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	for i := 0; i < len(testCase); i++ {
		stakingInfo, _ := newEmptyStakingInfo(uint64(1))

		weights := make([]float64, len(testCase[i].ccoToken))
		tokenListToCalcGini := make([]uint64, len(testCase[i].ccoToken))
		totalAmount := 0.0
		for j := 0; j < len(testCase[i].ccoToken); j++ {
			totalAmount += float64(testCase[i].ccoToken[j])
			tokenListToCalcGini[j] = testCase[i].ccoToken[j]
		}

		for j := 0; j < len(testCase[i].ccoToken); j++ {
			weights[j] = math.Round(float64(testCase[i].ccoToken[j]) * 100 / totalAmount)
			if weights[j] < 1 {
				weights[j] = 1
			}
			if weights[j] != testCase[i].beforeReflected[j] {
				t.Errorf("normal weight is incorrect. result : %v expected : %v", weights[j], testCase[i].beforeReflected[j])
			}
		}

		stakingAmountsGiniReflected := make([]float64, len(testCase[i].ccoToken))
		totalAmountGiniReflected := 0.0
		stakingInfo.Gini = CalcGiniCoefficient(tokenListToCalcGini)

		for j := 0; j < len(stakingAmountsGiniReflected); j++ {
			stakingAmountsGiniReflected[j] = math.Round(math.Pow(float64(testCase[i].ccoToken[j]), 1.0/(1+stakingInfo.Gini)))
			totalAmountGiniReflected += stakingAmountsGiniReflected[j]
		}

		for j := 0; j < len(testCase[i].ccoToken); j++ {
			if stakingAmountsGiniReflected[j] != testCase[i].adjustment[j] {
				t.Errorf("staking amount reflected gini is different. result : %v expected : %v", stakingAmountsGiniReflected[j], testCase[i].adjustment[j])
			}
		}

		for j := 0; j < len(testCase[i].ccoToken); j++ {
			stakingAmountsGiniReflected[j] = math.Round(stakingAmountsGiniReflected[j] * 100 / totalAmountGiniReflected)
			if stakingAmountsGiniReflected[j] != testCase[i].afterReflected[j] {
				t.Errorf("weight reflected gini is different. result : %v expected : %v", stakingAmountsGiniReflected[j], testCase[i].afterReflected[j])
			}
		}
	}
}
