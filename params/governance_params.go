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
	"math/big"
	"sync/atomic"
)

var (
	stakingUpdateInterval  uint64 = DefaultStakeUpdateInterval
	proposerUpdateInterval uint64 = DefaultProposerRefreshInterval
	lowerBoundBaseFee      uint64 = DefaultLowerBoundBaseFee
	upperBoundBaseFee      uint64 = DefaultUpperBoundBaseFee
	gasTarget              uint64 = DefaultGasTarget
	blockGasLimit          uint64 = DefaultBlockGasLimit
	baseFeeDenominator     uint64 = DefaultBaseFeeDenominator
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
	LowerBoundBaseFee
	UpperBoundBaseFee
	GasTarget
	BlockGasLimit
	BaseFeeDenominator
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

var (
	// Default Values: Constants used for getting default values for configuration
	DefaultGovernanceMode          = "none"
	DefaultGoverningNode           = "0x0000000000000000000000000000000000000000"
	DefaultEpoch                   = uint64(604800)
	DefaultProposerPolicy          = uint64(RoundRobin)
	DefaultSubGroupSize            = uint64(21)
	DefaultUnitPrice               = uint64(250000000000)
	DefaultLowerBoundBaseFee       = uint64(25000000000)
	DefaultUpperBoundBaseFee       = uint64(750000000000)
	DefaultGasTarget               = uint64(30000000)
	DefaultBlockGasLimit           = uint64(84000000)
	DefaultBaseFeeDenominator      = uint64(64)
	DefaultMintingAmount           = big.NewInt(0)
	DefaultRatio                   = "100/0/0"
	DefaultUseGiniCoeff            = false
	DefaultDefferedTxFee           = false
	DefaultMinimumStake            = big.NewInt(2000000)
	DefaultStakeUpdateInterval     = uint64(86400) // 1 day
	DefaultProposerRefreshInterval = uint64(3600)  // 1 hour
	DefaultPeriod                  = uint64(1)
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

func SetLowerBoundBaseFee(num uint64) {
	atomic.StoreUint64(&lowerBoundBaseFee, num)

}
func GetLowerBoundBaseFee() uint64 {
	ret := atomic.LoadUint64(&lowerBoundBaseFee)
	return ret
}

func SetUpperBoundBaseFee(num uint64) {
	atomic.StoreUint64(&upperBoundBaseFee, num)
}

func GetUpperBoundBaseFee() uint64 {
	ret := atomic.LoadUint64(&upperBoundBaseFee)
	return ret
}

func SetGasTarget(num uint64) {
	atomic.StoreUint64(&gasTarget, num)
}

func GetGasTarget() uint64 {
	ret := atomic.LoadUint64(&gasTarget)
	return ret
}

func SetBlockGasLimit(num uint64) {
	atomic.StoreUint64(&blockGasLimit, num)
}

func GetBlockGasLimit() uint64 {
	ret := atomic.LoadUint64(&blockGasLimit)
	return ret
}

func SetBaseFeeDenominator(num uint64) {
	atomic.StoreUint64(&baseFeeDenominator, num)
}

func GetBaseFeeDenominator() uint64 {
	ret := atomic.LoadUint64(&baseFeeDenominator)
	return ret
}
