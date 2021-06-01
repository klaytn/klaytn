// Modifications Copyright 2020 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// InternalTxTracer is a full blown transaction tracer that extracts and reports all
// the internal calls made by a transaction, along with any useful information.
//
// This file is derived from eth/tracers/internal/tracers/call_tracer.js (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/klaytn/klaytn/common"
)

func jsonMustCompact(data []byte) []byte {
	compactedBuffer := new(bytes.Buffer)
	err := json.Compact(compactedBuffer, data)
	if err != nil {
		panic(err)
	}
	return compactedBuffer.Bytes()
}

func genAddressPtr(addr string) *common.Address {
	ret := common.HexToAddress(addr)
	return &ret
}

func TestInternalTxTrace_MarshalJSON(t *testing.T) {
	type fields struct {
		Type     string
		From     *common.Address
		To       *common.Address
		Value    string
		Gas      uint64
		GasUsed  uint64
		Input    string
		Output   string
		Error    error
		Time     time.Duration
		Calls    []*InternalTxTrace
		Reverted *RevertedInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "revert_test",
			fields: fields{
				Type:    CALL.String(),
				From:    genAddressPtr("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"),
				To:      genAddressPtr("0xabbcd5b340c80b5f1c0545c04c987b87310296ae"),
				Value:   "0x0",
				Gas:     2971112,
				GasUsed: 195,
				Input:   "0x73b40a5c000000000000000000000000400de2e016bda6577407dfc379faba9899bc73ef0000000000000000000000002cc31912b2b0f3075a87b3640923d45a26cef3ee000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000064d79d8e6c7265636f76657279416464726573730000000000000000000000000000000000000000000000000000000000383e3ec32dc0f66d8fe60dbdc2f6815bdf73a988383e3ec32dc0f66d8fe60dbdc2f6815bdf73a98800000000000000000000000000000000000000000000000000000000000000000000000000000000",
				Error:   errExecutionReverted,
				Reverted: &RevertedInfo{
					Contract: genAddressPtr("0xabbcd5b340c80b5f1c0545c04c987b87310296ae"),
				},
			},
			want: jsonMustCompact([]byte(`{
	"type": "CALL",
	"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
	"to": "0xabbcd5b340c80b5f1c0545c04c987b87310296ae",
	"value": "0x0",
    "gas": "0x2d55e8",
    "gasUsed": "0xc3",
    "input": "0x73b40a5c000000000000000000000000400de2e016bda6577407dfc379faba9899bc73ef0000000000000000000000002cc31912b2b0f3075a87b3640923d45a26cef3ee000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000064d79d8e6c7265636f76657279416464726573730000000000000000000000000000000000000000000000000000000000383e3ec32dc0f66d8fe60dbdc2f6815bdf73a988383e3ec32dc0f66d8fe60dbdc2f6815bdf73a98800000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "error": "execution reverted",
    "reverted": {
      "contract": "0xabbcd5b340c80b5f1c0545c04c987b87310296ae"
    }
  }`)),
		},
		{
			name: "delegatecall_test",
			fields: fields{
				Type:    DELEGATECALL.String(),
				From:    genAddressPtr("0x2074f91470365144df78c333948814fe4158af5a"),
				To:      genAddressPtr("0x13cf0d18bcb898efd4e85ae2bb65a443ab86023c"),
				Gas:     179,
				GasUsed: 122,
				Input:   "0xc605f76c",
				Output:  "0x",
			},
			want: jsonMustCompact([]byte(`{
		"type": "DELEGATECALL",
        "from": "0x2074f91470365144df78c333948814fe4158af5a",
		"to": "0x13cf0d18bcb898efd4e85ae2bb65a443ab86023c",
        "gas": "0xb3",
        "gasUsed": "0x7a",
        "input": "0xc605f76c",
        "output": "0x"
}`)),
		},
		{
			name: "create2_test",
			fields: fields{
				Type:    CREATE2.String(),
				From:    genAddressPtr("0xe213d8b68ca3d01e51a6dba669de59ac9a8359ee"),
				To:      genAddressPtr("0xd3204138864ad97dbfe703cf7281f6b3397c6003"),
				Value:   "0x0",
				Gas:     11341,
				GasUsed: 10784,
				Input:   "0x6080604052348015600f57600080fd5b50604051602080608183398101806040526020811015602d57600080fd5b810190808051906020019092919050505050603580604c6000396000f3fe6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029",
				Output:  "0x6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029",
			},
			want: jsonMustCompact([]byte(`{
		"type": "CREATE2",
        "from": "0xe213d8b68ca3d01e51a6dba669de59ac9a8359ee",
		"to": "0xd3204138864ad97dbfe703cf7281f6b3397c6003",
		"value": "0x0",
        "gas": "0x2c4d",
        "gasUsed": "0x2a20",
        "input": "0x6080604052348015600f57600080fd5b50604051602080608183398101806040526020811015602d57600080fd5b810190808051906020019092919050505050603580604c6000396000f3fe6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029",
        "output": "0x6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029"
}`)),
		},
		{
			name: "create_test",
			fields: fields{
				Type:    CREATE.String(),
				From:    genAddressPtr("0x13e4acefe6a6700604929946e70e6443e4e73447"),
				To:      genAddressPtr("0x7dc9c9730689ff0b0fd506c67db815f12d90a448"),
				Value:   "0x0",
				Gas:     385286,
				GasUsed: 385286,
				Input:   "0x6080604052348015600f57600080fd5b50604051602080608183398101806040526020811015602d57600080fd5b810190808051906020019092919050505050603580604c6000396000f3fe6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029",
				Output:  "0x6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029",
			},
			want: jsonMustCompact([]byte(`{
    "type": "CREATE",
    "from": "0x13e4acefe6a6700604929946e70e6443e4e73447",
	"to": "0x7dc9c9730689ff0b0fd506c67db815f12d90a448",
	"value": "0x0",
    "gas": "0x5e106",
    "gasUsed": "0x5e106",
    "input": "0x6080604052348015600f57600080fd5b50604051602080608183398101806040526020811015602d57600080fd5b810190808051906020019092919050505050603580604c6000396000f3fe6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029",
    "output": "0x6080604052600080fdfea165627a7a72305820c577c5325e4c667e8439f998d616e74110a14fb5b185afe46d8443bbfd9c3d710029"
}`)),
		},
		{
			name: "call_test",
			fields: fields{
				Type:    CALL.String(),
				From:    genAddressPtr("0x3693da93b9d5e63cb4804b8813c8214a76724485"),
				To:      genAddressPtr("0x13cf0d18bcb898efd4e85ae2bb65a443ab86023c"),
				Value:   "0x0",
				Gas:     179,
				GasUsed: 122,
				Input:   "0xc605f76c",
				Output:  "0x",
			},
			want: jsonMustCompact([]byte(`{
	"type": "CALL",
	"from": "0x3693da93b9d5e63cb4804b8813c8214a76724485",
	"to": "0x13cf0d18bcb898efd4e85ae2bb65a443ab86023c",
	"value": "0x0",
	"gas": "0xb3",
	"gasUsed": "0x7a",
	"input": "0xc605f76c",
	"output": "0x"
}`)),
		},
		{
			name: "simple_call_test",
			fields: fields{
				Type:  CALL.String(),
				From:  genAddressPtr("0x3b873a919aa0512d5a0f09e6dcceaa4a6727fafe"),
				To:    genAddressPtr("0x0024f658a46fbb89d8ac105e98d7ac7cbbaf27c5"),
				Value: "0x6f05b59d3b20000",
				Input: "0x",
			},
			want: jsonMustCompact([]byte(`{
		"type": "CALL",
        "from": "0x3b873a919aa0512d5a0f09e6dcceaa4a6727fafe",
		"to": "0x0024f658a46fbb89d8ac105e98d7ac7cbbaf27c5",
		"value": "0x6f05b59d3b20000",
        "input": "0x"
}`)),
		},
		{
			name: "staticcall_test",
			fields: fields{
				Type:    OpCode(STATICCALL).String(),
				From:    genAddressPtr("0xbadc777c579b497ede07fa6ff93bdf4e31793f24"),
				To:      genAddressPtr("0x920fd5070602feaea2e251e9e7238b6c376bcae5"),
				Gas:     179,
				GasUsed: 122,
				Input:   "0xc605f76c",
				Output:  "0x",
			},
			want: jsonMustCompact([]byte(`{
		"type": "STATICCALL",
        "from": "0xbadc777c579b497ede07fa6ff93bdf4e31793f24",
		"to": "0x920fd5070602feaea2e251e9e7238b6c376bcae5",
        "gas": "0xb3",
        "gasUsed": "0x7a",
        "input": "0xc605f76c",
        "output": "0x"
}`)),
		},
		{
			name: "selfdestruct_test",
			fields: fields{
				Type: OpCode(SELFDESTRUCT).String(),
			},
			want: []byte(`{"type":"SELFDESTRUCT"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := InternalTxTrace{
				Type:     tt.fields.Type,
				From:     tt.fields.From,
				To:       tt.fields.To,
				Value:    tt.fields.Value,
				Gas:      tt.fields.Gas,
				GasUsed:  tt.fields.GasUsed,
				Input:    tt.fields.Input,
				Output:   tt.fields.Output,
				Error:    tt.fields.Error,
				Time:     tt.fields.Time,
				Calls:    tt.fields.Calls,
				Reverted: tt.fields.Reverted,
			}
			got, err := i.MarshalJSON()
			assert.Nil(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
