// Copyright 2018 The klaytn Authors
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
//go:build multidisktest
// +build multidisktest

package database

import (
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type multiDiskOption struct {
	numDisks  int
	numShards int
	numRows   int
	numBytes  int
}

// adjust diskPaths according to your setting
var diskPaths = [...]string{"", "/disk1", "/disk2", "/disk3", "/disk4", "/disk5", "/disk6", "/disk7"}

func genDirForMDPTest(b *testing.B, numDisks, numShards int) []string {
	if numDisks > len(diskPaths) {
		b.Fatalf("entered numDisks %v is larger than diskPaths %v", numDisks, len(diskPaths))
	}

	dirs := make([]string, numShards, numShards)
	for i := 0; i < numShards; i++ {
		diskNum := i % numDisks
		dir, err := ioutil.TempDir(diskPaths[diskNum], "klaytn-db-bench-mdp")
		if err != nil {
			b.Fatalf("cannot create temporary directory: %v", err)
		}
		dirs[i] = dir
	}
	return dirs
}

func benchmark_MDP_Put_GoRoutine(b *testing.B, mdo *multiDiskOption) {
	b.StopTimer()

	dirs := genDirForMDPTest(b, mdo.numDisks, mdo.numShards)
	defer removeDirs(dirs)

	opts := getKlayLDBOptions()
	databases := genDatabases(b, dirs, opts)
	defer closeDBs(databases)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		keys, values := genKeysAndValues(mdo.numBytes, mdo.numRows)
		b.StartTimer()

		var wait sync.WaitGroup
		wait.Add(mdo.numRows)

		for k := 0; k < mdo.numRows; k++ {
			shard := getShardForTest(keys, k, mdo.numShards)
			db := databases[shard]

			go func(currDB Database, idx int) {
				defer wait.Done()
				currDB.Put(keys[idx], values[idx])
			}(db, k)
		}
		wait.Wait()
	}
}

func benchmark_MDP_Put_NoGoRoutine(b *testing.B, mdo *multiDiskOption) {
	b.StopTimer()

	numDisks := mdo.numDisks
	numBytes := mdo.numBytes
	numRows := mdo.numRows
	numShards := mdo.numShards

	dirs := genDirForMDPTest(b, numDisks, numShards)
	defer removeDirs(dirs)

	opts := getKlayLDBOptions()
	databases := genDatabases(b, dirs, opts)
	defer closeDBs(databases)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		keys, values := genKeysAndValues(numBytes, numRows)
		b.StartTimer()

		for k := 0; k < numRows; k++ {
			shard := getShardForTest(keys, k, mdo.numShards)
			db := databases[shard]
			db.Put(keys[k], values[k])
		}
	}
}

// please change below rowSize to change the size of an input row for MDP_Put tests (GoRoutine & NoGoRoutine)
const rowSizePutMDP = 250
const numRowsPutMDP = 1000 * 10

var putMDPBenchmarks = [...]struct {
	name string
	mdo  multiDiskOption
}{
	// 10k Rows, 250 bytes, 1 disk, different number of shards
	{"10k_1D_1P_250", multiDiskOption{1, 1, numRowsPutMDP, rowSizePutMDP}},
	{"10k_1D_2P_250", multiDiskOption{1, 2, numRowsPutMDP, rowSizePutMDP}},
	{"10k_1D_4P_250", multiDiskOption{1, 4, numRowsPutMDP, rowSizePutMDP}},
	{"10k_1D_8P_250", multiDiskOption{1, 8, numRowsPutMDP, rowSizePutMDP}},

	// 10k Rows, 250 bytes, 8 shards (fixed), different number of disks
	{"10k_1D_8P_250", multiDiskOption{1, 8, numRowsPutMDP, rowSizePutMDP}},
	{"10k_2D_8P_250", multiDiskOption{2, 8, numRowsPutMDP, rowSizePutMDP}},
	{"10k_4D_8P_250", multiDiskOption{4, 8, numRowsPutMDP, rowSizePutMDP}},
	{"10k_8D_8P_250", multiDiskOption{8, 8, numRowsPutMDP, rowSizePutMDP}},

	// 10k Rows, 250 bytes, different number of disks & shards
	{"10k_1D_1P_250", multiDiskOption{1, 1, numRowsPutMDP, rowSizePutMDP}},
	{"10k_2D_2P_250", multiDiskOption{2, 2, numRowsPutMDP, rowSizePutMDP}},
	{"10k_4D_4P_250", multiDiskOption{4, 4, numRowsPutMDP, rowSizePutMDP}},
	{"10k_8D_8P_250", multiDiskOption{8, 8, numRowsPutMDP, rowSizePutMDP}},
}

func Benchmark_MDP_Put_GoRoutine(b *testing.B) {
	for _, bm := range putMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_MDP_Put_GoRoutine(b, &bm.mdo)
		})
	}
}

func Benchmark_MDP_Put_NoGoRoutine(b *testing.B) {
	for _, bm := range putMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_MDP_Put_NoGoRoutine(b, &bm.mdo)
		})
	}
}

func benchmark_MDP_Batch_GoRoutine(b *testing.B, mdo *multiDiskOption) {
	b.StopTimer()

	numDisks := mdo.numDisks
	numBytes := mdo.numBytes
	numRows := mdo.numRows
	numShards := mdo.numShards

	dirs := genDirForMDPTest(b, numDisks, numShards)
	defer removeDirs(dirs)

	opts := getKlayLDBOptions()
	databases := genDatabases(b, dirs, opts)
	defer closeDBs(databases)

	zeroSizeBatch := 0
	batchSizeSum := 0
	numBatches := 0
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// make same number of batches as numShards
		batches := make([]Batch, numShards, numShards)
		for k := 0; k < numShards; k++ {
			batches[k] = databases[k].NewBatch()
		}
		keys, values := genKeysAndValues(numBytes, numRows)
		b.StartTimer()
		for k := 0; k < numRows; k++ {
			shard := getShardForTest(keys, k, numShards)
			batches[shard].Put(keys[k], values[k])
		}

		for _, batch := range batches {
			if batch.ValueSize() == 0 {
				zeroSizeBatch++
			}
			batchSizeSum += batch.ValueSize()
			numBatches++
		}
		var wait sync.WaitGroup
		wait.Add(numShards)
		for _, batch := range batches {
			go func(currBatch Batch) {
				defer wait.Done()
				currBatch.Write()
			}(batch)
		}
		wait.Wait()
	}

	if zeroSizeBatch != 0 {
		b.Log("zeroSizeBatch: ", zeroSizeBatch)
	}
}

func benchmark_MDP_Batch_NoGoRoutine(b *testing.B, mdo *multiDiskOption) {
	b.StopTimer()

	numDisks := mdo.numDisks
	numBytes := mdo.numBytes
	numRows := mdo.numRows
	numShards := mdo.numShards

	dirs := genDirForMDPTest(b, numDisks, numShards)
	defer removeDirs(dirs)

	opts := getKlayLDBOptions()
	databases := genDatabases(b, dirs, opts)
	defer closeDBs(databases)

	zeroSizeBatch := 0
	batchSizeSum := 0
	numBatches := 0
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// make same number of batches as numShards
		batches := make([]Batch, numShards, numShards)
		for k := 0; k < numShards; k++ {
			batches[k] = databases[k].NewBatch()
		}
		keys, values := genKeysAndValues(numBytes, numRows)
		b.StartTimer()
		for k := 0; k < numRows; k++ {
			shard := getShardForTest(keys, k, numShards)
			batches[shard].Put(keys[k], values[k])
		}

		for _, batch := range batches {
			if batch.ValueSize() == 0 {
				zeroSizeBatch++
			}
			batchSizeSum += batch.ValueSize()
			numBatches++
		}

		for _, batch := range batches {
			batch.Write()
		}

	}

	if zeroSizeBatch != 0 {
		b.Log("zeroSizeBatch: ", zeroSizeBatch)
	}
}

// please change below rowSize to change the size of an input row for MDP_Batch tests (GoRoutine & NoGoRoutine)
const rowSizeBatchMDP = 250
const numRowsBatchMDP = 1000 * 10

var batchMDPBenchmarks = [...]struct {
	name string
	mdo  multiDiskOption
}{
	// 10k Rows, 250 bytes, 1 disk, different number of shards
	{"10k_1D_1P_250", multiDiskOption{1, 1, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_1D_2P_250", multiDiskOption{1, 2, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_1D_4P_250", multiDiskOption{1, 4, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_1D_8P_250", multiDiskOption{1, 8, numRowsBatchMDP, rowSizeBatchMDP}},

	// 10k Rows, 250 bytes, 8 shards (fixed), different number of disks
	{"10k_1D_8P_250", multiDiskOption{1, 8, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_2D_8P_250", multiDiskOption{2, 8, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_4D_8P_250", multiDiskOption{4, 8, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_8D_8P_250", multiDiskOption{8, 8, numRowsBatchMDP, rowSizeBatchMDP}},

	// 10k Rows, 250 bytes, different number of disks & shards
	{"10k_1D_1P_250", multiDiskOption{1, 1, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_2D_2P_250", multiDiskOption{2, 2, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_4D_4P_250", multiDiskOption{4, 4, numRowsBatchMDP, rowSizeBatchMDP}},
	{"10k_8D_8P_250", multiDiskOption{8, 8, numRowsBatchMDP, rowSizeBatchMDP}},
}

func Benchmark_MDP_Batch_GoRoutine(b *testing.B) {
	for _, bm := range batchMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_MDP_Batch_GoRoutine(b, &bm.mdo)
		})
	}
}

func Benchmark_MDP_Batch_NoGoRoutine(b *testing.B) {
	for _, bm := range batchMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_MDP_Batch_NoGoRoutine(b, &bm.mdo)
		})
	}
}

func benchmark_MDP_Get_NoGoRotine(b *testing.B, mdo *multiDiskOption, numReads int, readType func(int, int) int) {
	b.StopTimer()

	numDisks := mdo.numDisks
	numBytes := mdo.numBytes
	numRows := mdo.numRows
	numShards := mdo.numShards

	dirs := genDirForMDPTest(b, numDisks, numShards)
	defer removeDirs(dirs)

	opts := getKlayLDBOptions()
	databases := genDatabases(b, dirs, opts)
	defer closeDBs(databases)

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		keys, values := genKeysAndValues(numBytes, numRows)

		for k := 0; k < numRows; k++ {
			shard := getShardForTest(keys, k, numShards)
			db := databases[shard]
			db.Put(keys[k], values[k])
		}

		b.StartTimer()
		for k := 0; k < numReads; k++ {
			keyPos := readType(k, numRows)
			if keyPos >= len(keys) {
				b.Fatal("index out of range", keyPos)
			}
			shard := getShardForTest(keys, k, numShards)
			db := databases[shard]
			db.Get(keys[keyPos])
		}
	}
}

func benchmark_MDP_Get_GoRoutine(b *testing.B, mdo *multiDiskOption, numReads int, readType func(int, int) int) {
	b.StopTimer()

	numDisks := mdo.numDisks
	numBytes := mdo.numBytes
	numRows := mdo.numRows
	numShards := mdo.numShards

	dirs := genDirForMDPTest(b, numDisks, numShards)
	defer removeDirs(dirs)

	opts := getKlayLDBOptions()
	databases := genDatabases(b, dirs, opts)
	defer closeDBs(databases)

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		keys, values := genKeysAndValues(numBytes, numRows)

		for k := 0; k < numRows; k++ {
			shard := getShardForTest(keys, k, numShards)
			db := databases[shard]
			db.Put(keys[k], values[k])
		}

		b.StartTimer()
		var wg sync.WaitGroup
		wg.Add(numReads)
		for k := 0; k < numReads; k++ {
			keyPos := readType(k, numRows)
			if keyPos >= len(keys) {
				b.Fatalf("index out of range: keyPos: %v, k: %v, numRows: %v", keyPos, k, numRows)
			}

			shard := getShardForTest(keys, keyPos, numShards)
			db := databases[shard]

			go func(currDB Database, kPos int) {
				defer wg.Done()
				_, err := currDB.Get(keys[kPos])
				if err != nil {
					b.Fatalf("get failed: %v", err)
				}
			}(db, keyPos)

		}
		wg.Wait()
	}
}

// please change below rowSize to change the size of an input row for MDP_Get tests (GoRoutine & NoGoRoutine)
const rowSizeGetMDP = 250

const (
	insertedRowsBeforeGetMDP = 1000 * 100 // pre-insertion size before read
	numReadsMDP              = 1000
)

var getMDPBenchmarks = [...]struct {
	name     string
	mdo      multiDiskOption
	numReads int
}{
	// 10k Rows, 250 bytes, 1 disk, different number of shards
	{"10k_1D_1P_250", multiDiskOption{1, 1, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_1D_2P_250", multiDiskOption{1, 2, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_1D_4P_250", multiDiskOption{1, 4, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_1D_8P_250", multiDiskOption{1, 8, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},

	// 10k Rows, 250 bytes, 8 shards (fixed), different number of disks
	{"10k_1D_8P_250", multiDiskOption{1, 8, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_2D_8P_250", multiDiskOption{2, 8, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_4D_8P_250", multiDiskOption{4, 8, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_8D_8P_250", multiDiskOption{8, 8, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},

	// 10k Rows, 250 bytes, different number of disks & shards
	{"10k_1D_1P_250", multiDiskOption{1, 1, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_2D_2P_250", multiDiskOption{2, 2, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_4D_4P_250", multiDiskOption{4, 4, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
	{"10k_8D_8P_250", multiDiskOption{8, 8, insertedRowsBeforeGetMDP, rowSizeGetMDP}, numReadsMDP},
}

func Benchmark_MDP_Get_Random_1kRows_GoRoutine(b *testing.B) {
	for _, bm := range getMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_MDP_Get_GoRoutine(b, &bm.mdo, bm.numReads, randomRead)
		})
	}
}

func Benchmark_MDP_Get_Random_1kRows_NoGoRoutine(b *testing.B) {
	for _, bm := range getMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_MDP_Get_NoGoRotine(b, &bm.mdo, bm.numReads, randomRead)
		})
	}
}

func Benchmark_MDP_Parallel_Get(b *testing.B) {
	for _, bm := range getMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.StopTimer()
			mdo := bm.mdo
			numDisks := mdo.numDisks
			numBytes := mdo.numBytes
			numRows := mdo.numRows
			numShards := mdo.numShards

			dirs := genDirForMDPTest(b, numDisks, numShards)
			defer removeDirs(dirs)

			opts := getKlayLDBOptions()
			databases := genDatabases(b, dirs, opts)
			defer closeDBs(databases)

			keys, values := genKeysAndValues(numBytes, numRows)

			for k := 0; k < numRows; k++ {
				shard := getShardForTest(keys, k, numShards)
				db := databases[shard]
				db.Put(keys[k], values[k])
			}

			rand.Seed(time.Now().UnixNano())
			b.StartTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					idx := rand.Intn(numRows)
					shard := getShardForTest(keys, idx, numShards)
					_, err := databases[shard].Get(keys[idx])
					if err != nil {
						b.Fatalf("get failed: %v", err)
					}
				}
			})
		})
	}
}

func Benchmark_MDP_Parallel_Put(b *testing.B) {
	for _, bm := range putMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.StopTimer()
			mdo := bm.mdo
			numDisks := mdo.numDisks
			numBytes := mdo.numBytes
			numRows := mdo.numRows * 10 // extend candidate keys and values for randomness
			numShards := mdo.numShards

			dirs := genDirForMDPTest(b, numDisks, numShards)
			defer removeDirs(dirs)

			opts := getKlayLDBOptions()
			databases := genDatabases(b, dirs, opts)
			defer closeDBs(databases)

			keys, values := genKeysAndValues(numBytes, numRows)

			rand.Seed(time.Now().UnixNano())
			b.StartTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					idx := rand.Intn(numRows)
					shard := getShardForTest(keys, idx, numShards)
					db := databases[shard]
					db.Put(keys[idx], values[idx])
				}
			})
		})
	}
}

const parallelBatchSizeMDP = 100

func Benchmark_MDP_Parallel_Batch(b *testing.B) {
	for _, bm := range batchMDPBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.StopTimer()
			mdo := bm.mdo
			numDisks := mdo.numDisks
			numBytes := mdo.numBytes
			numRows := mdo.numRows * 10 // extend candidate keys and values for randomness
			numShards := mdo.numShards

			dirs := genDirForMDPTest(b, numDisks, numShards)
			defer removeDirs(dirs)

			opts := getKlayLDBOptions()
			databases := genDatabases(b, dirs, opts)
			defer closeDBs(databases)

			keys, values := genKeysAndValues(numBytes, numRows)

			rand.Seed(time.Now().UnixNano())
			b.StartTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					shard := rand.Intn(numShards)
					batch := databases[shard].NewBatch()
					for k := 0; k < parallelBatchSizeMDP; k++ {
						idx := rand.Intn(numRows)
						batch.Put(keys[idx], values[idx])
					}
					batch.Write()
				}
			})
		})
	}
}

// TODO-Klaytn: MAKE PRE-LOADED TEST FOR BATCH, PUT
