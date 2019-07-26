package reward

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// checking calculate blockNumber of stakingInfo and return the stakingInfo with the blockNumber correct when stakingInfo is stored in cache
func TestStakingManager_getStakingInfoFromStakingCache(t *testing.T) {
	testData := []uint64{
		0, 86400, 172800, 259200,
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
	stakingManager := newStakingManager(newTestBlockChain(), newDefaultTestGovernance())

	for i := 0; i < len(testData); i++ {
		testStakingInfo := newEmptyStakingInfo(testData[i])
		stakingManager.sic.add(testStakingInfo)
	}

	// should find a correct stakingInfo with a given block number
	for i := 0; i < len(testCases); i++ {
		resultStakingInfo := stakingManager.getStakingInfo(testCases[i].stakingNumber)
		assert.Equal(t, testCases[i].expectedNumber, resultStakingInfo.BlockNum)
	}
}
