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
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/contracts/servicechain_nft"
	"github.com/klaytn/klaytn/contracts/servicechain_token"
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

type testInfo struct {
	t                 *testing.T
	sim               *backends.SimulatedBackend
	sc                *SubBridge
	bm                *BridgeManager
	localInfo         *BridgeInfo
	remoteInfo        *BridgeInfo
	tokenLocalAddr    common.Address
	tokenRemoteAddr   common.Address
	tokenLocalBridge  *sctoken.ServiceChainToken
	tokenRemoteBridge *sctoken.ServiceChainToken
	nftLocalAddr      common.Address
	nftRemoteAddr     common.Address
	nftLocalBridge    *scnft.ServiceChainNFT
	nftRemoteBridge   *scnft.ServiceChainNFT
	nodeAuth          *bind.TransactOpts
	chainAuth         *bind.TransactOpts
	aliceAuth         *bind.TransactOpts
	recoveryCh        chan bool
	mu                sync.Mutex
	nftIndex          int64
}

const (
	testGasLimit     = 1000000
	testAmount       = 321
	testToken        = 123
	testNFT          = int64(7321)
	testChargeToken  = 100000
	testTimeout      = 10 * time.Second
	testTxCount      = 7
	testBlockOffset  = 3 // +2 for genesis and bridge contract, +1 by a hardcoded hint
	testPendingCount = 3
	testNonceOffset  = testTxCount - testPendingCount
)

type operations struct {
	request     func(*testInfo, *BridgeInfo)
	handle      func(*testInfo, *BridgeInfo, *RequestValueTransferEvent)
	dummyHandle func(*testInfo, *BridgeInfo)
}

var (
	ops = map[uint8]*operations{
		KLAY: {
			request:     requestKLAYTransfer,
			handle:      handleKLAYTransfer,
			dummyHandle: dummyHandleRequestKLAYTransfer,
		},
		TOKEN: {
			request:     requestTokenTransfer,
			handle:      handleTokenTransfer,
			dummyHandle: dummyHandleRequestTokenTransfer,
		},
		NFT: {
			request:     requestNFTTransfer,
			handle:      handleNFTTransfer,
			dummyHandle: dummyHandleRequestNFTTransfer,
		},
	}
)

// TestBasicKLAYTransferRecovery tests each methods of the value transfer recovery.
func TestBasicKLAYTransferRecovery(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	// 1. Init dummy chain and do some value transfers.
	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)

	// 2. Update recovery hint.
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	t.Log("value transfer hint", vtr.service2mainHint)
	assert.Equal(t, uint64(testTxCount), vtr.service2mainHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.service2mainHint.handleNonce)

	// 3. Request events by using the hint.
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}

	// 4. Check pending events.
	t.Log("check pending tx", "len", len(vtr.serviceChainEvents))
	var count = 0
	for index, evt := range vtr.serviceChainEvents {
		assert.Equal(t, info.nodeAuth.From, evt.From)
		assert.Equal(t, info.aliceAuth.From, evt.To)
		assert.Equal(t, big.NewInt(testAmount), evt.Amount)
		assert.Equal(t, uint64(index+testNonceOffset), evt.RequestNonce)
		assert.Condition(t, func() bool {
			return uint64(testBlockOffset) <= evt.Raw.BlockNumber
		})
		count++
	}
	assert.Equal(t, testPendingCount, count)

	// 5. Recover pending events
	info.recoveryCh <- true
	assert.Equal(t, nil, vtr.recoverPendingEvents())
	ops[KLAY].dummyHandle(info, info.remoteInfo)

	// 6. Check empty pending events.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}
	assert.Equal(t, 0, len(vtr.serviceChainEvents))

	assert.Equal(t, nil, vtr.Recover()) // nothing to recover
}

// TestBasicTokenTransferRecovery tests the token transfer recovery.
func TestBasicTokenTransferRecovery(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[TOKEN].request(info, info.localInfo)
		}
	})

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)
	t.Log("token transfer hint", vtr.service2mainHint)

	info.recoveryCh <- true
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	ops[TOKEN].dummyHandle(info, info.remoteInfo)

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.Equal(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)
}

// TestBasicNFTTransferRecovery tests the NFT transfer recovery.
// TODO-Klaytn-ServiceChain: implement NFT transfer.
func TestBasicNFTTransferRecovery(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[NFT].request(info, info.localInfo)
		}
	})

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)
	t.Log("token transfer hint", vtr.service2mainHint)

	info.recoveryCh <- true
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	ops[NFT].dummyHandle(info, info.remoteInfo)

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.Equal(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)
}

// TestMethodRecover tests the valueTransferRecovery.Recover() method.
func TestMethodRecover(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)

	info.recoveryCh <- true
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	ops[KLAY].dummyHandle(info, info.remoteInfo)

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.Equal(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)
}

// TestMethodStop tests the Stop method for stop the internal goroutine.
func TestMethodStop(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 1}, info.localInfo, info.remoteInfo)
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)

	info.recoveryCh <- true
	err = vtr.Start()
	if err != nil {
		t.Fatal("fail to start the value transfer")
	}
	assert.Equal(t, nil, vtr.WaitRunningStatus(true, 5*time.Second))
	err = vtr.Stop()
	if err != nil {
		t.Fatal("fail to stop the value transfer")
	}
	assert.Equal(t, false, vtr.isRunning)
}

// TestFlagVTRecovery tests the disabled vtrecovery option.
func TestFlagVTRecovery(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete the journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 60}, info.localInfo, info.remoteInfo)
	vtr.Start()
	assert.Equal(t, nil, vtr.WaitRunningStatus(true, 5*time.Second))
	vtr.Stop()

	vtr = NewValueTransferRecovery(&SCConfig{VTRecovery: false}, info.localInfo, info.remoteInfo)
	err := vtr.Start()
	assert.Equal(t, ErrVtrDisabled, err)
	vtr.Stop()
}

// TestAlreadyStartedVTRecovery tests the already started VTR error cases.
func TestAlreadyStartedVTRecovery(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete the journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 60}, info.localInfo, info.remoteInfo)
	err := vtr.Start()
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, vtr.WaitRunningStatus(true, 5*time.Second))

	err = vtr.Start()
	assert.Equal(t, ErrVtrAlreadyStarted, err)

	vtr.Stop()
}

// TestScenarioMainChainRecovery tests the value transfer recovery of the main chain to service chain value transfers.
func TestScenarioMainChainRecovery(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.remoteInfo)
		}
	})

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.main2serviceHint.requestNonce, vtr.main2serviceHint.handleNonce)

	info.recoveryCh <- true
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	ops[KLAY].dummyHandle(info, info.localInfo)

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.Equal(t, vtr.main2serviceHint.requestNonce, vtr.main2serviceHint.handleNonce)
}

// TestScenarioAutomaticRecovery tests the recovery of the internal goroutine.
func TestScenarioAutomaticRecovery(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), BridgeAddrJournal)); err != nil {
			t.Fatalf("fail to delete journal file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 1}, info.localInfo, info.remoteInfo)
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)

	info.recoveryCh <- true
	err = vtr.Start()
	if err != nil {
		t.Fatal("fail to start the value transfer")
	}
	assert.Equal(t, nil, vtr.WaitRunningStatus(true, 5*time.Second))
	ops[KLAY].dummyHandle(info, info.remoteInfo)

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	vtr.Stop()
	assert.Equal(t, vtr.service2mainHint.requestNonce, vtr.service2mainHint.handleNonce)
}

// prepare generates dummy blocks for testing value transfer recovery.
func prepare(t *testing.T, vtcallback func(*testInfo)) *testInfo {
	// Setup configuration.
	config := &SCConfig{}
	config.nodekey, _ = crypto.GenerateKey()
	config.chainkey, _ = crypto.GenerateKey()
	config.DataDir = os.TempDir()
	config.VTRecovery = true

	// Generate a new random account and a funded simulator.
	nodeAuth := bind.NewKeyedTransactor(config.nodekey)
	chainAuth := bind.NewKeyedTransactor(config.chainkey)
	aliceKey, _ := crypto.GenerateKey()
	aliceAuth := bind.NewKeyedTransactor(aliceKey)
	chainKeyAddr := crypto.PubkeyToAddress(config.chainkey.PublicKey)
	config.ServiceChainAccountAddr = &chainKeyAddr

	// Alloc genesis and create a simulator.
	alloc := blockchain.GenesisAlloc{
		nodeAuth.From:  {Balance: big.NewInt(params.KLAY)},
		chainAuth.From: {Balance: big.NewInt(params.KLAY)},
		aliceAuth.From: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	bam, _ := NewBridgeAccountManager(config.chainkey, config.nodekey)

	sc := &SubBridge{
		config:               config,
		peers:                newBridgePeerSet(),
		bridgeAccountManager: bam,
	}
	handler, err := NewSubBridgeHandler(sc.config, sc)
	if err != nil {
		log.Fatalf("Failed to initialize the bridgeHandler : %v", err)
		return nil
	}
	sc.handler = handler

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
	localInfo, _ := bm.GetBridgeInfo(localAddr)
	remoteInfo, _ := bm.GetBridgeInfo(remoteAddr)
	sim.Commit()

	// Prepare token contract
	tokenLocalAddr, _, tokenLocal, err := sctoken.DeployServiceChainToken(nodeAuth, sim, localAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainToken: %v", err)
	}
	tokenRemoteAddr, _, tokenRemote, err := sctoken.DeployServiceChainToken(chainAuth, sim, remoteAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainToken: %v", err)
	}
	sim.Commit()
	testToken := big.NewInt(testChargeToken)
	opts := &bind.TransactOpts{From: chainAuth.From, Signer: chainAuth.Signer, GasLimit: testGasLimit}
	_, err = tokenRemote.Transfer(opts, remoteAddr, testToken)
	if err != nil {
		log.Fatalf("Failed to Transfer for charging: %v", err)
	}
	sim.Commit()

	// Prepare NFT contract
	nftLocalAddr, _, nftLocal, err := scnft.DeployServiceChainNFT(nodeAuth, sim, localAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}
	nftRemoteAddr, _, nftRemote, err := scnft.DeployServiceChainNFT(chainAuth, sim, remoteAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}
	sim.Commit()

	// Register tokens on the bridge
	nodeOpts := &bind.TransactOpts{From: nodeAuth.From, Signer: nodeAuth.Signer, GasLimit: testGasLimit}
	chainOpts := &bind.TransactOpts{From: chainAuth.From, Signer: chainAuth.Signer, GasLimit: testGasLimit}
	_, err = localInfo.bridge.RegisterToken(nodeOpts, tokenLocalAddr, tokenRemoteAddr)
	_, err = localInfo.bridge.RegisterToken(nodeOpts, nftLocalAddr, nftRemoteAddr)
	_, err = remoteInfo.bridge.RegisterToken(chainOpts, tokenRemoteAddr, tokenLocalAddr)
	_, err = remoteInfo.bridge.RegisterToken(chainOpts, nftRemoteAddr, nftLocalAddr)
	sim.Commit()

	// Register an NFT to chain account (minting)
	for i := 0; i < testTxCount; i++ {
		opts := &bind.TransactOpts{From: nodeAuth.From, Signer: nodeAuth.Signer, GasLimit: testGasLimit}
		_, err = nftLocal.Register(opts, nodeAuth.From, big.NewInt(testNFT+int64(i)))
		if err != nil {
			log.Fatalf("Failed to Register NFT: %v", err)
		}
		opts = &bind.TransactOpts{From: chainAuth.From, Signer: chainAuth.Signer, GasLimit: testGasLimit}
		_, err = nftRemote.Register(opts, remoteAddr, big.NewInt(testNFT+int64(i)))
		if err != nil {
			log.Fatalf("Failed to Register NFT: %v", err)
		}
	}
	sim.Commit()

	// Subscribe events.
	err = bm.SubscribeEvent(localAddr)
	if err != nil {
		t.Fatalf("local bridge manager event subscription failed")
	}
	err = bm.SubscribeEvent(remoteAddr)
	if err != nil {
		t.Fatalf("remote bridge manager event subscription failed")
	}

	// Prepare channel for event handling.
	requestVTCh := make(chan RequestValueTransferEvent)
	handleVTCh := make(chan HandleValueTransferEvent)
	recoveryCh := make(chan bool)
	bm.SubscribeTokenReceived(requestVTCh)
	bm.SubscribeTokenWithDraw(handleVTCh)

	info := testInfo{
		t, sim, sc, bm, localInfo, remoteInfo,
		tokenLocalAddr, tokenRemoteAddr, tokenLocal, tokenRemote,
		nftLocalAddr, nftRemoteAddr, nftLocal, nftRemote,
		nodeAuth, chainAuth, aliceAuth, recoveryCh, sync.Mutex{}, testNFT,
	}

	// Start a event handling loop.
	wg := sync.WaitGroup{}
	wg.Add((2 * testTxCount) - testPendingCount)
	var isRecovery = false
	go func() {
		for {
			select {
			case ev := <-recoveryCh:
				isRecovery = ev
			case ev := <-requestVTCh:
				// Intentionally lost a single handle value transfer.
				// Since the increase of monotony in nonce is checked in the contract,
				// all subsequent handle transfer will be failed.
				if ev.RequestNonce == (testTxCount - testPendingCount) {
					t.Log("missing handle value transfer", "nonce", ev.RequestNonce)
				} else {
					switch ev.TokenType {
					case KLAY, TOKEN, NFT:
						break
					default:
						t.Fatal("received TokenType is unknown")
					}

					if ev.ContractAddr == info.localInfo.address {
						ops[ev.TokenType].handle(&info, info.remoteInfo, &ev)
					} else {
						ops[ev.TokenType].handle(&info, info.localInfo, &ev)
					}
				}

				if !isRecovery {
					wg.Done()
				}
			case _ = <-handleVTCh:
				if !isRecovery {
					wg.Done()
				}
			}
		}
	}()

	// Request value transfer.
	vtcallback(&info)
	WaitGroupWithTimeOut(&wg, testTimeout, t)

	return &info
}

func requestKLAYTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	opts := bi.account.GetTransactOpts()
	opts.Value = big.NewInt(testAmount)
	_, err := bi.bridge.RequestKLAYTransfer(opts, info.aliceAuth.From)
	if err != nil {
		log.Fatalf("Failed to RequestKLAYTransfer: %v", err)
	}
	info.sim.Commit()
}

func handleKLAYTransfer(info *testInfo, bi *BridgeInfo, ev *RequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	assert.Equal(info.t, new(big.Int).SetUint64(testAmount), ev.Amount)
	opts := bi.account.GetTransactOpts()
	_, err := bi.bridge.HandleKLAYTransfer(opts, ev.Amount, ev.To, ev.RequestNonce, ev.BlockNumber)
	if err != nil {
		log.Fatalf("\tFailed to HandleKLAYTransfer: %v", err)
	}
	info.sim.Commit()
}

// TODO-Klaytn-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestKLAYTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents(math.MaxUint64) {
		handleKLAYTransfer(info, bi, ev)
	}
	info.sim.Commit()
}

func requestTokenTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	testToken := big.NewInt(testToken)
	opts := bi.account.GetTransactOpts()

	var err error
	if bi.onServiceChain {
		_, err = info.tokenLocalBridge.RequestValueTransfer(opts, testToken, info.chainAuth.From)
	} else {
		_, err = info.tokenRemoteBridge.RequestValueTransfer(opts, testToken, info.nodeAuth.From)
	}

	if err != nil {
		log.Fatalf("Failed to RequestValueTransfer for charging: %v", err)
	}
	info.sim.Commit()
}

func handleTokenTransfer(info *testInfo, bi *BridgeInfo, ev *RequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	assert.Equal(info.t, new(big.Int).SetUint64(testToken), ev.Amount)
	_, err := bi.bridge.HandleTokenTransfer(
		bi.account.GetTransactOpts(),
		ev.Amount, ev.To, info.tokenRemoteAddr, ev.RequestNonce, ev.BlockNumber)
	if err != nil {
		log.Fatalf("Failed to HandleTokenTransfer: %v", err)
	}
	info.sim.Commit()
}

// TODO-Klaytn-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestTokenTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents(math.MaxUint64) {
		handleTokenTransfer(info, bi, ev)
	}
	info.sim.Commit()
}

func requestNFTTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	opts := bi.account.GetTransactOpts()
	// TODO-Klaytn need to separate service / main chain nftIndex.
	nftIndex := new(big.Int).SetInt64(info.nftIndex)

	var err error
	if bi.onServiceChain {
		_, err = info.nftLocalBridge.RequestValueTransfer(opts, nftIndex, info.aliceAuth.From)
	} else {
		_, err = info.nftRemoteBridge.RequestValueTransfer(opts, nftIndex, info.aliceAuth.From)
	}

	if err != nil {
		log.Fatalf("Failed to requestNFTTransfer for charging: %v", err)
	}
	info.nftIndex++
	info.sim.Commit()
}

func handleNFTTransfer(info *testInfo, bi *BridgeInfo, ev *RequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	var nftAddr common.Address

	if bi.onServiceChain {
		nftAddr = info.nftLocalAddr
	} else {
		nftAddr = info.nftRemoteAddr
	}

	_, err := bi.bridge.HandleNFTTransfer(
		bi.account.GetTransactOpts(),
		ev.Amount, ev.To, nftAddr, ev.RequestNonce, ev.BlockNumber)
	if err != nil {
		log.Fatalf("Failed to handleNFTTransfer: %v", err)
	}
	info.sim.Commit()
}

// TODO-Klaytn-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestNFTTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents(math.MaxUint64) {
		handleNFTTransfer(info, bi, ev)
	}
	info.sim.Commit()
}
