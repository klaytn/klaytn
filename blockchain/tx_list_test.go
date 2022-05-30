// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from core/tx_list_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/crypto"
	"github.com/stretchr/testify/assert"
)

// Tests that transactions can be added to strict lists and list contents and
// nonce boundaries are correctly maintained.
func TestStrictTxListAdd(t *testing.T) {
	// Generate a list of transactions to insert
	key, _ := crypto.GenerateKey()

	txs := make(types.Transactions, 1024)
	for i := 0; i < len(txs); i++ {
		txs[i] = transaction(uint64(i), 0, key)
	}
	// Insert the transactions in a random order
	list := newTxList(true)
	for _, v := range rand.Perm(len(txs)) {
		list.Add(txs[v], DefaultTxPoolConfig.PriceBump)
	}
	// Verify internal state
	if len(list.txs.items) != len(txs) {
		t.Errorf("transaction count mismatch: have %d, want %d", len(list.txs.items), len(txs))
	}
	for i, tx := range txs {
		if list.txs.items[tx.Nonce()] != tx {
			t.Errorf("item %d: transaction mismatch: have %v, want %v", i, list.txs.items[tx.Nonce()], tx)
		}
	}
}

// TestTxListReadyWithGasPrice check whether ReadyWithGasPrice() works well.
// It makes a slice of 10 transactions and executes ReadyWithGasPrice() with baseFee 30.
func TestTxListReadyWithGasPrice(t *testing.T) {
	// Start nonce : 3
	// baseFee : 30
	// Transaction[0:9] has gasPrice 50
	// The result of executing ReadyWithGasPrice will be Transaction[0:9].
	startNonce := 3
	expectedBaseFee := big.NewInt(30)
	nTxs := 10
	// Generate a list of transactions to insert
	key, _ := crypto.GenerateKey()

	txs := make(types.Transactions, nTxs)
	nonce := startNonce
	for i := 0; i < len(txs); i++ {
		txs[i] = pricedTransaction(uint64(nonce), 0, big.NewInt(50), key)
		nonce++
	}

	// Insert the transactions in a random order
	list := newTxList(true)
	rand.Seed(time.Now().UnixNano())
	for _, v := range rand.Perm(len(txs)) {
		list.Add(txs[v], DefaultTxPoolConfig.PriceBump)
	}

	ready := list.ReadyWithGasPrice(uint64(startNonce), expectedBaseFee)

	if ready.Len() != nTxs {
		t.Error("expected number of filtered txs", nTxs, "got", ready.Len())
	}

	nonce = startNonce
	for i := 0; i < len(ready); i++ {
		if result := reflect.DeepEqual(ready[i], txs[i]); !result {
			t.Error("ready hash : ", ready[i].Hash(), "tx[i] hash : ", txs[i].Hash())
		}

		if list.txs.items[uint64(nonce)] != nil {
			t.Error("Not deleted in txList. nonce : ", nonce, "tx hash : ", list.txs.items[uint64(nonce)].Hash())
		}
		nonce++
	}
}

// TestTxListReadyWithGasPricePartialFilter check whether ReadyWithGasPrice() works well.
// It makes a slice has 10 transactions and executes ReadyWithGasPrice() with baseFee 20.
//
// It checks whether filtering works well if there is a transaction
// with a gasPrice lower than the baseFee among Tx.
func TestTxListReadyWithGasPricePartialFilter(t *testing.T) {
	// Start nonce : 3
	// baseFee : 20
	// Transaction[0:6] has gasPrice 30
	// Transaction[7] has gasPrice 10
	// Transaction[8:9] has gasPrice 50
	// The result of executing ReadyWithGasPrice will be Transaction[0:6].
	startNonce := 3
	expectedBaseFee := big.NewInt(20)
	nTxs := 10

	// Generate a list of transactions to insert
	key, _ := crypto.GenerateKey()

	txs := make(types.Transactions, nTxs)
	nonce := startNonce
	for i := 0; i < len(txs); i++ {
		// Set the gasPrice of transaction lower than expectedBaseFee.
		if i == 7 {
			txs[i] = pricedTransaction(uint64(nonce), 0, big.NewInt(10), key)
		} else if i > 7 {
			txs[i] = pricedTransaction(uint64(nonce), 0, big.NewInt(50), key)
		} else {
			txs[i] = pricedTransaction(uint64(nonce), 0, big.NewInt(30), key)
		}

		nonce++
	}

	// Insert the transactions in a random order
	list := newTxList(true)
	for _, v := range rand.Perm(len(txs)) {
		list.Add(txs[v], DefaultTxPoolConfig.PriceBump)
	}

	ready := list.ReadyWithGasPrice(uint64(startNonce), expectedBaseFee)

	if ready.Len() != 7 {
		t.Error("expected filtered txs length", 7, "got", ready.Len())
	}

	nonce = startNonce
	for i := 0; i < len(ready); i++ {
		if result := reflect.DeepEqual(ready[i], txs[i]); result == false {
			t.Error("ready hash : ", ready[i].Hash(), "tx[i] hash : ", txs[i].Hash())
		}

		if list.txs.items[uint64(nonce)] != nil {
			t.Error("Not deleted in txList. nonce : ", nonce, "tx hash : ", list.txs.items[uint64(nonce)].Hash())
		}
		nonce++
	}
}

// TestSubstituteTxByGasPrice tests if the new tx has been successfully replaced
// if it has the same nonce and greater gas price as the old tx.
func TestSubstituteTxByGasPrice(t *testing.T) {
	// Generate a list of transactions to insert
	key, _ := crypto.GenerateKey()
	txList := newTxList(false)

	oldTx := pricedTransaction(0, 21000, big.NewInt(50), key)
	newTx := pricedTransaction(0, 21000, big.NewInt(60), key)

	if result, _ := txList.Add(oldTx, DefaultTxPoolConfig.PriceBump); !result {
		t.Error("it cannot add tx in tx list.")
	}

	result, replaced := txList.Add(newTx, DefaultTxPoolConfig.PriceBump)
	if !result {
		t.Error("it cannot replace tx in tx list.")
	}

	assert.Equal(t, replaced, oldTx)
}

// TestSubstituteTransactionAbort checks if a new tx aborts a new transaction
// if it has the same nonce and lower gas price as the old tx.
func TestSubstituteTransactionAbort(t *testing.T) {
	// Generate a list of transactions to insert
	key, _ := crypto.GenerateKey()
	txList := newTxList(false)

	oldTx := pricedTransaction(0, 21000, big.NewInt(50), key)
	newTx := pricedTransaction(0, 21000, big.NewInt(40), key)

	if result, _ := txList.Add(oldTx, DefaultTxPoolConfig.PriceBump); !result {
		t.Error("it cannot add tx in tx list.")
	}

	if result, replaced := txList.Add(newTx, DefaultTxPoolConfig.PriceBump); result || replaced != nil {
		t.Error("Expected to not substitute by a tx with lower gas price")
	}

}
