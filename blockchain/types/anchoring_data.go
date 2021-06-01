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
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
)

const (
	AnchoringDataType0    uint8 = 0
	AnchoringJSONDataType uint8 = 128
)

var (
	errUnknownAnchoringTxType = errors.New("unknown anchoring tx type")
)

type AnchoringDataInternal interface {
	GetBlockHash() common.Hash
	GetBlockNumber() *big.Int
}

type AnchoringData struct {
	Type uint8
	Data []byte
}

// AnchoringDataLegacy is an old anchoring type that does not support an data type.
type AnchoringDataLegacy struct {
	BlockHash     common.Hash `json:"blockHash"`
	TxHash        common.Hash `json:"transactionsRoot"`
	ParentHash    common.Hash `json:"parentHash"`
	ReceiptHash   common.Hash `json:"receiptsRoot"`
	StateRootHash common.Hash `json:"stateRoot"`
	BlockNumber   *big.Int    `json:"blockNumber"`
}

func (data *AnchoringDataLegacy) GetBlockHash() common.Hash {
	return data.BlockHash
}

func (data *AnchoringDataLegacy) GetBlockNumber() *big.Int {
	return data.BlockNumber
}

type AnchoringDataInternalType0 struct {
	BlockHash     common.Hash `json:"blockHash"`
	TxHash        common.Hash `json:"transactionsRoot"`
	ParentHash    common.Hash `json:"parentHash"`
	ReceiptHash   common.Hash `json:"receiptsRoot"`
	StateRootHash common.Hash `json:"stateRoot"`
	BlockNumber   *big.Int    `json:"blockNumber"`
	BlockCount    *big.Int    `json:"blockCount"`
	TxCount       *big.Int    `json:"txCount"`
}

func (data *AnchoringDataInternalType0) GetBlockHash() common.Hash {
	return data.BlockHash
}

func (data *AnchoringDataInternalType0) GetBlockNumber() *big.Int {
	return data.BlockNumber
}

func NewAnchoringDataType0(block *Block, blockCount uint64, txCount uint64) (*AnchoringData, error) {
	data := &AnchoringDataInternalType0{
		block.Hash(),
		block.Header().TxHash,
		block.Header().ParentHash,
		block.Header().ReceiptHash,
		block.Header().Root,
		block.Header().Number,
		new(big.Int).SetUint64(blockCount),
		new(big.Int).SetUint64(txCount),
	}
	encodedCCTxData, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil, err
	}
	return &AnchoringData{AnchoringDataType0, encodedCCTxData}, nil
}

func NewAnchoringJSONDataType(v interface{}) (*AnchoringData, error) {
	encoded, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &AnchoringData{AnchoringJSONDataType, encoded}, nil
}

// DecodeAnchoringData decodes an anchoring data used by main and sub bridges.
func DecodeAnchoringData(data []byte) (AnchoringDataInternal, error) {
	anchoringData := new(AnchoringData)
	if err := rlp.DecodeBytes(data, anchoringData); err != nil {
		anchoringDataLegacy := new(AnchoringDataLegacy)
		if err := rlp.DecodeBytes(data, anchoringDataLegacy); err != nil {
			return nil, err
		}
		logger.Trace("decoded legacy anchoring tx", "blockNum", anchoringDataLegacy.GetBlockNumber().String(), "blockHash", anchoringDataLegacy.GetBlockHash().String(), "txHash", anchoringDataLegacy.TxHash.String())
		return anchoringDataLegacy, nil
	}
	if anchoringData.Type == AnchoringDataType0 {
		anchoringDataInternal := new(AnchoringDataInternalType0)
		if err := rlp.DecodeBytes(anchoringData.Data, anchoringDataInternal); err != nil {
			return nil, err
		}
		logger.Trace("decoded type0 anchoring tx", "blockNum", anchoringDataInternal.BlockNumber.String(), "blockHash", anchoringDataInternal.BlockHash.String(), "txHash", anchoringDataInternal.TxHash.String(), "txCount", anchoringDataInternal.TxCount)
		return anchoringDataInternal, nil
	}
	return nil, errUnknownAnchoringTxType
}

// DecodeAnchoringDataToJSON decodes an anchoring data.
func DecodeAnchoringDataToJSON(data []byte) (interface{}, error) {
	anchoringData := new(AnchoringData)
	if err := rlp.DecodeBytes(data, anchoringData); err != nil {
		anchoringDataLegacy := new(AnchoringDataLegacy)
		if err := rlp.DecodeBytes(data, anchoringDataLegacy); err != nil {
			return nil, err
		}
		return anchoringDataLegacy, nil
	}
	if anchoringData.Type == AnchoringDataType0 {
		anchoringDataInternal := new(AnchoringDataInternalType0)
		if err := rlp.DecodeBytes(anchoringData.Data, anchoringDataInternal); err != nil {
			return nil, err
		}
		return anchoringDataInternal, nil
	}
	if anchoringData.Type == AnchoringJSONDataType {
		var v map[string]interface{}
		if err := json.Unmarshal(anchoringData.Data, &v); err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, fmt.Errorf("%w type=%v", errUnknownAnchoringTxType, anchoringData.Type)
}
