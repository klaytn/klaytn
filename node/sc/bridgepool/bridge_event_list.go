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
	"sync"
)

type itemWithNonce interface {
	Nonce() uint64
}

// EventSortedMap is a nonce->item map with a heap based index to allow
// iterating over the contents in a nonce-incrementing way.
type EventSortedMap struct {
	mu    *sync.Mutex
	items map[uint64]itemWithNonce // Hash map storing the item data
	index *nonceHeap               // Heap of nonces of all the stored items (non-strict mode)
}

// NewEventSortedMap creates a new nonce-sorted item map.
func NewEventSortedMap() *EventSortedMap {
	return &EventSortedMap{
		mu:    new(sync.Mutex),
		items: make(map[uint64]itemWithNonce),
		index: new(nonceHeap),
	}
}

// Get retrieves the current items associated with the given nonce.
func (m *EventSortedMap) Get(nonce uint64) itemWithNonce {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.items[nonce]
}

// Put inserts a new item into the map, also updating the map's nonce
// index. If a item already exists with the same nonce, it's overwritten.
func (m *EventSortedMap) Put(event itemWithNonce) {
	m.mu.Lock()
	defer m.mu.Unlock()

	nonce := event.Nonce()
	if m.items[nonce] == nil {
		heap.Push(m.index, nonce)
	}
	m.items[nonce] = event
}

// Ready retrieves a sequentially increasing list of events starting at the
// provided nonce that is ready for processing. The returned events will be
// removed from the list.
//
// Note, all events with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (m *EventSortedMap) Ready(start uint64) []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Short circuit if no events are available
	if m.index.Len() == 0 || (*m.index)[0] > start {
		return nil
	}
	// Otherwise start accumulating incremental events
	var ready []interface{}
	for next := (*m.index)[0]; m.index.Len() > 0 && (*m.index)[0] == next; next++ {
		ready = append(ready, m.items[next])
		delete(m.items, next)
		heap.Pop(m.index)
	}

	return ready
}

func (m *EventSortedMap) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.items)
}
