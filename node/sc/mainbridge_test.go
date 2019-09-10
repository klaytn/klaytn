package sc

import (
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const testNetVersion = uint64(8888)

var testProtocolVersion = int(SCProtocolVersion[0])

// testNewMainBridge returns a test MainBridge.
func testNewMainBridge(t *testing.T) *MainBridge {
	sCtx := node.NewServiceContext(&node.DefaultConfig, map[reflect.Type]node.Service{}, &event.TypeMux{}, &accounts.Manager{})
	mBridge, err := NewMainBridge(sCtx, &SCConfig{NetworkId: testNetVersion})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, mBridge)

	return mBridge
}

// TestCreateDB tests creation of chain database and proper working of database operation.
func TestCreateDB(t *testing.T) {
	sCtx := node.NewServiceContext(&node.DefaultConfig, map[reflect.Type]node.Service{}, &event.TypeMux{}, &accounts.Manager{})
	sConfig := &SCConfig{}
	name := "testDB"

	testKey := []byte{0x33}
	testValue := []byte{44}

	// Create DB
	dbManager := CreateDB(sCtx, sConfig, name)
	defer dbManager.Close()
	assert.NotNil(t, dbManager)

	// Check the type of DB
	db := dbManager.GetDB()
	assert.Equal(t, database.LevelDB, db.Type())

	// Check existence of a key
	hasKey, err := db.Has(testKey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, hasKey)

	// Put the key and check existence of the key again
	db.Put(testKey, testValue)
	hasKey, err = db.Has(testKey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, hasKey)

	// Get the value of the key
	val, err := db.Get(testKey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testValue, val)
}

// TestMainBridge tests some getters and basic operation of MainBridge.
func TestMainBridge_basic(t *testing.T) {
	// Create a test MainBridge
	mBridge := testNewMainBridge(t)

	// APIs returns default rpc APIs of MainBridge
	apis := mBridge.APIs()
	assert.Equal(t, 2, len(apis))
	assert.Equal(t, "mainbridge", apis[0].Namespace)
	assert.Equal(t, "mainbridge", apis[1].Namespace)

	// Getters for elements of MainBridge
	assert.Equal(t, true, mBridge.IsListening()) // Always returns `true`
	assert.Equal(t, testProtocolVersion, mBridge.ProtocolVersion())
	assert.Equal(t, testNetVersion, mBridge.NetVersion())

	// New components of MainBridge which will update old components
	bc := &blockchain.BlockChain{}
	txPool := &blockchain.TxPool{}
	compAPIs := []rpc.API{
		{
			Namespace: "klay",
			Version:   "1.0",
			Service:   api.NewPublicKlayAPI(&cn.ServiceChainAPIBackend{}),
			Public:    true,
		},
	}
	var comp []interface{}
	comp = append(comp, bc)
	comp = append(comp, txPool)
	comp = append(comp, compAPIs)

	// Check initial status of components
	assert.NotEqual(t, bc, mBridge.blockchain)
	assert.NotEqual(t, txPool, mBridge.txPool)
	assert.Nil(t, mBridge.rpcServer.GetServices()["klay"])

	// Update and check MainBridge components
	mBridge.SetComponents(comp)
	assert.Equal(t, bc, mBridge.blockchain)
	assert.Equal(t, txPool, mBridge.txPool)
	assert.NotNil(t, mBridge.rpcServer.GetServices()["klay"])

	// Start and stop MainBridge
	if err := mBridge.Start(p2p.SingleChannelServer{}); err != nil {
		t.Fatal(err)
	}
	defer mBridge.Stop()
}
