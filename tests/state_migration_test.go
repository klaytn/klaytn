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

package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

// continuous occurrence of state trie migration and node restart must success
func TestMigration_ContinuousRestartAndMigration(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	fullNode, node, validator, chainID, workspace, richAccount, _, _ := newSimpleBlockchain(t, 10)
	defer os.RemoveAll(workspace)

	stateTriePath := []byte("statetrie")

	numTxs := []int{10, 100}
	for i := 0; i < len(numTxs); i++ {
		numTx := numTxs[i%len(numTxs)]
		t.Log("attempt", strconv.Itoa(i), " : deployRandomTxs of", strconv.Itoa(numTx))
		deployRandomTxs(t, node.TxPool(), chainID, richAccount, numTx)
		time.Sleep(3 * time.Second) // wait until txpool is flushed

		startMigration(t, node)
		time.Sleep(1 * time.Second)

		t.Log("migration state before restart", node.ChainDB().InMigration())
		fullNode, node = restartNode(t, fullNode, node, workspace, validator)

		waitMigrationEnds(t, node)

		// check if migration succeeds (StateTrieDB changes when migration finishes)
		newPathKey := append([]byte("databaseDirectory"), common.Int64ToByteBigEndian(uint64(database.StateTrieDB))...)
		newStateTriePath, err := node.ChainDB().GetMiscDB().Get(newPathKey)
		assert.NoError(t, err)
		assert.NotEqual(t, stateTriePath, newStateTriePath, "migration failed")
		stateTriePath = newStateTriePath
	}

	stopNode(t, fullNode)
}

// state trie DB should be determined by the values of miscDB
func TestMigration_StartMigrationByMiscDB(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	fullNode, cn, validator, _, workspace, _, _, _ := newSimpleBlockchain(t, 10)
	defer os.RemoveAll(workspace)

	stateTriePathKey := append([]byte("databaseDirectory"), common.Int64ToByteBigEndian(uint64(database.StateTrieDB))...)

	// use the default StateTrie DB if it is not set on miscDB
	{
		// check if stateDB value is not stored in miscDB
		key, err := cn.ChainDB().GetMiscDB().Get(stateTriePathKey)
		assert.Error(t, err)
		assert.Len(t, key, 0)

		// write values in stateDB and check if the values are stored in DB
		entries := writeRandomValueToStateTrieDB(t, cn.ChainDB().NewBatch(database.StateTrieDB))
		stopNode(t, fullNode) // stop node to release DB lock
		checkIfStoredInDB(t, cn.ChainDB().GetDBConfig().NumStateTrieShards, filepath.Join(cn.ChainDB().GetDBConfig().Dir, "statetrie"), entries)
		fullNode, cn = startNode(t, workspace, validator)
	}

	// use the state trie DB that is specified in miscDB
	{
		// change stateDB value in miscDB
		newDBPath := "NEW_STATE_TRIE_DB_PATH"
		err := cn.ChainDB().GetMiscDB().Put(stateTriePathKey, []byte(newDBPath))
		assert.NoError(t, err)

		// an error expected on node start
		stopNode(t, fullNode)
		_, _, err = newKlaytnNode(t, workspace, validator)
		assert.Error(t, err, "start failure expected, changed state trie db has no data") // error expected
	}
}

func writeRandomValueToStateTrieDB(t *testing.T, batch database.Batch) map[string]string {
	entries := make(map[string]string, 10)

	for i := 0; i < 10; i++ {
		key, value := common.MakeRandomBytes(common.HashLength), common.MakeRandomBytes(400)
		err := batch.Put(key, value)
		assert.NoError(t, err)
		entries[string(key)] = string(value)
	}
	assert.NoError(t, batch.Write())

	return entries
}

func checkIfStoredInDB(t *testing.T, numShard uint, dir string, entries map[string]string) {
	dbs := make([]*leveldb.DB, numShard)
	for i := 0; i < 4; i++ {
		var err error
		dbs[i], err = leveldb.OpenFile(dir+"/"+strconv.Itoa(i), nil)
		assert.NoError(t, err)
		defer dbs[i].Close()
	}
	for k, v := range entries {
		datas := make([][]byte, 4)
		for i := 0; i < 4; i++ {
			datas[i], _ = dbs[i].Get([]byte(k), nil)
		}
		assert.Contains(t, datas, []byte(v), "value written in stateDB does not actually exist in DB")
	}
}

// if migration status is set on miscDB and a node is restarted, migration should start
func TestMigration_StartMigrationByMiscDBOnRestart(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	fullNode, node, validator, chainID, workspace, richAccount, _, _ := newSimpleBlockchain(t, 10)
	defer os.RemoveAll(workspace)
	miscDB := node.ChainDB().GetMiscDB()

	// size up state trie to be prepared for migration
	deployRandomTxs(t, node.TxPool(), chainID, richAccount, 100)
	time.Sleep(time.Second)

	// set migration status in miscDB
	migrationBlockNum := node.BlockChain().CurrentBlock().Header().Number.Uint64()
	err := miscDB.Put([]byte("migrationStatus"), common.Int64ToByteBigEndian(migrationBlockNum))
	assert.NoError(t, err)

	// set migration db path in miscDB
	migrationPathKey := append([]byte("databaseDirectory"), common.Int64ToByteBigEndian(uint64(database.StateTrieMigrationDB))...)
	migrationPath := []byte("statetrie_migrated_" + strconv.FormatUint(migrationBlockNum, 10))
	err = miscDB.Put(migrationPathKey, migrationPath)
	assert.NoError(t, err)

	// check migration Status in cache before restart
	assert.False(t, node.ChainDB().InMigration(), "migration has not started yet")
	assert.NotEqual(t, migrationBlockNum, node.ChainDB().MigrationBlockNumber(), "migration has not started yet")

	fullNode, node = restartNode(t, fullNode, node, workspace, validator)
	miscDB = node.ChainDB().GetMiscDB()

	// check migration Status in cache after restart
	if node.ChainDB().InMigration() {
		assert.Equal(t, migrationBlockNum, node.ChainDB().MigrationBlockNumber(), "migration block number should match")
		t.Log("Checked migration status while migration in on process")
	}

	waitMigrationEnds(t, node)

	// state trie path should not be "statetrie" in miscDB
	newPathKey := append([]byte("databaseDirectory"), common.Int64ToByteBigEndian(uint64(database.StateTrieDB))...)
	dir, err := miscDB.Get(newPathKey)
	assert.NoError(t, err)
	assert.NotEqual(t, "statetrie", string(dir), "migration failed")

	stopNode(t, fullNode)
}

func newSimpleBlockchain(t *testing.T, numAccounts int) (*node.Node, *cn.CN, *TestAccountType, *big.Int, string, *TestAccountType, []*TestAccountType, []*TestAccountType) {
	t.Log("=========== create blockchain ==============")
	fullNode, node, validator, chainID, workspace := newBlockchain(t)
	richAccount, accounts, contractAccounts := createAccount(t, numAccounts, validator)
	time.Sleep(5 * time.Second)

	return fullNode, node, validator, chainID, workspace, richAccount, accounts, contractAccounts
}

func startMigration(t *testing.T, node *cn.CN) {
	waitMigrationEnds(t, node)

	t.Log("=========== migrate trie ==============")
	err := node.BlockChain().PrepareStateMigration()
	assert.NoError(t, err)
}

func restartNode(t *testing.T, fullNode *node.Node, node *cn.CN, workspace string, validator *TestAccountType) (*node.Node, *cn.CN) {
	stopNode(t, fullNode)
	time.Sleep(2 * time.Second)
	newFullNode, newNode := startNode(t, workspace, validator)
	time.Sleep(2 * time.Second)

	return newFullNode, newNode
}

func startNode(t *testing.T, workspace string, validator *TestAccountType) (fullNode *node.Node, node *cn.CN) {
	t.Log("=========== starting node ==============")
	newFullNode, newNode, err := newKlaytnNode(t, workspace, validator)
	assert.NoError(t, err)
	if err := newNode.StartMining(false); err != nil {
		t.Fatal()
	}

	return newFullNode, newNode
}

func stopNode(t *testing.T, fullNode *node.Node) {
	if err := fullNode.Stop(); err != nil {
		t.Fatal(err)
	}
	t.Log("=========== stopped node ==============")
}

func waitMigrationEnds(t *testing.T, node *cn.CN) {
	for node.ChainDB().InMigration() {
		t.Log("state trie migration is processing; sleep for a second before a new migration")
		time.Sleep(time.Second)
	}
}
