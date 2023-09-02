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
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/klaytn/klaytn/crypto/bls/types"
	"github.com/pkg/errors"
)

type signature struct {
	// Pointer to underlying blst struct, hence the name 'p'
	p *blstSignature
}

func SignatureFromBytes(b []byte) (types.Signature, error) {
	if len(b) != types.SignatureLength {
		return nil, types.ErrSignatureLength(len(b))
	}

	if s, ok := signatureCache.Get(cacheKey(b)); ok {
		return s.(*signature).Copy(), nil
	}

	p := new(blstSignature).Uncompress(b)
	// Do not check for infinity since an aggregated signature could be infinite.
	if p == nil || !p.SigValidate(false) {
		return nil, types.ErrSignatureUnmarshal
	}

	s := &signature{p: p}
	signatureCache.Add(cacheKey(b), s.Copy())
	return s, nil
}

func MultipleSignaturesFromBytes(bs [][]byte) ([]types.Signature, error) {
	if len(bs) == 0 {
		return nil, nil
	}
	for _, b := range bs {
		if len(b) != types.SignatureLength {
			return nil, types.ErrSignatureLength(len(b))
		}
	}

	// Separate cached and uncached element
	sigs := make([]types.Signature, len(bs))
	var batchIndices []int
	var batchBytes [][]byte
	for i, b := range bs {
		if sig, ok := signatureCache.Get(cacheKey(b)); ok {
			sigs[i] = sig.(*signature).Copy()
		} else {
			batchIndices = append(batchIndices, i)
			batchBytes = append(batchBytes, b)
		}
	}

	// Compute on uncached elements
	batchPs := new(blstSignature).BatchUncompress(batchBytes)
	if batchPs == nil || len(batchPs) != len(batchBytes) {
		return nil, types.ErrSignatureUnmarshal
	}

	// Merge cached and uncached elements
	for i, outIdx := range batchIndices {
		p := batchPs[i]
		b := batchBytes[i]

		// Do not check for infinity since an aggregated signature could be infinite.
		if p == nil || !p.SigValidate(false) {
			return nil, types.ErrSignatureUnmarshal
		}

		sig := &signature{p: p}
		signatureCache.Add(cacheKey(b), sig.Copy())
		sigs[outIdx] = sig
	}

	return sigs, nil
}

// AggregateSignatures assumes that given signatures has passed the SigValidate() check
// i.e. they are on the right subgroup.
//
// Signature objects are expected to be returned by Sign(), SignatureFromBytes()
// and AggregateSignaturesFromBytes(), and they all should be valid.
// Therefore AggregateSignatures skips the validatity check.
func AggregateSignatures(sigs []types.Signature) (types.Signature, error) {
	if len(sigs) == 0 {
		return nil, types.ErrEmptyArray
	}

	ps := make([]*blstSignature, len(sigs))
	for i, sig := range sigs {
		ps[i] = sig.(*signature).p
	}

	agg := new(blstAggregateSignature)
	groupcheck := false // alreaay checked in *SignatureFromBytes()
	if !agg.Aggregate(ps, groupcheck) {
		return nil, types.ErrSignatureAggregate
	}
	return &signature{p: agg.ToAffine()}, nil
}

func AggregateSignaturesFromBytes(bs [][]byte) (types.Signature, error) {
	pks, err := MultipleSignaturesFromBytes(bs)
	if err != nil {
		return nil, err
	}
	return AggregateSignatures(pks)
}

func (s *signature) Marshal() []byte {
	return s.p.Compress()
}

func (s *signature) Copy() types.Signature {
	np := *s.p
	return &signature{p: &np}
}

func Sign(sk types.SecretKey, msg []byte) types.Signature {
	sig := new(blstSignature).Sign(
		sk.(*secretKey).p, msg, types.DomainSeparationTag)
	return &signature{p: sig}
}

func Verify(sig types.Signature, msg []byte, pk types.PublicKey) bool {
	sigGroupCheck := false // alreaay checked in *SignatureFromBytes()
	pkValidate := false    // alreaay checked in *PublicKeyFromBytes()
	return sig.(*signature).p.Verify(
		sigGroupCheck, pk.(*publicKey).p, pkValidate, msg, types.DomainSeparationTag)
}

func FastAggregateVerify(sig types.Signature, msg []byte, pks []types.PublicKey) bool {
	pubPs := make([]*blstPublicKey, len(pks))
	for i := 0; i < len(pks); i++ {
		pubPs[i] = pks[i].(*publicKey).p
	}

	sigGroupCheck := false // alreaay checked in *SignatureFromBytes()
	return sig.(*signature).p.FastAggregateVerify(
		sigGroupCheck, pubPs, msg, types.DomainSeparationTag)
}

func VerifySignature(sigb []byte, msg [32]byte, pk types.PublicKey) (bool, error) {
	sig, err := SignatureFromBytes(sigb)
	if err != nil {
		return false, err
	}
	sigGroupCheck := false // alreaay checked in SignatureFromBytes()
	pkValidate := false    // alreaay checked in *PublicKeyFromBytes()
	ok := sig.(*signature).p.Verify(
		sigGroupCheck, pk.(*publicKey).p, pkValidate, msg[:], types.DomainSeparationTag)
	return ok, nil
}

// https://ethresear.ch/t/fast-verification-of-multiple-bls-signatures/5407
func VerifyMultipleSignatures(sigbs [][]byte, msgs [][32]byte, pks []types.PublicKey) (bool, error) {
	if len(sigbs) == 0 {
		return false, nil
	}
	if len(sigbs) != len(msgs) || len(msgs) != len(pks) {
		return false, errors.Errorf("Unmatched sigs, msgs, pubkeys lengths: %d, %d, %d", len(sigbs), len(msgs), len(pks))
	}
	count := len(sigbs)

	// Benefit from cache
	sigs, err := MultipleSignaturesFromBytes(sigbs)
	if err != nil {
		return false, err
	}

	// Convert to blst pointers
	var (
		sigPs = make([]*blstSignature, count)
		pkPs  = make([]*blstPublicKey, count)
		msgPs = make([]blstMessage, count)
	)
	for i := 0; i < len(sigs); i++ {
		sigPs[i] = sigs[i].(*signature).p
		pkPs[i] = pks[i].(*publicKey).p
		msgPs[i] = msgs[i][:]
	}

	// Supply random source
	var randErr error
	randLock := new(sync.Mutex)
	randFunc := func(outScalar *blstScalar) {
		var r [blstScalarBytes]byte
		randLock.Lock()
		_, err := rand.Read(r[:])
		randLock.Unlock()

		if err != nil {
			randErr = err
		}
		outScalar.FromBEndian(r[:])
	}

	dummy := new(blstSignature)
	sigGroupCheck := false // alreaay checked in MultipleSignaturesFromBytes()
	pkValidate := false    // alreaay checked in *PublicKeyFromBytes()
	ok := dummy.MultipleAggregateVerify(
		sigPs, sigGroupCheck, pkPs, pkValidate, msgPs, types.DomainSeparationTag, randFunc, blstRandBits)

	if randErr != nil {
		return false, fmt.Errorf("verify aborted: %s", randErr.Error())
	} else {
		return ok, nil
	}
}
