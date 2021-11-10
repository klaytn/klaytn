package sendKlay

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultGasLimit = 5000000
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type sendKlayEOATestEnv struct {
	backend             *backends.SimulatedBackend
	sender              []*bind.TransactOpts
	receiver            []*bind.TransactOpts
	initialBalanceLimit *big.Int
}

func generateSendKlayEOATestEnv(t *testing.T) *sendKlayEOATestEnv {
	senderNum := 1
	receiverNum := 3
	accountNum := senderNum + receiverNum
	keys := make([]*ecdsa.PrivateKey, accountNum)
	accounts := make([]*bind.TransactOpts, accountNum)
	for i := 0; i < accountNum; i++ {
		keys[i], _ = crypto.GenerateKey()
		accounts[i] = bind.NewKeyedTransactor(keys[i])
		accounts[i].GasLimit = DefaultGasLimit
	}

	// generate backend with deployed
	alloc := blockchain.GenesisAlloc{}
	for i := 0; i < senderNum; i++ {
		alloc[accounts[i].From] = blockchain.GenesisAccount{
			Balance: big.NewInt(params.KLAY),
		}
	}
	backend := backends.NewSimulatedBackend(alloc)

	return &sendKlayEOATestEnv{
		backend:             backend,
		sender:              accounts[0:senderNum],
		receiver:            accounts[senderNum:],
		initialBalanceLimit: account.GetInitialBalanceLimit(),
	}
}

type sendKlayContractTestEnv struct {
	backend             *backends.SimulatedBackend
	sender              *bind.TransactOpts
	receiver            *bind.TransactOpts
	contract            *SendKlay
	contractAddress     common.Address
	initialBalanceLimit *big.Int
}

func generateSendKlayContractTestEnv(t *testing.T) *sendKlayContractTestEnv {
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

	return &sendKlayContractTestEnv{
		backend:             backend,
		sender:              sender,
		receiver:            receiver,
		contract:            contract,
		contractAddress:     contractAddress,
		initialBalanceLimit: account.GetInitialBalanceLimit(),
	}
}

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

func setBalanceLimit(backend *backends.SimulatedBackend, from *bind.TransactOpts, balanceLimit *big.Int, t *testing.T) *types.Transaction {
	ctx := context.Background()

	nonce, err := backend.NonceAt(ctx, from.From, nil)
	assert.NoError(t, err)

	chainID, err := backend.ChainID(ctx)
	assert.NoError(t, err)

	tx, err := types.NewTransactionWithMap(types.TxTypeBalanceLimitUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:        nonce,
		types.TxValueKeyGasLimit:     0,
		types.TxValueKeyGasPrice:     big.NewInt(0),
		types.TxValueKeyFrom:         from.From,
		types.TxValueKeyBalanceLimit: balanceLimit,
	})
	assert.NoError(t, err)

	signedTx, err := from.Signer(types.NewEIP155Signer(chainID), from.From, tx)
	assert.NoError(t, err)
	err = backend.SendTransaction(context.Background(), signedTx)
	assert.NoError(t, err)

	return signedTx
}

func setAccountStatus(backend *backends.SimulatedBackend, from *bind.TransactOpts, accountStatus account.AccountStatus, t *testing.T) (*types.Transaction, error) {
	ctx := context.Background()

	nonce, err := backend.NonceAt(ctx, from.From, nil)
	assert.Equal(t, err, nil)

	chainID, err := backend.ChainID(ctx)
	assert.Equal(t, err, nil)

	gasPrice, err := backend.SuggestGasPrice(ctx)
	assert.Equal(t, err, nil)

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          from.From,
		types.TxValueKeyGasLimit:      0,
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyAccountStatus: uint64(accountStatus),
	}

	tx, err := types.NewTransactionWithMap(types.TxTypeAccountStatusUpdate, values)
	if err != nil {
		return nil, err
	}

	signedTx, err := from.Signer(types.NewEIP155Signer(chainID), from.From, tx)
	assert.Equal(t, err, nil)
	err = backend.SendTransaction(ctx, signedTx)

	return signedTx, err
}
