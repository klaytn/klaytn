package balanceLimit

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
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

func ValueTransfer(backend *backends.SimulatedBackend, from *bind.TransactOpts, to common.Address, value *big.Int, t *testing.T) (*types.Transaction, error) {
	ctx := context.Background()

	nonce, err := backend.NonceAt(ctx, from.From, nil)
	assert.Equal(t, err, nil)

	chainID, err := backend.ChainID(ctx)
	assert.Equal(t, err, nil)

	gasPrice, err := backend.SuggestGasPrice(ctx)
	assert.Equal(t, err, nil)

	tx := types.NewTransaction(
		nonce,
		to,
		value,
		DefaultGasLimit,
		gasPrice,
		nil)

	signedTx, err := from.Signer(types.NewEIP155Signer(chainID), from.From, tx)
	assert.Equal(t, err, nil)
	err = backend.SendTransaction(ctx, signedTx)

	return signedTx, err
}

type balanceLimitContractTestEnv struct {
	backend             *backends.SimulatedBackend
	sender              *bind.TransactOpts
	receiver            *bind.TransactOpts
	contract            *SendKlay
	contractAddress     common.Address
	initialBalanceLimit *big.Int
}

func generateBalanceLimitContractTestEnv(t *testing.T) *balanceLimitContractTestEnv {
	senderKey, _ := crypto.GenerateKey()
	sender := bind.NewKeyedTransactor(senderKey)
	sender.GasLimit = DefaultGasLimit

	receiverKey, _ := crypto.GenerateKey()
	receiver := bind.NewKeyedTransactor(receiverKey)
	receiver.GasLimit = DefaultGasLimit

	// generate backend with deployed
	alloc := blockchain.GenesisAlloc{
		sender.From: {
			Balance: new(big.Int).Mul(big.NewInt(100), account.GetInitialBalanceLimit()),
		},
	}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy
	contractAddress, tx, contract, err := DeploySendKlay(sender, backend)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	return &balanceLimitContractTestEnv{
		backend:             backend,
		sender:              sender,
		receiver:            receiver,
		contract:            contract,
		contractAddress:     contractAddress,
		initialBalanceLimit: account.GetInitialBalanceLimit(),
	}
}

// contract에 대해 getBalanceLimit을 호출할 때, ErrNotEOA (“Not a EOA”)에러 발생
func TestBalanceLimit_contract_getBalanceLimit(t *testing.T) {
	env := generateBalanceLimitContractTestEnv(t)
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
	env := generateBalanceLimitContractTestEnv(t)
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
	env := generateBalanceLimitContractTestEnv(t)
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
	env := generateBalanceLimitContractTestEnv(t)
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
	env := generateBalanceLimitContractTestEnv(t)
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
	env := generateBalanceLimitContractTestEnv(t)
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
	env := generateBalanceLimitContractTestEnv(t)
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
	env := generateBalanceLimitContractTestEnv(t)
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
