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

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"
)

type Service struct{}

type Args struct {
	S string
}

func (s *Service) NoArgsRets() {
}

type Result struct {
	String string
	Int    int
	Args   *Args
}

func (s *Service) Echo(str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func (s *Service) EchoWithCtx(ctx context.Context, str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func (s *Service) Sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-time.After(duration):
	case <-ctx.Done():
	}
}

func (s *Service) Rets() (string, error) {
	return "", nil
}

func (s *Service) InvalidRets1() (error, string) {
	return nil, ""
}

func (s *Service) InvalidRets2() (string, string) {
	return "", ""
}

func (s *Service) InvalidRets3() (string, string, error) {
	return "", "", nil
}

func (s *Service) Subscription(ctx context.Context) (*Subscription, error) {
	return nil, nil
}

func TestServerRegisterName(t *testing.T) {
	server := NewServer()
	service := new(Service)

	if err := server.RegisterName("calc", service); err != nil {
		t.Fatalf("%v", err)
	}

	if len(server.services.services) != 2 {
		t.Fatalf("Expected 2 service entries, got %d", len(server.services.services))
	}

	svc, ok := server.services.services["calc"]
	if !ok {
		t.Fatalf("Expected service calc to be registered")
	}

	if len(svc.callbacks) != 5 {
		t.Errorf("Expected 5 callbacks for service 'calc', got %d", len(svc.callbacks))
	}

	if len(svc.subscriptions) != 1 {
		t.Errorf("Expected 1 subscription for service 'calc', got %d", len(svc.subscriptions))
	}
}

func testServerMethodExecution(t *testing.T, method string) {
	server := NewServer()
	service := new(Service)

	if err := server.RegisterName("test", service); err != nil {
		t.Fatalf("%v", err)
	}

	stringArg := "string arg"
	intArg := 1122
	argsArg := &Args{"abcde"}
	params := []interface{}{stringArg, intArg, argsArg}

	request := map[string]interface{}{
		"id":      12345,
		"method":  "test_" + method,
		"version": "2.0",
		"params":  params,
	}

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go server.ServeCodec(NewCodec(serverConn), 0)

	out := json.NewEncoder(clientConn)
	in := json.NewDecoder(clientConn)

	if err := out.Encode(request); err != nil {
		t.Fatal(err)
	}

	var msg jsonrpcMessage
	if err := in.Decode(&msg); err != nil {
		t.Fatal(err)
	}

	if !msg.isResponse() {
		t.Fatal("message is not response")
	}
}

func TestServerMethodExecution(t *testing.T) {
	testServerMethodExecution(t, "echo")
}

func TestServerMethodWithCtx(t *testing.T) {
	testServerMethodExecution(t, "echoWithCtx")
}

// This test checks that responses are delivered for very short-lived connections that
// only carry a single request.
func TestServerShortLivedConn(t *testing.T) {
	server := newTestServer("service", new(Service))
	defer server.Stop()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	defer listener.Close()
	go server.ServeListener(listener)

	var (
		request  = `{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n"
		wantResp = `{"jsonrpc":"2.0","id":1,"result":{"rpc":"1.0","service":"1.0"}}` + "\n"
		deadline = time.Now().Add(10 * time.Second)
	)
	for i := 0; i < 20; i++ {
		conn, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			t.Fatal("can't dial:", err)
		}
		defer conn.Close()
		conn.SetDeadline(deadline)
		// Write the request, then half-close the connection so the server stops reading.
		conn.Write([]byte(request))
		conn.(*net.TCPConn).CloseWrite()
		// Now try to get the response.
		buf := make([]byte, 2000)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal("read error:", err)
		}
		if !bytes.Equal(buf[:n], []byte(wantResp)) {
			t.Fatalf("wrong response: %s", buf[:n])
		}
	}
}
