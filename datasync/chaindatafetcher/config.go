// Copyright 2020 The klaytn Authors
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

package chaindatafetcher

import (
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas"
)

type ChainDataFetcherMode int

const (
	ModeKAS = ChainDataFetcherMode(iota)
	ModeKafka
)

const (
	DefaultNumHandlers      = 10
	DefaultJobChannelSize   = 50
	DefaultBlockChannelSize = 500
	DefaultDBPort           = "3306"
)

//go:generate gencodec -type ChainDataFetcherConfig -formats toml -out gen_config.go
type ChainDataFetcherConfig struct {
	EnabledChainDataFetcher bool
	Mode                    ChainDataFetcherMode
	NoDefaultStart          bool
	NumHandlers             int
	JobChannelSize          int
	BlockChannelSize        int

	KasConfig   *kas.KASConfig `json:"-"` // Deprecated: This configuration is not used anymore.
	KafkaConfig *kafka.KafkaConfig
}

var DefaultChainDataFetcherConfig = &ChainDataFetcherConfig{
	EnabledChainDataFetcher: false,
	Mode:                    ModeKAS,
	NoDefaultStart:          false,
	NumHandlers:             DefaultNumHandlers,
	JobChannelSize:          DefaultJobChannelSize,
	BlockChannelSize:        DefaultBlockChannelSize,

	KasConfig:   kas.DefaultKASConfig,
	KafkaConfig: kafka.GetDefaultKafkaConfig(),
}
