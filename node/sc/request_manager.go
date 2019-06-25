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

package sc

import (
	"crypto/ecdsa"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"math/big"
)

const (
	GasLimit = 5000000
)

// TODO-Klaytn need to refactor for chainID / gasPrice
func MakeTransactOpts(accountKey *ecdsa.PrivateKey, nonce *big.Int, chainID *big.Int, gasPrice *big.Int) *bind.TransactOpts {
	if accountKey == nil {
		return nil
	}
	auth := bind.NewKeyedTransactor(accountKey)
	auth.GasLimit = GasLimit
	auth.GasPrice = gasPrice
	auth.Nonce = nonce
	auth.Signer = func(signer types.Signer, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, signer, accountKey)
	}
	return auth
}
