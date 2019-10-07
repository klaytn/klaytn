package rpc

import "github.com/klaytn/klaytn/metrics"

var (
	rpcTotalRequestsCounter    = metrics.NewRegisteredCounter("rpc/counts/total", nil)
	rpcSuccessResponsesCounter = metrics.NewRegisteredCounter("rpc/counts/success", nil)
	rpcErrorResponsesCounter   = metrics.NewRegisteredCounter("rpc/counts/errors", nil)
	rpcPendingRequestCount     = metrics.NewRegisteredCounter("rpc/counts/pending", nil)
)
