// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/vm/interpreter.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"fmt"
	"hash"
	"sync/atomic"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
)

// Config are the configuration options for the Interpreter
type Config struct {
	Debug                   bool   // Enables debugging
	Tracer                  Tracer // Opcode logger
	NoRecursion             bool   // Disables call, callcode, delegate call and create
	EnablePreimageRecording bool   // Enables recording of SHA3/keccak preimages

	JumpTable [256]*operation // EVM instruction table, automatically populated if unset

	// RunningEVM is to indicate the running EVM and used to stop the EVM.
	RunningEVM chan *EVM

	// UseOpcodeComputationCost is to enable applying the opcode computation cost limit.
	UseOpcodeComputationCost bool

	// Enables collecting internal transaction data during processing a block
	EnableInternalTxTracing bool

	// Prefetching is true if the EVM is used for prefetching.
	Prefetching bool

	// Additional EIPs that are to be enabled
	ExtraEips []int
}

// keccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// Interpreter is used to run Klaytn based contracts and will utilise the
// passed environment to query external sources for state information.
// The Interpreter will run the byte code VM based on the passed
// configuration.
type Interpreter struct {
	evm *EVM
	cfg *Config

	intPool *intPool

	hasher    keccakState // Keccak256 hasher instance shared across opcodes
	hasherBuf common.Hash // Keccak256 hasher result array shared aross opcodes

	readOnly   bool   // Whether to throw on stateful modifications
	returnData []byte // Last CALL's return data for subsequent reuse
}

// NewEVMInterpreter returns a new instance of the Interpreter.
func NewEVMInterpreter(evm *EVM, cfg *Config) *Interpreter {
	// We use the STOP instruction whether to see
	// the jump table was initialised. If it was not
	// we'll set the default jump table.
	if cfg.JumpTable[STOP] == nil {
		var jt JumpTable
		switch {
		case evm.chainRules.IsLondon:
			jt = LondonInstructionSet
		case evm.chainRules.IsIstanbul:
			jt = IstanbulInstructionSet
		default:
			jt = ConstantinopleInstructionSet
		}
		for i, eip := range cfg.ExtraEips {
			if err := EnableEIP(eip, &jt); err != nil {
				// Disable it, so caller can check if it's activated or not
				cfg.ExtraEips = append(cfg.ExtraEips[:i], cfg.ExtraEips[i+1:]...)
				logger.Error("EIP activation failed", "eip", eip, "error", err)
			}
		}
		cfg.JumpTable = jt
	}

	return &Interpreter{
		evm: evm,
		cfg: cfg,
	}
}

///////////////////////////////////////////////////////
// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
//var (
//	prevOp OpCode
//	globalTimer = time.Now()
//	opCnt = make([]uint64, 256)
//	opTime = make([]uint64, 256)
//	precompiledCnt = make([]uint64, 16)
//	precompiledTime = make([]uint64, 16)
//	opDebug = true
//)
///////////////////////////////////////////////////////

// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operation except for
// ErrExecutionReverted which means revert-and-keep-gas-left.
func (in *Interpreter) Run(contract *Contract, input []byte) (ret []byte, err error) {
	if in.intPool == nil {
		in.intPool = poolOfIntPools.get()
		defer func() {
			poolOfIntPools.put(in.intPool)
			in.intPool = nil
		}()
	}

	///////////////////////////////////////////////////////
	// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
	//if opDebug {
	//	if in.evm.depth == 0 {
	//		for i := 0; i< 256; i++ {
	//			opCnt[i] = 0
	//			opTime[i] = 0
	//		}
	//		prevOp = 0
	//		defer func() {
	//			for i := 0; i < 256; i++ {
	//				if opCnt[i] > 0 {
	//					fmt.Println("op", OpCode(i).String(), "computationCost", in.cfg.JumpTable[i].computationCost, "cnt", opCnt[i], "avg", opTime[i]/opCnt[i])
	//				}
	//			}
	//			for i := 0; i < 16; i++ {
	//				if precompiledCnt[i] > 0 {
	//					fmt.Println("precompiled contract addr", i, "cnt", precompiledCnt[i], "avg", precompiledTime[i]/precompiledCnt[i])
	//				}
	//			}
	//		}()
	//	}
	//}
	///////////////////////////////////////////////////////

	// Increment the call depth which is restricted to 1024
	in.evm.depth++
	defer func() { in.evm.depth-- }()

	// Reset the previous call's return data. It's unimportant to preserve the old buffer
	// as every returning call will return new data anyway.
	in.returnData = nil

	// Don't bother with the execution if there's no code.
	if len(contract.Code) == 0 {
		return nil, nil
	}

	var (
		op    OpCode        // current opcode
		mem   = NewMemory() // bound memory
		stack = newstack()  // local stack
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc   = uint64(0) // program counter
		cost uint64
		// copies used by tracer
		pcCopy              uint64              // needed for the deferred Tracer
		gasCopy             uint64              // for Tracer to log gas remaining before execution
		logged              bool                // deferred Tracer should ignore already logged steps
		res                 []byte              // result of the opcode execution function
		allocatedMemorySize = uint64(mem.Len()) // Currently allocated memory size
	)
	contract.Input = input

	// Reclaim the stack as an int pool when the execution stops
	defer func() { in.intPool.put(stack.data...) }()

	if in.cfg.Debug {
		defer func() {
			if err != nil {
				if !logged {
					in.cfg.Tracer.CaptureState(in.evm, pcCopy, op, gasCopy, cost, mem, stack, contract, in.evm.depth, err)
				} else {
					in.cfg.Tracer.CaptureFault(in.evm, pcCopy, op, gasCopy, cost, mem, stack, contract, in.evm.depth, err)
				}
			}
		}()
	}

	// The Interpreter main run loop (contextual). This loop runs until either an
	// explicit STOP, RETURN or SELFDESTRUCT is executed, an error occurred during
	// the execution of one of the operations or until the done flag is set by the
	// parent context.
	for atomic.LoadInt32(&in.evm.abort) == 0 {
		if in.cfg.Debug {
			// Capture pre-execution values for tracing.
			logged, pcCopy, gasCopy = false, pc, contract.Gas
		}

		///////////////////////////////////////////////////////
		// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
		//if opDebug {
		//	prevOp = op
		//}
		///////////////////////////////////////////////////////
		// Get the operation from the jump table and validate the stack to ensure there are
		// enough stack items available to perform the operation.
		op = contract.GetOp(pc)
		operation := in.cfg.JumpTable[op]
		if operation == nil {
			return nil, fmt.Errorf("invalid opcode 0x%x", int(op)) // TODO-Klaytn-Issue615
		}
		// Validate stack
		if sLen := stack.len(); sLen < operation.minStack {
			return nil, fmt.Errorf("stack underflow (%d <=> %d)", sLen, operation.minStack)
		} else if sLen > operation.maxStack {
			return nil, fmt.Errorf("stack limit reached %d (%d)", sLen, operation.maxStack)
		}
		// If the operation is valid, enforce and write restrictions
		if in.readOnly {
			// If the interpreter is operating in readonly mode, make sure no
			// state-modifying operation is performed. The 3rd stack item
			// for a call operation is the value. Transferring value from one
			// account to the others means the state is modified and should also
			// return with an error.
			if operation.writes || (op == CALL && stack.Back(2).Sign() != 0) {
				return nil, ErrWriteProtection
			}
		}

		// Static portion of gas
		cost = operation.constantGas // For tracing
		if !contract.UseGas(operation.constantGas) {
			return nil, kerrors.ErrOutOfGas
		}

		// We limit tx's execution time using the sum of computation cost of opcodes.
		if in.evm.vmConfig.UseOpcodeComputationCost {
			///////////////////////////////////////////////////////
			// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
			//if opDebug && prevOp > 0 {
			//	elapsed := uint64(time.Since(globalTimer).Nanoseconds())
			//	fmt.Println("[", in.evm.depth, "]", "prevop", prevOp.String(), "-", op.String(),  "computationCost", in.cfg.JumpTable[prevOp].computationCost, "total", in.evm.opcodeComputationCostSum, "elapsed", elapsed)
			//	opTime[prevOp] += elapsed
			//	opCnt[prevOp] += 1
			//}
			//globalTimer = time.Now()
			///////////////////////////////////////////////////////
			in.evm.opcodeComputationCostSum += operation.computationCost
			if in.evm.opcodeComputationCostSum > params.OpcodeComputationCostLimit {
				return nil, ErrOpcodeComputationCostLimitReached
			}
		}
		var memorySize uint64
		var extraSize uint64
		// calculate the new memory size and expand the memory to fit
		// the operation
		// Memory check needs to be done prior to evaluating the dynamic gas portion,
		// to detect calculation overflows
		if operation.memorySize != nil {
			memSize, overflow := operation.memorySize(stack)
			if overflow {
				return nil, errGasUintOverflow // TODO-Klaytn-Issue615
			}
			// memory is expanded in words of 32 bytes. Gas
			// is also calculated in words.
			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				return nil, errGasUintOverflow // TODO-Klaytn-Issue615
			}
			if allocatedMemorySize < memorySize {
				extraSize = memorySize - allocatedMemorySize
			}
		}
		// Dynamic portion of gas
		// consume the gas and return an error if not enough gas is available.
		// cost is explicitly set so that the capture state defer method can get the proper cost
		if operation.dynamicGas != nil {
			var dynamicCost uint64
			dynamicCost, err = operation.dynamicGas(in.evm, contract, stack, mem, memorySize)
			cost += dynamicCost // total cost, for debug tracing
			if err != nil || !contract.UseGas(dynamicCost) {
				return nil, kerrors.ErrOutOfGas // TODO-Klaytn-Issue615
			}
		}
		if extraSize > 0 {
			mem.Increase(extraSize)
			allocatedMemorySize = uint64(mem.Len())
		}

		if in.cfg.Debug {
			in.cfg.Tracer.CaptureState(in.evm, pc, op, gasCopy, cost, mem, stack, contract, in.evm.depth, err)
			logged = true
		}

		// execute the operation
		res, err = operation.execute(&pc, in.evm, contract, mem, stack)
		// verifyPool is a build flag. Pool verification makes sure the integrity
		// of the integer pool by comparing values to a default value.
		if verifyPool {
			verifyIntegerPool(in.intPool)
		}
		// if the operation clears the return data (e.g. it has returning data)
		// set the last return to the result of the operation.
		if operation.returns {
			in.returnData = res
		}

		switch {
		case err != nil:
			return nil, err // TODO-Klaytn-Issue615
		case operation.reverts:
			return res, ErrExecutionReverted // TODO-Klaytn-Issue615
		case operation.halts:
			return res, nil
		case !operation.jumps:
			pc++
		}
	}

	abort := atomic.LoadInt32(&in.evm.abort)
	if (abort & CancelByTotalTimeLimit) != 0 {
		return nil, ErrTotalTimeLimitReached // TODO-Klaytn-Issue615
	}
	return nil, nil
}
