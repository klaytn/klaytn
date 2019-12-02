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

package dbsyncer

import (
	"github.com/klaytn/klaytn/blockchain/types"
)

type MakeQueryRequest struct {
	block       *types.Block
	blockBase   uint64
	txIndexBase int
	txs         []*types.Transaction
	receipts    []*types.Receipt
	inc         int

	result chan *MakeQueryResult
}

type BulkInsertRequest struct {
	bulkInserts []*BulkInsertQuery
	inc         int

	result chan *BulkInsertResult
}

type BulkInsertResult struct {
	err error
}

type MakeQueryResult struct {
	block *types.Block

	cols string
	val  []interface{}

	scols string
	sval  []interface{}
	count int

	tcols  string
	tval   []interface{}
	tcount int

	err error
}

// QueryEngine is a helper structure to concurrently making insert query
type QueryEngine struct {
	ds          *DBSyncer
	taskQueue   int
	insertQueue int
	tasks       chan *MakeQueryRequest
	inserts     chan *BulkInsertRequest
}

func newQueryEngine(ds *DBSyncer, taskQueue int, insertQueue int) *QueryEngine {
	queryEngine := &QueryEngine{
		ds:          ds,
		tasks:       make(chan *MakeQueryRequest, taskQueue),
		inserts:     make(chan *BulkInsertRequest, insertQueue),
		taskQueue:   taskQueue,
		insertQueue: insertQueue,
	}
	for i := 0; i < taskQueue; i++ {
		go queryEngine.processing()
	}
	for i := 0; i < insertQueue; i++ {
		go queryEngine.executing()
	}

	return queryEngine
}

func (qe *QueryEngine) make(block *types.Block, txKey uint64, tx *types.Transaction, receipt *types.Receipt, result chan *MakeQueryResult) {

	cols, val, txMapArg, summaryArg, err := MakeTxDBRow(block, txKey, tx, receipt)
	scols, sval, count, serr := MakeSummaryDBRow(summaryArg)
	tcols, tval, tcount, terr := MakeTxMappingRow(txMapArg)

	defer func() {
		// recover from panic caused by writing to a closed channel
		if r := recover(); r != nil {
			logger.Error("channel closed", "err", r)
			return
		}
	}()
	if err == nil && serr == nil && terr == nil {
		result <- &MakeQueryResult{block, cols, val, scols, sval, count, tcols, tval, tcount, nil}
	} else {
		if err != nil {
			logger.Error("fail to make row (tx)", "err", err)
			result <- &MakeQueryResult{block, cols, val, scols, sval, count, tcols, tval, tcount, err}
		}
		if serr != nil {
			logger.Error("fail to make row (summary)", "err", serr)
			result <- &MakeQueryResult{block, cols, val, scols, sval, count, tcols, tval, tcount, serr}
		}
		if terr != nil {
			logger.Error("fail to make row (senderHash)", "err", terr)
			result <- &MakeQueryResult{block, cols, val, scols, sval, count, tcols, tval, tcount, terr}
		}
	}
}

func (qe *QueryEngine) execute(insertQuery *BulkInsertQuery, result chan *BulkInsertResult) {
	defer func() {
		// recover from panic caused by writing to a closed channel
		if r := recover(); r != nil {
			logger.Error("channel closed", "err", r)
			return
		}
	}()
	if err := qe.ds.bulkInsert(insertQuery.parameters, insertQuery.vals, insertQuery.blockNumber, insertQuery.insertCount); err != nil {
		logger.Error("fail to bulkinsert (tx/summary/senderHash)", "err", err)
		result <- &BulkInsertResult{err}
	}
	result <- &BulkInsertResult{nil}
}

func (qe *QueryEngine) processing() {
	for task := range qe.tasks {
		for i := 0; i < len(task.txs); i += task.inc {
			txKey := task.blockBase + uint64(task.txIndexBase+i)
			qe.make(task.block, txKey, task.txs[i], task.receipts[i], task.result)
		}
	}
}

func (qe *QueryEngine) executing() {
	for insert := range qe.inserts {
		for i := 0; i < len(insert.bulkInserts); i += insert.inc {
			qe.execute(insert.bulkInserts[i], insert.result)
		}
	}
}

func (qe *QueryEngine) insertTransactions(block *types.Block, txs types.Transactions, receipts types.Receipts, result chan *MakeQueryResult) {
	// If there's nothing to recover, abort
	if len(txs) == 0 {
		return
	}

	blockBase := block.NumberU64() * TX_KEY_FACTOR

	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := qe.taskQueue
	if len(txs) < tasks*4 {
		tasks = (len(txs) + 3) / 4
	}
	for i := 0; i < tasks; i++ {
		qe.tasks <- &MakeQueryRequest{
			block:       block,
			blockBase:   blockBase,
			txIndexBase: i,
			txs:         txs[i:],
			receipts:    receipts[i:],
			inc:         tasks,
			result:      result,
		}
	}
}

func (qe *QueryEngine) executeBulkInsert(bulkInserts []*BulkInsertQuery, result chan *BulkInsertResult) {
	// If there's nothing to recover, abort
	if len(bulkInserts) == 0 {
		return
	}

	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := qe.insertQueue
	if len(bulkInserts) < tasks {
		tasks = len(bulkInserts)
	}
	for i := 0; i < tasks; i++ {
		qe.inserts <- &BulkInsertRequest{
			bulkInserts: bulkInserts[i:],
			inc:         tasks,
			result:      result,
		}
	}
}
