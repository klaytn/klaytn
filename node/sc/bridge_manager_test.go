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
	"fmt"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/contracts/sc_erc20"
	"github.com/klaytn/klaytn/contracts/sc_erc721"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"
)

// WaitGroupWithTimeOut waits the given wait group until the timout duration.
func WaitGroupWithTimeOut(wg *sync.WaitGroup, duration time.Duration, t *testing.T) {
	c := make(chan struct{})
	go func() {
		wg.Wait()
		c <- struct{}{}
	}()
	t.Log("start to wait group")
	select {
	case <-c:
		t.Log("waiting group is done")
	case <-time.After(duration):
		t.Fatal("timed out waiting group")
	}
}

// CheckReceipt can check if the tx receipt has expected status.
func CheckReceipt(b bind.DeployBackend, tx *types.Transaction, duration time.Duration, expectedStatus uint, t *testing.T) {
	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), duration)
	defer cancelTimeout()

	receipt, err := bind.WaitMined(timeoutContext, b, tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedStatus, receipt.Status)
}

// TestBridgeManager tests the event/method of Token/NFT/Bridge contracts.
// TODO-Klaytn-Servicechain needs to refine this test.
// - consider main/service chain simulated backend.
// - separate each test
func TestBridgeManager(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(6)

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bacc, _ := NewBridgeAccounts(config.DataDir)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	pAuth := bacc.cAccount.GetTransactOpts()
	cAuth := bacc.pAccount.GetTransactOpts()

	// Generate a new random account and a funded simulator
	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)

	bobKey, _ := crypto.GenerateKey()
	bob := bind.NewKeyedTransactor(bobKey)

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		alice.From:            {Balance: big.NewInt(params.KLAY)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KLAY)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
	}
	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	bridgeManager, err := NewBridgeManager(sc)

	testToken := big.NewInt(123)
	testKLAY := big.NewInt(321)

	// 1. Deploy Bridge Contract
	addr, err := bridgeManager.DeployBridgeTest(sim, false)
	if err != nil {
		log.Fatalf("Failed to deploy new bridge contract: %v", err)
	}
	bridgeInfo, _ := bridgeManager.GetBridgeInfo(addr)
	bridge := bridgeInfo.bridge
	fmt.Println("===== BridgeContract Addr ", addr.Hex())
	sim.Commit() // block

	// 2. Deploy Token Contract
	tokenAddr, tx, token, err := sctoken.DeployServiceChainToken(alice, sim, addr)
	if err != nil {
		log.Fatalf("Failed to DeployGXToken: %v", err)
	}
	sim.Commit() // block

	// 3. Deploy NFT Contract
	nftTokenID := uint64(4438)
	nftAddr, tx, nft, err := scnft.DeployServiceChainNFT(alice, sim, addr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}
	sim.Commit() // block

	// Register the owner as a signer
	_, err = bridge.RegisterOperator(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, cAuth.From)
	assert.NoError(t, err)
	sim.Commit() // block

	// Register tokens on the bridgeInfo
	bridgeInfo.RegisterToken(tokenAddr, tokenAddr)
	bridgeInfo.RegisterToken(nftAddr, nftAddr)

	// Register tokens on the bridge
	bridge.RegisterToken(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, tokenAddr, tokenAddr)
	bridge.RegisterToken(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, nftAddr, nftAddr)
	sim.Commit() // block

	cTokenAddr, err := bridge.AllowedTokens(nil, tokenAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cTokenAddr, tokenAddr)
	cNftAddr, err := bridge.AllowedTokens(nil, nftAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cNftAddr, nftAddr)

	balance, _ := sim.BalanceAt(context.Background(), pAuth.From, nil)
	fmt.Printf("auth(%v) KLAY balance : %v\n", pAuth.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), cAuth.From, nil)
	fmt.Printf("auth2(%v) KLAY balance : %v\n", cAuth.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), alice.From, nil)
	fmt.Printf("auth3(%v) KLAY balance : %v\n", alice.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), bob.From, nil)
	fmt.Printf("auth4(%v) KLAY balance : %v\n", bob.From.String(), balance)

	// 4. Subscribe Bridge Contract
	bridgeManager.SubscribeEvent(addr)

	requestValueTransferEventCh := make(chan *RequestValueTransferEvent)
	handleValueTransferEventCh := make(chan *HandleValueTransferEvent)
	bridgeManager.SubscribeRequestEvent(requestValueTransferEventCh)
	bridgeManager.SubscribeHandleEvent(handleValueTransferEventCh)

	go func() {
		for {
			select {
			case ev := <-requestValueTransferEventCh:
				fmt.Println("Request Event",
					"type", ev.TokenType,
					"amount", ev.ValueOrTokenId,
					"from", ev.From.String(),
					"to", ev.To.String(),
					"contract", ev.Raw.Address.String(),
					"token", ev.TokenAddress.String(),
					"requestNonce", ev.RequestNonce)

				done, err := bridge.HandledRequestTx(nil, ev.Raw.TxHash)
				assert.NoError(t, err)
				assert.Equal(t, false, done)

				// insert the value transfer request event to the bridge info's event list.
				bridgeInfo.AddRequestValueTransferEvents([]*RequestValueTransferEvent{ev})

				// handle the value transfer request event in the event list.
				bridgeInfo.processingPendingRequestEvents()

				sim.Commit() // block
				wg.Done()
				done, err = bridge.HandledRequestTx(nil, ev.Raw.TxHash)
				assert.NoError(t, err)
				assert.Equal(t, true, done)

			case ev := <-handleValueTransferEventCh:
				fmt.Println("Handle value transfer event",
					"bridgeAddr", ev.Raw.Address.Hex(),
					"type", ev.TokenType,
					"amount", ev.ValueOrTokenId,
					"owner", ev.To.String(),
					"contract", ev.Raw.Address.String(),
					"token", ev.TokenAddress.String(),
					"handleNonce", ev.HandleNonce)
				wg.Done()
			}
		}
	}()

	// 6. Register (Mint) an NFT to Alice
	{
		tx, err = nft.MintWithTokenURI(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, alice.From, big.NewInt(int64(nftTokenID)), "testURI")
		assert.NoError(t, err)
		fmt.Println("Register NFT Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		owner, err := nft.OwnerOf(nil, big.NewInt(int64(nftTokenID)))
		assert.Equal(t, nil, err)
		assert.Equal(t, alice.From, owner)
	}

	// 7. Request ERC20 Transfer from Alice to Bob
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, testToken, bob.From, big.NewInt(0), nil)
		assert.NoError(t, err)
		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	}

	// 8. RequestKLAYTransfer from Alice to Bob
	{
		tx, err = bridge.RequestKLAYTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, Value: testKLAY, GasLimit: testGasLimit}, bob.From, testKLAY, nil)
		assert.NoError(t, err)
		fmt.Println("DepositKLAY Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9. Request NFT transfer from Alice to Bob
	{
		tx, err = nft.RequestValueTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, big.NewInt(int64(nftTokenID)), bob.From, nil)
		assert.NoError(t, err)
		fmt.Println("nft.RequestValueTransfer Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// Wait a few second for wait group
	WaitGroupWithTimeOut(&wg, 3*time.Second, t)

	// 10. Check Token balance
	{
		balance, err = token.BalanceOf(nil, bob.From)
		assert.Equal(t, nil, err)
		assert.Equal(t, testToken.String(), balance.String())
	}

	// 11. Check KLAY balance
	{
		balance, err = sim.BalanceAt(context.Background(), bob.From, nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, testKLAY.String(), balance.String())
	}

	// 12. Check NFT owner
	{
		owner, err := nft.OwnerOf(nil, big.NewInt(int64(nftTokenID)))
		assert.Equal(t, nil, err)
		assert.Equal(t, bob.From, owner)
	}

	bridgeManager.Stop()
}

// TestBridgeManagerWithFee tests the KLAY/ERC20 transfer with fee.
func TestBridgeManagerWithFee(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(7 * 2)

	// Generate a new random account and a funded simulator
	AliceKey, _ := crypto.GenerateKey()
	Alice := bind.NewKeyedTransactor(AliceKey)

	BobKey, _ := crypto.GenerateKey()
	Bob := bind.NewKeyedTransactor(BobKey)

	receiverKey, _ := crypto.GenerateKey()
	receiver := bind.NewKeyedTransactor(receiverKey)

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bacc, _ := NewBridgeAccounts(config.DataDir)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	pAuth := bacc.cAccount.GetTransactOpts()
	cAuth := bacc.pAccount.GetTransactOpts()

	// Create Simulated backend
	initialValue := int64(10000000000)
	alloc := blockchain.GenesisAlloc{
		Alice.From:            {Balance: big.NewInt(initialValue)},
		bacc.cAccount.address: {Balance: big.NewInt(initialValue)},
		bacc.pAccount.address: {Balance: big.NewInt(initialValue)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
	}
	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	bridgeManager, err := NewBridgeManager(sc)

	testToken := int64(100000)
	testKLAY := int64(100000)
	KLAYFee := int64(500)
	ERC20Fee := int64(500)

	// 1. Deploy Bridge Contract
	pBridgeAddr, err := bridgeManager.DeployBridgeTest(sim, false)
	assert.NoError(t, err)
	pBridgeInfo, _ := bridgeManager.GetBridgeInfo(pBridgeAddr)
	pBridge := pBridgeInfo.bridge
	fmt.Println("===== BridgeContract Addr ", pBridgeAddr.Hex())
	sim.Commit() // block

	// 2. Deploy Token Contract
	tokenAddr, tx, token, err := sctoken.DeployServiceChainToken(pAuth, sim, pBridgeAddr)
	assert.NoError(t, err)
	sim.Commit() // block

	// Set value transfer fee
	{
		nilReceiver, err := pBridge.FeeReceiver(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, common.Address{}, nilReceiver)
	}

	pBridge.SetFeeReceiver(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, receiver.From)
	sim.Commit() // block

	{
		recv, err := pBridge.FeeReceiver(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, receiver.From, recv)
	}

	{
		fee, err := pBridge.FeeOfKLAY(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, big.NewInt(0).String(), fee.String())
	}

	{
		fee, err := pBridge.FeeOfERC20(nil, tokenAddr)
		assert.Equal(t, nil, err)
		assert.Equal(t, big.NewInt(0).String(), fee.String())
	}

	cn, err := pBridge.ConfigurationNonce(nil)
	assert.NoError(t, err)
	_, err = pBridge.RegisterOperator(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, cAuth.From)
	assert.NoError(t, err)
	pBridge.SetKLAYFee(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, big.NewInt(KLAYFee), cn)
	pBridge.SetERC20Fee(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, tokenAddr, big.NewInt(ERC20Fee), cn+1)
	sim.Commit() // block

	{
		fee, err := pBridge.FeeOfKLAY(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, KLAYFee, fee.Int64())
	}

	{
		fee, err := pBridge.FeeOfERC20(nil, tokenAddr)
		assert.Equal(t, nil, err)
		assert.Equal(t, ERC20Fee, fee.Int64())
	}

	// Register tokens on the bridgeInfo
	pBridgeInfo.RegisterToken(tokenAddr, tokenAddr)

	// Register tokens on the bridge
	pBridge.RegisterToken(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, tokenAddr, tokenAddr)
	sim.Commit() // block

	cTokenAddr, err := pBridge.AllowedTokens(nil, tokenAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cTokenAddr, tokenAddr)

	balance, _ := sim.BalanceAt(context.Background(), Alice.From, nil)
	fmt.Printf("Alice(%v) KLAY balance : %v\n", Alice.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), Bob.From, nil)
	fmt.Printf("Bob(%v) KLAY balance : %v\n", Bob.From.String(), balance)

	// 4. Subscribe Bridge Contract
	bridgeManager.SubscribeEvent(pBridgeAddr)

	requestValueTransferEventCh := make(chan *RequestValueTransferEvent)
	handleValueTransferEventCh := make(chan *HandleValueTransferEvent)
	bridgeManager.SubscribeRequestEvent(requestValueTransferEventCh)
	bridgeManager.SubscribeHandleEvent(handleValueTransferEventCh)

	go func() {
		for {
			select {
			case ev := <-requestValueTransferEventCh:
				fmt.Println("Request value transfer event",
					"type", ev.TokenType,
					"amount", ev.ValueOrTokenId,
					"from", ev.From.String(),
					"to", ev.To.String(),
					"contract", ev.Raw.Address.String(),
					"token", ev.TokenAddress.String(),
					"requestNonce", ev.RequestNonce,
					"fee", ev.Fee.String())

				// insert the value transfer request event to the bridge info's event list.
				pBridgeInfo.AddRequestValueTransferEvents([]*RequestValueTransferEvent{ev})

				// handle the value transfer request event in the event list.
				pBridgeInfo.processingPendingRequestEvents()

				sim.Commit() // block
				wg.Done()

			case ev := <-handleValueTransferEventCh:
				fmt.Println("Handle value transfer event",
					"bridgeAddr", ev.Raw.Address.Hex(),
					"type", ev.TokenType,
					"amount", ev.ValueOrTokenId,
					"owner", ev.To.String(),
					"contract", ev.Raw.Address.String(),
					"token", ev.TokenAddress.String(),
					"handleNonce", ev.HandleNonce)
				wg.Done()
			}
		}
	}()

	// 5. transfer from parentAcc to Alice for charging and check balances
	{
		tx, err = token.Transfer(&bind.TransactOpts{From: pAuth.From, Signer: pAuth.Signer, GasLimit: testGasLimit}, Alice.From, big.NewInt(initialValue))
		if err != nil {
			log.Fatalf("Failed to Transfer for charging: %v", err)
		}
		fmt.Println("Transfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		balance, err = token.BalanceOf(nil, pAuth.From)
		assert.Equal(t, nil, err)
		fmt.Println("parentAcc token balance", balance.String())

		balance, err = token.BalanceOf(nil, Alice.From)
		assert.Equal(t, nil, err)
		fmt.Println("Alice token balance", balance.String())

		balance, err = token.BalanceOf(nil, Bob.From)
		assert.Equal(t, nil, err)
		fmt.Println("Bob token balance", balance.String())
	}

	// 7-1. Request ERC20 Transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee), nil)
		assert.NoError(t, err)
		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 7-2. Request ERC20 Transfer from Alice to Bob with insufficient zero feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(0), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 7-3. Request ERC20 Transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee-1), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 7-4. Request ERC20 Transfer from Alice to Bob with enough feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee+1), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 8-1. Approve/Request ERC20 Transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee), nil)
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 8-2. Approve/Request ERC20 Transfer from Alice to Bob with insufficient zero feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(0), nil)
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 8-3. Approve/Request ERC20 Transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee-1))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee-1), nil)
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 8-4. Approve/Request ERC20 Transfer from Alice to Bob with enough feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee+1))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee+1), nil)
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-1. Request KLAY transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY + KLAYFee), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY), nil)
		if err != nil {
			log.Fatalf("Failed to RequestKLAYTransfer: %v", err)
		}
		fmt.Println("RequestKLAYTransfer Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-2. Request KLAY transfer from Alice to Bob with zero feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 9-3. Request KLAY transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY + (KLAYFee - 1)), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 9-4. Request KLAY transfer from Alice to Bob with enough feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY + (KLAYFee + 1)), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-4. Request KLAY transfer from Alice to Alice through () payable method
	{
		nonce, _ := sim.PendingNonceAt(context.Background(), Alice.From)
		gasPrice, _ := sim.SuggestGasPrice(context.Background())
		unsignedTx := types.NewTransaction(nonce, pBridgeAddr, big.NewInt(testKLAY+KLAYFee), testGasLimit, gasPrice, []byte{})

		chainID, _ := sim.ChainID(context.Background())
		tx, err = types.SignTx(unsignedTx, types.NewEIP155Signer(chainID), AliceKey)
		sim.SendTransaction(context.Background(), tx)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// Wait a few second for wait group
	WaitGroupWithTimeOut(&wg, 3*time.Second, t)

	// 10. Check Token balance
	{
		balance, err = token.BalanceOf(nil, Alice.From)
		assert.Equal(t, nil, err)
		fmt.Println("Alice token balance", balance.String())
		assert.Equal(t, initialValue-(testToken+ERC20Fee)*4, balance.Int64())

		balance, err = token.BalanceOf(nil, Bob.From)
		assert.Equal(t, nil, err)
		fmt.Println("Bob token balance", balance.String())
		assert.Equal(t, testToken*4, balance.Int64())

		balance, err = token.BalanceOf(nil, receiver.From)
		assert.Equal(t, nil, err)
		fmt.Println("Fee receiver token balance", balance.String())
		assert.Equal(t, ERC20Fee*4, balance.Int64())
	}

	// 11. Check KLAY balance
	{
		balance, _ = sim.BalanceAt(context.Background(), Alice.From, nil)
		fmt.Println("Alice KLAY balance :", balance)
		assert.Equal(t, initialValue-(testKLAY+KLAYFee)*2-KLAYFee, balance.Int64())

		balance, _ = sim.BalanceAt(context.Background(), Bob.From, nil)
		fmt.Println("Bob KLAY balance :", balance)
		assert.Equal(t, big.NewInt(testKLAY*2).String(), balance.String())

		balance, _ = sim.BalanceAt(context.Background(), receiver.From, nil)
		fmt.Println("receiver KLAY balance :", balance)
		assert.Equal(t, KLAYFee*3, balance.Int64())
	}

	bridgeManager.Stop()
}

// TestBasicJournal tests basic journal functionality.
func TestBasicJournal(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)

	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(os.TempDir())
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KLAY)},
		auth2.From:            {Balance: big.NewInt(params.KLAY)},
		auth4.From:            {Balance: big.NewInt(params.KLAY)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KLAY)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}
	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)

	localAddr, err := bm.DeployBridgeTest(sim, true)
	assert.NoError(t, err)
	remoteAddr, err := bm.DeployBridgeTest(sim, false)
	assert.NoError(t, err)

	bm.SetJournal(localAddr, remoteAddr)

	ps := sc.BridgePeerSet()
	ps.peers["test"] = nil

	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	localInfo, ok := bm.GetBridgeInfo(localAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, localInfo.subscribed)

	remoteInfo, ok := bm.GetBridgeInfo(remoteAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, remoteInfo.subscribed)
}

// TestMethodRestoreBridges tests restoring bridges from the journal.
func TestMethodRestoreBridges(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true
	config.VTRecoveryInterval = 60

	bacc, _ := NewBridgeAccounts(os.TempDir())
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KLAY)},
		auth2.From:            {Balance: big.NewInt(params.KLAY)},
		auth4.From:            {Balance: big.NewInt(params.KLAY)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KLAY)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager and deploy bridge contract.
	bm, _ := NewBridgeManager(sc)

	var bridgeAddrs [4]common.Address
	for i := 0; i < 4; i++ {
		if i%2 == 0 {
			bridgeAddrs[i], err = bm.DeployBridgeTest(sim, true)
		} else {
			bridgeAddrs[i], err = bm.DeployBridgeTest(sim, false)
		}
		if err != nil {
			t.Fatal("deploy bridge test failed", bridgeAddrs[i])
		}
		bm.DeleteBridgeInfo(bridgeAddrs[i])
	}
	sim.Commit()

	// Set journal
	bm.SetJournal(bridgeAddrs[0], bridgeAddrs[1])
	bm.journal.cache[bridgeAddrs[0]].Subscribed = true
	bm.SetJournal(bridgeAddrs[2], bridgeAddrs[3])
	bm.journal.cache[bridgeAddrs[2]].Subscribed = true

	ps := sc.BridgePeerSet()
	ps.peers["test"] = nil

	// Call RestoreBridges
	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	// Duplicated RestoreBridges
	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	// Case 1: check bridge contract creation.
	for i := 0; i < 4; i++ {
		info, _ := bm.GetBridgeInfo(bridgeAddrs[i])
		assert.NotEqual(t, nil, info.bridge)
	}

	// Case 2: check subscription
	for i := 0; i < 4; i++ {
		info, _ := bm.GetBridgeInfo(bridgeAddrs[i])
		assert.Equal(t, true, info.subscribed)
	}

	// Case 3: check recovery
	recovery1 := bm.recoveries[bridgeAddrs[0]]
	assert.NotEqual(t, nil, recovery1)
	recovery1.Start()
	assert.Equal(t, nil, recovery1.WaitRunningStatus(true, 5*time.Second))
	recovery2 := bm.recoveries[bridgeAddrs[2]]
	assert.NotEqual(t, nil, recovery2)
	recovery2.Start()
	assert.Equal(t, nil, recovery2.WaitRunningStatus(true, 5*time.Second))

	bm.stopAllRecoveries()
	bm.Stop()
}

// TestMethodGetAllBridge tests a method GetAllBridge.
func TestMethodGetAllBridge(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	testBridge1 := common.BytesToAddress([]byte("test1"))
	testBridge2 := common.BytesToAddress([]byte("test2"))

	bm.journal.insert(testBridge1, testBridge2)
	bm.journal.insert(testBridge2, testBridge1)

	bridges := bm.GetAllBridge()
	assert.Equal(t, 2, len(bridges))

	bm.Stop()
}

// TestErrorDuplication tests if duplication of journal insertion is ignored or not.
func TestErrorDuplication(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	localAddr := common.BytesToAddress([]byte("test1"))
	remoteAddr := common.BytesToAddress([]byte("test2"))

	err = bm.journal.insert(localAddr, remoteAddr)
	assert.Equal(t, nil, err)
	err = bm.journal.insert(remoteAddr, localAddr)
	assert.Equal(t, nil, err)

	// try duplicated insert.
	err = bm.journal.insert(localAddr, remoteAddr)
	assert.NotEqual(t, nil, err)
	err = bm.journal.insert(remoteAddr, localAddr)
	assert.NotEqual(t, nil, err)

	// check cache size for checking duplication
	bridges := bm.GetAllBridge()
	assert.Equal(t, 2, len(bridges))

	bm.Stop()
}

// TestMethodSetJournal tests if duplication of journal insertion is ignored or not.
func TestMethodSetJournal(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	localAddr := common.BytesToAddress([]byte("test1"))
	remoteAddr := common.BytesToAddress([]byte("test2"))

	// Simple insert case
	err = bm.SetJournal(localAddr, remoteAddr)
	assert.Equal(t, nil, err)

	// Error case
	err = bm.SetJournal(localAddr, remoteAddr)
	assert.NotEqual(t, nil, err)

	// Check the number of bridge elements for checking duplication
	bridges := bm.GetAllBridge()
	assert.Equal(t, 1, len(bridges))

	bm.Stop()
}

// TestErrorDuplicatedSetBridgeInfo tests if duplication of bridge info insertion.
func TestErrorDuplicatedSetBridgeInfo(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(os.TempDir())
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KLAY)},
		auth2.From:            {Balance: big.NewInt(params.KLAY)},
		auth4.From:            {Balance: big.NewInt(params.KLAY)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KLAY)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager
	bm, err := NewBridgeManager(sc)
	addr, err := bm.DeployBridgeTest(sim, false)
	bridgeInfo, _ := bm.GetBridgeInfo(addr)

	// Try to call duplicated SetBridgeInfo
	err = bm.SetBridgeInfo(addr, bridgeInfo.bridge, common.Address{}, nil, sc.bridgeAccounts.pAccount, false, false)
	assert.NotEqual(t, nil, err)
	bm.Stop()
}

// TestScenarioSubUnsub tests subscription and unsubscription scenario.
func TestScenarioSubUnsub(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(os.TempDir())
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KLAY)},
		auth2.From:            {Balance: big.NewInt(params.KLAY)},
		auth4.From:            {Balance: big.NewInt(params.KLAY)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KLAY)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)

	localAddr, err := bm.DeployBridgeTest(sim, true)
	if err != nil {
		t.Fatal("deploy bridge test failed", localAddr)
	}

	bm.SubscribeEvent(localAddr)
	bridgeInfo, ok := bm.GetBridgeInfo(localAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, true, bridgeInfo.subscribed)
	bm.UnsubscribeEvent(localAddr)
	assert.Equal(t, false, bridgeInfo.subscribed)

	// Journal is irrelevant to the bridge unsubscription.
	journal := bm.journal.cache[localAddr]
	assert.NotEqual(t, nil, journal)
}

// TestErrorEmptyAccount tests empty account error in case of journal insertion.
func TestErrorEmptyAccount(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	localAddr := common.BytesToAddress([]byte("test1"))
	remoteAddr := common.BytesToAddress([]byte("test2"))

	err = bm.journal.insert(localAddr, common.Address{})
	assert.NotEqual(t, nil, err)

	err = bm.journal.insert(common.Address{}, remoteAddr)
	assert.NotEqual(t, nil, err)

	bm.Stop()
}

// TestErrorDupSubscription tests duplicated subscription error.
func TestErrorDupSubscription(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(os.TempDir())
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KLAY)},
		auth2.From:            {Balance: big.NewInt(params.KLAY)},
		auth4.From:            {Balance: big.NewInt(params.KLAY)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KLAY)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	var err error
	sc.handler, err = NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// 1. Prepare manager and subscribe event
	bm, err := NewBridgeManager(sc)

	addr, err := bm.DeployBridgeTest(sim, false)
	bridgeInfo, _ := bm.GetBridgeInfo(addr)
	bridge := bridgeInfo.bridge
	fmt.Println("===== BridgeContract Addr ", addr.Hex())
	sim.Commit() // block

	bm.bridges[addr], err = NewBridgeInfo(nil, addr, bridge, common.Address{}, nil, bacc.cAccount, true, true)

	bm.journal.cache[addr] = &BridgeJournal{addr, addr, true}

	bm.SubscribeEvent(addr)
	err = bm.SubscribeEvent(addr)
	assert.NotEqual(t, nil, err)

	bm.Stop()
}

// for TestMethod
func (bm *BridgeManager) DeployBridgeTest(backend *backends.SimulatedBackend, local bool) (common.Address, error) {
	if local {
		acc := bm.subBridge.bridgeAccounts.cAccount
		addr, bridge, err := bm.deployBridgeTest(acc, backend, true)
		if err != nil {
			return common.Address{}, err
		}
		err = bm.SetBridgeInfo(addr, bridge, common.Address{}, nil, acc, local, false)
		if err != nil {
			return common.Address{}, err
		}
		return addr, err
	} else {
		acc := bm.subBridge.bridgeAccounts.pAccount
		addr, bridge, err := bm.deployBridgeTest(acc, backend, false)
		if err != nil {
			return common.Address{}, err
		}
		err = bm.SetBridgeInfo(addr, bridge, common.Address{}, nil, acc, local, false)
		if err != nil {
			return common.Address{}, err
		}
		return addr, err
	}
}

func (bm *BridgeManager) deployBridgeTest(acc *accountInfo, backend *backends.SimulatedBackend, modeMintBurn bool) (common.Address, *bridge.Bridge, error) {
	auth := acc.GetTransactOpts()
	auth.Value = big.NewInt(10000)
	addr, tx, contract, err := bridge.DeployBridge(auth, backend, modeMintBurn)
	if err != nil {
		logger.Error("", "err", err)
		return common.Address{}, nil, err
	}
	logger.Info("Bridge is deploying on CurrentChain", "addr", addr, "txHash", tx.Hash().String())

	backend.Commit()

	// TODO-Klaytn-Servicechain needs to support WaitMined
	//timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancelTimeout()
	//
	//receipt, err := bind.WaitMined(timeoutContext, backend, tx)
	//if err != nil {
	//	log.Fatal("Failed to deploy.", "err", err, "txHash", tx.Hash().String(), "status", receipt.Status)
	//	return common.Address{}, nil, err
	//}
	//fmt.Println("deployBridge is executed.", "addr", addr.String(), "txHash", tx.Hash().String())

	return addr, contract, nil
}
