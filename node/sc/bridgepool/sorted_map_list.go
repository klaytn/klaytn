// Modifications Copyright 2019 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from core/tx_list.go (2018/06/04).
// Modified and improved for the klaytn development.

package bridgepool

import (
	"container/heap"
	"errors"
	"fmt"
	"sort"
	"sync"
)

// nonceHeap is a heap.Interface implementation over 64bit unsigned integers for
// retrieving sorted transactions from the possibly gapped future queue.
type nonceHeap []uint64

func (h nonceHeap) Len() int           { return len(h) }
func (h nonceHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h nonceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nonceHeap) Push(x interface{}) {
	*h = append(*h, x.(uint64))
}

func (h *nonceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type itemWithNonce interface {
	Nonce() uint64
}

type items []itemWithNonce

func (s items) Len() int { return len(s) }
func (s items) Less(i, j int) bool {
	return s[i].Nonce() < s[j].Nonce()
}
func (s items) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ItemSortedMap is a nonce->item map with a heap based index to allow
// iterating over the contents in a nonce-incrementing way.
const (
	UnlimitedItemSortedMap = -1
)

var (
	ErrSizeLimit = errors.New("sorted map size limit")
)

type ItemSortedMap struct {
	mu        *sync.Mutex
	items     map[uint64]itemWithNonce // Hash map storing the item data
	index     *nonceHeap               // Heap of nonces of all the stored items (non-strict mode)
	cache     items                    // Cache of the items already sorted
	sizeLimit int                      // This is sizeLimit of the sorted map.
}

// NewItemSortedMap creates a new nonce-sorted item map.
func NewItemSortedMap(size int) *ItemSortedMap {
	return &ItemSortedMap{
		mu:        new(sync.Mutex),
		items:     make(map[uint64]itemWithNonce),
		index:     new(nonceHeap),
		sizeLimit: size,
	}
}

// Get retrieves the current items associated with the given nonce.
func (m *ItemSortedMap) Get(nonce uint64) itemWithNonce {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.items[nonce]
}

// Exist returns if the nonce exist.
func (m *ItemSortedMap) Exist(nonce uint64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exist := m.items[nonce]

	return exist
}

// Put inserts a new item into the map, also updating the map's nonce
// index. If a item already exists with the same nonce, it's overwritten.
func (m *ItemSortedMap) Put(event itemWithNonce) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	nonce := event.Nonce()
	if m.sizeLimit != UnlimitedItemSortedMap && len(m.items) >= m.sizeLimit && m.items[nonce] == nil {
		return fmt.Errorf("failed to put %v nonce : %w : %v", event.Nonce(), ErrSizeLimit, m.sizeLimit)
	}

	if m.items[nonce] == nil {
		heap.Push(m.index, nonce)
	}
	m.items[nonce], m.cache = event, nil

	return nil
}

// Pop removes given count number of minimum nonce items from the map.
// Every removed items is returned for any post-removal maintenance.
func (m *ItemSortedMap) Pop(count int) items {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Otherwise start accumulating incremental events
	var ready items
	for m.index.Len() > 0 && len(ready) < count {
		nonce := (*m.index)[0]
		ready = append(ready, m.items[nonce])
		delete(m.items, nonce)
		heap.Pop(m.index)
	}

	return ready
}

// Ready retrieves a sequentially increasing list of events starting at the
// provided nonce that is ready for processing. The returned events will be
// removed from the list.
//
// Note, all events with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (m *ItemSortedMap) Ready(start uint64) items {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Short circuit if no events are available
	if m.index.Len() == 0 || (*m.index)[0] > start {
		return nil
	}
	// Otherwise start accumulating incremental events
	var ready items
	for next := (*m.index)[0]; m.index.Len() > 0 && (*m.index)[0] == next; next++ {
		ready = append(ready, m.items[next])
		delete(m.items, next)
		heap.Pop(m.index)
	}

	return ready
}

// Filter iterates over the list of items and removes all of them for which
// the specified function evaluates to true.
func (m *ItemSortedMap) Filter(filter func(interface{}) bool) items {
	m.mu.Lock()
	defer m.mu.Unlock()

	var removed items

	// Collect all the items to filter out
	for nonce, tx := range m.items {
		if filter(tx) {
			removed = append(removed, tx)
			delete(m.items, nonce)
		}
	}
	// If items were removed, the heap and cache are ruined
	if len(removed) > 0 {
		*m.index = make([]uint64, 0, len(m.items))
		for nonce := range m.items {
			*m.index = append(*m.index, nonce)
		}
		heap.Init(m.index)

		m.cache = nil
	}
	return removed
}

// Remove deletes a item from the maintained map, returning whether the
// item was found.
func (m *ItemSortedMap) Remove(nonce uint64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Short circuit if no item is present
	_, ok := m.items[nonce]
	if !ok {
		return false
	}
	// Otherwise delete the item and fix the heap index
	for i := 0; i < m.index.Len(); i++ {
		if (*m.index)[i] == nonce {
			heap.Remove(m.index, i)
			break
		}
	}
	delete(m.items, nonce)
	m.cache = nil

	return true
}

// Forward removes all items from the map with a nonce lower than the
// provided threshold. Every removed transaction is returned for any post-removal
// maintenance.
func (m *ItemSortedMap) Forward(threshold uint64) items {
	m.mu.Lock()
	defer m.mu.Unlock()
	var removed items

	// Pop off heap items until the threshold is reached
	for m.index.Len() > 0 && (*m.index)[0] < threshold {
		nonce := heap.Pop(m.index).(uint64)
		removed = append(removed, m.items[nonce])
		delete(m.items, nonce)
	}
	// If we had a cached order, shift the front
	if m.cache != nil {
		m.cache = m.cache[len(removed):]
	}
	return removed
}

// Len returns the length of the map.
func (m *ItemSortedMap) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.items)
}

// Flatten returns a nonce-sorted slice of items based on the loosely
// sorted internal representation. The result of the sorting is cached in case
// it's requested again before any modifications are made to the contents.
func (m *ItemSortedMap) Flatten() items {
	return m.FlattenByCount(0)
}

// FlattenByCount returns requested number of nonce-sorted slice of cached
// items. The result of the sorting is cached like as Flatten method.
func (m *ItemSortedMap) FlattenByCount(count int) items {
	m.mu.Lock()
	defer m.mu.Unlock()
	// If the sorting was not cached yet, create and cache it
	if m.cache == nil {
		m.cache = make(items, 0, len(m.items))
		for _, tx := range m.items {
			m.cache = append(m.cache, tx)
		}
		sort.Sort(m.cache)
	}
	txLen := len(m.cache)
	if count != 0 && txLen > count {
		txLen = count
	}
	return m.cache[:txLen]
}
