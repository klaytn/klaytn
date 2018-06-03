package gxp

import (
	"time"
	"math/big"
	"ground-x/go-gxplatform/params"
	"ground-x/go-gxplatform/core"
	"ground-x/go-gxplatform/common"
	"os"
	"os/user"
	"runtime"
	"path/filepath"
	"ground-x/go-gxplatform/common/hexutil"
	"ground-x/go-gxplatform/consensus/gxhash"
	"ground-x/go-gxplatform/gxp/gasprice"
	"ground-x/go-gxplatform/gxp/downloader"
)

// DefaultConfig contains default settings for use on the GxPlatform main net.
var DefaultConfig = Config{
	//SyncMode: downloader.FastSync,
	Gxhash: gxhash.Config{
		CacheDir:       "gxhash",
		CachesInMem:    2,
		CachesOnDisk:   3,
		DatasetsInMem:  1,
		DatasetsOnDisk: 2,
	},
	NetworkId:     1,
	LightPeers:    100,
	DatabaseCache: 768,
	TrieCache:     256,
	TrieTimeout:   5 * time.Minute,
	GasPrice:      big.NewInt(18 * params.Shannon),

	TxPool: core.DefaultTxPoolConfig,
	GPO: gasprice.Config{
		Blocks:     20,
		Percentile: 60,
	},
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
	if runtime.GOOS == "windows" {
		DefaultConfig.Gxhash.DatasetDir = filepath.Join(home, "AppData", "Gxhash")
	} else {
		DefaultConfig.Gxhash.DatasetDir = filepath.Join(home, ".gxhash")
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the GXP main net block is used.
	Genesis *core.Genesis `toml:",omitempty"`

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to
	SyncMode  downloader.SyncMode
	NoPruning bool

	// Light client options
	LightServ  int `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightPeers int `toml:",omitempty"` // Maximum number of LES client peers

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	TrieCache          int
	TrieTimeout        time.Duration

	// Mining-related options
	Gxbase    common.Address 	`toml:",omitempty"`
	MinerThreads int            `toml:",omitempty"`
	ExtraData    []byte         `toml:",omitempty"`
	GasPrice     *big.Int

	// Gxhash options
	Gxhash gxhash.Config

	// Transaction pool options
	TxPool core.TxPoolConfig

	// Gas Price Oracle options
	GPO gasprice.Config

	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`
}

type configMarshaling struct {
	ExtraData hexutil.Bytes
}
