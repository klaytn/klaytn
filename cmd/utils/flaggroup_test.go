// Copyright 2020 The klaytn Authors
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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/urfave/cli.v1"
)

func TestFlagGroups_Duplication(t *testing.T) {
	testGroups := FlagGroups
	exist := make(map[string]bool)

	for _, group := range testGroups {
		for _, flag := range group.Flags {
			if exist[flag.GetName()] == true {
				t.Error("a flag belong to more than one group", "flag", flag.GetName())
			}
			exist[flag.GetName()] = true
		}
	}
}

func TestCategorizeFlags(t *testing.T) {
	var testFlags []cli.Flag

	// Sort flagGroups to compare with actualGroups which is a sorted slice.
	expectedGroups := sortFlagGroup(FlagGroups, uncategorized)

	// Prepare flags to be categorized
	for _, group := range expectedGroups {
		testFlags = append(testFlags, group.Flags...)
	}

	actualGroups := CategorizeFlags(testFlags)

	assert.Equal(t, expectedGroups, actualGroups)
}

func TestCategorizeFlags_Duplication(t *testing.T) {
	var testFlags []cli.Flag
	var testFlagsDup []cli.Flag

	// Prepare a slice of test flags and a slice of duplicated test flags.
	for _, group := range FlagGroups {
		testFlags = append(testFlags, group.Flags...)
		testFlagsDup = append(testFlagsDup, group.Flags...)
		testFlagsDup = append(testFlagsDup, group.Flags...)
	}

	expected := CategorizeFlags(testFlags)
	actual := CategorizeFlags(testFlagsDup)

	assert.Equal(t, expected, actual)
}
