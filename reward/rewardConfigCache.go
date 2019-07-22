// Copyright 2019 The klaytn Authors
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
	"errors"
	"github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/params"
	"math/big"
	"strconv"
	"strings"
)

var (
	FailGettingConfigure = errors.New("fail to get configure from governance")
)

const (
	maxRewardConfigCache = 5 // Must be positive integer
)

type rewardConfig struct {
	blockNum uint64

	mintingAmount *big.Int
	cnRatio       *big.Int
	pocRatio      *big.Int
	kirRatio      *big.Int
	totalRatio    *big.Int
	unitPrice     *big.Int
}

// Cache for parsed reward parameters from governance
type rewardConfigCache struct {
	cache            *lru.ARCCache
	governanceHelper governanceHelper
}

func newRewardConfigCache(governanceHelper governanceHelper) *rewardConfigCache {
	cache, _ := lru.NewARC(maxRewardConfigCache)
	return &rewardConfigCache{
		cache:            cache,
		governanceHelper: governanceHelper,
	}
}

func (rewardConfigCache *rewardConfigCache) get(blockNumber uint64) (*rewardConfig, error) {
	var epoch uint64
	if result, err := rewardConfigCache.governanceHelper.GetItemAtNumberByIntKey(blockNumber, params.Epoch); err == nil {
		epoch = result.(uint64)
	} else {
		logger.Error("Couldn't get epoch from governance", "blockNumber", blockNumber, "err", err)
		return nil, FailGettingConfigure
	}

	if blockNumber%epoch == 0 {
		blockNumber -= epoch
	} else {
		blockNumber -= (blockNumber % epoch)
	}

	if rewardConfigCache.cache.Contains(blockNumber) {
		config, _ := rewardConfigCache.cache.Get(blockNumber)
		return config.(*rewardConfig), nil
	}

	newConfig, err := rewardConfigCache.newRewardConfig(blockNumber)
	if err != nil {
		return nil, err
	}

	rewardConfigCache.add(blockNumber, newConfig)
	return newConfig, nil
}

func (rewardConfigCache *rewardConfigCache) newRewardConfig(blockNumber uint64) (*rewardConfig, error) {
	mintingAmount := big.NewInt(0)
	if result, err := rewardConfigCache.governanceHelper.GetItemAtNumberByIntKey(blockNumber, params.MintingAmount); err == nil {
		mintingAmount.SetString(result.(string), 10)
	} else {
		logger.Error("Couldn't get MintingAmount from governance", "blockNumber", blockNumber, "err", err)
		return nil, FailGettingConfigure
	}

	cnRatio := big.NewInt(0)
	pocRatio := big.NewInt(0)
	kirRatio := big.NewInt(0)
	totalRatio := big.NewInt(0)
	if result, err := rewardConfigCache.governanceHelper.GetItemAtNumberByIntKey(blockNumber, params.Ratio); err == nil {
		cn, poc, kir, parsingError := rewardConfigCache.parseRewardRatio(result.(string))
		if parsingError != nil {
			return nil, parsingError
		}
		cnRatio.SetInt64(int64(cn))
		pocRatio.SetInt64(int64(poc))
		kirRatio.SetInt64(int64(kir))
		totalRatio.SetInt64(int64(cn + poc + kir))
	} else {
		logger.Error("Couldn't get Ratio from governance", "blockNumber", blockNumber, "err", err)
		return nil, FailGettingConfigure
	}

	unitPrice := big.NewInt(0)
	if result, err := rewardConfigCache.governanceHelper.GetItemAtNumberByIntKey(blockNumber, params.UnitPrice); err == nil {
		unitPrice.SetUint64(result.(uint64))
	} else {
		logger.Error("Couldn't get MintingAmount from governance", "blockNumber", blockNumber, "err", err)
		return nil, FailGettingConfigure
	}

	rewardConfig := &rewardConfig{
		blockNum:      blockNumber,
		mintingAmount: mintingAmount,
		cnRatio:       cnRatio,
		pocRatio:      pocRatio,
		kirRatio:      kirRatio,
		totalRatio:    totalRatio,
		unitPrice:     unitPrice,
	}
	return rewardConfig, nil
}

func (rewardConfigCache *rewardConfigCache) add(blockNumber uint64, config *rewardConfig) {
	rewardConfigCache.cache.Add(blockNumber, config)
}

func (rewardConfigCache *rewardConfigCache) parseRewardRatio(ratio string) (int, int, int, error) {
	s := strings.Split(ratio, "/")
	if len(s) != 3 {
		return 0, 0, 0, errors.New("invalid format")
	}
	cn, err1 := strconv.Atoi(s[0])
	poc, err2 := strconv.Atoi(s[1])
	kir, err3 := strconv.Atoi(s[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, errors.New("parsing error")
	}
	return cn, poc, kir, nil
}
