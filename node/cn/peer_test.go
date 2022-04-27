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
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/crypto"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/stretchr/testify/assert"
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
			t.Error(err)
			return
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
			t.Error(t)
			return
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
	assert.Equal(t, sentTxs[0].Hash(), receivedTxs[0].Hash())

	sendTxBinary, _ := sentTxs[0].MarshalBinary()
	receivedTxBinary, _ := receivedTxs[0].MarshalBinary()
	assert.Equal(t, sendTxBinary, receivedTxBinary)
}

func TestBasePeer_ReSendTransactions(t *testing.T) {
	sentTxs := types.Transactions{tx1}

	basePeer, _, oppositePipe := newBasePeer()
	assert.False(t, basePeer.KnowsTx(sentTxs[0].Hash()))
	go func(t *testing.T) {
		if err := basePeer.ReSendTransactions(sentTxs); err != nil {
			t.Error(t)
			return
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
	assert.Equal(t, sentTxs[0].Hash(), receivedTxs[0].Hash())

	sendTxBinary, _ := sentTxs[0].MarshalBinary()
	receivedTxBinary, _ := receivedTxs[0].MarshalBinary()
	assert.Equal(t, sendTxBinary, receivedTxBinary)

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

func TestBasePeer_ConnType(t *testing.T) {
	basePeer, _, _ := newBasePeer()
	assert.Equal(t, common.CONSENSUSNODE, basePeer.ConnType())
}

func TestBasePeer_GetAndSetAddr(t *testing.T) {
	basePeer, _, _ := newBasePeer()
	assert.Equal(t, common.Address{}, basePeer.GetAddr())
	basePeer.SetAddr(addrs[0])
	assert.Equal(t, addrs[0], basePeer.GetAddr())
	basePeer.SetAddr(addrs[1])
	assert.Equal(t, addrs[1], basePeer.GetAddr())
}

func TestBasePeer_GetVersion(t *testing.T) {
	basePeer, _, _ := newBasePeer()
	assert.Equal(t, version, basePeer.GetVersion())
}

func TestBasePeer_RegisterConsensusMsgCode(t *testing.T) {
	basePeer, _, _ := newBasePeer()
	assert.True(t, strings.Contains(basePeer.RegisterConsensusMsgCode(NewBlockHashesMsg).Error(), errNotSupportedByBasePeer.Error()))
}

func TestBasePeer_GetRW(t *testing.T) {
	basePeer, pipe1, _ := newBasePeer()
	assert.Equal(t, pipe1, basePeer.GetRW())
}

func TestBasePeer_SendBlockHeaders(t *testing.T) {
	header1 := blocks[0].Header()
	header2 := blocks[1].Header()

	sentHeaders := []*types.Header{header1, header2}

	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.SendBlockHeaders(sentHeaders); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedHeaders []*types.Header
	if err := receivedMsg.Decode(&receivedHeaders); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentHeaders, receivedHeaders)
}

func TestBasePeer_SendFetchedBlockHeader(t *testing.T) {
	sentHeader := blocks[0].Header()

	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.SendFetchedBlockHeader(sentHeader); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedHeader types.Header
	if err := receivedMsg.Decode(&receivedHeader); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentHeader, &receivedHeader)
}

func TestBasePeer_SendNodeData(t *testing.T) {
	sentData := [][]byte{hashes[0][:], hashes[1][:]}

	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.SendNodeData(sentData); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedData [][]byte
	if err := receivedMsg.Decode(&receivedData); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentData, receivedData)
}

func TestBasePeer_FetchBlockHeader(t *testing.T) {
	sentHash := hashes[0]

	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.FetchBlockHeader(sentHash); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedHash common.Hash
	if err := receivedMsg.Decode(&receivedHash); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentHash, receivedHash)
}

func TestBasePeer_RequestBodies(t *testing.T) {
	sentHashes := hashes
	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.RequestBodies(sentHashes); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedHashes []common.Hash
	if err := receivedMsg.Decode(&receivedHashes); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentHashes, receivedHashes)
}

func TestBasePeer_FetchBlockBodies(t *testing.T) {
	sentHashes := hashes
	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.FetchBlockBodies(sentHashes); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedHashes []common.Hash
	if err := receivedMsg.Decode(&receivedHashes); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentHashes, receivedHashes)
}

func TestBasePeer_RequestNodeData(t *testing.T) {
	sentHashes := hashes
	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.RequestNodeData(sentHashes); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedHashes []common.Hash
	if err := receivedMsg.Decode(&receivedHashes); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentHashes, receivedHashes)
}

func TestBasePeer_RequestReceipts(t *testing.T) {
	sentHashes := hashes
	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.RequestReceipts(sentHashes); err != nil {
			t.Error(err)
			return
		}
	}(t)
	receivedMsg, err := oppositePipe.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedHashes []common.Hash
	if err := receivedMsg.Decode(&receivedHashes); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sentHashes, receivedHashes)
}

func TestBasePeer_SendTransactionWithSortedByTime(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := types.LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions.
	txs := types.Transactions{}
	for _, key := range keys {
		tx, _ := types.SignTx(types.NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)

		txs = append(txs, tx)
	}

	// Shuffle transactions.
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(txs), func(i, j int) {
		txs[i], txs[j] = txs[j], txs[i]
	})

	sortedTxs := make(types.Transactions, len(txs))
	copy(sortedTxs, txs)

	// Sort transaction by time.
	sort.Sort(types.TxByPriceAndTime(sortedTxs))

	basePeer, _, oppositePipe := newBasePeer()
	for _, tx := range txs {
		assert.False(t, basePeer.KnowsTx(tx.Hash()))
	}

	go func(t *testing.T) {
		if err := basePeer.SendTransactions(txs); err != nil {
			t.Error(t)
			return
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

	assert.Equal(t, len(txs), len(receivedTxs))

	// It should be received transaction with sorted by times.
	for i, tx := range receivedTxs {
		assert.True(t, basePeer.KnowsTx(tx.Hash()))
		assert.Equal(t, sortedTxs[i].Hash(), tx.Hash())
		assert.False(t, sortedTxs[i].Time().Equal(tx.Time()))
	}
}

func TestBasePeer_ReSendTransactionWithSortedByTime(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := types.LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions.
	txs := types.Transactions{}
	for _, key := range keys {
		tx, _ := types.SignTx(types.NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)

		txs = append(txs, tx)
	}

	// Shuffle transactions.
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(txs), func(i, j int) {
		txs[i], txs[j] = txs[j], txs[i]
	})

	sortedTxs := make(types.Transactions, len(txs))
	copy(sortedTxs, txs)

	// Sort transaction by time.
	sort.Sort(types.TxByPriceAndTime(sortedTxs))

	basePeer, _, oppositePipe := newBasePeer()
	go func(t *testing.T) {
		if err := basePeer.ReSendTransactions(txs); err != nil {
			t.Error(t)
			return
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

	assert.Equal(t, len(txs), len(receivedTxs))

	// It should be received transaction with sorted by times.
	for i, tx := range receivedTxs {
		assert.Equal(t, sortedTxs[i].Hash(), tx.Hash())
		assert.False(t, sortedTxs[i].Time().Equal(tx.Time()))
	}
}

func TestMultiChannelPeer_SendTransactionWithSortedByTime(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := types.LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions.
	txs := types.Transactions{}
	for _, key := range keys {
		tx, _ := types.SignTx(types.NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)

		txs = append(txs, tx)
	}

	// Shuffle transactions.
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(txs), func(i, j int) {
		txs[i], txs[j] = txs[j], txs[i]
	})

	sortedTxs := make(types.Transactions, len(txs))
	copy(sortedTxs, txs)

	// Sort transaction by time.
	sort.Sort(types.TxByPriceAndTime(sortedTxs))

	_, oppositePipe1, oppositePipe2 := newBasePeer()
	multiPeer, _ := newPeerWithRWs(version, p2pPeers[0], []p2p.MsgReadWriter{oppositePipe1, oppositePipe2})

	for _, tx := range txs {
		assert.False(t, multiPeer.KnowsTx(tx.Hash()))
	}

	go func(t *testing.T) {
		if err := multiPeer.SendTransactions(txs); err != nil {
			t.Error(t)
			return
		}
	}(t)

	receivedMsg, err := oppositePipe1.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedTxs types.Transactions
	if err := receivedMsg.Decode(&receivedTxs); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(txs), len(receivedTxs))

	// It should be received transaction with sorted by times.
	for i, tx := range receivedTxs {
		assert.True(t, multiPeer.KnowsTx(tx.Hash()))
		assert.Equal(t, sortedTxs[i].Hash(), tx.Hash())
		assert.False(t, sortedTxs[i].Time().Equal(tx.Time()))
	}
}

func TestMultiChannelPeer_ReSendTransactionWithSortedByTime(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := types.LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions.
	txs := types.Transactions{}
	for _, key := range keys {
		tx, _ := types.SignTx(types.NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)

		txs = append(txs, tx)
	}

	// Shuffle transactions.
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(txs), func(i, j int) {
		txs[i], txs[j] = txs[j], txs[i]
	})

	sortedTxs := make(types.Transactions, len(txs))
	copy(sortedTxs, txs)

	// Sort transaction by time.
	sort.Sort(types.TxByPriceAndTime(sortedTxs))

	_, oppositePipe1, oppositePipe2 := newBasePeer()
	multiPeer, _ := newPeerWithRWs(version, p2pPeers[0], []p2p.MsgReadWriter{oppositePipe1, oppositePipe2})

	for _, tx := range txs {
		assert.False(t, multiPeer.KnowsTx(tx.Hash()))
	}

	go func(t *testing.T) {
		if err := multiPeer.ReSendTransactions(txs); err != nil {
			t.Error(t)
			return
		}
	}(t)

	receivedMsg, err := oppositePipe1.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}

	var receivedTxs types.Transactions
	if err := receivedMsg.Decode(&receivedTxs); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(txs), len(receivedTxs))

	// It should be received transaction with sorted by times.
	for i, tx := range receivedTxs {
		assert.Equal(t, sortedTxs[i].Hash(), tx.Hash())
		assert.False(t, sortedTxs[i].Time().Equal(tx.Time()))
	}
}
