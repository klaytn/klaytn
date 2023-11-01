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
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKip113(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlWarn)

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

	// make a blockchain node
	fullNode, node, validator, _, workspace := newBlockchain(t, config, nil)
	defer func() {
		fullNode.Stop()
		os.RemoveAll(workspace)
	}()

	var (
		senderKey = validator.Keys[0]
		sender    = bind.NewKeyedTransactor(senderKey)

		chain      = node.BlockChain().(*blockchain.BlockChain)
		transactor = backends.NewBlockchainContractBackend(chain, node.TxPool().(*blockchain.TxPool), nil)

		nodeId        = common.HexToAddress("0xaaaa")
		_, pub1, pop1 = makeBlsKey()
	)

	contract, contractAddr := deployKip113Mock(t, sender, transactor, chain)

	// Register a BLS key
	tx, err := contract.Register(sender, nodeId, pub1, pop1)
	require.Nil(t, err)
	require.NotNil(t, waitReceipt(chain, tx.Hash()))

	infos, err := system.ReadKip113All(transactor, contractAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, pub1, infos[nodeId].PublicKey)
	assert.Equal(t, pop1, infos[nodeId].Pop)

	// TODO: test with Registry
}

func TestKIP113GenesisAlloc(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlWarn)

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

	var (
		ownerKey, _ = crypto.GenerateKey()
		owner       = bind.NewKeyedTransactor(ownerKey)

		KIP113MockAddr   = common.HexToAddress("0x0000000000000000000000000000000000000402")
		ERC1967ProxyAddr = common.HexToAddress("0x0000000000000000000000000000000000000403")

		nodeId        = common.HexToAddress("0xaaaa")
		nodeId2       = common.HexToAddress("0xbbbb")
		_, pub1, pop1 = makeBlsKey()
		_, pub2, pop2 = makeBlsKey()

		allocProxyStorage  = system.AllocProxy(KIP113MockAddr)
		allocKip113Storage = system.AllocKip113(system.AllocKip113Init{
			Infos: system.BlsPublicKeyInfos{
				nodeId:  {PublicKey: pub1, Pop: pop1},
				nodeId2: {PublicKey: pub2, Pop: pop2},
			},
			Owner: owner.From,
		})
	)

	// Create storage with AllocProxy() and AllocKip113()
	allocStorage := make(map[common.Hash]common.Hash)
	for k, v := range allocProxyStorage {
		allocStorage[k] = v
	}
	for k, v := range allocKip113Storage {
		allocStorage[k] = v
	}

	genesisAlloc := blockchain.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(params.KLAY),
		},
		KIP113MockAddr: {
			Code:    system.Kip113MockCode,
			Balance: common.Big0,
		},
		ERC1967ProxyAddr: {
			Code:    system.ERC1967ProxyCode,
			Storage: allocStorage,
			Balance: common.Big0,
		},
	}

	genesis := blockchain.DefaultGenesisBlock()
	for k, v := range genesisAlloc {
		genesis.Alloc[k] = v
	}

	// make a blockchain node
	fullNode, node, _, _, workspace := newBlockchain(t, config, genesis)
	defer func() {
		fullNode.Stop()
		os.RemoveAll(workspace)
	}()

	var (
		chain      = node.BlockChain().(*blockchain.BlockChain)
		transactor = backends.NewBlockchainContractBackend(chain, node.TxPool().(*blockchain.TxPool), nil)
	)

	// Check BLS keys have been allocated correctly
	infos, err := system.ReadKip113All(transactor, ERC1967ProxyAddr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(infos))
	assert.Equal(t, pub1, infos[nodeId].PublicKey)
	assert.Equal(t, pop1, infos[nodeId].Pop)
	assert.Equal(t, pub2, infos[nodeId2].PublicKey)
	assert.Equal(t, pop2, infos[nodeId2].Pop)
}

func deployKip113Mock(t *testing.T, sender *bind.TransactOpts, backend *backends.BlockchainContractBackend, chain *blockchain.BlockChain, params ...interface{}) (*contracts.KIP113MockTransactor, common.Address) {
	// Prepare input data for ERC1967Proxy constructor
	abi, err := contracts.KIP113MockMetaData.GetAbi()
	assert.Nil(t, err)
	data, err := abi.Pack("initialize")
	assert.Nil(t, err)

	// Deploy Proxy contract
	// 1. Deploy KIP113Mock implementation contract
	implAddr, tx, _, err := contracts.DeployKIP113Mock(sender, backend)
	require.Nil(t, err)
	require.NotNil(t, waitReceipt(chain, tx.Hash()))

	// 2. Deploy ERC1967Proxy(KIP113Mock.address, _data)
	contractAddr, tx, _, err := contracts.DeployERC1967Proxy(sender, backend, implAddr, data)
	require.Nil(t, err)
	require.NotNil(t, waitReceipt(chain, tx.Hash()))

	// 3. Attach KIP113Mock contract to the proxy
	contract, _ := contracts.NewKIP113MockTransactor(contractAddr, backend)

	return contract, contractAddr
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
