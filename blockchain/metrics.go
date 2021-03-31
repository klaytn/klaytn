// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from eth/metrics.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"github.com/rcrowley/go-metrics"
)

var (
	cacheGetFutureBlockMissMeter = metrics.NewRegisteredMeter("klay/cache/get/futureblock/miss", nil)
	cacheGetFutureBlockHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/futureblock/hit", nil)

	cacheGetBadBlockMissMeter = metrics.NewRegisteredMeter("klay/cache/get/badblock/miss", nil)
	cacheGetBadBlockHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/badblock/hit", nil)

	trieDBNodesSizeBytesGauge = metrics.NewRegisteredGauge("klay/triedb/nodessizebytes", nil)
	trieDBPreimagesSizeGauge  = metrics.NewRegisteredGauge("klay/triedb/preimagessizebytes", nil)

	headBlockNumberGauge = metrics.NewRegisteredGauge("blockchain/head/blocknumber", nil)
	blockTxCountsGauge   = metrics.NewRegisteredGauge("blockchain/block/tx/gauge", nil)
	blockTxCountsCounter = metrics.NewRegisteredCounter("blockchain/block/tx/counter", nil)
	// the counter to record a bad block, increases 1 if bad block occurs
	badBlockCounter = metrics.NewRegisteredCounter("blockchain/bad/block/counter", nil)

	txPoolPendingGauge = metrics.NewRegisteredGauge("tx/pool/pending/gauge", nil)
	txPoolQueueGauge   = metrics.NewRegisteredGauge("tx/pool/queue/gauge", nil)
)
