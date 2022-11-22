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
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
)

type IDeriveSha interface {
	DeriveSha(list types.DerivableList) common.Hash
}

type GovernanceEngine interface {
	ParamsAt(num uint64) (*params.GovParamSet, error)
}

var (
	config     *params.ChainConfig
	gov        GovernanceEngine
	impls      map[int]IDeriveSha
	emptyRoots map[int]common.Hash

	logger = log.NewModuleLogger(log.Blockchain)
)

func init() {
	impls = map[int]IDeriveSha{
		types.ImplDeriveShaOriginal: DeriveShaOrig{},
		types.ImplDeriveShaSimple:   DeriveShaSimple{},
		types.ImplDeriveShaConcat:   DeriveShaConcat{},
	}

	emptyRoots = make(map[int]common.Hash)
	for implType, impl := range impls {
		emptyRoots[implType] = impl.DeriveSha(types.Transactions{})
	}
}

func InitDeriveSha(chainConfig *params.ChainConfig, govEngine GovernanceEngine) {
	config = chainConfig
	gov = govEngine
	types.DeriveSha = DeriveShaMux
	types.EmptyRootHash = EmptyRootHashMux
	logger.Info("InitDeriveSha", "initial", config.DeriveShaImpl, "withGov", gov != nil)
}

func DeriveShaMux(list types.DerivableList, num *big.Int) common.Hash {
	return impls[getType(num)].DeriveSha(list)
}

func EmptyRootHashMux(num *big.Int) common.Hash {
	return emptyRoots[getType(num)]
}

func getType(num *big.Int) int {
	implType := config.DeriveShaImpl

	// gov == nil if blockchain.InitDeriveSha() is used, in genesis block manipulation
	// and unit tests. gov != nil if blockchain.InitDeriveShaWithGov is used,
	// in ordinary blockchain processing.
	if gov != nil {
		if pset, err := gov.ParamsAt(num.Uint64()); err != nil {
			logger.Crit("Cannot determine DeriveShaImpl", "num", num.Uint64(), "err", err)
		} else {
			implType = pset.DeriveShaImpl()
		}
	}

	if _, ok := impls[implType]; ok {
		return implType
	} else {
		logger.Error("Unrecognized deriveShaImpl, falling back to Orig", "impl", implType)
		return types.ImplDeriveShaOriginal
	}
}
