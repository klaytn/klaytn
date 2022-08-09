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
// This file is derived from eth/tracers/api_test.go (2022/08/08).
// Modified and improved for the klaytn development.
package tracers

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"testing"

	klaytnapi "github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
)

var (
	errStateNotFound       = errors.New("state not found")
	errBlockNotFound       = errors.New("block not found")
	errTransactionNotFound = errors.New("transaction not found")
)

type testBackend struct {
	chainConfig *params.ChainConfig
	engine      consensus.Engine
	chaindb     database.DBManager
	chain       *blockchain.BlockChain
}

func newTestBackend(t *testing.T, n int, gspec *blockchain.Genesis, generator func(i int, b *blockchain.BlockGen)) *testBackend {
	backend := &testBackend{
		chainConfig: params.TestChainConfig,
		engine:      gxhash.NewFaker(),
		chaindb:     database.NewMemoryDBManager(),
	}
	// Generate blocks for testing
	gspec.Config = backend.chainConfig
	var (
		gendb   = database.NewMemoryDBManager()
		genesis = gspec.MustCommit(gendb)
	)
	blocks, _ := blockchain.GenerateChain(backend.chainConfig, genesis, backend.engine, gendb, n, generator)
	// Import the canonical chain
	gspec.MustCommit(backend.chaindb)
	cacheConfig := &blockchain.CacheConfig{
		CacheSize:           512,
		BlockInterval:       blockchain.DefaultBlockInterval,
		TriesInMemory:       blockchain.DefaultTriesInMemory,
		TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
		SnapshotCacheSize:   512,
		ArchiveMode:         true, // Archive mode
	}
	chain, err := blockchain.NewBlockChain(backend.chaindb, cacheConfig, backend.chainConfig, backend.engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	backend.chain = chain
	return backend
}

func (b *testBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.chain.GetHeaderByHash(hash), nil
}

func (b *testBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentHeader(), nil
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}

func (b *testBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.chain.GetBlockByHash(hash), nil
}

func (b *testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock(), nil
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}

func (b *testBackend) GetTxAndLookupInfo(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	tx, hash, blockNumber, index := b.chain.GetTxAndLookupInfoInCache(txHash)
	if tx == nil {
		return nil, common.Hash{}, 0, 0
	}
	return tx, hash, blockNumber, index
}

func (b *testBackend) RPCGasCap() *big.Int {
	return big.NewInt(250 * params.Ston)
}

func (b *testBackend) ChainConfig() *params.ChainConfig {
	return b.chainConfig
}

func (b *testBackend) Engine() consensus.Engine {
	return b.engine
}

func (b *testBackend) ChainDB() database.DBManager {
	return b.chaindb
}

func (b *testBackend) StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (*state.StateDB, error) {
	statedb, err := b.chain.StateAt(block.Root())
	if err != nil {
		return nil, errStateNotFound
	}
	return statedb, nil
}

func (b *testBackend) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (blockchain.Message, vm.Context, *state.StateDB, error) {
	parent := b.chain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.Context{}, nil, errBlockNotFound
	}
	statedb, err := b.chain.StateAt(parent.Root())
	if err != nil {
		return nil, vm.Context{}, nil, errStateNotFound
	}
	if txIndex == 0 && len(block.Transactions()) == 0 {
		return nil, vm.Context{}, statedb, nil
	}
	// Recompute transactions up to the target index.
	signer := types.MakeSigner(b.chainConfig, block.Number())
	for idx, tx := range block.Transactions() {
		msg, _ := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		context := blockchain.NewEVMContext(msg, block.Header(), b.chain, nil)
		if idx == txIndex {
			return msg, context, statedb, nil
		}
		vmenv := vm.NewEVM(context, statedb, b.chainConfig, &vm.Config{Debug: true, EnableInternalTxTracing: true})
		if _, _, kerr := blockchain.ApplyMessage(vmenv, msg); kerr.ErrTxInvalid != nil {
			return nil, vm.Context{}, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), kerr.ErrTxInvalid)
		}
		statedb.Finalise(true, true)
	}
	return nil, vm.Context{}, nil, fmt.Errorf("transaction index %d out of range for block %#x", txIndex, block.Hash())
}

// TODO-tracer: implement TraceCall
//func TestTraceCall(t *testing.T) {
//	t.Parallel()
//
//	// Initialize test accounts
//	accounts := newAccounts(3)
//	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
//		accounts[0].addr: {Balance: big.NewInt(params.KLAY)},
//		accounts[1].addr: {Balance: big.NewInt(params.KLAY)},
//		accounts[2].addr: {Balance: big.NewInt(params.KLAY)},
//	}}
//	genBlocks := 10
//	signer := types.HomesteadSigner{}
//	api := NewAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *blockchain.BlockGen) {
//		// Transfer from account[0] to account[1]
//		//    value: 1000 wei
//		//    fee:   0 wei
//		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
//		b.AddTx(tx)
//	}))
//
//	var testSuite = []struct {
//		blockNumber rpc.BlockNumber
//		call        klaytnapi.CallArgs
//		config      *TraceConfig
//		expectErr   error
//		expect      interface{}
//	}{
//		// Standard JSON trace upon the genesis, plain transfer.
//		{
//			blockNumber: rpc.BlockNumber(0),
//			call: klaytnapi.CallArgs{
//				From:  accounts[0].addr,
//				To:    &accounts[1].addr,
//				Value: (hexutil.Big)(*big.NewInt(1000)),
//			},
//			config:    nil,
//			expectErr: nil,
//			expect: &klaytnapi.ExecutionResult{
//				Gas:         params.TxGas,
//				Failed:      false,
//				ReturnValue: "",
//				StructLogs:  []klaytnapi.StructLogRes{},
//			},
//		},
//		// Standard JSON trace upon the head, plain transfer.
//		{
//			blockNumber: rpc.BlockNumber(genBlocks),
//			call: klaytnapi.CallArgs{
//				From:  accounts[0].addr,
//				To:    &accounts[1].addr,
//				Value: (hexutil.Big)(*big.NewInt(1000)),
//			},
//			config:    nil,
//			expectErr: nil,
//			expect: &klaytnapi.ExecutionResult{
//				Gas:         params.TxGas,
//				Failed:      false,
//				ReturnValue: "",
//				StructLogs:  []klaytnapi.StructLogRes{},
//			},
//		},
//		// Standard JSON trace upon the non-existent block, error expects
//		{
//			blockNumber: rpc.BlockNumber(genBlocks + 1),
//			call: klaytnapi.CallArgs{
//				From:  accounts[0].addr,
//				To:    &accounts[1].addr,
//				Value: (hexutil.Big)(*big.NewInt(1000)),
//			},
//			config:    nil,
//			expectErr: fmt.Errorf("block #%d not found", genBlocks+1),
//			expect:    nil,
//		},
//		// Standard JSON trace upon the latest block
//		{
//			blockNumber: rpc.LatestBlockNumber,
//			call: klaytnapi.CallArgs{
//				From:  accounts[0].addr,
//				To:    &accounts[1].addr,
//				Value: (hexutil.Big)(*big.NewInt(1000)),
//			},
//			config:    nil,
//			expectErr: nil,
//			expect: &klaytnapi.ExecutionResult{
//				Gas:         params.TxGas,
//				Failed:      false,
//				ReturnValue: "",
//				StructLogs:  []klaytnapi.StructLogRes{},
//			},
//		},
//		// Standard JSON trace upon the pending block
//		{
//			blockNumber: rpc.PendingBlockNumber,
//			call: klaytnapi.CallArgs{
//				From:  accounts[0].addr,
//				To:    &accounts[1].addr,
//				Value: (hexutil.Big)(*big.NewInt(1000)),
//			},
//			config:    nil,
//			expectErr: nil,
//			expect: &klaytnapi.ExecutionResult{
//				Gas:         params.TxGas,
//				Failed:      false,
//				ReturnValue: "",
//				StructLogs:  []klaytnapi.StructLogRes{},
//			},
//		},
//	}
//	for _, testspec := range testSuite {
//		result, err := api.TraceCall(context.Background(), testspec.call, rpc.BlockNumberOrHash{BlockNumber: &testspec.blockNumber}, testspec.config)
//		if testspec.expectErr != nil {
//			if err == nil {
//				t.Errorf("Expect error %v, get nothing", testspec.expectErr)
//				continue
//			}
//			if !reflect.DeepEqual(err, testspec.expectErr) {
//				t.Errorf("Error mismatch, want %v, get %v", testspec.expectErr, err)
//			}
//		} else {
//			if err != nil {
//				t.Errorf("Expect no error, get %v", err)
//				continue
//			}
//			if !reflect.DeepEqual(result, testspec.expect) {
//				t.Errorf("Result mismatch, want %v, get %v", testspec.expect, result)
//			}
//		}
//	}
//}

func TestTraceTransaction(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(2)
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.KLAY)},
		accounts[1].addr: {Balance: big.NewInt(params.KLAY)},
	}}
	target := common.Hash{}
	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
	api := NewAPI(newTestBackend(t, 1, genesis, func(i int, b *blockchain.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(1), nil), signer, accounts[0].key)
		b.AddTx(tx)
		target = tx.Hash()
	}))
	result, err := api.TraceTransaction(context.Background(), target, nil)
	if err != nil {
		t.Errorf("Failed to trace transaction %v", err)
	}
	if !reflect.DeepEqual(result, &klaytnapi.ExecutionResult{
		Gas:         params.TxGas,
		Failed:      false,
		ReturnValue: "",
		StructLogs:  []klaytnapi.StructLogRes{},
	}) {
		t.Error("Transaction tracing result is different")
	}
}

func TestTraceBlock(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.KLAY)},
		accounts[1].addr: {Balance: big.NewInt(params.KLAY)},
		accounts[2].addr: {Balance: big.NewInt(params.KLAY)},
	}}
	genBlocks := 10
	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
	api := NewAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *blockchain.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
		b.AddTx(tx)
	}))

	testSuite := []struct {
		blockNumber rpc.BlockNumber
		config      *TraceConfig
		expect      interface{}
		expectErr   error
	}{
		// Trace genesis block, expect error
		{
			blockNumber: rpc.BlockNumber(0),
			config:      nil,
			expect:      nil,
			expectErr:   errors.New("genesis is not traceable"),
		},
		// Trace head block
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{
					Result: &klaytnapi.ExecutionResult{
						Gas:         params.TxGas,
						Failed:      false,
						ReturnValue: "",
						StructLogs:  []klaytnapi.StructLogRes{},
					},
				},
			},
		},
		// Trace non-existent block
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			config:      nil,
			expectErr:   fmt.Errorf("block #%d not found", genBlocks+1),
			expect:      nil,
		},
		// Trace latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{
					Result: &klaytnapi.ExecutionResult{
						Gas:         params.TxGas,
						Failed:      false,
						ReturnValue: "",
						StructLogs:  []klaytnapi.StructLogRes{},
					},
				},
			},
		},
		// Trace pending block
		{
			blockNumber: rpc.PendingBlockNumber,
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{
					Result: &klaytnapi.ExecutionResult{
						Gas:         params.TxGas,
						Failed:      false,
						ReturnValue: "",
						StructLogs:  []klaytnapi.StructLogRes{},
					},
				},
			},
		},
	}
	for _, testspec := range testSuite {
		result, err := api.TraceBlockByNumber(context.Background(), testspec.blockNumber, testspec.config)
		if testspec.expectErr != nil {
			if err == nil {
				t.Errorf("Expect error %v, get nothing", testspec.expectErr)
				continue
			}
			if !reflect.DeepEqual(err, testspec.expectErr) {
				t.Errorf("Error mismatch, want %v, get %v", testspec.expectErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("Expect no error, get %v", err)
				continue
			}
			if len(result) != len(testspec.expect.([]*txTraceResult)) {
				t.Errorf("Result length mismatch, want %v, get %v", len(result), len(testspec.expect.([]*txTraceResult)))
			}
			for idx, r := range result {
				if !reflect.DeepEqual(r.Result, testspec.expect.([]*txTraceResult)[idx].Result) {
					t.Errorf("Result mismatch, want %v, get %v", testspec.expect, result)
				}
			}
		}
	}
}

type Account struct {
	key  *ecdsa.PrivateKey
	addr common.Address
}

type Accounts []Account

func (a Accounts) Len() int           { return len(a) }
func (a Accounts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Accounts) Less(i, j int) bool { return bytes.Compare(a[i].addr.Bytes(), a[j].addr.Bytes()) < 0 }

func newAccounts(n int) (accounts Accounts) {
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		accounts = append(accounts, Account{key: key, addr: addr})
	}
	sort.Sort(accounts)
	return accounts
}
