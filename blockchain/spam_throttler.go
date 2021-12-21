// Copyright 2021 The klaytn Authors
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

package blockchain

import (
	"errors"
	"sync"

	"github.com/rcrowley/go-metrics"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// TODO-Klaytn: move these variables into txpool when blockchain struct contains a txpool interface
// spamThrottler need to be accessed by both of TxPool and BlockChain.
var (
	DisableSpamThrottlerAtRuntime            = false
	spamThrottler                 *throttler = nil
	spamThrottlerMu                          = new(sync.RWMutex)
)

var (
	thresholdGauge     = metrics.NewRegisteredGauge("txpool/throttler/threshold", nil)
	candidateSizeGauge = metrics.NewRegisteredGauge("txpool/throttler/candidate/size", nil)
	throttledSizeGauge = metrics.NewRegisteredGauge("txpool/throttler/throttled/size", nil)
	allowedSizeGauge   = metrics.NewRegisteredGauge("txpool/throttler/allowed/size", nil)
)

type throttler struct {
	config *throttlerConfig

	candidates map[*common.Address]int  // throttle candidates with spam weight. Not for concurrent use
	throttled  map[*common.Address]int  // throttled addresses. It requires mu.lock for concurrent use
	allowed    map[*common.Address]bool // white listed addresses. It requires mu.lock for concurrent use
	mu         *sync.RWMutex            // mutex for throttled and allowed

	threshold  int
	throttleCh chan *types.Transaction
	quitCh     chan struct{}
}

type throttlerConfig struct {
	activateTxPoolSize uint `json:"activate_tx_pool_size"`
	targetFailRatio    uint `json:"target_fail_ratio"`
	throttleTPS        uint `json:"throttle_tps"`
	weightMapSize      uint `json:"weight_map_size"`

	increaseWeight   int `json:"increase_weight"`
	decreaseWeight   int `json:"decrease_weight"`
	initialThreshold int `json:"initial_threshold"` // initialThreshold <= threshold <= throttledWeight
	throttledWeight  int `json:"throttled_weight"`
}

var DefaultSpamThrottlerConfig = &throttlerConfig{
	activateTxPoolSize: 1000,
	targetFailRatio:    10,
	throttleTPS:        10,    // len(throttleCh) = throttleTPS * 3 = 32KB * 10 * 3 = 960KB
	weightMapSize:      10000, // (20 + 4)B * 10000 = 240KB

	increaseWeight:   5,
	decreaseWeight:   1,
	initialThreshold: 100,
	throttledWeight:  300,
}

func GetSpamThrottler() *throttler {
	spamThrottlerMu.RLock()
	t := spamThrottler
	spamThrottlerMu.RUnlock()
	return t
}

func validateConfig(conf *throttlerConfig) error {
	if conf == nil {
		return errors.New("nil throttlerConfig")
	}

	if conf.targetFailRatio > 100 {
		return errors.New("invalid throttlerConfig. 0 <= targetFailRatio <= 100")
	}

	if conf.initialThreshold > conf.increaseWeight && conf.initialThreshold < conf.throttledWeight {
		return errors.New("invalid throttlerConfig. increaseWeight < initialThreshold < throttledWeight")
	}

	return nil
}

// adjustThreshold adjusts the spam weight threshold of throttler in an adaptive way.
func (t *throttler) adjustThreshold(ratio uint) {
	var newThreshold int
	// Decrease threshold if a fail ratio is bigger than target value to put more addresses in throttled map
	if ratio > t.config.targetFailRatio {
		newThreshold = t.threshold - t.config.increaseWeight

		// Set minimum threshold
		if newThreshold < t.config.increaseWeight {
			newThreshold = t.config.increaseWeight
		}

		// Increase threshold if a fail ratio is smaller than target ratio until it exceeds initialThreshold
	} else {
		newThreshold = t.threshold + t.config.increaseWeight

		if newThreshold > t.config.initialThreshold {
			newThreshold = t.config.initialThreshold
		}
	}

	t.threshold = newThreshold

	// Update metrics
	thresholdGauge.Update(int64(newThreshold))
}

// newAllowed generates a new allowed list of throttler.
func (t *throttler) newAllowed(allowed []*common.Address) {
	t.mu.Lock()

	a := make(map[*common.Address]bool, len(allowed))
	for _, addr := range allowed {
		a[addr] = true
	}
	t.allowed = a
	allowedSize := len(allowed)

	t.mu.Unlock()

	// Update metrics
	allowedSizeGauge.Update(int64(allowedSize))
}

// updateThrottled removes outdated addresses from the throttle list and adds new addresses to the list.
func (t *throttler) updateThrottled(newThrottled []*common.Address) {
	var removeThrottled []*common.Address
	t.mu.Lock()

	// Decrease spam weight for all throttled addresses.
	for addr, remained := range t.throttled {
		t.throttled[addr] = remained - t.config.decreaseWeight
		if t.throttled[addr] < 0 {
			removeThrottled = append(removeThrottled, addr)
		}
	}

	// Remove throttled addresses from throttled map.
	for _, addr := range removeThrottled {
		delete(t.throttled, addr)
	}

	for _, addr := range newThrottled {
		t.throttled[addr] = t.config.throttledWeight
	}

	size := len(t.throttled)
	t.mu.Unlock()

	// Update metrics
	throttledSizeGauge.Update(int64(size))
}

// updateThrottlerState updates the throttle list by calculating spam weight of candidates.
func (t *throttler) updateThrottlerState(txs types.Transactions, receipts types.Receipts) {
	var removeCandidate []*common.Address
	var newThrottled []*common.Address

	numFailed := 0
	failRatio := uint(0)
	mapSize := uint(len(t.candidates))

	// Increase spam weight of throttle candidates who generate failed txs.
	for i, receipt := range receipts {
		if receipt.Status != types.ReceiptStatusSuccessful {
			numFailed++

			toAddr := txs[i].To()
			if toAddr == nil {
				continue
			}

			weight := t.candidates[toAddr]
			if weight == 0 {
				if mapSize >= t.config.weightMapSize {
					continue
				}
				mapSize++
			}
			t.candidates[toAddr] = weight + t.config.increaseWeight
		}
	}

	// Decrease spam weight for all candidates and update throttle lists in throttled.
	for addr, weight := range t.candidates {
		newWeight := weight - t.config.decreaseWeight

		switch {
		case newWeight < 0:
			removeCandidate = append(removeCandidate, addr)

		case newWeight > t.threshold:
			removeCandidate = append(removeCandidate, addr)
			newThrottled = append(newThrottled, addr)

		default:
			t.candidates[addr] = newWeight
		}
	}

	// Remove throttle candidates from candidates map.
	for _, addr := range removeCandidate {
		delete(t.candidates, addr)
	}

	if len(receipts) != 0 {
		failRatio = uint(100 * numFailed / len(receipts))
	}

	// Update throttled and threshold
	t.updateThrottled(newThrottled)
	t.adjustThreshold(failRatio)

	// Update metrics
	candidateSizeGauge.Update(int64(len(t.candidates)))
}

// classifyTxs classifies given txs into allowTxs and throttleTxs.
// If to-address of tx is listed in the throttle list, it is classified as throttleTx.
func (t *throttler) classifyTxs(txs types.Transactions) (types.Transactions, types.Transactions) {
	allowTxs := txs[:0]
	throttleTxs := txs[:0]

	t.mu.RLock()
	for _, tx := range txs {
		if tx.To() != nil && t.throttled[tx.To()] > 0 && t.allowed[tx.To()] == false {
			throttleTxs = append(throttleTxs, tx)
		} else {
			allowTxs = append(allowTxs, tx)
		}
	}
	t.mu.RUnlock()

	return allowTxs, throttleTxs
}
