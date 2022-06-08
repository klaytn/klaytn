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
	"fmt"
	"sync"
	"time"

	klaytnmetrics "github.com/klaytn/klaytn/metrics"

	"github.com/klaytn/klaytn/common/fdlimit"
	"github.com/klaytn/klaytn/log"
	metricutils "github.com/klaytn/klaytn/metrics/utils"
	"github.com/rcrowley/go-metrics"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
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
	minWriteBufferSize             = 2 * opt.MiB
	minBlockCacheCapacity          = 2 * minWriteBufferSize
	MinOpenFilesCacheCapacity      = 16
	minBitsPerKeyForFilter         = 10
	minFileDescriptorsForDBManager = 2048
	minFileDescriptorsForLevelDB   = 16
)

var defaultLevelDBOption = &opt.Options{
	WriteBuffer:            minWriteBufferSize,
	BlockCacheCapacity:     minBlockCacheCapacity,
	OpenFilesCacheCapacity: MinOpenFilesCacheCapacity,
	Filter:                 filter.NewBloomFilter(minBitsPerKeyForFilter),
	DisableBufferPool:      false,
	DisableSeeksCompaction: true,
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
	if limit < minFileDescriptorsForDBManager {
		raised, err := fdlimit.Raise(minFileDescriptorsForDBManager)
		if err != nil || raised < minFileDescriptorsForDBManager {
			logger.Crit("Raised number of file descriptor is below the minimum value",
				"currFileDescriptorsLimit", limit, "minFileDescriptorsForDBManager", minFileDescriptorsForDBManager)
		}
		limit = int(raised)
	}
	return limit / 2 // Leave half for networking and other stuff
}

type levelDB struct {
	fn string      // filename for reporting
	db *leveldb.DB // LevelDB instance

	writeDelayCountMeter    metrics.Meter // Meter for measuring the cumulative number of write delays
	writeDelayDurationMeter metrics.Meter // Meter for measuring the cumulative duration of write delays

	aliveSnapshotsMeter metrics.Meter // Meter for measuring the number of alive snapshots
	aliveIteratorsMeter metrics.Meter // Meter for measuring the number of alive iterators

	compTimer              klaytnmetrics.HybridTimer // Meter for measuring the total time spent in database compaction
	compReadMeter          metrics.Meter             // Meter for measuring the data read during compaction
	compWriteMeter         metrics.Meter             // Meter for measuring the data written during compaction
	diskReadMeter          metrics.Meter             // Meter for measuring the effective amount of data read
	diskWriteMeter         metrics.Meter             // Meter for measuring the effective amount of data written
	blockCacheGauge        metrics.Gauge             // Gauge for measuring the current size of block cache
	openedTablesCountMeter metrics.Meter
	memCompGauge           metrics.Gauge // Gauge for tracking the number of memory compaction
	level0CompGauge        metrics.Gauge // Gauge for tracking the number of table compaction in level0
	nonlevel0CompGauge     metrics.Gauge // Gauge for tracking the number of table compaction in non0 level
	seekCompGauge          metrics.Gauge // Gauge for tracking the number of table compaction caused by read opt

	levelSizesGauge     []metrics.Gauge
	levelTablesGauge    []metrics.Gauge
	levelReadGauge      []metrics.Gauge
	levelWriteGauge     []metrics.Gauge
	levelDurationsGauge []metrics.Gauge

	perfCheck       bool
	getTimer        klaytnmetrics.HybridTimer
	putTimer        klaytnmetrics.HybridTimer
	batchWriteTimer klaytnmetrics.HybridTimer

	quitLock sync.Mutex      // Mutex protecting the quit channel access
	quitChan chan chan error // Quit channel to stop the metrics collection before closing the database

	prefix string     // prefix used for metrics
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
		DisableSeeksCompaction:        true,
	}

	return newOption
}

func NewLevelDB(dbc *DBConfig, entryType DBEntryType) (*levelDB, error) {
	localLogger := logger.NewWith("path", dbc.Dir)

	// Ensure we have some minimal caching and file guarantees
	if dbc.LevelDBCacheSize < 16 {
		dbc.LevelDBCacheSize = 16
	}
	if dbc.OpenFilesLimit < minFileDescriptorsForLevelDB {
		dbc.OpenFilesLimit = minFileDescriptorsForLevelDB
	}

	ldbOpts := getLevelDBOptions(dbc)
	ldbOpts.Compression = getCompressionType(dbc.LevelDBCompression, entryType)

	localLogger.Info("LevelDB configurations",
		"levelDBCacheSize", (ldbOpts.WriteBuffer+ldbOpts.BlockCacheCapacity)/opt.MiB, "openFilesLimit", ldbOpts.OpenFilesCacheCapacity,
		"useBufferPool", !ldbOpts.DisableBufferPool, "usePerfCheck", dbc.EnableDBPerfMetrics, "compressionType", ldbOpts.Compression,
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
		fn:        dbc.Dir,
		db:        db,
		logger:    localLogger,
		perfCheck: dbc.EnableDBPerfMetrics,
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
	// value = rle.Compress(value)
	if db.perfCheck {
		start := time.Now()
		err := db.put(key, value)
		db.putTimer.Update(time.Since(start))
		return err
	}
	return db.put(key, value)
}

func (db *levelDB) put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *levelDB) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get returns the given key if it's present.
func (db *levelDB) Get(key []byte) ([]byte, error) {
	if db.perfCheck {
		start := time.Now()
		val, err := db.get(key)
		db.getTimer.Update(time.Since(start))
		return val, err
	}
	return db.get(key)
	// return rle.Decompress(dat)
}

func (db *levelDB) get(key []byte) ([]byte, error) {
	dat, err := db.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, dataNotFoundErr
		}
		return nil, err
	}
	return dat, nil
}

// Delete deletes the key from the queue and database
func (db *levelDB) Delete(key []byte) error {
	// Execute the actual operation
	return db.db.Delete(key, nil)
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *levelDB) NewIterator(prefix []byte, start []byte) Iterator {
	return db.db.NewIterator(bytesPrefixRange(prefix, start), nil)
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
	db.prefix = prefix

	// Initialize all the metrics collector at the requested prefix
	db.writeDelayCountMeter = metrics.NewRegisteredMeter(prefix+"writedelay/count", nil)
	db.writeDelayDurationMeter = metrics.NewRegisteredMeter(prefix+"writedelay/duration", nil)
	db.aliveSnapshotsMeter = metrics.NewRegisteredMeter(prefix+"snapshots", nil)
	db.aliveIteratorsMeter = metrics.NewRegisteredMeter(prefix+"iterators", nil)
	db.compTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"compaction/time", nil)
	db.compReadMeter = metrics.NewRegisteredMeter(prefix+"compaction/read", nil)
	db.compWriteMeter = metrics.NewRegisteredMeter(prefix+"compaction/write", nil)
	db.diskReadMeter = metrics.NewRegisteredMeter(prefix+"disk/read", nil)
	db.diskWriteMeter = metrics.NewRegisteredMeter(prefix+"disk/write", nil)
	db.blockCacheGauge = metrics.NewRegisteredGauge(prefix+"blockcache", nil)

	db.openedTablesCountMeter = metrics.NewRegisteredMeter(prefix+"opendedtables", nil)

	db.getTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"get/time", nil)
	db.putTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"put/time", nil)
	db.batchWriteTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"batchwrite/time", nil)

	db.memCompGauge = metrics.NewRegisteredGauge(prefix+"compact/memory", nil)
	db.level0CompGauge = metrics.NewRegisteredGauge(prefix+"compact/level0", nil)
	db.nonlevel0CompGauge = metrics.NewRegisteredGauge(prefix+"compact/nonlevel0", nil)
	db.seekCompGauge = metrics.NewRegisteredGauge(prefix+"compact/seek", nil)

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

	// Write delay related stats
	var prevWriteDelayCount int32
	var prevWriteDelayDuration time.Duration

	// Alive snapshots/iterators
	var prevAliveSnapshots, prevAliveIterators int32

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
		// Write delay related stats
		db.writeDelayCountMeter.Mark(int64(s.WriteDelayCount - prevWriteDelayCount))
		db.writeDelayDurationMeter.Mark(int64(s.WriteDelayDuration - prevWriteDelayDuration))
		prevWriteDelayCount, prevWriteDelayDuration = s.WriteDelayCount, s.WriteDelayDuration

		// Alive snapshots/iterators
		db.aliveSnapshotsMeter.Mark(int64(s.AliveSnapshots - prevAliveSnapshots))
		db.aliveIteratorsMeter.Mark(int64(s.AliveIterators - prevAliveIterators))
		prevAliveSnapshots, prevAliveIterators = s.AliveSnapshots, s.AliveIterators

		// Compaction related stats
		var currCompRead, currCompWrite int64
		var currCompTime time.Duration
		for i := 0; i < len(s.LevelDurations); i++ {
			currCompTime += s.LevelDurations[i]
			currCompRead += s.LevelRead[i]
			currCompWrite += s.LevelWrite[i]

			db.updateLevelStats(s, i)
		}
		db.compTimer.Update(currCompTime - prevCompTime)
		db.compReadMeter.Mark(currCompRead - prevCompRead)
		db.compWriteMeter.Mark(currCompWrite - prevCompWrite)
		prevCompTime, prevCompRead, prevCompWrite = currCompTime, currCompRead, currCompWrite

		// IO related stats
		currRead, currWrite := s.IORead, s.IOWrite
		db.diskReadMeter.Mark(int64(currRead - prevRead))
		db.diskWriteMeter.Mark(int64(currWrite - prevWrite))
		prevRead, prevWrite = currRead, currWrite

		// BlockCache/OpenedTables related stats
		db.blockCacheGauge.Update(int64(s.BlockCacheSize))
		db.openedTablesCountMeter.Mark(int64(s.OpenedTablesCount))

		// Compaction related stats
		db.memCompGauge.Update(int64(s.MemComp))
		db.level0CompGauge.Update(int64(s.Level0Comp))
		db.nonlevel0CompGauge.Update(int64(s.NonLevel0Comp))
		db.seekCompGauge.Update(int64(s.SeekComp))

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

// updateLevelStats collects level-wise stats.
func (db *levelDB) updateLevelStats(s *leveldb.DBStats, lv int) {
	// dynamically creates a new metrics for a new level
	if len(db.levelSizesGauge) <= lv {
		prefix := db.prefix + fmt.Sprintf("level%v/", lv)
		db.levelSizesGauge = append(db.levelSizesGauge, metrics.NewRegisteredGauge(prefix+"size", nil))
		db.levelTablesGauge = append(db.levelTablesGauge, metrics.NewRegisteredGauge(prefix+"tables", nil))
		db.levelReadGauge = append(db.levelReadGauge, metrics.NewRegisteredGauge(prefix+"read", nil))
		db.levelWriteGauge = append(db.levelWriteGauge, metrics.NewRegisteredGauge(prefix+"write", nil))
		db.levelDurationsGauge = append(db.levelDurationsGauge, metrics.NewRegisteredGauge(prefix+"duration", nil))
	}

	db.levelSizesGauge[lv].Update(s.LevelSizes[lv])
	db.levelTablesGauge[lv].Update(int64(s.LevelTablesCounts[lv]))
	db.levelReadGauge[lv].Update(s.LevelRead[lv])
	db.levelWriteGauge[lv].Update(s.LevelWrite[lv])
	db.levelDurationsGauge[lv].Update(int64(s.LevelDurations[lv]))
}

func (db *levelDB) NewBatch() Batch {
	return &ldbBatch{b: new(leveldb.Batch), ldb: db}
}

// ldbBatch is a write-only leveldb batch that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type ldbBatch struct {
	b    *leveldb.Batch
	ldb  *levelDB
	size int
}

// Put inserts the given value into the batch for later committing.
func (b *ldbBatch) Put(key, value []byte) error {
	b.b.Put(key, value)
	b.size += len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *ldbBatch) Delete(key []byte) error {
	b.b.Delete(key)
	b.size++
	return nil
}

// Write flushes any accumulated data to disk.
func (b *ldbBatch) Write() error {
	if b.ldb.perfCheck {
		start := time.Now()
		err := b.write()
		b.ldb.batchWriteTimer.Update(time.Since(start))
		return err
	}
	return b.write()
}

func (b *ldbBatch) write() error {
	return b.ldb.db.Write(b.b, nil)
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *ldbBatch) ValueSize() int {
	return b.size
}

// Reset resets the batch for reuse.
func (b *ldbBatch) Reset() {
	b.b.Reset()
	b.size = 0
}

// bytesPrefixRange returns key range that satisfy
// - the given prefix, and
// - the given seek position
func bytesPrefixRange(prefix, start []byte) *util.Range {
	r := util.BytesPrefix(prefix)
	r.Start = append(r.Start, start...)
	return r
}

// Replay replays the batch contents.
func (b *ldbBatch) Replay(w KeyValueWriter) error {
	return b.b.Replay(&replayer{writer: w})
}

// replayer is a small wrapper to implement the correct replay methods.
type replayer struct {
	writer  KeyValueWriter
	failure error
}

// Put inserts the given value into the key-value data store.
func (r *replayer) Put(key, value []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Put(key, value)
}

// Delete removes the key from the key-value data store.
func (r *replayer) Delete(key []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Delete(key)
}
