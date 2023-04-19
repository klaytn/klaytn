// Copyright 2022 The klaytn Authors
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

package blockchain

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

func createLocalTestDB(t *testing.T) (string, database.DBManager) {
	dir, err := ioutil.TempDir("", "klaytn-test-migration")
	if err != nil {
		t.Fatalf("failed to create a database: %v", err)
	}
	dbc := &database.DBConfig{Dir: dir, DBType: database.LevelDB, LevelDBCacheSize: 128, OpenFilesLimit: 128}
	db := database.NewDBManager(dbc)
	return dir, db
}

func TestBlockChain_migrateState(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	dir, testdb := createLocalTestDB(t)
	defer os.RemoveAll(dir)

	var (
		gspec = &Genesis{Config: params.TestChainConfig}
		_     = gspec.MustCommit(testdb)
	)

	chain, err := NewBlockChain(testdb, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	if err != nil {
		t.Fatalf("Failed to create local chain, %v", err)
	}
	defer chain.Stop()

	chain.testMigrationHook = func() {
		// sleep 2 seconds to write byte codes while migrating statedb
		t.Log("taking 2 seconds for testing purpose...")
		time.Sleep(2 * time.Second)
	}

	if err := chain.PrepareStateMigration(); err != nil {
		t.Fatalf("failed to prepare state migration: %v", err)
	}

	b := chain.CurrentBlock()
	if err := chain.StartStateMigration(b.NumberU64(), b.Root()); err != nil {
		t.Fatalf("failed to start state migration: %v", err)
	}

	// write code while statedb migration
	expectedCode := []byte{0x60, 0x80, 0x60, 0x40}
	codeHash := common.BytesToHash(expectedCode).ToRootExtHash()
	chain.stateCache.TrieDB().DiskDB().WriteCode(codeHash, expectedCode)

	for chain.db.InMigration() {
		t.Log("It is in migration, sleeping 1 second")
		time.Sleep(1 * time.Second)
	}

	// read code after deleting the original db
	actualCode := chain.db.ReadCode(codeHash)
	if actualCode == nil || bytes.Compare(expectedCode, actualCode) != 0 {
		t.Fatalf("mismatch bytecodes: (expected: %v, actual: %v)", common.Bytes2Hex(expectedCode), common.Bytes2Hex(actualCode))
	}
}
