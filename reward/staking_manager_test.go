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
	"testing"

	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

type stakingManagerTestCase struct {
	blockNum    uint64       // Requested num in GetStakingInfo(num)
	stakingNum  uint64       // Corresponding staking info block number
	stakingInfo *StakingInfo // Expected GetStakingInfo() output
}

// Note that Golang will correctly initialize these globals according to dependency.
// https://go.dev/ref/spec#Order_of_evaluation

// Note the testdata must not exceed maxStakingCache because otherwise cache test will fail.
var stakingManagerTestData = []*StakingInfo{
	stakingInfoTestCases[0].stakingInfo,
	stakingInfoTestCases[1].stakingInfo,
	stakingInfoTestCases[2].stakingInfo,
	stakingInfoTestCases[3].stakingInfo,
}
var stakingManagerTestCases = generateStakingManagerTestCases()

func generateStakingManagerTestCases() []stakingManagerTestCase {
	s := stakingManagerTestData

	return []stakingManagerTestCase{
		{1, 0, s[0]},
		{100, 0, s[0]},
		{86400, 0, s[0]},
		{86401, 0, s[0]},
		{172800, 0, s[0]},
		{172801, 86400, s[1]},
		{200000, 86400, s[1]},
		{259200, 86400, s[1]},
		{259201, 172800, s[2]},
		{300000, 172800, s[2]},
		{345600, 172800, s[2]},
		{345601, 259200, s[3]},
		{400000, 259200, s[3]},
	}
}

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

// Check that appropriate StakingInfo is returned given various blockNum argument.
func checkGetStakingInfo(t *testing.T) {
	for _, testcase := range stakingManagerTestCases {
		expcectedInfo := testcase.stakingInfo
		actualInfo := GetStakingInfo(testcase.blockNum)

		assert.Equal(t, testcase.stakingNum, actualInfo.BlockNum)
		assert.Equal(t, expcectedInfo, actualInfo)
	}
}

// Check that StakinInfo are loaded from cache
func TestStakingManager_GetFromCache(t *testing.T) {
	stakingManager := NewStakingManager(newTestBlockChain(), newDefaultTestGovernance(), nil)

	for _, testdata := range stakingManagerTestData {
		stakingManager.stakingInfoCache.add(testdata)
	}

	checkGetStakingInfo(t)
}

// Check that StakinInfo are loaded from database
func TestStakingManager_GetFromDB(t *testing.T) {
	db := database.NewMemoryDBManager()
	NewStakingManager(newTestBlockChain(), newDefaultTestGovernance(), db)

	for _, testdata := range stakingManagerTestData {
		AddStakingInfoToDB(testdata)
	}

	checkGetStakingInfo(t)
}

// Even if Gini was -1 in the cache, GetStakingInfo returns valid Gini
func TestStakingManager_FillGiniFromCache(t *testing.T) {
	stakingManager := NewStakingManager(newTestBlockChain(), newDefaultTestGovernance(), nil)

	for _, testdata := range stakingManagerTestData {
		// Insert a modified copy of testdata to cache
		copydata := &StakingInfo{}
		json.Unmarshal([]byte(testdata.String()), copydata)
		copydata.Gini = -1 // Suppose Gini was -1 in the cache
		stakingManager.stakingInfoCache.add(copydata)
	}

	checkGetStakingInfo(t)
}

// Even if Gini was -1 in the DB, GetStakingInfo returns valid Gini
func TestStakingManager_FillGiniFromDB(t *testing.T) {
	db := database.NewMemoryDBManager()
	NewStakingManager(newTestBlockChain(), newDefaultTestGovernance(), db)

	for _, testdata := range stakingManagerTestData {
		// Insert a modified copy of testdata to cache
		copydata := &StakingInfo{}
		json.Unmarshal([]byte(testdata.String()), copydata)
		copydata.Gini = -1 // Suppose Gini was -1 in the cache
		AddStakingInfoToDB(testdata)
	}

	checkGetStakingInfo(t)
}
