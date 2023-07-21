package core

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/rcrowley/go-metrics"
)

type vrank struct {
	startTime            time.Time
	view                 istanbul.View
	committee            istanbul.Validators
	threshold            time.Duration
	commitArrivalTimeMap map[common.Address]time.Duration
}

var (
	// VRank metrics
	vrankFirstCommitArrivalTimeGauge           = metrics.NewRegisteredGauge("vrank/first_commit", nil)
	vrankQuorumCommitArrivalTimeGauge          = metrics.NewRegisteredGauge("vrank/quorum_commit", nil)
	vrankAvgCommitArrivalTimeWithinQuorumGauge = metrics.NewRegisteredGauge("vrank/avg_commit_within_quorum", nil)
	vrankLastCommitArrivalTimeGauge            = metrics.NewRegisteredGauge("vrank/last_commit", nil)

	vrankDefaultThreshold = "300ms"

	Vrank *vrank
)

const (
	vrankArrivedEarly = iota
	vrankArrivedLate
	vrankNotArrived
)

const (
	vrankNotArrivedPlaceholder = -1
)

func NewVrank(view istanbul.View, committee istanbul.Validators) *vrank {
	threshold, _ := time.ParseDuration(vrankDefaultThreshold)
	return &vrank{
		startTime:            time.Now(),
		view:                 view,
		committee:            committee,
		threshold:            threshold,
		commitArrivalTimeMap: make(map[common.Address]time.Duration),
	}
}

func (v *vrank) TimeSinceStart() time.Duration {
	return time.Now().Sub(v.startTime)
}

func (v *vrank) AddCommit(msg *istanbul.Subject, src istanbul.Validator) {
	if v.isTargetCommit(msg, src) {
		v.commitArrivalTimeMap[src.Address()] = v.TimeSinceStart()
	}
}

func (v *vrank) HandleCommitted(blockNum *big.Int) {
	if v.view.Sequence.Cmp(blockNum) != 0 {
		return
	}

	committedTime := v.TimeSinceStart()
	if v.threshold > committedTime {
		v.threshold = committedTime
	}

	vrankQuorumCommitArrivalTimeGauge.Update(int64(committedTime))
	sum := int64(0)
	for _, v := range v.commitArrivalTimeMap {
		sum += int64(v)
	}
	avg := sum / int64(len(v.commitArrivalTimeMap))
	vrankAvgCommitArrivalTimeWithinQuorumGauge.Update(avg)
}

// Stop logs accumulated data in a compressed form
func (v *vrank) Stop() {
	var (
		serialized = serialize(v.committee, v.commitArrivalTimeMap)

		assessBitmap = hex.EncodeToString(compress(assessAll(serialized, v.threshold)))

		lateCommits        = make([]time.Duration, 0)
		encodedLateCommits = make([]string, 0)
		lastCommit         = time.Duration(0)
	)

	for _, t := range serialized {
		if assess(t, v.threshold) == vrankArrivedLate {
			lateCommits = append(lateCommits, t)
		}
	}

	for _, t := range lateCommits {
		if lastCommit < t {
			lastCommit = t
		}
	}
	if lastCommit != time.Duration(0) {
		vrankLastCommitArrivalTimeGauge.Update(int64(lastCommit))
	}

	// encode late commits
	// if t >
	for _, t := range lateCommits {
		var s string
		if t > time.Second {
			s = fmt.Sprintf("%.1fs", t.Seconds())
		} else {
			s = fmt.Sprintf("%d", t.Milliseconds())
		}
		encodedLateCommits = append(encodedLateCommits, s)
	}

	logger.Info("VRank",
		"bitmap", assessBitmap,
		"late", encodedLateCommits,
	)
}

func (v *vrank) isTargetCommit(msg *istanbul.Subject, src istanbul.Validator) bool {
	if msg.View.Cmp(&v.view) != 0 {
		return false
	}
	_, ok := v.commitArrivalTimeMap[src.Address()]
	if ok {
		return false
	}
	return true
}

func assess(t, threshold time.Duration) int {
	if t == vrankNotArrivedPlaceholder {
		return vrankNotArrived
	}

	if t > threshold {
		return vrankArrivedLate
	} else {
		return vrankArrivedEarly
	}
}

func assessAll(t []time.Duration, threshold time.Duration) []int {
	ret := make([]int, 0)
	for _, v := range t {
		ret = append(ret, assess(v, threshold))
	}
	return ret
}

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
func compress(arr []int) []byte {
	zip := func(a, b, c, d int) byte {
		a &= 0b11
		b &= 0b11
		c &= 0b11
		d &= 0b11
		return byte(a<<6 | b<<4 | c<<2 | d<<0)
	}

	// pad zero to make len(arr)%4 == 0
	switch len(arr) % 4 {
	case 1:
		arr = append(arr, []int{0, 0, 0}...)
	case 2:
		arr = append(arr, []int{0, 0}...)
	case 3:
		arr = append(arr, []int{0}...)
	}

	ret := make([]byte, 0)

	for i := 0; i < len(arr)/4; i++ {
		chunk := arr[4*i : 4*(i+1)]
		ret = append(ret, zip(chunk[0], chunk[1], chunk[2], chunk[3]))
	}
	return ret
}
