// Copyright 2018 The klaytn Authors
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

package statedb

import (
	"fmt"
	"reflect"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////
// Additional member functions of Iterator (defined in iterator.go)
////////////////////////////////////////////////////////////////////////////////
func (it *Iterator) NextAny() bool {
	if it.nodeIt.Next(true) {
		return true
	}
	it.Err = it.nodeIt.Error()
	return false
}

////////////////////////////////////////////////////////////////////////////////
// Additional member functions of nodeIterator (defined in iterator.go)
////////////////////////////////////////////////////////////////////////////////
func (it *nodeIterator) GetType() string {
	if len(it.stack) == 0 {
		return ""
	}
	return reflect.TypeOf(it.stack[len(it.stack)-1].node).String()
}

func (it *nodeIterator) GetKeyNibbles() (key, key_nibbles string) {
	k := it.path
	k = k[:len(k)-(len(k)&1)]

	for i, n := range it.path {
		if i == len(it.path)-1 && n == 16 {
			key_nibbles += "T"
		} else {
			key_nibbles += indices[n]
		}
	}

	return string(hexToKeybytes(k)), key_nibbles
}

////////////////////////////////////////////////////////////////////////////////
// NodeIntMap
//
// Stores a mapping between node* and int
// This is required to make an integer ID of a node object for id in vis.js.
////////////////////////////////////////////////////////////////////////////////
type NodeIntMap struct {
	hashMap map[*node]int
	counter int
}

func NewHashIntMap() *NodeIntMap {
	return &NodeIntMap{
		hashMap: map[*node]int{},
		counter: 1,
	}
}

func (m *NodeIntMap) Get(h *node) int {
	if _, ok := m.hashMap[h]; !ok {
		m.hashMap[h] = m.counter
		m.counter++
	}

	return m.hashMap[h]
}

////////////////////////////////////////////////////////////////////////////////
// VisNode
//
// Describes a node object in vis.js.
////////////////////////////////////////////////////////////////////////////////
type VisNode struct {
	id       int
	label    string
	level    int
	addr     *node
	str      string
	typename string
}

func (v *VisNode) String() string {
	return fmt.Sprintf("{id:%d, label:'%s', level:%d, addr:'%p', typename:'%s', x:%d, y:%d}",
		v.id, v.label, v.level, v.addr, v.typename, v.id, v.id)
}

func SerializeNodes(nodes []VisNode) (ret string) {
	for _, n := range nodes {
		ret += fmt.Sprintf("%s, \n", n.String())
	}

	return
}

////////////////////////////////////////////////////////////////////////////////
// VisEdge
//
// Describes an edge object in vis.js.
////////////////////////////////////////////////////////////////////////////////
type VisEdge struct {
	from  int
	to    int
	label string
}

func (e *VisEdge) String() string {
	return fmt.Sprintf("{from:%d, to:%d, label:'%s'}", e.from, e.to, e.label)
}

func SerializeEdges(edges []VisEdge) (ret string) {
	for _, e := range edges {
		ret += fmt.Sprintf("%s, \n", e.String())
	}

	return
}

////////////////////////////////////////////////////////////////////////////////
// TestPrintTrie
//
// You can execute only this test by `go test -run TestPrintTrie`
////////////////////////////////////////////////////////////////////////////////
func TestPrintTrie(t *testing.T) {
	trie := newEmptyTrie()
	vals := []struct{ k, v string }{
		//{"klaytn", "wookiedoo"},
		//{"horse", "stallion"},
		//{"shaman", "horse"},
		//{"doge", "coin"},
		//{"dog", "puppy"},
		{"do", "verb"},
		{"dok", "puppyuyyy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
		{"barb", "ba"},
		{"bard", "bc"},
		{"bars", "bb"},
		{"bar", "b"},
		{"fab", "z"},
		{"food", "ab"},
		{"foos", "aa"},
		{"foo", "a"},
		{"aardvark", "c"},
		//{"bar", "b"},
		//{"barb", "bd"},
		//{"bars", "be"},
		//{"fab", "z"},
		//{"foo", "a"},
		//{"foos", "aa"},
		//{"food", "ab"},
		{"jars", "d"},
	}
	all := make(map[string]string)
	for _, val := range vals {
		all[val.k] = val.v
		trie.Update([]byte(val.k), []byte(val.v))
	}
	trie.Commit(nil)

	nodeIntMap := NewHashIntMap()
	var visNodes []VisNode
	var visEdges []VisEdge

	it := NewIterator(trie.NodeIterator(nil))
	for it.NextAny() {
		nodeIt, _ := it.nodeIt.(*nodeIterator)

		key, key_nibbles := nodeIt.GetKeyNibbles()

		edgeLabel := ""

		myId := nodeIntMap.Get(&nodeIt.stack[len(nodeIt.stack)-1].node)
		pId := 0
		if len(nodeIt.stack) > 1 {
			parent := &nodeIt.stack[len(nodeIt.stack)-2].node
			pId = nodeIntMap.Get(parent)
			switch (*parent).(type) {
			case *fullNode:
				edgeLabel = key_nibbles[len(key_nibbles)-1:]
			default:
			}
		}

		label := string("ROOT")
		if len(key_nibbles) > 0 {
			label = fmt.Sprintf("%s\\n%s", key, key_nibbles)
			if key_nibbles[len(key_nibbles)-1:] == "T" {
				label += fmt.Sprintf("\\nValue:%s", string(nodeIt.LeafBlob()))
			}
		}

		visNodes = append(visNodes, VisNode{
			id:       myId,
			addr:     &nodeIt.stack[len(nodeIt.stack)-1].node,
			str:      nodeIt.stack[len(nodeIt.stack)-1].node.fstring("0"),
			label:    label,
			level:    len(nodeIt.stack),
			typename: nodeIt.GetType(),
		})

		if pId > 0 {
			visEdges = append(visEdges, VisEdge{
				from:  pId,
				to:    myId,
				label: edgeLabel,
			})
		}
	}

	fmt.Printf("var nodes = new vis.DataSet([%s]);\n", SerializeNodes(visNodes))
	fmt.Printf("var edges = new vis.DataSet([%s]);\n", SerializeEdges(visEdges))
}
