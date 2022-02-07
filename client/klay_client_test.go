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
	"fmt"
	fastws "github.com/clevergo/websocket"
	//gorillaws "github.com/gorilla/websocket"
	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
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

//func TestDialWebsocketAuth(t *testing.T) {
//	url := "wss://KASKZCTSDT07NI1PM54OKL85:nPFDFf1Qh3Zy5VfNmYwl3WV_Vq_R_Dmo3cBtncbP@node-api.klaytnapi.com/v1/ws/open?chain-id=1001"
//	c, err := Dial(url)
//	assert.Nil(t, err)
//
//	ch := make(chan *types.Header, 2)
//	sub, err := c.SubscribeNewHead(context.Background(), ch)
//	assert.Nil(t, err)
//
//	select {
//	case ev := <-ch:
//		log.Printf("New block header: %v", ev.Number)
//	case err := <-sub.Err():
//		log.Printf("Error while subcription: %v", err)
//	}
//}

func TestDialWebsocketAuth(t *testing.T) {
	url := "wss://KASKZCTSDT07NI1PM54OKL85:nPFDFf1Qh3Zy5VfNmYwl3WV_Vq_R_Dmo3cBtncbP@node-api.klaytnapi.com/v1/ws/open?chain-id=1001"
	url = "wss://node-api.klaytnapi.com/v1/ws/open?chain-id=1001"
	auth := "Basic S0FTS1pDVFNEVDA3TkkxUE01NE9LTDg1Om5QRkRGZjFRaDNaeTVWZk5tWXdsM1dWX1ZxX1JfRG1vM2NCdG5jYlA="

	_ = auth
	/*		header := http.Header(make(map[string][]string))
			header.Add("Authorization", auth)
			ctx := context.Background()
			dialer := goril.Dialer{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			}
			conn, resp, err := dialer.DialContext(ctx, url, header)
			fmt.Println(conn)
			fmt.Println(resp)
			fmt.Println(err)
	*/
	if 1 == 1 {

		// {"jsonrpc":"2.0", "id": 1, "method": "klay_subscribe", "params": ["newHeads"]}
		ctx := context.Background()
		_ = ctx

		dialer := fastws.Dialer{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		header := http.Header(make(map[string][]string))
		header.Add("Authorization", auth)

		conn, resp, err := dialer.Dial(url, header)
		fmt.Println(conn)
		fmt.Println(resp)
		fmt.Println(err)

		err = conn.WriteJSON(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "klay_subscribe",
			"params":  []string{"newHeads"},
		})
		fmt.Println("write err", err)

		//		time.Sleep(time.Second * 2)

		//		var r interface{}
		//		conn.ReadJSON(r)
		//		fmt.Println("resp", r)

		_, msg, _ := conn.ReadMessage()
		fmt.Println("read bytes", string(msg))
		_, msg, _ = conn.ReadMessage()
		fmt.Println("read bytes", string(msg))
		_, msg, _ = conn.ReadMessage()
		fmt.Println("read bytes", string(msg))
		_, msg, _ = conn.ReadMessage()
		fmt.Println("read bytes", string(msg))
		_, msg, _ = conn.ReadMessage()
		fmt.Println("read bytes", string(msg))

	}

	if 1 == 0 {

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
}
