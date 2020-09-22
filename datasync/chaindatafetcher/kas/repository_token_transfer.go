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
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

var tokenTransferEventHash = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

// splitToWords divides log data to the words.
func splitToWords(data []byte) ([]common.Hash, error) {
	if len(data)%common.HashLength != 0 {
		return nil, fmt.Errorf("data length is not valid. want: %v, actual: %v", common.HashLength, len(data))
	}
	var words []common.Hash
	for i := 0; i < len(data); i += common.HashLength {
		words = append(words, common.BytesToHash(data[i:i+common.HashLength]))
	}
	return words, nil
}

// wordToAddress trims input word to get address field only.
func wordToAddress(word common.Hash) common.Address {
	return common.BytesToAddress(word[common.HashLength-common.AddressLength:])
}

// transformLogsToTokenTransfers converts the given event into Klaytn Compatible Token transfers.
func transformLogsToTokenTransfers(event blockchain.ChainEvent) ([]*KCTTransfer, map[common.Address]struct{}, error) {
	timestamp := event.Block.Time().Int64()
	var kctTransfers []*KCTTransfer
	mergedUpdatedEOAs := make(map[common.Address]struct{})
	for _, log := range event.Logs {
		if len(log.Topics) > 0 && log.Topics[0] == tokenTransferEventHash {
			transfer, updatedEOAs, err := transformLogToTokenTransfer(log)
			if err != nil {
				return nil, nil, err
			}
			transfer.Timestamp = timestamp
			kctTransfers = append(kctTransfers, transfer)
			for key := range updatedEOAs {
				mergedUpdatedEOAs[key] = struct{}{}
			}
		}
	}

	return kctTransfers, mergedUpdatedEOAs, nil
}

// transformLogToTokenTransfer converts the given log to Klaytn Compatible Token transfer.
func transformLogToTokenTransfer(log *types.Log) (*KCTTransfer, map[common.Address]struct{}, error) {
	// in case of token transfer,
	// case 1:
	//   log.LogTopics[0] = token transfer event hash
	//   log.LogData = concat(fromAddress, toAddress, value)
	// case 2:
	//   log.LogTopics[0] = token transfer event hash
	//   log.LogTopics[1] = fromAddress
	//   log.LogTopics[2] = toAddresss
	//   log.LogData = value
	words, err := splitToWords(log.Data)
	if err != nil {
		return nil, nil, err
	}
	data := append(log.Topics, words...)
	from := wordToAddress(data[1])
	to := wordToAddress(data[2])
	value := new(big.Int).SetBytes(data[3].Bytes())

	txLogId := int64(log.BlockNumber)*maxTxCountPerBlock*maxTxLogCountPerTx + int64(log.TxIndex)*maxTxLogCountPerTx + int64(log.Index)
	updatedEOAs := make(map[common.Address]struct{})
	updatedEOAs[from] = struct{}{}
	updatedEOAs[to] = struct{}{}

	return &KCTTransfer{
		ContractAddress:  log.Address.Bytes(),
		From:             from.Bytes(),
		To:               to.Bytes(),
		TransactionLogId: txLogId,
		Value:            "0x" + value.Text(16),
		TransactionHash:  log.TxHash.Bytes(),
	}, updatedEOAs, nil
}

// InsertTokenTransfers inserts token transfers in the given chain event into KAS database.
// The token transfers are divided into chunkUnit because of max number of place holders.
func (r *repository) InsertTokenTransfers(event blockchain.ChainEvent) error {
	tokenTransfers, updatedEOAs, err := transformLogsToTokenTransfers(event)
	if err != nil {
		return err
	}

	chunkUnit := maxPlaceholders / placeholdersPerKCTTransferItem
	var chunks []*KCTTransfer

	for tokenTransfers != nil {
		if placeholdersPerKCTTransferItem*len(tokenTransfers) > maxPlaceholders {
			chunks = tokenTransfers[:chunkUnit]
			tokenTransfers = tokenTransfers[chunkUnit:]
		} else {
			chunks = tokenTransfers
			tokenTransfers = nil
		}

		if err := r.bulkInsertTokenTransfers(chunks); err != nil {
			logger.Error("Failed to insertTokenTransfers", "err", err, "numTokenTransfers", len(chunks))
			return err
		}
	}

	go r.InvalidateCacheEOAList(updatedEOAs)
	return nil
}

// bulkInsertTokenTransfers inserts the given token transfers in multiple rows at once.
func (r *repository) bulkInsertTokenTransfers(tokenTransfers []*KCTTransfer) error {
	if len(tokenTransfers) == 0 {
		logger.Debug("the token transfer list is empty")
		return nil
	}
	var valueStrings []string
	var valueArgs []interface{}

	for _, transfer := range tokenTransfers {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?)")

		valueArgs = append(valueArgs, transfer.TransactionLogId)
		valueArgs = append(valueArgs, transfer.From)
		valueArgs = append(valueArgs, transfer.To)
		valueArgs = append(valueArgs, transfer.Value)
		valueArgs = append(valueArgs, transfer.ContractAddress)
		valueArgs = append(valueArgs, transfer.TransactionHash)
		valueArgs = append(valueArgs, transfer.Timestamp)
	}

	rawQuery := `
			INSERT INTO kct_transfers(transactionLogId, fromAddr, toAddr, value, contractAddress, transactionHash, timestamp)
			VALUES %s
			ON DUPLICATE KEY
			UPDATE transactionLogId=transactionLogId`
	query := fmt.Sprintf(rawQuery, strings.Join(valueStrings, ","))

	if _, err := r.db.DB().Exec(query, valueArgs...); err != nil {
		return err
	}
	return nil
}
