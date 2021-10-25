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
// This file is derived from core/state/snapshot/disklayer.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"sync"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
)

// diskLayer is a low level persistent snapshot built on top of a key-value store.
type diskLayer struct {
	diskdb database.Database // Key-value store containing the base snapshot
	triedb *statedb.Database // Trie node cache for reconstruction purposes
	cache  *fastcache.Cache  // Cache to avoid hitting the disk for direct access

	root  common.Hash // Root hash of the base snapshot
	stale bool        // Signals that the layer became stale (state progressed)

	genMarker  []byte        // Marker for the state that's indexed during initial layer generation
	genPending chan struct{} // Notification channel when generation is done (test synchronicity)
	// TODO-snapshot uncomment after porting generate.go
	//genAbort   chan chan *generatorStats // Notification channel to abort generating the snapshot in this layer

	lock sync.RWMutex
}
