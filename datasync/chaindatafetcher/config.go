package chaindatafetcher

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
	NumHandlers:             10,
	JobChannelSize:          50,

	BlockChannelSize: 500,

	DBPort: "3306",
}
