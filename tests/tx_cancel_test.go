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
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/log"
	"github.com/stretchr/testify/assert"
)

// TestTxCancel tests TxCancel transaction types:
// 1. Insert a value transfer transaction with nonce 0.
// 2. Insert a value transfer transaction with nonce 0. This should not be replaced.
// 3. Insert a TxCancel transaction with nonce 0. This should replace the tx with the same nonce.
// 4. Insert a TxCancel transaction with nonce 0 and different gas limit. This should replace the tx with the same nonce.
func TestTxCancel(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()
	opt := testOption{1000, 2000, 4, 1, []byte{}, makeNewTransactionsToRandom}

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(opt.numMaxAccounts, opt.numValidators)
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

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 1. Insert a value transfer transaction with nonce 0.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyTo:       *bcdata.addrs[0],
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 2. Insert a value transfer transaction with nonce 0. This should not be replaced.
	{
		var txs types.Transactions

		pending, queued := txpool.Content()
		oldtx := pending[*bcdata.addrs[0]][0]

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyTo:       *bcdata.addrs[1],
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued = txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, oldtx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 3. Insert a TxCancel transaction with nonce 0. This should replace the tx with the same nonce.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeCancel, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 4. Insert a TxCancel transaction with nonce 0 and different gas limit. This should replace the tx with the same nonce.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyGasLimit: gasLimit + 10,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeCancel, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestTxFeeDelegatedCancel tests TxCancel transaction types:
// 1. Insert a value transfer transaction with nonce 0.
// 2. Insert a value transfer transaction with nonce 0. This should not be replaced.
// 3. Insert a TxCancel transaction with nonce 0. This should replace the tx with the same nonce.
// 4. Insert a TxCancel transaction with nonce 0 and different gas limit. This should replace the tx with the same nonce.
func TestTxFeeDelegatedCancel(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()
	opt := testOption{1000, 2000, 4, 1, []byte{}, makeNewTransactionsToRandom}

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(opt.numMaxAccounts, opt.numValidators)
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

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 1. Insert a value transfer transaction with nonce 0.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyTo:       *bcdata.addrs[0],
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 2. Insert a value transfer transaction with nonce 0. This should not be replaced.
	{
		var txs types.Transactions

		pending, queued := txpool.Content()
		oldtx := pending[*bcdata.addrs[0]][0]

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyTo:       *bcdata.addrs[1],
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued = txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, oldtx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 3. Insert a TxCancel transaction with nonce 0. This should replace the tx with the same nonce.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: *bcdata.addrs[1],
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[1]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 4. Insert a TxCancel transaction with nonce 0 and different gas limit. This should replace the tx with the same nonce.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyGasLimit: gasLimit + 10,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: *bcdata.addrs[1],
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[1]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestTxFeeDelegatedCancelWithRatio tests TxCancel transaction types:
// 1. Insert a value transfer transaction with nonce 0.
// 2. Insert a value transfer transaction with nonce 0. This should not be replaced.
// 3. Insert a TxCancel transaction with nonce 0. This should replace the tx with the same nonce.
// 4. Insert a TxCancel transaction with nonce 0 and different gas limit. This should replace the tx with the same nonce.
func TestTxFeeDelegatedCancelWithRatio(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()
	opt := testOption{1000, 2000, 4, 1, []byte{}, makeNewTransactionsToRandom}

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(opt.numMaxAccounts, opt.numValidators)
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

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 1. Insert a value transfer transaction with nonce 0.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyTo:       *bcdata.addrs[0],
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 2. Insert a value transfer transaction with nonce 0. This should not be replaced.
	{
		var txs types.Transactions

		pending, queued := txpool.Content()
		oldtx := pending[*bcdata.addrs[0]][0]

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(0),
			types.TxValueKeyFrom:     *bcdata.addrs[0],
			types.TxValueKeyTo:       *bcdata.addrs[1],
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued = txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, oldtx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 3. Insert a TxCancel transaction with nonce 0. This should replace the tx with the same nonce.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(0),
			types.TxValueKeyFrom:               *bcdata.addrs[0],
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           *bcdata.addrs[1],
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[1]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	// 4. Insert a TxCancel transaction with nonce 0 and different gas limit. This should replace the tx with the same nonce.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(0),
			types.TxValueKeyFrom:               *bcdata.addrs[0],
			types.TxValueKeyGasLimit:           gasLimit + 10,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           *bcdata.addrs[1],
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[1]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		txpool.AddRemotes(txs)

		pending, queued := txpool.Content()
		assert.Equal(t, 0, len(queued))
		assert.Equal(t, 1, len(pending))
		assert.True(t, tx.Equal(pending[*bcdata.addrs[0]][0]))
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}
