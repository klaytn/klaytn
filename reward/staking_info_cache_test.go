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

// test cache limit of stakingInfoCache
func TestStakingInfoCache_Add_Limit(t *testing.T) {
	stakingInfoCache := newStakingInfoCache()

	for i := 1; i <= 10; i++ {
		testStakingInfo := newEmptyStakingInfo(uint64(i))
		stakingInfoCache.add(testStakingInfo)

		if len(stakingInfoCache.cells) > maxStakingCache {
			t.Errorf("over the max limit of stakingCache. Current Len : %v, MaxStakingCache : %v", len(stakingInfoCache.cells), maxStakingCache)
		}
	}
}

func TestStakingInfoCache_Add_SameNumber(t *testing.T) {
	stakingInfoCache := newStakingInfoCache()

	testStakingInfo1 := newEmptyStakingInfo(uint64(1))
	testStakingInfo2 := newEmptyStakingInfo(uint64(1))

	stakingInfoCache.add(testStakingInfo1)
	stakingInfoCache.add(testStakingInfo2)

	assert.Equal(t, 1, len(stakingInfoCache.cells), "StakingInfo with Same block number is saved to the stakingCache")
}

func TestStakingInfoCache_Add_SmallNumber(t *testing.T) {
	stakingInfoCache := newStakingInfoCache()

	for i := uint64(10); i > 0; i-- {
		testStakingInfo := newEmptyStakingInfo(i)
		stakingInfoCache.add(testStakingInfo)
		assert.Equal(t, i, stakingInfoCache.minBlockNum)
	}
}

// stakingInfo with minBlockNum should be deleted if add more than limit
func TestStakingInfoCache_Add_MinBlockNum(t *testing.T) {
	stakingInfoCache := newStakingInfoCache()

	for i := 1; i < 5; i++ {
		testStakingInfo := newEmptyStakingInfo(uint64(i))
		stakingInfoCache.add(testStakingInfo)
	}

	testStakingInfo := newEmptyStakingInfo(uint64(5))
	stakingInfoCache.add(testStakingInfo) // blockNum 1 should be deleted
	assert.Equal(t, uint64(2), stakingInfoCache.minBlockNum)

	testStakingInfo = newEmptyStakingInfo(uint64(6))
	stakingInfoCache.add(testStakingInfo) // blockNum 2 should be deleted
	assert.Equal(t, uint64(3), stakingInfoCache.minBlockNum)
}

func TestStakingInfoCache_Add(t *testing.T) {
	testCases := []struct {
		blockNumber       uint64
		expectedLen       int
		expectedMinNumber uint64
	}{
		{1, 1, 1},
		{5, 2, 1},
		{10, 3, 1},
		{7, 4, 1},
		{15, 4, 5},
		{30, 4, 7},
		{20, 4, 10},
		{3, 4, 3},
	}
	stakingInfoCache := newStakingInfoCache()
	for i := 0; i < len(testCases); i++ {
		testStakingInfo := newEmptyStakingInfo(testCases[i].blockNumber)
		stakingInfoCache.add(testStakingInfo)
		assert.Equal(t, testCases[i].expectedLen, len(stakingInfoCache.cells))
		assert.Equal(t, testCases[i].expectedMinNumber, stakingInfoCache.minBlockNum)
	}
}

func TestStakingInfoCache_Get(t *testing.T) {
	stakingInfoCache := newStakingInfoCache()

	for i := 1; i <= 4; i++ {
		testStakingInfo := newEmptyStakingInfo(uint64(i))
		stakingInfoCache.add(testStakingInfo)
	}

	// should find correct stakingInfo with a given block number
	for i := uint64(1); i <= 4; i++ {
		testStakingInfo := stakingInfoCache.get(i)
		assert.Equal(t, i, testStakingInfo.BlockNum)
	}

	// nothing should be found as no matched block number is in the cache
	for i := uint64(5); i < 10; i++ {
		testStakingInfo := stakingInfoCache.get(i)
		assert.Nil(t, testStakingInfo)
	}
}
