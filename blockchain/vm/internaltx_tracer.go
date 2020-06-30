package vm

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"math/big"
	"strconv"
	"time"
)

var errEvmExecutionReverted = errors.New("evm: execution reverted")
var errExecutionReverted = errors.New("execution reverted")
var errInternalFailure = errors.New("internal failure")
var emptyAddr = common.Address{}

// InternalTxTracer is a full blown transaction tracer that extracts and reports all
// the internal calls made by a transaction, along with any useful information.
// It is ported to golang from JS, specifically call_tracer.js
type InternalTxTracer struct {
	cfg LogConfig

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
}

// NewInternalTxLogger returns a new InternalTxTracer
func NewInternalTxLogger(cfg *LogConfig) *InternalTxTracer {
	logger := &InternalTxTracer{}
	if cfg != nil {
		logger.cfg = *cfg
	}
	return logger
}

// InternalCall is emitted to the EVM each cycle and lists information about the current internal state
// prior to the execution of the statement.
type InternalCall struct {
	Type          OpCode         `json:"op"`
	From          common.Address `json:"from"`
	To            common.Address `json:"to"`
	Input         string
	Gas           uint64                      `json:"gas"`
	GasIn         uint64                      `json:"gasIn"`
	GasUsed       uint64                      `json:"gasUsed"`
	GasCost       uint64                      `json:"gasCost"`
	Value         string                      `json:"value"`
	Err           error                       `json:"-"`
	OutOff        *big.Int
	OutLen        *big.Int
	Output        []byte

	calls []*InternalCall
}

// OpName formats the operand name in a human-readable format.
func (s *InternalCall) OpName() string {
	return s.Type.String()
}

// ErrorString formats the tracerLog's error as a string.
func (s *InternalCall) ErrorString() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}

// InternalTxTraceResult is returned data after the end of trace-collecting cycle.
// It implements an object returned by "result" function at call_tracer.js
type InternalTxTraceResult struct {
	callType OpCode
	from     common.Address
	to       common.Address
	value    string
	gas      uint64
	gasUsed  uint64
	input    string // hex string
	output   string // hex string
	time     time.Duration
	calls    []*InternalCall
	error    error
	reverted revertedInfo
}

type revertedInfo struct {
	contract common.Address
	message  string
}

// CaptureStart implements the Tracer interface to initialize the tracing operation.
func (this *InternalTxTracer) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) error {
	this.ctx["type"] = "CALL"
	if create {
		this.ctx["type"] = "CREATE"
	}
	this.ctx["from"] = from
	this.ctx["to"] = to
	this.ctx["input"] = hexutil.Encode(input)
	this.ctx["gas"] = gas
	this.ctx["value"] = value

	return nil
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

// CaptureState implements the Tracer interface to trace a single step of VM execution.
func (this *InternalTxTracer) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, logStack *Stack, contract *Contract, depth int, err error) error {
	if this.err == nil {
		// Initialize the context if it wasn't done yet
		if !this.initialized {
			this.ctx["block"] = env.BlockNumber.Uint64()
			this.initialized = true
		}
		// As this tracer is executed inside the node and there is no reason to stop this
		// while the node is running, let's skip stop logic from tracers.Tracer

		log := &tracerLog{
			env, pc, op, gas, cost,
			memory, logStack, contract, depth, err,
		}
		err := this.step(log)
		if err != nil {
			this.err = wrapError("step", err)
		}
	}
	return nil
}

func (this *InternalTxTracer) step(log *tracerLog) error {
	// Capture any errors immediately
	if log.err != nil {
		this.fault(log)
		return nil
	}

	// We only care about system opcodes, faster if we pre-check once
	sysCall := (log.op & 0xf0) == 0xf0
	op := log.op
	// If a new contract is being created, add to the call stack
	if sysCall && (op == CREATE || op == CREATE2) {
		inOff := log.stack.Back(1)
		inEnd := big.NewInt(0).Add(inOff, log.stack.Back(2)).Int64()

		// Assemble the internal call report and store for completion
		call := &InternalCall{
			Type:    op,
			From:    log.contract.Address(),
			Input:   hexutil.Encode(log.memory.store[inOff.Int64():inEnd]),
			Gas:     log.gas,
			GasCost: log.cost,
			Value:   "0x" + log.stack.peek().Text(16), // '0x' + tracerLog.stack.peek(0).toString(16)
		}
		this.callStack = append(this.callStack, call)
		this.descended = true
		return nil
	}
	// If a contract is being self destructed, gather that as a subcall too
	if sysCall && op == SELFDESTRUCT {
		left := this.callStackLength()
		if this.callStack[left-1] == nil {
			this.callStack[left-1] = &InternalCall{}
		}
		if this.callStack[left-1].calls == nil {
			this.callStack[left-1].calls = []*InternalCall{}
		}
		this.callStack = append(this.callStack, &InternalCall{Type: op})
		return nil
	}
	// If a new method invocation is being done, add to the call stack
	if sysCall && (op == CALL || op == CALLCODE || op == DELEGATECALL || op == STATICCALL) {
		// Skip any pre-compile invocations, those are just fancy opcodes
		toAddr := common.HexToAddress(log.stack.Back(1).String())
		if common.IsPrecompiledContractAddress(toAddr) {
			return nil
		}
		off := 1
		if op == DELEGATECALL || op == STATICCALL {
			off = 0
		}

		inOff := log.stack.Back(2 + off)
		inEnd := big.NewInt(0).Add(inOff, log.stack.Back(3+off)).Int64()

		// Assemble the internal call report and store for completion
		call := &InternalCall{
			Type:    op,
			From:    log.contract.Address(),
			To:      toAddr,
			Input:   hexutil.Encode(log.memory.store[inOff.Int64():inEnd]),
			GasIn:   log.gas,
			GasCost: log.cost,
			OutOff:  big.NewInt(log.stack.Back(4 + off).Int64()),
			OutLen:  big.NewInt(log.stack.Back(5 + off).Int64()),
		}
		if op != DELEGATECALL && op != STATICCALL {
			call.Value = "0x" + log.stack.Back(2).Text(16)
		}
		this.callStack = append(this.callStack, call)
		this.descended = true
		return nil
	}
	// If we've just descended into an inner call, retrieve it's true allowance. We
	// need to extract if from within the call as there may be funky gas dynamics
	// with regard to requested and actually given gas (2300 stipend, 63/64 rule).
	if this.descended {
		if log.depth >= this.callStackLength() {
			this.callStack[this.callStackLength()-1].Gas = log.gas
		} else {
			// TODO(karalabe): The call was made to a plain account. We currently don't
			// have access to the true gas amount inside the call and so any amount will
			// mostly be wrong since it depends on a lot of input args. Skip gas for now.
		}
		this.descended = false
	}
	// If an existing call is returning, pop off the call stack
	if sysCall && op == REVERT {
		this.callStack[this.callStackLength()-1].Err = errExecutionReverted
		if this.revertedContract == emptyAddr {
			if this.callStack[this.callStackLength()-1].To == emptyAddr {
				this.revertedContract = log.contract.Address()
			} else {
				this.revertedContract = this.callStack[this.callStackLength()-1].To
			}
		}
		return nil
	}
	if log.depth == this.callStackLength()-1 {
		// Pop off the last call and get the execution results
		call := this.callStackPop()

		if call.Type == CREATE || call.Type == CREATE2 {
			// If the call was a CREATE, retrieve the contract address and output code
			call.GasUsed = call.GasIn - call.GasCost - log.gas
			call.GasIn, call.GasCost = uint64(0), uint64(0)

			ret := log.stack.peek()
			if ret.Cmp(big.NewInt(0)) != 0 {
				call.To = common.HexToAddress(ret.String())
				call.Output = log.env.StateDB.GetCode(common.HexToAddress(ret.String()))
			} else if call.Err == nil {
				call.Err = errInternalFailure // TODO(karalabe): surface these faults somehow
			}
		} else {
			// If the call was a contract call, retrieve the gas usage and output
			if call.Gas != uint64(0) {
				call.GasUsed = call.GasIn - call.GasCost + call.Gas - log.gas

				ret := log.stack.peek()
				if ret.Cmp(big.NewInt(0)) != 0 {
					callOutOff, callOutLen := call.OutOff.Int64(), call.OutLen.Int64()
					call.Output = log.memory.store[callOutOff : callOutOff+callOutLen]
				} else if call.Err == nil {
					call.Err = errInternalFailure // TODO(karalabe): surface these faults somehow
				}
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
		left := this.callStackLength()
		if this.callStack[left-1] == nil {
			this.callStack[left-1] = &InternalCall{}
		}
		if len(this.callStack[left-1].calls) == 0 {
			this.callStack[left-1].calls = []*InternalCall{}
		}
		this.callStack[left-1].calls = append(this.callStack[left-1].calls, call)
	}

	return nil
}

// CaptureFault implements the Tracer interface to trace an execution fault
// while running an opcode.
func (this *InternalTxTracer) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, s *Stack, contract *Contract, depth int, err error) error {
	if this.err == nil {
		// Apart from the error, everything matches the previous invocation
		this.errValue = err.Error()

		log := &tracerLog{
			env, pc, op, gas, cost,
			memory, s, contract, depth, err,
		}
		// fault does not return an error
		this.fault(log)
	}
	return nil
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (this *InternalTxTracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	this.ctx["output"] = hexutil.Encode(output)
	this.ctx["gasUsed"] = gasUsed
	this.ctx["time"] = t

	if err != nil {
		this.ctx["error"] = err.Error()
	}
	return nil
}

func (this *InternalTxTracer) GetResult() (*InternalTxTraceResult, error) {
	result, err := this.result()
	if err != nil {
		this.err = wrapError("result", err)
	}
	// Clean up the JavaScript environment
	this.reset()
	return result, this.err
}

// reset clears data collected during the previous tracing.
// It should act like calling jst.vm.DestroyHeap() and jst.vm.Destroy() at tracers.Tracer
func (this *InternalTxTracer) reset() {
	this.callStack = []*InternalCall{}
	this.output = nil

	this.descended = false
	this.revertedContract = common.Address{}
	this.initialized = false
	this.revertString = ""
}

// result is invoked when all the opcodes have been iterated over and returns
// the final result of the tracing.
func (this *InternalTxTracer) result() (*InternalTxTraceResult, error) {
	if _, exist := this.ctx["type"]; !exist {
		this.ctx["type"] = 0
	}
	if _, exist := this.ctx["from"]; !exist {
		this.ctx["from"] = common.Address{}
	}
	if _, exist := this.ctx["to"]; !exist {
		this.ctx["to"] = common.Address{}
	}
	if _, exist := this.ctx["value"]; !exist {
		this.ctx["value"] = big.NewInt(0)
	}
	if _, exist := this.ctx["gas"]; !exist {
		this.ctx["gas"] = 0
	}
	if _, exist := this.ctx["gasUsed"]; !exist {
		this.ctx["gasUsed"] = 0
	}
	if _, exist := this.ctx["input"]; !exist {
		this.ctx["input"] = ""
	}
	if _, exist := this.ctx["output"]; !exist {
		this.ctx["output"] = ""
	}
	if _, exist := this.ctx["time"]; !exist {
		this.ctx["time"] = time.Duration(0)
	}
	if this.callStackLength() == 0 {
		this.callStack = []*InternalCall{{}}
	}
	result := &InternalTxTraceResult{
		callType: this.ctx["type"].(OpCode),
		from:     this.ctx["from"].(common.Address),
		to:       this.ctx["to"].(common.Address),
		value:    "0x" + this.ctx["value"].(*big.Int).Text(16),
		gas:      this.ctx["gas"].(uint64),
		gasUsed:  this.ctx["gasUsed"].(uint64),
		input:    this.ctx["input"].(string),
		output:   this.ctx["output"].(string),
		time:     this.ctx["time"].(time.Duration),
	}
	if this.callStack[0].calls != nil {
		result.calls = this.callStack[0].calls
	}
	if this.callStack[0].Err != nil {
		result.error = this.callStack[0].Err
	} else if ctxErr, _ := this.ctx["error"]; ctxErr != nil {
		result.error = ctxErr.(error)
	}
	if result.error != nil {
		result.output = "" // delete result.output;
	}
	if this.ctx["error"] == errEvmExecutionReverted {
		outputHex := this.ctx["output"].(string) // it is already a hex string
		if outputHex[2:10] == "08c379a0" {
			defaultOffset := 10

			stringOffset, err := strconv.ParseInt(outputHex[defaultOffset:defaultOffset+32*2], 16, 64)
			if err != nil {
				logger.Error("failed to parse hex string to get stringOffset",
					"err", err, "outputHex", outputHex)
				return nil, err
			}
			stringLength, err := strconv.ParseInt(outputHex[defaultOffset+32*2:defaultOffset+32*2+32*2], 16, 64)
			if err != nil {
				logger.Error("failed to parse hex string to get stringLength",
					"err", err, "outputHex", outputHex)
				return nil, err
			}
			start := defaultOffset + 32*2 + int(stringOffset*2)
			end := start + int(stringLength*2)
			asciiInBytes, err := hex.DecodeString(outputHex[start:end])
			if err != nil {
				logger.Error("failed to parse hex string to get ASCII representation",
					"err", err, "outputHex", outputHex)
				return nil, err
			}
			this.revertString = string(asciiInBytes)
		}
		result.reverted = revertedInfo{contract: this.revertedContract, message: this.revertString}
	}
	return result, nil
}

// InternalTxLogs returns the captured tracerLog entries.
func (this *InternalTxTracer) InternalTxLogs() []*InternalCall { return this.callStack }

// fault is invoked when the actual execution of an opcode fails.
func (this *InternalTxTracer) fault(log *tracerLog) {
	// If the topmost call already reverted, don't handle the additional fault again
	if this.callStack[this.callStackLength()-1].Err != nil {
		return
	}
	// Pop off the just failed call
	call := this.callStackPop()
	call.Err = log.err

	// Consume all available gas and clean any leftovers
	if call.Gas != uint64(0) {
		call.GasUsed = call.Gas
	}
	call.GasIn, call.GasCost = uint64(0), uint64(0)
	call.OutOff, call.OutLen = nil, nil

	// Flatten the failed call into its parent
	left := this.callStackLength()
	if left > 0 {
		if this.callStack[left-1] == nil {
			this.callStack[left-1] = &InternalCall{}
		}
		if len(this.callStack[left-1].calls) == 0 {
			this.callStack[left-1].calls = []*InternalCall{}
		}
		this.callStack[left-1].calls = append(this.callStack[left-1].calls, call)
		return
	}
	// Last call failed too, leave it in the stack
	this.callStack = append(this.callStack, call)
}

func (this *InternalTxTracer) callStackLength() int {
	return len(this.callStack)
}

func (this *InternalTxTracer) callStackPop() *InternalCall {
	topItem := this.callStack[this.callStackLength()-1]
	this.callStack = this.callStack[:this.callStackLength()-2]
	return topItem
}

func wrapError(context string, err error) error {
	var message string
	switch err := err.(type) {
	default:
		message = err.Error()
	}
	return fmt.Errorf("%v    in server-side tracer function '%v'", message, context)
}
