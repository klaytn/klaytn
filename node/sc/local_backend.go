// Modifications Copyright 2019 The klaytn Authors
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

package sc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/bloombits"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

const defaultGasPrice = 50 * params.Ston

var errBlockNumberUnsupported = errors.New("LocalBackend cannot access blocks other than the latest block")
var errGasEstimationFailed = errors.New("gas required exceeds allowance or always failing transaction")

// TODO-Klaytn currently LocalBackend is only for ServiceChain, especially Bridge SmartContract
type LocalBackend struct {
	subbridge *SubBridge

	events *filters.EventSystem // Event system for filtering log events live
	config *params.ChainConfig
}

func checkCtx(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func NewLocalBackend(main *SubBridge) (*LocalBackend, error) {
	return &LocalBackend{
		subbridge: main,
		config:    main.blockchain.Config(),
		events:    filters.NewEventSystem(main.EventMux(), &filterLocalBackend{main}, false),
	}, nil
}

func (lb *LocalBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	if blockNumber != nil && blockNumber.Cmp(lb.subbridge.blockchain.CurrentBlock().Number()) != 0 {
		return nil, errBlockNumberUnsupported
	}
	statedb, err := lb.subbridge.blockchain.State()
	if err != nil {
		return nil, err
	}
	return statedb.GetCode(contract), nil
}

func (lb *LocalBackend) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	if blockNumber != nil && blockNumber.Cmp(lb.subbridge.blockchain.CurrentBlock().Number()) != 0 {
		return nil, errBlockNumberUnsupported
	}
	statedb, err := lb.subbridge.blockchain.State()
	if err != nil {
		return nil, err
	}
	return statedb.GetBalance(account), nil
}

func (lb *LocalBackend) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	if blockNumber != nil && blockNumber.Cmp(lb.subbridge.blockchain.CurrentBlock().Number()) != 0 {
		return nil, errBlockNumberUnsupported
	}
	currentState, err := lb.subbridge.blockchain.State()
	if err != nil {
		return nil, err
	}
	rval, _, _, err := lb.callContract(ctx, call, lb.subbridge.blockchain.CurrentBlock(), currentState)
	return rval, err
}

func (b *LocalBackend) callContract(ctx context.Context, call klaytn.CallMsg, block *types.Block, statedb *state.StateDB) ([]byte, uint64, bool, error) {
	// Set default gas & gas price if none were set
	gas, gasPrice := uint64(call.Gas), call.GasPrice
	if gas == 0 {
		gas = math.MaxUint64 / 2
	}
	if gasPrice == nil || gasPrice.Sign() == 0 {
		gasPrice = new(big.Int).SetUint64(defaultGasPrice)
	}

	intrinsicGas, err := types.IntrinsicGas(call.Data, nil, call.To == nil, b.config.Rules(block.Number()))
	if err != nil {
		return nil, 0, false, err
	}

	// Create new call message
	msg := types.NewMessage(call.From, call.To, 0, call.Value, gas, gasPrice, call.Data, false, intrinsicGas)

	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	statedb.SetBalance(msg.ValidatedSender(), math.MaxBig256)
	vmError := func() error { return nil }

	context := blockchain.NewEVMContext(msg, block.Header(), b.subbridge.blockchain, nil)
	evm := vm.NewEVM(context, statedb, b.config, &vm.Config{})
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel(vm.CancelByCtxDone)
	}()

	res, gas, kerr := blockchain.ApplyMessage(evm, msg)
	err = kerr.ErrTxInvalid
	if err := vmError(); err != nil {
		return nil, 0, false, err
	}

	// Propagate error of Receipt as JSON RPC error
	if err == nil {
		err = blockchain.GetVMerrFromReceiptStatus(kerr.Status)
	}

	return res, gas, kerr.Status != types.ReceiptStatusSuccessful, err
}

func (lb *LocalBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	// TODO-Klaytn this is not pending code but latest code
	return lb.CodeAt(ctx, contract, lb.subbridge.blockchain.CurrentBlock().Number())
}

func (lb *LocalBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if err := checkCtx(ctx); err != nil {
		return 0, err
	}
	return lb.subbridge.txPool.GetPendingNonce(account), nil
}

func (lb *LocalBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(lb.config.UnitPrice), nil
}

func (lb *LocalBackend) EstimateGas(ctx context.Context, call klaytn.CallMsg) (gas uint64, err error) {
	if err := checkCtx(ctx); err != nil {
		return 0, err
	}
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	if uint64(call.Gas) >= params.TxGas {
		hi = uint64(call.Gas)
	} else {
		hi = params.UpperGasLimit
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) bool {
		call.Gas = gas

		currentState, err := lb.subbridge.blockchain.State()
		if err != nil {
			return false
		}
		_, _, failed, err := lb.callContract(ctx, call, lb.subbridge.blockchain.CurrentBlock(), currentState)
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
			return 0, fmt.Errorf("gas required exceeds allowance or always failing transaction")
		}
	}
	return hi, nil
}

func (lb *LocalBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := checkCtx(ctx); err != nil {
		return err
	}
	return lb.subbridge.txPool.AddLocal(tx)
}

// ChainID can return the chain ID of the chain.
func (lb *LocalBackend) ChainID(ctx context.Context) (*big.Int, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	return lb.config.ChainID, nil
}

func (lb *LocalBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	receipt := lb.subbridge.blockchain.GetReceiptByTxHash(txHash)
	if receipt != nil {
		return receipt, nil
	}
	return nil, errors.New("receipt is not exist")
}

func (lb *LocalBackend) FilterLogs(ctx context.Context, query klaytn.FilterQuery) ([]types.Log, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	// Convert the RPC block numbers into internal representations
	if query.FromBlock == nil {
		query.FromBlock = big.NewInt(rpc.LatestBlockNumber.Int64())
	}
	if query.ToBlock == nil {
		query.ToBlock = big.NewInt(rpc.LatestBlockNumber.Int64())
	}
	from := query.FromBlock.Int64()
	to := query.ToBlock.Int64()

	// Construct and execute the filter
	filter := filters.NewRangeFilter(&filterLocalBackend{lb.subbridge}, from, to, query.Addresses, query.Topics)

	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]types.Log, len(logs))
	for i, log := range logs {
		res[i] = *log
	}
	return res, nil
}

func (lb *LocalBackend) SubscribeFilterLogs(ctx context.Context, query klaytn.FilterQuery, ch chan<- types.Log) (klaytn.Subscription, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	// Subscribe to contract events
	sink := make(chan []*types.Log)

	sub, err := lb.events.SubscribeLogs(query, sink)
	if err != nil {
		return nil, err
	}
	// Since we're getting logs in batches, we need to flatten them into a plain stream
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case logs := <-sink:
				for _, log := range logs {
					select {
					case ch <- *log:
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
func (lb *LocalBackend) CurrentBlockNumber(ctx context.Context) (uint64, error) {
	if err := checkCtx(ctx); err != nil {
		return 0, err
	}
	return lb.subbridge.blockchain.CurrentBlock().NumberU64(), nil
}

type filterLocalBackend struct {
	subbridge *SubBridge
}

func (fb *filterLocalBackend) ChainDB() database.DBManager {
	// TODO-Klaytn consider chain's chainDB instead of bridge's chainDB currently.
	return fb.subbridge.chainDB
}
func (fb *filterLocalBackend) EventMux() *event.TypeMux {
	// TODO-Klaytn consider chain's eventMux instead of bridge's eventMux currently.
	return fb.subbridge.EventMux()
}

func (fb *filterLocalBackend) HeaderByNumber(ctx context.Context, block rpc.BlockNumber) (*types.Header, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	// TODO-Klaytn consider pendingblock instead of latest block
	if block == rpc.LatestBlockNumber {
		return fb.subbridge.blockchain.CurrentHeader(), nil
	}
	return fb.subbridge.blockchain.GetHeaderByNumber(uint64(block.Int64())), nil
}

func (fb *filterLocalBackend) GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts {
	if err := checkCtx(ctx); err != nil {
		return nil
	}
	return fb.subbridge.blockchain.GetReceiptsByBlockHash(hash)
}

func (fb *filterLocalBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	return fb.subbridge.blockchain.GetLogsByHash(hash), nil
}

func (fb *filterLocalBackend) SubscribeNewTxsEvent(ch chan<- blockchain.NewTxsEvent) event.Subscription {
	return fb.subbridge.txPool.SubscribeNewTxsEvent(ch)
}

func (fb *filterLocalBackend) SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription {
	return fb.subbridge.blockchain.SubscribeChainEvent(ch)
}

func (fb *filterLocalBackend) SubscribeRemovedLogsEvent(ch chan<- blockchain.RemovedLogsEvent) event.Subscription {
	return fb.subbridge.blockchain.SubscribeRemovedLogsEvent(ch)
}

func (fb *filterLocalBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return fb.subbridge.blockchain.SubscribeLogsEvent(ch)
}

func (fb *filterLocalBackend) BloomStatus() (uint64, uint64) {
	// TODO-Klaytn consider this number of sections.
	// BloomBitsBlocks (const : 4096), the number of processed sections maintained by the chain indexer
	return 4096, 0
}

func (fb *filterLocalBackend) ServiceFilter(_dummyCtx context.Context, session *bloombits.MatcherSession) {
	// TODO-Klaytn this method should implmentation to support indexed tag in solidity
	//for i := 0; i < bloomFilterThreads; i++ {
	//	go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, backend.bloomRequests)
	//}
}

func (fb *filterLocalBackend) ChainConfig() *params.ChainConfig {
	return fb.subbridge.blockchain.Config()
}
