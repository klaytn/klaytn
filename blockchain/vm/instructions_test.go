// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from core/vm/instructions_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

type TwoOperandTestcase struct {
	X        string
	Y        string
	Expected string
}

type twoOperandParams struct {
	x string
	y string
}

var commonParams []*twoOperandParams
var twoOpMethods map[string]executionFunc

func init() {

	// Params is a list of common edgecases that should be used for some common tests
	params := []string{
		"0000000000000000000000000000000000000000000000000000000000000000", // 0
		"0000000000000000000000000000000000000000000000000000000000000001", // +1
		"0000000000000000000000000000000000000000000000000000000000000005", // +5
		"7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", // + max -1
		"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", // + max
		"8000000000000000000000000000000000000000000000000000000000000000", // - max
		"8000000000000000000000000000000000000000000000000000000000000001", // - max+1
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffb", // - 5
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", // - 1
	}
	// Params are combined so each param is used on each 'side'
	commonParams = make([]*twoOperandParams, len(params)*len(params))
	for i, x := range params {
		for j, y := range params {
			commonParams[i*len(params)+j] = &twoOperandParams{x, y}
		}
	}
	twoOpMethods = map[string]executionFunc{
		"add":     opAdd,
		"sub":     opSub,
		"mul":     opMul,
		"div":     opDiv,
		"sdiv":    opSdiv,
		"mod":     opMod,
		"smod":    opSmod,
		"exp":     opExp,
		"signext": opSignExtend,
		"lt":      opLt,
		"gt":      opGt,
		"slt":     opSlt,
		"sgt":     opSgt,
		"eq":      opEq,
		"and":     opAnd,
		"or":      opOr,
		"xor":     opXor,
		"byte":    opByte,
		"shl":     opSHL,
		"shr":     opSHR,
		"sar":     opSAR,
	}
}

func testTwoOperandOp(t *testing.T, tests []TwoOperandTestcase, opFn executionFunc, name string) {

	var (
		env            = NewEVM(Context{}, nil, params.TestChainConfig, &Config{})
		stack          = newstack()
		pc             = uint64(0)
		evmInterpreter = env.interpreter
	)
	// Stuff a couple of nonzero bigints into pool, to ensure that ops do not rely on pooled integers to be zero
	evmInterpreter.intPool = poolOfIntPools.get()
	evmInterpreter.intPool.put(big.NewInt(-1337))
	evmInterpreter.intPool.put(big.NewInt(-1337))
	evmInterpreter.intPool.put(big.NewInt(-1337))

	for i, test := range tests {
		x := new(big.Int).SetBytes(common.Hex2Bytes(test.X))
		y := new(big.Int).SetBytes(common.Hex2Bytes(test.Y))
		expected := new(big.Int).SetBytes(common.Hex2Bytes(test.Expected))
		stack.push(x)
		stack.push(y)
		opFn(&pc, env, nil, nil, stack)
		actual := stack.pop()

		if actual.Cmp(expected) != 0 {
			t.Errorf("Testcase %v %d, %v(%x, %x): expected  %x, got %x", name, i, name, x, y, expected, actual)
		}
		// Check pool usage
		// 1.pool is not allowed to contain anything on the stack
		// 2.pool is not allowed to contain the same pointers twice
		if evmInterpreter.intPool.pool.len() > 0 {

			poolvals := make(map[*big.Int]struct{})
			poolvals[actual] = struct{}{}

			for evmInterpreter.intPool.pool.len() > 0 {
				key := evmInterpreter.intPool.get()
				if _, exist := poolvals[key]; exist {
					t.Errorf("Testcase %v %d, pool contains double-entry", name, i)
				}
				poolvals[key] = struct{}{}
			}
		}
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func TestByteOp(t *testing.T) {
	tests := []TwoOperandTestcase{
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", "00", "AB"},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", "01", "CD"},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", "00", "00"},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", "01", "CD"},
		{"0000000000000000000000000000000000000000000000000000000000102030", "1F", "30"},
		{"0000000000000000000000000000000000000000000000000000000000102030", "1E", "20"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "20", "00"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "FFFFFFFFFFFFFFFF", "00"},
	}
	testTwoOperandOp(t, tests, opByte, "byte")
}

func TestSHL(t *testing.T) {
	// Testcases from https://github.com/ethereum/EIPs/blob/master/EIPS/eip-145.md#shl-shift-left
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000002"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "ff", "8000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "0101", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "8000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
	}
	testTwoOperandOp(t, tests, opSHL, "shl")
}

func TestSHR(t *testing.T) {
	// Testcases from https://github.com/ethereum/EIPs/blob/master/EIPS/eip-145.md#shr-logical-shift-right
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "00", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "01", "4000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "ff", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0101", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
	}
	testTwoOperandOp(t, tests, opSHR, "shr")
}

func TestSAR(t *testing.T) {
	// Testcases from https://github.com/ethereum/EIPs/blob/master/EIPS/eip-145.md#sar-arithmetic-shift-right
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "00", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "01", "c000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "ff", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0100", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0101", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000000000000000000000000000", "fe", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "f8", "000000000000000000000000000000000000000000000000000000000000007f"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "fe", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
	}

	testTwoOperandOp(t, tests, opSAR, "sar")
}

// getResult is a convenience function to generate the expected values
func getResult(args []*twoOperandParams, opFn executionFunc) []TwoOperandTestcase {
	var (
		env         = NewEVM(Context{}, nil, params.TestChainConfig, &Config{})
		stack       = newstack()
		pc          = uint64(0)
		interpreter = env.interpreter
	)
	interpreter.intPool = poolOfIntPools.get()
	result := make([]TwoOperandTestcase, len(args))
	for i, param := range args {
		x := new(big.Int).SetBytes(common.Hex2Bytes(param.x))
		y := new(big.Int).SetBytes(common.Hex2Bytes(param.y))
		stack.push(x)
		stack.push(y)
		opFn(&pc, env, nil, nil, stack)
		actual := stack.pop()
		result[i] = TwoOperandTestcase{param.x, param.y, fmt.Sprintf("%064x", actual)}
	}
	return result
}

// utility function to fill the json-file with testcases
// Enable this test to generate the 'testcases_xx.json' files
func xTestWriteExpectedValues(t *testing.T) {
	for name, method := range twoOpMethods {
		data, err := json.Marshal(getResult(commonParams, method))
		if err != nil {
			t.Fatal(err)
		}
		_ = ioutil.WriteFile(fmt.Sprintf("testdata/testcases_%v.json", name), data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Fatal("This test should not be activated")
}

// TestJsonTestcases runs through all the testcases defined as json-files
func TestJsonTestcases(t *testing.T) {
	for name := range twoOpMethods {
		data, err := ioutil.ReadFile(fmt.Sprintf("testdata/testcases_%v.json", name))
		if err != nil {
			t.Fatal("Failed to read file", err)
		}
		var testcases []TwoOperandTestcase
		json.Unmarshal(data, &testcases)
		testTwoOperandOp(t, testcases, twoOpMethods[name], name)
	}
}

func initStateDB(db database.DBManager) *state.StateDB {
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb, nil)

	contractAddress := common.HexToAddress("0x18f30de96ce789fe778b9a5f420f6fdbbd9b34d8")
	code := "60ca60205260005b612710811015630000004557602051506020515060205150602051506020515060205150602051506020515060205150602051506001016300000007565b00"
	statedb.CreateSmartContractAccount(contractAddress, params.CodeFormatEVM, params.Rules{})
	statedb.SetCode(contractAddress, common.Hex2Bytes(code))
	stateHash := common.HexToHash("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")
	statedb.SetState(contractAddress, stateHash, stateHash)
	statedb.SetBalance(contractAddress, big.NewInt(1000))
	statedb.SetNonce(contractAddress, uint64(1))

	{
		// create a contract having a STOP operation.
		contractAddress := common.HexToAddress("0x18f30de96ce789fe778b9a5f420f6fdbbd9b34d9")
		code := "00"
		statedb.CreateSmartContractAccount(contractAddress, params.CodeFormatEVM, params.Rules{})
		statedb.SetCode(contractAddress, common.Hex2Bytes(code))
		statedb.SetBalance(contractAddress, big.NewInt(1000))
		statedb.SetNonce(contractAddress, uint64(1))
	}

	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false)
	statedb, _ = state.New(root, sdb, nil)

	return statedb
}

func opBenchmark(bench *testing.B, op func(pc *uint64, evm *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error), args ...string) {
	var (
		initialCall = true
		canTransfer = func(db StateDB, address common.Address, amount *big.Int) bool {
			if initialCall {
				initialCall = false
				return true
			}
			return db.GetBalance(address).Cmp(amount) >= 0
		}

		ctx = Context{
			BlockNumber: big1024,
			BlockScore:  big.NewInt(0),
			Coinbase:    common.HexToAddress("0xf4b0cb429b7d341bf467f2d51c09b64cd9add37c"),
			GasPrice:    big.NewInt(1),
			BaseFee:     big.NewInt(1000000000000000000),
			GasLimit:    uint64(1000000000000000),
			Time:        big.NewInt(1488928920),
			GetHash: func(num uint64) common.Hash {
				return common.BytesToHash(crypto.Keccak256([]byte(big.NewInt(int64(num)).String())))
			},
			CanTransfer: canTransfer,
			Transfer:    func(db StateDB, sender, recipient common.Address, amount *big.Int) {},
		}

		memDBManager = database.NewMemoryDBManager()
		statedb      = initStateDB(memDBManager)

		env            = NewEVM(ctx, statedb, params.TestChainConfig, &Config{})
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.Origin = common.HexToAddress("0x9d19bb4553940f422104b1d0c8e5704c5aab63c9")
	env.callGasTemp = uint64(100000000000)
	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	evmInterpreter.returnData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mem.Resize(64)
	mem.Set(uint64(1), uint64(20), common.Hex2Bytes("000164865d7db79197021449b4a6aa650193b09e"))

	// make contract
	senderAddress := common.HexToAddress("0x91fb186da8f327f999782d1ae1ceacbd4fbbf146")
	payerAddress := common.HexToAddress("0x18f30de96ce789fe778b9a5f420f6fdbbd9b34d8")

	caller := types.NewAccountRefWithFeePayer(senderAddress, payerAddress)
	object := types.NewAccountRefWithFeePayer(payerAddress, senderAddress)
	contract := NewContract(caller, object, big.NewInt(0), uint64(1000))
	contract.Input = senderAddress.Bytes()
	contract.Gas = uint64(1000)
	contract.Code = common.Hex2Bytes("60ca60205260005b612710811015630000004557602051506020515060205150602051506020515060205150602051506020515060205150602051506001016300000007565b00")

	// convert args
	byteArgs := make([][]byte, len(args))
	for i, arg := range args {
		byteArgs[i] = common.Hex2Bytes(arg)
	}
	pc := uint64(0)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		for _, arg := range byteArgs {
			a := new(big.Int).SetBytes(arg)
			stack.push(a)
		}
		op(&pc, env, contract, mem, stack)
		if stack.len() > 0 {
			stack.pop()
		}
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func BenchmarkOpAdd64(b *testing.B) {
	x := "ffffffff"
	y := "fd37f3e2bba2c4f"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpAdd128(b *testing.B) {
	x := "ffffffffffffffff"
	y := "f5470b43c6549b016288e9a65629687"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpAdd256(b *testing.B) {
	x := "0802431afcbce1fc194c9eaa417b2fb67dc75a95db0bc7ec6b1c8af11df6a1da9"
	y := "a1f5aac137876480252e5dcac62c354ec0d42b76b0642b6181ed099849ea1d57"

	opBenchmark(b, opAdd, x, y)
}

func BenchmarkOpSub64(b *testing.B) {
	x := "51022b6317003a9d"
	y := "a20456c62e00753a"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpSub128(b *testing.B) {
	x := "4dde30faaacdc14d00327aac314e915d"
	y := "9bbc61f5559b829a0064f558629d22ba"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpSub256(b *testing.B) {
	x := "4bfcd8bb2ac462735b48a17580690283980aa2d679f091c64364594df113ea37"
	y := "97f9b1765588c4e6b69142eb00d20507301545acf3e1238c86c8b29be227d46e"

	opBenchmark(b, opSub, x, y)
}

func BenchmarkOpMul(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMul, x, y)
}

func BenchmarkOpDiv256(b *testing.B) {
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpDiv128(b *testing.B) {
	x := "fdedc7f10142ff97"
	y := "fbdfda0e2ce356173d1993d5f70a2b11"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpDiv64(b *testing.B) {
	x := "fcb34eb3"
	y := "f97180878e839129"
	opBenchmark(b, opDiv, x, y)
}

func BenchmarkOpSdiv(b *testing.B) {
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"

	opBenchmark(b, opSdiv, x, y)
}

func BenchmarkOpMod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMod, x, y)
}

func BenchmarkOpSmod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSmod, x, y)
}

func BenchmarkOpExp(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opExp, x, y)
}

func BenchmarkOpSignExtend(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSignExtend, x, y)
}

func BenchmarkOpLt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opLt, x, y)
}

func BenchmarkOpGt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opGt, x, y)
}

func BenchmarkOpSlt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSlt, x, y)
}

func BenchmarkOpSgt(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opSgt, x, y)
}

func BenchmarkOpEq(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opEq, x, y)
}
func BenchmarkOpEq2(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201fffffffe"
	opBenchmark(b, opEq, x, y)
}
func BenchmarkOpAnd(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opAnd, x, y)
}

func BenchmarkOpOr(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opOr, x, y)
}

func BenchmarkOpXor(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opXor, x, y)
}

func BenchmarkOpByte(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opByte, x, y)
}

func BenchmarkOpAddmod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	z := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opAddmod, x, y, z)
}

func BenchmarkOpMulmod(b *testing.B) {
	x := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	y := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	z := "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"

	opBenchmark(b, opMulmod, x, y, z)
}

func BenchmarkOpSHL(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSHL, x, y)
}
func BenchmarkOpSHR(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSHR, x, y)
}
func BenchmarkOpSAR(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSAR, x, y)
}
func BenchmarkOpIsZero(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	opBenchmark(b, opIszero, x)
}

func TestOpMstore(t *testing.T) {
	var (
		env            = NewEVM(Context{}, nil, params.TestChainConfig, &Config{})
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	mem.Resize(64)
	pc := uint64(0)
	v := "abcdef00000000000000abba000000000deaf000000c0de00100000000133700"
	stack.pushN(new(big.Int).SetBytes(common.Hex2Bytes(v)), big.NewInt(0))
	opMstore(&pc, env, nil, mem, stack)
	if got := common.Bytes2Hex(mem.GetCopy(0, 32)); got != v {
		t.Fatalf("Mstore fail, got %v, expected %v", got, v)
	}
	stack.pushN(big.NewInt(0x1), big.NewInt(0))
	opMstore(&pc, env, nil, mem, stack)
	if common.Bytes2Hex(mem.GetCopy(0, 32)) != "0000000000000000000000000000000000000000000000000000000000000001" {
		t.Fatalf("Mstore failed to overwrite previous value")
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func BenchmarkOpMstore(bench *testing.B) {
	var (
		env            = NewEVM(Context{}, nil, params.TestChainConfig, &Config{})
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	mem.Resize(64)
	pc := uint64(0)
	memStart := big.NewInt(0)
	value := big.NewInt(0x1337)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		stack.pushN(value, memStart)
		opMstore(&pc, env, nil, mem, stack)
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func BenchmarkOpSHA3(bench *testing.B) {
	var (
		env            = NewEVM(Context{}, nil, params.TestChainConfig, &Config{})
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)
	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()
	mem.Resize(32)
	pc := uint64(0)
	start := big.NewInt(0)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		stack.pushN(big.NewInt(32), start)
		opSha3(&pc, env, nil, mem, stack)
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func TestCreate2Addreses(t *testing.T) {
	type testcase struct {
		origin   string
		salt     string
		code     string
		expected string
	}

	for i, tt := range []testcase{
		{
			origin:   "0x0000000000000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "0x4d1a2e2bb4f88f0250f26ffff098b0b30b26bf38",
		},
		{
			origin:   "0xdeadbeef00000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "0xB928f69Bb1D91Cd65274e3c79d8986362984fDA3",
		},
		{
			origin:   "0xdeadbeef00000000000000000000000000000000",
			salt:     "0xfeed000000000000000000000000000000000000",
			code:     "0x00",
			expected: "0xD04116cDd17beBE565EB2422F2497E06cC1C9833",
		},
		{
			origin:   "0x0000000000000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0xdeadbeef",
			expected: "0x70f2b2914A2a4b783FaEFb75f459A580616Fcb5e",
		},
		{
			origin:   "0x00000000000000000000000000000000deadbeef",
			salt:     "0xcafebabe",
			code:     "0xdeadbeef",
			expected: "0x60f3f640a8508fC6a86d45DF051962668E1e8AC7",
		},
		{
			origin:   "0x00000000000000000000000000000000deadbeef",
			salt:     "0xcafebabe",
			code:     "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			expected: "0x1d8bfDC5D46DC4f61D6b6115972536eBE6A8854C",
		},
		{
			origin:   "0x0000000000000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x",
			expected: "0xE33C0C7F7df4809055C3ebA6c09CFe4BaF1BD9e0",
		},
	} {

		origin := common.BytesToAddress(common.FromHex(tt.origin))
		salt := common.BytesToHash(common.FromHex(tt.salt))
		code := common.FromHex(tt.code)
		codeHash := crypto.Keccak256(code)
		address := crypto.CreateAddress2(origin, salt, codeHash)
		/*
			stack          := newstack()
			// salt, but we don't need that for this test
			stack.push(big.NewInt(int64(len(code)))) //size
			stack.push(big.NewInt(0)) // memstart
			stack.push(big.NewInt(0)) // value
			gas, _ := gasCreate2(params.GasTable{}, nil, nil, stack, nil, 0)
			fmt.Printf("Example %d\n* address `0x%x`\n* salt `0x%x`\n* init_code `0x%x`\n* gas (assuming no mem expansion): `%v`\n* result: `%s`\n\n", i,origin, salt, code, gas, address.String())
		*/
		expected := common.BytesToAddress(common.FromHex(tt.expected))
		if !bytes.Equal(expected.Bytes(), address.Bytes()) {
			t.Errorf("test %d: expected %s, got %s", i, expected.String(), address.String())
		}

	}
}

func BenchmarkOpStop(b *testing.B) {
	opBenchmark(b, opStop)
}

func BenchmarkOpNot(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	opBenchmark(b, opNot, x)
}

func BenchmarkOpAddress(b *testing.B) {
	opBenchmark(b, opAddress)
}

func BenchmarkOpBalance(b *testing.B) {
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d8"
	opBenchmark(b, opBalance, addr)
}

func BenchmarkOpOrigin(b *testing.B) {
	opBenchmark(b, opOrigin)
}

func BenchmarkOpCaller(b *testing.B) {
	opBenchmark(b, opCaller)
}

func BenchmarkOpCallValue(b *testing.B) {
	opBenchmark(b, opCallValue)
}

func BenchmarkOpCallDataLoad(b *testing.B) {
	x := "0000000000000000000000000000000000000000000000000000000000000001" // 1
	opBenchmark(b, opCallDataLoad, x)
}

func BenchmarkOpCallDataSize(b *testing.B) {
	opBenchmark(b, opCallDataSize)
}

func BenchmarkOpCallDataCopy(b *testing.B) {
	length := "0000000000000000000000000000000000000000000000000000000000000013"     // 19
	dataOffset := "0000000000000000000000000000000000000000000000000000000000000001" // 1
	memOffset := "0000000000000000000000000000000000000000000000000000000000000001"  // 1
	opBenchmark(b, opCallDataCopy, length, dataOffset, memOffset)
}

func BenchmarkOpReturnDataSize(b *testing.B) {
	opBenchmark(b, opReturnDataSize)
}

func BenchmarkOpReturnDataCopy(b *testing.B) {
	length := "0000000000000000000000000000000000000000000000000000000000000009"     // 9
	dataOffset := "0000000000000000000000000000000000000000000000000000000000000001" // 1
	memOffset := "0000000000000000000000000000000000000000000000000000000000000001"  // 1
	opBenchmark(b, opReturnDataCopy, length, dataOffset, memOffset)
}

func BenchmarkOpExtCodeSize(b *testing.B) {
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d8"
	opBenchmark(b, opExtCodeSize, addr)
}

func BenchmarkOpCodeSize(b *testing.B) {
	opBenchmark(b, opCodeSize)
}

func BenchmarkOpCodeCopy(b *testing.B) {
	length := "000000000000000000000000000000000000000000000000000000000000003c"     // 60
	dataOffset := "0000000000000000000000000000000000000000000000000000000000000001" // 1
	memOffset := "0000000000000000000000000000000000000000000000000000000000000001"  // 1
	opBenchmark(b, opCodeCopy, length, dataOffset, memOffset)
}

func BenchmarkOpExtCodeCopy(b *testing.B) {
	length := "000000000000000000000000000000000000000000000000000000000000003c"     // 60
	codeOffset := "0000000000000000000000000000000000000000000000000000000000000001" // 1
	memOffset := "0000000000000000000000000000000000000000000000000000000000000001"  // 1
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d8"
	opBenchmark(b, opExtCodeCopy, length, codeOffset, memOffset, addr)
}

func BenchmarkOpExtCodeHash(b *testing.B) {
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d8"
	opBenchmark(b, opExtCodeHash, addr)
}

func BenchmarkOpGasprice(b *testing.B) {
	opBenchmark(b, opGasprice)
}

func BenchmarkOpBlockhash(b *testing.B) {
	num := "00000000000000000000000000000000000000000000000000000000000003E8" // 1000
	opBenchmark(b, opBlockhash, num)
}

func BenchmarkOpCoinbase(b *testing.B) {
	opBenchmark(b, opCoinbase)
}

func BenchmarkOpTimestamp(b *testing.B) {
	opBenchmark(b, opTimestamp)
}

func BenchmarkOpNumber(b *testing.B) {
	opBenchmark(b, opNumber)
}

func BenchmarkOpDifficulty(b *testing.B) {
	opBenchmark(b, opDifficulty)
}

func BenchmarkOpGasLimit(b *testing.B) {
	opBenchmark(b, opGasLimit)
}

func BenchmarkOpPop(b *testing.B) {
	x := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	opBenchmark(b, opPop, x)
}

func BenchmarkOpMload(b *testing.B) {
	offset := "0000000000000000000000000000000000000000000000000000000000000001"
	opBenchmark(b, opMload, offset)
}

func BenchmarkOpMstore8(b *testing.B) {
	val := "7FFFFFFFFFFFFFFF"
	off := "0000000000000000000000000000000000000000000000000000000000000001"
	opBenchmark(b, opMstore8, val, off)
}

func BenchmarkOpSload(b *testing.B) {
	loc := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	opBenchmark(b, opSload, loc)
}

func BenchmarkOpSstore(bench *testing.B) {
	var (
		memDBManager = database.NewMemoryDBManager()
		statedb      = initStateDB(memDBManager)

		env            = NewEVM(Context{}, statedb, params.TestChainConfig, &Config{})
		stack          = newstack()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	evmInterpreter.intPool = poolOfIntPools.get()

	// make contract
	senderAddress := common.HexToAddress("0x91fb186da8f327f999782d1ae1ceacbd4fbbf146")
	payerAddress := common.HexToAddress("0x18f30de96ce789fe778b9a5f420f6fdbbd9b34d8")

	caller := types.NewAccountRefWithFeePayer(senderAddress, payerAddress)
	object := types.NewAccountRefWithFeePayer(payerAddress, senderAddress)
	contract := NewContract(caller, object, big.NewInt(0), uint64(1000))

	// convert args
	byteArgs := make([][]byte, 2)
	byteArgs[0] = common.Hex2Bytes("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd") // val
	byteArgs[1] = common.Hex2Bytes("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")  // loc

	pc := uint64(0)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		stack.push(new(big.Int).SetBytes(byteArgs[0]))
		stack.push(new(big.Int).SetBytes(append(byteArgs[1], byte(i))))
		opSstore(&pc, env, contract, nil, stack)
		if stack.len() > 0 {
			stack.pop()
		}
	}
	poolOfIntPools.put(evmInterpreter.intPool)
}

func BenchmarkOpJump(b *testing.B) {
	pos := "0000000000000000000000000000000000000000000000000000000000000045"
	opBenchmark(b, opJump, pos)
}

func BenchmarkOpJumpi(b *testing.B) {
	cond := "0000000000000000000000000000000000000000000000000000000000000001"
	pos := "0000000000000000000000000000000000000000000000000000000000000045"
	opBenchmark(b, opJumpi, cond, pos)
}

func BenchmarkOpJumpdest(b *testing.B) {
	opBenchmark(b, opJumpdest)
}

func BenchmarkOpPc(b *testing.B) {
	opBenchmark(b, opPc)
}

func BenchmarkOpMsize(b *testing.B) {
	opBenchmark(b, opMsize)
}

func BenchmarkOpGas(b *testing.B) {
	opBenchmark(b, opGas)
}

func BenchmarkOpCreate(b *testing.B) {
	size := "0000000000000014"
	offset := "0000000000000001"
	value := "7FFFFFFFFFFFFFFF"
	opBenchmark(b, opCreate, size, offset, value)
}

func BenchmarkOpCreate2(b *testing.B) {
	salt := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	size := "0000000000000014"
	offset := "0000000000000001"
	endowment := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	opBenchmark(b, opCreate2, salt, size, offset, endowment)
}

func BenchmarkOpCall(b *testing.B) {
	retSize := "0000000000000000000000000000000000000000000000000000000000000001"
	retOffset := "0000000000000000000000000000000000000000000000000000000000000001"
	inSize := "0000000000000014"
	inOffset := "0000000000000001"
	value := "0000000000000000000000000000000000000000000000000000000000000001"
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d9"
	gas := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	opBenchmark(b, opCall, retSize, retOffset, inSize, inOffset, value, addr, gas)
}

func BenchmarkOpCallCode(b *testing.B) {
	retSize := "0000000000000000000000000000000000000000000000000000000000000001"
	retOffset := "0000000000000000000000000000000000000000000000000000000000000001"
	inSize := "0000000000000014"
	inOffset := "0000000000000001"
	value := "0000000000000000000000000000000000000000000000000000000000000001"
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d9"
	gas := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	opBenchmark(b, opCallCode, retSize, retOffset, inSize, inOffset, value, addr, gas)
}

func BenchmarkOpDelegateCall(b *testing.B) {
	retSize := "0000000000000000000000000000000000000000000000000000000000000001"
	retOffset := "0000000000000000000000000000000000000000000000000000000000000001"
	inSize := "0000000000000014"
	inOffset := "0000000000000001"
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d9"
	gas := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	opBenchmark(b, opDelegateCall, retSize, retOffset, inSize, inOffset, addr, gas)
}

func BenchmarkOpStaticCall(b *testing.B) {
	retSize := "0000000000000000000000000000000000000000000000000000000000000001"
	retOffset := "0000000000000000000000000000000000000000000000000000000000000001"
	inSize := "0000000000000014"
	inOffset := "0000000000000001"
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d9"
	gas := "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	opBenchmark(b, opStaticCall, retSize, retOffset, inSize, inOffset, addr, gas)
}

func BenchmarkOpReturn(b *testing.B) {
	size := "0000000000000014"
	offset := "0000000000000001"
	opBenchmark(b, opReturn, size, offset)
}

func BenchmarkOpRevert(b *testing.B) {
	size := "0000000000000014"
	offset := "0000000000000001"
	opBenchmark(b, opRevert, size, offset)
}

func BenchmarkOpSuicide(b *testing.B) {
	addr := "18f30de96ce789fe778b9a5f420f6fdbbd9b34d8"
	opBenchmark(b, opSuicide, addr)
}

func BenchmarkOpPush1(b *testing.B) {
	opBenchmark(b, opPush1)
}

func BenchmarkOpPush2(b *testing.B) {
	opBenchmark(b, makePush(uint64(2), 2))
}

func BenchmarkOpPush3(b *testing.B) {
	opBenchmark(b, makePush(uint64(3), 3))
}

func BenchmarkOpPush4(b *testing.B) {
	opBenchmark(b, makePush(uint64(4), 4))
}

func BenchmarkOpPush5(b *testing.B) {
	opBenchmark(b, makePush(uint64(5), 5))
}

func BenchmarkOpPush6(b *testing.B) {
	opBenchmark(b, makePush(uint64(6), 6))
}

func BenchmarkOpPush7(b *testing.B) {
	opBenchmark(b, makePush(uint64(7), 7))
}

func BenchmarkOpPush8(b *testing.B) {
	opBenchmark(b, makePush(uint64(8), 8))
}

func BenchmarkOpPush9(b *testing.B) {
	opBenchmark(b, makePush(uint64(9), 9))
}

func BenchmarkOpPush10(b *testing.B) {
	opBenchmark(b, makePush(uint64(10), 10))
}

func BenchmarkOpPush11(b *testing.B) {
	opBenchmark(b, makePush(uint64(11), 11))
}

func BenchmarkOpPush12(b *testing.B) {
	opBenchmark(b, makePush(uint64(12), 12))
}

func BenchmarkOpPush13(b *testing.B) {
	opBenchmark(b, makePush(uint64(13), 13))
}

func BenchmarkOpPush14(b *testing.B) {
	opBenchmark(b, makePush(uint64(14), 14))
}

func BenchmarkOpPush15(b *testing.B) {
	opBenchmark(b, makePush(uint64(15), 15))
}

func BenchmarkOpPush16(b *testing.B) {
	opBenchmark(b, makePush(uint64(16), 16))
}

func BenchmarkOpPush17(b *testing.B) {
	opBenchmark(b, makePush(uint64(17), 17))
}

func BenchmarkOpPush18(b *testing.B) {
	opBenchmark(b, makePush(uint64(18), 18))
}

func BenchmarkOpPush19(b *testing.B) {
	opBenchmark(b, makePush(uint64(19), 19))
}

func BenchmarkOpPush20(b *testing.B) {
	opBenchmark(b, makePush(uint64(20), 20))
}

func BenchmarkOpPush21(b *testing.B) {
	opBenchmark(b, makePush(uint64(21), 21))
}

func BenchmarkOpPush22(b *testing.B) {
	opBenchmark(b, makePush(uint64(22), 22))
}

func BenchmarkOpPush23(b *testing.B) {
	opBenchmark(b, makePush(uint64(23), 23))
}

func BenchmarkOpPush24(b *testing.B) {
	opBenchmark(b, makePush(uint64(24), 24))
}

func BenchmarkOpPush25(b *testing.B) {
	opBenchmark(b, makePush(uint64(25), 25))
}

func BenchmarkOpPush26(b *testing.B) {
	opBenchmark(b, makePush(uint64(26), 26))
}

func BenchmarkOpPush27(b *testing.B) {
	opBenchmark(b, makePush(uint64(27), 27))
}

func BenchmarkOpPush28(b *testing.B) {
	opBenchmark(b, makePush(uint64(28), 28))
}

func BenchmarkOpPush29(b *testing.B) {
	opBenchmark(b, makePush(uint64(29), 29))
}

func BenchmarkOpPush30(b *testing.B) {
	opBenchmark(b, makePush(uint64(30), 30))
}

func BenchmarkOpPush31(b *testing.B) {
	opBenchmark(b, makePush(uint64(31), 31))
}

func BenchmarkOpPush32(b *testing.B) {
	opBenchmark(b, makePush(uint64(32), 32))
}

func BenchmarkOpDup1(b *testing.B) {
	size := 1
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup2(b *testing.B) {
	size := 2
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup3(b *testing.B) {
	size := 3
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup4(b *testing.B) {
	size := 4
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup5(b *testing.B) {
	size := 5
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup6(b *testing.B) {
	size := 6
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup7(b *testing.B) {
	size := 7
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup8(b *testing.B) {
	size := 8
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup9(b *testing.B) {
	size := 9
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup10(b *testing.B) {
	size := 10
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup11(b *testing.B) {
	size := 11
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup12(b *testing.B) {
	size := 12
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup13(b *testing.B) {
	size := 13
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup14(b *testing.B) {
	size := 14
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup15(b *testing.B) {
	size := 15
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpDup16(b *testing.B) {
	size := 16
	stacks := genStacksForDup(size)
	opBenchmark(b, makeDup(int64(size)), stacks...)
}

func BenchmarkOpSwap1(b *testing.B) {
	size := 1
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap2(b *testing.B) {
	size := 2
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap3(b *testing.B) {
	size := 3
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap4(b *testing.B) {
	size := 4
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap5(b *testing.B) {
	size := 5
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap6(b *testing.B) {
	size := 6
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap7(b *testing.B) {
	size := 7
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap8(b *testing.B) {
	size := 8
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap9(b *testing.B) {
	size := 9
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap10(b *testing.B) {
	size := 10
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap11(b *testing.B) {
	size := 11
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap12(b *testing.B) {
	size := 12
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap13(b *testing.B) {
	size := 13
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap14(b *testing.B) {
	size := 14
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap15(b *testing.B) {
	size := 15
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpSwap16(b *testing.B) {
	size := 16
	stacks := genStacksForSwap(size)
	opBenchmark(b, makeSwap(int64(size)), stacks...)
}

func BenchmarkOpLog0(b *testing.B) {
	size := 0
	stacks := genStacksForLog(size)
	opBenchmark(b, makeLog(size), stacks...)
}

func BenchmarkOpLog1(b *testing.B) {
	size := 1
	stacks := genStacksForLog(size)
	opBenchmark(b, makeLog(size), stacks...)
}

func BenchmarkOpLog2(b *testing.B) {
	size := 2
	stacks := genStacksForLog(size)
	opBenchmark(b, makeLog(size), stacks...)
}

func BenchmarkOpLog3(b *testing.B) {
	size := 3
	stacks := genStacksForLog(size)
	opBenchmark(b, makeLog(size), stacks...)
}

func BenchmarkOpLog4(b *testing.B) {
	size := 4
	stacks := genStacksForLog(size)
	opBenchmark(b, makeLog(size), stacks...)
}

func BenchmarkOpChainID(b *testing.B) {
	opBenchmark(b, opChainID)
}

func BenchmarkOpSelfBalance(b *testing.B) {
	opBenchmark(b, opSelfBalance)
}

func BenchmarkOpBaseFee(b *testing.B) {
	opBenchmark(b, opBaseFee)
}

func genStacksForDup(size int) []string {
	stacks := make([]string, size)
	return fillStacks(stacks, size)
}

func genStacksForSwap(size int) []string {
	stacks := make([]string, size+1)
	fillStacks(stacks, size)
	stacks[len(stacks)-1] = "ABCDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff"
	return stacks
}

func genStacksForLog(size int) []string {
	stacks := make([]string, size+2)
	fillStacks(stacks, size)
	stacks[len(stacks)-2] = "0000000000000000000000000000000000000000000000000000000000000014" // 20
	stacks[len(stacks)-1] = "0000000000000000000000000000000000000000000000000000000000000001" // 1
	return stacks
}

func fillStacks(stacks []string, n int) []string {
	for i := 0; i < n; i++ {
		stacks[i] = "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	}
	return stacks
}
