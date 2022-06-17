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
	"log"
	"math/big"
	"strconv"
	"testing"
	"time"

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
)

const (
	timeOut                 = 3 * time.Second // timeout of context and event loop for simulated backend.
	DefaultBridgeTxGasLimit = uint64(10000000)
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
		DefaultBridgeTxGasLimit,
		gasPrice,
		nil)

	signedTx, err := auth.Signer(types.LatestSignerForChainID(chainID), auth.From, tx)
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
	_, err := b.RequestKLAYTransfer(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: DefaultBridgeTxGasLimit, Value: new(big.Int).SetUint64(value)}, to, new(big.Int).SetUint64(value), nil)
	if err != nil {
		t.Fatalf("fail to RequestKLAYTransfer %v", err)
	}
}

// SendHandleKLAYTransfer send a handleValueTransfer transaction to the bridge contract.
func SendHandleKLAYTransfer(b *bridge.Bridge, auth *bind.TransactOpts, to common.Address, value uint64, nonce uint64, blockNum uint64, t *testing.T) *types.Transaction {
	tx, err := b.HandleKLAYTransfer(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: DefaultBridgeTxGasLimit}, common.Hash{10}, common.Address{0}, to, big.NewInt(int64(value)), nonce, blockNum, nil)
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
	defer backend.Close()

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
	defer backend.Close()

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
	requestSub, err := b.WatchRequestValueTransfer(nil, requestValueTransferEventCh, nil, nil, nil)
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
	defer backend.Close()

	chargeAmount := big.NewInt(10000000)
	bridgeAccount.Value = chargeAmount
	bridgeAddress, tx, b, err := bridge.DeployBridge(bridgeAccount, backend, false)
	if err != nil {
		t.Fatalf("fail to DeployBridge %v", err)
	}
	backend.Commit()
	WaitMined(tx, backend, t)
	t.Log("1. Bridge is deployed.", "bridgeAddress=", bridgeAddress.String(), "txHash=", tx.Hash().String())

	tx, err = b.RegisterOperator(&bind.TransactOpts{From: bridgeAccount.From, Signer: bridgeAccount.Signer, GasLimit: DefaultBridgeTxGasLimit}, bridgeAccount.From)
	assert.NoError(t, err)
	backend.Commit()
	WaitMined(tx, backend, t)

	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 100)
	handleSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh, nil, nil, nil)
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
	defer backend.Close()

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
	opts := bind.MakeTransactOpts(bridgeAccountKey, big.NewInt(int64(nonce)), DefaultBridgeTxGasLimit, gasPrice)

	tx, err = b.SetCounterPartBridge(opts, common.Address{2})
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, WaitMined(tx, backend, t))

	version, err := b.VERSION(nil)
	assert.Equal(t, uint64(0x1), version)

	allowedTokens, err := b.RegisteredTokens(nil, common.Address{1})
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
	aliceAcc.GasLimit = DefaultBridgeTxGasLimit

	bobKey, _ := crypto.GenerateKey()
	bobAcc := bind.NewKeyedTransactor(bobKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)
	defer backend.Close()

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
	requestSub, err := b.WatchRequestValueTransfer(nil, requestValueTransferEventCh, nil, nil, nil)
	assert.NoError(t, err)
	defer requestSub.Unsubscribe()

	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 10)
	handleSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh, nil, nil, nil)
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
	aliceAcc.GasLimit = DefaultBridgeTxGasLimit

	bobKey, _ := crypto.GenerateKey()
	bobAcc := bind.NewKeyedTransactor(bobKey)

	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)
	defer backend.Close()

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
	requestSub, err := b.WatchRequestValueTransfer(nil, requestValueTransferEventCh, nil, nil, nil)
	assert.NoError(t, err)
	defer requestSub.Unsubscribe()
	requestValueTransferEncodedEventCh := make(chan *bridge.BridgeRequestValueTransferEncoded, 10)
	requestEncodedEvSub, err := b.WatchRequestValueTransferEncoded(nil, requestValueTransferEncodedEventCh, nil, nil, nil)
	assert.NoError(t, err)
	defer requestEncodedEvSub.Unsubscribe()

	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 10)
	handleSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh, nil, nil, nil)
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

		// HandleERC721Transfer
		tx, err = b.HandleERC721Transfer(bridgeAccount, ev.Raw.TxHash, ev.From, ev.To, ev.TokenAddress, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber, "", ev.ExtraData)
		assert.NoError(t, err)
		backend.Commit()
		assert.Nil(t, bind.CheckWaitMined(backend, tx))

	case ev := <-requestValueTransferEncodedEventCh:
		assert.Equal(t, testToken.String(), ev.ValueOrTokenId.String())
		assert.Equal(t, rNonce, ev.RequestNonce)
		assert.Equal(t, erc721Addr, ev.TokenAddress)
		assert.Equal(t, ERC721, ev.TokenType)
		assert.Equal(t, bobAcc.From, ev.To)

		// HandleERC721Transfer
		uri := GetURI(RequestValueTransferEncodedEvent{ev})
		tx, err = b.HandleERC721Transfer(bridgeAccount, ev.Raw.TxHash, ev.From, ev.To, ev.TokenAddress, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber, uri, ev.ExtraData)
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

type bridgeTokenTestENV struct {
	backend    *backends.SimulatedBackend
	operator   *bind.TransactOpts
	tester     *bind.TransactOpts
	bridge     *bridge.Bridge
	erc20      *sctoken.ServiceChainToken
	erc721     *scnft.ServiceChainNFT
	erc20Addr  common.Address
	erc721Addr common.Address
}

func generateBridgeTokenTestEnv(t *testing.T) *bridgeTokenTestENV {
	key, _ := crypto.GenerateKey()
	operator := bind.NewKeyedTransactor(key)
	operator.GasLimit = DefaultBridgeTxGasLimit

	testKey, _ := crypto.GenerateKey()
	tester := bind.NewKeyedTransactor(testKey)
	tester.GasLimit = DefaultBridgeTxGasLimit

	alloc := blockchain.GenesisAlloc{operator.From: {Balance: big.NewInt(params.KLAY)}, tester.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy Bridge
	bridgeAddr, tx, b, err := bridge.DeployBridge(operator, backend, true)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Deploy ERC20
	erc20Addr, tx, erc20, err := sctoken.DeployServiceChainToken(operator, backend, bridgeAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Register ERC20
	tx, err = b.RegisterToken(operator, erc20Addr, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Charge token to tester
	tx, err = erc20.Transfer(operator, tester.From, big.NewInt(100))
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Deploy ERC721
	erc721Addr, tx, erc721, err := scnft.DeployServiceChainNFT(operator, backend, bridgeAddr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Register ERC721
	tx, err = b.RegisterToken(operator, erc721Addr, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	// Mint token to tester
	tx, err = erc721.RegisterBulk(operator, tester.From, big.NewInt(0), big.NewInt(10))
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	return &bridgeTokenTestENV{
		backend,
		operator,
		tester,
		b,
		erc20,
		erc721,
		erc20Addr,
		erc721Addr,
	}
}

// TestBridgeContract_RegisterToken checks belows.
// - RegisterToken works well
// - DeregisterToken works well
func TestBridgeContract_RegisterToken(t *testing.T) {
	env := generateBridgeTokenTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	operator := env.operator
	b := env.bridge
	erc20Addr := env.erc20Addr
	erc721Addr := env.erc721Addr

	// Deregister erc20
	tx, err := b.DeregisterToken(operator, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	cpToken, err := b.RegisteredTokens(nil, erc20Addr)
	assert.NoError(t, err)
	assert.Equal(t, common.Address{}, cpToken)

	cpToken, err = b.RegisteredTokens(nil, erc721Addr)
	assert.NoError(t, err)
	assert.Equal(t, erc721Addr, cpToken)

	// Deregister erc721Addr
	tx, err = b.DeregisterToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	cpToken, err = b.RegisteredTokens(nil, erc721Addr)
	assert.NoError(t, err)
	assert.Equal(t, common.Address{}, cpToken)

	tokens, err := b.GetRegisteredTokenList(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tokens))

	// Register erc20
	tx, err = b.RegisterToken(operator, erc20Addr, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	cpToken, err = b.RegisteredTokens(nil, erc20Addr)
	assert.NoError(t, err)
	assert.Equal(t, erc20Addr, cpToken)

	tx, err = b.RegisterToken(operator, erc20Addr, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// Register erc721
	tx, err = b.RegisterToken(operator, erc721Addr, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	cpToken, err = b.RegisteredTokens(nil, erc721Addr)
	assert.NoError(t, err)
	assert.Equal(t, erc721Addr, cpToken)

	tx, err = b.RegisterToken(operator, erc721Addr, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// Deregister erc721Addr
	tx, err = b.DeregisterToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	cpToken, err = b.RegisteredTokens(nil, erc721Addr)
	assert.NoError(t, err)
	assert.Equal(t, common.Address{}, cpToken)
}

// TestBridgeContract_InitStatus checks initial lock status.
func TestBridgeContract_InitStatus(t *testing.T) {
	env := generateBridgeTokenTestEnv(t)
	defer env.backend.Close()

	b := env.bridge
	erc20Addr := env.erc20Addr
	erc721Addr := env.erc721Addr

	// initial value check
	isLocked, err := b.LockedTokens(nil, erc20Addr)
	assert.NoError(t, err)
	assert.Equal(t, false, isLocked)

	isLocked, err = b.LockedTokens(nil, erc721Addr)
	assert.NoError(t, err)
	assert.Equal(t, false, isLocked)

	isLocked, err = b.IsLockedKLAY(nil)
	assert.NoError(t, err)
	assert.Equal(t, false, isLocked)
}

// TestBridgeContract_InitRequest checks the following:
// - the request value transfer can be allowed after registering it.
func TestBridgeContract_InitRequest(t *testing.T) {
	env := generateBridgeTokenTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	operator := env.operator
	tester := env.tester
	b := env.bridge
	erc20 := env.erc20
	erc721 := env.erc721

	// check to allow value transfer
	tx, err := erc20.RequestValueTransfer(tester, big.NewInt(1), operator.From, big.NewInt(0), nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	tx, err = erc721.RequestValueTransfer(tester, big.NewInt(1), operator.From, nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	tester.Value = big.NewInt(1)
	tx, err = b.RequestKLAYTransfer(tester, tester.From, big.NewInt(1), nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))
	tester.Value = nil
}

// TestBridgeContract_TokenLock checks the following:
// - the token can be lock to prevent value transfer requests.
func TestBridgeContract_TokenLock(t *testing.T) {
	env := generateBridgeTokenTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	operator := env.operator
	tester := env.tester
	b := env.bridge
	erc20 := env.erc20
	erc721 := env.erc721
	erc20Addr := env.erc20Addr
	erc721Addr := env.erc721Addr

	// lock token
	tx, err := b.LockToken(operator, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.LockToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.LockKLAY(operator)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// check value after locking
	isLocked, err := b.LockedTokens(nil, erc20Addr)
	assert.NoError(t, err)
	assert.Equal(t, true, isLocked)

	isLocked, err = b.LockedTokens(nil, erc721Addr)
	assert.NoError(t, err)
	assert.Equal(t, true, isLocked)

	isLocked, err = b.IsLockedKLAY(nil)
	assert.NoError(t, err)
	assert.Equal(t, true, isLocked)

	// check to prevent value transfer
	tx, err = erc20.RequestValueTransfer(tester, big.NewInt(1), operator.From, big.NewInt(0), nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.NotNil(t, bind.CheckWaitMined(backend, tx))

	tx, err = erc721.RequestValueTransfer(tester, big.NewInt(1), operator.From, nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.NotNil(t, bind.CheckWaitMined(backend, tx))

	tester.Value = big.NewInt(1)
	tx, err = b.RequestKLAYTransfer(tester, tester.From, big.NewInt(1), nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.NotNil(t, bind.CheckWaitMined(backend, tx))
	tester.Value = nil
}

// TestBridgeContract_TokenLockFail checks the following:
// - testing the case locking token is fail.
func TestBridgeContract_TokenLockFail(t *testing.T) {
	env := generateBridgeTokenTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	operator := env.operator
	tester := env.tester
	b := env.bridge
	erc20Addr := env.erc20Addr
	erc721Addr := env.erc721Addr

	// fail locking token by invalid owner.
	tx, err := b.LockToken(tester, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	tx, err = b.LockToken(tester, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	tx, err = b.LockKLAY(tester)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// fail locking unregistered token.
	testAddr := common.BytesToAddress([]byte("unregistered"))
	tx, err = b.LockToken(operator, testAddr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// fail locking token if it is already locked.
	tx, err = b.LockToken(operator, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.LockToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.LockKLAY(operator)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.LockToken(operator, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	tx, err = b.LockToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	tx, err = b.LockKLAY(operator)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
}

// TestBridgeContract_TokenUnlockFail checks the following:
// - testing the case unlocking token is fail.
func TestBridgeContract_TokenUnlockFail(t *testing.T) {
	env := generateBridgeTokenTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	operator := env.operator
	b := env.bridge
	erc20Addr := env.erc20Addr
	erc721Addr := env.erc721Addr

	// fail locking unregistered token.
	testAddr := common.BytesToAddress([]byte("unregistered"))
	tx, err := b.UnlockToken(operator, testAddr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// fail locking token if it is already unlocked.
	tx, err = b.UnlockToken(operator, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	tx, err = b.UnlockToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	tx, err = b.UnlockKLAY(operator)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
}

// TestBridgeContract_CheckValueTransferAfterUnLock checks the following:
// - the token can be unlock to allow value transfer requests.
func TestBridgeContract_CheckValueTransferAfterUnLock(t *testing.T) {
	env := generateBridgeTokenTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	operator := env.operator
	tester := env.tester
	b := env.bridge
	erc20 := env.erc20
	erc721 := env.erc721
	erc20Addr := env.erc20Addr
	erc721Addr := env.erc721Addr

	// lock token
	tx, err := b.LockToken(operator, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.LockToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.LockKLAY(operator)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// unlock token
	tx, err = b.UnlockToken(operator, erc20Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.UnlockToken(operator, erc721Addr)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = b.UnlockKLAY(operator)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// check value after unlocking
	isLocked, err := b.LockedTokens(nil, erc20Addr)
	assert.NoError(t, err)
	assert.Equal(t, false, isLocked)

	isLocked, err = b.LockedTokens(nil, erc721Addr)
	assert.NoError(t, err)
	assert.Equal(t, false, isLocked)

	isLocked, err = b.IsLockedKLAY(nil)
	assert.NoError(t, err)
	assert.Equal(t, false, isLocked)

	// check to allow value transfer
	tx, err = erc20.RequestValueTransfer(tester, big.NewInt(1), operator.From, big.NewInt(0), nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	tx, err = erc721.RequestValueTransfer(tester, big.NewInt(1), operator.From, nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	tester.Value = big.NewInt(1)
	tx, err = b.RequestKLAYTransfer(tester, tester.From, big.NewInt(1), nil)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))
	tester.Value = nil
}

// TestBridgeRequestHandleGasUsed tests the gas used of handle function
// with the gap of lowerHandle nonce and handle nonce.
func TestBridgeRequestHandleGasUsed(t *testing.T) {
	// Generate a new random account and a funded simulator
	authKey, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(authKey)
	auth.GasLimit = DefaultBridgeTxGasLimit

	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)

	bobKey, _ := crypto.GenerateKey()
	bob := bind.NewKeyedTransactor(bobKey)

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		alice.From: {Balance: big.NewInt(params.KLAY)},
		auth.From:  {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	var err error

	// Deploy a bridge contract
	auth.Value = big.NewInt(100000000000)
	_, _, b, err := bridge.DeployBridge(auth, sim, false)
	assert.NoError(t, err)
	sim.Commit() // block
	auth.Value = big.NewInt(0)

	// Register the owner as a signer
	_, err = b.RegisterOperator(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, auth.From)
	assert.NoError(t, err)
	sim.Commit() // block

	// Subscribe Bridge Contract
	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 10)
	handleValueTransferSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh, nil, nil, nil)
	defer handleValueTransferSub.Unsubscribe()

	handleFunc := func(nonce int) {
		hTx, err := b.HandleKLAYTransfer(auth, common.HexToHash(strconv.Itoa(nonce)), alice.From, bob.From, big.NewInt(1), uint64(nonce), uint64(1+nonce), nil)
		assert.NoError(t, err)
		sim.Commit()

		receipt, err := bind.WaitMined(context.Background(), sim, hTx)
		assert.NoError(t, err)
		assert.Equal(t, uint(0x1), receipt.Status)

		select {
		case ev := <-handleValueTransferEventCh:
			t.Log("Handle value transfer event",
				"handleNonce", ev.HandleNonce,
				"lowerHandleNonce", ev.LowerHandleNonce,
				"gasUsed", receipt.GasUsed,
				"status", receipt.Status)
		case <-time.After(1 * time.Second):
			if receipt != nil {
				t.Log("handle event omitted Tx gas used=", receipt.GasUsed)
			}
			t.Fatal("handle event omitted")
		}
	}

	// handle 0 ~ 499 nonce
	for i := 0; i < 500; i++ {
		handleFunc(i)
	}

	lowerHandleNonce, _ := b.LowerHandleNonce(nil)
	assert.Equal(t, uint64(500), lowerHandleNonce)
	upperHandleNonce, _ := b.UpperHandleNonce(nil)
	assert.Equal(t, uint64(499), upperHandleNonce)

	// handle 501 ~ 999 nonce
	for i := 501; i < 1000; i++ {
		handleFunc(i)
	}

	lowerHandleNonce, _ = b.LowerHandleNonce(nil)
	assert.Equal(t, uint64(500), lowerHandleNonce)
	upperHandleNonce, _ = b.UpperHandleNonce(nil)
	assert.Equal(t, uint64(999), upperHandleNonce)

	//This 500 nonce handle checks whether the handle transaction which has a loop failed.
	handleFunc(500)

	lowerHandleNonce, _ = b.LowerHandleNonce(nil)
	assert.Equal(t, uint64(701), lowerHandleNonce)
	upperHandleNonce, _ = b.UpperHandleNonce(nil)
	assert.Equal(t, uint64(999), upperHandleNonce)
}

// TestBridgeMaxOperator tests
// - the gas used of handle function with max operators.
// - preventing to add more operators than the limit.
func TestBridgeMaxOperatorHandleTxGasUsed(t *testing.T) {
	// Generate a new random account and a funded simulator
	maxOperator := 12

	var authList []*bind.TransactOpts
	for i := 0; i <= maxOperator; i++ {
		authKey, _ := crypto.GenerateKey()
		authList = append(authList, bind.NewKeyedTransactor(authKey))
		authList[i].GasLimit = DefaultBridgeTxGasLimit
	}
	auth := authList[0]

	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)

	bobKey, _ := crypto.GenerateKey()
	bob := bind.NewKeyedTransactor(bobKey)

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		alice.From: {Balance: big.NewInt(params.KLAY)},
		auth.From:  {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	var err error

	// Deploy a bridge contract
	auth.Value = big.NewInt(100000000000)
	_, _, b, err := bridge.DeployBridge(auth, sim, false)
	assert.NoError(t, err)
	sim.Commit() // block
	auth.Value = big.NewInt(0)

	// Register the owner as a signer
	for _, a := range authList {
		_, err = b.RegisterOperator(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, a.From)
		assert.NoError(t, err)
	}

	_, err = b.SetOperatorThreshold(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, voteTypeValueTransfer, uint8(maxOperator))
	assert.NoError(t, err)

	sim.Commit() // block

	// test preventing more operators than the limit.
	{
		tx, err := b.RegisterOperator(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, authList[maxOperator].From)
		assert.NoError(t, err)
		sim.Commit() // block

		receipt, err := bind.WaitMined(context.Background(), sim, tx)
		assert.NoError(t, err)
		assert.Equal(t, types.ReceiptStatusErrExecutionReverted, receipt.Status)
	}

	// Subscribe Bridge Contract
	handleValueTransferEventCh := make(chan *bridge.BridgeHandleValueTransfer, 10)
	handleValueTransferSub, err := b.WatchHandleValueTransfer(nil, handleValueTransferEventCh, nil, nil, nil)
	defer handleValueTransferSub.Unsubscribe()

	handleFunc := func(a *bind.TransactOpts, nonce int) {
		hTx, err := b.HandleKLAYTransfer(a, common.HexToHash(strconv.Itoa(nonce)), alice.From, bob.From, big.NewInt(1), uint64(nonce), uint64(1+nonce), nil)
		assert.NoError(t, err)
		sim.Commit()

		receipt, err := bind.WaitMined(context.Background(), sim, hTx)
		assert.NoError(t, err)
		assert.Equal(t, uint(0x1), receipt.Status)

		t.Log("Handle value transfer tx receipt", "gasUsed", receipt.GasUsed, "status", receipt.Status)
	}

	for i := 0; i < maxOperator; i++ {
		handleFunc(authList[i], 0)
	}

	select {
	case ev := <-handleValueTransferEventCh:
		t.Log("Handle value transfer event",
			"handleNonce", ev.HandleNonce,
			"lowerHandleNonce", ev.LowerHandleNonce)
	case <-time.After(1 * time.Second):
		t.Fatal("handle event omitted")
	}
}

// TestBridgeThresholdLimit tests preventing the invalid threshold value.
func TestBridgeThresholdLimit(t *testing.T) {
	// Generate a new random account and a funded simulator
	maxOperator := 12

	var authList []*bind.TransactOpts
	for i := 0; i < maxOperator; i++ {
		authKey, _ := crypto.GenerateKey()
		authList = append(authList, bind.NewKeyedTransactor(authKey))
		authList[i].GasLimit = DefaultBridgeTxGasLimit
	}
	auth := authList[0]

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		auth.From: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	var err error

	// Deploy a bridge contract
	auth.Value = big.NewInt(100000000000)
	_, _, b, err := bridge.DeployBridge(auth, sim, false)
	assert.NoError(t, err)
	sim.Commit() // block
	auth.Value = big.NewInt(0)

	// Register the owner as a signer
	for i, a := range authList {
		_, err = b.RegisterOperator(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, a.From)
		assert.NoError(t, err)
		sim.Commit() // block

		// bigger threshold than operators
		{
			tx, err := b.SetOperatorThreshold(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, voteTypeValueTransfer, uint8(i+2))
			assert.NoError(t, err)
			sim.Commit() // block

			receipt, err := bind.WaitMined(context.Background(), sim, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusErrExecutionReverted, receipt.Status)
		}

		// zero threshold
		{
			tx, err := b.SetOperatorThreshold(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, voteTypeValueTransfer, uint8(0))
			assert.NoError(t, err)
			sim.Commit() // block

			receipt, err := bind.WaitMined(context.Background(), sim, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusErrExecutionReverted, receipt.Status)
		}

		// same threshold with operators
		{
			tx, err := b.SetOperatorThreshold(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, voteTypeValueTransfer, uint8(i+1))
			assert.NoError(t, err)
			sim.Commit() // block

			receipt, err := bind.WaitMined(context.Background(), sim, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
		}

		// lower threshold than operator
		if i > 0 {
			tx, err := b.SetOperatorThreshold(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, voteTypeValueTransfer, uint8(i))
			assert.NoError(t, err)
			sim.Commit() // block

			receipt, err := bind.WaitMined(context.Background(), sim, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
		}
	}
}
