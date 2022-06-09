// Copyright 2018 The klaytn Authors
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
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/klaytn/klaytn/log"
)

const (
	gcThreshold      = int64(1 << 30) // GB
	sizeGCTickerTime = 1 * time.Minute
)

type badgerDB struct {
	fn string // filename for reporting
	db *badger.DB

	gcTicker *time.Ticker  // runs periodically and runs gc if db size exceeds the threshold.
	closeCh  chan struct{} // stops gc go routine when db closes.

	logger log.Logger // Contextual logger tracking the database path
}

func getBadgerDBOptions(dbDir string) badger.Options {
	opts := badger.DefaultOptions(dbDir)
	return opts
}

func NewBadgerDB(dbDir string) (*badgerDB, error) {
	localLogger := logger.NewWith("dbDir", dbDir)

	if fi, err := os.Stat(dbDir); err == nil {
		if !fi.IsDir() {
			return nil, fmt.Errorf("failed to make badgerDB while checking dbDir. Given dbDir is not a directory. dbDir: %v", dbDir)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to make badgerDB while making dbDir. dbDir: %v, err: %v", dbDir, err)
		}
	} else {
		return nil, fmt.Errorf("failed to make badgerDB while checking dbDir. dbDir: %v, err: %v", dbDir, err)
	}

	opts := getBadgerDBOptions(dbDir)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to make badgerDB while opening the DB. dbDir: %v, err: %v", dbDir, err)
	}

	badger := &badgerDB{
		fn:       dbDir,
		db:       db,
		logger:   localLogger,
		gcTicker: time.NewTicker(sizeGCTickerTime),
		closeCh:  make(chan struct{}),
	}

	go badger.runValueLogGC()

	return badger, nil
}

// runValueLogGC runs gc for two cases.
// It periodically checks the size of value log and runs gc if it exceeds gcThreshold.
func (bg *badgerDB) runValueLogGC() {
	_, lastValueLogSize := bg.db.Size()

	for {
		select {
		case <-bg.closeCh:
			bg.logger.Debug("Stopped value log GC", "dbDir", bg.fn)
			return
		case <-bg.gcTicker.C:
			_, currValueLogSize := bg.db.Size()
			if currValueLogSize-lastValueLogSize < gcThreshold {
				continue
			}

			err := bg.db.RunValueLogGC(0.5)
			if err != nil {
				bg.logger.Error("Error while runValueLogGC()", "err", err)
				continue
			}

			_, lastValueLogSize = bg.db.Size()
		}
	}
}

func (bg *badgerDB) Type() DBType {
	return BadgerDB
}

// Path returns the path to the database directory.
func (bg *badgerDB) Path() string {
	return bg.fn
}

// Put inserts the given key and value pair to the database.
func (bg *badgerDB) Put(key []byte, value []byte) error {
	txn := bg.db.NewTransaction(true)
	defer txn.Discard()
	err := txn.Set(key, value)
	if err != nil {
		return err
	}
	return txn.Commit()
}

// Has returns true if the corresponding value to the given key exists.
func (bg *badgerDB) Has(key []byte) (bool, error) {
	txn := bg.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		return false, err
	}
	err = item.Value(nil)
	return err == nil, err
}

// Get returns the corresponding value to the given key if exists.
func (bg *badgerDB) Get(key []byte) ([]byte, error) {
	txn := bg.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, dataNotFoundErr
		}
		return nil, err
	}

	var valCopy []byte
	err = item.Value(func(val []byte) error {
		valCopy = append([]byte{}, val...)
		return nil
	})
	return valCopy, err
}

// Delete deletes the key from the queue and database
func (bg *badgerDB) Delete(key []byte) error {
	txn := bg.db.NewTransaction(true)
	defer txn.Discard()
	err := txn.Delete(key)
	if err != nil {
		return err
	}
	return txn.Commit()
}

func (bg *badgerDB) NewIterator(prefix []byte, start []byte) Iterator {
	// txn := bg.db.NewTransaction(false)
	// return txn.NewIterator(badger.DefaultIteratorOptions)
	logger.CritWithStack("badgerDB doesn't support NewIterator")
	return nil
}

func (bg *badgerDB) Close() {
	close(bg.closeCh)
	err := bg.db.Close()
	if err == nil {
		bg.logger.Info("Database closed")
	} else {
		bg.logger.Error("Failed to close database", "err", err)
	}
}

func (bg *badgerDB) LDB() *badger.DB {
	return bg.db
}

func (bg *badgerDB) NewBatch() Batch {
	txn := bg.db.NewTransaction(true)
	return &badgerBatch{db: bg.db, txn: txn}
}

func (bg *badgerDB) Meter(prefix string) {
	logger.Warn("badgerDB does not support metrics!")
}

type badgerBatch struct {
	db   *badger.DB
	txn  *badger.Txn
	size int
}

// Put inserts the given value into the batch for later committing.
func (b *badgerBatch) Put(key, value []byte) error {
	err := b.txn.Set(key, value)
	b.size += len(value)
	return err
}

// Delete inserts the a key removal into the batch for later committing.
func (b *badgerBatch) Delete(key []byte) error {
	if err := b.txn.Delete(key); err != nil {
		return err
	}
	b.size += 1
	return nil
}

// Write flushes any accumulated data to disk.
func (b *badgerBatch) Write() error {
	return b.txn.Commit()
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *badgerBatch) ValueSize() int {
	return b.size
}

// Replay replays the batch contents.
func (b *badgerBatch) Reset() {
	b.txn = b.db.NewTransaction(true)
	b.size = 0
}

// Replay replays the batch contents.
func (b *badgerBatch) Replay(w KeyValueWriter) error {
	logger.CritWithStack("Replay is not implemented in badgerBatch!")
	return nil
}
