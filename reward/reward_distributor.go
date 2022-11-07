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
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
)

var CalcDeferredRewardTimer time.Duration

var logger = log.NewModuleLogger(log.Reward)

var (
	errInvalidFormat = errors.New("invalid ratio format")
	errParsingRatio  = errors.New("parsing ratio fail")
)

type BalanceAdder interface {
	AddBalance(addr common.Address, v *big.Int)
}

// Cannot use governance.Engine because of cyclic dependency.
// Instead declare only the methods used by this package.
type governanceHelper interface {
	Params() *params.GovParamSet
	ParamsAt(num uint64) (*params.GovParamSet, error)
}

type rewardConfig struct {
	// hardfork rules
	rules params.Rules

	// values calculated from block header
	totalFee *big.Int

	// values from ChainConfig
	mintingAmount *big.Int
	minimumStake  *big.Int
	deferredTxFee bool

	// parsed ratio
	cnRatio    *big.Int
	kgfRatio   *big.Int
	kirRatio   *big.Int
	totalRatio *big.Int

	// parsed KIP82 ratio
	cnProposerRatio *big.Int
	cnStakingRatio  *big.Int
	cnTotalRatio    *big.Int
}

type RewardSpec struct {
	Minted   *big.Int                    // the amount newly minted
	Fee      *big.Int                    // total tx fee spent
	Burnt    *big.Int                    // the amount burnt
	Proposer *big.Int                    // the amount allocated to the block proposer
	Stakers  *big.Int                    // total amount allocated to stakers
	Kgf      *big.Int                    // the amount allocated to KGF
	Kir      *big.Int                    // the amount allocated to KIR
	Rewards  map[common.Address]*big.Int // mapping from reward recipient to amounts
}

// TODO: this is for legacy, will be removed
type RewardDistributor struct{}

func NewRewardDistributor(gh governanceHelper) *RewardDistributor {
	return &RewardDistributor{}
}

// DistributeBlockReward distributes a given block's reward at the end of block processing
func DistributeBlockReward(b BalanceAdder, rewards map[common.Address]*big.Int) {
	for addr, amount := range rewards {
		b.AddBalance(addr, amount)
	}
}

func NewRewardConfig(header *types.Header, config *params.ChainConfig) (*rewardConfig, error) {
	if config.Governance == nil {
		return nil, errors.New("no config.Governance")
	}
	if config.Governance.Reward == nil {
		return nil, errors.New("no config.Governance.Reward")
	}

	rules := config.Rules(header.Number)

	cnRatio, kgfRatio, kirRatio, totalRatio, err := parseRewardRatio(config.Governance.Reward.Ratio)
	if err != nil {
		return nil, err
	}

	cnProposerRatio, cnStakingRatio, cnTotalRatio, err := parseRewardKip82Ratio(config.Governance.Reward.Kip82Ratio)
	if err != nil {
		return nil, err
	}

	return &rewardConfig{
		// hardfork rules
		rules: rules,

		// values calculated from block header
		totalFee: GetTotalTxFee(header, config),

		// values from ChainConfig
		mintingAmount: new(big.Int).Set(config.Governance.Reward.MintingAmount),
		minimumStake:  new(big.Int).Set(config.Governance.Reward.MinimumStake),
		deferredTxFee: config.Governance.Reward.DeferredTxFee,

		// parsed ratio
		cnRatio:    big.NewInt(int64(cnRatio)),
		kgfRatio:   big.NewInt(int64(kgfRatio)),
		kirRatio:   big.NewInt(int64(kirRatio)),
		totalRatio: big.NewInt(int64(totalRatio)),

		// parsed KIP82 ratio
		cnProposerRatio: big.NewInt(int64(cnProposerRatio)),
		cnStakingRatio:  big.NewInt(int64(cnStakingRatio)),
		cnTotalRatio:    big.NewInt(int64(cnTotalRatio)),
	}, nil
}

func GetTotalTxFee(header *types.Header, config *params.ChainConfig) *big.Int {
	totalFee := new(big.Int).SetUint64(header.GasUsed)
	if config.IsMagmaForkEnabled(header.Number) {
		totalFee = totalFee.Mul(totalFee, header.BaseFee)
	} else {
		totalFee = totalFee.Mul(totalFee, new(big.Int).SetUint64(config.UnitPrice))
	}
	return totalFee
}

// GetBlockReward returns the actual reward amounts paid in this block
// Used in klay_getReward RPC API
func GetBlockReward(header *types.Header, config *params.ChainConfig) (*RewardSpec, error) {
	var spec *RewardSpec
	var err error

	if config.Istanbul == nil {
		return nil, errors.New("no IstanbulConfig")
	}

	policy := config.Istanbul.ProposerPolicy
	if policy == uint64(istanbul.RoundRobin) || policy == uint64(istanbul.Sticky) {
		spec, err = CalcDeferredRewardSimple(header, config)
		if err != nil {
			return nil, err
		}
	} else {
		spec, err = CalcDeferredReward(header, config)
		if err != nil {
			return nil, err
		}

		// Compensate the difference between CalcDeferredReward() and actual payment.
		// If not DeferredTxFee, CalcDeferredReward() assumes 0 total_fee, but
		// some non-zero fee is paid to the proposer.
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

// CalcDeferredRewardSimple distributes rewards to proposer after optional fee burning
// this behaves similar to the previous MintKLAY
func CalcDeferredRewardSimple(header *types.Header, config *params.ChainConfig) (*RewardSpec, error) {
	rc, err := NewRewardConfig(header, config)
	if err != nil {
		return nil, err
	}

	minted := rc.mintingAmount
	totalFee := rc.totalFee
	rewardFee := new(big.Int).Set(totalFee)
	burntFee := big.NewInt(0)

	if rc.rules.IsMagma {
		burnt := getBurnAmountMagma(rewardFee)
		rewardFee = rewardFee.Sub(rewardFee, burnt)
		burntFee = burntFee.Add(burntFee, burnt)
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
	defer func(start time.Time) {
		CalcDeferredRewardTimer = time.Since(start)
	}(time.Now())

	rc, err := NewRewardConfig(header, config)
	if err != nil {
		return nil, err
	}

	var (
		minted      = rc.mintingAmount
		stakingInfo = GetStakingInfo(header.Number.Uint64())
	)

	totalFee, rewardFee, burntFee := calcDeferredFee(rc)
	proposer, stakers, kgf, kir, splitRem := calcSplit(rc, minted, rewardFee)
	shares, shareRem := calcShares(stakingInfo, stakers, rc.minimumStake.Uint64())

	// Remainder from (CN, KGF, KIR) split goes to KGF
	kgf = kgf.Add(kgf, splitRem)
	// Remainder from staker shares goes to Proposer
	proposer = proposer.Add(proposer, shareRem)

	// if KGF or KIR is not set, proposer gets the portion
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
		Minted:   minted,
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
	logger.Debug("CalcDeferredReward returns", "spec", spec)

	return spec, nil
}

// calcDeferredFee splits fee into (total, reward, burnt)
func calcDeferredFee(rc *rewardConfig) (*big.Int, *big.Int, *big.Int) {
	// If not DeferredTxFee, fees are already added to the proposer during TX execution.
	// Therefore, there are no fees to distribute here at the end of block processing.
	// However, the fees must be compensated to calculate actual rewards paid.
	if !rc.deferredTxFee {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}

	totalFee := rc.totalFee
	rewardFee := new(big.Int).Set(totalFee)
	burntFee := big.NewInt(0)

	// after magma, burn half of gas
	if rc.rules.IsMagma {
		burnt := getBurnAmountMagma(rewardFee)
		rewardFee = rewardFee.Sub(rewardFee, burnt)
		burntFee = burntFee.Add(burntFee, burnt)
	}

	// after kore, burn fees up to proposer's minted reward
	if rc.rules.IsKore {
		burnt := getBurnAmountKore(rc, rewardFee)
		rewardFee = rewardFee.Sub(rewardFee, burnt)
		burntFee = burntFee.Add(burntFee, burnt)
	}

	logger.Debug("calcDeferredFee returns",
		"totalFee", totalFee.Uint64(),
		"rewardFee", rewardFee.Uint64(),
		"burntFee", burntFee.Uint64(),
	)
	return totalFee, rewardFee, burntFee
}

func getBurnAmountMagma(fee *big.Int) *big.Int {
	return new(big.Int).Div(fee, big.NewInt(2))
}

func getBurnAmountKore(rc *rewardConfig, fee *big.Int) *big.Int {
	cn, _, _ := splitByRatio(rc, rc.mintingAmount)
	proposer, _ := splitByKip82Ratio(rc, cn)

	logger.Debug("getBurnAmountKore returns",
		"fee", fee.Uint64(),
		"proposer", proposer.Uint64(),
	)

	if fee.Cmp(proposer) >= 0 {
		return proposer
	} else {
		return new(big.Int).Set(fee) // return copy of the parameter
	}
}

// calcSplit splits fee into (proposer, stakers, kgf, kir, reamining)
// the sum of the output must be equal to (minted + fee)
func calcSplit(rc *rewardConfig, minted, fee *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {
	totalResource := big.NewInt(0)
	totalResource = totalResource.Add(minted, fee)

	if rc.rules.IsKore {
		cn, kgf, kir := splitByRatio(rc, minted)
		proposer, stakers := splitByKip82Ratio(rc, cn)

		proposer = proposer.Add(proposer, fee)

		remaining := new(big.Int).Set(totalResource)
		remaining = remaining.Sub(remaining, kgf)
		remaining = remaining.Sub(remaining, kir)
		remaining = remaining.Sub(remaining, proposer)
		remaining = remaining.Sub(remaining, stakers)

		logger.Debug("calcSplit after kore returns",
			"proposer", proposer.Uint64(),
			"stakers", stakers.Uint64(),
			"kgf", kgf.Uint64(),
			"kir", kir.Uint64(),
			"remaining", remaining.Uint64(),
		)
		return proposer, stakers, kgf, kir, remaining
	} else {
		source := big.NewInt(0)
		source = source.Add(minted, fee)
		cn, kgf, kir := splitByRatio(rc, source)

		remaining := new(big.Int).Set(totalResource)
		remaining = remaining.Sub(remaining, kgf)
		remaining = remaining.Sub(remaining, kir)
		remaining = remaining.Sub(remaining, cn)

		logger.Debug("calcSplit before kore returns",
			"cn", cn.Uint64(),
			"kgf", kgf.Uint64(),
			"kir", kir.Uint64(),
			"remaining", remaining.Uint64(),
		)
		return cn, big.NewInt(0), kgf, kir, remaining
	}
}

// splitByRatio splits by `ratio`. It ignores any remaining amounts.
func splitByRatio(rc *rewardConfig, source *big.Int) (*big.Int, *big.Int, *big.Int) {
	cn := new(big.Int).Mul(source, rc.cnRatio)
	cn = cn.Div(cn, rc.totalRatio)

	kgf := new(big.Int).Mul(source, rc.kgfRatio)
	kgf = kgf.Div(kgf, rc.totalRatio)

	kir := new(big.Int).Mul(source, rc.kirRatio)
	kir = kir.Div(kir, rc.totalRatio)

	return cn, kgf, kir
}

// splitByKip82Ratio splits by `kip82ratio`. It ignores any remaining amounts.
func splitByKip82Ratio(rc *rewardConfig, source *big.Int) (*big.Int, *big.Int) {
	proposer := new(big.Int).Mul(source, rc.cnProposerRatio)
	proposer = proposer.Div(proposer, rc.cnTotalRatio)

	stakers := new(big.Int).Mul(source, rc.cnStakingRatio)
	stakers = stakers.Div(stakers, rc.cnTotalRatio)

	return proposer, stakers
}

// calcShares distributes stake reward among staked CNs
func calcShares(stakingInfo *StakingInfo, stakeReward *big.Int, minStake uint64) (map[common.Address]*big.Int, *big.Int) {
	// if stakingInfo is nil, stakeReward goes to proposer
	if stakingInfo == nil {
		return make(map[common.Address]*big.Int), stakeReward
	}

	cns := stakingInfo.GetConsolidatedStakingInfo()

	totalStakesInt := uint64(0)

	for _, node := range cns.GetAllNodes() {
		if node.StakingAmount > minStake { // comparison in Klay
			totalStakesInt += (node.StakingAmount - minStake)
		}
	}

	totalStakes := new(big.Int).SetUint64(totalStakesInt)
	remaining := new(big.Int).Set(stakeReward)
	shares := make(map[common.Address]*big.Int)

	for _, node := range cns.GetAllNodes() {
		if node.StakingAmount > minStake {
			effectiveStake := new(big.Int).SetUint64(node.StakingAmount - minStake)
			// effectiveStake, totalStakes are in Klay, but will cancel out
			rewardAmount := new(big.Int).Mul(stakeReward, effectiveStake)
			rewardAmount = rewardAmount.Div(rewardAmount, totalStakes)
			remaining = remaining.Sub(remaining, rewardAmount)
			if rewardAmount.Cmp(big.NewInt(0)) > 0 {
				shares[node.RewardAddr] = rewardAmount
			}
		}
	}
	logger.Debug("calcShares returns",
		"stakeReward", stakeReward.Uint64(),
		"remaining", remaining.Uint64(),
		"shares", shares)

	return shares, remaining
}

// parseRewardRatio parses string `ratio` into ints
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

// parseRewardKip82Ratio parses string `kip82ratio` into ints
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
