package vm

import (
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

// Test condition
type TestField struct {
	codeFormat params.CodeFormat
	block      *big.Int
}

// Expected result
type ExpectField struct {
	gas uint64
	ret string
	err error
}

// TestData
var (
	// Test Condition
	ver1   = params.CodeFormatEVM                   // Indicates that a contract is deployed before istanbulCompatible change
	ver2   = params.CodeFormatEVMIstanbulCompatible // Indicates that a contract is deployed after istanbulCompatible change
	Block4 = big.NewInt(4)                          // Block number before the istanbulCompatible Change
	Block5 = big.NewInt(5)                          // Block number after the istanbulCompatible Change(also istanbulCompatible Block number)
	// Test Input
	vmLogInput    = []byte("Hello")
	feePayerInput = []byte("")
	validateInput = common.Hex2Bytes("00000000000000000000000000000001234567890000000000000000000000000000000000000000000000000000000000000000db38e383310a432f1cec0bfbd765fe67c9932b02295b1f2062a785eec288b2f962d578d163a8b13c9cee127e0b24c01d1906979ac3bc6a763bf95c50999ca97c01")
	// Expected Output
	vmLogOutput    = ""
	feePayerOutput = "000000000000000000000000636f6e7472616374"
	validateOutput = "00"
)

var PrecompiledContractAddressMappingTestData = []struct {
	addr  string
	input []byte
	TestField
	ExpectField
}{
	// Expect Normal Behavior
	// Condition 1. Caller Contract Deploy - before IstanbulCompatible Change, Call - before istanbulCompatible
	{"0x009", vmLogInput, TestField{ver1, Block4}, ExpectField{200, vmLogOutput, nil}},
	{"0x00a", feePayerInput, TestField{ver1, Block4}, ExpectField{300, feePayerOutput, nil}},
	{"0x00b", validateInput, TestField{ver1, Block4}, ExpectField{5000, validateOutput, nil}},
	// Condition 2. Caller Contract Deploy - before IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x009", vmLogInput, TestField{ver1, Block5}, ExpectField{200, vmLogOutput, nil}},
	{"0x00a", feePayerInput, TestField{ver1, Block5}, ExpectField{300, feePayerOutput, nil}},
	{"0x00b", validateInput, TestField{ver1, Block5}, ExpectField{5000, validateOutput, nil}},
	// Condition 3. Caller Contract Deploy - after IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x3fd", vmLogInput, TestField{ver2, Block5}, ExpectField{200, vmLogOutput, nil}},
	{"0x3fe", feePayerInput, TestField{ver2, Block5}, ExpectField{300, feePayerOutput, nil}},
	{"0x3ff", validateInput, TestField{ver2, Block5}, ExpectField{5000, validateOutput, nil}},

	// Expect Error Behavior
	// Condition 1. Caller Contract Deploy - before IstanbulCompatible Change, Call - before istanbulCompatible
	{"0x3fd", vmLogInput, TestField{ver1, Block4}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	{"0x3fe", feePayerInput, TestField{ver1, Block4}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	{"0x3ff", validateInput, TestField{ver1, Block4}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	// Condition 2. Caller Contract Deploy - before IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x3fd", vmLogInput, TestField{ver1, Block5}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	{"0x3fe", feePayerInput, TestField{ver1, Block5}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	{"0x3ff", validateInput, TestField{ver1, Block5}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	// Condition 2. Caller Contract Deploy - after IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x009", vmLogInput, TestField{ver2, Block5}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	{"0x00a", feePayerInput, TestField{ver2, Block5}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
	{"0x00b", validateInput, TestField{ver2, Block5}, ExpectField{0, "", kerrors.ErrPrecompiledContractAddress}},
}

func TestPrecompiledContractAddressMapping(t *testing.T) {
	config := params.TestChainConfig        // Make ChainConfig for the test
	config.IstanbulCompatibleBlock = Block5 // Set IstanbulCompatible block number as '5'

	// Run Tests
	for _, tt := range PrecompiledContractAddressMappingTestData {
		// Make StateDB
		callerAddr := common.BytesToAddress([]byte("contract"))
		statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()))
		statedb.CreateSmartContractAccount(callerAddr, tt.codeFormat)

		// Make EVM environment
		vmctx := Context{
			CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
			Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
			BlockNumber: tt.block,
		}
		vmenv := NewEVM(vmctx, statedb, config, &Config{})

		// run
		ret, gas, err := vmenv.Call(AccountRef(callerAddr), common.HexToAddress(tt.addr), tt.input, math.MaxUint64, new(big.Int))

		// compare with expected data
		assert.Equal(t, tt.ret, common.Bytes2Hex(ret))
		assert.Equal(t, tt.gas, math.MaxUint64-gas)
		assert.Equal(t, tt.err, err)
	}
}
