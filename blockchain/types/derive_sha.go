// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/types/derive_sha.go (2018/06/04).
// Modified and improved for the klaytn development.

package types

import (
	"github.com/klaytn/klaytn/common"
)

type DerivableList interface {
	Len() int
	GetRlp(i int) []byte
}

const (
	ImplDeriveShaOriginal int = iota
	ImplDeriveShaSimple
	ImplDeriveShaConcat
)

type IDeriveSha interface {
	DeriveSha(list DerivableList) common.Hash
}

var deriveShaObj IDeriveSha = nil

func InitDeriveSha(i IDeriveSha) {
	deriveShaObj = i

	// reset EmptyRootHash.
	EmptyRootHash = DeriveSha(Transactions{})
}

func DeriveSha(list DerivableList) common.Hash {
	return deriveShaObj.DeriveSha(list)
}
