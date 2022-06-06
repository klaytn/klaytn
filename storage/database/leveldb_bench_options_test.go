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

package database

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func genTempDirForTestDB(b *testing.B) string {
	dir, err := ioutil.TempDir("", "klaytn-db-bench")
	if err != nil {
		b.Fatalf("cannot create temporary directory: %v", err)
	}
	return dir
}

func getKlayLDBOptions() *opt.Options {
	return getLevelDBOptions(&DBConfig{LevelDBCacheSize: 128, OpenFilesLimit: 128})
}

func getKlayLDBOptionsForGetX(x int) *opt.Options {
	opts := getKlayLDBOptions()
	opts.WriteBuffer *= x
	opts.BlockCacheCapacity *= x
	opts.OpenFilesCacheCapacity *= x
	opts.DisableBlockCache = true

	return opts
}

func getKlayLDBOptionsForPutX(x int) *opt.Options {
	opts := getKlayLDBOptions()
	opts.BlockCacheCapacity *= x
	opts.BlockRestartInterval *= x

	opts.BlockSize *= x
	opts.CompactionExpandLimitFactor *= x
	opts.CompactionL0Trigger *= x
	opts.CompactionTableSize *= x

	opts.CompactionSourceLimitFactor *= x
	opts.Compression = opt.DefaultCompression

	return opts
}

func getKlayLDBOptionsForBatchX(x int) *opt.Options {
	opts := getKlayLDBOptions()
	opts.BlockCacheCapacity *= x
	opts.BlockRestartInterval *= x

	opts.BlockSize *= x
	opts.CompactionExpandLimitFactor *= x
	opts.CompactionL0Trigger *= x
	opts.CompactionTableSize *= x

	opts.CompactionSourceLimitFactor *= x
	opts.Compression = opt.DefaultCompression

	return opts
}

// readTypeFunc determines requested index
func benchmarkKlayOptionsGet(b *testing.B, opts *opt.Options, valueLength, numInsertions, numGets int, readTypeFunc func(int, int) int) {
	b.StopTimer()
	b.ReportAllocs()

	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	for i := 0; i < numInsertions; i++ {
		bs := []byte(strconv.Itoa(i))
		db.Put(bs, randStrBytes(valueLength))
	}

	b.StartTimer()
	for k := 0; k < b.N; k++ {
		for i := 0; i < numGets; i++ {
			bs := []byte(strconv.Itoa(readTypeFunc(i, numInsertions)))
			_, err := db.Get(bs)
			if err != nil {
				b.Fatalf("get failed: %v", err)
			}
		}
	}
}

func randomRead(currIndex, numInsertions int) int {
	return rand.Intn(numInsertions)
}

func sequentialRead(currIndex, numInsertions int) int {
	return numInsertions - currIndex - 1
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func zipfRead(currIndex, numInsertions int) int {
	zipf := rand.NewZipf(r, 3.14, 2.72, uint64(numInsertions))
	zipfNum := zipf.Uint64()
	return numInsertions - int(zipfNum) - 1
}

const (
	getKlayValueLegnth   = 250
	getKlayNumInsertions = 1000 * 100
	getKlaynumGets       = 1000
)

var getKlayOptions = [...]struct {
	name          string
	valueLength   int
	numInsertions int
	numGets       int
	opts          *opt.Options
}{
	{"X1", getKlayValueLegnth, getKlayNumInsertions, getKlaynumGets, getKlayLDBOptionsForGetX(1)},
	{"X2", getKlayValueLegnth, getKlayNumInsertions, getKlaynumGets, getKlayLDBOptionsForGetX(2)},
	{"X4", getKlayValueLegnth, getKlayNumInsertions, getKlaynumGets, getKlayLDBOptionsForGetX(4)},
	{"X8", getKlayValueLegnth, getKlayNumInsertions, getKlaynumGets, getKlayLDBOptionsForGetX(8)},
	//{"X16", getKlayValueLegnth, getKlayNumInsertions, getKlaynumGets, getKlayLDBOptionsForGetX(16)},
	//{"X32", getKlayValueLegnth, getKlayNumInsertions, getKlaynumGets, getKlayLDBOptionsForGetX(32)},
	//{"X64", getKlayValueLegnth, getKlayNumInsertions, getKlaynumGets, getKlayLDBOptionsForGetX(64)},
}

func Benchmark_KlayOptions_Get_Random(b *testing.B) {
	for _, bm := range getKlayOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKlayOptionsGet(b, bm.opts, bm.valueLength, bm.numInsertions, bm.numGets, randomRead)
		})
	}
}

func Benchmark_KlayOptions_Get_Sequential(b *testing.B) {
	for _, bm := range getKlayOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKlayOptionsGet(b, bm.opts, bm.valueLength, bm.numInsertions, bm.numGets, sequentialRead)
		})
	}
}

func Benchmark_KlayOptions_Get_Zipf(b *testing.B) {
	for _, bm := range getKlayOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKlayOptionsGet(b, bm.opts, bm.valueLength, bm.numInsertions, bm.numGets, zipfRead)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////// Put Insertion Tests Beginning //////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func benchmarkKlayOptionsPut(b *testing.B, opts *opt.Options, valueLength, numInsertions int) {
	b.StopTimer()

	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for k := 0; k < numInsertions; k++ {
			db.Put(randStrBytes(32), randStrBytes(valueLength))
		}
	}
}

func Benchmark_KlayOptions_Put(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	const (
		putKlayValueLegnth   = 250
		putKlayNumInsertions = 1000 * 10
	)

	putKlayOptions := [...]struct {
		name          string
		valueLength   int
		numInsertions int
		opts          *opt.Options
	}{
		{"X1", putKlayValueLegnth, putKlayNumInsertions, getKlayLDBOptionsForPutX(1)},
		{"X2", putKlayValueLegnth, putKlayNumInsertions, getKlayLDBOptionsForPutX(2)},
		{"X4", putKlayValueLegnth, putKlayNumInsertions, getKlayLDBOptionsForPutX(4)},
		{"X8", putKlayValueLegnth, putKlayNumInsertions, getKlayLDBOptionsForPutX(8)},
		//{"X16", putKlayValueLegnth, putKlayNumInsertions, getKlayLDBOptionsForPutX(16)},
		//{"X32", putKlayValueLegnth, putKlayNumInsertions, getKlayLDBOptionsForPutX(32)},
		//{"X64", putKlayValueLegnth, putKlayNumInsertions, getKlayLDBOptionsForPutX(64)},
	}
	for _, bm := range putKlayOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKlayOptionsPut(b, bm.opts, bm.valueLength, bm.numInsertions)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// Put Insertion Tests End /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////////////////
////////////////////////// SHARDED PUT INSERTION TESTS BEGINNING //////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func removeDirs(dirs []string) {
	for _, dir := range dirs {
		os.RemoveAll(dir)
	}
}

func genDatabases(b *testing.B, dirs []string, opts *opt.Options) []*levelDB {
	databases := make([]*levelDB, len(dirs), len(dirs))
	for i := 0; i < len(dirs); i++ {
		databases[i], _ = NewLevelDBWithOption(dirs[i], opts)
	}
	return databases
}

func closeDBs(databases []*levelDB) {
	for _, db := range databases {
		db.Close()
	}
}

func genKeysAndValues(valueLength, numInsertions int) ([][]byte, [][]byte) {
	keys := make([][]byte, numInsertions, numInsertions)
	values := make([][]byte, numInsertions, numInsertions)
	for i := 0; i < numInsertions; i++ {
		keys[i] = randStrBytes(32)
		values[i] = randStrBytes(valueLength)
	}
	return keys, values
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////// Batch Insertion Tests Beginning /////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func benchmarkKlayOptionsBatch(b *testing.B, opts *opt.Options, valueLength, numInsertions int) {
	b.StopTimer()
	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		keys, values := genKeysAndValues(valueLength, numInsertions)
		b.StartTimer()
		batch := db.NewBatch()
		for k := 0; k < numInsertions; k++ {
			batch.Put(keys[k], values[k])
		}
		batch.Write()
	}
}

func Benchmark_KlayOptions_Batch(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	const (
		batchValueLength       = 250
		batchKlayNumInsertions = 1000 * 10
	)

	putKlayOptions := [...]struct {
		name          string
		valueLength   int
		numInsertions int
		opts          *opt.Options
	}{
		{"X1", batchValueLength, batchKlayNumInsertions, getKlayLDBOptionsForBatchX(1)},
		{"X2", batchValueLength, batchKlayNumInsertions, getKlayLDBOptionsForBatchX(2)},
		{"X4", batchValueLength, batchKlayNumInsertions, getKlayLDBOptionsForBatchX(4)},
		{"X8", batchValueLength, batchKlayNumInsertions, getKlayLDBOptionsForBatchX(8)},
		//{"X16", batchValueLength, batchKlayNumInsertions, getKlayLDBOptionsForBatchX(16)},
		//{"X32", batchValueLength, batchKlayNumInsertions, getKlayLDBOptionsForBatchX(32)},
		//{"X64", batchValueLength, batchKlayNumInsertions, getKlayLDBOptionsForBatchX(64)},
	}
	for _, bm := range putKlayOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKlayOptionsBatch(b, bm.opts, bm.valueLength, bm.numInsertions)
		})
	}
}

// TODO-Klaytn = Add a test for checking GoRoutine Overhead

///////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////// Batch Insertion Tests End /////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////////////////
////////////////////////// Ideal Batch Size Tests Begins //////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

type idealBatchBM struct {
	name      string
	totalSize int
	batchSize int
	rowSize   int
}

func benchmarkIdealBatchSize(b *testing.B, bm idealBatchBM) {
	b.StopTimer()

	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	opts := getKlayLDBOptions()
	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	b.StartTimer()

	var wg sync.WaitGroup
	numBatches := bm.totalSize / bm.batchSize
	wg.Add(numBatches)
	for i := 0; i < numBatches; i++ {
		batch := db.NewBatch()
		for k := 0; k < bm.batchSize; k++ {
			batch.Put(randStrBytes(32), randStrBytes(bm.rowSize))
		}

		go func(currBatch Batch) {
			defer wg.Done()
			currBatch.Write()
		}(batch)
	}
	wg.Wait()
}

func Benchmark_IdealBatchSize(b *testing.B) {
	b.StopTimer()
	// please change below rowSize to change the size of an input row
	// key = 32 bytes, value = rowSize bytes
	const rowSize = 250

	benchmarks := []idealBatchBM{
		// to run test with total size smaller than 1,000 rows
		// go test -bench=Benchmark_IdealBatchSize/SmallBatches
		{"SmallBatches_100Rows_10Batch_250Bytes", 100, 10, rowSize},
		{"SmallBatches_100Rows_20Batch_250Bytes", 100, 20, rowSize},
		{"SmallBatches_100Rows_25Batch_250Bytes", 100, 25, rowSize},
		{"SmallBatches_100Rows_50Batch_250Bytes", 100, 50, rowSize},
		{"SmallBatches_100Rows_100Batch_250Bytes", 100, 100, rowSize},

		{"SmallBatches_200Rows_10Batch_250Bytes", 200, 10, rowSize},
		{"SmallBatches_200Rows_20Batch_250Bytes", 200, 20, rowSize},
		{"SmallBatches_200Rows_25Batch_250Bytes", 200, 25, rowSize},
		{"SmallBatches_200Rows_50Batch_250Bytes", 200, 50, rowSize},
		{"SmallBatches_200Rows_100Batch_250Bytes", 200, 100, rowSize},

		{"SmallBatches_400Rows_10Batch_250Bytes", 400, 10, rowSize},
		{"SmallBatches_400Rows_20Batch_250Bytes", 400, 20, rowSize},
		{"SmallBatches_400Rows_25Batch_250Bytes", 400, 25, rowSize},
		{"SmallBatches_400Rows_50Batch_250Bytes", 400, 50, rowSize},
		{"SmallBatches_400Rows_100Batch_250Bytes", 400, 100, rowSize},

		{"SmallBatches_800Rows_10Batch_250Bytes", 800, 10, rowSize},
		{"SmallBatches_800Rows_20Batch_250Bytes", 800, 20, rowSize},
		{"SmallBatches_800Rows_25Batch_250Bytes", 800, 25, rowSize},
		{"SmallBatches_800Rows_50Batch_250Bytes", 800, 50, rowSize},
		{"SmallBatches_800Rows_100Batch_250Bytes", 800, 100, rowSize},

		// to run test with total size between than 1k rows ~ 10k rows
		// go test -bench=Benchmark_IdealBatchSize/LargeBatches
		{"LargeBatches_1kRows_100Batch_250Bytes", 1000, 100, rowSize},
		{"LargeBatches_1kRows_200Batch_250Bytes", 1000, 200, rowSize},
		{"LargeBatches_1kRows_250Batch_250Bytes", 1000, 250, rowSize},
		{"LargeBatches_1kRows_500Batch_250Bytes", 1000, 500, rowSize},
		{"LargeBatches_1kRows_1000Batch_250Bytes", 1000, 1000, rowSize},

		{"LargeBatches_2kRows_100Batch_250Bytes", 2000, 100, rowSize},
		{"LargeBatches_2kRows_200Batch_250Bytes", 2000, 200, rowSize},
		{"LargeBatches_2kRows_250Batch_250Bytes", 2000, 250, rowSize},
		{"LargeBatches_2kRows_500Batch_250Bytes", 2000, 500, rowSize},
		{"LargeBatches_2kRows_1000Batch_250Bytes", 2000, 1000, rowSize},

		{"LargeBatches_4kRows_100Batch_250Bytes", 4000, 100, rowSize},
		{"LargeBatches_4kRows_200Batch_250Bytes", 4000, 200, rowSize},
		{"LargeBatches_4kRows_250Batch_250Bytes", 4000, 250, rowSize},
		{"LargeBatches_4kRows_500Batch_250Bytes", 4000, 500, rowSize},
		{"LargeBatches_4kRows_1000Batch_250Bytes", 4000, 1000, rowSize},

		{"LargeBatches_8kRows_100Batch_250Bytes", 8000, 100, rowSize},
		{"LargeBatches_8kRows_200Batch_250Bytes", 8000, 200, rowSize},
		{"LargeBatches_8kRows_250Batch_250Bytes", 8000, 250, rowSize},
		{"LargeBatches_8kRows_500Batch_250Bytes", 8000, 500, rowSize},
		{"LargeBatches_8kRows_1000Batch_250Bytes", 8000, 1000, rowSize},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for m := 0; m < b.N; m++ {
				benchmarkIdealBatchSize(b, bm)
			}
		})
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStrBytes(n int) []byte {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func getShardForTest(keys [][]byte, index, numShards int) int64 {
	return int64(index % numShards)
	// TODO-Klaytn: CHANGE BELOW LOGIC FROM ROUND-ROBIN TO USE getShardForTest
	//key := keys[index]
	//hashString := strings.TrimPrefix(common.Bytes2Hex(key),"0x")
	//if len(hashString) > 15 {
	//	hashString = hashString[:15]
	//}
	//seed, _ := strconv.ParseInt(hashString, 16, 64)
	//shard := seed % int64(numShards)
	//
	//return shard
}
