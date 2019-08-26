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

package types

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/ser/rlp"
	"math/big"
)

type ChainHashes struct {
	Type uint8
	Data []byte
}

type ChainHashesInternalType0 struct {
	BlockHash     common.Hash
	TxHash        common.Hash
	ParentHash    common.Hash
	ReceiptHash   common.Hash
	StateRootHash common.Hash
	BlockNumber   *big.Int
	Period        *big.Int
	TxCounts      *big.Int
}

func NewChainHashesType0(block *Block, period *big.Int, txCounts *big.Int) (*ChainHashes, error) {
	data := &ChainHashesInternalType0{block.Hash(), block.Header().TxHash,
		block.Header().ParentHash, block.Header().ReceiptHash,
		block.Header().Root, block.Header().Number, period, txCounts}
	encodedCCTxData, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil, err
	}
	return &ChainHashes{0, encodedCCTxData}, nil
}
