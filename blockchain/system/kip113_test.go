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
		senderAddr   = sender.From

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KLAY),
			},
		}
		backend = backends.NewSimulatedBackend(alloc)

		_, pub1, pop1 = makeBlsKey()
		_, pub2, _    = makeBlsKey()
	)

	contractAddr, _, contract, err := contracts.DeployKIP113Mock(sender, backend)
	backend.Commit()
	assert.Nil(t, err)
	t.Logf("KIP113Mock at %x", contractAddr)

	// With a valid record
	contract.RegisterPublicKey(sender, pub1, pop1)
	backend.Commit()

	infos, err := ReadKip113All(backend, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, pub1, infos[senderAddr].PublicKey)
	assert.Equal(t, pop1, infos[senderAddr].Pop)

	// With an invalid record
	// Another registerPublicKey() call from the same sender overwrites the existing info.
	contract.RegisterPublicKey(sender, pub2, pop1) // pub vs. pop mismatch
	backend.Commit()

	// Returns zero record because invalid records have been filtered out.
	infos, err = ReadKip113All(backend, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(infos))
}

func TestAllocKIP113(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)
		senderAddr   = sender.From

		senderKey2, _ = crypto.GenerateKey()
		sender2       = bind.NewKeyedTransactor(senderKey2)
		senderAddr2   = sender2.From

		KIP113MockAddr = common.HexToAddress("0x0000000000000000000000000000000000000402")

		_, pub1, pop1 = makeBlsKey()
		_, pub2, pop2 = makeBlsKey()
	)

	// 1. Create storage with AllocKIP113
	allocStorage := AllocKIP113(&AllocKIP113Init{
		Infos: map[common.Address]BlsPublicKeyInfo{
			senderAddr: {
				PublicKey: pub1,
				Pop:       pop1,
			},
			senderAddr2: {
				PublicKey: pub2,
				Pop:       pop2,
			},
		},
		Addrs: []common.Address{
			senderAddr,
			senderAddr2,
		},
	})

	// 2. Create storage by calling registerPublicKey()
	var (
		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KLAY),
			},
			sender2.From: {
				Balance: big.NewInt(params.KLAY),
			},
			KIP113MockAddr: {
				Code:    Kip113MockCode,
				Balance: common.Big0,
			},
		}
		backend     = backends.NewSimulatedBackend(alloc)
		contract, _ = contracts.NewKIP113MockTransactor(KIP113MockAddr, backend)
	)

	contract.RegisterPublicKey(sender, pub1, pop1)
	contract.RegisterPublicKey(sender2, pub2, pop2)
	backend.Commit()

	execStorage := make(map[common.Hash]common.Hash)
	stateDB, _ := backend.BlockChain().State()
	stateDB.ForEachStorage(KIP113MockAddr, func(key common.Hash, value common.Hash) bool {
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
