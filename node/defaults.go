// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from node/defaults.go (2018/06/04).
// Modified and improved for the klaytn development.

package node

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/klaytn/klaytn/networks/rpc"

	"github.com/klaytn/klaytn/storage/database"

	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/nat"
)

const (
	DefaultHTTPHost               = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort               = 8551        // Default TCP port for the HTTP RPC server
	DefaultWSHost                 = "localhost" // Default host interface for the websocket RPC server
	DefaultWSPort                 = 8552        // Default TCP port for the websocket RPC server
	DefaultGRPCHost               = "localhost" // Default host interface for the gRPC server
	DefaultGRPCPort               = 8553        // Default TCP port for the gRPC server
	DefaultP2PPort                = 32323
	DefaultP2PSubPort             = 32324
	DefaultMaxPhysicalConnections = 10 // Default the max number of node's physical connections
)

// DefaultConfig contains reasonable default settings.
var DefaultConfig = Config{
	DBType:           DefaultDBType(),
	DataDir:          DefaultDataDir(),
	HTTPPort:         DefaultHTTPPort,
	HTTPModules:      []string{"net", "web3"},
	HTTPVirtualHosts: []string{"localhost"},
	HTTPTimeouts:     rpc.DefaultHTTPTimeouts,
	WSPort:           DefaultWSPort,
	WSModules:        []string{"net", "web3"},
	GRPCPort:         DefaultGRPCPort,
	P2P: p2p.Config{
		ListenAddr:             fmt.Sprintf(":%d", DefaultP2PPort),
		MaxPhysicalConnections: DefaultMaxPhysicalConnections,
		NAT:                    nat.Any(),
	},
}

func DefaultDBType() database.DBType {
	return database.LevelDB
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	dirname := filepath.Base(os.Args[0])
	if dirname == "" {
		dirname = "klay"
	}
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", strings.ToUpper(dirname))
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", strings.ToUpper(dirname))
		} else {
			return filepath.Join(home, "."+dirname)
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
