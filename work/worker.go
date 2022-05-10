// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from miner/worker.go (2018/06/04).
// Modified and improved for the klaytn development.

package work

import (
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	klaytnmetrics "github.com/klaytn/klaytn/metrics"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/rcrowley/go-metrics"
)

const (
	resultQueueSize  = 10
	miningLogAtDepth = 5

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10
	// chainSideChanSize is the size of channel listening to ChainSideEvent.
	chainSideChanSize = 10
	// maxResendSize is the size of resending transactions to peer in order to prevent the txs from missing.
	maxResendTxSize = 1000
)

var (
	// Metrics for miner
	timeLimitReachedCounter = metrics.NewRegisteredCounter("miner/timelimitreached", nil)
	tooLongTxCounter        = metrics.NewRegisteredCounter("miner/toolongtx", nil)
	ResultChGauge           = metrics.NewRegisteredGauge("miner/resultch", nil)
	resentTxGauge           = metrics.NewRegisteredGauge("miner/tx/resend/gauge", nil)
	usedAllTxsCounter       = metrics.NewRegisteredCounter("miner/usedalltxs", nil)
	checkedTxsGauge         = metrics.NewRegisteredGauge("miner/checkedtxs", nil)
	tCountGauge             = metrics.NewRegisteredGauge("miner/tcount", nil)
	nonceTooLowTxsGauge     = metrics.NewRegisteredGauge("miner/nonce/low/txs", nil)
	nonceTooHighTxsGauge    = metrics.NewRegisteredGauge("miner/nonce/high/txs", nil)
	gasLimitReachedTxsGauge = metrics.NewRegisteredGauge("miner/limitreached/gas/txs", nil)
	strangeErrorTxsCounter  = metrics.NewRegisteredCounter("miner/strangeerror/txs", nil)

	blockMiningTimer          = klaytnmetrics.NewRegisteredHybridTimer("miner/block/mining/time", nil)
	blockMiningExecuteTxTimer = klaytnmetrics.NewRegisteredHybridTimer("miner/block/execute/time", nil)
	blockMiningCommitTxTimer  = klaytnmetrics.NewRegisteredHybridTimer("miner/block/commit/time", nil)
	blockMiningFinalizeTimer  = klaytnmetrics.NewRegisteredHybridTimer("miner/block/finalize/time", nil)

	accountReadTimer   = klaytnmetrics.NewRegisteredHybridTimer("miner/block/account/reads", nil)
	accountHashTimer   = klaytnmetrics.NewRegisteredHybridTimer("miner/block/account/hashes", nil)
	accountUpdateTimer = klaytnmetrics.NewRegisteredHybridTimer("miner/block/account/updates", nil)
	accountCommitTimer = klaytnmetrics.NewRegisteredHybridTimer("miner/block/account/commits", nil)

	storageReadTimer   = klaytnmetrics.NewRegisteredHybridTimer("miner/block/storage/reads", nil)
	storageHashTimer   = klaytnmetrics.NewRegisteredHybridTimer("miner/block/storage/hashes", nil)
	storageUpdateTimer = klaytnmetrics.NewRegisteredHybridTimer("miner/block/storage/updates", nil)
	storageCommitTimer = klaytnmetrics.NewRegisteredHybridTimer("miner/block/storage/commits", nil)

	snapshotAccountReadTimer = metrics.NewRegisteredTimer("miner/snapshot/account/reads", nil)
	snapshotStorageReadTimer = metrics.NewRegisteredTimer("miner/snapshot/storage/reads", nil)
	snapshotCommitTimer      = metrics.NewRegisteredTimer("miner/snapshot/commits", nil)
)

// Agent can register themself with the worker
type Agent interface {
	Work() chan<- *Task
	SetReturnCh(chan<- *Result)
	Stop()
	Start()
	GetHashRate() int64
}

// Task is the workers current environment and holds
// all of the current state information
type Task struct {
	config *params.ChainConfig
	signer types.Signer

	stateMu sync.RWMutex   // protects state
	state   *state.StateDB // apply state changes here
	tcount  int            // tx count in cycle

	Block *types.Block // the new block

	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt

	createdAt time.Time
}

type Result struct {
	Task  *Task
	Block *types.Block
}

// worker is the main object which takes care of applying messages to the new state
type worker struct {
	config *params.ChainConfig
	engine consensus.Engine

	mu sync.Mutex

	// update loop
	mux          *event.TypeMux
	txsCh        chan blockchain.NewTxsEvent
	txsSub       event.Subscription
	chainHeadCh  chan blockchain.ChainHeadEvent
	chainHeadSub event.Subscription
	chainSideCh  chan blockchain.ChainSideEvent
	chainSideSub event.Subscription
	wg           sync.WaitGroup

	agents map[Agent]struct{}
	recv   chan *Result

	backend Backend
	chain   BlockChain
	proc    blockchain.Validator
	chainDB database.DBManager

	extra []byte

	currentMu  sync.Mutex
	current    *Task
	rewardbase common.Address

	snapshotMu    sync.RWMutex
	snapshotBlock *types.Block
	snapshotState *state.StateDB

	// atomic status counters
	mining int32
	atWork int32

	nodetype common.ConnType
}

func newWorker(config *params.ChainConfig, engine consensus.Engine, rewardbase common.Address, backend Backend, mux *event.TypeMux, nodetype common.ConnType, TxResendUseLegacy bool) *worker {
	worker := &worker{
		config:      config,
		engine:      engine,
		backend:     backend,
		mux:         mux,
		txsCh:       make(chan blockchain.NewTxsEvent, txChanSize),
		chainHeadCh: make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
		chainSideCh: make(chan blockchain.ChainSideEvent, chainSideChanSize),
		chainDB:     backend.ChainDB(),
		recv:        make(chan *Result, resultQueueSize),
		chain:       backend.BlockChain(),
		proc:        backend.BlockChain().Validator(),
		agents:      make(map[Agent]struct{}),
		nodetype:    nodetype,
		rewardbase:  rewardbase,
	}

	// Subscribe NewTxsEvent for tx pool
	worker.txsSub = backend.TxPool().SubscribeNewTxsEvent(worker.txsCh)
	// Subscribe events for blockchain
	worker.chainHeadSub = backend.BlockChain().SubscribeChainHeadEvent(worker.chainHeadCh)
	worker.chainSideSub = backend.BlockChain().SubscribeChainSideEvent(worker.chainSideCh)
	go worker.update()

	go worker.wait(TxResendUseLegacy)
	return worker
}

func (self *worker) setExtra(extra []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.extra = extra
}

func (self *worker) pending() (*types.Block, *state.StateDB) {
	if atomic.LoadInt32(&self.mining) == 0 {
		// return a snapshot to avoid contention on currentMu mutex
		self.snapshotMu.RLock()
		defer self.snapshotMu.RUnlock()
		return self.snapshotBlock, self.snapshotState.Copy()
	}

	self.currentMu.Lock()
	defer self.currentMu.Unlock()
	self.current.stateMu.Lock()
	defer self.current.stateMu.Unlock()
	return self.current.Block, self.current.state.Copy()
}

func (self *worker) pendingBlock() *types.Block {
	if atomic.LoadInt32(&self.mining) == 0 {
		// return a snapshot to avoid contention on currentMu mutex
		self.snapshotMu.RLock()
		defer self.snapshotMu.RUnlock()
		return self.snapshotBlock
	}

	self.currentMu.Lock()
	defer self.currentMu.Unlock()
	return self.current.Block
}

func (self *worker) start() {
	self.mu.Lock()
	defer self.mu.Unlock()

	atomic.StoreInt32(&self.mining, 1)

	// istanbul BFT
	if istanbul, ok := self.engine.(consensus.Istanbul); ok {
		istanbul.Start(self.chain, self.chain.CurrentBlock, self.chain.HasBadBlock)
	}

	// spin up agents
	for agent := range self.agents {
		agent.Start()
	}
}

func (self *worker) stop() {
	self.wg.Wait()

	self.mu.Lock()
	defer self.mu.Unlock()
	if atomic.LoadInt32(&self.mining) == 1 {
		for agent := range self.agents {
			agent.Stop()
		}
	}

	// istanbul BFT
	if istanbul, ok := self.engine.(consensus.Istanbul); ok {
		istanbul.Stop()
	}

	atomic.StoreInt32(&self.mining, 0)
	atomic.StoreInt32(&self.atWork, 0)
}

func (self *worker) register(agent Agent) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.agents[agent] = struct{}{}
	agent.SetReturnCh(self.recv)
}

func (self *worker) unregister(agent Agent) {
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.agents, agent)
	agent.Stop()
}

func (self *worker) handleTxsCh(quitByErr chan bool) {
	defer self.txsSub.Unsubscribe()

	for {
		select {
		// Handle NewTxsEvent
		case <-self.txsCh:
			if atomic.LoadInt32(&self.mining) != 0 {
				// If we're mining, but nothing is being processed, wake on new transactions
				if self.config.Clique != nil && self.config.Clique.Period == 0 {
					self.commitNewWork()
				}
			}

		case <-quitByErr:
			return
		}
	}
}

func (self *worker) update() {
	defer self.chainHeadSub.Unsubscribe()
	defer self.chainSideSub.Unsubscribe()

	quitByErr := make(chan bool, 1)
	go self.handleTxsCh(quitByErr)

	for {
		// A real event arrived, process interesting content
		select {
		// Handle ChainHeadEvent
		case <-self.chainHeadCh:
			// istanbul BFT
			if h, ok := self.engine.(consensus.Handler); ok {
				h.NewChainHead()
			}
			self.commitNewWork()

			// TODO-Klaytn-Issue264 If we are using istanbul BFT, then we always have a canonical chain.
			//         Later we may be able to refine below code.
			// Handle ChainSideEvent
		case <-self.chainSideCh:

			// System stopped
		case <-self.txsSub.Err():
			quitByErr <- true
			return
		case <-self.chainHeadSub.Err():
			quitByErr <- true
			return
		case <-self.chainSideSub.Err():
			quitByErr <- true
			return
		}
	}
}

func (self *worker) wait(TxResendUseLegacy bool) {
	for {
		mustCommitNewWork := true
		for result := range self.recv {
			atomic.AddInt32(&self.atWork, -1)
			ResultChGauge.Update(ResultChGauge.Value() - 1)
			if result == nil {
				continue
			}

			// TODO-Klaytn drop or missing tx
			if self.nodetype != common.CONSENSUSNODE {
				if !TxResendUseLegacy {
					continue
				}
				pending, err := self.backend.TxPool().Pending()
				if err != nil {
					logger.Error("Failed to fetch pending transactions", "err", err)
					continue
				}

				if len(pending) > 0 {
					accounts := len(pending)
					resendTxSize := maxResendTxSize / accounts
					if resendTxSize == 0 {
						resendTxSize = 1
					}
					var resendTxs []*types.Transaction
					for _, sortedTxs := range pending {
						if len(sortedTxs) >= resendTxSize {
							resendTxs = append(resendTxs, sortedTxs[:resendTxSize]...)
						} else {
							resendTxs = append(resendTxs, sortedTxs...)
						}
					}
					if len(resendTxs) > 0 {
						resentTxGauge.Update(int64(len(resendTxs)))
						self.backend.ReBroadcastTxs(resendTxs)
					}
				}
				continue
			}

			block := result.Block
			work := result.Task

			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.
			for _, r := range work.receipts {
				for _, l := range r.Logs {
					l.BlockHash = block.Hash()
				}
			}
			work.stateMu.Lock()
			for _, log := range work.state.Logs() {
				log.BlockHash = block.Hash()
			}

			start := time.Now()
			result, err := self.chain.WriteBlockWithState(block, work.receipts, work.state)
			work.stateMu.Unlock()
			if err != nil {
				if err == blockchain.ErrKnownBlock {
					logger.Debug("Tried to insert already known block", "num", block.NumberU64(), "hash", block.Hash().String())
				} else {
					logger.Error("Failed writing block to chain", "err", err)
				}
				continue
			}
			blockWriteTime := time.Since(start)

			// TODO-Klaytn-Issue264 If we are using istanbul BFT, then we always have a canonical chain.
			//         Later we may be able to refine below code.

			// check if canon block and write transactions
			if result.Status == blockchain.CanonStatTy {
				// implicit by posting ChainHeadEvent
				mustCommitNewWork = false
			}

			// Broadcast the block and announce chain insertion event
			self.mux.Post(blockchain.NewMinedBlockEvent{Block: block})

			var events []interface{}

			work.stateMu.RLock()
			logs := work.state.Logs()
			work.stateMu.RUnlock()

			events = append(events, blockchain.ChainEvent{Block: block, Hash: block.Hash(), Logs: logs})
			if result.Status == blockchain.CanonStatTy {
				events = append(events, blockchain.ChainHeadEvent{Block: block})
			}

			// update governance CurrentSet if it is at an epoch block
			if err := self.engine.CreateSnapshot(self.chain, block.NumberU64(), block.Hash(), nil); err != nil {
				logger.Error("Failed to call snapshot", "err", err)
			}

			logger.Info("Successfully wrote mined block", "num", block.NumberU64(),
				"hash", block.Hash(), "txs", len(block.Transactions()), "elapsed", blockWriteTime)
			self.chain.PostChainEvents(events, logs)

			// TODO-Klaytn-Issue264 If we are using istanbul BFT, then we always have a canonical chain.
			//         Later we may be able to refine below code.
			if mustCommitNewWork {
				self.commitNewWork()
			}
		}
	}
}

// push sends a new work task to currently live work agents.
func (self *worker) push(work *Task) {
	if atomic.LoadInt32(&self.mining) != 1 {
		return
	}
	for agent := range self.agents {
		atomic.AddInt32(&self.atWork, 1)
		if ch := agent.Work(); ch != nil {
			ch <- work
		}
	}
}

// makeCurrent creates a new environment for the current cycle.
func (self *worker) makeCurrent(parent *types.Block, header *types.Header) error {
	stateDB, err := self.chain.StateAt(parent.Root())
	if err != nil {
		return err
	}
	work := NewTask(self.config, types.MakeSigner(self.config, header.Number), stateDB, header)
	if self.nodetype != common.CONSENSUSNODE {
		work.Block = parent
	}

	// Keep track of transactions which return errors so they can be removed
	work.tcount = 0
	self.current = work
	return nil
}

func (self *worker) commitNewWork() {
	var pending map[common.Address]types.Transactions
	var err error
	if self.nodetype == common.CONSENSUSNODE {
		// Check any fork transitions needed
		pending, err = self.backend.TxPool().Pending()
		if err != nil {
			logger.Error("Failed to fetch pending transactions", "err", err)
			return
		}
	}

	self.mu.Lock()
	defer self.mu.Unlock()
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	parent := self.chain.CurrentBlock()

	// TODO-Klaytn drop or missing tx
	tstart := time.Now()
	tstamp := tstart.Unix()
	if self.nodetype == common.CONSENSUSNODE {
		ideal := time.Unix(parent.Time().Int64()+params.BlockGenerationInterval, 0)
		// If a timestamp of this block is faster than the ideal timestamp,
		// wait for a while and get a new timestamp
		if tstart.Before(ideal) {
			wait := ideal.Sub(tstart)
			logger.Info("Mining too far in the future", "wait", common.PrettyDuration(wait))
			time.Sleep(wait)

			tstart = time.Now()    // refresh for metrics
			tstamp = tstart.Unix() // refresh for block timestamp
		}
	}

	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		Extra:      self.extra,
		Time:       big.NewInt(tstamp),
	}
	if err := self.engine.Prepare(self.chain, header); err != nil {
		logger.Error("Failed to prepare header for mining", "err", err)
		return
	}
	// Could potentially happen if starting to mine in an odd state.
	err = self.makeCurrent(parent, header)
	if err != nil {
		logger.Error("Failed to create mining context", "err", err)
		return
	}

	// Obtain current work's state lock after we receive new work assignment.
	self.current.stateMu.Lock()
	defer self.current.stateMu.Unlock()

	// Create the current work task
	work := self.current
	if self.nodetype == common.CONSENSUSNODE {
		txs := types.NewTransactionsByPriceAndNonce(self.current.signer, pending)
		work.commitTransactions(self.mux, txs, self.chain, self.rewardbase)
		finishedCommitTx := time.Now()

		// Create the new block to seal with the consensus engine
		if work.Block, err = self.engine.Finalize(self.chain, header, work.state, work.txs, work.receipts); err != nil {
			logger.Error("Failed to finalize block for sealing", "err", err)
			return
		}
		finishedFinalize := time.Now()

		// We only care about logging if we're actually mining.
		if atomic.LoadInt32(&self.mining) == 1 {
			// Update the metrics subsystem with all the measurements
			accountReadTimer.Update(work.state.AccountReads)
			accountHashTimer.Update(work.state.AccountHashes)
			accountUpdateTimer.Update(work.state.AccountUpdates)
			accountCommitTimer.Update(work.state.AccountCommits)

			storageReadTimer.Update(work.state.StorageReads)
			storageHashTimer.Update(work.state.StorageHashes)
			storageUpdateTimer.Update(work.state.StorageUpdates)
			storageCommitTimer.Update(work.state.StorageCommits)

			snapshotAccountReadTimer.Update(work.state.SnapshotAccountReads)
			snapshotStorageReadTimer.Update(work.state.SnapshotStorageReads)
			snapshotCommitTimer.Update(work.state.SnapshotCommits)

			trieAccess := work.state.AccountReads + work.state.AccountHashes + work.state.AccountUpdates + work.state.AccountCommits
			trieAccess += work.state.StorageReads + work.state.StorageHashes + work.state.StorageUpdates + work.state.StorageCommits

			tCountGauge.Update(int64(work.tcount))
			blockMiningTime := time.Since(tstart)
			commitTxTime := finishedCommitTx.Sub(tstart)
			finalizeTime := finishedFinalize.Sub(finishedCommitTx)

			blockMiningTimer.Update(blockMiningTime)
			blockMiningCommitTxTimer.Update(commitTxTime)
			blockMiningExecuteTxTimer.Update(commitTxTime - trieAccess)
			blockMiningFinalizeTimer.Update(finalizeTime)
			logger.Info("Commit new mining work",
				"number", work.Block.Number(), "hash", work.Block.Hash(),
				"txs", work.tcount, "elapsed", common.PrettyDuration(blockMiningTime),
				"commitTime", common.PrettyDuration(commitTxTime), "finalizeTime", common.PrettyDuration(finalizeTime))
		}
	}

	self.push(work)
	self.updateSnapshot()
}

func (self *worker) updateSnapshot() {
	self.snapshotMu.Lock()
	defer self.snapshotMu.Unlock()

	self.snapshotBlock = types.NewBlock(
		self.current.header,
		self.current.txs,
		self.current.receipts,
	)
	self.snapshotState = self.current.state.Copy()
}

func (env *Task) commitTransactions(mux *event.TypeMux, txs *types.TransactionsByPriceAndNonce, bc BlockChain, rewardbase common.Address) {
	coalescedLogs := env.ApplyTransactions(txs, bc, rewardbase)

	if len(coalescedLogs) > 0 || env.tcount > 0 {
		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		go func(logs []*types.Log, tcount int) {
			if len(logs) > 0 {
				mux.Post(blockchain.PendingLogsEvent{Logs: logs})
			}
			if tcount > 0 {
				mux.Post(blockchain.PendingStateEvent{})
			}
		}(cpy, env.tcount)
	}
}

func (env *Task) ApplyTransactions(txs *types.TransactionsByPriceAndNonce, bc BlockChain, rewardbase common.Address) []*types.Log {
	var coalescedLogs []*types.Log

	// Limit the execution time of all transactions in a block
	var abort int32 = 0       // To break the below commitTransaction for loop when timed out
	chDone := make(chan bool) // To stop the goroutine below when processing txs is completed

	// chEVM is used to notify the below goroutine of the running EVM so it can call evm.Cancel
	// when timed out.  We use a buffered channel to prevent the main EVM execution routine
	// from being blocked due to the channel communication.
	chEVM := make(chan *vm.EVM, 1)

	go func() {
		blockTimer := time.NewTimer(params.BlockGenerationTimeLimit)
		timeout := false
		var evm *vm.EVM

		for {
			select {
			case <-blockTimer.C:
				timeout = true
				atomic.StoreInt32(&abort, 1)

			case <-chDone:
				// Everything is done. Stop this goroutine.
				return

			case evm = <-chEVM:
			}

			if timeout && evm != nil {
				// Allow the first transaction to complete although it exceeds the time limit.
				if env.tcount > 0 {
					// The total time limit reached, thus we stop the currently running EVM.
					evm.Cancel(vm.CancelByTotalTimeLimit)
				}
				evm = nil
			}
		}
	}()

	vmConfig := &vm.Config{
		RunningEVM:               chEVM,
		UseOpcodeComputationCost: true,
	}

	var numTxsChecked int64 = 0
	var numTxsNonceTooLow int64 = 0
	var numTxsNonceTooHigh int64 = 0
	var numTxsGasLimitReached int64 = 0
CommitTransactionLoop:
	for atomic.LoadInt32(&abort) == 0 {
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			// To indicate that it does not have enough transactions for params.BlockGenerationTimeLimit.
			if numTxsChecked > 0 {
				usedAllTxsCounter.Inc(1)
			}
			break
		}
		numTxsChecked++
		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := types.Sender(env.signer, tx)

		// NOTE-Klaytn Since Klaytn is always in EIP155, the below replay protection code is not needed.
		// TODO-Klaytn-RemoveLater Remove the code commented below.
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		//if tx.Protected() && !env.config.IsEIP155(env.header.Number) {
		//	logger.Trace("Ignoring reply protected transaction", "hash", tx.Hash())
		//	//logger.Error("#### worker.commitTransaction","tx.protected",tx.Protected(),"tx.hash",tx.Hash(),"nonce",tx.Nonce(),"to",tx.To())
		//	txs.Pop()
		//	continue
		//}
		// Start executing the transaction
		env.state.Prepare(tx.Hash(), common.Hash{}, env.tcount)

		err, logs := env.commitTransaction(tx, bc, rewardbase, vmConfig)
		switch err {
		case blockchain.ErrGasLimitReached:
			// Pop the current out-of-gas transaction without shifting in the next from the account
			logger.Trace("Gas limit exceeded for current block", "sender", from)
			numTxsGasLimitReached++
			txs.Pop()

		case blockchain.ErrNonceTooLow:
			// New head notification data race between the transaction pool and miner, shift
			logger.Trace("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce())
			numTxsNonceTooLow++
			txs.Shift()

		case blockchain.ErrNonceTooHigh:
			// Reorg notification data race between the transaction pool and miner, skip account =
			logger.Trace("Skipping account with high nonce", "sender", from, "nonce", tx.Nonce())
			numTxsNonceTooHigh++
			txs.Pop()

		case vm.ErrTotalTimeLimitReached:
			logger.Warn("Transaction aborted due to time limit", "hash", tx.Hash().String())
			timeLimitReachedCounter.Inc(1)
			if env.tcount == 0 {
				logger.Error("A single transaction exceeds total time limit", "hash", tx.Hash().String())
				tooLongTxCounter.Inc(1)
			}
			// NOTE-Klaytn Exit for loop immediately without checking abort variable again.
			break CommitTransactionLoop

		case blockchain.ErrTxTypeNotSupported:
			// Pop the unsupported transaction without shifting in the next from the account
			logger.Trace("Skipping unsupported transaction type", "sender", from, "type", tx.Type())
			txs.Pop()

		case nil:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			env.tcount++
			txs.Shift()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			logger.Warn("Transaction failed, account skipped", "sender", from, "hash", tx.Hash().String(), "err", err)
			strangeErrorTxsCounter.Inc(1)
			txs.Shift()
		}
	}

	// Update the number of transactions checked and dropped during ApplyTransactions.
	checkedTxsGauge.Update(numTxsChecked)
	nonceTooLowTxsGauge.Update(numTxsNonceTooLow)
	nonceTooHighTxsGauge.Update(numTxsNonceTooHigh)
	gasLimitReachedTxsGauge.Update(numTxsGasLimitReached)

	// Stop the goroutine that has been handling the timer.
	chDone <- true

	return coalescedLogs
}

func (env *Task) commitTransaction(tx *types.Transaction, bc BlockChain, rewardbase common.Address, vmConfig *vm.Config) (error, []*types.Log) {
	snap := env.state.Snapshot()

	receipt, _, _, err := bc.ApplyTransaction(env.config, &rewardbase, env.state, env.header, tx, &env.header.GasUsed, vmConfig)
	if err != nil {
		if err != vm.ErrInsufficientBalance && err != vm.ErrTotalTimeLimitReached {
			tx.MarkUnexecutable(true)
		}
		env.state.RevertToSnapshot(snap)
		return err, nil
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)

	return nil, receipt.Logs
}

func NewTask(config *params.ChainConfig, signer types.Signer, statedb *state.StateDB, header *types.Header) *Task {
	return &Task{
		config:    config,
		signer:    signer,
		state:     statedb,
		header:    header,
		createdAt: time.Now(),
	}
}

func (env *Task) Transactions() []*types.Transaction { return env.txs }
func (env *Task) Receipts() []*types.Receipt         { return env.receipts }
