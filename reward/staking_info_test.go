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

package reward

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

func TestStakingInfo_GetIndexByNodeAddress(t *testing.T) {
	testdata := []common.Address{
		common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"),
		common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"),
		common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"),
		common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"),
	}
	testCases := []struct {
		address common.Address
		index   int
		err     error
	}{
		{common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"), 0, nil},
		{common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"), 1, nil},
		{common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"), 2, nil},
		{common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"), 3, nil},
		{common.StringToAddress("0x027AbB8c9f952cfFf01B1707fF14E2CB5D439502"), AddrNotFoundInCouncilNodes, ErrAddrNotInStakingInfo},
	}

	stakingInfo := newEmptyStakingInfo(0)
	stakingInfo.CouncilNodeAddrs = testdata

	for i := 0; i < len(testCases); i++ {
		result, err := stakingInfo.GetIndexByNodeAddress(testCases[i].address)
		assert.Equal(t, testCases[i].index, result)
		assert.Equal(t, testCases[i].err, err)
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
		err           error
	}{
		{common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"), 100, nil},
		{common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"), 200, nil},
		{common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"), 300, nil},
		{common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"), 400, nil},
		{common.StringToAddress("0x027AbB8c9f952cfFf01B1707fF14E2CB5D439502"), 0, ErrAddrNotInStakingInfo},
	}

	stakingInfo := newEmptyStakingInfo(0)
	stakingInfo.CouncilNodeAddrs = testdata.address
	stakingInfo.CouncilStakingAmounts = testdata.stakingAmount

	for i := 0; i < len(testCases); i++ {
		result, err := stakingInfo.GetStakingAmountByNodeId(testCases[i].address)
		assert.Equal(t, testCases[i].stakingAmount, result)
		assert.Equal(t, testCases[i].err, err)
	}
}

func TestStakingInfo_String(t *testing.T) {
	testCases := []*StakingInfo{
		newEmptyStakingInfo(0),
		{
			1,
			[]common.Address{
				common.StringToAddress("node address 1"),
				common.StringToAddress("node address 2"),
				common.StringToAddress("node address 3"),
				common.StringToAddress("node address 4"),
			},
			[]common.Address{
				common.StringToAddress("staking address 1"),
				common.StringToAddress("staking address 2"),
				common.StringToAddress("staking address 3"),
				common.StringToAddress("staking address 4"),
			},
			[]common.Address{
				common.StringToAddress("reward address 1"),
				common.StringToAddress("reward address 2"),
				common.StringToAddress("reward address 3"),
				common.StringToAddress("reward address 4"),
			},
			common.StringToAddress("kir address"),
			common.StringToAddress("poc address"),
			false,
			0.0,
			[]uint64{5000000, 5000000, 5000000, 5000000},
		},
		{
			86400,
			[]common.Address{
				common.HexToAddress("0x8aD8F547fa00f58A8c4fb3B671Ee5f1A75bA028a"),
				common.HexToAddress("0xB2AAda7943919e82143324296987f6091F3FDC9e"),
				common.HexToAddress("0xD95c70710f07A3DaF7ae11cFBa10c789da3564D0"),
				common.HexToAddress("0xC704765db1d21C2Ea6F7359dcB8FD5233DeD16b5"),
			},
			[]common.Address{
				common.HexToAddress("0x4dd324F9821485caE941640B32c3Bcf1fA6E93E6"),
				common.HexToAddress("0x0d5Df5086B5f86f748dFaed5779c3f862C075B1f"),
				common.HexToAddress("0xD3Ff05f00491571E86A3cc8b0c320aA76D7413A5"),
				common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f"),
			},
			[]common.Address{
				common.HexToAddress("0x241c793A9AD555f52f6C3a83afe6178408796ab2"),
				common.HexToAddress("0x79b427Fb77077A9716E08D049B0e8f36Abfc8E2E"),
				common.HexToAddress("0x62E47d858bf8513fc401886B94E33e7DCec2Bfb7"),
				common.HexToAddress("0xf275f9f4c0d375F9E3E50370f93b504A1e45dB09"),
			},
			common.HexToAddress("0x136807B12327a8AfF9831F09617dA1B9D398cda2"),
			common.HexToAddress("0x46bA8F7538CD0749e572b2631F9FB4Ce3653AFB8"),
			true,
			0.5,
			[]uint64{10000000, 20000000, 30000000, 40000000},
		},
	}

	for _, testStakingInfo := range testCases {
		resultStr := testStakingInfo.String()
		t.Logf("%s", resultStr)
		resultByteArr := []byte(resultStr)
		resultStakingInfo := &StakingInfo{}

		err := json.Unmarshal(resultByteArr, resultStakingInfo)
		assert.NoError(t, err)
		assert.Equal(t, testStakingInfo, resultStakingInfo)
	}
}

func TestCalcGiniCoefficient(t *testing.T) {
	testCase := []struct {
		testdata []float64
		result   float64
	}{
		{[]float64{1, 1, 1}, 0.0},
		{[]float64{0, 8, 0, 0, 0}, 0.8},
		{[]float64{5, 4, 3, 2, 1}, 0.27},
	}

	for i := 0; i < len(testCase); i++ {
		result := CalcGiniCoefficient(testCase[i].testdata)
		assert.Equal(t, testCase[i].result, result)
	}
}

func TestGiniReflectToExpectedCCO(t *testing.T) {
	testCase := []struct {
		ccoToken        []float64
		beforeReflected []float64
		adjustment      []float64
		afterReflected  []float64
	}{
		{
			[]float64{
				66666667, 233333333, 5000000, 5000000, 5000000,
				77777778, 5000000, 33333333, 20000000, 16666667,
				10000000, 5000000, 5000000, 5000000, 5000000,
				5000000, 5000000, 5000000, 5000000, 5000000,
				5000000,
			},
			[]float64{13, 44, 1, 1, 1, 15, 1, 6, 4, 3, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]float64{42612, 89426, 9202, 9202, 9202, 46682, 9202, 28275, 20900, 18762, 13868, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202},
			[]float64{11, 23, 2, 2, 2, 12, 2, 7, 5, 5, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		},
		{
			[]float64{
				400000000, 233333333, 233333333, 150000000, 108333333,
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
		stakingInfo := newEmptyStakingInfo(uint64(1))

		weights := make([]float64, len(testCase[i].ccoToken))
		tokenListToCalcGini := make([]float64, len(testCase[i].ccoToken))
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

// TestStakingInfoJSON tests marshalling and unmarshalling StakingInfo
// StakingInfo is marshaled before storing to DB.
func TestStakingInfoJSON(t *testing.T) {
	enc := newEmptyStakingInfo(1234)

	s, err := json.Marshal(enc)
	if err != nil {
		t.Fatal(err)
	}

	dec := new(StakingInfo)
	err = json.Unmarshal(s, dec)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(enc, dec) {
		t.Fatal(errors.New("problem while marshaling or unmarshaling"))
	}
}
