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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/ser/rlp"
	"math/big"
)

// PublicTransactionPoolAPI exposes methods for the RPC interface
type PublicTransactionPoolAPI struct {
	b         Backend
	nonceLock *AddrLocker
}

// NewPublicTransactionPoolAPI creates a new RPC service with methods specific for the transaction pool.
func NewPublicTransactionPoolAPI(b Backend, nonceLock *AddrLocker) *PublicTransactionPoolAPI {
	return &PublicTransactionPoolAPI{b, nonceLock}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*hexutil.Uint, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil && err == nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n, err
	}
	return nil, err
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) (*hexutil.Uint, error) {
	block, err := s.b.GetBlock(ctx, blockHash)
	if block != nil && err == nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n, err
	}
	return nil, err
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil && err == nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index)), nil
	}
	return nil, err
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.GetBlock(ctx, blockHash)
	if block != nil && err == nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index)), nil
	}
	return nil, err
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (hexutil.Bytes, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil && err == nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index)), nil
	}
	return nil, err
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (hexutil.Bytes, error) {
	block, err := s.b.GetBlock(ctx, blockHash)
	if block != nil && err == nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index)), nil
	}
	return nil, err
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (s *PublicTransactionPoolAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (*hexutil.Uint64, error) {
	// Ask transaction pool for the nonce which includes pending transactions
	if blockNr == rpc.PendingBlockNumber {
		nonce := s.b.GetPoolNonce(ctx, address)
		return (*hexutil.Uint64)(&nonce), nil
	}

	// Ask NonceCache for the nonce if blockNumber is lastestBlockNumber
	if blockNr == rpc.LatestBlockNumber {
		nonce, ok := s.b.GetNonceInCache(address)
		if ok {
			return (*hexutil.Uint64)(&nonce), nil
		}
	}

	// Resolve block number and use its state to ask for the nonce
	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
	if err != nil {
		return nil, err
	}
	nonce := state.GetNonce(address)
	return (*hexutil.Uint64)(&nonce), state.Error()
}

func (s *PublicTransactionPoolAPI) GetTransactionBySenderTxHash(ctx context.Context, senderTxHash common.Hash) map[string]interface{} {
	txhash := s.b.ChainDB().ReadTxHashFromSenderTxHash(senderTxHash)
	if common.EmptyHash(txhash) {
		txhash = senderTxHash
	}
	return s.GetTransactionByHash(ctx, txhash)
}

// GetTransactionByHash returns the transaction for the given hash
func (s *PublicTransactionPoolAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	// Try to return an already finalized transaction
	if tx, blockHash, blockNumber, index := s.b.ChainDB().ReadTxAndLookupInfo(hash); tx != nil {
		return newRPCTransaction(tx, blockHash, blockNumber, index)
	}
	// No finalized transaction, try to retrieve it from the pool
	if tx := s.b.GetPoolTransaction(hash); tx != nil {
		return newRPCPendingTransaction(tx)
	}
	// Transaction unknown, return as such
	return nil
}

// GetDecodedAnchoringTransactionByHash returns the decoded anchoring data of anchoring transaction for the given hash
func (s *PublicTransactionPoolAPI) GetDecodedAnchoringTransactionByHash(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	var tx *types.Transaction

	// Try to return an already finalized transaction
	if tx, _, _, _ = s.b.ChainDB().ReadTxAndLookupInfo(hash); tx != nil {
		goto decode
	}

	// No finalized transaction, try to retrieve it from the pool
	if tx = s.b.GetPoolTransaction(hash); tx != nil {
		goto decode
	}
	return nil, errors.New("can't find the transaction")

decode:

	if tx.Type() != types.TxTypeChainDataAnchoring {
		return nil, errors.New("invalid transaction type")
	}

	data, err := tx.AnchoredData()
	if err != nil {
		return nil, err
	}

	anchoringDataInternal, err := types.DecodeAnchoringData(data)
	if err != nil {
		return nil, err
	}

	str, err := json.Marshal(anchoringDataInternal)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal(str, &result)

	return result, nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (s *PublicTransactionPoolAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	var tx *types.Transaction

	// Retrieve a finalized transaction, or a pooled otherwise
	if tx, _, _, _ = s.b.ChainDB().ReadTxAndLookupInfo(hash); tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			// Transaction not found anywhere, abort
			return nil, nil
		}
	}

	// Serialize to RLP and return
	return rlp.EncodeToBytes(tx)
}

// RpcOutputReceipt converts a receipt to the RPC output with the associated information regarding to the
// block in which the receipt is included, the transaction that outputs the receipt, and the receipt itself.
func RpcOutputReceipt(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, receipt *types.Receipt) map[string]interface{} {
	if tx == nil || receipt == nil {
		return nil
	}
	fields := newRPCTransaction(tx, blockHash, blockNumber, index)

	if receipt.Status != types.ReceiptStatusSuccessful {
		fields["status"] = hexutil.Uint(types.ReceiptStatusFailed)
		fields["txError"] = hexutil.Uint(receipt.Status)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}

	fields["logsBloom"] = receipt.Bloom
	fields["gasUsed"] = hexutil.Uint64(receipt.GasUsed)

	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	} else {
		fields["logs"] = receipt.Logs
	}
	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	} else {
		fields["contractAddress"] = nil

	}

	// Rename field name `hash` to `transactionHash` since this function returns a JSON object of a receipt.
	fields["transactionHash"] = fields["hash"]
	delete(fields, "hash")

	return fields
}

func (s *PublicTransactionPoolAPI) GetTransactionReceiptBySenderTxHash(ctx context.Context, senderTxHash common.Hash) (map[string]interface{}, error) {
	txhash := s.b.ChainDB().ReadTxHashFromSenderTxHash(senderTxHash)
	if common.EmptyHash(txhash) {
		txhash = senderTxHash
	}
	return s.GetTransactionReceipt(ctx, txhash)
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (s *PublicTransactionPoolAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	return RpcOutputReceipt(s.b.GetTxLookupInfoAndReceipt(ctx, hash)), nil
}

// GetTransactionReceiptInCache returns the transaction receipt for the given transaction hash.
func (s *PublicTransactionPoolAPI) GetTransactionReceiptInCache(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	return RpcOutputReceipt(s.b.GetTxLookupInfoAndReceiptInCache(hash)), nil
}

// sign is a helper function that signs a transaction with the private key of the given address.
func (s *PublicTransactionPoolAPI) sign(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Request the wallet to sign the transaction
	return wallet.SignTx(account, tx, s.b.ChainConfig().ChainID)
}

type NewTxArgs interface {
	setDefaults(context.Context, Backend) error
	toTransaction() (*types.Transaction, error)
	from() common.Address
}

type ValueTransferTxArgs struct {
	From     common.Address  `json:"from"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	To       common.Address  `json:"to"`
	Value    *hexutil.Big    `json:"value"`
}

func (args *ValueTransferTxArgs) from() common.Address {
	return args.From
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *ValueTransferTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.Gas == nil {
		args.Gas = new(hexutil.Uint64)
		*(*uint64)(args.Gas) = 90000
	}
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Nonce == nil {
		nonce := b.GetPoolNonce(ctx, args.From)
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	return nil
}

func (args *ValueTransferTxArgs) toTransaction() (*types.Transaction, error) {
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    (uint64)(*args.Nonce),
		types.TxValueKeyGasLimit: (uint64)(*args.Gas),
		types.TxValueKeyGasPrice: (*big.Int)(args.GasPrice),
		types.TxValueKeyFrom:     args.From,
		types.TxValueKeyTo:       args.To,
		types.TxValueKeyAmount:   (*big.Int)(args.Value),
	})

	if err != nil {
		return nil, err
	}

	return tx, nil
}

type AccountUpdateTxArgs struct {
	From     common.Address  `json:"from"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	Key      *hexutil.Bytes  `json:"key"`
}

func (args *AccountUpdateTxArgs) from() common.Address {
	return args.From
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *AccountUpdateTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.Gas == nil {
		args.Gas = new(hexutil.Uint64)
		*(*uint64)(args.Gas) = 90000
	}
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Nonce == nil {
		nonce := b.GetPoolNonce(ctx, args.From)
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	return nil
}

func (args *AccountUpdateTxArgs) toTransaction() (*types.Transaction, error) {
	serializer := accountkey.NewAccountKeySerializer()

	if err := rlp.DecodeBytes(*args.Key, &serializer); err != nil {
		return nil, err
	}
	tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      (uint64)(*args.Nonce),
		types.TxValueKeyGasLimit:   (uint64)(*args.Gas),
		types.TxValueKeyGasPrice:   (*big.Int)(args.GasPrice),
		types.TxValueKeyFrom:       args.From,
		types.TxValueKeyAccountKey: serializer.GetKey(),
	})

	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *SendTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.Gas == nil {
		args.Gas = new(hexutil.Uint64)
		*(*uint64)(args.Gas) = 90000
	}
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	if args.Nonce == nil {
		nonce := b.GetPoolNonce(ctx, args.From)
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return errors.New(`Both "data" and "input" are set and not equal. Please use "input" to pass transaction call data.`)
	}
	if args.To == nil {
		// Contract creation
		var input []byte
		if args.Data != nil {
			input = *args.Data
		} else if args.Input != nil {
			input = *args.Input
		}
		if len(input) == 0 {
			return errors.New(`contract creation without any data provided`)
		}
	}
	return nil
}

func (args *SendTxArgs) toTransaction() *types.Transaction {
	var input []byte
	if args.Data != nil {
		input = *args.Data
	} else if args.Input != nil {
		input = *args.Input
	}
	if args.To == nil {
		return types.NewContractCreation(uint64(*args.Nonce), (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input)
	}
	return types.NewTransaction(uint64(*args.Nonce), *args.To, (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input)
}

var submitTxCount = 0

// submitTransaction is a helper function that submits tx to txPool and logs a message.
func submitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
	//submitTxCount++
	//log.Error("### submitTransaction","tx",submitTxCount)

	if tx.Type().IsRPCExcluded() {
		return common.Hash{}, fmt.Errorf("%s cannot be submitted via RPC!", tx.Type().String())
	}

	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}
	// TODO-Klaytn only enable on logging
	//if tx.To() == nil {
	//	signer := types.MakeSigner(b.ChainConfig(), b.CurrentBlock().Number())
	//	from, err := types.Sender(signer, tx)
	//	if err != nil {
	//		logger.Error("api.submitTransaction make from","err",err)
	//		return common.Hash{}, err
	//	}
	//	addr := crypto.CreateAddress(from, tx.Nonce())
	//	logger.Info("Submitted contract creation", "fullhash", tx.Hash().Hex(), "contract", addr.Hex())
	//} else {
	//	logger.Info("Submitted transaction", "fullhash", tx.Hash().Hex(), "recipient", tx.To())
	//}
	return tx.Hash(), nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (s *PublicTransactionPoolAPI) SendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: args.From}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}

	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	// Assemble the transaction and sign with the wallet
	tx := args.toTransaction()

	signed, err := wallet.SignTx(account, tx, s.b.ChainConfig().ChainID)
	if err != nil {
		return common.Hash{}, err
	}
	return submitTransaction(ctx, s.b, signed)
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTransactionPoolAPI) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return common.Hash{}, err
	}
	return submitTransaction(ctx, s.b, tx)
}

// Sign calculates an ECDSA signature for:
// keccack256("\x19Klaytn Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (s *PublicTransactionPoolAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Sign the requested hash with the wallet
	signature, err := wallet.SignHash(account, signHash(data))
	if err == nil {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	return signature, err
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (s *PublicTransactionPoolAPI) SignTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error) {
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx, err := s.sign(args.From, args.toTransaction())
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, tx}, nil
}

// PendingTransactions returns the transactions that are in the transaction pool and have a from address that is one of
// the accounts this node manages.
func (s *PublicTransactionPoolAPI) PendingTransactions() ([]map[string]interface{}, error) {
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}

	transactions := make([]map[string]interface{}, 0, len(pending))
	for _, tx := range pending {
		signer := types.NewEIP155Signer(tx.ChainId())
		from, _ := types.Sender(signer, tx)
		if _, err := s.b.AccountManager().Find(accounts.Account{Address: from}); err == nil {
			transactions = append(transactions, newRPCPendingTransaction(tx))
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (s *PublicTransactionPoolAPI) Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	if sendArgs.Nonce == nil {
		return common.Hash{}, fmt.Errorf("missing transaction nonce in transaction spec")
	}
	if err := sendArgs.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	matchTx := sendArgs.toTransaction()
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return common.Hash{}, err
	}

	for _, p := range pending {
		signer := types.NewEIP155Signer(p.ChainId())
		wantSigHash := signer.Hash(matchTx)

		if pFrom, err := types.Sender(signer, p); err == nil && pFrom == sendArgs.From && signer.Hash(p) == wantSigHash {
			// Match. Re-sign and send the transaction.
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				sendArgs.GasPrice = gasPrice
			}
			if gasLimit != nil && *gasLimit != 0 {
				sendArgs.Gas = gasLimit
			}
			signedTx, err := s.sign(sendArgs.From, sendArgs.toTransaction())
			if err != nil {
				return common.Hash{}, err
			}
			if err = s.b.SendTx(ctx, signedTx); err != nil {
				return common.Hash{}, err
			}
			return signedTx.Hash(), nil
		}
	}

	return common.Hash{}, fmt.Errorf("Transaction %#x not found", matchTx.Hash())
}
