package misc

import (
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/params"
)

func VerifyMagmaHeader(config *params.ChainConfig, parentHeader, header *types.Header) error {
	if header.BaseFee == nil {
		return fmt.Errorf("header is missing baseFee")
	}
	// Verify the baseFee is correct based on the parent header.
	expectedBaseFee := NextBlockBaseFee(parentHeader, config)
	if header.BaseFee.Cmp(expectedBaseFee) != 0 {
		return fmt.Errorf("invalid baseFee: have %s, want %s, parentBaseFee %s, parentGasUsed %d",
			header.BaseFee, expectedBaseFee, parentHeader.BaseFee, parentHeader.GasUsed)
	}
	return nil
}

// If the upperBoundBaseFee is Odd number then the baseFee can be upperBoundBaseFee + 1
func makeEven(baseFee *big.Int) *big.Int {
	if baseFee.Bit(0) != 0 {
		return baseFee.Add(baseFee, common.Big1)
	}
	return baseFee
}

func NextBlockBaseFee(parentHeader *types.Header, config *params.ChainConfig) *big.Int {
	// governance parameters
	lowerBoundBaseFee := new(big.Int).SetUint64(config.Governance.KIP71.LowerBoundBaseFee)
	upperBoundBaseFee := new(big.Int).SetUint64(config.Governance.KIP71.UpperBoundBaseFee)

	// If the parent is the magma disabled block or genesis, then return the lowerBoundBaseFee (default 25ston)
	if !config.IsMagmaForkEnabled(parentHeader.Number) || parentHeader.Number.Cmp(new(big.Int).SetUint64(0)) == 0 {
		return makeEven(lowerBoundBaseFee)
	}

	var baseFeeDenominator *big.Int
	if config.Governance.KIP71.BaseFeeDenominator == 0 {
		// To avoid panic, set the fluctuation range small
		baseFeeDenominator = new(big.Int).SetUint64(64)
	} else {
		baseFeeDenominator = new(big.Int).SetUint64(config.Governance.KIP71.BaseFeeDenominator)
	}
	gasTarget := config.Governance.KIP71.GasTarget
	upperGasLimit := config.Governance.KIP71.MaxBlockGasUsedForBaseFee

	// check the case of upper/lowerBoundBaseFee is updated by governance mechanism
	parentBaseFee := parentHeader.BaseFee
	if parentBaseFee.Cmp(upperBoundBaseFee) >= 0 {
		parentBaseFee = upperBoundBaseFee
	} else if parentBaseFee.Cmp(lowerBoundBaseFee) <= 0 {
		parentBaseFee = lowerBoundBaseFee
	}

	parentGasUsed := parentHeader.GasUsed
	// upper gas limit cut off the impulse of used gas to upper bound
	if parentGasUsed > upperGasLimit {
		parentGasUsed = upperGasLimit
	}
	if parentGasUsed == gasTarget {
		return makeEven(parentBaseFee)
	} else if parentGasUsed > gasTarget {
		// shortcut. If parentBaseFee is already reached upperbound, do not calculate.
		if parentBaseFee.Cmp(upperBoundBaseFee) == 0 {
			return makeEven(upperBoundBaseFee)
		}
		// If the parent block used more gas than its target,
		// the baseFee of the next block should increase.
		// baseFeeDelta = max(1, parentBaseFee * (parentGasUsed - gasTarget) / gasTarget / baseFeeDenominator)
		gasUsedDelta := new(big.Int).SetUint64(parentGasUsed - gasTarget)
		x := new(big.Int).Mul(parentBaseFee, gasUsedDelta)
		y := x.Div(x, new(big.Int).SetUint64(gasTarget))
		baseFeeDelta := math.BigMax(x.Div(y, baseFeeDenominator), common.Big1)

		nextBaseFee := x.Add(parentBaseFee, baseFeeDelta)
		if nextBaseFee.Cmp(upperBoundBaseFee) > 0 {
			return makeEven(upperBoundBaseFee)
		}
		return makeEven(nextBaseFee)
	} else {
		// shortcut. If parentBaseFee is already reached lower bound, do not calculate.
		if parentBaseFee.Cmp(lowerBoundBaseFee) == 0 {
			return makeEven(lowerBoundBaseFee)
		}
		// Otherwise if the parent block used less gas than its target,
		// the baseFee of the next block should decrease.
		// baseFeeDelta = parentBaseFee * (gasTarget - parentGasUsed) / gasTarget / baseFeeDenominator
		gasUsedDelta := new(big.Int).SetUint64(gasTarget - parentGasUsed)
		x := new(big.Int).Mul(parentBaseFee, gasUsedDelta)
		y := x.Div(x, new(big.Int).SetUint64(gasTarget))
		baseFeeDelta := x.Div(y, baseFeeDenominator)

		nextBaseFee := x.Sub(parentBaseFee, baseFeeDelta)
		if nextBaseFee.Cmp(lowerBoundBaseFee) < 0 {
			return makeEven(lowerBoundBaseFee)
		}
		return makeEven(nextBaseFee)
	}
}
