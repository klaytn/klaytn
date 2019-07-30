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
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"os"
	"path"
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
// And It tests the nonce error case of bridge deploy (#2284)
// TODO-Klaytn-Servicechain needs to refine this test.
// - consider main/service chain simulated backend.
// - separate each test
func TestBridgeManager(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(6)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key3, _ := crypto.GenerateKey()
	auth3 := bind.NewKeyedTransactor(key3)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: big.NewInt(params.KLAY)}, auth2.From: {Balance: big.NewInt(params.KLAY)}, auth4.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	balance2, _ := sim.BalanceAt(context.Background(), auth2.From, big.NewInt(0))
	fmt.Println("after reward, balance :", balance2)

	config := &SCConfig{}
	config.nodekey = key
	config.chainkey = key2
	config.DataDir = os.TempDir()

	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.MainChainAccountAddr = &chainKeyAddr

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		bridgeAccountManager: bam,
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
	tokenAddr, tx, token, err := sctoken.DeployServiceChainToken(auth, sim, addr)
	if err != nil {
		log.Fatalf("Failed to DeployGXToken: %v", err)
	}
	sim.Commit() // block

	// 3. Deploy NFT Contract
	nftTokenID := uint64(4438)
	nftAddr, tx, nft, err := scnft.DeployServiceChainNFT(auth, sim, addr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}
	sim.Commit() // block

	// Register tokens on the bridge
	rn, err := bridge.RequestNonce(nil)
	if err != nil {
		log.Fatalf("Failed to register token: %v", err)
	}
	bridge.RegisterToken(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, tokenAddr, tokenAddr, rn)
	sim.Commit() // block
	if err != nil {
		log.Fatalf("Failed to register token: %v", err)
	}
	bridge.RegisterToken(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, nftAddr, nftAddr, rn+1)
	sim.Commit() // block

	cTokenAddr, err := bridge.AllowedTokens(nil, tokenAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cTokenAddr, tokenAddr)
	cNftAddr, err := bridge.AllowedTokens(nil, nftAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cNftAddr, nftAddr)

	balance, _ := sim.BalanceAt(context.Background(), auth.From, nil)
	fmt.Printf("auth(%v) KLAY balance : %v\n", auth.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), auth2.From, nil)
	fmt.Printf("auth2(%v) KLAY balance : %v\n", auth2.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), auth3.From, nil)
	fmt.Printf("auth3(%v) KLAY balance : %v\n", auth3.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), auth4.From, nil)
	fmt.Printf("auth4(%v) KLAY balance : %v\n", auth4.From.String(), balance)

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

				switch ev.TokenType {
				case KLAY:
					tx, err := bridge.HandleKLAYTransfer(
						&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit},
						ev.From, ev.To, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber)
					if err != nil {
						log.Fatalf("Failed to HandleKLAYTransfer: %v", err)
					}
					fmt.Println("WithdrawKLAY Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block

				case ERC20:
					tx, err := bridge.HandleERC20Transfer(
						&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit},
						ev.From, ev.To, tokenAddr, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber)
					if err != nil {
						log.Fatalf("Failed to HandleERC20Transfer: %v", err)
					}
					fmt.Println("HandleERC20Transfer Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block

				case ERC721:
					owner, err := nft.OwnerOf(&bind.CallOpts{From: auth.From}, big.NewInt(int64(nftTokenID)))
					assert.Equal(t, nil, err)
					fmt.Println("NFT owner before HandleERC721Transfer: ", owner.String())

					tx, err := bridge.HandleERC721Transfer(
						&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit},
						ev.From, ev.To, nftAddr, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber, ev.Uri)
					if err != nil {
						log.Fatalf("Failed to HandleERC721Transfer: %v", err)
					}
					fmt.Println("HandleERC721Transfer Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block
				}

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

	// 5. transfer from auth to auth2 for charging and check balances
	{
		tx, err = token.Transfer(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, auth2.From, testToken)
		if err != nil {
			log.Fatalf("Failed to Transfer for charging: %v", err)
		}
		fmt.Println("Transfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth.From}, auth.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth token balance", balance.String())

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth2.From}, auth2.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth2 token balance", balance.String())

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth3.From}, auth3.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth3 token balance", balance.String())

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth4.From}, auth4.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth4 token balance", balance.String())
	}

	// 6. Register (Mint) an NFT to Auth4
	{
		tx, err = nft.MintWithTokenURI(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, auth4.From, big.NewInt(int64(nftTokenID)), "testURI")
		if err != nil {
			log.Fatalf("Failed to Register NFT: %v", err)
		}
		fmt.Println("Register NFT Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		balance, err = nft.BalanceOf(&bind.CallOpts{From: auth.From}, auth4.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth4 NFT balance", balance.String())
		fmt.Println("auth4 address", auth4.From.String())
		owner, err := nft.OwnerOf(&bind.CallOpts{From: auth.From}, big.NewInt(int64(nftTokenID)))
		assert.Equal(t, nil, err)
		fmt.Println("NFT owner after registering", owner.String())
	}

	// 7. RequestValueTransfer from auth2 to auth3
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, testToken, auth3.From, big.NewInt(0))
		if err != nil {
			log.Fatalf("Failed to SafeTransferAndCall: %v", err)
		}
		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	}

	// 8. DepositKLAY from auth to auth3
	{
		tx, err = bridge.RequestKLAYTransfer(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, Value: testKLAY, GasLimit: testGasLimit}, auth3.From, testKLAY)
		if err != nil {
			log.Fatalf("Failed to DepositKLAY: %v", err)
		}
		fmt.Println("DepositKLAY Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9. Request NFT value transfer from auth4 to auth3
	{
		tx, err = nft.RequestValueTransfer(&bind.TransactOpts{From: auth4.From, Signer: auth4.Signer, GasLimit: testGasLimit}, big.NewInt(int64(nftTokenID)), auth3.From)
		if err != nil {
			log.Fatalf("Failed to nft.RequestValueTransfer: %v", err)
		}
		fmt.Println("nft.RequestValueTransfer Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// Wait a few second for wait group
	WaitGroupWithTimeOut(&wg, 3*time.Second, t)

	// 10. Check Token balance
	{
		balance, err = token.BalanceOf(&bind.CallOpts{From: auth.From}, auth.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth token balance", balance.String())
		balance, err = token.BalanceOf(&bind.CallOpts{From: auth2.From}, auth2.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth2 token balance", balance.String())
		balance, err = token.BalanceOf(&bind.CallOpts{From: auth3.From}, auth3.From)
		assert.Equal(t, nil, err)
		fmt.Println("auth3 token balance", balance.String())

		if balance.Cmp(testToken) != 0 {
			t.Fatal("testToken is mismatched,", "expected", testToken.String(), "result", balance.String())
		}
	}

	// 11. Check KLAY balance
	{
		balance, _ = sim.BalanceAt(context.Background(), auth.From, nil)
		fmt.Println("auth KLAY balance :", balance)

		balance, _ = sim.BalanceAt(context.Background(), auth2.From, nil)
		fmt.Println("auth2 KLAY balance :", balance)

		balance, _ = sim.BalanceAt(context.Background(), auth3.From, nil)
		fmt.Println("auth3 KLAY balance :", balance)

		if balance.Cmp(testKLAY) != 0 {
			t.Fatal("testKLAY is mismatched,", "expected", testKLAY.String(), "result", balance.String())
		}
	}

	// 12. Check NFT owner
	{
		owner, err := nft.OwnerOf(&bind.CallOpts{From: auth.From}, big.NewInt(int64(nftTokenID)))
		assert.Equal(t, nil, err)
		fmt.Println("NFT owner", owner.String())
		if owner != auth3.From {
			t.Fatal("NFT owner is mismatched", "expeted", auth3.From.String(), "result", owner.String())
		}
	}

	// 13. Nonce check on deploy error
	{

		nonce, err := sim.NonceAt(context.Background(), bam.scAccount.address, nil)
		if err != nil {
			t.Fatal("failed to sim.NonceAt", err)
		}
		bam.scAccount.SetNonce(nonce)

		addr2, err := bridgeManager.DeployBridgeNonceTest(sim)
		if err != nil {
			log.Fatalf("Failed to deploy new bridge contract: %v %v", err, addr2)
		}
	}

	bridgeManager.Stop()
}

// TestBridgeManagerWithFee tests the KLAY/ERC20 transfer with fee.
func TestBridgeManagerWithFee(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(7 * 2)

	// Generate a new random account and a funded simulator
	parentKey, _ := crypto.GenerateKey()
	parentAcc := bind.NewKeyedTransactor(parentKey)

	childKey, _ := crypto.GenerateKey()
	childAcc := bind.NewKeyedTransactor(childKey)

	AliceKey, _ := crypto.GenerateKey()
	Alice := bind.NewKeyedTransactor(AliceKey)

	BobKey, _ := crypto.GenerateKey()
	Bob := bind.NewKeyedTransactor(BobKey)

	receiverKey, _ := crypto.GenerateKey()
	receiver := bind.NewKeyedTransactor(receiverKey)

	initialValue := int64(10000000000)
	alloc := blockchain.GenesisAlloc{
		parentAcc.From: {Balance: big.NewInt(initialValue)},
		childAcc.From:  {Balance: big.NewInt(initialValue)},
		Alice.From:     {Balance: big.NewInt(initialValue)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	config := &SCConfig{}
	config.nodekey = parentKey
	config.chainkey = childKey
	config.DataDir = os.TempDir()

	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.MainChainAccountAddr = &chainKeyAddr

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		bridgeAccountManager: bam,
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
	if err != nil {
		log.Fatalf("Failed to deploy new bridge contract: %v", err)
	}
	pBridgeInfo, _ := bridgeManager.GetBridgeInfo(pBridgeAddr)
	pBridge := pBridgeInfo.bridge
	fmt.Println("===== BridgeContract Addr ", pBridgeAddr.Hex())
	sim.Commit() // block

	// 2. Deploy Token Contract
	tokenAddr, tx, token, err := sctoken.DeployServiceChainToken(parentAcc, sim, pBridgeAddr)
	if err != nil {
		log.Fatalf("Failed to DeployGXToken: %v", err)
	}
	sim.Commit() // block

	// Set value transfer fee
	{
		nilReceiver, err := pBridge.FeeReceiver(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, common.Address{}, nilReceiver)
	}

	gn, err := pBridge.GovernanceNonces(nil, TxTypeGovernance)
	assert.NoError(t, err)
	pBridge.SetFeeReceiver(&bind.TransactOpts{From: childAcc.From, Signer: childAcc.Signer, GasLimit: testGasLimit}, receiver.From, gn)
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

	gnr, err := pBridge.GovernanceNonces(nil, TxTypeGovernanceRealtime)
	assert.NoError(t, err)
	_, err = pBridge.SetKLAYFee(&bind.TransactOpts{From: childAcc.From, Signer: childAcc.Signer, GasLimit: testGasLimit}, big.NewInt(KLAYFee), gnr)
	assert.NoError(t, err)
	_, err = pBridge.SetERC20Fee(&bind.TransactOpts{From: childAcc.From, Signer: childAcc.Signer, GasLimit: testGasLimit}, tokenAddr, big.NewInt(ERC20Fee), gnr+1)
	assert.NoError(t, err)
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

	// Register tokens on the bridge
	assert.NoError(t, err)
	pBridge.RegisterToken(&bind.TransactOpts{From: childAcc.From, Signer: childAcc.Signer, GasLimit: testGasLimit}, tokenAddr, tokenAddr, gn+1)
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

				switch ev.TokenType {
				case KLAY:
					assert.Equal(t, common.Address{}, ev.TokenAddress)
					assert.Equal(t, KLAYFee, ev.Fee.Int64())

					// HandleKLAYTransfer by Event
					tx, err := pBridge.HandleKLAYTransfer(&bind.TransactOpts{From: childAcc.From, Signer: childAcc.Signer, GasLimit: testGasLimit}, ev.From, ev.To, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber)
					if err != nil {
						log.Fatalf("Failed to HandleKLAYTransfer: %v", err)
					}
					fmt.Println("HandleKLAYTransfer Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block

				case ERC20:
					assert.Equal(t, tokenAddr, ev.TokenAddress)
					assert.Equal(t, ERC20Fee, ev.Fee.Int64())

					// HandleERC20Transfer by Event
					tx, err := pBridge.HandleERC20Transfer(&bind.TransactOpts{From: childAcc.From, Signer: childAcc.Signer, GasLimit: testGasLimit}, ev.From, ev.To, tokenAddr, ev.ValueOrTokenId, ev.RequestNonce, ev.Raw.BlockNumber)
					if err != nil {
						log.Fatalf("Failed to HandleERC20Transfer: %v", err)
					}
					fmt.Println("HandleERC20Transfer Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block
				}

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
		tx, err = token.Transfer(&bind.TransactOpts{From: parentAcc.From, Signer: parentAcc.Signer, GasLimit: testGasLimit}, Alice.From, big.NewInt(initialValue))
		if err != nil {
			log.Fatalf("Failed to Transfer for charging: %v", err)
		}
		fmt.Println("Transfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		balance, err = token.BalanceOf(nil, parentAcc.From)
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
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee))
		if err != nil {
			log.Fatalf("Failed to SafeTransferAndCall: %v", err)
		}
		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 7-2. Request ERC20 Transfer from Alice to Bob with insufficient zero feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(0))
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 7-3. Request ERC20 Transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee-1))
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 7-4. Request ERC20 Transfer from Alice to Bob with enough feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee+1))
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 8-1. Approve/Request ERC20 Transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee))
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 8-2. Approve/Request ERC20 Transfer from Alice to Bob with insufficient zero feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(0))
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 8-3. Approve/Request ERC20 Transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee-1))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee-1))
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 8-4. Approve/Request ERC20 Transfer from Alice to Bob with enough feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee+1))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee+1))
		assert.Equal(t, nil, err)

		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-1. Request KLAY transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY + KLAYFee), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY))
		if err != nil {
			log.Fatalf("Failed to RequestKLAYTransfer: %v", err)
		}
		fmt.Println("RequestKLAYTransfer Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-2. Request KLAY transfer from Alice to Bob with zero feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY))
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 9-3. Request KLAY transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY + (KLAYFee - 1)), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY))
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 9-4. Request KLAY transfer from Alice to Bob with enough feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKLAY + (KLAYFee + 1)), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKLAY))
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

	// 13. Nonce check on deploy error
	{

		nonce, err := sim.NonceAt(context.Background(), bam.scAccount.address, nil)
		if err != nil {
			t.Fatal("failed to sim.NonceAt", err)
		}
		bam.scAccount.SetNonce(nonce)

		addr2, err := bridgeManager.DeployBridgeNonceTest(sim)
		if err != nil {
			log.Fatalf("Failed to deploy new bridge contract: %v %v", err, addr2)
		}
	}

	bridgeManager.Stop()
}

// TestBasicJournal tests basic journal functionality.
func TestBasicJournal(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
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

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: big.NewInt(params.KLAY)}, auth2.From: {Balance: big.NewInt(params.KLAY)}, auth4.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	config := &SCConfig{}
	config.nodekey = key
	config.chainkey = key2
	config.DataDir = os.TempDir()
	config.VTRecovery = true

	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.MainChainAccountAddr = &chainKeyAddr

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		localBackend:         sim,
		remoteBackend:        sim,
		bridgeAccountManager: bam,
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
	remoteAddr, err := bm.DeployBridgeTest(sim, false)
	if err != nil {
		t.Fatal("deploy bridge test failed", remoteAddr)
	}

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
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
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

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: big.NewInt(params.KLAY)}, auth2.From: {Balance: big.NewInt(params.KLAY)}, auth4.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	config := &SCConfig{VTRecovery: true, VTRecoveryInterval: 60}
	config.nodekey = key
	config.chainkey = key2
	config.DataDir = os.TempDir()
	config.VTRecovery = true

	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.MainChainAccountAddr = &chainKeyAddr

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		localBackend:         sim,
		remoteBackend:        sim,
		bridgeAccountManager: bam,
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
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: os.TempDir(), VTRecovery: true},
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
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: os.TempDir(), VTRecovery: true},
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
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: os.TempDir(), VTRecovery: true},
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
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
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

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: big.NewInt(params.KLAY)}, auth2.From: {Balance: big.NewInt(params.KLAY)}, auth4.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	config := &SCConfig{}
	config.nodekey = key
	config.chainkey = key2
	config.DataDir = os.TempDir()
	config.VTRecovery = true

	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.MainChainAccountAddr = &chainKeyAddr

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		localBackend:         sim,
		remoteBackend:        sim,
		bridgeAccountManager: bam,
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
	err = bm.SetBridgeInfo(addr, bridgeInfo.bridge, common.Address{}, nil, sc.bridgeAccountManager.mcAccount, false, false)
	assert.NotEqual(t, nil, err)
	bm.Stop()
}

// TestScenarioSubUnsub tests subscription and unsubscription scenario.
func TestScenarioSubUnsub(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
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

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: big.NewInt(params.KLAY)}, auth2.From: {Balance: big.NewInt(params.KLAY)}, auth4.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	config := &SCConfig{}
	config.nodekey = key
	config.chainkey = key2
	config.DataDir = os.TempDir()
	config.VTRecovery = true

	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.MainChainAccountAddr = &chainKeyAddr

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		localBackend:         sim,
		remoteBackend:        sim,
		bridgeAccountManager: bam,
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
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: os.TempDir(), VTRecovery: true},
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
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
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

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: big.NewInt(params.KLAY)}, auth2.From: {Balance: big.NewInt(params.KLAY)}, auth4.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	config := &SCConfig{}
	config.nodekey = key
	config.chainkey = key2
	config.DataDir = os.TempDir()
	config.VTRecovery = true

	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.MainChainAccountAddr = &chainKeyAddr

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		localBackend:         sim,
		remoteBackend:        sim,
		bridgeAccountManager: bam,
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

	bm.bridges[addr], err = NewBridgeInfo(nil, addr, bridge, common.Address{}, nil, bam.scAccount, true, true)

	bm.journal.cache[addr] = &BridgeJournal{addr, addr, true}

	bm.SubscribeEvent(addr)
	err = bm.SubscribeEvent(addr)
	assert.NotEqual(t, nil, err)

	bm.Stop()
}

// for TestMethod
func (bm *BridgeManager) DeployBridgeTest(backend *backends.SimulatedBackend, local bool) (common.Address, error) {
	if local {
		acc := bm.subBridge.bridgeAccountManager.scAccount
		addr, bridge, err := bm.deployBridgeTest(acc, backend, true)
		err = bm.SetBridgeInfo(addr, bridge, common.Address{}, nil, acc, local, false)
		if err != nil {
			return common.Address{}, err
		}
		return addr, err
	} else {
		acc := bm.subBridge.bridgeAccountManager.mcAccount
		addr, bridge, err := bm.deployBridgeTest(acc, backend, false)
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

// Nonce should not be increased when error occurs
func (bm *BridgeManager) DeployBridgeNonceTest(backend bind.ContractBackend) (common.Address, error) {
	key := bm.subBridge.bridgeAccountManager.mcAccount.key
	nonce := bm.subBridge.bridgeAccountManager.mcAccount.GetNonce()
	bm.subBridge.bridgeAccountManager.mcAccount.key = nil
	_, addr, _ := bm.DeployBridge(backend, false)
	bm.subBridge.bridgeAccountManager.mcAccount.key = key

	if nonce != bm.subBridge.bridgeAccountManager.mcAccount.GetNonce() {
		return addr, errors.New("nonce is accidentally increased")
	}

	return addr, nil
}
