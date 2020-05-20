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
	"fmt"
	"path"
	"strconv"
)

var errKeyLengthZero = fmt.Errorf("database key for partitioned database should be greater than 0")

type partitionedDB struct {
	fn            string
	partitions    []Database
	numPartitions uint

	pdbBatchTaskCh chan pdbBatchTask
}

type pdbBatchTask struct {
	batch    Batch               // A batch that each worker executes.
	index    int                 // Index of given batch.
	resultCh chan pdbBatchResult // Batch result channel for each partitionedDBBatch.
}

type pdbBatchResult struct {
	index int   // Index of the batch result.
	err   error // Error from the batch write operation.
}

func newPartitionedDB(dbc *DBConfig, et DBEntryType, numPartitions uint) (*partitionedDB, error) {
	const numPartitionsLimit = 16

	if numPartitions == 0 {
		logger.Crit("numPartitions should be greater than 0!")
	}

	if numPartitions > numPartitionsLimit {
		logger.Crit(fmt.Sprintf("numPartitions should be equal to or smaller than %v, but it is %v.", numPartitionsLimit, numPartitions))
	}

	if !IsPow2(numPartitions) {
		logger.Crit(fmt.Sprintf("numPartitions should be power of two, but it is %v", numPartitions))
	}

	partitions := make([]Database, 0, numPartitions)
	pdbBatchTaskCh := make(chan pdbBatchTask, numPartitions*2)
	for i := 0; i < int(numPartitions); i++ {
		copiedDBC := *dbc
		copiedDBC.Dir = path.Join(copiedDBC.Dir, strconv.Itoa(i))
		copiedDBC.LevelDBCacheSize /= int(numPartitions)

		db, err := newDatabase(&copiedDBC, et)
		if err != nil {
			return nil, err
		}
		partitions = append(partitions, db)
		go batchWriteWorker(pdbBatchTaskCh)
	}

	return &partitionedDB{
		fn: dbc.Dir, partitions: partitions,
		numPartitions: numPartitions, pdbBatchTaskCh: pdbBatchTaskCh}, nil
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

// calcPartition returns partition index derived from the given key.
// If len(key) is zero, it returns errKeyLengthZero.
func calcPartition(key []byte, numPartitions uint) (int, error) {
	if len(key) == 0 {
		return 0, errKeyLengthZero
	}

	return int(key[0]) & (int(numPartitions) - 1), nil
}

// getPartition returns the partition corresponding to the given key.
func (pdb *partitionedDB) getPartition(key []byte) (Database, error) {
	if partitionIndex, err := calcPartition(key, uint(pdb.numPartitions)); err != nil {
		return nil, err
	} else {
		return pdb.partitions[partitionIndex], nil
	}
}

func (pdb *partitionedDB) Put(key []byte, value []byte) error {
	if partition, err := pdb.getPartition(key); err != nil {
		return err
	} else {
		return partition.Put(key, value)
	}
}

func (pdb *partitionedDB) Get(key []byte) ([]byte, error) {
	if partition, err := pdb.getPartition(key); err != nil {
		return nil, err
	} else {
		return partition.Get(key)
	}
}

func (pdb *partitionedDB) Has(key []byte) (bool, error) {
	if partition, err := pdb.getPartition(key); err != nil {
		return false, err
	} else {
		return partition.Has(key)
	}
}

func (pdb *partitionedDB) Delete(key []byte) error {
	if partition, err := pdb.getPartition(key); err != nil {
		return err
	} else {
		return partition.Delete(key)
	}
}

func (pdb *partitionedDB) Close() {
	close(pdb.pdbBatchTaskCh)

	for _, partition := range pdb.partitions {
		partition.Close()
	}
}

type partitionedDBIterator struct {
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
func (pdb *partitionedDB) NewIterator() Iterator {
	// TODO-Klaytn implement this later.
	return nil
}

// NewIteratorWithStart creates a binary-alphabetical iterator over a subset of
// database content starting at a particular initial key (or after, if it does
// not exist).
func (pdb *partitionedDB) NewIteratorWithStart(start []byte) Iterator {
	iterators := make([]Iterator, 0, pdb.numPartitions)
	for i := 0; i < int(pdb.numPartitions); i++ {
		iterators = append(iterators, pdb.partitions[i].NewIteratorWithStart(start))
	}

	for _, iter := range iterators {
		if iter != nil {
			if !iter.Next() {
				iter = nil
			}
		}
	}

	return &partitionedDBIterator{iterators, nil, nil}
}

// NewIteratorWithPrefix creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix.
func (pdb *partitionedDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	// TODO-Klaytn implement this later.
	return nil
}

func (pdi *partitionedDBIterator) Next() bool {
	var minIter Iterator
	minIdx := -1
	minKey := []byte{0}
	minKeyValue := []byte{0}

	for idx, iter := range pdi.iterators {
		if iter != nil {
			if bytes.Compare(minKey, iter.Key()) >= 0 {
				minIdx = idx
				minIter = iter
				minKey = iter.Key()
				minKeyValue = iter.Value()
			}
		}
	}

	if minIter == nil {
		return false
	}

	pdi.key = minKey
	pdi.value = minKeyValue

	if !minIter.Next() {
		pdi.iterators[minIdx] = nil
	}

	return false
}

func (pdi *partitionedDBIterator) Error() error {
	return nil
}

func (pdi *partitionedDBIterator) Key() []byte {
	return nil
}

func (pdi *partitionedDBIterator) Value() []byte {
	return nil
}

func (pdi *partitionedDBIterator) Release() {

}

func (pdb *partitionedDB) NewBatch() Batch {
	batches := make([]Batch, 0, pdb.numPartitions)
	for i := 0; i < int(pdb.numPartitions); i++ {
		batches = append(batches, pdb.partitions[i].NewBatch())
	}

	return &partitionedDBBatch{batches: batches, numBatches: pdb.numPartitions,
		taskCh: pdb.pdbBatchTaskCh, resultCh: make(chan pdbBatchResult, pdb.numPartitions)}
}

func (pdb *partitionedDB) Type() DBType {
	return PartitionedDB
}

func (pdb *partitionedDB) Meter(prefix string) {
	for index, partition := range pdb.partitions {
		partition.Meter(prefix + strconv.Itoa(index))
	}
}

type partitionedDBBatch struct {
	batches    []Batch
	numBatches uint

	taskCh   chan pdbBatchTask
	resultCh chan pdbBatchResult
}

func (pdbBatch *partitionedDBBatch) Put(key []byte, value []byte) error {
	if partitionIndex, err := calcPartition(key, uint(pdbBatch.numBatches)); err != nil {
		return err
	} else {
		return pdbBatch.batches[partitionIndex].Put(key, value)
	}
}

// ValueSize is called to determine whether to write batches when it exceeds
// certain limit. partitionedDB returns the largest size of its batches to
// write all batches at once when one of batch exceeds the limit.
func (pdbBatch *partitionedDBBatch) ValueSize() int {
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
func (pdbBatch *partitionedDBBatch) Write() error {
	for index, batch := range pdbBatch.batches {
		pdbBatch.taskCh <- pdbBatchTask{batch, index, pdbBatch.resultCh}
	}

	var err error
	for range pdbBatch.batches {
		if batchResult := <-pdbBatch.resultCh; batchResult.err != nil {
			logger.Error("Error while writing partitioned batch", "index", batchResult.index, "err", batchResult.err)
			err = batchResult.err
		}
	}
	// Leave logs for each error but only return the last one.
	return err
}

func (pdbBatch *partitionedDBBatch) Reset() {
	for _, batch := range pdbBatch.batches {
		batch.Reset()
	}
}
