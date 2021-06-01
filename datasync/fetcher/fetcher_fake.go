// Copyright 2020 The klaytn Authors
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

package fetcher

import (
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

type FakeFetcher struct{}

func NewFakeFetcher() *FakeFetcher {
	logger.Warn("fetcher is disabled; no data will be fetched from peers")
	return &FakeFetcher{}
}

func (*FakeFetcher) Enqueue(peer string, block *types.Block) error { return nil }
func (*FakeFetcher) FilterBodies(peer string, transactions [][]*types.Transaction, time time.Time) [][]*types.Transaction {
	return nil
}
func (*FakeFetcher) FilterHeaders(peer string, headers []*types.Header, time time.Time) []*types.Header {
	return nil
}
func (*FakeFetcher) Notify(peer string, hash common.Hash, number uint64, time time.Time, headerFetcher HeaderRequesterFn, bodyFetcher BodyRequesterFn) error {
	return nil
}
func (*FakeFetcher) Start() {}
func (*FakeFetcher) Stop()  {}
