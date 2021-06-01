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
	"github.com/klaytn/klaytn/storage/database"
)

type CheckpointDB struct {
	manager database.DBManager
}

func NewCheckpointDB() *CheckpointDB {
	return &CheckpointDB{}
}

func (db *CheckpointDB) ReadCheckpoint() (int64, error) {
	checkpoint, err := db.manager.ReadChainDataFetcherCheckpoint()
	return int64(checkpoint), err
}

func (db *CheckpointDB) WriteCheckpoint(checkpoint int64) error {
	return db.manager.WriteChainDataFetcherCheckpoint(uint64(checkpoint))
}

func (db *CheckpointDB) SetComponent(component interface{}) {
	switch c := component.(type) {
	case database.DBManager:
		db.manager = c
	}
}
