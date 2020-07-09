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
// MongoDB needs tuning for performance. Adjust the values below for your server.
//   connection : requests are queued per connection, a thread
//   collection : a collection lock is acquired for every request
//   PoolSize   : connections allowed in the driver's connection pool to each server

package database

import (
	"context"
	"github.com/klaytn/klaytn/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"path/filepath"

	"time"
)

const (
	STOP_TIMEOUT = 10 * time.Second
)

// mongoDBClient controls connections to MongoDB.
// Increasing connections can lead to
// A node creates only one connection.
type mongoDB struct {
	client *mongo.Client   // MongoDB connection
	ctx    context.Context // context for handling mongo connection

	dbName   string            // database name
	collName string            // collection name
	coll     *mongo.Collection // The collection instance of mongoDB

	logger log.Logger // Contextual logger tracking the database path
}

type result struct {
	D []byte
}

// TODO check WiredTiger for memory usage

func getMongoDBOptions() *options.ClientOptions {
	newOption := options.Client()

	newOption.ApplyURI("mongodb://localhost:27017")
	newOption.SetAppName("Klaytn")

	// TODO find an adequate number of pool size through benchmarking
	//      Pool size may be differ for every db type and the number for CPU core
	//newOption.SetMaxPoolSize(100)

	return newOption
}

func NewMongoDB(dbc *DBConfig, entryType DBEntryType) (*mongoDB, error) {
	collName := "0"
	if dbc.Partitioned {
		collName = filepath.Base(dbc.Dir)
	}

	mdb, err := NewMongoDBWithOption(dbBaseDirs[entryType], collName, getMongoDBOptions())
	if err != nil {
		return nil, err
	}

	return mdb, nil
}

func NewMongoDBWithOption(dbName, collName string, mdbOption *options.ClientOptions) (*mongoDB, error) {
	localLogger := logger.NewWith()

	ctx := context.Background()
	client, err := mongo.Connect(ctx, mdbOption)
	if err != nil {
		return nil, err
	}

	localLogger.Info("Allocated MongoDB", "database", dbName, "collection", collName)

	return &mongoDB{
		client: client,
		ctx:    ctx,

		dbName:   dbName,
		collName: collName,
		coll:     client.Database(dbName).Collection(collName),

		logger: localLogger,
	}, nil
}

// TODO find out how index works and apply them for performance
// https://docs.mongodb.com/manual/reference/limits/#indexes
// https://docs.mongodb.com/manual/faq/indexes/
func (db *mongoDB) Put(key []byte, value []byte) error {
	_, err := db.coll.InsertOne(db.ctx, bson.M{"_id": key, "d": value})
	return err
}

func (db *mongoDB) Get(key []byte) ([]byte, error) {
	var r result
	err := db.coll.FindOne(db.ctx, bson.M{"_id": key}).Decode(&r)
	return r.D, err
}

func (db *mongoDB) Has(key []byte) (bool, error) {
	num, err := db.coll.CountDocuments(db.ctx, bson.M{"_id": key})
	return num > 0, err
}

func (db *mongoDB) Delete(key []byte) error {
	var r result
	return db.coll.FindOne(db.ctx, bson.M{"_id": key}).Decode(&r)
}

func (db *mongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), STOP_TIMEOUT)
	defer cancel()
	if err := db.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (db *mongoDB) NewBatch() Batch {
	return &mdbBatch{ctx: db.ctx, coll: db.coll, b: []interface{}{}, logger: db.logger}
}

func (db *mongoDB) Type() DBType {
	return MongoDB
}

// TODO check if PATH can be fetched from mongo
func (db *mongoDB) Path() string {
	return ""
}

func (db *mongoDB) Meter(prefix string) {
	// TODO check mongodb stats
	//  https://github.com/mongodb/mongo-go-driver/tree/master/benchmark
	//  https://docs.mongodb.com/manual/reference/method/db.collection.stats/
}

func (db *mongoDB) NewIterator() Iterator {
	// TODO create iterator for mongoDB
	db.logger.Warn("NewIterator called, but no iterator implemented for mongo")
	return nil
}

func (db *mongoDB) NewIteratorWithStart(start []byte) Iterator {
	db.logger.Warn("NewIterator called, but no iterator implemented for mongo")
	return nil
}

func (db *mongoDB) NewIteratorWithPrefix(prefix []byte) Iterator {
	db.logger.Warn("NewIterator called, but no iterator implemented for mongo")
	return nil
}

type mdbBatch struct {
	ctx  context.Context
	coll *mongo.Collection

	b []interface{}

	logger log.Logger
}

func (b *mdbBatch) Put(key, value []byte) error {
	b.b = append(b.b, bson.M{"_id": key, "d": value})
	return nil
}

func (b *mdbBatch) Write() error {
	result, err := b.coll.InsertMany(b.ctx, b.b)
	if err != nil && len(result.InsertedIDs) != len(b.b) {
		b.logger.Info("failed to batch write mongoDB", "err", err,
			"expected write", len(b.b), "actual write", len(result.InsertedIDs),
			"expected docs", b.b, "write IDs", result.InsertedIDs)
	}
	return errors.Wrap(err, "failed to batch write to mongoDB")
}

func (b *mdbBatch) ValueSize() int {
	return len(b.b)
}

func (b *mdbBatch) Reset() {
	b.b = []interface{}{}
}
