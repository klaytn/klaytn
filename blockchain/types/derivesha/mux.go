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

package derivesha

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

// TODO-Klaytn: Make DeriveShaMux state-less
// As-is: InitDS(type) + DS(list) + ERH
// To-be: InitDS() + DS(list, type) + ERH(type)

type IDeriveSha interface {
	DeriveSha(list types.DerivableList) common.Hash
}

var deriveShaObj IDeriveSha = nil
var logger = log.NewModuleLogger(log.Blockchain)

func InitDeriveSha(implType int) {
	switch implType {
	case types.ImplDeriveShaOriginal:
		deriveShaObj = DeriveShaOrig{}
	case types.ImplDeriveShaSimple:
		deriveShaObj = DeriveShaSimple{}
	case types.ImplDeriveShaConcat:
		deriveShaObj = DeriveShaConcat{}
	default:
		logger.Error("Unrecognized deriveShaImpl, falling back to Orig", "impl", implType)
		deriveShaObj = DeriveShaOrig{}
	}

	types.DeriveSha = DeriveShaMux
	types.EmptyRootHash = DeriveShaMux(types.Transactions{})
}

func DeriveShaMux(list types.DerivableList) common.Hash {
	return deriveShaObj.DeriveSha(list)
}
