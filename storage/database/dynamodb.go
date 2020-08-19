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
// Database implementation of AWS DynamoDB.
//
// [WARN] Using this DB may cause pricing in your AWS account.
// [WARN] DynamoDB creates both Dynamo DB table and S3 bucket.
//
// You need to set AWS credentials to access to dynamoDB.
//    $ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
//    $ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET

package database

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/rcrowley/go-metrics"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/klaytn/klaytn/log"
)

var overSizedDataPrefix = []byte("oversizeditem")

// errors
var dataNotFoundErr = errors.New("data is not found with the given key")
var nilDynamoConfigErr = errors.New("attempt to create DynamoDB with nil configuration")
var noTableNameErr = errors.New("dynamoDB table name not provided")

// batch write size
const dynamoWriteSizeLimit = 399 * 1024 // The maximum write size is 400KB including attribute names and values
const dynamoBatchSize = 25
const dynamoMaxRetry = 5

// batch write
const WorkerNum = 10
const itemChanSize = WorkerNum * 2

var (
	dynamoOnceClient  sync.Once                   // makes sure dynamo client is created once
	dynamoDBClient    *dynamodb.DynamoDB          // handles dynamoDB connections
	dynamoOnceWorker  sync.Once                   // makes sure worker is created once
	dynamoWriteCh     chan *batchWriteWorkerInput // use global write channel for shared worker
	dynamoOpenedDBNum uint
)

type DynamoDBConfig struct {
	TableName          string
	Region             string // AWS region
	Endpoint           string // Where DynamoDB reside
	IsProvisioned      bool   // Billing mode
	ReadCapacityUnits  int64  // read capacity when provisioned
	WriteCapacityUnits int64  // write capacity when provisioned
}

type batchWriteWorkerInput struct {
	tableName string
	items     []*dynamodb.WriteRequest
	wg        *sync.WaitGroup
}

// TODO-Klaytn refactor the structure : there are common configs that are placed separated
type dynamoDB struct {
	config DynamoDBConfig
	fdb    fileDB     // where over size items are stored
	logger log.Logger // Contextual logger tracking the database path

	// metrics
	batchWriteTimeMeter       metrics.Meter
	batchWriteCountMeter      metrics.Meter
	batchWriteSizeMeter       metrics.Meter
	batchWriteSecPerItemMeter metrics.Meter
	batchWriteSecPerByteMeter metrics.Meter
}

type DynamoData struct {
	Key []byte `json:"Key" dynamodbav:"Key"`
	Val []byte `json:"Val" dynamodbav:"Val"`
}

// TODO-Klaytn: remove the test config when flag setting is completed
/*
 * Please Run DynamoDB local with docker
 * $ docker pull amazon/dynamodb-local
 * $ docker run -d -p 8000:8000 amazon/dynamodb-local
 */
func GetDefaultDynamoDBConfig() *DynamoDBConfig {
	return &DynamoDBConfig{
		Region:             "ap-northeast-2",
		Endpoint:           "https://dynamodb.ap-northeast-2.amazonaws.com",
		TableName:          "klaytn-default" + strconv.Itoa(time.Now().Nanosecond()),
		IsProvisioned:      false,
		ReadCapacityUnits:  10000,
		WriteCapacityUnits: 10000,
	}
}

func NewDynamoDB(config *DynamoDBConfig, dbtype DBEntryType) (*dynamoDB, error) {
	if config == nil {
		return nil, nilDynamoConfigErr
	}
	if len(config.TableName) == 0 {
		return nil, noTableNameErr
	}
	if len(config.Endpoint) == 0 {
		config.Endpoint = "https://dynamodb." + config.Region + ".amazonaws.com"
	}

	s3FileDB, err := newS3FileDB(config.Region, "https://s3."+config.Region+".amazonaws.com", config.TableName)
	if err != nil {
		logger.Error("Unable to create/get S3FileDB", "DB", config.TableName)
		return nil, err
	}

	dynamoOnceClient.Do(func() {
		dynamoDBClient = dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Endpoint: aws.String(config.Endpoint),
				Region:   aws.String(config.Region),
			},
		})))
	})

	dynamoDB := &dynamoDB{
		config: *config,
		fdb:    s3FileDB,
	}
	dynamoDB.config.TableName += "-" + dbBaseDirs[dbtype]
	dynamoDB.logger = logger.NewWith("region", config.Region, "tableName", dynamoDB.config.TableName)

	// Check if the table is ready to serve
	for {
		tableStatus, err := dynamoDB.tableStatus()
		if err != nil {
			if !strings.Contains(err.Error(), "ResourceNotFoundException") {
				dynamoDB.logger.Error("unable to get DynamoDB table status", "err", err.Error())
				return nil, err
			}

			dynamoDB.logger.Info("creating a DynamoDB table", "endPoint", config.Endpoint)
			if err := dynamoDB.createTable(); err != nil {
				dynamoDB.logger.Error("unable to create a DynamoDB table", "err", err.Error())
				return nil, err
			}
		}

		switch tableStatus {
		case dynamodb.TableStatusActive:
			dynamoDB.logger.Warn("Successfully created dynamoDB table. You will be CHARGED until the DB is deleted.", "endPoint", config.Endpoint)
			// count successful table creating
			dynamoOpenedDBNum++
			// create workers on the first successful table creation
			dynamoOnceWorker.Do(func() {
				createBatchWriteWorkerPool(config.Endpoint, config.Region)
			})
			return dynamoDB, nil
		case dynamodb.TableStatusDeleting, dynamodb.TableStatusArchiving, dynamodb.TableStatusArchived:
			return nil, errors.New("failed to get DynamoDB table, table status : " + tableStatus)
		default:
			dynamoDB.logger.Info("waiting for the table to be ready", "table status", tableStatus)
			time.Sleep(1 * time.Second)
		}
	}
}

func (dynamo *dynamoDB) createTable() error {
	input := &dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Key"),
				AttributeType: aws.String("B"), // B - the attribute is of type Binary
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       aws.String("HASH"), // HASH - partition key, RANGE - sort key
			},
		},

		TableName: aws.String(dynamo.config.TableName),
	}

	if dynamo.config.IsProvisioned {
		input.BillingMode = aws.String("PROVISIONED")
		input.ProvisionedThroughput = &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(dynamo.config.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(dynamo.config.WriteCapacityUnits),
		}
		dynamo.logger.Warn("Billing mode is provisioned. You will be charged every hour.", "RCU", dynamo.config.ReadCapacityUnits, "WRU", dynamo.config.WriteCapacityUnits)
	}

	_, err := dynamoDBClient.CreateTable(input)
	if err != nil {
		dynamo.logger.Error("Error while creating the DynamoDB table", "err", err, "tableName", dynamo.config.TableName)
		return err
	}
	dynamo.logger.Warn("Requesting create dynamoDB table. You will be charged until the table is deleted.")
	return nil
}

func (dynamo *dynamoDB) deleteTable() error {
	if _, err := dynamoDBClient.DeleteTable(&dynamodb.DeleteTableInput{TableName: &dynamo.config.TableName}); err != nil {
		dynamo.logger.Error("Error while deleting the DynamoDB table", "tableName", dynamo.config.TableName)
		return err
	}
	dynamo.logger.Info("Successfully deleted the DynamoDB table", "tableName", dynamo.config.TableName)
	return nil
}

func (dynamo *dynamoDB) tableStatus() (string, error) {
	desc, err := dynamo.tableDescription()
	if err != nil {
		return "", err
	}

	return *desc.TableStatus, nil
}

func (dynamo *dynamoDB) tableDescription() (*dynamodb.TableDescription, error) {
	describe, err := dynamoDBClient.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(dynamo.config.TableName)})
	if describe == nil {
		return nil, err
	}

	return describe.Table, err
}

func (dynamo *dynamoDB) Type() DBType {
	return DynamoDB
}

// Put inserts the given key and value pair to the database.
func (dynamo *dynamoDB) Put(key []byte, val []byte) error {
	if len(key) == 0 {
		return nil
	}

	if len(val) > dynamoWriteSizeLimit {
		_, err := dynamo.fdb.write(item{key: key, val: val})
		if err != nil {
			return err
		}
		return dynamo.Put(key, overSizedDataPrefix)
	}

	data := DynamoData{Key: key, Val: val}
	marshaledData, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(dynamo.config.TableName),
		Item:      marshaledData,
	}

	_, err = dynamoDBClient.PutItem(params)
	if err != nil {
		fmt.Printf("Put ERROR: %v\n", err.Error())
		return err
	}

	return nil
}

// Has returns true if the corresponding value to the given key exists.
func (dynamo *dynamoDB) Has(key []byte) (bool, error) {
	if _, err := dynamo.Get(key); err != nil {
		return false, err
	}
	return true, nil
}

// Get returns the corresponding value to the given key if exists.
func (dynamo *dynamoDB) Get(key []byte) ([]byte, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(dynamo.config.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				B: key,
			},
		},
		ConsistentRead: aws.Bool(true),
	}

	result, err := dynamoDBClient.GetItem(params)
	if err != nil {
		fmt.Printf("Get ERROR: %v\n", err.Error())
		return nil, err
	}

	if result.Item == nil {
		return nil, dataNotFoundErr
	}

	var data DynamoData
	if err := dynamodbattribute.UnmarshalMap(result.Item, &data); err != nil {
		return nil, err
	}

	if data.Val == nil {
		return []byte{}, nil
	}

	if bytes.Equal(data.Val, overSizedDataPrefix) {
		return dynamo.fdb.read(key)
	}

	return data.Val, nil
}

// Delete deletes the key from the queue and database
func (dynamo *dynamoDB) Delete(key []byte) error {
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(dynamo.config.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				B: key,
			},
		},
	}

	_, err := dynamoDBClient.DeleteItem(params)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		return err
	}
	return nil
}

func (dynamo *dynamoDB) Close() {
	if dynamoOpenedDBNum > 0 {
		dynamoOpenedDBNum--
	}
	if dynamoOpenedDBNum == 0 {
		close(dynamoWriteCh)
	}
}

func (dynamo *dynamoDB) Meter(prefix string) {
	// TODO-Klaytn: implement this later. Consider the return values of bathItemWrite
	dynamo.batchWriteTimeMeter = metrics.NewRegisteredMeter(prefix+"batchwrite/time", nil)
	dynamo.batchWriteCountMeter = metrics.NewRegisteredMeter(prefix+"batchwrite/count", nil)
	dynamo.batchWriteSizeMeter = metrics.NewRegisteredMeter(prefix+"batchwrite/size", nil)
	dynamo.batchWriteSecPerItemMeter = metrics.NewRegisteredMeter(prefix+"batchwrite/secperitem", nil)
	dynamo.batchWriteSecPerByteMeter = metrics.NewRegisteredMeter(prefix+"batchwrite/secperbyte", nil)
}

func (dynamo *dynamoDB) NewIterator() Iterator {
	// TODO-Klaytn: implement this later.
	return nil
}

func (dynamo *dynamoDB) NewIteratorWithStart(start []byte) Iterator {
	// TODO-Klaytn: implement this later.
	return nil
}

func (dynamo *dynamoDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	// TODO-Klaytn: implement this later.
	return nil
}

func createBatchWriteWorkerPool(endpoint, region string) {
	dynamoWriteCh = make(chan *batchWriteWorkerInput, itemChanSize)
	for i := 0; i < WorkerNum; i++ {
		go createBatchWriteWorker(dynamoWriteCh)
	}
	logger.Info("made dynamo batch write workers", "workerNum", WorkerNum)
}

func createBatchWriteWorker(writeCh <-chan *batchWriteWorkerInput) {
	failCount := 0
	batchWriteInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{},
	}
	logger.Debug("generate a dynamoDB batchWrite worker")

	for batchInput := range writeCh {
		batchWriteInput.RequestItems[batchInput.tableName] = batchInput.items

		BatchWriteItemOutput, err := dynamoDBClient.BatchWriteItem(batchWriteInput)
		numUnprocessed := len(BatchWriteItemOutput.UnprocessedItems[batchInput.tableName])
		for err != nil || numUnprocessed != 0 {
			if err != nil {
				failCount++
				logger.Warn("dynamoDB failed to write batch items",
					"tableName", batchInput.tableName, "err", err, "failCnt", failCount)
				if failCount > dynamoMaxRetry {
					logger.Error("dynamoDB failed many times. sleep a second and retry",
						"tableName", batchInput.tableName, "failCnt", failCount)
					time.Sleep(time.Second)
				}
			}

			if numUnprocessed != 0 {
				logger.Debug("dynamoDB batchWrite remains unprocessedItem",
					"tableName", batchInput.tableName, "numUnprocessedItem", numUnprocessed)
				batchWriteInput.RequestItems[batchInput.tableName] = BatchWriteItemOutput.UnprocessedItems[batchInput.tableName]
			}

			BatchWriteItemOutput, err = dynamoDBClient.BatchWriteItem(batchWriteInput)
			numUnprocessed = len(BatchWriteItemOutput.UnprocessedItems)
		}

		failCount = 0
		batchInput.wg.Done()
	}
	logger.Debug("close a dynamoDB batchWrite worker")
}

func (dynamo *dynamoDB) NewBatch() Batch {
	return &dynamoBatch{db: dynamo, tableName: dynamo.config.TableName, wg: &sync.WaitGroup{}}
}

type dynamoBatch struct {
	db         *dynamoDB
	tableName  string
	batchItems []*dynamodb.WriteRequest
	size       int
	wg         *sync.WaitGroup
}

func (batch *dynamoBatch) Put(key, val []byte) error {
	data := DynamoData{Key: key, Val: val}
	dataSize := len(val)
	if dataSize == 0 {
		return nil
	}

	// If the size of the item is larger than the limit, it should be handled in different way
	if dataSize > dynamoWriteSizeLimit {
		batch.wg.Add(1)
		go func() {
			failCnt := 0
			batch.db.logger.Debug("write large size data into fileDB")

			_, err := batch.db.fdb.write(item{key: key, val: val})
			for err != nil {
				failCnt++
				batch.db.logger.Error("cannot write an item into fileDB. check the status of s3",
					"err", err, "numRetry", failCnt)
				time.Sleep(time.Second)

				batch.db.logger.Warn("retrying write an item into fileDB")
				_, err = batch.db.fdb.write(item{key: key, val: val})
			}
			batch.wg.Done()
		}()
		data.Val = overSizedDataPrefix
		dataSize = len(data.Val)
	}

	marshaledData, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		batch.db.logger.Error("err while batch put", "err", err, "len(val)", len(val))
		return err
	}

	batch.batchItems = append(batch.batchItems, &dynamodb.WriteRequest{
		PutRequest: &dynamodb.PutRequest{Item: marshaledData},
	})
	batch.size += dataSize

	if len(batch.batchItems) == dynamoBatchSize {
		batch.wg.Add(1)
		dynamoWriteCh <- &batchWriteWorkerInput{batch.tableName, batch.batchItems, batch.wg}
		batch.Reset()
	}
	return nil
}

func (batch *dynamoBatch) Write() error {
	var writeRequest []*dynamodb.WriteRequest
	numRemainedItems := len(batch.batchItems)

	for numRemainedItems > 0 {
		if numRemainedItems > dynamoBatchSize {
			writeRequest = batch.batchItems[:dynamoBatchSize]
			batch.batchItems = batch.batchItems[dynamoBatchSize:]
		} else {
			writeRequest = batch.batchItems
		}
		batch.wg.Add(1)
		dynamoWriteCh <- &batchWriteWorkerInput{batch.tableName, writeRequest, batch.wg}
		numRemainedItems -= len(writeRequest)
	}

	batch.wg.Wait()
	return nil
}

func (batch *dynamoBatch) ValueSize() int {
	return batch.size
}

func (batch *dynamoBatch) Reset() {
	batch.batchItems = []*dynamodb.WriteRequest{}
	batch.size = 0
}
