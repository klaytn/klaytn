// Copyright 2019 The go-ethereum Authors
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

package vm

import (
	"fmt"

	"github.com/klaytn/klaytn/params"
)

// EnableEIP enables the given EIP on the config.
// This operation writes in-place, and callers need to ensure that the globally
// defined jump tables are not polluted.
func EnableEIP(eipNum int, jt *JumpTable) error {
	switch eipNum {
	case 4399:
		enable4399(jt)
	case 3529:
		enable3529(jt)
	case 2929:
		enable2929(jt)
	case 2200:
		enable2200(jt)
	case 1884:
		enable1884(jt)
	case 1344:
		enable1344(jt)
	default:
		return fmt.Errorf("undefined eip %d", eipNum)
	}
	return nil
}

// enable1884 applies EIP-1884 to the given jump table:
// - Increase cost of BALANCE to 700
// - Increase cost of EXTCODEHASH to 700
// - Increase cost of SLOAD to 800
// - Define SELFBALANCE, with cost GasFastStep (5)
func enable1884(jt *JumpTable) {
	// Gas cost changes
	jt[SLOAD].constantGas = params.SloadGasEIP1884
	jt[BALANCE].constantGas = params.BalanceGasEIP1884
	jt[EXTCODEHASH].constantGas = params.ExtcodeHashGasEIP1884

	// New opcode
	jt[SELFBALANCE] = &operation{
		execute:         opSelfBalance,
		constantGas:     GasFastStep,
		minStack:        minStack(0, 1),
		maxStack:        maxStack(0, 1),
		computationCost: params.SelfBalanceComputationCost,
	}
}

func opSelfBalance(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	balance := evm.interpreter.intPool.get().Set(evm.StateDB.GetBalance(contract.Address()))
	stack.push(balance)
	return nil, nil
}

// enable1344 applies EIP-1344 (ChainID Opcode)
// - Adds an opcode that returns the current chain’s EIP-155 unique identifier
func enable1344(jt *JumpTable) {
	// New opcode
	jt[CHAINID] = &operation{
		execute:         opChainID,
		constantGas:     GasQuickStep,
		minStack:        minStack(0, 1),
		maxStack:        maxStack(0, 1),
		computationCost: params.ChainIDComputationCost,
	}
}

// opChainID implements CHAINID opcode
func opChainID(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	chainId := evm.interpreter.intPool.get().Set(evm.chainConfig.ChainID)
	stack.push(chainId)
	return nil, nil
}

// enable2200 applies EIP-2200 (Rebalance net-metered SSTORE)
func enable2200(jt *JumpTable) {
	jt[SLOAD].constantGas = params.SloadGasEIP2200
	jt[SSTORE].dynamicGas = gasSStoreEIP2200
}

// enableIstanbulComputationCostModification modifies ADDMOD, MULMOD, NOT, XOR, SHL, SHR, SAR computation cost
// The modification is activated with istanbulCompatible change activation.
func enableIstanbulComputationCostModification(jt *JumpTable) {
	jt[ADDMOD].computationCost = params.AddmodComputationCostIstanbul
	jt[MULMOD].computationCost = params.MulmodComputationCostIstanbul
	jt[NOT].computationCost = params.NotComputationCostIstanbul
	jt[XOR].computationCost = params.XorComputationCostIstanbul
	jt[SHL].computationCost = params.ShlComputationCostIstanbul
	jt[SHR].computationCost = params.ShrComputationCostIstanbul
	jt[SAR].computationCost = params.SarComputationCostIstanbul
}

// enable3198 applies EIP-3198 (BASEFEE Opcode)
// - Adds an opcode that returns the current block's base fee.
func enable3198(jt *JumpTable) {
	// New opcode
	jt[BASEFEE] = &operation{
		execute:         opBaseFee,
		constantGas:     GasQuickStep,
		minStack:        minStack(0, 1),
		maxStack:        maxStack(0, 1),
		computationCost: params.BaseFeeComputationCost,
	}
}

// opBaseFee implements BASEFEE opcode
func opBaseFee(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	baseFee := evm.interpreter.intPool.get().Set(evm.Context.BaseFee)
	stack.push(baseFee)
	return nil, nil
}

// enable2929 enables "EIP-2929: Gas cost increases for state access opcodes"
// https://eips.ethereum.org/EIPS/eip-2929
func enable2929(jt *JumpTable) {
	jt[SSTORE].dynamicGas = gasSStoreEIP2929

	jt[SLOAD].constantGas = 0
	jt[SLOAD].dynamicGas = gasSLoadEIP2929

	jt[EXTCODECOPY].constantGas = params.WarmStorageReadCostEIP2929
	jt[EXTCODECOPY].dynamicGas = gasExtCodeCopyEIP2929

	jt[EXTCODESIZE].constantGas = params.WarmStorageReadCostEIP2929
	jt[EXTCODESIZE].dynamicGas = gasEip2929AccountCheck

	jt[EXTCODEHASH].constantGas = params.WarmStorageReadCostEIP2929
	jt[EXTCODEHASH].dynamicGas = gasEip2929AccountCheck

	jt[BALANCE].constantGas = params.WarmStorageReadCostEIP2929
	jt[BALANCE].dynamicGas = gasEip2929AccountCheck

	jt[CALL].constantGas = params.WarmStorageReadCostEIP2929
	jt[CALL].dynamicGas = gasCallEIP2929

	jt[CALLCODE].constantGas = params.WarmStorageReadCostEIP2929
	jt[CALLCODE].dynamicGas = gasCallCodeEIP2929

	jt[STATICCALL].constantGas = params.WarmStorageReadCostEIP2929
	jt[STATICCALL].dynamicGas = gasStaticCallEIP2929

	jt[DELEGATECALL].constantGas = params.WarmStorageReadCostEIP2929
	jt[DELEGATECALL].dynamicGas = gasDelegateCallEIP2929

	// This was previously part of the dynamic cost, but we're using it as a constantGas
	// factor here
	jt[SELFDESTRUCT].constantGas = params.SelfdestructGas
	jt[SELFDESTRUCT].dynamicGas = gasSelfdestructEIP2929
}

// enable3529 enabled "EIP-3529: Reduction in refunds":
// - Removes refunds for selfdestructs
// - Reduces refunds for SSTORE
// - Reduces max refunds to 20% gas
func enable3529(jt *JumpTable) {
	jt[SSTORE].dynamicGas = gasSStoreEIP3529
	jt[SELFDESTRUCT].dynamicGas = gasSelfdestructEIP3529
}

// enable4399 applies EIP-4399 (PREVRANDAO Opcode)
// - Change the 0x44 opcode from returning difficulty value to returning prev blockhash value
func enable4399(jt *JumpTable) {
	jt[PREVRANDAO] = &operation{
		execute:         opRandom,
		constantGas:     GasQuickStep,
		minStack:        minStack(0, 1),
		maxStack:        maxStack(0, 1),
		computationCost: params.RandomComputationCost,
	}
}
