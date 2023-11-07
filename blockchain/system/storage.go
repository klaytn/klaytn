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
	"math/big"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
)

// Calculate the contract storage slot for a mapping element.
// Returns keccak256(h(key) . base) + offset where h() differs by the key type.
//
// "The value corresponding to a mapping key k is located at keccak256(h(k) . p)
// where . is concatenation and h is a function that is applied to the key depending
// on its type"
//
// The baseSlot (p) can be of type:
// - *big.Int or common.Hash to represent a slot index
//
// The key can be of type:
// - *big.Int to represent Solidity uint8..uint256
// - common.Address to represent Solidity address
// - common.Hash to represent Solidity bytes1..bytes32, or another slot index
// - []byte to represent Solidity bytes
// - string to represent Solidity string
//
// See https://docs.soliditylang.org/en/v0.8.20/internals/layout_in_storage.html
func calcMappingSlot(baseSlot, key interface{}, offset int) common.Hash {
	baseBytes := lpad32(baseSlot).Bytes()

	var keyBytes []byte
	switch v := key.(type) {
	// For bytes and string types, h(k) = k, unpadded.
	case []byte:
		keyBytes = v
	case string:
		keyBytes = []byte(v)
	// For value types, h(k) = lpad32(k)
	default:
		keyBytes = lpad32(key).Bytes()
	}

	elemHash := crypto.Keccak256Hash(append(keyBytes, baseBytes...))

	slot := new(big.Int).SetBytes(elemHash.Bytes())
	slot.Add(slot, big.NewInt(int64(offset))) // slot += offset
	return common.BytesToHash(slot.Bytes())
}

// Calculate the contract storage slot for a dynamic array element.
// Returns keccak256(base) + (elemSize * index + offset)
//
// The baseSlot (p) can be of type:
// - *big.Int or common.Hash to represent a slot index
//
// "Array data is located starting at keccak256(p) and it is laid out in the same way
// as statically-sized array data would: One element after the other"
//
// See https://docs.soliditylang.org/en/v0.8.20/internals/layout_in_storage.html
func calcArraySlot(baseSlot interface{}, elemSize, index, offset int) common.Hash {
	baseBytes := lpad32(baseSlot).Bytes()
	baseHash := crypto.Keccak256Hash(baseBytes)

	slot := new(big.Int).SetBytes(baseHash.Bytes())

	mult := new(big.Int).Mul(
		big.NewInt(int64(elemSize)),
		big.NewInt(int64(index)))
	slot.Add(slot, mult) // slot += size * index

	slot.Add(slot, big.NewInt(int64(offset))) // slot += offset

	return common.BytesToHash(slot.Bytes())
}

// Allocates a dynamic string or bytes data that begins at `baseSlot`.
//   - If the data is less than 32-bytes long,
//     @ base         [data..., 0..., len*2+0] // even number indicates short string
//   - If the data is 32-bytes or longer,
//     @ base         [0...,          len*2+1] // odd number indicates long string
//     @ H(base) + i  [ data[i*32 : i*32+32] ] // right padded (left aligned)
func allocDynamicData(baseSlot interface{}, data []byte) map[common.Hash]common.Hash {
	storage := make(map[common.Hash]common.Hash)
	base := lpad32(baseSlot)

	// Short string
	if len(data) < 32 {
		slice := make([]byte, 32)
		copy(slice[:], data)
		slice[31] = byte(len(data) * 2)

		storage[base] = common.BytesToHash(slice)
		return storage
	}

	// Long string
	bigLen := big.NewInt(int64(len(data)*2 + 1))
	storage[base] = common.BytesToHash(bigLen.Bytes())

	baseHash := crypto.Keccak256Hash(base.Bytes())
	bigSlot := new(big.Int).SetBytes(baseHash.Bytes())
	for len(data) >= 32 {
		slot := common.BytesToHash(bigSlot.Bytes())
		storage[slot] = common.BytesToHash(data[0:32])

		data = data[32:]
		bigSlot = new(big.Int).Add(bigSlot, common.Big1)
	}
	if len(data) > 0 {
		slice := make([]byte, 32)
		copy(slice[:], data)

		slot := common.BytesToHash(bigSlot.Bytes())
		storage[slot] = common.BytesToHash(slice)
	}

	return storage
}

// MergeStorage merges multiple storage maps into one.
func MergeStorage(ss ...map[common.Hash]common.Hash) map[common.Hash]common.Hash {
	out := make(map[common.Hash]common.Hash)
	for _, s := range ss {
		for k, v := range s {
			if _, ok := out[k]; ok {
				logger.Crit("storage slot collision", "slot", k)
			} else {
				out[k] = v
			}
		}
	}
	return out
}

// Pad Solidity "value types" to 32 bytes.
// Only a few value types are implemented here.
// See https://docs.soliditylang.org/en/v0.8.20/types.html
func lpad32(value interface{}) common.Hash {
	switch v := value.(type) {
	case int:
		return common.BytesToHash(big.NewInt(int64(v)).Bytes())
	case uint64:
		return common.BytesToHash(big.NewInt(int64(v)).Bytes())
	case *big.Int:
		return common.BytesToHash(v.Bytes())
	case common.Address:
		return common.BytesToHash(v.Bytes())
	case common.Hash:
		return v
	case []byte:
		// Use allocDynamicData() for dynamic bytes types
		return common.BytesToHash(v)
	default:
		logger.Crit("not a slot value type", "value", value)
		return common.Hash{}
	}
}
