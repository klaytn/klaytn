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
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func testDynamoDB_Put(t *testing.T) {
	dynamo, err := NewDynamoDB(GetDefaultDynamoDBConfig())
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

func testDynamoBatch_Write(t *testing.T) {
	dynamo, err := NewDynamoDB(GetDefaultDynamoDBConfig())
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

func testDynamoBatch_WriteLargeData(t *testing.T) {
	dynamo, err := NewDynamoDB(GetDefaultDynamoDBConfig())
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
