// Copyright 2018 The klaytn Authors
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

/*
Package ecies implements the Elliptic Curve Integrated Encryption Scheme.

Package ecies provides functions to Encrypt / Decrypt messages using ECIES which are used in RLPx transport protocol in Klaytn.

Source Files

Each file contains following contents
 - ecies.go  : Provides encryption / decryption related functions used for exchanging messages in RLPx protocol
 - params.go : Contains parameters for ECIES encryption, specifying the symmetric encryption and HMAC parameters
*/
package ecies
