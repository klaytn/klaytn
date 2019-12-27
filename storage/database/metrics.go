// Copyright 2019 The klaytn Authors
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

package database

import (
	"github.com/rcrowley/go-metrics"
)

var (
	cacheGetBlockBodyMissMeter = metrics.NewRegisteredMeter("klay/cache/get/blockbody/miss", nil)
	cacheGetBlockBodyHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/blockbody/hit", nil)

	cacheGetBlockBodyRLPMissMeter = metrics.NewRegisteredMeter("klay/cache/get/blockbodyrlp/miss", nil)
	cacheGetBlockBodyRLPHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/blockbodyrlp/hit", nil)

	cacheGetBlockMissMeter = metrics.NewRegisteredMeter("klay/cache/get/block/miss", nil)
	cacheGetBlockHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/block/hit", nil)

	cacheGetRecentTransactionsMissMeter = metrics.NewRegisteredMeter("klay/cache/get/transactions/miss", nil)
	cacheGetRecentTransactionsHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/transactions/hit", nil)

	cacheGetRecentBlockReceiptsMissMeter = metrics.NewRegisteredMeter("klay/cache/get/blockreceipts/miss", nil)
	cacheGetRecentBlockReceiptsHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/blockreceipts/hit", nil)

	cacheGetRecentTxReceiptMissMeter = metrics.NewRegisteredMeter("klay/cache/get/txreceipt/miss", nil)
	cacheGetRecentTxReceiptHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/txreceipt/hit", nil)

	cacheGetHeaderMissMeter = metrics.NewRegisteredMeter("klay/cache/get/header/miss", nil)
	cacheGetHeaderHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/header/hit", nil)

	cacheGetTDMissMeter = metrics.NewRegisteredMeter("klay/cache/get/td/miss", nil)
	cacheGetTDHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/td/hit", nil)

	cacheGetBlockNumberMissMeter = metrics.NewRegisteredMeter("klay/cache/get/blocknumber/miss", nil)
	cacheGetBlockNumberHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/blocknumber/hit", nil)

	cacheGetCanonicalHashMissMeter = metrics.NewRegisteredMeter("klay/cache/get/canonicalhash/miss", nil)
	cacheGetCanonicalHashHitMeter  = metrics.NewRegisteredMeter("klay/cache/get/canonicalhash/hit", nil)
)
