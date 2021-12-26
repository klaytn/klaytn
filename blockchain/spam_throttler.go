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
	"time"

	"github.com/rcrowley/go-metrics"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// TODO-Klaytn: move these variables into TxPool when BlockChain struct contains a TxPool interface
// spamThrottler need to be accessed by both of TxPool and BlockChain.
var (
	spamThrottler   *throttler = nil
	spamThrottlerMu            = new(sync.RWMutex)
)

var (
	thresholdGauge           = metrics.NewRegisteredGauge("txpool/throttler/threshold", nil)
	candidateSizeGauge       = metrics.NewRegisteredGauge("txpool/throttler/candidate/size", nil)
	throttledSizeGauge       = metrics.NewRegisteredGauge("txpool/throttler/throttled/size", nil)
	allowedSizeGauge         = metrics.NewRegisteredGauge("txpool/throttler/allowed/size", nil)
	throttlerUpdateTimeGauge = metrics.NewRegisteredGauge("txpool/throttler/update/time", nil)
)

type throttler struct {
	config *ThrottlerConfig

	candidates map[common.Address]int  // throttle candidates with spam weight. Not for concurrent use
	throttled  map[common.Address]int  // throttled addresses. It requires mu.lock for concurrent use
	allowed    map[common.Address]bool // white listed addresses. It requires mu.lock for concurrent use
	mu         *sync.RWMutex           // mutex for throttled and allowed

	threshold  int
	throttleCh chan *types.Transaction
	quitCh     chan struct{}
}

type ThrottlerConfig struct {
	ActivateTxPoolSize uint `json:"activate_tx_pool_size"`
	TargetFailRatio    uint `json:"target_fail_ratio"`
	ThrottleTPS        uint `json:"throttle_tps"`
	WeightMapSize      uint `json:"weight_map_size"`

	IncreaseWeight   int `json:"increase_weight"`
	DecreaseWeight   int `json:"decrease_weight"`
	InitialThreshold int `json:"initial_threshold"` // InitialThreshold <= threshold <= ThrottledWeight
	ThrottledWeight  int `json:"throttled_weight"`
}

var DefaultSpamThrottlerConfig = &ThrottlerConfig{
	ActivateTxPoolSize: 1000,
	TargetFailRatio:    10,
	ThrottleTPS:        10,    // len(throttleCh) = ThrottleTPS * 3 = 32KB * 10 * 3 = 960KB
	WeightMapSize:      10000, // (20 + 4)B * 10000 = 240KB

	IncreaseWeight:   5,
	DecreaseWeight:   1,
	InitialThreshold: 100,
	ThrottledWeight:  300,
}

func GetSpamThrottler() *throttler {
	spamThrottlerMu.RLock()
	t := spamThrottler
	spamThrottlerMu.RUnlock()
	return t
}

func validateConfig(conf *ThrottlerConfig) error {
	if conf == nil {
		return errors.New("nil ThrottlerConfig")
	}

	if conf.TargetFailRatio > 100 {
		return errors.New("invalid ThrottlerConfig. 0 <= TargetFailRatio <= 100")
	}

	if conf.InitialThreshold <= conf.IncreaseWeight || conf.InitialThreshold >= conf.ThrottledWeight {
		return errors.New("invalid ThrottlerConfig. IncreaseWeight < InitialThreshold < ThrottledWeight")
	}

	return nil
}

// adjustThreshold adjusts the spam weight threshold of throttler in an adaptive way.
func (t *throttler) adjustThreshold(ratio uint) {
	var newThreshold int
	// Decrease threshold if a fail ratio is bigger than target value to put more addresses in throttled map
	if ratio > t.config.TargetFailRatio {
		newThreshold = t.threshold - t.config.IncreaseWeight

		// Set minimum threshold
		if newThreshold < t.config.IncreaseWeight {
			newThreshold = t.config.IncreaseWeight
		}

		// Increase threshold if a fail ratio is smaller than target ratio until it exceeds InitialThreshold
	} else {
		newThreshold = t.threshold + t.config.IncreaseWeight

		if newThreshold > t.config.InitialThreshold {
			newThreshold = t.config.InitialThreshold
		}
	}

	t.threshold = newThreshold

	// Update metrics
	thresholdGauge.Update(int64(newThreshold))
}

// newAllowed generates a new allowed list of throttler.
func (t *throttler) newAllowed(allowed []common.Address) {
	t.mu.Lock()

	a := make(map[common.Address]bool, len(allowed))
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
func (t *throttler) updateThrottled(newThrottled []common.Address) {
	var removeThrottled []common.Address
	t.mu.Lock()

	// Decrease spam weight for all throttled addresses.
	for addr, remained := range t.throttled {
		t.throttled[addr] = remained - t.config.DecreaseWeight
		if t.throttled[addr] < 0 {
			removeThrottled = append(removeThrottled, addr)
		}
	}

	// Remove throttled addresses from throttled map.
	for _, addr := range removeThrottled {
		delete(t.throttled, addr)
	}

	for _, addr := range newThrottled {
		t.throttled[addr] = t.config.ThrottledWeight
	}

	size := len(t.throttled)
	t.mu.Unlock()

	// Update metrics
	throttledSizeGauge.Update(int64(size))
}

// updateThrottlerState updates the throttle list by calculating spam weight of candidates.
func (t *throttler) updateThrottlerState(txs types.Transactions, receipts types.Receipts) {
	var removeCandidate []common.Address
	var newThrottled []common.Address

	startTime := time.Now()
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

			weight := t.candidates[*toAddr]
			if weight == 0 {
				if mapSize >= t.config.WeightMapSize {
					continue
				}
				mapSize++
			}

			t.candidates[*toAddr] = weight + t.config.IncreaseWeight
		}
	}

	// Decrease spam weight for all candidates and update throttle lists in throttled.
	for addr, weight := range t.candidates {
		newWeight := weight - t.config.DecreaseWeight

		switch {
		case newWeight <= 0:
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
	throttlerUpdateTimeGauge.Update(int64(time.Since(startTime)))
}

// classifyTxs classifies given txs into allowTxs and throttleTxs.
// If to-address of tx is listed in the throttle list, it is classified as throttleTx.
func (t *throttler) classifyTxs(txs types.Transactions) (types.Transactions, types.Transactions) {
	allowTxs := txs[:0]
	throttleTxs := txs[:0]

	t.mu.RLock()
	for _, tx := range txs {
		if tx.To() != nil && t.throttled[*tx.To()] > 0 && t.allowed[*tx.To()] == false {
			throttleTxs = append(throttleTxs, tx)
		} else {
			allowTxs = append(allowTxs, tx)
		}
	}
	t.mu.RUnlock()

	return allowTxs, throttleTxs
}

// SetAllowed resets the allowed list of throttler. The previous list will be abandoned.
func (t *throttler) SetAllowed(addrs []common.Address) {
	t.mu.Lock()
	t.allowed = make(map[common.Address]bool)
	for _, addr := range addrs {
		t.allowed[addr] = true
	}
	t.mu.Unlock()
}

func (t *throttler) GetAllowed() []common.Address {
	var addrs []common.Address

	t.mu.RLock()
	for addr := range t.allowed {
		addrs = append(addrs, addr)
	}
	t.mu.RUnlock()

	return addrs
}

func (t *throttler) GetThrottled() []common.Address {
	var addrs []common.Address

	t.mu.RLock()
	for addr := range t.throttled {
		addrs = append(addrs, addr)
	}
	t.mu.RUnlock()

	return addrs
}

func (t *throttler) GetCandidates() map[common.Address]int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.candidates
}

func (t *throttler) GetConfig() *ThrottlerConfig {
	return t.config
}
