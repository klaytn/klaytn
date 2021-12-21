// Modifications Copyright 2019 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from eth/api_tracer.go (2018/06/04).
// Modified and improved for the klaytn development.

package cn

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	klaytnapi "github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn/tracers"
	"github.com/klaytn/klaytn/rlp"
	statedb2 "github.com/klaytn/klaytn/storage/statedb"
)

const (
	// defaultTraceTimeout is the amount of time a single transaction can execute
	// by default before being forcefully aborted.
	defaultTraceTimeout = 5 * time.Second

	// defaultTraceReexec is the number of blocks the tracer is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	defaultTraceReexec = uint64(128)

	// fastCallTracer is the go-version callTracer which is lighter and faster than
	// Javascript version.
	fastCallTracer = "fastCallTracer"
)

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*vm.LogConfig
	Tracer  *string
	Timeout *string
	Reexec  *uint64
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

// TraceChain returns the structured logs created during the execution of EVM
// between two blocks (excluding start) and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceChain(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (*rpc.Subscription, error) {
	// Fetch the block interval that we want to trace
	var from, to *types.Block
	switch start {
	case rpc.PendingBlockNumber:
		from = api.cn.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		from = api.cn.blockchain.CurrentBlock()
	default:
		from = api.cn.blockchain.GetBlockByNumber(uint64(start))
	}
	switch end {
	case rpc.PendingBlockNumber:
		to = api.cn.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		to = api.cn.blockchain.CurrentBlock()
	default:
		to = api.cn.blockchain.GetBlockByNumber(uint64(end))
	}
	// Trace the chain if we've found all our blocks
	if from == nil {
		return nil, fmt.Errorf("starting block #%d not found", start)
	}
	if to == nil {
		return nil, fmt.Errorf("end block #%d not found", end)
	}
	if from.Number().Cmp(to.Number()) >= 0 {
		return nil, fmt.Errorf("end block #%d needs to come after start block #%d", end, start)
	}
	return api.traceChain(ctx, from, to, config)
}

// traceChain configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requestd tracer.
func (api *PrivateDebugAPI) traceChain(ctx context.Context, start, end *types.Block, config *TraceConfig) (*rpc.Subscription, error) {
	// Tracing a chain is a **long** operation, only do with subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	sub := notifier.CreateSubscription()

	// Ensure we have a valid starting state before doing any work
	origin := start.NumberU64()
	database := state.NewDatabaseWithExistingCache(api.cn.ChainDB(), api.cn.blockchain.StateCache().TrieDB().TrieNodeCache()) // Chain tracing will probably start at genesis

	if number := start.NumberU64(); number > 0 {
		start = api.cn.blockchain.GetBlock(start.ParentHash(), start.NumberU64()-1)
		if start == nil {
			return nil, fmt.Errorf("parent block #%d not found", number-1)
		}
	}
	statedb, err := state.New(start.Root(), database, nil)
	if err != nil {
		// If the starting state is missing, allow some number of blocks to be reexecuted
		reexec := defaultTraceReexec
		if config != nil && config.Reexec != nil {
			reexec = *config.Reexec
		}
		// Find the most recent block that has the state available
		for i := uint64(0); i < reexec; i++ {
			start = api.cn.blockchain.GetBlock(start.ParentHash(), start.NumberU64()-1)
			if start == nil {
				break
			}
			if statedb, err = state.New(start.Root(), database, nil); err == nil {
				break
			}
		}
		// If we still don't have the state available, bail out
		if err != nil {
			switch err.(type) {
			case *statedb2.MissingNodeError:
				return nil, errors.New("required historical state unavailable")
			default:
				return nil, err
			}
		}
	}
	// Execute all the transaction contained within the chain concurrently for each block
	blocks := int(end.NumberU64() - origin)

	threads := runtime.NumCPU()
	if threads > blocks {
		threads = blocks
	}
	var (
		pend    = new(sync.WaitGroup)
		tasks   = make(chan *blockTraceTask, threads)
		results = make(chan *blockTraceTask, threads)
	)
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()

			// Fetch and execute the next block trace tasks
			for task := range tasks {
				signer := types.MakeSigner(api.config, task.block.Number())

				// Trace all the transactions contained within
				for i, tx := range task.block.Transactions() {
					msg, err := tx.AsMessageWithAccountKeyPicker(signer, task.statedb, task.block.NumberU64())
					if err != nil {
						logger.Warn("Tracing failed", "hash", tx.Hash(), "block", task.block.NumberU64(), "err", err)
						task.results[i] = &txTraceResult{TxHash: tx.Hash(), Error: err.Error()}
						break
					}

					vmctx := blockchain.NewEVMContext(msg, task.block.Header(), api.cn.blockchain, nil)

					res, err := api.traceTx(ctx, msg, vmctx, task.statedb, config)
					if err != nil {
						task.results[i] = &txTraceResult{TxHash: tx.Hash(), Error: err.Error()}
						logger.Warn("Tracing failed", "hash", tx.Hash(), "block", task.block.NumberU64(), "err", err)
						break
					}
					task.statedb.Finalise(true, true)
					task.results[i] = &txTraceResult{TxHash: tx.Hash(), Result: res}
				}
				// Stream the result back to the user or abort on teardown
				select {
				case results <- task:
				case <-notifier.Closed():
					return
				}
			}
		}()
	}
	// Start a goroutine to feed all the blocks into the tracers
	begin := time.Now()

	go func() {
		var (
			logged time.Time
			number uint64
			traced uint64
			failed error
			proot  common.Hash
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
		// Feed all the blocks both into the tracer, as well as fast process concurrently
		for number = start.NumberU64() + 1; number <= end.NumberU64(); number++ {
			// Stop tracing if interruption was requested
			select {
			case <-notifier.Closed():
				return
			default:
			}
			// Print progress logs if long enough time elapsed
			if time.Since(logged) > log.StatsReportLimit {
				if number > origin {
					nodeSize, preimageSize := database.TrieDB().Size()
					logger.Info("Tracing chain segment", "start", origin, "end", end.NumberU64(), "current", number, "transactions", traced, "elapsed", time.Since(begin), "nodeSize", nodeSize, "preimageSize", preimageSize)
				} else {
					logger.Info("Preparing state for chain trace", "block", number, "start", origin, "elapsed", time.Since(begin))
				}
				logged = time.Now()
			}
			// Retrieve the next block to trace
			block := api.cn.blockchain.GetBlockByNumber(number)
			if block == nil {
				failed = fmt.Errorf("block #%d not found", number)
				break
			}
			// Send the block over to the concurrent tracers (if not in the fast-forward phase)
			if number > origin {
				txs := block.Transactions()

				select {
				case tasks <- &blockTraceTask{statedb: statedb.Copy(), block: block, rootref: proot, results: make([]*txTraceResult, len(txs))}:
				case <-notifier.Closed():
					return
				}
				traced += uint64(len(txs))
			}
			// Generate the next state snapshot fast without tracing
			_, _, _, _, _, err := api.cn.blockchain.Processor().Process(block, statedb, vm.Config{UseOpcodeComputationCost: true})
			if err != nil {
				failed = err
				break
			}
			// Finalize the state so any modifications are written to the trie
			root, err := statedb.Commit(true)
			if err != nil {
				failed = err
				break
			}
			if err := statedb.Reset(root); err != nil {
				failed = err
				break
			}
			// Reference the trie twice, once for us, once for the trancer
			database.TrieDB().Reference(root, common.Hash{})
			if number >= origin {
				database.TrieDB().Reference(root, common.Hash{})
			}
			// Dereference all past tries we ourselves are done working with
			if !common.EmptyHash(proot) {
				database.TrieDB().Dereference(proot)
			}
			proot = root
		}
	}()

	// Keep reading the trace results and stream the to the user
	go func() {
		var (
			done = make(map[uint64]*blockTraceResult)
			next = origin + 1
		)
		for res := range results {
			// Queue up next received result
			result := &blockTraceResult{
				Block:  hexutil.Uint64(res.block.NumberU64()),
				Hash:   res.block.Hash(),
				Traces: res.results,
			}
			done[uint64(result.Block)] = result

			// Dereference any paret tries held in memory by this task
			database.TrieDB().Dereference(res.rootref)

			// Stream completed traces to the user, aborting on the first error
			for result, ok := done[next]; ok; result, ok = done[next] {
				if len(result.Traces) > 0 || next == end.NumberU64() {
					notifier.Notify(sub.ID, result)
				}
				delete(done, next)
				next++
			}
		}
	}()
	return sub, nil
}

// TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]*txTraceResult, error) {
	// Fetch the block that we want to trace
	var block *types.Block

	switch number {
	case rpc.PendingBlockNumber:
		block = api.cn.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		block = api.cn.blockchain.CurrentBlock()
	default:
		block = api.cn.blockchain.GetBlockByNumber(uint64(number))
	}
	// Trace the block if it was found
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", number)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockByHash(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	block := api.cn.blockchain.GetBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("block #%x not found", hash)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlock returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlock(ctx context.Context, blob hexutil.Bytes, config *TraceConfig) ([]*txTraceResult, error) {
	block := new(types.Block)
	if err := rlp.Decode(bytes.NewReader(blob), block); err != nil {
		return nil, fmt.Errorf("could not decode block: %v", err)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockFromFile returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockFromFile(ctx context.Context, file string, config *TraceConfig) ([]*txTraceResult, error) {
	blob, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	return api.TraceBlock(ctx, common.Hex2Bytes(string(blob)), config)
}

// TraceBadBlock returns the structured logs created during the execution of
// EVM against a block pulled from the pool of bad ones and returns them as a JSON
// object.
func (api *PrivateDebugAPI) TraceBadBlock(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	blocks, err := api.cn.blockchain.BadBlocks()
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if block.Hash == hash {
			return api.traceBlock(ctx, block.Block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}

// StandardTraceBlockToFile dumps the structured logs created during the
// execution of EVM to the local file system and returns a list of files
// to the caller.
func (api *PrivateDebugAPI) StandardTraceBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([]string, error) {
	block := api.cn.blockchain.GetBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("block %#x not found", hash)
	}
	return api.standardTraceBlockToFile(ctx, block, config)
}

// StandardTraceBadBlockToFile dumps the structured logs created during the
// execution of EVM against a block pulled from the pool of bad ones to the
// local file system and returns a list of files to the caller.
func (api *PrivateDebugAPI) StandardTraceBadBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([]string, error) {
	blocks, err := api.cn.blockchain.BadBlocks()
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		if block.Hash == hash {
			return api.standardTraceBlockToFile(ctx, block.Block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}

// traceBlock configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requestd tracer.
func (api *PrivateDebugAPI) traceBlock(ctx context.Context, block *types.Block, config *TraceConfig) ([]*txTraceResult, error) {
	// Create the parent state database
	if err := api.cn.engine.VerifyHeader(api.cn.blockchain, block.Header(), true); err != nil {
		return nil, err
	}
	parent := api.cn.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, fmt.Errorf("parent %x not found", block.ParentHash())
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}

	statedb, deferFn, err := api.stateAt(parent, reexec)
	defer deferFn()
	if err != nil {
		return nil, fmt.Errorf("can not get the state of block %#x: %v", parent.Root(), err)
	}

	// Execute all the transaction contained within the block concurrently
	var (
		signer = types.MakeSigner(api.config, block.Number())

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

				vmctx := blockchain.NewEVMContext(msg, block.Header(), api.cn.blockchain, nil)

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

		vmctx := blockchain.NewEVMContext(msg, block.Header(), api.cn.blockchain, nil)

		vmenv := vm.NewEVM(vmctx, statedb, api.config, &vm.Config{UseOpcodeComputationCost: true})
		if _, _, kerr := blockchain.ApplyMessage(vmenv, msg); kerr.ErrTxInvalid != nil {
			failed = kerr.ErrTxInvalid
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
func (api *PrivateDebugAPI) standardTraceBlockToFile(ctx context.Context, block *types.Block, config *StdTraceConfig) ([]string, error) {
	// If we're tracing a single transaction, make sure it's present
	if config != nil && !common.EmptyHash(config.TxHash) {
		if !containsTx(block, config.TxHash) {
			return nil, fmt.Errorf("transaction %#x not found in block", config.TxHash)
		}
	}
	// Create the parent state database
	if err := api.cn.engine.VerifyHeader(api.cn.blockchain, block.Header(), true); err != nil {
		return nil, err
	}
	parent := api.cn.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, fmt.Errorf("parent %#x not found", block.ParentHash())
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}

	statedb, deferFn, err := api.stateAt(parent, reexec)
	defer deferFn()
	if err != nil {
		return nil, fmt.Errorf("can not get the state of block %#x: %v", parent.Root(), err)
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
		signer = types.MakeSigner(api.config, block.Number())
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
			vmctx = blockchain.NewEVMContext(msg, block.Header(), api.cn.blockchain, nil)

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
		vmenv := vm.NewEVM(vmctx, statedb, api.config, &vmConf)
		_, _, kerr := blockchain.ApplyMessage(vmenv, msg)

		if dump != nil {
			dump.Close()
			logger.Info("Wrote standard trace", "file", dump.Name())
		}
		if kerr.ErrTxInvalid != nil {
			return dumps, kerr.ErrTxInvalid
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

// computeStateDB retrieves the state database associated with a certain block.
// A number of blocks are attempted to be reexecuted to generate the desired state.
func (api *PrivateDebugAPI) computeStateDB(block *types.Block, reexec uint64) (*state.StateDB, error) {
	// try to reexec blocks until we find a state or reach our limit
	origin := block.NumberU64()
	database := state.NewDatabaseWithExistingCache(api.cn.ChainDB(), api.cn.blockchain.StateCache().TrieDB().TrieNodeCache())

	var statedb *state.StateDB
	var err error

	for i := uint64(0); i < reexec; i++ {
		if statedb, err = state.New(block.Root(), database, nil); err == nil {
			break
		}
		blockNumber := block.NumberU64()
		block = api.cn.blockchain.GetBlock(block.ParentHash(), blockNumber-1)
		if block == nil {
			return nil, fmt.Errorf("block #%d not found", blockNumber-1)
		}
	}
	if err != nil {
		switch err.(type) {
		case *statedb2.MissingNodeError:
			return nil, fmt.Errorf("required historical state unavailable (reexec=%d)", reexec)
		default:
			return nil, err
		}
	}
	// State was available at historical point, regenerate
	var (
		start  = time.Now()
		logged time.Time
		proot  common.Hash
	)
	for block.NumberU64() < origin {
		// Print progress logs if long enough time elapsed
		if time.Since(logged) > log.StatsReportLimit {
			logger.Info("Regenerating historical state", "block", block.NumberU64()+1, "target", origin, "remaining", origin-block.NumberU64()-1, "elapsed", time.Since(start))
			logged = time.Now()
		}
		// Retrieve the next block to regenerate and process it
		if block = api.cn.blockchain.GetBlockByNumber(block.NumberU64() + 1); block == nil {
			return nil, fmt.Errorf("block #%d not found", block.NumberU64()+1)
		}
		_, _, _, _, _, err := api.cn.blockchain.Processor().Process(block, statedb, vm.Config{UseOpcodeComputationCost: true})
		if err != nil {
			return nil, fmt.Errorf("processing block %d failed: %v", block.NumberU64(), err)
		}
		// Finalize the state so any modifications are written to the trie
		root, err := statedb.Commit(true)
		if err != nil {
			return nil, err
		}
		if err := statedb.Reset(root); err != nil {
			return nil, fmt.Errorf("state reset after block %d failed: %v", block.NumberU64(), err)
		}
		database.TrieDB().Reference(root, common.Hash{})
		if !common.EmptyHash(proot) {
			database.TrieDB().Dereference(proot)
		}
		proot = root
	}
	nodeSize, preimageSize := database.TrieDB().Size()
	logger.Info("Historical state regenerated", "block", block.NumberU64(), "elapsed", time.Since(start), "nodeSize", nodeSize, "preimageSize", preimageSize)
	return statedb, nil
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	// Retrieve the transaction and assemble its EVM context
	tx, blockHash, _, index := api.cn.ChainDB().ReadTxAndLookupInfo(hash)
	if tx == nil {
		return nil, fmt.Errorf("transaction %#x not found", hash)
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	msg, vmctx, statedb, err := api.computeTxEnv(blockHash, int(index), reexec)
	if err != nil {
		return nil, err
	}
	// Trace the transaction and return
	return api.traceTx(ctx, msg, vmctx, statedb, config)
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func (api *PrivateDebugAPI) traceTx(ctx context.Context, message blockchain.Message, vmctx vm.Context, statedb *state.StateDB, config *TraceConfig) (interface{}, error) {
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
			// Constuct the JavaScript tracer to execute with
			if tracer, err = tracers.New(*config.Tracer); err != nil {
				return nil, err
			}
		}
		// Handle timeouts and RPC cancellations
		deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
		go func() {
			<-deadlineCtx.Done()
			switch t := tracer.(type) {
			case *tracers.Tracer:
				t.Stop(errors.New("execution timeout"))
			case *vm.InternalTxTracer:
				t.Stop(errors.New("execution timeout"))
			default:
				logger.Warn("unknown tracer type", "type", reflect.TypeOf(t).String())
			}
		}()
		defer cancel()

	case config == nil:
		tracer = vm.NewStructLogger(nil)

	default:
		tracer = vm.NewStructLogger(config.LogConfig)
	}
	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(vmctx, statedb, api.config, &vm.Config{Debug: true, Tracer: tracer, UseOpcodeComputationCost: true})

	ret, gas, kerr := blockchain.ApplyMessage(vmenv, message)
	if kerr.ErrTxInvalid != nil {
		return nil, fmt.Errorf("tracing failed: %v", kerr.ErrTxInvalid)
	}
	// Depending on the tracer type, format and return the output
	switch tracer := tracer.(type) {
	case *vm.StructLogger:
		return &klaytnapi.ExecutionResult{
			Gas:         gas,
			Failed:      kerr.Status != types.ReceiptStatusSuccessful,
			ReturnValue: fmt.Sprintf("%x", ret),
			StructLogs:  klaytnapi.FormatLogs(tracer.StructLogs()),
		}, nil

	case *tracers.Tracer:
		return tracer.GetResult()
	case *vm.InternalTxTracer:
		return tracer.GetResult()

	default:
		panic(fmt.Sprintf("bad tracer type %T", tracer))
	}
}

// stateAt returns the given block's state from cached stateDB with the GC lock or by regenerating.
func (api *PrivateDebugAPI) stateAt(block *types.Block, reexec uint64) (*state.StateDB, func(), error) {
	var stateDB *state.StateDB

	// If we have the state fully available in cachedNode, use that.
	stateDB, err := api.cn.blockchain.StateAtWithGCLock(block.Root())
	if err == nil {
		logger.Debug("Get stateDB from stateCache", "block", block.NumberU64())
		// During this processing, this lock will prevent to evict the state.
		return stateDB, stateDB.UnlockGCCachedNode, nil
	}

	emptyFn := func() {}

	// If we have the state fully available in persistent, use that.
	stateDB, err = api.cn.blockchain.StateAt(block.Root())
	if err == nil {
		logger.Debug("Get stateDB from persistent DB or its cache", "block", block.NumberU64())
		return stateDB, emptyFn, nil
	}

	// If no state is locally available, the desired state will be generated.
	stateDB, err = api.computeStateDB(block, reexec)
	if err == nil {
		logger.Debug("Get stateDB by computeStateDB", "block", block.NumberU64(), "reexec", reexec)
		return stateDB, emptyFn, nil
	}

	return nil, emptyFn, err
}

// computeTxEnv returns the execution environment of a certain transaction.
func (api *PrivateDebugAPI) computeTxEnv(blockHash common.Hash, txIndex int, reexec uint64) (blockchain.Message, vm.Context, *state.StateDB, error) {
	// Create the parent state database
	block := api.cn.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return nil, vm.Context{}, nil, fmt.Errorf("block %#x not found", blockHash)
	}
	parent := api.cn.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.Context{}, nil, fmt.Errorf("parent %#x not found", block.ParentHash())
	}

	statedb, deferFn, err := api.stateAt(parent, reexec)
	defer deferFn()
	if err != nil {
		return nil, vm.Context{}, nil, fmt.Errorf("can not get the state of block %#x: %v", parent.Root(), err)
	}

	// Recompute transactions up to the target index.
	signer := types.MakeSigner(api.config, block.Number())

	for idx, tx := range block.Transactions() {
		// Assemble the transaction call message and return if the requested offset
		msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		if err != nil {
			logger.Warn("ComputeTxEnv failed", "hash", tx.Hash(), "block", block.NumberU64(), "err", err)
			return nil, vm.Context{}, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}

		context := blockchain.NewEVMContext(msg, block.Header(), api.cn.blockchain, nil)
		if idx == txIndex {
			return msg, context, statedb, nil
		}
		// Not yet the searched for transaction, execute on top of the current state
		vmenv := vm.NewEVM(context, statedb, api.config, &vm.Config{UseOpcodeComputationCost: true})
		if _, _, kerr := blockchain.ApplyMessage(vmenv, msg); kerr.ErrTxInvalid != nil {
			return nil, vm.Context{}, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}
		// Ensure any modifications are committed to the state
		statedb.Finalise(true, true)
	}
	return nil, vm.Context{}, nil, fmt.Errorf("transaction index %d out of range for block %#x", txIndex, blockHash)
}
