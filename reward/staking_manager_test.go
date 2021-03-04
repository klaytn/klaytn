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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	stakingInterval = uint64(86400)
	testData        = []uint64{
		0, 1, 2, 3,
	}
	testCases = []struct {
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
)

func TestStakingManager_NewStakingManager(t *testing.T) {
	// test if nil
	assert.Nil(t, GetStakingManager())
	assert.Nil(t, GetStakingInfo(123))

	st, err := updateStakingInfo(456)
	assert.Nil(t, st)
	assert.EqualError(t, err, ErrStakingManagerNotSet.Error())

	assert.EqualError(t, CheckStakingInfoStored(789), ErrStakingManagerNotSet.Error())

	// test if get same
	stNew := NewStakingManager(newTestBlockChain(), newDefaultTestGovernance(), nil)
	stGet := GetStakingManager()
	assert.NotNil(t, stNew)
	assert.Equal(t, stGet, stNew)
}

// checking calculate blockNumber of stakingInfo and return the stakingInfo with the blockNumber correct when stakingInfo is stored in cache
func TestStakingManager_getStakingInfoFromStakingCache(t *testing.T) {
	stakingManager := NewStakingManager(newTestBlockChain(), newDefaultTestGovernance(), nil)

	for i := 0; i < len(testData); i++ {
		testStakingInfo := newEmptyStakingInfo(testData[i] * stakingInterval)
		stakingManager.stakingInfoCache.add(testStakingInfo)
	}

	// should find a correct stakingInfo with a given block number
	for i := 0; i < len(testCases); i++ {
		resultStakingInfo := GetStakingInfo(testCases[i].stakingNumber)
		assert.Equal(t, testCases[i].expectedNumber, resultStakingInfo.BlockNum)
	}
}
