// Modifications Copyright 2019 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethclient/ethclient.go (2018/06/04).
// Modified and improved for the klaytn development.

package client

import (
	"context"
	"errors"
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/p2p"
	"math/big"
)

// BridgeAddPeerOnParentChain can add a static peer on bridge node for service chain.
func (ec *Client) BridgeAddPeerOnBridge(ctx context.Context, url string) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "subbridge_addPeer", url)
	return result, err
}

// BridgeRemovePeerOnParentChain can remove a static peer on bridge node.
func (ec *Client) BridgeRemovePeerOnBridge(ctx context.Context, url string) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "subbridge_removePeer", url)
	return result, err
}

// BridgePeersOnBridge returns the peer list of bridge node for service chain.
func (ec *Client) BridgePeersOnBridge(ctx context.Context) ([]*p2p.PeerInfo, error) {
	var result []*p2p.PeerInfo
	err := ec.c.CallContext(ctx, &result, "subbridge_peers")
	return result, err
}

// BridgeNodeInfo returns the node information
func (ec *Client) BridgeNodeInfo(ctx context.Context) (*p2p.NodeInfo, error) {
	var result p2p.NodeInfo
	err := ec.c.CallContext(ctx, &result, "subbridge_nodeInfo")
	return &result, err
}

// BridgePeersOnBridge returns the peer list of bridge node for service chain.
func (ec *Client) MainBridgePeersOnBridge(ctx context.Context) ([]*p2p.PeerInfo, error) {
	var result []*p2p.PeerInfo
	err := ec.c.CallContext(ctx, &result, "mainbridge_peers")
	return result, err
}

// BridgeNodeInfo returns the node information
func (ec *Client) MainBridgeNodeInfo(ctx context.Context) (*p2p.NodeInfo, error) {
	var result p2p.NodeInfo
	err := ec.c.CallContext(ctx, &result, "mainbridge_nodeInfo")
	return &result, err
}

// BridgeGetChildChainIndexingEnabled can get if child chain indexing is enabled or not.
func (ec *Client) BridgeGetChildChainIndexingEnabled(ctx context.Context) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "mainbridge_getChildChainIndexingEnabled")
	return result, err
}

// BridgeConvertServiceChainBlockHashToMainChainTxHash can convert service chain block hash to
// anchoring tx hash which contains anchored data.
func (ec *Client) BridgeConvertServiceChainBlockHashToMainChainTxHash(ctx context.Context, scBlockHash common.Hash) (common.Hash, error) {
	var txHash common.Hash
	err := ec.c.CallContext(ctx, &txHash, "mainbridge_convertServiceChainBlockHashToMainChainTxHash", scBlockHash)
	return txHash, err
}

// BridgeConvertRequestTxHashToHandleTxHash can convert a request value transfer tx hash to
// the corresponded handle value transfer tx hash.
func (ec *Client) BridgeConvertRequestTxHashToHandleTxHash(ctx context.Context, requestTxHash common.Hash) (common.Hash, error) {
	var handleTxHash common.Hash
	err := ec.c.CallContext(ctx, &handleTxHash, "subbridge_convertRequestTxHashToHandleTxHash", requestTxHash)
	return handleTxHash, err
}

// BridgeGetReceiptFromParentChain can get the receipt of child chain tx from parent node.
func (ec *Client) BridgeGetReceiptFromParentChain(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	var result *types.Receipt
	err := ec.c.CallContext(ctx, &result, "subbridge_getReceiptFromParentChain", hash)
	if err == nil && result == nil {
		return nil, klaytn.NotFound
	}
	return result, err
}

// BridgeGetMainChainAccountAddr can get a main chain bridge account address.
func (ec *Client) BridgeGetMainChainAccountAddr(ctx context.Context) (common.Address, error) {
	var result common.Address
	err := ec.c.CallContext(ctx, &result, "subbridge_getMainChainAccountAddr")
	return result, err
}

// BridgeGetServiceChainAccountAddr can get a service chain bridge account address.
func (ec *Client) BridgeGetServiceChainAccountAddr(ctx context.Context) (common.Address, error) {
	var result common.Address
	err := ec.c.CallContext(ctx, &result, "subbridge_getServiceChainAccountAddr")
	return result, err
}

// BridgeGetMainChainAccountNonce can get a main chain bridge account nonce.
func (ec *Client) BridgeGetMainChainAccountNonce(ctx context.Context) (uint64, error) {
	var result uint64
	err := ec.c.CallContext(ctx, &result, "subbridge_getMainChainAccountNonce")
	return result, err
}

// BridgeGetServiceChainAccountAddr can get a service chain bridge account nonce.
func (ec *Client) BridgeGetServiceChainAccountNonce(ctx context.Context) (uint64, error) {
	var result uint64
	err := ec.c.CallContext(ctx, &result, "subbridge_getServiceChainAccountNonce")
	return result, err
}

// BridgeGetLatestAnchoredBlockNumber can return the latest anchored block number.
func (ec *Client) BridgeGetLatestAnchoredBlockNumber(ctx context.Context) (uint64, error) {
	var result uint64
	err := ec.c.CallContext(ctx, &result, "subbridge_getLatestAnchoredBlockNumber")
	return result, err
}

// BridgeEnableAnchoring can enable anchoring function and return the set value.
func (ec *Client) BridgeEnableAnchoring(ctx context.Context) (bool, error) {
	return ec.setAnchoring(ctx, true)
}

// BridgeDisableAnchoring can disable anchoring function and return the set value.
func (ec *Client) BridgeDisableAnchoring(ctx context.Context) (bool, error) {
	return ec.setAnchoring(ctx, false)
}

// setAnchoring can set if anchoring is enabled or not.
func (ec *Client) setAnchoring(ctx context.Context, enable bool) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "subbridge_anchoring", enable)
	return result, err
}

// BridgeGetAnchoringPeriod can get the block period to anchor chain data.
func (ec *Client) BridgeGetAnchoringPeriod(ctx context.Context) (uint64, error) {
	var result uint64
	err := ec.c.CallContext(ctx, &result, "subbridge_getAnchoringPeriod")
	return result, err
}

// BridgeGetSentChainTxsLimit can get the maximum number of transaction which child peer can send to parent peer once.
func (ec *Client) BridgeGetSentChainTxsLimit(ctx context.Context) (uint64, error) {
	var result uint64
	err := ec.c.CallContext(ctx, &result, "subbridge_getSentChainTxsLimit")
	return result, err
}

// BridgeDeployBridge can deploy the pair of bridge for parent/child chain.
func (ec *Client) BridgeDeployBridge(ctx context.Context) (common.Address, common.Address, error) {
	var result []common.Address

	err := ec.c.CallContext(ctx, &result, "subbridge_deployBridge")
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	if len(result) != 2 {
		return common.Address{}, common.Address{}, errors.New("output arguments length err")
	}

	return result[0], result[1], nil
}

// BridgeRegisterBridge can register the given pair of deployed child/parent bridges.
func (ec *Client) BridgeRegisterBridge(ctx context.Context, scBridge common.Address, mcBridge common.Address) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "subbridge_registerBridge", scBridge, mcBridge)
	return result, err
}

// BridgeDeregisterBridge can deregister the given pair of deployed child/parent bridges.
func (ec *Client) BridgeDeregisterBridge(ctx context.Context, scBridge common.Address, mcBridge common.Address) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "subbridge_deregisterBridge", scBridge, mcBridge)
	return result, err
}

// TODO-Klaytn if client pkg is removed in sc pkg, this will be replaced origin struct.
type BridgeJournal struct {
	LocalAddress  common.Address `json:"localAddress"`
	RemoteAddress common.Address `json:"remoteAddress"`
	Subscribed    bool           `json:"subscribed"`
}

// BridgeListBridge can return the list of the bridge.
func (ec *Client) BridgeListBridge(ctx context.Context) ([]*BridgeJournal, error) {
	var result []*BridgeJournal
	err := ec.c.CallContext(ctx, &result, "subbridge_listBridge")
	return result, err
}

// BridgeSubscribeBridge can enable for service chain bridge to subscribe the event of given service/main chain bridges.
// If the subscribing is failed, it returns an error.
func (ec *Client) BridgeSubscribeBridge(ctx context.Context, scBridge common.Address, mcBridge common.Address) error {
	return ec.c.CallContext(ctx, nil, "subbridge_subscribeBridge", scBridge, mcBridge)
}

// BridgeUnsubscribeBridge disables the event subscription of the given service/main chain bridges.
// If the unsubscribing is failed, it returns an error.
func (ec *Client) BridgeUnsubscribeBridge(ctx context.Context, scBridge common.Address, mcBridge common.Address) error {
	return ec.c.CallContext(ctx, nil, "subbridge_unsubscribeBridge", scBridge, mcBridge)
}

// BridgeRegisterTokenContract can register the given pair of deployed service/main chain token contracts.
// If the registering is failed, it returns an error.
func (ec *Client) BridgeRegisterTokenContract(ctx context.Context, scBridge, mcBridge, scToken, mcToken common.Address) error {
	return ec.c.CallContext(ctx, nil, "subbridge_registerToken", scBridge, mcBridge, scToken, mcToken)
}

// BridgeDeregisterTokenContract can deregister the given pair of deployed service/main chain token contracts.
// If the registering is failed, it returns an error.
func (ec *Client) BridgeDeregisterTokenContract(ctx context.Context, scBridge, mcBridge, scToken, mcToken common.Address) error {
	return ec.c.CallContext(ctx, nil, "subbridge_deregisterToken", scBridge, mcBridge, scToken, mcToken)
}

// BridgeTxPendingCount can return the count of the pend tx in bridge txpool.
func (ec *Client) BridgeTxPendingCount(ctx context.Context) (int, error) {
	var result int
	err := ec.c.CallContext(ctx, &result, "subbridge_txPendingCount")
	return result, err
}

// BridgeGetTxPending can return the pend tx list mapped by address.
func (ec *Client) BridgeGetTxPending(ctx context.Context) (map[common.Address]types.Transactions, error) {
	var result map[common.Address]types.Transactions
	err := ec.c.CallContext(ctx, &result, "subbridge_txPending")
	return result, err
}

// BridgeSetERC20Fee can set the ERC20 transfer fee.
func (ec *Client) BridgeSetERC20Fee(ctx context.Context, bridgeAddr, tokenAddr common.Address, fee *big.Int) (common.Hash, error) {
	var result common.Hash
	err := ec.c.CallContext(ctx, &result, "subbridge_setERC20Fee", bridgeAddr, tokenAddr, fee)
	return result, err
}

// BridgeSetKLAYFee can set the KLAY transfer fee.
func (ec *Client) BridgeSetKLAYFee(ctx context.Context, bridgeAddr common.Address, fee *big.Int) (common.Hash, error) {
	var result common.Hash
	err := ec.c.CallContext(ctx, &result, "subbridge_setKLAYFee", bridgeAddr, fee)
	return result, err
}

// BridgeGetERC20Fee returns the ERC20 transfer fee.
func (ec *Client) BridgeGetERC20Fee(ctx context.Context, bridgeAddr, tokenAddr common.Address) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "subbridge_getERC20Fee", bridgeAddr, tokenAddr)
	return (*big.Int)(&result), err
}

// BridgeGetKLAYFee returns the KLAY transfer fee.
func (ec *Client) BridgeGetKLAYFee(ctx context.Context, bridgeAddr common.Address) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "subbridge_getKLAYFee", bridgeAddr)
	return (*big.Int)(&result), err
}

// BridgeSetFeeReceiver can set the fee receiver.
func (ec *Client) BridgeSetFeeReceiver(ctx context.Context, bridgeAddr, receiver common.Address) (common.Hash, error) {
	var result common.Hash
	err := ec.c.CallContext(ctx, &result, "subbridge_setFeeReceiver", bridgeAddr, receiver)
	return result, err
}

// BridgeGetFeeReceiver returns the fee receiver.
func (ec *Client) BridgeGetFeeReceiver(ctx context.Context, bridgeAddr common.Address) (common.Address, error) {
	var result common.Address
	err := ec.c.CallContext(ctx, &result, "subbridge_getFeeReceiver", bridgeAddr)
	return result, err
}
