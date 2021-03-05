// Modifications Copyright 2019 The klaytn Authors
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
	"crypto/rand"
	"net"
	"sync"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
)

const (
	tableIPLimit, tableSubnet   = 10, 24
	bucketIPLimit, bucketSubnet = 2, 24 // at most 2 addresses from the same /24
	// We keep buckets for the upper 1/15 of distances because
	// it's very unlikely we'll ever encounter a node that's closer.
	hashBits          = len(common.Hash{}) * 8
	nBuckets          = hashBits / 15       // Number of buckets
	bucketMinDistance = hashBits - nBuckets // Log distance of closest bucket

	seedMinTableTime = 5 * time.Minute
)

type KademliaStorage struct {
	targetType  NodeType
	tab         *Table
	buckets     [nBuckets]*bucket
	bucketsMu   sync.Mutex
	ips         netutil.DistinctNetSet
	noDiscover  bool // if noDiscover is true, doesn't lookup new node.
	localLogger log.Logger
}

func (s *KademliaStorage) init() {
	s.localLogger = logger.NewWith("Discover", "Kademlia")
	s.ips = netutil.DistinctNetSet{Subnet: tableSubnet, Limit: tableIPLimit}
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()
	for i := range s.buckets {
		s.buckets[i] = &bucket{
			ips: netutil.DistinctNetSet{Subnet: bucketSubnet, Limit: bucketIPLimit},
		}
	}
}

func (s *KademliaStorage) lookup(targetID NodeID, refreshIfEmpty bool, targetType NodeType) []*Node {
	s.localLogger.Debug("lookup start", "StorageName", s.name(), "targetID", targetID,
		"targetNodeType", nodeTypeName(targetType), "refreshIfEmpty", refreshIfEmpty)
	var (
		target = crypto.Keccak256Hash(targetID[:])
		result *nodesByDistance
	)

	for {
		// generate initial result set
		result = s.closest(target, bucketSize)
		if len(result.entries) > 0 || !refreshIfEmpty {
			break
		}
		// The result set is empty, all nodes were dropped, refresh.
		// We actually wait for the refresh to complete here. The very
		// first query will hit this case and run the bootstrapping
		// logic.
		<-s.tab.refresh()
		refreshIfEmpty = false
	}
	return s.tab.findNewNode(result, targetID, targetType, true, bucketSize)
}

func (s *KademliaStorage) getNodes(max int) []*Node {
	nbd := s.closest(crypto.Keccak256Hash(s.tab.self.ID[:]), max)
	var ret []*Node
	for _, nd := range nbd.entries {
		if nd.NType == s.targetType {
			ret = append(ret, nd)
		}
	}
	return ret
}

func (s *KademliaStorage) doRevalidate() {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()

	last, bi := s.nodeToRevalidate()
	if last == nil {
		// No non-empty bucket found.
		return
	}

	holdingTime := s.tab.db.bondTime(last.ID).Add(10 * time.Second)
	if time.Now().Before(holdingTime) {
		s.localLogger.Debug("skip revalidate", "StorageName", s.name())
		return
	}

	// Ping the selected node and wait for a pong.
	err := s.tab.ping(last.ID, last.addr())
	b := s.bucketByIdx(bi)
	if err == nil {
		// The node responded, move it to the front.
		s.localLogger.Debug("Revalidated node", "StorageName", s.name(), "bucketIdx", bi, "nodeId", last.ID)
		b.bump(last)
		return
	}
	// No reply received, pick a replacement or delete the node if there aren't
	// any replacements.
	if r := s.replace(b, last); r != nil {
		s.localLogger.Debug("Replaced the node without any response", "StorageName", s.name(), "bucketIdx", bi, "nodeId", last.ID, "ip", last.IP, "r", r.ID, "rip", r.IP)
	} else {
		s.localLogger.Debug("Removed the node without any response", "StorageName", s.name(), "bucketIdx", bi, "nodeId", last.ID, "ip", last.IP)
	}
}

func (s *KademliaStorage) setTargetNodeType(tType NodeType) {
	s.targetType = tType
}

func (s *KademliaStorage) doRefresh() {
	if s.noDiscover {
		return
	}

	// Run self lookup to discover new neighbor nodes.
	s.lookup(s.tab.self.ID, false, s.targetType)

	// The Kademlia paper specifies that the bucket refresh should
	// perform a lookup in the least recently used bucket. We cannot
	// adhere to this because the findnode target is a 512bit value
	// (not hash-sized) and it is not easily possible to generate a
	// sha3 preimage that falls into a chosen bucket.
	// We perform a few lookups with a random target instead.
	for i := 0; i < 3; i++ {
		var target NodeID
		rand.Read(target[:])
		s.lookup(target, false, s.targetType)
	}
}

func (s *KademliaStorage) nodeAll() (nodes []*Node) {
	for _, b := range s.buckets {
		nodes = append(nodes, b.entries...)
	}
	return nodes
}

// The caller must hold s.bucketMu
func (s *KademliaStorage) bucketByIdx(bi int) *bucket {
	// TODO-Klaytn-Node range check
	return s.buckets[bi]
}

// The caller must hold s.bucketMu
func (s *KademliaStorage) nodeToRevalidate() (n *Node, bi int) {
	s.tab.randMu.Lock()
	defer s.tab.randMu.Unlock()

	for _, bi = range s.tab.rand.Perm(len(s.buckets)) { // TODO
		b := s.buckets[bi]
		if len(b.entries) > 0 {
			last := b.entries[len(b.entries)-1]
			return last, bi
		}
	}
	return nil, 0
}

func (s *KademliaStorage) copyBondedNodes() {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()
	now := time.Now()
	for _, b := range &s.buckets {
		for _, n := range b.entries {
			if now.Sub(n.addedAt) >= seedMinTableTime {
				s.tab.db.updateNode(n)
			}
		}
	}
}

func (s *KademliaStorage) len() (n int) {
	for _, b := range &s.buckets {
		n += len(b.entries)
	}
	return n
}

func (s *KademliaStorage) getReplacements() []*Node {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()
	var nodes []*Node
	for i := 0; i < nBuckets; i++ {
		nodes = append(nodes, s.buckets[i].replacements...)
	}
	return nodes
}

func (s *KademliaStorage) getBucketEntries() []*Node {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()

	var nodes []*Node
	for i := 0; i < nBuckets; i++ {
		nodes = append(nodes, s.buckets[i].entries...)
	}
	return nodes
}

func (s *KademliaStorage) stuff(nodes []*Node) {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()

	for _, n := range nodes {
		if n.ID == s.tab.self.ID {
			continue // don't add self
		}
		b := s.bucket(n.sha)
		if len(b.entries) < bucketSize {
			s.bumpOrAdd(b, n)
		}

	}
}

// replace removes n from the replacement list and replaces 'last' with it if it is the
// last entry in the bucket. If 'last' isn't the last entry, it has either been replaced
// with someone else or became active.
func (s *KademliaStorage) replace(b *bucket, last *Node) *Node {
	if len(b.entries) == 0 || b.entries[len(b.entries)-1].ID != last.ID {
		// Entry has moved, don't replace it.
		return nil
	}
	// Still the last entry.
	if len(b.replacements) == 0 {
		s.deleteInBucket(b, last)
		return nil
	}
	s.tab.randMu.Lock()
	r := b.replacements[s.tab.rand.Intn(len(b.replacements))]
	s.tab.randMu.Unlock()
	b.replacements = deleteNode(b.replacements, r)
	b.entries[len(b.entries)-1] = r
	s.removeIP(b, last.IP)
	return r
}

func (s *KademliaStorage) delete(n *Node) {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()
	s.deleteInBucket(s.bucket(n.sha), n)
}

// The caller must hold tab.mutex.
func (s *KademliaStorage) deleteInBucket(b *bucket, n *Node) {
	b.entries = deleteNode(b.entries, n)
	s.removeIP(b, n.IP) // TODO-Klaytn-Node Does the IP is not lock?
}

// closest returns the n nodes in the table that are closest to the
// given id. The caller must hold s.bucketMu.
func (s *KademliaStorage) closest(target common.Hash, nresults int) *nodesByDistance {
	// This is a very wasteful way to find the closest nodes but
	// obviously correct. I believe that tree-based buckets would make
	// this easier to implement efficiently.
	// TODO-Klaytn-Node more efficient ways to obtain the closest nodes could be considered.
	close := &nodesByDistance{target: target}
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()

	for _, b := range &s.buckets {
		for _, n := range b.entries {
			close.push(n, nresults)
		}
	}
	return close
}

func (s *KademliaStorage) setTable(t *Table) {
	s.tab = t
}

func (s *KademliaStorage) add(n *Node) {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()
	b := s.bucket(n.sha)
	if !s.bumpOrAdd(b, n) {
		// Node is not in table. Add it to the replacement list.
		s.addReplacement(b, n)
	}
}

func (s *KademliaStorage) readRandomNodes(buf []*Node) (n int) {
	s.bucketsMu.Lock()
	defer s.bucketsMu.Unlock()

	// Find all non-empty buckets and get a fresh slice of their entries.
	var buckets [][]*Node
	for _, b := range &s.buckets {
		if len(b.entries) > 0 {
			buckets = append(buckets, b.entries[:])
		}
	}
	if len(buckets) == 0 {
		return 0
	}
	// Shuffle the buckets.
	s.tab.randMu.Lock()
	for i := len(buckets) - 1; i > 0; i-- {
		j := s.tab.rand.Intn(len(buckets))
		buckets[i], buckets[j] = buckets[j], buckets[i]
	}
	s.tab.randMu.Unlock()
	// Move head of each bucket into buf, removing buckets that become empty.
	var i, j int
	for ; i < len(buf); i, j = i+1, (j+1)%len(buckets) {
		b := buckets[j]
		buf[i] = &(*b[0])
		buckets[j] = b[1:]
		if len(b) == 1 {
			buckets = append(buckets[:j], buckets[j+1:]...)
		}
		if len(buckets) == 0 {
			break
		}
	}
	return i + 1
}

// The caller must hold s.bucketMu
func (s *KademliaStorage) bucket(sha common.Hash) *bucket {
	d := logdist(s.tab.self.sha, sha)
	if d <= bucketMinDistance {
		return s.buckets[0]
	}
	return s.buckets[d-bucketMinDistance-1]
}

// bumpOrAdd moves n to the front of the bucket entry list or adds it if the list isn't
// full. The return value is true if n is in the bucket.
// The caller must hold s.bucketMu
func (s *KademliaStorage) bumpOrAdd(b *bucket, n *Node) bool {
	if b.bump(n) {
		s.localLogger.Trace("Add(Bumped)", "StorageName", s.name(), "node", n)
		return true
	}
	if len(b.entries) >= bucketSize || !s.addIP(b, n.IP) {
		s.localLogger.Debug("Add(New) -Exceed Max", "StorageName", s.name(), "node", n)
		return false
	}
	s.localLogger.Trace("Add(New)", "StorageName", s.name(), "node", n)
	b.entries, _ = pushNode(b.entries, n, bucketSize)
	b.replacements = deleteNode(b.replacements, n)
	n.addedAt = time.Now()
	if s.tab.nodeAddedHook != nil {
		s.tab.nodeAddedHook(n)
	}
	return true
}

// The caller must hold s.bucketMu.
func (s *KademliaStorage) addReplacement(b *bucket, n *Node) {
	for _, e := range b.replacements {
		if e.ID == n.ID {
			return // already in list
		}
	}
	if !s.addIP(b, n.IP) {
		return
	}
	var removed *Node
	b.replacements, removed = pushNode(b.replacements, n, maxReplacements)
	if removed != nil {
		s.removeIP(b, removed.IP)
	}
}

func (s *KademliaStorage) addIP(b *bucket, ip net.IP) bool {
	if netutil.IsLAN(ip) {
		return true
	}
	if !s.ips.Add(ip) {
		s.localLogger.Debug("IP exceeds table limit", "StorageName", s.name(), "ip", ip)
		return false
	}
	if !b.ips.Add(ip) {
		s.localLogger.Debug("IP exceeds bucket limit", "StorageName", s.name(), "ip", ip)
		s.ips.Remove(ip)
		return false
	}
	return true
}

func (s *KademliaStorage) removeIP(b *bucket, ip net.IP) {
	if netutil.IsLAN(ip) {
		return
	}
	s.ips.Remove(ip)
	b.ips.Remove(ip)
}

func (s *KademliaStorage) name() string {
	return nodeTypeName(s.targetType)
}

// bucket contains nodes, ordered by their last activity. the entry
// that was most recently active is the first element in entries.
type bucket struct {
	entries      []*Node // live entries, sorted by time of last contact
	replacements []*Node // recently seen nodes to be used if revalidation fails
	ips          netutil.DistinctNetSet
}

// bump moves the given node to the front of the bucket entry list
// if it is contained in that list.
// caller
func (b *bucket) bump(n *Node) bool {
	for i := range b.entries {
		if b.entries[i].ID == n.ID {
			// move it to the front
			copy(b.entries[1:], b.entries[:i])
			b.entries[0] = n
			return true
		}
	}
	return false
}

func (s *KademliaStorage) isAuthorized(id NodeID) bool { return true }
func (s *KademliaStorage) getAuthorizedNodes() []*Node { return nil }
func (s *KademliaStorage) putAuthorizedNode(*Node)     {}
func (s *KademliaStorage) deleteAuthorizedNode(NodeID) {}
