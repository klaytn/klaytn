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

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

var logger = log.NewModuleLogger(log.Reward)

type BalanceAdder interface {
	AddBalance(addr common.Address, v *big.Int)
}

type governanceHelper interface {
	Epoch() uint64
	GetItemAtNumberByIntKey(num uint64, key int) (interface{}, error)
	DeferredTxFee() bool
	ProposerPolicy() uint64
	StakingUpdateInterval() uint64
}

type RewardDistributor struct {
	rcc *rewardConfigCache
	gh  governanceHelper
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
	totalTxFee := totalGasUsed.Mul(totalGasUsed, rewardConfig.unitPrice)
	return totalTxFee
}

// MintKLAY mints KLAY and gives the KLAY and the total transaction gas fee to the block proposer.
func (rd *RewardDistributor) MintKLAY(b BalanceAdder, header *types.Header) error {
	rewardConfig, err := rd.rcc.get(header.Number.Uint64())
	if err != nil {
		return err
	}

	totalTxFee := rd.getTotalTxFee(header, rewardConfig)
	blockReward := totalTxFee.Add(rewardConfig.mintingAmount, totalTxFee)

	b.AddBalance(header.Rewardbase, blockReward)
	return nil
}

// DistributeBlockReward distributes block reward to proposer, kirAddr and pocAddr.
func (rd *RewardDistributor) DistributeBlockReward(b BalanceAdder, header *types.Header, pocAddr common.Address, kirAddr common.Address) error {
	rewardConfig, err := rd.rcc.get(header.Number.Uint64())
	if err != nil {
		return err
	}

	// Calculate total tx fee
	totalTxFee := common.Big0
	if rd.gh.DeferredTxFee() {
		totalTxFee = rd.getTotalTxFee(header, rewardConfig)
	}

	rd.distributeBlockReward(b, header, totalTxFee, rewardConfig, pocAddr, kirAddr)
	return nil
}

// distributeBlockReward mints KLAY and distributes newly minted KLAY and transaction fee to proposer, kirAddr and pocAddr.
func (rd *RewardDistributor) distributeBlockReward(b BalanceAdder, header *types.Header, totalTxFee *big.Int, rewardConfig *rewardConfig, pocAddr common.Address, kirAddr common.Address) {
	proposer := header.Rewardbase
	// Block reward
	blockReward := big.NewInt(0).Add(rewardConfig.mintingAmount, totalTxFee)

	tmpInt := big.NewInt(0)

	tmpInt = tmpInt.Mul(blockReward, rewardConfig.cnRatio)
	cnReward := big.NewInt(0).Div(tmpInt, rewardConfig.totalRatio)

	tmpInt = tmpInt.Mul(blockReward, rewardConfig.pocRatio)
	pocIncentive := big.NewInt(0).Div(tmpInt, rewardConfig.totalRatio)

	tmpInt = tmpInt.Mul(blockReward, rewardConfig.kirRatio)
	kirIncentive := big.NewInt(0).Div(tmpInt, rewardConfig.totalRatio)

	remaining := tmpInt.Sub(blockReward, cnReward)
	remaining = tmpInt.Sub(remaining, pocIncentive)
	remaining = tmpInt.Sub(remaining, kirIncentive)
	pocIncentive = pocIncentive.Add(pocIncentive, remaining)

	// CN reward
	b.AddBalance(proposer, cnReward)

	// Proposer gets PoC incentive and KIR incentive, if there is no PoC/KIR address.
	// PoC
	if common.EmptyAddress(pocAddr) {
		pocAddr = proposer
	}
	b.AddBalance(pocAddr, pocIncentive)

	// KIR
	if common.EmptyAddress(kirAddr) {
		kirAddr = proposer
	}
	b.AddBalance(kirAddr, kirIncentive)

	logger.Debug("Block reward", "blockNumber", header.Number.Uint64(),
		"Reward address of a proposer", proposer, "CN reward amount", cnReward,
		"PoC address", pocAddr, "Poc incentive", pocIncentive,
		"KIR address", kirAddr, "KIR incentive", kirIncentive)
}
