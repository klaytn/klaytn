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
	logger.Debug("bridgeNode block number", "number", block.Number())
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

func (cce *ChildChainEventHandler) ProcessRequestEvent(ev RequestValueTransferEvent) error {
	handleBridgeAddr := cce.subbridge.bridgeManager.GetCounterPartBridgeAddr(ev.ContractAddr)
	if handleBridgeAddr == (common.Address{}) {
		return fmt.Errorf("there is no counter part bridge of the bridge(%v)", ev.ContractAddr.String())
	}

	handleBridgeInfo, ok := cce.subbridge.bridgeManager.GetBridgeInfo(handleBridgeAddr)
	if !ok {
		return fmt.Errorf("there is no counter part bridge info(%v) of the bridge(%v)", handleBridgeAddr.String(), ev.ContractAddr.String())
	}

	// TODO-Klaytn need to manage the size limitation of pending event list.
	handleBridgeInfo.AddRequestValueTransferEvents([]*RequestValueTransferEvent{&ev})
	return nil
}

func (cce *ChildChainEventHandler) ProcessHandleEvent(ev HandleValueTransferEvent) error {
	handleBridgeInfo, ok := cce.subbridge.bridgeManager.GetBridgeInfo(ev.ContractAddr)
	if !ok {
		return errors.New("there is no bridge")
	}

	handleBridgeInfo.UpdateHandledNonce(ev.HandleNonce + 1)

	tokenType := ev.TokenType
	tokenAddr := ev.TokenAddr

	logger.Trace("RequestValueTransfer Event", "bridgeAddr", ev.ContractAddr.String(), "handleNonce", ev.HandleNonce, "to", ev.Owner.String(), "valueType", tokenType, "token/NFT contract", tokenAddr, "value", ev.Amount.String())
	return nil
}

// ConvertServiceChainBlockHashToMainChainTxHash returns a transaction hash of a transaction which contains
// ChainHashes, with the key made with given service chain block hash.
// Index is built when service chain indexing is enabled.
func (cce *ChildChainEventHandler) ConvertServiceChainBlockHashToMainChainTxHash(scBlockHash common.Hash) common.Hash {
	return cce.subbridge.chainDB.ConvertServiceChainBlockHashToMainChainTxHash(scBlockHash)
}
