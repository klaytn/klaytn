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

package fork

import (
	"errors"
	"math/big"

	"github.com/klaytn/klaytn/params"
)

var (
	// hardForkBlockNumberConfig contains only hardFork block number
	hardForkBlockNumberConfig *params.ChainConfig
)

// Rules returns the hard fork information
// CAUTIOUS: Use it when chainConfig value is not reachable
func Rules(blockNumber *big.Int) (*params.Rules, error) {
	if hardForkBlockNumberConfig == nil {
		return nil, errors.New("hardForkBlockNumberConfig variable is not initialized")
	}
	rules := hardForkBlockNumberConfig.Rules(blockNumber)
	return &rules, nil
}

// SetHardForkBlockNumberConfig sets values in HardForkConfig if it is not nil.
// CAUTIOUS: Calling this function is can be dangerous, so avoid using it except tests
func SetHardForkBlockNumberConfig(h *params.ChainConfig) error {
	if h == nil {
		return errors.New("hardForkBlockNumberConfig cannot be initialized as nil")
	}
	hardForkBlockNumberConfig = h
	return nil
}
