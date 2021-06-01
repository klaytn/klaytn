package rpc

import "github.com/rcrowley/go-metrics"

var (
	rpcTotalRequestsCounter    = metrics.NewRegisteredCounter("rpc/counts/total", nil)
	rpcSuccessResponsesCounter = metrics.NewRegisteredCounter("rpc/counts/success", nil)
	rpcErrorResponsesCounter   = metrics.NewRegisteredCounter("rpc/counts/errors", nil)
	rpcPendingRequestsCount    = metrics.NewRegisteredCounter("rpc/counts/pending", nil)

	wsSubscriptionReqCounter   = metrics.NewRegisteredCounter("ws/counts/subscription/request", nil)
	wsUnsubscriptionReqCounter = metrics.NewRegisteredCounter("ws/counts/unsubscription/request", nil)
	wsConnCounter              = metrics.NewRegisteredCounter("ws/counts/connections/total", nil)
)
