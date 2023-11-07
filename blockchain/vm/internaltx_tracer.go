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
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/holiman/uint256"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

var (
	errExecutionReverted = errors.New("execution reverted")
	errInternalFailure   = errors.New("internal failure")
	emptyAddr            = common.Address{}
)

// InternalTxTracer is a full blown transaction tracer that extracts and reports all
// the internal calls made by a transaction, along with any useful information.
// It is ported to golang from JS, specifically call_tracer.js
type InternalTxTracer struct {
	callStack []*InternalCall
	output    []byte
	err       error
	errValue  string

	// Below are newly added fields to support call_tracer.js
	descended        bool
	revertedContract common.Address
	ctx              map[string]interface{} // Transaction context gathered throughout execution
	initialized      bool
	revertString     string

	interrupt uint32 // Atomic flag to signal execution interruption
	reason    error  // Textual reason for the interruption

	gasLimit uint64 // Amount of gas bought for the whole tx
}

// NewInternalTxTracer returns a new InternalTxTracer.
func NewInternalTxTracer() *InternalTxTracer {
	logger := &InternalTxTracer{
		callStack: []*InternalCall{{}},
		ctx:       map[string]interface{}{},
	}
	return logger
}

// InternalCall is emitted to the EVM each cycle and
// lists information about the current internal state
// prior to the execution of the statement.
type InternalCall struct {
	Type  string          `json:"type"`
	From  *common.Address `json:"from"`
	To    *common.Address `json:"to"`
	Value string          `json:"value"`

	Gas     uint64 `json:"gas"`
	GasIn   uint64 `json:"gasIn"`
	GasUsed uint64 `json:"gasUsed"`
	GasCost uint64 `json:"gasCost"`

	Input  string `json:"input"`  // hex string
	Output string `json:"output"` // hex string
	Error  error  `json:"err"`

	OutOff *big.Int `json:"outoff"`
	OutLen *big.Int `json:"outlen"`

	Calls []*InternalCall `json:"calls"`
}

// OpName formats the operand name in a human-readable format.
func (s *InternalCall) OpName() string {
	return s.Type
}

// ErrorString formats the tracerLog's error as a string.
func (s *InternalCall) ErrorString() string {
	if s.Error != nil {
		return s.Error.Error()
	}
	return ""
}

func (s *InternalCall) ToTrace() *InternalTxTrace {
	nestedCalls := []*InternalTxTrace{}
	for _, call := range s.Calls {
		nestedCalls = append(nestedCalls, call.ToTrace())
	}

	return &InternalTxTrace{
		Type:  s.Type,
		From:  s.From,
		To:    s.To,
		Value: s.Value,

		Gas:     s.Gas,
		GasUsed: s.GasUsed,

		Input:  s.Input,
		Output: s.Output,
		Error:  s.Error,

		Calls: nestedCalls,
	}
}

// InternalTxTrace is returned data after the end of trace-collecting cycle.
// It implements an object returned by "result" function at call_tracer.js
type InternalTxTrace struct {
	Type  string          `json:"type"`
	From  *common.Address `json:"from,omitempty"`
	To    *common.Address `json:"to,omitempty"`
	Value string          `json:"value,omitempty"`

	Gas     uint64 `json:"gas,omitempty"`
	GasUsed uint64 `json:"gasUsed,omitempty"`

	Input  string `json:"input,omitempty"`  // hex string
	Output string `json:"output,omitempty"` // hex string
	Error  error  `json:"error,omitempty"`

	Time  time.Duration      `json:"time,omitempty"`
	Calls []*InternalTxTrace `json:"calls,omitempty"`

	Reverted *RevertedInfo `json:"reverted,omitempty"`
}

type RevertedInfo struct {
	Contract *common.Address `json:"contract,omitempty"`
	Message  string          `json:"message,omitempty"`
}

// CaptureStart implements the Tracer interface to initialize the tracing operation.
func (t *InternalTxTracer) CaptureStart(env *EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	t.ctx["type"] = CALL.String()
	if create {
		t.ctx["type"] = CREATE.String()
	}
	t.ctx["from"] = from
	t.ctx["to"] = to
	t.ctx["input"] = hexutil.Encode(input)
	t.ctx["gas"] = gas
	t.ctx["gasPrice"] = env.TxContext.GasPrice
	t.ctx["value"] = value
}

// tracerLog is used to help comparing codes between this and call_tracer.js
// by following the conventions used in call_tracer.js
type tracerLog struct {
	env      *EVM
	pc       uint64
	op       OpCode
	gas      uint64
	cost     uint64
	memory   *Memory
	stack    *Stack
	contract *Contract
	depth    int
	err      error
}

func wrapError(context string, err error) error {
	return fmt.Errorf("%v    in server-side tracer function '%v'", err.Error(), context)
}

// Stop terminates execution of the tracer at the first opportune moment.
func (t *InternalTxTracer) Stop(err error) {
	t.reason = err
	atomic.StoreUint32(&t.interrupt, 1)
}

// CaptureState implements the Tracer interface to trace a single step of VM execution.
func (t *InternalTxTracer) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, scope *ScopeContext, depth int, err error) {
	if t.err == nil {
		// Initialize the context if it wasn't done yet
		if !t.initialized {
			t.ctx["block"] = env.Context.BlockNumber.Uint64()
			t.initialized = true
		}
		// If tracing was interrupted, set the error and stop
		if atomic.LoadUint32(&t.interrupt) > 0 {
			t.err = t.reason
			return
		}

		log := &tracerLog{
			env, pc, op, gas, cost,
			scope.Memory, scope.Stack, scope.Contract, depth, err,
		}
		err := t.step(log)
		if err != nil {
			t.err = wrapError("step", err)
		}
	}
}

func (t *InternalTxTracer) step(log *tracerLog) error {
	// Capture any errors immediately
	if log.err != nil {
		t.fault(log)
		return nil
	}

	// We only care about system opcodes, faster if we pre-check once
	sysCall := (log.op & 0xf0) == 0xf0
	op := log.op
	// If a new contract is being created, add to the call stack
	if sysCall && (op == CREATE || op == CREATE2) {
		inOff := log.stack.Back(1)
		inEnd := big.NewInt(0).Add(inOff.ToBig(), log.stack.Back(2).ToBig()).Int64()

		// Assemble the internal call report and store for completion
		fromAddr := log.contract.Address()
		call := &InternalCall{
			Type:    op.String(),
			From:    &fromAddr,
			Input:   hexutil.Encode(log.memory.Slice(int64(inOff.Uint64()), inEnd)),
			GasIn:   log.gas,
			GasCost: log.cost,
			Value:   "0x" + log.stack.peek().Hex(), // '0x' + tracerLog.stack.peek(0).toString(16)
		}
		t.callStack = append(t.callStack, call)
		t.descended = true
		return nil
	}
	// If a contract is being self destructed, gather that as a subcall too
	if sysCall && op == SELFDESTRUCT {
		left := t.callStackLength()
		if t.callStack[left-1] == nil {
			t.callStack[left-1] = &InternalCall{}
		}
		if t.callStack[left-1].Calls == nil {
			t.callStack[left-1].Calls = []*InternalCall{}
		}
		contractAddr := log.contract.Address()
		ret := log.stack.peek()
		toAddr := common.HexToAddress(ret.Hex())
		t.callStack[left-1].Calls = append(
			t.callStack[left-1].Calls,
			&InternalCall{
				Type:    op.String(),
				From:    &contractAddr,
				To:      &toAddr,
				Value:   "0x" + log.env.StateDB.GetBalance(contractAddr).Text(16),
				GasIn:   log.gas,
				GasCost: log.cost,
			},
		)
		return nil
	}
	// If a new method invocation is being done, add to the call stack
	if sysCall && (op == CALL || op == CALLCODE || op == DELEGATECALL || op == STATICCALL) {

		// Skip any pre-compile invocations, those are just fancy opcodes
		toAddr := common.HexToAddress(log.stack.Back(1).Hex())
		if _, ok := PrecompiledContractsByzantiumCompatible[toAddr]; ok {
			return nil
		}

		off := 1
		if op == DELEGATECALL || op == STATICCALL {
			off = 0
		}

		inOff := log.stack.Back(2 + off)
		inEnd := big.NewInt(0).Add(inOff.ToBig(), log.stack.Back(3+off).ToBig()).Int64()

		// Assemble the internal call report and store for completion
		fromAddr := log.contract.Address()
		call := &InternalCall{
			Type:    op.String(),
			From:    &fromAddr,
			To:      &toAddr,
			Input:   hexutil.Encode(log.memory.Slice(int64(inOff.Uint64()), inEnd)),
			GasIn:   log.gas,
			GasCost: log.cost,
			OutOff:  big.NewInt(int64(log.stack.Back(4 + off).Uint64())),
			OutLen:  big.NewInt(int64(log.stack.Back(5 + off).Uint64())),
		}
		if op != DELEGATECALL && op != STATICCALL {
			call.Value = "0x" + log.stack.Back(2).Hex()
		}
		t.callStack = append(t.callStack, call)
		t.descended = true

		return nil
	}
	// If we've just descended into an inner call, retrieve it's true allowance. We
	// need to extract if from within the call as there may be funky gas dynamics
	// with regard to requested and actually given gas (2300 stipend, 63/64 rule).
	if t.descended {
		if log.depth >= t.callStackLength() {
			t.callStack[t.callStackLength()-1].Gas = log.gas
		} else {
			// TODO(karalabe): The call was made to a plain account. We currently don't
			// have access to the true gas amount inside the call and so any amount will
			// mostly be wrong since it depends on a lot of input args. Skip gas for now.
		}
		t.descended = false
	}
	// If an existing call is returning, pop off the call stack
	if sysCall && op == REVERT && t.callStackLength() > 0 {
		t.callStack[t.callStackLength()-1].Error = errExecutionReverted
		if t.revertedContract == emptyAddr {
			if t.callStack[t.callStackLength()-1].To == nil {
				t.revertedContract = log.contract.Address()
			} else {
				t.revertedContract = *t.callStack[t.callStackLength()-1].To
			}
		}
		return nil
	}
	if log.depth == t.callStackLength()-1 {
		// Pop off the last call and get the execution results
		call := t.callStackPop()

		if call.Type == CREATE.String() || call.Type == CREATE2.String() {
			// If the call was a CREATE, retrieve the contract address and output code
			call.GasUsed = call.GasIn - call.GasCost - log.gas
			call.GasIn, call.GasCost = uint64(0), uint64(0)

			ret := log.stack.peek()
			if ret.Cmp(uint256.NewInt(0)) != 0 {
				toAddr := common.HexToAddress(ret.Hex())
				call.To = &toAddr
				call.Output = hexutil.Encode(log.env.StateDB.GetCode(common.HexToAddress(ret.Hex())))
			} else if call.Error == nil {
				call.Error = errInternalFailure // TODO(karalabe): surface these faults somehow
			}
		} else {
			// If the call was a contract call, retrieve the gas usage and output
			if call.Gas != uint64(0) {
				call.GasUsed = call.GasIn - call.GasCost + call.Gas - log.gas
			}
			ret := log.stack.peek()
			if ret == nil || ret.Cmp(uint256.NewInt(0)) != 0 {
				callOutOff, callOutLen := call.OutOff.Int64(), call.OutLen.Int64()
				call.Output = hexutil.Encode(log.memory.Slice(callOutOff, callOutOff+callOutLen))
			} else if call.Error == nil {
				call.Error = errInternalFailure // TODO(karalabe): surface these faults somehow
			}
			call.GasIn, call.GasCost = uint64(0), uint64(0)
			call.OutOff, call.OutLen = nil, nil
		}
		if call.Gas != uint64(0) {
			// TODO-ChainDataFetcher
			// Below is the original code, but it is just to convert the value into the hex string, nothing to do
			// call.gas = '0x' + bigInt(call.gas).toString(16);
		}
		// Inject the call into the previous one
		left := t.callStackLength()
		if left == 0 {
			left = 1 // added to avoid index out of range in golang
			t.callStack = []*InternalCall{{}}
		}
		if t.callStack[left-1] == nil {
			t.callStack[left-1] = &InternalCall{}
		}
		if len(t.callStack[left-1].Calls) == 0 {
			t.callStack[left-1].Calls = []*InternalCall{}
		}
		t.callStack[left-1].Calls = append(t.callStack[left-1].Calls, call)
	}

	return nil
}

// CaptureFault implements the Tracer interface to trace an execution fault
// while running an opcode.
func (t *InternalTxTracer) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, scope *ScopeContext, depth int, err error) {
	if t.err == nil {
		// Apart from the error, everything matches the previous invocation
		t.errValue = err.Error()

		log := &tracerLog{
			env, pc, op, gas, cost,
			scope.Memory, scope.Stack, scope.Contract, depth, err,
		}
		// fault does not return an error
		t.fault(log)
	}
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (t *InternalTxTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
	t.ctx["output"] = hexutil.Encode(output)
	t.ctx["gasUsed"] = gasUsed

	if err != nil {
		t.ctx["error"] = err
	}
}

func (t *InternalTxTracer) CaptureEnter(typ OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

func (t *InternalTxTracer) CaptureExit(output []byte, gasUsed uint64, err error) {}

func (t *InternalTxTracer) CaptureTxStart(gasLimit uint64) {
	t.gasLimit = gasLimit
}

func (t *InternalTxTracer) CaptureTxEnd(restGas uint64) {
	t.ctx["gasUsed"] = t.gasLimit - restGas
}

func (t *InternalTxTracer) GetResult() (*InternalTxTrace, error) {
	result := t.result()
	// Clean up the JavaScript environment
	t.reset()
	return result, t.err
}

// reset clears data collected during the previous tracing.
// It should act like calling jst.vm.DestroyHeap() and jst.vm.Destroy() at tracers.Tracer
func (t *InternalTxTracer) reset() {
	t.callStack = []*InternalCall{}
	t.output = nil

	t.descended = false
	t.revertedContract = common.Address{}
	t.initialized = false
	t.revertString = ""
}

// result is invoked when all the opcodes have been iterated over and returns
// the final result of the tracing.
func (t *InternalTxTracer) result() *InternalTxTrace {
	if _, exist := t.ctx["type"]; !exist {
		t.ctx["type"] = ""
	}
	if _, exist := t.ctx["from"]; !exist {
		t.ctx["from"] = nil
	}
	if _, exist := t.ctx["to"]; !exist {
		t.ctx["to"] = nil
	}
	if _, exist := t.ctx["value"]; !exist {
		t.ctx["value"] = big.NewInt(0)
	}
	if _, exist := t.ctx["gas"]; !exist {
		t.ctx["gas"] = uint64(0)
	}
	if _, exist := t.ctx["gasUsed"]; !exist {
		t.ctx["gasUsed"] = uint64(0)
	}
	if _, exist := t.ctx["input"]; !exist {
		t.ctx["input"] = ""
	}
	if _, exist := t.ctx["output"]; !exist {
		t.ctx["output"] = ""
	}
	if _, exist := t.ctx["time"]; !exist {
		t.ctx["time"] = time.Duration(0)
	}
	if t.callStackLength() == 0 {
		t.callStack = []*InternalCall{{}}
	}
	var from, to *common.Address
	if addr, ok := t.ctx["from"].(common.Address); ok {
		from = &addr
	}
	if addr, ok := t.ctx["to"].(common.Address); ok {
		to = &addr
	}

	result := &InternalTxTrace{
		Type:    t.ctx["type"].(string),
		From:    from,
		To:      to,
		Value:   "0x" + t.ctx["value"].(*big.Int).Text(16),
		Gas:     t.ctx["gas"].(uint64),
		GasUsed: t.ctx["gasUsed"].(uint64),
		Input:   t.ctx["input"].(string),
		Output:  t.ctx["output"].(string),
		Time:    t.ctx["time"].(time.Duration),
	}

	nestedCalls := []*InternalTxTrace{}
	for _, call := range t.callStack[0].Calls {
		nestedCalls = append(nestedCalls, call.ToTrace())
	}
	result.Calls = nestedCalls

	if t.callStack[0].Error != nil {
		result.Error = t.callStack[0].Error
	} else if ctxErr := t.ctx["error"]; ctxErr != nil {
		result.Error = ctxErr.(error)
	}
	if result.Error != nil && (result.Error.Error() != errExecutionReverted.Error() || result.Output == "0x") {
		result.Output = "" // delete result.output;
	}
	if err := t.ctx["error"]; err != nil && err.(error).Error() == ErrExecutionReverted.Error() {
		outputHex := t.ctx["output"].(string) // it is already a hex string

		if s, err := abi.UnpackRevert(common.FromHex(outputHex)); err == nil {
			t.revertString = s
		} else {
			t.revertString = ""
		}

		contract := t.revertedContract
		message := t.revertString
		result.Reverted = &RevertedInfo{Contract: &contract, Message: message}
	}
	return result
}

// InternalTxLogs returns the captured tracerLog entries.
func (t *InternalTxTracer) InternalTxLogs() []*InternalCall { return t.callStack }

// fault is invoked when the actual execution of an opcode fails.
func (t *InternalTxTracer) fault(log *tracerLog) {
	if t.callStackLength() == 0 {
		return
	}
	// If the topmost call already reverted, don't handle the additional fault again
	if t.callStack[t.callStackLength()-1].Error != nil {
		return
	}
	// Pop off the just failed call
	call := t.callStackPop()
	call.Error = log.err

	// Consume all available gas and clean any leftovers
	if call.Gas != uint64(0) {
		call.GasUsed = call.Gas
	}
	call.GasIn, call.GasCost = uint64(0), uint64(0)
	call.OutOff, call.OutLen = nil, nil

	// Flatten the failed call into its parent
	left := t.callStackLength()
	if left > 0 {
		if t.callStack[left-1] == nil {
			t.callStack[left-1] = &InternalCall{}
		}
		if len(t.callStack[left-1].Calls) == 0 {
			t.callStack[left-1].Calls = []*InternalCall{}
		}
		t.callStack[left-1].Calls = append(t.callStack[left-1].Calls, call)
		return
	}
	// Last call failed too, leave it in the stack
	t.callStack = append(t.callStack, call)
}

func (t *InternalTxTracer) callStackLength() int {
	return len(t.callStack)
}

func (t *InternalTxTracer) callStackPop() *InternalCall {
	if t.callStackLength() == 0 {
		return &InternalCall{}
	}

	topItem := t.callStack[t.callStackLength()-1]
	t.callStack = t.callStack[:t.callStackLength()-1]
	return topItem
}
