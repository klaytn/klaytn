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

package kas

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/log"
	"time"
)

const (
	maxTxCountPerBlock      = int64(1000000)
	maxTxLogCountPerTx      = int64(100000)
	maxInternalTxCountPerTx = int64(10000)

	maxPlaceholders = 65535

	placeholdersPerTxItem          = 13
	placeholdersPerKCTTransferItem = 7
	placeholdersPerRevertedTxItem  = 5
	placeholdersPerContractItem    = 1
	placeholdersPerFTItem          = 3
	placeholdersPerNFTItem         = 3

	maxDBRetryCount = 20
	DBRetryInterval = 1 * time.Second
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

type repository struct {
	db *gorm.DB

	contractCaller *contractCaller
	blockchainApi  *api.PublicBlockChainAPI
}

func getEndpoint(user, password, host, port, name string) string {
	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?parseTime=True&charset=utf8mb4"
}

func NewRepository(user, password, host, port, name string) (*repository, error) {
	endpoint := getEndpoint(user, password, host, port, name)
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < maxDBRetryCount; i++ {
		db, err = gorm.Open("mysql", endpoint)
		if err != nil {
			logger.Warn("Retrying to connect DB", "endpoint", endpoint, "err", err)
			time.Sleep(DBRetryInterval)
		} else {
			// TODO-ChainDataFetcher insert other options such as maxOpen, maxIdle, maxLifetime, etc.
			//db.DB().SetMaxOpenConns(maxOpen)
			//db.DB().SetMaxIdleConns(maxIdle)
			//db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

			return &repository{db: db}, nil
		}
	}
	logger.Error("Failed to connect to the database", "endpoint", endpoint, "err", err)
	return nil, err
}

func (r *repository) SetComponent(component interface{}) {
	switch c := component.(type) {
	case *api.PublicBlockChainAPI:
		r.blockchainApi = c
	}
}
