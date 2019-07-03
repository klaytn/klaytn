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
	"crypto/ecdsa"
	"fmt"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/ser/rlp"
	"math/big"
)

const (
	SyncRequestInterval = 10
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
	// chainKey is a private key for account in parent chain, owned by service chain admin.
	// Used for signing transaction executed on the parent chain.
	chainKey *ecdsa.PrivateKey
	// MainChainAccountAddr is a hex account address used for chain identification from parent chain.
	MainChainAccountAddr *common.Address
	// remoteGasPrice means gas price of parent chain, used to make a service chain transaction.
	// Therefore, for now, it is only used by child chain side.
	remoteGasPrice        uint64
	mainChainAccountNonce uint64
	nonceSynced           bool
	chainTxPeriod         uint64

	// This is the block number of the latest anchoring tx which is added into bridge txPool.
	latestAnchoredBlockNumber uint64

	nodeKey                 *ecdsa.PrivateKey
	ServiceChainAccountAddr *common.Address

	// TODO-Klaytn-ServiceChain Need to limit the number independently? Or just managing the size of sentServiceChainTxs?
	sentServiceChainTxsLimit uint64

	skipSyncBlockCount int32
}

func NewSubBridgeHandler(scc *SCConfig, main *SubBridge) (*SubBridgeHandler, error) {
	// initialize the main chain account
	var mainChainAccountAddr *common.Address
	if scc.MainChainAccountAddr != nil {
		mainChainAccountAddr = scc.MainChainAccountAddr
	} else {
		chainKeyAddr := crypto.PubkeyToAddress(scc.chainkey.PublicKey)
		mainChainAccountAddr = &chainKeyAddr
		scc.MainChainAccountAddr = mainChainAccountAddr
	}
	// initialize ServiceChainAccount
	var serviceChainAccountAddr *common.Address
	if scc.ServiceChainAccountAddr != nil {
		serviceChainAccountAddr = scc.ServiceChainAccountAddr
	} else {
		nodeKeyAddr := crypto.PubkeyToAddress(scc.nodekey.PublicKey)
		serviceChainAccountAddr = &nodeKeyAddr
		scc.ServiceChainAccountAddr = serviceChainAccountAddr
	}
	return &SubBridgeHandler{
		subbridge:                 main,
		parentChainID:             new(big.Int).SetUint64(scc.ParentChainID),
		MainChainAccountAddr:      mainChainAccountAddr,
		chainKey:                  scc.chainkey,
		remoteGasPrice:            uint64(0),
		mainChainAccountNonce:     uint64(0),
		nonceSynced:               false,
		chainTxPeriod:             scc.AnchoringPeriod,
		latestAnchoredBlockNumber: uint64(0),
		sentServiceChainTxsLimit:  scc.SentChainTxsLimit,
		ServiceChainAccountAddr:   serviceChainAccountAddr,
		nodeKey:                   scc.nodekey,
	}, nil
}

func (sbh *SubBridgeHandler) setParentChainID(chainId *big.Int) {
	sbh.parentChainID = chainId
	sbh.subbridge.bridgeAccountManager.mcAccount.SetChainID(chainId)
}

func (sbh *SubBridgeHandler) getParentChainID() *big.Int {
	return sbh.parentChainID
}

func (sbh *SubBridgeHandler) LockMainChainAccount() {
	sbh.subbridge.bridgeAccountManager.mcAccount.Lock()
}

func (sbh *SubBridgeHandler) UnLockMainChainAccount() {
	sbh.subbridge.bridgeAccountManager.mcAccount.UnLock()
}

// getMainChainAccountNonce returns the main chain account nonce of main chain account address.
func (sbh *SubBridgeHandler) getMainChainAccountNonce() uint64 {
	return sbh.subbridge.bridgeAccountManager.mcAccount.GetNonce()
}

// setMainChainAccountNonce sets the main chain account nonce of main chain account address.
func (sbh *SubBridgeHandler) setMainChainAccountNonce(newNonce uint64) {
	sbh.subbridge.bridgeAccountManager.mcAccount.SetNonce(newNonce)
}

// addMainChainAccountNonce increases nonce by number
func (sbh *SubBridgeHandler) addMainChainAccountNonce(number uint64) {
	sbh.subbridge.bridgeAccountManager.mcAccount.IncNonce()
}

// getMainChainAccountNonceSynced returns whether the main chain account nonce is synced or not.
func (sbh *SubBridgeHandler) getMainChainAccountNonceSynced() bool {
	return sbh.nonceSynced
}

// setMainChainAccountNonceSynced sets whether the main chain account nonce is synced or not.
func (sbh *SubBridgeHandler) setMainChainAccountNonceSynced(synced bool) {
	sbh.nonceSynced = synced
}

func (sbh *SubBridgeHandler) getServiceChainAccountNonce() uint64 {
	return sbh.subbridge.txPool.GetPendingNonce(*sbh.ServiceChainAccountAddr)
}

func (sbh *SubBridgeHandler) getRemoteGasPrice() uint64 {
	return sbh.remoteGasPrice
}

func (sbh *SubBridgeHandler) setRemoteGasPrice(gasPrice uint64) {
	sbh.subbridge.bridgeAccountManager.mcAccount.SetGasPrice(big.NewInt(int64(gasPrice)))
	sbh.remoteGasPrice = gasPrice
}

// GetMainChainAccountAddr returns a pointer of a hex address of an account used for parent chain.
// If given as a parameter, it will use it. If not given, it will use the address of the public key
// derived from chainKey.
func (sbh *SubBridgeHandler) GetMainChainAccountAddr() *common.Address {
	return &sbh.subbridge.bridgeAccountManager.mcAccount.address
}

// GetServiceChainAccountAddr returns a pointer of a hex address of an account used for service chain.
// If given as a parameter, it will use it. If not given, it will use the address of the public key
// derived from chainKey.
func (sbh *SubBridgeHandler) GetServiceChainAccountAddr() *common.Address {
	return &sbh.subbridge.bridgeAccountManager.scAccount.address
}

// getChainKey returns the private key used for signing parent chain tx.
func (sbh *SubBridgeHandler) getChainKey() *ecdsa.PrivateKey {
	return sbh.subbridge.bridgeAccountManager.mcAccount.key
}

// getNodeKey returns the private key used for signing service chain tx.
func (sbh *SubBridgeHandler) getNodeKey() *ecdsa.PrivateKey {
	return sbh.subbridge.bridgeAccountManager.scAccount.key
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
		sbh.subbridge.rpcConn.Write(data)
		return nil
	case StatusMsg:
		return nil
	case ServiceChainParentChainInfoResponseMsg:
		logger.Debug("received ServiceChainParentChainInfoResponseMsg")
		if err := sbh.handleServiceChainParentChainInfoResponseMsg(p, msg); err != nil {
			return err
		}

	case ServiceChainReceiptResponseMsg:
		logger.Debug("received ServiceChainReceiptResponseMsg")
		if err := sbh.handleServiceChainReceiptResponseMsg(p, msg); err != nil {
			return err
		}
	default:
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}
	return nil
}

// handleServiceChainParentChainInfoResponseMsg handles parent chain info response message from parent chain.
// It will update the mainChainAccountNonce and remoteGasPrice of ServiceChainProtocolManager.
func (sbh *SubBridgeHandler) handleServiceChainParentChainInfoResponseMsg(p BridgePeer, msg p2p.Msg) error {
	var pcInfo parentChainInfo
	if err := msg.Decode(&pcInfo); err != nil {
		logger.Error("failed to decode", "err", err)
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	sbh.LockMainChainAccount()
	defer sbh.UnLockMainChainAccount()

	poolNonce := sbh.subbridge.bridgeTxPool.GetMaxTxNonce(sbh.GetMainChainAccountAddr())
	if poolNonce > 0 {
		poolNonce += 1
		// just check
		if sbh.getMainChainAccountNonce() > poolNonce {
			logger.Error("main chain account nonce is bigger than the chain pool nonce.", "BridgeTxPoolNonce", poolNonce, "mainChainAccountNonce", sbh.getMainChainAccountNonce())
		}
		if poolNonce < pcInfo.Nonce {
			// bridgeTxPool journal miss txs which already sent to parent-chain
			logger.Error("chain pool nonce is less than the parent chain nonce.", "chainPoolNonce", poolNonce, "receivedNonce", pcInfo.Nonce)
			sbh.setMainChainAccountNonce(pcInfo.Nonce)
		} else {
			// bridgeTxPool journal has txs which don't receive receipt from parent-chain
			sbh.setMainChainAccountNonce(poolNonce)
		}
	} else if sbh.getMainChainAccountNonce() > pcInfo.Nonce {
		logger.Error("main chain account nonce is bigger than the received nonce.", "mainChainAccountNonce", sbh.getMainChainAccountNonce(), "receivedNonce", pcInfo.Nonce)
		sbh.setMainChainAccountNonce(pcInfo.Nonce)
	} else {
		// there is no tx in bridgetTxPool, so parent-chain's nonce is used
		sbh.setMainChainAccountNonce(pcInfo.Nonce)
	}
	sbh.setMainChainAccountNonceSynced(true)
	sbh.setRemoteGasPrice(pcInfo.GasPrice)
	logger.Info("ServiceChainNonceResponse", "receivedNonce", pcInfo.Nonce, "gasPrice", pcInfo.GasPrice, "mainChainAccountNonce", sbh.getMainChainAccountNonce())
	return nil
}

// handleServiceChainReceiptResponseMsg handles receipt response message from parent chain.
// It will store the received receipts and remove corresponding transaction in the resending list.
func (sbh *SubBridgeHandler) handleServiceChainReceiptResponseMsg(p BridgePeer, msg p2p.Msg) error {
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

// genUnsignedServiceChainTx generates an unsigned transaction, which type is TxTypeChainDataAnchoring.
// Nonce of account used for service chain transaction will be increased after the signing.
func (sbh *SubBridgeHandler) genUnsignedServiceChainTx(block *types.Block) (*types.Transaction, error) {
	chainHashes := types.NewChainHashes(block)
	encodedCCTxData, err := rlp.EncodeToBytes(chainHashes)
	if err != nil {
		return nil, err
	}

	values := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:        sbh.getMainChainAccountNonce(), // main chain account nonce will be increased after signing a transaction.
		types.TxValueKeyFrom:         *sbh.GetMainChainAccountAddr(),
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
	if sbh.getMainChainAccountNonceSynced() {
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
	txs := sbh.subbridge.GetBridgeTxPool().PendingTxsByAddress(sbh.MainChainAccountAddr, int(sbh.GetSentChainTxsLimit())) // TODO-Klaytn-Servicechain change GetSentChainTxsLimit type to int from uint64
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

// writeServiceChainTxReceipt writes the received receipts of service chain transactions.
func (sbh *SubBridgeHandler) writeServiceChainTxReceipts(bc *blockchain.BlockChain, receipts []*types.ReceiptForStorage) {
	for _, receipt := range receipts {
		txHash := receipt.TxHash
		if tx := sbh.subbridge.GetBridgeTxPool().Get(txHash); tx != nil {
			if tx.Type() == types.TxTypeChainDataAnchoring {
				chainHashes := new(types.ChainHashes)
				data, err := tx.AnchoredData()
				if err != nil {
					logger.Error("failed to get anchoring tx type from the tx", "txHash", txHash.String())
					return
				}
				if err := rlp.DecodeBytes(data, chainHashes); err != nil {
					logger.Error("failed to RLP decode ChainHashes", "txHash", txHash.String())
					return
				}
				sbh.WriteReceiptFromParentChain(chainHashes.BlockHash, (*types.Receipt)(receipt))
				sbh.WriteAnchoredBlockNumber(chainHashes.BlockNumber.Uint64())
				logger.Debug("received anchoring tx receipt", "blockNum", chainHashes.BlockNumber.String(), "blcokHash", chainHashes.BlockHash.String(), "txHash", txHash.String())
			}

			sbh.subbridge.GetBridgeTxPool().RemoveTx(tx)
		} else {
			logger.Trace("received service chain transaction receipt does not exist in sentServiceChainTxs", "txHash", txHash.String())
		}

		logger.Trace("received service chain transaction receipt", "txHash", txHash.String())
	}
}

func (sbh *SubBridgeHandler) RegisterNewPeer(p BridgePeer) error {
	if sbh.getParentChainID().Cmp(p.GetChainID()) != 0 {
		return fmt.Errorf("attempt to add a peer with different chainID failed! existing chainID: %v, new chainID: %v", sbh.getParentChainID(), p.GetChainID())
	}
	// sync nonce and gasprice with peer
	sbh.SyncNonceAndGasPrice()
	return nil
}

// broadcastServiceChainReceiptRequest broadcasts receipt requests for service chain transactions.
func (sbh *SubBridgeHandler) broadcastServiceChainReceiptRequest() {
	hashes := sbh.subbridge.GetBridgeTxPool().PendingTxHashesByAddress(sbh.GetMainChainAccountAddr(), int(sbh.GetSentChainTxsLimit())) // TODO-Klaytn-Servicechain change GetSentChainTxsLimit type to int from uint64
	for _, peer := range sbh.subbridge.BridgePeerSet().peers {
		peer.SendServiceChainReceiptRequest(hashes)
		logger.Debug("sent ServiceChainReceiptRequest", "peerID", peer.GetID(), "numReceiptsRequested", len(hashes))
	}
}

func (sbh *SubBridgeHandler) blockAnchoringManager(block *types.Block) {
	startBlkNum := sbh.GetNextAnchoringBlockNumber()

	var successCnt, cnt, blkNum uint64
	latestBlkNum := block.Number().Uint64()

	// TODO-Klaytn-Servicechain remove this code or define right confirmation block count.
	// To consider a non-absolute finality consensus like PoA (Clique), this code anchors past blocks.
	const confirmCnt = uint64(10)
	if latestBlkNum > confirmCnt {
		latestBlkNum = latestBlkNum - confirmCnt
	} else {
		return
	}

	for cnt, blkNum = 0, startBlkNum; cnt <= sbh.sentServiceChainTxsLimit && blkNum <= latestBlkNum; cnt, blkNum = cnt+1, blkNum+1 {
		block := sbh.subbridge.blockchain.GetBlockByNumber(blkNum)
		if block == nil {
			logger.Warn("blockAnchoringManager: break to generateAndAddAnchoringTxIntoTxPool by the missed block", "missedBlockNumber", blkNum)
			break
		}
		if err := sbh.generateAndAddAnchoringTxIntoTxPool(block); err == nil {
			sbh.UpdateLastestAnchoredBlockNumber(blkNum)
			successCnt++
		} else {
			logger.Trace("blockAnchoringManager: break to generateAndAddAnchoringTxIntoTxPool", "cnt", cnt, "startBlockNumber", startBlkNum, "FailedBlockNumber", blkNum, "latestBlockNum", block.NumberU64())
			break
		}
	}
	if successCnt > 0 {
		logger.Info("Generate anchoring txs", "txCount", successCnt, "startBlockNumber", startBlkNum, "endBlockNumber", blkNum-1)
	}
}

func (sbh *SubBridgeHandler) generateAndAddAnchoringTxIntoTxPool(block *types.Block) error {
	// Generating Anchoring Tx
	if block.NumberU64()%sbh.chainTxPeriod != 0 {
		return nil
	}
	sbh.LockMainChainAccount()
	defer sbh.UnLockMainChainAccount()

	unsignedTx, err := sbh.genUnsignedServiceChainTx(block)
	if err != nil {
		logger.Error("Failed to generate service chain transaction", "blockNum", block.NumberU64(), "err", err)
		return err
	}
	// TODO-Klaytn-ServiceChain Change types.NewEIP155Signer to types.MakeSigner using parent chain's chain config and block number
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(sbh.parentChainID), sbh.getChainKey())
	if err != nil {
		logger.Error("failed signing tx", "err", err)
		return err
	}
	if err := sbh.subbridge.GetBridgeTxPool().AddLocal(signedTx); err == nil {
		sbh.addMainChainAccountNonce(1)
	} else {
		logger.Debug("failed to add tx into bridge txpool", "err", err)
		return err
	}

	logger.Trace("blockAnchoringManager: Success to generate anchoring tx", "blockNum", block.NumberU64(), "blockhash", block.Hash().String(), "txHash", signedTx.Hash().String())

	return nil
}

// SyncNonceAndGasPrice requests the nonce of address used for service chain tx to parent chain peers.
func (scpm *SubBridgeHandler) SyncNonceAndGasPrice() {
	addr := scpm.GetMainChainAccountAddr()
	for _, peer := range scpm.subbridge.BridgePeerSet().peers {
		peer.SendServiceChainInfoRequest(addr)
	}
}

// GetLatestAnchoredBlockNumber returns the latest block number whose data has been anchored to the parent chain.
func (sbh *SubBridgeHandler) GetLatestAnchoredBlockNumber() uint64 {
	return sbh.subbridge.ChainDB().ReadAnchoredBlockNumber()
}

// GetNextAnchoringBlockNumber returns the next block number which is needed to be anchored.
func (sbh *SubBridgeHandler) GetNextAnchoringBlockNumber() uint64 {
	if sbh.latestAnchoredBlockNumber == 0 {
		sbh.latestAnchoredBlockNumber = sbh.subbridge.ChainDB().ReadAnchoredBlockNumber()
	}

	// If latestAnchoredBlockNumber == 0, there are two cases below.
	// 1) The last block number anchored is 0 block(genesis block).
	// 2) There is no block anchored, so this is the 1st time. (If there is no value in DB, it returns 0.)
	// To cover all cases without complex DB routine, the condition below is added.
	// Even if genesis block can be anchored more than 2 times,
	// this routine can guarantee anchoring genesis block.
	if sbh.latestAnchoredBlockNumber == 0 {
		return sbh.latestAnchoredBlockNumber
	}

	return sbh.latestAnchoredBlockNumber + 1
}

// UpdateLastestAnchoredBlockNumber set the latestAnchoredBlockNumber to the block number of the last anchoring tx which was added into bridge txPool.
func (sbh *SubBridgeHandler) UpdateLastestAnchoredBlockNumber(newLastestAnchoredBN uint64) {
	if sbh.latestAnchoredBlockNumber < newLastestAnchoredBN {
		sbh.latestAnchoredBlockNumber = newLastestAnchoredBN
	}
}

// WriteAnchoredBlockNumber writes the block number whose data has been anchored to the parent chain.
func (sbh *SubBridgeHandler) WriteAnchoredBlockNumber(blockNum uint64) {
	sbh.UpdateLastestAnchoredBlockNumber(blockNum)
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
