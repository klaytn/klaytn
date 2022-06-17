// Copyright 2019 The klaytn Authors
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

package accountkey

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/rlp"
)

var (
	errNotS256Curve = errors.New("key is not on the S256 curve")
	errNoXYValue    = errors.New("X or Y value of the public key does not exist")
)

// Since ecdsa.PublicKey does not provide RLP/JSON serialization,
// PublicKeySerializable provides RLP/JSON serialization.
// It is used for AccountKey as an internal structure.
type PublicKeySerializable ecdsa.PublicKey

type publicKeySerializableInternalJSON struct {
	X *hexutil.Big `json:"x"`
	Y *hexutil.Big `json:"y"`
}

// newPublicKeySerializable creates a PublicKeySerializable object.
// The object is initialized with default values.
// Curve = S256 curve
// X = 0
// Y = 0
func newPublicKeySerializable() *PublicKeySerializable {
	return &PublicKeySerializable{
		Curve: crypto.S256(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}
}

// EncodeRLP encodes ecdsa.PublicKey using RLP.
// For now, it supports S256 curve only.
// For that reason, this function serializes only X and Y using CompressPubkey().
func (p *PublicKeySerializable) EncodeRLP(w io.Writer) error {
	// Do not serialize if it is not on S256 curve.
	if !crypto.S256().IsOnCurve(p.X, p.Y) {
		return errNotS256Curve
	}
	return rlp.Encode(w, crypto.CompressPubkey((*ecdsa.PublicKey)(p)))
}

// DecodeRLP decodes PublicKeySerializable using RLP.
// For now, it supports S256 curve only.
// This function deserializes using UncompressPubkey().
func (p *PublicKeySerializable) DecodeRLP(s *rlp.Stream) error {
	b := []byte{}
	if err := s.Decode(&b); err != nil {
		return err
	}
	pubkey, err := crypto.DecompressPubkey(b)
	if err != nil {
		return err
	}
	*p = *((*PublicKeySerializable)(pubkey))

	return nil
}

// MarshalJSON encodes PublicKeySerializable using JSON.
// For now, it supports S256 curve only.
// For that reason, this function serializes only X and Y.
func (p *PublicKeySerializable) MarshalJSON() ([]byte, error) {
	// Do not serialize if it is not on S256 curve.
	if !crypto.S256().IsOnCurve(p.X, p.Y) {
		return nil, errNotS256Curve
	}
	return json.Marshal(&publicKeySerializableInternalJSON{
		(*hexutil.Big)(p.X), (*hexutil.Big)(p.Y)})
}

// UnmarshalJSON decodes PublicKeySerializable using JSON.
// For now, it supports S256 curve only.
// For that reason, this function deserializes only X and Y. Refer to MarshalJSON() above.
func (p *PublicKeySerializable) UnmarshalJSON(b []byte) error {
	var dec publicKeySerializableInternalJSON
	if err := json.Unmarshal(b, &dec); err != nil {
		return err
	}
	if dec.X == nil || dec.Y == nil {
		return errNoXYValue
	}
	p.X = (*big.Int)(dec.X)
	p.Y = (*big.Int)(dec.Y)

	return nil
}

// DeepCopy creates a new PublicKeySerializable object and newly allocates memory for all its attributes.
// Then, the values of the original object are copied to those of the new object.
func (p *PublicKeySerializable) DeepCopy() *PublicKeySerializable {
	pk := newPublicKeySerializable()
	pk.X = new(big.Int).Set(p.X)
	pk.Y = new(big.Int).Set(p.Y)

	return pk
}

// Equal returns true if all attributes between p and pk are the same.
// Otherwise, it returns false.
func (p *PublicKeySerializable) Equal(pk *PublicKeySerializable) bool {
	return p.X.Cmp(pk.X) == 0 &&
		p.Y.Cmp(pk.Y) == 0
}

// String returns a string containing information of all attributes.
func (p *PublicKeySerializable) String() string {
	b, _ := json.Marshal(p)

	return fmt.Sprintf("S256Pubkey:%s", string(b))
}
