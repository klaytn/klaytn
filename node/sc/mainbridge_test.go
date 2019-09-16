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
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

const testNetVersion = uint64(8888)

var testProtocolVersion = int(SCProtocolVersion[0])

// testNewMainBridge returns a test MainBridge.
func testNewMainBridge(t *testing.T) *MainBridge {
	sCtx := node.NewServiceContext(&node.DefaultConfig, map[reflect.Type]node.Service{}, &event.TypeMux{}, &accounts.Manager{})
	mBridge, err := NewMainBridge(sCtx, &SCConfig{NetworkId: testNetVersion})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, mBridge)

	return mBridge
}

// TestCreateDB tests creation of chain database and proper working of database operation.
func TestCreateDB(t *testing.T) {
	sCtx := node.NewServiceContext(&node.DefaultConfig, map[reflect.Type]node.Service{}, &event.TypeMux{}, &accounts.Manager{})
	sConfig := &SCConfig{}
	name := "testDB"

	// Create a DB Manager
	dbManager := CreateDB(sCtx, sConfig, name)
	defer dbManager.Close()
	assert.NotNil(t, dbManager)

	// Check initial DBConfig of `CreateDB()`
	dbConfig := dbManager.GetDBConfig()
	assert.True(t, strings.HasSuffix(dbConfig.Dir, name))
	assert.Equal(t, database.LevelDB, dbConfig.DBType)
}

// TestMainBridge tests some getters and basic operation of MainBridge.
func TestMainBridge_basic(t *testing.T) {
	// Create a test MainBridge
	mBridge := testNewMainBridge(t)

	// APIs returns default rpc APIs of MainBridge
	apis := mBridge.APIs()
	assert.Equal(t, 2, len(apis))
	assert.Equal(t, "mainbridge", apis[0].Namespace)
	assert.Equal(t, "mainbridge", apis[1].Namespace)

	// Test getters for elements of MainBridge
	assert.Equal(t, true, mBridge.IsListening()) // Always returns `true`
	assert.Equal(t, testProtocolVersion, mBridge.ProtocolVersion())
	assert.Equal(t, testNetVersion, mBridge.NetVersion())

	// New components of MainBridge which will update old components
	bc := &blockchain.BlockChain{}
	txPool := &blockchain.TxPool{}
	compAPIs := []rpc.API{
		{
			Namespace: "klay",
			Version:   "1.0",
			Service:   api.NewPublicKlayAPI(&cn.ServiceChainAPIBackend{}),
			Public:    true,
		},
	}
	var comp []interface{}
	comp = append(comp, bc)
	comp = append(comp, txPool)
	comp = append(comp, compAPIs)

	// Check initial status of components
	assert.Nil(t, mBridge.blockchain)
	assert.Nil(t, mBridge.txPool)
	assert.Nil(t, mBridge.rpcServer.GetServices()["klay"])

	// Update and check MainBridge components
	mBridge.SetComponents(comp)
	assert.Equal(t, bc, mBridge.blockchain)
	assert.Equal(t, txPool, mBridge.txPool)
	assert.NotNil(t, mBridge.rpcServer.GetServices()["klay"])

	// Start MainBridge and stop later
	if err := mBridge.Start(p2p.SingleChannelServer{}); err != nil {
		t.Fatal(err)
	}
	defer mBridge.Stop()

	//TODO more test
}

// TestMainBridge_removePeer tests correct removal of a peer from `MainBridge.peers`.
func TestMainBridge_removePeer(t *testing.T) {
	// Create a MainBridge (it may have 0 peers)
	mBridge := testNewMainBridge(t)
	defer mBridge.chainDB.Close()

	// Prepare a bridgePeer to be added and removed
	nodeID := "0x1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d"
	peer := p2p.NewPeer(discover.MustHexID(nodeID), "name", []p2p.Cap{})
	bridgePeer := mBridge.newPeer(1, peer, &p2p.MsgPipeRW{})

	// Add the bridgePeer
	mBridge.peers.Register(bridgePeer)
	peerNum := mBridge.peers.Len()

	// Try to remove a non-registered bridgePeer, and nothing happen
	mBridge.removePeer("0x11111111")
	assert.Equal(t, peerNum, mBridge.peers.Len())

	// Remove the registered bridgePeer
	mBridge.removePeer(bridgePeer.GetID())
	assert.Equal(t, peerNum-1, mBridge.peers.Len())
}

// TestMainBridge_handleMsg fails when a bridgePeer fails to read a message or reads a too long message.
func TestMainBridge_handleMsg(t *testing.T) {
	// Create a MainBridge
	mBridge := testNewMainBridge(t)
	defer mBridge.chainDB.Close()

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
				t.Fatal(err)
			}
		}()

		if err := mBridge.handleMsg(bridgePeer); err != nil {
			t.Fatal(err)
		}
	}

	// Case2. Send an invalid message having large size and fail to handle
	{
		data := strings.Repeat("a", ProtocolMaxMsgSize+1)
		go func() {
			if err := p2p.Send(pipe2, StatusMsg, data); err != nil {
				t.Fatal(err)
			}
		}()

		err := mBridge.handleMsg(bridgePeer)
		assert.True(t, strings.HasPrefix(err.Error(), "Message too long"))
	}

	// Case3. Return an error when it fails to read a message
	{
		_ = pipe2.Close()

		err := mBridge.handleMsg(bridgePeer)
		assert.Equal(t, p2p.ErrPipeClosed, err)

	}
	_ = pipe1.Close()
}
