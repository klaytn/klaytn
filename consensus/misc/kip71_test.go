package misc

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
)

func getTestConfig(forkedBlockNum *big.Int) *params.ChainConfig {
	testConfig := params.CypressChainConfig
	testConfig.UnitPrice = uint64(25000000000)
	testConfig.MagmaCompatibleBlock = forkedBlockNum
	testConfig.Governance = &params.GovernanceConfig{
		KIP71: params.GetDefaultKIP71Config(),
	}
	return testConfig
}

func TestEvenBaseFee(t *testing.T) {
	tests := []struct {
		upperBoundBaseFee         uint64
		lowerBoundBaseFee         uint64
		gasTarget                 uint64
		baseFeeDenominator        uint64
		maxBlockGasUsedForBaseFee uint64
		parentGasUsed             uint64
		parentBaseFee             uint64
	}{
		// Current default setting
		{750000000000, 25000000000, 30000000, 20, 60000000, 43212345, 13009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 34857284, 83009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 0, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 0, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 0, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 69856837, 34958439443},
		{750000000000, 25000000000, 30000000, 20, 60000000, 18275847, 43459029443},
		{750000000000, 25000000000, 30000000, 20, 60000000, 34857359, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 28914728, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 83597392, 43238573843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 28472874, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 34895734, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 17274858, 18432574843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 39093494, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 18981884, 28574384343},
		{750000000000, 25000000000, 30000000, 20, 60000000, 12873, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 18273, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 1, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 18949, 43009212843},
		{750000000000, 25000000000, 30000000, 20, 60000000, 19828, 43009212843},
		// monkey test 1
		{750000000001, 25123987327, 30000000, 20, 38573282, 83212345, 750000000000},
		{750000000003, 25000000000, 30000000, 20, 83478211, 43212345, 48572839123},
		{750000000005, 25000000000, 30000000, 20, 48471293, 43212345, 48572839123},
		{750000000006, 25000000000, 30000000, 20, 43852728, 43212345, 48572839123},
		{750000000007, 25000000000, 30000000, 20, 11282848, 43212345, 48572839123},
		{750000000009, 25000000000, 30000000, 20, 38574848, 43212345, 48572839123},
		{800000001233, 25000000000, 30000000, 20, 83839191, 43212345, 48572839123},
		{751231231233, 25000000000, 30000000, 20, 88391928, 43212345, 48572839123},
		{750321321321, 25000000000, 30000000, 20, 60000000, 43212345, 48572839123},
		{750000000777, 25000000000, 30000000, 20, 60000000, 43212345, 48572839123},
		// monkey test 2
		{750000000000, 28347289478, 30000000, 20, 60000000, 43212345, 75489234128},
		{750000000000, 28729584848, 30000000, 20, 60000000, 43212345, 75489234128},
		{0, 0, 30000000, 20, 60000000, 0, 43009212843},
		{750000000000, 29999999999, 30000000, 20, 60000000, 43212345, 75489234128},
		{750000000000, 19999999999, 30000000, 20, 60000000, 43212345, 75489234128},
		{750000000000, 0, 30000000, 20, 60000000, 43212345, 75489234128},
		{750000000000, 1, 30000000, 20, 60000000, 43212345, 75489234128},
		{750000000000, 0, 30000000, 20, 60000000, 43212345, 75489234128},
		{0, 0, 30000000, 20, 60000000, 43212345, 75489234128},
		{0, 0, 30000000, 20, 60000000, 43212345, 75489234128},
	}

	testConfig := getTestConfig(big.NewInt(3))
	for _, test := range tests {
		testConfig.Governance.KIP71.LowerBoundBaseFee = test.lowerBoundBaseFee
		testConfig.Governance.KIP71.UpperBoundBaseFee = test.upperBoundBaseFee
		testConfig.Governance.KIP71.GasTarget = test.gasTarget
		testConfig.Governance.KIP71.MaxBlockGasUsedForBaseFee = test.maxBlockGasUsedForBaseFee
		testConfig.Governance.KIP71.BaseFeeDenominator = test.baseFeeDenominator

		parent := &types.Header{
			Number:  common.Big3,
			GasUsed: test.parentGasUsed,
			BaseFee: new(big.Int).SetUint64(test.parentBaseFee),
		}
		if even := NextBlockBaseFee(parent, getTestConfig(big.NewInt(3))); even.Bit(0) != 0 {
			t.Errorf("NextBlockBaseFee:%d is not a even number", even)
		}
	}
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

func TestNextBlockBaseFeeWhenGovernanceUpdated(t *testing.T) {
	tests := []struct {
		upperBoundBaseFee uint64 // updated upper bound
		lowerBoundBaseFee uint64 // updated lower bound
		parentGasUsed     uint64
		parentBaseFee     int64
		nextBaseFee       int64
	}{
		{750000000000, 25000000000, 30000000, 750000000000, 750000000000},
		{700000000000, 25000000000, 30000000, 750000000000, 700000000000},
		{800000000000, 25000000000, 30000000, 750000000000, 750000000000},
		{750000000000, 25000000000, 30000000, 25000000000, 25000000000},
		{750000000000, 30000000000, 30000000, 25000000000, 30000000000},
		{750000000000, 20000000000, 30000000, 25000000000, 25000000000},
		{750000000000, 25000000000, 40000000 /* > gasTarget(30000000) */, 750000000000, 750000000000},
		{750000000000, 25000000000, 20000000 /* < gasTarget(30000000) */, 25000000000, 25000000000},
	}
	for i, test := range tests {
		config := getTestConfig(common.Big2)
		config.Governance.KIP71.UpperBoundBaseFee = test.upperBoundBaseFee
		config.Governance.KIP71.LowerBoundBaseFee = test.lowerBoundBaseFee
		parent := &types.Header{
			Number:  common.Big3,
			GasUsed: test.parentGasUsed,
			BaseFee: big.NewInt(test.parentBaseFee),
		}
		if have, want := NextBlockBaseFee(parent, config), big.NewInt(test.nextBaseFee); have.Cmp(want) != 0 {
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

func BenchmarkNextBlockBaseFeeRandom(b *testing.B) {
	parentBaseFee := big.NewInt(500000000000)
	parent := &types.Header{
		Number:  common.Big3,
		GasUsed: 10000000,
		BaseFee: parentBaseFee,
	}
	for i := 0; i < b.N; i++ {
		if rand.Int()%2 == 0 {
			parent.GasUsed = 10000000
		} else {
			parent.GasUsed = 40000000
		}
		_ = NextBlockBaseFee(parent, getTestConfig(big.NewInt(2)))
	}
}

func BenchmarkNextBlockBaseFeeUpperBound(b *testing.B) {
	parentBaseFee := big.NewInt(750000000000)
	parent := &types.Header{
		Number:  common.Big3,
		GasUsed: 40000000,
		BaseFee: parentBaseFee,
	}
	for i := 0; i < b.N; i++ {
		_ = NextBlockBaseFee(parent, getTestConfig(big.NewInt(2)))
	}
}

func BenchmarkNextBlockBaseFeeLowerBound(b *testing.B) {
	parentBaseFee := big.NewInt(25000000000)
	parent := &types.Header{
		Number:  common.Big3,
		GasUsed: 10000000,
		BaseFee: parentBaseFee,
	}
	for i := 0; i < b.N; i++ {
		_ = NextBlockBaseFee(parent, getTestConfig(big.NewInt(2)))
	}
}
