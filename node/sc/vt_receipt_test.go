package sc

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/cmd/homi/setup"
	"github.com/klaytn/klaytn/common"
	revertcontract "github.com/klaytn/klaytn/contracts/revert_test"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/secp256k1"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

const (
	SUCCESS = iota
	OUT_OF_GAS_TEST
	NOT_ENOUGH_CONTRACT_BALANCE_TEST
	REVERT_ON_THE_OTHER_ADDR_TEST
	UNEXECUTED_TEST

	TEST_AMOUNT_OF_KLAY = uint64(100)
)

func handle(t *testing.T, ev IRequestValueTransferEvent, bi *BridgeInfo, backend *backends.SimulatedBackend) {
	txHash := ev.GetRaw().TxHash
	bridge := bi.bridge
	done, err := bridge.HandledRequestTx(nil, txHash)
	assert.NoError(t, err)
	assert.Equal(t, false, done)

	// insert the value transfer request event to the bridge info's event list.
	bi.AddRequestValueTransferEvents([]IRequestValueTransferEvent{ev})

	// handle the value transfer request event in the event list.
	bi.processingPendingRequestEvents()

	backend.Commit() // block
}

func recovery(t *testing.T, sim *backends.SimulatedBackend, vtr *valueTransferRecovery) {
	err := vtr.updateRecoveryHint()
	assert.NoError(t, err)
	err = vtr.retrievePendingEvents()
	assert.NoError(t, err)
	vtr.lookupReceipt()
	sim.Commit()
}

func revertConfiguration(t *testing.T, sim *backends.SimulatedBackend, bm *BridgeManager, bi *BridgeInfo, TEST_CASE int, auth *bind.TransactOpts) {
	switch TEST_CASE {
	case OUT_OF_GAS_TEST:
		bm.subBridge.APIBackend.SetChildBridgeOperatorGasLimit(DefaultBridgeTxGasLimit)
	case NOT_ENOUGH_CONTRACT_BALANCE_TEST:
		chargeKLAY(t, sim, bi, auth)
	}
}

func getBalance(t *testing.T, sim *backends.SimulatedBackend, addr common.Address) uint64 {
	balance, err := sim.BalanceAt(context.Background(), addr, nil)
	assert.Equal(t, nil, err)
	return balance.Uint64()
}

func chargeKLAY(t *testing.T, sim *backends.SimulatedBackend, bi *BridgeInfo, auth *bind.TransactOpts) {
	tx, err := bi.bridge.ChargeWithoutEvent(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, Value: big.NewInt(int64(10000)), GasLimit: testGasLimit})
	assert.Equal(t, nil, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

func sendKLAY(t *testing.T, sim *backends.SimulatedBackend, from *bind.TransactOpts, to common.Address, value *big.Int) {
	nonce, _ := sim.PendingNonceAt(context.Background(), from.From)
	gasPrice, _ := sim.SuggestGasPrice(context.Background())
	unsignedTx := types.NewTransaction(nonce, to, value, testGasLimit, gasPrice, []byte{})

	chainID, _ := sim.ChainID(context.Background())
	tx, err := from.Signer(types.LatestSignerForChainID(chainID), from.From, unsignedTx)
	sim.SendTransaction(context.Background(), tx)
	assert.Equal(t, nil, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

func isExpected(t *testing.T, sim *backends.SimulatedBackend, TEST_CASE int, from, to common.Address, prevFromBalance, prevToBalance uint64) {
	fromBalance, toBalance := getBalance(t, sim, from), getBalance(t, sim, to)
	switch TEST_CASE {
	case SUCCESS, UNEXECUTED_TEST:
		assert.Equal(t, prevFromBalance-TEST_AMOUNT_OF_KLAY, fromBalance)
		assert.Equal(t, prevToBalance+TEST_AMOUNT_OF_KLAY, toBalance)
	case OUT_OF_GAS_TEST, NOT_ENOUGH_CONTRACT_BALANCE_TEST, REVERT_ON_THE_OTHER_ADDR_TEST:
		assert.Equal(t, prevFromBalance, fromBalance)
		assert.Equal(t, prevToBalance, toBalance)
	}
}

func balanceCheck(t *testing.T, sim *backends.SimulatedBackend, users []*bind.TransactOpts, expected uint64) {
	for _, user := range users {
		assert.Equal(t, getBalance(t, sim, user.From), expected)
	}
}

func sendRefundCall(t *testing.T, sim *backends.SimulatedBackend, operator *bind.TransactOpts, bi *BridgeInfo, reqNonce uint64) *types.Transaction {
	auth := &bind.TransactOpts{From: operator.From, Signer: operator.Signer, GasLimit: testGasLimit}
	tx, err := bi.bridge.HandleRefund(auth, reqNonce)
	assert.NoError(t, err)
	sim.Commit()
	return tx
}

func sendWithdrawKLAYCall(t *testing.T, sim *backends.SimulatedBackend, operator *bind.TransactOpts, bi *BridgeInfo, value *big.Int) *types.Transaction {
	auth := &bind.TransactOpts{From: operator.From, Signer: operator.Signer, GasLimit: testGasLimit}
	tx, err := bi.bridge.WithdrawKLAY(auth, value)
	assert.NoError(t, err)
	sim.Commit()
	return tx
}

func testKLAYOutofGasFailure(t *testing.T, bm *BridgeManager) {
	t.Log("parent operator gas limit", bm.subBridge.APIBackend.GetChildBridgeOperatorGasLimit())
	bm.subBridge.APIBackend.SetChildBridgeOperatorGasLimit(100000)
	t.Log("parent operator gas limit", bm.subBridge.APIBackend.GetChildBridgeOperatorGasLimit())
}

func testKLAYNotEnoughContractBalanceFailure(t *testing.T, sim *backends.SimulatedBackend, bi *BridgeInfo) {
	// Withdraw all KLAY to test not enough balance error
	auth := bi.account.GenerateTransactOpts()
	bridgeBalance := getBalance(t, sim, bi.address)
	tx, err := bi.bridge.WithdrawKLAY(auth, big.NewInt(int64(bridgeBalance)))
	assert.NoError(t, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	bridgeBalance = getBalance(t, sim, bi.address)
}

func testKLAYRevertOnTheOtherAddrFailure(t *testing.T, sim *backends.SimulatedBackend, bi *BridgeInfo) common.Address {
	auth := bi.account.GenerateTransactOpts()
	contractAddr, tx, _, err := revertcontract.DeployRevertContract(auth, sim)
	assert.NoError(t, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	return contractAddr
}

func getCounterpartBridgeInfo(t *testing.T, bm *BridgeManager, addr common.Address) *BridgeInfo {
	reqBridgeAddr := bm.GetCounterPartBridgeAddr(addr)
	assert.NotEqual(t, reqBridgeAddr, common.Address{})
	reqBridgeInfo, ok := bm.GetBridgeInfo(reqBridgeAddr)
	assert.Equal(t, ok, true)
	return reqBridgeInfo
}

func bridgeSetup(t *testing.T) (*backends.SimulatedBackend,
	*BridgeManager, *BridgeInfo, *BridgeInfo,
	*valueTransferRecovery, *bind.TransactOpts, *bind.TransactOpts, string,
) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Generate a new random account and a funded simulator
	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)
	bobKey, _ := crypto.GenerateKey()
	bob := bind.NewKeyedTransactor(bobKey)

	config := &SCConfig{}
	config.DataDir = tempDir

	bacc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	cAlloc := blockchain.GenesisAlloc{
		alice.From:            {Balance: big.NewInt(params.KLAY)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KLAY)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KLAY)},
	}
	sim := backends.NewSimulatedBackend(cAlloc)

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		blockchain:     sim.BlockChain(),
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	sc.APIBackend = &SubBridgeAPI{sc}
	sc.handler, err = NewSubBridgeHandler(sc)
	assert.NoError(t, err)

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)
	sc.handler.subbridge.bridgeManager = bm

	cBridgeAddr := deployBridge(t, bm, sim, true)
	pBridgeAddr := deployBridge(t, bm, sim, false)
	err = bm.subBridge.APIBackend.RegisterBridge(cBridgeAddr, pBridgeAddr, nil)
	assert.NoError(t, err)
	cBridgeAddrStr, pBridgeAddrStr := cBridgeAddr.String(), pBridgeAddr.String()
	err = bm.subBridge.APIBackend.SubscribeBridge(&cBridgeAddrStr, &pBridgeAddrStr)
	assert.NoError(t, err)
	cbi, _ := bm.GetBridgeInfo(cBridgeAddr)
	pbi, _ := bm.GetBridgeInfo(pBridgeAddr)
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, cbi, pbi)
	return sim, bm, cbi, pbi, vtr, alice, bob, node.DefaultConfig.ResolvePath("chaindata")
}

func handleLoop(t *testing.T, sim *backends.SimulatedBackend, bm *BridgeManager, bi *BridgeInfo, handled chan<- struct{}, unexecutedTest bool, exit chan struct{}) {
	reqVTevCh := make(chan RequestValueTransferEvent)
	handleValueTransferEventCh := make(chan *HandleValueTransferEvent)
	requestRefundEventCh := make(chan *RequestRefundEvent)
	handleRefundEventCh := make(chan *HandleRefundEvent)

	bm.SubscribeReqVTev(reqVTevCh)
	bm.SubscribeHandleVTev(handleValueTransferEventCh)
	bm.SubscribeReqRefundEv(requestRefundEventCh)
	bm.SubscribeHandleRefundEv(handleRefundEventCh)

	// Handle the request and handle events
	nReq := 0
	go func() {
		for {
			select {
			case ev := <-reqVTevCh:
				if nReq == UNEXECUTED_TEST && unexecutedTest {
					t.Log("Skip the handle value transfer execution")
				} else {
					handle(t, ev, bi, sim)
				}
				handled <- struct{}{}
				nReq++
			case ev := <-handleValueTransferEventCh:
				reqBridgeInfo := getCounterpartBridgeInfo(t, bm, bi.address)
				auth := reqBridgeInfo.account.GenerateTransactOpts()
				auth.GasLimit = params.UpperGasLimit

				tx, err := reqBridgeInfo.bridge.RemoveRefundLedger(auth, ev.HandleNonce)
				assert.NoError(t, err)
				sim.Commit()
				CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
			case ev := <-requestRefundEventCh:
				reqBridgeInfo := getCounterpartBridgeInfo(t, bm, bi.address)
				auth := reqBridgeInfo.account.GenerateTransactOpts()
				auth.GasLimit = params.UpperGasLimit

				tx, err := reqBridgeInfo.bridge.HandleRefund(auth, ev.RequestNonce)
				assert.NoError(t, err)
				sim.Commit()
				CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
			case ev := <-handleRefundEventCh:
				auth := bi.account.GenerateTransactOpts()
				auth.GasLimit = params.UpperGasLimit

				reqNonce, reqTxHash, reqBlkNum := ev.RequestNonce, ev.Raw.TxHash, ev.Raw.BlockNumber
				tx, err := bi.bridge.UpdateHandleStatus(auth, reqNonce, reqTxHash, reqBlkNum, true)
				assert.NoError(t, err)
				sim.Commit()
				CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
				handled <- struct{}{}
			case <-exit:
				exit <- struct{}{}
				return
			}
		}
	}()
}

func cleanup(t *testing.T, sim *backends.SimulatedBackend, dataDir string, exit chan struct{}) {
	exit <- struct{}{}
	<-exit
	sim.Close()
	if err := os.RemoveAll(dataDir); err != nil {
		t.Fatalf("fail to delete file %v", err)
	}
}

func TestHandleVTFailure(t *testing.T) {
	sim, bm, cbi, pbi, vtr, alice, bob, dataDir := bridgeSetup(t)
	exit := make(chan struct{}, 1)
	defer cleanup(t, sim, dataDir, exit)

	nRequest := 5
	handled := make(chan struct{}, 1)
	handleLoop(t, sim, bm, cbi, handled, true, exit)

	// Do request KLAY transfer
	for TEST_CASE := 0; TEST_CASE < nRequest; TEST_CASE++ {
		revertContractAddr := common.Address{}
		switch TEST_CASE {
		case OUT_OF_GAS_TEST:
			testKLAYOutofGasFailure(t, bm)
		case NOT_ENOUGH_CONTRACT_BALANCE_TEST:
			testKLAYNotEnoughContractBalanceFailure(t, sim, cbi)
		case REVERT_ON_THE_OTHER_ADDR_TEST:
			revertContractAddr = testKLAYRevertOnTheOtherAddrFailure(t, sim, cbi)
		}
		var to common.Address
		if revertContractAddr == (common.Address{}) {
			to = bob.From
		} else {
			to = revertContractAddr
		}
		prevAliceBalance, prevBobBalance := getBalance(t, sim, alice.From), getBalance(t, sim, bob.From)
		tx, err := pbi.bridge.RequestKLAYTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, Value: big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), GasLimit: testGasLimit}, to, big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), nil)
		assert.NoError(t, err)
		sim.Commit()
		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		<-handled
		recovery(t, sim, vtr)
		if TEST_CASE == OUT_OF_GAS_TEST || TEST_CASE == NOT_ENOUGH_CONTRACT_BALANCE_TEST || TEST_CASE == REVERT_ON_THE_OTHER_ADDR_TEST {
			<-handled
		}
		if TEST_CASE == UNEXECUTED_TEST {
			assert.Equal(t, 1, len(vtr.parentEvents))
			cbi.AddRequestValueTransferEvents(vtr.parentEvents)
			cbi.processingPendingRequestEvents()
			sim.Commit()
		}
		sim.Commit()
		isExpected(t, sim, TEST_CASE, alice.From, bob.From, prevAliceBalance, prevBobBalance)
		revertConfiguration(t, sim, bm, cbi, TEST_CASE, alice)
	}

	senderChainReqNonce, err := pbi.bridge.RequestNonce(nil)
	assert.NoError(t, err)
	senderChainNumberOfRefunds, err := pbi.bridge.NRefunds(nil)
	assert.NoError(t, err)
	receiverChainLowerHandleNonce, err := cbi.bridge.LowerHandleNonce(nil)
	assert.NoError(t, err)
	receiverChianFailedHandles, err := cbi.bridge.FailedHandles(nil)
	assert.NoError(t, err)
	t.Logf(
		"senderChainReqNonce %d\n senderChainNumberOfRefunds %d\n receiverChianFailedHandles %d\n receiverChainLowerHandleNonce %d\n",
		senderChainReqNonce, senderChainNumberOfRefunds, receiverChianFailedHandles, receiverChainLowerHandleNonce)
	assert.Equal(t, senderChainReqNonce, receiverChainLowerHandleNonce)
	assert.Equal(t, senderChainNumberOfRefunds, receiverChianFailedHandles)

	// An arbitrary sleep bypasses the access of deallocated memory.
	time.Sleep(time.Second * 5)
}

func TestWithdraw(t *testing.T) {
	{
		sim, _, _, pbi, _, alice, bob, dataDir := bridgeSetup(t)
		// 1. Send KLAY three times in the circustance of absnces of bridge node
		nRequest := 3
		for i := 0; i < nRequest; i++ {
			tx, err := pbi.bridge.RequestKLAYTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, Value: big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), GasLimit: testGasLimit}, bob.From, big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), nil)
			assert.NoError(t, err)
			sim.Commit()
			CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		}

		prevBridgeBalance := getBalance(t, sim, pbi.address)
		prevOperatorBalance := getBalance(t, sim, pbi.account.address)
		t.Log("bridge balance", prevBridgeBalance)
		t.Log("operator balance", prevOperatorBalance)

		// 2. Try to refund all the KLAY of bridge contract holds.
		auth := pbi.account.GenerateTransactOpts()
		tx, err := pbi.bridge.WithdrawKLAY(auth, big.NewInt(int64(prevBridgeBalance)))
		assert.NoError(t, err)
		sim.Commit()
		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		bridgeBalance := getBalance(t, sim, pbi.address)
		operatorBalance := getBalance(t, sim, pbi.account.address)
		t.Log("bridge balance", bridgeBalance)
		t.Log("operator balance", operatorBalance)

		assert.Equal(t, bridgeBalance, TEST_AMOUNT_OF_KLAY*3)
		assert.Equal(t, operatorBalance-prevOperatorBalance, prevBridgeBalance-bridgeBalance)
		exit := make(chan struct{}, 1)
		cleanup(t, sim, dataDir, exit)
	}

	{
		sim, bm, cbi, pbi, _, alice, bob, dataDir := bridgeSetup(t)
		exit := make(chan struct{}, 1)
		defer cleanup(t, sim, dataDir, exit)

		nRequest := 3
		handled := make(chan struct{}, 1)
		handleLoop(t, sim, bm, cbi, handled, false, exit)

		// 1. Send KLAY three times
		for i := 0; i < nRequest; i++ {
			tx, err := pbi.bridge.RequestKLAYTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, Value: big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), GasLimit: testGasLimit}, bob.From, big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), nil)
			assert.NoError(t, err)
			sim.Commit()
			CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
			<-handled
			sim.Commit()
		}

		// Wait until `RemoveRefundLedger` transaction is done
		time.Sleep(time.Second * 5)

		prevBridgeBalance := getBalance(t, sim, pbi.address)
		prevOperatorBalance := getBalance(t, sim, pbi.account.address)
		t.Log("bridge balance", prevBridgeBalance)
		t.Log("operator balance", prevOperatorBalance)

		// 2. Try to refund all the KLAY of bridge contract holds.
		auth := pbi.account.GenerateTransactOpts()
		tx, err := pbi.bridge.WithdrawKLAY(auth, big.NewInt(int64(prevBridgeBalance)))
		assert.NoError(t, err)
		sim.Commit()
		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		bridgeBalance := getBalance(t, sim, pbi.address)
		operatorBalance := getBalance(t, sim, pbi.account.address)
		t.Log("bridge balance", bridgeBalance)
		t.Log("operator balance", operatorBalance)

		assert.Equal(t, bridgeBalance, uint64(0))
		assert.Equal(t, operatorBalance, prevOperatorBalance+prevBridgeBalance)
	}
}

func TestSuggestLeastFeeForKLAYTransfer(t *testing.T) {
	sim, bm, cbi, pbi, _, _, _, dataDir := bridgeSetup(t)
	exit := make(chan struct{}, 1)
	defer cleanup(t, sim, dataDir, exit)

	ston25 := big.NewInt(1000000000 * 25)
	ston750 := big.NewInt(1000000000 * 750)
	cbi.account.gasPrice = ston25
	pbi.account.gasPrice = ston750

	cDetail, err := bm.subBridge.APIBackend.SuggestLeastFee(cbi.address, "KLAY")
	assert.NoError(t, err)
	t.Log(cDetail)

	pDetail, err := bm.subBridge.APIBackend.SuggestLeastFee(pbi.address, "KLAY")
	assert.NoError(t, err)
	t.Log(pDetail)

	cDetailMap := map[string]interface{}(cDetail)
	pDetailMap := map[string]interface{}(pDetail)

	// A calcaulted total cost from child bridge is larger than the parent bridge's one
	// because the child bridge contract has more contract calls for parent bridge contract.
	assert.Equal(t, cDetailMap["SumOfCost"].(uint64) > pDetailMap["SumOfCost"].(uint64), true)
	assert.Equal(t, cDetailMap["SumOfGasUsed"], pDetailMap["SumOfGasUsed"])
}

func TestQueryOfLeastAmountKLAY(t *testing.T) {
	sim, _, cbi, _, _, _, bob, dataDir := bridgeSetup(t)
	exit := make(chan struct{}, 1)
	defer cleanup(t, sim, dataDir, exit)

	wantToSend := big.NewInt(100)

	k, err := cbi.bridge.GetMinimumAmountOfKLAY(nil, wantToSend)
	assert.NoError(t, err)
	assert.Equal(t, k, wantToSend)

	configNonce, err := cbi.bridge.ConfigurationNonce(nil)
	assert.NoError(t, err)

	// Set fee receiver
	auth := cbi.account.GenerateTransactOpts()
	tx, err := cbi.bridge.SetFeeReceiver(auth, bob.From)
	assert.NoError(t, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	sim.Commit() // block

	// Set KLAY fee
	fee := big.NewInt(50)
	tx, err = cbi.bridge.SetKLAYFee(auth, fee, configNonce)
	assert.NoError(t, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	k, err = cbi.bridge.GetMinimumAmountOfKLAY(nil, wantToSend)
	assert.NoError(t, err)
	assert.Equal(t, k, new(big.Int).Add(wantToSend, fee))
}

func TestMultiBridgeOperation(t *testing.T) {
	sim, _, _, pbi, _, _, bob, dataDir := bridgeSetup(t)
	exit := make(chan struct{}, 1)
	defer cleanup(t, sim, dataDir, exit)

	// Prepare N of operators and N of EOAs
	nOperators := 3
	operators := make([]*bind.TransactOpts, nOperators)
	users := make([]*bind.TransactOpts, nOperators)
	auth := pbi.account.GenerateTransactOpts()
	for i := 0; i < nOperators; i++ {
		key, _ := crypto.GenerateKey()
		operators[i] = bind.NewKeyedTransactor(key)
		tx, err := pbi.bridge.RegisterOperator(auth, operators[i].From)
		assert.NoError(t, err)
		sim.Commit()
		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// Prepare refund ledger
		key, _ = crypto.GenerateKey()
		users[i] = bind.NewKeyedTransactor(key)
		sendKLAY(t, sim, auth, users[i].From, big.NewInt(100))
	}

	// Set threshold for refund vote
	tx, err := pbi.bridge.SetOperatorThreshold(auth, voteTypeHandleRefund, uint8(nOperators))
	assert.NoError(t, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// Set threshold for withdraw vote
	tx, err = pbi.bridge.SetOperatorThreshold(auth, voteTypeWithdraw, uint8(nOperators))
	assert.NoError(t, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// (1) Test of withdraw
	{
		expectedBalance := uint64(100)
		for _, user := range users {
			assert.Equal(t, getBalance(t, sim, user.From), expectedBalance)
		}
		tobeSent := big.NewInt(3)
		for _, user := range users {
			tx, err := pbi.bridge.RequestKLAYTransfer(&bind.TransactOpts{From: user.From, Signer: user.Signer, Value: tobeSent, GasLimit: testGasLimit}, bob.From, tobeSent, nil)
			assert.NoError(t, err)
			sim.Commit()
			CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		}
		balanceCheck(t, sim, users, expectedBalance-tobeSent.Uint64())

		// TEST 1-1 - Failure: quorum not satisfied
		reqNonce, err := pbi.bridge.RequestNonce(nil)
		assert.NoError(t, err)
		reqNonces := make([]uint64, reqNonce)
		for i := 0; i < int(reqNonce); i++ {
			reqNonces[i] = uint64(i)
		}
		for _, reqNonce := range reqNonces {
			for i := 0; i < len(operators)-1; i++ {
				tx := sendRefundCall(t, sim, operators[i], pbi, reqNonce)
				CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
			}
		}
		balanceCheck(t, sim, users, expectedBalance-tobeSent.Uint64())

		// TEST 1-2 - Success: quorum satisfied
		for _, reqNonce := range reqNonces {
			tx := sendRefundCall(t, sim, operators[len(operators)-1], pbi, reqNonce)
			CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		}
		balanceCheck(t, sim, users, expectedBalance)

		// TEST 1-3 - Failure: vote to closed nonce
		for _, reqNonce := range reqNonces {
			for _, operator := range operators {
				tx := sendRefundCall(t, sim, operator, pbi, reqNonce)
				CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
			}
		}
	}

	// (2) Test of refund KLAY
	{
		ownerBalance := getBalance(t, sim, pbi.account.address)
		bridgeBalance := getBalance(t, sim, pbi.address)

		// TEST 2-1 - Failure: quorum not satisfied
		for i := 0; i < len(operators)-1; i++ {
			tx := sendWithdrawKLAYCall(t, sim, operators[i], pbi, big.NewInt(int64(bridgeBalance)))
			CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		}
		assert.Equal(t, bridgeBalance, getBalance(t, sim, pbi.address))

		// TEST 2-2 - Success: quorum satisfied
		tx := sendWithdrawKLAYCall(t, sim, operators[len(operators)-1], pbi, big.NewInt(int64(bridgeBalance)))
		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		withdrawal := uint64(0)
		for _, operator := range operators {
			withdrawal += getBalance(t, sim, operator.From)
		}
		withdrawal += getBalance(t, sim, pbi.account.address) - ownerBalance
		assert.Equal(t, getBalance(t, sim, pbi.address), bridgeBalance-withdrawal)
	}
}

func makeParentOperatorKey(t *testing.T, path string) {
	password := setup.RandStringRunes(params.PasswordLength)
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount(password)
	assert.NoError(t, err)
	setup.WriteFile([]byte(password), path, acc.Address.String())
}

func getPubkey(t *testing.T, dir string, h []byte) (*accountInfo, []byte) {
	acc := parentAccInit(nil, dir, params.UpperGasLimit)
	sig, err := acc.keystore.SignHash(accounts.Account{Address: acc.address}, h)
	assert.NoError(t, err)
	pubkey, err := secp256k1.RecoverPubkey(h, sig)
	assert.NoError(t, err)
	return acc, pubkey
}

func isSamePubkey(t *testing.T, operatorPubkey, sig, data []byte, shouldBeTrue bool) {
	recovered, err := secp256k1.RecoverPubkey(data, sig)
	assert.NoError(t, err)
	if shouldBeTrue {
		assert.Equal(t, bytes.Compare(recovered, operatorPubkey), 0)
	} else {
		assert.NotEqual(t, bytes.Compare(recovered, operatorPubkey), 0)
	}
}

func makeHash() ([]byte, []byte, []byte) {
	txHash1 := [common.HashLength]byte(common.HexToHash("0x1"))
	txHash2 := [common.HashLength]byte(common.HexToHash("0x2"))
	txHash3 := [common.HashLength]byte(common.HexToHash("0x3"))
	return txHash1[:], txHash2[:], txHash3[:]
}

func TestMsgIntegrity(t *testing.T) {
	operatorKeyDir := "sc_test_temporal_key"
	attackerKeyDir := "sc_test_attacker_key"
	makeParentOperatorKey(t, path.Join(operatorKeyDir, ParentBridgeAccountName))
	makeParentOperatorKey(t, path.Join(attackerKeyDir, ParentBridgeAccountName))
	defer func() {
		if err := os.RemoveAll(operatorKeyDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
		if err := os.RemoveAll(attackerKeyDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	txHash1, txHash2, txHash3 := makeHash()
	operatorAcc, operatorPubkey := getPubkey(t, operatorKeyDir, txHash1)
	attackerAcc, _ := getPubkey(t, attackerKeyDir, txHash1)

	operatorSig, err := operatorAcc.keystore.SignHash(accounts.Account{Address: operatorAcc.address}, txHash2)
	assert.NoError(t, err)
	attackerSig, err := attackerAcc.keystore.SignHash(accounts.Account{Address: attackerAcc.address}, txHash2)
	assert.NoError(t, err)
	isSamePubkey(t, operatorPubkey, operatorSig, txHash2, true)
	isSamePubkey(t, operatorPubkey, attackerSig, txHash2, false)
	isSamePubkey(t, operatorPubkey, operatorSig, txHash3, false)
}

func TestHandleRefundRecovery(t *testing.T) {
	sim, bm, cbi, pbi, vtr, alice, bob, dataDir := bridgeSetup(t)
	exit := make(chan struct{}, 1)
	defer cleanup(t, sim, dataDir, exit)

	handled := make(chan struct{}, 1)
	testHandleRefundFailure(t, sim, vtr, bm, cbi, handled, true, exit)

	// Do request KLAY transfer
	testKLAYNotEnoughContractBalanceFailure(t, sim, cbi)
	to := bob.From
	prevAliceBalance, prevBobBalance := getBalance(t, sim, alice.From), getBalance(t, sim, bob.From)
	tx, err := pbi.bridge.RequestKLAYTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, Value: big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), GasLimit: testGasLimit}, to, big.NewInt(int64(TEST_AMOUNT_OF_KLAY)), nil)
	assert.NoError(t, err)
	sim.Commit()
	CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	<-handled
	recovery(t, sim, vtr)
	<-handled

	sim.Commit()
	isExpected(t, sim, NOT_ENOUGH_CONTRACT_BALANCE_TEST, alice.From, bob.From, prevAliceBalance, prevBobBalance)
	revertConfiguration(t, sim, bm, cbi, NOT_ENOUGH_CONTRACT_BALANCE_TEST, alice)

	// An arbitrary sleep bypasses the access of deallocated memory.
	time.Sleep(time.Second * 5)
}

func testHandleRefundFailure(t *testing.T, sim *backends.SimulatedBackend, vtr *valueTransferRecovery, bm *BridgeManager, bi *BridgeInfo, handled chan<- struct{}, unexecutedTest bool, exit chan struct{}) {
	reqVTevCh := make(chan RequestValueTransferEvent)
	handleValueTransferEventCh := make(chan *HandleValueTransferEvent)
	requestRefundEventCh := make(chan *RequestRefundEvent)
	handleRefundEventCh := make(chan *HandleRefundEvent)

	bm.SubscribeReqVTev(reqVTevCh)
	bm.SubscribeHandleVTev(handleValueTransferEventCh)
	bm.SubscribeReqRefundEv(requestRefundEventCh)
	bm.SubscribeHandleRefundEv(handleRefundEventCh)

	// Handle the request and handle events
	nReceivedHandleRefund := 0
	go func() {
		for {
			select {
			case ev := <-reqVTevCh:
				handle(t, ev, bi, sim)
				handled <- struct{}{}
			case ev := <-handleValueTransferEventCh:
				reqBridgeInfo := getCounterpartBridgeInfo(t, bm, bi.address)
				auth := reqBridgeInfo.account.GenerateTransactOpts()
				auth.GasLimit = params.UpperGasLimit

				tx, err := reqBridgeInfo.bridge.RemoveRefundLedger(auth, ev.HandleNonce)
				assert.NoError(t, err)
				sim.Commit()
				CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
			case ev := <-requestRefundEventCh:
				t.Log("RequestRefund event received")
				if nReceivedHandleRefund == 0 {
					t.Log("Do not execute the refund handle if its receive is the first time")
					nReceivedHandleRefund++
					recovery(t, sim, vtr)
				} else {
					reqBridgeInfo := getCounterpartBridgeInfo(t, bm, bi.address)
					auth := reqBridgeInfo.account.GenerateTransactOpts()
					auth.GasLimit = params.UpperGasLimit

					tx, err := reqBridgeInfo.bridge.HandleRefund(auth, ev.RequestNonce)
					assert.NoError(t, err)
					sim.Commit()
					CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
				}
			case ev := <-handleRefundEventCh:
				auth := bi.account.GenerateTransactOpts()
				auth.GasLimit = params.UpperGasLimit

				reqNonce, reqTxHash, reqBlkNum := ev.RequestNonce, ev.Raw.TxHash, ev.Raw.BlockNumber
				tx, err := bi.bridge.UpdateHandleStatus(auth, reqNonce, reqTxHash, reqBlkNum, true)
				assert.NoError(t, err)
				sim.Commit()
				CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
				handled <- struct{}{}
			case <-exit:
				exit <- struct{}{}
				return
			}
		}
	}()
}
