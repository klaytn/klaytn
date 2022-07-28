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
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/rlp"
)

const (
	VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED ValueTransferException = iota
	VALUE_TRANSFER_KLAY_SUCCEEDED_TX
	VALUE_TRANSFER_KLAY_OUT_OF_GAS
	VALUE_TRANSFER_KLAY_VM_ERR
	VALUE_TRANSFER_KLAY_REVERTED
	VALUE_TRANSFER_KLAY_UNEXECUTED // (1) if the operator does not have KLAY to execute its handle value transfer tx, or (2) Bridge node is not servicing (connection down or synching)
	VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE
	VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS // means potential attacker sent KLAY value transfer
)

var KLAYReasoingResultStringMap = map[ValueTransferException]string{
	VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED:     "VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED",
	VALUE_TRANSFER_KLAY_SUCCEEDED_TX:                "VALUE_TRANSFER_KLAY_SUCCEEDED_TX",
	VALUE_TRANSFER_KLAY_OUT_OF_GAS:                  "VALUE_TRANSFER_KLAY_OUT_OF_GAS",
	VALUE_TRANSFER_KLAY_VM_ERR:                      "VALUE_TRANSFER_KLAY_VM_ERR",
	VALUE_TRANSFER_KLAY_REVERTED:                    "VALUE_TRANSFER_KLAY_REVERTED",
	VALUE_TRANSFER_KLAY_UNEXECUTED:                  "VALUE_TRANSFER_KLAY_UNEXECUTED",
	VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE: "VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE",
	VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS: "VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS",
}

type RequestKLAYReasoningVT struct {
	HandleTxCost            *big.Int
	Value                   *big.Int
	RequestTxHash           common.Hash
	HandleTxHash            common.Hash
	ToAddr                  common.Address
	BridgeAddr              common.Address
	CounterpartBridgeAddr   common.Address
	CounterpartOperatorAddr common.Address
	IsChildSent             bool
}

// The same struct with `RequestKLAYReasoningVT`
type RequestKLAYReasoningVTForRLP struct {
	HandleTxCost            *big.Int
	Value                   *big.Int
	RequestTxHash           common.Hash
	HandleTxHash            common.Hash
	ToAddr                  common.Address
	BridgeAddr              common.Address
	CounterpartBridgeAddr   common.Address
	CounterpartOperatorAddr common.Address
	IsChildSent             bool
}

type ResponseKLAYReasoningVT struct {
	Value                *big.Int
	TxCost               *big.Int
	ContractBalance      *big.Int
	OperatorBalance      *big.Int
	StateDBError         string
	TraceError           string
	RevertedMsg          string
	Error                string
	GasLimit             uint64
	GasUsed              uint64
	RequestTxHash        common.Hash
	HandleTxHash         common.Hash
	BridgeAddr           common.Address
	RevertedAddr         common.Address
	OperatorAddr         common.Address
	ToAddrIsContractAddr bool
	IsChildSent          bool
	Reasoning            ValueTransferException
}

// The same struct with `ResponseKLAYReasoningVT`
type ResponseKLAYReasoningVTForRLP struct {
	Value                *big.Int
	TxCost               *big.Int
	ContractBalance      *big.Int
	OperatorBalance      *big.Int
	StateDBError         string
	TraceError           string
	RevertedMsg          string
	Error                string
	GasLimit             uint64
	GasUsed              uint64
	RequestTxHash        common.Hash
	HandleTxHash         common.Hash
	BridgeAddr           common.Address
	RevertedAddr         common.Address
	OperatorAddr         common.Address
	ToAddrIsContractAddr bool
	IsChildSent          bool
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
	requestTxHash, handleTxHash common.Hash,
) *RequestVTReasoningWrapper {
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

		err, gasLimit, gasUsed := getErrorWithGasUsed(traceResult)
		respKLAYReasoning.Error = err
		respKLAYReasoning.GasLimit = gasLimit
		respKLAYReasoning.GasUsed = gasUsed

		// TODO-hyunsooda: Add KLAY `transfer` failure error at EVM execution
		// Check 1. if its tx is not reverted and failed by out of gas
		if revertedAddr == (common.Address{}) {
			if err != "" {
				if err == kerrors.ErrOutOfGas.Error() {
					respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_OUT_OF_GAS
				} else {
					respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_VM_ERR
				}
			} else {
				respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_SUCCEEDED_TX
			}
		} else if revertedAddr == reqKLAYReasoning.CounterpartBridgeAddr {
			// Check 2. if the bridge contract's balance is not enough
			if revertedMsg == "" { // Do not compare the contract's balance due to the balance retrive timing
				respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE
			} else if revertedMsg != "" {
				// Check 3. The internal error of bridge contract must have revert message unless its failure is not by insufficient balance error
				respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_REVERTED
			}
		} else {
			// Check 4. If the `to` address is contract account and reverted on its fallback function
			respKLAYReasoning.Reasoning = VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS
		}
	} else {
		// `VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED` is set if the trace result is failed,
		respKLAYReasoning.TraceError = err.Error()
		// if err == cn.ErrNotFoundTx {
		if err.Error() == cn.MakeNotFoundTxErr(reqKLAYReasoning.HandleTxHash).Error() {
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
	var parentOrChild string
	if respKLAYReasoning.IsChildSent {
		parentOrChild = "parent"
	} else {
		parentOrChild = "child"
	}
	var cause string
	if respKLAYReasoning.OperatorBalance.Cmp(respKLAYReasoning.Value) < 1 {
		cause = "The bridge operator may not have enough balance to execute value transfer transaction"
	}
	if respKLAYReasoning.Reasoning == VALUE_TRANSFER_KLAY_OUT_OF_GAS {
		cause = "Please validate the operator's klay balance and gaslmiit value (`subbridge_getParentBridgeOperatorGasLimit` and `subbridge_getChildBridgeOperatorGasLimit`)"
	}

	logger.Error(reasoningWord,
		"chain", parentOrChild,
		"BridgeAddr", respKLAYReasoning.BridgeAddr.Hex(),
		"OperatorAddr", respKLAYReasoning.OperatorAddr.Hex(),
		"Value", respKLAYReasoning.Value,
		"ContractBalance", respKLAYReasoning.ContractBalance,
		"OperatorBalance", respKLAYReasoning.OperatorBalance,
		"RequestTxHash", respKLAYReasoning.RequestTxHash.Hex(),
		"HandleTxHash", respKLAYReasoning.HandleTxHash.Hex(),
		"RevertedAddr", respKLAYReasoning.RevertedAddr.Hex(),
		"RevertedMsg", respKLAYReasoning.RevertedMsg,
		"Error", respKLAYReasoning.Error,
		"GasLimit", respKLAYReasoning.GasLimit,
		"GasUsed", respKLAYReasoning.GasUsed,
		"StateDBErr", respKLAYReasoning.StateDBError,
		"TraceResultErr", respKLAYReasoning.TraceError,
		"TxCost", respKLAYReasoning.TxCost,
		"Reasoning", KLAYReasoingResultStringMap[respKLAYReasoning.Reasoning],
		"Cause", cause,
	)
}

func (respKLAYReasoning *ResponseKLAYReasoningVT) decideResend() bool {
	switch respKLAYReasoning.Reasoning {
	case VALUE_TRANSFER_KLAY_TRACETRANSACTION_FAILED:
		respKLAYReasoning.log("[SC][Reasoning] Failed to execute trace transaction call. No resend its transaction again")
		return false
	case VALUE_TRANSFER_KLAY_SUCCEEDED_TX:
		respKLAYReasoning.log("[SC][Reasoning] Succeeded tx. No resend its transaction")
		return false
	case VALUE_TRANSFER_KLAY_OUT_OF_GAS:
		respKLAYReasoning.log("[SC][Reasoning] Failed by out of gas. No resend its transaction")
		return false
	case VALUE_TRANSFER_KLAY_VM_ERR:
		respKLAYReasoning.log("[SC][Reasoning] Failed by vm execution error. No resend its transaction")
		return false
	case VALUE_TRANSFER_KLAY_REVERTED:
		respKLAYReasoning.log("[SC][Reasoning] Failed by revert. Expected to never happen this error. No resend its transaction")
		return false
	case VALUE_TRANSFER_KLAY_UNEXECUTED:
		respKLAYReasoning.log("[SC][Reasoning] Tx was not found. Try to Resend it if its reasoning timeout exceeded")
		return true
	case VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE:
		respKLAYReasoning.log("[SC][Reasoning] Failed by not enough bridge contract balance. No resend its transaction.")
		return false
	case VALUE_TRANSFER_KLAY_REVERT_ON_THE_OTHER_ADDRESS:
		respKLAYReasoning.log("[SC][Reasoning] The contract execution was reverted in the other contract address. No resend its transaction")
		return false
	}
	return false
}

func (respKLAYReasoning *ResponseKLAYReasoningVT) refundable() bool {
	switch respKLAYReasoning.Reasoning {
	case VALUE_TRANSFER_KLAY_OUT_OF_GAS, VALUE_TRANSFER_KLAY_NOT_ENOUGH_CONTRACT_BALANCE:
		return true
	default:
		return false
	}
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
