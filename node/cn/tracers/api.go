// Modifications Copyright 2022 The klaytn Authors
// Copyright 2021 The go-ethereum Authors
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
// This file is derived from eth/tracers/api.go (2022/08/08).
// Modified and improved for the klaytn development.

package tracers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	klaytnapi "github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
)

const (
	// defaultTraceTimeout is the amount of time a single transaction can execute
	// by default before being forcefully aborted.
	defaultTraceTimeout = 5 * time.Second

	// defaultLoggerTimeout is the amount of time a logger can aggregate trace logs
	defaultLoggerTimeout = 1 * time.Second

	// defaultTraceReexec is the number of blocks the tracer is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	defaultTraceReexec = uint64(128)

	// defaultTracechainMemLimit is the size of the triedb, at which traceChain
	// switches over and tries to use a disk-backed database instead of building
	// on top of memory.
	// For non-archive nodes, this limit _will_ be overblown, as disk-backed tries
	// will only be found every ~15K blocks or so.
	// For klaytn, this value is set to a value 4 times larger compared to the ethereum setting.
	defaultTracechainMemLimit = common.StorageSize(4 * 500 * 1024 * 1024)

	// fastCallTracer is the go-version callTracer which is lighter and faster than
	// Javascript version.
	fastCallTracer = "fastCallTracer"
)

var (
	HeavyAPIRequestLimit int32 = 500
	heavyAPIRequestCount int32 = 0
)

// Backend interface provides the common API services with access to necessary functions.
type Backend interface {
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error)
	GetTxAndLookupInfo(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	RPCGasCap() *big.Int
	ChainConfig() *params.ChainConfig
	ChainDB() database.DBManager
	Engine() consensus.Engine
	// StateAtBlock returns the state corresponding to the stateroot of the block.
	// N.B: For executing transactions on block N, the required stateRoot is block N-1,
	// so this method should be called with the parent.
	StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (*state.StateDB, error)
	StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (blockchain.Message, vm.Context, *state.StateDB, error)
}

// API is the collection of tracing APIs exposed over the private debugging endpoint.
type API struct {
	backend     Backend
	unsafeTrace bool
}

// NewAPIUnsafeDisabled creates a new API definition for the tracing methods of the CN service,
// only allowing predefined tracers.
func NewAPIUnsafeDisabled(backend Backend) *API {
	return &API{backend: backend, unsafeTrace: false}
}

// NewAPI creates a new API definition for the tracing methods of the CN service,
// allowing both predefined tracers and Javascript snippet based tracing.
func NewAPI(backend Backend) *API {
	return &API{backend: backend, unsafeTrace: true}
}

type chainContext struct {
	backend Backend
	ctx     context.Context
}

func (context *chainContext) Engine() consensus.Engine {
	return context.backend.Engine()
}

func (context *chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	header, err := context.backend.HeaderByNumber(context.ctx, rpc.BlockNumber(number))
	if err != nil {
		return nil
	}
	if header.Hash() == hash {
		return header
	}
	header, err = context.backend.HeaderByHash(context.ctx, hash)
	if err != nil {
		return nil
	}
	return header
}

// chainContext constructs the context reader which is used by the evm for reading
// the necessary chain context.
func newChainContext(ctx context.Context, backend Backend) blockchain.ChainContext {
	return &chainContext{backend: backend, ctx: ctx}
}

// blockByNumber is the wrapper of the chain access function offered by the backend.
// It will return an error if the block is not found.
func (api *API) blockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	return api.backend.BlockByNumber(ctx, number)
}

// blockByHash is the wrapper of the chain access function offered by the backend.
// It will return an error if the block is not found.
func (api *API) blockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return api.backend.BlockByHash(ctx, hash)
}

// blockByNumberAndHash is the wrapper of the chain access function offered by
// the backend. It will return an error if the block is not found.
func (api *API) blockByNumberAndHash(ctx context.Context, number rpc.BlockNumber, hash common.Hash) (*types.Block, error) {
	block, err := api.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	if block.Hash() == hash {
		return block, nil
	}
	return api.blockByHash(ctx, hash)
}

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*vm.LogConfig
	Tracer        *string
	Timeout       *string
	LoggerTimeout *string
	Reexec        *uint64
}

// StdTraceConfig holds extra parameters to standard-json trace functions.
type StdTraceConfig struct {
	*vm.LogConfig
	Reexec *uint64
	TxHash common.Hash
}

// txTraceResult is the result of a single transaction trace.
type txTraceResult struct {
	TxHash common.Hash `json:"txHash,omitempty"` // transaction hash
	Result interface{} `json:"result,omitempty"` // Trace results produced by the tracer
	Error  string      `json:"error,omitempty"`  // Trace failure produced by the tracer
}

// blockTraceTask represents a single block trace task when an entire chain is
// being traced.
type blockTraceTask struct {
	statedb *state.StateDB   // Intermediate state prepped for tracing
	block   *types.Block     // Block to trace the transactions from
	rootref common.Hash      // Trie root reference held for this task
	results []*txTraceResult // Trace results procudes by the task
}

// blockTraceResult represets the results of tracing a single block when an entire
// chain is being traced.
type blockTraceResult struct {
	Block  hexutil.Uint64   `json:"block"`  // Block number corresponding to this trace
	Hash   common.Hash      `json:"hash"`   // Block hash corresponding to this trace
	Traces []*txTraceResult `json:"traces"` // Trace results produced by the task
}

// txTraceTask represents a single transaction trace task when an entire block
// is being traced.
type txTraceTask struct {
	statedb *state.StateDB // Intermediate state prepped for tracing
	index   int            // Transaction offset in the block
}

func checkRangeAndReturnBlock(api *API, ctx context.Context, start, end rpc.BlockNumber) (*types.Block, *types.Block, error) {
	// Fetch the block interval that we want to trace
	from, err := api.blockByNumber(ctx, start)
	if err != nil {
		return nil, nil, err
	}
	to, err := api.blockByNumber(ctx, end)
	if err != nil {
		return nil, nil, err
	}

	// Trace the chain if we've found all our blocks
	if from == nil {
		return nil, nil, fmt.Errorf("starting block #%d not found", start)
	}
	if to == nil {
		return nil, nil, fmt.Errorf("end block #%d not found", end)
	}
	if from.Number().Cmp(to.Number()) >= 0 {
		return nil, nil, fmt.Errorf("end block #%d needs to come after start block #%d", end, start)
	}
	return from, to, nil
}

// TraceChain returns the structured logs created during the execution of EVM
// between two blocks (excluding start) and returns them as a JSON object.
func (api *API) TraceChain(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (*rpc.Subscription, error) {
	if !api.unsafeTrace {
		return nil, errors.New("TraceChain is disabled")
	}
	from, to, err := checkRangeAndReturnBlock(api, ctx, start, end)
	if err != nil {
		return nil, err
	}
	// Tracing a chain is a **long** operation, only do with subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	var sub *rpc.Subscription
	sub = notifier.CreateSubscription()
	_, err = api.traceChain(from, to, config, notifier, sub)
	return sub, err
}

// traceChain configures a new tracer according to the provided configuration, and
// executes all the transactions contained within.
// The traceChain operates in two modes: subscription mode and rpc mode
//   - if notifier and sub is not nil, it works as a subscription mode and returns nothing
//   - if those parameters are nil, it works as a rpc mode and returns the block trace results, so it can pass the result through rpc-call
func (api *API) traceChain(start, end *types.Block, config *TraceConfig, notifier *rpc.Notifier, sub *rpc.Subscription) (map[uint64]*blockTraceResult, error) {
	// Prepare all the states for tracing. Note this procedure can take very
	// long time. Timeout mechanism is necessary.
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	// Execute all the transaction contained within the chain concurrently for each block
	blocks := int(end.NumberU64() - start.NumberU64())
	threads := runtime.NumCPU()
	if threads > blocks {
		threads = blocks
	}
	var (
		pend     = new(sync.WaitGroup)
		tasks    = make(chan *blockTraceTask, threads)
		results  = make(chan *blockTraceTask, threads)
		localctx = context.Background()
	)
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()

			// Fetch and execute the next block trace tasks
			for task := range tasks {
				signer := types.MakeSigner(api.backend.ChainConfig(), task.block.Number())

				// Trace all the transactions contained within
				for i, tx := range task.block.Transactions() {
					msg, err := tx.AsMessageWithAccountKeyPicker(signer, task.statedb, task.block.NumberU64())
					if err != nil {
						logger.Warn("Tracing failed", "hash", tx.Hash(), "block", task.block.NumberU64(), "err", err)
						task.results[i] = &txTraceResult{TxHash: tx.Hash(), Error: err.Error()}
						break
					}

					vmctx := blockchain.NewEVMContext(msg, task.block.Header(), newChainContext(localctx, api.backend), nil)

					res, err := api.traceTx(localctx, msg, vmctx, task.statedb, config)
					if err != nil {
						task.results[i] = &txTraceResult{TxHash: tx.Hash(), Error: err.Error()}
						logger.Warn("Tracing failed", "hash", tx.Hash(), "block", task.block.NumberU64(), "err", err)
						break
					}
					task.statedb.Finalise(true, true)
					task.results[i] = &txTraceResult{TxHash: tx.Hash(), Result: res}
				}
				if notifier != nil {
					// Stream the result back to the user or abort on teardown
					select {
					case results <- task:
					case <-notifier.Closed():
						return
					}
				} else {
					results <- task
				}
			}
		}()
	}
	// Start a goroutine to feed all the blocks into the tracers
	begin := time.Now()

	go func() {
		var (
			logged  time.Time
			number  uint64
			traced  uint64
			failed  error
			parent  common.Hash
			statedb *state.StateDB
		)
		// Ensure everything is properly cleaned up on any exit path
		defer func() {
			close(tasks)
			pend.Wait()

			switch {
			case failed != nil:
				logger.Warn("Chain tracing failed", "start", start.NumberU64(), "end", end.NumberU64(), "transactions", traced, "elapsed", time.Since(begin), "err", failed)
			case number < end.NumberU64():
				logger.Warn("Chain tracing aborted", "start", start.NumberU64(), "end", end.NumberU64(), "abort", number, "transactions", traced, "elapsed", time.Since(begin))
			default:
				logger.Info("Chain tracing finished", "start", start.NumberU64(), "end", end.NumberU64(), "transactions", traced, "elapsed", time.Since(begin))
			}
			close(results)
		}()
		var preferDisk bool
		// Feed all the blocks both into the tracer, as well as fast process concurrently
		for number = start.NumberU64(); number < end.NumberU64(); number++ {
			if notifier != nil {
				// Stop tracing if interruption was requested
				select {
				case <-notifier.Closed():
					return
				default:
				}
			}
			// Print progress logs if long enough time elapsed
			if time.Since(logged) > log.StatsReportLimit {
				logged = time.Now()
				logger.Info("Tracing chain segment", "start", start.NumberU64(), "end", end.NumberU64(), "current", number, "transactions", traced, "elapsed", time.Since(begin))
			}
			// Retrieve the parent state to trace on top
			block, err := api.blockByNumber(localctx, rpc.BlockNumber(number))
			if err != nil {
				failed = err
				break
			}
			// Prepare the statedb for tracing. Don't use the live database for
			// tracing to avoid persisting state junks into the database.
			statedb, err = api.backend.StateAtBlock(localctx, block, reexec, statedb, false, preferDisk)
			if err != nil {
				failed = err
				break
			}
			if trieDb := statedb.Database().TrieDB(); trieDb != nil {
				// Hold the reference for tracer, will be released at the final stage
				trieDb.ReferenceRoot(block.Root())

				// Release the parent state because it's already held by the tracer
				if !common.EmptyHash(parent) {
					trieDb.Dereference(parent)
				}
				// Prefer disk if the trie db memory grows too much
				s1, s2, s3 := trieDb.Size()
				if !preferDisk && (s1+s2+s3) > defaultTracechainMemLimit {
					logger.Info("Switching to prefer-disk mode for tracing", "size", s1+s2+s3)
					preferDisk = true
				}
			}
			parent = block.Root()

			next, err := api.blockByNumber(localctx, rpc.BlockNumber(number+1))
			if err != nil {
				failed = err
				break
			}
			// Send the block over to the concurrent tracers (if not in the fast-forward phase)
			txs := next.Transactions()
			if notifier != nil {
				select {
				case tasks <- &blockTraceTask{statedb: statedb.Copy(), block: next, rootref: block.Root(), results: make([]*txTraceResult, len(txs))}:
				case <-notifier.Closed():
					return
				}
			} else {
				tasks <- &blockTraceTask{statedb: statedb.Copy(), block: next, rootref: block.Root(), results: make([]*txTraceResult, len(txs))}
			}
			traced += uint64(len(txs))
		}
	}()

	waitForResult := func() map[uint64]*blockTraceResult {
		// Keep reading the trace results and stream the to the user
		var (
			done = make(map[uint64]*blockTraceResult)
			next = start.NumberU64() + 1
		)
		for res := range results {
			// Queue up next received result
			result := &blockTraceResult{
				Block:  hexutil.Uint64(res.block.NumberU64()),
				Hash:   res.block.Hash(),
				Traces: res.results,
			}
			done[uint64(result.Block)] = result

			// Dereference any parent tries held in memory by this task
			if res.statedb.Database().TrieDB() != nil {
				res.statedb.Database().TrieDB().Dereference(res.rootref)
			}
			if notifier != nil {
				// Stream completed traces to the user, aborting on the first error
				for result, ok := done[next]; ok; result, ok = done[next] {
					if len(result.Traces) > 0 || next == end.NumberU64() {
						notifier.Notify(sub.ID, result)
					}
					delete(done, next)
					next++
				}
			} else {
				if len(done) == blocks {
					return done
				}
			}
		}
		return nil
	}

	if notifier != nil {
		go waitForResult()
		return nil, nil
	}

	return waitForResult(), nil
}

// TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *API) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]*txTraceResult, error) {
	block, err := api.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockByNumberRange returns the ranged blocks tracing results
// TODO-tracer: limit the result by the size of the return
func (api *API) TraceBlockByNumberRange(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (map[uint64]*blockTraceResult, error) {
	if !api.unsafeTrace {
		return nil, fmt.Errorf("TraceBlockByNumberRange is disabled")
	}
	// When the block range is [start,end], the actual tracing block would be [start+1,end]
	// this is the reason why we change the block range to [start-1, end] so that we can trace [start,end] blocks
	from, to, err := checkRangeAndReturnBlock(api, ctx, start-1, end)
	if err != nil {
		return nil, err
	}
	return api.traceChain(from, to, config, nil, nil)
}

// TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *API) TraceBlockByHash(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	block, err := api.blockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlock returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *API) TraceBlock(ctx context.Context, blob hexutil.Bytes, config *TraceConfig) ([]*txTraceResult, error) {
	block := new(types.Block)
	if err := rlp.Decode(bytes.NewReader(blob), block); err != nil {
		return nil, fmt.Errorf("could not decode block: %v", err)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockFromFile returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *API) TraceBlockFromFile(ctx context.Context, file string, config *TraceConfig) ([]*txTraceResult, error) {
	if !api.unsafeTrace {
		return nil, errors.New("TraceBlockFromFile is disabled")
	}
	blob, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	return api.TraceBlock(ctx, common.Hex2Bytes(string(blob)), config)
}

// TraceBadBlock returns the structured logs created during the execution of
// EVM against a block pulled from the pool of bad ones and returns them as a JSON
// object.
func (api *API) TraceBadBlock(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	blocks, err := api.backend.ChainDB().ReadAllBadBlocks()
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if block.Hash() == hash {
			return api.traceBlock(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}

// StandardTraceBlockToFile dumps the structured logs created during the
// execution of EVM to the local file system and returns a list of files
// to the caller.
func (api *API) StandardTraceBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([]string, error) {
	if !api.unsafeTrace {
		return nil, errors.New("StandardTraceBlockToFile is disabled")
	}
	block, err := api.blockByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("block %#x not found", hash)
	}
	return api.standardTraceBlockToFile(ctx, block, config)
}

// StandardTraceBadBlockToFile dumps the structured logs created during the
// execution of EVM against a block pulled from the pool of bad ones to the
// local file system and returns a list of files to the caller.
func (api *API) StandardTraceBadBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([]string, error) {
	if !api.unsafeTrace {
		return nil, errors.New("StandardTraceBadBlockToFile is disabled")
	}
	blocks, err := api.backend.ChainDB().ReadAllBadBlocks()
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if block.Hash() == hash {
			return api.standardTraceBlockToFile(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}

// traceBlock configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requestd tracer.
func (api *API) traceBlock(ctx context.Context, block *types.Block, config *TraceConfig) ([]*txTraceResult, error) {
	if !api.unsafeTrace {
		if atomic.LoadInt32(&heavyAPIRequestCount) >= HeavyAPIRequestLimit {
			return nil, fmt.Errorf("heavy debug api requests exceed the limit: %d", int64(HeavyAPIRequestLimit))
		}
		atomic.AddInt32(&heavyAPIRequestCount, 1)
		defer atomic.AddInt32(&heavyAPIRequestCount, -1)
	}
	if block.NumberU64() == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	// Create the parent state database
	parent, err := api.blockByNumberAndHash(ctx, rpc.BlockNumber(block.NumberU64()-1), block.ParentHash())
	if err != nil {
		return nil, err
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}

	statedb, err := api.backend.StateAtBlock(ctx, parent, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	// Execute all the transaction contained within the block concurrently
	var (
		signer  = types.MakeSigner(api.backend.ChainConfig(), block.Number())
		txs     = block.Transactions()
		results = make([]*txTraceResult, len(txs))

		pend = new(sync.WaitGroup)
		jobs = make(chan *txTraceTask, len(txs))
	)
	threads := runtime.NumCPU()
	if threads > len(txs) {
		threads = len(txs)
	}
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()

			// Fetch and execute the next transaction trace tasks
			for task := range jobs {
				msg, err := txs[task.index].AsMessageWithAccountKeyPicker(signer, task.statedb, block.NumberU64())
				if err != nil {
					logger.Warn("Tracing failed", "tx idx", task.index, "block", block.NumberU64(), "err", err)
					results[task.index] = &txTraceResult{TxHash: txs[task.index].Hash(), Error: err.Error()}
					continue
				}

				vmctx := blockchain.NewEVMContext(msg, block.Header(), newChainContext(ctx, api.backend), nil)
				res, err := api.traceTx(ctx, msg, vmctx, task.statedb, config)
				if err != nil {
					results[task.index] = &txTraceResult{TxHash: txs[task.index].Hash(), Error: err.Error()}
					continue
				}
				results[task.index] = &txTraceResult{TxHash: txs[task.index].Hash(), Result: res}
			}
		}()
	}
	// Feed the transactions into the tracers and return
	var failed error
	for i, tx := range txs {
		// Send the trace task over for execution
		jobs <- &txTraceTask{statedb: statedb.Copy(), index: i}

		// Generate the next state snapshot fast without tracing
		msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		if err != nil {
			logger.Warn("Tracing failed", "hash", tx.Hash(), "block", block.NumberU64(), "err", err)
			failed = err
			break
		}

		vmctx := blockchain.NewEVMContext(msg, block.Header(), newChainContext(ctx, api.backend), nil)
		vmenv := vm.NewEVM(vmctx, statedb, api.backend.ChainConfig(), &vm.Config{UseOpcodeComputationCost: true})
		if _, err = blockchain.ApplyMessage(vmenv, msg); err != nil {
			failed = err
			break
		}
		// Finalize the state so any modifications are written to the trie
		statedb.Finalise(true, true)
	}
	close(jobs)
	pend.Wait()

	// If execution failed in between, abort
	if failed != nil {
		return nil, failed
	}
	return results, nil
}

// standardTraceBlockToFile configures a new tracer which uses standard JSON output,
// and traces either a full block or an individual transaction. The return value will
// be one filename per transaction traced.
func (api *API) standardTraceBlockToFile(ctx context.Context, block *types.Block, config *StdTraceConfig) ([]string, error) {
	// If we're tracing a single transaction, make sure it's present
	if config != nil && !common.EmptyHash(config.TxHash) {
		if !containsTx(block, config.TxHash) {
			return nil, fmt.Errorf("transaction %#x not found in block", config.TxHash)
		}
	}
	if block.NumberU64() == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	parent, err := api.blockByNumberAndHash(ctx, rpc.BlockNumber(block.NumberU64()-1), block.ParentHash())
	if err != nil {
		return nil, err
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	statedb, err := api.backend.StateAtBlock(ctx, parent, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	// Retrieve the tracing configurations, or use default values
	var (
		logConfig vm.LogConfig
		txHash    common.Hash
	)
	if config != nil {
		if config.LogConfig != nil {
			logConfig = *config.LogConfig
		}
		txHash = config.TxHash
	}
	logConfig.Debug = true

	// Execute transaction, either tracing all or just the requested one
	var (
		signer = types.MakeSigner(api.backend.ChainConfig(), block.Number())
		dumps  []string
	)
	for i, tx := range block.Transactions() {
		// Prepare the transaction for un-traced execution
		msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		if err != nil {
			logger.Warn("Tracing failed", "hash", tx.Hash(), "block", block.NumberU64(), "err", err)
			return nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}

		var (
			vmctx = blockchain.NewEVMContext(msg, block.Header(), newChainContext(ctx, api.backend), nil)

			vmConf vm.Config
			dump   *os.File
		)

		// If the transaction needs tracing, swap out the configs
		if tx.Hash() == txHash || common.EmptyHash(txHash) {
			// Generate a unique temporary file to dump it into
			prefix := fmt.Sprintf("block_%#x-%d-%#x-", block.Hash().Bytes()[:4], i, tx.Hash().Bytes()[:4])

			dump, err = ioutil.TempFile(os.TempDir(), prefix)
			if err != nil {
				return nil, err
			}
			dumps = append(dumps, dump.Name())

			// Swap out the noop logger to the standard tracer
			vmConf = vm.Config{
				Debug:                    true,
				Tracer:                   vm.NewJSONLogger(&logConfig, bufio.NewWriter(dump)),
				EnablePreimageRecording:  true,
				UseOpcodeComputationCost: true,
			}
		}
		// Execute the transaction and flush any traces to disk
		vmenv := vm.NewEVM(vmctx, statedb, api.backend.ChainConfig(), &vmConf)
		_, err = blockchain.ApplyMessage(vmenv, msg)

		if dump != nil {
			dump.Close()
			logger.Info("Wrote standard trace", "file", dump.Name())
		}
		if err != nil {
			return dumps, err
		}
		// Finalize the state so any modifications are written to the trie
		statedb.Finalise(true, true)

		// If we've traced the transaction we were looking for, abort
		if tx.Hash() == txHash {
			break
		}
	}
	return dumps, nil
}

// containsTx reports whether the transaction with a certain hash
// is contained within the specified block.
func containsTx(block *types.Block, hash common.Hash) bool {
	for _, tx := range block.Transactions() {
		if tx.Hash() == hash {
			return true
		}
	}
	return false
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *API) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	if !api.unsafeTrace {
		if atomic.LoadInt32(&heavyAPIRequestCount) >= HeavyAPIRequestLimit {
			return nil, fmt.Errorf("heavy debug api requests exceed the limit: %d", int64(HeavyAPIRequestLimit))
		}
		atomic.AddInt32(&heavyAPIRequestCount, 1)
		defer atomic.AddInt32(&heavyAPIRequestCount, -1)
	}
	// Retrieve the transaction and assemble its EVM context
	tx, blockHash, blockNumber, index := api.backend.GetTxAndLookupInfo(hash)
	if tx == nil {
		return nil, fmt.Errorf("transaction %#x not found", hash)
	}
	// It shouldn't happen in practice.
	if blockNumber == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	block, err := api.blockByNumberAndHash(ctx, rpc.BlockNumber(blockNumber), blockHash)
	if err != nil {
		return nil, err
	}
	msg, vmctx, statedb, err := api.backend.StateAtTransaction(ctx, block, int(index), reexec)
	if err != nil {
		return nil, err
	}
	// Trace the transaction and return
	return api.traceTx(ctx, msg, vmctx, statedb, config)
}

// TraceCall lets you trace a given klay_call. It collects the structured logs
// created during the execution of EVM if the given transaction was added on
// top of the provided block and returns them as a JSON object.
func (api *API) TraceCall(ctx context.Context, args klaytnapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (interface{}, error) {
	if !api.unsafeTrace {
		if atomic.LoadInt32(&heavyAPIRequestCount) >= HeavyAPIRequestLimit {
			return nil, fmt.Errorf("heavy debug api requests exceed the limit: %d", int64(HeavyAPIRequestLimit))
		}
		atomic.AddInt32(&heavyAPIRequestCount, 1)
		defer atomic.AddInt32(&heavyAPIRequestCount, -1)
	}
	// Try to retrieve the specified block
	var (
		err   error
		block *types.Block
	)
	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err = api.blockByHash(ctx, hash)
	} else if number, ok := blockNrOrHash.Number(); ok {
		block, err = api.blockByNumber(ctx, number)
	} else {
		return nil, errors.New("invalid arguments; neither block nor hash specified")
	}
	if err != nil {
		return nil, err
	}
	// try to recompute the state
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	statedb, err := api.backend.StateAtBlock(ctx, block, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}

	// Execute the trace
	intrinsicGas, err := types.IntrinsicGas(args.InputData(), nil, args.To == nil, api.backend.ChainConfig().Rules(block.Number()))
	if err != nil {
		return nil, err
	}
	basefee := new(big.Int).SetUint64(params.ZeroBaseFee)
	if block.Header().BaseFee != nil {
		basefee = block.Header().BaseFee
	}
	gasCap := uint64(0)
	if rpcGasCap := api.backend.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}
	msg, err := args.ToMessage(gasCap, basefee, intrinsicGas)
	if err != nil {
		return nil, err
	}
	vmctx := blockchain.NewEVMContext(msg, block.Header(), newChainContext(ctx, api.backend), nil)

	// Add gas fee to sender for estimating gasLimit/computing cost or calling a function by insufficient balance sender.
	statedb.AddBalance(msg.ValidatedSender(), new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), basefee))

	return api.traceTx(ctx, msg, vmctx, statedb, config)
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func (api *API) traceTx(ctx context.Context, message blockchain.Message, vmctx vm.Context, statedb *state.StateDB, config *TraceConfig) (interface{}, error) {
	// Assemble the structured logger or the JavaScript tracer
	var (
		tracer vm.Tracer
		err    error
	)
	switch {
	case config != nil && config.Tracer != nil:
		// Define a meaningful timeout of a single transaction trace
		timeout := defaultTraceTimeout
		if config.Timeout != nil {
			if timeout, err = time.ParseDuration(*config.Timeout); err != nil {
				return nil, err
			}
		}

		if *config.Tracer == fastCallTracer {
			tracer = vm.NewInternalTxTracer()
		} else {
			// Construct the JavaScript tracer to execute with
			if tracer, err = New(*config.Tracer, api.unsafeTrace); err != nil {
				return nil, err
			}
		}
		// Handle timeouts and RPC cancellations
		deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
		go func() {
			<-deadlineCtx.Done()
			if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
				switch t := tracer.(type) {
				case *Tracer:
					t.Stop(errors.New("execution timeout"))
				case *vm.InternalTxTracer:
					t.Stop(errors.New("execution timeout"))
				default:
					logger.Warn("unknown tracer type", "type", reflect.TypeOf(t).String())
				}
			}
		}()
		defer cancel()

	case config == nil:
		tracer = vm.NewStructLogger(nil)

	default:
		tracer = vm.NewStructLogger(config.LogConfig)
	}
	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(vmctx, statedb, api.backend.ChainConfig(), &vm.Config{Debug: true, Tracer: tracer, UseOpcodeComputationCost: true})

	ret, err := blockchain.ApplyMessage(vmenv, message)
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %v", err)
	}
	// Depending on the tracer type, format and return the output
	switch tracer := tracer.(type) {
	case *vm.StructLogger:
		loggerTimeout := defaultLoggerTimeout
		if config != nil && config.LoggerTimeout != nil {
			if loggerTimeout, err = time.ParseDuration(*config.LoggerTimeout); err != nil {
				return nil, err
			}
		}
		if logs, err := klaytnapi.FormatLogs(loggerTimeout, tracer.StructLogs()); err == nil {
			return &klaytnapi.ExecutionResult{
				Gas:         ret.UsedGas,
				Failed:      ret.Failed(),
				ReturnValue: fmt.Sprintf("%x", ret.Return()),
				StructLogs:  logs,
			}, nil
		} else {
			return nil, err
		}

	case *Tracer:
		return tracer.GetResult()
	case *vm.InternalTxTracer:
		return tracer.GetResult()

	default:
		panic(fmt.Sprintf("bad tracer type %T", tracer))
	}
}
