package reward

import "sync"

const (
	maxStakingCache = 4
)

type stakingInfoCache struct {
	cells       map[uint64]*StakingInfo
	minBlockNum uint64
	lock        sync.RWMutex
}

func newStakingInfoCache() *stakingInfoCache {
	stakingCache := new(stakingInfoCache)
	stakingCache.cells = make(map[uint64]*StakingInfo)
	return stakingCache
}

func (sc *stakingInfoCache) get(blockNum uint64) *StakingInfo {
	sc.lock.RLock()
	defer sc.lock.RUnlock()

	if s, ok := sc.cells[blockNum]; ok {
		return s
	}
	return nil
}

func (sc *stakingInfoCache) add(stakingInfo *StakingInfo) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	// Assumption: stakingInfo should not be nil.

	if _, ok := sc.cells[stakingInfo.BlockNum]; ok {
		return
	}

	if len(sc.cells) < maxStakingCache {
		// empty room available
		sc.cells[stakingInfo.BlockNum] = stakingInfo
		if stakingInfo.BlockNum < sc.minBlockNum || len(sc.cells) == 1 {
			// new minBlockNum or newly inserted one is the first element
			sc.minBlockNum = stakingInfo.BlockNum
		}
		return
	}

	// evict one and insert new one
	delete(sc.cells, sc.minBlockNum)

	// update minBlockNum
	if stakingInfo.BlockNum < sc.minBlockNum {
		sc.minBlockNum = stakingInfo.BlockNum
	} else {
		min := stakingInfo.BlockNum
		for _, s := range sc.cells {
			if s.BlockNum < min {
				min = s.BlockNum
			}
		}
		sc.minBlockNum = min
	}
	sc.cells[stakingInfo.BlockNum] = stakingInfo
}
