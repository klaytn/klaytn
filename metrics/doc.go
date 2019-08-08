// Copyright 2018 The klaytn Authors
//
// This file is derived from metrics/metrics.go (2018/06/04).
// See LICENSE in the metrics directory for the original copyright and license.

/*
Go port of Coda Hale's Metrics library: <https://github.com/rcrowley/go-metrics>

Coda Hale's original work: <https://github.com/codahale/metrics>

Package metrics provides various kinds of metrics which can be used to collect
system statistics. Collected statistics can be exported via InfluxDB, Librato or
Prometheus. Different types of metrics have different characteristics.

  - Gauges: hold an int64 value that can be set arbitrarily.

  - EWMAs: continuously calculate an exponentially-weighted moving average based on an outside source of clock ticks.

  - Histograms: calculate distribution statistics from a series of int64 values.

  - Meters: count events to produce exponentially-weighted moving average rates and a mean rate.

  - Registries: hold references to a set of metrics by name and can iterate over them.

  - Timers: capture the duration and the rate of events.

*/
package metrics
