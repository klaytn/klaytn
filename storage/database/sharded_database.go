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
	"fmt"
	"path"
	"strconv"
)

var errKeyLengthZero = fmt.Errorf("database key for sharded database should be greater than 0")

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
		numShards: numShards, sdbBatchTaskCh: sdbBatchTaskCh}, nil
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

type shardedDBIterator struct {
	// TODO-Klaytn implement this later.
	iterators []Iterator
	key       []byte
	value     []byte

	//numBatches uint
	//
	//taskCh   chan pdbBatchTask
	//resultCh chan pdbBatchResult
}

// NewIterator creates a binary-alphabetical iterator over the entire keyspace
// contained within the key-value database.
func (pdb *shardedDB) NewIterator() Iterator {
	// TODO-Klaytn implement this later.
	return nil
}

// NewIteratorWithStart creates a binary-alphabetical iterator over a subset of
// database content starting at a particular initial key (or after, if it does
// not exist).
func (pdb *shardedDB) NewIteratorWithStart(start []byte) Iterator {
	// TODO-Klaytn implement this later.
	iterators := make([]Iterator, 0, pdb.numShards)
	for i := 0; i < int(pdb.numShards); i++ {
		iterators = append(iterators, pdb.shards[i].NewIteratorWithStart(start))
	}

	for _, iter := range iterators {
		if iter != nil {
			if !iter.Next() {
				iter = nil
			}
		}
	}

	return &shardedDBIterator{iterators, nil, nil}
}

// NewIteratorWithPrefix creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix.
func (pdb *shardedDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	// TODO-Klaytn implement this later.
	return nil
}

func (pdi *shardedDBIterator) Next() bool {
	// TODO-Klaytn implement this later.
	//var minIter Iterator
	//minIdx := -1
	//minKey := []byte{0}
	//minKeyValue := []byte{0}
	//
	//for idx, iter := range pdi.iterators {
	//	if iter != nil {
	//		if bytes.Compare(minKey, iter.Key()) >= 0 {
	//			minIdx = idx
	//			minIter = iter
	//			minKey = iter.Key()
	//			minKeyValue = iter.Value()
	//		}
	//	}
	//}
	//
	//if minIter == nil {
	//	return false
	//}
	//
	//pdi.key = minKey
	//pdi.value = minKeyValue
	//
	//if !minIter.Next() {
	//	pdi.iterators[minIdx] = nil
	//}
	//
	return true
}

func (pdi *shardedDBIterator) Error() error {
	// TODO-Klaytn implement this later.
	return nil
}

func (pdi *shardedDBIterator) Key() []byte {
	// TODO-Klaytn implement this later.
	return nil
}

func (pdi *shardedDBIterator) Value() []byte {
	// TODO-Klaytn implement this later.
	return nil
}

func (pdi *shardedDBIterator) Release() {
	// TODO-Klaytn implement this later.
}

func (db *shardedDB) NewBatch() Batch {
	batches := make([]Batch, 0, db.numShards)
	for i := 0; i < int(db.numShards); i++ {
		batches = append(batches, db.shards[i].NewBatch())
	}

	return &shardedDBBatch{batches: batches, numBatches: db.numShards,
		taskCh: db.sdbBatchTaskCh, resultCh: make(chan sdbBatchResult, db.numShards)}
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
	if ShardIndex, err := shardIndexByKey(key, uint(sdbBatch.numBatches)); err != nil {
		return err
	} else {
		return sdbBatch.batches[ShardIndex].Put(key, value)
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
