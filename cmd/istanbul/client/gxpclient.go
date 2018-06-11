package client

import (
	"math/big"
	"context"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/core/types"
	"ground-x/go-gxplatform"
)

// Blockchain Access

// BlockByHash returns the given full block.
//
// Note that loading full blocks requires two requests. Use HeaderByHash
// if you don't need all transactions or uncle headers.
func (c *client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return c.gxpClient.BlockByHash(ctx, hash)
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (c *client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return c.gxpClient.BlockByNumber(ctx, number)
}

// HeaderByHash returns the block header with the given hash.
func (c *client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return c.gxpClient.HeaderByHash(ctx, hash)
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (c *client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.gxpClient.HeaderByNumber(ctx, number)
}

// TransactionByHash returns the transaction with the given hash.
func (c *client) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return c.gxpClient.TransactionByHash(ctx, hash)
}

// TransactionCount returns the total number of transactions in the given block.
func (c *client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	return c.gxpClient.TransactionCount(ctx, blockHash)
}

// TransactionInBlock returns a single transaction at index in the given block.
func (c *client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	return c.gxpClient.TransactionInBlock(ctx, blockHash, index)
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (c *client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.gxpClient.TransactionReceipt(ctx, txHash)
}

// SyncProgress retrieves the current progress of the sync algorithm. If there's
// no sync currently running, it returns nil.
func (c *client) SyncProgress(ctx context.Context) (*gxplatform.SyncProgress, error) {
	return c.gxpClient.SyncProgress(ctx)
}

// SubscribeNewHead subscribes to notifications about the current blockchain head
// on the given channel.
func (c *client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (gxplatform.Subscription, error) {
	return c.gxpClient.SubscribeNewHead(ctx, ch)
}

// State Access

// NetworkID returns the network ID (also known as the chain ID) for this chain.
func (c *client) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.gxpClient.NetworkID(ctx)
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (c *client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.gxpClient.BalanceAt(ctx, account, blockNumber)
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (c *client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	return c.gxpClient.StorageAt(ctx, account, key, blockNumber)
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (c *client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.gxpClient.CodeAt(ctx, account, blockNumber)
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (c *client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return c.gxpClient.NonceAt(ctx, account, blockNumber)
}

// Filters

// FilterLogs executes a filter query.
func (c *client) FilterLogs(ctx context.Context, q gxplatform.FilterQuery) ([]types.Log, error) {
	return c.gxpClient.FilterLogs(ctx, q)
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (c *client) SubscribeFilterLogs(ctx context.Context, q gxplatform.FilterQuery, ch chan<- types.Log) (gxplatform.Subscription, error) {
	return c.gxpClient.SubscribeFilterLogs(ctx, q, ch)
}

// Pending State

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (c *client) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	return c.gxpClient.PendingBalanceAt(ctx, account)
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (c *client) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	return c.gxpClient.PendingStorageAt(ctx, account, key)
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (c *client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return c.gxpClient.PendingCodeAt(ctx, account)
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (c *client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return c.gxpClient.PendingNonceAt(ctx, account)
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (c *client) PendingTransactionCount(ctx context.Context) (uint, error) {
	return c.gxpClient.PendingTransactionCount(ctx)
}

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
//
// blockNumber selects the block height at which the call runs. It can be nil, in which
// case the code is taken from the latest known block. Note that state from very old
// blocks might not be available.
func (c *client) CallContract(ctx context.Context, msg gxplatform.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return c.gxpClient.CallContract(ctx, msg, blockNumber)
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (c *client) PendingCallContract(ctx context.Context, msg gxplatform.CallMsg) ([]byte, error) {
	return c.gxpClient.PendingCallContract(ctx, msg)
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (c *client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.gxpClient.SuggestGasPrice(ctx)
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (c *client) EstimateGas(ctx context.Context, msg gxplatform.CallMsg) (uint64, error) {
	return c.gxpClient.EstimateGas(ctx, msg)
}

// SendRawTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (c *client) SendRawTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.gxpClient.SendTransaction(ctx, tx)
}

