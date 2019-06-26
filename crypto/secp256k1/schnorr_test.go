// Copyright 2018 The klaytn Authors
//
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

package secp256k1

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestSingle(t *testing.T) {
	m := hash([]byte("hello world"))

	G := S256()
	x, _ := hex.DecodeString("4f52e337ad8bf1ce10cbb72ab91d9954474cea39811040df5558297df3e3c1bf")
	P := G.Marshal(G.ScalarBaseMult(x))

	R, s := SchnorrSignSingle(G, m, x, P)

	start := time.Now()
	result := SchnorrVerify(G, m, R, s, P)
	elapsed := time.Since(start)

	fmt.Printf("Go Imp: %d us\n", elapsed/1000)
	if !result {
		t.Fail()
	}

	start = time.Now()
	result = SchnorrVerifyNative(m, R, s, P)
	elapsed = time.Since(start)

	fmt.Printf("C Imp: %d us\n", elapsed/1000)
	if !result {
		t.Fail()
	}
}

func TestDouble(t *testing.T) {
	msg := hash([]byte("hello"))
	G := S256()

	// For two arbitrary SECP256k1 private keys
	x0, _ := hex.DecodeString("4f52e337ad8bf1ce10cbb72ab91d9954474cea39811040df5558297df3e3c1bf") // Alice
	x1, _ := hex.DecodeString("54390f30fbe7094a5080cd4dcc3f91869faaca23a43aed2eb20a426db067b7c9") // Bob

	// Derive public keys from the private keys; P0 = Alice's pubkey, P1 = Bob's pubkey
	P0x, P0y := G.ScalarBaseMult(x0)
	P0 := G.Marshal(P0x, P0y)

	P1x, P1y := G.ScalarBaseMult(x1)
	P1 := G.Marshal(P1x, P1y)

	// Begin the signing party
	Q0, R0, y0 := SchnorrSignMultiBootstrap(G, msg, x0, P0, P1) // Alice
	Q1, R1, y1 := SchnorrSignMultiBootstrap(G, msg, x1, P1, P0) // Bob

	// Unpack the masked public keys
	Q0x, Q0y := G.Unmarshal(Q0)
	Q1x, Q1y := G.Unmarshal(Q1)

	// Combine them to get the common public key P
	P := G.Marshal(G.Add(Q0x, Q0y, Q1x, Q1y))

	// Unpack R values
	R0x, R0y := G.Unmarshal(R0)
	R1x, R1y := G.Unmarshal(R1)

	// Combine them to get the R; this is the first chunk of the signature signed by both Alice and Bob
	R := G.Marshal(G.Add(R0x, R0y, R1x, R1y))

	// Now both Alice and Bob know P and R; let's get the rest --- s
	s0 := SchnorrSignMultiComputeS(msg, P, R, y0)
	s1 := SchnorrSignMultiComputeS(msg, P, R, y1)

	// Let's combine this together
	s := ScAdd(s0, s1)

	// At this point, we have (R, s), a complete Schnorr signature signed by both Alice and Bob.
	// We can verify (R, s) by checking R = s * G + H(m || P || R) * P.

	// Let's verify the signature
	i := int64(0)
	var start time.Time
	var elapsed, sumNative, sumGo time.Duration

	for ; i < 2000; i++ {
		// Run verification using pure Go impl
		start = time.Now()

		if !SchnorrVerify(G, msg, R, s, P) {
			t.Fail()
		}

		elapsed = time.Since(start)
		sumGo += elapsed

		// Run verification using native C impl
		start = time.Now()
		if !SchnorrVerifyNative(msg, R, s, P) {
			t.Fail()
		}

		elapsed = time.Since(start)
		sumNative += elapsed
	}

	fmt.Printf("schnorr_go = %d us\n", sumGo.Nanoseconds()/i/1000)
	fmt.Printf("schnorr_native = %d us\n", sumNative.Nanoseconds()/i/1000)
}

func TestScPointMulTest(t *testing.T) {
	scalar := hash([]byte("hello"))
	G := S256()

	P0x, P0y := G.ScalarBaseMult(scalar)
	P := G.Marshal(P0x, P0y)

	pureGo := G.Marshal(G.ScalarMult(P0x, P0y, scalar))
	cImp := ScPointMul(P, scalar)

	if !bytes.Equal(pureGo, cImp) {
		t.Fail()
	}
}

func TestScBaseMultTest(t *testing.T) {
	scalar := hash([]byte("hello"))
	G := S256()

	pureGo := G.Marshal(G.ScalarBaseMult(scalar))
	cImp := ScBaseMul(scalar)

	if !bytes.Equal(pureGo, cImp) {
		t.Fail()
	}
}

func TestComputeC(t *testing.T) {
	G := S256()

	x0, _ := hex.DecodeString("4f52e337ad8bf1ce10cbb72ab91d9954474cea39811040df5558297df3e3c1bf")
	x1, _ := hex.DecodeString("54390f30fbe7094a5080cd4dcc3f91869faaca23a43aed2eb20a426db067b7c9")

	P0x, P0y := G.ScalarBaseMult(x0)
	P0 := G.Marshal(P0x, P0y)

	P1x, P1y := G.ScalarBaseMult(x1)
	P1 := G.Marshal(P1x, P1y)

	C := ComputeC(P0, P1)
	K := ComputeC(P1, P0)

	if !bytes.Equal(C, K) {
		t.Fail()
	}
}
