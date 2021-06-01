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

package types

import (
	"encoding/json"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

func genRandomAddress() *common.Address {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return &addr
}

func genRandomHash() (h common.Hash) {
	rand.Seed(time.Now().UnixNano())
	hasher := sha3.NewKeccak256()

	r := rand.Uint64()
	rlp.Encode(hasher, r)
	hasher.Sum(h[:0])

	return h
}

func TestDecodingLegacyAnchoringTx(t *testing.T) {
	anchoringData := &AnchoringDataLegacy{
		BlockHash:     genRandomHash(),
		TxHash:        genRandomHash(),
		ParentHash:    genRandomHash(),
		ReceiptHash:   genRandomHash(),
		StateRootHash: genRandomHash(),
		BlockNumber:   new(big.Int).SetUint64(rand.Uint64()),
	}
	data, err := rlp.EncodeToBytes(anchoringData)
	assert.NoError(t, err)

	// Decoding the anchoring tx.
	decodedData, err := DecodeAnchoringData(data)
	assert.NoError(t, err)
	assert.Equal(t, anchoringData, decodedData)

	decodedDataJSON, err := DecodeAnchoringDataToJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, anchoringData, decodedDataJSON)
}

func TestDecodingAnchoringTxType0(t *testing.T) {
	block := genBlock()
	blockCnt := rand.Uint64()
	txCount := rand.Uint64()
	anchoringData, err := NewAnchoringDataType0(block, blockCnt, txCount)
	assert.NoError(t, err)

	data, err := rlp.EncodeToBytes(anchoringData)
	assert.NoError(t, err)

	// Decoding the anchoring tx.
	decodedData, err := DecodeAnchoringData(data)
	assert.NoError(t, err)
	decodedInternalData, ok := decodedData.(*AnchoringDataInternalType0)
	assert.True(t, ok)
	assert.Equal(t, txCount, decodedInternalData.TxCount.Uint64())
	assert.Equal(t, blockCnt, decodedInternalData.BlockCount.Uint64())
	assert.Equal(t, block.header.Root, decodedInternalData.StateRootHash)
	assert.Equal(t, block.header.TxHash, decodedInternalData.TxHash)
	assert.Equal(t, block.header.ReceiptHash, decodedInternalData.ReceiptHash)
	assert.Equal(t, block.header.Number, decodedInternalData.BlockNumber)
	assert.Equal(t, block.header.Hash(), decodedInternalData.BlockHash)
	assert.Equal(t, block.header.ParentHash, decodedInternalData.ParentHash)

	decodedDataJSON, err := DecodeAnchoringDataToJSON(data)
	assert.NoError(t, err)
	decodedInternalData, ok = decodedDataJSON.(*AnchoringDataInternalType0)
	assert.True(t, ok)
	assert.Equal(t, txCount, decodedInternalData.TxCount.Uint64())
	assert.Equal(t, blockCnt, decodedInternalData.BlockCount.Uint64())
	assert.Equal(t, block.header.Root, decodedInternalData.StateRootHash)
	assert.Equal(t, block.header.TxHash, decodedInternalData.TxHash)
	assert.Equal(t, block.header.ReceiptHash, decodedInternalData.ReceiptHash)
	assert.Equal(t, block.header.Number, decodedInternalData.BlockNumber)
	assert.Equal(t, block.header.Hash(), decodedInternalData.BlockHash)
	assert.Equal(t, block.header.ParentHash, decodedInternalData.ParentHash)

	var result AnchoringDataInternalType0
	err = rlp.DecodeBytes(anchoringData.Data, &result)
	assert.NoError(t, err)
	expResult, err := json.Marshal(result)

	actResult, err := json.Marshal(decodedInternalData)
	assert.NoError(t, err)
	assert.Equal(t, expResult, actResult)
}

func TestDecodingAnchoringTxJSONType(t *testing.T) {
	originalData := map[string]interface{}{
		"int":    1,
		"string": "string",
		"bigInt": big.NewInt(100),
		"hash":   genRandomHash(),
		"addr":   genRandomAddress(),
	}

	anchoringData, err := NewAnchoringJSONDataType(originalData)
	assert.NoError(t, err)

	data, err := rlp.EncodeToBytes(anchoringData)
	assert.NoError(t, err)

	// Decoding the anchoring tx.
	decodedDataJSON, err := DecodeAnchoringDataToJSON(data)
	assert.NoError(t, err)

	expResult, err := json.Marshal(originalData)
	assert.NoError(t, err)
	actResult, err := json.Marshal(decodedDataJSON)
	assert.NoError(t, err)
	assert.Equal(t, expResult, actResult)
}
