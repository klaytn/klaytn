package database

import (
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

func testDynamoDBReadOnly_Put(t *testing.T) {
	dynamo, err := newDynamoDBReadOnly(GetDefaultDynamoDBConfig())
	defer dynamo.deleteDB()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dynamoDB", dynamo.config.TableName)

	testKey := common.MakeRandomBytes(32)
	testVal := common.MakeRandomBytes(500)

	// check if not found for not put key, value
	val, err := dynamo.Get(testKey)
	assert.Nil(t, val)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), dataNotFoundErr.Error())

	// put key, value
	assert.NoError(t, dynamo.Put(testKey, testVal))

	// check if not found for put key, value
	val, err = dynamo.Get(testKey)
	assert.Nil(t, val)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), dataNotFoundErr.Error())
}

func testDynamoDBReadOnly_Write(t *testing.T) {
	dynamo, err := newDynamoDBReadOnly(GetDefaultDynamoDBConfig())
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
		val, err := dynamo.Get(testKeys[i])
		assert.Nil(t, val)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), dataNotFoundErr.Error())
	}
}
