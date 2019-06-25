// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"github.com/valyala/fasthttp"
	"net"
	"net/url"
	"strings"
	"time"
)

const (
	originHeader     = "Origin"
	protocolHeader   = "Sec-Websocket-Protocol"
	websocketVersion = "13"
)

// HandshakeError describes an error with the handshake from the peer.
type HandshakeError struct {
	message string
}

func (e HandshakeError) Error() string {
	return e.message
}

// Upgrader specifies parameters for upgrading an HTTP connection to a
// WebSocket connection.
type Upgrader struct {
	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
	// size is zero, then a default value of 4096 is used. The I/O buffer sizes
	// do not limit the size of the messages that can be sent or received.
	ReadBufferSize, WriteBufferSize int

	// Subprotocols specifies the server's supported protocols in order of
	// preference. If this field is set, then the Upgrade method negotiates a
	// subprotocol by selecting the first match in this list with a protocol
	// requested by the client.
	Subprotocols []string

	// Error specifies the function for generating HTTP error responses. If Error
	// is nil, then http.Error is used to generate the HTTP response.
	Error func(ctx *fasthttp.RequestCtx, status int, reason error)

	// CheckOrigin returns true if the request Origin header is acceptable. If
	// CheckOrigin is nil, the host in the Origin header must not be set or
	// must match the host of the request.
	CheckOrigin func(ctx *fasthttp.RequestCtx) bool

	// WebSocket connection handler
	Handler func(*Conn)
}

func (u *Upgrader) returnError(ctx *fasthttp.RequestCtx, status int, reason string) error {
	err := HandshakeError{reason}
	if u.Error != nil {
		u.Error(ctx, status, err)
	} else {
		ctx.SetStatusCode(status)
		ctx.Response.Header.Set(protocolHeader, websocketVersion)
		ctx.SetBodyString(fasthttp.StatusMessage(status))
	}
	return err
}

// checkSameOrigin returns true if the origin is not set or is equal to the request host.
func checkSameOrigin(ctx *fasthttp.RequestCtx) bool {
	origin := string(ctx.Request.Header.Peek(originHeader))
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	return u.Host == string(ctx.Host())
}

func (u *Upgrader) selectSubprotocol(ctx *fasthttp.RequestCtx) string {
	if u.Subprotocols != nil {
		clientProtocols := Subprotocols(ctx)
		for _, serverProtocol := range u.Subprotocols {
			for _, clientProtocol := range clientProtocols {
				if clientProtocol == serverProtocol {
					return clientProtocol
				}
			}
		}
	}

	return ""
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// application negotiated subprotocol (Sec-Websocket-Protocol).
//
// If the upgrade fails, then Upgrade replies to the client with an HTTP error
// response.
func (u *Upgrader) Upgrade(ctx *fasthttp.RequestCtx, handlers ...func(*Conn)) error {
	handler := u.Handler
	if len(handlers) > 0 {
		handler = handlers[0]
	}
	if handler == nil {
		panic("Upgrader's handler must be set.")
	}

	if !ctx.IsGet() {
		return u.returnError(ctx, fasthttp.StatusMethodNotAllowed, "websocket: method not GET")
	}
	if !tokenListContainsValue(string(ctx.Request.Header.Peek("Sec-WebSocket-Version")), websocketVersion) {
		return u.returnError(ctx, fasthttp.StatusBadRequest, "websocket: version != 13")
	}

	if !tokenListContainsValue(string(ctx.Request.Header.Peek("Connection")), "upgrade") {
		return u.returnError(ctx, fasthttp.StatusBadRequest, "websocket: could not find connection header with token 'upgrade'")
	}

	if !tokenListContainsValue(string(ctx.Request.Header.Peek("Upgrade")), "websocket") {
		return u.returnError(ctx, fasthttp.StatusBadRequest, "websocket: could not find upgrade header with token 'websocket'")
	}

	checkOrigin := u.CheckOrigin
	if checkOrigin == nil {
		checkOrigin = checkSameOrigin
	}
	if !checkOrigin(ctx) {
		return u.returnError(ctx, fasthttp.StatusForbidden, "websocket: origin not allowed")
	}

	challengeKey := ctx.Request.Header.Peek("Sec-WebSocket-Key")
	if len(challengeKey) == 0 {
		return u.returnError(ctx, fasthttp.StatusBadRequest, "websocket: key missing or blank")
	}

	ctx.SetStatusCode(fasthttp.StatusSwitchingProtocols)
	ctx.Response.Header.Set("Upgrade", "websocket")
	ctx.Response.Header.Set("Connection", "Upgrade")
	ctx.Response.Header.Set("Sec-WebSocket-Accept", computeAcceptKeyByte(challengeKey))

	// The subprotocol may have already been set in the response
	subprotocol := u.selectSubprotocol(ctx)

	ctx.Hijack(func(conn net.Conn) {
		c := newConn(conn, true, u.ReadBufferSize, u.WriteBufferSize)
		c.subprotocol = subprotocol
		handler(c)
	})

	return nil
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// This function is deprecated, use websocket.Upgrader instead.
//
// The application is responsible for checking the request origin before
// calling Upgrade. An example implementation of the same origin policy is:
//
//	if req.Header.Get("Origin") != "http://"+req.Host {
//		http.Error(w, "Origin not allowed", 403)
//		return
//	}
//
// If the endpoint supports subprotocols, then the application is responsible
// for negotiating the protocol used on the connection. Use the Subprotocols()
// function to get the subprotocols requested by the client. Use the
// Sec-Websocket-Protocol response header to specify the subprotocol selected
// by the application.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// negotiated subprotocol (Sec-Websocket-Protocol).
//
// The connection buffers IO to the underlying network connection. The
// readBufSize and writeBufSize parameters specify the size of the buffers to
// use. Messages can be larger than the buffers.
//
// If the request is not a valid WebSocket handshake, then Upgrade returns an
// error of type HandshakeError. Applications should handle this error by
// replying to the client with an HTTP error response.
func Upgrade(ctx *fasthttp.RequestCtx, readBufSize, writeBufSize int) error {
	u := Upgrader{ReadBufferSize: readBufSize, WriteBufferSize: writeBufSize}
	u.Error = func(ctx *fasthttp.RequestCtx, status int, reason error) {
		// don't return errors to maintain backwards compatibility
	}
	u.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool {
		// allow all connections by default
		return true
	}
	return u.Upgrade(ctx)
}

// Subprotocols returns the subprotocols requested by the client in the
// Sec-Websocket-Protocol header.
func Subprotocols(ctx *fasthttp.RequestCtx) []string {
	h := strings.TrimSpace(string(ctx.Request.Header.Peek(protocolHeader)))
	if h == "" {
		return nil
	}
	protocols := strings.Split(h, ",")
	for i := range protocols {
		protocols[i] = strings.TrimSpace(protocols[i])
	}
	return protocols
}

// IsWebSocketUpgrade returns true if the client requested upgrade to the
// WebSocket protocol.
func IsWebSocketUpgrade(ctx *fasthttp.RequestCtx) bool {
	return tokenListContainsValue(string(ctx.Request.Header.Peek("Connection")), "upgrade") &&
		tokenListContainsValue(string(ctx.Request.Header.Peek("Upgrade")), "websocket")
}
