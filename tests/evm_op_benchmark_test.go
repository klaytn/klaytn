// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/compiler"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BenchmarkEvmOpTestCase struct {
	testName string
	funcName string
	input    []interface{}
}

func BenchmarkEvmOp(t *testing.B) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// multisig10Initial has a initial key pair of multisig10 before the account key update
	multisig10Initial, err := createAnonymousAccount("bb113e82881499a7a361e8354a5b68f6c6885c7bcba09ea2b0891480396c3200")
	require.Equal(t, nil, err)

	multisig10, err := createMultisigAccount(uint(1),
		[]uint{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]string{
			"bb113e82881499a7a361e8354a5b68f6c6885c7bcba09ea2b0891480396c322e",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e989",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e98A",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e98B",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e98C",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e98D",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e98E",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e98F",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e999",
			"c32c471b732e2f56103e2f8e8cfd52792ef548f05f326e546a7d1fbf9d0419ec",
		},
		multisig10Initial.Addr)

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	gasPrice := new(big.Int).SetUint64(0)
	gasLimit := uint64(100000000000)

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	filename := string("../contracts/computationcost/opcodeBench.sol")
	contracts, err := compiler.CompileSolidityOrLoad("", filename)
	require.NoError(t, err)

	var c compiler.Contract
	for k, v := range contracts {
		if strings.Contains(k, "StopContract") {
			c = *v
			break
		}
	}

	contractCode := c.Code
	stopContractAbiJson, err := json.Marshal(c.Info.AbiDefinition)
	require.NoError(t, err)

	stopContractAbi, err := abi.JSON(strings.NewReader(string(stopContractAbiJson)))
	require.NoError(t, err)

	StopContractStopInput, err := stopContractAbi.Pack("Sstop")
	require.NoError(t, err)

	for k, v := range contracts {
		if strings.Contains(k, "OpCodeBenchmarkContract") {
			c = *v
			break
		}
	}

	abiJson, err := json.Marshal(c.Info.AbiDefinition)
	require.NoError(t, err)
	abiStr := string(abiJson)

	contractAddrs := make(map[string]common.Address)
	contractNames := make([]string, 0, len(contracts))

	// Deploy smart contract (reservoir -> contract), create multisig account.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)

		for name, contract := range contracts {
			fmt.Println("deploying contract ", name)
			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:         reservoir.Nonce,
				types.TxValueKeyFrom:          reservoir.Addr,
				types.TxValueKeyTo:            (*common.Address)(nil),
				types.TxValueKeyAmount:        amount,
				types.TxValueKeyGasLimit:      gasLimit,
				types.TxValueKeyGasPrice:      gasPrice,
				types.TxValueKeyHumanReadable: false,
				types.TxValueKeyData:          common.FromHex(contract.Code),
				types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
			require.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.Keys)
			require.Equal(t, nil, err)

			txs = append(txs, tx)
			reservoir.Nonce += 1

			contractNames = append(contractNames, name)
		}

		{
			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:      multisig10Initial.Nonce,
				types.TxValueKeyFrom:       multisig10Initial.Addr,
				types.TxValueKeyGasLimit:   gasLimit,
				types.TxValueKeyGasPrice:   gasPrice,
				types.TxValueKeyAccountKey: multisig10.AccKey,
				types.TxValueKeyFeePayer:   reservoir.Addr,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, multisig10Initial.Keys)
			assert.Equal(t, nil, err)

			err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
			assert.Equal(t, nil, err)

			txs = append(txs, tx)
			multisig10Initial.Nonce++
		}

		require.NoError(t, bcdata.GenABlockWithTransactions(accountMap, txs, prof))
	}

	blockHash := bcdata.bc.CurrentBlock().Hash()
	receipts := bcdata.bc.GetReceiptsByBlockHash(blockHash)
	for i := 0; i < len(contracts); i++ {
		addr := receipts[i].ContractAddress
		n := strings.Split(contractNames[i], ":")[1]
		contractAddrs[n] = addr
		fmt.Println("createdaddr", addr.String())
	}

	// prepare data for validate sender (reservoir)
	validateSenderData1 := make([]byte, 0, 20+32+65)
	{
		hash := crypto.Keccak256Hash(reservoir.Addr.Bytes())
		sigs, err := crypto.Sign(hash.Bytes(), reservoir.Keys[0])
		require.NoError(t, err)

		validateSenderData1 = append([]byte{}, reservoir.Addr.Bytes()...)
		validateSenderData1 = append(validateSenderData1, hash.Bytes()...)
		validateSenderData1 = append(validateSenderData1, sigs...)
	}

	// prepare data for validate sender (multisig10)
	validateSenderMultisig10_2 := make([]byte, 0, 20+32+65*2)
	validateSenderMultisig10_3 := make([]byte, 0, 20+32+65*3)
	validateSenderMultisig10_4 := make([]byte, 0, 20+32+65*4)
	validateSenderMultisig10_5 := make([]byte, 0, 20+32+65*5)
	validateSenderMultisig10_6 := make([]byte, 0, 20+32+65*6)
	validateSenderMultisig10_7 := make([]byte, 0, 20+32+65*7)
	validateSenderMultisig10_8 := make([]byte, 0, 20+32+65*8)
	validateSenderMultisig10_9 := make([]byte, 0, 20+32+65*9)
	validateSenderMultisig10_10 := make([]byte, 0, 20+32+65*10)
	{
		hash := crypto.Keccak256Hash(multisig10.Addr.Bytes())
		sigs := make([][]byte, 10)
		for i := 0; i < 10; i++ {
			s, err := crypto.Sign(hash.Bytes(), multisig10.Keys[i])
			require.NoError(t, err)
			sigs[i] = s
		}
		validateSenderMultisig := make([]byte, 0, 20+32)
		validateSenderMultisig = append(validateSenderMultisig, multisig10.Addr.Bytes()...)
		validateSenderMultisig = append(validateSenderMultisig, hash.Bytes()...)

		validateSenderMultisig10_2 = append(validateSenderMultisig10_2, validateSenderMultisig...)
		validateSenderMultisig10_2 = append(validateSenderMultisig10_2, sigs[0]...)
		validateSenderMultisig10_2 = append(validateSenderMultisig10_2, sigs[1]...)

		validateSenderMultisig10_3 = append(validateSenderMultisig10_3, validateSenderMultisig10_2...)
		validateSenderMultisig10_3 = append(validateSenderMultisig10_3, sigs[2]...)

		validateSenderMultisig10_4 = append(validateSenderMultisig10_4, validateSenderMultisig10_3...)
		validateSenderMultisig10_4 = append(validateSenderMultisig10_4, sigs[3]...)

		validateSenderMultisig10_5 = append(validateSenderMultisig10_5, validateSenderMultisig10_4...)
		validateSenderMultisig10_5 = append(validateSenderMultisig10_5, sigs[4]...)

		validateSenderMultisig10_6 = append(validateSenderMultisig10_6, validateSenderMultisig10_5...)
		validateSenderMultisig10_6 = append(validateSenderMultisig10_6, sigs[5]...)

		validateSenderMultisig10_7 = append(validateSenderMultisig10_7, validateSenderMultisig10_6...)
		validateSenderMultisig10_7 = append(validateSenderMultisig10_7, sigs[6]...)

		validateSenderMultisig10_8 = append(validateSenderMultisig10_8, validateSenderMultisig10_7...)
		validateSenderMultisig10_8 = append(validateSenderMultisig10_8, sigs[7]...)

		validateSenderMultisig10_9 = append(validateSenderMultisig10_9, validateSenderMultisig10_8...)
		validateSenderMultisig10_9 = append(validateSenderMultisig10_9, sigs[8]...)

		validateSenderMultisig10_10 = append(validateSenderMultisig10_10, validateSenderMultisig10_9...)
		validateSenderMultisig10_10 = append(validateSenderMultisig10_10, sigs[9]...)
	}

	loopCnt := big.NewInt(1000000)
	// loopCnt := big.NewInt(10000)
	// loopCnt := big.NewInt(1)

	testcases := []struct {
		testName string
		funcName string
		input    []interface{}
	}{
		{
			"Add/low",
			"Add",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Add/high",
			"Add",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Sub/low",
			"Sub",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Sub/high",
			"Sub",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Mul/low",
			"Mul",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Mul/high",
			"Mul",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Div/low",
			"Div",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Div/high",
			"Div",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Sdiv/low",
			"Sdiv",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Sdiv/high",
			"Sdiv",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Mod/low",
			"Mod",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Mod/high",
			"Mod",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Smod/low",
			"Smod",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Smod/high",
			"Smod",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Exp/low",
			"Exp",
			[]interface{}{
				loopCnt, big.NewInt(10), big.NewInt(10),
			},
		},
		{
			"Exp/high",
			"Exp",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Addmod/low",
			"Addmod",
			[]interface{}{
				loopCnt, big.NewInt(10), big.NewInt(10), big.NewInt(10),
			},
		},
		{
			"Addmod/high",
			"Addmod",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Mulmod/low",
			"Mulmod",
			[]interface{}{
				loopCnt, big.NewInt(10), big.NewInt(10), big.NewInt(10),
			},
		},
		{
			"Mulmod/high",
			"Mulmod",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Not",
			"Not",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Lt",
			"Lt",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Gt",
			"Gt",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Slt",
			"Slt",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Sgt",
			"Sgt",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Eq",
			"Eq",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Iszero",
			"Iszero",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"And",
			"And",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Or",
			"Or",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Xor",
			"Xor",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Byte",
			"Byte",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Shl/low",
			"Shl",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Shl/high",
			"Shl",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Shr/low",
			"Shr",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Shr/high",
			"Shr",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Sar/low",
			"Sar",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"Sar/high",
			"Sar",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"SignExtend/low",
			"SignExtend",
			[]interface{}{
				loopCnt, big.NewInt(1), big.NewInt(1),
			},
		},
		{
			"SignExtend/high",
			"SignExtend",
			[]interface{}{
				loopCnt,
				new(big.Int).SetBytes(common.FromHex("32")),
				new(big.Int).SetBytes(common.FromHex("7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe")),
			},
		},
		{
			"Sha3",
			"Sha3",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Pc",
			"Pc",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup1",
			"Dup1",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup2",
			"Dup2",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup3",
			"Dup3",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup4",
			"Dup4",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup5",
			"Dup5",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup6",
			"Dup6",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup7",
			"Dup7",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup8",
			"Dup8",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup9",
			"Dup9",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup10",
			"Dup10",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup11",
			"Dup11",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup12",
			"Dup12",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup13",
			"Dup13",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup14",
			"Dup14",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup15",
			"Dup15",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Dup16",
			"Dup16",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap1",
			"Swap1",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap2",
			"Swap2",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap3",
			"Swap3",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap4",
			"Swap4",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap5",
			"Swap5",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap6",
			"Swap6",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap7",
			"Swap7",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap8",
			"Swap8",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap9",
			"Swap9",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap10",
			"Swap10",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap11",
			"Swap11",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap12",
			"Swap12",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap13",
			"Swap13",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap14",
			"Swap14",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap15",
			"Swap15",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Swap16",
			"Swap16",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Call",
			"Call",
			[]interface{}{
				loopCnt,
				contractAddrs["StopContract"],
			},
		},
		{
			"CallCode",
			"CallCode",
			[]interface{}{
				loopCnt,
				contractAddrs["StopContract"],
			},
		},
		{
			"StaticCall",
			"StaticCall",
			[]interface{}{
				loopCnt,
				contractAddrs["StopContract"], StopContractStopInput,
			},
		},
		{
			"DelegateCall",
			"DelegateCall",
			[]interface{}{
				loopCnt,
				contractAddrs["StopContract"],
			},
		},
		{
			"Create",
			"Create",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Create2",
			"Create2",
			[]interface{}{
				loopCnt, common.FromHex(contractCode), big.NewInt(1234),
			},
		},
		{
			"Sstore",
			"Sstore",
			[]interface{}{
				loopCnt,
				big.NewInt(0),
			},
		},
		{
			"Sload",
			"Sload",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Mstore/10",
			"Mstore",
			[]interface{}{
				loopCnt, big.NewInt(10),
			},
		},
		{
			"Mstore/100",
			"Mstore",
			[]interface{}{
				loopCnt, big.NewInt(100),
			},
		},
		{
			"Mstore/100000",
			"Mstore",
			[]interface{}{
				loopCnt, big.NewInt(100000),
			},
		},
		{
			"Mstore/1000000",
			"Mstore",
			[]interface{}{
				loopCnt, big.NewInt(1000000),
			},
		},
		{
			"Mload/10",
			"Mload",
			[]interface{}{
				loopCnt, big.NewInt(10),
			},
		},
		{
			"Mload/100",
			"Mload",
			[]interface{}{
				loopCnt, big.NewInt(100),
			},
		},
		{
			"Mload/100000",
			"Mload",
			[]interface{}{
				loopCnt, big.NewInt(100000),
			},
		},
		{
			"Mload/1000000",
			"Mload",
			[]interface{}{
				loopCnt, big.NewInt(1000000),
			},
		},
		{
			"Msize",
			"Msize",
			[]interface{}{
				loopCnt, big.NewInt(1000),
			},
		},
		{
			"Gas",
			"Gas",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Address",
			"Address",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"Balance",
			"Balance",
			[]interface{}{
				loopCnt, reservoir.Addr,
			},
		},
		{
			"Caller",
			"Caller",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"CallValue",
			"CallValue",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"CallDataLoad",
			"CallDataLoad",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"CallDataSize",
			"CallDataSize",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"CallDataCopy/1",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("01"),
				big.NewInt(1),
			},
		},
		{
			"CallDataCopy/16",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516"),
				big.NewInt(16),
			},
		},
		{
			"CallDataCopy/32",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(32),
			},
		},
		{
			"CallDataCopy/64",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(64),
			},
		},
		{
			"CallDataCopy/128",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(128),
			},
		},
		{
			"CallDataCopy/256",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(256),
			},
		},
		{
			"CallDataCopy/512",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(512),
			},
		},
		{
			"CallDataCopy/1024",
			"CallDataCopy",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1024),
			},
		},
		{
			"CodeSize",
			"CodeSize",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"CodeCopy/32",
			"CodeCopy",
			[]interface{}{
				loopCnt, big.NewInt(32),
			},
		},
		{
			"CodeCopy/64",
			"CodeCopy",
			[]interface{}{
				loopCnt, big.NewInt(64),
			},
		},
		{
			"CodeCopy/128",
			"CodeCopy",
			[]interface{}{
				loopCnt, big.NewInt(128),
			},
		},
		{
			"CodeCopy/256",
			"CodeCopy",
			[]interface{}{
				loopCnt, big.NewInt(256),
			},
		},
		{
			"CodeCopy/512",
			"CodeCopy",
			[]interface{}{
				loopCnt, big.NewInt(512),
			},
		},
		{
			"CodeCopy/1024",
			"CodeCopy",
			[]interface{}{
				loopCnt, big.NewInt(1024),
			},
		},
		{
			"ExtCodeSize",
			"ExtCodeSize",
			[]interface{}{
				loopCnt, contractAddrs["OpCodeBenchmarkContract"],
			},
		},
		{
			"ExtCodeCopy/32",
			"ExtCodeCopy",
			[]interface{}{
				loopCnt, contractAddrs["OpCodeBenchmarkContract"], big.NewInt(32),
			},
		},
		{
			"ExtCodeCopy/64",
			"ExtCodeCopy",
			[]interface{}{
				loopCnt, contractAddrs["OpCodeBenchmarkContract"], big.NewInt(64),
			},
		},
		{
			"ExtCodeCopy/128",
			"ExtCodeCopy",
			[]interface{}{
				loopCnt, contractAddrs["OpCodeBenchmarkContract"], big.NewInt(128),
			},
		},
		{
			"ExtCodeCopy/256",
			"ExtCodeCopy",
			[]interface{}{
				loopCnt, contractAddrs["OpCodeBenchmarkContract"], big.NewInt(256),
			},
		},
		{
			"ExtCodeCopy/512",
			"ExtCodeCopy",
			[]interface{}{
				loopCnt, contractAddrs["OpCodeBenchmarkContract"], big.NewInt(512),
			},
		},
		{
			"ExtCodeCopy/1024",
			"ExtCodeCopy",
			[]interface{}{
				loopCnt, contractAddrs["OpCodeBenchmarkContract"], big.NewInt(1024),
			},
		},
		{
			"ReturnDataSize",
			"ReturnDataSize",
			[]interface{}{
				loopCnt,
			},
		},
		{
			"ReturnDataCopy/1",
			"ReturnDataCopy",
			[]interface{}{
				loopCnt, big.NewInt(1), contractAddrs["StopContract"],
			},
		},
		{
			"ReturnDataCopy/32",
			"ReturnDataCopy",
			[]interface{}{
				loopCnt, big.NewInt(32), contractAddrs["StopContract"],
			},
		},
		{
			"ReturnDataCopy/64",
			"ReturnDataCopy",
			[]interface{}{
				loopCnt, big.NewInt(32), contractAddrs["StopContract"],
			},
		},
		{
			"ReturnDataCopy/128",
			"ReturnDataCopy",
			[]interface{}{
				loopCnt, big.NewInt(128), contractAddrs["StopContract"],
			},
		},
		{
			"ReturnDataCopy/256",
			"ReturnDataCopy",
			[]interface{}{
				loopCnt, big.NewInt(256), contractAddrs["StopContract"],
			},
		},
		{
			"ReturnDataCopy/512",
			"ReturnDataCopy",
			[]interface{}{
				loopCnt, big.NewInt(512), contractAddrs["StopContract"],
			},
		},
		{
			"ReturnDataCopy/1024",
			"ReturnDataCopy",
			[]interface{}{
				loopCnt, big.NewInt(1024), contractAddrs["StopContract"],
			},
		},
		// Not supported in solc-0.4.24
		//{
		//	"ExtCodeHash",
		//	"ExtCodeHash",
		//	[]interface{}{
		//		loopCnt, contractAddrs["OpCodeBenchmarkContract"],
		//	},
		//},

		{
			"Log0/32",
			"Log0",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"Log0/64",
			"Log0",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"Log0/128",
			"Log0",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"Log0/256",
			"Log0",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"Log0/512",
			"Log0",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"Log0/1024",
			"Log0",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},

		{
			"Log1/32",
			"Log1",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
			},
		},
		{
			"Log1/64",
			"Log1",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
			},
		},
		{
			"Log1/128",
			"Log1",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
			},
		},
		{
			"Log1/256",
			"Log1",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
			},
		},
		{
			"Log1/512",
			"Log1",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
			},
		},
		{
			"Log1/1024",
			"Log1",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
			},
		},

		{
			"Log2/32",
			"Log2",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
			},
		},
		{
			"Log2/64",
			"Log2",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
			},
		},
		{
			"Log2/128",
			"Log2",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
			},
		},
		{
			"Log2/256",
			"Log2",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
			},
		},
		{
			"Log2/512",
			"Log2",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
			},
		},
		{
			"Log2/1024",
			"Log2",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
			},
		},

		{
			"Log3/32",
			"Log3",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
			},
		},
		{
			"Log3/64",
			"Log3",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
			},
		},
		{
			"Log3/128",
			"Log3",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
			},
		},
		{
			"Log3/256",
			"Log3",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
			},
		},
		{
			"Log3/512",
			"Log3",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
			},
		},
		{
			"Log3/1024",
			"Log3",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
			},
		},

		{
			"Log4/32",
			"Log4",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(4),
			},
		},
		{
			"Log4/64",
			"Log4",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(4),
			},
		},
		{
			"Log4/128",
			"Log4",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(4),
			},
		},
		{
			"Log4/256",
			"Log4",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(4),
			},
		},
		{
			"Log4/512",
			"Log4",
			[]interface{}{
				loopCnt,
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(4),
			},
		},
		{
			"Log4/1024",
			"Log4",
			[]interface{}{
				loopCnt,
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
				big.NewInt(4),
			},
		},

		{
			"Origin",
			"Origin",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"GasPrice",
			"GasPrice",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"BlockHash",
			"BlockHash",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"Coinbase",
			"Coinbase",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"Timestamp",
			"Timestamp",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"Number",
			"Number",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"Difficulty",
			"Difficulty",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"GasLimit",
			"GasLimit",
			[]interface{}{
				loopCnt,
			},
		},

		{
			"Combination",
			"Combination",
			[]interface{}{
				loopCnt, big.NewInt(100),
			},
		},

		{
			"ecrecover",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x01"),
				common.FromHex("38d18acb67d25c8bb9942764b62f18e17054f66a817bd4295423adf9ed98873e000000000000000000000000000000000000000000000000000000000000001b38d18acb67d25c8bb9942764b62f18e17054f66a817bd4295423adf9ed98873e789d1dd423d25f0772d2748d60f7e4b81bb14d086eba8e8e8efb6dcff8a4ae02"),
			},
		},

		{
			"sha256/32",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x02"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"sha256/64",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x02"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"sha256/128",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x02"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"sha256/256",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x02"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},

		// ripemd160hash
		{
			"ripemd160/32",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x03"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"ripemd160/64",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x03"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"ripemd160/128",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x03"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"ripemd160/256",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x03"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},

		// datacopy
		{
			"dataCopy/32",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x04"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"dataCopy/64",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x04"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"dataCopy/128",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x04"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"dataCopy/256",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x04"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},

		// bigExpMod
		{
			"bigExpMod/eip_example1",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("0000000000000000000000000000000000000000000000000000000000000001" +
					"0000000000000000000000000000000000000000000000000000000000000020" +
					"0000000000000000000000000000000000000000000000000000000000000020" +
					"03" +
					"fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2e" +
					"fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"),
			},
		},
		{
			"bigExpMod/eip_example2",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("0000000000000000000000000000000000000000000000000000000000000000" +
					"0000000000000000000000000000000000000000000000000000000000000020" +
					"0000000000000000000000000000000000000000000000000000000000000020" +
					"fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2e" +
					"fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"),
			},
		},
		{
			"bigExpMod/eip_example2/mod",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("0000000000000000000000000000000000000000000000000000000000000001" +
					"0000000000000000000000000000000000000000000000000000000000000020" +
					"0000000000000000000000000000000000000000000000000000000000000020" +
					"fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2e" +
					"fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"),
			},
		},
		{
			"bigExpMod/modexp_modsize0_returndatasize.json/Byzantium/4",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000f3f14010101"),
			},
		},
		{
			"bigExpMod/nagydani-1-square",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000040e09ad9675465c53a109fac66a445c91b292d2bb2c5268addb30cd82f80fcb0033ff97c80a5fc6f39193ae969c6ede6710a6b7ac27078a06d90ef1c72e5c85fb502fc9e1f6beb81516545975218075ec2af118cd8798df6e08a147c60fd6095ac2bb02c2908cf4dd7c81f11c289e4bce98f3553768f392a80ce22bf5c4f4a248c6b"),
			},
		},
		{
			"bigExpMod/nagydani-1-qube",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000040e09ad9675465c53a109fac66a445c91b292d2bb2c5268addb30cd82f80fcb0033ff97c80a5fc6f39193ae969c6ede6710a6b7ac27078a06d90ef1c72e5c85fb503fc9e1f6beb81516545975218075ec2af118cd8798df6e08a147c60fd6095ac2bb02c2908cf4dd7c81f11c289e4bce98f3553768f392a80ce22bf5c4f4a248c6b"),
			},
		},
		{
			"bigExpMod/nagydani-1-pow0x10001",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000040e09ad9675465c53a109fac66a445c91b292d2bb2c5268addb30cd82f80fcb0033ff97c80a5fc6f39193ae969c6ede6710a6b7ac27078a06d90ef1c72e5c85fb5010001fc9e1f6beb81516545975218075ec2af118cd8798df6e08a147c60fd6095ac2bb02c2908cf4dd7c81f11c289e4bce98f3553768f392a80ce22bf5c4f4a248c6b"),
			},
		},
		{
			"bigExpMod/nagydani-2-square",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000080cad7d991a00047dd54d3399b6b0b937c718abddef7917c75b6681f40cc15e2be0003657d8d4c34167b2f0bbbca0ccaa407c2a6a07d50f1517a8f22979ce12a81dcaf707cc0cebfc0ce2ee84ee7f77c38b9281b9822a8d3de62784c089c9b18dcb9a2a5eecbede90ea788a862a9ddd9d609c2c52972d63e289e28f6a590ffbf5102e6d893b80aeed5e6e9ce9afa8a5d5675c93a32ac05554cb20e9951b2c140e3ef4e433068cf0fb73bc9f33af1853f64aa27a0028cbf570d7ac9048eae5dc7b28c87c31e5810f1e7fa2cda6adf9f1076dbc1ec1238560071e7efc4e9565c49be9e7656951985860a558a754594115830bcdb421f741408346dd5997bb01c287087"),
			},
		},
		{
			"bigExpMod/nagydani-2-qube",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000080cad7d991a00047dd54d3399b6b0b937c718abddef7917c75b6681f40cc15e2be0003657d8d4c34167b2f0bbbca0ccaa407c2a6a07d50f1517a8f22979ce12a81dcaf707cc0cebfc0ce2ee84ee7f77c38b9281b9822a8d3de62784c089c9b18dcb9a2a5eecbede90ea788a862a9ddd9d609c2c52972d63e289e28f6a590ffbf5103e6d893b80aeed5e6e9ce9afa8a5d5675c93a32ac05554cb20e9951b2c140e3ef4e433068cf0fb73bc9f33af1853f64aa27a0028cbf570d7ac9048eae5dc7b28c87c31e5810f1e7fa2cda6adf9f1076dbc1ec1238560071e7efc4e9565c49be9e7656951985860a558a754594115830bcdb421f741408346dd5997bb01c287087"),
			},
		},
		{
			"bigExpMod/nagydani-2-pow0x10001",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000080cad7d991a00047dd54d3399b6b0b937c718abddef7917c75b6681f40cc15e2be0003657d8d4c34167b2f0bbbca0ccaa407c2a6a07d50f1517a8f22979ce12a81dcaf707cc0cebfc0ce2ee84ee7f77c38b9281b9822a8d3de62784c089c9b18dcb9a2a5eecbede90ea788a862a9ddd9d609c2c52972d63e289e28f6a590ffbf51010001e6d893b80aeed5e6e9ce9afa8a5d5675c93a32ac05554cb20e9951b2c140e3ef4e433068cf0fb73bc9f33af1853f64aa27a0028cbf570d7ac9048eae5dc7b28c87c31e5810f1e7fa2cda6adf9f1076dbc1ec1238560071e7efc4e9565c49be9e7656951985860a558a754594115830bcdb421f741408346dd5997bb01c287087"),
			},
		},
		{
			"bigExpMod/nagydani-3-square",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000100c9130579f243e12451760976261416413742bd7c91d39ae087f46794062b8c239f2a74abf3918605a0e046a7890e049475ba7fbb78f5de6490bd22a710cc04d30088179a919d86c2da62cf37f59d8f258d2310d94c24891be2d7eeafaa32a8cb4b0cfe5f475ed778f45907dc8916a73f03635f233f7a77a00a3ec9ca6761a5bbd558a2318ecd0caa1c5016691523e7e1fa267dd35e70c66e84380bdcf7c0582f540174e572c41f81e93da0b757dff0b0fe23eb03aa19af0bdec3afb474216febaacb8d0381e631802683182b0fe72c28392539850650b70509f54980241dc175191a35d967288b532a7a8223ce2440d010615f70df269501944d4ec16fe4a3cb02d7a85909174757835187cb52e71934e6c07ef43b4c46fc30bbcd0bc72913068267c54a4aabebb493922492820babdeb7dc9b1558fcf7bd82c37c82d3147e455b623ab0efa752fe0b3a67ca6e4d126639e645a0bf417568adbb2a6a4eef62fa1fa29b2a5a43bebea1f82193a7dd98eb483d09bb595af1fa9c97c7f41f5649d976aee3e5e59e2329b43b13bea228d4a93f16ba139ccb511de521ffe747aa2eca664f7c9e33da59075cc335afcd2bf3ae09765f01ab5a7c3e3938ec168b74724b5074247d200d9970382f683d6059b94dbc336603d1dfee714e4b447ac2fa1d99ecb4961da2854e03795ed758220312d101e1e3d87d5313a6d052aebde75110363d"),
			},
		},
		{
			"bigExpMod/nagydani-3-qube",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000100c9130579f243e12451760976261416413742bd7c91d39ae087f46794062b8c239f2a74abf3918605a0e046a7890e049475ba7fbb78f5de6490bd22a710cc04d30088179a919d86c2da62cf37f59d8f258d2310d94c24891be2d7eeafaa32a8cb4b0cfe5f475ed778f45907dc8916a73f03635f233f7a77a00a3ec9ca6761a5bbd558a2318ecd0caa1c5016691523e7e1fa267dd35e70c66e84380bdcf7c0582f540174e572c41f81e93da0b757dff0b0fe23eb03aa19af0bdec3afb474216febaacb8d0381e631802683182b0fe72c28392539850650b70509f54980241dc175191a35d967288b532a7a8223ce2440d010615f70df269501944d4ec16fe4a3cb03d7a85909174757835187cb52e71934e6c07ef43b4c46fc30bbcd0bc72913068267c54a4aabebb493922492820babdeb7dc9b1558fcf7bd82c37c82d3147e455b623ab0efa752fe0b3a67ca6e4d126639e645a0bf417568adbb2a6a4eef62fa1fa29b2a5a43bebea1f82193a7dd98eb483d09bb595af1fa9c97c7f41f5649d976aee3e5e59e2329b43b13bea228d4a93f16ba139ccb511de521ffe747aa2eca664f7c9e33da59075cc335afcd2bf3ae09765f01ab5a7c3e3938ec168b74724b5074247d200d9970382f683d6059b94dbc336603d1dfee714e4b447ac2fa1d99ecb4961da2854e03795ed758220312d101e1e3d87d5313a6d052aebde75110363d"),
			},
		},
		{
			"bigExpMod/nagydani-3-pow0x10001",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000100c9130579f243e12451760976261416413742bd7c91d39ae087f46794062b8c239f2a74abf3918605a0e046a7890e049475ba7fbb78f5de6490bd22a710cc04d30088179a919d86c2da62cf37f59d8f258d2310d94c24891be2d7eeafaa32a8cb4b0cfe5f475ed778f45907dc8916a73f03635f233f7a77a00a3ec9ca6761a5bbd558a2318ecd0caa1c5016691523e7e1fa267dd35e70c66e84380bdcf7c0582f540174e572c41f81e93da0b757dff0b0fe23eb03aa19af0bdec3afb474216febaacb8d0381e631802683182b0fe72c28392539850650b70509f54980241dc175191a35d967288b532a7a8223ce2440d010615f70df269501944d4ec16fe4a3cb010001d7a85909174757835187cb52e71934e6c07ef43b4c46fc30bbcd0bc72913068267c54a4aabebb493922492820babdeb7dc9b1558fcf7bd82c37c82d3147e455b623ab0efa752fe0b3a67ca6e4d126639e645a0bf417568adbb2a6a4eef62fa1fa29b2a5a43bebea1f82193a7dd98eb483d09bb595af1fa9c97c7f41f5649d976aee3e5e59e2329b43b13bea228d4a93f16ba139ccb511de521ffe747aa2eca664f7c9e33da59075cc335afcd2bf3ae09765f01ab5a7c3e3938ec168b74724b5074247d200d9970382f683d6059b94dbc336603d1dfee714e4b447ac2fa1d99ecb4961da2854e03795ed758220312d101e1e3d87d5313a6d052aebde75110363d"),
			},
		},
		{
			"bigExpMod/nagydani-4-square",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000200db34d0e438249c0ed685c949cc28776a05094e1c48691dc3f2dca5fc3356d2a0663bd376e4712839917eb9a19c670407e2c377a2de385a3ff3b52104f7f1f4e0c7bf7717fb913896693dc5edbb65b760ef1b00e42e9d8f9af17352385e1cd742c9b006c0f669995cb0bb21d28c0aced2892267637b6470d8cee0ab27fc5d42658f6e88240c31d6774aa60a7ebd25cd48b56d0da11209f1928e61005c6eb709f3e8e0aaf8d9b10f7d7e296d772264dc76897ccdddadc91efa91c1903b7232a9e4c3b941917b99a3bc0c26497dedc897c25750af60237aa67934a26a2bc491db3dcc677491944bc1f51d3e5d76b8d846a62db03dedd61ff508f91a56d71028125035c3a44cbb041497c83bf3e4ae2a9613a401cc721c547a2afa3b16a2969933d3626ed6d8a7428648f74122fd3f2a02a20758f7f693892c8fd798b39abac01d18506c45e71432639e9f9505719ee822f62ccbf47f6850f096ff77b5afaf4be7d772025791717dbe5abf9b3f40cff7d7aab6f67e38f62faf510747276e20a42127e7500c444f9ed92baf65ade9e836845e39c4316d9dce5f8e2c8083e2c0acbb95296e05e51aab13b6b8f53f06c9c4276e12b0671133218cc3ea907da3bd9a367096d9202128d14846cc2e20d56fc8473ecb07cecbfb8086919f3971926e7045b853d85a69d026195c70f9f7a823536e2a8f4b3e12e94d9b53a934353451094b8102df3143a0057457d75e8c708b6337a6f5a4fd1a06727acf9fb93e2993c62f3378b37d56c85e7b1e00f0145ebf8e4095bd723166293c60b6ac1252291ef65823c9e040ddad14969b3b340a4ef714db093a587c37766d68b8d6b5016e741587e7e6bf7e763b44f0247e64bae30f994d248bfd20541a333e5b225ef6a61199e301738b1e688f70ec1d7fb892c183c95dc543c3e12adf8a5e8b9ca9d04f9445cced3ab256f29e998e69efaa633a7b60e1db5a867924ccab0a171d9d6e1098dfa15acde9553de599eaa56490c8f411e4985111f3d40bddfc5e301edb01547b01a886550a61158f7e2033c59707789bf7c854181d0c2e2a42a93cf09209747d7082e147eb8544de25c3eb14f2e35559ea0c0f5877f2f3fc92132c0ae9da4e45b2f6c866a224ea6d1f28c05320e287750fbc647368d41116e528014cc1852e5531d53e4af938374daba6cee4baa821ed07117253bb3601ddd00d59a3d7fb2ef1f5a2fbba7c429f0cf9a5b3462410fd833a69118f8be9c559b1000cc608fd877fb43f8e65c2d1302622b944462579056874b387208d90623fcdaf93920ca7a9e4ba64ea208758222ad868501cc2c345e2d3a5ea2a17e5069248138c8a79c0251185d29ee73e5afab5354769142d2bf0cb6712727aa6bf84a6245fcdae66e4938d84d1b9dd09a884818622080ff5f98942fb20acd7e0c916c2d5ea7ce6f7e173315384518f"),
			},
		},
		{
			"bigExpMod/nagydani-4-qube",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000200db34d0e438249c0ed685c949cc28776a05094e1c48691dc3f2dca5fc3356d2a0663bd376e4712839917eb9a19c670407e2c377a2de385a3ff3b52104f7f1f4e0c7bf7717fb913896693dc5edbb65b760ef1b00e42e9d8f9af17352385e1cd742c9b006c0f669995cb0bb21d28c0aced2892267637b6470d8cee0ab27fc5d42658f6e88240c31d6774aa60a7ebd25cd48b56d0da11209f1928e61005c6eb709f3e8e0aaf8d9b10f7d7e296d772264dc76897ccdddadc91efa91c1903b7232a9e4c3b941917b99a3bc0c26497dedc897c25750af60237aa67934a26a2bc491db3dcc677491944bc1f51d3e5d76b8d846a62db03dedd61ff508f91a56d71028125035c3a44cbb041497c83bf3e4ae2a9613a401cc721c547a2afa3b16a2969933d3626ed6d8a7428648f74122fd3f2a02a20758f7f693892c8fd798b39abac01d18506c45e71432639e9f9505719ee822f62ccbf47f6850f096ff77b5afaf4be7d772025791717dbe5abf9b3f40cff7d7aab6f67e38f62faf510747276e20a42127e7500c444f9ed92baf65ade9e836845e39c4316d9dce5f8e2c8083e2c0acbb95296e05e51aab13b6b8f53f06c9c4276e12b0671133218cc3ea907da3bd9a367096d9202128d14846cc2e20d56fc8473ecb07cecbfb8086919f3971926e7045b853d85a69d026195c70f9f7a823536e2a8f4b3e12e94d9b53a934353451094b8103df3143a0057457d75e8c708b6337a6f5a4fd1a06727acf9fb93e2993c62f3378b37d56c85e7b1e00f0145ebf8e4095bd723166293c60b6ac1252291ef65823c9e040ddad14969b3b340a4ef714db093a587c37766d68b8d6b5016e741587e7e6bf7e763b44f0247e64bae30f994d248bfd20541a333e5b225ef6a61199e301738b1e688f70ec1d7fb892c183c95dc543c3e12adf8a5e8b9ca9d04f9445cced3ab256f29e998e69efaa633a7b60e1db5a867924ccab0a171d9d6e1098dfa15acde9553de599eaa56490c8f411e4985111f3d40bddfc5e301edb01547b01a886550a61158f7e2033c59707789bf7c854181d0c2e2a42a93cf09209747d7082e147eb8544de25c3eb14f2e35559ea0c0f5877f2f3fc92132c0ae9da4e45b2f6c866a224ea6d1f28c05320e287750fbc647368d41116e528014cc1852e5531d53e4af938374daba6cee4baa821ed07117253bb3601ddd00d59a3d7fb2ef1f5a2fbba7c429f0cf9a5b3462410fd833a69118f8be9c559b1000cc608fd877fb43f8e65c2d1302622b944462579056874b387208d90623fcdaf93920ca7a9e4ba64ea208758222ad868501cc2c345e2d3a5ea2a17e5069248138c8a79c0251185d29ee73e5afab5354769142d2bf0cb6712727aa6bf84a6245fcdae66e4938d84d1b9dd09a884818622080ff5f98942fb20acd7e0c916c2d5ea7ce6f7e173315384518f"),
			},
		},
		{
			"bigExpMod/nagydani-4-pow0x10001",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000200db34d0e438249c0ed685c949cc28776a05094e1c48691dc3f2dca5fc3356d2a0663bd376e4712839917eb9a19c670407e2c377a2de385a3ff3b52104f7f1f4e0c7bf7717fb913896693dc5edbb65b760ef1b00e42e9d8f9af17352385e1cd742c9b006c0f669995cb0bb21d28c0aced2892267637b6470d8cee0ab27fc5d42658f6e88240c31d6774aa60a7ebd25cd48b56d0da11209f1928e61005c6eb709f3e8e0aaf8d9b10f7d7e296d772264dc76897ccdddadc91efa91c1903b7232a9e4c3b941917b99a3bc0c26497dedc897c25750af60237aa67934a26a2bc491db3dcc677491944bc1f51d3e5d76b8d846a62db03dedd61ff508f91a56d71028125035c3a44cbb041497c83bf3e4ae2a9613a401cc721c547a2afa3b16a2969933d3626ed6d8a7428648f74122fd3f2a02a20758f7f693892c8fd798b39abac01d18506c45e71432639e9f9505719ee822f62ccbf47f6850f096ff77b5afaf4be7d772025791717dbe5abf9b3f40cff7d7aab6f67e38f62faf510747276e20a42127e7500c444f9ed92baf65ade9e836845e39c4316d9dce5f8e2c8083e2c0acbb95296e05e51aab13b6b8f53f06c9c4276e12b0671133218cc3ea907da3bd9a367096d9202128d14846cc2e20d56fc8473ecb07cecbfb8086919f3971926e7045b853d85a69d026195c70f9f7a823536e2a8f4b3e12e94d9b53a934353451094b81010001df3143a0057457d75e8c708b6337a6f5a4fd1a06727acf9fb93e2993c62f3378b37d56c85e7b1e00f0145ebf8e4095bd723166293c60b6ac1252291ef65823c9e040ddad14969b3b340a4ef714db093a587c37766d68b8d6b5016e741587e7e6bf7e763b44f0247e64bae30f994d248bfd20541a333e5b225ef6a61199e301738b1e688f70ec1d7fb892c183c95dc543c3e12adf8a5e8b9ca9d04f9445cced3ab256f29e998e69efaa633a7b60e1db5a867924ccab0a171d9d6e1098dfa15acde9553de599eaa56490c8f411e4985111f3d40bddfc5e301edb01547b01a886550a61158f7e2033c59707789bf7c854181d0c2e2a42a93cf09209747d7082e147eb8544de25c3eb14f2e35559ea0c0f5877f2f3fc92132c0ae9da4e45b2f6c866a224ea6d1f28c05320e287750fbc647368d41116e528014cc1852e5531d53e4af938374daba6cee4baa821ed07117253bb3601ddd00d59a3d7fb2ef1f5a2fbba7c429f0cf9a5b3462410fd833a69118f8be9c559b1000cc608fd877fb43f8e65c2d1302622b944462579056874b387208d90623fcdaf93920ca7a9e4ba64ea208758222ad868501cc2c345e2d3a5ea2a17e5069248138c8a79c0251185d29ee73e5afab5354769142d2bf0cb6712727aa6bf84a6245fcdae66e4938d84d1b9dd09a884818622080ff5f98942fb20acd7e0c916c2d5ea7ce6f7e173315384518f"),
			},
		},
		{
			"bigExpMod/nagydani-5-square",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000400c5a1611f8be90071a43db23cc2fe01871cc4c0e8ab5743f6378e4fef77f7f6db0095c0727e20225beb665645403453e325ad5f9aeb9ba99bf3c148f63f9c07cf4fe8847ad5242d6b7d4499f93bd47056ddab8f7dee878fc2314f344dbee2a7c41a5d3db91eff372c730c2fdd3a141a4b61999e36d549b9870cf2f4e632c4d5df5f024f81c028000073a0ed8847cfb0593d36a47142f578f05ccbe28c0c06aeb1b1da027794c48db880278f79ba78ae64eedfea3c07d10e0562668d839749dc95f40467d15cf65b9cfc52c7c4bcef1cda3596dd52631aac942f146c7cebd46065131699ce8385b0db1874336747ee020a5698a3d1a1082665721e769567f579830f9d259cec1a836845109c21cf6b25da572512bf3c42fd4b96e43895589042ab60dd41f497db96aec102087fe784165bb45f942859268fd2ff6c012d9d00c02ba83eace047cc5f7b2c392c2955c58a49f0338d6fc58749c9db2155522ac17914ec216ad87f12e0ee95574613942fa615898c4d9e8a3be68cd6afa4e7a003dedbdf8edfee31162b174f965b20ae752ad89c967b3068b6f722c16b354456ba8e280f987c08e0a52d40a2e8f3a59b94d590aeef01879eb7a90b3ee7d772c839c85519cbeaddc0c193ec4874a463b53fcaea3271d80ebfb39b33489365fc039ae549a17a9ff898eea2f4cb27b8dbee4c17b998438575b2b8d107e4a0d66ba7fca85b41a58a8d51f191a35c856dfbe8aef2b00048a694bbccff832d23c8ca7a7ff0b6c0b3011d00b97c86c0628444d267c951d9e4fb8f83e154b8f74fb51aa16535e498235c5597dac9606ed0be3173a3836baa4e7d756ffe1e2879b415d3846bccd538c05b847785699aefde3e305decb600cd8fb0e7d8de5efc26971a6ad4e6d7a2d91474f1023a0ac4b78dc937da0ce607a45974d2cac1c33a2631ff7fe6144a3b2e5cf98b531a9627dea92c1dc82204d09db0439b6a11dd64b484e1263aa45fd9539b6020b55e3baece3986a8bffc1003406348f5c61265099ed43a766ee4f93f5f9c5abbc32a0fd3ac2b35b87f9ec26037d88275bd7dd0a54474995ee34ed3727f3f97c48db544b1980193a4b76a8a3ddab3591ce527f16d91882e67f0103b5cda53f7da54d489fc4ac08b6ab358a5a04aa9daa16219d50bd672a7cb804ed769d218807544e5993f1c27427104b349906a0b654df0bf69328afd3013fbe430155339c39f236df5557bf92f1ded7ff609a8502f49064ec3d1dbfb6c15d3a4c11a4f8acd12278cbf68acd5709463d12e3338a6eddb8c112f199645e23154a8e60879d2a654e3ed9296aa28f134168619691cd2c6b9e2eba4438381676173fc63c2588a3c5910dc149cf3760f0aa9fa9c3f5faa9162b0bf1aac9dd32b706a60ef53cbdb394b6b40222b5bc80eea82ba8958386672564cae3794f977871ab62337cf02e30049201ec12937e7ce79d0f55d9c810e20acf52212aca1d3888949e0e4830aad88d804161230eb89d4d329cc83570fe257217d2119134048dd2ed167646975fc7d77136919a049ea74cf08ddd2b896890bb24a0ba18094a22baa351bf29ad96c66bbb1a598f2ca391749620e62d61c3561a7d3653ccc8892c7b99baaf76bf836e2991cb06d6bc0514568ff0d1ec8bb4b3d6984f5eaefb17d3ea2893722375d3ddb8e389a8eef7d7d198f8e687d6a513983df906099f9a2d23f4f9dec6f8ef2f11fc0a21fac45353b94e00486f5e17d386af42502d09db33cf0cf28310e049c07e88682aeeb00cb833c5174266e62407a57583f1f88b304b7c6e0c84bbe1c0fd423072d37a5bd0aacf764229e5c7cd02473460ba3645cd8e8ae144065bf02d0dd238593d8e230354f67e0b2f23012c23274f80e3ee31e35e2606a4a3f31d94ab755e6d163cff52cbb36b6d0cc67ffc512aeed1dce4d7a0d70ce82f2baba12e8d514dc92a056f994adfb17b5b9712bd5186f27a2fda1f7039c5df2c8587fdc62f5627580c13234b55be4df3056050e2d1ef3218f0dd66cb05265fe1acfb0989d8213f2c19d1735a7cf3fa65d88dad5af52dc2bba22b7abf46c3bc77b5091baab9e8f0ddc4d5e581037de91a9f8dcbc69309be29cc815cf19a20a7585b8b3073edf51fc9baeb3e509b97fa4ecfd621e0fd57bd61cac1b895c03248ff12bdbc57509250df3517e8a3fe1d776836b34ab352b973d932ef708b14f7418f9eceb1d87667e61e3e758649cb083f01b133d37ab2f5afa96d6c84bcacf4efc3851ad308c1e7d9113624fce29fab460ab9d2a48d92cdb281103a5250ad44cb2ff6e67ac670c02fdafb3e0f1353953d6d7d5646ca1568dea55275a050ec501b7c6250444f7219f1ba7521ba3b93d089727ca5f3bbe0d6c1300b423377004954c5628fdb65770b18ced5c9b23a4a5a6d6ef25fe01b4ce278de0bcc4ed86e28a0a68818ffa40970128cf2c38740e80037984428c1bd5113f40ff47512ee6f4e4d8f9b8e8e1b3040d2928d003bd1c1329dc885302fbce9fa81c23b4dc49c7c82d29b52957847898676c89aa5d32b5b0e1c0d5a2b79a19d67562f407f19425687971a957375879d90c5f57c857136c17106c9ab1b99d80e69c8c954ed386493368884b55c939b8d64d26f643e800c56f90c01079d7c534e3b2b7ae352cefd3016da55f6a85eb803b85e2304915fd2001f77c74e28746293c46e4f5f0fd49cf988aafd0026b8e7a3bab2da5cdce1ea26c2e29ec03f4807fac432662b2d6c060be1c7be0e5489de69d0a6e03a4b9117f9244b34a0f1ecba89884f781c6320412413a00c4980287409a2a78c2cd7e65cecebbe4ec1c28cac4dd95f6998e78fc6f1392384331c9436aa10e10e2bf8ad2c4eafbcf276aa7bae64b74428911b3269c749338b0fc5075ad"),
			},
		},
		{
			"bigExpMod/nagydani-5-pow0x10001",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x05"),
				common.FromHex("000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000400c5a1611f8be90071a43db23cc2fe01871cc4c0e8ab5743f6378e4fef77f7f6db0095c0727e20225beb665645403453e325ad5f9aeb9ba99bf3c148f63f9c07cf4fe8847ad5242d6b7d4499f93bd47056ddab8f7dee878fc2314f344dbee2a7c41a5d3db91eff372c730c2fdd3a141a4b61999e36d549b9870cf2f4e632c4d5df5f024f81c028000073a0ed8847cfb0593d36a47142f578f05ccbe28c0c06aeb1b1da027794c48db880278f79ba78ae64eedfea3c07d10e0562668d839749dc95f40467d15cf65b9cfc52c7c4bcef1cda3596dd52631aac942f146c7cebd46065131699ce8385b0db1874336747ee020a5698a3d1a1082665721e769567f579830f9d259cec1a836845109c21cf6b25da572512bf3c42fd4b96e43895589042ab60dd41f497db96aec102087fe784165bb45f942859268fd2ff6c012d9d00c02ba83eace047cc5f7b2c392c2955c58a49f0338d6fc58749c9db2155522ac17914ec216ad87f12e0ee95574613942fa615898c4d9e8a3be68cd6afa4e7a003dedbdf8edfee31162b174f965b20ae752ad89c967b3068b6f722c16b354456ba8e280f987c08e0a52d40a2e8f3a59b94d590aeef01879eb7a90b3ee7d772c839c85519cbeaddc0c193ec4874a463b53fcaea3271d80ebfb39b33489365fc039ae549a17a9ff898eea2f4cb27b8dbee4c17b998438575b2b8d107e4a0d66ba7fca85b41a58a8d51f191a35c856dfbe8aef2b00048a694bbccff832d23c8ca7a7ff0b6c0b3011d00b97c86c0628444d267c951d9e4fb8f83e154b8f74fb51aa16535e498235c5597dac9606ed0be3173a3836baa4e7d756ffe1e2879b415d3846bccd538c05b847785699aefde3e305decb600cd8fb0e7d8de5efc26971a6ad4e6d7a2d91474f1023a0ac4b78dc937da0ce607a45974d2cac1c33a2631ff7fe6144a3b2e5cf98b531a9627dea92c1dc82204d09db0439b6a11dd64b484e1263aa45fd9539b6020b55e3baece3986a8bffc1003406348f5c61265099ed43a766ee4f93f5f9c5abbc32a0fd3ac2b35b87f9ec26037d88275bd7dd0a54474995ee34ed3727f3f97c48db544b1980193a4b76a8a3ddab3591ce527f16d91882e67f0103b5cda53f7da54d489fc4ac08b6ab358a5a04aa9daa16219d50bd672a7cb804ed769d218807544e5993f1c27427104b349906a0b654df0bf69328afd3013fbe430155339c39f236df5557bf92f1ded7ff609a8502f49064ec3d1dbfb6c15d3a4c11a4f8acd12278cbf68acd5709463d12e3338a6eddb8c112f199645e23154a8e60879d2a654e3ed9296aa28f134168619691cd2c6b9e2eba4438381676173fc63c2588a3c5910dc149cf3760f0aa9fa9c3f5faa9162b0bf1aac9dd32b706a60ef53cbdb394b6b40222b5bc80eea82ba8958386672564cae3794f977871ab62337cf010001e30049201ec12937e7ce79d0f55d9c810e20acf52212aca1d3888949e0e4830aad88d804161230eb89d4d329cc83570fe257217d2119134048dd2ed167646975fc7d77136919a049ea74cf08ddd2b896890bb24a0ba18094a22baa351bf29ad96c66bbb1a598f2ca391749620e62d61c3561a7d3653ccc8892c7b99baaf76bf836e2991cb06d6bc0514568ff0d1ec8bb4b3d6984f5eaefb17d3ea2893722375d3ddb8e389a8eef7d7d198f8e687d6a513983df906099f9a2d23f4f9dec6f8ef2f11fc0a21fac45353b94e00486f5e17d386af42502d09db33cf0cf28310e049c07e88682aeeb00cb833c5174266e62407a57583f1f88b304b7c6e0c84bbe1c0fd423072d37a5bd0aacf764229e5c7cd02473460ba3645cd8e8ae144065bf02d0dd238593d8e230354f67e0b2f23012c23274f80e3ee31e35e2606a4a3f31d94ab755e6d163cff52cbb36b6d0cc67ffc512aeed1dce4d7a0d70ce82f2baba12e8d514dc92a056f994adfb17b5b9712bd5186f27a2fda1f7039c5df2c8587fdc62f5627580c13234b55be4df3056050e2d1ef3218f0dd66cb05265fe1acfb0989d8213f2c19d1735a7cf3fa65d88dad5af52dc2bba22b7abf46c3bc77b5091baab9e8f0ddc4d5e581037de91a9f8dcbc69309be29cc815cf19a20a7585b8b3073edf51fc9baeb3e509b97fa4ecfd621e0fd57bd61cac1b895c03248ff12bdbc57509250df3517e8a3fe1d776836b34ab352b973d932ef708b14f7418f9eceb1d87667e61e3e758649cb083f01b133d37ab2f5afa96d6c84bcacf4efc3851ad308c1e7d9113624fce29fab460ab9d2a48d92cdb281103a5250ad44cb2ff6e67ac670c02fdafb3e0f1353953d6d7d5646ca1568dea55275a050ec501b7c6250444f7219f1ba7521ba3b93d089727ca5f3bbe0d6c1300b423377004954c5628fdb65770b18ced5c9b23a4a5a6d6ef25fe01b4ce278de0bcc4ed86e28a0a68818ffa40970128cf2c38740e80037984428c1bd5113f40ff47512ee6f4e4d8f9b8e8e1b3040d2928d003bd1c1329dc885302fbce9fa81c23b4dc49c7c82d29b52957847898676c89aa5d32b5b0e1c0d5a2b79a19d67562f407f19425687971a957375879d90c5f57c857136c17106c9ab1b99d80e69c8c954ed386493368884b55c939b8d64d26f643e800c56f90c01079d7c534e3b2b7ae352cefd3016da55f6a85eb803b85e2304915fd2001f77c74e28746293c46e4f5f0fd49cf988aafd0026b8e7a3bab2da5cdce1ea26c2e29ec03f4807fac432662b2d6c060be1c7be0e5489de69d0a6e03a4b9117f9244b34a0f1ecba89884f781c6320412413a00c4980287409a2a78c2cd7e65cecebbe4ec1c28cac4dd95f6998e78fc6f1392384331c9436aa10e10e2bf8ad2c4eafbcf276aa7bae64b74428911b3269c749338b0fc5075ad"),
			},
		},

		// bn256Add
		{
			"bn256Add/chfast1",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x06"),
				common.FromHex("18b18acfb4c2c30276db5411368e7185b311dd124691610c5d3b74034e093dc9063c909c4720840cb5134cb9f59fa749755796819658d32efc0d288198f3726607c2b7f58a84bd6145f00c9c2bc0bb1a187f20ff2c92963a88019e7c6a014eed06614e20c147e940f2d70da3f74c9a17df361706a4485c742bd6788478fa17d7"),
			},
		},
		{
			"bn256Add/cdetrio3",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x06"),
				common.FromHex("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			},
		},

		// bn256ScalarMul
		{
			"bn256ScalarMul/chfast1",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x07"),
				common.FromHex("2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb721611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb20400000000000000000000000000000000000000000000000011138ce750fa15c2"),
			},
		},
		{
			"bn256ScalarMul/cdetrio2",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x07"),
				common.FromHex("1a87b0584ce92f4593d161480614f2989035225609f08058ccfa3d0f940febe31a2f3c951f6dadcc7ee9007dff81504b0fcd6d7cf59996efdc33d92bf7f9f8f630644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000000"),
			},
		},

		// bn256Paring
		{
			"bn256Paring/one_point",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x08"),
				common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa"),
			},
		},
		{
			"bn256Paring/two_point_match_2",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x08"),
				common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d"),
			},
		},
		{
			"bn256Paring/two_point_match_3",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x08"),
				common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002203e205db4f19b37b60121b83a7333706db86431c6d835849957ed8c3928ad7927dc7234fd11d3e8c36c59277c3e6f149d5cd3cfa9a62aee49f8130962b4b3b9195e8aa5b7827463722b8c153931579d3505566b4edf48d498e185f0509de15204bb53b8977e5f92a0bc372742c4830944a59b4fe6b1c0466e2a6dad122b5d2e030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd31a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa"),
			},
		},
		{
			"bn256Paring/ten_point_match_1",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x08"),
				common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d"),
			},
		},
		{
			"bn256Paring/ten_point_match_2",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x08"),
				common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002203e205db4f19b37b60121b83a7333706db86431c6d835849957ed8c3928ad7927dc7234fd11d3e8c36c59277c3e6f149d5cd3cfa9a62aee49f8130962b4b3b9195e8aa5b7827463722b8c153931579d3505566b4edf48d498e185f0509de15204bb53b8977e5f92a0bc372742c4830944a59b4fe6b1c0466e2a6dad122b5d2e030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd31a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002203e205db4f19b37b60121b83a7333706db86431c6d835849957ed8c3928ad7927dc7234fd11d3e8c36c59277c3e6f149d5cd3cfa9a62aee49f8130962b4b3b9195e8aa5b7827463722b8c153931579d3505566b4edf48d498e185f0509de15204bb53b8977e5f92a0bc372742c4830944a59b4fe6b1c0466e2a6dad122b5d2e030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd31a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002203e205db4f19b37b60121b83a7333706db86431c6d835849957ed8c3928ad7927dc7234fd11d3e8c36c59277c3e6f149d5cd3cfa9a62aee49f8130962b4b3b9195e8aa5b7827463722b8c153931579d3505566b4edf48d498e185f0509de15204bb53b8977e5f92a0bc372742c4830944a59b4fe6b1c0466e2a6dad122b5d2e030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd31a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002203e205db4f19b37b60121b83a7333706db86431c6d835849957ed8c3928ad7927dc7234fd11d3e8c36c59277c3e6f149d5cd3cfa9a62aee49f8130962b4b3b9195e8aa5b7827463722b8c153931579d3505566b4edf48d498e185f0509de15204bb53b8977e5f92a0bc372742c4830944a59b4fe6b1c0466e2a6dad122b5d2e030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd31a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002203e205db4f19b37b60121b83a7333706db86431c6d835849957ed8c3928ad7927dc7234fd11d3e8c36c59277c3e6f149d5cd3cfa9a62aee49f8130962b4b3b9195e8aa5b7827463722b8c153931579d3505566b4edf48d498e185f0509de15204bb53b8977e5f92a0bc372742c4830944a59b4fe6b1c0466e2a6dad122b5d2e030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd31a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa"),
			},
		},
		{
			"bn256Paring/ten_point_match_3",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x08"),
				common.FromHex("105456a333e6d636854f987ea7bb713dfd0ae8371a72aea313ae0c32c0bf10160cf031d41b41557f3e7e3ba0c51bebe5da8e6ecd855ec50fc87efcdeac168bcc0476be093a6d2b4bbf907172049874af11e1b6267606e00804d3ff0037ec57fd3010c68cb50161b7d1d96bb71edfec9880171954e56871abf3d93cc94d745fa114c059d74e5b6c4ec14ae5864ebe23a71781d86c29fb8fb6cce94f70d3de7a2101b33461f39d9e887dbb100f170a2345dde3c07e256d1dfa2b657ba5cd030427000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000021a2c3013d2ea92e13c800cde68ef56a294b883f6ac35d25f587c09b1b3c635f7290158a80cd3d66530f74dc94c94adb88f5cdb481acca997b6e60071f08a115f2f997f3dbd66a7afe07fe7862ce239edba9e05c5afff7f8a1259c9733b2dfbb929d1691530ca701b4a106054688728c9972c8512e9789e9567aae23e302ccd75"),
			},
		},

		// VmLog
		{
			"vmLog/1",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("48"),
			},
		},
		{
			"vmLog/13",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("48656c6c6f204b6c6179746e21"),
			},
		},
		{
			"vmLog/40",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("48484848484848484848484848484848484848484848484848484848484848484848484848484848"),
			},
		},
		{
			"vmLog/64",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"vmLog/128",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"vmLog/256",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"vmLog/1024",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("0102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},
		{
			"vmLog/2048",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x09"),
				common.FromHex("01020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132010203040506070809101112131415161718192021222324252627282930313201020304050607080910111213141516171819202122232425262728293031320102030405060708091011121314151617181920212223242526272829303132"),
			},
		},

		// feePayer
		{
			"feePayer",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0A"),
				common.FromHex(""),
			},
		},

		// validateSender
		{
			"validateSender/1",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderData1,
			},
		},
		{
			"validateSender/2",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_2,
			},
		},
		{
			"validateSender/3",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_3,
			},
		},
		{
			"validateSender/4",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_4,
			},
		},
		{
			"validateSender/5",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_5,
			},
		},
		{
			"validateSender/6",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_6,
			},
		},
		{
			"validateSender/7",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_7,
			},
		},
		{
			"validateSender/8",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_8,
			},
		},
		{
			"validateSender/9",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_9,
			},
		},
		{
			"validateSender/10",
			"precompiledContractTest",
			[]interface{}{
				loopCnt,
				common.HexToAddress("0x0B"),
				validateSenderMultisig10_10,
			},
		},
	}

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	assert.Equal(t, nil, err)
	amount := new(big.Int).SetUint64(0)

	for _, tc := range testcases {
		t.Run(tc.testName, func(t *testing.B) {
			if testing.Verbose() {
				fmt.Printf("----------------------testing %s...\n", tc.testName)
			}
			input := tc.input
			if tc.funcName == "BlockHash" {
				input = append(input, new(big.Int).SetUint64(bcdata.bc.CurrentBlock().NumberU64()-2))
			}

			// for i := 0; i < 1000; i++ {
			var txs types.Transactions

			// tc.input[1] = big.NewInt(int64(i) * 10000)

			data, err := abii.Pack(tc.funcName, input...)
			assert.Equal(t, nil, err)

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:    reservoir.Nonce,
				types.TxValueKeyFrom:     reservoir.Addr,
				types.TxValueKeyTo:       contractAddrs["OpCodeBenchmarkContract"],
				types.TxValueKeyAmount:   amount,
				types.TxValueKeyGasLimit: gasLimit,
				types.TxValueKeyGasPrice: gasPrice,
				types.TxValueKeyData:     data,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.Keys)
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			require.NoError(t, bcdata.GenABlockWithTransactions(accountMap, txs, prof))
			reservoir.Nonce += 1
			//}
		})
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}
