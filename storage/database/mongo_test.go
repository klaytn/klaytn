package database

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/log/term"
	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io"
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

func enableLog() {
	usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	glogger := log.NewGlogHandler(log.StreamHandler(output, log.TerminalFormat(usecolor)))
	log.PrintOrigins(true)
	log.ChangeGlobalLogLevel(glogger, log.Lvl(5))
	glogger.Vmodule("")
	glogger.BacktraceAt("")
	log.Root().SetHandler(glogger)
}

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
	}
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

func TestMongoDB_NewBatch(t *testing.T) {
	enableLog()

	mongoDB, close := newMongoDB()
	defer close()

	type entry struct {
		key   []byte
		value []byte
	}

	e := make([]entry, 10)
	b := mongoDB.NewBatch()
	for i := 0; i < 10; i++ {
		key, value := append([]byte("security-key-"), common.MakeRandomBytes(256)...), common.MakeRandomBytes(600)
		assert.NoError(t, b.Put(key, value))
		e[i] = entry{key, value}
	}
	assert.NoError(t, b.Write())

	for _, kv := range e {
		result, err := mongoDB.Get(kv.key)
		assert.NoError(t, err)
		assert.Equal(t, kv.value, result)
	}
}
