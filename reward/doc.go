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

/*
Reward Package
Package Reward implements the Klaytn Reward System.
Klaytn Reward System manages stakingInfo and distributes block reward.

Managing stakingInfo
Klaytn uses weightedRandom policy to choose a block proposer.
It means, the percentage of becoming a block proposer depends on how much klay a node has staked.
Therefore, a node staked more than others will have more opportunities than other nodes.

StakingInfo is a data including stakingAmount and addresses (node, staking, reward, poc and kir).
stakingAmount is a balance how much each node has staked.
addresses are included stakingAddress which store stakingAmount and reward, poc, kir which get a block reward when a block has been created.
StakingInfo is made every 86400 blocks (stakingInterval) and used in a next interval.

type StakingInfo struct {
	BlockNum              uint64
	CouncilNodeAddrs      []common.Address // Address of Council
	CouncilStakingAddrs   []common.Address // Address of Staking contract which holds staking balance
	CouncilRewardAddrs    []common.Address // Address of Council account which will get block reward
	KIRAddr               common.Address   // Address of KIR contract
	PoCAddr               common.Address   // Address of PoC contract
	UseGini               bool             // configure whether Gini is used or not
	Gini                  float64          // Gini coefficient
	CouncilStakingAmounts []uint64         // Staking amounts of Council. It is derived from Staking address of council
}

StakingInfo is managed by stakingManager which has a cache for saving StakingInfo.
stakingManager calculates block number and interval to find a stakingInfo for current block
and returns correct stakingInfo to use.

related struct
 - RewardManager
 - stakingManager
 - addressBookManager
 - stakingInfoCache
 - stakingInfo

Distributing Reward
Klaytn distributes the reward of a block to proposer, PoC and KIR.

The detail information of PoC and KIR is available on Klaytn docs.
 PoC - https://docs.klaytn.com/klaytn/token_economy#proof-of-contribution
 KIR - https://docs.klaytn.com/klaytn/token_economy#klaytn-improvement-reserve

Configures related to reward system such as mintingAmount, ratio and unitPrice are determined by Klaytn governance.
All configures are saved as rewardConfig every epoch block (default 604800 blocks) and managed by rewardConfigCache.

A proposer which has made a current block will get the reward of the block.
A block reward is calculated by following steps.
 1. calculate totalReward by adding mintingAmount and totalTxFee (unitPrice * gasUsed).
 2. divide totalReward by ratio (default 34/54/12).
 3. distribute reward to each address(proposer, PoC, KIR).

related struct
 - RewardManager
 - rewardConfigCache
*/
package reward
