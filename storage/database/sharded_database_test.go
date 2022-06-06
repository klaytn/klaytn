// Copyright 2021 The klaytn Authors
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
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/klaytn/klaytn/common"

	"github.com/stretchr/testify/assert"
)

var ShardedDBConfig = []*DBConfig{
	{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 2, ParallelDBWrite: true},
	{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 4, ParallelDBWrite: true},
}

// testIterator tests if given iterator iterates all entries
func testIterator(t *testing.T, checkOrder bool, entryNums []uint, dbConfig []*DBConfig, entriesFromIterator func(t *testing.T, db shardedDB, entryNum uint) []common.Entry) {
	for _, entryNum := range entryNums {
		entries := common.CreateEntries(int(entryNum))
		dbs := make([]shardedDB, len(dbConfig))

		// create DB and write data for testing
		for i, config := range dbConfig {
			config.Dir, _ = ioutil.TempDir(os.TempDir(), "test-shardedDB-iterator")
			defer func(dir string) {
				if err := os.RemoveAll(dir); err != nil {
					t.Fatalf("fail to delete file %v", err)
				}
			}(config.Dir)

			// create sharded DB
			db, err := newShardedDB(config, 0, config.NumStateTrieShards)
			if err != nil {
				t.Log("Error occured while creating DB :", err)
				t.FailNow()
			}
			dbs[i] = *db

			// write entries data in DB
			batch := db.NewBatch()
			for _, entry := range entries {
				assert.NoError(t, batch.Put(entry.Key, entry.Val))
			}
			assert.NoError(t, batch.Write())
		}

		// sort entries for each compare
		sort.Slice(entries, func(i, j int) bool { return bytes.Compare(entries[i].Key, entries[j].Key) < 0 })

		// get data from iterator and compare
		for _, db := range dbs {
			// create iterator
			entriesFromIt := entriesFromIterator(t, db, entryNum)
			if !checkOrder {
				sort.Slice(entriesFromIt, func(i, j int) bool { return bytes.Compare(entriesFromIt[i].Key, entriesFromIt[j].Key) < 0 })
			}

			// compare if entries generated and entries from iterator is same
			assert.Equal(t, len(entries), len(entriesFromIt))
			assert.True(t, reflect.DeepEqual(entries, entriesFromIt))
		}
	}
}

// TestShardedDBIterator tests if shardedDBIterator iterates all entries with diverse shard size
func TestShardedDBIterator(t *testing.T) {
	testIterator(t, true, []uint{100}, ShardedDBConfig, newShardedDBIterator)
}

// TestShardedDBIteratorUnsorted tests if shardedDBIteratorUnsorted iterates all entries with diverse shard size
func TestShardedDBIteratorUnsorted(t *testing.T) {
	testIterator(t, false, []uint{100}, ShardedDBConfig, newShardedDBIteratorUnsorted)
}

// TestShardedDBParallelIterator tests if shardedDBParallelIterator iterates all entries with diverse shard size
func TestShardedDBParallelIterator(t *testing.T) {
	testIterator(t, false, []uint{100}, ShardedDBConfig, newShardedDBParallelIterator)
}

// TestShardedDBIteratorSize tests if shardedDBIterator iterates all entries for different
// entry sizes
func TestShardedDBIteratorSize(t *testing.T) {
	config := ShardedDBConfig[0]
	size := config.NumStateTrieShards
	testIterator(t, true, []uint{size - 1, size, size + 1}, []*DBConfig{config}, newShardedDBIterator)
}

// TestShardedDBIteratorUnsortedSize tests if shardedDBIteratorUnsorted iterates all entries
func TestShardedDBIteratorUnsortedSize(t *testing.T) {
	config := ShardedDBConfig[0]
	size := config.NumStateTrieShards
	testIterator(t, false, []uint{size - 1, size, size + 1}, []*DBConfig{config}, newShardedDBIteratorUnsorted)
}

// TestShardedDBParallelIteratorSize tests if shardedDBParallelIterator iterates all entries
func TestShardedDBParallelIteratorSize(t *testing.T) {
	config := ShardedDBConfig[0]
	size := config.NumStateTrieShards
	testIterator(t, false, []uint{size - 1, size, size + 1}, []*DBConfig{config}, newShardedDBParallelIterator)
}

func newShardedDBIterator(t *testing.T, db shardedDB, entryNum uint) []common.Entry {
	entries := make([]common.Entry, 0, entryNum)
	it := db.NewIterator(nil, nil)

	for it.Next() {
		entries = append(entries, common.Entry{Key: it.Key(), Val: it.Value()})
	}
	it.Release()
	assert.NoError(t, it.Error())
	return entries
}

func newShardedDBIteratorUnsorted(t *testing.T, db shardedDB, entryNum uint) []common.Entry {
	entries := make([]common.Entry, 0, entryNum)
	it := db.NewIteratorUnsorted(nil, nil)

	for it.Next() {
		entries = append(entries, common.Entry{Key: it.Key(), Val: it.Value()})
	}
	it.Release()
	assert.NoError(t, it.Error())
	return entries
}

func newShardedDBParallelIterator(t *testing.T, db shardedDB, entryNum uint) []common.Entry {
	entries := make([]common.Entry, 0, entryNum) // store all items
	var l sync.RWMutex                           // mutex for entries

	// create chan Iterator and get channels
	it := db.NewParallelIterator(context.Background(), nil, nil, nil)
	chans := it.Channels()

	// listen all channels and get key/value
	done := make(chan struct{})
	for _, ch := range chans {
		go func(ch chan common.Entry) {
			for e := range ch {
				l.Lock()
				entries = append(entries, e)
				l.Unlock()
			}
			done <- struct{}{} // tell
		}(ch)
	}
	// wait for all iterators to finish
	for range chans {
		<-done
	}
	close(done)
	it.Release()
	return entries
}

func testShardedIterator_Release(t *testing.T, entryNum int, checkFunc func(db shardedDB)) {
	entries := common.CreateEntries(entryNum)

	// create DB and write data for testing
	for _, config := range ShardedDBConfig {
		config.Dir, _ = ioutil.TempDir(os.TempDir(), "test-shardedDB-iterator")
		defer func(dir string) {
			if err := os.RemoveAll(dir); err != nil {
				t.Fatalf("fail to delete file %v", err)
			}
		}(config.Dir)

		// create sharded DB
		db, err := newShardedDB(config, MiscDB, config.NumStateTrieShards)
		if err != nil {
			t.Log("Error occured while creating DB :", err)
			t.FailNow()
		}

		// write entries data in DB
		batch := db.NewBatch()
		for _, entry := range entries {
			assert.NoError(t, batch.Put(entry.Key, entry.Val))
		}
		assert.NoError(t, batch.Write())

		// check if Release quits iterator
		checkFunc(*db)
	}
}

func TestShardedDBIterator_Release(t *testing.T) {
	testShardedIterator_Release(t, shardedDBCombineChanSize+10, func(db shardedDB) {
		// Next() returns True if Release() is not called
		{
			it := db.NewIterator(nil, nil)
			defer it.Release()

			// check if data exists
			for i := 0; i < shardedDBCombineChanSize+1; i++ {
				assert.True(t, it.Next())
			}
		}

		//  Next() returns False if Release() is called
		{
			it := db.NewIterator(nil, nil)
			it.Release()

			// flush data in channel
			for i := 0; i < shardedDBCombineChanSize+1; i++ {
				it.Next()
			}

			// check if Next returns false
			assert.False(t, it.Next())
		}
	})
}

func TestShardedDBIteratorUnsorted_Release(t *testing.T) {
	testShardedIterator_Release(t, shardedDBCombineChanSize+10, func(db shardedDB) {
		// Next() returns True if Release() is not called
		{
			it := db.NewIteratorUnsorted(nil, nil)
			defer it.Release()

			// check if data exists
			for i := 0; i < shardedDBCombineChanSize+1; i++ {
				assert.True(t, it.Next())
			}
		}

		//  Next() returns False if Release() is called
		{
			it := db.NewIteratorUnsorted(nil, nil)
			it.Release()

			// flush data in channel
			for i := 0; i < shardedDBCombineChanSize+1; i++ {
				it.Next()
			}

			// check if Next returns false
			assert.False(t, it.Next())
		}
	})
}

func TestShardedDBParallelIterator_Release(t *testing.T) {
	testShardedIterator_Release(t,
		int(ShardedDBConfig[len(ShardedDBConfig)-1].NumStateTrieShards*shardedDBSubChannelSize*2),
		func(db shardedDB) {
			// Next() returns True if Release() is not called
			{
				it := db.NewParallelIterator(context.Background(), nil, nil, nil)
				defer it.Release()

				for _, ch := range it.Channels() {
					// check if channel is not closed
					for i := 0; i < shardedDBSubChannelSize+1; i++ {
						e, ok := <-ch
						assert.NotNil(t, e)
						assert.True(t, ok)
					}
				}
			}

			//  Next() returns False if Release() is called
			{
				it := db.NewParallelIterator(context.Background(), nil, nil, nil)
				it.Release()
				for _, ch := range it.Channels() {

					// flush data in channel
					for i := 0; i < shardedDBSubChannelSize+1; i++ {
						<-ch
					}

					// check if channel is closed
					_, ok := <-ch
					assert.False(t, ok)
				}
			}
		})
}
