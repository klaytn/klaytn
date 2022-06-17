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
	"errors"
	"strings"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
)

func (ds *DBSyncer) HandleChainEventParallel(block *types.Block) error {
	logger.Info("dbsyncer(multi) HandleChainEvent", "number", block.Number(), "txs", block.Transactions().Len())

	bulkInsertQuerys := []*BulkInsertQuery{}

	startblock := time.Now()
	bulkInsertQuerys, err := ds.parallelSyncBlockHeader(block)
	if err != nil {
		logger.Error("fail to sync block", "block", block.Number(), "err", err)
		return err
	}
	blocktime := time.Since(startblock)

	starttx := time.Now()
	receipttime := time.Duration(0)
	if block.Transactions().Len() > 0 {

		startrceipt := time.Now()
		receipts := ds.blockchain.GetReceiptsByBlockHash(block.Hash())
		receipttime = time.Since(startrceipt)

		bulkInsertQuerys, err = ds.parallelSyncTransactions(block, receipts, bulkInsertQuerys)
		if err != nil {
			logger.Error("fail to sync transaction", "block", block.Number(), "err", err)
			return err
		}
	}
	txtime := time.Since(starttx)

	startInsert := time.Now()
	queryResult := make(chan *BulkInsertResult, len(bulkInsertQuerys))
	ds.queryEngine.executeBulkInsert(bulkInsertQuerys, queryResult)

	totalResult := 0
	defer close(queryResult)

RESULT:
	for result := range queryResult {
		if result.err != nil {
			return result.err
		}

		totalResult++
		if totalResult == len(bulkInsertQuerys) {
			break RESULT
		}
	}
	insertTime := time.Since(startInsert)

	totalTime := time.Since(startblock)
	if ds.logMode {
		logger.Info("dbsyncer(multi) time", "number", block.Number(), "block", blocktime, "txs", txtime, "receipts", receipttime, "insert", insertTime, "total", totalTime)
	}

	return nil
}

func (ds *DBSyncer) parallelSyncBlockHeader(block *types.Block) ([]*BulkInsertQuery, error) {
	blockStr := ds.blockInsertQuery + ","
	vals := []interface{}{}

	bulkInsertQuerys := []*BulkInsertQuery{}

	proposerAddr, committeeAddrs, err := getProposerAndValidatorsFromBlock(block)
	if err != nil {
		return nil, err
	}

	totalTx := block.Transactions().Len()
	committee := strings.ToLower(committeeAddrs)
	gasUsed := block.Header().GasUsed
	gasPrice := ds.blockchain.Config().UnitPrice
	hash := block.Header().Hash().Hex()
	number := block.Header().Number.Uint64()
	parentHash := block.Header().ParentHash.Hex()
	proposer := strings.ToLower(proposerAddr)
	reward := block.Header().Rewardbase.Hex()
	size := block.Size()
	timestamp := block.Header().Time.String()
	timestampFos := block.Header().TimeFoS

	vals = append(vals, totalTx, committee, gasUsed, gasPrice, hash,
		number, parentHash, proposer, reward, size, timestamp, timestampFos)

	bulkInsertQuerys = append(bulkInsertQuerys, &BulkInsertQuery{blockStr, vals, block.NumberU64(), 1})

	return bulkInsertQuerys, nil
}

func (ds *DBSyncer) parallelSyncTransactions(block *types.Block, receipts types.Receipts, bulkInsertQuerys []*BulkInsertQuery) ([]*BulkInsertQuery, error) {

	txStr, vals, insertCount := ds.resetTxParameter()
	summaryStr, summaryVals, summaryInsertCount := ds.resetSummaryParameter()
	txMapStr, txMapVals, txMapInsertCount := ds.resetTxMapParameter()

	txLen := block.Transactions().Len()
	result := make(chan *MakeQueryResult, txLen)

	if txLen != receipts.Len() {
		logger.Error("transactions is not matched receipts", "txs", txLen, "receipts", receipts.Len())
		return nil, errors.New("transaction count is not matched receipts")
	}

	ds.queryEngine.insertTransactions(block, block.Transactions(), receipts, result)
	totalTxs := 0

	defer close(result)

QUERY:
	for record := range result {

		if record.err != nil {
			return nil, record.err
		}

		txStr += record.cols + ","
		vals = append(vals, record.val...)
		insertCount++

		if insertCount >= ds.bulkInsertSize {
			bulkInsertQuerys = append(bulkInsertQuerys, &BulkInsertQuery{txStr, vals, block.NumberU64(), insertCount})

			txStr, vals, insertCount = ds.resetTxParameter()
		}

		if record.count == 1 {
			summaryStr += record.scols + ","
			summaryVals = append(summaryVals, record.sval...)
			summaryInsertCount++
		}

		if summaryInsertCount >= ds.bulkInsertSize {
			bulkInsertQuerys = append(bulkInsertQuerys, &BulkInsertQuery{summaryStr, summaryVals, block.NumberU64(), summaryInsertCount})

			summaryStr, summaryVals, summaryInsertCount = ds.resetTxParameter()
		}

		if record.tcount == 1 {
			txMapStr += record.tcols + ","
			txMapVals = append(txMapVals, record.tval...)
			txMapInsertCount++
		}

		if txMapInsertCount >= ds.bulkInsertSize {
			bulkInsertQuerys = append(bulkInsertQuerys, &BulkInsertQuery{txMapStr, txMapVals, block.NumberU64(), txMapInsertCount})

			txMapStr, txMapVals, txMapInsertCount = ds.resetTxMapParameter()
		}

		totalTxs++
		if totalTxs == block.Transactions().Len() {
			break QUERY
		}
	}

	if insertCount > 0 {
		bulkInsertQuerys = append(bulkInsertQuerys, &BulkInsertQuery{txStr, vals, block.NumberU64(), insertCount})
	}

	if summaryInsertCount > 0 {
		bulkInsertQuerys = append(bulkInsertQuerys, &BulkInsertQuery{summaryStr, summaryVals, block.NumberU64(), summaryInsertCount})
	}

	if txMapInsertCount > 0 {
		bulkInsertQuerys = append(bulkInsertQuerys, &BulkInsertQuery{txMapStr, txMapVals, block.NumberU64(), txMapInsertCount})
	}

	return bulkInsertQuerys, nil
}
