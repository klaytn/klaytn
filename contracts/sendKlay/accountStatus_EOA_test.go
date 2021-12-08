package sendKlay

import (
	"context"
	"math"
	"math/big"
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
