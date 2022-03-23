// Copyright 2021 The klaytn Authors
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

package api

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/klaytn/klaytn/rlp"

	"github.com/klaytn/klaytn/accounts/abi"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node/cn/filters"
	"github.com/klaytn/klaytn/params"
)

const (
	// EmptySha3Uncles always have value which is the result of
	// `crypto.Keccak256Hash(rlp.EncodeToBytes([]*types.Header(nil)).String())`
	// because there is no uncles in Klaytn.
	// Just use const value because we don't have to calculate it everytime which always be same result.
	EmptySha3Uncles = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	// ZeroHashrate exists for supporting Ethereum compatible data structure.
	// There is no POW mining mechanism in Klaytn.
	ZeroHashrate uint64 = 0
	// ZeroUncleCount is always zero because there is no uncle blocks in Klaytn.
	ZeroUncleCount uint = 0
)

var (
	errNoMiningWork = errors.New("no mining work available yet")
)

// EthereumAPI provides an API to access the Klaytn through the `eth` namespace.
// TODO-Klaytn: Removed unused variable
type EthereumAPI struct {
	publicFilterAPI   *filters.PublicFilterAPI
	governanceKlayAPI *governance.GovernanceKlayAPI

	publicKlayAPI            *PublicKlayAPI
	publicBlockChainAPI      *PublicBlockChainAPI
	publicTransactionPoolAPI *PublicTransactionPoolAPI
	publicAccountAPI         *PublicAccountAPI
	publicGovernanceAPI      *governance.PublicGovernanceAPI
}

// NewEthereumAPI creates a new ethereum API.
// EthereumAPI operates using Klaytn's API internally without overriding.
// Therefore, it is necessary to use APIs defined in two different packages(cn and api),
// so those apis will be defined through a setter.
func NewEthereumAPI() *EthereumAPI {
	return &EthereumAPI{nil, nil, nil, nil, nil, nil, nil}
}

// SetPublicFilterAPI sets publicFilterAPI
func (api *EthereumAPI) SetPublicFilterAPI(publicFilterAPI *filters.PublicFilterAPI) {
	api.publicFilterAPI = publicFilterAPI
}

// SetGovernanceKlayAPI sets governanceKlayAPI
func (api *EthereumAPI) SetGovernanceKlayAPI(governanceKlayAPI *governance.GovernanceKlayAPI) {
	api.governanceKlayAPI = governanceKlayAPI
}

// SetPublicKlayAPI sets publicKlayAPI
func (api *EthereumAPI) SetPublicKlayAPI(publicKlayAPI *PublicKlayAPI) {
	api.publicKlayAPI = publicKlayAPI
}

// SetPublicBlockChainAPI sets publicBlockChainAPI
func (api *EthereumAPI) SetPublicBlockChainAPI(publicBlockChainAPI *PublicBlockChainAPI) {
	api.publicBlockChainAPI = publicBlockChainAPI
}

// SetPublicTransactionPoolAPI sets publicTransactionPoolAPI
func (api *EthereumAPI) SetPublicTransactionPoolAPI(publicTransactionPoolAPI *PublicTransactionPoolAPI) {
	api.publicTransactionPoolAPI = publicTransactionPoolAPI
}

// SetPublicAccountAPI sets publicAccountAPI
func (api *EthereumAPI) SetPublicAccountAPI(publicAccountAPI *PublicAccountAPI) {
	api.publicAccountAPI = publicAccountAPI
}

// SetPublicGovernanceAPI sets publicGovernanceAPI
func (api *EthereumAPI) SetPublicGovernanceAPI(publicGovernanceAPI *governance.PublicGovernanceAPI) {
	api.publicGovernanceAPI = publicGovernanceAPI
}

// Etherbase is the address of operating node.
// Unlike Ethereum, it only returns the node address because Klaytn does not have a POW mechanism.
func (api *EthereumAPI) Etherbase() (common.Address, error) {
	return api.publicGovernanceAPI.NodeAddress(), nil
}

// Coinbase is the address of operating node (alias for Etherbase).
func (api *EthereumAPI) Coinbase() (common.Address, error) {
	return api.Etherbase()
}

// Hashrate returns the POW hashrate.
// Unlike Ethereum, it always returns ZeroHashrate because Klaytn does not have a POW mechanism.
func (api *EthereumAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(ZeroHashrate)
}

// Mining returns an indication if this node is currently mining.
// Unlike Ethereum, it always returns false because Klaytn does not have a POW mechanism,
func (api *EthereumAPI) Mining() bool {
	return false
}

// GetWork returns an errNoMiningWork because klaytn does not have a POW mechanism.
func (api *EthereumAPI) GetWork() ([4]string, error) {
	return [4]string{}, errNoMiningWork
}

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

// SubmitWork returns false because klaytn does not have a POW mechanism.
func (api *EthereumAPI) SubmitWork(nonce BlockNonce, hash, digest common.Hash) bool {
	return false
}

// SubmitHashrate returns false because klaytn does not have a POW mechanism.
func (api *EthereumAPI) SubmitHashrate(rate hexutil.Uint64, id common.Hash) bool {
	return false
}

// GetHashrate returns ZeroHashrate because klaytn does not have a POW mechanism.
func (api *EthereumAPI) GetHashrate() uint64 {
	return ZeroHashrate
}

// NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
// as transactions enter the pending state.
//
// It is part of the filter package because this filter can be used through the
// `eth_getFilterChanges` polling method that is also used for log filters.
//
// https://eth.wiki/json-rpc/API#eth_newpendingtransactionfilter
func (api *EthereumAPI) NewPendingTransactionFilter() rpc.ID {
	return api.publicFilterAPI.NewPendingTransactionFilter()
}

// NewPendingTransactions creates a subscription that is triggered each time a transaction
// enters the transaction pool and was signed from one of the transactions this nodes manages.
func (api *EthereumAPI) NewPendingTransactions(ctx context.Context) (*rpc.Subscription, error) {
	return api.publicFilterAPI.NewPendingTransactions(ctx)
}

// NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
// It is part of the filter package since polling goes with eth_getFilterChanges.
//
// https://eth.wiki/json-rpc/API#eth_newblockfilter
func (api *EthereumAPI) NewBlockFilter() rpc.ID {
	return api.publicFilterAPI.NewBlockFilter()
}

// NewHeads send a notification each time a new (header) block is appended to the chain.
func (api *EthereumAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()
	go func() {
		headers := make(chan *types.Header)
		headersSub := api.publicFilterAPI.Events().SubscribeNewHeads(headers)

		for {
			select {
			case h := <-headers:
				header, err := api.rpcMarshalHeader(h)
				if err != nil {
					logger.Error("Failed to marshal header during newHeads subscription", "err", err)
					headersSub.Unsubscribe()
					return
				}
				notifier.Notify(rpcSub.ID, header)
			case <-rpcSub.Err():
				headersSub.Unsubscribe()
				return
			case <-notifier.Closed():
				headersSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}

// Logs creates a subscription that fires for all new log that match the given filter criteria.
func (api *EthereumAPI) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
	return api.publicFilterAPI.Logs(ctx, crit)
}

// NewFilter creates a new filter and returns the filter id. It can be
// used to retrieve logs when the state changes. This method cannot be
// used to fetch logs that are already stored in the state.
//
// Default criteria for the from and to block are "latest".
// Using "latest" as block number will return logs for mined blocks.
// Using "pending" as block number returns logs for not yet mined (pending) blocks.
// In case logs are removed (chain reorg) previously returned logs are returned
// again but with the removed property set to true.
//
// In case "fromBlock" > "toBlock" an error is returned.
//
// https://eth.wiki/json-rpc/API#eth_newfilter
func (api *EthereumAPI) NewFilter(crit filters.FilterCriteria) (rpc.ID, error) {
	return api.publicFilterAPI.NewFilter(crit)
}

// GetLogs returns logs matching the given argument that are stored within the state.
//
// https://eth.wiki/json-rpc/API#eth_getlogs
func (api *EthereumAPI) GetLogs(ctx context.Context, crit filters.FilterCriteria) ([]*types.Log, error) {
	return api.publicFilterAPI.GetLogs(ctx, crit)
}

// UninstallFilter removes the filter with the given filter id.
//
// https://eth.wiki/json-rpc/API#eth_uninstallfilter
func (api *EthereumAPI) UninstallFilter(id rpc.ID) bool {
	return api.publicFilterAPI.UninstallFilter(id)
}

// GetFilterLogs returns the logs for the filter with the given id.
// If the filter could not be found an empty array of logs is returned.
//
// https://eth.wiki/json-rpc/API#eth_getfilterlogs
func (api *EthereumAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([]*types.Log, error) {
	return api.publicFilterAPI.GetFilterLogs(ctx, id)
}

// GetFilterChanges returns the logs for the filter with the given id since
// last time it was called. This can be used for polling.
//
// For pending transaction and block filters the result is []common.Hash.
// (pending)Log filters return []Log.
//
// https://eth.wiki/json-rpc/API#eth_getfilterchanges
func (api *EthereumAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	return api.publicFilterAPI.GetFilterChanges(id)
}

// GasPrice returns a suggestion for a gas price.
func (api *EthereumAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	return api.publicKlayAPI.GasPrice(ctx)
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (api *EthereumAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	return api.publicKlayAPI.MaxPriorityFeePerGas(ctx)
}

// DecimalOrHex unmarshals a non-negative decimal or hex parameter into a uint64.
type DecimalOrHex uint64

// UnmarshalJSON implements json.Unmarshaler.
func (dh *DecimalOrHex) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	value, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		value, err = hexutil.DecodeUint64(input)
	}
	if err != nil {
		return err
	}
	*dh = DecimalOrHex(value)
	return nil
}

func (api *EthereumAPI) FeeHistory(ctx context.Context, blockCount DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	return api.publicKlayAPI.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (api *EthereumAPI) Syncing() (interface{}, error) {
	return api.publicKlayAPI.Syncing()
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (api *EthereumAPI) ChainId() (*hexutil.Big, error) {
	return api.publicBlockChainAPI.ChainId(), nil
}

// BlockNumber returns the block number of the chain head.
func (api *EthereumAPI) BlockNumber() hexutil.Uint64 {
	return api.publicBlockChainAPI.BlockNumber()
}

// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (api *EthereumAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	return api.publicBlockChainAPI.GetBalance(ctx, address, blockNrOrHash)
}

// EthAccountResult structs for GetProof
// AccountResult in go-ethereum has been renamed to EthAccountResult.
// AccountResult is defined in go-ethereum's internal package, so AccountResult is redefined here as EthAccountResult.
type EthAccountResult struct {
	Address      common.Address     `json:"address"`
	AccountProof []string           `json:"accountProof"`
	Balance      *hexutil.Big       `json:"balance"`
	CodeHash     common.Hash        `json:"codeHash"`
	Nonce        hexutil.Uint64     `json:"nonce"`
	StorageHash  common.Hash        `json:"storageHash"`
	StorageProof []EthStorageResult `json:"storageProof"`
}

// StorageResult in go-ethereum has been renamed to EthStorageResult.
// StorageResult is defined in go-ethereum's internal package, so StorageResult is redefined here as EthStorageResult.
type EthStorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// GetProof returns the Merkle-proof for a given account and optionally some storage keys.
// This feature is not supported in Klaytn yet. It just returns account information from state trie.
func (api *EthereumAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*EthAccountResult, error) {
	state, _, err := api.publicKlayAPI.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	storageTrie := state.StorageTrie(address)
	storageHash := types.EmptyRootHash
	codeHash := state.GetCodeHash(address)
	storageProof := make([]EthStorageResult, len(storageKeys))

	// if we have a storageTrie, (which means the account exists), we can update the storagehash
	if storageTrie != nil {
		storageHash = storageTrie.Hash()
	} else {
		// no storageTrie means the account does not exist, so the codeHash is the hash of an empty bytearray.
		codeHash = crypto.Keccak256Hash(nil)
	}

	return &EthAccountResult{
		Address:      address,
		AccountProof: []string{},
		Balance:      (*hexutil.Big)(state.GetBalance(address)),
		CodeHash:     codeHash,
		Nonce:        hexutil.Uint64(state.GetNonce(address)),
		StorageHash:  storageHash,
		StorageProof: storageProof,
	}, state.Error()
}

// GetHeaderByNumber returns the requested canonical block header.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
func (api *EthereumAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	// In Ethereum, err is always nil because the backend of Ethereum always return nil.
	klaytnHeader, err := api.publicBlockChainAPI.b.HeaderByNumber(ctx, number)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}
		return nil, err
	}
	response, err := api.rpcMarshalHeader(klaytnHeader)
	if err != nil {
		return nil, err
	}
	if number == rpc.PendingBlockNumber {
		// Pending header need to nil out a few fields
		for _, field := range []string{"hash", "nonce", "miner"} {
			response[field] = nil
		}
	}
	return response, nil
}

// GetHeaderByHash returns the requested header by hash.
func (api *EthereumAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	// In Ethereum, err is always nil because the backend of Ethereum always return nil.
	klaytnHeader, _ := api.publicBlockChainAPI.b.HeaderByHash(ctx, hash)
	if klaytnHeader != nil {
		response, err := api.rpcMarshalHeader(klaytnHeader)
		if err != nil {
			return nil
		}
		return response
	}
	return nil
}

// GetBlockByNumber returns the requested canonical block.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
// * When fullTx is true all transactions in the block are returned, otherwise
//   only the transaction hash is returned.
func (api *EthereumAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	// Klaytn backend returns error when there is no matched block but
	// Ethereum returns it as nil without error, so we should return is as nil when there is no matched block.
	klaytnBlock, err := api.publicBlockChainAPI.b.BlockByNumber(ctx, number)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}
		return nil, err
	}
	response, err := api.rpcMarshalBlock(klaytnBlock, true, fullTx)
	if err == nil && number == rpc.PendingBlockNumber {
		// Pending blocks need to nil out a few fields
		for _, field := range []string{"hash", "nonce", "miner"} {
			response[field] = nil
		}
	}
	return response, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (api *EthereumAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	// Klaytn backend returns error when there is no matched block but
	// Ethereum returns it as nil without error, so we should return is as nil when there is no matched block.
	klaytnBlock, err := api.publicBlockChainAPI.b.BlockByHash(ctx, hash)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}
		return nil, err
	}
	return api.rpcMarshalBlock(klaytnBlock, true, fullTx)
}

// GetUncleByBlockNumberAndIndex returns nil because there is no uncle block in Klaytn.
func (api *EthereumAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	return nil, nil
}

// GetUncleByBlockHashAndIndex returns nil because there is no uncle block in Klaytn.
func (api *EthereumAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	return nil, nil
}

// GetUncleCountByBlockNumber returns 0 when given blockNr exists because there is no uncle block in Klaytn.
func (api *EthereumAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := api.publicBlockChainAPI.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(ZeroUncleCount)
		return &n
	}
	return nil
}

// GetUncleCountByBlockHash returns 0 when given blockHash exists because there is no uncle block in Klaytn.
func (api *EthereumAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := api.publicBlockChainAPI.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(ZeroUncleCount)
		return &n
	}
	return nil
}

// GetCode returns the code stored at the given address in the state for the given block number.
func (api *EthereumAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return api.publicBlockChainAPI.GetCode(ctx, address, blockNrOrHash)
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers are also allowed.
func (api *EthereumAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return api.publicBlockChainAPI.GetStorageAt(ctx, address, key, blockNrOrHash)
}

// EthOverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
// OverrideAccount in go-ethereum has been renamed to EthOverrideAccount.
// OverrideAccount is defined in go-ethereum's internal package, so OverrideAccount is redefined here as EthOverrideAccount.
type EthOverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// EthStateOverride is the collection of overridden accounts.
// StateOverride in go-ethereum has been renamed to EthStateOverride.
// StateOverride is defined in go-ethereum's internal package, so StateOverride is redefined here as EthStateOverride.
type EthStateOverride map[common.Address]EthOverrideAccount

func (diff *EthStateOverride) Apply(state *state.StateDB) error {
	if diff == nil {
		return nil
	}
	for addr, account := range *diff {
		// Override account nonce.
		if account.Nonce != nil {
			state.SetNonce(addr, uint64(*account.Nonce))
		}
		// Override account(contract) code.
		if account.Code != nil {
			state.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil {
			state.SetBalance(addr, (*big.Int)(*account.Balance))
		}
		if account.State != nil && account.StateDiff != nil {
			return fmt.Errorf("account %s has both 'state' and 'stateDiff'", addr.Hex())
		}
		// Replace entire state if caller requires.
		if account.State != nil {
			state.SetStorage(addr, *account.State)
		}
		// Apply state diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				state.SetState(addr, key, value)
			}
		}
	}
	return nil
}

// Call executes the given transaction on the state for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the state/blockchain and is
// useful to execute and retrieve values.
func (api *EthereumAPI) Call(ctx context.Context, args EthTransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *EthStateOverride) (hexutil.Bytes, error) {
	gasCap := uint64(0)
	if rpcGasCap := api.publicBlockChainAPI.b.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}
	result, _, status, err := EthDoCall(ctx, api.publicBlockChainAPI.b, args, blockNrOrHash, overrides, localTxExecutionTime, gasCap)
	if err != nil {
		return nil, err
	}

	err = blockchain.GetVMerrFromReceiptStatus(status)
	if err != nil && isReverted(err) && len(result) > 0 {
		return nil, newRevertError(result)
	}
	return common.CopyBytes(result), err
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (api *EthereumAPI) EstimateGas(ctx context.Context, args EthTransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	bNrOrHash := rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	gasCap := uint64(0)
	if rpcGasCap := api.publicBlockChainAPI.b.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}
	return EthDoEstimateGas(ctx, api.publicBlockChainAPI.b, args, bNrOrHash, gasCap)
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (api *EthereumAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	transactionCount, _ := api.publicTransactionPoolAPI.GetBlockTransactionCountByNumber(ctx, blockNr)
	return transactionCount
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (api *EthereumAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	transactionCount, _ := api.publicTransactionPoolAPI.GetBlockTransactionCountByHash(ctx, blockHash)
	return transactionCount
}

// CreateAccessList creates a EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.
func (api *EthereumAPI) CreateAccessList(ctx context.Context, args EthTransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*AccessListResult, error) {
	// To use CreateAccess of PublicBlockChainAPI, we need to convert the EthTransactionArgs to SendTxArgs.
	// However, since SendTxArgs does not yet support MaxFeePerGas and MaxPriorityFeePerGas, the conversion logic is bound to be incomplete.
	// Since this parameter is not actually used and currently only returns an empty result value, implement the logic to return an empty result separately,
	// and later, when the API is actually implemented, add the relevant fields to SendTxArgs and call the function in PublicBlockChainAPI.
	// TODO-Klaytn: Modify below logic to use api.publicBlockChainAPI.CreateAccessList
	result := &AccessListResult{Accesslist: &types.AccessList{}, GasUsed: hexutil.Uint64(0)}
	return result, nil
}

// EthRPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
// RPCTransaction in go-ethereum has been renamed to EthRPCTransaction.
// RPCTransaction is defined in go-ethereum's internal package, so RPCTransaction is redefined here as EthRPCTransaction.
type EthRPCTransaction struct {
	BlockHash        *common.Hash      `json:"blockHash"`
	BlockNumber      *hexutil.Big      `json:"blockNumber"`
	From             common.Address    `json:"from"`
	Gas              hexutil.Uint64    `json:"gas"`
	GasPrice         *hexutil.Big      `json:"gasPrice"`
	GasFeeCap        *hexutil.Big      `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big      `json:"maxPriorityFeePerGas,omitempty"`
	Hash             common.Hash       `json:"hash"`
	Input            hexutil.Bytes     `json:"input"`
	Nonce            hexutil.Uint64    `json:"nonce"`
	To               *common.Address   `json:"to"`
	TransactionIndex *hexutil.Uint64   `json:"transactionIndex"`
	Value            *hexutil.Big      `json:"value"`
	Type             hexutil.Uint64    `json:"type"`
	Accesses         *types.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big      `json:"chainId,omitempty"`
	V                *hexutil.Big      `json:"v"`
	R                *hexutil.Big      `json:"r"`
	S                *hexutil.Big      `json:"s"`
}

// ethTxJSON is the JSON representation of Ethereum transaction.
// ethTxJSON is used by eth namespace APIs which returns Transaction object as it is.
// Because every transaction in Klaytn, implements json.Marshaler interface (MarshalJSON), but
// it is marshaled for Klaytn format only.
// e.g. Ethereum transaction have V, R, and S field for signature but,
// Klaytn transaction have types.TxSignaturesJSON which includes array of signatures which is not
// applicable for Ethereum transaction.
type ethTxJSON struct {
	Type hexutil.Uint64 `json:"type"`

	// Common transaction fields:
	Nonce                *hexutil.Uint64 `json:"nonce"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	Value                *hexutil.Big    `json:"value"`
	Data                 *hexutil.Bytes  `json:"input"`
	V                    *hexutil.Big    `json:"v"`
	R                    *hexutil.Big    `json:"r"`
	S                    *hexutil.Big    `json:"s"`
	To                   *common.Address `json:"to"`

	// Access list transaction fields:
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
	AccessList *types.AccessList `json:"accessList,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

// newEthRPCTransactionFromBlockIndex creates an EthRPCTransaction from block and index parameters.
func newEthRPCTransactionFromBlockIndex(b *types.Block, index uint64) *EthRPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		logger.Error("invalid transaction index", "given index", index, "length of txs", len(txs))
		return nil
	}
	return newEthRPCTransaction(txs[index], b.Hash(), b.NumberU64(), index)
}

// newEthRPCTransactionFromBlockHash returns a transaction that will serialize to the RPC representation.
func newEthRPCTransactionFromBlockHash(b *types.Block, hash common.Hash) *EthRPCTransaction {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newEthRPCTransactionFromBlockIndex(b, uint64(idx))
		}
	}
	return nil
}

// resolveToField returns value which fits to `to` field based on transaction types.
// This function is used when converting Klaytn transactions to Ethereum transaction types.
func resolveToField(tx *types.Transaction) *common.Address {
	switch tx.Type() {
	case types.TxTypeAccountUpdate, types.TxTypeFeeDelegatedAccountUpdate, types.TxTypeFeeDelegatedAccountUpdateWithRatio,
		types.TxTypeCancel, types.TxTypeFeeDelegatedCancel, types.TxTypeFeeDelegatedCancelWithRatio,
		types.TxTypeChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		// These type of transactions actually do not have `to` address, but Ethereum always have `to` field,
		// so we Klaytn developers decided to fill the `to` field with `from` address value in these case.
		from := getFrom(tx)
		return &from
	}
	return tx.To()
}

// newEthRPCTransaction creates an EthRPCTransaction from Klaytn transaction.
func newEthRPCTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber, index uint64) *EthRPCTransaction {
	// When an unknown transaction is requested through rpc call,
	// nil is returned by Klaytn API, and it is handled.
	if tx == nil {
		return nil
	}

	typeInt := tx.Type()
	// If tx is not Ethereum transaction, the type is converted to TxTypeLegacyTransaction.
	if !tx.IsEthereumTransaction() {
		typeInt = types.TxTypeLegacyTransaction
	}

	signature := tx.GetTxInternalData().RawSignatureValues()[0]

	result := &EthRPCTransaction{
		Type:     hexutil.Uint64(byte(typeInt)),
		From:     getFrom(tx),
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    tx.Data(),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       resolveToField(tx),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(signature.V),
		R:        (*hexutil.Big)(signature.R),
		S:        (*hexutil.Big)(signature.S),
	}

	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	switch typeInt {
	case types.TxTypeEthereumAccessList:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case types.TxTypeEthereumDynamicFee:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// TODO-Klaytn: If we change the gas price policy from fixed to dynamic,
		// we should change the params.BaseFee(fixed value) to dynamic value.
		price := math.BigMin(new(big.Int).Add(tx.GasTipCap(), new(big.Int).SetUint64(params.BaseFee)), tx.GasFeeCap())
		result.GasPrice = (*hexutil.Big)(price)
	}
	return result
}

// newEthRPCPendingTransaction creates an EthRPCTransaction for pending tx.
func newEthRPCPendingTransaction(tx *types.Transaction) *EthRPCTransaction {
	return newEthRPCTransaction(tx, common.Hash{}, 0, 0)
}

// formatTxToEthTxJSON formats types.Transaction to ethTxJSON.
// Use this function for only Ethereum typed transaction.
func formatTxToEthTxJSON(tx *types.Transaction) *ethTxJSON {
	var enc ethTxJSON

	enc.Type = hexutil.Uint64(byte(tx.Type()))
	// If tx is not Ethereum transaction, the type is converted to TxTypeLegacyTransaction.
	if !tx.IsEthereumTransaction() {
		enc.Type = hexutil.Uint64(types.TxTypeLegacyTransaction)
	}
	enc.Hash = tx.Hash()
	signature := tx.GetTxInternalData().RawSignatureValues()[0]
	// Initialize signature values when it is nil.
	if signature.V == nil {
		signature.V = new(big.Int)
	}
	if signature.R == nil {
		signature.R = new(big.Int)
	}
	if signature.S == nil {
		signature.S = new(big.Int)
	}
	nonce := tx.Nonce()
	gas := tx.Gas()
	enc.Nonce = (*hexutil.Uint64)(&nonce)
	enc.Gas = (*hexutil.Uint64)(&gas)
	enc.Value = (*hexutil.Big)(tx.Value())
	data := tx.Data()
	enc.Data = (*hexutil.Bytes)(&data)
	enc.To = tx.To()
	enc.V = (*hexutil.Big)(signature.V)
	enc.R = (*hexutil.Big)(signature.R)
	enc.S = (*hexutil.Big)(signature.S)

	switch tx.Type() {
	case types.TxTypeEthereumAccessList:
		al := tx.AccessList()
		enc.AccessList = &al
		enc.ChainID = (*hexutil.Big)(tx.ChainId())
		enc.GasPrice = (*hexutil.Big)(tx.GasPrice())
	case types.TxTypeEthereumDynamicFee:
		al := tx.AccessList()
		enc.AccessList = &al
		enc.ChainID = (*hexutil.Big)(tx.ChainId())
		enc.MaxFeePerGas = (*hexutil.Big)(tx.GasFeeCap())
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(tx.GasTipCap())
	default:
		enc.GasPrice = (*hexutil.Big)(tx.GasPrice())
	}
	return &enc
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (api *EthereumAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *EthRPCTransaction {
	block, err := api.publicTransactionPoolAPI.b.BlockByNumber(ctx, blockNr)
	if err != nil {
		return nil
	}

	return newEthRPCTransactionFromBlockIndex(block, uint64(index))
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (api *EthereumAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *EthRPCTransaction {
	block, err := api.publicTransactionPoolAPI.b.BlockByHash(ctx, blockHash)
	if err != nil || block == nil {
		return nil
	}
	return newEthRPCTransactionFromBlockIndex(block, uint64(index))
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (api *EthereumAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	rawTx, err := api.publicTransactionPoolAPI.GetRawTransactionByBlockNumberAndIndex(ctx, blockNr, index)
	if err != nil {
		return nil
	}
	if rawTx[0] == byte(types.EthereumTxTypeEnvelope) {
		rawTx = rawTx[1:]
	}
	return rawTx
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (api *EthereumAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	rawTx, err := api.publicTransactionPoolAPI.GetRawTransactionByBlockHashAndIndex(ctx, blockHash, index)
	if err != nil {
		return nil
	}
	if rawTx[0] == byte(types.EthereumTxTypeEnvelope) {
		rawTx = rawTx[1:]
	}
	return rawTx
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number.
func (api *EthereumAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	return api.publicTransactionPoolAPI.GetTransactionCount(ctx, address, blockNrOrHash)
}

// GetTransactionByHash returns the transaction for the given hash.
func (api *EthereumAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*EthRPCTransaction, error) {
	// Try to return an already finalized transaction
	if tx, blockHash, blockNumber, index := api.publicTransactionPoolAPI.b.ChainDB().ReadTxAndLookupInfo(hash); tx != nil {
		return newEthRPCTransaction(tx, blockHash, blockNumber, index), nil
	}
	// No finalized transaction, try to retrieve it from the pool
	if tx := api.publicTransactionPoolAPI.b.GetPoolTransaction(hash); tx != nil {
		return newEthRPCPendingTransaction(tx), nil
	}
	// Transaction unknown, return as such
	return nil, nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (api *EthereumAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	rawTx, err := api.publicTransactionPoolAPI.GetRawTransactionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if rawTx[0] == byte(types.EthereumTxTypeEnvelope) {
		return rawTx[1:], nil
	}
	return rawTx, nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (api *EthereumAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	// Formats return Klaytn Transaction Receipt to the Ethereum Transaction Receipt.
	tx, blockHash, blockNumber, index, receipt := api.publicTransactionPoolAPI.b.GetTxLookupInfoAndReceipt(ctx, hash)

	if tx == nil {
		return nil, nil
	}
	receipts := api.publicTransactionPoolAPI.b.GetBlockReceipts(ctx, blockHash)
	cumulativeGasUsed := uint64(0)
	for i := uint64(0); i <= index; i++ {
		cumulativeGasUsed += receipts[i].GasUsed
	}

	ethTx, err := newEthTransactionReceipt(tx, api.publicTransactionPoolAPI.b, blockHash, blockNumber, index, cumulativeGasUsed, receipt)
	if err != nil {
		return nil, err
	}
	return ethTx, nil
}

// newEthTransactionReceipt creates a transaction receipt in Ethereum format.
func newEthTransactionReceipt(tx *types.Transaction, b Backend, blockHash common.Hash, blockNumber, index, cumulativeGasUsed uint64, receipt *types.Receipt) (map[string]interface{}, error) {
	// When an unknown transaction receipt is requested through rpc call,
	// nil is returned by Klaytn API, and it is handled.
	if tx == nil || receipt == nil {
		return nil, nil
	}

	typeInt := tx.Type()
	// If tx is not Ethereum transaction, the type is converted to TxTypeLegacyTransaction.
	if !tx.IsEthereumTransaction() {
		typeInt = types.TxTypeLegacyTransaction
	}

	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   tx.Hash(),
		"transactionIndex":  hexutil.Uint64(index),
		"from":              getFrom(tx),
		"to":                resolveToField(tx),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(byte(typeInt)),
	}

	if b.ChainConfig().IsEthTxTypeForkEnabled(new(big.Int).SetUint64(blockNumber)) {
		// TODO-Klaytn: Klaytn is using fixed BaseFee(0) as now but
		// if we apply dynamic BaseFee, we should add calculated BaseFee instead of using params.BaseFee.
		baseFee := new(big.Int).SetUint64(params.BaseFee)
		fields["effectiveGasPrice"] = hexutil.Uint64(tx.EffectiveGasPrice(baseFee).Uint64())
	} else {
		fields["effectiveGasPrice"] = hexutil.Uint64(tx.GasPrice().Uint64())
	}

	// Always use the "status" field and Ignore the "root" field.
	if receipt.Status != types.ReceiptStatusSuccessful {
		// In Ethereum, status field can have 0(=Failure) or 1(=Success) only.
		fields["status"] = hexutil.Uint(types.ReceiptStatusFailed)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}

	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	}
	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}

	return fields, nil
}

// EthTransactionArgs represents the arguments to construct a new transaction
// or a message call.
// TransactionArgs in go-ethereum has been renamed to EthTransactionArgs.
// TransactionArgs is defined in go-ethereum's internal package, so TransactionArgs is redefined here as EthTransactionArgs.
type EthTransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

// from retrieves the transaction sender address.
func (args *EthTransactionArgs) from() common.Address {
	if args.From == nil {
		return common.Address{}
	}
	return *args.From
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *EthTransactionArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

// setDefaults fills in default values for unspecified tx fields.
func (args *EthTransactionArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// After london, default to 1559 uncles gasPrice is set
	head := b.CurrentBlock().Header()
	// TODO-Klaytn: Klaytn is using fixed BaseFee(0) as now but
	// if we apply dynamic BaseFee, we should add calculated BaseFee instead of using params.BaseFee.
	fixedBaseFee := new(big.Int).SetUint64(params.BaseFee)
	// Klaytn uses fixed gasPrice policy determined by Governance, so
	// only fixedGasPrice value is allowed to be used as args.MaxFeePerGas and args.MaxPriorityFeePerGas.
	fixedGasPrice, err := b.SuggestPrice(ctx)
	if err != nil {
		return err
	}

	// If user specifies both maxPriorityFee and maxFee, then we do not
	// need to consult the chain for defaults. It's definitely a London tx.
	if args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil {
		if b.ChainConfig().IsEthTxTypeForkEnabled(head.Number) && args.GasPrice == nil {
			if args.MaxPriorityFeePerGas == nil {
				// TODO-Klaytn: Original logic of Ethereum uses b.SuggestTipCap which suggests TipCap, not a GasPrice.
				// But Klaytn currently uses fixed unit price determined by Governance, so using b.SuggestPrice
				// is fine as now.
				args.MaxPriorityFeePerGas = (*hexutil.Big)(fixedGasPrice)
			}
			if args.MaxFeePerGas == nil {
				// TODO-Klaytn: Calculating formula of gasFeeCap is same with Ethereum except for
				// using fixedBaseFee which means gasFeeCap is always same with args.MaxPriorityFeePerGas as now.
				gasFeeCap := new(big.Int).Add(
					(*big.Int)(args.MaxPriorityFeePerGas),
					new(big.Int).Mul(fixedBaseFee, big.NewInt(2)),
				)
				args.MaxFeePerGas = (*hexutil.Big)(gasFeeCap)
			}
			if args.MaxPriorityFeePerGas.ToInt().Cmp(fixedGasPrice) != 0 || args.MaxFeePerGas.ToInt().Cmp(fixedGasPrice) != 0 {
				return fmt.Errorf("only %s is allowed to be used as maxFeePerGas and maxPriorityPerGas", fixedGasPrice.Text(16))
			}
			if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
				return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
			}
		} else {
			if args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil {
				return errors.New("maxFeePerGas or maxPriorityFeePerGas specified but london is not active yet")
			}
			if args.GasPrice == nil {
				// TODO-Klaytn: Original logic of Ethereum uses b.SuggestTipCap which suggests TipCap, not a GasPrice.
				// But Klaytn currently uses fixed unit price determined by Governance, so using b.SuggestPrice
				// is fine as now.
				if b.ChainConfig().IsEthTxTypeForkEnabled(head.Number) {
					// TODO-Klaytn: Klaytn is using fixed BaseFee(0) as now but
					// if we apply dynamic BaseFee, we should add calculated BaseFee instead of params.BaseFee.
					fixedGasPrice.Add(fixedGasPrice, new(big.Int).SetUint64(params.BaseFee))
				}
				args.GasPrice = (*hexutil.Big)(fixedGasPrice)
			}
		}
	} else {
		// Both maxPriorityFee and maxFee set by caller. Sanity-check their internal relation
		if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
		}
	}
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	if args.Nonce == nil {
		nonce := b.GetPoolNonce(ctx, args.from())
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
	}
	if args.To == nil && len(args.data()) == 0 {
		return errors.New(`contract creation without any data provided`)
	}
	// Estimate the gas usage if necessary.
	if args.Gas == nil {
		// These fields are immutable during the estimation, safe to
		// pass the pointer directly.
		data := args.data()
		callArgs := EthTransactionArgs{
			From:                 args.From,
			To:                   args.To,
			GasPrice:             args.GasPrice,
			MaxFeePerGas:         args.MaxFeePerGas,
			MaxPriorityFeePerGas: args.MaxPriorityFeePerGas,
			Value:                args.Value,
			Data:                 (*hexutil.Bytes)(&data),
			AccessList:           args.AccessList,
		}
		pendingBlockNr := rpc.NewBlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
		gasCap := uint64(0)
		if rpcGasCap := b.RPCGasCap(); rpcGasCap != nil {
			gasCap = rpcGasCap.Uint64()
		}
		estimated, err := EthDoEstimateGas(ctx, b, callArgs, pendingBlockNr, gasCap)
		if err != nil {
			return err
		}
		args.Gas = &estimated
		logger.Trace("Estimate gas usage automatically", "gas", args.Gas)
	}
	if args.ChainID == nil {
		id := (*hexutil.Big)(b.ChainConfig().ChainID)
		args.ChainID = id
	}
	return nil
}

// ToMessage change EthTransactionArgs to types.Transaction in Klaytn.
func (args *EthTransactionArgs) ToMessage(globalGasCap uint64, baseFee *big.Int, intrinsicGas uint64) (*types.Transaction, error) {
	// Reject invalid combinations of pre- and post-1559 fee styles
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// Set sender address or use zero address if none specified.
	addr := args.from()

	// Set default gas & gas price if none were set
	gas := globalGasCap
	if gas == 0 {
		gas = uint64(math.MaxUint64 / 2)
	}
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != 0 && globalGasCap < gas {
		logger.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap
	}
	var (
		gasPrice  *big.Int
		gasFeeCap *big.Int
		gasTipCap *big.Int
	)
	if baseFee == nil {
		// If there's no basefee, then it must be a non-1559 execution
		gasPrice = new(big.Int)
		if args.GasPrice != nil {
			gasPrice = args.GasPrice.ToInt()
		}
		gasFeeCap, gasTipCap = gasPrice, gasPrice
	} else {
		// A basefee is provided, necessitating 1559-type execution
		if args.GasPrice != nil {
			// User specified the legacy gas field, convert to 1559 gas typing
			gasPrice = args.GasPrice.ToInt()
			gasFeeCap, gasTipCap = gasPrice, gasPrice
		} else {
			// User specified 1559 gas fields (or none), use those
			gasFeeCap = new(big.Int)
			if args.MaxFeePerGas != nil {
				gasFeeCap = args.MaxFeePerGas.ToInt()
			}
			gasTipCap = new(big.Int)
			if args.MaxPriorityFeePerGas != nil {
				gasTipCap = args.MaxPriorityFeePerGas.ToInt()
			}
			// Backfill the legacy gasPrice for EVM execution, unless we're all zeros
			gasPrice = new(big.Int)
			if gasFeeCap.BitLen() > 0 || gasTipCap.BitLen() > 0 {
				gasPrice = math.BigMin(new(big.Int).Add(gasTipCap, baseFee), gasFeeCap)
			}
		}
	}
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}
	data := args.data()

	// TODO-Klaytn: Klaytn does not support accessList yet.
	// var accessList types.AccessList
	// if args.AccessList != nil {
	//	 accessList = *args.AccessList
	// }
	return types.NewMessage(addr, args.To, 0, value, gas, gasPrice, data, false, intrinsicGas), nil
}

// toTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *EthTransactionArgs) toTransaction() *types.Transaction {
	var tx *types.Transaction
	switch {
	case args.MaxFeePerGas != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		tx = types.NewTx(&types.TxInternalDataEthereumDynamicFee{
			ChainID:      (*big.Int)(args.ChainID),
			AccountNonce: uint64(*args.Nonce),
			GasTipCap:    (*big.Int)(args.MaxPriorityFeePerGas),
			GasFeeCap:    (*big.Int)(args.MaxFeePerGas),
			GasLimit:     uint64(*args.Gas),
			Recipient:    args.To,
			Amount:       (*big.Int)(args.Value),
			Payload:      args.data(),
			AccessList:   al,
		})
	case args.AccessList != nil:
		tx = types.NewTx(&types.TxInternalDataEthereumAccessList{
			ChainID:      (*big.Int)(args.ChainID),
			AccountNonce: uint64(*args.Nonce),
			Recipient:    args.To,
			GasLimit:     uint64(*args.Gas),
			Price:        (*big.Int)(args.GasPrice),
			Amount:       (*big.Int)(args.Value),
			Payload:      args.data(),
			AccessList:   *args.AccessList,
		})
	default:
		tx = types.NewTx(&types.TxInternalDataLegacy{
			AccountNonce: uint64(*args.Nonce),
			Price:        (*big.Int)(args.GasPrice),
			GasLimit:     uint64(*args.Gas),
			Recipient:    args.To,
			Amount:       (*big.Int)(args.Value),
			Payload:      args.data(),
		})
	}
	return tx
}

func (args *EthTransactionArgs) ToTransaction() *types.Transaction {
	return args.toTransaction()
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (api *EthereumAPI) SendTransaction(ctx context.Context, args EthTransactionArgs) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the addresses mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		api.publicTransactionPoolAPI.nonceLock.LockAddr(args.from())
		defer api.publicTransactionPoolAPI.nonceLock.UnlockAddr(args.from())
	}
	if err := args.setDefaults(ctx, api.publicTransactionPoolAPI.b); err != nil {
		return common.Hash{}, err
	}
	tx := args.toTransaction()
	signedTx, err := api.publicTransactionPoolAPI.sign(args.from(), tx)
	if err != nil {
		return common.Hash{}, err
	}
	// Check if signedTx is RLP-Encodable format.
	_, err = rlp.EncodeToBytes(signedTx)
	if err != nil {
		return common.Hash{}, err
	}
	return submitTransaction(ctx, api.publicTransactionPoolAPI.b, signedTx)
}

// EthSignTransactionResult represents a RLP encoded signed transaction.
// SignTransactionResult in go-ethereum has been renamed to EthSignTransactionResult.
// SignTransactionResult is defined in go-ethereum's internal package, so SignTransactionResult is redefined here as EthSignTransactionResult.
type EthSignTransactionResult struct {
	Raw hexutil.Bytes `json:"raw"`
	Tx  *ethTxJSON    `json:"tx"`
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (api *EthereumAPI) FillTransaction(ctx context.Context, args EthTransactionArgs) (*EthSignTransactionResult, error) { // Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, api.publicTransactionPoolAPI.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and obtain rlp
	tx := args.toTransaction()
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	if tx.IsEthTypedTransaction() {
		// Return rawTx except types.TxTypeEthEnvelope: 0x78(= 1 byte)
		return &EthSignTransactionResult{data[1:], formatTxToEthTxJSON(tx)}, nil
	}
	return &EthSignTransactionResult{data, formatTxToEthTxJSON(tx)}, nil
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (api *EthereumAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	if 0 < input[0] && input[0] < 0x7f {
		inputBytes := []byte{byte(types.EthereumTxTypeEnvelope)}
		inputBytes = append(inputBytes, input...)
		return api.publicTransactionPoolAPI.SendRawTransaction(ctx, inputBytes)
	}
	// legacy transaction
	return api.publicTransactionPoolAPI.SendRawTransaction(ctx, input)
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
func (api *EthereumAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	return api.publicTransactionPoolAPI.Sign(addr, data)
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (api *EthereumAPI) SignTransaction(ctx context.Context, args EthTransactionArgs) (*EthSignTransactionResult, error) {
	if args.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}
	if args.GasPrice == nil && (args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil) {
		return nil, fmt.Errorf("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
	}
	if args.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}
	if err := args.setDefaults(ctx, api.publicTransactionPoolAPI.b); err != nil {
		return nil, err
	}
	// Before actually sign the transaction, ensure the transaction fee is reasonable.
	tx := args.toTransaction()
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), api.publicTransactionPoolAPI.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	signed, err := api.publicTransactionPoolAPI.sign(args.from(), tx)
	if err != nil {
		return nil, err
	}
	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if tx.IsEthTypedTransaction() {
		// Return rawTx except types.TxTypeEthEnvelope: 0x78(= 1 byte)
		return &EthSignTransactionResult{data[1:], formatTxToEthTxJSON(tx)}, nil
	}
	return &EthSignTransactionResult{data, formatTxToEthTxJSON(tx)}, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (api *EthereumAPI) PendingTransactions() ([]*EthRPCTransaction, error) {
	pending, err := api.publicTransactionPoolAPI.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	accounts := getAccountsFromWallets(api.publicTransactionPoolAPI.b.AccountManager().Wallets())
	transactions := make([]*EthRPCTransaction, 0, len(pending))
	for _, tx := range pending {
		from := getFrom(tx)
		if _, exists := accounts[from]; exists {
			ethTx := newEthRPCPendingTransaction(tx)
			if ethTx == nil {
				return nil, nil
			}
			transactions = append(transactions, ethTx)
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (api *EthereumAPI) Resend(ctx context.Context, sendArgs EthTransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	return common.Hash{}, errors.New("this api is not supported by Klaytn because Klaytn use fixed gasPrice policy")
}

// Accounts returns the collection of accounts this node manages.
func (api *EthereumAPI) Accounts() []common.Address {
	return api.publicAccountAPI.Accounts()
}

// rpcMarshalHeader marshal block header as Ethereum compatible format.
// It returns error when fetching Author which is block proposer is failed.
func (api *EthereumAPI) rpcMarshalHeader(head *types.Header) (map[string]interface{}, error) {
	var proposer common.Address
	var err error
	if head.Number.Sign() != 0 {
		proposer, err = api.publicKlayAPI.b.Engine().Author(head)
		if err != nil {
			// miner is the field Klaytn should provide the correct value. It's not the field dummy value is allowed.
			logger.Error("Failed to fetch author during marshaling header", "err", err.Error())
			return nil, err
		}
	}
	result := map[string]interface{}{
		"number":          (*hexutil.Big)(head.Number),
		"hash":            head.Hash(),
		"parentHash":      head.ParentHash,
		"nonce":           BlockNonce{},  // There is no block nonce concept in Klaytn, so it must be empty.
		"mixHash":         common.Hash{}, // Klaytn does not use mixHash, so it must be empty.
		"sha3Uncles":      common.HexToHash(EmptySha3Uncles),
		"logsBloom":       head.Bloom,
		"stateRoot":       head.Root,
		"miner":           proposer,
		"difficulty":      (*hexutil.Big)(head.BlockScore),
		"totalDifficulty": (*hexutil.Big)(api.publicKlayAPI.b.GetTd(head.Hash())),
		// extraData always return empty Bytes because actual value of extraData in Klaytn header cannot be used as meaningful way because
		// we cannot provide original header of Klaytn and this field is used as consensus info which is encoded value of validators addresses, validators signatures, and proposer signature in Klaytn.
		"extraData": hexutil.Bytes{},
		"size":      hexutil.Uint64(head.Size()),
		// There is no gas limit mechanism in Klaytn, check details in https://docs.klaytn.com/klaytn/design/computation/computation-cost.
		"gasLimit":         hexutil.Uint64(params.UpperGasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Big(*head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}

	if api.publicBlockChainAPI.b.ChainConfig().IsEthTxTypeForkEnabled(head.Number) {
		result["baseFeePerGas"] = (*hexutil.Big)(new(big.Int).SetUint64(params.BaseFee))
	}

	return result, nil
}

// rpcMarshalBlock marshal block as Ethereum compatible format
func (api *EthereumAPI) rpcMarshalBlock(block *types.Block, inclTx, fullTx bool) (map[string]interface{}, error) {
	fields, err := api.rpcMarshalHeader(block.Header())
	if err != nil {
		return nil, err
	}
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(tx *types.Transaction) (interface{}, error) {
			return tx.Hash(), nil
		}
		if fullTx {
			formatTx = func(tx *types.Transaction) (interface{}, error) {
				return newEthRPCTransactionFromBlockHash(block, tx.Hash()), nil
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range txs {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}
	// There is no uncles in Klaytn
	fields["uncles"] = []common.Hash{}

	return fields, nil
}

func EthDoCall(ctx context.Context, b Backend, args EthTransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *EthStateOverride, timeout time.Duration, globalGasCap uint64) ([]byte, uint64, uint, error) {
	defer func(start time.Time) { logger.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	st, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if st == nil || err != nil {
		return nil, 0, 0, err
	}
	if err := overrides.Apply(st); err != nil {
		return nil, 0, 0, err
	}

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

	// TODO-Klaytn: Klaytn is using fixed baseFee as now but, if we change this fixed baseFee as dynamic baseFee, we should update this logic too.
	fixedBaseFee := new(big.Int).SetUint64(params.BaseFee)
	intrinsicGas, err := types.IntrinsicGas(args.data(), nil, args.To == nil, b.ChainConfig().Rules(header.Number))
	if err != nil {
		return nil, 0, 0, err
	}
	msg, err := args.ToMessage(globalGasCap, fixedBaseFee, intrinsicGas)
	if err != nil {
		return nil, 0, 0, err
	}
	// The intrinsicGas is checked again later in the blockchain.ApplyMessage function,
	// but we check in advance here in order to keep StateTransition.TransactionDb method as unchanged as possible
	// and to clarify error reason correctly to serve eth namespace APIs.
	// This case is handled by EthDoEstimateGas function.
	if msg.Gas() < intrinsicGas {
		return nil, 0, 0, fmt.Errorf("%w: msg.gas %d, want %d", blockchain.ErrIntrinsicGas, msg.Gas(), intrinsicGas)
	}
	evm, vmError, err := b.GetEVM(ctx, msg, st, header, vm.Config{})
	if err != nil {
		return nil, 0, 0, err
	}
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel(vm.CancelByCtxDone)
	}()

	// Execute the message.
	res, gas, kerr := blockchain.ApplyMessage(evm, msg)
	err = kerr.ErrTxInvalid
	if err := vmError(); err != nil {
		return nil, 0, 0, err
	}
	// If the timer caused an abort, return an appropriate error message
	if evm.Cancelled() {
		return nil, 0, 0, fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}
	if err != nil {
		return res, 0, 0, fmt.Errorf("err: %w (supplied gas %d)", err, msg.Gas())
	}
	// TODO-Klaytn-Interface: Introduce ExecutionResult struct from geth to return more detail information
	return res, gas, kerr.Status, nil
}

func EthDoEstimateGas(ctx context.Context, b Backend, args EthTransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, gasCap uint64) (hexutil.Uint64, error) {
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	// Use zero address if sender unspecified.
	if args.From == nil {
		args.From = new(common.Address)
	}
	// Determine the highest gas limit can be used during the estimation.
	if args.Gas != nil && uint64(*args.Gas) >= params.TxGas {
		hi = uint64(*args.Gas)
	} else {
		// Ethereum set hi as gas ceiling of the block but,
		// there is no actual gas limit in Klaytn, so we set it as params.UpperGasLimit.
		hi = params.UpperGasLimit
	}
	// Normalize the max fee per gas the call is willing to spend.
	var feeCap *big.Int
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return 0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		feeCap = args.MaxFeePerGas.ToInt()
	} else {
		feeCap = common.Big0
	}
	// recap the highest gas limit with account's available balance.
	if feeCap.BitLen() != 0 {
		state, _, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return 0, err
		}
		balance := state.GetBalance(*args.From) // from can't be nil
		available := new(big.Int).Set(balance)
		if args.Value != nil {
			if args.Value.ToInt().Cmp(available) >= 0 {
				return 0, errors.New("insufficient funds for transfer")
			}
			available.Sub(available, args.Value.ToInt())
		}
		allowance := new(big.Int).Div(available, feeCap)

		// If the allowance is larger than maximum uint64, skip checking
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := args.Value
			if transfer == nil {
				transfer = new(hexutil.Big)
			}
			logger.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", transfer.ToInt(), "maxFeePerGas", feeCap, "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	// Recap the highest gas allowance with specified gascap.
	if gasCap != 0 && hi > gasCap {
		logger.Warn("Caller gas above allowance, capping", "requested", hi, "cap", gasCap)
		hi = gasCap
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction.
	// executable returns
	// - bool: true when a call with given args and gas is executable.
	// - []byte: EVM execution result.
	// - error: error occurred during EVM execution.
	// - error: consensus error which is not EVM related error (less balance of caller, wrong nonce, etc...).
	executable := func(gas uint64) (bool, []byte, error, error) {
		args.Gas = (*hexutil.Uint64)(&gas)
		ret, _, status, err := EthDoCall(ctx, b, args, rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber), nil, 0, gasCap)
		if err != nil {
			if errors.Is(err, blockchain.ErrIntrinsicGas) {
				// Special case, raise gas limit
				return false, ret, nil, nil
			}
			// Returns error when it is not VM error (less balance or wrong nonce, etc...).
			return false, nil, nil, err
		}
		// If err is vmError, return vmError with returned data
		vmErr := blockchain.GetVMerrFromReceiptStatus(status)
		if vmErr != nil {
			return false, ret, vmErr, nil
		}
		return true, ret, vmErr, nil
	}

	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		isExecutable, _, _, err := executable(mid)
		if err != nil {
			return 0, err
		}

		if !isExecutable {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		isExecutable, ret, vmErr, err := executable(hi)
		if err != nil {
			return 0, err
		}
		if !isExecutable {
			if vmErr != nil {
				// Treat vmErr as RevertError only when there was returned data from call.
				if isReverted(vmErr) && len(ret) > 0 {
					return 0, newRevertError(ret)
				}
				return 0, vmErr
			}
			// Otherwise, the specified gas cap is too low
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}
	return hexutil.Uint64(hi), nil
}

// isReverted checks given error is vm.ErrExecutionReverted
func isReverted(err error) bool {
	if errors.Is(err, vm.ErrExecutionReverted) {
		return true
	}
	return false
}

// newRevertError wraps data returned when EVM execution was reverted.
// Make sure that data is returned when execution reverted situation.
func newRevertError(data []byte) *revertError {
	reason, errUnpack := abi.UnpackRevert(data)
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(data),
	}
}

// revertError is an API error that encompassas an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func checkTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	feeEth := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))), new(big.Float).SetInt(big.NewInt(params.KLAY)))
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f klay) exceeds the configured cap (%.2f klay)", feeFloat, cap)
	}
	return nil
}
