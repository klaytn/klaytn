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

//go:build rocksdb
// +build rocksdb

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

func filterPolicyStrToNative(t string) *grocksdb.NativeFilterPolicy {
	switch t {
	case "bloom":
		return grocksdb.NewBloomFilter(defaultBitsPerKey)
	case "ribbon":
		return grocksdb.NewRibbonFilterPolicy(defaultBitsPerKey)
	default:
		return nil
	}
}

func compressionStrToType(t string) grocksdb.CompressionType {
	switch t {
	case "snappy":
		return grocksdb.SnappyCompression
	case "zlib":
		return grocksdb.ZLibCompression
	case "bz2":
		return grocksdb.Bz2Compression
	case "lz4":
		return grocksdb.LZ4Compression
	case "lz4hc":
		return grocksdb.LZ4HCCompression
	case "xpress":
		return grocksdb.XpressCompression
	case "zstd":
		return grocksdb.ZSTDCompression
	default:
		return grocksdb.NoCompression
	}
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

	// Ensure we have some minimal caching and file guarantees
	if config.CacheSize < minCacheSizeForRocksDB {
		logger.Warn("Cache size too small, increasing to minimum recommended", "oldCacheSize", config.CacheSize, "newCacheSize", minCacheSizeForRocksDB)
		config.CacheSize = minCacheSizeForRocksDB
	}
	if config.MaxOpenFiles < minOpenFilesForRocksDB {
		logger.Warn("Max open files too small, increasing to minimum recommended", "oldMaxOpenFiles", config.MaxOpenFiles, "newMaxOpenFiles", minOpenFilesForRocksDB)
		config.MaxOpenFiles = minOpenFilesForRocksDB
	}

	blockCacheSize := config.CacheSize / 2 * 1024 * 1024 // half of cacheSize in MiB
	bufferSize := config.CacheSize / 2 * 1024 * 1024     // half of cacheSize in MiB

	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(blockCacheSize))
	if cacheIndexAndFilter := config.CacheIndexAndFilter; cacheIndexAndFilter {
		bbto.SetCacheIndexAndFilterBlocks(cacheIndexAndFilter)
		bbto.SetPinL0FilterAndIndexBlocksInCache(cacheIndexAndFilter)
	}

	policy := filterPolicyStrToNative(config.FilterPolicy)
	if policy != nil {
		bbto.SetFilterPolicy(policy)
		bbto.SetOptimizeFiltersForMemory(true)
	}

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.SetWriteBufferSize(bufferSize)
	opts.SetDumpMallocStats(config.DumpMallocStat)
	opts.SetCompression(compressionStrToType(config.CompressionType))
	opts.SetBottommostCompression(compressionStrToType(config.BottommostCompressionType))
	opts.SetMaxOpenFiles(config.MaxOpenFiles)

	logger.Info("RocksDB configuration", "blockCacheSize", blockCacheSize, "bufferSize", bufferSize, "enableDumpMallocStat", config.DumpMallocStat, "compressionType", config.CompressionType, "bottommostCompressionType", config.BottommostCompressionType, "filterPolicy", config.FilterPolicy, "disableMetrics", config.DisableMetrics, "maxOpenFiles", config.MaxOpenFiles, "cacheIndexAndFilter", config.CacheIndexAndFilter)

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
	if !db.config.DisableMetrics {
		start := time.Now()
		defer db.putTimer.Update(time.Since(start))
	}
	return db.db.Put(db.wo, key, value)
}

func (db *rocksDB) Has(key []byte) (bool, error) {
	dat, err := db.db.GetBytes(db.ro, key)
	if dat == nil || err != nil {
		return false, err
	}

	return true, nil
}

func (db *rocksDB) Get(key []byte) ([]byte, error) {
	if !db.config.DisableMetrics {
		start := time.Now()
		defer db.getTimer.Update(time.Since(start))
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

func (db *rocksDB) TryCatchUpWithPrimary() error {
	return db.db.TryCatchUpWithPrimary()
}

type rdbIter struct {
	first  bool
	iter   *grocksdb.Iterator
	prefix []byte
	db     *rocksDB
}

// Next moves the iterator to the next key/value pair. It returns whether the
// iterator is exhausted.
func (i *rdbIter) Next() bool {
	if i.first {
		i.first = false
	} else {
		i.iter.Next()
	}
	return i.iter.ValidForPrefix(i.prefix)
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (i *rdbIter) Error() error {
	if i.first {
		return nil
	}
	return i.iter.Err()
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (i *rdbIter) Key() []byte {
	if i.first {
		return nil
	}
	key := i.iter.Key()
	defer key.Free()
	return key.Data()
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (i *rdbIter) Value() []byte {
	if i.first {
		return nil
	}
	val := i.iter.Value()
	defer val.Free()
	return val.Data()
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
	firstKey := append(prefix, start...)
	iter.Seek(firstKey)
	return &rdbIter{first: true, iter: iter, prefix: prefix, db: db}
}

func (db *rocksDB) Close() {
	close(db.quitCh)
	db.db.CancelAllBackgroundWork(true)
	db.db.Close()
	db.wo.Destroy()
	db.ro.Destroy()
	db.logger.Info("RocksDB is closed")
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
	if !metricutils.Enabled || db.config.DisableMetrics {
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
	if !b.db.config.DisableMetrics {
		start := time.Now()
		defer b.db.batchWriteTimer.Update(time.Since(start))
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
	b.db.logger.Crit("rocksdb batch does not implement Replay method")
	return nil
}
