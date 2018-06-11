package client

import (
	"math/big"
	"context"
	"ground-x/go-gxplatform/p2p"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/core/types"
	"ground-x/go-gxplatform"
)

type Client interface {
	Close()
	AddPeer(ctx context.Context, nodeURL string) error
	AdminPeers(ctx context.Context) ([]*p2p.PeerInfo, error)
	NodeInfo(ctx context.Context) (*p2p.PeerInfo, error)
	BlockNumber(ctx context.Context) (*big.Int, error)
	StartMining(ctx context.Context) error
	StopMining(ctx context.Context) error
	SendTransaction(ctx context.Context, from, to common.Address, value *big.Int) (string, error)
	CreateContract(ctx context.Context, from common.Address, bytecode string, gas *big.Int) (string, error)
	CreatePrivateContract(ctx context.Context, from common.Address, bytecode string, gas *big.Int, privateFor []string) (string, error)
	ProposeValidator(ctx context.Context, address common.Address, auth bool) error
	GetValidators(ctx context.Context, blockNumbers *big.Int) ([]common.Address, error)

	// gxp client
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	SyncProgress(ctx context.Context) (*gxplatform.SyncProgress, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (gxplatform.Subscription, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	FilterLogs(ctx context.Context, q gxplatform.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q gxplatform.FilterQuery, ch chan<- types.Log) (gxplatform.Subscription, error)
	PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
	PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	PendingTransactionCount(ctx context.Context) (uint, error)
	CallContract(ctx context.Context, msg gxplatform.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCallContract(ctx context.Context, msg gxplatform.CallMsg) ([]byte, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg gxplatform.CallMsg) (uint64, error)
	SendRawTransaction(ctx context.Context, tx *types.Transaction) error
}

