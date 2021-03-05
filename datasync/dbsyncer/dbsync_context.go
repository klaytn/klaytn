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
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
)

// HandleChainEventContext supports 2PC Commit (insert block + insert txs) for data consistency
// @TODO-Klaytn improve performance, too slower than HanleChainEvent()
func (ds *DBSyncer) HandleChainEventContext(block *types.Block) error {
	logger.Info("dbsyncer HandleChainEvent", "number", block.Number(), "txs", block.Transactions().Len())
	startblock := time.Now()

	ctx, cancel := context.WithTimeout(ds.ctx, 90*time.Second)
	defer cancel()

	tx, err := ds.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		logger.Error("fail to begin tx", "err", err)
		return err
	}

	if err := ds.syncBlockHeaderContext(ctx, tx, block); err != nil {
		logger.Error("fail to sync block", "block", block.Number(), "err", err)
		if rerr := tx.Rollback(); rerr != nil {
			logger.Error("fail to rollback tx", "block", block.Number(), "err", rerr)
		}
		return err
	}

	blocktime := time.Since(startblock)
	starttx := time.Now()

	if block.Transactions().Len() > 0 {

		if err := ds.syncTransactionsContext(ctx, tx, block); err != nil {
			logger.Error("fail to sync transaction", "block", block.Number(), "err", err)
			if rerr := tx.Rollback(); rerr != nil {
				logger.Error("fail to rollback tx", "block", block.Number(), "err", rerr)
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("fail to commit tx", "block", block.Number(), "err", err)
		return err
	}

	txtime := time.Since(starttx)
	totalTime := time.Since(startblock)
	if ds.logMode {
		logger.Info("dbsync time", "number", block.Number(), "block", blocktime, "txs", txtime, "total", totalTime)
	}

	return nil
}

func (ds *DBSyncer) syncBlockHeaderContext(ctx context.Context, tx *sql.Tx, block *types.Block) error {
	proposerAddr, committeeAddrs, err := getProposerAndValidatorsFromBlock(block)

	totalTx := block.Transactions().Len()
	committee := strings.ToLower(committeeAddrs)
	gasUsed := block.Header().GasUsed
	gasPrice := ds.blockchain.Config().UnitPrice
	hash := block.Header().Hash().Hex()
	number := block.Header().Number.Uint64()
	parentHash := block.Header().ParentHash.Hex()
	proposer := strings.ToLower(proposerAddr)
	reward := strings.ToLower(block.Header().Rewardbase.Hex())
	size := block.Size()
	timestamp := block.Header().Time.String()
	timestampFos := block.Header().TimeFoS

	stmtIns, err := ds.db.Prepare(ds.blockInsertQuery)

	if err != nil {
		logger.Error("fail to prepare (block)", "query", ds.blockInsertQuery)
		return err
	}
	tStmtIns := tx.StmtContext(ctx, stmtIns)

	defer func() {
		if err := tStmtIns.Close(); err != nil {
			logger.Error("fail to close stmt", "err", err)
		}
	}()

	if _, err := tStmtIns.Exec(totalTx, committee, gasUsed, gasPrice, hash,
		number, parentHash, proposer, reward, size, timestamp, timestampFos); err != nil {
		logger.Error("fail to insert DB (block)", "number", block.Number(), "err", err)
		return err
	}

	return nil
}

func (ds *DBSyncer) syncTransactionsContext(ctx context.Context, syncTx *sql.Tx, block *types.Block) error {
	txKey := block.NumberU64() * TX_KEY_FACTOR
	txStr, vals, insertCount := ds.resetTxParameter()
	summaryStr, summaryVals, summaryInsertCount := ds.resetSummaryParameter()
	txMapStr, txMapVals, txMapInsertCount := ds.resetTxMapParameter()

	receipts := ds.blockchain.GetReceiptsByBlockHash(block.Hash())

	for index, tx := range block.Transactions() {
		txKey += uint64(index)
		cols, val, txMapArg, summaryArg, err := MakeTxDBRow(block, txKey, tx, receipts[index])
		if err != nil {
			return err
		}

		txStr += cols + ","
		vals = append(vals, val...)
		insertCount++

		if insertCount >= ds.bulkInsertSize {
			if err := ds.bulkInsertContext(ctx, syncTx, txStr, vals, block, insertCount); err != nil {
				return err
			}
			txStr, vals, insertCount = ds.resetTxParameter()
		}

		scols, sval, count, err := MakeSummaryDBRow(summaryArg)

		if count == 1 {
			summaryStr += scols + ","
			summaryVals = append(summaryVals, sval...)
			summaryInsertCount++
		}

		if summaryInsertCount >= ds.bulkInsertSize {
			if err := ds.bulkInsertContext(ctx, syncTx, summaryStr, summaryVals, block, summaryInsertCount); err != nil {
				return err
			}
			summaryStr, summaryVals, summaryInsertCount = ds.resetTxParameter()
		}

		tcols, tval, tcount, err := MakeTxMappingRow(txMapArg)
		if err != nil {
			return err
		}

		if tcount == 1 {
			txMapStr += tcols + ","
			txMapVals = append(txMapVals, tval...)
			txMapInsertCount++
		}

		if txMapInsertCount >= ds.bulkInsertSize {
			if err := ds.bulkInsertContext(ctx, syncTx, txMapStr, txMapVals, block, txMapInsertCount); err != nil {
				return err
			}
			txMapStr, txMapVals, txMapInsertCount = ds.resetTxMapParameter()
		}
	}

	if insertCount > 0 {
		if err := ds.bulkInsertContext(ctx, syncTx, txStr, vals, block, insertCount); err != nil {
			return err
		}
	}

	if summaryInsertCount > 0 {
		if err := ds.bulkInsertContext(ctx, syncTx, summaryStr, summaryVals, block, summaryInsertCount); err != nil {
			return err
		}
	}

	if txMapInsertCount > 0 {
		if err := ds.bulkInsertContext(ctx, syncTx, txMapStr, txMapVals, block, txMapInsertCount); err != nil {
			return err
		}
	}

	return nil
}

func (ds *DBSyncer) bulkInsertContext(ctx context.Context, tx *sql.Tx, sqlStr string, vals []interface{}, block *types.Block, insertCount int) error {
	start := time.Now()
	// trim the last
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmtTxs, err := ds.db.Prepare(sqlStr)
	if err != nil {
		logger.Error("fail to create prepare", "sql", sqlStr, "err", err)
		return err
	}
	tStmtTxs := tx.StmtContext(ctx, stmtTxs)

	if _, err := tStmtTxs.Exec(vals...); err != nil {
		logger.Error("fail to insert DB (tx)", "number", block.Number(), "err", err)
		return err
	}
	defer func() {
		if err := tStmtTxs.Close(); err != nil {
			logger.Error("fail to close stmt", "err", err)
		}
	}()

	if ds.logMode {
		txTime := time.Since(start)
		logger.Info("TX INSERT", "counts", insertCount, "time", txTime)
	}

	return nil
}
