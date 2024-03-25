// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from eth/gasprice/gasprice.go (2018/06/04).
// Modified and improved for the klaytn development.

package gasprice

import (
	"context"
	"math/big"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"golang.org/x/exp/slices"
)

const sampleNumber = 3 // Number of transactions sampled in a block

var maxPrice = big.NewInt(500 * params.Ston)

type Config struct {
	Blocks           int
	Percentile       int
	MaxHeaderHistory int
	MaxBlockHistory  int
	Default          *big.Int `toml:",omitempty"`
}

// OracleBackend includes all necessary background APIs for oracle.
type OracleBackend interface {
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error)
	GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts
	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block
}

type TxPool interface {
	GasPrice() *big.Int
}

// Oracle recommends gas prices based on the content of recent
// blocks. Suitable for both light and full clients.
type Oracle struct {
	backend   OracleBackend
	lastHead  common.Hash
	lastPrice *big.Int
	cacheLock sync.RWMutex
	fetchLock sync.Mutex
	txPool    TxPool

	checkBlocks, maxEmpty, maxBlocks  int
	percentile                        int
	maxHeaderHistory, maxBlockHistory int

	historyCache *lru.Cache
}

// NewOracle returns a new oracle.
func NewOracle(backend OracleBackend, params Config, txPool TxPool) *Oracle {
	blocks := params.Blocks
	if blocks < 1 {
		blocks = 1
	}
	percent := params.Percentile
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	maxHeaderHistory := params.MaxHeaderHistory
	if maxHeaderHistory < 1 {
		maxHeaderHistory = 1
		logger.Warn("Sanitizing invalid gasprice oracle max header history", "provided", params.MaxHeaderHistory, "updated", maxHeaderHistory)
	}
	maxBlockHistory := params.MaxBlockHistory
	if maxBlockHistory < 1 {
		maxBlockHistory = 1
		logger.Warn("Sanitizing invalid gasprice oracle max block history", "provided", params.MaxBlockHistory, "updated", maxBlockHistory)
	}
	cache, _ := lru.New(2048)
	
	return &Oracle{
		backend:          backend,
		lastPrice:        params.Default,
		checkBlocks:      blocks,
		maxEmpty:         blocks / 2,
		maxBlocks:        blocks * 5,
		percentile:       percent,
		maxHeaderHistory: maxHeaderHistory,
		maxBlockHistory:  maxBlockHistory,
		txPool:           txPool,
		historyCache:     cache,
	}
}

// Tx gas price requirements has changed over the hardforks
//
// | Fork              | gasPrice                     | maxFeePerGas                 | maxPriorityFeePerGas         |
// |------------------ |----------------------------- |----------------------------- |----------------------------- |
// | Before EthTxType  | must be fixed UnitPrice (1)  | N/A (2)                      | N/A (2)                      |
// | After EthTxType   | must be fixed UnitPrice (1)  | must be fixed UnitPrice (3)  | must be fixed UnitPrice (3)  |
// | After Magma       | BaseFee or higher (4)        | BaseFee or higher (4)        | Ignored (4)                  |
//
// (1) If tx.type != 2 && !rules.IsMagma: https://github.com/klaytn/klaytn/blob/v1.11.1/blockchain/tx_pool.go#L729
// (2) If tx.type == 2 && !rules.IsEthTxType: https://github.com/klaytn/klaytn/blob/v1.11.1/blockchain/tx_pool.go#L670
// (3) If tx.type == 2 && !rules.IsMagma: https://github.com/klaytn/klaytn/blob/v1.11.1/blockchain/tx_pool.go#L710
// (4) If tx.type == 2 && rules.IsMagma: https://github.com/klaytn/klaytn/blob/v1.11.1/blockchain/tx_pool.go#L703
//
// The suggested prices needs to match the requirements.
//
// | Fork              | SuggestPrice (for gasPrice and maxFeePerGas)                | SuggestTipCap (for maxPriorityFeePerGas) |
// |------------------ |------------------------------------------------------------ |----------------------------- |
// | Before Magma      | Fixed UnitPrice                                             | Fixed UnitPrice              |
// | After Magma       | BaseFee * 2                                                 | Zero                         |

// SuggestPrice returns the recommended gas price.
// This value is intended to be used as gasPrice or maxFeePerGas.
func (gpo *Oracle) SuggestPrice(ctx context.Context) (*big.Int, error) {
	if gpo.txPool == nil {
		// If txpool is not set, just return 0. This is used for testing.
		return common.Big0, nil
	}

	nextNum := new(big.Int).Add(gpo.backend.CurrentBlock().Number(), common.Big1)
	if gpo.backend.ChainConfig().IsDragonForkEnabled(nextNum) {
		// After Dragon, include suggested tip
		baseFee := gpo.txPool.GasPrice()
		suggestedTip, err := gpo.SuggestTipCap(ctx)
		if err != nil {
			return nil, err
		}
		return new(big.Int).Add(new(big.Int).Mul(baseFee, common.Big2), suggestedTip), nil
	} else if gpo.backend.ChainConfig().IsMagmaForkEnabled(nextNum) {
		// After Magma, return the twice of BaseFee as a buffer.
		baseFee := gpo.txPool.GasPrice()
		return new(big.Int).Mul(baseFee, common.Big2), nil
	} else {
		// Before Magma, return the fixed UnitPrice.
		unitPrice := gpo.txPool.GasPrice()
		return unitPrice, nil
	}
}

// SuggestTipCap returns the recommended gas tip cap.
// This value is intended to be used as maxPriorityFeePerGas.
func (gpo *Oracle) SuggestTipCap(ctx context.Context) (*big.Int, error) {
	if gpo.txPool == nil {
		// If txpool is not set, just return 0. This is used for testing.
		return common.Big0, nil
	}

	nextNum := new(big.Int).Add(gpo.backend.CurrentBlock().Number(), common.Big1)
	if gpo.backend.ChainConfig().IsDragonForkEnabled(nextNum) {
		// After Dragon, return using fee history.
		// By default config, this will return 60% percentile of last 20 blocks
		// See node/cn/config.go for the default config.
		return gpo.suggestTipCapUsingFeeHistory(ctx)
	} else if gpo.backend.ChainConfig().IsMagmaForkEnabled(nextNum) {
		// After Magma, return zero
		return common.Big0, nil
	} else {
		// Before Magma, return the fixed UnitPrice.
		unitPrice := gpo.txPool.GasPrice()
		return unitPrice, nil
	}
}

// suggestTipCapUsingFeeHistory returns a tip cap based on fee history.
func (oracle *Oracle) suggestTipCapUsingFeeHistory(ctx context.Context) (*big.Int, error) {
	head, _ := oracle.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
	headHash := head.Hash()

	// If the latest gasprice is still available, return it.
	oracle.cacheLock.RLock()
	lastHead, lastPrice := oracle.lastHead, oracle.lastPrice
	oracle.cacheLock.RUnlock()
	if headHash == lastHead {
		return new(big.Int).Set(lastPrice), nil
	}
	oracle.fetchLock.Lock()
	defer oracle.fetchLock.Unlock()

	// Try checking the cache again, maybe the last fetch fetched what we need
	oracle.cacheLock.RLock()
	lastHead, lastPrice = oracle.lastHead, oracle.lastPrice
	oracle.cacheLock.RUnlock()
	if headHash == lastHead {
		return new(big.Int).Set(lastPrice), nil
	}
	var (
		sent, exp int
		number    = head.Number.Uint64()
		result    = make(chan results, oracle.checkBlocks)
		quit      = make(chan struct{})
		results   []*big.Int
	)
	for sent < oracle.checkBlocks && number > 0 {
		go oracle.getBlockValues(ctx, number, sampleNumber, result, quit)
		sent++
		exp++
		number--
	}
	for exp > 0 {
		res := <-result
		if res.err != nil {
			close(quit)
			return new(big.Int).Set(lastPrice), res.err
		}
		exp--
		// Nothing returned. There are two special cases here:
		// - The block is empty
		// - All the transactions included are sent by the miner itself.
		// In these cases, use the latest calculated price for sampling.
		if len(res.values) == 0 {
			res.values = []*big.Int{lastPrice}
		}
		// Besides, in order to collect enough data for sampling, if nothing
		// meaningful returned, try to query more blocks. But the maximum
		// is 2*checkBlocks.
		if len(res.values) == 1 && len(results)+1+exp < oracle.checkBlocks*2 && number > 0 {
			go oracle.getBlockValues(ctx, number, sampleNumber, result, quit)
			sent++
			exp++
			number--
		}
		results = append(results, res.values...)
	}
	price := lastPrice
	if len(results) > 0 {
		slices.SortFunc(results, func(a, b *big.Int) int { return a.Cmp(b) })
		price = results[(len(results)-1)*oracle.percentile/100]
	}
	if price.Cmp(maxPrice) > 0 {
		price = new(big.Int).Set(maxPrice)
	}
	oracle.cacheLock.Lock()
	oracle.lastHead = headHash
	oracle.lastPrice = price
	oracle.cacheLock.Unlock()

	return new(big.Int).Set(price), nil
}

type results struct {
	values []*big.Int
	err    error
}

// getBlockValues calculates the lowest transaction gas price in a given block
// and sends it to the result channel. If the block is empty or all transactions
// are sent by the miner itself(it doesn't make any sense to include this kind of
// transaction prices for sampling), nil gasprice is returned.
func (oracle *Oracle) getBlockValues(ctx context.Context, blockNum uint64, limit int, result chan results, quit chan struct{}) {
	block, err := oracle.backend.BlockByNumber(ctx, rpc.BlockNumber(blockNum))
	if block == nil {
		select {
		case result <- results{nil, err}:
		case <-quit:
		}
		return
	}
	signer := types.MakeSigner(oracle.backend.ChainConfig(), block.Number())

	// Sort the transaction by effective tip in ascending sort.
	txs := block.Transactions()
	sortedTxs := make([]*types.Transaction, len(txs))
	copy(sortedTxs, txs)
	baseFee := block.Header().BaseFee
	slices.SortFunc(sortedTxs, func(a, b *types.Transaction) int {
		tip1 := a.EffectiveGasTip(baseFee)
		tip2 := b.EffectiveGasTip(baseFee)
		return tip1.Cmp(tip2)
	})

	var prices []*big.Int
	for _, tx := range sortedTxs {
		tip := tx.EffectiveGasTip(baseFee)
		sender, err := types.Sender(signer, tx)
		if err == nil && sender != block.Rewardbase() {
			prices = append(prices, tip)
			if len(prices) >= limit {
				break
			}
		}
	}
	select {
	case result <- results{prices, nil}:
	case <-quit:
	}
}

func (oracle *Oracle) PurgeCache() {
	oracle.cacheLock.Lock()
	oracle.historyCache.Purge()
	oracle.cacheLock.Unlock()
}