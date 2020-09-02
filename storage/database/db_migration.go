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
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const (
	dbMigrationFetchNum = 500
	reportCycle         = dbMigrationFetchNum * 100
)

// StartDBMigration migrates a DB to another DB.
// (e.g. LevelDB -> LevelDB, LevelDB -> BadgerDB, LevelDB -> DynamoDB)
//
// This feature uses Iterator. A src DB should have implementation of Iteratee to use this function.
// Do not use db migration while a node is executing.
func (dbm *databaseManager) StartDBMigration(dstdbm DBManager) error {
	// settings for quit signal from os
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// TODO enable for all dbs
	srcDB := dbm.getDatabase(MiscDB) // first DB
	dstDB := dstdbm.getDatabase(MiscDB)

	// create src iterator and dst batch
	srcIter := srcDB.NewIterator()
	dstBatch := dstDB.NewBatch()

	// vars for log
	start := time.Now()
	cycleNum, fetchedTotal, previousFetchedTotal, fetched := 0, 0, 0, 1

loop:
	for fetched > 0 {
		cycleStart := time.Now()
		for fetched = 0; fetched < dbMigrationFetchNum && srcIter.Next(); fetched++ {
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
		}
		fetchedTotal += fetched

		if fetchedTotal%reportCycle == 0 {
			logger.Info("DB migrated",
				"fetched", fetchedTotal-previousFetchedTotal, "elapsedIter", time.Since(cycleStart),
				"fetchedTotal", fetchedTotal, "elapsedTotal", time.Since(start))
			cycleStart = time.Now()
			previousFetchedTotal = fetchedTotal
		}

		// check for quit signal from OS
		select {
		case <-sigQuit:
			logger.Info("exit called", "iterNum", cycleNum, "fetchedTotal", fetchedTotal, "elapsedTotal", time.Since(start))
			break loop
		default:
		}

		cycleNum++
	}

	if err := dstBatch.Write(); err != nil {
		return errors.Wrap(err, "failed to write items")
	}

	logger.Info("Finish DB migration", "iterNum", cycleNum, "fetched", fetchedTotal, "elapsedTotal", time.Since(start))

	srcIter.Release()
	if err := srcIter.Error(); err != nil { // any accumulated error from iterator
		return errors.Wrap(err, "failed to iterate")
	}

	return nil
}
