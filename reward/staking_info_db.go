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

	"github.com/klaytn/klaytn/common"
)

var ErrStakingDBNotSet = errors.New("stakingInfoDB is not set")

type stakingInfoDB interface {
	HasStakingInfo(blockNum uint64) (bool, error)
	ReadStakingInfo(blockNum uint64) ([]byte, error)
	WriteStakingInfo(blockNum uint64, stakingInfo []byte) error

	ReadCanonicalHash(number uint64) common.Hash
}

// HasStakingInfoFromDB returns existence of staking information from miscdb with blockhash.
// Note that blockhash is also returned in order to check the validity of the block.
func HasStakingInfoFromDB(blockNumber uint64) (common.Hash, bool, error) {
	if stakingManager.stakingInfoDB == nil {
		return common.Hash{}, false, ErrStakingDBNotSet
	}

	hash := stakingManager.stakingInfoDB.ReadCanonicalHash(blockNumber)
	has, err := stakingManager.stakingInfoDB.HasStakingInfo(blockNumber)

	return hash, has, err
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

func AddStakingInfoToDB(stakingInfo *StakingInfo) error {
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
