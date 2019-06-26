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
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/contracts/servicechain_nft"
	"github.com/klaytn/klaytn/contracts/servicechain_token"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/node/sc/bridgepool"
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
	bridge.RegisterToken(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, tokenAddr, tokenAddr)
	bridge.RegisterToken(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, nftAddr, nftAddr)
	sim.Commit() // block

	cTokenAddr, err := bridge.AllowedTokens(nil, tokenAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cTokenAddr, tokenAddr)
	cNftAddr, err := bridge.AllowedTokens(nil, nftAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cNftAddr, nftAddr)

	// TODO-Klaytn-Servicechain needs to support WaitDeployed
	//timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancelTimeout()
	//
	//addr, err = bind.WaitDeployed(timeoutContext, sim, tx)
	//if err != nil {
	//	log.Fatal("Failed to DeployGXToken.", "err", err, "txHash", tx.Hash().String())
	//
	//}
	//fmt.Println("GXToken is deployed.", "addr", addr.String(), "txHash", tx.Hash().String())

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

	tokenCh := make(chan RequestValueTransferEvent)
	tokenSendCh := make(chan HandleValueTransferEvent)
	bridgeManager.SubscribeTokenReceived(tokenCh)
	bridgeManager.SubscribeTokenWithDraw(tokenSendCh)

	go func() {
		for {
			select {
			case ev := <-tokenCh:
				fmt.Println("Deposit Event",
					"type", ev.TokenType,
					"amount", ev.Amount,
					"from", ev.From.String(),
					"to", ev.To.String(),
					"contract", ev.ContractAddr.String(),
					"token", ev.TokenAddr.String(),
					"requestNonce", ev.RequestNonce)

				switch ev.TokenType {
				case 0:
					// WithdrawKLAY by Event
					tx, err := bridge.HandleKLAYTransfer(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, ev.Amount, ev.To, ev.RequestNonce, ev.BlockNumber)
					if err != nil {
						log.Fatalf("Failed to WithdrawKLAY: %v", err)
					}
					fmt.Println("WithdrawKLAY Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block

				case 1:
					// WithdrawToken by Event
					tx, err := bridge.HandleTokenTransfer(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, ev.Amount, ev.To, tokenAddr, ev.RequestNonce, ev.BlockNumber)
					if err != nil {
						log.Fatalf("Failed to WithdrawToken: %v", err)
					}
					fmt.Println("WithdrawToken Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block

				case 2:
					owner, err := nft.OwnerOf(&bind.CallOpts{From: auth.From}, big.NewInt(int64(nftTokenID)))
					if err != nil {
						t.Fatal(err)
					}
					fmt.Println("NFT owner before WithdrawERC721: ", owner.String())

					// WithdrawToken by Event
					tx, err := bridge.HandleNFTTransfer(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, ev.Amount, ev.To, nftAddr, ev.RequestNonce, ev.BlockNumber)
					if err != nil {
						log.Fatalf("Failed to WithdrawERC721: %v", err)
					}
					fmt.Println("WithdrawERC721 Transaction by event ", tx.Hash().Hex())
					sim.Commit() // block
				}

				wg.Done()

			case ev := <-tokenSendCh:
				fmt.Println("receive token withdraw event ", ev.ContractAddr.Hex())
				fmt.Println("Withdraw Event",
					"type", ev.TokenType,
					"amount", ev.Amount,
					"owner", ev.Owner.String(),
					"contract", ev.ContractAddr.String(),
					"token", ev.TokenAddr.String(),
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

		// TODO-Klaytn-Servicechain needs to support WaitMined
		//timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancelTimeout()
		//
		//receipt, err := bind.WaitMined(timeoutContext, sim, tx)
		//if err != nil {
		//	log.Fatal("Failed to WithdrawToken.", "err", err, "txHash", tx.Hash().String(),"status",receipt.Status)
		//}
		//fmt.Println("WithdrawToken is executed.", "addr", addr.String(), "txHash", tx.Hash().String())

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth.From}, auth.From)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("auth token balance", balance.String())

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth2.From}, auth2.From)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("auth2 token balance", balance.String())

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth3.From}, auth3.From)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("auth3 token balance", balance.String())

		balance, err = token.BalanceOf(&bind.CallOpts{From: auth4.From}, auth4.From)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("auth4 token balance", balance.String())
	}

	// 6. Register (Mint) an NFT to Auth4
	{
		tx, err = nft.Register(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: testGasLimit}, auth4.From, big.NewInt(int64(nftTokenID)))
		if err != nil {
			log.Fatalf("Failed to Register NFT: %v", err)
		}
		fmt.Println("Register NFT Transaction", tx.Hash().Hex())
		sim.Commit() // block

		balance, err = nft.BalanceOf(&bind.CallOpts{From: auth.From}, auth4.From)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("auth4 NFT balance", balance.String())
		fmt.Println("auth4 address", auth4.From.String())
		owner, err := nft.OwnerOf(&bind.CallOpts{From: auth.From}, big.NewInt(int64(nftTokenID)))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("NFT owner after registering", owner.String())
	}

	// 7. RequestValueTransfer from auth2 to auth3
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, GasLimit: testGasLimit}, testToken, auth3.From)
		if err != nil {
			log.Fatalf("Failed to SafeTransferAndCall: %v", err)
		}
		fmt.Println("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		// TODO-Klaytn-Servicechain needs to support WaitMined
		//timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancelTimeout()
		//
		//receipt, err := bind.WaitMined(timeoutContext, sim, tx)
		//if err != nil {
		//	log.Fatal("Failed to RequestValueTransfer.", "err", err, "txHash", tx.Hash().String(), "status", receipt.Status)
		//
		//}
		//fmt.Println("RequestValueTransfer is executed.", "addr", addr.String(), "txHash", tx.Hash().String())

	}

	// 8. DepositKLAY from auth to auth3
	{
		tx, err = bridge.RequestKLAYTransfer(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, Value: testKLAY, GasLimit: testGasLimit}, auth3.From)
		if err != nil {
			log.Fatalf("Failed to DepositKLAY: %v", err)
		}
		fmt.Println("DepositKLAY Transaction", tx.Hash().Hex())

		sim.Commit() // block
	}

	// 9. Request NFT value transfer from auth4 to auth3
	{
		tx, err = nft.RequestValueTransfer(&bind.TransactOpts{From: auth4.From, Signer: auth4.Signer, GasLimit: testGasLimit}, big.NewInt(int64(nftTokenID)), auth3.From)
		if err != nil {
			log.Fatalf("Failed to nft.RequestValueTransfer: %v", err)
		}
		fmt.Println("nft.RequestValueTransfer Transaction", tx.Hash().Hex())

		sim.Commit() // block

		timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()

		receipt, err := bind.WaitMined(timeoutContext, sim, tx)
		if err != nil {
			log.Fatal("Failed to nft.RequestValueTransfer.", "err", err, "txHash", tx.Hash().String(), "status", receipt.Status)

		}
		fmt.Println("nft.RequestValueTransfer is executed.", "addr", addr.String(), "txHash", tx.Hash().String())
	}

	// Wait a few second for wait group
	WaitGroupWithTimeOut(&wg, 3*time.Second, t)

	// 10. Check Token balance
	{
		balance, err = token.BalanceOf(&bind.CallOpts{From: auth.From}, auth.From)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("auth token balance", balance.String())
		balance, err = token.BalanceOf(&bind.CallOpts{From: auth2.From}, auth2.From)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("auth2 token balance", balance.String())
		balance, err = token.BalanceOf(&bind.CallOpts{From: auth3.From}, auth3.From)
		if err != nil {
			t.Fatal(err)
		}
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
		if err != nil {
			t.Fatal(err)
		}
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
	if sc.addressManager, err = NewAddressManager(); err != nil {
		t.Fatal("new address manager is failed")
	}

	localAddr, err := bm.DeployBridgeTest(sim, true)
	if err != nil {
		t.Fatal("deploy bridge test failed", localAddr)
	}
	remoteAddr, err := bm.DeployBridgeTest(sim, false)
	if err != nil {
		t.Fatal("deploy bridge test failed", remoteAddr)
	}

	sc.addressManager.AddBridge(localAddr, remoteAddr)
	bm.SetJournal(localAddr, remoteAddr)

	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	localInfo, ok := bm.GetBridgeInfo(localAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, localInfo.subscribed)
	assert.Equal(t, sc.addressManager.GetCounterPartBridge(localAddr), remoteAddr)

	remoteInfo, ok := bm.GetBridgeInfo(remoteAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, remoteInfo.subscribed)
	assert.Equal(t, sc.addressManager.GetCounterPartBridge(remoteAddr), localAddr)
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
	bm, err := NewBridgeManager(sc)
	if sc.addressManager, err = NewAddressManager(); err != nil {
		t.Fatal("new address manager is failed")
	}

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
	sc.addressManager.AddBridge(bridgeAddrs[0], bridgeAddrs[1])
	bm.SetJournal(bridgeAddrs[0], bridgeAddrs[1])
	bm.journal.cache[bridgeAddrs[0]].Subscribed = true
	sc.addressManager.AddBridge(bridgeAddrs[2], bridgeAddrs[3])
	bm.SetJournal(bridgeAddrs[2], bridgeAddrs[3])
	bm.journal.cache[bridgeAddrs[2]].Subscribed = true

	// Call RestoreBridges
	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	// Case 1: check bridge contract creation.
	for i := 0; i < 4; i++ {
		info, _ := bm.GetBridgeInfo(bridgeAddrs[i])
		assert.NotEqual(t, nil, info.bridge)
	}

	// Case 2: check address manager
	am := bm.subBridge.addressManager
	assert.Equal(t, bridgeAddrs[1], am.GetCounterPartBridge(bridgeAddrs[0]))
	assert.Equal(t, bridgeAddrs[3], am.GetCounterPartBridge(bridgeAddrs[2]))

	// Case 3: check subscription
	for i := 0; i < 4; i++ {
		info, _ := bm.GetBridgeInfo(bridgeAddrs[i])
		assert.Equal(t, true, info.subscribed)
	}

	// Case 4: check recovery
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
	if sc.addressManager, err = NewAddressManager(); err != nil {
		t.Fatal("new address manager is failed")
	}

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

	bm.bridges[addr] = &BridgeInfo{
		nil,
		addr,
		common.Address{},
		nil,
		bridge,
		nil,
		true,
		true,
		bridgepool.NewEventSortedMap(),
		0,
		true,
		0,
		0,
		0,
		make(chan struct{}),
		make(chan struct{}),
	}
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
		addr, bridge, err := bm.deployBridgeTest(acc, backend)
		err = bm.SetBridgeInfo(addr, bridge, common.Address{}, nil, acc, local, false)
		if err != nil {
			return common.Address{}, err
		}
		return addr, err
	} else {
		acc := bm.subBridge.bridgeAccountManager.mcAccount
		addr, bridge, err := bm.deployBridgeTest(acc, backend)
		err = bm.SetBridgeInfo(addr, bridge, common.Address{}, nil, acc, local, false)
		if err != nil {
			return common.Address{}, err
		}
		return addr, err
	}
}

func (bm *BridgeManager) deployBridgeTest(acc *accountInfo, backend *backends.SimulatedBackend) (common.Address, *bridge.Bridge, error) {
	auth := acc.GetTransactOpts()
	auth.Value = big.NewInt(10000)
	addr, tx, contract, err := bridge.DeployBridge(auth, backend)
	if err != nil {
		logger.Error("", "err", err)
		return common.Address{}, nil, err
	}
	logger.Info("Bridge is deploying on CurrentChain", "addr", addr, "txHash", tx.Hash().String())

	// TODO-Klaytn-Servicechain needs to support WaitMined
	//backend.Commit()
	//
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
