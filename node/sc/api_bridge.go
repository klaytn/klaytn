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
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/node"
	"github.com/pkg/errors"
)

var (
	ErrInvalidBridgePair             = errors.New("Invalid bridge pair")
	ErrBridgeContractVersionMismatch = errors.New("Bridge contract version mismatch")
)

// MainBridgeAPI Implementation for main-bridge node
type MainBridgeAPI struct {
	mainBridge *MainBridge
}

func (mb *MainBridgeAPI) GetChildChainIndexingEnabled() bool {
	return mb.mainBridge.eventhandler.GetChildChainIndexingEnabled()
}

func (mb *MainBridgeAPI) ConvertChildChainBlockHashToParentChainTxHash(scBlockHash common.Hash) common.Hash {
	return mb.mainBridge.eventhandler.ConvertChildChainBlockHashToParentChainTxHash(scBlockHash)
}

// Peers retrieves all the information we know about each individual peer at the
// protocol granularity.
func (mb *MainBridgeAPI) Peers() ([]*p2p.PeerInfo, error) {
	server := mb.mainBridge.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.PeersInfo(), nil
}

// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.
func (mb *MainBridgeAPI) NodeInfo() (*p2p.NodeInfo, error) {
	server := mb.mainBridge.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.NodeInfo(), nil
}

// SubBridgeAPI Implementation for sub-bridge node
type SubBridgeAPI struct {
	subBridge *SubBridge
}

func (sb *SubBridgeAPI) ConvertChildChainBlockHashToParentChainTxHash(cBlockHash common.Hash) common.Hash {
	return sb.subBridge.eventhandler.ConvertChildChainBlockHashToParentChainTxHash(cBlockHash)
}

func (sb *SubBridgeAPI) GetLatestAnchoredBlockNumber() uint64 {
	return sb.subBridge.handler.GetLatestAnchoredBlockNumber()
}

func (sb *SubBridgeAPI) GetReceiptFromParentChain(blockHash common.Hash) *types.Receipt {
	return sb.subBridge.handler.GetReceiptFromParentChain(blockHash)
}

func (sb *SubBridgeAPI) GetAnchoringTxHashByBlockNumber(bn uint64) common.Hash {
	block := sb.subBridge.blockchain.GetBlockByNumber(bn)
	if block == nil {
		return common.Hash{}
	}
	receipt := sb.subBridge.handler.GetReceiptFromParentChain(block.Hash())
	if receipt == nil {
		return common.Hash{}
	}
	return receipt.TxHash
}

func (sb *SubBridgeAPI) RegisterOperator(bridgeAddr, operatorAddr common.Address) (common.Hash, error) {
	return sb.subBridge.bridgeManager.RegisterOperator(bridgeAddr, operatorAddr)
}

func (sb *SubBridgeAPI) GetRegisteredOperators(bridgeAddr common.Address) ([]common.Address, error) {
	return sb.subBridge.bridgeManager.GetOperators(bridgeAddr)
}

func (sb *SubBridgeAPI) SetValueTransferOperatorThreshold(bridgeAddr common.Address, threshold uint8) (common.Hash, error) {
	return sb.subBridge.bridgeManager.SetValueTransferOperatorThreshold(bridgeAddr, threshold)
}

func (sb *SubBridgeAPI) GetValueTransferOperatorThreshold(bridgeAddr common.Address) (uint8, error) {
	return sb.subBridge.bridgeManager.GetValueTransferOperatorThreshold(bridgeAddr)
}

func (sb *SubBridgeAPI) DeployBridge() ([]common.Address, error) {
	cAcc := sb.subBridge.bridgeAccounts.cAccount
	pAcc := sb.subBridge.bridgeAccounts.pAccount

	cBridge, cBridgeAddr, err := sb.subBridge.bridgeManager.DeployBridge(cAcc.GenerateTransactOpts(), sb.subBridge.localBackend, true)
	if err != nil {
		logger.Error("Failed to deploy child bridge.", "err", err)
		return nil, err
	}
	pBridge, pBridgeAddr, err := sb.subBridge.bridgeManager.DeployBridge(pAcc.GenerateTransactOpts(), sb.subBridge.remoteBackend, false)
	if err != nil {
		logger.Error("Failed to deploy parent bridge.", "err", err)
		return nil, err
	}

	err = sb.subBridge.bridgeManager.SetBridgeInfo(cBridgeAddr, cBridge, pBridgeAddr, pBridge, cAcc, true, false)
	if err != nil {
		return nil, err
	}

	err = sb.subBridge.bridgeManager.SetBridgeInfo(pBridgeAddr, pBridge, cBridgeAddr, cBridge, pAcc, false, false)
	if err != nil {
		return nil, err
	}

	err = sb.subBridge.bridgeManager.SetJournal(cBridgeAddr, pBridgeAddr)
	if err != nil {
		return nil, err
	}

	return []common.Address{cBridgeAddr, pBridgeAddr}, nil
}

// SubscribeBridge enables the given child/parent chain bridges to subscribe the events.
func (sb *SubBridgeAPI) SubscribeBridge(cBridgeAddr, pBridgeAddr common.Address) error {
	if !sb.subBridge.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	err := sb.subBridge.bridgeManager.SubscribeEvent(cBridgeAddr)
	if err != nil {
		logger.Error("Failed to SubscribeEvent child bridge", "addr", cBridgeAddr, "err", err)
		return err
	}

	err = sb.subBridge.bridgeManager.SubscribeEvent(pBridgeAddr)
	if err != nil {
		logger.Error("Failed to SubscribeEvent parent bridge", "addr", pBridgeAddr, "err", err)
		sb.subBridge.bridgeManager.UnsubscribeEvent(cBridgeAddr)
		return err
	}

	sb.subBridge.bridgeManager.journal.cache[cBridgeAddr].Subscribed = true

	// Update the journal's subscribed flag.
	sb.subBridge.bridgeManager.journal.rotate(sb.subBridge.bridgeManager.GetAllBridge())

	err = sb.subBridge.bridgeManager.AddRecovery(cBridgeAddr, pBridgeAddr)
	if err != nil {
		sb.subBridge.bridgeManager.UnsubscribeEvent(cBridgeAddr)
		sb.subBridge.bridgeManager.UnsubscribeEvent(pBridgeAddr)
		return err
	}
	return nil
}

// UnsubscribeBridge disables the event subscription of the given child/parent chain bridges.
func (sb *SubBridgeAPI) UnsubscribeBridge(cBridgeAddr, pBridgeAddr common.Address) error {
	if !sb.subBridge.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	sb.subBridge.bridgeManager.UnsubscribeEvent(cBridgeAddr)
	sb.subBridge.bridgeManager.UnsubscribeEvent(pBridgeAddr)

	sb.subBridge.bridgeManager.journal.cache[cBridgeAddr].Subscribed = false

	sb.subBridge.bridgeManager.journal.rotate(sb.subBridge.bridgeManager.GetAllBridge())
	return nil
}

func (sb *SubBridgeAPI) ConvertRequestTxHashToHandleTxHash(hash common.Hash) common.Hash {
	return sb.subBridge.chainDB.ReadHandleTxHashFromRequestTxHash(hash)
}

func (sb *SubBridgeAPI) TxPendingCount() int {
	return sb.subBridge.GetBridgeTxPool().Stats()
}

func (sb *SubBridgeAPI) TxPending() map[common.Address]types.Transactions {
	return sb.subBridge.GetBridgeTxPool().Pending()
}

func (sb *SubBridgeAPI) ListBridge() []*BridgeJournal {
	return sb.subBridge.bridgeManager.GetAllBridge()
}

func (sb *SubBridgeAPI) GetBridgeInformation(bridgeAddr common.Address) (map[string]interface{}, error) {
	if ctBridge := sb.subBridge.bridgeManager.GetCounterPartBridgeAddr(bridgeAddr); ctBridge == (common.Address{}) {
		return nil, ErrInvalidBridgePair
	}

	bi, ok := sb.subBridge.bridgeManager.GetBridgeInfo(bridgeAddr)
	if !ok {
		return nil, ErrNoBridgeInfo
	}

	bi.UpdateInfo()

	return map[string]interface{}{
		"isRunning":        bi.isRunning,
		"requestNonce":     bi.requestNonceFromCounterPart,
		"handleNonce":      bi.handleNonce,
		"lowerHandleNonce": bi.lowerHandleNonce,
		"counterPart":      bi.counterpartAddress,
		"onServiceChain":   bi.onChildChain,
		"isSubscribed":     bi.subscribed,
		"pendingEventSize": bi.pendingRequestEvent.Len(),
	}, nil
}

func (sb *SubBridgeAPI) KASAnchor(blkNum uint64) error {
	block := sb.subBridge.blockchain.GetBlockByNumber(blkNum)
	if block != nil {
		if err := sb.subBridge.kasAnchor.AnchorBlock(block); err != nil {
			logger.Error("Failed to anchor a block via KAS", "blkNum", block.NumberU64(), "err", err)
			return err
		}
		return nil
	}
	return ErrInvalidBlock
}

func (sb *SubBridgeAPI) Anchoring(flag bool) bool {
	return sb.subBridge.SetAnchoringTx(flag)
}

func (sb *SubBridgeAPI) GetAnchoring() bool {
	return sb.subBridge.GetAnchoringTx()
}

func (sb *SubBridgeAPI) RegisterBridge(cBridgeAddr common.Address, pBridgeAddr common.Address) error {
	cBridge, err := bridge.NewBridge(cBridgeAddr, sb.subBridge.localBackend)
	if err != nil {
		return err
	}
	pBridge, err := bridge.NewBridge(pBridgeAddr, sb.subBridge.remoteBackend)
	if err != nil {
		return err
	}

	bm := sb.subBridge.bridgeManager
	err = bm.SetBridgeInfo(cBridgeAddr, cBridge, pBridgeAddr, pBridge, sb.subBridge.bridgeAccounts.cAccount, true, false)
	if err != nil {
		return err
	}
	err = bm.SetBridgeInfo(pBridgeAddr, pBridge, cBridgeAddr, cBridge, sb.subBridge.bridgeAccounts.pAccount, false, false)
	if err != nil {
		bm.DeleteBridgeInfo(cBridgeAddr)
		return err
	}

	err = bm.SetJournal(cBridgeAddr, pBridgeAddr)
	if err != nil {
		return err
	}

	return nil
}

func (sb *SubBridgeAPI) DeregisterBridge(cBridgeAddr common.Address, pBridgeAddr common.Address) error {
	if !sb.subBridge.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	bm := sb.subBridge.bridgeManager
	journal := bm.journal.cache[cBridgeAddr]

	if journal.Subscribed {
		bm.UnsubscribeEvent(journal.ChildAddress)
		bm.UnsubscribeEvent(journal.ParentAddress)

		bm.DeleteRecovery(cBridgeAddr, pBridgeAddr)
	}

	delete(bm.journal.cache, cBridgeAddr)

	if err := bm.journal.rotate(bm.GetAllBridge()); err != nil {
		logger.Warn("failed to rotate bridge journal", "err", err, "cBridge", cBridgeAddr.String(), "pBridge", pBridgeAddr.String())
	}

	if err := bm.DeleteBridgeInfo(cBridgeAddr); err != nil {
		logger.Warn("failed to Delete child chain bridge info", "err", err, "bridge", cBridgeAddr.String())
	}

	if err := bm.DeleteBridgeInfo(pBridgeAddr); err != nil {
		logger.Warn("failed to Delete parent chain bridge info", "err", err, "bridge", pBridgeAddr.String())
	}

	return nil
}

func (sb *SubBridgeAPI) RegisterToken(cBridgeAddr, pBridgeAddr, cTokenAddr, pTokenAddr common.Address) error {
	if !sb.subBridge.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	cBi, cExist := sb.subBridge.bridgeManager.GetBridgeInfo(cBridgeAddr)
	pBi, pExist := sb.subBridge.bridgeManager.GetBridgeInfo(pBridgeAddr)

	if !cExist || !pExist {
		return ErrNoBridgeInfo
	}

	err := cBi.RegisterToken(cTokenAddr, pTokenAddr)
	if err != nil {
		return err
	}

	err = pBi.RegisterToken(pTokenAddr, cTokenAddr)
	if err != nil {
		return err
	}

	cBi.account.Lock()
	tx, err := cBi.bridge.RegisterToken(cBi.account.GenerateTransactOpts(), cTokenAddr, pTokenAddr)
	if err != nil {
		cBi.account.UnLock()
		return err
	}
	cBi.account.IncNonce()
	cBi.account.UnLock()
	logger.Debug("cBridge registered token", "txHash", tx.Hash().String(), "cToken", cTokenAddr.String(), "pToken", pTokenAddr.String())

	pBi.account.Lock()
	tx, err = pBi.bridge.RegisterToken(pBi.account.GenerateTransactOpts(), pTokenAddr, cTokenAddr)
	if err != nil {
		pBi.account.UnLock()
		return err
	}
	pBi.account.IncNonce()
	pBi.account.UnLock()
	logger.Debug("pBridge registered token", "txHash", tx.Hash().String(), "cToken", cTokenAddr.String(), "pToken", pTokenAddr.String())

	logger.Info("Register token", "cToken", cTokenAddr.String(), "pToken", pTokenAddr.String())
	return nil
}

func (sb *SubBridgeAPI) GetParentTransactionReceipt(txHash common.Hash) (map[string]interface{}, error) {
	ctx := context.Background()
	return sb.subBridge.remoteBackend.(RemoteBackendInterface).TransactionReceiptRpcOutput(ctx, txHash)
}

func (sb *SubBridgeAPI) SetERC20Fee(bridgeAddr, tokenAddr common.Address, fee *big.Int) (common.Hash, error) {
	return sb.subBridge.bridgeManager.SetERC20Fee(bridgeAddr, tokenAddr, fee)
}

func (sb *SubBridgeAPI) SetKLAYFee(bridgeAddr common.Address, fee *big.Int) (common.Hash, error) {
	return sb.subBridge.bridgeManager.SetKLAYFee(bridgeAddr, fee)
}

func (sb *SubBridgeAPI) SetFeeReceiver(bridgeAddr, receiver common.Address) (common.Hash, error) {
	return sb.subBridge.bridgeManager.SetFeeReceiver(bridgeAddr, receiver)
}

func (sb *SubBridgeAPI) GetERC20Fee(bridgeAddr, tokenAddr common.Address) (*big.Int, error) {
	return sb.subBridge.bridgeManager.GetERC20Fee(bridgeAddr, tokenAddr)
}

func (sb *SubBridgeAPI) GetKLAYFee(bridgeAddr common.Address) (*big.Int, error) {
	return sb.subBridge.bridgeManager.GetKLAYFee(bridgeAddr)
}

func (sb *SubBridgeAPI) GetFeeReceiver(bridgeAddr common.Address) (common.Address, error) {
	return sb.subBridge.bridgeManager.GetFeeReceiver(bridgeAddr)
}

func (sb *SubBridgeAPI) DeregisterToken(cBridgeAddr, pBridgeAddr, cTokenAddr, pTokenAddr common.Address) error {
	cBi, cExist := sb.subBridge.bridgeManager.GetBridgeInfo(cBridgeAddr)
	pBi, pExist := sb.subBridge.bridgeManager.GetBridgeInfo(pBridgeAddr)

	if !cExist || !pExist {
		return ErrNoBridgeInfo
	}

	pTokenAddrCheck := cBi.GetCounterPartToken(cTokenAddr)
	cTokenAddrCheck := pBi.GetCounterPartToken(pTokenAddr)

	if pTokenAddr != pTokenAddrCheck || cTokenAddr != cTokenAddrCheck {
		return errors.New("invalid toke pair")
	}

	cBi.DeregisterToken(cTokenAddr, pTokenAddr)
	pBi.DeregisterToken(pTokenAddr, cTokenAddr)

	cBi.account.Lock()
	defer cBi.account.UnLock()
	tx, err := cBi.bridge.DeregisterToken(cBi.account.GenerateTransactOpts(), cTokenAddr)
	if err != nil {
		return err
	}
	cBi.account.IncNonce()
	logger.Debug("cBridge deregistered token", "txHash", tx.Hash().String(), "cToken", cTokenAddr.String(), "pToken", pTokenAddr.String())

	pBi.account.Lock()
	defer pBi.account.UnLock()
	tx, err = pBi.bridge.DeregisterToken(pBi.account.GenerateTransactOpts(), pTokenAddr)
	if err != nil {
		return err
	}
	pBi.account.IncNonce()
	logger.Debug("pBridge deregistered token", "txHash", tx.Hash().String(), "cToken", cTokenAddr.String(), "pToken", pTokenAddr.String())
	return err
}

// AddPeer requests connecting to a remote node, and also maintaining the new
// connection at all times, even reconnecting if it is lost.
func (sb *SubBridgeAPI) AddPeer(url string) (bool, error) {
	// Make sure the server is running, fail otherwise
	server := sb.subBridge.bridgeServer
	if server == nil {
		return false, node.ErrNodeStopped
	}
	// TODO-Klaytn Refactoring this to check whether the url is valid or not by dialing and return it.
	if _, err := addPeerInternal(server, url); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// addPeerInternal does common part for AddPeer.
func addPeerInternal(server p2p.Server, url string) (*discover.Node, error) {
	// Try to add the url as a static peer and return
	node, err := discover.ParseNode(url)
	if err != nil {
		return nil, fmt.Errorf("invalid kni: %v", err)
	}
	server.AddPeer(node)
	return node, nil
}

// RemovePeer disconnects from a a remote node if the connection exists
func (sb *SubBridgeAPI) RemovePeer(url string) (bool, error) {
	// Make sure the server is running, fail otherwise
	server := sb.subBridge.bridgeServer
	if server == nil {
		return false, node.ErrNodeStopped
	}
	// Try to remove the url as a static peer and return
	node, err := discover.ParseNode(url)
	if err != nil {
		return false, fmt.Errorf("invalid kni: %v", err)
	}
	server.RemovePeer(node)
	return true, nil
}

// Peers retrieves all the information we know about each individual peer at the
// protocol granularity.
func (sb *SubBridgeAPI) Peers() ([]*p2p.PeerInfo, error) {
	server := sb.subBridge.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.PeersInfo(), nil
}

// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.
func (sb *SubBridgeAPI) NodeInfo() (*p2p.NodeInfo, error) {
	server := sb.subBridge.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.NodeInfo(), nil
}

func (sb *SubBridgeAPI) GetParentOperatorAddr() common.Address {
	return sb.subBridge.bridgeAccounts.pAccount.address
}

func (sb *SubBridgeAPI) GetChildOperatorAddr() common.Address {
	return sb.subBridge.bridgeAccounts.cAccount.address
}

func (sb *SubBridgeAPI) GetParentOperatorNonce() uint64 {
	return sb.subBridge.handler.getParentOperatorNonce()
}

func (sb *SubBridgeAPI) GetChildOperatorNonce() uint64 {
	return sb.subBridge.handler.getChildOperatorNonce()
}

func (sb *SubBridgeAPI) GetParentOperatorBalance() (*big.Int, error) {
	return sb.subBridge.handler.getParentOperatorBalance()
}

// GetParentBridgeContractBalance returns the balance of the bridge contract in the parent chain.
func (sb *SubBridgeAPI) GetParentBridgeContractBalance(addr common.Address) (*big.Int, error) {
	return sb.subBridge.handler.getParentBridgeContractBalance(addr)
}

func (sb *SubBridgeAPI) GetChildOperatorBalance() (*big.Int, error) {
	return sb.subBridge.handler.getChildOperatorBalance()
}

// GetChildBridgeContractBalance returns the balance of the bridge contract in the child chain.
func (sb *SubBridgeAPI) GetChildBridgeContractBalance(addr common.Address) (*big.Int, error) {
	return sb.subBridge.handler.getChildBridgeContractBalance(addr)
}

// GetOperators returns the information of bridge operators.
func (sb *SubBridgeAPI) GetOperators() map[string]interface{} {
	return sb.subBridge.bridgeAccounts.GetBridgeOperators()
}

// LockParentOperator can lock the parent bridge operator.
func (sb *SubBridgeAPI) LockParentOperator() error {
	return sb.subBridge.bridgeAccounts.pAccount.LockAccount()
}

// LockChildOperator can lock the child bridge operator.
func (sb *SubBridgeAPI) LockChildOperator() error {
	return sb.subBridge.bridgeAccounts.cAccount.LockAccount()
}

// UnlockParentOperator can unlock the parent bridge operator.
func (sb *SubBridgeAPI) UnlockParentOperator(passphrase string, duration *uint64) error {
	return sb.subBridge.bridgeAccounts.pAccount.UnLockAccount(passphrase, duration)
}

// UnlockChildOperator can unlock the child bridge operator.
func (sb *SubBridgeAPI) UnlockChildOperator(passphrase string, duration *uint64) error {
	return sb.subBridge.bridgeAccounts.cAccount.UnLockAccount(passphrase, duration)
}

func (sb *SubBridgeAPI) GetAnchoringPeriod() uint64 {
	return sb.subBridge.config.AnchoringPeriod
}

func (sb *SubBridgeAPI) GetSentChainTxsLimit() uint64 {
	return sb.subBridge.config.SentChainTxsLimit
}

// SetParentOperatorFeePayer can set the parent bridge operator's fee payer.
func (sb *SubBridgeAPI) SetParentOperatorFeePayer(feePayer common.Address) error {
	return sb.subBridge.bridgeAccounts.SetParentOperatorFeePayer(feePayer)
}

// SetChildOperatorFeePayer can set the child bridge operator's fee payer.
func (sb *SubBridgeAPI) SetChildOperatorFeePayer(feePayer common.Address) error {
	return sb.subBridge.bridgeAccounts.SetChildOperatorFeePayer(feePayer)
}

// SetParentBridgeOperatorGasLimit changes value of bridge parent operator's gaslimit.
func (sb *SubBridgeAPI) SetParentBridgeOperatorGasLimit(fee uint64) {
	sb.subBridge.bridgeAccounts.SetParentBridgeOperatorGasLimit(fee)
}

// SetChildBridgeOperatorGasLimit changes value of bridge child operator's gaslimit.
func (sb *SubBridgeAPI) SetChildBridgeOperatorGasLimit(fee uint64) {
	sb.subBridge.bridgeAccounts.SetChildBridgeOperatorGasLimit(fee)
}

// GetParentOperatorFeePayer can return the parent bridge operator's fee payer.
func (sb *SubBridgeAPI) GetParentOperatorFeePayer() common.Address {
	return sb.subBridge.bridgeAccounts.GetParentOperatorFeePayer()
}

// GetChildOperatorFeePayer can return the child bridge operator's fee payer.
func (sb *SubBridgeAPI) GetChildOperatorFeePayer() common.Address {
	return sb.subBridge.bridgeAccounts.GetChildOperatorFeePayer()
}

// GetParentBridgeOperatorGasLimit gets value of bridge parent operator's gaslimit.
func (sb *SubBridgeAPI) GetParentBridgeOperatorGasLimit() uint64 {
	return sb.subBridge.bridgeAccounts.GetParentBridgeOperatorGasLimit()
}

// GetChildBridgeOperatorGasLimit gets value of bridge child operator's gaslimit.
func (sb *SubBridgeAPI) GetChildBridgeOperatorGasLimit() uint64 {
	return sb.subBridge.bridgeAccounts.GetChildBridgeOperatorGasLimit()
}
