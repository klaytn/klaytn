// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethclient/ethclient.go (2018/06/04).
// Modified and improved for the klaytn development.

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/rlp"
)

// TODO-Klaytn Needs to separate APIs along with each namespaces.

// Client defines typed wrappers for the Klaytn RPC API.
type Client struct {
	c       *rpc.Client
	chainID *big.Int
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client) *Client {
	return &Client{c, nil}
}

func (ec *Client) Close() {
	ec.c.Close()
}

// Blockchain Access

// BlockByHash returns the given full block.
//
// Note that loading full blocks requires two requests. Use HeaderByHash
// if you don't need all transactions.
func (ec *Client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return ec.getBlock(ctx, "klay_getBlockByHash", hash, true)
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions.
func (ec *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return ec.getBlock(ctx, "klay_getBlockByNumber", toBlockNumArg(number), true)
}

type rpcBlock struct {
	Hash         common.Hash      `json:"hash"`
	Transactions []rpcTransaction `json:"transactions"`
}

func (ec *Client) getBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
	var raw json.RawMessage
	err := ec.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, klaytn.NotFound
	}
	// Decode header and transactions.
	var head *types.Header
	var body rpcBlock
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// TODO-Klaytn Enable the below error checks after having a way to get the correct EmptyRootHash
	// Quick-verify transaction lists. This mostly helps with debugging the server.
	//if head.TxHash == types.EmptyRootHash && len(body.Transactions) > 0 {
	//	return nil, fmt.Errorf("server returned non-empty transaction list but block header indicates no transactions")
	//}
	//if head.TxHash != types.EmptyRootHash && len(body.Transactions) == 0 {
	//	return nil, fmt.Errorf("server returned empty transaction list but block header indicates transactions")
	//}
	// Fill the sender cache of transactions in the block.
	txs := make([]*types.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		if tx.From != nil {
			setSenderFromServer(tx.tx, *tx.From, body.Hash)
		}
		txs[i] = tx.tx
	}
	return types.NewBlockWithHeader(head).WithBody(txs), nil
}

// HeaderByHash returns the block header with the given hash.
func (ec *Client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	var head *types.Header
	err := ec.c.CallContext(ctx, &head, "klay_getBlockByHash", hash, false)
	if err == nil && head == nil {
		err = klaytn.NotFound
	}
	return head, err
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var head *types.Header
	err := ec.c.CallContext(ctx, &head, "klay_getBlockByNumber", toBlockNumArg(number), false)
	if err == nil && head == nil {
		err = klaytn.NotFound
	}
	return head, err
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

// TransactionByHash returns the transaction with the given hash.
func (ec *Client) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	var json *rpcTransaction
	err = ec.c.CallContext(ctx, &json, "klay_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	} else if json == nil {
		return nil, false, klaytn.NotFound
	} else if sigs := json.tx.RawSignatureValues(); sigs[0].V == nil {
		return nil, false, fmt.Errorf("server returned transaction without signature")
	}
	if json.From != nil && json.BlockHash != nil {
		setSenderFromServer(json.tx, *json.From, *json.BlockHash)
	}
	return json.tx, json.BlockNumber == nil, nil
}

// TransactionSender returns the sender address of the given transaction. The transaction
// must be known to the remote node and included in the blockchain at the given block and
// index. The sender is the one derived by the protocol at the time of inclusion.
//
// There is a fast-path for transactions retrieved by TransactionByHash and
// TransactionInBlock. Getting their sender address can be done without an RPC interaction.
func (ec *Client) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	// Try to load the address from the cache.
	sender, err := types.Sender(&senderFromServer{blockhash: block}, tx)
	if err == nil {
		return sender, nil
	}
	var meta struct {
		Hash common.Hash
		From common.Address
	}
	if err = ec.c.CallContext(ctx, &meta, "klay_getTransactionByBlockHashAndIndex", block, hexutil.Uint64(index)); err != nil {
		return common.Address{}, err
	}
	if meta.Hash == (common.Hash{}) || meta.Hash != tx.Hash() {
		return common.Address{}, errors.New("wrong inclusion block/index")
	}
	return meta.From, nil
}

// TransactionCount returns the total number of transactions in the given block.
func (ec *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	var num hexutil.Uint
	err := ec.c.CallContext(ctx, &num, "klay_getBlockTransactionCountByHash", blockHash)
	return uint(num), err
}

// TransactionInBlock returns a single transaction at index in the given block.
func (ec *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	var json *rpcTransaction
	err := ec.c.CallContext(ctx, &json, "klay_getTransactionByBlockHashAndIndex", blockHash, hexutil.Uint64(index))
	if err == nil {
		if json == nil {
			return nil, klaytn.NotFound
		} else if sigs := json.tx.RawSignatureValues(); sigs[0].V == nil {
			return nil, fmt.Errorf("server returned transaction without signature")
		}
	}
	if json.From != nil && json.BlockHash != nil {
		setSenderFromServer(json.tx, *json.From, *json.BlockHash)
	}
	return json.tx, err
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (ec *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var r *types.Receipt
	err := ec.c.CallContext(ctx, &r, "klay_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, klaytn.NotFound
		}
	}
	return r, err
}

// TransactionReceiptRpcOutput returns the receipt of a transaction by transaction hash as a rpc output.
func (ec *Client) TransactionReceiptRpcOutput(ctx context.Context, txHash common.Hash) (r map[string]interface{}, err error) {
	err = ec.c.CallContext(ctx, &r, "klay_getTransactionReceipt", txHash)
	if err == nil && r == nil {
		return nil, klaytn.NotFound
	}
	return
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

type rpcProgress struct {
	StartingBlock hexutil.Uint64
	CurrentBlock  hexutil.Uint64
	HighestBlock  hexutil.Uint64
	PulledStates  hexutil.Uint64
	KnownStates   hexutil.Uint64
}

// SyncProgress retrieves the current progress of the sync algorithm. If there's
// no sync currently running, it returns nil.
func (ec *Client) SyncProgress(ctx context.Context) (*klaytn.SyncProgress, error) {
	var raw json.RawMessage
	if err := ec.c.CallContext(ctx, &raw, "klay_syncing"); err != nil {
		return nil, err
	}
	// Handle the possible response types
	var syncing bool
	if err := json.Unmarshal(raw, &syncing); err == nil {
		return nil, nil // Not syncing (always false)
	}
	var progress *rpcProgress
	if err := json.Unmarshal(raw, &progress); err != nil {
		return nil, err
	}
	return &klaytn.SyncProgress{
		StartingBlock: uint64(progress.StartingBlock),
		CurrentBlock:  uint64(progress.CurrentBlock),
		HighestBlock:  uint64(progress.HighestBlock),
		PulledStates:  uint64(progress.PulledStates),
		KnownStates:   uint64(progress.KnownStates),
	}, nil
}

// SubscribeNewHead subscribes to notifications about the current blockchain head
// on the given channel.
func (ec *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (klaytn.Subscription, error) {
	return ec.c.KlaySubscribe(ctx, ch, "newHeads")
}

// State Access

// NetworkID returns the network ID (also known as the chain ID) for this chain.
func (ec *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	version := new(big.Int)
	var ver string
	if err := ec.c.CallContext(ctx, &ver, "net_version"); err != nil {
		return nil, err
	}
	if _, ok := version.SetString(ver, 10); !ok {
		return nil, fmt.Errorf("invalid net_version result %q", ver)
	}
	return version, nil
}

// BalanceAt returns the peb balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (ec *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "klay_getBalance", account, toBlockNumArg(blockNumber))
	return (*big.Int)(&result), err
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *Client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "klay_getStorageAt", account, key, toBlockNumArg(blockNumber))
	return result, err
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *Client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "klay_getCode", account, toBlockNumArg(blockNumber))
	return result, err
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	err := ec.c.CallContext(ctx, &result, "klay_getTransactionCount", account, toBlockNumArg(blockNumber))
	return uint64(result), err
}

// Filters

// FilterLogs executes a filter query.
func (ec *Client) FilterLogs(ctx context.Context, q klaytn.FilterQuery) ([]types.Log, error) {
	var result []types.Log
	err := ec.c.CallContext(ctx, &result, "klay_getLogs", toFilterArg(q))
	return result, err
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (ec *Client) SubscribeFilterLogs(ctx context.Context, q klaytn.FilterQuery, ch chan<- types.Log) (klaytn.Subscription, error) {
	return ec.c.KlaySubscribe(ctx, ch, "logs", toFilterArg(q))
}

func toFilterArg(q klaytn.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

// Pending State

// PendingBalanceAt returns the peb balance of the given account in the pending state.
func (ec *Client) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "klay_getBalance", account, "pending")
	return (*big.Int)(&result), err
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *Client) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "klay_getStorageAt", account, key, "pending")
	return result, err
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *Client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "klay_getCode", account, "pending")
	return result, err
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := ec.c.CallContext(ctx, &result, "klay_getTransactionCount", account, "pending")
	return uint64(result), err
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (ec *Client) PendingTransactionCount(ctx context.Context) (uint, error) {
	var num hexutil.Uint
	err := ec.c.CallContext(ctx, &num, "klay_getBlockTransactionCountByNumber", "pending")
	return uint(num), err
}

// TODO: SubscribePendingTransactions (needs server side)

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
//
// blockNumber selects the block height at which the call runs. It can be nil, in which
// case the code is taken from the latest known block. Note that state from very old
// blocks might not be available.
func (ec *Client) CallContract(ctx context.Context, msg klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.c.CallContext(ctx, &hex, "klay_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *Client) PendingCallContract(ctx context.Context, msg klaytn.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.c.CallContext(ctx, &hex, "klay_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var hex hexutil.Big
	if err := ec.c.CallContext(ctx, &hex, "klay_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the latest state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *Client) EstimateGas(ctx context.Context, msg klaytn.CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := ec.c.CallContext(ctx, &hex, "klay_estimateGas", toCallArg(msg))
	if err != nil {
		return 0, err
	}
	return uint64(hex), nil
}

// SendTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ec *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	_, err := ec.SendRawTransaction(ctx, tx)
	return err
	//data, err := rlp.EncodeToBytes(tx)
	//if err != nil {
	//	return err
	//}
	//return ec.c.CallContext(ctx, nil, "klay_sendRawTransaction", common.ToHex(data))
}

// SendRawTransaction injects a signed transaction into the pending pool for execution.
//
// This function can return the transaction hash and error.
func (ec *Client) SendRawTransaction(ctx context.Context, tx *types.Transaction) (common.Hash, error) {
	var hex hexutil.Bytes
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return common.Hash{}, err
	}
	if err := ec.c.CallContext(ctx, &hex, "klay_sendRawTransaction", hexutil.Encode(data)); err != nil {
		return common.Hash{}, err
	}
	hash := common.BytesToHash(hex)
	return hash, nil
}

// SendUnsignedTransaction injects a unsigned transaction into the pending pool for execution.
//
// This function can return the transaction hash and error.
func (ec *Client) SendUnsignedTransaction(ctx context.Context, from common.Address, to common.Address, gas uint64, gasPrice uint64, value *big.Int, data []byte, input []byte) (common.Hash, error) {
	var hex hexutil.Bytes

	tGas := hexutil.Uint64(gas)
	bigGasPrice := new(big.Int).SetUint64(gasPrice)
	tGasPrice := (*hexutil.Big)(bigGasPrice)
	hValue := (*hexutil.Big)(value)
	tData := hexutil.Bytes(data)
	tInput := hexutil.Bytes(input)

	unsignedTx := api.SendTxArgs{
		From:      from,
		Recipient: &to,
		GasLimit:  &tGas,
		Price:     tGasPrice,
		Amount:    hValue,
		//Nonce : nonce,	Nonce will be determined by Klay node.
		Data:    &tData,
		Payload: &tInput,
	}

	if err := ec.c.CallContext(ctx, &hex, "klay_sendTransaction", toSendTxArgs(unsignedTx)); err != nil {
		return common.Hash{}, err
	}
	hash := common.BytesToHash(hex)
	return hash, nil
}

// ImportRawKey can create key store from private key string on Klaytn node.
func (ec *Client) ImportRawKey(ctx context.Context, key string, password string) (common.Address, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "personal_importRawKey", key, password)
	address := common.BytesToAddress(result)
	return address, err
}

// UnlockAccount can unlock the account on Klaytn node.
func (ec *Client) UnlockAccount(ctx context.Context, address common.Address, password string, time uint) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "personal_unlockAccount", address, password, time)
	return result, err
}

func toCallArg(msg klaytn.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

func toSendTxArgs(msg api.SendTxArgs) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.Recipient,
	}
	if *msg.GasLimit != 0 {
		arg["gas"] = (*hexutil.Uint64)(msg.GasLimit)
	}
	if msg.Price != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.Price)
	}
	if msg.Amount != nil {
		arg["value"] = (*hexutil.Big)(msg.Amount)
	}
	if len(*msg.Data) > 0 {
		arg["data"] = (*hexutil.Bytes)(msg.Data)
	}
	if len(*msg.Payload) > 0 {
		arg["input"] = (*hexutil.Bytes)(msg.Payload)
	}

	return arg
}

// BlockNumber can get the latest block number.
func (ec *Client) BlockNumber(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "klay_blockNumber")
	return (*big.Int)(&result), err
}

// ChainID can return the chain ID of the chain.
func (ec *Client) ChainID(ctx context.Context) (*big.Int, error) {
	if ec.chainID != nil {
		return ec.chainID, nil
	}

	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "klay_chainID")
	if err == nil {
		ec.chainID = (*big.Int)(&result)
	}
	return ec.chainID, err
}

// AddPeer can add a static peer on Klaytn node.
func (ec *Client) AddPeer(ctx context.Context, url string) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "admin_addPeer", url)
	return result, err
}

// RemovePeer can remove a static peer on Klaytn node.
func (ec *Client) RemovePeer(ctx context.Context, url string) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "admin_removePeer", url)
	return result, err
}

// CreateAccessList tries to create an access list for a specific transaction based on the
// current pending state of the blockchain.
func (ec *Client) CreateAccessList(ctx context.Context, msg klaytn.CallMsg) (*types.AccessList, uint64, string, error) {
	type AccessListResult struct {
		Accesslist *types.AccessList `json:"accessList"`
		Error      string            `json:"error,omitempty"`
		GasUsed    hexutil.Uint64    `json:"gasUsed"`
	}
	var result AccessListResult
	if err := ec.c.CallContext(ctx, &result, "klay_createAccessList", toCallArg(msg)); err != nil {
		return nil, 0, "", err
	}
	return result.Accesslist, uint64(result.GasUsed), result.Error, nil
}
