// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from eth/api.go (2018/06/04).
// Modified and improved for the klaytn development.

package cn

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/klaytn/klaytn/work"
	"io"
	"math/big"
	"os"
	"strings"
)

// PublicKlayAPI provides an API to access Klaytn CN-related
// information.
type PublicKlayServiceChainAPI struct {
	sc *ServiceChain
}

// NewPublicKlayAPI creates a new Klaytn protocol API for full nodes.
func NewPublicKlayServiceChainAPI(sc *ServiceChain) *PublicKlayServiceChainAPI {
	return &PublicKlayServiceChainAPI{sc}
}

// Signer is the address that signing the block.
func (api *PublicKlayServiceChainAPI) Signer() (common.Address, error) {
	return api.sc.Signer()
}

// Hashrate returns the POW hashrate
func (api *PublicKlayServiceChainAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(api.sc.Miner().HashRate())
}

// PublicMinerAPI provides an API to control the miner.
// It offers only methods that operate on data that pose no security risk when it is publicly accessible.
type PublicServiceChainMinerAPI struct {
	sc    *ServiceChain
	agent *work.RemoteAgent
}

// NewPublicMinerAPI create a new PublicMinerAPI instance.
func NewPublicServiceChainMinerAPI(e *ServiceChain) *PublicServiceChainMinerAPI {
	agent := work.NewRemoteAgent(e.BlockChain(), e.Engine())
	e.Miner().Register(agent)

	return &PublicServiceChainMinerAPI{e, agent}
}

// Mining returns an indication if this node is currently mining.
func (api *PublicServiceChainMinerAPI) Mining() bool {
	return api.sc.IsMining()
}

// SubmitWork can be used by external miner to submit their POW solution. It returns an indication if the work was
// accepted. Note, this is not an indication if the provided work was valid!
func (api *PublicServiceChainMinerAPI) SubmitWork(solution common.Hash) bool {
	return api.agent.SubmitWork(solution)
}

// GetWork returns a work package for external miner. The work package consists of 3 strings
// result[0], 32 bytes hex encoded current block header pow-hash
// result[1], 32 bytes hex encoded seed hash used for DAG
// result[2], 32 bytes hex encoded boundary condition ("target"), 2^256/blockScore
func (api *PublicServiceChainMinerAPI) GetWork() ([3]string, error) {
	if !api.sc.IsMining() {
		if err := api.sc.StartMining(false); err != nil {
			return [3]string{}, err
		}
	}
	work, err := api.agent.GetWork()
	if err != nil {
		return work, fmt.Errorf("mining not ready: %v", err)
	}
	return work, nil
}

// PrivateMinerAPI provides private RPC methods to control the miner.
// These methods can be abused by external users and must be considered insecure for use by untrusted users.
type PrivateServiceChainMinerAPI struct {
	e *ServiceChain
}

// NewPrivateMinerAPI create a new RPC service which controls the miner of this node.
func NewPrivateServiceChainMinerAPI(e *ServiceChain) *PrivateServiceChainMinerAPI {
	return &PrivateServiceChainMinerAPI{e: e}
}

// Start the miner with the given number of threads. If threads is nil the number
// of workers started is equal to the number of logical CPUs that are usable by
// this process. If mining is already running, this method adjust the number of
// threads allowed to use.
func (api *PrivateServiceChainMinerAPI) Start(threads *int) error {
	// Set the number of threads if the seal engine supports it
	if threads == nil {
		threads = new(int)
	} else if *threads == 0 {
		*threads = -1 // Disable the work from within
	}
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := api.e.engine.(threaded); ok {
		logger.Info("Updated mining threads", "threads", *threads)
		th.SetThreads(*threads)
	}
	// Start the miner and return
	if !api.e.IsMining() {
		// Propagate the initial price point to the transaction pool
		api.e.lock.RLock()
		price := api.e.gasPrice
		api.e.lock.RUnlock()

		if price.Cmp(api.e.txPool.GasPrice()) == 0 {
			return api.e.StartMining(true)
		} else {
			logger.Error("PrivateMinerAPI Start: Invalid unit price from API", "TxPool UnitPrice", api.e.txPool.GasPrice(), "API UnitPrice", price)
			return blockchain.ErrInvalidUnitPrice
		}
	}
	return nil
}

// Stop the miner
func (api *PrivateServiceChainMinerAPI) Stop() bool {
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := api.e.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	api.e.StopMining()
	return true
}

// SetExtra sets the extra data string that is included when this miner mines a block.
func (api *PrivateServiceChainMinerAPI) SetExtra(extra string) (bool, error) {
	if err := api.e.Miner().SetExtra([]byte(extra)); err != nil {
		return false, err
	}
	return true, nil
}

// SetGasPrice sets the minimum accepted gas price for the miner.
func (api *PrivateServiceChainMinerAPI) SetGasPrice(gasPrice hexutil.Big) bool {
	if api.e.txPool.GasPrice().Cmp((*big.Int)(&gasPrice)) != 0 {
		logger.Debug("PrivateMinerAPI.SetGasPrice", "TxPool UnitPrice", api.e.txPool.GasPrice(), "Given UnitPrice", gasPrice.ToInt())
		return false
	}

	api.e.lock.Lock()
	api.e.gasPrice = (*big.Int)(&gasPrice)
	api.e.lock.Unlock()

	api.e.txPool.SetGasPrice((*big.Int)(&gasPrice))
	return true
}

// GetHashrate returns the current hashrate of the miner.
func (api *PrivateServiceChainMinerAPI) GetHashrate() uint64 {
	return uint64(api.e.miner.HashRate())
}

// PrivateAdminAPI is the collection of sc full node-related APIs
// exposed over the private admin endpoint.
type PrivateServiceChainAdminAPI struct {
	sc *ServiceChain
}

// NewPrivateAdminAPI creates a new API definition for the full node private
// admin methods of the sc service.
func NewPrivateServiceChainAdminAPI(sc *ServiceChain) *PrivateServiceChainAdminAPI {
	return &PrivateServiceChainAdminAPI{sc: sc}
}

// ExportChain exports the current blockchain into a local file.
func (api *PrivateServiceChainAdminAPI) ExportChain(file string) (bool, error) {
	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if err := api.sc.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}

// ImportChain imports a blockchain from a local file.
func (api *PrivateServiceChainAdminAPI) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}

	// Run actual the import in pre-configured batches
	stream := rlp.NewStream(reader, 0)

	blocks, index := make([]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.sc.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.sc.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}

// PublicDebugAPI is the collection of Klaytn full node APIs exposed
// over the public debugging endpoint.
type PublicServiceChainDebugAPI struct {
	sc *ServiceChain
}

// NewPublicDebugAPI creates a new API definition for the full node-
// related public debug methods of the Klaytn service.
func NewPublicServiceChainDebugAPI(sc *ServiceChain) *PublicServiceChainDebugAPI {
	return &PublicServiceChainDebugAPI{sc: sc}
}

// DumpBlock retrieves the entire state of the database at a given block.
func (api *PublicServiceChainDebugAPI) DumpBlock(blockNr rpc.BlockNumber) (state.Dump, error) {
	if blockNr == rpc.PendingBlockNumber {
		// If we're dumping the pending state, we need to request
		// both the pending block as well as the pending state from
		// the miner and operate on those
		_, stateDb := api.sc.miner.Pending()
		return stateDb.RawDump(), nil
	}
	var block *types.Block
	if blockNr == rpc.LatestBlockNumber {
		block = api.sc.blockchain.CurrentBlock()
	} else {
		block = api.sc.blockchain.GetBlockByNumber(uint64(blockNr))
	}
	if block == nil {
		return state.Dump{}, fmt.Errorf("block #%d not found", blockNr)
	}
	stateDb, err := api.sc.BlockChain().StateAt(block.Root())
	if err != nil {
		return state.Dump{}, err
	}
	return stateDb.RawDump(), nil
}

// PrivateDebugAPI is the collection of sc full node APIs exposed over
// the private debugging endpoint.
type PrivateServiceChainDebugAPI struct {
	config *params.ChainConfig
	sc     *ServiceChain
}

// NewPrivateDebugAPI creates a new API definition for the full node-related
// private debug methods of the sc service.
func NewPrivateServiceChainDebugAPI(config *params.ChainConfig, sc *ServiceChain) *PrivateServiceChainDebugAPI {
	return &PrivateServiceChainDebugAPI{config: config, sc: sc}
}

// Preimage is a debug API function that returns the preimage for a sha3 hash, if known.
func (api *PrivateServiceChainDebugAPI) Preimage(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	if preimage := api.sc.ChainDB().ReadPreimage(hash); preimage != nil {
		return preimage, nil
	}
	return nil, errors.New("unknown preimage")
}

// GetBadBLocks returns a list of the last 'bad blocks' that the client has seen on the network
// and returns them as a JSON list of block-hashes
func (api *PrivateServiceChainDebugAPI) GetBadBlocks(ctx context.Context) ([]blockchain.BadBlockArgs, error) {
	return api.sc.BlockChain().BadBlocks()
}

// GetModifiedAccountsByumber returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateServiceChainDebugAPI) GetModifiedAccountsByNumber(startNum uint64, endNum *uint64) ([]common.Address, error) {
	var startBlock, endBlock *types.Block

	startBlock = api.sc.blockchain.GetBlockByNumber(startNum)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startNum)
	}

	if endNum == nil {
		endBlock = startBlock
		startBlock = api.sc.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.sc.blockchain.GetBlockByNumber(*endNum)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %d not found", *endNum)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

// GetModifiedAccountsByHash returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateServiceChainDebugAPI) GetModifiedAccountsByHash(startHash common.Hash, endHash *common.Hash) ([]common.Address, error) {
	var startBlock, endBlock *types.Block
	startBlock = api.sc.blockchain.GetBlockByHash(startHash)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startHash)
	}

	if endHash == nil {
		endBlock = startBlock
		startBlock = api.sc.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.sc.blockchain.GetBlockByHash(*endHash)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %x not found", *endHash)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

func (api *PrivateServiceChainDebugAPI) getModifiedAccounts(startBlock, endBlock *types.Block) ([]common.Address, error) {
	if startBlock.Number().Uint64() >= endBlock.Number().Uint64() {
		return nil, fmt.Errorf("start block height (%d) must be less than end block height (%d)", startBlock.Number().Uint64(), endBlock.Number().Uint64())
	}

	oldTrie, err := statedb.NewSecureTrie(startBlock.Root(), statedb.NewDatabase(api.sc.chainDB))
	if err != nil {
		return nil, err
	}
	newTrie, err := statedb.NewSecureTrie(endBlock.Root(), statedb.NewDatabase(api.sc.chainDB))
	if err != nil {
		return nil, err
	}

	diff, _ := statedb.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := statedb.NewIterator(diff)

	var dirty []common.Address
	for iter.Next() {
		key := newTrie.GetKey(iter.Key)
		if key == nil {
			return nil, fmt.Errorf("no preimage found for hash %x", iter.Key)
		}
		dirty = append(dirty, common.BytesToAddress(key))
	}
	return dirty, nil
}
