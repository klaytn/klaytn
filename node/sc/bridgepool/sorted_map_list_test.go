// Modifications Copyright 2019 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/tx_list.go (2018/06/04).
// Modified and improved for the klaytn development.

package bridgepool

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type item struct {
	nonce uint64
}

func newItem(nonce uint64) *item {
	return &item{
		nonce: nonce,
	}
}

func (i *item) Nonce() uint64 {
	return i.nonce
}

func TestSortedMap_Put(t *testing.T) {
	sizeOfSortedMap := UnlimitedItemSortedMap
	sortedMap := NewItemSortedMap(sizeOfSortedMap)

	// test a sortedMap without a size limit.
	for i := uint64(0); i < 10000000; i++ {
		err := sortedMap.Put(newItem(i))
		assert.NoError(t, err)
		assert.Equal(t, int(i+1), sortedMap.Len())
	}

	// test a sortedMap with a size limit.
	sizeOfSortedMap = 10
	sortedMap = NewItemSortedMap(sizeOfSortedMap)

	for i := uint64(0); i < uint64(sizeOfSortedMap); i++ {
		err := sortedMap.Put(newItem(i))
		assert.NoError(t, err)
		assert.Equal(t, int(i+1), sortedMap.Len())
	}

	// test that an item cannot be put into a full sortedMap.
	err := sortedMap.Put(newItem(uint64(sizeOfSortedMap)))
	t.Logf("returned err: %v", err)
	assert.True(t, errors.Is(err, ErrSizeLimit))
	assert.Equal(t, sizeOfSortedMap, sortedMap.Len())

	// test that an existing nonce item can be put (overwrite) into a full sortedMap.
	err = sortedMap.Put(newItem(uint64(sizeOfSortedMap - 1)))
	assert.NoError(t, err)
	assert.Equal(t, sizeOfSortedMap, sortedMap.Len())
}
