// Copyright 2018 The klaytn Authors
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

/*
Package database implements various types of databases used in Klaytn.
This package is used to read/write data from/to the persistent layer.

Overview of database package

DBManager is the interface used by the consumers of database package.
databaseManager is the implementation of DBManager interface. It contains cacheManager and a list of Database interfaces.
cacheManager caches data stored in the persistent layer, to decrease the direct access to the persistent layer.
Database is the interface for persistent layer implementation. Currently there are 4 implementations, levelDB, memDB,
badgerDB, dynamoDB and shardedDB.

Source Files

  - badger_database.go       : implementation of badgerDB, which wraps github.com/dgraph-io/badger
  - cache_manager.go         : implementation of cacheManager, which manages cache layer over persistent layer
  - db_manager.go            : contains DBManager and databaseManager
  - dynamodb.go              : implementation of dynamoDB, which wraps github.com/aws/aws-sdk-go/service/dynamodb
  - interface.go             : interfaces used outside database package
  - leveldb_database.go      : implementation of levelDB, which wraps github.com/syndtr/goleveldb
  - memory_database.go       : implementation of MemDB, which wraps go native map structure
  - metrics.go               : metrics used in database package, mostly related to cacheManager
  - sharded_database.go      : implementation of shardedDB, which wraps a list of Database interface
  - schema.go                : prefixes and suffixes for database keys and database key generating functions
*/
package database
