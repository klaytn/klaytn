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
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"sync"
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

type ChainDataFetcher struct {
	config *ChainDataFetcherConfig

	blockchain *blockchain.BlockChain

	chainCh  chan blockchain.ChainEvent
	chainSub event.Subscription

	reqCh  chan *request // TODO-ChainDataFetcher add logic to insert new requests from APIs to this channel
	resCh  chan *response
	stopCh chan struct{}

	numHandlers int

	checkpoint    int64
	checkpointMap map[int64]struct{}

	wg sync.WaitGroup

	repo Repository
}

func NewChainDataFetcher(ctx *node.ServiceContext, cfg *ChainDataFetcherConfig) (*ChainDataFetcher, error) {
	repo, err := kas.NewRepository(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		logger.Error("Failed to create new Repository", "err", err, "user", cfg.DBUser, "password", cfg.DBPassword, "host", cfg.DBHost, "port", cfg.DBPort, "name", cfg.DBName)
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
		resCh:         make(chan *response, cfg.JobChannelSize),
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
	// TODO-ChainDataFetcher add APIs to start or stop chaindata fetcher
	return []rpc.API{}
}

func (f *ChainDataFetcher) Start(server p2p.Server) error {
	// launch multiple goroutines to handle new blocks
	for i := 0; i < f.numHandlers; i++ {
		go f.handleRequest()
	}

	// subscribe chain head event
	f.chainSub = f.blockchain.SubscribeChainEvent(f.chainCh)
	go f.reqLoop()
	go f.resLoop()

	return nil
}

func (f *ChainDataFetcher) Stop() error {
	f.chainSub.Unsubscribe()
	close(f.stopCh)

	logger.Info("wait for all goroutines to be terminated...")
	f.wg.Wait()
	logger.Info("terminated all goroutines for chaindatafetcher")
	return nil
}

func (f *ChainDataFetcher) Components() []interface{} {
	return nil
}

func (f *ChainDataFetcher) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			f.blockchain = v
		}
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
		case req := <-f.reqCh:
			res := &response{
				reqType:     requestTypeTransaction,
				blockNumber: req.event.Block.Number(),
				err:         nil,
			}

			res.err = f.repo.InsertTransactions(req.event)
			// TODO-ChainDataFetcher insert other types of data
			f.resCh <- res
		}
	}
}

func (f *ChainDataFetcher) reqLoop() {
	f.wg.Add(1)
	defer f.wg.Done()
	for {
		select {
		case <-f.stopCh:
			logger.Info("stopped reqLoop for chaindatafetcher")
			return
		case ev := <-f.chainCh:
			f.reqCh <- &request{
				reqType: requestTypeTransaction,
				event:   ev,
			}
		}
	}
}

func (f *ChainDataFetcher) resLoop() {
	f.wg.Add(1)
	defer f.wg.Done()
	for {
		select {
		case <-f.stopCh:
			logger.Info("stopped resLoop for chaindatafetcher")
			return
		case res := <-f.resCh:
			if res.err != nil {
				logger.Error("db insertion is failed", "blockNumber", res.blockNumber, "reqType", res.reqType, "err", res.err)
				// TODO-ChainDataFetcher add retry logic when data insertion is failed
			} else {
				if err := f.updateCheckpoint(res.blockNumber.Int64()); err != nil {
					logger.Error("Failed to update checkpoint", "checkpoint", res.blockNumber.Int64())
				}
			}
		}
	}
}

func (f *ChainDataFetcher) updateCheckpoint(num int64) error {
	f.checkpointMap[num] = struct{}{}

	oldCheckpoint := f.checkpoint
	newCheckpoint := oldCheckpoint
	for {
		if _, ok := f.checkpointMap[newCheckpoint]; !ok {
			break
		}
		delete(f.checkpointMap, newCheckpoint)
		newCheckpoint++
	}

	if oldCheckpoint != newCheckpoint {
		f.checkpoint = newCheckpoint
		return f.repo.WriteCheckpoint(newCheckpoint)
	}
	return nil
}
