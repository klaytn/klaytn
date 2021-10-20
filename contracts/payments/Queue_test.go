package payments

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

type queueTestENV struct {
	backend              *backends.SimulatedBackend
	queueContract        *Queue
	queueContractAddress common.Address
	owner                *bind.TransactOpts
}

func generateQueueTestEnv(t *testing.T) *queueTestENV {
	// create account
	key, _ := crypto.GenerateKey()
	account := bind.NewKeyedTransactor(key)
	account.GasLimit = DefaultGasLimit

	// generate backend with deployed
	alloc := blockchain.GenesisAlloc{
		account.From: {
			Balance: big.NewInt(params.KLAY),
		},
	}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy
	contractAddress, tx, contract, err := DeployQueue(account, backend)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, bind.CheckWaitMined(backend, tx))

	return &queueTestENV{
		backend:              backend,
		owner:                account,
		queueContract:        contract,
		queueContractAddress: contractAddress,
	}
}

// FIFO로 값이 나오는지 확인
func TestQueue_TestQueue(t *testing.T) {
	env := generateQueueTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	owner := env.owner
	contract := env.queueContract

	// create Hashes
	valNum := 10
	vals := make([][32]byte, valNum)
	for i := 0; i < valNum; i++ {
		copy(vals[i][:], common.MakeRandomBytes(32))
	}

	// enqueue
	for i := 0; i < valNum; i++ {
		tx, err := contract.Enqueue(owner, vals[i])
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// dequeue
	for i := 0; i < valNum; i++ {
		// Peek를 통해 FIFO 순서대로 값이 나오는지 확인
		txHash, err := contract.Peek(&bind.CallOpts{From: owner.From})
		assert.NoError(t, err)
		assert.Equal(t, vals[i], txHash)

		tx, err := contract.Dequeue(owner)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}
}

// queue에 10000개 이상이 들어가는지 확인
func TestQueue_enqueque_size(t *testing.T) {
	env := generateQueueTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	owner := env.owner
	contract := env.queueContract

	// create Hashes
	valNum := 1000
	vals := make([][32]byte, valNum)
	for i := 0; i < valNum; i++ {
		copy(vals[i][:], common.MakeRandomBytes(32))
	}

	// enqueue
	for i := 0; i < valNum; i++ {
		tx, err := contract.Enqueue(owner, vals[i])
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}
}

func TestQueue_dequeue_empty(t *testing.T) {
	env := generateQueueTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	owner := env.owner
	contract := env.queueContract

	tx, err := contract.Dequeue(owner)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
}

// min 이상 max 이히의 값
func getRandomNum(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func TestQueue_getAll(t *testing.T) {
	env := generateQueueTestEnv(t)
	defer env.backend.Close()

	backend := env.backend
	owner := env.owner
	contract := env.queueContract

	var vals [][32]byte

	for i := 0; i < 3; i++ {
		// enqueue
		for j := 0; j < getRandomNum(0, 10); j++ {
			// random Hash 생성
			var val [32]byte
			copy(val[:], common.MakeRandomBytes(32))
			// vals 값 업데이트
			vals = append(vals, val)

			// enqueue
			tx, err := contract.Enqueue(owner, val)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			// 안의 값 확인
			getVals, err := contract.GetAll(&bind.CallOpts{From: owner.From})
			assert.NoError(t, err)
			assert.Equal(t, len(vals), len(getVals))
			assert.Equal(t, vals, getVals)
		}
		// dequeue
		for j := 0; j < getRandomNum(0, len(vals)); j++ {
			// Peek를 통해 FIFO 순서대로 값이 나오는지 확인
			txHash, err := contract.Peek(&bind.CallOpts{From: owner.From})
			assert.NoError(t, err)
			assert.Equal(t, vals[0], txHash)

			// dequeue
			tx, err := contract.Dequeue(owner)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			// vals 값 업데이트
			vals = vals[1:]

			// 안의 값 확인
			getVals, err := contract.GetAll(&bind.CallOpts{From: owner.From})
			assert.NoError(t, err)
			assert.Equal(t, len(vals), len(getVals))
			assert.Equal(t, vals, getVals)
		}
	}
}
