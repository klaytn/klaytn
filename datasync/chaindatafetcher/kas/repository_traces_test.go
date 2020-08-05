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
	assert.False(t, isInternalTxResult(&testResult))
}
