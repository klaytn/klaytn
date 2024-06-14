// Modifications Copyright 2024 The Kaia Authors
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
// Modified and improved for the Kaia development.

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
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
	block, err := s.b.BlockByHash(ctx, blockHash)
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
		return newRPCTransactionFromBlockIndex(block, uint64(index), s.b.ChainConfig()), nil
	}
	return nil, err
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, blockHash)
	if block != nil && err == nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index), s.b.ChainConfig()), nil
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
	block, err := s.b.BlockByHash(ctx, blockHash)
	if block != nil && err == nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index)), nil
	}
	return nil, err
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number or hash
func (s *PublicTransactionPoolAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	// Ask transaction pool for the nonce which includes pending transactions
	if blockNr, ok := blockNrOrHash.Number(); ok && blockNr == rpc.PendingBlockNumber {
		nonce := s.b.GetPoolNonce(ctx, address)
		return (*hexutil.Uint64)(&nonce), nil
	}

	// Resolve block number and use its state to ask for the nonce
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
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
		header, err := s.b.HeaderByHash(ctx, blockHash)
		if err != nil {
			return nil
		}
		return newRPCTransaction(header, tx, blockHash, blockNumber, index, s.b.ChainConfig())
	}
	// No finalized transaction, try to retrieve it from the pool
	if tx := s.b.GetPoolTransaction(hash); tx != nil {
		return newRPCPendingTransaction(tx, s.b.ChainConfig())
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

	if !tx.Type().IsChainDataAnchoring() {
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
func RpcOutputReceipt(header *types.Header, tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, receipt *types.Receipt, config *params.ChainConfig) map[string]interface{} {
	if tx == nil || receipt == nil {
		return nil
	}

	fields := newRPCTransaction(header, tx, blockHash, blockNumber, index, config)

	if receipt.Status != types.ReceiptStatusSuccessful {
		fields["status"] = hexutil.Uint(types.ReceiptStatusFailed)
		fields["txError"] = hexutil.Uint(receipt.Status)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}

	fields["logsBloom"] = receipt.Bloom
	fields["gasUsed"] = hexutil.Uint64(receipt.GasUsed)

	fields["effectiveGasPrice"] = hexutil.Uint64(tx.EffectiveGasPrice(header, config).Uint64())

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
	tx, blockHash, blockNumber, index, receipt := s.b.GetTxLookupInfoAndReceipt(ctx, hash)
	return s.getTransactionReceipt(ctx, tx, blockHash, blockNumber, index, receipt)
}

// GetTransactionReceiptInCache returns the transaction receipt for the given transaction hash.
func (s *PublicTransactionPoolAPI) GetTransactionReceiptInCache(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	tx, blockHash, blockNumber, index, receipt := s.b.GetTxLookupInfoAndReceiptInCache(hash)
	return s.getTransactionReceipt(ctx, tx, blockHash, blockNumber, index, receipt)
}

// getTransactionReceipt returns the transaction receipt for the given transaction hash.
func (s *PublicTransactionPoolAPI) getTransactionReceipt(ctx context.Context, tx *types.Transaction, blockHash common.Hash,
	blockNumber uint64, index uint64, receipt *types.Receipt,
) (map[string]interface{}, error) {
	// No error handling is required here.
	// Header is checked in the following RpcOutputReceipt function
	header, _ := s.b.HeaderByHash(ctx, blockHash)
	return RpcOutputReceipt(header, tx, blockHash, blockNumber, index, receipt, s.b.ChainConfig()), nil
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

// signAsFeePayer is a helper function that signs a transaction as a fee payer with the private key of the given address.
func (s *PublicTransactionPoolAPI) signAsFeePayer(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Request the wallet to sign the transaction
	return wallet.SignTxAsFeePayer(account, tx, s.b.ChainConfig().ChainID)
}

var submitTxCount = 0

// submitTransaction is a helper function that submits tx to txPool and logs a message.
func submitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
	// submitTxCount++
	// log.Error("### submitTransaction","tx",submitTxCount)

	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}
	// TODO-Kaia only enable on logging
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
	if args.AccountNonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	signedTx, err := s.SignTransaction(ctx, args)
	if err != nil {
		return common.Hash{}, err
	}

	return submitTransaction(ctx, s.b, signedTx.Tx)
}

// SendTransactionAsFeePayer creates a transaction for the given argument, sign it as a fee payer
// and submit it to the transaction pool.
func (s *PublicTransactionPoolAPI) SendTransactionAsFeePayer(ctx context.Context, args SendTxArgs) (common.Hash, error) {
	// Don't allow dynamic assign of values from the setDefaults function since the sender already signed on specific values.
	if args.TypeInt == nil {
		return common.Hash{}, errTxArgNilTxType
	}
	if args.AccountNonce == nil {
		return common.Hash{}, errTxArgNilNonce
	}
	if args.GasLimit == nil {
		return common.Hash{}, errTxArgNilGas
	}
	if args.Price == nil {
		return common.Hash{}, errTxArgNilGasPrice
	}

	if args.TxSignatures == nil {
		return common.Hash{}, errTxArgNilSenderSig
	}

	feePayerSignedTx, err := s.SignTransactionAsFeePayer(ctx, args)
	if err != nil {
		return common.Hash{}, err
	}

	return submitTransaction(ctx, s.b, feePayerSignedTx.Tx)
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
		signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
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
	if args.TypeInt != nil && args.TypeInt.IsEthTypedTransaction() {
		if args.Price == nil && (args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil) {
			return nil, fmt.Errorf("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
		}
	}

	// No need to obtain the noncelock mutex, since we won't be sending this
	// tx into the transaction pool, but right back to the user
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx, err := args.toTransaction()
	if err != nil {
		return nil, err
	}
	signedTx, err := s.sign(args.From, tx)
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signedTx}, nil
}

// SignTransactionAsFeePayer will sign the given transaction as a fee payer
// with the from account. The node needs to have the private key of the account
// corresponding with the given from address and it needs to be unlocked.
func (s *PublicTransactionPoolAPI) SignTransactionAsFeePayer(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error) {
	// Allows setting a default nonce value of the sender just for the case the fee payer tries to sign a tx earlier than the sender.
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx, err := args.toTransaction()
	if err != nil {
		return nil, err
	}
	// Don't return errors for nil signature allowing the fee payer to sign a tx earlier than the sender.
	if args.TxSignatures != nil {
		tx.SetSignature(args.TxSignatures.ToTxSignatures())
	}
	feePayer, err := tx.FeePayer()
	if err != nil {
		return nil, errTxArgInvalidFeePayer
	}
	feePayerSignedTx, err := s.signAsFeePayer(feePayer, tx)
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(feePayerSignedTx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, feePayerSignedTx}, nil
}

func getAccountsFromWallets(wallets []accounts.Wallet) map[common.Address]struct{} {
	accounts := make(map[common.Address]struct{})
	for _, wallet := range wallets {
		for _, account := range wallet.Accounts() {
			accounts[account.Address] = struct{}{}
		}
	}
	return accounts
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (s *PublicTransactionPoolAPI) PendingTransactions() ([]map[string]interface{}, error) {
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	accounts := getAccountsFromWallets(s.b.AccountManager().Wallets())
	transactions := make([]map[string]interface{}, 0, len(pending))
	for _, tx := range pending {
		from := getFrom(tx)
		if _, exists := accounts[from]; exists {
			transactions = append(transactions, newRPCPendingTransaction(tx, s.b.ChainConfig()))
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (s *PublicTransactionPoolAPI) Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	return resend(s, ctx, &sendArgs, gasPrice, gasLimit)
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func resend(s *PublicTransactionPoolAPI, ctx context.Context, sendArgs NewTxArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	if sendArgs.nonce() == nil {
		return common.Hash{}, fmt.Errorf("missing transaction nonce in transaction spec")
	}
	if err := sendArgs.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	matchTx, err := sendArgs.toTransaction()
	if err != nil {
		return common.Hash{}, err
	}
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return common.Hash{}, err
	}

	for _, p := range pending {
		signer := types.LatestSignerForChainID(p.ChainId())
		wantSigHash := signer.Hash(matchTx)

		if pFrom, err := types.Sender(signer, p); err == nil && pFrom == sendArgs.from() && signer.Hash(p) == wantSigHash {
			// Match. Re-sign and send the transaction.
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				sendArgs.setGasPrice(gasPrice)
			}
			if gasLimit != nil && *gasLimit != 0 {
				sendArgs.setGas(gasLimit)
			}
			tx, err := sendArgs.toTransaction()
			if err != nil {
				return common.Hash{}, err
			}
			signedTx, err := s.sign(sendArgs.from(), tx)
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

// RecoverFromTransaction recovers the sender address from a signed raw transaction.
// The signature is validated against the sender account's key configuration at the given block number.
func (s *PublicTransactionPoolAPI) RecoverFromTransaction(ctx context.Context, encodedTx hexutil.Bytes, blockNumber rpc.BlockNumber) (common.Address, error) {
	if len(encodedTx) == 0 {
		return common.Address{}, fmt.Errorf("Empty input")
	}
	if 0 < encodedTx[0] && encodedTx[0] < 0x8 { // TODO: Add helper to distinguish eth vs Kaia txtypes
		ethEnvelope := []byte{byte(types.EthereumTxTypeEnvelope)}
		encodedTx = append(ethEnvelope, encodedTx...)
	}

	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return common.Address{}, err
	}

	var bn uint64
	if blockNumber == rpc.LatestBlockNumber || blockNumber == rpc.PendingBlockNumber {
		bn = s.b.CurrentBlock().NumberU64()
	} else {
		bn = blockNumber.Uint64()
	}

	signer := types.MakeSigner(s.b.ChainConfig(), new(big.Int).SetUint64(bn))
	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return common.Address{}, err
	}
	_, err = tx.ValidateSender(signer, state, bn)
	if err != nil {
		return common.Address{}, err
	}
	return tx.From()
}

// RecoverFromMessage validates that the message is signed by one of the keys in the given account.
// Returns the recovered signer address, which may be different from the account address.
func (s *PublicTransactionPoolAPI) RecoverFromMessage(
	ctx context.Context, address common.Address, data, sig hexutil.Bytes, blockNumber rpc.BlockNumber,
) (common.Address, error) {
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, fmt.Errorf("invalid Klaytn signature (V is not 27 or 28)")
	}

	// Transform yellow paper V from 27/28 to 0/1
	sig[crypto.RecoveryIDOffset] -= 27

	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return common.Address{}, err
	}
	key := state.GetKey(address)

	// We cannot identify if the signature has signed with kaia or eth prefix without the signer's address.
	// Even though a user signed message with eth prefix, it will return invalid something in klayEcRecover.
	// We should call each rcrecover function separately and the actual result will be checked in ValidateMember.
	var recoverErr error
	if pubkey, err := klayEcRecover(data, sig); err == nil {
		if key.ValidateMember(pubkey, address) {
			return crypto.PubkeyToAddress(*pubkey), nil
		}
	} else {
		recoverErr = err
	}
	if pubkey, err := ethEcRecover(data, sig); err == nil {
		if key.ValidateMember(pubkey, address) {
			return crypto.PubkeyToAddress(*pubkey), nil
		}
	} else {
		recoverErr = err
	}
	if recoverErr != nil {
		return common.Address{}, recoverErr
	} else {
		return common.Address{}, errors.New("Invalid signature")
	}
}
