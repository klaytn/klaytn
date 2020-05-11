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

package reward

import (
	"encoding/binary"
	"encoding/json"
	"github.com/klaytn/klaytn/storage/database"
)

type stakingInfoDB struct {
	db database.Database
}

func (sdb *stakingInfoDB) get(blockNum uint64) *StakingInfo {
	if sdb.db == nil {
		logger.Debug("stakingInfoDB.get() is called but stakingInfoDB is not set.")
		return nil
	}

	key := intToByte(blockNum)
	value, err := sdb.db.Get(key)
	if err != nil {
		logger.Error("Failed to get staking info from stakingInfoDB.", "blockNum", blockNum, "err", err)
		return nil
	}

	stakingInfo := new(StakingInfo)
	err = json.Unmarshal(value, stakingInfo)
	if err != nil {
		logger.Error("Failed to unmarshal staking info.", "blockNum", blockNum, "err", err, "value", value)
		return nil
	}

	return stakingInfo
}

func (sdb *stakingInfoDB) add(stakingInfo *StakingInfo) {
	if sdb.db == nil {
		logger.Debug("stakingInfoDB.add() is called but stakingInfoDB is not set.")
		return
	}
	key := intToByte(stakingInfo.BlockNum)

	value, err := json.Marshal(stakingInfo)
	if err != nil {
		logger.Error("Failed to marshal staking info before adding.", "err", err, "stakingInfo", stakingInfo)
	}

	err = sdb.db.Put(key, value)
	if err != nil {
		logger.Error("Failed to put staking info to DB.", "blockNum", stakingInfo.BlockNum, "err", err, "stakingInfo", stakingInfo)
	}
	logger.Info("Add a new stakingInfo to stakingInfoDB", "stakingInfo", stakingInfo)
}

// intToByte converts value of int to []byte
// This is used for making keys in stakingInfoDB
func intToByte(num uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, num)
	return bs
}
