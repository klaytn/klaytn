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

package dbsyncer

import "time"

const (
	BLOCK_MODE = "block"
	HEAD_MODE  = "head"
)

//go:generate gencodec -type DBConfig -formats toml -out gen_config.go
type DBConfig struct {
	EnabledDBSyncer bool
	EnabledLogMode  bool

	// DB Config
	DBHost     string `toml:",omitempty"`
	DBPort     string `toml:",omitempty"`
	DBUser     string `toml:",omitempty"`
	DBPassword string `toml:",omitempty"`
	DBName     string `toml:",omitempty"`

	MaxIdleConns     int           `toml:",omitempty"`
	MaxOpenConns     int           `toml:",omitempty"`
	ConnMaxLifetime  time.Duration `toml:",omitempty"`
	BlockChannelSize int           `toml:",omitempty"`

	GenQueryThread int `toml:",omitempty"`
	InsertThread   int `toml:",omitempty"`

	BulkInsertSize int `toml:",omitempty"`

	Mode      string `toml:",omitempty"`
	EventMode string `toml:",omitempty"`

	MaxBlockDiff uint64 `toml:",omitempty"`
}

var DefaultDBConfig = &DBConfig{

	EnabledDBSyncer: false,
	EnabledLogMode:  false,

	DBPort: "3306",

	MaxIdleConns:     50,
	MaxOpenConns:     30,
	ConnMaxLifetime:  1 * time.Hour,
	BlockChannelSize: 5,

	GenQueryThread: 100,
	InsertThread:   30,

	BulkInsertSize: 200,

	Mode:      "multi",
	EventMode: HEAD_MODE,

	MaxBlockDiff: 0,
}
