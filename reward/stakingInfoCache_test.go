package reward

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

// stakingInfo with minBlockNum should be deleted if add more than limit
func TestStakingInfoCache_Add(t *testing.T) {
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
