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
// This file is derived from core/state/snapshot/iterator.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
)

// Iterator is an iterator to step over all the accounts or the specific
// storage in a snapshot which may or may not be composed of multiple layers.
type Iterator interface {
	// Next steps the iterator forward one element, returning false if exhausted,
	// or an error if iteration failed for some reason (e.g. root being iterated
	// becomes stale and garbage collected).
	Next() bool

	// Error returns any failure that occurred during iteration, which might have
	// caused a premature iteration exit (e.g. snapshot stack becoming stale).
	Error() error

	// Hash returns the hash of the account or storage slot the iterator is
	// currently at.
	Hash() common.Hash

	// Release releases associated resources. Release should always succeed and
	// can be called multiple times without causing error.
	Release()
}

// AccountIterator is an iterator to step over all the accounts in a snapshot,
// which may or may not be composed of multiple layers.
type AccountIterator interface {
	Iterator

	// Account returns the RLP encoded slim account the iterator is currently at.
	// An error will be returned if the iterator becomes invalid
	Account() []byte
}

// StorageIterator is an iterator to step over the specific storage in a snapshot,
// which may or may not be composed of multiple layers.
type StorageIterator interface {
	Iterator

	// Slot returns the storage slot the iterator is currently at. An error will
	// be returned if the iterator becomes invalid
	Slot() []byte
}

// diskAccountIterator is an account iterator that steps over the live accounts
// contained within a disk layer.
type diskAccountIterator struct {
	layer *diskLayer
	it    database.Iterator
}

// AccountIterator creates an account iterator over a disk layer.
func (dl *diskLayer) AccountIterator(seek common.Hash) AccountIterator {
	pos := common.TrimRightZeroes(seek[:])
	return &diskAccountIterator{
		layer: dl,
		it:    dl.diskdb.GetSnapshotDB().NewIterator(database.SnapshotAccountPrefix, pos),
	}
}

// Next steps the iterator forward one element, returning false if exhausted.
func (it *diskAccountIterator) Next() bool {
	// If the iterator was already exhausted, don't bother
	if it.it == nil {
		return false
	}
	// Try to advance the iterator and release it if we reached the end
	for {
		if !it.it.Next() {
			it.it.Release()
			it.it = nil
			return false
		}
		if len(it.it.Key()) == len(database.SnapshotAccountPrefix)+common.HashLength {
			break
		}
	}
	return true
}

// Error returns any failure that occurred during iteration, which might have
// caused a premature iteration exit (e.g. snapshot stack becoming stale).
//
// A diff layer is immutable after creation content wise and can always be fully
// iterated without error, so this method always returns nil.
func (it *diskAccountIterator) Error() error {
	if it.it == nil {
		return nil // Iterator is exhausted and released
	}
	return it.it.Error()
}

// Hash returns the hash of the account the iterator is currently at.
func (it *diskAccountIterator) Hash() common.Hash {
	return common.BytesToHash(it.it.Key()) // The prefix will be truncated
}

// Account returns the RLP encoded slim account the iterator is currently at.
func (it *diskAccountIterator) Account() []byte {
	return it.it.Value()
}

// Release releases the database snapshot held during iteration.
func (it *diskAccountIterator) Release() {
	// The iterator is auto-released on exhaustion, so make sure it's still alive
	if it.it != nil {
		it.it.Release()
		it.it = nil
	}
}

// diskStorageIterator is a storage iterator that steps over the live storage
// contained within a disk layer.
type diskStorageIterator struct {
	layer   *diskLayer
	account common.Hash
	it      database.Iterator
}

// StorageIterator creates a storage iterator over a disk layer.
// If the whole storage is destructed, then all entries in the disk
// layer are deleted already. So the "destructed" flag returned here
// is always false.
func (dl *diskLayer) StorageIterator(account common.Hash, seek common.Hash) (StorageIterator, bool) {
	pos := common.TrimRightZeroes(seek[:])
	return &diskStorageIterator{
		layer:   dl,
		account: account,
		it:      dl.diskdb.GetSnapshotDB().NewIterator(append(database.SnapshotStoragePrefix, account.Bytes()...), pos),
	}, false
}

// Next steps the iterator forward one element, returning false if exhausted.
func (it *diskStorageIterator) Next() bool {
	// If the iterator was already exhausted, don't bother
	if it.it == nil {
		return false
	}
	// Try to advance the iterator and release it if we reached the end
	for {
		if !it.it.Next() {
			it.it.Release()
			it.it = nil
			return false
		}
		if len(it.it.Key()) == len(database.SnapshotStoragePrefix)+common.HashLength+common.HashLength {
			break
		}
	}
	return true
}

// Error returns any failure that occurred during iteration, which might have
// caused a premature iteration exit (e.g. snapshot stack becoming stale).
//
// A diff layer is immutable after creation content wise and can always be fully
// iterated without error, so this method always returns nil.
func (it *diskStorageIterator) Error() error {
	if it.it == nil {
		return nil // Iterator is exhausted and released
	}
	return it.it.Error()
}

// Hash returns the hash of the storage slot the iterator is currently at.
func (it *diskStorageIterator) Hash() common.Hash {
	return common.BytesToHash(it.it.Key()) // The prefix will be truncated
}

// Slot returns the raw storage slot content the iterator is currently at.
func (it *diskStorageIterator) Slot() []byte {
	return it.it.Value()
}

// Release releases the database snapshot held during iteration.
func (it *diskStorageIterator) Release() {
	// The iterator is auto-released on exhaustion, so make sure it's still alive
	if it.it != nil {
		it.it.Release()
		it.it = nil
	}
}

// TODO-snapshot implement diffAccountIterator / diffStroageIterator
