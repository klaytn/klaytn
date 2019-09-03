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

package sc

import (
	"context"
	"errors"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/contracts/extbridge"
	sctoken "github.com/klaytn/klaytn/contracts/sc_erc20"
	scnft "github.com/klaytn/klaytn/contracts/sc_erc721"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"testing"
	"time"
)

const (
	gasLimit uint64 = 3000000         // gasLimit for contract transaction.
	timeOut         = 3 * time.Second // timeout of context and event loop for simulated backend.
)

// WaitMined waits the tx receipt until the timeout.
func WaitMined(tx *types.Transaction, backend bind.DeployBackend, t *testing.T) error {
	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), timeOut)
	defer cancelTimeout()

	receipt, err := bind.WaitMined(timeoutContext, backend, tx)
	if err != nil {
		t.Fatal("Failed to WaitMined.", "err", err, "txHash", tx.Hash().String(), "status", receipt.Status)
		return err
	}

	if receipt == nil {
		return errors.New("receipt not found")
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Log("receipt", "status", receipt.Status)
		return errors.New("not successful receipt")
	}
	return nil
}

// TransferSignedTx sends the transaction to transfer KLAY from auth to `to` and waits the execution of the transaction.
func TransferSignedTx(auth *bind.TransactOpts, backend *backends.SimulatedBackend, to common.Address, value *big.Int, t *testing.T) (common.Hash, *big.Int, error) {
	ctx := context.Background()

	nonce, err := backend.NonceAt(ctx, auth.From, nil)
	assert.Equal(t, err, nil)

	chainID, err := backend.ChainID(ctx)
	assert.Equal(t, err, nil)

	gasPrice, err := backend.SuggestGasPrice(ctx)
	assert.Equal(t, err, nil)

	tx := types.NewTransaction(
		nonce,
		to,
		value,
		gasLimit,
		gasPrice,
		nil)

	signedTx, err := auth.Signer(types.NewEIP155Signer(chainID), auth.From, tx)
	assert.Equal(t, err, nil)

	fee := big.NewInt(0)

	err = backend.SendTransaction(ctx, signedTx)
	assert.Equal(t, err, nil)

	backend.Commit()

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, backend, signedTx)
	if err != nil {
		log.Fatalf("WaitMined time out %v", err)
	}

	fee.Mul(big.NewInt(int64(receipt.GasUsed)), gasPrice)

	return tx.Hash(), fee, nil
}

// RequestKLAYTransfer sends a requestValueTransfer transaction to the bridge contract.
func RequestKLAYTransfer(b *bridge.Bridge, auth *bind.TransactOpts, to common.Address, value uint64, t *testing.T) {
	_, err := b.RequestKLAYTransfer(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: gasLimit, Value: new(big.Int).SetUint64(value)}, to, new(big.Int).SetUint64(value), nil)
	if err != nil {
		t.Fatalf("fail to RequestKLAYTransfer %v", err)
	}
}

// SendHandleKLAYTransfer send a handleValueTransfer transaction to the bridge contract.
func SendHandleKLAYTransfer(b *bridge.Bridge, auth *bind.TransactOpts, to common.Address, value uint64, nonce uint64, blockNum uint64, t *testing.T) *types.Transaction {
	tx, err := b.HandleKLAYTransfer(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: gasLimit}, common.Hash{10}, common.Address{0}, to, big.NewInt(int64(value)), nonce, blockNum, nil)
	if err != nil {
		t.Fatalf("fail to SendHandleKLAYTransfer %v", err)
		return nil
	}
	return tx
}

// TestBridgeDeployWithKLAY checks to the state/contract balance of the bridge deployed.
func TestBridgeDeployWithKLAY(t *testing.T) {
	bridgeAccountKey, _ := crypto.GenerateKey()
	bridgeAccount := bind.NewKeyedTransactor(bridgeAccountKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)

	chargeAmount := big.NewInt(10000000)
	bridgeAccount.Value = chargeAmount
	bridgeAddress, tx, _, err := bridge.DeployBridge(bridgeAccount, backend, false)
	if err != nil {
		t.Fatalf("fail to DeployBridge %v", err)
	}
	backend.Commit()
	WaitMined(tx, backend, t)

	balanceContract, err := backend.BalanceAt(nil, bridgeAddress, nil)
	if err != nil {
		t.Fatalf("fail to GetKLAY %v", err)
	}

	balanceState, err := backend.BalanceAt(context.Background(), bridgeAddress, nil)
	if err != nil {
		t.Fatal("failed to BalanceAt")
	}

	assert.Equal(t, chargeAmount, balanceState)
	assert.Equal(t, chargeAmount, balanceContract)
}

// TestBridgeRequestValueTransferNonce checks the bridge emit events with serialized nonce.
func TestBridgeRequestValueTransferNonce(t *testing.T) {
	bridgeAccountKey, _ := crypto.GenerateKey()
	bridgeAccount := bind.NewKeyedTransactor(bridgeAccountKey)

	testAccKey, _ := crypto.GenerateKey()
	testAcc := bind.NewKeyedTransactor(testAccKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)

	chargeAmount := big.NewInt(10000000)
	bridgeAccount.Value = chargeAmount
	addr, tx, b, err := bridge.DeployBridge(bridgeAccount, backend, false)
	if err != nil {
		t.Fatalf("fail to DeployBridge %v", err)
	}
	backend.Commit()
	WaitMined(tx, backend, t)
	t.Log("1. Bridge is deployed.", "addr=", addr.String(), "txHash=", tx.Hash().String())

	requestValueTransferEventCh := make(chan *bridge.BridgeRequestValueTransfer, 100)
	requestSub, err := b.WatchRequestValueTransfer(nil, requestValueTransferEventCh)
	defer requestSub.Unsubscribe()
	if err != nil {
		t.Fatalf("fail to WatchHandleValueTransfer %v", err)
	}
	t.Log("2. Bridge is subscribed.")

	RequestKLAYTransfer(b, bridgeAccount, testAcc.From, 1, t)
	backend.Commit()

	expectedNonce := uint64(0)

loop:
	for {
		select {
		case ev := <-requestValueTransferEventCh:
			assert.Equal(t, expectedNonce, ev.RequestNonce)

			if expectedNonce == 1000 {
				return
			}
			expectedNonce++

			// TODO-Klaytn added more request token/NFT transfer cases,
			RequestKLAYTransfer(b, bridgeAccount, testAcc.From, 1, t)
			backend.Commit()

		case err := <-requestSub.Err():
			t.Log("Contract Event Loop Running Stop by requestSub.Err()", "err", err)
			break loop

		case <-time.After(timeOut):
			t.Log("Contract Event Loop Running Stop by timeout")
			break loop
		}
	}

	t.Fatal("fail to check monotone increasing nonce", "lastNonce", expectedNonce)
}

// TestBridgeHandleValueTransferNonceAndBlockNumber checks the following:
// - the bridge allows the handle value transfer with an arbitrary nonce.
// - the bridge keeps lower handle nonce for the recovery.
// - the bridge correctly stores and returns the block number.
func TestBridgeHandleValueTransferNonceAndBlockNumber(t *testing.T) {
	bridgeAccountKey, _ := crypto.GenerateKey()
	bridgeAccount := bind.NewKeyedTransactor(bridgeAccountKey)

	testAccKey, _ := crypto.GenerateKey()
	testAcc := bind.NewKeyedTransactor(testAccKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)

	chargeAmount := big.NewInt(10000000)
	bridgeAccount.Value = chargeAmount
	bridgeAddress, tx, b, err := bridge.DeployBridge(bridgeAccount, backend, false)
	if err != nil {
		t.Fatalf("fail to DeployBridge %v", err)
	}
	backend.Commit()
	WaitMined(tx, backend, t)
	t.Log("1. Bridge is deployed.", "bridgeAddress=", bridgeAddress.String(), "txHash=", tx.Hash().String())

	tx, err = b.RegisterOperator(&bind.TransactOpts{From: bridgeAccount.From, Signer: bridgeAccount.Signer, GasLimit: gasLimit}, bridgeAccount.From)
	assert.NoError(t, err)
	backend.Commit()
	WaitMined(tx, backend, t)

	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 100)
	handleSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh)
	assert.NoError(t, err)
	defer handleSub.Unsubscribe()
	if err != nil {
		t.Fatalf("fail to DepositKLAY %v", err)
	}
	t.Log("2. Bridge is subscribed.")

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	testCount := uint64(1000)
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)
	tx = SendHandleKLAYTransfer(b, bridgeAccount, testAcc.From, transferAmount, sentNonce, sentBlockNumber, t)
	backend.Commit()

	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), timeOut)
	defer cancelTimeout()

	receipt, err := bind.WaitMined(timeoutContext, backend, tx)

	if err != nil {
		t.Fatal("Failed to WaitMined.", "err", err, "txHash", tx.Hash().String(), "status", receipt.Status)
	}

loop:
	for {
		select {
		case ev := <-handleValueTransferEventCh:
			assert.Equal(t, sentNonce, ev.HandleNonce)

			if sentNonce == testCount {
				bal, err := backend.BalanceAt(context.Background(), testAcc.From, nil)
				assert.NoError(t, err)
				assert.Equal(t, bal, big.NewInt(int64(transferAmount*(testCount-nonceOffset+1))))

				lowerHandleNonce, err := b.LowerHandleNonce(nil)
				assert.NoError(t, err)
				assert.Equal(t, lowerHandleNonce, uint64(0))
				return
			}
			sentNonce++
			sentBlockNumber++

			SendHandleKLAYTransfer(b, bridgeAccount, testAcc.From, transferAmount, sentNonce, sentBlockNumber, t)
			backend.Commit()

			resultBlockNumber, err := b.RecoveryBlockNumber(nil)
			assert.NoError(t, err)

			resultHandleNonce, err := b.UpperHandleNonce(nil)
			assert.NoError(t, err)

			assert.Equal(t, sentNonce, resultHandleNonce)
			assert.Equal(t, uint64(1), resultBlockNumber)

		case err := <-handleSub.Err():
			t.Log("Contract Event Loop Running Stop by handleSub.Err()", "err", err)
			break loop

		case <-time.After(timeOut):
			t.Log("Contract Event Loop Running Stop by timeout")
			break loop
		}
	}

	t.Fatal("fail to check monotone increasing nonce", "lastNonce", sentNonce)
}

// TestBridgePublicVariables checks the results of the public variables.
func TestBridgePublicVariables(t *testing.T) {
	bridgeAccountKey, _ := crypto.GenerateKey()
	bridgeAccount := bind.NewKeyedTransactor(bridgeAccountKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)

	chargeAmount := big.NewInt(10000000)
	bridgeAccount.Value = chargeAmount
	bridgeAddress, tx, b, err := bridge.DeployBridge(bridgeAccount, backend, false)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, WaitMined(tx, backend, t))

	balanceContract, err := backend.BalanceAt(nil, bridgeAddress, nil)
	assert.NoError(t, err)
	assert.Equal(t, chargeAmount, balanceContract)

	ctx := context.Background()
	nonce, err := backend.NonceAt(ctx, bridgeAccount.From, nil)
	gasPrice, err := backend.SuggestGasPrice(ctx)
	opts := bind.MakeTransactOpts(bridgeAccountKey, big.NewInt(int64(nonce)), gasLimit, gasPrice)

	tx, err = b.SetCounterPartBridge(opts, common.Address{2})
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, WaitMined(tx, backend, t))

	version, err := b.VERSION(nil)
	assert.Equal(t, uint64(0x1), version)

	allowedTokens, err := b.AllowedTokens(nil, common.Address{1})
	assert.Equal(t, common.Address{0}, allowedTokens)

	counterpartBridge, err := b.CounterpartBridge(nil)
	assert.Equal(t, common.Address{2}, counterpartBridge)

	hnonce, err := b.LowerHandleNonce(nil)
	assert.Equal(t, uint64(0), hnonce)

	owner, err := b.IsOwner(&bind.CallOpts{From: bridgeAccount.From})
	assert.Equal(t, true, owner)

	notOwner, err := b.IsOwner(&bind.CallOpts{From: common.Address{1}})
	assert.Equal(t, false, notOwner)

	isRunning, err := b.IsRunning(nil)
	assert.Equal(t, true, isRunning)

	lastBN, err := b.RecoveryBlockNumber(nil)
	assert.Equal(t, uint64(0x1), lastBN)

	bridgeOwner, err := b.Owner(nil)
	assert.Equal(t, bridgeAccount.From, bridgeOwner)

	rnonce, err := b.RequestNonce(nil)
	assert.Equal(t, uint64(0), rnonce)
}

// TestExtendedBridgeAndCallbackERC20 checks the following:
// - the extBridge can call a callback contract method from ERC20 value transfer.
func TestExtendedBridgeAndCallbackERC20(t *testing.T) {
	bridgeAccountKey, _ := crypto.GenerateKey()
	bridgeAccount := bind.NewKeyedTransactor(bridgeAccountKey)

	aliceKey, _ := crypto.GenerateKey()
	aliceAcc := bind.NewKeyedTransactor(aliceKey)
	aliceAcc.GasLimit = gasLimit

	bobKey, _ := crypto.GenerateKey()
	bobAcc := bind.NewKeyedTransactor(bobKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy extBridge
	bridgeAddr, tx, eb, err := extbridge.DeployExtBridge(bridgeAccount, backend, true)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Deploy token
	erc20Addr, tx, erc20, err := sctoken.DeployServiceChainToken(bridgeAccount, backend, bridgeAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Register token
	tx, err = eb.RegisterToken(bridgeAccount, erc20Addr, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Charge token to Alice
	testToken := big.NewInt(100000)
	tx, err = erc20.Transfer(bridgeAccount, aliceAcc.From, testToken)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Give minter role to bridge contract
	tx, err = erc20.AddMinter(bridgeAccount, bridgeAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Deploy callback contract
	callbackAddr, tx, cb, err := extbridge.DeployCallback(bridgeAccount, backend)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Set callback address to ExtBridge contract
	tx, err = eb.SetCallback(bridgeAccount, callbackAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Subscribe callback contract event
	registerOfferEventCh := make(chan *extbridge.CallbackRegisteredOffer, 10)
	registerOfferEventSub, err := cb.WatchRegisteredOffer(nil, registerOfferEventCh)
	assert.NoError(t, err)
	defer registerOfferEventSub.Unsubscribe()

	// Subscribe bridge contract events
	b, err := bridge.NewBridge(bridgeAddr, backend) // create base bridge contract object, not extBridge object
	assert.NoError(t, err)

	requestValueTransferEventCh := make(chan *bridge.BridgeRequestValueTransfer, 10)
	requestSub, err := b.WatchRequestValueTransfer(nil, requestValueTransferEventCh)
	assert.NoError(t, err)
	defer requestSub.Unsubscribe()

	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 10)
	handleSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh)
	assert.NoError(t, err)
	defer handleSub.Unsubscribe()

	// Approve / RequestSellERC20
	rNonce := uint64(0)
	amount := big.NewInt(1000)
	offerPrice := big.NewInt(1006)

	tx, err = erc20.Approve(aliceAcc, bridgeAddr, amount)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Fail case
	tx, err = eb.RequestERC20Transfer(aliceAcc, erc20Addr, bobAcc.From, amount, common.Big0, nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Error(t, bind.CheckWaitMined(backend, tx))

	// Success case
	tx, err = eb.RequestSellERC20(aliceAcc, erc20Addr, bobAcc.From, amount, common.Big0, offerPrice)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Check request request event
	select {
	case ev := <-requestValueTransferEventCh:
		assert.Equal(t, amount.String(), ev.ValueOrTokenId.String())
		assert.Equal(t, rNonce, ev.RequestNonce)
		assert.Equal(t, erc20Addr, ev.TokenAddress)
		assert.Equal(t, ERC20, ev.TokenType)
		assert.Equal(t, bobAcc.From, ev.To)

		// HandleERC20Transfer
		tx, err = b.HandleERC20Transfer(bridgeAccount, ev.Raw.TxHash, ev.From, ev.To, ev.TokenAddress, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber, ev.ExtraData)
		assert.NoError(t, err)
		backend.Commit()
		assert.Nil(t, bind.CheckWaitMined(backend, tx))

	case <-time.After(time.Second):
		t.Fatalf("requestValueTransferEvent was not found.")
	}

	// Check handle request event
	select {
	case ev := <-handleValueTransferEventCh:
		assert.Equal(t, amount.String(), ev.ValueOrTokenId.String())
		assert.Equal(t, rNonce, ev.HandleNonce)
		assert.Equal(t, erc20Addr, ev.TokenAddress)
		assert.Equal(t, ERC20, ev.TokenType)
		assert.Equal(t, callbackAddr, ev.To)

	case <-time.After(time.Second):
		t.Fatalf("handleValueTransferEvent was not found.")
	}

	// Check Callback event
	select {
	case ev := <-registerOfferEventCh:
		assert.Equal(t, amount.String(), ev.ValueOrID.String())
		assert.Equal(t, offerPrice.String(), ev.Price.String())
		assert.Equal(t, erc20Addr, ev.TokenAddress)

	case <-time.After(time.Second):
		t.Fatalf("registerOfferEvent was not found.")
	}
}

// TestExtendedBridgeAndCallbackERC721 checks the following:
// - the extBridge can call a callback contract method from ERC721 value transfer.
func TestExtendedBridgeAndCallbackERC721(t *testing.T) {
	bridgeAccountKey, _ := crypto.GenerateKey()
	bridgeAccount := bind.NewKeyedTransactor(bridgeAccountKey)

	aliceKey, _ := crypto.GenerateKey()
	aliceAcc := bind.NewKeyedTransactor(aliceKey)
	aliceAcc.GasLimit = gasLimit

	bobKey, _ := crypto.GenerateKey()
	bobAcc := bind.NewKeyedTransactor(bobKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy extBridge
	bridgeAddr, tx, eb, err := extbridge.DeployExtBridge(bridgeAccount, backend, true)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Deploy token
	erc721Addr, tx, erc721, err := scnft.DeployServiceChainNFT(bridgeAccount, backend, bridgeAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Register token
	tx, err = eb.RegisterToken(bridgeAccount, erc721Addr, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Charge token to Alice
	testToken := big.NewInt(100000)
	tx, err = erc721.MintWithTokenURI(bridgeAccount, aliceAcc.From, testToken, "")
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Give minter role to bridge contract
	tx, err = erc721.AddMinter(bridgeAccount, bridgeAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Deploy callback contract
	callbackAddr, tx, cb, err := extbridge.DeployCallback(bridgeAccount, backend)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Set callback address to ExtBridge contract
	tx, err = eb.SetCallback(bridgeAccount, callbackAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Subscribe callback contract event
	registerOfferEventCh := make(chan *extbridge.CallbackRegisteredOffer, 10)
	registerOfferEventSub, err := cb.WatchRegisteredOffer(nil, registerOfferEventCh)
	assert.NoError(t, err)
	defer registerOfferEventSub.Unsubscribe()

	// Subscribe bridge contract events
	b, err := bridge.NewBridge(bridgeAddr, backend) // create base bridge contract object, not extBridge object
	assert.NoError(t, err)

	requestValueTransferEventCh := make(chan *bridge.BridgeRequestValueTransfer, 10)
	requestSub, err := b.WatchRequestValueTransfer(nil, requestValueTransferEventCh)
	assert.NoError(t, err)
	defer requestSub.Unsubscribe()

	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 10)
	handleSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh)
	assert.NoError(t, err)
	defer handleSub.Unsubscribe()

	// Approve / RequestSellERC721
	rNonce := uint64(0)
	offerPrice := big.NewInt(1006)

	tx, err = erc721.Approve(aliceAcc, bridgeAddr, testToken)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Fail case
	tx, err = eb.RequestERC721Transfer(aliceAcc, erc721Addr, bobAcc.From, testToken, nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Error(t, bind.CheckWaitMined(backend, tx))

	// Success case
	tx, err = eb.RequestSellERC721(aliceAcc, erc721Addr, bobAcc.From, testToken, offerPrice)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Check request request event
	select {
	case ev := <-requestValueTransferEventCh:
		assert.Equal(t, testToken.String(), ev.ValueOrTokenId.String())
		assert.Equal(t, rNonce, ev.RequestNonce)
		assert.Equal(t, erc721Addr, ev.TokenAddress)
		assert.Equal(t, ERC721, ev.TokenType)
		assert.Equal(t, bobAcc.From, ev.To)

		// HandleERC20Transfer
		tx, err = b.HandleERC721Transfer(bridgeAccount, ev.Raw.TxHash, ev.From, ev.To, ev.TokenAddress, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber, ev.Uri, ev.ExtraData)
		assert.NoError(t, err)
		backend.Commit()
		assert.Nil(t, bind.CheckWaitMined(backend, tx))

	case <-time.After(time.Second):
		t.Fatalf("requestValueTransferEvent was not found.")
	}

	// Check handle request event
	select {
	case ev := <-handleValueTransferEventCh:
		assert.Equal(t, testToken.String(), ev.ValueOrTokenId.String())
		assert.Equal(t, rNonce, ev.HandleNonce)
		assert.Equal(t, erc721Addr, ev.TokenAddress)
		assert.Equal(t, ERC721, ev.TokenType)
		assert.Equal(t, callbackAddr, ev.To)

	case <-time.After(time.Second):
		t.Fatalf("handleValueTransferEvent was not found.")
	}

	// Check Callback event
	select {
	case ev := <-registerOfferEventCh:
		assert.Equal(t, testToken.String(), ev.ValueOrID.String())
		assert.Equal(t, offerPrice.String(), ev.Price.String())
		assert.Equal(t, erc721Addr, ev.TokenAddress)

	case <-time.After(time.Second):
		t.Fatalf("registerOfferEvent was not found.")
	}
}
