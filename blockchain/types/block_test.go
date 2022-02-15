// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/types/block_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package types

import (
	"bytes"
	"io/ioutil"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
)

func genHeader() *Header {
	return &Header{
		ParentHash:  common.HexToHash("6e3826cd2407f01ceaad3cebc1235102001c0bb9a0f8c915ab2958303bc0972c"),
		Rewardbase:  common.HexToAddress("5A0043070275d9f6054307Ee7348bD660849D90f"),
		Root:        common.HexToHash("f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb97"),
		TxHash:      common.HexToHash("f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb99"),
		ReceiptHash: common.HexToHash("f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fba9"),
		Bloom:       BytesToBloom([]byte("f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fba9f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fba9f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fba9")),
		BlockScore:  big.NewInt(10),
		Number:      big.NewInt(9),
		GasUsed:     uint64(100),
		Time:        big.NewInt(1549606547),
		TimeFoS:     20,
		Extra:       common.Hex2Bytes("0xd7820404846b6c617988676f312e31302e33856c696e75780000000000000000f90164f854942525dbdbb7ed59b8e02a6c4d3fb2a75b8b07e25094718aabda0f016e6127db6575cf0a803da7d4087b94c9ead9f875f4adc261a4b5dc264ee58039f281a794d8408db804ab30691e984e8623e2edb4cba853dfb8419da17c0fe3fdecbc32f3a2fbedf8300693067d0f944014cf575076df888709b2057869f36edc299542be1372d2b582bd8dc8e2c220059270fa37b2a2fe287ffb00f8c9b841e8765ffc1bfda438115f9bfa912f39bcc2a286fdb67c71229c9fe4084db5dd942d2076e244a4faf915aeb51a5ea097706e5421e2a7985425d0f9d6fa446c378d00b8411511d9bbd78f6a8b2151406c3c5071bcbe7a452a2ad4eebe1f9a15494ef8ff3b63b41e033de9a02c48e640d51944d7d20a462f7785525c6b26c401177521808101b841c7abdfa3ef691e8a306c8fedc32ee6af44b2fc82b921358466db9948ffce42f221a6870e01eda5ab4f54b6ee68798631e2d46a090c76c8a5d507453acaec48c401"),
		Governance:  common.Hex2Bytes("b8dc7b22676f7665726e696e676e6f6465223a22307865373333636234643237396461363936663330643437306638633034646563623534666362306432222c22676f7665726e616e63656d6f6465223a2273696e676c65222c22726577617264223a7b226d696e74696e67616d6f756e74223a393630303030303030303030303030303030302c22726174696f223a2233342f33332f3333227d2c22626674223a7b2265706f6368223a33303030302c22706f6c696379223a302c22737562223a32317d2c22756e69745072696365223a32353030303030303030307d"),
		Vote:        common.Hex2Bytes("e194e733cb4d279da696f30d470f8c04decb54fcb0d28565706f6368853330303030"),
	}
}

func genBlock() *Block {
	return &Block{
		header:       genHeader(),
		transactions: Transactions{},
	}
}

// TestBlockHash tests if the hash calculator can return the correct block hash.
func TestBlockHash(t *testing.T) {
	/*
		> klay.getBlock(49921251)
		{
			blockscore: "0x1",
			extraData: "0xd883010503846b6c617988676f312e31342e36856c696e757800000000000000f90604f901ce94bca8ffa45cc8e30bbc0522cdf1a1e0ebf540dfe29452d41ca72af615a1ac3301b0a93efa222ecc7541941782834bf8847e235f21f2c1f13fca4d5dff66219436ff2aa21d5c6828ee12cd2bc3de0e987bc0d4e7949419fa2e3b9eb1158de31be66c586a52f49c5de794e3d92072d8b9a59a0427485a1b5f459271df457c94b9456fd65a6810b19df24832c50b2e61a41867f8948a88a093c05376886754a9b70b0d0a826a5e64be9456e8c1463c341abf8b168c3079ea41ce8a387e18946873352021fe9226884616dc6f189f289aeb0cc5940b59cae1f03534209fdb9ddf5ea65b310cd7060c9453970bc504cbc41c2a0e6460aef7d866551862849416c192585a0ab24b552783b4bf7d8dc9f6855c3594e783fc94fddaeebef7293d6c5864cff280f121e194a2ba8f7798649a778a1fd66d3035904949fec555947b065fbb3a9b6e97f2ba788d63700bd4b8b408bc945e59db28cef098d5a2e877f84127aed10d7378f294ed6ee8a1877f9582858dbe2509abb0ac33e5f24e9456e3a565e31f8fb0ba0b12c03355518c6437212094386ca3cb8bb13f48d1a6adc1fb8df09e7bb7f9c8949f10d38e650184142c1c791e1b8d03e5f14ae47f94f8c9c61c5e7f2b6219d1c28b94e5cb3cdc802594b841f1cbe853c3df5fd75d25d7343e80455e103477f8b381aa0cbac56245da14b9093bcd2134bc12a7b99c02386a2a40bfd016d207c75fd0d2e34b476885448cdcd400f903edb8419e719dd155e9cc7a23f5c4e4a4ac43e4ab38758f5aa9d864267306bb0afafa8c3eaab8933cd4f0e585da5281fb4acfdd5e8ab5901ffa00dfc8bcc1146490455101b84105df89ca0f0ecebfee0358c68c9b14b9cf29af27fb70f7fefcd1d10e11c1503732280fbffec53c559bff1366cfc6e11153fab9e45d0b4230b4f418dbc7fe8d4301b841464414eadf064c3354e05b17188b8a12bb6831a7957e1e43c9f6417189863ed43b1efdd86e932fca444c80b43760adf2fa403d0938a403d7a880e9a25861caab00b841d92d24405ef652f939c91b0f12153de3d0f79f20ccff9d19b5b1907eaa37b4cb54b370cb097b8a3b3eedd484b5e59dcdbb739635a756ebd16ff67079f279292100b8410019ea8f699c67554fed7769438562bb9afcacce87d4b68aa73d79765adf60a554da08109bcae270f472e893154b2b1f5d6db11273a0fd3e386ab8a9044fb0d601b841730b78e6506ad3d2e6aa96e398a50f16350d52e67f8de48a0f1e6eb0543125515707371c9ae1896ac1ad234c6a8b69d4626e37ddfeab77b1d908762ba24bbe2300b841bc1049d862282e4a28b383c71ce08e036e486d633908dc3c4c982d3ef6209168750ffb00886e704d527cc4cb30a0f598d4888e6b645085614920192e781da89001b841949ea5ec5d8e8e3da20a54be04452a239f77e3d540a0bb1dbdcecdc7b952616d08959b9fbfa6529a968e0c9ec7b023b5271aa546799f944210ec8ed786a8efa101b84184407503c7b15a33d115721ea924e3eff4e131b4272d4c6eaab33e88390d9fbf3572f45c8ab5e7d894c4f99df6b7b52da32c0e8a4fb8903c9e81f40732503db900b841f372fb55ac864896d7d78219fc1aacd77641b4e78762e979325fff2b23b194d855939e37bc7eaaad6084dfc4cecf04806c38419ef505fd4114ee9da864bc5e2900b8417806616afaa7d21e49be276a086722cc81d8a3812c19dc80c5e32fa7dd46b3b37a387c6cf83873e2c0b114e1b833d49db2934666e2860d9b57a72fd6093118b900b841662361882fddbc22cd58eb794f06fcfeca39ccdffd642f7450037bab61920d5c06e4d65fa0397dc88875605d2821668a0ebcb94e43d72f5e27eee45860ba95b200b841cc9ba8706b8711e1a9cce58077688273709b0fdfb820a4649f624dd3eff6d6006ce999397e514553348334649b2baed97d7e0e1705459cc22c12ac3842a6e42a00b8415d7cd03c4948afd6dd29662145c661cb2a28a5da75771eddf3436db1e095237c26223e1034b4856f077f5d09a0c62442c5e8d7147caef5ab0f5af41a2edcdbe800b841762ebace1d091ec9eb161d5bfb32bdd232ecf638b4e567333cb8f38783c8973d79a7784cab13a6ceecd077b60be847113b2f8b52be2bbb16a6a496dd166d697e00",
			gasUsed: 974936,
			governanceData: "0x",
			hash: "0x5c0ba5050c597bbe3edfa4434b1bd59ef2c7e3d695bd023469e649dbea6aa02d",
			logsBloom: "0x00000000000000000000000000000000000000100000000000000000400000000004000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000040000000000000000000000000000000000010000000000000000008400000000000000080000000000000000000000000000000000000000000020008000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000008000000",
			number: 49921251,
			parentHash: "0x17ef3bfcd861e9f42f06616262196e3af53bbb03ef8ce1958434b22f2287aa52",
			receiptsRoot: "0x01125072b53afaed30d7c90b2e28b77abb2eb2af2516a2feb00fce865dff8f92",
			reward: "0x179679457f93094a4e7186abcb2089661e92fc22",
			size: 5654,
			stateRoot: "0xf7cabcd8242901bc774ca7af1d9858e7f2adfa1c8d3a44402e270ca88ba6b07f",
			timestamp: 1611578519,
			timestampFoS: 26,
			totalBlockScore: "0x2f9bce4",
			transactions: ["0x5ad7e8803c3da290f232c6ac65fa6be31846d411f57ef62f7a2a661727fcfe4a", "0x896f2db7c56efdc74c972e18c24892ddb1924cf15081e60924252a07f3df1458", "0x4a1c6f03b17ad92ad12076d94a79dfcaecf45f64d4594c3eed7cfdc52d26f0cd", "0xbe1e2342cfc7d342c6618c02e2c37b6f721275bc8b2206b3157999a3a6a2ff66", "0x9baf8ae069c1706cb3eceb2bc33404387c66e26234ddbdd0bb5eca7b794e41fe", "0x799a2110933169ce224962b86184826c1c892f4f1f2b3f21485baa69a3cc18d3", "0x0d7870349032a3cdd0bff1cba1dfef2f09b5346efb67bb69a16dcd7d2caaac58"],
			transactionsRoot: "0xfb13666763ea4562bdb3f6b288d32c963a0b05edbab59f93677aa04306784070",
			voteData: "0x"
		}
	*/

	header := &Header{
		ParentHash:  common.HexToHash("0x17ef3bfcd861e9f42f06616262196e3af53bbb03ef8ce1958434b22f2287aa52"),
		Rewardbase:  common.HexToAddress("0x179679457f93094a4e7186abcb2089661e92fc22"),
		Root:        common.HexToHash("0xf7cabcd8242901bc774ca7af1d9858e7f2adfa1c8d3a44402e270ca88ba6b07f"),
		TxHash:      common.HexToHash("0xfb13666763ea4562bdb3f6b288d32c963a0b05edbab59f93677aa04306784070"),
		ReceiptHash: common.HexToHash("0x01125072b53afaed30d7c90b2e28b77abb2eb2af2516a2feb00fce865dff8f92"),
		Bloom:       BytesToBloom(common.Hex2Bytes("00000000000000000000000000000000000000100000000000000000400000000004000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000040000000000000000000000000000000000010000000000000000008400000000000000080000000000000000000000000000000000000000000020008000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000008000000")),
		BlockScore:  big.NewInt(1),
		Number:      big.NewInt(49921251),
		GasUsed:     uint64(974936),
		Time:        big.NewInt(1611578519),
		TimeFoS:     26,
		Extra:       common.Hex2Bytes("d883010503846b6c617988676f312e31342e36856c696e757800000000000000f90604f901ce94bca8ffa45cc8e30bbc0522cdf1a1e0ebf540dfe29452d41ca72af615a1ac3301b0a93efa222ecc7541941782834bf8847e235f21f2c1f13fca4d5dff66219436ff2aa21d5c6828ee12cd2bc3de0e987bc0d4e7949419fa2e3b9eb1158de31be66c586a52f49c5de794e3d92072d8b9a59a0427485a1b5f459271df457c94b9456fd65a6810b19df24832c50b2e61a41867f8948a88a093c05376886754a9b70b0d0a826a5e64be9456e8c1463c341abf8b168c3079ea41ce8a387e18946873352021fe9226884616dc6f189f289aeb0cc5940b59cae1f03534209fdb9ddf5ea65b310cd7060c9453970bc504cbc41c2a0e6460aef7d866551862849416c192585a0ab24b552783b4bf7d8dc9f6855c3594e783fc94fddaeebef7293d6c5864cff280f121e194a2ba8f7798649a778a1fd66d3035904949fec555947b065fbb3a9b6e97f2ba788d63700bd4b8b408bc945e59db28cef098d5a2e877f84127aed10d7378f294ed6ee8a1877f9582858dbe2509abb0ac33e5f24e9456e3a565e31f8fb0ba0b12c03355518c6437212094386ca3cb8bb13f48d1a6adc1fb8df09e7bb7f9c8949f10d38e650184142c1c791e1b8d03e5f14ae47f94f8c9c61c5e7f2b6219d1c28b94e5cb3cdc802594b841f1cbe853c3df5fd75d25d7343e80455e103477f8b381aa0cbac56245da14b9093bcd2134bc12a7b99c02386a2a40bfd016d207c75fd0d2e34b476885448cdcd400f903edb8419e719dd155e9cc7a23f5c4e4a4ac43e4ab38758f5aa9d864267306bb0afafa8c3eaab8933cd4f0e585da5281fb4acfdd5e8ab5901ffa00dfc8bcc1146490455101b84105df89ca0f0ecebfee0358c68c9b14b9cf29af27fb70f7fefcd1d10e11c1503732280fbffec53c559bff1366cfc6e11153fab9e45d0b4230b4f418dbc7fe8d4301b841464414eadf064c3354e05b17188b8a12bb6831a7957e1e43c9f6417189863ed43b1efdd86e932fca444c80b43760adf2fa403d0938a403d7a880e9a25861caab00b841d92d24405ef652f939c91b0f12153de3d0f79f20ccff9d19b5b1907eaa37b4cb54b370cb097b8a3b3eedd484b5e59dcdbb739635a756ebd16ff67079f279292100b8410019ea8f699c67554fed7769438562bb9afcacce87d4b68aa73d79765adf60a554da08109bcae270f472e893154b2b1f5d6db11273a0fd3e386ab8a9044fb0d601b841730b78e6506ad3d2e6aa96e398a50f16350d52e67f8de48a0f1e6eb0543125515707371c9ae1896ac1ad234c6a8b69d4626e37ddfeab77b1d908762ba24bbe2300b841bc1049d862282e4a28b383c71ce08e036e486d633908dc3c4c982d3ef6209168750ffb00886e704d527cc4cb30a0f598d4888e6b645085614920192e781da89001b841949ea5ec5d8e8e3da20a54be04452a239f77e3d540a0bb1dbdcecdc7b952616d08959b9fbfa6529a968e0c9ec7b023b5271aa546799f944210ec8ed786a8efa101b84184407503c7b15a33d115721ea924e3eff4e131b4272d4c6eaab33e88390d9fbf3572f45c8ab5e7d894c4f99df6b7b52da32c0e8a4fb8903c9e81f40732503db900b841f372fb55ac864896d7d78219fc1aacd77641b4e78762e979325fff2b23b194d855939e37bc7eaaad6084dfc4cecf04806c38419ef505fd4114ee9da864bc5e2900b8417806616afaa7d21e49be276a086722cc81d8a3812c19dc80c5e32fa7dd46b3b37a387c6cf83873e2c0b114e1b833d49db2934666e2860d9b57a72fd6093118b900b841662361882fddbc22cd58eb794f06fcfeca39ccdffd642f7450037bab61920d5c06e4d65fa0397dc88875605d2821668a0ebcb94e43d72f5e27eee45860ba95b200b841cc9ba8706b8711e1a9cce58077688273709b0fdfb820a4649f624dd3eff6d6006ce999397e514553348334649b2baed97d7e0e1705459cc22c12ac3842a6e42a00b8415d7cd03c4948afd6dd29662145c661cb2a28a5da75771eddf3436db1e095237c26223e1034b4856f077f5d09a0c62442c5e8d7147caef5ab0f5af41a2edcdbe800b841762ebace1d091ec9eb161d5bfb32bdd232ecf638b4e567333cb8f38783c8973d79a7784cab13a6ceecd077b60be847113b2f8b52be2bbb16a6a496dd166d697e00"),
		Governance:  common.Hex2Bytes(""),
		Vote:        common.Hex2Bytes(""),
	}

	b := &Block{
		header: header,
	}

	assert.NotEqual(t, b.Hash(), common.HexToAddress("0x5c0ba5050c597bbe3edfa4434b1bd59ef2c7e3d695bd023469e649dbea6aa02d"))
	t.Log(b.Hash().String())
}

func TestBlockEncoding(t *testing.T) {
	b := genBlock()

	// To make block encoded bytes, uncomment below.
	//encB, err := rlp.EncodeToBytes(b)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(common.Bytes2Hex(encB))
	//return

	blockEnc := common.FromHex("f902adf902a9a06e3826cd2407f01ceaad3cebc1235102001c0bb9a0f8c915ab2958303bc0972c945a0043070275d9f6054307ee7348bd660849d90fa0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb97a0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb99a0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fba9b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006634313261313563623634373762643162306534386538666332643130313239326135633162623963306237386637613131323966656134663836356662613966343132613135636236343737626431623065343865386663326431303132393261356331626239633062373866376131313239666561346638363566626139663431326131356362363437376264316230653438653866633264313031323932613563316262396330623738663761313132396665613466383635666261390a0964845c5d1e931480b8deb8dc7b22676f7665726e696e676e6f6465223a22307865373333636234643237396461363936663330643437306638633034646563623534666362306432222c22676f7665726e616e63656d6f6465223a2273696e676c65222c22726577617264223a7b226d696e74696e67616d6f756e74223a393630303030303030303030303030303030302c22726174696f223a2233342f33332f3333227d2c22626674223a7b2265706f6368223a33303030302c22706f6c696379223a302c22737562223a32317d2c22756e69745072696365223a32353030303030303030307da2e194e733cb4d279da696f30d470f8c04decb54fcb0d28565706f6368853330303030c0")
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatal("decode error: ", err)
	}

	header := block.header
	println(header.String())

	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}

	// Comparing the hash from  Header and ToHeader()
	fHeader := Header{
		header.ParentHash,
		header.Rewardbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.BlockScore,
		header.Number,
		header.GasUsed,
		header.Time,
		header.TimeFoS,
		header.Extra,
		header.Governance,
		header.Vote,
	}

	resHash := rlpHash(fHeader)
	resCopiedBlockHeader := rlpHash(header)
	check("Hash", resHash, resCopiedBlockHeader)

	// Check the field value of example block.
	check("ParentHash", block.ParentHash(), b.ParentHash())
	check("Rewardbase", block.Rewardbase(), b.Rewardbase())
	check("Root", block.Root(), b.Root())
	check("TxHash", block.TxHash(), b.TxHash())
	check("ReceiptHash", block.ReceiptHash(), b.ReceiptHash())
	check("Bloom", block.Bloom(), b.Bloom())
	check("BlockScore", block.BlockScore(), b.BlockScore())
	check("NUmber", block.Number(), b.Number())
	check("GasUsed", block.GasUsed(), b.GasUsed())
	check("Time", block.Time(), b.Time())
	check("TimeFoS", block.TimeFoS(), b.TimeFoS())
	check("Extra", block.Extra(), b.Extra())
	check("Hash", block.Hash(), b.Hash())
	check("Size", block.Size(), common.StorageSize(len(blockEnc)))

	// TODO-Klaytn Consider to use new block with some transactions
	//tx1 := NewTransaction(0, common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(10), 50000, big.NewInt(10), nil)
	//
	//tx1, _ = tx1.WithSignature(HomesteadSigner{}, common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100"))
	//fmt.Println(block.Transactions()[0].Hash())
	//fmt.Println(tx1.data)
	//fmt.Println(tx1.Hash())
	//check("len(Transactions)", len(block.Transactions()), 1)
	//check("Transactions[0].Hash", block.Transactions()[0].Hash(), tx1.Hash())

	ourBlockEnc, err := rlp.EncodeToBytes(&block)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockEnc, blockEnc) {
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockEnc, blockEnc)
	}
}

func TestEIP2718BlockEncoding(t *testing.T) {
	b := genBlock()

	blockEnc := common.FromHex("f903aff902a9a06e3826cd2407f01ceaad3cebc1235102001c0bb9a0f8c915ab2958303bc0972c945a0043070275d9f6054307ee7348bd660849d90fa0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb97a0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb99a0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fba9b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006634313261313563623634373762643162306534386538666332643130313239326135633162623963306237386637613131323966656134663836356662613966343132613135636236343737626431623065343865386663326431303132393261356331626239633062373866376131313239666561346638363566626139663431326131356362363437376264316230653438653866633264313031323932613563316262396330623738663761313132396665613466383635666261390a0964845c5d1e931480b8deb8dc7b22676f7665726e696e676e6f6465223a22307865373333636234643237396461363936663330643437306638633034646563623534666362306432222c22676f7665726e616e63656d6f6465223a2273696e676c65222c22726577617264223a7b226d696e74696e67616d6f756e74223a393630303030303030303030303030303030302c22726174696f223a2233342f33332f3333227d2c22626674223a7b2265706f6368223a33303030302c22706f6c696379223a302c22737562223a32317d2c22756e69745072696365223a32353030303030303030307da2e194e733cb4d279da696f30d470f8c04decb54fcb0d28565706f6368853330303030f90100f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a8025a09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b17801f89b01800a8301e24194095e7baea6a6c7c4c2dfeb977efac326af552d878080f838f7940000000000000000000000000000000000000001e1a0000000000000000000000000000000000000000000000000000000000000000001a03dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335a0476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef14")
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatal("decode error: ", err)
	}
	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}

	check("ParentHash", block.ParentHash(), b.ParentHash())
	check("Rewardbase", block.Rewardbase(), b.Rewardbase())
	check("Root", block.Root(), b.Root())
	check("TxHash", block.TxHash(), b.TxHash())
	check("ReceiptHash", block.ReceiptHash(), b.ReceiptHash())
	check("Bloom", block.Bloom(), b.Bloom())
	check("BlockScore", block.BlockScore(), b.BlockScore())
	check("NUmber", block.Number(), b.Number())
	check("GasUsed", block.GasUsed(), b.GasUsed())
	check("Time", block.Time(), b.Time())
	check("TimeFoS", block.TimeFoS(), b.TimeFoS())
	check("Extra", block.Extra(), b.Extra())
	check("Hash", block.Hash(), b.Hash())
	check("Size", block.Size(), common.StorageSize(len(blockEnc)))

	// Create legacy tx.
	to := common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	tx1 := NewTx(&TxInternalDataLegacy{
		AccountNonce: 0,
		Recipient:    &to,
		Amount:       big.NewInt(10),
		GasLimit:     50000,
		Price:        big.NewInt(10),
	})
	sig := common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100")
	tx1, _ = tx1.WithSignature(LatestSignerForChainID(big.NewInt(1)), sig)

	// Create ACL tx.
	addr := common.HexToAddress("0x0000000000000000000000000000000000000001")
	tx2 := NewTx(&TxInternalDataEthereumAccessList{
		ChainID:      big.NewInt(1),
		AccountNonce: 0,
		Recipient:    &to,
		GasLimit:     123457,
		Price:        big.NewInt(10),
		AccessList:   AccessList{{Address: addr, StorageKeys: []common.Hash{{0}}}},
	})
	sig2 := common.Hex2Bytes("3dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef1401")
	tx2, _ = tx2.WithSignature(LatestSignerForChainID(big.NewInt(1)), sig2)
	check("len(Transactions)", len(block.Transactions()), 2)
	check("Transactions[0].Hash", block.Transactions()[0].Hash(), tx1.Hash())
	check("Transactions[1].Hash", block.Transactions()[1].Hash(), tx2.Hash())
	check("Transactions[1].Type()", block.Transactions()[1].Type(), TxType(TxTypeEthereumAccessList))

	ourBlockEnc, err := rlp.EncodeToBytes(&block)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockEnc, blockEnc) {
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockEnc, blockEnc)
	}
}

func TestEIP1559BlockEncoding(t *testing.T) {
	b := genBlock()

	blockEnc := common.FromHex("f9044ff902a9a06e3826cd2407f01ceaad3cebc1235102001c0bb9a0f8c915ab2958303bc0972c945a0043070275d9f6054307ee7348bd660849d90fa0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb97a0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fb99a0f412a15cb6477bd1b0e48e8fc2d101292a5c1bb9c0b78f7a1129fea4f865fba9b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006634313261313563623634373762643162306534386538666332643130313239326135633162623963306237386637613131323966656134663836356662613966343132613135636236343737626431623065343865386663326431303132393261356331626239633062373866376131313239666561346638363566626139663431326131356362363437376264316230653438653866633264313031323932613563316262396330623738663761313132396665613466383635666261390a0964845c5d1e931480b8deb8dc7b22676f7665726e696e676e6f6465223a22307865373333636234643237396461363936663330643437306638633034646563623534666362306432222c22676f7665726e616e63656d6f6465223a2273696e676c65222c22726577617264223a7b226d696e74696e67616d6f756e74223a393630303030303030303030303030303030302c22726174696f223a2233342f33332f3333227d2c22626674223a7b2265706f6368223a33303030302c22706f6c696379223a302c22737562223a32317d2c22756e69745072696365223a32353030303030303030307da2e194e733cb4d279da696f30d470f8c04decb54fcb0d28565706f6368853330303030f901a0f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a8025a09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b17801f89b01800a8301e24194095e7baea6a6c7c4c2dfeb977efac326af552d878080f838f7940000000000000000000000000000000000000001e1a0000000000000000000000000000000000000000000000000000000000000000001a03dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335a0476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef147802f89c01800a0a8301e24194095e7baea6a6c7c4c2dfeb977efac326af552d878080f838f7940000000000000000000000000000000000000002e1a0000000000000000000000000000000000000000000000000000000000000000001a03dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335a0476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef14")
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatal("decode error: ", err)
	}
	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}

	check("ParentHash", block.ParentHash(), b.ParentHash())
	check("Rewardbase", block.Rewardbase(), b.Rewardbase())
	check("Root", block.Root(), b.Root())
	check("TxHash", block.TxHash(), b.TxHash())
	check("ReceiptHash", block.ReceiptHash(), b.ReceiptHash())
	check("Bloom", block.Bloom(), b.Bloom())
	check("BlockScore", block.BlockScore(), b.BlockScore())
	check("NUmber", block.Number(), b.Number())
	check("GasUsed", block.GasUsed(), b.GasUsed())
	check("Time", block.Time(), b.Time())
	check("TimeFoS", block.TimeFoS(), b.TimeFoS())
	check("Extra", block.Extra(), b.Extra())
	check("Hash", block.Hash(), b.Hash())
	check("Size", block.Size(), common.StorageSize(len(blockEnc)))

	// Create legacy tx.
	to := common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	tx1 := NewTx(&TxInternalDataLegacy{
		AccountNonce: 0,
		Recipient:    &to,
		Amount:       big.NewInt(10),
		GasLimit:     50000,
		Price:        big.NewInt(10),
	})
	sig := common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100")
	tx1, _ = tx1.WithSignature(LatestSignerForChainID(big.NewInt(1)), sig)

	// Create ACL tx.
	addr := common.HexToAddress("0x0000000000000000000000000000000000000001")
	tx2 := NewTx(&TxInternalDataEthereumAccessList{
		ChainID:      big.NewInt(1),
		AccountNonce: 0,
		Recipient:    &to,
		GasLimit:     123457,
		Price:        big.NewInt(10),
		AccessList:   AccessList{{Address: addr, StorageKeys: []common.Hash{{0}}}},
	})
	sig2 := common.Hex2Bytes("3dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef1401")
	tx2, _ = tx2.WithSignature(LatestSignerForChainID(big.NewInt(1)), sig2)

	// Create DynamicFee tx.
	addr2 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	tx3 := NewTx(&TxInternalDataEthereumDynamicFee{
		ChainID:      big.NewInt(1),
		AccountNonce: 0,
		Recipient:    &to,
		GasLimit:     123457,
		GasFeeCap:    big.NewInt(10),
		GasTipCap:    big.NewInt(10),
		AccessList:   AccessList{{Address: addr2, StorageKeys: []common.Hash{{0}}}},
	})
	sig3 := common.Hex2Bytes("3dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef1401")
	tx3, _ = tx3.WithSignature(LatestSignerForChainID(big.NewInt(1)), sig3)

	check("len(Transactions)", len(block.Transactions()), 3)
	check("Transactions[0].Hash", block.Transactions()[0].Hash(), tx1.Hash())
	check("Transactions[1].Hash", block.Transactions()[1].Hash(), tx2.Hash())
	check("Transactions[1].Type()", block.Transactions()[1].Type(), TxType(TxTypeEthereumAccessList))
	check("Transactions[2].Hash", block.Transactions()[2].Hash(), tx3.Hash())
	check("Transactions[2].Type()", block.Transactions()[2].Type(), TxType(TxTypeEthereumDynamicFee))

	ourBlockEnc, err := rlp.EncodeToBytes(&block)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockEnc, blockEnc) {
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockEnc, blockEnc)
	}
}

func BenchmarkBlockEncodingHashWithInterface(b *testing.B) {
	data, err := ioutil.ReadFile("../../tests/b1.rlp")
	if err != nil {
		b.Fatal("Failed to read a block file: ", err)
	}

	blockEnc := common.FromHex(string(data))
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		b.Fatal("decode error: ", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		block.header.HashNoNonce()
	}
}

func BenchmarkBlockEncodingRlpHash(b *testing.B) {
	data, err := ioutil.ReadFile("../../tests/b1.rlp")
	if err != nil {
		b.Fatal("Failed to read a block file: ", err)
	}

	blockEnc := common.FromHex(string(data))
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		b.Fatal("decode error: ", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rlpHash(block.header)
	}
}

func BenchmarkBlockEncodingCopiedBlockHeader(b *testing.B) {
	data, err := ioutil.ReadFile("../../tests/b1.rlp")
	if err != nil {
		b.Fatal("Failed to read a block file: ", err)
	}

	blockEnc := common.FromHex(string(data))
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		b.Fatal("decode error: ", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rlpHash(block.header)
	}
}

// TODO-Klaytn-FailedTest Test fails. Analyze and enable it later.
/*
// from bcValidBlockTest.json, "SimpleTx"
func TestBlockEncoding(t *testing.T) {
	blockEnc := common.FromHex("f90260f901f9a083cafc574e1f51ba9dc0568fc617a08ea2429fb384059c972f13b19fa1c8dd55a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a05fe50b260da6308036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67a0bc37d79753ad738a6dac4921e57392f145d8887476de3f783dfa7edae9283e52b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000001832fefd8825208845506eb0780a0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff49888a13a5a8c8f2bb1c4f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1c0")
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatal("decode error: ", err)
	}

	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}
	check("BlockScore", block.BlockScore(), big.NewInt(131072))
	check("GasUsed", block.GasUsed(), uint64(21000))
	check("Root", block.Root(), common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017"))
	check("Hash", block.Hash(), common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e"))
	check("Nonce", block.Nonce(), uint64(0xa13a5a8c8f2bb1c4))
	check("Time", block.Time(), big.NewInt(1426516743))
	check("Size", block.Size(), common.StorageSize(len(blockEnc)))

	tx1 := NewTransaction(0, common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(10), 50000, big.NewInt(10), nil)

	tx1, _ = tx1.WithSignature(HomesteadSigner{}, common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100"))
	fmt.Println(block.Transactions()[0].Hash())
	fmt.Println(tx1.data)
	fmt.Println(tx1.Hash())
	check("len(Transactions)", len(block.Transactions()), 1)
	check("Transactions[0].Hash", block.Transactions()[0].Hash(), tx1.Hash())

	ourBlockEnc, err := rlp.EncodeToBytes(&block)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockEnc, blockEnc) {
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockEnc, blockEnc)
	}
}
*/
