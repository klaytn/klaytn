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
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/rpc"
)

// PrivateDebugAPI is the collection of Klaytn APIs exposed over the private
// debugging endpoint.
type PrivateDebugAPI struct {
	b Backend
}

// NewPrivateDebugAPI creates a new API definition for the private debug methods
// of the Klaytn service.
func NewPrivateDebugAPI(b Backend) *PrivateDebugAPI {
	return &PrivateDebugAPI{b: b}
}

// ChaindbProperty returns leveldb properties of the chain database.
func (api *PrivateDebugAPI) ChaindbProperty(property string) (string, error) {
	return api.b.ChainDB().Stat(property)
}

// ChaindbCompact flattens the entire key-value database into a single level,
// removing all unused slots and merging all keys.
func (api *PrivateDebugAPI) ChaindbCompact() error {
	for b := 0; b <= 255; b++ {
		var (
			start = []byte{byte(b)}
			end   = []byte{byte(b + 1)}
		)
		if b == 255 {
			end = nil
		}
		logger.Info("Compacting database started", "range", fmt.Sprintf("%#X-%#X", start, end))
		cstart := time.Now()
		if err := api.b.ChainDB().Compact(start, end); err != nil {
			logger.Error("Database compaction failed", "err", err)
			return err
		}
		logger.Info("Compacting database completed", "range", fmt.Sprintf("%#X-%#X", start, end), "elapsed", common.PrettyDuration(time.Since(cstart)))
	}
	return nil
}

// SetHead rewinds the head of the blockchain to a previous block.
func (api *PrivateDebugAPI) SetHead(number rpc.BlockNumber) error {
	if number == rpc.PendingBlockNumber ||
		number == rpc.LatestBlockNumber ||
		number.Uint64() > api.b.CurrentBlock().NumberU64() {
		return errors.New("Cannot rewind to future")
	}
	return api.b.SetHead(uint64(number))
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *PrivateDebugAPI) PrintBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (string, error) {
	block, _ := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if block == nil {
		blockNumberOrHashString, _ := blockNrOrHash.NumberOrHashString()
		return "", fmt.Errorf("block %v not found", blockNumberOrHashString)
	}
	return spew.Sdump(block), nil
}
