// Copyright 2022 The klaytn Authors
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
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

// genRandomString returns ASCII string
func genRandomString(t *testing.T) string {
	const length = 10
	result := "v"
	for {
		if len(result) >= length {
			return result

		}
		num, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(127)))
		assert.Nil(t, err)

		n := num.Int64()
		if n > 32 && n < 127 {
			result += strconv.FormatInt(n, 16)
		}
	}
}

// genRandomNum returns random string of float number
func genRandomNum() string {
	return fmt.Sprintf("%f", rand.Float64()*10000)
}

// makePeer creates a new peer.
func makePeer() *p2p.Peer {
	key, _ := crypto.GenerateKey()
	nodeID := discover.PubkeyID(&key.PublicKey)
	return p2p.NewPeer(nodeID, "name", []p2p.Cap{})
}

// getHandshakeRet runs handsake between two peers and return the result of it.
func getHandshakeRet(bridgePeer1, bridgePeer2 BridgePeer) (error, error) {
	var (
		errHandshakePeer1 error
		errHandshakePeer2 error
	)
	networkID := uint64(8888)
	chainID := big.NewInt(9999)
	td := big.NewInt(1)
	hash := common.Hash{}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		errHandshakePeer1 = bridgePeer1.Handshake(networkID, chainID, td, hash)
	}()
	go func() {
		defer wg.Done()
		errHandshakePeer2 = bridgePeer2.Handshake(networkID, chainID, td, hash)
	}()
	wg.Wait()
	return errHandshakePeer1, errHandshakePeer2
}

// compareVersion checks version matching based on property-based testing.
func compareVersion(t *testing.T, pv1, pv2, nv1, nv2 string, errHandshakePeer1, errHandshakePeer2 error) {
	if pv1 == pv2 && nv1 == nv2 {
		assert.Nil(t, errHandshakePeer1)
		assert.Nil(t, errHandshakePeer2)
	} else {
		assert.NotNil(t, errHandshakePeer1)
		assert.NotNil(t, errHandshakePeer2)
	}
}

// TestVersionCompare tests version correctness. Both peers must have same version of SC protocol, kscn binary version, and contract version.
func TestVersionCompare(t *testing.T) {
	// Property-based testing
	peer1, peer2 := makePeer(), makePeer()
	pipe1, pipe2 := p2p.MsgPipe()

	{
		// Case 1 - Success (Exactly matched)
		testProtocolVersion1, testProtocolVersion2 := "1.0", "1.0"
		testNodeVersion1, testNodeVersion2 := "v1.8.0", "v1.8.0"
		bridgePeer1 := newBridgePeer(testProtocolVersion1, testNodeVersion1, peer1, pipe1)
		bridgePeer2 := newBridgePeer(testProtocolVersion2, testNodeVersion2, peer2, pipe2)

		errHandshakePeer1, errHandshakePeer2 := getHandshakeRet(bridgePeer1, bridgePeer2)
		compareVersion(t,
			testProtocolVersion1, testProtocolVersion2,
			testNodeVersion1, testNodeVersion2,
			errHandshakePeer1, errHandshakePeer2)
	}

	{
		// Case 2 - Failure (ProtocolVersion diff)
		testProtocolVersion1, testProtocolVersion2 := "1.0", "1.1"
		testNodeVersion1, testNodeVersion2 := "v1.8.0", "v1.8.0"
		bridgePeer1 := newBridgePeer(testProtocolVersion1, testNodeVersion1, peer1, pipe1)
		bridgePeer2 := newBridgePeer(testProtocolVersion2, testNodeVersion2, peer2, pipe2)

		errHandshakePeer1, errHandshakePeer2 := getHandshakeRet(bridgePeer1, bridgePeer2)
		compareVersion(t,
			testProtocolVersion1, testProtocolVersion2,
			testNodeVersion1, testNodeVersion2,
			errHandshakePeer1, errHandshakePeer2)
	}

	{
		// Case 3 - Failure (NodeVersion diff)
		testProtocolVersion1, testProtocolVersion2 := "1.0", "1.0"
		testNodeVersion1, testNodeVersion2 := "v1.8.0", "v1.8.1"
		bridgePeer1 := newBridgePeer(testProtocolVersion1, testNodeVersion1, peer1, pipe1)
		bridgePeer2 := newBridgePeer(testProtocolVersion2, testNodeVersion2, peer2, pipe2)

		errHandshakePeer1, errHandshakePeer2 := getHandshakeRet(bridgePeer1, bridgePeer2)
		compareVersion(t,
			testProtocolVersion1, testProtocolVersion2,
			testNodeVersion1, testNodeVersion2,
			errHandshakePeer1, errHandshakePeer2)
	}

	{
		// Case 4 - Failure (All diff)
		testProtocolVersion1, testProtocolVersion2 := "1.0", "1.1"
		testNodeVersion1, testNodeVersion2 := "v1.8.0", "v1.8.1"
		bridgePeer1 := newBridgePeer(testProtocolVersion1, testNodeVersion1, peer1, pipe1)
		bridgePeer2 := newBridgePeer(testProtocolVersion2, testNodeVersion2, peer2, pipe2)

		errHandshakePeer1, errHandshakePeer2 := getHandshakeRet(bridgePeer1, bridgePeer2)
		compareVersion(t,
			testProtocolVersion1, testProtocolVersion2,
			testNodeVersion1, testNodeVersion2,
			errHandshakePeer1, errHandshakePeer2)
	}

	{
		// Case 5 - ? (Random)
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 100; i++ {
			testProtocolVersion1, testProtocolVersion2 := genRandomNum(), genRandomNum()
			testNodeVersion1, testNodeVersion2 := genRandomString(t), genRandomString(t)
			bridgePeer1 := newBridgePeer(testProtocolVersion1, testNodeVersion1, peer1, pipe1)
			bridgePeer2 := newBridgePeer(testProtocolVersion2, testNodeVersion2, peer2, pipe2)

			errHandshakePeer1, errHandshakePeer2 := getHandshakeRet(bridgePeer1, bridgePeer2)
			compareVersion(t,
				testProtocolVersion1, testProtocolVersion2,
				testNodeVersion1, testNodeVersion2,
				errHandshakePeer1, errHandshakePeer2)
		}
	}
}

// makeBackend makes backand and returns it
func makeBackend(bridgeAccount *bind.TransactOpts) *backends.SimulatedBackend {
	alloc := blockchain.GenesisAlloc{bridgeAccount.From: {Balance: big.NewInt(params.KLAY)}}
	backend := backends.NewSimulatedBackend(alloc)
	return backend
}

// deployBridgeContract deploys bridge contract
func deployBridgeContract(t *testing.T, bridgeAccount *bind.TransactOpts, backend *backends.SimulatedBackend) *bridge.Bridge {
	_, tx, b, err := bridge.DeployBridge(bridgeAccount, backend, false)
	assert.NoError(t, err)
	backend.Commit()
	assert.Nil(t, WaitMined(tx, backend, t))
	return b
}

// deployAnotherBridgeContract deploys another bridge contract
func deployAnotherBridgeContract(t *testing.T, bridgeAccount *bind.TransactOpts, backend *backends.SimulatedBackend) *bridge.BridgeAnotherVersion {
	_, tx, b, err := bridge.DeployBridgeAnotherVersion(bridgeAccount, backend)
	assert.Nil(t, err)
	backend.Commit()
	assert.Nil(t, WaitMined(tx, backend, t))
	return b
}

// getBridgeContractVersion queries a contract version and returns it
func getBridgeContractVersion(t *testing.T, bridgeContract *bridge.Bridge) uint64 {
	version, err := bridgeContract.VERSION(nil)
	assert.Nil(t, err)
	return version
}

// getAnotherBridgeContractVersion queries a contract version and returns it
func getAnotherBridgeContractVersion(t *testing.T, bridgeContract *bridge.BridgeAnotherVersion) uint64 {
	version, err := bridgeContract.VERSION(nil)
	assert.Nil(t, err)
	return version
}

// TestCompareBridgeContractVersion compares version of two contracts
func TestCompareBridgeContractVersion(t *testing.T) {
	bridgeAccountKey, _ := crypto.GenerateKey()
	bridgeAccount := bind.NewKeyedTransactor(bridgeAccountKey)
	backend := makeBackend(bridgeAccount)
	defer backend.Close()

	{
		anotherBridgeContract := deployAnotherBridgeContract(t, bridgeAccount, backend)
		bridgeContract := deployBridgeContract(t, bridgeAccount, backend)

		// Case 1 - Success (matched with same contracts)
		assert.Equal(t, getBridgeContractVersion(t, bridgeContract), getBridgeContractVersion(t, bridgeContract))
		assert.Equal(t, getAnotherBridgeContractVersion(t, anotherBridgeContract), getAnotherBridgeContractVersion(t, anotherBridgeContract))

		// Case 2 - Failure (not matched)
		assert.NotEqual(t, getBridgeContractVersion(t, bridgeContract), getAnotherBridgeContractVersion(t, anotherBridgeContract))
	}
}
