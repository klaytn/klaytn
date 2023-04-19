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
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test single signature verification.
// Usage example: validating one IBFT Commit message in handleCommit()
func TestVerify(t *testing.T) {
	sk, _ := RandKey()
	pk := sk.PublicKey()

	var msg [32]byte
	sigb := Sign(sk, msg[:]).Marshal()

	ok, err := VerifySignature(sigb, msg, pk)
	assert.Nil(t, err)
	assert.True(t, ok)
}

// Test aggregated signature verification for the same message.
// Usage example: validating aggregated block signature in VerifyHeader()
func TestAggregateVerify(t *testing.T) {
	var (
		sk1, _ = RandKey()
		sk2, _ = RandKey()

		pkb1 = sk1.PublicKey().Marshal()
		pkb2 = sk2.PublicKey().Marshal()
		pkbs = [][]byte{pkb1, pkb2}

		msg = [32]byte{'a', 'b', 'c'}

		sigb1 = Sign(sk1, msg[:]).Marshal()
		sigb2 = Sign(sk2, msg[:]).Marshal()
		sigbs = [][]byte{sigb1, sigb2}
	)

	apk, _ := AggregatePublicKeys(pkbs)
	asig, _ := AggregateCompressedSignatures(sigbs)
	asigb := asig.Marshal()
	ok, err := VerifySignature(asigb, msg, apk)
	assert.Nil(t, err)
	assert.True(t, ok)
}

// Test aggregated signatrue verification for distinct messages.
// Usage example: validating BLS registry contract
func TestMultipleVerify(t *testing.T) {
	var (
		sk1, _ = RandKey()
		sk2, _ = RandKey()

		pkb1 = sk1.PublicKey().Marshal()
		pkb2 = sk2.PublicKey().Marshal()
		pkbs = [][]byte{pkb1, pkb2}

		msg1 = [32]byte{'1', '2', '3'}
		msg2 = [32]byte{'4', '5', '6'}
		msgs = [][32]byte{msg1, msg2}

		sigb1 = Sign(sk1, msg1[:]).Marshal()
		sigb2 = Sign(sk2, msg2[:]).Marshal()
		sigbs = [][]byte{sigb1, sigb2}
	)

	pks, _ := MultiplePublicKeysFromBytes(pkbs)
	ok, err := VerifyMultipleSignatures(sigbs, msgs, pks)
	assert.Nil(t, err)
	assert.True(t, ok)
}
