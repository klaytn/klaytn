package governance

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	govcontract "github.com/klaytn/klaytn/contracts/gov"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareSimulatedContract(t *testing.T) ([]*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *govcontract.GovParam) {
	// Create accounts and simulated blockchain
	accounts := []*bind.TransactOpts{}
	alloc := blockchain.GenesisAlloc{}
	for i := 0; i < 1; i++ {
		key, _ := crypto.GenerateKey()
		account := bind.NewKeyedTransactor(key)
		accounts = append(accounts, account)
		alloc[account.From] = blockchain.GenesisAccount{Balance: big.NewInt(params.KLAY)}
	}
	sim := backends.NewSimulatedBackend(alloc)

	// Deploy contract
	owner := accounts[0]
	address, _, contract, err := govcontract.DeployGovParam(owner, sim)
	require.Nil(t, err)
	sim.Commit()

	tx, err := contract.Initialize(owner, owner.From)
	require.Nil(t, err)
	sim.Commit()

	receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	return accounts, sim, address, contract
}

func TestContractEngine_Simulated(t *testing.T) {
	accounts, sim, _, contract := prepareSimulatedContract(t)

	var (
		owner  = accounts[0]
		name   = "istanbul.committeesize"
		valueA = []byte{0xa}
		valueB = []byte{0xb}
	)

	// Empty array before SetParam()
	names, values, err := contract.GetAllParams(nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(names))
	assert.Equal(t, 0, len(values))

	// Call SetParam()
	tx, err := contract.SetParamByOwner(owner, name, valueA)
	require.Nil(t, err)
	sim.Commit()

	receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Value exists after SetParam()
	names, values, err = contract.GetAllParams(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, 1, len(values))
	assert.Equal(t, name, names[0])
	assert.Equal(t, valueA, values[0].Value)

	// Call SetParam() again
	tx, err = contract.SetParamByOwner(owner, name, valueB)
	require.Nil(t, err)
	sim.Commit()

	receipt, _ = sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Value changed after SetParam()
	names, values, err = contract.GetAllParams(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, 1, len(values))
	assert.Equal(t, name, names[0])
	assert.Equal(t, valueB, values[0].Value)
}
