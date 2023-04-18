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
	sk1, _ := RandKey()
	sk2, _ := RandKey()

	pkb1 := sk1.PublicKey().Marshal()
	pkb2 := sk2.PublicKey().Marshal()

	var msg [32]byte
	sigb1 := Sign(sk1, msg[:]).Marshal()
	sigb2 := Sign(sk2, msg[:]).Marshal()

	apk, _ := AggregatePublicKeys([][]byte{pkb1, pkb2})
	asig, _ := AggregateCompressedSignatures([][]byte{sigb1, sigb2})

	asigb := asig.Marshal()
	ok, err := VerifySignature(asigb, msg, apk)
	assert.Nil(t, err)
	assert.True(t, ok)
}
