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

package database

import (
	"math/big"
	"reflect"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
)

// NOTE-Klaytn-Cache BlockChain Caches
// Below is the list of the constants for cache size.
// TODO-Klaytn: Below should be handled by ini or other configurations.
const (
	maxHeaderCache        = 512
	maxTdCache            = 1024
	maxBlockNumberCache   = 2048
	maxCanonicalHashCache = 2048

	maxBodyCache            = 256
	maxBlockCache           = 256
	maxRecentTransactions   = 30000
	maxRecentBlockReceipts  = 30
	maxRecentTxReceipt      = 30000
	maxSenderTxHashToTxHash = 30000
)

const (
	numShardsHeaderCache        = 4096
	numShardsTdCache            = 4096
	numShardsBlockNumberCache   = 4096
	numShardsCanonicalHashCache = 4096

	numShardsBodyCache            = 4096
	numShardsBlockCache           = 4096
	numShardsRecentTransactions   = 4096
	numShardsRecentBlockReceipts  = 4096
	numShardsRecentTxReceipt      = 4096
	numShardsSenderTxHashToTxHash = 4096
)

type cacheKey int

const (
	headerCacheIndex cacheKey = iota
	tdCacheIndex
	blockNumberCacheIndex
	canonicalCacheIndex

	bodyCacheIndex
	bodyRLPCacheIndex
	blockCacheIndex
	recentTxAndLookupInfoIndex
	recentBlockReceiptsIndex
	recentTxReceiptIndex
	senderTxHashToTxHashIndex

	cacheKeySize
)

var lruCacheConfig = [cacheKeySize]common.CacheConfiger{
	headerCacheIndex:      common.LRUConfig{CacheSize: maxHeaderCache, IsScaled: true},
	tdCacheIndex:          common.LRUConfig{CacheSize: maxTdCache, IsScaled: true},
	blockNumberCacheIndex: common.LRUConfig{CacheSize: maxBlockNumberCache, IsScaled: true},
	canonicalCacheIndex:   common.LRUConfig{CacheSize: maxCanonicalHashCache, IsScaled: true},

	bodyCacheIndex:             common.LRUConfig{CacheSize: maxBodyCache, IsScaled: true},
	bodyRLPCacheIndex:          common.LRUConfig{CacheSize: maxBodyCache, IsScaled: true},
	blockCacheIndex:            common.LRUConfig{CacheSize: maxBlockCache, IsScaled: true},
	recentTxAndLookupInfoIndex: common.LRUConfig{CacheSize: maxRecentTransactions, IsScaled: true},
	recentBlockReceiptsIndex:   common.LRUConfig{CacheSize: maxRecentBlockReceipts, IsScaled: true},
	recentTxReceiptIndex:       common.LRUConfig{CacheSize: maxRecentTxReceipt, IsScaled: true},
	senderTxHashToTxHashIndex:  common.LRUConfig{CacheSize: maxSenderTxHashToTxHash, IsScaled: true},
}

var lruShardCacheConfig = [cacheKeySize]common.CacheConfiger{
	headerCacheIndex:      common.LRUShardConfig{CacheSize: maxHeaderCache, NumShards: numShardsHeaderCache, IsScaled: true},
	tdCacheIndex:          common.LRUShardConfig{CacheSize: maxTdCache, NumShards: numShardsTdCache, IsScaled: true},
	blockNumberCacheIndex: common.LRUShardConfig{CacheSize: maxBlockNumberCache, NumShards: numShardsBlockNumberCache, IsScaled: true},
	canonicalCacheIndex:   common.LRUShardConfig{CacheSize: maxCanonicalHashCache, NumShards: numShardsCanonicalHashCache, IsScaled: true},

	bodyCacheIndex:             common.LRUShardConfig{CacheSize: maxBodyCache, NumShards: numShardsBodyCache, IsScaled: true},
	bodyRLPCacheIndex:          common.LRUShardConfig{CacheSize: maxBodyCache, NumShards: numShardsBodyCache, IsScaled: true},
	blockCacheIndex:            common.LRUShardConfig{CacheSize: maxBlockCache, NumShards: numShardsBlockCache, IsScaled: true},
	recentTxAndLookupInfoIndex: common.LRUShardConfig{CacheSize: maxRecentTransactions, NumShards: numShardsRecentTransactions, IsScaled: true},
	recentBlockReceiptsIndex:   common.LRUShardConfig{CacheSize: maxRecentBlockReceipts, NumShards: numShardsRecentBlockReceipts, IsScaled: true},
	recentTxReceiptIndex:       common.LRUShardConfig{CacheSize: maxRecentTxReceipt, NumShards: numShardsRecentTxReceipt, IsScaled: true},
	senderTxHashToTxHashIndex:  common.LRUShardConfig{CacheSize: maxSenderTxHashToTxHash, NumShards: numShardsSenderTxHashToTxHash, IsScaled: true},
}

var fifoCacheConfig = [cacheKeySize]common.CacheConfiger{
	headerCacheIndex:      common.FIFOCacheConfig{CacheSize: maxHeaderCache, IsScaled: true},
	tdCacheIndex:          common.FIFOCacheConfig{CacheSize: maxTdCache, IsScaled: true},
	blockNumberCacheIndex: common.FIFOCacheConfig{CacheSize: maxBlockNumberCache, IsScaled: true},
	canonicalCacheIndex:   common.FIFOCacheConfig{CacheSize: maxCanonicalHashCache, IsScaled: true},

	bodyCacheIndex:             common.FIFOCacheConfig{CacheSize: maxBodyCache, IsScaled: true},
	bodyRLPCacheIndex:          common.FIFOCacheConfig{CacheSize: maxBodyCache, IsScaled: true},
	blockCacheIndex:            common.FIFOCacheConfig{CacheSize: maxBlockCache, IsScaled: true},
	recentTxAndLookupInfoIndex: common.FIFOCacheConfig{CacheSize: maxRecentTransactions, IsScaled: true},
	recentBlockReceiptsIndex:   common.FIFOCacheConfig{CacheSize: maxRecentBlockReceipts, IsScaled: true},
	recentTxReceiptIndex:       common.FIFOCacheConfig{CacheSize: maxRecentTxReceipt, IsScaled: true},
	senderTxHashToTxHashIndex:  common.FIFOCacheConfig{CacheSize: maxSenderTxHashToTxHash, IsScaled: true},
}

func newCache(cacheNameKey cacheKey, cacheType common.CacheType) common.Cache {
	var cache common.Cache

	switch cacheType {
	case common.FIFOCacheType:
		cache = common.NewCache(fifoCacheConfig[cacheNameKey])
	case common.LRUCacheType:
		cache = common.NewCache(lruCacheConfig[cacheNameKey])
	case common.LRUShardCacheType:
		cache = common.NewCache(lruShardCacheConfig[cacheNameKey])
	default:
		cache = common.NewCache(fifoCacheConfig[cacheNameKey])
	}
	return cache
}

type TransactionLookup struct {
	Tx *types.Transaction
	*TxLookupEntry
}

// cacheManager handles caches of data structures stored in Database.
// Previously, most of them were handled by blockchain.HeaderChain or
// blockchain.BlockChain.
type cacheManager struct {
	// caches from blockchain.HeaderChain
	headerCache        common.Cache
	tdCache            common.Cache
	blockNumberCache   common.Cache
	canonicalHashCache common.Cache

	// caches from blockchain.BlockChain
	bodyCache             common.Cache // Cache for the most recent block bodies
	bodyRLPCache          common.Cache // Cache for the most recent block bodies in RLP encoded format
	blockCache            common.Cache // Cache for the most recent entire blocks
	recentTxAndLookupInfo common.Cache // recent TX and LookupInfo cache
	recentBlockReceipts   common.Cache // recent block receipts cache
	recentTxReceipt       common.Cache // recent TX receipt cache

	senderTxHashToTxHashCache common.Cache
}

// newCacheManager returns a pointer of cacheManager with predefined configurations.
func newCacheManager() *cacheManager {
	cm := &cacheManager{
		headerCache:        newCache(headerCacheIndex, common.DefaultCacheType),
		tdCache:            newCache(tdCacheIndex, common.DefaultCacheType),
		blockNumberCache:   newCache(blockNumberCacheIndex, common.DefaultCacheType),
		canonicalHashCache: newCache(canonicalCacheIndex, common.DefaultCacheType),

		bodyCache:    newCache(bodyCacheIndex, common.DefaultCacheType),
		bodyRLPCache: newCache(bodyRLPCacheIndex, common.DefaultCacheType),
		blockCache:   newCache(blockCacheIndex, common.DefaultCacheType),

		recentTxAndLookupInfo: newCache(recentTxAndLookupInfoIndex, common.DefaultCacheType),
		recentBlockReceipts:   newCache(recentBlockReceiptsIndex, common.DefaultCacheType),
		recentTxReceipt:       newCache(recentTxReceiptIndex, common.DefaultCacheType),

		senderTxHashToTxHashCache: newCache(recentTxReceiptIndex, common.DefaultCacheType),
	}
	return cm
}

// clearHeaderChainCache flushes out 1) headerCache, 2) tdCache and 3) blockNumberCache.
func (cm *cacheManager) clearHeaderChainCache() {
	cm.headerCache.Purge()
	cm.tdCache.Purge()
	cm.blockNumberCache.Purge()
	cm.canonicalHashCache.Purge()
}

// clearBlockChainCache flushes out 1) bodyCache, 2) bodyRLPCache, 3) blockCache,
// 4) recentTxAndLookupInfo, 5) recentBlockReceipts and 6) recentTxReceipt.
func (cm *cacheManager) clearBlockChainCache() {
	cm.bodyCache.Purge()
	cm.bodyRLPCache.Purge()
	cm.blockCache.Purge()
	cm.recentTxAndLookupInfo.Purge()
	cm.recentBlockReceipts.Purge()
	cm.recentTxReceipt.Purge()
	cm.senderTxHashToTxHashCache.Purge()
}

// readHeaderCache looks for cached header in headerCache.
// It returns nil if not found.
func (cm *cacheManager) readHeaderCache(hash common.Hash) *types.Header {
	if header, ok := cm.headerCache.Get(hash); ok && header != nil {
		cacheGetHeaderHitMeter.Mark(1)
		return header.(*types.Header)
	}
	cacheGetHeaderMissMeter.Mark(1)
	return nil
}

// writeHeaderCache writes header as a value, headerHash as a key.
func (cm *cacheManager) writeHeaderCache(hash common.Hash, header *types.Header) {
	if header == nil {
		return
	}
	cm.headerCache.Add(hash, header)
}

// deleteHeaderCache writes nil as a value, headerHash as a key, to indicate given
// headerHash is deleted in headerCache.
func (cm *cacheManager) deleteHeaderCache(hash common.Hash) {
	cm.headerCache.Add(hash, nil)
}

// hasHeaderInCache returns if a cachedHeader exists with given headerHash.
func (cm *cacheManager) hasHeaderInCache(hash common.Hash) bool {
	if cached, ok := cm.headerCache.Get(hash); ok && cached != nil {
		return true
	}
	return false
}

// readTdCache looks for cached total blockScore in tdCache.
// It returns nil if not found.
func (cm *cacheManager) readTdCache(hash common.Hash) *big.Int {
	if cached, ok := cm.tdCache.Get(hash); ok && cached != nil {
		cacheGetTDHitMeter.Mark(1)
		return cached.(*big.Int)
	}
	cacheGetTDMissMeter.Mark(1)
	return nil
}

// writeHeaderCache writes total blockScore as a value, headerHash as a key.
func (cm *cacheManager) writeTdCache(hash common.Hash, td *big.Int) {
	if td == nil {
		return
	}
	cm.tdCache.Add(hash, td)
}

// deleteTdCache writes nil as a value, headerHash as a key, to indicate given
// headerHash is deleted in TdCache.
func (cm *cacheManager) deleteTdCache(hash common.Hash) {
	cm.tdCache.Add(hash, nil)
}

// readBlockNumberCache looks for cached headerNumber in blockNumberCache.
// It returns nil if not found.
func (cm *cacheManager) readBlockNumberCache(hash common.Hash) *uint64 {
	if cached, ok := cm.blockNumberCache.Get(hash); ok && cached != nil {
		cacheGetBlockNumberHitMeter.Mark(1)
		blockNumber := cached.(uint64)
		return &blockNumber
	}
	cacheGetBlockNumberMissMeter.Mark(1)
	return nil
}

// writeHeaderCache writes headerNumber as a value, headerHash as a key.
func (cm *cacheManager) writeBlockNumberCache(hash common.Hash, number uint64) {
	cm.blockNumberCache.Add(hash, number)
}

// deleteBlockNumberCache deletes headerNumber with a headerHash as a key.
func (cm *cacheManager) deleteBlockNumberCache(hash common.Hash) {
	cm.blockNumberCache.Add(hash, nil)
}

// readCanonicalHashCache looks for cached canonical hash in canonicalHashCache.
// It returns empty hash if not found.
func (cm *cacheManager) readCanonicalHashCache(number uint64) common.Hash {
	if cached, ok := cm.canonicalHashCache.Get(common.CacheKeyUint64(number)); ok {
		cacheGetCanonicalHashHitMeter.Mark(1)
		canonicalHash := cached.(common.Hash)
		return canonicalHash
	}
	cacheGetCanonicalHashMissMeter.Mark(1)
	return common.Hash{}
}

// writeCanonicalHashCache writes canonical hash as a value, headerNumber as a key.
func (cm *cacheManager) writeCanonicalHashCache(number uint64, hash common.Hash) {
	cm.canonicalHashCache.Add(common.CacheKeyUint64(number), hash)
}

// readBodyCache looks for cached blockBody in bodyCache.
// It returns nil if not found.
func (cm *cacheManager) readBodyCache(hash common.Hash) *types.Body {
	if cachedBody, ok := cm.bodyCache.Get(hash); ok && cachedBody != nil {
		cacheGetBlockBodyHitMeter.Mark(1)
		return cachedBody.(*types.Body)
	}
	cacheGetBlockBodyMissMeter.Mark(1)
	return nil
}

// writeBodyCache writes blockBody as a value, blockHash as a key.
func (cm *cacheManager) writeBodyCache(hash common.Hash, body *types.Body) {
	if body == nil {
		return
	}
	cm.bodyCache.Add(hash, body)
}

// deleteBodyCache writes nil as a value, blockHash as a key, to indicate given
// txHash is deleted in bodyCache and bodyRLPCache.
func (cm *cacheManager) deleteBodyCache(hash common.Hash) {
	cm.bodyCache.Add(hash, nil)
	cm.bodyRLPCache.Add(hash, nil)
}

// readBodyRLPCache looks for cached RLP-encoded blockBody in bodyRLPCache.
// It returns nil if not found.
func (cm *cacheManager) readBodyRLPCache(hash common.Hash) rlp.RawValue {
	if cachedBodyRLP, ok := cm.bodyRLPCache.Get(hash); ok && cachedBodyRLP != nil {
		cacheGetBlockBodyRLPHitMeter.Mark(1)
		return cachedBodyRLP.(rlp.RawValue)
	}
	cacheGetBlockBodyRLPMissMeter.Mark(1)
	return nil
}

// writeBodyRLPCache writes RLP-encoded blockBody as a value, blockHash as a key.
func (cm *cacheManager) writeBodyRLPCache(hash common.Hash, bodyRLP rlp.RawValue) {
	if bodyRLP == nil {
		return
	}
	cm.bodyRLPCache.Add(hash, bodyRLP)
}

// readBlockCache looks for cached block in blockCache.
// It returns nil if not found.
func (cm *cacheManager) readBlockCache(hash common.Hash) *types.Block {
	if cachedBlock, ok := cm.blockCache.Get(hash); ok && cachedBlock != nil {
		cacheGetBlockHitMeter.Mark(1)
		return cachedBlock.(*types.Block)
	}
	cacheGetBlockMissMeter.Mark(1)
	return nil
}

// hasBlockInCache returns if given hash exists in blockCache.
func (cm *cacheManager) hasBlockInCache(hash common.Hash) bool {
	if cachedBlock, ok := cm.blockCache.Get(hash); ok && cachedBlock != nil {
		return true
	}
	return false
}

// writeBlockCache writes block as a value, blockHash as a key.
func (cm *cacheManager) writeBlockCache(hash common.Hash, block *types.Block) {
	if block == nil {
		return
	}
	cm.blockCache.Add(hash, block)
}

// deleteBlockCache writes nil as a value, blockHash as a key,
// to indicate given blockHash is deleted in recentBlockReceipts.
func (cm *cacheManager) deleteBlockCache(hash common.Hash) {
	cm.blockCache.Add(hash, nil)
}

// readTxAndLookupInfoInCache looks for cached tx and its look up information in recentTxAndLookupInfo.
// It returns nil and empty values if not found.
func (cm *cacheManager) readTxAndLookupInfoInCache(txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	if value, ok := cm.recentTxAndLookupInfo.Get(txHash); ok && value != nil {
		cacheGetRecentTransactionsHitMeter.Mark(1)
		txLookup, ok := value.(*TransactionLookup)
		if !ok {
			logger.Error("invalid type in recentTxAndLookupInfo. expected=*TransactionLookup", "actual=", reflect.TypeOf(value))
			return nil, common.Hash{}, 0, 0
		}
		return txLookup.Tx, txLookup.BlockHash, txLookup.BlockIndex, txLookup.Index
	}
	cacheGetRecentTransactionsMissMeter.Mark(1)
	return nil, common.Hash{}, 0, 0
}

// writeTxAndLookupInfoCache writes a tx and its lookup information as a value, txHash as a key.
func (cm *cacheManager) writeTxAndLookupInfoCache(txHash common.Hash, txLookup *TransactionLookup) {
	if txLookup == nil {
		return
	}
	cm.recentTxAndLookupInfo.Add(txHash, txLookup)
}

// readBlockReceiptsInCache looks for cached blockReceipts in recentBlockReceipts.
// It returns nil if not found.
func (cm *cacheManager) readBlockReceiptsInCache(blockHash common.Hash) types.Receipts {
	if cachedBlockReceipts, ok := cm.recentBlockReceipts.Get(blockHash); ok && cachedBlockReceipts != nil {
		cacheGetRecentBlockReceiptsHitMeter.Mark(1)
		return cachedBlockReceipts.(types.Receipts)
	}
	cacheGetRecentBlockReceiptsMissMeter.Mark(1)
	return nil
}

// writeBlockReceiptsCache writes blockReceipts as a value, blockHash as a key.
func (cm *cacheManager) writeBlockReceiptsCache(blockHash common.Hash, receipts types.Receipts) {
	if receipts == nil {
		return
	}
	cm.recentBlockReceipts.Add(blockHash, receipts)
}

// deleteBlockReceiptsCache writes nil as a value, blockHash as a key, to indicate given
// blockHash is deleted in recentBlockReceipts.
func (cm *cacheManager) deleteBlockReceiptsCache(blockHash common.Hash) {
	cm.recentBlockReceipts.Add(blockHash, nil)
}

// readTxReceiptInCache looks for cached txReceipt in recentTxReceipt.
// It returns nil if not found.
func (cm *cacheManager) readTxReceiptInCache(txHash common.Hash) *types.Receipt {
	if cachedReceipt, ok := cm.recentTxReceipt.Get(txHash); ok && cachedReceipt != nil {
		cacheGetRecentTxReceiptHitMeter.Mark(1)
		return cachedReceipt.(*types.Receipt)
	}
	cacheGetRecentTxReceiptMissMeter.Mark(1)
	return nil
}

// writeTxReceiptCache writes txReceipt as a value, txHash as a key.
func (cm *cacheManager) writeTxReceiptCache(txHash common.Hash, receipt *types.Receipt) {
	if receipt == nil {
		return
	}
	cm.recentTxReceipt.Add(txHash, receipt)
}

// deleteTxReceiptCache writes nil as a value, blockHash as a key, to indicate given
// txHash is deleted in recentTxReceipt.
func (cm *cacheManager) deleteTxReceiptCache(txHash common.Hash) {
	cm.recentTxReceipt.Add(txHash, nil)
}

// writeSenderTxHashToTxHashCache writes senderTxHash to txHash mapping information to cache.
func (cm *cacheManager) writeSenderTxHashToTxHashCache(senderTxHash, txHash common.Hash) {
	cm.senderTxHashToTxHashCache.Add(senderTxHash, txHash)
}

// readSenderTxHashToTxHashCache looks for matching txHash from senderTxHash.
// If txHash does not exist in the cache, it returns an empty hash.
func (cm *cacheManager) readSenderTxHashToTxHashCache(senderTxHash common.Hash) common.Hash {
	if matchedTxHash, ok := cm.senderTxHashToTxHashCache.Get(senderTxHash); ok && matchedTxHash != nil {
		return matchedTxHash.(common.Hash)
	}
	return common.Hash{}
}
