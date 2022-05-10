// Copyright 2019 The klaytn Authors
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
	"errors"
	"fmt"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

var (
	ErrGetServiceChainPHInCCEH = errors.New("ServiceChainPH isn't set in ChildChainEventHandler")
)

type ChildChainEventHandler struct {
	subbridge *SubBridge

	handler *SubBridgeHandler
}

const (
	// TODO-Klaytn need to define proper value.
	errorDiffRequestHandleNonce = 10000
)

func NewChildChainEventHandler(bridge *SubBridge, handler *SubBridgeHandler) (*ChildChainEventHandler, error) {
	return &ChildChainEventHandler{subbridge: bridge, handler: handler}, nil
}

func (cce *ChildChainEventHandler) HandleChainHeadEvent(block *types.Block) error {
	logger.Trace("bridgeNode block number", "number", block.Number())
	cce.handler.LocalChainHeadEvent(block)

	// Logging information of value transfer
	cce.subbridge.bridgeManager.LogBridgeStatus()

	return nil
}

func (cce *ChildChainEventHandler) HandleTxEvent(tx *types.Transaction) error {
	//TODO-Klaytn event handle
	return nil
}

func (cce *ChildChainEventHandler) HandleTxsEvent(txs []*types.Transaction) error {
	//TODO-Klaytn event handle
	return nil
}

func (cce *ChildChainEventHandler) HandleLogsEvent(logs []*types.Log) error {
	//TODO-Klaytn event handle
	return nil
}

func (cce *ChildChainEventHandler) ProcessRequestEvent(ev IRequestValueTransferEvent) error {
	addr := ev.GetRaw().Address

	handleBridgeAddr := cce.subbridge.bridgeManager.GetCounterPartBridgeAddr(addr)
	if handleBridgeAddr == (common.Address{}) {
		return fmt.Errorf("there is no counter part bridge of the bridge(%v)", addr.String())
	}

	handleBridgeInfo, ok := cce.subbridge.bridgeManager.GetBridgeInfo(handleBridgeAddr)
	if !ok {
		return fmt.Errorf("there is no counter part bridge info(%v) of the bridge(%v)", handleBridgeAddr.String(), addr.String())
	}

	// TODO-Klaytn need to manage the size limitation of pending event list.
	handleBridgeInfo.AddRequestValueTransferEvents([]IRequestValueTransferEvent{ev})
	return nil
}

func (cce *ChildChainEventHandler) ProcessHandleEvent(ev *HandleValueTransferEvent) error {
	handleBridgeInfo, ok := cce.subbridge.bridgeManager.GetBridgeInfo(ev.Raw.Address)
	if !ok {
		return errors.New("there is no bridge")
	}

	handleBridgeInfo.MarkHandledNonce(ev.HandleNonce)
	handleBridgeInfo.UpdateLowerHandleNonce(ev.LowerHandleNonce)

	logger.Trace("RequestValueTransfer Event",
		"bridgeAddr", ev.Raw.Address.String(),
		"handleNonce", ev.HandleNonce,
		"to", ev.To.String(),
		"valueType", ev.TokenType,
		"token/NFT contract", ev.TokenAddress,
		"value", ev.ValueOrTokenId.String())
	return nil
}

// ConvertChildChainBlockHashToParentChainTxHash returns a transaction hash of a transaction which contains
// AnchoringData, with the key made with given child chain block hash.
// Index is built when child chain indexing is enabled.
func (cce *ChildChainEventHandler) ConvertChildChainBlockHashToParentChainTxHash(scBlockHash common.Hash) common.Hash {
	return cce.subbridge.chainDB.ConvertChildChainBlockHashToParentChainTxHash(scBlockHash)
}
