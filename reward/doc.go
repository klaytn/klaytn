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
Package reward implements the Klaytn Reward System.
The Klaytn reward system manages stakingInfo and distributes block rewards.

Managing stakingInfo

Klaytn uses WeightedRandom policy to choose a block proposer.
It means, the percentage of becoming a block proposer depends on how much KLAY a node has staked.
Therefore, a node with stakes more than others will have more opportunities than other nodes.

StakingInfo is a data including stakingAmount and addresses (node, staking, reward, PoC and KIR).
StakingAmount is a balance how much KLAY each node has staked. Only CCO can stake KLAY.
Each CCO stakes KLAY by a smart contract called staking contract.
StakingAmount is derived by checking balance of staking contract.
StakingInfo has 5 types of addresses. All addresses are obtained by the addressBookContract which is pre-deployed in the genesis block.
stakingAddress are addresses of staking contracts of CCO. reward, PoC and KIR addresses are the addresses which get a block reward when a block has been created.
StakingInfo is made every 86400 blocks (stakingInterval) and used in a next interval.

	type StakingInfo struct {
		BlockNum              uint64
		CouncilNodeAddrs      []common.Address // Address of Council
		CouncilStakingAddrs   []common.Address // Address of Staking contract which holds staking balance
		CouncilRewardAddrs    []common.Address // Address of Council account which will get a block reward
		KIRAddr               common.Address   // Address of KIR contract
		PoCAddr               common.Address   // Address of PoC contract
		UseGini               bool             // configure whether Gini is used or not
		Gini                  float64          // Gini coefficient
		CouncilStakingAmounts []uint64         // StakingAmounts of Council. They are derived from Staking addresses of council
	}

StakingInfo is managed by a StakingManager which has a cache for saving StakingInfos.
The StakingManager calculates block number with interval to find a stakingInfo for current block
and returns correct stakingInfo to use.


 related struct
 - RewardDistributor
 - StakingManager
 - addressBookConnector
 - stakingInfoCache
 - stakingInfo


Distributing Reward

Klaytn distributes the reward of a block to proposer, PoC and KIR.
The detail information of PoC and KIR is available on Klaytn docs.

PoC - https://docs.klaytn.com/klaytn/token_economy#proof-of-contribution

KIR - https://docs.klaytn.com/klaytn/token_economy#klaytn-improvement-reserve

Configurations related to the reward system such as mintingAmount, ratio and unitPrice are determined by the Klaytn governance.
All configurations are saved as rewardConfig on every epoch block (default 604,800 blocks) and managed by a rewardConfigCache.

A proposer which has made a current block will get the reward of the block.
A block reward is calculated by following steps.
First, calculate totalReward by adding mintingAmount and totalTxFee (unitPrice * gasUsed).
Second, divide totalReward by ratio (default 34/54/12 - proposer/PoC/KIR).
Last, distribute reward to each address (proposer, PoC, KIR).

 related struct
 - RewardDistributor
 - rewardConfigCache
*/
package reward
