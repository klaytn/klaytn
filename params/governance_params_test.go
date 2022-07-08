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
	"testing"
)

func TestSetProposerUpdateInterval(t *testing.T) {
	testData := []uint64{3600, 100, 500, 7200, 10}

	for i := 0; i < len(testData); i++ {
		SetProposerUpdateInterval(testData[i])

		if ProposerUpdateInterval() != testData[i] {
			t.Errorf("ProposerUadateInterval is different from the given testData. Result : %v, Expected : %v", ProposerUpdateInterval(), testData[i])
		}
	}
}

func TestSetStakingUpdateInterval(t *testing.T) {
	testData := []uint64{1000, 3600, 86400, 100, 500, 7200, 10}

	for i := 0; i < len(testData); i++ {
		SetStakingUpdateInterval(testData[i])

		if StakingUpdateInterval() != testData[i] {
			t.Errorf("StakingUpdateInterval is different from the given testData. Result : %v, Expected : %v", StakingUpdateInterval(), testData[i])
		}
	}
}

func TestIsProposerUpdateInterval(t *testing.T) {
	testCase := []struct {
		interval uint64
		blockNu  uint64
		result   bool
	}{
		{10, 3, false},
		{10, 10, true},
		{10, 11, false},
		{10, 20, true},
		{10, 21, false},
		{100, 99, false},
		{100, 100, true},
		{100, 101, false},
		{100, 200, true},
		{100, 201, false},
		{3600, 3599, false},
		{3600, 3600, true},
		{3600, 3601, false},
		{3600, 7200, true},
		{3600, 7201, false},
		{3600, 36000, true},
		{3600, 36001, false},
	}

	for i := 0; i < len(testCase); i++ {
		SetProposerUpdateInterval(testCase[i].interval)

		result, _ := IsProposerUpdateInterval(testCase[i].blockNu)

		if result != testCase[i].result {
			t.Errorf("The result is different from the expected result. Result : %v, Expected : %v, block number : %v, update interval : %v",
				result, testCase[i].result, testCase[i].blockNu, testCase[i].interval)
		}
	}
}

func TestIsStakingUpdatePossible(t *testing.T) {
	testCase := []struct {
		interval uint64
		blockNu  uint64
		result   bool
	}{
		{10, 3, false},
		{10, 10, true},
		{10, 11, false},
		{10, 20, true},
		{10, 21, false},
		{3600, 3599, false},
		{3600, 3600, true},
		{3600, 3601, false},
		{3600, 7200, true},
		{3600, 7201, false},
		{3600, 36000, true},
		{3600, 36001, false},
		{86400, 86399, false},
		{86400, 86400, true},
		{86400, 86401, false},
		{86400, 864000, true},
		{86400, 864001, false},
	}

	for i := 0; i < len(testCase); i++ {
		SetStakingUpdateInterval(testCase[i].interval)

		result := IsStakingUpdateInterval(testCase[i].blockNu)

		if result != testCase[i].result {
			t.Errorf("The result is different from the expected result. Result : %v, Expected : %v, block number : %v, update interval : %v",
				result, testCase[i].result, testCase[i].blockNu, testCase[i].interval)
		}
	}
}

func TestCalcProposerBlockNumber(t *testing.T) {
	testCase := []struct {
		interval uint64
		blockNu  uint64
		result   uint64
	}{
		{10, 3, 0},
		{10, 10, 0},
		{10, 11, 10},
		{10, 20, 10},
		{10, 21, 20},
		{3600, 3599, 0},
		{3600, 3600, 0},
		{3600, 3601, 3600},
		{3600, 7199, 3600},
		{3600, 7200, 3600},
		{3600, 7201, 7200},
	}

	for i := 0; i < len(testCase); i++ {
		SetProposerUpdateInterval(testCase[i].interval)

		result := CalcProposerBlockNumber(testCase[i].blockNu)

		if result != testCase[i].result {
			t.Errorf("The result is different from the expected result. Result : %v, Expected : %v, block number : %v, update interval : %v",
				result, testCase[i].result, testCase[i].blockNu, testCase[i].interval)
		}
	}
}

func TestCalcStakingBlockNumber(t *testing.T) {
	testCase := []struct {
		interval uint64
		blockNu  uint64
		result   uint64
	}{
		{10, 3, 0},
		{10, 10, 0},
		{10, 11, 0},
		{10, 20, 0},
		{10, 21, 10},
		{10, 30, 10},
		{10, 31, 20},
		{10, 40, 20},
		{3600, 3600, 0},
		{3600, 5000, 0},
		{3600, 7200, 0},
		{3600, 7201, 3600},
		{3600, 10800, 3600},
		{3600, 10801, 7200},
		{86400, 3600, 0},
		{86400, 10000, 0},
		{86400, 86400, 0},
		{86400, 172800, 0},
		{86400, 172801, 86400},
		{86400, 259200, 86400},
		{86400, 259201, 172800},
	}

	for i := 0; i < len(testCase); i++ {
		SetStakingUpdateInterval(testCase[i].interval)

		result := CalcStakingBlockNumber(testCase[i].blockNu)

		if result != testCase[i].result {
			t.Errorf("The result is different from the expected result. Result : %v, Expected : %v, block number : %v, update interval : %v",
				result, testCase[i].result, testCase[i].blockNu, testCase[i].interval)
		}
	}
}
