package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/system"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test full Randao hardfork scenario under the condition similar to the Cypress network.
func TestRandaoFork_Deploy(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	// Test parameters
	var (
		numNodes   = 1
		forkNum    = big.NewInt(15)
		owner      = bind.NewKeyedTransactor(deriveTestAccount(5))
		kip113Addr = crypto.CreateAddress(owner.From, uint64(1)) // predict deployed address.
		config     = testRandao_config(forkNum, owner.From, kip113Addr)
	)

	// Start the chain
	ctx, err := newBlockchainTestContext(&blockchainTestOverrides{
		numNodes:    numNodes,
		numAccounts: 8,
		config:      config,
	})
	require.Nil(t, err)
	ctx.Subscribe(t, func(ev *blockchain.ChainEvent) {
		b := ev.Block
		t.Logf("block[%3d] txs=%d mixHash=%x", b.NumberU64(), b.Transactions().Len(), b.Header().MixHash)
	})
	ctx.Start()
	defer ctx.Cleanup()

	// Deploy KIP113 before hardfork.
	// Note that its address is already configured in RandaoRegistry.
	_, actualKip113Addr := testRandao_deployKip113(t, ctx)
	assert.Equal(t, kip113Addr, actualKip113Addr) // check the prediced address

	// Pass the hardfork block
	ctx.WaitBlock(t, forkNum.Uint64())

	// Inspect the chain
	testRandao_checkRegistry(t, ctx)
	testRandao_checkKip113(t, ctx)

	// Propose by each node for once
	ctx.WaitBlock(t, forkNum.Uint64()+uint64(numNodes))
}

// Test Randao hardfork scenario where it's enabled from the genesis
func TestRandaoFork_Genesis(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	// Test parameters
	var (
		numNodes   = 1
		forkNum    = big.NewInt(5)
		owner      = bind.NewKeyedTransactor(deriveTestAccount(5))
		kip113Addr = common.HexToAddress("0x0000000000000000000000000000000000000403")
		config     = testRandao_config(forkNum, owner.From, kip113Addr)
		alloc      = testRandao_alloc(numNodes, kip113Addr)
	)

	// Start the chain
	ctx, err := newBlockchainTestContext(&blockchainTestOverrides{
		numNodes:    numNodes,
		numAccounts: 8,
		config:      config,
		alloc:       alloc,
	})
	require.Nil(t, err)
	ctx.Subscribe(t, func(ev *blockchain.ChainEvent) {
		b := ev.Block
		t.Logf("block[%3d] txs=%d mixHash=%x", b.NumberU64(), b.Transactions().Len(), b.Header().MixHash)
	})
	ctx.Start()
	defer ctx.Cleanup()

	// Pass the hardfork block
	ctx.WaitBlock(t, forkNum.Uint64())

	// Inspect the chain
	testRandao_checkRegistry(t, ctx)
	testRandao_checkKip113(t, ctx)

	// Propose by each node for once
	ctx.WaitBlock(t, forkNum.Uint64()+uint64(numNodes))
}

// Make ChainConfig that hardforks at `forkNum` and the Registry owner be `owner`.
func testRandao_config(forkNum *big.Int, owner, kip113Addr common.Address) *params.ChainConfig {
	config := blockchainTestChainConfig.Copy()
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
			system.Kip113Name: kip113Addr,
		},
		Owner: owner,
	}
	return config
}

// Make GenesisAlloc that contains node BLS public keys
func testRandao_alloc(numNodes int, kip113Addr common.Address) blockchain.GenesisAlloc {
	infos := make(system.BlsPublicKeyInfos)
	for i := 0; i < numNodes; i++ {
		var (
			key   = deriveTestAccount(i)
			addr  = crypto.PubkeyToAddress(key.PublicKey)
			sk, _ = bls.DeriveFromECDSA(key)
			pk    = sk.PublicKey().Marshal()
			pop   = bls.PopProve(sk).Marshal()
		)
		infos[addr] = system.BlsPublicKeyInfo{PublicKey: pk, Pop: pop}
	}

	var (
		implAddr = common.HexToAddress("0x0000000000000000000000000000000000000402")
		owner    = crypto.PubkeyToAddress(deriveTestAccount(5).PublicKey)

		proxyStorage = system.AllocProxy(implAddr)
		implStorage  = system.AllocKip113(system.AllocKip113Init{
			Infos: infos,
			Owner: owner,
		})
		storage = system.MergeStorage(proxyStorage, implStorage)
	)

	return blockchain.GenesisAlloc{
		implAddr: {
			Code:    system.Kip113MockCode,
			Balance: common.Big0,
		},
		kip113Addr: {
			Code:    system.ERC1967ProxyCode,
			Storage: storage,
			Balance: common.Big0,
		},
	}
}

// Deploy KIP-113 contract
func testRandao_deployKip113(t *testing.T, ctx *blockchainTestContext) (*system_contracts.KIP113Mock, common.Address) {
	var (
		sender = bind.NewKeyedTransactor(deriveTestAccount(5))

		abi, _      = system_contracts.KIP113MockMetaData.GetAbi()
		initData, _ = abi.Pack("initialize")

		chain   = ctx.nodes[0].cn.BlockChain()
		txpool  = ctx.nodes[0].cn.TxPool().(*blockchain.TxPool)
		backend = backends.NewBlockchainContractBackend(chain, txpool, nil)
	)

	// Deploy implementation and proxy
	implAddr, tx, _, err := system_contracts.DeployKIP113Mock(sender, backend)
	assert.Nil(t, err)
	ctx.WaitTx(t, tx.Hash())

	proxyAddr, tx, _, err := system_contracts.DeployERC1967Proxy(sender, backend, implAddr, initData)
	assert.Nil(t, err)
	ctx.WaitTx(t, tx.Hash())

	t.Logf("Kip113 impl=%s proxy=%s", implAddr.Hex(), proxyAddr.Hex())
	kip113, _ := system_contracts.NewKIP113Mock(proxyAddr, backend)

	// Register node BLS public keys
	var txs []*types.Transaction
	for i := 0; i < ctx.numNodes; i++ {
		var (
			addr  = ctx.accountAddrs[i]
			sk, _ = bls.DeriveFromECDSA(ctx.accountKeys[i])
			pk    = sk.PublicKey().Marshal()
			pop   = bls.PopProve(sk).Marshal()
		)
		t.Logf("node[%2d] addr=%x blsPub=%x", i, addr, pk)

		tx, err := kip113.Register(sender, addr, pk, pop)
		txs = append(txs, tx)
		assert.Nil(t, err)
	}
	for _, tx := range txs {
		ctx.WaitTx(t, tx.Hash())
	}

	infos, _ := system.ReadKip113All(backend, proxyAddr, nil)
	t.Logf("Kip113 getAllBlsInfo().length=%d", len(infos))

	return kip113, proxyAddr
}

// Inspect the given chain
func testRandao_checkRegistry(t *testing.T, ctx *blockchainTestContext) {
	var (
		forkNum         = int64(ctx.config.RandaoCompatibleBlock.Uint64())
		forkParentNum   = big.NewInt(forkNum - 1)
		forkParentNum_1 = big.NewInt(forkNum - 2)
		atForkParent    = &bind.CallOpts{BlockNumber: forkParentNum}
		kip113Addr      = ctx.config.RandaoRegistry.Records[system.Kip113Name]
		ownerAddr       = ctx.config.RandaoRegistry.Owner

		bgctx       = context.Background()
		chain       = ctx.nodes[0].cn.BlockChain()
		backend     = backends.NewBlockchainContractBackend(chain, nil, nil)
		registry, _ = system_contracts.NewRegistryCaller(system.RegistryAddr, backend)
	)

	// Registry code is installed exactly at forkParentNum
	code, err := backend.CodeAt(bgctx, system.RegistryAddr, forkParentNum_1)
	assert.Nil(t, err)
	assert.Empty(t, code)

	code, err = backend.CodeAt(bgctx, system.RegistryAddr, forkParentNum)
	assert.Nil(t, err)
	assert.NotNil(t, code)

	// Registry contents are correct
	names, err := registry.GetAllNames(atForkParent)
	t.Logf("Registry.getAllNames()=%v", names)
	assert.Nil(t, err)
	assert.Equal(t, []string{system.Kip113Name}, names)

	addr, err := registry.GetActiveAddr(atForkParent, system.Kip113Name)
	t.Logf("Registry.getActiveAddr('KIP113')=%s", addr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, kip113Addr, addr)

	addr, err = registry.Owner(atForkParent)
	t.Logf("Registry.owner()=%s", ownerAddr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, ownerAddr, addr)

	// Registry accessors are correct
	addr, err = system.ReadRegistryActiveAddr(backend, system.Kip113Name, forkParentNum_1)
	assert.ErrorIs(t, err, system.ErrRegistryNotInstalled)
	assert.Empty(t, addr)

	addr, err = system.ReadRegistryActiveAddr(backend, system.Kip113Name, forkParentNum)
	assert.Nil(t, err)
	assert.Equal(t, kip113Addr, addr)
}

func testRandao_checkKip113(t *testing.T, ctx *blockchainTestContext) {
	var (
		forkNum = ctx.config.RandaoCompatibleBlock
		chain   = ctx.nodes[0].cn.BlockChain()
		backend = backends.NewBlockchainContractBackend(chain, nil, nil)
	)

	kip113Addr, err := system.ReadRegistryActiveAddr(backend, system.Kip113Name, forkNum)
	assert.Nil(t, err)

	// BLS public keys of every nodes are registered
	infos, err := system.ReadKip113All(backend, kip113Addr, forkNum)
	t.Logf("Kip113.getAllBlsInfo()=%v", infos.String())
	assert.Nil(t, err)
	assert.Len(t, infos, ctx.numNodes)
	for i := 0; i < ctx.numNodes; i++ {
		addr := ctx.accountAddrs[i]
		assert.Contains(t, infos, addr)
	}
}
