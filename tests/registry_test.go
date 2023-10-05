package tests

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain/system"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

// Test Registry contract installation at Cancun hardfork
// under the condition similar to the Cypress network.
func TestRegistryFork(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlInfo)

	forkParentNum := big.NewInt(4)
	forkNum := big.NewInt(5)
	kip103Addr := common.HexToAddress("0x103103")

	// Start blockchain node
	config := params.CypressChainConfig.Copy()
	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.CancunCompatibleBlock = forkNum
	config.Kip103CompatibleBlock = common.Big1
	config.Kip103ContractAddress = kip103Addr

	fullNode, node, _, _, workspace := newBlockchain(t, config)
	defer func() {
		fullNode.Stop()
		os.RemoveAll(workspace)
	}()

	var (
		chain        = node.BlockChain()
		backend      = backends.NewBlockchainContractBackend(chain, nil, nil)
		contract, _  = system_contracts.NewRegistryCaller(system.RegistryAddr, backend)
		atForkParent = &bind.CallOpts{BlockNumber: forkParentNum}
	)
	// Wait for hardfork to pass
	waitBlock(chain, forkParentNum.Uint64())

	// Registry code is installed
	code, err := backend.CodeAt(context.Background(), system.RegistryAddr, forkParentNum)
	assert.Nil(t, err)
	assert.NotNil(t, code)

	// Registry contents are correct
	names, err := contract.GetAllNames(atForkParent)
	t.Log("Registry.getAllNames()", names)
	assert.Nil(t, err)
	assert.Equal(t, []string{system.AddressBookName, system.Kip103Name}, names)

	addr, err := contract.GetActiveAddr(atForkParent, system.AddressBookName)
	t.Log("Registry.getActiveAddr('AddressBook')", addr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, system.AddressBookAddr, addr)

	addr, err = contract.GetActiveAddr(atForkParent, system.Kip103Name)
	t.Log("Registry.getActiveAddr('TreasuryRebalance')", addr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, kip103Addr, addr)
}
