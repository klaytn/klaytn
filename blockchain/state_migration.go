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

package blockchain

import (
	"errors"
	"fmt"
	"github.com/alecthomas/units"
	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/mclock"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type stateTrieMigrationDB struct {
	database.DBManager
}

func (td *stateTrieMigrationDB) ReadCachedTrieNode(hash common.Hash) ([]byte, error) {
	return td.ReadCachedTrieNodeFromNew(hash)
}
func (td *stateTrieMigrationDB) ReadCachedTrieNodePreimage(secureKey []byte) ([]byte, error) {
	return td.ReadCachedTrieNodePreimageFromNew(secureKey)
}

func (td *stateTrieMigrationDB) ReadStateTrieNode(key []byte) ([]byte, error) {
	return td.ReadStateTrieNodeFromNew(key)
}

func (td *stateTrieMigrationDB) HasStateTrieNode(key []byte) (bool, error) {
	return td.HasStateTrieNodeFromNew(key)
}

func (td *stateTrieMigrationDB) ReadPreimage(hash common.Hash) []byte {
	return td.ReadPreimageFromNew(hash)
}

func (bc *BlockChain) stateMigrationCommit(s *statedb.TrieSync, batch database.Batch) (int, error) {
	written, err := s.Commit(batch)
	if written == 0 || err != nil {
		return written, err
	}

	if batch.ValueSize() > database.IdealBatchSize {
		if err := batch.Write(); err != nil {
			return 0, fmt.Errorf("DB write error: %v", err)
		}
		batch.Reset()
	}

	return written, nil
}

func (bc *BlockChain) concurrentRead(db *statedb.Database, quitCh chan struct{}, hashCh chan common.Hash, resultCh chan statedb.SyncResult) {
	for {
		select {
		case <-quitCh:
			return
		case hash := <-hashCh:
			data, err := db.NodeFromOld(hash)
			if err != nil {
				resultCh <- statedb.SyncResult{Hash: hash, Err: err}
				continue
			}
			resultCh <- statedb.SyncResult{Hash: hash, Data: data}
		}
	}
}

// migrateState is the core implementation of state trie migration.
// This migrates a trie from StateTrieDB to StateTrieMigrationDB.
// Reading StateTrieDB happens in parallel and writing StateTrieMigrationDB happens in batch write.
//
// Before this function is called, StateTrieMigrationDB should be set.
// After the migration finish, the original StateTrieDB is removed and StateTrieMigrationDB becomes a new StateTrieDB.
func (bc *BlockChain) migrateState(rootHash common.Hash) (returnErr error) {
	bc.wg.Add(1)
	defer func() {
		bc.db.FinishStateMigration(returnErr == nil)
		bc.wg.Done()
	}()

	start := time.Now()

	srcState := bc.StateCache()
	targetState := state.NewDatabase(&stateTrieMigrationDB{bc.db})

	// NOTE: lruCache is mendatory when state migration and block processing are executed simultaneously
	lruCache, _ := lru.New(int(2 * units.Giga / common.HashLength)) // 2GB for 62,500,000 common.Hash key values
	trieSync := state.NewStateSync(rootHash, targetState.TrieDB().DiskDB(), nil, lruCache)
	var queue []common.Hash
	committedCnt := 0

	quitCh := make(chan struct{})
	defer close(quitCh)

	// Prepare concurrent read goroutines
	threads := runtime.NumCPU()
	hashCh := make(chan common.Hash, threads)
	resultCh := make(chan statedb.SyncResult, threads)

	for th := 0; th < threads; th++ {
		go bc.concurrentRead(srcState.TrieDB(), quitCh, hashCh, resultCh)
	}

	stateTrieBatch := targetState.TrieDB().DiskDB().NewBatch(database.StateTrieDB)
	stats := migrationStats{initialStartTime: start, startTime: mclock.Now()}

	// Migration main loop
	for trieSync.Pending() > 0 {
		bc.committedCnt, bc.pendingCnt = committedCnt, trieSync.Pending()
		queue = append(queue[:0], trieSync.Missing(1024)...)
		results := make([]statedb.SyncResult, len(queue))

		// Read the trie nodes
		startRead := time.Now()
		go func() {
			for _, hash := range queue {
				hashCh <- hash
			}
		}()

		for i := 0; i < len(queue); i++ {
			result := <-resultCh
			if result.Err != nil {
				logger.Error("State migration is failed by resultCh",
					"result.hash", result.Hash.String(), "result.Err", result.Err)
				return fmt.Errorf("failed to retrieve node data for %x: %v", result.Hash, result.Err)
			}
			results[i] = result
		}
		stats.read += len(queue)
		stats.readElapsed += time.Since(startRead)

		// Process trie nodes
		startProcess := time.Now()
		if _, index, err := trieSync.Process(results); err != nil {
			logger.Error("State migration is failed by process error", "err", err)
			return fmt.Errorf("failed to process result #%d: %v", index, err)
		}
		stats.processElapsed += time.Since(startProcess)

		// Commit trie nodes
		startWrite := time.Now()
		written, err := bc.stateMigrationCommit(trieSync, stateTrieBatch)
		if err != nil {
			logger.Error("State migration is failed by commit error", "err", err)
			return fmt.Errorf("failed to commit data #%d: %v", written, err)
		}
		stats.committed += written
		stats.writeElapsed += time.Since(startWrite)

		// Report progress
		stats.stateMigrationReport(false, trieSync.Pending(), trieSync.CalcProgressPercentage())

		select {
		case <-bc.stopStateMigration:
			logger.Error("State migration is failed by stop")
			return errors.New("stop state migration")
		case <-bc.quit:
			return nil
		default:
		}
	}

	// Flush trie nodes which is not written yet.
	if err := stateTrieBatch.Write(); err != nil {
		logger.Error("State migration is failed by commit error", "err", err)
		return fmt.Errorf("DB write error: %v", err)
	}

	stats.stateMigrationReport(true, trieSync.Pending(), trieSync.CalcProgressPercentage())

	// Clear memory of trieSync
	trieSync = nil

	elapsed := time.Since(start)
	speed := float64(stats.totalCommitted) / elapsed.Seconds()
	logger.Info("State migration : Copy is done",
		"totalRead", stats.totalRead, "totalCommitted", stats.totalCommitted,
		"totalElapsed", elapsed, "committed per second", speed)

	startCheck := time.Now()
	if err := state.CheckStateConsistency(srcState, targetState, rootHash, bc.committedCnt, bc.quit); err != nil {
		logger.Error("State migration : copied stateDB is invalid", "err", err)
		return err
	}
	checkElapsed := time.Since(startCheck)
	logger.Info("State migration is completed", "copyElapsed", elapsed, "checkElapsed", checkElapsed)
	return nil
}

// migrationStats tracks and reports on state migration.
type migrationStats struct {
	read, committed, totalRead, totalCommitted, pending int
	progress                                            float64
	initialStartTime                                    time.Time
	startTime                                           mclock.AbsTime
	readElapsed                                         time.Duration
	processElapsed                                      time.Duration
	writeElapsed                                        time.Duration
}

func (st *migrationStats) stateMigrationReport(force bool, pending int, progress float64) {
	var (
		now     = mclock.Now()
		elapsed = time.Duration(now) - time.Duration(st.startTime)
	)

	if force || elapsed >= log.StatsReportLimit {
		st.totalRead += st.read
		st.totalCommitted += st.committed
		st.pending, st.progress = pending, progress

		progressStr := strconv.FormatFloat(st.progress, 'f', 4, 64)
		progressStr = strings.TrimRight(progressStr, "0")
		progressStr = strings.TrimRight(progressStr, ".") + "%"

		logger.Info("State migration progress",
			"progress", progressStr,
			"totalRead", st.totalRead, "totalCommitted", st.totalCommitted, "pending", st.pending,
			"read", st.read, "readElapsed", st.readElapsed, "processElapsed", st.processElapsed,
			"written", st.committed, "writeElapsed", st.writeElapsed,
			"elapsed", common.PrettyDuration(elapsed),
			"totalElapsed", time.Since(st.initialStartTime))

		st.read, st.committed = 0, 0
		st.startTime = now
	}
}
func (bc *BlockChain) checkTrieContents(oldDB, newDB *statedb.Database, root common.Hash) ([]common.Address, error) {
	oldTrie, err := statedb.NewSecureTrie(root, oldDB)
	if err != nil {
		return nil, err
	}

	newTrie, err := statedb.NewSecureTrie(root, newDB)
	if err != nil {
		return nil, err
	}

	diff, _ := statedb.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := statedb.NewIterator(diff)

	var dirty []common.Address

	for iter.Next() {
		key := newTrie.GetKey(iter.Key)
		if key == nil {
			return nil, fmt.Errorf("no preimage found for hash %x", iter.Key)
		}

		dirty = append(dirty, common.BytesToAddress(key))
	}

	return dirty, nil
}

// restartStateMigration is called when a server is restarted while migration. The migration continues.
func (bc *BlockChain) restartStateMigration() {
	if bc.db.InMigration() {
		number := bc.db.MigrationBlockNumber()

		block := bc.GetBlockByNumber(number)
		if block == nil {
			logger.Error("failed to get migration block number", "blockNumber", number)
			return
		}

		root := block.Root()
		logger.Warn("State migration : Restarted", "blockNumber", number, "root", root.String())

		go bc.migrateState(root)
	}
}

// PrepareStateMigration sets prepareStateMigration to be called in checkStartStateMigration.
func (bc *BlockChain) PrepareStateMigration() error {
	if bc.db.InMigration() || bc.prepareStateMigration {
		return errors.New("migration already started")
	}

	bc.prepareStateMigration = true

	return nil
}

func (bc *BlockChain) checkStartStateMigration(number uint64, root common.Hash) bool {
	if bc.prepareStateMigration {
		logger.Info("State migration is started", "block", number, "root", root)

		if err := bc.StartStateMigration(number, root); err != nil {
			logger.Error("Failed to start state migration", "err", err)
		}

		bc.prepareStateMigration = false

		return true
	}

	return false
}

// migrationPrerequisites is a collection of functions that needs to be run
// before state trie migration. If one of the functions fails to run,
// the migration will not start.
var migrationPrerequisites []func(uint64) error

func RegisterMigrationPrerequisites(f func(uint64) error) {
	migrationPrerequisites = append(migrationPrerequisites, f)
}

// StartStateMigration checks prerequisites, configures DB and starts migration.
func (bc *BlockChain) StartStateMigration(number uint64, root common.Hash) error {
	if bc.db.InMigration() {
		return errors.New("migration already started")
	}

	for _, f := range migrationPrerequisites {
		if err := f(number); err != nil {
			return err
		}
	}

	if err := bc.db.CreateMigrationDBAndSetStatus(number); err != nil {
		return err
	}

	go bc.migrateState(root)

	return nil
}

func (bc *BlockChain) StopStateMigration() error {
	if !bc.db.InMigration() {
		return errors.New("not in migration")
	}

	bc.stopStateMigration <- struct{}{}

	return nil
}

// StatusStateMigration returns if it is in migration, the block number of in migration,
// number of committed blocks and number of pending blocks
func (bc *BlockChain) StatusStateMigration() (bool, uint64, int, int, float64) {
	return bc.db.InMigration(), bc.db.MigrationBlockNumber(), bc.committedCnt, bc.pendingCnt, bc.progress
}
