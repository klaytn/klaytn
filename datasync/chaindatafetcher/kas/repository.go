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
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
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

	maxOpenConnection = 100
	maxIdleConnection = 10
	connMaxLifetime   = 24 * time.Hour
	maxDBRetryCount   = 20
	DBRetryInterval   = 1 * time.Second
	apiCtxTimeout     = 200 * time.Millisecond
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

type repository struct {
	db *gorm.DB

	config *KASConfig

	contractCaller *contractCaller
	blockchainApi  BlockchainAPI
}

func getEndpoint(user, password, host, port, name string) string {
	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?parseTime=True&charset=utf8mb4"
}

func NewRepository(config *KASConfig) (*repository, error) {
	endpoint := getEndpoint(config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
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
			db.DB().SetMaxOpenConns(maxOpenConnection)
			db.DB().SetMaxIdleConns(maxIdleConnection)
			db.DB().SetConnMaxLifetime(connMaxLifetime)

			return &repository{db: db, config: config}, nil
		}
	}
	logger.Error("Failed to connect to the database", "endpoint", endpoint, "err", err)
	return nil, err
}

func (r *repository) setBlockchainAPI(apis []rpc.API) {
	for _, a := range apis {
		switch s := a.Service.(type) {
		case *api.PublicBlockChainAPI:
			r.blockchainApi = s
		}
	}
}

func (r *repository) SetComponent(component interface{}) {
	switch c := component.(type) {
	case []rpc.API:
		r.setBlockchainAPI(c)
	}
}

func (r *repository) HandleChainEvent(event blockchain.ChainEvent, reqType types.RequestType) error {
	switch reqType {
	case types.RequestTypeTransaction:
		return r.InsertTransactions(event)
	case types.RequestTypeTrace:
		return r.InsertTraceResults(event)
	case types.RequestTypeTokenTransfer:
		return r.InsertTokenTransfers(event)
	case types.RequestTypeContract:
		return r.InsertContracts(event)
	default:
		return fmt.Errorf("unsupported data type. [blockNumber: %v, reqType: %v]", event.Block.NumberU64(), reqType)
	}
}

func makeEOAListStr(eoaList map[common.Address]struct{}) string {
	eoaStrList := ""
	for key := range eoaList {
		eoaStrList += "\""
		eoaStrList += strings.ToLower(key.String())
		eoaStrList += "\""
		eoaStrList += ","
	}
	return eoaStrList[:len(eoaStrList)-1]
}

func makeBasicAuthWithParam(param string) string {
	return "Basic " + param
}

func (r *repository) InvalidateCacheEOAList(eoaList map[common.Address]struct{}) {
	numEOAs := len(eoaList)
	logger.Trace("the number of EOA list for KAS cache invalidation", "numEOAs", numEOAs)
	if numEOAs <= 0 || !r.config.CacheUse {
		return
	}

	url := r.config.CacheInvalidationURL
	payloadStr := fmt.Sprintf(`{"type": "stateChange","payload": {"addresses": [%v]}}`, makeEOAListStr(eoaList))
	payload := strings.NewReader(payloadStr)

	// set up timeout for API call
	ctx, cancel := context.WithTimeout(context.Background(), apiCtxTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		logger.Error("Creating a new http request is failed", "err", err, "url", url, "payload", payloadStr)
		return
	}
	req.Header.Add("x-chain-id", r.config.XChainId)
	req.Header.Add("Authorization", makeBasicAuthWithParam(r.config.BasicAuthParam))
	req.Header.Add("Content-Type", "text/plain")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Client do method is failed", "err", err, "url", url, "payload", payloadStr, "header", req.Header)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Error("Reading response body is failed", "err", err, "body", res.Body)
			return
		}
		logger.Error("cache invalidation is failed", "response", string(body))
	}
}
