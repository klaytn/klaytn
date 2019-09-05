// Copyright 2019 The klaytn Authors
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

package cn

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	mocks2 "github.com/klaytn/klaytn/consensus/mocks"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

const blockNum1 = 20190902

var td1 = big.NewInt(123)

var addr1 common.Address
var addr2 common.Address
var addr3 common.Address

var key1 *ecdsa.PrivateKey
var key2 *ecdsa.PrivateKey
var key3 *ecdsa.PrivateKey

var nodeID1 discover.NodeID
var nodeID2 discover.NodeID
var nodeID3 discover.NodeID

var tx *types.Transaction
var txs types.Transactions

func init() {
	key1, _ = crypto.GenerateKey()
	key2, _ = crypto.GenerateKey()
	key3, _ = crypto.GenerateKey()

	addr1 = crypto.PubkeyToAddress(key1.PublicKey)
	addr2 = crypto.PubkeyToAddress(key2.PublicKey)
	addr3 = crypto.PubkeyToAddress(key3.PublicKey)

	nodeID1 = discover.PubkeyID(&key1.PublicKey)
	nodeID2 = discover.PubkeyID(&key2.PublicKey)
	nodeID3 = discover.PubkeyID(&key3.PublicKey)

	tx = types.NewTransaction(111, addr1, big.NewInt(111), 111, big.NewInt(111), nil)
	txs = types.Transactions{tx}
}

func newMocks(t *testing.T) (*gomock.Controller, *mocks2.MockEngine, *mocks.MockBlockChain, *mocks.MockTxPool) {
	mockCtrl := gomock.NewController(t)
	mockEngine := mocks2.NewMockEngine(mockCtrl)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockTxPool := mocks.NewMockTxPool(mockCtrl)

	return mockCtrl, mockEngine, mockBlockChain, mockTxPool
}

func newBlock(blockNum int) *types.Block {
	header := &types.Header{Number: big.NewInt(int64(blockNum))}
	header.Hash()
	return types.NewBlockWithHeader(header)
}

func TestNewProtocolManager(t *testing.T) {
	//1. If consensus.Engine returns an empty Protocol, NewProtocolManager throws an error.
	{
		mockCtrl, mockEngine, mockBlockChain, mockTxPool := newMocks(t)
		defer mockCtrl.Finish()

		block := newBlock(blockNum1)
		mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)
		mockEngine.EXPECT().Protocol().Return(consensus.Protocol{}).Times(1)

		pm, err := NewProtocolManager(nil, downloader.FastSync, 0, nil, mockTxPool,
			mockEngine, mockBlockChain, nil, -1, &Config{})

		assert.Nil(t, pm)
		assert.Equal(t, errIncompatibleConfig, err)
	}
}

func TestSampleSize(t *testing.T) {
	peers := make([]Peer, minNumPeersToSendBlock-1)
	assert.Equal(t, len(peers), sampleSize(peers))

	peers = make([]Peer, 4)
	assert.Equal(t, minNumPeersToSendBlock, sampleSize(peers))

	peers = make([]Peer, 16)
	assert.Equal(t, 4, sampleSize(peers))
}

func TestSamplingPeers(t *testing.T) {
	peers := make([]Peer, 10)
	assert.Equal(t, peers, samplingPeers(peers, 20))
	assert.Equal(t, peers[:5], samplingPeers(peers, 5))
}

func TestBroadcastBlock_NoParentExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	block := newBlock(blockNum1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	td := int64(100)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().GetBlock(block.ParentHash(), block.NumberU64()-1).Return(nil).Times(1)
	mockBlockChain.EXPECT().GetTd(block.ParentHash(), block.NumberU64()-1).Return(big.NewInt(td)).Times(0)
	pm.blockchain = mockBlockChain

	mockPeers := NewMockPeerSet(mockCtrl)
	pm.peers = mockPeers

	mockPeer := NewMockPeer(mockCtrl)
	mockPeers.EXPECT().SamplePeersToSendBlock(block, pm.nodetype).Return([]Peer{mockPeer}).Times(0)
	mockPeer.EXPECT().AsyncSendNewBlock(block, new(big.Int).Add(block.BlockScore(), big.NewInt(td))).Times(0)

	pm.BroadcastBlock(block)
}

func TestBroadcastBlock_ParentExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	block := newBlock(blockNum1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	td := int64(100)
	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	mockBlockChain.EXPECT().GetBlock(block.ParentHash(), block.NumberU64()-1).Return(block).Times(1)
	mockBlockChain.EXPECT().GetTd(block.ParentHash(), block.NumberU64()-1).Return(big.NewInt(td)).Times(1)
	pm.blockchain = mockBlockChain

	mockPeers := NewMockPeerSet(mockCtrl)
	pm.peers = mockPeers

	mockPeer := NewMockPeer(mockCtrl)
	mockPeers.EXPECT().SamplePeersToSendBlock(block, pm.nodetype).Return([]Peer{mockPeer}).Times(1)
	mockPeer.EXPECT().AsyncSendNewBlock(block, new(big.Int).Add(block.BlockScore(), big.NewInt(td))).Times(1)

	pm.BroadcastBlock(block)
}

func TestBroadcastBlockHash(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	block := newBlock(blockNum1)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// When the given block doesn't exist.
	{
		mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
		mockBlockChain.EXPECT().HasBlock(block.Hash(), block.NumberU64()).Return(false).Times(1)
		pm.blockchain = mockBlockChain
		pm.BroadcastBlockHash(block)
	}

	// When the given block exists.
	{
		mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
		mockBlockChain.EXPECT().HasBlock(block.Hash(), block.NumberU64()).Return(true).Times(1)
		pm.blockchain = mockBlockChain

		mockPeer := NewMockPeer(mockCtrl)
		mockPeer.EXPECT().AsyncSendNewBlockHash(block).Times(1)

		mockPeers := NewMockPeerSet(mockCtrl)
		mockPeers.EXPECT().PeersWithoutBlock(block.Hash()).Return([]Peer{mockPeer}).Times(1)
		pm.peers = mockPeers

		pm.BroadcastBlockHash(block)
	}
}

func TestBroadcastTxsFromCN_CN_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	// Using gomock.Any() for AsyncSendTransactions calls,
	// since transactions are put into a new list inside broadcastCNTx.
	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).Times(1)
	cnPeer.EXPECT().AsyncSendTransactions(gomock.Any()).Times(0)
	pnPeer.EXPECT().AsyncSendTransactions(gomock.Any()).Times(0)
	enPeer.EXPECT().AsyncSendTransactions(gomock.Any()).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromCN_CN_Exists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	// Using gomock.Any() for AsyncSendTransactions calls,
	// since transactions are put into a new list inside broadcastCNTx.
	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)
	cnPeer.EXPECT().AsyncSendTransactions(gomock.Any()).Times(1)
	pnPeer.EXPECT().AsyncSendTransactions(gomock.Any()).Times(0)
	enPeer.EXPECT().AsyncSendTransactions(gomock.Any()).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromPN_PN_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.PROXYNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)

	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Any()).Times(1)
	pnPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)
	enPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromPN_PN_Exists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.PROXYNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)

	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Any()).Times(1)
	pnPeer.EXPECT().SendTransactions(gomock.Any()).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromEN_EN_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)

	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)
	enPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)
	pnPeer.EXPECT().SendTransactions(gomock.Any()).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromEN_EN_Exists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)

	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)
	enPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)
	pnPeer.EXPECT().SendTransactions(gomock.Any()).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Any()).Times(1)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFromEN_PN_NotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	cnPeer, pnPeer, enPeer := createAndRegisterPeers(mockCtrl, peers)

	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)

	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).Times(1)
	enPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)
	pnPeer.EXPECT().SendTransactions(gomock.Any()).Times(0)
	enPeer.EXPECT().SendTransactions(gomock.Any()).Times(1)

	pm.BroadcastTxs(txs)
}

func TestBroadcastTxsFrom_DefaultCase(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.BOOTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	createAndRegisterPeers(mockCtrl, peers)

	// There are no expected calls for the mocks.
	pm.nodetype = node.BOOTNODE
	pm.BroadcastTxs(txs)

	pm.nodetype = node.UNKNOWNNODE
	pm.BroadcastTxs(txs)
}

func TestUseTxResend(t *testing.T) {
	testSet := [...]struct {
		pm     *ProtocolManager
		result bool
	}{
		{&ProtocolManager{nodetype: node.CONSENSUSNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: node.ENDPOINTNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: node.PROXYNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: node.BOOTNODE, txResendUseLegacy: true}, false},
		{&ProtocolManager{nodetype: node.UNKNOWNNODE, txResendUseLegacy: true}, false},

		{&ProtocolManager{nodetype: node.CONSENSUSNODE, txResendUseLegacy: false}, false},
		{&ProtocolManager{nodetype: node.ENDPOINTNODE, txResendUseLegacy: false}, true},
		{&ProtocolManager{nodetype: node.PROXYNODE, txResendUseLegacy: false}, true},
		{&ProtocolManager{nodetype: node.BOOTNODE, txResendUseLegacy: false}, true},
		{&ProtocolManager{nodetype: node.UNKNOWNNODE, txResendUseLegacy: false}, true},
	}

	for _, tc := range testSet {
		assert.Equal(t, tc.result, tc.pm.useTxResend())
	}
}

func TestNodeInfo(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBlockChain := mocks.NewMockBlockChain(mockCtrl)
	pm.blockchain = mockBlockChain

	genesis := newBlock(0)
	block := newBlock(blockNum1)
	config := &params.ChainConfig{ChainID: td1}

	pm.networkId = 1234
	mockBlockChain.EXPECT().CurrentBlock().Return(block).Times(1)
	mockBlockChain.EXPECT().GetTd(block.Hash(), block.NumberU64()).Return(td1).Times(1)
	mockBlockChain.EXPECT().Genesis().Return(genesis).Times(1)
	mockBlockChain.EXPECT().Config().Return(config).Times(1)

	expected := &NodeInfo{
		Network:    pm.networkId,
		BlockScore: td1,
		Genesis:    genesis.Hash(),
		Config:     config,
		Head:       block.Hash(),
	}

	assert.Equal(t, *expected, *pm.NodeInfo())
}

func TestGetCNPeersAndGetENPeers(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peers.cnpeers[addr1] = cnPeer
	peers.pnpeers[addr2] = pnPeer
	peers.enpeers[addr3] = enPeer

	cnPeers := pm.GetCNPeers()
	enPeers := pm.GetENPeers()

	assert.Equal(t, 1, len(cnPeers))
	assert.Equal(t, 1, len(enPeers))

	assert.Equal(t, cnPeer, cnPeers[addr1])
	assert.Equal(t, enPeer, enPeers[addr3])
}

func TestFindPeers_AddrExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(1)

	targets := make(map[common.Address]bool)
	targets[addr1] = true
	targets[addr2] = true
	targets[addr3] = false

	foundPeers := pm.FindPeers(targets)

	assert.Equal(t, 2, len(foundPeers))
	assert.EqualValues(t, cnPeer, foundPeers[addr1])
	assert.EqualValues(t, pnPeer, foundPeers[addr2])
	assert.Nil(t, foundPeers[addr3])
}

func TestFindPeers_AddrNotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	pnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	enPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)

	cnPeer.EXPECT().GetP2PPeerID().Return(nodeID1).Times(1)
	pnPeer.EXPECT().GetP2PPeerID().Return(nodeID2).Times(1)
	enPeer.EXPECT().GetP2PPeerID().Return(nodeID3).Times(1)

	cnPeer.EXPECT().SetAddr(addr1).Times(1)
	pnPeer.EXPECT().SetAddr(addr2).Times(1)
	enPeer.EXPECT().SetAddr(addr3).Times(1)

	targets := make(map[common.Address]bool)
	targets[addr1] = true
	targets[addr2] = true
	targets[addr3] = false

	foundPeers := pm.FindPeers(targets)

	assert.Equal(t, 2, len(foundPeers))
	assert.EqualValues(t, cnPeer, foundPeers[addr1])
	assert.EqualValues(t, pnPeer, foundPeers[addr2])
	assert.Nil(t, foundPeers[addr3])
}

func TestFindCNPeers(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers

	cnPeer1 := NewMockPeer(mockCtrl)
	cnPeer2 := NewMockPeer(mockCtrl)
	cnPeer3 := NewMockPeer(mockCtrl)

	peers.cnpeers[addr1] = cnPeer1
	peers.cnpeers[addr2] = cnPeer2
	peers.cnpeers[addr3] = cnPeer3

	targets := make(map[common.Address]bool)
	targets[addr1] = true
	targets[addr2] = true
	targets[addr3] = false

	foundCNPeers := pm.FindCNPeers(targets)

	assert.Equal(t, 2, len(foundCNPeers))
	assert.EqualValues(t, cnPeer1, foundCNPeers[addr1])
	assert.EqualValues(t, cnPeer2, foundCNPeers[addr2])
	assert.Nil(t, foundCNPeers[addr3])
}

func TestGetPeers_AddrExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(1)

	foundAddrs := pm.GetPeers()

	assert.Equal(t, 3, len(foundAddrs))
	assert.True(t, contains(foundAddrs, addr1))
	assert.True(t, contains(foundAddrs, addr2))
	assert.True(t, contains(foundAddrs, addr3))
}

func TestGetPeers_AddrNotExists(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.ENDPOINTNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := NewMockPeerSet(mockCtrl)
	pm.peers = peers

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peersResult := map[string]Peer{"cnPeer": cnPeer, "pnPeer": pnPeer, "enPeer": enPeer}

	peers.EXPECT().Peers().Return(peersResult).Times(1)
	cnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	pnPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)
	enPeer.EXPECT().GetAddr().Return(common.Address{}).Times(1)

	cnPeer.EXPECT().GetP2PPeerID().Return(nodeID1).Times(1)
	pnPeer.EXPECT().GetP2PPeerID().Return(nodeID2).Times(1)
	enPeer.EXPECT().GetP2PPeerID().Return(nodeID3).Times(1)

	cnPeer.EXPECT().SetAddr(addr1).Times(1)
	pnPeer.EXPECT().SetAddr(addr2).Times(1)
	enPeer.EXPECT().SetAddr(addr3).Times(1)

	foundAddrs := pm.GetPeers()

	assert.Equal(t, 3, len(foundAddrs))
	assert.True(t, contains(foundAddrs, addr1))
	assert.True(t, contains(foundAddrs, addr2))
	assert.True(t, contains(foundAddrs, addr3))
}

func contains(addrs []common.Address, item common.Address) bool {
	for _, a := range addrs {
		if a == item {
			return true
		}
	}
	return false
}

func createAndRegisterPeers(mockCtrl *gomock.Controller, peers *peerSet) (*MockPeer, *MockPeer, *MockPeer) {
	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	peers.cnpeers[addr1] = cnPeer
	peers.pnpeers[addr2] = pnPeer
	peers.enpeers[addr3] = enPeer

	peers.peers[fmt.Sprintf("%x", nodeID1[:8])] = cnPeer
	peers.peers[fmt.Sprintf("%x", nodeID2[:8])] = pnPeer
	peers.peers[fmt.Sprintf("%x", nodeID3[:8])] = enPeer

	return cnPeer, pnPeer, enPeer
}
