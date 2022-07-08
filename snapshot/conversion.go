// Modifications Copyright 2021 The klaytn Authors
// Copyright 2020 The go-ethereum Authors
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
// This file is derived from core/state/snapshot/conversion.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"encoding/binary"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
)

// trieKV represents a trie key-value pair
type trieKV struct {
	key   common.Hash
	value []byte
}

type (
	// trieGeneratorFn is the interface of trie generation which can
	// be implemented by different trie algorithm.
	trieGeneratorFn func(in chan (trieKV), out chan (common.Hash))

	// leafCallbackFn is the callback invoked at the leaves of the trie,
	// returns the subtrie root with the specified subtrie identifier.
	leafCallbackFn func(accountHash, codeHash common.Hash, stat *generateStats) (common.Hash, error)
)

// TODO-Klaytn-Snapshot port GenerateAccountTrieRoot/GenerateStorageTrieRoot/GenerateTrie

// generateStats is a collection of statistics gathered by the trie generator
// for logging purposes.
type generateStats struct {
	head  common.Hash
	start time.Time

	accounts uint64 // Number of accounts done (including those being crawled)
	slots    uint64 // Number of storage slots done (including those being crawled)

	slotsStart map[common.Hash]time.Time   // Start time for account slot crawling
	slotsHead  map[common.Hash]common.Hash // Slot head for accounts being crawled

	lock sync.RWMutex
}

// newGenerateStats creates a new generator stats.
func newGenerateStats() *generateStats {
	return &generateStats{
		slotsStart: make(map[common.Hash]time.Time),
		slotsHead:  make(map[common.Hash]common.Hash),
		start:      time.Now(),
	}
}

// progressAccounts updates the generator stats for the account range.
func (stat *generateStats) progressAccounts(account common.Hash, done uint64) {
	stat.lock.Lock()
	defer stat.lock.Unlock()

	stat.accounts += done
	stat.head = account
}

// finishAccounts updates the gemerator stats for the finished account range.
func (stat *generateStats) finishAccounts(done uint64) {
	stat.lock.Lock()
	defer stat.lock.Unlock()

	stat.accounts += done
}

// progressContract updates the generator stats for a specific in-progress contract.
func (stat *generateStats) progressContract(account common.Hash, slot common.Hash, done uint64) {
	stat.lock.Lock()
	defer stat.lock.Unlock()

	stat.slots += done
	stat.slotsHead[account] = slot
	if _, ok := stat.slotsStart[account]; !ok {
		stat.slotsStart[account] = time.Now()
	}
}

// finishContract updates the generator stats for a specific just-finished contract.
func (stat *generateStats) finishContract(account common.Hash, done uint64) {
	stat.lock.Lock()
	defer stat.lock.Unlock()

	stat.slots += done
	delete(stat.slotsHead, account)
	delete(stat.slotsStart, account)
}

// report prints the cumulative progress statistic smartly.
func (stat *generateStats) report() {
	stat.lock.RLock()
	defer stat.lock.RUnlock()

	ctx := []interface{}{
		"accounts", stat.accounts,
		"slots", stat.slots,
		"elapsed", common.PrettyDuration(time.Since(stat.start)),
	}
	if stat.accounts > 0 {
		// If there's progress on the account trie, estimate the time to finish crawling it
		if done := binary.BigEndian.Uint64(stat.head[:8]) / stat.accounts; done > 0 {
			var (
				left  = (math.MaxUint64 - binary.BigEndian.Uint64(stat.head[:8])) / stat.accounts
				speed = done/uint64(time.Since(stat.start)/time.Millisecond+1) + 1 // +1s to avoid division by zero
				eta   = time.Duration(left/speed) * time.Millisecond
			)
			// If there are large contract crawls in progress, estimate their finish time
			for acc, head := range stat.slotsHead {
				start := stat.slotsStart[acc]
				if done := binary.BigEndian.Uint64(head[:8]); done > 0 {
					var (
						left  = math.MaxUint64 - binary.BigEndian.Uint64(head[:8])
						speed = done/uint64(time.Since(start)/time.Millisecond+1) + 1 // +1s to avoid division by zero
					)
					// Override the ETA if larger than the largest until now
					if slotETA := time.Duration(left/speed) * time.Millisecond; eta < slotETA {
						eta = slotETA
					}
				}
			}
			ctx = append(ctx, []interface{}{
				"eta", common.PrettyDuration(eta),
			}...)
		}
	}
	logger.Info("Iterating state snapshot", ctx...)
}

// reportDone prints the last log when the whole generation is finished.
func (stat *generateStats) reportDone() {
	stat.lock.RLock()
	defer stat.lock.RUnlock()

	var ctx []interface{}
	ctx = append(ctx, []interface{}{"accounts", stat.accounts}...)
	if stat.slots != 0 {
		ctx = append(ctx, []interface{}{"slots", stat.slots}...)
	}
	ctx = append(ctx, []interface{}{"elapsed", common.PrettyDuration(time.Since(stat.start))}...)
	logger.Info("Iterated snapshot", ctx...)
}

// runReport periodically prints the progress information.
func runReport(stats *generateStats, stop chan bool) {
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			stats.report()
			timer.Reset(time.Second * 8)
		case success := <-stop:
			if success {
				stats.reportDone()
			}
			return
		}
	}
}

// generateTrieRoot generates the trie hash based on the snapshot iterator.
// It can be used for generating account trie, storage trie or even the
// whole state which connects the accounts and the corresponding storages.
func generateTrieRoot(it Iterator, accountHash common.Hash, generatorFn trieGeneratorFn, leafCallback leafCallbackFn, stats *generateStats, report bool) (common.Hash, error) {
	var (
		in      = make(chan trieKV)         // chan to pass leaves
		out     = make(chan common.Hash, 1) // chan to collect result
		stoplog = make(chan bool, 1)        // 1-size buffer, works when logging is not enabled
		wg      sync.WaitGroup
	)
	// Spin up a go-routine for trie hash re-generation
	wg.Add(1)
	go func() {
		defer wg.Done()
		generatorFn(in, out)
	}()
	// Spin up a go-routine for progress logging
	if report && stats != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runReport(stats, stoplog)
		}()
	}
	// Create a semaphore to assign tasks and collect results through. We'll pre-
	// fill it with nils, thus using the same channel for both limiting concurrent
	// processing and gathering results.
	threads := runtime.NumCPU()
	results := make(chan error, threads)
	for i := 0; i < threads; i++ {
		results <- nil // fill the semaphore
	}
	// stop is a helper function to shutdown the background threads
	// and return the re-generated trie hash.
	stop := func(fail error) (common.Hash, error) {
		close(in)
		result := <-out
		for i := 0; i < threads; i++ {
			if err := <-results; err != nil && fail == nil {
				fail = err
			}
		}
		stoplog <- fail == nil

		wg.Wait()
		return result, fail
	}
	var (
		logged    = time.Now()
		processed = uint64(0)
		leaf      trieKV
	)
	// Start to feed leaves
	for it.Next() {
		if accountHash == (common.Hash{}) {
			var (
				err      error
				fullData []byte
			)
			if leafCallback == nil {
				fullData = it.(AccountIterator).Account()
			} else {
				// Wait until the semaphore allows us to continue, aborting if
				// a sub-task failed
				if err := <-results; err != nil {
					results <- nil // stop will drain the results, add a noop back for this error we just consumed
					return stop(err)
				}
				// Fetch the next account and process it concurrently
				serializer := account.NewAccountSerializer()
				if err := rlp.DecodeBytes(it.(AccountIterator).Account(), serializer); err != nil {
					logger.Error("Failed to decode an account from iterator", "err", err)
					return stop(err)
				}
				acc := serializer.GetAccount()
				go func(hash common.Hash) {
					contract, ok := acc.(*account.SmartContractAccount)
					if !ok {
						results <- nil
						return
					}
					subroot, err := leafCallback(hash, common.BytesToHash(contract.GetCodeHash()), stats)
					if err != nil {
						results <- err
						return
					}
					rootHash := contract.GetStorageRoot()
					if rootHash != subroot {
						results <- fmt.Errorf("invalid subroot(path %x), want %x, have %x", hash, rootHash, subroot)
						return
					}
					results <- nil
				}(it.Hash())
				fullData, err = rlp.EncodeToBytes(serializer)
				if err != nil {
					return stop(err)
				}
			}
			leaf = trieKV{it.Hash(), fullData}
		} else {
			leaf = trieKV{it.Hash(), common.CopyBytes(it.(StorageIterator).Slot())}
		}
		in <- leaf

		// Accumulate the generation statistic if it's required.
		processed++
		if time.Since(logged) > 3*time.Second && stats != nil {
			if accountHash == (common.Hash{}) {
				stats.progressAccounts(it.Hash(), processed)
			} else {
				stats.progressContract(accountHash, it.Hash(), processed)
			}
			logged, processed = time.Now(), 0
		}
	}
	// Commit the last part statistic.
	if processed > 0 && stats != nil {
		if accountHash == (common.Hash{}) {
			stats.finishAccounts(processed)
		} else {
			stats.finishContract(accountHash, processed)
		}
	}
	return stop(nil)
}

func trieGenerate(in chan trieKV, out chan common.Hash) {
	db := statedb.NewDatabase(database.NewMemoryDBManager())
	t, _ := statedb.NewTrie(common.Hash{}, db)
	for leaf := range in {
		t.TryUpdate(leaf.key[:], leaf.value)
	}
	var root common.Hash
	if db == nil {
		root = t.Hash()
	} else {
		root, _ = t.Commit(nil)
	}
	out <- root
}
