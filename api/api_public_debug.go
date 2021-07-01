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

	"github.com/davecgh/go-spew/spew"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/rlp"
)

// PublicDebugAPI is the collection of Klaytn APIs exposed over the public
// debugging endpoint.
type PublicDebugAPI struct {
	b Backend
}

// NewPublicDebugAPI creates a new API definition for the public debug methods
// of the Klaytn service.
func NewPublicDebugAPI(b Backend) *PublicDebugAPI {
	return &PublicDebugAPI{b: b}
}

// GetBlockRlp retrieves the RLP encoded for of a single block.
func (api *PublicDebugAPI) GetBlockRlp(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (string, error) {
	block, _ := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if block == nil {
		blockNumberOrHashString, _ := blockNrOrHash.NumberOrHashString()
		return "", fmt.Errorf("block %v not found", blockNumberOrHashString)
	}
	encoded, err := rlp.EncodeToBytes(block)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", encoded), nil
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *PublicDebugAPI) PrintBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (string, error) {
	block, _ := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if block == nil {
		blockNumberOrHashString, _ := blockNrOrHash.NumberOrHashString()
		return "", fmt.Errorf("block %v not found", blockNumberOrHashString)
	}
	return spew.Sdump(block), nil
}
