package system

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func TestReadRegistry(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KLAY),
			},
			RegistryAddr: {
				Code:    RegistryMockCode,
				Balance: common.Big0,
			},
		}
		backend = backends.NewSimulatedBackend(alloc)

		recordName = "AcmeContract"
		recordAddr = common.HexToAddress("0xaaaa")
	)

	// Without a record
	addr, err := ReadRegistryActiveAddr(backend, recordName, common.Big0)
	assert.Nil(t, err)
	assert.Equal(t, common.Address{}, addr)

	// Register a record
	contract, err := contracts.NewRegistryMockTransactor(RegistryAddr, backend)
	_, err = contract.Register(sender, recordName, recordAddr, common.Big1)
	assert.Nil(t, err)
	backend.Commit()

	// With the record
	addr, err = ReadRegistryActiveAddr(backend, recordName, common.Big1)
	assert.Nil(t, err)
	assert.Equal(t, recordAddr, addr)
}

// Test that AllocRegistry correctly reproduces the storage state
// identical to the state after a series of register() call.
func TestAllocRegistry(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	// 1. Create storage with AllocRegistry
	allocStorage := AllocRegistry(&AllocRegistryInit{
		Contracts: map[string]common.Address{
			"AcmeContract": common.HexToAddress("0xaaaa"),
			"TestContract": common.HexToAddress("0xcccc"),
		},
		Owner: common.HexToAddress("0xffff"),
	})

	// 2. Create storage by calling register()
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KLAY),
			},
			RegistryAddr: {
				Code:    RegistryMockCode,
				Balance: common.Big0,
			},
		}
		backend     = backends.NewSimulatedBackend(alloc)
		contract, _ = contracts.NewRegistryMockTransactor(RegistryAddr, backend)
	)

	contract.Register(sender, "AcmeContract", common.HexToAddress("0xaaaa"), common.Big0)
	contract.Register(sender, "TestContract", common.HexToAddress("0xcccc"), common.Big0)
	contract.TransferOwnership(sender, common.HexToAddress("0xffff"))
	backend.Commit()

	execStorage := make(map[common.Hash]common.Hash)
	stateDB, _ := backend.BlockChain().State()
	stateDB.ForEachStorage(RegistryAddr, func(key common.Hash, value common.Hash) bool {
		execStorage[key] = value
		return true
	})

	// 3. Compare the two states
	for k, v := range allocStorage {
		assert.Equal(t, v.Hex(), execStorage[k].Hex(), k.Hex())
		t.Logf("%x %x\n", k, v)
	}
	for k, v := range execStorage {
		assert.Equal(t, v.Hex(), allocStorage[k].Hex(), k.Hex())
	}
}
