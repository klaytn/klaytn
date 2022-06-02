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

func TestCalcBaseFee(t *testing.T) {
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
		if have, want := CalcBaseFee(parent, getTestConfig(big.NewInt(3))), big.NewInt(test.nextBaseFee); have.Cmp(want) != 0 {
			t.Errorf("test %d: have %d  want %d, ", i, have, want)
		}
	}
}

func TestBlockNumReachDoubleBaseFee(t *testing.T) {
	blockNum := 0
	parentBaseFee := big.NewInt(25000000000)
	for i := 0; i < 15; i++ {
		parent := &types.Header{
			Number:  common.Big3,
			GasUsed: 84000000,
			BaseFee: parentBaseFee,
		}
		parentBaseFee = CalcBaseFee(parent, getTestConfig(big.NewInt(3)))
		blockNum = i
		// t.Logf("test %d: have %d, ", i, parentBaseFee)
	}
	if parentBaseFee.Cmp(big.NewInt(25000000000*2)) < 0 {
		t.Errorf("block number %d: have %d want < %d", blockNum, parentBaseFee, 25000000000*2)
	}
}

func TestBlockNumReachHalfBaseFee(t *testing.T) {
	blockNum := 0
	parentBaseFee := big.NewInt(60000000000)
	for i := 0; i < 749; i++ {
		parent := &types.Header{
			Number:  common.Big3,
			GasUsed: 29000000,
			BaseFee: parentBaseFee,
		}
		parentBaseFee = CalcBaseFee(parent, getTestConfig(big.NewInt(3)))
		blockNum = i
		// t.Logf("test %d: have %d, ", i, parentBaseFee)
	}
	if parentBaseFee.Cmp(big.NewInt(60000000000/2)) > 0 {
		t.Errorf("block number %d: have %d want > %d", blockNum, parentBaseFee, 60000000000/2)
	}
}

func TestBlockNumReacheLowerToMaxBaseFee(t *testing.T) {
	blockNum := 0
	parentBaseFee := big.NewInt(25000000000)
	for i := 0; i < 70; i++ {
		parent := &types.Header{
			Number:  common.Big3,
			GasUsed: 84000000,
			BaseFee: parentBaseFee,
		}
		parentBaseFee = CalcBaseFee(parent, getTestConfig(big.NewInt(3)))
		blockNum = i
		// t.Logf("test %d: have %d, ", i, parentBaseFee)
	}
	if parentBaseFee.Cmp(big.NewInt(750000000000)) != 0 {
		t.Errorf("block number %d: have %d want %d", blockNum, parentBaseFee, 750000000000)
	}
}

func TestBlockNumReachMaxToLowerBaseFee(t *testing.T) {
	blockNum := 0
	parentBaseFee := big.NewInt(750000000000)
	for i := 0; i < 2040; i++ {
		parent := &types.Header{
			Number:  common.Big3,
			GasUsed: 29000000,
			BaseFee: parentBaseFee,
		}
		parentBaseFee = CalcBaseFee(parent, getTestConfig(big.NewInt(3)))
		blockNum = i
		// t.Logf("test %d: have %d, ", i, parentBaseFee)
	}
	if parentBaseFee.Cmp(big.NewInt(25000000000)) != 0 {
		t.Errorf("block number %d: have %d want %d", blockNum, parentBaseFee, 25000000000)
	}
}

func TestInactieDynamicPolicyBeforeForkedBlock(t *testing.T) {
	parentBaseFee := big.NewInt(25000000000)
	parent := &types.Header{
		Number:  common.Big3,
		GasUsed: 84000000,
		BaseFee: parentBaseFee,
	}
	nextBaseFee := CalcBaseFee(parent, getTestConfig(big.NewInt(5)))
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
	nextBaseFee := CalcBaseFee(parent, getTestConfig(big.NewInt(2)))
	if parentBaseFee.Cmp(nextBaseFee) > 0 {
		t.Errorf("after fork, dynamic base fee policy should be active, current base fee: %d  next base fee: %d", parentBaseFee, nextBaseFee)
	}
}
