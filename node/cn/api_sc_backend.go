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
// This file is derived from eth/api_backend.go (2018/06/04).
// Modified and improved for the klaytn development.

package cn

import (
	"context"
	"fmt"
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/bloombits"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn/gasprice"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"math/big"
)

// ServiceChainAPIBackend implements api.Backend for full nodes
type ServiceChainAPIBackend struct {
	sc  *ServiceChain
	gpo *gasprice.Oracle
}

// GetNonceInCache returns (cachedNonce, true) if nonce exists in cache.
// If not, it returns (0, false).
func (b *ServiceChainAPIBackend) GetNonceInCache(addr common.Address) (uint64, bool) {
	return b.sc.blockchain.GetNonceInCache(addr)
}

// GetTxLookupInfoAndReceipt retrieves a tx and lookup info and receipt for a given transaction hash.
func (b *ServiceChainAPIBackend) GetTxLookupInfoAndReceipt(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	return b.sc.blockchain.GetTxLookupInfoAndReceipt(txHash)
}

// GetTxAndLookupInfoInCache retrieves a tx and lookup info for a given transaction hash in cache.
func (b *ServiceChainAPIBackend) GetTxAndLookupInfoInCache(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return b.sc.blockchain.GetTxAndLookupInfoInCache(txHash)
}

// GetBlockReceiptsInCache retrieves receipts for a given block hash in cache.
func (b *ServiceChainAPIBackend) GetBlockReceiptsInCache(blockHash common.Hash) types.Receipts {
	return b.sc.blockchain.GetBlockReceiptsInCache(blockHash)
}

// GetTxLookupInfoAndReceiptInCache retrieves a tx and lookup info and receipt for a given transaction hash in cache.
func (b *ServiceChainAPIBackend) GetTxLookupInfoAndReceiptInCache(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	return b.sc.blockchain.GetTxLookupInfoAndReceiptInCache(txHash)
}

func (b *ServiceChainAPIBackend) ChainConfig() *params.ChainConfig {
	return b.sc.chainConfig
}

func (b *ServiceChainAPIBackend) CurrentBlock() *types.Block {
	return b.sc.blockchain.CurrentBlock()
}

func (b *ServiceChainAPIBackend) SetHead(number uint64) {
	//b.sc.protocolManager.downloader.Cancel()
	b.sc.blockchain.SetHead(number)
}

func (b *ServiceChainAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.sc.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.sc.blockchain.CurrentBlock().Header(), nil
	}
	header := b.sc.blockchain.GetHeaderByNumber(uint64(blockNr))
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNr)
	}
	return header, nil
}

func (b *ServiceChainAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.sc.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.sc.blockchain.CurrentBlock(), nil
	}
	block := b.sc.blockchain.GetBlockByNumber(uint64(blockNr))
	if block == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNr)
	}
	return block, nil
}

func (b *ServiceChainAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.sc.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.sc.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *ServiceChainAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	block := b.sc.blockchain.GetBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("the block does not exist (block hash: %s)", hash.String())
	}
	return block, nil
}

// GetTxAndLookupInfo retrieves a tx and lookup info for a given transaction hash.
func (b *ServiceChainAPIBackend) GetTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return b.sc.blockchain.GetTxAndLookupInfo(hash)
}

// GetBlockReceipts retrieves the receipts for all transactions with given block hash.
func (b *ServiceChainAPIBackend) GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts {
	return b.sc.blockchain.GetReceiptsByBlockHash(hash)
}

func (b *ServiceChainAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	return b.sc.blockchain.GetLogsByHash(hash), nil
}

func (b *ServiceChainAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.sc.blockchain.GetTdByHash(blockHash)
}

func (b *ServiceChainAPIBackend) GetEVM(ctx context.Context, msg blockchain.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.ValidatedSender(), math.MaxBig256)
	vmError := func() error { return nil }

	context := blockchain.NewEVMContext(msg, header, b.sc.BlockChain(), nil)
	return vm.NewEVM(context, state, b.sc.chainConfig, &vmCfg), vmError, nil
}

func (b *ServiceChainAPIBackend) SubscribeRemovedLogsEvent(ch chan<- blockchain.RemovedLogsEvent) event.Subscription {
	return b.sc.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *ServiceChainAPIBackend) SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription {
	return b.sc.BlockChain().SubscribeChainEvent(ch)
}

func (b *ServiceChainAPIBackend) SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription {
	return b.sc.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *ServiceChainAPIBackend) SubscribeChainSideEvent(ch chan<- blockchain.ChainSideEvent) event.Subscription {
	return b.sc.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *ServiceChainAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.sc.BlockChain().SubscribeLogsEvent(ch)
}

func (b *ServiceChainAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.sc.txPool.AddLocal(signedTx)
}

func (b *ServiceChainAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.sc.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *ServiceChainAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.sc.txPool.Get(hash)
}

func (b *ServiceChainAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) uint64 {
	return b.sc.txPool.GetPendingNonce(addr)
}

func (b *ServiceChainAPIBackend) Stats() (pending int, queued int) {
	return b.sc.txPool.Stats()
}

func (b *ServiceChainAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.sc.TxPool().Content()
}

func (b *ServiceChainAPIBackend) SubscribeNewTxsEvent(ch chan<- blockchain.NewTxsEvent) event.Subscription {
	return b.sc.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *ServiceChainAPIBackend) Progress() klaytn.SyncProgress {
	return b.sc.Progress()
}

func (b *ServiceChainAPIBackend) ProtocolVersion() int {
	return b.sc.ProtocolVersion()
}

func (b *ServiceChainAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *ServiceChainAPIBackend) ChainDB() database.DBManager {
	return b.sc.ChainDB()
}

func (b *ServiceChainAPIBackend) EventMux() *event.TypeMux {
	return b.sc.EventMux()
}

func (b *ServiceChainAPIBackend) AccountManager() accounts.AccountManager {
	return b.sc.AccountManager()
}

func (b *ServiceChainAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.sc.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *ServiceChainAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.sc.bloomRequests)
	}
}

func (b *ServiceChainAPIBackend) IsParallelDBWrite() bool {
	return b.sc.BlockChain().IsParallelDBWrite()
}

func (b *ServiceChainAPIBackend) IsSenderTxHashIndexingEnabled() bool {
	return b.sc.BlockChain().IsSenderTxHashIndexingEnabled()
}
