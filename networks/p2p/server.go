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
// This file is derived from p2p/server.go (2018/06/04).
// Modified and improved for the klaytn development.

package p2p

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/klaytn/klaytn/common/mclock"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/nat"
	"github.com/klaytn/klaytn/networks/p2p/netutil"

	"github.com/klaytn/klaytn/common"
)

const (
	defaultDialTimeout = 15 * time.Second

	// Connectivity defaults.
	maxActiveDialTasks     = 16
	defaultMaxPendingPeers = 50
	defaultDialRatio       = 3

	// Maximum time allowed for reading a complete message.
	// This is effectively the amount of time a connection can be idle.
	frameReadTimeout = 30 * time.Second

	// Maximum amount of time allowed for writing a complete message.
	frameWriteTimeout = 20 * time.Second
)

var errServerStopped = errors.New("server stopped")

// Config holds Server options.
type Config struct {
	// This field must be set to a valid secp256k1 private key.
	PrivateKey *ecdsa.PrivateKey `toml:"-"`

	// MaxPhysicalConnections is the maximum number of physical connections.
	// A peer uses one connection if single channel peer and uses two connections if
	// multi channel peer. It must be greater than zero.
	MaxPhysicalConnections int

	// ConnectionType is a type of connection like Consensus or Normal
	// described at common.ConnType
	// When the connection is established, each peer exchange each connection type
	ConnectionType common.ConnType

	// MaxPendingPeers is the maximum number of peers that can be pending in the
	// handshake phase, counted separately for inbound and outbound connections.
	// Zero defaults to preset values.
	MaxPendingPeers int `toml:",omitempty"`

	// DialRatio controls the ratio of inbound to dialed connections.
	// Example: a DialRatio of 2 allows 1/2 of connections to be dialed.
	// Setting DialRatio to zero defaults it to 3.
	DialRatio int `toml:",omitempty"`

	// NoDiscovery can be used to disable the peer discovery mechanism.
	// Disabling is useful for protocol debugging (manual topology).
	NoDiscovery bool

	// Name sets the node name of this server.
	// Use common.MakeName to create a name that follows existing conventions.
	Name string `toml:"-"`

	// BootstrapNodes are used to establish connectivity
	// with the rest of the network.
	BootstrapNodes []*discover.Node

	//// BootstrapNodesV5 are used to establish connectivity
	//// with the rest of the network using the V5 discovery
	//// protocol.
	//BootstrapNodesV5 []*discv5.Node `toml:",omitempty"`

	// Static nodes are used as pre-configured connections which are always
	// maintained and re-connected on disconnects.
	StaticNodes []*discover.Node

	// Trusted nodes are used as pre-configured connections which are always
	// allowed to connect, even above the peer limit.
	TrustedNodes []*discover.Node

	// Connectivity can be restricted to certain IP networks.
	// If this option is set to a non-nil value, only hosts which match one of the
	// IP networks contained in the list are considered.
	NetRestrict *netutil.Netlist `toml:",omitempty"`

	// NodeDatabase is the path to the database containing the previously seen
	// live nodes in the network.
	NodeDatabase string `toml:",omitempty"`

	// Protocols should contain the protocols supported
	// by the server. Matching protocols are launched for
	// each peer.
	Protocols []Protocol `toml:"-"`

	// If ListenAddr is set to a non-nil address, the server
	// will listen for incoming connections.
	//
	// If the port is zero, the operating system will pick a port. The
	// ListenAddr field will be updated with the actual address when
	// the server is started.
	ListenAddr string

	// NoListen can be used to disable the listening for incoming connections.
	NoListen bool

	// SubListenAddr is the list of the secondary listen address used for peer-to-peer connections.
	SubListenAddr []string

	// If EnableMultiChannelServer is true, multichannel can communicate with other nodes
	EnableMultiChannelServer bool

	// If set to a non-nil value, the given NAT port mapper
	// is used to make the listening port available to the
	// Internet.
	NAT nat.Interface `toml:",omitempty"`

	// If Dialer is set to a non-nil value, the given Dialer
	// is used to dial outbound peer connections.
	Dialer NodeDialer `toml:"-"`

	// If NoDial is true, the server will not dial any peers.
	NoDial bool `toml:",omitempty"`

	// If EnableMsgEvents is set then the server will emit PeerEvents
	// whenever a message is sent to or received from a peer
	EnableMsgEvents bool

	// Logger is a custom logger to use with the p2p.Server.
	Logger log.Logger `toml:",omitempty"`

	// RWTimerConfig is a configuration for interval based timer for rw.
	// It checks if a rw successfully writes its task in given time.
	RWTimerConfig RWTimerConfig

	// NetworkID to use for selecting peers to connect to
	NetworkID uint64
}

// NewServer returns a new Server interface.
func NewServer(config Config) Server {
	bServer := &BaseServer{
		Config: config,
	}

	if config.EnableMultiChannelServer {
		listeners := make([]net.Listener, 0, len(config.SubListenAddr)+1)
		listenAddrs := make([]string, 0, len(config.SubListenAddr)+1)
		listenAddrs = append(listenAddrs, config.ListenAddr)
		listenAddrs = append(listenAddrs, config.SubListenAddr...)
		return &MultiChannelServer{
			BaseServer:     bServer,
			listeners:      listeners,
			ListenAddrs:    listenAddrs,
			CandidateConns: make(map[discover.NodeID][]*conn),
		}
	} else {
		return &SingleChannelServer{
			BaseServer: bServer,
		}
	}
}

// Server manages all peer connections.
type Server interface {
	// GetProtocols returns a slice of protocols.
	GetProtocols() []Protocol

	// AddProtocols adds protocols to the server.
	AddProtocols(p []Protocol)

	// SetupConn runs the handshakes and attempts to add the connection
	// as a peer. It returns when the connection has been added as a peer
	// or the handshakes have failed.
	SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error

	// AddLastLookup adds lastLookup to duration.
	AddLastLookup() time.Time

	// SetLastLookupToNow sets LastLookup to the current time.
	SetLastLookupToNow()

	// CheckNilNetworkTable returns whether network table is nil.
	CheckNilNetworkTable() bool

	// GetNodes returns up to max alive nodes which a NodeType is nType
	GetNodes(nType discover.NodeType, max int) []*discover.Node

	// Lookup performs a network search for nodes close
	// to the given target. It approaches the target by querying
	// nodes that are closer to it on each iteration.
	// The given target does not need to be an actual node
	// identifier.
	Lookup(target discover.NodeID, nType discover.NodeType) []*discover.Node

	// Resolve searches for a specific node with the given ID and NodeType.
	// It returns nil if the node could not be found.
	Resolve(target discover.NodeID, nType discover.NodeType) *discover.Node

	// Start starts running the server.
	// Servers can not be re-used after stopping.
	Start() (err error)

	// Stop terminates the server and all active peer connections.
	// It blocks until all active connections are closed.
	Stop()

	// AddPeer connects to the given node and maintains the connection until the
	// server is shut down. If the connection fails for any reason, the server will
	// attempt to reconnect the peer.
	AddPeer(node *discover.Node)

	// RemovePeer disconnects from the given node.
	RemovePeer(node *discover.Node)

	// SubscribePeers subscribes the given channel to peer events.
	SubscribeEvents(ch chan *PeerEvent) event.Subscription

	// PeersInfo returns an array of metadata objects describing connected peers.
	PeersInfo() []*PeerInfo

	// NodeInfo gathers and returns a collection of metadata known about the host.
	NodeInfo() *NodeInfo

	// Name returns name of server.
	Name() string

	// PeerCount returns the number of connected peers.
	PeerCount() int

	// PeerCountByType returns the number of connected specific tyeps of peers.
	PeerCountByType() map[string]uint

	// MaxPhysicalConnections returns maximum count of peers.
	MaxPeers() int

	// Disconnect tries to disconnect peer.
	Disconnect(destID discover.NodeID)

	// GetListenAddress returns the listen address list of the server.
	GetListenAddress() []string

	// Peers returns all connected peers.
	Peers() []*Peer

	// NodeDialer is used to connect to nodes in the network, typically by using
	// an underlying net.Dialer but also using net.Pipe in tests.
	NodeDialer
}

// MultiChannelServer is a server that uses a multi channel.
type MultiChannelServer struct {
	*BaseServer
	listeners      []net.Listener
	ListenAddrs    []string
	CandidateConns map[discover.NodeID][]*conn
}

// Start starts running the MultiChannelServer.
// MultiChannelServer can not be re-used after stopping.
func (srv *MultiChannelServer) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	srv.logger = srv.Config.Logger
	if srv.logger == nil {
		srv.logger = logger.NewWith()
	}
	srv.logger.Info("Starting P2P networking")

	// static fields
	if srv.PrivateKey == nil {
		return fmt.Errorf("Server.PrivateKey must be set to a non-nil key")
	}

	if !srv.ConnectionType.Valid() {
		return fmt.Errorf("Invalid connection type speficied")
	}

	if srv.newTransport == nil {
		srv.newTransport = newRLPX
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.addpeer = make(chan *conn)
	srv.delpeer = make(chan peerDrop)
	srv.posthandshake = make(chan *conn)
	srv.addstatic = make(chan *discover.Node)
	srv.removestatic = make(chan *discover.Node)
	srv.peerOp = make(chan peerOpFunc)
	srv.peerOpDone = make(chan struct{})
	srv.discpeer = make(chan discover.NodeID)

	var (
		conn      *net.UDPConn
		realaddr  *net.UDPAddr
		unhandled chan discover.ReadPacket
	)

	if !srv.NoDiscovery {
		addr, err := net.ResolveUDPAddr("udp", srv.ListenAddrs[ConnDefault])
		if err != nil {
			return err
		}
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}
		realaddr = conn.LocalAddr().(*net.UDPAddr)
		if srv.NAT != nil {
			if !realaddr.IP.IsLoopback() {
				go nat.Map(srv.NAT, srv.quit, "udp", realaddr.Port, realaddr.Port, "klaytn discovery")
			}
			// TODO: react to external IP changes over time.
			if ext, err := srv.NAT.ExternalIP(); err == nil {
				realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
			}
		}
	}

	// node table
	if !srv.NoDiscovery {
		cfg := discover.Config{
			PrivateKey:   srv.PrivateKey,
			AnnounceAddr: realaddr,
			NodeDBPath:   srv.NodeDatabase,
			NetRestrict:  srv.NetRestrict,
			Bootnodes:    srv.BootstrapNodes,
			Unhandled:    unhandled,
			Conn:         conn,
			Addr:         realaddr,
			Id:           discover.PubkeyID(&srv.PrivateKey.PublicKey),
			NodeType:     ConvertNodeType(srv.ConnectionType),
			NetworkID:    srv.NetworkID,
		}

		ntab, err := discover.ListenUDP(&cfg)
		if err != nil {
			return err
		}
		srv.ntab = ntab
	}

	dialer := newDialState(srv.StaticNodes, srv.BootstrapNodes, srv.ntab, srv.maxDialedConns(), srv.NetRestrict, srv.PrivateKey, srv.getTypeStatics())

	// handshake
	srv.ourHandshake = &protoHandshake{Version: baseProtocolVersion, Name: srv.Name(), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey), Multichannel: true}
	for _, p := range srv.Protocols {
		srv.ourHandshake.Caps = append(srv.ourHandshake.Caps, p.cap())
	}
	for _, l := range srv.ListenAddrs {
		s := strings.Split(l, ":")
		if len(s) == 2 {
			if port, err := strconv.Atoi(s[1]); err == nil {
				srv.ourHandshake.ListenPort = append(srv.ourHandshake.ListenPort, uint64(port))
			}
		}
	}

	// listen/dial
	if srv.NoDial && srv.NoListen {
		srv.logger.Error("P2P server will be useless, neither dialing nor listening")
	}
	if !srv.NoListen {
		if srv.ListenAddrs != nil && len(srv.ListenAddrs) != 0 && srv.ListenAddrs[ConnDefault] != "" {
			if err := srv.startListening(); err != nil {
				return err
			}
		} else {
			srv.logger.Error("P2P server might be useless, listening address is missing")
		}
	}

	srv.loopWG.Add(1)
	go srv.run(dialer)
	srv.running = true
	srv.logger.Info("Started P2P server", "id", discover.PubkeyID(&srv.PrivateKey.PublicKey), "multichannel", true)
	return nil
}

// startListening starts listening on the specified port on the server.
func (srv *MultiChannelServer) startListening() error {
	// Launch the TCP listener.
	for i, listenAddr := range srv.ListenAddrs {
		listener, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return err
		}
		laddr := listener.Addr().(*net.TCPAddr)
		srv.ListenAddrs[i] = laddr.String()
		srv.listeners = append(srv.listeners, listener)
		srv.loopWG.Add(1)
		go srv.listenLoop(listener)
		// Map the TCP listening port if NAT is configured.
		if !laddr.IP.IsLoopback() && srv.NAT != nil {
			srv.loopWG.Add(1)
			go func() {
				nat.Map(srv.NAT, srv.quit, "tcp", laddr.Port, laddr.Port, "klaytn p2p")
				srv.loopWG.Done()
			}()
		}
	}
	return nil
}

// listenLoop waits for an external connection and connects it.
func (srv *MultiChannelServer) listenLoop(listener net.Listener) {
	defer srv.loopWG.Done()
	srv.logger.Info("RLPx listener up", "self", srv.makeSelf(listener, srv.ntab))

	tokens := defaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	for {
		// Wait for a handshake slot before accepting.
		<-slots

		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = listener.Accept()
			if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
				srv.logger.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				srv.logger.Debug("Read error", "err", err)
				return
			}
			break
		}

		// Reject connections that do not match NetRestrict.
		if srv.NetRestrict != nil {
			if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok && !srv.NetRestrict.Contains(tcp.IP) {
				srv.logger.Debug("Rejected conn (not whitelisted in NetRestrict)", "addr", fd.RemoteAddr())
				fd.Close()
				slots <- struct{}{}
				continue
			}
		}

		fd = newMeteredConn(fd, true)
		srv.logger.Trace("Accepted connection", "addr", fd.RemoteAddr())
		go func() {
			srv.SetupConn(fd, inboundConn, nil)
			slots <- struct{}{}
		}()
	}
}

// SetupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *MultiChannelServer) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}

	c := &conn{fd: fd, transport: srv.newTransport(fd), flags: flags, conntype: common.ConnTypeUndefined, cont: make(chan error), portOrder: PortOrderUndefined}
	if dialDest != nil {
		c.portOrder = PortOrder(dialDest.PortOrder)
	} else {
		for i, addr := range srv.ListenAddrs {
			s1 := strings.Split(addr, ":")                    // string format example, [::]:30303 or 123.123.123.123:30303
			s2 := strings.Split(fd.LocalAddr().String(), ":") // string format example, 123.123.123.123:30303
			if s1[len(s1)-1] == s2[len(s2)-1] {
				c.portOrder = PortOrder(i)
				break
			}
		}
	}

	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.close(err)
		srv.logger.Trace("close connection", "id", c.id, "err", err)
	}
	return err
}

// setupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *MultiChannelServer) setupConn(c *conn, flags connFlag, dialDest *discover.Node) error {
	// Prevent leftover pending conns from entering the handshake.
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}

	var err error
	// Run the connection type handshake
	if c.conntype, err = c.doConnTypeHandshake(srv.ConnectionType); err != nil {
		srv.logger.Warn("Failed doConnTypeHandshake", "addr", c.fd.RemoteAddr(), "conn", c.flags,
			"conntype", c.conntype, "err", err)
		return err
	}
	srv.logger.Trace("Connection Type Trace", "addr", c.fd.RemoteAddr(), "conn", c.flags, "ConnType", c.conntype.String())

	// Run the encryption handshake.
	if c.id, err = c.doEncHandshake(srv.PrivateKey, dialDest); err != nil {
		srv.logger.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return err
	}

	clog := srv.logger.NewWith("id", c.id, "addr", c.fd.RemoteAddr(), "conn", c.flags)
	// For dialed connections, check that the remote public key matches.
	if dialDest != nil && c.id != dialDest.ID {
		clog.Trace("Dialed identity mismatch", "want", c, dialDest.ID)
		return DiscUnexpectedIdentity
	}
	err = srv.checkpoint(c, srv.posthandshake)
	if err != nil {
		clog.Trace("Rejected peer before protocol handshake", "err", err)
		return err
	}
	// Run the protocol handshake
	phs, err := c.doProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed protobuf handshake", "err", err)
		return err
	}
	if phs.ID != c.id {
		clog.Trace("Wrong devp2p handshake identity", "err", phs.ID)
		return DiscUnexpectedIdentity
	}
	c.caps, c.name, c.multiChannel = phs.Caps, phs.Name, phs.Multichannel

	if c.multiChannel && dialDest != nil && (dialDest.TCPs == nil || len(dialDest.TCPs) < 2) && len(dialDest.TCPs) < len(phs.ListenPort) {
		logger.Debug("[Dial] update and retry the dial candidate as a multichannel",
			"id", dialDest.ID, "addr", dialDest.IP, "previous", dialDest.TCPs, "new", phs.ListenPort)

		dialDest.TCPs = make([]uint16, 0, len(phs.ListenPort))
		for _, listenPort := range phs.ListenPort {
			dialDest.TCPs = append(dialDest.TCPs, uint16(listenPort))
		}
		return errUpdateDial
	}

	err = srv.checkpoint(c, srv.addpeer)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

// run is the main loop that the server runs.
func (srv *MultiChannelServer) run(dialstate dialer) {
	logger.Debug("[p2p.Server] start MultiChannel p2p server")
	defer srv.loopWG.Done()
	var (
		peers         = make(map[discover.NodeID]*Peer)
		inboundCount  = 0
		outboundCount = 0
		trusted       = make(map[discover.NodeID]bool, len(srv.TrustedNodes))
		taskdone      = make(chan task, maxActiveDialTasks)
		runningTasks  []task
		queuedTasks   []task // tasks that can't run yet
	)
	// Put trusted nodes into a map to speed up checks.
	// Trusted peers are loaded on startup and cannot be
	// modified while the server is running.
	for _, n := range srv.TrustedNodes {
		trusted[n.ID] = true
	}

	// removes t from runningTasks
	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	// starts until max number of active tasks is satisfied
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < maxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			srv.logger.Trace("New dial task", "task", t)
			go func() { t.Do(srv); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		// Start from queue first.
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		// Query dialer for new tasks and start as many as possible now.
		if len(runningTasks) < maxActiveDialTasks {
			nt := dialstate.newTasks(len(runningTasks)+len(queuedTasks), peers, time.Now())
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}

running:
	for {
		scheduleTasks()

		select {
		case <-srv.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case n := <-srv.addstatic:
			// This channel is used by AddPeer to add to the
			// ephemeral static peer list. Add it to the dialer,
			// it will keep the node connected.
			srv.logger.Debug("Adding static node", "node", n)
			dialstate.addStatic(n)
		case n := <-srv.removestatic:
			// This channel is used by RemovePeer to send a
			// disconnect request to a peer and begin the
			// stop keeping the node connected
			srv.logger.Debug("Removing static node", "node", n)
			dialstate.removeStatic(n)
			if p, ok := peers[n.ID]; ok {
				p.Disconnect(DiscRequested)
			}
		case op := <-srv.peerOp:
			// This channel is used by Peers and PeerCount.
			op(peers)
			srv.peerOpDone <- struct{}{}
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.logger.Trace("Dial task done", "task", t)
			dialstate.taskDone(t, time.Now())
			delTask(t)
		case c := <-srv.posthandshake:
			// A connection has passed the encryption handshake so
			// the remote identity is known (but hasn't been verified yet).
			if trusted[c.id] {
				// Ensure that the trusted flag is set before checking against MaxPhysicalConnections.
				c.flags |= trustedConn
			}
			// TODO: track in-progress inbound node IDs (pre-Peer) to avoid dialing them.
			select {
			case c.cont <- srv.encHandshakeChecks(peers, inboundCount, c):
			case <-srv.quit:
				break running
			}
		case c := <-srv.addpeer:
			var p *Peer
			var e error
			// At this point the connection is past the protocol handshake.
			// Its capabilities are known and the remote identity is verified.
			err := srv.protoHandshakeChecks(peers, inboundCount, c)
			if err == nil {
				if c.multiChannel {
					connSet := srv.CandidateConns[c.id]
					if connSet == nil {
						connSet = make([]*conn, len(srv.ListenAddrs))
						srv.CandidateConns[c.id] = connSet
					}

					if int(c.portOrder) < len(connSet) {
						connSet[c.portOrder] = c
					}

					count := len(connSet)
					for _, conn := range connSet {
						if conn != nil {
							count--
						}
					}

					if count == 0 {
						p, e = newPeer(connSet, srv.Protocols, srv.Config.RWTimerConfig)
						srv.CandidateConns[c.id] = nil
					}
				} else {
					// The handshakes are done and it passed all checks.
					p, e = newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
				}

				if e != nil {
					srv.logger.Error("Fail make a new peer", "err", e)
				} else if p != nil {
					// If message events are enabled, pass the peerFeed
					// to the peer
					if srv.EnableMsgEvents {
						p.events = &srv.peerFeed
					}
					name := truncateName(c.name)
					srv.logger.Debug("Adding p2p peer", "name", name, "addr", c.fd.RemoteAddr(), "peers", len(peers)+1)
					go srv.runPeer(p)
					peers[c.id] = p

					peerCountGauge.Update(int64(len(peers)))
					inboundCount, outboundCount = increasesConnectionMetric(inboundCount, outboundCount, p)
				}
			}
			// The dialer logic relies on the assumption that
			// dial tasks complete after the peer has been added or
			// discarded. Unblock the task last.
			select {
			case c.cont <- err:
			case <-srv.quit:
				break running
			}
		case pd := <-srv.delpeer:
			// A peer disconnected.
			d := common.PrettyDuration(mclock.Now() - pd.created)
			pd.logger.Debug("Removing p2p peer", "duration", d, "peers", len(peers)-1, "req", pd.requested, "err", pd.err)
			delete(peers, pd.ID())

			peerCountGauge.Update(int64(len(peers)))
			inboundCount, outboundCount = decreasesConnectionMetric(inboundCount, outboundCount, pd.Peer)
		case nid := <-srv.discpeer:
			if p, ok := peers[nid]; ok {
				p.Disconnect(DiscRequested)
				p.logger.Debug(fmt.Sprintf("disconnect peer"))
			}
		}
	}

	srv.logger.Trace("P2P networking is spinning down")

	// Terminate discovery. If there is a running lookup it will terminate soon.
	if srv.ntab != nil {
		srv.ntab.Close()
	}
	//if srv.DiscV5 != nil {
	//	srv.DiscV5.Close()
	//}
	// Disconnect all peers.
	for _, p := range peers {
		p.Disconnect(DiscQuitting)
	}
	// Wait for peers to shut down. Pending connections and tasks are
	// not handled here and will terminate soon-ish because srv.quit
	// is closed.
	for len(peers) > 0 {
		p := <-srv.delpeer
		p.logger.Trace("<-delpeer (spindown)", "remainingTasks", len(runningTasks))
		delete(peers, p.ID())
	}
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections are closed.
func (srv *MultiChannelServer) Stop() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}
	srv.running = false
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}
	for _, listener := range srv.listeners {
		listener.Close()
	}
	close(srv.quit)
	srv.loopWG.Wait()
}

// GetListenAddress returns the listen addresses of the server.
func (srv *MultiChannelServer) GetListenAddress() []string {
	return srv.ListenAddrs
}

// decreasesConnectionMetric decreases the metric of the number of peer connections.
func decreasesConnectionMetric(inboundCount int, outboundCount int, p *Peer) (int, int) {
	pInbound, pOutbound := p.GetNumberInboundAndOutbound()
	inboundCount -= pInbound
	outboundCount -= pOutbound

	updatesConnectionMetric(inboundCount, outboundCount)
	return inboundCount, outboundCount
}

// increasesConnectionMetric increases the metric of the number of peer connections.
func increasesConnectionMetric(inboundCount int, outboundCount int, p *Peer) (int, int) {
	pInbound, pOutbound := p.GetNumberInboundAndOutbound()
	inboundCount += pInbound
	outboundCount += pOutbound

	updatesConnectionMetric(inboundCount, outboundCount)
	return inboundCount, outboundCount
}

// updatesConnectionMetric updates the metric of the number of peer connections.
func updatesConnectionMetric(inboundCount int, outboundCount int) {
	connectionCountGauge.Update(int64(outboundCount + inboundCount))
	connectionInCountGauge.Update(int64(inboundCount))
	connectionOutCountGauge.Update(int64(outboundCount))
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *MultiChannelServer) runPeer(p *Peer) {
	if srv.newPeerHook != nil {
		srv.newPeerHook(p)
	}

	// broadcast peer add
	srv.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.ID(),
	})

	// run the protocol
	remoteRequested, err := p.runWithRWs()

	// broadcast peer drop
	srv.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	srv.delpeer <- peerDrop{p, err, remoteRequested}
}

// SingleChannelServer is a server that uses a single channel.
type SingleChannelServer struct {
	*BaseServer
}

// AddLastLookup adds lastLookup to duration.
func (srv *BaseServer) AddLastLookup() time.Time {
	srv.lastLookupMu.Lock()
	defer srv.lastLookupMu.Unlock()
	return srv.lastLookup.Add(lookupInterval)
}

// SetLastLookupToNow sets LastLookup to the current time.
func (srv *BaseServer) SetLastLookupToNow() {
	srv.lastLookupMu.Lock()
	defer srv.lastLookupMu.Unlock()
	srv.lastLookup = time.Now()
}

// Dial creates a TCP connection to the node.
func (srv *BaseServer) Dial(dest *discover.Node) (net.Conn, error) {
	return srv.Dialer.Dial(dest)
}

// Dial creates a TCP connection to the node.
func (srv *BaseServer) DialMulti(dest *discover.Node) ([]net.Conn, error) {
	return srv.Dialer.DialMulti(dest)
}

// BaseServer is a common data structure used by implementation of Server.
type BaseServer struct {
	// Config fields may not be modified while the server is running.
	Config

	// Hooks for testing. These are useful because we can inhibit
	// the whole protocol stack.
	newTransport func(net.Conn) transport
	newPeerHook  func(*Peer)

	lock    sync.Mutex // protects running
	running bool

	ntab         discover.Discovery
	listener     net.Listener
	ourHandshake *protoHandshake
	lastLookup   time.Time
	lastLookupMu sync.Mutex
	//DiscV5       *discv5.Network

	// These are for Peers, PeerCount (and nothing else).
	peerOp     chan peerOpFunc
	peerOpDone chan struct{}

	quit          chan struct{}
	addstatic     chan *discover.Node
	removestatic  chan *discover.Node
	posthandshake chan *conn
	addpeer       chan *conn
	delpeer       chan peerDrop
	discpeer      chan discover.NodeID
	loopWG        sync.WaitGroup // loop, listenLoop
	peerFeed      event.Feed
	logger        log.Logger
}

type peerOpFunc func(map[discover.NodeID]*Peer)

type peerDrop struct {
	*Peer
	err       error
	requested bool // true if signaled by the peer
}

type connFlag int

const (
	dynDialedConn connFlag = 1 << iota
	staticDialedConn
	inboundConn
	trustedConn
)

type PortOrder int

const (
	PortOrderUndefined PortOrder = -1
)

// conn wraps a network connection with information gathered
// during the two handshakes.
type conn struct {
	fd net.Conn
	transport
	flags        connFlag
	conntype     common.ConnType // valid after the encryption handshake at the inbound connection case
	cont         chan error      // The run loop uses cont to signal errors to SetupConn.
	id           discover.NodeID // valid after the encryption handshake
	caps         []Cap           // valid after the protocol handshake
	name         string          // valid after the protocol handshake
	portOrder    PortOrder       // portOrder is the order of the ports that should be connected in multi-channel.
	multiChannel bool            // multiChannel is whether the peer is using multi-channel.
}

type transport interface {
	doConnTypeHandshake(myConnType common.ConnType) (common.ConnType, error)
	// The two handshakes.
	doEncHandshake(prv *ecdsa.PrivateKey, dialDest *discover.Node) (discover.NodeID, error)
	doProtoHandshake(our *protoHandshake) (*protoHandshake, error)
	// The MsgReadWriter can only be used after the encryption
	// handshake has completed. The code uses conn.id to track this
	// by setting it to a non-nil value after the encryption handshake.
	MsgReadWriter
	// transports must provide Close because we use MsgPipe in some of
	// the tests. Closing the actual network connection doesn't do
	// anything in those tests because NsgPipe doesn't use it.
	close(err error)
}

func (c *conn) String() string {
	s := c.flags.String()
	s += " " + c.conntype.String()
	if (c.id != discover.NodeID{}) {
		s += " " + c.id.String()
	}
	s += " " + c.fd.RemoteAddr().String()
	return s
}

func (c *conn) Inbound() bool {
	return c.flags&inboundConn != 0
}

func (f connFlag) String() string {
	s := ""
	if f&trustedConn != 0 {
		s += "-trusted"
	}
	if f&dynDialedConn != 0 {
		s += "-dyndial"
	}
	if f&staticDialedConn != 0 {
		s += "-staticdial"
	}
	if f&inboundConn != 0 {
		s += "-inbound"
	}
	if s != "" {
		s = s[1:]
	}
	return s
}

func (c *conn) is(f connFlag) bool {
	return c.flags&f != 0
}

// GetProtocols returns a slice of protocols.
func (srv *BaseServer) GetProtocols() []Protocol {
	return srv.Protocols
}

// AddProtocols adds protocols to the server.
func (srv *BaseServer) AddProtocols(p []Protocol) {
	srv.Protocols = append(srv.Protocols, p...)
}

// Peers returns all connected peers.
func (srv *BaseServer) Peers() []*Peer {
	var ps []*Peer
	select {
	// Note: We'd love to put this function into a variable but
	// that seems to cause a weird compiler error in some
	// environments.
	case srv.peerOp <- func(peers map[discover.NodeID]*Peer) {
		for _, p := range peers {
			ps = append(ps, p)
		}
	}:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return ps
}

// PeerCount returns the number of connected peers.
func (srv *BaseServer) PeerCount() int {
	var count int
	select {
	case srv.peerOp <- func(ps map[discover.NodeID]*Peer) { count = len(ps) }:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return count
}

func (srv *BaseServer) PeerCountByType() map[string]uint {
	pc := make(map[string]uint)
	pc["total"] = 0
	select {
	case srv.peerOp <- func(ps map[discover.NodeID]*Peer) {
		for _, peer := range ps {
			key := ConvertConnTypeToString(peer.ConnType())
			pc[key]++
			pc["total"]++
		}
	}:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return pc
}

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *BaseServer) AddPeer(node *discover.Node) {
	select {
	case srv.addstatic <- node:
	case <-srv.quit:
	}
}

// RemovePeer disconnects from the given node.
func (srv *BaseServer) RemovePeer(node *discover.Node) {
	select {
	case srv.removestatic <- node:
	case <-srv.quit:
	}
}

// SubscribePeers subscribes the given channel to peer events.
func (srv *BaseServer) SubscribeEvents(ch chan *PeerEvent) event.Subscription {
	return srv.peerFeed.Subscribe(ch)
}

// Self returns the local node's endpoint information.
func (srv *BaseServer) Self() *discover.Node {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return &discover.Node{IP: net.ParseIP("0.0.0.0")}
	}
	return srv.makeSelf(srv.listener, srv.ntab)
}

func (srv *BaseServer) makeSelf(listener net.Listener, discovery discover.Discovery) *discover.Node {
	// If the server's not running, return an empty node.
	// If the node is running but discovery is off, manually assemble the node infos.
	if discovery == nil {
		// Inbound connections disabled, use zero address.
		if listener == nil {
			return &discover.Node{IP: net.ParseIP("0.0.0.0"), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey)}
		}
		// Otherwise inject the listener address too
		addr := listener.Addr().(*net.TCPAddr)
		return &discover.Node{
			ID:  discover.PubkeyID(&srv.PrivateKey.PublicKey),
			IP:  addr.IP,
			TCP: uint16(addr.Port),
		}
	}
	// Otherwise return the discovery node.
	return discovery.Self()
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections are closed.
func (srv *BaseServer) Stop() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}
	srv.running = false
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}
	close(srv.quit)
	srv.loopWG.Wait()
}

// GetListenAddress returns the listen address of the server.
func (srv *BaseServer) GetListenAddress() []string {
	return []string{srv.ListenAddr}
}

// sharedUDPConn implements a shared connection. Write sends messages to the underlying connection while read returns
// messages that were found unprocessable and sent to the unhandled channel by the primary listener.
type sharedUDPConn struct {
	*net.UDPConn
	unhandled chan discover.ReadPacket
}

// ReadFromUDP implements discv5.conn
func (s *sharedUDPConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	packet, ok := <-s.unhandled
	if !ok {
		return 0, nil, fmt.Errorf("Connection was closed")
	}
	l := len(packet.Data)
	if l > len(b) {
		l = len(b)
	}
	copy(b[:l], packet.Data[:l])
	return l, packet.Addr, nil
}

// Close implements discv5.conn
func (s *sharedUDPConn) Close() error {
	return nil
}

// Start starts running the server.
// Servers can not be re-used after stopping.
func (srv *BaseServer) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	srv.logger = srv.Config.Logger
	if srv.logger == nil {
		srv.logger = logger.NewWith()
	}
	srv.logger.Info("Starting P2P networking")

	// static fields
	if srv.PrivateKey == nil {
		return fmt.Errorf("Server.PrivateKey must be set to a non-nil key")
	}

	if !srv.ConnectionType.Valid() {
		return fmt.Errorf("Invalid connection type speficied")
	}

	if srv.newTransport == nil {
		srv.newTransport = newRLPX
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.addpeer = make(chan *conn)
	srv.delpeer = make(chan peerDrop)
	srv.posthandshake = make(chan *conn)
	srv.addstatic = make(chan *discover.Node)
	srv.removestatic = make(chan *discover.Node)
	srv.peerOp = make(chan peerOpFunc)
	srv.peerOpDone = make(chan struct{})
	srv.discpeer = make(chan discover.NodeID)

	var (
		conn      *net.UDPConn
		realaddr  *net.UDPAddr
		unhandled chan discover.ReadPacket
	)

	if !srv.NoDiscovery {
		addr, err := net.ResolveUDPAddr("udp", srv.ListenAddr)
		if err != nil {
			return err
		}
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}
		realaddr = conn.LocalAddr().(*net.UDPAddr)
		if srv.NAT != nil {
			if !realaddr.IP.IsLoopback() {
				go nat.Map(srv.NAT, srv.quit, "udp", realaddr.Port, realaddr.Port, "klaytn discovery")
			}
			// TODO: react to external IP changes over time.
			if ext, err := srv.NAT.ExternalIP(); err == nil {
				realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
			}
		}
	}

	// node table
	if !srv.NoDiscovery {
		cfg := discover.Config{
			PrivateKey:   srv.PrivateKey,
			AnnounceAddr: realaddr,
			NodeDBPath:   srv.NodeDatabase,
			NetRestrict:  srv.NetRestrict,
			Bootnodes:    srv.BootstrapNodes,
			Unhandled:    unhandled,
			Conn:         conn,
			Addr:         realaddr,
			Id:           discover.PubkeyID(&srv.PrivateKey.PublicKey),
			NodeType:     ConvertNodeType(srv.ConnectionType),
			NetworkID:    srv.NetworkID,
		}

		cfgForLog := cfg
		cfgForLog.PrivateKey = nil

		logger.Info("Create udp", "config", cfgForLog)

		ntab, err := discover.ListenUDP(&cfg)
		if err != nil {
			return err
		}
		srv.ntab = ntab
	}

	dialer := newDialState(srv.StaticNodes, srv.BootstrapNodes, srv.ntab, srv.maxDialedConns(), srv.NetRestrict, srv.PrivateKey, srv.getTypeStatics())

	// handshake
	srv.ourHandshake = &protoHandshake{Version: baseProtocolVersion, Name: srv.Name(), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey), Multichannel: false}
	for _, p := range srv.Protocols {
		srv.ourHandshake.Caps = append(srv.ourHandshake.Caps, p.cap())
	}
	// listen/dial
	if srv.NoDial && srv.NoListen {
		srv.logger.Error("P2P server will be useless, neither dialing nor listening")
	}
	if !srv.NoListen {
		if srv.ListenAddr != "" {
			if err := srv.startListening(); err != nil {
				return err
			}
		} else {
			srv.logger.Error("P2P server might be useless, listening address is missing")
		}
	}

	srv.loopWG.Add(1)
	go srv.run(dialer)
	srv.running = true
	srv.logger.Info("Started P2P server", "id", discover.PubkeyID(&srv.PrivateKey.PublicKey), "multichannel", false)
	return nil
}

func (srv *BaseServer) startListening() error {
	// Launch the TCP listener.
	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	srv.ListenAddr = laddr.String()
	srv.listener = listener
	srv.loopWG.Add(1)
	go srv.listenLoop()
	// Map the TCP listening port if NAT is configured.
	if !laddr.IP.IsLoopback() && srv.NAT != nil {
		srv.loopWG.Add(1)
		go func() {
			nat.Map(srv.NAT, srv.quit, "tcp", laddr.Port, laddr.Port, "klaytn p2p")
			srv.loopWG.Done()
		}()
	}
	return nil
}

type dialer interface {
	newTasks(running int, peers map[discover.NodeID]*Peer, now time.Time) []task
	taskDone(task, time.Time)
	addStatic(*discover.Node)
	removeStatic(*discover.Node)
}

func (srv *BaseServer) run(dialstate dialer) {
	defer srv.loopWG.Done()
	var (
		peers        = make(map[discover.NodeID]*Peer)
		inboundCount = 0
		trusted      = make(map[discover.NodeID]bool, len(srv.TrustedNodes))
		taskdone     = make(chan task, maxActiveDialTasks)
		runningTasks []task
		queuedTasks  []task // tasks that can't run yet
	)
	// Put trusted nodes into a map to speed up checks.
	// Trusted peers are loaded on startup and cannot be
	// modified while the server is running.
	for _, n := range srv.TrustedNodes {
		trusted[n.ID] = true
	}

	// removes t from runningTasks
	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	// starts until max number of active tasks is satisfied
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < maxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			srv.logger.Trace("New dial task", "task", t)
			go func() { t.Do(srv); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		// Start from queue first.
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		// Query dialer for new tasks and start as many as possible now.
		if len(runningTasks) < maxActiveDialTasks {
			nt := dialstate.newTasks(len(runningTasks)+len(queuedTasks), peers, time.Now())
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}

running:
	for {
		scheduleTasks()

		select {
		case <-srv.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case n := <-srv.addstatic:
			// This channel is used by AddPeer to add to the
			// ephemeral static peer list. Add it to the dialer,
			// it will keep the node connected.
			srv.logger.Debug("Adding static node", "node", n)
			dialstate.addStatic(n)
		case n := <-srv.removestatic:
			// This channel is used by RemovePeer to send a
			// disconnect request to a peer and begin the
			// stop keeping the node connected
			srv.logger.Debug("Removing static node", "node", n)
			dialstate.removeStatic(n)
			if p, ok := peers[n.ID]; ok {
				p.Disconnect(DiscRequested)
			}
		case op := <-srv.peerOp:
			// This channel is used by Peers and PeerCount.
			op(peers)
			srv.peerOpDone <- struct{}{}
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.logger.Trace("Dial task done", "task", t)
			dialstate.taskDone(t, time.Now())
			delTask(t)
		case c := <-srv.posthandshake:
			// A connection has passed the encryption handshake so
			// the remote identity is known (but hasn't been verified yet).
			if trusted[c.id] {
				// Ensure that the trusted flag is set before checking against MaxPhysicalConnections.
				c.flags |= trustedConn
			}
			// TODO: track in-progress inbound node IDs (pre-Peer) to avoid dialing them.
			select {
			case c.cont <- srv.encHandshakeChecks(peers, inboundCount, c):
			case <-srv.quit:
				break running
			}
		case c := <-srv.addpeer:
			// At this point the connection is past the protocol handshake.
			// Its capabilities are known and the remote identity is verified.
			var err error
			err = srv.protoHandshakeChecks(peers, inboundCount, c)
			if err == nil {
				// The handshakes are done and it passed all checks.
				p, err := newPeer([]*conn{c}, srv.Protocols, srv.Config.RWTimerConfig)
				if err != nil {
					srv.logger.Error("Fail make a new peer", "err", err)
				} else {
					// If message events are enabled, pass the peerFeed
					// to the peer
					if srv.EnableMsgEvents {
						p.events = &srv.peerFeed
					}
					name := truncateName(c.name)
					srv.logger.Debug("Adding p2p peer", "name", name, "addr", c.fd.RemoteAddr(), "peers", len(peers)+1)
					go srv.runPeer(p)
					peers[c.id] = p

					if p.Inbound() {
						inboundCount++
					}
					peerCountGauge.Update(int64(len(peers)))
					peerInCountGauge.Update(int64(inboundCount))
					peerOutCountGauge.Update(int64(len(peers) - inboundCount))
				}
			}
			// The dialer logic relies on the assumption that
			// dial tasks complete after the peer has been added or
			// discarded. Unblock the task last.
			select {
			case c.cont <- err:
			case <-srv.quit:
				break running
			}
		case pd := <-srv.delpeer:
			// A peer disconnected.
			d := common.PrettyDuration(mclock.Now() - pd.created)
			pd.logger.Debug("Removing p2p peer", "duration", d, "peers", len(peers)-1, "req", pd.requested, "err", pd.err)
			delete(peers, pd.ID())

			if pd.Inbound() {
				inboundCount--
			}

			peerCountGauge.Update(int64(len(peers)))
			peerInCountGauge.Update(int64(inboundCount))
			peerOutCountGauge.Update(int64(len(peers) - inboundCount))
		case nid := <-srv.discpeer:
			if p, ok := peers[nid]; ok {
				p.Disconnect(DiscRequested)
				p.logger.Debug(fmt.Sprintf("disconnect peer"))
			}
		}
	}

	srv.logger.Trace("P2P networking is spinning down")

	// Terminate discovery. If there is a running lookup it will terminate soon.
	if srv.ntab != nil {
		srv.ntab.Close()
	}
	//if srv.DiscV5 != nil {
	//	srv.DiscV5.Close()
	//}
	// Disconnect all peers.
	for _, p := range peers {
		p.Disconnect(DiscQuitting)
	}
	// Wait for peers to shut down. Pending connections and tasks are
	// not handled here and will terminate soon-ish because srv.quit
	// is closed.
	for len(peers) > 0 {
		p := <-srv.delpeer
		p.logger.Trace("<-delpeer (spindown)", "remainingTasks", len(runningTasks))
		delete(peers, p.ID())
	}
}

func (srv *BaseServer) protoHandshakeChecks(peers map[discover.NodeID]*Peer, inboundCount int, c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.Protocols) > 0 && countMatchingProtocols(srv.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	// Repeat the encryption handshake checks because the
	// peer set might have changed between the handshakes.
	return srv.encHandshakeChecks(peers, inboundCount, c)
}

func (srv *BaseServer) encHandshakeChecks(peers map[discover.NodeID]*Peer, inboundCount int, c *conn) error {
	switch {
	case !c.is(trustedConn|staticDialedConn) && len(peers) >= srv.Config.MaxPhysicalConnections:
		return DiscTooManyPeers
	case !c.is(trustedConn) && c.is(inboundConn) && inboundCount >= srv.maxInboundConns():
		return DiscTooManyPeers
	case peers[c.id] != nil:
		return DiscAlreadyConnected
	case c.id == srv.Self().ID:
		return DiscSelf
	default:
		return nil
	}
}

func (srv *BaseServer) maxInboundConns() int {
	return srv.Config.MaxPhysicalConnections - srv.maxDialedConns()
}

func (srv *BaseServer) maxDialedConns() int {
	switch srv.ConnectionType {
	case common.CONSENSUSNODE:
		return 0
	case common.PROXYNODE:
		return 0
	case common.ENDPOINTNODE:
		if srv.NoDiscovery || srv.NoDial {
			return 0
		}
		r := srv.DialRatio
		if r == 0 {
			r = defaultDialRatio
		}
		return srv.Config.MaxPhysicalConnections / r
	case common.BOOTNODE:
		return 0 // TODO check the bn for en
	default:
		logger.Crit("[p2p.Server] UnSupported Connection Type:", "ConnectionType", srv.ConnectionType)
		return 0
	}
}

func (srv *BaseServer) getTypeStatics() map[dialType]typedStatic {
	switch srv.ConnectionType {
	case common.CONSENSUSNODE:
		tsMap := make(map[dialType]typedStatic)
		tsMap[DT_CN] = typedStatic{100, 3} // TODO-Klaytn-Node Change to literal to constant (maxNodeCount, MaxTry)
		return tsMap
	case common.PROXYNODE:
		tsMap := make(map[dialType]typedStatic)
		tsMap[DT_PN] = typedStatic{1, 3} // // TODO-Klaytn-Node Change to literal to constant (maxNodeCount, MaxTry)
		return tsMap
	case common.ENDPOINTNODE:
		tsMap := make(map[dialType]typedStatic)
		tsMap[DT_PN] = typedStatic{2, 3} // // TODO-Klaytn-Node Change to literal to constant (maxNodeCount, MaxTry)
		return tsMap
	case common.BOOTNODE:
		return nil
	default:
		logger.Crit("[p2p.Server] UnSupported Connection Type:", "ConnectionType", srv.ConnectionType)
		return nil
	}
}

type tempError interface {
	Temporary() bool
}

// listenLoop runs in its own goroutine and accepts
// inbound connections.
func (srv *BaseServer) listenLoop() {
	defer srv.loopWG.Done()
	srv.logger.Info("RLPx listener up", "self", srv.makeSelf(srv.listener, srv.ntab))

	tokens := defaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	for {
		// Wait for a handshake slot before accepting.
		<-slots

		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = srv.listener.Accept()
			if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
				srv.logger.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				srv.logger.Debug("Read error", "err", err)
				return
			}
			break
		}

		// Reject connections that do not match NetRestrict.
		if srv.NetRestrict != nil {
			if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok && !srv.NetRestrict.Contains(tcp.IP) {
				srv.logger.Debug("Rejected conn (not whitelisted in NetRestrict)", "addr", fd.RemoteAddr())
				fd.Close()
				slots <- struct{}{}
				continue
			}
		}

		fd = newMeteredConn(fd, true)
		srv.logger.Trace("Accepted connection", "addr", fd.RemoteAddr())
		go func() {
			srv.SetupConn(fd, inboundConn, nil)
			slots <- struct{}{}
		}()
	}
}

// SetupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *BaseServer) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}

	c := &conn{fd: fd, transport: srv.newTransport(fd), flags: flags, conntype: common.ConnTypeUndefined, cont: make(chan error), portOrder: ConnDefault}
	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.close(err)
		srv.logger.Trace("Setting up connection failed", "id", c.id, "err", err)
	}
	return err
}

func (srv *BaseServer) setupConn(c *conn, flags connFlag, dialDest *discover.Node) error {
	// Prevent leftover pending conns from entering the handshake.
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}

	var err error
	// Run the connection type handshake
	if c.conntype, err = c.doConnTypeHandshake(srv.ConnectionType); err != nil {
		srv.logger.Warn("Failed doConnTypeHandshake", "addr", c.fd.RemoteAddr(), "conn", c.flags,
			"conntype", c.conntype, "err", err)
		return err
	}
	srv.logger.Trace("Connection Type Trace", "addr", c.fd.RemoteAddr(), "conn", c.flags, "ConnType", c.conntype.String())

	// Run the encryption handshake.
	if c.id, err = c.doEncHandshake(srv.PrivateKey, dialDest); err != nil {
		srv.logger.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return err
	}

	clog := srv.logger.NewWith("id", c.id, "addr", c.fd.RemoteAddr(), "conn", c.flags)
	// For dialed connections, check that the remote public key matches.
	if dialDest != nil && c.id != dialDest.ID {
		clog.Trace("Dialed identity mismatch", "want", c, dialDest.ID)
		return DiscUnexpectedIdentity
	}
	err = srv.checkpoint(c, srv.posthandshake)
	if err != nil {
		clog.Trace("Rejected peer before protocol handshake", "err", err)
		return err
	}
	// Run the protocol handshake
	phs, err := c.doProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed protobuf handshake", "err", err)
		return err
	}
	if phs.ID != c.id {
		clog.Trace("Wrong devp2p handshake identity", "err", phs.ID)
		return DiscUnexpectedIdentity
	}
	c.caps, c.name, c.multiChannel = phs.Caps, phs.Name, phs.Multichannel

	err = srv.checkpoint(c, srv.addpeer)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

func truncateName(s string) string {
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

// checkpoint sends the conn to run, which performs the
// post-handshake checks for the stage (posthandshake, addpeer).
func (srv *BaseServer) checkpoint(c *conn, stage chan<- *conn) error {
	select {
	case stage <- c:
	case <-srv.quit:
		return errServerStopped
	}
	select {
	case err := <-c.cont:
		return err
	case <-srv.quit:
		return errServerStopped
	}
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *BaseServer) runPeer(p *Peer) {
	if srv.newPeerHook != nil {
		srv.newPeerHook(p)
	}

	// broadcast peer add
	srv.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.ID(),
	})

	// run the protocol
	remoteRequested, err := p.run()

	// broadcast peer drop
	srv.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	srv.delpeer <- peerDrop{p, err, remoteRequested}
}

// NodeInfo represents a short summary of the information known about the host.
type NodeInfo struct {
	ID    string `json:"id"`   // Unique node identifier (also the encryption key)
	Name  string `json:"name"` // Name of the node, including client type, version, OS, custom data
	Enode string `json:"kni"`  // Enode URL for adding this peer from remote peers
	IP    string `json:"ip"`   // IP address of the node
	Ports struct {
		Discovery int `json:"discovery"` // UDP listening port for discovery protocol
		Listener  int `json:"listener"`  // TCP listening port for RLPx
	} `json:"ports"`
	ListenAddr string                 `json:"listenAddr"`
	Protocols  map[string]interface{} `json:"protocols"`
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (srv *BaseServer) NodeInfo() *NodeInfo {
	node := srv.Self()

	// Gather and assemble the generic node infos
	info := &NodeInfo{
		Name:       srv.Name(),
		Enode:      node.String(),
		ID:         node.ID.String(),
		IP:         node.IP.String(),
		ListenAddr: srv.ListenAddr,
		Protocols:  make(map[string]interface{}),
	}
	info.Ports.Discovery = int(node.UDP)
	info.Ports.Listener = int(node.TCP)

	// Gather all the running protocol infos (only once per protocol type)
	for _, proto := range srv.Protocols {
		if _, ok := info.Protocols[proto.Name]; !ok {
			nodeInfo := interface{}("unknown")
			if query := proto.NodeInfo; query != nil {
				nodeInfo = proto.NodeInfo()
			}
			info.Protocols[proto.Name] = nodeInfo
		}
	}
	return info
}

// PeersInfo returns an array of metadata objects describing connected peers.
func (srv *BaseServer) PeersInfo() []*PeerInfo {
	// Gather all the generic and sub-protocol specific infos
	infos := make([]*PeerInfo, 0, srv.PeerCount())
	for _, peer := range srv.Peers() {
		if peer != nil {
			infos = append(infos, peer.Info())
		}
	}
	// Sort the result array alphabetically by node identifier
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].ID > infos[j].ID {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}
	return infos
}

// Disconnect tries to disconnect peer.
func (srv *BaseServer) Disconnect(destID discover.NodeID) {
	srv.discpeer <- destID
}

// CheckNilNetworkTable returns whether network table is nil.
func (srv *BaseServer) CheckNilNetworkTable() bool {
	return srv.ntab == nil
}

// Lookup performs a network search for nodes close
// to the given target. It approaches the target by querying
// nodes that are closer to it on each iteration.
// The given target does not need to be an actual node
// identifier.
func (srv *BaseServer) Lookup(target discover.NodeID, nType discover.NodeType) []*discover.Node {
	return srv.ntab.Lookup(target, nType)
}

// Resolve searches for a specific node with the given ID and NodeType.
// It returns nil if the node could not be found.
func (srv *BaseServer) Resolve(target discover.NodeID, nType discover.NodeType) *discover.Node {
	return srv.ntab.Resolve(target, nType)
}

func (srv *BaseServer) GetNodes(nType discover.NodeType, max int) []*discover.Node {
	return srv.ntab.GetNodes(nType, max)
}

// Name returns name of server.
func (srv *BaseServer) Name() string {
	return srv.Config.Name
}

// MaxPhysicalConnections returns maximum count of peers.
func (srv *BaseServer) MaxPeers() int {
	return srv.Config.MaxPhysicalConnections
}

func ConvertNodeType(ct common.ConnType) discover.NodeType {
	switch ct {
	case common.CONSENSUSNODE:
		return discover.NodeTypeCN
	case common.PROXYNODE:
		return discover.NodeTypePN
	case common.ENDPOINTNODE:
		return discover.NodeTypeEN
	case common.BOOTNODE:
		return discover.NodeTypeBN
	default:
		return discover.NodeTypeUnknown // TODO-Klaytn-Node Maybe, call panic() func or Crit()
	}
}

func ConvertConnType(nt discover.NodeType) common.ConnType {
	switch nt {
	case discover.NodeTypeCN:
		return common.CONSENSUSNODE
	case discover.NodeTypePN:
		return common.PROXYNODE
	case discover.NodeTypeEN:
		return common.ENDPOINTNODE
	case discover.NodeTypeBN:
		return common.BOOTNODE
	default:
		return common.UNKNOWNNODE
	}
}

func ConvertConnTypeToString(ct common.ConnType) string {
	switch ct {
	case common.CONSENSUSNODE:
		return "cn"
	case common.PROXYNODE:
		return "pn"
	case common.ENDPOINTNODE:
		return "en"
	case common.BOOTNODE:
		return "bn"
	default:
		return "unknown"
	}
}

func ConvertStringToConnType(s string) common.ConnType {
	st := strings.ToLower(s)
	switch st {
	case "cn":
		return common.CONSENSUSNODE
	case "pn":
		return common.PROXYNODE
	case "en":
		return common.ENDPOINTNODE
	case "bn":
		return common.BOOTNODE
	default:
		return common.UNKNOWNNODE
	}
}
