// Copyright 2020 The klaytn Authors
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

package kas

import (
	"fmt"
	"strings"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

// TxFilteringTypes filters types which are only stored in KAS database.
var TxFilteringTypes = map[types.TxType]bool{
	types.TxTypeLegacyTransaction: true,

	types.TxTypeValueTransfer:                      true,
	types.TxTypeFeeDelegatedValueTransfer:          true,
	types.TxTypeFeeDelegatedValueTransferWithRatio: true,

	types.TxTypeValueTransferMemo:                      true,
	types.TxTypeFeeDelegatedValueTransferMemo:          true,
	types.TxTypeFeeDelegatedValueTransferMemoWithRatio: true,

	types.TxTypeSmartContractDeploy:                      true,
	types.TxTypeFeeDelegatedSmartContractDeploy:          true,
	types.TxTypeFeeDelegatedSmartContractDeployWithRatio: true,

	types.TxTypeSmartContractExecution:                      true,
	types.TxTypeFeeDelegatedSmartContractExecution:          true,
	types.TxTypeFeeDelegatedSmartContractExecutionWithRatio: true,
}

// transformToTxs transforms a chain event to transaction list according to the KAS database scheme.
func transformToTxs(event blockchain.ChainEvent) ([]*Tx, map[common.Address]struct{}) {
	var txs []*Tx
	updatedEOAs := make(map[common.Address]struct{})

	block := event.Block
	head := block.Header()
	receipts := event.Receipts

	for idx, rawTx := range block.Transactions() {
		txId := head.Number.Int64()*maxTxCountPerBlock*maxTxLogCountPerTx + int64(idx)*maxInternalTxCountPerTx

		// from
		var from common.Address
		if rawTx.IsEthereumTransaction() {
			signer := types.LatestSignerForChainID(rawTx.ChainId())
			from, _ = types.Sender(signer, rawTx)
		} else {
			from, _ = rawTx.From()
		}

		// to
		to := rawTx.To()
		if to == nil {
			to = &common.Address{}
		}

		// value
		value := hexutil.EncodeBig(rawTx.Value())

		// status
		var status int
		if receipts[idx].Status != types.ReceiptStatusSuccessful {
			status = int(types.ReceiptStatusFailed)
		} else {
			status = int(receipts[idx].Status)
		}

		// fee payer
		// fee ratio
		var (
			feePayer []byte
			feeRatio uint
		)

		if rawTx.IsFeeDelegatedTransaction() {
			payer, _ := rawTx.FeePayer()
			feePayer = payer.Bytes()

			ratio, ok := rawTx.FeeRatio()
			if ok {
				feeRatio = uint(ratio)
			}
		}

		updatedEOAs[from] = struct{}{}
		updatedEOAs[*to] = struct{}{}

		tx := &Tx{
			Timestamp:       head.Time.Int64(),
			TransactionId:   txId,
			FromAddr:        from.Bytes(),
			ToAddr:          to.Bytes(),
			Value:           value,
			TransactionHash: rawTx.Hash().Bytes(),
			Status:          status,
			TypeInt:         int(rawTx.Type()),
			GasPrice:        rawTx.GasPrice().Uint64(),
			GasUsed:         receipts[idx].GasUsed,
			FeePayer:        feePayer,
			FeeRatio:        feeRatio,
		}

		txs = append(txs, tx)
	}
	return txs, updatedEOAs
}

// InsertTransactions inserts transactions in the given chain event into KAS database.
func (r *repository) InsertTransactions(event blockchain.ChainEvent) error {
	txs, updatedEOAs := transformToTxs(event)
	if err := r.insertTransactions(txs); err != nil {
		logger.Error("Failed to insertTransactions", "err", err, "blockNumber", event.Block.NumberU64(), "numTxs", len(txs))
		return err
	}
	go r.InvalidateCacheEOAList(updatedEOAs)
	return nil
}

// insertTransactions inserts the given transactions divided into chunkUnit because of the max number of placeholders.
func (r *repository) insertTransactions(txs []*Tx) error {
	chunkUnit := maxPlaceholders / placeholdersPerTxItem
	var chunks []*Tx

	for txs != nil {
		if placeholdersPerTxItem*len(txs) > maxPlaceholders {
			chunks = txs[:chunkUnit]
			txs = txs[chunkUnit:]
		} else {
			chunks = txs
			txs = nil
		}

		if err := r.bulkInsertTransactions(chunks); err != nil {
			return err
		}
	}

	return nil
}

// bulkInsertTransactions inserts the given transactions in multiple rows at once.
func (r *repository) bulkInsertTransactions(txs []*Tx) error {
	var valueStrings []string
	var valueArgs []interface{}

	for _, tx := range txs {
		if _, ok := TxFilteringTypes[types.TxType(tx.TypeInt)]; ok {
			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?,?,?,?,?)")

			valueArgs = append(valueArgs, tx.TransactionId)
			valueArgs = append(valueArgs, tx.FromAddr)
			valueArgs = append(valueArgs, tx.ToAddr)
			valueArgs = append(valueArgs, tx.Value)
			valueArgs = append(valueArgs, tx.TransactionHash)
			valueArgs = append(valueArgs, tx.Status)
			valueArgs = append(valueArgs, tx.Timestamp)
			valueArgs = append(valueArgs, tx.TypeInt)
			valueArgs = append(valueArgs, tx.GasPrice)
			valueArgs = append(valueArgs, tx.GasUsed)
			valueArgs = append(valueArgs, tx.FeePayer)
			valueArgs = append(valueArgs, tx.FeeRatio)
			valueArgs = append(valueArgs, tx.Internal)
		}
	}

	if len(valueStrings) == 0 {
		return nil
	}

	rawQuery := `
			INSERT INTO klay_transfers(transactionId, fromAddr, toAddr, value, transactionHash, status, timestamp, typeInt, gasPrice, gasUsed, feePayer, feeRatio, internal)
			VALUES %s
			ON DUPLICATE KEY
			UPDATE transactionId=transactionId`
	query := fmt.Sprintf(rawQuery, strings.Join(valueStrings, ","))

	if _, err := r.db.DB().Exec(query, valueArgs...); err != nil {
		return err
	}

	return nil
}
