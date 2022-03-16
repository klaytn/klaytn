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
// This file is derived from accounts/abi/bind/backends/simulated.go (2018/06/04).
// Modified and improved for the klaytn development.

package backends

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/bloombits"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

// This nil assignment ensures compile time that SimulatedBackend implements bind.ContractBackend.
var _ bind.ContractBackend = (*SimulatedBackend)(nil)

var (
	errBlockNumberUnsupported  = errors.New("simulatedBackend cannot access blocks other than the latest block")
	errBlockDoesNotExist       = errors.New("block does not exist in blockchain")
	errTransactionDoesNotExist = errors.New("transaction does not exist")
	errGasEstimationFailed     = errors.New("gas required exceeds allowance or always failing transaction")
)

// SimulatedBackend implements bind.ContractBackend, simulating a blockchain in
// the background. Its main purpose is to allow for easy testing of contract bindings.
// Simulated backend implements the following interfaces:
// ChainReader, ChainStateReader, ContractBackend, ContractCaller, ContractFilterer, ContractTransactor,
// DeployBackend, GasEstimator, GasPricer, LogFilterer, PendingContractCaller, TransactionReader, and TransactionSender
type SimulatedBackend struct {
	database   database.DBManager     // In memory database to store our testing data
	blockchain *blockchain.BlockChain // Klaytn blockchain to handle the consensus

	mu           sync.Mutex
	pendingBlock *types.Block   // Currently pending block that will be imported on request
	pendingState *state.StateDB // Currently pending state that will be the active on request

	events *filters.EventSystem // Event system for filtering log events live

	config *params.ChainConfig
}

// NewSimulatedBackendWithDatabase creates a new binding backend based on the given database
// and uses a simulated blockchain for testing purposes.
func NewSimulatedBackendWithDatabase(database database.DBManager, alloc blockchain.GenesisAlloc, cfg *params.ChainConfig) *SimulatedBackend {
	genesis := blockchain.Genesis{Config: cfg, Alloc: alloc}
	genesis.MustCommit(database)
	blockchain, _ := blockchain.NewBlockChain(database, nil, genesis.Config, gxhash.NewFaker(), vm.Config{})

	backend := &SimulatedBackend{
		database:   database,
		blockchain: blockchain,
		config:     genesis.Config,
		events:     filters.NewEventSystem(new(event.TypeMux), &filterBackend{database, blockchain}, false),
	}
	backend.rollback()
	return backend
}

// NewSimulatedBackendWithGasPrice creates a new binding backend using a simulated blockchain with a given unitPrice.
// for testing purposes.
func NewSimulatedBackendWithGasPrice(alloc blockchain.GenesisAlloc, unitPrice uint64) *SimulatedBackend {
	// Without changing `params.AllGxhashProtocolChanges`,
	// the copied config is used for no side effect of other tests
	cfg := *params.AllGxhashProtocolChanges
	cfg.UnitPrice = unitPrice
	return NewSimulatedBackendWithDatabase(database.NewMemoryDBManager(), alloc, &cfg)
}

// NewSimulatedBackend creates a new binding backend using a simulated blockchain
// for testing purposes.
func NewSimulatedBackend(alloc blockchain.GenesisAlloc) *SimulatedBackend {
	return NewSimulatedBackendWithDatabase(database.NewMemoryDBManager(), alloc, params.AllGxhashProtocolChanges)
}

// Close terminates the underlying blockchain's update loop.
func (b *SimulatedBackend) Close() error {
	b.blockchain.Stop()
	return nil
}

// Commit imports all the pending transactions as a single block and starts a
// fresh new state.
func (b *SimulatedBackend) Commit() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, err := b.blockchain.InsertChain([]*types.Block{b.pendingBlock}); err != nil {
		panic(err) // This cannot happen unless the simulator is wrong, fail in that case
	}
	b.rollback()
}

// Rollback aborts all pending transactions, reverting to the last committed state.
func (b *SimulatedBackend) Rollback() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.rollback()
}

func (b *SimulatedBackend) rollback() {
	blocks, _ := blockchain.GenerateChain(b.config, b.blockchain.CurrentBlock(), gxhash.NewFaker(), b.database, 1, func(int, *blockchain.BlockGen) {})
	stateDB, _ := b.blockchain.State()

	b.pendingBlock = blocks[0]
	b.pendingState, _ = state.New(b.pendingBlock.Root(), stateDB.Database(), nil)
}

// stateByBlockNumber retrieves a state by a given blocknumber.
func (b *SimulatedBackend) stateByBlockNumber(ctx context.Context, blockNumber *big.Int) (*state.StateDB, error) {
	if blockNumber == nil || blockNumber.Cmp(b.blockchain.CurrentBlock().Number()) == 0 {
		return b.blockchain.State()
	}
	block, err := b.blockByNumberNoLock(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	return b.blockchain.StateAt(block.Root())
}

// CodeAt returns the code associated with a certain account in the blockchain.
func (b *SimulatedBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	return stateDB.GetCode(contract), nil
}

// BalanceAt returns the peb balance of a certain account in the blockchain.
func (b *SimulatedBackend) BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	return stateDB.GetBalance(contract), nil
}

// NonceAt returns the nonce of a certain account in the blockchain.
func (b *SimulatedBackend) NonceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return 0, err
	}

	return stateDB.GetNonce(contract), nil
}

// StorageAt returns the value of key in the storage of an account in the blockchain.
func (b *SimulatedBackend) StorageAt(ctx context.Context, contract common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	val := stateDB.GetState(contract, key)
	return val[:], nil
}

// TransactionReceipt returns the receipt of a transaction.
func (b *SimulatedBackend) TransactionReceipt(_ context.Context, txHash common.Hash) (*types.Receipt, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	receipt, _, _, _ := b.database.ReadReceipt(txHash)
	return receipt, nil
}

// TransactionByHash checks the pool of pending transactions in addition to the
// blockchain. The isPending return value indicates whether the transaction has been
// mined yet. Note that the transaction may not be part of the canonical chain even if
// it's not pending.
func (b *SimulatedBackend) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	tx := b.pendingBlock.Transaction(txHash)
	if tx != nil {
		return tx, true, nil
	}
	tx, _, _, _ = b.database.ReadTxAndLookupInfo(txHash)
	if tx != nil {
		return tx, false, nil
	}
	return nil, false, klaytn.NotFound
}

// BlockByHash retrieves a block based on the block hash.
func (b *SimulatedBackend) BlockByHash(_ context.Context, hash common.Hash) (*types.Block, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if hash == b.pendingBlock.Hash() {
		return b.pendingBlock, nil
	}

	block := b.blockchain.GetBlockByHash(hash)
	if block != nil {
		return block, nil
	}

	return nil, errBlockDoesNotExist
}

// BlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (b *SimulatedBackend) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.blockByNumberNoLock(ctx, number)
}

// blockByNumberNoLock retrieves a block from the database by number, caching it
// (associated with its hash) if found without Lock.
func (b *SimulatedBackend) blockByNumberNoLock(_ context.Context, number *big.Int) (*types.Block, error) {
	if number == nil || number.Cmp(b.pendingBlock.Number()) == 0 {
		return b.blockchain.CurrentBlock(), nil
	}

	block := b.blockchain.GetBlockByNumber(uint64(number.Int64()))
	if block == nil {
		return nil, errBlockDoesNotExist
	}

	return block, nil
}

// HeaderByHash returns a block header from the current canonical chain.
func (b *SimulatedBackend) HeaderByHash(_ context.Context, hash common.Hash) (*types.Header, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if hash == b.pendingBlock.Hash() {
		return b.pendingBlock.Header(), nil
	}

	header := b.blockchain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errBlockDoesNotExist
	}

	return header, nil
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (b *SimulatedBackend) HeaderByNumber(_ context.Context, block *big.Int) (*types.Header, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if block == nil || block.Cmp(b.pendingBlock.Number()) == 0 {
		return b.blockchain.CurrentHeader(), nil
	}

	return b.blockchain.GetHeaderByNumber(uint64(block.Int64())), nil
}

// TransactionCount returns the number of transactions in a given block.
func (b *SimulatedBackend) TransactionCount(_ context.Context, blockHash common.Hash) (uint, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if blockHash == b.pendingBlock.Hash() {
		return uint(b.pendingBlock.Transactions().Len()), nil
	}

	block := b.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return uint(0), errBlockDoesNotExist
	}

	return uint(block.Transactions().Len()), nil
}

// TransactionInBlock returns the transaction for a specific block at a specific index.
func (b *SimulatedBackend) TransactionInBlock(_ context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if blockHash == b.pendingBlock.Hash() {
		transactions := b.pendingBlock.Transactions()
		if uint(len(transactions)) < index+1 {
			return nil, errTransactionDoesNotExist
		}

		return transactions[index], nil
	}

	block := b.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return nil, errBlockDoesNotExist
	}

	transactions := block.Transactions()
	if uint(len(transactions)) < index+1 {
		return nil, errTransactionDoesNotExist
	}

	return transactions[index], nil
}

// PendingCodeAt returns the code associated with an account in the pending state.
func (b *SimulatedBackend) PendingCodeAt(_ context.Context, contract common.Address) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.pendingState.GetCode(contract), nil
}

// CallContract executes a contract call.
func (b *SimulatedBackend) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if blockNumber != nil && blockNumber.Cmp(b.blockchain.CurrentBlock().Number()) != 0 {
		return nil, errBlockNumberUnsupported
	}
	stateDB, err := b.blockchain.State()
	if err != nil {
		return nil, err
	}
	res, _, _, err := b.callContract(ctx, call, b.blockchain.CurrentBlock(), stateDB)
	return res, err
}

// PendingCallContract executes a contract call on the pending state.
func (b *SimulatedBackend) PendingCallContract(ctx context.Context, call klaytn.CallMsg) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	defer b.pendingState.RevertToSnapshot(b.pendingState.Snapshot())

	rval, _, _, err := b.callContract(ctx, call, b.pendingBlock, b.pendingState)
	return rval, err
}

// PendingNonceAt implements PendingStateReader.PendingNonceAt, retrieving
// the nonce currently pending for the account.
func (b *SimulatedBackend) PendingNonceAt(_ context.Context, account common.Address) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.pendingState.GetOrNewStateObject(account).Nonce(), nil
}

// SuggestGasPrice implements ContractTransactor.SuggestGasPrice. Since the simulated
// chain doesn't have miners, we just return a gas price of 1 for any call.
func (b *SimulatedBackend) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	return new(big.Int).SetUint64(b.config.UnitPrice), nil
}

// EstimateGas executes the requested code against the latest block/state and
// returns the used amount of gas.
func (b *SimulatedBackend) EstimateGas(ctx context.Context, call klaytn.CallMsg) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Determine the lowest and highest possible gas limits to binary search in between
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	if call.Gas >= params.TxGas {
		hi = call.Gas
	} else {
		hi = params.UpperGasLimit
	}

	// Recap the highest gas allowance with account's balance.
	if call.GasPrice != nil && call.GasPrice.BitLen() != 0 {
		balance := b.pendingState.GetBalance(call.From) // from can't be nil
		available := new(big.Int).Set(balance)
		if call.Value != nil {
			if call.Value.Cmp(available) >= 0 {
				return 0, errors.New("insufficient funds for transfer")
			}
			available.Sub(available, call.Value)
		}
		allowance := new(big.Int).Div(available, call.GasPrice)
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := call.Value
			if transfer == nil {
				transfer = new(big.Int)
			}
			bind.Logger.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", transfer, "gasprice", call.GasPrice, "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) bool {
		call.Gas = gas

		currentState, err := b.blockchain.State()
		if err != nil {
			return false
		}
		_, _, failed, err := b.callContract(ctx, call, b.blockchain.CurrentBlock(), currentState)
		if err != nil || failed {
			return false
		}
		return true
	}
	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		if !executable(mid) {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		if !executable(hi) {
			return 0, errGasEstimationFailed
		}
	}
	return hi, nil
}

// callContract implements common code between normal and pending contract calls.
// state is modified during execution, make sure to copy it if necessary.
func (b *SimulatedBackend) callContract(_ context.Context, call klaytn.CallMsg, block *types.Block, stateDB *state.StateDB) ([]byte, uint64, bool, error) {
	// Ensure message is initialized properly.
	if call.GasPrice == nil {
		call.GasPrice = big.NewInt(1)
	}
	if call.Gas == 0 {
		call.Gas = 50000000
	}
	if call.Value == nil {
		call.Value = new(big.Int)
	}
	// Set infinite balance to the fake caller account.
	from := stateDB.GetOrNewStateObject(call.From)
	from.SetBalance(math.MaxBig256)
	// Execute the call.
	nonce := from.Nonce()
	intrinsicGas, _ := types.IntrinsicGas(call.Data, nil, call.To == nil, b.config.Rules(block.Number()))
	msg := types.NewMessage(call.From, call.To, nonce, call.Value, call.Gas, call.GasPrice, call.Data, true, intrinsicGas)

	evmContext := blockchain.NewEVMContext(msg, block.Header(), b.blockchain, nil)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(evmContext, stateDB, b.config, &vm.Config{})

	ret, usedGas, kerr := blockchain.NewStateTransition(vmenv, msg).TransitionDb()

	// Propagate error of Receipt
	err := kerr.ErrTxInvalid
	if err == nil {
		err = blockchain.GetVMerrFromReceiptStatus(kerr.Status)
	}

	return ret, usedGas, kerr.Status != types.ReceiptStatusSuccessful, err
}

// SendTransaction updates the pending block to include the given transaction.
// It panics if the transaction is invalid.
func (b *SimulatedBackend) SendTransaction(_ context.Context, tx *types.Transaction) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check transaction validity
	block := b.blockchain.CurrentBlock()
	signer := types.MakeSigner(b.blockchain.Config(), block.Number())
	sender, err := types.Sender(signer, tx)

	if err != nil {
		panic(fmt.Errorf("invalid transaction: %v", err))
	}
	nonce := b.pendingState.GetNonce(sender)
	if tx.Nonce() != nonce {
		panic(fmt.Errorf("invalid transaction nonce: got %d, want %d", tx.Nonce(), nonce))
	}

	// Include tx in chain.
	blocks, _ := blockchain.GenerateChain(b.config, block, gxhash.NewFaker(), b.database, 1, func(number int, block *blockchain.BlockGen) {
		for _, tx := range b.pendingBlock.Transactions() {
			block.AddTxWithChain(b.blockchain, tx)
		}
		block.AddTxWithChain(b.blockchain, tx)
	})
	stateDB, _ := b.blockchain.State()

	b.pendingBlock = blocks[0]
	b.pendingState, _ = state.New(b.pendingBlock.Root(), stateDB.Database(), nil)
	return nil
}

// FilterLogs executes a log filter operation, blocking during execution and
// returning all the results in one batch.
//
// TODO(karalabe): Deprecate when the subscription one can return past data too.
func (b *SimulatedBackend) FilterLogs(ctx context.Context, query klaytn.FilterQuery) ([]types.Log, error) {
	// Initialize unset filter boundaries to run from genesis to chain head
	from := int64(0)
	if query.FromBlock != nil {
		from = query.FromBlock.Int64()
	}
	to := int64(-1)
	if query.ToBlock != nil {
		to = query.ToBlock.Int64()
	}
	// Construct and execute the filter
	filter := filters.NewRangeFilter(&filterBackend{b.database, b.blockchain}, from, to, query.Addresses, query.Topics)

	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]types.Log, len(logs))
	for i, nLog := range logs {
		res[i] = *nLog
	}
	return res, nil
}

// ChainID can return the chain ID of the chain.
func (b *SimulatedBackend) ChainID(ctx context.Context) (*big.Int, error) {
	return b.blockchain.Config().ChainID, nil
}

// SubscribeFilterLogs creates a background log filtering operation, returning a
// subscription immediately, which can be used to stream the found events.
func (b *SimulatedBackend) SubscribeFilterLogs(_ context.Context, query klaytn.FilterQuery, ch chan<- types.Log) (klaytn.Subscription, error) {
	// Subscribe to contract events
	sink := make(chan []*types.Log)

	sub, err := b.events.SubscribeLogs(query, sink)
	if err != nil {
		return nil, err
	}
	// Since we're getting logs in batches, we need to flatten them into a plain stream
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case logs := <-sink:
				for _, nlog := range logs {
					select {
					case ch <- *nlog:
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// CurrentBlockNumber returns a current block number.
func (b *SimulatedBackend) CurrentBlockNumber(ctx context.Context) (uint64, error) {
	return b.blockchain.CurrentBlock().NumberU64(), nil
}

// SubscribeNewHead returns an event subscription for a new header
func (b *SimulatedBackend) SubscribeNewHead(_ context.Context, ch chan<- *types.Header) (klaytn.Subscription, error) {
	// subscribe to a new head
	sink := make(chan *types.Header)
	sub := b.events.SubscribeNewHeads(sink)

	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case head := <-sink:
				select {
				case ch <- head:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// AdjustTime adds a time shift to the simulated clock.
// It can only be called on empty blocks.
func (b *SimulatedBackend) AdjustTime(adjustment time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.pendingBlock.Transactions()) != 0 {
		return errors.New("Could not adjust time on non-empty block")
	}

	blocks, _ := blockchain.GenerateChain(b.config, b.blockchain.CurrentBlock(), gxhash.NewFaker(), b.database, 1, func(number int, block *blockchain.BlockGen) {
		block.OffsetTime(int64(adjustment.Seconds()))
	})
	stateDB, _ := b.blockchain.State()

	b.pendingBlock = blocks[0]
	b.pendingState, _ = state.New(b.pendingBlock.Root(), stateDB.Database(), nil)

	return nil
}

// Blockchain returns the underlying blockchain.
func (b *SimulatedBackend) BlockChain() *blockchain.BlockChain {
	return b.blockchain
}

func (b *SimulatedBackend) PendingBlock() *types.Block {
	return b.pendingBlock
}

// filterBackend implements filters.Backend to support filtering for logs without
// taking bloom-bits acceleration structures into account.
type filterBackend struct {
	db database.DBManager
	bc *blockchain.BlockChain
}

func (fb *filterBackend) ChainDB() database.DBManager { return fb.db }
func (fb *filterBackend) EventMux() *event.TypeMux    { panic("not supported") }

func (fb *filterBackend) HeaderByNumber(_ context.Context, block rpc.BlockNumber) (*types.Header, error) {
	if block == rpc.LatestBlockNumber {
		return fb.bc.CurrentHeader(), nil
	}
	return fb.bc.GetHeaderByNumber(uint64(block.Int64())), nil
}

func (fb *filterBackend) GetBlockReceipts(_ context.Context, hash common.Hash) types.Receipts {
	return fb.bc.GetReceiptsByBlockHash(hash)
}

func (fb *filterBackend) GetLogs(_ context.Context, hash common.Hash) ([][]*types.Log, error) {
	return fb.bc.GetLogsByHash(hash), nil
}

func (fb *filterBackend) SubscribeNewTxsEvent(_ chan<- blockchain.NewTxsEvent) event.Subscription {
	return nullSubscription()
}

func (fb *filterBackend) SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription {
	return fb.bc.SubscribeChainEvent(ch)
}

func (fb *filterBackend) SubscribeRemovedLogsEvent(ch chan<- blockchain.RemovedLogsEvent) event.Subscription {
	return fb.bc.SubscribeRemovedLogsEvent(ch)
}

func (fb *filterBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return fb.bc.SubscribeLogsEvent(ch)
}

func (fb *filterBackend) SubscribePendingLogsEvent(_ chan<- []*types.Log) event.Subscription {
	return nullSubscription()
}

func (fb *filterBackend) BloomStatus() (uint64, uint64) { return 4096, 0 }

func (fb *filterBackend) ServiceFilter(_ context.Context, _ *bloombits.MatcherSession) {
	panic("not supported")
}

func (fb *filterBackend) ChainConfig() *params.ChainConfig {
	return fb.bc.Config()
}

func nullSubscription() event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		<-quit
		return nil
	})
}
