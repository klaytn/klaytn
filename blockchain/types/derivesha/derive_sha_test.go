package derivesha

import (
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"gotest.tools/assert"
)

func TestEmptyRoot(t *testing.T) {
	assert.Equal(t,
		DeriveShaOrig{}.DeriveSha(types.Transactions{}).Hex(),
		"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	assert.Equal(t,
		DeriveShaSimple{}.DeriveSha(types.Transactions{}).Hex(),
		"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
	assert.Equal(t,
		DeriveShaConcat{}.DeriveSha(types.Transactions{}).Hex(),
		"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
}
