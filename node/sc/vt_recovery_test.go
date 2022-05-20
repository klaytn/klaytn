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
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	sctoken "github.com/klaytn/klaytn/contracts/sc_erc20"
	scnft "github.com/klaytn/klaytn/contracts/sc_erc721"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
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
)

type operations struct {
	request     func(*testInfo, *BridgeInfo)
	handle      func(*testInfo, *BridgeInfo, IRequestValueTransferEvent)
	dummyHandle func(*testInfo, *BridgeInfo)
}

var (
	ops = map[uint8]*operations{
		KLAY: {
			request:     requestKLAYTransfer,
			handle:      handleKLAYTransfer,
			dummyHandle: dummyHandleRequestKLAYTransfer,
		},
		ERC20: {
			request:     requestTokenTransfer,
			handle:      handleTokenTransfer,
			dummyHandle: dummyHandleRequestTokenTransfer,
		},
		ERC721: {
			request:     requestNFTTransfer,
			handle:      handleNFTTransfer,
			dummyHandle: dummyHandleRequestNFTTransfer,
		},
	}
)

// TestBasicKLAYTransferRecovery tests each methods of the value transfer recovery.
func TestBasicKLAYTransferRecovery(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// 1. Init dummy chain and do some value transfers.
	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)

	// 2. Update recovery hint.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	t.Log("value transfer hint", vtr.child2parentHint)
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.child2parentHint.handleNonce)

	// 3. Request events by using the hint.
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}

	// 4. Check pending events.
	t.Log("check pending tx", "len", len(vtr.childEvents))
	var count = 0
	for _, ev := range vtr.childEvents {
		assert.Equal(t, info.nodeAuth.From, ev.GetFrom())
		assert.Equal(t, info.aliceAuth.From, ev.GetTo())
		assert.Equal(t, big.NewInt(testAmount), ev.GetValueOrTokenId())
		assert.Condition(t, func() bool {
			return uint64(testBlockOffset) <= ev.GetRaw().BlockNumber
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
	assert.Equal(t, 0, len(vtr.childEvents))

	assert.Equal(t, nil, vtr.Recover()) // nothing to recover
}

// TestKLAYTransferLongRangeRecovery tests a long block range recovery.
func TestKLAYTransferLongRangeRecovery(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	oldMaxPendingTxs := maxPendingTxs
	maxPendingTxs = 2
	defer func() {
		maxPendingTxs = oldMaxPendingTxs
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// 1. Init dummy chain and do some value transfers.
	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
			for i := uint64(0); i < filterLogsStride; i++ {
				info.sim.Commit()
			}
		}
	})
	defer info.sim.Close()
	// TODO-Klaytn need to remove sleep
	time.Sleep(1 * time.Second)
	info.sim.Commit()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)

	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.child2parentHint.handleNonce)

	// 2. first recovery.
	info.recoveryCh <- true
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	// TODO-Klaytn need to remove sleep
	time.Sleep(1 * time.Second)
	info.sim.Commit()

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount+maxPendingTxs), vtr.child2parentHint.handleNonce)

	// 3. second recovery.
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	// TODO-Klaytn need to remove sleep
	time.Sleep(1 * time.Second)
	info.sim.Commit()

	// 4. Check if recovery is done.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.handleNonce)
}

// TestBasicTokenTransferRecovery tests the token transfer recovery.
func TestBasicTokenTransferRecovery(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[ERC20].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
	t.Log("token transfer hint", vtr.child2parentHint)

	info.recoveryCh <- true
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	ops[ERC20].dummyHandle(info, info.remoteInfo)

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.Equal(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
}

// TestBasicNFTTransferRecovery tests the NFT transfer recovery.
// TODO-Klaytn-ServiceChain: implement NFT transfer.
func TestBasicNFTTransferRecovery(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[ERC721].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
	t.Log("token transfer hint", vtr.child2parentHint)

	info.recoveryCh <- true
	err = vtr.Recover()
	if err != nil {
		t.Fatal("fail to recover the value transfer")
	}
	ops[ERC721].dummyHandle(info, info.remoteInfo)

	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.Equal(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
}

// TestMethodRecover tests the valueTransferRecovery.Recover() method.
func TestMethodRecover(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)

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
	assert.Equal(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
}

// TestMethodStop tests the Stop method for stop the internal goroutine.
func TestMethodStop(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 1}, info.localInfo, info.remoteInfo)
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)

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
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 60}, info.localInfo, info.remoteInfo)
	vtr.Start()
	assert.Equal(t, nil, vtr.WaitRunningStatus(true, 5*time.Second))
	vtr.Stop()

	vtr = NewValueTransferRecovery(&SCConfig{VTRecovery: false}, info.localInfo, info.remoteInfo)
	err = vtr.Start()
	assert.Equal(t, ErrVtrDisabled, err)
	vtr.Stop()
}

// TestAlreadyStartedVTRecovery tests the already started VTR error cases.
func TestAlreadyStartedVTRecovery(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()
	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 60}, info.localInfo, info.remoteInfo)
	err = vtr.Start()
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, vtr.WaitRunningStatus(true, 5*time.Second))

	err = vtr.Start()
	assert.Equal(t, ErrVtrAlreadyStarted, err)

	vtr.Stop()
}

// TestScenarioMainChainRecovery tests the value transfer recovery of the parent chain to child chain value transfers.
func TestScenarioMainChainRecovery(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.remoteInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.parent2childHint.requestNonce, vtr.parent2childHint.handleNonce)

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
	assert.Equal(t, vtr.parent2childHint.requestNonce, vtr.parent2childHint.handleNonce)
}

// TestScenarioAutomaticRecovery tests the recovery of the internal goroutine.
func TestScenarioAutomaticRecovery(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 1}, info.localInfo, info.remoteInfo)
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update a value transfer hint")
	}
	assert.NotEqual(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)

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
	assert.Equal(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
}

// TestMultiOperatorRequestRecovery tests value transfer recovery for the multi-operator.
func TestMultiOperatorRequestRecovery(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// 1. Init dummy chain and do some value transfers.
	info := prepare(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KLAY].request(info, info.localInfo)
		}
	})
	defer info.sim.Close()

	// 2. Set multi-operator.
	cAcc := info.nodeAuth
	pAcc := info.chainAuth
	opts := &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
	_, err = info.localInfo.bridge.RegisterOperator(opts, pAcc.From)
	assert.NoError(t, err)
	opts = &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
	_, err = info.remoteInfo.bridge.RegisterOperator(opts, cAcc.From)
	assert.NoError(t, err)
	info.sim.Commit()

	// 3. Set operator threshold.
	opts = &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
	_, err = info.localInfo.bridge.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	opts = &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
	_, err = info.remoteInfo.bridge.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)

	// 4. Update recovery hint.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	t.Log("value transfer hint", vtr.child2parentHint)
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.child2parentHint.handleNonce)

	// 5. Request events by using the hint.
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}

	// 6. Check pending events.
	t.Log("check pending tx", "len", len(vtr.childEvents))
	var count = 0
	for _, ev := range vtr.childEvents {
		assert.Equal(t, info.nodeAuth.From, ev.GetFrom())
		assert.Equal(t, info.aliceAuth.From, ev.GetTo())
		assert.Equal(t, big.NewInt(testAmount), ev.GetValueOrTokenId())
		assert.Condition(t, func() bool {
			return uint64(testBlockOffset) <= ev.GetRaw().BlockNumber
		})
		count++
	}
	assert.Equal(t, testPendingCount, count)

	// 7. Recover pending events
	info.recoveryCh <- true
	assert.Equal(t, nil, vtr.recoverPendingEvents())
	ops[KLAY].dummyHandle(info, info.remoteInfo)

	// 8. Recover from the other operator (value transfer is not recovered yet).
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	t.Log("value transfer hint", vtr.child2parentHint)
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.child2parentHint.handleNonce)

	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}
	assert.Equal(t, testPendingCount, len(vtr.childEvents))
	assert.Equal(t, nil, vtr.recoverPendingEvents())
	info.remoteInfo.account = info.localInfo.account // other operator
	ops[KLAY].dummyHandle(info, info.remoteInfo)

	// 9. Check results.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}
	assert.Equal(t, 0, len(vtr.childEvents))
	assert.Equal(t, nil, vtr.Recover()) // nothing to recover
}

// prepare generates dummy blocks for testing value transfer recovery.
func prepare(t *testing.T, vtcallback func(*testInfo)) *testInfo {
	// Setup configuration.
	config := &SCConfig{}
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, err := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	assert.NoError(t, err)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	cAcc := bacc.cAccount.GenerateTransactOpts()
	pAcc := bacc.pAccount.GenerateTransactOpts()

	// Generate a new random account and a funded simulator.
	aliceKey, _ := crypto.GenerateKey()
	aliceAuth := bind.NewKeyedTransactor(aliceKey)

	// Alloc genesis and create a simulator.
	alloc := blockchain.GenesisAlloc{
		cAcc.From:      {Balance: big.NewInt(params.KLAY)},
		pAcc.From:      {Balance: big.NewInt(params.KLAY)},
		aliceAuth.From: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
		localBackend:   sim,
		remoteBackend:  sim,
	}
	handler, err := NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize the bridgeHandler : %v", err)
		return nil
	}
	sc.handler = handler
	sc.blockchain = sim.BlockChain()

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)
	localAddr, err := bm.DeployBridgeTest(sim, 10000, true)
	assert.NoError(t, err)
	remoteAddr, err := bm.DeployBridgeTest(sim, 10000, false)
	assert.NoError(t, err)

	localInfo, _ := bm.GetBridgeInfo(localAddr)
	remoteInfo, _ := bm.GetBridgeInfo(remoteAddr)
	sim.Commit()

	// Prepare token contract
	tokenLocalAddr, tx, tokenLocal, err := sctoken.DeployServiceChainToken(cAcc, sim, localAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainToken: %v", err)
	}
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	tokenRemoteAddr, tx, tokenRemote, err := sctoken.DeployServiceChainToken(pAcc, sim, remoteAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainToken: %v", err)
	}
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	testToken := big.NewInt(testChargeToken)
	opts := &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
	tx, err = tokenRemote.Transfer(opts, remoteAddr, testToken)
	if err != nil {
		log.Fatalf("Failed to Transfer for charging: %v", err)
	}
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	// Prepare NFT contract
	nftLocalAddr, tx, nftLocal, err := scnft.DeployServiceChainNFT(cAcc, sim, localAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	nftRemoteAddr, tx, nftRemote, err := scnft.DeployServiceChainNFT(pAcc, sim, remoteAddr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	// Register tokens on the bridge
	nodeOpts := &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
	chainOpts := &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
	tx, err = localInfo.bridge.RegisterToken(nodeOpts, tokenLocalAddr, tokenRemoteAddr)
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	tx, err = localInfo.bridge.RegisterToken(nodeOpts, nftLocalAddr, nftRemoteAddr)
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	tx, err = remoteInfo.bridge.RegisterToken(chainOpts, tokenRemoteAddr, tokenLocalAddr)
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	tx, err = remoteInfo.bridge.RegisterToken(chainOpts, nftRemoteAddr, nftLocalAddr)
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))

	// Register an NFT to chain account (minting)
	for i := 0; i < testTxCount; i++ {
		opts := &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
		tx, err = nftLocal.MintWithTokenURI(opts, cAcc.From, big.NewInt(testNFT+int64(i)), "testURI")
		assert.NoError(t, err)
		sim.Commit()
		assert.Nil(t, bind.CheckWaitMined(sim, tx))

		opts = &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
		tx, err = nftRemote.MintWithTokenURI(opts, remoteAddr, big.NewInt(testNFT+int64(i)), "testURI")
		assert.NoError(t, err)
		sim.Commit()
		assert.Nil(t, bind.CheckWaitMined(sim, tx))
	}

	// Register the owner as a signer
	_, err = localInfo.bridge.RegisterOperator(&bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}, cAcc.From)
	assert.NoError(t, err)
	_, err = remoteInfo.bridge.RegisterOperator(&bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}, pAcc.From)
	assert.NoError(t, err)
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
	requestVTencodedCh := make(chan RequestValueTransferEncodedEvent)
	handleVTCh := make(chan *HandleValueTransferEvent)
	recoveryCh := make(chan bool)
	bm.SubscribeReqVTev(requestVTCh)
	bm.SubscribeReqVTencodedEv(requestVTencodedCh)
	bm.SubscribeHandleVTev(handleVTCh)

	info := testInfo{
		t, sim, sc, bm, localInfo, remoteInfo,
		tokenLocalAddr, tokenRemoteAddr, tokenLocal, tokenRemote,
		nftLocalAddr, nftRemoteAddr, nftLocal, nftRemote,
		cAcc, pAcc, aliceAuth, recoveryCh, sync.Mutex{}, testNFT,
	}

	// Start a event handling loop.
	wg := sync.WaitGroup{}
	wg.Add((2 * testTxCount) - testPendingCount)
	var isRecovery = false
	reqHandler := func(ev IRequestValueTransferEvent) {
		t.Log("request value transfer", "nonce", ev.GetRequestNonce())
		if ev.GetRequestNonce() >= (testTxCount - testPendingCount) {
			t.Log("missing handle value transfer", "nonce", ev.GetRequestNonce())
		} else {
			switch ev.GetTokenType() {
			case KLAY, ERC20, ERC721:
				break
			default:
				t.Errorf("received ev.TokenType is unknown: %v", ev.GetTokenType())
				return
			}

			if ev.GetRaw().Address == info.localInfo.address {
				ops[ev.GetTokenType()].handle(&info, info.remoteInfo, ev)
			} else {
				ops[ev.GetTokenType()].handle(&info, info.localInfo, ev)
			}
		}
	}
	go func() {
		for {
			select {
			case ev := <-recoveryCh:
				isRecovery = ev
			case ev := <-requestVTCh:
				reqHandler(ev)
				if !isRecovery {
					wg.Done()
				}
			case ev := <-requestVTencodedCh:
				reqHandler(ev)
				if !isRecovery {
					wg.Done()
				}
			case ev := <-handleVTCh:
				t.Log("handle value transfer", "nonce", ev.HandleNonce)
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

	opts := bi.account.GenerateTransactOpts()
	opts.Value = big.NewInt(testAmount)
	tx, err := bi.bridge.RequestKLAYTransfer(opts, info.aliceAuth.From, big.NewInt(testAmount), nil)
	if err != nil {
		log.Fatalf("Failed to RequestKLAYTransfer: %v", err)
	}
	info.sim.Commit()
	assert.Nil(info.t, bind.CheckWaitMined(info.sim, tx))
}

func handleKLAYTransfer(info *testInfo, bi *BridgeInfo, ev IRequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	assert.Equal(info.t, new(big.Int).SetUint64(testAmount), ev.GetValueOrTokenId())
	opts := bi.account.GenerateTransactOpts()
	tx, err := bi.bridge.HandleKLAYTransfer(opts, ev.GetRaw().TxHash, ev.GetFrom(), ev.GetTo(), ev.GetValueOrTokenId(), ev.GetRequestNonce(), ev.GetRaw().BlockNumber, ev.GetExtraData())
	if err != nil {
		log.Fatalf("\tFailed to HandleKLAYTransfer: %v", err)
	}
	info.sim.Commit()
	assert.Nil(info.t, bind.CheckWaitMined(info.sim, tx))
}

// TODO-Klaytn-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestKLAYTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents() {
		handleKLAYTransfer(info, bi, ev.(RequestValueTransferEvent))
	}
	info.sim.Commit()
}

func requestTokenTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	var tx *types.Transaction

	testToken := big.NewInt(testToken)
	opts := bi.account.GenerateTransactOpts()

	var err error
	if bi.onChildChain {
		tx, err = info.tokenLocalBridge.RequestValueTransfer(opts, testToken, info.chainAuth.From, big.NewInt(0), nil)
	} else {
		tx, err = info.tokenRemoteBridge.RequestValueTransfer(opts, testToken, info.nodeAuth.From, big.NewInt(0), nil)
	}

	if err != nil {
		log.Fatalf("Failed to RequestValueTransfer for charging: %v", err)
	}
	info.sim.Commit()
	assert.Nil(info.t, bind.CheckWaitMined(info.sim, tx))
}

func handleTokenTransfer(info *testInfo, bi *BridgeInfo, ev IRequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	assert.Equal(info.t, new(big.Int).SetUint64(testToken), ev.GetValueOrTokenId())
	tx, err := bi.bridge.HandleERC20Transfer(
		bi.account.GenerateTransactOpts(), ev.GetRaw().TxHash, ev.GetFrom(), ev.GetTo(), info.tokenRemoteAddr, ev.GetValueOrTokenId(), ev.GetRequestNonce(), ev.GetRaw().BlockNumber, ev.GetExtraData())
	if err != nil {
		log.Fatalf("Failed to HandleERC20Transfer: %v", err)
	}
	info.sim.Commit()
	assert.Nil(info.t, bind.CheckWaitMined(info.sim, tx))
}

// TODO-Klaytn-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestTokenTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents() {
		handleTokenTransfer(info, bi, ev.(RequestValueTransferEvent))
	}
	info.sim.Commit()
}

func requestNFTTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	var tx *types.Transaction

	opts := bi.account.GenerateTransactOpts()
	// TODO-Klaytn need to separate child / parent chain nftIndex.
	nftIndex := new(big.Int).SetInt64(info.nftIndex)

	var err error
	if bi.onChildChain {
		tx, err = info.nftLocalBridge.RequestValueTransfer(opts, nftIndex, info.aliceAuth.From, nil)
	} else {
		tx, err = info.nftRemoteBridge.RequestValueTransfer(opts, nftIndex, info.aliceAuth.From, nil)
	}

	if err != nil {
		log.Fatalf("Failed to requestNFTTransfer for charging: %v", err)
	}
	info.nftIndex++
	info.sim.Commit()
	assert.Nil(info.t, bind.CheckWaitMined(info.sim, tx))
}

func handleNFTTransfer(info *testInfo, bi *BridgeInfo, ev IRequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	var nftAddr common.Address

	if bi.onChildChain {
		nftAddr = info.nftLocalAddr
	} else {
		nftAddr = info.nftRemoteAddr
	}
	uri := GetURI(ev)
	tx, err := bi.bridge.HandleERC721Transfer(
		bi.account.GenerateTransactOpts(),
		ev.GetRaw().TxHash, ev.GetFrom(), ev.GetTo(), nftAddr, ev.GetValueOrTokenId(),
		ev.GetRequestNonce(), ev.GetRaw().BlockNumber, uri, ev.GetExtraData())
	if err != nil {
		log.Fatalf("Failed to handleERC721Transfer: %v", err)
	}
	info.sim.Commit()
	assert.Nil(info.t, bind.CheckWaitMined(info.sim, tx))
}

// TODO-Klaytn-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestNFTTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents() {
		handleNFTTransfer(info, bi, ev)
	}
	info.sim.Commit()
}
