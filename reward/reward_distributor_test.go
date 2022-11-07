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
	"github.com/klaytn/klaytn/log"
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
		params.MinimumStake:        "5000000",
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

var (
	cnBaseAddr     = 500
	stakeBaseAddr  = 600
	rewardBaseAddr = 700
	minStaking     = uint64(2000000) // changing this value will not change the governance's min staking
	minted, _      = big.NewInt(0).SetString("9600000000000000000", 10)
	proposer       = common.StringToAddress("0x1552F52D459B713E0C4558e66C8c773a75615FA8")
	kir            = intToAddress(1000)
	kgf            = intToAddress(2000)
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
		KIRAddr:               kir,
		PoCAddr:               kgf,
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

func TestRewardDistributor_GetActualReward(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	header := &types.Header{
		Number:     big.NewInt(1),
		GasUsed:    1000,
		BaseFee:    big.NewInt(1),
		Rewardbase: proposer,
	}

	testcases := []struct {
		header      *types.Header
		config      *params.ChainConfig
		stakingInfo *StakingInfo
		expected    *RewardSpec
	}{
		{
			header: header,
			config: noMagma(roundrobin(getTestConfig())),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(0),
				Proposer: big.NewInt(0).SetUint64(9600000000000001000),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(0).SetUint64(9600000000000001000),
				},
			},
		},
		{
			header: header,
			config: roundrobin(getTestConfig()),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(500),
				Proposer: big.NewInt(0).SetUint64(9600000000000000500),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(0).SetUint64(9600000000000000500),
				},
			},
		},
		{ // before magma, default stakings
			header:      header,
			config:      noMagma(getTestConfig()),
			stakingInfo: genStakingInfo(5, nil, nil),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(0),
				Proposer: big.NewInt(0).SetUint64(3264000000000000340),
				Stakers:  big.NewInt(0),
				Kgf:      big.NewInt(5184000000000000540),
				Kir:      big.NewInt(1152000000000000120),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(3264000000000000340),
					kir:      big.NewInt(1152000000000000120),
					kgf:      big.NewInt(5184000000000000540),
				},
			},
		},
		{ // before magma, more-than-default staking
			header: header,
			config: noMagma(getTestConfig()),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 1,
				1: minStaking + 1,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(0),
				Proposer: big.NewInt(3264000000000000340),
				Stakers:  big.NewInt(0),
				Kgf:      big.NewInt(5184000000000000540),
				Kir:      big.NewInt(1152000000000000120),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(3264000000000000340),
					kir:      big.NewInt(1152000000000000120),
					kgf:      big.NewInt(5184000000000000540),
				},
			},
		},
		{ // before kore, more-than-default staking
			header: header,
			config: noKore(getTestConfig()),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 2,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(500),
				Proposer: big.NewInt(3264000000000000170),
				Stakers:  big.NewInt(0),
				Kgf:      big.NewInt(5184000000000000270),
				Kir:      big.NewInt(1152000000000000060),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(3264000000000000170),
					kir:      big.NewInt(1152000000000000060),
					kgf:      big.NewInt(5184000000000000270),
				},
			},
		},
		{ // after kore, more-than-default staking, small fee
			header: header,
			config: getTestConfig(),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 1,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(1000),
				Proposer: big.NewInt(652800000000000000),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(5184000000000000000),
				Kir:      big.NewInt(1152000000000000000),
				Rewards: map[common.Address]*big.Int{
					proposer:                         big.NewInt(652800000000000000),
					kgf:                              big.NewInt(5184000000000000000),
					kir:                              big.NewInt(1152000000000000000),
					intToAddress(rewardBaseAddr):     big.NewInt(1740800000000000000),
					intToAddress(rewardBaseAddr + 1): big.NewInt(870400000000000000),
				},
			},
		},
		{ // after kore, more-than-default staking, large fee
			header: &types.Header{
				Number:     big.NewInt(1),
				GasUsed:    3e18,
				BaseFee:    big.NewInt(1),
				Rewardbase: proposer,
			},
			config: getTestConfig(),
			stakingInfo: genStakingInfo(7, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 1,
				2: minStaking + 2,
				3: minStaking + 1,
				4: minStaking + 2,
				5: minStaking + 1,
				6: minStaking + 2,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(3e18),
				Burnt:    big.NewInt(2152800000000000000),
				Proposer: big.NewInt(1500000000000000005),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(5184000000000000000),
				Kir:      big.NewInt(1152000000000000000),
				Rewards: map[common.Address]*big.Int{
					proposer:                         big.NewInt(1500000000000000005),
					kgf:                              big.NewInt(5184000000000000000),
					kir:                              big.NewInt(1152000000000000000),
					intToAddress(rewardBaseAddr):     big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 1): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 2): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 3): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 4): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 5): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 6): big.NewInt(474763636363636363),
				},
			},
		},
		{ // after kore, more-than-default staking, large fee, proposer = rewardbase
			header: &types.Header{
				Number:     big.NewInt(1),
				GasUsed:    3e18,
				BaseFee:    big.NewInt(1),
				Rewardbase: intToAddress(rewardBaseAddr),
			},
			config: getTestConfig(),
			stakingInfo: genStakingInfo(7, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 1,
				2: minStaking + 2,
				3: minStaking + 1,
				4: minStaking + 2,
				5: minStaking + 1,
				6: minStaking + 2,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(3e18),
				Burnt:    big.NewInt(2152800000000000000),
				Proposer: big.NewInt(1500000000000000005),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(5184000000000000000),
				Kir:      big.NewInt(1152000000000000000),
				Rewards: map[common.Address]*big.Int{
					kgf:                              big.NewInt(5184000000000000000),
					kir:                              big.NewInt(1152000000000000000),
					intToAddress(rewardBaseAddr):     big.NewInt(1974763636363636368),
					intToAddress(rewardBaseAddr + 1): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 2): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 3): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 4): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 5): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 6): big.NewInt(474763636363636363),
				},
			},
		},
		{ // after kore, more-than-default staking, small fee, kgf = kir = 0
			header: header,
			config: getTestConfig(),
			stakingInfo: &StakingInfo{
				BlockNum:              0,
				CouncilNodeAddrs:      []common.Address{intToAddress(cnBaseAddr)},
				CouncilStakingAddrs:   []common.Address{intToAddress(stakeBaseAddr)},
				CouncilRewardAddrs:    []common.Address{intToAddress(rewardBaseAddr)},
				KIRAddr:               common.Address{},
				PoCAddr:               common.Address{},
				UseGini:               false,
				CouncilStakingAmounts: []uint64{minStaking + 1},
			},
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(1000),
				Proposer: big.NewInt(6988800000000000000),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(0),
				Kir:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					proposer:                     big.NewInt(6988800000000000000),
					intToAddress(rewardBaseAddr): big.NewInt(2611200000000000000),
				},
			},
		},
		{ // no staking info
			header:      header,
			config:      getTestConfig(),
			stakingInfo: nil,
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(1000),
				Proposer: minted,
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(0),
				Kir:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					proposer: minted,
				},
			},
		},
	}

	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	for i, tc := range testcases {
		if tc.stakingInfo == nil {
			SetTestStakingManager(nil)
		} else {
			SetTestStakingManagerWithStakingInfoCache(tc.stakingInfo)
		}
		spec, err := GetBlockReward(tc.header, tc.config)
		assert.Nil(t, err, "testcases[%d] failed", i)
		assert.Equal(t, tc.expected, spec, "testcases[%d] failed", i)
	}
}

func TestRewardDistributor_CalcSimpleReward(t *testing.T) {
	header := &types.Header{
		Number:     big.NewInt(1),
		GasUsed:    1000,
		BaseFee:    big.NewInt(1),
		Rewardbase: proposer,
	}

	testcases := []struct {
		header   *types.Header
		config   *params.ChainConfig
		expected *RewardSpec
	}{
		{ // before magma
			header: header,
			config: noMagma(roundrobin(getTestConfig())),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(0),
				Proposer: big.NewInt(0).SetUint64(9600000000000001000),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(0).SetUint64(9600000000000001000),
				},
			},
		},
		{ // after magma
			header: header,
			config: roundrobin(getTestConfig()),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(500), // 50% of tx fee burnt
				Proposer: big.NewInt(0).SetUint64(9600000000000000500),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(0).SetUint64(9600000000000000500),
				},
			},
		},
	}

	for i, tc := range testcases {
		spec, err := CalcDeferredRewardSimple(tc.header, tc.config)
		assert.Nil(t, err, "testcases[%d] failed", i)
		assert.Equal(t, tc.expected, spec, "testcases[%d] failed", i)
	}
}

func TestRewardDistributor_CalcDeferredReward(t *testing.T) {
	header := &types.Header{
		Number:     big.NewInt(1),
		GasUsed:    1000,
		BaseFee:    big.NewInt(1),
		Rewardbase: proposer,
	}

	testcases := []struct {
		header      *types.Header
		config      *params.ChainConfig
		stakingInfo *StakingInfo
		expected    *RewardSpec
	}{
		{ // before magma, default stakings
			header:      header,
			config:      noMagma(getTestConfig()),
			stakingInfo: genStakingInfo(5, nil, nil),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(0),
				Proposer: big.NewInt(0).SetUint64(3264000000000000340),
				Stakers:  big.NewInt(0),
				Kgf:      big.NewInt(5184000000000000540),
				Kir:      big.NewInt(1152000000000000120),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(3264000000000000340),
					kir:      big.NewInt(1152000000000000120),
					kgf:      big.NewInt(5184000000000000540),
				},
			},
		},
		{ // before magma, more-than-default staking
			header: header,
			config: noMagma(getTestConfig()),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 1,
				1: minStaking + 1,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(0),
				Proposer: big.NewInt(3264000000000000340),
				Stakers:  big.NewInt(0),
				Kgf:      big.NewInt(5184000000000000540),
				Kir:      big.NewInt(1152000000000000120),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(3264000000000000340),
					kir:      big.NewInt(1152000000000000120),
					kgf:      big.NewInt(5184000000000000540),
				},
			},
		},
		{ // before kore, more-than-default staking
			header: header,
			config: noKore(getTestConfig()),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 2,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(500),
				Proposer: big.NewInt(3264000000000000170),
				Stakers:  big.NewInt(0),
				Kgf:      big.NewInt(5184000000000000270),
				Kir:      big.NewInt(1152000000000000060),
				Rewards: map[common.Address]*big.Int{
					proposer: big.NewInt(3264000000000000170),
					kir:      big.NewInt(1152000000000000060),
					kgf:      big.NewInt(5184000000000000270),
				},
			},
		},
		{ // after kore, more-than-default staking, small fee
			header: header,
			config: getTestConfig(),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 1,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(1000),
				Proposer: big.NewInt(652800000000000000),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(5184000000000000000),
				Kir:      big.NewInt(1152000000000000000),
				Rewards: map[common.Address]*big.Int{
					proposer:                         big.NewInt(652800000000000000),
					kgf:                              big.NewInt(5184000000000000000),
					kir:                              big.NewInt(1152000000000000000),
					intToAddress(rewardBaseAddr):     big.NewInt(1740800000000000000),
					intToAddress(rewardBaseAddr + 1): big.NewInt(870400000000000000),
				},
			},
		},
		{ // after kore, more-than-default staking, large fee
			header: &types.Header{
				Number:     big.NewInt(1),
				GasUsed:    3e18,
				BaseFee:    big.NewInt(1),
				Rewardbase: proposer,
			},
			config: getTestConfig(),
			stakingInfo: genStakingInfo(7, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 1,
				2: minStaking + 2,
				3: minStaking + 1,
				4: minStaking + 2,
				5: minStaking + 1,
				6: minStaking + 2,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(3e18),
				Burnt:    big.NewInt(2152800000000000000),
				Proposer: big.NewInt(1500000000000000005),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(5184000000000000000),
				Kir:      big.NewInt(1152000000000000000),
				Rewards: map[common.Address]*big.Int{
					proposer:                         big.NewInt(1500000000000000005),
					kgf:                              big.NewInt(5184000000000000000),
					kir:                              big.NewInt(1152000000000000000),
					intToAddress(rewardBaseAddr):     big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 1): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 2): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 3): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 4): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 5): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 6): big.NewInt(474763636363636363),
				},
			},
		},
		{ // after kore, more-than-default staking, large fee, proposer = rewardbase
			header: &types.Header{
				Number:     big.NewInt(1),
				GasUsed:    3e18,
				BaseFee:    big.NewInt(1),
				Rewardbase: intToAddress(rewardBaseAddr),
			},
			config: getTestConfig(),
			stakingInfo: genStakingInfo(7, nil, map[int]uint64{
				0: minStaking + 2,
				1: minStaking + 1,
				2: minStaking + 2,
				3: minStaking + 1,
				4: minStaking + 2,
				5: minStaking + 1,
				6: minStaking + 2,
			}),
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(3e18),
				Burnt:    big.NewInt(2152800000000000000),
				Proposer: big.NewInt(1500000000000000005),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(5184000000000000000),
				Kir:      big.NewInt(1152000000000000000),
				Rewards: map[common.Address]*big.Int{
					kgf:                              big.NewInt(5184000000000000000),
					kir:                              big.NewInt(1152000000000000000),
					intToAddress(rewardBaseAddr):     big.NewInt(1974763636363636368),
					intToAddress(rewardBaseAddr + 1): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 2): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 3): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 4): big.NewInt(474763636363636363),
					intToAddress(rewardBaseAddr + 5): big.NewInt(237381818181818181),
					intToAddress(rewardBaseAddr + 6): big.NewInt(474763636363636363),
				},
			},
		},
		{ // after kore, more-than-default staking, small fee, kgf = kir = 0
			header: header,
			config: getTestConfig(),
			stakingInfo: &StakingInfo{
				BlockNum:              0,
				CouncilNodeAddrs:      []common.Address{intToAddress(cnBaseAddr)},
				CouncilStakingAddrs:   []common.Address{intToAddress(stakeBaseAddr)},
				CouncilRewardAddrs:    []common.Address{intToAddress(rewardBaseAddr)},
				KIRAddr:               common.Address{},
				PoCAddr:               common.Address{},
				UseGini:               false,
				CouncilStakingAmounts: []uint64{minStaking + 1},
			},
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(1000),
				Proposer: big.NewInt(6988800000000000000),
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(0),
				Kir:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					proposer:                     big.NewInt(6988800000000000000),
					intToAddress(rewardBaseAddr): big.NewInt(2611200000000000000),
				},
			},
		},
		{ // no staking info
			header:      header,
			config:      getTestConfig(),
			stakingInfo: nil,
			expected: &RewardSpec{
				Minted:   minted,
				Fee:      big.NewInt(1000),
				Burnt:    big.NewInt(1000),
				Proposer: minted,
				Stakers:  big.NewInt(2611200000000000000),
				Kgf:      big.NewInt(0),
				Kir:      big.NewInt(0),
				Rewards: map[common.Address]*big.Int{
					proposer: minted,
				},
			},
		},
	}

	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	for i, tc := range testcases {
		if tc.stakingInfo == nil {
			SetTestStakingManager(nil)
		} else {
			SetTestStakingManagerWithStakingInfoCache(tc.stakingInfo)
		}
		spec, err := CalcDeferredReward(tc.header, tc.config)
		assert.Nil(t, err, "testcases[%d] failed", i)
		assert.Equal(t, tc.expected, spec, "testcases[%d] failed", i)
	}
}

func TestRewardDistributor_calcDeferredFee(t *testing.T) {
	type Result struct{ total, reward, burnt uint64 }

	header := &types.Header{
		Number:     big.NewInt(1),
		GasUsed:    1000,
		BaseFee:    big.NewInt(1),
		Rewardbase: proposer,
	}

	testcases := []struct {
		header   *types.Header
		config   *params.ChainConfig
		expected *Result
	}{
		{
			header: header,
			config: noDeferred(getTestConfig()),
			expected: &Result{
				total:  0,
				reward: 0,
				burnt:  0,
			},
		},
		{
			header: header,
			config: noMagma(getTestConfig()),
			expected: &Result{
				total:  1000,
				reward: 1000,
				burnt:  0,
			},
		},
		{
			header: header,
			config: noKore(getTestConfig()),
			expected: &Result{
				total:  1000,
				reward: 500,
				burnt:  500,
			},
		},
		{
			header: header,
			config: getTestConfig(),
			expected: &Result{
				total:  1000,
				reward: 0,
				burnt:  1000,
			},
		},
		{
			header: &types.Header{
				Number:     big.NewInt(1),
				GasUsed:    10e18,
				BaseFee:    big.NewInt(1),
				Rewardbase: proposer,
			},
			config: getTestConfig(),
			expected: &Result{
				total:  10e18,
				reward: 4347200000000000000, // 5 klay - 9.6*0.34*0.2 klay
				burnt:  5652800000000000000, // 5 klay + 9.6*0.34*0.2 klay
			},
		},
	}

	for i, tc := range testcases {
		rc, err := NewRewardConfig(tc.header, tc.config)
		assert.Nil(t, err)

		total, reward, burnt := calcDeferredFee(rc)
		actual := &Result{
			total:  total.Uint64(),
			reward: reward.Uint64(),
			burnt:  burnt.Uint64(),
		}
		assert.Equal(t, tc.expected, actual, "testcases[%d] failed", i)
	}
}

func TestRewardDistributor_calcSplit(t *testing.T) {
	type Result struct{ proposer, stakers, kgf, kir, remaining uint64 }

	header := &types.Header{
		Number:  big.NewInt(1),
		BaseFee: big.NewInt(0), // placeholder
	}

	testcases := []struct {
		header   *types.Header
		config   *params.ChainConfig
		fee      *big.Int
		expected *Result
	}{
		{
			header: header,
			config: noKore(getTestConfig()),
			fee:    big.NewInt(0),
			expected: &Result{
				proposer:  3.264e18, // 9.6e18 * 0.34
				stakers:   0,
				kgf:       5.184e18, // 9.6e18 * 0.54
				kir:       1.152e18, // 9.6e18 * 0.12
				remaining: 0,
			},
		},
		{
			header: header,
			config: noKore(getTestConfig()),
			fee:    big.NewInt(55555),
			expected: &Result{
				proposer:  3264000000000018888, // (9.6e18 + 55555) * 0.34
				stakers:   0,
				kgf:       5184000000000029999, // (9.6e18 + 55555) * 0.54
				kir:       1152000000000006666, // (9.6e18 + 55555) * 0.12
				remaining: 2,
			},
		},
		{
			header: header,
			config: getTestConfig(),
			fee:    big.NewInt(0),
			expected: &Result{
				proposer:  652800000000000000,  // 9.6e18 * 0.34 * 0.2
				stakers:   2611200000000000000, // 9.6e18 * 0.34 * 0.8
				kgf:       5184000000000000000, // 9.6e18 * 0.54
				kir:       1152000000000000000, // 9.6e18 * 0.12
				remaining: 0,
			},
		},
		{
			header: header,
			config: getTestConfig(),
			fee:    big.NewInt(55555),
			expected: &Result{
				proposer:  652800000000055555,  // 9.6e18 * 0.34 * 0.2 + 55555
				stakers:   2611200000000000000, // 9.6e18 * 0.34 * 0.8
				kgf:       5184000000000000000, // 9.6e18 * 0.54
				kir:       1152000000000000000, // 9.6e18 * 0.12
				remaining: 0,
			},
		},
	}

	for i, tc := range testcases {
		rc, err := NewRewardConfig(tc.header, tc.config)
		assert.Nil(t, err)

		proposer, stakers, kgf, kir, remaining := calcSplit(rc, minted, tc.fee)
		actual := &Result{
			proposer:  proposer.Uint64(),
			stakers:   stakers.Uint64(),
			kgf:       kgf.Uint64(),
			kir:       kir.Uint64(),
			remaining: remaining.Uint64(),
		}
		assert.Equal(t, tc.expected, actual, "testcases[%d] failed", i)

		expectedTotalAmount := big.NewInt(0)
		expectedTotalAmount = expectedTotalAmount.Add(minted, tc.fee)

		actualTotalAmount := big.NewInt(0)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, proposer)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, stakers)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, kgf)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, kir)
		actualTotalAmount = actualTotalAmount.Add(actualTotalAmount, remaining)
		assert.Equal(t, expectedTotalAmount, actualTotalAmount, "testcases[%d] failed", i)
	}
}

func TestRewardDistributor_calcShares(t *testing.T) {
	type Result struct {
		shares    map[common.Address]*big.Int
		remaining uint64
	}

	testcases := []struct {
		config      *params.ChainConfig
		stakingInfo *StakingInfo
		stakeReward *big.Int
		expected    *Result
	}{
		{
			config:      getTestConfig(),
			stakingInfo: genStakingInfo(5, nil, nil),
			stakeReward: big.NewInt(500),
			expected: &Result{
				shares:    map[common.Address]*big.Int{},
				remaining: 500,
			},
		},
		{
			config:      getTestConfig(),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{0: minStaking * 2}),
			stakeReward: big.NewInt(500),
			expected: &Result{
				shares: map[common.Address]*big.Int{
					intToAddress(rewardBaseAddr): big.NewInt(500),
				},
				remaining: 0,
			},
		},
		{
			config: getTestConfig(),
			stakingInfo: genStakingInfo(5, nil, map[int]uint64{
				0: minStaking * 2,
				1: minStaking * 2,
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
			config: getTestConfig(),
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
	}

	for i, tc := range testcases {
		shares, remaining := calcShares(tc.stakingInfo, tc.stakeReward, minStaking)
		actual := &Result{
			shares:    shares,
			remaining: remaining.Uint64(),
		}
		assert.Equal(t, tc.expected, actual, "testcases[%d] failed", i)
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
