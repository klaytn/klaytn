package rpc

import "github.com/rcrowley/go-metrics"

var (
	rpcTotalRequestsCounter    = metrics.NewRegisteredCounter("rpc/counts/total", nil)
	rpcSuccessResponsesCounter = metrics.NewRegisteredCounter("rpc/counts/success", nil)
	rpcErrorResponsesCounter   = metrics.NewRegisteredCounter("rpc/counts/errors", nil)
	rpcPendingRequestsCount    = metrics.NewRegisteredCounter("rpc/counts/pending", nil)
)
