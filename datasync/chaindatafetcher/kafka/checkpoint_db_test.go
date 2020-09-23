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

package kafka

import (
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestCheckpointDB_ReadCheckpoint(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	testDB := database.NewMemoryDBManager()
	db := NewCheckpointDB()
	db.SetComponent(testDB)

	// 1. nothing is stored
	checkpoint, err := db.ReadCheckpoint()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), checkpoint)

	// 2. store a checkpoint
	expected := rand.Int63()
	assert.NoError(t, db.WriteCheckpoint(expected))

	// 3. assert the expected
	actual, err := db.ReadCheckpoint()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// 4. overwrite checkpoint
	expected2 := rand.Int63()
	assert.NoError(t, db.WriteCheckpoint(expected2))

	// 5. assert the expected 2
	actual2, err := db.ReadCheckpoint()
	assert.NoError(t, err)
	assert.Equal(t, expected2, actual2)
}
