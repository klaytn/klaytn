package chaindatafetcher

import (
	"github.com/klaytn/klaytn/blockchain"
	"math/big"
)

type requestType uint

const (
	requestTypeTransaction = requestType(1) << iota
)

type request struct {
	reqType requestType
	event   blockchain.ChainEvent
}

type response struct {
	reqType     requestType
	blockNumber *big.Int
	err         error
}
