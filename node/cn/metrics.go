// Modifications Copyright 2018 The klaytn Authors
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

package cn

import (
	"github.com/klaytn/klaytn/consensus/istanbul/backend"
	metricutils "github.com/klaytn/klaytn/metrics/utils"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/rcrowley/go-metrics"
)

var (
	propTxnInPacketsMeter                = metrics.NewRegisteredMeter("klay/prop/txns/in/packets", nil)
	propTxnInTrafficMeter                = metrics.NewRegisteredMeter("klay/prop/txns/in/traffic", nil)
	propTxnOutPacketsMeter               = metrics.NewRegisteredMeter("klay/prop/txns/out/packets", nil)
	propTxnOutTrafficMeter               = metrics.NewRegisteredMeter("klay/prop/txns/out/traffic", nil)
	propTxPeersGauge                     = metrics.NewRegisteredGauge("klay/prop/tx/peers/gauge", nil)
	propHashInPacketsMeter               = metrics.NewRegisteredMeter("klay/prop/hashes/in/packets", nil)
	propHashInTrafficMeter               = metrics.NewRegisteredMeter("klay/prop/hashes/in/traffic", nil)
	propHashOutPacketsMeter              = metrics.NewRegisteredMeter("klay/prop/hashes/out/packets", nil)
	propHashOutTrafficMeter              = metrics.NewRegisteredMeter("klay/prop/hashes/out/traffic", nil)
	propBlockInPacketsMeter              = metrics.NewRegisteredMeter("klay/prop/blocks/in/packets", nil)
	propBlockInTrafficMeter              = metrics.NewRegisteredMeter("klay/prop/blocks/in/traffic", nil)
	propBlockOutPacketsMeter             = metrics.NewRegisteredMeter("klay/prop/blocks/out/packets", nil)
	propBlockOutTrafficMeter             = metrics.NewRegisteredMeter("klay/prop/blocks/out/traffic", nil)
	reqHeaderInPacketsMeter              = metrics.NewRegisteredMeter("klay/req/headers/in/packets", nil)
	reqHeaderInTrafficMeter              = metrics.NewRegisteredMeter("klay/req/headers/in/traffic", nil)
	reqHeaderOutPacketsMeter             = metrics.NewRegisteredMeter("klay/req/headers/out/packets", nil)
	reqHeaderOutTrafficMeter             = metrics.NewRegisteredMeter("klay/req/headers/out/traffic", nil)
	reqBodyInPacketsMeter                = metrics.NewRegisteredMeter("klay/req/bodies/in/packets", nil)
	reqBodyInTrafficMeter                = metrics.NewRegisteredMeter("klay/req/bodies/in/traffic", nil)
	reqBodyOutPacketsMeter               = metrics.NewRegisteredMeter("klay/req/bodies/out/packets", nil)
	reqBodyOutTrafficMeter               = metrics.NewRegisteredMeter("klay/req/bodies/out/traffic", nil)
	reqStateInPacketsMeter               = metrics.NewRegisteredMeter("klay/req/states/in/packets", nil)
	reqStateInTrafficMeter               = metrics.NewRegisteredMeter("klay/req/states/in/traffic", nil)
	reqStateOutPacketsMeter              = metrics.NewRegisteredMeter("klay/req/states/out/packets", nil)
	reqStateOutTrafficMeter              = metrics.NewRegisteredMeter("klay/req/states/out/traffic", nil)
	reqReceiptInPacketsMeter             = metrics.NewRegisteredMeter("klay/req/receipts/in/packets", nil)
	reqReceiptInTrafficMeter             = metrics.NewRegisteredMeter("klay/req/receipts/in/traffic", nil)
	reqReceiptOutPacketsMeter            = metrics.NewRegisteredMeter("klay/req/receipts/out/packets", nil)
	reqReceiptOutTrafficMeter            = metrics.NewRegisteredMeter("klay/req/receipts/out/traffic", nil)
	miscInPacketsMeter                   = metrics.NewRegisteredMeter("klay/misc/in/packets", nil)
	miscInTrafficMeter                   = metrics.NewRegisteredMeter("klay/misc/in/traffic", nil)
	miscOutPacketsMeter                  = metrics.NewRegisteredMeter("klay/misc/out/packets", nil)
	miscOutTrafficMeter                  = metrics.NewRegisteredMeter("klay/misc/out/traffic", nil)
	txReceiveCounter                     = metrics.NewRegisteredCounter("klay/tx/recv/counter", nil)
	txResendCounter                      = metrics.NewRegisteredCounter("klay/tx/resend/counter", nil)
	txSendCounter                        = metrics.NewRegisteredCounter("klay/tx/send/counter", nil)
	txResendRoutineGauge                 = metrics.NewRegisteredGauge("klay/tx/resend/routine/gauge", nil)
	cnPeerCountGauge                     = metrics.NewRegisteredGauge("p2p/CNPeerCountGauge", nil)
	pnPeerCountGauge                     = metrics.NewRegisteredGauge("p2p/PNPeerCountGauge", nil)
	enPeerCountGauge                     = metrics.NewRegisteredGauge("p2p/ENPeerCountGauge", nil)
	propConsensusIstanbulInPacketsMeter  = metrics.NewRegisteredMeter("klay/prop/consensus/istanbul/in/packets", nil)
	propConsensusIstanbulInTrafficMeter  = metrics.NewRegisteredMeter("klay/prop/consensus/istanbul/in/traffic", nil)
	propConsensusIstanbulOutPacketsMeter = metrics.NewRegisteredMeter("klay/prop/consensus/istanbul/out/packets", nil)
	propConsensusIstanbulOutTrafficMeter = metrics.NewRegisteredMeter("klay/prop/consensus/istanbul/out/traffic", nil)
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
	case msg.Code == BlockHeadersMsg:
		packets, traffic = reqHeaderInPacketsMeter, reqHeaderInTrafficMeter
	case msg.Code == BlockBodiesMsg:
		packets, traffic = reqBodyInPacketsMeter, reqBodyInTrafficMeter

	case rw.version >= klay63 && msg.Code == NodeDataMsg:
		packets, traffic = reqStateInPacketsMeter, reqStateInTrafficMeter
	case rw.version >= klay63 && msg.Code == ReceiptsMsg:
		packets, traffic = reqReceiptInPacketsMeter, reqReceiptInTrafficMeter

	case msg.Code == NewBlockHashesMsg:
		packets, traffic = propHashInPacketsMeter, propHashInTrafficMeter
	case msg.Code == NewBlockMsg:
		packets, traffic = propBlockInPacketsMeter, propBlockInTrafficMeter
	case msg.Code == TxMsg:
		packets, traffic = propTxnInPacketsMeter, propTxnInTrafficMeter
	case msg.Code == backend.IstanbulMsg:
		packets, traffic = propConsensusIstanbulInPacketsMeter, propConsensusIstanbulInTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	return msg, err
}

func (rw *meteredMsgReadWriter) WriteMsg(msg p2p.Msg) error {
	// Account for the data traffic
	packets, traffic := miscOutPacketsMeter, miscOutTrafficMeter
	switch {
	case msg.Code == BlockHeadersMsg:
		packets, traffic = reqHeaderOutPacketsMeter, reqHeaderOutTrafficMeter
	case msg.Code == BlockBodiesMsg:
		packets, traffic = reqBodyOutPacketsMeter, reqBodyOutTrafficMeter

	case rw.version >= klay63 && msg.Code == NodeDataMsg:
		packets, traffic = reqStateOutPacketsMeter, reqStateOutTrafficMeter
	case rw.version >= klay63 && msg.Code == ReceiptsMsg:
		packets, traffic = reqReceiptOutPacketsMeter, reqReceiptOutTrafficMeter

	case msg.Code == NewBlockHashesMsg:
		packets, traffic = propHashOutPacketsMeter, propHashOutTrafficMeter
	case msg.Code == NewBlockMsg:
		packets, traffic = propBlockOutPacketsMeter, propBlockOutTrafficMeter
	case msg.Code == TxMsg:
		packets, traffic = propTxnOutPacketsMeter, propTxnOutTrafficMeter
	case msg.Code == backend.IstanbulMsg:
		packets, traffic = propConsensusIstanbulOutPacketsMeter, propConsensusIstanbulOutTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	// Send the packet to the p2p layer
	return rw.MsgReadWriter.WriteMsg(msg)
}
