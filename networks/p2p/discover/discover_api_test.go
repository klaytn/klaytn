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
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

var (
	rawUrls = []string{
		"kni://a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02@0.0.0.1:3000",
		"kni://bf0d5b99945d5b94bc8430af9ec8019dd3535be586809d04422aa1bacf16b9161cc868c63d11423e30525d8ea41407c28854b1e12058297f8e0ef3c357c8d6c0@0.0.0.2:4000?discport=0",
		"kni://418211f65eb3c14ddab49f690ed291f7681bff08742889a518ca58d2176b4afaa8fdaff1cd45957d7a7c14965d1ffd811bd677e70dc2e309cf8baacf73eb2bdc@0.0.0.3:5000?discport=32323",
		"kni://7ea7fae31c0d1a0a3208607f948b7b8d2eb5c410fb07d27bbcad8359fef1e78ec10e3f9132feaa69b0438ac70d5a3a13af041e47089bc43d90f68ba5fc72a23c@0.0.0.4:6000",
		"kni://e92f2fcd87ad42b9f6e55aebc2e35c73e5443fc3128a82dc11cc84acf445a14beee202140c5004d73c981d8bf423cb05714e9f7b1cabfd2720e34ed22e819dc3@0.0.0.5:7000",
		"kni://01a63da798a08b6322ae59d398628addebd87eb582a801976d7d1b2a1143574a15adc4a8ae08d1286e8cb36ef6d904e25ae932d3cb4e695944b48afb58981550@0.0.0.6:8000?discport=0",
	}
)

type mockDiscoveryStorage struct {
	data []*Node
}

func (mds *mockDiscoveryStorage) name() string                        { return "mockDiscoveryStorage" }
func (mds *mockDiscoveryStorage) setTable(t *Table)                   {}
func (mds *mockDiscoveryStorage) setTargetNodeType(tType NodeType)    {}
func (mds *mockDiscoveryStorage) init()                               {}
func (mds *mockDiscoveryStorage) add(n *Node)                         { mds.data = append(mds.data, n) }
func (mds *mockDiscoveryStorage) delete(n *Node)                      { mds.data = deleteNode(mds.data, n) }
func (mds *mockDiscoveryStorage) len() (n int)                        { return 0 }
func (mds *mockDiscoveryStorage) nodeAll() []*Node                    { return mds.data }
func (mds *mockDiscoveryStorage) readRandomNodes(buf []*Node) (n int) { return 0 }
func (mds *mockDiscoveryStorage) closest(target common.Hash, nresults int) *nodesByDistance {
	return nil
}
func (mds *mockDiscoveryStorage) stuff(nodes []*Node) {}
func (mds *mockDiscoveryStorage) copyBondedNodes()    {}
func (mds *mockDiscoveryStorage) lookup(targetID NodeID, refreshIfEmpty bool, targetType NodeType) []*Node {
	return nil
}
func (mds *mockDiscoveryStorage) getNodes(max int) []*Node { return nil }
func (mds *mockDiscoveryStorage) doRevalidate()            {}
func (mds *mockDiscoveryStorage) doRefresh()               {}

func (mds *mockDiscoveryStorage) isAuthorized(id NodeID) bool    { return true }
func (mds *mockDiscoveryStorage) getBucketEntries() []*Node      { return mds.data }
func (mds *mockDiscoveryStorage) getAuthorizedNodes() []*Node    { return nil }
func (mds *mockDiscoveryStorage) putAuthorizedNode(node *Node)   {}
func (mds *mockDiscoveryStorage) deleteAuthorizedNode(id NodeID) {}

func createTestData() *mockDiscoveryStorage {
	var d []*Node
	for _, url := range rawUrls {
		node, _ := ParseNode(url)
		d = append(d, node)
	}
	mds := &mockDiscoveryStorage{
		data: d,
	}
	return mds
}

func createTestTable() (*Table, error) {
	id := MustHexID("904afadc3587310adad21971c28a9213d03b04c4bd5ec73f7609861eb57b60b26809a52e72cbc6af9dfdd780a6d6629898ef197149f840e03f237660bc0dbb48")
	tmpDB, _ := newNodeDB("", Version, id)
	tab := &Table{
		db:          tmpDB,
		storages:    make(map[NodeType]discoverStorage),
		localLogger: log.NewModuleLogger(log.NetworksP2PDiscover).NewWith("Discover", "Test"),
	}
	storage := createTestData()
	tab.addStorage(NodeTypeUnknown, storage)
	for _, node := range storage.data {
		err := tab.db.updateNode(node)
		if err != nil {
			return nil, err
		}
	}
	return tab, nil
}

func TestTable_CreateUpdateNodeOnDB(t *testing.T) {
	tab, _ := createTestTable()
	n, _ := ParseNode("kni://31a03810a315b200053ea5905d4e0a32a50eafcf2613e90e17ee877652dbb2aba1c725016ebe44d48529c8755203046a927fb01b97d758c009d248c4d7dfdd8c@0.0.0.0:32323?discport=0")
	err := tab.CreateUpdateNodeOnDB(n)
	if err != nil {
		t.Error("CreateUpdateNodeOnDB is failed", "err", err)
	}
	id, _ := HexID("31a03810a315b200053ea5905d4e0a32a50eafcf2613e90e17ee877652dbb2aba1c725016ebe44d48529c8755203046a927fb01b97d758c009d248c4d7dfdd8c")
	n2 := tab.db.node(id)
	if !n.CompareNode(n2) {
		t.Error("nodes are different", "inserted:", n, "retrieved", n2)
	}
}

func TestTable_CreateUpdateNodeOnTable(t *testing.T) {
	tab, _ := createTestTable()
	n, _ := ParseNode("kni://31a03810a315b200053ea5905d4e0a32a50eafcf2613e90e17ee877652dbb2aba1c725016ebe44d48529c8755203046a927fb01b97d758c009d248c4d7dfdd8c@0.0.0.0:32323?discport=0")
	err := tab.CreateUpdateNodeOnTable(n)
	if err != nil {
		t.Error("CreateUpdateNodeOnTable is failed", "err", err)
	}
	entries := tab.GetBucketEntries()
	replacements := tab.GetReplacements()
	if !isIn(n, append(entries, replacements...)) {
		t.Errorf("the node does not exist. nodeid: %v", n.ID)
	}
}

func TestTable_GetNodeFromDB(t *testing.T) {
	tab, _ := createTestTable()
	id0 := MustHexID("a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02")
	id1 := MustHexID("bf0d5b99945d5b94bc8430af9ec8019dd3535be586809d04422aa1bacf16b9161cc868c63d11423e30525d8ea41407c28854b1e12058297f8e0ef3c357c8d6c0")
	id2 := MustHexID("418211f65eb3c14ddab49f690ed291f7681bff08742889a518ca58d2176b4afaa8fdaff1cd45957d7a7c14965d1ffd811bd677e70dc2e309cf8baacf73eb2bdc")
	id3 := MustHexID("7ea7fae31c0d1a0a3208607f948b7b8d2eb5c410fb07d27bbcad8359fef1e78ec10e3f9132feaa69b0438ac70d5a3a13af041e47089bc43d90f68ba5fc72a23c")
	id4 := MustHexID("e92f2fcd87ad42b9f6e55aebc2e35c73e5443fc3128a82dc11cc84acf445a14beee202140c5004d73c981d8bf423cb05714e9f7b1cabfd2720e34ed22e819dc3")
	id5 := MustHexID("01a63da798a08b6322ae59d398628addebd87eb582a801976d7d1b2a1143574a15adc4a8ae08d1286e8cb36ef6d904e25ae932d3cb4e695944b48afb58981550")
	ids := []NodeID{id0, id1, id2, id3, id4, id5}
	for idx, id := range ids {
		node1, err := tab.GetNodeFromDB(id)
		if err != nil {
			t.Error("get node is failed.", "err:", err)
		}
		node2, _ := ParseNode(rawUrls[idx])
		if !node1.CompareNode(node2) {
			t.Error("node1 is different from node2.", "node1:", node1, "node2:", node2)
		}
	}
}

func TestTable_GetBucketEntries(t *testing.T) {
	tab, _ := createTestTable()
	entries := tab.GetBucketEntries()
	if storage, ok := tab.storages[NodeTypeUnknown].(*mockDiscoveryStorage); ok {
		if len(storage.data) != len(entries) {
			t.Error("the length of actual data is different from the result of GetBucketEntries")
		}
		for _, node := range storage.data {
			existed := false
			for _, node2 := range entries {
				if node.CompareNode(node2) {
					existed = true
					break
				}
			}
			if !existed {
				t.Error("node does not existed", "nodeid", node.ID)
			}
		}
		for _, node := range storage.data {
			faultyNode, _ := ParseNode("kni://71cd148865fd6dd9658a233aef8e7d844d279f4d04c9783c91d63c664109709609900ee3bc6e6b8ca7f4bff555be5331c10801e1e6e6713ebfd339674ccc377b@0.0.0.0:32323?discport=0")
			if node.CompareNode(faultyNode) {
				t.Error("the node existed although it is not updated", "nodeid", faultyNode.ID)
			}
		}
	}

}

func TestTable_DeleteNodeFromDB(t *testing.T) {
	testData1 := MustParseNode(rawUrls[0])
	testData2 := MustParseNode(rawUrls[1])
	tab, _ := createTestTable()
	if n := tab.db.node(testData1.ID); n == nil {
		t.Error("node does not exist", "nodeid", testData1.ID)
	}
	if n := tab.db.node(testData2.ID); n == nil {
		t.Error("node does not exist", "nodeid", testData2.ID)
	}
	err := tab.DeleteNodeFromDB(testData1)
	if err != nil {
		t.Error("delete node from table is failed", "err", err)
	}
	err = tab.DeleteNodeFromDB(testData2)
	if err != nil {
		t.Error("delete node from table is failed", "err", err)
	}
	if n := tab.db.node(testData1.ID); n != nil {
		t.Error("node is not removed", "nodeid", n.ID)
	}
	if n := tab.db.node(testData2.ID); n != nil {
		t.Error("node is not removed", "nodeid", n.ID)
	}
}

func TestTable_DeleteNodeFromTable(t *testing.T) {
	tab, _ := createTestTable()
	storage, _ := tab.storages[NodeTypeUnknown].(*mockDiscoveryStorage)
	node0 := storage.data[0]
	node1 := storage.data[1]
	err := tab.DeleteNodeFromTable(node0)
	if err != nil {
		t.Error("delete node from table is failed", "err", err)
	}
	err = tab.DeleteNodeFromTable(node1)
	if err != nil {
		t.Error("delete node from table is failed", "err", err)
	}
	if isIn(node0, storage.data) {
		t.Error("the node is not removed", "nodeid", node0.ID)
	}
	if isIn(node1, storage.data) {
		t.Error("the node is not removed", "nodeid", node1.ID)
	}
}
