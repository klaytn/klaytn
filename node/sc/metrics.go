// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/metrics.go (2018/06/04).
// Modified and improved for the klaytn development.

package sc

import (
	metricutils "github.com/klaytn/klaytn/metrics/utils"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/rcrowley/go-metrics"
)

var (
	propTxnInPacketsMeter     = metrics.NewRegisteredMeter("klay/bridge/prop/txns/in/packets", nil)
	propTxnInTrafficMeter     = metrics.NewRegisteredMeter("klay/bridge/prop/txns/in/traffic", nil)
	propTxnOutPacketsMeter    = metrics.NewRegisteredMeter("klay/bridge/prop/txns/out/packets", nil)
	propTxnOutTrafficMeter    = metrics.NewRegisteredMeter("klay/bridge/prop/txns/out/traffic", nil)
	reqReceiptInPacketsMeter  = metrics.NewRegisteredMeter("klay/bridge/req/receipts/in/packets", nil)
	reqReceiptInTrafficMeter  = metrics.NewRegisteredMeter("klay/bridge/req/receipts/in/traffic", nil)
	reqReceiptOutPacketsMeter = metrics.NewRegisteredMeter("klay/bridge/req/receipts/out/packets", nil)
	reqReceiptOutTrafficMeter = metrics.NewRegisteredMeter("klay/bridge/req/receipts/out/traffic", nil)
	miscInPacketsMeter        = metrics.NewRegisteredMeter("klay/bridge/misc/in/packets", nil)
	miscInTrafficMeter        = metrics.NewRegisteredMeter("klay/bridge/misc/in/traffic", nil)
	miscOutPacketsMeter       = metrics.NewRegisteredMeter("klay/bridge/misc/out/packets", nil)
	miscOutTrafficMeter       = metrics.NewRegisteredMeter("klay/bridge/misc/out/traffic", nil)

	vtRequestEventMeter = metrics.NewRegisteredMeter("klay/bridge/vt/event/request", nil)
	vtHandleEventMeter  = metrics.NewRegisteredMeter("klay/bridge/vt/event/handle", nil)

	vtRecoveredRequestEventMeter = metrics.NewRegisteredMeter("klay/bridge/vt/event/recovery/request", nil)
	vtPendingRequestEventCounter = metrics.NewRegisteredCounter("klay/bridge/vt/event/pend/request", nil)

	vtRequestNonceCount     = metrics.NewRegisteredCounter("klay/bridge/vt/nonce/request", nil)
	vtHandleNonceCount      = metrics.NewRegisteredCounter("klay/bridge/vt/nonce/handle", nil)
	vtLowerHandleNonceCount = metrics.NewRegisteredCounter("klay/bridge/vt/nonce/lowerhandle", nil)

	lastAnchoredBlockNumGauge = metrics.NewRegisteredGauge("klay/bridge/anchroing/blocknumber", nil)

	// TODO-Klaytn-Servicechain need to add below metrics
	//txReceiveCounter     = metrics.NewRegisteredCounter("klay/bridge/tx/recv/counter", nil)
	//txResendCounter      = metrics.NewRegisteredCounter("klay/bridge/tx/resend/counter", nil)
	//txResendGauge        = metrics.NewRegisteredGauge("klay/bridge/tx/resend/gauge", nil)
	//txSendCounter        = metrics.NewRegisteredCounter("klay/bridge/tx/send/counter", nil)
	//txResendRoutineGauge = metrics.NewRegisteredGauge("klay/bridge/tx/resend/routine/gauge", nil)
)

// meteredMsgReadWriter is a wrapper around a p2p.MsgReadWriter, capable of
// accumulating the above defined metrics based on the data stream contents.
type meteredMsgReadWriter struct {
	p2p.MsgReadWriter     // Wrapped message stream to meter
	version           int // Protocol version to select correct meters
}

// newMeteredMsgWriter wraps a p2p MsgReadWriter with metering support. If the
// metrics system is disabled, this function returns the original object.
func newMeteredMsgWriter(rw p2p.MsgReadWriter) p2p.MsgReadWriter {
	if !metricutils.Enabled {
		return rw
	}
	return &meteredMsgReadWriter{MsgReadWriter: rw}
}

// Init sets the protocol version used by the stream to know which meters to
// increment in case of overlapping message ids between protocol versions.
func (rw *meteredMsgReadWriter) Init(version int) {
	rw.version = version
}

func (rw *meteredMsgReadWriter) ReadMsg() (p2p.Msg, error) {
	// Read the message and short circuit in case of an error
	msg, err := rw.MsgReadWriter.ReadMsg()
	if err != nil {
		return msg, err
	}
	// Account for the data traffic
	packets, traffic := miscInPacketsMeter, miscInTrafficMeter
	switch {
	case msg.Code == ServiceChainTxsMsg: // If version check is needed, add `rw.version >= klay63`.
		packets, traffic = propTxnInPacketsMeter, propTxnInTrafficMeter
	case msg.Code == ServiceChainReceiptResponseMsg:
		packets, traffic = reqReceiptInPacketsMeter, reqReceiptInTrafficMeter
	case msg.Code == ServiceChainReceiptRequestMsg:
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	return msg, err
}

func (rw *meteredMsgReadWriter) WriteMsg(msg p2p.Msg) error {
	// Account for the data traffic
	packets, traffic := miscOutPacketsMeter, miscOutTrafficMeter
	switch {
	case msg.Code == ServiceChainTxsMsg: // If version check is needed, add `rw.version >= klay63`.
		packets, traffic = propTxnOutPacketsMeter, propTxnOutTrafficMeter
	case msg.Code == ServiceChainReceiptResponseMsg:
		packets, traffic = reqReceiptOutPacketsMeter, reqReceiptOutTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	// Send the packet to the p2p layer
	return rw.MsgReadWriter.WriteMsg(msg)
}
