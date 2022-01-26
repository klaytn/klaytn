// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/types/transaction_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// The values in those tests are from the Transaction Tests
// at github.com/ethereum/tests.
var (
	emptyTx = NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)

	rightvrsTx, _ = NewTransaction(
		3,
		common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b"),
		big.NewInt(10),
		2000,
		big.NewInt(1),
		common.FromHex("5544"),
	).WithSignature(
		NewEIP155Signer(common.Big1),
		common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
	)
)

func TestTransactionSigHash(t *testing.T) {
	signer := NewEIP155Signer(common.Big1)
	if signer.Hash(emptyTx) != common.HexToHash("a715f8447b97e3105d2cc0a8aca1466fa3a02f7cc6d2f9a3fe89f2581c9111c5") {
		t.Errorf("empty transaction hash mismatch, ɡot %x", signer.Hash(emptyTx))
	}
	if signer.Hash(rightvrsTx) != common.HexToHash("bd63ce94e66c7ffbce3b61023bbf9ee6df36047525b123201dcb5c4332f105ae") {
		t.Errorf("RightVRS transaction hash mismatch, ɡot %x", signer.Hash(rightvrsTx))
	}
}

func TestTransactionEncode(t *testing.T) {
	txb, err := rlp.EncodeToBytes(rightvrsTx)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	should := common.FromHex("f86103018207d094b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a82554426a098ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4aa08887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a3")
	if !bytes.Equal(txb, should) {
		t.Errorf("encoded RLP mismatch, ɡot %x", txb)
	}
}

func decodeTx(data []byte) (*Transaction, error) {
	var tx Transaction
	t, err := &tx, rlp.Decode(bytes.NewReader(data), &tx)

	return t, err
}

func defaultTestKey() (*ecdsa.PrivateKey, common.Address) {
	key, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

func TestRecipientEmpty(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(common.Hex2Bytes("f84980808080800126a0f18ba0124c1ed46fef6673ff5f614bafbb33a23ad92874bfa3cb3abad56d9a72a046690eb704a07384224d12e991da61faceefede59c6741d85e7d72e097855eaf"))
	if err != nil {
		t.Fatal(err)
	}

	signer := NewEIP155Signer(common.Big1)

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}
	if addr != from {
		t.Error("derived address doesn't match")
	}
}

func TestRecipientNormal(t *testing.T) {
	_, addr := defaultTestKey()

	tx, err := decodeTx(common.Hex2Bytes("f85d808080940000000000000000000000000000000000000000800126a0c1f2953a2277033c693f3d352b740479788672ba21e76d567557aa069b7e5061a06e798331dbd58c7438fe0e0a64b3b17c8378c726da3613abae8783b5dccc9944"))
	if err != nil {
		t.Fatal(err)
	}

	signer := NewEIP155Signer(common.Big1)

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}

	if addr != from {
		t.Fatal("derived address doesn't match")
	}
}

// Tests that transactions can be correctly sorted according to their price in
// decreasing order, but at the same time with increasing nonces when issued by
// the same account.
func TestTransactionPriceNonceSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 25)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	signer := NewEIP155Signer(common.Big1)
	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		for i := 0; i < 25; i++ {
			tx, _ := SignTx(NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(start+i)), nil), signer, key)
			groups[addr] = append(groups[addr], tx)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != 25*25 {
		t.Errorf("expected %d transactions, found %d", 25*25, len(txs))
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)

		// Make sure the nonce order is valid
		for j, txj := range txs[i+1:] {
			fromj, _ := Sender(signer, txj)

			if fromi == fromj && txi.Nonce() > txj.Nonce() {
				t.Errorf("invalid nonce ordering: tx #%d (A=%x N=%v) < tx #%d (A=%x N=%v)", i, fromi[:4], txi.Nonce(), i+j, fromj[:4], txj.Nonce())
			}
		}

		// If the next tx has different from account, the price must be lower than the current one
		if i+1 < len(txs) {
			next := txs[i+1]
			fromNext, _ := Sender(signer, next)
			if fromi != fromNext && txi.GasPrice().Cmp(next.GasPrice()) < 0 {
				t.Errorf("invalid ɡasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", i, fromi[:4], txi.GasPrice(), i+1, fromNext[:4], next.GasPrice())
			}
		}
	}
}

func TestGasOverflow(t *testing.T) {
	// AccountCreation
	// calculate gas for account creation
	numKeys := new(big.Int).SetUint64(accountkey.MaxNumKeysForMultiSig)
	gasPerKey := new(big.Int).SetUint64(params.TxAccountCreationGasPerKey)
	defaultGas := new(big.Int).SetUint64(params.TxAccountCreationGasDefault)
	txGas := new(big.Int).SetUint64(params.TxGasAccountCreation)
	totalGas := new(big.Int).Add(txGas, new(big.Int).Add(defaultGas, new(big.Int).Mul(numKeys, gasPerKey)))
	assert.Equal(t, true, totalGas.BitLen() <= 64)

	// ValueTransfer
	// calculate gas for validation of multisig accounts.
	gasPerKey = new(big.Int).SetUint64(params.TxValidationGasPerKey)
	defaultGas = new(big.Int).SetUint64(params.TxValidationGasDefault)
	txGas = new(big.Int).SetUint64(params.TxGas)
	totalGas = new(big.Int).Add(txGas, new(big.Int).Add(defaultGas, new(big.Int).Mul(numKeys, gasPerKey)))
	assert.Equal(t, true, totalGas.BitLen() <= 64)

	// TODO-Klaytn-Gas: Need to find a way of checking integer overflow for smart contract execution.
}

// TODO-Klaytn-FailedTest This test is failed in Klaytn
/*
// TestTransactionJSON tests serializing/de-serializing to/from JSON.
func TestTransactionJSON(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	signer := NewEIP155Signer(common.Big1)

	transactions := make([]*Transaction, 0, 50)
	for i := uint64(0); i < 25; i++ {
		var tx *Transaction
		switch i % 2 {
		case 0:
			tx = NewTransaction(i, common.Address{1}, common.Big0, 1, common.Big2, []byte("abcdef"))
		case 1:
			tx = NewContractCreation(i, common.Big0, 1, common.Big2, []byte("abcdef"))
		}
		transactions = append(transactions, tx)

		signedTx, err := SignTx(tx, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}

		transactions = append(transactions, signedTx)
	}

	for _, tx := range transactions {
		data, err := json.Marshal(tx)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		var parsedTx *Transaction
		if err := json.Unmarshal(data, &parsedTx); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}

		// compare nonce, price, gaslimit, recipient, amount, payload, V, R, S
		if tx.Hash() != parsedTx.Hash() {
			t.Errorf("parsed tx differs from original tx, want %v, ɡot %v", tx, parsedTx)
		}
		if tx.ChainId().Cmp(parsedTx.ChainId()) != 0 {
			t.Errorf("invalid chain id, want %d, ɡot %d", tx.ChainId(), parsedTx.ChainId())
		}
	}
}
*/

func TestIntrinsicGas(t *testing.T) {
	// testData contains two kind of members
	// inputString - test input data
	// expectGas - expect gas according to the specific condition.
	//            it differs depending on whether the contract is created or not,
	//            or whether it has passed through the Istanbul compatible block.
	testData := []struct {
		inputString string
		expectGas1  uint64 // contractCreate - false, isIstanbul - false
		expectGas2  uint64 // contractCreate - false, isIstanbul - true
		expectGas3  uint64 // contractCreate - true,  isIstanbul - false
		expectGas4  uint64 // contractCreate - true,  isIstanbul - true
	}{
		{"0000", 21008, 21200, 53008, 53200},
		{"1000", 21072, 21200, 53072, 53200},
		{"0100", 21072, 21200, 53072, 53200},
		{"ff3d", 21136, 21200, 53136, 53200},
		{"0000a6bc", 21144, 21400, 53144, 53400},
		{"fd00fd00", 21144, 21400, 53144, 53400},
		{"", 21000, 21000, 53000, 53000},
	}
	for _, tc := range testData {
		var (
			data []byte // input data entered through the tx argument
			gas  uint64 // the gas varies depending on what comes in as a condition(contractCreate & IsIstanbul)
			err  error  // in this unittest, every testcase returns nil error.
		)

		data, err = hex.DecodeString(tc.inputString) // decode input string to hex data
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, false, params.Rules{IsIstanbul: false})
		assert.Equal(t, tc.expectGas1, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, false, params.Rules{IsIstanbul: true})
		assert.Equal(t, tc.expectGas2, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, true, params.Rules{IsIstanbul: false})
		assert.Equal(t, tc.expectGas3, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, true, params.Rules{IsIstanbul: true})
		assert.Equal(t, tc.expectGas4, gas)
		assert.Equal(t, nil, err)
	}
}
