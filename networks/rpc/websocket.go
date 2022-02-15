// Modifications Copyright 2022 The klaytn Authors
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
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/valyala/fasthttp"
	"gopkg.in/fatih/set.v0"

	//"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set"
	//"github.com/ethereum/go-ethereum/log"
	fastws "github.com/clevergo/websocket"
	"github.com/gorilla/websocket"
)

const (
	wsReadBuffer  = 1024
	wsWriteBuffer = 1024
)

var wsBufferPool = new(sync.Pool)

// NewWSServer creates a new websocket RPC server around an API provider.
//
// Deprecated: use Server.WebsocketHandler
func NewWSServer(allowedOrigins []string, srv *Server) *http.Server {
	return &http.Server{Handler: srv.WebsocketHandler(allowedOrigins)}
}

// WebsocketHandler returns a handler that serves JSON-RPC to WebSocket connections.
//
// allowedOrigins should be a comma-separated list of allowed origin URLs.
// To allow connections with any origin, pass "*".
func (s *Server) WebsocketHandler(allowedOrigins []string) http.Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  wsReadBuffer,
		WriteBufferSize: wsWriteBuffer,
		WriteBufferPool: wsBufferPool,
		CheckOrigin:     wsHandshakeValidator(allowedOrigins),
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			//log.Debug("WebSocket upgrade failed", "err", err)
			return
		}
		codec := newWebsocketCodec(conn)
		s.ServeCodec(codec, OptionMethodInvocation|OptionSubscriptions)
	})
}

// wsHandshakeValidator returns a handler that verifies the origin during the
// websocket upgrade process. When a '*' is specified as an allowed origins all
// connections are accepted.
func wsHandshakeValidator(allowedOrigins []string) func(*http.Request) bool {
	origins := mapset.NewSet()
	allowAllOrigins := false

	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAllOrigins = true
		}
		if origin != "" {
			origins.Add(strings.ToLower(origin))
		}
	}
	// allow localhost if no allowedOrigins are specified.
	if len(origins.ToSlice()) == 0 {
		origins.Add("http://localhost")
		if hostname, err := os.Hostname(); err == nil {
			origins.Add("http://" + strings.ToLower(hostname))
		}
	}
	//log.Debug(fmt.Sprintf("Allowed origin(s) for WS RPC interface %v", origins.ToSlice()))

	f := func(req *http.Request) bool {
		// Skip origin verification if no Origin header is present. The origin check
		// is supposed to protect against browser based attacks. Browsers always set
		// Origin. Non-browser software can put anything in origin and checking it doesn't
		// provide additional security.
		if _, ok := req.Header["Origin"]; !ok {
			return true
		}
		// Verify origin against whitelist.
		origin := strings.ToLower(req.Header.Get("Origin"))
		if allowAllOrigins || origins.Contains(origin) {
			return true
		}
		//log.Warn("Rejected WebSocket connection", "origin", origin)
		return false
	}

	return f
}

type wsHandshakeError struct {
	err    error
	status string
}

func (e wsHandshakeError) Error() string {
	s := e.err.Error()
	if e.status != "" {
		s += " (HTTP status " + e.status + ")"
	}
	return s
}

// DialWebsocket creates a new RPC client that communicates with a JSON-RPC server
// that is listening on the given endpoint.
//
// The context is used for the initial connection establishment. It does not
// affect subsequent interactions with the client.
func DialWebsocket(ctx context.Context, endpoint, origin string) (*Client, error) {
	endpoint, header, err := wsClientHeaders(endpoint, origin)
	if err != nil {
		return nil, err
	}
	dialer := websocket.Dialer{
		ReadBufferSize:  wsReadBuffer,
		WriteBufferSize: wsWriteBuffer,
		WriteBufferPool: wsBufferPool,
	}
	//return newClient(ctx, func(ctx context.Context) (ServerCodec, error) {
	return NewClient(ctx, func(ctx context.Context) (ServerCodec, error) {
		conn, resp, err := dialer.DialContext(ctx, endpoint, header)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			hErr := wsHandshakeError{err: err}
			if resp != nil {
				hErr.status = resp.Status

			}
			return nil, hErr
		}
		return newWebsocketCodec(conn), nil
	})
}

func wsClientHeaders(endpoint, origin string) (string, http.Header, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return endpoint, nil, err
	}

	header := make(http.Header)
	if origin != "" {
		header.Add("origin", origin)
	}
	if endpointURL.User != nil {
		b64auth := base64.StdEncoding.EncodeToString([]byte(endpointURL.User.String()))
		header.Add("authorization", "Basic "+b64auth)
		endpointURL.User = nil
	}
	return endpointURL.String(), header, nil
}

//func newWebsocketCodec(conn *websocket.Conn) ServerCodec {
func newWebsocketCodec(conn *websocket.Conn) ServerCodec {
	conn.SetReadLimit(maxRequestContentLength)
	return NewCodec(conn, conn.WriteJSON, conn.ReadJSON)
}

var upgrader = fastws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewFastWSServer(allowedOrigins []string, srv *Server) *fasthttp.Server {
	upgrader.CheckOrigin = wsFastHandshakeValidator(allowedOrigins)

	// TODO-Klaytn concurreny default (256 * 1024), goroutine limit (8192)
	return &fasthttp.Server{
		Concurrency:        ConcurrencyLimit,
		MaxRequestBodySize: common.MaxRequestContentLength,
		Handler:            srv.FastWebsocketHandler,
	}
}

func wsFastHandshakeValidator(allowedOrigins []string) func(ctx *fasthttp.RequestCtx) bool {
	origins := set.New()
	allowAllOrigins := false

	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAllOrigins = true
		}
		if origin != "" {
			origins.Add(strings.ToLower(origin))
		}
	}

	// allow localhost if no allowedOrigins are specified.
	if len(origins.List()) == 0 {
		origins.Add("http://localhost")
		if hostname, err := os.Hostname(); err == nil {
			origins.Add("http://" + strings.ToLower(hostname))
		}
	}

	logger.Debug(fmt.Sprintf("Allowed origin(s) for WS RPC interface %v\n", origins.List()))

	f := func(ctx *fasthttp.RequestCtx) bool {
		origin := strings.ToLower(string(ctx.Request.Header.Peek("Origin")))
		if allowAllOrigins || origins.Has(origin) {
			return true
		}
		logger.Warn(fmt.Sprintf("origin '%s' not allowed on WS-RPC interface\n", origin))
		return false
	}

	return f
}

func (srv *Server) FastWebsocketHandler(ctx *fasthttp.RequestCtx) {

	// TODO-Klaytn handle websocket protocol
	protocol := ctx.Request.Header.Peek("Sec-WebSocket-Protocol")
	if protocol != nil {
		ctx.Response.Header.Set("Sec-WebSocket-Protocol", string(protocol))
	}

	err := upgrader.Upgrade(ctx, func(conn *fastws.Conn) {
		if atomic.LoadInt32(&srv.wsConnCount) >= MaxWebsocketConnections {
			return
		}
		atomic.AddInt32(&srv.wsConnCount, 1)
		wsConnCounter.Inc(1)
		defer func() {
			atomic.AddInt32(&srv.wsConnCount, -1)
			wsConnCounter.Dec(1)
		}()
		if WebsocketReadDeadline != 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(WebsocketReadDeadline) * time.Second))
		}
		if WebsocketWriteDeadline != 0 {
			conn.SetWriteDeadline(time.Now().Add(time.Duration(WebsocketWriteDeadline) * time.Second))
		}
		//Create a custom encode/decode pair to enforce payload size and number encoding
		encoder := func(v interface{}) error {
			msg, err := json.Marshal(v)
			if err != nil {
				return err
			}
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return err
			}
			//return fastws.WriteJSON(conn, v)
			return err
		}
		decoder := func(v interface{}) error {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			dec := json.NewDecoder(bytes.NewReader(data))
			dec.UseNumber()
			return dec.Decode(v)
			//return fastws.ReadJSON(conn, v)
		}

		reader := bufio.NewReaderSize(bytes.NewReader(ctx.Request.Body()), common.MaxRequestContentLength)
		srv.ServeCodec(NewCodec(&httpReadWriteNopCloser{reader, ctx.Response.BodyWriter()}, encoder, decoder), OptionMethodInvocation|OptionSubscriptions)
	})
	if err != nil {
		logger.Error("FastWebsocketHandler fail to upgrade message", "err", err)
		return
	}
}

// NewWSServer creates a new websocket RPC server around an API provider.
//
// Deprecated: use Server.WebsocketHandler
//func NewWSServer(allowedOrigins []string, srv *Server) *http.Server {
//	return &http.Server{
//		Handler: srv.WebsocketHandler(allowedOrigins),
//	}
//}
