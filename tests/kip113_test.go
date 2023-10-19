// Copyright 2018 The klaytn Authors
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

package tests

import (
	"math/big"
	"math/rand"
	"os"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/system"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func TestKIP113(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlInfo)

	// prepare chain configuration
	config := params.CypressChainConfig.Copy()
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
	config.Governance.Reward.MintingAmount = new(big.Int).Mul(big.NewInt(9000000000000000000), big.NewInt(params.KLAY))

	// make a blockchain node
	fullNode, node, validator, _, workspace := newBlockchain(t, config)
	defer func() {
		fullNode.Stop()
		os.RemoveAll(workspace)
	}()

	var (
		senderKey  = validator.Keys[0]
		sender     = bind.NewKeyedTransactor(senderKey)
		senderAddr = sender.From

		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)

		_, pub1, pop1 = makeBlsKey()
		_, pub2, _    = makeBlsKey()
	)
	contractAddr, tx, contract, err := contracts.DeployKIP113Mock(sender, transactor)
	if err != nil {
		t.Fatal(err)
	}
	checkReceipt(t, node.BlockChain().(*blockchain.BlockChain), tx.Hash())

	// Register BLS key for sender
	if tx, err = contract.RegisterPublicKey(sender, pub1, pop1); err != nil {
		t.Fatal(err)
	}
	checkReceipt(t, node.BlockChain().(*blockchain.BlockChain), tx.Hash())

	infos, err := system.ReadKip113All(transactor, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, pub1, infos[senderAddr].PublicKey)
	assert.Equal(t, pop1, infos[senderAddr].Pop)

	// Unregister BLS key for sender
	if tx, err = contract.UnregisterPublicKey(sender, senderAddr); err != nil {
		t.Fatal(err)
	}
	checkReceipt(t, node.BlockChain().(*blockchain.BlockChain), tx.Hash())

	infos, err = system.ReadKip113All(transactor, contractAddr, nil)
	assert.Nil(t, err)
	assert.Nil(t, infos[senderAddr].PublicKey)
	assert.Nil(t, infos[senderAddr].Pop)
	assert.Equal(t, 0, len(infos))

	// Register invalid BLS key for sender
	if tx, err = contract.RegisterPublicKey(sender, pub2, pop1); err != nil {
		t.Fatal(err)
	}
	checkReceipt(t, node.BlockChain().(*blockchain.BlockChain), tx.Hash())

	// Invalid BLS key will be filtered out in REadKip113All function
	infos, err = system.ReadKip113All(transactor, contractAddr, nil)
	assert.Nil(t, err)
	assert.Nil(t, infos[senderAddr].PublicKey)
	assert.Nil(t, infos[senderAddr].Pop)
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

func checkReceipt(t *testing.T, chain *blockchain.BlockChain, txHash common.Hash) {
	receipt := waitReceipt(chain, txHash)
	if receipt == nil {
		t.Fatal("timeout")
	}
}
