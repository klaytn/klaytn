// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from accounts/keystore/keystore_passphrase_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package keystore

import (
	"crypto/ecdsa"
	"io/ioutil"
	"testing"

	"github.com/klaytn/klaytn/crypto"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"

	"github.com/klaytn/klaytn/common"
)

const (
	veryLightScryptN = 2
	veryLightScryptP = 1
)

// Tests that a json key file can be decrypted and encrypted in multiple rounds.
func TestKeyEncryptDecrypt(t *testing.T) {
	keyjson, err := ioutil.ReadFile("testdata/very-light-scrypt.json")
	if err != nil {
		t.Fatal(err)
	}
	password := ""
	address := common.HexToAddress("45dea0fb0bba44f4fcf290bba71fd57d7117cbb8")

	// Do a few rounds of decryption and encryption
	for i := 0; i < 3; i++ {
		// Try a bad password first
		if _, err := DecryptKey(keyjson, password+"bad"); err == nil {
			t.Errorf("test %d: json key decrypted with bad password", i)
		}
		// Decrypt with the correct password
		key, err := DecryptKey(keyjson, password)
		if err != nil {
			t.Fatalf("test %d: json key failed to decrypt: %v", i, err)
		}
		if key.GetAddress() != address {
			t.Errorf("test %d: key address mismatch: have %x, want %x", i, key.GetAddress(), address)
		}
		// Recrypt with a new password and start over
		password += "new data appended"
		if keyjson, err = EncryptKeyV3(key, password, veryLightScryptN, veryLightScryptP); err != nil {
			t.Errorf("test %d: failed to recrypt key %v", i, err)
		}
	}
}

// Tests encoding and decoding of a keystore v4 object with a single private key.
func TestEncryptDecryptV4Single(t *testing.T) {
	pk, err := crypto.GenerateKey()
	require.NoError(t, err)

	auth := string("password")

	key := &KeyV4{
		Id:          uuid.NewRandom(),
		Address:     crypto.PubkeyToAddress(pk.PublicKey),
		PrivateKeys: [][]*ecdsa.PrivateKey{{pk}},
	}

	keyjson, err := EncryptKey(key, auth, 4096, 1)
	require.NoError(t, err)

	k, err := DecryptKey(keyjson, "password")

	require.Equal(t, key, k)
}

// Tests encoding and decoding of a keystore v4 object with multiple private keys.
func TestEncryptDecryptV4Multiple(t *testing.T) {
	pks := make([]*ecdsa.PrivateKey, 10)

	for i := range pks {
		k, err := crypto.GenerateKey()
		require.NoError(t, err)
		pks[i] = k
	}

	auth := "password"

	key := &KeyV4{
		Id:          uuid.NewRandom(),
		Address:     crypto.PubkeyToAddress(pks[0].PublicKey),
		PrivateKeys: [][]*ecdsa.PrivateKey{pks},
	}

	keyjson, err := EncryptKey(key, auth, 4096, 1)
	require.NoError(t, err)

	k, err := DecryptKey(keyjson, "password")
	require.NoError(t, err)

	require.Equal(t, key, k)
}

// Tests encoding and decoding of a keystore v4 object with role-based keys.
func TestEncryptDecryptV4RoleBased(t *testing.T) {
	pks := make([][]*ecdsa.PrivateKey, 3)

	for i := range pks {
		pks[i] = make([]*ecdsa.PrivateKey, 10)
		for j := range pks[i] {
			k, err := crypto.GenerateKey()
			require.NoError(t, err)
			pks[i][j] = k
		}
	}

	auth := "password"

	key := &KeyV4{
		Id:          uuid.NewRandom(),
		Address:     crypto.PubkeyToAddress(pks[0][0].PublicKey),
		PrivateKeys: pks,
	}

	keyjson, err := EncryptKey(key, auth, 4096, 1)
	require.NoError(t, err)

	k, err := DecryptKey(keyjson, auth)
	require.NoError(t, err)

	require.Equal(t, key, k)
}

// Tests decoding of a hard-coded keystore v4 JSON object storing a private key.
func TestKeyDecryptV4Single(t *testing.T) {
	keyjson := []byte(`{
    "version": 4,
	"id": "7a0a8557-22a5-4c90-b554-d6f3b13783ea",
	"address": "0x86bce8c859f5f304aa30adb89f2f7b6ee5a0d6e2",
    "keyring": [
		{
			"ciphertext": "696d0e8e8bd21ff1f82f7c87b6964f0f17f8bfbd52141069b59f084555f277b7",
			"cipherparams": { "iv": "1fd13e0524fa1095c5f80627f1d24cbd" },
			"cipher": "aes-128-ctr",
			"kdf": "scrypt",
			"kdfparams": {
				"dklen": 32,
				"salt": "7ee980925cef6a60553cda3e91cb8e3c62733f64579f633d0f86ce050c151e26",
				"n": 4096,
				"r": 8,
				"p": 1
			},
			"mac": "8684d8dc4bf17318cd46c85dbd9a9ec5d9b290e04d78d4f6b5be9c413ff30ea4"
		}
    ]
}`)

	k, err := DecryptKey(keyjson, "password")
	require.NoError(t, err)

	require.Equal(t, common.HexToAddress("0x86bce8c859f5f304aa30adb89f2f7b6ee5a0d6e2"), k.GetAddress())
	key, err := crypto.ToECDSA(common.Hex2Bytes("36e0a792553f94a7660e5484cfc8367e7d56a383261175b9abced7416a5d87df"))
	require.NoError(t, err)

	require.Equal(t, key, k.GetPrivateKeys()[0][0])
}

// Tests decoding of a hard-coded keystore v4 JSON object storing multiple private keys.
func TestKeyDecryptV4Multiple(t *testing.T) {
	keyjson := []byte(
		`
{
	"version": 4,
	"id": "55da3f9c-6444-4fc1-abfa-f2eabfc57501",
	"address": "0x86bce8c859f5f304aa30adb89f2f7b6ee5a0d6e2",
	"keyring": [
		{
			"ciphertext": "93dd2c777abd9b80a0be8e1eb9739cbf27c127621a5d3f81e7779e47d3bb22f6",
			"cipherparams": { "iv": "84f90907f3f54f53d19cbd6ae1496b86" },
			"cipher": "aes-128-ctr",
			"kdf": "scrypt",
			"kdfparams": {
				"dklen": 32,
				"salt": "69bf176a136c67a39d131912fb1e0ada4be0ed9f882448e1557b5c4233006e10",
				"n": 4096,
				"r": 8,
				"p": 1
			},
			"mac": "8f6d1d234f4a87162cf3de0c7fb1d4a8421cd8f5a97b86b1a8e576ffc1eb52d2"
		},
		{
			"ciphertext": "53d50b4e86b550b26919d9b8cea762cd3c637dfe4f2a0f18995d3401ead839a6",
			"cipherparams": { "iv": "d7a6f63558996a9f99e7daabd289aa2c" },
			"cipher": "aes-128-ctr",
			"kdf": "scrypt",
			"kdfparams": {
				"dklen": 32,
				"salt": "966116898d90c3e53ea09e4850a71e16df9533c1f9e1b2e1a9edec781e1ad44f",
				"n": 4096,
				"r": 8,
				"p": 1
			},
			"mac": "bca7125e17565c672a110ace9a25755847d42b81aa7df4bb8f5ce01ef7213295"
		}
	]
}
`)

	k, err := DecryptKey(keyjson, "password")
	require.NoError(t, err)

	require.Equal(t, common.HexToAddress("0x86bce8c859f5f304aa30adb89f2f7b6ee5a0d6e2"), k.GetAddress())
	require.Equal(t, k.GetPrivateKeys(), [][]*ecdsa.PrivateKey{
		{
			crypto.ToECDSAUnsafe(common.Hex2Bytes("d1e9f8f00ef9f93365f5eabccccb3f3c5783001b61a40f0f74270e50158c163d")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("4bd8d0b0c1575a7a35915f9af3ef8beb11ad571337ec9b6aca7c88ca7458ef5c")),
		},
	})
}

// Tests decoding of a hard-coded keystore v4 JSON object storing role-based keys.
func TestKeyDecryptV4RoleBased(t *testing.T) {
	keyjson := []byte(
		`
{
  "address": "5bf1a459fcac4dae0910c420d9c4643a89855c4f",
  "keyring": [
    [
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "deed1b95fca6f6b780e938333b16524a9bfaef9ef46a4301db9eb1b53a2110d7",
        "cipherparams": {
          "iv": "ab043dbfc0436b58cba0ec8b7da5ab27"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "3d34e949f5aaaf2e18681d9830c09c2e3c777ba6d715caf8ecef66e2ebe44107"
        },
        "mac": "d9be9c5d91ce9e966735f6e472561bf2b06fb58dea4ae91439fa88265e69b0ae"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "ad1cb5c41160f92561bdd426c0db93dcfeba9e4dedf4431cb57a0e04788d0de6",
        "cipherparams": {
          "iv": "14fae90baf0c62fa18dfbf99b8025ca4"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "045e100b13f78567476c1ebbdb851c68536e641c5f2ed778f8d6d9e2c2ba8ef3"
        },
        "mac": "a4a3abba6b9b8f6bd8d72034e0fc3c5a1af7382851f317381ec1a5e12b8e4b7d"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "c25cdb53ff80fafabc27019f160a230b1a08534a350fdb8cb15909719df24d0f",
        "cipherparams": {
          "iv": "c3825ea9ef8bc478110ad0d383876ecb"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "4daaea8d0ac83fad1a48f36132c38407aec215b492c3813a2c5639808cad5c42"
        },
        "mac": "77224060cfdadeb98dfd7239816a70a3e6f74eec02d027041bf6eac0303b4528"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "b2210f7d7c9062f75b043081cc04d4556ed562f9efa8b75b0778bee53782dcae",
        "cipherparams": {
          "iv": "0a5c983343d3e7f6039a5689029dd57f"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "1f3a5fd0d9d62313ab354b341f300758f8fc6393ab244b50b1873c5d0ee176dc"
        },
        "mac": "c5e0ad4f33e05db790b74fa21fa3204c3f2b1ccc1bef1fb08a9d6c1f1a4cdb08"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "0b4a310f71a26469d5e4a2b3609cf248f1e7e54dedfb6f08d37f817286c9d84e",
        "cipherparams": {
          "iv": "aad04a9dc73250291ec1dae634420596"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "185f65cfbbd835cc06300f1f3dd34b2a6175ff8080f249759d98b3095dd27119"
        },
        "mac": "7d7a74dc791ac7ded99574940b37411f36eb24fdd3f6bf1c82133f7ced8984ea"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "8e50ef3fd6fa53a22e7d10a868935593b36a17e9cba4dc16722fb15c48b55023",
        "cipherparams": {
          "iv": "f56d0975e48b4a3936b8f87d8556c483"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "d0a0a1f4b172ba78b86f6193b8863d0a3e792195894883099980c03ed70893b7"
        },
        "mac": "091f83a526afce5287141b012fadc20dad93de8f834663facde6b54f7d64efa4"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "73966e70348072d59b2c9b4574a57b2df69de70afc2a892936251b898869c640",
        "cipherparams": {
          "iv": "180ce870ed381eaefae79210675d0531"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "f70885b79f5c90963848d46a62de378852d94bb6786709872cfac6bafc7b9ce2"
        },
        "mac": "0da0e9ddc7004311356654efb358dbd62186b30ffc97e0e672c6938ec793457a"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "bdee7419c58744b5b15d3d92666f20340d29e6ba716735b75a43c87b61d71d76",
        "cipherparams": {
          "iv": "d341bc56c431b105cdd1654bc822435e"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "1d9c3d932197ad3632ca07aaea04aa65f9347e903e310345d4741b52b6c28066"
        },
        "mac": "451eaf374f7c39ea3e063081b854e9a9b196def73d07e3bb9d41b4c919f4ebc2"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "8319770114706fc99ec5d8ba0e0f8762eab974d8ebe64ad3ea9896b0b08947dd",
        "cipherparams": {
          "iv": "d0796a9ee6eac10526035c2abfe1e179"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "9e019be17ac9caab6868b6e6ea8bf57941ef107bb4a43937c46222474cdaff37"
        },
        "mac": "e00630bee72f16bf67585d68cc5e77b8405b661aa523ec0212285938e3454332"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "a18a95a97a51ed6062643f16c6e40c210c266d3822e794bf615b1ac422ce9d11",
        "cipherparams": {
          "iv": "39d922c4546c29d2c1aa8a5143f509af"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "2fd3874e028c323cdfe15d6c21060befb1ba6f6e38bf51059f48f301dbe024ab"
        },
        "mac": "4c4ee1929b17d1317fb5cd8724b76065b581abab74e155d30e26a6996504d946"
      }
    ],
    [
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "d2b634d1c215bef632783883e6d9920784df64a0debf72987fcbdc816abe3da2",
        "cipherparams": {
          "iv": "2ba6644ba2c2b912e7dfebdfa9086f4d"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "8d3cab2cf30badb153944d0a2d0fd1214ffeed2738082f6acd76ed2c2e048b44"
        },
        "mac": "28e025d8a7117ecef98150fb0aef4ead8c8b54815b1820ffc8a04f69ea82e277"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "a069d776f9c6b433f87c30d3284936a97a82d6204fad1adeb7e1f65c1c89aa1c",
        "cipherparams": {
          "iv": "edb7c6159f5d9f142fe2f3a17c326c0c"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "5bb28deccadfd207061d8a7b82788a9483e4d5f504aa8c4197190f971b4fea10"
        },
        "mac": "418f0a887920ffd8a01f647d55f25f45ea0ef53b2977411b45dbd1a178c690f7"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "78ecf565643f0bb57c818242731f25ca333fc6a80c0f6e4889835fd25f79fdf8",
        "cipherparams": {
          "iv": "08776afa6efb9ae8324d728658b75958"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "7d51d79e0c960f562582ad657cb778225159b2930e08a4ed88402596d71e974a"
        },
        "mac": "3fbe366bea655afdb42a05c7a3c06fd5b52eabbfb5cf5fbbfcb3bb42d6e2bee3"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "9b3c5f7378095d2fa7f33986d409f28a1ef2f9b85c52a08ad8f70a36ccaf4b4f",
        "cipherparams": {
          "iv": "b44a3d369f9a9016f61f9c7510ee7027"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "2422aad4295f8949cd9b48728e0595fa8289b7f9faa3f57d8a2c63bfa8dafa81"
        },
        "mac": "1f85b0bf8dffe75fbdde00e9ac0a63458945281888a43cca63d65e3af607b3f1"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "90843bc47482a39046a6666fc0ed8973d033a156ee7ff1c28e91604665f0fb42",
        "cipherparams": {
          "iv": "792197f7f310b609154401723bf630c9"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "b6c5eef485a1f568d2ab1f3a234da5a9ad0592c30341372c9ece08e07db276e0"
        },
        "mac": "4a7e6c6f7fa22accd2fe8f87e8150e28419c12039fe82ad966195a4110fc3546"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "55a4902bc4aaed432628eabe10833169ce943bb51c2d528b2a861730c113fffd",
        "cipherparams": {
          "iv": "c2be7a4fbdd3fdd01ac562875e1da297"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "5e933ddb6018844356889e56edc56b5e9179b388fed67b919d3421251faa955c"
        },
        "mac": "9545baedfa7d2044e9d894b82bb32ea54a8f028dee5645c2673faba1574f0d71"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "9bb56fbaf51422cfd2fcf6b18fa6e42980c26280f01c6443879f7f032cc64b28",
        "cipherparams": {
          "iv": "98566dbe29712239e4dc28e225d85552"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "8b57fcfd31b1d4db459b9c79f4e617848c8744abb3de4686d7b881420bddaf9b"
        },
        "mac": "dcd4254e4a7136c0df4367b9411b52ebc1e35e3d82059fb708981518bdb37f91"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "4113ab4bc8c8a40ee013c2ec9831f6bc5e1b8468366bab56ce56bd17fddc0c7b",
        "cipherparams": {
          "iv": "6701ace43092aa749d2cfe2248a314bf"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "53f63657a5213aa416d35a8f082aef0fd8e398af068b9e6f45d893efb5df309a"
        },
        "mac": "70d15a4afcb1882e85de547cd709b3cf160ea8b4350aa6e413b43189f2a4fffd"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "3c9023e0a4023e06c466ce785269449265b7a6bb0e5eb32a2cbb1242d67b73ab",
        "cipherparams": {
          "iv": "a3761b149738e168b43bc6dbd6a5a9f3"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "ca8408e34bfff90c062bb0ab8959b8cdb25871fb5c13aa16b4e9f9917137adc3"
        },
        "mac": "29283d50e91268d365ca5f178fa73d3d7bb9ee34ce81dad3ce15929f665fbe96"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "8e2a0e9297165abf0db979af50fd01297c64b97673d0c049def685ed05e20cd7",
        "cipherparams": {
          "iv": "2e64c9eb8a42992dd194285b490eee64"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "86f569d5dd84b058ed5a618b33c14b723d446162af2b23c959249b0b2df491b6"
        },
        "mac": "6d7212972ded83c5128ad0b828189c95c77ac4c7784edf403d77688908797f45"
      }
    ],
    [
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "666e90f61a8d9d8b612961cf1f5b44ed718343b5b60778c51d95a53d3b70bfd9",
        "cipherparams": {
          "iv": "ce3f2076a56d29e606511272d54bb373"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "aad67bd8f3960f585e03067fd806bd23f531acbfe8ade4e01a0866952da016a3"
        },
        "mac": "29ab9cd731bb3df2401db1ff5db1bcfaf718a66c5dd5d0a62d2591455b806937"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "091cad02f7c0c45e297288cbdf4c60f7ca7914ec619e027f39c2800fce342948",
        "cipherparams": {
          "iv": "8f3f7c1e08a8d7281ab92a4c2816f7d1"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "38ee78865647b149c9d5eda6554919e32babf25badf2dc36638f87d25e0ba741"
        },
        "mac": "792f6fe079e2a0d7b2314bf136dac3c295f59d097f9972053539aa3ed63334ee"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "846df61cf31705a256d64c2544a5f3e604c73f02a3bed2c5405f0fba5234ce4f",
        "cipherparams": {
          "iv": "da45e5028fc8fe8dc80fdedbc7529d67"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "b1a342841710d91165e25dbd0bd2aaeb84d6ee7ceeacaf7401f6691b50ae03b8"
        },
        "mac": "6c836be73c01eb8f9bb0c244a9bf693dc603680f02f5eddeb1f638a073a3b97e"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "517373acafb851df6eedc4d419d0968e08184dad31d13617f5c4f2a76ee0b325",
        "cipherparams": {
          "iv": "a75dca682365c8be0d486549afe3b7b2"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "bce019b9dccc24a1cd2a612d7b24d9b0113d561c0f69f9579873632503809db6"
        },
        "mac": "4f5165f71a41b04d12efabc842cb26815567e091dc839c0822098feb2ccd460b"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "8900c93b208c50a6606bc4d3cb70fbbd7e2082ac20a50dad2fe522a5f25be0e4",
        "cipherparams": {
          "iv": "e7291589cb96f1b0aa535e7e44af00f1"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "c8303cbff39201327fbbe3a7c51e404a877d0656c363a8fa7b9708340b67b8fa"
        },
        "mac": "31fe42941a9c4167d09214b57613db9844d0a98a85443bb1b6c2456d7c91d61a"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "7fabb89609d7df7fd4175987076e06b29a423be259686d3084874cf957b323f3",
        "cipherparams": {
          "iv": "2cf25f05f2138177342002836475dd22"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "9eb81363707f3664718e635007f2335fd4167670bdb3c9eece4c0cd1f245b4be"
        },
        "mac": "cb1477c23a67e6460bbd27d9859f073d9c7934f6bd138cdddd2399f2062190ac"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "903630b6226e94d12e6360d3f23c38d78f56771f5c78369ff320a53d1fa2f104",
        "cipherparams": {
          "iv": "c6d864616798267f07acd0fb96161289"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "c6164c7bb00390a24c205944ab533cb146b334703b02433542c93c5016b9f921"
        },
        "mac": "4345a92c4357a682c1156059b50d033119ecadb1cc1b2695347c6992bfe643c7"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "8a6b8854ac09d9cd59fd0037ff51180f9776005c1b22e49548af1fa5fe6d1e8a",
        "cipherparams": {
          "iv": "07827df9d1466f58aeb8dbb66c85e7fd"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "e34a7fdd70748293b5ae8eeb636209cb8cba5cd8f4ac8460f7da680d45dc4c06"
        },
        "mac": "97fe95ea956834fce2e15abd46c7e43c3ba8aa2feb1cbbe6534bc3b417b7a606"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "7822994b653b69148a87fc839cb5236c8d9895bd1c5988cf9eb99801efc7a489",
        "cipherparams": {
          "iv": "5343a4e8955132e00ba88732ad7eacd9"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "751ecdff9c96fee6be20850d5e6ce7f3f9518216b639a4872d550d1bbe9d6f47"
        },
        "mac": "7d37f4756b7ed25d453bd48b7ac66446a1b9da146f88925bfa02d8c0da6a512a"
      },
      {
        "cipher": "aes-128-ctr",
        "ciphertext": "a2a995a549e5a10e3b6f18ef536a02a13ea1ffb819de51b0cf73b77a6d037f01",
        "cipherparams": {
          "iv": "46448adacdd9364cbc5472fc754a310c"
        },
        "kdf": "scrypt",
        "kdfparams": {
          "dklen": 32,
          "n": 4096,
          "p": 1,
          "r": 8,
          "salt": "52009f694f5f387cd6daab8aabea4a4a83a68a4ac7610faab06b6cc24f13498e"
        },
        "mac": "353c4e2fe08482e35d1ec5191ddcca1884886e2842e01bc78836bd5636c14a0c"
      }
    ]
  ],
  "id": "51be2d49-f4f3-4215-b7c3-7e3a31191fb2",
  "version": 4
}
`)

	k, err := DecryptKey(keyjson, "password")
	require.NoError(t, err)

	require.Equal(t, common.HexToAddress("0x5bf1a459fcac4dae0910c420d9c4643a89855c4f"), k.GetAddress())
	require.Equal(t, k.GetPrivateKeys(), [][]*ecdsa.PrivateKey{
		{
			crypto.ToECDSAUnsafe(common.Hex2Bytes("402b411c437893748081751a27ed310b71b4da49770d28338f020c4663adf720")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("cf9e08f724dd666aaf38e0aae689466b2b1d15e731aa7b80aacff57b294eb634")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("8f5505154a6b5662892dd890b66843da6a4ee2941db74224c80b6064532d6623")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("a6952ea0ac137efda3474460a5fbf9a6e24258f7e0e43325661e62637e30eb15")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("69ae72dda34c0d83f8eb8729c80f3c13dbf4db27c94b3f8dba40d7e167538f3f")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("ae45c060a174b652ecf2bcc9b5d8b937554902c94803c5d3e35444b2a6153268")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("0bb459b227b82e19eb5865bec204a125996de314d52b78eca86e2f448c83e432")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("775685cfda8c4ecba96eab877c206e53a21b736c9df4a55eb045dcb90d8d7498")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("34267a26113e165fccefc317ca364e4cc3413eb2cac2412a68b191f22e057195")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("01d6cbfcf6e560eb853759002b16b698bdcca1c9d1941a03845fc088b5967811")),
		},
		{
			crypto.ToECDSAUnsafe(common.Hex2Bytes("ceb6db41a55290dcac1fcbd82a31e6916a80beb65fd681227a61c8fbf92a0305")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("cdb5b8f88adbbe7b42e7d9747d27916b29b92fd44b1f9974ea65546e7c977a73")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("f2f80b09e34bdec35663c43dc2926684c89c1a0f8e59b8ae116f114e30da71f8")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("7e88d048e2b37131ad050b44b2716309763fa31d310adc818e57e9da0cabd6c0")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("6a7757f6af99910275e98e847c9fe8989dccecaf4b9f72176459d1d2503e758f")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("d98dd9a175e92e6b5a975085a81c480516aaae70c3ed04603cb4b62a06e5ef29")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("a4edb7600e7609a666de0cedd89cdb6e374e7c8e26effa89d6b7d28d66596179")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("5037c508aa70ae3a2f6808c136e1aa65dc33591ce3bdb26a87a396bb9e2119f3")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("8b639d5d2108d6217e94503b74427cd35dd72d06698072dbaa460f8a05cf391b")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("daafea5b869f0f6519001697fcafdfd3ea22e429e7e9bf14862a8995a953fe98")),
		},
		{
			crypto.ToECDSAUnsafe(common.Hex2Bytes("e80f49e9103ff90c825e1a9cfc57a93fa0c9ba8fc75f65b2a146763ee15ae46e")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("39c39da79dc642e7476dc442b39473e72d24aae67075240460f245c9b541a686")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("17540767a9de215153ec230e8027ce8c23bf399fba163a71cbe0e8b9600b9ab1")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("d20f6b99f838221edb8e7651fad2bb3f6c7c07321976b636f5c7dba57b921ff1")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("5939e2b6afff1060a28db9f98034c576814b9ae9f9d1f71df2aebfe365d660e7")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("e3e763873c0dd23b1310c0f57228c8c9adf2db1097d8b53d8b2f16cec86b5c3e")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("964af45ee7f833403cc24696dc06574ded7af5f7efbb1430db996960d0746a4b")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("17433d30efc4bd7bd3131e3c46b6419fd77ec078bff902c878eb8bab00b089cc")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("49e27a5bb241b8673fe13c72bcf914ee335e1eb526e4ef3d335af7de28b28a18")),
			crypto.ToECDSAUnsafe(common.Hex2Bytes("13855eb62d3a425fe5f224975dcdcc561f782ffdc76df3b0d14c7d980f1348bf")),
		},
	})
}

// Tests decoding of a hard-coded keystore v3 JSON object storing a private key.
func TestKeyDecryptV3(t *testing.T) {
	keyjson := []byte(`{
    "version": 3,
	"id": "7a0a8557-22a5-4c90-b554-d6f3b13783ea",
	"address": "0x86bce8c859f5f304aa30adb89f2f7b6ee5a0d6e2",
    "crypto": {
		"ciphertext": "696d0e8e8bd21ff1f82f7c87b6964f0f17f8bfbd52141069b59f084555f277b7",
		"cipherparams": { "iv": "1fd13e0524fa1095c5f80627f1d24cbd" },
		"cipher": "aes-128-ctr",
		"kdf": "scrypt",
		"kdfparams": {
			"dklen": 32,
			"salt": "7ee980925cef6a60553cda3e91cb8e3c62733f64579f633d0f86ce050c151e26",
			"n": 4096,
			"r": 8,
			"p": 1
		},
		"mac": "8684d8dc4bf17318cd46c85dbd9a9ec5d9b290e04d78d4f6b5be9c413ff30ea4"
	}
}`)

	k, err := DecryptKey(keyjson, "password")
	require.NoError(t, err)

	require.Equal(t, common.HexToAddress("0x86bce8c859f5f304aa30adb89f2f7b6ee5a0d6e2"), k.GetAddress())
	key, err := crypto.ToECDSA(common.Hex2Bytes("36e0a792553f94a7660e5484cfc8367e7d56a383261175b9abced7416a5d87df"))
	require.NoError(t, err)

	require.Equal(t, key, k.GetPrivateKeys()[0][0])
}
