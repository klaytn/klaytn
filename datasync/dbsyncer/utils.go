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

package dbsyncer

import (
	"strings"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/rlp"
)

const (
	// mysql driver (go) has query parameter max size (65535)
	// transaction record has 15 parameters
	BULK_INSERT_SIZE = 3000
	TX_KEY_FACTOR    = 100000
)

func getProposerAndValidatorsFromBlock(block *types.Block) (proposer string, validators string, err error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return "", "", nil
	}
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		return "", "", err
	}

	sigHash, err := sigHash(block.Header())
	if err != nil {
		return "", "", err
	}
	proposerAddr, err := istanbul.GetSignatureAddress(sigHash.Bytes(), istanbulExtra.Seal)
	if err != nil {
		return "", "", err
	}
	commiteeAddrs := make([]common.Address, len(istanbulExtra.Validators))
	for i, addr := range istanbulExtra.Validators {
		commiteeAddrs[i] = addr
	}
	var strValidators []string
	for _, validator := range istanbulExtra.Validators {
		strValidators = append(strValidators, validator.Hex())
	}

	return proposerAddr.Hex(), strings.Join(strValidators, ","), nil
}

// ecrecover extracts the Klaytn account address from a signed header.
func ecrecover(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	sigHash, err := sigHash(header)
	if err != nil {
		return common.Address{}, err
	}
	addr, err := istanbul.GetSignatureAddress(sigHash.Bytes(), istanbulExtra.Seal)
	if err != nil {
		return addr, err
	}
	return addr, nil
}

func sigHash(header *types.Header) (hash common.Hash, err error) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	if err := rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false)); err != nil {
		logger.Error("fail to encode", "err", err)
		return common.Hash{}, err
	}
	hasher.Sum(hash[:0])
	return hash, nil
}
