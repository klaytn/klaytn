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
		account.GasLimit = 10000000
		accounts = append(accounts, account)
		alloc[account.From] = blockchain.GenesisAccount{Balance: big.NewInt(params.KLAY)}
	}
	sim := backends.NewSimulatedBackend(alloc)

	// Deploy contract
	owner := accounts[0]
	address, _, contract, err := govcontract.DeployGovParam(owner, sim)
	require.Nil(t, err)
	sim.Commit()

	return accounts, sim, address, contract
}

func prepareSimulatedContractWithParams(t *testing.T, p map[string][]byte) ([]*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *govcontract.GovParam) {
	accounts, sim, address, contract := prepareSimulatedContract(t)

	// Empty array before addParam()
	names, values, err := contract.GetAllParams(nil)
	require.Nil(t, err)
	require.Equal(t, 0, len(names))
	require.Equal(t, 0, len(values))

	// Call addParam()
	owner := accounts[0]
	txs := make(types.Transactions, 0)
	for k, v := range p {
		tx, err := contract.AddParam(owner, k, v)
		require.Nil(t, err)
		txs = append(txs, tx)
	}

	sim.Commit() // mine addParam
	sim.Commit() // mine activation block so getParam returns added param

	for _, tx := range txs {
		receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	// Value exists after addParam
	for k, v := range p {
		value, err := contract.GetParam(nil, k)
		require.Nil(t, err)
		require.Equal(t, v, value)
	}

	return accounts, sim, address, contract
}

func TestContractEngine_ContractConnector(t *testing.T) {
	var (
		name   = "istanbul.committeesize"
		valueA = []byte{0xa}
		valueB = []byte{0xbb, 0xbb}
		p      = map[string][]byte{
			name: valueA,
		}
	)

	accounts, sim, _, contract := prepareSimulatedContractWithParams(t, p)

	owner := accounts[0]

	// Value exists after SetParam()
	names, values, err := contract.GetAllParams(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, 1, len(values))
	assert.Equal(t, name, names[0])
	assert.Equal(t, valueA, values[0])

	// Call SetParam() again
	ab := sim.BlockChain().CurrentHeader().Number.Uint64() + 2
	tx, err := contract.SetParam(owner, name, valueB, ab)
	require.Nil(t, err)

	// increase block number to reach activation block
	for sim.BlockChain().CurrentHeader().Number.Uint64() < ab {
		sim.Commit()
	}

	receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Value changed after SetParam()
	names, values, err = contract.GetAllParams(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, 1, len(values))
	assert.Equal(t, name, names[0])
	assert.Equal(t, valueB, values[0])
}

func prepareContractEngine(t *testing.T, bc *blockchain.BlockChain, addr common.Address) *ContractEngine {
	config := params.CypressChainConfig.Copy()
	config.Governance.GovParamContract = addr

	e := NewContractEngine(config)
	e.SetBlockchain(bc)
	err := e.UpdateParams()
	require.Nil(t, err)

	return e
}

func TestContractEngine_Params(t *testing.T) {
	initialParam := map[string][]byte{
		"istanbul.committeesize": {0xa},
		"governance.unitprice":   {0xb},
	}
	accounts, sim, addr, contract := prepareSimulatedContractWithParams(t, initialParam)
	e := prepareContractEngine(t, sim.BlockChain(), addr)

	//     start          setparam       activation         end
	// Block |---------------|---------------|---------------|
	//               ^               ^               ^
	//               t0              t1              t2
	// At num = activation - 1, Params() = prev
	// At num = activation, Params() = next
	var (
		start      = sim.BlockChain().CurrentHeader().Number.Uint64()
		setparam   = start + 5
		activation = setparam + 5
		end        = activation + 5
		key        = "governance.unitprice"
		val        = []byte{0xff, 0xff, 0xff, 0xff}
		update, _  = params.NewGovParamSetBytesMap(map[string][]byte{
			key: val,
		})
		psetPrev, _ = params.NewGovParamSetBytesMap(initialParam)   // for t0 & t1
		psetNext    = params.NewGovParamSetMerged(psetPrev, update) // for t2
		owner       = accounts[0]
	)

	for num := start; num < end; num++ {
		if num == setparam { // setParam
			contract.SetParam(owner, key, val, activation)
		}

		var expected *params.GovParamSet

		if num < activation { // t0 & t1
			expected = psetPrev
		} else { // t2
			expected = psetNext
		}

		require.Equal(t, expected, e.Params())
		sim.Commit()
		err := e.UpdateParams()
		require.Nil(t, err)
	}
}

func TestContractEngine_ParamsAt(t *testing.T) {
	initialParam := map[string][]byte{
		"istanbul.committeesize": {0xa},
		"governance.unitprice":   {0xbb, 0xbb, 0xbb, 0xbb},
	}
	accounts, sim, addr, contract := prepareSimulatedContractWithParams(t, initialParam)
	e := prepareContractEngine(t, sim.BlockChain(), addr)

	//     start          setparam       activation         end
	// Block |---------------|---------------|---------------|
	//               ^               ^               ^
	//               t0              t1              t2
	// ParamsAt(activation) = prev
	// ParamsAt(activation + 1) = next
	var (
		start      = sim.BlockChain().CurrentHeader().Number.Uint64()
		setparam   = start + 5
		activation = setparam + 5
		end        = activation + 5
		key        = "governance.unitprice"
		val        = []byte{0xff, 0xff, 0xff, 0xff}
		update, _  = params.NewGovParamSetBytesMap(map[string][]byte{
			key: val,
		})
		psetPrev, _ = params.NewGovParamSetBytesMap(initialParam)   // for t0 & t1
		psetNext    = params.NewGovParamSetMerged(psetPrev, update) // for t2
		owner       = accounts[0]
	)

	for num := start; num < end; num++ {
		if num == setparam { // setParam
			contract.SetParam(owner, key, val, activation)
		}

		for iter := start + 1; iter <= num; iter++ {
			var expected *params.GovParamSet

			if iter < activation { // t0 & t1
				expected = psetPrev
			} else { // t2
				expected = psetNext
			}

			result, _ := e.ParamsAt(iter + 1)
			require.Equal(t, expected, result)
		}

		sim.Commit()
		err := e.UpdateParams()
		require.Nil(t, err)
	}
}
