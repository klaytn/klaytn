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
	"github.com/klaytn/klaytn/blockchain/types/derivesha"
	"github.com/klaytn/klaytn/params"
)

// InitDeriveSha makes DeriveSha() and EmptyRootHash() depend on only ChainConfig.DeriveShaImpl.
// Call InitDeriveSha when you work exclusivly with genesis block (e.g. initGenesis)
func InitDeriveSha(config *params.ChainConfig) {
	derivesha.InitDeriveSha(config, nil)
}

// InitDeriveShaWithGov makes DeriveSha() and EmptyRootHash depend on the
// governance parameters. For any given block number, correct DeriveShaImpl will be used.
// Call InitDeriveShaWithGov before processing blocks.
func InitDeriveShaWithGov(config *params.ChainConfig, gov derivesha.GovernanceEngine) {
	derivesha.InitDeriveSha(config, gov)
}
