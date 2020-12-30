// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
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
// This file is derived from p2p/peer.go (2018/06/04).
// Modified and improved for the klaytn development.

package p2p

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/mclock"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/rlp"
)

var logger = log.NewModuleLogger(log.NetworksP2P)

const (
	baseProtocolVersion    = 5
	baseProtocolLength     = uint64(16)
	baseProtocolMaxMsgSize = 2 * 1024

	snappyProtocolVersion = 5

	pingInterval = 15 * time.Second
)

const (
	// devp2p message codes
	handshakeMsg = 0x00
	discMsg      = 0x01
	pingMsg      = 0x02
	pongMsg      = 0x03
)

const (
	ConnDefault = iota
	ConnTxMsg
)

// protoHandshake is the RLP structure of the protocol handshake.
type protoHandshake struct {
	Version      uint64
	Name         string
	Caps         []Cap
	ListenPort   []uint64
	ID           discover.NodeID
	Multichannel bool

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// PeerEventType is the type of peer events emitted by a p2p.Server
type PeerEventType string

const (
	// PeerEventTypeAdd is the type of event emitted when a peer is added
	// to a p2p.Server
	PeerEventTypeAdd PeerEventType = "add"

	// PeerEventTypeDrop is the type of event emitted when a peer is
	// dropped from a p2p.Server
	PeerEventTypeDrop PeerEventType = "drop"

	// PeerEventTypeMsgSend is the type of event emitted when a
	// message is successfully sent to a peer
	PeerEventTypeMsgSend PeerEventType = "msgsend"

	// PeerEventTypeMsgRecv is the type of event emitted when a
	// message is received from a peer
	PeerEventTypeMsgRecv PeerEventType = "msgrecv"
)

// PeerEvent is an event emitted when peers are either added or dropped from
// a p2p.Server or when a message is sent or received on a peer connection
type PeerEvent struct {
	Type     PeerEventType   `json:"type"`
	Peer     discover.NodeID `json:"peer"`
	Error    string          `json:"error,omitempty"`
	Protocol string          `json:"protocol,omitempty"`
	MsgCode  *uint64         `json:"msg_code,omitempty"`
	MsgSize  *uint32         `json:"msg_size,omitempty"`
}

// Peer represents a connected remote node.
type Peer struct {
	rws     []*conn
	running map[string][]*protoRW
	logger  log.Logger
	created mclock.AbsTime

	wg       sync.WaitGroup
	protoErr chan error
	closed   chan struct{}
	disc     chan DiscReason

	// events receives message send / receive events if set
	events *event.Feed
}

// NewPeer returns a peer for testing purposes.
func NewPeer(id discover.NodeID, name string, caps []Cap) *Peer {
	pipe, _ := net.Pipe()
	conn := []*conn{{fd: pipe, transport: nil, id: id, caps: caps, name: name}}
	peer, _ := newPeer(conn, nil, defaultRWTimerConfig)
	close(peer.closed) // ensures Disconnect doesn't block
	return peer
}

// ID returns the node's public key.
func (p *Peer) ID() discover.NodeID {
	return p.rws[ConnDefault].id
}

// Name returns the node name that the remote node advertised.
func (p *Peer) Name() string {
	return p.rws[ConnDefault].name
}

// Caps returns the capabilities (supported subprotocols) of the remote peer.
func (p *Peer) Caps() []Cap {
	// TODO: maybe return copy
	return p.rws[ConnDefault].caps
}

// RemoteAddr returns the remote address of the network connection.
func (p *Peer) RemoteAddr() net.Addr {
	return p.rws[ConnDefault].fd.RemoteAddr()
}

// LocalAddr returns the local address of the network connection.
func (p *Peer) LocalAddr() net.Addr {
	return p.rws[ConnDefault].fd.LocalAddr()
}

// Disconnect terminates the peer connection with the given reason.
// It returns immediately and does not wait until the connection is closed.
func (p *Peer) Disconnect(reason DiscReason) {
	select {
	case p.disc <- reason:
	case <-p.closed:
	}
}

// String implements fmt.Stringer.
func (p *Peer) String() string {
	return fmt.Sprintf("Peer %x %v", p.rws[ConnDefault].id[:8], p.RemoteAddr())
}

// Inbound returns true if the peer is an inbound connection
func (p *Peer) Inbound() bool {
	return p.rws[ConnDefault].flags&inboundConn != 0
}

// GetNumberInboundAndOutbound returns the number of
// inbound and outbound connections connected to the peer.
func (p *Peer) GetNumberInboundAndOutbound() (int, int) {
	inbound, outbound := 0, 0
	for _, rw := range p.rws {
		if rw.flags&inboundConn != 0 {
			inbound++
		} else {
			outbound++
		}
	}
	return inbound, outbound
}

// newPeer should be called to create a peer.
func newPeer(conns []*conn, protocols []Protocol, tc RWTimerConfig) (*Peer, error) {
	if conns == nil || len(conns) < 1 || conns[ConnDefault] == nil {
		return nil, errors.New("conn is invalid")
	}
	msgReadWriters := make([]MsgReadWriter, len(conns))
	for i, c := range conns {
		msgReadWriters[i] = c
	}
	protomap := matchProtocols(protocols, conns[ConnDefault].caps, msgReadWriters, tc)
	p := &Peer{
		rws:      conns,
		running:  protomap,
		created:  mclock.Now(),
		disc:     make(chan DiscReason),
		protoErr: make(chan error, len(protomap)+len(conns)), // protocols + pingLoop
		closed:   make(chan struct{}),
		logger:   logger.NewWith("id", conns[ConnDefault].id, "conn", conns[ConnDefault].flags),
	}
	return p, nil
}

func (p *Peer) Log() log.Logger {
	return p.logger
}

func (p *Peer) run() (remoteRequested bool, err error) {
	var (
		writeStart = make(chan struct{}, 1)
		writeErr   = make(chan error, 1)
		readErr    = make(chan error, 1)
		reason     DiscReason // sent to the peer
	)
	if len(p.rws) != 1 {
		return false, errors.New("The size of rws should be 1")
	}
	p.wg.Add(2)
	go p.readLoop(ConnDefault, p.rws[ConnDefault], readErr)
	go p.pingLoop(p.rws[ConnDefault])

	// Start all protocol handlers.
	writeStart <- struct{}{}
	p.startProtocols(writeStart, writeErr)

	// Wait for an error or disconnect.
loop:
	for {
		select {
		case err = <-writeErr:
			// A write finished. Allow the next write to start if
			// there was no error.
			if err != nil {
				reason = DiscNetworkError
				break loop
			}
			writeStart <- struct{}{}
		case err = <-readErr:
			if r, ok := err.(DiscReason); ok {
				remoteRequested = true
				reason = r
			} else {
				reason = DiscNetworkError
			}
			break loop
		case err = <-p.protoErr:
			reason = discReasonForError(err)
			break loop
		case err = <-p.disc:
			reason = discReasonForError(err)
			break loop
		}
	}

	close(p.closed)
	for _, rw := range p.rws {
		rw.close(reason)
	}
	p.wg.Wait()
	logger.Debug(fmt.Sprintf("run stopped, peer: %v", p.ID()))
	return remoteRequested, err
}

// ErrorPeer is a peer error
type ErrorPeer struct {
	remoteRequested bool
	err             error
}

// runWithRWs Runs peer
func (p *Peer) runWithRWs() (remoteRequested bool, err error) {
	resultErr := make(chan ErrorPeer, len(p.rws))
	var errs ErrorPeer

	var (
		writeStarts = make([]chan struct{}, 0, len(p.rws))
		writeErr    = make([]chan error, 0, 1)
		readErr     = make([]chan error, 0, 1)
	)

	for range p.rws {
		writeStarts = append(writeStarts, make(chan struct{}, 1))
		writeErr = append(writeErr, make(chan error, 1))
		readErr = append(readErr, make(chan error, 1))
	}

	for i, rw := range p.rws {
		p.wg.Add(2)
		go p.readLoop(i, rw, readErr[i])
		go p.pingLoop(rw)
		writeStarts[i] <- struct{}{}
	}

	// Start all protocol handlers.
	p.startProtocolsWithRWs(writeStarts, writeErr)

	for i, rw := range p.rws {
		p.wg.Add(1)
		go p.handleError(rw, resultErr, writeErr[i], writeStarts[i], readErr[i])
	}

	select {
	case errs = <-resultErr:
		close(p.closed)
	}
	p.wg.Wait()
	logger.Debug(fmt.Sprintf("run stopped, peer: %v", p.ID()))
	return errs.remoteRequested, errs.err
}

//handleError handles read, write, and protocol errors on rw
func (p *Peer) handleError(rw *conn, errCh chan<- ErrorPeer, writeErr <-chan error, writeStart chan<- struct{}, readErr <-chan error) {
	defer p.wg.Done()
	var errRW ErrorPeer
	var reason DiscReason // sent to the peer

	// Wait for an error or disconnect.
loop:
	for {
		select {
		case errRW.err = <-writeErr:
			// A write finished. Allow the next write to start if
			// there was no error.
			if errRW.err != nil {
				reason = DiscNetworkError
				break loop
			}
			writeStart <- struct{}{}
		case errRW.err = <-readErr:
			if r, ok := errRW.err.(DiscReason); ok {
				errRW.remoteRequested = true
				reason = r
			} else {
				reason = DiscNetworkError
			}
			break loop
		case errRW.err = <-p.protoErr:
			reason = discReasonForError(errRW.err)
			break loop
		case errRW.err = <-p.disc:
			reason = discReasonForError(errRW.err)
			break loop
		case <-p.closed:
			reason = DiscQuitting
			break loop
		}
	}

	rw.close(reason)
	errCh <- errRW
}

func (p *Peer) pingLoop(rw *conn) {
	ping := time.NewTimer(pingInterval)
	defer p.wg.Done()
	defer ping.Stop()
	for {
		select {
		case <-ping.C:
			if err := SendItems(rw, pingMsg); err != nil {
				p.protoErr <- err
				logger.Debug(fmt.Sprintf("pingLoop stopped, peer: %v", p.ID()))
				return
			}
			ping.Reset(pingInterval)
		case <-p.closed:
			logger.Debug(fmt.Sprintf("pingLoop stopped, peer: %v", p.ID()))
			return
		}
	}
}

func (p *Peer) readLoop(connectionOrder int, rw *conn, errc chan<- error) {
	defer p.wg.Done()
	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			errc <- err
			logger.Debug(fmt.Sprintf("readLoop stopped, peer: %v", p.ID()))
			return
		}
		msg.ReceivedAt = time.Now()
		if err = p.handle(connectionOrder, rw, msg); err != nil {
			errc <- err
			logger.Debug(fmt.Sprintf("readLoop stopped, peer: %v", p.ID()))
			return
		}
	}
}

func (p *Peer) handle(connectionOrder int, rw *conn, msg Msg) error {
	switch {
	case msg.Code == pingMsg:
		msg.Discard()
		go SendItems(rw, pongMsg)
	case msg.Code == discMsg:
		var reason [1]DiscReason
		// This is the last message. We don't need to discard or
		// check errors because, the connection will be closed after it.
		rlp.Decode(msg.Payload, &reason)
		return reason[0]
	case msg.Code < baseProtocolLength:
		// ignore other base protocol messages
		return msg.Discard()
	default:
		// it's a subprotocol message
		proto, err := p.getProto(connectionOrder, msg.Code)
		if err != nil {
			return fmt.Errorf("msg code out of range: %v", msg.Code)
		}
		select {
		case proto.in <- msg:
			return nil
		case <-p.closed:
			return io.EOF
		}
	}
	return nil
}

func countMatchingProtocols(protocols []Protocol, caps []Cap) int {
	n := 0
	for _, cap := range caps {
		for _, proto := range protocols {
			if proto.Name == cap.Name && proto.Version == cap.Version {
				n++
			}
		}
	}
	return n
}

// matchProtocols creates structures for matching named subprotocols.
func matchProtocols(protocols []Protocol, caps []Cap, rws []MsgReadWriter, tc RWTimerConfig) map[string][]*protoRW {
	sort.Sort(capsByNameAndVersion(caps))
	offset := baseProtocolLength
	result := make(map[string][]*protoRW)

outer:
	for _, cap := range caps {
		for _, proto := range protocols {
			if proto.Name == cap.Name && proto.Version == cap.Version {
				// If an old protocol version matched, revert it
				if old := result[cap.Name]; old != nil {
					offset -= old[ConnDefault].Length
				}
				// Assign the new match
				protoRWs := make([]*protoRW, 0, len(rws))
				if rws == nil || len(rws) == 0 {
					protoRWs = []*protoRW{{Protocol: proto, offset: offset, in: make(chan Msg), w: nil, tc: tc}}
				} else {
					for _, rw := range rws {
						protoRWs = append(protoRWs, &protoRW{Protocol: proto, offset: offset, in: make(chan Msg), w: rw, tc: tc})
					}
				}
				result[cap.Name] = protoRWs
				offset += proto.Length

				continue outer
			}
		}
	}
	return result
}

func (p *Peer) startProtocols(writeStart <-chan struct{}, writeErr chan<- error) {
	p.wg.Add(len(p.running))
	for _, protos := range p.running {
		if len(protos) != 1 {
			p.logger.Error("The size of protos should be 1", "size", len(protos))
			p.protoErr <- errProtocolReturned
			return
		}
		proto := protos[ConnDefault]
		proto.closed = p.closed
		proto.wstart = writeStart
		proto.werr = writeErr
		proto.tc = defaultRWTimerConfig
		var rw MsgReadWriter = proto
		if p.events != nil {
			rw = newMsgEventer(rw, p.events, p.ID(), proto.Name)
		}
		p.logger.Trace(fmt.Sprintf("Starting protocol %s/%d", proto.Name, proto.Version))
		go func() {
			//p.wg.Add(1)
			defer p.wg.Done()
			err := proto.Run(p, rw)
			if err == nil {
				p.logger.Trace(fmt.Sprintf("Protocol %s/%d returned", proto.Name, proto.Version))
				err = errProtocolReturned
			} else if err != io.EOF {
				p.logger.Error(fmt.Sprintf("Protocol %s/%d failed", proto.Name, proto.Version), "err", err)
			}
			p.protoErr <- err
			p.logger.Debug(fmt.Sprintf("Protocol go routine stopped, peer: %v", p.ID()))
			//p.wg.Done()
		}()
	}
}

// startProtocolsWithRWs run the protocol using several RWs.
func (p *Peer) startProtocolsWithRWs(writeStarts []chan struct{}, writeErrs []chan error) {
	p.wg.Add(len(p.running))

	for _, protos := range p.running {
		rws := make([]MsgReadWriter, 0, len(protos))
		protos := protos
		for i, proto := range protos {
			proto.closed = p.closed
			if len(writeStarts) > i {
				proto.wstart = writeStarts[i]
			} else {
				writeErrs[i] <- errors.New("WriteStartsChannelSize")
			}
			proto.werr = writeErrs[i]

			var rw MsgReadWriter = proto
			if p.events != nil {
				rw = newMsgEventer(rw, p.events, p.ID(), proto.Name)
			}
			rws = append(rws, rw)
		}

		p.logger.Trace(fmt.Sprintf("Starting protocol %s/%d", protos[ConnDefault].Name, protos[ConnDefault].Version))
		go func() {
			//p.wg.Add(1)
			defer p.wg.Done()
			err := protos[ConnDefault].RunWithRWs(p, rws)
			if err == nil {
				p.logger.Trace(fmt.Sprintf("Protocol %s/%d returned", protos[ConnDefault].Name, protos[ConnDefault].Version))
				err = errProtocolReturned
			} else if err != io.EOF {
				p.logger.Error(fmt.Sprintf("Protocol %s/%d failed", protos[ConnDefault].Name, protos[ConnDefault].Version), "err", err)
			}
			p.protoErr <- err
			p.logger.Debug(fmt.Sprintf("Protocol go routine stopped, peer: %v", p.ID()))
			//p.wg.Done()
		}()
	}
}

// getProto finds the protocol responsible for handling
// the given message code.
func (p *Peer) getProto(connectionOrder int, code uint64) (*protoRW, error) {
	for _, proto := range p.running {
		if code >= proto[connectionOrder].offset && code < proto[connectionOrder].offset+proto[connectionOrder].Length {
			return proto[connectionOrder], nil
		}
	}
	return nil, newPeerError(errInvalidMsgCode, "%d", code)
}

type RWTimerConfig struct {
	Interval uint64
	WaitTime time.Duration
}

var defaultRWTimerConfig = RWTimerConfig{1000, 15 * time.Second}

type protoRW struct {
	Protocol
	in     chan Msg        // receices read messages
	closed <-chan struct{} // receives when peer is shutting down
	wstart <-chan struct{} // receives when write may start
	werr   chan<- error    // for write results
	offset uint64
	w      MsgWriter
	count  uint64 // count the number of WriteMsg calls
	tc     RWTimerConfig
}

func (rw *protoRW) WriteMsg(msg Msg) (err error) {
	if msg.Code >= rw.Length {
		return newPeerError(errInvalidMsgCode, "not handled, (code %x) (size %d)", msg.Code, msg.Size)
	}
	msg.Code += rw.offset
	rwCount := atomic.AddUint64(&rw.count, 1)
	if rwCount%rw.tc.Interval == 0 {
		timer := time.NewTimer(rw.tc.WaitTime)
		defer timer.Stop()
		select {
		case <-rw.wstart:
			err = rw.w.WriteMsg(msg)
			// Report write status back to Peer.run. It will initiate
			// shutdown if the error is non-nil and unblock the next write
			// otherwise. The calling protocol code should exit for errors
			// as well but we don't want to rely on that.
		case <-rw.closed:
			err = fmt.Errorf("shutting down")
			return err
		case <-timer.C:
			writeMsgTimeOutCounter.Inc(1)
			err = fmt.Errorf("failed to write message for %v", rw.tc.WaitTime)
		}
	} else {
		select {
		case <-rw.wstart:
			err = rw.w.WriteMsg(msg)
		case <-rw.closed:
			err = fmt.Errorf("shutting down")
			return err
		}
	}
	select {
	case rw.werr <- err:
	default:
	}
	return err
}

func (rw *protoRW) ReadMsg() (Msg, error) {
	select {
	case msg := <-rw.in:
		msg.Code -= rw.offset
		return msg, nil
	case <-rw.closed:
		return Msg{}, io.EOF
	}
}

// NetworkInfo represents the connection information with the peer.
type NetworkInfo struct {
	LocalAddress  string `json:"localAddress"`  // Local endpoint of the TCP data connection
	RemoteAddress string `json:"remoteAddress"` // Remote endpoint of the TCP data connection
	Inbound       bool   `json:"inbound"`
	Trusted       bool   `json:"trusted"`
	Static        bool   `json:"static"`
	NodeType      string `json:"nodeType"`
}

// PeerInfo represents a short summary of the information known about a connected
// peer. Sub-protocol independent fields are contained and initialized here, with
// protocol specifics delegated to all connected sub-protocols.
type PeerInfo struct {
	ID        string                 `json:"id"`        // Unique node identifier (also the encryption key)
	Name      string                 `json:"name"`      // Name of the node, including client type, version, OS, custom data
	Caps      []string               `json:"caps"`      // Sum-protocols advertised by this particular peer
	Networks  []NetworkInfo          `json:"networks"`  // Networks is all the NetworkInfo associated with the peer
	Protocols map[string]interface{} `json:"protocols"` // Sub-protocol specific metadata fields
}

// Info gathers and returns a collection of metadata known about a peer.
func (p *Peer) Info() *PeerInfo {
	// Gather the protocol capabilities
	var caps []string
	for _, cap := range p.Caps() {
		caps = append(caps, cap.String())
	}
	// Assemble the generic peer metadata
	info := &PeerInfo{
		ID:        p.ID().String(),
		Name:      p.Name(),
		Caps:      caps,
		Protocols: make(map[string]interface{}),
	}

	for _, rw := range p.rws {
		var network NetworkInfo
		network.LocalAddress = rw.fd.LocalAddr().String()
		network.RemoteAddress = rw.fd.RemoteAddr().String()
		network.Inbound = rw.is(inboundConn)
		network.Trusted = rw.is(trustedConn)
		network.Static = rw.is(staticDialedConn)
		switch rw.conntype {
		case common.CONSENSUSNODE:
			network.NodeType = "cn"
		case common.ENDPOINTNODE:
			network.NodeType = "en"
		case common.PROXYNODE:
			network.NodeType = "pn"
		case common.BOOTNODE:
			network.NodeType = "bn"
		default:
			network.NodeType = "unknown"
		}
		info.Networks = append(info.Networks, network)
	}

	// Gather all the running protocol infos
	for _, proto := range p.running {
		protoInfo := interface{}("unknown")
		if query := proto[ConnDefault].Protocol.PeerInfo; query != nil {
			if metadata := query(p.ID()); metadata != nil {
				protoInfo = metadata
			} else {
				protoInfo = "handshake"
			}
		}
		info.Protocols[proto[ConnDefault].Name] = protoInfo
	}
	return info
}

func (p *Peer) ConnType() common.ConnType {
	return p.rws[ConnDefault].conntype
}

type PeerTypeValidator interface {
	// ValidatePeerType returns nil if successful. Otherwise, it returns an error object.
	ValidatePeerType(addr common.Address) error
}
