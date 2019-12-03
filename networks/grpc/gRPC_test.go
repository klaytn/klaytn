package grpc

import (
	"encoding/json"
	"github.com/klaytn/klaytn/networks/rpc"
	"sync"
	"testing"
	"time"
)

const (
	TEST_BLOCK_NUMBER = float64(123456789)
)

type APIgRPC struct{}

func (a APIgRPC) BlockNumber() float64 {
	return TEST_BLOCK_NUMBER
}

func TestGRPC(t *testing.T) {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(2)

	addr := "127.0.0.1:4000"
	handler := rpc.NewServer()

	handler.RegisterName("klay", &APIgRPC{})

	listener := &Listener{Addr: addr}
	listener.SetRPCServer(handler)
	go listener.Start()

	time.Sleep(2 * time.Second)

	go testCall(t, addr, waitGroup)
	go testBiCall(t, addr, waitGroup)
	waitGroup.Wait()
}
func testCall(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	kclient, _ := NewgKlaytnClient(addr)
	defer kclient.Close()

	knclient, err := kclient.makeKlaytnClient(timeout)
	if err != nil {
		t.Errorf("fail to make klaytnNodeClient: err=%v\n", err)
		return
	}

	request, err := kclient.makeRPCRequest("klay", "klay_blockNumber", nil)
	if err != nil {
		t.Errorf("fail to make RPCRequest: err=%v\n", err)
		return
	}

	response, err := knclient.Call(kclient.ctx, request)
	if err != nil {
		t.Errorf("fail to call: err=%v\n", err)
		return
	}
	if err := kclient.handleRPCResponse(response); err != nil {
		t.Errorf("fail to handle RPCResponse: err=%v\n", err)
		return
	}
	var out jsonSuccessResponse
	if err := json.Unmarshal(response.Payload, &out); err != nil {
		t.Errorf("fail to handle RPCResponse: err=%v\n", err)
	}

	if out.Result != TEST_BLOCK_NUMBER {
		t.Errorf("result expected:%f, actual:%f", TEST_BLOCK_NUMBER, out.Result)
	}
}

func testBiCall(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	kclient, _ := NewgKlaytnClient(addr)
	defer kclient.Close()

	knclient, err := kclient.makeKlaytnClient(timeout)
	if err != nil {
		t.Errorf("fail to make klaytnNodeClient: err=%v\n", err)
		return
	}

	stream, _ := knclient.BiCall(kclient.ctx)
	go kclient.handleBiCall(stream, func() (request *RPCRequest, e error) {
		request, err := kclient.makeRPCRequest("klay", "klay_blockNumber", nil)
		if err != nil {
			t.Errorf("fail to make RPCRequest: err=%v\n", err)
			return request, err
		}
		return request, nil
	}, func(response *RPCResponse) error {
		var out jsonSuccessResponse
		if err := json.Unmarshal(response.Payload, &out); err != nil {
			t.Errorf("fail to handle RPCResponse: err=%v\n", err)
		}

		if out.Result != TEST_BLOCK_NUMBER {
			t.Errorf("result expected:%f, actual:%f", TEST_BLOCK_NUMBER, out.Result)
		}
		return nil
	})

	select {
	case <-time.After(3 * time.Second):
	}
}
