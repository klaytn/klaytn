// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from rpc/json.go (2018/06/04).
// Modified and improved for the klaytn development.

package grpc

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type gKlaytnClient struct {
	addr   string
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn
}

// json message
type jsonRequest struct {
	Method  string          `json:"method"`
	Version string          `json:"jsonrpc"`
	Id      json.RawMessage `json:"id,omitempty"`
	Payload json.RawMessage `json:"params,omitempty"`
}

type jsonSuccessResponse struct {
	Version string      `json:"jsonrpc"`
	Id      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result"`
}

type jsonSubscription struct {
	Subscription string      `json:"subscription"`
	Result       interface{} `json:"result,omitempty"`
}

func NewgKlaytnClient(addr string) (*gKlaytnClient, error) {
	return &gKlaytnClient{addr: addr}, nil
}

const timeout = 5 * time.Minute

func (gkc *gKlaytnClient) makeKlaytnClient(timeout time.Duration) (KlaytnNodeClient, error) {
	gkc.ctx, gkc.cancel = context.WithTimeout(context.Background(), timeout)
	conn, err := grpc.DialContext(gkc.ctx, gkc.addr, grpc.WithInsecure())
	if err != nil {
		logger.Error("failed to dial server", "err", err)
		return nil, err
	}
	gkc.conn = conn

	return NewKlaytnNodeClient(gkc.conn), nil
}

func (gkc *gKlaytnClient) makeRPCRequest(service string, method string, args []interface{}) (*RPCRequest, error) {
	payload, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	id, err := json.Marshal(1)
	if err != nil {
		return nil, err
	}

	arguments := &jsonRequest{method, "2.0", id, payload}
	params, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	return &RPCRequest{Service: service, Method: method, Params: params}, nil
}

func (gkc *gKlaytnClient) handleRPCResponse(response *RPCResponse) error {
	var out jsonSuccessResponse
	if err := json.Unmarshal(response.Payload, &out); err != nil {
		logger.Error("failed to handle response", "err", err)
		return err
	}

	//fmt.Println(out.Result)
	return nil
}

func (gkc *gKlaytnClient) handleSubscribe(client KlaytnNode_SubscribeClient, handle func(response *RPCResponse) error) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	ticker := time.NewTicker(1 * time.Second)

loop:
	for {
		select {
		case <-ticker.C:
			rev, err := client.Recv()
			if err == io.EOF {
				logger.Debug("close conn")
				waitGroup.Done()
				break loop
			}
			if rev != nil {
				if err := handle(rev); err != nil {
					logger.Warn("fail to handle", "err", err)
					waitGroup.Done()
					break loop
				}
			}
		}
	}

	waitGroup.Wait()
}

func (gkc *gKlaytnClient) handleBiCall(stream KlaytnNode_BiCallClient, request func() (*RPCRequest, error), handle func(response *RPCResponse) error) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		for {
			req, err := request()
			if err != nil {
				logger.Warn("fail to make request", "err", err)
				waitGroup.Done()
				return
			}
			if err := stream.Send(req); err != nil {
				logger.Warn("fail to send request", "err", err)
				waitGroup.Done()
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			var recv RPCResponse
			if err := stream.RecvMsg(&recv); err != nil {
				logger.Warn("fail to recv response", "err", err)
				waitGroup.Done()
			}

			if err := handle(&recv); err != nil {
				logger.Warn("fail to handle response", "err", err)
				waitGroup.Done()
			}
		}
	}()

	waitGroup.Wait()
}

func (gkc *gKlaytnClient) Close() {
	if gkc.cancel != nil {
		gkc.cancel()
	}

	if gkc.conn != nil {
		if err := gkc.conn.Close(); err != nil {
			logger.Warn("fail to close conn", "err", err)
		}
	}
}
