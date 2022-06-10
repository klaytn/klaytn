// Copyright 2019 The klaytn Authors
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
	"bytes"
	"container/heap"
	"context"
	"fmt"
	"path"
	"strconv"
	"sync"

	"github.com/klaytn/klaytn/common"
)

var errKeyLengthZero = fmt.Errorf("database key for sharded database should be greater than 0")

const numShardsLimit = 256

type shardedDB struct {
	fn        string
	shards    []Database
	numShards uint

	sdbBatchTaskCh chan sdbBatchTask
}

type sdbBatchTask struct {
	batch    Batch               // A batch that each worker executes.
	index    int                 // Index of given batch.
	resultCh chan sdbBatchResult // Batch result channel for each shardedDBBatch.
}

type sdbBatchResult struct {
	index int   // Index of the batch result.
	err   error // Error from the batch write operation.
}

// newShardedDB creates database with numShards shards, or partitions.
// The type of database is specified DBConfig.DBType.
func newShardedDB(dbc *DBConfig, et DBEntryType, numShards uint) (*shardedDB, error) {
	if numShards == 0 {
		logger.Crit("numShards should be greater than 0!")
	}

	if numShards > numShardsLimit {
		logger.Crit(fmt.Sprintf("numShards should be equal to or smaller than %v, but it is %v.", numShardsLimit, numShards))
	}

	if !IsPow2(numShards) {
		logger.Crit(fmt.Sprintf("numShards should be power of two, but it is %v", numShards))
	}

	shards := make([]Database, 0, numShards)
	sdbBatchTaskCh := make(chan sdbBatchTask, numShards*2)
	sdbLevelDBCacheSize := dbc.LevelDBCacheSize / int(numShards)
	sdbOpenFilesLimit := dbc.OpenFilesLimit / int(numShards)
	for i := 0; i < int(numShards); i++ {
		copiedDBC := *dbc
		copiedDBC.Dir = path.Join(copiedDBC.Dir, strconv.Itoa(i))
		copiedDBC.LevelDBCacheSize = sdbLevelDBCacheSize
		copiedDBC.OpenFilesLimit = sdbOpenFilesLimit

		db, err := newDatabase(&copiedDBC, et)
		if err != nil {
			return nil, err
		}
		shards = append(shards, db)
		go batchWriteWorker(sdbBatchTaskCh)
	}

	logger.Info("Created a sharded database", "dbType", et, "numShards", numShards)
	return &shardedDB{
		fn: dbc.Dir, shards: shards,
		numShards: numShards, sdbBatchTaskCh: sdbBatchTaskCh,
	}, nil
}

// batchWriteWorker executes passed batch tasks.
func batchWriteWorker(batchTasks <-chan sdbBatchTask) {
	for task := range batchTasks {
		task.resultCh <- sdbBatchResult{task.index, task.batch.Write()}
	}
}

// IsPow2 checks if the given number is power of two or not.
func IsPow2(num uint) bool {
	return (num & (num - 1)) == 0
}

// shardIndexByKey returns shard index derived from the given key.
// If len(key) is zero, it returns errKeyLengthZero.
func shardIndexByKey(key []byte, numShards uint) (int, error) {
	if len(key) == 0 {
		return 0, errKeyLengthZero
	}

	return int(key[0]) & (int(numShards) - 1), nil
}

// getShardByKey returns the shard corresponding to the given key.
func (db *shardedDB) getShardByKey(key []byte) (Database, error) {
	if shardIndex, err := shardIndexByKey(key, uint(db.numShards)); err != nil {
		return nil, err
	} else {
		return db.shards[shardIndex], nil
	}
}

func (db *shardedDB) Put(key []byte, value []byte) error {
	if shard, err := db.getShardByKey(key); err != nil {
		return err
	} else {
		return shard.Put(key, value)
	}
}

func (db *shardedDB) Get(key []byte) ([]byte, error) {
	if shard, err := db.getShardByKey(key); err != nil {
		return nil, err
	} else {
		return shard.Get(key)
	}
}

func (db *shardedDB) Has(key []byte) (bool, error) {
	if shard, err := db.getShardByKey(key); err != nil {
		return false, err
	} else {
		return shard.Has(key)
	}
}

func (db *shardedDB) Delete(key []byte) error {
	if shard, err := db.getShardByKey(key); err != nil {
		return err
	} else {
		return shard.Delete(key)
	}
}

func (db *shardedDB) Close() {
	close(db.sdbBatchTaskCh)

	for _, shard := range db.shards {
		shard.Close()
	}
}

// Not enough size of channel slows down the iterator
const shardedDBCombineChanSize = 1024 // Size of resultCh
const shardedDBSubChannelSize = 128   // Size of each sub-channel of resultChs

// shardedDBIterator iterates all items of each shardDB.
// This is useful when you want to get items in serial in binary-alphabetigcal order.
type shardedDBIterator struct {
	parallelIterator shardedDBParallelIterator

	resultCh chan common.Entry
	key      []byte // current key
	value    []byte // current value
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *shardedDB) NewIterator(prefix []byte, start []byte) Iterator {
	it := &shardedDBIterator{
		parallelIterator: db.NewParallelIterator(context.TODO(), prefix, start, nil),
		resultCh:         make(chan common.Entry, shardedDBCombineChanSize),
	}

	go it.runCombineWorker()

	return it
}

// NewIteratorUnsorted creates a iterator over the entire keyspace contained within
// the key-value database. This is useful when you want to get items fast in serial.
// If you want to get ordered items in serial, checkout shardedDB.NewIterator()
// If you want to get items in parallel from channels, checkout shardedDB.NewParallelIterator()
// IteratorUnsorted is a implementation of Iterator and data are accessed with
// Next(), Key() and Value() methods. With ChanIterator, data can be accessed with
// channels. The channels are gained with Channels() method.
func (db *shardedDB) NewIteratorUnsorted(prefix []byte, start []byte) Iterator {
	resultCh := make(chan common.Entry, shardedDBCombineChanSize)
	return &shardedDBIterator{
		parallelIterator: db.NewParallelIterator(context.TODO(), prefix, start, resultCh),
		resultCh:         resultCh,
	}
}

// runCombineWorker fetches any key/value from resultChs and put the data in resultCh
// in binary-alphabetical order.
func (it *shardedDBIterator) runCombineWorker() {
	// creates min-priority queue smallest values from each iterators
	entries := &entryHeap{}
	heap.Init(entries)
	for i, ch := range it.parallelIterator.resultChs {
		if e, ok := <-ch; ok {
			heap.Push(entries, entryWithShardNum{e, i})
		}
	}

chanIter:
	for len(*entries) != 0 {
		// check if done
		select {
		case <-it.parallelIterator.ctx.Done():
			logger.Trace("[shardedDBIterator] combine worker ends due to ctx")
			break chanIter
		default:
		}

		// look for smallest key
		minEntry := heap.Pop(entries).(entryWithShardNum)

		// fill resultCh with smallest key
		it.resultCh <- minEntry.Entry

		// fill used entry with new entry
		// skip this if channel is closed
		if e, ok := <-it.parallelIterator.resultChs[minEntry.shardNum]; ok {
			heap.Push(entries, entryWithShardNum{e, minEntry.shardNum})
		}
	}
	logger.Trace("[shardedDBIterator] combine worker finished")
	close(it.resultCh)
}

// Next gets the next item from iterators.
func (it *shardedDBIterator) Next() bool {
	e, ok := <-it.resultCh
	if !ok {
		logger.Debug("[shardedDBIterator] Next is called on closed channel")
		return false
	}
	it.key, it.value = e.Key, e.Val
	return true
}

func (it *shardedDBIterator) Error() error {
	for i, iter := range it.parallelIterator.iterators {
		if iter.Error() != nil {
			logger.Error("[shardedDBIterator] error from iterator",
				"err", iter.Error(), "shardNum", i, "key", it.key, "val", it.value)
			return iter.Error()
		}
	}
	return nil
}

func (it *shardedDBIterator) Key() []byte {
	return it.key
}

func (it *shardedDBIterator) Value() []byte {
	return it.value
}

func (it *shardedDBIterator) Release() {
	it.parallelIterator.cancel()
}

type entryWithShardNum struct {
	common.Entry
	shardNum int
}

type entryHeap []entryWithShardNum

func (e entryHeap) Len() int {
	return len(e)
}

func (e entryHeap) Less(i, j int) bool {
	return bytes.Compare(e[i].Key, e[j].Key) < 0
}

func (e entryHeap) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *entryHeap) Push(x interface{}) {
	*e = append(*e, x.(entryWithShardNum))
}

func (e *entryHeap) Pop() interface{} {
	old := *e
	n := len(old)
	element := old[n-1]
	*e = old[0 : n-1]
	return element
}

// shardedDBParallelIterator creates iterators for each shard DB.
// Channels subscribing each iterators can be gained.
// Each iterators fetch values in binary-alphabetical order.
// This is useful when you want to operate on each items in parallel.
type shardedDBParallelIterator struct {
	ctx    context.Context
	cancel context.CancelFunc

	iterators []Iterator

	combinedChan bool // all workers put items to one resultChan
	shardNum     int  // num of shards left to iterate
	shardNumMu   *sync.Mutex
	resultChs    []chan common.Entry
}

// NewParallelIterator creates iterators for each shard DB. This is useful when you
// want to operate on each items in parallel.
// If `resultCh` is given, all items are written to `resultCh`, unsorted with a
// particular key prefix, starting at a particular initial key. If `resultCh`
// is not given, new channels are created for each DB. Items are written to
// corresponding channels in binary-alphabetical order. The channels can be
// gained by calling `Channels()`.
//
// If you want to get ordered items in serial, checkout shardedDB.NewIterator()
// If you want to get unordered items in serial with Iterator Interface,
// checkout shardedDB.NewIteratorUnsorted().
func (db *shardedDB) NewParallelIterator(ctx context.Context, prefix []byte, start []byte, resultCh chan common.Entry) shardedDBParallelIterator {
	if ctx == nil {
		ctx = context.TODO()
	}

	it := shardedDBParallelIterator{
		ctx:          ctx,
		cancel:       nil,
		iterators:    make([]Iterator, len(db.shards)),
		combinedChan: resultCh != nil,
		shardNum:     len(db.shards),
		shardNumMu:   &sync.Mutex{},
		resultChs:    make([]chan common.Entry, len(db.shards)),
	}
	it.ctx, it.cancel = context.WithCancel(ctx)

	for i, shard := range db.shards {
		it.iterators[i] = shard.NewIterator(prefix, start)
		if resultCh == nil {
			it.resultChs[i] = make(chan common.Entry, shardedDBSubChannelSize)
		} else {
			it.resultChs[i] = resultCh
		}
		go it.runChanWorker(it.ctx, it.iterators[i], it.resultChs[i])
	}

	return it
}

// runChanWorker runs a worker. The worker gets key/value pair from
// `it` and push the value to `resultCh`.
// `iterator.Release()` is called after all iterating is finished.
// `resultCh` is closed after the iterating is finished.
func (sit *shardedDBParallelIterator) runChanWorker(ctx context.Context, it Iterator, resultCh chan common.Entry) {
iter:
	for it.Next() {
		select {
		case <-ctx.Done():
			break iter
		default:
		}
		key := make([]byte, len(it.Key()))
		val := make([]byte, len(it.Value()))
		copy(key, it.Key())
		copy(val, it.Value())
		resultCh <- common.Entry{Key: key, Val: val}
	}
	// Release the iterator. There is nothing to iterate anymore.
	it.Release()
	// Close `resultCh`. If it is `combinedChan`, the close only happens
	// when this is the last living worker.
	sit.shardNumMu.Lock()
	defer sit.shardNumMu.Unlock()
	if sit.shardNum--; sit.combinedChan && sit.shardNum > 0 {
		return
	}
	close(resultCh)
}

// Channels returns channels that can subscribe on.
func (it *shardedDBParallelIterator) Channels() []chan common.Entry {
	return it.resultChs
}

// Release stops all iterators, channels and workers
// Even Release() is called, there could be some items left in the channel.
// Each iterator.Release() is called in `runChanWorker`.
func (it *shardedDBParallelIterator) Release() {
	it.cancel()
}

func (db *shardedDB) NewBatch() Batch {
	batches := make([]Batch, 0, db.numShards)
	for i := 0; i < int(db.numShards); i++ {
		batches = append(batches, db.shards[i].NewBatch())
	}

	return &shardedDBBatch{
		batches: batches, numBatches: db.numShards,
		taskCh: db.sdbBatchTaskCh, resultCh: make(chan sdbBatchResult, db.numShards),
	}
}

func (db *shardedDB) Type() DBType {
	return ShardedDB
}

func (db *shardedDB) Meter(prefix string) {
	for index, shard := range db.shards {
		shard.Meter(prefix + strconv.Itoa(index))
	}
}

type shardedDBBatch struct {
	batches    []Batch
	numBatches uint

	taskCh   chan sdbBatchTask
	resultCh chan sdbBatchResult
}

func (sdbBatch *shardedDBBatch) Put(key []byte, value []byte) error {
	if ShardIndex, err := shardIndexByKey(key, sdbBatch.numBatches); err != nil {
		return err
	} else {
		return sdbBatch.batches[ShardIndex].Put(key, value)
	}
}

func (sdbBatch *shardedDBBatch) Delete(key []byte) error {
	if ShardIndex, err := shardIndexByKey(key, sdbBatch.numBatches); err != nil {
		return err
	} else {
		return sdbBatch.batches[ShardIndex].Delete(key)
	}
}

// ValueSize is called to determine whether to write batches when it exceeds
// certain limit. shardedDB returns the largest size of its batches to
// write all batches at once when one of batch exceeds the limit.
func (sdbBatch *shardedDBBatch) ValueSize() int {
	maxSize := 0
	for _, batch := range sdbBatch.batches {
		if batch.ValueSize() > maxSize {
			maxSize = batch.ValueSize()
		}
	}
	return maxSize
}

// Write passes the list of batch tasks to taskCh so batch can be processed
// by underlying workers. Write waits until all workers return the result.
func (sdbBatch *shardedDBBatch) Write() error {
	for index, batch := range sdbBatch.batches {
		sdbBatch.taskCh <- sdbBatchTask{batch, index, sdbBatch.resultCh}
	}

	var err error
	for range sdbBatch.batches {
		if batchResult := <-sdbBatch.resultCh; batchResult.err != nil {
			logger.Error("Error while writing sharded batch", "index", batchResult.index, "err", batchResult.err)
			err = batchResult.err
		}
	}
	// Leave logs for each error but only return the last one.
	return err
}

func (sdbBatch *shardedDBBatch) Reset() {
	for _, batch := range sdbBatch.batches {
		batch.Reset()
	}
}

func (sdbBatch *shardedDBBatch) Replay(w KeyValueWriter) error {
	for _, batch := range sdbBatch.batches {
		if err := batch.Replay(w); err != nil {
			return err
		}
	}
	return nil
}
