package sendKlay

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/stretchr/testify/assert"
)

// contract에 대해 getAccountStatus을 호출할 때, ErrNotEOA (“Not a EOA”)에러 발생
func TestAccountStatus_contract_getAccountStatus(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	contractAddress := env.contractAddress

	// contract에 대해 getAccountStatus을 호출할 때, ErrNotEOA (“Not a EOA”)에러 발생
	status, err := backend.AccountStatusAt(context.Background(), contractAddress, nil)
	assert.Equal(t, account.ErrNotEOA, err)
	assert.Equal(t, account.AccountStatusUndefined, status)
}

// contract → EOA로 cbdc 전송 시 receiver의 balanceLimit을 지키는 것을 확인
func TestAccountStatus_contract_transfer(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	receiver := env.receiver
	contract := env.contract

	// receiver의 AccountStatus를 Stop로 설정함
	tx, err := setAccountStatus(backend, receiver, account.AccountStatusStop, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// contract에 충분한 CBDC 전송
	sender.Value = new(big.Int).Mul(big.NewInt(100), big.NewInt(10))
	tx, err = contract.ContractPayable(sender)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	sender.Value = big.NewInt(0)

	// 해당 contract를 이용하여 receiver에게 value transfer 시 실패하는 것을 확인
	sendValue := big.NewInt(100)
	tx, err = contract.ContractTransfer(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
}

// contract → EOA로 cbdc 전송 시 receiver의 balanceLimit을 지키는 것을 확인
func TestAccountStatus_contract_send(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	receiver := env.receiver
	contract := env.contract

	// receiver의 AccountStatus를 Stop로 설정함
	tx, err := setAccountStatus(backend, receiver, account.AccountStatusStop, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// contract에 충분한 CBDC 전송
	sender.Value = new(big.Int).Mul(big.NewInt(100), big.NewInt(10))
	tx, err = contract.ContractPayable(sender)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	sender.Value = big.NewInt(0)

	// 해당 contract를 이용하여 receiver에게 value transfer 시 실패하는 것을 확인
	sendValue := big.NewInt(100)
	tx, err = contract.ContractSend(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
}

// contract → EOA로 cbdc 전송 시 receiver의 balanceLimit을 지키는 것을 확인
func TestAccountStatus_contract_call(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	receiver := env.receiver
	contract := env.contract

	// receiver의 AccountStatus를 Stop로 설정함
	tx, err := setAccountStatus(backend, receiver, account.AccountStatusStop, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// contract에 충분한 CBDC 전송
	sender.Value = new(big.Int).Mul(big.NewInt(100), big.NewInt(10))
	tx, err = contract.ContractPayable(sender)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	sender.Value = big.NewInt(0)

	// 해당 contract를 이용하여 receiver에게 value transfer 시 실패하는 것을 확인
	sendValue := big.NewInt(100)
	tx, err = contract.ContractCall(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
}
