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
	"reflect"
	"strings"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
)

const selfDestructType = "SELFDESTRUCT"

var emptyTraceResult = &vm.InternalTxTrace{
	Value: "0x0",
	Calls: []*vm.InternalTxTrace{},
}

func isEmptyTraceResult(trace *vm.InternalTxTrace) bool {
	return reflect.DeepEqual(trace, emptyTraceResult)
}

// getEntryTx returns a entry transaction which may call internal transactions.
func getEntryTx(block *types.Block, txIdx int, tx *types.Transaction) *Tx {
	head := block.Header()
	txId := head.Number.Int64()*maxTxCountPerBlock*maxTxLogCountPerTx + int64(txIdx)*maxInternalTxCountPerTx
	return &Tx{
		TransactionId:   txId,
		TransactionHash: tx.Hash().Bytes(),
		Status:          int(types.ReceiptStatusSuccessful),
		Timestamp:       block.Time().Int64(),
		TypeInt:         int(tx.Type()),
		Internal:        true,
	}
}

// transformToInternalTx converts the result of call tracer into the internal transaction list according to the KAS database scheme.
func transformToInternalTx(trace *vm.InternalTxTrace, offset *int64, entryTx *Tx, isFirstCall bool) ([]*Tx, error) {
	if trace.Error != nil {
		// ignore the internal transaction if it has an error field
		return nil, nil
	}

	if trace.Type == "" {
		return nil, noOpcodeError
	} else if trace.Type == selfDestructType {
		// TODO-ChainDataFetcher currently, skip it when self-destruct is encountered.
		return nil, nil
	}

	if trace.From == nil {
		return nil, noFromFieldError
	}

	if trace.To == nil {
		return nil, noToFieldError
	}

	var txs []*Tx
	if !isFirstCall && trace.Value != "" && trace.Value != "0x0" { // filter the internal tx if the value is 0
		*offset++ // adding 1 to calculate the transaction id in the order of internal transactions
		txs = append(txs, &Tx{
			TransactionId:   entryTx.TransactionId + *offset,
			FromAddr:        trace.From.Bytes(),
			ToAddr:          trace.To.Bytes(),
			Value:           trace.Value,
			TransactionHash: entryTx.TransactionHash,
			Status:          entryTx.Status,
			Timestamp:       entryTx.Timestamp,
			TypeInt:         entryTx.TypeInt,
			Internal:        true,
		})
	}

	for _, call := range trace.Calls {
		nestedCalls, err := transformToInternalTx(call, offset, entryTx, false)
		if err != nil {
			return nil, err
		}
		txs = append(txs, nestedCalls...)
	}

	return txs, nil
}

// transformToRevertedTx converts the result of call tracer into the reverted transaction information according to the KAS database scheme.
func transformToRevertedTx(trace *vm.InternalTxTrace, block *types.Block, entryTx *types.Transaction) (*RevertedTx, error) {
	return &RevertedTx{
		TransactionHash: entryTx.Hash().Bytes(),
		BlockNumber:     block.Number().Int64(),
		RevertMessage:   trace.Reverted.Message,
		ContractAddress: trace.Reverted.Contract.Bytes(),
		Timestamp:       block.Time().Int64(),
	}, nil
}

// transformToTraceResults converts the chain event to internal transactions as well as reverted transactions.
func transformToTraceResults(event blockchain.ChainEvent) ([]*Tx, []*RevertedTx, error) {
	var (
		internalTxs []*Tx
		revertedTxs []*RevertedTx
	)
	for txIdx, trace := range event.InternalTxTraces {
		if isEmptyTraceResult(trace) {
			continue
		}

		tx := event.Block.Transactions()[txIdx]
		receipt := event.Receipts[txIdx]

		entryTx := getEntryTx(event.Block, txIdx, tx)
		offset := int64(0)

		// transforms the result into internal transaction which is associated with KLAY transfer recursively.
		if receipt.Status == types.ReceiptStatusSuccessful {
			internalTx, err := transformToInternalTx(trace, &offset, entryTx, true)
			if err != nil {
				logger.Error("Failed to transform tracing result into internal tx", "err", err, "txHash", common.BytesToHash(entryTx.TransactionHash).String())
				return nil, nil, err
			}
			internalTxs = append(internalTxs, internalTx...)
		}

		// transforms the result into an evm reverted transaction.
		if receipt.Status == types.ReceiptStatusErrExecutionReverted {
			revertedTx, err := transformToRevertedTx(trace, event.Block, tx)
			if err != nil {
				logger.Error("Failed to transform tracing result into reverted tx", "err", err, "txHash", common.BytesToHash(entryTx.TransactionHash).String())
				return nil, nil, err
			}
			revertedTxs = append(revertedTxs, revertedTx)
		}
	}
	return internalTxs, revertedTxs, nil
}

// InsertTraceResults inserts internal and reverted transactions in the given chain event into KAS database.
func (r *repository) InsertTraceResults(event blockchain.ChainEvent) error {
	internalTxs, revertedTxs, err := transformToTraceResults(event)
	if err != nil {
		logger.Error("Failed to transform the given event to tracing results", "err", err, "blockNumber", event.Block.NumberU64())
		return err
	}

	if err := r.insertTransactions(internalTxs); err != nil {
		logger.Error("Failed to insert internal transactions", "err", err, "blockNumber", event.Block.NumberU64(), "numInternalTxs", len(internalTxs))
		return err
	}

	if err := r.insertRevertedTransactions(revertedTxs); err != nil {
		logger.Error("Failed to insert reverted transactions", "err", err, "blockNumber", event.Block.NumberU64(), "numRevertedTxs", len(revertedTxs))
		return err
	}
	return nil
}

// insertRevertedTransactions inserts the given reverted transactions divided into chunkUnit because of the max number of placeholders.
func (r *repository) insertRevertedTransactions(revertedTxs []*RevertedTx) error {
	chunkUnit := maxPlaceholders / placeholdersPerRevertedTxItem
	var chunks []*RevertedTx

	for revertedTxs != nil {
		if placeholdersPerRevertedTxItem*len(revertedTxs) > maxPlaceholders {
			chunks = revertedTxs[:chunkUnit]
			revertedTxs = revertedTxs[chunkUnit:]
		} else {
			chunks = revertedTxs
			revertedTxs = nil
		}

		if err := r.bulkInsertRevertedTransactions(chunks); err != nil {
			return err
		}
	}

	return nil
}

// bulkInsertRevertedTransactions inserts the given reverted transactions in multiple rows at once.
func (r *repository) bulkInsertRevertedTransactions(revertedTxs []*RevertedTx) error {
	if len(revertedTxs) == 0 {
		return nil
	}
	var valueStrings []string
	var valueArgs []interface{}

	for _, revertedTx := range revertedTxs {
		valueStrings = append(valueStrings, "(?,?,?,?,?)")

		valueArgs = append(valueArgs, revertedTx.TransactionHash)
		valueArgs = append(valueArgs, revertedTx.BlockNumber)
		valueArgs = append(valueArgs, revertedTx.ContractAddress)
		valueArgs = append(valueArgs, revertedTx.RevertMessage)
		valueArgs = append(valueArgs, revertedTx.Timestamp)
	}

	rawQuery := `
			INSERT INTO reverted_transactions(transactionHash, blockNumber, contractAddress, revertMessage, timestamp)
			VALUES %s
			ON DUPLICATE KEY
			UPDATE transactionHash=transactionHash`
	query := fmt.Sprintf(rawQuery, strings.Join(valueStrings, ","))

	if _, err := r.db.DB().Exec(query, valueArgs...); err != nil {
		return err
	}
	return nil
}
