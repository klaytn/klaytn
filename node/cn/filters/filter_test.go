// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/filters/filter_system_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package filters

import (
	"context"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	cn "github.com/klaytn/klaytn/node/cn/filters/mock"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var addr1 = common.HexToAddress("111")
var addr2 = common.HexToAddress("222")
var addrs []common.Address

var topic1 common.Hash
var topic2 common.Hash
var topics [][]common.Hash

var begin = int64(12345)
var end = int64(12345)

var header *types.Header

var someErr = errors.New("some error")

func init() {
	addrs = []common.Address{addr1, addr2}
	topics = [][]common.Hash{{topic1}, {topic2}}
	header = &types.Header{
		Number:     big.NewInt(int64(123)),
		BlockScore: big.NewInt(int64(1)),
		Extra:      addrs[0][:],
		Governance: addrs[0][:],
		Vote:       addrs[0][:],
	}
}

func genFilter(t *testing.T) (*gomock.Controller, *cn.MockBackend, *Filter) {
	mockCtrl := gomock.NewController(t)
	mockBackend := cn.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().BloomStatus().Return(uint64(123), uint64(321)).Times(1)
	newFilter := NewRangeFilter(mockBackend, begin, end, addrs, topics)
	return mockCtrl, mockBackend, newFilter
}

func TestFilter_New(t *testing.T) {
	mockCtrl, mockBackend, newFilter := genFilter(t)
	defer mockCtrl.Finish()

	assert.NotNil(t, newFilter)
	assert.Equal(t, mockBackend, newFilter.backend)
	assert.Equal(t, begin, newFilter.begin)
	assert.Equal(t, end, newFilter.end)
	assert.Equal(t, topics, newFilter.topics)
	assert.Equal(t, addrs, newFilter.addresses)
	assert.NotNil(t, newFilter.matcher)
}

func TestFilter_Logs(t *testing.T) {
	ctx := context.Background()
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().HeaderByNumber(ctx, rpc.LatestBlockNumber).Times(1).Return(nil, nil)
		logs, err := newFilter.Logs(ctx)
		assert.Nil(t, logs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestFilter_unindexedLogs(t *testing.T) {
	ctx := context.Background()
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().HeaderByNumber(ctx, rpc.BlockNumber(newFilter.begin)).Times(1).Return(nil, nil)
		logs, err := newFilter.unindexedLogs(ctx, uint64(newFilter.end))
		assert.Nil(t, logs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().HeaderByNumber(ctx, rpc.BlockNumber(newFilter.begin)).Times(1).Return(header, nil)
		logs, err := newFilter.unindexedLogs(ctx, uint64(newFilter.end))
		assert.Nil(t, logs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestFilter_checkMatches(t *testing.T) {
	ctx := context.Background()
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().GetLogs(ctx, header.Hash()).Times(1).Return(nil, someErr)
		logs, err := newFilter.checkMatches(ctx, header)
		assert.Nil(t, logs)
		assert.Equal(t, someErr, err)
		mockCtrl.Finish()
	}
	{
		mockCtrl, mockBackend, newFilter := genFilter(t)
		mockBackend.EXPECT().GetLogs(ctx, header.Hash()).Times(1).Return(nil, nil)
		logs, err := newFilter.checkMatches(ctx, header)
		assert.Nil(t, logs)
		assert.NoError(t, err)
		mockCtrl.Finish()
	}
}

func TestFilter_bloomFilter(t *testing.T) {
	{
		assert.True(t, bloomFilter(types.Bloom{}, nil, nil))
	}
	{
		assert.False(t, bloomFilter(types.Bloom{}, nil, [][]common.Hash{{topic1}}))
	}
	{
		assert.False(t, bloomFilter(types.Bloom{}, []common.Address{addr1}, nil))
	}
}

func makeReceipt(addr common.Address) *types.Receipt {
	receipt := genReceipt(false, 0)
	receipt.Logs = []*types.Log{
		{Address: addr},
	}
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt
}

func BenchmarkFilters(b *testing.B) {
	var (
		db         = database.NewMemoryDBManager()
		mux        = new(event.TypeMux)
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig}
		key1, _    = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1      = crypto.PubkeyToAddress(key1.PublicKey)
		addr2      = common.BytesToAddress([]byte("jeff"))
		addr3      = common.BytesToAddress([]byte("ethereum"))
		addr4      = common.BytesToAddress([]byte("random addresses please"))
	)
	defer db.Close()

	genesis := blockchain.GenesisBlockForTesting(db, addr1, big.NewInt(1000000))
	chain, receipts := blockchain.GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), db, 100010, func(i int, gen *blockchain.BlockGen) {
		switch i {
		case 2403:
			receipt := makeReceipt(addr1)
			gen.AddUncheckedReceipt(receipt)
		case 1034:
			receipt := makeReceipt(addr2)
			gen.AddUncheckedReceipt(receipt)
		case 34:
			receipt := makeReceipt(addr3)
			gen.AddUncheckedReceipt(receipt)
		case 99999:
			receipt := makeReceipt(addr4)
			gen.AddUncheckedReceipt(receipt)

		}
	})
	for i, block := range chain {
		db.WriteBlock(block)
		db.WriteCanonicalHash(block.Hash(), block.NumberU64())
		db.WriteHeadBlockHash(block.Hash())
		db.WriteReceipts(block.Hash(), block.NumberU64(), receipts[i])
	}
	b.ResetTimer()

	filter := NewRangeFilter(backend, 0, -1, []common.Address{addr1, addr2, addr3, addr4}, nil)

	for i := 0; i < b.N; i++ {
		logs, _ := filter.Logs(context.Background())
		if len(logs) != 4 {
			b.Fatal("expected 4 logs, got", len(logs))
		}
	}
}

func genReceipt(failed bool, cumulativeGasUsed uint64) *types.Receipt {
	r := &types.Receipt{GasUsed: cumulativeGasUsed}
	if failed {
		r.Status = types.ReceiptStatusFailed
	} else {
		r.Status = types.ReceiptStatusSuccessful
	}
	return r
}

func TestFilters(t *testing.T) {
	var (
		db         = database.NewMemoryDBManager()
		mux        = new(event.TypeMux)
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig}
		key1, _    = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr       = crypto.PubkeyToAddress(key1.PublicKey)

		hash1 = common.BytesToHash([]byte("topic1"))
		hash2 = common.BytesToHash([]byte("topic2"))
		hash3 = common.BytesToHash([]byte("topic3"))
		hash4 = common.BytesToHash([]byte("topic4"))
	)
	defer db.Close()

	genesis := blockchain.GenesisBlockForTesting(db, addr, big.NewInt(1000000))
	chain, receipts := blockchain.GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), db, 1000, func(i int, gen *blockchain.BlockGen) {
		switch i {
		case 1:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash1},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), 1, big.NewInt(1), nil))
		case 2:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash2},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(2, common.HexToAddress("0x2"), big.NewInt(2), 2, big.NewInt(2), nil))

		case 998:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash3},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(998, common.HexToAddress("0x998"), big.NewInt(998), 998, big.NewInt(998), nil))
		case 999:
			receipt := genReceipt(false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash4},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(999, common.HexToAddress("0x999"), big.NewInt(999), 999, big.NewInt(999), nil))
		}
	})
	for i, block := range chain {
		db.WriteBlock(block)
		db.WriteCanonicalHash(block.Hash(), block.NumberU64())
		db.WriteHeadBlockHash(block.Hash())
		db.WriteReceipts(block.Hash(), block.NumberU64(), receipts[i])
	}

	filter := NewRangeFilter(backend, 0, -1, []common.Address{addr}, [][]common.Hash{{hash1, hash2, hash3, hash4}})

	logs, _ := filter.Logs(context.Background())
	if len(logs) != 4 {
		t.Error("expected 4 log, got", len(logs))
	}

	filter = NewRangeFilter(backend, 900, 999, []common.Address{addr}, [][]common.Hash{{hash3}})
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = NewRangeFilter(backend, 990, -1, []common.Address{addr}, [][]common.Hash{{hash3}})
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = NewRangeFilter(backend, 1, 10, nil, [][]common.Hash{{hash1, hash2}})

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 2 {
		t.Error("expected 2 log, got", len(logs))
	}

	failHash := common.BytesToHash([]byte("fail"))
	filter = NewRangeFilter(backend, 0, -1, nil, [][]common.Hash{{failHash}})

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}

	failAddr := common.BytesToAddress([]byte("failmenow"))
	filter = NewRangeFilter(backend, 0, -1, []common.Address{failAddr}, nil)

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}

	filter = NewRangeFilter(backend, 0, -1, nil, [][]common.Hash{{failHash}, {hash1}})

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}
}
