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

package system

import (
	"context"
	"math/big"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/common"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
)

func InstallRegistry(state *state.StateDB) error {
	return state.SetCode(RegistryAddr, RegistryCode)
}

func ReadRegistryActiveAddr(backend bind.ContractCaller, name string, num *big.Int) (common.Address, error) {
	caller, err := contracts.NewRegistryCaller(RegistryAddr, backend)
	if err != nil {
		return common.Address{}, err
	}

	code, err := backend.CodeAt(context.Background(), RegistryAddr, nil)
	if err != nil {
		return common.Address{}, err
	}
	if code == nil {
		return common.Address{}, ErrRegistryNotInstalled
	}

	opts := &bind.CallOpts{BlockNumber: num}
	return caller.GetActiveContract(opts, name)
}
