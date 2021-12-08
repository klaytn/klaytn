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
	"github.com/stretchr/testify/assert"
)

// 존재하지 않는 account에 대해 getAccountStatus을 호출할 때, ErrNilAccount(Account not set) 에러 발생
func TestAccountStatus_getAccountStatus_newAccount(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	receiver := env.receiver[0]

	// getBalanceLimit 호출 시 실패
	status, err := backend.AccountStatusAt(context.Background(), receiver.From, nil)
	assert.Equal(t, account.ErrNilAccount, err)
	assert.Equal(t, account.AccountStatusUndefined, status)
}

// 존재하지 않는 account에 대해 setBalanceLimit을 호출할 때, BalanceLimit 설정하는 것을 확인
func TestAccountStatus_setAccountStatus_newAccount(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	setStatus := account.AccountStatusStop

	// setBalanceLimit 호출
	tx, err := setAccountStatus(backend, sender, setStatus, math.MaxUint64, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	//klay_BalanceLimit 실행 시 설정한 값이 셋팅된 것을 확인
	status, err := backend.AccountStatusAt(context.Background(), sender.From, nil)
	assert.NoError(t, err)
	assert.Equal(t, setStatus, status)
}

// txError로 ReceiptStatusErrStoppedAccountFrom 발생하는지 확인
func TestAccountStatus_EOA_ReceiptStatus_From(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	receiver := env.receiver[0]

	tx, err := setAccountStatus(backend, sender, account.AccountStatusStop, math.MaxUint64, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = ValueTransfer(backend, sender, receiver.From, big.NewInt(0), math.MaxUint64, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrStoppedAccountFrom, t)
}

// txError로 ReceiptStatusErrStoppedAccountTo 발생하는지 확인
func TestAccountStatus_EOA_ReceiptStatus_To(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]
	receiver := env.receiver[0]

	tx, err := setAccountStatus(backend, receiver, account.AccountStatusStop, math.MaxUint64, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = ValueTransfer(backend, sender, receiver.From, big.NewInt(0), math.MaxUint64, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrStoppedAccountTo, t)
}

// 한 블록에 setAccountStatus과 ValueTransfer를 반복적으로 호출했을 때 AccountStatus가 Active일 때 tx가 성공한 것을 확인
func TestAccountStatus_pending_setAccountStatus_valueTransfer(t *testing.T) {
	env := generateSendKlayEOATestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender[0]

	var successTxs []*types.Transaction
	var failTxs []*types.Transaction

	senderNonce := uint64(0)

	tryNum := 100
	transferMax := math.MaxInt32

	for i := 0; i <= tryNum; i++ {
		transferAmount := rand.Intn(transferMax)

		// AccountStatusUpdate tx (Active)
		txAccountStatusActive, err := setAccountStatus(backend, sender, account.AccountStatusActive, senderNonce, t)
		assert.NoError(t, err)
		successTxs = append(successTxs, txAccountStatusActive)
		senderNonce++

		// value transfer tx (성공)
		txSuccess, err := ValueTransfer(backend, sender, sender.From, big.NewInt(int64(transferAmount)), senderNonce, t)
		assert.NoError(t, err)
		successTxs = append(successTxs, txSuccess)
		senderNonce++

		// AccountStatusUpdate tx (Stop)
		txAccountStatusStop, err := setAccountStatus(backend, sender, account.AccountStatusStop, senderNonce, t)
		assert.NoError(t, err)
		successTxs = append(successTxs, txAccountStatusStop)
		senderNonce++

		// value transfer tx (실패)
		txFail, err := ValueTransfer(backend, sender, sender.From, big.NewInt(int64(transferAmount)), senderNonce, t)
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
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrStoppedAccountFrom, t)
	}
}
