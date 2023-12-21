// Copyright 2022 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

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
	"github.com/klaytn/klaytn/storage/database"
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
	config := &params.ChainConfig{}
	config.SetDefaults()
	config.UnitPrice = 25e9
	config.IstanbulCompatibleBlock = common.Big0
	config.LondonCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0

	sim := backends.NewSimulatedBackendWithDatabase(database.NewMemoryDBManager(), alloc, config)

	// Deploy contract
	owner := accounts[0]
	address, tx, contract, err := govcontract.DeployGovParam(owner, sim)
	require.Nil(t, err)
	sim.Commit()

	receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	return accounts, sim, address, contract
}

func prepareSimulatedContractWithParams(t *testing.T, params map[string][]byte) ([]*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *govcontract.GovParam) {
	// Create accounts and simulated blockchain
	accounts, sim, address, contract := prepareSimulatedContract(t)
	owner := accounts[0]

	for name, val := range params {
		tx, err := contract.SetParamIn(owner, name, true, val, big.NewInt(1))
		require.Nil(t, err)
		sim.Commit()

		// check tx success
		receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	ab := new(big.Int).Set(sim.BlockChain().CurrentHeader().Number)
	ab = ab.Add(ab, big.NewInt(1))

	// check with govcontract
	names, values, err := contract.GetAllParamsAt(nil, ab)
	require.Nil(t, err)
	require.Equal(t, len(params), len(names))
	require.Equal(t, len(params), len(values))

	returned := make(map[string][]byte)
	for i := 0; i < len(names); i++ {
		returned[names[i]] = values[i]
	}

	require.Equal(t, params, returned)

	return accounts, sim, address, contract
}

func prepareContractEngine(t *testing.T, bc *blockchain.BlockChain, addr common.Address) *ContractEngine {
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	dbm.WriteGovernance(map[string]interface{}{
		"governance.govparamcontract": addr,
	}, 0)
	gov := NewGovernance(bc.Config(), dbm)
	pset, err := gov.EffectiveParams(0)
	require.Nil(t, err)
	require.Equal(t, addr, pset.GovParamContract())

	gov.SetBlockchain(bc)

	e := NewContractEngine(gov)
	err = e.UpdateParams(bc.CurrentBlock().NumberU64())
	require.Nil(t, err)

	return e
}

// TestContractEngine_Params tests if CurrentParams() returns the parameters required
// for generating the next block. That is,
//
//	start          setparam       activation-1       end
//
// Block |---------------|---------------|---------------|
//
// ..............^               ^               ^
// ..............t0              t1              t2
//
// At num = activation - 2, CurrentParams() = prev
// At num = activation - 1, CurrentParams() = next
//
//	because next is for generating "activation" block
func TestContractEngine_Params(t *testing.T) {
	initialParam := map[string][]byte{
		"istanbul.committeesize": {0xa},
		"governance.unitprice":   {0xb},
	}
	accounts, sim, addr, contract := prepareSimulatedContractWithParams(t, initialParam)

	e := prepareContractEngine(t, sim.BlockChain(), addr)

	var (
		start      = sim.BlockChain().CurrentBlock().NumberU64()
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
			contract.SetParam(owner, key, true, val, new(big.Int).SetUint64(activation))
		}

		var expected *params.GovParamSet

		if num < activation-1 { // t0 & t1
			expected = psetPrev
		} else { // t2
			expected = psetNext
		}

		assert.Equal(t, expected, e.CurrentParams(), "CurrentParams() on block %d failed", num)
		sim.Commit()
		err := e.UpdateParams(sim.BlockChain().CurrentBlock().NumberU64())
		assert.Nil(t, err)
	}
}

// TestContractEngine_ParamsAt tests if EffectiveParams(num) returns the parameters
// required for generating the "num" block. That is,
//
//	start          setparam       activation         end
//
// Block |---------------|---------------|---------------|
//
// ..............^               ^               ^
// ..............t0              t1              t2
//
// EffectiveParams(activation - 1) = prev
// EffectiveParams(activation)     = next
func TestContractEngine_ParamsAt(t *testing.T) {
	initialParam := map[string][]byte{
		"istanbul.committeesize": {0xa},
		"governance.unitprice":   {0xbb, 0xbb, 0xbb, 0xbb},
	}
	accounts, sim, addr, contract := prepareSimulatedContractWithParams(t, initialParam)

	e := prepareContractEngine(t, sim.BlockChain(), addr)

	var (
		start      = sim.BlockChain().CurrentBlock().NumberU64()
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
			contract.SetParam(owner, key, true, val, new(big.Int).SetUint64(activation))
		}

		for iter := start + 1; iter <= num; iter++ {
			var expected *params.GovParamSet

			if iter < activation { // t0 & t1
				expected = psetPrev
			} else { // t2
				expected = psetNext
			}

			result, _ := e.EffectiveParams(iter)
			assert.Equal(t, expected, result, "EffectiveParams(%d) on block %d failed", iter, num)
		}

		sim.Commit()
		err := e.UpdateParams(sim.BlockChain().CurrentBlock().NumberU64())
		assert.Nil(t, err)
	}
}
