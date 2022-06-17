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
	"math/big"
	"net"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/pkg/errors"
)

var (
	NoParentPeerErr = errors.New("no parent peer")
)

const timeout = 30 * time.Second

// TODO-Klaytn currently RemoteBackend is only for ServiceChain, especially Bridge SmartContract
type RemoteBackend struct {
	subBridge *SubBridge

	rpcClient *rpc.Client
	chainID   *big.Int
}

func NewRpcClientP2P(sb *SubBridge) *rpc.Client {
	initctx := context.Background()
	c, _ := rpc.NewClient(initctx, func(context.Context) (net.Conn, error) {
		p1, p2 := net.Pipe()
		sb.SetRPCConn(p1)
		return p2, nil
	})
	return c
}

func NewRemoteBackend(sb *SubBridge) (*RemoteBackend, error) {
	rCli := NewRpcClientP2P(sb)

	return &RemoteBackend{
		subBridge: sb,
		rpcClient: rCli,
	}, nil
}

func (rb *RemoteBackend) checkParentPeer() bool {
	return rb.subBridge.peers.Len() > 0
}

func (rb *RemoteBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	var result hexutil.Bytes
	err := rb.rpcClient.CallContext(ctx, &result, "klay_getCode", contract, toBlockNumArg(blockNumber))
	return result, err
}

func (rb *RemoteBackend) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	var hex hexutil.Big
	err := rb.rpcClient.CallContext(ctx, &hex, "klay_getBalance", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (rb *RemoteBackend) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	var hex hexutil.Bytes
	err := rb.rpcClient.CallContext(ctx, &hex, "klay_call", toCallArg(call), toBlockNumArg(blockNumber))
	return hex, err
}

func (rb *RemoteBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	var result hexutil.Bytes
	err := rb.rpcClient.CallContext(ctx, &result, "klay_getCode", contract, "pending")
	return result, err
}

func (rb *RemoteBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if !rb.checkParentPeer() {
		return 0, NoParentPeerErr
	}
	var result hexutil.Uint64
	err := rb.rpcClient.CallContext(ctx, &result, "klay_getTransactionCount", account, "pending")
	return uint64(result), err
}

func (rb *RemoteBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	var hex hexutil.Big
	if err := rb.rpcClient.CallContext(ctx, &hex, "klay_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (rb *RemoteBackend) EstimateGas(ctx context.Context, msg klaytn.CallMsg) (uint64, error) {
	if !rb.checkParentPeer() {
		return 0, NoParentPeerErr
	}

	var hex hexutil.Uint64
	err := rb.rpcClient.CallContext(ctx, &hex, "klay_estimateGas", toCallArg(msg))
	if err != nil {
		return 0, err
	}
	return uint64(hex), nil
}

func (rb *RemoteBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if !rb.checkParentPeer() {
		return NoParentPeerErr
	}
	return rb.subBridge.bridgeTxPool.AddLocal(tx)
}

func (rb *RemoteBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	var r *types.Receipt
	err := rb.rpcClient.CallContext(ctx, &r, "klay_getTransactionReceipt", txHash)
	if err == nil && r == nil {
		return nil, klaytn.NotFound
	}
	return r, err
}

func (rb *RemoteBackend) TransactionReceiptRpcOutput(ctx context.Context, txHash common.Hash) (r map[string]interface{}, err error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}

	err = rb.rpcClient.CallContext(ctx, &r, "klay_getTransactionReceipt", txHash)
	if err == nil && r == nil {
		return nil, klaytn.NotFound
	}
	return
}

// ChainID returns the chain ID of the sub-bridge configuration.
func (rb *RemoteBackend) ChainID(ctx context.Context) (*big.Int, error) {
	return big.NewInt(int64(rb.subBridge.config.ParentChainID)), nil
}

func (rb *RemoteBackend) FilterLogs(ctx context.Context, query klaytn.FilterQuery) (result []types.Log, err error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	err = rb.rpcClient.CallContext(ctx, &result, "klay_getLogs", toFilterArg(query))
	return
}

func (rb *RemoteBackend) SubscribeFilterLogs(ctx context.Context, query klaytn.FilterQuery, ch chan<- types.Log) (klaytn.Subscription, error) {
	if !rb.checkParentPeer() {
		return nil, NoParentPeerErr
	}
	return rb.rpcClient.KlaySubscribe(ctx, ch, "logs", toFilterArg(query))
}

// CurrentBlockNumber returns a current block number.
func (rb *RemoteBackend) CurrentBlockNumber(ctx context.Context) (uint64, error) {
	if !rb.checkParentPeer() {
		return 0, NoParentPeerErr
	}
	var result hexutil.Uint64
	err := rb.rpcClient.CallContext(ctx, &result, "klay_blockNumber")
	return uint64(result), err
}

func toFilterArg(q klaytn.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

func toCallArg(msg klaytn.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}
