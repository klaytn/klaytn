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

package grpc

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_BLOCK_NUMBER = float64(123456789)
)

type APIgRPC struct{}

func (a APIgRPC) BlockNumber() float64 {
	return TEST_BLOCK_NUMBER
}

func TestGRPC(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	addr := "127.0.0.1:4000"
	handler := rpc.NewServer()

	handler.RegisterName("klay", &APIgRPC{})

	listener := &Listener{Addr: addr}
	listener.SetRPCServer(handler)
	go listener.Start()

	time.Sleep(2 * time.Second)

	go testCall(t, addr, wg)
	go testBiCall(t, addr, wg)
	wg.Wait()
}
func testCall(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	kclient, _ := NewgKlaytnClient(addr)
	defer kclient.Close()

	knclient, err := kclient.makeKlaytnClient(timeout)
	assert.NoError(t, err)

	request, err := kclient.makeRPCRequest("klay", "klay_blockNumber", nil)
	assert.NoError(t, err)

	response, err := knclient.Call(kclient.ctx, request)
	assert.NoError(t, err)

	err = kclient.handleRPCResponse(response)
	assert.NoError(t, err)

	var out jsonSuccessResponse
	json.Unmarshal(response.Payload, &out)
	assert.NoError(t, err)
	assert.Equal(t, out.Result, TEST_BLOCK_NUMBER)
}

func testBiCall(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	kclient, _ := NewgKlaytnClient(addr)
	defer kclient.Close()

	knclient, err := kclient.makeKlaytnClient(timeout)
	assert.NoError(t, err)

	stream, _ := knclient.BiCall(kclient.ctx)
	go kclient.handleBiCall(stream, func() (request *RPCRequest, e error) {
		request, err := kclient.makeRPCRequest("klay", "klay_blockNumber", nil)
		if assert.NoError(t, err) {
			return request, err
		}

		return request, nil
	}, func(response *RPCResponse) error {
		var out jsonSuccessResponse
		err = json.Unmarshal(response.Payload, &out)
		assert.NoError(t, err)
		assert.Equal(t, out.Result, TEST_BLOCK_NUMBER)

		return nil
	})

	time.Sleep(3 * time.Second)
}
