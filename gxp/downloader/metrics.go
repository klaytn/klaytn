package downloader

import "ground-x/go-gxplatform/matrics"

var (
	headerInMeter      = metrics.NewRegisteredMeter("gxp/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("gxp/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("gxp/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("gxp/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("gxp/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("gxp/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("gxp/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("gxp/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("gxp/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("gxp/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("gxp/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("gxp/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("gxp/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("gxp/downloader/states/drop", nil)
)
