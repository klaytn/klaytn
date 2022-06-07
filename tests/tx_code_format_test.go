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

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

type genTxWithCodeFormat func(t *testing.T, signer types.Signer, from TestAccount, payer TestAccount, gasPrice *big.Int, codeFormat params.CodeFormat) *types.Transaction

func TestCodeFormat(t *testing.T) {
	testFunctions := []struct {
		Name  string
		genTx genTxWithCodeFormat
	}{
		{"SmartContractDeployWithValidCodeFormat", genSmartContractDeployWithCodeFormat},
		{"FeeDelegatedSmartContractDeployWithValidCodeFormat", genFeeDelegatedSmartContractDeployWithCodeFormat},
		{"FeeDelegatedWithRatioSmartContractDeployWithValidCodeFormat", genFeeDelegatedWithRatioSmartContractDeployWithCodeFormat},
		{"SmartContractDeployWithInvalidCodeFormat", genSmartContractDeployWithCodeFormat},
		{"FeeDelegatedSmartContractDeployWithInvalidCodeFormat", genFeeDelegatedSmartContractDeployWithCodeFormat},
		{"FeeDelegatedWithRatioSmartContractDeployWithInvalidCodeFormat", genFeeDelegatedWithRatioSmartContractDeployWithCodeFormat},
	}

	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	assert.Equal(t, nil, err)
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

	testCodeFormat := func(tx *types.Transaction, state error) {
		receipt, _, err := applyTransaction(t, bcdata, tx)
		if err != nil {
			assert.Equal(t, state, err)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
		} else {
			assert.Equal(t, nil, err)
			assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
		}
	}

	for _, f := range testFunctions {
		codeFormat := params.CodeFormatEVM
		state := error(nil)
		var tx *types.Transaction

		if strings.Contains(f.Name, "WithInvalid") {
			codeFormat = params.CodeFormatLast
			state = kerrors.ErrInvalidCodeFormat
		}

		if strings.Contains(f.Name, "FeeDelegated") {
			tx = f.genTx(t, signer, reservoir, reservoir, gasPrice, codeFormat)
		} else {
			tx = f.genTx(t, signer, reservoir, nil, gasPrice, codeFormat)
		}

		t.Run(f.Name, func(t *testing.T) {
			testCodeFormat(tx, state)
		})
	}
}

func genSmartContractDeployWithCodeFormat(t *testing.T, signer types.Signer, from TestAccount, payer TestAccount, gasPrice *big.Int, codeFormat params.CodeFormat) *types.Transaction {
	values := genMapForDeployWithCodeFormat(t, from, gasPrice, codeFormat)

	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	return tx
}

func genFeeDelegatedSmartContractDeployWithCodeFormat(t *testing.T, signer types.Signer, from TestAccount, payer TestAccount, gasPrice *big.Int, codeFormat params.CodeFormat) *types.Transaction {
	values := genMapForDeployWithCodeFormat(t, from, gasPrice, codeFormat)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx
}

func genFeeDelegatedWithRatioSmartContractDeployWithCodeFormat(t *testing.T, signer types.Signer, from TestAccount, payer TestAccount, gasPrice *big.Int, codeFormat params.CodeFormat) *types.Transaction {
	values := genMapForDeployWithCodeFormat(t, from, gasPrice, codeFormat)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetTxKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx
}

func genMapForDeployWithCodeFormat(t *testing.T, from TestAccount, gasPrice *big.Int, codeFormat params.CodeFormat) map[types.TxValueKeyType]interface{} {
	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         from.GetNonce(),
		types.TxValueKeyAmount:        new(big.Int).SetUint64(0),
		types.TxValueKeyGasLimit:      gasLimit,
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyTo:            (*common.Address)(nil),
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyFrom:          from.GetAddr(),
		types.TxValueKeyData:          common.FromHex(code),
		types.TxValueKeyCodeFormat:    codeFormat,
	}

	return values
}
