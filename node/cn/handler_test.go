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
	consensusmocks "github.com/klaytn/klaytn/consensus/mocks"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn/mocks"
	"github.com/klaytn/klaytn/params"
	workmocks "github.com/klaytn/klaytn/work/mocks"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

const blockNum1 = 20190902

var td1 = big.NewInt(123)

const numVals = 6

var addrs []common.Address
var keys []*ecdsa.PrivateKey
var nodeids []discover.NodeID

var tx1 *types.Transaction
var txs types.Transactions

var hash1 common.Hash

var signer types.Signer

func init() {
	addrs = make([]common.Address, numVals)
	keys = make([]*ecdsa.PrivateKey, numVals)
	nodeids = make([]discover.NodeID, numVals)

	for i := range keys {
		keys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
		nodeids[i] = discover.PubkeyID(&keys[i].PublicKey)
	}

	signer := types.MakeSigner(params.BFTTestChainConfig, big.NewInt(2019))
	tx1 = types.NewTransaction(111, addrs[0], big.NewInt(111), 111, big.NewInt(111), addrs[0][:])

	tx1.Sign(signer, keys[0])
	tx1.Size()
	txs = types.Transactions{tx1}

	hash1 = tx1.Hash()
}

func newMocks(t *testing.T) (*gomock.Controller, *consensusmocks.MockEngine, *workmocks.MockBlockChain, *workmocks.MockTxPool) {
	mockCtrl := gomock.NewController(t)
	mockEngine := consensusmocks.NewMockEngine(mockCtrl)
	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
	mockTxPool := workmocks.NewMockTxPool(mockCtrl)

	return mockCtrl, mockEngine, mockBlockChain, mockTxPool
}

func newBlock(blockNum int) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(blockNum)),
		BlockScore: big.NewInt(int64(1)),
		Extra:      addrs[0][:],
		Governance: addrs[0][:],
		Vote:       addrs[0][:],
	}
	header.Hash()
	block := types.NewBlockWithHeader(header)
	block = block.WithBody(types.Transactions{})
	block.Hash()
	block.Size()
	block.BlockScore()
	return block
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
	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
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
	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
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
		mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
		mockBlockChain.EXPECT().HasBlock(block.Hash(), block.NumberU64()).Return(false).Times(1)
		pm.blockchain = mockBlockChain
		pm.BroadcastBlockHash(block)
	}

	// When the given block exists.
	{
		mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
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

	// Using gomock.Eq(txs) for AsyncSendTransactions calls,
	// since transactions are put into a new list inside broadcastCNTx.
	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)
	cnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	pnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)

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

	// Using gomock.Eq(txs) for AsyncSendTransactions calls,
	// since transactions are put into a new list inside broadcastCNTx.
	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)
	cnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(1)
	pnPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().AsyncSendTransactions(gomock.Eq(txs)).Times(0)

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

	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)

	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)

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

	cnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)

	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)

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

	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)
	enPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)

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

	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)
	enPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

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

	pnPeer.EXPECT().KnowsTx(tx1.Hash()).Return(true).Times(1)
	enPeer.EXPECT().KnowsTx(tx1.Hash()).Return(false).Times(1)

	cnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(0)
	enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

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

func TestReBroadcastTxs_CN(t *testing.T) {
	pm := &ProtocolManager{}
	pm.nodetype = node.CONSENSUSNODE
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	peers := newPeerSet()
	pm.peers = peers
	createAndRegisterPeers(mockCtrl, peers)

	pm.ReBroadcastTxs(txs)
}

func TestReBroadcastTxs_PN(t *testing.T) {
	// CN Peer=0, PN Peer=1
	{
		pm := &ProtocolManager{}
		pm.nodetype = node.PROXYNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		enPeer := NewMockPeer(mockCtrl)
		enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.enpeers[addrs[2]] = enPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = enPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
	// CN Peer=1, PN Peer=0
	{
		pm := &ProtocolManager{}
		pm.nodetype = node.PROXYNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
		pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.pnpeers[addrs[2]] = pnPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = pnPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
}

func TestReBroadcastTxs_EN(t *testing.T) {
	// PN Peer=0, EN Peer=1
	{
		pm := &ProtocolManager{}
		pm.nodetype = node.ENDPOINTNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		enPeer := NewMockPeer(mockCtrl)
		enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)
		enPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.enpeers[addrs[2]] = enPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = enPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
	// PN Peer=1, EN Peer=0
	{
		pm := &ProtocolManager{}
		pm.nodetype = node.ENDPOINTNODE
		mockCtrl := gomock.NewController(t)

		peers := newPeerSet()
		pm.peers = peers

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
		pnPeer.EXPECT().SendTransactions(gomock.Eq(txs)).Times(1)

		peers.pnpeers[addrs[2]] = pnPeer
		peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = pnPeer

		pm.ReBroadcastTxs(txs)

		mockCtrl.Finish()
	}
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

	mockBlockChain := workmocks.NewMockBlockChain(mockCtrl)
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

	peers.cnpeers[addrs[0]] = cnPeer
	peers.pnpeers[addrs[1]] = pnPeer
	peers.enpeers[addrs[2]] = enPeer

	cnPeers := pm.GetCNPeers()
	enPeers := pm.GetENPeers()

	assert.Equal(t, 1, len(cnPeers))
	assert.Equal(t, 1, len(enPeers))

	assert.Equal(t, cnPeer, cnPeers[addrs[0]])
	assert.Equal(t, enPeer, enPeers[addrs[2]])
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
	cnPeer.EXPECT().GetAddr().Return(addrs[0]).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addrs[1]).Times(1)
	enPeer.EXPECT().GetAddr().Return(addrs[2]).Times(1)

	targets := make(map[common.Address]bool)
	targets[addrs[0]] = true
	targets[addrs[1]] = true
	targets[addrs[2]] = false

	foundPeers := pm.FindPeers(targets)

	assert.Equal(t, 2, len(foundPeers))
	assert.EqualValues(t, cnPeer, foundPeers[addrs[0]])
	assert.EqualValues(t, pnPeer, foundPeers[addrs[1]])
	assert.Nil(t, foundPeers[addrs[2]])
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

	cnPeer.EXPECT().GetP2PPeerID().Return(nodeids[0]).Times(1)
	pnPeer.EXPECT().GetP2PPeerID().Return(nodeids[1]).Times(1)
	enPeer.EXPECT().GetP2PPeerID().Return(nodeids[2]).Times(1)

	cnPeer.EXPECT().SetAddr(addrs[0]).Times(1)
	pnPeer.EXPECT().SetAddr(addrs[1]).Times(1)
	enPeer.EXPECT().SetAddr(addrs[2]).Times(1)

	targets := make(map[common.Address]bool)
	targets[addrs[0]] = true
	targets[addrs[1]] = true
	targets[addrs[2]] = false

	foundPeers := pm.FindPeers(targets)

	assert.Equal(t, 2, len(foundPeers))
	assert.EqualValues(t, cnPeer, foundPeers[addrs[0]])
	assert.EqualValues(t, pnPeer, foundPeers[addrs[1]])
	assert.Nil(t, foundPeers[addrs[2]])
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

	peers.cnpeers[addrs[0]] = cnPeer1
	peers.cnpeers[addrs[1]] = cnPeer2
	peers.cnpeers[addrs[2]] = cnPeer3

	targets := make(map[common.Address]bool)
	targets[addrs[0]] = true
	targets[addrs[1]] = true
	targets[addrs[2]] = false

	foundCNPeers := pm.FindCNPeers(targets)

	assert.Equal(t, 2, len(foundCNPeers))
	assert.EqualValues(t, cnPeer1, foundCNPeers[addrs[0]])
	assert.EqualValues(t, cnPeer2, foundCNPeers[addrs[1]])
	assert.Nil(t, foundCNPeers[addrs[2]])
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
	cnPeer.EXPECT().GetAddr().Return(addrs[0]).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addrs[1]).Times(1)
	enPeer.EXPECT().GetAddr().Return(addrs[2]).Times(1)

	foundAddrs := pm.GetPeers()

	assert.Equal(t, 3, len(foundAddrs))
	assert.True(t, contains(foundAddrs, addrs[0]))
	assert.True(t, contains(foundAddrs, addrs[1]))
	assert.True(t, contains(foundAddrs, addrs[2]))
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

	cnPeer.EXPECT().GetP2PPeerID().Return(nodeids[0]).Times(1)
	pnPeer.EXPECT().GetP2PPeerID().Return(nodeids[1]).Times(1)
	enPeer.EXPECT().GetP2PPeerID().Return(nodeids[2]).Times(1)

	cnPeer.EXPECT().SetAddr(addrs[0]).Times(1)
	pnPeer.EXPECT().SetAddr(addrs[1]).Times(1)
	enPeer.EXPECT().SetAddr(addrs[2]).Times(1)

	foundAddrs := pm.GetPeers()

	assert.Equal(t, 3, len(foundAddrs))
	assert.True(t, contains(foundAddrs, addrs[0]))
	assert.True(t, contains(foundAddrs, addrs[1]))
	assert.True(t, contains(foundAddrs, addrs[2]))
}

func TestEnqueue(t *testing.T) {
	pm := &ProtocolManager{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fetcherMock := mocks.NewMockProtocolManagerFetcher(mockCtrl)
	pm.fetcher = fetcherMock

	block := newBlock(blockNum1)
	id := nodeids[0].String()

	fetcherMock.EXPECT().Enqueue(id, block).Times(1)
	pm.Enqueue(id, block)
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

	peers.cnpeers[addrs[0]] = cnPeer
	peers.pnpeers[addrs[1]] = pnPeer
	peers.enpeers[addrs[2]] = enPeer

	peers.peers[fmt.Sprintf("%x", nodeids[0][:8])] = cnPeer
	peers.peers[fmt.Sprintf("%x", nodeids[1][:8])] = pnPeer
	peers.peers[fmt.Sprintf("%x", nodeids[2][:8])] = enPeer

	return cnPeer, pnPeer, enPeer
}
