package database

import (
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

var (
	// Deprecated
	valueTransferTxHashPrefix = []byte("vt-tx-hash-key-") // Prefix + hash -> hash

	requestValueTransferTxPrefix = []byte("vt-request-tx-hash-key-")      // Prefix + bridge address + counterpart bridge addr + request hash -> key
	refundTxKeyPrefix            = []byte("vt-refund-request-nonce-key-") // Prefix + bridge address + counterpart bridge addr + request hash -> key
	failedRequestNonceKeyPrefix  = []byte("vt-failed-request-nonce-")     // Prefix + bridge address + counterpart bridge addr + request hash -> key
)

type BridgeRequestEvent struct {
	TokenType      uint8
	From           common.Address
	To             common.Address
	TokenAddr      common.Address
	ValueOrTokenId *big.Int
	RequestNonce   uint64
	Fee            *big.Int
	ExtraData      []byte
}

type HandleInfo struct {
	RequestEvent  BridgeRequestEvent
	RequestTxHash common.Hash
	HandleTx      types.Transaction
}

type FailedHandleInfo struct {
	RequestEvent  BridgeRequestEvent
	RequestTxHash common.Hash
}

type RefundInfo struct {
	RequestEvent  BridgeRequestEvent
	RequestTxHash common.Hash
	RefundTx      types.Transaction
}

func NewBridgeRequestEvent(tokenType uint8, from, to, tokenAddr common.Address,
	valueOrTokenId *big.Int, reqNonce uint64, fee *big.Int, extraData []byte) BridgeRequestEvent {
	return BridgeRequestEvent{
		tokenType,
		from,
		to,
		tokenAddr,
		valueOrTokenId,
		reqNonce,
		fee,
		extraData,
	}
}

func (ev BridgeRequestEvent) makeRPCOutput() map[string]interface{} {
	// TODO-hyunsooda: Assign thru reflect loop
	m := make(map[string]interface{})
	var tokenType string
	switch ev.TokenType {
	case 0:
		tokenType = "KLAY"
	case 1:
		tokenType = "ERC20"
	case 2:
		tokenType = "ERC721"
	}
	m["TokenType"] = tokenType
	m["From"] = ev.From
	m["To"] = ev.To
	m["TokenAddr"] = ev.TokenAddr
	m["ValueOrTokenId"] = ev.ValueOrTokenId
	m["RequestNonce"] = ev.RequestNonce
	m["Fee"] = ev.Fee
	m["ExtraData"] = ev.ExtraData
	return m
}

func makeTxRPCOutput(tx *types.Transaction) map[string]interface{} {
	if tx == nil {
		return nil
	}
	output := tx.MakeRPCOutput()
	output["hash"] = tx.Hash()
	output["cost"] = tx.Cost()
	output["fee"] = tx.Fee()
	output["size"] = tx.Size()
	return output
}

func makeRPCOutput(ev *BridgeRequestEvent, reqTxHash common.Hash, handleTx *types.Transaction) map[string]interface{} {
	output := ev.makeRPCOutput()
	output["RequestTxHash"] = reqTxHash
	if handleTx != nil {
		for prop, value := range makeTxRPCOutput(handleTx) {
			output[prop] = value
		}
	}
	return output
}

func (handleInfo HandleInfo) MakeRPCOutput() map[string]interface{} {
	return makeRPCOutput(&handleInfo.RequestEvent, handleInfo.RequestTxHash, &handleInfo.HandleTx)
}

func (handleInfo FailedHandleInfo) MakeRPCOutput() map[string]interface{} {
	return makeRPCOutput(&handleInfo.RequestEvent, handleInfo.RequestTxHash, nil)
}

func (refundInfo RefundInfo) MakeRPCOutput() map[string]interface{} {
	return makeRPCOutput(&refundInfo.RequestEvent, refundInfo.RequestTxHash, &refundInfo.RefundTx)
}

// Deprecated
func valueTransferTxHashKey(rTxHash common.Hash) []byte {
	return append(valueTransferTxHashPrefix, rTxHash.Bytes()...)
}

func reuqestValueTransferHashKey(bridgeAddr, counterpartBridgeAddr common.Address, rTxHash common.Hash) []byte {
	return append(append(append(requestValueTransferTxPrefix, bridgeAddr.Bytes()...), counterpartBridgeAddr.Bytes()...), rTxHash.Bytes()...)
}

func refundTxKey(bridgeAddr, counterpartBridgeAddr common.Address, rTxHash common.Hash) []byte {
	return append(append(append(refundTxKeyPrefix, bridgeAddr.Bytes()...), counterpartBridgeAddr.Bytes()...), rTxHash.Bytes()...)
}

func failedRequestNonceKey(bridgeAddr, counterpartBridgeAddr common.Address, rTxHash common.Hash) []byte {
	return append(append(append(failedRequestNonceKeyPrefix, bridgeAddr.Bytes()...), counterpartBridgeAddr.Bytes()...), rTxHash.Bytes()...)
}
