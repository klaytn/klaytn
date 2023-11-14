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

package bls

import (
	"crypto/ecdsa"
	"encoding/hex"
	"os"
	"strings"

	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls/blst"
	"github.com/klaytn/klaytn/crypto/bls/types"
)

type (
	SecretKey = types.SecretKey
	PublicKey = types.PublicKey
	Signature = types.Signature
)

// Applications are expected to use below top-level functions,
// Do not directly use 'blst' or other underlying implementations.
//
// Some functions are named after the equivalent functions in prysm Ethereum CL.
// (https://github.com/prysmaticlabs/prysm/blob/v4.0.2/crypto/bls/bls.go)
// Such naming should provide compatiblity with prysm code snippets,
// in case prysm code snippets are integrated to klaytn.
//
// ikm -> SK:  GenerateKey
// ()  -> SK:  RandKey
// ec  -> SK:  DeriveFromECDSA
// path -> SK: LoadKey
// b32 -> SK:  SecretKeyFromBytes
// b48 -> PK:  PublicKeyFromBytes
// b96 -> Sig: SignatureFromBytes
//
// []b48 -> []PK:  MultiplePublicKeysFromBytes
// []b96 -> []Sig: MultipleSignaturesFromBytes
//
// []PK  -> PK:    AggregateMultiplePubkeys
// []Sig -> Sig:   AggregateSignatures
//
// []b48 -> PK:    AggregatePublicKeys
// []b96 -> Sig:   AggregateCompressedSignatures
//
// Sign(SK, msg) -> Sig
// VerifySignature(b96, msg, PK) -> ok, err
// PopProve(SK) -> Proof
// PoPVerify(PK, Proof) -> ok, err

// GenerateKey generates a BLS secret key from the initial key material (IKM).
// It is deterministic process. Same IKM yields the same secret key.
func GenerateKey(ikm []byte) (SecretKey, error) {
	return blst.GenerateKey(ikm)
}

// RandKey generates a random BLS secret key.
func RandKey() (SecretKey, error) {
	return blst.RandKey()
}

// DeriveFromECDSA generates a BLS secret key from the given EC private key.
// It is deterministic process. Same EC private key yields the same secret key.
func DeriveFromECDSA(priv *ecdsa.PrivateKey) (SecretKey, error) {
	return GenerateKey(crypto.FromECDSA(priv))
}

// LoadKey loads a BLS secret key from the given file.
func LoadKey(path string) (SecretKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(b)
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "0x")
	b, err = hex.DecodeString(content)
	if err != nil {
		return nil, err
	}

	return SecretKeyFromBytes(b)
}

// SaveKey stores a BLS secret key to the given file.
func SaveKey(path string, sk SecretKey) error {
	b := hex.EncodeToString(sk.Marshal())
	return os.WriteFile(path, []byte(b), 0o600)
}

// SecretKeyFromBytes unmarshals and validates a BLS secret key from bytes.
func SecretKeyFromBytes(b []byte) (SecretKey, error) {
	return blst.SecretKeyFromBytes(b)
}

// PublicKeyFromBytes unmarshals and validates a BLS public key from bytes.
func PublicKeyFromBytes(b []byte) (PublicKey, error) {
	return blst.PublicKeyFromBytes(b)
}

// SignatureFromBytes unmarshals and validates a BLS signature from bytes.
func SignatureFromBytes(b []byte) (Signature, error) {
	return blst.SignatureFromBytes(b)
}

// MultiplePublicKeysFromBytes unmarshals and validates multiple BLS public keys
// from bytes. Returns an empty slice if an empty slice is given.
func MultiplePublicKeysFromBytes(bs [][]byte) ([]PublicKey, error) {
	return blst.MultiplePublicKeysFromBytes(bs)
}

// MultipleSignaturesFromBytes unmarshals multiple BLS signatures from bytes.
// Returns an empty slice if an empty slice is given.
func MultipleSignaturesFromBytes(bs [][]byte) ([]Signature, error) {
	return blst.MultipleSignaturesFromBytes(bs)
}

// AggregateMultiplePubkeys aggregates multiple BLS public keys.
// Assumes that all given public keys are previously validated.
// Returns error if an empty slice is given.
func AggregateMultiplePubkeys(pks []PublicKey) (PublicKey, error) {
	return blst.AggregatePublicKeys(pks)
}

// AggregateSignatures aggregates multiple BLS signatures.
// Assumes that all given signatures are previously validated.
// Returns error if an empty slice is given.
func AggregateSignatures(sigs []Signature) (Signature, error) {
	return blst.AggregateSignatures(sigs)
}

// AggregatePublicKeys unmarshals and validates multiple BLS public key from bytes
// and then aggregates them. Returns error if an empty slice is given.
func AggregatePublicKeys(bs [][]byte) (PublicKey, error) {
	return blst.AggregatePublicKeysFromBytes(bs)
}

// AggregateCompressedSignatures unmarshals and validates multiple BLS signatures from bytes
// and then aggregates them. Returns error if an empty slice is given.
func AggregateCompressedSignatures(bs [][]byte) (Signature, error) {
	return blst.AggregateSignaturesFromBytes(bs)
}

// Sign calculates a signature.
func Sign(sk SecretKey, msg []byte) Signature {
	return blst.Sign(sk, msg)
}

// VerifySignature checks a signature. To perform aggregate verify, supply the
// aggregate signature and aggregate public key.
func VerifySignature(sig []byte, msg [32]byte, pk PublicKey) (bool, error) {
	return blst.VerifySignature(sig, msg, pk)
}

// VerifyMultipleSignatures verifies multiple signatures for distinct messages securely.
func VerifyMultipleSignatures(sigs [][]byte, msgs [][32]byte, pubKeys []PublicKey) (bool, error) {
	return blst.VerifyMultipleSignatures(sigs, msgs, pubKeys)
}

// PopProve calculates the proof-of-possession for the secret key,
// which is the signature with its public key as message.
func PopProve(sk SecretKey) Signature {
	// draft-irtf-cfrg-bls-signature-05 section 3.3.2. PopProve
	msg := sk.PublicKey().Marshal()
	return blst.Sign(sk, msg)
}

// PopVerify verifies the proof-of-possession for the public key.
func PopVerify(pk PublicKey, proof Signature) bool {
	// draft-irtf-cfrg-bls-signature-05 section 3.3.3. PopVerify
	msg := pk.Marshal()
	return blst.Verify(proof, msg, pk)
}
