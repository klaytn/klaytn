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
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/klaytn/klaytn/node/cn/filters"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

const defaultGasPrice = 25 * params.Ston
const localTxExecutionTime = 5 * time.Second

var logger = log.NewModuleLogger(log.API)

// PublicBlockChainAPI provides an API to access the Klaytn blockchain.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicBlockChainAPI struct {
	b Backend
}

// NewPublicBlockChainAPI creates a new Klaytn blockchain API.
func NewPublicBlockChainAPI(b Backend) *PublicBlockChainAPI {
	return &PublicBlockChainAPI{b}
}

// BlockNumber returns the block number of the chain head.
func (s *PublicBlockChainAPI) BlockNumber() hexutil.Uint64 {
	header, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber) // latest header should always be available
	return hexutil.Uint64(header.Number.Uint64())
}

// ChainID returns the chain ID of the chain from genesis file.
func (s *PublicBlockChainAPI) ChainID() *hexutil.Big {
	return s.ChainId()
}

// ChainId returns the chain ID of the chain from genesis file.
// This is for compatibility with ethereum client
func (s *PublicBlockChainAPI) ChainId() *hexutil.Big {
	if s.b.ChainConfig() != nil {
		return (*hexutil.Big)(s.b.ChainConfig().ChainID)
	}
	return nil
}

// IsContractAccount returns true if the account associated with addr has a non-empty codeHash.
// It returns false otherwise.
func (s *PublicBlockChainAPI) IsContractAccount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (bool, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return false, err
	}
	return state.IsContractAccount(address), state.Error()
}

// IsHumanReadable returns true if the account associated with addr is a human-readable account.
// It returns false otherwise.
//func (s *PublicBlockChainAPI) IsHumanReadable(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (bool, error) {
//	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
//	if err != nil {
//		return false, err
//	}
//	return state.IsHumanReadable(address), state.Error()
//}

// GetBlockReceipts returns all the transaction receipts for the given block hash.
func (s *PublicBlockChainAPI) GetBlockReceipts(ctx context.Context, blockHash common.Hash) ([]map[string]interface{}, error) {
	receipts := s.b.GetBlockReceipts(ctx, blockHash)
	block, err := s.b.BlockByHash(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	txs := block.Transactions()
	if receipts.Len() != txs.Len() {
		return nil, fmt.Errorf("the size of transactions and receipts is different in the block (%s)", blockHash.String())
	}
	fieldsList := make([]map[string]interface{}, 0, len(receipts))
	for index, receipt := range receipts {
		fields := RpcOutputReceipt(txs[index], blockHash, block.NumberU64(), uint64(index), receipt)
		fieldsList = append(fieldsList, fields)
	}
	return fieldsList, nil
}

// GetBalance returns the amount of peb for the given address in the state of the
// given block number or hash. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers and hash are also allowed.
func (s *PublicBlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(state.GetBalance(address)), state.Error()
}

// AccountCreated returns true if the account associated with the address is created.
// It returns false otherwise.
func (s *PublicBlockChainAPI) AccountCreated(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (bool, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return false, err
	}
	return state.Exist(address), state.Error()
}

// GetAccount returns account information of an input address.
func (s *PublicBlockChainAPI) GetAccount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*account.AccountSerializer, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return &account.AccountSerializer{}, err
	}
	acc := state.GetAccount(address)
	if acc == nil {
		return &account.AccountSerializer{}, err
	}
	serAcc := account.NewAccountSerializerWithAccount(acc)
	return serAcc, state.Error()
}

// rpcMarshalHeader converts the given header to the RPC output.
func (s *PublicBlockChainAPI) rpcMarshalHeader(header *types.Header) map[string]interface{} {
	fields := filters.RPCMarshalHeader(header, s.b.ChainConfig().IsEthTxTypeForkEnabled(header.Number))
	return fields
}

// GetHeaderByNumber returns the requested canonical block header.
func (s *PublicBlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	header, err := s.b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return s.rpcMarshalHeader(header), nil
}

// GetHeaderByHash returns the requested header by hash.
func (s *PublicBlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	header, err := s.b.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return s.rpcMarshalHeader(header), nil
}

// GetBlockByNumber returns the requested block. When blockNr is -1 the chain head is returned. When fullTx is true all
// transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI) GetBlockByNumber(ctx context.Context, blockNr rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil && err == nil {
		response, err := s.rpcOutputBlock(block, true, fullTx)
		if err == nil && blockNr == rpc.PendingBlockNumber {
			// Pending blocks need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI) GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	return s.rpcOutputBlock(block, true, fullTx)
}

// GetCode returns the code stored at the given address in the state for the given block number or hash.
func (s *PublicBlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers and hash are also allowed.
func (s *PublicBlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	res := state.GetState(address, common.HexToHash(key))
	return res[:], state.Error()
}

// GetAccountKey returns the account key of EOA at a given address.
// If the account of the given address is a Legacy Account or a Smart Contract Account, it will return nil.
func (s *PublicBlockChainAPI) GetAccountKey(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*accountkey.AccountKeySerializer, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return &accountkey.AccountKeySerializer{}, err
	}
	if state.Exist(address) == false {
		return nil, nil
	}
	accountKey := state.GetKey(address)
	serAccKey := accountkey.NewAccountKeySerializerWithAccountKey(accountKey)
	return serAccKey, state.Error()
}

// IsParallelDBWrite returns if parallel write is enabled or not.
// If enabled, data written in WriteBlockWithState is being written in parallel manner.
func (s *PublicBlockChainAPI) IsParallelDBWrite() bool {
	return s.b.IsParallelDBWrite()
}

// IsSenderTxHashIndexingEnabled returns if senderTxHash to txHash mapping information
// indexing is enabled or not.
func (s *PublicBlockChainAPI) IsSenderTxHashIndexingEnabled() bool {
	return s.b.IsSenderTxHashIndexingEnabled()
}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      hexutil.Uint64  `json:"gas"`
	GasPrice hexutil.Big     `json:"gasPrice"`
	Value    hexutil.Big     `json:"value"`
	Data     hexutil.Bytes   `json:"data"`
}

func DoCall(ctx context.Context, b Backend, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash, vmCfg vm.Config, timeout time.Duration, globalGasCap *big.Int) ([]byte, uint64, uint64, bool, error) {
	defer func(start time.Time) { logger.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, 0, 0, false, err
	}
	// Set sender address or use a default if none specified
	addr := args.From
	if addr == (common.Address{}) {
		if wallets := b.AccountManager().Wallets(); len(wallets) > 0 {
			if accounts := wallets[0].Accounts(); len(accounts) > 0 {
				addr = accounts[0].Address
			}
		}
	}
	// Set default gas & gas price if none were set
	gas, gasPrice := uint64(args.Gas), args.GasPrice.ToInt()
	if gas == 0 {
		gas = math.MaxUint64 / 2
	}
	if globalGasCap != nil && globalGasCap.Uint64() < gas {
		logger.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap.Uint64()
	}
	if gasPrice.Sign() == 0 {
		gasPrice = new(big.Int).SetUint64(defaultGasPrice)
	}

	intrinsicGas, err := types.IntrinsicGas(args.Data, nil, args.To == nil, b.ChainConfig().Rules(header.Number))
	if err != nil {
		return nil, 0, 0, false, err
	}

	// Create new call message
	msg := types.NewMessage(addr, args.To, 0, args.Value.ToInt(), gas, gasPrice, args.Data, false, intrinsicGas)

	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Add gas fee to sender for estimating gasLimit/computing cost or calling a function by insufficient balance sender.
	state.AddBalance(msg.ValidatedSender(), new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), msg.GasPrice()))
	// Get a new instance of the EVM.
	evm, vmError, err := b.GetEVM(ctx, msg, state, header, vmCfg)
	if err != nil {
		return nil, 0, 0, false, err
	}
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel(vm.CancelByCtxDone)
	}()

	res, gas, kerr := blockchain.ApplyMessage(evm, msg)
	err = kerr.ErrTxInvalid
	if err := vmError(); err != nil {
		return nil, 0, 0, false, err
	}
	// If the timer caused an abort, return an appropriate error message
	if evm.Cancelled() {
		return nil, 0, 0, false, fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}

	// Propagate error of Receipt as JSON RPC error
	if err == nil {
		err = blockchain.GetVMerrFromReceiptStatus(kerr.Status)
	}

	return res, gas, evm.GetOpCodeComputationCost(), kerr.Status != types.ReceiptStatusSuccessful, err
}

// Call executes the given transaction on the state for the given block number or hash.
// It doesn't make and changes in the state/blockchain and is useful to execute and retrieve values.
func (s *PublicBlockChainAPI) Call(ctx context.Context, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	result, _, _, _, err := DoCall(ctx, s.b, args, blockNrOrHash, vm.Config{}, localTxExecutionTime, s.b.RPCGasCap())
	return (hexutil.Bytes)(result), err
}

func (s *PublicBlockChainAPI) EstimateComputationCost(ctx context.Context, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	_, _, computationCost, _, err := DoCall(ctx, s.b, args, blockNrOrHash, vm.Config{UseOpcodeComputationCost: true}, localTxExecutionTime, s.b.RPCGasCap())
	return (hexutil.Uint64)(computationCost), err
}

// EstimateGas returns an estimate of the amount of gas needed to execute the given transaction against the latest block.
func (s *PublicBlockChainAPI) EstimateGas(ctx context.Context, args CallArgs) (hexutil.Uint64, error) {
	return s.DoEstimateGas(ctx, s.b, args, s.b.RPCGasCap())
}

func (s *PublicBlockChainAPI) DoEstimateGas(ctx context.Context, b Backend, args CallArgs, gasCap *big.Int) (hexutil.Uint64, error) {
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	if uint64(args.Gas) >= params.TxGas {
		hi = uint64(args.Gas)
	} else {
		// Retrieve the current pending block to act as the gas ceiling
		hi = params.UpperGasLimit
	}
	if gasCap != nil && hi > gasCap.Uint64() {
		logger.Warn("Caller gas above allowance, capping", "requested", hi, "cap", gasCap)
		hi = gasCap.Uint64()
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) bool {
		args.Gas = hexutil.Uint64(gas)

		_, _, _, failed, err := DoCall(ctx, b, args, rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber), vm.Config{UseOpcodeComputationCost: true}, localTxExecutionTime, gasCap)
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
	return hexutil.Uint64(hi), nil
}

// ExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
// execution status, the amount of gas used and the return value
type ExecutionResult struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// accessListResult returns an optional accesslist
// Its the result of the `debug_createAccessList` RPC call.
// It contains an error if the transaction itself failed.
type AccessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

// CreateAccessList creates a EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.
// TODO-Klaytn: Have to implement logic. For now, Klaytn does not implement actual access list logic, so return empty access list result.
func (s *PublicBlockChainAPI) CreateAccessList(ctx context.Context, args SendTxArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*AccessListResult, error) {
	result := &AccessListResult{Accesslist: &types.AccessList{}, GasUsed: hexutil.Uint64(0)}
	return result, nil
}

// StructLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type StructLogRes struct {
	Pc      uint64             `json:"pc"`
	Op      string             `json:"op"`
	Gas     uint64             `json:"gas"`
	GasCost uint64             `json:"gasCost"`
	Depth   int                `json:"depth"`
	Error   error              `json:"error,omitempty"`
	Stack   *[]string          `json:"stack,omitempty"`
	Memory  *[]string          `json:"memory,omitempty"`
	Storage *map[string]string `json:"storage,omitempty"`
}

// formatLogs formats EVM returned structured logs for json output
func FormatLogs(logs []vm.StructLog) []StructLogRes {
	formatted := make([]StructLogRes, len(logs))
	for index, trace := range logs {
		formatted[index] = StructLogRes{
			Pc:      trace.Pc,
			Op:      trace.Op.String(),
			Gas:     trace.Gas,
			GasCost: trace.GasCost,
			Depth:   trace.Depth,
			Error:   trace.Err,
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = fmt.Sprintf("%x", math.PaddedBigBytes(stackValue, 32))
			}
			formatted[index].Stack = &stack
		}
		if trace.Memory != nil {
			memory := make([]string, 0, (len(trace.Memory)+31)/32)
			for i := 0; i+32 <= len(trace.Memory); i += 32 {
				memory = append(memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
			}
			formatted[index].Memory = &memory
		}
		if trace.Storage != nil {
			storage := make(map[string]string)
			for i, storageValue := range trace.Storage {
				storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
			}
			formatted[index].Storage = &storage
		}
	}
	return formatted
}

func RpcOutputBlock(b *types.Block, td *big.Int, inclTx bool, fullTx bool, isEnabledEthTxTypeFork bool) (map[string]interface{}, error) {
	head := b.Header() // copies the header once
	fields := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             b.Hash(),
		"parentHash":       head.ParentHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"reward":           head.Rewardbase,
		"blockscore":       (*hexutil.Big)(head.BlockScore),
		"totalBlockScore":  (*hexutil.Big)(td),
		"extraData":        hexutil.Bytes(head.Extra),
		"governanceData":   hexutil.Bytes(head.Governance),
		"voteData":         hexutil.Bytes(head.Vote),
		"size":             hexutil.Uint64(b.Size()),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        (*hexutil.Big)(head.Time),
		"timestampFoS":     (hexutil.Uint)(head.TimeFoS),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}

	if inclTx {
		formatTx := func(tx *types.Transaction) (interface{}, error) {
			return tx.Hash(), nil
		}

		if fullTx {
			formatTx = func(tx *types.Transaction) (interface{}, error) {
				return newRPCTransactionFromBlockHash(b, tx.Hash()), nil
			}
		}

		txs := b.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range b.Transactions() {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}

	if isEnabledEthTxTypeFork {
		fields["baseFeePerGas"] = (*hexutil.Big)(new(big.Int).SetUint64(params.BaseFee))
	}

	return fields, nil
}

// rpcOutputBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func (s *PublicBlockChainAPI) rpcOutputBlock(b *types.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	return RpcOutputBlock(b, s.b.GetTd(b.Hash()), inclTx, fullTx, s.b.ChainConfig().IsEthTxTypeForkEnabled(b.Header().Number))
}

func getFrom(tx *types.Transaction) common.Address {
	var from common.Address
	if tx.IsEthereumTransaction() {
		signer := types.LatestSignerForChainID(tx.ChainId())
		from, _ = types.Sender(signer, tx)
	} else {
		from, _ = tx.From()
	}
	return from
}

// newRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func newRPCTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64) map[string]interface{} {
	from := getFrom(tx)

	output := tx.MakeRPCOutput()

	output["senderTxHash"] = tx.SenderTxHashAll()
	output["blockHash"] = blockHash
	output["blockNumber"] = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
	output["from"] = from
	output["hash"] = tx.Hash()
	output["transactionIndex"] = hexutil.Uint(index)

	return output
}

// newRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func newRPCPendingTransaction(tx *types.Transaction) map[string]interface{} {
	return newRPCTransaction(tx, common.Hash{}, 0, 0)
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex(b *types.Block, index uint64) map[string]interface{} {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return newRPCTransaction(txs[index], b.Hash(), b.NumberU64(), index)
}

// newRPCRawTransactionFromBlockIndex returns the bytes of a transaction given a block and a transaction index.
func newRPCRawTransactionFromBlockIndex(b *types.Block, index uint64) hexutil.Bytes {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	blob, _ := rlp.EncodeToBytes(txs[index])
	return blob
}

// newRPCTransactionFromBlockHash returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockHash(b *types.Block, hash common.Hash) map[string]interface{} {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newRPCTransactionFromBlockIndex(b, uint64(idx))
		}
	}
	return nil
}
