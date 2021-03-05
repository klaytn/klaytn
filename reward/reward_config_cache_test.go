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
	"errors"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

type testGovernance struct {
	epoch           uint64
	mintingAmount   string
	ratio           string
	unitPrice       uint64
	useGiniCoeff    bool
	policy          uint64
	stakingInterval uint64
	deferredTxFee   bool
}

func newDefaultTestGovernance() *testGovernance {
	return &testGovernance{
		epoch:           604800,
		mintingAmount:   "9600000000000000000",
		ratio:           "34/54/12",
		unitPrice:       25000000000,
		useGiniCoeff:    true,
		policy:          params.WeightedRandom,
		stakingInterval: 86400,
		deferredTxFee:   true,
	}
}

func newTestGovernance(epoch uint64, mintingAmount string, ratio string, unitPrice uint64, useGiniCoeff bool, stakingInterval uint64, deferredTxFee bool) *testGovernance {
	return &testGovernance{
		epoch:           epoch,
		mintingAmount:   mintingAmount,
		ratio:           ratio,
		unitPrice:       unitPrice,
		useGiniCoeff:    useGiniCoeff,
		stakingInterval: stakingInterval,
		deferredTxFee:   deferredTxFee,
	}
}

func (governance *testGovernance) Epoch() uint64 {
	return governance.epoch
}

func (governance *testGovernance) GetItemAtNumberByIntKey(num uint64, key int) (interface{}, error) {
	switch key {
	case params.MintingAmount:
		return governance.mintingAmount, nil
	case params.Ratio:
		return governance.ratio, nil
	case params.UnitPrice:
		return governance.unitPrice, nil
	case params.Epoch:
		return governance.epoch, nil
	default:
		return nil, errors.New("Unhandled key on testGovernance")
	}
}

func (governance *testGovernance) ProposerPolicy() uint64 {
	return governance.policy
}

func (governance *testGovernance) DeferredTxFee() bool {
	return governance.deferredTxFee
}

func (governance *testGovernance) StakingUpdateInterval() uint64 {
	return governance.stakingInterval
}

func (governance *testGovernance) setTestGovernance(epoch uint64, mintingAmount string, ratio string, unitprice uint64, useGiniCoeff bool, deferredTxFee bool) {
	governance.epoch = epoch
	governance.mintingAmount = mintingAmount
	governance.ratio = ratio
	governance.unitPrice = unitprice
	governance.useGiniCoeff = useGiniCoeff
	governance.deferredTxFee = deferredTxFee
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
		testGovernance testGovernance
		result         rewardConfig
	}{
		{
			testGovernance{
				epoch:         604800,
				mintingAmount: "9600000000000000000",
				ratio:         "34/54/12",
				unitPrice:     25000000000,
				useGiniCoeff:  true,
				deferredTxFee: true,
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
			testGovernance{
				epoch:         30,
				mintingAmount: "10000",
				ratio:         "50/30/20",
				unitPrice:     50000000000,
				useGiniCoeff:  true,
				deferredTxFee: false,
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
			testGovernance{
				epoch:         3000,
				mintingAmount: "100000000",
				ratio:         "10/35/55",
				unitPrice:     1500000000,
				useGiniCoeff:  false,
				deferredTxFee: true,
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
		testGovernance := &testCases[i].testGovernance
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
		blockNumber   uint64
		epoch         uint64
		mintingAmount string
		ratio         string
		unitprice     uint64
		useGiniCoeff  bool
		deferredTxFee bool
		result        rewardConfig
	}{
		{
			blockNumber:   1,
			epoch:         604800,
			mintingAmount: "9600000000000000000",
			ratio:         "34/54/12",
			unitprice:     25000000000,
			useGiniCoeff:  true,
			deferredTxFee: true,
			result: rewardConfig{
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
			blockNumber:   604805,
			epoch:         604800,
			mintingAmount: "9600000000000000",
			ratio:         "40/25/35",
			unitprice:     250000000,
			useGiniCoeff:  true,
			deferredTxFee: false,
			result: rewardConfig{
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
			blockNumber:   1210000,
			epoch:         604800,
			mintingAmount: "100000000000000000",
			ratio:         "34/33/33",
			unitprice:     100000000000,
			useGiniCoeff:  false,
			deferredTxFee: true,
			result: rewardConfig{
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
		testGovernance.setTestGovernance(testCases[i].epoch, testCases[i].mintingAmount, testCases[i].ratio, testCases[i].unitprice, testCases[i].useGiniCoeff, testCases[i].deferredTxFee)
		blockNumber := testCases[i].blockNumber
		if blockNumber%testCases[i].epoch == 0 {
			blockNumber -= testCases[i].epoch
		} else {
			blockNumber -= (blockNumber % testCases[i].epoch)
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
