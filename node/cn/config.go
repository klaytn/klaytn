// Modifications Copyright 2018 The klaytn Authors
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

package cn

import (
	"math/big"
	"os"
	"os/user"
	"time"

	"github.com/klaytn/klaytn/storage/statedb"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/node/cn/gasprice"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

var logger = log.NewModuleLogger(log.NodeCN)

// GetDefaultConfig returns default settings for use on the Klaytn main net.
func GetDefaultConfig() *Config {
	return &Config{
		SyncMode:            downloader.FullSync,
		NetworkId:           params.CypressNetworkId,
		LevelDBCacheSize:    768,
		TrieCacheSize:       512,
		TrieTimeout:         5 * time.Minute,
		TrieBlockInterval:   blockchain.DefaultBlockInterval,
		TrieNodeCacheConfig: *statedb.GetEmptyTrieNodeCacheConfig(),
		TriesInMemory:       blockchain.DefaultTriesInMemory,
		GasPrice:            big.NewInt(18 * params.Ston),

		TxPool: blockchain.DefaultTxPoolConfig,
		GPO: gasprice.Config{
			Blocks:           20,
			Percentile:       60,
			MaxHeaderHistory: 1024,
			MaxBlockHistory:  1024,
		},
		WsEndpoint: "localhost:8546",

		Istanbul: *istanbul.DefaultConfig,
	}
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Klaytn main net block is used.
	Genesis *blockchain.Genesis `toml:",omitempty"`

	// Protocol options
	NetworkId     uint64 // Network ID to use for selecting peers to connect to
	SyncMode      downloader.SyncMode
	NoPruning     bool
	WorkerDisable bool // disables worker and does not start istanbul

	// KES options
	DownloaderDisable bool
	FetcherDisable    bool

	// Service chain options
	ParentOperatorAddr *common.Address `toml:",omitempty"` // A hex account address in the parent chain used to sign a child chain transaction.
	AnchoringPeriod    uint64          // Period when child chain sends an anchoring transaction to the parent chain. Default value is 1.
	SentChainTxsLimit  uint64          // Number of chain transactions stored for resending. Default value is 1000.

	// Light client options
	//LightServ  int `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	//LightPeers int `toml:",omitempty"` // Maximum number of LES client peers

	OverwriteGenesis bool
	StartBlockNumber uint64

	// Database options
	DBType               database.DBType
	SkipBcVersionCheck   bool `toml:"-"`
	SingleDB             bool
	NumStateTrieShards   uint
	EnableDBPerfMetrics  bool
	LevelDBCompression   database.LevelDBCompressionType
	LevelDBBufferPool    bool
	LevelDBCacheSize     int
	DynamoDBConfig       database.DynamoDBConfig
	TrieCacheSize        int
	TrieTimeout          time.Duration
	TrieBlockInterval    uint
	TriesInMemory        uint64
	SenderTxHashIndexing bool
	ParallelDBWrite      bool
	TrieNodeCacheConfig  statedb.TrieNodeCacheConfig
	SnapshotCacheSize    int

	// Mining-related options
	ServiceChainSigner common.Address `toml:",omitempty"`
	ExtraData          []byte         `toml:",omitempty"`
	GasPrice           *big.Int

	// Reward
	Rewardbase common.Address `toml:",omitempty"`

	// Transaction pool options
	TxPool blockchain.TxPoolConfig

	// Gas Price Oracle options
	GPO gasprice.Config

	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool
	// Enables collecting internal transaction data during processing a block
	EnableInternalTxTracing bool
	// Istanbul options
	Istanbul istanbul.Config

	// Miscellaneous options
	DocRoot string `toml:"-"`

	WsEndpoint string `toml:",omitempty"`

	// Tx Resending options
	TxResendInterval  uint64
	TxResendCount     int
	TxResendUseLegacy bool

	// Service Chain
	NoAccountCreation bool

	// use separate network different from baobab or cypress
	IsPrivate bool

	// Restart
	AutoRestartFlag    bool
	RestartTimeOutFlag time.Duration
	DaemonPathFlag     string

	// RPCGasCap is the global gas cap for eth-call variants.
	RPCGasCap *big.Int `toml:",omitempty"`

	// RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for
	// send-transction variants. The unit is klay.
	// This is used by eth namespace RPC APIs
	RPCTxFeeCap float64
}

type configMarshaling struct {
	ExtraData hexutil.Bytes
}

func (c *Config) getVMConfig() vm.Config {
	return vm.Config{
		EnablePreimageRecording: c.EnablePreimageRecording,
		EnableInternalTxTracing: c.EnableInternalTxTracing,
	}
}
