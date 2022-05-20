package misc

import (
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/params"
)

func VerifyKIP71Header(config *params.ChainConfig, parentHeader, header *types.Header) error {
	if header.BaseFee == nil {
		return fmt.Errorf("header is missing baseFee")
	}
	// Verify the baseFee is correct based on the parent header.
	expectedBaseFee := CalcBaseFee(parentHeader, config)
	if header.BaseFee.Cmp(expectedBaseFee) != 0 {
		return fmt.Errorf("invalid baseFee: have %s, want %s, parentBaseFee %s, parentGasUsed %d",
			header.BaseFee, expectedBaseFee, parentHeader.BaseFee, parentHeader.GasUsed)
	}
	return nil
}

func CalcBaseFee(parentHeader *types.Header, config *params.ChainConfig) *big.Int {
	// If the current is the kip71 disabled block, then return default base fee (250ston)
	if !config.IsKIP71ForkEnabled(parentHeader.Number) {
		return new(big.Int).SetUint64(config.UnitPrice)
	}

	// governance parameters
	lowerBoundBaseFee := new(big.Int).SetUint64(config.Governance.KIP71.LowerBoundBaseFee)
	upperBoundBaseFee := new(big.Int).SetUint64(config.Governance.KIP71.UpperBoundBaseFee)
	baseFeeDenominator := new(big.Int).SetUint64(config.Governance.KIP71.BaseFeeDenominator)
	gasTarget := config.Governance.KIP71.GasTarget
	upperGasLimit := config.Governance.KIP71.BlockGasLimit

	parentBaseFee := parentHeader.BaseFee
	parentGasUsed := parentHeader.GasUsed
	// upper gas limit cut off the impulse of used gas to upper bound
	if parentGasUsed > upperGasLimit {
		parentGasUsed = upperGasLimit
	}
	if parentGasUsed == gasTarget {
		return new(big.Int).Set(parentHeader.BaseFee)
	} else if parentGasUsed > gasTarget {
		// If the parent block used more gas than its target,
		// the baseFee of the next block should increase.
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
		// Otherwise if the parent block used less gas than its target,
		// the baseFee of the next block should decrease.
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
