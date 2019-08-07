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
	"github.com/stretchr/testify/assert"
	"testing"
)

// checking calculate blockNumber of stakingInfo and return the stakingInfo with the blockNumber correct when stakingInfo is stored in cache
func TestStakingManager_getStakingInfoFromStakingCache(t *testing.T) {
	stakingInterval := uint64(86400)
	testData := []uint64{
		0, 1, 2, 3,
	}
	testCases := []struct {
		stakingNumber  uint64
		expectedNumber uint64
	}{
		{1, 0},
		{100, 0},
		{86400, 0},
		{86401, 0},
		{172800, 0},
		{172801, 86400},
		{200000, 86400},
		{259200, 86400},
		{259201, 172800},
		{300000, 172800},
		{345600, 172800},
		{345601, 259200},
		{400000, 259200},
	}
	stakingManager := NewStakingManager(newTestBlockChain(), newDefaultTestGovernance())

	for i := 0; i < len(testData); i++ {
		testStakingInfo := newEmptyStakingInfo(testData[i] * stakingInterval)
		stakingManager.sic.add(testStakingInfo)
	}

	// should find a correct stakingInfo with a given block number
	for i := 0; i < len(testCases); i++ {
		resultStakingInfo := stakingManager.GetStakingInfo(testCases[i].stakingNumber)
		assert.Equal(t, testCases[i].expectedNumber, resultStakingInfo.BlockNum)
	}
}
