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
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	klaytnmetrics "github.com/klaytn/klaytn/metrics"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/log"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
)

var overSizedDataPrefix = []byte("oversizeditem")

// Performance of batch operations of DynamoDB are collected by default.
var dynamoBatchWriteTimeMeter metrics.Meter = &metrics.NilMeter{}

// errors
var dataNotFoundErr = errors.New("data is not found with the given key")

var (
	nilDynamoConfigErr = errors.New("attempt to create DynamoDB with nil configuration")
	noTableNameErr     = errors.New("dynamoDB table name not provided")
)

// batch write size
const dynamoWriteSizeLimit = 399 * 1024 // The maximum write size is 400KB including attribute names and values
const (
	dynamoBatchSize = 25
	dynamoMaxRetry  = 20
	dynamoTimeout   = 10 * time.Second
)

// batch write
const WorkerNum = 10
const itemChanSize = WorkerNum * 2

var (
	dynamoDBClient    *dynamodb.DynamoDB          // handles dynamoDB connections
	dynamoWriteCh     chan *batchWriteWorkerInput // use global write channel for shared worker
	dynamoOnceWorker  = &sync.Once{}              // makes sure worker is created once
	dynamoOpenedDBNum uint
)

type DynamoDBConfig struct {
	TableName          string
	Region             string // AWS region
	Endpoint           string // Where DynamoDB reside (Used to specify the localstack endpoint on the test)
	S3Endpoint         string // Where S3 reside
	IsProvisioned      bool   // Billing mode
	ReadCapacityUnits  int64  // read capacity when provisioned
	WriteCapacityUnits int64  // write capacity when provisioned
	ReadOnly           bool   // disables write
	PerfCheck          bool
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
	getTimer klaytnmetrics.HybridTimer
	putTimer klaytnmetrics.HybridTimer
}

type DynamoData struct {
	Key []byte `json:"Key" dynamodbav:"Key"`
	Val []byte `json:"Val" dynamodbav:"Val"`
}

// CustomRetryer wraps AWS SDK's built in DefaultRetryer adding additional custom features.
// DefaultRetryer of AWS SDK has its own standard of retryable situation,
// but it's not proper when network environment is not stable.
// CustomRetryer conservatively retry in all error cases because DB failure of Klaytn is critical.
type CustomRetryer struct {
	client.DefaultRetryer
}

// ShouldRetry overrides AWS SDK's built in DefaultRetryer to retry in all error cases.
func (r CustomRetryer) ShouldRetry(req *request.Request) bool {
	logger.Debug("dynamoDB client retry", "error", req.Error, "retryCnt", req.RetryCount, "retryDelay",
		req.RetryDelay, "maxRetry", r.MaxRetries())
	return req.Error != nil && req.RetryCount < r.MaxRetries()
}

// GetTestDynamoConfig gets dynamo config for actual aws DynamoDB test
//
// If you use this config, you will be charged for what you use.
// You need to set AWS credentials to access to dynamoDB.
//    $ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
//    $ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET
func GetDefaultDynamoDBConfig() *DynamoDBConfig {
	return &DynamoDBConfig{
		Region:             "ap-northeast-2",
		Endpoint:           "", // nil or "" means the default generated endpoint
		TableName:          "klaytn-default" + strconv.Itoa(time.Now().Nanosecond()),
		IsProvisioned:      false,
		ReadCapacityUnits:  10000,
		WriteCapacityUnits: 10000,
		ReadOnly:           false,
		PerfCheck:          true,
	}
}

// NewDynamoDB creates either dynamoDB or dynamoDBReadOnly depending on config.ReadOnly.
func NewDynamoDB(config *DynamoDBConfig) (Database, error) {
	if config.ReadOnly {
		return newDynamoDBReadOnly(config)
	}
	return newDynamoDB(config)
}

// newDynamoDB creates dynamoDB. dynamoDB can be used to create dynamoDBReadOnly.
func newDynamoDB(config *DynamoDBConfig) (*dynamoDB, error) {
	if config == nil {
		return nil, nilDynamoConfigErr
	}
	if len(config.TableName) == 0 {
		return nil, noTableNameErr
	}

	config.TableName = strings.ReplaceAll(config.TableName, "_", "-")

	s3FileDB, err := newS3FileDB(config.Region, config.S3Endpoint, config.TableName)
	if err != nil {
		logger.Error("Unable to create/get S3FileDB", "DB", config.TableName)
		return nil, err
	}

	if dynamoDBClient == nil {
		dynamoDBClient = dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Retryer: CustomRetryer{
					DefaultRetryer: client.DefaultRetryer{
						NumMaxRetries:    dynamoMaxRetry,
						MaxRetryDelay:    time.Second,
						MaxThrottleDelay: time.Second,
					},
				},
				Endpoint:         aws.String(config.Endpoint),
				Region:           aws.String(config.Region),
				S3ForcePathStyle: aws.Bool(true),
				MaxRetries:       aws.Int(dynamoMaxRetry),
				HTTPClient:       &http.Client{Timeout: dynamoTimeout}, // default client is &http.Client{}
			},
		})))
	}
	dynamoDB := &dynamoDB{
		config: *config,
		fdb:    s3FileDB,
	}

	dynamoDB.logger = logger.NewWith("region", config.Region, "tableName", dynamoDB.config.TableName)

	// Check if the table is ready to serve
	for {
		tableStatus, err := dynamoDB.tableStatus()
		if err != nil {
			if !strings.Contains(err.Error(), "ResourceNotFoundException") {
				dynamoDB.logger.Error("unable to get DynamoDB table status", "err", err.Error())
				return nil, err
			}

			dynamoDB.logger.Warn("creating a DynamoDB table. You will be CHARGED until the DB is deleted")
			if err := dynamoDB.createTable(); err != nil {
				dynamoDB.logger.Error("unable to create a DynamoDB table", "err", err.Error())
				return nil, err
			}
		}

		switch tableStatus {
		case dynamodb.TableStatusActive:
			if !dynamoDB.config.ReadOnly {
				// count successful table creating
				dynamoOpenedDBNum++
				// create workers on the first successful table creation
				dynamoOnceWorker.Do(func() {
					createBatchWriteWorkerPool()
				})
			}
			dynamoDB.logger.Info("successfully created dynamoDB session")
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
	if dynamo.config.PerfCheck {
		start := time.Now()
		err := dynamo.put(key, val)
		dynamo.putTimer.Update(time.Since(start))
		return err
	}
	return dynamo.put(key, val)
}

func (dynamo *dynamoDB) put(key []byte, val []byte) error {
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
		dynamo.logger.Crit("failed to put an item", "err", err, "key", hexutil.Encode(data.Key))
		return err
	}

	return nil
}

// Has returns true if the corresponding value to the given key exists.
func (dynamo *dynamoDB) Has(key []byte) (bool, error) {
	if _, err := dynamo.Get(key); err != nil {
		if err == dataNotFoundErr {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Get returns the corresponding value to the given key if exists.
func (dynamo *dynamoDB) Get(key []byte) ([]byte, error) {
	if dynamo.config.PerfCheck {
		start := time.Now()
		val, err := dynamo.get(key)
		dynamo.getTimer.Update(time.Since(start))
		return val, err
	}
	return dynamo.get(key)
}

func (dynamo *dynamoDB) get(key []byte) ([]byte, error) {
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
		dynamo.logger.Crit("failed to get an item", "err", err, "key", hexutil.Encode(key))
		return nil, err
	}

	if result.Item == nil {
		return nil, dataNotFoundErr
	}

	var data DynamoData
	if err := dynamodbattribute.UnmarshalMap(result.Item, &data); err != nil {
		dynamo.logger.Crit("failed to unmarshal dynamodb data", "err", err)
		return nil, err
	}

	if data.Val == nil {
		return []byte{}, nil
	}

	if bytes.Equal(data.Val, overSizedDataPrefix) {
		ret, err := dynamo.fdb.read(key)
		if err != nil {
			dynamo.logger.Crit("failed to read filedb data", "err", err, "key", hexutil.Encode(key))
		}
		return ret, err
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
		dynamo.logger.Crit("failed to delete an item", "err", err, "key", hexutil.Encode(key))
		return err
	}
	return nil
}

func (dynamo *dynamoDB) Close() {
	if dynamoOpenedDBNum > 0 {
		dynamoOpenedDBNum--
	}
	if dynamoOpenedDBNum == 0 && dynamoWriteCh != nil {
		close(dynamoWriteCh)
	}
}

func (dynamo *dynamoDB) Meter(prefix string) {
	dynamo.getTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"get/time", nil)
	dynamo.putTimer = klaytnmetrics.NewRegisteredHybridTimer(prefix+"put/time", nil)
	dynamoBatchWriteTimeMeter = metrics.NewRegisteredMeter(prefix+"batchwrite/time", nil)
}

func (dynamo *dynamoDB) NewIterator(prefix []byte, start []byte) Iterator {
	// TODO-Klaytn: implement this later.
	return nil
}

func createBatchWriteWorkerPool() {
	dynamoWriteCh = make(chan *batchWriteWorkerInput, itemChanSize)
	for i := 0; i < WorkerNum; i++ {
		go createBatchWriteWorker(dynamoWriteCh)
	}
	logger.Info("made dynamo batch write workers", "workerNum", WorkerNum)
}

func createBatchWriteWorker(writeCh <-chan *batchWriteWorkerInput) {
	failCount := 0
	logger.Debug("generate a dynamoDB batchWrite worker")

	for batchInput := range writeCh {
		batchWriteInput := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{},
		}
		batchWriteInput.RequestItems[batchInput.tableName] = batchInput.items

		BatchWriteItemOutput, err := dynamoDBClient.BatchWriteItem(batchWriteInput)
		numUnprocessed := len(BatchWriteItemOutput.UnprocessedItems[batchInput.tableName])
		for err != nil || numUnprocessed != 0 {
			if err != nil {
				// ValidationException occurs when a required parameter is missing, a value is out of range,
				// or data types mismatch and so on. If this is the case, check if there is a duplicated key,
				// batch length out of range, null value and so on.
				// When ValidationException occurs, retrying won't fix the problem.
				if strings.Contains(err.Error(), "ValidationException") {
					logger.Crit("Invalid input for dynamoDB BatchWrite",
						"err", err, "tableName", batchInput.tableName, "itemNum", len(batchInput.items))
				}
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

			start := time.Now()
			BatchWriteItemOutput, err = dynamoDBClient.BatchWriteItem(batchWriteInput)
			dynamoBatchWriteTimeMeter.Mark(int64(time.Since(start)))
			numUnprocessed = len(BatchWriteItemOutput.UnprocessedItems)
		}

		failCount = 0
		batchInput.wg.Done()
	}
	logger.Debug("close a dynamoDB batchWrite worker")
}

func (dynamo *dynamoDB) NewBatch() Batch {
	return &dynamoBatch{db: dynamo, tableName: dynamo.config.TableName, wg: &sync.WaitGroup{}, keyMap: map[string]struct{}{}}
}

type dynamoBatch struct {
	db         *dynamoDB
	tableName  string
	batchItems []*dynamodb.WriteRequest
	keyMap     map[string]struct{} // checks duplication of keys
	size       int
	wg         *sync.WaitGroup
}

// Put adds an item to dynamo batch.
// If the number of items in batch reaches dynamoBatchSize, a write request to dynamoDB is made.
// Each batch write is executed in thread. (There is an worker pool for dynamo batch write)
//
// Note: If there is a duplicated key in a batch, only the first value is written.
func (batch *dynamoBatch) Put(key, val []byte) error {
	// if there is an duplicated key in batch, skip
	if _, exist := batch.keyMap[string(key)]; exist {
		return nil
	}
	batch.keyMap[string(key)] = struct{}{}

	data := DynamoData{Key: key, Val: val}
	dataSize := len(val)

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

// Delete inserts the a key removal into the batch for later committing.
func (batch *dynamoBatch) Delete(key []byte) error {
	logger.CritWithStack("Delete should not be called when using dynamodb batch")
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
	batch.keyMap = map[string]struct{}{}
	batch.size = 0
}

func (batch *dynamoBatch) Replay(w KeyValueWriter) error {
	logger.CritWithStack("Replay should not be called when using dynamodb batch")
	return nil
}
