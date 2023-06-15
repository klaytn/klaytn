// Copyright 2023 The klaytn Authors
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

	"github.com/linxGnu/grocksdb"
)

const defaultRocksDBCacheSize = 2 // 2MB

type rocksDB struct {
	config *RocksDBConfig
	db     *grocksdb.DB // rocksDB instance

	wo *grocksdb.WriteOptions
	ro *grocksdb.ReadOptions
}

type RocksDBConfig struct {
	Secondary bool
	CacheSize uint64
}

func GetDefaultRocksDBConfig() *RocksDBConfig {
	return &RocksDBConfig{
		Secondary:      false,
		CacheSize:      defaultRocksDBCacheSize,
		DumpMallocStat: false,
	}
}

// openFile checks if the path is valid directory or not. If not exists, the path directory is created.
func openFile(path string, needToMake bool) error {
	if fi, err := os.Stat(path); err == nil {
		if !fi.IsDir() {
			return fmt.Errorf("rocksdb: open %s: not a directory", path)
		}
	} else if os.IsNotExist(err) && needToMake {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func NewRocksDB(path string, config *RocksDBConfig) (*rocksDB, error) {
	if err := openFile(path, !config.Secondary); err != nil {
		return nil, err
	}

	blockCacheSize := config.CacheSize / 2 * 1024 * 1024 // half of cacheSize in MiB
	bufferSize := config.CacheSize / 2 * 1024 * 1024     // half of cacheSize in MiB

	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(blockCacheSize))

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.SetWriteBufferSize(bufferSize)

	logger.Info("RocksDB configuration", "blockCacheSize", blockCacheSize, "bufferSize", bufferSize)

	var (
		db  *grocksdb.DB
		err error
	)

	if config.Secondary {
		db, err = grocksdb.OpenDbAsSecondary(opts, path, path)
	} else {
		db, err = grocksdb.OpenDb(opts, path)
	}
	if err != nil {
		return nil, err
	}
	return &rocksDB{config: config, db: db, wo: grocksdb.NewDefaultWriteOptions(), ro: grocksdb.NewDefaultReadOptions()}, nil
}

func (db *rocksDB) Type() DBType {
	return RocksDB
}

func (db *rocksDB) Put(key []byte, value []byte) error {
	if db.config.Secondary {
		return nil
	}
	return db.db.Put(db.wo, key, value)
}

func (db *rocksDB) Has(key []byte) (bool, error) {
	if db.config.Secondary {
		db.db.TryCatchUpWithPrimary()
	}

	dat, err := db.db.GetBytes(db.ro, key)
	if dat == nil || err != nil {
		return false, err
	}

	return true, nil
}

func (db *rocksDB) Get(key []byte) ([]byte, error) {
	if db.config.Secondary {
		db.db.TryCatchUpWithPrimary()
	}
	return db.get(key)
}

func (db *rocksDB) get(key []byte) ([]byte, error) {
	dat, err := db.db.GetBytes(db.ro, key)
	if dat == nil {
		return nil, dataNotFoundErr
	}

	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (db *rocksDB) Delete(key []byte) error {
	if db.config.Secondary {
		return nil
	}
	return db.db.Delete(db.wo, key)
}

type rdbIter struct {
	initialized bool
	iter        *grocksdb.Iterator
	prefix      []byte
	db          *rocksDB
}

// Next moves the iterator to the next key/value pair. It returns whether the
// iterator is exhausted.
func (i *rdbIter) Next() bool {
	if !i.initialized {
		i.initialized = i.iter.ValidForPrefix(i.prefix)
		return i.initialized
	}
	i.iter.Next()
	return i.iter.ValidForPrefix(i.prefix)
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (i *rdbIter) Error() error {
	if !i.initialized {
		return nil
	}
	return i.iter.Err()
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (i *rdbIter) Key() []byte {
	if !i.initialized {
		return []byte{}
	}
	return i.iter.Key().Data()
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (i *rdbIter) Value() []byte {
	if !i.initialized {
		return []byte{}
	}
	return i.iter.Value().Data()
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (i *rdbIter) Release() {
	i.iter.Close()
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *rocksDB) NewIterator(prefix []byte, start []byte) Iterator {
	iter := db.db.NewIterator(db.ro)
	if len(start) > 0 {
		iter.Seek(start)
	}
	return &rdbIter{initialized: false, iter: iter, prefix: prefix, db: db}
}

func (db *rocksDB) Close() {
	db.db.CancelAllBackgroundWork(true)
	db.db.Close()
	db.wo.Destroy()
	db.ro.Destroy()
	logger.Info("Rocksdb is closed")
}

// Meter configures the database metrics collectors and
func (db *rocksDB) Meter(prefix string) {
	// TODO-Klaytn-RocksDB implement metrics for rocksdb
}

func (db *rocksDB) NewBatch() Batch {
	return &rdbBatch{b: grocksdb.NewWriteBatch(), db: db}
}

// rdbBatch is a write-only rocksdb batch that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type rdbBatch struct {
	b    *grocksdb.WriteBatch
	db   *rocksDB
	size int
}

// Put inserts the given value into the batch for later committing.
func (b *rdbBatch) Put(key, value []byte) error {
	if b.db.config.Secondary {
		return nil
	}
	b.b.Put(key, value)
	b.size += len(value)
	return nil
}

// Delete inserts a key removal into the batch for later committing.
func (b *rdbBatch) Delete(key []byte) error {
	if b.db.config.Secondary {
		return nil
	}
	b.b.Delete(key)
	b.size++
	return nil
}

// Write flushes any accumulated data to disk.
func (b *rdbBatch) Write() error {
	if b.db.config.Secondary {
		return nil
	}
	return b.write()
}

func (b *rdbBatch) write() error {
	return b.db.db.Write(b.db.wo, b.b)
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *rdbBatch) ValueSize() int {
	return b.size
}

// Reset resets the batch for reuse.
func (b *rdbBatch) Reset() {
	b.b.Clear()
	b.size = 0
}

// Release free memory allocated to rocksdb batch object.
func (b *rdbBatch) Release() {
	b.b.Destroy()
}

// Replay replays the batch contents.
func (b *rdbBatch) Replay(w KeyValueWriter) error {
	logger.Crit("rocksdb batch does not implement Replay method")
	return nil
}
