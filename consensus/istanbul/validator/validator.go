// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/validator/validator.go (2018/06/04).
// Modified and improved for the klaytn development.

package validator

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/log"
)

var logger = log.NewModuleLogger(log.ConsensusIstanbulValidator)

func New(addr common.Address) istanbul.Validator {
	return &defaultValidator{
		address: addr,
	}
}

func NewValidatorSet(addrs []common.Address, proposerPolicy istanbul.ProposerPolicy, subGroupSize uint64, chain consensus.ChainReader) istanbul.ValidatorSet {
	var valSet istanbul.ValidatorSet
	if proposerPolicy == istanbul.WeightedRandom {
		valSet = NewWeightedCouncil(addrs, nil, nil, nil, proposerPolicy, subGroupSize, 0, 0, chain)
	} else {
		valSet = NewSubSet(addrs, proposerPolicy, subGroupSize)
	}

	return valSet
}

func NewSet(addrs []common.Address, policy istanbul.ProposerPolicy) istanbul.ValidatorSet {
	return newDefaultSet(addrs, policy)
}

func NewSubSet(addrs []common.Address, policy istanbul.ProposerPolicy, subSize uint64) istanbul.ValidatorSet {
	return newDefaultSubSet(addrs, policy, subSize)
}

func ExtractValidators(extraData []byte) []common.Address {
	// get the validator addresses
	addrs := make([]common.Address, (len(extraData) / common.AddressLength))
	for i := 0; i < len(addrs); i++ {
		copy(addrs[i][:], extraData[i*common.AddressLength:])
	}

	return addrs
}
