package payments

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/params"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"

	"github.com/klaytn/klaytn/crypto"

	"github.com/klaytn/klaytn/accounts/abi/bind/backends"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultGasLimit = 5000000
)

// CheckReceipt can check if the tx receipt has expected status.
func CheckReceipt(b bind.DeployBackend, tx *types.Transaction, duration time.Duration, expectedStatus uint, t *testing.T) {
	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), duration)
	defer cancelTimeout()

	receipt, err := bind.WaitMined(timeoutContext, b, tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedStatus, receipt.Status)
}

type paymentTestENV struct {
	backend                *backends.SimulatedBackend
	accounts               []*bind.TransactOpts
	paymentContract        *Payment
	paymentContractAddress common.Address
	settleInterval         int
}

func generatePaymentTestEnv(t *testing.T) *paymentTestENV {
	// create account
	accountNum := 20
	keys := make([]*ecdsa.PrivateKey, accountNum)
	accounts := make([]*bind.TransactOpts, accountNum)
	for i := 0; i < accountNum; i++ {
		keys[i], _ = crypto.GenerateKey()
		accounts[i] = bind.NewKeyedTransactor(keys[i])
		accounts[i].GasLimit = DefaultGasLimit
	}

	// generate backend with deployed
	alloc := blockchain.GenesisAlloc{}
	for i := 0; i < accountNum; i++ {
		alloc[accounts[i].From] = blockchain.GenesisAccount{
			Balance: big.NewInt(params.KLAY),
		}
	}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy
	settleInterval := 100
	contractAddress, tx, contract, err := DeployPayment(accounts[0], backend, big.NewInt(int64(settleInterval)))
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	return &paymentTestENV{
		backend:                backend,
		accounts:               accounts,
		paymentContract:        contract,
		paymentContractAddress: contractAddress,
		settleInterval:         settleInterval,
	}
}

func TestPayment_sendPayment(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	thirdParty := env.accounts[2]
	contract := env.paymentContract

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()

	// A → B에게 X만큼 sendPayment를 실행
	sender.Value = big.NewInt(300)
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// SendPayment Event의 Hash값 확인
	item := <-sink
	hash := item.Hash

	// B는 getPayments 실행 시 결과를 확인할 수 있음
	getPayment, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getPayment))
	assert.Equal(t, getPayment, [][32]byte{hash})

	// C는 getPayments 실행 시 결과를 확인할 수 없음
	getPayment, err = contract.GetPayments(&bind.CallOpts{From: thirdParty.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayment))
}

func TestPayment_sendPayment_noZero(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract

	// value == 0 이면 트랜잭션 실패
	sender.Value = big.NewInt(0)
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// value > 0 이면 트랜잭션 성공
	sender.Value = big.NewInt(1)
	tx, err = contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	sender.Value = big.NewInt(100000)
	tx, err = contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

func TestPayment_sendPayment_balance(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	contractAddress := env.paymentContractAddress
	sendAmount := big.NewInt(300)

	// contract, sender, receiver의 초기 balance 조회
	senderPreviousBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractPreviousBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverPreviousBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	// sendPayment 실행
	sender.Value = sendAmount
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// contract, sender, receiver balance 변화 확인
	senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	assert.Equal(t, big.NewInt(0).Sub(senderPreviousBalance, sendAmount), senderAfterBalance)
	assert.Equal(t, big.NewInt(0).Add(contractPreviousBalance, sendAmount), contractAfterBalance)
	assert.Equal(t, receiverPreviousBalance, receiverAfterBalance)
}

func TestPayment_isAbleToSettle(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	settleInterval := env.settleInterval

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()

	sender.Value = big.NewInt(300)
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// SendPayment Event의 Hash값 확인
	item := <-sink
	hash := item.Hash

	// SettleInterval가 지나기 전에는 isAbleToSettle 실행 시 False인 것을 확인
	{
		isAbleToSettle, err := contract.IsAbleToSettle(&bind.CallOpts{From: receiver.From}, hash)
		assert.NoError(t, err)
		assert.False(t, isAbleToSettle)
	}

	// SettleInterval가 지난 이후에는 isAbleToSettle 실행 시 True인 것을 확인
	{
		// settleInterval 만큼 블럭 생성
		for i := 0; i < settleInterval; i++ {
			backend.Commit()
		}

		isAbleToSettle, err := contract.IsAbleToSettle(&bind.CallOpts{From: receiver.From}, hash)
		assert.NoError(t, err)
		assert.True(t, isAbleToSettle)
	}
}

func TestPayment_settlePayment(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	settleInterval := env.settleInterval

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()
	var previousHash [32]byte
	var afterHash [32]byte

	// 블록 생성 전 sendPayment 실행
	{
		sender.Value = big.NewInt(300)
		tx, err := contract.SendPayment(sender, receiver.From)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// SendPayment Event의 Hash값 확인
		item := <-sink
		previousHash = item.Hash
	}
	// settleInterval 만큼 블록 생성
	for i := 0; i < settleInterval; i++ {
		backend.Commit()
	}

	// 블록 생성 후 sendPayment 실행
	{
		sender.Value = big.NewInt(300)
		tx, err := contract.SendPayment(sender, receiver.From)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// SendPayment Event의 Hash값 확인
		item := <-sink
		afterHash = item.Hash
	}
	// getPayments를 호출하여 두 트랜잭션이 리턴되는 것을 확인
	{
		getPayments, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(getPayments))
		assert.Contains(t, getPayments, previousHash)
		assert.Contains(t, getPayments, afterHash)
	}
	// settlePayment를 실행
	{
		tx, err := contract.SettlePayment(receiver)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}
	// getPayments를 호출하여 조건을 만족하지 못한 (블록 생성 후) 트랜잭션이 리턴되는 것을 확인
	{
		getPayment, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(getPayment))
		assert.Equal(t, afterHash, getPayment[0])
	}
}

func TestPayment_settlePayment_balance(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	contractAddress := env.paymentContractAddress
	settleInterval := env.settleInterval
	sendAmount := big.NewInt(300)

	// contract와 receiver의 초기 balance 조회
	senderPreviousBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractPreviousBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverPreviousBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	// sendPayment와 settlePayment 실행
	{
		// sendPayment 실행
		sender.Value = sendAmount
		tx, err := contract.SendPayment(sender, receiver.From)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// settleInterval 만큼 블럭 생성
		for i := 0; i < settleInterval; i++ {
			backend.Commit()
		}

		// settlePayment 실행
		tx, err = contract.SettlePayment(receiver)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// contract, sender, receiver balance 변화 확인
	senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	assert.Equal(t, big.NewInt(0).Sub(senderPreviousBalance, sendAmount), senderAfterBalance)
	assert.Equal(t, contractPreviousBalance, contractAfterBalance)
	assert.Equal(t, big.NewInt(0).Add(receiverPreviousBalance, sendAmount), receiverAfterBalance)
}

func TestPayment_settlePayment_empty(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	receiver := env.accounts[0]
	contract := env.paymentContract

	// getPayments를 호출하여 받을 Payment가 없는 것을 확인
	getPayment, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayment))

	// receiver가 settlePayment 실행
	tx, err := contract.SettlePayment(receiver)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// getPayments를 호출하여 상태 변화가 없음을 확인
	getPayment, err = contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayment))
}

func TestPayment_settlePayment_empty_balance(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	contractAddress := env.paymentContractAddress

	// contract, sender, receiver의 초기 balance 확인
	senderPreviousBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractPreviousBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverPreviousBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	// receiver가 settlePayment 실행
	tx, err := contract.SettlePayment(receiver)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// contract, sender, receiver의 balance 변화 확인
	senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	assert.Equal(t, senderPreviousBalance, senderAfterBalance)
	assert.Equal(t, contractPreviousBalance, contractAfterBalance)
	assert.Equal(t, receiverPreviousBalance, receiverAfterBalance)
}

func TestPayment_settlePayment_pending(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	settleInterval := env.settleInterval

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()
	var hash [32]byte

	// sendPayment 실행
	sender.Value = big.NewInt(300)
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// SendPayment Event의 Hash값 확인
	item := <-sink
	hash = item.Hash

	// SettleInterval가 지나기 전에는 SettlePayment를 호출하여도 getPayments의 상태변화가 없음을 확인
	for i := 0; i < settleInterval-1; i++ {
		tx, err := contract.SettlePayment(receiver)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		getPayments, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(getPayments))
		assert.Contains(t, getPayments, hash)
	}

	// SettleInterval가 지난 후에는 SettlePayment를 호출하면 getPayments에 아무것도 조회되지 않는 것을 확인
	{
		tx, err := contract.SettlePayment(receiver)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		getPayments, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(getPayments))
	}
}

func TestPayment_settlePayment_pending_balance(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	contractAddress := env.paymentContractAddress
	settleInterval := env.settleInterval
	sendAmount := big.NewInt(300)

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()

	// contract, sender, receiver의 초기 balance 확인
	senderPreviousBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractPreviousBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverPreviousBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	// sendPayment 실행
	sender.Value = sendAmount
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// SettleInterval가 지나기 전에는 SettlePayment를 호출하여도 getPayments의 상태변화가 없음을 확인
	for i := 0; i < settleInterval-1; i++ {
		tx, err := contract.SettlePayment(receiver)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// contract, sender, receiver의 balance 변화 확인
		senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
		assert.NoError(t, err)
		contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
		assert.NoError(t, err)
		receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
		assert.NoError(t, err)

		assert.Equal(t, big.NewInt(0).Sub(senderPreviousBalance, sendAmount), senderAfterBalance)
		assert.Equal(t, big.NewInt(0).Add(contractPreviousBalance, sendAmount), contractAfterBalance)
		assert.Equal(t, receiverPreviousBalance, receiverAfterBalance)
	}

	// SettleInterval가 지난 후에는 SettlePayment를 호출하면 getPayments에 아무것도 조회되지 않는 것을 확인
	{
		tx, err := contract.SettlePayment(receiver)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// contract, sender, receiver의 balance 변화 확인
		senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
		assert.NoError(t, err)
		contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
		assert.NoError(t, err)
		receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
		assert.NoError(t, err)

		assert.Equal(t, big.NewInt(0).Sub(senderPreviousBalance, sendAmount), senderAfterBalance)
		assert.Equal(t, contractPreviousBalance, contractAfterBalance)
		assert.Equal(t, big.NewInt(0).Add(receiverPreviousBalance, sendAmount), receiverAfterBalance)
	}
}

func TestPayment_settlePayment_cancelEmpty(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()
	var hash [32]byte

	// sendPayment 실행
	sender.Value = big.NewInt(300)
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// SendPayment Event의 Hash값 확인
	item := <-sink
	hash = item.Hash

	// 받을 Payment를 cancel
	tx, err = contract.CancelPayment(receiver, hash)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// getPayments를 호출하여 받은 Payment가 없음을 확인
	getPayments, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayments))

	// settlePayment를 호출함
	tx, err = contract.SettlePayment(receiver)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// getPayments를 호출하여 받은 Payment가 없음을 확인
	getPayments, err = contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayments))
}

func TestPayment_settlePayment_cancelEmpty_balance(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	contractAddress := env.paymentContractAddress

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()
	var hash [32]byte

	// contract, sender, receiver의 초기 balance 확인
	senderPreviousBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractPreviousBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverPreviousBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	// sendPayment 실행
	sender.Value = big.NewInt(300)
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// SendPayment Event의 Hash값 확인
	item := <-sink
	hash = item.Hash

	// 받을 Payment를 cancel
	tx, err = contract.CancelPayment(receiver, hash)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// contract, sender, receiver의 balance 변화 없음 확인
	senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	assert.Equal(t, senderPreviousBalance, senderAfterBalance)
	assert.Equal(t, contractPreviousBalance, contractAfterBalance)
	assert.Equal(t, receiverPreviousBalance, receiverAfterBalance)

	// settlePayment를 호출함
	tx, err = contract.SettlePayment(receiver)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// contract, sender, receiver의 balance 변화 없음 확인
	senderAfterBalance, err = backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractAfterBalance, err = backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverAfterBalance, err = backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	assert.Equal(t, senderPreviousBalance, senderAfterBalance)
	assert.Equal(t, contractPreviousBalance, contractAfterBalance)
	assert.Equal(t, receiverPreviousBalance, receiverAfterBalance)
}

func TestPayment_cancelPayment(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()
	var hash [32]byte

	// sendPayment 실행
	sender.Value = big.NewInt(300)
	tx, err := contract.SendPayment(sender, receiver.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// SendPayment Event의 Hash값 확인
	item := <-sink
	hash = item.Hash

	// getPayments를 호출하여 트랜잭션이 리턴되는 것을 확인
	getPayments, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getPayments))
	assert.Equal(t, hash, getPayments[0])

	// 일부의 트랜잭션에 대해 cancelPayment 실행
	tx, err = contract.CancelPayment(receiver, hash)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// getPayments를 호출하여 취소하지 않은 트랜잭션만 리턴되는 것을 확인
	getPayments, err = contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayments))
}

func TestPayment_cancelPayment_balance(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	contractAddress := env.paymentContractAddress
	sendAmount := big.NewInt(300)

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()

	// contract, sender, receiver의 초기 balance 확인
	senderPreviousBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractPreviousBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverPreviousBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	// sendPayment와 cancelPayment 실행
	{
		sender.Value = sendAmount
		tx, err := contract.SendPayment(sender, receiver.From)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// SendPayment Event의 Hash값 확인
		item := <-sink
		hash := item.Hash

		tx, err = contract.CancelPayment(receiver, hash)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// contract, sender, receiver의 balance 변화 확인
	senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	assert.Equal(t, senderPreviousBalance, senderAfterBalance)
	assert.Equal(t, contractPreviousBalance, contractAfterBalance)
	assert.Equal(t, receiverPreviousBalance, receiverAfterBalance)
}

func TestPayment_cancelPayment_empty(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	receiver := env.accounts[0]
	contract := env.paymentContract

	// getPayments를 호출
	getPayments, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayments))

	// cancel 할 Payment가 없는 receiver가 cancelPayment 실행
	var hash [32]byte
	hashDynamicSize := common.MakeRandomBytes(32)
	copy(hash[:], hashDynamicSize)

	tx, err := contract.CancelPayment(receiver, hash)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// getPayments를 호출하여 상태 변화가 없음을 확인
	getPayments, err = contract.GetPayments(&bind.CallOpts{From: receiver.From})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getPayments))
}

func TestPayment_cancelPayment_empty_balance(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	contractAddress := env.paymentContractAddress

	// contract, sender, receiver의 초기 balance 확인
	senderPreviousBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractPreviousBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverPreviousBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	// cancel 할 Payment가 없는 receiver가 cancelPayment 실행
	var hash [32]byte
	hashDynamicSize := common.MakeRandomBytes(32)
	copy(hash[:], hashDynamicSize)

	tx, err := contract.CancelPayment(receiver, hash)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// contract, sender, receiver의 balance 변화 확인
	senderAfterBalance, err := backend.BalanceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	contractAfterBalance, err := backend.BalanceAt(context.Background(), contractAddress, nil)
	assert.NoError(t, err)
	receiverAfterBalance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)

	assert.Equal(t, senderPreviousBalance, senderAfterBalance)
	assert.Equal(t, contractPreviousBalance, contractAfterBalance)
	assert.Equal(t, receiverPreviousBalance, receiverAfterBalance)
}

func TestPayment_settlePayment_cancelPayment(t *testing.T) {
	env := generatePaymentTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.accounts[0]
	receiver := env.accounts[1]
	contract := env.paymentContract
	settleInterval := env.settleInterval

	// subscribe SendPayment event
	var (
		sink     = make(chan *PaymentSendPayment)
		sub, err = contract.WatchSendPayment(nil, sink)
	)
	assert.NoError(t, err)
	defer func() {
		sub.Unsubscribe()
	}()

	sendAmounts := []*big.Int{big.NewInt(300), big.NewInt(400), big.NewInt(1), big.NewInt(500)}
	var hashes [][32]byte

	// settlePayment 4회 실행 [1, 2, 3, 4]
	for _, sendAmount := range sendAmounts {
		sender.Value = sendAmount
		tx, err := contract.SendPayment(sender, receiver.From)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// SendPayment Event의 Hash값 확인
		item := <-sink
		hashes = append(hashes, item.Hash)
	}

	// cancelPayment 로 2번째 payment 취소함 [1, X, 3, 4]
	{
		tx, err := contract.CancelPayment(receiver, hashes[1])
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 트랜잭션 1, 3, 4번이 있는 것을 확인
	{
		getPayment, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
		assert.NoError(t, err)
		assert.Equal(t, 3, len(getPayment))
		assert.Contains(t, getPayment, hashes[0])
		assert.Contains(t, getPayment, hashes[2])
		assert.Contains(t, getPayment, hashes[3])
	}

	// 트랜잭션 settle
	{
		for i := 0; i < settleInterval*10; i++ {
			backend.Commit()
		}

		tx, err := contract.SettlePayment(receiver)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 1, 3, 4번째 payment가 settle되는 것을 확인
	{
		getPayment, err := contract.GetPayments(&bind.CallOpts{From: receiver.From})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(getPayment))
	}
}
