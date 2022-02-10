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
	"github.com/klaytn/klaytn/storage"
	"github.com/stretchr/testify/assert"
)

func TestDynamoDBReadOnly_Put(t *testing.T) {
	storage.SkipLocalTest(t)

	dynamo, err := newDynamoDBReadOnly(GetTestDynamoConfig())
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

func TestDynamoDBReadOnly_Write(t *testing.T) {
	storage.SkipLocalTest(t)

	dynamo, err := newDynamoDBReadOnly(GetTestDynamoConfig())
	defer dynamo.deleteDB()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dynamoDB", dynamo.config.TableName)

	var testKeys [][]byte
	var testVals [][]byte
	batch := dynamo.NewBatch()

	itemNum := 25
	// put items
	for i := 0; i < itemNum; i++ {
		testKey := common.MakeRandomBytes(32)
		testVal := common.MakeRandomBytes(500)

		testKeys = append(testKeys, testKey)
		testVals = append(testVals, testVal)

		assert.NoError(t, batch.Put(testKey, testVal))
	}
	assert.NoError(t, batch.Write())

	// check if not exist
	for i := 0; i < itemNum; i++ {
		val, err := dynamo.Get(testKeys[i])
		assert.Nil(t, val)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), dataNotFoundErr.Error())
	}
}

func (dynamo *dynamoDBReadOnly) deleteDB() {
	dynamo.Close()
	dynamo.deleteTable()
	dynamo.fdb.deleteBucket()
}
