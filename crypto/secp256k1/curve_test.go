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
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAddSamePoint(t *testing.T) {
	/*
		This test is intended to highlight the bug in klaytn/crypto/secp256k1/curve.go#affineFromJacobian.
		When passed with same points, BitCurve.Add invokes affineFromJacobian(0, 0, 0) which then invokes
		(big.Int).Mul(nil, nil).

		Although executing (big.Int).Mul(nil, nil) is not problematic in Go 1.10.3, it has been found invoking the same
		is fatal in Go 1.11 (causing SIGSEGV; terminating the program).

	*/
	x0, _ := hex.DecodeString("4f52e337ad8bf1ce10cbb72ab91d9954474cea39811040df5558297df3e3c1bf") // Alice
	x1, _ := hex.DecodeString("4f52e337ad8bf1ce10cbb72ab91d9954474cea39811040df5558297df3e3c1bf") // Alice

	G := S256()
	P0x, P0y := G.ScalarBaseMult(x0)
	P1x, P1y := G.ScalarBaseMult(x1)

	Qx, Qy := G.Add(P0x, P0y, P1x, P1y)

	fmt.Println(Qx, Qy) // should print 0 0 with no fatal error
}
