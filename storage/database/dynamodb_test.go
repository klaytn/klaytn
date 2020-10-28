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
// For local test, please run the below.
//    $ docker run -d -p 4566:4566 localstack/localstack:0.11.5

package database

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/log/term"
	"github.com/mattn/go-colorable"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
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

type SuiteDynamoDB struct {
	suite.Suite
	dynamoDBs []*dynamoDB
}

func (s *SuiteDynamoDB) TearDownSuite() {
	for _, dynamo := range s.dynamoDBs {
		dynamo.deleteDB()
	}
}

func TestDynamoDB(t *testing.T) {
	suite.Run(t, new(SuiteDynamoDB))
}

func (s *SuiteDynamoDB) TestDynamoDB_Put() {
	dynamo, err := newDynamoDB(GetTestDynamoConfig())
	if err != nil {
		s.FailNow("failed to create dynamoDB", err)
	}
	s.dynamoDBs = append(s.dynamoDBs, dynamo)

	testKey := common.MakeRandomBytes(32)
	testVal := common.MakeRandomBytes(500)

	val, err := dynamo.Get(testKey)

	s.Nil(val)
	s.Error(err)
	s.Equal(err.Error(), dataNotFoundErr.Error())

	s.NoError(dynamo.Put(testKey, testVal))
	returnedVal, returnedErr := dynamo.Get(testKey)
	s.Equal(testVal, returnedVal)
	s.NoError(returnedErr)
}

func (s *SuiteDynamoDB) TestDynamoBatch_Write() {
	dynamo, err := newDynamoDB(GetTestDynamoConfig())
	if err != nil {
		s.FailNow("failed to create dynamoDB", err)
	}
	s.dynamoDBs = append(s.dynamoDBs, dynamo)

	var testKeys [][]byte
	var testVals [][]byte
	batch := dynamo.NewBatch()

	itemNum := 25
	for i := 0; i < itemNum; i++ {
		testKey := common.MakeRandomBytes(32)
		testVal := common.MakeRandomBytes(500)

		testKeys = append(testKeys, testKey)
		testVals = append(testVals, testVal)

		s.NoError(batch.Put(testKey, testVal))
	}
	s.NoError(batch.Write())

	// check if exist
	for i := 0; i < itemNum; i++ {
		returnedVal, returnedErr := dynamo.Get(testKeys[i])
		s.NoError(returnedErr)
		s.Equal(hexutil.Encode(testVals[i]), hexutil.Encode(returnedVal))
	}
}

func (s *SuiteDynamoDB) TestDynamoBatch_WriteLargeData() {
	dynamo, err := newDynamoDB(GetTestDynamoConfig())
	if err != nil {
		s.FailNow("failed to create dynamoDB", err)
	}
	s.dynamoDBs = append(s.dynamoDBs, dynamo)

	var testKeys [][]byte
	var testVals [][]byte
	batch := dynamo.NewBatch()

	itemNum := 26
	for i := 0; i < itemNum; i++ {
		testKey := common.MakeRandomBytes(32)
		testVal := common.MakeRandomBytes(500 * 1024)

		testKeys = append(testKeys, testKey)
		testVals = append(testVals, testVal)

		s.NoError(batch.Put(testKey, testVal))
	}
	s.NoError(batch.Write())

	// check if exist
	for i := 0; i < itemNum; i++ {
		returnedVal, returnedErr := dynamo.Get(testKeys[i])
		s.NoError(returnedErr)
		s.Equal(hexutil.Encode(testVals[i]), hexutil.Encode(returnedVal))
	}
}

func (s *SuiteDynamoDB) TestDynamoBatch_DuplicatedKey() {
	dynamo, err := newDynamoDB(GetTestDynamoConfig())
	if err != nil {
		s.FailNow("failed to create dynamoDB", err)
	}
	s.dynamoDBs = append(s.dynamoDBs, dynamo)

	var testKeys [][]byte
	var testVals [][]byte
	batch := dynamo.NewBatch()

	itemNum := 25
	for i := 0; i < itemNum; i++ {
		testKey := common.MakeRandomBytes(256)
		testVal := common.MakeRandomBytes(600)

		testKeys = append(testKeys, testKey)
		testVals = append(testVals, testVal)

		s.NoError(batch.Put(testKey, testVal))
		s.NoError(batch.Put(testKey, testVal))
		s.NoError(batch.Put(testKey, testVal))
	}
	s.NoError(batch.Write())

	// check if exist
	for i := 0; i < itemNum; i++ {
		returnedVal, returnedErr := dynamo.Get(testKeys[i])
		s.NoError(returnedErr)
		s.Equal(hexutil.Encode(testVals[i]), hexutil.Encode(returnedVal))
	}
}

// testDynamoBatch_WriteMutliTables checks if there is no error when working with more than one tables.
// This also checks if shared workers works as expected.
func (s *SuiteDynamoDB) TestDynamoBatch_WriteMutliTables() {
	// this test might end with Crit, enableLog to find out the log
	//enableLog()

	// create DynamoDB1
	dynamo, err := newDynamoDB(GetTestDynamoConfig())
	if err != nil {
		s.FailNow("failed to create dynamoDB", err)
	}
	s.dynamoDBs = append(s.dynamoDBs, dynamo)

	// create DynamoDB2
	dynamo2, err := newDynamoDB(GetTestDynamoConfig())
	if err != nil {
		s.FailNow("failed to create dynamoDB", err)
	}
	s.dynamoDBs = append(s.dynamoDBs, dynamo2)

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

			s.NoError(batch.Put(testKey, testVal))

			// write key2 and val2 to db2
			testKey2 := common.MakeRandomBytes(10)
			testVal2 := common.MakeRandomBytes(20)

			testKeys2 = append(testKeys2, testKey2)
			testVals2 = append(testVals2, testVal2)

			s.NoError(batch2.Put(testKey2, testVal2))
		}
	}
	s.NoError(batch.Write())
	s.NoError(batch2.Write())

	// check if exist
	for i := 0; i < itemNum; i++ {
		// dynamodb 1 - check if wrote key and val
		returnedVal, returnedErr := dynamo.Get(testKeys[i])
		s.NoError(returnedErr)
		s.Equal(hexutil.Encode(testVals[i]), hexutil.Encode(returnedVal))
		// dynamodb 1 - check if not wrote key2 and val2
		returnedVal, returnedErr = dynamo.Get(testKeys2[i])
		s.Nil(returnedVal, "the entry should not be put in this table")

		// dynamodb 2 - check if wrote key2 and val2
		returnedVal, returnedErr = dynamo2.Get(testKeys2[i])
		s.NoError(returnedErr)
		s.Equal(hexutil.Encode(testVals2[i]), hexutil.Encode(returnedVal))
		// dynamodb 2 - check if not wrote key and val
		returnedVal, returnedErr = dynamo2.Get(testKeys[i])
		s.Nil(returnedVal, "the entry should not be put in this table")
	}
}

func (dynamo *dynamoDB) deleteDB() {
	dynamo.Close()
	dynamo.deleteTable()
	dynamo.fdb.deleteBucket()
}
