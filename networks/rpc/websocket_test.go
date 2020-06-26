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
	"github.com/klaytn/klaytn/common"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

	var (
		srv     = newTestServer("service", new(Service))
		httpsrv = httptest.NewServer(srv.WebsocketHandler([]string{"*"}))
		wsAddr  = "ws:" + strings.TrimPrefix(httpsrv.URL, "http:")
	)
	defer srv.Stop()
	defer httpsrv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := DialWebsocket(ctx, wsAddr, "")
	if err != nil {
		t.Fatalf("can't dial: %v", err)
	}
	defer client.Close()

	// This call sends slightly less than the limit and should work.
	var result echoResult
	arg := strings.Repeat("x", common.MaxRequestContentLength-200)
	if err := client.Call(&result, "service_echo", arg, 1); err != nil {
		t.Fatalf("valid call didn't work: %v", err)
	}
	if result.String != arg {
		t.Fatal("wrong string echoed")
	}

	// This call sends twice the allowed size and shouldn't work.
	arg = strings.Repeat("x", common.MaxRequestContentLength*2)
	err = client.Call(&result, "test_echo", arg)
	if err == nil {
		t.Fatal("no error for too large call")
	}
}
