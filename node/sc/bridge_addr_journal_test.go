// Copyright 2019 The klaytn Authors
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

package sc

import (
	"os"
	"path"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

func TestBridgeJournal(t *testing.T) {

	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), "test.rlp")); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	journal := newBridgeAddrJournal(path.Join(os.TempDir(), "test.rlp"))
	if err := journal.load(func(journal BridgeJournal) error {
		t.Log("Local address ", journal.ChildAddress.Hex())
		t.Log("Remote address ", journal.ParentAddress.Hex())
		t.Log("Subscribed", journal.Subscribed)
		return nil
	}); err != nil {
		t.Fatalf("fail to load journal %v", err)
	}
	if err := journal.rotate([]*BridgeJournal{}); err != nil {
		t.Fatalf("fail to rotate journal %v", err)
	}

	err := journal.insert(common.BytesToAddress([]byte("test1")), common.BytesToAddress([]byte("test2")))
	if err != nil {
		t.Fatalf("fail to insert address %v", err)
	}
	err = journal.insert(common.BytesToAddress([]byte("test2")), common.BytesToAddress([]byte("test3")))
	if err != nil {
		t.Fatalf("fail to insert address %v", err)
	}
	err = journal.insert(common.BytesToAddress([]byte("test3")), common.BytesToAddress([]byte("test1")))
	if err != nil {
		t.Fatalf("fail to insert address %v", err)
	}

	if err := journal.close(); err != nil {
		t.Fatalf("fail to close file %v", err)
	}

	journal = newBridgeAddrJournal(path.Join(os.TempDir(), "test.rlp"))

	if err := journal.load(func(journal BridgeJournal) error {
		switch address := journal.ChildAddress.Hex(); address {
		case "0x0000000000000000000000000000007465737431":
			if journal.ParentAddress.Hex() != "0x0000000000000000000000000000007465737432" {
				t.Fatalf("unknown remoteAddress")
			}
		case "0x0000000000000000000000000000007465737432":
			if journal.ParentAddress.Hex() != "0x0000000000000000000000000000007465737433" {
				t.Fatalf("unknown remoteAddress")
			}
		case "0x0000000000000000000000000000007465737433":
			if journal.ParentAddress.Hex() != "0x0000000000000000000000000000007465737431" {
				t.Fatalf("unknown remoteAddress")
			}
		default:
			t.Fatalf("unknown localAddress")
		}
		return nil
	}); err != nil {
		t.Fatalf("fail to load journal %v", err)
	}
}

// TestBridgeJournalCache tests the journal cache.
func TestBridgeJournalCache(t *testing.T) {
	defer func() {
		if err := os.Remove(path.Join(os.TempDir(), "test.rlp")); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Step 1: Make new journal
	journals := newBridgeAddrJournal(path.Join(os.TempDir(), "test.rlp"))

	if err := journals.load(func(journal BridgeJournal) error { return nil }); err != nil {
		t.Fatalf("fail to load journal %v", err)
	}
	if err := journals.rotate([]*BridgeJournal{}); err != nil {
		t.Fatalf("fail to rotate journal %v", err)
	}

	localAddr := common.BytesToAddress([]byte("test1"))
	remoteAddr := common.BytesToAddress([]byte("test2"))

	err := journals.insert(localAddr, remoteAddr)
	if err != nil {
		t.Fatalf("fail to insert address %v", err)
	}

	assert.Equal(t, 1, len(journals.cache))
	for _, journal := range journals.cache {
		assert.Equal(t, localAddr, journal.ChildAddress)
		assert.Equal(t, remoteAddr, journal.ParentAddress)
	}
}
