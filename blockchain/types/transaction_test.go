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
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

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
	testAddr = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")

	emptyTx = NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)

	rightvrsTx, _ = NewTransaction(
		3,
		testAddr,
		big.NewInt(10),
		2000,
		big.NewInt(1),
		common.FromHex("5544"),
	).WithSignature(
		LatestSignerForChainID(common.Big1),
		common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
	)

	accessListTx = TxInternalDataEthereumAccessList{
		ChainID:      big.NewInt(1),
		AccountNonce: 3,
		Recipient:    &testAddr,
		Amount:       big.NewInt(10),
		GasLimit:     25000,
		Price:        big.NewInt(1),
		Payload:      common.FromHex("5544"),
	}

	accessAddr   = common.HexToAddress("0x0000000000000000000000000000000000000001")
	dynamicFeeTx = TxInternalDataEthereumDynamicFee{
		ChainID:      big.NewInt(1),
		AccountNonce: 3,
		Recipient:    &testAddr,
		Amount:       big.NewInt(10),
		GasLimit:     25000,
		GasFeeCap:    big.NewInt(1),
		GasTipCap:    big.NewInt(1),
		Payload:      common.FromHex("5544"),
		AccessList:   AccessList{{Address: accessAddr, StorageKeys: []common.Hash{{0}}}},
	}

	emptyEip2718Tx = &Transaction{
		data: &accessListTx,
	}

	emptyEip1559Tx = &Transaction{
		data: &dynamicFeeTx,
	}

	signedEip2718Tx, _ = emptyEip2718Tx.WithSignature(
		NewEIP2930Signer(big.NewInt(1)),
		common.Hex2Bytes("c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b266032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d3752101"),
	)

	signedEip1559Tx, _ = emptyEip1559Tx.WithSignature(
		NewLondonSigner(big.NewInt(1)),
		common.Hex2Bytes("c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b266032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d3752101"))
)

func TestTransactionSigHash(t *testing.T) {
	signer := LatestSignerForChainID(common.Big1)
	if signer.Hash(emptyTx) != common.HexToHash("a715f8447b97e3105d2cc0a8aca1466fa3a02f7cc6d2f9a3fe89f2581c9111c5") {
		t.Errorf("empty transaction hash mismatch, ɡot %x", signer.Hash(emptyTx))
	}
	if signer.Hash(rightvrsTx) != common.HexToHash("bd63ce94e66c7ffbce3b61023bbf9ee6df36047525b123201dcb5c4332f105ae") {
		t.Errorf("RightVRS transaction hash mismatch, ɡot %x", signer.Hash(rightvrsTx))
	}
}

func TestEIP2718TransactionSigHash(t *testing.T) {
	s := NewEIP2930Signer(big.NewInt(1))
	if s.Hash(emptyEip2718Tx) != common.HexToHash("49b486f0ec0a60dfbbca2d30cb07c9e8ffb2a2ff41f29a1ab6737475f6ff69f3") {
		t.Errorf("empty EIP-2718 transaction hash mismatch, got %x", s.Hash(emptyEip2718Tx))
	}
	if s.Hash(signedEip2718Tx) != common.HexToHash("49b486f0ec0a60dfbbca2d30cb07c9e8ffb2a2ff41f29a1ab6737475f6ff69f3") {
		t.Errorf("signed EIP-2718 transaction hash mismatch, got %x", s.Hash(signedEip2718Tx))
	}
}

func TestEIP1559TransactionSigHash(t *testing.T) {
	s := NewLondonSigner(big.NewInt(1))
	if s.Hash(emptyEip1559Tx) != common.HexToHash("a52ce25a7d108740bce8fbb2dfa1f26793b2e8eea94a7700bedbae13cbdd8a0f") {
		t.Errorf("empty EIP-1559 transaction hash mismatch, got %x", s.Hash(emptyEip2718Tx))
	}
	if s.Hash(signedEip1559Tx) != common.HexToHash("a52ce25a7d108740bce8fbb2dfa1f26793b2e8eea94a7700bedbae13cbdd8a0f") {
		t.Errorf("signed EIP-1559 transaction hash mismatch, got %x", s.Hash(signedEip2718Tx))
	}
}

// This test checks signature operations on access list transactions.
func TestEIP2930Signer(t *testing.T) {
	var (
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		keyAddr = crypto.PubkeyToAddress(key.PublicKey)
		signer1 = NewEIP2930Signer(big.NewInt(1))
		signer2 = NewEIP2930Signer(big.NewInt(2))
		tx0     = NewTx(&TxInternalDataEthereumAccessList{AccountNonce: 1, ChainID: new(big.Int)})
		tx1     = NewTx(&TxInternalDataEthereumAccessList{ChainID: big.NewInt(1), AccountNonce: 1, V: new(big.Int), R: new(big.Int), S: new(big.Int)})
		tx2, _  = SignTx(NewTx(&TxInternalDataEthereumAccessList{ChainID: big.NewInt(2), AccountNonce: 1}), signer2, key)
	)

	tests := []struct {
		tx             *Transaction
		signer         Signer
		wantSignerHash common.Hash
		wantSenderErr  error
		wantSignErr    error
		wantHash       common.Hash // after signing
	}{
		{
			tx:             tx0,
			signer:         signer1,
			wantSignerHash: common.HexToHash("846ad7672f2a3a40c1f959cd4a8ad21786d620077084d84c8d7c077714caa139"),
			wantSenderErr:  ErrInvalidChainId,
			wantHash:       common.HexToHash("1ccd12d8bbdb96ea391af49a35ab641e219b2dd638dea375f2bc94dd290f2549"),
		},
		{
			tx:             tx1,
			signer:         signer1,
			wantSenderErr:  ErrInvalidSig,
			wantSignerHash: common.HexToHash("846ad7672f2a3a40c1f959cd4a8ad21786d620077084d84c8d7c077714caa139"),
			wantHash:       common.HexToHash("1ccd12d8bbdb96ea391af49a35ab641e219b2dd638dea375f2bc94dd290f2549"),
		},
		{
			// This checks what happens when trying to sign an unsigned tx for the wrong chain.
			tx:             tx1,
			signer:         signer2,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("846ad7672f2a3a40c1f959cd4a8ad21786d620077084d84c8d7c077714caa139"),
			wantSignErr:    ErrInvalidChainId,
		},
		{
			// This checks what happens when trying to re-sign a signed tx for the wrong chain.
			tx:             tx2,
			signer:         signer1,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("367967247499343401261d718ed5aa4c9486583e4d89251afce47f4a33c33362"),
			wantSignErr:    ErrInvalidChainId,
		},
	}

	for i, test := range tests {
		sigHash := test.signer.Hash(test.tx)
		if sigHash != test.wantSignerHash {
			t.Errorf("test %d: wrong sig hash: got %x, want %x", i, sigHash, test.wantSignerHash)
		}
		sender, err := Sender(test.signer, test.tx)
		if err != test.wantSenderErr {
			t.Errorf("test %d: wrong Sender error %q", i, err)
		}
		if err == nil && sender != keyAddr {
			t.Errorf("test %d: wrong sender address %x", i, sender)
		}
		signedTx, err := SignTx(test.tx, test.signer, key)
		if err != test.wantSignErr {
			t.Fatalf("test %d: wrong SignTx error %q", i, err)
		}
		if signedTx != nil {
			if signedTx.Hash() != test.wantHash {
				t.Errorf("test %d: wrong tx hash after signing: got %x, want %x", i, signedTx.Hash(), test.wantHash)
			}
		}
	}
}

// This test checks signature operations on dynamic fee transactions.
func TestLondonSigner(t *testing.T) {
	var (
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		keyAddr = crypto.PubkeyToAddress(key.PublicKey)
		signer1 = NewLondonSigner(big.NewInt(1))
		signer2 = NewLondonSigner(big.NewInt(2))
		tx0     = NewTx(&TxInternalDataEthereumDynamicFee{AccountNonce: 1, ChainID: new(big.Int)})
		tx1     = NewTx(&TxInternalDataEthereumDynamicFee{ChainID: big.NewInt(1), AccountNonce: 1, V: new(big.Int), R: new(big.Int), S: new(big.Int)})
		tx2, _  = SignTx(NewTx(&TxInternalDataEthereumDynamicFee{ChainID: big.NewInt(2), AccountNonce: 1}), signer2, key)
	)

	tests := []struct {
		tx             *Transaction
		signer         Signer
		wantSignerHash common.Hash
		wantSenderErr  error
		wantSignErr    error
		wantHash       common.Hash // after signing
	}{
		{
			tx:             tx0,
			signer:         signer1,
			wantSignerHash: common.HexToHash("b6afee4d44e0392fb5d3204b350596d6677440bced7ebd998db73c9671527c57"),
			wantSenderErr:  ErrInvalidChainId,
			wantHash:       common.HexToHash("a2c6373b7eed946fd4165a0d8503aa26afc8e99f09e2be58b332fbbedc279f7a"),
		},
		{
			tx:             tx1,
			signer:         signer1,
			wantSenderErr:  ErrInvalidSig,
			wantSignerHash: common.HexToHash("b6afee4d44e0392fb5d3204b350596d6677440bced7ebd998db73c9671527c57"),
			wantHash:       common.HexToHash("a2c6373b7eed946fd4165a0d8503aa26afc8e99f09e2be58b332fbbedc279f7a"),
		},
		{
			// This checks what happens when trying to sign an unsigned tx for the wrong chain.
			tx:             tx1,
			signer:         signer2,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("b6afee4d44e0392fb5d3204b350596d6677440bced7ebd998db73c9671527c57"),
			wantSignErr:    ErrInvalidChainId,
		},
		{
			// This checks what happens when trying to re-sign a signed tx for the wrong chain.
			tx:             tx2,
			signer:         signer1,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("b0759fc55582f3e60ded82843dcc17733d8c65f543d2cf2613a47a5c6ac9fc48"),
			wantSignErr:    ErrInvalidChainId,
		},
	}

	for i, test := range tests {
		sigHash := test.signer.Hash(test.tx)
		if sigHash != test.wantSignerHash {
			t.Errorf("test %d: wrong sig hash: got %x, want %x", i, sigHash, test.wantSignerHash)
		}
		sender, err := Sender(test.signer, test.tx)
		if err != test.wantSenderErr {
			t.Errorf("test %d: wrong Sender error %q", i, err)
		}
		if err == nil && sender != keyAddr {
			t.Errorf("test %d: wrong sender address %x", i, sender)
		}
		signedTx, err := SignTx(test.tx, test.signer, key)
		if err != test.wantSignErr {
			t.Fatalf("test %d: wrong SignTx error %q", i, err)
		}
		if signedTx != nil {
			if signedTx.Hash() != test.wantHash {
				t.Errorf("test %d: wrong tx hash after signing: got %x, want %x", i, signedTx.Hash(), test.wantHash)
			}
		}
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

func TestEIP2718TransactionEncode(t *testing.T) {
	// RLP representation
	{
		have, err := rlp.EncodeToBytes(signedEip2718Tx)
		if err != nil {
			t.Fatalf("encode error: %v", err)
		}
		want := common.FromHex("7801f8630103018261a894b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a825544c001a0c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b2660a032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d37521")
		if !bytes.Equal(have, want) {
			t.Errorf("encoded RLP mismatch, got %x", have)
		}
	}
}

func TestEIP1559TransactionEncode(t *testing.T) {
	// RLP representation
	{
		have, err := rlp.EncodeToBytes(signedEip1559Tx)
		if err != nil {
			t.Fatalf("encode error: %v", err)
		}
		want := common.FromHex("7802f89d010301018261a894b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a825544f838f7940000000000000000000000000000000000000001e1a0000000000000000000000000000000000000000000000000000000000000000001a0c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b2660a032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d37521")
		if !bytes.Equal(have, want) {
			t.Errorf("encoded RLP mismatch, got %x", have)
		}
	}
}

func TestEffectiveGasPrice(t *testing.T) {
	legacyTx := NewTx(&TxInternalDataLegacy{Price: big.NewInt(1000)})
	dynamicTx := NewTx(&TxInternalDataEthereumDynamicFee{GasFeeCap: big.NewInt(4000), GasTipCap: big.NewInt(1000)})

	baseFee := big.NewInt(2000)
	have := legacyTx.EffectiveGasPrice(baseFee)
	want := big.NewInt(1000)
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasPrice(baseFee)
	want = big.NewInt(3000)
	assert.Equal(t, want, have)

	baseFee = big.NewInt(0)
	have = legacyTx.EffectiveGasPrice(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasPrice(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)
}

func TestEffectiveGasTip(t *testing.T) {
	legacyTx := NewTx(&TxInternalDataLegacy{Price: big.NewInt(1000)})
	dynamicTx := NewTx(&TxInternalDataEthereumDynamicFee{GasFeeCap: big.NewInt(4000), GasTipCap: big.NewInt(1000)})

	baseFee := big.NewInt(2000)
	have := legacyTx.EffectiveGasTip(baseFee)
	want := big.NewInt(1000)
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasTip(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)

	baseFee = big.NewInt(0)
	have = legacyTx.EffectiveGasTip(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasTip(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)

	a := new(big.Int)
	assert.Equal(t, 0, a.BitLen())
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

	signer := LatestSignerForChainID(common.Big1)

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

	signer := LatestSignerForChainID(common.Big1)

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

	signer := LatestSignerForChainID(common.Big1)
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
			gas  uint64 // the gas varies depending on what comes in as a condition(contractCreate & IsIstanbulForkEnabled)
			err  error  // in this unittest, every testcase returns nil error.
		)

		data, err = hex.DecodeString(tc.inputString) // decode input string to hex data
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, false, params.Rules{IsIstanbul: false})
		assert.Equal(t, tc.expectGas1, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, false, params.Rules{IsIstanbul: true})
		assert.Equal(t, tc.expectGas2, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, true, params.Rules{IsIstanbul: false})
		assert.Equal(t, tc.expectGas3, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, true, params.Rules{IsIstanbul: true})
		assert.Equal(t, tc.expectGas4, gas)
		assert.Equal(t, nil, err)
	}
}

// Tests that if multiple transactions have the same price, the ones seen earlier
// are prioritized to avoid network spam attacks aiming for a specific ordering.
func TestTransactionTimeSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)

		tx, _ := SignTx(NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
		tx.time = time.Unix(0, int64(len(keys)-start))

		groups[addr] = append(groups[addr], tx)
	}
	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != len(keys) {
		t.Errorf("expected %d transactions, found %d", len(keys), len(txs))
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)
		if i+1 < len(txs) {
			next := txs[i+1]
			fromNext, _ := Sender(signer, next)

			if txi.GasPrice().Cmp(next.GasPrice()) < 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", i, fromi[:4], txi.GasPrice(), i+1, fromNext[:4], next.GasPrice())
			}
			// Make sure time order is ascending if the txs have the same gas price
			if txi.GasPrice().Cmp(next.GasPrice()) == 0 && txi.time.After(next.time) {
				t.Errorf("invalid received time ordering: tx #%d (A=%x T=%v) > tx #%d (A=%x T=%v)", i, fromi[:4], txi.time, i+1, fromNext[:4], next.time)
			}
		}
	}
}

// TestTransactionCoding tests serializing/de-serializing to/from rlp and JSON.
func TestTransactionCoding(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	var (
		signer    = LatestSignerForChainID(common.Big1)
		addr      = common.HexToAddress("0x0000000000000000000000000000000000000001")
		recipient = common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
		accesses  = AccessList{{Address: addr, StorageKeys: []common.Hash{{0}}}}
	)
	for i := uint64(0); i < 500; i++ {
		var txData TxInternalData
		switch i % 5 {
		case 0:
			// Legacy tx.
			txData = &TxInternalDataLegacy{
				AccountNonce: i,
				Recipient:    &recipient,
				GasLimit:     1,
				Price:        big.NewInt(2),
				Payload:      []byte("abcdef"),
			}
		case 1:
			// Legacy tx contract creation.
			txData = &TxInternalDataLegacy{
				AccountNonce: i,
				GasLimit:     1,
				Price:        big.NewInt(2),
				Payload:      []byte("abcdef"),
			}
		case 2:
			// Tx with non-zero access list.
			txData = &TxInternalDataEthereumAccessList{
				ChainID:      big.NewInt(1),
				AccountNonce: i,
				Recipient:    &recipient,
				GasLimit:     123457,
				Price:        big.NewInt(10),
				AccessList:   accesses,
				Payload:      []byte("abcdef"),
			}
		case 3:
			// Tx with empty access list.
			txData = &TxInternalDataEthereumAccessList{
				ChainID:      big.NewInt(1),
				AccountNonce: i,
				Recipient:    &recipient,
				GasLimit:     123457,
				Price:        big.NewInt(10),
				Payload:      []byte("abcdef"),
			}
		case 4:
			// Contract creation with access list.
			txData = &TxInternalDataEthereumAccessList{
				ChainID:      big.NewInt(1),
				AccountNonce: i,
				GasLimit:     123457,
				Price:        big.NewInt(10),
				AccessList:   accesses,
			}
		case 5:
			// Tx with non-zero access list.
			txData = &TxInternalDataEthereumDynamicFee{
				ChainID:      big.NewInt(1),
				AccountNonce: i,
				Recipient:    &recipient,
				GasLimit:     123457,
				GasFeeCap:    big.NewInt(10),
				GasTipCap:    big.NewInt(10),
				AccessList:   accesses,
				Payload:      []byte("abcdef"),
			}
		case 6:
			// Tx with dynamic fee.
			txData = &TxInternalDataEthereumDynamicFee{
				ChainID:      big.NewInt(1),
				AccountNonce: i,
				Recipient:    &recipient,
				GasLimit:     123457,
				GasFeeCap:    big.NewInt(10),
				GasTipCap:    big.NewInt(10),
				Payload:      []byte("abcdef"),
			}
		case 7:
			// Contract creation with dynamic fee tx.
			txData = &TxInternalDataEthereumDynamicFee{
				ChainID:      big.NewInt(1),
				AccountNonce: i,
				GasLimit:     123457,
				GasFeeCap:    big.NewInt(10),
				GasTipCap:    big.NewInt(10),
				AccessList:   accesses,
			}
		}

		transaction := Transaction{data: txData}
		tx, err := SignTx(&transaction, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}
		// RLP
		parsedTx, err := encodeDecodeBinary(tx)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(parsedTx, tx)

		// JSON
		parsedTx, err = encodeDecodeJSON(tx)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(parsedTx, tx)
	}
}

func encodeDecodeJSON(tx *Transaction) (*Transaction, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("json encoding failed: %v", err)
	}
	var parsedTx = &Transaction{}
	if err := json.Unmarshal(data, &parsedTx); err != nil {
		return nil, fmt.Errorf("json decoding failed: %v", err)
	}
	return parsedTx, nil
}

func encodeDecodeBinary(tx *Transaction) (*Transaction, error) {
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("rlp encoding failed: %v", err)
	}
	var parsedTx = &Transaction{}
	if err := parsedTx.UnmarshalBinary(data); err != nil {
		return nil, fmt.Errorf("rlp decoding failed: %v", err)
	}
	return parsedTx, nil
}

func assertEqual(orig *Transaction, cpy *Transaction) error {
	// compare nonce, price, gaslimit, recipient, amount, payload, V, R, S
	if want, got := orig.Hash(), cpy.Hash(); want != got {
		return fmt.Errorf("parsed tx differs from original tx, want %v, got %v", want, got)
	}
	if want, got := orig.ChainId(), cpy.ChainId(); want.Cmp(got) != 0 {
		return fmt.Errorf("invalid chain id, want %d, got %d", want, got)
	}

	if orig.Type().IsEthTypedTransaction() && cpy.Type().IsEthTypedTransaction() {
		tOrig := orig.data.(TxInternalDataEthTyped)
		tCpy := cpy.data.(TxInternalDataEthTyped)

		if !reflect.DeepEqual(tOrig.GetAccessList(), tCpy.GetAccessList()) {
			return fmt.Errorf("access list wrong!")
		}
	}

	return nil
}

func TestIsSorted(t *testing.T) {
	signer := LatestSignerForChainID(big.NewInt(1))

	key, _ := crypto.GenerateKey()
	batches := make(Transactions, 10)

	for i := 0; i < 10; i++ {
		batches[i], _ = SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(i)), nil), signer, key)
	}

	// Shuffle transactions.
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(batches), func(i, j int) {
		batches[i], batches[j] = batches[j], batches[i]
	})

	sort.Sort(TxByPriceAndTime(batches))
	assert.True(t, sort.IsSorted(TxByPriceAndTime(batches)))
}

func BenchmarkTxSortByTime30000(b *testing.B) { benchmarkTxSortByTime(b, 30000) }
func BenchmarkTxSortByTime20000(b *testing.B) { benchmarkTxSortByTime(b, 20000) }
func benchmarkTxSortByTime(b *testing.B, size int) {
	signer := LatestSignerForChainID(big.NewInt(1))

	key, _ := crypto.GenerateKey()
	batches := make(Transactions, size)

	for i := 0; i < size; i++ {
		batches[i], _ = SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(i)), nil), signer, key)
	}

	// Shuffle transactions.
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(batches), func(i, j int) {
		batches[i], batches[j] = batches[j], batches[i]
	})

	// Benchmark importing the transactions into the queue
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sort.Sort(TxByPriceAndTime(batches))
	}
}
