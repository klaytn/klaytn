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
	"encoding/json"
	"errors"
)

var ErrStakingDBNotSet = errors.New("stakingInfoDB is not set")

type stakingInfoDB interface {
	ReadStakingInfo(blockNum uint64) ([]byte, error)
	WriteStakingInfo(blockNum uint64, stakingInfo []byte) error
}

func getStakingInfoFromDB(blockNum uint64) (*StakingInfo, error) {
	if stakingManager.stakingInfoDB == nil {
		return nil, ErrStakingDBNotSet
	}

	jsonByte, err := stakingManager.stakingInfoDB.ReadStakingInfo(blockNum)
	if err != nil {
		return nil, err
	}

	stakingInfo := new(StakingInfo)
	err = json.Unmarshal(jsonByte, stakingInfo)
	if err != nil {
		return nil, err
	}

	return stakingInfo, nil
}

func addStakingInfoToDB(stakingInfo *StakingInfo) error {
	if stakingManager.stakingInfoDB == nil {
		return ErrStakingDBNotSet
	}

	marshaledStakingInfo, err := json.Marshal(stakingInfo)
	if err != nil {
		return err
	}

	err = stakingManager.stakingInfoDB.WriteStakingInfo(stakingInfo.BlockNum, marshaledStakingInfo)
	if err != nil {
		return err
	}

	return nil
}
