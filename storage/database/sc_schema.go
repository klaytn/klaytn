package database

import (
	"fmt"
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
	HandleTx      *types.Transaction
}

type FailedHandleInfo struct {
	RequestNonce  uint64
	RequestTxHash common.Hash
}

type RefundInfo struct {
	RequestNonce  uint64
	RequestTxHash common.Hash
	RefundTx      *types.Transaction
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
	output["cost"] = tx.Cost()
	output["fee"] = tx.Fee()
	output["size"] = tx.Size()
	return output
}

func makeRPCOutput(ev *BridgeRequestEvent, reqTxHash common.Hash, handleTx *types.Transaction) map[string]interface{} {
	var output map[string]interface{}
	if ev != nil {
		output = ev.makeRPCOutput()
	} else {
		output = make(map[string]interface{})
	}
	output["RequestTxHash"] = reqTxHash
	if handleTx != nil {
		for prop, value := range makeTxRPCOutput(handleTx) {
			output[prop] = value
		}
	}
	return output
}

func (handleInfo HandleInfo) MakeRPCOutput() map[string]interface{} {
	m := makeRPCOutput(&handleInfo.RequestEvent, handleInfo.RequestTxHash, handleInfo.HandleTx)
	m["HandleTxHash"] = handleInfo.HandleTx.Hash()
	return m
}

func (failedHandleInfo FailedHandleInfo) MakeRPCOutput() map[string]interface{} {
	m := makeRPCOutput(nil, failedHandleInfo.RequestTxHash, nil)
	m["RequestNonce"] = failedHandleInfo.RequestNonce
	return m
}

func (refundInfo RefundInfo) MakeRPCOutput() map[string]interface{} {
	m := makeRPCOutput(nil, refundInfo.RequestTxHash, refundInfo.RefundTx)
	m["RequestNonce"] = refundInfo.RequestNonce
	m["RefundTxHash"] = refundInfo.RefundTx.Hash()
	return m
}

// Deprecated
func valueTransferTxHashKey(rTxHash common.Hash) []byte {
	return append(valueTransferTxHashPrefix, rTxHash.Bytes()...)
}

func reuqestValueTransferHashKey(bridgeAddr, counterpartBridgeAddr common.Address, rTxHash common.Hash) []byte {
	b1, b2 := align(bridgeAddr, counterpartBridgeAddr)
	return append(append(append(requestValueTransferTxPrefix, b1...), b2...), rTxHash.Bytes()...)
}

func refundTxKey(bridgeAddr, counterpartBridgeAddr common.Address, rTxHash common.Hash) []byte {
	fmt.Println("WWWW refund", rTxHash.Hex())
	b1, b2 := align(bridgeAddr, counterpartBridgeAddr)
	return append(append(append(refundTxKeyPrefix, b1...), b2...), rTxHash.Bytes()...)
}

func failedRequestNonceKey(bridgeAddr, counterpartBridgeAddr common.Address, rTxHash common.Hash) []byte {
	b1, b2 := align(bridgeAddr, counterpartBridgeAddr)
	return append(append(append(failedRequestNonceKeyPrefix, b1...), b2...), rTxHash.Bytes()...)
}

func align(addr1, addr2 common.Address) ([]byte, []byte) {
	b1 := addr1.Bytes()
	b2 := addr2.Bytes()
	if b1[0] <= b2[0] {
		b1, b2 = b2, b1
	}
	return b1, b2
}
