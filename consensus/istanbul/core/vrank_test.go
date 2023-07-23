package core

import (
	"encoding/hex"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/consensus/istanbul/validator"
	"github.com/stretchr/testify/assert"
)

func genCommitteeFromAddrs(addrs []common.Address) istanbul.Validators {
	committee := []istanbul.Validator{}
	for _, addr := range addrs {
		committee = append(committee, validator.New(addr))
	}
	return committee
}

func TestVrank(t *testing.T) {
	var (
		N         = 6
		quorum    = 4
		addrs, _  = genValidators(N)
		committee = genCommitteeFromAddrs(addrs)
		view      = istanbul.View{Sequence: big.NewInt(0), Round: big.NewInt(0)}
		msg       = &istanbul.Subject{View: &view}
		vrank     = NewVrank(view, committee)

		expectedAssessList  []int
		expectedLateCommits []time.Duration
	)

	sort.Sort(committee)
	for i := 0; i < quorum; i++ {
		vrank.AddCommit(msg, committee[i])
		expectedAssessList = append(expectedAssessList, vrankArrivedEarly)
	}
	vrank.HandleCommitted(view.Sequence)
	for i := quorum; i < N; i++ {
		t := vrank.AddCommit(msg, committee[i])
		expectedAssessList = append(expectedAssessList, vrankArrivedLate)
		expectedLateCommits = append(expectedLateCommits, t)
	}

	bitmap, late := vrank.Log()
	assert.Equal(t, hex.EncodeToString(compress(expectedAssessList)), bitmap)
	assert.Equal(t, encodeDurationBatch(expectedLateCommits), late)
}

func TestVrankAssessBatch(t *testing.T) {
	arr := []time.Duration{0, 4, 1, vrankNotArrivedPlaceholder, 2}
	threshold := time.Duration(2)
	expected := []int{
		vrankArrivedEarly, vrankArrivedLate, vrankArrivedEarly, vrankNotArrived, vrankArrivedEarly,
	}
	assert.Equal(t, expected, assessBatch(arr, threshold))
}

func TestVrankSerialize(t *testing.T) {
	var (
		N         = 6
		addrs, _  = genValidators(N)
		committee = genCommitteeFromAddrs(addrs)
		timeMap   = make(map[common.Address]time.Duration)
		expected  = make([]time.Duration, len(addrs))
	)

	sort.Sort(committee)
	for i, val := range committee {
		t := time.Duration((i * 100) * int(time.Millisecond))
		timeMap[val.Address()] = t
		expected[i] = t
	}

	assert.Equal(t, expected, serialize(committee, timeMap))
}

func TestVrankCompress(t *testing.T) {
	tcs := []struct {
		input    []int
		expected []byte
	}{
		{
			input:    []int{2},
			expected: []byte{0b10_000000},
		},
		{
			input:    []int{2, 1},
			expected: []byte{0b10_01_0000},
		},
		{
			input:    []int{0, 2, 1},
			expected: []byte{0b00_10_01_00},
		},
		{
			input:    []int{0, 2, 1, 1},
			expected: []byte{0b00_10_01_01},
		},
		{
			input:    []int{1, 2, 1, 2, 1},
			expected: []byte{0b01_10_01_10, 0b01_000000},
		},
		{
			input:    []int{1, 2, 1, 2, 1, 2},
			expected: []byte{0b01_10_01_10, 0b01_10_0000},
		},
		{
			input:    []int{1, 2, 1, 2, 1, 2, 1},
			expected: []byte{0b01_10_01_10, 0b01_10_01_00},
		},
		{
			input:    []int{1, 2, 1, 2, 1, 2, 0, 2},
			expected: []byte{0b01_10_01_10, 0b01_10_00_10},
		},
		{
			input:    []int{1, 1, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2, 1, 1, 2, 2, 1, 1},
			expected: []byte{0b01011010, 0b01011010, 0b01011010, 0b01011010, 0b01010000},
		},
	}
	for i, tc := range tcs {
		assert.Equal(t, tc.expected, compress(tc.input), "tc %d failed", i)
	}
}
