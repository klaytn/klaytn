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
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
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
	CurrentParams() *params.GovParamSet
	EffectiveParams(num uint64) (*params.GovParamSet, error)
}

type rewardConfig struct {
	// hardfork rules
	rules params.Rules

	// values calculated from block header
	totalFee *big.Int

	// values from GovParamSet
	mintingAmount *big.Int
	minimumStake  *big.Int
	deferredTxFee bool

	// parsed ratio
	cnRatio    *big.Int
	kffRatio   *big.Int
	kcfRatio   *big.Int
	totalRatio *big.Int

	// parsed KIP82 ratio
	cnProposerRatio *big.Int
	cnStakingRatio  *big.Int
	cnTotalRatio    *big.Int
}

type RewardSpec struct {
	Minted   *big.Int                    `json:"minted"`   // the amount newly minted
	TotalFee *big.Int                    `json:"totalFee"` // total tx fee spent
	BurntFee *big.Int                    `json:"burntFee"` // the amount burnt
	Proposer *big.Int                    `json:"proposer"` // the amount allocated to the block proposer
	Stakers  *big.Int                    `json:"stakers"`  // total amount allocated to stakers
	KFF      *big.Int                    `json:"kff"`      // the amount allocated to KFF
	KCF      *big.Int                    `json:"kcf"`      // the amount allocated to KCF
	Rewards  map[common.Address]*big.Int `json:"rewards"`  // mapping from reward recipient to amounts
}

func NewRewardSpec() *RewardSpec {
	return &RewardSpec{
		Minted:   big.NewInt(0),
		TotalFee: big.NewInt(0),
		BurntFee: big.NewInt(0),
		Proposer: big.NewInt(0),
		Stakers:  big.NewInt(0),
		KFF:      big.NewInt(0),
		KCF:      big.NewInt(0),
		Rewards:  make(map[common.Address]*big.Int),
	}
}

func (spec *RewardSpec) Add(delta *RewardSpec) {
	spec.Minted.Add(spec.Minted, delta.Minted)
	spec.TotalFee.Add(spec.TotalFee, delta.TotalFee)
	spec.BurntFee.Add(spec.BurntFee, delta.BurntFee)
	spec.Proposer.Add(spec.Proposer, delta.Proposer)
	spec.Stakers.Add(spec.Stakers, delta.Stakers)
	spec.KFF.Add(spec.KFF, delta.KFF)
	spec.KCF.Add(spec.KCF, delta.KCF)

	for addr, amount := range delta.Rewards {
		incrementRewardsMap(spec.Rewards, addr, amount)
	}
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

func NewRewardConfig(header *types.Header, rules params.Rules, pset *params.GovParamSet) (*rewardConfig, error) {
	cnRatio, kffRatio, kcfRatio, totalRatio, err := parseRewardRatio(pset.Ratio())
	if err != nil {
		return nil, err
	}

	var cnProposerRatio, cnStakingRatio, cnTotalRatio int64
	if rules.IsKore {
		cnProposerRatio, cnStakingRatio, cnTotalRatio, err = parseRewardKip82Ratio(pset.Kip82Ratio())
		if err != nil {
			return nil, err
		}
	}

	return &rewardConfig{
		// hardfork rules
		rules: rules,

		// values calculated from block header
		totalFee: GetTotalTxFee(header, rules, pset),

		// values from GovParamSet
		mintingAmount: new(big.Int).Set(pset.MintingAmountBig()),
		minimumStake:  new(big.Int).Set(pset.MinimumStakeBig()),
		deferredTxFee: pset.DeferredTxFee(),

		// parsed ratio
		cnRatio:    big.NewInt(cnRatio),
		kffRatio:   big.NewInt(kffRatio),
		kcfRatio:   big.NewInt(kcfRatio),
		totalRatio: big.NewInt(totalRatio),

		// parsed KIP82 ratio
		cnProposerRatio: big.NewInt(cnProposerRatio),
		cnStakingRatio:  big.NewInt(cnStakingRatio),
		cnTotalRatio:    big.NewInt(cnTotalRatio),
	}, nil
}

func GetTotalTxFee(header *types.Header, rules params.Rules, pset *params.GovParamSet) *big.Int {
	totalFee := new(big.Int).SetUint64(header.GasUsed)
	if rules.IsMagma {
		totalFee = totalFee.Mul(totalFee, header.BaseFee)
	} else {
		totalFee = totalFee.Mul(totalFee, new(big.Int).SetUint64(pset.UnitPrice()))
	}
	return totalFee
}

// config.Istanbul must have been set
func IsRewardSimple(pset *params.GovParamSet) bool {
	return pset.Policy() != uint64(istanbul.WeightedRandom)
}

// CalcRewardParamBlock returns the block number with which governance parameters must be fetched
// This mimics the legacy reward config cache before Kore
func CalcRewardParamBlock(num, epoch uint64, rules params.Rules) uint64 {
	if !rules.IsKore && num%epoch == 0 {
		return num - epoch
	}
	return num
}

// GetBlockReward returns the actual reward amounts paid in this block
// Used in klay_getReward RPC API
func GetBlockReward(header *types.Header, rules params.Rules, pset *params.GovParamSet) (*RewardSpec, error) {
	var spec *RewardSpec
	var err error

	if IsRewardSimple(pset) {
		spec, err = CalcDeferredRewardSimple(header, rules, pset)
		if err != nil {
			return nil, err
		}
	} else {
		spec, err = CalcDeferredReward(header, rules, pset)
		if err != nil {
			return nil, err
		}
	}

	// Compensate the difference between CalcDeferredReward() and actual payment.
	// If not DeferredTxFee, CalcDeferredReward() assumes 0 total_fee, but
	// some non-zero fee already has been paid to the proposer.
	if !pset.DeferredTxFee() {
		if rules.IsMagma {
			txFee := GetTotalTxFee(header, rules, pset)
			txFeeBurn := getBurnAmountMagma(txFee)
			txFeeRemained := new(big.Int).Sub(txFee, txFeeBurn)
			spec.BurntFee = txFeeBurn

			spec.Proposer = spec.Proposer.Add(spec.Proposer, txFeeRemained)
			spec.TotalFee = spec.TotalFee.Add(spec.TotalFee, txFee)
			incrementRewardsMap(spec.Rewards, header.Rewardbase, txFeeRemained)
		} else if rules.IsServiceChainRewardFix {
			txFee := GetTotalTxFee(header, rules, pset)
			spec.Proposer = spec.Proposer.Add(spec.Proposer, txFee)
			spec.TotalFee = spec.TotalFee.Add(spec.TotalFee, txFee)
			incrementRewardsMap(spec.Rewards, header.Rewardbase, txFee)
		} else {
			txFee := GetTotalTxFee(header, rules, pset)
			spec.Proposer = spec.Proposer.Add(spec.Proposer, txFee)
			spec.TotalFee = spec.TotalFee.Add(spec.TotalFee, txFee)
			// get the proposer of this block.
			proposer, err := ecrecover(header)
			if err != nil {
				return nil, err
			}
			incrementRewardsMap(spec.Rewards, proposer, txFee)

		}
	}

	return spec, nil
}

// CalcDeferredRewardSimple distributes rewards to proposer after optional fee burning
// this behaves similar to the previous MintKLAY
// MintKLAY has been superseded because we need to split reward distribution
// logic into (1) calculation, and (2) actual distribution.
// CalcDeferredRewardSimple does the former and DistributeBlockReward does the latter
func CalcDeferredRewardSimple(header *types.Header, rules params.Rules, pset *params.GovParamSet) (*RewardSpec, error) {
	rc, err := NewRewardConfig(header, rules, pset)
	if err != nil {
		return nil, err
	}

	minted := rc.mintingAmount

	// If not DeferredTxFee, fees are already added to the proposer during TX execution.
	// Therefore, there are no fees to distribute here at the end of block processing.
	// However, there was a bug causing the fees to be distributed again.
	// See https://github.com/klaytn/klaytn/issues/1692.
	// To maintain backward compatibility, we only fix the buggy logic after Magma
	// and leave the buggy logic before Magma.
	// However, if you wish to address only this issue instead, apply the ServiceChainRewardFixHardfork.
	// Before the hardforks, tx fees must be compensated to calculate actual rewards paid.
	if !rc.deferredTxFee && (rc.rules.IsMagma || rc.rules.IsServiceChainRewardFix) {
		proposer := new(big.Int).Set(minted)
		spec := NewRewardSpec()
		spec.Minted = minted
		spec.TotalFee = big.NewInt(0)
		spec.BurntFee = big.NewInt(0)
		spec.Proposer = proposer
		incrementRewardsMap(spec.Rewards, header.Rewardbase, proposer)
		return spec, nil
	}

	totalFee := rc.totalFee
	rewardFee := new(big.Int).Set(totalFee)
	burntFee := big.NewInt(0)

	if rc.rules.IsMagma {
		burnt := getBurnAmountMagma(rewardFee)
		rewardFee = rewardFee.Sub(rewardFee, burnt)
		burntFee = burntFee.Add(burntFee, burnt)
	}

	proposer := big.NewInt(0).Add(minted, rewardFee)

	logger.Debug("CalcDeferredRewardSimple returns",
		"proposer", proposer.Uint64(),
		"totalFee", totalFee.Uint64(),
		"burntFee", burntFee.Uint64(),
	)

	spec := NewRewardSpec()
	spec.Minted = minted
	spec.TotalFee = totalFee
	spec.BurntFee = burntFee
	spec.Proposer = proposer
	incrementRewardsMap(spec.Rewards, header.Rewardbase, proposer)
	return spec, nil
}

// CalcDeferredReward calculates the deferred rewards,
// which are determined at the end of block processing.
func CalcDeferredReward(header *types.Header, rules params.Rules, pset *params.GovParamSet) (*RewardSpec, error) {
	defer func(start time.Time) {
		CalcDeferredRewardTimer = time.Since(start)
	}(time.Now())

	rc, err := NewRewardConfig(header, rules, pset)
	if err != nil {
		return nil, err
	}

	var (
		minted      = rc.mintingAmount
		stakingInfo = GetStakingInfo(header.Number.Uint64())
	)

	totalFee, rewardFee, burntFee := calcDeferredFee(rc)
	proposer, stakers, kff, kcf, splitRem := calcSplit(rc, minted, rewardFee)
	shares, shareRem := calcShares(stakingInfo, stakers, rc.minimumStake.Uint64())

	// Remainder from (CN, KFF, KCF) split goes to KFF
	kff = kff.Add(kff, splitRem)
	// Remainder from staker shares goes to Proposer
	// Then, deduct it from stakers so that `minted + totalFee - burntFee = proposer + stakers + kff + kcf`
	proposer = proposer.Add(proposer, shareRem)
	stakers = stakers.Sub(stakers, shareRem)

	// if KFF or KCF is not set, proposer gets the portion
	if stakingInfo == nil || common.EmptyAddress(stakingInfo.KFFAddr) {
		logger.Debug("KFF empty, proposer gets its portion", "kff", kff)
		proposer = proposer.Add(proposer, kff)
		kff = big.NewInt(0)
	}
	if stakingInfo == nil || common.EmptyAddress(stakingInfo.KCFAddr) {
		logger.Debug("KCF empty, proposer gets its portion", "kcf", kcf)
		proposer = proposer.Add(proposer, kcf)
		kcf = big.NewInt(0)
	}

	spec := NewRewardSpec()
	spec.Minted = minted
	spec.TotalFee = totalFee
	spec.BurntFee = burntFee
	spec.Proposer = proposer
	spec.Stakers = stakers
	spec.KFF = kff
	spec.KCF = kcf

	incrementRewardsMap(spec.Rewards, header.Rewardbase, proposer)

	if stakingInfo != nil && !common.EmptyAddress(stakingInfo.KFFAddr) {
		incrementRewardsMap(spec.Rewards, stakingInfo.KFFAddr, kff)
	}
	if stakingInfo != nil && !common.EmptyAddress(stakingInfo.KCFAddr) {
		incrementRewardsMap(spec.Rewards, stakingInfo.KCFAddr, kcf)
	}

	for rewardAddr, rewardAmount := range shares {
		incrementRewardsMap(spec.Rewards, rewardAddr, rewardAmount)
	}
	logger.Debug("CalcDeferredReward() returns", "spec", spec)

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

	logger.Debug("calcDeferredFee()",
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

	logger.Debug("getBurnAmountKore()",
		"fee", fee.Uint64(),
		"proposer", proposer.Uint64(),
	)

	if fee.Cmp(proposer) >= 0 {
		return proposer
	} else {
		return new(big.Int).Set(fee) // return copy of the parameter
	}
}

// calcSplit splits fee into (proposer, stakers, kff, kcf, remaining)
// the sum of the output must be equal to (minted + fee)
func calcSplit(rc *rewardConfig, minted, fee *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {
	totalResource := big.NewInt(0)
	totalResource = totalResource.Add(minted, fee)

	if rc.rules.IsKore {
		cn, kff, kcf := splitByRatio(rc, minted)
		proposer, stakers := splitByKip82Ratio(rc, cn)

		proposer = proposer.Add(proposer, fee)

		remaining := new(big.Int).Set(totalResource)
		remaining = remaining.Sub(remaining, kff)
		remaining = remaining.Sub(remaining, kcf)
		remaining = remaining.Sub(remaining, proposer)
		remaining = remaining.Sub(remaining, stakers)

		logger.Debug("calcSplit() after kore",
			"[in] minted", minted.Uint64(),
			"[in] fee", fee.Uint64(),
			"[out] proposer", proposer.Uint64(),
			"[out] stakers", stakers.Uint64(),
			"[out] kff", kff.Uint64(),
			"[out] kcf", kcf.Uint64(),
			"[out] remaining", remaining.Uint64(),
		)
		return proposer, stakers, kff, kcf, remaining
	} else {
		cn, kff, kcf := splitByRatio(rc, totalResource)

		remaining := new(big.Int).Set(totalResource)
		remaining = remaining.Sub(remaining, kff)
		remaining = remaining.Sub(remaining, kcf)
		remaining = remaining.Sub(remaining, cn)

		logger.Debug("calcSplit() before kore",
			"[in] minted", minted.Uint64(),
			"[in] fee", fee.Uint64(),
			"[out] cn", cn.Uint64(),
			"[out] kff", kff.Uint64(),
			"[out] kcf", kcf.Uint64(),
			"[out] remaining", remaining.Uint64(),
		)
		return cn, big.NewInt(0), kff, kcf, remaining
	}
}

// splitByRatio splits by `ratio`. It ignores any remaining amounts.
func splitByRatio(rc *rewardConfig, source *big.Int) (*big.Int, *big.Int, *big.Int) {
	cn := new(big.Int).Mul(source, rc.cnRatio)
	cn = cn.Div(cn, rc.totalRatio)

	kff := new(big.Int).Mul(source, rc.kffRatio)
	kff = kff.Div(kff, rc.totalRatio)

	kcf := new(big.Int).Mul(source, rc.kcfRatio)
	kcf = kcf.Div(kcf, rc.totalRatio)

	return cn, kff, kcf
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
			// The KLAY unit will cancel out:
			// rewardAmount (peb) = stakeReward (peb) * effectiveStake (KLAY) / totalStakes (KLAY)
			rewardAmount := new(big.Int).Mul(stakeReward, effectiveStake)
			rewardAmount = rewardAmount.Div(rewardAmount, totalStakes)
			remaining = remaining.Sub(remaining, rewardAmount)
			if rewardAmount.Cmp(big.NewInt(0)) > 0 {
				shares[node.RewardAddr] = rewardAmount
			}
		}
	}
	logger.Debug("calcShares()",
		"[in] stakeReward", stakeReward.Uint64(),
		"[out] remaining", remaining.Uint64(),
		"[out] shares", shares,
	)

	return shares, remaining
}

// parseRewardRatio parses string `ratio` into ints
func parseRewardRatio(ratio string) (int64, int64, int64, int64, error) {
	s := strings.Split(ratio, "/")
	if len(s) != params.RewardSliceCount {
		logger.Error("Invalid ratio format", "ratio", ratio)
		return 0, 0, 0, 0, errInvalidFormat
	}
	cn, err1 := strconv.ParseInt(s[0], 10, 64)
	kff, err2 := strconv.ParseInt(s[1], 10, 64)
	kcf, err3 := strconv.ParseInt(s[2], 10, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		logger.Error("Could not parse ratio", "ratio", ratio)
		return 0, 0, 0, 0, errParsingRatio
	}
	return cn, kff, kcf, cn + kff + kcf, nil
}

// parseRewardKip82Ratio parses string `kip82ratio` into ints
func parseRewardKip82Ratio(ratio string) (int64, int64, int64, error) {
	s := strings.Split(ratio, "/")
	if len(s) != params.RewardKip82SliceCount {
		logger.Error("Invalid kip82ratio format", "ratio", ratio)
		return 0, 0, 0, errInvalidFormat
	}
	proposer, err1 := strconv.ParseInt(s[0], 10, 64)
	stakers, err2 := strconv.ParseInt(s[1], 10, 64)

	if err1 != nil || err2 != nil {
		logger.Error("Could not parse kip82ratio", "ratio", ratio)
		return 0, 0, 0, errParsingRatio
	}
	return proposer, stakers, proposer + stakers, nil
}

func incrementRewardsMap(m map[common.Address]*big.Int, addr common.Address, amount *big.Int) {
	_, ok := m[addr]
	if !ok {
		m[addr] = big.NewInt(0)
	}

	m[addr] = m[addr].Add(m[addr], amount)
}

// ecrecover extracts the Klaytn account address from a signed header.
func ecrecover(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return common.Address{}, nil
	}

	sigHash, err := sigHash(header)
	if err != nil {
		return common.Address{}, err
	}
	addr, err := istanbul.GetSignatureAddress(sigHash.Bytes(), istanbulExtra.Seal)
	if err != nil {
		return addr, err
	}
	return addr, nil
}

func sigHash(header *types.Header) (hash common.Hash, err error) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	if err := rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false)); err != nil {
		logger.Error("fail to encode", "err", err)
		return common.Hash{}, err
	}
	hasher.Sum(hash[:0])
	return hash, nil
}
