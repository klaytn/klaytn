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

package chaindatafetcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasDataType(t *testing.T) {
	testdata := []requestType{
		// single type
		requestTypeTransaction,
		requestTypeTokenTransfer,
		requestTypeTraces,
		requestTypeContracts,

		// 2 composite type
		requestTypeTransaction | requestTypeTokenTransfer,
		requestTypeTransaction | requestTypeTraces,
		requestTypeTransaction | requestTypeContracts,
		requestTypeTokenTransfer | requestTypeTraces,
		requestTypeTokenTransfer | requestTypeContracts,
		requestTypeTraces | requestTypeContracts,

		// 3 composite type
		requestTypeTransaction | requestTypeTokenTransfer | requestTypeTraces,
		requestTypeTransaction | requestTypeTokenTransfer | requestTypeContracts,
		requestTypeTokenTransfer | requestTypeTraces | requestTypeContracts,

		// all type
		requestTypeAll,
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
		{false, true, true, true},

		// all type
		{true, true, true, true},
	}

	checkFuncs := []func(requestType) bool{
		hasTransactions, hasTokenTransfers, hasTraces, hasContracts,
	}

	for idx, types := range testdata {
		expectedTypes := expected[idx]
		for idx, check := range checkFuncs {
			assert.Equal(t, expectedTypes[idx], check(types))
		}
	}
}
