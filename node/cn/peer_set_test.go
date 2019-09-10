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

func TestPeerSet_Register(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(3)
	cnPeer.EXPECT().Broadcast().AnyTimes()

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(3)
	pnPeer.EXPECT().Broadcast().AnyTimes()

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(3)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(3)
	enPeer.EXPECT().Broadcast().AnyTimes()

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	assert.Equal(t, errAlreadyRegistered, peerSet.Register(cnPeer))
	assert.Equal(t, errAlreadyRegistered, peerSet.Register(pnPeer))
	assert.Equal(t, errAlreadyRegistered, peerSet.Register(enPeer))

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
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(4)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()
	cnPeer.EXPECT().Close().Times(1)

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(4)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()
	pnPeer.EXPECT().Close().Times(1)

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(4)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()
	enPeer.EXPECT().Close().Times(1)

	assert.Equal(t, errNotRegistered, peerSet.Unregister(nodeID1.String()))
	assert.Equal(t, errNotRegistered, peerSet.Unregister(nodeID2.String()))
	assert.Equal(t, errNotRegistered, peerSet.Unregister(nodeID3.String()))

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	assert.NoError(t, peerSet.Unregister(nodeID1.String()))
	assert.NoError(t, peerSet.Unregister(nodeID2.String()))
	assert.NoError(t, peerSet.Unregister(nodeID3.String()))
}

func TestPeerSet_Peers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(3)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	peers := peerSet.Peers()
	expectedPeers := map[string]Peer{nodeID1.String(): cnPeer, nodeID2.String(): pnPeer, nodeID3.String(): enPeer}
	assert.EqualValues(t, expectedPeers, peers)
}

func TestPeerSet_CNPeers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()

	assert.EqualValues(t, map[common.Address]Peer{}, peerSet.CNPeers())
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.EqualValues(t, map[common.Address]Peer{addr1: cnPeer}, peerSet.CNPeers())
}

func TestPeerSet_PNPeers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()

	assert.EqualValues(t, map[common.Address]Peer{}, peerSet.PNPeers())
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.EqualValues(t, map[common.Address]Peer{addr1: cnPeer}, peerSet.PNPeers())
}

func TestPeerSet_ENPeers(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()

	assert.EqualValues(t, map[common.Address]Peer{}, peerSet.ENPeers())
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.EqualValues(t, map[common.Address]Peer{addr1: cnPeer}, peerSet.ENPeers())
}

func TestPeerSet_Peer_And_Len(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(3)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()

	assert.Equal(t, 0, peerSet.Len())
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.Equal(t, 1, peerSet.Len())
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.Equal(t, 2, peerSet.Len())
	assert.NoError(t, peerSet.Register(enPeer))
	assert.Equal(t, 3, peerSet.Len())

	assert.Equal(t, cnPeer, peerSet.Peer(nodeID1.String()))
	assert.Equal(t, pnPeer, peerSet.Peer(nodeID2.String()))
	assert.Equal(t, enPeer, peerSet.Peer(nodeID3.String()))
}

func TestPeerSet_PeersWithoutBlock(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	block := newBlock(blockNum1)

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()
	cnPeer.EXPECT().KnowsBlock(block.Hash()).Return(false).Times(2)

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()
	pnPeer.EXPECT().KnowsBlock(block.Hash()).Return(true).Times(2)

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(3)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()
	enPeer.EXPECT().KnowsBlock(block.Hash()).Return(false).Times(2)

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
	tx := types.NewTransaction(111, addr1, big.NewInt(111), 111, big.NewInt(111), nil)

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()
	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()
	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).Times(1)

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(3)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()
	enPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

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
	tx := types.NewTransaction(111, addr1, big.NewInt(111), 111, big.NewInt(111), nil)

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(4)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()
	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(4)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()
	pnPeer.EXPECT().KnowsTx(tx.Hash()).Return(true).Times(1)

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(4)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(3)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()
	enPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

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
	tx := types.NewTransaction(111, addr1, big.NewInt(111), 111, big.NewInt(111), nil)

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()
	cnPeer.EXPECT().KnowsTx(tx.Hash()).Return(false).Times(1)

	assert.EqualValues(t, []Peer{}, peerSet.CNWithoutTx(tx.Hash()))
	assert.NoError(t, peerSet.Register(cnPeer))
	assert.EqualValues(t, []Peer{cnPeer}, peerSet.CNWithoutTx(tx.Hash()))
}

func TestPeerSet_BestPeer(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()
	cnPeer.EXPECT().Head().Return(common.Hash{}, big.NewInt(111)).Times(2)

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()
	pnPeer.EXPECT().Head().Return(common.Hash{}, big.NewInt(222)).Times(2)

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(4)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()
	enPeer.EXPECT().Head().Return(common.Hash{}, big.NewInt(333)).Times(1)
	enPeer.EXPECT().Close().Times(1)

	assert.Nil(t, peerSet.BestPeer())

	assert.NoError(t, peerSet.Register(cnPeer))
	assert.NoError(t, peerSet.Register(pnPeer))
	assert.NoError(t, peerSet.Register(enPeer))

	assert.Equal(t, enPeer, peerSet.BestPeer())

	assert.NoError(t, peerSet.Unregister(nodeID3.String()))

	assert.Equal(t, pnPeer, peerSet.BestPeer())
}

func TestPeerSet_Close(t *testing.T) {
	peerSet := newPeerSet()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cnPeer := NewMockPeer(mockCtrl)
	cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(1)
	cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
	cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
	cnPeer.EXPECT().Broadcast().AnyTimes()
	cnPeer.EXPECT().DisconnectP2PPeer(p2p.DiscQuitting).Times(1)

	pnPeer := NewMockPeer(mockCtrl)
	pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(1)
	pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
	pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
	pnPeer.EXPECT().Broadcast().AnyTimes()
	pnPeer.EXPECT().DisconnectP2PPeer(p2p.DiscQuitting).Times(1)

	enPeer := NewMockPeer(mockCtrl)
	enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(1)
	enPeer.EXPECT().GetAddr().Return(addr3).Times(3)
	enPeer.EXPECT().GetID().Return(nodeID3.String()).Times(2)
	enPeer.EXPECT().Broadcast().AnyTimes()
	enPeer.EXPECT().DisconnectP2PPeer(p2p.DiscQuitting).Times(1)

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
		cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
		cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
		cnPeer.EXPECT().Broadcast().AnyTimes()

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
		pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
		pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
		pnPeer.EXPECT().Broadcast().AnyTimes()

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
		cnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
		cnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
		cnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
		cnPeer.EXPECT().Broadcast().AnyTimes()

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		pnPeer.EXPECT().GetAddr().Return(addr2).Times(3)
		pnPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
		pnPeer.EXPECT().Broadcast().AnyTimes()

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
		cnPeer1.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
		cnPeer1.EXPECT().GetAddr().Return(addr1).Times(3)
		cnPeer1.EXPECT().GetID().Return(nodeID1.String()).Times(2)
		cnPeer1.EXPECT().Broadcast().AnyTimes()

		cnPeer2 := NewMockPeer(mockCtrl)
		cnPeer2.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
		cnPeer2.EXPECT().GetAddr().Return(addr2).Times(3)
		cnPeer2.EXPECT().GetID().Return(nodeID2.String()).Times(2)
		cnPeer2.EXPECT().Broadcast().AnyTimes()

		cnPeer3 := NewMockPeer(mockCtrl)
		cnPeer3.EXPECT().ConnType().Return(p2p.ConnType(node.CONSENSUSNODE)).Times(2)
		cnPeer3.EXPECT().GetAddr().Return(addr3).Times(3)
		cnPeer3.EXPECT().GetID().Return(nodeID3.String()).Times(2)
		cnPeer3.EXPECT().Broadcast().AnyTimes()

		pnPeer := NewMockPeer(mockCtrl)
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		pnPeer.EXPECT().GetAddr().Return(addr4).Times(3)
		pnPeer.EXPECT().GetID().Return(nodeID4.String()).Times(2)
		pnPeer.EXPECT().Broadcast().AnyTimes()

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
		pnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
		pnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
		pnPeer.EXPECT().Broadcast().AnyTimes()

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
		enPeer.EXPECT().GetAddr().Return(addr2).Times(3)
		enPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
		enPeer.EXPECT().Broadcast().AnyTimes()

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
		pnPeer.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		pnPeer.EXPECT().GetAddr().Return(addr1).Times(3)
		pnPeer.EXPECT().GetID().Return(nodeID1.String()).Times(2)
		pnPeer.EXPECT().Broadcast().AnyTimes()

		enPeer := NewMockPeer(mockCtrl)
		enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)
		enPeer.EXPECT().GetAddr().Return(addr2).Times(3)
		enPeer.EXPECT().GetID().Return(nodeID2.String()).Times(2)
		enPeer.EXPECT().Broadcast().AnyTimes()

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
		pnPeer1.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		pnPeer1.EXPECT().GetAddr().Return(addr1).Times(3)
		pnPeer1.EXPECT().GetID().Return(nodeID1.String()).Times(2)
		pnPeer1.EXPECT().Broadcast().AnyTimes()

		pnPeer2 := NewMockPeer(mockCtrl)
		pnPeer2.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		pnPeer2.EXPECT().GetAddr().Return(addr2).Times(3)
		pnPeer2.EXPECT().GetID().Return(nodeID2.String()).Times(2)
		pnPeer2.EXPECT().Broadcast().AnyTimes()

		pnPeer3 := NewMockPeer(mockCtrl)
		pnPeer3.EXPECT().ConnType().Return(p2p.ConnType(node.PROXYNODE)).Times(2)
		pnPeer3.EXPECT().GetAddr().Return(addr3).Times(3)
		pnPeer3.EXPECT().GetID().Return(nodeID3.String()).Times(2)
		pnPeer3.EXPECT().Broadcast().AnyTimes()

		enPeer := NewMockPeer(mockCtrl)
		enPeer.EXPECT().ConnType().Return(p2p.ConnType(node.ENDPOINTNODE)).Times(2)
		enPeer.EXPECT().GetAddr().Return(addr4).Times(3)
		enPeer.EXPECT().GetID().Return(nodeID4.String()).Times(2)
		enPeer.EXPECT().Broadcast().AnyTimes()

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
