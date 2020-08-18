package chaindatafetcher

import "github.com/rcrowley/go-metrics"

var (
	checkpointGauge = metrics.NewRegisteredGauge("chaindatafetcher/checkpoint/gauge", nil)

	txInsertionTimeGauge            = metrics.NewRegisteredGauge("chaindatafetcher/insertion/tx/gauge", nil)
	tokenTransferInsertionTimeGauge = metrics.NewRegisteredGauge("chaindatafetcher/insertion/tokentransfer/gauge", nil)
	contractsInsertionTimeGauge     = metrics.NewRegisteredGauge("chaindatafetcher/insertion/contracts/gauge", nil)
	tracesInsertionTimeGauge        = metrics.NewRegisteredGauge("chaindatafetcher/insertion/traces/gauge", nil)
	totalInsertionTimeGauge         = metrics.NewRegisteredGauge("chaindatafetcher/insertion/total/gauge", nil)

	handledBlockNumberGauge = metrics.NewRegisteredGauge("chaindatafetcher/handle/blocknumber/gauge", nil)

	numChainEventGauge = metrics.NewRegisteredGauge("chaindatafetcher/chainevent/gauge", nil)
	numRequestsGauge   = metrics.NewRegisteredGauge("chaindatafetcher/requests/gauge", nil)
)
