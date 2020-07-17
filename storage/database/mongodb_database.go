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
	"path/filepath"

	"github.com/klaytn/klaytn/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"time"
)

const (
	STOP_TIMEOUT = 10 * time.Second
	URI          = "mongodb://localhost:27017"
	//URI          = "mongodb://winnie:password@mongodb.cluster-cnuopt13avbx.ap-northeast-2.docdb.amazonaws.com:27017/?ssl=true&ssl_ca_certs=rds-combined-ca-bundle.pem&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false"
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

// TODO Set WiredTiger to improve the performance on behalf of memory usage (Server config)
// https://docs.mongodb.com/v3.6/reference/configuration-options/#storage-wiredtiger-options

func getMongoDBOptions() *options.ClientOptions {
	newOption := options.Client()

	newOption.ApplyURI(URI)
	newOption.SetAppName("Klaytn")
	newOption.SetRetryReads(true)
	newOption.SetRetryWrites(true)

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

func (db *mongoDB) Put(key []byte, value []byte) error {
	newKey, newValue := make([]byte, len(key)), make([]byte, len(value))
	copy(newKey, key)
	copy(newValue, value)
	entry := bson.M{
		"_id": newKey,
		"d":   newValue,
	}

	upsert := true
	option := options.ReplaceOptions{Upsert: &upsert}
	_, err := db.coll.ReplaceOne(db.ctx, bson.M{"_id": newKey}, entry, &option)
	// TODO more than one of values in resultof ReplaceOne should be 1

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
	return &mdbBatch{ctx: db.ctx, coll: db.coll, b: []mongo.WriteModel{}, logger: db.logger}
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

	b []mongo.WriteModel

	logger log.Logger
}

func (b *mdbBatch) Put(key, value []byte) error {
	newKey, newValue := make([]byte, len(key)), make([]byte, len(value))
	copy(newKey, key)
	copy(newValue, value)
	entry := bson.M{
		"_id": newKey,
		"d":   newValue,
	}

	upsert := true
	b.b = append(b.b, &mongo.ReplaceOneModel{
		Upsert:      &upsert,
		Filter:      bson.M{"_id": newKey},
		Replacement: entry,
	})
	return nil
}

func (b *mdbBatch) Write() error {
	if len(b.b) == 0 {
		return nil
	}

	r, err := b.coll.BulkWrite(b.ctx, b.b)
	b.logger.Debug("Batch write on mongodb", "InsertedCount", r.InsertedCount,
		"MatchedCount", r.MatchedCount, "ModifiedCount", r.ModifiedCount,
		"DeletedCount", r.DeletedCount, "UpsertedCount", r.UpsertedCount)

	return errors.Wrap(err, "failed to batch write to mongoDB")
}

func (b *mdbBatch) ValueSize() int {
	return len(b.b)
}

func (b *mdbBatch) Reset() {
	b.b = []mongo.WriteModel{}
}
