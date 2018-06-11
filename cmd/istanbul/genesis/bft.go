package genesis

import (
	"math/big"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/core"
	"ground-x/go-gxplatform/params"
	"ground-x/go-gxplatform/common/math"
	"ground-x/go-gxplatform/common/hexutil"
)

//go:generate gencodec -type BFTGenesis -field-override genesisSpecMarshaling -out gen_bft_genesis.go

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce      math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Difficulty *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]core.GenesisAccount
}

type BFTChainConfig struct {
	*params.ChainConfig
	IsBFT bool `json:"isBFT,omitempty"`
}

type BFTGenesis struct {
	Config     *BFTChainConfig `json:"config"`
	Nonce      uint64             `json:"nonce"`
	Timestamp  uint64             `json:"timestamp"`
	ExtraData  []byte             `json:"extraData"`
	GasLimit   uint64             `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int           `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash        `json:"mixHash"`
	Coinbase   common.Address     `json:"coinbase"`
	Alloc      core.GenesisAlloc  `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

// ToBFT converts standard genesis to bft genesis
func ToBFT(g *core.Genesis, isBFT bool) *BFTGenesis {
	return &BFTGenesis{
		Config: &BFTChainConfig{
			ChainConfig: g.Config,
			IsBFT:    isBFT,
		},
		Nonce:      g.Nonce,
		Timestamp:  g.Timestamp,
		ExtraData:  g.ExtraData,
		GasLimit:   g.GasLimit,
		Difficulty: g.Difficulty,
		Mixhash:    g.Mixhash,
		Coinbase:   g.Coinbase,
		Alloc:      g.Alloc,
		Number:     g.Number,
		GasUsed:    g.GasUsed,
		ParentHash: g.ParentHash,
	}
}
