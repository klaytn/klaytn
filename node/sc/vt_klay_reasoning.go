// Copyright 2022 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package sc

import (
	"context"
	"io"
	"math/big"
	"reflect"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/rlp"
)

const (
	VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED ValueTransferException = iota
	VALUE_TRANSFER_KLAY_SUCCEEDED_TX
	VALUE_TRANSFER_KLAY_INVALID_VOTE_OR_NONCE
	VALUE_TRANSFER_KLAY_UNEXECUTED // (1) if the operator does not have KLAY to execute its handle value transfer tx, or (2) Bridge node is not servicing (connection down or synching)
	VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE
	VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS // means potential attacker sent KLAY value transfer
	VALUE_TRANSFER_UNREACHABLE
)

var KLAYReasoingResultStringMap = map[ValueTransferException]string{
	VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED:     "VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED",
	VALUE_TRANSFER_KLAY_SUCCEEDED_TX:                "VALUE_TRANSFER_KLAY_SUCCEEDED_TX",
	VALUE_TRANSFER_KLAY_INVALID_VOTE_OR_NONCE:       "VALUE_TRANSFER_KLAY_INVALID_VOTE_OR_NONCE",
	VALUE_TRANSFER_KLAY_UNEXECUTED:                  "VALUE_TRANSFER_KLAY_UNEXECUTED",
	VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE: "VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE",
	VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS: "VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS",
}

type RequestKLAYReasoningVT struct {
	BridgeAddr              common.Address
	CounterpartBridgeAddr   common.Address
	CounterpartOperatorAddr common.Address
	Value                   *big.Int // Amount of requested value
	IsChildSent             bool
	ToAddr                  common.Address // `to` address in the event of `reuqestValueTransfer` and `requestValueTransferEncoded`
	HandleTxCost            *big.Int
	RequestTxHash           common.Hash
	HandleTxHash            common.Hash
}

// The same struct with `RequestKLAYReasoningVT`
type RequestKLAYReasoningVTForRLP struct {
	BridgeAddr              common.Address
	CounterpartBridgeAddr   common.Address
	CounterpartOperatorAddr common.Address
	Value                   *big.Int
	IsChildSent             bool
	ToAddr                  common.Address
	HandleTxCost            *big.Int
	RequestTxHash           common.Hash
	HandleTxHash            common.Hash
}

type ResponseKLAYReasoningVT struct {
	BridgeAddr           common.Address
	OperatorAddr         common.Address
	Value                *big.Int // Amount of requested value
	IsChildSent          bool
	ContractBalance      *big.Int // The bridge contract balance in the counterpart chain
	OperatorBalance      *big.Int
	ToAddrIsContractAddr bool
	RequestTxHash        common.Hash
	HandleTxHash         common.Hash
	RevertedAddr         common.Address // Wherefrom the valeu transfer is reverted
	RevertedMsg          string
	TraceError           string
	StateDBError         string
	TxCost               *big.Int
	Reasoning            ValueTransferException
}

// The same struct with `ResponseKLAYReasoningVT`
type ResponseKLAYReasoningVTForRLP struct {
	BridgeAddr           common.Address
	OperatorAddr         common.Address
	Value                *big.Int
	IsChildSent          bool
	ContractBalance      *big.Int
	OperatorBalance      *big.Int
	ToAddrIsContractAddr bool
	RequestTxHash        common.Hash
	HandleTxHash         common.Hash
	RevertedAddr         common.Address
	RevertedMsg          string
	TraceError           string
	StateDBError         string
	TxCost               *big.Int
	Reasoning            ValueTransferException
}

func (reqKLAYReasoning *RequestKLAYReasoningVT) EncodeRLP(w io.Writer) error {
	enc := NewObj(reqKLAYReasoning, reflect.TypeOf(&RequestKLAYReasoningVTForRLP{}))
	return rlp.Encode(w, &enc)
}

func (reqKLAYReasoning *RequestKLAYReasoningVT) DecodeRLP(s *rlp.Stream) error {
	var encodedReqTxReasoning RequestKLAYReasoningVTForRLP
	if err := s.Decode(&encodedReqTxReasoning); err != nil {
		return err
	}
	newObj := NewObj(&encodedReqTxReasoning, reflect.TypeOf(&RequestKLAYReasoningVT{}))
	*reqKLAYReasoning = newObj.(RequestKLAYReasoningVT)
	return nil
}

func (respKLAYReasoning *ResponseKLAYReasoningVT) EncodeRLP(w io.Writer) error {
	enc := NewObj(respKLAYReasoning, reflect.TypeOf(&ResponseKLAYReasoningVTForRLP{}))
	return rlp.Encode(w, &enc)
}

func (respKLAYReasoning *ResponseKLAYReasoningVT) DecodeRLP(s *rlp.Stream) error {
	var encodedRespTxReasoning ResponseKLAYReasoningVTForRLP
	if err := s.Decode(&encodedRespTxReasoning); err != nil {
		return err
	}
	newObj := NewObj(&encodedRespTxReasoning, reflect.TypeOf(&ResponseKLAYReasoningVT{}))
	*respKLAYReasoning = newObj.(ResponseKLAYReasoningVT)
	return nil
}

func NewRequestKLAYReasoningVT(bridgeAddr, counterpartAddr, counterpartOperatorAddr, toAddr common.Address,
	valueOrTokenId, handleTxCost *big.Int, isChildSent bool,
	requestTxHash, handleTxHash common.Hash) *RequestVTReasoningWrapper {
	return &RequestVTReasoningWrapper{
		TokenType: KLAY,
		KLAYReasoning: &RequestKLAYReasoningVT{
			BridgeAddr:              bridgeAddr,
			CounterpartBridgeAddr:   counterpartAddr,
			CounterpartOperatorAddr: counterpartOperatorAddr,
			Value:                   valueOrTokenId,
			IsChildSent:             isChildSent,
			ToAddr:                  toAddr,
			RequestTxHash:           requestTxHash,
			HandleTxHash:            handleTxHash,
			HandleTxCost:            handleTxCost,
		},
	}
}

func (reqKLAYReasoning *RequestKLAYReasoningVT) reasoning(blockchain *blockchain.BlockChain, debugAPI *cn.PrivateDebugAPI) ResponseVTReasoningWrapper {
	var respKLAYReasoning ResponseKLAYReasoningVT
	respKLAYReasoning.BridgeAddr = reqKLAYReasoning.BridgeAddr
	respKLAYReasoning.OperatorAddr = reqKLAYReasoning.CounterpartOperatorAddr
	respKLAYReasoning.RequestTxHash = reqKLAYReasoning.RequestTxHash
	respKLAYReasoning.HandleTxHash = reqKLAYReasoning.HandleTxHash
	respKLAYReasoning.Value = reqKLAYReasoning.Value
	respKLAYReasoning.IsChildSent = reqKLAYReasoning.IsChildSent
	respKLAYReasoning.TxCost = reqKLAYReasoning.HandleTxCost
	if statedb, err := blockchain.State(); err == nil {
		respKLAYReasoning.ContractBalance = statedb.GetBalance(reqKLAYReasoning.CounterpartBridgeAddr)
		respKLAYReasoning.OperatorBalance = statedb.GetBalance(reqKLAYReasoning.CounterpartOperatorAddr)
		isContractAccount := statedb.IsContractAccount(reqKLAYReasoning.ToAddr)
		respKLAYReasoning.ToAddrIsContractAddr = isContractAccount
	} else {
		respKLAYReasoning.StateDBError = err.Error()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fastCallTracer := "fastCallTracer"
	traceResult, err := debugAPI.TraceTransaction(ctx, reqKLAYReasoning.HandleTxHash, &cn.TraceConfig{Tracer: &fastCallTracer})
	if err == nil {
		revertedAddr, revertedMsg := getRevertedAddrAndMsg(traceResult)
		respKLAYReasoning.RevertedAddr = revertedAddr
		respKLAYReasoning.RevertedMsg = revertedMsg

		// TODO-hyunsooda: Add KLAY `transfer` failure error at EVM execution
		// Check 1. if its tx is not reverted
		if revertedAddr == (common.Address{}) {
			respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_SUCCEEDED_TX
		} else if revertedAddr == reqKLAYReasoning.CounterpartBridgeAddr {
			// Check 2. if the bridge contract's balance is not enough
			if revertedMsg == "" { // Do not compare the contract's balance due to the balance retrive timing
				respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE
			} else if revertedMsg != "" {
				// Check 3. The internal error of bridge contract must have revert message unless its failure is not by insufficient balance error
				respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_INVALID_VOTE_OR_NONCE
			}
		} else {
			// Check 4. If the `to` address is contract account and reverted on its fallback function
			respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS
		}
	} else {
		// `VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED` is set if the trace result is failed,
		respKLAYReasoning.TraceError = err.Error()
		if err == cn.ErrNotFoundTx {
			respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_UNEXECUTED
		} else {
			respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED
		}
	}
	return ResponseVTReasoningWrapper{
		TokenType:     KLAY,
		KLAYReasoning: &respKLAYReasoning,
	}
}

func (respKLAYReasoning *ResponseKLAYReasoningVT) log(reasoningWord string) {
	operatorAddr := respKLAYReasoning.OperatorAddr
	operatorBalance := respKLAYReasoning.OperatorBalance
	txCost := respKLAYReasoning.TxCost
	var parentOrChild string
	if respKLAYReasoning.IsChildSent {
		parentOrChild = "parent"
	} else {
		parentOrChild = "child"
	}

	logger.Error(reasoningWord,
		"chain", parentOrChild,
		"BridgeAddr", respKLAYReasoning.BridgeAddr.Hex(),
		"OperatorAddr", respKLAYReasoning.OperatorAddr.Hex(),
		"Value", respKLAYReasoning.Value,
		"ContractBalance", respKLAYReasoning.ContractBalance,
		"OperatorBalance", operatorBalance,
		"RequestTxHash", respKLAYReasoning.RequestTxHash.Hex(),
		"HandleTxHash", respKLAYReasoning.HandleTxHash.Hex(),
		"RevertedAddr", respKLAYReasoning.RevertedAddr.Hex(),
		"RevertedMsg", respKLAYReasoning.RevertedMsg,
		"StateDBErr", respKLAYReasoning.StateDBError,
		"TraceResultErr", respKLAYReasoning.TraceError,
		"TxCost", txCost,
		"Reasoning", KLAYReasoingResultStringMap[respKLAYReasoning.Reasoning],
	)

	if operatorBalance.Cmp(txCost) < 1 {
		logger.Error("[SC][Reasoning] The bridge operator may not have enough balance to execute value transfer transaction",
			"chain", parentOrChild,
			"operaotrAddr", operatorAddr.Hex(),
			"txCost", txCost,
			"operatorBalance", operatorBalance,
		)
	}
}

func (respKLAYReasoning *ResponseKLAYReasoningVT) decideResend() (ValueTransferException, bool) {
	switch respKLAYReasoning.Reasoning {
	case VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED:
		respKLAYReasoning.log("[SC][Reasoning] Failed to execute trace transaction call. No resend its transaction again")
		return VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED, false
	case VALUE_TRANSFER_KLAY_SUCCEEDED_TX:
		respKLAYReasoning.log("[SC][Reasoning] Succeeded tx. No resend its transaction")
		return VALUE_TRANSFER_KLAY_SUCCEEDED_TX, false
	case VALUE_TRANSFER_KLAY_INVALID_VOTE_OR_NONCE:
		respKLAYReasoning.log("[SC][Reasoning] Failed by some of internal errors. Expected to never happen this error. No resend its transaction")
		return VALUE_TRANSFER_KLAY_INVALID_VOTE_OR_NONCE, false
	case VALUE_TRANSFER_KLAY_UNEXECUTED:
		respKLAYReasoning.log("[SC][Reasoning] Tx was not found. Send its transaction again")
		return VALUE_TRANSFER_KLAY_UNEXECUTED, true
	case VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE:
		respKLAYReasoning.log("[SC][Reasoning] Failed by not enough bridge contract balance. No resend its transaction")
		return VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE, false
	case VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS:
		respKLAYReasoning.log("[SC][Reasoning] The contract execution was reverted in the other contract address. No resend its transaction")
		return VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS, false
	}
	return VALUE_TRANSFER_UNREACHABLE, false
}

func (respKLAYReasoning *ResponseKLAYReasoningVT) hardening(ev IRequestValueTransferEvent) {
	switch reqEv := ev.(type) {
	case RequestValueTransferEvent:
		reqEv.To = common.HexToAddress("0x0000000000000000000000000000000000000000")
		reqEv.ValueOrTokenId = common.Big0
	case RequestValueTransferEncodedEvent:
		reqEv.To = common.HexToAddress("0x0000000000000000000000000000000000000000")
		reqEv.ValueOrTokenId = common.Big0
	}
	logger.Warn("[SC][Reasoning] Make hardening against handle value transfer tx.",
		"Reasoning", KLAYReasoingResultStringMap[respKLAYReasoning.Reasoning],
		"RequestTxHash", respKLAYReasoning.RequestTxHash.Hex(),
		"HandleTxHash", respKLAYReasoning.HandleTxHash.Hex(),
	)
}

func NewObj(s interface{}, targetType reflect.Type) interface{} {
	oldObj := reflect.ValueOf(s).Elem()
	newObj := reflect.New(targetType.Elem()).Elem()
	for i := 0; i < newObj.NumField(); i++ {
		fieldName := newObj.Type().Field(i).Name
		fieldVal := reflect.Indirect(newObj).FieldByName(fieldName)
		oldVal := reflect.Indirect(oldObj).FieldByName(fieldName)
		fieldVal.Set(reflect.ValueOf(oldVal.Interface()))
	}
	return newObj.Interface()
}
