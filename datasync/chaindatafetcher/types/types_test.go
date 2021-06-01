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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasDataType(t *testing.T) {
	testdata := []RequestType{
		// single type
		RequestTypeTransaction,
		RequestTypeTokenTransfer,
		RequestTypeTrace,
		RequestTypeContract,

		// 2 composite type
		RequestTypeTransaction | RequestTypeTokenTransfer,
		RequestTypeTransaction | RequestTypeTrace,
		RequestTypeTransaction | RequestTypeContract,
		RequestTypeTokenTransfer | RequestTypeTrace,
		RequestTypeTokenTransfer | RequestTypeContract,
		RequestTypeTrace | RequestTypeContract,

		// 3 composite type
		RequestTypeTransaction | RequestTypeTokenTransfer | RequestTypeTrace,
		RequestTypeTransaction | RequestTypeTokenTransfer | RequestTypeContract,
		RequestTypeTransaction | RequestTypeTrace | RequestTypeContract,
		RequestTypeTokenTransfer | RequestTypeTrace | RequestTypeContract,

		// all type
		RequestTypeAll,
	}

	expected := [][]bool{
		// single type
		{true, false, false, false},
		{false, true, false, false},
		{false, false, true, false},
		{false, false, false, true},

		// composite type
		{true, true, false, false},
		{true, false, true, false},
		{true, false, false, true},
		{false, true, true, false},
		{false, true, false, true},
		{false, false, true, true},

		// 3 composite type
		{true, true, true, false},
		{true, true, false, true},
		{true, false, true, true},
		{false, true, true, true},

		// all type
		{true, true, true, true},
	}

	targetType := []RequestType{
		RequestTypeTransaction, RequestTypeTokenTransfer, RequestTypeTrace, RequestTypeContract,
	}

	for idx, types := range testdata {
		expectedTypes := expected[idx]
		for idx, target := range targetType {
			assert.Equal(t, expectedTypes[idx], CheckRequestType(types, target))
		}
	}
}
