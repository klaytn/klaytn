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

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

type testBalanceAdder struct {
	accounts map[common.Address]*big.Int
}

func newTestBalanceAdder() *testBalanceAdder {
	balanceAdder := &testBalanceAdder{}
	balanceAdder.accounts = make(map[common.Address]*big.Int)
	return balanceAdder
}

func (balanceAdder *testBalanceAdder) AddBalance(addr common.Address, v *big.Int) {
	balance, ok := balanceAdder.accounts[addr]
	if ok {
		balanceAdder.accounts[addr] = big.NewInt(0).Add(balance, v)
	} else {
		balanceAdder.accounts[addr] = v
	}
}

func (balanceAdder *testBalanceAdder) GetBalance(addr common.Address) *big.Int {
	balance, ok := balanceAdder.accounts[addr]
	if ok {
		return balance
	} else {
		return nil
	}
}

func Test_isEmptyAddress(t *testing.T) {
	testCases := []struct {
		address common.Address
		result  bool
	}{
		{
			common.Address{},
			true,
		},
		{
			common.HexToAddress("0x0000000000000000000000000000000000000000"),
			true,
		},
		{
			common.StringToAddress("0xA75Ed91f789BF9dc121DACB822849955ca3AB6aD"),
			false,
		},
		{
			common.StringToAddress("0x4bCDd8E3F9776d16056815E189EcB5A8bF8E4CBb"),
			false,
		},
	}
	for _, testCase := range testCases {
		assert.Equal(t, testCase.result, common.EmptyAddress(testCase.address))
	}
}

func TestRewardDistributor_getTotalTxFee(t *testing.T) {
	testCases := []struct {
		gasUsed            uint64
		unitPrice          uint64
		expectedTotalTxFee *big.Int
	}{
		{0, 25000000000, big.NewInt(0)},
		{200000, 25000000000, big.NewInt(5000000000000000)},
		{129346, 10000000000, big.NewInt(1293460000000000)},
		{9236192, 50000, big.NewInt(461809600000)},
		{12936418927364923, 0, big.NewInt(0)},
	}
	rewardDistributor := NewRewardDistributor(newDefaultTestGovernance())
	rewardConfig := &rewardConfig{}

	header := &types.Header{}
	unitPrice := big.NewInt(0)

	for _, testCase := range testCases {
		header.GasUsed = testCase.gasUsed
		rewardConfig.unitPrice = unitPrice.SetUint64(testCase.unitPrice)

		result := rewardDistributor.getTotalTxFee(header, rewardConfig)

		assert.Equal(t, testCase.expectedTotalTxFee.Uint64(), result.Uint64())
	}
}

func TestRewardDistributor_MintKLAY(t *testing.T) {
	BalanceAdder := newTestBalanceAdder()
	header := &types.Header{}
	header.Number = big.NewInt(0)
	header.Rewardbase = common.StringToAddress("0x1552F52D459B713E0C4558e66C8c773a75615FA8")
	governance := newDefaultTestGovernance()
	rewardDistributor := NewRewardDistributor(governance)

	err := rewardDistributor.MintKLAY(BalanceAdder, header)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, BalanceAdder.GetBalance(header.Rewardbase).Int64())
	assert.Equal(t, governance.mintingAmount, BalanceAdder.GetBalance(header.Rewardbase).String())
}

func TestRewardDistributor_distributeBlockReward(t *testing.T) {
	testCases := []struct {
		totalTxFee         *big.Int
		rewardConfig       *rewardConfig
		expectedCnBalance  *big.Int
		expectedPocBalance *big.Int
		expectedKirBalance *big.Int
	}{
		{
			totalTxFee: big.NewInt(0),
			rewardConfig: &rewardConfig{
				blockNum:      1,
				mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
				cnRatio:       big.NewInt(0).SetInt64(34),
				pocRatio:      big.NewInt(0).SetInt64(54),
				kirRatio:      big.NewInt(0).SetInt64(12),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(25000000000),
			},
			expectedCnBalance:  big.NewInt(0).SetUint64(3264000000000000000),
			expectedPocBalance: big.NewInt(0).SetUint64(5184000000000000000),
			expectedKirBalance: big.NewInt(0).SetUint64(1152000000000000000),
		},
		{
			totalTxFee: big.NewInt(1000000),
			rewardConfig: &rewardConfig{
				blockNum:      1,
				mintingAmount: big.NewInt(0).SetUint64(10000000000),
				cnRatio:       big.NewInt(0).SetInt64(60),
				pocRatio:      big.NewInt(0).SetInt64(30),
				kirRatio:      big.NewInt(0).SetInt64(10),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(25000000000),
			},
			expectedCnBalance:  big.NewInt(0).SetUint64(6000600000),
			expectedPocBalance: big.NewInt(0).SetUint64(3000300000),
			expectedKirBalance: big.NewInt(0).SetUint64(1000100000),
		},
	}

	header := &types.Header{}
	header.Number = big.NewInt(0)
	header.Rewardbase = common.StringToAddress("0x1552F52D459B713E0C4558e66C8c773a75615FA8")
	pocAddress := common.StringToAddress("0x4bCDd8E3F9776d16056815E189EcB5A8bF8E4CBb")
	kirAddress := common.StringToAddress("0xd38A08AD21B44681f5e75D0a3CA4793f3E6c03e7")
	governance := newDefaultTestGovernance()

	for _, testCase := range testCases {
		BalanceAdder := newTestBalanceAdder()
		rewardDistributor := NewRewardDistributor(governance)
		rewardDistributor.distributeBlockReward(BalanceAdder, header, testCase.totalTxFee, testCase.rewardConfig, pocAddress, kirAddress)

		assert.Equal(t, testCase.expectedCnBalance.Uint64(), BalanceAdder.GetBalance(header.Rewardbase).Uint64())
		assert.Equal(t, testCase.expectedPocBalance.Uint64(), BalanceAdder.GetBalance(pocAddress).Uint64())
		assert.Equal(t, testCase.expectedKirBalance.Uint64(), BalanceAdder.GetBalance(kirAddress).Uint64())
	}
}

func TestRewardDistributor_DistributeBlockReward(t *testing.T) {
	testCases := []struct {
		gasUsed            uint64
		epoch              uint64
		mintingAmount      string
		ratio              string
		unitprice          uint64
		useGiniCoeff       bool
		deferredTxFee      bool
		expectedCnBalance  *big.Int
		expectedPocBalance *big.Int
		expectedKirBalance *big.Int
	}{
		{
			gasUsed:            100,
			epoch:              30,
			mintingAmount:      "50000",
			ratio:              "40/50/10",
			unitprice:          500,
			useGiniCoeff:       true,
			deferredTxFee:      true,
			expectedCnBalance:  big.NewInt(0).SetUint64(40000),
			expectedPocBalance: big.NewInt(0).SetUint64(50000),
			expectedKirBalance: big.NewInt(0).SetUint64(10000),
		},
		{
			gasUsed:            0,
			epoch:              604800,
			mintingAmount:      "9600000000000000000",
			ratio:              "34/54/12",
			unitprice:          25000000000,
			useGiniCoeff:       true,
			deferredTxFee:      true,
			expectedCnBalance:  big.NewInt(0).SetUint64(3264000000000000000),
			expectedPocBalance: big.NewInt(0).SetUint64(5184000000000000000),
			expectedKirBalance: big.NewInt(0).SetUint64(1152000000000000000),
		},
		{
			gasUsed:            0,
			epoch:              3600,
			mintingAmount:      "0",
			ratio:              "100/0/0",
			unitprice:          0,
			useGiniCoeff:       true,
			deferredTxFee:      true,
			expectedCnBalance:  big.NewInt(0).SetUint64(0),
			expectedPocBalance: big.NewInt(0).SetUint64(0),
			expectedKirBalance: big.NewInt(0).SetUint64(0),
		},
	}

	header := &types.Header{}
	header.Number = big.NewInt(0)
	header.Rewardbase = common.StringToAddress("0x1552F52D459B713E0C4558e66C8c773a75615FA8")
	pocAddress := common.StringToAddress("0x4bCDd8E3F9776d16056815E189EcB5A8bF8E4CBb")
	kirAddress := common.StringToAddress("0xd38A08AD21B44681f5e75D0a3CA4793f3E6c03e7")
	governance := newDefaultTestGovernance()

	for _, testCase := range testCases {
		BalanceAdder := newTestBalanceAdder()
		governance.setTestGovernance(testCase.epoch, testCase.mintingAmount, testCase.ratio, testCase.unitprice, testCase.useGiniCoeff, testCase.deferredTxFee)
		header.GasUsed = testCase.gasUsed
		rewardDistributor := NewRewardDistributor(governance)

		err := rewardDistributor.DistributeBlockReward(BalanceAdder, header, pocAddress, kirAddress)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.NotNil(t, BalanceAdder.GetBalance(header.Rewardbase).Int64())
		assert.Equal(t, testCase.expectedCnBalance.Uint64(), BalanceAdder.GetBalance(header.Rewardbase).Uint64())
		assert.Equal(t, testCase.expectedPocBalance.Uint64(), BalanceAdder.GetBalance(pocAddress).Uint64())
		assert.Equal(t, testCase.expectedKirBalance.Uint64(), BalanceAdder.GetBalance(kirAddress).Uint64())
	}
}
