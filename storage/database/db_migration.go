// Modifications Copyright 2020 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
)

// StartDBMigration migrates a DB to a different kind of DB.
// (e.g. LevelDB -> BadgerDB, LevelDB -> DynamoDB)
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
	srcDB := dbm.dbs[0]
	dstDB := dstdbm.getDatabase(DBEntryType(MiscDB))

	// create src iterator and dst batch
	srcIter := srcDB.NewIterator()
	dstBatch := dstDB.NewBatch()

	// vars for log
	start := time.Now()
	iterateNum, fetchedTotal, fetched := 0, 0, 1

loop:
	for fetched > 0 {
		// fetch keys and values
		var keys, vals [][]byte
		iterStart := time.Now()
		keys, vals, fetched = iterateDB(srcIter, dbMigrationFetchNum)
		elapsedFetchTime := time.Since(iterStart)
		fetchedTotal += fetched

		// write fetched keys and values to DB
		putStart := time.Now()
		for i := 0; i < fetched; i++ {
			err := dstBatch.Put(keys[i], vals[i])
			if err != nil {
				return errors.Wrap(err, "failed to put batch")
			}
		}
		logger.Info("DB migrated", "iterNum", iterateNum, "fetched", fetched,
			"elapsedFetch", elapsedFetchTime, "elapsedPut", time.Since(putStart), "elapsedIter", time.Since(iterStart),
			"fetchedTotal", fetchedTotal, "elapsedTotal", time.Since(start))

		// check for quit signal from OS
		select {
		case <-sigQuit:
			logger.Info("exit called", "iterNum", iterateNum, "fetched", fetchedTotal, "elapsedTotal", time.Since(start))
			break loop
		default:
		}

		iterateNum++
	}

	err := dstBatch.Write()
	if err != nil {
		return errors.Wrap(err, "failed to write items")
	}

	logger.Info("Finish DB migration", "iterNum", iterateNum, "fetched", fetchedTotal, "elapsedTotal", time.Since(start))

	srcIter.Release()
	err = srcIter.Error()
	if err != nil {
		return errors.Wrap(err, "failed to iterate")
	}

	return nil
}

func iterateDB(iter Iterator, num int) ([][]byte, [][]byte, int) {
	keys := make([][]byte, num)
	vals := make([][]byte, num)

	var i int
	for i = 0; i < num && iter.Next(); i++ {
		// Contents of iter.Key() and iter.Value() should not be modified, and
		// only valid until the next call to Next.
		keys[i] = make([]byte, len(iter.Key()))
		vals[i] = make([]byte, len(iter.Value()))
		copy(keys[i], iter.Key())
		copy(vals[i], iter.Value())
	}

	return keys, vals, i
}
