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
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/stretchr/testify/assert"
)

// TestFeePayerContract tests a direct call of precompiled contract 0xa (feepayer).
// - It tests a contract `FeePayer` in fee_payer_test.sol.
// - It directly calls the precompiled contract 0xa.
func TestFeePayerContract(t *testing.T) {
	contractFunctions := []string{"GetFeePayerDirect", "GetFeePayer"}

	for _, c := range contractFunctions {
		t.Run(c, func(t *testing.T) {
			testFeePayerContract(t, c)
		})
	}
}

// TestFeePayerContract tests an indirect call of precompiled contract 0xa (feepayer).
// - It tests a contract `FeePayerIndirect` in fee_payer_indirect_test.sol.
// - It calls a deployed contract calling the precompiled contract 0xa.
// TODO-Klaytn-FeePayer: need more test cases for other calls such as delegatecall, etc.
func TestFeePayerContractIndirect(t *testing.T) {
	contractFunctions := []string{"TestCall"}

	for _, c := range contractFunctions {
		t.Run(c, func(t *testing.T) {
			testFeePayerContractIndirect(t, c)
		})
	}
}

// testFeePayerContract tests the fee-payer precompiled contract.
// 1. Deploy the FeePayer contract.
// 2. Make an input data for the contract call. The function name is given as a parameter.
// 3. Call the given function `fn`.
// 4. Check the returned value with the specified fee payer.
func testFeePayerContract(t *testing.T, fn string) {
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(2000, 4)
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

	// 1. Deploy the contract `FeePayer`.
	start = time.Now()
	filepath := "../contracts/feepayer/fee_payer_test.sol"
	contracts, err := deployContract(filepath, bcdata, accountMap, prof)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_deployContract", time.Now().Sub(start))

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)
	for _, c := range contracts {
		abii, err := abi.JSON(strings.NewReader(c.abi))

		// 2. Make an input data for the contract call. The function name is given as a parameter.
		data, err := abii.Pack(fn)
		if err != nil {
			t.Fatal(err)
		}

		n := accountMap.GetNonce(*bcdata.addrs[0])

		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    n,
			types.TxValueKeyFeePayer: *bcdata.addrs[1],
			types.TxValueKeyGasPrice: big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(5000000),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyTo:       c.address,
			types.TxValueKeyData:     data,
		})
		assert.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		assert.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		assert.Equal(t, nil, err)

		// 3. Call the given function `fn`.
		ret, err := callContract(bcdata, tx)
		assert.Equal(t, nil, err)

		// 4. Check the returned value with the specified fee payer.
		var feePayer common.Address
		if err := abii.Unpack(&feePayer, fn, ret); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *bcdata.addrs[1], feePayer)
	}
}

// testFeePayerContract tests an indirect call of the fee-payer precompiled contract.
// 1. Deploy the FeePayer contract.
// 2. Deploy the FeePayerIndirect contract.
// 3. Make an input data for the contract call. The function name is given as a parameter.
// 4. Call the given function `fn`.
// 5. Check the returned value with the specified fee payer.
func testFeePayerContractIndirect(t *testing.T, fn string) {
	prof := profile.NewProfiler()

	callee_path := "../contracts/feepayer/fee_payer_test.sol"
	caller_path := "../contracts/feepayer/fee_payer_indirect_test.sol"

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(2000, 4)
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

	// 1. Deploy the FeePayer contract.
	callee_contracts, err := deployContract(callee_path, bcdata, accountMap, prof)
	assert.Equal(t, nil, err)
	var callee deployedContract
	for _, v := range callee_contracts {
		callee = *v
		break
	}

	// 2. Deploy the FeePayerIndirect contract.
	start = time.Now()
	caller_contracts, err := deployContract(caller_path, bcdata, accountMap, prof)
	assert.Equal(t, nil, err)

	prof.Profile("main_deployContract", time.Now().Sub(start))

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)
	{
		var c deployedContract
		for k, v := range caller_contracts {
			if strings.Contains(k, "FeePayerIndirect") {
				c = *v
				break
			}
		}
		abii, err := abi.JSON(strings.NewReader(c.abi))

		// 3. Make an input data for the contract call. The function name is given as a parameter.
		data, err := abii.Pack(fn, callee.address)
		if err != nil {
			t.Fatal(err)
		}

		n := accountMap.GetNonce(*bcdata.addrs[0])

		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    n,
			types.TxValueKeyFeePayer: *bcdata.addrs[1],
			types.TxValueKeyGasPrice: big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(5000000),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyTo:       c.address,
			types.TxValueKeyData:     data,
		})
		assert.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		assert.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		assert.Equal(t, nil, err)

		// 4. Call the given function `fn`.
		ret, err := callContract(bcdata, tx)
		assert.Equal(t, nil, err)

		// 5. Check the returned value with the specified fee payer.
		var feePayer common.Address
		if err := abii.Unpack(&feePayer, fn, ret); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *bcdata.addrs[1], feePayer)
	}
}
