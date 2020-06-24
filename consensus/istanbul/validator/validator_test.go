package validator

import (
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/klaytn/klaytn/common"
)

func TestCalcSeed(t *testing.T) {
	type testCase struct {
		hash        common.Hash
		expected    int64
		expectedErr error
	}
	testCases := []testCase{
		{
			hash:        common.HexToHash("0x1111"),
			expected:    0,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234000000000000000000000000000000000000000000000000000000000000"),
			expected:    81979586966978560,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341230000000000000000000000000000000000000000000000000"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341234000000000000000000000000000000000000000000000000"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341234123412341234123412341234123412341234123412341234"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x0034123412341234123412341234123412341234123412341234123412341234"),
			expected:    916044602622243,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0xabcdef3412341234123412341234123412341234123412341234123412341234"),
			expected:    773738372352131363,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		actual, err := ConvertHashToSeed(tc.hash)
		assert.Equal(t, err, tc.expectedErr)
		assert.Equal(t, actual, tc.expected)
	}
}
