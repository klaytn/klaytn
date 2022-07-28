package sc

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/common"
	bridgecontract "github.com/klaytn/klaytn/contracts/bridge"
)

// suggestLeastFee calculates least fee to be set as value transfer service fee.
// It returns by addition of four estimated gas of contract calls (`BridgeTransfer.sol/removeRefundLedger`, `BridgeTransfer.sol/updateHandleStatus`, 'handleValueTransfer')
func (bi *BridgeInfo) suggestLeastFee(ctbi *BridgeInfo, handleValueTransfers []string) (map[string]interface{}, error) {
	reqNonce, err := bi.bridge.RequestNonce(nil)
	if err != nil {
		return nil, err
	}
	unhandledReqNonce := reqNonce + 10000
	// TODO-hyunsooda: Consider two options below
	// (1) Return the gratest cost from `handleKLAYTransfer`, `handleValueTransfer for ERC20 and ERC721`.
	// (2) Or, consider a fine-grained fee per types of value (token)
	methods := []string{"removeRefundLedger", "updateHandleStatus"}
	methods = append(methods, handleValueTransfers...)
	params := [][]interface{}{
		{unhandledReqNonce},
		{unhandledReqNonce, common.HexToHash("0"), uint64(1), false},
		{common.HexToHash("0"), common.HexToAddress("0x1"), common.HexToAddress("0x2"), common.Big0, unhandledReqNonce, uint64(1), []byte{}},
	}
	bridgeInfos := []*BridgeInfo{bi, ctbi, ctbi}
	totalCost, totalGasUsed := uint64(0), uint64(0)
	output := make(map[string]interface{})
	for idx, method := range methods {
		if leastGasLimit, gasPrice, err := bridgeInfos[idx].requestGasEstimate(method, params[idx]...); err == nil {
			output[method] = map[string]interface{}{
				"requiredGas":                       leastGasLimit,
				"gasPrice":                          gasPrice,
				"totalCost(requiredGas * gasPrice)": leastGasLimit * gasPrice,
			}
			totalGasUsed += leastGasLimit
			totalCost += leastGasLimit * gasPrice
		} else {
			return nil, err
		}
	}
	output["SumOfGasUsed"] = totalGasUsed
	output["SumOfCost"] = totalCost
	return output, nil
}

func (bi *BridgeInfo) requestGasEstimate(method string, params ...interface{}) (uint64, uint64, error) {
	parsed, err := abi.JSON(strings.NewReader(bridgecontract.BridgeABI))
	if err != nil {
		return 0, 0, err
	}
	input, err := parsed.Pack(method, params...)
	if err != nil {
		return 0, 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	from := common.Address{}
	gasPrice := common.Big0
	var backend Backend
	if bi.onChildChain {
		acc := bi.subBridge.bridgeAccounts.cAccount
		if bi.subBridge.blockchain.Config().IsMagmaForkEnabled(bi.subBridge.blockchain.CurrentHeader().Number) {
			gasPrice = new(big.Int).SetUint64(bi.subBridge.blockchain.Config().Governance.KIP71.UpperBoundBaseFee)
		} else {
			gasPrice = new(big.Int).SetUint64(bi.subBridge.blockchain.Config().UnitPrice)
		}
		from = acc.address
		backend = bi.subBridge.localBackend
	} else {
		acc := bi.subBridge.bridgeAccounts.pAccount
		if acc.kip71Config.UpperBoundBaseFee != 0 {
			gasPrice = new(big.Int).SetUint64(acc.kip71Config.UpperBoundBaseFee)
		} else {
			gasPrice = acc.gasPrice
		}
		from = acc.address
		backend = bi.subBridge.remoteBackend
	}

	msg := klaytn.CallMsg{From: from, To: &bi.address, GasPrice: gasPrice, Value: common.Big0, Data: input}
	leastGasLimit, err := backend.EstimateGas(ctx, msg)
	return leastGasLimit, gasPrice.Uint64(), err
}
