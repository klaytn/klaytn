package chaindatafetcher

import (
	"github.com/klaytn/klaytn/blockchain"
	"math/big"
)

// requestType informs which data should be exported such as block, transaction, transaction log, etc.
type requestType uint

const (
	requestTypeTransaction = requestType(1) << iota
)

// request contains a raw block which should be handled and the type of data which should be exported.
type request struct {
	reqType requestType
	event   blockchain.ChainEvent
}

// response contains the result of handling the requested block including request type, block number, and error if exists.
type response struct {
	reqType     requestType
	blockNumber *big.Int
	err         error
}
