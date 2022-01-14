package blockchain

import (
	"math/big"
	"sync"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func newTestThrottler(config *ThrottlerConfig) *throttler {
	return &throttler{
		config:     config,
		candidates: make(map[common.Address]int),
		throttled:  make(map[common.Address]int),
		allowed:    make(map[common.Address]bool),
		mu:         new(sync.RWMutex),
		threshold:  config.InitialThreshold,
		throttleCh: make(chan *types.Transaction, config.ThrottleTPS*5),
		quitCh:     make(chan struct{}),
	}
}

func TestThrottler_updateThrottlerState(t *testing.T) {
	testConfig := DefaultSpamThrottlerConfig
	testCases := []struct {
		tcName          string
		txFailNum       int
		candidateNum    int
		candidateWeight int
		throttledNum    int
		throttledWeight int
	}{
		{
			tcName:          "one address generates enough fail txs to be throttled",
			txFailNum:       testConfig.InitialThreshold/testConfig.IncreaseWeight + 1,
			candidateNum:    0,
			candidateWeight: 0,
			throttledNum:    1,
			throttledWeight: testConfig.ThrottleSeconds,
		},
		{
			tcName:       "one address generates not enough fail txs to be throttled",
			txFailNum:    testConfig.InitialThreshold / testConfig.IncreaseWeight,
			candidateNum: 1,
			// th.config.IncreaseWeight * th.config.InitialThreshold / th.config.IncreaseWeight - th.config.DecreaseWeight
			candidateWeight: testConfig.InitialThreshold - testConfig.DecreaseWeight,
			throttledNum:    0,
			throttledWeight: 0,
		},
	}

	txNum := 1000
	amount := big.NewInt(0)
	gasLimit := uint64(10000)
	gasPrice := big.NewInt(25 * params.Ston)
	toFail := common.BytesToAddress(common.MakeRandomBytes(20))
	toSuccess := common.BytesToAddress(common.MakeRandomBytes(20))

	for _, tc := range testCases {
		var txs types.Transactions
		var receipts types.Receipts

		// Reset throttler
		th := newTestThrottler(testConfig)

		// Generates fail txs
		for i := 0; i < tc.txFailNum; i++ {
			tx := types.NewTransaction(0, toFail, amount, gasLimit, gasPrice, nil)
			receipt := &types.Receipt{
				Status: types.ReceiptStatusFailed,
			}
			txs = append(txs, tx)
			receipts = append(receipts, receipt)
		}

		// Generate success txs
		for i := 0; i < txNum-tc.txFailNum; i++ {
			tx := types.NewTransaction(0, toSuccess, amount, gasLimit, gasPrice, nil)
			receipt := &types.Receipt{
				Status: types.ReceiptStatusSuccessful,
			}
			txs = append(txs, tx)
			receipts = append(receipts, receipt)
		}

		th.updateThrottlerState(txs, receipts)

		assert.Equal(t, tc.candidateNum, len(th.candidates))
		assert.Equal(t, tc.candidateWeight, th.candidates[toFail])
		assert.Equal(t, tc.throttledNum, len(th.throttled))
		assert.Equal(t, tc.throttledWeight, th.throttled[toFail])
	}
}
