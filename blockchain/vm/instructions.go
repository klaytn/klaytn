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
// This file is derived from core/vm/instructions.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/params"
)

var (
	bigZero = new(big.Int)
	tt255   = math.BigPow(2, 255)
)

func opAdd(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	math.U256(y.Add(x, y))

	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opSub(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	math.U256(y.Sub(x, y))

	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opMul(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.pop()
	scope.Stack.push(math.U256(x.Mul(x, y)))

	evm.interpreter.intPool.put(y)

	return nil, nil
}

func opDiv(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	if y.Sign() != 0 {
		math.U256(y.Div(x, y))
	} else {
		y.SetUint64(0)
	}
	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opSdiv(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := math.S256(scope.Stack.pop()), math.S256(scope.Stack.pop())
	res := evm.interpreter.intPool.getZero()

	if y.Sign() == 0 || x.Sign() == 0 {
		scope.Stack.push(res)
	} else {
		if x.Sign() != y.Sign() {
			res.Div(x.Abs(x), y.Abs(y))
			res.Neg(res)
		} else {
			res.Div(x.Abs(x), y.Abs(y))
		}
		scope.Stack.push(math.U256(res))
	}
	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opMod(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.pop()
	if y.Sign() == 0 {
		scope.Stack.push(x.SetUint64(0))
	} else {
		scope.Stack.push(math.U256(x.Mod(x, y)))
	}
	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opSmod(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := math.S256(scope.Stack.pop()), math.S256(scope.Stack.pop())
	res := evm.interpreter.intPool.getZero()

	if y.Sign() == 0 {
		scope.Stack.push(res)
	} else {
		if x.Sign() < 0 {
			res.Mod(x.Abs(x), y.Abs(y))
			res.Neg(res)
		} else {
			res.Mod(x.Abs(x), y.Abs(y))
		}
		scope.Stack.push(math.U256(res))
	}
	evm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opExp(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	base, exponent := scope.Stack.pop(), scope.Stack.pop()
	// some shortcuts
	cmpToOne := exponent.Cmp(big1)
	if cmpToOne < 0 { // Exponent is zero
		// x ^ 0 == 1
		scope.Stack.push(base.SetUint64(1))
	} else if base.Sign() == 0 {
		// 0 ^ y, if y != 0, == 0
		scope.Stack.push(base.SetUint64(0))
	} else if cmpToOne == 0 { // Exponent is one
		// x ^ 1 == x
		scope.Stack.push(base)
	} else {
		scope.Stack.push(math.Exp(base, exponent))
		evm.interpreter.intPool.put(base)
	}
	evm.interpreter.intPool.put(exponent)

	return nil, nil
}

func opSignExtend(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	back := scope.Stack.pop()
	if back.Cmp(big.NewInt(31)) < 0 {
		bit := uint(back.Uint64()*8 + 7)
		num := scope.Stack.pop()
		mask := back.Lsh(common.Big1, bit)
		mask.Sub(mask, common.Big1)
		if num.Bit(int(bit)) > 0 {
			num.Or(num, mask.Not(mask))
		} else {
			num.And(num, mask)
		}

		scope.Stack.push(math.U256(num))
	}

	evm.interpreter.intPool.put(back)
	return nil, nil
}

func opNot(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x := scope.Stack.Peek()
	math.U256(x.Not(x))
	return nil, nil
}

func opLt(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	if x.Cmp(y) < 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opGt(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	if x.Cmp(y) > 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opSlt(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()

	xSign := x.Cmp(tt255)
	ySign := y.Cmp(tt255)

	switch {
	case xSign >= 0 && ySign < 0:
		y.SetUint64(1)

	case xSign < 0 && ySign >= 0:
		y.SetUint64(0)

	default:
		if x.Cmp(y) < 0 {
			y.SetUint64(1)
		} else {
			y.SetUint64(0)
		}
	}
	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opSgt(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()

	xSign := x.Cmp(tt255)
	ySign := y.Cmp(tt255)

	switch {
	case xSign >= 0 && ySign < 0:
		y.SetUint64(0)

	case xSign < 0 && ySign >= 0:
		y.SetUint64(1)

	default:
		if x.Cmp(y) > 0 {
			y.SetUint64(1)
		} else {
			y.SetUint64(0)
		}
	}
	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opEq(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	if x.Cmp(y) == 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opIszero(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x := scope.Stack.Peek()
	if x.Sign() > 0 {
		x.SetUint64(0)
	} else {
		x.SetUint64(1)
	}
	return nil, nil
}

func opAnd(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.pop()
	scope.Stack.push(x.And(x, y))

	evm.interpreter.intPool.put(y)
	return nil, nil
}

func opOr(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	y.Or(x, y)

	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opXor(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y := scope.Stack.pop(), scope.Stack.Peek()
	y.Xor(x, y)

	evm.interpreter.intPool.put(x)
	return nil, nil
}

func opByte(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	th, val := scope.Stack.pop(), scope.Stack.Peek()
	if th.Cmp(common.Big32) < 0 {
		b := math.Byte(val, 32, int(th.Int64()))
		val.SetUint64(uint64(b))
	} else {
		val.SetUint64(0)
	}
	evm.interpreter.intPool.put(th)
	return nil, nil
}

func opAddmod(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y, z := scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop()
	if z.Cmp(bigZero) > 0 {
		x.Add(x, y)
		x.Mod(x, z)
		scope.Stack.push(math.U256(x))
	} else {
		scope.Stack.push(x.SetUint64(0))
	}
	evm.interpreter.intPool.put(y, z)
	return nil, nil
}

func opMulmod(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	x, y, z := scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop()
	if z.Cmp(bigZero) > 0 {
		x.Mul(x, y)
		x.Mod(x, z)
		scope.Stack.push(math.U256(x))
	} else {
		scope.Stack.push(x.SetUint64(0))
	}
	evm.interpreter.intPool.put(y, z)
	return nil, nil
}

// opSHL implements Shift Left
// The SHL instruction (shift left) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the left by arg1 number of bits.
func opSHL(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := math.U256(scope.Stack.pop()), math.U256(scope.Stack.Peek())
	defer evm.interpreter.intPool.put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		value.SetUint64(0)
		return nil, nil
	}
	n := uint(shift.Uint64())
	math.U256(value.Lsh(value, n))

	return nil, nil
}

// opSHR implements Logical Shift Right
// The SHR instruction (logical shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with zero fill.
func opSHR(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := math.U256(scope.Stack.pop()), math.U256(scope.Stack.Peek())
	defer evm.interpreter.intPool.put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		value.SetUint64(0)
		return nil, nil
	}
	n := uint(shift.Uint64())
	math.U256(value.Rsh(value, n))

	return nil, nil
}

// opSAR implements Arithmetic Shift Right
// The SAR instruction (arithmetic shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with sign extension.
func opSAR(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// Note, S256 returns (potentially) a new bigint, so we're popping, not peeking this one
	shift, value := math.U256(scope.Stack.pop()), math.S256(scope.Stack.pop())
	defer evm.interpreter.intPool.put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		if value.Sign() >= 0 {
			value.SetUint64(0)
		} else {
			value.SetInt64(-1)
		}
		scope.Stack.push(math.U256(value))
		return nil, nil
	}
	n := uint(shift.Uint64())
	value.Rsh(value, n)
	scope.Stack.push(math.U256(value))

	return nil, nil
}

func opSha3(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	offset, size := scope.Stack.pop(), scope.Stack.pop()
	data := scope.Memory.GetPtr(offset.Int64(), size.Int64())

	if evm.interpreter.hasher == nil {
		evm.interpreter.hasher = sha3.NewKeccak256().(keccakState)
	} else {
		evm.interpreter.hasher.Reset()
	}
	evm.interpreter.hasher.Write(data)
	evm.interpreter.hasher.Read(evm.interpreter.hasherBuf[:])

	if evm.Config.EnablePreimageRecording {
		evm.StateDB.AddPreimage(evm.interpreter.hasherBuf, data)
	}
	scope.Stack.push(evm.interpreter.intPool.get().SetBytes(evm.interpreter.hasherBuf[:]))

	evm.interpreter.intPool.put(offset, size)
	return nil, nil
}

func opAddress(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetBytes(scope.Contract.Address().Bytes()))
	return nil, nil
}

func opBalance(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	slot := scope.Stack.Peek()
	slot.Set(evm.StateDB.GetBalance(common.BigToAddress(slot)))
	return nil, nil
}

func opOrigin(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetBytes(evm.Origin.Bytes()))
	return nil, nil
}

func opCaller(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetBytes(scope.Contract.Caller().Bytes()))
	return nil, nil
}

func opCallValue(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().Set(scope.Contract.value))
	return nil, nil
}

func opCallDataLoad(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetBytes(getDataBig(scope.Contract.Input, scope.Stack.pop(), big32)))
	return nil, nil
}

func opCallDataSize(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetInt64(int64(len(scope.Contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	var (
		memOffset  = scope.Stack.pop()
		dataOffset = scope.Stack.pop()
		length     = scope.Stack.pop()
	)
	scope.Memory.Set(memOffset.Uint64(), length.Uint64(), getDataBig(scope.Contract.Input, dataOffset, length))

	evm.interpreter.intPool.put(memOffset, dataOffset, length)
	return nil, nil
}

func opReturnDataSize(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetUint64(uint64(len(evm.interpreter.returnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	var (
		memOffset  = scope.Stack.pop()
		dataOffset = scope.Stack.pop()
		length     = scope.Stack.pop()

		end = evm.interpreter.intPool.get().Add(dataOffset, length)
	)
	defer evm.interpreter.intPool.put(memOffset, dataOffset, length, end)

	if !end.IsUint64() || uint64(len(evm.interpreter.returnData)) < end.Uint64() {
		return nil, ErrReturnDataOutOfBounds // TODO-Klaytn-Issue615
	}
	scope.Memory.Set(memOffset.Uint64(), length.Uint64(), evm.interpreter.returnData[dataOffset.Uint64():end.Uint64()])

	return nil, nil
}

func opExtCodeSize(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	slot := scope.Stack.Peek()
	slot.SetUint64(uint64(evm.StateDB.GetCodeSize(common.BigToAddress(slot))))

	return nil, nil
}

func opCodeSize(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	l := evm.interpreter.intPool.get().SetInt64(int64(len(scope.Contract.Code)))
	scope.Stack.push(l)

	return nil, nil
}

func opCodeCopy(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	var (
		memOffset  = scope.Stack.pop()
		codeOffset = scope.Stack.pop()
		length     = scope.Stack.pop()
	)
	codeCopy := getDataBig(scope.Contract.Code, codeOffset, length)
	scope.Memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	evm.interpreter.intPool.put(memOffset, codeOffset, length)
	return nil, nil
}

func opExtCodeCopy(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	var (
		addr       = common.BigToAddress(scope.Stack.pop())
		memOffset  = scope.Stack.pop()
		codeOffset = scope.Stack.pop()
		length     = scope.Stack.pop()
	)
	codeCopy := getDataBig(evm.StateDB.GetCode(addr), codeOffset, length)
	scope.Memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	evm.interpreter.intPool.put(memOffset, codeOffset, length)
	return nil, nil
}

// opExtCodeHash returns the code hash of a specified account.
// There are several cases when the function is called, while we can relay everything
// to `state.GetCodeHash` function to ensure the correctness.
//
//	(1) Caller tries to get the code hash of a normal contract account, state
//
// should return the relative code hash and set it as the result.
//
//	(2) Caller tries to get the code hash of a non-existent account, state should
//
// return common.Hash{} and zero will be set as the result.
//
//	(3) Caller tries to get the code hash for an account without contract code,
//
// state should return emptyCodeHash(0xc5d246...) as the result.
//
//	(4) Caller tries to get the code hash of a precompiled account, the result
//
// should be zero or emptyCodeHash.
//
// It is worth noting that in order to avoid unnecessary create and clean,
// all precompile accounts on mainnet have been transferred 1 wei, so the return
// here should be emptyCodeHash.
// If the precompile account is not transferred any amount on a private or
// customized chain, the return value will be zero.
//
//	(5) Caller tries to get the code hash for an account which is marked as self-destructed
//
// in the current transaction, the code hash of this account should be returned.
//
//	(6) Caller tries to get the code hash for an account which is marked as deleted,
//
// this account should be regarded as a non-existent account and zero should be returned.
func opExtCodeHash(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	slot := scope.Stack.Peek()
	slot.SetBytes(evm.StateDB.GetCodeHash(common.BigToAddress(slot)).Bytes())
	return nil, nil
}

func opGasprice(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().Set(evm.GasPrice))
	return nil, nil
}

func opBlockhash(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	num := scope.Stack.pop()

	n := evm.interpreter.intPool.get().Sub(evm.Context.BlockNumber, common.Big257)
	if num.Cmp(n) > 0 && num.Cmp(evm.Context.BlockNumber) < 0 {
		scope.Stack.push(evm.Context.GetHash(num.Uint64()).Big())
	} else {
		scope.Stack.push(evm.interpreter.intPool.getZero())
	}
	evm.interpreter.intPool.put(num, n)
	return nil, nil
}

func opCoinbase(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetBytes(evm.Context.Coinbase.Bytes()))
	return nil, nil
}

func opTimestamp(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(math.U256(evm.interpreter.intPool.get().Set(evm.Context.Time)))
	return nil, nil
}

func opNumber(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(math.U256(evm.interpreter.intPool.get().Set(evm.Context.BlockNumber)))
	return nil, nil
}

func opDifficulty(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(math.U256(evm.interpreter.intPool.get().Set(evm.Context.BlockScore)))
	return nil, nil
}

func opRandom(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// evm.BlockNumber.Uint64() is always greater than or equal to 1
	// since evm will not run on the genesis block
	scope.Stack.push(evm.Context.GetHash(evm.Context.BlockNumber.Uint64() - 1).Big())
	return nil, nil
}

func opGasLimit(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(math.U256(evm.interpreter.intPool.get().SetUint64(evm.Context.GasLimit)))
	return nil, nil
}

func opPop(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	evm.interpreter.intPool.put(scope.Stack.pop())
	return nil, nil
}

func opMload(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	v := scope.Stack.Peek()
	offset := v.Int64()
	v.SetBytes(scope.Memory.GetPtr(offset, 32))

	return nil, nil
}

func opMstore(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// pop value of the stack
	mStart, val := scope.Stack.pop(), scope.Stack.pop()
	scope.Memory.Set32(mStart.Uint64(), val)

	evm.interpreter.intPool.put(mStart, val)
	return nil, nil
}

func opMstore8(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	off, val := scope.Stack.pop().Int64(), scope.Stack.pop().Int64()
	scope.Memory.store[off] = byte(val & 0xff)

	return nil, nil
}

func opSload(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	loc := scope.Stack.Peek()
	val := evm.StateDB.GetState(scope.Contract.Address(), common.BigToHash(loc))
	loc.SetBytes(val.Bytes())
	return nil, nil
}

func opSstore(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	loc := common.BigToHash(scope.Stack.pop())
	val := scope.Stack.pop()
	evm.StateDB.SetState(scope.Contract.Address(), loc, common.BigToHash(val))

	evm.interpreter.intPool.put(val)
	return nil, nil
}

func opJump(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	pos := scope.Stack.pop()
	if !scope.Contract.validJumpdest(pos) {
		return nil, ErrInvalidJump
	}
	*pc = pos.Uint64()

	evm.interpreter.intPool.put(pos)
	return nil, nil
}

func opJumpi(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	pos, cond := scope.Stack.pop(), scope.Stack.pop()
	if cond.Sign() != 0 {
		if !scope.Contract.validJumpdest(pos) {
			return nil, ErrInvalidJump
		}
		*pc = pos.Uint64()
	} else {
		*pc++
	}

	evm.interpreter.intPool.put(pos, cond)
	return nil, nil
}

func opJumpdest(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	return nil, nil
}

func opPc(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetInt64(int64(scope.Memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	scope.Stack.push(evm.interpreter.intPool.get().SetUint64(scope.Contract.Gas))
	return nil, nil
}

func opCreate(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	var (
		value        = scope.Stack.pop()
		offset, size = scope.Stack.pop(), scope.Stack.pop()
		input        = scope.Memory.GetCopy(offset.Int64(), size.Int64())
		gas          = scope.Contract.Gas
	)

	// This is from EIP150.
	gas -= gas / 64

	scope.Contract.UseGas(gas)
	res, addr, returnGas, suberr := evm.Create(scope.Contract, input, gas, value, params.CodeFormatEVM)
	// Push item on the stack based on the returned error. If the ruleset is
	// homestead we must check for CodeStoreOutOfGasError (homestead only
	// rule) and treat as an error, if the ruleset is frontier we must
	// ignore this error and pretend the operation was successful.
	if suberr != nil {
		scope.Stack.push(evm.interpreter.intPool.getZero())
	} else {
		scope.Stack.push(evm.interpreter.intPool.get().SetBytes(addr.Bytes()))
	}
	scope.Contract.Gas += returnGas
	evm.interpreter.intPool.put(value, offset, size)

	if suberr == ErrExecutionReverted {
		return res, nil
	}
	return nil, nil
}

func opCreate2(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	var (
		endowment    = scope.Stack.pop()
		offset, size = scope.Stack.pop(), scope.Stack.pop()
		salt         = scope.Stack.pop()
		input        = scope.Memory.GetCopy(offset.Int64(), size.Int64())
		gas          = scope.Contract.Gas
	)

	// Apply EIP150
	gas -= gas / 64
	scope.Contract.UseGas(gas)
	res, addr, returnGas, suberr := evm.Create2(scope.Contract, input, gas, endowment, salt, params.CodeFormatEVM)
	// Push item on the stack based on the returned error.
	if suberr != nil {
		scope.Stack.push(evm.interpreter.intPool.getZero())
	} else {
		scope.Stack.push(evm.interpreter.intPool.get().SetBytes(addr.Bytes()))
	}
	scope.Contract.Gas += returnGas
	evm.interpreter.intPool.put(endowment, offset, size, salt)

	if suberr == ErrExecutionReverted {
		return res, nil
	}
	return nil, nil
}

func opCall(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// Pop gas. The actual gas in evm.callGasTemp.
	evm.interpreter.intPool.put(scope.Stack.pop())
	gas := evm.callGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop()
	toAddr := common.BigToAddress(addr)
	value = math.U256(value)
	// Get the arguments from the memory.
	args := scope.Memory.GetPtr(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	ret, returnGas, err := evm.Call(scope.Contract, toAddr, args, gas, value)
	if err != nil {
		scope.Stack.push(evm.interpreter.intPool.getZero())
	} else {
		scope.Stack.push(evm.interpreter.intPool.get().SetUint64(1))
	}
	if err == nil || err == ErrExecutionReverted {
		ret = common.CopyBytes(ret)
		scope.Memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	scope.Contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opCallCode(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// Pop gas. The actual gas is in evm.callGasTemp.
	evm.interpreter.intPool.put(scope.Stack.pop())
	gas := evm.callGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop()
	toAddr := common.BigToAddress(addr)
	value = math.U256(value)
	// Get arguments from the memory.
	args := scope.Memory.GetPtr(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	ret, returnGas, err := evm.CallCode(scope.Contract, toAddr, args, gas, value)
	if err != nil {
		scope.Stack.push(evm.interpreter.intPool.getZero())
	} else {
		scope.Stack.push(evm.interpreter.intPool.get().SetUint64(1))
	}
	if err == nil || err == ErrExecutionReverted {
		ret = common.CopyBytes(ret)
		scope.Memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	scope.Contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opDelegateCall(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// Pop gas. The actual gas is in evm.callGasTemp.
	evm.interpreter.intPool.put(scope.Stack.pop())
	gas := evm.callGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop()
	toAddr := common.BigToAddress(addr)
	// Get arguments from the memory.
	args := scope.Memory.GetPtr(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := evm.DelegateCall(scope.Contract, toAddr, args, gas)
	if err != nil {
		scope.Stack.push(evm.interpreter.intPool.getZero())
	} else {
		scope.Stack.push(evm.interpreter.intPool.get().SetUint64(1))
	}
	if err == nil || err == ErrExecutionReverted {
		ret = common.CopyBytes(ret)
		scope.Memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	scope.Contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opStaticCall(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// Pop gas. The actual gas is in evm.callGasTemp.
	evm.interpreter.intPool.put(scope.Stack.pop())
	gas := evm.callGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop(), scope.Stack.pop()
	toAddr := common.BigToAddress(addr)
	// Get arguments from the memory.
	args := scope.Memory.GetPtr(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := evm.StaticCall(scope.Contract, toAddr, args, gas)
	if err != nil {
		scope.Stack.push(evm.interpreter.intPool.getZero())
	} else {
		scope.Stack.push(evm.interpreter.intPool.get().SetUint64(1))
	}
	if err == nil || err == ErrExecutionReverted {
		ret = common.CopyBytes(ret)
		scope.Memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	scope.Contract.Gas += returnGas

	evm.interpreter.intPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opReturn(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	offset, size := scope.Stack.pop(), scope.Stack.pop()
	ret := scope.Memory.GetPtr(offset.Int64(), size.Int64())

	evm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opRevert(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	offset, size := scope.Stack.pop(), scope.Stack.pop()
	ret := scope.Memory.GetPtr(offset.Int64(), size.Int64())

	evm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opStop(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	return nil, nil
}

func opSelfdestruct(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	// TODO-klaytn: call frame tracing https://github.com/ethereum/go-ethereum/pull/23087
	// beneficiary := scope.Stack.pop()
	balance := evm.StateDB.GetBalance(scope.Contract.Address())
	evm.StateDB.AddBalance(common.BigToAddress(scope.Stack.pop()), balance)
	evm.StateDB.SelfDestruct(scope.Contract.Address())

	// if evm.interpreter.cfg.Debug {
	// 	evm.interpreter.cfg.Tracer.CaptureEnter(SELFDESTRUCT, scope.Contract.Address(), common.BigToAddress(beneficiary), []byte{}, 0, balance)
	// 	evm.interpreter.cfg.Tracer.CaptureExit([]byte{}, 0, nil)
	// }

	return nil, nil
}

// opPush1 is a specialized version of pushN
func opPush1(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	var (
		codeLen = uint64(len(scope.Contract.Code))
		integer = evm.interpreter.intPool.get()
	)
	*pc += 1
	if *pc < codeLen {
		scope.Stack.push(integer.SetUint64(uint64(scope.Contract.Code[*pc])))
	} else {
		scope.Stack.push(integer.SetUint64(0))
	}
	return nil, nil
}

func opSelfdestruct6780(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
	if evm.interpreter.readOnly {
		return nil, ErrWriteProtection
	}
	beneficiary := scope.Stack.pop()
	balance := evm.StateDB.GetBalance(scope.Contract.Address())
	evm.StateDB.SubBalance(scope.Contract.Address(), balance)
	evm.StateDB.AddBalance(common.BigToAddress(beneficiary), balance)
	evm.StateDB.SelfDestruct6780(scope.Contract.Address())

	// TODO-klaytn: call frame tracing https://github.com/ethereum/go-ethereum/pull/23087
	// if tracer := interpreter.evm.Config.Tracer; tracer != nil {
	// 	tracer.CaptureEnter(SELFDESTRUCT, scope.Contract.Address(), common.BigToAddress(beneficiary), []byte{}, 0, balance)
	// 	tracer.CaptureExit([]byte{}, 0, nil)
	// }

	return nil, nil
}

// following functions are used by the instruction jump  table

// make log instruction function
func makeLog(size int) executionFunc {
	return func(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
		topics := make([]common.Hash, size)
		mStart, mSize := scope.Stack.pop(), scope.Stack.pop()
		for i := 0; i < size; i++ {
			topics[i] = common.BigToHash(scope.Stack.pop())
		}

		d := scope.Memory.GetCopy(mStart.Int64(), mSize.Int64())
		evm.StateDB.AddLog(&types.Log{
			Address: scope.Contract.Address(),
			Topics:  topics,
			Data:    d,
			// This is a non-consensus field, but assigned here because
			// blockchain/state doesn't know the current block number.
			BlockNumber: evm.Context.BlockNumber.Uint64(),
		})

		evm.interpreter.intPool.put(mStart, mSize)
		return nil, nil
	}
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
		codeLen := len(scope.Contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := evm.interpreter.intPool.get()
		scope.Stack.push(integer.SetBytes(common.RightPadBytes(scope.Contract.Code[startMin:endMin], pushByteSize)))

		*pc += size
		return nil, nil
	}
}

// make dup instruction function
func makeDup(size int64) executionFunc {
	return func(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
		scope.Stack.dup(evm.interpreter.intPool, int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size++
	return func(pc *uint64, evm *EVM, scope *ScopeContext) ([]byte, error) {
		scope.Stack.swap(int(size))
		return nil, nil
	}
}
