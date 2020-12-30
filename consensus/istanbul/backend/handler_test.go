// Copyright 2020 The klaytn Authors
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

package backend

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/rlp"
)

func TestBackend_HandleMsg(t *testing.T) {
	_, backend := newBlockChain(1)
	defer backend.Stop()
	eventSub := backend.istanbulEventMux.Subscribe(istanbul.MessageEvent{})

	addr := common.StringToAddress("test addr")
	data := &istanbul.ConsensusMsg{
		PrevHash: common.HexToHash("0x1234"),
		Payload:  []byte("test data"),
	}
	hash := istanbul.RLPHash(data.Payload)
	size, payload, _ := rlp.EncodeToReader(data)

	// Success case
	{
		msg := p2p.Msg{
			Code:    IstanbulMsg,
			Size:    uint32(size),
			Payload: payload,
		}
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Nil(t, err)
		assert.True(t, isHandled)

		if err != nil {
			t.Fatalf("handle message failed: %v", err)
		}

		recentMsg, ok := backend.recentMessages.Get(addr)
		assert.True(t, ok)

		cachedMsg, ok := recentMsg.(*lru.ARCCache)
		assert.True(t, ok)

		value, ok := cachedMsg.Get(hash)
		assert.True(t, ok)
		assert.True(t, value.(bool))

		value, ok = backend.knownMessages.Get(hash)
		assert.True(t, ok)
		assert.True(t, value.(bool))

		evTimer := time.NewTimer(3 * time.Second)
		defer evTimer.Stop()

		select {
		case event := <-eventSub.Chan():
			switch ev := event.Data.(type) {
			case istanbul.MessageEvent:
				assert.Equal(t, data.Payload, ev.Payload)
				assert.Equal(t, data.PrevHash, ev.Hash)
			default:
				t.Fatal("unexpected message type")
			}
		case <-evTimer.C:
			t.Fatal("failed to subscribe istanbul message event")
		}
	}

	// Failure case - undefined message code
	{
		msg := p2p.Msg{
			Code:    0x99,
			Size:    uint32(size),
			Payload: payload,
		}
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Equal(t, nil, err)
		assert.False(t, isHandled)
	}

	// Failure case - invalid message data
	{
		size, payload, _ := rlp.EncodeToReader([]byte{0x1, 0x2})
		msg := p2p.Msg{
			Code:    IstanbulMsg,
			Size:    uint32(size),
			Payload: payload,
		}
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Equal(t, errDecodeFailed, err)
		assert.True(t, isHandled)
	}

	// Failure case - stopped istanbul engine
	{
		msg := p2p.Msg{
			Code:    IstanbulMsg,
			Size:    uint32(size),
			Payload: payload,
		}
		_ = backend.Stop()
		isHandled, err := backend.HandleMsg(addr, msg)
		assert.Equal(t, istanbul.ErrStoppedEngine, err)
		assert.True(t, isHandled)
	}
}

func TestBackend_Protocol(t *testing.T) {
	backend := newTestBackend()
	assert.Equal(t, istanbulProtocol, backend.Protocol())
}

func TestBackend_ValidatePeerType(t *testing.T) {
	_, backend := newBlockChain(1)
	defer backend.Stop()

	// Return nil if the input address is a validator
	{
		err := backend.ValidatePeerType(backend.address)
		assert.Nil(t, err)
	}

	// Return an error if the input address is invalid
	{
		err := backend.ValidatePeerType(common.Address{})
		assert.Equal(t, errInvalidPeerAddress, err)
	}

	// Return an error if backend.chain is not set
	{
		backend.chain = nil
		err := backend.ValidatePeerType(backend.address)
		assert.Equal(t, errNoChainReader, err)
	}
}
