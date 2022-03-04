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

package gasprice

import (
	"context"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"

	"github.com/golang/mock/gomock"
	mock_api "github.com/klaytn/klaytn/api/mocks"
	"github.com/stretchr/testify/assert"
)

const testHead = 32

type testBackend struct {
	chain *blockchain.BlockChain
}

func (b *testBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number > testHead {
		return nil, nil
	}
	if number == rpc.LatestBlockNumber {
		number = testHead
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}

func (b *testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number > testHead {
		return nil, nil
	}
	if number == rpc.LatestBlockNumber {
		number = testHead
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}

func (b *testBackend) GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts {
	return b.chain.GetReceiptsByBlockHash(hash)
}

func (b *testBackend) ChainConfig() *params.ChainConfig {
	return b.chain.Config()
}

func newTestBackend(t *testing.T) *testBackend {
	var (
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr   = crypto.PubkeyToAddress(key.PublicKey)

		gspec = &blockchain.Genesis{
			Config: params.TestChainConfig,
			Alloc:  blockchain.GenesisAlloc{addr: {Balance: big.NewInt(math.MaxInt64)}},
		}
		db      = database.NewMemoryDBManager()
		genesis = gspec.MustCommit(db)
	)
	// Generate testing blocks
	blocks, _ := blockchain.GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, testHead+1, func(i int, b *blockchain.BlockGen) {
		b.SetRewardbase(addr)

		toaddr := common.Address{}
		data := make([]byte, 1)
		gas, _ := types.IntrinsicGas(data, nil, false, params.TestChainConfig.Rules(big.NewInt(0)))
		signer := types.NewEIP155Signer(params.TestChainConfig.ChainID)
		tx, _ := types.SignTx(types.NewTransaction(b.TxNonce(addr), toaddr, big.NewInt(1), gas, nil, data), signer, key)
		b.AddTx(tx)
	})
	// Construct testing chain
	chain, err := blockchain.NewBlockChain(db, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	if err != nil {
		t.Fatalf("Failed to create local chain, %v", err)
	}
	chain.InsertChain(blocks)
	return &testBackend{chain: chain}
}

func TestGasPrice_NewOracle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBackend := mock_api.NewMockBackend(mockCtrl)
	params := Config{}
	oracle := NewOracle(mockBackend, params, nil)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Blocks: 2}
	oracle = NewOracle(mockBackend, params, nil)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 2, oracle.checkBlocks)
	assert.Equal(t, 1, oracle.maxEmpty)
	assert.Equal(t, 10, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Percentile: -1}
	oracle = NewOracle(mockBackend, params, nil)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Percentile: 101}
	oracle = NewOracle(mockBackend, params, nil)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 100, oracle.percentile)

	params = Config{Percentile: 101, Default: big.NewInt(123)}
	oracle = NewOracle(mockBackend, params, nil)

	assert.Equal(t, big.NewInt(123), oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 100, oracle.percentile)
}

func TestGasPrice_SuggestPrice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBackend := mock_api.NewMockBackend(mockCtrl)
	params := Config{}
	testBackend := newTestBackend(t)
	chainConfig := testBackend.ChainConfig()
	chainConfig.UnitPrice = 0
	txPoolWith0 := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, chainConfig, testBackend.chain)
	oracle := NewOracle(mockBackend, params, txPoolWith0)

	price, err := oracle.SuggestPrice(nil)
	assert.Equal(t, price, common.Big0)
	assert.Nil(t, err)

	params = Config{Default: big.NewInt(123)}
	chainConfig.UnitPrice = 25
	txPoolWith25 := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, chainConfig, testBackend.chain)
	oracle = NewOracle(mockBackend, params, txPoolWith25)

	price, err = oracle.SuggestPrice(nil)
	assert.Equal(t, big.NewInt(25), price)
	assert.Nil(t, err)
}
