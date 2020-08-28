package database

import (
	"os"
	"os/signal"
	"syscall"

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

	// fetch keys and values
	keys, vals, fetched := iterateDB(srcIter, dbMigrationFetchNum)
	iterateNum := 0

loop:
	for fetched > 0 {
		// write fetched keys and values to DB
		for i := 0; i < fetched; i++ {
			err := dstBatch.Put(keys[i], vals[i])
			if err != nil {
				return errors.Wrap(err, "failed to put batch")
			}
		}

		// check for quit signal from OS
		select {
		case <-sigQuit:
			logger.Info("exit called")
			break loop
		default:
		}

		// fetch keys and values
		keys, vals, fetched = iterateDB(srcIter, dbMigrationFetchNum)
		iterateNum++
	}

	err := dstBatch.Write()
	if err != nil {
		return errors.Wrap(err, "failed to write items")
	}

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
