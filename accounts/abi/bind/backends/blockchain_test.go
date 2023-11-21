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
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/node/cn/filters"
	mock_filter "github.com/klaytn/klaytn/node/cn/filters/mock"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

var (
	testAddr  = crypto.PubkeyToAddress(testKey.PublicKey)
	code1Addr = common.HexToAddress("0x1111111111111111111111111111111111111111")
	code2Addr = common.HexToAddress("0x2222222222222222222222222222222222222222")

	parsedAbi1, _ = abi.JSON(strings.NewReader(abiJSON))
	parsedAbi2, _ = abi.JSON(strings.NewReader(reverterABI))
	code1Bytes    = common.FromHex(deployedCode)
	code2Bytes    = common.FromHex(reverterDeployedBin)
)

func newTestBlockchain() *blockchain.BlockChain {
	config := params.TestChainConfig.Copy()
	return newTestBlockchainWithConfig(config)
}

func newTestBlockchainWithConfig(config *params.ChainConfig) *blockchain.BlockChain {
	alloc := blockchain.GenesisAlloc{
		testAddr:  {Balance: big.NewInt(10000000000)},
		code1Addr: {Balance: big.NewInt(0), Code: code1Bytes},
		code2Addr: {Balance: big.NewInt(0), Code: code2Bytes},
	}

	db := database.NewMemoryDBManager()
	genesis := blockchain.Genesis{Config: config, Alloc: alloc}
	genesis.MustCommit(db)

	bc, _ := blockchain.NewBlockChain(db, nil, genesis.Config, gxhash.NewFaker(), vm.Config{})

	// Append 10 blocks to test with block numbers other than 0
	block := bc.CurrentBlock()
	blocks, _ := blockchain.GenerateChain(config, block, gxhash.NewFaker(), db, 10, func(i int, b *blockchain.BlockGen) {})
	bc.InsertChain(blocks)

	return bc
}

func TestBlockchainCodeAt(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractBackend(bc, nil, nil)

	// Normal cases
	code, err := c.CodeAt(context.Background(), code1Addr, nil)
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	code, err = c.CodeAt(context.Background(), code2Addr, nil)
	assert.Nil(t, err)
	assert.Equal(t, code2Bytes, code)

	code, err = c.CodeAt(context.Background(), code1Addr, common.Big0)
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	code, err = c.CodeAt(context.Background(), code1Addr, common.Big1)
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	code, err = c.CodeAt(context.Background(), code1Addr, big.NewInt(10))
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	// Non-code address
	code, err = c.CodeAt(context.Background(), testAddr, nil)
	assert.True(t, code == nil && err == nil)

	// Invalid block number
	code, err = c.CodeAt(context.Background(), code1Addr, big.NewInt(11))
	assert.True(t, code == nil && err == errBlockDoesNotExist)
}

func TestBlockchainCallContract(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractBackend(bc, nil, nil)

	data_receive, _ := parsedAbi1.Pack("receive", []byte("X"))
	data_revertString, _ := parsedAbi2.Pack("revertString")
	data_revertNoString, _ := parsedAbi2.Pack("revertNoString")

	// Normal case
	ret, err := c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code1Addr,
		Gas:  1000000,
		Data: data_receive,
	}, nil)
	assert.Nil(t, err)
	assert.Equal(t, expectedReturn, ret)

	// Error outside VM - Intrinsic Gas
	ret, err = c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code1Addr,
		Gas:  20000,
		Data: data_receive,
	}, nil)
	assert.True(t, errors.Is(err, blockchain.ErrIntrinsicGas))

	// VM revert error - empty reason
	ret, err = c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code2Addr,
		Gas:  100000,
		Data: data_revertNoString,
	}, nil)
	assert.Equal(t, "execution reverted: ", err.Error())

	// VM revert error - string reason
	ret, err = c.CallContract(context.Background(), klaytn.CallMsg{
		From: testAddr,
		To:   &code2Addr,
		Gas:  100000,
		Data: data_revertString,
	}, nil)
	assert.Equal(t, "execution reverted: some error", err.Error())
}

func TestBlockchainPendingCodeAt(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractBackend(bc, nil, nil)

	// Normal cases
	code, err := c.PendingCodeAt(context.Background(), code1Addr)
	assert.Nil(t, err)
	assert.Equal(t, code1Bytes, code)

	code, err = c.PendingCodeAt(context.Background(), code2Addr)
	assert.Nil(t, err)
	assert.Equal(t, code2Bytes, code)

	// Non-code address
	code, err = c.PendingCodeAt(context.Background(), testAddr)
	assert.True(t, code == nil && err == nil)
}

func TestBlockChainSuggestGasPrice(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractBackend(bc, nil, nil)

	// Normal case
	gasPrice, err := c.SuggestGasPrice(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, params.TestChainConfig.UnitPrice, gasPrice.Uint64())

	config := params.TestChainConfig.Copy()
	config.IstanbulCompatibleBlock = common.Big0
	config.LondonCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0
	config.Governance = params.GetDefaultGovernanceConfig()
	config.Governance.KIP71.LowerBoundBaseFee = 0
	bc = newTestBlockchainWithConfig(config)
	c = NewBlockchainContractBackend(bc, nil, nil)

	// Normal case
	gasPrice, err = c.SuggestGasPrice(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, bc.CurrentBlock().Header().BaseFee.Uint64()*2, gasPrice.Uint64())
}

func TestBlockChainEstimateGas(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractBackend(bc, nil, nil)

	// Normal case
	gas, err := c.EstimateGas(context.Background(), klaytn.CallMsg{
		From:  testAddr,
		To:    &testAddr,
		Value: big.NewInt(1000),
	})
	assert.Nil(t, err)
	assert.Equal(t, uint64(params.TxGas), gas)

	// Error case - simple transfer with insufficient funds with zero gasPrice
	gas, err = c.EstimateGas(context.Background(), klaytn.CallMsg{
		From:  code1Addr,
		To:    &code1Addr,
		Value: big.NewInt(1),
	})
	assert.Contains(t, err.Error(), "insufficient balance for transfer")
	assert.Zero(t, gas)
}

func TestBlockChainSendTransaction(t *testing.T) {
	bc := newTestBlockchain()
	block := bc.CurrentBlock()
	state, err := bc.State()
	txPoolConfig := blockchain.DefaultTxPoolConfig
	txPoolConfig.Journal = "/dev/null" // disable journaling to file
	txPool := blockchain.NewTxPool(txPoolConfig, bc.Config(), bc)
	defer txPool.Stop()
	assert.Nil(t, err)
	c := NewBlockchainContractBackend(bc, txPool, nil)

	// create a signed transaction to send
	nonce := state.GetNonce(testAddr)
	tx := types.NewTransaction(nonce, testAddr, big.NewInt(1000), params.TxGas, big.NewInt(1), nil)
	chainId, err := c.ChainID(context.Background())
	if err != nil {
		t.Errorf("could not get chain ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainId), testKey)
	if err != nil {
		t.Errorf("could not sign tx: %v", err)
	}

	// send tx to simulated backend
	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Errorf("could not add tx to pending block: %v", err)
	}

	blocks, _ := blockchain.GenerateChain(bc.Config(), block, gxhash.NewFaker(), state.Database().TrieDB().DiskDB(), 1, func(i int, b *blockchain.BlockGen) {
		txs, err := txPool.Pending()
		if err != nil {
			t.Errorf("could not get pending txs: %v", err)
		}
		for _, v := range txs {
			for _, v2 := range v {
				b.AddTx(v2)
			}
		}
	})
	bc.InsertChain(blocks)

	block = bc.GetBlockByNumber(11)
	if block == nil {
		t.Errorf("could not get block at height 11")
	}

	assert.True(t, len(block.Transactions()) != 0)
	if signedTx.Hash() != block.Transactions()[0].Hash() {
		t.Errorf("did not commit sent transaction. expected hash %v got hash %v", block.Transactions()[0].Hash(), signedTx.Hash())
	}
	assert.False(t, block.Header().EmptyReceipts())
}

func TestBlockChainChainID(t *testing.T) {
	bc := newTestBlockchain()
	c := NewBlockchainContractBackend(bc, nil, nil)

	// Normal case
	chainId, err := c.ChainID(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, params.TestChainConfig.ChainID, chainId)
}

func initBackendForFiltererTests(t *testing.T, bc *blockchain.BlockChain) *BlockchainContractBackend {
	block := bc.CurrentBlock()
	state, _ := bc.State()

	// Add one block with contract execution to generate logs
	data_receive, _ := parsedAbi1.Pack("receive", []byte("X"))
	tx := types.NewTransaction(uint64(0), code1Addr, big.NewInt(0), 50000, big.NewInt(1), data_receive)
	chainId := bc.Config().ChainID

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainId), testKey)
	if err != nil {
		t.Errorf("could not sign tx: %v", err)
	}

	blocks, _ := blockchain.GenerateChain(bc.Config(), block, gxhash.NewFaker(), state.Database().TrieDB().DiskDB(), 1, func(i int, b *blockchain.BlockGen) {
		b.AddTx(signedTx)
	})
	bc.InsertChain(blocks)

	// mock filterer backend
	mockCtrl := gomock.NewController(t)
	mockBackend := mock_filter.NewMockBackend(mockCtrl)

	any := gomock.Any()
	txPoolConfig := blockchain.DefaultTxPoolConfig
	txPoolConfig.Journal = "/dev/null" // disable journaling to file
	txPool := blockchain.NewTxPool(txPoolConfig, bc.Config(), bc)
	subscribeNewTxsEvent := func(ch chan<- blockchain.NewTxsEvent) klaytn.Subscription {
		return txPool.SubscribeNewTxsEvent(ch)
	}
	subscribeLogsEvent := func(ch chan<- []*types.Log) klaytn.Subscription {
		return bc.SubscribeLogsEvent(ch)
	}
	subscribeRemovedLogsEvent := func(ch chan<- blockchain.RemovedLogsEvent) klaytn.Subscription {
		return bc.SubscribeRemovedLogsEvent(ch)
	}
	subscribeChainEvent := func(ch chan<- blockchain.ChainEvent) klaytn.Subscription {
		return bc.SubscribeChainEvent(ch)
	}
	mockBackend.EXPECT().SubscribeNewTxsEvent(any).DoAndReturn(subscribeNewTxsEvent).AnyTimes()
	mockBackend.EXPECT().SubscribeLogsEvent(any).DoAndReturn(subscribeLogsEvent).AnyTimes()
	mockBackend.EXPECT().SubscribeRemovedLogsEvent(any).DoAndReturn(subscribeRemovedLogsEvent).AnyTimes()
	mockBackend.EXPECT().SubscribeChainEvent(any).DoAndReturn(subscribeChainEvent).AnyTimes()

	f := filters.NewEventSystem(&event.TypeMux{}, mockBackend, false)
	c := NewBlockchainContractBackend(bc, nil, f)

	return c
}

func TestBlockChainFilterLogs(t *testing.T) {
	bc := newTestBlockchain()
	c := initBackendForFiltererTests(t, bc)

	// Normal case
	logs, err := c.FilterLogs(context.Background(), klaytn.FilterQuery{
		FromBlock: big.NewInt(10),
		ToBlock:   big.NewInt(11),
		Addresses: []common.Address{code1Addr},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(logs))

	// No logs exist for code2Addr
	logs, err = c.FilterLogs(context.Background(), klaytn.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(11),
		Addresses: []common.Address{code2Addr},
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(logs))
}

func TestBlockChainSubscribeFilterLogs(t *testing.T) {
	bc := newTestBlockchain()
	c := initBackendForFiltererTests(t, bc)

	logs := make(chan types.Log)
	sub, err := c.SubscribeFilterLogs(context.Background(), klaytn.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(20),
		Addresses: []common.Address{code1Addr},
	}, logs)
	assert.Nil(t, err)
	assert.NotNil(t, sub)

	// Insert a block with contract execution to generate logs
	go func() {
		state, _ := bc.State()
		nonce := state.GetNonce(testAddr)
		data_receive, _ := parsedAbi1.Pack("receive", []byte("X"))
		tx := types.NewTransaction(nonce, code1Addr, big.NewInt(0), 50000, big.NewInt(1), data_receive)
		chainId, _ := c.ChainID(context.Background())

		signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainId), testKey)
		if err != nil {
			t.Errorf("could not sign tx: %v", err)
		}

		block := bc.CurrentBlock()

		blocks, _ := blockchain.GenerateChain(c.bc.Config(), block, gxhash.NewFaker(), state.Database().TrieDB().DiskDB(), 1, func(i int, b *blockchain.BlockGen) {
			b.AddTx(signedTx)
		})
		bc.InsertChain(blocks)
	}()

	// Wait for 2 logs
	for i := 0; i < 2; i++ {
		select {
		case log := <-logs:
			assert.Equal(t, code1Addr, log.Address)
			assert.Contains(t, hex.EncodeToString(log.Data), hex.EncodeToString(testAddr.Bytes()))
		case <-time.After(3 * time.Second):
			t.Fatal("timeout while waiting for logs")
		}
	}
}
