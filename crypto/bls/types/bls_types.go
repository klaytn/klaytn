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

package types

import (
	"errors"
	"fmt"
)

var (
	// For immediate clarity, using literal numbers instead of blst named constants
	SecretKeyLength = 32
	PublicKeyLength = 48
	SignatureLength = 96
	// draft-irtf-cfrg-pairing-friendly-curves-11#4.2.1 BLS12_381
	CurveOrderHex = "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"
	// draft-irtf-cfrg-bls-signature-05#4.2.3
	DomainSeparationTag = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")
)

var (
	ErrSecretKeyGen       = errors.New("BLS secret keygen failed")
	ErrSecretKeyUnmarshal = errors.New("BLS secret key unmarshal failed")
	ErrPublicKeyUnmarshal = errors.New("BLS public key unmarshal failed")
	ErrPublicKeyAggregate = errors.New("BLS public key aggregation failed")
	ErrSignatureUnmarshal = errors.New("BLS signature unmarshal failed")
	ErrSignatureAggregate = errors.New("BLS signature aggregation failed")
	ErrEmptyArray         = errors.New("BLS aggregation failed due to empty array")
)

func ErrSecretKeyLength(have int) error {
	return fmt.Errorf("BLS secret key length mismatch: want: %d have: %d", SecretKeyLength, have)
}

func ErrPublicKeyLength(have int) error {
	return fmt.Errorf("BLS public key length mismatch: want: %d have: %d", PublicKeyLength, have)
}

func ErrSignatureLength(have int) error {
	return fmt.Errorf("BLS signature length mismatch: want: %d have: %d", SignatureLength, have)
}

type SecretKey interface {
	PublicKey() PublicKey
	Marshal() []byte
}

type PublicKey interface {
	Marshal() []byte
	Copy() PublicKey
}

type Signature interface {
	Marshal() []byte
	Copy() Signature
}
