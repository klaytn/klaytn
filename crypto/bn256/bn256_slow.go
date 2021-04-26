// Copyright 2018 The klaytn Authors
// Copyright 2018 Péter Szilágyi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
//
// This file is derived from crypto/bn256/bn256_slow.go (2018/06/04).
// Modified and improved for the klaytn development.

// +build !amd64,!arm64

package bn256

import bn256 "github.com/klaytn/klaytn/crypto/bn256/cloudflare"

// G1 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G1 = bn256.G1

// G2 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G2 = bn256.G2

// PairingCheck calculates the Optimal Ate pairing for a set of points.
func PairingCheck(a []*G1, b []*G2) bool {
	return bn256.PairingCheck(a, b)
}
