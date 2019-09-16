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
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/ser/rlp"
	"math/big"
)

const (
	SyncRequestInterval = 10
)

var (
	errUnknownAnchoringTxType = errors.New("unknown anchoring tx type")
)

// parentChainInfo handles the information of parent chain, which is needed from child chain.
type parentChainInfo struct {
	Nonce    uint64
	GasPrice uint64
}

type SubBridgeHandler struct {
	subbridge *SubBridge
	// parentChainID is the first received chainID from parent chain peer.
	// It will be reset to nil if there's no parent peer.
	parentChainID *big.Int
	// remoteGasPrice means gas price of parent chain, used to make a service chain transaction.
	// Therefore, for now, it is only used by child chain side.
	remoteGasPrice        uint64
	mainChainAccountNonce uint64
	nonceSynced           bool
	chainTxPeriod         uint64

	latestTxCountAddedBlockNumber uint64
	txCountEnabledBlockNumber     uint64
	txCount                       uint64 // accumulated tx counts in blocks for each anchoring period.

	// TODO-Klaytn-ServiceChain Need to limit the number independently? Or just managing the size of sentServiceChainTxs?
	sentServiceChainTxsLimit uint64

	skipSyncBlockCount int32
}

func NewSubBridgeHandler(main *SubBridge) (*SubBridgeHandler, error) {
	return &SubBridgeHandler{
		subbridge:                     main,
		parentChainID:                 new(big.Int).SetUint64(main.config.ParentChainID),
		remoteGasPrice:                uint64(0),
		mainChainAccountNonce:         uint64(0),
		nonceSynced:                   false,
		chainTxPeriod:                 main.config.AnchoringPeriod,
		latestTxCountAddedBlockNumber: uint64(0),
		sentServiceChainTxsLimit:      main.config.SentChainTxsLimit,
	}, nil
}

func (sbh *SubBridgeHandler) setParentChainID(chainId *big.Int) {
	sbh.parentChainID = chainId
	sbh.subbridge.bridgeAccounts.pAccount.SetChainID(chainId)
}

func (sbh *SubBridgeHandler) getParentChainID() *big.Int {
	return sbh.parentChainID
}

func (sbh *SubBridgeHandler) LockParentOperator() {
	sbh.subbridge.bridgeAccounts.pAccount.Lock()
}

func (sbh *SubBridgeHandler) UnLockParentOperator() {
	sbh.subbridge.bridgeAccounts.pAccount.UnLock()
}

// getParentOperatorNonce returns the parent chain operator nonce of parent chain operator address.
func (sbh *SubBridgeHandler) getParentOperatorNonce() uint64 {
	return sbh.subbridge.bridgeAccounts.pAccount.GetNonce()
}

// setParentOperatorNonce sets the parent chain operator nonce of parent chain operator address.
func (sbh *SubBridgeHandler) setParentOperatorNonce(newNonce uint64) {
	sbh.subbridge.bridgeAccounts.pAccount.SetNonce(newNonce)
}

// addParentOperatorNonce increases nonce by number
func (sbh *SubBridgeHandler) addParentOperatorNonce(number uint64) {
	sbh.subbridge.bridgeAccounts.pAccount.IncNonce()
}

// getParentOperatorNonceSynced returns whether the parent chain operator account nonce is synced or not.
func (sbh *SubBridgeHandler) getParentOperatorNonceSynced() bool {
	return sbh.nonceSynced
}

// setParentOperatorNonceSynced sets whether the parent chain operator account nonce is synced or not.
func (sbh *SubBridgeHandler) setParentOperatorNonceSynced(synced bool) {
	sbh.nonceSynced = synced
}

func (sbh *SubBridgeHandler) getChildOperatorNonce() uint64 {
	return sbh.subbridge.txPool.GetPendingNonce(sbh.subbridge.bridgeAccounts.cAccount.address)
}

func (sbh *SubBridgeHandler) getRemoteGasPrice() uint64 {
	return sbh.remoteGasPrice
}

func (sbh *SubBridgeHandler) setRemoteGasPrice(gasPrice uint64) {
	sbh.subbridge.bridgeAccounts.pAccount.SetGasPrice(big.NewInt(int64(gasPrice)))
	sbh.remoteGasPrice = gasPrice
}

// GetParentOperatorAddr returns a pointer of a hex address of an account used for parent chain.
// If given as a parameter, it will use it. If not given, it will use the address of the public key
// derived from chainKey.
func (sbh *SubBridgeHandler) GetParentOperatorAddr() *common.Address {
	return &sbh.subbridge.bridgeAccounts.pAccount.address
}

// GetChildOperatorAddr returns a pointer of a hex address of an account used for child chain.
// If given as a parameter, it will use it. If not given, it will use the address of the public key
// derived from chainKey.
func (sbh *SubBridgeHandler) GetChildOperatorAddr() *common.Address {
	return &sbh.subbridge.bridgeAccounts.cAccount.address
}

// GetAnchoringPeriod returns the period to make and send a chain transaction to parent chain.
func (sbh *SubBridgeHandler) GetAnchoringPeriod() uint64 {
	return sbh.chainTxPeriod
}

// GetSentChainTxsLimit returns the maximum number of stored chain transactions for resending.
func (sbh *SubBridgeHandler) GetSentChainTxsLimit() uint64 {
	return sbh.sentServiceChainTxsLimit
}

func (sbh *SubBridgeHandler) HandleMainMsg(p BridgePeer, msg p2p.Msg) error {
	// Handle the message depending on its contents
	switch msg.Code {
	case ServiceChainResponse:
		logger.Trace("received rpc ServiceChainResponse")
		data := make([]byte, msg.Size)
		err := msg.Decode(&data)
		if err != nil {
			logger.Error("failed to decode the p2p ServiceChainResponse message", "err", err)
			return nil
		}
		logger.Trace("send rpc response to the rpc client")
		_, err = sbh.subbridge.rpcConn.Write(data)
		if err != nil {
			return err
		}
		return nil
	case StatusMsg:
		return nil
	case ServiceChainParentChainInfoResponseMsg:
		logger.Debug("received ServiceChainParentChainInfoResponseMsg")
		if err := sbh.handleParentChainInfoResponseMsg(p, msg); err != nil {
			return err
		}

	case ServiceChainReceiptResponseMsg:
		logger.Debug("received ServiceChainReceiptResponseMsg")
		if err := sbh.handleParentChainReceiptResponseMsg(p, msg); err != nil {
			return err
		}
	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}
	return nil
}

// handleParentChainInfoResponseMsg handles parent chain info response message from parent chain.
// It will update the ParentOperatorNonce and remoteGasPrice of ServiceChainProtocolManager.
func (sbh *SubBridgeHandler) handleParentChainInfoResponseMsg(p BridgePeer, msg p2p.Msg) error {
	var pcInfo parentChainInfo
	if err := msg.Decode(&pcInfo); err != nil {
		logger.Error("failed to decode", "err", err)
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	sbh.LockParentOperator()
	defer sbh.UnLockParentOperator()

	poolNonce := sbh.subbridge.bridgeTxPool.GetMaxTxNonce(sbh.GetParentOperatorAddr())
	if poolNonce > 0 {
		poolNonce += 1
		// just check
		if sbh.getParentOperatorNonce() > poolNonce {
			logger.Error("parent chain operator nonce is bigger than the chain pool nonce.", "BridgeTxPoolNonce", poolNonce, "mainChainAccountNonce", sbh.getParentOperatorNonce())
		}
		if poolNonce < pcInfo.Nonce {
			// bridgeTxPool journal miss txs which already sent to parent-chain
			logger.Error("chain pool nonce is less than the parent chain nonce.", "chainPoolNonce", poolNonce, "receivedNonce", pcInfo.Nonce)
			sbh.setParentOperatorNonce(pcInfo.Nonce)
		} else {
			// bridgeTxPool journal has txs which don't receive receipt from parent-chain
			sbh.setParentOperatorNonce(poolNonce)
		}
	} else if sbh.getParentOperatorNonce() > pcInfo.Nonce {
		logger.Error("parent chain operator nonce is bigger than the received nonce.", "mainChainAccountNonce", sbh.getParentOperatorNonce(), "receivedNonce", pcInfo.Nonce)
		sbh.setParentOperatorNonce(pcInfo.Nonce)
	} else {
		// there is no tx in bridgetTxPool, so parent-chain's nonce is used
		sbh.setParentOperatorNonce(pcInfo.Nonce)
	}
	sbh.setParentOperatorNonceSynced(true)
	sbh.setRemoteGasPrice(pcInfo.GasPrice)
	logger.Info("ParentChainNonceResponse", "receivedNonce", pcInfo.Nonce, "gasPrice", pcInfo.GasPrice, "mainChainAccountNonce", sbh.getParentOperatorNonce())
	return nil
}

// handleParentChainReceiptResponseMsg handles receipt response message from parent chain.
// It will store the received receipts and remove corresponding transaction in the resending list.
func (sbh *SubBridgeHandler) handleParentChainReceiptResponseMsg(p BridgePeer, msg p2p.Msg) error {
	// TODO-Klaytn-ServiceChain Need to add an option, not to write receipts.
	// Decode the retrieval message
	var receipts []*types.ReceiptForStorage
	if err := msg.Decode(&receipts); err != nil && err != rlp.EOL {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// Stores receipt and remove tx from sentServiceChainTxs only if the tx is successfully executed.
	sbh.writeServiceChainTxReceipts(sbh.subbridge.blockchain, receipts)
	return nil
}

// genUnsignedChainDataAnchoringTx generates an unsigned transaction, which type is TxTypeChainDataAnchoring.
// Nonce of account used for service chain transaction will be increased after the signing.
func (sbh *SubBridgeHandler) genUnsignedChainDataAnchoringTx(block *types.Block) (*types.Transaction, error) {
	anchoringData, err := types.NewAnchoringDataType0(block, new(big.Int).SetUint64(sbh.chainTxPeriod), new(big.Int).SetUint64(sbh.txCount))
	if err != nil {
		return nil, err
	}
	encodedCCTxData, err := rlp.EncodeToBytes(anchoringData)
	if err != nil {
		return nil, err
	}

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:        sbh.getParentOperatorNonce(), // parent chain operator nonce will be increased after signing a transaction.
		types.TxValueKeyFrom:         *sbh.GetParentOperatorAddr(),
		types.TxValueKeyGasLimit:     uint64(100000), // TODO-Klaytn-ServiceChain should define proper gas limit
		types.TxValueKeyGasPrice:     new(big.Int).SetUint64(sbh.remoteGasPrice),
		types.TxValueKeyAnchoredData: encodedCCTxData,
	}

	if tx, err := types.NewTransactionWithMap(types.TxTypeChainDataAnchoring, values); err != nil {
		return nil, err
	} else {
		return tx, nil
	}
}

// LocalChainHeadEvent deals with servicechain feature to generate/broadcast service chain transactions and request receipts.
func (sbh *SubBridgeHandler) LocalChainHeadEvent(block *types.Block) {
	if sbh.getParentOperatorNonceSynced() {
		// TODO-Klaytn if other feature use below chainTx, this condition should be refactored to use it for other feature.
		if sbh.subbridge.GetAnchoringTx() {
			sbh.blockAnchoringManager(block)
		}
		sbh.broadcastServiceChainTx()
		sbh.broadcastServiceChainReceiptRequest()

		sbh.skipSyncBlockCount = 0
	} else {
		if sbh.skipSyncBlockCount%SyncRequestInterval == 0 {
			// TODO-Klaytn too many request while sync main-net
			sbh.SyncNonceAndGasPrice()
			// check tx's receipts which parent-chain already executed in bridgeTxPool
			go sbh.broadcastServiceChainReceiptRequest()
		}
		sbh.skipSyncBlockCount++
	}
}

// broadcastServiceChainTx broadcasts service chain transactions to parent chain peers.
// It signs the given unsigned transaction with parent chain ID and then send it to its
// parent chain peers.
func (sbh *SubBridgeHandler) broadcastServiceChainTx() {
	parentChainID := sbh.parentChainID
	if parentChainID == nil {
		logger.Error("unexpected nil parentChainID while broadcastServiceChainTx")
	}
	txs := sbh.subbridge.GetBridgeTxPool().PendingTxsByAddress(&sbh.subbridge.bridgeAccounts.pAccount.address, int(sbh.GetSentChainTxsLimit())) // TODO-Klaytn-Servicechain change GetSentChainTxsLimit type to int from uint64
	peers := sbh.subbridge.BridgePeerSet().peers

	for _, peer := range peers {
		if peer.GetChainID().Cmp(parentChainID) != 0 {
			logger.Error("parent peer with different parent chainID", "peerID", peer.GetID(), "peer chainID", peer.GetChainID(), "parent chainID", parentChainID)
			continue
		}
		peer.SendServiceChainTxs(txs)
		logger.Trace("sent ServiceChainTxData", "peerID", peer.GetID())
	}
	logger.Debug("broadcastServiceChainTx ServiceChainTxData", "len(txs)", len(txs), "len(peers)", len(peers))
}

// decodeAnchoringTx decodes an anchoring transaction.
func (sbh *SubBridgeHandler) decodeAnchoringTx(data []byte) (common.Hash, *big.Int, error) {
	anchoringData := new(types.AnchoringData)
	if err := rlp.DecodeBytes(data, anchoringData); err != nil {
		anchoringDataLegacy := new(types.AnchoringDataLegacy)
		if err := rlp.DecodeBytes(data, anchoringDataLegacy); err != nil {
			return common.Hash{}, nil, err
		}
		logger.Trace("decoded legacy anchoring tx", "blockNum", anchoringDataLegacy.BlockNumber.String(), "blockHash", anchoringDataLegacy.BlockHash.String(), "txHash", anchoringDataLegacy.TxHash.String())
		return anchoringDataLegacy.BlockHash, anchoringDataLegacy.BlockNumber, nil
	}
	if anchoringData.Type == types.AnchoringDataType0 {
		anchoringDataInternal := new(types.AnchoringDataInternalType0)
		if err := rlp.DecodeBytes(anchoringData.Data, anchoringDataInternal); err != nil {
			return common.Hash{}, nil, err
		}
		logger.Trace("decoded type0 anchoring tx", "blockNum", anchoringDataInternal.BlockNumber.String(), "blockHash", anchoringDataInternal.BlockHash.String(), "txHash", anchoringDataInternal.TxHash.String(), "txCount", anchoringDataInternal.TxCount)
		return anchoringDataInternal.BlockHash, anchoringDataInternal.BlockNumber, nil
	} else {
		return common.Hash{}, nil, errUnknownAnchoringTxType
	}
}

// writeServiceChainTxReceipts writes the received receipts of service chain transactions.
func (sbh *SubBridgeHandler) writeServiceChainTxReceipts(bc *blockchain.BlockChain, receipts []*types.ReceiptForStorage) {
	for _, receipt := range receipts {
		txHash := receipt.TxHash
		if tx := sbh.subbridge.GetBridgeTxPool().Get(txHash); tx != nil {
			if tx.Type() == types.TxTypeChainDataAnchoring {
				data, err := tx.AnchoredData()
				if err != nil {
					logger.Error("failed to get anchoring data", "txHash", txHash.String(), "err", err)
					continue
				}
				blockHash, blockNumber, err := sbh.decodeAnchoringTx(data)
				if err != nil {
					logger.Error("failed to decode anchoring tx", "txHash", txHash.String(), "err", err)
					continue
				}
				sbh.WriteReceiptFromParentChain(blockHash, (*types.Receipt)(receipt))
				sbh.WriteAnchoredBlockNumber(blockNumber.Uint64())
			}
			// TODO-Klaytn-ServiceChain: support other tx types if needed.
			sbh.subbridge.GetBridgeTxPool().RemoveTx(tx)
		} else {
			logger.Trace("received service chain transaction receipt does not exist in sentServiceChainTxs", "txHash", txHash.String())
		}
		logger.Trace("received service chain transaction receipt", "txHash", txHash.String())
	}
}

func (sbh *SubBridgeHandler) RegisterNewPeer(p BridgePeer) error {
	sbh.subbridge.addPeerCh <- struct{}{}

	if sbh.getParentChainID().Cmp(p.GetChainID()) != 0 {
		return fmt.Errorf("attempt to add a peer with different chainID failed! existing chainID: %v, new chainID: %v", sbh.getParentChainID(), p.GetChainID())
	}
	// sync nonce and gasprice with peer
	sbh.SyncNonceAndGasPrice()

	return nil
}

// broadcastServiceChainReceiptRequest broadcasts receipt requests for service chain transactions.
func (sbh *SubBridgeHandler) broadcastServiceChainReceiptRequest() {
	hashes := sbh.subbridge.GetBridgeTxPool().PendingTxHashesByAddress(sbh.GetParentOperatorAddr(), int(sbh.GetSentChainTxsLimit())) // TODO-Klaytn-Servicechain change GetSentChainTxsLimit type to int from uint64
	for _, peer := range sbh.subbridge.BridgePeerSet().peers {
		peer.SendServiceChainReceiptRequest(hashes)
		logger.Debug("sent ServiceChainReceiptRequest", "peerID", peer.GetID(), "numReceiptsRequested", len(hashes))
	}
}

// updateTxCount update txCount to insert into anchoring tx. It skips first remnant txCount.
func (sbh *SubBridgeHandler) updateTxCount(block *types.Block) {
	if sbh.txCountEnabledBlockNumber == 0 {
		sbh.txCount = 0 // reset for the next anchoring period
		sbh.txCountEnabledBlockNumber = block.NumberU64()
		if sbh.chainTxPeriod > 1 {
			remnant := block.NumberU64() % sbh.chainTxPeriod
			if remnant < 2 {
				// A small trick to start tx counting quickly.
				sbh.txCountEnabledBlockNumber += 1 - remnant
			} else {
				sbh.txCountEnabledBlockNumber += (sbh.chainTxPeriod - remnant) + 1
			}
		}
	}

	var startBlkNum uint64
	if sbh.latestTxCountAddedBlockNumber == 0 {
		startBlkNum = block.NumberU64()
	} else {
		startBlkNum = sbh.latestTxCountAddedBlockNumber + 1
	}

	if startBlkNum < sbh.txCountEnabledBlockNumber {
		startBlkNum = sbh.txCountEnabledBlockNumber
	}

	for i := startBlkNum; i <= block.NumberU64(); i++ {
		b := sbh.subbridge.blockchain.GetBlockByNumber(i)
		if b == nil {
			logger.Warn("blockAnchoringManager: break to generateAndAddAnchoringTxIntoTxPool by the missed block", "missedBlockNumber", i)
			break
		}
		sbh.txCount += uint64(b.Transactions().Len())
		sbh.UpdateLatestTxCountAddedBlockNumber(i)
	}
}

// blockAnchoringManager generates anchoring transactions and updates transaction count.
func (sbh *SubBridgeHandler) blockAnchoringManager(block *types.Block) {
	sbh.updateTxCount(block)
	sbh.generateAndAddAnchoringTxIntoTxPool(block)
}

func (sbh *SubBridgeHandler) generateAndAddAnchoringTxIntoTxPool(block *types.Block) error {
	// Generating Anchoring Tx
	if block.NumberU64()%sbh.chainTxPeriod != 0 {
		return nil
	}
	sbh.LockParentOperator()
	defer sbh.UnLockParentOperator()

	unsignedTx, err := sbh.genUnsignedChainDataAnchoringTx(block)
	if err != nil {
		logger.Error("Failed to generate service chain transaction", "blockNum", block.NumberU64(), "err", err)
		return err
	}
	txCount := sbh.txCount
	sbh.txCount = 0 // reset for the next anchoring period

	signedTx, err := sbh.subbridge.bridgeAccounts.pAccount.SignTx(unsignedTx)
	if err != nil {
		logger.Error("failed signing tx", "err", err)
		return err
	}
	if err := sbh.subbridge.GetBridgeTxPool().AddLocal(signedTx); err == nil {
		sbh.addParentOperatorNonce(1)
	} else {
		logger.Debug("failed to add tx into bridge txpool", "err", err)
		return err
	}

	logger.Info("generate an anchoring tx", "blockNum", block.NumberU64(), "blockhash", block.Hash().String(), "txCount", txCount, "txHash", signedTx.Hash().String())

	return nil
}

// SyncNonceAndGasPrice requests the nonce of address used for service chain tx to parent chain peers.
func (scpm *SubBridgeHandler) SyncNonceAndGasPrice() {
	addr := scpm.GetParentOperatorAddr()
	for _, peer := range scpm.subbridge.BridgePeerSet().peers {
		peer.SendServiceChainInfoRequest(addr)
	}
}

// GetLatestAnchoredBlockNumber returns the latest block number whose data has been anchored to the parent chain.
func (sbh *SubBridgeHandler) GetLatestAnchoredBlockNumber() uint64 {
	return sbh.subbridge.ChainDB().ReadAnchoredBlockNumber()
}

// UpdateLatestTxCountAddedBlockNumber sets the latestTxCountAddedBlockNumber to the block number of the last anchoring tx which was added into bridge txPool.
func (sbh *SubBridgeHandler) UpdateLatestTxCountAddedBlockNumber(newLatestAnchoredBN uint64) {
	if sbh.latestTxCountAddedBlockNumber < newLatestAnchoredBN {
		sbh.latestTxCountAddedBlockNumber = newLatestAnchoredBN
	}
}

// WriteAnchoredBlockNumber writes the block number whose data has been anchored to the parent chain.
func (sbh *SubBridgeHandler) WriteAnchoredBlockNumber(blockNum uint64) {
	if sbh.GetLatestAnchoredBlockNumber() < blockNum {
		sbh.subbridge.chainDB.WriteAnchoredBlockNumber(blockNum)
		lastAnchoredBlockNumGauge.Update(int64(blockNum))
	}
}

// WriteReceiptFromParentChain writes a receipt received from parent chain to child chain
// with corresponding block hash. It assumes that a child chain has only one parent chain.
func (sbh *SubBridgeHandler) WriteReceiptFromParentChain(blockHash common.Hash, receipt *types.Receipt) {
	sbh.subbridge.chainDB.WriteReceiptFromParentChain(blockHash, receipt)
}

// GetReceiptFromParentChain returns a receipt received from parent chain to child chain
// with corresponding block hash. It assumes that a child chain has only one parent chain.
func (sbh *SubBridgeHandler) GetReceiptFromParentChain(blockHash common.Hash) *types.Receipt {
	return sbh.subbridge.chainDB.ReadReceiptFromParentChain(blockHash)
}
