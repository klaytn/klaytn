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

package account

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// Compile time interface checks
var (
	_ Account = (*LegacyAccount)(nil)
	_ Account = (*ExternallyOwnedAccount)(nil)
	_ Account = (*SmartContractAccount)(nil)

	_ ProgramAccount = (*SmartContractAccount)(nil)

	_ AccountWithKey = (*ExternallyOwnedAccount)(nil)
	_ AccountWithKey = (*SmartContractAccount)(nil)
)

// TestAccountSerialization tests serialization of various account types.
func TestAccountSerialization(t *testing.T) {
	accs := []struct {
		Name string
		acc  Account
	}{
		{"EOA", genEOA()},
		{"EOAWithPublic", genEOAWithPublicKey()},
		{"SCA", genSCA()},
		{"SCAWithPublic", genSCAWithPublicKey()},
	}
	testcases := []struct {
		Name string
		fn   func(t *testing.T, acc Account)
	}{
		{"RLP", testAccountRLP},
		{"JSON", testAccountJSON},
	}
	for _, test := range testcases {
		for _, acc := range accs {
			Name := test.Name + "/" + acc.Name
			t.Run(Name, func(t *testing.T) {
				test.fn(t, acc.acc)
			})
		}
	}
}

func testAccountRLP(t *testing.T, acc Account) {
	enc := NewAccountSerializerWithAccount(acc)

	b, err := rlp.EncodeToBytes(enc)
	if err != nil {
		panic(err)
	}

	dec := NewAccountSerializer()

	if err := rlp.DecodeBytes(b, &dec); err != nil {
		panic(err)
	}

	if !acc.Equal(dec.account) {
		fmt.Println("acc")
		fmt.Println(acc)
		fmt.Println("dec.account")
		fmt.Println(dec.account)
		t.Errorf("acc != dec.account")
	}
}

func testAccountJSON(t *testing.T, acc Account) {
	enc := NewAccountSerializerWithAccount(acc)

	b, err := json.Marshal(enc)
	if err != nil {
		panic(err)
	}

	dec := NewAccountSerializer()

	if err := json.Unmarshal(b, &dec); err != nil {
		panic(err)
	}

	if !acc.Equal(dec.account) {
		fmt.Println("acc")
		fmt.Println(acc)
		fmt.Println("dec.account")
		fmt.Println(dec.account)
		t.Errorf("acc != dec.account")
	}
}

func genRandomHash() (h common.Hash) {
	hasher := sha3.NewKeccak256()

	r := rand.Uint64()
	rlp.Encode(hasher, r)
	hasher.Sum(h[:0])

	return h
}

func genEOA() *ExternallyOwnedAccount {
	humanReadable := false

	return newExternallyOwnedAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyLegacy(),
	})
}

func genEOAWithPublicKey() *ExternallyOwnedAccount {
	humanReadable := false

	k, _ := crypto.GenerateKey()

	return newExternallyOwnedAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyPublicWithValue(&k.PublicKey),
	})
}

func genSCA() *SmartContractAccount {
	humanReadable := false

	return newSmartContractAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyLegacy(),
		AccountValueKeyStorageRoot:   genRandomHash(),
		AccountValueKeyCodeHash:      genRandomHash().Bytes(),
		AccountValueKeyCodeInfo:      params.CodeInfo(0),
	})
}

func genSCAWithPublicKey() *SmartContractAccount {
	humanReadable := false

	k, _ := crypto.GenerateKey()

	return newSmartContractAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyPublicWithValue(&k.PublicKey),
		AccountValueKeyStorageRoot:   genRandomHash(),
		AccountValueKeyCodeHash:      genRandomHash().Bytes(),
		AccountValueKeyCodeInfo:      params.CodeInfo(0),
	})
}

// Tests RLP encoding against manually generated strings.
func TestSmartContractAccountExt(t *testing.T) {
	// To create testcases,
	// - Install https://github.com/ethereumjs/rlp
	//     npm install -g rlp
	// - In bash, run
	//     maketc(){ echo $(rlp encode "$1")$(rlp encode "$2" | cut -b3-); }
	//     maketc 2 '["0x1234","0x5678"]'
	var (
		commonFields = &AccountCommon{
			nonce:         0x1234,
			balance:       big.NewInt(0x5678),
			humanReadable: false,
			key:           accountkey.NewAccountKeyLegacy(),
		}
		hash     = common.HexToHash("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff")
		exthash  = common.HexToExtHash("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01")
		codehash = common.HexToHash("aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd").Bytes()
		codeinfo = params.CodeInfo(0x10)

		// StorageRoot is hash32:  maketc 2 '[["0x1234","0x5678","","0x01",[]],"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","0xaaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd","0x10"]'
		scaLegacyRLP = "0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10"
		scaLegacy    = &SmartContractAccount{
			AccountCommon: commonFields,
			storageRoot:   hash.ExtendLegacy(),
			codeHash:      codehash,
			codeInfo:      codeinfo,
		}
		// StorageRoot is exthash: maketc 2 '[["0x1234","0x5678","","0x01",[]],"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01","0xaaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd","0x10"]'
		scaExtRLP = "0x02f854c98212348256788001c0a700112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01a0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10"
		scaExt    = &SmartContractAccount{
			AccountCommon: commonFields,
			storageRoot:   exthash,
			codeHash:      codehash,
			codeInfo:      codeinfo,
		}
	)
	checkEncode := func(account Account, encoded string) {
		enc := NewAccountSerializerWithAccount(account)
		b, err := rlp.EncodeToBytes(enc)
		assert.Nil(t, err)
		assert.Equal(t, encoded, hexutil.Encode(b))
	}
	checkEncodeExt := func(account Account, encoded string) {
		enc := NewAccountSerializerExtWithAccount(account)
		b, err := rlp.EncodeToBytes(enc)
		assert.Nil(t, err)
		assert.Equal(t, encoded, hexutil.Encode(b))
	}
	checkEncode(scaLegacy, scaLegacyRLP)
	checkEncodeExt(scaLegacy, scaLegacyRLP) // LegacyExt are always unextended

	checkEncode(scaExt, scaLegacyRLP) // Regular encoding still results in hash32. Use it for merkle hash.
	checkEncodeExt(scaExt, scaExtRLP) // Must use SerializeExt to preserve exthash. Use it for disk storage.

	checkDecode := func(encoded string, account Account) {
		b := common.FromHex(encoded)
		dec := NewAccountSerializer()
		err := rlp.DecodeBytes(b, &dec)
		assert.Nil(t, err)
		assert.True(t, dec.GetAccount().Equal(account))
	}
	checkDecodeExt := func(encoded string, account Account) {
		b := common.FromHex(encoded)
		dec := NewAccountSerializerExt()
		err := rlp.DecodeBytes(b, &dec)
		assert.Nil(t, err)
		assert.True(t, dec.GetAccount().Equal(account))
	}

	checkDecode(scaLegacyRLP, scaLegacy)
	checkDecodeExt(scaLegacyRLP, scaLegacy)

	checkDecode(scaExtRLP, scaExt)
	checkDecodeExt(scaExtRLP, scaExt)
}

func TestUnextendRLP(t *testing.T) {
	// storage slot
	testcases := []struct {
		extended   string
		unextended string
	}{
		{ // storage slot (33B) kept as-is
			"0xa06700000000000000000000000000000000000000000000000000000000000000",
			"0xa06700000000000000000000000000000000000000000000000000000000000000",
		},
		{ // Short EOA (<=ExtHashLength) kept as-is
			"0x01c98212348256788001c0",
			"0x01c98212348256788001c0",
		},
		{ // Long EOA (>ExtHashLength) kept as-is
			"0x01ea8212348256788002a1030bc77753515dd61c66df6445ffffbedfc16b6b46c73eb09f01a970cb3bf0a8de",
			"0x01ea8212348256788002a1030bc77753515dd61c66df6445ffffbedfc16b6b46c73eb09f01a970cb3bf0a8de",
		},
		{ // SCA with Hash keps as-is
			"0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			"0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
		},
		{ // SCA with ExtHash unextended
			"0x02f854c98212348256788001c0a700112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01a0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			"0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
		},
		{ // Long malform data kept as-is
			"0xdead0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"0xdead0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{ // Short malformed data kept as-is
			"0x80",
			"0x80",
		},
		{ // Short malformed data kept as-is
			"0x00",
			"0x00",
		},
		{ // Legacy account may crash DecodeRLP, but must not crash UnextendSerializedAccount.
			"0xf8448080a00000000000000000000000000000000000000000000000000000000000000000a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
			"0xf8448080a00000000000000000000000000000000000000000000000000000000000000000a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
	}

	for _, tc := range testcases {
		unextended := UnextendSerializedAccount(common.FromHex(tc.extended))
		assert.Equal(t, tc.unextended, hexutil.Encode(unextended))
	}
}
