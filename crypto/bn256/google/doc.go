// Copyright 2018 The klaytn Authors
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file is derived from crypto/bn256/google/bn256.go (2018/06/04).
// Modified and improved for the klaytn development.

// Package bn256 implements a particular bilinear group at the 128-bit security level.
//
// Bilinear groups are the basis of many of the new cryptographic protocols
// that have been proposed over the past decade. They consist of a triplet of
// groups (G₁, G₂ and GT) such that there exists a function e(g₁ˣ,g₂ʸ)=gTˣʸ
// (where gₓ is a generator of the respective group). That function is called
// a pairing function.
//
// This package specifically implements the Optimal Ate pairing over a 256-bit
// Barreto-Naehrig curve as described in
// http://cryptojedi.org/papers/dclxvi-20100714.pdf. Its output is compatible
// with the implementation described in that paper.
//
//This package previously claimed to operate at a 128-bit security level. However, recent improvements in attacks mean that it is no longer true. See https://moderncrypto.org/mail-archive/curves/2016/000740.html.
//
//Deprecated: due to its weakened security, new systems should not rely on this elliptic curve. This package is frozen, and not implemented in constant time. There is a more complete implementation at github.com/cloudflare/bn256, but note that it suffers from the same security issues of the underlying curve.
package bn256
