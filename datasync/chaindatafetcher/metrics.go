package chaindatafetcher

import "github.com/rcrowley/go-metrics"

var (
	checkpointGauge = metrics.NewRegisteredGauge("chaindatafetcher/checkpoint/gauge", nil)

	totalInsertionTimeGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/total/gauge", nil)

	// KAS specific metrics
	txsInsertionTimeGauge            = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/txs/gauge", nil)
	tokenTransfersInsertionTimeGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/tokentransfers/gauge", nil)
	contractsInsertionTimeGauge      = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/contracts/gauge", nil)
	tracesInsertionTimeGauge         = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/traces/gauge", nil)

	txsInsertionRetryGauge            = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/txs/gauge", nil)
	tokenTransfersInsertionRetryGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/tokentransfers/gauge", nil)
	contractsInsertionRetryGauge      = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/contracts/gauge", nil)
	tracesInsertionRetryGauge         = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/traces/gauge", nil)

	// Kafka specific metrics
	blockGroupInsertionTimeGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/blockgroup/gauge", nil)
	traceGroupInsertionTimeGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/tracegroup/gauge", nil)

	blockGroupInsertionRetryGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/blockgroup/gauge", nil)
	traceGroupInsertionRetryGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/tracegroup/gauge", nil)

	handledBlockNumberGauge = metrics.NewRegisteredGauge("chaindatafetcher/handle/blocknumber/gauge", nil)

	numChainEventGauge = metrics.NewRegisteredGauge("chaindatafetcher/chainevent/gauge", nil)
	numRequestsGauge   = metrics.NewRegisteredGauge("chaindatafetcher/requests/gauge", nil)

	traceAPIErrorCounter = metrics.NewRegisteredCounter("chaindatafetcher/trace/error", nil)
)
