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
	"fmt"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/node"
	"github.com/pkg/errors"
)

var (
	ErrInvalidBridgePair = errors.New("invalid bridge pair")
)

// MainBridgeAPI Implementation for main-bridge node
type MainBridgeAPI struct {
	sc *MainBridge
}

func (mbapi *MainBridgeAPI) GetChildChainIndexingEnabled() bool {
	return mbapi.sc.eventhandler.GetChildChainIndexingEnabled()
}

func (mbapi *MainBridgeAPI) ConvertServiceChainBlockHashToMainChainTxHash(scBlockHash common.Hash) common.Hash {
	return mbapi.sc.eventhandler.ConvertServiceChainBlockHashToMainChainTxHash(scBlockHash)
}

// Peers retrieves all the information we know about each individual peer at the
// protocol granularity.
func (mbapi *MainBridgeAPI) Peers() ([]*p2p.PeerInfo, error) {
	server := mbapi.sc.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.PeersInfo(), nil
}

// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.
func (mbapi *MainBridgeAPI) NodeInfo() (*p2p.NodeInfo, error) {
	server := mbapi.sc.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.NodeInfo(), nil
}

// SubBridgeAPI Implementation for sub-bridge node
type SubBridgeAPI struct {
	sc *SubBridge
}

func (sbapi *SubBridgeAPI) ConvertServiceChainBlockHashToMainChainTxHash(scBlockHash common.Hash) common.Hash {
	return sbapi.sc.eventhandler.ConvertServiceChainBlockHashToMainChainTxHash(scBlockHash)
}

func (sbapi *SubBridgeAPI) GetLatestAnchoredBlockNumber() uint64 {
	return sbapi.sc.handler.GetLatestAnchoredBlockNumber()
}

func (sbapi *SubBridgeAPI) GetReceiptFromParentChain(blockHash common.Hash) *types.Receipt {
	return sbapi.sc.handler.GetReceiptFromParentChain(blockHash)
}

func (sbapi *SubBridgeAPI) DeployBridge() ([]common.Address, error) {
	cBridge, cBridgeAddr, err := sbapi.sc.bridgeManager.DeployBridge(sbapi.sc.localBackend, true)
	if err != nil {
		logger.Error("Failed to deploy service chain bridge.", "err", err)
		return nil, err
	}
	pBridge, pBridgeAddr, err := sbapi.sc.bridgeManager.DeployBridge(sbapi.sc.remoteBackend, false)
	if err != nil {
		logger.Error("Failed to deploy main chain bridge.", "err", err)
		return nil, err
	}

	pAcc := sbapi.sc.bridgeAccountManager.mcAccount
	cAcc := sbapi.sc.bridgeAccountManager.scAccount

	err = sbapi.sc.bridgeManager.SetBridgeInfo(cBridgeAddr, cBridge, pBridgeAddr, pBridge, cAcc, true, false)
	if err != nil {
		return nil, err
	}

	err = sbapi.sc.bridgeManager.SetBridgeInfo(pBridgeAddr, pBridge, cBridgeAddr, cBridge, pAcc, false, false)
	if err != nil {
		return nil, err
	}

	err = sbapi.sc.bridgeManager.SetJournal(cBridgeAddr, pBridgeAddr)
	if err != nil {
		return nil, err
	}

	return []common.Address{cBridgeAddr, pBridgeAddr}, nil
}

// SubscribeBridge enables the given service/main chain bridges to subscribe the events.
func (sbapi *SubBridgeAPI) SubscribeBridge(cBridgeAddr, pBridgeAddr common.Address) error {
	if !sbapi.sc.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	err := sbapi.sc.bridgeManager.SubscribeEvent(cBridgeAddr)
	if err != nil {
		logger.Error("Failed to SubscribeEvent child bridge", "addr", cBridgeAddr, "err", err)
		return err
	}

	err = sbapi.sc.bridgeManager.SubscribeEvent(pBridgeAddr)
	if err != nil {
		logger.Error("Failed to SubscribeEvent parent bridge", "addr", pBridgeAddr, "err", err)
		sbapi.sc.bridgeManager.UnsubscribeEvent(cBridgeAddr)
		return err
	}

	sbapi.sc.bridgeManager.journal.cache[cBridgeAddr].Subscribed = true

	// Update the journal's subscribed flag.
	sbapi.sc.bridgeManager.journal.rotate(sbapi.sc.bridgeManager.GetAllBridge())

	err = sbapi.sc.bridgeManager.AddRecovery(cBridgeAddr, pBridgeAddr)
	if err != nil {
		sbapi.sc.bridgeManager.UnsubscribeEvent(cBridgeAddr)
		sbapi.sc.bridgeManager.UnsubscribeEvent(pBridgeAddr)
		return err
	}
	return nil
}

// UnsubscribeBridge disables the event subscription of the given service/main chain bridges.
func (sbapi *SubBridgeAPI) UnsubscribeBridge(cBridgeAddr, pBridgeAddr common.Address) error {
	if !sbapi.sc.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	sbapi.sc.bridgeManager.UnsubscribeEvent(cBridgeAddr)
	sbapi.sc.bridgeManager.UnsubscribeEvent(pBridgeAddr)

	sbapi.sc.bridgeManager.journal.cache[cBridgeAddr].Subscribed = false

	sbapi.sc.bridgeManager.journal.rotate(sbapi.sc.bridgeManager.GetAllBridge())
	return nil
}

func (sbapi *SubBridgeAPI) ConvertRequestTxHashToHandleTxHash(hash common.Hash) common.Hash {
	return sbapi.sc.chainDB.ReadHandleTxHashFromRequestTxHash(hash)
}

func (sbapi *SubBridgeAPI) TxPendingCount() int {
	return sbapi.sc.GetBridgeTxPool().Stats()
}

func (sbapi *SubBridgeAPI) TxPending() map[common.Address]types.Transactions {
	return sbapi.sc.GetBridgeTxPool().Pending()
}

func (sbapi *SubBridgeAPI) ListBridge() []*BridgeJournal {
	return sbapi.sc.bridgeManager.GetAllBridge()
}

func (sbapi *SubBridgeAPI) GetBridgeInformation(bridgeAddr common.Address) (map[string]interface{}, error) {
	if ctBridge := sbapi.sc.bridgeManager.GetCounterPartBridgeAddr(bridgeAddr); ctBridge == (common.Address{}) {
		return nil, ErrInvalidBridgePair
	}

	bi, ok := sbapi.sc.bridgeManager.GetBridgeInfo(bridgeAddr)
	if !ok {
		return nil, ErrNoBridgeInfo
	}

	bi.UpdateInfo()

	return map[string]interface{}{
		"isRunning":        bi.isRunning,
		"requestNonce":     bi.requestNonceFromCounterPart,
		"handleNonce":      bi.handleNonce,
		"counterPart":      bi.counterpartAddress,
		"onServiceChain":   bi.onServiceChain,
		"isSubscribed":     bi.subscribed,
		"pendingEventSize": bi.pendingRequestEvent.Len(),
	}, nil
}

func (sbapi *SubBridgeAPI) Anchoring(flag bool) bool {
	return sbapi.sc.SetAnchoringTx(flag)
}

func (sbapi *SubBridgeAPI) GetAnchoring() bool {
	return sbapi.sc.GetAnchoringTx()
}

func (sbapi *SubBridgeAPI) RegisterBridge(cBridgeAddr common.Address, pBridgeAddr common.Address) error {
	cBridge, err := bridge.NewBridge(cBridgeAddr, sbapi.sc.localBackend)
	if err != nil {
		return err
	}
	pBridge, err := bridge.NewBridge(pBridgeAddr, sbapi.sc.remoteBackend)
	if err != nil {
		return err
	}

	bm := sbapi.sc.bridgeManager
	err = bm.SetBridgeInfo(cBridgeAddr, cBridge, pBridgeAddr, pBridge, sbapi.sc.bridgeAccountManager.scAccount, true, false)
	if err != nil {
		return err
	}
	err = bm.SetBridgeInfo(pBridgeAddr, pBridge, cBridgeAddr, cBridge, sbapi.sc.bridgeAccountManager.mcAccount, false, false)
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

func (sbapi *SubBridgeAPI) DeregisterBridge(cBridgeAddr common.Address, pBridgeAddr common.Address) error {
	if !sbapi.sc.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	bm := sbapi.sc.bridgeManager
	journal := bm.journal.cache[cBridgeAddr]

	if journal.Subscribed {
		bm.UnsubscribeEvent(journal.LocalAddress)
		bm.UnsubscribeEvent(journal.RemoteAddress)

		bm.DeleteRecovery(cBridgeAddr, pBridgeAddr)
	}

	delete(bm.journal.cache, cBridgeAddr)

	if err := bm.journal.rotate(bm.GetAllBridge()); err != nil {
		logger.Warn("failed to rotate bridge journal", "err", err, "scBridge", cBridgeAddr.String(), "mcBridge", pBridgeAddr.String())
	}

	if err := bm.DeleteBridgeInfo(cBridgeAddr); err != nil {
		logger.Warn("failed to Delete service chain bridge info", "err", err, "bridge", cBridgeAddr.String())
	}

	if err := bm.DeleteBridgeInfo(pBridgeAddr); err != nil {
		logger.Warn("failed to Delete main chain bridge info", "err", err, "bridge", pBridgeAddr.String())
	}

	return nil
}

func (sbapi *SubBridgeAPI) RegisterToken(cBridgeAddr, pBridgeAddr, cTokenAddr, pTokenAddr common.Address) error {
	if !sbapi.sc.bridgeManager.IsValidBridgePair(cBridgeAddr, pBridgeAddr) {
		return ErrInvalidBridgePair
	}

	cBi, cExist := sbapi.sc.bridgeManager.GetBridgeInfo(cBridgeAddr)
	pBi, pExist := sbapi.sc.bridgeManager.GetBridgeInfo(pBridgeAddr)

	if !cExist || !pExist {
		return errors.New("bridge does not exist")
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
	tx, err := cBi.bridge.RegisterToken(cBi.account.GetTransactOpts(), cTokenAddr, pTokenAddr)
	if err != nil {
		cBi.account.UnLock()
		return err
	}
	cBi.account.IncNonce()
	cBi.account.UnLock()
	logger.Debug("scBridge registered token", "txHash", tx.Hash().String(), "scToken", cTokenAddr.String(), "mcToken", pTokenAddr.String())

	pBi.account.Lock()
	tx, err = pBi.bridge.RegisterToken(pBi.account.GetTransactOpts(), pTokenAddr, cTokenAddr)
	if err != nil {
		pBi.account.UnLock()
		return err
	}
	pBi.account.IncNonce()
	pBi.account.UnLock()
	logger.Debug("mcBridge registered token", "txHash", tx.Hash().String(), "scToken", cTokenAddr.String(), "mcToken", pTokenAddr.String())

	logger.Info("Register token", "scToken", cTokenAddr.String(), "mcToken", pTokenAddr.String())
	return nil
}

func (sbapi *SubBridgeAPI) DeregisterToken(cBridgeAddr, pBridgeAddr, cTokenAddr, pTokenAddr common.Address) error {
	cBi, cExist := sbapi.sc.bridgeManager.GetBridgeInfo(cBridgeAddr)
	pBi, pExist := sbapi.sc.bridgeManager.GetBridgeInfo(pBridgeAddr)

	if !cExist || !pExist {
		return errors.New("bridge does not exist")
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
	tx, err := cBi.bridge.DeregisterToken(cBi.account.GetTransactOpts(), cTokenAddr)
	if err != nil {
		return err
	}
	cBi.account.IncNonce()
	logger.Debug("scBridge deregistered token", "txHash", tx.Hash().String(), "scToken", cTokenAddr.String(), "mcToken", pTokenAddr.String())

	pBi.account.Lock()
	defer pBi.account.UnLock()
	tx, err = pBi.bridge.DeregisterToken(pBi.account.GetTransactOpts(), pTokenAddr)
	if err != nil {
		return err
	}
	pBi.account.IncNonce()
	logger.Debug("mcBridge deregistered token", "txHash", tx.Hash().String(), "scToken", cTokenAddr.String(), "mcToken", pTokenAddr.String())
	return err
}

// AddPeer requests connecting to a remote node, and also maintaining the new
// connection at all times, even reconnecting if it is lost.
func (sbapi *SubBridgeAPI) AddPeer(url string) (bool, error) {
	// Make sure the server is running, fail otherwise
	server := sbapi.sc.bridgeServer
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
	server.AddPeer(node, false)
	return node, nil
}

// RemovePeer disconnects from a a remote node if the connection exists
func (sbapi *SubBridgeAPI) RemovePeer(url string) (bool, error) {
	// Make sure the server is running, fail otherwise
	server := sbapi.sc.bridgeServer
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
func (sbapi *SubBridgeAPI) Peers() ([]*p2p.PeerInfo, error) {
	server := sbapi.sc.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.PeersInfo(), nil
}

// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.
func (sbapi *SubBridgeAPI) NodeInfo() (*p2p.NodeInfo, error) {
	server := sbapi.sc.bridgeServer
	if server == nil {
		return nil, node.ErrNodeStopped
	}
	return server.NodeInfo(), nil
}

func (sbapi *SubBridgeAPI) GetMainChainAccountAddr() string {
	return sbapi.sc.config.MainChainAccountAddr.Hex()
}

func (sbapi *SubBridgeAPI) GetServiceChainAccountAddr() string {
	return sbapi.sc.config.ServiceChainAccountAddr.Hex()
}

func (sbapi *SubBridgeAPI) GetMainChainAccountNonce() uint64 {
	return sbapi.sc.handler.getMainChainAccountNonce()
}

func (sbapi *SubBridgeAPI) GetServiceChainAccountNonce() uint64 {
	return sbapi.sc.handler.getServiceChainAccountNonce()
}

func (sbapi *SubBridgeAPI) GetAnchoringPeriod() uint64 {
	return sbapi.sc.config.AnchoringPeriod
}

func (sbapi *SubBridgeAPI) GetSentChainTxsLimit() uint64 {
	return sbapi.sc.config.SentChainTxsLimit
}
