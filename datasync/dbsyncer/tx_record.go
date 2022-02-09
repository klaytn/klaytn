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

package dbsyncer

import (
	"strings"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

type SummaryArguments struct {
	from            string
	txHash          string
	to              string
	contractAddress string

	tx *types.Transaction
}

type TxMapArguments struct {
	senderTxHash string
	txHash       string
}

type TxReceiptArgument struct {
	tx      *types.Transaction
	receipt *types.Receipt
}

type BulkInsertQuery struct {
	parameters  string
	vals        []interface{}
	blockNumber uint64
	insertCount int
}

func MakeTxDBRow(block *types.Block, txKey uint64, tx *types.Transaction, receipt *types.Receipt) (string, []interface{}, TxMapArguments, SummaryArguments, error) {

	cols := ""
	vals := []interface{}{}

	blockHash := block.Hash().Hex()
	blockNumber := block.Number().Uint64()

	contractAddress := ""
	if (receipt.ContractAddress != common.Address{}) {
		contractAddress = strings.ToLower(receipt.ContractAddress.Hex())
	}

	from := ""
	if tx.IsEthereumTransaction() {
		signer := types.LatestSignerForChainID(tx.ChainId())
		addr, err := types.Sender(signer, tx)
		if err != nil {
			logger.Error("fail to tx.From", "err", err)
			return "", []interface{}{}, TxMapArguments{}, SummaryArguments{}, err
		}
		from = strings.ToLower(addr.Hex())
	} else {
		addr, err := tx.From()
		if err != nil {
			logger.Error("fail to tx.From", "err", err)
			return "", []interface{}{}, TxMapArguments{}, SummaryArguments{}, err
		}
		from = strings.ToLower(addr.Hex())
	}

	gas := tx.Gas()
	txGasPrice := tx.GasPrice().String()
	txGasUsed := receipt.GasUsed

	input := hexutil.Bytes(tx.Data()).String()
	if input == "0x" {
		input = ""
	}

	nonce := tx.Nonce()
	status := receipt.Status
	to := "" // '' means that address doesn't exist
	if tx.To() != nil {
		to = strings.ToLower(tx.To().Hex())
	}
	timestamp := block.Time().Uint64()
	txHash := tx.Hash().Hex()
	txtype := tx.Type().String()
	value := tx.Value().String()

	feePayer := ""
	feeRatio := uint8(0)
	if tx.IsFeeDelegatedTransaction() {
		feeAddr, err := tx.FeePayer()
		if err != nil {
			logger.Error("fail to tx.FeePayer", "err", err)
			return "", []interface{}{}, TxMapArguments{}, SummaryArguments{}, err
		}
		ratio, _ := tx.FeeRatio()
		// ok is false in TxTypeFeeDelegatedValueTransfer
		// ok is true in TxTypeFeeDelegatedValueTransferWithRatio

		feePayer = strings.ToLower(feeAddr.Hex())
		feeRatio = uint8(ratio)
	}

	senderHash := ""
	if senderTxHash, ok := tx.SenderTxHash(); ok {
		senderHash = senderTxHash.Hex()
	}

	cols = "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	vals = append(vals, txKey, blockHash, blockNumber, contractAddress, from, gas, txGasPrice, txGasUsed, input, nonce, status, to,
		timestamp, txHash, txtype, value, feePayer, feeRatio, senderHash)

	TxMapArg := TxMapArguments{senderHash, txHash}
	summaryArg := SummaryArguments{from, txHash, to, contractAddress, tx}

	return cols, vals, TxMapArg, summaryArg, nil
}

func MakeSummaryDBRow(sa SummaryArguments) (cols string, vals []interface{}, count int, err error) {
	// insert account summary for creation and deploy
	if !sa.tx.IsEthereumTransaction() {
		if sa.tx.Type().IsAccountCreation() {
			accountType := 0
			creator := sa.from
			createdTx := sa.txHash

			hra := true
			//TODO-Klaytn need to use humanreable field
			internalTx, ok := sa.tx.GetTxInternalData().(*types.TxInternalDataAccountCreation)
			if !ok {
				logger.Error("fail to convert TxInternalDataAccountCreation", "txhash", sa.tx.Hash().Hex())
				hra = false
			} else {
				hra = internalTx.HumanReadable
			}

			cols = "(?,?,?,?,?)"
			vals = append(vals, sa.to, accountType, creator, createdTx, hra)

			count = 1

		} else if sa.tx.Type().IsContractDeploy() {
			accountType := 1
			creator := sa.from
			createdTx := sa.txHash

			//TODO-Klaytn need to use humanreable field
			hra := false

			switch internalTx := sa.tx.GetTxInternalData().(type) {
			case *types.TxInternalDataSmartContractDeploy:
				hra = internalTx.HumanReadable
			case *types.TxInternalDataFeeDelegatedSmartContractDeploy:
				hra = internalTx.HumanReadable
			case *types.TxInternalDataFeeDelegatedSmartContractDeployWithRatio:
				hra = internalTx.HumanReadable
			default:
				logger.Error("fail to convert ", "type", internalTx, "txhash", sa.tx.Hash().Hex())
			}

			cols = "(?,?,?,?,?)"
			vals = append(vals, sa.contractAddress, accountType, creator, createdTx, hra)

			count = 1
		}
	} else if sa.contractAddress != "" {
		// AccountCreation for LegacyTransaction
		accountType := 1
		creator := sa.from
		createdTx := sa.txHash

		cols = "(?,?,?,?,?)"
		vals = append(vals, sa.contractAddress, accountType, creator, createdTx, false)

		count = 1
	}

	return cols, vals, count, nil
}

func MakeTxMappingRow(tm TxMapArguments) (cols string, vals []interface{}, count int, err error) {

	if tm.senderTxHash != "" {
		cols = "(?,?)"
		vals = append(vals, tm.senderTxHash, tm.txHash)

		count = 1
	}

	return cols, vals, count, nil
}
