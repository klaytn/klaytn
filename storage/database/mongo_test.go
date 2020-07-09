package database

import (
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"testing"
	"time"
)

//const (
//	START_TIMEOUT = 10 * time.Second
//	STOP_TIMEOUT  = 10 * time.Second
//)

const (
	mongoDBURI = "mongodb://localhost:27017"
)

func setUpConnection(t *testing.T, f func(*mongo.Client)) {
	// make connection
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		t.Fatal(err, "Unable to connection MongoDB")
	}

	time.Sleep(5 * time.Second)

	// test connection
	err = client.Ping(nil, readpref.Primary())
	assert.NoError(t, err)

	// run configured function
	f(client)

	// close connection
	if err := client.Disconnect(nil); err != nil {
		panic(err)
	}
}

func newMongoDB() (*mongoDB, func()) {
	dirName := "test"
	db, err := NewMongoDBWithOption(dirName, "0", getMongoDBOptions())
	if err != nil {
		panic("failed to create test database: " + err.Error())
	}

	return db, func() {
		db.Close()
		os.RemoveAll(dirName)
	}
}

func TestConnection(t *testing.T) {
	for i := 0; i < 10; i++ {
		newMongoDB()
		//setUpConnection(t, func(m *mongo.Client) {})
	}
}

func TestDB(t *testing.T) {

}

func TestReadWrite(t *testing.T) {
	mongoDB, close := newMongoDB()
	defer close()

	for i := 0; i < 100; i++ {
		key, value := common.MakeRandomBytes(256), common.MakeRandomBytes(600)
		assert.NoError(t, mongoDB.Put(key, value))
		result, err := mongoDB.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	}
}
