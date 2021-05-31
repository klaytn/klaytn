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
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas"
	cfTypes "github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/rcrowley/go-metrics"
)

const (
	stopped = uint32(iota)
	running
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)
var errUnsupportedMode = errors.New("the given chaindatafetcher mode is not supported")
var errMaxRetryExceeded = errors.New("the number of retries is exceeded over max")

//go:generate mockgen -destination=./mocks/blockchain_mock.go -package=mocks github.com/klaytn/klaytn/datasync/chaindatafetcher BlockChain
type BlockChain interface {
	SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription
	CurrentHeader() *types.Header
	GetBlockByNumber(number uint64) *types.Block
	GetReceiptsByBlockHash(blockHash common.Hash) types.Receipts
}

type ChainDataFetcher struct {
	config *ChainDataFetcherConfig

	blockchain BlockChain
	debugAPI   *cn.PrivateDebugAPI

	chainCh  chan blockchain.ChainEvent
	chainSub event.Subscription

	reqCh  chan *cfTypes.Request // TODO-ChainDataFetcher add logic to insert new requests from APIs to this channel
	stopCh chan struct{}

	numHandlers int

	checkpointMu  sync.RWMutex
	checkpoint    int64
	checkpointMap map[int64]struct{}

	wg sync.WaitGroup

	repo         Repository
	checkpointDB CheckpointDB
	setters      []ComponentSetter

	fetchingStarted      uint32
	fetchingStopCh       chan struct{}
	fetchingWg           sync.WaitGroup
	rangeFetchingStarted uint32
	rangeFetchingStopCh  chan struct{}
	rangeFetchingWg      sync.WaitGroup
}

func NewChainDataFetcher(ctx *node.ServiceContext, cfg *ChainDataFetcherConfig) (*ChainDataFetcher, error) {
	var (
		repo         Repository
		checkpointDB CheckpointDB
		setters      []ComponentSetter
		err          error
	)
	switch cfg.Mode {
	case ModeKAS:
		repo, checkpointDB, setters, err = getKasComponents(cfg.KasConfig)
		if err != nil {
			return nil, err
		}
	case ModeKafka:
		repo, checkpointDB, setters, err = getKafkaComponents(cfg.KafkaConfig)
		if err != nil {
			return nil, err
		}
	default:
		logger.Error("the chaindatafetcher mode is not supported", "mode", cfg.Mode)
		return nil, errUnsupportedMode
	}
	return &ChainDataFetcher{
		config:        cfg,
		chainCh:       make(chan blockchain.ChainEvent, cfg.BlockChannelSize),
		reqCh:         make(chan *cfTypes.Request, cfg.JobChannelSize),
		stopCh:        make(chan struct{}),
		numHandlers:   cfg.NumHandlers,
		checkpointMap: make(map[int64]struct{}),
		repo:          repo,
		checkpointDB:  checkpointDB,
		setters:       setters,
	}, nil
}

func getKasComponents(cfg *kas.KASConfig) (Repository, CheckpointDB, []ComponentSetter, error) {
	repo, err := kas.NewRepository(cfg)
	if err != nil {
		return nil, nil, nil, err
	}
	return repo, repo, []ComponentSetter{repo}, nil
}

func getKafkaComponents(cfg *kafka.KafkaConfig) (Repository, CheckpointDB, []ComponentSetter, error) {
	repo, err := kafka.NewRepository(cfg)
	if err != nil {
		return nil, nil, nil, err
	}
	checkpointDB := kafka.NewCheckpointDB()
	return repo, checkpointDB, []ComponentSetter{repo, checkpointDB}, nil
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
			logger.Error("start fetching is failed", "err", err)
			return err
		}
	}
	logger.Info("chaindata fetcher is started", "numHandlers", f.numHandlers)
	return nil
}

func (f *ChainDataFetcher) Stop() error {
	f.stopFetching()
	f.stopRangeFetching()
	logger.Info("wait for all goroutines to be terminated...", "numGoroutines", f.config.NumHandlers)
	close(f.stopCh)
	f.wg.Wait()
	logger.Info("chaindata fetcher is stopped")
	return nil
}

func (f *ChainDataFetcher) sendRequests(startBlock, endBlock uint64, reqType cfTypes.RequestType, shouldUpdateCheckpoint bool, stopCh chan struct{}) {
	logger.Info("sending requests is started", "startBlock", startBlock, "endBlock", endBlock)
	for i := startBlock; i <= endBlock; i++ {
		select {
		case <-stopCh:
			logger.Info("stopped making requests", "startBlock", startBlock, "endBlock", endBlock, "stoppedBlock", i)
			return
		case f.reqCh <- cfTypes.NewRequest(reqType, shouldUpdateCheckpoint, i):
		}
	}
	logger.Info("sending requests is finished", "startBlock", startBlock, "endBlock", endBlock)
}

func (f *ChainDataFetcher) startFetching() error {
	if !atomic.CompareAndSwapUint32(&f.fetchingStarted, stopped, running) {
		return errors.New("fetching is already started")
	}

	// subscribe chain event in order to handle new blocks.
	f.chainSub = f.blockchain.SubscribeChainEvent(f.chainCh)
	checkpoint := uint64(f.checkpoint)
	currentBlock := f.blockchain.CurrentHeader().Number.Uint64()

	f.fetchingStopCh = make(chan struct{})
	f.fetchingWg.Add(1)

	// lanuch a goroutine to handle from checkpoint to the head block.
	go func() {
		defer f.fetchingWg.Done()
		switch f.config.Mode {
		case ModeKAS:
			f.sendRequests(uint64(f.checkpoint), currentBlock, cfTypes.RequestTypeAll, true, f.fetchingStopCh)
		case ModeKafka:
			f.sendRequests(uint64(f.checkpoint), currentBlock, cfTypes.RequestTypeGroupAll, true, f.fetchingStopCh)
		default:
			logger.Error("the chaindatafetcher mode is not supported", "mode", f.config.Mode, "checkpoint", f.checkpoint, "currentBlock", currentBlock)
		}
	}()
	logger.Info("fetching is started", "startedCheckpoint", checkpoint, "currentBlock", currentBlock)
	return nil
}

func (f *ChainDataFetcher) stopFetching() error {
	if !atomic.CompareAndSwapUint32(&f.fetchingStarted, running, stopped) {
		return errors.New("fetching is not running")
	}

	f.chainSub.Unsubscribe()
	close(f.fetchingStopCh)
	f.fetchingWg.Wait()
	logger.Info("fetching is stopped")
	return nil
}

func (f *ChainDataFetcher) startRangeFetching(startBlock, endBlock uint64, reqType cfTypes.RequestType) error {
	if !atomic.CompareAndSwapUint32(&f.rangeFetchingStarted, stopped, running) {
		return errors.New("range fetching is already started")
	}

	f.rangeFetchingStopCh = make(chan struct{})
	f.rangeFetchingWg.Add(1)
	go func() {
		defer f.rangeFetchingWg.Done()
		f.sendRequests(startBlock, endBlock, reqType, false, f.rangeFetchingStopCh)
		atomic.StoreUint32(&f.rangeFetchingStarted, stopped)
	}()
	logger.Info("range fetching is started", "startBlock", startBlock, "endBlock", endBlock)
	return nil
}

func (f *ChainDataFetcher) stopRangeFetching() error {
	if !atomic.CompareAndSwapUint32(&f.rangeFetchingStarted, running, stopped) {
		return errors.New("range fetching is not running")
	}
	close(f.rangeFetchingStopCh)
	f.rangeFetchingWg.Wait()
	logger.Info("range fetching is stopped")
	return nil
}

func (f *ChainDataFetcher) makeChainEvent(blockNumber uint64) (blockchain.ChainEvent, error) {
	var logs []*types.Log
	block := f.blockchain.GetBlockByNumber(blockNumber)
	if block == nil {
		return blockchain.ChainEvent{}, fmt.Errorf("GetBlockByNumber is failed. blockNumber=%v", blockNumber)
	}
	receipts := f.blockchain.GetReceiptsByBlockHash(block.Hash())
	if receipts == nil {
		return blockchain.ChainEvent{}, fmt.Errorf("GetReceiptsByBlockHash is failed. blockNumber=%v", blockNumber)
	}
	for _, r := range receipts {
		logs = append(logs, r.Logs...)
	}
	var internalTraces []*vm.InternalTxTrace
	if block.Transactions().Len() > 0 {
		fct := "fastCallTracer"
		timeout := "24h"
		results, err := f.debugAPI.TraceBlockByNumber(context.Background(), rpc.BlockNumber(block.Number().Int64()), &cn.TraceConfig{
			Tracer:  &fct,
			Timeout: &timeout,
		})
		if err != nil {
			traceAPIErrorCounter.Inc(1)
			logger.Error("Failed to call trace block by number", "err", err, "blockNumber", block.NumberU64())
			return blockchain.ChainEvent{}, err
		}
		for _, r := range results {
			if r.Result != nil {
				internalTraces = append(internalTraces, r.Result.(*vm.InternalTxTrace))
			} else {
				traceAPIErrorCounter.Inc(1)
				logger.Error("the trace result is nil", "err", r.Error, "blockNumber", blockNumber)
				internalTraces = append(internalTraces, &vm.InternalTxTrace{Value: "0x0", Calls: []*vm.InternalTxTrace{}})
			}
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

func (f *ChainDataFetcher) setDebugAPI(apis []rpc.API) {
	for _, a := range apis {
		switch s := a.Service.(type) {
		case *cn.PrivateDebugAPI:
			f.debugAPI = s
		}
	}
}

func (f *ChainDataFetcher) setCheckpoint() {
	checkpoint, err := f.checkpointDB.ReadCheckpoint()
	if err != nil {
		logger.Crit("ReadCheckpoint is failed", "err", err)
	}

	if checkpoint == 0 {
		checkpoint = f.blockchain.CurrentHeader().Number.Int64()
	}
	f.checkpoint = checkpoint
	logger.Info("Chaindatafetcher initial checkpoint is set", "checkpoint", f.checkpoint)
}

func (f *ChainDataFetcher) setComponent(component interface{}) {
	switch v := component.(type) {
	case *blockchain.BlockChain:
		f.blockchain = v
	case []rpc.API:
		f.setDebugAPI(v)
	}
}

func (f *ChainDataFetcher) SetComponents(components []interface{}) {
	for _, component := range components {
		f.setComponent(component)
		for _, setter := range f.setters {
			setter.SetComponent(component)
		}
	}
	f.setCheckpoint()
}

func (f *ChainDataFetcher) handleRequestByType(reqType cfTypes.RequestType, shouldUpdateCheckpoint bool, ev blockchain.ChainEvent) error {
	now := time.Now()
	// TODO-ChainDataFetcher parallelize handling data

	// iterate over all types of requests
	// - RequestTypeTransaction
	// - RequestTypeTokenTransfer
	// - RequestTypeContract
	// - RequestTypeTrace
	// - RequestTypeBlockGroup
	// - RequestTypeTraceGroup
	for targetType := cfTypes.RequestTypeTransaction; targetType < cfTypes.RequestTypeLength; targetType = targetType << 1 {
		if cfTypes.CheckRequestType(reqType, targetType) {
			if err := f.updateInsertionTimeGauge(f.retryFunc(f.repo.HandleChainEvent))(ev, targetType); err != nil {
				logger.Error("handling chain event is failed", "blockNumber", ev.Block.NumberU64(), "err", err, "reqType", reqType, "targetType", targetType)
				return err
			}
		}
	}
	elapsed := time.Since(now)
	totalInsertionTimeGauge.Update(elapsed.Milliseconds())

	if shouldUpdateCheckpoint {
		f.updateCheckpoint(ev.Block.Number().Int64())
	}
	handledBlockNumberGauge.Update(ev.Block.Number().Int64())
	return nil
}

func (f *ChainDataFetcher) resetChainCh() {
	for {
		select {
		case <-f.chainCh:
		default:
			return
		}
	}
}

func (f *ChainDataFetcher) resetRequestCh() {
	for {
		select {
		case <-f.reqCh:
		default:
			return
		}
	}
}

func (f *ChainDataFetcher) pause() {
	f.stopFetching()
	f.stopRangeFetching()
	f.resetChainCh()
	f.resetRequestCh()
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
			numChainEventGauge.Update(int64(len(f.chainCh)))
			var err error
			switch f.config.Mode {
			case ModeKAS:
				err = f.handleRequestByType(cfTypes.RequestTypeAll, true, ev)
			case ModeKafka:
				err = f.handleRequestByType(cfTypes.RequestTypeGroupAll, true, ev)
			default:
				logger.Error("the chaindatafetcher mode is not supported", "mode", f.config.Mode, "blockNumber", ev.Block.NumberU64())
			}

			if err != nil && err == errMaxRetryExceeded {
				logger.Error("the chaindatafetcher reaches the maximum retries. it pauses fetching and clear the channels", "blockNum", ev.Block.NumberU64())
				f.pause()
			}
		case req := <-f.reqCh:
			numRequestsGauge.Update(int64(len(f.reqCh)))
			ev, err := f.makeChainEvent(req.BlockNumber)
			if err != nil {
				// TODO-ChainDataFetcher handle error
				logger.Error("making chain event is failed", "err", err)
				break
			}
			err = f.handleRequestByType(req.ReqType, req.ShouldUpdateCheckpoint, ev)
			if err != nil && err == errMaxRetryExceeded {
				logger.Error("the chaindatafetcher reaches the maximum retries. it pauses fetching and clear the channels", "blockNum", ev.Block.NumberU64())
				f.pause()
			}
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
		checkpointGauge.Update(f.checkpoint)
		return f.checkpointDB.WriteCheckpoint(newCheckpoint)
	}
	return nil
}

func getInsertionTimeGauge(reqType cfTypes.RequestType) metrics.Gauge {
	switch reqType {
	case cfTypes.RequestTypeTransaction:
		return txsInsertionTimeGauge
	case cfTypes.RequestTypeTokenTransfer:
		return tokenTransfersInsertionTimeGauge
	case cfTypes.RequestTypeContract:
		return contractsInsertionTimeGauge
	case cfTypes.RequestTypeTrace:
		return tracesInsertionTimeGauge
	case cfTypes.RequestTypeBlockGroup:
		return blockGroupInsertionTimeGauge
	case cfTypes.RequestTypeTraceGroup:
		return traceGroupInsertionTimeGauge
	default:
		logger.Warn("the request type is not supported", "type", reqType)
		return metrics.NilGauge{}
	}
}

func (f *ChainDataFetcher) updateInsertionTimeGauge(insert HandleChainEventFn) HandleChainEventFn {
	return func(chainEvent blockchain.ChainEvent, reqType cfTypes.RequestType) error {
		now := time.Now()
		if err := insert(chainEvent, reqType); err != nil {
			return err
		}
		elapsed := time.Since(now)
		gauge := getInsertionTimeGauge(reqType)
		gauge.Update(elapsed.Milliseconds())
		return nil
	}
}

func getInsertionRetryGauge(reqType cfTypes.RequestType) metrics.Gauge {
	switch reqType {
	case cfTypes.RequestTypeTransaction:
		return txsInsertionRetryGauge
	case cfTypes.RequestTypeTokenTransfer:
		return tokenTransfersInsertionRetryGauge
	case cfTypes.RequestTypeContract:
		return contractsInsertionRetryGauge
	case cfTypes.RequestTypeTrace:
		return tracesInsertionRetryGauge
	case cfTypes.RequestTypeBlockGroup:
		return blockGroupInsertionRetryGauge
	case cfTypes.RequestTypeTraceGroup:
		return traceGroupInsertionRetryGauge
	default:
		logger.Warn("the request type is not supported", "type", reqType)
		return metrics.NilGauge{}
	}
}

func (f *ChainDataFetcher) retryFunc(insert HandleChainEventFn) HandleChainEventFn {
	return func(event blockchain.ChainEvent, reqType cfTypes.RequestType) error {
		i := 0
		for err := insert(event, reqType); err != nil; err = insert(event, reqType) {
			select {
			case <-f.stopCh:
				return err
			default:
				if i > InsertMaxRetry {
					return errMaxRetryExceeded
				}
				i++
				gauge := getInsertionRetryGauge(reqType)
				gauge.Update(int64(i))
				logger.Warn("retrying...", "blockNumber", event.Block.NumberU64(), "retryCount", i, "err", err)
				time.Sleep(InsertRetryInterval)
			}
		}
		return nil
	}
}

func (f *ChainDataFetcher) status() string {
	return fmt.Sprintf("{fetching: %v, rangeFetching: %v}", atomic.LoadUint32(&f.fetchingStarted), atomic.LoadUint32(&f.rangeFetchingStarted))
}
