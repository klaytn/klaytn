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
// This file is derived from internal/ethapi/backend.go (2018/06/04).
// Modified and improved for the klaytn development.

package api

import (
	"context"
	"math/big"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

//go:generate mockgen -destination=api/mocks/backend_mock.go github.com/klaytn/klaytn/api Backend
// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	// General Klaytn API
	Progress() klaytn.SyncProgress
	ProtocolVersion() int
	SuggestPrice(ctx context.Context) (*big.Int, error)
	ChainDB() database.DBManager
	EventMux() *event.TypeMux
	AccountManager() accounts.AccountManager
	RPCGasCap() *big.Int  // global gas cap for klay_call over rpc: DoS protection
	RPCTxFeeCap() float64 // global tx fee cap for all transaction related APIs
	Engine() consensus.Engine
	FeeHistory(ctx context.Context, blockCount int, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error)

	// BlockChain API
	SetHead(number uint64)
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error)
	BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error)
	BlockByHash(ctx context.Context, blockHash common.Hash) (*types.Block, error)
	BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error)
	StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error)
	StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error)
	GetBlockReceipts(ctx context.Context, blockHash common.Hash) types.Receipts
	GetTxLookupInfoAndReceipt(ctx context.Context, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt)
	GetTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	GetTd(blockHash common.Hash) *big.Int
	GetEVM(ctx context.Context, msg blockchain.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error)
	SubscribeChainEvent(ch chan<- blockchain.ChainEvent) event.Subscription
	SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription
	SubscribeChainSideEvent(ch chan<- blockchain.ChainSideEvent) event.Subscription
	IsParallelDBWrite() bool

	IsSenderTxHashIndexingEnabled() bool

	// TxPool API
	SendTx(ctx context.Context, signedTx *types.Transaction) error
	GetPoolTransactions() (types.Transactions, error)
	GetPoolTransaction(txHash common.Hash) *types.Transaction
	GetPoolNonce(ctx context.Context, addr common.Address) uint64
	Stats() (pending int, queued int)
	TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions)
	SubscribeNewTxsEvent(chan<- blockchain.NewTxsEvent) event.Subscription

	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block

	GetTxAndLookupInfoInCache(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64)
	GetBlockReceiptsInCache(blockHash common.Hash) types.Receipts
	GetTxLookupInfoAndReceiptInCache(Hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt)
}

func GetAPIs(apiBackend Backend) ([]rpc.API, *EthereumAPI) {
	nonceLock := new(AddrLocker)

	ethAPI := NewEthereumAPI()

	publicKlayAPI := NewPublicKlayAPI(apiBackend)
	publicBlockChainAPI := NewPublicBlockChainAPI(apiBackend)
	publicTransactionPoolAPI := NewPublicTransactionPoolAPI(apiBackend, nonceLock)
	publicAccountAPI := NewPublicAccountAPI(apiBackend.AccountManager())

	ethAPI.SetPublicKlayAPI(publicKlayAPI)
	ethAPI.SetPublicBlockChainAPI(publicBlockChainAPI)
	ethAPI.SetPublicTransactionPoolAPI(publicTransactionPoolAPI)
	ethAPI.SetPublicAccountAPI(publicAccountAPI)

	return []rpc.API{
		{
			Namespace: "klay",
			Version:   "1.0",
			Service:   publicKlayAPI,
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   publicBlockChainAPI,
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   publicTransactionPoolAPI,
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(apiBackend),
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   publicAccountAPI,
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountAPI(apiBackend, nonceLock),
			Public:    false,
		},
	}, ethAPI
}
