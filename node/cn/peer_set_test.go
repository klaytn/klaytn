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
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/node"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func setMockPeers(mockPeers []*MockPeer) {
	for i, mp := range mockPeers {
		mp.EXPECT().GetAddr().Return(addrs[i]).AnyTimes()
		mp.EXPECT().GetID().Return(nodeids[i].String()).AnyTimes()
		mp.EXPECT().Broadcast().AnyTimes()
	}
}

func setMockPeersConnType(cnPeer, pnPeer, enPeer *MockPeer) {
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).AnyTimes()
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).AnyTimes()
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).AnyTimes()
}

func TestPeerSet_Register(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))
	assert.Equal(t, 3, len(peerSet.peers))

	assert.Equal(t, errAlreadyRegistered, peerSet.Register(cnPeer))
	assert.Equal(t, errAlreadyRegistered, peerSet.Register(pnPeer))
	assert.Equal(t, errAlreadyRegistered, peerSet.Register(enPeer))
	assert.Equal(t, 3, len(peerSet.peers))

	peerSet.closed = true
	assert.Equal(t, errClosed, peerSet.Register(cnPeer))
	assert.Equal(t, errClosed, peerSet.Register(pnPeer))
	assert.Equal(t, errClosed, peerSet.Register(enPeer))
}

func TestPeerSet_Unregister(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	cnPeer.EXPECT().Close().Times(1)
	pnPeer.EXPECT().Close().Times(1)
	enPeer.EXPECT().Close().Times(1)

	assert.Equal(t, errNotRegistered, peerSet.Unregister(nodeids[0].String()))
	assert.Equal(t, errNotRegistered, peerSet.Unregister(nodeids[1].String()))
	assert.Equal(t, errNotRegistered, peerSet.Unregister(nodeids[2].String()))

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))
	assert.Equal(t, 3, len(peerSet.peers))

	assert.NoError(t, peerSet.Unregister(nodeids[0].String()))
	assert.NoError(t, peerSet.Unregister(nodeids[1].String()))
	assert.NoError(t, peerSet.Unregister(nodeids[2].String()))
}

func TestPeerSet_Peers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	peers := peerSet.Peers()
	expectedPeers := map[string]Peer{nodeids[0].String(): cnPeer, nodeids[1].String(): pnPeer, nodeids[2].String(): enPeer}
	assert.EqualValues(t, expectedPeers, peers)
}

func TestPeerSet_CNPeers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	setMockPeers([]*MockPeer{cnPeer})

	assert.EqualValues(t, map[common.Address]Peer{}, peerSet.CNPeers())
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.EqualValues(t, map[common.Address]Peer{addrs[0]: cnPeer}, peerSet.CNPeers())
}

func TestPeerSet_PNPeers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	setMockPeers([]*MockPeer{pnPeer})

	assert.EqualValues(t, map[common.Address]Peer{}, peerSet.PNPeers())
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.EqualValues(t, map[common.Address]Peer{addrs[0]: pnPeer}, peerSet.PNPeers())
}

func TestPeerSet_ENPeers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	setMockPeers([]*MockPeer{enPeer})

	assert.EqualValues(t, map[common.Address]Peer{}, peerSet.ENPeers())
	assert.NoError(t, peerSet.Register(enPeer))
	assert.EqualValues(t, map[common.Address]Peer{addrs[0]: enPeer}, peerSet.ENPeers())
}

func TestPeerSet_Peer_And_Len(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	assert.Equal(t, 0, peerSet.Len())
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.Equal(t, 1, peerSet.Len())
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.Equal(t, 2, peerSet.Len())
	assert.NoError(t, peerSet.Register(enPeer))
	assert.Equal(t, 3, peerSet.Len())

	assert.Equal(t, cnPeer, peerSet.Peer(nodeids[0].String()))
	assert.Equal(t, pnPeer, peerSet.Peer(nodeids[1].String()))
	assert.Equal(t, enPeer, peerSet.Peer(nodeids[2].String()))
}

func TestPeerSet_PeersWithoutBlock(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	block := newBlock(blockNum1)

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	cnPeer.EXPECT().KnowsBlock(block.Hash()).Return(false).AnyTimes()
	pnPeer.EXPECT().KnowsBlock(block.Hash()).Return(true).AnyTimes()
	enPeer.EXPECT().KnowsBlock(block.Hash()).Return(false).AnyTimes()

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	peersWithoutBlock := peerSet.PeersWithoutBlock(block.Hash())
	assert.Equal(t, 2, len(peersWithoutBlock))
	assert.EqualValues(t, []Peer{cnPeer, enPeer}, peersWithoutBlock)

	cnWithoutBlock := peerSet.CNWithoutBlock(block.Hash())
	assert.Equal(t, 1, len(cnWithoutBlock))
	assert.Equal(t, []Peer{cnPeer}, cnWithoutBlock)

	pnWithoutBlock := peerSet.PNWithoutBlock(block.Hash())
	assert.Equal(t, 0, len(pnWithoutBlock))
	assert.Equal(t, []Peer{}, pnWithoutBlock)

	enWithoutBlock := peerSet.ENWithoutBlock(block.Hash())
	assert.Equal(t, 1, len(enWithoutBlock))
	assert.Equal(t, []Peer{enPeer}, enWithoutBlock)
}

func TestPeerSet_PeersWithoutTx(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tx := types.NewTransaction(111, addrs[0], big.NewInt(111), 111, big.NewInt(111), nil)

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).AnyTimes()
	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).AnyTimes()
	enPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).AnyTimes()

	assert.EqualValues(t, []Peer{}, peerSet.PeersWithoutTx(tx.Hash()))

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	peersWithoutTx := peerSet.PeersWithoutTx(tx.Hash())

	assert.Equal(t, 2, len(peersWithoutTx))
	assert.EqualValues(t, []Peer{cnPeer, enPeer}, peersWithoutTx)
}

func TestPeerSet_TypePeersWithoutTx(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tx := types.NewTransaction(111, addrs[0], big.NewInt(111), 111, big.NewInt(111), nil)

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).AnyTimes()
	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).AnyTimes()
	enPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).AnyTimes()

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	assert.EqualValues(t, []Peer{}, peerSet.TypePeersWithoutTx(tx.Hash(), node.CONSENSUSNODE))
	assert.EqualValues(t, []Peer{}, peerSet.TypePeersWithoutTx(tx.Hash(), node.PROXYNODE))
	assert.EqualValues(t, []Peer{}, peerSet.TypePeersWithoutTx(tx.Hash(), node.ENDPOINTNODE))

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	assert.EqualValues(t, []Peer{cnPeer}, peerSet.TypePeersWithoutTx(tx.Hash(), node.CONSENSUSNODE))
	assert.EqualValues(t, []Peer{}, peerSet.TypePeersWithoutTx(tx.Hash(), node.PROXYNODE))
	assert.EqualValues(t, []Peer{enPeer}, peerSet.TypePeersWithoutTx(tx.Hash(), node.ENDPOINTNODE))
}

func TestPeerSet_CNWithoutTx(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tx := types.NewTransaction(111, addrs[0], big.NewInt(111), 111, big.NewInt(111), nil)

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).AnyTimes()
	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)
	setMockPeers([]*MockPeer{cnPeer})

	assert.EqualValues(t, []Peer{}, peerSet.CNWithoutTx(tx.Hash()))
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.EqualValues(t, []Peer{cnPeer}, peerSet.CNWithoutTx(tx.Hash()))
}

func TestPeerSet_BestPeer(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	cnPeer.EXPECT().Head().Return(common.Hash{}, big.NewInt(111)).Times(2)
	pnPeer.EXPECT().Head().Return(common.Hash{}, big.NewInt(222)).Times(2)
	enPeer.EXPECT().Head().Return(common.Hash{}, big.NewInt(333)).Times(1)
	enPeer.EXPECT().Close().Times(1)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	assert.Nil(t, peerSet.BestPeer())

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	assert.Equal(t, enPeer, peerSet.BestPeer())

	assert.NoError(t, peerSet.Unregister(nodeids[2].String()))

	assert.Equal(t, pnPeer, peerSet.BestPeer())
}

func TestPeerSet_Close(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	pnPeer := NewMockPeer(mockCtrl)
	enPeer := NewMockPeer(mockCtrl)

	cnPeer.EXPECT().DisconnectP2PPeer(p2p.DiscQuitting).Times(1)
	pnPeer.EXPECT().DisconnectP2PPeer(p2p.DiscQuitting).Times(1)
	enPeer.EXPECT().DisconnectP2PPeer(p2p.DiscQuitting).Times(1)

	setMockPeersConnType(cnPeer, pnPeer, enPeer)
	setMockPeers([]*MockPeer{cnPeer, pnPeer, enPeer})

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	assert.False(t, peerSet.closed)
	peerSet.Close()
	assert.True(t, peerSet.closed)
}

func TestPeerSet_SampleResendPeersByType_PN(t *testing.T) {
	// CN Peer=1, PN Peer=0
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		cnPeer := NewMockPeer(mockCtrl)
		cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
		setMockPeers([]*MockPeer{cnPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.PROXYNODE))
		assert.NoError(t, peerSet.Register(cnPeer))
		assert.Equal(t, []Peer{cnPeer}, peerSet.SampleResendPeersByType(node.PROXYNODE))

		mockCtrl.Finish()
	}
	// CN Peer=0, PN Peer=1
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(3)
		setMockPeers([]*MockPeer{pnPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.PROXYNODE))
		assert.NoError(t, peerSet.Register(pnPeer))
		assert.Equal(t, []Peer{pnPeer}, peerSet.SampleResendPeersByType(node.PROXYNODE))

		mockCtrl.Finish()
	}
	// CN Peer=1, PN Peer=1
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		cnPeer := NewMockPeer(mockCtrl)
		pnPeer := NewMockPeer(mockCtrl)

		cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).AnyTimes()
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).AnyTimes()

		setMockPeers([]*MockPeer{cnPeer, pnPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.PROXYNODE))
		assert.NoError(t, peerSet.Register(cnPeer))
		assert.NoError(t, peerSet.Register(pnPeer))
		assert.Equal(t, []Peer{cnPeer}, peerSet.SampleResendPeersByType(node.PROXYNODE))

		mockCtrl.Finish()
	}
	// CN Peer=3, PN Peer=1
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		cnPeer1 := NewMockPeer(mockCtrl)
		cnPeer2 := NewMockPeer(mockCtrl)
		cnPeer3 := NewMockPeer(mockCtrl)
		pnPeer := NewMockPeer(mockCtrl)

		cnPeer1.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).AnyTimes()
		cnPeer2.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).AnyTimes()
		cnPeer3.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).AnyTimes()
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).AnyTimes()

		setMockPeers([]*MockPeer{cnPeer1, cnPeer2, cnPeer3, pnPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.PROXYNODE))
		assert.NoError(t, peerSet.Register(cnPeer1))
		assert.NoError(t, peerSet.Register(cnPeer2))
		assert.NoError(t, peerSet.Register(cnPeer3))
		assert.NoError(t, peerSet.Register(pnPeer))
		resendPeers := peerSet.SampleResendPeersByType(node.PROXYNODE)

		assert.Equal(t, 2, len(resendPeers))
		assert.False(t, containsPeer(pnPeer, resendPeers))

		mockCtrl.Finish()
	}
}

func TestPeerSet_SampleResendPeersByType_EN(t *testing.T) {
	// PN Peer=1, EN Peer=0
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		setMockPeers([]*MockPeer{pnPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.ENDPOINTNODE))
		assert.NoError(t, peerSet.Register(pnPeer))
		assert.Equal(t, []Peer{pnPeer}, peerSet.SampleResendPeersByType(node.ENDPOINTNODE))

		mockCtrl.Finish()
	}
	// PN Peer=0, EN Peer=1
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		enPeer := NewMockPeer(mockCtrl)
		enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(3)
		setMockPeers([]*MockPeer{enPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.ENDPOINTNODE))
		assert.NoError(t, peerSet.Register(enPeer))
		assert.Equal(t, []Peer{enPeer}, peerSet.SampleResendPeersByType(node.ENDPOINTNODE))

		mockCtrl.Finish()
	}
	// PN Peer=1, EN Peer=1
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		pnPeer := NewMockPeer(mockCtrl)
		enPeer := NewMockPeer(mockCtrl)

		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)

		setMockPeers([]*MockPeer{pnPeer, enPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.ENDPOINTNODE))
		assert.NoError(t, peerSet.Register(pnPeer))
		assert.NoError(t, peerSet.Register(enPeer))
		assert.Equal(t, []Peer{pnPeer}, peerSet.SampleResendPeersByType(node.ENDPOINTNODE))

		mockCtrl.Finish()
	}
	// PN Peer=3, EN Peer=1
	{
		peerSet := newPeerSet()
		mockCtrl := gomock.NewController(t)

		pnPeer1 := NewMockPeer(mockCtrl)
		pnPeer2 := NewMockPeer(mockCtrl)
		pnPeer3 := NewMockPeer(mockCtrl)
		enPeer := NewMockPeer(mockCtrl)

		pnPeer1.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).AnyTimes()
		pnPeer2.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).AnyTimes()
		pnPeer3.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).AnyTimes()
		enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).AnyTimes()

		setMockPeers([]*MockPeer{pnPeer1, pnPeer2, pnPeer3, enPeer})

		assert.Equal(t, []Peer{}, peerSet.SampleResendPeersByType(node.ENDPOINTNODE))
		assert.NoError(t, peerSet.Register(pnPeer1))
		assert.NoError(t, peerSet.Register(pnPeer2))
		assert.NoError(t, peerSet.Register(pnPeer3))
		assert.NoError(t, peerSet.Register(enPeer))
		resendPeers := peerSet.SampleResendPeersByType(node.ENDPOINTNODE)

		assert.Equal(t, 2, len(resendPeers))
		assert.False(t, containsPeer(enPeer, resendPeers))

		mockCtrl.Finish()
	}
}

func containsPeer(target Peer, peers []Peer) bool {
	for _, peer := range peers {
		if target == peer {
			return true
		}
	}
	return false
}
