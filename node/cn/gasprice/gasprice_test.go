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
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type backendMock struct{}

func (mock *backendMock) Progress() klaytn.SyncProgress {
	return klaytn.SyncProgress{}
}

func (mock *backendMock) ProtocolVersion() int {
	return 0
}
func (mock *backendMock) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return nil, nil
}

func (mock *backendMock) ChainDB() database.DBManager {
	return nil
}
func (mock *backendMock) EventMux() *event.TypeMux {
	return nil
}

func (mock *backendMock) AccountManager() *accounts.Manager {
	return nil
}

// BlockChain API
func (mock *backendMock) SetHead(number uint64) {}

func (mock *backendMock) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	return nil, nil
}

func (mock *backendMock) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	return nil, nil
}

func (mock *backendMock) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	return nil, nil, nil
}

func (mock *backendMock) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return nil, nil
}

func (mock *backendMock) GetBlockReceipts(ctx context.Context, blockHash common.Hash) types.Receipts {
	return nil
}

func (mock *backendMock) GetTxLookupInfoAndReceipt(ctx context.Context, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	return nil, common.Hash{}, 0, 0, nil
}

func (mock *backendMock) GetTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return nil, common.Hash{}, 0, 0
}

func (mock *backendMock) GetTd(blockHash common.Hash) *big.Int {
	return nil
}

func (mock *backendMock) GetEVM(ctx context.Context, msg blockchain.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	return nil, nil, nil
}

func (mock *backendMock) SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription {
	return nil
}

func (mock *backendMock) SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription {
	return nil
}

func (mock *backendMock) SubscribeChainSideEvent(ch chan<- blockchain.ChainSideEvent) event.Subscription {
	return nil
}

func (mock *backendMock) IsParallelDBWrite() bool {
	return false
}

func (mock *backendMock) GetNonceInCache(address common.Address) (uint64, bool) {
	return 0, false
}

func (mock *backendMock) IsSenderTxHashIndexingEnabled() bool {
	return false
}

// TxPool API
func (mock *backendMock) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return nil
}

func (mock *backendMock) GetPoolTransactions() (types.Transactions, error) {
	return nil, nil
}

func (mock *backendMock) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	return nil
}

func (mock *backendMock) GetPoolNonce(ctx context.Context, addr common.Address) uint64 {
	return 0
}

func (mock *backendMock) Stats() (pending int, queued int) {
	return 0, 0
}

func (mock *backendMock) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return nil, nil
}

func (mock *backendMock) SubscribeNewTxsEvent(chan<- blockchain.NewTxsEvent) event.Subscription {
	return nil
}

func (mock *backendMock) ChainConfig() *params.ChainConfig {
	return nil
}

func (mock *backendMock) CurrentBlock() *types.Block {
	return nil
}

func (mock *backendMock) GetTxAndLookupInfoInCache(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return nil, common.Hash{}, 0, 0
}

func (mock *backendMock) GetBlockReceiptsInCache(blockHash common.Hash) types.Receipts {
	return nil
}

func (mock *backendMock) GetTxLookupInfoAndReceiptInCache(Hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	return nil, common.Hash{}, 0, 0, nil
}

func TestGasPrice_NewOracle(t *testing.T) {
	mock := &backendMock{}
	params := Config{}
	oracle := NewOracle(mock, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Blocks: 2}
	oracle = NewOracle(mock, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 2, oracle.checkBlocks)
	assert.Equal(t, 1, oracle.maxEmpty)
	assert.Equal(t, 10, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Percentile: -1}
	oracle = NewOracle(mock, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Percentile: 101}
	oracle = NewOracle(mock, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 100, oracle.percentile)

	params = Config{Percentile: 101, Default: big.NewInt(123)}
	oracle = NewOracle(mock, params)

	assert.Equal(t, big.NewInt(123), oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 100, oracle.percentile)
}

func TestGasPrice_SuggestPrice(t *testing.T) {
	mock := &backendMock{}
	params := Config{}
	oracle := NewOracle(mock, params)

	price, err := oracle.SuggestPrice(nil)
	assert.Nil(t, price)
	assert.Nil(t, err)

	params = Config{Default: big.NewInt(123)}
	oracle = NewOracle(mock, params)

	price, err = oracle.SuggestPrice(nil)
	assert.Equal(t, big.NewInt(123), price)
	assert.Nil(t, err)
}
