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
// This file is derived from core/blockchain.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	mrand "math/rand"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/snapshot"

	"github.com/go-redis/redis/v7"
	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/mclock"
	"github.com/klaytn/klaytn/common/prque"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/log"
	klaytnmetrics "github.com/klaytn/klaytn/metrics"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/rcrowley/go-metrics"
)

// If total insertion time of a block exceeds insertTimeLimit,
// that time will be logged by blockLongInsertTimeGauge.
const insertTimeLimit = common.PrettyDuration(time.Second)

var (
	accountReadTimer   = klaytnmetrics.NewRegisteredHybridTimer("state/account/reads", nil)
	accountHashTimer   = klaytnmetrics.NewRegisteredHybridTimer("state/account/hashes", nil)
	accountUpdateTimer = klaytnmetrics.NewRegisteredHybridTimer("state/account/updates", nil)
	accountCommitTimer = klaytnmetrics.NewRegisteredHybridTimer("state/account/commits", nil)

	storageReadTimer   = klaytnmetrics.NewRegisteredHybridTimer("state/storage/reads", nil)
	storageHashTimer   = klaytnmetrics.NewRegisteredHybridTimer("state/storage/hashes", nil)
	storageUpdateTimer = klaytnmetrics.NewRegisteredHybridTimer("state/storage/updates", nil)
	storageCommitTimer = klaytnmetrics.NewRegisteredHybridTimer("state/storage/commits", nil)

	snapshotAccountReadTimer = metrics.NewRegisteredTimer("state/snapshot/account/reads", nil)
	snapshotStorageReadTimer = metrics.NewRegisteredTimer("state/snapshot/storage/reads", nil)
	snapshotCommitTimer      = metrics.NewRegisteredTimer("state/snapshot/commits", nil)

	blockInsertTimer    = klaytnmetrics.NewRegisteredHybridTimer("chain/inserts", nil)
	blockProcessTimer   = klaytnmetrics.NewRegisteredHybridTimer("chain/process", nil)
	blockExecutionTimer = klaytnmetrics.NewRegisteredHybridTimer("chain/execution", nil)
	blockFinalizeTimer  = klaytnmetrics.NewRegisteredHybridTimer("chain/finalize", nil)
	blockValidateTimer  = klaytnmetrics.NewRegisteredHybridTimer("chain/validate", nil)
	blockAgeTimer       = klaytnmetrics.NewRegisteredHybridTimer("chain/age", nil)

	blockPrefetchExecuteTimer   = klaytnmetrics.NewRegisteredHybridTimer("chain/prefetch/executes", nil)
	blockPrefetchInterruptMeter = metrics.NewRegisteredMeter("chain/prefetch/interrupts", nil)

	ErrNoGenesis            = errors.New("genesis not found in chain")
	ErrNotExistNode         = errors.New("the node does not exist in cached node")
	ErrQuitBySignal         = errors.New("quit by signal")
	ErrNotInWarmUp          = errors.New("not in warm up")
	logger                  = log.NewModuleLogger(log.Blockchain)
	kesCachePrefixBlockLogs = []byte("blockLogs")
)

// Below is the list of the constants for cache size.
// TODO-Klaytn: Below should be handled by ini or other configurations.
const (
	maxFutureBlocks     = 256
	maxTimeFutureBlocks = 30
	maxBadBlocks        = 10
	// TODO-Klaytn-Issue1911  This flag needs to be adjusted to the appropriate value.
	//  Currently, this value is taken to cache all 10 million accounts
	//  and should be optimized considering memory size and performance.
	maxAccountForCache = 10000000
)

const (
	DefaultTriesInMemory = 128
	// BlockChainVersion ensures that an incompatible database forces a resync from scratch.
	BlockChainVersion    = 3
	DefaultBlockInterval = 128
	MaxPrefetchTxs       = 20000
)

// CacheConfig contains the configuration values for the 1) stateDB caching and
// 2) trie caching/pruning resident in a blockchain.
type CacheConfig struct {
	// TODO-Klaytn-Issue1666 Need to check the benefit of trie caching.
	ArchiveMode          bool                         // If true, state trie is not pruned and always written to database
	CacheSize            int                          // Size of in-memory cache of a trie (MiB) to flush matured singleton trie nodes to disk
	BlockInterval        uint                         // Block interval to flush the trie. Each interval state trie will be flushed into disk
	TriesInMemory        uint64                       // Maximum number of recent state tries according to its block number
	SenderTxHashIndexing bool                         // Enables saving senderTxHash to txHash mapping information to database and cache
	TrieNodeCacheConfig  *statedb.TrieNodeCacheConfig // Configures trie node cache
	SnapshotCacheSize    int                          // Memory allowance (MB) to use for caching snapshot entries in memory
}

// gcBlock is used for priority queue for GC.
type gcBlock struct {
	root     common.Hash
	blockNum uint64
}

// BlockChain represents the canonical chain given a database with a genesis
// block. The Blockchain manages chain imports, reverts, chain reorganisations.
//
// Importing blocks in to the block chain happens according to the set of rules
// defined by the two stage Validator. Processing of blocks is done using the
// Processor which processes the included transaction. The validation of the state
// is done in the second part of the Validator. Failing results in aborting of
// the import.
//
// The BlockChain also helps in returning blocks from **any** chain included
// in the database as well as blocks that represents the canonical chain. It's
// important to note that GetBlock can return any block and does not need to be
// included in the canonical one where as GetBlockByNumber always represents the
// canonical chain.
type BlockChain struct {
	chainConfig   *params.ChainConfig // Chain & network configuration
	chainConfigMu *sync.RWMutex
	cacheConfig   *CacheConfig // stateDB caching and trie caching/pruning configuration

	db      database.DBManager // Low level persistent database to store final content in
	snaps   *snapshot.Tree     // Snapshot tree for fast trie leaf access
	triegc  *prque.Prque       // Priority queue mapping block numbers to tries to gc
	chBlock chan gcBlock       // chPushBlockGCPrque is a channel for delivering the gc item to gc loop.

	hc            *HeaderChain
	rmLogsFeed    event.Feed
	chainFeed     event.Feed
	chainSideFeed event.Feed
	chainHeadFeed event.Feed
	logsFeed      event.Feed
	scope         event.SubscriptionScope
	genesisBlock  *types.Block

	mu sync.RWMutex // global mutex for locking chain operations

	checkpoint       int          // checkpoint counts towards the new checkpoint
	currentBlock     atomic.Value // Current head of the block chain
	currentFastBlock atomic.Value // Current head of the fast-sync chain (may be above the block chain!)

	stateCache   state.Database // State database to reuse between imports (contains state cache)
	futureBlocks *lru.Cache     // future blocks are blocks added for later processing

	quit    chan struct{} // blockchain quit channel
	running int32         // running must be called atomically
	// procInterrupt must be atomically called
	procInterrupt int32          // interrupt signaler for block processing
	wg            sync.WaitGroup // chain processing wait group for shutting down

	engine     consensus.Engine
	processor  Processor  // block processor interface
	prefetcher Prefetcher // Block state prefetcher interface
	validator  Validator  // block and state validator interface
	vmConfig   vm.Config

	badBlocks *lru.Cache // Bad block cache

	parallelDBWrite bool // TODO-Klaytn-Storage parallelDBWrite will be replaced by number of goroutines when worker pool pattern is introduced.

	// State migration
	prepareStateMigration bool
	stopStateMigration    chan struct{}
	readCnt               int
	committedCnt          int
	pendingCnt            int
	progress              float64
	migrationErr          error

	// Warm up
	lastCommittedBlock uint64
	quitWarmUp         chan struct{}

	prefetchTxCh chan prefetchTx
}

// prefetchTx is used to prefetch transactions, when fetcher works.
type prefetchTx struct {
	ti                int
	block             *types.Block
	followupInterrupt *uint32
}

// NewBlockChain returns a fully initialised block chain using information
// available in the database. It initialises the default Klaytn validator and
// Processor.
func NewBlockChain(db database.DBManager, cacheConfig *CacheConfig, chainConfig *params.ChainConfig, engine consensus.Engine, vmConfig vm.Config) (*BlockChain, error) {
	if cacheConfig == nil {
		cacheConfig = &CacheConfig{
			ArchiveMode:         false,
			CacheSize:           512,
			BlockInterval:       DefaultBlockInterval,
			TriesInMemory:       DefaultTriesInMemory,
			TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
			SnapshotCacheSize:   512,
		}
	}

	if cacheConfig.TrieNodeCacheConfig == nil {
		cacheConfig.TrieNodeCacheConfig = statedb.GetEmptyTrieNodeCacheConfig()
	}

	state.EnabledExpensive = db.GetDBConfig().EnableDBPerfMetrics

	// Initialize DeriveSha implementation
	InitDeriveSha(chainConfig.DeriveShaImpl)

	futureBlocks, _ := lru.New(maxFutureBlocks)
	badBlocks, _ := lru.New(maxBadBlocks)

	bc := &BlockChain{
		chainConfig:        chainConfig,
		chainConfigMu:      new(sync.RWMutex),
		cacheConfig:        cacheConfig,
		db:                 db,
		triegc:             prque.New(),
		chBlock:            make(chan gcBlock, 1000),
		stateCache:         state.NewDatabaseWithNewCache(db, cacheConfig.TrieNodeCacheConfig),
		quit:               make(chan struct{}),
		futureBlocks:       futureBlocks,
		engine:             engine,
		vmConfig:           vmConfig,
		badBlocks:          badBlocks,
		parallelDBWrite:    db.IsParallelDBWrite(),
		stopStateMigration: make(chan struct{}),
		prefetchTxCh:       make(chan prefetchTx, MaxPrefetchTxs),
	}

	// set hardForkBlockNumberConfig which will be used as a global variable
	if err := fork.SetHardForkBlockNumberConfig(bc.chainConfig); err != nil {
		return nil, err
	}

	bc.validator = NewBlockValidator(chainConfig, bc, engine)
	bc.prefetcher = newStatePrefetcher(chainConfig, bc, engine)
	bc.processor = NewStateProcessor(chainConfig, bc, engine)

	var err error
	bc.hc, err = NewHeaderChain(db, chainConfig, engine, bc.getProcInterrupt)
	if err != nil {
		return nil, err
	}
	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		return nil, ErrNoGenesis
	}
	var nilBlock *types.Block
	bc.currentBlock.Store(nilBlock)
	bc.currentFastBlock.Store(nilBlock)

	if err := bc.loadLastState(); err != nil {
		return nil, err
	}
	// Make sure the state associated with the block is available
	head := bc.CurrentBlock()
	if _, err := state.New(head.Root(), bc.stateCache, bc.snaps); err != nil {
		// Head state is missing, before the state recovery, find out the
		// disk layer point of snapshot(if it's enabled). Make sure the
		// rewound point is lower than disk layer.
		var diskRoot common.Hash
		if bc.cacheConfig.SnapshotCacheSize > 0 {
			diskRoot = bc.db.ReadSnapshotRoot()
		}
		if diskRoot != (common.Hash{}) {
			logger.Warn("Head state missing, repairing", "number", head.Number(), "hash", head.Hash(), "snaproot", diskRoot)

			snapDisk, err := bc.setHeadBeyondRoot(head.NumberU64(), diskRoot, true)
			if err != nil {
				return nil, err
			}

			// Chain rewound, persist old snapshot number to indicate recovery procedure
			if snapDisk != 0 {
				bc.db.WriteSnapshotRecoveryNumber(snapDisk)
			}
		} else {
			// Dangling block without a state associated, init from scratch
			logger.Warn("Head state missing, repairing chain",
				"number", head.NumberU64(), "hash", head.Hash().String())
			if _, err := bc.setHeadBeyondRoot(head.NumberU64(), common.Hash{}, true); err != nil {
				return nil, err
			}
		}
	}
	// Check the current state of the block hashes and make sure that we do not have any of the bad blocks in our chain
	for hash := range BadHashes {
		if header := bc.GetHeaderByHash(hash); header != nil {
			// get the canonical block corresponding to the offending header's number
			headerByNumber := bc.GetHeaderByNumber(header.Number.Uint64())
			// make sure the headerByNumber (if present) is in our current canonical chain
			if headerByNumber != nil && headerByNumber.Hash() == header.Hash() {
				logger.Error("Found bad hash, rewinding chain", "number", header.Number, "hash", header.ParentHash)
				bc.SetHead(header.Number.Uint64() - 1)
				logger.Error("Chain rewind was successful, resuming normal operation")
			}
		}
	}

	// Load any existing snapshot, regenerating it if loading failed
	if bc.cacheConfig.SnapshotCacheSize > 0 {
		// If the chain was rewound past the snapshot persistent layer (causing
		// a recovery block number to be persisted to disk), check if we're still
		// in recovery mode and in that case, don't invalidate the snapshot on a
		// head mismatch.
		var recover bool

		head := bc.CurrentBlock()
		if layer := bc.db.ReadSnapshotRecoveryNumber(); layer != nil && *layer > head.NumberU64() {
			logger.Warn("Enabling snapshot recovery", "chainhead", head.NumberU64(), "diskbase", *layer)
			recover = true
		}
		bc.snaps, _ = snapshot.New(bc.db, bc.stateCache.TrieDB(), bc.cacheConfig.SnapshotCacheSize, head.Root(), false, true, recover)
	}

	for i := 1; i <= bc.cacheConfig.TrieNodeCacheConfig.NumFetcherPrefetchWorker; i++ {
		bc.wg.Add(1)
		go bc.prefetchTxWorker(i)
	}
	logger.Info("prefetchTxWorkers are started", "num", bc.cacheConfig.TrieNodeCacheConfig.NumFetcherPrefetchWorker)

	// Take ownership of this particular state
	go bc.update()
	bc.gcCachedNodeLoop()
	bc.restartStateMigration()

	if cacheConfig.TrieNodeCacheConfig.DumpPeriodically() {
		logger.Info("LocalCache is used for trie node cache, start saving cache to file periodically",
			"dir", bc.cacheConfig.TrieNodeCacheConfig.FastCacheFileDir,
			"period", bc.cacheConfig.TrieNodeCacheConfig.FastCacheSavePeriod)
		trieDB := bc.stateCache.TrieDB()
		bc.wg.Add(1)
		go func() {
			defer bc.wg.Done()
			trieDB.SaveCachePeriodically(bc.cacheConfig.TrieNodeCacheConfig, bc.quit)
		}()
	}

	return bc, nil
}

// prefetchTxWorker receives a block and a transaction index, which it pre-executes
// to retrieve and cache the data for the actual block processing.
func (bc *BlockChain) prefetchTxWorker(index int) {
	defer bc.wg.Done()

	logger.Debug("prefetchTxWorker is started", "index", index)
	var snaps *snapshot.Tree
	if bc.cacheConfig.TrieNodeCacheConfig.UseSnapshotForPrefetch {
		snaps = bc.snaps
	}
	for followup := range bc.prefetchTxCh {
		stateDB, err := state.NewForPrefetching(bc.CurrentBlock().Root(), bc.stateCache, snaps)
		if err != nil {
			logger.Debug("failed to retrieve stateDB for prefetchTxWorker", "err", err)
			continue
		}
		vmCfg := bc.vmConfig
		vmCfg.Prefetching = true
		bc.prefetcher.PrefetchTx(followup.block, followup.ti, stateDB, vmCfg, followup.followupInterrupt)
	}
	logger.Debug("prefetchTxWorker is terminated", "index", index)
}

// SetCanonicalBlock resets the canonical as the block with the given block number.
// It works as rewinding the head block to the previous one, but does not delete the data.
func (bc *BlockChain) SetCanonicalBlock(blockNum uint64) {
	// If the given block number is zero (it is zero by default), it does nothing
	if blockNum == 0 {
		return
	}
	// Read the block with the given block number and set it as canonical block
	targetBlock := bc.db.ReadBlockByNumber(blockNum)
	if targetBlock == nil {
		logger.Error("failed to retrieve the block", "blockNum", blockNum)
		return
	}
	bc.insert(targetBlock)
	if err := bc.loadLastState(); err != nil {
		logger.Error("failed to load last state after setting the canonical block", "err", err)
		return
	}
	// Make sure the state associated with the block is available
	head := bc.CurrentBlock()
	if _, err := state.New(head.Root(), bc.stateCache, bc.snaps); err != nil {
		// Dangling block without a state associated, init from scratch
		logger.Warn("Head state missing, repairing chain",
			"number", head.NumberU64(), "hash", head.Hash().String())
		if _, err := bc.setHeadBeyondRoot(head.NumberU64(), common.Hash{}, true); err != nil {
			logger.Error("Repairing chain is failed", "number", head.NumberU64(), "hash", head.Hash().String(), "err", err)
			return
		}
	}
	logger.Info("successfully set the canonical block", "blockNum", blockNum)
}

func (bc *BlockChain) UseGiniCoeff() bool {
	bc.chainConfigMu.RLock()
	defer bc.chainConfigMu.RUnlock()

	return bc.chainConfig.Governance.Reward.UseGiniCoeff
}

func (bc *BlockChain) SetUseGiniCoeff(val bool) {
	bc.chainConfigMu.Lock()
	defer bc.chainConfigMu.Unlock()

	bc.chainConfig.Governance.Reward.UseGiniCoeff = val
}

func (bc *BlockChain) ProposerPolicy() uint64 {
	bc.chainConfigMu.RLock()
	defer bc.chainConfigMu.RUnlock()

	return bc.chainConfig.Istanbul.ProposerPolicy
}

func (bc *BlockChain) SetProposerPolicy(val uint64) {
	bc.chainConfigMu.Lock()
	defer bc.chainConfigMu.Unlock()

	bc.chainConfig.Istanbul.ProposerPolicy = val
}

func (bc *BlockChain) getProcInterrupt() bool {
	return atomic.LoadInt32(&bc.procInterrupt) == 1
}

// loadLastState loads the last known chain state from the database. This method
// assumes that the chain manager mutex is held.
func (bc *BlockChain) loadLastState() error {
	// Restore the last known head block
	head := bc.db.ReadHeadBlockHash()
	if head == (common.Hash{}) {
		// Corrupt or empty database, init from scratch
		logger.Info("Empty database, resetting chain")
		return bc.Reset()
	}
	// Make sure the entire head block is available
	currentBlock := bc.GetBlockByHash(head)
	if currentBlock == nil {
		// Corrupt or empty database, init from scratch
		logger.Error("Head block missing, resetting chain", "hash", head.String())
		return bc.Reset()
	}
	// Everything seems to be fine, set as the head block
	bc.currentBlock.Store(currentBlock)
	bc.lastCommittedBlock = currentBlock.NumberU64()

	// Restore the last known head header
	currentHeader := currentBlock.Header()
	if head := bc.db.ReadHeadHeaderHash(); head != (common.Hash{}) {
		if header := bc.GetHeaderByHash(head); header != nil {
			currentHeader = header
		}
	}
	bc.hc.SetCurrentHeader(currentHeader)

	// Restore the last known head fast block
	bc.currentFastBlock.Store(currentBlock)
	if head := bc.db.ReadHeadFastBlockHash(); head != (common.Hash{}) {
		if block := bc.GetBlockByHash(head); block != nil {
			bc.currentFastBlock.Store(block)
		}
	}

	// Issue a status log for the user
	currentFastBlock := bc.CurrentFastBlock()

	headerTd := bc.GetTd(currentHeader.Hash(), currentHeader.Number.Uint64())
	blockTd := bc.GetTd(currentBlock.Hash(), currentBlock.NumberU64())
	fastTd := bc.GetTd(currentFastBlock.Hash(), currentFastBlock.NumberU64())

	logger.Info("Loaded most recent local header", "number", currentHeader.Number, "hash", currentHeader.Hash(), "td", headerTd, "age", common.PrettyAge(time.Unix(int64(currentHeader.Time.Uint64()), 0)))
	logger.Info("Loaded most recent local full block", "number", currentBlock.Number(), "hash", currentBlock.Hash(), "td", blockTd, "age", common.PrettyAge(time.Unix(int64(currentHeader.Time.Uint64()), 0)))
	logger.Info("Loaded most recent local fast block", "number", currentFastBlock.Number(), "hash", currentFastBlock.Hash(), "td", fastTd, "age", common.PrettyAge(time.Unix(int64(currentHeader.Time.Uint64()), 0)))

	return nil
}

// SetHead rewinds the local chain to a new head with the extra condition
// that the rewind must pass the specified state root. The method will try to
// delete minimal data from disk whilst retaining chain consistency.
func (bc *BlockChain) SetHead(head uint64) error {
	_, err := bc.setHeadBeyondRoot(head, common.Hash{}, false)
	return err
}

// setHeadBeyondRoot rewinds the local chain to a new head with the extra condition
// that the rewind must pass the specified state root. This method is meant to be
// used when rewinding with snapshots enabled to ensure that we go back further than
// persistent disk layer. Depending on whether the node was fast synced or full, and
// in which state, the method will try to delete minimal data from disk whilst
// retaining chain consistency.
//
// The method returns the block number where the requested root cap was found.
func (bc *BlockChain) setHeadBeyondRoot(head uint64, root common.Hash, repair bool) (uint64, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Track the block number of the requested root hash
	var rootNumber uint64 // (no root == always 0)

	updateFn := func(header *types.Header) error {
		// Rewind the block chain, ensuring we don't end up with a stateless head block
		if currentBlock := bc.CurrentBlock(); currentBlock != nil && header.Number.Uint64() <= currentBlock.NumberU64() {
			newHeadBlock := bc.GetBlock(header.Hash(), header.Number.Uint64())
			if newHeadBlock == nil {
				logger.Error("Gap in the chain, rewinding to genesis", "number", header.Number, "hash", header.Hash())
				newHeadBlock = bc.genesisBlock
			} else {
				// Block exists, keep rewinding until we find one with state,
				// keeping rewinding until we exceed the optional threshold
				// root hash
				beyondRoot := (root == common.Hash{}) // Flag whether we're beyond the requested root (no root, always true)

				for {
					// If a root threshold was requested but not yet crossed, check
					if root != (common.Hash{}) && !beyondRoot && newHeadBlock.Root() == root {
						beyondRoot, rootNumber = true, newHeadBlock.NumberU64()
					}
					if _, err := state.New(newHeadBlock.Root(), bc.stateCache, bc.snaps); err != nil {
						// Rewound state missing, rolled back to the parent block, reset to genesis
						logger.Trace("Block state missing, rewinding further", "number", newHeadBlock.NumberU64(), "hash", newHeadBlock.Hash())
						parent := bc.GetBlock(newHeadBlock.ParentHash(), newHeadBlock.NumberU64()-1)
						if parent != nil {
							newHeadBlock = parent
							continue
						}
						logger.Error("Missing block in the middle, aiming genesis", "number", newHeadBlock.NumberU64()-1, "hash", newHeadBlock.ParentHash())
						newHeadBlock = bc.genesisBlock
					}
					if beyondRoot || newHeadBlock.NumberU64() == 0 {
						logger.Debug("Rewound to block with state", "number", newHeadBlock.NumberU64(), "hash", newHeadBlock.Hash().String())
						break
					}
					// if newHeadBlock has state, then rewind first
					logger.Debug("Skipping block with threshold state", "number", newHeadBlock.NumberU64(), "hash", newHeadBlock.Hash().String(), "root", newHeadBlock.Root().String())
					newHeadBlock = bc.GetBlock(newHeadBlock.ParentHash(), newHeadBlock.NumberU64()-1) // Keep rewinding
				}

			}
			if newHeadBlock.NumberU64() == 0 {
				return errors.New("rewound to block number 0, but repair failed")
			}
			bc.db.WriteHeadBlockHash(newHeadBlock.Hash())

			// Degrade the chain markers if they are explicitly reverted.
			// In theory we should update all in-memory markers in the
			// last step, however the direction of SetHead is from high
			// to low, so it's safe the update in-memory markers directly.
			bc.currentBlock.Store(newHeadBlock)
			headBlockNumberGauge.Update(int64(newHeadBlock.NumberU64()))
		}

		// Rewind the fast block in a simpleton way to the target head
		if currentFastBlock := bc.CurrentFastBlock(); currentFastBlock != nil && header.Number.Uint64() < currentFastBlock.NumberU64() {
			newHeadFastBlock := bc.GetBlock(header.Hash(), header.Number.Uint64())
			// If either blocks reached nil, reset to the genesis state
			if newHeadFastBlock == nil {
				newHeadFastBlock = bc.genesisBlock
			}
			bc.db.WriteHeadFastBlockHash(newHeadFastBlock.Hash())

			// Degrade the chain markers if they are explicitly reverted.
			// In theory we should update all in-memory markers in the
			// last step, however the direction of SetHead is from high
			// to low, so it's safe the update in-memory markers directly.
			bc.currentFastBlock.Store(newHeadFastBlock)
		}
		return nil
	}

	// Rewind the header chain, deleting all block bodies until then
	delFn := func(hash common.Hash, num uint64) {
		// Remove relative body and receipts from the active store.
		// The header, total difficulty and canonical hash will be
		// removed in the hc.SetHead function.
		bc.db.DeleteBody(hash, num)
		bc.db.DeleteReceipts(hash, num)

	}

	// If SetHead was only called as a chain reparation method, try to skip
	// touching the header chain altogether
	if repair {
		if err := updateFn(bc.CurrentBlock().Header()); err != nil {
			return 0, err
		}
	} else {
		// Rewind the chain to the requested head and keep going backwards until a
		// block with a state is found
		logger.Warn("Rewinding blockchain", "target", head)
		if err := bc.hc.SetHead(head, updateFn, delFn); err != nil {
			return 0, err
		}
	}

	// Clear out any stale content from the caches
	bc.futureBlocks.Purge()
	bc.db.ClearBlockChainCache()
	//TODO-Klaytn add governance DB deletion logic.

	return rootNumber, bc.loadLastState()
}

// FastSyncCommitHead sets the current head block to the one defined by the hash
// irrelevant what the chain contents were prior.
func (bc *BlockChain) FastSyncCommitHead(hash common.Hash) error {
	// Make sure that both the block as well at its state trie exists
	block := bc.GetBlockByHash(hash)
	if block == nil {
		return fmt.Errorf("non existent block [%x…]", hash[:4])
	}
	if _, err := statedb.NewSecureTrie(block.Root(), bc.stateCache.TrieDB()); err != nil {
		return err
	}
	// If all checks out, manually set the head block
	bc.mu.Lock()
	bc.currentBlock.Store(block)
	bc.lastCommittedBlock = block.NumberU64()
	bc.mu.Unlock()

	logger.Info("Committed new head block", "number", block.Number(), "hash", hash)
	return nil
}

// CurrentBlock retrieves the current head block of the canonical chain. The
// block is retrieved from the blockchain's internal cache.
func (bc *BlockChain) CurrentBlock() *types.Block {
	return bc.currentBlock.Load().(*types.Block)
}

// CurrentFastBlock retrieves the current fast-sync head block of the canonical
// chain. The block is retrieved from the blockchain's internal cache.
func (bc *BlockChain) CurrentFastBlock() *types.Block {
	return bc.currentFastBlock.Load().(*types.Block)
}

// Validator returns the current validator.
func (bc *BlockChain) Validator() Validator {
	return bc.validator
}

// Processor returns the current processor.
func (bc *BlockChain) Processor() Processor {
	return bc.processor
}

// State returns a new mutable state based on the current HEAD block.
func (bc *BlockChain) State() (*state.StateDB, error) {
	return bc.StateAt(bc.CurrentBlock().Root())
}

// StateAt returns a new mutable state based on a particular point in time.
func (bc *BlockChain) StateAt(root common.Hash) (*state.StateDB, error) {
	return state.New(root, bc.stateCache, bc.snaps)
}

// StateAtWithPersistent returns a new mutable state based on a particular point in time with persistent trie nodes.
func (bc *BlockChain) StateAtWithPersistent(root common.Hash) (*state.StateDB, error) {
	exist := bc.stateCache.TrieDB().DoesExistNodeInPersistent(root)
	if !exist {
		return nil, ErrNotExistNode
	}
	return state.New(root, bc.stateCache, bc.snaps)
}

// StateAtWithGCLock returns a new mutable state based on a particular point in time with read lock of the state nodes.
func (bc *BlockChain) StateAtWithGCLock(root common.Hash) (*state.StateDB, error) {
	bc.RLockGCCachedNode()

	exist := bc.stateCache.TrieDB().DoesExistCachedNode(root)
	if !exist {
		bc.RUnlockGCCachedNode()
		return nil, ErrNotExistNode
	}

	stateDB, err := state.New(root, bc.stateCache, bc.snaps)
	if err != nil {
		bc.RUnlockGCCachedNode()
		return nil, err
	}

	return stateDB, nil
}

// StateCache returns the caching database underpinning the blockchain instance.
func (bc *BlockChain) StateCache() state.Database {
	return bc.stateCache
}

// Reset purges the entire blockchain, restoring it to its genesis state.
func (bc *BlockChain) Reset() error {
	return bc.ResetWithGenesisBlock(bc.genesisBlock)
}

// ResetWithGenesisBlock purges the entire blockchain, restoring it to the
// specified genesis state.
func (bc *BlockChain) ResetWithGenesisBlock(genesis *types.Block) error {
	// Dump the entire block chain and purge the caches
	if err := bc.SetHead(0); err != nil {
		return err
	}
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Prepare the genesis block and reinitialise the chain
	bc.hc.WriteTd(genesis.Hash(), genesis.NumberU64(), genesis.BlockScore())
	bc.db.WriteBlock(genesis)

	bc.genesisBlock = genesis
	bc.insert(bc.genesisBlock)
	bc.currentBlock.Store(bc.genesisBlock)
	bc.hc.SetGenesis(bc.genesisBlock.Header())
	bc.hc.SetCurrentHeader(bc.genesisBlock.Header())
	bc.currentFastBlock.Store(bc.genesisBlock)

	return nil
}

// repair tries to repair the current blockchain by rolling back the current block
// until one with associated state is found. This is needed to fix incomplete db
// writes caused either by crashes/power outages, or simply non-committed tries.
//
// This method only rolls back the current block. The current header and current
// fast block are left intact.
// Deprecated: in order to repair chain, please use SetHead or setHeadBeyondRoot methods
func (bc *BlockChain) repair(head **types.Block) error {
	for {
		// Abort if we've rewound to a head block that does have associated state
		if _, err := state.New((*head).Root(), bc.stateCache, bc.snaps); err == nil {
			logger.Info("Rewound blockchain to past state", "number", (*head).Number(), "hash", (*head).Hash())
			return nil
		} else {
			// Should abort and return error, otherwise it will fall into infinite loop
			if (*head).NumberU64() == 0 {
				return errors.New("rewound to block number 0, but repair failed")
			} else {
				// If headBlockNumber > 0, rewind one block and recheck state availability there
				block := bc.GetBlock((*head).ParentHash(), (*head).NumberU64()-1)
				if block == nil {
					return fmt.Errorf("missing block %d [%x]", (*head).NumberU64()-1, (*head).ParentHash())
				}
				*head = block
			}
		}
	}
}

// Export writes the active chain to the given writer.
func (bc *BlockChain) Export(w io.Writer) error {
	return bc.ExportN(w, uint64(0), bc.CurrentBlock().NumberU64())
}

// ExportN writes a subset of the active chain to the given writer.
func (bc *BlockChain) ExportN(w io.Writer, first uint64, last uint64) error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if first > last {
		return fmt.Errorf("export failed: first (%d) is greater than last (%d)", first, last)
	}
	logger.Info("Exporting batch of blocks", "count", last-first+1)

	start, reported := time.Now(), time.Now()
	for nr := first; nr <= last; nr++ {
		block := bc.GetBlockByNumber(nr)
		if block == nil {
			return fmt.Errorf("export failed on #%d: not found", nr)
		}
		if err := block.EncodeRLP(w); err != nil {
			return err
		}
		if time.Since(reported) >= log.StatsReportLimit {
			logger.Info("Exporting blocks", "exported", block.NumberU64()-first, "elapsed", common.PrettyDuration(time.Since(start)))
			reported = time.Now()
		}
	}

	return nil
}

// insert injects a new head block into the current block chain. This method
// assumes that the block is indeed a true head. It will also reset the head
// header and the head fast sync block to this very same block if they are older
// or if they are on a different side chain.
//
// Note, this function assumes that the `mu` mutex is held!
func (bc *BlockChain) insert(block *types.Block) {

	// If the block is on a side chain or an unknown one, force other heads onto it too
	updateHeads := bc.db.ReadCanonicalHash(block.NumberU64()) != block.Hash()

	// Add the block to the canonical chain number scheme and mark as the head
	bc.db.WriteCanonicalHash(block.Hash(), block.NumberU64())
	bc.db.WriteHeadBlockHash(block.Hash())

	bc.currentBlock.Store(block)

	// If the block is better than our head or is on a different chain, force update heads
	if updateHeads {
		bc.hc.SetCurrentHeader(block.Header())
		bc.db.WriteHeadFastBlockHash(block.Hash())

		bc.currentFastBlock.Store(block)
	}
}

// Genesis retrieves the chain's genesis block.
func (bc *BlockChain) Genesis() *types.Block {
	return bc.genesisBlock
}

// GetBodyRLP retrieves a block body in RLP encoding from the database by hash,
// caching it if found.
func (bc *BlockChain) GetBodyRLP(hash common.Hash) rlp.RawValue {
	return bc.db.ReadBodyRLPByHash(hash)
}

// HasBlock checks if a block is fully present in the database or not.
func (bc *BlockChain) HasBlock(hash common.Hash, number uint64) bool {
	return bc.db.HasBlock(hash, number)
}

// HasState checks if state trie is fully present in the database or not.
func (bc *BlockChain) HasState(hash common.Hash) bool {
	_, err := bc.stateCache.OpenTrie(hash)
	return err == nil
}

// HasBlockAndState checks if a block and associated state trie is fully present
// in the database or not, caching it if present.
func (bc *BlockChain) HasBlockAndState(hash common.Hash, number uint64) bool {
	// Check first that the block itself is known
	block := bc.GetBlock(hash, number)
	if block == nil {
		return false
	}
	return bc.HasState(block.Root())
}

// ShouldTryInserting returns the state whether the block should be inserted.
// If a node doesn't have the given block or the block number of given block is higher than the node's head block, it can try inserting the block.
func (bc *BlockChain) ShouldTryInserting(block *types.Block) bool {
	return !bc.HasBlockAndState(block.Hash(), block.NumberU64()) || bc.CurrentBlock().NumberU64() < block.NumberU64()
}

// GetBlock retrieves a block from the database by hash and number,
// caching it if found.
func (bc *BlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.db.ReadBlock(hash, number)
}

// GetBlockByHash retrieves a block from the database by hash, caching it if found.
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	return bc.db.ReadBlockByHash(hash)
}

// GetBlockNumber retrieves a blockNumber from the database by hash, caching it if found.
func (bc *BlockChain) GetBlockNumber(hash common.Hash) *uint64 {
	return bc.hc.GetBlockNumber(hash)
}

// GetBlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (bc *BlockChain) GetBlockByNumber(number uint64) *types.Block {
	return bc.db.ReadBlockByNumber(number)
}

// GetTxAndLookupInfo retrieves a tx and lookup info for a given transaction hash.
func (bc *BlockChain) GetTxAndLookupInfo(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	tx, blockHash, blockNumber, index := bc.GetTxAndLookupInfoInCache(txHash)
	if tx == nil {
		tx, blockHash, blockNumber, index = bc.db.ReadTxAndLookupInfo(txHash)
	}
	return tx, blockHash, blockNumber, index
}

// GetTxLookupInfoAndReceipt retrieves a tx and lookup info and receipt for a given transaction hash.
func (bc *BlockChain) GetTxLookupInfoAndReceipt(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	tx, blockHash, blockNumber, index := bc.GetTxAndLookupInfo(txHash)
	if tx == nil {
		return nil, common.Hash{}, 0, 0, nil
	}

	receipt := bc.GetReceiptByTxHash(txHash)
	if receipt == nil {
		return nil, common.Hash{}, 0, 0, nil
	}

	return tx, blockHash, blockNumber, index, receipt
}

// GetTxLookupInfoAndReceiptInCache retrieves a tx and lookup info and receipt for a given transaction hash in cache.
func (bc *BlockChain) GetTxLookupInfoAndReceiptInCache(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	tx, blockHash, blockNumber, index := bc.GetTxAndLookupInfoInCache(txHash)
	if tx == nil {
		return nil, common.Hash{}, 0, 0, nil
	}

	receipt := bc.GetTxReceiptInCache(txHash)
	if receipt == nil {
		return nil, common.Hash{}, 0, 0, nil
	}

	return tx, blockHash, blockNumber, index, receipt
}

// GetReceiptsByBlockHash retrieves the receipts for all transactions with given block hash.
func (bc *BlockChain) GetReceiptsByBlockHash(blockHash common.Hash) types.Receipts {
	return bc.db.ReadReceiptsByBlockHash(blockHash)
}

// GetReceiptByTxHash retrieves a receipt for a given transaction hash.
func (bc *BlockChain) GetReceiptByTxHash(txHash common.Hash) *types.Receipt {
	receipt := bc.GetTxReceiptInCache(txHash)
	if receipt != nil {
		return receipt
	}

	tx, blockHash, _, index := bc.GetTxAndLookupInfo(txHash)
	if tx == nil {
		return nil
	}

	receipts := bc.GetReceiptsByBlockHash(blockHash)
	if len(receipts) <= int(index) {
		logger.Error("receipt index exceeds the size of receipts", "receiptIndex", index, "receiptsSize", len(receipts))
		return nil
	}
	return receipts[index]
}

// GetLogsByHash retrieves the logs for all receipts in a given block.
func (bc *BlockChain) GetLogsByHash(hash common.Hash) [][]*types.Log {
	receipts := bc.GetReceiptsByBlockHash(hash)
	if receipts == nil {
		return nil
	}

	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs
}

// TrieNode retrieves a blob of data associated with a trie node (or code hash)
// either from ephemeral in-memory cache, or from persistent storage.
func (bc *BlockChain) TrieNode(hash common.Hash) ([]byte, error) {
	return bc.stateCache.TrieDB().Node(hash)
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (bc *BlockChain) Stop() {
	if !atomic.CompareAndSwapInt32(&bc.running, 0, 1) {
		return
	}
	// Unsubscribe all subscriptions registered from blockchain
	bc.scope.Close()
	if bc.cacheConfig.TrieNodeCacheConfig.RedisSubscribeBlockEnable {
		bc.CloseBlockSubscriptionLoop()
	}

	close(bc.prefetchTxCh)
	close(bc.quit)
	atomic.StoreInt32(&bc.procInterrupt, 1)

	bc.wg.Wait()

	// Ensure that the entirety of the state snapshot is journalled to disk.
	var snapBase common.Hash
	if bc.snaps != nil {
		var err error
		if snapBase, err = bc.snaps.Journal(bc.CurrentBlock().Root()); err != nil {
			logger.Error("Failed to journal state snapshot", "err", err)
		}
	}

	triedb := bc.stateCache.TrieDB()
	if !bc.isArchiveMode() {
		number := bc.CurrentBlock().NumberU64()
		recent := bc.GetBlockByNumber(number)
		if recent == nil {
			logger.Error("Failed to find recent block from persistent", "blockNumber", number)
			return
		}

		logger.Info("Writing cached state to disk", "block", recent.Number(), "hash", recent.Hash(), "root", recent.Root())
		if err := triedb.Commit(recent.Root(), true, number); err != nil {
			logger.Error("Failed to commit recent state trie", "err", err)
		}
		if snapBase != (common.Hash{}) {
			logger.Info("Writing snapshot state to disk", "root", snapBase)
			if err := triedb.Commit(snapBase, true, number); err != nil {
				logger.Error("Failed to commit recent state trie", "err", err)
			}
		}
		for !bc.triegc.Empty() {
			triedb.Dereference(bc.triegc.PopItem().(common.Hash))
		}
		if size, _ := triedb.Size(); size != 0 {
			logger.Error("Dangling trie nodes after full cleanup")
		}
	}
	if triedb.TrieNodeCache() != nil {
		_ = triedb.TrieNodeCache().Close()
	}

	logger.Info("Blockchain manager stopped")
}

func (bc *BlockChain) procFutureBlocks() {
	blocks := make([]*types.Block, 0, bc.futureBlocks.Len())
	for _, hash := range bc.futureBlocks.Keys() {
		hashKey, ok := hash.(common.CacheKey)
		if !ok {
			logger.Error("invalid key type", "expect", "common.CacheKey", "actual", reflect.TypeOf(hash))
			continue
		}

		if block, exist := bc.futureBlocks.Peek(hashKey); exist {
			cacheGetFutureBlockHitMeter.Mark(1)
			blocks = append(blocks, block.(*types.Block))
		} else {
			cacheGetFutureBlockMissMeter.Mark(1)
		}
	}
	if len(blocks) > 0 {
		types.BlockBy(types.Number).Sort(blocks)

		// Insert one by one as chain insertion needs contiguous ancestry between blocks
		for i := range blocks {
			bc.InsertChain(blocks[i : i+1])
		}
	}
}

// WriteStatus status of write
type WriteStatus byte

// TODO-Klaytn-Issue264 If we are using istanbul BFT, then we always have a canonical chain.
//                  Later we may be able to remove SideStatTy.
const (
	NonStatTy WriteStatus = iota
	CanonStatTy
	SideStatTy
)

// WriteResult includes the block write status and related statistics.
type WriteResult struct {
	Status         WriteStatus
	TotalWriteTime time.Duration
	TrieWriteTime  time.Duration
}

// Rollback is designed to remove a chain of links from the database that aren't
// certain enough to be valid.
func (bc *BlockChain) Rollback(chain []common.Hash) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for i := len(chain) - 1; i >= 0; i-- {
		hash := chain[i]

		currentHeader := bc.CurrentHeader()
		if currentHeader.Hash() == hash {
			bc.hc.SetCurrentHeader(bc.GetHeader(currentHeader.ParentHash, currentHeader.Number.Uint64()-1))
		}
		if currentFastBlock := bc.CurrentFastBlock(); currentFastBlock.Hash() == hash {
			newFastBlock := bc.GetBlock(currentFastBlock.ParentHash(), currentFastBlock.NumberU64()-1)
			bc.currentFastBlock.Store(newFastBlock)
			bc.db.WriteHeadFastBlockHash(newFastBlock.Hash())
		}
		if currentBlock := bc.CurrentBlock(); currentBlock.Hash() == hash {
			newBlock := bc.GetBlock(currentBlock.ParentHash(), currentBlock.NumberU64()-1)
			bc.currentBlock.Store(newBlock)
			bc.db.WriteHeadBlockHash(newBlock.Hash())
		}
	}
}

// SetReceiptsData computes all the non-consensus fields of the receipts
func SetReceiptsData(config *params.ChainConfig, block *types.Block, receipts types.Receipts) error {
	signer := types.MakeSigner(config, block.Number())

	transactions, logIndex := block.Transactions(), uint(0)
	if len(transactions) != len(receipts) {
		return errors.New("transaction and receipt count mismatch")
	}

	for j := 0; j < len(receipts); j++ {
		// The transaction hash can be retrieved from the transaction itself
		receipts[j].TxHash = transactions[j].Hash()

		// The contract address can be derived from the transaction itself
		if transactions[j].To() == nil {
			// Deriving the signer is expensive, only do if it's actually needed
			from, _ := types.Sender(signer, transactions[j])
			receipts[j].ContractAddress = crypto.CreateAddress(from, transactions[j].Nonce())
		}
		// The derived log fields can simply be set from the block and transaction
		for k := 0; k < len(receipts[j].Logs); k++ {
			receipts[j].Logs[k].BlockNumber = block.NumberU64()
			receipts[j].Logs[k].BlockHash = block.Hash()
			receipts[j].Logs[k].TxHash = receipts[j].TxHash
			receipts[j].Logs[k].TxIndex = uint(j)
			receipts[j].Logs[k].Index = logIndex
			logIndex++
		}
	}
	return nil
}

// InsertReceiptChain attempts to complete an already existing header chain with
// transaction and receipt data.
func (bc *BlockChain) InsertReceiptChain(blockChain types.Blocks, receiptChain []types.Receipts) (int, error) {
	bc.wg.Add(1)
	defer bc.wg.Done()

	// Do a sanity check that the provided chain is actually ordered and linked
	for i := 1; i < len(blockChain); i++ {
		if blockChain[i].NumberU64() != blockChain[i-1].NumberU64()+1 || blockChain[i].ParentHash() != blockChain[i-1].Hash() {
			logger.Error("Non contiguous receipt insert", "number", blockChain[i].Number(), "hash", blockChain[i].Hash(), "parent", blockChain[i].ParentHash(),
				"prevnumber", blockChain[i-1].Number(), "prevhash", blockChain[i-1].Hash())
			return 0, fmt.Errorf("non contiguous insert: item %d is #%d [%x…], item %d is #%d [%x…] (parent [%x…])", i-1, blockChain[i-1].NumberU64(),
				blockChain[i-1].Hash().Bytes()[:4], i, blockChain[i].NumberU64(), blockChain[i].Hash().Bytes()[:4], blockChain[i].ParentHash().Bytes()[:4])
		}
	}

	var (
		stats = struct{ processed, ignored int32 }{}
		start = time.Now()
		bytes = 0

		// TODO-Klaytn Needs to roll back if any one of batches fails
		bodyBatch            = bc.db.NewBatch(database.BodyDB)
		receiptsBatch        = bc.db.NewBatch(database.ReceiptsDB)
		txLookupEntriesBatch = bc.db.NewBatch(database.TxLookUpEntryDB)
	)
	for i, block := range blockChain {
		receipts := receiptChain[i]
		// Short circuit insertion if shutting down or processing failed
		if atomic.LoadInt32(&bc.procInterrupt) == 1 {
			return 0, nil
		}
		// Short circuit if the owner header is unknown
		if !bc.HasHeader(block.Hash(), block.NumberU64()) {
			return i, fmt.Errorf("containing header #%d [%x…] unknown", block.Number(), block.Hash().Bytes()[:4])
		}
		// Skip if the entire data is already known
		if bc.HasBlock(block.Hash(), block.NumberU64()) {
			stats.ignored++
			continue
		}
		// Compute all the non-consensus fields of the receipts
		if err := SetReceiptsData(bc.chainConfig, block, receipts); err != nil {
			return i, fmt.Errorf("failed to set receipts data: %v", err)
		}
		// Write all the data out into the database
		bc.db.PutBodyToBatch(bodyBatch, block.Hash(), block.NumberU64(), block.Body())
		bc.db.PutReceiptsToBatch(receiptsBatch, block.Hash(), block.NumberU64(), receipts)
		bc.db.PutTxLookupEntriesToBatch(txLookupEntriesBatch, block)

		stats.processed++

		totalBytes, err := database.WriteBatchesOverThreshold(bodyBatch, receiptsBatch, txLookupEntriesBatch)
		if err != nil {
			return 0, err
		} else {
			bytes += totalBytes
		}
	}

	totalBytes, err := database.WriteBatches(bodyBatch, receiptsBatch, txLookupEntriesBatch)
	if err != nil {
		return 0, err
	} else {
		bytes += totalBytes
	}

	// Update the head fast sync block if better
	bc.mu.Lock()
	head := blockChain[len(blockChain)-1]
	if td := bc.GetTd(head.Hash(), head.NumberU64()); td != nil { // Rewind may have occurred, skip in that case
		currentFastBlock := bc.CurrentFastBlock()
		if bc.GetTd(currentFastBlock.Hash(), currentFastBlock.NumberU64()).Cmp(td) < 0 {
			bc.db.WriteHeadFastBlockHash(head.Hash())
			bc.currentFastBlock.Store(head)
		}
	}
	bc.mu.Unlock()

	logger.Info("Imported new block receipts",
		"count", stats.processed,
		"elapsed", common.PrettyDuration(time.Since(start)),
		"number", head.Number(),
		"hash", head.Hash(),
		"size", common.StorageSize(bytes),
		"ignored", stats.ignored)
	return 0, nil
}

// WriteBlockWithoutState writes only the block and its metadata to the database,
// but does not write any state. This is used to construct competing side forks
// up to the point where they exceed the canonical total blockscore.
func (bc *BlockChain) WriteBlockWithoutState(block *types.Block, td *big.Int) {
	bc.wg.Add(1)
	defer bc.wg.Done()

	bc.hc.WriteTd(block.Hash(), block.NumberU64(), td)
	bc.writeBlock(block)
}

type TransactionLookup struct {
	Tx *types.Transaction
	*database.TxLookupEntry
}

// writeBlock writes block to persistent database.
// If write through caching is enabled, it also writes block to the cache.
func (bc *BlockChain) writeBlock(block *types.Block) {
	bc.db.WriteBlock(block)
}

// writeReceipts writes receipts to persistent database.
// If write through caching is enabled, it also writes blockReceipts to the cache.
func (bc *BlockChain) writeReceipts(hash common.Hash, number uint64, receipts types.Receipts) {
	bc.db.WriteReceipts(hash, number, receipts)
}

// writeStateTrie writes state trie to database if possible.
// If an archiving node is running, it always flushes state trie to DB.
// If not, it flushes state trie to DB periodically. (period = bc.cacheConfig.BlockInterval)
func (bc *BlockChain) writeStateTrie(block *types.Block, state *state.StateDB) error {
	state.LockGCCachedNode()
	defer state.UnlockGCCachedNode()

	root, err := state.Commit(true)
	if err != nil {
		return err
	}
	trieDB := bc.stateCache.TrieDB()
	trieDB.UpdateMetricNodes()

	// If we're running an archive node, always flush
	if bc.isArchiveMode() {
		if err := trieDB.Commit(root, false, block.NumberU64()); err != nil {
			return err
		}

		bc.checkStartStateMigration(block.NumberU64(), root)
		bc.lastCommittedBlock = block.NumberU64()
	} else {
		// Full but not archive node, do proper garbage collection
		trieDB.Reference(root, common.Hash{}) // metadata reference to keep trie alive

		// If we exceeded our memory allowance, flush matured singleton nodes to disk
		var (
			nodesSize, preimagesSize = trieDB.Size()
			nodesSizeLimit           = common.StorageSize(bc.cacheConfig.CacheSize) * 1024 * 1024
		)

		trieDBNodesSizeBytesGauge.Update(int64(nodesSize))
		trieDBPreimagesSizeGauge.Update(int64(preimagesSize))

		if nodesSize > nodesSizeLimit || preimagesSize > 4*1024*1024 {
			// NOTE-Klaytn Not to change the original behavior, error is not returned.
			// Error should be returned if it is thought to be safe in the future.
			if err := trieDB.Cap(nodesSizeLimit - database.IdealBatchSize); err != nil {
				logger.Error("Error from trieDB.Cap", "err", err, "limit", nodesSizeLimit-database.IdealBatchSize)
			}
		}

		if isCommitTrieRequired(bc, block.NumberU64()) {
			logger.Trace("Commit the state trie into the disk", "blocknum", block.NumberU64())
			if err := trieDB.Commit(block.Header().Root, true, block.NumberU64()); err != nil {
				return err
			}

			if bc.checkStartStateMigration(block.NumberU64(), root) {
				// flush referenced trie nodes out to new stateTrieDB
				if err := trieDB.Cap(0); err != nil {
					logger.Error("Error from trieDB.Cap by state migration", "err", err)
				}
			}

			bc.lastCommittedBlock = block.NumberU64()
		}

		bc.chBlock <- gcBlock{root, block.NumberU64()}
	}
	return nil
}

// RLockGCCachedNode locks the GC lock of CachedNode.
func (bc *BlockChain) RLockGCCachedNode() {
	bc.stateCache.RLockGCCachedNode()
}

// RUnlockGCCachedNode unlocks the GC lock of CachedNode.
func (bc *BlockChain) RUnlockGCCachedNode() {
	bc.stateCache.RUnlockGCCachedNode()
}

// DefaultTriesInMemory returns the number of tries residing in the memory.
func (bc *BlockChain) triesInMemory() uint64 {
	return bc.cacheConfig.TriesInMemory
}

// gcCachedNodeLoop runs a loop to gc.
func (bc *BlockChain) gcCachedNodeLoop() {
	trieDB := bc.stateCache.TrieDB()

	bc.wg.Add(1)
	go func() {
		defer bc.wg.Done()
		for {
			select {
			case block := <-bc.chBlock:
				bc.triegc.Push(block.root, -int64(block.blockNum))
				logger.Trace("Push GC block", "blkNum", block.blockNum, "hash", block.root.String())

				blkNum := block.blockNum
				if blkNum <= bc.triesInMemory() {
					continue
				}

				// Garbage collect anything below our required write retention
				chosen := blkNum - bc.triesInMemory()
				cnt := 0
				for !bc.triegc.Empty() {
					root, number := bc.triegc.Pop()
					if uint64(-number) > chosen {
						bc.triegc.Push(root, number)
						break
					}
					trieDB.Dereference(root.(common.Hash))
					cnt++
				}
				logger.Debug("GC cached node", "currentBlk", blkNum, "chosenBlk", chosen, "deferenceCnt", cnt)
			case <-bc.quit:
				return
			}
		}
	}()
}

func isCommitTrieRequired(bc *BlockChain, blockNum uint64) bool {
	if bc.prepareStateMigration {
		return true
	}

	// TODO-Klaytn-Issue1602 Introduce a simple and more concise way to determine commit trie requirements from governance
	if blockNum%uint64(bc.cacheConfig.BlockInterval) == 0 {
		return true
	}

	if bc.chainConfig.Istanbul != nil {
		return bc.ProposerPolicy() == params.WeightedRandom &&
			params.IsStakingUpdateInterval(blockNum)
	}
	return false
}

// isReorganizationRequired returns if reorganization is required or not based on total blockscore.
func isReorganizationRequired(localTd, externTd *big.Int, currentBlock, block *types.Block) bool {
	reorg := externTd.Cmp(localTd) > 0
	if !reorg && externTd.Cmp(localTd) == 0 {
		// Split same-blockscore blocks by number, then at random
		reorg = block.NumberU64() < currentBlock.NumberU64() || (block.NumberU64() == currentBlock.NumberU64() && mrand.Float64() < 0.5)
	}
	return reorg
}

// WriteBlockWithState writes the block and all associated state to the database.
// If we are to use writeBlockWithState alone, we should use mutex to protect internal state.
func (bc *BlockChain) WriteBlockWithState(block *types.Block, receipts []*types.Receipt, stateDB *state.StateDB) (WriteResult, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	return bc.writeBlockWithState(block, receipts, stateDB)
}

// writeBlockWithState writes the block and all associated state to the database.
// If BlockChain.parallelDBWrite is true, it calls writeBlockWithStateParallel.
// If not, it calls writeBlockWithStateSerial.
func (bc *BlockChain) writeBlockWithState(block *types.Block, receipts []*types.Receipt, stateDB *state.StateDB) (WriteResult, error) {
	var status WriteResult
	var err error
	if bc.parallelDBWrite {
		status, err = bc.writeBlockWithStateParallel(block, receipts, stateDB)
	} else {
		status, err = bc.writeBlockWithStateSerial(block, receipts, stateDB)
	}

	if err != nil {
		return status, err
	}

	// Publish the committed block to the redis cache of stateDB.
	// The cache uses the block to distinguish the latest state.
	if bc.cacheConfig.TrieNodeCacheConfig.RedisPublishBlockEnable {
		blockLogsKey := append(kesCachePrefixBlockLogs, block.Number().Bytes()...)
		bc.writeBlockLogsToRemoteCache(blockLogsKey, receipts)

		blockRlp, err := rlp.EncodeToBytes(block)
		if err != nil {
			logger.Error("failed to encode lastCommittedBlock", "blockNumber", block.NumberU64(), "err", err)
		}

		pubSub, ok := bc.stateCache.TrieDB().TrieNodeCache().(statedb.BlockPubSub)
		if ok {
			if err := pubSub.PublishBlock(hexutil.Encode(blockRlp)); err != nil {
				logger.Error("failed to publish block to redis", "blockNumber", block.NumberU64(), "err", err)
			}
		} else {
			logger.Error("invalid TrieNodeCache type", "trieNodeCacheConfig", bc.cacheConfig.TrieNodeCacheConfig)
		}
	}

	return status, err
}

// writeBlockLogsToRemoteCache writes block logs to remote cache.
// The stored logs will be used by KES service nodes to subscribe log events.
// This method is only for KES nodes.
func (bc *BlockChain) writeBlockLogsToRemoteCache(blockLogsKey []byte, receipts []*types.Receipt) {
	var entireBlockLogs []*types.LogForStorage
	for _, receipt := range receipts {
		for _, log := range receipt.Logs {
			// convert Log to LogForStorage to encode entire data
			entireBlockLogs = append(entireBlockLogs, (*types.LogForStorage)(log))
		}
	}
	encodedBlockLogs, err := rlp.EncodeToBytes(entireBlockLogs)
	if err != nil {
		logger.Error("rlp encoding error", "err", err)
		return
	}
	// TODO-Klaytn-KES: refine this not to use trieNodeCache
	cache, ok := bc.stateCache.TrieDB().TrieNodeCache().(*statedb.HybridCache)
	if !ok {
		logger.Error("only HybridCache supports block logs writing",
			"TrieNodeCacheType", reflect.TypeOf(bc.stateCache.TrieDB().TrieNodeCache()))
	} else {
		cache.Remote().Set(blockLogsKey, encodedBlockLogs)
	}
}

// writeBlockWithStateSerial writes the block and all associated state to the database in serial manner.
func (bc *BlockChain) writeBlockWithStateSerial(block *types.Block, receipts []*types.Receipt, state *state.StateDB) (WriteResult, error) {
	start := time.Now()
	bc.wg.Add(1)
	defer bc.wg.Done()

	var status WriteStatus
	// Calculate the total blockscore of the block
	ptd := bc.GetTd(block.ParentHash(), block.NumberU64()-1)
	if ptd == nil {
		logger.Error("unknown ancestor (writeBlockWithStateSerial)", "num", block.NumberU64(),
			"hash", block.Hash(), "parentHash", block.ParentHash())
		return WriteResult{Status: NonStatTy}, consensus.ErrUnknownAncestor
	}

	if !bc.ShouldTryInserting(block) {
		return WriteResult{Status: NonStatTy}, ErrKnownBlock
	}

	currentBlock := bc.CurrentBlock()
	localTd := bc.GetTd(currentBlock.Hash(), currentBlock.NumberU64())
	externTd := new(big.Int).Add(block.BlockScore(), ptd)

	// Irrelevant of the canonical status, write the block itself to the database
	bc.hc.WriteTd(block.Hash(), block.NumberU64(), externTd)

	// Write other block data.
	bc.writeBlock(block)

	trieWriteStart := time.Now()
	if err := bc.writeStateTrie(block, state); err != nil {
		return WriteResult{Status: NonStatTy}, err
	}
	trieWriteTime := time.Since(trieWriteStart)

	bc.writeReceipts(block.Hash(), block.NumberU64(), receipts)

	// TODO-Klaytn-Issue264 If we are using istanbul BFT, then we always have a canonical chain.
	//         Later we may be able to refine below code.

	// If the total blockscore is higher than our known, add it to the canonical chain
	// Second clause in the if statement reduces the vulnerability to selfish mining.
	// Please refer to http://www.cs.cornell.edu/~ie53/publications/btcProcFC.pdf
	currentBlock = bc.CurrentBlock()
	reorg := isReorganizationRequired(localTd, externTd, currentBlock, block)
	if reorg {
		// Reorganise the chain if the parent is not the head block
		if block.ParentHash() != currentBlock.Hash() {
			if err := bc.reorg(currentBlock, block); err != nil {
				return WriteResult{Status: NonStatTy}, err
			}
		}
		// Write the positional metadata for transaction/receipt lookups and preimages
		if err := bc.writeTxLookupEntries(block); err != nil {
			return WriteResult{Status: NonStatTy}, err
		}
		bc.db.WritePreimages(block.NumberU64(), state.Preimages())
		status = CanonStatTy
	} else {
		status = SideStatTy
	}

	return bc.finalizeWriteBlockWithState(block, status, start, trieWriteTime)
}

// writeBlockWithStateParallel writes the block and all associated state to the database using goroutines.
func (bc *BlockChain) writeBlockWithStateParallel(block *types.Block, receipts []*types.Receipt, state *state.StateDB) (WriteResult, error) {
	start := time.Now()
	bc.wg.Add(1)
	defer bc.wg.Done()

	var status WriteStatus
	// Calculate the total blockscore of the block
	ptd := bc.GetTd(block.ParentHash(), block.NumberU64()-1)
	if ptd == nil {
		logger.Error("unknown ancestor (writeBlockWithStateParallel)", "num", block.NumberU64(),
			"hash", block.Hash(), "parentHash", block.ParentHash())
		return WriteResult{Status: NonStatTy}, consensus.ErrUnknownAncestor
	}

	if !bc.ShouldTryInserting(block) {
		return WriteResult{Status: NonStatTy}, ErrKnownBlock
	}

	currentBlock := bc.CurrentBlock()
	localTd := bc.GetTd(currentBlock.Hash(), currentBlock.NumberU64())
	externTd := new(big.Int).Add(block.BlockScore(), ptd)

	parallelDBWriteWG := sync.WaitGroup{}
	parallelDBWriteErrCh := make(chan error, 2)
	// Irrelevant of the canonical status, write the block itself to the database
	// TODO-Klaytn-Storage Implementing worker pool pattern instead of generating goroutines every time.
	parallelDBWriteWG.Add(4)
	go func() {
		defer parallelDBWriteWG.Done()
		bc.hc.WriteTd(block.Hash(), block.NumberU64(), externTd)
	}()

	// Write other block data.
	go func() {
		defer parallelDBWriteWG.Done()
		bc.writeBlock(block)
	}()

	var trieWriteTime time.Duration
	trieWriteStart := time.Now()
	go func() {
		defer parallelDBWriteWG.Done()
		if err := bc.writeStateTrie(block, state); err != nil {
			parallelDBWriteErrCh <- err
		}
		trieWriteTime = time.Since(trieWriteStart)
	}()

	go func() {
		defer parallelDBWriteWG.Done()
		bc.writeReceipts(block.Hash(), block.NumberU64(), receipts)
	}()

	// Wait until all writing goroutines are terminated.
	parallelDBWriteWG.Wait()
	select {
	case err := <-parallelDBWriteErrCh:
		return WriteResult{Status: NonStatTy}, err
	default:
	}

	// TODO-Klaytn-Issue264 If we are using istanbul BFT, then we always have a canonical chain.
	//         Later we may be able to refine below code.

	// If the total blockscore is higher than our known, add it to the canonical chain
	// Second clause in the if statement reduces the vulnerability to selfish mining.
	// Please refer to http://www.cs.cornell.edu/~ie53/publications/btcProcFC.pdf
	currentBlock = bc.CurrentBlock()
	reorg := isReorganizationRequired(localTd, externTd, currentBlock, block)
	if reorg {
		// Reorganise the chain if the parent is not the head block
		if block.ParentHash() != currentBlock.Hash() {
			if err := bc.reorg(currentBlock, block); err != nil {
				return WriteResult{Status: NonStatTy}, err
			}
		}

		parallelDBWriteWG.Add(2)

		go func() {
			defer parallelDBWriteWG.Done()
			// Write the positional metadata for transaction/receipt lookups
			if err := bc.writeTxLookupEntries(block); err != nil {
				parallelDBWriteErrCh <- err
			}
		}()

		go func() {
			defer parallelDBWriteWG.Done()
			bc.db.WritePreimages(block.NumberU64(), state.Preimages())
		}()

		// Wait until all writing goroutines are terminated.
		parallelDBWriteWG.Wait()

		status = CanonStatTy
	} else {
		status = SideStatTy
	}

	select {
	case err := <-parallelDBWriteErrCh:
		return WriteResult{Status: NonStatTy}, err
	default:
	}

	return bc.finalizeWriteBlockWithState(block, status, start, trieWriteTime)
}

// finalizeWriteBlockWithState updates metrics and inserts block when status is CanonStatTy.
func (bc *BlockChain) finalizeWriteBlockWithState(block *types.Block, status WriteStatus, startTime time.Time, trieWriteTime time.Duration) (WriteResult, error) {
	// Set new head.
	if status == CanonStatTy {
		bc.insert(block)
		headBlockNumberGauge.Update(block.Number().Int64())
		blockTxCountsGauge.Update(int64(block.Transactions().Len()))
		blockTxCountsCounter.Inc(int64(block.Transactions().Len()))
	}
	bc.futureBlocks.Remove(block.Hash())
	return WriteResult{status, time.Since(startTime), trieWriteTime}, nil
}

func (bc *BlockChain) writeTxLookupEntries(block *types.Block) error {
	return bc.db.WriteAndCacheTxLookupEntries(block)
}

// GetTxAndLookupInfoInCache retrieves a tx and lookup info for a given transaction hash in cache.
func (bc *BlockChain) GetTxAndLookupInfoInCache(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	return bc.db.ReadTxAndLookupInfoInCache(hash)
}

// GetBlockReceiptsInCache returns receipt of txHash in cache.
func (bc *BlockChain) GetBlockReceiptsInCache(blockHash common.Hash) types.Receipts {
	return bc.db.ReadBlockReceiptsInCache(blockHash)
}

// GetTxReceiptInCache returns receipt of txHash in cache.
func (bc *BlockChain) GetTxReceiptInCache(txHash common.Hash) *types.Receipt {
	return bc.db.ReadTxReceiptInCache(txHash)
}

// InsertChain attempts to insert the given batch of blocks in to the canonical
// chain or, otherwise, create a fork. If an error is returned it will return
// the index number of the failing block as well an error describing what went
// wrong.
//
// After insertion is done, all accumulated events will be fired.
func (bc *BlockChain) InsertChain(chain types.Blocks) (int, error) {
	n, events, logs, err := bc.insertChain(chain)
	bc.PostChainEvents(events, logs)
	return n, err
}

// insertChain will execute the actual chain insertion and event aggregation. The
// only reason this method exists as a separate one is to make locking cleaner
// with deferred statements.
func (bc *BlockChain) insertChain(chain types.Blocks) (int, []interface{}, []*types.Log, error) {
	// Sanity check that we have something meaningful to import
	if len(chain) == 0 {
		return 0, nil, nil, nil
	}
	// Do a sanity check that the provided chain is actually ordered and linked
	for i := 1; i < len(chain); i++ {
		if chain[i].NumberU64() != chain[i-1].NumberU64()+1 || chain[i].ParentHash() != chain[i-1].Hash() {
			// Chain broke ancestry, log a messge (programming error) and skip insertion
			logger.Error("Non contiguous block insert", "number", chain[i].Number(), "hash", chain[i].Hash(),
				"parent", chain[i].ParentHash(), "prevnumber", chain[i-1].Number(), "prevhash", chain[i-1].Hash())

			return 0, nil, nil, fmt.Errorf("non contiguous insert: item %d is #%d [%x…], item %d is #%d [%x…] (parent [%x…])", i-1, chain[i-1].NumberU64(),
				chain[i-1].Hash().Bytes()[:4], i, chain[i].NumberU64(), chain[i].Hash().Bytes()[:4], chain[i].ParentHash().Bytes()[:4])
		}
	}

	// Pre-checks passed, start the full block imports
	bc.wg.Add(1)
	defer bc.wg.Done()

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// A queued approach to delivering events. This is generally
	// faster than direct delivery and requires much less mutex
	// acquiring.
	var (
		stats         = insertStats{startTime: mclock.Now()}
		events        = make([]interface{}, 0, len(chain))
		lastCanon     *types.Block
		coalescedLogs []*types.Log
	)
	// Start the parallel header verifier
	headers := make([]*types.Header, len(chain))
	seals := make([]bool, len(chain))

	for i, block := range chain {
		headers[i] = block.Header()
		seals[i] = true
	}

	var (
		abort   chan<- struct{}
		results <-chan error
	)
	if bc.engine.CanVerifyHeadersConcurrently() {
		abort, results = bc.engine.VerifyHeaders(bc, headers, seals)
	} else {
		abort, results = bc.engine.PreprocessHeaderVerification(headers)
	}
	defer close(abort)

	// Start a parallel signature recovery (signer will fluke on fork transition, minimal perf loss)
	senderCacher.recoverFromBlocks(types.MakeSigner(bc.chainConfig, chain[0].Number()), chain)

	// Iterate over the blocks and insert when the verifier permits
	for i, block := range chain {
		// If the chain is terminating, stop processing blocks
		if atomic.LoadInt32(&bc.procInterrupt) == 1 {
			logger.Debug("Premature abort during blocks processing")
			break
		}
		// Create a new trie using the parent block and report an
		// error if it fails.
		var parent *types.Block
		if i == 0 {
			parent = bc.GetBlock(block.ParentHash(), block.NumberU64()-1)
		} else {
			parent = chain[i-1]
		}

		// If we have a followup block, run that against the current state to pre-cache
		// transactions and probabilistically some of the account/storage trie nodes.
		var followupInterrupt uint32

		if bc.cacheConfig.TrieNodeCacheConfig.NumFetcherPrefetchWorker > 0 && parent != nil {
			var snaps *snapshot.Tree
			if bc.cacheConfig.TrieNodeCacheConfig.UseSnapshotForPrefetch {
				snaps = bc.snaps
			}

			// Tx prefetcher is enabled for all cases (both single and multiple block insertion).
			for ti := range block.Transactions() {
				select {
				case bc.prefetchTxCh <- prefetchTx{ti, block, &followupInterrupt}:
				default:
				}
			}
			if i < len(chain)-1 {
				// current block is not the last one, so prefetch the right next block
				followup := chain[i+1]
				go func(start time.Time) {
					defer func() {
						if err := recover(); err != nil {
							logger.Error("Got panic and recovered from prefetcher", "err", err)
						}
					}()

					throwaway, err := state.NewForPrefetching(parent.Root(), bc.stateCache, snaps)
					if throwaway == nil || err != nil {
						logger.Warn("failed to get StateDB for prefetcher", "err", err,
							"parentBlockNum", parent.NumberU64(), "currBlockNum", bc.CurrentBlock().NumberU64())
						return
					}

					vmCfg := bc.vmConfig
					vmCfg.Prefetching = true
					bc.prefetcher.Prefetch(followup, throwaway, vmCfg, &followupInterrupt)

					blockPrefetchExecuteTimer.Update(time.Since(start))
					if atomic.LoadUint32(&followupInterrupt) == 1 {
						blockPrefetchInterruptMeter.Mark(1)
					}
				}(time.Now())
			}
		}
		// If the header is a banned one, straight out abort
		if BadHashes[block.Hash()] {
			bc.reportBlock(block, nil, ErrBlacklistedHash)
			return i, events, coalescedLogs, ErrBlacklistedHash
		}
		// Wait for the block's verification to complete
		bstart := time.Now()

		err := <-results
		if !bc.engine.CanVerifyHeadersConcurrently() && err == nil {
			err = bc.engine.VerifyHeader(bc, block.Header(), true)
		}

		if err == nil {
			err = bc.validator.ValidateBody(block)
		}

		switch {
		case err == ErrKnownBlock:
			// Block and state both already known. However if the current block is below
			// this number we did a rollback and we should reimport it nonetheless.
			if bc.CurrentBlock().NumberU64() >= block.NumberU64() {
				stats.ignored++
				continue
			}

		case err == consensus.ErrFutureBlock:
			// Allow up to MaxFuture second in the future blocks. If this limit is exceeded
			// the chain is discarded and processed at a later time if given.
			max := big.NewInt(time.Now().Unix() + maxTimeFutureBlocks)
			if block.Time().Cmp(max) > 0 {
				return i, events, coalescedLogs, fmt.Errorf("future block: %v > %v", block.Time(), max)
			}
			bc.futureBlocks.Add(block.Hash(), block)
			stats.queued++
			continue

		case err == consensus.ErrUnknownAncestor && bc.futureBlocks.Contains(block.ParentHash()):
			bc.futureBlocks.Add(block.Hash(), block)
			stats.queued++
			continue

		case err == consensus.ErrPrunedAncestor:
			// Block competing with the canonical chain, store in the db, but don't process
			// until the competitor TD goes above the canonical TD
			currentBlock := bc.CurrentBlock()
			localTd := bc.GetTd(currentBlock.Hash(), currentBlock.NumberU64())
			externTd := new(big.Int).Add(bc.GetTd(block.ParentHash(), block.NumberU64()-1), block.BlockScore())
			if localTd.Cmp(externTd) > 0 {
				bc.WriteBlockWithoutState(block, externTd)
				continue
			}
			// Competitor chain beat canonical, gather all blocks from the common ancestor
			var winner []*types.Block

			parent := bc.GetBlock(block.ParentHash(), block.NumberU64()-1)
			for !bc.HasState(parent.Root()) {
				winner = append(winner, parent)
				parent = bc.GetBlock(parent.ParentHash(), parent.NumberU64()-1)
			}
			for j := 0; j < len(winner)/2; j++ {
				winner[j], winner[len(winner)-1-j] = winner[len(winner)-1-j], winner[j]
			}
			// Import all the pruned blocks to make the state available
			bc.mu.Unlock()
			_, evs, logs, err := bc.insertChain(winner)
			bc.mu.Lock()
			events, coalescedLogs = evs, logs

			if err != nil {
				return i, events, coalescedLogs, err
			}

		case err != nil:
			bc.futureBlocks.Remove(block.Hash())
			bc.reportBlock(block, nil, err)
			return i, events, coalescedLogs, err
		}

		stateDB, err := bc.StateAt(parent.Root())
		if err != nil {
			return i, events, coalescedLogs, err
		}

		// Process block using the parent state as reference point.
		receipts, logs, usedGas, internalTxTraces, procStats, err := bc.processor.Process(block, stateDB, bc.vmConfig)
		if err != nil {
			bc.reportBlock(block, receipts, err)
			atomic.StoreUint32(&followupInterrupt, 1)
			return i, events, coalescedLogs, err
		}

		// Validate the state using the default validator
		err = bc.validator.ValidateState(block, parent, stateDB, receipts, usedGas)
		if err != nil {
			bc.reportBlock(block, receipts, err)
			atomic.StoreUint32(&followupInterrupt, 1)
			return i, events, coalescedLogs, err
		}
		afterValidate := time.Now()

		// Write the block to the chain and get the writeResult.
		writeResult, err := bc.writeBlockWithState(block, receipts, stateDB)
		if err != nil {
			atomic.StoreUint32(&followupInterrupt, 1)
			if err == ErrKnownBlock {
				logger.Debug("Tried to insert already known block", "num", block.NumberU64(), "hash", block.Hash().String())
				continue
			}
			return i, events, coalescedLogs, err
		}
		atomic.StoreUint32(&followupInterrupt, 1)

		// Update to-address based spam throttler when spamThrottler is enabled and a single block is fetched.
		spamThrottler := GetSpamThrottler()
		if spamThrottler != nil && len(chain) == 1 {
			spamThrottler.updateThrottlerState(block.Transactions(), receipts)
		}

		// Update the metrics subsystem with all the measurements
		accountReadTimer.Update(stateDB.AccountReads)
		accountHashTimer.Update(stateDB.AccountHashes)
		accountUpdateTimer.Update(stateDB.AccountUpdates)
		accountCommitTimer.Update(stateDB.AccountCommits)

		storageReadTimer.Update(stateDB.StorageReads)
		storageHashTimer.Update(stateDB.StorageHashes)
		storageUpdateTimer.Update(stateDB.StorageUpdates)
		storageCommitTimer.Update(stateDB.StorageCommits)

		snapshotAccountReadTimer.Update(stateDB.SnapshotAccountReads)
		snapshotStorageReadTimer.Update(stateDB.SnapshotStorageReads)
		snapshotCommitTimer.Update(stateDB.SnapshotCommits)

		trieAccess := stateDB.AccountReads + stateDB.AccountHashes + stateDB.AccountUpdates + stateDB.AccountCommits
		trieAccess += stateDB.StorageReads + stateDB.StorageHashes + stateDB.StorageUpdates + stateDB.StorageCommits

		blockAgeTimer.Update(time.Since(time.Unix(int64(block.Time().Uint64()), 0)))

		switch writeResult.Status {
		case CanonStatTy:
			processTxsTime := common.PrettyDuration(procStats.AfterApplyTxs.Sub(procStats.BeforeApplyTxs))
			processFinalizeTime := common.PrettyDuration(procStats.AfterFinalize.Sub(procStats.AfterApplyTxs))
			validateTime := common.PrettyDuration(afterValidate.Sub(procStats.AfterFinalize))
			totalTime := common.PrettyDuration(time.Since(bstart))
			logger.Info("Inserted a new block", "number", block.Number(), "hash", block.Hash(),
				"txs", len(block.Transactions()), "gas", block.GasUsed(), "elapsed", totalTime,
				"processTxs", processTxsTime, "finalize", processFinalizeTime, "validateState", validateTime,
				"totalWrite", writeResult.TotalWriteTime, "trieWrite", writeResult.TrieWriteTime)

			blockProcessTimer.Update(time.Duration(processTxsTime))
			blockExecutionTimer.Update(time.Duration(processTxsTime) - trieAccess)
			blockFinalizeTimer.Update(time.Duration(processFinalizeTime))
			blockValidateTimer.Update(time.Duration(validateTime))
			blockInsertTimer.Update(time.Duration(totalTime))

			coalescedLogs = append(coalescedLogs, logs...)
			events = append(events, ChainEvent{
				Block:            block,
				Hash:             block.Hash(),
				Logs:             logs,
				Receipts:         receipts,
				InternalTxTraces: internalTxTraces,
			})
			lastCanon = block

		case SideStatTy:
			logger.Debug("Inserted forked block", "number", block.Number(), "hash", block.Hash(), "diff", block.BlockScore(), "elapsed",
				common.PrettyDuration(time.Since(bstart)), "txs", len(block.Transactions()), "gas", block.GasUsed())

			events = append(events, ChainSideEvent{block})
		}
		stats.processed++
		stats.usedGas += usedGas

		cache, _ := bc.stateCache.TrieDB().Size()
		stats.report(chain, i, cache)

		// update governance CurrentSet if it is at an epoch block
		if bc.engine.CreateSnapshot(bc, block.NumberU64(), block.Hash(), nil) != nil {
			return i, events, coalescedLogs, err
		}
	}
	// Append a single chain head event if we've progressed the chain
	if lastCanon != nil && bc.CurrentBlock().Hash() == lastCanon.Hash() {
		events = append(events, ChainHeadEvent{lastCanon})
	}
	return 0, events, coalescedLogs, nil
}

// BlockSubscriptionLoop subscribes blocks from a redis server and processes them.
// This method is only for KES nodes.
func (bc *BlockChain) BlockSubscriptionLoop(pool *TxPool) {
	var ch <-chan *redis.Message
	logger.Info("subscribe blocks from redis cache")

	pubSub, ok := bc.stateCache.TrieDB().TrieNodeCache().(statedb.BlockPubSub)
	if !ok || pubSub == nil {
		logger.Crit("invalid block pub/sub configure", "trieNodeCacheConfig",
			bc.stateCache.TrieDB().GetTrieNodeCacheConfig())
	}

	ch = pubSub.SubscribeBlockCh()
	if ch == nil {
		logger.Crit("failed to create redis subscription channel")
	}

	for msg := range ch {
		logger.Debug("msg from redis subscription channel", "msg", msg.Payload)

		blockRlp, err := hexutil.Decode(msg.Payload)
		if err != nil {
			logger.Error("failed to decode redis subscription msg", "msg", msg.Payload)
			continue
		}

		block := &types.Block{}
		if err := rlp.DecodeBytes(blockRlp, block); err != nil {
			logger.Error("failed to rlp decode block", "msg", msg.Payload, "block", string(blockRlp))
			continue
		}

		oldHead := bc.CurrentHeader()
		bc.replaceCurrentBlock(block)
		pool.lockedReset(oldHead, bc.CurrentHeader())

		// just in case the block number jumps up more than one, iterates all missed blocks
		for blockNum := oldHead.Number.Uint64() + 1; blockNum < block.Number().Uint64(); blockNum++ {
			retrievedBlock := bc.GetBlockByNumber(blockNum)
			bc.sendKESSubscriptionData(retrievedBlock)
		}
		bc.sendKESSubscriptionData(block)
	}

	logger.Info("closed the block subscription loop")
}

// sendKESSubscriptionData sends data to chainFeed and logsFeed.
// ChainEvent containing only Block and Hash is sent to chainFeed.
// []*types.Log containing entire logs of a block is set to logsFeed.
// The logs are expected to be delivered from remote cache.
// If it failed to read log data from remote cache, it will read the data from database.
// This method is only for KES nodes.
func (bc *BlockChain) sendKESSubscriptionData(block *types.Block) {
	bc.chainFeed.Send(ChainEvent{
		Block: block,
		Hash:  block.Hash(),
		// TODO-Klaytn-KES: fill the following data if needed
		Receipts:         types.Receipts{},
		Logs:             []*types.Log{},
		InternalTxTraces: []*vm.InternalTxTrace{},
	})

	// TODO-Klaytn-KES: refine this not to use trieNodeCache
	logKey := append(kesCachePrefixBlockLogs, block.Number().Bytes()...)
	encodedLogs := bc.stateCache.TrieDB().TrieNodeCache().Get(logKey)
	if encodedLogs == nil {
		logger.Warn("cannot get a block log from the remote cache", "blockNum", block.NumberU64())

		// read log data from database and send it
		logsList := bc.GetLogsByHash(block.Header().Hash())
		var logs []*types.Log
		for _, list := range logsList {
			logs = append(logs, list...)
		}
		bc.logsFeed.Send(logs)
		return
	}

	entireLogs := []*types.LogForStorage{}
	if err := rlp.DecodeBytes(encodedLogs, &entireLogs); err != nil {
		logger.Warn("failed to decode a block log", "blockNum", block.NumberU64(), "err", err)

		// read log data from database and send it
		logsList := bc.GetLogsByHash(block.Header().Hash())
		var logs []*types.Log
		for _, list := range logsList {
			logs = append(logs, list...)
		}
		bc.logsFeed.Send(logs)
		return
	}

	// convert LogForStorage to Log
	logs := make([]*types.Log, len(entireLogs))
	for i, log := range entireLogs {
		logs[i] = (*types.Log)(log)
	}
	bc.logsFeed.Send(logs)
}

// CloseBlockSubscriptionLoop closes BlockSubscriptionLoop.
func (bc *BlockChain) CloseBlockSubscriptionLoop() {
	pubSub, ok := bc.stateCache.TrieDB().TrieNodeCache().(statedb.BlockPubSub)
	if ok {
		if err := pubSub.UnsubscribeBlock(); err != nil {
			logger.Error("failed to unsubscribe blocks", "err", err, "trieNodeCacheConfig",
				bc.stateCache.TrieDB().GetTrieNodeCacheConfig())
		}
	}
}

// replaceCurrentBlock replaces bc.currentBlock to the given block.
func (bc *BlockChain) replaceCurrentBlock(latestBlock *types.Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if latestBlock == nil {
		logger.Error("no latest block")
		return
	}

	// Don't update current block if the latest block is not newer than current block.
	currentBlock := bc.CurrentBlock()
	if currentBlock.NumberU64() >= latestBlock.NumberU64() {
		logger.Debug("ignore an old block", "currentBlockNumber", currentBlock.NumberU64(), "oldBlockNumber",
			latestBlock.NumberU64())
		return
	}

	// Insert a new block and update metrics
	bc.insert(latestBlock)
	bc.hc.SetCurrentHeader(latestBlock.Header())

	headBlockNumberGauge.Update(latestBlock.Number().Int64())
	blockTxCountsGauge.Update(int64(latestBlock.Transactions().Len()))
	blockTxCountsCounter.Inc(int64(latestBlock.Transactions().Len()))
	bc.stateCache.TrieDB().UpdateMetricNodes()

	logger.Info("Replaced the current block",
		"blkNum", latestBlock.NumberU64(), "blkHash", latestBlock.Hash().String())
}

// insertStats tracks and reports on block insertion.
type insertStats struct {
	queued, processed, ignored int
	usedGas                    uint64
	lastIndex                  int
	startTime                  mclock.AbsTime
}

// report prints statistics if some number of blocks have been processed
// or more than a few seconds have passed since the last message.
func (st *insertStats) report(chain []*types.Block, index int, cache common.StorageSize) {
	// report will leave a log only if inserting two or more blocks at once
	if len(chain) <= 1 {
		return
	}
	// Fetch the timings for the batch
	var (
		now     = mclock.Now()
		elapsed = time.Duration(now) - time.Duration(st.startTime)
	)
	// If we're at the last block of the batch or report period reached, log
	if index == len(chain)-1 || elapsed >= log.StatsReportLimit {
		var (
			end = chain[index]
			txs = countTransactions(chain[st.lastIndex : index+1])
		)
		context := []interface{}{
			"number", end.Number(), "hash", end.Hash(), "blocks", st.processed, "txs", txs, "elapsed", common.PrettyDuration(elapsed),
			"trieDBSize", cache, "mgas", float64(st.usedGas) / 1000000, "mgasps", float64(st.usedGas) * 1000 / float64(elapsed),
		}

		timestamp := time.Unix(int64(end.Time().Uint64()), 0)
		context = append(context, []interface{}{"age", common.PrettyAge(timestamp)}...)

		if st.queued > 0 {
			context = append(context, []interface{}{"queued", st.queued}...)
		}
		if st.ignored > 0 {
			context = append(context, []interface{}{"ignored", st.ignored}...)
		}
		logger.Info("Imported new chain segment", context...)

		*st = insertStats{startTime: now, lastIndex: index + 1}
	}
}

func countTransactions(chain []*types.Block) (c int) {
	for _, b := range chain {
		c += len(b.Transactions())
	}
	return c
}

// reorgs takes two blocks, an old chain and a new chain and will reconstruct the blocks and inserts them
// to be part of the new canonical chain and accumulates potential missing transactions and post an
// event about them
func (bc *BlockChain) reorg(oldBlock, newBlock *types.Block) error {
	var (
		newChain    types.Blocks
		oldChain    types.Blocks
		commonBlock *types.Block
		deletedTxs  types.Transactions
		deletedLogs []*types.Log
		// collectLogs collects the logs that were generated during the
		// processing of the block that corresponds with the given hash.
		// These logs are later announced as deleted.
		collectLogs = func(hash common.Hash) {
			// Coalesce logs and set 'Removed'.
			number := bc.GetBlockNumber(hash)
			if number == nil {
				return
			}
			receipts := bc.db.ReadReceipts(hash, *number)
			for _, receipt := range receipts {
				for _, log := range receipt.Logs {
					del := *log
					del.Removed = true
					deletedLogs = append(deletedLogs, &del)
				}
			}
		}
	)

	// first reduce whoever is higher bound
	if oldBlock.NumberU64() > newBlock.NumberU64() {
		// reduce old chain
		for ; oldBlock != nil && oldBlock.NumberU64() != newBlock.NumberU64(); oldBlock = bc.GetBlock(oldBlock.ParentHash(), oldBlock.NumberU64()-1) {
			oldChain = append(oldChain, oldBlock)
			deletedTxs = append(deletedTxs, oldBlock.Transactions()...)

			collectLogs(oldBlock.Hash())
		}
	} else {
		// reduce new chain and append new chain blocks for inserting later on
		for ; newBlock != nil && newBlock.NumberU64() != oldBlock.NumberU64(); newBlock = bc.GetBlock(newBlock.ParentHash(), newBlock.NumberU64()-1) {
			newChain = append(newChain, newBlock)
		}
	}
	if oldBlock == nil {
		return fmt.Errorf("Invalid old chain")
	}
	if newBlock == nil {
		return fmt.Errorf("Invalid new chain")
	}

	for {
		if oldBlock.Hash() == newBlock.Hash() {
			commonBlock = oldBlock
			break
		}

		oldChain = append(oldChain, oldBlock)
		newChain = append(newChain, newBlock)
		deletedTxs = append(deletedTxs, oldBlock.Transactions()...)
		collectLogs(oldBlock.Hash())

		oldBlock, newBlock = bc.GetBlock(oldBlock.ParentHash(), oldBlock.NumberU64()-1), bc.GetBlock(newBlock.ParentHash(), newBlock.NumberU64()-1)
		if oldBlock == nil {
			return fmt.Errorf("Invalid old chain")
		}
		if newBlock == nil {
			return fmt.Errorf("Invalid new chain")
		}
	}
	// Ensure the user sees large reorgs
	if len(oldChain) > 0 && len(newChain) > 0 {
		logFn := logger.Debug
		if len(oldChain) > 63 {
			logFn = logger.Warn
		}
		logFn("Chain split detected", "number", commonBlock.Number(), "hash", commonBlock.Hash(),
			"drop", len(oldChain), "dropfrom", oldChain[0].Hash(), "add", len(newChain), "addfrom", newChain[0].Hash())
	} else {
		logger.Error("Impossible reorg, please file an issue", "oldnum", oldBlock.Number(), "oldhash", oldBlock.Hash(), "newnum", newBlock.Number(), "newhash", newBlock.Hash())
	}
	// Insert the new chain, taking care of the proper incremental order
	var addedTxs types.Transactions
	for i := len(newChain) - 1; i >= 0; i-- {
		// insert the block in the canonical way, re-writing history
		bc.insert(newChain[i])
		// write lookup entries for hash based transaction/receipt searches
		bc.db.WriteTxLookupEntries(newChain[i])
		addedTxs = append(addedTxs, newChain[i].Transactions()...)
	}
	// calculate the difference between deleted and added transactions
	diff := types.TxDifference(deletedTxs, addedTxs)
	// When transactions get deleted from the database that means the
	// receipts that were created in the fork must also be deleted
	for _, tx := range diff {
		bc.db.DeleteTxLookupEntry(tx.Hash())
	}
	if len(deletedLogs) > 0 {
		go bc.rmLogsFeed.Send(RemovedLogsEvent{deletedLogs})
	}
	if len(oldChain) > 0 {
		go func() {
			for _, block := range oldChain {
				bc.chainSideFeed.Send(ChainSideEvent{Block: block})
			}
		}()
	}

	return nil
}

// PostChainEvents iterates over the events generated by a chain insertion and
// posts them into the event feed.
// TODO: Should not expose PostChainEvents. The chain events should be posted in WriteBlock.
func (bc *BlockChain) PostChainEvents(events []interface{}, logs []*types.Log) {
	// post event logs for further processing
	if logs != nil {
		bc.logsFeed.Send(logs)
	}
	for _, event := range events {
		switch ev := event.(type) {
		case ChainEvent:
			bc.chainFeed.Send(ev)

		case ChainHeadEvent:
			bc.chainHeadFeed.Send(ev)

		case ChainSideEvent:
			bc.chainSideFeed.Send(ev)
		}
	}
}

func (bc *BlockChain) update() {
	futureTimer := time.NewTicker(5 * time.Second)
	defer futureTimer.Stop()
	for {
		select {
		case <-futureTimer.C:
			bc.procFutureBlocks()
		case <-bc.quit:
			return
		}
	}
}

// BadBlockArgs represents the entries in the list returned when bad blocks are queried.
type BadBlockArgs struct {
	Hash  common.Hash  `json:"hash"`
	Block *types.Block `json:"block"`
}

// BadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
func (bc *BlockChain) BadBlocks() ([]BadBlockArgs, error) {
	blocks := make([]BadBlockArgs, 0, bc.badBlocks.Len())
	for _, hash := range bc.badBlocks.Keys() {
		hashKey, ok := hash.(common.CacheKey)
		if !ok {
			logger.Error("invalid key type", "expect", "common.CacheKey", "actual", reflect.TypeOf(hash))
			continue
		}

		if blk, exist := bc.badBlocks.Peek(hashKey); exist {
			cacheGetBadBlockHitMeter.Mark(1)
			block := blk.(*types.Block)
			blocks = append(blocks, BadBlockArgs{block.Hash(), block})
		} else {
			cacheGetBadBlockMissMeter.Mark(1)
		}
	}
	return blocks, nil
}

// istanbul BFT
func (bc *BlockChain) HasBadBlock(hash common.Hash) bool {
	return bc.badBlocks.Contains(hash)
}

// addBadBlock adds a bad block to the bad-block LRU cache
func (bc *BlockChain) addBadBlock(block *types.Block) {
	bc.badBlocks.Add(block.Header().Hash(), block)
}

// reportBlock logs a bad block error.
func (bc *BlockChain) reportBlock(block *types.Block, receipts types.Receipts, err error) {
	badBlockCounter.Inc(1)
	bc.addBadBlock(block)

	var receiptString string
	for i, receipt := range receipts {
		receiptString += fmt.Sprintf("\t %d: tx: %v status: %v gas: %v contract: %v bloom: %x logs: %v\n",
			i, receipt.TxHash.Hex(), receipt.Status, receipt.GasUsed, receipt.ContractAddress.Hex(),
			receipt.Bloom, receipt.Logs)
	}
	logger.Error(fmt.Sprintf(`
########## BAD BLOCK #########
Chain config: %v

Number: %v
Hash: 0x%x
%v

Error: %v
##############################
`, bc.chainConfig, block.Number(), block.Hash(), receiptString, err))
}

// InsertHeaderChain attempts to insert the given header chain in to the local
// chain, possibly creating a reorg. If an error is returned, it will return the
// index number of the failing header as well an error describing what went wrong.
//
// The verify parameter can be used to fine tune whether nonce verification
// should be done or not. The reason behind the optional check is because some
// of the header retrieval mechanisms already need to verify nonces, as well as
// because nonces can be verified sparsely, not needing to check each.
func (bc *BlockChain) InsertHeaderChain(chain []*types.Header, checkFreq int) (int, error) {
	start := time.Now()
	if i, err := bc.hc.ValidateHeaderChain(chain, checkFreq); err != nil {
		return i, err
	}

	// Make sure only one thread manipulates the chain at once
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.wg.Add(1)
	defer bc.wg.Done()

	whFunc := func(header *types.Header) error {

		_, err := bc.hc.WriteHeader(header)
		return err
	}

	return bc.hc.InsertHeaderChain(chain, whFunc, start)
}

// CurrentHeader retrieves the current head header of the canonical chain. The
// header is retrieved from the HeaderChain's internal cache.
func (bc *BlockChain) CurrentHeader() *types.Header {
	return bc.hc.CurrentHeader()
}

// GetTd retrieves a block's total blockscore in the canonical chain from the
// database by hash and number, caching it if found.
func (bc *BlockChain) GetTd(hash common.Hash, number uint64) *big.Int {
	return bc.hc.GetTd(hash, number)
}

// GetTdByHash retrieves a block's total blockscore in the canonical chain from the
// database by hash, caching it if found.
func (bc *BlockChain) GetTdByHash(hash common.Hash) *big.Int {
	return bc.hc.GetTdByHash(hash)
}

// GetHeader retrieves a block header from the database by hash and number,
// caching it if found.
func (bc *BlockChain) GetHeader(hash common.Hash, number uint64) *types.Header {
	return bc.hc.GetHeader(hash, number)
}

// GetHeaderByHash retrieves a block header from the database by hash, caching it if
// found.
func (bc *BlockChain) GetHeaderByHash(hash common.Hash) *types.Header {
	return bc.hc.GetHeaderByHash(hash)
}

// HasHeader checks if a block header is present in the database or not, caching
// it if present.
func (bc *BlockChain) HasHeader(hash common.Hash, number uint64) bool {
	return bc.hc.HasHeader(hash, number)
}

// GetBlockHashesFromHash retrieves a number of block hashes starting at a given
// hash, fetching towards the genesis block.
func (bc *BlockChain) GetBlockHashesFromHash(hash common.Hash, max uint64) []common.Hash {
	return bc.hc.GetBlockHashesFromHash(hash, max)
}

// GetHeaderByNumber retrieves a block header from the database by number,
// caching it (associated with its hash) if found.
func (bc *BlockChain) GetHeaderByNumber(number uint64) *types.Header {
	return bc.hc.GetHeaderByNumber(number)
}

// Config retrieves the blockchain's chain configuration.
func (bc *BlockChain) Config() *params.ChainConfig { return bc.chainConfig }

// Engine retrieves the blockchain's consensus engine.
func (bc *BlockChain) Engine() consensus.Engine { return bc.engine }

// SubscribeRemovedLogsEvent registers a subscription of RemovedLogsEvent.
func (bc *BlockChain) SubscribeRemovedLogsEvent(ch chan<- RemovedLogsEvent) event.Subscription {
	return bc.scope.Track(bc.rmLogsFeed.Subscribe(ch))
}

// SubscribeChainEvent registers a subscription of ChainEvent.
func (bc *BlockChain) SubscribeChainEvent(ch chan<- ChainEvent) event.Subscription {
	return bc.scope.Track(bc.chainFeed.Subscribe(ch))
}

// SubscribeChainHeadEvent registers a subscription of ChainHeadEvent.
func (bc *BlockChain) SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription {
	return bc.scope.Track(bc.chainHeadFeed.Subscribe(ch))
}

// SubscribeChainSideEvent registers a subscription of ChainSideEvent.
func (bc *BlockChain) SubscribeChainSideEvent(ch chan<- ChainSideEvent) event.Subscription {
	return bc.scope.Track(bc.chainSideFeed.Subscribe(ch))
}

// SubscribeLogsEvent registers a subscription of []*types.Log.
func (bc *BlockChain) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return bc.scope.Track(bc.logsFeed.Subscribe(ch))
}

// isArchiveMode returns whether current blockchain is in archiving mode or not.
// cacheConfig.ArchiveMode means trie caching is disabled.
func (bc *BlockChain) isArchiveMode() bool {
	return bc.cacheConfig.ArchiveMode
}

// IsParallelDBWrite returns if parallel write is enabled or not.
// If enabled, data written in WriteBlockWithState is being written in parallel manner.
func (bc *BlockChain) IsParallelDBWrite() bool {
	return bc.parallelDBWrite
}

// IsSenderTxHashIndexingEnabled returns if storing senderTxHash to txHash mapping information
// is enabled or not.
func (bc *BlockChain) IsSenderTxHashIndexingEnabled() bool {
	return bc.cacheConfig.SenderTxHashIndexing
}

func (bc *BlockChain) SaveTrieNodeCacheToDisk() error {
	if err := bc.stateCache.TrieDB().CanSaveTrieNodeCacheToFile(); err != nil {
		return err
	}
	go bc.stateCache.TrieDB().SaveTrieNodeCacheToFile(bc.cacheConfig.TrieNodeCacheConfig.FastCacheFileDir, runtime.NumCPU()/2)
	return nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func (bc *BlockChain) ApplyTransaction(chainConfig *params.ChainConfig, author *common.Address, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, vmConfig *vm.Config) (*types.Receipt, uint64, *vm.InternalTxTrace, error) {

	// TODO-Klaytn We reject transactions with unexpected gasPrice and do not put the transaction into TxPool.
	//         And we run transactions regardless of gasPrice if we push transactions in the TxPool.
	/*
		// istanbul BFT
		if tx.GasPrice() != nil && tx.GasPrice().Cmp(common.Big0) > 0 {
			return nil, uint64(0), ErrInvalidGasPrice
		}
	*/

	blockNumber := header.Number.Uint64()

	// validation for each transaction before execution
	if err := tx.Validate(statedb, blockNumber); err != nil {
		return nil, 0, nil, err
	}

	msg, err := tx.AsMessageWithAccountKeyPicker(types.MakeSigner(chainConfig, header.Number), statedb, blockNumber)
	if err != nil {
		return nil, 0, nil, err
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, chainConfig, vmConfig)
	// Apply the transaction to the current state (included in the env)
	_, gas, kerr := ApplyMessage(vmenv, msg)
	err = kerr.ErrTxInvalid
	if err != nil {
		return nil, 0, nil, err
	}

	var internalTrace *vm.InternalTxTrace
	if vmConfig.EnableInternalTxTracing {
		internalTrace, err = GetInternalTxTrace(vmConfig.Tracer)
		if err != nil {
			logger.Error("failed to get tracing result from a transaction", "txHash", tx.Hash().String(), "err", err)
			return nil, 0, nil, err
		}
	}
	// Update the state with pending changes
	statedb.Finalise(true, false)
	*usedGas += gas

	receipt := types.NewReceipt(kerr.Status, tx.Hash(), gas)
	// if the transaction created a contract, store the creation address in the receipt.
	msg.FillContractAddress(vmenv.Context.Origin, receipt)
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, internalTrace, err
}

func GetInternalTxTrace(tracer vm.Tracer) (*vm.InternalTxTrace, error) {
	var (
		internalTxTrace *vm.InternalTxTrace
		err             error
	)
	switch tracer := tracer.(type) {
	case *vm.InternalTxTracer:
		internalTxTrace, err = tracer.GetResult()
		if err != nil {
			return nil, err
		}
	default:
		logger.Error("To trace internal transactions, VM tracer type should be vm.InternalTxTracer", "actualType", reflect.TypeOf(tracer).String())
		return nil, ErrInvalidTracer
	}
	return internalTxTrace, nil
}

// CheckBlockChainVersion checks the version of the current database and upgrade if possible.
func CheckBlockChainVersion(chainDB database.DBManager) error {
	bcVersion := chainDB.ReadDatabaseVersion()
	if bcVersion != nil && *bcVersion > BlockChainVersion {
		return fmt.Errorf("database version is v%d, Klaytn %s only supports v%d", *bcVersion, params.Version, BlockChainVersion)
	} else if bcVersion == nil || *bcVersion < BlockChainVersion {
		bcVersionStr := "N/A"
		if bcVersion != nil {
			bcVersionStr = strconv.Itoa(int(*bcVersion))
		}
		logger.Warn("Upgrade database version", "from", bcVersionStr, "to", BlockChainVersion)
		chainDB.WriteDatabaseVersion(BlockChainVersion)
	}
	return nil
}
