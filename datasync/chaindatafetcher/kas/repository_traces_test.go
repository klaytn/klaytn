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
	"testing"

	"github.com/klaytn/klaytn/common"

	"github.com/klaytn/klaytn/blockchain/types"

	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/stretchr/testify/assert"
)

// A specific types of transaction returns empty trace result which is defined as a variable `emptyTraceResult`.
// As a result of the tracing, `reflect.DeepEqual` returns an error while comparing with the not-initialized slice.
func TestRepository_isEmptyTraceResult(t *testing.T) {
	// right empty result.
	data := &vm.InternalTxTrace{
		Value: "0x0",
		Calls: []*vm.InternalTxTrace{},
	}
	assert.True(t, isEmptyTraceResult(data))

	// wrong empty result.
	data = &vm.InternalTxTrace{
		Value: "0x0",
	}
	assert.False(t, isEmptyTraceResult(data))
}

func makeOffset(offset int64) *int64 {
	return &offset
}

func makeEntryTx() *Tx {
	txhash := genRandomHash()
	return &Tx{
		TransactionId:   100000000000,
		TransactionHash: txhash.Bytes(),
		Status:          int(types.ReceiptStatusSuccessful),
		Timestamp:       1,
		TypeInt:         int(types.TxTypeLegacyTransaction),
	}
}

func makeInternalTrace(callType, value string, from *common.Address, to *common.Address) *vm.InternalTxTrace {
	return &vm.InternalTxTrace{
		Type:  callType,
		From:  from,
		To:    to,
		Value: value,
	}
}

func makeExpectedInternalTx(offset int64, entryTx *Tx, trace *vm.InternalTxTrace) *Tx {
	return &Tx{
		TransactionId:   entryTx.TransactionId + offset,
		FromAddr:        trace.From.Bytes(),
		ToAddr:          trace.To.Bytes(),
		Value:           trace.Value,
		TransactionHash: entryTx.TransactionHash,
		Status:          entryTx.Status,
		Timestamp:       entryTx.Timestamp,
		TypeInt:         entryTx.TypeInt,
		Internal:        true,
	}
}

func TestRepository_transformToInternalTx(t *testing.T) {
	type args struct {
		trace       *vm.InternalTxTrace
		offset      *int64
		entryTx     *Tx
		isFirstCall bool
	}

	// valid test case
	entryTx := makeEntryTx()
	trace := makeInternalTrace("TEST", "0x1", genRandomAddress(), genRandomAddress())
	args1 := args{trace, makeOffset(1), entryTx, false}
	expected1 := []*Tx{makeExpectedInternalTx(2, entryTx, trace)}

	// valid test case 2
	entryTx2 := makeEntryTx()
	trace2 := makeInternalTrace("TEST", "0x1", genRandomAddress(), genRandomAddress())
	innerTrace := makeInternalTrace("TEST", "0x2", genRandomAddress(), genRandomAddress())
	trace2.Calls = []*vm.InternalTxTrace{innerTrace}
	args2 := args{trace2, makeOffset(0), entryTx2, false}

	expected2 := []*Tx{
		makeExpectedInternalTx(1, entryTx2, trace2),
		makeExpectedInternalTx(2, entryTx2, innerTrace),
	}

	tests := []struct {
		name     string
		args     args
		expected []*Tx
		err      error
	}{
		{
			name:     "success_valid_internal_tx",
			args:     args1,
			expected: expected1,
			err:      nil,
		},
		{
			name:     "success_valid_internal_tx_with_inner_calls",
			args:     args2,
			expected: expected2,
			err:      nil,
		},
		{
			name: "fail_noOpcodeError",
			args: args{
				trace:       makeInternalTrace("", "0x1", genRandomAddress(), genRandomAddress()),
				offset:      makeOffset(0),
				entryTx:     makeEntryTx(),
				isFirstCall: false,
			},
			expected: nil,
			err:      noOpcodeError,
		},
		{
			name: "fail_noFromFieldError",
			args: args{
				trace:       makeInternalTrace("TEST", "0x1", nil, genRandomAddress()),
				offset:      makeOffset(0),
				entryTx:     makeEntryTx(),
				isFirstCall: false,
			},
			expected: nil,
			err:      noFromFieldError,
		},
		{
			name: "fail_noToFieldError",
			args: args{
				trace:       makeInternalTrace("TEST", "0x1", genRandomAddress(), nil),
				offset:      makeOffset(0),
				entryTx:     makeEntryTx(),
				isFirstCall: false,
			},
			expected: nil,
			err:      noToFieldError,
		},
		{
			name: "success_ignore_selfdestruct",
			args: args{
				trace:       makeInternalTrace(selfDestructType, "0x1", genRandomAddress(), genRandomAddress()),
				offset:      makeOffset(0),
				entryTx:     makeEntryTx(),
				isFirstCall: false,
			},
			expected: nil,
			err:      nil,
		},
		{
			name: "success_ignore_firstCall",
			args: args{
				trace:       makeInternalTrace("TEST", "0x1", genRandomAddress(), genRandomAddress()),
				offset:      makeOffset(0),
				entryTx:     makeEntryTx(),
				isFirstCall: true,
			},
			expected: nil,
			err:      nil,
		},
		{
			name: "success_empty_value",
			args: args{
				trace:       makeInternalTrace("TEST", "", genRandomAddress(), genRandomAddress()),
				offset:      makeOffset(0),
				entryTx:     makeEntryTx(),
				isFirstCall: false,
			},
			expected: nil,
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformToInternalTx(tt.args.trace, tt.args.offset, tt.args.entryTx, tt.args.isFirstCall)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
