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
	"reflect"
	"testing"

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

func makeInternalTrace() *vm.InternalTxTrace {
	from := genRandomAddress()
	to := genRandomAddress()
	return &vm.InternalTxTrace{
		Type:  "TEST",
		From:  &from,
		To:    &to,
		Value: "0x1",
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

func TestRepository_transformToInternalTx_Success_ValidInternalTx(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()
	expected := makeExpectedInternalTx(2, entryTx, trace)

	offset := int64(1)
	txs, err := transformToInternalTx(trace, &offset, entryTx, false)
	assert.NoError(t, err)
	assert.True(t, len(txs) == 1)
	assert.True(t, reflect.DeepEqual(txs[0], expected))
}

func TestRepository_transformToInternalTx_Success_ValidInternalTxWithInnerTrace(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()
	innerTrace := makeInternalTrace()
	innerTrace.Value = "0x2"
	trace.Calls = []*vm.InternalTxTrace{innerTrace}

	expected := makeExpectedInternalTx(1, entryTx, trace)
	expected2 := makeExpectedInternalTx(2, entryTx, innerTrace)

	offset := int64(0)
	txs, err := transformToInternalTx(trace, &offset, entryTx, false)
	assert.NoError(t, err)
	assert.True(t, len(txs) == 2)
	assert.True(t, reflect.DeepEqual(txs[0], expected))
	assert.True(t, reflect.DeepEqual(txs[1], expected2))
}

func TestRepository_transformToInternalTx_Fail_NoOpcode(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()
	trace.Type = ""

	offset := int64(0)
	_, err := transformToInternalTx(trace, &offset, entryTx, false)
	assert.Error(t, err)
	assert.Equal(t, noOpcodeError, err)
}

func TestRepository_transformToInternalTx_Fail_NoFromField(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()
	trace.From = nil

	offset := int64(0)
	_, err := transformToInternalTx(trace, &offset, entryTx, false)
	assert.Error(t, err)
	assert.Equal(t, noFromFieldError, err)
}

func TestRepository_transformToInternalTx_Fail_NoToField(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()
	trace.To = nil

	offset := int64(0)
	_, err := transformToInternalTx(trace, &offset, entryTx, false)
	assert.Error(t, err)
	assert.Equal(t, noToFieldError, err)
}

func TestRepository_transformToInternalTx_Success_SELFDESTRUCT(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()
	trace.Type = selfDestructType

	offset := int64(0)
	txs, err := transformToInternalTx(trace, &offset, entryTx, false)
	assert.Nil(t, txs)
	assert.Nil(t, err)
}

func TestRepository_transformToInternalTx_Success_FirstCall(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()

	offset := int64(0)
	txs, err := transformToInternalTx(trace, &offset, entryTx, true)
	assert.Nil(t, txs)
	assert.Nil(t, err)
}

func TestRepository_transformToInternalTx_Success_EmptyValue(t *testing.T) {
	entryTx := makeEntryTx()
	trace := makeInternalTrace()
	trace.Value = ""

	offset := int64(0)
	txs, err := transformToInternalTx(trace, &offset, entryTx, true)
	assert.Nil(t, txs)
	assert.Nil(t, err)
}
