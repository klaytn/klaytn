package database

import (
	"math/big"
	"reflect"

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
	valueOrTokenId *big.Int, reqNonce uint64, fee *big.Int, extraData []byte,
) BridgeRequestEvent {
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
	m := make(map[string]interface{})
	val := reflect.ValueOf(&ev).Elem()
	for i := 0; i < reflect.TypeOf(&ev).Elem().NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		fieldVal := reflect.Indirect(val).FieldByName(fieldName)
		m[fieldName] = fieldVal.Interface()
	}
	switch ev.TokenType {
	case 0:
		m["TokenType"] = "KLAY"
	case 1:
		m["TokenType"] = "ERC20"
	case 2:
		m["TokenType"] = "ERC721"
	}
	return m
}

func makeTxRPCOutput(tx *types.Transaction) map[string]interface{} {
	if tx == nil {
		return nil
	}
	output := tx.MakeRPCOutput()
	output["HandleTxHash"] = tx.Hash()
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
