// Modifications Copyright 2019 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from eth/config.go (2018/06/04).
// Modified and improved for the klaytn development.

package sc

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/discover"
)

const (
	datadirMainChainBridgeNodes = "main-bridges.json" // Path within the datadir to the static node list
)

var logger = log.NewModuleLogger(log.ServiceChain)

// DefaultConfig contains default settings for use on the Klaytn main net.
var DefaultConfig = SCConfig{
	NetworkId: 1,
	MaxPeer:   1, // Only a single main-bridge and sub-bridge pair is allowed.
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
}

//go:generate gencodec -type SCConfig -formats toml -out gen_config.go
type SCConfig struct {
	// Name sets the instance name of the node. It must not contain the / character and is
	// used in the devp2p node identifier. The instance name is "kscn". If no
	// value is specified, the basename of the current executable is used.
	Name string `toml:"-"`

	// BridgeService
	EnabledMainBridge bool
	EnabledSubBridge  bool
	DataDir           string

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	LevelDBCacheSize   int
	TrieCacheSize      int
	TrieTimeout        time.Duration
	TrieBlockInterval  uint
	ChildChainIndexing bool

	// Network
	MainBridgePort string
	SubBridgePort  string
	MaxPeer        int

	// ServiceChain
	ServiceChainConsensus string
	AnchoringPeriod       uint64
	SentChainTxsLimit     uint64

	ParentChainID                      uint64
	VTRecovery                         bool
	VTRecoveryInterval                 uint64
	Anchoring                          bool
	ServiceChainParentOperatorGasLimit uint64
	ServiceChainChildOperatorGasLimit  uint64

	// KAS
	KASAnchor               bool
	KASAnchorUrl            string
	KASAnchorPeriod         uint64
	KASAnchorOperator       string
	KASAccessKey            string
	KASSecretKey            string
	KASXChainId             string
	KASAnchorRequestTimeout time.Duration
}

// NodeName returns the devp2p node identifier.
func (c *SCConfig) NodeName() string {
	name := c.name()
	// Backwards compatibility: previous versions used title-cased "Klaytn", keep that.
	if name == "klay" || name == "klay-testnet" {
		name = "Klaytn"
	}
	name += "/" + runtime.GOOS + "-" + runtime.GOARCH
	name += "/" + runtime.Version()
	return name
}

func (c *SCConfig) name() string {
	if c.Name == "" {
		progname := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
		if progname == "" {
			panic("empty executable name, set Config.Name")
		}
		return progname
	}
	return c.Name
}

// StaticNodes returns a list of node enode URLs configured as static nodes.
func (c *SCConfig) MainBridges() []*discover.Node {
	return c.parsePersistentNodes(filepath.Join(c.DataDir, datadirMainChainBridgeNodes))
}

func (c *SCConfig) parsePersistentNodes(path string) []*discover.Node {
	// Short circuit if no node config is present
	if c.DataDir == "" {
		return nil
	}
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	// Load the nodes from the config file.
	var nodelist []string
	if err := common.LoadJSON(path, &nodelist); err != nil {
		logger.Error(fmt.Sprintf("Can't load node file %s: %v", path, err))
		return nil
	}
	// Interpret the list as a discovery node array
	var nodes []*discover.Node
	for _, url := range nodelist {
		if url == "" {
			continue
		}
		node, err := discover.ParseNode(url)
		if err != nil {
			logger.Error(fmt.Sprintf("Node URL %s: %v\n", url, err))
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// ResolvePath resolves path in the instance directory.
func (c *SCConfig) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if c.DataDir == "" {
		return ""
	}
	return filepath.Join(c.instanceDir(), path)
}

func (c *SCConfig) instanceDir() string {
	if c.DataDir == "" {
		return ""
	}
	return filepath.Join(c.DataDir, c.name())
}

func (c *SCConfig) getKey(path string) *ecdsa.PrivateKey {
	keyFile := c.ResolvePath(path)
	if key, err := crypto.LoadECDSA(keyFile); err == nil {
		return key
	}
	// No persistent key found, generate and store a new one.
	key, err := crypto.GenerateKey()
	if err != nil {
		logger.Crit("Failed to generate chain key", "err", err)
	}
	instanceDir := filepath.Join(c.DataDir, c.name())
	if err := os.MkdirAll(instanceDir, 0700); err != nil {
		logger.Crit("Failed to make dir to persist chain key", "err", err)
	}
	keyFile = filepath.Join(instanceDir, path)
	if err := crypto.SaveECDSA(keyFile, key); err != nil {
		logger.Crit("Failed to persist chain key", "err", err)
	}
	return key
}
