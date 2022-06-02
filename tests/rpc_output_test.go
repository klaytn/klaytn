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

// Basically, this test is disabled for the CI test. To run this,
// $ go test -run TestRPCOutput -tags RPCOutput
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
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/consensus/istanbul/backend"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// TestRPCOutput prints the RPC output in a JSON format by executing GetBlockWithConsensusInfoByNumber().
// This function makes two blocks.
// The first block contains the following:
// - TxTypeLegacyTransaction
// - TxTypeValueTransfer
// - TxTypeFeeDelegatedValueTransfer
// - TxTypeFeeDelegatedValueTransferWithRatio
// - TxTypeValueTransferMemo
// - TxTypeFeeDelegatedValueTransferMemo
// - TxTypeFeeDelegatedValueTransferMemoWithRatio
// - TxTypeAccountCreation
// The Second block contains the following:
// - TxTypeAccountUpdate
// - TxTypeFeeDelegatedAccountUpdate
// - TxTypeFeeDelegatedAccountUpdateWithRatio
// - TxTypeSmartContractDeploy
// - TxTypeFeeDelegatedSmartContractDeploy
// - TxTypeFeeDelegatedSmartContractDeployWithRatio
// - TxTypeSmartContractExecution
// - TxTypeFeeDelegatedSmartContractExecution
// - TxTypeFeeDelegatedSmartContractExecutionWithRatio
// - TxTypeCancel
// - TxTypeFeeDelegatedCancel
// - TxTypeFeeDelegatedCancelWithRatio
// - TxTypeChainDataAnchoring
func BenchmarkRPCOutput(t *testing.B) {
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

	reservoir2 := &TestAccountType{
		Addr:  *bcdata.addrs[1],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[1]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	assert.Equal(t, nil, err)

	// decoupled account
	decoupled, err := createDecoupledAccount("c64f2cd1196e2a1791365b00c4bc07ab8f047b73152e4617c6ed06ac221a4b0c",
		common.HexToAddress("0x75c3098be5e4b63fbac05838daaee378dd48098d"))
	assert.Equal(t, nil, err)

	// contract address
	contractAddr := common.Address{}

	if testing.Verbose() {
		fmt.Println("ChainID", (*hexutil.Big)(bcdata.bc.Config().ChainID))
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("reservoirPrvKey = ", (*hexutil.Big)(reservoir.Keys[0].D))
		fmt.Println("reservoir2Addr = ", reservoir2.Addr.String())
		fmt.Println("reservoir2PrvKey= ", (*hexutil.Big)(reservoir2.Keys[0].D))
		fmt.Println("anonAddr = ", anon.Addr.String())
		fmt.Println("anonPrvKey= ", (*hexutil.Big)(anon.Keys[0].D))
		fmt.Println("decoupledAddr = ", decoupled.Addr.String())
		fmt.Println("decoupledPrvKey = ", (*hexutil.Big)(decoupled.Keys[0].D))
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	var txs types.Transactions

	// TxTypeLegacyTransaction
	{
		amount := new(big.Int).SetUint64(100000000000)
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeValueTransfer
	{
		amount := new(big.Int).Mul(big.NewInt(10000), new(big.Int).SetUint64(params.KLAY))
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       decoupled.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedValueTransfer
	{
		amount := new(big.Int).Mul(big.NewInt(10000), new(big.Int).SetUint64(params.KLAY))
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       decoupled.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: reservoir2.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedValueTransferWithRatio
	{
		amount := new(big.Int).SetUint64(10000000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 decoupled.Addr,
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeValueTransferMemo
	{
		amount := new(big.Int).SetUint64(10000000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       decoupled.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte("hello"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedValueTransferMemo
	{
		amount := new(big.Int).SetUint64(10000000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       decoupled.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte("hello"),
			types.TxValueKeyFeePayer: reservoir2.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemo, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedValueTransferMemoWithRatio
	{
		amount := new(big.Int).SetUint64(10000000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 decoupled.Addr,
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               []byte("hello"),
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemoWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// Generate the first block!
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}

	// Make a new transaction slice.
	txs = make(types.Transactions, 0, 10)

	// TxTypeAccountUpdate
	{
		newKey, err := crypto.HexToECDSA("41bd2b972564206658eab115f26ff4db617e6eb39c81a557adc18d8305d2f867")
		if err != nil {
			t.Fatal(err)
		}

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&newKey.PublicKey),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		anon.Nonce += 1

		anon.Keys = []*ecdsa.PrivateKey{newKey}
	}

	// TxTypeFeeDelegatedAccountUpdate
	{
		newKey, err := crypto.HexToECDSA("41bd2b972564206658eab115f26ff4db617e6eb39c81a557adc18d8305d2f867")
		if err != nil {
			t.Fatal(err)
		}

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&newKey.PublicKey),
			types.TxValueKeyFeePayer:   reservoir.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		anon.Nonce += 1

		anon.Keys = []*ecdsa.PrivateKey{newKey}
	}

	// TxTypeFeeDelegatedAccountUpdateWithRatio
	{
		newKey, err := crypto.HexToECDSA("ed580f5bd71a2ee4dae5cb43e331b7d0318596e561e6add7844271ed94156b20")
		if err != nil {
			t.Fatal(err)
		}

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              anon.Nonce,
			types.TxValueKeyFrom:               anon.Addr,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyAccountKey:         accountkey.NewAccountKeyPublicWithValue(&newKey.PublicKey),
			types.TxValueKeyFeePayer:           reservoir.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(11),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		anon.Nonce += 1

		anon.Keys = []*ecdsa.PrivateKey{newKey}
	}

	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
	abiStr := `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`

	// TxTypeSmartContractDeploy
	{
		amount := new(big.Int).SetUint64(0)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      common.Big0,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyData:          common.FromHex(code),
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedSmartContractDeploy
	{
		amount := new(big.Int).SetUint64(0)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      common.Big0,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyData:          common.FromHex(code),
			types.TxValueKeyFeePayer:      reservoir2.Addr,
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedSmartContractDeployWithRatio
	{
		amount := new(big.Int).SetUint64(0)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 (*common.Address)(nil),
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           common.Big0,
			types.TxValueKeyHumanReadable:      false,
			types.TxValueKeyData:               common.FromHex(code),
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(33),
			types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeSmartContractExecution
	{
		amountToSend := new(big.Int).SetUint64(10)

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("reward", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       contractAddr,
			types.TxValueKeyAmount:   amountToSend,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     data,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedSmartContractExecution
	{
		amountToSend := new(big.Int).SetUint64(10)

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("reward", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       contractAddr,
			types.TxValueKeyAmount:   amountToSend,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: common.Big0,
			types.TxValueKeyData:     data,
			types.TxValueKeyFeePayer: reservoir2.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedSmartContractExecutionWithRatio
	{
		amountToSend := new(big.Int).SetUint64(10)

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("reward", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 contractAddr,
			types.TxValueKeyAmount:             amountToSend,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           common.Big0,
			types.TxValueKeyData:               data,
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(66),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeCancel
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeCancel, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedCancel
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: reservoir2.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeFeeDelegatedCancelWithRatio
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           big.NewInt(0),
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(88),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// TxTypeChainDataAnchoring
	{
		data := &types.AnchoringDataInternalType0{
			BlockHash:     common.HexToHash("0"),
			TxHash:        common.HexToHash("1"),
			ParentHash:    common.HexToHash("2"),
			ReceiptHash:   common.HexToHash("3"),
			StateRootHash: common.HexToHash("4"),
			BlockNumber:   big.NewInt(5),
			TxCount:       big.NewInt(6),
		}
		encodedCCTxData, err := rlp.EncodeToBytes(data)
		if err != nil {
			panic(err)
		}
		blockTxData := &types.AnchoringData{Type: 0, Data: encodedCCTxData}

		anchoredData, err := rlp.EncodeToBytes(blockTxData)
		if err != nil {
			panic(err)
		}

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:        reservoir.Nonce,
			types.TxValueKeyFrom:         reservoir.Addr,
			types.TxValueKeyGasLimit:     gasLimit,
			types.TxValueKeyGasPrice:     gasPrice,
			types.TxValueKeyAnchoredData: anchoredData,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeChainDataAnchoring, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		reservoir.Nonce += 1
	}

	// Generate the second block!
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}

	// Get APIExtension to execute GetBlockWithConsensusInfoByNumber.
	apis := bcdata.bc.Engine().APIs(bcdata.bc)
	apiExtension, ok := apis[1].Service.(*backend.APIExtension)
	if !ok {
		// checkout the code `consensus/istanbul/backend/engine.go` if it fails.
		t.Fatalf("APIExetension is not the second item of apis. check out the code!")
	}

	// Print the JSON output of the first block.
	blkNum := rpc.BlockNumber(1)
	out, err := apiExtension.GetBlockWithConsensusInfoByNumber(&blkNum)
	b, _ := json.MarshalIndent(out, "", "\t")
	fmt.Println(string(b))

	// Print the JSON output of the second block.
	blkNum = rpc.BlockNumber(2)
	out, err = apiExtension.GetBlockWithConsensusInfoByNumber(&blkNum)
	b, _ = json.MarshalIndent(out, "", "\t")
	fmt.Println(string(b))

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}
