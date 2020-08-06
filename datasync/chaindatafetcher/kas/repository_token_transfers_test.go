package kas

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitToWords_Success(t *testing.T) {
	data := "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000850f0263a87af6dd51acb8baab96219041e28fda00000000000000000000000000000000000000000000d3c21bcecceda1000000"
	bytes, err := hexutil.Decode(data)
	assert.NoError(t, err)

	t.Log(bytes)
	hashes, _ := splitToWords(bytes)
	assert.Equal(t, common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), hashes[0])
	assert.Equal(t, common.HexToHash("0x000000000000000000000000850f0263a87af6dd51acb8baab96219041e28fda"), hashes[1])
	assert.Equal(t, common.HexToHash("0x00000000000000000000000000000000000000000000d3c21bcecceda1000000"), hashes[2])
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
