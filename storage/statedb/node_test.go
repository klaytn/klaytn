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

func checkDecodeNode(t *testing.T, idx int, tc *testNodeEncodingTC) {
	n, err := decodeNode(tc.hash, tc.encoded)
	t.Logf("tc[%02d] %s", idx, n)
	assert.Nil(t, err, idx)

	// Result from decodeNode must contain .flags.hash
	// But tc.expanded may not contain .flags.hash. Fill it here.
	switch exp := tc.expanded.(type) {
	case *fullNode:
		exp.flags = nodeFlag{hash: tc.hash}
	case *shortNode:
		exp.flags = nodeFlag{hash: tc.hash}
	}

	assert.Equal(t, tc.expanded, n, idx)
}

func collapsedNodeTCs() []*testNodeEncodingTC {
	var (
		hex_abc = common.FromHex("0a0b0c")
		com_abc = common.FromHex("1abc")
		hex_de  = common.FromHex("0d0e10")
		com_de  = common.FromHex("20de")
		hex_f   = common.FromHex("0f10")
		com_f   = common.FromHex("3f")

		val2   = valueNode("hi")
		val4   = valueNode("tiny")
		val24  = valueNode("data_little_less_than_32")
		val31  = valueNode("data_little_less_than_32_bytes")
		val36  = valueNode("data_definitely_longer_than_32_bytes")
		hash32 = hashNode(common.FromHex("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"))
	)

	// To create testcases,
	// - Install https://github.com/ethereumjs/rlp
	//     npm install -g rlp
	// - In bash, run
	//     maketc(){ rlp encode "$1"; rlp encode "$1" | xxd -r -p | keccak-256sum; }
	//     maketc '["0x20de", "hello"]'
	testcases := []*testNodeEncodingTC{
		{ // leaf node with len(data) < 32 but lenEncoded > 32. ["0x20de", "data_little_less_than_32_bytes"]
			hash:     common.FromHex("b39ab6f7fe5f3fffa2f24c1ef22a70965186732c72bc7417ac24e7ed0e3afa38"),
			encoded:  common.FromHex("0xe28220de9e646174615f6c6974746c655f6c6573735f7468616e5f33325f6279746573"),
			expanded: &shortNode{Key: hex_de, Val: val31},
			inserted: &rawShortNode{Key: com_de, Val: val31},
		},
		{ // leaf node with len(data) >= 32. ["0x20de", "data_definitely_longer_than_32_bytes"]
			hash:     common.FromHex("916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb"),
			encoded:  common.FromHex("0xe88220dea4646174615f646566696e6974656c795f6c6f6e6765725f7468616e5f33325f6279746573"),
			expanded: &shortNode{Key: hex_de, Val: val36},
			inserted: &rawShortNode{Key: com_de, Val: val36},
		},
		{ // extension node with a hash. ["0x1abc", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"]
			hash:     common.FromHex("07eb19f1985a7e529b62bd44decd98a18b26faa10370a8b5f08bdf88cfce573e"),
			encoded:  common.FromHex("0xe4821abca000112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"),
			expanded: &shortNode{Key: hex_abc, Val: hash32},
			inserted: &rawShortNode{Key: com_abc, Val: hash32},
		},
		{ // extension node with an embedded leaf. ["0x1abc", ["0x20de", "data_little_less_than_32"]]
			hash:     common.FromHex("7604c7ac8fc0e0ac99c455cd2f05e5b4ac6929f40b937e6fc164b1d032d65173"),
			encoded:  common.FromHex("0xe0821abcdc8220de98646174615f6c6974746c655f6c6573735f7468616e5f3332"),
			expanded: &shortNode{Key: hex_abc, Val: &shortNode{Key: hex_de, Val: val24}},
			inserted: &rawShortNode{Key: com_abc, Val: &rawShortNode{Key: com_de, Val: val24}},
		},
		{ // extension node with an embedded branch. ["0x1abc", [["0x20de", "tiny"], ["0x3f", "hi"], "","","","","","","","","","","","","","",""]]
			hash:    common.FromHex("e5eb732b1c680825db277638d7c055020de26823a6d1cc16fe7b12c6cbc83620"),
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
		{ // full node with a hash and a value. ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","","","","","","","","","","","","","","","","data_little_less_than_32_bytes"]
			hash:     common.FromHex("16cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2"),
			encoded:  common.FromHex("0xf84fa000112233445566778899aabbccddeeff00112233445566778899aabbccddeeff8080808080808080808080808080809e646174615f6c6974746c655f6c6573735f7468616e5f33325f6279746573"),
			expanded: &fullNode{Children: [17]node{0: hash32, 16: val31}},
			inserted: rawFullNode([17]node{0: hash32, 16: val31}),
		},
		{ // full node with an embedded leaf. [["0x20de","tiny"],"","","","","","","","","","","","","","","","data_little_less_than_32"]
			hash:    common.FromHex("9b112f76f00fd35094375d08864e936fd9a00bb92509006fb020aeb0f2e10714"),
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
	}
	return testcases
}

func resolvedNodeTCs() []*testNodeEncodingTC {
	var (
		hex_abc = common.FromHex("0a0b0c")
		com_abc = common.FromHex("1abc")
		hex_de  = common.FromHex("0d0e10")

		val31  = valueNode("data_little_less_than_32_bytes")
		val36  = valueNode("data_definitely_longer_than_32_bytes")
		hash32 = hashNode(common.FromHex("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"))
	)

	var (
		// leaf node with more than 32 bytes. ["0x20de", "data_definitely_longer_than_32_bytes"]
		shortChildHash = hashNode(common.FromHex("916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb"))
		// full node with more than 32 bytes. ["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","","","","","","","","","","","","","","","","data_little_less_than_32_bytes"]
		fullChildHash = hashNode(common.FromHex("16cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2"))
	)

	testcases := []*testNodeEncodingTC{
		{ // full parent, short child. ["0x916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb","","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("f240b93556d06a01f7220537a11bcf4c02c690e7ff5a97fa6c98157e7298b0f1"),
			encoded: common.FromHex("0xf1a0916d1af8295fc9ad621dc43ecfa4c9f25286a619715b270929d0425cd37c80eb80808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				0: &shortNode{Key: hex_de, Val: val36},
			}},
			inserted: rawFullNode([17]node{
				0: shortChildHash,
			}),
		},
		{ // short parent, full child. ["0x1abc", "0x16cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2"]
			hash:    common.FromHex("e94b214debb6a08829ad81dcb171a7dfda58b3c382dfffd9d929b8a832fc1e68"),
			encoded: common.FromHex("0xe4821abca016cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2"),
			expanded: &shortNode{
				Key: hex_abc,
				Val: &fullNode{Children: [17]node{
					0:  hash32,
					16: val31,
				}},
			},
			inserted: &rawShortNode{
				Key: com_abc,
				Val: fullChildHash,
			},
		},
		{ // full parent, full child. ["","0x16cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2","","","","","","","","","","","","","","",""]
			hash:    common.FromHex("284f0dc1379f6c36d8c49ac41ed5e091d39257047a8d6d9f5d652dc98994e14f"),
			encoded: common.FromHex("0xf180a016cf3239c032f06e83a71ee519f4c58a41e5e3e609f641841d65bb49c09d59d2808080808080808080808080808080"),
			expanded: &fullNode{Children: [17]node{
				1: &fullNode{Children: [17]node{
					0:  hash32,
					16: val31,
				}},
			}},
			inserted: rawFullNode([17]node{
				1: fullChildHash,
			}),
		},
	}
	return testcases
}

func TestDecodeNodeTC(t *testing.T) {
	for idx, tc := range collapsedNodeTCs() {
		checkDecodeNode(t, idx, tc)
	}
}
