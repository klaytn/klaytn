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

const (
	defaultRocksDBCacheSize    = 2 // 2MB
	defaultBitsPerKey          = 10
	minCacheSizeForRocksDB     = 16
	defaultOpenFilesForRocksDB = 1024
	minOpenFilesForRocksDB     = 16
)

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

type RocksDBConfig struct {
	Secondary                 bool
	DumpMallocStat            bool
	DisableMetrics            bool
	CacheSize                 uint64
	CompressionType           string
	BottommostCompressionType string
	FilterPolicy              string
	MaxOpenFiles              int
	CacheIndexAndFilter       bool
}

func GetDefaultRocksDBConfig() *RocksDBConfig {
	return &RocksDBConfig{
		Secondary:                 false,
		CacheSize:                 defaultRocksDBCacheSize,
		DumpMallocStat:            false,
		CompressionType:           "lz4",
		BottommostCompressionType: "zstd",
		FilterPolicy:              "ribbon",
		DisableMetrics:            false,
		MaxOpenFiles:              defaultOpenFilesForRocksDB,
		CacheIndexAndFilter:       true,
	}
}
