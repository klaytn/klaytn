package sendKlay

import (
	"context"
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

// 존재하지 않는 account에 대해 getBalanceLimit을 호출할 때, ErrNilAccount(Account not set) 에러 발생
func TestBalanceLimit_EOA_getBalanceLimit_newAccount(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	receiver := env.receiver[0]

	// getBalanceLimit 호출 시 실패
	limit, err := backend.BalanceLimitAt(context.Background(), receiver.From, nil)
	assert.Equal(t, account.ErrNilAccount, err)
	assert.Equal(t, big.NewInt(0), limit)
}

// 존재하지 않는 account에 대해 setBalanceLimit을 호출할 때, BalanceLimit 설정하는 것을 확인
func TestBalanceLimit_EOA_setBalanceLimit_newAccount(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	balanceLimit := big.NewInt(100)

	// setBalanceLimit 호출
	tx := setBalanceLimit(backend, sender, balanceLimit, math.MaxUint64, t)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	//klay_BalanceLimit 실행 시 설정한 값이 셋팅된 것을 확인
	limit, err := backend.BalanceLimitAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	assert.Equal(t, balanceLimit, limit)
}

// 존재하는 account에 대해 getBalanceLimit을 호출할 때, 초기값 출력 InitialBalanceLimit 확인
func TestBalanceLimit_EOA_InitialBalanceLimit(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	receiver := env.receiver[0]

	// A → B로 value transfer (B를 EOA로 설정)
	min := 0
	max := math.MaxInt32
	transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
	tx, err := ValueTransfer(backend, sender, receiver.From, transferAmount, math.MaxUint64, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// B에 대해 getBalanceLimit 실행하여 InitialBalanceLimit 조회
	limit, err := backend.BalanceLimitAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	assert.Equal(t, account.GetInitialBalanceLimit(), limit)
}

// txError로 ReceiptStatusErrExceedBalanceLimit가 발생하는지 확인
func TestBalanceLimit_EOA_ReceiptStatus(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	receiver := env.receiver[0]

	balanceLimit := big.NewInt(100)
	transferAmount := big.NewInt(200)

	tx := setBalanceLimit(backend, receiver, balanceLimit, math.MaxUint64, t)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err := ValueTransfer(backend, sender, receiver.From, transferAmount, math.MaxUint64, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExceedBalanceLimit, t)
}

// 여러 개의 트랜잭션을 한 블록에 커밋했을 때, 일부 트랜잭션만 실패하는지 확인하기 위함
func TestBalanceLimit_EOA_pending_valueTransfer(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	receiver := env.receiver[0]

	var successTxs []*types.Transaction
	var failTxs []*types.Transaction

	senderNonce := uint64(0)
	receiverNonce := uint64(0)

	tryNum := 100

	// setBalanceLimit 호출
	balanceLimit := new(big.Int).Mul(big.NewInt(int64(tryNum)), big.NewInt(10))
	tx := setBalanceLimit(backend, receiver, balanceLimit, receiverNonce, t)
	successTxs = append(successTxs, tx)
	receiverNonce++

	// 다음 트랜젝션 100쌍 생성
	for i := 0; i < tryNum; i++ {
		// 성공하는 value transfer
		transferAmountSuccess := big.NewInt(10)
		txSuccess, err := ValueTransfer(backend, sender, receiver.From, transferAmountSuccess, senderNonce, t)
		successTxs = append(successTxs, txSuccess)
		assert.NoError(t, err)
		senderNonce++

		// 실패하는 value transfer
		transferAmountFail := new(big.Int).Mul(balanceLimit, big.NewInt(2))
		txFail, err := ValueTransfer(backend, sender, receiver.From, transferAmountFail, senderNonce, t)
		failTxs = append(failTxs, txFail)
		assert.NoError(t, err)
		senderNonce++
	}

	// 블록 생성 (모든 tx 하나의 블록에)
	backend.Commit()

	// 트랜잭션 성공/실패값 확인
	for _, tx := range successTxs {
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}
	for _, tx := range failTxs {
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExceedBalanceLimit, t)
	}
}

// 한 블록에 setBalanceLimit과 ValueTransfer를 여러번 호출했을 때 BalanceLimit이 안넘는 tx만 반영된 것을 확인
func TestBalanceLimit_pending_setBalanceLimit_valueTransfer(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	receiver := env.receiver[0]

	var successTxs []*types.Transaction
	var failTxs []*types.Transaction

	senderNonce := uint64(0)
	receiverNonce := uint64(0)

	// 한 account에 대하여 다음 tx를 100개 생성
	tryNum := 100
	balanceLimit := 0
	transferMax := params.KLAY / tryNum
	transferSum := 0

	for i := 0; i <= tryNum; i++ {
		transferAmountSuccess := rand.Intn(transferMax)
		transferAmountFail := rand.Intn(transferMax)
		balanceLimit += transferAmountSuccess
		transferSum += transferAmountSuccess

		// BalanceLimitUpdate tx
		tx := setBalanceLimit(backend, receiver, big.NewInt(int64(balanceLimit)), receiverNonce, t)
		successTxs = append(successTxs, tx)
		receiverNonce++

		// 성공하는 value transfer
		txSuccess, err := ValueTransfer(backend, sender, receiver.From, big.NewInt(int64(transferAmountSuccess)), senderNonce, t)
		assert.NoError(t, err)
		successTxs = append(successTxs, txSuccess)
		senderNonce++

		// 실패하는 value transfer
		txFail, err := ValueTransfer(backend, sender, receiver.From, big.NewInt(int64(transferAmountFail)), senderNonce, t)
		assert.NoError(t, err)
		failTxs = append(failTxs, txFail)
		senderNonce++
	}

	// 블록 생성
	backend.Commit()

	// 트랜잭션 성공/실패값 확인
	for _, tx := range successTxs {
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}
	for _, tx := range failTxs {
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExceedBalanceLimit, t)
	}

	// 최종 balanceLimit, balance, Nonce 값 확인
	limit, err := backend.BalanceLimitAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(int64(balanceLimit)), limit)

	balance, err := backend.BalanceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(int64(transferSum)), balance)

	getSenderNonce, err := backend.NonceAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	assert.Equal(t, senderNonce, getSenderNonce)

	getReceiverNonce, err := backend.NonceAt(context.Background(), receiver.From, nil)
	assert.NoError(t, err)
	assert.Equal(t, receiverNonce, getReceiverNonce)
}
