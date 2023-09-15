// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from ethdb/database_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package database

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func newTestLDB() (Database, func(), string) {
	dirName, err := ioutil.TempDir(os.TempDir(), "klay_leveldb_test_")
	if err != nil {
		panic("failed to create test file: " + err.Error())
	}
	db, err := NewLevelDBWithOption(dirName, GetDefaultLevelDBOption())
	if err != nil {
		panic("failed to create test database: " + err.Error())
	}

	return db, func() {
		db.Close()
		os.RemoveAll(dirName)
	}, "ldb"
}

func newTestBadgerDB() (Database, func(), string) {
	dirName, err := ioutil.TempDir(os.TempDir(), "klay_badgerdb_test_")
	if err != nil {
		panic("failed to create test file: " + err.Error())
	}
	db, err := NewBadgerDB(dirName)
	if err != nil {
		panic("failed to create test database: " + err.Error())
	}

	return db, func() {
		db.Close()
		os.RemoveAll(dirName)
	}, "badger"
}

func newTestMemDB() (Database, func(), string) {
	db := NewMemDB()
	return db, func() {
		db.Close()
	}, "memdb"
}

func newTestDynamoS3DB() (Database, func(), string) {
	// to start test with DynamoDB singletons
	oldDynamoDBClient := dynamoDBClient
	dynamoDBClient = nil

	oldDynamoOnceWorker := dynamoOnceWorker
	dynamoOnceWorker = &sync.Once{}

	oldDynamoWriteCh := dynamoWriteCh
	dynamoWriteCh = nil

	db, err := newDynamoDB(GetTestDynamoConfig())
	if err != nil {
		panic("failed to create test DynamoS3 database: " + err.Error())
	}
	return db, func() {
		db.Close()
		db.deleteTable()
		db.fdb.deleteBucket()

		// to finish test with DynamoDB singletons
		dynamoDBClient = oldDynamoDBClient
		dynamoOnceWorker = oldDynamoOnceWorker
		dynamoWriteCh = oldDynamoWriteCh
	}, "dynamos3db"
}

type commonDatabaseTestSuite struct {
	suite.Suite
	newFn    func() (Database, func(), string)
	removeFn func()
	database Database
}

var testDatabases []func() (Database, func(), string)

func TestDatabaseTestSuite(t *testing.T) {
	// If you want to include dynamo test, use below line
	// var testDatabases = []func() (Database, func()){newTestLDB, newTestBadgerDB, newTestMemDB, newTestDynamoS3DB}

	// TODO-Klaytn-Database Need to add DynamoDB to the below list.
	testDatabases = append(testDatabases, newTestLDB, newTestBadgerDB, newTestMemDB)
	for _, newFn := range testDatabases {
		suite.Run(t, &commonDatabaseTestSuite{newFn: newFn})
	}
}

func (ts *commonDatabaseTestSuite) BeforeTest(suiteName, testName string) {
	var dbname string
	ts.database, ts.removeFn, dbname = ts.newFn()
	ts.T().Logf("before test - dbname: %v, suiteName: %v, testName: %v", dbname, suiteName, testName)
}

func (ts *commonDatabaseTestSuite) AfterTest(suiteName, testName string) {
	ts.T().Logf("after test - suiteName: %v, testName: %v", suiteName, testName)
	ts.removeFn()
	ts.database, ts.removeFn = nil, nil
}

// TestNilValue checks if all database write/read nil value in the same way.
func (ts *commonDatabaseTestSuite) TestNilValue() {
	db, t := ts.database, ts.T()

	// non-batch
	{
		// write nil value
		key := common.MakeRandomBytes(32)
		assert.Nil(t, db.Put(key, nil))

		// get nil value
		ret, err := db.Get(key)
		assert.Equal(t, []byte{}, ret)
		assert.Nil(t, err)

		// check existence
		exist, err := db.Has(key)
		assert.Equal(t, true, exist)
		assert.Nil(t, err)

		val, err := db.Get(randStrBytes(100))
		assert.Nil(t, val)
		assert.Error(t, err)
		assert.Equal(t, dataNotFoundErr, err)
	}

	// batch
	{
		batch := db.NewBatch()

		// write nil value
		key := common.MakeRandomBytes(32)
		assert.Nil(t, batch.Put(key, nil))
		assert.NoError(t, batch.Write())

		// get nil value
		ret, err := db.Get(key)
		assert.Equal(t, []byte{}, ret)
		assert.Nil(t, err)

		// check existence
		exist, err := db.Has(key)
		assert.Equal(t, true, exist)
		assert.Nil(t, err)

		val, err := db.Get(randStrBytes(100))
		assert.Nil(t, val)
		assert.Error(t, err)
		assert.Equal(t, dataNotFoundErr, err)
	}
}

// TestNotFoundErr checks if an empty database returns DataNotFoundErr for the given random key.
func (ts *commonDatabaseTestSuite) TestNotFoundErr() {
	db, t := ts.database, ts.T()

	val, err := db.Get(randStrBytes(100))
	assert.Nil(t, val)
	assert.Error(t, err)
	assert.Equal(t, dataNotFoundErr, err)
}

// TestPutGet tests the basic put and get operations.
func (ts *commonDatabaseTestSuite) TestPutGet() {
	db, t := ts.database, ts.T()

	// Since badgerDB can't store empty key, testValues is modified. Below line is the original testValues.
	// var testValues = []string{"", "a", "1251", "\x00123\x00"}
	testValues := []string{"a", "1251", "\x00123\x00"}

	// put
	for _, v := range testValues {
		err := db.Put([]byte(v), []byte(v))
		if err != nil {
			t.Fatalf("put failed: %v", err)
		}
	}

	// get
	for _, v := range testValues {
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data, []byte(v)) {
			t.Fatalf("get returned wrong result, got %q expected %q", string(data), v)
		}
	}

	// override with "?"
	for _, v := range testValues {
		err := db.Put([]byte(v), []byte("?"))
		if err != nil {
			t.Fatalf("put override failed: %v", err)
		}
	}

	// get "?" by key
	for _, v := range testValues {
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data, []byte("?")) {
			t.Fatalf("get returned wrong result, got %q expected ?", string(data))
		}
	}

	// override returned value
	for _, v := range testValues {
		orig, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		orig[0] = byte(0xff)
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data, []byte("?")) {
			t.Fatalf("get returned wrong result, got %q expected ?", string(data))
		}
	}

	// delete
	for _, v := range testValues {
		err := db.Delete([]byte(v))
		if err != nil {
			t.Fatalf("delete %q failed: %v", v, err)
		}
	}

	// try to get deleted values
	for _, v := range testValues {
		_, err := db.Get([]byte(v))
		if err == nil {
			t.Fatalf("got deleted value %q", v)
		}
	}
}

func TestShardDB(t *testing.T) {
	key := common.Hex2Bytes("0x91d6f7d2537d8a0bd7d487dcc59151ebc00da306")

	hashstring := strings.TrimPrefix("0x93d6f3d2537d8a0bd7d485dcc59151ebc00da306", "0x")
	if len(hashstring) > 15 {
		hashstring = hashstring[:15]
	}
	seed, _ := strconv.ParseInt(hashstring, 16, 64)

	shard := seed % int64(12)

	idx := common.BytesToHash(key).Big().Mod(common.BytesToHash(key).Big(), big.NewInt(4))

	fmt.Printf("idx %d   %d   %d\n", idx, shard, seed)
}

// TestParallelPutGet tests the parallel put and get operations.
func (ts *commonDatabaseTestSuite) TestParallelPutGet() {
	db := ts.database
	const n = 8
	var pending sync.WaitGroup

	pending.Add(n)
	for i := 0; i < n; i++ {
		go func(key string) {
			defer pending.Done()
			err := db.Put([]byte(key), []byte("v"+key))
			if err != nil {
				panic("put failed: " + err.Error())
			}
		}(strconv.Itoa(i))
	}
	pending.Wait()

	pending.Add(n)
	for i := 0; i < n; i++ {
		go func(key string) {
			defer pending.Done()
			data, err := db.Get([]byte(key))
			if err != nil {
				panic("get failed: " + err.Error())
			}
			if !bytes.Equal(data, []byte("v"+key)) {
				panic(fmt.Sprintf("get failed, got %q expected %q", []byte(data), []byte("v"+key)))
			}
		}(strconv.Itoa(i))
	}
	pending.Wait()

	pending.Add(n)
	for i := 0; i < n; i++ {
		go func(key string) {
			defer pending.Done()
			err := db.Delete([]byte(key))
			if err != nil {
				panic("delete failed: " + err.Error())
			}
		}(strconv.Itoa(i))
	}
	pending.Wait()

	pending.Add(n)
	for i := 0; i < n; i++ {
		go func(key string) {
			defer pending.Done()
			_, err := db.Get([]byte(key))
			if err == nil {
				panic("got deleted value")
			}
		}(strconv.Itoa(i))
	}
	pending.Wait()
}

// TestDBEntryLengthCheck checks if dbDirs and dbConfigRatio are
// specified for every DBEntryType.
func TestDBEntryLengthCheck(t *testing.T) {
	dbRatioSum := 0
	for i := 0; i < int(databaseEntryTypeSize); i++ {
		if dbBaseDirs[i] == "" {
			t.Fatalf("Database directory should be specified! index: %v", i)
		}

		if dbConfigRatio[i] == 0 {
			t.Fatalf("Database configuration ratio should be specified! index: %v", i)
		}

		dbRatioSum += dbConfigRatio[i]
	}

	if dbRatioSum != 100 {
		t.Fatalf("Sum of database configuration ratio should be 100! actual: %v", dbRatioSum)
	}
}

type testData struct {
	k, v []byte
}

type testDataSlice []*testData

func (d testDataSlice) Len() int {
	return len(d)
}

func (d testDataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d testDataSlice) Less(i, j int) bool {
	return bytes.Compare(d[i].k, d[j].k) < 0
}

func insertRandomData(db KeyValueWriter, prefix []byte, num int) (testDataSlice, error) {
	ret := testDataSlice{}
	for i := 0; i < num; i++ {
		key := common.MakeRandomBytes(32)
		val := append(key, key...)
		if len(prefix) > 0 {
			key = append(prefix, key...)
		}
		ret = append(ret, &testData{k: key, v: val})

		if err := db.Put(key, val); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (ts *commonDatabaseTestSuite) Test_Put() {
	num := 100

	// test put
	data, err := insertRandomData(ts.database, nil, num)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), num, len(data))
}

func (ts *commonDatabaseTestSuite) Test_Get() {
	num, db := 100, ts.database

	data, _ := insertRandomData(ts.database, nil, num)

	for idx := range data {
		actual, err := db.Get(data[idx].k)
		assert.NoError(ts.T(), err)
		assert.Equal(ts.T(), data[idx].v, actual)
	}

	notExistKey := []byte{0x1}
	actual, err := db.Get(notExistKey)
	assert.Equal(ts.T(), dataNotFoundErr, err)
	assert.Nil(ts.T(), actual)
}

func (ts *commonDatabaseTestSuite) Test_Has() {
	num, db := 100, ts.database

	data, _ := insertRandomData(ts.database, nil, num)

	for idx := range data {
		has, err := db.Has(data[idx].k)
		assert.NoError(ts.T(), err)
		assert.True(ts.T(), has)
	}

	notExistKey := []byte{0x1}
	has, err := db.Has(notExistKey)
	assert.NoError(ts.T(), err)
	assert.False(ts.T(), has)
}

func (ts *commonDatabaseTestSuite) Test_Delete() {
	num, db := 100, ts.database

	data, _ := insertRandomData(ts.database, nil, num)

	for idx := range data {
		if idx%2 == 0 {
			err := db.Delete(data[idx].k)
			assert.NoError(ts.T(), err)
		}
	}

	for idx := range data {
		has, _ := db.Has(data[idx].k)
		if idx%2 == 0 {
			assert.False(ts.T(), has)
		} else {
			assert.True(ts.T(), has)
		}
	}
}

func (ts *commonDatabaseTestSuite) Test_Iterator_NoData() {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	db := ts.database
	if _, ok := db.(*badgerDB); ok {
		ts.T().Skip()
	}

	// testing iterator without prefix nor specific-starting key
	it := db.NewIterator(nil, nil)
	defer it.Release()

	assert.False(ts.T(), it.Next())
}

func (ts *commonDatabaseTestSuite) Test_Iterator_WithoutPrefixAndStart() {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	num, db := 100, ts.database
	if _, ok := db.(*badgerDB); ok {
		ts.T().Skip()
	}

	data, _ := insertRandomData(ts.database, nil, num)
	sort.Sort(data)

	// testing iterator without prefix nor specific-starting key
	it := db.NewIterator(nil, nil)
	defer it.Release()

	idx := 0
	for it.Next() {
		key, val := it.Key(), it.Value()
		assert.Equal(ts.T(), data[idx].k, key)
		assert.Equal(ts.T(), data[idx].v, val)
		idx++
	}
	assert.Equal(ts.T(), len(data), idx)
}

func (ts *commonDatabaseTestSuite) Test_Iterator_WithPrefix() {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	num, prefix, db := 10, common.Hex2Bytes("deaddeaf"), ts.database
	if _, ok := db.(*badgerDB); ok {
		ts.T().Skip()
	}

	insertRandomData(ts.database, nil, num)
	prefixData, _ := insertRandomData(ts.database, prefix, num)
	sort.Sort(prefixData)

	// testing iterator with key-prefix
	it := db.NewIterator(prefix, nil)
	defer it.Release()
	assert.Nil(ts.T(), it.Key())
	assert.Nil(ts.T(), it.Value())

	idx := 0
	for it.Next() {
		key, val := it.Key(), it.Value()
		assert.Equal(ts.T(), prefixData[idx].k, key)
		assert.Equal(ts.T(), prefixData[idx].v, val)
		idx++
	}
	assert.Equal(ts.T(), len(prefixData), idx)
}

func (ts *commonDatabaseTestSuite) Test_Iterator_WithStart() {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	num, db := 100, ts.database
	if _, ok := db.(*badgerDB); ok {
		ts.T().Skip()
	}

	data, _ := insertRandomData(ts.database, nil, num)
	sort.Sort(data)

	startIdx := len(data) / 3

	// testing iterator with specific starting key
	it := db.NewIterator(nil, data[startIdx].k)
	defer it.Release()
	assert.Nil(ts.T(), it.Key())
	assert.Nil(ts.T(), it.Value())

	idx := startIdx
	for it.Next() {
		key, val := it.Key(), it.Value()
		assert.Equal(ts.T(), data[idx].k, key)
		assert.Equal(ts.T(), data[idx].v, val)
		idx++
	}
	assert.Equal(ts.T(), len(data), idx)
}

func (ts *commonDatabaseTestSuite) Test_Iterator_WithPrefixAndStart() {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	num, prefix, db := 10, common.Hex2Bytes("deaddeaf"), ts.database
	if _, ok := db.(*badgerDB); ok {
		ts.T().Skip()
	}

	insertRandomData(ts.database, common.Hex2Bytes("aaaabbbb"), num)
	data, _ := insertRandomData(ts.database, prefix, num)
	sort.Sort(data)

	startIdx := len(data) / 3

	// testing iterator with prefix and specific-starting key
	it := db.NewIterator(prefix, data[startIdx].k[4:])
	defer it.Release()
	assert.Nil(ts.T(), it.Key())
	assert.Nil(ts.T(), it.Value())

	idx := startIdx
	for it.Next() {
		key, val := it.Key(), it.Value()
		assert.Equal(ts.T(), data[idx].k, key)
		assert.Equal(ts.T(), data[idx].v, val)
		idx++
	}
	assert.Equal(ts.T(), len(data), idx)
}

func (ts *commonDatabaseTestSuite) Test_BatchWrite() {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	numData, numIter, db := 1000, 100, ts.database
	if _, ok := db.(*badgerDB); ok {
		ts.T().Skip()
	}

	batch := db.NewBatch()
	defer batch.Release()
	data := testDataSlice{}
	for i := 0; i < numIter; i++ {
		inserted, _ := insertRandomData(batch, nil, numData)
		batch.Write()
		batch.Reset()
		data = append(data, inserted...)
	}

	assert.Equal(ts.T(), numData*numIter, len(data))
	for _, d := range data {
		actual, err := db.Get(d.k)
		assert.NoError(ts.T(), err)
		assert.Equal(ts.T(), d.v, actual)
	}
}
