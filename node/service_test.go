// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from node/service_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/storage/database"
)

// Tests that databases are correctly created persistent or ephemeral based on
// the configured service context.
func TestContextDatabases(t *testing.T) {
	// Create a temporary folder and ensure no database is contained within
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temporary data directory: %v", err)
	}
	defer os.RemoveAll(dir)

	if _, err := os.Stat(filepath.Join(dir, "database")); err == nil {
		t.Fatalf("non-created database already exists")
	}
	// Request the opening/creation of a database and ensure it persists to disk
	ctx := NewServiceContext(&Config{Name: "unit-test", DataDir: dir}, map[reflect.Type]Service{}, &event.TypeMux{}, &accounts.Manager{})
	dbc := &database.DBConfig{Dir: "persistent", DBType: database.LevelDB,
		LevelDBCacheSize: 0, OpenFilesLimit: 0}
	db := ctx.OpenDatabase(dbc)
	db.Close()

	if _, err := os.Stat(filepath.Join(dir, "unit-test", "persistent")); err != nil {
		t.Fatalf("persistent database doesn't exists: %v", err)
	}
	// Request th opening/creation of an ephemeral database and ensure it's not persisted
	ctx = NewServiceContext(&Config{DataDir: ""}, map[reflect.Type]Service{}, &event.TypeMux{}, &accounts.Manager{})
	dbc = &database.DBConfig{Dir: "ephemeral", DBType: database.LevelDB,
		LevelDBCacheSize: 0, OpenFilesLimit: 0}

	db = ctx.OpenDatabase(dbc)
	db.Close()

	if _, err := os.Stat(filepath.Join(dir, "ephemeral")); err == nil {
		t.Fatalf("ephemeral database exists")
	}
}

// Tests that already constructed services can be retrieves by later ones.
func TestContextServices(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
	// Define a verifier that ensures a NoopA is before it and NoopB after
	verifier := func(ctx *ServiceContext) (Service, error) {
		var objA *NoopServiceA
		if ctx.Service(&objA) != nil {
			return nil, fmt.Errorf("former service not found")
		}
		var objB *NoopServiceB
		if err := ctx.Service(&objB); err != ErrServiceUnknown {
			return nil, fmt.Errorf("latters lookup error mismatch: have %v, want %v", err, ErrServiceUnknown)
		}
		return new(NoopService), nil
	}
	// Register the collection of services
	if err := stack.Register(NewNoopServiceA); err != nil {
		t.Fatalf("former failed to register service: %v", err)
	}
	if err := stack.Register(verifier); err != nil {
		t.Fatalf("failed to register service verifier: %v", err)
	}
	if err := stack.Register(NewNoopServiceB); err != nil {
		t.Fatalf("latter failed to register service: %v", err)
	}
	// Start the protocol stack and ensure services are constructed in order
	if err := stack.Start(); err != nil {
		t.Fatalf("failed to start stack: %v", err)
	}
	defer stack.Stop()
}
