package core

import (
	"testing"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/consensus/istanbul/validator"
	"github.com/stretchr/testify/assert"
)

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
		addrs = []common.Address{
			common.HexToAddress("0x6666666666666666666666666666666666666666"),
			common.HexToAddress("0x5555555555555555555555555555555555555555"),
			common.HexToAddress("0x4444444444444444444444444444444444444444"),
			common.HexToAddress("0x3333333333333333333333333333333333333333"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
		}
		committee = istanbul.Validators{}
		timeMap   = make(map[common.Address]time.Duration)
		expected  = make([]time.Duration, len(addrs))
	)

	for i, addr := range addrs {
		committee = append(committee, validator.New(addr))
		t := time.Duration((i * 100) * int(time.Millisecond))
		timeMap[addr] = t
		expected[len(expected)-1-i] = t
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
		assert.Equal(t, compress(tc.input), tc.expected, "tc %d failed", i)
	}
}
