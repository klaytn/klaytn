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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
)

const (
	AddrNotFoundInCouncilNodes = -1
	maxStakingLimit            = uint64(100000000000)
	DefaultGiniCoefficient     = -1.0
)

var (
	maxStakingLimitBigInt = big.NewInt(0).SetUint64(maxStakingLimit)

	ErrAddrNotInStakingInfo = errors.New("Address is not in stakingInfo")
)

// StakingInfo contains staking information.
type StakingInfo struct {
	BlockNum uint64 // Block number where staking information of Council is fetched

	// Information retrieved from AddressBook smart contract
	CouncilNodeAddrs    []common.Address // NodeIds of Council
	CouncilStakingAddrs []common.Address // Address of Staking account which holds staking balance
	CouncilRewardAddrs  []common.Address // Address of Council account which will get block reward
	KIRAddr             common.Address   // Address of KIR contract
	PoCAddr             common.Address   // Address of PoC contract

	UseGini bool
	Gini    float64 // gini coefficient

	// Derived from CouncilStakingAddrs
	CouncilStakingAmounts []uint64 // Staking amounts of Council
}

func newEmptyStakingInfo(blockNum uint64) *StakingInfo {
	stakingInfo := &StakingInfo{
		BlockNum:              blockNum,
		CouncilNodeAddrs:      make([]common.Address, 0, 0),
		CouncilStakingAddrs:   make([]common.Address, 0, 0),
		CouncilRewardAddrs:    make([]common.Address, 0, 0),
		KIRAddr:               common.Address{},
		PoCAddr:               common.Address{},
		CouncilStakingAmounts: make([]uint64, 0, 0),
		Gini:                  DefaultGiniCoefficient,
		UseGini:               false,
	}
	return stakingInfo
}

func newStakingInfo(bc blockChain, helper governanceHelper, blockNum uint64, nodeAddrs []common.Address, stakingAddrs []common.Address, rewardAddrs []common.Address, KIRAddr common.Address, PoCAddr common.Address) (*StakingInfo, error) {
	intervalBlock := bc.GetBlockByNumber(blockNum)
	if intervalBlock == nil {
		logger.Trace("Failed to get the block by the given number", "blockNum", blockNum)
		return nil, errors.New(fmt.Sprintf("Failed to get the block by the given number. blockNum: %d", blockNum))
	}
	statedb, err := bc.StateAt(intervalBlock.Root())
	if err != nil {
		logger.Trace("Failed to make a state for interval block", "interval blockNum", blockNum, "err", err)
		return nil, err
	}

	// Get balance of stakingAddrs
	stakingAmounts := make([]uint64, len(stakingAddrs))
	for i, stakingAddr := range stakingAddrs {
		tempStakingAmount := big.NewInt(0).Div(statedb.GetBalance(stakingAddr), big.NewInt(0).SetUint64(params.KLAY))
		if tempStakingAmount.Cmp(maxStakingLimitBigInt) > 0 {
			tempStakingAmount.SetUint64(maxStakingLimit)
		}
		stakingAmounts[i] = tempStakingAmount.Uint64()
	}

	var useGini bool
	if res, err := helper.GetItemAtNumberByIntKey(blockNum, params.UseGiniCoeff); err != nil {
		logger.Trace("Failed to get useGiniCoeff from governance", "blockNum", blockNum, "err", err)
		return nil, err
	} else {
		useGini = res.(bool)
	}
	gini := DefaultGiniCoefficient

	stakingInfo := &StakingInfo{
		BlockNum:              blockNum,
		CouncilNodeAddrs:      nodeAddrs,
		CouncilStakingAddrs:   stakingAddrs,
		CouncilRewardAddrs:    rewardAddrs,
		KIRAddr:               KIRAddr,
		PoCAddr:               PoCAddr,
		CouncilStakingAmounts: stakingAmounts,
		Gini:                  gini,
		UseGini:               useGini,
	}
	return stakingInfo, nil
}

func (s *StakingInfo) GetIndexByNodeAddress(nodeAddress common.Address) (int, error) {
	for i, addr := range s.CouncilNodeAddrs {
		if addr == nodeAddress {
			return i, nil
		}
	}
	return AddrNotFoundInCouncilNodes, ErrAddrNotInStakingInfo
}

func (s *StakingInfo) GetStakingAmountByNodeId(nodeAddress common.Address) (uint64, error) {
	i, err := s.GetIndexByNodeAddress(nodeAddress)
	if err != nil {
		return 0, err
	}
	return s.CouncilStakingAmounts[i], nil
}

func (s *StakingInfo) String() string {
	j, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

type float64Slice []float64

func (p float64Slice) Len() int           { return len(p) }
func (p float64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func CalcGiniCoefficient(stakingAmount float64Slice) float64 {
	sort.Sort(stakingAmount)

	// calculate gini coefficient
	sumOfAbsoluteDifferences := float64(0)
	subSum := float64(0)

	for i, x := range stakingAmount {
		temp := x*float64(i) - subSum
		sumOfAbsoluteDifferences = sumOfAbsoluteDifferences + temp
		subSum = subSum + x
	}

	result := sumOfAbsoluteDifferences / subSum / float64(len(stakingAmount))
	result = math.Round(result*100) / 100

	return result
}
