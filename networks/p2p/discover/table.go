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
// This file is derived from p2p/discover/table.go (2018/06/04).
// Modified and improved for the klaytn development.

package discover

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	mrand "math/rand"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
)

const (
	alpha           = 3  // Kademlia concurrency factor
	bucketSize      = 16 // Kademlia bucket size
	maxReplacements = 10 // Size of per-bucket replacement list

	maxBondingPingPongs = 16 // Limit on the number of concurrent ping/pong interactions
	maxFindnodeFailures = 5  // Nodes exceeding this limit are dropped

	refreshInterval    = 30 * time.Minute
	revalidateInterval = 10 * time.Second
	copyNodesInterval  = 30 * time.Second

	seedCount  = 30
	seedMaxAge = 5 * 24 * time.Hour
)

type DiscoveryType uint8

type Discovery interface {
	Self() *Node
	Close()
	Resolve(target NodeID, targetType NodeType) *Node
	Lookup(target NodeID, targetType NodeType) []*Node
	GetNodes(targetType NodeType, max int) []*Node
	ReadRandomNodes([]*Node, NodeType) int
	RetrieveNodes(target common.Hash, nType NodeType, nresults int) []*Node // replace of closest():Table

	HasBond(id NodeID) bool
	Bond(pinged bool, id NodeID, addr *net.UDPAddr, tcpPort uint16, nType NodeType) (*Node, error)
	IsAuthorized(fromID NodeID, nType NodeType) bool

	// interfaces for API
	Name() string
	CreateUpdateNodeOnDB(n *Node) error
	CreateUpdateNodeOnTable(n *Node) error
	GetNodeFromDB(id NodeID) (*Node, error)
	DeleteNodeFromDB(n *Node) error
	DeleteNodeFromTable(n *Node) error
	GetBucketEntries() []*Node
	GetReplacements() []*Node

	GetAuthorizedNodes() []*Node
	PutAuthorizedNodes(nodes []*Node)
	DeleteAuthorizedNodes(nodes []*Node)
}

type Table struct {
	nursery []*Node     // bootstrap nodes
	rand    *mrand.Rand // source of randomness, periodically reseeded
	randMu  sync.Mutex
	ips     netutil.DistinctNetSet

	db         *nodeDB // database of known nodes
	refreshReq chan chan struct{}
	initDone   chan struct{}
	closeReq   chan struct{}
	closed     chan struct{}

	bondmu    sync.Mutex
	bonding   map[NodeID]*bondproc
	bondslots chan struct{} // limits total number of active bonding processes

	nodeAddedHook func(*Node) // for testing

	net  transport
	self *Node // metadata of the local node

	storages   map[NodeType]discoverStorage
	storagesMu sync.RWMutex

	localLogger log.Logger
}

type bondproc struct {
	err  error
	n    *Node
	done chan struct{}
}

// transport is implemented by the UDP transport.
// it is an interface so we can test without opening lots of UDP
// sockets and without generating a private key.
type transport interface {
	ping(toid NodeID, toaddr *net.UDPAddr) error
	waitping(NodeID) error
	findnode(toid NodeID, toaddr *net.UDPAddr, target NodeID, targetNT NodeType, max int) ([]*Node, error)
	close()
}

func newTable(cfg *Config) (Discovery, error) {
	// If no node database was given, use an in-memory one
	db, err := newNodeDB(cfg.NodeDBPath, Version, cfg.Id)
	if err != nil {
		return nil, err
	}

	tab := &Table{
		net:         cfg.udp,
		db:          db,
		self:        NewNode(cfg.Id, cfg.Addr.IP, uint16(cfg.Addr.Port), uint16(cfg.Addr.Port), nil, cfg.NodeType),
		bonding:     make(map[NodeID]*bondproc),
		bondslots:   make(chan struct{}, maxBondingPingPongs),
		refreshReq:  make(chan chan struct{}),
		initDone:    make(chan struct{}),
		closeReq:    make(chan struct{}),
		closed:      make(chan struct{}),
		rand:        mrand.New(mrand.NewSource(0)),
		storages:    make(map[NodeType]discoverStorage),
		localLogger: logger.NewWith("Discover", "Table"),
	}

	switch cfg.NodeType {
	case NodeTypeCN:
		tab.addStorage(NodeTypeCN, &simpleStorage{targetType: NodeTypeCN, noDiscover: true, max: 100})
		tab.addStorage(NodeTypeBN, &simpleStorage{targetType: NodeTypeBN, noDiscover: true, max: 3})
	case NodeTypePN:
		tab.addStorage(NodeTypePN, &simpleStorage{targetType: NodeTypePN, noDiscover: true, max: 1})
		tab.addStorage(NodeTypeEN, &KademliaStorage{targetType: NodeTypeEN, noDiscover: true})
		tab.addStorage(NodeTypeBN, &simpleStorage{targetType: NodeTypeBN, noDiscover: true, max: 3})
	case NodeTypeEN:
		tab.addStorage(NodeTypePN, &simpleStorage{targetType: NodeTypePN, noDiscover: true, max: 2})
		tab.addStorage(NodeTypeEN, &KademliaStorage{targetType: NodeTypeEN})
		tab.addStorage(NodeTypeBN, &simpleStorage{targetType: NodeTypeBN, noDiscover: true, max: 3})
	case NodeTypeBN:
		tab.addStorage(NodeTypeCN, NewSimpleStorage(NodeTypeCN, true, 100, cfg.AuthorizedNodes))
		tab.addStorage(NodeTypePN, NewSimpleStorage(NodeTypePN, true, 100, cfg.AuthorizedNodes))
		tab.addStorage(NodeTypeEN, &KademliaStorage{targetType: NodeTypeEN, noDiscover: true})
		tab.addStorage(NodeTypeBN, &simpleStorage{targetType: NodeTypeBN, max: 3})
	}

	if err := tab.setFallbackNodes(cfg.Bootnodes); err != nil {
		return nil, err
	}
	for i := 0; i < cap(tab.bondslots); i++ {
		tab.bondslots <- struct{}{}
	}

	tab.seedRand()
	tab.loadSeedNodes(false)
	// Start the background expiration goroutine after loading seeds so that the search for
	// seed nodes also considers older nodes that would otherwise be removed by the
	// expiration.
	tab.db.ensureExpirer()
	tab.localLogger.Debug("new "+tab.Name()+" created", "err", nil)
	return tab, nil
}

func (tab *Table) IsAuthorized(fromID NodeID, nType NodeType) bool {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	if tab.storages[nType] != nil {
		return tab.storages[nType].isAuthorized(fromID)
	}
	return true
}

// setFallbackNodes sets the initial points of contact. These nodes
// are used to connect to the network if the table is empty and there
// are no known nodes in the database.
func (tab *Table) setFallbackNodes(nodes []*Node) error {
	for _, n := range nodes {
		if err := n.validateComplete(); err != nil {
			return fmt.Errorf("bad bootstrap/fallback node %q (%v)", n, err)
		}
	}
	tab.nursery = make([]*Node, 0, len(nodes))
	for _, n := range nodes {
		cpy := *n
		// Recompute cpy.sha because the node might not have been
		// created by NewNode or ParseNode.
		cpy.sha = crypto.Keccak256Hash(n.ID[:])
		tab.nursery = append(tab.nursery, &cpy)
	}
	return nil
}

func (tab *Table) findNewNode(seeds *nodesByDistance, targetID NodeID, targetNT NodeType, recursiveFind bool, max int) []*Node {
	var (
		asked          = make(map[NodeID]bool)
		seen           = make(map[NodeID]bool)
		reply          = make(chan []*Node, alpha)
		pendingQueries = 0
	)

	// don't query further if we hit ourself.
	// unlikely to happen often in practice.
	asked[tab.self.ID] = true
	for _, e := range seeds.entries {
		seen[e.ID] = true
	}

	for {
		// ask the alpha closest nodes that we haven't asked yet
		for i := 0; i < len(seeds.entries) && pendingQueries < alpha; i++ {
			n := seeds.entries[i]
			if !asked[n.ID] {
				asked[n.ID] = true
				pendingQueries++
				go func() {
					// Find potential neighbors to bond with
					r, err := tab.net.findnode(n.ID, n.addr(), targetID, targetNT, max)
					if err != nil {
						// Bump the failure counter to detect and evacuate non-bonded entries
						fails := tab.db.findFails(n.ID) + 1
						tab.db.updateFindFails(n.ID, fails)
						tab.localLogger.Trace("Bumping findnode failure counter", "id", n.ID, "failcount", fails)

						if fails >= maxFindnodeFailures {
							tab.localLogger.Trace("Too many findnode failures, dropping", "id", n.ID, "failcount", fails)
							tab.delete(n)
						}
					}
					if targetNT != NodeTypeBN {
						r = removeBn(r)
					}
					reply <- tab.bondall(r)
				}()
			}
		}
		if pendingQueries == 0 {
			// we have asked all closest nodes, stop the search
			break
		}

		if recursiveFind {
			// wait for the next reply
			for _, n := range <-reply {
				if n != nil && !seen[n.ID] {
					seen[n.ID] = true
					seeds.push(n, max)
				}
			}
			pendingQueries--
		} else {
			for i := 0; i < pendingQueries; i++ {
				for _, n := range <-reply {
					if n != nil && !seen[n.ID] {
						seen[n.ID] = true
						seeds.push(n, max)
					}
				}
			}
			break
		}
	}
	if targetNT != NodeTypeBN {
		seeds.entries = removeBn(seeds.entries)
	}
	tab.localLogger.Debug("findNewNode: found nodes", "length", len(seeds.entries), "nodeType", targetNT)
	return seeds.entries
}

func (tab *Table) addStorage(nType NodeType, s discoverStorage) {
	tab.storagesMu.Lock()
	defer tab.storagesMu.Unlock()
	s.setTable(tab)
	tab.storages[nType] = s
	s.init()
}

func (tab *Table) seedRand() {
	var b [8]byte
	crand.Read(b[:])

	//tab.mutex.Lock()
	tab.randMu.Lock()
	tab.rand.Seed(int64(binary.BigEndian.Uint64(b[:])))
	tab.randMu.Unlock()
	//tab.mutex.Unlock()
}

// Self returns the local node.
// The returned node should not be modified by the caller.
func (tab *Table) Self() *Node {
	return tab.self
}

// ReadRandomNodes fills the given slice with random nodes from the
// table. It will not write the same node more than once. The nodes in
// the slice are copies and can be modified by the caller.
func (tab *Table) ReadRandomNodes(buf []*Node, nType NodeType) (n int) {
	if !tab.isInitDone() {
		return 0
	}

	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	if tab.storages[nType] == nil {
		tab.localLogger.Warn("ReadRandomNodes: Not Supported NodeType", "NodeType", nType)
		return 0
	}

	return tab.storages[nType].readRandomNodes(buf)
}

// Close terminates the network listener and flushes the node database.
func (tab *Table) Close() {
	select {
	case <-tab.closed:
		// already closed.
	case tab.closeReq <- struct{}{}:
		<-tab.closed // wait for refreshLoop to end.
	}
}

// isInitDone returns whether the table's initial seeding procedure has completed.
func (tab *Table) isInitDone() bool {
	select {
	case <-tab.initDone:
		return true
	default:
		return false
	}
}

// Resolve searches for a specific node with the given ID.
// It returns nil if the node could not be found.
func (tab *Table) Resolve(targetID NodeID, targetType NodeType) *Node {
	// If the node is present in the local table, no
	// network interaction is required.
	hash := crypto.Keccak256Hash(targetID[:])
	cl := tab.closest(hash, targetType, 1)
	if len(cl.entries) > 0 && cl.entries[0].ID == targetID {
		return cl.entries[0]
	}
	// Otherwise, do a network lookup.
	result := tab.Lookup(targetID, targetType)
	for _, n := range result {
		if n.ID == targetID {
			return n
		}
	}
	return nil
}

// Lookup performs a network search for nodes close
// to the given target. It approaches the target by querying
// nodes that are closer to it on each iteration.
// The given target does not need to be an actual node
// identifier.
func (tab *Table) Lookup(targetID NodeID, targetType NodeType) []*Node {
	return tab.lookup(targetID, true, targetType)
}

func (tab *Table) lookup(targetID NodeID, refreshIfEmpty bool, targetNT NodeType) []*Node {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()

	if tab.storages[targetNT] == nil {
		tab.localLogger.Warn("lookup: Not Supported NodeType", "NodeType", targetNT)
		return []*Node{}
	}
	return tab.storages[targetNT].lookup(targetID, refreshIfEmpty, targetNT)
}

func (tab *Table) GetNodes(targetNT NodeType, max int) []*Node {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()

	if tab.storages[targetNT] == nil {
		tab.localLogger.Warn("getNodes: Not Supported NodeType", "NodeType", targetNT)
		return []*Node{}
	}
	return tab.storages[targetNT].getNodes(max)
}

func removeBn(nodes []*Node) []*Node {
	tmp := nodes[:0]
	for _, n := range nodes {
		if n.NType != NodeTypeBN {
			tmp = append(tmp, n)
		}
	}
	return tmp
}

func (tab *Table) refresh() <-chan struct{} {
	done := make(chan struct{})
	select {
	case tab.refreshReq <- done:
	case <-tab.closed:
		close(done)
	}
	return done
}

// loop schedules refresh, revalidate runs and coordinates shutdown.
func (tab *Table) loop() {
	var (
		revalidate     = time.NewTimer(tab.nextRevalidateTime())
		refresh        = time.NewTicker(refreshInterval)
		copyNodes      = time.NewTicker(copyNodesInterval)
		revalidateDone = make(chan struct{})
		refreshDone    = make(chan struct{})           // where doRefresh reports completion
		waiting        = []chan struct{}{tab.initDone} // holds waiting callers while doRefresh runs
	)
	defer refresh.Stop()
	defer revalidate.Stop()
	defer copyNodes.Stop()

	// Start initial refresh.
	go tab.doRefresh(refreshDone)

loop:
	for {
		select {
		case <-refresh.C:
			tab.seedRand()
			if refreshDone == nil {
				refreshDone = make(chan struct{})
				go tab.doRefresh(refreshDone)
			}
		case req := <-tab.refreshReq:
			waiting = append(waiting, req)
			if refreshDone == nil {
				refreshDone = make(chan struct{})
				go tab.doRefresh(refreshDone)
			}
		case <-refreshDone:
			for _, ch := range waiting {
				close(ch)
			}
			waiting, refreshDone = nil, nil
		case <-revalidate.C:
			go tab.doRevalidate(revalidateDone)
		case <-revalidateDone:
			tt := tab.nextRevalidateTime()
			revalidate.Reset(tt)
		case <-copyNodes.C:
			go tab.copyBondedNodes()
		case <-tab.closeReq:
			break loop
		}
	}

	if tab.net != nil {
		tab.net.close()
	}
	if refreshDone != nil {
		<-refreshDone
	}
	for _, ch := range waiting {
		close(ch)
	}
	tab.db.close()
	close(tab.closed)
}

// doRefresh performs a lookup for a random target to keep buckets
// full. seed nodes are inserted if the table is empty (initial
// bootstrap or discarded faulty peers).
func (tab *Table) doRefresh(done chan struct{}) {
	tab.localLogger.Trace("doRefresh()")
	defer close(done)

	// Load nodes from the database and insert
	// them. This should yield a few previously seen nodes that are
	// (hopefully) still alive.
	tab.loadSeedNodes(true)

	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	for _, ds := range tab.storages {
		ds.doRefresh()
	}
}

func (tab *Table) loadSeedNodes(bond bool) {
	// TODO-Klaytn-Node Separate logic to storages.
	seeds := tab.db.querySeeds(seedCount, seedMaxAge)
	seeds = removeBn(seeds)
	seeds = append(seeds, tab.nursery...)
	if bond {
		seeds = tab.bondall(seeds)
	}
	for i := range seeds {
		seed := seeds[i]
		age := log.Lazy{Fn: func() interface{} { return time.Since(tab.db.bondTime(seed.ID)) }}
		tab.localLogger.Debug("Found seed node in database", "id", seed.ID, "addr", seed.addr(), "age", age)
		tab.add(seed)
	}
}

// doRevalidate checks that the last node in a random bucket is still live
// and replaces or deletes the node if it isn't.
func (tab *Table) doRevalidate(done chan<- struct{}) {
	defer func() { done <- struct{}{} }()

	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	for _, ds := range tab.storages {
		ds.doRevalidate()
	}
}

func (tab *Table) nextRevalidateTime() time.Duration {
	tab.randMu.Lock()
	defer tab.randMu.Unlock()

	return time.Duration(tab.rand.Int63n(int64(revalidateInterval)))
}

// copyBondedNodes adds nodes from the table to the database if they have been in the table
// longer then minTableTime.
func (tab *Table) copyBondedNodes() {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	for _, ds := range tab.storages {
		ds.copyBondedNodes()
	}
}

// closest returns the n nodes in the table that are closest to the
// given id. The caller must hold tab.mutex.
func (tab *Table) closest(target common.Hash, nType NodeType, nresults int) *nodesByDistance {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()

	if tab.storages[nType] == nil {
		tab.localLogger.Warn("closest(): Not Supported NodeType", "NodeType", nType)
		return &nodesByDistance{}
	}
	return tab.storages[nType].closest(target, nresults)
}

// RetrieveNodes returns node list except bootnode. This method is used to make a result of FINDNODE request.
func (tab *Table) RetrieveNodes(target common.Hash, nType NodeType, nresults int) []*Node {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()

	if tab.storages[nType] == nil {
		tab.localLogger.Warn("RetrieveNodes: Not Supported NodeType", "NodeType", nType)
		return []*Node{}
	}
	nodes := tab.storages[nType].closest(target, nresults).entries
	if nType != NodeTypeBN {
		nodes = removeBn(nodes)
	}
	return nodes
}

func (tab *Table) len() (n int) {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()

	for _, ds := range tab.storages {
		n += ds.len()
	}
	return n
}

func (tab *Table) nodes() (n []*Node) {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()

	for _, ds := range tab.storages {
		n = append(n, ds.nodeAll()...)
	}
	return n
}

// bondall bonds with all given nodes concurrently and returns
// those nodes for which bonding has probably succeeded.
func (tab *Table) bondall(nodes []*Node) (result []*Node) {
	rc := make(chan *Node, len(nodes))
	for i := range nodes {
		go func(n *Node) {
			nn, _ := tab.Bond(false, n.ID, n.addr(), n.TCP, n.NType)
			rc <- nn
		}(nodes[i])
	}
	for range nodes {
		if n := <-rc; n != nil {
			result = append(result, n)
		}
	}
	return result
}

// Bond ensures the local node has a bond with the given remote node.
// It also attempts to insert the node into the table if bonding succeeds.
// The caller must not hold tab.mutex.
//
// A bond is must be established before sending findnode requests.
// Both sides must have completed a ping/pong exchange for a bond to
// exist. The total number of active bonding processes is limited in
// order to restrain network use.
//
// bond is meant to operate idempotently in that bonding with a remote
// node which still remembers a previously established bond will work.
// The remote node will simply not send a ping back, causing waitping
// to time out.
//
// If pinged is true, the remote node has just pinged us and one half
// of the process can be skipped.
func (tab *Table) Bond(pinged bool, id NodeID, addr *net.UDPAddr, tcpPort uint16, nType NodeType) (*Node, error) {
	if id == tab.self.ID {
		return nil, errors.New("is self")
	}
	if pinged && !tab.isInitDone() {
		return nil, errors.New("still initializing")
	}
	// Start bonding if we haven't seen this node for a while or if it failed findnode too often.
	node, fails := tab.db.node(id), tab.db.findFails(id)
	age := time.Since(tab.db.bondTime(id))
	var result error
	// A Bootnode always add node(cn, pn, en) to table.
	if fails > 0 || age > nodeDBNodeExpiration || (node == nil && tab.self.NType == NodeTypeBN) {
		tab.localLogger.Trace("Bond - Starting bonding ping/pong", "id", id, "known", node != nil, "failcount", fails, "age", age)

		tab.bondmu.Lock()
		w := tab.bonding[id]
		if w != nil {
			// Wait for an existing bonding process to complete.
			tab.bondmu.Unlock()
			<-w.done
		} else {
			// Register a new bonding process.
			w = &bondproc{done: make(chan struct{})}
			tab.bonding[id] = w
			tab.bondmu.Unlock()
			// Do the ping/pong. The result goes into w.
			tab.pingpong(w, pinged, id, addr, tcpPort, nType)
			// Unregister the process after it's done.
			tab.bondmu.Lock()
			delete(tab.bonding, id)
			tab.bondmu.Unlock()
		}
		// Retrieve the bonding results
		result = w.err
		tab.localLogger.Trace("Bond", "error", result)
		if result == nil {
			node = w.n
		}
	}
	// Add the node to the table even if the bonding ping/pong
	// fails. It will be replaced quickly if it continues to be
	// unresponsive.
	if node != nil {
		tab.localLogger.Trace("Bond - Add", "id", node.ID, "type", node.NType, "sha", node.sha)
		tab.add(node)
		tab.db.updateFindFails(id, 0)
		lenEntries := len(tab.GetBucketEntries())
		lenReplacements := len(tab.GetReplacements())
		bucketEntriesGauge.Update(int64(lenEntries))
		bucketReplacementsGauge.Update(int64(lenReplacements))
	}
	return node, result
}

func (tab *Table) pingpong(w *bondproc, pinged bool, id NodeID, addr *net.UDPAddr, tcpPort uint16, nType NodeType) {
	// Request a bonding slot to limit network usage
	<-tab.bondslots
	defer func() { tab.bondslots <- struct{}{} }()

	// Ping the remote side and wait for a pong.
	if w.err = tab.ping(id, addr); w.err != nil {
		close(w.done)
		return
	}
	if !pinged {
		// Give the remote node a chance to ping us before we start
		// sending findnode requests. If they still remember us,
		// waitping will simply time out.
		tab.localLogger.Trace("pingpong-waitping", "to", id)
		tab.net.waitping(id)
	}
	// Bonding succeeded, update the node database.
	w.n = NewNode(id, addr.IP, uint16(addr.Port), tcpPort, nil, nType)
	tab.localLogger.Trace("pingpong-success, make new node", "node", w.n)
	close(w.done)
}

// ping a remote endpoint and wait for a reply, also updating the node
// database accordingly.
func (tab *Table) ping(id NodeID, addr *net.UDPAddr) error {
	tab.localLogger.Trace("ping", "to", id)
	tab.db.updateLastPing(id, time.Now())
	if err := tab.net.ping(id, addr); err != nil {
		return err
	}
	tab.db.updateBondTime(id, time.Now())
	return nil
}

// bucket returns the bucket for the given node ID hash.
// This method is for only unit tests.
func (tab *Table) bucket(sha common.Hash, nType NodeType) *bucket {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()

	if tab.storages[nType] == nil {
		tab.localLogger.Warn("bucket(): Not Supported NodeType", "NodeType", nType)
		return &bucket{}
	}
	if _, ok := tab.storages[nType].(*KademliaStorage); !ok {
		tab.localLogger.Warn("bucket(): bucket() only allowed to use at KademliaStorage", "NodeType", nType)
		return &bucket{}
	}
	ks := tab.storages[nType].(*KademliaStorage)

	ks.bucketsMu.Lock()
	defer ks.bucketsMu.Unlock()
	return ks.bucket(sha)
}

// add attempts to add the given node its corresponding bucket. If the
// bucket has space available, adding the node succeeds immediately.
// Otherwise, the node is added if the least recently active node in
// the bucket does not respond to a ping packet.
//
// The caller must not hold tab.mutex.
func (tab *Table) add(new *Node) {
	tab.localLogger.Trace("add(node)", "NodeType", new.NType, "node", new, "sha", new.sha)
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	if new.NType == NodeTypeBN {
		for _, ds := range tab.storages {
			ds.add(new)
		}
	} else {
		if tab.storages[new.NType] == nil {
			tab.localLogger.Warn("add(): Not Supported NodeType", "NodeType", new.NType)
			return
		}
		tab.storages[new.NType].add(new)
	}
}

// stuff adds nodes the table to the end of their corresponding bucket
// if the bucket is not full.
func (tab *Table) stuff(nodes []*Node, nType NodeType) {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	if tab.storages[nType] == nil {
		tab.localLogger.Warn("stuff(): Not Supported NodeType", "NodeType", nType)
		return
	}
	tab.storages[nType].stuff(nodes)
}

// delete removes an entry from the node table (used to evacuate
// failed/non-bonded discovery peers).
func (tab *Table) delete(node *Node) {
	tab.storagesMu.RLock()
	defer tab.storagesMu.RUnlock()
	for _, ds := range tab.storages {
		ds.delete(node)
	}
}

func (tab *Table) HasBond(id NodeID) bool {
	return tab.db.hasBond(id)
}

// nodesByDistance is a list of nodes, ordered by
// distance to target.
type nodesByDistance struct {
	entries []*Node
	target  common.Hash
}

// push adds the given node to the list, keeping the total size below maxElems.
func (h *nodesByDistance) push(n *Node, maxElems int) {
	ix := sort.Search(len(h.entries), func(i int) bool {
		return distcmp(h.target, h.entries[i].sha, n.sha) > 0
	})
	if len(h.entries) < maxElems {
		h.entries = append(h.entries, n)
	}
	if ix == len(h.entries) {
		// farther away than all nodes we already have.
		// if there was room for it, the node is now the last element.
	} else {
		// slide existing entries down to make room
		// this will overwrite the entry we just appended.
		copy(h.entries[ix+1:], h.entries[ix:])
		h.entries[ix] = n
	}
}

func (h *nodesByDistance) String() string {
	return fmt.Sprintf("nodeByDistance target: %s, entries: %s", h.target, h.entries)
}
