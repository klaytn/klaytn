package misc

import (
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/params"
)

func VerifyMagmaHeader(parentHeader, header *types.Header, kip71Config *params.KIP71Config) error {
	if header.BaseFee == nil {
		return fmt.Errorf("header is missing baseFee")
	}
	// Verify the baseFee is correct based on the parent header.
	expectedBaseFee := NextMagmaBlockBaseFee(parentHeader, kip71Config)
	if header.BaseFee.Cmp(expectedBaseFee) != 0 {
		return fmt.Errorf("invalid baseFee: have %s, want %s, parentBaseFee %s, parentGasUsed %d",
			header.BaseFee, expectedBaseFee, parentHeader.BaseFee, parentHeader.GasUsed)
	}
	return nil
}

func makeEvenByFloor(baseFee *big.Int) *big.Int {
	if baseFee.Bit(0) != 0 {
		baseFee.Sub(baseFee, common.Big1)
	}
	return baseFee
}

func makeEvenByCeil(baseFee *big.Int) *big.Int {
	if baseFee.Bit(0) != 0 {
		baseFee.Add(baseFee, common.Big1)
	}
	return baseFee
}

func NextMagmaBlockBaseFee(parentHeader *types.Header, kip71Config *params.KIP71Config) *big.Int {
	// governance parameters
	lowerBoundBaseFee := new(big.Int).SetUint64(kip71Config.LowerBoundBaseFee)
	upperBoundBaseFee := new(big.Int).SetUint64(kip71Config.UpperBoundBaseFee)
	makeEvenByCeil(lowerBoundBaseFee)
	makeEvenByFloor(upperBoundBaseFee)

	nextFee := nextBlockBaseFee(parentHeader, kip71Config, lowerBoundBaseFee, upperBoundBaseFee)
	return makeEvenByFloor(nextFee)
}

func nextBlockBaseFee(parentHeader *types.Header, kip71Config *params.KIP71Config, lowerBoundBaseFee, upperBoundBaseFee *big.Int) *big.Int {
	// If the parent is the magma disabled block or genesis, then return the lowerBoundBaseFee (default 25ston)
	if parentHeader.Number.Cmp(new(big.Int).SetUint64(0)) == 0 || parentHeader.BaseFee == nil {
		return lowerBoundBaseFee
	}

	var baseFeeDenominator *big.Int
	if kip71Config.BaseFeeDenominator == 0 {
		// To avoid panic, set the fluctuation range small
		baseFeeDenominator = new(big.Int).SetUint64(64)
	} else {
		baseFeeDenominator = new(big.Int).SetUint64(kip71Config.BaseFeeDenominator)
	}
	gasTarget := kip71Config.GasTarget
	upperGasLimit := kip71Config.MaxBlockGasUsedForBaseFee

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
		return parentBaseFee
	} else if parentGasUsed > gasTarget {
		// shortcut. If parentBaseFee is already reached upperbound, do not calculate.
		if parentBaseFee.Cmp(upperBoundBaseFee) == 0 {
			return upperBoundBaseFee
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
			return upperBoundBaseFee
		}
		return nextBaseFee
	} else {
		// shortcut. If parentBaseFee is already reached lower bound, do not calculate.
		if parentBaseFee.Cmp(lowerBoundBaseFee) == 0 {
			return lowerBoundBaseFee
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
			return lowerBoundBaseFee
		}
		return nextBaseFee
	}
}
