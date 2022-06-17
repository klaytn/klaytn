package misc

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
)

func getTestConfig(forkedBlockNum *big.Int) *params.ChainConfig {
	testConfig := params.CypressChainConfig
	testConfig.UnitPrice = uint64(25000000000)
	testConfig.KIP71CompatibleBlock = forkedBlockNum
	testConfig.Governance = &params.GovernanceConfig{
		KIP71: params.GetDefaultKip71Config(),
	}
	return testConfig
}

func TestNextBlockBaseFee(t *testing.T) {
	tests := []struct {
		parentBaseFee int64
		parentGasUsed uint64
		nextBaseFee   int64
	}{
		{750000000000, 30000000, 750000000000}, // usage == target
		{30000000000, 20000000, 29500000000},   // usage below target
		{300000000000, 40000000, 305000000000}, // usage above target
	}
	for i, test := range tests {
		parent := &types.Header{
			Number:  common.Big3,
			GasUsed: test.parentGasUsed,
			BaseFee: big.NewInt(test.parentBaseFee),
		}
		if have, want := NextBlockBaseFee(parent, getTestConfig(big.NewInt(3))), big.NewInt(test.nextBaseFee); have.Cmp(want) != 0 {
			t.Errorf("test %d: have %d  want %d, ", i, have, want)
		}
	}
}

type BaseFeeTestCase struct {
	genesisParentBaseFee *big.Int
	hardforkedNum        *big.Int
	GasUsed              uint64
	compMethod           func(*big.Int, *big.Int) bool
	expectedBaseFee      *big.Int
	expectedNum          int
}

func TestBlocksToReachExpectedBaseFee(t *testing.T) {
	testCases := []BaseFeeTestCase{
		{
			big.NewInt(25000000000),
			common.Big3,
			84000000,
			func(a *big.Int, b *big.Int) bool { return a.Cmp(b) > 0 },
			big.NewInt(25000000000 * 2),
			15,
		},
		{
			big.NewInt(60000000000),
			common.Big3,
			29000000,
			func(a *big.Int, b *big.Int) bool { return a.Cmp(b) < 0 },
			big.NewInt(60000000000 / 2),
			416,
		},
		{
			big.NewInt(25000000000),
			common.Big3,
			84000000,
			func(a *big.Int, b *big.Int) bool { return a.Cmp(b) == 0 },
			big.NewInt(750000000000),
			70,
		},
		{
			big.NewInt(750000000000),
			common.Big3,
			29000000,
			func(a *big.Int, b *big.Int) bool { return a.Cmp(b) == 0 },
			big.NewInt(25000000000),
			2040,
		},
	}

	for _, testCase := range testCases {
		blocksToReachExpectedBaseFee(t, testCase)
	}
}

func blocksToReachExpectedBaseFee(t *testing.T, testCase BaseFeeTestCase) {
	testConfig := getTestConfig(big.NewInt(3))
	blockNum := 0
	parentBaseFee := testCase.genesisParentBaseFee
	for {
		blockNum++
		parent := &types.Header{
			Number:  testCase.hardforkedNum,
			GasUsed: testCase.GasUsed,
			BaseFee: parentBaseFee,
		}
		parentBaseFee = NextBlockBaseFee(parent, testConfig)

		if testCase.compMethod(parentBaseFee, testCase.expectedBaseFee) {
			break
		}
	}
	if blockNum != int(testCase.expectedNum) {
		t.Errorf("block number %d expected block number %d, have %d didn't reach to %d", blockNum, testCase.expectedNum, parentBaseFee, testCase.expectedBaseFee)
	}
}

func TestInactieDynamicPolicyBeforeForkedBlock(t *testing.T) {
	parentBaseFee := big.NewInt(25000000000)
	parent := &types.Header{
		Number:  common.Big3,
		GasUsed: 84000000,
		BaseFee: parentBaseFee,
	}
	nextBaseFee := NextBlockBaseFee(parent, getTestConfig(big.NewInt(5)))
	if parentBaseFee.Cmp(nextBaseFee) < 0 {
		t.Errorf("before fork, dynamic base fee policy should be inactive, current base fee: %d  next base fee: %d", parentBaseFee, nextBaseFee)
	}
}

func TestActieDynamicPolicyAfterForkedBlock(t *testing.T) {
	parentBaseFee := big.NewInt(25000000000)
	parent := &types.Header{
		Number:  common.Big3,
		GasUsed: 84000000,
		BaseFee: parentBaseFee,
	}
	nextBaseFee := NextBlockBaseFee(parent, getTestConfig(big.NewInt(2)))
	if parentBaseFee.Cmp(nextBaseFee) > 0 {
		t.Errorf("after fork, dynamic base fee policy should be active, current base fee: %d  next base fee: %d", parentBaseFee, nextBaseFee)
	}
}
