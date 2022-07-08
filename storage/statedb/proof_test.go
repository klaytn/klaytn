// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
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
// This file is derived from trie/proof_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package statedb

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"sort"
	"testing"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/storage/database"
)

func init() {
	mrand.Seed(time.Now().Unix())
}

func TestProof(t *testing.T) {
	trie, vals := randomTrie(500)
	root := trie.Hash()
	for _, kv := range vals {
		proofs := database.NewMemoryDBManager()
		if trie.Prove(kv.k, 0, proofs) != nil {
			t.Fatalf("missing key %x while constructing proof", kv.k)
		}
		val, err, _ := VerifyProof(root, kv.k, proofs)
		if err != nil {
			t.Fatalf("VerifyProof error for key %x: %v\nraw proof: %v", kv.k, err, proofs)
		}
		if !bytes.Equal(val, kv.v) {
			t.Fatalf("VerifyProof returned wrong value for key %x: got %x, want %x", kv.k, val, kv.v)
		}
	}
}

func TestOneElementProof(t *testing.T) {
	trie := new(Trie)
	updateString(trie, "k", "v")
	memDBManager := database.NewMemoryDBManager()
	proofs := memDBManager.GetMemDB()
	trie.Prove([]byte("k"), 0, memDBManager)
	if len(proofs.Keys()) != 1 {
		t.Error("proof should have one element")
	}
	val, err, _ := VerifyProof(trie.Hash(), []byte("k"), memDBManager)
	if err != nil {
		t.Fatalf("VerifyProof error: %v\nproof hashes: %v", err, proofs.Keys())
	}
	if !bytes.Equal(val, []byte("v")) {
		t.Fatalf("VerifyProof returned wrong value: got %x, want 'k'", val)
	}
}

func TestVerifyBadProof(t *testing.T) {
	trie, vals := randomTrie(800)
	root := trie.Hash()
	for _, kv := range vals {
		memDBManager := database.NewMemoryDBManager()
		proofs := memDBManager.GetMemDB()
		trie.Prove(kv.k, 0, memDBManager)
		if len(proofs.Keys()) == 0 {
			t.Fatal("zero length proof")
		}
		keys := proofs.Keys()
		key := keys[mrand.Intn(len(keys))]
		node, _ := proofs.Get(key)
		proofs.Delete(key)
		mutateByte(node)
		proofs.Put(crypto.Keccak256(node), node)
		if _, err, _ := VerifyProof(root, kv.k, memDBManager); err == nil {
			t.Fatalf("expected proof to fail for key %x", kv.k)
		}
	}
}

type entrySlice []*kv

func (p entrySlice) Len() int           { return len(p) }
func (p entrySlice) Less(i, j int) bool { return bytes.Compare(p[i].k, p[j].k) < 0 }
func (p entrySlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// TestRangeProof tests normal range proof with both edge proofs
// as the existent proof. The test cases are generated randomly.
func TestRangeProof(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)
	for i := 0; i < 500; i++ {
		start := mrand.Intn(len(entries))
		end := mrand.Intn(len(entries)-start) + start + 1

		proof := database.NewMemoryDBManager()
		if err := trie.Prove(entries[start].k, 0, proof); err != nil {
			t.Fatalf("Failed to prove the first node %v", err)
		}
		if err := trie.Prove(entries[end-1].k, 0, proof); err != nil {
			t.Fatalf("Failed to prove the last node %v", err)
		}
		var keys [][]byte
		var vals [][]byte
		for i := start; i < end; i++ {
			keys = append(keys, entries[i].k)
			vals = append(vals, entries[i].v)
		}
		_, err := VerifyRangeProof(trie.Hash(), keys[0], keys[len(keys)-1], keys, vals, proof)
		if err != nil {
			t.Fatalf("Case %d(%d->%d) expect no error, got %v", i, start, end-1, err)
		}
	}
}

// TestRangeProof tests normal range proof with two non-existent proofs.
// The test cases are generated randomly.
func TestRangeProofWithNonExistentProof(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)
	for i := 0; i < 500; i++ {
		start := mrand.Intn(len(entries))
		end := mrand.Intn(len(entries)-start) + start + 1
		proof := database.NewMemoryDBManager()

		// Short circuit if the decreased key is same with the previous key
		first := decreseKey(common.CopyBytes(entries[start].k))
		if start != 0 && bytes.Equal(first, entries[start-1].k) {
			continue
		}
		// Short circuit if the decreased key is underflow
		if bytes.Compare(first, entries[start].k) > 0 {
			continue
		}
		// Short circuit if the increased key is same with the next key
		last := increseKey(common.CopyBytes(entries[end-1].k))
		if end != len(entries) && bytes.Equal(last, entries[end].k) {
			continue
		}
		// Short circuit if the increased key is overflow
		if bytes.Compare(last, entries[end-1].k) < 0 {
			continue
		}
		if err := trie.Prove(first, 0, proof); err != nil {
			t.Fatalf("Failed to prove the first node %v", err)
		}
		if err := trie.Prove(last, 0, proof); err != nil {
			t.Fatalf("Failed to prove the last node %v", err)
		}
		var keys [][]byte
		var vals [][]byte
		for i := start; i < end; i++ {
			keys = append(keys, entries[i].k)
			vals = append(vals, entries[i].v)
		}
		_, err := VerifyRangeProof(trie.Hash(), first, last, keys, vals, proof)
		if err != nil {
			t.Fatalf("Case %d(%d->%d) expect no error, got %v", i, start, end-1, err)
		}
	}
	// Special case, two edge proofs for two edge key.
	proof := database.NewMemoryDBManager()
	first := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000").Bytes()
	last := common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff").Bytes()
	if err := trie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	var k [][]byte
	var v [][]byte
	for i := 0; i < len(entries); i++ {
		k = append(k, entries[i].k)
		v = append(v, entries[i].v)
	}
	_, err := VerifyRangeProof(trie.Hash(), first, last, k, v, proof)
	if err != nil {
		t.Fatal("Failed to verify whole rang with non-existent edges")
	}
}

// TestRangeProofWithInvalidNonExistentProof tests such scenarios:
// - There exists a gap between the first element and the left edge proof
// - There exists a gap between the last element and the right edge proof
func TestRangeProofWithInvalidNonExistentProof(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)

	// Case 1
	start, end := 100, 200
	first := decreseKey(common.CopyBytes(entries[start].k))

	proof := database.NewMemoryDBManager()
	if err := trie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(entries[end-1].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	start = 105 // Gap created
	k := make([][]byte, 0)
	v := make([][]byte, 0)
	for i := start; i < end; i++ {
		k = append(k, entries[i].k)
		v = append(v, entries[i].v)
	}
	_, err := VerifyRangeProof(trie.Hash(), first, k[len(k)-1], k, v, proof)
	if err == nil {
		t.Fatalf("Expected to detect the error, got nil")
	}

	// Case 2
	start, end = 100, 200
	last := increseKey(common.CopyBytes(entries[end-1].k))
	proof = database.NewMemoryDBManager()
	if err := trie.Prove(entries[start].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	end = 195 // Capped slice
	k = make([][]byte, 0)
	v = make([][]byte, 0)
	for i := start; i < end; i++ {
		k = append(k, entries[i].k)
		v = append(v, entries[i].v)
	}
	_, err = VerifyRangeProof(trie.Hash(), k[0], last, k, v, proof)
	if err == nil {
		t.Fatalf("Expected to detect the error, got nil")
	}
}

// TestOneElementRangeProof tests the proof with only one
// element. The first edge proof can be existent one or
// non-existent one.
func TestOneElementRangeProof(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)

	// One element with existent edge proof, both edge proofs
	// point to the SAME key.
	start := 1000
	proof := database.NewMemoryDBManager()
	if err := trie.Prove(entries[start].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	_, err := VerifyRangeProof(trie.Hash(), entries[start].k, entries[start].k, [][]byte{entries[start].k}, [][]byte{entries[start].v}, proof)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// One element with left non-existent edge proof
	start = 1000
	first := decreseKey(common.CopyBytes(entries[start].k))
	proof = database.NewMemoryDBManager()
	if err := trie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(entries[start].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err = VerifyRangeProof(trie.Hash(), first, entries[start].k, [][]byte{entries[start].k}, [][]byte{entries[start].v}, proof)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// One element with right non-existent edge proof
	start = 1000
	last := increseKey(common.CopyBytes(entries[start].k))
	proof = database.NewMemoryDBManager()
	if err := trie.Prove(entries[start].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err = VerifyRangeProof(trie.Hash(), entries[start].k, last, [][]byte{entries[start].k}, [][]byte{entries[start].v}, proof)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// One element with two non-existent edge proofs
	start = 1000
	first, last = decreseKey(common.CopyBytes(entries[start].k)), increseKey(common.CopyBytes(entries[start].k))
	proof = database.NewMemoryDBManager()
	if err := trie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err = VerifyRangeProof(trie.Hash(), first, last, [][]byte{entries[start].k}, [][]byte{entries[start].v}, proof)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test the mini trie with only a single element.
	tinyTrie := new(Trie)
	entry := &kv{randBytes(32), randBytes(20), false}
	tinyTrie.Update(entry.k, entry.v)

	first = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000").Bytes()
	last = entry.k
	proof = database.NewMemoryDBManager()
	if err := tinyTrie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := tinyTrie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err = VerifyRangeProof(tinyTrie.Hash(), first, last, [][]byte{entry.k}, [][]byte{entry.v}, proof)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestAllElementsProof tests the range proof with all elements.
// The edge proofs can be nil.
func TestAllElementsProof(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)

	var k [][]byte
	var v [][]byte
	for i := 0; i < len(entries); i++ {
		k = append(k, entries[i].k)
		v = append(v, entries[i].v)
	}
	_, err := VerifyRangeProof(trie.Hash(), nil, nil, k, v, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// With edge proofs, it should still work.
	proof := database.NewMemoryDBManager()
	if err := trie.Prove(entries[0].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(entries[len(entries)-1].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err = VerifyRangeProof(trie.Hash(), k[0], k[len(k)-1], k, v, proof)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Even with non-existent edge proofs, it should still work.
	proof = database.NewMemoryDBManager()
	first := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000").Bytes()
	last := common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff").Bytes()
	if err := trie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err = VerifyRangeProof(trie.Hash(), first, last, k, v, proof)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestSingleSideRangeProof tests the range starts from zero.
func TestSingleSideRangeProof(t *testing.T) {
	for i := 0; i < 64; i++ {
		trie := new(Trie)
		entries := entrySlice{}
		for i := 0; i < 4096; i++ {
			value := &kv{randBytes(32), randBytes(20), false}
			trie.Update(value.k, value.v)
			entries = append(entries, value)
		}
		sort.Sort(entries)

		cases := []int{0, 1, 50, 100, 1000, 2000, len(entries) - 1}
		for _, pos := range cases {
			proof := database.NewMemoryDBManager()
			if err := trie.Prove(common.Hash{}.Bytes(), 0, proof); err != nil {
				t.Fatalf("Failed to prove the first node %v", err)
			}
			if err := trie.Prove(entries[pos].k, 0, proof); err != nil {
				t.Fatalf("Failed to prove the first node %v", err)
			}
			k := make([][]byte, 0)
			v := make([][]byte, 0)
			for i := 0; i <= pos; i++ {
				k = append(k, entries[i].k)
				v = append(v, entries[i].v)
			}
			_, err := VerifyRangeProof(trie.Hash(), common.Hash{}.Bytes(), k[len(k)-1], k, v, proof)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		}
	}
}

// TestReverseSingleSideRangeProof tests the range ends with 0xffff...fff.
func TestReverseSingleSideRangeProof(t *testing.T) {
	for i := 0; i < 64; i++ {
		trie := new(Trie)
		entries := entrySlice{}
		for i := 0; i < 4096; i++ {
			value := &kv{randBytes(32), randBytes(20), false}
			trie.Update(value.k, value.v)
			entries = append(entries, value)
		}
		sort.Sort(entries)

		cases := []int{0, 1, 50, 100, 1000, 2000, len(entries) - 1}
		for _, pos := range cases {
			proof := database.NewMemoryDBManager()
			if err := trie.Prove(entries[pos].k, 0, proof); err != nil {
				t.Fatalf("Failed to prove the first node %v", err)
			}
			last := common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
			if err := trie.Prove(last.Bytes(), 0, proof); err != nil {
				t.Fatalf("Failed to prove the last node %v", err)
			}
			k := make([][]byte, 0)
			v := make([][]byte, 0)
			for i := pos; i < len(entries); i++ {
				k = append(k, entries[i].k)
				v = append(v, entries[i].v)
			}
			_, err := VerifyRangeProof(trie.Hash(), k[0], last.Bytes(), k, v, proof)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		}
	}
}

// TestBadRangeProof tests a few cases which the proof is wrong.
// The prover is expected to detect the error.
func TestBadRangeProof(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)

	for i := 0; i < 500; i++ {
		start := mrand.Intn(len(entries))
		end := mrand.Intn(len(entries)-start) + start + 1
		proof := database.NewMemoryDBManager()
		if err := trie.Prove(entries[start].k, 0, proof); err != nil {
			t.Fatalf("Failed to prove the first node %v", err)
		}
		if err := trie.Prove(entries[end-1].k, 0, proof); err != nil {
			t.Fatalf("Failed to prove the last node %v", err)
		}
		var keys [][]byte
		var vals [][]byte
		for i := start; i < end; i++ {
			keys = append(keys, entries[i].k)
			vals = append(vals, entries[i].v)
		}
		first, last := keys[0], keys[len(keys)-1]
		testcase := mrand.Intn(6)
		var index int
		switch testcase {
		case 0:
			// Modified key
			index = mrand.Intn(end - start)
			keys[index] = randBytes(32) // In theory it can't be same
		case 1:
			// Modified val
			index = mrand.Intn(end - start)
			vals[index] = randBytes(20) // In theory it can't be same
		case 2:
			// Gapped entry slice
			index = mrand.Intn(end - start)
			if (index == 0 && start < 100) || (index == end-start-1 && end <= 100) {
				continue
			}
			keys = append(keys[:index], keys[index+1:]...)
			vals = append(vals[:index], vals[index+1:]...)
		case 3:
			// Out of order
			index1 := mrand.Intn(end - start)
			index2 := mrand.Intn(end - start)
			if index1 == index2 {
				continue
			}
			keys[index1], keys[index2] = keys[index2], keys[index1]
			vals[index1], vals[index2] = vals[index2], vals[index1]
		case 4:
			// Set random key to nil, do nothing
			index = mrand.Intn(end - start)
			keys[index] = nil
		case 5:
			// Set random value to nil, deletion
			index = mrand.Intn(end - start)
			vals[index] = nil
		}
		_, err := VerifyRangeProof(trie.Hash(), first, last, keys, vals, proof)
		if err == nil {
			t.Fatalf("%d Case %d index %d range: (%d->%d) expect error, got nil", i, testcase, index, start, end-1)
		}
	}
}

// TestGappedRangeProof focuses on the small trie with embedded nodes.
// If the gapped node is embedded in the trie, it should be detected too.
func TestGappedRangeProof(t *testing.T) {
	trie := new(Trie)
	var entries []*kv // Sorted entries
	for i := byte(0); i < 10; i++ {
		value := &kv{common.LeftPadBytes([]byte{i}, 32), []byte{i}, false}
		trie.Update(value.k, value.v)
		entries = append(entries, value)
	}
	first, last := 2, 8
	proof := database.NewMemoryDBManager()
	if err := trie.Prove(entries[first].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(entries[last-1].k, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	var keys [][]byte
	var vals [][]byte
	for i := first; i < last; i++ {
		if i == (first+last)/2 {
			continue
		}
		keys = append(keys, entries[i].k)
		vals = append(vals, entries[i].v)
	}
	_, err := VerifyRangeProof(trie.Hash(), keys[0], keys[len(keys)-1], keys, vals, proof)
	if err == nil {
		t.Fatal("expect error, got nil")
	}
}

// TestSameSideProofs tests the element is not in the range covered by proofs
func TestSameSideProofs(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)

	pos := 1000
	first := decreseKey(common.CopyBytes(entries[pos].k))
	first = decreseKey(first)
	last := decreseKey(common.CopyBytes(entries[pos].k))

	proof := database.NewMemoryDBManager()
	if err := trie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err := VerifyRangeProof(trie.Hash(), first, last, [][]byte{entries[pos].k}, [][]byte{entries[pos].v}, proof)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	first = increseKey(common.CopyBytes(entries[pos].k))
	last = increseKey(common.CopyBytes(entries[pos].k))
	last = increseKey(last)

	proof = database.NewMemoryDBManager()
	if err := trie.Prove(first, 0, proof); err != nil {
		t.Fatalf("Failed to prove the first node %v", err)
	}
	if err := trie.Prove(last, 0, proof); err != nil {
		t.Fatalf("Failed to prove the last node %v", err)
	}
	_, err = VerifyRangeProof(trie.Hash(), first, last, [][]byte{entries[pos].k}, [][]byte{entries[pos].v}, proof)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestHasRightElement(t *testing.T) {
	trie := new(Trie)
	var entries entrySlice
	for i := 0; i < 4096; i++ {
		value := &kv{randBytes(32), randBytes(20), false}
		trie.Update(value.k, value.v)
		entries = append(entries, value)
	}
	sort.Sort(entries)

	cases := []struct {
		start   int
		end     int
		hasMore bool
	}{
		{-1, 1, true}, // single element with non-existent left proof
		{0, 1, true},  // single element with existent left proof
		{0, 10, true},
		{50, 100, true},
		{50, len(entries), false},               // No more element expected
		{len(entries) - 1, len(entries), false}, // Single last element with two existent proofs(point to same key)
		{len(entries) - 1, -1, false},           // Single last element with non-existent right proof
		{0, len(entries), false},                // The whole set with existent left proof
		{-1, len(entries), false},               // The whole set with non-existent left proof
		{-1, -1, false},                         // The whole set with non-existent left/right proof
	}
	for _, c := range cases {
		var (
			firstKey []byte
			lastKey  []byte
			start    = c.start
			end      = c.end
			proof    = database.NewMemoryDBManager()
		)
		if c.start == -1 {
			firstKey, start = common.Hash{}.Bytes(), 0
			if err := trie.Prove(firstKey, 0, proof); err != nil {
				t.Fatalf("Failed to prove the first node %v", err)
			}
		} else {
			firstKey = entries[c.start].k
			if err := trie.Prove(entries[c.start].k, 0, proof); err != nil {
				t.Fatalf("Failed to prove the first node %v", err)
			}
		}
		if c.end == -1 {
			lastKey, end = common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff").Bytes(), len(entries)
			if err := trie.Prove(lastKey, 0, proof); err != nil {
				t.Fatalf("Failed to prove the first node %v", err)
			}
		} else {
			lastKey = entries[c.end-1].k
			if err := trie.Prove(entries[c.end-1].k, 0, proof); err != nil {
				t.Fatalf("Failed to prove the first node %v", err)
			}
		}
		k := make([][]byte, 0)
		v := make([][]byte, 0)
		for i := start; i < end; i++ {
			k = append(k, entries[i].k)
			v = append(v, entries[i].v)
		}
		hasMore, err := VerifyRangeProof(trie.Hash(), firstKey, lastKey, k, v, proof)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if hasMore != c.hasMore {
			t.Fatalf("Wrong hasMore indicator, want %t, got %t", c.hasMore, hasMore)
		}
	}
}

// TestEmptyRangeProof tests the range proof with "no" element.
// The first edge proof must be a non-existent proof.
func TestEmptyRangeProof(t *testing.T) {
	trie, vals := randomTrie(4096)
	entries := entrySlice{}
	for _, kv := range vals {
		entries = append(entries, kv)
	}
	sort.Sort(entries)

	cases := []struct {
		pos int
		err bool
	}{
		{len(entries) - 1, false},
		{500, true},
	}
	for _, c := range cases {
		proof := database.NewMemoryDBManager()
		first := increseKey(common.CopyBytes(entries[c.pos].k))
		if err := trie.Prove(first, 0, proof); err != nil {
			t.Fatalf("Failed to prove the first node %v", err)
		}
		_, err := VerifyRangeProof(trie.Hash(), first, nil, nil, nil, proof)
		if c.err && err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !c.err && err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
}

// TestBloatedProof tests a malicious proof, where the proof is more or less the
// whole trie. Previously we didn't accept such packets, but the new APIs do, so
// lets leave this test as a bit weird, but present.
func TestBloatedProof(t *testing.T) {
	// Use a small trie
	trie, kvs := nonRandomTrie(100)
	entries := entrySlice{}
	for _, kv := range kvs {
		entries = append(entries, kv)
	}
	sort.Sort(entries)
	var keys [][]byte
	var vals [][]byte

	proof := database.NewMemoryDBManager()
	// In the 'malicious' case, we add proofs for every single item
	// (but only one key/value pair used as leaf)
	for i, entry := range entries {
		trie.Prove(entry.k, 0, proof)
		if i == 50 {
			keys = append(keys, entry.k)
			vals = append(vals, entry.v)
		}
	}
	// For reference, we use the same function, but _only_ prove the first
	// and last element
	want := database.NewMemoryDBManager()
	trie.Prove(keys[0], 0, want)
	trie.Prove(keys[len(keys)-1], 0, want)

	if _, err := VerifyRangeProof(trie.Hash(), keys[0], keys[len(keys)-1], keys, vals, proof); err != nil {
		t.Fatalf("expected bloated proof to succeed, got %v", err)
	}
}

// mutateByte changes one byte in b.
func mutateByte(b []byte) {
	for r := mrand.Intn(len(b)); ; {
		new := byte(mrand.Intn(255))
		if new != b[r] {
			b[r] = new
			break
		}
	}
}

func increseKey(key []byte) []byte {
	for i := len(key) - 1; i >= 0; i-- {
		key[i]++
		if key[i] != 0x0 {
			break
		}
	}
	return key
}

func decreseKey(key []byte) []byte {
	for i := len(key) - 1; i >= 0; i-- {
		key[i]--
		if key[i] != 0xff {
			break
		}
	}
	return key
}

func BenchmarkProve(b *testing.B) {
	trie, vals := randomTrie(100)
	var keys []string
	for k := range vals {
		keys = append(keys, k)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kv := vals[keys[i%len(keys)]]
		memDBManager := database.NewMemoryDBManager()
		proofs := memDBManager.GetMemDB()
		if trie.Prove(kv.k, 0, memDBManager); len(proofs.Keys()) == 0 {
			b.Fatalf("zero length proof for %x", kv.k)
		}
	}
}

func BenchmarkVerifyProof(b *testing.B) {
	trie, vals := randomTrie(100)
	root := trie.Hash()
	var keys []string
	var proofs []database.DBManager
	for k := range vals {
		keys = append(keys, k)
		proof := database.NewMemoryDBManager()
		trie.Prove([]byte(k), 0, proof)
		proofs = append(proofs, proof)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % len(keys)
		if _, err, _ := VerifyProof(root, []byte(keys[im]), proofs[im]); err != nil {
			b.Fatalf("key %x: %v", keys[im], err)
		}
	}
}

func randomTrie(n int) (*Trie, map[string]*kv) {
	trie := new(Trie)
	vals := make(map[string]*kv)
	for i := byte(0); i < 100; i++ {
		value := &kv{common.LeftPadBytes([]byte{i}, 32), []byte{i}, false}
		value2 := &kv{common.LeftPadBytes([]byte{i + 10}, 32), []byte{i}, false}
		trie.Update(value.k, value.v)
		trie.Update(value2.k, value2.v)
		vals[string(value.k)] = value
		vals[string(value2.k)] = value2
	}
	for i := 0; i < n; i++ {
		value := &kv{randBytes(32), randBytes(20), false}
		trie.Update(value.k, value.v)
		vals[string(value.k)] = value
	}
	return trie, vals
}

func randBytes(n int) []byte {
	r := make([]byte, n)
	crand.Read(r)
	return r
}

func nonRandomTrie(n int) (*Trie, map[string]*kv) {
	trie := new(Trie)
	vals := make(map[string]*kv)
	max := uint64(0xffffffffffffffff)
	for i := uint64(0); i < uint64(n); i++ {
		value := make([]byte, 32)
		key := make([]byte, 32)
		binary.LittleEndian.PutUint64(key, i)
		binary.LittleEndian.PutUint64(value, i-max)
		// value := &kv{common.LeftPadBytes([]byte{i}, 32), []byte{i}, false}
		elem := &kv{key, value, false}
		trie.Update(elem.k, elem.v)
		vals[string(elem.k)] = elem
	}
	return trie, vals
}
