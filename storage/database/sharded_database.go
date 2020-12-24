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
	"context"
	"fmt"
	"path"
	"strconv"

	"github.com/klaytn/klaytn/common/prque"
)

var errKeyLengthZero = fmt.Errorf("database key for sharded database should be greater than 0")

type shardedDB struct {
	fn        string
	shards    []Database
	numShards uint

	pdbBatchTaskCh chan pdbBatchTask
}

type pdbBatchTask struct {
	batch    Batch               // A batch that each worker executes.
	index    int                 // Index of given batch.
	resultCh chan pdbBatchResult // Batch result channel for each shardedDBBatch.
}

type pdbBatchResult struct {
	index int   // Index of the batch result.
	err   error // Error from the batch write operation.
}

// newShardedDB creates database with numShards shards, or partitions.
// The type of database is specified DBConfig.DBType.
func newShardedDB(dbc *DBConfig, et DBEntryType, numShards uint) (*shardedDB, error) {
	const numShardsLimit = 16

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
	pdbBatchTaskCh := make(chan pdbBatchTask, numShards*2)
	for i := 0; i < int(numShards); i++ {
		copiedDBC := *dbc
		copiedDBC.Dir = path.Join(copiedDBC.Dir, strconv.Itoa(i))
		copiedDBC.LevelDBCacheSize /= int(numShards)

		db, err := newDatabase(&copiedDBC, et)
		if err != nil {
			return nil, err
		}
		shards = append(shards, db)
		go batchWriteWorker(pdbBatchTaskCh)
	}

	return &shardedDB{
		fn: dbc.Dir, shards: shards,
		numShards: numShards, pdbBatchTaskCh: pdbBatchTaskCh}, nil
}

// batchWriteWorker executes passed batch tasks.
func batchWriteWorker(batchTasks <-chan pdbBatchTask) {
	for task := range batchTasks {
		task.resultCh <- pdbBatchResult{task.index, task.batch.Write()}
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
func (pdb *shardedDB) getShardByKey(key []byte) (Database, error) {
	if shardIndex, err := shardIndexByKey(key, uint(pdb.numShards)); err != nil {
		return nil, err
	} else {
		return pdb.shards[shardIndex], nil
	}
}

func (pdb *shardedDB) Put(key []byte, value []byte) error {
	if shard, err := pdb.getShardByKey(key); err != nil {
		return err
	} else {
		return shard.Put(key, value)
	}
}

func (pdb *shardedDB) Get(key []byte) ([]byte, error) {
	if shard, err := pdb.getShardByKey(key); err != nil {
		return nil, err
	} else {
		return shard.Get(key)
	}
}

func (pdb *shardedDB) Has(key []byte) (bool, error) {
	if shard, err := pdb.getShardByKey(key); err != nil {
		return false, err
	} else {
		return shard.Has(key)
	}
}

func (pdb *shardedDB) Delete(key []byte) error {
	if shard, err := pdb.getShardByKey(key); err != nil {
		return err
	} else {
		return shard.Delete(key)
	}
}

func (pdb *shardedDB) Close() {
	close(pdb.pdbBatchTaskCh)

	for _, shard := range pdb.shards {
		shard.Close()
	}
}

// Not enough size of channel slows down the iterator
const shardedDBCombineChanSize = 1024 // Size of resultCh
const shardedDBSubChannelSize = 128   // Size of each channel of resultChs

// shardedDBIterator iterates all items of each shardDB.
// This is useful when you want to get items in serial in binary-alphabetical order.
type shardedDBIterator struct {
	shardedDBChanIterator

	resultCh chan entry
	key      []byte // current key
	value    []byte // current value
}

// NewIterator creates a binary-alphabetical iterator over the entire keyspace
// contained within the key-value database.
// If you want to get unordered items faster in serial, checkout shardedDB.NewIteratorUnsorted().
// If you want to get items in parallel from channels, checkout shardedDB.NewChanIterator()
func (db *shardedDB) NewIterator() Iterator {
	return db.newIterator(func(db Database) Iterator { return db.NewIterator() })
}

// NewIteratorWithStart creates a binary-alphabetical iterator over a subset of
// database content starting at a particular initial key (or after, if it does
// not exist).
func (db *shardedDB) NewIteratorWithStart(start []byte) Iterator {
	return db.newIterator(func(db Database) Iterator { return db.NewIteratorWithStart(start) })
}

// NewIteratorWithPrefix creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix.
func (db *shardedDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	return db.newIterator(func(db Database) Iterator { return db.NewIteratorWithPrefix(prefix) })
}

func (db *shardedDB) newIterator(newIterator func(Database) Iterator) Iterator {
	it := &shardedDBIterator{
		db.NewChanIterator(context.Background(), nil, newIterator),
		make(chan entry, shardedDBCombineChanSize),
		nil, nil}

	go it.runCombineWorker()

	return it
}

// runCombineWorker fetches any key/value from resultChs and put the data in resultCh
// in binary-alphabetical order.
func (it shardedDBIterator) runCombineWorker() {
	type entryWithShardNum struct {
		entry
		shardNum int
	}

	// creates min-priority queue smallest values from each iterators
	entries := prque.NewByteSlice(true)
	for i, ch := range it.resultChs {
		if e, ok := <-ch; ok {
			entries.Push(entryWithShardNum{e, i}, e.key)
		}
	}

chanIter:
	for !entries.Empty() {
		// check if done
		select {
		case <-it.ctx.Done():
			logger.Trace("[shardedDBIterator] combine worker ends due to ctx")
			break chanIter
		default:
		}

		// look for smallest key
		minEntry := entries.PopItem().(entryWithShardNum)

		// fill resultCh with smallest key
		it.resultCh <- minEntry.entry

		// fill used entry with new entry
		// skip this if channel is closed
		if e, ok := <-it.resultChs[minEntry.shardNum]; ok {
			entries.Push(entryWithShardNum{e, minEntry.shardNum}, e.key)
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
	it.key, it.value = e.key, e.val
	return true
}

func (it *shardedDBIterator) Error() error {
	for i, iter := range it.iterators {
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
	it.cancel()
}

type entry struct {
	key, val []byte
}

// shardedDBIteratorUnsorted iterates all items of each shardDB.
// This is useful when you want to get items fast in serial.
type shardedDBIteratorUnsorted struct {
	shardedDBIterator
}

// NewIteratorUnsorted creates a iterator over the entire keyspace contained within
// the key-value database.
// If you want to get ordered items in serial, checkout shardedDB.NewIterator()
// If you want to get items in parallel from channels, checkout shardedDB.NewChanIterator()
func (db *shardedDB) NewIteratorUnsorted() Iterator {
	return db.newIteratorUnsorted(func(db Database) Iterator { return db.NewIterator() })
}

// NewIteratorWithStartUnsorted creates a iterator over a subset of database content
// starting at a particular initial key (or after, if it does not exist).
func (db *shardedDB) NewIteratorWithStartUnsorted(start []byte) Iterator {
	return db.newIteratorUnsorted(func(db Database) Iterator { return db.NewIteratorWithStart(start) })
}

// NewIteratorWithPrefixUnsorted creates a iterator over a subset of database content
// with a particular key prefix.
func (db *shardedDB) NewIteratorWithPrefixUnsorted(prefix []byte) Iterator {
	return db.newIteratorUnsorted(func(db Database) Iterator { return db.NewIteratorWithPrefix(prefix) })
}

func (db *shardedDB) newIteratorUnsorted(newIterator func(Database) Iterator) Iterator {
	resultCh := make(chan entry, shardedDBCombineChanSize)
	it := &shardedDBIteratorUnsorted{
		shardedDBIterator{
			db.NewChanIterator(context.Background(), resultCh, newIterator),
			resultCh,
			nil, nil}}
	return it
}

// shardedDBChanIterator creates iterators for each shard DB.
// Channels subscribing each iterators can be gained.
// Each iterators fetch values in binary-alphabetical order.
// This is useful when you want to operate on each items in parallel.
type shardedDBChanIterator struct {
	ctx    context.Context
	cancel context.CancelFunc

	iterators []Iterator

	combinedChan bool // all workers put items to one resultChan
	shardNum     int  // num of shards left to iterate
	resultChs    []chan entry
}

// NewChanIterator creates iterators for each shard DB. This is useful when you
// want to operate on each items in parallel.
// If `resultCh` is given, all items are written to `resultCh`, unsorted. If
// `resultCh` is not given, new channels are created for each DB. Items are
// written to corresponding channels in binary-alphabetical order. The channels
// can be gained by calling `Channels()`.
//
// If you want to get ordered items in serial, checkout shardedDB.NewIterator()
// If you want to get unordered items in serial with Iterator Interface,
// checkout shardedDB.NewIteratorUnsorted().
func (db *shardedDB) NewChanIterator(ctx context.Context, resultCh chan entry, newIterator func(Database) Iterator) shardedDBChanIterator {
	if ctx == nil {
		ctx = context.Background()
	}

	it := shardedDBChanIterator{ctx, nil,
		make([]Iterator, len(db.shards)),
		resultCh != nil,
		len(db.shards),
		make([]chan entry, len(db.shards))}
	it.ctx, it.cancel = context.WithCancel(ctx)

	for i := 0; i < len(db.shards); i++ {
		it.iterators[i] = newIterator(db.shards[i])
		if resultCh == nil {
			it.resultChs[i] = make(chan entry, shardedDBSubChannelSize)
		} else {
			it.resultChs[i] = resultCh
		}
		go it.runChanWorker(it.iterators[i], it.resultChs[i], it.ctx)
	}

	return it
}

func (sit *shardedDBChanIterator) runChanWorker(it Iterator, resultCh chan entry, ctx context.Context) {
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
		resultCh <- entry{key, val}
	}
	it.Release()
	if sit.shardNum--; sit.combinedChan && sit.shardNum > 0 {
		return
	}
	close(resultCh)
}

// Channels returns channels that can subscribe on.
func (it *shardedDBChanIterator) Channels() []chan entry {
	return it.resultChs
}

// Release stops all iterators, channels and workers
// Even Release() is called, there could be some items left in the channel.
func (it *shardedDBChanIterator) Release() {
	it.cancel()
}

func (pdb *shardedDB) NewBatch() Batch {
	batches := make([]Batch, 0, pdb.numShards)
	for i := 0; i < int(pdb.numShards); i++ {
		batches = append(batches, pdb.shards[i].NewBatch())
	}

	return &shardedDBBatch{batches: batches, numBatches: pdb.numShards,
		taskCh: pdb.pdbBatchTaskCh, resultCh: make(chan pdbBatchResult, pdb.numShards)}
}

func (pdb *shardedDB) Type() DBType {
	return ShardedDB
}

func (pdb *shardedDB) Meter(prefix string) {
	for index, shard := range pdb.shards {
		shard.Meter(prefix + strconv.Itoa(index))
	}
}

type shardedDBBatch struct {
	batches    []Batch
	numBatches uint

	taskCh   chan pdbBatchTask
	resultCh chan pdbBatchResult
}

func (pdbBatch *shardedDBBatch) Put(key []byte, value []byte) error {
	if ShardIndex, err := shardIndexByKey(key, uint(pdbBatch.numBatches)); err != nil {
		return err
	} else {
		return pdbBatch.batches[ShardIndex].Put(key, value)
	}
}

// ValueSize is called to determine whether to write batches when it exceeds
// certain limit. shardedDB returns the largest size of its batches to
// write all batches at once when one of batch exceeds the limit.
func (pdbBatch *shardedDBBatch) ValueSize() int {
	maxSize := 0
	for _, batch := range pdbBatch.batches {
		if batch.ValueSize() > maxSize {
			maxSize = batch.ValueSize()
		}
	}
	return maxSize
}

// Write passes the list of batch tasks to taskCh so batch can be processed
// by underlying workers. Write waits until all workers return the result.
func (pdbBatch *shardedDBBatch) Write() error {
	for index, batch := range pdbBatch.batches {
		pdbBatch.taskCh <- pdbBatchTask{batch, index, pdbBatch.resultCh}
	}

	var err error
	for range pdbBatch.batches {
		if batchResult := <-pdbBatch.resultCh; batchResult.err != nil {
			logger.Error("Error while writing sharded batch", "index", batchResult.index, "err", batchResult.err)
			err = batchResult.err
		}
	}
	// Leave logs for each error but only return the last one.
	return err
}

func (pdbBatch *shardedDBBatch) Reset() {
	for _, batch := range pdbBatch.batches {
		batch.Reset()
	}
}
