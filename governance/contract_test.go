package governance

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
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
	address, _, contract, err := govcontract.DeployGovParam(owner, sim, owner.From)
	require.Nil(t, err)
	sim.Commit()

	return accounts, sim, address, contract
}

func TestContractEngine_Simulated(t *testing.T) {
	accounts, sim, _, contract := prepareSimulatedContract(t)

	var (
		owner      = accounts[0]
		name       = "istanbul.committeesize"
		value      = []byte{22}
		now        = sim.BlockChain().CurrentBlock().Number().Uint64()
		activation = now + 10
		prevParams = []govcontract.GovParamParamView{{Name: name, Value: []byte{}}}
		nextParams = []govcontract.GovParamParamView{{Name: name, Value: value}}
	)
	contract.SetParam(owner, name, value, activation)
	sim.Commit()

	for {
		sim.Commit() // Advance a block

		latest, err := contract.GetAllParams(nil)
		assert.Nil(t, err)
		pending, err := contract.GetAllParamsAtNextBlock(nil)
		assert.Nil(t, err)

		now = sim.BlockChain().CurrentBlock().Number().Uint64()
		if now < activation-1 {
			assert.Equal(t, prevParams, latest)
			assert.Equal(t, prevParams, pending)
		}
		if now == activation-1 {
			assert.Equal(t, prevParams, latest)
			assert.Equal(t, nextParams, pending)
		}
		if now > activation-1 {
			assert.Equal(t, nextParams, latest)
			assert.Equal(t, nextParams, pending)
			break
		}
	}
}
