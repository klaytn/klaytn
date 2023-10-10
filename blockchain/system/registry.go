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

type RegistryRecord struct {
	Addr       common.Address
	Activation *big.Int
}

// AllocRegistryInit Specifies the initial Registry state.
// This struct only represents a special case of Registry state where:
// - there is only one record for each name
// - the activation of all records is zero
// - the names array is lexicographically sorted
type AllocRegistryInit struct {
	Contracts map[string]common.Address
	Owner     common.Address
}

func NewAllocRegistryInit() *AllocRegistryInit {
	return &AllocRegistryInit{
		Contracts: make(map[string]common.Address),
	}
}

// Create storage state from the given initial values.
// The storage slots are calculated according to the solidity layout rule.
// https://docs.soliditylang.org/en/v0.8.20/internals/layout_in_storage.html
func AllocRegistry(init *AllocRegistryInit) map[common.Hash]common.Hash {
	if init == nil {
		return nil
	}
	storage := make(map[common.Hash]common.Hash)

	// slot[0]: mapping(string => Record[]) records;
	// In AllocRegistry, records[name] is always Record[] of one element.
	// - records[x].length @ Hash(x, 0)
	// - records[x][i].addr @ Hash(Hash(x, 0)) + (2*i)
	// - records[x][i].activation @ Hash(Hash(x, 0)) + (2*i + 1)
	for name, addr := range init.Contracts {
		arraySlot := calcMappingSlot(0, name)
		storage[arraySlot] = lpad32(1)

		addrSlot := calcArraySlot(arraySlot, 2, 0, 0)
		activationSlot := calcArraySlot(arraySlot, 2, 0, 1)

		storage[addrSlot] = lpad32(addr)
		storage[activationSlot] = lpad32(0)
	}

	names := make([]string, 0)
	for name, _ := range init.Contracts {
		names = append(names, name)
	}
	sort.Strings(names)

	// slot[1]: string[] names;
	// - names.length @ 1
	// - names[i] @ Hash(1) + i
	storage[lpad32(1)] = lpad32(len(names))
	for i, name := range names {
		nameSlot := calcArraySlot(1, 1, i, 0)
		storage[nameSlot] = encodeShortString(name)
	}

	// slot[2]: address _owner;
	storage[lpad32(2)] = lpad32(init.Owner)

	return storage
}

func InstallRegistry(state *state.StateDB, init *AllocRegistryInit) error {
	if err := state.SetCode(RegistryAddr, RegistryCode); err != nil {
		return err
	}
	storage := AllocRegistry(init)
	for key, value := range storage {
		state.SetState(RegistryAddr, key, value)
	}
	return nil
}

// Install Registry at the state with the initial records.
// The initial records contains pre-deployed system contracts that could have
// existed before the Cancun hardfork.
//
// - "AddressBook" at the constant AddressBookAdddr if code exists
// - "GovParam" at govParamAddr if nonzero
// - "TreasuryRebalance" at config.Kip103ContractAddress if nonzero
//
// Because these system contracts must exist in the Registry right after the
// Cancun hardfork block number, add them to Registry by modifying the state.
func InstallRegistryAtCancunFork(state *state.StateDB, config *params.ChainConfig, govParamAddr common.Address) error {
	return InstallRegistry(state, cancunRegistryInit(state, config, govParamAddr))
}

func cancunRegistryInit(state *state.StateDB, config *params.ChainConfig, govParamAddr common.Address) *AllocRegistryInit {
	init := NewAllocRegistryInit()

	if len(state.GetCode(AddressBookAddr)) > 0 {
		init.Contracts[AddressBookName] = AddressBookAddr
	}

	if (govParamAddr != common.Address{}) {
		init.Contracts[GovParamName] = govParamAddr
	}

	if (config.Kip103ContractAddress != common.Address{}) {
		init.Contracts[Kip103Name] = config.Kip103ContractAddress
	}

	return init
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
	return caller.GetActiveAddr(opts, name)
}
