package system

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
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

	// Returns zero record because the only record is invalid.
	infos, err = ReadKip113All(backend, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(infos))
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
