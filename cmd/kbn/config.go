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

package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/nat"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"gopkg.in/urfave/cli.v1"
)

const (
	DefaultHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 8551        // Default TCP port for the HTTP RPC server
	DefaultWSHost   = "localhost" // Default host interface for the websocket RPC server
	DefaultWSPort   = 8552        // Default TCP port for the websocket RPC server
	DefaultGRPCHost = "localhost" // Default host interface for the gRPC server
	DefaultGRPCPort = 8553        // Default TCP port for the gRPC server
)

type bootnodeConfig struct {
	// Parameter variables
	networkID    uint64
	addr         string
	genKeyPath   string
	nodeKeyFile  string
	nodeKeyHex   string
	natFlag      string
	netrestrict  string
	writeAddress bool

	// Context
	restrictList *netutil.Netlist
	nodeKey      *ecdsa.PrivateKey
	natm         nat.Interface
	listenAddr   string

	// Authorized Nodes are used as pre-configured nodes list which are only
	// bonded with this bootnode.
	AuthorizedNodes []*discover.Node

	// DataDir is the file system folder the node should use for any data storage
	// requirements. The configured data directory will not be directly shared with
	// registered services, instead those can use utility methods to create/access
	// databases or flat files. This enables ephemeral nodes which can fully reside
	// in memory.
	DataDir string

	// IPCPath is the requested location to place the IPC endpoint. If the path is
	// a simple file name, it is placed inside the data directory (or on the root
	// pipe path on Windows), whereas if it's a resolvable path name (absolute or
	// relative), then that specific path is enforced. An empty path disables IPC.
	IPCPath string `toml:",omitempty"`

	// HTTP module type is http server module type (fasthttp and http)
	HTTPServerType string `toml:",omitempty"`

	// HTTPHost is the host interface on which to start the HTTP RPC server. If this
	// field is empty, no HTTP API endpoint will be started.
	HTTPHost string `toml:",omitempty"`

	// HTTPPort is the TCP port number on which to start the HTTP RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful
	// for ephemeral nodes).
	HTTPPort int `toml:",omitempty"`

	// HTTPCors is the Cross-Origin Resource Sharing header to send to requesting
	// clients. Please be aware that CORS is a browser enforced security, it's fully
	// useless for custom HTTP clients.
	HTTPCors []string `toml:",omitempty"`

	// HTTPVirtualHosts is the list of virtual hostnames which are allowed on incoming requests.
	// This is by default {'localhost'}. Using this prevents attacks like
	// DNS rebinding, which bypasses SOP by simply masquerading as being within the same
	// origin. These attacks do not utilize CORS, since they are not cross-domain.
	// By explicitly checking the Host-header, the server will not allow requests
	// made against the server with a malicious host domain.
	// Requests using ip address directly are not affected
	HTTPVirtualHosts []string `toml:",omitempty"`

	// HTTPModules is a list of API modules to expose via the HTTP RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	HTTPModules []string `toml:",omitempty"`

	// HTTPTimeouts allows for customization of the timeout values used by the HTTP RPC
	// interface.
	HTTPTimeouts rpc.HTTPTimeouts

	// WSHost is the host interface on which to start the websocket RPC server. If
	// this field is empty, no websocket API endpoint will be started.
	WSHost string `toml:",omitempty"`

	// WSPort is the TCP port number on which to start the websocket RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful for
	// ephemeral nodes).
	WSPort int `toml:",omitempty"`

	// WSOrigins is the list of domain to accept websocket requests from. Please be
	// aware that the server can only act upon the HTTP request the client sends and
	// cannot verify the validity of the request header.
	WSOrigins []string `toml:",omitempty"`

	// WSModules is a list of API modules to expose via the websocket RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	WSModules []string `toml:",omitempty"`

	// WSExposeAll exposes all API modules via the WebSocket RPC interface rather
	// than just the public ones.
	//
	// *WARNING* Only set this if the node is running in a trusted network, exposing
	// private APIs to untrusted users is a major security risk.
	WSExposeAll bool `toml:",omitempty"`

	// GRPCHost is the host interface on which to start the gRPC server. If
	// this field is empty, no gRPC API endpoint will be started.
	GRPCHost string `toml:",omitempty"`

	// GRPCPort is the TCP port number on which to start the gRPC server. The
	// default zero value is valid and will pick a port number randomly (useful for
	// ephemeral nodes).
	GRPCPort int `toml:",omitempty"`

	// Logger is a custom logger to use with the p2p.Server.
	Logger log.Logger `toml:",omitempty"`
}

// splitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func splitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
}

// checkExclusive verifies that only a single instance of the provided flags was
// set by the user. Each flag might optionally be followed by a string type to
// specialize it further.
func checkExclusive(ctx *cli.Context, args ...interface{}) {
	set := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		// Make sure the next argument is a flag and skip if not set
		flag, ok := args[i].(cli.Flag)
		if !ok {
			panic(fmt.Sprintf("invalid argument, not cli.Flag type: %T", args[i]))
		}
		// Check if next arg extends current and expand its name if so
		name := flag.GetName()

		if i+1 < len(args) {
			switch option := args[i+1].(type) {
			case string:
				// Extended flag, expand the name and shift the arguments
				if ctx.GlobalString(flag.GetName()) == option {
					name += "=" + option
				}
				i++

			case cli.Flag:
			default:
				panic(fmt.Sprintf("invalid argument, not cli.Flag or string extension: %T", args[i+1]))
			}
		}
		// Mark the flag if it's set
		if ctx.GlobalIsSet(flag.GetName()) {
			set = append(set, "--"+name)
		}
	}
	if len(set) > 1 {
		log.Fatalf("Flags %v can't be used at the same time", strings.Join(set, ", "))
	}
}

func setAuthorizedNodes(ctx *cli.Context, cfg *bootnodeConfig) {
	if !ctx.GlobalIsSet(utils.AuthorizedNodesFlag.Name) {
		return
	}
	urls := ctx.GlobalString(utils.AuthorizedNodesFlag.Name)
	splitedUrls := strings.Split(urls, ",")
	cfg.AuthorizedNodes = make([]*discover.Node, 0, len(splitedUrls))
	for _, url := range splitedUrls {
		node, err := discover.ParseNode(url)
		if err != nil {
			logger.Error("URL is invalid", "kni", url, "err", err)
			continue
		}
		cfg.AuthorizedNodes = append(cfg.AuthorizedNodes, node)
	}
}

// setHTTP creates the HTTP RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setHTTP(ctx *cli.Context, cfg *bootnodeConfig) {
	if ctx.GlobalBool(utils.RPCEnabledFlag.Name) && cfg.HTTPHost == "" {
		cfg.HTTPHost = "127.0.0.1"
		if ctx.GlobalIsSet(utils.RPCListenAddrFlag.Name) {
			cfg.HTTPHost = ctx.GlobalString(utils.RPCListenAddrFlag.Name)
		}
	}

	if ctx.GlobalIsSet(utils.RPCPortFlag.Name) {
		cfg.HTTPPort = ctx.GlobalInt(utils.RPCPortFlag.Name)
	}
	if ctx.GlobalIsSet(utils.RPCCORSDomainFlag.Name) {
		cfg.HTTPCors = splitAndTrim(ctx.GlobalString(utils.RPCCORSDomainFlag.Name))
	}
	if ctx.GlobalIsSet(utils.RPCApiFlag.Name) {
		cfg.HTTPModules = splitAndTrim(ctx.GlobalString(utils.RPCApiFlag.Name))
	}
	if ctx.GlobalIsSet(utils.RPCVirtualHostsFlag.Name) {
		cfg.HTTPVirtualHosts = splitAndTrim(ctx.GlobalString(utils.RPCVirtualHostsFlag.Name))
	}
}

// setWS creates the WebSocket RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setWS(ctx *cli.Context, cfg *bootnodeConfig) {
	if ctx.GlobalBool(utils.WSEnabledFlag.Name) && cfg.WSHost == "" {
		cfg.WSHost = "127.0.0.1"
		if ctx.GlobalIsSet(utils.WSListenAddrFlag.Name) {
			cfg.WSHost = ctx.GlobalString(utils.WSListenAddrFlag.Name)
		}
	}

	if ctx.GlobalIsSet(utils.WSPortFlag.Name) {
		cfg.WSPort = ctx.GlobalInt(utils.WSPortFlag.Name)
	}
	if ctx.GlobalIsSet(utils.WSAllowedOriginsFlag.Name) {
		cfg.WSOrigins = splitAndTrim(ctx.GlobalString(utils.WSAllowedOriginsFlag.Name))
	}
	if ctx.GlobalIsSet(utils.WSApiFlag.Name) {
		cfg.WSModules = splitAndTrim(ctx.GlobalString(utils.WSApiFlag.Name))
	}
}

// setIPC creates an IPC path configuration from the set command line flags,
// returning an empty string if IPC was explicitly disabled, or the set path.
func setIPC(ctx *cli.Context, cfg *bootnodeConfig) {
	checkExclusive(ctx, utils.IPCDisabledFlag, utils.IPCPathFlag)
	switch {
	case ctx.GlobalBool(utils.IPCDisabledFlag.Name):
		cfg.IPCPath = ""
	case ctx.GlobalIsSet(utils.IPCPathFlag.Name):
		cfg.IPCPath = ctx.GlobalString(utils.IPCPathFlag.Name)
	}
}

// setgRPC creates the gRPC listener interface string from the set
// command line flags, returning empty if the gRPC endpoint is disabled.
func setgRPC(ctx *cli.Context, cfg *bootnodeConfig) {
	if ctx.GlobalBool(utils.GRPCEnabledFlag.Name) && cfg.GRPCHost == "" {
		cfg.GRPCHost = "127.0.0.1"
		if ctx.GlobalIsSet(utils.GRPCListenAddrFlag.Name) {
			cfg.GRPCHost = ctx.GlobalString(utils.GRPCListenAddrFlag.Name)
		}
	}

	if ctx.GlobalIsSet(utils.GRPCPortFlag.Name) {
		cfg.GRPCPort = ctx.GlobalInt(utils.GRPCPortFlag.Name)
	}
}

func (ctx *bootnodeConfig) checkCMDState() int {
	if ctx.genKeyPath != "" {
		return generateNodeKeySpecified
	}
	if ctx.nodeKeyFile == "" && ctx.nodeKeyHex == "" {
		return noPrivateKeyPathSpecified
	}
	if ctx.nodeKeyFile != "" && ctx.nodeKeyHex != "" {
		return nodeKeyDuplicated
	}
	if ctx.writeAddress {
		return writeOutAddress
	}
	return goodToGo
}

func (ctx *bootnodeConfig) generateNodeKey() {
	nodeKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("could not generate key: %v", err)
	}
	if err = crypto.SaveECDSA(ctx.genKeyPath, nodeKey); err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(0)
}

func (ctx *bootnodeConfig) doWriteOutAddress() {
	err := ctx.readNodeKey()
	if err != nil {
		log.Fatalf("Failed to read node key: %v", err)
	}
	fmt.Printf("%v\n", discover.PubkeyID(&(ctx.nodeKey).PublicKey))
	os.Exit(0)
}

func (ctx *bootnodeConfig) readNodeKey() error {
	var err error
	if ctx.nodeKeyFile != "" {
		ctx.nodeKey, err = crypto.LoadECDSA(ctx.nodeKeyFile)
		return err
	}
	if ctx.nodeKeyHex != "" {
		ctx.nodeKey, err = crypto.LoadECDSA(ctx.nodeKeyHex)
		return err
	}
	return nil
}

func (ctx *bootnodeConfig) validateNetworkParameter() error {
	var err error
	if ctx.natFlag != "" {
		ctx.natm, err = nat.Parse(ctx.natFlag)
		if err != nil {
			return err
		}
	}

	if ctx.netrestrict != "" {
		ctx.restrictList, err = netutil.ParseNetlist(ctx.netrestrict)
		if err != nil {
			return err
		}
	}

	if ctx.addr[0] != ':' {
		ctx.listenAddr = ":" + ctx.addr
	} else {
		ctx.listenAddr = ctx.addr
	}

	return nil
}

func (c *bootnodeConfig) HttpServerType() string {
	if c.HTTPServerType == "" {
		return "http"
	}
	return c.HTTPServerType
}

func (c *bootnodeConfig) IsFastHTTP() bool {
	return c.HttpServerType() == "fasthttp"
}

// IPCEndpoint resolves an IPC endpoint based on a configured value, taking into
// account the set data folders as well as the designated platform we're currently
// running on.
func (c *bootnodeConfig) IPCEndpoint() string {
	// Short circuit if IPC has not been enabled
	if c.IPCPath == "" {
		return ""
	}
	// On windows we can only use plain top-level pipes
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(c.IPCPath, `\\.\pipe\`) {
			return c.IPCPath
		}
		return `\\.\pipe\` + c.IPCPath
	}
	// Resolve names into the data directory full paths otherwise
	if filepath.Base(c.IPCPath) == c.IPCPath {
		if c.DataDir == "" {
			return filepath.Join(os.TempDir(), c.IPCPath)
		}
		return filepath.Join(c.DataDir, c.IPCPath)
	}
	return c.IPCPath
}

func DefaultIPCEndpoint(clientIdentifier string) string {
	if clientIdentifier == "" {
		clientIdentifier = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
		if clientIdentifier == "" {
			panic("empty executable name")
		}
	}
	config := &bootnodeConfig{DataDir: os.TempDir(), IPCPath: clientIdentifier + ".ipc"}
	return config.IPCEndpoint()
}

// HTTPEndpoint resolves an HTTP endpoint based on the configured host interface
// and port parameters.
func (c *bootnodeConfig) HTTPEndpoint() string {
	if c.HTTPHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.HTTPHost, c.HTTPPort)
}

// DefaultHTTPEndpoint returns the HTTP endpoint used by default.
func DefaultHTTPEndpoint() string {
	config := &bootnodeConfig{HTTPHost: DefaultHTTPHost, HTTPPort: DefaultHTTPPort}
	return config.HTTPEndpoint()
}

// WSEndpoint resolves a websocket endpoint based on the configured host interface
// and port parameters.
func (c *bootnodeConfig) WSEndpoint() string {
	if c.WSHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.WSHost, c.WSPort)
}

// DefaultWSEndpoint returns the websocket endpoint used by default.
func DefaultWSEndpoint() string {
	config := &bootnodeConfig{WSHost: DefaultWSHost, WSPort: DefaultWSPort}
	return config.WSEndpoint()
}

// GRPCEndpoint resolves a gRPC endpoint based on the configured host interface
// and port parameters.
func (c *bootnodeConfig) GRPCEndpoint() string {
	if c.GRPCHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.GRPCHost, c.GRPCPort)
}

// DefaultGRPCEndpoint returns the gRPC endpoint used by default.
func DefaultGRPCEndpoint() string {
	config := &bootnodeConfig{GRPCHost: DefaultGRPCHost, GRPCPort: DefaultGRPCPort}
	return config.GRPCEndpoint()
}
