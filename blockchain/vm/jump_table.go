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
// This file is derived from core/vm/jump_table.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"errors"

	"github.com/klaytn/klaytn/params"
)

type (
	executionFunc func(pc *uint64, env *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error)
	gasFunc       func(*EVM, *Contract, *Stack, *Memory, uint64) (uint64, error) // last parameter is the requested memory size as a uint64
	// memorySizeFunc returns the required size, and whether the operation overflowed a uint64
	memorySizeFunc func(*Stack) (size uint64, overflow bool)
)

var errGasUintOverflow = errors.New("gas uint64 overflow")

type operation struct {
	// execute is the operation function
	execute     executionFunc
	constantGas uint64
	dynamicGas  gasFunc
	// minStack tells how many stack items are required
	minStack int
	// maxStack specifies the max length the stack can have for this operation
	// to not overflow the stack.
	maxStack int

	// memorySize returns the memory size required for the operation
	memorySize memorySizeFunc

	// computationCost represents approximated execution time of an operation.
	// This value will be used to limit the execution time of a transaction on EVM.
	computationCost uint64

	halts   bool // indicates whether the operation should halt further execution
	jumps   bool // indicates whether the program counter should not increment
	writes  bool // determines whether this a state modifying operation
	reverts bool // determines whether the operation reverts state (implicitly halts)
	returns bool // determines whether the operations sets the return data content
}

var (
	ConstantinopleInstructionSet = newConstantinopleInstructionSet()
	IstanbulInstructionSet       = newIstanbulInstructionSet()
	LondonInstructionSet         = newLondonInstructionSet()
)

// JumpTable contains the EVM opcodes supported at a given fork.
type JumpTable [256]*operation

// newLondonInstructionSet returns the frontier, homestead, byzantium,
// constantinople, istanbul, petersburg, berlin and london instructions.
func newLondonInstructionSet() JumpTable {
	instructionSet := newIstanbulInstructionSet()
	enable3198(&instructionSet) // Base fee opcode https://eips.ethereum.org/EIPS/eip-3198
	return instructionSet
}

// newIstanbulInstructionSet returns the frontier, homestead, byzantium,
// constantinople, istanbul and petersburg instructions.
func newIstanbulInstructionSet() JumpTable {
	instructionSet := newConstantinopleInstructionSet()
	enable1344(&instructionSet)
	enable1884(&instructionSet)
	enable2200(&instructionSet)
	enableIstanbulComputationCostModification(&instructionSet)
	return instructionSet
}

// newConstantinopleInstructionSet returns the frontier, homestead
// byzantium and constantinople instructions.
func newConstantinopleInstructionSet() JumpTable {
	instructionSet := newByzantiumInstructionSet()
	instructionSet[SHL] = &operation{
		execute:         opSHL,
		constantGas:     GasFastestStep,
		minStack:        minStack(2, 1),
		maxStack:        maxStack(2, 1),
		computationCost: params.ShlComputationCost,
	}
	instructionSet[SHR] = &operation{
		execute:         opSHR,
		constantGas:     GasFastestStep,
		minStack:        minStack(2, 1),
		maxStack:        maxStack(2, 1),
		computationCost: params.ShrComputationCost,
	}
	instructionSet[SAR] = &operation{
		execute:         opSAR,
		constantGas:     GasFastestStep,
		minStack:        minStack(2, 1),
		maxStack:        maxStack(2, 1),
		computationCost: params.SarComputationCost,
	}
	instructionSet[EXTCODEHASH] = &operation{
		execute:         opExtCodeHash,
		constantGas:     params.ExtcodeHashGasConstantinople,
		minStack:        minStack(1, 1),
		maxStack:        maxStack(1, 1),
		computationCost: params.ExtCodeHashComputationCost,
	}
	instructionSet[CREATE2] = &operation{
		execute:         opCreate2,
		constantGas:     params.Create2Gas,
		dynamicGas:      gasCreate2,
		minStack:        minStack(4, 1),
		maxStack:        maxStack(4, 1),
		memorySize:      memoryCreate2,
		writes:          true,
		returns:         true,
		computationCost: params.Create2ComputationCost,
	}
	return instructionSet
}

// newByzantiumInstructionSet returns the frontier, homestead and
// byzantium instructions.
func newByzantiumInstructionSet() JumpTable {
	// instructions that can be executed during the homestead phase.
	instructionSet := newHomesteadInstructionSet()
	instructionSet[STATICCALL] = &operation{
		execute:         opStaticCall,
		constantGas:     params.CallGas,
		dynamicGas:      gasStaticCall,
		minStack:        minStack(6, 1),
		maxStack:        maxStack(6, 1),
		memorySize:      memoryStaticCall,
		returns:         true,
		computationCost: params.StaticCallComputationCost,
	}
	instructionSet[RETURNDATASIZE] = &operation{
		execute:         opReturnDataSize,
		constantGas:     GasQuickStep,
		minStack:        minStack(0, 1),
		maxStack:        maxStack(0, 1),
		computationCost: params.ReturnDataSizeComputationCost,
	}
	instructionSet[RETURNDATACOPY] = &operation{
		execute:         opReturnDataCopy,
		constantGas:     GasFastestStep,
		dynamicGas:      gasReturnDataCopy,
		minStack:        minStack(3, 0),
		maxStack:        maxStack(3, 0),
		memorySize:      memoryReturnDataCopy,
		computationCost: params.ReturnDataCopyComputationCost,
	}
	instructionSet[REVERT] = &operation{
		execute:         opRevert,
		dynamicGas:      gasRevert,
		minStack:        minStack(2, 0),
		maxStack:        maxStack(2, 0),
		memorySize:      memoryRevert,
		reverts:         true,
		returns:         true,
		computationCost: params.RevertComputationCost,
	}
	return instructionSet
}

// newHomesteadInstructionSet returns the frontier and homestead
// instructions that can be executed during the homestead phase.
func newHomesteadInstructionSet() JumpTable {
	instructionSet := newFrontierInstructionSet()
	instructionSet[DELEGATECALL] = &operation{
		execute:         opDelegateCall,
		constantGas:     params.CallGas,
		dynamicGas:      gasDelegateCall,
		minStack:        minStack(6, 1),
		maxStack:        maxStack(6, 1),
		memorySize:      memoryDelegateCall,
		returns:         true,
		computationCost: params.DelegateCallComputationCost,
	}
	return instructionSet
}

// newFrontierInstructionSet returns the frontier instructions
// that can be executed during the frontier phase.
func newFrontierInstructionSet() JumpTable {
	return JumpTable{
		STOP: {
			execute:         opStop,
			constantGas:     0,
			minStack:        minStack(0, 0),
			maxStack:        maxStack(0, 0),
			halts:           true,
			computationCost: params.StopComputationCost,
		},
		ADD: {
			execute:         opAdd,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.AddComputationCost,
		},
		MUL: {
			execute:         opMul,
			constantGas:     GasFastStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.MulComputationCost,
		},
		SUB: {
			execute:         opSub,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.SubComputationCost,
		},
		DIV: {
			execute:         opDiv,
			constantGas:     GasFastStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.DivComputationCost,
		},
		SDIV: {
			execute:         opSdiv,
			constantGas:     GasFastStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.SdivComputationCost,
		},
		MOD: {
			execute:         opMod,
			constantGas:     GasFastStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.ModComputationCost,
		},
		SMOD: {
			execute:         opSmod,
			constantGas:     GasFastStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.SmodComputationCost,
		},
		ADDMOD: {
			execute:         opAddmod,
			constantGas:     GasMidStep,
			minStack:        minStack(3, 1),
			maxStack:        maxStack(3, 1),
			computationCost: params.AddmodComputationCost,
		},
		MULMOD: {
			execute:         opMulmod,
			constantGas:     GasMidStep,
			minStack:        minStack(3, 1),
			maxStack:        maxStack(3, 1),
			computationCost: params.MulmodComputationCost,
		},
		EXP: {
			execute:         opExp,
			dynamicGas:      gasExp,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.ExpComputationCost,
		},
		SIGNEXTEND: {
			execute:         opSignExtend,
			constantGas:     GasFastStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.SignExtendComputationCost,
		},
		LT: {
			execute:         opLt,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.LtComputationCost,
		},
		GT: {
			execute:         opGt,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.GtComputationCost,
		},
		SLT: {
			execute:         opSlt,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.SltComputationCost,
		},
		SGT: {
			execute:         opSgt,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.SgtComputationCost,
		},
		EQ: {
			execute:         opEq,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.EqComputationCost,
		},
		ISZERO: {
			execute:         opIszero,
			constantGas:     GasFastestStep,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			computationCost: params.IszeroComputationCost,
		},
		AND: {
			execute:         opAnd,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.AndComputationCost,
		},
		XOR: {
			execute:         opXor,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.XorComputationCost,
		},
		OR: {
			execute:         opOr,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.OrComputationCost,
		},
		NOT: {
			execute:         opNot,
			constantGas:     GasFastestStep,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			computationCost: params.NotComputationCost,
		},
		BYTE: {
			execute:         opByte,
			constantGas:     GasFastestStep,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			computationCost: params.ByteComputationCost,
		},
		SHA3: {
			execute:         opSha3,
			constantGas:     params.Sha3Gas,
			dynamicGas:      gasSha3,
			minStack:        minStack(2, 1),
			maxStack:        maxStack(2, 1),
			memorySize:      memorySha3,
			computationCost: params.Sha3ComputationCost,
		},
		ADDRESS: {
			execute:         opAddress,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.AddressComputationCost,
		},
		BALANCE: {
			execute:         opBalance,
			constantGas:     params.BalanceGasEIP150,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			computationCost: params.BalanceComputationCost,
		},
		ORIGIN: {
			execute:         opOrigin,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.OriginComputationCost,
		},
		CALLER: {
			execute:         opCaller,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.CallerComputationCost,
		},
		CALLVALUE: {
			execute:         opCallValue,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.CallValueComputationCost,
		},
		CALLDATALOAD: {
			execute:         opCallDataLoad,
			constantGas:     GasFastestStep,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			computationCost: params.CallDataLoadComputationCost,
		},
		CALLDATASIZE: {
			execute:         opCallDataSize,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.CallDataSizeComputationCost,
		},
		CALLDATACOPY: {
			execute:         opCallDataCopy,
			constantGas:     GasFastestStep,
			dynamicGas:      gasCallDataCopy,
			minStack:        minStack(3, 0),
			maxStack:        maxStack(3, 0),
			memorySize:      memoryCallDataCopy,
			computationCost: params.CallDataCopyComputationCost,
		},
		CODESIZE: {
			execute:         opCodeSize,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.CodeSizeComputationCost,
		},
		CODECOPY: {
			execute:         opCodeCopy,
			constantGas:     GasFastestStep,
			dynamicGas:      gasCodeCopy,
			minStack:        minStack(3, 0),
			maxStack:        maxStack(3, 0),
			memorySize:      memoryCodeCopy,
			computationCost: params.CodeCopyComputationCost,
		},
		GASPRICE: {
			execute:         opGasprice,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.GasPriceComputationCost,
		},
		EXTCODESIZE: {
			execute:         opExtCodeSize,
			constantGas:     params.ExtcodeSizeGas,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			computationCost: params.ExtCodeSizeComputationCost,
		},
		EXTCODECOPY: {
			execute:         opExtCodeCopy,
			constantGas:     params.ExtcodeCopyBase,
			dynamicGas:      gasExtCodeCopy,
			minStack:        minStack(4, 0),
			maxStack:        maxStack(4, 0),
			memorySize:      memoryExtCodeCopy,
			computationCost: params.ExtCodeCopyComputationCost,
		},
		BLOCKHASH: {
			execute:         opBlockhash,
			constantGas:     GasExtStep,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			computationCost: params.BlockHashComputationCost,
		},
		COINBASE: {
			execute:         opCoinbase,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.CoinbaseComputationCost,
		},
		TIMESTAMP: {
			execute:         opTimestamp,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.TimestampComputationCost,
		},
		NUMBER: {
			execute:         opNumber,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.NumberComputationCost,
		},
		DIFFICULTY: {
			execute:         opDifficulty,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.DifficultyComputationCost,
		},
		GASLIMIT: {
			execute:         opGasLimit,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.GasLimitComputationCost,
		},
		POP: {
			execute:         opPop,
			constantGas:     GasQuickStep,
			minStack:        minStack(1, 0),
			maxStack:        maxStack(1, 0),
			computationCost: params.PopComputationCost,
		},
		MLOAD: {
			execute:         opMload,
			constantGas:     GasFastestStep,
			dynamicGas:      gasMLoad,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			memorySize:      memoryMLoad,
			computationCost: params.MloadComputationCost,
		},
		MSTORE: {
			execute:         opMstore,
			constantGas:     GasFastestStep,
			dynamicGas:      gasMStore,
			minStack:        minStack(2, 0),
			maxStack:        maxStack(2, 0),
			memorySize:      memoryMStore,
			computationCost: params.MstoreComputationCost,
		},
		MSTORE8: {
			execute:         opMstore8,
			constantGas:     GasFastestStep,
			dynamicGas:      gasMStore8,
			memorySize:      memoryMStore8,
			minStack:        minStack(2, 0),
			maxStack:        maxStack(2, 0),
			computationCost: params.Mstore8ComputationCost,
		},
		SLOAD: {
			execute:         opSload,
			constantGas:     params.SloadGasEIP150,
			minStack:        minStack(1, 1),
			maxStack:        maxStack(1, 1),
			computationCost: params.SloadComputationCost,
		},
		SSTORE: {
			execute:         opSstore,
			dynamicGas:      gasSStore,
			minStack:        minStack(2, 0),
			maxStack:        maxStack(2, 0),
			writes:          true,
			computationCost: params.SstoreComputationCost,
		},
		JUMP: {
			execute:         opJump,
			constantGas:     GasMidStep,
			minStack:        minStack(1, 0),
			maxStack:        maxStack(1, 0),
			jumps:           true,
			computationCost: params.JumpComputationCost,
		},
		JUMPI: {
			execute:         opJumpi,
			constantGas:     GasSlowStep,
			minStack:        minStack(2, 0),
			maxStack:        maxStack(2, 0),
			jumps:           true,
			computationCost: params.JumpiComputationCost,
		},
		PC: {
			execute:         opPc,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PcComputationCost,
		},
		MSIZE: {
			execute:         opMsize,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.MsizeComputationCost,
		},
		GAS: {
			execute:         opGas,
			constantGas:     GasQuickStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.GasComputationCost,
		},
		JUMPDEST: {
			execute:         opJumpdest,
			constantGas:     params.JumpdestGas,
			minStack:        minStack(0, 0),
			maxStack:        maxStack(0, 0),
			computationCost: params.JumpDestComputationCost,
		},
		PUSH1: {
			execute:         opPush1,
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH2: {
			execute:         makePush(2, 2),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH3: {
			execute:         makePush(3, 3),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH4: {
			execute:         makePush(4, 4),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH5: {
			execute:         makePush(5, 5),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH6: {
			execute:         makePush(6, 6),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH7: {
			execute:         makePush(7, 7),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH8: {
			execute:         makePush(8, 8),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH9: {
			execute:         makePush(9, 9),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH10: {
			execute:         makePush(10, 10),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH11: {
			execute:         makePush(11, 11),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH12: {
			execute:         makePush(12, 12),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH13: {
			execute:         makePush(13, 13),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH14: {
			execute:         makePush(14, 14),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH15: {
			execute:         makePush(15, 15),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH16: {
			execute:         makePush(16, 16),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH17: {
			execute:         makePush(17, 17),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH18: {
			execute:         makePush(18, 18),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH19: {
			execute:         makePush(19, 19),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH20: {
			execute:         makePush(20, 20),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH21: {
			execute:         makePush(21, 21),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH22: {
			execute:         makePush(22, 22),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH23: {
			execute:         makePush(23, 23),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH24: {
			execute:         makePush(24, 24),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH25: {
			execute:         makePush(25, 25),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH26: {
			execute:         makePush(26, 26),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH27: {
			execute:         makePush(27, 27),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH28: {
			execute:         makePush(28, 28),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH29: {
			execute:         makePush(29, 29),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH30: {
			execute:         makePush(30, 30),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH31: {
			execute:         makePush(31, 31),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		PUSH32: {
			execute:         makePush(32, 32),
			constantGas:     GasFastestStep,
			minStack:        minStack(0, 1),
			maxStack:        maxStack(0, 1),
			computationCost: params.PushComputationCost,
		},
		DUP1: {
			execute:         makeDup(1),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(1),
			maxStack:        maxDupStack(1),
			computationCost: params.Dup1ComputationCost,
		},
		DUP2: {
			execute:         makeDup(2),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(2),
			maxStack:        maxDupStack(2),
			computationCost: params.Dup2ComputationCost,
		},
		DUP3: {
			execute:         makeDup(3),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(3),
			maxStack:        maxDupStack(3),
			computationCost: params.Dup3ComputationCost,
		},
		DUP4: {
			execute:         makeDup(4),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(4),
			maxStack:        maxDupStack(4),
			computationCost: params.Dup4ComputationCost,
		},
		DUP5: {
			execute:         makeDup(5),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(5),
			maxStack:        maxDupStack(5),
			computationCost: params.Dup5ComputationCost,
		},
		DUP6: {
			execute:         makeDup(6),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(6),
			maxStack:        maxDupStack(6),
			computationCost: params.Dup6ComputationCost,
		},
		DUP7: {
			execute:         makeDup(7),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(7),
			maxStack:        maxDupStack(7),
			computationCost: params.Dup7ComputationCost,
		},
		DUP8: {
			execute:         makeDup(8),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(8),
			maxStack:        maxDupStack(8),
			computationCost: params.Dup8ComputationCost,
		},
		DUP9: {
			execute:         makeDup(9),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(9),
			maxStack:        maxDupStack(9),
			computationCost: params.Dup9ComputationCost,
		},
		DUP10: {
			execute:         makeDup(10),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(10),
			maxStack:        maxDupStack(10),
			computationCost: params.Dup10ComputationCost,
		},
		DUP11: {
			execute:         makeDup(11),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(11),
			maxStack:        maxDupStack(11),
			computationCost: params.Dup11ComputationCost,
		},
		DUP12: {
			execute:         makeDup(12),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(12),
			maxStack:        maxDupStack(12),
			computationCost: params.Dup12ComputationCost,
		},
		DUP13: {
			execute:         makeDup(13),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(13),
			maxStack:        maxDupStack(13),
			computationCost: params.Dup13ComputationCost,
		},
		DUP14: {
			execute:         makeDup(14),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(14),
			maxStack:        maxDupStack(14),
			computationCost: params.Dup14ComputationCost,
		},
		DUP15: {
			execute:         makeDup(15),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(15),
			maxStack:        maxDupStack(15),
			computationCost: params.Dup15ComputationCost,
		},
		DUP16: {
			execute:         makeDup(16),
			constantGas:     GasFastestStep,
			minStack:        minDupStack(16),
			maxStack:        maxDupStack(16),
			computationCost: params.Dup16ComputationCost,
		},
		SWAP1: {
			execute:         makeSwap(1),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(2),
			maxStack:        maxSwapStack(2),
			computationCost: params.Swap1ComputationCost,
		},
		SWAP2: {
			execute:         makeSwap(2),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(3),
			maxStack:        maxSwapStack(3),
			computationCost: params.Swap2ComputationCost,
		},
		SWAP3: {
			execute:         makeSwap(3),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(4),
			maxStack:        maxSwapStack(4),
			computationCost: params.Swap3ComputationCost,
		},
		SWAP4: {
			execute:         makeSwap(4),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(5),
			maxStack:        maxSwapStack(5),
			computationCost: params.Swap4ComputationCost,
		},
		SWAP5: {
			execute:         makeSwap(5),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(6),
			maxStack:        maxSwapStack(6),
			computationCost: params.Swap5ComputationCost,
		},
		SWAP6: {
			execute:         makeSwap(6),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(7),
			maxStack:        maxSwapStack(7),
			computationCost: params.Swap6ComputationCost,
		},
		SWAP7: {
			execute:         makeSwap(7),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(8),
			maxStack:        maxSwapStack(8),
			computationCost: params.Swap7ComputationCost,
		},
		SWAP8: {
			execute:         makeSwap(8),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(9),
			maxStack:        maxSwapStack(9),
			computationCost: params.Swap8ComputationCost,
		},
		SWAP9: {
			execute:         makeSwap(9),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(10),
			maxStack:        maxSwapStack(10),
			computationCost: params.Swap9ComputationCost,
		},
		SWAP10: {
			execute:         makeSwap(10),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(11),
			maxStack:        maxSwapStack(11),
			computationCost: params.Swap10ComputationCost,
		},
		SWAP11: {
			execute:         makeSwap(11),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(12),
			maxStack:        maxSwapStack(12),
			computationCost: params.Swap11ComputationCost,
		},
		SWAP12: {
			execute:         makeSwap(12),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(13),
			maxStack:        maxSwapStack(13),
			computationCost: params.Swap12ComputationCost,
		},
		SWAP13: {
			execute:         makeSwap(13),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(14),
			maxStack:        maxSwapStack(14),
			computationCost: params.Swap13ComputationCost,
		},
		SWAP14: {
			execute:         makeSwap(14),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(15),
			maxStack:        maxSwapStack(15),
			computationCost: params.Swap14ComputationCost,
		},
		SWAP15: {
			execute:         makeSwap(15),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(16),
			maxStack:        maxSwapStack(16),
			computationCost: params.Swap15ComputationCost,
		},
		SWAP16: {
			execute:         makeSwap(16),
			constantGas:     GasFastestStep,
			minStack:        minSwapStack(17),
			maxStack:        maxSwapStack(17),
			computationCost: params.Swap16ComputationCost,
		},
		LOG0: {
			execute:         makeLog(0),
			dynamicGas:      makeGasLog(0),
			minStack:        minStack(2, 0),
			maxStack:        maxStack(2, 0),
			memorySize:      memoryLog,
			writes:          true,
			computationCost: params.Log0ComputationCost,
		},
		LOG1: {
			execute:         makeLog(1),
			dynamicGas:      makeGasLog(1),
			minStack:        minStack(3, 0),
			maxStack:        maxStack(3, 0),
			memorySize:      memoryLog,
			writes:          true,
			computationCost: params.Log1ComputationCost,
		},
		LOG2: {
			execute:         makeLog(2),
			dynamicGas:      makeGasLog(2),
			minStack:        minStack(4, 0),
			maxStack:        maxStack(4, 0),
			memorySize:      memoryLog,
			writes:          true,
			computationCost: params.Log2ComputationCost,
		},
		LOG3: {
			execute:         makeLog(3),
			dynamicGas:      makeGasLog(3),
			minStack:        minStack(5, 0),
			maxStack:        maxStack(5, 0),
			memorySize:      memoryLog,
			writes:          true,
			computationCost: params.Log3ComputationCost,
		},
		LOG4: {
			execute:         makeLog(4),
			dynamicGas:      makeGasLog(4),
			minStack:        minStack(6, 0),
			maxStack:        maxStack(6, 0),
			memorySize:      memoryLog,
			writes:          true,
			computationCost: params.Log4ComputationCost,
		},
		CREATE: {
			execute:         opCreate,
			constantGas:     params.CreateGas,
			dynamicGas:      gasCreate,
			minStack:        minStack(3, 1),
			maxStack:        maxStack(3, 1),
			memorySize:      memoryCreate,
			writes:          true,
			returns:         true,
			computationCost: params.CreateComputationCost,
		},
		CALL: {
			execute:         opCall,
			constantGas:     params.CallGas,
			dynamicGas:      gasCall,
			minStack:        minStack(7, 1),
			maxStack:        maxStack(7, 1),
			memorySize:      memoryCall,
			returns:         true,
			computationCost: params.CallComputationCost,
		},
		CALLCODE: {
			execute:         opCallCode,
			constantGas:     params.CallGas,
			dynamicGas:      gasCallCode,
			minStack:        minStack(7, 1),
			maxStack:        maxStack(7, 1),
			memorySize:      memoryCall,
			returns:         true,
			computationCost: params.CallCodeComputationCost,
		},
		RETURN: {
			execute:         opReturn,
			dynamicGas:      gasReturn,
			minStack:        minStack(2, 0),
			maxStack:        maxStack(2, 0),
			memorySize:      memoryReturn,
			halts:           true,
			computationCost: params.ReturnComputationCost,
		},
		SELFDESTRUCT: {
			execute:         opSuicide,
			dynamicGas:      gasSelfdestruct,
			minStack:        minStack(1, 0),
			maxStack:        maxStack(1, 0),
			halts:           true,
			writes:          true,
			computationCost: params.SelfDestructComputationCost,
		},
	}
}
