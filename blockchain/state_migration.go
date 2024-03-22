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
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/klaytn/klaytn/blockchain/types"

	"github.com/alecthomas/units"
	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/mclock"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
)

const (
	DefaultWarmUpMinLoad        = 90
	DefaultWarmUpReportInterval = time.Second * 10
)

var (
	stopWarmUpErr           = errors.New("warm-up terminate by StopWarmUp")
	blockChainStopWarmUpErr = errors.New("warm-up terminate as blockchain stopped")
)

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

func (bc *BlockChain) concurrentRead(db state.Database, quitCh chan struct{}, hashCh chan common.Hash, resultCh chan statedb.SyncResult) {
	for {
		select {
		case <-quitCh:
			return
		case hash := <-hashCh:
			data, err := db.TrieDB().NodeFromOld(hash.ExtendZero())
			if err != nil {
				data, err = db.ContractCode(hash)
			}
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
	bc.migrationErr = nil
	defer func() {
		bc.migrationErr = returnErr
		// If migration stops by quit signal, it doesn't finish migration and it it will restart again.
		if returnErr != ErrQuitBySignal {
			// lock to prevent from a conflict of state DB close and state DB write
			bc.mu.Lock()
			bc.db.FinishStateMigration(returnErr == nil)
			bc.mu.Unlock()
		}
	}()

	start := time.Now()

	srcState := bc.StateCache()
	dstState := state.NewDatabase(bc.db)

	// collect preimages
	preimages := bc.db.CollectPreimages(func(k common.Hash) []byte {
		return dstState.TrieDB().Preimage(k)
	})

	// NOTE: lruCache is mandatory when state migration and block processing are executed simultaneously
	lruCache, _ := lru.New(int(2 * units.Giga / common.HashLength)) // 2GB for 62,500,000 common.Hash key values
	trieSync := state.NewStateSync(rootHash, dstState.TrieDB().DiskDB(), nil, lruCache, nil)
	var queue []common.Hash

	quitCh := make(chan struct{})
	defer close(quitCh)

	// Prepare concurrent read goroutines
	threads := runtime.NumCPU()
	hashCh := make(chan common.Hash, threads)
	resultCh := make(chan statedb.SyncResult, threads)

	for th := 0; th < threads; th++ {
		go bc.concurrentRead(srcState, quitCh, hashCh, resultCh)
	}

	stateTrieBatch := dstState.TrieDB().DiskDB().NewBatch(database.StateTrieDB)
	defer stateTrieBatch.Release()
	stats := migrationStats{initialStartTime: start, startTime: mclock.Now()}

	if bc.testMigrationHook != nil {
		bc.testMigrationHook()
	}

	// Migration main loop
	for trieSync.Pending() > 0 {
		nodes, _, codes := trieSync.Missing(1024)
		queue = append(queue[:0], append(nodes, codes...)...)
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
		for index, result := range results {
			if err := trieSync.Process(result); err != nil {
				logger.Error("State migration is failed by process error", "err", err)
				return fmt.Errorf("failed to process result #%d: %v", index, err)
			}
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
			logger.Info("State migration terminated by request")
			return errors.New("stop state migration")
		case <-bc.quit:
			logger.Info("State migration stopped by quit signal; should continue on node restart")
			return ErrQuitBySignal
		default:
		}

		bc.readCnt, bc.committedCnt, bc.pendingCnt, bc.progress = stats.totalRead, stats.totalCommitted, trieSync.Pending(), stats.progress
	}

	// Flush trie nodes which is not written yet.
	if err := stateTrieBatch.Write(); err != nil {
		logger.Error("State migration is failed by commit error", "err", err)
		return fmt.Errorf("DB write error: %v", err)
	}

	// Reinsert the collected images which are retrieved before migration
	bc.db.WritePreimages(0, preimages)

	stats.stateMigrationReport(true, trieSync.Pending(), trieSync.CalcProgressPercentage())
	bc.readCnt, bc.committedCnt, bc.pendingCnt, bc.progress = stats.totalRead, stats.totalCommitted, trieSync.Pending(), stats.progress

	// Clear memory of trieSync
	trieSync = nil

	elapsed := time.Since(start)
	speed := float64(stats.totalCommitted) / elapsed.Seconds()
	logger.Info("State migration : Copy is done",
		"totalRead", stats.totalRead, "totalCommitted", stats.totalCommitted,
		"totalElapsed", elapsed, "committed per second", speed)

	startCheck := time.Now()
	if err := state.CheckStateConsistencyParallel(srcState, dstState, rootHash, bc.quit); err != nil {
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
	oldTrie, err := statedb.NewSecureTrie(root, oldDB, nil)
	if err != nil {
		return nil, err
	}

	newTrie, err := statedb.NewSecureTrie(root, newDB, nil)
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

		bc.wg.Add(1)
		go func() {
			bc.migrateState(root)
			bc.wg.Done()
		}()
	}
}

// PrepareStateMigration sets prepareStateMigration to be called in checkStartStateMigration.
func (bc *BlockChain) PrepareStateMigration() error {
	if bc.db.ReadPruningEnabled() {
		return errors.New("state migration not supported with live pruning enabled")
	}

	if bc.db.InMigration() || bc.prepareStateMigration {
		return errors.New("migration already started")
	}

	bc.prepareStateMigration = true
	logger.Info("State migration is prepared", "expectedMigrationStartingBlockNumber", bc.CurrentBlock().NumberU64()+1)

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

// For tests starting and stopping node instances, clear residual migrationPrerequisites
// that might no longer work.
// TODO: remove this function when we have a better way to handle this.
// e.g. StartStateMigration directly calls CheckStakingInfoStored instead of this callback.
func ClearMigrationPrerequisites() {
	migrationPrerequisites = nil
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

	bc.wg.Add(1)
	go func() {
		bc.migrateState(root)
		bc.wg.Done()
	}()

	return nil
}

func (bc *BlockChain) StopStateMigration() error {
	if !bc.db.InMigration() {
		return errors.New("not in migration")
	}

	bc.stopStateMigration <- struct{}{}

	return nil
}

// StateMigrationStatus returns if it is in migration, the block number of in migration,
// number of committed blocks and number of pending blocks
func (bc *BlockChain) StateMigrationStatus() (bool, uint64, int, int, int, float64, error) {
	return bc.db.InMigration(), bc.db.MigrationBlockNumber(), bc.readCnt, bc.committedCnt, bc.pendingCnt, bc.progress, bc.migrationErr
}

// trieWarmUp runs state.Iterator, generated from the given state or storage trie node hash,
// until it reaches end. If it reaches end, it will send a nil error to errCh to indicate that
// it has been finished.
func (bc *BlockChain) trieWarmUp(next func() bool, resultCh chan int, errCh chan error) {
	var (
		resultErr    error
		reportTicker = time.NewTicker(DefaultWarmUpReportInterval)
		nReads       = 0
	)

	defer func() {
		resultCh <- nReads
		errCh <- resultErr
		reportTicker.Stop()
	}()

	for next() {
		select {
		case <-bc.quitWarmUp:
			return
		case <-bc.quit:
			resultErr = blockChainStopWarmUpErr
			return
		case <-reportTicker.C:
			resultCh <- nReads
			nReads = 1
		default:
			nReads++
		}
	}
}

// warmUpTrieCache receives errors from each warm-up goroutine.
// If it receives a nil error, it means a child goroutine is successfully terminated.
// It also periodically checks and logs warm-up progress.
func (bc *BlockChain) warmUpTrieCache(mainTrieDB *statedb.Database, minLoad uint, nChildren int,
	resultCh chan int, errCh chan error,
) {
	defer func() { bc.quitWarmUp = nil }()

	cache := mainTrieDB.TrieNodeCache()
	mainTrieCacheLimit := mainTrieDB.GetTrieNodeLocalCacheByteLimit()
	if minLoad == 0 {
		minLoad = DefaultWarmUpMinLoad
	}

	var (
		resultErr   error
		started     = time.Now()
		logged      = time.Now()
		context     []interface{}
		progress    uint64
		trieReadCnt int
	)

	updateContext := func() {
		switch c := cache.(type) {
		case *statedb.FastCache:
			stats := c.UpdateStats().(fastcache.Stats)
			progress = stats.BytesSize * 100 / mainTrieCacheLimit
			context = []interface{}{
				"warmUpCnt", trieReadCnt,
				"cacheLimit", units.Base2Bytes(mainTrieCacheLimit).String(),
				"cachedSize", units.Base2Bytes(stats.BytesSize).String(),
				"progress", progress,
				"elapsed", time.Since(started),
			}
		default:
			context = []interface{}{
				"warmUpCnt", trieReadCnt,
				"cacheLimit", units.Base2Bytes(mainTrieCacheLimit).String(),
			}
		}
	}

	for childCnt := 0; childCnt < nChildren; {
		select {
		case nReadNodes := <-resultCh:
			trieReadCnt += nReadNodes

			if time.Since(logged) < log.StatsReportLimit {
				continue
			}
			logged = time.Now()
			updateContext()
			if progress > uint64(minLoad) { // more than 90%
				close(bc.quitWarmUp)
				logger.Info("Warm up is completed", context...)
				return
			}
			logger.Info("Warm up progress", context...)
		case err := <-errCh:
			// if errCh returns nil, it means success.
			if err != nil {
				resultErr = err
				logger.Warn("Warm up got an error", "err", err)
			}

			logger.Debug("Warm up a child trie is finished", "childCnt", childCnt, "err", err)
			childCnt++
		}
	}

	updateContext()
	context = append(context, "resultErr", resultErr)
	logger.Info("Warm up is completed", context...)
}

// StartWarmUp retrieves all state/storage tries of the latest state root and caches the tries.
func (bc *BlockChain) StartWarmUp(minLoad uint) error {
	block, db, mainTrieDB, err := bc.prepareWarmUp()
	if err != nil {
		return err
	}
	// retrieve children nodes of state trie root node
	children, err := db.TrieDB().NodeChildren(block.Root().ExtendZero())
	if err != nil {
		return err
	}
	// run goroutine for each child node
	resultCh := make(chan int, len(children))
	errCh := make(chan error)
	bc.quitWarmUp = make(chan struct{})
	for idx, child := range children {
		childHash := child.Unextend()
		stateDB, err := state.New(childHash, db, nil, nil)
		if err != nil {
			logger.Warn("[WarmUp] Failed to get state",
				"rootHash", children, "childIdx", idx, "childHash", childHash.Hex())
			continue
		}
		it := state.NewNodeIterator(stateDB)
		go bc.trieWarmUp(it.Next, resultCh, errCh)
	}
	// run a warm-up checker routine
	go bc.warmUpTrieCache(mainTrieDB, minLoad, len(children), resultCh, errCh)
	logger.Info("State trie warm-up is started", "blockNum", block.NumberU64(),
		"root", block.Root().String(), "len(children)", len(children))
	return nil
}

// StopWarmUp stops the warming up process.
func (bc *BlockChain) StopWarmUp() error {
	if bc.quitWarmUp == nil {
		return ErrNotInWarmUp
	}

	close(bc.quitWarmUp)

	return nil
}

// StartCollectingTrieStats collects state or storage trie statistics.
func (bc *BlockChain) StartCollectingTrieStats(contractAddr common.Address) error {
	block := bc.GetBlockByNumber(bc.lastCommittedBlock)
	if block == nil {
		return fmt.Errorf("Block #%d not found", bc.lastCommittedBlock)
	}

	mainTrieDB := bc.StateCache().TrieDB()
	cache := mainTrieDB.TrieNodeCache()
	if cache == nil {
		return fmt.Errorf("target cache is nil")
	}
	db := state.NewDatabaseWithExistingCache(bc.db, cache)

	startNode := block.Root().ExtendZero()
	// If the contractAddr is given, start collecting stats from the root of storage trie
	if !common.EmptyAddress(contractAddr) {
		var err error
		startNode, err = bc.GetContractStorageRoot(block, db, contractAddr)
		if err != nil {
			logger.Error("Failed to get the contract storage root",
				"contractAddr", contractAddr.String(), "rootHash", block.Root().String(),
				"err", err)
			return err
		}
	}

	children, err := db.TrieDB().NodeChildren(startNode)
	if err != nil {
		logger.Error("Failed to retrieve the children of start node", "err", err)
		return err
	}

	logger.Info("Started collecting trie statistics",
		"blockNum", block.NumberU64(), "root", block.Root().String(), "len(children)", len(children))
	go collectTrieStats(db, startNode)

	return nil
}

// collectChildrenStats wraps CollectChildrenStats, in order to send finish signal to resultCh.
func collectChildrenStats(db state.Database, child common.ExtHash, resultCh chan<- statedb.NodeInfo) {
	db.TrieDB().CollectChildrenStats(child, 2, resultCh)
	resultCh <- statedb.NodeInfo{Finished: true}
}

// collectTrieStats is the main function of collecting trie statistics.
// It spawns goroutines for the upper-most children and iterates each sub-trie.
func collectTrieStats(db state.Database, startNode common.ExtHash) {
	children, err := db.TrieDB().NodeChildren(startNode)
	if err != nil {
		logger.Error("Failed to retrieve the children of start node", "err", err)
	}

	// collecting statistics by running individual goroutines for each child
	resultCh := make(chan statedb.NodeInfo, 10000)
	for _, child := range children {
		go collectChildrenStats(db, child, resultCh)
	}

	numGoRoutines := len(children)
	ticker := time.NewTicker(1 * time.Minute)

	numNodes, numLeafNodes, maxDepth := 0, 0, 0
	depthCounter := make(map[int]int)
	begin := time.Now()
	for {
		select {
		case result := <-resultCh:
			if result.Finished {
				numGoRoutines--
				if numGoRoutines == 0 {
					logger.Info("Finished collecting trie statistics", "elapsed", time.Since(begin),
						"numNodes", numNodes, "numLeafNodes", numLeafNodes, "maxDepth", maxDepth)
					printDepthStats(depthCounter)
					return
				}
				continue
			}
			numNodes++
			// if a leaf node, collect the depth data
			if result.Depth != 0 {
				numLeafNodes++
				depthCounter[result.Depth]++
				if result.Depth > maxDepth {
					maxDepth = result.Depth
				}
			}
		case <-ticker.C:
			// leave a periodic log
			logger.Info("Collecting trie statistics is in progress...", "elapsed", time.Since(begin),
				"numGoRoutines", numGoRoutines, "numNodes", numNodes, "numLeafNodes", numLeafNodes, "maxDepth", maxDepth)
			printDepthStats(depthCounter)
		}
	}
}

// printDepthStats leaves logs containing the depth and the number of nodes in the depth.
func printDepthStats(depthCounter map[int]int) {
	// max depth 20 is set by heuristically
	for depth := 2; depth < 20; depth++ {
		if depthCounter[depth] == 0 {
			continue
		}
		logger.Info("number of leaf nodes in a depth",
			"depth", depth, "numNodes", depthCounter[depth])
	}
}

// GetContractStorageRoot returns the storage root of a contract based on the given block.
func (bc *BlockChain) GetContractStorageRoot(block *types.Block, db state.Database, contractAddr common.Address) (common.ExtHash, error) {
	stateDB, err := state.New(block.Root(), db, nil, nil)
	if err != nil {
		return common.ExtHash{}, fmt.Errorf("failed to get StateDB - %w", err)
	}
	return stateDB.GetContractStorageRoot(contractAddr)
}

// prepareWarmUp creates and returns resources needed for state warm-up.
func (bc *BlockChain) prepareWarmUp() (*types.Block, state.Database, *statedb.Database, error) {
	// There is a chance of concurrent access to quitWarmUp, though not likely to happen.
	if bc.quitWarmUp != nil {
		return nil, nil, nil, fmt.Errorf("already warming up")
	}

	block := bc.GetBlockByNumber(bc.lastCommittedBlock)
	if block == nil {
		return nil, nil, nil, fmt.Errorf("block #%d not found", bc.lastCommittedBlock)
	}

	mainTrieDB := bc.StateCache().TrieDB()
	cache := mainTrieDB.TrieNodeCache()
	if cache == nil {
		return nil, nil, nil, fmt.Errorf("target cache is nil")
	}
	db := state.NewDatabaseWithExistingCache(bc.db, cache)
	return block, db, mainTrieDB, nil
}

func prepareContractWarmUp(block *types.Block, db state.Database, contractAddr common.Address) (common.ExtHash, state.Trie, error) {
	stateDB, err := state.New(block.Root(), db, nil, nil)
	if err != nil {
		return common.ExtHash{}, nil, fmt.Errorf("failed to get StateDB, err: %w", err)
	}
	storageTrieRoot, err := stateDB.GetContractStorageRoot(contractAddr)
	if err != nil {
		return common.ExtHash{}, nil, err
	}
	storageTrie, err := db.OpenStorageTrie(storageTrieRoot, nil)
	if err != nil {
		return common.ExtHash{}, nil, err
	}
	return storageTrieRoot, storageTrie, nil
}

// StartContractWarmUp retrieves a storage trie of the latest state root and caches the trie
// corresponding to the given contract address.
func (bc *BlockChain) StartContractWarmUp(contractAddr common.Address, minLoad uint) error {
	block, db, mainTrieDB, err := bc.prepareWarmUp()
	if err != nil {
		return err
	}
	// prepare contract storage trie specific resources - storageTrieRoot and storageTrie
	storageTrieRoot, storageTrie, err := prepareContractWarmUp(block, db, contractAddr)
	if err != nil {
		return fmt.Errorf("failed to prepare contract warm-up, err: %w", err)
	}
	// retrieve children nodes of contract storage trie root node
	children, err := db.TrieDB().NodeChildren(storageTrieRoot)
	if err != nil {
		return err
	}
	// run goroutine for each child node
	resultCh := make(chan int, len(children))
	errCh := make(chan error)
	bc.quitWarmUp = make(chan struct{})
	for _, child := range children {
		it := statedb.NewIterator(storageTrie.NodeIterator(child[:]))
		go bc.trieWarmUp(it.Next, resultCh, errCh)
	}
	// run a warm-up checker routine
	go bc.warmUpTrieCache(mainTrieDB, minLoad, len(children), resultCh, errCh)
	logger.Info("Contract storage trie warm-up is started",
		"blockNum", block.NumberU64(), "root", block.Root().String(), "contractAddr", contractAddr.String(),
		"contractStorageRoot", storageTrieRoot.String(), "len(children)", len(children))
	return nil
}
