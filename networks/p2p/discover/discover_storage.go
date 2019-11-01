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
	"github.com/klaytn/klaytn/common"
)

type discoverStorage interface {
	name() string
	setTable(t *Table)
	setTargetNodeType(tType NodeType)
	init()

	add(n *Node)
	delete(n *Node)
	len() (n int)
	nodeAll() []*Node
	readRandomNodes(buf []*Node) (n int)
	closest(target common.Hash, nresults int) *nodesByDistance
	stuff(nodes []*Node)
	copyBondedNodes()

	lookup(targetID NodeID, refreshIfEmpty bool, targetType NodeType) []*Node
	getNodes(max int) []*Node
	doRevalidate()
	doRefresh()

	isAuthorized(id NodeID) bool

	// API
	getBucketEntries() []*Node
	getAuthorizedNodes() []*Node
	putAuthorizedNode(node *Node)
	deleteAuthorizedNode(id NodeID)
}

// pushNode adds n to the front of list, keeping at most max items.
func pushNode(list []*Node, n *Node, max int) ([]*Node, *Node) {
	if len(list) < max {
		list = append(list, nil)
	}
	removed := list[len(list)-1]
	copy(list[1:], list)
	list[0] = n
	return list, removed
}

// deleteNode removes n from list.
func deleteNode(list []*Node, n *Node) []*Node {
	for i := range list {
		if list[i].ID == n.ID {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// nodeTypeName converts NodeType to string.
func nodeTypeName(nt NodeType) string { // TODO-Klaytn-Node Consolidate p2p.NodeType and common.ConnType
	switch nt {
	case NodeTypeCN:
		return "CN"
	case NodeTypePN:
		return "PN"
	case NodeTypeEN:
		return "EN"
	case NodeTypeBN:
		return "BN"
	default:
		return "Unknown Node Type"
	}
}
