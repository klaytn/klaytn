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
	elapsedFetch, elapsedPut, elapsedTotal := time.Now(), time.Now(), time.Now()
	var elapsedFetchTime time.Duration
	iterateNum, fetchedTotal := 0, 0

	// fetch keys and values
	keys, vals, fetched := iterateDB(srcIter, dbMigrationFetchNum)
	elapsedFetchTime = time.Since(elapsedFetch)
	fetchedTotal += fetched

loop:
	for fetched > 0 {
		// write fetched keys and values to DB
		elapsedPut = time.Now()
		for i := 0; i < fetched; i++ {
			err := dstBatch.Put(keys[i], vals[i])
			if err != nil {
				return errors.Wrap(err, "failed to put batch")
			}
		}
		logger.Info("Finish DB migration", "iterNum", iterateNum, "fetched", fetched,
			"elapsedFetch", elapsedFetchTime, "elapsedPut", time.Since(elapsedPut), "elapsedIter", time.Since(elapsedFetch),
			"fetchedTotal", fetchedTotal, "elapsedTotal", time.Since(elapsedTotal))

		// check for quit signal from OS
		select {
		case <-sigQuit:
			logger.Info("exit called", "iterNum", iterateNum, "fetched", fetchedTotal, "elapsedTotal", time.Since(elapsedTotal))
			break loop
		default:
		}

		// fetch keys and values
		elapsedFetch = time.Now()
		keys, vals, fetched = iterateDB(srcIter, dbMigrationFetchNum)
		elapsedFetchTime = time.Since(elapsedFetch)

		iterateNum++
		fetchedTotal += fetched
	}

	err := dstBatch.Write()
	if err != nil {
		return errors.Wrap(err, "failed to write items")
	}

	logger.Info("Finish DB migration", "iterNum", iterateNum, "fetched", fetchedTotal, "elapsedTotal", time.Since(elapsedTotal))

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
