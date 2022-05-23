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

func opSelfBalance(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, errStack *ErrStack) ([]byte, error) {
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
func opChainID(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, errStack *ErrStack) ([]byte, error) {
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
func opBaseFee(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack, errStack *ErrStack) ([]byte, error) {
	baseFee := evm.interpreter.intPool.get().Set(evm.Context.BaseFee)
	stack.push(baseFee)
	return nil, nil
}
