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
// Returns keccak256(h(key) . base) where h() differs by the key type.
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
func calcMappingSlot(baseSlot, key interface{}) common.Hash {
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

	return crypto.Keccak256Hash(append(keyBytes, baseBytes...))
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

// Calculate the contract storage slot for a struct element.
// Returns keccak256(base) + index
func calcStructSlot(baseSlot interface{}, index int) common.Hash {
	baseBytes := lpad32(baseSlot).Bytes()
	baseHash := crypto.Keccak256Hash(baseBytes)

	slot := new(big.Int).SetBytes(baseHash.Bytes())

	slot.Add(slot, big.NewInt(int64(index))) // slot += index

	return common.BytesToHash(slot.Bytes())
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
	default:
		logger.Crit("not a slot value type", "value", value)
		return common.Hash{}
	}
}

// Pad a dynamic string to a single slot.
// Only strings at most 31 bytes long are implemented here.
// See https://docs.soliditylang.org/en/v0.8.20/internals/layout_in_storage.html
func encodeShortString(s string) common.Hash {
	if len(s) > 31 {
		logger.Crit("encoding more than 31 byte string is not implemented")
		return common.Hash{}
	}

	// "In particular: if the data is at most 31 bytes long,
	// the elements are stored in the higher-order bytes (left aligned)
	// and the lowest-order byte stores the value length * 2."
	b := make([]byte, 32)
	copy(b, []byte(s))
	b[31] = byte(len(s)) * 2
	return common.BytesToHash(b)
}

// Add a int to a hash.
func addIntToHash(h common.Hash, i int) common.Hash {
	slot := new(big.Int).SetBytes(h.Bytes())

	slot.Add(slot, big.NewInt(int64(i))) // slot += index

	return common.BytesToHash(slot.Bytes())
}
