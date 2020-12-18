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
	{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 4, ParallelDBWrite: true},
	{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 8, ParallelDBWrite: true},
	{DBType: LevelDB, SingleDB: false, NumStateTrieShards: 16, ParallelDBWrite: true},
}

func createEntries(entryNum int) []entry {
	entries := make([]entry, entryNum)
	for i := 0; i < entryNum; i++ {
		//entries[i] = entry{common.MakeRandomBytes(10), common.MakeRandomBytes(10)}
		entries[i] = entry{common.MakeRandomBytes(256), common.MakeRandomBytes(600)}
	}
	return entries
}

// testIterator tests if given iterator iterates all entries
func testIterator(t *testing.T, checkOrder bool, entriesFromIterator func(db shardedDB, entryNum int) []entry) {
	entryNum := 500
	entries := createEntries(entryNum)
	dbs := make([]shardedDB, len(ShardedDBConfig))

	// create DB and write data for testing
	for i, config := range ShardedDBConfig {
		config.Dir, _ = ioutil.TempDir(os.TempDir(), "test-shardedDB-iterator")
		defer func(dir string) {
			if err := os.RemoveAll(dir); err != nil {
				t.Fatalf("fail to delete file %v", err)
			}
		}(config.Dir)

		// create sharded DB
		db, err := newShardedDB(config, 0, config.NumStateTrieShards)
		dbs[i] = *db
		if err != nil {
			t.Log("Error occured while creating DB")
			t.FailNow()
		}

		// write entries data in DB
		batch := db.NewBatch()
		for _, entry := range entries {
			batch.Put(entry.key, entry.val)
		}
		assert.NoError(t, batch.Write())
	}

	// sort entries for each compare
	sort.Slice(entries, func(i, j int) bool { return bytes.Compare(entries[i].key, entries[j].key) < 0 })

	// get data from iterator and compare
	for _, db := range dbs {
		// create iterator
		entriesFromIt := entriesFromIterator(db, entryNum)
		if !checkOrder {
			sort.Slice(entriesFromIt, func(i, j int) bool { return bytes.Compare(entriesFromIt[i].key, entriesFromIt[j].key) < 0 })
		}

		// compare if entries generated and entries from iterator is same
		assert.Equal(t, len(entries), len(entriesFromIt))
		assert.Equal(t, entries[0], entriesFromIt[0])
		assert.True(t, reflect.DeepEqual(entries, entriesFromIt))
	}
}

// TestShardedDBChanIterator tests if shardedDBIterator iterates all entries
// TODO implement TestShardedDBIteratorWithStart and TestShardedDBIteratorWithPrefix
func TestShardedDBIterator(t *testing.T) {
	testIterator(t, true, func(db shardedDB, entryNum int) []entry {
		entries := make([]entry, 0, entryNum)
		it := db.NewIterator()

		for it.Next() {
			entries = append(entries, entry{it.Key(), it.Value()})
		}
		it.Release()
		assert.NoError(t, it.Error())
		return entries
	})
}

// TestShardedDBChanIterator tests if shardedDBChanIterator iterates all entries
// TODO implement TestShardedDBChanIteratorWithStart and TestShardedDBChanIteratorWithStartWitPrefix
func TestShardedDBChanIterator(t *testing.T) {
	testIterator(t, false, func(db shardedDB, entryNum int) []entry {
		entries := make([]entry, 0, entryNum) // store all items
		var l sync.RWMutex                    // mutex for entries

		// create chan Iterator and get channels
		it := db.NewChanIterator(context.Background(), func(db Database) Iterator { return db.NewIterator() })
		chans := it.Channels()

		// listen all channels and get key/value
		done := make(chan struct{})
		for _, ch := range chans {
			go func(ch chan entry) {
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
	})
}

// TestShardedDBChanIterator tests if shardedDBIteratorUnsorted iterates all entries
// TODO implement TestShardedDBIteratorWithStartUnsorted and TestShardedDBIteratorWithPrefixUnsorted
func TestShardedDBIteratorUnsorted(t *testing.T) {
	testIterator(t, false, func(db shardedDB, entryNum int) []entry {
		entries := make([]entry, 0, entryNum)
		it := db.NewIteratorUnsorted()

		for it.Next() {
			entries = append(entries, entry{it.Key(), it.Value()})
		}
		it.Release()
		assert.NoError(t, it.Error())
		return entries
	})
}
