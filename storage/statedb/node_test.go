// Copyright 2023 The klaytn Authors
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
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

// testNodeEncodingTC represents a trie node in various forms.
type testNodeEncodingTC struct {
	hash     []byte // Keccak(encoded)
	encoded  []byte // RLP(inserted)
	expanded node   // node{..., flags: {Hash}}
	inserted node   // rawNode{...}
}

var nodeEncodingTCs map[string]testNodeEncodingTC

func hexToHashNode(s string) hashNode {
	return hashNode(common.HexToHash(s).Bytes())
}

func hexToExtHashNode(s string) hashNode {
	return hashNode(common.BytesToExtHash(common.FromHex(s)).Bytes())
}

func init() {
	var (
		// subpaths in HEX and COMPACT encodings
		hex_abc = common.FromHex("0a0b0c")
		com_abc = common.FromHex("1abc")
		hex_de  = common.FromHex("0d0e10")
		com_de  = common.FromHex("20de")
		hex_f   = common.FromHex("0f10")
		com_f   = common.FromHex("3f")

		// plain byte values
		val2  = valueNode("hi")
		val4  = valueNode("tiny")
		val24 = valueNode("data_little_less_than_32")
		val31 = valueNode("data_little_less_than_32_bytes")
		val36 = valueNode("data_definitely_longer_than_32_bytes")

		// dummy ref values
		hash32    = hexToHashNode("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff")
		extlegacy = hexToExtHashNode("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00000000000000")
		exthash   = hexToExtHashNode("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344")

		// computed ref from hasher
		hLegacy = "00000000000000"
		hExt    = "ccccddddeeee01"
		hExt2   = "ccccddddeeee02"
	)

	// To create testcases,
	// - Install https://github.com/ethereumjs/rlp
	//     npm install -g rlp
	// - In bash, run
	//     maketc(){ rlp encode "$1"; rlp encode "$1" | xxd -r -p | keccak-256sum; }
	//     maketc '["0x20de", "hello"]'

	// TC names:
	// - collapsed/*: single-depth node as stored in databased. decodeNode() can handle it.
	// - resolved/*: multi-depth node exists in memory. decodeNode() cannot handle it.
	// - */legacy: only contains hash32
	// - */extroot: contains exthash, and the node is state root
	// - */exthash: contains exthash, and the node is not state root

	nodeEncodingTCs = map[string]testNodeEncodingTC{
		// Leaf nodes
		"collapsed/leaf_short/legacy": {
			// leaf node with len(data) < 32 but lenEncoded > 32. ["0x20de", "data_little_less_than_32_bytes"]
			hash:     common.FromHex("b39ab6f7fe5f3fffa2f24c1ef22a70965186732c72bc7417ac24e7ed0e3afa38" + hLegacy),
			encoded:  common.FromHex("0xe28220de9e646174615f6c6974746c655f6c6573735f7468616e5f33325f6279746573"),
			expanded: &shortNode{Key: hex_de, Val: val31},
			inserted: &rawShortNode{Key: com_de, Val: val31},
		},
		"collapsed/leaf_long/legacy": {
			// leaf node with len(data) >= 32. ["0x20de", "data_definitely_longer_than_32_bytes"]
			hash:     common.FromHex("916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb" + hLegacy),
			encoded:  common.FromHex("0xe88220dea4646174615f646566696e6974656c795f6c6f6e6765725f7468616e5f33325f6279746573"),
			expanded: &shortNode{Key: hex_de, Val: val36},
			inserted: &rawShortNode{Key: com_de, Val: val36},
		},

		// Extension nodes without ref
		"collapsed/extension_embedded_leaf/legacy": {
			// extension node with an embedded leaf. ["0x1abc", ["0x20de", "data_little_less_than_32"]]
			hash:     common.FromHex("7604c7ac8fc0e0ac99c455cd2f05e5b4ac6929f40b937e6fc164b1d032d65173" + hLegacy),
			encoded:  common.FromHex("0xe0821abcdc8220de98646174615f6c6974746c655f6c6573735f7468616e5f3332"),
			expanded: &shortNode{Key: hex_abc, Val: &shortNode{Key: hex_de, Val: val24}},
			inserted: &rawShortNode{Key: com_abc, Val: &rawShortNode{Key: com_de, Val: val24}},
		},
		"collapsed/extension_embedded_branch/legacy": {
			// extension node with an embedded branch. ["0x1abc", [["0x20de", "tiny"], ["0x3f", "hi"], "","","","","","","","","","","","","","",""]]
			hash:    common.FromHex("e5eb732b1c680825db277638d7c055020de26823a6d1cc16fe7b12c6cbc83620" + hLegacy),
			encoded: common.FromHex("0xe1821abcddc88220de8474696e79c43f826869808080808080808080808080808080"),
			expanded: &shortNode{
				Key: hex_abc,
				Val: &fullNode{Children: [17]node{
					0: &shortNode{Key: hex_de, Val: val4},
					1: &shortNode{Key: hex_f, Val: val2},
				}},
			},
			inserted: &rawShortNode{
				Key: com_abc,
				Val: rawFullNode([17]node{
					0: &rawShortNode{Key: com_de, Val: val4},
					1: &rawShortNode{Key: com_f, Val: val2},
				}),
			},
		},

		// Branch nodes without ref
		"collapsed/branch_embedded_leaf/legacy": {
			// full node with an embedded leaf. [["0x20de","tiny"],"","","","","","","","","","","","","","","","data_little_less_than_32"]
			hash:    common.FromHex("9b112f76f00fd35094375d08864e936fd9a00bb92509006fb020aeb0f2e10714" + hLegacy),
			encoded: common.FromHex("0xf1c88220de8474696e7980808080808080808080808080808098646174615f6c6974746c655f6c6573735f7468616e5f3332"),
			expanded: &fullNode{Children: [17]node{
				0:  &shortNode{Key: hex_de, Val: val4},
				16: val24,
			}},
			inserted: rawFullNode([17]node{
				0:  &rawShortNode{Key: com_de, Val: val4},
				16: val24,
			}),
		},

		// Extension nodes with ref
		"collapsed/extension_ref/legacy": {
			// extension node with a hash. ["0x1abc", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"]
			hash:     common.FromHex("07eb19f1985a7e529b62bd44decd98a18b26faa10370a8b5f08bdf88cfce573e" + hLegacy),
			encoded:  common.FromHex("0xe4821abca000112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"),
			expanded: &shortNode{Key: hex_abc, Val: extlegacy}, // decodeNode will hash32.ExtendLegacy()
			inserted: &rawShortNode{Key: com_abc, Val: hash32}, // hasher will canonicalRef() to hash32
		},
		"collapsed/extension_ref/extroot": {
			// extension node with a hash.
			// hashing ["0x1abc", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"]
			// storing ["0x1abc", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344"]
			hash:     common.FromHex("07eb19f1985a7e529b62bd44decd98a18b26faa10370a8b5f08bdf88cfce573e" + hLegacy),
			encoded:  common.FromHex("0xeb821abca700112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344"),
			expanded: &shortNode{Key: hex_abc, Val: exthash},    // decodeNode will leave exthash as-is
			inserted: &rawShortNode{Key: com_abc, Val: exthash}, // hasher will leave exthash as-is
		},
		"collapsed/extension_ref/exthash": {
			// extension node with a hash.
			// hashing ["0x1abc", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"]
			// storing ["0x1abc", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344"]
			hash:     common.FromHex("07eb19f1985a7e529b62bd44decd98a18b26faa10370a8b5f08bdf88cfce573e" + hExt),
			encoded:  common.FromHex("0xeb821abca700112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344"),
			expanded: &shortNode{Key: hex_abc, Val: exthash},    // decodeNode will leave exthash as-is
			inserted: &rawShortNode{Key: com_abc, Val: exthash}, // hasher will leave exthash as-is
		},

		// Branch nodes with ref
		"collapsed/branch_ref/legacy": {
			// branch node with a hash and a value. ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","","","","","","","","","","","","","","","","data_little_less_than_32_bytes"]
			hash:     common.FromHex("16cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2" + hLegacy),
			encoded:  common.FromHex("0xf84fa000112233445566778899aabbccddeeff00112233445566778899aabbccddeeff8080808080808080808080808080809e646174615f6c6974746c655f6c6573735f7468616e5f33325f6279746573"),
			expanded: &fullNode{Children: [17]node{0: extlegacy, 16: val31}}, // decodeNode will hash32.ExtendLegacy()
			inserted: rawFullNode([17]node{0: hash32, 16: val31}),            // hasher will canonicalRef() to hash32
		},
		"collapsed/branch_ref/extroot": {
			// hashing ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","","","","","","","","","","","","","","","","data_little_less_than_32_bytes"]
			// storing ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344","","","","","","","","","","","","","","","","data_little_less_than_32_bytes"]
			hash:     common.FromHex("16cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2" + hLegacy),
			encoded:  common.FromHex("0xf856a700112233445566778899aabbccddeeff00112233445566778899aabbccddeeff111122223333448080808080808080808080808080809e646174615f6c6974746c655f6c6573735f7468616e5f33325f6279746573"),
			expanded: &fullNode{Children: [17]node{0: exthash, 16: val31}}, // decodeNode will leave exthash as-is
			inserted: rawFullNode([17]node{0: exthash, 16: val31}),         // hasher will leave exthash as-is
		},
		"collapsed/branch_ref/exthash": {
			// hashing ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","","","","","","","","","","","","","","","","data_little_less_than_32_bytes"]
			// storing ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344","","","","","","","","","","","","","","","","data_little_less_than_32_bytes"]
			hash:     common.FromHex("16cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2" + hExt),
			encoded:  common.FromHex("0xf856a700112233445566778899aabbccddeeff00112233445566778899aabbccddeeff111122223333448080808080808080808080808080809e646174615f6c6974746c655f6c6573735f7468616e5f33325f6279746573"),
			expanded: &fullNode{Children: [17]node{0: exthash, 16: val31}}, // decodeNode will leave exthash as-is
			inserted: rawFullNode([17]node{0: exthash, 16: val31}),         // hasher will leave exthash as-is
		},

		// Branch and extension nodes with a resovled child node
		"resolved/branch_leaf/legacy": {
			// parent = ["0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb","","","","","","","","","","","","","","","",""]
			// child  = ["0x20de", "data_definitely_longer_than_32_bytes"]
			hash:    common.FromHex("f240b93556d06a01f7220537a11bcf4c02c690e7ff5a97fa6c98157e7298b0f1" + hLegacy),
			encoded: common.FromHex("0xf1a0916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb80808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				0: &shortNode{Key: hex_de, Val: val36},
			}},
			inserted: rawFullNode([17]node{
				0: hashNode(common.FromHex("0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb")),
			}),
		},
		"resolved/branch_leaf/extroot": {
			// parent hashing = ["0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb","","","","","","","","","","","","","","","",""]
			// parent storing = ["0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80ebccccddddeeee01","","","","","","","","","","","","","","","",""]
			// child          = ["0x20de", "data_definitely_longer_than_32_bytes"]
			hash:    common.FromHex("f240b93556d06a01f7220537a11bcf4c02c690e7ff5a97fa6c98157e7298b0f1" + hLegacy),
			encoded: common.FromHex("0xf838a7916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80ebccccddddeeee0180808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				0: &shortNode{Key: hex_de, Val: val36},
			}},
			inserted: rawFullNode([17]node{
				0: hashNode(common.FromHex("0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb" + hExt)),
			}),
		},
		"resolved/branch_leaf/exthash": {
			// parent hashing = ["0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb","","","","","","","","","","","","","","","",""]
			// parent storing = ["0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80ebccccddddeeee01","","","","","","","","","","","","","","","",""]
			// child          = ["0x20de", "data_definitely_longer_than_32_bytes"]
			hash:    common.FromHex("f240b93556d06a01f7220537a11bcf4c02c690e7ff5a97fa6c98157e7298b0f1" + hExt2), // parent node hashed after child
			encoded: common.FromHex("0xf838a7916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80ebccccddddeeee0180808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				0: &shortNode{Key: hex_de, Val: val36},
			}},
			inserted: rawFullNode([17]node{
				0: hashNode(common.FromHex("0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb" + hExt)), // child node first hashed
			}),
		},

		"resolved/extension_branch/legacy": {
			// parent = ["0x1abc", "0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7"]
			// child  = ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("3cd041acf721fcf611e8a06c39e347cb7daed5e6e9e9e9951473f8f920082659" + hLegacy),
			encoded: common.FromHex("0xe4821abca0df5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7"),
			expanded: &shortNode{
				Key: hex_abc,
				Val: &fullNode{Children: [17]node{0: hash32}},
			},
			inserted: &rawShortNode{
				Key: com_abc,
				Val: hashNode(common.FromHex("0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7")),
			},
		},
		"resolved/extension_branch/extroot": {
			// parent hashing = ["0x1abc", "0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7"]
			// parent storing = ["0x1abc", "0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01"]
			// child          = ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344","","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("3cd041acf721fcf611e8a06c39e347cb7daed5e6e9e9e9951473f8f920082659" + hLegacy),
			encoded: common.FromHex("0xeb821abca7df5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01"),
			expanded: &shortNode{
				Key: hex_abc,
				Val: &fullNode{Children: [17]node{0: exthash}},
			},
			inserted: &rawShortNode{
				Key: com_abc,
				Val: hashNode(common.FromHex("0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7" + hExt)),
			},
		},
		"resolved/extension_branch/exthash": {
			// parent hashing = ["0x1abc", "0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7"]
			// parent storing = ["0x1abc", "0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01"]
			// child          = ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344","","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("3cd041acf721fcf611e8a06c39e347cb7daed5e6e9e9e9951473f8f920082659" + hExt2), // parent node hashed after child
			encoded: common.FromHex("0xeb821abca7df5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01"),
			expanded: &shortNode{
				Key: hex_abc,
				Val: &fullNode{Children: [17]node{0: exthash}},
			},
			inserted: &rawShortNode{
				Key: com_abc,
				Val: hashNode(common.FromHex("0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7" + hExt)), // child node first hashed
			},
		},

		"resolved/branch_branch/legacy": {
			// parent = ["","0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7","","","","","","","","","","","","","","",""]
			// child  = ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("e03dce39399489a1bf77c4d766668fddb2fff56f52090c71e16689e75ac97a25" + hLegacy),
			encoded: common.FromHex("0xf180a0df5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				1: &fullNode{Children: [17]node{0: hash32}},
			}},
			inserted: rawFullNode([17]node{
				1: hashNode(common.FromHex("0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7")),
			}),
		},
		"resolved/branch_branch/extroot": {
			// parent hashing = ["","0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7","","","","","","","","","","","","","","",""]
			// parent storing = ["","0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01","","","","","","","","","","","","","","",""]
			// child          = ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344","","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("e03dce39399489a1bf77c4d766668fddb2fff56f52090c71e16689e75ac97a25" + hLegacy),
			encoded: common.FromHex("0xf83880a7df5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				1: &fullNode{Children: [17]node{0: exthash}},
			}},
			inserted: rawFullNode([17]node{
				1: hashNode(common.FromHex("0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7" + hExt)),
			}),
		},
		"resolved/branch_branch/exthash": {
			// parent hashing = ["","0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7","","","","","","","","","","","","","","",""]
			// parent storing = ["","0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01","","","","","","","","","","","","","","",""]
			// child          = ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff11112222333344","","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("e03dce39399489a1bf77c4d766668fddb2fff56f52090c71e16689e75ac97a25" + hExt2), // parent node hashed after child
			encoded: common.FromHex("0xf83880a7df5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7ccccddddeeee01808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				1: &fullNode{Children: [17]node{0: hash32}},
			}},
			inserted: rawFullNode([17]node{
				1: hashNode(common.FromHex("0xdf5551a5661b69abe1ad76df12ceb6471f01b2639a24dbd2f8fe72d26e0dffb7" + hExt)), // child node first hashed
			}),
		},
	}
}

func selectNodeEncodingTCs(names []string) map[string]testNodeEncodingTC {
	testcases := map[string]testNodeEncodingTC{}
	for _, name := range names {
		testcases[name] = nodeEncodingTCs[name]
	}
	return testcases
}

func collapsedNodeTCs_legacy() map[string]testNodeEncodingTC {
	return selectNodeEncodingTCs([]string{
		"collapsed/leaf_short/legacy",
		"collapsed/leaf_long/legacy",
		"collapsed/extension_embedded_leaf/legacy",
		"collapsed/extension_embedded_branch/legacy",
		"collapsed/branch_embedded_leaf/legacy",
		"collapsed/extension_ref/legacy",
		"collapsed/branch_ref/legacy",
	})
}

func collapsedNodeTCs_extroot() map[string]testNodeEncodingTC {
	return selectNodeEncodingTCs([]string{
		"collapsed/extension_ref/extroot",
		"collapsed/branch_ref/extroot",
	})
}

func collapsedNodeTCs_exthash() map[string]testNodeEncodingTC {
	return selectNodeEncodingTCs([]string{
		"collapsed/extension_ref/exthash",
		"collapsed/branch_ref/exthash",
	})
}

func resolvedNodeTCs_legacy() map[string]testNodeEncodingTC {
	return selectNodeEncodingTCs([]string{
		"resolved/branch_leaf/legacy",
		"resolved/extension_branch/legacy",
		"resolved/branch_branch/legacy",
	})
}

func resolvedNodeTCs_extroot() map[string]testNodeEncodingTC {
	return selectNodeEncodingTCs([]string{
		"resolved/branch_leaf/extroot",
		"resolved/extension_branch/extroot",
		"resolved/branch_branch/extroot",
	})
}

func resolvedNodeTCs_exthash() map[string]testNodeEncodingTC {
	return selectNodeEncodingTCs([]string{
		"resolved/branch_leaf/exthash",
		"resolved/extension_branch/exthash",
		"resolved/branch_branch/exthash",
	})
}

func checkDecodeNode(t *testing.T, name string, tc testNodeEncodingTC) {
	n, err := decodeNode(tc.hash, tc.encoded)
	t.Logf("tc[%s]\n %s", name, n)
	assert.Nil(t, err, name)

	// Result from decodeNode must contain .flags.hash
	// But tc.expanded may not contain .flags.hash. Fill it here.
	switch exp := tc.expanded.(type) {
	case *fullNode:
		exp.flags = nodeFlag{hash: tc.hash}
	case *shortNode:
		exp.flags = nodeFlag{hash: tc.hash}
	}

	assert.Equal(t, tc.expanded, n, name)
}

func TestDecodeNodeTC(t *testing.T) {
	for name, tc := range collapsedNodeTCs_legacy() {
		checkDecodeNode(t, name, tc)
	}
	for name, tc := range collapsedNodeTCs_exthash() {
		checkDecodeNode(t, name, tc)
	}
	for name, tc := range collapsedNodeTCs_extroot() {
		checkDecodeNode(t, name, tc)
	}
}
