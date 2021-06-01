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

import "strings"

// Code using batches should try to add this much data to the batch.
// The value was determined empirically.

type DBType string

const (
	LevelDB   DBType = "LevelDB"
	BadgerDB         = "BadgerDB"
	MemoryDB         = "MemoryDB"
	DynamoDB         = "DynamoDBS3"
	ShardedDB        = "ShardedDB"
)

// ToValid converts DBType to a valid one.
// If it is unable to convert, "" is returned.
func (db DBType) ToValid() DBType {
	validDBType := []DBType{LevelDB, BadgerDB, MemoryDB, DynamoDB}

	for _, vdb := range validDBType {
		if strings.ToLower(string(vdb)) == strings.ToLower(string(db)) {
			return vdb
		}
	}

	return ""
}

// selfShardable returns if the db is able to shard by itself or not
func (db DBType) selfShardable() bool {
	switch db {
	case DynamoDB:
		return true
	}
	return false
}

// KeyValueWriter wraps the Put method of a backing data store.
type KeyValueWriter interface {
	// Put inserts the given value into the key-value data store.
	Put(key []byte, value []byte) error

	// Delete removes the key from the key-value data store.
	Delete(key []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	KeyValueWriter
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Close()
	NewBatch() Batch
	Type() DBType
	Meter(prefix string)
	Iteratee
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
