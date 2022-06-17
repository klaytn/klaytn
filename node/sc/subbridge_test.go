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

package sc

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/node"
	"github.com/stretchr/testify/assert"
)

// testNewSubBridge returns a test SubBridge.
func testNewSubBridge(t *testing.T) *SubBridge {
	tempDir, err := ioutil.TempDir(os.TempDir(), "klaytn-test-sb-")
	if err != nil {
		t.Fatal(err)
	}

	sCtx := node.NewServiceContext(&node.DefaultConfig, map[reflect.Type]node.Service{}, &event.TypeMux{}, &accounts.Manager{})
	sBridge, err := NewSubBridge(sCtx, &SCConfig{NetworkId: testNetVersion, DataDir: tempDir})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, sBridge)

	return sBridge
}

// removeData removes blockchain data generated during the SubBridge test.
func (sb *SubBridge) removeData(t *testing.T) {
	if err := os.RemoveAll(sb.config.DataDir); err != nil {
		t.Fatal(err)
	}
}

// TestSubBridge_basic tests some getters and basic operation of SubBridge.
func TestSubBridge_basic(t *testing.T) {
	// Create a test SubBridge
	sBridge := testNewSubBridge(t)
	defer sBridge.removeData(t)

	// APIs returns default rpc APIs of MainBridge
	apis := sBridge.APIs()
	assert.Equal(t, 2, len(apis))
	assert.Equal(t, "subbridge", apis[0].Namespace)
	assert.Equal(t, "subbridge", apis[1].Namespace)

	// Test getters for elements of SubBridge
	assert.Equal(t, true, sBridge.IsListening()) // Always returns `true`
	assert.Equal(t, testProtocolVersion, sBridge.ProtocolVersion())
	assert.Equal(t, testNetVersion, sBridge.NetVersion())

	// New components of MainBridge which will update old components
	bc := testBlockChain(t)
	txPool := testTxPool(sBridge.config.DataDir, bc)

	var comp []interface{}
	comp = append(comp, bc)
	comp = append(comp, txPool)

	// Check initial status of components
	assert.Nil(t, sBridge.blockchain)
	assert.Nil(t, sBridge.txPool)

	// Update and check MainBridge components
	sBridge.SetComponents(comp)
	assert.Equal(t, bc, sBridge.blockchain)
	assert.Equal(t, txPool, sBridge.txPool)

	// Start MainBridge and stop later
	if err := sBridge.Start(p2p.SingleChannelServer{}); err != nil {
		t.Fatal(err)
	}
	defer sBridge.Stop()

	//TODO more test
}

// TestSubBridge_removePeer tests correct removal of a peer from `SubBridge.peers`.
func TestSubBridge_removePeer(t *testing.T) {
	// Create a test SubBridge (it may have 0 peers)
	sBridge := testNewSubBridge(t)
	defer sBridge.removeData(t)
	defer sBridge.chainDB.Close()

	// Set components of SubBridge
	bc := testBlockChain(t)
	txPool := testTxPool(sBridge.config.DataDir, bc)

	var comp []interface{}
	comp = append(comp, bc)
	comp = append(comp, txPool)

	sBridge.SetComponents(comp)

	// Prepare information of bridgePeer to be added and removed
	nodeID := "0x1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d"
	peer := p2p.NewPeer(discover.MustHexID(nodeID), "name", []p2p.Cap{})
	peerID := peer.ID()
	bridgePeerID := fmt.Sprintf("%x", peerID[:8])

	// mockBridgePeer mocks BridgePeer
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBridgePeer := NewMockBridgePeer(mockCtrl)
	mockBridgePeer.EXPECT().GetID().Return(bridgePeerID).AnyTimes()
	mockBridgePeer.EXPECT().GetChainID().Return(big.NewInt(0)).AnyTimes()
	mockBridgePeer.EXPECT().GetP2PPeer().Return(peer).AnyTimes()
	mockBridgePeer.EXPECT().SendServiceChainInfoRequest(gomock.Any()).Return(nil).AnyTimes()
	mockBridgePeer.EXPECT().Close().Return().AnyTimes()

	// Add the bridgePeer
	if err := sBridge.peers.Register(mockBridgePeer); err != nil {
		t.Fatal(err)
	}
	if err := sBridge.handler.RegisterNewPeer(mockBridgePeer); err != nil {
		t.Fatal(err)
	}

	peerNum := sBridge.peers.Len()

	// Try to remove a non-registered bridgePeer, and nothing happen
	sBridge.removePeer("0x11111111")
	assert.Equal(t, peerNum, sBridge.peers.Len())

	// Remove the registered bridgePeer
	sBridge.removePeer(mockBridgePeer.GetID())
	assert.Equal(t, peerNum-1, sBridge.peers.Len())
}

// TestSubBridge_handleMsg fails when a bridgePeer fails to read a message or reads a too long message.
func TestSubBridge_handleMsg(t *testing.T) {
	// Create a test SubBridge
	sBridge := testNewSubBridge(t)
	defer sBridge.removeData(t)
	defer sBridge.chainDB.Close()

	// Elements for a bridgePeer
	key, _ := crypto.GenerateKey()
	nodeID := discover.PubkeyID(&key.PublicKey)
	peer := p2p.NewPeer(nodeID, "name", []p2p.Cap{})
	pipe1, pipe2 := p2p.MsgPipe()

	// bridgePeer will receive a message through rw1
	bridgePeer := newBridgePeer(testProtocolVersion, peer, pipe1)

	// Case1. Send a valid message and handle it successfully
	{
		data := "valid message"
		go func() {
			if err := p2p.Send(pipe2, StatusMsg, data); err != nil {
				t.Error(err)
				return
			}
		}()

		if err := sBridge.handleMsg(bridgePeer); err != nil {
			t.Fatal(err)
		}
	}

	// Case2. Send an invalid message having large size and fail to handle
	{
		data := strings.Repeat("a", ProtocolMaxMsgSize+1)
		go func() {
			if err := p2p.Send(pipe2, StatusMsg, data); err != nil {
				t.Error(err)
				return
			}
		}()

		err := sBridge.handleMsg(bridgePeer)
		assert.True(t, strings.HasPrefix(err.Error(), "Message too long"))
	}

	// Case3. Return an error when it fails to read a message
	{
		_ = pipe2.Close()

		err := sBridge.handleMsg(bridgePeer)
		assert.Equal(t, p2p.ErrPipeClosed, err)

	}
	_ = pipe1.Close()
}

// TestSubBridge_handle tests the fail cases of `handle` function.
// There are no success cases in this test since `handle` has a infinite loop inside.
func TestSubBridge_handle(t *testing.T) {
	// Create a test SubBridge
	sBridge := testNewSubBridge(t)
	defer sBridge.removeData(t)
	defer sBridge.chainDB.Close()

	// Set components of SubBridge
	bc := testBlockChain(t)
	txPool := testTxPool(sBridge.config.DataDir, bc)

	var comp []interface{}
	comp = append(comp, bc)
	comp = append(comp, txPool)

	sBridge.SetComponents(comp)

	// Variables will be used as return values of mockBridgePeer
	key, _ := crypto.GenerateKey()
	nodeID := discover.PubkeyID(&key.PublicKey)
	peer := p2p.NewPeer(nodeID, "name", []p2p.Cap{})
	peerID := peer.ID()
	bridgePeerID := fmt.Sprintf("%x", peerID[:8])
	pipe, _ := p2p.MsgPipe()

	// mockBridgePeer mocks BridgePeer
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBridgePeer := NewMockBridgePeer(mockCtrl)
	mockBridgePeer.EXPECT().GetID().Return(bridgePeerID).AnyTimes()
	mockBridgePeer.EXPECT().GetChainID().Return(big.NewInt(0)).AnyTimes()
	mockBridgePeer.EXPECT().GetP2PPeer().Return(peer).AnyTimes()
	mockBridgePeer.EXPECT().GetP2PPeerID().Return(peerID).AnyTimes()
	mockBridgePeer.EXPECT().GetRW().Return(pipe).AnyTimes()
	mockBridgePeer.EXPECT().SendServiceChainInfoRequest(gomock.Any()).Return(nil).AnyTimes()
	mockBridgePeer.EXPECT().Close().Return().AnyTimes()

	// Case 1 - Error if `sBridge.peers.Len()` was equal or bigger than `sBridge.maxPeers`
	{
		// Set maxPeers to make the test fail
		sBridge.maxPeers = sBridge.peers.Len()

		err := sBridge.handle(mockBridgePeer)
		assert.Equal(t, p2p.DiscTooManyPeers, err)
	}
	// Resolve the above failure condition by increasing maxPeers
	sBridge.maxPeers += 5

	// Case 2 - Error if handshake of BridgePeer failed
	{
		// Make handshake fail
		mockBridgePeer.EXPECT().Handshake(sBridge.networkId, sBridge.getChainID(), gomock.Any(), sBridge.blockchain.CurrentHeader().Hash()).Return(p2p.ErrPipeClosed).Times(1)

		err := sBridge.handle(mockBridgePeer)
		assert.Equal(t, p2p.ErrPipeClosed, err)
	}
	// Resolve the above failure condition by making handshake success
	mockBridgePeer.EXPECT().Handshake(sBridge.networkId, sBridge.getChainID(), gomock.Any(), sBridge.blockchain.CurrentHeader().Hash()).Return(nil).AnyTimes()

	// Case 3 - Error when the same peer was registered before
	{
		// Pre-register a peer which will be added again
		sBridge.peers.peers[bridgePeerID] = &baseBridgePeer{}

		err := sBridge.handle(mockBridgePeer)
		assert.Equal(t, errAlreadyRegistered, err)
	}
	// Resolve the above failure condition by deleting the registered peer
	delete(sBridge.peers.peers, bridgePeerID)

	// Case 4 - Error if `sBridge.handleMsg` failed
	{
		// Close of the peer's pipe make `sBridge.handleMsg` fail
		_ = pipe.Close()

		err := sBridge.handle(mockBridgePeer)
		assert.Equal(t, p2p.ErrPipeClosed, err)
	}
}

// TestSubBridge_SendRPCData tests SendRPCResponseData function of SubBridge.
// The function sends RPC response data to SubBridge's peers.
func TestSubBridge_SendRPCData(t *testing.T) {
	// Create a test SubBridge
	sBridge := testNewSubBridge(t)
	defer sBridge.removeData(t)
	defer sBridge.chainDB.Close()

	// Test data used as a parameter of SendResponseRPC function
	data := []byte{0x11, 0x22, 0x33}

	// mockBridgePeer mocks BridgePeer
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBridgePeer := NewMockBridgePeer(mockCtrl)
	mockBridgePeer.EXPECT().GetID().Return("testID").AnyTimes() // for `sBridge.BridgePeerSet().Register(mockBridgePeer)`

	// Register mockBridgePeer as a peer of `sBridge.BridgePeerSet`
	if err := sBridge.BridgePeerSet().Register(mockBridgePeer); err != nil {
		t.Fatal(err)
	}

	// Case 1 - Error if SendResponseRPC of mockBridgePeer failed
	{
		// Make mockBridgePeer return an error
		mockBridgePeer.EXPECT().SendRequestRPC(data).Return(p2p.ErrPipeClosed).Times(1)

		err := sBridge.SendRPCData(data)
		assert.Equal(t, p2p.ErrPipeClosed, err)
	}

	// Case 2 - Success if SendResponseRPC of mockBridgePeer succeeded
	{
		// Make mockBridgePeer return nil
		mockBridgePeer.EXPECT().SendRequestRPC(data).Return(nil).Times(1)

		err := sBridge.SendRPCData(data)
		assert.Equal(t, nil, err)
	}
}
