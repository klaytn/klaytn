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

package blockchain

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	mock_statedb "github.com/klaytn/klaytn/storage/statedb/mocks"
)

func getTestBlock(t *testing.T) types.Block {
	var block types.Block
	blockDump, err := ioutil.ReadFile("../tests/b1.rlp")
	if err != nil {
		t.Fatal(err)
	}
	if err := rlp.DecodeBytes(common.FromHex(string(blockDump)), &block); err != nil {
		t.Fatal("decode error: ", err)
	}
	return block
}

func getTestLogs(t *testing.T) ([]*types.Receipt, []*types.LogForStorage, []byte) {
	testReceipts := []*types.Receipt{
		{
			Logs: []*types.Log{
				{
					Address:     common.HexToAddress("0xecf8f87f810ecf450940c9f60066b4a7a501d6a7"),
					BlockHash:   common.HexToHash("0x1430cbc43787f54eb0b77dfd4de8246a8c256aeac79c1e0c94dc7c59ba4e7c57"),
					BlockNumber: 1,
					Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000001a055690d9db80000"),
					Index:       0,
					TxIndex:     0,
					TxHash:      common.HexToHash("0x7bea0019816d8146eec8cebd7dea862aa6e97fb5387c354f7d484e1d1e79efec"),
					Topics: []common.Hash{
						common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
						common.HexToHash("0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"),
					},
				},
			},
		},
		{
			Logs: []*types.Log{
				{
					Address:     common.HexToAddress("0x91d6f7d2537d8a0bd7d487dcc59151ebc00da306"),
					BlockHash:   common.HexToHash("0x1430cbc43787f54eb0b77dfd4de8246a8c256aeac79c1e0c94dc7c59ba4e7c57"),
					BlockNumber: 1,
					Data:        []byte{},
					Index:       1,
					TxIndex:     3,
					TxHash:      common.HexToHash("0xb8e08da739c350f7e11e392cbac5b48d992e91990b643726be6c8ccf12a07e2b"),
					Topics: []common.Hash{
						common.HexToHash("0x1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace113"),
						common.HexToHash("0x84d9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace243"),
					},
				},
			},
		},
	}

	testStorageLogs := []*types.LogForStorage{}
	for _, receipts := range testReceipts {
		for _, log := range receipts.Logs {
			testStorageLogs = append(testStorageLogs, (*types.LogForStorage)(log))
		}
	}

	encodedBlockLogs, err := rlp.EncodeToBytes(testStorageLogs)
	if err != nil {
		t.Fatal(err)
	}
	return testReceipts, testStorageLogs, encodedBlockLogs
}

func TestBlockChain_sendKESSubscriptionData(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTrieNodeCache := mock_statedb.NewMockTrieNodeCache(mockCtrl)

	chainCh := make(chan ChainEvent)
	logsCh := make(chan []*types.Log)
	defer close(chainCh)
	defer close(logsCh)

	// prepare Blockchain to be tested
	bc := BlockChain{stateCache: state.NewDatabaseWithExistingCache(database.NewMemoryDBManager(), mockTrieNodeCache)}
	bc.SubscribeChainEvent(chainCh)
	bc.SubscribeLogsEvent(logsCh)

	// prepare test data to be send by sendKESSubscriptionData
	block := getTestBlock(t)
	blockLogsKey := append(kesCachePrefixBlockLogs, block.Number().Bytes()...)
	receipts, blockLogs, encodedBlockLogs := getTestLogs(t)

	// resultCh receive data when checkSubscribedData finished successfully
	resultCh := make(chan struct{})
	defer close(resultCh)

	// check whether expected data are delivered to the subscription channel
	checkSubscribedData := func() {
		subscribedChain := <-chainCh
		if block.Hash() != subscribedChain.Block.Hash() ||
			!deep.Equal(block.Header(), subscribedChain.Block.Header()) ||
			!deep.Equal(block.Transactions(), subscribedChain.Block.Transactions()) {
			t.Errorf("\nexpected block=%v\nactual block=%v", block.String(), subscribedChain.Block.String())
			// to terminate blocking main thread
			panic("block doesn't match")
		}

		subscribedLogs := <-logsCh
		if len(subscribedLogs) != len(blockLogs) {
			t.Errorf("expected length=%d, actual length=%d", len(blockLogs), len(subscribedLogs))
			// to terminate blocking main thread
			panic("log length doesn't match")
		}

		for i, testLog := range blockLogs {
			expectedLog := (*types.Log)(testLog)
			if !deep.Equal(expectedLog, subscribedLogs[i]) {
				t.Errorf("%dth log doen't match\nexpected log=%v\nactual log=%v",
					i, expectedLog.String(), subscribedLogs[i].String())
				// to terminate blocking main thread
				panic("log doesn't match")
			}
		}
		resultCh <- struct{}{}
	}

	// in normal case, block logs are retrieved by remote cache
	{
		mockTrieNodeCache.EXPECT().Get(blockLogsKey).Return(encodedBlockLogs).Times(1)
		bc.db = nil

		go checkSubscribedData()
		bc.sendKESSubscriptionData(&block)

		select {
		case <-resultCh:
		case <-time.NewTimer(100 * time.Millisecond).C:
			t.Fatal("timeout")
		}
	}

	// if remote cache has invalid block logs, block logs are retrieved by database
	{
		mockTrieNodeCache.EXPECT().Get(blockLogsKey).Return([]byte{0x1, 0x2}).Times(1)
		bc.db = database.NewMemoryDBManager()
		bc.db.WriteReceipts(block.Hash(), block.NumberU64(), receipts)

		go checkSubscribedData()
		bc.sendKESSubscriptionData(&block)

		select {
		case <-resultCh:
		case <-time.NewTimer(100 * time.Millisecond).C:
			t.Fatal("timeout")
		}
	}

	// if remote cache doesn't have block logs, block logs are retrieved by database
	{
		mockTrieNodeCache.EXPECT().Get(blockLogsKey).Return(nil).Times(1)
		bc.db = database.NewMemoryDBManager()
		bc.db.WriteReceipts(block.Hash(), block.NumberU64(), receipts)

		go checkSubscribedData()
		bc.sendKESSubscriptionData(&block)

		select {
		case <-resultCh:
		case <-time.NewTimer(100 * time.Millisecond).C:
			t.Fatal("timeout")
		}
	}
}
