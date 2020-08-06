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
	"encoding/json"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fastCallTracerResult = []byte(`{
      "from": "0x0000000000000000000000000000000000000000",
      "gas": 0,
      "gasUsed": 0,
      "input": "",
      "output": "",
      "reverted": {
        "contract": "0x0000000000000000000000000000000000000000",
        "message": ""
      },
      "time": 0,
      "to": "0x0000000000000000000000000000000000000000",
      "type": "",
      "value": "0x0"
}`)

func TestIsInternalTxResult_Success(t *testing.T) {
	var testResult vm.InternalTxTrace
	assert.True(t, json.Valid(fastCallTracerResult))
	assert.NoError(t, json.Unmarshal(fastCallTracerResult, &testResult))
	assert.True(t, isEmptyTraceResult(&testResult))
}
