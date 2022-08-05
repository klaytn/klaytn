package sc

import (
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	bridgecontract "github.com/klaytn/klaytn/contracts/bridge"
)

type IRequestValueTransferEvent interface {
	Nonce() uint64
	GetTokenType() uint8
	GetFrom() common.Address
	GetTo() common.Address
	GetTokenAddress() common.Address
	GetValueOrTokenId() *big.Int
	GetRequestNonce() uint64
	GetFee() *big.Int
	GetExtraData() []byte

	GetRaw() types.Log

	String() string
}

//////////////////// type RequestValueTransferEvent struct ////////////////////
// RequestValueTransferEvent from Bridge contract
type RequestValueTransferEvent struct {
	*bridgecontract.BridgeRequestValueTransfer
}

func (rEv RequestValueTransferEvent) Nonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEvent) GetTokenType() uint8 {
	return rEv.TokenType
}

func (rEv RequestValueTransferEvent) GetFrom() common.Address {
	return rEv.From
}

func (rEv RequestValueTransferEvent) GetTo() common.Address {
	return rEv.To
}

func (rEv RequestValueTransferEvent) GetTokenAddress() common.Address {
	return rEv.TokenAddress
}

func (rEv RequestValueTransferEvent) GetValueOrTokenId() *big.Int {
	return rEv.ValueOrTokenId
}

func (rEv RequestValueTransferEvent) GetRequestNonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEvent) GetFee() *big.Int {
	return rEv.Fee
}

func (rEv RequestValueTransferEvent) GetExtraData() []byte {
	return rEv.ExtraData
}

func (rEv RequestValueTransferEvent) GetRaw() types.Log {
	return rEv.Raw
}

func (rEv RequestValueTransferEvent) String() string {
	return fmt.Sprintf("requestNonce = %v, tokenType = %v, from = %v, to = %v,tokenAddr = %v, valueOrTokenId = %v, fee = %v, extraData = %v", rEv.RequestNonce, rEv.TokenType, rEv.From.Hex(), rEv.To.Hex(), rEv.TokenAddress.Hex(), rEv.ValueOrTokenId, rEv.Fee, rEv.ExtraData)
}

//////////////////// type RequestValueTransferEncodedEvent struct ////////////////////
type RequestValueTransferEncodedEvent struct {
	*bridgecontract.BridgeRequestValueTransferEncoded
}

func (rEv RequestValueTransferEncodedEvent) Nonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEncodedEvent) GetTokenType() uint8 {
	return rEv.TokenType
}

func (rEv RequestValueTransferEncodedEvent) GetFrom() common.Address {
	return rEv.From
}

func (rEv RequestValueTransferEncodedEvent) GetTo() common.Address {
	return rEv.To
}

func (rEv RequestValueTransferEncodedEvent) GetTokenAddress() common.Address {
	return rEv.TokenAddress
}

func (rEv RequestValueTransferEncodedEvent) GetValueOrTokenId() *big.Int {
	return rEv.ValueOrTokenId
}

func (rEv RequestValueTransferEncodedEvent) GetRequestNonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEncodedEvent) GetFee() *big.Int {
	return rEv.Fee
}

func (rEv RequestValueTransferEncodedEvent) GetExtraData() []byte {
	return rEv.ExtraData
}

func (rEv RequestValueTransferEncodedEvent) GetRaw() types.Log {
	return rEv.Raw
}

func (rEv RequestValueTransferEncodedEvent) String() string {
	return fmt.Sprintf("requestNonce = %v, tokenType = %v, from = %v, to = %v,tokenAddr = %v, valueOrTokenId = %v, fee = %v, extraData = %v, encodingVer = %v, encodingData = %v", rEv.RequestNonce, rEv.TokenType, rEv.From.Hex(), rEv.To.Hex(), rEv.TokenAddress.Hex(), rEv.ValueOrTokenId, rEv.Fee, rEv.ExtraData, rEv.EncodingVer, rEv.EncodedData)
}

// HandleValueTransferEvent from Bridge contract
type HandleValueTransferEvent struct {
	*bridgecontract.BridgeHandleValueTransfer
}

// `RequestRefundEvent` event` from Bridge contract
type RequestRefundEvent struct {
	*bridgecontract.BridgeRequestRefund
}

// `HandleRefundEvent` event` from Bridge contract
type HandleRefundEvent struct {
	*bridgecontract.BridgeHandleRefund
}
