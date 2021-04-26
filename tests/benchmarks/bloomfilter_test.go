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

package benchmarks

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var addrs = []string{
	"0x229E467d46913DfdB70D8f13c2461665F7FeB2a9",
	"0x312446899c6d5b05AD45a71a762bba7D0DEa074a",
	"0x6705CeeE1e40CC5120170ea022e4f91de63fB379",
	"0x2b789F8b96D2D8cd95814754032e6183b898DDcA",
	"0x26Facf5aA382Eb51dCCa1742Eb0a3cc9d8FE9392",
	"0x545CD2e0E3ddf6F5632439b1e655B6d0e9a73B8b",
	"0x34f9A1bE6A66177395469f6E982d4563D0870827",
	"0x2bf3B29c19C9EF545E27f6085DB9b9B26Db8B656",
	"0xA02732AaFAF6883278339A22039CaFac5e059282",
	"0x1d3cdb28cf75d00b63E9e95b2bcBda767ea7d231",
	"0x756C7bf861b64cf5F1B23a231a90d4D096126A68",
	"0x66ee129A97D1C87e492d50B34569b8324F8600DF",
	"0x3f46db8390BC22FAf7f34675B58d35b959C8e0C9",
	"0x30E06cebB43b776091EFfd1dE182e36AF2eb4FA3",
	"0x230e7142Fd5FF14d8A177c6E6fbF76aadbF4bC5F",
	"0xE0e3C8752b368c921b3c520797E0bBC2EB50E7Ed",
	"0xCC7d575D6d538A987082ad5a1DCAe9e1c469Bd95",
	"0xCD2DC73D736c31D7aa98209dD4068979F19e486B",
	"0xC97e706e31c3E22E8d17C62Ef886046C8BC09864",
	"0x48a3f5Af534C01Ad35aFB2b86fAE3e0cDd61f6F4",
	"0xa26B94C2c708A76DA440Dc245Dc03e031c1AfD05",
	"0x30831D05a46f5C62F209Fe0524d1a9F731d74247",
	"0x7fA4780fAe30D9a0b0e2131F993598C19524f7D1",
	"0xED18B2a4C14BA903ad3156FBE5306a23D5594218",
	"0x9B80efC4eD1EC0B012526cb3407758ddD6b9B640",
	"0x438A7fDC25b0a8880c51529A4C1459B32405496B",
	"0xe5cA882eB33B9fB1c3D33FE482b79Db3D2f4B473",
	"0x21f58c4d24CE65279a9457bc96c0ABC2a581208c",
	"0x6288bfB0da903a4eb4641dB008C736588Cd60DCB",
	"0xc10320Af5D9CC75ab00e6BED91dcA817f2A4cBDf",
	"0x66307bC0A5Db3Cd87bcD5fc5281bdDD937D09b6D",
	"0x782ED08B2C25A3979e32c15d191B00f20f5A80d7",
	"0xD6E3cdb9Bb23ef2E1268722e2e91c0509e8E9fA7",
	"0x2D39aB465327e1bD7962F4eaCD75843a555140A1",
	"0x3c8B900D9E8ee88232Dc836F039a88d36e8446cB",
	"0x629dd6B9B57D4CeA689804373bA6B562D984cD93",
	"0x6381F02d3Da9823EBD7FF1BF2c198Fb5F0Ff8460",
	"0xd24151F529a594Eb0fE7A751261A2CD39ec76845",
	"0xAFdaeDf63E67228B0ddc472F2266B0715779F852",
	"0x9EbD57f9611eD745C08B00b6cc7c657C910e3c11",
	"0xF648E5500218a5012FBE2E65c184d136E58aA2ED",
	"0x76fD23E0995BfAD8E3896827D471478C122541Ea",
	"0xEe562Fef72128e530C0Ae4C86F16D714eAd1A07E",
	"0x37197FFD3b1c0BE35103D883519ECe09Af15Ca35",
	"0xf94eC4C22724557d94e69f67B9D3567eD65f3d37",
	"0xfA50A75cEAACceB321e3f92b293e121737C79Ea6",
	"0x9DFF8E6327D692a31Bb2fB8dc68F88D054b4aa3a",
	"0xD569e27092a0CB89B6bf3E7D645cb6c9a350474a",
	"0x963cF4eA4AAFA550Ba100F1e22703b80a1DBdFE1",
	"0x05B12922B33D8e294e8C70e49BC1DAb261C1E517",
	"0x50dc2304b1e1A6bC6619e7Eb523b6D2cdf255961",
	"0xbdB837580481862e2b12336f4A26Af4b23A6dF30",
	"0xda1d187aEEdfBF91E6D8e62bBF568C3D96FA3af1",
	"0xAC99De0933A6124e3ac260706C1d34f0651f27A5",
	"0x2964d12F17742FBAccAB7FCE9a790fA168a87325",
	"0x88E39847786D39A7cdAA7Aa53e9E9B1A023ab844",
	"0x6d981748AC668B68d803DA5C69562F4ae6A1452d",
	"0xf886C291a646b6fDc54441D9BEEd11a9f78AEB15",
	"0xF0f54c01Fe4D2363C1D2F377d8Cca525c6C130dC",
	"0x8650a24e65d42a32bda84AC9DD02B8902eA821c2",
	"0x70EcB9089533C3c2Ea53E7Df300F0A33f7bFC61b",
	"0x56a619032B2999D0e104519bD0bb74B0c46E378E",
	"0xfd80977D9550A5B5319f7dadb3674D29193b8C60",
	"0xaFBAb54E1999eC71086fac001deef33d4E156cAE",
	"0xd1448E5CFC6440c2437912d51f2080c421c2C638",
	"0x7D13Fb871E42556D2Dcf5b6A6126b68345Cc347A",
	"0x7849729e108631224c197dB4C1EBf272D98F793e",
	"0x3b8822D9e34F1090670c64F99d4cb6a2b4795FAa",
	"0x09271289E82d292Bb8d2DAb48eb7c5ABB591A630",
	"0xaE21c5ACE28dD7393f770b5218aB9bAea5456352",
	"0x84F4351c067dd981d6e8Db6dEC96bA5aed0CF010",
	"0xd941E547Ce39d654505cFB5211D6B87ff2811519",
	"0x7069dc629305D5f02fAdB542c0634f3BD124137E",
	"0x4ee5Ef6fcC0E1A15b84547101064d70dD2AD22df",
	"0x956171aEeb695A125b770D5AAD6C6c0Ceb6bdBe8",
	"0x395FbF1dD980c8Cf10586D0df4945016526BE786",
	"0xAb456bc354b1003F10BeD63C001CaeE55b4Be675",
	"0x277570CDcf52dB43E7F307CbF1E58748f5F3f8E8",
	"0x5F048d235ddc6E727C0B9D1FF11510767Eb172Cb",
	"0x29299F84810af7e005880B65f96485e2251051e5",
	"0xA9e0e5d8B704Ec921E8a55c92353E23910bAdbC8",
	"0x9C7ee419706779C4233C0e6a547f9a67AF1AE04E",
	"0xDb91330c323961937b73586E6Cbd19a2642e08a9",
	"0x4e2bB263D843cCDeeB8c3269bca732B9158FaCA1",
	"0xC3c33Ec0C89F0E64F84F3D561BF6259BedB28A43",
	"0x56a2639569035159BfC6BDEb1654F4FD57a016FB",
	"0xA747aba1bc48eAb2e0e68a310C81F0cD679cAf74",
	"0x6e86f162Be0117f9F921D4328f2DC8ada2d345B0",
	"0x088A4f7d637998d8d7ac53D2f5e9Bf89C0b308eF",
	"0xda90F4bBa7ACE86cFB097d90024A58fB59c7D5b5",
	"0x974ed4327fDb8786A3C543410858840DD1674631",
	"0xe92Be80a0bb0749d30544964d7E4437cF9b02f57",
	"0x0F07efB587C192BF4666990F72AD7f853e7F85C6",
	"0x4984c77B4752B38B86E2CBBAD34625c763c2DD57",
	"0x3aB6E0759a9d5Ba5EdD15A69F154dE57998829ee",
	"0xbcF2Da0fc549e1160555104EE24158B866Dec35C",
	"0xed0cB5774014bc66F03dfDdfeEf9790e2097A320",
	"0x9c8A0e62E4b48a9aB8f4f8cd5A31fF36Ddbdf82F",
	"0xf284f64f865d379f06b2e9D7d5714DFA3860e8de",
	"0x7647D15CC3C58FCeD12a8cB40c5603345aD7f37d",
}

var badAddrs = []string{
	"0xc4E9fe2E3f040964770Fa0A2AE2f549292F318fc",
	"0x5FdC97cD5ae70500B8A01b474c0f77a0c3845575",
	"0xEB5cFBDd1c69C05E7FD12b214c07141F8413fe31",
	"0x0F2A66C73fbbff5fD21Cd4F3ACb0491B61c58bA9",
	"0x80F26DC89A19428E79156E66753E602e4b421dA8",
	"0x277570CDcf52dB43E7F307CbF1E58748f5F3f8E7",
	"0x5F048d235ddc6E727C0B9D1FF11510767Eb172C0",
	"0x29299F84810af7e005880B65f96485e2251051e0",
	"0xA9e0e5d8B704Ec921E8a55c92353E23910bAdbC6",
	"0x9C7ee419706779C4233C0e6a547f9a67AF1AE04A",
}

const bitsPerKey = 16

type testOption struct {
	testAddrs []string
	expected  bool
	msg       string
}

var lookupBenchmarks = [...]struct {
	name string
	opt  testOption
}{
	{"True", testOption{addrs, true, "expected"}},
	{"False", testOption{badAddrs, false, "did not expect"}},
}

var compareBenchmarks = [...]struct {
	name string
	opt  testOption
}{
	{"Same", testOption{addrs, true, "to be the same as"}},
	{"Different", testOption{badAddrs, false, "to be different from"}},
}

var bitandBenchmarks = [...]struct {
	name string
	opt  testOption
}{
	{"True", testOption{addrs, true, "expected"}},
	{"False", testOption{badAddrs, false, "did not expect"}},
}

func bloomBitAND(bloom1, bloom2 types.Bloom) types.Bloom {
	bloom1Bytes := bloom1.Bytes()
	bloom2Bytes := bloom2.Bytes()
	var resultBytes types.Bloom

	for i := 0; i < types.BloomByteLength; i++ {
		resultBytes[i] = bloom1Bytes[i] & bloom2Bytes[i]
	}

	return resultBytes
}

func ldbBloomBitAND(bloom1, bloom2 []byte) []byte {
	lenFilter := len(bloom1)
	resultBytes := make([]byte, lenFilter)

	for i := 0; i < lenFilter; i++ {
		resultBytes[i] = bloom1[i] & bloom2[i]
	}

	return resultBytes
}

func Benchmark_Map_Add(b *testing.B) {
	toset := make(map[string]int)
	numAddrs := len(addrs)

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		addr := addrs[k%numAddrs]
		toset[addr] = k
	}
}

func benchmark_Map_Lookup(b *testing.B, toset map[string]int, opt *testOption) {
	numAddrs := len(opt.testAddrs)

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		addr := opt.testAddrs[k%numAddrs]
		if _, exists := toset[addr]; exists != opt.expected {
			b.Error(opt.msg, addr, "to test true")
		}
	}
}

func Benchmark_Map_Lookup(b *testing.B) {
	toset := make(map[string]int)
	for i, addr := range addrs {
		toset[addr] = i
	}

	for _, bm := range lookupBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_Map_Lookup(b, toset, &bm.opt)
		})
	}
}

func Benchmark_Bloom_Add(b *testing.B) {
	numAddrs := len(addrs)

	b.ResetTimer()
	for k := 0; k < b.N/numAddrs; k++ {
		var bloom types.Bloom
		for _, addr := range addrs {
			bloom.Add(new(big.Int).SetBytes([]byte(addr)))
		}
	}
}

func benchmark_Bloom_Lookup(b *testing.B, bloom types.Bloom, opt *testOption) {
	numAddrs := len(opt.testAddrs)

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		addr := opt.testAddrs[k%numAddrs]
		if bloom.TestBytes([]byte(addr)) != opt.expected {
			b.Error(opt.msg, addr, "to test true")
		}
	}
}

func Benchmark_Bloom_Lookup(b *testing.B) {
	var bloom types.Bloom
	for _, addr := range addrs {
		bloom.Add(new(big.Int).SetBytes([]byte(addr)))
	}

	for _, bm := range lookupBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_Bloom_Lookup(b, bloom, &bm.opt)
		})
	}
}

func benchmark_Bloom_Compare(b *testing.B, bloom1 types.Bloom, opt *testOption) {
	var bloom2 types.Bloom
	for _, addr := range opt.testAddrs {
		bloom2.Add(new(big.Int).SetBytes([]byte(addr)))
	}

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		if (bloom1 == bloom2) != opt.expected {
			b.Error("expected", bloom2.Big(), opt.msg, bloom1.Big())
		}
	}
}

func Benchmark_Bloom_Compare(b *testing.B) {
	var bloom1 types.Bloom
	for _, addr := range addrs {
		bloom1.Add(new(big.Int).SetBytes([]byte(addr)))
	}

	for _, bm := range compareBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_Bloom_Compare(b, bloom1, &bm.opt)
		})
	}
}

func benchmark_Bloom_BitAND(b *testing.B, bloom1 types.Bloom, opt *testOption) {
	var bloom2 types.Bloom
	if len(opt.testAddrs) > 0 {
		bloom2.Add(new(big.Int).SetBytes([]byte(opt.testAddrs[0])))
	}

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		result := bloomBitAND(bloom1, bloom2)
		if (result == bloom2) != opt.expected {
			b.Error(opt.msg, bloom2, "to be part of", bloom1)
		}
	}
}

func Benchmark_Bloom_BitAND(b *testing.B) {
	var bloom1 types.Bloom
	for _, addr := range addrs {
		bloom1.Add(new(big.Int).SetBytes([]byte(addr)))
	}

	for _, bm := range bitandBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_Bloom_BitAND(b, bloom1, &bm.opt)
		})
	}
}

func Benchmark_Bloom_BitAND_TrueNaive(b *testing.B) {
	var bloom1 types.Bloom
	var bloom2 types.Bloom

	for _, addr := range addrs {
		bloom1.Add(new(big.Int).SetBytes([]byte(addr)))
	}

	if len(addrs) > 0 {
		bloom2.Add(new(big.Int).SetBytes([]byte(addrs[0])))
	}

	bloomAnd := new(big.Int)

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		bloom1Big := bloom1.Big()
		bloom2Big := bloom2.Big()
		if bloomAnd.And(bloom1Big, bloom2Big).Cmp(bloom2Big) != 0 {
			b.Error("expected", bloom2Big, "to be part of", bloom1Big)
		}
	}
}

func Benchmark_LDBBloom_Add(b *testing.B) {
	numAddrs := len(addrs)

	b.ResetTimer()
	for k := 0; k < b.N/numAddrs; k++ {
		bloomFilter := filter.NewBloomFilter(bitsPerKey)
		generator := bloomFilter.NewGenerator()
		for _, addr := range addrs {
			generator.Add([]byte(addr))
		}
		buf := &util.Buffer{}
		generator.Generate(buf)
	}
}

func benchmark_LDBBloom_Lookup(b *testing.B, bloomFilter filter.Filter, filterInstance []byte, opt *testOption) {
	numAddrs := len(opt.testAddrs)

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		addr := opt.testAddrs[k%numAddrs]
		if bloomFilter.Contains(filterInstance, []byte(addr)) != opt.expected {
			b.Error(opt.msg, addr, "to test true")
		}
	}
}

func Benchmark_LDBBloom_Lookup(b *testing.B) {
	bloomFilter := filter.NewBloomFilter(bitsPerKey)
	generator := bloomFilter.NewGenerator()
	for _, addr := range addrs {
		generator.Add([]byte(addr))
	}
	buf := &util.Buffer{}
	generator.Generate(buf)
	filterInstance := buf.Bytes()

	for _, bm := range lookupBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_LDBBloom_Lookup(b, bloomFilter, filterInstance, &bm.opt)
		})
	}
}

func benchmark_LDBBloom_Compare(b *testing.B, bloomFilter filter.Filter, filterInstance1 []byte, opt *testOption) {
	numAddrs := len(opt.testAddrs)
	generator2 := bloomFilter.NewGenerator()
	for i := 0; i < len(addrs); i++ {
		generator2.Add([]byte(opt.testAddrs[i%numAddrs]))
	}
	buf2 := &util.Buffer{}
	generator2.Generate(buf2)
	filterInstance2 := buf2.Bytes()

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		if bytes.Equal(filterInstance1, filterInstance2) != opt.expected {
			b.Error("expected", filterInstance2, opt.msg, filterInstance1)
		}
	}
}

func Benchmark_LDBBloom_Compare(b *testing.B) {
	bloomFilter := filter.NewBloomFilter(bitsPerKey)

	generator1 := bloomFilter.NewGenerator()
	for _, addr := range addrs {
		generator1.Add([]byte(addr))
	}
	buf1 := &util.Buffer{}
	generator1.Generate(buf1)
	filterInstance1 := buf1.Bytes()

	for _, bm := range compareBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_LDBBloom_Compare(b, bloomFilter, filterInstance1, &bm.opt)
		})
	}
}

func benchmark_LDBBloom_BitAND(b *testing.B, bloomFilter filter.Filter, filterInstance1 []byte, opt *testOption) {
	generator2 := bloomFilter.NewGenerator()
	for i := 0; i < len(addrs); i++ {
		generator2.Add([]byte(opt.testAddrs[0]))
	}
	buf2 := &util.Buffer{}
	generator2.Generate(buf2)
	filterInstance2 := buf2.Bytes()

	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		result := ldbBloomBitAND(filterInstance1, filterInstance2)
		if bytes.Equal(result, filterInstance2) != opt.expected {
			b.Error(opt.msg, filterInstance2, "to be part of", filterInstance1)
		}
	}
}

func Benchmark_LDBBloom_BitAND(b *testing.B) {
	bloomFilter := filter.NewBloomFilter(bitsPerKey)

	generator1 := bloomFilter.NewGenerator()
	for _, addr := range addrs {
		generator1.Add([]byte(addr))
	}
	buf1 := &util.Buffer{}
	generator1.Generate(buf1)
	filterInstance1 := buf1.Bytes()

	for _, bm := range bitandBenchmarks {
		b.Run(bm.name, func(b *testing.B) {
			benchmark_LDBBloom_BitAND(b, bloomFilter, filterInstance1, &bm.opt)
		})
	}
}
