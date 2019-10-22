// Modifications Copyright 2018 The klaytn Authors
// Copyright 2018 The go-ethereum Authors
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
// This file is derived from core/rawdb/schema.go (2018/06/04).
// Modified and improved for the klaytn development.

package database

import (
	"encoding/binary"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/metrics"
)

// The fields below define the low level database schema prefixing.
var (
	// databaseVerisionKey tracks the current database version.
	databaseVerisionKey = []byte("DatabaseVersion")

	// headHeaderKey tracks the latest know header's hash.
	headHeaderKey = []byte("LastHeader")

	// headBlockKey tracks the latest know full block's hash.
	headBlockKey = []byte("LastBlock")

	// headFastBlockKey tracks the latest known incomplete block's hash duirng fast sync.
	headFastBlockKey = []byte("LastFast")

	// fastTrieProgressKey tracks the number of trie entries imported during fast sync.
	fastTrieProgressKey = []byte("TrieSync")

	validSectionKey = []byte("count")

	sectionHeadKeyPrefix = []byte("shead")

	snapshotKeyPrefix = []byte("snapshot")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`, used for indexes).
	headerPrefix       = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	headerTDSuffix     = []byte("t") // headerPrefix + num (uint64 big endian) + hash + headerTDSuffix -> td
	headerHashSuffix   = []byte("n") // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	headerNumberPrefix = []byte("H") // headerNumberPrefix + hash -> num (uint64 big endian)

	blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts

	txLookupPrefix = []byte("l") // txLookupPrefix + hash -> transaction/receipt lookup metadata

	preimagePrefix = []byte("secure-key-")  // preimagePrefix + hash -> preimage
	configPrefix   = []byte("klay-config-") // config prefix for the db

	// Chain index prefixes (use `i` + single byte to avoid mixing data types).
	BloomBitsIndexPrefix = []byte("iB") // BloomBitsIndexPrefix is the data table of a chain indexer to track its progress

	preimageCounter    = metrics.NewRegisteredCounter("db/preimage/total", nil)
	preimageHitCounter = metrics.NewRegisteredCounter("db/preimage/hits", nil)

	childChainTxHashPrefix          = []byte("ccTxHash")
	lastServiceChainTxReceiptKey    = []byte("LastServiceChainTxReceipt")
	lastIndexedBlockKey             = []byte("LastIndexedBlockKey")
	receiptFromParentChainKeyPrefix = []byte("receiptFromParentChain")

	valueTransferTxHashPrefix = []byte("vt-tx-hash-key-") // Prefix + hash -> hash

	// bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits
	bloomBitsPrefix = []byte("B")

	senderTxHashToTxHashPrefix = []byte("SenderTxHash")

	governancePrefix     = []byte("governance")
	governanceHistoryKey = []byte("governanceIdxHistory")
	governanceStateKey   = []byte("governanceState")

	databaseDirPrefix  = []byte("databaseDirectory")
	migrationStatusKey = []byte("migrationStatus")
)

// TxLookupEntry is a positional metadata to help looking up the data content of
// a transaction or receipt given only its hash.
type TxLookupEntry struct {
	BlockHash  common.Hash
	BlockIndex uint64
	Index      uint64
}

// encodeUint64 encodes a number as big endian uint64
func encodeUint64(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerKey = headerPrefix + num (uint64 big endian) + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(headerPrefix, encodeUint64(number)...), hash.Bytes()...)
}

// headerTDKey = headerPrefix + num (uint64 big endian) + hash + headerTDSuffix
func headerTDKey(number uint64, hash common.Hash) []byte {
	return append(headerKey(number, hash), headerTDSuffix...)
}

// headerHashKey = headerPrefix + num (uint64 big endian) + headerHashSuffix
func headerHashKey(number uint64) []byte {
	return append(append(headerPrefix, encodeUint64(number)...), headerHashSuffix...)
}

// headerNumberKey = headerNumberPrefix + hash
func headerNumberKey(hash common.Hash) []byte {
	return append(headerNumberPrefix, hash.Bytes()...)
}

// blockBodyKey = blockBodyPrefix + num (uint64 big endian) + hash
func blockBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(blockBodyPrefix, encodeUint64(number)...), hash.Bytes()...)
}

// blockReceiptsKey = blockReceiptsPrefix + num (uint64 big endian) + hash
func blockReceiptsKey(number uint64, hash common.Hash) []byte {
	return append(append(blockReceiptsPrefix, encodeUint64(number)...), hash.Bytes()...)
}

// TxLookupKey = txLookupPrefix + hash
func TxLookupKey(hash common.Hash) []byte {
	return append(txLookupPrefix, hash.Bytes()...)
}

func SenderTxHashToTxHashKey(senderTxHash common.Hash) []byte {
	return append(senderTxHashToTxHashPrefix, senderTxHash.Bytes()...)
}

// preimageKey = preimagePrefix + hash
func preimageKey(hash common.Hash) []byte {
	return append(preimagePrefix, hash.Bytes()...)
}

// configKey = configPrefix + hash
func configKey(hash common.Hash) []byte {
	return append(configPrefix, hash.Bytes()...)
}

func sectionHeadKey(encodedSection []byte) []byte {
	return append(sectionHeadKeyPrefix, encodedSection...)
}

func snapshotKey(hash common.Hash) []byte {
	return append(snapshotKeyPrefix, hash[:]...)
}

func childChainTxHashKey(ccBlockHash common.Hash) []byte {
	return append(append(childChainTxHashPrefix, ccBlockHash.Bytes()...))
}

func receiptFromParentChainKey(blockHash common.Hash) []byte {
	return append(receiptFromParentChainKeyPrefix, blockHash.Bytes()...)
}

func valueTransferTxHashKey(rTxHash common.Hash) []byte {
	return append(valueTransferTxHashPrefix, rTxHash.Bytes()...)
}

// bloomBitsKey = bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash
func BloomBitsKey(bit uint, section uint64, hash common.Hash) []byte {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), hash.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return key
}

func governanceKey(num uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return append(governancePrefix[:], b[:]...)
}

func databaseDirKey(dbEntryType uint64) []byte {
	return append(databaseDirPrefix, encodeUint64(dbEntryType)...)
}
