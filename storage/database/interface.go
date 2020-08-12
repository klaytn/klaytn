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
// This file is derived from ethdb/interface.go (2018/06/04).
// Modified and improved for the klaytn development.

package database

// Code using batches should try to add this much data to the batch.
// The value was determined empirically.

type DBType uint8

const (
	_ DBType = iota
	LevelDB
	BadgerDB
	MemoryDB
	PartitionedDB
	DynamoDB
)

func (dbType DBType) String() string {
	switch dbType {
	case LevelDB:
		return "LevelDB"
	case BadgerDB:
		return "BadgerDB"
	case MemoryDB:
		return "MemoryDB"
	case PartitionedDB:
		return "PartitionedDB"
	case DynamoDB:
		return "DynamoDB"
	default:
		logger.Error("Undefined DBType entered.", "entered DBType", dbType)
		return "undefined"
	}
}

const IdealBatchSize = 100 * 1024

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(key []byte, value []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	Putter
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close()
	NewBatch() Batch
	Type() DBType
	Meter(prefix string)
	Iteratee
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type Batch interface {
	Putter
	ValueSize() int // amount of data in the batch
	Write() error
	// Reset resets the batch for reuse
	Reset()
}

func WriteBatches(batches ...Batch) (int, error) {
	bytes := 0
	for _, batch := range batches {
		if batch.ValueSize() > 0 {
			bytes += batch.ValueSize()
			if err := batch.Write(); err != nil {
				return 0, err
			}
			batch.Reset()
		}
	}
	return bytes, nil
}

func WriteBatchesParallel(batches ...Batch) (int, error) {
	type result struct {
		bytes int
		err   error
	}

	resultCh := make(chan result, len(batches))
	for _, batch := range batches {
		go func(batch Batch) {
			bytes := batch.ValueSize()
			err := batch.Write()
			if err != nil {
				batch.Reset()
			}
			resultCh <- result{bytes, err}
		}(batch)
	}

	var bytes int
	for range batches {
		rst := <-resultCh
		if rst.err != nil {
			return bytes, rst.err
		}
		bytes += rst.bytes
	}

	return bytes, nil
}

func WriteBatchesOverThreshold(batches ...Batch) (int, error) {
	bytes := 0
	for _, batch := range batches {
		if batch.ValueSize() >= IdealBatchSize {
			if err := batch.Write(); err != nil {
				return 0, err
			}
			bytes += batch.ValueSize()
			batch.Reset()
		}
	}
	return bytes, nil
}

func PutAndWriteBatchesOverThreshold(batch Batch, key, val []byte) error {
	if err := batch.Put(key, val); err != nil {
		return err
	}

	if batch.ValueSize() >= IdealBatchSize {
		if err := batch.Write(); err != nil {
			return err
		}
		batch.Reset()
	}

	return nil
}
