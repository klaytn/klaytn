// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/tx_pool.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/prque"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
	"github.com/rcrowley/go-metrics"
)

const (
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10
	// demoteUnexecutablesFullValidationTxLimit is the number of txs will be fully validated in demoteUnexecutables.
	demoteUnexecutablesFullValidationTxLimit = 1000
	// txMsgCh is the number of list of transactions can be queued.
	txMsgChSize = 100
	// MaxTxDataSize is a heuristic limit of tx data size, and txPool rejects transactions over 32KB to prevent DOS attacks.
	MaxTxDataSize = 32 * 1024
)

var (
	evictionInterval    = time.Minute     // Time interval to check for evictable transactions
	statsReportInterval = 8 * time.Second // Time interval to report transaction pool stats

	txPoolIsFullErr = fmt.Errorf("txpool is full")

	errNotAllowedAnchoringTx = errors.New("locally anchoring chaindata tx is not allowed in this node")
)

var (
	// Metrics for the pending pool
	pendingDiscardCounter   = metrics.NewRegisteredCounter("txpool/pending/discard", nil)
	pendingReplaceCounter   = metrics.NewRegisteredCounter("txpool/pending/replace", nil)
	pendingRateLimitCounter = metrics.NewRegisteredCounter("txpool/pending/ratelimit", nil) // Dropped due to rate limiting
	pendingNofundsCounter   = metrics.NewRegisteredCounter("txpool/pending/nofunds", nil)   // Dropped due to out-of-funds

	// Metrics for the queued pool
	queuedDiscardCounter   = metrics.NewRegisteredCounter("txpool/queued/discard", nil)
	queuedReplaceCounter   = metrics.NewRegisteredCounter("txpool/queued/replace", nil)
	queuedRateLimitCounter = metrics.NewRegisteredCounter("txpool/queued/ratelimit", nil) // Dropped due to rate limiting
	queuedNofundsCounter   = metrics.NewRegisteredCounter("txpool/queued/nofunds", nil)   // Dropped due to out-of-funds

	// General tx metrics
	invalidTxCounter     = metrics.NewRegisteredCounter("txpool/invalid", nil)
	underpricedTxCounter = metrics.NewRegisteredCounter("txpool/underpriced", nil)
	refusedTxCounter     = metrics.NewRegisteredCounter("txpool/refuse", nil)
)

// TxStatus is the current status of a transaction as seen by the pool.
type TxStatus uint

const (
	TxStatusUnknown TxStatus = iota
	TxStatusQueued
	TxStatusPending
	// for Les
	TxStatusIncluded
)

// blockChain provides the state of blockchain and current gas limit to do
// some pre checks in tx pool and event subscribers.
type blockChain interface {
	CurrentBlock() *types.Block
	GetBlock(hash common.Hash, number uint64) *types.Block
	StateAt(root common.Hash) (*state.StateDB, error)

	SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription
}

// TxPoolConfig are the configuration parameters of the transaction pool.
type TxPoolConfig struct {
	NoLocals           bool          // Whether local transaction handling should be disabled
	AllowLocalAnchorTx bool          // if this is true, the txpool allow locally submitted anchor transactions
	DenyRemoteTx       bool          // Denies remote transactions receiving from other peers
	Journal            string        // Journal of local transactions to survive node restarts
	JournalInterval    time.Duration // Time interval to regenerate the local transaction journal

	PriceLimit uint64 // Minimum gas price to enforce for acceptance into the pool
	PriceBump  uint64 // Minimum price bump percentage to replace an already existing transaction (nonce)

	ExecSlotsAccount    uint64 // Number of executable transaction slots guaranteed per account
	ExecSlotsAll        uint64 // Maximum number of executable transaction slots for all accounts
	NonExecSlotsAccount uint64 // Maximum number of non-executable transaction slots permitted per account
	NonExecSlotsAll     uint64 // Maximum number of non-executable transaction slots for all accounts

	KeepLocals bool          // Disables removing timed-out local transactions
	Lifetime   time.Duration // Maximum amount of time non-executable transaction are queued

	NoAccountCreation            bool // Whether account creation transactions should be disabled
	EnableSpamThrottlerAtRuntime bool // Enable txpool spam throttler at runtime
}

// DefaultTxPoolConfig contains the default configurations for the transaction
// pool.
var DefaultTxPoolConfig = TxPoolConfig{
	Journal:         "transactions.rlp",
	JournalInterval: time.Hour,

	PriceLimit: 1,
	PriceBump:  10,

	ExecSlotsAccount:    16,
	ExecSlotsAll:        4096,
	NonExecSlotsAccount: 64,
	NonExecSlotsAll:     1024,

	KeepLocals: false,
	Lifetime:   5 * time.Minute,
}

// sanitize checks the provided user configurations and changes anything that's
// unreasonable or unworkable.
func (config *TxPoolConfig) sanitize() TxPoolConfig {
	conf := *config
	if conf.JournalInterval < time.Second {
		logger.Error("Sanitizing invalid txpool journal time", "provided", conf.JournalInterval, "updated", time.Second)
		conf.JournalInterval = time.Second
	}
	if conf.PriceLimit < 1 {
		logger.Error("Sanitizing invalid txpool price limit", "provided", conf.PriceLimit, "updated", DefaultTxPoolConfig.PriceLimit)
		conf.PriceLimit = DefaultTxPoolConfig.PriceLimit
	}
	if conf.PriceBump < 1 {
		logger.Error("Sanitizing invalid txpool price bump", "provided", conf.PriceBump, "updated", DefaultTxPoolConfig.PriceBump)
		conf.PriceBump = DefaultTxPoolConfig.PriceBump
	}
	return conf
}

// TxPool contains all currently known transactions. Transactions
// enter the pool when they are received from the network or submitted
// locally. They exit the pool when they are included in the blockchain.
//
// The pool separates processable transactions (which can be applied to the
// current state) and future transactions. Transactions move between those
// two states over time as they are received and processed.
type TxPool struct {
	config       TxPoolConfig
	chainconfig  *params.ChainConfig
	chain        blockChain
	gasPrice     *big.Int
	txFeed       event.Feed
	scope        event.SubscriptionScope
	chainHeadCh  chan ChainHeadEvent
	chainHeadSub event.Subscription
	signer       types.Signer
	mu           sync.RWMutex

	currentBlockNumber uint64                    // Current block number
	currentState       *state.StateDB            // Current state in the blockchain head
	pendingNonce       map[common.Address]uint64 // Pending nonce tracking virtual nonces

	locals  *accountSet // Set of local transaction to exempt from eviction rules
	journal *txJournal  // Journal of local transaction to back up to disk

	//TODO-Klaytn
	txMu sync.RWMutex

	pending map[common.Address]*txList         // All currently processable transactions
	queue   map[common.Address]*txList         // Queued but non-processable transactions
	beats   map[common.Address]time.Time       // Last heartbeat from each known account
	all     map[common.Hash]*types.Transaction // All transactions to allow lookups
	priced  *txPricedList                      // All transactions sorted by price

	wg sync.WaitGroup // for shutdown sync

	txMsgCh chan types.Transactions

	eip2718 bool // Fork indicator whether we are using EIP-2718 type transactions.
	eip1559 bool // Fork indicator whether we are using EIP-1559 type transactions.
}

// NewTxPool creates a new transaction pool to gather, sort and filter inbound
// transactions from the network.
func NewTxPool(config TxPoolConfig, chainconfig *params.ChainConfig, chain blockChain) *TxPool {
	// Sanitize the input to ensure no vulnerable gas prices are set
	config = (&config).sanitize()

	// Create the transaction pool with its initial settings
	pool := &TxPool{
		config:       config,
		chainconfig:  chainconfig,
		chain:        chain,
		signer:       types.LatestSignerForChainID(chainconfig.ChainID),
		pending:      make(map[common.Address]*txList),
		queue:        make(map[common.Address]*txList),
		beats:        make(map[common.Address]time.Time),
		all:          make(map[common.Hash]*types.Transaction),
		pendingNonce: make(map[common.Address]uint64),
		chainHeadCh:  make(chan ChainHeadEvent, chainHeadChanSize),
		gasPrice:     new(big.Int).SetUint64(chainconfig.UnitPrice),
		txMsgCh:      make(chan types.Transactions, txMsgChSize),
	}
	pool.locals = newAccountSet(pool.signer)
	pool.priced = newTxPricedList(&pool.all)
	pool.reset(nil, chain.CurrentBlock().Header())

	// If local transactions and journaling is enabled, load from disk
	if !config.NoLocals && config.Journal != "" {
		pool.journal = newTxJournal(config.Journal)

		if err := pool.journal.load(pool.AddLocals); err != nil {
			logger.Error("Failed to load transaction journal", "err", err)
		}
		if err := pool.journal.rotate(pool.local(), pool.signer); err != nil {
			logger.Error("Failed to rotate transaction journal", "err", err)
		}
	}
	// Subscribe events from blockchain
	pool.chainHeadSub = pool.chain.SubscribeChainHeadEvent(pool.chainHeadCh)

	// Start the event loop and return
	pool.wg.Add(2)
	go pool.loop()
	go pool.handleTxMsg()

	if config.EnableSpamThrottlerAtRuntime {
		if err := pool.StartSpamThrottler(DefaultSpamThrottlerConfig); err != nil {
			logger.Error("Failed to start spam throttler", "err", err)
		}
	}

	return pool
}

// loop is the transaction pool's main event loop, waiting for and reacting to
// outside blockchain events as well as for various reporting and transaction
// eviction events.
func (pool *TxPool) loop() {
	defer pool.wg.Done()

	// Start the stats reporting and transaction eviction tickers
	var prevPending, prevQueued, prevStales int

	report := time.NewTicker(statsReportInterval)
	defer report.Stop()

	evict := time.NewTicker(evictionInterval)
	defer evict.Stop()

	journal := time.NewTicker(pool.config.JournalInterval)
	defer journal.Stop()

	// Track the previous head headers for transaction reorgs
	head := pool.chain.CurrentBlock()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle ChainHeadEvent
		case ev := <-pool.chainHeadCh:
			if ev.Block != nil {
				pool.mu.Lock()
				currBlock := pool.chain.CurrentBlock()
				if ev.Block.Root() != currBlock.Root() {
					pool.mu.Unlock()
					logger.Debug("block from ChainHeadEvent is different from the CurrentBlock",
						"receivedNum", ev.Block.NumberU64(), "receivedHash", ev.Block.Hash().String(),
						"currNum", currBlock.NumberU64(), "currHash", currBlock.Hash().String())
					continue
				}
				pool.reset(head.Header(), ev.Block.Header())
				head = ev.Block
				pool.mu.Unlock()
			}
			// Be unsubscribed due to system stopped
		case <-pool.chainHeadSub.Err():
			return

			// Handle stats reporting ticks
		case <-report.C:
			pool.mu.RLock()
			pending, queued := pool.stats()
			stales := pool.priced.stales
			pool.mu.RUnlock()

			if pending != prevPending || queued != prevQueued || stales != prevStales {
				logger.Debug("Transaction pool status report", "executable", pending, "queued", queued, "stales", stales)
				prevPending, prevQueued, prevStales = pending, queued, stales
				txPoolPendingGauge.Update(int64(pending))
				txPoolQueueGauge.Update(int64(queued))
			}

			// Handle inactive account transaction eviction
		case <-evict.C:
			pool.mu.Lock()
			for addr, beat := range pool.beats {
				// Skip local transactions from the eviction mechanism
				if pool.config.KeepLocals && pool.locals.contains(addr) {
					delete(pool.beats, addr)
					continue
				}

				// Any non-locals old enough should be removed
				if time.Since(beat) > pool.config.Lifetime {
					if pool.queue[addr] != nil {
						for _, tx := range pool.queue[addr].Flatten() {
							pool.removeTx(tx.Hash(), true)
						}
					}
					delete(pool.beats, addr)
				}
			}
			pool.mu.Unlock()

			// Handle local transaction journal rotation
		case <-journal.C:
			if pool.journal != nil {
				pool.mu.Lock()
				if err := pool.journal.rotate(pool.local(), pool.signer); err != nil {
					logger.Error("Failed to rotate local tx journal", "err", err)
				}
				pool.mu.Unlock()
			}
		}
	}
}

// lockedReset is a wrapper around reset to allow calling it in a thread safe
// manner. This method is only ever used in the tester!
func (pool *TxPool) lockedReset(oldHead, newHead *types.Header) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.reset(oldHead, newHead)
}

// reset retrieves the current state of the blockchain and ensures the content
// of the transaction pool is valid with regard to the chain state.
func (pool *TxPool) reset(oldHead, newHead *types.Header) {
	// If we're reorging an old state, reinject all dropped transactions
	var reinject types.Transactions

	if oldHead != nil && oldHead.Hash() != newHead.ParentHash {
		// If the reorg is too deep, avoid doing it (will happen during fast sync)
		oldNum := oldHead.Number.Uint64()
		newNum := newHead.Number.Uint64()

		if depth := uint64(math.Abs(float64(oldNum) - float64(newNum))); depth > 64 {
			logger.Debug("Skipping deep transaction reorg", "depth", depth)
		} else {
			// Reorg seems shallow enough to pull in all transactions into memory
			var discarded, included types.Transactions

			var (
				rem = pool.chain.GetBlock(oldHead.Hash(), oldHead.Number.Uint64())
				add = pool.chain.GetBlock(newHead.Hash(), newHead.Number.Uint64())
			)
			for rem.NumberU64() > add.NumberU64() {
				discarded = append(discarded, rem.Transactions()...)
				if rem = pool.chain.GetBlock(rem.ParentHash(), rem.NumberU64()-1); rem == nil {
					logger.Error("Unrooted old chain seen by tx pool", "block", oldHead.Number, "hash", oldHead.Hash())
					return
				}
			}
			for add.NumberU64() > rem.NumberU64() {
				included = append(included, add.Transactions()...)
				if add = pool.chain.GetBlock(add.ParentHash(), add.NumberU64()-1); add == nil {
					logger.Error("Unrooted new chain seen by tx pool", "block", newHead.Number, "hash", newHead.Hash())
					return
				}
			}
			for rem.Hash() != add.Hash() {
				discarded = append(discarded, rem.Transactions()...)
				if rem = pool.chain.GetBlock(rem.ParentHash(), rem.NumberU64()-1); rem == nil {
					logger.Error("Unrooted old chain seen by tx pool", "block", oldHead.Number, "hash", oldHead.Hash())
					return
				}
				included = append(included, add.Transactions()...)
				if add = pool.chain.GetBlock(add.ParentHash(), add.NumberU64()-1); add == nil {
					logger.Error("Unrooted new chain seen by tx pool", "block", newHead.Number, "hash", newHead.Hash())
					return
				}
			}
			reinject = types.TxDifference(discarded, included)
		}
	}
	// Initialize the internal state to the current head
	if newHead == nil {
		newHead = pool.chain.CurrentBlock().Header() // Special case during testing
	}
	stateDB, err := pool.chain.StateAt(newHead.Root)
	if err != nil {
		logger.Error("Failed to reset txpool state", "err", err)
		return
	}
	pool.currentState = stateDB
	pool.pendingNonce = make(map[common.Address]uint64)
	pool.currentBlockNumber = newHead.Number.Uint64()

	// Inject any transactions discarded due to reorgs
	logger.Debug("Reinjecting stale transactions", "count", len(reinject))
	senderCacher.recover(pool.signer, reinject)

	//pool.mu.Lock()
	//defer pool.mu.Unlock()

	pool.addTxsLocked(reinject, false)

	// validate the pool of pending transactions, this will remove
	// any transactions that have been included in the block or
	// have been invalidated because of another transaction (e.g.
	// higher gas price)
	pool.demoteUnexecutables()

	pool.txMu.Lock()
	// Update all accounts to the latest known pending nonce
	for addr, list := range pool.pending {
		txs := list.Flatten()
		if len(txs) > 0 {
			// Heavy but will be cached and is needed by the miner anyway
			pool.setPendingNonce(addr, txs[len(txs)-1].Nonce()+1)
		}
	}
	pool.txMu.Unlock()
	// Check the queue and move transactions over to the pending if possible
	// or remove those that have become invalid
	pool.promoteExecutables(nil)

	// Update all fork indicator by next pending block number.
	next := new(big.Int).Add(newHead.Number, big.NewInt(1))

	// Enable Ethereum tx type transactions
	pool.eip2718 = pool.chainconfig.IsEthTxTypeForkEnabled(next)
	pool.eip1559 = pool.chainconfig.IsEthTxTypeForkEnabled(next)
}

// Stop terminates the transaction pool.
func (pool *TxPool) Stop() {
	// Unsubscribe all subscriptions registered from txpool
	pool.scope.Close()

	// Unsubscribe subscriptions registered from blockchain
	pool.chainHeadSub.Unsubscribe()
	pool.wg.Wait()

	if pool.journal != nil {
		pool.journal.close()
	}

	pool.StopSpamThrottler()
	logger.Info("Transaction pool stopped")
}

// SubscribeNewTxsEvent registers a subscription of NewTxsEvent and
// starts sending event to the given channel.
func (pool *TxPool) SubscribeNewTxsEvent(ch chan<- NewTxsEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}

// GasPrice returns the current gas price enforced by the transaction pool.
func (pool *TxPool) GasPrice() *big.Int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return new(big.Int).Set(pool.gasPrice)
}

// SetGasPrice updates the gas price of the transaction pool for new transactions, and drops all old transactions.
func (pool *TxPool) SetGasPrice(price *big.Int) {
	if pool.gasPrice.Cmp(price) != 0 {
		pool.mu.Lock()

		logger.Info("TxPool.SetGasPrice", "before", pool.gasPrice, "after", price)

		pool.gasPrice = price
		pool.pending = make(map[common.Address]*txList)
		pool.queue = make(map[common.Address]*txList)
		pool.beats = make(map[common.Address]time.Time)
		pool.all = make(map[common.Hash]*types.Transaction)
		pool.pendingNonce = make(map[common.Address]uint64)
		pool.locals = newAccountSet(pool.signer)
		pool.priced = newTxPricedList(&pool.all)

		pool.mu.Unlock()
	}
}

// Stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) Stats() (int, int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.stats()
}

// stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) stats() (int, int) {
	pending := 0
	for _, list := range pool.pending {
		pending += list.Len()
	}
	queued := 0
	for _, list := range pool.queue {
		queued += list.Len()
	}
	return pending, queued
}

// Content retrieves the data content of the transaction pool, returning all the
// pending as well as queued transactions, grouped by account and sorted by nonce.
func (pool *TxPool) Content() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.txMu.Lock()
	defer pool.txMu.Unlock()

	pending := make(map[common.Address]types.Transactions)
	for addr, list := range pool.pending {
		pending[addr] = list.Flatten()
	}
	queued := make(map[common.Address]types.Transactions)
	for addr, list := range pool.queue {
		queued[addr] = list.Flatten()
	}
	return pending, queued
}

// Pending retrieves all currently processable transactions, groupped by origin
// account and sorted by nonce. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) Pending() (map[common.Address]types.Transactions, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.txMu.Lock()
	defer pool.txMu.Unlock()

	pending := make(map[common.Address]types.Transactions)
	for addr, list := range pool.pending {
		pending[addr] = list.Flatten()
	}
	return pending, nil
}

// CachedPendingTxByCount retrieves about number of currently processable transactions
// by requested count, grouped by origin account and sorted by nonce.
func (pool *TxPool) CachedPendingTxsByCount(count int) types.Transactions {
	if count <= 0 {
		return nil
	}

	// It retrieves the half of the requested transaction recursively for returned
	// transactions much as possible.
	txPerAddr := count / 2
	if txPerAddr == 0 {
		txPerAddr = 1
	}

	pending := make(types.Transactions, 0, count)

	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.txMu.Lock()
	defer pool.txMu.Unlock()

	if len(pool.pending) == 0 {
		return nil
	}

	for _, list := range pool.pending {
		pendingPerAccount := list.CachedTxsFlattenByCount(txPerAddr)

		pending = append(pending, pendingPerAccount...)
		if len(pending) >= count {
			break
		}

		if len(pendingPerAccount) >= txPerAddr {
			if txPerAddr > 1 {
				txPerAddr = txPerAddr / 2
			}
		}
	}
	return pending
}

// local retrieves all currently known local transactions, groupped by origin
// account and sorted by nonce. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) local() map[common.Address]types.Transactions {
	txs := make(map[common.Address]types.Transactions)
	for addr := range pool.locals.accounts {
		if pending := pool.pending[addr]; pending != nil {
			txs[addr] = append(txs[addr], pending.Flatten()...)
		}
		if queued := pool.queue[addr]; queued != nil {
			txs[addr] = append(txs[addr], queued.Flatten()...)
		}
	}
	return txs
}

// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) validateTx(tx *types.Transaction) error {
	// Accept only legacy transactions until EIP-2718/2930 activates.
	if !pool.eip2718 && tx.IsEthTypedTransaction() {
		return ErrTxTypeNotSupported
	}
	// Reject dynamic fee transactions until EIP-1559 activates.
	if !pool.eip1559 && tx.Type() == types.TxTypeEthereumDynamicFee {
		return ErrTxTypeNotSupported
	}

	gasFeePayer := uint64(0)

	// Check chain Id first.
	if tx.ChainId().Cmp(pool.chainconfig.ChainID) != 0 {
		return ErrInvalidChainId
	}

	// NOTE-Klaytn Drop transactions with unexpected gasPrice
	// If the transaction type is DynamicFee tx, Compare transaction's GasFeeCap(MaxFeePerGas) and GasTipCap with tx pool's gasPrice to check to have same value.
	if tx.Type() == types.TxTypeEthereumDynamicFee {
		// Sanity check for extremely large numbers
		if tx.GasTipCap().BitLen() > 256 {
			return ErrTipVeryHigh
		}

		if tx.GasFeeCap().BitLen() > 256 {
			return ErrFeeCapVeryHigh
		}

		// Ensure gasFeeCap is greater than or equal to gasTipCap.
		if tx.GasFeeCap().Cmp(tx.GasTipCap()) < 0 {
			return ErrTipAboveFeeCap
		}

		if pool.gasPrice.Cmp(tx.GasTipCap()) != 0 {
			logger.Trace("fail to validate maxPriorityFeePerGas", "unitprice", pool.gasPrice, "maxPriorityFeePerGas", tx.GasFeeCap())
			return ErrInvalidGasTipCap
		}

		if pool.gasPrice.Cmp(tx.GasFeeCap()) != 0 {
			logger.Trace("fail to validate maxFeePerGas", "unitprice", pool.gasPrice, "maxFeePerGas", tx.GasTipCap())
			return ErrInvalidGasFeeCap
		}
	} else {
		if pool.gasPrice.Cmp(tx.GasPrice()) != 0 {
			logger.Trace("fail to validate unitprice", "unitprice", pool.gasPrice, "txUnitPrice", tx.GasPrice())
			return ErrInvalidUnitPrice
		}
	}

	// Heuristic limit, reject transactions over 32KB to prevent DOS attacks
	if tx.Size() > MaxTxDataSize {
		return ErrOversizedData
	}

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return ErrNegativeValue
	}

	// Make sure the transaction is signed properly
	gasFrom, err := tx.ValidateSender(pool.signer, pool.currentState, pool.currentBlockNumber)
	if err != nil {
		return err
	}
	from := tx.ValidatedSender()

	// Ensure the transaction adheres to nonce ordering
	if pool.getNonce(from) > tx.Nonce() {
		return ErrNonceTooLow
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	senderBalance := pool.getBalance(from)
	if tx.IsFeeDelegatedTransaction() {
		// balance check for fee-delegated tx
		gasFeePayer, err = tx.ValidateFeePayer(pool.signer, pool.currentState, pool.currentBlockNumber)
		if err != nil {
			return ErrInvalidFeePayer
		}
		feePayer := tx.ValidatedFeePayer()
		feePayerBalance := pool.getBalance(feePayer)
		feeRatio, isRatioTx := tx.FeeRatio()
		if isRatioTx {
			// Check fee ratio range
			if !feeRatio.IsValid() {
				return kerrors.ErrFeeRatioOutOfRange
			}

			feeByFeePayer, feeBySender := types.CalcFeeWithRatio(feeRatio, tx.Fee())

			if senderBalance.Cmp(new(big.Int).Add(tx.Value(), feeBySender)) < 0 {
				logger.Trace("[tx_pool] insufficient funds for feeBySender", "from", from, "balance", senderBalance, "feeBySender", feeBySender)
				return ErrInsufficientFundsFrom
			}

			if feePayerBalance.Cmp(feeByFeePayer) < 0 {
				logger.Trace("[tx_pool] insufficient funds for feeByFeePayer", "feePayer", feePayer, "balance", feePayerBalance, "feeByFeePayer", feeByFeePayer)
				return ErrInsufficientFundsFeePayer
			}
		} else {
			if senderBalance.Cmp(tx.Value()) < 0 {
				logger.Trace("[tx_pool] insufficient funds for cost(value)", "from", from, "balance", senderBalance, "value", tx.Value())
				return ErrInsufficientFundsFrom
			}

			if feePayerBalance.Cmp(tx.Fee()) < 0 {
				logger.Trace("[tx_pool] insufficient funds for cost(gas * price)", "feePayer", feePayer, "balance", feePayerBalance, "fee", tx.Fee())
				return ErrInsufficientFundsFeePayer
			}
		}
	} else {
		// balance check for non-fee-delegated tx
		if senderBalance.Cmp(tx.Cost()) < 0 {
			logger.Trace("[tx_pool] insufficient funds for cost(gas * price + value)", "from", from, "balance", senderBalance, "cost", tx.Cost())
			return ErrInsufficientFundsFrom
		}
	}

	intrGas, err := tx.IntrinsicGas(pool.currentBlockNumber)
	intrGas += gasFrom + gasFeePayer
	if err != nil {
		return err
	}
	if tx.Gas() < intrGas {
		return ErrIntrinsicGas
	}

	// "tx.Validate()" conducts additional validation for each new txType.
	// Validate humanReadable address when this tx has "true" in the humanReadable field.
	// Validate accountKey when the this create or update an account
	// Validate the existence of the address which will be created though this Tx
	// Validate a contract account whether it is executable
	if err := tx.Validate(pool.currentState, pool.currentBlockNumber); err != nil {
		return err
	}

	return nil
}

// getMaxTxFromQueueWhenNonceIsMissing finds and returns a trasaction with max nonce in queue when a given Tx has missing nonce.
// Otherwise it returns a given Tx itself.
func (pool *TxPool) getMaxTxFromQueueWhenNonceIsMissing(tx *types.Transaction, from *common.Address) *types.Transaction {
	txs := pool.queue[*from].txs

	maxTx := tx
	if txs.Get(tx.Nonce()) != nil {
		return maxTx
	}

	for _, t := range txs.items {
		if maxTx.Nonce() < t.Nonce() {
			maxTx = t
		}
	}
	return maxTx
}

// add validates a transaction and inserts it into the non-executable queue for
// later pending promotion and execution. If the transaction is a replacement for
// an already pending or queued one, it overwrites the previous and returns this
// so outer code doesn't uselessly call promote.
//
// If a newly added transaction is marked as local, its sending account will be
// whitelisted, preventing any associated transaction from being dropped out of
// the pool due to pricing constraints.
func (pool *TxPool) add(tx *types.Transaction, local bool) (bool, error) {
	// If the transaction is already known, discard it
	hash := tx.Hash()
	if pool.all[hash] != nil {
		logger.Trace("Discarding already known transaction", "hash", hash)
		return false, fmt.Errorf("known transaction: %x", hash)
	}
	// If the transaction fails basic validation, discard it
	if err := pool.validateTx(tx); err != nil {
		logger.Trace("Discarding invalid transaction", "hash", hash, "err", err)
		invalidTxCounter.Inc(1)
		return false, err
	}

	// If the transaction pool is full and new Tx is valid,
	// (1) discard a new Tx if there is no room for the account of the Tx
	// (2) remove an old Tx with the largest nonce from queue to make a room for a new Tx with missing nonce
	// (3) discard a new Tx if the new Tx does not have a missing nonce
	// (4) discard underpriced transactions
	if uint64(len(pool.all)) >= pool.config.ExecSlotsAll+pool.config.NonExecSlotsAll {
		// (1) discard a new Tx if there is no room for the account of the Tx
		from, _ := types.Sender(pool.signer, tx)
		if pool.queue[from] == nil {
			logger.Trace("Rejecting a new Tx, because TxPool is full and there is no room for the account", "hash", tx.Hash(), "account", from)
			refusedTxCounter.Inc(1)
			return false, fmt.Errorf("txpool is full: %d", uint64(len(pool.all)))
		}

		maxTx := pool.getMaxTxFromQueueWhenNonceIsMissing(tx, &from)
		if maxTx != tx {
			// (2) remove an old Tx with the largest nonce from queue to make a room for a new Tx with missing nonce
			pool.removeTx(maxTx.Hash(), true)
			logger.Trace("Removing an old Tx with the max nonce to insert a new Tx with missing nonce, because TxPool is full", "account", from, "new nonce(previously missing)", tx.Nonce(), "removed max nonce", maxTx.Nonce())
		} else {
			// (3) discard a new Tx if the new Tx does not have a missing nonce
			logger.Trace("Rejecting a new Tx, because TxPool is full and a new TX does not have missing nonce", "hash", tx.Hash())
			refusedTxCounter.Inc(1)
			return false, fmt.Errorf("txpool is full and the new tx does not have missing nonce: %d", uint64(len(pool.all)))
		}

		// (4) discard underpriced transactions
		// If the new transaction is underpriced, don't accept it
		if !local && pool.priced.Underpriced(tx, pool.locals) {
			logger.Trace("Discarding underpriced transaction", "hash", hash, "price", tx.GasPrice())
			underpricedTxCounter.Inc(1)
			return false, ErrUnderpriced
		}
		// New transaction is better than our worse ones, make room for it
		drop := pool.priced.Discard(len(pool.all)-int(pool.config.ExecSlotsAll+pool.config.NonExecSlotsAll-1), pool.locals)
		for _, tx := range drop {
			logger.Trace("Discarding freshly underpriced transaction", "hash", tx.Hash(), "price", tx.GasPrice())
			underpricedTxCounter.Inc(1)
			pool.removeTx(tx.Hash(), false)
		}
	}
	// If the transaction is replacing an already pending one, do directly
	from, _ := types.Sender(pool.signer, tx) // already validated
	if list := pool.pending[from]; list != nil && list.Overlaps(tx) {
		// Nonce already pending, check if required price bump is met
		inserted, old := list.Add(tx, pool.config.PriceBump)
		if !inserted {
			pendingDiscardCounter.Inc(1)
			return false, ErrAlreadyNonceExistInPool
		}
		// New transaction is better, replace old one
		if old != nil {
			delete(pool.all, old.Hash())
			pool.priced.Removed()
			pendingReplaceCounter.Inc(1)
		}
		pool.all[tx.Hash()] = tx
		pool.priced.Put(tx)
		pool.journalTx(from, tx)

		logger.Trace("Pooled new executable transaction", "hash", hash, "from", from, "to", tx.To())

		// We've directly injected a replacement transaction, notify subsystems
		go pool.txFeed.Send(NewTxsEvent{types.Transactions{tx}})

		return old != nil, nil
	}
	// New transaction isn't replacing a pending one, push into queue
	replace, err := pool.enqueueTx(hash, tx)
	if err != nil {
		return false, err
	}
	// Mark local addresses and journal local transactions
	if local {
		pool.locals.add(from)
	}
	pool.journalTx(from, tx)

	logger.Trace("Pooled new future transaction", "hash", hash, "from", from, "to", tx.To())
	return replace, nil
}

// enqueueTx inserts a new transaction into the non-executable transaction queue.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) enqueueTx(hash common.Hash, tx *types.Transaction) (bool, error) {
	// Try to insert the transaction into the future queue
	from, _ := types.Sender(pool.signer, tx) // already validated
	if pool.queue[from] == nil {
		pool.queue[from] = newTxList(false)
	}
	inserted, old := pool.queue[from].Add(tx, pool.config.PriceBump)
	if !inserted {
		// An older transaction was better, discard this
		queuedDiscardCounter.Inc(1)
		return false, ErrAlreadyNonceExistInPool
	}
	// Discard any previous transaction and mark this
	if old != nil {
		delete(pool.all, old.Hash())
		pool.priced.Removed()
		queuedReplaceCounter.Inc(1)
	}
	if pool.all[hash] == nil {
		pool.all[hash] = tx
		pool.priced.Put(tx)
	}

	pool.checkAndSetBeat(from)
	return old != nil, nil
}

// journalTx adds the specified transaction to the local disk journal if it is
// deemed to have been sent from a local account.
func (pool *TxPool) journalTx(from common.Address, tx *types.Transaction) {
	// Only journal if it's enabled and the transaction is local
	if pool.journal == nil || !pool.locals.contains(from) {
		return
	}
	if err := pool.journal.insert(tx); err != nil {
		logger.Error("Failed to journal local transaction", "err", err)
	}
}

// promoteTx adds a transaction to the pending (processable) list of transactions
// and returns whether it was inserted or an older was better.
//
// Note, this method assumes the pool lock is held!
func (pool *TxPool) promoteTx(addr common.Address, hash common.Hash, tx *types.Transaction) bool {
	// Try to insert the transaction into the pending queue
	if pool.pending[addr] == nil {
		pool.pending[addr] = newTxList(true)
	}
	list := pool.pending[addr]

	inserted, old := list.Add(tx, pool.config.PriceBump)
	if !inserted {
		// An older transaction was better, discard this
		delete(pool.all, hash)
		pool.priced.Removed()

		pendingDiscardCounter.Inc(1)
		return false
	}
	// Otherwise discard any previous transaction and mark this
	if old != nil {
		delete(pool.all, old.Hash())
		pool.priced.Removed()

		pendingReplaceCounter.Inc(1)
	}
	// Failsafe to work around direct pending inserts (tests)
	if pool.all[hash] == nil {
		pool.all[hash] = tx
		pool.priced.Put(tx)
	}
	// Set the potentially new pending nonce and notify any subsystems of the new tx
	pool.beats[addr] = time.Now()
	pool.setPendingNonce(addr, tx.Nonce()+1)

	return true
}

// HandleTxMsg transfers transactions to a channel where handleTxMsg calls AddRemotes
// to handle them. This is made not to wait from the results from TxPool.AddRemotes.
func (pool *TxPool) HandleTxMsg(txs types.Transactions) {
	if pool.config.DenyRemoteTx {
		return
	}

	// Filter spam txs based on to-address of failed txs
	spamThrottler := GetSpamThrottler()
	if spamThrottler != nil {
		pool.mu.RLock()
		poolSize := uint64(len(pool.all))
		pool.mu.RUnlock()

		// Activate spam throttler when pool has enough txs
		if poolSize > uint64(spamThrottler.config.ActivateTxPoolSize) {
			allowTxs, throttleTxs := spamThrottler.classifyTxs(txs)

			for _, tx := range throttleTxs {
				select {
				case spamThrottler.throttleCh <- tx:
				default:
					logger.Trace("drop a tx when throttleTxs channel is full", "txHash", tx.Hash())
					throttlerDropCount.Inc(1)
				}
			}

			txs = allowTxs
		}
	}

	// TODO-Klaytn: Consider removing the next line and move the above logic to `addTx` or `AddRemotes`
	senderCacher.recover(pool.signer, txs)
	pool.txMsgCh <- txs
}

func (pool *TxPool) throttleLoop(spamThrottler *throttler) {
	ticker := time.Tick(time.Second)
	throttleNum := int(spamThrottler.config.ThrottleTPS)

	for {
		select {
		case <-spamThrottler.quitCh:
			logger.Info("Stop spam throttler loop")
			return

		case <-ticker:
			txs := types.Transactions{}

			iterNum := len(spamThrottler.throttleCh)
			if iterNum > throttleNum {
				iterNum = throttleNum
			}

			for i := 0; i < iterNum; i++ {
				tx := <-spamThrottler.throttleCh
				txs = append(txs, tx)
			}

			if len(txs) > 0 {
				pool.AddRemotes(txs)
			}
		}
	}
}

func (pool *TxPool) StartSpamThrottler(conf *ThrottlerConfig) error {
	spamThrottlerMu.Lock()
	defer spamThrottlerMu.Unlock()

	if spamThrottler != nil {
		return errors.New("spam throttler was already running")
	}

	if conf == nil {
		conf = DefaultSpamThrottlerConfig
	}

	if err := validateConfig(conf); err != nil {
		return err
	}

	t := &throttler{
		config:     conf,
		candidates: make(map[common.Address]int),
		throttled:  make(map[common.Address]int),
		allowed:    make(map[common.Address]bool),
		mu:         new(sync.RWMutex),
		threshold:  conf.InitialThreshold,
		throttleCh: make(chan *types.Transaction, conf.ThrottleTPS*5),
		quitCh:     make(chan struct{}),
	}

	go pool.throttleLoop(t)

	spamThrottler = t
	logger.Info("Start spam throttler", "config", *conf)
	return nil
}

func (pool *TxPool) StopSpamThrottler() {
	spamThrottlerMu.Lock()
	defer spamThrottlerMu.Unlock()

	if spamThrottler != nil {
		close(spamThrottler.quitCh)
	}

	spamThrottler = nil
	candidateSizeGauge.Update(0)
	throttledSizeGauge.Update(0)
	allowedSizeGauge.Update(0)
	throttlerUpdateTimeGauge.Update(0)
	throttlerDropCount.Clear()
}

// handleTxMsg calls TxPool.AddRemotes by retrieving transactions from TxPool.txMsgCh.
func (pool *TxPool) handleTxMsg() {
	defer pool.wg.Done()

	for {
		select {
		case txs := <-pool.txMsgCh:
			pool.AddRemotes(txs)
		case <-pool.chainHeadSub.Err():
			return
		}
	}
}

// AddLocal enqueues a single transaction into the pool if it is valid, marking
// the sender as a local one in the mean time, ensuring it goes around the local
// pricing constraints.
func (pool *TxPool) AddLocal(tx *types.Transaction) error {
	if tx.Type().IsChainDataAnchoring() && !pool.config.AllowLocalAnchorTx {
		return errNotAllowedAnchoringTx
	}

	pool.mu.RLock()
	poolSize := uint64(len(pool.all))
	pool.mu.RUnlock()
	if poolSize >= pool.config.ExecSlotsAll+pool.config.NonExecSlotsAll {
		return fmt.Errorf("txpool is full: %d", poolSize)
	}
	return pool.addTx(tx, !pool.config.NoLocals)
}

// AddRemote enqueues a single transaction into the pool if it is valid. If the
// sender is not among the locally tracked ones, full pricing constraints will
// apply.
func (pool *TxPool) AddRemote(tx *types.Transaction) error {
	return pool.addTx(tx, false)
}

// AddLocals enqueues a batch of transactions into the pool if they are valid,
// marking the senders as a local ones in the mean time, ensuring they go around
// the local pricing constraints.
func (pool *TxPool) AddLocals(txs []*types.Transaction) []error {
	return pool.checkAndAddTxs(txs, !pool.config.NoLocals)
}

// AddRemotes enqueues a batch of transactions into the pool if they are valid.
// If the senders are not among the locally tracked ones, full pricing constraints
// will apply.
func (pool *TxPool) AddRemotes(txs []*types.Transaction) []error {
	return pool.checkAndAddTxs(txs, false)
}

// checkAndAddTxs compares the size of given transactions and the capacity of TxPool.
// If given transactions exceed the capacity of TxPool, it slices the given transactions
// so it can fit into TxPool's capacity.
func (pool *TxPool) checkAndAddTxs(txs []*types.Transaction, local bool) []error {
	pool.mu.RLock()
	poolSize := uint64(len(pool.all))
	pool.mu.RUnlock()
	poolCapacity := int(pool.config.ExecSlotsAll + pool.config.NonExecSlotsAll - poolSize)
	numTxs := len(txs)

	if poolCapacity < numTxs {
		txs = txs[:poolCapacity]
	}

	errs := pool.addTxs(txs, local)

	if poolCapacity < numTxs {
		for i := 0; i < numTxs-poolCapacity; i++ {
			errs = append(errs, txPoolIsFullErr)
		}
	}

	return errs
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *TxPool) addTx(tx *types.Transaction, local bool) error {
	senderCacher.recover(pool.signer, []*types.Transaction{tx})

	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Try to inject the transaction and update any state
	replace, err := pool.add(tx, local)
	if err != nil {
		return err
	}
	// If we added a new transaction, run promotion checks and return
	if !replace {
		from, _ := types.Sender(pool.signer, tx) // already validated
		pool.promoteExecutables([]common.Address{from})
	}
	return nil
}

// addTxs attempts to queue a batch of transactions if they are valid.
func (pool *TxPool) addTxs(txs []*types.Transaction, local bool) []error {
	senderCacher.recover(pool.signer, txs)

	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTxsLocked(txs, local)
}

// addTxsLocked attempts to queue a batch of transactions if they are valid,
// whilst assuming the transaction pool lock is already held.
func (pool *TxPool) addTxsLocked(txs []*types.Transaction, local bool) []error {
	// Add the batch of transaction, tracking the accepted ones
	dirty := make(map[common.Address]struct{})
	errs := make([]error, len(txs))

	for i, tx := range txs {
		var replace bool
		if replace, errs[i] = pool.add(tx, local); errs[i] == nil {
			if !replace {
				from, _ := types.Sender(pool.signer, tx) // already validated
				dirty[from] = struct{}{}
			}
		}
	}

	// Only reprocess the internal state if something was actually added
	if len(dirty) > 0 {
		addrs := make([]common.Address, 0, len(dirty))
		for addr := range dirty {
			addrs = append(addrs, addr)
		}
		pool.promoteExecutables(addrs)
	}
	return errs
}

// Status returns the status (unknown/pending/queued) of a batch of transactions
// identified by their hashes.
func (pool *TxPool) Status(hashes []common.Hash) []TxStatus {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	status := make([]TxStatus, len(hashes))
	for i, hash := range hashes {
		if tx := pool.all[hash]; tx != nil {
			from, _ := types.Sender(pool.signer, tx) // already validated
			if pool.pending[from] != nil && pool.pending[from].txs.items[tx.Nonce()] != nil {
				status[i] = TxStatusPending
			} else {
				status[i] = TxStatusQueued
			}
		}
	}
	return status
}

// Get returns a transaction if it is contained in the pool
// and nil otherwise.
func (pool *TxPool) Get(hash common.Hash) *types.Transaction {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.all[hash]
}

// checkAndSetBeat sets the beat of the account if there is no beat of the account.
func (pool *TxPool) checkAndSetBeat(addr common.Address) {
	_, exist := pool.beats[addr]

	if !exist {
		pool.beats[addr] = time.Now()
	}
}

// removeTx removes a single transaction from the queue, moving all subsequent
// transactions back to the future queue.
func (pool *TxPool) removeTx(hash common.Hash, outofbound bool) {
	// Fetch the transaction we wish to delete
	tx, ok := pool.all[hash]
	if !ok {
		return
	}
	addr, _ := types.Sender(pool.signer, tx) // already validated during insertion

	// Remove it from the list of known transactions
	delete(pool.all, hash)
	if outofbound {
		pool.priced.Removed()
	}
	// Remove the transaction from the pending lists and reset the account nonce
	if pending := pool.pending[addr]; pending != nil {
		if removed, invalids := pending.Remove(tx); removed {
			// If no more pending transactions are left, remove the list
			if pending.Empty() {
				delete(pool.pending, addr)
			}
			// Postpone any invalidated transactions
			for _, tx := range invalids {
				pool.enqueueTx(tx.Hash(), tx)
			}
			pool.updatePendingNonce(addr, tx.Nonce())
			return
		}
	}
	// Transaction is in the future queue
	if future := pool.queue[addr]; future != nil {
		future.Remove(tx)
		if future.Empty() {
			delete(pool.queue, addr)
		}
	}
}

// promoteExecutables moves transactions that have become processable from the
// future queue to the set of pending transactions. During this process, all
// invalidated transactions (low nonce, low balance) are deleted.
func (pool *TxPool) promoteExecutables(accounts []common.Address) {
	pool.txMu.Lock()
	defer pool.txMu.Unlock()
	// Track the promoted transactions to broadcast them at once
	var promoted []*types.Transaction

	// Gather all the accounts potentially needing updates
	if accounts == nil {
		accounts = make([]common.Address, 0, len(pool.queue))
		for addr := range pool.queue {
			accounts = append(accounts, addr)
		}
	}

	// Iterate over all accounts and promote any executable transactions
	for _, addr := range accounts {
		list := pool.queue[addr]
		if list == nil {
			continue // Just in case someone calls with a non existing account
		}
		// Drop all transactions that are deemed too old (low nonce)
		for _, tx := range list.Forward(pool.getNonce(addr)) {
			hash := tx.Hash()
			logger.Trace("Removed old queued transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
		}
		// Drop all transactions that are too costly (low balance)
		drops, _ := list.Filter(pool.getBalance(addr), pool)
		for _, tx := range drops {
			hash := tx.Hash()
			logger.Trace("Removed unpayable queued transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
			queuedNofundsCounter.Inc(1)
		}

		// Gather all executable transactions and promote them
		for _, tx := range list.Ready(pool.getPendingNonce(addr)) {
			hash := tx.Hash()
			if pool.promoteTx(addr, hash, tx) {
				logger.Trace("Promoting queued transaction", "hash", hash)
				promoted = append(promoted, tx)
			}
		}
		// Drop all transactions over the allowed limit
		if !pool.locals.contains(addr) {
			for _, tx := range list.Cap(int(pool.config.NonExecSlotsAccount)) {
				hash := tx.Hash()
				delete(pool.all, hash)
				pool.priced.Removed()
				queuedRateLimitCounter.Inc(1)
				logger.Trace("Removed cap-exceeding queued transaction", "hash", hash)
			}
		}
		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.queue, addr)
		}
	}
	// Notify subsystem for new promoted transactions.
	if len(promoted) > 0 {
		pool.txFeed.Send(NewTxsEvent{promoted})
	}
	// If the pending limit is overflown, start equalizing allowances
	pending := uint64(0)
	for _, list := range pool.pending {
		pending += uint64(list.Len())
	}

	if pending > pool.config.ExecSlotsAll {
		pendingBeforeCap := pending
		// Assemble a spam order to penalize large transactors first
		spammers := prque.New()
		for addr, list := range pool.pending {
			// Only evict transactions from high rollers
			if !pool.locals.contains(addr) && uint64(list.Len()) > pool.config.ExecSlotsAccount {
				spammers.Push(addr, int64(list.Len()))
			}
		}
		// Gradually drop transactions from offenders
		offenders := []common.Address{}
		for pending > pool.config.ExecSlotsAll && !spammers.Empty() {
			// Retrieve the next offender if not local address
			offender, _ := spammers.Pop()
			offenders = append(offenders, offender.(common.Address))

			// Equalize balances until all the same or below threshold
			if len(offenders) > 1 {
				// Calculate the equalization threshold for all current offenders
				threshold := pool.pending[offender.(common.Address)].Len()

				// Iteratively reduce all offenders until below limit or threshold reached
				for pending > pool.config.ExecSlotsAll && pool.pending[offenders[len(offenders)-2]].Len() > threshold {
					for i := 0; i < len(offenders)-1; i++ {
						list := pool.pending[offenders[i]]
						for _, tx := range list.Cap(list.Len() - 1) {
							// Drop the transaction from the global pools too
							hash := tx.Hash()
							delete(pool.all, hash)
							pool.priced.Removed()

							// Update the account nonce to the dropped transaction
							pool.updatePendingNonce(offenders[i], tx.Nonce())
							logger.Trace("Removed fairness-exceeding pending transaction", "hash", hash)
						}
						pending--
					}
				}
			}
		}
		// If still above threshold, reduce to limit or min allowance
		if pending > pool.config.ExecSlotsAll && len(offenders) > 0 {
			for pending > pool.config.ExecSlotsAll && uint64(pool.pending[offenders[len(offenders)-1]].Len()) > pool.config.ExecSlotsAccount {
				for _, addr := range offenders {
					list := pool.pending[addr]
					for _, tx := range list.Cap(list.Len() - 1) {
						// Drop the transaction from the global pools too
						hash := tx.Hash()
						delete(pool.all, hash)
						pool.priced.Removed()

						// Update the account nonce to the dropped transaction
						pool.updatePendingNonce(addr, tx.Nonce())
						logger.Trace("Removed fairness-exceeding pending transaction", "hash", hash)
					}
					pending--
				}
			}
		}
		pendingRateLimitCounter.Inc(int64(pendingBeforeCap - pending))
	}
	// If we've queued more transactions than the hard limit, drop oldest ones
	queued := uint64(0)
	for _, list := range pool.queue {
		queued += uint64(list.Len())
	}

	if queued > pool.config.NonExecSlotsAll {
		// Sort all accounts with queued transactions by heartbeat
		addresses := make(addresssByHeartbeat, 0, len(pool.queue))
		for addr := range pool.queue {
			if !pool.locals.contains(addr) { // don't drop locals
				addresses = append(addresses, addressByHeartbeat{addr, pool.beats[addr]})
			}
		}
		sort.Sort(addresses)

		// Drop transactions until the total is below the limit or only locals remain
		for drop := queued - pool.config.NonExecSlotsAll; drop > 0 && len(addresses) > 0; {
			addr := addresses[len(addresses)-1]
			list := pool.queue[addr.address]

			addresses = addresses[:len(addresses)-1]

			// Drop all transactions if they are less than the overflow
			if size := uint64(list.Len()); size <= drop {
				for _, tx := range list.Flatten() {
					pool.removeTx(tx.Hash(), true)
				}
				drop -= size
				queuedRateLimitCounter.Inc(int64(size))
				continue
			}
			// Otherwise drop only last few transactions
			txs := list.Flatten()
			for i := len(txs) - 1; i >= 0 && drop > 0; i-- {
				pool.removeTx(txs[i].Hash(), true)
				drop--
				queuedRateLimitCounter.Inc(1)
			}
		}
	}
}

// demoteUnexecutables removes invalid and processed transactions from the pools
// executable/pending queue and any subsequent transactions that become unexecutable
// are moved back into the future queue.
func (pool *TxPool) demoteUnexecutables() {
	pool.txMu.Lock()
	defer pool.txMu.Unlock()

	// full-validation count. demoteUnexecutables does full-validation for a limited number of txs.
	cnt := 0
	// Iterate over all accounts and demote any non-executable transactions
	for addr, list := range pool.pending {
		nonce := pool.getNonce(addr)
		var drops, invalids types.Transactions

		// Drop all transactions that are deemed too old (low nonce)
		for _, tx := range list.Forward(nonce) {
			hash := tx.Hash()
			logger.Trace("Removed old pending transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
		}

		// demoteUnexecutables does full-validation for a limited number of txs. Otherwise, it only validate nonce.
		// The logic below loosely checks the tx count for the efficiency and the simplicity.
		if cnt < demoteUnexecutablesFullValidationTxLimit {
			cnt += list.Len()
			drops, invalids = list.Filter(pool.getBalance(addr), pool)
		} else {
			drops, invalids = list.FilterUnexecutable()
		}

		// Drop all transactions that are unexecutable, and queue any invalids back for later
		for _, tx := range drops {
			hash := tx.Hash()
			logger.Trace("Removed unexecutable pending transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
			pendingNofundsCounter.Inc(1)
		}

		for _, tx := range invalids {
			hash := tx.Hash()
			logger.Trace("Demoting pending transaction", "hash", hash)
			pool.enqueueTx(hash, tx)
		}
		// If there's a gap in front, warn (should never happen) and postpone all transactions
		if list.Len() > 0 && list.txs.Get(nonce) == nil {
			for _, tx := range list.Cap(0) {
				hash := tx.Hash()
				logger.Error("Demoting invalidated transaction", "hash", hash)
				pool.enqueueTx(hash, tx)
			}
		}

		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.pending, addr)
		}
	}
}

// getNonce returns the nonce of the account from the cache. If it is not in the cache, it gets the nonce from the stateDB.
func (pool *TxPool) getNonce(addr common.Address) uint64 {
	return pool.currentState.GetNonce(addr)
}

// getBalance returns the balance of the account from the cache. If it is not in the cache, it gets the balance from the stateDB.
func (pool *TxPool) getBalance(addr common.Address) *big.Int {
	return pool.currentState.GetBalance(addr)
}

// GetPendingNonce is a method to check the last nonce value of pending in external API.
// Use getPendingNonce to get the nonce value inside txpool because it catches the lock.
func (pool *TxPool) GetPendingNonce(addr common.Address) uint64 {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.getPendingNonce(addr)
}

// getPendingNonce returns the canonical nonce for the managed or unmanaged account.
func (pool *TxPool) getPendingNonce(addr common.Address) uint64 {
	cNonce := pool.getNonce(addr)
	if pNonce, exist := pool.pendingNonce[addr]; !exist || pNonce < cNonce {
		pool.pendingNonce[addr] = cNonce
	}

	return pool.pendingNonce[addr]
}

// setPendingNonce sets the new canonical nonce for the managed state.
func (pool *TxPool) setPendingNonce(addr common.Address, nonce uint64) {
	pool.pendingNonce[addr] = nonce
}

// updatePendingNonce updates the account nonce to the dropped transaction.
func (pool *TxPool) updatePendingNonce(addr common.Address, nonce uint64) {
	if pool.getPendingNonce(addr) > nonce {
		pool.setPendingNonce(addr, nonce)
	}
}

// addressByHeartbeat is an account address tagged with its last activity timestamp.
type addressByHeartbeat struct {
	address   common.Address
	heartbeat time.Time
}

type addresssByHeartbeat []addressByHeartbeat

func (a addresssByHeartbeat) Len() int           { return len(a) }
func (a addresssByHeartbeat) Less(i, j int) bool { return a[i].heartbeat.Before(a[j].heartbeat) }
func (a addresssByHeartbeat) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// accountSet is simply a set of addresses to check for existence, and a signer
// capable of deriving addresses from transactions.
type accountSet struct {
	accounts map[common.Address]struct{}
	signer   types.Signer
}

// newAccountSet creates a new address set with an associated signer for sender
// derivations.
func newAccountSet(signer types.Signer) *accountSet {
	return &accountSet{
		accounts: make(map[common.Address]struct{}),
		signer:   signer,
	}
}

// contains checks if a given address is contained within the set.
func (as *accountSet) contains(addr common.Address) bool {
	_, exist := as.accounts[addr]
	return exist
}

// containsTx checks if the sender of a given tx is within the set. If the sender
// cannot be derived, this method returns false.
func (as *accountSet) containsTx(tx *types.Transaction) bool {
	if addr, err := types.Sender(as.signer, tx); err == nil {
		return as.contains(addr)
	}
	return false
}

// add inserts a new address into the set to track.
func (as *accountSet) add(addr common.Address) {
	as.accounts[addr] = struct{}{}
}

// txLookup is used internally by TxPool to track transactions while allowing lookup without
// mutex contention.
//
// Note, although this type is properly protected against concurrent access, it
// is **not** a type that should ever be mutated or even exposed outside of the
// transaction pool, since its internal state is tightly coupled with the pools
// internal mechanisms. The sole purpose of the type is to permit out-of-bound
// peeking into the pool in TxPool.Get without having to acquire the widely scoped
// TxPool.mu mutex.
type txLookup struct {
	all  map[common.Hash]*types.Transaction
	lock sync.RWMutex
}

// newTxLookup returns a new txLookup structure.
func newTxLookup() *txLookup {
	return &txLookup{
		all: make(map[common.Hash]*types.Transaction),
	}
}

// Range calls f on each key and value present in the map.
func (t *txLookup) Range(f func(hash common.Hash, tx *types.Transaction) bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for key, value := range t.all {
		if !f(key, value) {
			break
		}
	}
}

// Get returns a transaction if it exists in the lookup, or nil if not found.
func (t *txLookup) Get(hash common.Hash) *types.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.all[hash]
}

// Count returns the current number of items in the lookup.
func (t *txLookup) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.all)
}

// Add adds a transaction to the lookup.
func (t *txLookup) Add(tx *types.Transaction) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.all[tx.Hash()] = tx
}

// Remove removes a transaction from the lookup.
func (t *txLookup) Remove(hash common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	delete(t.all, hash)
}
