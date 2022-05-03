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

package params

import (
	"sync/atomic"
)

var (
	stakingUpdateInterval  uint64 = 86400 // About 1 day. 86400 blocks = (24 hrs) * (3600 secs/hr) * (1 block/sec)
	proposerUpdateInterval uint64 = 3600  // About 1 hour. 3600 blocks = (1 hr) * (3600 secs/hr) * (1 block/sec)
)

const (
	// Block reward will be separated by three pieces and distributed
	RewardSliceCount = 3
	// GovernanceConfig is stored in a cache which has below capacity
	GovernanceCacheLimit    = 512
	GovernanceIdxCacheLimit = 1000
	// The prefix for governance cache
	GovernanceCachePrefix = "governance"
)

type EngineType int

const (
	// Engine type
	UseIstanbul EngineType = iota
	UseClique
	Unknown
)

const (
	// Governance Key
	GovernanceMode int = iota
	GoverningNode
	Epoch
	Policy
	CommitteeSize
	UnitPrice
	MintingAmount
	Ratio
	UseGiniCoeff
	DeferredTxFee
	MinimumStake
	AddValidator
	RemoveValidator
	StakeUpdateInterval
	ProposerRefreshInterval
	ConstTxGasHumanReadable
	CliqueEpoch
	Timeout
)

const (
	GovernanceMode_None = iota
	GovernanceMode_Single
	GovernanceMode_Ballot
)

const (
	// Proposer policy
	// At the moment this is duplicated in istanbul/config.go, not to make a cross reference
	// TODO-Klatn-Governance: Find a way to manage below constants at single location
	RoundRobin = iota
	Sticky
	WeightedRandom
)

const (
	// Default Values: Constants used for getting default values for configuration
	DefaultGovernanceMode = "none"
	DefaultGoverningNode  = "0x0000000000000000000000000000000000000000"
	DefaultEpoch          = uint64(604800)
	DefaultProposerPolicy = uint64(0)
	DefaultSubGroupSize   = uint64(21)
	DefaultMintingAmount  = 0
	DefaultRatio          = "100/0/0"
	DefaultUseGiniCoeff   = false
	DefaultDefferedTxFee  = false
	DefaultUnitPrice      = uint64(250000000000)
	DefaultPeriod         = 1
)

func IsStakingUpdateInterval(blockNum uint64) bool {
	return (blockNum % StakingUpdateInterval()) == 0
}

// CalcStakingBlockNumber returns number of block which contains staking information required to make a new block with blockNum.
func CalcStakingBlockNumber(blockNum uint64) uint64 {
	stakingInterval := StakingUpdateInterval()
	if blockNum <= 2*stakingInterval {
		// Just return genesis block number.
		return 0
	}

	var number uint64
	if (blockNum % stakingInterval) == 0 {
		number = blockNum - 2*stakingInterval
	} else {
		number = blockNum - stakingInterval - (blockNum % stakingInterval)
	}
	return number
}

func IsProposerUpdateInterval(blockNum uint64) (bool, uint64) {
	proposerInterval := ProposerUpdateInterval()
	return (blockNum % proposerInterval) == 0, proposerInterval
}

// CalcProposerBlockNumber returns number of block where list of proposers is updated for block blockNum
func CalcProposerBlockNumber(blockNum uint64) uint64 {
	var number uint64
	if isInterval, proposerInterval := IsProposerUpdateInterval(blockNum); isInterval {
		number = blockNum - proposerInterval
	} else {
		number = blockNum - (blockNum % proposerInterval)

	}
	return number
}

func SetStakingUpdateInterval(num uint64) {
	atomic.StoreUint64(&stakingUpdateInterval, num)
}

func StakingUpdateInterval() uint64 {
	ret := atomic.LoadUint64(&stakingUpdateInterval)
	return ret
}

func SetProposerUpdateInterval(num uint64) {
	atomic.StoreUint64(&proposerUpdateInterval, num)
}

func ProposerUpdateInterval() uint64 {
	ret := atomic.LoadUint64(&proposerUpdateInterval)
	return ret
}
