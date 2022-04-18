// Modifications Copyright 2019 The klaytn Authors
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

package bridgepool

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/rcrowley/go-metrics"
)

var logger = log.NewModuleLogger(log.ServiceChain)

var (
	ErrKnownTx           = errors.New("Known Transaction")
	ErrUnknownTx         = errors.New("Unknown Transaction")
	ErrDuplicatedNonceTx = errors.New("Duplicated Nonce Transaction")
)

// TODO-Klaytn-Servicechain Add Metrics
var (
	// Metrics for the pending pool
	refusedTxCounter = metrics.NewRegisteredCounter("bridgeTxpool/refuse", nil)
)

// BridgeTxPoolConfig are the configuration parameters of the transaction pool.
type BridgeTxPoolConfig struct {
	ParentChainID *big.Int
	Journal       string        // Journal of local transactions to survive node restarts
	Rejournal     time.Duration // Time interval to regenerate the local transaction journal

	GlobalQueue uint64 // Maximum number of non-executable transaction slots for all accounts
}

// DefaultBridgeTxPoolConfig contains the default configurations for the transaction
// pool.
var DefaultBridgeTxPoolConfig = BridgeTxPoolConfig{
	ParentChainID: big.NewInt(2018),
	Journal:       "bridge_transactions.rlp",
	Rejournal:     10 * time.Minute,
	GlobalQueue:   8192,
}

// sanitize checks the provided user configurations and changes anything that's
// unreasonable or unworkable.
func (config *BridgeTxPoolConfig) sanitize() BridgeTxPoolConfig {
	conf := *config
	if conf.Rejournal < time.Second {
		logger.Error("Sanitizing invalid bridgetxpool journal time", "provided", conf.Rejournal, "updated", time.Second)
		conf.Rejournal = time.Second
	}

	if conf.Journal == "" {
		logger.Error("Sanitizing invalid bridgetxpool journal file name", "updated", DefaultBridgeTxPoolConfig.Journal)
		conf.Journal = DefaultBridgeTxPoolConfig.Journal
	}

	return conf
}

// BridgeTxPool contains all currently known chain transactions.
type BridgeTxPool struct {
	config BridgeTxPoolConfig
	// TODO-Klaytn-Servicechain consider to remove singer. For now, caused of value transfer tx which don't have `from` value, I leave it.
	signer types.Signer
	mu     sync.RWMutex
	//txMu   sync.RWMutex // TODO-Klaytn-Servicechain: implement fine-grained locks

	journal *bridgeTxJournal // Journal of transaction to back up to disk

	queue map[common.Address]*ItemSortedMap // Queued but non-processable transactions
	// TODO-Klaytn-Servicechain refine heartbeat for the tx not for account.
	all map[common.Hash]*types.Transaction // All transactions to allow lookups

	wg     sync.WaitGroup // for shutdown sync
	closed chan struct{}
}

// NewBridgeTxPool creates a new transaction pool to gather, sort and filter inbound
// transactions from the network.
func NewBridgeTxPool(config BridgeTxPoolConfig) *BridgeTxPool {
	// Sanitize the input to ensure no vulnerable gas prices are set
	config = (&config).sanitize()

	// Create the transaction pool with its initial settings
	pool := &BridgeTxPool{
		config: config,
		queue:  make(map[common.Address]*ItemSortedMap),
		all:    make(map[common.Hash]*types.Transaction),
		closed: make(chan struct{}),
	}

	pool.SetLatestSigner(config.ParentChainID)

	// load from disk
	pool.journal = newBridgeTxJournal(config.Journal)

	if err := pool.journal.load(pool.AddLocals); err != nil {
		logger.Error("Failed to load chain transaction journal", "err", err)
	}
	if err := pool.journal.rotate(pool.Pending()); err != nil {
		logger.Error("Failed to rotate chain transaction journal", "err", err)
	}

	// Start the event loop and return
	pool.wg.Add(1)
	go pool.loop()

	return pool
}

// Deprecated: This function is deprecated. Use SetLatestSigner instead.
// SetEIP155Signer set signer of txpool.
func (pool *BridgeTxPool) SetEIP155Signer(chainID *big.Int) {
	pool.signer = types.NewEIP155Signer(chainID)
}

// SetLatestSigner set latest signer to txpool
func (pool *BridgeTxPool) SetLatestSigner(chainID *big.Int) {
	pool.signer = types.LatestSignerForChainID(chainID)
}

// loop is the transaction pool's main event loop, waiting for and reacting to
// outside blockchain events as well as for various reporting and transaction
// eviction events.
func (pool *BridgeTxPool) loop() {
	defer pool.wg.Done()

	journal := time.NewTicker(pool.config.Rejournal)
	defer journal.Stop()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle local transaction journal rotation
		case <-journal.C:
			if pool.journal != nil {
				pool.mu.Lock()
				if err := pool.journal.rotate(pool.pending()); err != nil {
					logger.Error("Failed to rotate local tx journal", "err", err)
				}
				pool.mu.Unlock()
			}
		case <-pool.closed:
			// update journal file with txs in pool.
			// if txpool close without the rotate process,
			// when loading the txpool with the journal file again,
			// there is a limit to the size of pool so that not all tx will be
			// loaded and especially the latest tx will not be loaded
			if pool.journal != nil {
				pool.mu.Lock()
				if err := pool.journal.rotate(pool.pending()); err != nil {
					logger.Error("Failed to rotate local tx journal", "err", err)
				}
				pool.mu.Unlock()
			}
			logger.Info("BridgeTxPool loop is closing")
			return
		}
	}
}

// Stop terminates the transaction pool.
func (pool *BridgeTxPool) Stop() {
	close(pool.closed)
	pool.wg.Wait()

	if pool.journal != nil {
		pool.journal.close()
	}
	logger.Info("Transaction pool stopped")
}

// Stats retrieves the current pool stats, namely the number of pending transactions.
func (pool *BridgeTxPool) Stats() int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	queued := 0
	for _, list := range pool.queue {
		queued += list.Len()
	}
	return queued
}

// Content retrieves the data content of the transaction pool, returning all the
// queued transactions, grouped by account and sorted by nonce.
func (pool *BridgeTxPool) Content() map[common.Address]types.Transactions {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	queued := make(map[common.Address]types.Transactions)
	for addr, list := range pool.queue {
		var queuedTxs []*types.Transaction
		txs := list.Flatten()
		for _, tx := range txs {
			queuedTxs = append(queuedTxs, tx.(*types.Transaction))
		}
		queued[addr] = queuedTxs
	}
	return queued
}

// GetTx get the tx by tx hash.
func (pool *BridgeTxPool) GetTx(txHash common.Hash) (*types.Transaction, error) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	tx, ok := pool.all[txHash]

	if ok {
		return tx, nil
	} else {
		return nil, ErrUnknownTx
	}
}

// Pending returns all pending transactions by calling internal pending method.
func (pool *BridgeTxPool) Pending() map[common.Address]types.Transactions {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pending := pool.pending()
	return pending
}

// Pending retrieves all pending transactions, grouped by origin
// account and sorted by nonce.
func (pool *BridgeTxPool) pending() map[common.Address]types.Transactions {
	pending := make(map[common.Address]types.Transactions)
	for addr, list := range pool.queue {
		var pendingTxs []*types.Transaction
		txs := list.Flatten()
		for _, tx := range txs {
			pendingTxs = append(pendingTxs, tx.(*types.Transaction))
		}
		pending[addr] = pendingTxs
	}
	return pending
}

// PendingTxsByAddress retrieves pending transactions of from. They are sorted by nonce.
func (pool *BridgeTxPool) PendingTxsByAddress(from *common.Address, limit int) types.Transactions {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	var pendingTxs types.Transactions

	if list, exist := pool.queue[*from]; exist {
		txs := list.FlattenByCount(limit)
		for _, tx := range txs {
			pendingTxs = append(pendingTxs, tx.(*types.Transaction))
		}
		return pendingTxs
	}
	return nil
}

// PendingTxHashesByAddress retrieves pending transaction hashes of from. They are sorted by nonce.
func (pool *BridgeTxPool) PendingTxHashesByAddress(from *common.Address, limit int) []common.Hash {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if list, exist := pool.queue[*from]; exist {
		pendingTxHashes := make([]common.Hash, limit)
		txs := list.FlattenByCount(limit)
		for _, tx := range txs {
			pendingTxHashes = append(pendingTxHashes, tx.(*types.Transaction).Hash())
		}
		return pendingTxHashes
	}
	return nil
}

// GetMaxTxNonce finds max nonce of the address.
func (pool *BridgeTxPool) GetMaxTxNonce(from *common.Address) uint64 {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	maxNonce := uint64(0)
	if list, exist := pool.queue[*from]; exist {
		for _, t := range list.items {
			if maxNonce < t.Nonce() {
				maxNonce = t.Nonce()
			}
		}
	}
	return maxNonce
}

// add validates a transaction and inserts it into the non-executable queue for
// later pending promotion and execution. If the transaction is a replacement for
// an already pending or queued one, it overwrites the previous and returns this
// so outer code doesn't uselessly call promote.
func (pool *BridgeTxPool) add(tx *types.Transaction) error {
	// If the transaction is already known, discard it
	hash := tx.Hash()
	if pool.all[hash] != nil {
		logger.Trace("Discarding already known transaction", "hash", hash)
		return ErrKnownTx
	}

	from, err := types.Sender(pool.signer, tx)
	if err != nil {
		return err
	}

	if uint64(len(pool.all)) >= pool.config.GlobalQueue {
		logger.Trace("Rejecting a new Tx, because BridgeTxPool is full and there is no room for the account", "hash", tx.Hash(), "account", from)
		refusedTxCounter.Inc(1)
		return fmt.Errorf("txpool is full: %d", uint64(len(pool.all)))
	}

	if pool.queue[from] == nil {
		pool.queue[from] = NewItemSortedMap(UnlimitedItemSortedMap)
	} else {
		if pool.queue[from].Get(tx.Nonce()) != nil {
			return ErrDuplicatedNonceTx
		}
	}

	pool.queue[from].Put(tx)

	if pool.all[hash] == nil {
		pool.all[hash] = tx
	}

	// Mark journal transactions
	pool.journalTx(from, tx)

	logger.Trace("Pooled new future transaction", "hash", hash, "from", from, "to", tx.To())
	return nil
}

// journalTx adds the specified transaction to the local disk journal if it is
// deemed to have been sent from a service chain account.
func (pool *BridgeTxPool) journalTx(from common.Address, tx *types.Transaction) {
	// Only journal if it's enabled
	if pool.journal == nil {
		return
	}
	if err := pool.journal.insert(tx); err != nil {
		logger.Error("Failed to journal local transaction", "err", err)
	}
}

// AddLocal enqueues a single transaction into the pool if it is valid, marking
// the sender as a local one.
func (pool *BridgeTxPool) AddLocal(tx *types.Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTx(tx)
}

// AddLocals enqueues a batch of transactions into the pool if they are valid,
// marking the senders as a local ones.
func (pool *BridgeTxPool) AddLocals(txs []*types.Transaction) []error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTxs(txs)
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *BridgeTxPool) addTx(tx *types.Transaction) error {
	//senderCacher.recover(pool.signer, []*types.Transaction{tx})
	// Try to inject the transaction and update any state
	return pool.add(tx)
}

// addTxs attempts to queue a batch of transactions if they are valid.
func (pool *BridgeTxPool) addTxs(txs []*types.Transaction) []error {
	//senderCacher.recover(pool.signer, txs)
	// Add the batch of transaction, tracking the accepted ones
	errs := make([]error, len(txs))

	for i, tx := range txs {
		errs[i] = pool.add(tx)
	}

	return errs
}

// Get returns a transaction if it is contained in the pool
// and nil otherwise.
func (pool *BridgeTxPool) Get(hash common.Hash) *types.Transaction {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.all[hash]
}

// removeTx removes a single transaction from the queue.
func (pool *BridgeTxPool) removeTx(hash common.Hash) error {
	// Fetch the transaction we wish to delete
	tx, ok := pool.all[hash]
	if !ok {
		return ErrUnknownTx
	}

	addr, err := types.Sender(pool.signer, tx)
	if err != nil {
		return err
	}

	// Remove it from the list of known transactions
	delete(pool.all, hash)

	// Transaction is in the future queue
	if future := pool.queue[addr]; future != nil {
		future.Remove(tx.Nonce())
		if future.Len() == 0 {
			delete(pool.queue, addr)
		}
	}

	return nil
}

// Remove removes transactions from the queue.
func (pool *BridgeTxPool) Remove(txs types.Transactions) []error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	errs := make([]error, len(txs))
	for i, tx := range txs {
		errs[i] = pool.removeTx(tx.Hash())
	}
	return errs
}

// RemoveTx removes a single transaction from the queue.
func (pool *BridgeTxPool) RemoveTx(tx *types.Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	err := pool.removeTx(tx.Hash())
	if err != nil {
		logger.Error("RemoveTx", "err", err)
		return err
	}
	return nil
}
