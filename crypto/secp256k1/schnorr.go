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
	"crypto/sha256"
	"sort"
)

// SchnorrSignSingle digitally signs the input message using Schnorr signature scheme.
func SchnorrSignSingle(G *BitCurve, msg, x, P []byte) ([]byte, []byte) {
	// [Schnorr Signature 101]
	//
	// 1) Deriving public key from private key
	// : P = x * G
	// : P is a corresponding public key
	// : x is a private key
	// : G is a curve base point
	//
	// 2) Creating a random scalar for signing
	// : k = H(m || x)
	//
	// 3) Compute the first part of Schnorr signature, R
	// : R = k * G
	//
	// The goal is to make R = s * G + e * P for some s and e.
	//
	// 4) Compute e
	// : e = H(m || P || R)
	//
	// 5) Find s
	// Recall P = x * G; therefore
	// R     = s * G + e * P
	//       = s * G + e * x * G
	//       = (s + e * x) * G
	// k * G = (s + e * x) * G
	// k     = s + e * x
	// s     = k - e * x
	//
	// Finally, (R, s) is your Schnorr signature

	k := hash(msg, x)                   // k = H(m || x)
	R := G.Marshal(G.ScalarBaseMult(k)) // R = k * G

	e := hash(msg, P, R) // e = H(m || P || R)

	ex := ScMul(e, x) // e * x
	s := ScSub(k, ex) // s = k - e * x

	return R, s
}

// SchnorrVerifySingle verifies a Schnorr signature.
// Note that this implementation is relatively slow compared to the C implementation. Use this sparingly.
func SchnorrVerify(G *BitCurve, msg, R, s, P []byte) bool {
	e := hash(msg, P, R)

	sx, sy := G.ScalarBaseMult(s)

	Px, Py := G.Unmarshal(P)
	ePx, ePy := G.ScalarMult(Px, Py, e)

	V := G.Marshal(G.Add(sx, sy, ePx, ePy))

	return bytes.Equal(R, V) // R == s*G + e*P
}

/*
	[Brief Overview on Schnorr 2-out-of-2 signature scheme]

	In Schnorr signature, Alice shares (R, s) with R = k * G
	where k is an arbitrary number and G is an elliptic curve group.

	R can be computed from (e, s) via R = s * G + e * P for some e ans s and
	e can be computed from (R, s) via e = H(m || P || R)

	To sign a message, Alice needs to pick a value for k. For security reason,
	k should not be reused over multiple messages.

	Given a message m and private key x, we use k = H(m || x) where H is a secure hash function.

	R	=	k * G
		=	H(m || x) * G

	Recall that R = s * G + e * P. Hence,

	R	=	s * G + e * P      = k * G
			s * G + e * x * G  = k * G
			(s + e * x) * G    = k * G
			s + e * x          = k             (*)

	By (*), s = k - e * x                      (**)

	Suppose Alice and Bob have the following keys:

	Alice:	P0 = x0 * G
	Bob:	P1 = x1 * G

	Alice and Bob can safely create a multisig on a message, m, by using generated keys.

	1. 	C = H(P0 || P1)         Let C be the common value that both Alice and Bob use to generate safe keys

	2. 	Q0 = H(C || P0) * P0    Let Q0 be Alice's public key
	3. 	Q1 = H(C || P1) * P1    Let Q1 be Bob's public key
	4. 	P = Q0 + Q1             Let P be the public key for the target multisig

	5. 	y0 = x0 * H(C || P0)    Let y0 be Alice's private key
	6. 	y1 = x1 * H(C || P1)    Let y1 be Bob's private key

	7. 	k0 = H(m || y0)         Alice's random value for m
	8. 	k1 = H(m || y1)         Bob's random value for m
	9. 	R0 = k0 * G
	10.	R1 = k1 * G
	11.	R = R0 + R1

	We use e = H(m || P || R) to harden the security, preventing either Alice or Bob from stealing
	the control over the other party's key.

	Assuming Alice and Bob have already exchanged Q0, Q1, R0 and R1, Alice and Bob can compute the followings:

	e = H(m || P || R)										by (4),  (11)

	Using the computed e, Alice and Bob can compute s0 and s1.

	Alice:   s0 = k0 - e * y0
	Bob:     s1 = k1 - e * y1

	Once they exchange s0 and s1, they can compute s.

	s = s0 + s1

	Thus, Alice and Bob can publish (R, s) as a signature for m.
	Verifying (R, s) can be done by computing R = s * G + H(m || P || R) * P.
*/

// ComputeC derives the common value for multiple public keys.
func ComputeC(keys ...[]byte) []byte {
	sort.SliceStable(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	return hash(keys...)
}

// SchnorrSignMultiBootstrap computes an individual share of a Schnorr multi-signature given all public keys.
// Q: a dedicated, security hardened public key for this multi-signature party
// R: a part of the generating multi-signature for the input publickey
// y: a dedicated, security hardened private key for this multi-signature party
func SchnorrSignMultiBootstrap(G *BitCurve, msg, privateKey, publicKey []byte, othersPublicKeys ...[]byte) (Q, R, y []byte) {
	C := ComputeC(append([][]byte{publicKey}, othersPublicKeys...)...)

	z := hash(C, publicKey)

	Px, Py := G.Unmarshal(publicKey)
	Qx, Qy := G.ScalarMult(Px, Py, z)

	y = ScMul(privateKey, z)

	k := hash(msg, y)

	Rx, Ry := G.ScalarBaseMult(k)

	Q = G.Marshal(Qx, Qy)
	R = G.Marshal(Rx, Ry)
	return
}

// SchnorrSignMultiComputeS computes the s part of a Schnorr multi-signature (i.e., s in (R, s)).
func SchnorrSignMultiComputeS(msg, P, R, y []byte) []byte {
	e := hash(msg, P, R)
	k := hash(msg, y)
	return ScSub(k, ScMul(e, y))
}

// Simple helper function concatenating all input bytes and summing it up with SHA256.
func hash(bytes ...[]byte) []byte {
	h := sha256.New()
	for _, b := range bytes {
		h.Write(b)
	}
	result := h.Sum(nil)
	return result[:]
}
