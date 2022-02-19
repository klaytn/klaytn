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
	gorillaws "github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestWebsocketLargeCallFastws(t *testing.T) {
	t.Parallel()

	// create server
	var (
		srv    = newTestServer("service", new(Service))
		ln     = newTestListener()
		wsAddr = "ws://" + ln.Addr().String()
	)
	defer srv.Stop()
	defer ln.Close()

	go NewFastWSServer([]string{"*"}, srv).Serve(ln)
	time.Sleep(100 * time.Millisecond)

	//defer srv.Stop()
	//defer httpsrv.Close()
	//fmt.Println("server", httpsrv.Listener.Addr())
	//time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := DialWebsocket(ctx, wsAddr, "")
	//fmt.Println("dial web socket ", client, err)
	if err != nil {
		t.Fatalf("can't dial: %v", err)
	}
	defer client.Close()

	// set configurations before testing
	var result echoResult
	method := "service_echo"

	// set message size
	messageSize := 200

	messageSize, err = client.getMessageSize(method)
	assert.NoError(t, err)
	requestMaxLen := common.MaxRequestContentLength - messageSize
	// This call sends slightly less than the limit and should work.
	arg := strings.Repeat("x", requestMaxLen-1)

	assert.NoError(t, client.Call(&result, method, arg, 1), "valid call didn't work")
	assert.Equal(t, arg, result.String, "wrong string echoed")

	// This call sends slightly larger than the allowed size and shouldn't work.
	arg = strings.Repeat("x", requestMaxLen)
	assert.Error(t, client.Call(&result, method, arg), "no error for too large call")
}

func TestWebsocketLargeCall(t *testing.T) {
	t.Parallel()

	// create server
	var (
		srv     = newTestServer("service", new(Service))
		httpsrv = httptest.NewServer(srv.GSWebsocketHandler([]string{"*"}))
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

	// set message size
	messageSize := 200
	messageSize, err = client.getMessageSize(method)
	assert.NoError(t, err)
	requestMaxLen := common.MaxRequestContentLength - messageSize

	// This call sends slightly less than the limit and should work.
	arg := strings.Repeat("x", requestMaxLen-1)
	assert.NoError(t, client.Call(&result, method, arg, 1), "valid call didn't work")
	assert.Equal(t, arg, result.String, "wrong string echoed")

	// This call sends slightly larger than the allowed size and shouldn't work.
	arg = strings.Repeat("x", requestMaxLen)
	assert.Error(t, client.Call(&result, method, arg), "no error for too large call")
}

func newTestListener() net.Listener {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	return ln
}

func TestWSServer_MaxConnections(t *testing.T) {
	// create server
	var (
		srv = newTestServer("service", new(Service))
		ln  = newTestListener()
	)
	defer srv.Stop()
	defer ln.Close()

	go NewGSWSServer([]string{"*"}, srv).Serve(ln)
	time.Sleep(100 * time.Millisecond)

	// set max websocket connections
	MaxWebsocketConnections = 3
	testWebsocketMaxConnections(t, "ws://"+ln.Addr().String(), int(MaxWebsocketConnections))
}

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
	var expectedGSwsError = gorillaws.CloseError{
		Code: gorillaws.CloseGoingAway,
		Text: "unexpected EOF",
	}
	var expectedFastwsError = gorillaws.CloseError{
		Code: gorillaws.CloseAbnormalClosure,
		Text: "unexpected EOF",
	}

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
			//assert.Equal(t, "EOF", err.Error())
			if t.Name() == "TestWSServer_MaxConnections" {
				assert.Equal(t, expectedGSwsError.Error(), err.Error())
			}
			if t.Name() == "TestFastWSServer_MaxConnections" {
				assert.Equal(t, expectedFastwsError.Error(), err.Error())
			}
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

func TestWebsocketAuthCheck(t *testing.T) {
	t.Parallel()

	var (
		srv     = newTestServer("websocket test", new(Service))
		httpsrv = httptest.NewServer(srv.GSWebsocketHandler([]string{"http://example.com"}))
		wsURL   = "ws://testuser:test-PASS_01@" + strings.TrimPrefix(httpsrv.URL, "http://")
	)
	connect := false
	origHandler := httpsrv.Config.Handler
	httpsrv.Config.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			//fmt.Println("Received auth header = ", auth)
			expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:test-PASS_01"))
			//fmt.Println("expected auth  = ", expectedAuth)
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

	client, err := DialWebsocket(context.Background(), wsURL, "")
	if err == nil {
		client.Close()
		t.Fatal("no error for connect with auth header")
	}
	//if err != websocket.ErrBadHandshake {
	if err.Error() != "websocket: bad handshake" {
		t.Fatalf("wrong error for header: %q", err)
	}
}
