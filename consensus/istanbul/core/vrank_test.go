package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
