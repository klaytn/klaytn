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
	"errors"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/mocks"
	cfTypes "github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
	eventMocks "github.com/klaytn/klaytn/event/mocks"
	"github.com/stretchr/testify/assert"
)

func newTestChainDataFetcher() *ChainDataFetcher {
	return &ChainDataFetcher{
		config:        DefaultChainDataFetcherConfig,
		chainCh:       make(chan blockchain.ChainEvent),
		reqCh:         make(chan *cfTypes.Request),
		stopCh:        make(chan struct{}),
		numHandlers:   3,
		checkpoint:    0,
		checkpointMap: make(map[int64]struct{}),
		repo:          nil,
	}
}

func TestChainDataFetcher_Success_StartAndStop(t *testing.T) {
	fetcher := newTestChainDataFetcher()
	fetcher.config.NoDefaultStart = true
	// Start launches several goroutines.
	assert.NoError(t, fetcher.Start(nil))

	// Stop waits the all goroutines to be terminated.
	assert.NoError(t, fetcher.Stop())
}

func TestChainDataFetcher_Success_sendRequests(t *testing.T) {
	fetcher := newTestChainDataFetcher()
	stopCh := make(chan struct{})

	startBlock := uint64(0)
	endBlock := uint64(10)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fetcher.sendRequests(startBlock, endBlock, cfTypes.RequestTypeAll, false, stopCh)
	}()

	// take all the items from the reqCh and check them.
	for i := startBlock; i <= endBlock; i++ {
		r := <-fetcher.reqCh
		assert.Equal(t, i, r.BlockNumber)
		assert.Equal(t, cfTypes.RequestTypeAll, r.ReqType)
		assert.Equal(t, false, r.ShouldUpdateCheckpoint)
	}
	wg.Wait()
}

func TestChainDataFetcher_Success_sendRequestsStop(t *testing.T) {
	fetcher := newTestChainDataFetcher()
	stopCh := make(chan struct{})

	startBlock := uint64(0)
	endBlock := uint64(10)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fetcher.sendRequests(startBlock, endBlock, cfTypes.RequestTypeAll, false, stopCh)
	}()

	stopCh <- struct{}{}
	wg.Wait()
}

func TestChainDataFetcher_Success_fetchingStartAndStop(t *testing.T) {
	fetcher := newTestChainDataFetcher()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bc := mocks.NewMockBlockChain(ctrl)
	bc.EXPECT().SubscribeChainEvent(gomock.Any()).Return(nil).Times(1)
	// TODO-ChainDataFetcher the below statement is not working, so find out why.
	//bc.EXPECT().SubscribeChainEvent(gomock.Eq(fetcher.chainCh)).Return(nil).Times(1)
	bc.EXPECT().CurrentHeader().Return(&types.Header{Number: big.NewInt(1)}).Times(1)

	fetcher.blockchain = bc
	assert.NoError(t, fetcher.startFetching())

	sub := eventMocks.NewMockSubscription(ctrl)
	sub.EXPECT().Unsubscribe().Times(1)
	fetcher.chainSub = sub
	assert.NoError(t, fetcher.stopFetching())
}

func TestChainDataFetcher_Success_rangeFetchingStartAndStop(t *testing.T) {
	fetcher := newTestChainDataFetcher()
	// start range fetching and the method is waiting for the items to be taken from reqCh
	assert.NoError(t, fetcher.startRangeFetching(0, 10, cfTypes.RequestTypeAll))

	// take only parts of the requests
	<-fetcher.reqCh
	<-fetcher.reqCh
	<-fetcher.reqCh

	// stop fetching while waiting
	assert.NoError(t, fetcher.stopRangeFetching())
}

func TestChainDataFetcher_Success_rangeFetchingStartAndFinishedAlready(t *testing.T) {
	fetcher := newTestChainDataFetcher()
	assert.NoError(t, fetcher.startRangeFetching(0, 0, cfTypes.RequestTypeAll))
	// skip the request
	<-fetcher.reqCh

	// sleep to finish sending the request
	time.Sleep(100 * time.Millisecond)

	// already finished range fetching
	err := fetcher.stopRangeFetching()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "range fetching is not running"))
}

func TestChainDataFetcher_updateCheckpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockCheckpointDB(ctrl)

	fetcher := &ChainDataFetcher{
		checkpoint:    0,
		checkpointMap: make(map[int64]struct{}),
		checkpointDB:  m,
	}

	// update checkpoint as follows.
	// done order: 1, 0, 2, 3, 5, 7, 9, 8, 4, 6, 10
	// checkpoint: 0, 2, 3, 4, 4, 4, 4, 4, 6, 10, 11
	assert.NoError(t, fetcher.updateCheckpoint(1))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(2)))
	assert.NoError(t, fetcher.updateCheckpoint(0))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(3)))
	assert.NoError(t, fetcher.updateCheckpoint(2))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(4)))
	assert.NoError(t, fetcher.updateCheckpoint(3))
	assert.NoError(t, fetcher.updateCheckpoint(5))
	assert.NoError(t, fetcher.updateCheckpoint(7))
	assert.NoError(t, fetcher.updateCheckpoint(9))
	assert.NoError(t, fetcher.updateCheckpoint(8))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(6)))
	assert.NoError(t, fetcher.updateCheckpoint(4))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(10)))
	assert.NoError(t, fetcher.updateCheckpoint(6))

	m.EXPECT().WriteCheckpoint(gomock.Eq(int64(11)))
	assert.NoError(t, fetcher.updateCheckpoint(10))
}

func TestChainDataFetcher_retryFunc(t *testing.T) {
	fetcher := &ChainDataFetcher{}
	header := &types.Header{
		Number: big.NewInt(1),
	}
	ev := blockchain.ChainEvent{
		Block: types.NewBlockWithHeader(header),
	}

	i := 0
	f := func(event blockchain.ChainEvent, reqType cfTypes.RequestType) error {
		i++
		if i == 5 {
			return nil
		} else {
			return errors.New("test")
		}
	}

	assert.NoError(t, fetcher.retryFunc(f)(ev, cfTypes.RequestTypeAll))
	assert.True(t, i == 5)
}

func TestChainDataFetcher_handleRequestByType_WhileRetrying(t *testing.T) {
	fetcher := &ChainDataFetcher{
		config:        &ChainDataFetcherConfig{NumHandlers: 1}, // prevent panic with nil reference
		checkpoint:    1,
		checkpointMap: make(map[int64]struct{}), // in order to call CheckpointDB WriteCheckpoint method
		stopCh:        make(chan struct{}),      // in order to stop retrying
	}
	header1 := &types.Header{Number: big.NewInt(1)} // next block to be handled is 1
	block1 := blockchain.ChainEvent{Block: types.NewBlockWithHeader(header1)}
	header2 := &types.Header{Number: big.NewInt(2)} // next block to be handled is 2
	block2 := blockchain.ChainEvent{Block: types.NewBlockWithHeader(header2)}
	testError := errors.New("test-error") // fake error to call retrying infinitely

	// set up mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo, checkpointDB := mocks.NewMockRepository(ctrl), mocks.NewMockCheckpointDB(ctrl)

	// update checkpoint to 2
	precall := mockRepo.EXPECT().HandleChainEvent(gomock.Any(), gomock.Any()).Return(nil).Times(6)
	checkpointDB.EXPECT().WriteCheckpoint(gomock.Eq(int64(2))).Return(nil).Times(1)

	// retrying indefinitely
	mockRepo.EXPECT().HandleChainEvent(gomock.Any(), gomock.Any()).Return(testError).AnyTimes().After(precall)

	wg := sync.WaitGroup{}
	wg.Add(1)
	// trigger stop function after 3 seconds
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		fetcher.Stop()
	}()

	fetcher.repo, fetcher.checkpointDB = mockRepo, checkpointDB
	fetcher.handleRequestByType(cfTypes.RequestTypeAll, true, block1)
	fetcher.handleRequestByType(cfTypes.RequestTypeAll, true, block2)

	wg.Wait()
	assert.Equal(t, int64(2), fetcher.checkpoint)
}

func TestChainDataFetcher_setComponents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bc, db := mocks.NewMockBlockChain(ctrl), mocks.NewMockCheckpointDB(ctrl)

	fetcher := &ChainDataFetcher{
		blockchain:   bc,
		checkpointDB: db,
	}

	// if checkpoint no exist
	testBlockNumber := new(big.Int).SetInt64(5)
	db.EXPECT().ReadCheckpoint().Return(int64(0), nil).Times(1)
	testHeader := &types.Header{Number: testBlockNumber}
	bc.EXPECT().CurrentHeader().Return(testHeader).Times(1)
	fetcher.setCheckpoint()
	assert.Equal(t, testBlockNumber.Int64(), fetcher.checkpoint)

	// if checkpoint exist
	testCheckpoint := int64(10)
	db.EXPECT().ReadCheckpoint().Return(int64(10), nil).Times(1)
	fetcher.setCheckpoint()
	assert.Equal(t, testCheckpoint, fetcher.checkpoint)
}
