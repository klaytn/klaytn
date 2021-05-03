// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const (
	dbMigrationFetchNum = 500
	reportCycle         = dbMigrationFetchNum * 100
)

func copyDB(name string, srcDB, dstDB Database) error {
	// settings for quit signal from os
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// create src iterator and dst batch
	srcIter := srcDB.NewIterator(nil, nil)
	dstBatch := dstDB.NewBatch()

	// vars for log
	start := time.Now()
	previousFetched, fetched := 0, 0

	cycleStart := time.Now()
	for fetched = 0; srcIter.Next(); fetched++ {
		// fetch keys and values
		// Contents of srcIter.Key() and srcIter.Value() should not be modified, and
		// only valid until the next call to Next.
		key := make([]byte, len(srcIter.Key()))
		val := make([]byte, len(srcIter.Value()))
		copy(key, srcIter.Key())
		copy(val, srcIter.Value())

		// write fetched keys and values to DB
		// If dstDB is dynamoDB, Put will Write when the number items reach dynamoBatchSize.
		if err := dstBatch.Put(key, val); err != nil {
			return errors.Wrap(err, "failed to put batch")
		}

		if dstBatch.ValueSize() > IdealBatchSize {
			if err := dstBatch.Write(); err != nil {
				return err
			}
			dstBatch.Reset()
		}

		// make a report
		if fetched%(IdealBatchSize*10) == 0 {
			logger.Info("DB migrated",
				"db", name, "fetched", fetched-previousFetched, "elapsedIter", time.Since(cycleStart),
				"fetchedTotal", fetched, "elapsedTotal", time.Since(start))
			cycleStart = time.Now()
			previousFetched = fetched
		}

		// check for quit signal from OS
		select {
		case <-sigQuit:
			logger.Warn("exit called", "db", name, "fetchedTotal", fetched, "elapsedTotal", time.Since(start))
			return errors.New("sigQuit")
		default:
		}
	}

	if err := dstBatch.Write(); err != nil {
		return errors.Wrap(err, "failed to write items")
	}
	dstBatch.Reset()

	logger.Info("Finish DB migration", "db", name, "fetchedTotal", fetched, "elapsedTotal", time.Since(start))

	srcIter.Release()
	if err := srcIter.Error(); err != nil { // any accumulated error from iterator
		return errors.Wrap(err, "failed to iterate")
	}

	return nil
}

// StartDBMigration migrates a DB to another DB.
// (e.g. LevelDB -> LevelDB, LevelDB -> BadgerDB, LevelDB -> DynamoDB)
//
// This feature uses Iterator. A src DB should have implementation of Iteratee to use this function.
// Do not use db migration while a node is executing.
func (dbm *databaseManager) StartDBMigration(dstdbm DBManager) error {
	// from non single DB
	if !dbm.config.SingleDB {
		for et := StateTrieDB; et < databaseEntryTypeSize; et++ {
			srcDB := dbm.getDatabase(et)

			dstDB := dstdbm.getDatabase(MiscDB)
			if !dstdbm.GetDBConfig().SingleDB {
				dstDB = dstdbm.getDatabase(et)
			}

			if srcDB == nil {
				logger.Warn("skip nil src db", "db", dbBaseDirs[et])
				continue
			}

			if dstDB == nil {
				logger.Warn("skip nil dst db", "db", dbBaseDirs[et])
				continue
			}

			if err := copyDB(dbBaseDirs[et], srcDB, dstDB); err != nil {
				return err
			}

			logger.Info("complete copy db", "db", dbBaseDirs[et])
		}

		return nil
	}

	// single DB -> single DB
	srcDB := dbm.getDatabase(MiscDB)
	dstDB := dstdbm.getDatabase(MiscDB)

	if err := copyDB("single", srcDB, dstDB); err != nil {
		return err
	}

	// If the current src DB is misc DB, clear all db dir on dst
	// TODO: If DB Migration supports non-single db, change the checking logic
	if path.Base(dbm.config.Dir) == dbBaseDirs[MiscDB] {
		for i := uint8(MiscDB); i < uint8(databaseEntryTypeSize); i++ {
			dstdbm.setDBDir(DBEntryType(i), "")
		}
	}

	return nil
}
