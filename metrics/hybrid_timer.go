// Copyright 2021 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package metrics

import (
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
)

const gaugeSuffix = "/maxgauge"

var mu sync.Mutex
var gauges = make(map[string]metrics.Gauge)

// ResetMaxGauges sets the value of registered gauges to 0.
func ResetMaxGauges() {
	mu.Lock()
	defer mu.Unlock()
	for _, g := range gauges {
		g.Update(0)
	}
}

// registerHybridGauge registers the given metric under the given name.
// It returns a DuplicateMetric if a metric by the given name is already registered.
func registerHybridGauge(name string, g metrics.Gauge) {
	mu.Lock()
	defer mu.Unlock()

	if _, exist := gauges[name]; exist {
		return
	}
	gauges[name] = g
}

// HybridTimer holds both metrics.Meter and metrics.Gauge to track
// meter-wise value and temporal maximum value during the certain period.
type HybridTimer interface {
	Update(d time.Duration)
}

type hybridTimer struct {
	m metrics.Meter
	g metrics.Gauge
}

// NewRegisteredHybridTimer constructs and registers a new HybridTimer.
// `name` is used by meter and `name`+"/maxgauge" is used by gauge.
func NewRegisteredHybridTimer(name string, r metrics.Registry) HybridTimer {
	meter := metrics.NewRegisteredMeter(name, r)
	gaugeName := name + gaugeSuffix

	g := metrics.NewRegisteredGauge(gaugeName, r)
	registerHybridGauge(gaugeName, g)
	return &hybridTimer{m: meter, g: g}
}

// Update updates the value of meter and gauge.
// The value of gauge is updated only if the current value
// is greater than the current value.
func (mg *hybridTimer) Update(d time.Duration) {
	if mg.g.Value() < int64(d) {
		mg.g.Update(int64(d))
	}
	mg.m.Mark(int64(d))
}
