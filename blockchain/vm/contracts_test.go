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
// This file is derived from core/vm/contracts_test.go (2018/06/04).
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
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// precompiledTest defines the input/output pairs for precompiled contract tests.
type precompiledTest struct {
	Input, Expected string
	Gas             uint64
	Name            string
	NoBenchmark     bool // Benchmark primarily the worst-cases
}

// precompiledFailureTest defines the input/error pairs for precompiled
// contract failure tests.
type precompiledFailureTest struct {
	Input         string
	ExpectedError string
	Name          string
}

// EIP-152 test vectors
var blake2FMalformedInputTests = []precompiledFailureTest{
	{
		Input:         "",
		ExpectedError: errBlake2FInvalidInputLength.Error(),
		Name:          "vector 0: empty input",
	},
	{
		Input:         "00000c48c9bdf267e6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa5d182e6ad7f520e511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b61626300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000001",
		ExpectedError: errBlake2FInvalidInputLength.Error(),
		Name:          "vector 1: less than 213 bytes input",
	},
	{
		Input:         "000000000c48c9bdf267e6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa5d182e6ad7f520e511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b61626300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000001",
		ExpectedError: errBlake2FInvalidInputLength.Error(),
		Name:          "vector 2: more than 213 bytes input",
	},
	{
		Input:         "0000000c48c9bdf267e6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa5d182e6ad7f520e511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b61626300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000002",
		ExpectedError: errBlake2FInvalidFinalFlag.Error(),
		Name:          "vector 3: malformed final block indicator flag",
	},
}

// This function prepares background environment for running vmLog, feePayer, validateSender tests.
// It generates contract, evm, EOA test object.
func prepare(reqGas uint64) (*Contract, *EVM, error) {
	// Generate Contract
	contract := NewContract(types.NewAccountRefWithFeePayer(common.HexToAddress("1337"), common.HexToAddress("133773")),
		nil, new(big.Int), reqGas)

	// Generate EVM
	stateDb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil)
	txhash := common.HexToHash("0xc6a37e155d3fa480faea012a68ad35fd53c8cc3cd8263a434c697755985a6577")
	stateDb.Prepare(txhash, common.Hash{}, 0)
	evm := NewEVM(Context{BlockNumber: big.NewInt(0)}, stateDb, &params.ChainConfig{IstanbulCompatibleBlock: big.NewInt(0)}, &Config{})

	// Only stdout logging is tested to avoid file handling. It is used at vmLog test.
	params.VMLogTarget = params.VMLogToStdout

	// Generate EOA. It is used at validateSender test.
	k, err := crypto.HexToECDSA("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	accKey := accountkey.NewAccountKeyPublicWithValue(&k.PublicKey)
	stateDb.CreateEOA(common.HexToAddress("0x123456789"), false, accKey)

	return contract, evm, err
}

func testPrecompiled(addr string, test precompiledTest, t *testing.T) {
	p := PrecompiledContractsIstanbul[common.HexToAddress(addr)]
	in := common.Hex2Bytes(test.Input)
	reqGas, _ := p.GetRequiredGasAndComputationCost(in)

	contract, evm, err := prepare(reqGas)
	require.NoError(t, err)

	t.Run(fmt.Sprintf("%s-Gas=%d", test.Name, reqGas), func(t *testing.T) {
		if res, _, err := RunPrecompiledContract(p, in, contract, evm); err != nil {
			t.Error(err)
		} else if common.Bytes2Hex(res) != test.Expected {
			t.Errorf("Expected %v, got %v", test.Expected, common.Bytes2Hex(res))
		}
	})
}

func testPrecompiledOOG(addr string, test precompiledTest, t *testing.T) {
	p := PrecompiledContractsIstanbul[common.HexToAddress(addr)]
	in := common.Hex2Bytes(test.Input)
	reqGas, _ := p.GetRequiredGasAndComputationCost(in)
	reqGas -= 1

	contract, evm, _ := prepare(reqGas)

	t.Run(fmt.Sprintf("%s-Gas=%d", test.Name, reqGas), func(t *testing.T) {
		_, _, err := RunPrecompiledContract(p, in, contract, evm)
		if err.Error() != "out of gas" {
			t.Errorf("Expected error [out of gas], got [%v]", err)
		}
		// Verify that the precompile did not touch the input buffer
		exp := common.Hex2Bytes(test.Input)
		if !bytes.Equal(in, exp) {
			t.Errorf("Precompiled %v modified input data", addr)
		}
	})
}

func testPrecompiledFailure(addr string, test precompiledFailureTest, t *testing.T) {
	p := PrecompiledContractsIstanbul[common.HexToAddress(addr)]
	in := common.Hex2Bytes(test.Input)
	reqGas, _ := p.GetRequiredGasAndComputationCost(in)

	contract, evm, _ := prepare(reqGas)

	t.Run(test.Name, func(t *testing.T) {
		_, _, err := RunPrecompiledContract(p, in, contract, evm)
		if err.Error() != test.ExpectedError {
			t.Errorf("Expected error [%v], got [%v]", test.ExpectedError, err)
		}
		// Verify that the precompile did not touch the input buffer
		exp := common.Hex2Bytes(test.Input)
		if !bytes.Equal(in, exp) {
			t.Errorf("Precompiled %v modified input data", addr)
		}
	})
}

func benchmarkPrecompiled(addr string, test precompiledTest, bench *testing.B) {
	if test.NoBenchmark {
		return
	}
	p := PrecompiledContractsIstanbul[common.HexToAddress(addr)]
	in := common.Hex2Bytes(test.Input)
	reqGas, _ := p.GetRequiredGasAndComputationCost(in)

	contract, evm, _ := prepare(reqGas)

	var (
		res  []byte
		err  error
		data = make([]byte, len(in))
	)

	bench.Run(fmt.Sprintf("%s-Gas=%d", test.Name, contract.Gas), func(bench *testing.B) {
		bench.ResetTimer()
		for i := 0; i < bench.N; i++ {
			contract.Gas = reqGas
			copy(data, in)
			res, _, err = RunPrecompiledContract(p, data, contract, evm)
		}
		bench.StopTimer()
		//Check if it is correct
		if err != nil {
			bench.Error(err)
			return
		}
		if common.Bytes2Hex(res) != test.Expected {
			bench.Error(fmt.Sprintf("Expected %v, got %v", test.Expected, common.Bytes2Hex(res)))
			return
		}
	})
}

// Tests the sample inputs of the ecrecover
func TestPrecompiledEcrecover(t *testing.T)      { testJson("ecRecover", "01", t) }
func BenchmarkPrecompiledEcrecover(b *testing.B) { benchJson("ecRecover", "01", b) }

// Benchmarks the sample inputs from the SHA256 precompile.
func BenchmarkPrecompiledSha256(b *testing.B) { benchJson("sha256", "02", b) }

// Benchmarks the sample inputs from the RIPEMD precompile.
func BenchmarkPrecompiledRipeMD(b *testing.B) { benchJson("ripeMD", "03", b) }

// Benchmarks the sample inputs from the identity precompile.
func BenchmarkPrecompiledIdentity(b *testing.B) { benchJson("identity", "04", b) }

// Tests the sample inputs from the ModExp EIP 198.
func TestPrecompiledModExp(t *testing.T)      { testJson("modexp", "05", t) }
func BenchmarkPrecompiledModExp(b *testing.B) { benchJson("modexp", "05", b) }

// Tests the sample inputs from the elliptic curve addition EIP 213.
func TestPrecompiledBn256Add(t *testing.T)      { testJson("bn256Add", "06", t) }
func BenchmarkPrecompiledBn256Add(b *testing.B) { benchJson("bn256Add", "06", b) }

// Tests the sample inputs from the elliptic curve scalar multiplication EIP 213.
func TestPrecompiledBn256ScalarMul(t *testing.T)      { testJson("bn256ScalarMul", "07", t) }
func BenchmarkPrecompiledBn256ScalarMul(b *testing.B) { benchJson("bn256ScalarMul", "07", b) }

// Tests the sample inputs from the elliptic curve pairing check EIP 197.
func TestPrecompiledBn256Pairing(t *testing.T)      { testJson("bn256Pairing", "08", t) }
func BenchmarkPrecompiledBn256Pairing(b *testing.B) { benchJson("bn256Pairing", "08", b) }

func TestPrecompiledBlake2F(t *testing.T)              { testJson("blake2F", "09", t) }
func BenchmarkPrecompiledBlake2F(b *testing.B)         { benchJson("blake2F", "09", b) }
func TestPrecompileBlake2FMalformedInput(t *testing.T) { testJsonFail("blake2F", "09", t) }

// Tests the sample inputs of the vmLog
func TestPrecompiledVmLog(t *testing.T)      { testJson("vmLog", "3fd", t) }
func BenchmarkPrecompiledVmLog(b *testing.B) { benchJson("vmLog", "3fd", b) }

// Tests the sample inputs of the feePayer
func TestFeePayerContract(t *testing.T)         { testJson("feePayer", "3fe", t) }
func BenchmarkPrecompiledFeePayer(b *testing.B) { benchJson("feePayer", "3fe", b) }

// Tests the sample inputs of the validateSender
func TestValidateSenderContract(t *testing.T)         { testJson("validateSender", "3ff", t) }
func BenchmarkPrecompiledValidateSender(b *testing.B) { benchJson("validateSender", "3ff", b) }

// Tests OOG (out-of-gas) of modExp
func TestPrecompiledModExpOOG(t *testing.T) {
	modexpTests, err := loadJson("modexp")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range modexpTests {
		testPrecompiledOOG("05", test, t)
	}
}

func loadJson(name string) ([]precompiledTest, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/precompiles/%v.json", name))
	if err != nil {
		return nil, err
	}
	var testcases []precompiledTest
	err = json.Unmarshal(data, &testcases)
	return testcases, err
}

func loadJsonFail(name string) ([]precompiledFailureTest, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/precompiles/fail-%v.json", name))
	if err != nil {
		return nil, err
	}
	var testcases []precompiledFailureTest
	err = json.Unmarshal(data, &testcases)
	return testcases, err
}

func testJson(name, addr string, t *testing.T) {
	tests, err := loadJson(name)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		testPrecompiled(addr, test, t)
	}
}

func testJsonFail(name, addr string, t *testing.T) {
	tests, err := loadJsonFail(name)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		testPrecompiledFailure(addr, test, t)
	}
}

func benchJson(name, addr string, b *testing.B) {
	tests, err := loadJson(name)
	if err != nil {
		b.Fatal(err)
	}
	for _, test := range tests {
		benchmarkPrecompiled(addr, test, b)
	}
}

// TestEVM_CVE_2021_39137 tests an EVM vulnerability described in https://nvd.nist.gov/vuln/detail/CVE-2021-39137.
// The vulnerable EVM bytecode exploited in Ethereum is used as a test code.
// Test code reference: https://etherscan.io/tx/0x1cb6fb36633d270edefc04d048145b4298e67b8aa82a9e5ec4aa1435dd770ce4.
func TestEVM_CVE_2021_39137(t *testing.T) {
	fromAddr := common.HexToAddress("0x1a02a619e51cc5f8a2a61d2a60f6c80476ee8ead")
	contractAddr := common.HexToAddress("0x8eae784e072e961f76948a785b62c9a950fb17ae")

	testCases := []struct {
		name           string
		expectedResult []byte
		testCode       []byte
	}{
		{
			"staticCall test",
			contractAddr.Bytes(),
			hexutil.MustDecode("0x3034526020600760203460045afa602034343e604034f3"),
			/*
				// Pseudo code of the decompiled testCode
				memory[0:0x20] = address(this); // put contract address with padding left zero into the memory
				staticCall(gas, 0x04, 0x0, 0x20, 0x07, 0x20);  // operands: gas, to, in offset, in size, out offset, out size
				memory[0:0x20] = returnDataCopy(); // put the returned data from staticCall into the memory
				return memory[0:0x40];

				// Disassembly
				0000    30  ADDRESS
				0001    34  CALLVALUE
				0002    52  MSTORE
				0003    60  PUSH1 0x20
				0005    60  PUSH1 0x07
				0007    60  PUSH1 0x20
				0009    34  CALLVALUE
				000A    60  PUSH1 0x04
				000C    5A  GAS
				000D    FA  STATICCALL
				000E    60  PUSH1 0x20
				0010    34  CALLVALUE
				0011    34  CALLVALUE
				0012    3E  RETURNDATACOPY
				0013    60  PUSH1 0x40
				0015    34  CALLVALUE
				0016    F3  *RETURN
			*/
		},
		{
			"call test",
			contractAddr.Bytes(),
			hexutil.MustDecode("0x30345260206007602060003460045af1602034343e604034f3"),
		},
		{
			"callCode test",
			contractAddr.Bytes(),
			hexutil.MustDecode("0x30345260206007602060003460045af2602034343e604034f3"),
		},
		{
			"delegateCall test",
			contractAddr.Bytes(),
			hexutil.MustDecode("0x3034526020600760203460045af4602034343e604034f3"),
		},
	}

	gasLimit := uint64(99999999)
	tracer := NewStructLogger(nil)
	vmctx := Context{
		CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
	}
	stateDb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil)

	for _, tc := range testCases {
		stateDb.SetCode(contractAddr, tc.testCode)

		evm := NewEVM(vmctx, stateDb, params.TestChainConfig, &Config{Debug: true, Tracer: tracer})
		ret, _, err := evm.Call(AccountRef(fromAddr), contractAddr, nil, gasLimit, new(big.Int))
		if err != nil {
			t.Fatal(err)
		}

		if testing.Verbose() {
			buf := new(bytes.Buffer)
			WriteTrace(buf, tracer.StructLogs())
			if buf.Len() == 0 {
				t.Log("no EVM operation logs generated")
			} else {
				t.Log("EVM operation log:\n" + buf.String())
			}
			t.Logf("EVM output: 0x%x", tracer.Output())
			t.Logf("EVM error: %v", tracer.Error())
		}

		assert.Equal(t, tc.expectedResult, ret[12:32])
	}
}
