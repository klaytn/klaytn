// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/backend/api.go (2018/06/04).
// Modified and improved for the klaytn development.

package backend

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	klaytnApi "github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/networks/rpc"
)

// API is a user facing RPC API to dump Istanbul state
type API struct {
	chain    consensus.ChainReader
	istanbul *backend
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.istanbul.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.istanbul.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetValidators retrieves the list of authorized validators with the given block number.
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	snap, err := api.GetSnapshot(number)
	if err != nil {
		return nil, err
	}
	return snap.validators(), nil
}

// GetValidatorsAtHash retrieves the list of authorized validators with the given block hash.
func (api *API) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	snap, err := api.GetSnapshotAtHash(hash)
	if err != nil {
		return nil, err
	}
	return snap.validators(), nil
}

// GetDemotedValidators retrieves the list of authorized, but demoted validators with the given block number.
func (api *API) GetDemotedValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	snap, err := api.GetSnapshot(number)
	if err != nil {
		return nil, err
	}
	return snap.demotedValidators(), nil
}

// GetDemotedValidatorsAtHash retrieves the list of authorized, but demoted validators with the given block hash.
func (api *API) GetDemotedValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	snap, err := api.GetSnapshotAtHash(hash)
	if err != nil {
		return nil, err
	}
	return snap.demotedValidators(), nil
}

// Candidates returns the current candidates the node tries to uphold and vote on.
func (api *API) Candidates() map[common.Address]bool {
	api.istanbul.candidatesLock.RLock()
	defer api.istanbul.candidatesLock.RUnlock()

	proposals := make(map[common.Address]bool)
	for address, auth := range api.istanbul.candidates {
		proposals[address] = auth
	}
	return proposals
}

// Propose injects a new authorization candidate that the validator will attempt to
// push through.
func (api *API) Propose(address common.Address, auth bool) {
	api.istanbul.candidatesLock.Lock()
	defer api.istanbul.candidatesLock.Unlock()

	api.istanbul.candidates[address] = auth
}

// Discard drops a currently running candidate, stopping the validator from casting
// further votes (either for or against).
func (api *API) Discard(address common.Address) {
	api.istanbul.candidatesLock.Lock()
	defer api.istanbul.candidatesLock.Unlock()

	delete(api.istanbul.candidates, address)
}

// API extended by Klaytn developers
type APIExtension struct {
	chain    consensus.ChainReader
	istanbul *backend
}

var (
	errPendingNotAllowed       = errors.New("pending is not allowed")
	errInternalError           = errors.New("internal error")
	errStartNotPositive        = errors.New("start block number should be positive")
	errEndLargetThanLatest     = errors.New("end block number should be smaller than the latest block number")
	errStartLargerThanEnd      = errors.New("start should be smaller than end")
	errRequestedBlocksTooLarge = errors.New("number of requested blocks should be smaller than 50")
	errRangeNil                = errors.New("range values should not be nil")
	errExtractIstanbulExtra    = errors.New("extract Istanbul Extra from block header of the given block number")
	errNoBlockExist            = errors.New("block with the given block number is not existed")
	errNoBlockNumber           = errors.New("block number is not assigned")
)

// GetCouncil retrieves the list of authorized validators at the specified block.
func (api *APIExtension) GetCouncil(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else if *number == rpc.PendingBlockNumber {
		logger.Trace("Cannot get council of the pending block.", "number", number)
		return nil, errPendingNotAllowed
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the council from its snapshot
	if header == nil {
		return nil, errNoBlockExist // return nil if block is not found.
	}

	// Since v1.6.1, extra.validators represents a list of council
	// TODO-Klaytn : replace this to the below calculation logic
	//istanbulExtra, err := types.ExtractIstanbulExtra(header)
	//if err == nil {
	//	return istanbulExtra.Validators, nil
	//} else {
	//	return nil, errExtractIstanbulExtra
	//}

	// Calculate council list from snapshot
	snap, err := api.istanbul.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		logger.Error("Failed to get snapshot.", "hash", header.Hash(), "err", err)
		return nil, errInternalError
	}

	return append(snap.validators(), snap.demotedValidators()...), nil
}

func (api *APIExtension) GetCouncilSize(number *rpc.BlockNumber) (int, error) {
	council, err := api.GetCouncil(number)
	if err == nil {
		return len(council), nil
	} else {
		return -1, err
	}
}

func (api *APIExtension) GetCommittee(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else if *number == rpc.PendingBlockNumber {
		logger.Trace("Cannot get validators of the pending block.", "number", number)
		return nil, errPendingNotAllowed
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}

	if header == nil {
		return nil, errNoBlockExist
	}

	blockNumber := header.Number.Uint64()
	if blockNumber == 0 {
		// The committee of genesis block can not be calculated because it requires a previous block.
		istanbulExtra, err := types.ExtractIstanbulExtra(header)
		if err != nil {
			return nil, errExtractIstanbulExtra
		}
		return istanbulExtra.Validators, nil
	}

	round := header.Round()
	view := &istanbul.View{
		Sequence: new(big.Int).SetUint64(blockNumber),
		Round:    new(big.Int).SetUint64(uint64(round)),
	}

	// get the proposer of this block.
	proposer, err := ecrecover(header)
	if err != nil {
		return nil, err
	}

	// get the snapshot of the previous block.
	parentHash := header.ParentHash
	snap, err := api.istanbul.snapshot(api.chain, blockNumber-1, parentHash, nil)
	if err != nil {
		return nil, err
	}

	// get the committee list of this block at the view (blockNumber, round)
	committee := snap.ValSet.SubListWithProposer(parentHash, proposer, view, api.chain.Config().IsIstanbul(view.Sequence))
	addresses := make([]common.Address, len(committee))
	for i, v := range committee {
		addresses[i] = v.Address()
	}

	return addresses, nil
}

func (api *APIExtension) GetCommitteeSize(number *rpc.BlockNumber) (int, error) {
	committee, err := api.GetCommittee(number)
	if err == nil {
		return len(committee), nil
	} else {
		return -1, err
	}
}

type ConsensusInfo struct {
	proposer       common.Address
	originProposer common.Address // the proposal of 0 round at the same block number
	committee      []common.Address
	round          byte
}

func (api *APIExtension) getConsensusInfo(block *types.Block) (ConsensusInfo, error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return ConsensusInfo{}, nil
	}

	round := block.Header().Round()
	view := &istanbul.View{
		Sequence: new(big.Int).Set(block.Number()),
		Round:    new(big.Int).SetInt64(int64(round)),
	}

	// get the proposer of this block.
	proposer, err := ecrecover(block.Header())
	if err != nil {
		return ConsensusInfo{}, err
	}

	// get the snapshot of the previous block.
	parentHash := block.ParentHash()
	snap, err := api.istanbul.snapshot(api.chain, blockNumber-1, parentHash, nil)
	if err != nil {
		return ConsensusInfo{}, err
	}

	// get origin proposer at 0 round.
	originProposer := common.Address{}
	lastProposer := api.istanbul.GetProposer(blockNumber - 1)

	newValSet := snap.ValSet.Copy()
	newValSet.CalcProposer(lastProposer, 0)
	originProposer = newValSet.GetProposer().Address()

	// get the committee list of this block at the view (blockNumber, round)
	committee := snap.ValSet.SubListWithProposer(parentHash, proposer, view, api.chain.Config().IsIstanbul(view.Sequence))
	committeeAddrs := make([]common.Address, len(committee))
	for i, v := range committee {
		committeeAddrs[i] = v.Address()
	}

	// verify the committee list of the block using istanbul
	//proposalSeal := istanbulCore.PrepareCommittedSeal(block.Hash())
	//extra, err := types.ExtractIstanbulExtra(block.Header())
	//istanbulAddrs := make([]common.Address, len(committeeAddrs))
	//for i, seal := range extra.CommittedSeal {
	//	addr, err := istanbul.GetSignatureAddress(proposalSeal, seal)
	//	istanbulAddrs[i] = addr
	//	if err != nil {
	//		return proposer, []common.Address{}, err
	//	}
	//
	//	var found bool = false
	//	for _, v := range committeeAddrs {
	//		if addr == v {
	//			found = true
	//			break
	//		}
	//	}
	//	if found == false {
	//		logger.Trace("validator is different!", "snap", committeeAddrs, "istanbul", istanbulAddrs)
	//		return proposer, committeeAddrs, errors.New("validator set is different from Istanbul engine!!")
	//	}
	//}

	cInfo := ConsensusInfo{
		proposer:       proposer,
		originProposer: originProposer,
		committee:      committeeAddrs,
		round:          round,
	}

	return cInfo, nil
}

func (api *APIExtension) makeRPCBlockOutput(b *types.Block,
	cInfo ConsensusInfo, transactions types.Transactions, receipts types.Receipts) map[string]interface{} {
	head := b.Header() // copies the header once
	hash := head.Hash()

	td := big.NewInt(0)
	if bc, ok := api.chain.(*blockchain.BlockChain); ok {
		td = bc.GetTd(hash, b.NumberU64())
	}
	r, err := klaytnApi.RpcOutputBlock(b, td, false, false)
	if err != nil {
		logger.Error("failed to RpcOutputBlock", "err", err)
		return nil
	}

	// make transactions
	numTxs := len(transactions)
	rpcTransactions := make([]map[string]interface{}, numTxs)
	for i, tx := range transactions {
		rpcTransactions[i] = klaytnApi.RpcOutputReceipt(tx, hash, head.Number.Uint64(), uint64(i), receipts[i])
	}

	r["committee"] = cInfo.committee
	r["proposer"] = cInfo.proposer
	r["round"] = cInfo.round
	r["originProposer"] = cInfo.originProposer
	r["transactions"] = rpcTransactions

	return r
}

// TODO-Klaytn: This API functions should be managed with API functions with namespace "klay"
func (api *APIExtension) GetBlockWithConsensusInfoByNumber(number *rpc.BlockNumber) (map[string]interface{}, error) {
	b, ok := api.chain.(*blockchain.BlockChain)
	if !ok {
		logger.Error("chain is not a type of blockchain.BlockChain", "type", reflect.TypeOf(api.chain))
		return nil, errInternalError
	}
	var block *types.Block
	var blockNumber uint64

	if number == nil {
		logger.Trace("block number is not assigned")
		return nil, errNoBlockNumber
	}

	if *number == rpc.PendingBlockNumber {
		logger.Trace("Cannot get consensus information of the PendingBlock.")
		return nil, errPendingNotAllowed
	}

	if *number == rpc.LatestBlockNumber {
		block = b.CurrentBlock()
		blockNumber = block.NumberU64()
	} else {
		// rpc.EarliestBlockNumber == 0, no need to treat it as a special case.
		blockNumber = uint64(number.Int64())
		block = b.GetBlockByNumber(blockNumber)
	}

	if block == nil {
		logger.Trace("Finding a block by number failed.", "blockNum", blockNumber)
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNumber)
	}
	blockHash := block.Hash()

	cInfo, err := api.getConsensusInfo(block)
	if err != nil {
		logger.Error("Getting the proposer and validators failed.", "blockHash", blockHash, "err", err)
		return nil, errInternalError
	}

	receipts := b.GetBlockReceiptsInCache(blockHash)
	if receipts == nil {
		receipts = b.GetReceiptsByBlockHash(blockHash)
	}

	return api.makeRPCBlockOutput(block, cInfo, block.Transactions(), receipts), nil
}

func (api *APIExtension) GetBlockWithConsensusInfoByNumberRange(start *rpc.BlockNumber, end *rpc.BlockNumber) (map[string]interface{}, error) {
	blocks := make(map[string]interface{})

	if start == nil || end == nil {
		logger.Trace("the range values should not be nil.", "start", start, "end", end)
		return nil, errRangeNil
	}

	// check error status.
	s := start.Int64()
	e := end.Int64()
	if s < 0 {
		logger.Trace("start should be positive", "start", s)
		return nil, errStartNotPositive
	}

	eChain := api.chain.CurrentHeader().Number.Int64()
	if e > eChain {
		logger.Trace("end should be smaller than the lastest block number", "end", end, "eChain", eChain)
		return nil, errEndLargetThanLatest
	}

	if s > e {
		logger.Trace("start should be smaller than end", "start", s, "end", e)
		return nil, errStartLargerThanEnd
	}

	if (e - s) > 50 {
		logger.Trace("number of requested blocks should be smaller than 50", "start", s, "end", e)
		return nil, errRequestedBlocksTooLarge
	}

	// gather s~e blocks
	for i := s; i <= e; i++ {
		strIdx := fmt.Sprintf("0x%x", i)

		blockNum := rpc.BlockNumber(i)
		b, err := api.GetBlockWithConsensusInfoByNumber(&blockNum)
		if err != nil {
			logger.Error("error on GetBlockWithConsensusInfoByNumber", "err", err)
			blocks[strIdx] = nil
		} else {
			blocks[strIdx] = b
		}
	}

	return blocks, nil
}

func (api *APIExtension) GetBlockWithConsensusInfoByHash(blockHash common.Hash) (map[string]interface{}, error) {
	b, ok := api.chain.(*blockchain.BlockChain)
	if !ok {
		logger.Error("chain is not a type of blockchain.Blockchain, returning...", "type", reflect.TypeOf(api.chain))
		return nil, errInternalError
	}

	block := b.GetBlockByHash(blockHash)
	if block == nil {
		logger.Trace("Finding a block failed.", "blockHash", blockHash)
		return nil, fmt.Errorf("the block does not exist (block hash: %s)", blockHash.String())
	}

	cInfo, err := api.getConsensusInfo(block)
	if err != nil {
		logger.Error("Getting the proposer and validators failed.", "blockHash", blockHash, "err", err)
		return nil, errInternalError
	}

	receipts := b.GetBlockReceiptsInCache(blockHash)
	if receipts == nil {
		receipts = b.GetReceiptsByBlockHash(blockHash)
	}

	return api.makeRPCBlockOutput(block, cInfo, block.Transactions(), receipts), nil
}

func (api *API) GetTimeout() uint64 {
	return istanbul.DefaultConfig.Timeout
}
