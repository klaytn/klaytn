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
	"io"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/node/cn"
	"github.com/klaytn/klaytn/rlp"
)

type ValueTransferException uint8

type RequestVTReasoningWrapper struct {
	TokenType     uint8
	KLAYReasoning *RequestKLAYReasoningVT
	// ERC20Reasoning  *RequestERC20ReasoningVT
	// ERC721Reasoning *RequestERC721ReasoningVt
}

type ResponseVTReasoningWrapper struct {
	TokenType     uint8
	KLAYReasoning *ResponseKLAYReasoningVT
	// ERC20Reasoning  *ResponseERC20ReasoningVT
	// ERC721Reasoning *ResponseERC721ReasoningVT
}

// The same struct with `RequestTxReasoning`
type RequestVTReasoningWrapperForRLP struct {
	TokenType     uint8
	KLAYReasoning *RequestKLAYReasoningVT
	// ERC20Reasoning  *RequestERC20ReasoningVT
	// ERC721Reasoning *RequestERC721ReasoningVt
}

// The same struct with `ResponseTxReasoning`
type ResponseVTReasoningWrapperForRLP struct {
	TokenType     uint8
	KLAYReasoning *ResponseKLAYReasoningVT
	// ERC20Reasoning  *ResponseERC20ReasoningVT
	// ERC721Reasoning *ResponseERC721ReasoningVT
}

func (reqVTReasoningWrapper *RequestVTReasoningWrapper) EncodeRLP(w io.Writer) error {
	enc := &RequestVTReasoningWrapperForRLP{
		TokenType:     reqVTReasoningWrapper.TokenType,
		KLAYReasoning: reqVTReasoningWrapper.KLAYReasoning,
	}
	return rlp.Encode(w, enc)
}

func (respVTReasoningWrapper *ResponseVTReasoningWrapper) DecodeRLP(s *rlp.Stream) error {
	var encodedRespWrapper ResponseVTReasoningWrapperForRLP
	if err := s.Decode(&encodedRespWrapper); err != nil {
		return err
	}
	respVTReasoningWrapper.TokenType = encodedRespWrapper.TokenType
	respVTReasoningWrapper.KLAYReasoning = encodedRespWrapper.KLAYReasoning
	return nil
}

func (reqVTReasoningWrapper *RequestVTReasoningWrapper) reasoning(blockchain *blockchain.BlockChain, debugAPI *cn.PrivateDebugAPI) ResponseVTReasoningWrapper {
	switch reqVTReasoningWrapper.TokenType {
	case KLAY:
		return reqVTReasoningWrapper.KLAYReasoning.reasoning(blockchain, debugAPI)
		/*
			// TODO-hyunsooda: Add reasoning implementation for ERC20 and ERC721 tokens
			case ERC20:
			case ERC721:
		*/
	default:
		return ResponseVTReasoningWrapper{}
	}
}

func (respVTReasoningWrapper *ResponseVTReasoningWrapper) decideResend() bool {
	switch respVTReasoningWrapper.TokenType {
	case KLAY:
		return respVTReasoningWrapper.KLAYReasoning.decideResend()
	default:
		return false
	}
}

func (respVTReasoningWrapper *ResponseVTReasoningWrapper) refundable() bool {
	switch respVTReasoningWrapper.TokenType {
	case KLAY:
		return respVTReasoningWrapper.KLAYReasoning.refunadable()
	default:
		return false
	}
}

/*
func (respVTReasoningWrapper *ResponseVTReasoningWrapper) hardening(ev IRequestValueTransferEvent) {
	switch respVTReasoningWrapper.TokenType {
	case KLAY:
		respVTReasoningWrapper.KLAYReasoning.hardening(ev)
	}
}
*/

func (respVTReasoningWrapper *ResponseVTReasoningWrapper) getBridgeAddr() common.Address {
	switch respVTReasoningWrapper.TokenType {
	case KLAY:
		return respVTReasoningWrapper.KLAYReasoning.BridgeAddr
	default:
		return common.Address{}
	}
}

func getRevertedAddrAndMsg(traceResult interface{}) (common.Address, string) {
	if tr, ok := traceResult.(*vm.InternalTxTrace); ok {
		var revertedAddr common.Address
		if tr.Reverted != nil {
			revertedMsg := tr.Reverted.Message
			if tr.Reverted.Contract != nil {
				revertedAddr = *tr.Reverted.Contract
			} else {
				return common.Address{}, revertedMsg
			}
			for _, call := range tr.Calls {
				if call.To != nil {
					if revertedAddr == *call.To { // Tx was reverted in the `to` contract address
						return *call.To, revertedMsg
					}
				}
			}
			return revertedAddr, revertedMsg
		}
	}
	return common.Address{}, ""
}

func getErrorWithGasUsed(traceResult interface{}) (string, uint64, uint64) {
	if tr, ok := traceResult.(*vm.InternalTxTrace); ok {
		if tr.Error != nil {
			return tr.Error.Error(), tr.Gas, tr.GasUsed
		}
	}
	return "", 0, 0
}

func makeRequestKLAYHandleDebug(bi *BridgeInfo, ev IRequestValueTransferEvent) *RequestVTReasoningWrapper {
	isOnChild := bi.onChildChain
	reqTxHash := ev.GetRaw().TxHash
	handleTx := bi.bridgeDB.ReadAllHandleInfo(bi.address, bi.counterpartAddress, reqTxHash).HandleTx
	var counterpartOperatorAddr common.Address
	if isOnChild {
		counterpartOperatorAddr = bi.subBridge.APIBackend.GetParentOperatorAddr()
	} else {
		counterpartOperatorAddr = bi.subBridge.APIBackend.GetChildOperatorAddr()
	}
	return NewRequestKLAYReasoningVT(
		bi.address,
		bi.counterpartAddress,
		counterpartOperatorAddr,
		ev.GetTo(),
		ev.GetValueOrTokenId(),
		handleTx.Cost(),
		isOnChild,
		reqTxHash,
		handleTx.Hash(),
	)
}
