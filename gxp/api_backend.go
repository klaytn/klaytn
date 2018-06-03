package gxp

import (
	"groundx.xyz/go-gxplatform/core"
	"groundx.xyz/go-gxplatform/event"
	"context"
	"math/big"
	"groundx.xyz/go-gxplatform/accounts"
	"groundx.xyz/go-gxplatform/params"
	"groundx.xyz/go-gxplatform/gxp/gasprice"
	"groundx.xyz/go-gxplatform/core/types"
	"groundx.xyz/go-gxplatform/rpc"
	"groundx.xyz/go-gxplatform/core/state"
	"groundx.xyz/go-gxplatform/common"
	"groundx.xyz/go-gxplatform/core/rawdb"
	"groundx.xyz/go-gxplatform/core/vm"
	"groundx.xyz/go-gxplatform/common/math"
	"groundx.xyz/go-gxplatform/gxdb"
	"groundx.xyz/go-gxplatform/core/bloombits"
	"groundx.xyz/go-gxplatform/gxp/downloader"
)

// GxpAPIBackend implements gxpapi.Backend for full nodes
type GxpAPIBackend struct {
	gxp *GXP
	gpo *gasprice.Oracle
}

func (b *GxpAPIBackend) ChainConfig() *params.ChainConfig {
	return b.gxp.chainConfig
}

func (b *GxpAPIBackend) CurrentBlock() *types.Block {
	return b.gxp.blockchain.CurrentBlock()
}

func (b *GxpAPIBackend) SetHead(number uint64) {
	//b.gxp.protocolManager.downloader.Cancel()
	b.gxp.blockchain.SetHead(number)
}

func (b *GxpAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.gxp.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.gxp.blockchain.CurrentBlock().Header(), nil
	}
	return b.gxp.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *GxpAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.gxp.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.gxp.blockchain.CurrentBlock(), nil
	}
	return b.gxp.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *GxpAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.gxp.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.gxp.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *GxpAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.gxp.blockchain.GetBlockByHash(hash), nil
}

func (b *GxpAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if number := rawdb.ReadHeaderNumber(b.gxp.chainDb, hash); number != nil {
		return rawdb.ReadReceipts(b.gxp.chainDb, hash, *number), nil
	}
	return nil, nil
}

func (b *GxpAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := rawdb.ReadHeaderNumber(b.gxp.chainDb, hash)
	if number == nil {
		return nil, nil
	}
	receipts := rawdb.ReadReceipts(b.gxp.chainDb, hash, *number)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *GxpAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.gxp.blockchain.GetTdByHash(blockHash)
}

func (b *GxpAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.gxp.BlockChain(), nil)
	return vm.NewEVM(context, state, b.gxp.chainConfig, vmCfg), vmError, nil
}

func (b *GxpAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.gxp.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *GxpAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.gxp.BlockChain().SubscribeChainEvent(ch)
}

func (b *GxpAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.gxp.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *GxpAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.gxp.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *GxpAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.gxp.BlockChain().SubscribeLogsEvent(ch)
}

func (b *GxpAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.gxp.txPool.AddLocal(signedTx)
}

func (b *GxpAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.gxp.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *GxpAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.gxp.txPool.Get(hash)
}

func (b *GxpAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.gxp.txPool.State().GetNonce(addr), nil
}

func (b *GxpAPIBackend) Stats() (pending int, queued int) {
	return b.gxp.txPool.Stats()
}

func (b *GxpAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.gxp.TxPool().Content()
}

func (b *GxpAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.gxp.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *GxpAPIBackend) Downloader() *downloader.Downloader {
	return b.gxp.Downloader()
}

func (b *GxpAPIBackend) ProtocolVersion() int {
	return b.gxp.EthVersion()
}

func (b *GxpAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *GxpAPIBackend) ChainDb() gxdb.Database {
	return b.gxp.ChainDb()
}

func (b *GxpAPIBackend) EventMux() *event.TypeMux {
	return b.gxp.EventMux()
}

func (b *GxpAPIBackend) AccountManager() *accounts.Manager {
	return b.gxp.AccountManager()
}

func (b *GxpAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.gxp.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *GxpAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.gxp.bloomRequests)
	}
}
