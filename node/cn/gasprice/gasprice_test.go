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

	"github.com/golang/mock/gomock"
	mock_api "github.com/klaytn/klaytn/api/mocks"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/event/mocks"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
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

func (b *testBackend) CurrentBlock() *types.Block {
	return b.chain.CurrentBlock()
}

func (b *testBackend) SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription {
	return b.chain.SubscribeChainHeadEvent(ch)
}

func (b *testBackend) teardown() {
	b.chain.Stop()
}

// newTestBackend creates a test backend. OBS: don't forget to invoke tearDown
// after use, otherwise the blockchain instance will mem-leak via goroutines.
func newTestBackend(t *testing.T, magmaBlock, dragonBlock *big.Int) *testBackend {
	var (
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr   = crypto.PubkeyToAddress(key.PublicKey)
		config = params.TestChainConfig.Copy() // needs copy because it is modified below
		gspec  = &blockchain.Genesis{
			Config: config,
			Alloc:  blockchain.GenesisAlloc{addr: {Balance: big.NewInt(math.MaxInt64)}},
		}
		signer = types.LatestSignerForChainID(gspec.Config.ChainID)
		db      = database.NewMemoryDBManager()
		genesis = gspec.MustCommit(db)
	)

	config.EthTxTypeCompatibleBlock = magmaBlock
	config.IstanbulCompatibleBlock = magmaBlock
	config.MagmaCompatibleBlock = magmaBlock
	config.KoreCompatibleBlock = dragonBlock
	config.ShanghaiCompatibleBlock = dragonBlock
	config.CancunCompatibleBlock = dragonBlock
	config.DragonCompatibleBlock = dragonBlock
	if magmaBlock != nil {
		config.Governance = params.GetDefaultGovernanceConfig()
		config.Istanbul = params.GetDefaultIstanbulConfig()
		config.Governance.KIP71.LowerBoundBaseFee = 0
	}
	blocks, _ := blockchain.GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, testHead+1, func(i int, b *blockchain.BlockGen) {
		toaddr := common.Address{}
		// To test fee history, rewardbase should be different from the sender address
		b.SetRewardbase(toaddr)

		var txdata types.TxInternalData
		if config.MagmaCompatibleBlock != nil && b.Number().Cmp(config.MagmaCompatibleBlock) >= 0 {
			txdata = &types.TxInternalDataEthereumDynamicFee{
				ChainID:       gspec.Config.ChainID,
				AccountNonce:  b.TxNonce(addr),
				Recipient:     &common.Address{},
				GasLimit:      30000,
				GasFeeCap:     big.NewInt(100 * params.Ston),
				GasTipCap:     big.NewInt(int64(i+1) * params.Ston),
				Payload:       []byte{},
				Amount:        big.NewInt(100),
			}
		} else {
			txdata = &types.TxInternalDataLegacy{
				AccountNonce:  b.TxNonce(addr),
				Recipient:     &common.Address{},
				GasLimit:      21000,
				Price:         big.NewInt(int64(i+1) * params.Ston),
				Amount:        big.NewInt(100),
				Payload:       []byte{},
			}
		}
		tx, _ := types.SignTx(types.NewTx(txdata), signer, key)

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

	sub := mocks.NewMockSubscription(mockCtrl)
	mockBackend.EXPECT().SubscribeChainHeadEvent(gomock.Any()).Return(sub).Times(5)

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
	testBackend := newTestBackend(t, nil, nil)
	defer testBackend.teardown()
	chainConfig := testBackend.ChainConfig()
	chainConfig.UnitPrice = 0
	txPoolWith0 := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, chainConfig, testBackend.chain)

	sub := mocks.NewMockSubscription(mockCtrl)
	mockBackend.EXPECT().SubscribeChainHeadEvent(gomock.Any()).Return(sub).Times(2)
	oracle := NewOracle(mockBackend, params, txPoolWith0)

	currentBlock := testBackend.CurrentBlock()
	mockBackend.EXPECT().ChainConfig().Return(chainConfig).Times(2)
	mockBackend.EXPECT().CurrentBlock().Return(currentBlock).Times(2)

	price, err := oracle.SuggestPrice(nil)
	assert.Equal(t, price, common.Big0)
	assert.Nil(t, err)

	params = Config{Default: big.NewInt(123)}
	chainConfig.UnitPrice = 25
	mockBackend.EXPECT().ChainConfig().Return(chainConfig).Times(2)
	txPoolWith25 := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, chainConfig, testBackend.chain)
	oracle = NewOracle(mockBackend, params, txPoolWith25)

	price, err = oracle.SuggestPrice(nil)
	assert.Equal(t, big.NewInt(25), price)
	assert.Nil(t, err)
}

func TestSuggestTipCap(t *testing.T) {
	config := Config{
		Blocks:     3,
		Percentile: 60,
		Default:    big.NewInt(params.Ston),
		MaxHeaderHistory: 30,
		MaxBlockHistory: 30,
	}

	var cases = []struct {
		magmaBlock *big.Int // Magma fork block number
		dragonBlock *big.Int // Dragon fork block number
		expect *big.Int // Expected gasprice suggestion
	}{
		{nil, nil, big.NewInt(1)}, // If not Magma forked, should return unitPrice (which is 1 for test)

		{big.NewInt(0), nil, common.Big0}, // After Magma fork and before Dragon fork, should return 0

		// After Dragon fork
		{big.NewInt(0), big.NewInt(0), big.NewInt(params.Ston * int64(30))},   // Fork point in genesis
		{big.NewInt(1), big.NewInt(1), big.NewInt(params.Ston * int64(30))},   // Fork point in first block
		{big.NewInt(32), big.NewInt(32), big.NewInt(params.Ston * int64(30))}, // Fork point in last block
		{big.NewInt(33), big.NewInt(33), big.NewInt(params.Ston * int64(30))}, // Fork point in the future
	}
	for _, c := range cases {
		testBackend := newTestBackend(t, c.magmaBlock, c.dragonBlock)
		chainConfig := testBackend.ChainConfig()
		txPool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, chainConfig, testBackend.chain)
		oracle := NewOracle(testBackend, config, txPool)

		// The gas price sampled is: 32G, 31G, 30G, 29G, 28G, 27G
		got, err := oracle.SuggestTipCap(context.Background())
		testBackend.teardown()
		if err != nil {
			t.Fatalf("Failed to retrieve recommended gas price: %v", err)
		}
		if got.Cmp(c.expect) != 0 {
			t.Fatalf("Gas price mismatch, want %d, got %d",c.expect, got)
		}
	}
}
