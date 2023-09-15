// Copyright 2023 The klaytn Authors
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

package backends

import (
	"context"
	"errors"
	"math/big"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/params"
)

// Maintain separate minimal interfaces of blockchain.BlockChain because ContractBackend are used
// in various situations. BlockChain instances are often passed down as different interfaces such as
// consensus.ChainReader, governance.blockChain, work.BlockChain.
type BlockChainForCaller interface {
	// Required by NewEVMContext
	blockchain.ChainContext

	// Below is a subset of consensus.ChainReader
	// Only using the vocabulary of consensus.ChainReader for potential
	// usability within consensus package.
	Config() *params.ChainConfig
	GetHeaderByNumber(number uint64) *types.Header
	GetBlock(hash common.Hash, number uint64) *types.Block
	State() (*state.StateDB, error)
	StateAt(root common.Hash) (*state.StateDB, error)
	CurrentBlock() *types.Block
}

// BlockchainContractBackend implements bind.Contract* and bind.DeployBackend, based on
// a user-supplied blockchain.BlockChain instance.
// Its intended purpose is reading system contracts during block processing.
//
// Note that SimulatedBackend creates a new temporary BlockChain for testing,
// whereas BlockchainContractBackend uses an existing BlockChain with existing database.
type BlockchainContractBackend struct {
	bc     BlockChainForCaller
	txPool *blockchain.TxPool
	events *filters.EventSystem
}

// This nil assignment ensures at compile time that BlockchainContractBackend implements bind.Contract* and bind.DeployBackend.
var (
	_ bind.ContractCaller     = (*BlockchainContractBackend)(nil)
	_ bind.ContractTransactor = (*BlockchainContractBackend)(nil)
	_ bind.ContractFilterer   = (*BlockchainContractBackend)(nil)
	_ bind.DeployBackend      = (*BlockchainContractBackend)(nil)
	_ bind.ContractBackend    = (*BlockchainContractBackend)(nil)
)

// `txPool` is required for bind.ContractTransactor methods and `events` is required for bind.ContractFilterer methods.
// If `tp=nil`, bind.ContractTransactor methods could return errors.
// If `es=nil`, bind.ContractFilterer methods could return errors.
func NewBlockchainContractBackend(bc BlockChainForCaller, tp *blockchain.TxPool, es *filters.EventSystem) *BlockchainContractBackend {
	return &BlockchainContractBackend{bc, tp, es}
}

// bind.ContractCaller defined methods

func (b *BlockchainContractBackend) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	if _, state, err := b.getBlockAndState(blockNumber); err != nil {
		return nil, err
	} else {
		return state.GetCode(account), nil
	}
}

// Executes a read-only function call with respect to the specified block's state, or latest state if not specified.
//
// Returns call result in []byte.
// Returns error when:
// - cannot find the corresponding block or stateDB
// - VM revert error
// - VM other errors (e.g. NotProgramAccount, OutOfGas)
// - Error outside VM
func (b *BlockchainContractBackend) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	block, state, err := b.getBlockAndState(blockNumber)
	if err != nil {
		return nil, err
	}

	res, err := b.callContract(call, block, state)
	if err != nil {
		return nil, err
	}
	if len(res.Revert()) > 0 {
		return nil, blockchain.NewRevertError(res)
	}
	return res.Return(), res.Unwrap()
}

func (b *BlockchainContractBackend) callContract(call klaytn.CallMsg, block *types.Block, state *state.StateDB) (*blockchain.ExecutionResult, error) {
	if call.Gas == 0 {
		call.Gas = uint64(3e8) // enough gas for ordinary contract calls
	}

	intrinsicGas, err := types.IntrinsicGas(call.Data, nil, call.To == nil, b.bc.Config().Rules(block.Number()))
	if err != nil {
		return nil, err
	}

	msg := types.NewMessage(call.From, call.To, 0, call.Value, call.Gas, call.GasPrice, call.Data,
		false, intrinsicGas)

	txContext := blockchain.NewEVMTxContext(msg, block.Header())
	blockContext := blockchain.NewEVMBlockContext(block.Header(), b.bc, nil)

	// EVM demands the sender to have enough KLAY balance (gasPrice * gasLimit) in buyGas()
	// After KIP-71, gasPrice is nonzero baseFee, regardless of the msg.gasPrice (usually 0)
	// But our sender (usually 0x0) won't have enough balance. Instead we override gasPrice = 0 here
	txContext.GasPrice = big.NewInt(0)
	evm := vm.NewEVM(blockContext, txContext, state, b.bc.Config(), &vm.Config{})

	return blockchain.ApplyMessage(evm, msg)
}

func (b *BlockchainContractBackend) getBlockAndState(num *big.Int) (*types.Block, *state.StateDB, error) {
	var block *types.Block
	if num == nil {
		block = b.bc.CurrentBlock()
	} else {
		header := b.bc.GetHeaderByNumber(num.Uint64())
		if header == nil {
			return nil, nil, errBlockDoesNotExist
		}
		block = b.bc.GetBlock(header.Hash(), header.Number.Uint64())
	}
	if block == nil {
		return nil, nil, errBlockDoesNotExist
	}

	state, err := b.bc.StateAt(block.Root())
	return block, state, err
}

// bind.ContractTransactor defined methods

func (b *BlockchainContractBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	// TODO-Klaytn this is not pending code but latest code
	state, err := b.bc.State()
	if err != nil {
		return nil, err
	}
	return state.GetCode(account), nil
}

func (b *BlockchainContractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if b.txPool != nil {
		return b.txPool.GetPendingNonce(account), nil
	}
	// TODO-Klaytn this is not pending nonce but latest nonce
	state, err := b.bc.State()
	if err != nil {
		return 0, err
	}
	return state.GetNonce(account), nil
}

func (b *BlockchainContractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if b.bc.Config().IsMagmaForkEnabled(b.bc.CurrentBlock().Number()) {
		if b.txPool != nil {
			return new(big.Int).Mul(b.txPool.GasPrice(), big.NewInt(2)), nil
		} else {
			return new(big.Int).Mul(b.bc.CurrentBlock().Header().BaseFee, big.NewInt(2)), nil
		}
	} else {
		return new(big.Int).SetUint64(b.bc.Config().UnitPrice), nil
	}
}

func (b *BlockchainContractBackend) EstimateGas(ctx context.Context, call klaytn.CallMsg) (uint64, error) {
	state, err := b.bc.State()
	if err != nil {
		return 0, err
	}
	balance := state.GetBalance(call.From) // from can't be nil

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) (bool, *blockchain.ExecutionResult, error) {
		call.Gas = gas

		currentState, err := b.bc.State()
		if err != nil {
			return true, nil, nil
		}
		res, err := b.callContract(call, b.bc.CurrentBlock(), currentState)
		if err != nil {
			if errors.Is(err, blockchain.ErrIntrinsicGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			return true, nil, err // Bail out
		}
		return res.Failed(), res, nil
	}

	estimated, err := blockchain.DoEstimateGas(ctx, call.Gas, 0, call.Value, call.GasPrice, balance, executable)
	return uint64(estimated), err
}

func (b *BlockchainContractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if b.txPool == nil {
		return errors.New("tx pool not configured")
	}
	return b.txPool.AddLocal(tx)
}

func (b *BlockchainContractBackend) ChainID(ctx context.Context) (*big.Int, error) {
	return b.bc.Config().ChainID, nil
}

// bind.ContractFilterer defined methods

func (b *BlockchainContractBackend) FilterLogs(ctx context.Context, query klaytn.FilterQuery) ([]types.Log, error) {
	// Convert the current block numbers into internal representations
	if query.FromBlock == nil {
		query.FromBlock = big.NewInt(b.bc.CurrentBlock().Number().Int64())
	}
	if query.ToBlock == nil {
		query.ToBlock = big.NewInt(b.bc.CurrentBlock().Number().Int64())
	}
	from := query.FromBlock.Int64()
	to := query.ToBlock.Int64()

	state, err := b.bc.State()
	if err != nil {
		return nil, err
	}
	bc, ok := b.bc.(*blockchain.BlockChain)
	if !ok {
		return nil, errors.New("BlockChainForCaller is not blockchain.BlockChain")
	}
	filter := filters.NewRangeFilter(&filterBackend{state.Database().TrieDB().DiskDB(), bc}, from, to, query.Addresses, query.Topics)

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

func (b *BlockchainContractBackend) SubscribeFilterLogs(ctx context.Context, query klaytn.FilterQuery, ch chan<- types.Log) (klaytn.Subscription, error) {
	// Subscribe to contract events
	sink := make(chan []*types.Log)

	if b.events == nil {
		return nil, errors.New("events system not configured")
	}
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

// bind.DeployBackend defined methods

func (b *BlockchainContractBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	bc, ok := b.bc.(*blockchain.BlockChain)
	if !ok {
		return nil, errors.New("BlockChainForCaller is not blockchain.BlockChain")
	}
	receipt := bc.GetReceiptByTxHash(txHash)
	if receipt != nil {
		return receipt, nil
	}
	return nil, errors.New("receipt does not exist")
}

// sc.Backend requires BalanceAt and CurrentBlockNumber

func (b *BlockchainContractBackend) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	if _, state, err := b.getBlockAndState(blockNumber); err != nil {
		return nil, err
	} else {
		return state.GetBalance(account), nil
	}
}

func (b *BlockchainContractBackend) CurrentBlockNumber(ctx context.Context) (uint64, error) {
	return b.bc.CurrentBlock().NumberU64(), nil
}
