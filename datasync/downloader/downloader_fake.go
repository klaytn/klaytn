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

package downloader

import (
	"math/big"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// fakeDownloader do nothing
type FakeDownloader struct{}

func NewFakeDownloader() *FakeDownloader {
	logger.Warn("downloader is disabled; no data will be downloaded from peers")
	return &FakeDownloader{}
}

func (*FakeDownloader) RegisterPeer(id string, version int, peer Peer) error { return nil }
func (*FakeDownloader) UnregisterPeer(id string) error                       { return nil }

func (*FakeDownloader) DeliverBodies(id string, transactions [][]*types.Transaction) error {
	return nil
}
func (*FakeDownloader) DeliverHeaders(id string, headers []*types.Header) error      { return nil }
func (*FakeDownloader) DeliverNodeData(id string, data [][]byte) error               { return nil }
func (*FakeDownloader) DeliverReceipts(id string, receipts [][]*types.Receipt) error { return nil }

func (*FakeDownloader) Terminate() {}
func (*FakeDownloader) Synchronise(id string, head common.Hash, td *big.Int, mode SyncMode) error {
	return nil
}
func (*FakeDownloader) Progress() klaytn.SyncProgress { return klaytn.SyncProgress{} }
func (*FakeDownloader) Cancel()                       {}
