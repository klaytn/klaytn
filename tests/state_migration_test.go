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
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/node"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"
)

// continues occurrence of state trie migration and node restart must success
func TestMigration_ContinuesRestartAndMigration(t *testing.T) {
	fullNode, node, validator, chainId, workspace, richAccount, _, _ := newSimpleBlockchain(t, 100)
	defer os.RemoveAll(workspace)

	numTxs := []int{100, 1000, 10000}
	for i := 0; i < len(numTxs)*5; i++ {
		numTx := numTxs[i%len(numTxs)]
		t.Log("attempt", strconv.Itoa(i), "deployRandomTxs of", strconv.Itoa(numTx))
		deployRandomTxs(t, node.TxPool(), chainId, richAccount, numTx)
		time.Sleep(10 * time.Second) // wait until txpool is flushed

		startMigration(t, node)
		time.Sleep(1 * time.Second)
		t.Log("migration state before restart", node.ChainDB().InMigration())
		fullNode, node = restartNode(t, fullNode, node, workspace, validator)
	}
}

// if migration status is set on miscDB and a node is restarted, migration should start
func TestMigration_StartMigrationByMiscDB(t *testing.T) {
	fullNode, node, validator, chainId, workspace, richAccount, _, _ := newSimpleBlockchain(t, 100)
	defer os.RemoveAll(workspace)
	miscDB := node.ChainDB().GetMiscDB()

	// size up state trie to be prepared for migration
	deployRandomTxs(t, node.TxPool(), chainId, richAccount, 1000)

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
	assert.False(t, node.ChainDB().InMigration())
	assert.NotEqual(t, migrationBlockNum, node.ChainDB().MigrationBlockNumber())

	fullNode, node = restartNode(t, fullNode, node, workspace, validator)
	miscDB = node.ChainDB().GetMiscDB()

	// check migration Status in cache after restart
	if node.ChainDB().InMigration() {
		assert.Equal(t, migrationBlockNum, node.ChainDB().MigrationBlockNumber())
		t.Log("Checked migration status while migration in on process")
	}

	waitMigrationEnds(t, node)

	// state trie path should not be "statetrie" in miscDB
	newPathKey := append([]byte("databaseDirectory"), common.Int64ToByteBigEndian(uint64(database.StateTrieDB))...)
	dir, err := miscDB.Get(newPathKey)
	assert.NoError(t, err)
	assert.NotEqual(t, "statetrie", string(dir))
}

func newSimpleBlockchain(t *testing.T, numAccounts int) (*node.Node, *cn.CN, *TestAccountType, *big.Int, string, *TestAccountType, []*TestAccountType, []*TestAccountType) {
	if testing.Verbose() {
		enableLog() // Change verbosity level in the function if needed
	}

	t.Log("=========== create blockchain ==============")
	fullNode, node, validator, chainId, workspace := newBlockchain(t)
	richAccount, accounts, contractAccounts := createAccount(t, numAccounts, validator)
	time.Sleep(10 * time.Second)

	return fullNode, node, validator, chainId, workspace, richAccount, accounts, contractAccounts
}

func startMigration(t *testing.T, node *cn.CN) {
	waitMigrationEnds(t, node)

	t.Log("=========== migrate trie ==============")
	currentHeader := node.BlockChain().CurrentBlock().Header()
	err := node.BlockChain().StateCache().TrieDB().Commit(currentHeader.Root, false, currentHeader.Number.Uint64())
	assert.NoError(t, err)

	err = node.BlockChain().StartStateMigration(currentHeader.Number.Uint64(), currentHeader.Root)
	assert.NoError(t, err)

	err = node.BlockChain().StateCache().TrieDB().Cap(0)
	assert.NoError(t, err)
}

func restartNode(t *testing.T, fullNode *node.Node, node *cn.CN, workspace string, validator *TestAccountType) (*node.Node, *cn.CN) {
	if err := fullNode.Stop(); err != nil {
		t.Fatal(err)
	}
	t.Log("=========== stopped node ==============")
	time.Sleep(5 * time.Second)
	t.Log("=========== starting node ==============")
	newFullNode, newNode := newKlaytnNode(t, workspace, validator)
	if err := newNode.StartMining(false); err != nil {
		t.Fatal()
	}
	time.Sleep(5 * time.Second)

	return newFullNode, newNode
}

func waitMigrationEnds(t *testing.T, node *cn.CN) {
	for node.ChainDB().InMigration() {
		t.Log("state trie migration is processing; sleep for a second before a new migration")
		time.Sleep(time.Second)
	}
}
