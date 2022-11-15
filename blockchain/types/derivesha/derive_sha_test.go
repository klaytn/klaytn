// Copyright 2022 The klaytn Authors
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
