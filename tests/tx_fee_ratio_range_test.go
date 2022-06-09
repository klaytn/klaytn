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
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

// TestTxFeeRatioRange checks the range of the fee ratio.
func TestTxFeeRatioRange(t *testing.T) {
	testCases := []struct {
		feeRatio types.FeeRatio
		expected error
	}{
		{0, kerrors.ErrFeeRatioOutOfRange},
		{1, nil},
		{99, nil},
		{100, kerrors.ErrFeeRatioOutOfRange},
	}

	for _, test := range testCases {
		name := fmt.Sprintf("%d", test.feeRatio)
		t.Run(name, func(t *testing.T) {
			testTxFeeRatioRange(t, test.feeRatio, test.expected)
		})
	}
}

func testTxFeeRatioRange(t *testing.T, feeRatio types.FeeRatio, expected error) {
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

	reservoir3 := &TestAccountType{
		Addr:  *bcdata.addrs[2],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[2]},
		Nonce: uint64(0),
	}

	contract, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e989")
	assert.Equal(t, nil, err)

	code := common.FromHex("0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029")

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 0. Initial setting: deploy a smart contract to execute later
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyData:          code,
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contract.Addr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 1. TxTypeFeeDelegatedValueTransferWithRatio
	{
		amount := new(big.Int).SetUint64(10000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 reservoir3.Addr,
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: feeRatio,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		err = txpool.AddRemote(tx)
		assert.Equal(t, expected, err)

		reservoir.Nonce += 1
	}

	// 2. TxTypeFeeDelegatedValueTransferMemoWithRatio
	{
		amount := new(big.Int).SetUint64(10000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 reservoir3.Addr,
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               []byte("hello"),
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: feeRatio,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemoWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		err = txpool.AddRemote(tx)
		assert.Equal(t, expected, err)

		reservoir.Nonce += 1
	}

	// 3. TxTypeFeeDelegatedAccountUpdateWithRatio
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyAccountKey:         accountkey.NewAccountKeyPublicWithValue(&reservoir.Keys[0].PublicKey),
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: feeRatio,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		err = txpool.AddRemote(tx)
		assert.Equal(t, expected, err)

		reservoir.Nonce += 1
	}

	// 4, TxTypeFeeDelegatedSmartContractDeployWithRatio
	{
		amount := new(big.Int).SetUint64(10000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 (*common.Address)(nil),
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyHumanReadable:      false,
			types.TxValueKeyData:               []byte{0x80},
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: feeRatio,
			types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		err = txpool.AddRemote(tx)
		assert.Equal(t, expected, err)

		reservoir.Nonce += 1
	}

	// 5. TxTypeFeeDelegatedSmartContractExecutionWithRatio
	{
		amount := new(big.Int).SetUint64(10000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 contract.Addr,
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               []byte{0x80},
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: feeRatio,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		err = txpool.AddRemote(tx)
		assert.Equal(t, expected, err)

		reservoir.Nonce += 1
	}

	// 6. TxTypeFeeDelegatedCancelWithRatio
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: feeRatio,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		err = txpool.AddRemote(tx)
		assert.Equal(t, expected, err)

		reservoir.Nonce += 1
	}
}
