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

package chaindatafetcher

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

type ChainDataFetcher struct {
	config *ChainDataFetcherConfig

	blockchain    *blockchain.BlockChain
	blockchainAPI *api.PublicBlockChainAPI
	debugAPI      *cn.PrivateDebugAPI

	chainCh  chan blockchain.ChainEvent
	chainSub event.Subscription

	reqCh  chan *request // TODO-ChainDataFetcher add logic to insert new requests from APIs to this channel
	stopCh chan struct{}

	numHandlers int

	checkpointMu  sync.RWMutex
	checkpoint    int64
	checkpointMap map[int64]struct{}

	wg sync.WaitGroup

	repo Repository

	fetchingStarted      bool
	fetchingStopCh       chan struct{}
	fetchingWg           sync.WaitGroup
	rangeFetchingStarted bool
	rangeFetchingStopCh  chan struct{}
	rangeFetchingWg      sync.WaitGroup
}

func NewChainDataFetcher(ctx *node.ServiceContext, cfg *ChainDataFetcherConfig) (*ChainDataFetcher, error) {
	repo, err := kas.NewRepository(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		logger.Error("Failed to create new Repository", "err", err, "user", cfg.DBUser, "host", cfg.DBHost, "port", cfg.DBPort, "name", cfg.DBName)
		return nil, err
	}
	checkpoint, err := repo.ReadCheckpoint()
	if err != nil {
		logger.Error("Failed to get checkpoint", "err", err)
		return nil, err
	}
	return &ChainDataFetcher{
		config:        cfg,
		chainCh:       make(chan blockchain.ChainEvent, cfg.BlockChannelSize),
		reqCh:         make(chan *request, cfg.JobChannelSize),
		stopCh:        make(chan struct{}),
		numHandlers:   cfg.NumHandlers,
		checkpoint:    checkpoint,
		checkpointMap: make(map[int64]struct{}),
		repo:          repo,
	}, nil
}

func (f *ChainDataFetcher) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (f *ChainDataFetcher) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "chaindatafetcher",
			Version:   "1.0",
			Service:   NewPublicChainDataFetcherAPI(f),
			Public:    true,
		},
	}
}

func (f *ChainDataFetcher) Start(server p2p.Server) error {
	// launch multiple goroutines to handle new blocks
	for i := 0; i < f.numHandlers; i++ {
		go f.handleRequest()
	}

	if !f.config.NoDefaultStart {
		if err := f.startFetching(); err != nil {
			return err
		}
	}
	return nil
}

func (f *ChainDataFetcher) Stop() error {
	f.stopFetching()
	f.stopRangeFetching()
	close(f.stopCh)
	logger.Info("wait for all goroutines to be terminated...")
	f.wg.Wait()
	logger.Info("terminated all goroutines for chaindatafetcher")
	return nil
}

func (f *ChainDataFetcher) sendRequests(start, end uint64, reqType requestType, shouldUpdateCheckpoint bool, stopCh chan struct{}) {
	for i := start; i <= end; i++ {
		select {
		case <-stopCh:
			return
		case f.reqCh <- newRequest(reqType, shouldUpdateCheckpoint, i):
		}
	}
}

func (f *ChainDataFetcher) startFetching() error {
	if f.fetchingStarted {
		return errors.New("fetching is already started")
	}

	f.chainSub = f.blockchain.SubscribeChainEvent(f.chainCh)
	currentBlock := f.blockchain.CurrentHeader().Number.Uint64()

	f.fetchingStopCh = make(chan struct{})
	f.fetchingWg.Add(1)
	go func() {
		defer f.fetchingWg.Done()
		f.sendRequests(uint64(f.checkpoint), currentBlock, requestTypeAll, true, f.fetchingStopCh)
	}()
	f.fetchingStarted = true
	return nil
}

func (f *ChainDataFetcher) stopFetching() error {
	if !f.fetchingStarted {
		return errors.New("fetching is not running")
	}

	f.chainSub.Unsubscribe()
	close(f.fetchingStopCh)
	f.fetchingWg.Wait()
	f.fetchingStarted = false
	return nil
}

func (f *ChainDataFetcher) startRangeFetching(start, end uint64, reqType requestType) error {
	if f.rangeFetchingStarted {
		return errors.New("range fetching is already started")
	}
	f.rangeFetchingStopCh = make(chan struct{})
	f.rangeFetchingWg.Add(1)
	go func() {
		defer f.rangeFetchingWg.Done()
		f.sendRequests(start, end, reqType, false, f.rangeFetchingStopCh)
	}()
	f.rangeFetchingStarted = true
	return nil
}

func (f *ChainDataFetcher) stopRangeFetching() error {
	if !f.rangeFetchingStarted {
		return errors.New("range fetching is not running")
	}
	close(f.rangeFetchingStopCh)
	f.rangeFetchingWg.Wait()
	f.rangeFetchingStarted = false
	return nil
}

func (f *ChainDataFetcher) makeChainEvent(blockNumber uint64) (blockchain.ChainEvent, error) {
	var logs []*types.Log
	block := f.blockchain.GetBlockByNumber(blockNumber)
	receipts := f.blockchain.GetReceiptsByBlockHash(block.Hash())
	for _, r := range receipts {
		logs = append(logs, r.Logs...)
	}
	var internalTraces []*vm.InternalTxTrace
	if block.Transactions().Len() > 0 {
		fct := "fastCallTracer"
		results, err := f.debugAPI.TraceBlockByNumber(context.Background(), rpc.BlockNumber(block.Number().Int64()), &cn.TraceConfig{
			Tracer: &fct,
		})
		if err != nil {
			return blockchain.ChainEvent{}, err
		}
		for _, r := range results {
			// TODO-ChainDataFetcher Assume that the input parameters are valid always.
			internalTraces = append(internalTraces, r.Result.(*vm.InternalTxTrace))
		}
	}

	return blockchain.ChainEvent{
		Block:            block,
		Hash:             block.Hash(),
		Receipts:         receipts,
		Logs:             logs,
		InternalTxTraces: internalTraces,
	}, nil
}

func (f *ChainDataFetcher) Components() []interface{} {
	return nil
}

func (f *ChainDataFetcher) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			f.blockchain = v
		case []rpc.API:
			for _, a := range v {
				switch s := a.Service.(type) {
				case *api.PublicBlockChainAPI:
					f.repo.SetComponent(s)
				case *cn.PrivateDebugAPI:
					f.debugAPI = s
				}
			}
		}
	}
}

func (f *ChainDataFetcher) handleRequestByType(reqType requestType, shouldUpdateCheckpoint bool, ev blockchain.ChainEvent) {
	// TODO-ChainDataFetcher parallelize handling data
	if checkRequestType(reqType, requestTypeTransaction) {
		f.retryFunc(f.repo.InsertTransactions)(ev)
	}
	if checkRequestType(reqType, requestTypeTokenTransfer) {
		f.retryFunc(f.repo.InsertTokenTransfers)(ev)
	}
	if checkRequestType(reqType, requestTypeContracts) {
		f.retryFunc(f.repo.InsertContracts)(ev)
	}
	if checkRequestType(reqType, requestTypeTraces) {
		f.retryFunc(f.repo.InsertTraceResults)(ev)
	}
	if shouldUpdateCheckpoint {
		f.updateCheckpoint(ev.Block.Number().Int64())
	}
}

func (f *ChainDataFetcher) handleRequest() {
	f.wg.Add(1)
	defer f.wg.Done()
	for {
		select {
		case <-f.stopCh:
			logger.Info("handleRequest is stopped")
			return
		case ev := <-f.chainCh:
			f.handleRequestByType(requestTypeAll, true, ev)
		case req := <-f.reqCh:
			ev, err := f.makeChainEvent(req.blockNumber)
			if err != nil {
				// TODO-ChainDataFetcher handle error
				logger.Error("making chain event is failed", "err", err)
				break
			}
			f.handleRequestByType(req.reqType, req.shouldUpdateCheckpoint, ev)
		}
	}
}

func (f *ChainDataFetcher) updateCheckpoint(num int64) error {
	f.checkpointMu.Lock()
	defer f.checkpointMu.Unlock()
	f.checkpointMap[num] = struct{}{}

	updated := false
	newCheckpoint := f.checkpoint
	for {
		if _, ok := f.checkpointMap[newCheckpoint]; !ok {
			break
		}
		delete(f.checkpointMap, newCheckpoint)
		newCheckpoint++
		updated = true
	}

	if updated {
		f.checkpoint = newCheckpoint
		return f.repo.WriteCheckpoint(newCheckpoint)
	}
	return nil
}

func (f *ChainDataFetcher) retryFunc(insert func(blockchain.ChainEvent) error) func(blockchain.ChainEvent) {
	return func(event blockchain.ChainEvent) {
		i := 0
		for err := insert(event); err != nil; {
			select {
			case <-f.stopCh:
				return
			default:
				i++
				logger.Warn("retrying...", "blockNumber", event.Block.NumberU64(), "retryCount", i)
				time.Sleep(DBInsertRetryInterval)
			}
		}
	}
}
