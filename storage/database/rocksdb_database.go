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
	"strings"
	"time"

	"github.com/klaytn/klaytn/log"
	klaytnmetrics "github.com/klaytn/klaytn/metrics"
	metricutils "github.com/klaytn/klaytn/metrics/utils"
	"github.com/linxGnu/grocksdb"
	"github.com/rcrowley/go-metrics"
)

const defaultRocksDBCacheSize = 2 // 2MB

var properties = []string{
	"rocksdb.num-immutable-mem-table",         // returns number of immutable memtables that have not yet been flushed.
	"rocksdb.mem-table-flush-pending",         // returns 1 if a memtable flush is pending; otherwise, returns 0.
	"rocksdb.compaction-pending",              // returns 1 if at least one compaction is pending; otherwise, returns 0
	"rocksdb.background-errors",               // returns accumulated number of background errors.
	"rocksdb.cur-size-active-mem-table",       // returns approximate size of active memtable (bytes).
	"rocksdb.cur-size-all-mem-tables",         // returns approximate size of active and unflushed immutable memtables (bytes)
	"rocksdb.size-all-mem-tables",             // returns approximate size of active, unflushed immutable, and pinned immutable memtables (bytes)
	"rocksdb.num-entries-active-mem-table",    // returns total number of entries in the active memtable.
	"rocksdb.num-entries-imm-mem-tables",      // returns total number of entries in the unflushed immutable memtables.
	"rocksdb.num-deletes-active-mem-table",    // returns total number of delete entries in the active memtable.
	"rocksdb.num-deletes-imm-mem-tables",      // returns total number of delete entries in the unflushed immutable memtables.
	"rocksdb.estimate-num-keys",               // returns estimated number of total keys in the active and unflushed immutable memtables and storage.
	"rocksdb.estimate-table-readers-mem",      // returns estimated memory used for reading SST tables, excluding memory used in block cache (e.g.filter and index blocks).
	"rocksdb.is-file-deletions-enabled",       // returns 0 if deletion of obsolete files is enabled; otherwise, returns a non-zero number. This name may be misleading because true(non-zero) means disable, but we keep the name for backward compatibility.
	"rocksdb.num-snapshots",                   // returns number of unreleased snapshots of the database.
	"rocksdb.oldest-snapshot-time",            // returns number representing unix timestamp of oldest unreleased snapshot.
	"rocksdb.num-live-versions",               // returns number of live versions. `Version` is an internal data structure. See version_set.h for details. More live versions often mean more SST files are held from being deleted, by iterators or unfinished compactions.
	"rocksdb.current-super-version-number",    // returns number of current LSM version. It is a uint64_t integer number, incremented after there is any change to the LSM tree. The number is not preserved after restarting the DB. After DB restart, it will start from 0 again.
	"rocksdb.estimate-live-data-size",         // returns an estimate of the amount of live data in bytes. For BlobDB, it also includes the exact value of live bytes in the blob files of the version.
	"rocksdb.min-log-number-to-keep",          // return the minimum log number of the log files that should be kept.
	"rocksdb.min-obsolete-sst-number-to-keep", // return the minimum file number for an obsolete SST to be kept. The max value of `uint64_t` will be returned if all obsolete files can be deleted.
	//"rocksdb.total-sst-files-size", // returns total size (bytes) of all SST files belonging to any of the CF's versions. WARNING: may slow down online queries if there are too many files.
	"rocksdb.live-sst-files-size",               // returns total size (bytes) of all SST files belong to the CF's current version.
	"rocksdb.obsolete-sst-files-size",           // returns total size (bytes) of all SST files that became obsolete but have not yet been deleted or scheduled for deletion. SST files can end up in this state when using `DisableFileDeletions()`, for example. N.B. Unlike the other "*SstFilesSize" properties, this property includes SST files that originated in any of the DB's CFs.
	"rocksdb.base-level",                        // returns number of level to which L0 data will be compacted.
	"rocksdb.estimate-pending-compaction-bytes", // returns estimated total number of bytes compaction needs to rewrite to get all levels down to under target size. Not valid for other compactions than level-based.
	"rocksdb.num-running-compactions",           // returns the number of currently running compactions.
	"rocksdb.num-running-flushes",               // returns the number of currently running flushes.
	"rocksdb.actual-delayed-write-rate",         // returns the current actual delayed write rate. 0 means no delay.
	"rocksdb.is-write-stopped",                  // return 1 if write has been stopped.
	"rocksdb.estimate-oldest-key-time",          // returns an estimation of oldest key timestamp in the DB. Currently only available for FIFO compaction with compaction_options_fifo.allow_compaction = false.

	// Properties dedicated for BlobDB
	"rocksdb.num-blob-files",          // returns number of blob files in the current version.
	"rocksdb.total-blob-file-size",    // returns the total size of all blob files over all versions.
	"rocksdb.live-blob-file-size",     // returns the total size of all blob files in the current version.
	"rocksdb.blob-cache-capacity",     // returns blob cache capacity.
	"rocksdb.blob-cache-usage",        // returns the memory size for the entries residing in blob cache.
	"rocksdb.blob-cache-pinned-usage", // returns the memory size for the entries being pinned in blob cache.
}

type rocksDB struct {
	config *RocksDBConfig
	db     *grocksdb.DB // rocksDB instance

	wo *grocksdb.WriteOptions
	ro *grocksdb.ReadOptions

	quitCh          chan struct{}
	metrics         []metrics.Meter
	getTimer        klaytnmetrics.HybridTimer
	putTimer        klaytnmetrics.HybridTimer
	batchWriteTimer klaytnmetrics.HybridTimer

	prefix string
	logger log.Logger
}

type RocksDBConfig struct {
	Secondary      bool
	CacheSize      uint64
	DumpMallocStat bool
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
	localLogger := logger.NewWith("path", path)

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
	opts.SetDumpMallocStats(config.DumpMallocStat)

	logger.Info("RocksDB configuration", "blockCacheSize", blockCacheSize, "bufferSize", bufferSize, "enableDumpMallocStat", config.DumpMallocStat)

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
	return &rocksDB{
		config: config,
		db:     db,
		wo:     grocksdb.NewDefaultWriteOptions(),
		ro:     grocksdb.NewDefaultReadOptions(),
		logger: localLogger,
		quitCh: make(chan struct{}),
	}, nil
}

func (db *rocksDB) Type() DBType {
	return RocksDB
}

func (db *rocksDB) Put(key []byte, value []byte) error {
	if db.config.Secondary {
		return nil
	}
	start := time.Now()
	defer db.putTimer.Update(time.Since(start))
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
	start := time.Now()
	defer db.getTimer.Update(time.Since(start))
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

func (db *rocksDB) GetProperty(name string) string {
	return db.db.GetProperty(name)
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
	close(db.quitCh)
	db.db.CancelAllBackgroundWork(true)
	db.db.Close()
	db.wo.Destroy()
	db.ro.Destroy()
	db.logger.Info("Rocksdb is closed")
}

func (db *rocksDB) updateMeter(name string, meter metrics.Meter) {
	v, s := db.db.GetIntProperty(name)
	if s {
		meter.Mark(int64(v))
	}
}

// Meter configures the database metrics collectors and
func (db *rocksDB) Meter(prefix string) {
	db.prefix = prefix

	for _, property := range properties {
		splited := strings.Split(property, ".")
		name := strings.ReplaceAll(splited[1], "-", "/")
		db.metrics = append(db.metrics, metrics.NewRegisteredMeter(prefix+name, nil))
	}
	db.getTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"get/time", nil)
	db.putTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"put/time", nil)
	db.batchWriteTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"batchwrite/time", nil)

	// Short circuit metering if the metrics system is disabled
	// Above meters are initialized by NilMeter if metricutils.Enabled == false
	if !metricutils.Enabled {
		return
	}

	go db.meter(3 * time.Second)
}

func (db *rocksDB) meter(t time.Duration) {
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		select {
		case <-db.quitCh:
			return
		case <-ticker.C:
			for idx, property := range properties {
				db.updateMeter(property, db.metrics[idx])
			}
		}
	}
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
	start := time.Now()
	defer b.db.batchWriteTimer.Update(time.Since(start))
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
	b.db.logger.Crit("rocksdb batch does not implement Replay method")
	return nil
}
