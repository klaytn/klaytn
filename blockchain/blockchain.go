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
	"github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/mclock"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/rcrowley/go-metrics"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
	"io"
	"math/big"
	mrand "math/rand"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	blockInsertTimeGauge = metrics.NewRegisteredGauge("chain/inserts", nil)
	ErrNoGenesis         = errors.New("Genesis not found in chain")
	ErrNotExistNode      = errors.New("the node does not exist in cached node")
	logger               = log.NewModuleLogger(log.Blockchain)
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
	triesInMemory = 4
	// BlockChainVersion ensures that an incompatible database forces a resync from scratch.
	BlockChainVersion    = 3
	DefaultBlockInterval = 128
)

// CacheConfig contains the configuration values for the 1) stateDB caching and
// 2) trie caching/pruning resident in a blockchain.
type CacheConfig struct {
	// TODO-Klaytn-Issue1666 Need to check the benefit of trie caching.
	StateDBCaching       bool // Enables caching of state objects in stateDB.
	TxPoolStateCache     bool // Enables caching of nonce and balance for txpool.
	ArchiveMode          bool // If true, state trie is not pruned and always written to database.
	CacheSize            int  // Size of in-memory cache of a trie (MiB) to flush matured singleton trie nodes to disk
	BlockInterval        uint // Block interval to flush the trie. Each interval state trie will be flushed into disk.
	TrieCacheLimit       int  // Memory allowance (MB) to use for caching trie nodes in memory
	SenderTxHashIndexing bool // Enables saving senderTxHash to txHash mapping information to database and cache.
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

	db database.DBManager // Low level persistent database to store final content in

	triegc  *prque.Prque // Priority queue mapping block numbers to tries to gc
	chBlock chan gcBlock // chPushBlockGCPrque is a channel for delivering the gc item to gc loop.

	hc            *HeaderChain
	rmLogsFeed    event.Feed
	chainFeed     event.Feed
	chainSideFeed event.Feed
	chainHeadFeed event.Feed
	logsFeed      event.Feed
	scope         event.SubscriptionScope
	genesisBlock  *types.Block

	mu      sync.RWMutex // global mutex for locking chain operations
	chainmu sync.RWMutex // blockchain insertion lock
	procmu  sync.RWMutex // block processor lock

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

	engine    consensus.Engine
	processor Processor // block processor interface
	validator Validator // block and state validator interface
	vmConfig  vm.Config

	badBlocks *lru.Cache // Bad block cache

	parallelDBWrite bool // TODO-Klaytn-Storage parallelDBWrite will be replaced by number of goroutines when worker pool pattern is introduced.

	cachedStateDB       *state.StateDB
	lastUpdatedRootHash common.Hash

	nonceCache   common.Cache
	balanceCache common.Cache

	// State migration
	prepareStateMigration bool
	stopStateMigration    chan struct{}
	committedCnt          int
	pendingCnt            int
	progress              float64
}

// NewBlockChain returns a fully initialised block chain using information
// available in the database. It initialises the default Klaytn validator and
// Processor.
func NewBlockChain(db database.DBManager, cacheConfig *CacheConfig, chainConfig *params.ChainConfig, engine consensus.Engine, vmConfig vm.Config) (*BlockChain, error) {
	if cacheConfig == nil {
		cacheConfig = &CacheConfig{
			StateDBCaching: false,
			ArchiveMode:    false,
			CacheSize:      512,
			BlockInterval:  DefaultBlockInterval,
			TrieCacheLimit: 0,
		}
	}
	// Initialize DeriveSha implementation
	InitDeriveSha(chainConfig.DeriveShaImpl)

	futureBlocks, _ := lru.New(maxFutureBlocks)
	badBlocks, _ := lru.New(maxBadBlocks)

	var nonceCache common.Cache
	var balanceCache common.Cache
	if cacheConfig.TxPoolStateCache {
		nonceCache = common.NewCache(common.FIFOCacheConfig{CacheSize: maxAccountForCache})
		balanceCache = common.NewCache(common.FIFOCacheConfig{CacheSize: maxAccountForCache})
	}

	bc := &BlockChain{
		chainConfig:        chainConfig,
		chainConfigMu:      new(sync.RWMutex),
		cacheConfig:        cacheConfig,
		db:                 db,
		triegc:             prque.New(),
		chBlock:            make(chan gcBlock, 1000),
		stateCache:         state.NewDatabaseWithCache(db, cacheConfig.TrieCacheLimit),
		quit:               make(chan struct{}),
		futureBlocks:       futureBlocks,
		engine:             engine,
		vmConfig:           vmConfig,
		badBlocks:          badBlocks,
		parallelDBWrite:    db.IsParallelDBWrite(),
		nonceCache:         nonceCache,
		balanceCache:       balanceCache,
		stopStateMigration: make(chan struct{}),
	}

	bc.SetValidator(NewBlockValidator(chainConfig, bc, engine))
	bc.SetProcessor(NewStateProcessor(chainConfig, bc, engine))

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
	// Take ownership of this particular state
	go bc.update()
	go bc.gcCachedNodeLoop()
	go bc.restartStateMigration()

	return bc, nil
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
		logger.Error("Head block missing, resetting chain", "hash", head)
		return bc.Reset()
	}
	// Make sure the state associated with the block is available
	if _, err := state.New(currentBlock.Root(), bc.stateCache); err != nil {
		// Dangling block without a state associated, init from scratch
		logger.Error("Head state missing, repairing chain", "number", currentBlock.Number(), "hash", currentBlock.Hash())
		if err := bc.repair(&currentBlock); err != nil {
			return err
		}
	}
	// Everything seems to be fine, set as the head block
	bc.currentBlock.Store(currentBlock)

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

	logger.Info("Loaded most recent local header", "number", currentHeader.Number, "hash", currentHeader.Hash(), "td", headerTd)
	logger.Info("Loaded most recent local full block", "number", currentBlock.Number(), "hash", currentBlock.Hash(), "td", blockTd)
	logger.Info("Loaded most recent local fast block", "number", currentFastBlock.Number(), "hash", currentFastBlock.Hash(), "td", fastTd)

	return nil
}

// SetHead rewinds the local chain to a new head. In the case of headers, everything
// above the new head will be deleted and the new one set. In the case of blocks
// though, the head may be further rewound if block bodies are missing (non-archive
// nodes after a fast sync).
func (bc *BlockChain) SetHead(head uint64) error {
	logger.Info("Rewinding blockchain", "target", head)

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Rewind the header chain, deleting all block bodies until then
	delFn := func(hash common.Hash, num uint64) {
		bc.db.DeleteBody(hash, num)
	}
	bc.hc.SetHead(head, delFn)
	currentHeader := bc.CurrentHeader()

	// Clear out any stale content from the caches
	bc.futureBlocks.Purge()
	bc.db.ClearBlockChainCache()

	// Rewind the block chain, ensuring we don't end up with a stateless head block
	if currentBlock := bc.CurrentBlock(); currentBlock != nil && currentHeader.Number.Uint64() < currentBlock.NumberU64() {
		bc.currentBlock.Store(bc.GetBlock(currentHeader.Hash(), currentHeader.Number.Uint64()))
	}
	if currentBlock := bc.CurrentBlock(); currentBlock != nil {
		if _, err := state.New(currentBlock.Root(), bc.stateCache); err != nil {
			// Rewound state missing, rolled back to before pivot, reset to genesis
			bc.currentBlock.Store(bc.genesisBlock)
		}
	}
	// Rewind the fast block in a simpleton way to the target head
	if currentFastBlock := bc.CurrentFastBlock(); currentFastBlock != nil && currentHeader.Number.Uint64() < currentFastBlock.NumberU64() {
		bc.currentFastBlock.Store(bc.GetBlock(currentHeader.Hash(), currentHeader.Number.Uint64()))
	}
	// If either blocks reached nil, reset to the genesis state
	if currentBlock := bc.CurrentBlock(); currentBlock == nil {
		bc.currentBlock.Store(bc.genesisBlock)
	}
	if currentFastBlock := bc.CurrentFastBlock(); currentFastBlock == nil {
		bc.currentFastBlock.Store(bc.genesisBlock)
	}
	currentBlock := bc.CurrentBlock()
	currentFastBlock := bc.CurrentFastBlock()

	bc.db.WriteHeadBlockHash(currentBlock.Hash())
	bc.db.WriteHeadFastBlockHash(currentFastBlock.Hash())

	return bc.loadLastState()
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

// SetProcessor sets the processor required for making state modifications.
func (bc *BlockChain) SetProcessor(processor Processor) {
	bc.procmu.Lock()
	defer bc.procmu.Unlock()
	bc.processor = processor
}

// SetValidator sets the validator which is used to validate incoming blocks.
func (bc *BlockChain) SetValidator(validator Validator) {
	bc.procmu.Lock()
	defer bc.procmu.Unlock()
	bc.validator = validator
}

// Validator returns the current validator.
func (bc *BlockChain) Validator() Validator {
	bc.procmu.RLock()
	defer bc.procmu.RUnlock()
	return bc.validator
}

// Processor returns the current processor.
func (bc *BlockChain) Processor() Processor {
	bc.procmu.RLock()
	defer bc.procmu.RUnlock()
	return bc.processor
}

// State returns a new mutable state based on the current HEAD block.
func (bc *BlockChain) State() (*state.StateDB, error) {
	return bc.StateAt(bc.CurrentBlock().Root())
}

// StateAt returns a new mutable state based on a particular point in time.
func (bc *BlockChain) StateAt(root common.Hash) (*state.StateDB, error) {
	return state.New(root, bc.stateCache)
}

// StateAtWithGCLock returns a new mutable state based on a particular point in time with read lock of the state nodes.
func (bc *BlockChain) StateAtWithGCLock(root common.Hash) (*state.StateDB, error) {
	bc.RLockGCCachedNode()

	exist := bc.stateCache.TrieDB().DoesExistCachedNode(root)
	if !exist {
		bc.RUnlockGCCachedNode()
		return nil, ErrNotExistNode
	}

	stateDB, err := state.New(root, bc.stateCache)
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

// StateAtWithCache returns a new mutable state based on a particular point in time.
// If different from StateAt() in that it uses state object caching.
func (bc *BlockChain) StateAtWithCache(root common.Hash) (*state.StateDB, error) {
	if bc.cachedStateDB == nil {
		return state.NewWithCache(root, bc.stateCache, state.NewCachedStateObjects())
	} else {
		return state.NewWithCache(root, bc.stateCache, bc.cachedStateDB.GetCachedStateObjects())
	}
}

// TryGetCachedStateDB tries to get cachedStateDB, if StateDBCaching flag is true and it exists.
// It checks the validity of cachedStateDB by comparing saved lastUpdatedRootHash and passed headRootHash.
func (bc *BlockChain) TryGetCachedStateDB(rootHash common.Hash) (*state.StateDB, error) {
	if !bc.cacheConfig.StateDBCaching {
		return bc.StateAt(rootHash)
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// When cachedStateDB is nil, set cachedStateDB with a new StateDB.
	if bc.cachedStateDB == nil {
		if !common.EmptyHash(bc.lastUpdatedRootHash) {
			logger.Error("cachedStateDB is nil, but lastUpdatedRootHash is not common.Hash{}!",
				"lastUpdatedRootHash", bc.lastUpdatedRootHash.String())
			bc.lastUpdatedRootHash = common.Hash{}
		}
		cacheGetStateDBMissMeter.Mark(1)
		return bc.StateAtWithCache(rootHash)
	}

	// If cachedStateDB exists, check if we can use cachedStateDB.
	// If given rootHash is different from lastUpdatedRootHash, return stateDB without cache.
	if rootHash != bc.lastUpdatedRootHash {
		logger.Trace("Given rootHash is different from lastUpdatedRootHash",
			"givenRootHash", rootHash, "lastUpdatedRootHash", bc.lastUpdatedRootHash)
		cacheGetStateDBMissMeter.Mark(1)
		return bc.StateAt(rootHash)
	}
	cacheGetStateDBHitMeter.Mark(1)
	return bc.StateAtWithCache(rootHash)
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
func (bc *BlockChain) repair(head **types.Block) error {
	for {
		// Abort if we've rewound to a head block that does have associated state
		if _, err := state.New((*head).Root(), bc.stateCache); err == nil {
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

	for nr := first; nr <= last; nr++ {
		block := bc.GetBlockByNumber(nr)
		if block == nil {
			return fmt.Errorf("export failed on #%d: not found", nr)
		}

		if err := block.EncodeRLP(w); err != nil {
			return err
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
	close(bc.quit)
	atomic.StoreInt32(&bc.procInterrupt, 1)

	bc.wg.Wait()

	// Ensure the state of a recent block is also stored to disk before exiting.
	// We're writing three different states to catch different restart scenarios:
	//  - HEAD:     So we don't need to reprocess any blocks in the general case
	//  - HEAD-127: So we have a hard limit on the number of blocks reexecuted
	if !bc.isArchiveMode() {
		triedb := bc.stateCache.TrieDB()

		for _, offset := range []uint64{0, triesInMemory - 1} {
			if number := bc.CurrentBlock().NumberU64(); number > offset {
				recent := bc.GetBlockByNumber(number - offset)

				if recent == nil {
					logger.Error("Failed to find recent block from persistent", "blockNumber", number-offset)
					continue
				}

				logger.Info("Writing cached state to disk", "block", recent.Number(), "hash", recent.Hash(), "root", recent.Root())
				if err := triedb.Commit(recent.Root(), true, number-offset); err != nil {
					logger.Error("Failed to commit recent state trie", "err", err)
				}
			}
		}
		for !bc.triegc.Empty() {
			triedb.Dereference(bc.triegc.PopItem().(common.Hash))
		}
		if size, _ := triedb.Size(); size != 0 {
			logger.Error("Dangling trie nodes after full cleanup")
		}
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
	} else {
		// Full but not archive node, do proper garbage collection
		trieDB.Reference(root, common.Hash{}) // metadata reference to keep trie alive

		// If we exceeded our memory allowance, flush matured singleton nodes to disk
		var (
			nodesSize, preimagesSize = trieDB.Size()
			nodesSizeLimit           = common.StorageSize(bc.cacheConfig.CacheSize) * 1024 * 1024
		)
		if nodesSize > nodesSizeLimit || preimagesSize > 4*1024*1024 {
			// NOTE-Klaytn Not to change the original behavior, error is not returned.
			// Error should be returned if it is thought to be safe in the future.
			if err := trieDB.Cap(nodesSizeLimit - database.IdealBatchSize); err != nil {
				logger.Error("Error from trieDB.Cap", "limit", nodesSizeLimit-database.IdealBatchSize)
			}
		}

		if isCommitTrieRequired(bc, block.NumberU64()) {
			logger.Trace("Commit the state trie into the disk", "blocknum", block.NumberU64())
			if err := trieDB.Commit(block.Header().Root, true, block.NumberU64()); err != nil {
				return err
			}

			bc.checkStartStateMigration(block.NumberU64(), root)
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

// gcCachedNodeLoop runs a loop to gc.
func (bc *BlockChain) gcCachedNodeLoop() {
	trieDB := bc.stateCache.TrieDB()

	bc.wg.Add(1)
	go func() {
		defer bc.wg.Done()
		for {
			select {
			case block := <-bc.chBlock:
				bc.triegc.Push(block.root, -float32(block.blockNum))
				logger.Trace("Push GC block", "blkNum", block.blockNum, "hash", block.root.String())

				blkNum := block.blockNum
				if blkNum <= triesInMemory {
					continue
				}

				// Garbage collect anything below our required write retention
				chosen := blkNum - triesInMemory
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
// If BlockChain.parallelDBWrite is true, it calls writeBlockWithStateParallel.
// If not, it calls writeBlockWithStateSerial.
func (bc *BlockChain) WriteBlockWithState(block *types.Block, receipts []*types.Receipt, stateDB *state.StateDB) (WriteStatus, error) {
	var status WriteStatus
	var err error
	if bc.parallelDBWrite {
		status, err = bc.writeBlockWithStateParallel(block, receipts, stateDB)
	} else {
		status, err = bc.writeBlockWithStateSerial(block, receipts, stateDB)
	}

	// TODO-Klaytn-Issue1911 After reviewing the behavior and performance, change the UpdateCacheStateObjects to update at the same time.
	if err != nil {
		return status, err
	}

	if bc.cacheConfig.TxPoolStateCache {
		stateDB.UpdateTxPoolStateCache(bc.nonceCache, bc.balanceCache)
	}

	// Update lastUpdatedRootHash and cachedStateDB after successful WriteBlockWithState.
	if stateDB.UseCachedStateObjects() {
		bc.mu.Lock()
		defer bc.mu.Unlock()

		logger.Trace("Update cached StateDB information", "prevRootHash", bc.lastUpdatedRootHash.String(),
			"newRootHash", block.Root().String(), "newBlockNum", block.NumberU64())

		bc.lastUpdatedRootHash = block.Root()
		stateDB.UpdateCachedStateObjects(block.Root())
		bc.cachedStateDB = stateDB
	}

	return status, err
}

// writeBlockWithStateSerial writes the block and all associated state to the database in serial manner.
func (bc *BlockChain) writeBlockWithStateSerial(block *types.Block, receipts []*types.Receipt, state *state.StateDB) (WriteStatus, error) {
	start := time.Now()
	bc.wg.Add(1)
	defer bc.wg.Done()

	var status WriteStatus
	// Calculate the total blockscore of the block
	ptd := bc.GetTd(block.ParentHash(), block.NumberU64()-1)
	if ptd == nil {
		logger.Error("unknown ancestor (writeBlockWithStateSerial)", "num", block.NumberU64(),
			"hash", block.Hash(), "parentHash", block.ParentHash())
		return NonStatTy, consensus.ErrUnknownAncestor
	}
	// Make sure no inconsistent state is leaked during insertion
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if !bc.ShouldTryInserting(block) {
		return NonStatTy, ErrKnownBlock
	}

	currentBlock := bc.CurrentBlock()
	localTd := bc.GetTd(currentBlock.Hash(), currentBlock.NumberU64())
	externTd := new(big.Int).Add(block.BlockScore(), ptd)

	// Irrelevant of the canonical status, write the block itself to the database
	bc.hc.WriteTd(block.Hash(), block.NumberU64(), externTd)

	// Write other block data.
	bc.writeBlock(block)

	if err := bc.writeStateTrie(block, state); err != nil {
		return NonStatTy, err
	}

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
				return NonStatTy, err
			}
		}
		// Write the positional metadata for transaction/receipt lookups and preimages
		if err := bc.writeTxLookupEntries(block); err != nil {
			return NonStatTy, err
		}
		bc.db.WritePreimages(block.NumberU64(), state.Preimages())
		status = CanonStatTy
	} else {
		status = SideStatTy
	}

	return bc.finalizeWriteBlockWithState(block, status, start)
}

// writeBlockWithStateParallel writes the block and all associated state to the database using goroutines.
func (bc *BlockChain) writeBlockWithStateParallel(block *types.Block, receipts []*types.Receipt, state *state.StateDB) (WriteStatus, error) {
	start := time.Now()
	bc.wg.Add(1)
	defer bc.wg.Done()

	var status WriteStatus
	// Calculate the total blockscore of the block
	ptd := bc.GetTd(block.ParentHash(), block.NumberU64()-1)
	if ptd == nil {
		logger.Error("unknown ancestor (writeBlockWithStateParallel)", "num", block.NumberU64(),
			"hash", block.Hash(), "parentHash", block.ParentHash())
		return NonStatTy, consensus.ErrUnknownAncestor
	}
	// Make sure no inconsistent state is leaked during insertion
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if !bc.ShouldTryInserting(block) {
		return NonStatTy, ErrKnownBlock
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

	go func() {
		defer parallelDBWriteWG.Done()
		if err := bc.writeStateTrie(block, state); err != nil {
			parallelDBWriteErrCh <- err
		}
	}()

	go func() {
		defer parallelDBWriteWG.Done()
		bc.writeReceipts(block.Hash(), block.NumberU64(), receipts)
	}()

	// Wait until all writing goroutines are terminated.
	parallelDBWriteWG.Wait()
	select {
	case err := <-parallelDBWriteErrCh:
		return NonStatTy, err
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
				return NonStatTy, err
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
		return NonStatTy, err
	default:
	}

	return bc.finalizeWriteBlockWithState(block, status, start)
}

// finalizeWriteBlockWithState updates metrics and inserts block when status is CanonStatTy.
func (bc *BlockChain) finalizeWriteBlockWithState(block *types.Block, status WriteStatus, startTime time.Time) (WriteStatus, error) {
	// Set new head.
	if status == CanonStatTy {
		bc.insert(block)
		headBlockNumberGauge.Update(block.Number().Int64())
		blockTxCountsGauge.Update(int64(block.Transactions().Len()))
		blockTxCountsCounter.Inc(int64(block.Transactions().Len()))
	}
	bc.futureBlocks.Remove(block.Hash())

	elapsed := time.Since(startTime)
	logger.Debug("WriteBlockWithState", "blockNum", block.Number(), "parentHash", block.Header().ParentHash, "txs", len(block.Transactions()), "elapsed", elapsed)

	return status, nil
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

	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

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

	abort, results := bc.engine.VerifyHeaders(bc, headers, seals)
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
		// If the header is a banned one, straight out abort
		if BadHashes[block.Hash()] {
			bc.reportBlock(block, nil, ErrBlacklistedHash)
			return i, events, coalescedLogs, ErrBlacklistedHash
		}
		// Wait for the block's verification to complete
		bstart := time.Now()

		err := <-results
		if err == nil {
			err = bc.Validator().ValidateBody(block)
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
			bc.chainmu.Unlock()
			_, evs, logs, err := bc.insertChain(winner)
			bc.chainmu.Lock()
			events, coalescedLogs = evs, logs

			if err != nil {
				return i, events, coalescedLogs, err
			}

		case err != nil:
			bc.reportBlock(block, nil, err)
			return i, events, coalescedLogs, err
		}
		// Create a new trie using the parent block and report an
		// error if it fails.
		var parent *types.Block
		if i == 0 {
			parent = bc.GetBlock(block.ParentHash(), block.NumberU64()-1)
		} else {
			parent = chain[i-1]
		}

		logger.Info("TryGetCachedStateDB", "parent.Root()", parent.Root())
		stateDB, err := bc.TryGetCachedStateDB(parent.Root())
		if err != nil {
			return i, events, coalescedLogs, err
		}

		// for debug
		start := time.Now()

		// Process block using the parent state as reference point.
		receipts, logs, usedGas, err := bc.processor.Process(block, stateDB, bc.vmConfig)
		if err != nil {
			bc.reportBlock(block, receipts, err)
			return i, events, coalescedLogs, err
		}

		if block.Transactions().Len() > 0 {

			elapsed := time.Since(start)
			logger.Debug("blockchain.blockchain processing block", "elapsed", elapsed, "txs", block.Transactions().Len())
		}

		// Validate the state using the default validator
		err = bc.Validator().ValidateState(block, parent, stateDB, receipts, usedGas)
		if err != nil {

			bc.reportBlock(block, receipts, err)
			return i, events, coalescedLogs, err
		}

		// Write the block to the chain and get the status.
		status, err := bc.WriteBlockWithState(block, receipts, stateDB)
		if err != nil {
			if err == ErrKnownBlock {
				logger.Debug("Tried to insert already known block", "num", block.NumberU64(), "hash", block.Hash().String())
				continue
			}
			return i, events, coalescedLogs, err
		}

		switch status {
		case CanonStatTy:
			logger.Debug("Inserted new block", "number", block.Number(), "hash", block.Hash(),
				"txs", len(block.Transactions()), "gas", block.GasUsed(), "elapsed", common.PrettyDuration(time.Since(bstart)))

			coalescedLogs = append(coalescedLogs, logs...)
			events = append(events, ChainEvent{block, block.Hash(), logs})
			lastCanon = block

		case SideStatTy:
			logger.Debug("Inserted forked block", "number", block.Number(), "hash", block.Hash(), "diff", block.BlockScore(), "elapsed",
				common.PrettyDuration(time.Since(bstart)), "txs", len(block.Transactions()), "gas", block.GasUsed())

			events = append(events, ChainSideEvent{block})
		}
		blockInsertTimeGauge.Update(int64(time.Since(bstart)))
		stats.processed++
		stats.usedGas += usedGas

		cache, _ := bc.stateCache.TrieDB().Size()
		stats.report(chain, i, cache)
	}
	// Append a single chain head event if we've progressed the chain
	if lastCanon != nil && bc.CurrentBlock().Hash() == lastCanon.Hash() {
		events = append(events, ChainHeadEvent{lastCanon})
	}
	return 0, events, coalescedLogs, nil
}

// insertStats tracks and reports on block insertion.
type insertStats struct {
	queued, processed, ignored int
	usedGas                    uint64
	lastIndex                  int
	startTime                  mclock.AbsTime
}

// statsReportLimit is the time limit during import after which we always print
// out progress. This avoids the user wondering what's going on.
const statsReportLimit = 8 * time.Second

// report prints statistics if some number of blocks have been processed
// or more than a few seconds have passed since the last message.
func (st *insertStats) report(chain []*types.Block, index int, cache common.StorageSize) {
	// Fetch the timings for the batch
	var (
		now     = mclock.Now()
		elapsed = time.Duration(now) - time.Duration(st.startTime)
	)
	// If we're at the last block of the batch or report period reached, log
	if index == len(chain)-1 || elapsed >= statsReportLimit {
		var (
			end = chain[index]
			txs = countTransactions(chain[st.lastIndex : index+1])
		)
		context := []interface{}{
			"number", end.Number(), "hash", end.Hash().String(), "blocks", st.processed, "txs", txs, "elapsed", common.PrettyDuration(elapsed),
			"trieDBSize", cache, "mgas", float64(st.usedGas) / 1000000, "mgasps", float64(st.usedGas) * 1000 / float64(elapsed),
		}
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
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	bc.wg.Add(1)
	defer bc.wg.Done()

	whFunc := func(header *types.Header) error {
		bc.mu.Lock()
		defer bc.mu.Unlock()

		_, err := bc.hc.WriteHeader(header)
		return err
	}

	return bc.hc.InsertHeaderChain(chain, whFunc, start)
}

// writeHeader writes a header into the local chain, given that its parent is
// already known. If the total blockscore of the newly inserted header becomes
// greater than the current known TD, the canonical chain is re-routed.
//
// Note: This method is not concurrent-safe with inserting blocks simultaneously
// into the chain, as side effects caused by reorganisations cannot be emulated
// without the real blocks. Hence, writing headers directly should only be done
// in two scenarios: pure-header mode of operation (light clients), or properly
// separated header/block phases (non-archive clients).
func (bc *BlockChain) writeHeader(header *types.Header) error {
	bc.wg.Add(1)
	defer bc.wg.Done()

	bc.mu.Lock()
	defer bc.mu.Unlock()

	_, err := bc.hc.WriteHeader(header)
	return err
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

// GetNonceCache returns a nonceCache.
func (bc *BlockChain) GetNonceCache() common.Cache {
	return bc.nonceCache
}

// GetBalanceCache returns a balanceCache.
func (bc *BlockChain) GetBalanceCache() common.Cache {
	return bc.balanceCache
}

// GetNonceInCache returns (cachedNonce, true) if nonce exists in cache.
// If not, it returns (0, false).
func (bc *BlockChain) GetNonceInCache(addr common.Address) (uint64, bool) {
	nonceCache := bc.GetNonceCache()

	if nonceCache != nil {
		if obj, exist := nonceCache.Get(addr); exist && obj != nil {
			nonce, _ := obj.(uint64)
			return nonce, true
		}
	}
	return 0, false
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func (bc *BlockChain) ApplyTransaction(config *params.ChainConfig, author *common.Address, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg *vm.Config) (*types.Receipt, uint64, error) {

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
		return nil, 0, err
	}

	msg, err := tx.AsMessageWithAccountKeyPicker(types.MakeSigner(config, header.Number), statedb, blockNumber)
	if err != nil {
		return nil, 0, err
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	// Apply the transaction to the current state (included in the env)
	_, gas, kerr := ApplyMessage(vmenv, msg)
	err = kerr.ErrTxInvalid
	if err != nil {
		return nil, 0, err
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

	return receipt, gas, err
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
