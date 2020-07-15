package database

import (
	"os"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

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
