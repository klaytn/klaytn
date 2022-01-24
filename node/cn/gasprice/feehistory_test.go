package gasprice

import (
	"context"
	"errors"
	"testing"

	"github.com/klaytn/klaytn/networks/rpc"
)

func TestFeeHistory(t *testing.T) {
	var cases = []struct {
		maxHeader, maxBlock int
		count               int
		last                rpc.BlockNumber
		percent             []float64
		expFirst            uint64
		expCount            int
		expErr              error
	}{
		{1000, 1000, 10, 30, nil, 21, 10, nil},
		{1000, 1000, 10, 30, []float64{0, 10}, 21, 10, nil},
		{1000, 1000, 10, 30, []float64{20, 10}, 0, 0, errInvalidPercentile},
		{1000, 1000, 1000000000, 30, nil, 0, 31, nil},
		{1000, 1000, 1000000000, rpc.LatestBlockNumber, nil, 0, 33, nil},
		{1000, 1000, 10, 40, nil, 0, 0, errRequestBeyondHead},
		{1000, 1000, 10, 40, nil, 0, 0, errRequestBeyondHead},
		{20, 2, 100, rpc.LatestBlockNumber, nil, 13, 20, nil},
		{20, 2, 100, rpc.LatestBlockNumber, []float64{0, 10}, 31, 2, nil},
		{20, 2, 100, 32, []float64{0, 10}, 31, 2, nil},
		{1000, 1000, 1, rpc.PendingBlockNumber, nil, 0, 0, nil},
		{1000, 1000, 2, rpc.PendingBlockNumber, nil, 32, 1, nil},
	}
	for i, c := range cases {
		config := Config{
			MaxHeaderHistory: c.maxHeader,
			MaxBlockHistory:  c.maxBlock,
		}
		backend := newTestBackend(t)
		oracle := NewOracle(backend, config)

		first, reward, baseFee, ratio, err := oracle.FeeHistory(context.Background(), c.count, c.last, c.percent)

		expReward := c.expCount
		if len(c.percent) == 0 {
			expReward = 0
		}
		expBaseFee := c.expCount
		if expBaseFee != 0 {
			expBaseFee++
		}

		if first.Uint64() != c.expFirst {
			t.Fatalf("Test case %d: first block mismatch, want %d, got %d", i, c.expFirst, first)
		}
		if len(reward) != expReward {
			t.Fatalf("Test case %d: reward array length mismatch, want %d, got %d", i, expReward, len(reward))
		}
		if len(baseFee) != expBaseFee {
			t.Fatalf("Test case %d: baseFee array length mismatch, want %d, got %d", i, expBaseFee, len(baseFee))
		}
		if len(ratio) != c.expCount {
			t.Fatalf("Test case %d: gasUsedRatio array length mismatch, want %d, got %d", i, c.expCount, len(ratio))
		}
		if err != c.expErr && !errors.Is(err, c.expErr) {
			t.Fatalf("Test case %d: error mismatch, want %v, got %v", i, c.expErr, err)
		}
	}
}
