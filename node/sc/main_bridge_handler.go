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
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/pkg/errors"
	"math/big"
)

var (
	ErrNoChildChainID = errors.New("There is no childChainID")
)

type MainBridgeHandler struct {
	mainbridge *MainBridge
	// childChainID is the first received chainID from child chain peer.
	childChainIDs map[common.Address]*big.Int
}

func NewMainBridgeHandler(scc *SCConfig, main *MainBridge) (*MainBridgeHandler, error) {
	return &MainBridgeHandler{
		mainbridge:    main,
		childChainIDs: make(map[common.Address]*big.Int),
	}, nil
}

func (mbh *MainBridgeHandler) HandleSubMsg(p BridgePeer, msg p2p.Msg) error {

	// Handle the message depending on its contents
	switch msg.Code {
	case StatusMsg:
		return nil
	case ServiceChainTxsMsg:
		logger.Debug("received ServiceChainTxsMsg")
		// TODO-Klaytn how to check acceptTxs
		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		//if atomic.LoadUint32(&pm.acceptTxs) == 0 {
		//	break
		//}
		if err := mbh.handleServiceChainTxDataMsg(p, msg); err != nil {
			return err
		}
	case ServiceChainParentChainInfoRequestMsg:
		logger.Debug("received ServiceChainParentChainInfoRequestMsg")
		if err := mbh.handleServiceChainParentChainInfoRequestMsg(p, msg); err != nil {
			return err
		}
	case ServiceChainReceiptRequestMsg:
		logger.Debug("received ServiceChainReceiptRequestMsg")
		if err := mbh.handleServiceChainReceiptRequestMsg(p, msg); err != nil {
			return err
		}
	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}
	return nil
}

// handleServiceChainTxDataMsg handles service chain transactions from child chain.
// It will return an error if given tx is not TxTypeChainDataAnchoring type.
func (mbh *MainBridgeHandler) handleServiceChainTxDataMsg(p BridgePeer, msg p2p.Msg) error {
	//pm.txMsgLock.Lock()
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs []*types.Transaction
	if err := msg.Decode(&txs); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}

	// Only valid txs should be pushed into the pool.
	validTxs := make([]*types.Transaction, 0, len(txs))
	//validTxs := []*types.Transaction{}
	var err error
	for i, tx := range txs {
		if tx == nil {
			err = errResp(ErrDecode, "tx %d is nil", i)
			continue
		}
		validTxs = append(validTxs, tx)
	}
	mbh.mainbridge.txPool.AddRemotes(validTxs)
	return err
}

// handleServiceChainParentChainInfoRequestMsg handles parent chain info request message from child chain.
// It will send the nonce of the account and its gas price to the child chain peer who requested.
func (mbh *MainBridgeHandler) handleServiceChainParentChainInfoRequestMsg(p BridgePeer, msg p2p.Msg) error {
	var addr common.Address
	if err := msg.Decode(&addr); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	nonce := mbh.mainbridge.txPool.GetPendingNonce(addr)
	pcInfo := parentChainInfo{nonce, mbh.mainbridge.blockchain.Config().UnitPrice}
	p.SendServiceChainInfoResponse(&pcInfo)
	logger.Info("SendServiceChainInfoResponse", "mcBridgeAccoount", addr, "nonce", pcInfo.Nonce, "gasPrice", pcInfo.GasPrice)
	return nil
}

// handleServiceChainReceiptRequestMsg handles receipt request message from child chain.
// It will find and send corresponding receipts with given transaction hashes.
func (mbh *MainBridgeHandler) handleServiceChainReceiptRequestMsg(p BridgePeer, msg p2p.Msg) error {
	// Decode the retrieval message
	msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
	if _, err := msgStream.List(); err != nil {
		return err
	}
	// Gather state data until the fetch or network limits is reached
	var (
		hash               common.Hash
		receiptsForStorage []*types.ReceiptForStorage
	)
	for len(receiptsForStorage) < downloader.MaxReceiptFetch {
		// Retrieve the hash of the next block
		if err := msgStream.Decode(&hash); err == rlp.EOL {
			break
		} else if err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Retrieve the receipt of requested service chain tx, skip if unknown.
		receipt := mbh.mainbridge.blockchain.GetReceiptByTxHash(hash)
		if receipt == nil {
			continue
		}

		receiptsForStorage = append(receiptsForStorage, (*types.ReceiptForStorage)(receipt))
	}
	if len(receiptsForStorage) == 0 {
		return nil
	}
	return p.SendServiceChainReceiptResponse(receiptsForStorage)
}
