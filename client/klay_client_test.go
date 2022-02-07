// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from ethclient/ethclient_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package client

import (
	"context"
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// Verify that Client implements the Klaytn interfaces.
var (
	// _ = klaytn.Subscription(&Client{})
	_ = klaytn.ChainReader(&Client{})
	_ = klaytn.TransactionReader(&Client{})
	_ = klaytn.ChainStateReader(&Client{})
	_ = klaytn.ChainSyncReader(&Client{})
	_ = klaytn.ContractCaller(&Client{})
	_ = klaytn.LogFilterer(&Client{})
	_ = klaytn.TransactionSender(&Client{})
	_ = klaytn.GasPricer(&Client{})
	_ = klaytn.PendingStateReader(&Client{})
	_ = klaytn.PendingContractCaller(&Client{})
	_ = klaytn.GasEstimator(&Client{})
	// _ = klaytn.PendingStateEventer(&Client{})
)

func TestDialWebsocketAuth(t *testing.T) {
	url := "wss://KASKZCTSDT07NI1PM54OKL85:nPFDFf1Qh3Zy5VfNmYwl3WV_Vq_R_Dmo3cBtncbP@node-api.klaytnapi.com/v1/ws/open?chain-id=1001"
	c, err := Dial(url)
	assert.Nil(t, err)

	ch := make(chan *types.Header, 2)
	sub, err := c.SubscribeNewHead(context.Background(), ch)
	assert.Nil(t, err)

	select {
	case ev := <-ch:
		log.Printf("New block header: %v", ev.Number)
	case err := <-sub.Err():
		log.Printf("Error while subcription: %v", err)
	}
}
