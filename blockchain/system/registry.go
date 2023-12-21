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
	"sort"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/common"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/params"
)

// Create storage state from the given initial values.
// The storage slots are calculated according to the solidity layout rule.
// https://docs.soliditylang.org/en/v0.8.20/internals/layout_in_storage.html
func AllocRegistry(init *params.RegistryConfig) map[common.Hash]common.Hash {
	if init == nil {
		return nil
	}
	if init.Records == nil {
		init.Records = make(map[string]common.Address)
	}
	storage := make(map[common.Hash]common.Hash)

	// slot[0]: mapping(string => Record[]) records;
	// In AllocRegistry, records[name] is always Record[] of one element.
	// - records[x].length @ Hash(x, 0)
	// - records[x][i].addr @ Hash(Hash(x, 0)) + (2*i)
	// - records[x][i].activation @ Hash(Hash(x, 0)) + (2*i + 1)
	for name, addr := range init.Records {
		arraySlot := calcMappingSlot(0, name, 0)
		storage[arraySlot] = lpad32(1) // records[name].length = 1

		addrSlot := calcArraySlot(arraySlot, 2, 0, 0)
		activationSlot := calcArraySlot(arraySlot, 2, 0, 1)

		storage[addrSlot] = lpad32(addr)    // records[name][0].addr
		storage[activationSlot] = lpad32(0) // records[name][0].activation
	}

	names := make([]string, 0)
	for name := range init.Records {
		names = append(names, name)
	}
	sort.Strings(names)

	// slot[1]: string[] names;
	// - names.length @ 1
	// - names[i] @ Hash(1) + i
	storage[lpad32(1)] = lpad32(len(names))
	for i, name := range names {
		nameSlot := calcArraySlot(1, 1, i, 0) // Hash(1) + 1*i + 0
		for k, v := range allocDynamicData(nameSlot, []byte(name)) {
			storage[k] = v
		}
	}

	// slot[2]: address _owner;
	storage[lpad32(2)] = lpad32(init.Owner)

	return storage
}

func InstallRegistry(state *state.StateDB, init *params.RegistryConfig) error {
	if err := state.SetCode(RegistryAddr, RegistryCode); err != nil {
		return err
	}
	storage := AllocRegistry(init)
	for key, value := range storage {
		state.SetState(RegistryAddr, key, value)
	}
	return nil
}

func ReadActiveAddressFromRegistry(backend bind.ContractCaller, name string, num *big.Int) (common.Address, error) {
	code, err := backend.CodeAt(context.Background(), RegistryAddr, num)
	if err != nil {
		return common.Address{}, err
	}
	if code == nil {
		return common.Address{}, ErrRegistryNotInstalled
	}

	caller, err := contracts.NewRegistryCaller(RegistryAddr, backend)
	if err != nil {
		return common.Address{}, err
	}

	opts := &bind.CallOpts{BlockNumber: num}
	return caller.GetActiveAddr(opts, name)
}
