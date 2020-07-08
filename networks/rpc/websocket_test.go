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
	"github.com/stretchr/testify/assert"
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

	// create server
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
