// Modifications Copyright 2021 The klaytn Authors
// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/state/snapshot/snapshot.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"bytes"
	"errors"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/rcrowley/go-metrics"
)

var logger = log.NewModuleLogger(log.Snapshot)

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// emptyCode is the known hash of the empty EVM bytecode.
	emptyCode = crypto.Keccak256Hash(nil)
)

var (
	snapshotCleanAccountHitMeter   = metrics.NewRegisteredMeter("state/snapshot/clean/account/hit", nil)
	snapshotCleanAccountMissMeter  = metrics.NewRegisteredMeter("state/snapshot/clean/account/miss", nil)
	snapshotCleanAccountInexMeter  = metrics.NewRegisteredMeter("state/snapshot/clean/account/inex", nil)
	snapshotCleanAccountReadMeter  = metrics.NewRegisteredMeter("state/snapshot/clean/account/read", nil)
	snapshotCleanAccountWriteMeter = metrics.NewRegisteredMeter("state/snapshot/clean/account/write", nil)

	snapshotCleanStorageHitMeter   = metrics.NewRegisteredMeter("state/snapshot/clean/storage/hit", nil)
	snapshotCleanStorageMissMeter  = metrics.NewRegisteredMeter("state/snapshot/clean/storage/miss", nil)
	snapshotCleanStorageInexMeter  = metrics.NewRegisteredMeter("state/snapshot/clean/storage/inex", nil)
	snapshotCleanStorageReadMeter  = metrics.NewRegisteredMeter("state/snapshot/clean/storage/read", nil)
	snapshotCleanStorageWriteMeter = metrics.NewRegisteredMeter("state/snapshot/clean/storage/write", nil)

	snapshotDirtyAccountHitMeter   = metrics.NewRegisteredMeter("state/snapshot/dirty/account/hit", nil)
	snapshotDirtyAccountMissMeter  = metrics.NewRegisteredMeter("state/snapshot/dirty/account/miss", nil)
	snapshotDirtyAccountInexMeter  = metrics.NewRegisteredMeter("state/snapshot/dirty/account/inex", nil)
	snapshotDirtyAccountReadMeter  = metrics.NewRegisteredMeter("state/snapshot/dirty/account/read", nil)
	snapshotDirtyAccountWriteMeter = metrics.NewRegisteredMeter("state/snapshot/dirty/account/write", nil)

	snapshotDirtyStorageHitMeter   = metrics.NewRegisteredMeter("state/snapshot/dirty/storage/hit", nil)
	snapshotDirtyStorageMissMeter  = metrics.NewRegisteredMeter("state/snapshot/dirty/storage/miss", nil)
	snapshotDirtyStorageInexMeter  = metrics.NewRegisteredMeter("state/snapshot/dirty/storage/inex", nil)
	snapshotDirtyStorageReadMeter  = metrics.NewRegisteredMeter("state/snapshot/dirty/storage/read", nil)
	snapshotDirtyStorageWriteMeter = metrics.NewRegisteredMeter("state/snapshot/dirty/storage/write", nil)

	snapshotDirtyAccountHitDepthHist = metrics.NewRegisteredHistogram("state/snapshot/dirty/account/hit/depth", nil, metrics.NewExpDecaySample(1028, 0.015))
	snapshotDirtyStorageHitDepthHist = metrics.NewRegisteredHistogram("state/snapshot/dirty/storage/hit/depth", nil, metrics.NewExpDecaySample(1028, 0.015))

	snapshotFlushAccountItemMeter = metrics.NewRegisteredMeter("state/snapshot/flush/account/item", nil)
	snapshotFlushAccountSizeMeter = metrics.NewRegisteredMeter("state/snapshot/flush/account/size", nil)
	snapshotFlushStorageItemMeter = metrics.NewRegisteredMeter("state/snapshot/flush/storage/item", nil)
	snapshotFlushStorageSizeMeter = metrics.NewRegisteredMeter("state/snapshot/flush/storage/size", nil)

	// TODO-Klaytn-Snapshot update snapshotBloomIndexTimer
	//snapshotBloomIndexTimer = metrics.NewRegisteredResettingTimer("state/snapshot/bloom/index", nil)
	snapshotBloomErrorGauge = metrics.NewRegisteredGaugeFloat64("state/snapshot/bloom/error", nil)

	snapshotBloomAccountTrueHitMeter  = metrics.NewRegisteredMeter("state/snapshot/bloom/account/truehit", nil)
	snapshotBloomAccountFalseHitMeter = metrics.NewRegisteredMeter("state/snapshot/bloom/account/falsehit", nil)
	snapshotBloomAccountMissMeter     = metrics.NewRegisteredMeter("state/snapshot/bloom/account/miss", nil)

	snapshotBloomStorageTrueHitMeter  = metrics.NewRegisteredMeter("state/snapshot/bloom/storage/truehit", nil)
	snapshotBloomStorageFalseHitMeter = metrics.NewRegisteredMeter("state/snapshot/bloom/storage/falsehit", nil)
	snapshotBloomStorageMissMeter     = metrics.NewRegisteredMeter("state/snapshot/bloom/storage/miss", nil)

	// ErrSnapshotStale is returned from data accessors if the underlying snapshot
	// layer had been invalidated due to the chain progressing forward far enough
	// to not maintain the layer's original state.
	ErrSnapshotStale = errors.New("snapshot stale")

	// ErrNotCoveredYet is returned from data accessors if the underlying snapshot
	// is being generated currently and the requested data item is not yet in the
	// range of accounts covered.
	ErrNotCoveredYet = errors.New("not covered yet")

	// ErrNotConstructed is returned if the callers want to iterate the snapshot
	// while the generation is not finished yet.
	ErrNotConstructed = errors.New("snapshot is not constructed")

	// errSnapshotCycle is returned if a snapshot is attempted to be inserted
	// that forms a cycle in the snapshot tree.
	errSnapshotCycle = errors.New("snapshot cycle")
)

// Snapshot represents the functionality supported by a snapshot storage layer.
type Snapshot interface {
	// Root returns the root hash for which this snapshot was made.
	Root() common.Hash

	// Account directly retrieves the account associated with a particular hash in
	// the snapshot slim data format.
	Account(hash common.Hash) (*Account, error)

	// AccountRLP directly retrieves the account RLP associated with a particular
	// hash in the snapshot slim data format.
	AccountRLP(hash common.Hash) ([]byte, error)

	// Storage directly retrieves the storage data associated with a particular hash,
	// within a particular account.
	Storage(accountHash, storageHash common.Hash) ([]byte, error)
}

// snapshot is the internal version of the snapshot data layer that supports some
// additional methods compared to the public API.
type snapshot interface {
	Snapshot

	// Parent returns the subsequent layer of a snapshot, or nil if the base was
	// reached.
	//
	// Note, the method is an internal helper to avoid type switching between the
	// disk and diff layers. There is no locking involved.
	Parent() snapshot

	// Update creates a new layer on top of the existing snapshot diff tree with
	// the specified data items.
	//
	// Note, the maps are retained by the method to avoid copying everything.
	Update(blockRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte) *diffLayer

	// TODO-Klaytn-Snapshot port journal.go first
	//// Journal commits an entire diff hierarchy to disk into a single journal entry.
	//// This is meant to be used during shutdown to persist the snapshot without
	//// flattening everything down (bad for reorgs).
	//Journal(buffer *bytes.Buffer) (common.Hash, error)

	// Stale return whether this layer has become stale (was flattened across) or
	// if it's still live.
	Stale() bool

	// AccountIterator creates an account iterator over an arbitrary layer.
	AccountIterator(seek common.Hash) AccountIterator

	// StorageIterator creates a storage iterator over an arbitrary layer.
	StorageIterator(account common.Hash, seek common.Hash) (StorageIterator, bool)
}

// diffToDisk merges a bottom-most diff into the persistent disk layer underneath
// it. The method will panic if called onto a non-bottom-most diff layer.
//
// The disk layer persistence should be operated in an atomic way. All updates should
// be discarded if the whole transition if not finished.
func diffToDisk(bottom *diffLayer) *diskLayer {
	var (
		base  = bottom.parent.(*diskLayer)
		batch = base.diskdb.NewSnapshotDBBatch()
		stats *generatorStats
	)
	// If the disk layer is running a snapshot generator, abort it
	if base.genAbort != nil {
		abort := make(chan *generatorStats)
		base.genAbort <- abort
		stats = <-abort
	}
	// Put the deletion in the batch writer, flush all updates in the final step.
	batch.DeleteSnapshotRoot()

	// Mark the original base as stale as we're going to create a new wrapper
	base.lock.Lock()
	if base.stale {
		panic("parent disk layer is stale") // we've committed into the same base from two children, boo
	}
	base.stale = true
	base.lock.Unlock()

	// Destroy all the destructed accounts from the database
	for hash := range bottom.destructSet {
		// Skip any account not covered yet by the snapshot
		if base.genMarker != nil && bytes.Compare(hash[:], base.genMarker) > 0 {
			continue
		}
		// Remove all storage slots
		batch.DeleteAccountSnapshot(hash)
		base.cache.Set(hash[:], nil)

		it := base.diskdb.NewSnapshotDBIterator(database.StorageSnapshotsKey(hash), nil)
		for it.Next() {
			if key := it.Key(); len(key) == 65 { // TODO(karalabe): Yuck, we should move this into the iterator
				batch.Delete(key)
				base.cache.Del(key[1:])
				snapshotFlushStorageItemMeter.Mark(1)

				// Ensure we don't delete too much data blindly (contract can be
				// huge). It's ok to flush, the root will go missing in case of a
				// crash and we'll detect and regenerate the snapshot.
				if batch.ValueSize() > database.IdealBatchSize {
					if err := batch.Write(); err != nil {
						logger.Crit("Failed to write storage deletions", "err", err)
					}
					batch.Reset()
				}
			}
		}
		it.Release()
	}
	// Push all updated accounts into the database
	for hash, data := range bottom.accountData {
		// Skip any account not covered yet by the snapshot
		if base.genMarker != nil && bytes.Compare(hash[:], base.genMarker) > 0 {
			continue
		}
		// Push the account to disk
		batch.WriteAccountSnapshot(hash, data)
		base.cache.Set(hash[:], data)
		snapshotCleanAccountWriteMeter.Mark(int64(len(data)))

		snapshotFlushAccountItemMeter.Mark(1)
		snapshotFlushAccountSizeMeter.Mark(int64(len(data)))

		// Ensure we don't write too much data blindly. It's ok to flush, the
		// root will go missing in case of a crash and we'll detect and regen
		// the snapshot.
		if batch.ValueSize() > database.IdealBatchSize {
			if err := batch.Write(); err != nil {
				logger.Crit("Failed to write storage deletions", "err", err)
			}
			batch.Reset()
		}
	}
	// Push all the storage slots into the database
	for accountHash, storage := range bottom.storageData {
		// Skip any account not covered yet by the snapshot
		if base.genMarker != nil && bytes.Compare(accountHash[:], base.genMarker) > 0 {
			continue
		}
		// Generation might be mid-account, track that case too
		midAccount := base.genMarker != nil && bytes.Equal(accountHash[:], base.genMarker[:common.HashLength])

		for storageHash, data := range storage {
			// Skip any slot not covered yet by the snapshot
			if midAccount && bytes.Compare(storageHash[:], base.genMarker[common.HashLength:]) > 0 {
				continue
			}
			if len(data) > 0 {
				batch.WriteStorageSnapshot(accountHash, storageHash, data)
				base.cache.Set(append(accountHash[:], storageHash[:]...), data)
				snapshotCleanStorageWriteMeter.Mark(int64(len(data)))
			} else {
				batch.DeleteStorageSnapshot(accountHash, storageHash)
				base.cache.Set(append(accountHash[:], storageHash[:]...), nil)
			}
			snapshotFlushStorageItemMeter.Mark(1)
			snapshotFlushStorageSizeMeter.Mark(int64(len(data)))
		}
	}
	// Update the snapshot block marker and write any remainder data
	batch.WriteSnapshotRoot(bottom.root)

	// Write out the generator progress marker and report
	journalProgress(batch, base.genMarker, stats)

	// Flush all the updates in the single db operation. Ensure the
	// disk layer transition is atomic.
	if err := batch.Write(); err != nil {
		logger.Crit("Failed to write leftover snapshot", "err", err)
	}
	logger.Debug("Journalled disk layer", "root", bottom.root, "complete", base.genMarker == nil)
	res := &diskLayer{
		root:       bottom.root,
		cache:      base.cache,
		diskdb:     base.diskdb,
		triedb:     base.triedb,
		genMarker:  base.genMarker,
		genPending: base.genPending,
	}
	// If snapshot generation hasn't finished yet, port over all the starts and
	// continue where the previous round left off.
	//
	// Note, the `base.genAbort` comparison is not used normally, it's checked
	// to allow the tests to play with the marker without triggering this path.
	if base.genMarker != nil && base.genAbort != nil {
		res.genMarker = base.genMarker
		res.genAbort = make(chan chan *generatorStats)
		go res.generate(stats)
	}
	return res
}
