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
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/rlp"
)

const (
	SyncRequestInterval = 10
)

var (
	ErrInvalidBlock              = errors.New("block is invalid")
	ErrUnknownBridgeContractAddr = errors.New("The given address was not found in the bridge contract list")
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
	txCountStartingBlockNumber    uint64
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

// getParentOperatorBalance returns the parent chain operator balance.
func (sbh *SubBridgeHandler) getParentOperatorBalance() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return sbh.subbridge.remoteBackend.BalanceAt(ctx, sbh.subbridge.bridgeAccounts.pAccount.address, nil)
}

// getParentBridgeContractBalance returns the parent bridge contract's balance.
func (sbh *SubBridgeHandler) getParentBridgeContractBalance(addr common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if !sbh.subbridge.bridgeManager.IsInParentAddrs(addr) {
		logger.Error(ErrUnknownBridgeContractAddr.Error(), "addr", addr)
		return common.Big0, ErrUnknownBridgeContractAddr
	}
	return sbh.subbridge.remoteBackend.BalanceAt(ctx, addr, nil)
}

// setParentOperatorNonceSynced sets whether the parent chain operator account nonce is synced or not.
func (sbh *SubBridgeHandler) setParentOperatorNonceSynced(synced bool) {
	sbh.nonceSynced = synced
}

func (sbh *SubBridgeHandler) getChildOperatorNonce() uint64 {
	return sbh.subbridge.txPool.GetPendingNonce(sbh.subbridge.bridgeAccounts.cAccount.address)
}

// getChildOperatorBalance returns the child chain operator balance.
func (sbh *SubBridgeHandler) getChildOperatorBalance() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return sbh.subbridge.localBackend.BalanceAt(ctx, sbh.subbridge.bridgeAccounts.cAccount.address, nil)
}

// getChildBridgeContractBalance returns the child bridge contract's balance.
func (sbh *SubBridgeHandler) getChildBridgeContractBalance(addr common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if !sbh.subbridge.bridgeManager.IsInChildAddrs(addr) {
		logger.Error(ErrUnknownBridgeContractAddr.Error(), "addr", addr)
		return common.Big0, ErrUnknownBridgeContractAddr
	}
	return sbh.subbridge.localBackend.BalanceAt(ctx, addr, nil)
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
			// BridgeTxPool journal miss txs which already sent to parent-chain
			logger.Error("chain pool nonce is less than the parent chain nonce.", "chainPoolNonce", poolNonce, "receivedNonce", pcInfo.Nonce)
			sbh.setParentOperatorNonce(pcInfo.Nonce)
		} else {
			// BridgeTxPool journal has txs which don't receive receipt from parent-chain
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
	anchoringData, err := types.NewAnchoringDataType0(block, block.NumberU64()-sbh.txCountStartingBlockNumber+1, sbh.txCount)
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

	txType := types.TxTypeChainDataAnchoring

	if feePayer := sbh.subbridge.bridgeAccounts.GetParentOperatorFeePayer(); feePayer != (common.Address{}) {
		values[types.TxValueKeyFeePayer] = feePayer
		txType = types.TxTypeFeeDelegatedChainDataAnchoring
	}

	if tx, err := types.NewTransactionWithMap(txType, values); err != nil {
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
		sbh.txCountStartingBlockNumber = 0
		if sbh.skipSyncBlockCount%SyncRequestInterval == 0 {
			// TODO-Klaytn too many request while sync main-net
			sbh.SyncNonceAndGasPrice()
			// check tx's receipts which parent-chain already executed in BridgeTxPool
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
	logger.Trace("broadcastServiceChainTx ServiceChainTxData", "len(txs)", len(txs), "len(peers)", len(peers))
}

// writeServiceChainTxReceipts writes the received receipts of service chain transactions.
func (sbh *SubBridgeHandler) writeServiceChainTxReceipts(bc *blockchain.BlockChain, receipts []*types.ReceiptForStorage) {
	for _, receipt := range receipts {
		txHash := receipt.TxHash
		if tx := sbh.subbridge.GetBridgeTxPool().Get(txHash); tx != nil {
			if tx.Type().IsChainDataAnchoring() {
				data, err := tx.AnchoredData()
				if err != nil {
					logger.Error("failed to get anchoring data", "txHash", txHash.String(), "err", err)
					continue
				}
				decodedData, err := types.DecodeAnchoringData(data)
				if err != nil {
					logger.Error("failed to decode anchoring tx", "txHash", txHash.String(), "err", err)
					continue
				}
				sbh.WriteReceiptFromParentChain(decodedData.GetBlockHash(), (*types.Receipt)(receipt))
				sbh.WriteAnchoredBlockNumber(decodedData.GetBlockNumber().Uint64())
			}
			// TODO-Klaytn-ServiceChain: support other tx types if needed.
			sbh.subbridge.GetBridgeTxPool().RemoveTx(tx)
		} else {
			logger.Trace("received service chain transaction receipt does not exist in sentServiceChainTxs", "txHash", txHash.String())
		}
		logger.Trace("received service chain transaction receipt", "anchoring txHash", txHash.String())
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

// updateTxCount update txCount to insert into anchoring tx.
func (sbh *SubBridgeHandler) updateTxCount(block *types.Block) error {
	if block == nil {
		return ErrInvalidBlock
	}

	if sbh.txCountStartingBlockNumber == 0 {
		sbh.txCount = 0 // reset for the next anchoring period
		sbh.txCountStartingBlockNumber = block.NumberU64()
	}

	var startBlkNum uint64
	if sbh.latestTxCountAddedBlockNumber == 0 {
		startBlkNum = block.NumberU64()
	} else {
		startBlkNum = sbh.latestTxCountAddedBlockNumber + 1
	}

	if startBlkNum < sbh.txCountStartingBlockNumber {
		startBlkNum = sbh.txCountStartingBlockNumber
	}

	for i := startBlkNum; i <= block.NumberU64(); i++ {
		b := sbh.subbridge.blockchain.GetBlockByNumber(i)
		if b == nil {
			logger.Error("blockAnchoringManager: break to generateAndAddAnchoringTxIntoTxPool by the missed block", "missedBlockNumber", i)
			break
		}
		sbh.txCount += uint64(b.Transactions().Len())
		sbh.UpdateLatestTxCountAddedBlockNumber(i)
	}

	return nil
}

// blockAnchoringManager generates anchoring transactions and updates transaction count.
func (sbh *SubBridgeHandler) blockAnchoringManager(block *types.Block) error {
	if err := sbh.updateTxCount(block); err != nil {
		return err
	}
	return sbh.generateAndAddAnchoringTxIntoTxPool(block)
}

func (sbh *SubBridgeHandler) generateAndAddAnchoringTxIntoTxPool(block *types.Block) error {
	if block == nil {
		return ErrInvalidBlock
	}

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
	// Reset for the next anchoring period.
	sbh.txCount = 0
	sbh.txCountStartingBlockNumber = block.NumberU64() + 1

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

	logger.Info("Generate an anchoring tx", "blockNum", block.NumberU64(), "blockhash", block.Hash().String(), "txCount", txCount, "txHash", signedTx.Hash().String())

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
