// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from p2p/dial.go (2018/06/04).
// Modified and improved for the klaytn development.

package p2p

import (
	"container/heap"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
)

const (
	// This is the amount of time spent waiting in between
	// redialing a certain node.
	dialHistoryExpiration = 30 * time.Second

	// Discovery lookups are throttled and can only run
	// once every few seconds.
	lookupInterval = 4 * time.Second

	// If no peers are found for this amount of time, the initial bootnodes are
	// attempted to be connected.
	fallbackInterval = 20 * time.Second

	// Endpoint resolution is throttled with bounded backoff.
	initialResolveDelay = 60 * time.Second
	maxResolveDelay     = time.Hour
)

// NodeDialer is used to connect to nodes in the network, typically by using
// an underlying net.Dialer but also using net.Pipe in tests.
type NodeDialer interface {
	Dial(*discover.Node) (net.Conn, error)
	DialMulti(*discover.Node) ([]net.Conn, error)
}

// TCPDialer implements the NodeDialer interface by using a net.Dialer to
// create TCP connections to nodes in the network.
type TCPDialer struct {
	*net.Dialer
}

// Dial creates a TCP connection to the node.
func (t TCPDialer) Dial(dest *discover.Node) (net.Conn, error) {
	addr := &net.TCPAddr{IP: dest.IP, Port: int(dest.TCP)}
	return t.Dialer.Dial("tcp", addr.String())
}

// DialMulti creates TCP connections to the node.
func (t TCPDialer) DialMulti(dest *discover.Node) ([]net.Conn, error) {
	var conns []net.Conn
	if dest.TCPs != nil || len(dest.TCPs) != 0 {
		conns = make([]net.Conn, 0, len(dest.TCPs))
		for _, tcp := range dest.TCPs {
			addr := &net.TCPAddr{IP: dest.IP, Port: int(tcp)}
			conn, err := t.Dialer.Dial("tcp", addr.String())
			conns = append(conns, conn)
			if err != nil {
				return nil, err
			}
		}
	}
	return conns, nil
}

// dialstate schedules dials and discovery lookups.
// it get's a chance to compute new tasks on every iteration
// of the main loop in Server.run.
type dialstate struct {
	maxDynDials int
	ntab        discover.Discovery
	netrestrict *netutil.Netlist

	lookupRunning      bool
	typedLookupRunning map[dialType]bool
	dialing            map[discover.NodeID]connFlag
	lookupBuf          []*discover.Node // current discovery lookup results
	randomNodes        []*discover.Node // filled from Table
	static             map[discover.NodeID]*dialTask
	hist               *dialHistory

	start     time.Time        // time when the dialer was first used
	bootnodes []*discover.Node // default dials when there are no peers

	tsMap map[dialType]typedStatic // tsMap holds typedStaticDial per dialType(discovery name)
}

// the dial history remembers recent dials.
type dialHistory []pastDial

// pastDial is an entry in the dial history.
type pastDial struct {
	id  discover.NodeID
	exp time.Time
}

type task interface {
	Do(Server)
}

// A dialTask is generated for each node that is dialed. Its
// fields cannot be accessed while the task is running.
type dialTask struct {
	flags        connFlag
	dest         *discover.Node
	lastResolved time.Time
	resolveDelay time.Duration
	failedTry    int
	dialType     dialType
}

// discoverTask runs discovery table operations.
// Only one discoverTask is active at any time.
// discoverTask.Do performs a random lookup.
type discoverTask struct {
	results []*discover.Node
}

// A waitExpireTask is generated if there are no other tasks
// to keep the loop in Server.run ticking.
type waitExpireTask struct {
	time.Duration
}

type typedStatic struct {
	maxNodeCount int
	maxTry       int
}

type discoverTypedStaticTask struct {
	name    dialType
	max     int
	results []*discover.Node
}

func (t *discoverTypedStaticTask) Do(srv Server) {
	// newTasks generates a lookup task whenever typed static dials are
	// necessary. Lookups need to take some time, otherwise the
	// event loop spins too fast.
	next := srv.AddLastLookup()
	if now := time.Now(); now.Before(next) {
		time.Sleep(next.Sub(now))
	}
	srv.SetLastLookupToNow()
	t.results = srv.GetNodes(convertDialT2NodeT(t.name), t.max)
	logger.Trace("discoverTypedStaticTask", "result", len(t.results))
}

func (t *discoverTypedStaticTask) String() string {
	s := fmt.Sprintf("discover TypedStaticTask: max: %d", t.max)
	if len(t.results) > 0 {
		s += fmt.Sprintf(" (%d results)", len(t.results))
	}
	return s
}

const (
	DT_UNLIMITED = dialType("DIAL_TYPE_UNLIMITED")
	DT_CN        = dialType("CN")
	DT_PN        = dialType("PN")
	DT_EN        = dialType("EN")
)

type dialType string

func convertDialT2NodeT(dt dialType) discover.NodeType {
	switch dt {
	case DT_CN:
		return discover.NodeTypeCN
	case DT_PN:
		return discover.NodeTypePN
	default:
		logger.Crit("Support only CN, PN for typed static dial", "DialType", dt)
	}
	return discover.NodeTypeUnknown
}

func newDialState(static []*discover.Node, bootnodes []*discover.Node, ntab discover.Discovery, maxdyn int,
	netrestrict *netutil.Netlist, privateKey *ecdsa.PrivateKey, tsMap map[dialType]typedStatic) *dialstate {

	if tsMap == nil {
		tsMap = make(map[dialType]typedStatic)
	}
	tsMap[DT_UNLIMITED] = typedStatic{maxTry: math.MaxInt64, maxNodeCount: math.MaxInt64}

	s := &dialstate{
		maxDynDials:        maxdyn,
		ntab:               ntab,
		netrestrict:        netrestrict,
		static:             make(map[discover.NodeID]*dialTask),
		dialing:            make(map[discover.NodeID]connFlag),
		bootnodes:          make([]*discover.Node, len(bootnodes)),
		randomNodes:        make([]*discover.Node, maxdyn/2),
		hist:               new(dialHistory),
		tsMap:              tsMap,
		typedLookupRunning: make(map[dialType]bool),
	}
	copy(s.bootnodes, bootnodes)

	if privateKey != nil {
		selfNodeID := discover.PubkeyID(&privateKey.PublicKey)

		for _, n := range static {
			if selfNodeID != n.ID {
				s.addStatic(n)
			} else {
				logger.Debug("[Dial] Ignored static node which has same id with myself", "mySelfID", selfNodeID)
			}
		}
		return s
	}

	for _, n := range static {
		s.addStatic(n)
	}
	return s
}

func (s *dialstate) addStatic(n *discover.Node) {
	s.addTypedStatic(n, DT_UNLIMITED)
}

func (s *dialstate) addTypedStatic(n *discover.Node, dType dialType) {
	// This overwrites the task instead of updating an existing
	// entry, giving users the opportunity to force a resolve operation.
	if s.static[n.ID] == nil {
		logger.Trace("[Dial] Add TypedStatic", "node", n, "dialType", dType)
		if dType != DT_UNLIMITED {
			s.static[n.ID] = &dialTask{flags: staticDialedConn | trustedConn, dest: n, dialType: dType}
		} else {
			s.static[n.ID] = &dialTask{flags: staticDialedConn, dest: n, dialType: dType}
		}

	}
}

func (s *dialstate) removeStatic(n *discover.Node) {
	// This removes a task so future attempts to connect will not be made.
	delete(s.static, n.ID)
	// This removes a previous dial timestamp so that application
	// can force a server to reconnect with chosen peer immediately.
	s.hist.remove(n.ID)

}

func (s *dialstate) newTasks(nRunning int, peers map[discover.NodeID]*Peer, now time.Time) []task {
	if s.start.IsZero() {
		s.start = now
	}

	var newtasks []task
	addDialTask := func(flag connFlag, n *discover.Node) bool {
		logger.Trace("[Dial] Try to add dialTask", "connFlag", flag, "node", n)
		if err := s.checkDial(n, peers); err != nil {
			logger.Trace("[Dial] Skipping dial candidate from discovery nodes", "id", n.ID,
				"addr", &net.TCPAddr{IP: n.IP, Port: int(n.TCP)}, "err", err)
			return false
		}
		s.dialing[n.ID] = flag
		logger.Debug("[Dial] Add dial candidate from discovery nodes", "id", n.ID, "addr",
			&net.TCPAddr{IP: n.IP, Port: int(n.TCP)})
		newtasks = append(newtasks, &dialTask{flags: flag, dest: n})
		return true
	}

	var needDynDials int
	calcNeedDynDials := func() {
		needDynDials = s.maxDynDials
		for _, p := range peers {
			if p.rws[ConnDefault].is(dynDialedConn) {
				needDynDials--
			}
		}
		for _, flag := range s.dialing {
			if flag&dynDialedConn != 0 {
				needDynDials--
			}
		}
	}

	addStaticDialTasks := func() {
		cnt := make(map[dialType]int)
		for _, t := range s.static {
			cnt[t.dialType]++
		}

		checkStaticDial := func(dt *dialTask, peers map[discover.NodeID]*Peer) error {
			err := s.checkDial(dt.dest, peers)
			if err != nil {
				return err
			}

			sd := s.static[dt.dest.ID]
			if sd.flags&staticDialedConn == 0 {
				err := fmt.Errorf("dialer: can't check conntype except staticconn [connType : %d]", sd.flags)
				logger.Error("[Dial] ", "err", err)
				return err
			}

			typeSpec, ok := s.tsMap[dt.dialType]
			if !ok {
				err := fmt.Errorf("dialer: no data for typespec [%s]", dt.dialType)
				logger.Error("[Dial] ", "err", err)
				return err
			}

			if dt.failedTry > typeSpec.maxTry {
				return errExpired
			}

			if cnt[dt.dialType] > s.tsMap[dt.dialType].maxNodeCount {
				return errExceedMaxTypedDial
			}
			return nil
		}

		for id, t := range s.static {
			err := checkStaticDial(t, peers)
			switch err {
			case errNotWhitelisted, errSelf:
				logger.Info("[Dial] Removing static dial candidate from static nodes", "id",
					t.dest.ID, "addr", &net.TCPAddr{IP: t.dest.IP, Port: int(t.dest.TCP)}, "err", err)
				delete(s.static, t.dest.ID)
				cnt[t.dialType]--
			case errExpired:
				logger.Info("[Dial] Removing expired dial candidate from static nodes", "id",
					t.dest.ID, "addr", &net.TCPAddr{IP: t.dest.IP, Port: int(t.dest.TCP)}, "dialType", t.dialType,
					"dialCount", cnt[t.dialType], "err", err)
				delete(s.static, t.dest.ID)
				cnt[t.dialType]--
			case errExceedMaxTypedDial:
				logger.Info("[Dial] Removing exceeded dial candidate from static nodes", "id",
					t.dest.ID, "addr", &net.TCPAddr{IP: t.dest.IP, Port: int(t.dest.TCP)}, "dialType", t.dialType,
					"dialCount", cnt[t.dialType], "err", err,
				)
				delete(s.static, t.dest.ID)
				cnt[t.dialType]--
			case nil:
				s.dialing[id] = t.flags
				newtasks = append(newtasks, t)
				logger.Info("[Dial] Add dial candidate from static nodes", "id", t.dest.ID,
					"NodeType", t.dest.NType, "ip", t.dest.IP, "mainPort", t.dest.TCP, "port", t.dest.TCPs)
			default:
				logger.Trace("[Dial] Skipped addStaticDial", "reason", err, "to", t.dest)
			}
		}
		// 2. add typedStaticDiscoverTask
		if s.ntab != nil { // Run DiscoveryTasks when only Discovery Mode
			for k, ts := range s.tsMap {
				if k != DT_UNLIMITED && !s.typedLookupRunning[k] && cnt[k] < ts.maxNodeCount {
					logger.Debug("[Dial] Add new discoverTypedStaticTask", "name", k)
					s.typedLookupRunning[k] = true
					maxDiscover := ts.maxNodeCount - cnt[k]
					newtasks = append(newtasks, &discoverTypedStaticTask{name: k, max: maxDiscover})
				}
			}
		}
	}

	// Compute number of dynamic dials necessary at this point.
	calcNeedDynDials()
	logger.Trace("[Dial] Dynamic Dials Remained Capacity", "needDynDials", needDynDials, "maxDynDials", s.maxDynDials)

	// Expire the dial history on every invocation.
	s.hist.expire(now)

	// Create dials for static nodes if they are not connected.
	addStaticDialTasks()

	// Use random nodes from the table for half of the necessary
	// dynamic dials.
	randomCandidates := needDynDials / 2
	if randomCandidates > 0 {
		n := s.ntab.ReadRandomNodes(s.randomNodes, discover.NodeTypeEN)
		for i := 0; i < randomCandidates && i < n; i++ {
			if addDialTask(dynDialedConn, s.randomNodes[i]) {
				needDynDials--
			}
		}
	}
	// Create dynamic dials from random lookup results, removing tried
	// items from the result buffer.
	i := 0
	for ; i < len(s.lookupBuf) && needDynDials > 0; i++ {
		if addDialTask(dynDialedConn, s.lookupBuf[i]) {
			needDynDials--
		}
	}
	s.lookupBuf = s.lookupBuf[:copy(s.lookupBuf, s.lookupBuf[i:])]
	// Launch a discovery lookup if more candidates are needed.
	if len(s.lookupBuf) < needDynDials && !s.lookupRunning {
		s.lookupRunning = true
		newtasks = append(newtasks, &discoverTask{})
	}

	// Launch a timer to wait for the next node to expire if all
	// candidates have been tried and no task is currently active.
	// This should prevent cases where the dialer logic is not ticked
	// because there are no pending events.
	if nRunning == 0 && len(newtasks) == 0 && s.hist.Len() > 0 {
		t := &waitExpireTask{s.hist.min().exp.Sub(now)}
		newtasks = append(newtasks, t)
	}
	return newtasks
}

var (
	errSelf               = errors.New("is self")
	errAlreadyDialing     = errors.New("already dialing")
	errAlreadyConnected   = errors.New("already connected")
	errRecentlyDialed     = errors.New("recently dialed")
	errNotWhitelisted     = errors.New("not contained in netrestrict whitelist")
	errExpired            = errors.New("is expired")
	errExceedMaxTypedDial = errors.New("exceeded max typed dial")
	errUpdateDial         = errors.New("updated to be multichannel peer")
)

func (s *dialstate) checkDial(n *discover.Node, peers map[discover.NodeID]*Peer) error {
	_, dialing := s.dialing[n.ID]
	switch {
	case dialing:
		return errAlreadyDialing
	case peers[n.ID] != nil:
		return errAlreadyConnected
	case s.ntab != nil && n.ID == s.ntab.Self().ID:
		return errSelf
	case s.netrestrict != nil && !s.netrestrict.Contains(n.IP):
		return errNotWhitelisted
	case s.hist.contains(n.ID):
		return errRecentlyDialed
	}
	return nil
}

func (s *dialstate) taskDone(t task, now time.Time) {
	switch t := t.(type) {
	case *dialTask:
		s.hist.add(t.dest.ID, now.Add(dialHistoryExpiration))
		delete(s.dialing, t.dest.ID)
	case *discoverTask:
		s.lookupRunning = false
		s.lookupBuf = append(s.lookupBuf, t.results...)
	case *discoverTypedStaticTask:
		logger.Trace("[Dial] discoverTypedStaticTask - done", "t.name", t.name,
			"result count", len(t.results))
		s.typedLookupRunning[t.name] = false
		for _, r := range t.results {
			s.addTypedStatic(r, t.name)
		}
	}
}

func (t *dialTask) Do(srv Server) {
	if t.dest.Incomplete() {
		if !t.resolve(srv, t.dest.NType) {
			return
		}
	}
	var err error
	if len(t.dest.TCPs) > 1 {
		err = t.dialMulti(srv, t.dest)
	} else {
		err = t.dial(srv, t.dest)
	}

	if err != nil {
		logger.Debug("[Dial] Failed dialing", "task", t, "err", err)
		// Try resolving the ID of static nodes if dialing failed.
		if _, ok := err.(*dialError); ok && t.flags&staticDialedConn != 0 && t.resolve(srv, t.dest.NType) {
			if len(t.dest.TCPs) > 1 {
				err = t.dialMulti(srv, t.dest)
			} else {
				err = t.dial(srv, t.dest)
			}
		}

		// redial with updated connection
		if err == errUpdateDial {
			err = t.dialMulti(srv, t.dest)
		}

		if err != nil {
			t.failedTry++
		}
	}
}

// resolve attempts to find the current endpoint for the destination
// using discovery.
//
// Resolve operations are throttled with backoff to avoid flooding the
// discovery network with useless queries for nodes that don't exist.
// The backoff delay resets when the node is found.
func (t *dialTask) resolve(srv Server, nType discover.NodeType) bool {
	if srv.CheckNilNetworkTable() {
		logger.Debug("Can't resolve node", "id", t.dest.ID, "NodeType", nType,
			"err", "discovery is disabled")
		return false
	}
	if t.resolveDelay == 0 {
		t.resolveDelay = initialResolveDelay
	}
	if time.Since(t.lastResolved) < t.resolveDelay {
		return false
	}
	resolved := srv.Resolve(t.dest.ID, nType)
	t.lastResolved = time.Now()
	if resolved == nil {
		t.resolveDelay *= 2
		if t.resolveDelay > maxResolveDelay {
			t.resolveDelay = maxResolveDelay
		}
		logger.Debug("Resolving node failed", "id", t.dest.ID, "newdelay", t.resolveDelay)
		return false
	}
	// The node was found.
	t.resolveDelay = initialResolveDelay
	t.dest = resolved
	logger.Debug("Resolved node", "id", t.dest.ID, "addr", &net.TCPAddr{IP: t.dest.IP, Port: int(t.dest.TCP)})
	return true
}

type dialError struct {
	error
}

// dial performs the actual connection attempt.
func (t *dialTask) dial(srv Server, dest *discover.Node) error {
	dialTryCounter.Inc(1)
	logger.Debug("[Dial] Dialing node", "id", dest.ID, "addr", &net.TCPAddr{IP: dest.IP, Port: int(dest.TCP)})

	fd, err := srv.Dial(dest)
	if err != nil {
		dialFailCounter.Inc(1)
		return &dialError{err}
	}
	mfd := newMeteredConn(fd, false)
	return srv.SetupConn(mfd, t.flags, dest)
}

// dialMulti performs the actual connection attempt.
func (t *dialTask) dialMulti(srv Server, dest *discover.Node) error {
	dialTryCounter.Inc(1)
	addresses := make([]*net.TCPAddr, 0, len(dest.TCPs))
	for _, tcp := range dest.TCPs {
		addresses = append(addresses, &net.TCPAddr{IP: dest.IP, Port: int(tcp)})
	}
	logger.Debug("[Dial] Dialing node", "id", dest.ID, "addresses", addresses)

	fds, err := srv.DialMulti(dest)
	if err != nil {
		dialFailCounter.Inc(1)
		return &dialError{err}
	}

	var errorBackup error
	for portOrder, fd := range fds {
		mfd := newMeteredConn(fd, false)
		dest.PortOrder = uint16(portOrder)
		err := srv.SetupConn(mfd, t.flags, dest)
		if err != nil {
			errorBackup = err
		}
	}
	if errorBackup != nil {
		for _, fd := range fds {
			fd.Close()
		}
	}
	return errorBackup
}

func (t *dialTask) String() string {
	return fmt.Sprintf("%v %x %v:%d", t.flags, t.dest.ID[:8], t.dest.IP, t.dest.TCP)
}

func (t *discoverTask) Do(srv Server) {
	// newTasks generates a lookup task whenever dynamic dials are
	// necessary. Lookups need to take some time, otherwise the
	// event loop spins too fast.
	next := srv.AddLastLookup()
	if now := time.Now(); now.Before(next) {
		logger.Trace("discoverTask sleep", "period", next.Sub(now))
		time.Sleep(next.Sub(now))
	}
	logger.Trace("discoverTask wakeup")
	srv.SetLastLookupToNow()
	var target discover.NodeID
	rand.Read(target[:])
	t.results = srv.Lookup(target, discover.NodeTypeEN) // TODO-Klaytn Supposed dynamicDial discover only en, but type have to get from argument.
}

func (t *discoverTask) String() string {
	s := "discovery lookup"
	if len(t.results) > 0 {
		s += fmt.Sprintf(" (%d results)", len(t.results))
	}
	return s
}

func (t waitExpireTask) Do(Server) {
	time.Sleep(t.Duration)
}
func (t waitExpireTask) String() string {
	return fmt.Sprintf("wait for dial hist expire (%v)", t.Duration)
}

// Use only these methods to access or modify dialHistory.
func (h dialHistory) min() pastDial {
	return h[0]
}
func (h *dialHistory) add(id discover.NodeID, exp time.Time) {
	heap.Push(h, pastDial{id, exp})

}
func (h *dialHistory) remove(id discover.NodeID) bool {
	for i, v := range *h {
		if v.id == id {
			heap.Remove(h, i)
			return true
		}
	}
	return false
}
func (h dialHistory) contains(id discover.NodeID) bool {
	for _, v := range h {
		if v.id == id {
			return true
		}
	}
	return false
}
func (h *dialHistory) expire(now time.Time) {
	for h.Len() > 0 && h.min().exp.Before(now) {
		heap.Pop(h)
	}
}

// heap.Interface boilerplate
func (h dialHistory) Len() int           { return len(h) }
func (h dialHistory) Less(i, j int) bool { return h[i].exp.Before(h[j].exp) }
func (h dialHistory) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *dialHistory) Push(x interface{}) {
	*h = append(*h, x.(pastDial))
}
func (h *dialHistory) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
