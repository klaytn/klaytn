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
	"math/big"
)

// requestType informs which data should be exported such as block, transaction, transaction log, etc.
type requestType uint

const (
	requestTypeTransaction = requestType(1) << iota
	requestTypeTokenTransfer
	requestTypeContracts
	requestTypeTraces

	requestTypeLength
)

const (
	requestTypeAll = requestTypeTransaction | requestTypeTokenTransfer | requestTypeContracts | requestTypeTraces
)

// request contains a raw block which should be handled and the type of data which should be exported.
type request struct {
	reqType requestType
	block   uint64
}

func checkRequestType(rt requestType, targetType requestType) bool {
	return rt&targetType == targetType
}

func newRequest(reqType requestType, block uint64) *request {
	return &request{
		reqType: reqType,
		block:   block,
	}
}

// response contains the result of handling the requested block including request type, block number, and error if exists.
type response struct {
	reqType     requestType
	blockNumber *big.Int
	err         error
}

func newResponse(reqType requestType, blockNumber *big.Int, err error) *response {
	return &response{
		reqType:     reqType,
		blockNumber: blockNumber,
		err:         err,
	}
}
