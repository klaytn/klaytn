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
	log.EnableLogForTest(log.LvlError, log.LvlWarn)

	var (
		forkParentNum_1 = big.NewInt(3)
		forkParentNum   = big.NewInt(4)
		forkNum         = big.NewInt(5)
		bgctx           = context.Background()
		acmeAddr        = common.HexToAddress("0xaaaa")
		ownerAddr       = common.HexToAddress("0xbbbb")
	)

	// Start blockchain node
	config := params.CypressChainConfig.Copy()
	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
	config.LondonCompatibleBlock = common.Big0
	config.IstanbulCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0
	config.ShanghaiCompatibleBlock = common.Big0
	config.CancunCompatibleBlock = forkNum
	config.RandaoCompatibleBlock = forkNum
	config.RandaoRegistry = &params.RegistryConfig{
		Records: map[string]common.Address{
			"AcmeContract": acmeAddr,
		},
		Owner: ownerAddr,
	}

	fullNode, node, _, _, workspace := newBlockchain(t, config, nil)
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

	// Registry code is installed exactly at forkParentNum
	code, err := backend.CodeAt(bgctx, system.RegistryAddr, forkParentNum_1)
	assert.Nil(t, err)
	assert.Empty(t, code)

	code, err = backend.CodeAt(context.Background(), system.RegistryAddr, forkParentNum)
	assert.Nil(t, err)
	assert.NotNil(t, code)

	// Registry contents are correct
	names, err := contract.GetAllNames(atForkParent)
	t.Log("Registry.getAllNames()", names)
	assert.Nil(t, err)
	assert.Equal(t, []string{"AcmeContract"}, names)

	addr, err := contract.GetActiveAddr(atForkParent, "AcmeContract")
	t.Log("Registry.getActiveAddr('AcmeContract')", addr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, acmeAddr, addr)

	addr, err = contract.Owner(atForkParent)
	t.Log("Registry.owner()", ownerAddr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, ownerAddr, addr)

	// Registry accessors are correct
	addr, err = system.ReadRegistryActiveAddr(backend, "AcmeContract", forkParentNum_1)
	assert.ErrorIs(t, err, system.ErrRegistryNotInstalled)
	assert.Empty(t, addr)

	addr, err = system.ReadRegistryActiveAddr(backend, "AcmeContract", forkParentNum)
	assert.Nil(t, err)
	assert.Equal(t, acmeAddr, addr)
}
