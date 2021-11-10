package sendKlay

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

// contract에 대해 getBalanceLimit을 호출할 때, ErrNotEOA (“Not a EOA”)에러 발생
func TestBalanceLimit_contract_getBalanceLimit(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	contractAddress := env.contractAddress

	// contract에 대해 getBalanceLimit 호출하여 에러 확인
	limit, err := backend.BalanceLimitAt(context.Background(), contractAddress, nil)
	assert.Equal(t, account.ErrNotEOA, err)
	assert.Equal(t, big.NewInt(0), limit)
}

// contract를 생성할 때 balanceLimit의 제한이 없는 것을 확인
func TestBalanceLimit_contract_create(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	initialBalanceLimit := env.initialBalanceLimit

	// contract deploy 시 InitialBalanceLimit보다 큰 값을 전송
	sender.Value = new(big.Int).Mul(initialBalanceLimit, big.NewInt(10))
	_, tx, _, err := DeploySendKlay(sender, backend)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))
}

// EOA → contract로 cbdc를 전송할 때, 제한이 없는 것을 확인
func TestBalanceLimit_contract_payable(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	//contract := env.contract
	contractAddress := env.contractAddress
	initialBalanceLimit := env.initialBalanceLimit

	// contract에 InitialBalanceLimit보다 큰 값을 전송
	transferValue := new(big.Int).Mul(initialBalanceLimit, big.NewInt(10))
	tx, err := ValueTransfer(backend, sender, contractAddress, transferValue, t)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

// EOA → contract로 cbdc를 전송할 때, 제한이 없는 것을 확인
func TestBalanceLimit_contract_receive(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	contract := env.contract
	initialBalanceLimit := env.initialBalanceLimit

	//contract에 InitialBalanceLimit보다 큰 값을 전송
	sender.Value = new(big.Int).Mul(initialBalanceLimit, big.NewInt(10))
	tx, err := contract.Receive(sender)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

// EOA → contract로 cbdc를 전송할 때, 제한이 없는 것을 확인
func TestBalanceLimit_contract_fallback(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	contract := env.contract
	initialBalanceLimit := env.initialBalanceLimit

	//contract에 InitialBalanceLimit보다 큰 값을 전송
	sender.Value = new(big.Int).Mul(initialBalanceLimit, big.NewInt(10))
	tx, err := contract.Fallback(sender, common.MakeRandomBytes(10))
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

// contract → EOA로 cbdc 전송 시 receiver의 balanceLimit을 지키는 것을 확인
func TestBalanceLimit_contract_transfer(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	receiver := env.receiver
	contract := env.contract
	initialBalanceLimit := env.initialBalanceLimit

	// contract에 충분한 CBDC 전송
	sender.Value = new(big.Int).Mul(initialBalanceLimit, big.NewInt(10))
	tx, err := contract.ContractPayable(sender)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	sender.Value = big.NewInt(0)

	// 해당 contract를 이용하여 InitialBalanceLimit보다 큰 cbdc 전송시 에러 발생하는 것을 확인
	sendValue := new(big.Int).Mul(initialBalanceLimit, big.NewInt(2))
	tx, err = contract.ContractTransfer(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// 해당 contract를 이용하여 InitialBalanceLimit보다 작은 cbdc 전송시 성공하는 것을 확인
	sender.Value = big.NewInt(0)
	sendValue = big.NewInt(1)
	tx, err = contract.ContractTransfer(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

// contract → EOA로 cbdc 전송 시 receiver의 balanceLimit을 지키는 것을 확인
func TestBalanceLimit_contract_send(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	receiver := env.receiver
	contract := env.contract
	initialBalanceLimit := env.initialBalanceLimit

	// contract에 충분한 CBDC 전송
	sender.Value = new(big.Int).Mul(initialBalanceLimit, big.NewInt(10))
	tx, err := contract.ContractPayable(sender)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	sender.Value = big.NewInt(0)

	// 해당 contract를 이용하여 InitialBalanceLimit보다 큰 cbdc 전송시 에러 발생하는 것을 확인
	sendValue := new(big.Int).Mul(initialBalanceLimit, big.NewInt(2))
	tx, err = contract.ContractSend(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// 해당 contract를 이용하여 InitialBalanceLimit보다 작은 cbdc 전송시 성공하는 것을 확인
	sender.Value = big.NewInt(0)
	sendValue = big.NewInt(1)
	tx, err = contract.ContractSend(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

// contract → EOA로 cbdc 전송 시 receiver의 balanceLimit을 지키는 것을 확인
func TestBalanceLimit_contract_call(t *testing.T) {
	env := generateSendKlayContractTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	sender := env.sender
	receiver := env.receiver
	contract := env.contract
	initialBalanceLimit := env.initialBalanceLimit

	// contract에 충분한 CBDC 전송
	sender.Value = new(big.Int).Mul(initialBalanceLimit, big.NewInt(10))
	tx, err := contract.ContractPayable(sender)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	sender.Value = big.NewInt(0)

	// 해당 contract를 이용하여 InitialBalanceLimit보다 큰 cbdc 전송시 에러 발생하는 것을 확인
	sendValue := new(big.Int).Mul(initialBalanceLimit, big.NewInt(2))
	tx, err = contract.ContractCall(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	// 해당 contract를 이용하여 InitialBalanceLimit보다 작은 cbdc 전송시 성공하는 것을 확인
	sender.Value = big.NewInt(0)
	sendValue = big.NewInt(1)
	tx, err = contract.ContractCall(sender, receiver.From, sendValue)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}
