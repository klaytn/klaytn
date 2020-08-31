// Copyright 2020 The klaytn Authors
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
//
// You need to set AWS credentials to access to dynamoDB.
//    sh$ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
//    sh$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET

package database

import (
	"io"
	"os"
	"testing"

	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/log/term"
	"github.com/mattn/go-colorable"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func enableLog() {
	usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	glogger := log.NewGlogHandler(log.StreamHandler(output, log.TerminalFormat(usecolor)))
	log.PrintOrigins(true)
	log.ChangeGlobalLogLevel(glogger, log.Lvl(5))
	glogger.Vmodule("")
	glogger.BacktraceAt("")
	log.Root().SetHandler(glogger)
}

func testDynamoDB_Put(t *testing.T) {
	dynamo, err := NewDynamoDB(GetDefaultDynamoDBConfig())
	defer dynamo.deleteDB()
	if err != nil {
		t.Fatal(err)
	}
	testKey := common.MakeRandomBytes(32)
	testVal := common.MakeRandomBytes(500)

	val, err := dynamo.Get(testKey)

	assert.Nil(t, val)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), dataNotFoundErr.Error())

	assert.NoError(t, dynamo.Put(testKey, testVal))
	returnedVal, returnedErr := dynamo.Get(testKey)
	assert.Equal(t, testVal, returnedVal)
	assert.NoError(t, returnedErr)
}

func testDynamoBatch_Write(t *testing.T) {
	dynamo, err := NewDynamoDB(GetDefaultDynamoDBConfig())
	defer dynamo.deleteDB()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dynamoDB", dynamo.config.TableName)

	var testKeys [][]byte
	var testVals [][]byte
	batch := dynamo.NewBatch()

	itemNum := 25
	for i := 0; i < itemNum; i++ {
		testKey := common.MakeRandomBytes(32)
		testVal := common.MakeRandomBytes(500)

		testKeys = append(testKeys, testKey)
		testVals = append(testVals, testVal)

		assert.NoError(t, batch.Put(testKey, testVal))
	}
	assert.NoError(t, batch.Write())

	// check if exist
	for i := 0; i < itemNum; i++ {
		returnedVal, returnedErr := dynamo.Get(testKeys[i])
		assert.NoError(t, returnedErr)
		assert.Equal(t, hexutil.Encode(testVals[i]), hexutil.Encode(returnedVal))
	}
}

func testDynamoBatch_WriteLargeData(t *testing.T) {
	dynamo, err := NewDynamoDB(GetDefaultDynamoDBConfig())
	defer dynamo.deleteDB()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dynamoDB", dynamo.config.TableName)

	var testKeys [][]byte
	var testVals [][]byte
	batch := dynamo.NewBatch()

	itemNum := 26
	for i := 0; i < itemNum; i++ {
		testKey := common.MakeRandomBytes(32)
		testVal := common.MakeRandomBytes(500 * 1024)

		testKeys = append(testKeys, testKey)
		testVals = append(testVals, testVal)

		assert.NoError(t, batch.Put(testKey, testVal))
	}
	assert.NoError(t, batch.Write())

	// check if exist
	for i := 0; i < itemNum; i++ {
		returnedVal, returnedErr := dynamo.Get(testKeys[i])
		assert.NoError(t, returnedErr)
		assert.Equal(t, hexutil.Encode(testVals[i]), hexutil.Encode(returnedVal))
	}
}

// testDynamoBatch_WriteMutliTables checks there is no error when working with more than on tables.
// This also checks if shared workers works as expected.
func testDynamoBatch_WriteMutliTables(t *testing.T) {
	// this test might end with Crit, enableLog to find out the log
	//enableLog()

	// creat DynamoDB1
	dynamo, err := NewDynamoDB(GetDefaultDynamoDBConfig())
	defer dynamo.deleteDB()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dynamoDB1", dynamo.config.TableName)

	// creat DynamoDB2
	dynamo2, err := NewDynamoDB(GetDefaultDynamoDBConfig())
	defer dynamo2.deleteDB()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dynamoDB2", dynamo2.config.TableName)

	var testKeys, testVals [][]byte
	var testKeys2, testVals2 [][]byte

	batch := dynamo.NewBatch()
	batch2 := dynamo2.NewBatch()

	itemNum := WorkerNum * 2
	for i := 0; i < itemNum; i++ {
		// write batch items to db1 and db2 in turn
		for j := 0; j < dynamoBatchSize; j++ {
			// write key and val to db1
			testKey := common.MakeRandomBytes(10)
			testVal := common.MakeRandomBytes(20)

			testKeys = append(testKeys, testKey)
			testVals = append(testVals, testVal)

			assert.NoError(t, batch.Put(testKey, testVal))

			// write key2 and val2 to db2
			testKey2 := common.MakeRandomBytes(10)
			testVal2 := common.MakeRandomBytes(20)

			testKeys2 = append(testKeys2, testKey2)
			testVals2 = append(testVals2, testVal2)

			assert.NoError(t, batch2.Put(testKey2, testVal2))
		}
	}
	assert.NoError(t, batch.Write())
	assert.NoError(t, batch2.Write())

	// check if exist
	for i := 0; i < itemNum; i++ {
		// dynamodb 1 - check if wrote key and val
		returnedVal, returnedErr := dynamo.Get(testKeys[i])
		assert.NoError(t, returnedErr)
		assert.Equal(t, hexutil.Encode(testVals[i]), hexutil.Encode(returnedVal))
		// dynamodb 1 - check if not wrote key2 and val2
		returnedVal, returnedErr = dynamo.Get(testKeys2[i])
		assert.Nil(t, returnedVal, "the entry should not be put in this table")

		// dynamodb 2 - check if wrote key2 and val2
		returnedVal, returnedErr = dynamo2.Get(testKeys2[i])
		assert.NoError(t, returnedErr)
		assert.Equal(t, hexutil.Encode(testVals2[i]), hexutil.Encode(returnedVal))
		// dynamodb 2 - check if not wrote key and val
		returnedVal, returnedErr = dynamo2.Get(testKeys[i])
		assert.Nil(t, returnedVal, "the entry should not be put in this table")
	}
}

func (dynamo *dynamoDB) deleteDB() {
	dynamo.Close()
	dynamo.deleteTable()
	dynamo.fdb.deleteBucket()
}
