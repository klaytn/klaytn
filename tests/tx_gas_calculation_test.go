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
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

var code = "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

type TestAccount interface {
	GetAddr() common.Address
	GetTxKeys() []*ecdsa.PrivateKey
	GetUpdateKeys() []*ecdsa.PrivateKey
	GetFeeKeys() []*ecdsa.PrivateKey
	GetNonce() uint64
	GetAccKey() accountkey.AccountKey
	GetValidationGas(r accountkey.RoleType) uint64
	AddNonce()
	SetNonce(uint64)
	SetAddr(common.Address)
}

type genTransaction func(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64)

func TestGasCalculation(t *testing.T) {
	testFunctions := []struct {
		Name  string
		genTx genTransaction
	}{
		{"LegacyTransaction", genLegacyTransaction},
		{"AccessListTransaction", genAccessListTransaction},
		{"DynamicFeeTransaction", genDynamicFeeTransaction},

		{"ValueTransfer", genValueTransfer},
		{"ValueTransferWithMemo", genValueTransferWithMemo},
		{"AccountUpdate", genAccountUpdate},
		{"SmartContractDeploy", genSmartContractDeploy},
		{"SmartContractExecution", genSmartContractExecution},
		{"Cancel", genCancel},
		{"ChainDataAnchoring", genChainDataAnchoring},

		{"FeeDelegatedValueTransfer", genFeeDelegatedValueTransfer},
		{"FeeDelegatedValueTransferWithMemo", genFeeDelegatedValueTransferWithMemo},
		{"FeeDelegatedAccountUpdate", genFeeDelegatedAccountUpdate},
		{"FeeDelegatedSmartContractDeploy", genFeeDelegatedSmartContractDeploy},
		{"FeeDelegatedSmartContractExecution", genFeeDelegatedSmartContractExecution},
		{"FeeDelegatedCancel", genFeeDelegatedCancel},
		{"FeeDelegatedChainDataAnchoring", genFeeDelegatedChainDataAnchoring},

		{"FeeDelegatedWithRatioValueTransfer", genFeeDelegatedWithRatioValueTransfer},
		{"FeeDelegatedWithRatioValueTransferWithMemo", genFeeDelegatedWithRatioValueTransferWithMemo},
		{"FeeDelegatedWithRatioAccountUpdate", genFeeDelegatedWithRatioAccountUpdate},
		{"FeeDelegatedWithRatioSmartContractDeploy", genFeeDelegatedWithRatioSmartContractDeploy},
		{"FeeDelegatedWithRatioSmartContractExecution", genFeeDelegatedWithRatioSmartContractExecution},
		{"FeeDelegatedWithRatioCancel", genFeeDelegatedWithRatioCancel},
		{"FeeDelegatedWithRatioChainDataAnchoring", genFeeDelegatedWithRatioChainDataAnchoring},
	}

	accountTypes := []struct {
		Type    string
		account TestAccount
	}{
		{"KlaytnLegacy", genKlaytnLegacyAccount(t)},
		{"Public", genPublicAccount(t)},
		{"MultiSig", genMultiSigAccount(t)},
		{"RoleBasedWithPublic", genRoleBasedWithPublicAccount(t)},
		{"RoleBasedWithMultiSig", genRoleBasedWithMultiSigAccount(t)},
	}

	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	assert.Equal(t, nil, err)
	bcdata.bc.Config().IstanbulCompatibleBlock = big.NewInt(0)
	bcdata.bc.Config().LondonCompatibleBlock = big.NewInt(0)
	bcdata.bc.Config().EthTxTypeCompatibleBlock = big.NewInt(0)
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
	var reservoir TestAccount
	reservoir = &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// Preparing step. Send KLAY to a KlaytnAcount.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.GetNonce(),
			accountTypes[0].account.GetAddr(), amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.GetTxKeys())
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.AddNonce()
	}

	// Preparing step. Send KLAY to KlaytnAcounts.
	for i := 1; i < len(accountTypes); i++ {
		// create an account which account key will be replaced to one of account key types.
		anon, err := createAnonymousAccount(getRandomPrivateKeyString(t))
		assert.Equal(t, nil, err)

		{
			var txs types.Transactions

			amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
			tx := types.NewTransaction(reservoir.GetNonce(),
				anon.Addr, amount, gasLimit, gasPrice, []byte{})

			err := tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)
			txs = append(txs, tx)

			if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
				t.Fatal(err)
			}
			reservoir.AddNonce()
		}

		// update the account's key
		{
			var txs types.Transactions

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:      anon.Nonce,
				types.TxValueKeyFrom:       anon.Addr,
				types.TxValueKeyGasLimit:   gasLimit,
				types.TxValueKeyGasPrice:   gasPrice,
				types.TxValueKeyAccountKey: accountTypes[i].account.GetAccKey(),
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, anon.Keys)
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
				t.Fatal(err)
			}
			anon.AddNonce()
		}

		accountTypes[i].account.SetAddr(anon.Addr)
		accountTypes[i].account.SetNonce(anon.Nonce)
	}

	// For smart contract
	contract, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e989")
	contract.Addr = common.Address{}

	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.GetNonce(),
			types.TxValueKeyFrom:          reservoir.GetAddr(),
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyData:          common.FromHex(code),
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.GetTxKeys())
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contract.Addr = crypto.CreateAddress(reservoir.GetAddr(), reservoir.GetNonce())

		reservoir.AddNonce()
	}

	for _, f := range testFunctions {
		for _, sender := range accountTypes {
			toAccount := reservoir
			senderRole := accountkey.RoleTransaction

			// LegacyTransaction can be used only by the KlaytnAccount with AccountKeyLegacy.
			if sender.Type != "KlaytnLegacy" && (strings.Contains(f.Name, "Legacy") || strings.Contains(f.Name, "Access") || strings.Contains(f.Name, "Dynamic")) {
				continue
			}

			if strings.Contains(f.Name, "AccountUpdate") {
				senderRole = accountkey.RoleAccountUpdate
			}

			// Set contract's address with SmartContractExecution
			if strings.Contains(f.Name, "SmartContractExecution") {
				toAccount = contract
			}

			if !strings.Contains(f.Name, "FeeDelegated") {
				// For NonFeeDelegated Transactions
				Name := f.Name + "/" + sender.Type + "Sender"
				t.Run(Name, func(t *testing.T) {
					tx, intrinsic := f.genTx(t, signer, sender.account, toAccount, nil, gasPrice)
					acocuntValidationGas := sender.account.GetValidationGas(senderRole)
					testGasValidation(t, bcdata, tx, intrinsic+acocuntValidationGas)
				})
			} else {
				// For FeeDelegated(WithRatio) Transactions
				for _, payer := range accountTypes {
					Name := f.Name + "/" + sender.Type + "Sender/" + payer.Type + "Payer"
					t.Run(Name, func(t *testing.T) {
						tx, intrinsic := f.genTx(t, signer, sender.account, toAccount, payer.account, gasPrice)
						acocuntsValidationGas := sender.account.GetValidationGas(senderRole) + payer.account.GetValidationGas(accountkey.RoleFeePayer)
						testGasValidation(t, bcdata, tx, intrinsic+acocuntsValidationGas)
					})
				}
			}

		}
	}
}

func testGasValidation(t *testing.T, bcdata *BCData, tx *types.Transaction, validationGas uint64) {
	receipt, gas, err := applyTransaction(t, bcdata, tx)
	assert.Equal(t, nil, err)

	assert.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)

	assert.Equal(t, validationGas, gas)
}

func genLegacyTransaction(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	intrinsic := getIntrinsicGas(types.TxTypeLegacyTransaction)
	amount := big.NewInt(100000)
	tx := types.NewTransaction(from.GetNonce(), to.GetAddr(), amount, gasLimit, gasPrice, []byte{})

	err := tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genAccessListTransaction(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForAccessListTransaction(from, to, gasPrice, types.TxTypeEthereumAccessList)
	tx, err := types.NewTransactionWithMap(types.TxTypeEthereumAccessList, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genDynamicFeeTransaction(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForDynamicFeeTransaction(from, to, gasPrice, types.TxTypeEthereumDynamicFee)
	tx, err := types.NewTransactionWithMap(types.TxTypeEthereumDynamicFee, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genValueTransfer(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForValueTransfer(from, to, gasPrice, types.TxTypeValueTransfer)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedValueTransfer(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForValueTransfer(from, to, gasPrice, types.TxTypeFeeDelegatedValueTransfer)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedWithRatioValueTransfer(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForValueTransfer(from, to, gasPrice, types.TxTypeFeeDelegatedValueTransferWithRatio)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genValueTransferWithMemo(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, gasPayloadWithGas := genMapForValueTransferWithMemo(from, to, gasPrice, types.TxTypeValueTransferMemo)

	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, gasPayloadWithGas
}

func genFeeDelegatedValueTransferWithMemo(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, gasPayloadWithGas := genMapForValueTransferWithMemo(from, to, gasPrice, types.TxTypeFeeDelegatedValueTransferMemo)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemo, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, gasPayloadWithGas
}

func genFeeDelegatedWithRatioValueTransferWithMemo(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, gasPayloadWithGas := genMapForValueTransferWithMemo(from, to, gasPrice, types.TxTypeFeeDelegatedValueTransferMemoWithRatio)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemoWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, gasPayloadWithGas
}

func genAccountUpdate(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	newAccount, gasKey, _ := genNewAccountWithGas(t, from)

	values, intrinsic := genMapForUpdate(from, to, gasPrice, newAccount.GetAccKey(), types.TxTypeAccountUpdate)

	tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetUpdateKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic + gasKey
}

func genFeeDelegatedAccountUpdate(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	newAccount, gasKey, _ := genNewAccountWithGas(t, from)

	values, intrinsic := genMapForUpdate(from, to, gasPrice, newAccount.GetAccKey(), types.TxTypeFeeDelegatedAccountUpdate)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetUpdateKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic + gasKey
}

func genFeeDelegatedWithRatioAccountUpdate(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	newAccount, gasKey, _ := genNewAccountWithGas(t, from)

	values, intrinsic := genMapForUpdate(from, to, gasPrice, newAccount.GetAccKey(), types.TxTypeFeeDelegatedAccountUpdateWithRatio)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetUpdateKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic + gasKey
}

func genSmartContractDeploy(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsicGas := genMapForDeploy(from, to, gasPrice, types.TxTypeSmartContractDeploy)
	if values == nil {
		t.Fatalf("failed to genMapForDeploy")
	}

	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsicGas
}

func genFeeDelegatedSmartContractDeploy(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsicGas := genMapForDeploy(from, to, gasPrice, types.TxTypeFeeDelegatedSmartContractDeploy)
	if values == nil {
		t.Fatalf("failed to genMapForDeploy")
	}

	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsicGas
}

func genFeeDelegatedWithRatioSmartContractDeploy(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsicGas := genMapForDeploy(from, to, gasPrice, types.TxTypeFeeDelegatedSmartContractDeployWithRatio)
	if values == nil {
		t.Fatalf("failed to genMapForDeploy")
	}

	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsicGas
}

func genSmartContractExecution(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsicGas := genMapForExecution(from, to, gasPrice, types.TxTypeSmartContractExecution)
	if values == nil {
		t.Fatalf("failed to genMapForExecution")
	}

	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)

	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsicGas
}

func genFeeDelegatedSmartContractExecution(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsicGas := genMapForExecution(from, to, gasPrice, types.TxTypeFeeDelegatedSmartContractExecution)
	if values == nil {
		t.Fatalf("failed to genMapForExecution")
	}

	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, values)

	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsicGas
}

func genFeeDelegatedWithRatioSmartContractExecution(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsicGas := genMapForExecution(from, to, gasPrice, types.TxTypeFeeDelegatedSmartContractExecutionWithRatio)
	if values == nil {
		t.Fatalf("failed to genMapForExecution")
	}

	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, values)

	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsicGas
}

func genCancel(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForCancel(from, gasPrice, types.TxTypeCancel)

	tx, err := types.NewTransactionWithMap(types.TxTypeCancel, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedCancel(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForCancel(from, gasPrice, types.TxTypeFeeDelegatedCancel)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedWithRatioCancel(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForCancel(from, gasPrice, types.TxTypeFeeDelegatedCancelWithRatio)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genChainDataAnchoring(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForChainDataAnchoring(from, gasPrice, types.TxTypeChainDataAnchoring)

	tx, err := types.NewTransactionWithMap(types.TxTypeChainDataAnchoring, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedChainDataAnchoring(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForChainDataAnchoring(from, gasPrice, types.TxTypeFeeDelegatedChainDataAnchoring)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedChainDataAnchoring, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedWithRatioChainDataAnchoring(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForChainDataAnchoring(from, gasPrice, types.TxTypeFeeDelegatedChainDataAnchoringWithRatio)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedChainDataAnchoringWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

// Generate map functions
func genMapForLegacyTransaction(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)
	amount := big.NewInt(100000)
	data := []byte{0x11, 0x22}
	gasPayload := uint64(len(data)) * params.TxDataGas

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    from.GetNonce(),
		types.TxValueKeyTo:       to.GetAddr(),
		types.TxValueKeyAmount:   amount,
		types.TxValueKeyData:     data,
		types.TxValueKeyGasLimit: gasLimit,
		types.TxValueKeyGasPrice: gasPrice,
	}
	return values, intrinsic + gasPayload
}

func genMapForAccessListTransaction(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)
	amount := big.NewInt(100000)
	data := []byte{0x11, 0x22}
	gasPayload := uint64(len(data)) * params.TxDataGas
	accessList := types.AccessList{{Address: common.HexToAddress("0x0000000000000000000000000000000000000001"), StorageKeys: []common.Hash{{0}}}}
	toAddress := to.GetAddr()

	gasPayload += uint64(len(accessList)) * params.TxAccessListAddressGas
	gasPayload += uint64(accessList.StorageKeys()) * params.TxAccessListStorageKeyGas

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      from.GetNonce(),
		types.TxValueKeyTo:         &toAddress,
		types.TxValueKeyAmount:     amount,
		types.TxValueKeyData:       data,
		types.TxValueKeyGasLimit:   gasLimit,
		types.TxValueKeyGasPrice:   gasPrice,
		types.TxValueKeyAccessList: accessList,
		types.TxValueKeyChainID:    big.NewInt(1),
	}
	return values, intrinsic + gasPayload
}

func genMapForDynamicFeeTransaction(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)
	amount := big.NewInt(100000)
	data := []byte{0x11, 0x22}
	gasPayload := uint64(len(data)) * params.TxDataGas
	accessList := types.AccessList{{Address: common.HexToAddress("0x0000000000000000000000000000000000000001"), StorageKeys: []common.Hash{{0}}}}
	toAddress := to.GetAddr()

	gasPayload += uint64(len(accessList)) * params.TxAccessListAddressGas
	gasPayload += uint64(accessList.StorageKeys()) * params.TxAccessListStorageKeyGas

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      from.GetNonce(),
		types.TxValueKeyTo:         &toAddress,
		types.TxValueKeyAmount:     amount,
		types.TxValueKeyData:       data,
		types.TxValueKeyGasLimit:   gasLimit,
		types.TxValueKeyGasFeeCap:  gasPrice,
		types.TxValueKeyGasTipCap:  gasPrice,
		types.TxValueKeyAccessList: accessList,
		types.TxValueKeyChainID:    big.NewInt(1),
	}
	return values, intrinsic + gasPayload
}

func genMapForValueTransfer(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)
	amount := big.NewInt(100000)

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    from.GetNonce(),
		types.TxValueKeyFrom:     from.GetAddr(),
		types.TxValueKeyTo:       to.GetAddr(),
		types.TxValueKeyAmount:   amount,
		types.TxValueKeyGasLimit: gasLimit,
		types.TxValueKeyGasPrice: gasPrice,
	}
	return values, intrinsic
}

func genMapForValueTransferWithMemo(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)

	nonZeroData := []byte{1, 2, 3, 4}
	zeroData := []byte{0, 0, 0, 0}

	data := append(nonZeroData, zeroData...)

	amount := big.NewInt(100000)

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    from.GetNonce(),
		types.TxValueKeyFrom:     from.GetAddr(),
		types.TxValueKeyTo:       to.GetAddr(),
		types.TxValueKeyAmount:   amount,
		types.TxValueKeyGasLimit: gasLimit,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
	}

	gasPayload := uint64(len(data)) * params.TxDataGas

	return values, intrinsic + gasPayload
}

func genMapForCreate(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)
	amount := big.NewInt(0)

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         from.GetNonce(),
		types.TxValueKeyFrom:          from.GetAddr(),
		types.TxValueKeyTo:            to.GetAddr(),
		types.TxValueKeyAmount:        amount,
		types.TxValueKeyGasLimit:      gasLimit,
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyAccountKey:    to.GetAccKey(),
	}
	return values, intrinsic + uint64(len(to.GetTxKeys()))*params.TxAccountCreationGasPerKey
}

func genMapForUpdate(from TestAccount, to TestAccount, gasPrice *big.Int, newKeys accountkey.AccountKey, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      from.GetNonce(),
		types.TxValueKeyFrom:       from.GetAddr(),
		types.TxValueKeyGasLimit:   gasLimit,
		types.TxValueKeyGasPrice:   gasPrice,
		types.TxValueKeyAccountKey: newKeys,
	}
	return values, intrinsic
}

func genMapForDeploy(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	amount := new(big.Int).SetUint64(0)
	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         from.GetNonce(),
		types.TxValueKeyAmount:        amount,
		types.TxValueKeyGasLimit:      gasLimit,
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyFrom:          from.GetAddr(),
		types.TxValueKeyData:          common.FromHex(code),
		types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		types.TxValueKeyTo:            (*common.Address)(nil),
	}

	intrinsicGas := getIntrinsicGas(txType)
	intrinsicGas += uint64(0x175fd)

	gasPayloadWithGas, err := types.IntrinsicGasPayload(intrinsicGas, common.FromHex(code))
	if err != nil {
		return nil, 0
	}

	return values, gasPayloadWithGas
}

func genMapForExecution(from TestAccount, to TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	abiStr := `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`
	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		return nil, 0
	}

	data, err := abii.Pack("reward", to.GetAddr())
	if err != nil {
		return nil, 0
	}

	amount := new(big.Int).SetUint64(10)

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    from.GetNonce(),
		types.TxValueKeyFrom:     from.GetAddr(),
		types.TxValueKeyTo:       to.GetAddr(),
		types.TxValueKeyAmount:   amount,
		types.TxValueKeyGasLimit: gasLimit,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
	}

	intrinsicGas := getIntrinsicGas(txType)
	intrinsicGas += uint64(0x9ec4)

	gasPayloadWithGas, err := types.IntrinsicGasPayload(intrinsicGas, data)
	if err != nil {
		return nil, 0
	}

	return values, gasPayloadWithGas
}

func genMapForCancel(from TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    from.GetNonce(),
		types.TxValueKeyFrom:     from.GetAddr(),
		types.TxValueKeyGasLimit: gasLimit,
		types.TxValueKeyGasPrice: gasPrice,
	}
	return values, intrinsic
}

func genMapForChainDataAnchoring(from TestAccount, gasPrice *big.Int, txType types.TxType) (map[types.TxValueKeyType]interface{}, uint64) {
	intrinsic := getIntrinsicGas(txType)
	data := []byte{0x11, 0x22}
	gasPayload := uint64(len(data)) * params.TxDataGas

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:        from.GetNonce(),
		types.TxValueKeyFrom:         from.GetAddr(),
		types.TxValueKeyGasLimit:     gasLimit,
		types.TxValueKeyGasPrice:     gasPrice,
		types.TxValueKeyAnchoredData: data,
	}
	return values, intrinsic + gasPayload
}

func genKlaytnLegacyAccount(t *testing.T) TestAccount {
	// For KlaytnLegacy
	klaytnLegacy, err := createAnonymousAccount(getRandomPrivateKeyString(t))
	assert.Equal(t, nil, err)

	return klaytnLegacy
}

func genPublicAccount(t *testing.T) TestAccount {
	// For AccountKeyPublic
	var newAddress common.Address
	newAddress.SetBytesFromFront([]byte(getRandomPrivateKeyString(t)))

	publicAccount, err := createDecoupledAccount(getRandomPrivateKeyString(t), newAddress)
	assert.Equal(t, nil, err)

	return publicAccount
}

func genMultiSigAccount(t *testing.T) TestAccount {
	// For AccountKeyWeightedMultiSig
	var newAddress common.Address
	newAddress.SetBytesFromFront([]byte(getRandomPrivateKeyString(t)))

	multisigAccount, err := createMultisigAccount(uint(2),
		[]uint{1, 1, 1},
		[]string{getRandomPrivateKeyString(t), getRandomPrivateKeyString(t), getRandomPrivateKeyString(t)}, newAddress)
	assert.Equal(t, nil, err)

	return multisigAccount
}

func genRoleBasedWithPublicAccount(t *testing.T) TestAccount {
	// For AccountKeyRoleBased With AccountKeyPublic
	var roleBasedWithPublicAddr common.Address
	roleBasedWithPublicAddr.SetBytesFromFront([]byte(getRandomPrivateKeyString(t)))

	roleBasedWithPublic, err := createRoleBasedAccountWithAccountKeyPublic(
		[]string{getRandomPrivateKeyString(t), getRandomPrivateKeyString(t), getRandomPrivateKeyString(t)}, roleBasedWithPublicAddr)
	assert.Equal(t, nil, err)

	return roleBasedWithPublic
}

func genRoleBasedWithMultiSigAccount(t *testing.T) TestAccount {
	// For AccountKeyRoleBased With AccountKeyWeightedMultiSig
	p := genMultiSigParamForRoleBased(t)

	var roleBasedWithMultiSigAddr common.Address
	roleBasedWithMultiSigAddr.SetBytesFromFront([]byte(getRandomPrivateKeyString(t)))

	roleBasedWithMultiSig, err := createRoleBasedAccountWithAccountKeyWeightedMultiSig(
		[]TestCreateMultisigAccountParam{p[0], p[1], p[2]}, roleBasedWithMultiSigAddr)
	assert.Equal(t, nil, err)

	return roleBasedWithMultiSig
}

// Generate new Account functions for testing AccountUpdate
func genNewAccountWithGas(t *testing.T, testAccount TestAccount) (TestAccount, uint64, bool) {
	var newAccount TestAccount
	gas := uint64(0)
	readableGas := uint64(0)
	readable := false

	// AccountKeyLegacy
	if testAccount.GetAccKey() == nil || testAccount.GetAccKey().Type() == accountkey.AccountKeyTypeLegacy {
		anon, err := createAnonymousAccount(getRandomPrivateKeyString(t))
		assert.Equal(t, nil, err)

		return anon, gas + readableGas, readable
	}

	// newAccount
	newAccount, err := createAnonymousAccount(getRandomPrivateKeyString(t))
	assert.Equal(t, nil, err)

	switch testAccount.GetAccKey().Type() {
	case accountkey.AccountKeyTypePublic:
		publicAccount, err := createDecoupledAccount(getRandomPrivateKeyString(t), newAccount.GetAddr())
		assert.Equal(t, nil, err)

		newAccount = publicAccount
		gas += uint64(len(newAccount.GetTxKeys())) * params.TxAccountCreationGasPerKey
	case accountkey.AccountKeyTypeWeightedMultiSig:
		multisigAccount, err := createMultisigAccount(uint(2), []uint{1, 1, 1},
			[]string{getRandomPrivateKeyString(t), getRandomPrivateKeyString(t), getRandomPrivateKeyString(t)}, newAccount.GetAddr())
		assert.Equal(t, nil, err)

		newAccount = multisigAccount
		gas += uint64(len(newAccount.GetTxKeys())) * params.TxAccountCreationGasPerKey
	case accountkey.AccountKeyTypeRoleBased:
		p := genMultiSigParamForRoleBased(t)

		newRoleBasedWithMultiSig, err := createRoleBasedAccountWithAccountKeyWeightedMultiSig(
			[]TestCreateMultisigAccountParam{p[0], p[1], p[2]}, newAccount.GetAddr())
		assert.Equal(t, nil, err)

		newAccount = newRoleBasedWithMultiSig
		gas = uint64(len(newAccount.GetTxKeys())) * params.TxAccountCreationGasPerKey
		gas += uint64(len(newAccount.GetUpdateKeys())) * params.TxAccountCreationGasPerKey
		gas += uint64(len(newAccount.GetFeeKeys())) * params.TxAccountCreationGasPerKey
	}

	return newAccount, gas + readableGas, readable
}

func getRandomPrivateKeyString(t *testing.T) string {
	key, err := crypto.GenerateKey()
	assert.Equal(t, nil, err)
	keyBytes := crypto.FromECDSA(key)

	return common.Bytes2Hex(keyBytes)
}

// Return multisig parameters for creating RoleBased with MultiSig
func genMultiSigParamForRoleBased(t *testing.T) []TestCreateMultisigAccountParam {
	var params []TestCreateMultisigAccountParam
	param1 := TestCreateMultisigAccountParam{
		Threshold: uint(2),
		Weights:   []uint{1, 1, 1},
		PrvKeys:   []string{getRandomPrivateKeyString(t), getRandomPrivateKeyString(t), getRandomPrivateKeyString(t)},
	}
	params = append(params, param1)

	param2 := TestCreateMultisigAccountParam{
		Threshold: uint(2),
		Weights:   []uint{1, 1, 1},
		PrvKeys:   []string{getRandomPrivateKeyString(t), getRandomPrivateKeyString(t), getRandomPrivateKeyString(t)},
	}
	params = append(params, param2)

	param3 := TestCreateMultisigAccountParam{
		Threshold: uint(2),
		Weights:   []uint{1, 1, 1},
		PrvKeys:   []string{getRandomPrivateKeyString(t), getRandomPrivateKeyString(t), getRandomPrivateKeyString(t)},
	}
	params = append(params, param3)

	return params
}

func getIntrinsicGas(txType types.TxType) uint64 {
	var intrinsic uint64

	switch txType {
	case types.TxTypeLegacyTransaction:
		intrinsic = params.TxGas
	case types.TxTypeEthereumAccessList:
		intrinsic = params.TxGas
	case types.TxTypeEthereumDynamicFee:
		intrinsic = params.TxGas
	case types.TxTypeValueTransfer:
		intrinsic = params.TxGasValueTransfer
	case types.TxTypeFeeDelegatedValueTransfer:
		intrinsic = params.TxGasValueTransfer + params.TxGasFeeDelegated
	case types.TxTypeFeeDelegatedValueTransferWithRatio:
		intrinsic = params.TxGasValueTransfer + params.TxGasFeeDelegatedWithRatio
	case types.TxTypeValueTransferMemo:
		intrinsic = params.TxGasValueTransfer
	case types.TxTypeFeeDelegatedValueTransferMemo:
		intrinsic = params.TxGasValueTransfer + params.TxGasFeeDelegated
	case types.TxTypeFeeDelegatedValueTransferMemoWithRatio:
		intrinsic = params.TxGasValueTransfer + params.TxGasFeeDelegatedWithRatio
	case types.TxTypeAccountUpdate:
		intrinsic = params.TxGasAccountUpdate
	case types.TxTypeFeeDelegatedAccountUpdate:
		intrinsic = params.TxGasAccountUpdate + params.TxGasFeeDelegated
	case types.TxTypeFeeDelegatedAccountUpdateWithRatio:
		intrinsic = params.TxGasAccountUpdate + params.TxGasFeeDelegatedWithRatio
	case types.TxTypeSmartContractDeploy:
		intrinsic = params.TxGasContractCreation
	case types.TxTypeFeeDelegatedSmartContractDeploy:
		intrinsic = params.TxGasContractCreation + params.TxGasFeeDelegated
	case types.TxTypeFeeDelegatedSmartContractDeployWithRatio:
		intrinsic = params.TxGasContractCreation + params.TxGasFeeDelegatedWithRatio
	case types.TxTypeSmartContractExecution:
		intrinsic = params.TxGasContractExecution
	case types.TxTypeFeeDelegatedSmartContractExecution:
		intrinsic = params.TxGasContractExecution + params.TxGasFeeDelegated
	case types.TxTypeFeeDelegatedSmartContractExecutionWithRatio:
		intrinsic = params.TxGasContractExecution + params.TxGasFeeDelegatedWithRatio
	case types.TxTypeChainDataAnchoring:
		intrinsic = params.TxChainDataAnchoringGas
	case types.TxTypeFeeDelegatedChainDataAnchoring:
		intrinsic = params.TxChainDataAnchoringGas + params.TxGasFeeDelegated
	case types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		intrinsic = params.TxChainDataAnchoringGas + params.TxGasFeeDelegatedWithRatio
	case types.TxTypeCancel:
		intrinsic = params.TxGasCancel
	case types.TxTypeFeeDelegatedCancel:
		intrinsic = params.TxGasCancel + params.TxGasFeeDelegated
	case types.TxTypeFeeDelegatedCancelWithRatio:
		intrinsic = params.TxGasCancel + params.TxGasFeeDelegatedWithRatio
	}

	return intrinsic
}

// Implement TestAccount interface for TestAccountType
func (t *TestAccountType) GetAddr() common.Address {
	return t.Addr
}

func (t *TestAccountType) GetTxKeys() []*ecdsa.PrivateKey {
	return t.Keys
}

func (t *TestAccountType) GetUpdateKeys() []*ecdsa.PrivateKey {
	return t.Keys
}

func (t *TestAccountType) GetFeeKeys() []*ecdsa.PrivateKey {
	return t.Keys
}

func (t *TestAccountType) GetNonce() uint64 {
	return t.Nonce
}

func (t *TestAccountType) GetAccKey() accountkey.AccountKey {
	return t.AccKey
}

func (t *TestAccountType) SetNonce(nonce uint64) {
	t.Nonce = nonce
}

func (t *TestAccountType) SetAddr(addr common.Address) {
	t.Addr = addr
}

// Return SigValidationGas depends on AccountType
func (t *TestAccountType) GetValidationGas(r accountkey.RoleType) uint64 {
	if t.GetAccKey() == nil {
		return 0
	}

	var gas uint64

	switch t.GetAccKey().Type() {
	case accountkey.AccountKeyTypeLegacy:
		gas = 0
	case accountkey.AccountKeyTypePublic:
		gas = (1 - 1) * params.TxValidationGasPerKey
	case accountkey.AccountKeyTypeWeightedMultiSig:
		gas = uint64(len(t.GetTxKeys())-1) * params.TxValidationGasPerKey
	}

	return gas
}

func (t *TestAccountType) AddNonce() {
	t.Nonce += 1
}

// Implement TestAccount interface for TestRoleBasedAccountType
func (t *TestRoleBasedAccountType) GetAddr() common.Address {
	return t.Addr
}

func (t *TestRoleBasedAccountType) GetTxKeys() []*ecdsa.PrivateKey {
	return t.TxKeys
}

func (t *TestRoleBasedAccountType) GetUpdateKeys() []*ecdsa.PrivateKey {
	return t.UpdateKeys
}

func (t *TestRoleBasedAccountType) GetFeeKeys() []*ecdsa.PrivateKey {
	return t.FeeKeys
}

func (t *TestRoleBasedAccountType) GetNonce() uint64 {
	return t.Nonce
}

func (t *TestRoleBasedAccountType) GetAccKey() accountkey.AccountKey {
	return t.AccKey
}

func (t *TestRoleBasedAccountType) SetNonce(nonce uint64) {
	t.Nonce = nonce
}

func (t *TestRoleBasedAccountType) SetAddr(addr common.Address) {
	t.Addr = addr
}

// Return SigValidationGas depends on AccountType
func (t *TestRoleBasedAccountType) GetValidationGas(r accountkey.RoleType) uint64 {
	if t.GetAccKey() == nil {
		return 0
	}

	var gas uint64

	switch r {
	case accountkey.RoleTransaction:
		gas = uint64(len(t.GetTxKeys())-1) * params.TxValidationGasPerKey
	case accountkey.RoleAccountUpdate:
		gas = uint64(len(t.GetUpdateKeys())-1) * params.TxValidationGasPerKey
	case accountkey.RoleFeePayer:
		gas = uint64(len(t.GetFeeKeys())-1) * params.TxValidationGasPerKey
	}

	return gas
}

func (t *TestRoleBasedAccountType) AddNonce() {
	t.Nonce += 1
}
