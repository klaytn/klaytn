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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasDataType(t *testing.T) {
	testdata := []requestType{
		// single type
		requestTypeTransaction,
		requestTypeTokenTransfer,
		requestTypeTrace,
		requestTypeContract,

		// 2 composite type
		requestTypeTransaction | requestTypeTokenTransfer,
		requestTypeTransaction | requestTypeTrace,
		requestTypeTransaction | requestTypeContract,
		requestTypeTokenTransfer | requestTypeTrace,
		requestTypeTokenTransfer | requestTypeContract,
		requestTypeTrace | requestTypeContract,

		// 3 composite type
		requestTypeTransaction | requestTypeTokenTransfer | requestTypeTrace,
		requestTypeTransaction | requestTypeTokenTransfer | requestTypeContract,
		requestTypeTransaction | requestTypeTrace | requestTypeContract,
		requestTypeTokenTransfer | requestTypeTrace | requestTypeContract,

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
		{true, false, true, true},
		{false, true, true, true},

		// all type
		{true, true, true, true},
	}

	targetType := []requestType{
		requestTypeTransaction, requestTypeTokenTransfer, requestTypeTrace, requestTypeContract,
	}

	for idx, types := range testdata {
		expectedTypes := expected[idx]
		for idx, target := range targetType {
			assert.Equal(t, expectedTypes[idx], checkRequestType(types, target))
		}
	}
}
