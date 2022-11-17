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
	"strconv"
	"strings"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
)

var logger = log.NewModuleLogger(log.Reward)

type BalanceAdder interface {
	AddBalance(addr common.Address, v *big.Int)
}

// Cannot use governance.Engine because of cyclic dependency.
// Instead declare only the methods used by this package.
type governanceHelper interface {
	Params() *params.GovParamSet
	ParamsAt(num uint64) (*params.GovParamSet, error)
}

type RewardDistributor struct {
	rcc *rewardConfigCache
	gh  governanceHelper
}

type RewardSpec struct {
	Minted   *big.Int
	Fee      *big.Int
	Burnt    *big.Int
	Proposer *big.Int
	Stakers  *big.Int
	Kgf      *big.Int
	Kir      *big.Int
	Rewards  map[common.Address]*big.Int
}

func NewRewardDistributor(gh governanceHelper) *RewardDistributor {
	return &RewardDistributor{
		rcc: newRewardConfigCache(gh),
		gh:  gh,
	}
}

// getTotalTxFee returns the total transaction gas fee of the block.
func (rd *RewardDistributor) getTotalTxFee(header *types.Header, rewardConfig *rewardConfig) *big.Int {
	totalGasUsed := big.NewInt(0).SetUint64(header.GasUsed)
	if header.BaseFee != nil {
		// magma hardfork
		return totalGasUsed.Mul(totalGasUsed, header.BaseFee)
	} else {
		return totalGasUsed.Mul(totalGasUsed, rewardConfig.unitPrice)
	}
}

func (rd *RewardDistributor) txFeeBurning(txFee *big.Int) *big.Int {
	return txFee.Div(txFee, big.NewInt(2))
}

// DistributeBlockReward distributes a given block's reward at the end of block processing
func (rd *RewardDistributor) DistributeBlockReward(b BalanceAdder, rewards map[common.Address]*big.Int) {
	for addr, amount := range rewards {
		b.AddBalance(addr, amount)
	}
}

// Function hierarchy
// - get_actual_reward
//   - calc_simple_reward
//   - calc_deferred_reward
//     - calc_deferred_fee
//     - calc_split
//     - calc_shares

// GetActualReward returns the actual reward amounts paid in this block
// Used in klay_getReward RPC
func GetActualReward(header *types.Header, config *params.ChainConfig) (*RewardSpec, error) {
	var spec *RewardSpec
	var err error
	if config.Istanbul == nil {
		return nil, errors.New("no IstanbulConfig")
	}

	policy := config.Istanbul.ProposerPolicy
	if policy == uint64(istanbul.RoundRobin) || policy == uint64(istanbul.Sticky) {
		spec, err = CalcSimpleReward(header, config)
		if err != nil {
			return nil, err
		}
	} else {
		spec, err = CalcDeferredReward(header, config)
		if err != nil {
			return nil, err
		}

		if !config.Governance.Reward.DeferredTxFee {
			var blockFee *big.Int

			if config.IsMagmaForkEnabled(header.Number) {
				blockFee = new(big.Int).Mul(
					big.NewInt(0).SetUint64(header.GasUsed),
					header.BaseFee)
			} else {
				blockFee = new(big.Int).Mul(
					big.NewInt(0).SetUint64(header.GasUsed),
					big.NewInt(0).SetUint64(config.UnitPrice))
			}
			spec.Proposer = spec.Proposer.Add(spec.Proposer, blockFee)
			spec.Rewards[header.Rewardbase] = spec.Rewards[header.Rewardbase].Add(
				spec.Rewards[header.Rewardbase], blockFee)
		}
	}

	return spec, nil
}

// CalcSimpleReward distributes rewards to proposer after optional fee burning
// this behaves similar to the previous MintKLAY
func CalcSimpleReward(header *types.Header, config *params.ChainConfig) (*RewardSpec, error) {
	rewardConfig := config.Governance.Reward
	if rewardConfig == nil {
		return nil, errors.New("no rewardConfig")
	}

	minted := rewardConfig.MintingAmount
	var totalFee, rewardFee, burntFee *big.Int
	totalGasUsed := big.NewInt(0).SetUint64(header.GasUsed)

	if config.IsMagmaForkEnabled(header.Number) {
		totalFee = big.NewInt(0).Mul(totalGasUsed, header.BaseFee)
		rewardFee = big.NewInt(0).Div(totalFee, big.NewInt(2))
		burntFee = big.NewInt(0).Div(totalFee, big.NewInt(2))
	} else {
		unitPrice := big.NewInt(0).SetUint64(config.UnitPrice)
		totalFee = big.NewInt(0).Mul(totalGasUsed, unitPrice)
		rewardFee = big.NewInt(0)
		rewardFee.Set(totalFee)
		burntFee = big.NewInt(0)
	}

	proposer := big.NewInt(0).Add(minted, rewardFee)

	return &RewardSpec{
		Minted:   minted,
		Fee:      totalFee,
		Burnt:    burntFee,
		Proposer: proposer,
		Rewards:  map[common.Address]*big.Int{header.Rewardbase: proposer},
	}, nil
}

// CalcDeferredReward calculates the deferred rewards,
// which are determined at the end of block processing.
func CalcDeferredReward(header *types.Header, config *params.ChainConfig) (*RewardSpec, error) {
	rewardConfig := config.Governance.Reward
	if rewardConfig == nil {
		return nil, errors.New("no rewardConfig")
	}

	totalFee, rewardFee, burntFee := calcDeferredFee(header, config)
	stakingInfo := GetStakingInfo(header.Number.Uint64())

	proposer, stakers, kgf, kir, splitRem := calcSplit(header, config, rewardFee)
	shares, shareRem := calcShares(rewardConfig, stakingInfo, stakers)

	kgf = kgf.Add(kgf, splitRem)
	proposer = proposer.Add(proposer, shareRem)

	if stakingInfo == nil || common.EmptyAddress(stakingInfo.PoCAddr) {
		logger.Debug("KGF empty, proposer gets its portion", "kgf", kgf)
		proposer = proposer.Add(proposer, kgf)
		kgf = big.NewInt(0)
	}
	if stakingInfo == nil || common.EmptyAddress(stakingInfo.KIRAddr) {
		logger.Debug("KIR empty, proposer gets its portion", "kir", kir)
		proposer = proposer.Add(proposer, kir)
		kir = big.NewInt(0)
	}

	spec := &RewardSpec{
		Minted:   rewardConfig.MintingAmount,
		Fee:      totalFee,
		Burnt:    burntFee,
		Proposer: proposer,
		Stakers:  stakers,
		Kgf:      kgf,
		Kir:      kir,
	}

	spec.Rewards = make(map[common.Address]*big.Int)
	increment(spec.Rewards, header.Rewardbase, proposer)

	if stakingInfo != nil && !common.EmptyAddress(stakingInfo.PoCAddr) {
		increment(spec.Rewards, stakingInfo.PoCAddr, kgf)
	}
	if stakingInfo != nil && !common.EmptyAddress(stakingInfo.KIRAddr) {
		increment(spec.Rewards, stakingInfo.KIRAddr, kir)
	}

	for rewardAddr, rewardAmount := range shares {
		increment(spec.Rewards, rewardAddr, rewardAmount)
	}
	logger.Debug("calcDeferredReward returns", "spec", spec)

	return spec, nil
}

// calcDeferredFee splits fee into (total, reward, burnt)
func calcDeferredFee(header *types.Header, config *params.ChainConfig) (*big.Int, *big.Int, *big.Int) {
	rewardConfig := config.Governance.Reward

	if !rewardConfig.DeferredTxFee {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}

	var totalFee, rewardFee, burntFee *big.Int
	totalGasUsed := new(big.Int).SetUint64(header.GasUsed)

	if config.IsMagmaForkEnabled(header.Number) {
		totalFee = new(big.Int).Mul(totalGasUsed, header.BaseFee)
	} else {
		unitPrice := new(big.Int).SetUint64(config.UnitPrice)
		totalFee = new(big.Int).Mul(totalGasUsed, unitPrice)
	}

	rewardFee = new(big.Int).Add(big.NewInt(0), totalFee)
	burntFee = big.NewInt(0)
	if config.IsMagmaForkEnabled(header.Number) {
		halfFee := new(big.Int).Div(rewardFee, big.NewInt(2))
		rewardFee = rewardFee.Sub(rewardFee, halfFee)
		burntFee = burntFee.Add(burntFee, halfFee)
	}

	if config.IsKoreForkEnabled(header.Number) {
		minted := rewardConfig.MintingAmount
		cnRatio, _, _, totalRatio, _ := parseRewardRatio(rewardConfig.Ratio)
		cnMinted := new(big.Int).Mul(minted, big.NewInt(int64(cnRatio)))
		cnMinted = cnMinted.Div(cnMinted, big.NewInt(int64(totalRatio)))

		cnMintedBasicRatio, _, cnMintedTotalRatio, _ := parseRewardKip82Ratio(rewardConfig.Kip82Ratio)
		basicReward := new(big.Int).Mul(cnMinted, big.NewInt(int64(cnMintedBasicRatio)))
		basicReward = basicReward.Div(basicReward, big.NewInt(int64(cnMintedTotalRatio)))

		burntKip82 := new(big.Int)

		if rewardFee.Cmp(basicReward) < 0 {
			burntKip82.Set(rewardFee)
		} else {
			burntKip82.Set(basicReward)
		}

		rewardFee = rewardFee.Sub(rewardFee, burntKip82)
		burntFee = burntFee.Add(burntFee, burntKip82)
	}

	return totalFee, rewardFee, burntFee
}

// calcSplit splits fee into (proposer, stakers, kgf, kir, reamining)
func calcSplit(header *types.Header, config *params.ChainConfig, fee *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {
	rewardConfig := config.Governance.Reward

	var proposer, stakers, kgf, kir, remaining *big.Int

	if config.IsKoreForkEnabled(header.Number) {
		resource := new(big.Int)
		resource.Set(rewardConfig.MintingAmount)
		cnRatio, kgfRatio, kirRatio, totalRatio, _ := parseRewardRatio(rewardConfig.Ratio)
		cn := new(big.Int).Mul(resource, big.NewInt(int64(cnRatio)))
		cn = cn.Div(cn, big.NewInt(int64(totalRatio)))

		kgf = new(big.Int).Mul(resource, big.NewInt(int64(kgfRatio)))
		kgf = kgf.Div(kgf, big.NewInt(int64(totalRatio)))

		kir = new(big.Int).Mul(resource, big.NewInt(int64(kirRatio)))
		kir = kir.Div(kir, big.NewInt(int64(totalRatio)))

		cnMintedBasicRatio, cnMintedStakeRatio, cnMintedTotalRatio, _ := parseRewardKip82Ratio(rewardConfig.Kip82Ratio)
		cnBasic := new(big.Int).Mul(cn, big.NewInt(int64(cnMintedBasicRatio)))
		cnBasic = cnBasic.Div(cnBasic, big.NewInt(int64(cnMintedTotalRatio)))

		cnStake := new(big.Int).Mul(cn, big.NewInt(int64(cnMintedStakeRatio)))
		cnStake = cnStake.Div(cnStake, big.NewInt(int64(cnMintedTotalRatio)))

		remaining = new(big.Int)
		remaining.Set(resource)
		remaining = remaining.Sub(remaining, kgf)
		remaining = remaining.Sub(remaining, kir)
		remaining = remaining.Sub(remaining, cnBasic)
		remaining = remaining.Sub(remaining, cnStake)

		proposer = new(big.Int).Add(cnBasic, fee)
		stakers = cnStake
	} else {
		resource := new(big.Int).Add(rewardConfig.MintingAmount, fee)
		cnRatio, kgfRatio, kirRatio, totalRatio, _ := parseRewardRatio(rewardConfig.Ratio)
		cn := new(big.Int).Mul(resource, big.NewInt(int64(cnRatio)))
		cn = cn.Div(cn, big.NewInt(int64(totalRatio)))

		kgf = new(big.Int).Mul(resource, big.NewInt(int64(kgfRatio)))
		kgf = kgf.Div(kgf, big.NewInt(int64(totalRatio)))

		kir = new(big.Int).Mul(resource, big.NewInt(int64(kirRatio)))
		kir = kir.Div(kir, big.NewInt(int64(totalRatio)))

		remaining = new(big.Int)
		remaining.Set(resource)
		remaining = remaining.Sub(remaining, kgf)
		remaining = remaining.Sub(remaining, kir)
		remaining = remaining.Sub(remaining, cn)

		proposer = cn
		stakers = big.NewInt(0)
	}

	return proposer, stakers, kgf, kir, remaining
}

// calcShares distributes stake reward among staked CNs
func calcShares(rewardConfig *params.RewardConfig, stakingInfo *StakingInfo, stakeReward *big.Int) (map[common.Address]*big.Int, *big.Int) {
	// if stakingInfo is nil, stakeReward goes to proposer
	if stakingInfo == nil {
		return make(map[common.Address]*big.Int), stakeReward
	}

	cns := stakingInfo.GetConsolidatedStakingInfo()

	minStakeInt := rewardConfig.MinimumStake.Uint64()
	totalStakesInt := uint64(0)

	for _, node := range cns.GetAllNodes() {
		if node.StakingAmount > minStakeInt { // comparison in Klay
			totalStakesInt += (node.StakingAmount - minStakeInt)
		}
	}

	totalStakes := new(big.Int).SetUint64(totalStakesInt)
	remaining := new(big.Int).Set(stakeReward)
	shares := make(map[common.Address]*big.Int)

	for _, node := range cns.GetAllNodes() {
		if node.StakingAmount > minStakeInt {
			effectiveStake := new(big.Int).SetUint64(node.StakingAmount - minStakeInt)
			// effectiveStake, totalStakes are in Klay, but will cancel out
			rewardAmount := new(big.Int).Mul(stakeReward, effectiveStake)
			rewardAmount = rewardAmount.Div(rewardAmount, totalStakes)
			remaining = remaining.Sub(remaining, rewardAmount)
			if rewardAmount.Cmp(big.NewInt(0)) > 0 {
				shares[node.RewardAddr] = rewardAmount
			}
		}
	}
	logger.Debug("calcShares", "minStakeInt", minStakeInt,
		"stakeReward", stakeReward.Uint64(),
		"remaining", remaining.Uint64(),
		"shares", shares)

	return shares, remaining
}

func parseRewardRatio(ratio string) (int, int, int, int, error) {
	s := strings.Split(ratio, "/")
	if len(s) != params.RewardSliceCount {
		return 0, 0, 0, 0, errInvalidFormat
	}
	cn, err1 := strconv.Atoi(s[0])
	poc, err2 := strconv.Atoi(s[1])
	kir, err3 := strconv.Atoi(s[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, 0, errParsingRatio
	}
	return cn, poc, kir, cn + poc + kir, nil
}

func parseRewardKip82Ratio(ratio string) (int, int, int, error) {
	s := strings.Split(ratio, "/")
	if len(s) != params.RewardKip82SliceCount {
		return 0, 0, 0, errInvalidFormat
	}
	basic, err1 := strconv.Atoi(s[0])
	stake, err2 := strconv.Atoi(s[1])

	if err1 != nil || err2 != nil {
		return 0, 0, 0, errParsingRatio
	}
	return basic, stake, basic + stake, nil
}

func increment(m map[common.Address]*big.Int, addr common.Address, amount *big.Int) {
	_, ok := m[addr]
	if !ok {
		m[addr] = big.NewInt(0)
	}

	m[addr] = m[addr].Add(m[addr], amount)
}
