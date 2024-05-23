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
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
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

var (
	addr1 = common.HexToAddress("111")
	addr2 = common.HexToAddress("222")
	addrs []common.Address
)

var (
	topic1 common.Hash
	topic2 common.Hash
	topics [][]common.Hash
)

var (
	begin = int64(12345)
	end   = int64(12345)
)

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
		db = database.NewMemoryDBManager()

		mux        = new(event.TypeMux)
		txFeed     = new(event.Feed)
		rmLogsFeed = new(event.Feed)
		logsFeed   = new(event.Feed)
		chainFeed  = new(event.Feed)
		backend    = &testBackend{mux, db, 0, txFeed, rmLogsFeed, logsFeed, chainFeed, params.TestChainConfig}
		key1, _    = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr       = crypto.PubkeyToAddress(key1.PublicKey)
		signer     = types.NewLondonSigner(big.NewInt(1))
		// Logging contract
		contract  = common.Address{0xfe}
		contract2 = common.Address{0xff}
		abiStr    = `[{"inputs":[],"name":"log0","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"t1","type":"uint256"}],"name":"log1","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"t1","type":"uint256"},{"internalType":"uint256","name":"t2","type":"uint256"}],"name":"log2","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"t1","type":"uint256"},{"internalType":"uint256","name":"t2","type":"uint256"},{"internalType":"uint256","name":"t3","type":"uint256"}],"name":"log3","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"t1","type":"uint256"},{"internalType":"uint256","name":"t2","type":"uint256"},{"internalType":"uint256","name":"t3","type":"uint256"},{"internalType":"uint256","name":"t4","type":"uint256"}],"name":"log4","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
		bytecode  = common.FromHex("608060405234801561001057600080fd5b50600436106100575760003560e01c80630aa731851461005c5780632a4c08961461006657806378b9a1f314610082578063c670f8641461009e578063c683d6a3146100ba575b600080fd5b6100646100d6565b005b610080600480360381019061007b9190610143565b6100dc565b005b61009c60048036038101906100979190610196565b6100e8565b005b6100b860048036038101906100b391906101d6565b6100f2565b005b6100d460048036038101906100cf9190610203565b6100fa565b005b600080a0565b808284600080a3505050565b8082600080a25050565b80600080a150565b80828486600080a450505050565b600080fd5b6000819050919050565b6101208161010d565b811461012b57600080fd5b50565b60008135905061013d81610117565b92915050565b60008060006060848603121561015c5761015b610108565b5b600061016a8682870161012e565b935050602061017b8682870161012e565b925050604061018c8682870161012e565b9150509250925092565b600080604083850312156101ad576101ac610108565b5b60006101bb8582860161012e565b92505060206101cc8582860161012e565b9150509250929050565b6000602082840312156101ec576101eb610108565b5b60006101fa8482850161012e565b91505092915050565b6000806000806080858703121561021d5761021c610108565b5b600061022b8782880161012e565b945050602061023c8782880161012e565b935050604061024d8782880161012e565b925050606061025e8782880161012e565b9150509295919450925056fea264697066735822122073a4b156f487e59970dc1ef449cc0d51467268f676033a17188edafcee861f9864736f6c63430008110033")
		genesis   = blockchain.Genesis{Config: params.TestChainConfig, Alloc: blockchain.GenesisAlloc{
			addr:      {Balance: big.NewInt(0).Mul(big.NewInt(100), big.NewInt(params.KAIA))},
			contract:  {Balance: big.NewInt(0), Code: bytecode},
			contract2: {Balance: big.NewInt(0), Code: bytecode},
		}}
		hash1 = common.BytesToHash([]byte("topic1"))
		hash2 = common.BytesToHash([]byte("topic2"))
		hash3 = common.BytesToHash([]byte("topic3"))
		hash4 = common.BytesToHash([]byte("topic4"))
	)
	defer db.Close()
	contractABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		t.Fatal(err)
	}
	chain, _ := blockchain.GenerateChain(params.TestChainConfig, genesis.MustCommit(db), gxhash.NewFaker(), db, 1000, func(i int, gen *blockchain.BlockGen) {
		switch i {
		case 1:
			data, err := contractABI.Pack("log1", hash1.Big())
			if err != nil {
				t.Fatal(err)
			}
			tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataLegacy{
				AccountNonce: 0,
				Price:        big.NewInt(30000),
				GasLimit:     30000,
				Recipient:    &contract,
				Payload:      data,
				Amount:       big.NewInt(0),
			}), signer, key1)
			gen.AddTx(tx)
			tx2, _ := types.SignTx(types.NewTx(&types.TxInternalDataLegacy{
				AccountNonce: 1,
				Price:        big.NewInt(30000),
				GasLimit:     30000,
				Recipient:    &contract2,
				Payload:      data,
				Amount:       big.NewInt(0),
			}), signer, key1)
			gen.AddTx(tx2)
		case 2:
			data, err := contractABI.Pack("log2", hash2.Big(), hash1.Big())
			if err != nil {
				t.Fatal(err)
			}
			tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataLegacy{
				AccountNonce: 2,
				Price:        big.NewInt(30000),
				GasLimit:     30000,
				Recipient:    &contract2,
				Payload:      data,
				Amount:       big.NewInt(0),
			}), signer, key1)
			gen.AddTx(tx)

		case 998:
			data, err := contractABI.Pack("log1", hash3.Big())
			if err != nil {
				t.Fatal(err)
			}
			tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataLegacy{
				AccountNonce: 3,
				Price:        big.NewInt(30000),
				GasLimit:     30000,
				Recipient:    &contract2,
				Payload:      data,
				Amount:       big.NewInt(0),
			}), signer, key1)
			gen.AddTx(tx)
		case 999:
			data, err := contractABI.Pack("log1", hash4.Big())
			if err != nil {
				t.Fatal(err)
			}
			tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataLegacy{
				AccountNonce: 4,
				Price:        big.NewInt(30000),
				GasLimit:     30000,
				Recipient:    &contract2,
				Payload:      data,
				Amount:       big.NewInt(0),
			}), signer, key1)
			gen.AddTx(tx)
		}
	})
	bc, err := blockchain.NewBlockChain(db, nil, params.TestChainConfig, gxhash.NewFaker(), vm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = bc.InsertChain(chain)
	if err != nil {
		t.Fatal(err)
	}

	for i, block := range chain {
		readBlock := db.ReadBlock(block.Hash(), block.NumberU64())
		if readBlock == nil {
			fmt.Println(i, ": block num is", block.NumberU64(), ", and it doesn't have block inside")
		}
	}
	filter := NewRangeFilter(backend, 0, -1, []common.Address{contract2}, [][]common.Hash{{hash1, hash2, hash3, hash4}})

	logs, _ := filter.Logs(context.Background())
	if len(logs) != 4 {
		t.Error("expected 4 log, got", len(logs))
	}

	filter = NewRangeFilter(backend, 900, 999, []common.Address{contract2}, [][]common.Hash{{hash3}})
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = NewRangeFilter(backend, 990, -1, []common.Address{contract2}, [][]common.Hash{{hash3}})
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = NewRangeFilter(backend, 1, 10, nil, [][]common.Hash{{hash1, hash2}})

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 3 {
		t.Error("expected 3 log, got", len(logs))
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
