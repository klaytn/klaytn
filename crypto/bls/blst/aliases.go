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

package blst

import blst "github.com/supranational/blst/bindings/go"

// Aliases to underlying blst go binding symbols
//
// Klaytn uses the "minimal-signature-size" variant as per
// draft-irtf-cfrg-bls-signature-05#2.1
// where public keys are points in G2, signatures are points in G1.
type blstSecretKey = blst.SecretKey
type blstPublicKey = blst.P2Affine
type blstSignature = blst.P1Affine
type blstAggregatePublicKey = blst.P2Aggregate
type blstAggregateSignature = blst.P1Aggregate
