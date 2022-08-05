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
	"fmt"
	"io"
	"reflect"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
)

func newRequestReceiptHandle(bridgeAddr, counterpartBridgeAddr common.Address, handleTxHash common.Hash, ev IRequestValueTransferEvent) *RequestHandleReceipt {
	return &RequestHandleReceipt{
		RequestNonce:          ev.GetRequestNonce(),
		RequestTxHash:         ev.GetRaw().TxHash,
		BridgeAddr:            bridgeAddr,
		CounterpartBridgeAddr: counterpartBridgeAddr,
		HandleTxHash:          handleTxHash,
	}
}

func getHandleTx(bi *BridgeInfo, ev IRequestValueTransferEvent) (*types.Transaction, error) {
	reqTxHash := ev.GetRaw().TxHash
	handleInfo := bi.bridgeDB.ReadAllHandleInfo(bi.address, bi.counterpartAddress, reqTxHash)
	if handleInfo != nil {
		return handleInfo.HandleTx, nil
	}
	return nil, fmt.Errorf("No handle info was found (reqNonce = %d, handleTx = %x)", ev.GetRequestNonce(), reqTxHash)
}

type RequestHandleReceipt struct {
	RequestNonce          uint64
	BridgeAddr            common.Address
	CounterpartBridgeAddr common.Address
	RequestTxHash         common.Hash
	HandleTxHash          common.Hash
}

// The same struct with `RequestHandleReceipt`
type RequestHandleReceiptForRLP struct {
	RequestNonce          uint64
	BridgeAddr            common.Address
	CounterpartBridgeAddr common.Address
	RequestTxHash         common.Hash
	HandleTxHash          common.Hash
}

func (reqHandleReceipt *RequestHandleReceipt) EncodeRLP(w io.Writer) error {
	enc := NewObj(reqHandleReceipt, reflect.TypeOf(&RequestHandleReceiptForRLP{}))
	return rlp.Encode(w, &enc)
}

func (reqHandleReceipt *RequestHandleReceipt) DecodeRLP(s *rlp.Stream) error {
	var decoded RequestHandleReceiptForRLP
	if err := s.Decode(&decoded); err != nil {
		return err
	}
	newObj := NewObj(&decoded, reflect.TypeOf(&RequestHandleReceipt{}))
	*reqHandleReceipt = newObj.(RequestHandleReceipt)
	return nil
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

func getReceipt(blockchain *blockchain.BlockChain, handleTxHash common.Hash) *types.Receipt {
	return blockchain.GetReceiptByTxHash(handleTxHash)
}
