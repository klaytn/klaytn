// Modifications Copyright 2020 The klaytn Authors
// Copyright 2018 The go-ethereum Authors
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
// This file is derived from rpc/websocket_test.go (2020/04/03).
// Modified and improved for the klaytn development.

package rpc

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

type echoArgs struct {
	S string
}

type echoResult struct {
	String string
	Int    int
	Args   *echoArgs
}

func TestWebsocketLargeCall(t *testing.T) {
	t.Parallel()

	// create server
	var (
		srv     = newTestServer("service", new(Service))
		httpsrv = httptest.NewServer(srv.WebsocketHandler([]string{"*"}))
		wsAddr  = "ws:" + strings.TrimPrefix(httpsrv.URL, "http:")
	)
	defer srv.Stop()
	defer httpsrv.Close()
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := DialWebsocket(ctx, wsAddr, "")
	if err != nil {
		t.Fatalf("can't dial: %v", err)
	}
	defer client.Close()

	// set configurations before testing
	var result echoResult
	method := "service_echo"

	// This call sends slightly less than the limit and should work.
	arg := strings.Repeat("x", common.MaxRequestContentLength-200)
	assert.NoError(t, client.Call(&result, method, arg, 1), "valid call didn't work")
	assert.Equal(t, arg, result.String, "wrong string echoed")

	// This call sends slightly larger than the allowed size and shouldn't work.
	arg = strings.Repeat("x", common.MaxRequestContentLength)
	assert.Error(t, client.Call(&result, method, arg, 1), "no error for too large call")
}

func newTestListener() net.Listener {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	return ln
}

/*
func TestWSServer_MaxConnections(t *testing.T) {
	// create server
	var (
		srv = newTestServer("service", new(Service))
		ln  = newTestListener()
	)
	defer srv.Stop()
	defer ln.Close()

	go NewWSServer([]string{"*"}, srv).Serve(ln)
	time.Sleep(100 * time.Millisecond)

	// set max websocket connections
	MaxWebsocketConnections = 3
	testWebsocketMaxConnections(t, "ws://"+ln.Addr().String(), int(MaxWebsocketConnections))
}
*/

func TestFastWSServer_MaxConnections(t *testing.T) {
	// create server
	var (
		srv = newTestServer("service", new(Service))
		ln  = newTestListener()
	)
	defer srv.Stop()
	defer ln.Close()

	go NewFastWSServer([]string{"*"}, srv).Serve(ln)
	time.Sleep(100 * time.Millisecond)

	// set max websocket connections
	MaxWebsocketConnections = 3
	testWebsocketMaxConnections(t, "ws://"+ln.Addr().String(), int(MaxWebsocketConnections))
}

func testWebsocketMaxConnections(t *testing.T, addr string, maxConnections int) {
	var closers []*Client

	for i := 0; i <= maxConnections; i++ {
		client, err := DialWebsocket(context.Background(), addr, "")
		if err != nil {
			t.Fatal(err)
		}
		closers = append(closers, client)

		var result echoResult
		method := "service_echo"
		arg := strings.Repeat("x", i)
		err = client.Call(&result, method, arg, 1)
		if i < int(MaxWebsocketConnections) {
			assert.NoError(t, err)
			assert.Equal(t, arg, result.String, "wrong string echoed")
		} else {
			assert.Error(t, err)
			// assert.Equal(t, "EOF", err.Error())
		}
	}

	for _, client := range closers {
		client.Close()
	}
}

func TestWebsocketClientHeaders(t *testing.T) {
	t.Parallel()

	endpoint, header, err := wsClientHeaders("wss://testuser:test-PASS_01@example.com:1234", "https://example.com")
	if err != nil {
		t.Fatalf("wsGetConfig failed: %s", err)
	}
	if endpoint != "wss://example.com:1234" {
		t.Fatal("User should have been stripped from the URL")
	}
	if header.Get("authorization") != "Basic dGVzdHVzZXI6dGVzdC1QQVNTXzAx" {
		t.Fatal("Basic auth header is incorrect")
	}
	if header.Get("origin") != "https://example.com" {
		t.Fatal("Origin not set")
	}
}

// This test checks that the server rejects connections from disallowed origins.
func TestWebsocketOriginCheck(t *testing.T) {
	t.Parallel()

	var (
		srv     = newTestServer("service", new(Service))
		httpsrv = httptest.NewServer(srv.WebsocketHandler([]string{"http://example.com"}))
		wsURL   = "ws:" + strings.TrimPrefix(httpsrv.URL, "http:")
	)
	defer srv.Stop()
	defer httpsrv.Close()

	client, err := DialWebsocket(context.Background(), wsURL, "http://ekzample.com")
	if err == nil {
		client.Close()
		t.Fatal("no error for wrong origin")
	}
	wantErr := wsHandshakeError{websocket.ErrBadHandshake, "403 Forbidden"}
	if !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("wrong error for wrong origin: %q", err)
	}
}

func TestClientWebsocketPing(t *testing.T) {
	t.Parallel()

	var (
		sendPing    = make(chan struct{})
		server      = wsPingTestServer(t, sendPing)
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	)
	defer cancel()
	defer server.Shutdown(ctx)

	client, err := DialContext(ctx, "ws://"+server.Addr)
	if err != nil {
		t.Fatalf("client dial error: %v", err)
	}
	resultChan := make(chan int)
	sub, err := client.KlaySubscribe(ctx, resultChan, "foo")
	if err != nil {
		t.Fatalf("client subscribe error: %v", err)
	}

	// Wait for the context's deadline to be reached before proceeding.
	// This is important for reproducing https://github.com/ethereum/go-ethereum/issues/19798
	<-ctx.Done()
	close(sendPing)

	// Wait for the subscription result.
	timeout := time.NewTimer(5 * time.Second)
	for {
		select {
		case err := <-sub.Err():
			t.Error("client subscription error:", err)
		case result := <-resultChan:
			t.Log("client got result:", result)
			return
		case <-timeout.C:
			t.Error("didn't get any result within the test timeout")
			return
		}
	}
}

// wsPingTestServer runs a WebSocket server which accepts a single subscription request.
// When a value arrives on sendPing, the server sends a ping frame, waits for a matching
// pong and finally delivers a single subscription result.
func wsPingTestServer(t *testing.T, sendPing <-chan struct{}) *http.Server {
	var srv http.Server
	shutdown := make(chan struct{})
	srv.RegisterOnShutdown(func() {
		close(shutdown)
	})
	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade to WebSocket.
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("server WS upgrade error: %v", err)
			return
		}
		defer conn.Close()

		// Handle the connection.
		wsPingTestHandler(t, conn, shutdown, sendPing)
	})

	// Start the server.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	srv.Addr = listener.Addr().String()
	go srv.Serve(listener)
	return &srv
}

func wsPingTestHandler(t *testing.T, conn *websocket.Conn, shutdown, sendPing <-chan struct{}) {
	// Canned responses for the eth_subscribe call in TestClientWebsocketPing.
	const (
		subResp   = `{"jsonrpc":"2.0","id":1,"result":"0x00"}`
		subNotify = `{"jsonrpc":"2.0","method":"eth_subscription","params":{"subscription":"0x00","result":1}}`
	)

	// Handle subscribe request.
	if _, _, err := conn.ReadMessage(); err != nil {
		t.Errorf("server read error: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte(subResp)); err != nil {
		t.Errorf("server write error: %v", err)
		return
	}

	// Read from the connection to process control messages.
	pongCh := make(chan string)
	conn.SetPongHandler(func(d string) error {
		t.Logf("server got pong: %q", d)
		pongCh <- d
		return nil
	})
	go func() {
		for {
			typ, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			t.Logf("server got message (%d): %q", typ, msg)
		}
	}()

	// Write messages.
	var (
		sendResponse <-chan time.Time
		wantPong     string
	)
	for {
		select {
		case _, open := <-sendPing:
			if !open {
				sendPing = nil
			}
			t.Logf("server sending ping")
			conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			wantPong = "ping"
		case data := <-pongCh:
			if wantPong == "" {
				t.Errorf("unexpected pong")
			} else if data != wantPong {
				t.Errorf("got pong with wrong data %q", data)
			}
			wantPong = ""
			sendResponse = time.NewTimer(200 * time.Millisecond).C
		case <-sendResponse:
			t.Logf("server sending response")
			conn.WriteMessage(websocket.TextMessage, []byte(subNotify))
			sendResponse = nil
		case <-shutdown:
			conn.Close()
			return
		}
	}
}

func TestWebsocketAuthCheck(t *testing.T) {
	t.Parallel()

	var (
		srv     = newTestServer("service", new(Service))
		httpsrv = httptest.NewServer(srv.WebsocketHandler([]string{"http://example.com"}))
		wsURL   = "ws://testuser:test-PASS_01@" + strings.TrimPrefix(httpsrv.URL, "http://")
	)
	connect := false
	origHandler := httpsrv.Config.Handler
	httpsrv.Config.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:test-PASS_01"))
			if r.Method == http.MethodGet && auth == expectedAuth {
				connect = true
				w.WriteHeader(http.StatusSwitchingProtocols)
				return
			}
			if !connect {
				http.Error(w, "connect with authorization not received", http.StatusMethodNotAllowed)
				return
			}
			origHandler.ServeHTTP(w, r)
		})
	defer srv.Stop()
	defer httpsrv.Close()

	client, err := DialWebsocket(context.Background(), wsURL, "http://example.com")
	if err == nil {
		client.Close()
		t.Fatal("no error for connect with auth header")
	}
	wantErr := wsHandshakeError{websocket.ErrBadHandshake, "101 Switching Protocols"}
	if !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("wrong error for auth header: %q", err)
	}
}
