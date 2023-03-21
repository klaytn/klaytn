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

package types

// RequestType informs which data should be exported such as block, transaction, transaction log, etc.
type RequestType uint

const (
	// RequestTypes for KAS (deprecated)
	RequestTypeTransaction = RequestType(1) << iota
	RequestTypeTokenTransfer
	RequestTypeContract
	RequestTypeTrace

	// RequestTypes for Kafka
	RequestTypeBlockGroup
	RequestTypeTraceGroup

	RequestTypeLength
)

const (
	RequestTypeAll      = RequestTypeTransaction | RequestTypeTokenTransfer | RequestTypeContract | RequestTypeTrace
	RequestTypeGroupAll = RequestTypeBlockGroup | RequestTypeTraceGroup
)

func (t RequestType) IsValid() bool {
	return t == RequestTypeBlockGroup || t == RequestTypeTraceGroup || t == RequestTypeGroupAll
}

func (t RequestType) String() string {
	switch t {
	case RequestTypeGroupAll:
		return "all"
	case RequestTypeBlockGroup:
		return "block"
	case RequestTypeTraceGroup:
		return "trace"
	default:
		return "unknown"
	}
}

// Request contains a blockNumber which should be handled and the type of data which should be exported.
type Request struct {
	ReqType                RequestType
	ShouldUpdateCheckpoint bool
	BlockNumber            uint64
}

func CheckRequestType(rt RequestType, targetType RequestType) bool {
	return rt&targetType == targetType
}

func NewRequest(reqType RequestType, shouldUpdateCheckpoint bool, block uint64) *Request {
	return &Request{
		ReqType:                reqType,
		ShouldUpdateCheckpoint: shouldUpdateCheckpoint,
		BlockNumber:            block,
	}
}
