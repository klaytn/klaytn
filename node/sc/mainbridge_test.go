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
	"math/big"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/consensus/istanbul/backend"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
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

// testBlockChain returns a test BlockChain with initial values
func testBlockChain(t *testing.T) *blockchain.BlockChain {
	db := database.NewMemoryDBManager()

	gov := governance.NewMixedEngine(&params.ChainConfig{
		ChainID:       big.NewInt(2018),
		UnitPrice:     25000000000,
		DeriveShaImpl: 0,
		Istanbul: &params.IstanbulConfig{
			Epoch:          istanbul.DefaultConfig.Epoch,
			ProposerPolicy: uint64(istanbul.DefaultConfig.ProposerPolicy),
			SubGroupSize:   istanbul.DefaultConfig.SubGroupSize,
		},
		Governance: params.GetDefaultGovernanceConfig(params.UseIstanbul),
	}, db)

	prvKey, _ := crypto.GenerateKey()
	engine := backend.New(common.Address{}, istanbul.DefaultConfig, prvKey, db, gov, common.CONSENSUSNODE)

	var genesis *blockchain.Genesis
	genesis = blockchain.DefaultGenesisBlock()
	genesis.BlockScore = big.NewInt(1)
	genesis.Config.Governance = params.GetDefaultGovernanceConfig(params.UseIstanbul)
	genesis.Config.Istanbul = params.GetDefaultIstanbulConfig()
	genesis.Config.UnitPrice = 25 * params.Ston

	chainConfig, _, err := blockchain.SetupGenesisBlock(db, genesis, params.UnusedNetworkId, false, false)
	if _, ok := err.(*params.ConfigCompatError); err != nil && !ok {
		t.Fatal(err)
	}

	bc, err := blockchain.NewBlockChain(db, nil, chainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return bc
}

func testTxPool(dataDir string, bc *blockchain.BlockChain) *blockchain.TxPool {
	blockchain.DefaultTxPoolConfig.Journal = path.Join(dataDir, blockchain.DefaultTxPoolConfig.Journal)
	return blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bc.Config(), bc)
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

// TestMainBridge_basic tests some getters and basic operation of MainBridge.
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
			Service:   api.NewPublicKlayAPI(&cn.CNAPIBackend{}),
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
	if err := mBridge.peers.Register(bridgePeer); err != nil {
		t.Fatal(err)
	}
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
				t.Error(err)
				return
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
				t.Error(err)
				return
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

// TestMainBridge_handle tests the fail cases of `handle` function.
// There are no success cases in this test since `handle` has a infinite loop inside.
func TestMainBridge_handle(t *testing.T) {
	// Create a MainBridge
	mBridge := testNewMainBridge(t)
	defer mBridge.chainDB.Close()

	// Set testBlockChain to MainBridge.blockchain
	mBridge.blockchain = testBlockChain(t)

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
	mockBridgePeer.EXPECT().GetP2PPeer().Return(peer).AnyTimes()
	mockBridgePeer.EXPECT().GetP2PPeerID().Return(peerID).AnyTimes()
	mockBridgePeer.EXPECT().GetRW().Return(pipe).AnyTimes()
	mockBridgePeer.EXPECT().Close().Return().AnyTimes()

	// Case 1 - Error if `mBridge.peers.Len()` was equal or bigger than `mBridge.maxPeers`
	{
		// Set maxPeers to make the test fail
		mBridge.maxPeers = mBridge.peers.Len()

		err := mBridge.handle(mockBridgePeer)
		assert.Equal(t, p2p.DiscTooManyPeers, err)
	}
	// Resolve the above failure condition by increasing maxPeers
	mBridge.maxPeers += 5

	// Case 2 - Error if handshake of BridgePeer failed
	{
		// Make handshake fail
		mockBridgePeer.EXPECT().Handshake(mBridge.networkId, mBridge.getChainID(), gomock.Any(), mBridge.blockchain.CurrentHeader().Hash()).Return(p2p.ErrPipeClosed).Times(1)

		err := mBridge.handle(mockBridgePeer)
		assert.Equal(t, p2p.ErrPipeClosed, err)
	}
	// Resolve the above failure condition by making handshake success
	mockBridgePeer.EXPECT().Handshake(mBridge.networkId, mBridge.getChainID(), gomock.Any(), mBridge.blockchain.CurrentHeader().Hash()).Return(nil).AnyTimes()

	// Case 3 - Error when the same peer was registered before
	{
		// Pre-register a peer which will be added again
		mBridge.peers.peers[bridgePeerID] = &baseBridgePeer{}

		err := mBridge.handle(mockBridgePeer)
		assert.Equal(t, errAlreadyRegistered, err)
	}
	// Resolve the above failure condition by deleting the registered peer
	delete(mBridge.peers.peers, bridgePeerID)

	// Case 4 - Error if `mBridge.handleMsg` failed
	{
		// Close of the peer's pipe make `mBridge.handleMsg` fail
		_ = pipe.Close()

		err := mBridge.handle(mockBridgePeer)
		assert.Equal(t, p2p.ErrPipeClosed, err)
	}
}

// TestMainBridge_SendRPCResponseData tests SendRPCResponseData function of MainBridge.
// The function sends RPC response data to MainBridge's peers.
func TestMainBridge_SendRPCResponseData(t *testing.T) {
	// Create a MainBridge
	mBridge := testNewMainBridge(t)
	defer mBridge.chainDB.Close()

	// Test data used as a parameter of SendResponseRPC function
	data := []byte{0x11, 0x22, 0x33}

	// mockBridgePeer mocks BridgePeer
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBridgePeer := NewMockBridgePeer(mockCtrl)
	mockBridgePeer.EXPECT().GetID().Return("testID").AnyTimes() // for `mBridge.BridgePeerSet().Register(mockBridgePeer)`

	// Register mockBridgePeer as a peer of `mBridge.BridgePeerSet`
	if err := mBridge.BridgePeerSet().Register(mockBridgePeer); err != nil {
		t.Fatal(err)
	}

	// Case 1 - Error if SendResponseRPC of mockBridgePeer failed
	{
		// Make mockBridgePeer return an error
		mockBridgePeer.EXPECT().SendResponseRPC(data).Return(p2p.ErrPipeClosed).Times(1)

		err := mBridge.SendRPCResponseData(data)
		assert.Equal(t, p2p.ErrPipeClosed, err)
	}

	// Case 2 - Success if SendResponseRPC of mockBridgePeer succeeded
	{
		// Make mockBridgePeer return nil
		mockBridgePeer.EXPECT().SendResponseRPC(data).Return(nil).Times(1)

		err := mBridge.SendRPCResponseData(data)
		assert.Equal(t, nil, err)
	}
}
