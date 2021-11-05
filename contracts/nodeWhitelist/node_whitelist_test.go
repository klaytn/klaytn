package nodewhitelist

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultGasLimit = 5000000
	AccountNum      = 5
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CheckReceipt can check if the tx receipt has expected status.
func CheckReceipt(b bind.DeployBackend, tx *types.Transaction, duration time.Duration, expectedStatus uint, t *testing.T) {
	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), duration)
	defer cancelTimeout()

	receipt, err := bind.WaitMined(timeoutContext, b, tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedStatus, receipt.Status)
}

type nodeWhitelistEnv struct {
	backend               *backends.SimulatedBackend
	admin                 *bind.TransactOpts
	accounts              []*bind.TransactOpts
	nodeWhitelistContract *NodeWhitelist
}

func generateMintBurnTestEnv(t *testing.T, accountNum uint) *nodeWhitelistEnv {
	keys := make([]*ecdsa.PrivateKey, accountNum+4)
	accounts := make([]*bind.TransactOpts, accountNum+4)
	for i := uint(0); i < accountNum+3; i++ {
		keys[i], _ = crypto.GenerateKey()
		accounts[i] = bind.NewKeyedTransactor(keys[i])
		accounts[i].GasLimit = DefaultGasLimit
	}

	// generate backend with deployed
	alloc := blockchain.GenesisAlloc{
		blockchain.NodeWhitelistContractAddr: {
			Code:    common.FromHex(NodeWhitelistBinRuntime),
			Balance: big.NewInt(0),
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"): common.BytesToHash(accounts[0].From.Bytes()),
			},
		},
	}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy MintBurn contract
	mintBurnContract, err := NewNodeWhitelist(blockchain.NodeWhitelistContractAddr, backend)
	assert.NoError(t, err)

	return &nodeWhitelistEnv{
		backend:               backend,
		admin:                 accounts[0],
		accounts:              accounts[3:],
		nodeWhitelistContract: mintBurnContract,
	}
}

// TestNodeWhitelist_GetAdmin checks if GetAdmin returns the stored admin address.
func TestNodeWhitelist_GetAdmin(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	admin := env.admin
	contract := env.nodeWhitelistContract

	adminFromContract, err := contract.GetAdmin(nil)
	assert.NoError(t, err)
	assert.Equal(t, admin.From, adminFromContract)
}

// TestNodeWhitelist_SetAdmin_ByAdmin checks if SetAdmin is called successfully when called by admin.
func TestNodeWhitelist_SetAdmin_ByAdmin(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	admin := env.admin
	newAdmin := env.accounts[0]
	contract := env.nodeWhitelistContract

	tx, err := contract.SetAdmin(admin, newAdmin.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	adminFromContract, err := contract.GetAdmin(nil)
	assert.NoError(t, err)
	assert.Equal(t, newAdmin.From, adminFromContract)
}

// TestNodeWhitelist_SetAdmin_ByAdmin checks if SetAdmin is not called successfully when called by non-admin.
func TestNodeWhitelist_SetAdmin_ByNonAdmin(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	admin := env.admin
	newAdmin := env.accounts[0]
	contract := env.nodeWhitelistContract

	tx, err := contract.SetAdmin(newAdmin, newAdmin.From)
	assert.NoError(t, err)
	backend.Commit()
	CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

	adminFromContract, err := contract.GetAdmin(nil)
	assert.NoError(t, err)
	assert.Equal(t, admin.From, adminFromContract)
}

// TestNodeWhitelist_AddNode_ByAdmin checks if AddNode is called successfully when called by admin.
func TestNodeWhitelist_AddNode_ByAdmin(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	admin := env.admin
	contract := env.nodeWhitelistContract

	adminFromContract, err := contract.GetAdmin(nil)
	assert.NoError(t, err)
	assert.Equal(t, admin.From, adminFromContract)

	// randomly generates 50 nodes and add them to the contract by "admin"
	nodes := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		randomNode := fmt.Sprint(rand.Int63())
		nodes[randomNode] = struct{}{}

		tx, err := contract.AddNode(admin, randomNode)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// checks if the added nodes are returned
	whitelist, err := contract.GetWhitelist(nil)
	assert.NoError(t, err)
	assert.Equal(t, len(nodes), len(whitelist))

	for _, nodeInWhitelist := range whitelist {
		_, exist := nodes[nodeInWhitelist]
		assert.True(t, exist)
	}

	// check if the BlockChain.GetWhitelist returns the same result
	assert.Equal(t, whitelist, backend.BlockChain().GetNodeWhitelist())
}

// TestNodeWhitelist_AddNode_ByNonAdmin checks if AddNode is not called successfully when called by non-admin.
func TestNodeWhitelist_AddNode_ByNonAdmin(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	nonAdmin := env.accounts[0]
	contract := env.nodeWhitelistContract

	// randomly generates 50 nodes and add them to the contract by "non-admin"
	nodes := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		randomNode := fmt.Sprint(rand.Int63())
		nodes[randomNode] = struct{}{}

		tx, err := contract.AddNode(nonAdmin, randomNode)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// checks if the added nodes are returned
	whitelist, err := contract.GetWhitelist(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(whitelist))

	// check if the BlockChain.GetWhitelist returns the same result
	assert.Equal(t, whitelist, backend.BlockChain().GetNodeWhitelist())
}

// TestNodeWhitelist_DelNode_ByAdmin checks if DelNode is called successfully when called by admin.
func TestNodeWhitelist_DelNode_ByAdmin(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	admin := env.admin
	contract := env.nodeWhitelistContract

	adminFromContract, err := contract.GetAdmin(nil)
	assert.NoError(t, err)
	assert.Equal(t, admin.From, adminFromContract)

	// randomly generates 50 nodes and add them to the contract by "admin"
	nodes := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		randomNode := fmt.Sprint(rand.Int63())
		nodes[randomNode] = struct{}{}

		tx, err := contract.AddNode(admin, randomNode)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// delete all the added nodes by "admin"
	for node, _ := range nodes {
		tx, err := contract.DelNode(admin, node)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// checks if the deleted nodes are returned
	whitelist, err := contract.GetWhitelist(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(whitelist))

	// check if the BlockChain.GetWhitelist returns the same result
	assert.Equal(t, whitelist, backend.BlockChain().GetNodeWhitelist())
}

// TestNodeWhitelist_DelNode_ByNonAdmin checks if DelNode is not called successfully when called by non-admin.
func TestNodeWhitelist_DelNode_ByNonAdmin(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	admin := env.admin
	nonAdmin := env.accounts[0]
	contract := env.nodeWhitelistContract

	adminFromContract, err := contract.GetAdmin(nil)
	assert.NoError(t, err)
	assert.Equal(t, admin.From, adminFromContract)

	// randomly generates 50 nodes and add them to the contract by "admin"
	nodes := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		randomNode := fmt.Sprint(rand.Int63())
		nodes[randomNode] = struct{}{}

		tx, err := contract.AddNode(admin, randomNode)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// delete all the added nodes by "non-admin"
	for node, _ := range nodes {
		tx, err := contract.DelNode(nonAdmin, node)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// checks if the deleted nodes are returned
	whitelist, err := contract.GetWhitelist(nil)
	assert.NoError(t, err)
	assert.Equal(t, len(nodes), len(whitelist))

	for _, nodeInWhitelist := range whitelist {
		_, exist := nodes[nodeInWhitelist]
		assert.True(t, exist)
	}

	// check if the BlockChain.GetWhitelist returns the same result
	assert.Equal(t, whitelist, backend.BlockChain().GetNodeWhitelist())
}
