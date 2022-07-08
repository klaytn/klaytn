// Modifications Copyright 2020 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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

package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBType_ToValid(t *testing.T) {
	validDBType := []DBType{LevelDB, BadgerDB, MemoryDB, DynamoDB}

	// check acceptable dbtype values
	acceptableDBTypes := []DBType{
		"LevelDB", "leveldb", "levelDB", "LevelDB",
		"BadgerDB", "badgerdb", "Badgerdb", "badgerDB",
		"MemoryDB", "memorydb", "Memorydb", "memoryDB",
		"DynamoDBS3", "dynamodbs3", "Dynamodbs3", "DynamodbS3", "dynamodbS3",
	}

	for _, dbtype := range acceptableDBTypes {
		newType := dbtype.ToValid()
		assert.Contains(t, validDBType, newType, "dbtype not acceptable: "+dbtype)
	}

	// check not acceptable dbtype values
	notAcceptableDBTypes := []DBType{
		"level", "ldb", "badger", "memory",
		"dynamo", "dynamodb", "dynamos3", "ddb", "ShardedDB",
	}

	for _, dbtype := range notAcceptableDBTypes {
		newType := dbtype.ToValid()
		assert.Equal(t, DBType(""), newType, "dbtype should not acceptable:"+dbtype)
	}
}
