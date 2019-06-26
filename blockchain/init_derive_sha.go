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

package blockchain

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/storage/statedb"
)

func InitDeriveSha(deriveShaImpl int) {
	switch deriveShaImpl {
	case types.ImplDeriveShaOriginal:
		logger.Info("Using DeriveShaOrig!")
		types.InitDeriveSha(statedb.DeriveShaOrig{})
	case types.ImplDeriveShaSimple:
		logger.Info("Using DeriveShaSimple!")
		types.InitDeriveSha(types.DeriveShaSimple{})
	case types.ImplDeriveShaConcat:
		logger.Info("Using DeriveShaConcat!")
		types.InitDeriveSha(types.DeriveShaConcat{})
	default:
		logger.Error("Undefined deriveShaImpl!! use DeriveShaOrig!")
		types.InitDeriveSha(statedb.DeriveShaOrig{})
	}
}
