// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from tests/transaction_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package tests

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	txt := new(testMatcher)
	txt.config(`^Byzantium/`, params.ChainConfig{})

	txt.walk(t, transactionTestDir, func(t *testing.T, name string, test *TransactionTest) {
		cfg := txt.findConfig(name)
		if err := txt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}

// TestAccountCreationDisable tries to use accountCreation tx types which is disabled now.
// The tx should be invalided in txPool and execution process.
func TestAccountCreationDisable(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	// the same with types.errUndefinedTxType
	errUndefinedTxType := errors.New("undefined tx type")

	// Initialize blockchain
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	anon, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e989")
	assert.Equal(t, nil, err)

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	{
		// generate an accountCreation tx
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.GetNonce(),
			types.TxValueKeyFrom:          reservoir.GetAddr(),
			types.TxValueKeyTo:            anon.Addr,
			types.TxValueKeyAmount:        big.NewInt(0),
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      big.NewInt(25 * params.Ston),
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyAccountKey:    anon.AccKey,
		}

		tx, err := types.NewAccountCreationTransactionWithMap(values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		// fail to add tx in txPool
		err = txpool.AddRemote(tx)
		assert.Equal(t, errUndefinedTxType, err)

		// fail to execute tx
		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, errUndefinedTxType, err)
		assert.Equal(t, (*types.Receipt)(nil), receipt)
	}
}

// TestContractDeployWithDisabledAddress tests invalid contract deploy transactions.
// 1. If the humanReadable field of an tx is 'true', it should fail.
// 2. If the recipient field of an tx is not nil, it should fail.
func TestContractDeployWithDisabledAddress(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	testTxTypes := []types.TxType{
		types.TxTypeSmartContractDeploy,
		types.TxTypeFeeDelegatedSmartContractDeploy,
		types.TxTypeFeeDelegatedSmartContractDeployWithRatio,
	}

	// Initialize blockchain
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	contract, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e989")
	assert.Equal(t, nil, err)

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	for _, txType := range testTxTypes {
		// generate an invalid contract deploy tx with humanReadable flag as true
		{
			values, _ := genMapForTxTypes(reservoir, nil, txType)
			values[types.TxValueKeyHumanReadable] = true

			tx, err := types.NewTransactionWithMap(txType, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.Keys)
			assert.Equal(t, nil, err)

			if txType.IsFeeDelegatedTransaction() {
				err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
				assert.Equal(t, nil, err)
			}
			// fail to add tx in txPool
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrHumanReadableNotSupported, err)

			// fail to execute tx
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, kerrors.ErrHumanReadableNotSupported, err)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
		}

		// generate an invalid contract deploy tx with an recipient address not nil
		{
			values, _ := genMapForTxTypes(reservoir, nil, txType)
			values[types.TxValueKeyTo] = &contract.Addr

			tx, err := types.NewTransactionWithMap(txType, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.Keys)
			assert.Equal(t, nil, err)

			if txType.IsFeeDelegatedTransaction() {
				err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
				assert.Equal(t, nil, err)
			}
			// fail to add tx in txPool
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrInvalidContractAddress, err)

			// fail to execute tx
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, kerrors.ErrInvalidContractAddress, err)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
		}
	}
}
