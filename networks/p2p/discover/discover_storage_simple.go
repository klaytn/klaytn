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

import (
	"math/rand"
	"sync"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
)

type simpleStorage struct {
	tab         *Table
	targetType  NodeType
	nodes       []*Node
	noDiscover  bool // if noDiscover is true, don't lookup new node.
	nodesMutex  sync.Mutex
	max         int
	rand        *rand.Rand
	localLogger log.Logger

	lock               sync.RWMutex
	hasAuthorizedNodes bool
	authorizedNodes    map[NodeID]*Node
}

func NewSimpleStorage(nodeType NodeType, noDiscover bool, max int, authorizedNodes []*Node) *simpleStorage {
	storage := &simpleStorage{targetType: nodeType, noDiscover: noDiscover, max: max}
	storage.authorizedNodes = make(map[NodeID]*Node)
	for _, node := range authorizedNodes {
		if node.NType == nodeType {
			storage.putAuthorizedNode(node)
		}
	}
	return storage
}

func (s *simpleStorage) init() {
	// TODO
	s.localLogger = logger.NewWith("Discover", "Simple")
	now := time.Now().UnixNano()
	s.rand = rand.New(rand.NewSource(now))
}

func (s *simpleStorage) lookup(targetID NodeID, refreshIfEmpty bool, targetType NodeType) []*Node {
	// check exist alive bn
	var seeds []*Node
	s.nodesMutex.Lock()
	for _, n := range s.nodes {
		if n.NType == NodeTypeBN {
			seeds = append(seeds, n)
		}
	}
	s.nodesMutex.Unlock()

	if len(seeds) == 0 {
		seeds = append([]*Node{}, s.tab.nursery...)
		seeds = s.tab.bondall(seeds)
		for _, n := range seeds {
			s.add(n)
		}
	}
	s.localLogger.Debug("lookup", "StorageName", s.name(), "targetId", targetID, "targetType", nodeTypeName(targetType))
	return s.tab.findNewNode(&nodesByDistance{entries: seeds}, targetID, targetType, false, s.max)
}

func (s *simpleStorage) shuffle(vals []*Node) []*Node {
	if len(vals) == 0 {
		return vals
	}
	ret := make([]*Node, len(vals))
	perm := s.rand.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func (s *simpleStorage) getNodes(max int) []*Node {
	s.nodesMutex.Lock()
	nodes := s.shuffle(s.nodes)

	var ret []*Node
	for _, nd := range nodes {
		if nd.NType == s.targetType {
			ret = append(ret, nd)
		}
		if len(ret) >= max {
			break
		}
	}
	s.nodesMutex.Unlock()

	if len(ret) < max {
		ret = s.lookup(NodeID{}, true, s.targetType)
	}

	if len(ret) < max {
		return ret
	} else {
		return ret[:max]
	}
}

func (s *simpleStorage) doRevalidate() {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()

	if len(s.nodes) == 0 {
		return
	}

	oldest := s.nodes[len(s.nodes)-1]

	holdingTime := s.tab.db.bondTime(oldest.ID).Add(10 * time.Second) // TODO-Klaytn-Node Make sleep time as configurable
	if time.Now().Before(holdingTime) {
		return
	}

	err := s.tab.ping(oldest.ID, oldest.addr())

	if err != nil {
		s.localLogger.Info("Removed the node without any response", "StorageName", s.name(), "NodeID", oldest.ID, "NodeType", nodeTypeName(oldest.NType))
		s.deleteWithoutLock(oldest)
		return
	}
	copy(s.nodes[1:], s.nodes[:len(s.nodes)-1])
	s.nodes[0] = oldest
}

func (s *simpleStorage) setTargetNodeType(tType NodeType) {
	s.targetType = tType
}

func (s *simpleStorage) doRefresh() {
	if s.noDiscover {
		return
	}
	s.lookup(s.tab.self.ID, false, s.targetType)
}

func (s *simpleStorage) nodeAll() []*Node {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()
	return s.nodes
}

func (s *simpleStorage) len() (n int) {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()
	return len(s.nodes)
}

func (s *simpleStorage) copyBondedNodes() {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()
	for _, n := range s.nodes {
		s.tab.db.updateNode(n)
	}
}

func (s *simpleStorage) getBucketEntries() []*Node {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()

	var ret []*Node
	for _, n := range s.nodes {
		ret = append(ret, n)
	}
	return ret
}

// The caller must not hold tab.mutex.
func (s *simpleStorage) stuff(nodes []*Node) {
	panic("implement me")
}

// The caller must hold s.nodesMutex.
func (s *simpleStorage) delete(n *Node) {
	s.deleteWithLock(n)
}

func (s *simpleStorage) deleteWithLock(n *Node) {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()
	s.deleteWithoutLock(n)
}

func (s *simpleStorage) deleteWithoutLock(n *Node) {
	s.nodes = deleteNode(s.nodes, n)

	s.tab.db.deleteNode(n.ID)
	if netutil.IsLAN(n.IP) {
		return
	}
	s.tab.ips.Remove(n.IP)
}

func (s *simpleStorage) closest(target common.Hash, nresults int) *nodesByDistance {
	s.nodesMutex.Lock()
	defer s.nodesMutex.Unlock()
	// TODO-Klaytn-Node nodesByDistance is not suitable for SimpleStorage. Because there is no concept for distance
	// in the SimpleStorage. Change it
	cNodes := &nodesByDistance{target: target}
	nodes := s.shuffle(s.nodes)
	if len(nodes) > s.max {
		cNodes.entries = nodes[:s.max]
	} else {
		cNodes.entries = nodes
	}
	return cNodes
}

func (s *simpleStorage) setTable(t *Table) {
	s.tab = t
}

func (s *simpleStorage) readRandomNodes(buf []*Node) (n int) {
	panic("implement me")
}

func (s *simpleStorage) add(n *Node) {
	s.nodesMutex.Lock()
	s.bumpOrAdd(n)
	s.nodesMutex.Unlock()
}

// The caller must hold s.nodesMutex.
func (s *simpleStorage) bumpOrAdd(n *Node) bool {
	if s.bump(n) {
		s.localLogger.Trace("Add(Bumped)", "StorageName", s.name(), "node", n)
		return true
	}

	s.localLogger.Trace("Add(New)", "StorageName", s.name(), "node", n)
	s.nodes, _ = pushNode(s.nodes, n, math.MaxInt64) // TODO-Klaytn-Node Change Max value for more reasonable one.
	n.addedAt = time.Now()
	if s.tab.nodeAddedHook != nil {
		s.tab.nodeAddedHook(n)
	}
	return true
}

// The caller must hold s.nodesMutex.
func (s *simpleStorage) bump(n *Node) bool {
	for i := range s.nodes {
		if s.nodes[i].ID == n.ID {
			// move it to the front
			copy(s.nodes[1:], s.nodes[:i])
			s.nodes[0] = n
			return true
		}
	}
	return false
}

func (s *simpleStorage) name() string {
	return nodeTypeName(s.targetType)
}

func (s *simpleStorage) isAuthorized(id NodeID) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if !s.hasAuthorizedNodes {
		return true
	}
	for _, n := range s.authorizedNodes {
		if n.ID == id {
			return true
		}
	}
	return false
}

func (s *simpleStorage) getAuthorizedNodes() []*Node {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var ret []*Node
	for _, val := range s.authorizedNodes {
		ret = append(ret, val)
	}
	return ret
}

func (s *simpleStorage) putAuthorizedNode(node *Node) {
	if node.NType != s.targetType {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.authorizedNodes[node.ID] = node
	s.hasAuthorizedNodes = true
}

func (s *simpleStorage) deleteAuthorizedNode(id NodeID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.authorizedNodes[id]; ok {
		delete(s.authorizedNodes, id)
		s.hasAuthorizedNodes = len(s.authorizedNodes) != 0
	} else {
		logger.Debug("No node to be removed", "nodeid", id)
	}
}
