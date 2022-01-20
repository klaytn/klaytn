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
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/stretchr/testify/suite"
)

// GetTestDynamoConfig gets dynamo config for local test
//
// Please Run DynamoDB local with docker
//    $ docker run -d -p 4566:4566 localstack/localstack:0.11.5
func GetTestDynamoConfig() *DynamoDBConfig {
	return &DynamoDBConfig{
		Region:             "us-east-1",
		Endpoint:           "http://localhost:4566",
		S3Endpoint:         "http://localhost:4566",
		TableName:          "klaytn-default" + strconv.Itoa(time.Now().Nanosecond()),
		IsProvisioned:      false,
		ReadCapacityUnits:  10000,
		WriteCapacityUnits: 10000,
		ReadOnly:           false,
		PerfCheck:          false,
	}
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

func testDynamoDB(t *testing.T) {
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

// TestDynamoDB_Timeout tests if a timeout error occurs.
// When there is no answer from DynamoDB server due to network failure,
// a timeout error should occur.
// A fake server is setup to simulate a server with a response latency.
func (s *SuiteDynamoDB) TestDynamoDB_Timeout() {
	// fakeEndpoint allows TCP handshakes, but doesn't answer anything to client.
	// The fake server is used to produce a network failure scenario.
	fakeEndpoint := "127.0.0.1:14566"
	maxRetries := 2
	timeout := 1 * time.Second
	serverReadyWg := sync.WaitGroup{}
	serverReadyWg.Add(1)

	go func() {
		tcpAddr, err := net.ResolveTCPAddr("tcp", fakeEndpoint)
		if err != nil {
			s.FailNow(err.Error())
		}

		listen, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			s.FailNow(err.Error())
		}
		defer listen.Close()
		serverReadyWg.Done()

		for i := 0; i < maxRetries+1; i++ {
			// Deadline prevents infinite waiting of the fake server
			// Wait longer than (maxRetries+1) * timeout
			if err := listen.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
				s.FailNow(err.Error())
			}
			if _, err := listen.AcceptTCP(); err != nil {
				// the fake server ends silently when it meets deadline
				if strings.Contains(err.Error(), "timeout") {
					return
				}
				s.FailNow(err.Error())
			}
		}
	}()

	conf := GetTestDynamoConfig() // dummy values to create dynamoDB
	conf.Endpoint = fakeEndpoint
	conf.TableName = "test-table"

	dynamoDBClient = dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint:   aws.String(conf.Endpoint),
			Region:     aws.String(conf.Region),
			MaxRetries: aws.Int(maxRetries),
			HTTPClient: &http.Client{Timeout: timeout},
		},
	})))
	expectedElapsed := time.Duration(maxRetries+1) * timeout

	dynamoDB := &dynamoDB{
		config: *conf,
		logger: logger.NewWith("logger"),
	}
	serverReadyWg.Wait()

	start := time.Now()
	_, err := dynamoDB.get([]byte{0x1, 0x2})
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), "Timeout"))
	s.Equal(expectedElapsed, time.Since(start).Round(time.Second))
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
