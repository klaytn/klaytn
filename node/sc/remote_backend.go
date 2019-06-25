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
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/pkg/errors"
	"math/big"
	"sync/atomic"
)

var (
	ConnectionFailErr = errors.New("fail to connect remote chain")
)

// TODO-Klaytn currently RemoteBackend is only for ServiceChain, especially Bridge SmartContract
type RemoteBackend struct {
	subBrige  *SubBridge
	targetUrl string

	klayClient *client.Client
}

func NewRemoteBackend(main *SubBridge, rawUrl string) (*RemoteBackend, error) {
	client, err := client.Dial(rawUrl)
	if err != nil {
		logger.Error("fail to connect RemoteChain", "url", rawUrl, "err", err)
		client = nil
	} else {
		logger.Info("success to connect RemoteChain", "url", rawUrl)
	}

	return &RemoteBackend{
		subBrige:   main,
		targetUrl:  rawUrl,
		klayClient: client,
	}, nil
}

func (rb *RemoteBackend) checkConnection() bool {
	if rb.klayClient == nil {
		logger.Error("klayclient is nil so try to reconnect")
		return rb.tryReconnect()
	}
	if atomic.CompareAndSwapInt64(&rb.subBrige.checkConnection, 1, 0) {
		rb.klayClient.Close()
		logger.Error("klayclient is disconnected so try to reconnect")
		connected := rb.tryReconnect()
		if !connected {
			atomic.StoreInt64(&rb.subBrige.checkConnection, 1)
		} else {
			rb.subBrige.bridgeManager.ResetAllSubscribedEvents()
		}
		return connected
	}
	return true
}

func (rb *RemoteBackend) tryReconnect() bool {
	client, err := client.Dial(rb.targetUrl)
	if err != nil {
		logger.Error("fail to reconnect RemoteChain", "url", rb.targetUrl, "err", err)
		return false
	}
	logger.Info("success to reconnect RemoteChain", "url", rb.targetUrl)

	rb.klayClient = client
	return true
}

func (rb *RemoteBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.CodeAt(ctx, contract, blockNumber)
}

func (rb *RemoteBackend) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.CallContract(ctx, call, blockNumber)
}

func (rb *RemoteBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.PendingCodeAt(ctx, contract)
}

func (rb *RemoteBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if !rb.checkConnection() {
		return 0, ConnectionFailErr
	}
	return rb.klayClient.PendingNonceAt(ctx, account)
}

func (rb *RemoteBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.SuggestGasPrice(ctx)
}

func (rb *RemoteBackend) EstimateGas(ctx context.Context, call klaytn.CallMsg) (gas uint64, err error) {
	if !rb.checkConnection() {
		return 0, ConnectionFailErr
	}
	return rb.klayClient.EstimateGas(ctx, call)
}

func (rb *RemoteBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if !rb.checkConnection() {
		return ConnectionFailErr
	}
	return rb.subBrige.bridgeTxPool.AddLocal(tx)
	// return rb.klayClient.SendTransaction(ctx, tx)
}

func (rb *RemoteBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.TransactionReceipt(ctx, txHash)
}

// ChainID can return the chain ID of the chain.
func (rb *RemoteBackend) ChainID(ctx context.Context) (*big.Int, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.ChainID(ctx)
}

func (rb *RemoteBackend) FilterLogs(ctx context.Context, query klaytn.FilterQuery) ([]types.Log, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.FilterLogs(ctx, query)
}

func (rb *RemoteBackend) SubscribeFilterLogs(ctx context.Context, query klaytn.FilterQuery, ch chan<- types.Log) (klaytn.Subscription, error) {
	if !rb.checkConnection() {
		return nil, ConnectionFailErr
	}
	return rb.klayClient.SubscribeFilterLogs(ctx, query, ch)
}
