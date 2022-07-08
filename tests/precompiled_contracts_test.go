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
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ripemd160"
)

func TestPrecompiledContract(t *testing.T) {
	prof := profile.NewProfiler()

	if isCompilerAvailable() == false {
		if testing.Verbose() {
			fmt.Printf("TestFeePayerContract is skipped due to the lack of solc.")
		}
		return
	}

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

	// Deploy the contract.
	start = time.Now()
	filepath := "../contracts/precompiledContracts/precompiled.sol"
	contracts, err := deployContract(filepath, bcdata, accountMap, prof)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_deployContract", time.Now().Sub(start))

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	// PrecompiledContracts 0x01: ecrecover
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledEcrecover"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		hash := crypto.Keccak256Hash([]byte{1, 2, 3, 4})
		sig, err := crypto.Sign(hash[:], bcdata.privKeys[0])
		require.NoError(t, err)

		v := sig[crypto.RecoveryIDOffset] + 27
		var r, s [32]byte
		copy(r[:], sig[:32])
		copy(s[:], sig[32:64])

		// Make an input data for the contract call.
		data, err := abii.Pack("callEcrecover", hash, v, r, s)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var addr common.Address
		if err := abii.Unpack(&addr, "callEcrecover", ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, *bcdata.addrs[0], addr)
	}

	// PrecompiledContracts 0x02: sha256
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledSha256Hash"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		fn := "callSha256"

		dd := []byte{1, 2, 3, 4}
		hash := sha256.Sum256(dd)

		// Make an input data for the contract call.
		data, err := abii.Pack(fn, dd)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var addr [32]byte
		var bb [32]byte
		copy(bb[:], hash[:])
		if err := abii.Unpack(&addr, fn, ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, bb, addr)
	}

	// PrecompiledContracts 0x03: ripemd160
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledRipemd160Hash"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		fn := "callRipemd160"

		dd := []byte{1, 2, 3, 4}
		ripemd := ripemd160.New()
		ripemd.Write(dd)
		hash := ripemd.Sum(nil)

		// Make an input data for the contract call.
		data, err := abii.Pack(fn, dd)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var addr [32]byte
		var bb [32]byte
		copy(bb[:], hash[:])
		if err := abii.Unpack(&addr, fn, ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, bb, addr)
	}

	// PrecompiledContracts 0x04: datacopy (identity)
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledDatacopy"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		fn := "callDatacopy"

		dd := []byte{1, 2, 3, 4}

		// Make an input data for the contract call.
		data, err := abii.Pack(fn, dd)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var addr []byte
		if err := abii.Unpack(&addr, fn, ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, dd, addr)
	}

	// PrecompiledContracts 0x05: bigModExp
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledBigModExp"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		fn := "callBigModExp"

		var base, exp, mod [32]byte
		copy(base[:], hexutil.MustDecode("0x0000000000000000000000000000000462e4ded88953a39ce849a8a7fa163fa9"))
		copy(exp[:], hexutil.MustDecode("0x1f4a3123ff1223a1b0d040057af8a9fe70baa9258e0b959273ffc5718c6d4cc7"))
		copy(mod[:], hexutil.MustDecode("0x00000000000000000000000000077d29a9c710b7e616683f194f18c43b43b869"))

		// Make an input data for the contract call.
		data, err := abii.Pack(fn, base, exp, mod)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var result [32]byte
		if err := abii.Unpack(&result, fn, ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, hexutil.MustDecode("0x0000000000000000000000000002dc17ba6bb47224fccc9a5ece97e2f691b4aa"), result[:])
	}

	// PrecompiledContracts 0x06: bn256Add
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledBn256Add"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		fn := "callBn256Add"

		var ax, ay, bx, by [32]byte
		copy(ax[:], hexutil.MustDecode("0x17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa9"))
		copy(ay[:], hexutil.MustDecode("0x01e0559bacb160664764a357af8a9fe70baa9258e0b959273ffc5718c6d4cc7c"))
		copy(bx[:], hexutil.MustDecode("0x039730ea8dff1254c0fee9c0ea777d29a9c710b7e616683f194f18c43b43b869"))
		copy(by[:], hexutil.MustDecode("0x073a5ffcc6fc7a28c30723d6e58ce577356982d65b833a5a5c15bf9024b43d98"))

		// Make an input data for the contract call.
		data, err := abii.Pack(fn, ax, ay, bx, by)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var result [2][32]byte
		if err := abii.Unpack(&result, fn, ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, hexutil.MustDecode("0x15bf2bb17880144b5d1cd2b1f46eff9d617bffd1ca57c37fb5a49bd84e53cf66"), result[0][:])
		require.Equal(t, hexutil.MustDecode("0x049c797f9ce0d17083deb32b5e36f2ea2a212ee036598dd7624c168993d1355f"), result[1][:])
	}

	// PrecompiledContracts 0x07: bn256ScalarMul
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledBn256ScalarMul"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		fn := "callBn256ScalarMul"

		var ax, ay, scalar [32]byte
		copy(ax[:], hexutil.MustDecode("0x2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb7"))
		copy(ay[:], hexutil.MustDecode("0x21611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb204"))
		copy(scalar[:], hexutil.MustDecode("0x00000000000000000000000000000000000000000000000011138ce750fa15c2"))

		// Make an input data for the contract call.
		data, err := abii.Pack(fn, ax, ay, scalar)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var result [2][32]byte
		if err := abii.Unpack(&result, fn, ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, hexutil.MustDecode("0x070a8d6a982153cae4be29d434e8faef8a47b274a053f5a4ee2a6c9c13c31e5c"), result[0][:])
		require.Equal(t, hexutil.MustDecode("0x031b8ce914eba3a9ffb989f9cdd5b0f01943074bf4f0f315690ec3cec6981afc"), result[1][:])
	}

	// PrecompiledContracts 0x08: bn256Paring
	{
		c := contracts["../contracts/precompiledContracts/precompiled.sol:PrecompiledBn256Pairing"]
		abii, err := abi.JSON(strings.NewReader(c.abi))

		fn := "callBn256Pairing"

		x1 := hexutil.MustDecode("0x1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f59")
		y1 := hexutil.MustDecode("0x3034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41")
		x2 := hexutil.MustDecode("0x209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf7")
		y2 := hexutil.MustDecode("0x04bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a41678")
		x3 := hexutil.MustDecode("0x2bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d")
		y3 := hexutil.MustDecode("0x120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550")
		x4 := hexutil.MustDecode("0x111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c")
		y4 := hexutil.MustDecode("0x2032c61a830e3c17286de9462bf242fca2883585b93870a73853face6a6bf411")
		x5 := hexutil.MustDecode("0x198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2")
		y5 := hexutil.MustDecode("0x1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed")
		x6 := hexutil.MustDecode("0x090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b")
		y6 := hexutil.MustDecode("0x12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa")

		input := make([]byte, 32*2*6)
		copy(input[32*0:32*1], x1)
		copy(input[32*1:32*2], y1)
		copy(input[32*2:32*3], x2)
		copy(input[32*3:32*4], y2)
		copy(input[32*4:32*5], x3)
		copy(input[32*5:32*6], y3)
		copy(input[32*6:32*7], x4)
		copy(input[32*7:32*8], y4)
		copy(input[32*8:32*9], x5)
		copy(input[32*9:32*10], y5)
		copy(input[32*10:32*11], x6)
		copy(input[32*11:32*12], y6)

		// Make an input data for the contract call.
		data, err := abii.Pack(fn, input)
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
		require.Equal(t, nil, err)

		err = tx.Sign(signer, bcdata.privKeys[0])
		require.Equal(t, nil, err)

		err = tx.SignFeePayer(signer, bcdata.privKeys[1])
		require.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		require.Equal(t, nil, err)

		// Check the returned value.
		var result [32]byte
		if err := abii.Unpack(&result, fn, ret); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"), result[:])
	}
}
