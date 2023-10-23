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
	"errors"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/log"
)

var (
	logger = log.NewModuleLogger(log.Blockchain)

	// Canonical system contract names registered in Registry.
	AddressBookName = "AddressBook"
	GovParamName    = "GovParam"
	Kip103Name      = "TreasuryRebalance"

	AllContractNames = []string{
		AddressBookName,
		GovParamName,
		Kip103Name,
	}

	// Some system contracts are allocated at special addresses.
	AddressBookAddr = common.HexToAddress("0x0000000000000000000000000000000000000400") // TODO: replace contracts/reward/contract/utils.go
	RegistryAddr    = common.HexToAddress("0x0000000000000000000000000000000000000401")

	// System contract binaries to be injected at hardfork or used in testing.
	RegistryCode     = hexutil.MustDecode("0x" + contracts.RegistryBinRuntime)
	RegistryMockCode = hexutil.MustDecode("0x" + contracts.RegistryMockBinRuntime)

	// Errors
	ErrRegistryNotInstalled = errors.New("Registry contract not installed")
)
