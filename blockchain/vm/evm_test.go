package vm

import (
	"errors"
	"math"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

// Block data
var (
	Block4 = big.NewInt(4) // Block number before the istanbulHF
	Block5 = big.NewInt(5) // Block number after the istanbulHF(also istanbulHF Block number)
)

type TestData struct {
	addr  string
	input []byte
	// Test condition
	isDeployedAfterIstanbulHF bool     // contract deployment time
	block                     *big.Int // block when execution is done
	// Expect field
	expectGas uint64
	expectRet string
	expectErr error
}

func runPrecompiledContractTestWithHFCondition(t *testing.T, config *params.ChainConfig, testData []TestData) {
	for _, tc := range testData {
		// Make StateDB
		callerAddr := common.BytesToAddress([]byte("contract"))
		statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil)
		statedb.CreateSmartContractAccount(callerAddr, params.CodeFormatEVM, params.Rules{IsIstanbul: tc.isDeployedAfterIstanbulHF})

		// Make EVM environment
		vmctx := Context{
			CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
			Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
			BlockNumber: tc.block,
		}
		vmenv := NewEVM(vmctx, statedb, config, &Config{})

		// run
		ret, gas, err := vmenv.Call(AccountRef(callerAddr), common.HexToAddress(tc.addr), tc.input, math.MaxUint64, new(big.Int))

		// compare with expected data
		assert.Equal(t, tc.expectRet, common.Bytes2Hex(ret))
		assert.Equal(t, tc.expectGas, math.MaxUint64-gas)
		assert.Equal(t, tc.expectErr, err)
	}
}

func TestPrecompiledContractAddressMapping(t *testing.T) {
	var (
		// Test Input
		vmLogInput    = []byte("Hello")
		feePayerInput = []byte("")
		validateInput = common.Hex2Bytes("00000000000000000000000000000001234567890000000000000000000000000000000000000000000000000000000000000000db38e383310a432f1cec0bfbd765fe67c9932b02295b1f2062a785eec288b2f962d578d163a8b13c9cee127e0b24c01d1906979ac3bc6a763bf95c50999ca97c01")
		// Expected Output
		vmLogOutput    = ""
		feePayerOutput = "000000000000000000000000636f6e7472616374"
		validateOutput = "00"
		// Test ChainConfig
		config = &params.ChainConfig{IstanbulCompatibleBlock: Block5} // Set IstanbulCompatible block number as '5'
	)

	runPrecompiledContractTestWithHFCondition(t, config, []TestData{
		// Expect Normal Behavior
		// Condition 1. Caller Contract Deploy - before IstanbulCompatible, Call - before istanbulCompatible
		{"0x009", vmLogInput, false, Block4, 200, vmLogOutput, nil},
		{"0x00a", feePayerInput, false, Block4, 300, feePayerOutput, nil},
		{"0x00b", validateInput, false, Block4, 5000, validateOutput, nil},
		// Condition 2. Caller Contract Deploy - before IstanbulCompatible, Call - after istanbulCompatible
		{"0x009", vmLogInput, false, Block5, 200, vmLogOutput, nil},
		{"0x00a", feePayerInput, false, Block5, 300, feePayerOutput, nil},
		{"0x00b", validateInput, false, Block5, 5000, validateOutput, nil},
		// Condition 3. Caller Contract Deploy - after IstanbulCompatible, Call - after istanbulCompatible
		{"0x3fd", vmLogInput, true, Block5, 200, vmLogOutput, nil},
		{"0x3fe", feePayerInput, true, Block5, 300, feePayerOutput, nil},
		{"0x3ff", validateInput, true, Block5, 5000, validateOutput, nil},

		// Expect Error Behavior
		// Condition 1. Caller Contract Deploy - before IstanbulCompatible, Call - before istanbulCompatible
		{"0x3fd", vmLogInput, false, Block4, 0, "", kerrors.ErrPrecompiledContractAddress},
		{"0x3fe", feePayerInput, false, Block4, 0, "", kerrors.ErrPrecompiledContractAddress},
		{"0x3ff", validateInput, false, Block4, 0, "", kerrors.ErrPrecompiledContractAddress},
		// Condition 2. Caller Contract Deploy - before IstanbulCompatible, Call - after istanbulCompatible
		{"0x3fd", vmLogInput, false, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
		{"0x3fe", feePayerInput, false, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
		{"0x3ff", validateInput, false, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
		// Condition 3. Caller Contract Deploy - after IstanbulCompatible, Call - after istanbulCompatible
		{"0x009", vmLogInput, true, Block5, math.MaxUint64, "", errors.New("invalid input length")},
		{"0x00a", feePayerInput, true, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
		{"0x00b", validateInput, true, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
	})
}

func TestBn256GasCost(t *testing.T) {
	var (
		// Test Input
		bn256AddInput       = common.Hex2Bytes("18b18acfb4c2c30276db5411368e7185b311dd124691610c5d3b74034e093dc9063c909c4720840cb5134cb9f59fa749755796819658d32efc0d288198f3726607c2b7f58a84bd6145f00c9c2bc0bb1a187f20ff2c92963a88019e7c6a014eed06614e20c147e940f2d70da3f74c9a17df361706a4485c742bd6788478fa17d7")
		bn256ScalarMulInput = common.Hex2Bytes("2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb721611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb20400000000000000000000000000000000000000000000000011138ce750fa15c2")
		bn256PairingInput   = common.Hex2Bytes("1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f593034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf704bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a416782bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c2032c61a830e3c17286de9462bf242fca2883585b93870a73853face6a6bf411198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa")
		// Expected Output
		bn256AddOutput       = "2243525c5efd4b9c3d3c45ac0ca3fe4dd85e830a4ce6b65fa1eeaee202839703301d1d33be6da8e509df21cc35964723180eed7532537db9ae5e7d48f195c915"
		bn256ScalarMulOutput = "070a8d6a982153cae4be29d434e8faef8a47b274a053f5a4ee2a6c9c13c31e5c031b8ce914eba3a9ffb989f9cdd5b0f01943074bf4f0f315690ec3cec6981afc"
		bn256PairingOutput   = "0000000000000000000000000000000000000000000000000000000000000001"
		// Test ChainConfig
		config = &params.ChainConfig{IstanbulCompatibleBlock: Block5}
	)

	runPrecompiledContractTestWithHFCondition(t, config, []TestData{
		// Test whether appropriate gas cost is returned
		// Condition 1. Caller Contract Deploy - before IstanbulCompatible, Call - before istanbulCompatible
		{"0x006", bn256AddInput, false, Block4, params.Bn256AddGasConstantinople, bn256AddOutput, nil},
		{"0x007", bn256ScalarMulInput, false, Block4, params.Bn256ScalarMulGasConstantinople, bn256ScalarMulOutput, nil},
		{"0x008", bn256PairingInput, false, Block4, params.Bn256PairingBaseGasConstantinople + params.Bn256PairingPerPointGasConstantinople*uint64(len(bn256PairingInput)/192), bn256PairingOutput, nil},
		// Condition 2. Caller Contract Deploy - before IstanbulCompatible, Call - after istanbulCompatible
		{"0x006", bn256AddInput, false, Block5, params.Bn256AddGasConstantinople, bn256AddOutput, nil},
		{"0x007", bn256ScalarMulInput, false, Block5, params.Bn256ScalarMulGasConstantinople, bn256ScalarMulOutput, nil},
		{"0x008", bn256PairingInput, false, Block5, params.Bn256PairingBaseGasConstantinople + params.Bn256PairingPerPointGasConstantinople*uint64(len(bn256PairingInput)/192), bn256PairingOutput, nil},
		// Condition 3. Caller Contract Deploy - after IstanbulCompatible, Call - after istanbulCompatible
		{"0x006", bn256AddInput, true, Block5, params.Bn256AddGasIstanbul, bn256AddOutput, nil},
		{"0x007", bn256ScalarMulInput, true, Block5, params.Bn256ScalarMulGasIstanbul, bn256ScalarMulOutput, nil},
		{"0x008", bn256PairingInput, true, Block5, params.Bn256PairingBaseGasIstanbul + params.Bn256PairingPerPointGasIstanbul*uint64(len(bn256PairingInput)/192), bn256PairingOutput, nil},
	})
}
