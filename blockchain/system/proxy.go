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
	"github.com/klaytn/klaytn/common"
)

// This is the keccak-256 hash of "eip1967.proxy.implementation" subtracted by 1 used in the
// EIP-1967 proxy contract. See https://eips.ethereum.org/EIPS/eip-1967#implementation-slot
var ImplementationSlot = common.Hex2Bytes("360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")

func AllocProxy(impl common.Address) map[common.Hash]common.Hash {
	if impl == (common.Address{}) {
		return nil
	}
	storage := make(map[common.Hash]common.Hash)

	storage[common.BytesToHash(ImplementationSlot)] = lpad32(impl)

	return storage
}
