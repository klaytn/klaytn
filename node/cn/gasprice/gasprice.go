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
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
)

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
	SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription
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
	headEvent := make(chan blockchain.ChainHeadEvent, 1)
	backend.SubscribeChainHeadEvent(headEvent)
	go func() {
		var lastHead common.Hash
		for ev := range headEvent {
			if ev.Block.ParentHash() != lastHead {
				cache.Purge()
			}
			lastHead = ev.Block.Hash()
		}
	}()
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
// Though Klaytn does not recognize gas tip, this function returns some value for compatibility.
func (gpo *Oracle) SuggestTipCap(ctx context.Context) (*big.Int, error) {
	if gpo.txPool == nil {
		// If txpool is not set, just return 0. This is used for testing.
		return common.Big0, nil
	}

	nextNum := new(big.Int).Add(gpo.backend.CurrentBlock().Number(), common.Big1)
	if gpo.backend.ChainConfig().IsDragonForkEnabled(nextNum) {
		// TODO: After Dragon, return 60% percentile of last 20 blocks
		return common.Big0, nil
	} else if gpo.backend.ChainConfig().IsMagmaForkEnabled(nextNum) {
		// After Magma, return zero
		return common.Big0, nil
	} else {
		// Before Magma, return the fixed UnitPrice.
		unitPrice := gpo.txPool.GasPrice()
		return unitPrice, nil
	}
}
