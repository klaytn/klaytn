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
	"fmt"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

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

var (
	cnBaseAddr     = 500 // Dummy addresses goes like 0x000..5nn
	stakeBaseAddr  = 600
	rewardBaseAddr = 700
	minStaking     = uint64(2000000) // changing this value will not change the governance's min staking
	minted, _      = big.NewInt(0).SetString("9600000000000000000", 10)
	proposerAddr   = common.StringToAddress("0x1552F52D459B713E0C4558e66C8c773a75615FA8")
	kirAddr        = intToAddress(1000)
	kgfAddr        = intToAddress(2000)
)

// 500 -> 0x00000...0500
func intToAddress(x int) common.Address {
	return common.HexToAddress(fmt.Sprintf("0x%040d", x))
}

// rewardOverride[i] = j means rewards[i] = rewards[j]
// amountOverride[i] = j means amount[i] = j
func genStakingInfo(cnNum int, rewardOverride map[int]int, amountOverride map[int]uint64) *StakingInfo {
	cns := make([]common.Address, 0)
	stakes := make([]common.Address, 0)
	rewards := make([]common.Address, 0)
	amounts := make([]uint64, 0)

	for i := 0; i < cnNum; i++ {
		cns = append(cns, intToAddress(cnBaseAddr+i))
		stakes = append(stakes, intToAddress(stakeBaseAddr+i))
		rewards = append(rewards, intToAddress(rewardBaseAddr+i))
		amounts = append(amounts, minStaking)
	}

	for i := range rewardOverride {
		rewards[i] = rewards[rewardOverride[i]]
	}

	for i := range amountOverride {
		amounts[i] = amountOverride[i]
	}

	return &StakingInfo{
		BlockNum:              0,
		CouncilNodeAddrs:      cns,
		CouncilStakingAddrs:   stakes,
		CouncilRewardAddrs:    rewards,
		KIRAddr:               kirAddr,
		PoCAddr:               kgfAddr,
		UseGini:               false,
		CouncilStakingAmounts: amounts,
	}
}

type testBalanceAdder struct {
	accounts map[common.Address]*big.Int
}

func newTestBalanceAdder() *testBalanceAdder {
	balanceAdder := &testBalanceAdder{}
	balanceAdder.accounts = make(map[common.Address]*big.Int)
	return balanceAdder
}

func getTestConfig() *params.ChainConfig {
	return &params.ChainConfig{
		MagmaCompatibleBlock: big.NewInt(0),
		KoreCompatibleBlock:  big.NewInt(0),
		UnitPrice:            1,
		Governance: &params.GovernanceConfig{
			Reward: &params.RewardConfig{
				MintingAmount: minted,
				Ratio:         "34/54/12",
				Kip82Ratio:    "20/80",
				DeferredTxFee: true,
				MinimumStake:  big.NewInt(0).SetUint64(minStaking),
			},
		},
		Istanbul: &params.IstanbulConfig{
			ProposerPolicy: 2,
		},
	}
}

func (balanceAdder *testBalanceAdder) AddBalance(addr common.Address, v *big.Int) {
	balance, ok := balanceAdder.accounts[addr]
	if ok {
		balanceAdder.accounts[addr] = big.NewInt(0).Add(balance, v)
	} else {
		balanceAdder.accounts[addr] = v
	}
}

func noMagma(p *params.ChainConfig) *params.ChainConfig {
	p.MagmaCompatibleBlock = big.NewInt(100000000)
	p.KoreCompatibleBlock = big.NewInt(100000000)
	return p
}

func noKore(p *params.ChainConfig) *params.ChainConfig {
	p.KoreCompatibleBlock = big.NewInt(100000000)
	return p
}

func noDeferred(p *params.ChainConfig) *params.ChainConfig {
	p.Governance.Reward.DeferredTxFee = false
	return p
}

func roundrobin(p *params.ChainConfig) *params.ChainConfig {
	p.Istanbul.ProposerPolicy = 0
	return p
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

func TestRewardDistributor_GetTotalTxFee(t *testing.T) {
	testCases := []struct {
		gasUsed            uint64
		unitPrice          uint64
		baseFee            *big.Int
		expectedTotalTxFee *big.Int
	}{
		// before magma hardfork
		{0, 25000000000, nil, big.NewInt(0)},
		{200000, 25000000000, nil, big.NewInt(5000000000000000)},
		{200000, 25000000000, nil, big.NewInt(5000000000000000)},
		{129346, 10000000000, nil, big.NewInt(1293460000000000)},
		{129346, 10000000000, nil, big.NewInt(1293460000000000)},
		{9236192, 50000, nil, big.NewInt(461809600000)},
		{9236192, 50000, nil, big.NewInt(461809600000)},
		{12936418927364923, 0, nil, big.NewInt(0)},
		// after magma hardfork, unitprice ignored
		{0, 25000000000, big.NewInt(25000000000), big.NewInt(0)},
		{200000, 25000000000, big.NewInt(25000000000), big.NewInt(5000000000000000)},
		{200000, 25000000000, big.NewInt(25000000000), big.NewInt(5000000000000000)},
		{129346, 25000000000, big.NewInt(10000000000), big.NewInt(1293460000000000)},
		{129346, 250, big.NewInt(10000000000), big.NewInt(1293460000000000)},
		{9236192, 9876, big.NewInt(50000), big.NewInt(461809600000)},
		{9236192, 25000000000, big.NewInt(50000), big.NewInt(461809600000)},
		{12936418927364923, 25000000000, big.NewInt(0), big.NewInt(0)},
	}

	for _, testCase := range testCases {
		header := &types.Header{
			Number:  big.NewInt(0),
			GasUsed: testCase.gasUsed,
			BaseFee: testCase.baseFee,
		}
		config := &params.ChainConfig{
			UnitPrice: testCase.unitPrice,
		}
		if testCase.baseFee != nil {
			// enable Magma
			config.MagmaCompatibleBlock = big.NewInt(0)
		}

		result := GetTotalTxFee(header, config)
		assert.Equal(t, testCase.expectedTotalTxFee.Uint64(), result.Uint64())
	}
}

func TestRewardDistributor_getBurnAmountMagma(t *testing.T) {
	testCases := []struct {
		gasUsed            uint64
		baseFee            *big.Int
		expectedTotalTxFee *big.Int
	}{
		{0, big.NewInt(25000000000), big.NewInt(0)},
		{200000, big.NewInt(25000000000), big.NewInt(5000000000000000 / 2)},
		{200000, big.NewInt(25000000000), big.NewInt(5000000000000000 / 2)},
		{129346, big.NewInt(10000000000), big.NewInt(1293460000000000 / 2)},
		{129346, big.NewInt(10000000000), big.NewInt(1293460000000000 / 2)},
		{9236192, big.NewInt(50000), big.NewInt(461809600000 / 2)},
		{9236192, big.NewInt(50000), big.NewInt(461809600000 / 2)},
		{12936418927364923, big.NewInt(0), big.NewInt(0)},
	}

	header := &types.Header{
		Number: big.NewInt(1),
	}
	config := &params.ChainConfig{}
	config.MagmaCompatibleBlock = big.NewInt(0)

	for _, testCase := range testCases {
		header.GasUsed = testCase.gasUsed
		header.BaseFee = testCase.baseFee

		txFee := GetTotalTxFee(header, config)
		burnedTxFee := getBurnAmountMagma(txFee)
		assert.Equal(t, testCase.expectedTotalTxFee.Uint64(), burnedTxFee.Uint64())
	}
}

func TestRewardDistributor_GetBlockReward(t *testing.T) {
	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	var (
		header = &types.Header{
			Number:     big.NewInt(1),
			GasUsed:    1000,
			BaseFee:    big.NewInt(1),
			Rewardbase: proposerAddr,
		}

		stakingInfo = genStakingInfo(5, nil, map[int]uint64{
			0: minStaking + 4,
			1: minStaking + 3,
		})
	)

	testcases := []struct {
		policy        istanbul.ProposerPolicy
		deferredTxFee bool
		expected      *RewardSpec
	}{
		{
			policy:        istanbul.RoundRobin,
			deferredTxFee: false,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(1000),
				BurntFee: new(big.Int).SetUint64(500),
				Proposer: new(big.Int).SetUint64(9.6e18 + 500),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: new(big.Int).SetUint64(9.6e18 + 500),
				},
			},
		},
		{
			policy:        istanbul.WeightedRandom,
			deferredTxFee: true,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(1000),
				BurntFee: new(big.Int).SetUint64(1000),
				Proposer: new(big.Int).SetUint64(0.6528e18 + 1),
				Stakers:  new(big.Int).SetUint64(2.6112e18),
				Kgf:      new(big.Int).SetUint64(5.184e18),
				Kir:      new(big.Int).SetUint64(1.152e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr:                     new(big.Int).SetUint64(0.6528e18 + 1),
					kgfAddr:                          new(big.Int).SetUint64(5.184e18),
					kirAddr:                          new(big.Int).SetUint64(1.152e18),
					intToAddress(rewardBaseAddr):     new(big.Int).SetUint64(1492114285714285714),
					intToAddress(rewardBaseAddr + 1): new(big.Int).SetUint64(1119085714285714285),
				},
			},
		},
		{
			policy:        istanbul.WeightedRandom,
			deferredTxFee: false,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(1000),
				BurntFee: new(big.Int).SetUint64(0),
				Proposer: new(big.Int).SetUint64(0.6528e18 + 1),
				Stakers:  new(big.Int).SetUint64(2.6112e18),
				Kgf:      new(big.Int).SetUint64(5.184e18),
				Kir:      new(big.Int).SetUint64(1.152e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr:                     new(big.Int).SetUint64(0.6528e18 + 1000 + 1),
					kgfAddr:                          new(big.Int).SetUint64(5.184e18),
					kirAddr:                          new(big.Int).SetUint64(1.152e18),
					intToAddress(rewardBaseAddr):     new(big.Int).SetUint64(1492114285714285714),
					intToAddress(rewardBaseAddr + 1): new(big.Int).SetUint64(1119085714285714285),
				},
			},
		},
	}

	SetTestStakingManagerWithStakingInfoCache(stakingInfo)

	for i, tc := range testcases {
		config := getTestConfig()
		if !tc.deferredTxFee {
			config = noDeferred(config)
		}
		config.Istanbul.ProposerPolicy = uint64(tc.policy)
		spec, err := GetBlockReward(header, config)
		assert.Nil(t, err, "testcases[%d] failed", i)
		assert.Equal(t, tc.expected, spec, "testcases[%d] failed", i)
	}
}

func TestRewardDistributor_CalcDeferredRewardSimple(t *testing.T) {
	header := &types.Header{
		Number:     big.NewInt(1),
		GasUsed:    1000,
		BaseFee:    big.NewInt(1),
		Rewardbase: proposerAddr,
	}

	testcases := []struct {
		isMagma  bool
		expected *RewardSpec
	}{
		{
			isMagma: false,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(1000),
				BurntFee: new(big.Int).SetUint64(0),
				Proposer: new(big.Int).SetUint64(9.6e18 + 1000),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: new(big.Int).SetUint64(9.6e18 + 1000),
				},
			},
		},
		{
			isMagma: true,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(1000),
				BurntFee: new(big.Int).SetUint64(500), // 50% of tx fee burnt
				Proposer: new(big.Int).SetUint64(9.6e18 + 500),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: new(big.Int).SetUint64(9.6e18 + 500),
				},
			},
		},
	}

	for i, tc := range testcases {
		config := roundrobin(getTestConfig())
		if !tc.isMagma {
			config = noMagma(config)
		}

		spec, err := CalcDeferredRewardSimple(header, config)
		assert.Nil(t, err, "testcases[%d] failed", i)
		assert.Equal(t, tc.expected, spec, "testcases[%d] failed", i)
	}
}

func TestRewardDistributor_CalcDeferredReward(t *testing.T) {
	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	stakingInfo := genStakingInfo(5, nil, map[int]uint64{
		0: minStaking + 4,
		1: minStaking + 3,
	})

	testcases := []struct {
		desc     string
		isKore   bool
		isMagma  bool
		fee      uint64
		expected *RewardSpec
	}{
		{
			desc:    "isKore=false, isMagma=false, fee=1000 [000]",
			isKore:  false,
			isMagma: false,
			fee:     1000,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: big.NewInt(1000),
				BurntFee: big.NewInt(0),
				Proposer: big.NewInt(0).SetUint64(3.264e18 + 340),
				Stakers:  big.NewInt(0),
				Kgf:      big.NewInt(5.184e18 + 540),
				Kir:      big.NewInt(1.152e18 + 120),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: big.NewInt(3.264e18 + 340),
					kgfAddr:      big.NewInt(5.184e18 + 540),
					kirAddr:      big.NewInt(1.152e18 + 120),
				},
			},
		},
		{
			desc:    "isKore=false, isMagma=false, fee=10e18 [001]",
			isKore:  false,
			isMagma: false,
			fee:     10e18,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(10e18),
				BurntFee: new(big.Int).SetUint64(0),
				Proposer: new(big.Int).SetUint64(6.664e18),
				Stakers:  new(big.Int).SetUint64(0),
				Kgf:      new(big.Int).SetUint64(10.584e18),
				Kir:      new(big.Int).SetUint64(2.352e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: new(big.Int).SetUint64(6.664e18),
					kgfAddr:      new(big.Int).SetUint64(10.584e18),
					kirAddr:      new(big.Int).SetUint64(2.352e18),
				},
			},
		},
		{
			desc:    "isKore=false, isMagma=true, fee=1000 [010]",
			isKore:  false,
			isMagma: true,
			fee:     1000,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(1000),
				BurntFee: new(big.Int).SetUint64(500),
				Proposer: new(big.Int).SetUint64(3.264e18 + 170),
				Stakers:  new(big.Int).SetUint64(0),
				Kgf:      new(big.Int).SetUint64(5.184e18 + 270),
				Kir:      new(big.Int).SetUint64(1.152e18 + 60),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: new(big.Int).SetUint64(3.264e18 + 170),
					kirAddr:      new(big.Int).SetUint64(1.152e18 + 60),
					kgfAddr:      new(big.Int).SetUint64(5.184e18 + 270),
				},
			},
		},
		{
			desc:    "isKore=false, isMagma=true, fee=10e18 [011]",
			isKore:  false,
			isMagma: true,
			fee:     10e18,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(10e18),
				BurntFee: new(big.Int).SetUint64(5e18),
				Proposer: new(big.Int).SetUint64(4.964e18),
				Stakers:  new(big.Int).SetUint64(0),
				Kgf:      new(big.Int).SetUint64(7.884e18),
				Kir:      new(big.Int).SetUint64(1.752e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: new(big.Int).SetUint64(4.964e18),
					kgfAddr:      new(big.Int).SetUint64(7.884e18),
					kirAddr:      new(big.Int).SetUint64(1.752e18),
				},
			},
		},
		{
			desc:    "isKore=true, isMagma=true, fee=1000 [110]",
			isKore:  true,
			isMagma: true,
			fee:     1000,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(1000),
				BurntFee: new(big.Int).SetUint64(1000),
				Proposer: new(big.Int).SetUint64(0.6528e18 + 1),
				Stakers:  new(big.Int).SetUint64(2.6112e18),
				Kgf:      new(big.Int).SetUint64(5.184e18),
				Kir:      new(big.Int).SetUint64(1.152e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr:                     new(big.Int).SetUint64(0.6528e18 + 1),
					kgfAddr:                          new(big.Int).SetUint64(5.184e18),
					kirAddr:                          new(big.Int).SetUint64(1.152e18),
					intToAddress(rewardBaseAddr):     new(big.Int).SetUint64(1492114285714285714),
					intToAddress(rewardBaseAddr + 1): new(big.Int).SetUint64(1119085714285714285),
				},
			},
		},
		{ // after kore, more-than-default staking, large fee, proposer = rewardbase
			desc:    "isKore=true, isMagma=true, fee=10e18 [111]",
			isKore:  true,
			isMagma: true,
			fee:     10e18,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: new(big.Int).SetUint64(10e18),
				BurntFee: new(big.Int).SetUint64(5e18 + 0.6528e18),
				Proposer: new(big.Int).SetUint64(5e18 + 1),
				Stakers:  new(big.Int).SetUint64(2.6112e18),
				Kgf:      new(big.Int).SetUint64(5.184e18),
				Kir:      new(big.Int).SetUint64(1.152e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr:                     new(big.Int).SetUint64(5e18 + 1),
					kgfAddr:                          new(big.Int).SetUint64(5.184e18),
					kirAddr:                          new(big.Int).SetUint64(1.152e18),
					intToAddress(rewardBaseAddr):     new(big.Int).SetUint64(1492114285714285714),
					intToAddress(rewardBaseAddr + 1): new(big.Int).SetUint64(1119085714285714285),
				},
			},
		},
	}

	SetTestStakingManagerWithStakingInfoCache(stakingInfo)

	for _, tc := range testcases {
		header := &types.Header{
			Number:     big.NewInt(1),
			GasUsed:    tc.fee,
			BaseFee:    big.NewInt(1),
			Rewardbase: proposerAddr,
		}

		config := getTestConfig()
		if !tc.isKore {
			config = noKore(config)
		}
		if !tc.isMagma {
			config = noMagma(config)
		}

		spec, err := CalcDeferredReward(header, config)
		assert.Nil(t, err, "failed tc: %s", tc.desc)
		assert.Equal(t, tc.expected, spec, "failed tc: %s", tc.desc)
	}
}

func TestRewardDistributor_CalcDeferredReward_StakingInfos(t *testing.T) {
	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	var (
		header = &types.Header{
			Number:     big.NewInt(1),
			GasUsed:    1000,
			BaseFee:    big.NewInt(1),
			Rewardbase: proposerAddr,
		}
		config = getTestConfig()
	)

	testcases := []struct {
		desc        string
		stakingInfo *StakingInfo
		expected    *RewardSpec
	}{
		{
			desc:        "stakingInfo is nil, its portion goes to proposer",
			stakingInfo: nil,
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: big.NewInt(1000),
				BurntFee: big.NewInt(1000),
				Proposer: minted,
				Stakers:  big.NewInt(2.6112e18),
				Kgf:      big.NewInt(0),
				Kir:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: minted,
				},
			},
		},
		{
			desc: "stakingInfo has no poc, its portion goes to proposer",
			stakingInfo: &StakingInfo{
				KIRAddr: kirAddr,
				PoCAddr: common.Address{},
			},
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: big.NewInt(1000),
				BurntFee: big.NewInt(1000),
				Proposer: big.NewInt(8.448e18),
				Stakers:  big.NewInt(2.6112e18),
				Kgf:      big.NewInt(0),
				Kir:      big.NewInt(1.152e18), // minted * 0.12
				Rewards: map[common.Address]*big.Int{
					proposerAddr: big.NewInt(8.448e18),
					kirAddr:      big.NewInt(1.152e18),
				},
			},
		},
		{
			desc: "stakingInfo has no kir, its portion goes to proposer",
			stakingInfo: &StakingInfo{
				KIRAddr: common.Address{},
				PoCAddr: kgfAddr,
			},
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: big.NewInt(1000),
				BurntFee: big.NewInt(1000),
				Proposer: big.NewInt(4.416e18),
				Stakers:  big.NewInt(2.6112e18),
				Kgf:      big.NewInt(5.184e18), // minted * 0.54
				Kir:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: big.NewInt(4.416e18),
					kgfAddr:      big.NewInt(5.184e18),
				},
			},
		},
		{
			desc: "stakingInfo has the same poc and kir",
			stakingInfo: &StakingInfo{
				KIRAddr: kgfAddr,
				PoCAddr: kgfAddr,
			},
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: big.NewInt(1000),
				BurntFee: big.NewInt(1000),
				Proposer: big.NewInt(3.264e18),
				Stakers:  big.NewInt(2.6112e18),
				Kgf:      big.NewInt(5.184e18),
				Kir:      big.NewInt(1.152e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr: big.NewInt(3.264e18),
					kgfAddr:      big.NewInt(6.336e18),
				},
			},
		},
	}

	for i, tc := range testcases {
		if tc.stakingInfo == nil {
			SetTestStakingManager(nil)
		} else {
			SetTestStakingManagerWithStakingInfoCache(tc.stakingInfo)
		}
		spec, err := CalcDeferredReward(header, config)
		assert.Nil(t, err, "testcases[%d] failed", i)
		assert.Equal(t, tc.expected, spec, "testcases[%d] failed: %s", i, tc.desc)
	}
}

func TestRewardDistributor_CalcDeferredReward_Remainings(t *testing.T) {
	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	var (
		header = &types.Header{
			Number:     big.NewInt(1),
			GasUsed:    1000,
			BaseFee:    big.NewInt(1),
			Rewardbase: proposerAddr,
		}

		stakingInfo = genStakingInfo(5, nil, map[int]uint64{
			0: minStaking + 4,
			1: minStaking + 3,
		})
		splitRemainingConfig = getTestConfig()
	)
	splitRemainingConfig.Governance.Reward.MintingAmount = big.NewInt(333)

	testcases := []struct {
		desc     string
		config   *params.ChainConfig
		expected *RewardSpec
	}{
		{
			desc:   "split remaining goes to kgf",
			config: splitRemainingConfig,
			expected: &RewardSpec{
				Minted:   big.NewInt(333),
				TotalFee: big.NewInt(1000),
				BurntFee: big.NewInt(522),
				Proposer: big.NewInt(501), // proposer=22, rewardFee=478, shareRem=1
				Stakers:  big.NewInt(90),
				Kgf:      big.NewInt(182), // splitRem=3
				Kir:      big.NewInt(39),
				Rewards: map[common.Address]*big.Int{
					proposerAddr:                     big.NewInt(501),
					kgfAddr:                          big.NewInt(182),
					kirAddr:                          big.NewInt(39),
					intToAddress(rewardBaseAddr):     big.NewInt(51), // stakers * 4/7
					intToAddress(rewardBaseAddr + 1): big.NewInt(38), // stakers * 3/7
				},
			},
		},
		{
			desc:   "share remaining goes to proposer",
			config: getTestConfig(),
			expected: &RewardSpec{
				Minted:   minted,
				TotalFee: big.NewInt(1000),
				BurntFee: big.NewInt(1000),
				Proposer: big.NewInt(0.6528e18 + 1),
				Stakers:  big.NewInt(2.6112e18),
				Kgf:      big.NewInt(5.184e18),
				Kir:      big.NewInt(1.152e18),
				Rewards: map[common.Address]*big.Int{
					proposerAddr:                     big.NewInt(0.6528e18 + 1),
					kgfAddr:                          big.NewInt(5.184e18),
					kirAddr:                          big.NewInt(1.152e18),
					intToAddress(rewardBaseAddr):     big.NewInt(1492114285714285714), // stakers * 4/7
					intToAddress(rewardBaseAddr + 1): big.NewInt(1119085714285714285), // stakers * 3/7
				},
			},
		},
	}

	SetTestStakingManagerWithStakingInfoCache(stakingInfo)

	for _, tc := range testcases {
		spec, err := CalcDeferredReward(header, tc.config)
		assert.Nil(t, err, "failed tc: %s", tc.desc)
		assert.Equal(t, tc.expected, spec, "failed tc: %s", tc.desc)
	}
}

func TestRewardDistributor_calcDeferredFee(t *testing.T) {
	type Result struct{ total, reward, burnt uint64 }

	testcases := []struct {
		desc     string
		isKore   bool
		isMagma  bool
		fee      uint64
		expected *Result
	}{
		{
			desc:    "isKore=false, isMagma=false, fee=1000 [000]",
			isKore:  false,
			isMagma: false,
			fee:     1000,
			expected: &Result{
				total:  1000,
				reward: 1000,
				burnt:  0,
			},
		},
		{
			desc:    "isKore=false, isMagma=false, fee=10e18 [001]",
			isKore:  false,
			isMagma: false,
			fee:     10e18,
			expected: &Result{
				total:  10e18,
				reward: 10e18,
				burnt:  0,
			},
		},
		{
			desc:    "isKore=false, isMagma=true, fee=1000 [010]",
			isKore:  false,
			isMagma: true,
			fee:     1000,
			expected: &Result{
				total:  1000,
				reward: 500,
				burnt:  500,
			},
		},
		{
			desc:    "isKore=false, isMagma=true, fee=10e18 [011]",
			isKore:  false,
			isMagma: true,
			fee:     10e18,
			expected: &Result{
				total:  10e18,
				reward: 5e18,
				burnt:  5e18,
			},
		},
		{
			desc:    "isKore=true, isMagma=true, fee=1000 [110]",
			isKore:  true,
			isMagma: true,
			fee:     1000,
			expected: &Result{
				total:  1000,
				reward: 0,
				burnt:  1000,
			},
		},
		{
			desc:    "isKore=true, isMagma=true, fee=10e18 [111]",
			isKore:  true,
			isMagma: true,
			fee:     10e18,
			expected: &Result{
				total:  10e18,
				reward: 4.3472e18, // 5 - minted*0.34*0.2
				burnt:  5.6528e18, // 5 + minted*0.34*0.8
			},
		},
	}

	for _, tc := range testcases {
		header := &types.Header{
			Number:     big.NewInt(1),
			GasUsed:    tc.fee,
			BaseFee:    big.NewInt(1),
			Rewardbase: proposerAddr,
		}

		config := getTestConfig()
		if !tc.isKore {
			config = noKore(config)
		}
		if !tc.isMagma {
			config = noMagma(config)
		}

		rc, err := NewRewardConfig(header, config)
		assert.Nil(t, err)

		total, reward, burnt := calcDeferredFee(rc)
		actual := &Result{
			total:  total.Uint64(),
			reward: reward.Uint64(),
			burnt:  burnt.Uint64(),
		}
		assert.Equal(t, tc.expected, actual, "failed tc: %s", tc.desc)
	}
}

func TestRewardDistributor_calcDeferredFee_nodeferred(t *testing.T) {
	header := &types.Header{
		Number:     big.NewInt(1),
		GasUsed:    1000,
		BaseFee:    big.NewInt(1),
		Rewardbase: proposerAddr,
	}

	rc, err := NewRewardConfig(header, noDeferred(getTestConfig()))
	assert.Nil(t, err)

	total, reward, burnt := calcDeferredFee(rc)
	assert.Equal(t, uint64(0), total.Uint64())
	assert.Equal(t, uint64(0), reward.Uint64())
	assert.Equal(t, uint64(0), burnt.Uint64())
}

func TestRewardDistributor_calcSplit(t *testing.T) {
	type Result struct{ proposer, stakers, kgf, kir, remaining uint64 }

	header := &types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(0), // placeholder
	}

	testcases := []struct {
		desc     string
		isKore   bool
		fee      uint64
		expected *Result
	}{
		{
			desc:   "kore=false, fee=0",
			isKore: false,
			fee:    0,
			expected: &Result{
				proposer:  3.264e18, // minted * 0.34
				stakers:   0,
				kgf:       5.184e18, // minted * 0.54
				kir:       1.152e18, // minted * 0.12
				remaining: 0,
			},
		},
		{
			desc:   "kore=false, fee=55555",
			isKore: false,
			fee:    55555,
			expected: &Result{
				proposer:  3.264e18 + 18888, // (minted + fee) * 0.34
				stakers:   0,
				kgf:       5.184e18 + 29999, // (minted + fee) * 0.54
				kir:       1.152e18 + 6666,  // (minted + fee) * 0.12
				remaining: 2,
			},
		},
		{
			desc:   "kore=true, fee=0",
			isKore: true,
			fee:    0,
			expected: &Result{
				proposer:  0.6528e18, // minted * 0.34 * 0.2
				stakers:   2.6112e18, // minted * 0.34 * 0.8
				kgf:       5.184e18,  // minted * 0.54
				kir:       1.152e18,  // minted * 0.12
				remaining: 0,
			},
		},
		{
			desc:   "kore=true, fee=55555",
			isKore: true,
			fee:    55555,
			expected: &Result{
				proposer:  0.6528e18 + 55555, // minted * 0.34 * 0.2 + fee
				stakers:   2.6112e18,         // minted * 0.34 * 0.8
				kgf:       5.184e18,          // minted * 0.54
				kir:       1.152e18,          // minted * 0.12
				remaining: 0,
			},
		},
	}

	for _, tc := range testcases {
		config := getTestConfig()
		if !tc.isKore {
			config = noKore(config)
		}
		rc, err := NewRewardConfig(header, config)
		assert.Nil(t, err)

		fee := new(big.Int).SetUint64(tc.fee)
		proposer, stakers, kgf, kir, remaining := calcSplit(rc, minted, fee)
		actual := &Result{
			proposer:  proposer.Uint64(),
			stakers:   stakers.Uint64(),
			kgf:       kgf.Uint64(),
			kir:       kir.Uint64(),
			remaining: remaining.Uint64(),
		}
		assert.Equal(t, tc.expected, actual, "failed tc: %s", tc.desc)

		expectedTotalAmount := big.NewInt(0)
		expectedTotalAmount = expectedTotalAmount.Add(minted, fee)

		actualTotalAmount := big.NewInt(0)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, proposer)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, stakers)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, kgf)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, kir)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, remaining)
		assert.Equal(t, expectedTotalAmount, actualTotalAmount, "failed tc: %s", tc.desc)
	}
}

func TestRewardDistributor_calcShares(t *testing.T) {
	type Result struct {
		shares    map[common.Address]*big.Int
		remaining uint64
	}

	testcases := []struct {
		desc        string
		stakingInfo *StakingInfo
		stakeReward *big.Int
		expected    *Result
	}{
		{
			desc:        "all nodes 0%",
			stakingInfo: genStakingInfo(5, nil, nil),
			stakeReward: big.NewInt(500),
			expected: &Result{
				shares:    map[common.Address]*big.Int{},
				remaining: 500,
			},
		},
		{
			desc:        "no staking info",
			stakingInfo: nil,
			stakeReward: big.NewInt(500),
			expected: &Result{
				shares:    map[common.Address]*big.Int{},
				remaining: 500,
			},
		},
		{
			desc:        "CN0: 100%",
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{0: minStaking + 1}),
			stakeReward: big.NewInt(500),
			expected: &Result{
				shares: map[common.Address]*big.Int{
					intToAddress(rewardBaseAddr): big.NewInt(500),
				},
				remaining: 0,
			},
		},
		{
			desc: "CN0, CN1: 50%",
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 1,
				1: minStaking + 1,
			}),
			stakeReward: big.NewInt(500),
			expected: &Result{
				shares: map[common.Address]*big.Int{
					intToAddress(rewardBaseAddr):     big.NewInt(250),
					intToAddress(rewardBaseAddr + 1): big.NewInt(250),
				},
				remaining: 0,
			},
		},
		{
			desc: "CN0: 66%, CN1: 33%",
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 1,
			}),
			stakeReward: big.NewInt(500),
			expected: &Result{
				shares: map[common.Address]*big.Int{
					intToAddress(rewardBaseAddr):     big.NewInt(333),
					intToAddress(rewardBaseAddr + 1): big.NewInt(166),
				},
				remaining: 1,
			},
		},
		{
			desc: "CN0: 66/97, CN1: 17/97, CN2: 11/97, CN3: 2/97, CN4: 1/97",
			stakingInfo: genStakingInfo(7, nil, map[int]uint64{
				0: minStaking + 66,
				1: minStaking + 17,
				2: minStaking + 11,
				3: minStaking + 2,
				4: minStaking + 1, // total: 97
			}),
			stakeReward: big.NewInt(555),
			expected: &Result{
				shares: map[common.Address]*big.Int{
					intToAddress(rewardBaseAddr):     big.NewInt(377),
					intToAddress(rewardBaseAddr + 1): big.NewInt(97),
					intToAddress(rewardBaseAddr + 2): big.NewInt(62),
					intToAddress(rewardBaseAddr + 3): big.NewInt(11),
					intToAddress(rewardBaseAddr + 4): big.NewInt(5),
				},
				remaining: 3,
			},
		},
	}

	for _, tc := range testcases {
		shares, remaining := calcShares(tc.stakingInfo, tc.stakeReward, minStaking)
		actual := &Result{
			shares:    shares,
			remaining: remaining.Uint64(),
		}
		assert.Equal(t, tc.expected, actual, "failed tc: %s", tc.desc)
	}
}

func benchSetup() (*types.Header, *params.ChainConfig) {
	// in the worst case, distribute stake shares among N
	amounts := make(map[int]uint64)
	N := 50
	for i := 0; i < N; i++ {
		amounts[i] = minStaking + 1
	}

	stakingInfo := genStakingInfo(N, nil, amounts)
	SetTestStakingManagerWithStakingInfoCache(stakingInfo)

	config := getTestConfig()

	header := &types.Header{}
	header.BaseFee = big.NewInt(30000000000)
	header.Number = big.NewInt(0)
	header.Rewardbase = intToAddress(rewardBaseAddr)
	return header, config
}

func Benchmark_CalcDeferredReward(b *testing.B) {
	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	header, config := benchSetup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalcDeferredReward(header, config)
	}
}

func TestRewardConfigCache_parseRewardRatio(t *testing.T) {
	testCases := []struct {
		s   string
		cn  int64
		poc int64
		kir int64
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

	for i := 0; i < len(testCases); i++ {
		cn, poc, kir, total, err := parseRewardRatio(testCases[i].s)

		assert.Equal(t, testCases[i].cn, cn)
		assert.Equal(t, testCases[i].poc, poc)
		assert.Equal(t, testCases[i].kir, kir)
		assert.Equal(t, testCases[i].err, err)

		expectedTotal := testCases[i].cn + testCases[i].poc + testCases[i].kir
		assert.Equal(t, expectedTotal, total)
	}
}

func TestRewardConfigCache_parseRewardKip82Ratio(t *testing.T) {
	testCases := []struct {
		s        string
		proposer int64
		staking  int64
		err      error
	}{
		{"34/54", 34, 54, nil},
		{"20/80", 20, 80, nil},
		{"0/100", 0, 100, nil},
		{"34,54", 0, 0, errInvalidFormat},
		{"", 0, 0, errInvalidFormat},
		{"//", 0, 0, errInvalidFormat},
		{"1/", 0, 0, errParsingRatio},
		{"/1", 0, 0, errParsingRatio},
		{"1/2/", 0, 0, errInvalidFormat},
		{"3.3/3.3", 0, 0, errParsingRatio},
		{"a/b", 0, 0, errParsingRatio},
	}

	for i := 0; i < len(testCases); i++ {
		proposer, staking, total, err := parseRewardKip82Ratio(testCases[i].s)

		assert.Equal(t, testCases[i].proposer, proposer)
		assert.Equal(t, testCases[i].staking, staking)
		assert.Equal(t, testCases[i].err, err, "tc[%d] failed", i)

		expectedTotal := testCases[i].proposer + testCases[i].staking
		assert.Equal(t, expectedTotal, total)
	}
}
