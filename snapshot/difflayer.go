package snapshot

import (
	"sync"

	"github.com/klaytn/klaytn/common"
	"github.com/steakknife/bloomfilter"
)

// diffLayer represents a collection of modifications made to a state snapshot
// after running a block on top. It contains one sorted list for the account trie
// and one-one list for each storage tries.
//
// The goal of a diff layer is to act as a journal, tracking recent modifications
// made to the state, that have not yet graduated into a semi-immutable state.
type diffLayer struct {
	origin *diskLayer // Base disk layer to directly use on bloom misses
	parent snapshot   // Parent snapshot modified by this one, never nil
	memory uint64     // Approximate guess as to how much memory we use

	root  common.Hash // Root hash to which this snapshot diff belongs to
	stale uint32      // Signals that the layer became stale (state progressed)

	// destructSet is a very special helper marker. If an account is marked as
	// deleted, then it's recorded in this set. However it's allowed that an account
	// is included here but still available in other sets(e.g. storageData). The
	// reason is the diff layer includes all the changes in a *block*. It can
	// happen that in the tx_1, account A is self-destructed while in the tx_2
	// it's recreated. But we still need this marker to indicate the "old" A is
	// deleted, all data in other set belongs to the "new" A.
	destructSet map[common.Hash]struct{}               // Keyed markers for deleted (and potentially) recreated accounts
	accountList []common.Hash                          // List of account for iteration. If it exists, it's sorted, otherwise it's nil
	accountData map[common.Hash][]byte                 // Keyed accounts for direct retrieval (nil means deleted)
	storageList map[common.Hash][]common.Hash          // List of storage slots for iterated retrievals, one per account. Any existing lists are sorted if non-nil
	storageData map[common.Hash]map[common.Hash][]byte // Keyed storage slots for direct retrieval. one per account (nil means deleted)

	diffed *bloomfilter.Filter // Bloom filter tracking all the diffed items up to the disk layer

	lock sync.RWMutex
}
