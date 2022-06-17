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

	_ "github.com/go-sql-driver/mysql"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/work"
	"github.com/pkg/errors"
)

var logger = log.NewModuleLogger(log.Node)

type DBSyncer struct {
	cfg        *DBConfig
	dataSource string

	blockchain *blockchain.BlockChain

	// chain event
	chainCh     chan blockchain.ChainEvent
	chainHeadCh chan blockchain.ChainHeadEvent
	chainSub    event.Subscription
	logsCh      chan []*types.Log
	logsSub     event.Subscription

	ctx  context.Context
	stop context.CancelFunc
	db   *sql.DB

	logMode bool

	blockInsertQuery     string
	txInsertQuery        string
	summaryInsertQuery   string
	txHashMapInsertQuery string

	HandleBlock func(block *types.Block) error
	queryEngine *QueryEngine

	bulkInsertSize int

	eventMode string

	maxBlockDiff uint64
}

func NewDBSyncer(ctx *node.ServiceContext, cfg *DBConfig) (*DBSyncer, error) {

	logger.Info("initialize DBSyncer", "db.host",
		cfg.DBHost, "db.port", cfg.DBPort, "db.name", cfg.DBName, "db.user", cfg.DBUser, "db.max.idle",
		cfg.MaxIdleConns, "db.password", cfg.DBPassword, "db.max.open", cfg.MaxOpenConns, "db.max.lifetime",
		cfg.ConnMaxLifetime, "block.ch.size", cfg.BlockChannelSize, "mode", cfg.Mode, "genquery.th",
		cfg.GenQueryThread, "insert.th", cfg.InsertThread, "bulk.size", cfg.BulkInsertSize, "event.mode",
		cfg.EventMode, "max.block.diff", cfg.MaxBlockDiff)

	if cfg.DBHost == "" {
		return nil, errors.New("db config must be set (db.host)")
	} else if cfg.DBName == "" {
		return nil, errors.New("db config must be set (db.name)")
	} else if cfg.DBUser == "" {
		return nil, errors.New("db config must be set (db.user)")
	} else if cfg.DBPassword == "" {
		return nil, errors.New("db config must be set (db.password)")
	}

	return &DBSyncer{
		cfg:            cfg,
		logMode:        cfg.EnabledLogMode,
		bulkInsertSize: cfg.BulkInsertSize,
		eventMode:      cfg.EventMode,
		maxBlockDiff:   cfg.MaxBlockDiff,
	}, nil
}

func (ds *DBSyncer) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

func (ds *DBSyncer) APIs() []rpc.API {
	return []rpc.API{}
}

func (ds *DBSyncer) Start(server p2p.Server) error {

	ds.dataSource = ds.cfg.DBUser + ":" + ds.cfg.DBPassword + "@tcp(" + ds.cfg.DBHost + ":" +
		ds.cfg.DBPort + ")/" + ds.cfg.DBName + "?writeTimeout=10s&timeout=10s"

	db, err := sql.Open("mysql", ds.dataSource)
	if err != nil {
		logger.Error("fail to connect database", "target", ds.dataSource)
		return err
	}
	ds.db = db
	ds.db.SetMaxIdleConns(ds.cfg.MaxIdleConns)
	ds.db.SetConnMaxLifetime(ds.cfg.ConnMaxLifetime)
	ds.db.SetMaxOpenConns(ds.cfg.MaxOpenConns)

	// initialize context
	ds.ctx, ds.stop = context.WithCancel(context.Background())

	// query
	ds.blockInsertQuery = "INSERT INTO " + ds.cfg.DBName + ".block " + "(totalTx, " +
		"committee, gasUsed, gasPrice, hash, " +
		"number, parentHash, proposer, reward, size, " +
		"timestamp, timestampFoS)" + "VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

	ds.txInsertQuery = "INSERT INTO " + ds.cfg.DBName + ".transaction " + "(id, blockHash, blockNumber, " +
		"contractAddress, `from`, gas, gasPrice, gasUsed, input, nonce, status, `to`, " +
		"timestamp, txHash, type, value, feePayer, feeRatio, senderTxHash) VALUES "

	ds.summaryInsertQuery = "INSERT INTO " + ds.cfg.DBName + ".account_summary " + "(address, type, " +
		"creator, created_tx, hra) VALUES "

	ds.txHashMapInsertQuery = "INSERT INTO " + ds.cfg.DBName + ".sendertxhash_map " + "(senderTxHash, txHash) VALUES "

	if ds.cfg.Mode == "single" {
		ds.HandleBlock = ds.HandleChainEvent
	} else if ds.cfg.Mode == "multi" {
		ds.HandleBlock = ds.HandleChainEventParallel
		ds.queryEngine = newQueryEngine(ds, ds.cfg.GenQueryThread, ds.cfg.InsertThread)
	} else if ds.cfg.Mode == "context" {
		ds.HandleBlock = ds.HandleChainEventContext
	} else {
		ds.HandleBlock = ds.HandleChainEventParallel
		ds.queryEngine = newQueryEngine(ds, ds.cfg.GenQueryThread, ds.cfg.InsertThread)
	}

	return nil
}

func (ds *DBSyncer) Stop() error {
	if ds.db != nil {
		if err := ds.db.Close(); err != nil {
			logger.Error("fail to close db", "err", err)
		}
	}

	return nil
}

func (ds *DBSyncer) Components() []interface{} {
	return nil
}

func (ds *DBSyncer) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			ds.blockchain = v
			// event from core-service
			if ds.eventMode == BLOCK_MODE {
				// handle all blocks when many blocks create
				ds.chainCh = make(chan blockchain.ChainEvent, ds.cfg.BlockChannelSize)
				ds.chainSub = ds.blockchain.SubscribeChainEvent(ds.chainCh)
				// eventMode == "head"
			} else if ds.eventMode == HEAD_MODE {
				// handle last block when many blocks create
				ds.chainHeadCh = make(chan blockchain.ChainHeadEvent, ds.cfg.BlockChannelSize)
				ds.chainSub = ds.blockchain.SubscribeChainHeadEvent(ds.chainHeadCh)
			} else {
				logger.Error("unknown event.mode (block,head)", "current mode", ds.eventMode)
			}
			//ds.logsSub = ds.blockchain.SubscribeLogsEvent(ds.logsCh)
		case *blockchain.TxPool:
		case *work.Miner:
		}
	}

	go ds.loop()
}

func (ds *DBSyncer) loop() {
	report := time.NewTicker(1 * time.Minute)
	defer report.Stop()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle ChainEvent
		case ev := <-ds.chainCh:
			if ev.Block != nil {
				ds.HandleDiffBlock(ev.Block)
			} else {
				logger.Error("dbsyncer block event is nil")
			}
		case ev := <-ds.chainHeadCh:
			if ev.Block != nil {
				ds.HandleDiffBlock(ev.Block)
			} else {
				logger.Error("dbsyncer block event is nil")
			}
		case <-report.C:
			// check db health
			go ds.Ping()
		case err := <-ds.chainSub.Err():
			if err != nil {
				logger.Error("dbsyncer block subscription ", "err", err)
			}
			return
		}
	}
}

func (ds *DBSyncer) HandleDiffBlock(block *types.Block) {

	diff := ds.blockchain.CurrentBlock().NumberU64() - block.NumberU64()

	if ds.maxBlockDiff > 0 && diff > ds.maxBlockDiff {
		logger.Info("there are many block number difference (skip block)", "diff", diff, "skip-block", block.NumberU64())
	} else {
		if err := ds.HandleBlock(block); err != nil {
			logger.Error("dbsyncer block event", "block", block.Number(), "err", err)
		}
	}
}

func (ds *DBSyncer) Ping() {
	logger.Info("check database", "target", ds.dataSource)
	ctx, cancel := context.WithTimeout(ds.ctx, 10*time.Second)
	defer cancel()

	err := ds.db.PingContext(ctx)
	if err != nil {
		logger.Error("database down", "target", ds.dataSource, "err", err)
	}
}

func (ds *DBSyncer) HandleChainEvent(block *types.Block) error {
	logger.Info("dbsyncer HandleChainEvent", "number", block.Number(), "txs", block.Transactions().Len())
	startblock := time.Now()

	if err := ds.syncBlockHeader(block); err != nil {
		logger.Error("fail to sync block", "block", block.Number(), "err", err)
		return err
	}

	blocktime := time.Since(startblock)
	starttx := time.Now()

	if block.Transactions().Len() > 0 {
		if err := ds.SyncTransactions(block); err != nil {
			logger.Error("fail to sync transaction", "block", block.Number(), "err", err)
			return err
		}
	}

	txtime := time.Since(starttx)
	totalTime := time.Since(startblock)
	if ds.logMode {
		logger.Info("dbsync time", "number", block.Number(), "block", blocktime, "txs", txtime, "total", totalTime)
	}

	return nil
}

func (ds *DBSyncer) syncBlockHeader(block *types.Block) error {
	proposerAddr, committeeAddrs, err := getProposerAndValidatorsFromBlock(block)

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

	stmtIns, err := ds.db.Prepare(ds.blockInsertQuery)

	if err != nil {
		logger.Error("fail to prepare (block)", "query", ds.blockInsertQuery)
		return err
	}

	defer func() {
		if err := stmtIns.Close(); err != nil {
			logger.Error("fail to close stmt", "err", err)
		}
	}()

	if _, err := stmtIns.Exec(totalTx, committee, gasUsed, gasPrice, hash,
		number, parentHash, proposer, reward, size, timestamp, timestampFos); err != nil {
		logger.Error("fail to insert DB (block)", "number", block.Number(), "err", err)
		return err
	}

	return nil
}

func (ds *DBSyncer) SyncTransactions(block *types.Block) error {

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
			if err := ds.bulkInsert(txStr, vals, block.NumberU64(), insertCount); err != nil {
				return err
			}
			txStr, vals, insertCount = ds.resetTxParameter()
		}

		scols, sval, count, err := MakeSummaryDBRow(summaryArg)
		if err != nil {
			return err
		}

		if count == 1 {
			summaryStr += scols + ","
			summaryVals = append(summaryVals, sval...)
			summaryInsertCount++
		}

		if summaryInsertCount >= ds.bulkInsertSize {
			if err := ds.bulkInsert(summaryStr, summaryVals, block.NumberU64(), summaryInsertCount); err != nil {
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
			if err := ds.bulkInsert(txMapStr, txMapVals, block.NumberU64(), txMapInsertCount); err != nil {
				return err
			}
			txMapStr, txMapVals, txMapInsertCount = ds.resetTxMapParameter()
		}
	}

	if insertCount > 0 {
		if err := ds.bulkInsert(txStr, vals, block.NumberU64(), insertCount); err != nil {
			return err
		}
	}

	if summaryInsertCount > 0 {
		if err := ds.bulkInsert(summaryStr, summaryVals, block.NumberU64(), summaryInsertCount); err != nil {
			return err
		}
	}

	if txMapInsertCount > 0 {
		if err := ds.bulkInsert(txMapStr, txMapVals, block.NumberU64(), txMapInsertCount); err != nil {
			return err
		}
	}

	return nil
}

func (ds *DBSyncer) resetTxParameter() (txStr string, vals []interface{}, insertCount int) {
	txStr = ds.txInsertQuery
	vals = []interface{}{}
	insertCount = 0

	return txStr, vals, insertCount
}

func (ds *DBSyncer) resetSummaryParameter() (summaryStr string, vals []interface{}, insertCount int) {
	summaryStr = ds.summaryInsertQuery
	vals = []interface{}{}
	insertCount = 0

	return summaryStr, vals, insertCount
}

func (ds *DBSyncer) resetTxMapParameter() (txMapStr string, vals []interface{}, insertCount int) {
	txMapStr = ds.txHashMapInsertQuery
	vals = []interface{}{}
	insertCount = 0

	return txMapStr, vals, insertCount
}

func (ds *DBSyncer) bulkInsert(sqlStr string, vals []interface{}, blockNumber uint64, insertCount int) error {
	start := time.Now()
	// trim the last
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	stmtTxs, err := ds.db.Prepare(sqlStr)
	if err != nil {
		logger.Error("fail to create prepare", "sql", sqlStr, "err", err)
		return err
	}

	if _, err := stmtTxs.Exec(vals...); err != nil {
		logger.Error("fail to insert DB (tx)", "number", blockNumber, "err", err)
		return err
	}
	defer func() {
		if err := stmtTxs.Close(); err != nil {
			logger.Error("fail to close stmt", "err", err)
		}
	}()

	if ds.logMode {
		txTime := time.Since(start)
		logger.Info("TX/BLOCK INSERT", "block", blockNumber, "counts", insertCount, "time", txTime)
	}

	return nil
}

func (ds *DBSyncer) HandleLogsEvent(logs []*types.Log) error {
	return nil
}
