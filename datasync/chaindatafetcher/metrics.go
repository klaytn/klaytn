package chaindatafetcher

import "github.com/rcrowley/go-metrics"

var (
	checkpointGauge = metrics.NewRegisteredGauge("chaindatafetcher/checkpoint/gauge", nil)

	txInsertionTimeGauge            = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/tx/gauge", nil)
	tokenTransferInsertionTimeGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/tokentransfer/gauge", nil)
	contractsInsertionTimeGauge     = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/contracts/gauge", nil)
	tracesInsertionTimeGauge        = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/traces/gauge", nil)
	totalInsertionTimeGauge         = metrics.NewRegisteredGauge("chaindatafetcher/insertion/time/total/gauge", nil)

	txInsertionRetryGauge            = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/tx/gauge", nil)
	tokenTransferInsertionRetryGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/tokentransfer/gauge", nil)
	contractsInsertionRetryGauge     = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/contracts/gauge", nil)
	tracesInsertionRetryGauge        = metrics.NewRegisteredGauge("chaindatafetcher/insertion/retry/traces/gauge", nil)

	handledBlockNumberGauge = metrics.NewRegisteredGauge("chaindatafetcher/handle/blocknumber/gauge", nil)

	numChainEventGauge = metrics.NewRegisteredGauge("chaindatafetcher/chainevent/gauge", nil)
	numRequestsGauge   = metrics.NewRegisteredGauge("chaindatafetcher/requests/gauge", nil)
)
