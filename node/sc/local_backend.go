// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
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
// This file is derived from accounts/abi/bind/backends/simulated.go (2018/06/04).
// Modified and improved for the klaytn development.

package sc

import (
	"context"
	"fmt"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/bloombits"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

func checkCtx(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

type filterLocalBackend struct {
	subbridge *SubBridge
}

func (fb *filterLocalBackend) ChainDB() database.DBManager {
	// TODO-Klaytn consider chain's chainDB instead of bridge's chainDB currently.
	return fb.subbridge.chainDB
}

func (fb *filterLocalBackend) EventMux() *event.TypeMux {
	// TODO-Klaytn consider chain's eventMux instead of bridge's eventMux currently.
	return fb.subbridge.EventMux()
}

func (fb *filterLocalBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	if header := fb.subbridge.blockchain.GetHeaderByHash(hash); header != nil {
		return header, nil
	}
	return nil, fmt.Errorf("the header does not exist (hash: %d)", hash)
}

func (fb *filterLocalBackend) HeaderByNumber(ctx context.Context, block rpc.BlockNumber) (*types.Header, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	// TODO-Klaytn consider pendingblock instead of latest block
	if block == rpc.LatestBlockNumber {
		return fb.subbridge.blockchain.CurrentHeader(), nil
	}
	return fb.subbridge.blockchain.GetHeaderByNumber(uint64(block.Int64())), nil
}

func (fb *filterLocalBackend) GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts {
	if err := checkCtx(ctx); err != nil {
		return nil
	}
	return fb.subbridge.blockchain.GetReceiptsByBlockHash(hash)
}

func (fb *filterLocalBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	return fb.subbridge.blockchain.GetLogsByHash(hash), nil
}

func (fb *filterLocalBackend) SubscribeNewTxsEvent(ch chan<- blockchain.NewTxsEvent) event.Subscription {
	return fb.subbridge.txPool.SubscribeNewTxsEvent(ch)
}

func (fb *filterLocalBackend) SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription {
	return fb.subbridge.blockchain.SubscribeChainEvent(ch)
}

func (fb *filterLocalBackend) SubscribeRemovedLogsEvent(ch chan<- blockchain.RemovedLogsEvent) event.Subscription {
	return fb.subbridge.blockchain.SubscribeRemovedLogsEvent(ch)
}

func (fb *filterLocalBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return fb.subbridge.blockchain.SubscribeLogsEvent(ch)
}

func (fb *filterLocalBackend) BloomStatus() (uint64, uint64) {
	// TODO-Klaytn consider this number of sections.
	// BloomBitsBlocks (const : 4096), the number of processed sections maintained by the chain indexer
	return 4096, 0
}

func (fb *filterLocalBackend) ServiceFilter(_dummyCtx context.Context, session *bloombits.MatcherSession) {
	// TODO-Klaytn this method should implmentation to support indexed tag in solidity
	//for i := 0; i < bloomFilterThreads; i++ {
	//	go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, backend.bloomRequests)
	//}
}

func (fb *filterLocalBackend) ChainConfig() *params.ChainConfig {
	return fb.subbridge.blockchain.Config()
}
