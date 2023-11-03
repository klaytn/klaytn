package system

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func TestReadKip113(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KLAY),
			},
		}
		backend = backends.NewSimulatedBackend(alloc)

		nodeId        = common.HexToAddress("0xaaaa")
		_, pub1, pop1 = makeBlsKey()
		_, pub2, _    = makeBlsKey()
	)

	// Deploy Proxy contract
	transactor, contractAddr := deployKip113Mock(t, sender, backend)

	// With a valid record
	transactor.Register(sender, nodeId, pub1, pop1)
	backend.Commit()

	caller, _ := contracts.NewKIP113Caller(contractAddr, backend)

	opts := &bind.CallOpts{BlockNumber: nil}
	owner, _ := caller.Owner(opts)
	assert.Equal(t, sender.From, owner)
	t.Logf("owner: %x", owner)

	infos, err := ReadKip113All(backend, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, pub1, infos[nodeId].PublicKey)
	assert.Equal(t, pop1, infos[nodeId].Pop)

	// With an invalid record
	// Another register() call for the same nodeId overwrites the existing info.
	transactor.Register(sender, nodeId, pub2, pop1) // pub vs. pop mismatch
	backend.Commit()

	// Returns zero record because invalid records have been filtered out.
	infos, err = ReadKip113All(backend, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(infos))
}

func TestAllocKip113(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)

		KIP113MockAddr = common.HexToAddress("0x0000000000000000000000000000000000000402")

		nodeId1       = common.HexToAddress("0xaaaa")
		nodeId2       = common.HexToAddress("0xbbbb")
		_, pub1, pop1 = makeBlsKey()
		_, pub2, pop2 = makeBlsKey()

		abi, _   = contracts.KIP113MockMetaData.GetAbi()
		input, _ = abi.Pack("initialize")

		allocProxyStorage  = AllocProxy(KIP113MockAddr)
		allocKip113Storage = AllocKip113(AllocKip113Init{
			Infos: BlsPublicKeyInfos{
				nodeId1: {PublicKey: pub1, Pop: pop1},
				nodeId2: {PublicKey: pub2, Pop: pop2},
			},
			Owner: sender.From,
		})
	)

	// 1. Merge two storage maps
	allocStorage := MergeStorage(allocProxyStorage, allocKip113Storage)

	// 2. Create storage by calling register()
	var (
		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KLAY),
			},
			KIP113MockAddr: {
				Code:    Kip113MockCode,
				Balance: common.Big0,
			},
		}
		backend               = backends.NewSimulatedBackend(alloc)
		contractAddr, _, _, _ = contracts.DeployERC1967Proxy(sender, backend, KIP113MockAddr, input)
	)
	backend.Commit()

	contract, _ := contracts.NewKIP113MockTransactor(contractAddr, backend)

	contract.Register(sender, nodeId1, pub1, pop1)
	contract.Register(sender, nodeId2, pub2, pop2)
	backend.Commit()

	execStorage := make(map[common.Hash]common.Hash)
	stateDB, _ := backend.BlockChain().State()
	stateDB.ForEachStorage(contractAddr, func(key common.Hash, value common.Hash) bool {
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

func deployKip113Mock(t *testing.T, sender *bind.TransactOpts, backend *backends.SimulatedBackend, params ...interface{}) (*contracts.KIP113MockTransactor, common.Address) {
	// Prepare input data for ERC1967Proxy constructor
	abi, err := contracts.KIP113MockMetaData.GetAbi()
	assert.Nil(t, err)
	data, err := abi.Pack("initialize")
	assert.Nil(t, err)

	// Deploy Proxy contract
	// 1. Deploy KIP113Mock implementation contract
	implAddr, _, _, err := contracts.DeployKIP113Mock(sender, backend)
	backend.Commit()
	assert.Nil(t, err)
	t.Logf("KIP113Mock impl at %x", implAddr)

	// 2. Deploy ERC1967Proxy(KIP113Mock.address, _data)
	contractAddr, _, _, err := contracts.DeployERC1967Proxy(sender, backend, implAddr, data)
	backend.Commit()
	assert.Nil(t, err)
	t.Logf("ERC1967Proxy at %x", contractAddr)

	// 3. Attach KIP113Mock contract to the proxy
	transactor, _ := contracts.NewKIP113MockTransactor(contractAddr, backend)

	return transactor, contractAddr
}

func makeBlsKey() (priv, pub, pop []byte) {
	ikm := make([]byte, 32)
	rand.Read(ikm)

	sk, _ := bls.GenerateKey(ikm)
	pk := sk.PublicKey()
	sig := bls.PopProve(sk)

	priv = sk.Marshal()
	pub = pk.Marshal()
	pop = sig.Marshal()
	if len(priv) != 32 || len(pub) != 48 || len(pop) != 96 {
		panic("bad bls key")
	}
	return
}
