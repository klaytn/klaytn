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
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// TestEncodeAndDecode_Receipt_Slice_StreamDecoding checks if encoding and decoding works
// as expected for types.Receipt by decoding with rlp.NewStream. rlp.NewStream is needed when you
// want to 1) limit the amount of data to be decoded (e.g., first 100 receipts from given receipts) or
// 2) apply different decoding method to given byte slice (e.g., first decode receipt and second decode block)
// TxHash should be an empty hash for Receipt encoding and decoding.
func TestEncodeAndDecode_Receipt_Slice_StreamDecoding(t *testing.T) {
	rct := &Receipt{}
	rct.TxHash = common.BigToHash(big.NewInt(12345))
	rct.GasUsed = uint64(12345)
	rct.Status = ReceiptStatusSuccessful

	var rctList []*Receipt
	rctList = append(rctList, rct)

	size, r, err := rlp.EncodeToReader(rctList)
	if err != nil {
		t.Fatalf("Error while rlp.EncodeToReader. err: %v\n", err)
	}

	msgStream := rlp.NewStream(r, uint64(size))
	if _, err := msgStream.List(); err != nil {
		t.Fatalf("Error while msgStream.List(). err: %v\n", err)
	}

	var rctEncAndDec *Receipt
	if err := msgStream.Decode(&rctEncAndDec); err != nil {
		t.Fatalf("Error while msgStream.Decode(). err: %v\n", err)
	}

	assert.Equal(t, rct.Status, rctEncAndDec.Status)
	assert.Equal(t, rct.GasUsed, rctEncAndDec.GasUsed)
	assert.Equal(t, common.Hash{}, rctEncAndDec.TxHash) // TxHash should be an empty hash.
}

// TestEncodeAndDecode_Receipt_Slice_DirectDecoding checks if encoding and decoding works
// as expected for types.Receipt by decoding with rlp.DecodeBytes. TxHash should be an empty
// hash for Receipt encoding and decoding.
func TestEncodeAndDecode_Receipt_Slice_DirectDecoding(t *testing.T) {
	rct := &Receipt{}
	rct.TxHash = common.BigToHash(big.NewInt(12345))
	rct.GasUsed = uint64(12345)
	rct.Status = ReceiptStatusSuccessful

	var rctList []*Receipt
	rctList = append(rctList, rct)

	bytes, err := rlp.EncodeToBytes(rctList)
	if err != nil {
		t.Fatalf("Error while rlp.EncodeToReader. err: %v\n", err)
	}

	var encAndDecReceipts []*Receipt
	rlp.DecodeBytes(bytes, &encAndDecReceipts)

	assert.Equal(t, rct.Status, encAndDecReceipts[0].Status)
	assert.Equal(t, rct.GasUsed, encAndDecReceipts[0].GasUsed)
	assert.Equal(t, common.Hash{}, encAndDecReceipts[0].TxHash) // TxHash should be an empty hash.
}

// TestEncodeAndDecode_ReceiptForStorage_Slice_StreamDecoding checks if encoding and decoding works
// as expected for types.ReceiptForStorage by decoding with rlp.NewStream. rlp.NewStream is needed when you
// want to 1) limit the amount of data to be decoded (e.g., first 100 receipts from given receipts) or
// 2) apply different decoding method to given byte slice (e.g., first decode receipt and second decode block)
// TxHash must appear for ReceiptForStorage encoding and decoding.
func TestEncodeAndDecode_ReceiptForStorage_Slice_StreamDecoding(t *testing.T) {
	rct := &Receipt{}
	rct.TxHash = common.BigToHash(big.NewInt(12345))
	rct.GasUsed = uint64(12345)
	rct.Status = ReceiptStatusSuccessful

	// ReceiptForStorage is a simple wrapper to change the way of RLP encoding/decoding.
	rctfs := (*ReceiptForStorage)(rct)

	var rctfsList []*ReceiptForStorage
	rctfsList = append(rctfsList, rctfs)

	size, r, err := rlp.EncodeToReader(rctfsList)
	if err != nil {
		t.Fatalf("Error while rlp.EncodeToReader. err: %v\n", err)
	}

	msgStream := rlp.NewStream(r, uint64(size))
	if _, err := msgStream.List(); err != nil {
		t.Fatalf("Error while msgStream.List(). err: %v\n", err)
	}

	var rctfsEndAndDec *ReceiptForStorage
	if err := msgStream.Decode(&rctfsEndAndDec); err != nil {
		t.Fatalf("Error while msgStream.Decode(). err: %v\n", err)
	}

	assert.Equal(t, rct.Status, rctfsEndAndDec.Status)
	assert.Equal(t, rct.GasUsed, rctfsEndAndDec.GasUsed)
	assert.Equal(t, rct.TxHash, rctfsEndAndDec.TxHash) // TxHash should be equal to the original one.
}

// TestEncodeAndDecode_ReceiptForStorage_Slice_DirectDecoding checks if encoding and decoding works
// as expected for types.ReceiptForStorage by decoding with rlp.DecodeBytes. TxHash must appear for
// ReceiptForStorage encoding and decoding.
func TestEncodeAndDecode_ReceiptForStorage_Slice_DirectDecoding(t *testing.T) {
	rct := &Receipt{}
	rct.TxHash = common.BigToHash(big.NewInt(12345))
	rct.GasUsed = uint64(12345)
	rct.Status = ReceiptStatusSuccessful

	// ReceiptForStorage is a simple wrapper to change the way of RLP encoding/decoding.
	rctfs := (*ReceiptForStorage)(rct)

	var rctfsList []*ReceiptForStorage
	rctfsList = append(rctfsList, rctfs)

	bytes, err := rlp.EncodeToBytes(rctfsList)
	if err != nil {
		t.Fatalf("Error while rlp.EncodeToBytes. err: %v\n", err)
	}

	var encAndDecReceiptForStorage []*ReceiptForStorage
	rlp.DecodeBytes(bytes, &encAndDecReceiptForStorage)

	assert.Equal(t, rct.Status, encAndDecReceiptForStorage[0].Status)
	assert.Equal(t, rct.GasUsed, encAndDecReceiptForStorage[0].GasUsed)
	assert.Equal(t, rct.TxHash, encAndDecReceiptForStorage[0].TxHash) // TxHash should be equal to the original one.
}
