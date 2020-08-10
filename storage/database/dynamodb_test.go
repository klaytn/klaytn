package database

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDynamoDB_Put(t *testing.T) {
	dynamo, err := NewDynamoDB(createTestDynamoDBConfig())
	defer dynamo.deletedDB()
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

func TestDynamoBatch_Write(t *testing.T) {
	dynamo, err := NewDynamoDB(createTestDynamoDBConfig())
	defer dynamo.deletedDB()
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

func TestDynamoBatch_WriteLargeData(t *testing.T) {
	dynamo, err := NewDynamoDB(createTestDynamoDBConfig())
	defer dynamo.deletedDB()
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

func (dynamo *dynamoDB) deletedDB() {
	dynamo.Close()
	dynamo.deleteTable()
	dynamo.fdb.deleteBucket()
}
