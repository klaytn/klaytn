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

// TestData
var (
	// Test Condition
	Block4 = big.NewInt(4) // Block number before the istanbulHF
	Block5 = big.NewInt(5) // Block number after the istanbulHF(also istanbulHF Block number)
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
	// Test condition
	isDeployedAfterIstanbulHF bool     // contract deployment time
	block                     *big.Int // block when execution is done
	// Expect field
	expectGas uint64
	expectRet string
	expectErr error
}{
	// Expect Normal Behavior
	// Condition 1. Caller Contract Deploy - before IstanbulCompatible Change, Call - before istanbulCompatible
	{"0x009", vmLogInput, false, Block4, 200, vmLogOutput, nil},
	{"0x00a", feePayerInput, false, Block4, 300, feePayerOutput, nil},
	{"0x00b", validateInput, false, Block4, 5000, validateOutput, nil},
	// Condition 2. Caller Contract Deploy - before IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x009", vmLogInput, false, Block5, 200, vmLogOutput, nil},
	{"0x00a", feePayerInput, false, Block5, 300, feePayerOutput, nil},
	{"0x00b", validateInput, false, Block5, 5000, validateOutput, nil},
	// Condition 3. Caller Contract Deploy - after IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x3fd", vmLogInput, true, Block5, 200, vmLogOutput, nil},
	{"0x3fe", feePayerInput, true, Block5, 300, feePayerOutput, nil},
	{"0x3ff", validateInput, true, Block5, 5000, validateOutput, nil},

	// Expect Error Behavior
	// Condition 1. Caller Contract Deploy - before IstanbulCompatible Change, Call - before istanbulCompatible
	{"0x3fd", vmLogInput, false, Block4, 0, "", kerrors.ErrPrecompiledContractAddress},
	{"0x3fe", feePayerInput, false, Block4, 0, "", kerrors.ErrPrecompiledContractAddress},
	{"0x3ff", validateInput, false, Block4, 0, "", kerrors.ErrPrecompiledContractAddress},
	// Condition 2. Caller Contract Deploy - before IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x3fd", vmLogInput, false, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
	{"0x3fe", feePayerInput, false, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
	{"0x3ff", validateInput, false, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
	// Condition 2. Caller Contract Deploy - after IstanbulCompatible Change, Call - after istanbulCompatible
	{"0x009", vmLogInput, true, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
	{"0x00a", feePayerInput, true, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
	{"0x00b", validateInput, true, Block5, 0, "", kerrors.ErrPrecompiledContractAddress},
}

func TestPrecompiledContractAddressMapping(t *testing.T) {
	config := params.TestChainConfig        // Make ChainConfig for the test
	config.IstanbulCompatibleBlock = Block5 // Set IstanbulCompatible block number as '5'

	// Run Tests
	for _, tc := range PrecompiledContractAddressMappingTestData {
		// Make StateDB
		callerAddr := common.BytesToAddress([]byte("contract"))
		statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()))
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
