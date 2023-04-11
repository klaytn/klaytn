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

import (
	"github.com/klaytn/klaytn/crypto/bls/types"
)

type publicKey struct {
	// Pointer to underlying blst struct, hence the name 'p'
	p *blstPublicKey
}

func PublicKeyFromBytes(b []byte) (types.PublicKey, error) {
	if len(b) != types.PublicKeyLength {
		return nil, types.ErrPublicKeyLength(len(b))
	}

	if pk, ok := publicKeyCache().Get(cacheKey(b)); ok {
		return pk.(*publicKey), nil
	}

	p := new(blstPublicKey).Uncompress(b)
	if p == nil || !p.KeyValidate() {
		return nil, types.ErrPublicKeyUnmarshal
	}

	pk := &publicKey{p: p}
	publicKeyCache().Add(cacheKey(b), pk)
	return pk, nil
}

func MultiplePublicKeysFromBytes(bs [][]byte) ([]types.PublicKey, error) {
	if len(bs) == 0 {
		return nil, nil
	}
	for _, b := range bs {
		if len(b) != types.PublicKeyLength {
			return nil, types.ErrPublicKeyLength(len(b))
		}
	}

	// Separate cached and uncached element
	pks := make([]types.PublicKey, len(bs))
	var batchIndices []int
	var batchBytes [][]byte
	for i, b := range bs {
		if pk, ok := publicKeyCache().Get(cacheKey(b)); ok {
			pks[i] = pk.(*publicKey)
		} else {
			batchIndices = append(batchIndices, i)
			batchBytes = append(batchBytes, b)
		}
	}

	// Compute on uncached elements
	batchPs := new(blstPublicKey).BatchUncompress(batchBytes)
	if batchPs == nil || len(batchPs) != len(batchBytes) {
		return nil, types.ErrPublicKeyUnmarshal
	}

	// Merge cached and uncached elements
	for i, outIdx := range batchIndices {
		b := batchBytes[i]
		p := batchPs[i]

		if p == nil || !p.KeyValidate() {
			return nil, types.ErrPublicKeyUnmarshal
		}

		pk := &publicKey{p: p}
		publicKeyCache().Add(cacheKey(b), pk)
		pks[outIdx] = pk
	}

	return pks, nil
}

func AggregatePublicKeys(pks []types.PublicKey) (types.PublicKey, error) {
	if len(pks) == 0 {
		return nil, types.ErrEmptyArray
	}

	ps := make([]*blstPublicKey, len(pks))
	for i, pk := range pks {
		ps[i] = pk.(*publicKey).p
	}

	agg := new(blstAggregatePublicKey)
	groupcheck := false // alreaay checked in *PublicKeyFromBytes()
	if !agg.Aggregate(ps, groupcheck) {
		return nil, types.ErrPublicKeyAggregate
	}
	return &publicKey{p: agg.ToAffine()}, nil
}

func AggregatePublicKeysFromBytes(bs [][]byte) (types.PublicKey, error) {
	pks, err := MultiplePublicKeysFromBytes(bs)
	if err != nil {
		return nil, err
	}
	return AggregatePublicKeys(pks)
}

func (pk *publicKey) Marshal() []byte {
	return pk.p.Compress()
}
