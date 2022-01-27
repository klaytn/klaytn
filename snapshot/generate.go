// Modifications Copyright 2021 The klaytn Authors
// Copyright 2019 The go-ethereum Authors
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
// This file is derived from core/state/snapshot/generate.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/rcrowley/go-metrics"
)

var (
	// errMissingTrie is returned if the target trie is missing while the generation
	// is running. In this case the generation is aborted and wait the new signal.
	errMissingTrie = errors.New("missing trie")
)

var (
	// snapAccountProveCounter measures time spent on the account proving
	snapAccountProveCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/prove", nil)
	// snapAccountTrieReadCounter measures time spent on the account trie iteration
	snapAccountTrieReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/trieread", nil)
	// snapAccountSnapReadCounter measues time spent on the snapshot account iteration
	snapAccountSnapReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/snapread", nil)
	// snapAccountWriteCounter measures time spent on writing/updating/deleting accounts
	snapAccountWriteCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/write", nil)
	// snapStorageProveCounter measures time spent on storage proving
	snapStorageProveCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/prove", nil)
	// snapStorageTrieReadCounter measures time spent on the storage trie iteration
	snapStorageTrieReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/trieread", nil)
	// snapStorageSnapReadCounter measures time spent on the snapshot storage iteration
	snapStorageSnapReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/snapread", nil)
	// snapStorageWriteCounter measures time spent on writing/updating/deleting storages
	snapStorageWriteCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/write", nil)
)

// generatorStats is a collection of statistics gathered by the snapshot generator
// for logging purposes.
type generatorStats struct {
	origin   uint64             // Origin prefix where generation started
	start    time.Time          // Timestamp when generation started
	accounts uint64             // Number of accounts indexed(generated or recovered)
	slots    uint64             // Number of storage slots indexed(generated or recovered)
	storage  common.StorageSize // Total account and storage slot size(generation or recovery)
}

// Log creates an contextual log with the given message and the context pulled
// from the internally maintained statistics.
func (gs *generatorStats) Log(msg string, root common.Hash, marker []byte) {
	var ctx []interface{}
	if root != (common.Hash{}) {
		ctx = append(ctx, []interface{}{"root", root}...)
	}
	// Figure out whether we're after or within an account
	switch len(marker) {
	case common.HashLength:
		ctx = append(ctx, []interface{}{"at", common.BytesToHash(marker)}...)
	case 2 * common.HashLength:
		ctx = append(ctx, []interface{}{
			"in", common.BytesToHash(marker[:common.HashLength]),
			"at", common.BytesToHash(marker[common.HashLength:]),
		}...)
	}
	// Add the usual measurements
	ctx = append(ctx, []interface{}{
		"accounts", gs.accounts,
		"slots", gs.slots,
		"storage", gs.storage,
		"elapsed", common.PrettyDuration(time.Since(gs.start)),
	}...)
	// Calculate the estimated indexing time based on current stats
	if len(marker) > 0 {
		if done := binary.BigEndian.Uint64(marker[:8]) - gs.origin; done > 0 {
			left := math.MaxUint64 - binary.BigEndian.Uint64(marker[:8])

			speed := done/uint64(time.Since(gs.start)/time.Millisecond+1) + 1 // +1s to avoid division by zero
			ctx = append(ctx, []interface{}{
				"eta", common.PrettyDuration(time.Duration(left/speed) * time.Millisecond),
			}...)
		}
	}
	logger.Info(msg, ctx...)
}

// proofResult contains the output of range proving which can be used
// for further processing regardless if it is successful or not.
type proofResult struct {
	keys     [][]byte      // The key set of all elements being iterated, even proving is failed
	vals     [][]byte      // The val set of all elements being iterated, even proving is failed
	diskMore bool          // Set when the database has extra snapshot states since last iteration
	trieMore bool          // Set when the trie has extra snapshot states(only meaningful for successful proving)
	proofErr error         // Indicator whether the given state range is valid or not
	tr       *statedb.Trie // The trie, in case the trie was resolved by the prover (may be nil)
}

// valid returns the indicator that range proof is successful or not.
func (result *proofResult) valid() bool {
	return result.proofErr == nil
}

// last returns the last verified element key regardless of whether the range proof is
// successful or not. Nil is returned if nothing involved in the proving.
func (result *proofResult) last() []byte {
	var last []byte
	if len(result.keys) > 0 {
		last = result.keys[len(result.keys)-1]
	}
	return last
}

// forEach iterates all the visited elements and applies the given callback on them.
// The iteration is aborted if the callback returns non-nil error.
func (result *proofResult) forEach(callback func(key []byte, val []byte) error) error {
	for i := 0; i < len(result.keys); i++ {
		key, val := result.keys[i], result.vals[i]
		if err := callback(key, val); err != nil {
			return err
		}
	}
	return nil
}

// proveRange proves the snapshot segment with particular prefix is "valid".
// The iteration start point will be assigned if the iterator is restored from
// the last interruption. Max will be assigned in order to limit the maximum
// amount of data involved in each iteration.
//
// The proof result will be returned if the range proving is finished, otherwise
// the error will be returned to abort the entire procedure.
func (dl *diskLayer) proveRange(stats *generatorStats, root common.Hash, prefix []byte, kind string, origin []byte, max int, valueConvertFn func([]byte) ([]byte, error)) (*proofResult, error) {
	var (
		keys     [][]byte
		vals     [][]byte
		proof    = database.NewMemoryDBManager()
		diskMore = false
	)

	iter := dl.diskdb.GetSnapshotDB().NewIterator(prefix, origin)
	defer iter.Release()

	var start = time.Now()
	for iter.Next() {
		key := iter.Key()
		if len(key) != len(prefix)+common.HashLength {
			continue
		}
		if len(keys) == max {
			// Break if we've reached the max size, and signal that we're not
			// done yet.
			diskMore = true
			break
		}
		keys = append(keys, common.CopyBytes(key[len(prefix):]))

		if valueConvertFn == nil {
			vals = append(vals, common.CopyBytes(iter.Value()))
		} else {
			val, err := valueConvertFn(iter.Value())
			if err != nil {
				// Special case, the state data is corrupted (invalid slim-format account),
				// don't abort the entire procedure directly. Instead, let the fallback
				// generation to heal the invalid data.
				//
				// Here append the original value to ensure that the number of key and
				// value are the same.
				vals = append(vals, common.CopyBytes(iter.Value()))
				logger.Error("Failed to convert account state data", "err", err)
			} else {
				vals = append(vals, val)
			}
		}
	}
	// Update metrics for database iteration and merkle proving
	if kind == "storage" {
		snapStorageSnapReadCounter.Inc(time.Since(start).Nanoseconds())
	} else {
		snapAccountSnapReadCounter.Inc(time.Since(start).Nanoseconds())
	}
	defer func(start time.Time) {
		if kind == "storage" {
			snapStorageProveCounter.Inc(time.Since(start).Nanoseconds())
		} else {
			snapAccountProveCounter.Inc(time.Since(start).Nanoseconds())
		}
	}(time.Now())

	// The snap state is exhausted, pass the entire key/val set for verification
	if origin == nil && !diskMore {
		var (
			dbm    = database.NewMemoryDBManager()
			triedb = statedb.NewDatabase(dbm)
		)
		tr, _ := statedb.NewTrie(common.Hash{}, triedb)
		for i, key := range keys {
			tr.TryUpdate(key, vals[i])
		}
		if gotRoot := tr.Hash(); gotRoot != root {
			return &proofResult{
				keys:     keys,
				vals:     vals,
				proofErr: fmt.Errorf("wrong root: have %#x want %#x", gotRoot, root),
			}, nil
		}
		return &proofResult{keys: keys, vals: vals}, nil
	}
	// Snap state is chunked, generate edge proofs for verification.
	tr, err := statedb.NewTrie(root, dl.triedb)
	if err != nil {
		stats.Log("Trie missing, state snapshotting paused", dl.root, dl.genMarker)
		return nil, errMissingTrie
	}
	// Firstly find out the key of last iterated element.
	var last []byte
	if len(keys) > 0 {
		last = keys[len(keys)-1]
	}
	// Generate the Merkle proofs for the first and last element
	if origin == nil {
		origin = common.Hash{}.Bytes()
	}
	if err := tr.Prove(origin, 0, proof); err != nil {
		logger.Debug("Failed to prove range", "kind", kind, "origin", origin, "err", err)
		return &proofResult{
			keys:     keys,
			vals:     vals,
			diskMore: diskMore,
			proofErr: err,
			tr:       tr,
		}, nil
	}
	if last != nil {
		if err := tr.Prove(last, 0, proof); err != nil {
			logger.Debug("Failed to prove range", "kind", kind, "last", last, "err", err)
			return &proofResult{
				keys:     keys,
				vals:     vals,
				diskMore: diskMore,
				proofErr: err,
				tr:       tr,
			}, nil
		}
	}
	// Verify the snapshot segment with range prover, ensure that all flat states
	// in this range correspond to merkle trie.
	// TODO-Klaytn-Snapshot this comment has to be uncommented after VerifyRangeProof implementation
	//cont, err := statedb.VerifyRangeProof(root, origin, last, keys, vals, proof)
	return &proofResult{
			keys:     keys,
			vals:     vals,
			diskMore: diskMore,
			// TODO-Klaytn-Snapshot this comment has to be uncommented after VerifyRangeProof implementation
			//trieMore: cont,
			proofErr: err,
			tr:       tr},
		nil
}
