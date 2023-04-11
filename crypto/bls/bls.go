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
	"github.com/klaytn/klaytn/crypto/bls/blst"
	"github.com/klaytn/klaytn/crypto/bls/types"
)

type (
	SecretKey = types.SecretKey
	PublicKey = types.PublicKey
	Signature = types.Signature
)

// RandKey generates a random BLS secret key.
func RandKey() (SecretKey, error) {
	return blst.RandKey()
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

// AggregatePublicKeys aggregates multiple BLS public keys.
// Assumes that all given public keys are previously validated.
// Returns error if an empty slice is given.
func AggregatePublicKeys(pks []PublicKey) (PublicKey, error) {
	return blst.AggregatePublicKeys(pks)
}

// AggregatePublicKeysFromBytes unmarshals and validates multiple BLS public key from bytes
// and then aggregates them. Returns error if an empty slice is given.
func AggregatePublicKeysFromBytes(bs [][]byte) (PublicKey, error) {
	return blst.AggregatePublicKeysFromBytes(bs)
}

// AggregatePublicKeys aggregates multiple BLS signatures.
// Assumes that all given signatures are previously validated.
// Returns error if an empty slice is given.
func AggregateSignatures(sigs []Signature) (Signature, error) {
	return blst.AggregateSignatures(sigs)
}

// AggregateSignaturesFromBytes unmarshals and validates multiple BLS signatures from bytes
// and then aggregates them. Returns error if an empty slice is given.
func AggregateSignaturesFromBytes(bs [][]byte) (Signature, error) {
	return blst.AggregateSignaturesFromBytes(bs)
}

// Sign calculates a signature.
func Sign(sk SecretKey, msg []byte) Signature {
	return blst.Sign(sk, msg)
}

// Verify checks a signature. To perform aggregate verify, supply the
// aggregate signature and aggregate public key.
func Verify(sig Signature, msg []byte, pk PublicKey) bool {
	return blst.Verify(sig, msg, pk)
}
