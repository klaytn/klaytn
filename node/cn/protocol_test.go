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

package cn

import (
	"bytes"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

func TestHashOrNumber_Success(t *testing.T) {
	// Hash is used for hashOrNumber.
	{
		data := &hashOrNumber{Hash: common.HexToHash("111")}
		b := new(bytes.Buffer)
		assert.NoError(t, rlp.Encode(b, data))

		decodedHashData := &hashOrNumber{}
		assert.NoError(t, rlp.DecodeBytes(b.Bytes(), decodedHashData))
		assert.EqualValues(t, data, decodedHashData)
	}
	// Number is used for hashOrNumber.
	{
		data := &hashOrNumber{Number: uint64(111)}
		b := new(bytes.Buffer)
		assert.NoError(t, rlp.Encode(b, data))

		decodedHashData := &hashOrNumber{}
		assert.NoError(t, rlp.DecodeBytes(b.Bytes(), decodedHashData))
		assert.EqualValues(t, data, decodedHashData)
	}
}

func TestHashOrNumber_Fail(t *testing.T) {
	// If both hash and number is provided, it should throw an error.
	{
		data := &hashOrNumber{Hash: common.HexToHash("111"), Number: uint64(111)}
		b := new(bytes.Buffer)
		assert.Error(t, rlp.Encode(b, data))
	}
	// If trying to decode invalid type into hashOrNumber, it should throw an error.
	{
		data := &basePeer{id: "TestBasePeer"}
		b := new(bytes.Buffer)
		assert.NoError(t, rlp.Encode(b, data))

		decodedHashData := &hashOrNumber{}
		assert.Error(t, rlp.DecodeBytes(b.Bytes(), decodedHashData))
	}
}
