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

package database

import (
	"encoding/json"
	"github.com/klaytn/klaytn/common"
)

// ReadStakingInfo reads staking information from database. It returns
// (StakingInfo, nil) if it succeeds to read and (nil, error) if it fails.
// StakingInfo is stored in MiscDB.
// Be sure to use the right block number before calling this function.
// (Refer to CalcStakingBlockNumber() in params/governance_params.go)
func (dbm *databaseManager) ReadStakingInfo(blockNum uint64) (interface{}, error) {
	db := dbm.getDatabase(MiscDB)

	key := makeStakingInfoKey(blockNum)
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	stakingInfo := new(interface{})
	err = json.Unmarshal(value, stakingInfo)
	if err != nil {
		return nil, err
	}

	return stakingInfo, nil
}

// WriteStakingInfo writes staking information to database. It returns
// nil if it succeeds to write and error if it fails.
// Key should be the blockNum of stakingIfo. Value is marshaled
// stakingInfo. StakingInfo is stored in MiscDB.
// stakingInfo should be type StakingInfo defined in reward/staking_info.go
// Be sure to use the right block number before calling this function.
// (Refer to CalcStakingBlockNumber() in params/governance_params.go)
func (dbm *databaseManager) WriteStakingInfo(blockNum uint64, stakingInfo interface{}) error {
	db := dbm.getDatabase(MiscDB)

	key := makeStakingInfoKey(blockNum)
	value, err := json.Marshal(stakingInfo)
	if err != nil {
		return err
	}

	return db.Put(key, value)
}

// makeStakingInfoKey is used for making keys for staking info
func makeStakingInfoKey(num uint64) []byte {
	key := append(stakingInfoPrefix, common.Int64ToByteLittleEndian(num)...)
	return key
}
