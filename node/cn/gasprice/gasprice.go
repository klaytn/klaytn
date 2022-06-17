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

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/networks/rpc"

	"github.com/klaytn/klaytn/common"
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
	return &Oracle{
		backend:          backend,
		lastPrice:        params.Default,
		checkBlocks:      blocks,
		maxEmpty:         blocks / 2,
		maxBlocks:        blocks * 5,
		percentile:       percent,
		maxHeaderHistory: params.MaxHeaderHistory,
		maxBlockHistory:  params.MaxBlockHistory,
		txPool:           txPool,
	}
}

// SuggestPrice returns the recommended gas price.
func (gpo *Oracle) SuggestPrice(ctx context.Context) (*big.Int, error) {

	if gpo.txPool == nil {
		// If txpool is not set, just return 0. This is used for testing.
		return common.Big0, nil
	}
	// Since we have fixed gas price, we can directly get this value from TxPool.
	return gpo.txPool.GasPrice(), nil
	/*
		// TODO-Klaytn-RemoveLater Later remove below obsolete code if we don't need them anymore.
		gpo.cacheLock.RLock()
		lastHead := gpo.lastHead
		lastPrice := gpo.lastPrice
		gpo.cacheLock.RUnlock()

		head, _ := gpo.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
		headHash := head.Hash()
		if headHash == lastHead {
			return lastPrice, nil
		}

		gpo.fetchLock.Lock()
		defer gpo.fetchLock.Unlock()

		// try checking the cache again, maybe the last fetch fetched what we need
		gpo.cacheLock.RLock()
		lastHead = gpo.lastHead
		lastPrice = gpo.lastPrice
		gpo.cacheLock.RUnlock()
		if headHash == lastHead {
			return lastPrice, nil
		}

		blockNum := head.Number.Uint64()
		ch := make(chan getBlockPricesResult, gpo.checkBlocks)
		sent := 0
		exp := 0
		var blockPrices []*big.Int
		for sent < gpo.checkBlocks && blockNum > 0 {
			go gpo.getBlockPrices(ctx, types.MakeSigner(gpo.backend.ChainConfig(), big.NewInt(int64(blockNum))), blockNum, ch)
			sent++
			exp++
			blockNum--
		}
		maxEmpty := gpo.maxEmpty
		for exp > 0 {
			res := <-ch
			if res.err != nil {
				return lastPrice, res.err
			}
			exp--
			if res.price != nil {
				blockPrices = append(blockPrices, res.price)
				continue
			}
			if maxEmpty > 0 {
				maxEmpty--
				continue
			}
			if blockNum > 0 && sent < gpo.maxBlocks {
				go gpo.getBlockPrices(ctx, types.MakeSigner(gpo.backend.ChainConfig(), big.NewInt(int64(blockNum))), blockNum, ch)
				sent++
				exp++
				blockNum--
			}
		}
		price := lastPrice
		if len(blockPrices) > 0 {
			sort.Sort(bigIntArray(blockPrices))
			price = blockPrices[(len(blockPrices)-1)*gpo.percentile/100]
		}
		if price.Cmp(maxPrice) > 0 {
			price = new(big.Int).Set(maxPrice)
		}

		gpo.cacheLock.Lock()
		gpo.lastHead = headHash
		gpo.lastPrice = price
		gpo.cacheLock.Unlock()
		return price, nil
	*/
}

// TODO-Klaytn-RemoveLater Later remove below obsolete code if we don't need them anymore.
//type getBlockPricesResult struct {
//	price *big.Int
//	err   error
//}
//
//type transactionsByGasPrice []*types.Transaction
//
//func (t transactionsByGasPrice) Len() int           { return len(t) }
//func (t transactionsByGasPrice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
//func (t transactionsByGasPrice) Less(i, j int) bool { return t[i].GasPrice().Cmp(t[j].GasPrice()) < 0 }
//
//type bigIntArray []*big.Int
//
//func (s bigIntArray) Len() int           { return len(s) }
//func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
//func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
