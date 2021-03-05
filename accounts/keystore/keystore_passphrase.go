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
// This file is derived from accounts/keystore/keystore_passphrase.go (2018/06/04).
// Modified and improved for the klaytn development.

package keystore

import (
	"bytes"
	"crypto/aes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/crypto"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

const (
	keyHeaderKDF = "scrypt"

	// StandardScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptN = 1 << 18

	// StandardScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptP = 1

	// LightScryptN is the N parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptN = 1 << 12

	// LightScryptP is the P parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	LightScryptP = 6

	scryptR     = 8
	scryptDKLen = 32
)

type keyStorePassphrase struct {
	keysDirPath string
	scryptN     int
	scryptP     int
	// skipKeyFileVerification disables the security-feature which does
	// reads and decrypts any newly created keyfiles. This should be 'false' in all
	// cases except tests -- setting this to 'true' is not recommended.
	skipKeyFileVerification bool
}

func (ks keyStorePassphrase) GetKey(addr common.Address, filename, auth string) (Key, error) {
	// Load the key from the keystore and decrypt its contents
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	key, err := DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.GetAddress() != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.GetAddress(), addr)
	}
	return key, nil
}

// StoreKey generates a key, encrypts with 'auth' and stores in the given directory
func StoreKey(dir, auth string, scryptN, scryptP int) (common.Address, error) {
	_, a, err := storeNewKey(&keyStorePassphrase{dir, scryptN, scryptP, false}, rand.Reader, auth)
	return a.Address, err
}

func (ks keyStorePassphrase) StoreKey(filename string, key Key, auth string) error {
	keyjson, err := EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	if err != nil {
		return err
	}
	// Write into temporary file
	tmpName, err := writeTemporaryKeyFile(filename, keyjson)
	if err != nil {
		return err
	}
	if ks.skipKeyFileVerification == false { //do not skip file verification
		// Verify that we can decrypt the file with the given password.
		_, err = ks.GetKey(key.GetAddress(), tmpName, auth)
		if err != nil {
			msg := "An error was encountered when saving and verifying the keystore file. \n" +
				"This indicates that the keystore is corrupted. \n" +
				"The corrupted file is stored at \n%v\n" +
				"Please file a ticket at:\n\n" +
				"https://github.com/klaytn/klaytn/issues" +
				"The error was : %w"
			return fmt.Errorf(msg, tmpName, err)
		}
	}
	return os.Rename(tmpName, filename)
}

func (ks keyStorePassphrase) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ks.keysDirPath, filename)
}

// encryptCrypto encrypts a private key to a cryptoJSON object.
func encryptCrypto(keyBytes []byte, auth string, scryptN, scryptP int) (*cryptoJSON, error) {
	authArray := []byte(auth)

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	derivedKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	encryptKey := derivedKey[:16]

	iv := make([]byte, aes.BlockSize) // 16
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return nil, err
	}
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	return &cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}, nil
}

// EncryptKey encrypts a key using the specified scrypt parameters into a json
// blob that can be decrypted later on. It uses the keystore v4 format.
func EncryptKey(key Key, auth string, scryptN, scryptP int) ([]byte, error) {
	pks := key.GetPrivateKeys()
	crypto := make([][]cryptoJSON, len(pks))
	for i, keys := range pks {
		crypto[i] = make([]cryptoJSON, len(keys))
		for j, k := range keys {
			keyBytes := math.PaddedBigBytes(k.D, 32)
			c, err := encryptCrypto(keyBytes, auth, scryptN, scryptP)
			if err != nil {
				return nil, err
			}
			crypto[i][j] = *c
		}
	}
	encryptedKeyJSONV4 := encryptedKeyJSONV4{
		hex.EncodeToString(key.GetAddress().Bytes()),
		crypto,
		key.GetId().String(),
		4,
	}
	return json.Marshal(encryptedKeyJSONV4)
}

// EncryptKeyV3 encrypts a key using the specified scrypt parameters into a json
// blob that can be decrypted later on. It uses the keystore v3 format.
func EncryptKeyV3(key Key, auth string, scryptN, scryptP int) ([]byte, error) {
	authArray := []byte(auth)

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	derivedKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	encryptKey := derivedKey[:16]
	keyBytes := math.PaddedBigBytes(key.GetPrivateKey().D, 32)

	iv := make([]byte, aes.BlockSize) // 16
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return nil, err
	}
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	cryptoStruct := cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	encryptedKeyJSONV3 := encryptedKeyJSONV3{
		hex.EncodeToString(key.GetAddress().Bytes()),
		cryptoStruct,
		key.GetId().String(),
		3,
	}
	return json.Marshal(encryptedKeyJSONV3)
}

// DecryptKey decrypts a key from a json blob, returning the private key itself.
// TODO: use encryptedKeyJSON object directly instead of double unmarshalling.
func DecryptKey(keyjson []byte, auth string) (Key, error) {
	// Parse the json into a simple map to fetch the key version
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &m); err != nil {
		return nil, err
	}
	// Depending on the version try to parse one way or another
	var (
		keyBytes [][][]byte
		keyId    []byte
		err      error
		address  common.Address
	)

	switch v := m["version"].(type) {
	case string:
		if v == "1" {
			k := new(encryptedKeyJSONV1)
			if err := json.Unmarshal(keyjson, k); err != nil {
				return nil, err
			}
			keyBytes = make([][][]byte, 1)
			keyBytes[0] = make([][]byte, 1)
			keyBytes[0][0], keyId, err = decryptKeyV1(k, auth)
			address = common.HexToAddress(k.Address)
		}

	case float64:
		switch v {
		case 3:
			k := new(encryptedKeyJSONV3)
			if err := json.Unmarshal(keyjson, k); err != nil {
				return nil, err
			}
			keyBytes = make([][][]byte, 1)
			keyBytes[0] = make([][]byte, 1)
			keyBytes[0][0], keyId, err = decryptKeyV3(k, auth)
			address = common.HexToAddress(k.Address)
		case 4:
			// At first, try to decrypt using encryptedKeyJSONV4
			k := new(encryptedKeyJSONV4)
			if err := json.Unmarshal(keyjson, k); err != nil {
				// If it fails, try to decrypt using encryptedKeyJSONV4Single
				kSingle := new(encryptedKeyJSONV4Single)
				if err = json.Unmarshal(keyjson, kSingle); err != nil {
					return nil, err
				}

				// If succeeded, copy the values of kSingle to k
				k.Id = kSingle.Id
				k.Address = kSingle.Address
				k.Keyring = [][]cryptoJSON{kSingle.Keyring}
				k.Version = kSingle.Version
			}

			keyBytes, keyId, err = decryptKeyV4(k, auth)
			address = common.HexToAddress(k.Address)
		default:
			return nil, fmt.Errorf("undefined version: %f", v)
		}

	default:
		return nil, fmt.Errorf("undefined type of version: %s", m)
	}

	// Handle any decryption errors and return the key
	if err != nil {
		return nil, err
	}
	privateKeys := make([][]*ecdsa.PrivateKey, len(keyBytes))
	for i, keys := range keyBytes {
		privateKeys[i] = make([]*ecdsa.PrivateKey, len(keys))
		for j, key := range keys {
			privateKeys[i][j], err = crypto.ToECDSA(key)
			if err != nil {
				return nil, err
			}
		}
	}

	return &KeyV4{
		Id:          uuid.UUID(keyId),
		Address:     address,
		PrivateKeys: privateKeys,
	}, nil
}

func decryptKeyV4(keyProtected *encryptedKeyJSONV4, auth string) (keyBytes [][][]byte, keyId []byte, err error) {
	if keyProtected.Version != 4 {
		return nil, nil, fmt.Errorf("version not supported: %v (should be 4)", keyProtected.Version)
	}

	keyId = uuid.Parse(keyProtected.Id)
	keyBytes = make([][][]byte, len(keyProtected.Keyring))

	for i, keys := range keyProtected.Keyring {
		keyBytes[i] = make([][]byte, len(keys))
		for j, key := range keys {
			keyBytes[i][j], err = decryptKey(key, auth)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return keyBytes, keyId, err
}

func decryptKeyV3(keyProtected *encryptedKeyJSONV3, auth string) (keyBytes []byte, keyId []byte, err error) {
	if keyProtected.Version != 3 {
		return nil, nil, fmt.Errorf("Version not supported: %v", keyProtected.Version)
	}
	keyId = uuid.Parse(keyProtected.Id)

	plainText, err := decryptKey(keyProtected.Crypto, auth)
	if err != nil {
		return nil, nil, err
	}
	return plainText, keyId, err
}

func decryptKey(cryptoJson cryptoJSON, auth string) (keyBytes []byte, err error) {
	if cryptoJson.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("Cipher not supported: %v", cryptoJson.Cipher)
	}

	mac, err := hex.DecodeString(cryptoJson.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(cryptoJson.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(cryptoJson.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(cryptoJson, auth)
	if err != nil {
		return nil, err
	}

	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, ErrDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}

	return plainText, err
}

func decryptKeyV1(keyProtected *encryptedKeyJSONV1, auth string) (keyBytes []byte, keyId []byte, err error) {
	keyId = uuid.Parse(keyProtected.Id)
	mac, err := hex.DecodeString(keyProtected.Crypto.MAC)
	if err != nil {
		return nil, nil, err
	}

	iv, err := hex.DecodeString(keyProtected.Crypto.CipherParams.IV)
	if err != nil {
		return nil, nil, err
	}

	cipherText, err := hex.DecodeString(keyProtected.Crypto.CipherText)
	if err != nil {
		return nil, nil, err
	}

	derivedKey, err := getKDFKey(keyProtected.Crypto, auth)
	if err != nil {
		return nil, nil, err
	}

	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, nil, ErrDecrypt
	}

	plainText, err := aesCBCDecrypt(crypto.Keccak256(derivedKey[:16])[:16], cipherText, iv)
	if err != nil {
		return nil, nil, err
	}
	return plainText, keyId, err
}

func getKDFKey(cryptoJSON cryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == keyHeaderKDF {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(authArray, salt, n, r, p, dkLen)

	} else if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("Unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("Unsupported KDF: %s", cryptoJSON.KDF)
}

// TODO: can we do without this when unmarshalling dynamic JSON?
// why do integers in KDF params end up as float64 and not int after
// unmarshal?
func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}
