// Copyright 2018 The klaytn Authors
//
// This file is derived from metrics/metrics.go (2018/06/04).
// See LICENSE in the top directory for the original copyright and license.

package metricutils

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/klaytn/klaytn/log"
	prometheusmetrics "github.com/klaytn/klaytn/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
	"gopkg.in/urfave/cli.v1"
)

// Enabled is checked by the constructor functions for all of the
// standard metrics.  If it is true, the metric returned is a stub.
//
// This global kill-switch helps quantify the observer effect and makes
// for less cluttered pprof profiles.
var Enabled = false
var EnabledPrometheusExport = false
var logger = log.NewModuleLogger(log.Metrics)

// MetricsEnabledFlag is the CLI flag name to use to enable metrics collections.
const MetricsEnabledFlag = "metrics"
const DashboardEnabledFlag = "dashboard"
const PrometheusExporterFlag = "prometheus"
const PrometheusExporterPortFlag = "prometheusport"

// Init enables or disables the metrics system. Since we need this to run before
// any other code gets to create meters and timers, we'll actually do an ugly hack
// and peek into the command line args for the metrics flag.
func init() {
	for _, arg := range os.Args {
		if flag := strings.TrimLeft(arg, "-"); flag == MetricsEnabledFlag || flag == DashboardEnabledFlag {
			Enabled = true
		}
		if flag := strings.TrimLeft(arg, "-"); flag == PrometheusExporterFlag {
			EnabledPrometheusExport = true
		}
	}
}

// StartMetricCollectionAndExport starts exporting to prometheus and collects process metrics.
func StartMetricCollectionAndExport(ctx *cli.Context) {
	metricsCollectionInterval := 3 * time.Second
	if Enabled {
		logger.Info("Enabling metrics collection")
		if EnabledPrometheusExport {
			logger.Info("Enabling Prometheus Exporter")
			pClient := prometheusmetrics.NewPrometheusProvider(metrics.DefaultRegistry, "klaytn",
				"", prometheus.DefaultRegisterer, metricsCollectionInterval)
			go pClient.UpdatePrometheusMetrics()
			http.Handle("/metrics", promhttp.Handler())
			port := ctx.GlobalInt(PrometheusExporterPortFlag)

			go func() {
				err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
				if err != nil {
					logger.Error("PrometheusExporter starting failed:", "port", port, "err", err)
				}
			}()
		}
	}
	go CollectProcessMetrics(metricsCollectionInterval)
}

// CollectProcessMetrics periodically collects various metrics about the running process.
func CollectProcessMetrics(refresh time.Duration) {
	// Short circuit if the metrics system is disabled
	if !Enabled {
		return
	}
	// Create the various data collectors
	memstats := make([]*runtime.MemStats, 2)
	diskstats := make([]*DiskStats, 2)
	for i := 0; i < len(memstats); i++ {
		memstats[i] = new(runtime.MemStats)
		diskstats[i] = new(DiskStats)
	}
	// Define the various metrics to collect
	memAllocs := metrics.GetOrRegisterMeter("system/memory/allocs", metrics.DefaultRegistry)
	memFrees := metrics.GetOrRegisterMeter("system/memory/frees", metrics.DefaultRegistry)
	memInuse := metrics.GetOrRegisterMeter("system/memory/inuse", metrics.DefaultRegistry)
	memPauses := metrics.GetOrRegisterMeter("system/memory/pauses", metrics.DefaultRegistry)

	var diskReads, diskReadBytes, diskWrites, diskWriteBytes metrics.Meter
	if err := ReadDiskStats(diskstats[0]); err == nil {
		diskReads = metrics.GetOrRegisterMeter("system/disk/readcount", metrics.DefaultRegistry)
		diskReadBytes = metrics.GetOrRegisterMeter("system/disk/readdata", metrics.DefaultRegistry)
		diskWrites = metrics.GetOrRegisterMeter("system/disk/writecount", metrics.DefaultRegistry)
		diskWriteBytes = metrics.GetOrRegisterMeter("system/disk/writedata", metrics.DefaultRegistry)
	} else {
		logger.Debug("Failed to read disk metrics", "err", err)
	}
	// Iterate loading the different stats and updating the meters
	for i := 1; ; i++ {
		runtime.ReadMemStats(memstats[i%2])
		memAllocs.Mark(int64(memstats[i%2].Mallocs - memstats[(i-1)%2].Mallocs))
		memFrees.Mark(int64(memstats[i%2].Frees - memstats[(i-1)%2].Frees))
		memInuse.Mark(int64(memstats[i%2].Alloc - memstats[(i-1)%2].Alloc))
		memPauses.Mark(int64(memstats[i%2].PauseTotalNs - memstats[(i-1)%2].PauseTotalNs))

		if ReadDiskStats(diskstats[i%2]) == nil {
			diskReads.Mark(diskstats[i%2].ReadCount - diskstats[(i-1)%2].ReadCount)
			diskReadBytes.Mark(diskstats[i%2].ReadBytes - diskstats[(i-1)%2].ReadBytes)
			diskWrites.Mark(diskstats[i%2].WriteCount - diskstats[(i-1)%2].WriteCount)
			diskWriteBytes.Mark(diskstats[i%2].WriteBytes - diskstats[(i-1)%2].WriteBytes)
		}
		time.Sleep(refresh)
	}
}
