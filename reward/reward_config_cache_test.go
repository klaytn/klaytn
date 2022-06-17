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
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

type testGovernance struct {
	p *params.GovParamSet
}

func newTestGovernance(intMap map[int]interface{}) *testGovernance {
	p, _ := params.NewGovParamSetIntMap(intMap)
	return &testGovernance{p}
}

func newDefaultTestGovernance() *testGovernance {
	return newTestGovernance(map[int]interface{}{
		params.Epoch:               604800,
		params.Policy:              params.WeightedRandom,
		params.UnitPrice:           25000000000,
		params.MintingAmount:       "9600000000000000000",
		params.Ratio:               "34/54/12",
		params.UseGiniCoeff:        true,
		params.DeferredTxFee:       true,
		params.StakeUpdateInterval: 86400,
	})
}

func (governance *testGovernance) Params() *params.GovParamSet {
	return governance.p
}

func (governance *testGovernance) ParamsAt(num uint64) (*params.GovParamSet, error) {
	return governance.p, nil
}

func (governance *testGovernance) setTestGovernance(intMap map[int]interface{}) {
	p, _ := params.NewGovParamSetIntMap(intMap)
	governance.p = p
}

func TestRewardConfigCache_parseRewardRatio(t *testing.T) {
	testCases := []struct {
		s   string
		cn  int
		poc int
		kir int
		err error
	}{
		{"34/54/12", 34, 54, 12, nil},
		{"3/3/3", 3, 3, 3, nil},
		{"10/20/30", 10, 20, 30, nil},
		{"34,54,12", 0, 0, 0, errInvalidFormat},
		{"/", 0, 0, 0, errInvalidFormat},
		{"///", 0, 0, 0, errInvalidFormat},
		{"1//", 0, 0, 0, errParsingRatio},
		{"/1/", 0, 0, 0, errParsingRatio},
		{"//1", 0, 0, 0, errParsingRatio},
		{"1/2/3/4/", 0, 0, 0, errInvalidFormat},
		{"3.3/3.3/3.3", 0, 0, 0, errParsingRatio},
		{"a/b/c", 0, 0, 0, errParsingRatio},
	}
	rewardConfigCache := newRewardConfigCache(newDefaultTestGovernance())

	for i := 0; i < len(testCases); i++ {
		cn, poc, kir, err := rewardConfigCache.parseRewardRatio(testCases[i].s)

		assert.Equal(t, testCases[i].cn, cn)
		assert.Equal(t, testCases[i].poc, poc)
		assert.Equal(t, testCases[i].kir, kir)
		assert.Equal(t, testCases[i].err, err)
	}
}

func TestRewardConfigCache_newRewardConfig(t *testing.T) {
	testCases := []struct {
		config map[int]interface{}
		result rewardConfig
	}{
		{
			map[int]interface{}{
				params.Epoch:         604800,
				params.MintingAmount: "9600000000000000000",
				params.Ratio:         "34/54/12",
				params.UnitPrice:     25000000000,
				params.UseGiniCoeff:  true,
				params.DeferredTxFee: true,
			},
			rewardConfig{
				blockNum:      0,
				mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
				cnRatio:       big.NewInt(0).SetInt64(34),
				pocRatio:      big.NewInt(0).SetInt64(54),
				kirRatio:      big.NewInt(0).SetInt64(12),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(25000000000),
			},
		},
		{
			map[int]interface{}{
				params.Epoch:         30,
				params.MintingAmount: "10000",
				params.Ratio:         "50/30/20",
				params.UnitPrice:     50000000000,
				params.UseGiniCoeff:  true,
				params.DeferredTxFee: false,
			},
			rewardConfig{
				blockNum:      1,
				mintingAmount: big.NewInt(0).SetInt64(10000),
				cnRatio:       big.NewInt(0).SetInt64(50),
				pocRatio:      big.NewInt(0).SetInt64(30),
				kirRatio:      big.NewInt(0).SetInt64(20),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(50000000000),
			},
		},
		{
			map[int]interface{}{
				params.Epoch:         3000,
				params.MintingAmount: "100000000",
				params.Ratio:         "10/35/55",
				params.UnitPrice:     1500000000,
				params.UseGiniCoeff:  false,
				params.DeferredTxFee: true,
			},
			rewardConfig{
				blockNum:      2,
				mintingAmount: big.NewInt(0).SetInt64(100000000),
				cnRatio:       big.NewInt(0).SetInt64(10),
				pocRatio:      big.NewInt(0).SetInt64(35),
				kirRatio:      big.NewInt(0).SetInt64(55),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(1500000000),
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		testGovernance := newTestGovernance(testCases[i].config)
		rewardConfigCache := newRewardConfigCache(testGovernance)

		rewardConfig, error := rewardConfigCache.newRewardConfig(uint64(i))

		if error != nil {
			t.Errorf("error has occurred err : %v", error)
		}

		expectedResult := &testCases[i].result
		assert.Equal(t, expectedResult.blockNum, rewardConfig.blockNum)
		assert.Equal(t, expectedResult.mintingAmount, rewardConfig.mintingAmount)
		assert.Equal(t, expectedResult.cnRatio, rewardConfig.cnRatio)
		assert.Equal(t, expectedResult.pocRatio, rewardConfig.pocRatio)
		assert.Equal(t, expectedResult.kirRatio, rewardConfig.kirRatio)
		assert.Equal(t, expectedResult.totalRatio, rewardConfig.totalRatio)
		assert.Equal(t, expectedResult.unitPrice, rewardConfig.unitPrice)
	}
}

func TestRewardConfigCache_add(t *testing.T) {
	testCases := []rewardConfig{
		{
			blockNum:      1,
			mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
			cnRatio:       big.NewInt(0).SetInt64(34),
			pocRatio:      big.NewInt(0).SetInt64(54),
			kirRatio:      big.NewInt(0).SetInt64(12),
			totalRatio:    big.NewInt(0).SetInt64(100),
			unitPrice:     big.NewInt(0).SetInt64(25000000000),
		},
		{
			blockNum:      2,
			mintingAmount: big.NewInt(0).SetInt64(10000),
			cnRatio:       big.NewInt(0).SetInt64(50),
			pocRatio:      big.NewInt(0).SetInt64(30),
			kirRatio:      big.NewInt(0).SetInt64(20),
			totalRatio:    big.NewInt(0).SetInt64(100),
			unitPrice:     big.NewInt(0).SetInt64(50000000000),
		},
		{
			blockNum:      3,
			mintingAmount: big.NewInt(0).SetInt64(100000000),
			cnRatio:       big.NewInt(0).SetInt64(10),
			pocRatio:      big.NewInt(0).SetInt64(35),
			kirRatio:      big.NewInt(0).SetInt64(55),
			totalRatio:    big.NewInt(0).SetInt64(100),
			unitPrice:     big.NewInt(0).SetInt64(1500000000),
		},
	}

	testGovernance := newDefaultTestGovernance()
	rewardConfigCache := newRewardConfigCache(testGovernance)

	for i := 0; i < len(testCases); i++ {
		rewardConfigCache.add(uint64(i)+1, &testCases[i])
		assert.Equal(t, i+1, rewardConfigCache.cache.Len())
	}
}

func TestRewardConfigCache_add_sameNumber(t *testing.T) {
	rewardConfig := rewardConfig{
		blockNum:      1,
		mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
		cnRatio:       big.NewInt(0).SetInt64(34),
		pocRatio:      big.NewInt(0).SetInt64(54),
		kirRatio:      big.NewInt(0).SetInt64(12),
		totalRatio:    big.NewInt(0).SetInt64(100),
		unitPrice:     big.NewInt(0).SetInt64(25000000000),
	}

	testGovernance := newDefaultTestGovernance()
	rewardConfigCache := newRewardConfigCache(testGovernance)

	rewardConfigCache.add(1, &rewardConfig)
	rewardConfigCache.add(1, &rewardConfig)
	assert.Equal(t, 1, rewardConfigCache.cache.Len())
}

func TestRewardConfigCache_get_exist(t *testing.T) {
	testCases := []struct {
		blockNumber uint64
		config      map[int]interface{}
		result      rewardConfig
	}{
		{
			1,
			map[int]interface{}{
				params.Epoch:         604800,
				params.MintingAmount: "9600000000000000000",
				params.Ratio:         "34/54/12",
				params.UnitPrice:     25000000000,
				params.UseGiniCoeff:  true,
				params.DeferredTxFee: true,
			},
			rewardConfig{
				blockNum:      0,
				mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
				cnRatio:       big.NewInt(0).SetInt64(34),
				pocRatio:      big.NewInt(0).SetInt64(54),
				kirRatio:      big.NewInt(0).SetInt64(12),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(25000000000),
			},
		},
		{
			604805,
			map[int]interface{}{
				params.Epoch:         604800,
				params.MintingAmount: "9600000000000000",
				params.Ratio:         "40/25/35",
				params.UnitPrice:     250000000,
				params.UseGiniCoeff:  true,
				params.DeferredTxFee: false,
			},
			rewardConfig{
				blockNum:      604800,
				mintingAmount: big.NewInt(0).SetUint64(9600000000000000),
				cnRatio:       big.NewInt(0).SetInt64(40),
				pocRatio:      big.NewInt(0).SetInt64(25),
				kirRatio:      big.NewInt(0).SetInt64(35),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(250000000),
			},
		},
		{
			1210000,
			map[int]interface{}{
				params.Epoch:         604800,
				params.MintingAmount: "100000000000000000",
				params.Ratio:         "34/33/33",
				params.UnitPrice:     100000000000,
				params.UseGiniCoeff:  false,
				params.DeferredTxFee: true,
			},
			rewardConfig{
				blockNum:      1209600,
				mintingAmount: big.NewInt(0).SetUint64(100000000000000000),
				cnRatio:       big.NewInt(0).SetInt64(34),
				pocRatio:      big.NewInt(0).SetInt64(33),
				kirRatio:      big.NewInt(0).SetInt64(33),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(100000000000),
			},
		},
	}

	testGovernance := newDefaultTestGovernance()
	rewardConfigCache := newRewardConfigCache(testGovernance)

	for i := 0; i < len(testCases); i++ {
		blockNumber := testCases[i].blockNumber
		testGovernance.setTestGovernance(testCases[i].config)
		epoch := testGovernance.Params().Epoch()

		if blockNumber%epoch == 0 {
			blockNumber -= epoch
		} else {
			blockNumber -= (blockNumber % epoch)
		}
		rewardConfig, _ := rewardConfigCache.newRewardConfig(blockNumber)
		rewardConfigCache.add(blockNumber, rewardConfig)
	}
	for i := 0; i < len(testCases); i++ {
		rewardConfig, err := rewardConfigCache.get(testCases[i].blockNumber)
		if err != nil {
			t.Errorf("error has occurred. err : %v", err)
		}
		assert.Equal(t, testCases[i].result.blockNum, rewardConfig.blockNum)
		assert.Equal(t, testCases[i].result.mintingAmount, rewardConfig.mintingAmount)
		assert.Equal(t, testCases[i].result.cnRatio, rewardConfig.cnRatio)
		assert.Equal(t, testCases[i].result.pocRatio, rewardConfig.pocRatio)
		assert.Equal(t, testCases[i].result.kirRatio, rewardConfig.kirRatio)
		assert.Equal(t, testCases[i].result.totalRatio, rewardConfig.totalRatio)
		assert.Equal(t, testCases[i].result.unitPrice, rewardConfig.unitPrice)
	}
}
