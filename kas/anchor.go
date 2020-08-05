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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"math/big"
	"net/http"
)

var (
	codeOK                = 0
	errNotFoundBlock      = errors.New("not found block")
	errInvalidBlockNumber = errors.New("invalid block number")
)

//go:generate mockgen -destination=./mocks/anchordb_mock.go -package=mocks github.com/klaytn/klaytn/kas AnchorDB
type AnchorDB interface {
	WriteAnchoredBlockNumber(blockNum uint64)
	ReadAnchoredBlockNumber() uint64
}

//go:generate mockgen -destination=./mocks/blockchain_mock.go -package=mocks github.com/klaytn/klaytn/kas BlockChain
type BlockChain interface {
	GetBlockByNumber(number uint64) *types.Block
}

//go:generate mockgen -destination=./mocks/client_mock.go -package=mocks github.com/klaytn/klaytn/kas HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Anchor struct {
	kasConfig *KASConfig
	db        AnchorDB
	bc        BlockChain
	client    HTTPClient
}

func NewKASAnchor(kasConfig *KASConfig, db AnchorDB, bc BlockChain) *Anchor {
	return &Anchor{
		kasConfig: kasConfig,
		db:        db,
		bc:        bc,
		client:    &http.Client{},
	}
}


type respBody struct {
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}

type reqBody struct {
	Operator common.Address `json:"operator"`
	Payload  interface{}    `json:"Payload"`
}

type Payload struct {
	Id string `json:"id"`
	types.AnchoringDataInternalType0
}

func dataToPayload(anchorData *types.AnchoringDataInternalType0) *Payload {
	payload := &Payload{
		Id:                         anchorData.BlockNumber.String(),
		AnchoringDataInternalType0: *anchorData,
	}

	return payload
}
func (anchor *Anchor) sendRequest(payload interface{}) (*respBody, error) {
	header := map[string]string{
		"Content-Type": "application/json",
		"X-Krn":        anchor.kasConfig.Xkrn,
	}

	bodyData := reqBody{
		Operator: anchor.kasConfig.Operator,
		Payload:  payload,
	}

	bodyDataBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(bodyDataBytes)

	req, err := http.NewRequest("POST", anchor.kasConfig.Url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(anchor.kasConfig.User, anchor.kasConfig.Pwd)
	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := anchor.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := respBody{}
	json.NewDecoder(resp.Body).Decode(&v)

	return &v, nil
}
