// Copyright 2018 The klaytn Authors
//
// This file is derived from crypto/bn256 package (2018/06/04).
// See LICENSE in the crypto/bn256 directory for the original copyright and license.

/*
Package bn256 implements a particular bilinear group.

Bilinear groups are the basis of many of the new cryptographic protocols that have been proposed over the past decade. They consist of a triplet of groups (G₁, G₂ and GT) such that there exists a function e(g₁ˣ,g₂ʸ)=gTˣʸ (where gₓ is a generator of the respective group). That function is called a pairing function.

This package specifically implements the Optimal Ate pairing over a 256-bit Barreto-Naehrig curve as described in http://cryptojedi.org/papers/dclxvi-20100714.pdf. Its output is compatible with the implementation described in that paper.

This package previously claimed to operate at a 128-bit security level. However, recent improvements in attacks mean that it is no longer true. See https://moderncrypto.org/mail-archive/curves/2016/000740.html.

Initial package from Google is deprecated due to its weakened security and Klaytn is using the more complete implementation from Cloudflare

Source Files

Each file contains following contents
 - bn256_fast.go : Provides PairingCheck function for amd64 and arm64 architecture
 - bn256_fuzz.go : Provides functions to check if Google's library and Cloudflare's library are giving the same results
 - bn256_slow.go : Provides PairingCheck function for non-amd64 and non-arm64 architecture
*/
package bn256
