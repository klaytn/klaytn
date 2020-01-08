// Copyright 2019 The klaytn Authors
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

package discover

import "github.com/rcrowley/go-metrics"

var (
	bucketEntriesGauge      = metrics.NewRegisteredGauge("discover/bucketEntries", nil)      // the closest nodes list gauge
	bucketReplacementsGauge = metrics.NewRegisteredGauge("discover/bucketReplacements", nil) // the replacement nodes list gauge (the nodes are found, but the entries are full)
	udpPacketCounter        = metrics.NewRegisteredCounter("discover/udpPacket", nil)        // the received udp packet counter
	pingMeter               = metrics.NewRegisteredMeter("discover/ping", nil)               // sending ping packet meter
	pendingPongCounter      = metrics.NewRegisteredCounter("discover/pendingPong", nil)      // pending pong packet counter at the moment
	pongMeter               = metrics.NewRegisteredMeter("discover/pong", nil)               // received pong packet meter
	findNodesMeter          = metrics.NewRegisteredMeter("discover/findnodes", nil)          // sending findnode packet meter
	pendingNeighborsCounter = metrics.NewRegisteredCounter("discover/pendingNeighbors", nil) // pending neighbors counter at the moment
	neighborsMeter          = metrics.NewRegisteredMeter("discover/neighbors", nil)          // received neighbors packet meter
	mismatchNetworkCounter  = metrics.NewRegisteredMeter("discover/mismatchNetwork", nil)    // mismatch network ping packet counter
)
