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
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChainDataFetcher_updateCheckpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	fetcher := &ChainDataFetcher{
		checkpoint:    0,
		checkpointMap: make(map[int64]struct{}),
		repo:          m,
	}

	// update checkpoint as follows.
	// done order: 1, 0, 2, 3, 5, 7, 9, 8, 4, 6, 10
	// checkpoint: 0, 2, 3, 4, 4, 4, 4, 4, 6, 10, 11
	assert.NoError(t, fetcher.updateCheckpoint(1))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(2)))
	assert.NoError(t, fetcher.updateCheckpoint(0))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(3)))
	assert.NoError(t, fetcher.updateCheckpoint(2))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(4)))
	assert.NoError(t, fetcher.updateCheckpoint(3))
	assert.NoError(t, fetcher.updateCheckpoint(5))
	assert.NoError(t, fetcher.updateCheckpoint(7))
	assert.NoError(t, fetcher.updateCheckpoint(9))
	assert.NoError(t, fetcher.updateCheckpoint(8))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(6)))
	assert.NoError(t, fetcher.updateCheckpoint(4))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(10)))
	assert.NoError(t, fetcher.updateCheckpoint(6))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(11)))
	assert.NoError(t, fetcher.updateCheckpoint(10))
}
