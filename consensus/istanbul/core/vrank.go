// Copyright 2023 The klaytn Authors
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

package core

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"sort"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/rcrowley/go-metrics"
)

type Vrank struct {
	startTime             time.Time
	view                  istanbul.View
	committee             istanbul.Validators
	threshold             time.Duration
	firstCommit           int64
	quorumCommit          int64
	avgCommitWithinQuorum int64
	lastCommit            int64
	commitArrivalTimeMap  map[common.Address]time.Duration
}

var (
	// VRank metrics
	vrankFirstCommitArrivalTimeGauge           = metrics.NewRegisteredGauge("vrank/first_commit", nil)
	vrankQuorumCommitArrivalTimeGauge          = metrics.NewRegisteredGauge("vrank/quorum_commit", nil)
	vrankAvgCommitArrivalTimeWithinQuorumGauge = metrics.NewRegisteredGauge("vrank/avg_commit_within_quorum", nil)
	vrankLastCommitArrivalTimeGauge            = metrics.NewRegisteredGauge("vrank/last_commit", nil)

	vrankDefaultThreshold = "300ms" // the time to receive 2f+1 commits in an ideal network

	vrank *Vrank
)

const (
	vrankArrivedEarly = iota
	vrankArrivedLate
	vrankNotArrived
)

const (
	vrankNotArrivedPlaceholder = -1
)

func NewVrank(view istanbul.View, committee istanbul.Validators) *Vrank {
	threshold, _ := time.ParseDuration(vrankDefaultThreshold)
	return &Vrank{
		startTime:             time.Now(),
		view:                  view,
		committee:             committee,
		threshold:             threshold,
		firstCommit:           int64(0),
		quorumCommit:          int64(0),
		avgCommitWithinQuorum: int64(0),
		lastCommit:            int64(0),
		commitArrivalTimeMap:  make(map[common.Address]time.Duration),
	}
}

func (v *Vrank) TimeSinceStart() time.Duration {
	return time.Now().Sub(v.startTime)
}

func (v *Vrank) AddCommit(msg *istanbul.Subject, src istanbul.Validator) {
	if v.isTargetCommit(msg, src) {
		t := v.TimeSinceStart()
		v.commitArrivalTimeMap[src.Address()] = t
	}
}

func (v *Vrank) HandleCommitted(blockNum *big.Int) {
	if v.view.Sequence.Cmp(blockNum) != 0 {
		return
	}

	if len(v.commitArrivalTimeMap) != 0 {
		sum := int64(0)
		firstCommitTime := time.Duration(math.MaxInt64)
		quorumCommitTime := time.Duration(0)
		for _, arrivalTime := range v.commitArrivalTimeMap {
			sum += int64(arrivalTime)
			if firstCommitTime > arrivalTime {
				firstCommitTime = arrivalTime
			}
			if quorumCommitTime < arrivalTime {
				quorumCommitTime = arrivalTime
			}
		}
		avg := sum / int64(len(v.commitArrivalTimeMap))
		v.avgCommitWithinQuorum = avg
		v.firstCommit = int64(firstCommitTime)
		v.quorumCommit = int64(quorumCommitTime)

		if quorumCommitTime != time.Duration(0) && v.threshold > quorumCommitTime {
			v.threshold = quorumCommitTime
		}
	}
}

func (v *Vrank) Bitmap() string {
	serialized := serialize(v.committee, v.commitArrivalTimeMap)
	assessed := assessBatch(serialized, v.threshold)
	compressed := compress(assessed)
	return hex.EncodeToString(compressed)
}

func (v *Vrank) LateCommits() []time.Duration {
	serialized := serialize(v.committee, v.commitArrivalTimeMap)
	lateCommits := make([]time.Duration, 0)
	for _, t := range serialized {
		if assess(t, v.threshold) == vrankArrivedLate {
			lateCommits = append(lateCommits, t)
		}
	}
	return lateCommits
}

// Log logs accumulated data in a compressed form
func (v *Vrank) Log() {
	var (
		lastCommit  = time.Duration(0)
		lateCommits = v.LateCommits()
	)

	// lastCommit = max(lateCommits)
	for _, t := range lateCommits {
		if lastCommit < t {
			lastCommit = t
		}
	}
	v.lastCommit = int64(lastCommit)

	v.updateMetrics()

	logger.Info("VRank", "seq", v.view.Sequence.Int64(),
		"round", v.view.Round.Int64(),
		"bitmap", v.Bitmap(),
		"late", encodeDurationBatch(lateCommits),
	)
}

func (v *Vrank) updateMetrics() {
	if v.firstCommit != int64(0) {
		vrankFirstCommitArrivalTimeGauge.Update(v.firstCommit)
	}
	if v.quorumCommit != int64(0) {
		vrankQuorumCommitArrivalTimeGauge.Update(v.quorumCommit)
	}
	if v.avgCommitWithinQuorum != int64(0) {
		vrankAvgCommitArrivalTimeWithinQuorumGauge.Update(v.avgCommitWithinQuorum)
	}
	if v.lastCommit != int64(0) {
		vrankLastCommitArrivalTimeGauge.Update(v.lastCommit)
	}
}

func (v *Vrank) isTargetCommit(msg *istanbul.Subject, src istanbul.Validator) bool {
	if msg.View == nil || msg.View.Sequence == nil || msg.View.Round == nil {
		return false
	}
	if msg.View.Cmp(&v.view) != 0 {
		return false
	}
	_, ok := v.commitArrivalTimeMap[src.Address()]
	if ok {
		return false
	}
	return true
}

// assess determines if given time is early, late, or not arrived
func assess(t, threshold time.Duration) uint8 {
	if t == vrankNotArrivedPlaceholder {
		return vrankNotArrived
	}

	if t > threshold {
		return vrankArrivedLate
	} else {
		return vrankArrivedEarly
	}
}

func assessBatch(ts []time.Duration, threshold time.Duration) []uint8 {
	ret := make([]uint8, len(ts))
	for i, t := range ts {
		ret[i] = assess(t, threshold)
	}
	return ret
}

// serialize serializes arrivalTime hashmap into array.
// If committee is sorted, we can simply figure out the validator position in the output array
// by sorting the output of `klay.getCommittee()`
func serialize(committee istanbul.Validators, arrivalTimeMap map[common.Address]time.Duration) []time.Duration {
	sortedCommittee := make(istanbul.Validators, len(committee))
	copy(sortedCommittee[:], committee[:])
	sort.Sort(sortedCommittee)

	serialized := make([]time.Duration, len(sortedCommittee))
	for i, v := range sortedCommittee {
		val, ok := arrivalTimeMap[v.Address()]
		if ok {
			serialized[i] = val
		} else {
			serialized[i] = vrankNotArrivedPlaceholder
		}

	}
	return serialized
}

// compress compresses data into 2-bit bitmap
// e.g., [1, 0, 2] => [0b01_00_10_00]
func compress(arr []uint8) []byte {
	zip := func(a, b, c, d uint8) byte {
		a &= 0b11
		b &= 0b11
		c &= 0b11
		d &= 0b11
		return byte(a<<6 | b<<4 | c<<2 | d<<0)
	}

	// pad zero to make len(arr)%4 == 0
	for len(arr)%4 != 0 {
		arr = append(arr, 0)
	}

	ret := make([]byte, len(arr)/4)

	for i := 0; i < len(arr)/4; i++ {
		chunk := arr[4*i : 4*(i+1)]
		ret[i] = zip(chunk[0], chunk[1], chunk[2], chunk[3])
	}
	return ret
}

// encodeDuration encodes given duration into string
// The returned string is at most 4 bytes
func encodeDuration(d time.Duration) string {
	if d > 10*time.Second {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d > time.Second {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else {
		return fmt.Sprintf("%d", d.Milliseconds())
	}
}

func encodeDurationBatch(ds []time.Duration) []string {
	ret := make([]string, len(ds))
	for i, d := range ds {
		ret[i] = encodeDuration(d)
	}
	return ret
}
