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

func TestSystemRegistry(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)
		senderAddr   = sender.From

		alloc = blockchain.GenesisAlloc{
			senderAddr: {
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
