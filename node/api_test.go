// Modifications Copyright 2021 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from node/api_test.go (2021/05/17).
// Modified and improved for the klaytn development.

package node

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/stretchr/testify/assert"
)

// This test uses the admin_startRPC and admin_startWS APIs,
// checking whether the HTTP server is started correctly.
type test struct {
	name string
	cfg  Config
	fn   func(*testing.T, *Node, *PrivateAdminAPI)

	// Checks. These run after the node is configured and all API calls have been made.
	wantReachable bool // whether the HTTP server should be reachable at all
	wantRPC       bool // whether JSON-RPC/HTTP should be accessible
	wantWS        bool // whether JSON-RPC/WS should be accessible
}

func TestStartRPC(t *testing.T) {
	tests := []test{
		{
			name: "all off",
			cfg:  Config{},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
			},
			wantReachable: false,
			wantRPC:       false,
			wantWS:        false,
		},
		{
			name: "rpc enabled through config",
			cfg:  Config{HTTPHost: "127.0.0.1"},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
			},
			wantReachable: true,
			wantRPC:       true,
			wantWS:        false,
		},
		{
			name: "rpc enabled through API",
			cfg:  Config{},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
				_, err := api.StartHTTP(sp("127.0.0.1"), ip(0), nil, nil, nil)
				assert.NoError(t, err)
			},
			wantReachable: true,
			wantRPC:       true,
			wantWS:        false,
		},
		{
			name: "rpc start again after failure",
			cfg:  Config{},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
				// Listen on a random port.
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				if err != nil {
					t.Fatal("can't listen:", err)
				}
				defer listener.Close()
				port := listener.Addr().(*net.TCPAddr).Port

				// Now try to start RPC on that port. This should fail.
				_, err = api.StartHTTP(sp("127.0.0.1"), ip(port), nil, nil, nil)
				if err == nil {
					t.Fatal("StartHTTP should have failed on port", port)
				}

				// Try again after unblocking the port. It should work this time.
				listener.Close()
				_, err = api.StartHTTP(sp("127.0.0.1"), ip(port), nil, nil, nil)
				assert.NoError(t, err)
			},
			wantReachable: true,
			wantRPC:       true,
			wantWS:        false,
		},
		{
			name: "rpc stopped through API",
			cfg:  Config{HTTPHost: "127.0.0.1"},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
				_, err := api.StopHTTP()
				assert.NoError(t, err)
			},
			wantReachable: false,
			wantRPC:       false,
			wantWS:        false,
		},
		{
			name:          "ws enabled through config",
			cfg:           Config{WSHost: "127.0.0.1"},
			wantReachable: false,
			wantRPC:       false,
			wantWS:        true,
		},
		{
			name: "ws enabled through API",
			cfg:  Config{},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
				_, err := api.StartWS(sp("127.0.0.1"), ip(0), nil, nil)
				assert.NoError(t, err)
			},
			wantReachable: false,
			wantRPC:       false,
			wantWS:        true,
		},
		{
			name: "ws stopped through API",
			cfg:  Config{WSHost: "127.0.0.1"},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
				_, err := api.StopWS()
				assert.NoError(t, err)
			},
			wantReachable: false,
			wantRPC:       false,
			wantWS:        false,
		},
		{
			name: "ws enabled after RPC",
			cfg:  Config{HTTPHost: "127.0.0.1"},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
				_, err := api.StartWS(sp("127.0.0.1"), ip(0), nil, nil)
				assert.NoError(t, err)
			},
			wantReachable: true,
			wantRPC:       true,
			wantWS:        true,
		},
		{
			name: "ws enabled after RPC then stopped",
			cfg:  Config{HTTPHost: "127.0.0.1"},
			fn: func(t *testing.T, n *Node, api *PrivateAdminAPI) {
				_, err := api.StartWS(sp("127.0.0.1"), ip(0), nil, nil)
				assert.NoError(t, err)

				_, err = api.StopWS()
				assert.NoError(t, err)
			},
			wantReachable: true,
			wantRPC:       true,
			wantWS:        false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runTestWithServerType(t, test, "fasthttp")
			runTestWithServerType(t, test, "http")
		})
	}
}

func runTestWithServerType(t *testing.T, test test, httpServerType string) {
	// Setting test node config
	config := test.cfg
	config.P2P.NoDiscovery = true

	// Create Node.
	stack, err := New(&config)
	if err != nil {
		t.Fatal("can't create node:", err)
	}

	// Register the test config.
	stack.config.HTTPServerType = httpServerType
	stack.config.HTTPPort = 0
	stack.config.WSPort = 0

	if err := stack.Start(); err != nil {
		t.Fatal("can't start node:", err)
	}

	defer stack.Stop()

	// Run the API call hook.
	if test.fn != nil {
		test.fn(t, stack, &PrivateAdminAPI{stack})
	}

	httpBaseURL := "http://" + stack.httpEndpoint
	if stack.httpListener != nil {
		httpBaseURL = "http://" + stack.httpListener.Addr().String()
	}

	wsBaseURL := "ws://" + stack.wsEndpoint
	if stack.wsListener != nil {
		wsBaseURL = "ws://" + stack.wsListener.Addr().String()
	}

	reachable := checkReachable(httpBaseURL)
	rpcAvailable := checkRPC(httpBaseURL)
	wsAvailable := checkRPC(wsBaseURL)

	if reachable != test.wantReachable {
		t.Errorf("HTTP server is %sreachable, want it %sreachable", not(reachable), not(test.wantReachable))
	}
	if rpcAvailable != test.wantRPC {
		t.Errorf("HTTP RPC %savailable, want it %savailable", not(rpcAvailable), not(test.wantRPC))
	}
	if wsAvailable != test.wantWS {
		t.Errorf("WS RPC %savailable, want it %savailable", not(wsAvailable), not(test.wantWS))
	}
}

// checkReachable checks if the TCP endpoint in rawurl is open.
func checkReachable(rawurl string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// checkBodyOK checks whether the given HTTP URL responds with 200 OK and body "OK".
func checkBodyOK(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}
	buf := make([]byte, 2)
	if _, err = io.ReadFull(resp.Body, buf); err != nil {
		return false
	}
	return bytes.Equal(buf, []byte("OK"))
}

// checkRPC checks whether JSON-RPC works against the given URL.
func checkRPC(url string) bool {
	c, err := rpc.Dial(url)
	if err != nil {
		return false
	}
	defer c.Close()

	_, err = c.SupportedModules()
	return err == nil
}

// string/int pointer helpers.
func sp(s string) *string { return &s }
func ip(i int) *int       { return &i }

func not(ok bool) string {
	if ok {
		return ""
	}
	return "not "
}
