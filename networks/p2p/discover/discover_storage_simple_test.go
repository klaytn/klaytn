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
	"crypto/rand"
	rand2 "math/rand"
	"net"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/pkg/errors"
)

var (
	testData = [5][]*Node{
		NodeTypeUnknown: {
			MustParseNode("kni://2d2d43be39c40e1b104952cc351127fb3783b66ca065ba4b8c46f6e73e603e511203e399154c5b96c8ca13f8dd9086f6d29c74867a3b7ea6bd4f0205b25522b6@0.0.0.1:32323?discport=32323"),
			MustParseNode("kni://e2fc0988b6286fd15a9c208bdf283fb456575357911d618e2f47326ce534db1f94c3a1edc6cb7a399797b35ad7dd82a2647ba9ec43b33948302cefb9edd2c9b2@0.0.0.2:32323?discport=32323"),
			MustParseNode("kni://e898d53588ed46888ace5a6a61c2ca71034ae23aa004b8525a5045a5a51d43dd72b9ca49b346ed0155cc6e2cd143486109a6fe21845a59778012994ea7d9e128@0.0.0.3:32323?discport=32323"),
			MustParseNode("kni://79ce147a6e955cbf43004e19480c4d8139eeb71ec99ee872499e4cb05e37ec711049781702200f958698246475e7bcf898acad261d1950c280abf5090ea00ac0@0.0.0.4:32323?discport=32323"),
			MustParseNode("kni://693d678f8497a0a6019acdcc6388c489a0078387bb0dafa27e36b856765b20abe6a7c6e69b6bb4c0028babf490f7aea5481c2c39f69153c4f670af720ec59f67@0.0.0.5:32323?discport=32323"),
			MustParseNode("kni://850555871dd0eeea187d7ad0f219b4def2a377cee8a6b40dff985b907687c4a36afa337a8499aff157ef41dadba8dcfb202786bb86d0a13371c8a517db3bfade@0.0.0.6:32323?discport=32323"),
			MustParseNode("kni://dda3ce1adf19a1d88fba5fe91d856aa8c5a9fb71892b4cf75ffc19f28a6a38cc8eb3d4cdd83c2013ba8e8112fe44a43de0198e639171c4d5887452f0f7b21712@0.0.0.7:32323?discport=32323"),
			MustParseNode("kni://473e2e359e2de9922eff61269a719004768c96733da823aaabbe6c926540f40156107b7fac931703d9f01a9cfc33652889c98fb62dbffe6a04a9ce2c93bc7512@0.0.0.8:32323?discport=32323"),
			MustParseNode("kni://26a8bd88faadfe401787e99ca7909793cbe388562a0e43286488c2306770a73173ed0a50c5c894ae9e390da21c9091f4edf924d3e4299b303434628198b8a38a@0.0.0.9:32323?discport=32323"),
			MustParseNode("kni://dbc67e54732ca531c8c3b2175858d712d1a1cac4a59ad6b639c6fe72bc261916c2236a38227502975f10e5353e4554dc50daaede6f8c9203beb36438af356d99@0.0.0.10:32323?discport=32323"),
		},
		NodeTypeCN: {
			MustParseNode("kni://385730c3ca13a8a469fe8530216a3f17f1befe2f4a32394b58a63744f24170a135e1b9d22f34dcc6f52c98aa4248f8978a6ae9e22dc4a5869f8bf3e4154902fd@0.0.0.20:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://9114a018b91ed64c5e8b3e41494079d447bef4de19b3b41d36e46241cd06a69799a2fc6ed38a38f444ba508724efdea10ea4f3d95d331bb797c8601d0074fe5b@0.0.0.21:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://d3bc248af75576d88a8f8077e124f2d6d3980a40e0eeba093a06ba43ccd32b79ad8abb69ac38d349d3b52565033ebe917104b6a37108ce9c0fbdfcb52ad3c2b8@0.0.0.22:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://1914b4f6ee1eb996ee56551d0a95496ef617afeedc264f02db7abb97bee731f21e8fcd2354362a7ca3e043ace4541a54714e0287430228a802b4882154b941ba@0.0.0.23:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://6648a7c1c479c31e7c3410464f624f1a28819c9eb008930bf80a5614e0cf67d1a656024efca51971ac7f6a202355957e0574e6f02383bfeab444327291b335af@0.0.0.24:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://860a1057975bfd40bde03bdb376b33f290af988f378e6a6a2c78d7fd931ba84a16c9f4af193f77d7b642e5faa7f9333b2e980960b10f04cd915ce47837bb109f@0.0.0.25:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://dee90a411b0746300614041dde915979ff3763e7191414b4ee03ab5ab75ceabfd575cb85b1a033731535beb4a06b67c129df6e7ee7558c8775c8f1cf6fd6c289@0.0.0.26:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://7cf836b7f3518e632b939197fb1836547aa6d751ad48b31b5afb06a85404cd6731d2fcb3b123715c9efef996468fe3d3201b86a1a05d0f767c69ac6495d8b49b@0.0.0.27:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://2416d306b78af7de6f1ffec1e047adced5355d00d027442a247c48c34f5aa1693521dc8fe3b08bc50c7ccfecfc42350e67ccf48e834fda8f8ca42fb0125a2d08@0.0.0.28:32323?discport=32323&ntype=cn"),
			MustParseNode("kni://da4bcb9afda40ef3b23c548387c464d656843640d0be14fd64c487e42c009922dcd3a1a2d7b76e86b8266666c826e9b5f462cddd89eab202336aec7762988cf7@0.0.0.29:32323?discport=32323&ntype=cn"),
		},
		NodeTypePN: {
			MustParseNode("kni://1240e51c2afb39b00983a8b2ccefb6dbc69106573172495fb939c093e1dbf8f2465e0fc6e8d92f4a4392dd3e790b4fefe99644e21593bcc71f1ded132e701cf1@0.0.0.30:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://df9adb43f4f4c9d3ec8916658165af95ab152615917ebdd683fcc9e72bf91b43562ff0292eeb06a5beef643106605b99fda64583ad9367129b18ffd1b2150745@0.0.0.31:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://913e6c1fc9b07e9491449e004533f981930897270aeb662a3699d97fe597744c098741f7507339cdd3a664bdbc2030e0518166634350d662f2193701bf66789d@0.0.0.32:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://9ca4fc863979534d25ea380f090d14b9b06834d8b9dd0a99b5b6bd96490cffbab643a11a57d468969a9aca1e7f5b4e18a021ea96480c2f1c609be932af10f9ca@0.0.0.33:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://4075e26855d5f36c2bc1f8c925015b541a95a8119786058a931a1faa9a29288c82d8eec4a93637ce93245926832054d1f18237a6f12d6964566b1a8b9229955a@0.0.0.34:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://9b919bb71ba339d7085bb2a134c4848822ed3d30874cbf39e990bb6c99dd40e2b476370e81fe477a92d609be5e6e085592fbaa5b171b7b43b41b26478ec87071@0.0.0.35:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://652061a27cfb14000cf3d7cf3e9a1eff345ecfc6926159aa4c113cfac351825f85173747fae780579365cb7ec27563583048be0aa47e0e86bdca22c7bf8e961c@0.0.0.36:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://dd3928a718ead22ba43586fcb552b796534aec09d066a3a966bf1f3baf6cfb02fff9490795388a7c59db796507b3267aa8dbd74ecc42266b0a668b67c0058a66@0.0.0.37:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://3cd62e0d8ab04c1af5bfa843bbb06ef45b88e86373b1851f8876cffdd30442f1d245789416d284b55a97ec8000b420e98ac305bb2f89d99be3822257e23f83e4@0.0.0.38:32323?discport=32323&ntype=pn"),
			MustParseNode("kni://1ae4d5bd01661dc3e12831b5224f5938b71d2e62b00da4f6f2923d273dbcd47512df1861941a394af21827e69f9712820e8fd281ef3f8e8db23af94865f9378e@0.0.0.39:32323?discport=32323&ntype=pn"),
		},
		NodeTypeEN: {
			MustParseNode("kni://7dc47bde51af36d33655e6ddc5e447921fa89a573451ea2f5b19a95de863ea6bf8d9ecb70f58a83b181d8e2a4d2e6552e5cd93b9bb5d6afd89b343c8a1f5e9cb@0.0.0.40:32323?discport=32323&ntype=en"),
			MustParseNode("kni://a9eac3b14f6331a2f37146f32d8f408b9bf55137268c4b5b01bbcd7778d435e0c3b54422cb0b34d09c505bbf4e45b231c59c40a0fe61d1a0b42ecbc4f7e029cc@0.0.0.41:32323?discport=32323&ntype=en"),
			MustParseNode("kni://9bc13d089dd9875c1f3bd5875beed7ba347c2e94c30a3ad2a3dfba5924041d8146d28af61e5f715223eade8852c34959e95d525b6293f8bb4a485d226c3bc8c3@0.0.0.42:32323?discport=32323&ntype=en"),
			MustParseNode("kni://3ab832146617d1c65b0ea9b6c84741899a592a3cc4c9cb932088db35ae22484f32f283c70552d1978e49e3ad2d26c7676e88d8c65bdea1e54e85bdd6cdd8c576@0.0.0.43:32323?discport=32323&ntype=en"),
			MustParseNode("kni://3619dd865d84d0d56f4e2ada241002338216eba72ddc1f5e2fa4fe814afd2823e3470b5818980cfc2af223be63f4087a37bbaa1e7b20ca9b52306da6170fe185@0.0.0.44:32323?discport=32323&ntype=en"),
			MustParseNode("kni://63001c32619da5364b46e7cdb18cac53b90e9b9348674d581a71376c1b09eda637564e3c491c116c8a99c6d4e56c6ee8b45c4e87ed430467a79202cf2de4e9ef@0.0.0.45:32323?discport=32323&ntype=en"),
			MustParseNode("kni://4c3bf139758bea0392f243cb09273894446a7966ddddc434be9d093ba1365535f19251e2f10e3c452b782bb23e7f68e340bbf0109eba9d9ebd4dbefb7288771a@0.0.0.46:32323?discport=32323&ntype=en"),
			MustParseNode("kni://3034d5f70c79bb0bd7b6b7dedf39b2d60f3c75ac5ece83dd23d293cd7d5c68f84c56fa4cbde7eeafb4389d4c9c553e8c0eec14ac6c6ffa9a76ff6225e88406cd@0.0.0.47:32323?discport=32323&ntype=en"),
			MustParseNode("kni://3ba9d852e21065f03d56b0fefc7e91d192d8bd8e4cc154733a6a273771b2074cec7fa612a77f454c8a43709e82e7e359bde86388f133b4eaf197370f839ce160@0.0.0.48:32323?discport=32323&ntype=en"),
			MustParseNode("kni://3837b626353aa94a787bf4b12e6da93aaf6763fc2e825c5a43c5d59255b48a8bb73ed27a35b76171c979b947c167ddaa22eec398ea5a38b58c8027be1b792aca@0.0.0.49:32323?discport=32323&ntype=en"),
		},
		NodeTypeBN: {
			MustParseNode("kni://cea98c8c8bfb04113d217b018f341b49c348a2332eb4a8838e1d3cd1042b78dfc7155b538e9d24ac06f4b43b768f1cd2e7fe02f4e4f90da288c43af691c3b76f@0.0.0.10:32323?discport=32323&ntype=bn"),
			MustParseNode("kni://e95e84d61398a4fc36b901004a7afb223b2761aa3aefaf875c9d119d308ab37a2c21632db51a0109e9d4bd368591ba9a9d814ec0c8ff39e8b5c6f42b74b6c731@0.0.0.11:32323?discport=32323&ntype=bn"),
			MustParseNode("kni://8ec3c47e34a0e6c14ff7a82a16da98c155a57c197c4baefa2b3cebce260b6eb26b60f17d775474f859585ade32ede2c3dc6b916a7266ae10aa09f967546fa708@0.0.0.12:32323?discport=32323&ntype=bn"),
		},
	}
	testStorages = [6]*simpleStorage{
		NodeTypeUnknown: {targetType: NodeTypeUnknown, noDiscover: true, max: 100, nodes: testData[NodeTypeUnknown]},
		NodeTypeCN:      {targetType: NodeTypeCN, noDiscover: true, max: 100, nodes: testData[NodeTypeCN]},
		NodeTypePN:      {targetType: NodeTypePN, noDiscover: true, max: 100, nodes: testData[NodeTypePN]},
		NodeTypeEN:      {targetType: NodeTypeEN, noDiscover: true, max: 100, nodes: testData[NodeTypeEN]},
		NodeTypeBN:      {targetType: NodeTypeBN, noDiscover: true, max: 100, nodes: testData[NodeTypeBN]},
	}
	testNet = &simpleTestnet{
		network: testData,
	}
)

type simpleTestnet struct {
	network [5][]*Node
}

func (tn *simpleTestnet) findnode(toid NodeID, toaddr *net.UDPAddr, target NodeID, nType NodeType, max int) ([]*Node, error) {
	switch nType {
	case NodeTypeUnknown:
		return testData[NodeTypeUnknown], nil
	case NodeTypeCN:
		return testData[NodeTypeCN], nil
	case NodeTypePN:
		return testData[NodeTypePN], nil
	case NodeTypeEN:
		return testData[NodeTypeEN], nil
	case NodeTypeBN:
		return testData[NodeTypeBN], nil
	default:
		return nil, errors.New("No node type exist")
	}
}

func (*simpleTestnet) close()                                      {}
func (*simpleTestnet) waitping(from NodeID) error                  { return nil }
func (*simpleTestnet) ping(toid NodeID, toaddr *net.UDPAddr) error { return nil }

func isIn(candidate *Node, list []*Node) bool {
	for _, node := range list {
		if candidate.CompareNode(node) {
			return true
		}
	}
	return false
}

func TestShuffle(t *testing.T) {
	testStorages[NodeTypeUnknown].init()
	// shuffle empty list
	empty := testStorages[NodeTypeUnknown].shuffle([]*Node{})
	if len(empty) != 0 {
		t.Errorf("the length of the shuffled empty list. expected: %v, actual: %v\n", 0, len(empty))
	}

	// shuffle an element list
	oneElement := testStorages[NodeTypeUnknown].shuffle([]*Node{testData[NodeTypeUnknown][0]})
	if len(oneElement) != 1 {
		t.Errorf("the length of shuffled an element list. expected: %v, actual: %v\n", 1, len(oneElement))
	}
	if !oneElement[0].CompareNode(testData[NodeTypeUnknown][0]) {
		t.Errorf("the shuffled result is wrong. expected: %v, acutal: %v\n", testData[NodeTypeUnknown][0].String(), oneElement[0].String())
	}

	// shuffle the predefined list
	list := testStorages[NodeTypeUnknown].shuffle(testData[NodeTypeUnknown])
	if len(list) != len(testData[NodeTypeUnknown]) {
		t.Errorf("the length of shuffled list is wrong. expected: %v, actual: %v\n", len(testData[NodeTypeUnknown]), len(list))
	}
	isOrderChanged := false
	for idx, shuffled := range list {
		if !testData[NodeTypeUnknown][idx].CompareNode(shuffled) {
			isOrderChanged = true
		}
	}
	if !isOrderChanged {
		t.Error("the order of the list is not changed.")
	}
	for _, shuffled := range list {
		if !isIn(shuffled, testData[NodeTypeUnknown]) {
			t.Errorf("one of the elements does not exist after shuffling. missing: %v\n", shuffled.String())
		}
	}
}

func TestSimple_lookup(t *testing.T) {
	self := nodeAtDistance(common.Hash{}, 0, NodeTypeUnknown)
	conf := Config{
		udp:        testNet,
		Id:         self.ID,
		Addr:       &net.UDPAddr{},
		Bootnodes:  testData[NodeTypeBN],
		NodeDBPath: "",
		NodeType:   NodeTypeUnknown,
	}
	discv, _ := newTable(&conf)
	tab := discv.(*Table)

	// lookup on empty table returns no nodes
	typeList := []NodeType{NodeTypeUnknown, NodeTypeCN, NodeTypePN, NodeTypeEN, NodeTypeBN}
	for _, ntype := range typeList {
		storage := &simpleStorage{targetType: ntype, noDiscover: true, max: 100, nodes: testData[ntype]}
		tab.addStorage(ntype, storage)
	}

	for _, nType := range typeList {
		storage := tab.storages[nType]
		results := storage.lookup(NodeID{}, true, nType) // second parameter is not used
		if len(results) != len(testData[nType]) {
			t.Errorf("the size of lookup result is wrong. expected: %v, actual: %v\n", len(testData[nType]), len(results))
		}
		for _, actual := range results {
			if !isIn(actual, testData[nType]) {
				t.Errorf("the node does not exist in the lookup result. expected: %v", actual)
			}
		}
	}
}

func TestSimple_doRevalidate(t *testing.T) {
	self := nodeAtDistance(common.Hash{}, 0, NodeTypeUnknown)
	conf := Config{
		udp:        testNet,
		Id:         self.ID,
		Addr:       &net.UDPAddr{},
		Bootnodes:  testData[NodeTypeBN],
		NodeDBPath: "",
		NodeType:   NodeTypeUnknown,
	}
	discv, _ := newTable(&conf)
	tab := discv.(*Table)
	size := len(testData[NodeTypeUnknown])
	candidate := make([]*Node, size)
	copy(candidate, testData[NodeTypeUnknown])
	testStorages[NodeTypeUnknown].nodes = candidate
	tab.addStorage(NodeTypeUnknown, testStorages[NodeTypeUnknown])

	for idx, node := range testData[NodeTypeUnknown] {
		if !node.CompareNode(candidate[idx]) {
			t.Errorf("copy nodelist is wrong. expected: %v, actual: %v", node, candidate[idx])
		}
	}

	for i := 1; i < size; i++ {
		testStorages[NodeTypeUnknown].doRevalidate()
		for idx, node := range testData[NodeTypeUnknown] {
			if !node.CompareNode(candidate[(idx+i)%size]) {
				t.Fatalf("the result of doRevalidate is wrong. expected: %v, actual: %v", node, candidate[(idx+1)%size])
			}
		}
	}
}

func TestSimple_getNodes(t *testing.T) {
	self := nodeAtDistance(common.Hash{}, 0, NodeTypeUnknown)
	conf := Config{
		udp:        testNet,
		Id:         self.ID,
		Addr:       &net.UDPAddr{},
		Bootnodes:  testData[NodeTypeBN],
		NodeDBPath: "",
		NodeType:   NodeTypeUnknown,
	}
	discv, _ := newTable(&conf)
	tab := discv.(*Table)
	nodeTypes := []NodeType{NodeTypeUnknown, NodeTypeCN, NodeTypePN, NodeTypeEN}
	for _, nodeType := range nodeTypes {
		results := tab.GetNodes(nodeType, 1)
		if len(results) != 0 {
			t.Errorf("Returns something although there is nothing. expected: 0, actual: %v", len(results))
		}
	}
	for _, ntype := range nodeTypes {
		storage := &simpleStorage{targetType: ntype, noDiscover: true, max: 100, nodes: testData[ntype]}
		tab.addStorage(ntype, storage)
	}

	for _, nodeType := range nodeTypes {
		size := len(testData[nodeType])
		for i := 0; i <= 2*size; i++ {
			results := tab.GetNodes(nodeType, i)
			if i <= size {
				if len(results) != i {
					t.Errorf("the length of getNodes is wrong. expected: %v, acutal: %v", i, len(results))
				}
				for _, node := range results {
					if !isIn(node, testData[nodeType]) {
						t.Errorf("the result does not exist in the test data. wrong output: %v", node)
					}
				}
			} else {
				if len(results) != size {
					t.Errorf("the length of getNodes is wrong. expected: %v, acutal: %v", i, len(results))
				}
				for _, node := range results {
					if !isIn(node, testData[nodeType]) {
						t.Errorf("the result does not exist in the test data. wrong output: %v", node)
					}
				}
			}
			// Post processing in order to bond each node again
			for _, node := range results {
				tab.db.deleteNode(node.ID)
			}
		}
	}
}

func TestSimple_closest(t *testing.T) {
	typeList := []NodeType{NodeTypeUnknown, NodeTypeCN, NodeTypePN, NodeTypeEN}
	for _, ntype := range typeList {
		var target NodeID
		rand.Read(target[:])
		hash := crypto.Keccak256Hash(target[:]) // randome target hash
		storage := testStorages[ntype]
		storage.init()
		storage.max = 5
		results := storage.closest(hash, rand2.Int()) // second parameter is not used
		if len(results.entries) != 5 {
			t.Errorf("the length of the results.entries is wrong. expected: %v, actual: %v", 5, len(results.entries))
		}
		for _, node := range results.entries {
			if !isIn(node, storage.nodes) {
				t.Errorf("node does not exist in the storage. unknown node: %v", node)
			}
		}
		storage.max = 50
		results = storage.closest(hash, rand2.Int())
		if len(results.entries) != 10 {
			t.Errorf("the length of the results.entries is wrong. expected: %v, actual: %v", 10, len(results.entries))
		}
		for _, node := range results.entries {
			if !isIn(node, storage.nodes) {
				t.Errorf("node does not exist in the storage. unknown node: %v", node)
			}
		}
	}
}
