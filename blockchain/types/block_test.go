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
