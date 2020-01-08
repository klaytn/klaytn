// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethdb/database.go (2018/06/04).
// Modified and improved for the klaytn development.

package database

import (
	"github.com/klaytn/klaytn/common/fdlimit"
	metricutils "github.com/klaytn/klaytn/metrics/utils"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/klaytn/klaytn/log"
	"github.com/rcrowley/go-metrics"
)

var OpenFileLimit = 64

type LevelDBCompressionType uint8

const (
	AllNoCompression LevelDBCompressionType = iota
	ReceiptOnlySnappyCompression
	StateTrieOnlyNoCompression
	AllSnappyCompression
)

const (
	minWriteBufferSize        = 2 * opt.MiB
	minBlockCacheCapacity     = 2 * minWriteBufferSize
	MinOpenFilesCacheCapacity = 16
	minBitsPerKeyForFilter    = 10
)

var defaultLevelDBOption = &opt.Options{
	WriteBuffer:            minWriteBufferSize,
	BlockCacheCapacity:     minBlockCacheCapacity,
	OpenFilesCacheCapacity: MinOpenFilesCacheCapacity,
	Filter:                 filter.NewBloomFilter(minBitsPerKeyForFilter),
	DisableBufferPool:      false,
}

// GetDefaultLevelDBOption returns default LevelDB option copied from defaultLevelDBOption.
// defaultLevelDBOption has fields with minimum values.
func GetDefaultLevelDBOption() *opt.Options {
	copiedOption := *defaultLevelDBOption
	return &copiedOption
}

// GetOpenFilesLimit raises out the number of allowed file handles per process
// for Klaytn and returns half of the allowance to assign to the database.
func GetOpenFilesLimit() int {
	limit, err := fdlimit.Current()
	if err != nil {
		logger.Crit("Failed to retrieve file descriptor allowance", "err", err)
	}
	if limit < 2048 {
		if err := fdlimit.Raise(2048); err != nil {
			logger.Crit("Failed to raise file descriptor allowance", "err", err)
		}
	}
	if limit > 2048 { // cap database file descriptors even if more is available
		limit = 2048
	}
	return limit / 2 // Leave half for networking and other stuff
}

type levelDB struct {
	fn string      // filename for reporting
	db *leveldb.DB // LevelDB instance

	compTimeMeter   metrics.Meter // Meter for measuring the total time spent in database compaction
	compReadMeter   metrics.Meter // Meter for measuring the data read during compaction
	compWriteMeter  metrics.Meter // Meter for measuring the data written during compaction
	diskReadMeter   metrics.Meter // Meter for measuring the effective amount of data read
	diskWriteMeter  metrics.Meter // Meter for measuring the effective amount of data written
	blockCacheGauge metrics.Gauge // Gauge for measuring the current size of block cache

	quitLock sync.Mutex      // Mutex protecting the quit channel access
	quitChan chan chan error // Quit channel to stop the metrics collection before closing the database

	logger log.Logger // Contextual logger tracking the database path
}

func getLevelDBOptions(dbc *DBConfig) *opt.Options {
	newOption := &opt.Options{
		OpenFilesCacheCapacity:        dbc.OpenFilesLimit,
		BlockCacheCapacity:            dbc.LevelDBCacheSize / 2 * opt.MiB,
		WriteBuffer:                   dbc.LevelDBCacheSize / 2 * opt.MiB,
		Filter:                        filter.NewBloomFilter(10),
		DisableBufferPool:             !dbc.LevelDBBufferPool,
		CompactionTableSize:           2 * opt.MiB,
		CompactionTableSizeMultiplier: 1.0,
	}

	return newOption
}

func NewLevelDB(dbc *DBConfig, entryType DBEntryType) (*levelDB, error) {
	localLogger := logger.NewWith("path", dbc.Dir)

	// Ensure we have some minimal caching and file guarantees
	if dbc.LevelDBCacheSize < 16 {
		dbc.LevelDBCacheSize = 16
	}
	if dbc.OpenFilesLimit < 16 {
		dbc.OpenFilesLimit = 16
	}

	ldbOpts := getLevelDBOptions(dbc)
	ldbOpts.Compression = getCompressionType(dbc.LevelDBCompression, entryType)

	localLogger.Info("LevelDB configurations",
		"levelDBCacheSize", (ldbOpts.WriteBuffer+ldbOpts.BlockCacheCapacity)/opt.MiB, "openFilesLimit", ldbOpts.OpenFilesCacheCapacity,
		"useBufferPool", !ldbOpts.DisableBufferPool, "compressionType", ldbOpts.Compression,
		"compactionTableSize(MB)", ldbOpts.CompactionTableSize/opt.MiB, "compactionTableSizeMultiplier", ldbOpts.CompactionTableSizeMultiplier)

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(dbc.Dir, ldbOpts)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(dbc.Dir, nil)
	}
	// (Re)check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	return &levelDB{
		fn:     dbc.Dir,
		db:     db,
		logger: localLogger,
	}, nil
}

// setMinLevelDBOption sets some value of options if they are smaller than minimum value.
func setMinLevelDBOption(ldbOption *opt.Options) {
	if ldbOption.WriteBuffer < minWriteBufferSize {
		ldbOption.WriteBuffer = minWriteBufferSize
	}

	if ldbOption.BlockCacheCapacity < minBlockCacheCapacity {
		ldbOption.BlockCacheCapacity = minBlockCacheCapacity
	}

	if ldbOption.OpenFilesCacheCapacity < MinOpenFilesCacheCapacity {
		ldbOption.OpenFilesCacheCapacity = MinOpenFilesCacheCapacity
	}
}

func getCompressionType(ct LevelDBCompressionType, dbEntryType DBEntryType) opt.Compression {
	if ct == AllSnappyCompression {
		return opt.SnappyCompression
	}

	if ct == AllNoCompression {
		return opt.NoCompression
	}

	if ct == ReceiptOnlySnappyCompression {
		if dbEntryType == ReceiptsDB {
			return opt.SnappyCompression
		} else {
			return opt.NoCompression
		}
	}

	if ct == StateTrieOnlyNoCompression {
		if dbEntryType == StateTrieDB {
			return opt.NoCompression
		} else {
			return opt.SnappyCompression
		}
	}
	return opt.NoCompression
}

// NewLevelDBWithOption explicitly receives LevelDB option to construct a LevelDB object.
func NewLevelDBWithOption(dbPath string, ldbOption *opt.Options) (*levelDB, error) {
	// TODO-Klaytn-Database Replace `NewLevelDB` with `NewLevelDBWithOption`

	localLogger := logger.NewWith("path", dbPath)

	setMinLevelDBOption(ldbOption)

	localLogger.Info("Allocated LevelDB",
		"WriteBuffer (MB)", ldbOption.WriteBuffer/opt.MiB, "OpenFilesCacheCapacity", ldbOption.OpenFilesCacheCapacity, "BlockCacheCapacity (MB)", ldbOption.BlockCacheCapacity/opt.MiB,
		"CompactionTableSize (MB)", ldbOption.CompactionTableSize/opt.MiB, "CompactionTableSizeMultiplier", ldbOption.CompactionTableSizeMultiplier, "DisableBufferPool", ldbOption.DisableBufferPool)

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(dbPath, ldbOption)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(dbPath, nil)
	}
	// (Re)check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	return &levelDB{
		fn:     dbPath,
		db:     db,
		logger: localLogger,
	}, nil

}

func (db *levelDB) Type() DBType {
	return LevelDB
}

// Path returns the path to the database directory.
func (db *levelDB) Path() string {
	return db.fn
}

// Put puts the given key / value to the queue
func (db *levelDB) Put(key []byte, value []byte) error {
	// Generate the data to write to disk, update the meter and write
	//value = rle.Compress(value)

	return db.db.Put(key, value, nil)
}

func (db *levelDB) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get returns the given key if it's present.
func (db *levelDB) Get(key []byte) ([]byte, error) {
	// Retrieve the key and increment the miss counter if not found
	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
	//return rle.Decompress(dat)
}

// Delete deletes the key from the queue and database
func (db *levelDB) Delete(key []byte) error {
	// Execute the actual operation
	return db.db.Delete(key, nil)
}

func (db *levelDB) NewIterator() iterator.Iterator {
	return db.db.NewIterator(nil, nil)
}

// NewIteratorWithPrefix returns a iterator to iterate over subset of database content with a particular prefix.
func (db *levelDB) NewIteratorWithPrefix(prefix []byte) iterator.Iterator {
	return db.db.NewIterator(util.BytesPrefix(prefix), nil)
}

func (db *levelDB) Close() {
	// Stop the metrics collection to avoid internal database races
	db.quitLock.Lock()
	defer db.quitLock.Unlock()

	if db.quitChan != nil {
		errc := make(chan error)
		db.quitChan <- errc
		if err := <-errc; err != nil {
			db.logger.Error("Metrics collection failed", "err", err)
		}
		db.quitChan = nil
	}
	err := db.db.Close()
	if err == nil {
		db.logger.Info("Database closed")
	} else {
		db.logger.Error("Failed to close database", "err", err)
	}
}

func (db *levelDB) LDB() *leveldb.DB {
	return db.db
}

// Meter configures the database metrics collectors and
func (db *levelDB) Meter(prefix string) {
	// Initialize all the metrics collector at the requested prefix
	db.compTimeMeter = metrics.NewRegisteredMeter(prefix+"compaction/time", nil)
	db.compReadMeter = metrics.NewRegisteredMeter(prefix+"compaction/read", nil)
	db.compWriteMeter = metrics.NewRegisteredMeter(prefix+"compaction/write", nil)
	db.diskReadMeter = metrics.NewRegisteredMeter(prefix+"disk/read", nil)
	db.diskWriteMeter = metrics.NewRegisteredMeter(prefix+"disk/write", nil)
	db.blockCacheGauge = metrics.NewRegisteredGauge(prefix+"blockcache", nil)

	// Short circuit metering if the metrics system is disabled
	// Above meters are initialized by NilMeter if metricutils.Enabled == false
	if !metricutils.Enabled {
		return
	}

	// Create a quit channel for the periodic collector and run it
	db.quitLock.Lock()
	db.quitChan = make(chan chan error)
	db.quitLock.Unlock()

	go db.meter(3 * time.Second)
}

// meter periodically retrieves internal leveldb counters and reports them to
// the metrics subsystem.
//
// This is how a stats table look like (currently):
//   Compactions
//    Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
//   -------+------------+---------------+---------------+---------------+---------------
//      0   |          0 |       0.00000 |       1.27969 |       0.00000 |      12.31098
//      1   |         85 |     109.27913 |      28.09293 |     213.92493 |     214.26294
//      2   |        523 |    1000.37159 |       7.26059 |      66.86342 |      66.77884
//      3   |        570 |    1113.18458 |       0.00000 |       0.00000 |       0.00000
//
// This is how the iostats look like (currently):
// Read(MB):3895.04860 Write(MB):3654.64712
func (db *levelDB) meter(refresh time.Duration) {
	s := new(leveldb.DBStats)

	// Compaction related stats
	var prevCompRead, prevCompWrite int64
	var prevCompTime time.Duration

	// IO related stats
	var prevRead, prevWrite uint64

	var (
		errc chan error
		merr error
	)

	// Keep collecting stats unless an error occurs
hasError:
	for {
		merr = db.db.Stats(s)
		if merr != nil {
			break
		}

		// Compaction related stats
		var currCompRead, currCompWrite int64
		var currCompTime time.Duration
		for i := 0; i < len(s.LevelDurations); i++ {
			currCompTime += s.LevelDurations[i]
			currCompRead += s.LevelRead[i]
			currCompWrite += s.LevelWrite[i]
		}

		db.compTimeMeter.Mark(int64(currCompTime.Seconds() - prevCompTime.Seconds()))
		db.compReadMeter.Mark(int64(currCompRead - prevCompRead))
		db.compWriteMeter.Mark(int64(currCompWrite - prevCompWrite))

		prevCompTime = currCompTime
		prevCompRead = currCompRead
		prevCompWrite = currCompWrite

		// IO related stats
		currRead, currWrite := s.IORead, s.IOWrite

		db.diskReadMeter.Mark(int64(currRead - prevRead))
		db.diskWriteMeter.Mark(int64(currWrite - prevWrite))

		prevRead, prevWrite = currRead, currWrite

		// BlockCache size
		db.blockCacheGauge.Update(int64(s.BlockCacheSize))

		// Sleep a bit, then repeat the stats collection
		select {
		case errc = <-db.quitChan:
			// Quit requesting, stop hammering the database
			break hasError
		case <-time.After(refresh):
			// Timeout, gather a new set of stats
		}
	}

	if errc == nil {
		errc = <-db.quitChan
	}
	errc <- merr
}

func (db *levelDB) NewBatch() Batch {
	return &ldbBatch{db: db.db, b: new(leveldb.Batch)}
}

type ldbBatch struct {
	db   *leveldb.DB
	b    *leveldb.Batch
	size int
}

func (b *ldbBatch) Put(key, value []byte) error {
	b.b.Put(key, value)
	b.size += len(value)
	return nil
}

func (b *ldbBatch) Write() error {
	return b.db.Write(b.b, nil)
}

func (b *ldbBatch) ValueSize() int {
	return b.size
}

func (b *ldbBatch) Reset() {
	b.b.Reset()
	b.size = 0
}
