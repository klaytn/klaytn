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

package kas

import (
	"strings"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestSplitToWords_Success(t *testing.T) {
	data := "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000850f0263a87af6dd51acb8baab96219041e28fda00000000000000000000000000000000000000000000d3c21bcecceda1000000"
	bytes, err := hexutil.Decode(data)
	assert.NoError(t, err)

	hashes, err := splitToWords(bytes)
	assert.NoError(t, err)
	assert.Equal(t, common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), hashes[0])
	assert.Equal(t, common.HexToHash("0x000000000000000000000000850f0263a87af6dd51acb8baab96219041e28fda"), hashes[1])
	assert.Equal(t, common.HexToHash("0x00000000000000000000000000000000000000000000d3c21bcecceda1000000"), hashes[2])
}

func TestSplitToWords_Fail_DataLengthError(t *testing.T) {
	data := "0x0000000000000000000000000000000000850f0263a87af6dd51acb8baab96219041e28fda00000000000000000000000000000000000000000000d3c21bcecceda1000000"
	bytes, err := hexutil.Decode(data)
	assert.NoError(t, err)

	_, err = splitToWords(bytes)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "data length is not valid"))
}

func TestWordsToAddress_Success(t *testing.T) {
	hash1 := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	hash2 := common.HexToHash("0x000000000000000000000000850f0263a87af6dd51acb8baab96219041e28fda")
	hash3 := common.HexToHash("0x00000000000000000000000000000000000000000000d3c21bcecceda1000000")

	expected1 := common.HexToAddress("0x0000000000000000000000000000000000000000")
	expected2 := common.HexToAddress("0x850f0263a87af6dd51acb8baab96219041e28fda")
	expected3 := common.HexToAddress("0x00000000000000000000d3c21bcecceda1000000")

	addr1 := wordToAddress(hash1)
	addr2 := wordToAddress(hash2)
	addr3 := wordToAddress(hash3)

	assert.Equal(t, expected1, addr1)
	assert.Equal(t, expected2, addr2)
	assert.Equal(t, expected3, addr3)
}
