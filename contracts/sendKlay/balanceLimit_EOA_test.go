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
	tx := setBalanceLimit(backend, sender, balanceLimit, t)
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
	tx, err := ValueTransfer(backend, sender, receiver.From, transferAmount, t)
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

	tx := setBalanceLimit(backend, receiver, balanceLimit, t)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err := ValueTransfer(backend, sender, receiver.From, transferAmount, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExceedBalanceLimit, t)
}
