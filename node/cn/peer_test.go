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
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var version = 63

func newBasePeer() (Peer, *p2p.MsgPipeRW, *p2p.MsgPipeRW) {
	pipe1, pipe2 := p2p.MsgPipe()

	return newPeer(version, p2pPeers[0], pipe1), pipe1, pipe2
}

func TestBasePeer_AddToKnownBlocks(t *testing.T) {
	basePeer, _, _ := newBasePeer()
	assert.False(t, basePeer.KnowsBlock(hash1))
	basePeer.AddToKnownBlocks(hash1)
	assert.True(t, basePeer.KnowsBlock(hash1))
}

func TestBasePeer_AddToKnownTxs(t *testing.T) {
	basePeer, _, _ := newBasePeer()
	assert.False(t, basePeer.KnowsTx(hash1))
	basePeer.AddToKnownTxs(hash1)
	assert.True(t, basePeer.KnowsTx(hash1))
}

func TestBasePeer_Send(t *testing.T) {
	basePeer, _, oppositePipe := newBasePeer()
	data := "a message data"
	expectedMsg := generateMsg(t, NewBlockHashesMsg, data)
	go func(t *testing.T) {
		if err := basePeer.Send(NewBlockHashesMsg, data); err != nil {
			t.Fatal(err)
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedMsg.Code, receivedMsg.Code)
	assert.Equal(t, expectedMsg.Size, receivedMsg.Size)

	var decodedStr string
	if err := receivedMsg.Decode(&decodedStr); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, data, decodedStr)
}

func TestBasePeer_SendTransactions(t *testing.T) {
	sentTxs := types.Transactions{tx1}

	basePeer, _, oppositePipe := newBasePeer()
	assert.False(t, basePeer.KnowsTx(sentTxs[0].Hash()))
	go func(t *testing.T) {
		if err := basePeer.SendTransactions(sentTxs); err != nil {
			t.Fatal(t)
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedTxs types.Transactions
	if err := receivedMsg.Decode(&receivedTxs); err != nil {
		t.Fatal(err)
	}

	assert.True(t, basePeer.KnowsTx(tx1.Hash()))
	assert.Equal(t, len(sentTxs), len(receivedTxs))
	receivedTxs[0].Hash()
	assert.Equal(t, sentTxs[0], receivedTxs[0])
}

func TestBasePeer_ReSendTransactions(t *testing.T) {
	sentTxs := types.Transactions{tx1}

	basePeer, _, oppositePipe := newBasePeer()
	assert.False(t, basePeer.KnowsTx(sentTxs[0].Hash()))
	go func(t *testing.T) {
		if err := basePeer.ReSendTransactions(sentTxs); err != nil {
			t.Fatal(t)
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedTxs types.Transactions
	if err := receivedMsg.Decode(&receivedTxs); err != nil {
		t.Fatal(err)
	}

	assert.False(t, basePeer.KnowsTx(tx1.Hash()))
	assert.Equal(t, len(sentTxs), len(receivedTxs))
	receivedTxs[0].Hash()
	assert.Equal(t, sentTxs[0], receivedTxs[0])
}

func TestBasePeer_AsyncSendTransactions(t *testing.T) {
	sentTxs := types.Transactions{tx1}
	lastTxs := types.Transactions{types.NewTransaction(333, addrs[0], big.NewInt(333), 333, big.NewInt(333), addrs[0][:])}

	basePeer, _, _ := newBasePeer()

	// To queuedTxs be filled with transactions
	for i := 0; i < maxQueuedTxs; i++ {
		basePeer.AsyncSendTransactions(sentTxs)
	}
	// lastTxs shouldn't go into the queuedTxs
	basePeer.AsyncSendTransactions(lastTxs)

	assert.True(t, basePeer.KnowsTx(tx1.Hash()))
	assert.False(t, basePeer.KnowsTx(lastTxs[0].Hash()))
}
