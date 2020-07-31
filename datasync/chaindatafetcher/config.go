package chaindatafetcher

const (
	DefaultNumHandlers      = 10
	DefaultJobChannelSize   = 50
	DefaultBlockChannelSize = 500
	DefaultDBPort           = "3306"
)

//go:generate gencodec -type ChainDataFetcherConfig -formats toml -out gen_config.go
type ChainDataFetcherConfig struct {
	EnabledChainDataFetcher bool
	NoDefaultStart          bool
	NumHandlers             int
	JobChannelSize          int
	BlockChannelSize        int

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

var DefaultChainDataFetcherConfig = &ChainDataFetcherConfig{
	EnabledChainDataFetcher: false,
	NoDefaultStart:          false,
	NumHandlers:             DefaultNumHandlers,
	JobChannelSize:          DefaultJobChannelSize,
	BlockChannelSize:        DefaultBlockChannelSize,

	DBPort: DefaultDBPort,
}
