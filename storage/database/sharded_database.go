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
	"errors"
	"fmt"
	"path"
	"strconv"
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

type shardedDBIterator struct {
	iterators []Iterator // iterators for each shard
	finish    []bool     // check if each iterator is finished
	minIdx    int        // iterator index that has min key

	keys [][]byte // keys of each iterator
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
	// create iterators for each shard
	iterators := make([]Iterator, 0, pdb.numShards)
	for i := 0; i < int(pdb.numShards); i++ {
		iterators = append(iterators, pdb.shards[i].NewIteratorWithStart(start))
	}

	return &shardedDBIterator{iterators, make([]bool, pdb.numShards), -1, make([][]byte, pdb.numShards)}
}

// NewIteratorWithPrefix creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix.
func (pdb *shardedDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	// TODO-Klaytn implement this later.
	return nil
}

func (pdi *shardedDBIterator) Next() bool {
	if pdi.minIdx == -1 { // it this is the first call if Next, call Next on each one
		// check if next exists for each iterator
		for idx, iter := range pdi.iterators {
			if !iter.Next() {
				pdi.finish[idx] = true
			} else {
				//pdi.keys[idx] = make([]byte, len(iter.Value()))
				//copy(pdi.keys[idx], iter.Value())
			}
		}
	} else if !pdi.iterators[pdi.minIdx].Next() { // go to next one in previous min iterator
		//pdi.keys[pdi.minIdx] = pdi.iterators[pdi.minIdx].Value()
	} else {
		//pdi.keys[pdi.minIdx] = []byte{255}
		pdi.finish[pdi.minIdx] = true // if there is no item in iterator, finish it
	}

	// find the iterator with min key
	minKey := []byte{255, 255, 255, 255, 255} // TODO default value should be max
	minIdx := -1

	for idx, iter := range pdi.iterators {
		if !pdi.finish[idx] {
			if bytes.Compare(minKey, iter.Key()) >= 0 {
				minIdx = idx
				minKey = iter.Key()
			}
		}
	}

	if minIdx == -1 {
		return false
	}

	pdi.minIdx = minIdx

	return true
}

func (pdi *shardedDBIterator) Error() error {
	errs := errors.New("")
	for idx, iter := range pdi.iterators {
		if iter != nil {
			errs = fmt.Errorf("%v [iter%d]%v", errs.Error(), idx, iter.Error().Error())
		}
	}
	return nil
}

func (pdi *shardedDBIterator) Key() []byte {
	return pdi.iterators[pdi.minIdx].Key()
}

func (pdi *shardedDBIterator) Value() []byte {
	return pdi.iterators[pdi.minIdx].Value()
}

func (pdi *shardedDBIterator) Release() {
	for _, iter := range pdi.iterators {
		if iter != nil {
			iter.Release()
		}
	}
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
